package data_test

import (
	"testing"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestMakeProcedure(t *testing.T) {
	as := assert.New(t)

	f1 := data.MakeProcedure(func(...ale.Value) ale.Value {
		return S("hello!")
	})

	as.Contains(":type procedure", f1)
	as.NoError(f1.CheckArity(99))
}

func TestProcedureEquality(t *testing.T) {
	as := assert.New(t)
	f1 := data.MakeProcedure(func(...ale.Value) ale.Value { return nil })
	f2 := data.MakeProcedure(func(...ale.Value) ale.Value { return nil })
	as.True(f1.Equal(f1))
	as.False(f1.Equal(f2))
	as.False(f1.Equal(I(42)))
}
