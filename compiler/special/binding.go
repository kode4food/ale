package special

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
)

type uniqueNames map[data.Name]bool

// Error messages
const (
	ErrUnpairedBindings    = "binding must be a paired vector"
	ErrUnexpectedLetSyntax = "unexpected binding syntax: %s"
	ErrNameAlreadyBound    = "name is already bound in local scope: %s"
)

// Let encodes a binding form. Binding values are evaluated first, and are then
// bound to fresh names, meaning that mutual recursion is not supported
func Let(e encoder.Encoder, args ...data.Value) {
	bindings, body := parseLet(args...)
	generate.Locals(e, bindings, func(e encoder.Encoder) {
		generate.Block(e, body)
	})
}

// LetMutual encodes a binding form. First fresh names are introduced,
// and then binding values are evaluated with access to those names via
// the MutualScope
func LetMutual(e encoder.Encoder, args ...data.Value) {
	bindings, body := parseLet(args...)
	generate.MutualLocals(e, bindings, func(e encoder.Encoder) {
		generate.Block(e, body)
	})
}

func parseLet(args ...data.Value) (generate.Bindings, data.Vector) {
	data.AssertMinimum(2, len(args))
	bindings := parseLetBindings(args[0])
	return bindings, data.NewVector(args[1:]...)
}

func parseLetBindings(v data.Value) generate.Bindings {
	switch v := v.(type) {
	case data.List:
		names := uniqueNames{}
		res := generate.Bindings{}
		for f, r, ok := v.Split(); ok; f, r, ok = r.Split() {
			b := parseLetBinding(f.(data.Vector))
			names.see(b.Name)
			res = append(res, b)
		}
		return res
	case data.Vector:
		return generate.Bindings{
			parseLetBinding(v),
		}
	default:
		panic(fmt.Errorf(ErrUnexpectedLetSyntax, v))
	}
}

func parseLetBinding(b data.Vector) *generate.Binding {
	v := b.Values()
	if len(v) != 2 {
		panic(errors.New(ErrUnpairedBindings))
	}
	return &generate.Binding{
		Name:  v[0].(data.LocalSymbol).Name(),
		Value: v[1],
	}
}

func (u uniqueNames) see(n data.Name) {
	if _, ok := u[n]; ok {
		panic(fmt.Errorf(ErrNameAlreadyBound, n))
	}
	u[n] = true
}
