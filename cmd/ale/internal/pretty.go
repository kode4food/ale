package internal

import (
	"strings"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/cmd/ale/internal/console"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/sequence"
)

const (
	defaultIndent = "  "
)

func PrettyPrint(v ale.Value) string {
	return PrettyPrintAt(v, 0)
}

func PrettyPrintAt(v ale.Value, offset int) string {
	return prettyPrint(v, offset, 0)
}

func prettyPrint(v ale.Value, offset, depth int) string {
	switch val := v.(type) {
	case *data.List:
		return prettyList(offset, val, depth)
	case data.Vector:
		return prettyVector(offset, val, depth)
	case *data.Object:
		return prettyObject(offset, val, depth)
	case *data.Cons:
		return prettyCons(offset, val, depth)
	default:
		return data.ToQuotedString(val)
	}
}

func prettyList(offset int, list *data.List, depth int) string {
	if list == nil {
		return "()"
	}

	elements := basics.Map(sequence.ToVector(list), func(v ale.Value) string {
		return prettyPrint(v, offset, depth+1)
	})
	return formatSequence("(", ")", elements, offset, depth)
}

func prettyVector(offset int, vec data.Vector, depth int) string {
	elements := basics.Map(vec, func(v ale.Value) string {
		return prettyPrint(v, offset, depth+1)
	})
	return formatSequence("[", "]", elements, offset, depth)
}

func prettyObject(offset int, obj *data.Object, depth int) string {
	if obj == nil {
		return "{}"
	}

	pairs := obj.Pairs()
	formattedPairs := basics.Map(pairs, func(pair data.Pair) [2]string {
		key := prettyPrint(pair.Car(), offset, depth+1)
		value := prettyPrint(pair.Cdr(), offset, depth+1)
		return [2]string{key, value}
	})

	maxKeyWidth := findMaxKeyWidth(formattedPairs)
	elements := basics.Map(formattedPairs, func(fp [2]string) string {
		return alignPair(fp[0], fp[1], maxKeyWidth)
	})

	return formatObject("{", "}", elements, offset, depth)
}

func findMaxKeyWidth(formattedPairs [][2]string) int {
	maxWidth := 0
	for _, fp := range formattedPairs {
		key := fp[0]
		if !strings.Contains(key, "\n") && len(key) > maxWidth {
			maxWidth = len(key)
		}
	}
	return maxWidth
}

func alignPair(key, value string, maxKeyWidth int) string {
	if strings.Contains(key, "\n") {
		return key + " " + value
	}
	padding := strings.Repeat(" ", maxKeyWidth-len(key))
	return key + padding + " " + value
}

func prettyCons(offset int, cons *data.Cons, depth int) string {
	car := prettyPrint(cons.Car(), offset, depth+1)
	
	// Calculate the new offset for the cdr to align with the opening brace
	prefix := "(" + car + " . "
	cdrOffset := offset + len(prefix)
	cdr := prettyPrint(cons.Cdr(), cdrOffset, depth+1)
	
	return prefix + cdr + ")"
}

func makeIndent(offset, depth int) (base string, prompt string) {
	base = strings.Repeat(defaultIndent, depth+1)
	if depth == 0 {
		prompt = strings.Repeat(" ", offset)
	}
	return
}

func lineWidth(offset, depth int) int {
	screenWidth := console.GetScreenWidth()
	if depth == 0 {
		return screenWidth - offset
	}
	return screenWidth
}

func contentWidth(offset, depth int) int {
	indentWidth := len(strings.Repeat(defaultIndent, depth+1))
	return lineWidth(offset, depth) - indentWidth
}

func formatObject(
	start, end string, elements []string, offset, depth int,
) string {
	if len(elements) == 0 {
		return start + end
	}
	if len(elements) == 1 {
		return start + elements[0] + end
	}

	baseIndent, promptOffset := makeIndent(offset, depth)
	var buf strings.Builder
	buf.WriteString(start)

	for _, elem := range elements {
		buf.WriteString("\n" + promptOffset + baseIndent + elem)
	}

	buf.WriteString(
		"\n" + promptOffset + strings.Repeat(defaultIndent, depth) + end,
	)
	return buf.String()
}

func formatSequence(
	start, end string, elements []string, offset, depth int,
) string {
	if len(elements) == 0 {
		return start + end
	}

	inline := start + strings.Join(elements, " ") + end
	if len(inline) <= lineWidth(offset, depth) && !strings.Contains(inline, "\n") {
		return inline
	}

	return formatMultiline(start, end, elements, offset, depth)
}

func formatMultiline(
	start, end string, elements []string, offset, depth int,
) string {
	baseIndent, promptOffset := makeIndent(offset, depth)
	contentWidth := contentWidth(offset, depth)

	var buf strings.Builder
	buf.WriteString(start)

	for i := 0; i < len(elements); {
		buf.WriteString("\n" + promptOffset + baseIndent)
		i = writeLineElements(&buf, elements, i, contentWidth)
	}

	buf.WriteString(
		"\n" + promptOffset + strings.Repeat(defaultIndent, depth) + end,
	)
	return buf.String()
}

func writeLineElements(
	buf *strings.Builder, elements []string, startIdx, availableWidth int,
) int {
	lineWidth := 0

	for i := startIdx; i < len(elements); i++ {
		elem := elements[i]
		needsSpace := i > startIdx
		elemWidth := len(elem) + boolToInt(needsSpace)

		if lineWidth+elemWidth > availableWidth && i > startIdx {
			return i
		}

		if needsSpace {
			buf.WriteString(" ")
		}
		buf.WriteString(elem)
		lineWidth += elemWidth
	}

	return len(elements)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
