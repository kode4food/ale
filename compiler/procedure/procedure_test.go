package procedure_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/procedure"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/runtime/isa"
)

func TestFromEncoder(t *testing.T) {
	as := assert.New(t)

	e1 := assert.GetTestEncoder()
	e1.Emit(isa.ArgLen)
	e1.Emit(isa.Return)

	l := procedure.FromEncoder(e1)
	as.NotNil(l)
	as.Nil(l.CheckArity(-1))

	c, ok := l.Call().(data.Procedure)
	as.True(ok)
	as.NotNil(c)

	as.Equal(I(4), c.Call(S("one"), S("two"), S("three"), S("four")))
	as.Contains(":type procedure", c)
}
