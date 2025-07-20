package data_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/lang/lex"
	rdata "github.com/kode4food/ale/read/data"
)

func TestFromString(t *testing.T) {
	as := assert.New(t)

	ns := assert.GetTestNamespace()
	d1 := rdata.MustFromString(ns, `(1 2 3)`)
	as.Equal(L(I(1), I(2), I(3)), d1.Car())

	d2 := rdata.MustFromString(ns, `(#include "hello")`)
	as.Equal(L(LS("#include"), S("hello")), d2.Car())

	defer as.ExpectPanic(fmt.Errorf(lex.ErrUnexpectedCharacters, "'"))
	rdata.MustFromString(ns, `(1 2 '3)`).Car()
}
