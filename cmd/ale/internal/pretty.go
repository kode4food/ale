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

// PrettyPrintAt pretty formats a Value at the given indentation offset
func PrettyPrintAt(v ale.Value, offset int) string {
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

func prettyList(list *data.List, off int) string {
	if list == nil {
		return lang.ListStart + lang.ListEnd
	}
	v := sequence.ToVector(list)
	return prettySeq(v, lang.ListStart, lang.ListEnd, off)
}

func prettyVector(vec data.Vector, off int) string {
	return prettySeq(vec, lang.VectorStart, lang.VectorEnd, off)
}

func prettySeq(elements data.Vector, start, end string, off int) string {
	if len(elements) == 0 {
		return start + end
	}

	elementOffset := off + indentSize
	formatted := basics.Map(elements, func(v ale.Value) string {
		return PrettyPrintAt(v, elementOffset)
	})
	return formatSeq(start, end, formatted, off)
}

func prettyObject(obj *data.Object, off int) string {
	if obj == nil {
		return lang.ObjectStart + lang.ObjectEnd
	}

	pairs := obj.Pairs()
	elementOffset := off + indentSize
	formattedPairs := basics.Map(pairs, func(pair data.Pair) [2]string {
		key := PrettyPrintAt(pair.Car(), elementOffset)
		value := PrettyPrintAt(pair.Cdr(), elementOffset)
		return [2]string{key, value}
	})

	sortedPairs := basics.SortedFunc(formattedPairs, func(l, r [2]string) int {
		if l[0] < r[0] {
			return -1
		} else if l[0] > r[0] {
			return 1
		}
		return 0
	})

	maxKeyWidth := maxKeyWidth(sortedPairs)
	elements := basics.Map(sortedPairs, func(fp [2]string) string {
		return alignPair(fp[0], fp[1], maxKeyWidth)
	})

	return formatObject(lang.ObjectStart, lang.ObjectEnd, elements, off)
}

func maxKeyWidth(formattedPairs [][2]string) int {
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

func prettyCons(cons *data.Cons, off int) string {
	car := PrettyPrintAt(cons.Car(), off)
	cdr := PrettyPrintAt(cons.Cdr(), off)
	return lang.ListStart + car + SP + lang.Dot + SP + cdr + lang.ListEnd
}

func makeIndent(off int) string {
	return strings.Repeat(SP, off)
}

func lineWidth(off int) int {
	screenWidth := console.GetScreenWidth()
	return screenWidth - off
}

func formatObject(start, end string, elements []string, off int) string {
	if len(elements) == 0 {
		return start + end
	}
	if len(elements) == 1 {
		return start + elements[0] + end
	}

	indent := makeIndent(off)
	elementIndent := makeIndent(off + indentSize)
	var buf strings.Builder
	buf.WriteString(start)

	for _, elem := range elements {
		buf.WriteString(NL + elementIndent + elem)
	}

	buf.WriteString(NL + indent + end)
	return buf.String()
}

func formatSeq(start, end string, elements []string, off int) string {
	if len(elements) == 0 {
		return start + end
	}

	inline := start + strings.Join(elements, SP) + end
	if len(inline) <= lineWidth(off) && !strings.Contains(inline, NL) {
		return inline
	}

	return formatMultiline(start, end, elements, off)
}

func formatMultiline(start, end string, elements []string, off int) string {
	indent := makeIndent(off)
	elemIndent := makeIndent(off + indentSize)
	availWidth := lineWidth(off + indentSize)

	var buf strings.Builder
	buf.WriteString(start)

	for i := 0; i < len(elements); {
		buf.WriteString(NL + elemIndent)
		i = writeElems(&buf, elements, i, availWidth)
	}

	buf.WriteString(NL + indent + end)
	return buf.String()
}

func writeElems(
	buf *strings.Builder, elements []string, startIdx, availWidth int,
) int {
	lineWidth := 0

	for i := startIdx; i < len(elements); i++ {
		elem := elements[i]
		spaceWidth := 0
		if i > startIdx {
			spaceWidth = 1
		}
		elemWidth := len(elem) + spaceWidth

		if lineWidth+elemWidth > availWidth && i > startIdx {
			return i
		}

		if i > startIdx {
			buf.WriteString(SP)
		}
		buf.WriteString(elem)
		lineWidth += elemWidth
	}

	return len(elements)
}
