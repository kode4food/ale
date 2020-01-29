package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/core/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/namespace"
)

func TestSymbols(t *testing.T) {
	as := assert.New(t)

	s1 := data.NewQualifiedSymbol("hello", "ale")
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

func TestResolveEval(t *testing.T) {
	as := assert.New(t)

	as.EvalTo(`(let [x 99] x)`, data.Integer(99))

	err := fmt.Errorf(namespace.ErrSymbolNotDeclared, "hello")
	as.PanicWith(`hello`, err)
	as.PanicWith(`(let [hello 99] hello) hello`, err)
}
