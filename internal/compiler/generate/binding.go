package generate

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
)

type (
	Binding struct {
		data.Value
		Name data.Local
	}

	Bindings []*Binding
)

func Locals(e encoder.Encoder, bindings Bindings, body Builder) {
	e.PushLocals()
	// Push the evaluated expressions to be bound
	for _, b := range bindings {
		Value(e, b.Value)
	}

	// Bind the popped expression results to names
	for i := len(bindings) - 1; i >= 0; i-- {
		b := bindings[i]
		l := e.AddLocal(b.Name, encoder.ValueCell)
		e.Emit(isa.Store, l.Index)
	}

	body(e)
	e.PopLocals()
}

func MutualLocals(e encoder.Encoder, bindings Bindings, body Builder) {
	e.PushLocals()
	// Create references
	cells := make(encoder.IndexedCells, len(bindings))
	for i, b := range bindings {
		c := e.AddLocal(b.Name, encoder.ReferenceCell)
		e.Emit(isa.NewRef)
		e.Emit(isa.Store, c.Index)
		cells[i] = c
	}

	// Push the evaluated expressions to be bound
	for _, b := range bindings {
		Value(e, b.Value)
	}

	// Bind the references
	for i := len(cells) - 1; i >= 0; i-- {
		c := cells[i]
		e.Emit(isa.Load, c.Index)
		e.Emit(isa.BindRef)
	}

	body(e)
	e.PopLocals()
}
