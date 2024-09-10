package builtin

import (
	"strings"

	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
)

const emptyString = data.String("")

// Str converts the provided arguments to an undelimited string
var Str = data.MakeProcedure(func(args ...data.Value) data.Value {
	v := data.Vector(args)
	return sequence.ToString(v)
})

// ReaderStr converts the provided arguments to a delimited string
var ReaderStr = data.MakeProcedure(func(args ...data.Value) data.Value {
	if len(args) == 0 {
		return emptyString
	}

	var b strings.Builder
	b.WriteString(data.ToQuotedString(args[0]))
	for _, f := range args[1:] {
		b.WriteString(" ")
		b.WriteString(data.ToQuotedString(f))
	}
	return data.String(b.String())
})
