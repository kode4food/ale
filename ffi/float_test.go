package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestFloatWrapper(t *testing.T) {
	as := assert.New(t)
	f := ffi.MustWrap(func(f1 float32, f2 float64) (float32, float64) {
		return f1 * 2, f2 * 3
	}).(data.Function)
	r := f.Call(F(9), F(15)).(data.Vector)
	as.Equal(F(18), r[0])
	as.Equal(F(45), r[1])
}

func TestFloatEval(t *testing.T) {
	as := NewWrapped(t)

	as.EvalTo(
		`(d 2.5 2.4)`,
		Env{
			"d": func(f32 float32, f64 float64) (float32, float64) {
				return f32 * 2, f64 * 2
			},
		},
		V(F(5.0), F(4.8)),
	)
}
