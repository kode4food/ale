package data_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/data"
)

func TestBytes(t *testing.T) {
	as := assert.New(t)

	v1 := data.NewBytes(I(85), I(95), I(28), I(25))
	as.Number(4, v1.Count())
	as.Number(85, v1.Car())
	as.Number(3, v1.Cdr().(data.Counted).Count())
	as.False(v1.IsEmpty())

	r, ok := v1.ElementAt(2)
	as.True(ok)
	as.Number(28, r)
	as.String(`#b[85 95 28 25]`, v1)

	idx, ok := v1.IndexOf(I(28))
	as.True(ok)
	as.Number(2, idx)

	_, ok = v1.IndexOf(I(100))
	as.False(ok)

	v2 := v1.Append(I(14)).(data.Bytes)
	as.Number(5, v2.Count())
	as.Number(4, v1.Count())

	r, ok = v2.ElementAt(4)
	as.True(ok)
	as.Number(14, r)

	v3 := v2.Append(data.Bytes{101, 102, 103}).(data.Bytes)
	as.Equal(8, v3.Count())
	as.String(`#b[85 95 28 25 14 101 102 103]`, v3)
}

func TestBytesReverse(t *testing.T) {
	as := assert.New(t)

	as.String("#b[4 3 2 1]", data.NewBytes(I(1), I(2), I(3), I(4)).Reverse())
	as.String("#b[]", data.EmptyBytes.Reverse())
}

func TestEmptyBytes(t *testing.T) {
	as := assert.New(t)

	v := data.EmptyBytes
	as.Nil(v.Car())
	as.String("#b[]", v)
	as.String("#b[]", v.Cdr())

	as.True(v.IsEmpty())
}

func TestBytesCall(t *testing.T) {
	as := assert.New(t)

	v1 := data.NewBytes(I(85), I(95), I(28), I(25))
	as.Equal(v1, v1.Call(I(0)))
	as.Equal(v1[1:], v1.Call(I(1)))
	as.Equal(data.EmptyBytes, v1.Call(I(4)))

	as.Equal(v1[1:2], v1.Call(I(1), I(2)))
	as.Equal(v1[:2], v1.Call(I(0), I(2)))
	as.Equal(v1[3:4], v1.Call(I(3), I(4)))

	testSequenceCallInterface(as, v1)
}

func TestBytesEquality(t *testing.T) {
	as := assert.New(t)

	v1 := data.NewBytes(I(85), I(95), I(28), I(25))
	v2 := data.NewBytes(I(85), I(95), I(28), I(25))
	v3 := data.NewBytes(I(85), I(28), I(25), I(95))
	v4 := data.NewBytes(I(85), I(95), I(28))

	as.True(v1.Equal(v1))
	as.True(v1.Equal(v2))
	as.False(v1.Equal(v3))
	as.False(v1.Equal(v4))
	as.False(v1.Equal(I(32)))
}

func TestBytesSplit(t *testing.T) {
	as := assert.New(t)

	v1 := data.NewBytes(I(85), I(95), I(28), I(25))
	f, r, ok := v1.Split()

	as.True(ok)
	as.Equal(I(85), f)
	as.Equal(
		data.Bytes{95, 28, 25},
		r.(data.Bytes),
	)

	v2 := data.NewBytes(I(85))
	f, r, ok = v2.Split()
	as.True(ok)
	as.Equal(I(85), f)
	as.Equal(data.Bytes{}, r.(data.Bytes))

	v3 := data.NewBytes()
	_, _, ok = v3.Split()
	as.False(ok)
}

func TestBytesAsKey(t *testing.T) {
	as := assert.New(t)

	o1, err := data.ValuesToObject(
		V(I(85), I(20)), I(42),
		V(I(85)), I(96),
		V(I(20)), I(128),
	)

	if as.NoError(err) {
		v, ok := o1.Get(V(I(85)))
		as.True(ok)
		as.Equal(I(96), v)

		v, ok = o1.Get(V(I(85), I(20)))
		as.True(ok)
		as.Equal(I(42), v)
	}
}

func TestBytesAppendIsolation(t *testing.T) {
	as := assert.New(t)

	v1 := data.NewBytes(I(1), I(2), I(3))
	v2 := v1.Append(I(4)).(data.Bytes).Append(I(5))
	v3 := v1.Append(I(6)).(data.Bytes).Append(I(7))

	as.String("#b[1 2 3]", v1)
	as.String("#b[1 2 3 4 5]", v2)
	as.String("#b[1 2 3 6 7]", v3)
}
