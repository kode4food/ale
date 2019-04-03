package main_test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	main "gitlab.com/kode4food/ale/cmd/ale"
	"gitlab.com/kode4food/ale/internal/assert"
)

func TestREPL(t *testing.T) {
	as := assert.New(t)

	r := main.NewREPL()
	as.NotNil(r)
}

func asFunction(as *assert.Wrapper, v api.Value) api.Call {
	if f, ok := v.(*api.Function); ok {
		return f.Call
	}
	as.Fail("value is not a function")
	return nil
}

func TestBuiltInUse(t *testing.T) {
	as := assert.New(t)

	ns1 := main.GetNS()
	v, ok := ns1.Resolve("use")
	as.True(ok)
	as.NotNil(v)
	use := asFunction(as, v)

	nsName := api.NewLocalSymbol("test-ns")
	nothing := use(nsName)
	as.NotNil(nothing)

	ns2 := main.GetNS()
	as.NotNil(ns2)
	as.String("test-ns", ns2.Domain())
}
