package data

type (
	// ArityChecker is the interface for arity checks
	ArityChecker func(int) error

	// Convention describes the way a function should be called
	Convention int

	// Call is the type of function that can be turned into a Function
	Call func(...Value) Value

	// Caller provides the necessary methods for performing a runtime call
	Caller interface {
		Call(...Value) Value
		CheckArity(int) error
		Convention() Convention
	}

	// Function is a Value that provides a Caller interface
	Function interface {
		Value
		Caller
	}

	function struct {
		call  Call
		arity ArityChecker
		conv  Convention
	}
)

//go:generate stringer -output function_string.go -type Convention -linecomment
const (
	ApplicativeCall Convention = iota // applicative
	NormalCall                        // normal
)

// Applicative constructs an Applicative Function from a func that matches the
// standard calling signature
func Applicative(c Call, arity ...int) Function {
	return MakeApplicative(c, MakeChecker(arity...))
}

func makeFunction(c Call, conv Convention, arity ArityChecker) Function {
	return &function{
		call:  c,
		arity: arity,
		conv:  conv,
	}
}

// MakeApplicative constructs an applicative Function from a Caller
// and ArityChecker
func MakeApplicative(c Call, arity ArityChecker) Function {
	return makeFunction(c, ApplicativeCall, arity)
}

// Normal constructs a normal Function with
func Normal(c Call, arity ...int) Function {
	return MakeNormal(c, MakeChecker(arity...))
}

// MakeNormal constructs a normal Function from a Caller and ArityChecker
func MakeNormal(c Call, arity ArityChecker) Function {
	return makeFunction(c, NormalCall, arity)
}

// IsApplicative returns whether the function is applicative
func IsApplicative(f Function) bool {
	return f.Convention() == ApplicativeCall
}

// IsNormal returns whether the function is normal
func IsNormal(f Function) bool {
	return f.Convention() == NormalCall
}

func (f *function) CheckArity(argCount int) error {
	if a := f.arity; a != nil {
		return a(argCount)
	}
	return nil
}

func (f *function) Call(args ...Value) Value {
	return f.call(args...)
}

func (f *function) Convention() Convention {
	return f.conv
}

func (f *function) Type() Name {
	return Name(f.conv.String())
}

func (f *function) Equal(v Value) bool {
	if v, ok := v.(*function); ok {
		return f == v
	}
	return false
}

func (f *function) String() string {
	return DumpString(f)
}
