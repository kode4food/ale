package stream_test

import (
	"sync"
	"testing"
	"time"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/stream"
)

func TestChannel(t *testing.T) {
	as := assert.New(t)

	ch := stream.NewChannel(0)

	var emit data.Function
	v, _ := ch.Get(stream.EmitKey)
	emit = v.(data.Function)

	var cl data.Function
	v, _ = ch.Get(stream.CloseKey)
	cl = v.(data.Function)

	var seq data.Sequence
	v, _ = ch.Get(stream.SequenceKey)
	seq = v.(data.Prepender).Prepend(F(1))
	as.Contains(":type channel-sequence", seq)

	var wg sync.WaitGroup

	gen := func() {
		emit.Call(F(2))
		time.Sleep(time.Millisecond * 50)
		emit.Call(F(3))
		time.Sleep(time.Millisecond * 30)
		emit.Call(S("foo"))
		time.Sleep(time.Millisecond * 10)
		emit.Call(S("bar"))
		cl.Call()
		wg.Done()
	}

	check := func() {
		f, _, ok := seq.Split()
		as.Number(1, f)
		as.True(ok)

		as.Number(1, seq.Car())
		as.Number(2, seq.Cdr().(data.Pair).Car())
		as.Number(3, seq.Cdr().(data.Pair).Cdr().(data.Pair).Car())
		as.False(seq.Cdr().(data.Pair).Cdr().(data.Pair).Cdr().(data.Sequence).
			IsEmpty())
		as.String("foo", seq.Cdr().(data.Pair).Cdr().(data.Pair).
			Cdr().(data.Pair).Car())
		as.False(seq.Cdr().(data.Pair).Cdr().(data.Pair).Cdr().(data.Pair).
			Cdr().(data.Sequence).IsEmpty())
		as.String("bar", seq.Cdr().(data.Pair).Cdr().(data.Pair).
			Cdr().(data.Pair).Cdr().(data.Pair).Car())
		as.True(seq.Cdr().(data.Pair).Cdr().(data.Pair).Cdr().(data.Pair).
			Cdr().(data.Pair).Cdr().(data.Sequence).IsEmpty())
		wg.Done()
	}

	wg.Add(4)
	go check()
	go check()
	go gen()
	go check()
	wg.Wait()
}
