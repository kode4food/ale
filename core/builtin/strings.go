package builtin

import (
	"strings"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/sequence"
)

const emptyString = data.String("")

// Str converts the provided arguments to an undelimited string
var Str = data.MakeProcedure(func(args ...ale.Value) ale.Value {
	v := data.Vector(args)
	return sequence.ToString(v)
})

// ReaderStr converts the provided arguments to a delimited string
var ReaderStr = data.MakeProcedure(func(args ...ale.Value) ale.Value {
	if len(args) == 0 {
		return emptyString
	}

	var b strings.Builder
	b.WriteString(data.ToQuotedString(args[0]))
	for _, f := range args[1:] {
		b.WriteString(lang.Space)
		b.WriteString(data.ToQuotedString(f))
	}
	return data.String(b.String())
})
