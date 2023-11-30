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
		Code() isa.Instructions
		StackSize() isa.Operand

		NewLabel() isa.Operand

		Globals() env.Namespace
		Constants() data.Vector
		AddConstant(data.Value) isa.Operand

		Closure() IndexedCells
		ResolveClosure(data.Local) (*IndexedCell, bool)

		PushParams(data.Locals, bool)
		PopParams()
		ResolveParam(data.Local) (*IndexedCell, bool)

		LocalCount() isa.Operand
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

// Code returns the encoder's resulting abstract machine Instructions
func (e *encoder) Code() isa.Instructions {
	return slices.Clone(e.code)
}

// StackSize returns the encoder's calculated stack size
func (e *encoder) StackSize() isa.Operand {
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
