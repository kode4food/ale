package builtin

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/stdlib"
)

const (
	// ChannelKey is the key used to identify a Channel
	ChannelKey = api.Keyword("channel")

	// EmitKey is the key used to emit to a Channel
	EmitKey = api.Keyword("emit")

	// SequenceKey is the key used to retrieve the Sequence from a Channel
	SequenceKey = api.Keyword("seq")
)

var channelPrototype = api.Object{
	ChannelKey: api.True,
}

// Go runs the provided function asynchronously
func Go(args ...api.Value) api.Value {
	fn := args[0].(api.Caller)
	restArgs := args[1:]
	go fn.Caller()(restArgs...)
	return api.Nil
}

// Chan instantiates a new go channel
func Chan(_ ...api.Value) api.Value {
	e, s := stdlib.NewChannel()

	return channelPrototype.Extend(api.Object{
		EmitKey:     bindWriter(e),
		CloseKey:    bindCloser(e),
		SequenceKey: s,
	})
}

// Promise instantiates a new eventually-fulfilled promise
func Promise(args ...api.Value) api.Value {
	if len(args) == 0 {
		return stdlib.NewPromise()
	}
	p := stdlib.NewPromise()
	p.Deliver(args[0])
	return p
}

// IsPromise returns whether or not the specified value is a promise
func IsPromise(args ...api.Value) api.Value {
	_, ok := args[0].(stdlib.Promise)
	return api.Bool(ok)
}
