package data

// Keyword is a Value that represents a Name that resolves to itself
type Keyword Name

// Name returns the name of the Keyword
func (k Keyword) Name() Name {
	return Name(k)
}

// Call turns Keyword into a callable type
func (k Keyword) Call() Call {
	return func(args ...Value) Value {
		m := args[0].(Mapped)
		res, ok := m.Get(k)
		if !ok && len(args) > 1 {
			return args[1]
		}
		return res
	}
}

// Convention returns the function's calling convention
func (k Keyword) Convention() Convention {
	return ApplicativeCall
}

// CheckArity performs a compile-time arity check for the function
func (k Keyword) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

// String converts Keyword into a string
func (k Keyword) String() string {
	return ":" + string(k)
}
