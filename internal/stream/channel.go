package stream

import (
	"runtime"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/do"
	"github.com/kode4food/ale/types"
)

type (
	// Emitter is an interface that is used to emit values to a Channel
	Emitter interface {
		Writer
		Closer
		Error(any)
	}

	channelResult struct {
		value data.Value
		error any
	}

	channelEmitter struct {
		ch chan<- channelResult
	}

	channelSequence struct {
		once do.Action
		ch   <-chan channelResult

		result channelResult
		rest   data.Sequence
		ok     bool
	}
)

const (
	// EmitKey is the key used to emit to a Channel
	EmitKey = data.Keyword("emit")

	// SequenceKey is the key used to retrieve the Sequence from a Channel
	SequenceKey = data.Keyword("seq")
)

var (
	emptyResult = channelResult{value: data.Nil, error: nil}

	channelEmitterType  = types.Basic("channel-emitter")
	channelSequenceType = types.Basic("channel-sequence")
)

// NewChannel produces an Emitter and Sequence pair
func NewChannel(size int) (Emitter, data.Sequence) {
	ch := make(chan channelResult, size)
	return NewChannelEmitter(ch), NewChannelSequence(ch)
}

// NewChannelEmitter produces an Emitter for sending values to a Go chan
func NewChannelEmitter(ch chan<- channelResult) Emitter {
	r := &channelEmitter{
		ch: ch,
	}
	runtime.SetFinalizer(r, func(e *channelEmitter) {
		defer func() { recover() }()
		close(ch)
	})
	return r
}

// Write will send a Value to the Go chan
func (e *channelEmitter) Write(v data.Value) {
	e.ch <- channelResult{v, nil}
}

// Error will send an Error to the Go chan
func (e *channelEmitter) Error(err any) {
	e.ch <- channelResult{data.Nil, err}
}

// Close will Close the Go chan
func (e *channelEmitter) Close() {
	runtime.SetFinalizer(e, nil)
	close(e.ch)
}

func (e *channelEmitter) Type() types.Type {
	return channelEmitterType
}

func (e *channelEmitter) Equal(v data.Value) bool {
	return e == v
}

func (e *channelEmitter) String() string {
	return data.DumpString(e)
}

// NewChannelSequence produces a new Sequence whose values come from a Go chan
func NewChannelSequence(ch <-chan channelResult) data.Sequence {
	return &channelSequence{
		once:   do.Once(),
		ch:     ch,
		result: emptyResult,
		rest:   data.EmptyList,
	}
}

func (c *channelSequence) resolve() *channelSequence {
	c.once(func() {
		if result, ok := <-c.ch; ok {
			c.ok = ok
			c.result = result
			if c.result.error == nil {
				c.rest = NewChannelSequence(c.ch)
			}
		}
	})
	if e := c.result.error; e != nil {
		panic(e)
	}
	return c
}

func (c *channelSequence) IsEmpty() bool {
	return !c.resolve().ok
}

func (c *channelSequence) First() data.Value {
	return c.resolve().result.value
}

func (c *channelSequence) Rest() data.Sequence {
	return c.resolve().rest
}

func (c *channelSequence) Split() (data.Value, data.Sequence, bool) {
	r := c.resolve()
	return r.result.value, r.rest, r.ok
}

func (c *channelSequence) Prepend(v data.Value) data.Sequence {
	return &channelSequence{
		once:   do.Never(),
		ok:     true,
		result: channelResult{value: v, error: nil},
		rest:   c,
	}
}

func (c *channelSequence) Type() types.Type {
	return channelSequenceType
}

func (c *channelSequence) Equal(v data.Value) bool {
	return c == v
}

func (c *channelSequence) String() string {
	return data.DumpString(c)
}
