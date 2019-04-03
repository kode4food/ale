package encoder

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/namespace"
	"gitlab.com/kode4food/ale/internal/runtime/isa"
)

type (
	// Type exposes an interface for stateful compiler encoding
	Type interface {
		api.Value

		Name() api.Name
		Parent() Type
		Child() Type
		NamedChild(api.Name) Type

		Globals() namespace.Type
		Append(...isa.Coder)
		Code() []isa.Code
		StackSize() int

		NewLabel() *Label

		Constants() api.Values
		AddConstant(api.Value) isa.Index

		Closure() api.Names
		ResolveClosure(api.LocalSymbol) (isa.Index, bool)

		PushArgs(api.Names, bool)
		PopArgs()
		ResolveArg(api.LocalSymbol) (isa.Index, bool, bool)

		LocalCount() int
		PushLocals()
		PopLocals()
		AddLocal(api.Name) isa.Index
		ResolveLocal(api.LocalSymbol) (isa.Index, bool)

		ResolveScope(api.LocalSymbol) (Scope, bool)
		InScope(api.LocalSymbol) bool
	}

	encoder struct {
		parent    Type
		globals   namespace.Type
		constants api.Values
		closure   api.Names
		name      api.Name
		args      argsStack
		locals    []Locals
		code      []isa.Code
		maxLocal  int
		nextLocal int
	}
)

func newEncoder(globals namespace.Type) *encoder {
	return &encoder{
		globals:   globals,
		constants: api.Values{},
		closure:   api.Names{},
		args:      argsStack{},
		locals:    []Locals{{}},
		code:      []isa.Code{},
	}
}

// NewEncoder instantiates a new Encoder
func NewEncoder(globals namespace.Type) Type {
	return newEncoder(globals)
}

func (e *encoder) child() *encoder {
	return &encoder{
		parent:    e,
		constants: api.Values{},
		closure:   api.Names{},
		args:      argsStack{},
		locals:    []Locals{{}},
		code:      []isa.Code{},
	}
}

// Child creates a child Type
func (e *encoder) Child() Type {
	return e.child()
}

func (e *encoder) NamedChild(name api.Name) Type {
	res := e.child()
	res.name = name
	return res
}

// Name returns the name of this encoder (example: a function's name)
func (e *encoder) Name() api.Name {
	return e.name
}

// Parent returns the parent of this encoder
func (e *encoder) Parent() Type {
	return e.parent
}

// Globals returns the global name/value map
func (e *encoder) Globals() namespace.Type {
	if e.globals != nil {
		return e.globals
	}
	if e.parent != nil {
		return e.parent.Globals()
	}
	return nil
}

// Append adds instructions to the Type's eventual output
func (e *encoder) Append(code ...isa.Coder) {
	for _, c := range code {
		e.code = append(e.code, c.Code())
	}
}

// Code returns the encoder's resulting VM instructions
func (e *encoder) Code() []isa.Code {
	ec := e.code
	res := make([]isa.Code, len(ec))
	copy(res, ec)
	isa.Verify(res)
	return res
}

// StackSize returns the encoder's calculated stack size
func (e *encoder) StackSize() int {
	res, _ := isa.CalculateStackSize(e.code)
	return res
}

func (e *encoder) String() string {
	return "encoder"
}
