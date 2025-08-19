package internal

import (
	"strings"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/cmd/ale/internal/console"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/sequence"
)

const (
	SP = " "
	NL = "\n"

	indentSize = 2
)

func PrettyPrintAt(v ale.Value, offset int) string {
	return prettyPrint(v, offset)
}

func prettyPrint(v ale.Value, offset int) string {
	switch val := v.(type) {
	case *data.List:
		return prettyList(val, offset)
	case data.Vector:
		return prettyVector(val, offset)
	case *data.Object:
		return prettyObject(val, offset)
	case *data.Cons:
		return prettyCons(val, offset)
	default:
		return data.ToQuotedString(val)
	}
}

func prettyList(list *data.List, offset int) string {
	if list == nil {
		return lang.ListStart + lang.ListEnd
	}

	elementOffset := offset + indentSize
	elements := basics.Map(sequence.ToVector(list), func(v ale.Value) string {
		return prettyPrint(v, elementOffset)
	})
	return formatSequence(lang.ListStart, lang.ListEnd, elements, offset)
}

func prettyVector(vec data.Vector, offset int) string {
	elementOffset := offset + indentSize
	elements := basics.Map(vec, func(v ale.Value) string {
		return prettyPrint(v, elementOffset)
	})
	return formatSequence(lang.VectorStart, lang.VectorEnd, elements, offset)
}

func prettyObject(obj *data.Object, offset int) string {
	if obj == nil {
		return lang.ObjectStart + lang.ObjectEnd
	}

	pairs := obj.Pairs()
	elementOffset := offset + indentSize
	formattedPairs := basics.Map(pairs, func(pair data.Pair) [2]string {
		key := prettyPrint(pair.Car(), elementOffset)
		value := prettyPrint(pair.Cdr(), elementOffset)
		return [2]string{key, value}
	})

	maxKeyWidth := findMaxKeyWidth(formattedPairs)
	elements := basics.Map(formattedPairs, func(fp [2]string) string {
		return alignPair(fp[0], fp[1], maxKeyWidth)
	})

	return formatObject(lang.ObjectStart, lang.ObjectEnd, elements, offset)
}

func findMaxKeyWidth(formattedPairs [][2]string) int {
	maxWidth := 0
	for _, fp := range formattedPairs {
		key := fp[0]
		if !strings.Contains(key, NL) && len(key) > maxWidth {
			maxWidth = len(key)
		}
	}
	return maxWidth
}

func alignPair(key, value string, maxKeyWidth int) string {
	if strings.Contains(key, NL) {
		return key + SP + value
	}
	padding := strings.Repeat(SP, maxKeyWidth-len(key))
	return key + padding + SP + value
}

func prettyCons(cons *data.Cons, offset int) string {
	elementOffset := offset + indentSize
	car := prettyPrint(cons.Car(), elementOffset)
	cdr := prettyPrint(cons.Cdr(), elementOffset)
	return lang.ListStart + car + SP + lang.Dot + SP + cdr + lang.ListEnd
}

func makeIndent(offset int) string {
	return strings.Repeat(SP, offset)
}

func lineWidth(offset int) int {
	screenWidth := console.GetScreenWidth()
	return screenWidth - offset
}

func formatObject(start, end string, elements []string, offset int) string {
	if len(elements) == 0 {
		return start + end
	}
	if len(elements) == 1 {
		return start + elements[0] + end
	}

	indent := makeIndent(offset)
	elementIndent := makeIndent(offset + indentSize)
	var buf strings.Builder
	buf.WriteString(start)

	for _, elem := range elements {
		buf.WriteString(NL + elementIndent + elem)
	}

	buf.WriteString(NL + indent + end)
	return buf.String()
}

func formatSequence(start, end string, elements []string, offset int) string {
	if len(elements) == 0 {
		return start + end
	}

	inline := start + strings.Join(elements, SP) + end
	if len(inline) <= lineWidth(offset) && !strings.Contains(inline, NL) {
		return inline
	}

	return formatMultiline(start, end, elements, offset)
}

func formatMultiline(start, end string, elements []string, offset int) string {
	indent := makeIndent(offset)
	elementIndent := makeIndent(offset + indentSize)
	availableWidth := lineWidth(offset + indentSize)

	var buf strings.Builder
	buf.WriteString(start)

	for i := 0; i < len(elements); {
		buf.WriteString(NL + elementIndent)
		i = writeElements(&buf, elements, i, availableWidth)
	}

	buf.WriteString(NL + indent + end)
	return buf.String()
}

func writeElements(
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
			buf.WriteString(SP)
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
