package data_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/read/data"
	"github.com/kode4food/ale/read/lex"
)

func TestFromString(t *testing.T) {
	as := assert.New(t)

	d1 := data.FromString(`(1 2 3)`)
	as.Equal(L(I(1), I(2), I(3)), d1.Car())

	defer as.ExpectPanic(fmt.Errorf(lex.ErrUnexpectedCharacters, "'"))
	data.FromString(`(1 2 '3)`).Car()
}
