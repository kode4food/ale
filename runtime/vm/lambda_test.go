package vm_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/ale/runtime/vm"
)

func TestLambdaFromEncoder(t *testing.T) {
	as := assert.New(t)

	e1 := assert.GetTestEncoder()
	e1.Emit(isa.ArgLen)
	e1.Emit(isa.Return)

	l := vm.LambdaFromEncoder(e1)
	as.NotNil(l)

	c, ok := l.Call().(data.Function)
	as.True(ok)
	as.NotNil(c)

	as.Equal(I(4), c.Call(S("one"), S("two"), S("three"), S("four")))
	as.Contains(":type lambda", c)
}
