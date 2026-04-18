package data

import (
	"cmp"
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
	"sync/atomic"

	"github.com/kode4food/ale"
	data "github.com/kode4food/ale/data/internal"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/types"
)

const (
	setBucketBits = 5
	setBucketMask = (1 << setBucketBits) - 1
)

// Set represents a persistent collection of unique Values
type Set struct {
	value     ale.Value
	valueHash uint64
	children  *data.SparseSlice[*Set]
	count     int
	hash      atomic.Uint64
}

var (
	// EmptySet represents an empty Set
	EmptySet *Set

	setSalt = rand.Uint64()

	// compile-time checks for interface implementation
	_ interface {
		Appender
		Counted
		Hashed
		Mapped
		Procedure
		ale.Typed
		fmt.Stringer
	} = (*Set)(nil)
)

// NewSet instantiates a new Set instance
func NewSet(vals ...ale.Value) *Set {
	res := EmptySet
	for _, v := range vals {
		res = res.Append(v).(*Set)
	}
	return res
}

// ValuesToSet interprets a set of Values as a Set
func ValuesToSet(vals ...ale.Value) *Set {
	return NewSet(vals...)
}

func (s *Set) Get(v ale.Value) (ale.Value, bool) {
	if s == nil {
		return Null, false
	}
	h := HashCode(v)
	return s.get(v, h, h)
}

func (s *Set) get(v ale.Value, vh, shifted uint64) (ale.Value, bool) {
	if s.valueHash == vh && s.value.Equal(v) {
		return s.value, true
	}

	idx := int(shifted & setBucketMask)
	if bucket, ok := s.children.Get(idx); ok {
		return bucket.get(v, vh, shifted>>setBucketBits)
	}
	return Null, false
}

func (s *Set) Append(v ale.Value) Sequence {
	h := HashCode(v)
	if s == nil {
		return &Set{
			value:     v,
			valueHash: h,
			count:     1,
		}
	}
	return s.put(v, h, h)
}

func (s *Set) put(v ale.Value, vh, shifted uint64) *Set {
	if s.valueHash == vh && s.value.Equal(v) {
		return s
	}

	idx := int(shifted & setBucketMask)
	bucket, ok := s.children.Get(idx)
	if ok {
		next := bucket.put(v, vh, shifted>>setBucketBits)
		if next == bucket {
			return s
		}
		bucket = next
	} else {
		bucket = &Set{value: v, valueHash: vh, count: 1}
	}

	children := s.children.Set(idx, bucket)
	return &Set{
		value:     s.value,
		valueHash: s.valueHash,
		children:  children,
		count:     1 + sumSetCount(children),
	}
}

func (s *Set) Remove(v ale.Value) (ale.Value, *Set, bool) {
	if s == nil {
		return Null, EmptySet, false
	}
	h := HashCode(v)
	if value, next, ok := s.remove(v, h, h); ok {
		if next != nil {
			return value, next, true
		}
		return value, EmptySet, true
	}
	return Null, s, false
}

func (s *Set) remove(v ale.Value, vh, shifted uint64) (ale.Value, *Set, bool) {
	if s.valueHash == vh && s.value.Equal(v) {
		return s.value, s.promote(), true
	}

	idx := int(shifted & setBucketMask)
	if bucket, ok := s.children.Get(idx); ok {
		if value, next, ok := bucket.remove(v, vh, shifted>>setBucketBits); ok {
			return value, s.copyWithChildAt(idx, next), true
		}
	}
	return nil, nil, false
}

func (s *Set) copyWithChildAt(idx int, child *Set) *Set {
	var children *data.SparseSlice[*Set]
	if child != nil {
		children = s.children.Set(idx, child)
	} else {
		children = s.children.Unset(idx)
	}

	return &Set{
		value:     s.value,
		valueHash: s.valueHash,
		children:  children,
		count:     1 + sumSetCount(children),
	}
}

func (s *Set) promote() *Set {
	if s == nil || s.children.IsEmpty() {
		return EmptySet
	}

	low := s.children.LowIndex()
	child, _ := s.children.Get(low)
	res := s.copyWithChildAt(low, child.promote())
	res.value = child.value
	res.valueHash = child.valueHash
	return res
}

func (s *Set) Car() ale.Value {
	if s == nil {
		return Null
	}
	return s.value
}

func (s *Set) Cdr() ale.Value {
	return s.promote()
}

func (s *Set) Split() (ale.Value, Sequence, bool) {
	if s == nil {
		return Null, EmptySet, false
	}
	return s.value, s.promote(), true
}

func (s *Set) Count() int {
	if s == nil {
		return 0
	}
	return s.count
}

func (s *Set) IsEmpty() bool {
	return s == nil
}

func (s *Set) CheckArity(argc int) error {
	return CheckRangedArity(1, 2, argc)
}

func (s *Set) Call(args ...ale.Value) ale.Value {
	res, ok := s.Get(args[0])
	if !ok && len(args) > 1 {
		return args[1]
	}
	return res
}

func (s *Set) Equal(other ale.Value) bool {
	if other, ok := other.(*Set); ok {
		if s == nil || other == nil || s == other {
			return s == other
		}
		if s.count != other.count {
			return false
		}
		lh := s.hash.Load()
		rh := other.hash.Load()
		if lh != 0 && rh != 0 && lh != rh {
			return false
		}
		return s.isIn(other)
	}
	return false
}

func (s *Set) isIn(other *Set) bool {
	if _, ok := other.Get(s.value); !ok {
		return false
	}
	for _, child := range s.childSets() {
		if !child.isIn(other) {
			return false
		}
	}
	return true
}

func (s *Set) Type() ale.Type {
	return types.MakeLiteral(types.BasicSet, s)
}

func (s *Set) HashCode() uint64 {
	if s == nil {
		return setSalt
	}
	return setSalt ^ s.hashCode()
}

func (s *Set) hashCode() uint64 {
	if h := s.hash.Load(); h != 0 {
		return h
	}
	res := s.valueHash
	for _, child := range s.childSets() {
		res ^= child.hashCode()
	}
	s.hash.Store(res)
	return res
}

func (s *Set) Members() Vector {
	if s == nil {
		return EmptyVector
	}
	return s.members(make(Vector, 0, s.count))
}

func (s *Set) members(v Vector) Vector {
	v = append(v, s.value)
	for _, child := range s.childSets() {
		v = child.members(v)
	}
	return v
}

func (s *Set) sortedMembers() Vector {
	res := s.Members()
	slices.SortFunc(res, func(left, right ale.Value) int {
		return cmp.Compare(ToString(left), ToString(right))
	})
	return res
}

func (s *Set) String() string {
	var buf strings.Builder
	buf.WriteString(lang.SetStart)
	for i, member := range s.sortedMembers() {
		if i > 0 {
			buf.WriteString(lang.Space)
		}
		buf.WriteString(ToQuotedString(member))
	}
	buf.WriteString(lang.ObjectEnd)
	return buf.String()
}

func (s *Set) childSets() []*Set {
	if s == nil {
		return nil
	}
	res, _ := s.children.RawData()
	return res
}

func sumSetCount(c *data.SparseSlice[*Set]) int {
	if c == nil {
		return 0
	}
	var res int
	raw, _ := c.RawData()
	for _, child := range raw {
		res += child.count
	}
	return res
}
