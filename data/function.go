package data

type (
	// Call is the basic function type
	Call func(...Value) Value

	// Caller is any value that returns a calling interface
	Caller interface {
		Value
		Caller() Call
	}

	// ArityChecker is the interface for arity checks
	ArityChecker func(int) error

	// Convention describes the way a function should be called
	Convention int

	// Function provides a caller, arity check, and calling convention
	Function interface {
		Caller
		CheckArity(int) error
		Convention() Convention
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

// Caller turns Call into a callable type
func (c Call) Caller() Call {
	return c
}

func (c Call) String() string {
	return "call"
}

func makeFunction(conv Convention, call Call, arity ArityChecker) Function {
	return &function{
		call:  call,
		arity: arity,
		conv:  conv,
	}
}

// MakeApplicative constructs an applicative function
func MakeApplicative(call Call, arity ArityChecker) Function {
	return makeFunction(ApplicativeCall, call, arity)
}

// MakeNormal constructs a normal function
func MakeNormal(call Call, arity ArityChecker) Function {
	return makeFunction(NormalCall, call, arity)
}

// IsApplicative returns whether or not the function is applicative
func IsApplicative(f Function) bool {
	return f.Convention() == ApplicativeCall
}

// IsNormal returns whether or not the function is normal
func IsNormal(f Function) bool {
	return f.Convention() == NormalCall
}

// Caller returns the calling interface for a function
func (f *function) Caller() Call {
	return f.call
}

// CheckArity checks to see if the argument count is valid
func (f *function) CheckArity(count int) error {
	if a := f.arity; a != nil {
		return a(count)
	}
	return nil
}

func (f *function) Convention() Convention {
	return f.conv
}

func (f *function) Type() Name {
	return Name(f.conv.String())
}

func (f *function) String() string {
	return DumpString(f)
}
