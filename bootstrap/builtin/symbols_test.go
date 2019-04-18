package builtin_test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/bootstrap/builtin"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestSymbols(t *testing.T) {
	as := assert.New(t)

	s1 := api.NewQualifiedSymbol("hello", "ale")
	as.True(builtin.IsSymbol(s1))
	as.False(builtin.IsLocal(s1))
	as.True(builtin.IsQualified(s1))

	s2 := builtin.Sym(s1)
	as.Identical(s1, s2)

	s3 := builtin.Sym(S("ale/hello"))
	as.Equal(s1, s3)

	s4 := builtin.Sym(S("howdy"))
	as.True(builtin.IsLocal(s4))
	as.False(builtin.IsQualified(s4))
}

func TestGenerated(t *testing.T) {
	as := assert.New(t)

	s1 := builtin.GenSym()
	as.True(builtin.IsSymbol(s1))
	as.True(builtin.IsLocal(s1))
	as.False(builtin.IsQualified(s1))
	as.Contains("x-anon-gensym-", s1)

	s2 := builtin.GenSym(S("blah"))
	as.Contains("x-blah-gensym-", s2)
}
