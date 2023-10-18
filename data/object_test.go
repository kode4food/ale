package data_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestObject(t *testing.T) {
	as := assert.New(t)

	o1 := data.NewObject(
		C(K("parent"), S("i am the parent")),
		C(K("name"), S("parent")),
	)

	o2 := o1.Put(
		C(K("child"), S("i am the child")),
	).(*data.Object).Put(
		C(K("name"), S("child")),
	).(*data.Object)

	as.String("i am the parent", as.MustGet(o2, K("parent")))
	as.String("child", as.MustGet(o2, K("name")))
	as.String("parent", as.MustGet(o1, K("name")))

	as.Contains(`:name "child"`, o2)
	as.Contains(`:child "i am the child"`, o2)
	as.Contains(`:parent "i am the parent"`, o2)

	v, r, ok := o2.Remove(K("not-found"))
	as.False(ok)
	as.Nil(v)
	as.Equal(o2, r)

	defer as.ExpectPanic(fmt.Sprintf(assert.ErrValueNotFound, ":missing"))
	as.MustGet(o2, K("missing"))
}

func TestEmptyObject(t *testing.T) {
	as := assert.New(t)

	o := data.EmptyObject
	as.True(o.IsEmpty())

	v, ok := o.Get(K("word"))
	as.Nil(v)
	as.False(ok)

	v, r, ok := o.Remove(K("nothing"))
	as.Nil(v)
	as.Equal(r, o)
	as.False(ok)

	as.Nil(o.Car())
	as.Nil(o.Cdr())
}

func TestValuesToObject(t *testing.T) {
	as := assert.New(t)

	o, err := data.ValuesToObject()
	as.Nil(o)
	as.Nil(err)
	as.Number(0, o.Count())
	as.True(o.IsEmpty())

	o, err = data.ValuesToObject(K("kwd"), S("value"))
	as.NotNil(o)
	as.Nil(err)
	as.Number(1, o.Count())
	as.False(o.IsEmpty())

	o, err = data.ValuesToObject(K("kwd"))
	as.Nil(o)
	as.NotNil(err)
	as.EqualError(err, data.ErrMapNotPaired)
}

func TestObjectRemoval(t *testing.T) {
	as := assert.New(t)

	// Load it
	o1 := data.EmptyObject
	for i := 0; i < 1000; i++ {
		k := K(fmt.Sprintf("key-%d", i))
		v := S(fmt.Sprintf("value-%d", i))
		o1 = o1.Put(C(k, v)).(*data.Object)
	}
	as.Equal(1000, o1.Count())

	// Remove half of it
	for i := 0; i < 1000; i += 2 {
		k := K(fmt.Sprintf("key-%d", i))
		v, r, ok := o1.Remove(k)
		o1 = r.(*data.Object)
		as.True(ok)
		as.String(fmt.Sprintf("value-%d", i), v)
	}
	as.False(o1 == data.EmptyObject)
	as.Equal(500, o1.Count())

	// Remove the other half
	for i := 1; i < 1000; i += 2 {
		k := K(fmt.Sprintf("key-%d", i))
		v, r, ok := o1.Remove(k)
		o1 = r.(*data.Object)
		as.True(ok)
		as.String(fmt.Sprintf("value-%d", i), v)
	}
	as.True(o1 == data.EmptyObject)
	as.Equal(0, o1.Count())
}

func TestObjectCall(t *testing.T) {
	as := assert.New(t)

	o1 := data.NewObject(
		C(K("parent"), S("i am the parent")),
		C(K("name"), S("parent")),
	)

	as.String("i am the parent", o1.Call(K("parent")))
	as.Nil(o1.Call(K("missing")))
	as.String("defaulted", o1.Call(K("missing"), S("defaulted")))

	as.EvalTo(`({:first 1} :first)`, I(1))
	as.EvalTo(`({:first 1} :second)`, data.Null)
	as.EvalTo(`({:first 1} :second 2)`, I(2))

	testSequenceCallInterface(as, o1)
}

func TestObjectIterate(t *testing.T) {
	as := assert.New(t)

	o1 := data.NewObject(
		C(K("first"), S("first value")),
		C(K("second"), S("second value")),
	)
	as.Equal(2, o1.Count())

	f1, r1, ok := o1.Split()
	as.True(ok)
	as.Equal(K("first"), f1.(*data.Cons).Car())
	as.Equal(S("first value"), f1.(*data.Cons).Cdr())
	as.Equal(1, r1.(*data.Object).Count())

	f2, r2, ok := r1.Split()
	as.True(ok)
	as.Equal(K("second"), f2.(*data.Cons).Car())
	as.Equal(S("second value"), f2.(*data.Cons).Cdr())

	_, _, ok = r2.Split()
	as.False(ok)
}

func TestObjectSplitDeterminism(t *testing.T) {
	as := assert.New(t)
	o := data.NewObject(
		C(K("z"), I(1024)),
		C(K("x"), I(5)),
		C(K("y"), I(99)),
	)
	f1, r1, ok := o.Split()
	r1Str := r1.String()
	as.True(ok)
	for i := 0; i < 50; i++ {
		f2, r2, ok := o.Split()
		as.True(ok)
		as.Equal(f1, f2)
		as.Equal(r1Str, r2.String())
	}
}

func TestObjectCarCdr(t *testing.T) {
	as := assert.New(t)
	o := data.NewObject(
		C(K("z"), I(1024)),
		C(K("x"), I(5)),
		C(K("y"), I(99)),
	)
	a1 := o.Car()
	d1 := o.Cdr()
	dStr := d1.String()
	for i := 0; i < 50; i++ {
		a2 := o.Car()
		d2 := o.Cdr()
		as.Equal(a1, a2)
		as.Equal(dStr, d2.String())
	}
}

func TestObjectEquality(t *testing.T) {
	as := assert.New(t)
	o1 := data.NewObject(
		C(K("z"), I(1024)),
		C(K("x"), I(5)),
		C(K("y"), I(99)),
	)
	o2 := data.NewObject( // Content same
		C(K("z"), I(1024)),
		C(K("x"), I(5)),
		C(K("y"), I(99)),
	)
	o3 := data.NewObject( // Missing key
		C(K("z"), I(1024)),
		C(K("y"), I(99)),
	)
	o4 := data.NewObject( // Additional Key
		C(K("z"), I(1024)),
		C(K("x"), I(5)),
		C(K("y"), I(99)),
		C(K("g"), I(1024)),
	)
	o5 := data.NewObject( // Modified Value in x
		C(K("z"), I(1024)),
		C(K("x"), I(6)),
		C(K("y"), I(99)),
	)
	as.True(o1.Equal(o1))
	as.True(o1.Equal(o2))
	as.False(o1.Equal(o3))
	as.False(o1.Equal(o4))
	as.False(o1.Equal(o5))
	as.False(o1.Equal(I(32)))
}

func TestObjectHash(t *testing.T) {
	as := assert.New(t)
	o1 := data.NewObject(
		C(K("z"), I(1024)),
		C(K("x"), I(5)),
		C(K("y"), I(99)),
	)
	o2 := data.NewObject(
		C(K("y"), I(99)),
		C(K("x"), I(5)),
		C(K("z"), I(1024)),
	)
	o3 := data.NewObject(
		C(K("y"), I(99)),
		C(K("z"), I(1024)),
	)
	o4 := data.NewObject(
		C(K("y"), I(99)),
	)
	o5 := data.NewObject()
	as.Equal(o1.HashCode(), o2.HashCode())
	as.NotEqual(o1.HashCode(), o3.HashCode())
	as.NotEqual(uint64(0), o4.HashCode())
	as.NotEqual(uint64(0), o5.HashCode())
}
