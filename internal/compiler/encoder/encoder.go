package encoder

import (
	"slices"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/basics"
	"github.com/kode4food/ale/internal/runtime/isa"
)

type (
	// Encoder exposes an interface for stateful compiler encoding
	Encoder interface {
		// Child creates a child encoder, allowing closure resolution
		Child() Encoder

		// Emit an instruction with the given opcode and operands
		Emit(isa.Opcode, ...isa.Operand)

		// Encode returns the encoded bytecode
		Encode() *Encoded

		// Globals returns the global namespace for this Encoder
		Globals() env.Namespace

		// NewLabel creates a new label for jump instructions
		NewLabel() isa.Operand

		// AddConstant adds a constant value and returns its index
		AddConstant(ale.Value) isa.Operand

		// PushParams pushes a new parameter frame
		PushParams(data.Locals, bool)

		// PopParams pops the current parameter frame
		PopParams()

		// PushLocals pushes a new local variable frame
		PushLocals()

		// PopLocals pops the current local variable frame
		PopLocals() error

		// AddLocal adds a local variable and returns its cell
		AddLocal(data.Local, CellType) (*IndexedCell, error)

		// ResolveScoped resolves a scoped variable
		ResolveScoped(data.Local) (*ScopedCell, bool)

		// ResolveClosure resolves a closure variable
		ResolveClosure(data.Local) (*IndexedCell, bool)

		// ResolveParam resolves a parameter variable
		ResolveParam(data.Local) (*IndexedCell, bool)

		// ResolveLocal resolves a local variable
		ResolveLocal(data.Local) (*IndexedCell, bool)
	}

	WrappedEncoder interface {
		Encoder

		// Wrapped returns the wrapped encoder
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
