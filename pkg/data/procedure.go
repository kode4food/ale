package data

import (
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/types"
)

type (
	// ArityChecker is the interface for arity checks
	ArityChecker func(int) error

	// Call is the type of function that can be turned into a Procedure
	Call func(...Value) Value

	// Caller provides the necessary methods for performing a runtime call
	Caller interface {
		Call(...Value) Value
		CheckArity(int) error
	}

	// Procedure is a Value that provides a Caller interface
	Procedure interface {
		Value
		Caller
	}

	procedure struct {
		call  Call
		arity ArityChecker
	}
)

// compile-time check for interface implementation
var _ Procedure = &procedure{}

// MakeProcedure constructs a Procedure from a func that matches the standard
// calling signature
func MakeProcedure(c Call, arity ...int) Procedure {
	check, err := MakeArityChecker(arity...)
	if err != nil {
		panic(debug.ProgrammerError(err.Error()))
	}
	return &procedure{
		call:  c,
		arity: check,
	}
}

func (p *procedure) CheckArity(argc int) error {
	return p.arity(argc)
}

func (p *procedure) Call(args ...Value) Value {
	return p.call(args...)
}

func (p *procedure) Type() types.Type {
	return types.BasicProcedure
}

func (p *procedure) Equal(other Value) bool {
	return p == other
}

func (p *procedure) Get(key Value) (Value, bool) {
	return DumpMapped(p).Get(key)
}
