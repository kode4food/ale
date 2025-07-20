package internal_test

import (
	"testing"

	"github.com/kode4food/ale"
	main "github.com/kode4food/ale/cmd/ale/internal"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/internal/compiler/encoder"
)

func TestREPL(t *testing.T) {
	as := assert.New(t)

	r := main.NewREPL()
	as.NotNil(r)
}

func asEncoder(t *testing.T, v ale.Value) compiler.Call {
	t.Helper()
	if f, ok := v.(compiler.Call); ok {
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
	v := as.IsBound(ns1, "use")
	use := asEncoder(t, v)
	nsName := LS("test-ns")
	as.Nil(use(encoder.NewEncoder(ns1), nsName))

	ns2 := repl.GetNS()
	as.NotNil(ns2)
	as.String("test-ns", ns2.Domain())
}
