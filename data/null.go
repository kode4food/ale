package data

type (
	// Null represents a null value, which is also the empty list
	Null interface {
		null() // marker
		List
	}

	nilValue struct{}
)

// EmptyList represents an empty List
var (
	EmptyList *nilValue
	Nil       = EmptyList
)

func (*nilValue) null() {}
func (*nilValue) list() {}

func (*nilValue) First() Value {
	return Nil
}

func (*nilValue) Rest() Sequence {
	return EmptyList
}

func (*nilValue) Split() (Value, Sequence, bool) {
	return Nil, EmptyList, false
}

func (*nilValue) IsEmpty() bool {
	return true
}

func (*nilValue) Reverse() Sequence {
	return EmptyList
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

func (*nilValue) Equal(v Value) bool {
	if _, ok := v.(*nilValue); ok {
		return true
	}
	return false
}

func (*nilValue) String() string {
	return "()"
}

func (*nilValue) HashCode() uint64 {
	return 0
}
