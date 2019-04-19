package data_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
)

func TestLocalSymbolInterning(t *testing.T) {
	as := assert.New(t)

	sym1 := data.NewLocalSymbol("hello")
	sym2 := data.NewLocalSymbol("there")
	sym3 := data.NewLocalSymbol("hello")

	as.Identical(sym1, sym3)
	as.NotIdentical(sym1, sym2)
}

func TestQualifiedSymbolInterning(t *testing.T) {
	as := assert.New(t)

	sym1 := data.NewQualifiedSymbol("hello", "domain")
	sym2 := data.NewQualifiedSymbol("there", "domain")
	sym3 := data.NewQualifiedSymbol("hello", "domain")

	as.Identical(sym1, sym3)
	as.NotIdentical(sym1, sym2)
}

func TestSymbolParsing(t *testing.T) {
	as := assert.New(t)

	s1 := data.ParseSymbol("domain/name1").(data.QualifiedSymbol)
	as.String("domain", string(s1.Domain()))
	as.String("name1", string(s1.Name()))

	s2 := data.ParseSymbol("/name2")
	if _, ok := s2.(data.QualifiedSymbol); ok {
		as.Fail("symbol should not be qualified")
	}

	s3 := data.ParseSymbol("name3")
	as.String("name3", string(s3.Name()))

	s4 := data.ParseSymbol("one/too/").(data.QualifiedSymbol)
	as.String("one", string(s4.Domain()))
	as.String("too/", string(s4.Name()))
}
