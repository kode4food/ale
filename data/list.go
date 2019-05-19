package data

// List contains a node to a singly-linked List
type List struct {
	first Value
	rest  *List
	size  int
}

// EmptyList represents an empty List
var EmptyList = &List{}

// NewList creates a new List instance
func NewList(v ...Value) *List {
	r := EmptyList
	for i := len(v) - 1; i >= 0; i-- {
		r = &List{
			first: v[i],
			rest:  r,
			size:  r.size + 1,
		}
	}
	return r
}

// First returns the first element of the List
func (l *List) First() Value {
	return l.first
}

// Rest returns the elements of the List that follow the first
func (l *List) Rest() Sequence {
	return l.rest
}

// IsEmpty returns whether or not this sequence is empty
func (l *List) IsEmpty() bool {
	return l.size == 0
}

// Split breaks the List into its components (first, rest, ok)
func (l *List) Split() (Value, Sequence, bool) {
	return l.first, l.rest, l.size != 0
}

// Prepend inserts an element at the beginning of the List
func (l *List) Prepend(v Value) Sequence {
	return &List{
		first: v,
		rest:  l,
		size:  l.size + 1,
	}
}

// Reverse returns a reversed copy of this List
func (l *List) Reverse() Sequence {
	if l.size <= 1 {
		return l
	}
	res := EmptyList
	for cur, cnt := l, 1; cur.size > 0; cur, cnt = cur.rest, cnt+1 {
		res = &List{
			first: cur.first,
			rest:  res,
			size:  cnt,
		}
	}
	return res
}

// Size returns the number of elements in the List
func (l *List) Size() int {
	return l.size
}

// ElementAt returns a specific element of the List
func (l *List) ElementAt(index int) (Value, bool) {
	if index > l.size-1 || index < 0 {
		return Nil, false
	}

	e := l
	for i := 0; i < index; i++ {
		e = e.rest
	}
	return e.first, true
}

// Caller turns List into a callable type
func (l *List) Caller() Call {
	return makeIndexedCall(l)
}

// String converts this List to a string
func (l *List) String() string {
	return MakeSequenceStr(l)
}
