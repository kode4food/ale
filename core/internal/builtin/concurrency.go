package builtin

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/async"
	"github.com/kode4food/ale/internal/stream"
)

const (
	// ChannelType is the type name for a channel
	ChannelType = data.String("channel")

	// EmitKey is the key used to emit to a Channel
	EmitKey = data.Keyword("emit")

	// SequenceKey is the key used to retrieve the Sequence from a Channel
	SequenceKey = data.Keyword("seq")
)

// Go runs the provided function asynchronously
var Go = data.Applicative(func(args ...data.Value) data.Value {
	fn := args[0].(data.Function)
	restArgs := args[1:]
	go fn.Call(restArgs...)
	return data.Nil
}, 1)

// Chan instantiates a new go channel
var Chan = data.Applicative(func(args ...data.Value) data.Value {
	var size int
	if len(args) != 0 {
		size = int(args[0].(data.Integer))
	}
	e, s := stream.NewChannel(size)

	return data.NewObject(
		data.NewCons(data.TypeKey, ChannelType),
		data.NewCons(EmitKey, bindWriter(e)),
		data.NewCons(CloseKey, bindCloser(e)),
		data.NewCons(SequenceKey, s),
	)
}, 0, 1)

// Promise instantiates a new eventually-fulfilled promise
var Promise = data.Applicative(func(args ...data.Value) data.Value {
	resolver := args[0].(data.Function)
	return async.NewPromise(resolver)
}, 1)

// IsPromise returns whether the specified value is a promise
var IsPromise = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(async.Promise)
	return data.Bool(ok)
}, 1)

// IsResolved returns whether the specified promise has been resolved
var IsResolved = data.Applicative(func(args ...data.Value) data.Value {
	p := args[0].(async.Promise)
	return data.Bool(p.IsResolved())
}, 1)
