package api_test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/assert"
)

func TestLocalSymbolInterning(t *testing.T) {
	as := assert.New(t)

	sym1 := api.NewLocalSymbol("hello")
	sym2 := api.NewLocalSymbol("there")
	sym3 := api.NewLocalSymbol("hello")

	as.Identical(sym1, sym3)
	as.NotIdentical(sym1, sym2)
}

func TestQualifiedSymbolInterning(t *testing.T) {
	as := assert.New(t)

	sym1 := api.NewQualifiedSymbol("hello", "domain")
	sym2 := api.NewQualifiedSymbol("there", "domain")
	sym3 := api.NewQualifiedSymbol("hello", "domain")

	as.Identical(sym1, sym3)
	as.NotIdentical(sym1, sym2)
}

func TestSymbolParsing(t *testing.T) {
	as := assert.New(t)

	s1 := api.ParseSymbol("domain/name1").(api.QualifiedSymbol)
	as.String("domain", string(s1.Domain()))
	as.String("name1", string(s1.Name()))

	s2 := api.ParseSymbol("/name2")
	if _, ok := s2.(api.QualifiedSymbol); ok {
		as.Fail("symbol should not be qualified")
	}

	s3 := api.ParseSymbol("name3")
	as.String("name3", string(s3.Name()))

	s4 := api.ParseSymbol("one/too/").(api.QualifiedSymbol)
	as.String("one", string(s4.Domain()))
	as.String("too/", string(s4.Name()))
}
