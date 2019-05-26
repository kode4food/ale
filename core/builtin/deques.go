package builtin

import "gitlab.com/kode4food/ale/data"

// Deque constructs a new deque
func Deque(args ...data.Value) data.Value {
	return data.NewDeque(args...)
}

// IsDeque returns whether or not the provided value is a deque
func IsDeque(args ...data.Value) data.Value {
	_, ok := args[0].(*data.Deque)
	return data.Bool(ok)
}
