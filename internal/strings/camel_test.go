package strings_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/strings"
)

func TestCamelToSnake(t *testing.T) {
	as := assert.New(t)

	as.Equal("this-was-camel", strings.CamelToSnake("thisWasCamel"))
	as.Equal("this-was-camel", strings.CamelToSnake("ThisWasCamel"))
	as.Equal("this-was-a-camel", strings.CamelToSnake("thisWasACamel"))
}

func TestCamelToWords(t *testing.T) {
	as := assert.New(t)

	as.Equal("this was camel", strings.CamelToWords("thisWasCamel"))
	as.Equal("this was camel", strings.CamelToWords("ThisWasCamel"))
	as.Equal("this was a camel", strings.CamelToWords("thisWasACamel"))
}
