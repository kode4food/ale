package stream

import (
	"runtime"
	"sync/atomic"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/do"
)

type (
	// Emitter is an interface that is used to emit values to a Channel
	Emitter interface {
		Writer
		Closer
		Error(interface{})
	}

	channelResult struct {
		value data.Value
		error interface{}
	}

	channelWrapper struct {
		seq    chan channelResult
		status uint32
	}

	channelEmitter struct {
		ch *channelWrapper
	}

	channelSequence struct {
		once do.Action
		ch   *channelWrapper

		result channelResult
		rest   data.Sequence
		ok     bool
	}
)

const (
	// ChannelType is the type name for a channel
	ChannelType = data.String("channel")

	// EmitKey is the key used to emit to a Channel
	EmitKey = data.Keyword("emit")

	// SequenceKey is the key used to retrieve the Sequence from a Channel
	SequenceKey = data.Keyword("seq")
)

const (
	channelReady uint32 = iota
	channelCloseRequested
	channelClosed
)

var emptyResult = channelResult{value: data.Nil, error: nil}

func (ch *channelWrapper) Close() {
	if atomic.LoadUint32(&ch.status) != channelClosed {
		atomic.StoreUint32(&ch.status, channelClosed)
		close(ch.seq)
	}
}

// NewChannel produces a Emitter and Sequence pair
func NewChannel(size int) (Emitter, data.Sequence) {
	seq := make(chan channelResult, size)
	ch := &channelWrapper{
		seq:    seq,
		status: channelReady,
	}
	return NewChannelEmitter(ch), NewChannelSequence(ch)
}

// NewChannelEmitter produces an Emitter for sending values to a Go chan
func NewChannelEmitter(ch *channelWrapper) Emitter {
	r := &channelEmitter{
		ch: ch,
	}
	runtime.SetFinalizer(r, func(e *channelEmitter) {
		defer func() { recover() }()
		if atomic.LoadUint32(&ch.status) != channelClosed {
			e.Close()
		}
	})
	return r
}

// Write will send a Value to the Go chan
func (e *channelEmitter) Write(v data.Value) {
	if atomic.LoadUint32(&e.ch.status) == channelReady {
		e.ch.seq <- channelResult{v, nil}
	}
	if atomic.LoadUint32(&e.ch.status) == channelCloseRequested {
		e.Close()
	}
}

// Error will send an Error to the Go chan
func (e *channelEmitter) Error(err interface{}) {
	if atomic.LoadUint32(&e.ch.status) == channelReady {
		e.ch.seq <- channelResult{data.Nil, err}
	}
	e.Close()
}

// Close will Close the Go chan
func (e *channelEmitter) Close() {
	runtime.SetFinalizer(e, nil)
	e.ch.Close()
}

func (e *channelEmitter) Type() data.Name {
	return "channel-emitter"
}

func (e *channelEmitter) Equal(v data.Value) bool {
	if v, ok := v.(*channelEmitter); ok {
		return e == v
	}
	return false
}

func (e *channelEmitter) String() string {
	return data.DumpString(e)
}

// NewChannelSequence produces a new Sequence whose values come from a Go chan
func NewChannelSequence(ch *channelWrapper) data.Sequence {
	r := &channelSequence{
		once:   do.Once(),
		ch:     ch,
		result: emptyResult,
		rest:   data.EmptyList,
	}
	runtime.SetFinalizer(r, func(c *channelSequence) {
		defer func() { recover() }()
		if atomic.LoadUint32(&c.ch.status) == channelReady {
			atomic.StoreUint32(&c.ch.status, channelCloseRequested)
			<-ch.seq // consume whatever is there
		}
	})
	return r
}

func (c *channelSequence) resolve() *channelSequence {
	c.once(func() {
		runtime.SetFinalizer(c, nil)
		ch := c.ch
		if result, ok := <-ch.seq; ok {
			c.ok = ok
			c.result = result
			c.rest = NewChannelSequence(ch)
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

func (c *channelSequence) Type() data.Name {
	return "channel-sequence"
}

func (c *channelSequence) Equal(v data.Value) bool {
	if v, ok := v.(*channelSequence); ok {
		return c == v
	}
	return false
}

func (c *channelSequence) String() string {
	return data.DumpString(c)
}
