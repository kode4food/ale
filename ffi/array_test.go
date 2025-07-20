package ffi_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
)

type ArrayTest struct {
	Bytes [3]byte `ale:"bytes"`
	Ints  [3]int  `ale:"ints"`
}

func TestArrayWrap(t *testing.T) {
	as := assert.New(t)
	a1 := [...]int{1, 2, 3}
	d1 := ffi.MustWrap(a1).(data.Vector)
	as.Equal(3, len(d1))
	as.Equal(I(1), d1[0])
	as.Equal(I(2), d1[1])
	as.Equal(I(3), d1[2])
}

func TestArrayWrapEquality(t *testing.T) {
	as := assert.New(t)
	a1 := []int{1, 2, 3}
	a2 := []int{4, 5, 6}
	a3 := []int{1, 2, 3}
	a4 := [][]int{a1, a2, a1, a3}
	w := ffi.MustWrap(a4).(data.Vector)
	as.Equal(w[0], w[2])
	as.NotEqual(w[0], w[1])
	as.Equal(w[2], w[3])
}

func TestArrayUnwrap(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(a [3]int) [3]int {
		res := [3]int{}
		for i, v := range a {
			res[i] = v * 2
		}
		return res
	}).(data.Procedure)
	out := f.Call(V(I(1), I(2), I(3))).(data.Vector)
	as.NotNil(out)
	as.Equal(3, len(out))
	as.Equal(I(2), out[0])
	as.Equal(I(4), out[1])
	as.Equal(I(6), out[2])
}

func TestArrayEval(t *testing.T) {
	as := NewWrapped(t)

	as.EvalTo(
		`[(first x) (rest x) (length x)]`,
		Env{
			"x": [...]int{10, 9, 8},
		},
		V(I(10), V(I(9), I(8)), I(3)),
	)

	as.EvalTo(
		`(d [1 2 3])`,
		Env{
			"d": func(in [3]int) (res [3]int) {
				for i, x := range in {
					res[i] = x * 2
				}
				return
			},
		},
		V(I(2), I(4), I(6)),
	)
}

func TestByteArrayUnwrap(t *testing.T) {
	as := assert.New(t)
	f1 := ffi.MustWrap(func(a [3]byte) [3]byte {
		res := [3]byte{}
		for i, v := range a {
			res[i] = v * 2
		}
		return res
	}).(data.Procedure)
	out1 := f1.Call(data.Bytes{1, 2, 3}).(data.Bytes)
	as.NotNil(out1)
	as.Equal(3, len(out1))
	as.Equal(byte(2), out1[0])
	as.Equal(byte(4), out1[1])
	as.Equal(byte(6), out1[2])

	f2 := ffi.MustWrap(func(a [3]byte) string {
		res := ""
		for _, v := range a {
			res += fmt.Sprintf("%d", v)
		}
		return res
	}).(data.Procedure)

	out2 := f2.Call(data.Bytes{1, 2, 3}).(data.String)
	as.NotNil(out2)
	as.Equal(3, len(out2))
	as.Equal("123", string(out2))

	out3 := f2.Call(S("123")).(data.String)
	as.NotNil(out3)
	as.Equal(6, len(out3))
	as.Equal("495051", string(out3))
}

func TestByteArrayStruct(t *testing.T) {
	as := NewWrapped(t)
	as.EvalTo(
		`(d {:bytes "123"})`,
		Env{
			"d": func(in ArrayTest) (res ArrayTest) {
				for i, x := range in.Bytes {
					res.Bytes[i] = x * 2
				}
				return
			},
		},
		O(
			C(K("bytes"), data.Bytes{98, 100, 102}),
			C(K("ints"), V(I(0), I(0), I(0))),
		),
	)
}
