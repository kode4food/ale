package stream

import (
	"runtime"

	"github.com/kode4food/ale/internal/do"
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/pkg/data"
)

type (
	chanEmitter struct {
		ch chan<- data.Value
	}

	chanSequence struct {
		once do.Action
		ch   <-chan data.Value

		result data.Value
		rest   data.Sequence
		ok     bool
	}
)

var chanSequenceType = types.MakeBasic("channel-sequence")

// NewChannel produces an Emitter and Sequence pair
func NewChannel(size int) *data.Object {
	ch := make(chan data.Value, size)
	e := newEmitter(ch)
	s := NewChannelSequence(ch)

	return data.NewObject(
		data.NewCons(EmitKey, bindWriter(e.Write)),
		data.NewCons(CloseKey, bindCloser(e)),
		data.NewCons(SequenceKey, s),
	)
}

// newEmitter produces an Emitter for sending values to a Go chan
func newEmitter(ch chan<- data.Value) *chanEmitter {
	r := &chanEmitter{ch: ch}
	runtime.SetFinalizer(r, func(e *chanEmitter) {
		defer func() { _ = recover() }()
		close(ch)
	})
	return r
}

// Write will send a Value to the Go chan
func (e *chanEmitter) Write(v data.Value) {
	e.ch <- v
}

// Close will Close the Go chan
func (e *chanEmitter) Close() error {
	runtime.SetFinalizer(e, nil)
	close(e.ch)
	return nil
}

// NewChannelSequence produces a new Sequence whose values come from a Go chan
func NewChannelSequence(ch <-chan data.Value) data.Sequence {
	return &chanSequence{
		once:   do.Once(),
		ch:     ch,
		result: data.Null,
		rest:   data.Null,
	}
}

func (c *chanSequence) resolve() *chanSequence {
	c.once(func() {
		result, ok := <-c.ch
		if !ok {
			return
		}
		c.ok = ok
		c.result = result
		c.rest = NewChannelSequence(c.ch)
	})

	return c
}

func (c *chanSequence) IsEmpty() bool {
	return !c.resolve().ok
}

func (c *chanSequence) Car() data.Value {
	return c.resolve().result
}

func (c *chanSequence) Cdr() data.Value {
	return c.resolve().rest
}

func (c *chanSequence) Split() (data.Value, data.Sequence, bool) {
	r := c.resolve()
	return r.result, r.rest, r.ok
}

func (c *chanSequence) Prepend(v data.Value) data.Sequence {
	return &chanSequence{
		once:   do.Never(),
		ok:     true,
		result: v,
		rest:   c,
	}
}

func (c *chanSequence) Type() types.Type {
	return chanSequenceType
}

func (c *chanSequence) Equal(v data.Value) bool {
	return c == v
}

func (c *chanSequence) Get(key data.Value) (data.Value, bool) {
	return data.DumpMapped(c).Get(key)
}
