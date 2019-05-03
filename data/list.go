package data

// List contains a node to a singly-linked List
type List struct {
	first Value
	rest  *List
	count int
}

// EmptyList represents an empty List
var EmptyList *List

// NewList creates a new List instance
func NewList(v ...Value) *List {
	r := EmptyList
	for i := len(v) - 1; i >= 0; i-- {
		r = &List{
			first: v[i],
			rest:  r,
			count: r.count + 1,
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
	return l == EmptyList
}

// Split breaks the List into its components (first, rest, ok)
func (l *List) Split() (Value, Sequence, bool) {
	return l.first, l.rest, l != EmptyList
}

// Prepend inserts an element at the beginning of the List
func (l *List) Prepend(v Value) Sequence {
	return &List{
		first: v,
		rest:  l,
		count: l.count + 1,
	}
}

// Reverse returns a reversed copy of this List
func (l *List) Reverse() Sequence {
	res := EmptyList
	for cur, cnt := l, 0; cur.count > 0; cur, cnt = cur.rest, cnt+1 {
		res = &List{
			first: cur.first,
			rest:  res,
			count: cnt,
		}
	}
	return res
}

// Count returns the number of elements in the List
func (l *List) Count() int {
	return l.count
}

// ElementAt returns a specific element of the List
func (l *List) ElementAt(index int) (Value, bool) {
	if index > l.count-1 || index < 0 {
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

func init() {
	EmptyList = &List{
		first: Nil,
		count: 0,
	}
	EmptyList.rest = EmptyList
}
