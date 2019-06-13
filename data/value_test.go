package data_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
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

	as.Falsey(data.Null)
	as.Falsey(data.False)
}
