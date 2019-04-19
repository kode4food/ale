package stdlib_test

import (
	"sync"
	"testing"
	"time"

	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/stdlib"
)

func TestChannel(t *testing.T) {
	as := assert.New(t)

	e, seq := stdlib.NewChannel()
	seq = seq.Prepend(F(1))
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
		as.Float(1, seq.First())
		as.Float(2, seq.Rest().First())
		as.Float(3, seq.Rest().Rest().First())
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

func TestPromise(t *testing.T) {
	as := assert.New(t)
	p1 := stdlib.NewPromise()

	go func() {
		time.Sleep(time.Millisecond * 50)
		p1.Deliver(S("hello"))
	}()

	as.Contains(":type promise", p1)
	as.String("hello", p1.Resolve())
	p1.Deliver(S("hello"))
	as.String("hello", p1.Resolve())

	defer as.ExpectPanic(stdlib.ExpectedUndelivered)
	p1.Deliver(S("goodbye"))
}
