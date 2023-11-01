package data

import "github.com/kode4food/ale/internal/types"

type (
	// ArityChecker is the interface for arity checks
	ArityChecker func(int) error

	// Call is the type of lambda that can be turned into a Lambda
	Call func(...Value) Value

	// Caller provides the necessary methods for performing a runtime call
	Caller interface {
		Call(...Value) Value
		CheckArity(int) error
	}

	// Lambda is a Value that provides a Caller interface
	Lambda interface {
		Value
		Caller
	}

	lambda struct {
		call  Call
		arity ArityChecker
	}
)

// MakeLambda constructs a Lambda from a func that matches the standard
// calling signature
func MakeLambda(c Call, arity ...int) Lambda {
	return &lambda{
		call:  c,
		arity: MakeChecker(arity...),
	}
}

func (l *lambda) CheckArity(argCount int) error {
	return l.arity(argCount)
}

func (l *lambda) Call(args ...Value) Value {
	return l.call(args...)
}

func (l *lambda) Type() types.Type {
	return types.BasicLambda
}

func (l *lambda) Equal(v Value) bool {
	return l == v
}

func (l *lambda) String() string {
	return DumpString(l)
}
