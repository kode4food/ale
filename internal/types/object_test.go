package types_test

import (
	"testing"

	"github.com/kode4food/ale/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestObjectAccepts(t *testing.T) {
	as := assert.New(t)

	o1 := types.MakeObject(types.BasicString, types.BasicProcedure)
	o2 := types.MakeObject(types.BasicNumber, types.BasicKeyword)
	o3 := types.MakeObject(types.BasicString, types.BasicProcedure)

	as.Equal("object(string->procedure)", o1.Name())
	as.Equal("object(number->keyword)", o2.Name())
	as.Equal("object(string->procedure)", o3.Name())

	as.True(types.Accepts(o1, o1))
	as.False(types.Accepts(o1, o2))
	as.True(types.Accepts(o1, o3))

	as.True(types.Accepts(types.BasicObject, o1))
	as.False(types.Accepts(o1, types.BasicObject))
	as.False(types.Accepts(o1, types.BasicBoolean))
}

func TestObjectEqual(t *testing.T) {
	as := assert.New(t)

	o1 := types.MakeObject(types.BasicString, types.BasicProcedure)
	o2 := types.MakeObject(types.BasicNumber, types.BasicKeyword)
	o3 := types.MakeObject(types.BasicString, types.BasicProcedure)

	as.True(o1.Equal(o1))
	as.False(o1.Equal(o2))
	as.True(o1.Equal(o3))

	as.False(types.BasicObject.Equal(o1))
	as.False(o1.Equal(types.BasicObject))
	as.False(o1.Equal(types.BasicBoolean))
}
