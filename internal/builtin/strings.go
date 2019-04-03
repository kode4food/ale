package builtin

import (
	"bytes"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/stdlib"
)

const emptyString = api.String("")

// Str converts the provided arguments to an undelimited string
func Str(args ...api.Value) api.Value {
	return stdlib.SequenceToStr(api.Vector(args))
}

// ReaderStr converts the provided arguments to a delimited string
func ReaderStr(args ...api.Value) api.Value {
	if len(args) == 0 {
		return emptyString
	}

	var b bytes.Buffer
	b.WriteString(api.MaybeQuoteString(args[0]))
	for _, f := range args[1:] {
		b.WriteString(" ")
		b.WriteString(api.MaybeQuoteString(f))
	}
	return api.String(b.String())
}

// IsStr returns whether or not the provided value is a string
func IsStr(args ...api.Value) api.Value {
	_, ok := args[0].(api.String)
	return api.Bool(ok)
}
