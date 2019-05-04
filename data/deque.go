package data

// Deque is a persistent double-ended queue
type Deque struct {
	head *List
	tail *List
}

// EmptyDeque represents an empty deque
var EmptyDeque = &Deque{
	head: EmptyList,
	tail: EmptyList,
}

// NewDeque creates a new double-ended queue
func NewDeque(v ...Value) *Deque {
	return &Deque{
		head: NewList(v...),
		tail: EmptyList,
	}
}

// First returns the first element of the deque
func (d *Deque) First() Value {
	f, _, _ := d.Split()
	return f
}

// Rest returns the elements of the deque that follow the first
func (d *Deque) Rest() Sequence {
	_, r, _ := d.Split()
	return r
}

// IsEmpty returns whether or not this sequence is empty
func (d *Deque) IsEmpty() bool {
	return d.head.IsEmpty() && d.tail.IsEmpty()
}

// Split breaks the deque into its components (first, rest, ok)
func (d *Deque) Split() (Value, Sequence, bool) {
	if f, r, ok := d.head.Split(); ok {
		rest := &Deque{
			head: r.(*List),
			tail: d.tail,
		}
		return f, rest, true
	}
	if f, r, ok := d.tail.Reverse().Split(); ok {
		rest := &Deque{
			head: r.(*List),
			tail: EmptyList,
		}
		return f, rest, true
	}
	return Nil, EmptyDeque, false
}

// Count returns the number of elements in the deque
func (d *Deque) Count() int {
	return d.head.count + d.tail.count
}

// Reverse returns a reversed version of the sequence
func (d *Deque) Reverse() Sequence {
	return &Deque{
		head: d.tail,
		tail: d.head,
	}
}

// Prepend inserts an element at the beginning of the deque
func (d *Deque) Prepend(v Value) Sequence {
	return &Deque{
		head: d.head.Prepend(v).(*List),
		tail: d.tail,
	}
}

// Append appends elements to the end of the deque
func (d *Deque) Append(args ...Value) Sequence {
	var newTail Sequence = d.tail
	for _, v := range args {
		newTail = newTail.Prepend(v)
	}
	return &Deque{
		head: d.head,
		tail: newTail.(*List),
	}
}

func (d *Deque) String() string {
	return MakeSequenceStr(d)
}
