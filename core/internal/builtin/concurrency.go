package builtin

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/async"
	"github.com/kode4food/ale/internal/stream"
	"github.com/kode4food/ale/runtime"
)

// Go runs the provided function asynchronously
var Go = data.Applicative(func(args ...data.Value) data.Value {
	fn := args[0].(data.Function)
	restArgs := args[1:]
	go func() {
		defer runtime.NormalizeGoRuntimeErrors()
		fn.Call(restArgs...)
	}()
	return data.Null
}, 1)

// Chan instantiates a new go channel
var Chan = data.Applicative(func(args ...data.Value) data.Value {
	var size int
	if len(args) != 0 {
		size = int(args[0].(data.Integer))
	}
	return stream.NewChannel(size)

}, 0, 1)

// Promise instantiates a new eventually-fulfilled promise
var Promise = data.Applicative(func(args ...data.Value) data.Value {
	resolver := args[0].(data.Function)
	return async.NewPromise(resolver)
}, 1)

// isResolved returns whether the specified promise has been resolved
func isResolved(v data.Value) bool {
	if p, ok := v.(*async.Promise); ok {
		return p.IsResolved()
	}
	return true
}
