package encoder_test

import (
	"testing"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestArgs(t *testing.T) {
	as := NewWrapped(t)

	e := getTestEncoder()
	e.PushArgs(data.Names{"arg0"}, false)
	e.PushArgs(data.Names{"arg1", "arg2", "arg3"}, true)

	c, ok := e.ResolveArg("arg2")
	as.True(ok)
	as.Equal(N("arg2"), c.Name)
	as.Equal(encoder.ValueCell, c.Type)

	c, ok = e.ResolveArg("arg3")
	as.True(ok)
	as.Equal(N("arg3"), c.Name)
	as.Equal(encoder.RestCell, c.Type)

	c, ok = e.ResolveArg("arg0")
	as.True(ok)
	as.Equal(N("arg0"), c.Name)
	as.Equal(encoder.ValueCell, c.Type)

	e.PopArgs()
	_, ok = e.ResolveArg("arg2")
	as.False(ok)
	_, ok = e.ResolveArg("arg0")
	as.True(ok)

	e.PopArgs()
	_, ok = e.ResolveArg("arg0")
	as.False(ok)
}
