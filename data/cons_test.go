package data_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestCons(t *testing.T) {
	as := assert.New(t)
	as.String(`(1 . 2)`, C(I(1), I(2)))
	as.String(`(1 2 . 3)`, C(I(1), C(I(2), I(3))))
}

func TestConsEquality(t *testing.T) {
	as := assert.New(t)
	c1 := C(I(1), I(2))
	c2 := C(I(1), I(2))
	c3 := C(I(1), I(3))

	as.True(c1.Equal(c1))
	as.True(c1.Equal(c2))
	as.False(c1.Equal(c3))
	as.False(c1.Equal(I(2)))
}
