package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/core/builtin"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

func TestSymbols(t *testing.T) {
	as := assert.New(t)

	s1 := data.NewQualifiedSymbol("hello", "ale")
	as.True(getPredicate(builtin.SymbolKey).Call(s1))
	as.False(getPredicate(builtin.LocalKey).Call(s1))
	as.True(getPredicate(builtin.QualifiedKey).Call(s1))

	s2 := builtin.Sym.Call(s1)
	as.Identical(s1, s2)

	s3 := builtin.Sym.Call(S("ale/hello"))
	as.Equal(s1, s3)

	s4 := builtin.Sym.Call(S("howdy"))
	as.True(getPredicate(builtin.LocalKey).Call(s4))
	as.False(getPredicate(builtin.QualifiedKey).Call(s4))
}

func TestGenerated(t *testing.T) {
	as := assert.New(t)

	s1 := builtin.GenSym.Call()
	as.True(getPredicate(builtin.SymbolKey).Call(s1))
	as.True(getPredicate(builtin.LocalKey).Call(s1))
	as.False(getPredicate(builtin.QualifiedKey).Call(s1))
	as.Contains("x-anon-gensym-", s1)

	s2 := builtin.GenSym.Call(LS("blah"))
	as.Contains("x-blah-gensym-", s2)
}

func TestResolveEval(t *testing.T) {
	as := assert.New(t)

	as.MustEvalTo(`(let [x 99] x)`, I(99))

	err := fmt.Errorf(env.ErrNameNotDeclared, "hello")
	as.ErrorWith(`hello`, err)
	as.ErrorWith(`(let [hello 99] hello) hello`, err)
}
