package special

import (
	"fmt"

	"gitlab.com/kode4food/ale/compiler/arity"
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/compiler/generate"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Error messages
const (
	UnpairedBindings = "bindings must be paired"
	NameAlreadyBound = "name is already bound in local scope: %s"
)

type (
	letBinding struct {
		name  data.Name
		value data.Value
	}

	letBindings []*letBinding

	uniqueNames map[data.Name]bool
)

// Let encodes a binding form. Binding values are evaluated first, and
// are then bound to fresh names, meaning that mutual recursion is not
// supported
func Let(e encoder.Type, args ...data.Value) {
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
func LetMutual(e encoder.Type, args ...data.Value) {
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
	arity.AssertMinimum(2, len(args))
	b := args[0].(data.Vector)
	lb := len(b)
	if lb%2 != 0 {
		panic(fmt.Errorf(UnpairedBindings))
	}
	names := uniqueNames{}
	bindings := letBindings{}
	for i := 0; i < lb; i += 2 {
		name := b[i].(data.LocalSymbol).Name()
		names.see(name)
		value := b[i+1]
		bindings = append(bindings, newLetBinding(name, value))
	}
	return bindings, args[1:]
}

func newLetBinding(name data.Name, value data.Value) *letBinding {
	return &letBinding{
		name:  name,
		value: value,
	}
}

func (u uniqueNames) see(n data.Name) {
	if _, ok := u[n]; ok {
		panic(fmt.Errorf(NameAlreadyBound, n))
	}
	u[n] = true
}
