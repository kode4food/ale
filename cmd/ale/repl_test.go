package main_test

import (
	"testing"

	main "github.com/kode4food/ale/cmd/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestREPL(t *testing.T) {
	as := assert.New(t)

	r := main.NewREPL()
	as.NotNil(r)
}

func asCaller(t *testing.T, v data.Value) data.Function {
	t.Helper()
	if f, ok := v.(data.Function); ok {
		return f
	}
	as := assert.New(t)
	as.Fail("value is not a function")
	return nil
}

func TestBuiltInUse(t *testing.T) {
	as := assert.New(t)

	ns1 := main.GetNS()
	e, ok := ns1.Resolve("use")
	as.True(ok && e.IsBound())
	as.NotNil(e.Value())
	use := asCaller(t, e.Value())

	nsName := LS("test-ns")
	nothing := use.Call(nsName)
	as.NotNil(nothing)

	ns2 := main.GetNS()
	as.NotNil(ns2)
	as.String("test-ns", ns2.Domain())
}
