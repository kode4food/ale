package builtin

import (
	"bytes"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

const emptyString = data.String("")

// Str converts the provided arguments to an undelimited string
var Str = data.Applicative(func(args ...data.Value) data.Value {
	v := data.NewVector(args...)
	return sequence.ToStr(v)
})

// ReaderStr converts the provided arguments to a delimited string
var ReaderStr = data.Applicative(func(args ...data.Value) data.Value {
	if len(args) == 0 {
		return emptyString
	}

	var b bytes.Buffer
	b.WriteString(data.MaybeQuoteString(args[0]))
	for _, f := range args[1:] {
		b.WriteString(" ")
		b.WriteString(data.MaybeQuoteString(f))
	}
	return data.String(b.String())
})

// IsString returns whether the provided value is a string
var IsString = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.String)
	return data.Bool(ok)
}, 1)
