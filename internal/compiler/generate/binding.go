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
	Binder   func(encoder.Encoder, Bindings, Builder) error

	bindEncoder struct {
		encoder.Encoder
		cell *encoder.IndexedCell
	}
)

func Locals(e encoder.Encoder, bindings Bindings, body Builder) error {
	e.PushLocals()
	// Push the evaluated expressions to be bound
	for _, b := range bindings {
		if err := Value(e, b.Value); err != nil {
			return err
		}
	}

	// Bind the popped expression results to names
	for i := len(bindings) - 1; i >= 0; i-- {
		b := bindings[i]
		l, err := e.AddLocal(b.Name, encoder.ValueCell)
		if err != nil {
			return err
		}
		e.Emit(isa.Store, l.Index)
	}

	if err := body(e); err != nil {
		return err
	}
	return e.PopLocals()
}

func MutualLocals(e encoder.Encoder, bindings Bindings, body Builder) error {
	e.PushLocals()
	// Create references
	cells := make(encoder.IndexedCells, len(bindings))
	for i, b := range bindings {
		c, err := e.AddLocal(b.Name, encoder.ReferenceCell)
		if err != nil {
			return err
		}
		e.Emit(isa.NewRef)
		e.Emit(isa.Store, c.Index)
		cells[i] = c
	}

	// Push the evaluated expressions to be bound
	for i, b := range bindings {
		if err := BoundValue(e, cells[i], b.Value); err != nil {
			return err
		}
	}

	// Bind the references
	for i := len(cells) - 1; i >= 0; i-- {
		c := cells[i]
		e.Emit(isa.Load, c.Index)
		e.Emit(isa.BindRef)
	}

	if err := body(e); err != nil {
		return err
	}
	return e.PopLocals()
}

func (b *bindEncoder) Wrapped() encoder.Encoder {
	return b.Encoder
}

func (b *bindEncoder) Child() encoder.Encoder {
	res := *b
	res.Encoder = b.Encoder.Child()
	return &res
}

func BoundValue(e encoder.Encoder, c *encoder.IndexedCell, v data.Value) error {
	be := &bindEncoder{
		Encoder: e,
		cell:    c,
	}
	if err := Value(be, v); err != nil {
		return err
	}
	return nil
}
