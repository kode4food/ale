package data

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/types"
)

type (
	// ArityChecker is the interface for arity checks
	ArityChecker func(int) error

	// Call is the type of function that can be turned into a Procedure
	Call func(...ale.Value) ale.Value

	// Procedure is any Value that provides a calling interface
	Procedure interface {
		ale.Typed
		Call(...ale.Value) ale.Value
		CheckArity(int) error
	}

	procedure struct {
		call  Call
		arity ArityChecker
	}
)

// compile-time check for interface implementation
var _ interface {
	Mapped
	Procedure
} = (*procedure)(nil)

// MakeProcedure constructs a Procedure from a func that matches the standard
// calling signature
func MakeProcedure(c Call, arity ...int) Procedure {
	check, err := MakeArityChecker(arity...)
	if err != nil {
		panic(debug.ProgrammerErrorf("%w", err))
	}
	return &procedure{
		call:  c,
		arity: check,
	}
}

func (p *procedure) CheckArity(argc int) error {
	return p.arity(argc)
}

func (p *procedure) Call(args ...ale.Value) ale.Value {
	return p.call(args...)
}

func (p *procedure) Type() ale.Type {
	return types.BasicProcedure
}

func (p *procedure) Equal(other ale.Value) bool {
	return p == other
}

func (p *procedure) Get(key ale.Value) (ale.Value, bool) {
	return DumpMapped(p).Get(key)
}
