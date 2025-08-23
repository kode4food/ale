package internal

import (
	"cmp"
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

func prettyList(l *data.List, off int) string {
	if l == nil {
		return lang.ListStart + lang.ListEnd
	}
	v := sequence.ToVector(l)
	return prettySeq(v, lang.ListStart, lang.ListEnd, off)
}

func prettyVector(v data.Vector, off int) string {
	return prettySeq(v, lang.VectorStart, lang.VectorEnd, off)
}

func prettySeq(elems data.Vector, start, end string, off int) string {
	if len(elems) == 0 {
		return start + end
	}

	elementOffset := off + indentSize
	formatted := basics.Map(elems, func(v ale.Value) string {
		return PrettyPrintAt(v, elementOffset)
	})
	return formatSeq(start, end, formatted, off)
}

func prettyObject(o *data.Object, off int) string {
	if o == nil {
		return lang.ObjectStart + lang.ObjectEnd
	}

	pairs := o.Pairs()
	elementOffset := off + indentSize
	formattedPairs := basics.Map(pairs, func(pair data.Pair) [2]string {
		key := PrettyPrintAt(pair.Car(), elementOffset)
		value := PrettyPrintAt(pair.Cdr(), elementOffset)
		return [2]string{key, value}
	})

	sorted := basics.SortedFunc(formattedPairs, func(l, r [2]string) int {
		return cmp.Compare(l[0], r[0])
	})

	maxWidth := maxKeyWidth(sorted)
	elems := basics.Map(sorted, func(fp [2]string) string {
		return alignPair(fp[0], fp[1], maxWidth)
	})

	return formatObject(lang.ObjectStart, lang.ObjectEnd, elems, off)
}

func maxKeyWidth(pairs [][2]string) int {
	maxWidth := 0
	for _, fp := range pairs {
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
	return lang.ListStart +
		PrettyPrintAt(cons.Car(), off) +
		SP + lang.Dot + SP +
		PrettyPrintAt(cons.Cdr(), off) +
		lang.ListEnd
}

func makeIndent(off int) string {
	return strings.Repeat(SP, off)
}

func lineWidth(off int) int {
	screenWidth := console.GetScreenWidth()
	return screenWidth - off
}

func formatObject(start, end string, elems []string, off int) string {
	switch len(elems) {
	case 0:
		return start + end
	case 1:
		return start + elems[0] + end
	default:
		indent := makeIndent(off)
		elemIndent := makeIndent(off + indentSize)
		var buf strings.Builder
		buf.WriteString(start)

		for _, elem := range elems {
			buf.WriteString(NL + elemIndent + elem)
		}

		buf.WriteString(NL + indent + end)
		return buf.String()
	}
}

func formatSeq(start, end string, elems []string, off int) string {
	switch len(elems) {
	case 0:
		return start + end
	default:
		inline := start + strings.Join(elems, SP) + end
		if len(inline) <= lineWidth(off) && !strings.Contains(inline, NL) {
			return inline
		}

		return formatLines(start, end, elems, off)
	}
}

func formatLines(start, end string, elems []string, off int) string {
	indent := makeIndent(off)
	elemIndent := makeIndent(off + indentSize)
	availWidth := lineWidth(off + indentSize)

	var buf strings.Builder
	buf.WriteString(start)

	for i := 0; i < len(elems); {
		buf.WriteString(NL + elemIndent)
		i = writeElems(&buf, elems, i, availWidth)
	}

	buf.WriteString(NL + indent + end)
	return buf.String()
}

func writeElems(
	buf *strings.Builder, elems []string, startIdx, availWidth int,
) int {
	lineWidth := 0

	for i := startIdx; i < len(elems); i++ {
		elem := elems[i]
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

	return len(elems)
}
