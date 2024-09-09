package data

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/kode4food/ale/internal/types"
)

// Object maps a set of Values, known as keys, to another set of Values
type Object struct {
	pair     Pair
	children [bucketSize]*Object
}

const (
	bucketBits = 5
	bucketSize = 1 << bucketBits
	bucketMask = bucketSize - 1
)

// ErrMapNotPaired is raised when a call to ValuesToObject receives an odd
// number of args, meaning it won't be capable of zipping them into an Object
const ErrMapNotPaired = "map does not contain an even number of elements"

// EmptyObject represents an empty Object
var (
	EmptyObject *Object
	emptyPairs  = Pairs{}
	objectHash  = rand.Uint64()

	// compile-time checks for interface implementation
	_ interface {
		Caller
		CountedSequence
		Hashed
		Mapper
		Typed
	} = &Object{}
)

// NewObject instantiates a new Object instance. Based on Phil Bagwell's Hashed
// Array Mapped Trie data structure, though not as space efficient. More
// information on HAMT's can be found at:
//
//	http://lampwww.epfl.ch/papers/idealhashtrees.pdf
func NewObject(pairs ...Pair) *Object {
	res := EmptyObject
	for _, p := range pairs {
		res = res.Put(p).(*Object)
	}
	return res
}

// ValuesToObject interprets a set of Values as an Object
func ValuesToObject(v ...Value) (*Object, error) {
	if len(v)%2 != 0 {
		return nil, errors.New(ErrMapNotPaired)
	}
	res := EmptyObject
	for i := len(v) - 2; i >= 0; i -= 2 {
		res = res.Put(NewCons(v[i], v[i+1])).(*Object)
	}
	return res, nil
}

func (o *Object) Get(k Value) (Value, bool) {
	if o == nil {
		return Null, false
	}
	h := HashCode(k)
	return o.get(k, h)
}

func (o *Object) get(k Value, hash uint64) (Value, bool) {
	if o.pair.Car().Equal(k) {
		return o.pair.Cdr(), true
	}
	bucket := o.children[hash&bucketMask]
	if bucket != nil {
		return bucket.get(k, hash>>bucketBits)
	}
	return Null, false
}

func (o *Object) Put(p Pair) Sequence {
	if o == nil {
		return &Object{
			pair: p,
		}
	}
	h := HashCode(p.Car())
	return o.put(p, h)
}

func (o *Object) put(p Pair, hash uint64) *Object {
	if o.pair.Car().Equal(p.Car()) {
		return &Object{
			pair:     p,
			children: o.children,
		}
	}

	idx := hash & bucketMask
	bucket := o.children[idx]
	if bucket == nil {
		bucket = &Object{pair: p}
	} else {
		bucket = bucket.put(p, hash>>bucketBits)
	}

	// return a copy with the new bucket
	res := *o
	res.children[idx] = bucket
	return &res
}

func (o *Object) Remove(k Value) (Value, Sequence, bool) {
	if o == nil {
		return Null, EmptyObject, false
	}
	h := HashCode(k)
	if v, r, ok := o.remove(k, h); ok {
		if r != nil {
			return v, r, true
		}
		return v, EmptyObject, true
	}
	return Null, o, false
}

func (o *Object) remove(k Value, hash uint64) (Value, *Object, bool) {
	if o.pair.Car().Equal(k) {
		return o.pair.Cdr(), o.promote(), true
	}
	idx := hash & bucketMask
	if bucket := o.children[idx]; bucket != nil {
		if v, r, ok := bucket.remove(k, hash>>bucketBits); ok {
			res := *o
			res.children[idx] = r
			return v, &res, true
		}
	}
	return nil, nil, false
}

func (o *Object) promote() *Object {
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

func (o *Object) Car() Value {
	if o == nil {
		return Null
	}
	if f := o.pair; f != nil {
		return f
	}
	return Null
}

func (o *Object) Cdr() Value {
	if o == nil {
		return EmptyObject
	}
	if r := o.promote(); r != nil {
		return r
	}
	return EmptyObject
}

func (o *Object) Split() (Value, Sequence, bool) {
	if o == nil {
		return Null, EmptyObject, false
	}
	if f := o.pair; f != nil {
		if r := o.promote(); r != nil {
			return f, r, true
		}
		return f, EmptyObject, true
	}
	return Null, EmptyObject, false
}

func (o *Object) Count() Integer {
	if o == nil {
		return 0
	}
	res := Integer(1)
	for _, c := range o.children {
		if c != nil {
			res += c.Count()
		}
	}
	return res
}

func (o *Object) IsEmpty() bool {
	return o == nil
}

func (o *Object) Call(args ...Value) Value {
	return mappedCall(o, args)
}

func (o *Object) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

func (o *Object) Equal(other Value) bool {
	if other, ok := other.(*Object); ok {
		if o == nil || o == other {
			return o == other
		}
		lp := o.Pairs()
		rp := other.Pairs()
		if len(lp) != len(rp) {
			return false
		}
		rs := rp.sorted()
		for i, l := range lp.sorted() {
			if !l.Equal(rs[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (*Object) Type() types.Type {
	return types.BasicObject
}

func (o *Object) HashCode() uint64 {
	if o == nil {
		return objectHash
	}
	return o.hashCode(objectHash)
}

func (o *Object) hashCode(acc uint64) uint64 {
	h := acc ^ HashCode(o.pair.Car()) ^ HashCode(o.pair.Cdr())
	for _, c := range o.children {
		if c != nil {
			h = c.hashCode(h)
		}
	}
	return h
}

func (o *Object) Pairs() Pairs {
	if o == nil {
		return emptyPairs
	}
	return o.pairs(Pairs{})
}

func (o *Object) pairs(p Pairs) Pairs {
	p = append(p, o.pair)
	for _, c := range o.children {
		if c != nil {
			p = c.pairs(p)
		}
	}
	return p
}

func (o *Object) String() string {
	var buf strings.Builder
	buf.WriteString("{")
	for i, p := range o.Pairs().sorted() {
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(ToQuotedString(p.Car()))
		buf.WriteString(" ")
		buf.WriteString(ToQuotedString(p.Cdr()))
	}
	buf.WriteString("}")
	return buf.String()
}
