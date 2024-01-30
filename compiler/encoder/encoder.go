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
		Globals() env.Namespace
		NewLabel() isa.Operand

		AddConstant(data.Value) isa.Operand

		PushParams(data.Locals, bool)
		PopParams()

		PushLocals()
		PopLocals()
		AddLocal(data.Local, CellType) *IndexedCell

		ResolveScoped(data.Local) (*ScopedCell, bool)
		ResolveClosure(data.Local) (*IndexedCell, bool)
		ResolveParam(data.Local) (*IndexedCell, bool)
		ResolveLocal(data.Local) (*IndexedCell, bool)
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

	// Runnable is a flattened representation of the Encoded state that can be
	// executed by the abstract machine
	Runnable Encoded
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
		Code:       slices.Clone(e.code),
		Globals:    e.Globals(),
		Constants:  slices.Clone(e.constants),
		Closure:    slices.Clone(e.closure),
		LocalCount: e.maxLocal,
		StackSize:  analysis.MustCalculateStackSize(e.code),
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

// Runnable returns a flattened representation of the Encoded state that can be
// executed by the abstract machine
func (e *Encoded) Runnable() *Runnable {
	res := *e
	res.Code = isa.Flatten(e.Code)
	return (*Runnable)(&res)
}
