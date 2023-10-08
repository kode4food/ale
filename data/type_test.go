package data_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestTypeCall(t *testing.T) {
	as := assert.New(t)
	l1 := data.NewList(I(1), I(2), I(3))
	pred := data.TypeOf(l1)
	as.NotNil(pred)

	l2 := data.NewList(I(9))
	v1 := data.NewVector(I(10))
	as.True(pred.Call(l1))
	as.True(pred.Call(l2))
	as.False(pred.Call(v1))
}

func TestTypeEqual(t *testing.T) {
	as := assert.New(t)
	l1 := data.NewList(I(1), I(2), I(3))
	l2 := data.NewList(I(9))
	v1 := data.NewVector(I(10))
	p1 := data.TypeOf(l1)
	p2 := data.TypeOf(l2)
	p3 := data.TypeOf(v1)
	as.True(p1.Equal(p2))
	as.False(p1.Equal(p3))
}
