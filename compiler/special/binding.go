package special

import (
	"errors"
	"fmt"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

type (
	letBinding struct {
		name  data.Name
		value data.Value
	}

	letBindings []*letBinding

	uniqueNames map[data.Name]bool
)

// Error messages
const (
	ErrUnpairedBindings    = "binding must be a paired vector"
	ErrUnexpectedLetSyntax = "unexpected binding syntax: %s"
	ErrNameAlreadyBound    = "name is already bound in local scope: %s"
)

// Let encodes a binding form. Binding values are evaluated first, and
// are then bound to fresh names, meaning that mutual recursion is not
// supported
func Let(e encoder.Encoder, args ...data.Value) {
	bindings, body := parseLet(args...)

	e.PushLocals()
	// Push the evaluated expressions to be bound
	for _, b := range bindings {
		generate.Value(e, b.value)
	}

	// Bind the popped expression results to names
	for i := len(bindings) - 1; i >= 0; i-- {
		b := bindings[i]
		l := e.AddLocal(b.name, encoder.ValueCell)
		e.Emit(isa.Store, l.Index)
	}

	generate.Block(e, body)
	e.PopLocals()
}

// LetMutual encodes a binding form. First fresh names are introduced,
// and then binding values are evaluated with access to those names via
// the MutualScope
func LetMutual(e encoder.Encoder, args ...data.Value) {
	bindings, body := parseLet(args...)

	e.PushLocals()
	// Create references
	cells := make(encoder.IndexedCells, len(bindings))
	for i, b := range bindings {
		c := e.AddLocal(b.name, encoder.ReferenceCell)
		e.Emit(isa.NewRef)
		e.Emit(isa.Store, c.Index)
		cells[i] = c
	}

	// Push the evaluated expressions to be bound
	for _, b := range bindings {
		generate.Value(e, b.value)
	}

	// Bind the references
	for i := len(cells) - 1; i >= 0; i-- {
		c := cells[i]
		e.Emit(isa.Load, c.Index)
		e.Emit(isa.BindRef)
	}

	generate.Block(e, body)
	e.PopLocals()
}

func parseLet(args ...data.Value) (letBindings, data.Vector) {
	data.AssertMinimum(2, len(args))
	bindings := parseLetBindings(args[0])
	return bindings, args[1:]
}

func parseLetBindings(v data.Value) letBindings {
	switch v := v.(type) {
	case data.List:
		names := uniqueNames{}
		res := letBindings{}
		for f, r, ok := v.Split(); ok; f, r, ok = r.Split() {
			b := parseLetBinding(f.(data.Vector))
			names.see(b.name)
			res = append(res, b)
		}
		return res
	case data.Vector:
		return letBindings{
			parseLetBinding(v),
		}
	default:
		panic(fmt.Errorf(ErrUnexpectedLetSyntax, v))
	}
}

func parseLetBinding(b data.Vector) *letBinding {
	if len(b) != 2 {
		panic(errors.New(ErrUnpairedBindings))
	}
	return &letBinding{
		name:  b[0].(data.LocalSymbol).Name(),
		value: b[1],
	}
}

func (u uniqueNames) see(n data.Name) {
	if _, ok := u[n]; ok {
		panic(fmt.Errorf(ErrNameAlreadyBound, n))
	}
	u[n] = true
}
