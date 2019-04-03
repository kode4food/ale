package stdlib_test

import (
	"testing"

	"gitlab.com/kode4food/ale/internal/assert"
	"gitlab.com/kode4food/ale/stdlib"
)

func TestConditionals(t *testing.T) {
	as := assert.New(t)

	i := 0
	inc := func() {
		i++
	}

	once := stdlib.Once()
	never := stdlib.Never()
	always := stdlib.Always()

	as.Integer(0, i)
	once(inc)
	as.Integer(1, i)
	once(inc)
	as.Integer(1, i)

	never(inc)
	as.Integer(1, i)
	never(inc)
	as.Integer(1, i)

	always(inc)
	as.Integer(2, i)
	always(inc)
	as.Integer(3, i)
	always(inc)
	as.Integer(4, i)
}
