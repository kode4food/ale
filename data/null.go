package data

type (
	// NullType represents the null value, which is also the empty list
	NullType interface {
		List
		Nil()
	}

	nullValue struct{}
)

// EmptyList represents an empty List
var (
	EmptyList *nullValue
	Null      = EmptyList
)

func (*nullValue) Nil()  {}
func (*nullValue) List() {}

func (*nullValue) First() Value {
	return Null
}

func (*nullValue) Rest() Sequence {
	return Null
}

func (*nullValue) Split() (Value, Sequence, bool) {
	return Null, Null, false
}

func (*nullValue) IsEmpty() bool {
	return true
}

func (*nullValue) Reverse() Sequence {
	return Null
}

func (*nullValue) Prepend(value Value) Sequence {
	return NewList(value)
}

func (*nullValue) ElementAt(int) (Value, bool) {
	return Null, false
}

func (*nullValue) Count() int {
	return 0
}

func (*nullValue) String() string {
	return "()"
}
