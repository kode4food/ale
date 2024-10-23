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

// Parse parses the kind of Markdown document that might be processed by a
// static site generator. It will parse any prologue parameters into the
// resulting object and return the remaining content as individual lines
func Parse(doc string) (*Header, []string, error) {
	obj, rest, err := parseDocument(doc)
	if err != nil {
		return nil, nil, err
	}
	return obj, skipEmptyLines(rest), nil
}

// ParseHeader parses the header of a Markdown document that might be processed
// by a static site generator, returning the prologue parameters as a Header
func ParseHeader(doc string) (*Header, error) {
	res, _, err := parseDocument(doc)
	return res, err
}

func skipEmptyLines(lines []string) []string {
	var first = 0
	for first < len(lines) && strings.TrimSpace(lines[first]) == "" {
		first++
	}
	return lines[first:]
}

func parseDocument(doc string) (*Header, []string, error) {
	lines := strings.Split(doc, "\n")
	if strings.TrimSpace(lines[0]) != "---" {
		return nil, lines, nil
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
	res, err := parseHeader(y)
	if err != nil {
		return nil, nil, err
	}
	return res, lines[rest:], nil
}

func parseHeader(b string) (*Header, error) {
	m := new(Header)
	if err := yaml.Unmarshal([]byte(b), &m); err != nil {
		return nil, err
	}
	return m, nil
}
