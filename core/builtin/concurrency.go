package builtin

import (
	"slices"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/runtime"
	"github.com/kode4food/ale/internal/stream"
	"github.com/kode4food/ale/internal/sync"
)

// Go runs the provided function asynchronously
var Go = data.MakeProcedure(func(args ...ale.Value) ale.Value {
	fn := args[0].(data.Procedure)
	callArgs := slices.Clone(args[1:])
	go func() {
		defer runtime.NormalizeGoRuntimeErrors()
		fn.Call(callArgs...)
	}()
	return data.Null
}, 1, data.OrMore)

// Chan instantiates a new go channel
var Chan = data.MakeProcedure(func(args ...ale.Value) ale.Value {
	var size int
	if len(args) != 0 {
		size = int(args[0].(data.Integer))
	}
	return stream.NewChannel(size)
}, 0, 1)

// isResolved returns whether the specified promise has been resolved
func isResolved(v ale.Value) bool {
	if p, ok := v.(*sync.Promise); ok {
		return p.IsResolved()
	}
	return true
}
