package markdown

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/kode4food/ale/internal/console"
)

// This is *not* a full-featured markdown formatter, or even a compliant one
// for that matter. It only supports the productions that are used by the
// system's documentation strings and will likely not evolve much beyond that

type (
	formatter func(string) string

	patternFormatter struct {
		pattern *regexp.Regexp
		format  formatter
	}
)

const blockPrefix = "```"

var (
	header = regexp.MustCompile(`^#+\s(.*)$`)

	indent  = regexp.MustCompile(`^#+ |^\s\s+`)
	comment = regexp.MustCompile(`;.*(\n|$)`)
	hashes  = regexp.MustCompile(`^#+ `)
	ticks   = regexp.MustCompile("`[^`]*`")
	unders  = regexp.MustCompile(`_[^_]*_`)
	double  = regexp.MustCompile(`[*][*][^*]*[*][*]`)
	stars   = regexp.MustCompile(`[*][^*]*[*]`)

	lineFormatters = map[*regexp.Regexp]formatter{
		regexp.MustCompile(`^#\s.*$`):   formatHeader1,
		regexp.MustCompile(`^##+\s.*$`): formatHeader2,
		regexp.MustCompile(`^\s\s.*$`):  formatIndent,
	}

	docFormatters = []*patternFormatter{
		{comment, formatComment},
		{ticks, trimmedFormatter("`", console.Code)},
		{ticks, trimmedFormatter("`", console.Code)},
		{double, trimmedFormatter(`**`, console.Bold)},
		{stars, trimmedFormatter(`*`, console.Italic)},
		{unders, trimmedFormatter(`_`, console.Result)},
	}

	escaped = regexp.MustCompile(`\\.`)
)

// FormatMarkdown formats a markdown asset for REPL display
func FormatMarkdown(s string) string {
	doc := strings.TrimSpace(s)
	header, lines := Parse(doc)

	var pre []string
	if desc := header.Description; desc != "" {
		pre = append(pre, fmt.Sprintf("# %s", desc))
	}
	if usage := header.Usage; usage != "" {
		pre = append(pre, "")
		pre = append(pre, fmt.Sprintf("Usage: `%s`", usage))
		pre = append(pre, "")
	}

	lines = append(pre, lines...)
	lines = formatCode(lines)
	lines = formatLines(lines)
	res := strings.Join(lines, "\n")
	return formatContent(res)
}

func formatCode(lines []string) []string {
	var res []string
	for i := 0; i < len(lines); i++ {
		if l := lines[i]; !strings.HasPrefix(l, blockPrefix) {
			res = append(res, l)
			continue
		}
		for i = i + 1; i < len(lines); i++ {
			c := lines[i]
			if strings.HasPrefix(c, blockPrefix) {
				break
			}
			res = append(res, fmt.Sprintf("  %s", c))
		}
	}
	return res
}

func formatLines(lines []string) []string {
	var out []string
	for _, l := range wrapLines(lines) {
		res := formatLine(l)
		out = append(out, res)
	}
	return out
}

func formatLine(l string) string {
	var res = l
	for r, f := range lineFormatters {
		res = r.ReplaceAllStringFunc(res, f)
	}
	return res
}

func formatContent(doc string) string {
	var buf bytes.Buffer
	var src = strings.TrimSpace(doc)
	for len(src) > 0 {
		f, sm, ok := firstFormatterMatch(src)
		if !ok {
			buf.WriteString(src)
			src = ""
			continue
		}
		start := sm[0]
		end := sm[1]
		fragment := src[start:end]
		buf.WriteString(src[0:start])
		buf.WriteString(f(fragment))
		src = src[end:]
	}
	return escaped.ReplaceAllStringFunc(buf.String(), func(s string) string {
		return s[1:]
	})
}

func firstFormatterMatch(src string) (formatter, []int, bool) {
	var idx = -1
	var match []int
	var format formatter
	for _, f := range docFormatters {
		if sm := f.pattern.FindStringSubmatchIndex(src); sm != nil {
			pos := sm[0]
			if idx != -1 && pos >= idx {
				continue
			}
			if pos > 0 && src[pos-1] == '\\' {
				continue
			}
			idx = pos
			format = f.format
			match = sm
		}
	}
	return format, match, idx != -1
}

func getFormatWidth() int {
	if res := console.GetScreenWidth(); res != -1 {
		return res - 4
	}
	return 76
}

func wrapLines(s []string) []string {
	var r []string
	w := getFormatWidth()
	for _, e := range s {
		r = append(r, wrapLine(e, w)...)
	}
	return r
}

func wrapLine(s string, w int) []string {
	var r []string
	i, s := lineIndent(s)
	il := strippedLen(i)

	var b bytes.Buffer
	b.WriteString(i)
	for _, e := range strings.Split(s, " ") {
		bl := strippedLen(b.String())
		if bl > il {
			if bl+strippedLen(e) >= w {
				r = append(r, b.String())
				b.Reset()
				b.WriteString(i)
			} else if !isEmptyString(b.String()) {
				b.WriteString(" ")
			}
		}
		b.WriteString(e)
	}
	return append(r, b.String())
}

func isEmptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func lineIndent(s string) (string, string) {
	l := indent.FindString(s)
	return l, s[len(l):]
}

func stripDelimited(s string) string {
	s = ticks.ReplaceAllStringFunc(s, trimmer(1))
	s = double.ReplaceAllStringFunc(s, trimmer(1))
	return hashes.ReplaceAllString(s, "")
}

func strippedLen(s string) int {
	return len(stripDelimited(s))
}

func trimmer(size int) formatter {
	return func(s string) string {
		if len(s) > (size * 2) {
			return s[size : len(s)-size]
		}
		return s[:size]
	}
}

func formatHeader1(s string) string {
	return console.Header1 + stripHashes(s) + console.Reset
}

func formatHeader2(s string) string {
	return console.Header2 + stripHashes(s) + console.Reset
}

func formatIndent(s string) string {
	return console.Code + s + console.Reset
}

func formatComment(s string) string {
	return console.Comment + s + console.Reset
}

func stripHashes(s string) string {
	if sm := header.FindStringSubmatch(s); sm != nil {
		return sm[1]
	}
	return s
}

func trimmedFormatter(delim, prefix string) formatter {
	trim := trimmer(len(delim))
	return func(s string) string {
		t := trim(s)
		if t == delim {
			return t
		}
		return prefix + t + console.Reset
	}
}
