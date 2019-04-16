package api

import "fmt"

type (
	// Call is the basic function type
	Call func(...Value) Value

	// Caller is implemented by types that can return a call interface
	Caller interface {
		Caller() Call
	}

	// Convention describes the way a function should be called
	Convention int

	// ArityChecker returns an error if valid argument count isn't provided
	ArityChecker func(int) error

	// Function wraps a function call with additional context
	Function struct {
		Call
		Convention
		ArityChecker
	}
)

//go:generate stringer -output function_string.go -type Convention -linecomment
const (
	ApplicativeCall Convention = iota // Applicative
	NormalCall                        // Normal
	MacroCall                         // Macro
)

// ApplicativeFunction wraps an applicative function call
func ApplicativeFunction(c Call) *Function {
	return &Function{
		Call:       c,
		Convention: ApplicativeCall,
	}
}

// NormalFunction wraps a normal function call
func NormalFunction(c Call) *Function {
	return &Function{
		Call:       c,
		Convention: NormalCall,
	}
}

// MacroFunction wraps a macro
func MacroFunction(c Call) *Function {
	return &Function{
		Call:       c,
		Convention: MacroCall,
	}
}

// Caller returns the calling interface for a call
func (c Call) Caller() Call {
	return c
}

func (c Call) String() string {
	return "function"
}

// Caller returns the calling interface for a function
func (v *Function) Caller() Call {
	return v.Call
}

// Type returns the type for a function
func (v *Function) Type() Name {
	return Name(fmt.Sprintf("%s", v.Convention))
}

// IsApplicative returns whether or not this function is applicative
func (v *Function) IsApplicative() bool {
	return v.Convention == ApplicativeCall
}

// IsNormal returns whether or not this function is normal
func (v *Function) IsNormal() bool {
	return v.Convention == NormalCall
}

// IsMacro returns whether or not this function is a macro
func (v *Function) IsMacro() bool {
	return v.Convention == MacroCall
}

// CheckArity checks to see if the argument count is valid
func (v *Function) CheckArity(count int) error {
	if a := v.ArityChecker; a != nil {
		return a(count)
	}
	return nil
}

func (v *Function) String() string {
	return DumpString(v)
}
