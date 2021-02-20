package data

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"sort"
)

type (
	// Object maps a set of Values, known as keys, to another set of Values
	Object interface {
		object() // marker
		Sequence
		Mapped
		Counted
		Caller
	}

	object struct {
		pair     Pair
		children [32]*object
	}

	emptyObject struct{}
)

// Standard Keys
const (
	NameKey     = Keyword("name")
	TypeKey     = Keyword("type")
	CountKey    = Keyword("count")
	InstanceKey = Keyword("instance")
)

// Error messages
const (
	ErrMapNotPaired = "map does not contain an even number of elements"
)

// EmptyObject represents an empty Object
var (
	EmptyObject *emptyObject

	objectHash = rand.Uint64()
)

// NewObject instantiates a new Object instance. Based on Phil Bagwell's
// Hashed Array Mapped Trie data structure. More information can be
// found at http://lampwww.epfl.ch/papers/idealhashtrees.pdf
func NewObject(pairs ...Pair) Object {
	var res Object = EmptyObject
	for _, p := range pairs {
		res = res.Put(p).(Object)
	}
	return res
}

// ValuesToObject interprets a set of Values as an Object
func ValuesToObject(v ...Value) (Object, error) {
	if len(v)%2 != 0 {
		return nil, errors.New(ErrMapNotPaired)
	}
	var res Object = EmptyObject
	for i := len(v) - 2; i >= 0; i -= 2 {
		res = res.Put(NewCons(v[i], v[i+1])).(Object)
	}
	return res, nil
}

func (*object) object() {}

func (o *object) Get(k Value) (Value, bool) {
	h := HashCode(k)
	return o.get(k, h)
}

func (o *object) get(k Value, hash uint64) (Value, bool) {
	if o.pair.Car().Equal(k) {
		return o.pair.Cdr(), true
	}
	bucket := o.children[hash&0x1f]
	if bucket != nil {
		return bucket.get(k, hash>>5)
	}
	return Nil, false
}

func (o *object) Put(p Pair) Sequence {
	h := HashCode(p.Car())
	return o.put(p, h)
}

func (o *object) put(p Pair, hash uint64) *object {
	if o.pair.Car().Equal(p.Car()) {
		return &object{
			pair:     p,
			children: o.children,
		}
	}

	idx := hash & 0x1f
	bucket := o.children[idx]
	if bucket == nil {
		bucket = &object{pair: p}
	} else {
		bucket = bucket.put(p, hash>>5)
	}

	// return a copy with the new bucket
	res := *o
	res.children[idx] = bucket
	return &res
}

func (o *object) Remove(k Value) (Value, Sequence, bool) {
	h := HashCode(k)
	if v, r, ok := o.remove(k, h); ok {
		if r != nil {
			return v, r, true
		}
		return v, EmptyObject, true
	}
	return Nil, o, false
}

func (o *object) remove(k Value, hash uint64) (Value, *object, bool) {
	if o.pair.Car().Equal(k) {
		return o.pair.Cdr(), o.promote(), true
	}
	idx := hash & 0x1f
	if bucket := o.children[idx]; bucket != nil {
		if v, r, ok := bucket.remove(k, hash>>5); ok {
			res := *o
			res.children[idx] = r
			return v, &res, true
		}
	}
	return nil, nil, false
}

func (o *object) promote() *object {
	for i, c := range o.children {
		if c != nil {
			res := *o
			res.pair = c.pair
			res.children[i] = c.promote()
			return &res
		}
	}
	return nil
}

func (o *object) First() Value {
	f, _, _ := o.Split()
	return f
}

func (o *object) Rest() Sequence {
	_, r, _ := o.Split()
	return r
}

func (o *object) Split() (Value, Sequence, bool) {
	if f := o.pair; f != nil {
		if r := o.promote(); r != nil {
			return f, r, true
		}
		return f, EmptyObject, true
	}
	return Nil, EmptyObject, false
}

func (o *object) Count() int {
	res := 1
	for _, c := range o.children {
		if c != nil {
			res += c.Count()
		}
	}
	return res
}

func (o *object) IsEmpty() bool {
	return false
}

func (o *object) Call(args ...Value) Value {
	return mappedCall(o, args)
}

func (o *object) Convention() Convention {
	return ApplicativeCall
}

func (o *object) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

func (o *object) Equal(v Value) bool {
	if v, ok := v.(*object); ok {
		lp := sortedPairs(o.Pairs())
		rp := sortedPairs(v.Pairs())
		if len(lp) != len(rp) {
			return false
		}
		for i, l := range lp {
			if !l.Equal(rp[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (o *object) HashCode() uint64 {
	return o.hashCode(objectHash)
}

func (o *object) hashCode(acc uint64) uint64 {
	h := acc * HashCode(o.pair.Car()) * HashCode(o.pair.Cdr())
	for _, c := range o.children {
		if c != nil {
			h = c.hashCode(h)
		}
	}
	return h
}

func (o *object) Pairs() Pairs {
	return o.pairs(Pairs{})
}

func (o *object) pairs(p Pairs) Pairs {
	p = append(p, o.pair)
	for _, c := range o.children {
		if c != nil {
			p = c.pairs(p)
		}
	}
	return p
}

func (o *object) String() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	for i, p := range sortedPairs(o.Pairs()) {
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(MaybeQuoteString(p.Car()))
		buf.WriteString(" ")
		buf.WriteString(MaybeQuoteString(p.Cdr()))
	}
	buf.WriteString("}")
	return buf.String()
}

func sortedPairs(p Pairs) Pairs {
	sort.Slice(p, func(l, r int) bool {
		ls := fmt.Sprintf("%s", p[l].Car().String())
		rs := fmt.Sprintf("%s", p[r].Car().String())
		return ls < rs
	})
	return p
}

func (*emptyObject) object() {}

func (*emptyObject) Get(Value) (Value, bool) {
	return Nil, false
}

func (*emptyObject) Put(p Pair) Sequence {
	return &object{
		pair: p,
	}
}

func (*emptyObject) Remove(Value) (Value, Sequence, bool) {
	return Nil, EmptyObject, false
}

func (*emptyObject) IsEmpty() bool {
	return true
}

func (*emptyObject) Count() int {
	return 0
}

func (*emptyObject) Split() (Value, Sequence, bool) {
	return Nil, EmptyObject, false
}

func (*emptyObject) First() Value {
	return Nil
}

func (*emptyObject) Rest() Sequence {
	return EmptyObject
}

func (*emptyObject) Call(args ...Value) Value {
	return mappedCall(EmptyObject, args)
}

func (*emptyObject) Convention() Convention {
	return ApplicativeCall
}

func (*emptyObject) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

func (*emptyObject) Equal(v Value) bool {
	if _, ok := v.(*emptyObject); ok {
		return true
	}
	return false
}

func (*emptyObject) String() string {
	return "{}"
}

func (*emptyObject) HashCode() uint64 {
	return objectHash
}
