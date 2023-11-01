package data_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestMakeLambda(t *testing.T) {
	as := assert.New(t)

	f1 := data.MakeLambda(func(...data.Value) data.Value {
		return S("hello!")
	})

	as.Contains(":type lambda", f1)
	as.Nil(f1.CheckArity(99))
}

func TestLambdaEquality(t *testing.T) {
	as := assert.New(t)
	f1 := data.MakeLambda(func(...data.Value) data.Value { return nil })
	f2 := data.MakeLambda(func(...data.Value) data.Value { return nil })
	as.True(f1.Equal(f1))
	as.False(f1.Equal(f2))
	as.False(f1.Equal(I(42)))
}
