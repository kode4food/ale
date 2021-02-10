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

	e, seq := stream.NewChannel(0)
	seq = seq.(data.PrependerSequence).Prepend(F(1))
	as.Contains(":type channel-emitter", e)
	as.Contains(":type channel-sequence", seq)

	var wg sync.WaitGroup

	gen := func() {
		e.Write(F(2))
		time.Sleep(time.Millisecond * 50)
		e.Write(F(3))
		time.Sleep(time.Millisecond * 30)
		e.Write(S("foo"))
		time.Sleep(time.Millisecond * 10)
		e.Write(S("bar"))
		e.Close()
		wg.Done()
	}

	check := func() {
		f, _, ok := seq.Split()
		as.Number(1, f)
		as.True(ok)

		as.Number(1, seq.First())
		as.Number(2, seq.Rest().First())
		as.Number(3, seq.Rest().Rest().First())
		as.False(seq.Rest().Rest().Rest().IsEmpty())
		as.String("foo", seq.Rest().Rest().Rest().First())
		as.False(seq.Rest().Rest().Rest().Rest().IsEmpty())
		as.String("bar", seq.Rest().Rest().Rest().Rest().First())
		as.True(seq.Rest().Rest().Rest().Rest().Rest().IsEmpty())
		wg.Done()
	}

	wg.Add(4)
	go check()
	go check()
	go gen()
	go check()
	wg.Wait()
}
