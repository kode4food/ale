package data_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestSet(t *testing.T) {
	as := assert.New(t)

	s1 := data.NewSet(K("parent"), K("name"))
	s2 := s1.Append(K("child")).(*data.Set).Append(K("name")).(*data.Set)

	v, ok := s2.Get(K("parent"))
	as.True(ok)
	as.Equal(K("parent"), v)

	as.Equal(2, s1.Count())
	as.Equal(3, s2.Count())

	as.Contains(":child", s2)
	as.Contains(":name", s2)
	as.Contains(":parent", s2)

	v, r, ok := s2.Remove(K("not-found"))
	as.False(ok)
	as.Nil(v)
	as.Equal(s2, r)
}

func TestSetDuplicateAppendIdentity(t *testing.T) {
	as := assert.New(t)

	s1 := data.NewSet(K("parent"), K("name"), K("child"))
	s2 := s1.Append(K("name")).(*data.Set)
	s3 := s1.Append(K("child")).(*data.Set)

	as.True(s1 == s2)
	as.True(s1 == s3)
}

func TestEmptySet(t *testing.T) {
	as := assert.New(t)

	s := data.EmptySet
	as.True(s.IsEmpty())

	v, ok := s.Get(K("word"))
	as.Nil(v)
	as.False(ok)

	v, r, ok := s.Remove(K("nothing"))
	as.Nil(v)
	as.Equal(r, s)
	as.False(ok)

	as.Nil(s.Car())
	as.Nil(s.Cdr())
}

func TestValuesToSet(t *testing.T) {
	as := assert.New(t)

	s := data.ValuesToSet()
	as.Nil(s)
	as.Number(0, s.Count())
	as.True(s.IsEmpty())

	s = data.ValuesToSet(K("kwd"), K("kwd"), S("value"))
	as.NotNil(s)
	as.Number(2, s.Count())
	as.False(s.IsEmpty())
}

func TestSetRemoval(t *testing.T) {
	as := assert.New(t)

	s1 := data.EmptySet
	for i := range 1000 {
		s1 = s1.Append(I(int64(i))).(*data.Set)
	}
	as.Equal(1000, s1.Count())

	for i := 0; i < 1000; i += 2 {
		v, r, ok := s1.Remove(I(int64(i)))
		s1 = r
		as.True(ok)
		as.Equal(I(int64(i)), v)
	}
	as.False(s1 == data.EmptySet)
	as.Equal(500, s1.Count())

	for i := 1; i < 1000; i += 2 {
		v, r, ok := s1.Remove(I(int64(i)))
		s1 = r
		as.True(ok)
		as.Equal(I(int64(i)), v)
	}
	as.True(s1 == data.EmptySet)
	as.Equal(0, s1.Count())
}

func TestSetCall(t *testing.T) {
	as := assert.New(t)

	s1 := data.NewSet(K("parent"), K("name"))

	as.Equal(K("parent"), s1.Call(K("parent")))
	as.Nil(s1.Call(K("missing")))
	as.String("defaulted", s1.Call(K("missing"), S("defaulted")))

	as.MustEvalTo(`(#{:first 1} :first)`, K("first"))
	as.MustEvalTo(`(#{:first 1} :second)`, data.Null)
	as.MustEvalTo(`(#{:first 1} :second 2)`, I(2))

	testSequenceCallInterface(as, s1)
}

func TestSetIterate(t *testing.T) {
	as := assert.New(t)

	s1 := data.NewSet(K("first"), K("second"))
	as.Equal(2, s1.Count())

	f1, r1, ok := s1.Split()
	as.True(ok)
	as.True(s1.Call(f1).Equal(f1))
	as.Equal(1, r1.(*data.Set).Count())

	f2, r2, ok := r1.Split()
	as.True(ok)
	as.True(s1.Call(f2).Equal(f2))

	_, _, ok = r2.Split()
	as.False(ok)
}

func TestSetSplitDeterminism(t *testing.T) {
	as := assert.New(t)
	s := data.NewSet(K("z"), K("x"), K("y"))
	f1, r1, ok := s.Split()
	r1Str := data.ToString(r1)
	as.True(ok)
	for range 50 {
		f2, r2, ok := s.Split()
		as.True(ok)
		as.Equal(f1, f2)
		as.Equal(r1Str, data.ToString(r2))
	}
}

func TestSetCarCdr(t *testing.T) {
	as := assert.New(t)
	s := data.NewSet(K("z"), K("x"), K("y"))
	a1 := s.Car()
	d1 := s.Cdr()
	dStr := data.ToString(d1)
	for range 50 {
		a2 := s.Car()
		d2 := s.Cdr()
		as.Equal(a1, a2)
		as.Equal(dStr, data.ToString(d2))
	}
}

func TestSetEquality(t *testing.T) {
	as := assert.New(t)
	s1 := data.NewSet(K("z"), K("x"), K("y"))
	s2 := data.NewSet(K("z"), K("x"), K("y"))
	s3 := data.NewSet(K("z"), K("y"))
	s4 := data.NewSet(K("z"), K("x"), K("y"), K("g"))
	as.True(s1.Equal(s1))
	as.True(s1.Equal(s2))
	as.False(s1.Equal(s3))
	as.False(s1.Equal(s4))
	as.False(s1.Equal(I(32)))
}

func TestSetHash(t *testing.T) {
	as := assert.New(t)
	s1 := data.NewSet(K("z"), K("x"), K("y"))
	s2 := data.NewSet(K("y"), K("x"), K("z"))
	s3 := data.NewSet(K("y"), K("z"))
	s4 := data.NewSet(K("y"))
	s5 := data.NewSet()
	as.Equal(s1.HashCode(), s2.HashCode())
	as.NotEqual(s1.HashCode(), s3.HashCode())
	as.NotEqual(uint64(0), s4.HashCode())
	as.NotEqual(uint64(0), s5.HashCode())
}
