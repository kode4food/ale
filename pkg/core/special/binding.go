package special

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/pkg/data"
)

type uniqueNames map[data.Local]bool

const (
	// ErrUnpairedBindings is raised when a Let binding Vector has fewer or
	// more than two elements
	ErrUnpairedBindings = "binding must be a paired vector"

	// ErrUnexpectedLetSyntax is raised when the Let bindings are not in the
	// form of a List or Vector
	ErrUnexpectedLetSyntax = "unexpected binding syntax: %s"

	// ErrNameAlreadyBound is raised when there's an attempt to bind the same
	// name more than once in a single Let scope
	ErrNameAlreadyBound = "name is already bound in local scope: %s"
)

// Let encodes a binding form. Binding values are evaluated first, and are then
// bound to fresh names, meaning that mutual recursion is not supported
func Let(e encoder.Encoder, args ...data.Value) error {
	return performBinding(e, generate.Locals, args...)
}

// LetMutual encodes a binding form. First fresh names are introduced,
// and then binding values are evaluated with access to those names via
// the MutualScope
func LetMutual(e encoder.Encoder, args ...data.Value) error {
	return performBinding(e, generate.MutualLocals, args...)
}

func performBinding(
	e encoder.Encoder, b generate.Binder, args ...data.Value,
) error {
	bindings, body := parseLet(args...)
	return b(e, bindings, func(e encoder.Encoder) error {
		return generate.Block(e, body)
	})
}

func parseLet(args ...data.Value) (generate.Bindings, data.Vector) {
	data.AssertMinimum(2, len(args))
	bindings := parseLetBindings(args[0])
	return bindings, args[1:]
}

func parseLetBindings(v data.Value) generate.Bindings {
	switch v := v.(type) {
	case *data.List:
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

func parseLetBinding(v data.Vector) *generate.Binding {
	if len(v) != 2 {
		panic(errors.New(ErrUnpairedBindings))
	}
	return &generate.Binding{
		Name:  v[0].(data.Local),
		Value: v[1],
	}
}

func (u uniqueNames) see(n data.Local) {
	if _, ok := u[n]; ok {
		panic(fmt.Errorf(ErrNameAlreadyBound, n))
	}
	u[n] = true
}
