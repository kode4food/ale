package markdown

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// Header stores the fields of a Markdown document header (at least the ones we
// care about for processing)
type Header struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Usage       string   `yaml:"usage"`
	Names       []string `yaml:"names"`
	Tags        []string `yaml:"tags"`
}

// Parse parses the kind of markdown document that might be processed by a
// static site generator. It will parse any prologue parameters into the
// resulting object and return the remaining content as individual lines
func Parse(doc string) (*Header, []string) {
	obj, rest := parseDocument(doc)
	return obj, skipEmptyLines(rest)
}

// ParseHeader parses the header of a markdown document that might be processed
// by a static site generator, returning the prologue parameters as a Header
func ParseHeader(doc string) *Header {
	res, _ := parseDocument(doc)
	return res
}

func skipEmptyLines(lines []string) []string {
	var first = 0
	for first < len(lines) && strings.TrimSpace(lines[first]) == "" {
		first++
	}
	return lines[first:]
}

func parseDocument(doc string) (*Header, []string) {
	lines := strings.Split(doc, "\n")
	if strings.TrimSpace(lines[0]) != "---" {
		return nil, lines
	}

	lines = lines[1:]
	var rest = 1
	for i, l := range lines {
		if strings.TrimSpace(l) == "---" {
			rest = i + 1
			break
		}
	}

	head := lines[:rest-1]
	y := strings.Join(head, "\n")
	return parseHeader(y), lines[rest:]
}

func parseHeader(b string) *Header {
	m := new(Header)
	if err := yaml.Unmarshal([]byte(b), &m); err != nil {
		panic(err)
	}
	return m
}
