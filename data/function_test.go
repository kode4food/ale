package data_test

import (
	"testing"

	"gitlab.com/kode4food/ale/compiler/arity"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestApplicativeFunction(t *testing.T) {
	as := assert.New(t)

	f1 := data.ApplicativeFunction(func(_ ...data.Value) data.Value {
		return S("hello!")
	})

	as.True(f1.IsApplicative())
	as.False(f1.IsNormal())
	as.Contains(":type Applicative", f1)

	c1 := f1.Caller()
	as.Equal(c1, c1.Caller())
	as.String("function", c1)
	as.String("hello!", c1())

	as.Nil(f1.CheckArity(99))
}

func TestNormalFunction(t *testing.T) {
	as := assert.New(t)

	f1 := data.NormalFunction(func(_ ...data.Value) data.Value {
		return S("hello!")
	})
	f1.ArityChecker = arity.MakeFixedChecker(0)

	as.True(f1.IsNormal())
	as.False(f1.IsApplicative())
	as.Contains(":type Normal", f1)

	as.Nil(f1.CheckArity(0))
	err := f1.CheckArity(2)
	as.EqualError(err, "got 2 arguments, expected 0")
}
