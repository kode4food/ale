package data_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestBools(t *testing.T) {
	as := assert.New(t)

	as.True(data.True.Equal(data.True))
	as.False(data.True.Equal(data.False))
	as.False(data.True.Equal(L()))
	as.False(data.False.Equal(L()))

	as.String("#t", data.True)
	as.String("#f", data.False)

	obj := O(
		C(data.True, S("is true")),
		C(data.False, S("is false")),
	)
	as.String("is true", as.MustGet(obj, data.True))
	as.String("is false", as.MustGet(obj, data.False))
}
