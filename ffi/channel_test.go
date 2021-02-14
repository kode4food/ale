package ffi_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/ffi"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/internal/stream"
)

func TestChannelTypes(t *testing.T) {
	as := assert.New(t)

	both := make(chan int, 0)
	recv := (<-chan int)(both)
	send := (chan<- int)(both)

	bw := ffi.MustWrap(both)
	as.Contains(`:seq`, bw)
	as.Contains(`:emit`, bw)
	as.Contains(`:close`, bw)

	rw := ffi.MustWrap(recv)
	as.Contains(`:seq`, rw)
	as.NotContains(`:emit`, rw)
	as.Contains(`:close`, rw)

	sw := ffi.MustWrap(send)
	as.NotContains(`:seq`, sw)
	as.Contains(`:emit`, sw)
	as.Contains(`:close`, sw)

	close(both)
}

func TestChannelSequence(t *testing.T) {
	as := assert.New(t)

	ch := make(chan int, 0)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	o := ffi.MustWrap(ch).(data.Object)
	s := as.MustGet(o, stream.SequenceKey).(data.Sequence)
	as.String("(0 1 2 3 4 5 6 7 8 9)", sequence.ToList(s))
}

func TestChannelEmit(t *testing.T) {
	as := assert.New(t)

	ch := make(chan int, 0)
	go func() {
		w := ffi.MustWrap(ch).(data.Object)
		emitFunc := as.MustGet(w, stream.EmitKey).(data.Function)
		emitFunc.Call(I(1), I(2))
		emitFunc.Call(I(3), I(4))
		closeFunc := as.MustGet(w, stream.CloseKey).(data.Function)
		closeFunc.Call()
	}()

	as.Equal(1, <-ch)
	as.Equal(2, <-ch)
	as.Equal(3, <-ch)
	as.Equal(4, <-ch)

	// Check for close
	_, ok := <-ch
	as.False(ok)
}
