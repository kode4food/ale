package encoder_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestArgs(t *testing.T) {
	as := assert.New(t)

	e := assert.GetTestEncoder()
	e.PushParams(data.Locals{"param0"}, false)
	e.PushParams(data.Locals{"param1", "param2", "param3"}, true)

	c, ok := e.ResolveParam("param2")
	as.True(ok)
	as.Equal(LS("param2"), c.Name)
	as.Equal(encoder.ValueCell, c.Type)

	c, ok = e.ResolveParam("param3")
	as.True(ok)
	as.Equal(LS("param3"), c.Name)
	as.Equal(encoder.RestCell, c.Type)

	c, ok = e.ResolveParam("param0")
	as.True(ok)
	as.Equal(LS("param0"), c.Name)
	as.Equal(encoder.ValueCell, c.Type)

	e.PopParams()
	_, ok = e.ResolveParam("param2")
	as.False(ok)
	_, ok = e.ResolveParam("param0")
	as.True(ok)

	e.PopParams()
	_, ok = e.ResolveParam("param0")
	as.False(ok)
}
