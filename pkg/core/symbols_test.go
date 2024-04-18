package core_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/core"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

func TestSymbols(t *testing.T) {
	as := assert.New(t)

	s1 := data.NewQualifiedSymbol("hello", "ale")
	as.True(getPredicate(core.SymbolKey).Call(s1))
	as.False(getPredicate(core.LocalKey).Call(s1))
	as.True(getPredicate(core.QualifiedKey).Call(s1))

	s2 := core.Sym.Call(s1)
	as.Identical(s1, s2)

	s3 := core.Sym.Call(S("ale/hello"))
	as.Equal(s1, s3)

	s4 := core.Sym.Call(S("howdy"))
	as.True(getPredicate(core.LocalKey).Call(s4))
	as.False(getPredicate(core.QualifiedKey).Call(s4))
}

func TestGenerated(t *testing.T) {
	as := assert.New(t)

	s1 := core.GenSym.Call()
	as.True(getPredicate(core.SymbolKey).Call(s1))
	as.True(getPredicate(core.LocalKey).Call(s1))
	as.False(getPredicate(core.QualifiedKey).Call(s1))
	as.Contains("x-anon-gensym-", s1)

	s2 := core.GenSym.Call(LS("blah"))
	as.Contains("x-blah-gensym-", s2)
}

func TestResolveEval(t *testing.T) {
	as := assert.New(t)

	as.EvalTo(`(let [x 99] x)`, I(99))

	err := fmt.Errorf(env.ErrSymbolNotDeclared, "hello")
	as.PanicWith(`hello`, err)
	as.PanicWith(`(let [hello 99] hello) hello`, err)
}
