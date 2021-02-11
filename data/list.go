package data

type (
	// List represents a singly-linked List
	List interface {
		list() // marker
		Sequence
		RandomAccess
		Prepender
		Reverser
	}

	list struct {
		first Value
		rest  List
		count int
	}
)

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

func (*list) list() {}

func (l *list) First() Value {
	return l.first
}

func (l *list) Rest() Sequence {
	return l.rest
}

func (l *list) IsEmpty() bool {
	return false
}

func (l *list) Split() (Value, Sequence, bool) {
	return l.first, l.rest, true
}

func (l *list) Car() Value {
	return l.first
}

func (l *list) Cdr() Value {
	return l.rest
}

func (l *list) Prepend(v Value) Sequence {
	return &list{
		first: v,
		rest:  l,
		count: l.count + 1,
	}
}

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

func (l *list) Count() int {
	return l.count
}

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

func (l *list) Call(args ...Value) Value {
	return indexedCall(l, args)
}

func (l *list) Convention() Convention {
	return ApplicativeCall
}

func (l *list) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

func (l *list) Equal(v Value) bool {
	if v, ok := v.(*list); ok {
		if l == v {
			return true
		}
		if l.count != v.count || !l.first.Equal(v.first) {
			return false
		}
		return l.rest.Equal(v.rest)
	}
	return false
}

func (l *list) String() string {
	return MakeSequenceStr(l)
}

func (l *list) HashCode() uint64 {
	var h uint64
	for f, r, ok := l.Split(); ok; f, r, ok = r.Split() {
		h ^= HashCode(f)
	}
	return h
}
