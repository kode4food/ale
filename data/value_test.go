package data_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestNames(t *testing.T) {
	as := assert.New(t)

	n := N("hello")
	as.Equal(N("hello"), n.Name())
}

func TestTruthy(t *testing.T) {
	as := assert.New(t)

	as.Truthy(data.True)
	as.Truthy(L(S("Hello")))
	as.Truthy(S("hello"))

	as.Falsey(data.Nil)
	as.Falsey(data.False)
}
