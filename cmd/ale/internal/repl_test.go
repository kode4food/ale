package internal_test

import (
	"testing"

	main "github.com/kode4food/ale/cmd/ale/internal"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/compiler/encoder"
	"github.com/kode4food/ale/pkg/compiler/special"
	"github.com/kode4food/ale/pkg/data"
)

func TestREPL(t *testing.T) {
	as := assert.New(t)

	r := main.NewREPL()
	as.NotNil(r)
}

func asEncoder(t *testing.T, v data.Value) special.Call {
	t.Helper()
	if f, ok := v.(special.Call); ok {
		return f
	}
	as := assert.New(t)
	as.Fail("value is not an encoder")
	return nil
}

func TestBuiltInUse(t *testing.T) {
	as := assert.New(t)

	repl := main.NewREPL()
	ns1 := repl.GetNS()
	e, ok := ns1.Resolve("use")
	as.True(ok && e.IsBound())
	as.NotNil(e.Value())
	use := asEncoder(t, e.Value())
	nsName := LS("test-ns")
	use(encoder.NewEncoder(ns1), nsName)

	ns2 := repl.GetNS()
	as.NotNil(ns2)
	as.String("test-ns", ns2.Domain())
}
