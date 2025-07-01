package data_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/pkg/data"
)

func TestSymbolType(t *testing.T) {
	as := assert.New(t)

	as.True(types.BasicSymbol.Equal(LS("hello").(data.Typed).Type()))
	as.True(types.BasicSymbol.Equal(QS("ale", "hello").(data.Typed).Type()))
	as.False(types.BasicBoolean.Equal(QS("ale", "hello").(data.Typed).Type()))
}

func TestLocalSymbolEquality(t *testing.T) {
	as := assert.New(t)

	sym1 := LS("hello")
	sym2 := LS("there")
	sym3 := LS("hello")

	as.True(sym1.Equal(sym1))
	as.False(sym1.Equal(sym2))
	as.True(sym1.Equal(sym3))
	as.False(sym1.Equal(I(32)))
	as.False(sym1.Equal(data.NewQualifiedSymbol("hello", "")))
}

func TestQualifiedSymbolEquality(t *testing.T) {
	as := assert.New(t)

	sym1 := data.NewQualifiedSymbol("hello", "domain")
	sym2 := data.NewQualifiedSymbol("there", "domain")
	sym3 := data.NewQualifiedSymbol("hello", "domain")

	as.True(sym1.Equal(sym1))
	as.False(sym1.Equal(sym2))
	as.True(sym1.Equal(sym3))
	as.False(sym1.Equal(I(32)))
	as.False(sym1.Equal(LS("hello")))
}

func TestSymbolParsing(t *testing.T) {
	as := assert.New(t)

	s, err := data.ParseSymbol("domain/name1")
	if as.NoError(err) {
		as.String("domain", s.(data.Qualified).Domain())
		as.String("name1", s.(data.Qualified).Name())
	}

	s, err = data.ParseSymbol("some space")
	as.Nil(s)
	as.EqualError(err, fmt.Sprintf(data.ErrInvalidSymbol, "some space"))

	s, err = data.ParseSymbol("domain/")
	as.Nil(s)
	as.EqualError(err, fmt.Sprintf(data.ErrInvalidSymbol, "domain/"))

	s, err = data.ParseSymbol("/name2")
	as.Nil(s)
	as.EqualError(err, fmt.Sprintf(data.ErrInvalidSymbol, "/name2"))

	s, err = data.ParseSymbol("one/too/")
	as.Nil(s)
	as.EqualError(err, fmt.Sprintf(data.ErrInvalidSymbol, "one/too/"))
}

func TestMustSymbolParsing(t *testing.T) {
	as := assert.New(t)

	s1 := data.MustParseSymbol("domain/name1").(data.Qualified)
	as.String("domain", s1.Domain())
	as.String("name1", s1.Name())
	as.String("domain/name1", s1)

	s2 := data.MustParseSymbol("name2").(data.Local)
	as.String("name2", s2.Name())
}

func TestSymbolGeneration(t *testing.T) {
	as := assert.New(t)

	prefix := "x-hello-gensym-%s-%s"

	gen := data.NewSymbolGenerator()
	for _, sfx := range data.SymbolGenDigits {
		s1 := gen.Local("hello")
		as.String(fmt.Sprintf(prefix, gen.Prefix(), string(sfx)), s1)
	}

	as.String(fmt.Sprintf(prefix, gen.Prefix(), "10"), gen.Local("hello"))
}

func TestSymbolHashing(t *testing.T) {
	as := assert.New(t)

	obj := O(
		C(LS("first"), I(1)),
		C(LS("second"), I(2)),
		C(QS("ale", "third"), I(3)),
	)

	as.Number(1, as.MustGet(obj, LS("first")))
	as.Number(2, as.MustGet(obj, LS("second")))
	as.Number(3, as.MustGet(obj, QS("ale", "third")))
}
