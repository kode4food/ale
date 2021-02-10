package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestArrayWrap(t *testing.T) {
	as := assert.New(t)
	a1 := [3]int{1, 2, 3}
	d1 := ffi.MustWrap(a1).(data.Vector).Values()
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
	w := ffi.MustWrap(a4).(data.Vector).Values()
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
	}).(data.Function)
	out := f.Call(V(I(1), I(2), I(3))).(data.Vector).Values()
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
			"x": [3]int{10, 9, 8},
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
