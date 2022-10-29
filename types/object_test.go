package types_test

import (
	"testing"

	"github.com/kode4food/ale/types"
	"github.com/stretchr/testify/assert"
)

func TestObject(t *testing.T) {
	as := assert.New(t)

	o1 := types.Object(types.String, types.Lambda)
	o2 := types.Object(types.Number, types.Keyword)
	o3 := types.Object(types.String, types.Lambda)

	as.Equal("object(string->lambda)", o1.Name())
	as.Equal("object(number->keyword)", o2.Name())
	as.Equal("object(string->lambda)", o3.Name())

	as.True(types.Accepts(o1, o1))
	as.False(types.Accepts(o1, o2))
	as.True(types.Accepts(o1, o3))

	as.True(types.Accepts(types.AnyObject, o1))
	as.False(types.Accepts(o1, types.AnyObject))
	as.False(types.Accepts(o1, types.Bool))
}
