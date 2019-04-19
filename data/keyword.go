package data

// Keyword is a Value that represents a Name that resolves to itself
type Keyword Name

// Name returns the name of the Keyword
func (k Keyword) Name() Name {
	return Name(k)
}

// Caller turns Keyword into a callable type
func (k Keyword) Caller() Call {
	return func(args ...Value) Value {
		m := args[0].(Mapped)
		res, ok := m.Get(k)
		if !ok && len(args) > 1 {
			return args[1]
		}
		return res
	}
}

// String converts Keyword into a string
func (k Keyword) String() string {
	return ":" + string(k)
}
