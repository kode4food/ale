package encoder

import (
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
		Constants() data.Values
		AddConstant(data.Value) isa.Operand

		Closure() IndexedCells
		ResolveClosure(data.Name) (*IndexedCell, bool)

		PushParams(data.Names, bool)
		PopParams()
		ResolveParam(data.Name) (*IndexedCell, bool)

		LocalCount() isa.Operand
		PushLocals()
		PopLocals()
		AddLocal(data.Name, CellType) *IndexedCell
		ResolveLocal(data.Name) (*IndexedCell, bool)

		ResolveScoped(data.Name) (*ScopedCell, bool)
	}

	encoder struct {
		parent    Encoder
		globals   env.Namespace
		constants data.Values
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

// Child creates a child Type
func (e *encoder) Child() Encoder {
	return e.child()
}

// Emit adds instructions to the Type's eventual output
func (e *encoder) Emit(oc isa.Opcode, args ...isa.Operand) {
	e.code = append(e.code, oc.New(args...))
}

// Code returns the encoder's resulting abstract machine Instructions
func (e *encoder) Code() isa.Instructions {
	code := e.code
	analysis.Verify(code)
	res := make(isa.Instructions, len(code))
	copy(res, code)
	return res
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
