package encoder

import (
	"slices"

	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

type (
	// Encoder exposes an interface for stateful compiler encoding
	Encoder interface {
		Child() Encoder

		Emit(isa.Opcode, ...isa.Operand)
		Encode() *Encoded
		Globals() env.Namespace
		NewLabel() isa.Operand

		AddConstant(data.Value) isa.Operand

		PushParams(data.Locals, bool)
		PopParams()

		PushLocals()
		PopLocals() error
		AddLocal(data.Local, CellType) (*IndexedCell, error)

		ResolveScoped(data.Local) (*ScopedCell, bool)
		ResolveClosure(data.Local) (*IndexedCell, bool)
		ResolveParam(data.Local) (*IndexedCell, bool)
		ResolveLocal(data.Local) (*IndexedCell, bool)
	}

	WrappedEncoder interface {
		Encoder
		Wrapped() Encoder
	}

	encoder struct {
		locals    []Locals
		code      isa.Instructions
		parent    Encoder
		globals   env.Namespace
		constants data.Vector
		closure   IndexedCells
		params    paramStack
		nextLabel isa.Operand
		nextLocal isa.Operand
	}
)

// NewEncoder instantiates a new Encoder
func NewEncoder(globals env.Namespace) Encoder {
	return &encoder{
		globals: globals,
		locals:  []Locals{{}},
	}
}

// Child creates a child Encoder
func (e *encoder) Child() Encoder {
	return &encoder{
		parent: e,
		locals: []Locals{{}},
	}
}

// Emit adds instructions to the Encoder's eventual output
func (e *encoder) Emit(oc isa.Opcode, args ...isa.Operand) {
	e.code = append(e.code, oc.New(args...))
}

// Encode returns the encoder's resulting abstract machine Instructions
func (e *encoder) Encode() *Encoded {
	return &Encoded{
		Code:      e.code,
		Globals:   e.Globals(),
		Constants: slices.Clone(e.constants),
		Closure: basics.Map(e.closure, func(elem *IndexedCell) data.Local {
			return elem.Name
		}),
	}
}

// Globals returns the global name/value map
func (e *encoder) Globals() env.Namespace {
	if e.globals != nil {
		return e.globals
	}
	if e.parent != nil {
		return e.parent.Globals()
	}
	return nil
}
