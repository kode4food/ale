package data_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestApplicativeFunction(t *testing.T) {
	as := assert.New(t)

	f1 := data.Applicative(func(_ ...data.Value) data.Value {
		return S("hello!")
	})

	as.True(data.IsApplicative(f1))
	as.False(data.IsNormal(f1))
	as.Contains(":type applicative", f1)

	as.Nil(f1.CheckArity(99))
}

func TestNormalFunction(t *testing.T) {
	as := assert.New(t)

	f1 := data.Normal(func(_ ...data.Value) data.Value {
		return S("hello!")
	}, 0)

	as.True(data.IsNormal(f1))
	as.False(data.IsApplicative(f1))
	as.Contains(":type normal", f1)

	as.Nil(f1.CheckArity(0))
	err := f1.CheckArity(2)
	as.EqualError(err, fmt.Sprintf(data.ErrFixedArity, 0, 2))
}

func TestFunctionEquality(t *testing.T) {
	as := assert.New(t)
	f1 := data.Applicative(func(...data.Value) data.Value { return nil })
	f2 := data.Applicative(func(...data.Value) data.Value { return nil })
	as.True(f1.Equal(f1))
	as.False(f1.Equal(f2))
	as.False(f1.Equal(I(42)))
}
