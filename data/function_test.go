package data_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/compiler/arity"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestApplicativeFunction(t *testing.T) {
	as := assert.New(t)

	f1 := data.MakeApplicative(func(_ ...data.Value) data.Value {
		return S("hello!")
	}, nil)

	as.True(data.IsApplicative(f1))
	as.False(data.IsNormal(f1))
	as.Contains(":type applicative", f1)

	as.Nil(f1.CheckArity(99))
}

func TestNormalFunction(t *testing.T) {
	as := assert.New(t)

	f1 := data.MakeNormal(func(_ ...data.Value) data.Value {
		return S("hello!")
	}, arity.MakeFixedChecker(0))

	as.True(data.IsNormal(f1))
	as.False(data.IsApplicative(f1))
	as.Contains(":type normal", f1)

	as.Nil(f1.CheckArity(0))
	err := f1.CheckArity(2)
	as.EqualError(err, fmt.Sprintf(arity.ErrFixedArity, 0, 2))
}
