package encoder

import (
	"slices"

	"github.com/kode4food/ale/compiler/ir/analysis"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/runtime/isa"
)

type (
	// Encoder exposes an interface for stateful compiler encoding
	Encoder interface {
		Child() Encoder

		Emit(isa.Opcode, ...isa.Operand)
		Encode() *Encoded

		NewLabel() isa.Operand

		Globals() env.Namespace
		AddConstant(data.Value) isa.Operand
		ResolveClosure(data.Local) (*IndexedCell, bool)

		PushParams(data.Locals, bool)
		PopParams()
		ResolveParam(data.Local) (*IndexedCell, bool)

		PushLocals()
		PopLocals()
		AddLocal(data.Local, CellType) *IndexedCell
		ResolveLocal(data.Local) (*IndexedCell, bool)

		ResolveScoped(data.Local) (*ScopedCell, bool)
	}

	encoder struct {
		parent    Encoder
		globals   env.Namespace
		constants data.Vector
		closure   IndexedCells
		params    paramStack
		locals    []Locals
		code      isa.Instructions
		nextLabel isa.Operand
		nextLocal isa.Operand
		maxLocal  isa.Operand
	}

	// Encoded is a snapshot of the current Encoder's state. It is used as an
	// intermediate step in the compilation process, particularly as input to
	// the optimizer.
	Encoded struct {
		Code       isa.Instructions
		Globals    env.Namespace
		Constants  data.Vector
		Closure    IndexedCells
		LocalCount isa.Operand
		StackSize  isa.Operand
	}
)

// NewEncoder instantiates a new Encoder
func NewEncoder(globals env.Namespace) Encoder {
	return &encoder{
		globals: globals,
		locals:  []Locals{{}},
	}
}

func (e *encoder) child() *encoder {
	return &encoder{
		parent: e,
		locals: []Locals{{}},
	}
}

// Child creates a child Encoder
func (e *encoder) Child() Encoder {
	return e.child()
}

// Emit adds instructions to the Encoder's eventual output
func (e *encoder) Emit(oc isa.Opcode, args ...isa.Operand) {
	e.code = append(e.code, oc.New(args...))
}

// Encode returns the encoder's resulting abstract machine Instructions
func (e *encoder) Encode() *Encoded {
	return &Encoded{
		Code:       slices.Clone(e.code),
		Globals:    e.Globals(),
		Constants:  slices.Clone(e.constants),
		Closure:    slices.Clone(e.closure),
		LocalCount: e.maxLocal,
		StackSize:  e.stackSize(),
	}
}

// StackSize returns the encoder's calculated stack size
func (e *encoder) stackSize() isa.Operand {
	res, _ := analysis.CalculateStackSize(e.code)
	return res
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
