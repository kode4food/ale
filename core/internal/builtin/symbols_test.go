package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/core/internal/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestSymbols(t *testing.T) {
	as := assert.New(t)

	s1 := data.NewQualifiedSymbol("hello", "ale")
	as.True(builtin.IsSymbol.Call(s1))
	as.False(builtin.IsLocal.Call(s1))
	as.True(builtin.IsQualified.Call(s1))

	s2 := builtin.Sym.Call(s1)
	as.Identical(s1, s2)

	s3 := builtin.Sym.Call(S("ale/hello"))
	as.Equal(s1, s3)

	s4 := builtin.Sym.Call(S("howdy"))
	as.True(builtin.IsLocal.Call(s4))
	as.False(builtin.IsQualified.Call(s4))
}

func TestGenerated(t *testing.T) {
	as := assert.New(t)

	s1 := builtin.GenSym.Call()
	as.True(builtin.IsSymbol.Call(s1))
	as.True(builtin.IsLocal.Call(s1))
	as.False(builtin.IsQualified.Call(s1))
	as.Contains("x-anon-gensym-", s1)

	s2 := builtin.GenSym.Call(S("blah"))
	as.Contains("x-blah-gensym-", s2)
}

func TestResolveEval(t *testing.T) {
	as := assert.New(t)

	as.EvalTo(`(let [x 99] x)`, I(99))

	err := fmt.Errorf(env.ErrSymbolNotDeclared, "hello")
	as.PanicWith(`hello`, err)
	as.PanicWith(`(let [hello 99] hello) hello`, err)
}
