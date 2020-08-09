package do_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/do"
)

func TestConditionals(t *testing.T) {
	as := assert.New(t)

	i := 0
	inc := func() {
		i++
	}

	once := do.Once()
	never := do.Never()
	always := do.Always()

	as.Number(0, i)
	once(inc)
	as.Number(1, i)
	once(inc)
	as.Number(1, i)

	never(inc)
	as.Number(1, i)
	never(inc)
	as.Number(1, i)

	always(inc)
	as.Number(2, i)
	always(inc)
	as.Number(3, i)
	always(inc)
	as.Number(4, i)
}
