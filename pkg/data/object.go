package data

import (
	"errors"
	"math/rand"
	"strings"
	"sync/atomic"

	"github.com/kode4food/ale/internal/data"
	"github.com/kode4food/ale/internal/types"
)

// Object maps a set of Values, known as keys, to another set of Values
type Object struct {
	pair     Pair
	keyHash  uint64
	children *data.SparseSlice[*Object]
	hash     uint64
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
	} = (*Object)(nil)
)

// NewObject instantiates a new Object instance. Based on Phil Bagwell's Hashed
// Array Mapped Trie data structure. More information on HAMT's can be found at
// http://lampwww.epfl.ch/papers/idealhashtrees.pdf
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
	return o.get(k, h, h)
}

func (o *Object) get(k Value, kh, shifted uint64) (Value, bool) {
	if o.keyHash == kh && o.pair.Car().Equal(k) {
		return o.pair.Cdr(), true
	}

	idx := int(shifted & bucketMask)
	if bucket, ok := o.children.Get(idx); ok {
		return bucket.get(k, kh, shifted>>bucketBits)
	}
	return Null, false
}

func (o *Object) Put(p Pair) Sequence {
	h := HashCode(p.Car())
	if o == nil {
		return &Object{
			pair:    p,
			keyHash: h,
		}
	}
	return o.put(p, h, h)
}

func (o *Object) put(p Pair, kh, shifted uint64) *Object {
	if o.keyHash == kh && o.pair.Car().Equal(p.Car()) {
		return &Object{
			pair:     p,
			keyHash:  kh,
			children: o.children,
		}
	}

	idx := int(shifted & bucketMask)
	bucket, ok := o.children.Get(idx)
	if ok {
		bucket = bucket.put(p, kh, shifted>>bucketBits)
	} else {
		bucket = &Object{pair: p, keyHash: kh}
	}

	// return a copy with the new bucket
	res := *o
	res.children = res.children.Set(idx, bucket)
	res.hash = 0
	return &res
}

func (o *Object) Remove(k Value) (Value, Sequence, bool) {
	if o == nil {
		return Null, EmptyObject, false
	}
	h := HashCode(k)
	if v, r, ok := o.remove(k, h, h); ok {
		if r != nil {
			return v, r, true
		}
		return v, EmptyObject, true
	}
	return Null, o, false
}

func (o *Object) remove(k Value, kh, shifted uint64) (Value, *Object, bool) {
	if o.keyHash == kh && o.pair.Car().Equal(k) {
		return o.pair.Cdr(), o.promote(), true
	}

	idx := int(shifted & bucketMask)
	if bucket, ok := o.children.Get(idx); ok {
		if v, r, ok := bucket.remove(k, kh, shifted>>bucketBits); ok {
			return v, o.copyWithChildAt(idx, r), true
		}
	}
	return nil, nil, false
}

func (o *Object) copyWithChildAt(idx int, child *Object) *Object {
	res := *o
	res.hash = 0
	if child != nil {
		res.children = res.children.Set(idx, child)
	} else {
		res.children = res.children.Unset(idx)
	}
	return &res
}

func (o *Object) promote() *Object {
	if o == nil || o.children.IsEmpty() {
		return EmptyObject
	}

	low := o.children.LowIndex()
	c, _ := o.children.Get(low)
	res := o.copyWithChildAt(low, c.promote())
	res.pair = c.pair
	res.keyHash = c.keyHash
	return res
}

func (o *Object) Car() Value {
	if o == nil {
		return Null
	}
	return o.pair
}

func (o *Object) Cdr() Value {
	return o.promote()
}

func (o *Object) Split() (Value, Sequence, bool) {
	if o == nil {
		return Null, EmptyObject, false
	}
	return o.pair, o.promote(), true
}

func (o *Object) Count() Integer {
	if o == nil {
		return 0
	}
	res := Integer(1)
	for _, c := range o.childObjects() {
		res += c.Count()
	}
	return res
}

func (o *Object) IsEmpty() bool {
	return o == nil
}

func (o *Object) Call(args ...Value) Value {
	return mappedCall(o, args)
}

func (o *Object) CheckArity(argc int) error {
	return CheckRangedArity(1, 2, argc)
}

func (o *Object) Equal(other Value) bool {
	if other, ok := other.(*Object); ok {
		if o == nil || other == nil || o == other {
			return o == other
		}
		lh := atomic.LoadUint64(&o.hash)
		rh := atomic.LoadUint64(&other.hash)
		if lh != 0 && rh != 0 && lh != rh {
			return false
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
	return o.hashCode() ^ objectHash
}

func (o *Object) hashCode() uint64 {
	if h := atomic.LoadUint64(&o.hash); h != 0 {
		return h
	}
	res := o.keyHash ^ HashCode(o.pair.Cdr())
	for _, c := range o.childObjects() {
		res ^= c.hashCode()
	}
	atomic.StoreUint64(&o.hash, res)
	return res
}

func (o *Object) Pairs() Pairs {
	if o == nil {
		return emptyPairs
	}
	return o.pairs(Pairs{})
}

func (o *Object) pairs(p Pairs) Pairs {
	p = append(p, o.pair)
	for _, c := range o.childObjects() {
		p = c.pairs(p)
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

func (o *Object) childObjects() []*Object {
	if o == nil {
		return nil
	}
	res, _ := o.children.RawData()
	return res
}
