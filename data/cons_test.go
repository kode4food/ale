package data_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestCons(t *testing.T) {
	as := assert.New(t)
	as.String(`(1 . 2)`, data.NewCons(I(1), I(2)))
	as.String(`(1 2 . 3)`, data.NewCons(I(1), data.NewCons(I(2), I(3))))
}

func TestConsEquality(t *testing.T) {
	as := assert.New(t)
	c1 := data.NewCons(I(1), I(2))
	c2 := data.NewCons(I(1), I(2))
	c3 := data.NewCons(I(1), I(3))

	as.True(c1.Equal(c1))
	as.True(c1.Equal(c2))
	as.False(c1.Equal(c3))
	as.False(c1.Equal(I(2)))
}

func TestConsStringify(t *testing.T) {
	as := assert.New(t)
	c1 := data.NewCons(I(1), V(I(2), I(3), I(4)))
	c2 := data.NewCons(I(1), L(I(2), I(3), I(4)))
	c3 := data.NewCons(I(1), I(2))
	c4 := data.NewCons(I(1), data.EmptyList)
	as.String("(1 2 3 4)", c1)
	as.String("(1 2 3 4)", c2)
	as.String("(1 . 2)", c3)
	as.String("(1)", c4)
}
