package main_test

import (
	"testing"

	main "github.com/kode4food/ale/cmd/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
)

func TestREPL(t *testing.T) {
	as := assert.New(t)

	r := main.NewREPL()
	as.NotNil(r)
}

func asFunction(as *assert.Wrapper, v data.Value) data.Call {
	if f, ok := v.(data.Caller); ok {
		return f.Call()
	}
	as.Fail("value is not a function")
	return nil
}

func TestBuiltInUse(t *testing.T) {
	as := assert.New(t)

	ns1 := main.GetNS()
	e, ok := ns1.Resolve("use")
	as.True(ok && e.IsBound())
	as.NotNil(e.Value())
	use := asFunction(as, e.Value())

	nsName := data.NewLocalSymbol("test-ns")
	nothing := use(nsName)
	as.NotNil(nothing)

	ns2 := main.GetNS()
	as.NotNil(ns2)
	as.String("test-ns", ns2.Domain())
}
