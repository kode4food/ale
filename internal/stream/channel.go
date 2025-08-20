package stream

import (
	"runtime"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sync"
	"github.com/kode4food/ale/internal/types"
)

type (
	chanEmitter struct {
		ch chan<- ale.Value
		cl runtime.Cleanup
	}

	chanSequence struct {
		once sync.Action
		ch   <-chan ale.Value

		result ale.Value
		rest   data.Sequence
		ok     bool
	}
)

var (
	chanSequenceType = types.MakeBasic("channel-sequence")

	// compile-time check for interface implementation
	_ data.Prepender = (*chanSequence)(nil)
)

// NewChannel produces an Emitter and Sequence pair
func NewChannel(size int) *data.Object {
	ch := make(chan ale.Value, size)
	e := newEmitter(ch)
	s := NewChannelSequence(ch)

	return data.NewObject(
		data.NewCons(EmitKey, bindWriter(e.Write)),
		data.NewCons(CloseKey, bindCloser(e)),
		data.NewCons(SequenceKey, s),
	)
}

// newEmitter produces an Emitter for sending values to a Go chan
func newEmitter(ch chan<- ale.Value) *chanEmitter {
	r := &chanEmitter{ch: ch}
	r.cl = runtime.AddCleanup(r, func(c chan<- ale.Value) {
		defer func() { _ = recover() }()
		close(c)
	}, r.ch)
	return r
}

// Write will send a Value to the Go chan
func (e *chanEmitter) Write(v ale.Value) {
	e.ch <- v
}

// Close will Close the Go chan
func (e *chanEmitter) Close() (err error) {
	defer func() { _ = recover() }()
	e.cl.Stop()
	close(e.ch)
	return nil
}

// NewChannelSequence produces a new Sequence whose values come from a Go chan
func NewChannelSequence(ch <-chan ale.Value) data.Sequence {
	return &chanSequence{
		once: sync.Once(),
		ch:   ch,
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

func (c *chanSequence) Car() ale.Value {
	return c.resolve().result
}

func (c *chanSequence) Cdr() ale.Value {
	return c.resolve().rest
}

func (c *chanSequence) Split() (ale.Value, data.Sequence, bool) {
	r := c.resolve()
	return r.result, r.rest, r.ok
}

func (c *chanSequence) Prepend(v ale.Value) data.Sequence {
	return &chanSequence{
		once:   sync.Never(),
		ok:     true,
		result: v,
		rest:   c,
	}
}

func (c *chanSequence) Type() ale.Type {
	return types.MakeLiteral(chanSequenceType, c)
}

func (c *chanSequence) Equal(other ale.Value) bool {
	return c == other
}

func (c *chanSequence) Get(key ale.Value) (ale.Value, bool) {
	return data.DumpMapped(c).Get(key)
}
