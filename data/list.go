package data

// List represents a singly-linked List
type List interface {
	List()
	Sequence
	Prepend(Value) Sequence
	Reverse() Sequence
	Indexed
	Counted
}

type list struct {
	first Value
	rest  List
	count int
}

// NewList creates a new List instance
func NewList(v ...Value) List {
	var res List = EmptyList
	for i, u := len(v)-1, 1; i >= 0; i, u = i-1, u+1 {
		res = &list{
			first: v[i],
			rest:  res,
			count: u,
		}
	}
	return res
}

func (l *list) List() {}

// First returns the first element of the List
func (l *list) First() Value {
	return l.first
}

// Rest returns the elements of the List that follow the first
func (l *list) Rest() Sequence {
	return l.rest
}

// IsEmpty returns whether this sequence is empty
func (l *list) IsEmpty() bool {
	return l.count == 0
}

// Split breaks the List into its components (first, rest, ok)
func (l *list) Split() (Value, Sequence, bool) {
	return l.first, l.rest, l.count != 0
}

// Car returns the first element of a Pair
func (l *list) Car() Value {
	return SequenceCar(l)
}

// Cdr returns the second element of a Pair
func (l *list) Cdr() Value {
	return SequenceCdr(l)
}

// Prepend inserts an element at the beginning of the List
func (l *list) Prepend(v Value) Sequence {
	return &list{
		first: v,
		rest:  l,
		count: l.count + 1,
	}
}

// Reverse returns a reversed copy of this List
func (l *list) Reverse() Sequence {
	if l.count <= 1 {
		return l
	}

	var res List = EmptyList
	var e List = l
	for d, u := e.Count(), 1; d > 0; e, d, u = e.Rest().(List), d-1, u+1 {
		res = &list{
			first: e.First(),
			rest:  res,
			count: u,
		}
	}
	return res
}

// Count returns the number of elements in the List
func (l *list) Count() int {
	return l.count
}

// ElementAt returns a specific element of the List
func (l *list) ElementAt(index int) (Value, bool) {
	if index > l.count-1 || index < 0 {
		return Nil, false
	}

	var e List = l
	for i := 0; i < index; i++ {
		e = e.Rest().(List)
	}
	return e.First(), true
}

// Call turns List into a callable type
func (l *list) Call() Call {
	return makeIndexedCall(l)
}

// Convention returns the function's calling convention
func (l *list) Convention() Convention {
	return ApplicativeCall
}

// CheckArity performs a compile-time arity check for the function
func (l *list) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

// String converts this List to a string
func (l *list) String() string {
	return MakeSequenceStr(l)
}
