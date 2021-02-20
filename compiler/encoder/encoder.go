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
		data.Value

		Parent() Encoder
		Child() Encoder

		Emit(isa.Opcode, ...isa.Coder)
		Code() isa.Instructions
		StackSize() int

		NewLabel() *Label

		Globals() env.Namespace
		Constants() data.Values
		AddConstant(data.Value) isa.Index

		Closure() IndexedCells
		ResolveClosure(data.Name) (*IndexedCell, bool)

		PushArgs(data.Names, bool)
		PopArgs()
		ResolveArg(data.Name) (*IndexedCell, bool)

		LocalCount() int
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
		args      argsStack
		locals    []Locals
		code      isa.Instructions
		maxLocal  int
		nextLocal int
		nextLabel int
	}
)

func newEncoder(globals env.Namespace) *encoder {
	return &encoder{
		globals: globals,
	}
}

// NewEncoder instantiates a new Encoder
func NewEncoder(globals env.Namespace) Encoder {
	return newEncoder(globals)
}

func (e *encoder) child() *encoder {
	return &encoder{
		parent: e,
	}
}

// Child creates a child Type
func (e *encoder) Child() Encoder {
	return e.child()
}

// Parent returns the parent of this encoder
func (e *encoder) Parent() Encoder {
	return e.parent
}

// Emit adds instructions to the Type's eventual output
func (e *encoder) Emit(oc isa.Opcode, args ...isa.Coder) {
	words := make([]isa.Word, len(args))
	for i, a := range args {
		words[i] = a.Word()
	}
	e.code = append(e.code, isa.New(oc, words...))
}

// Word returns the encoder's resulting VM instructions
func (e *encoder) Code() isa.Instructions {
	code := e.code
	analysis.Verify(code)
	res := make(isa.Instructions, len(code))
	copy(res, code)
	return res
}

// StackSize returns the encoder's calculated stack size
func (e *encoder) StackSize() int {
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

func (e *encoder) Equal(v data.Value) bool {
	if v, ok := v.(*encoder); ok {
		return e == v
	}
	return false
}

func (e *encoder) String() string {
	return "encoder"
}
