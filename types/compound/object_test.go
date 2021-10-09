package compound_test

import (
	"testing"

	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/compound"
	"github.com/stretchr/testify/assert"
)

func TestObject(t *testing.T) {
	as := assert.New(t)

	o1 := compound.Object(basic.String, basic.Lambda)
	o2 := compound.Object(basic.Number, basic.Keyword)
	o3 := compound.Object(basic.String, basic.Lambda)

	as.Equal("object(string->lambda)", o1.Name())
	as.Equal("object(number->keyword)", o2.Name())
	as.Equal("object(string->lambda)", o3.Name())

	as.True(o1.Accepts(o1))
	as.False(o1.Accepts(o2))
	as.True(o1.Accepts(o3))

	as.True(basic.Object.Accepts(o1))
	as.False(o1.Accepts(basic.Object))
	as.False(o1.Accepts(basic.Bool))
}
