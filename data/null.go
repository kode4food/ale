package data

type (
	// NullType represents a null value, which is also the empty list
	NullType interface {
		List
		Null()
	}

	nilValue struct{}
)

// EmptyList represents an empty List
var (
	EmptyList *nilValue
	Nil       = EmptyList
)

func (*nilValue) Null() {}
func (*nilValue) List() {}

func (*nilValue) First() Value {
	return Nil
}

func (*nilValue) Rest() Sequence {
	return Nil
}

func (*nilValue) Split() (Value, Sequence, bool) {
	return Nil, Nil, false
}

func (*nilValue) IsEmpty() bool {
	return true
}

func (*nilValue) Reverse() Sequence {
	return Nil
}

func (*nilValue) Prepend(value Value) Sequence {
	return NewList(value)
}

func (*nilValue) ElementAt(int) (Value, bool) {
	return Nil, false
}

func (*nilValue) Count() int {
	return 0
}

func (*nilValue) String() string {
	return "()"
}
