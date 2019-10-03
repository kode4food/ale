package stdlib_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/stdlib"
)

func TestChannel(t *testing.T) {
	as := assert.New(t)

	e, seq := stdlib.NewChannel(0)
	seq = seq.(data.Prepender).Prepend(F(1))
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

func TestPromiseCaller(t *testing.T) {
	as := assert.New(t)
	p1 := stdlib.NewPromise(func(_ ...data.Value) data.Value {
		return S("hello")
	})
	c1 := p1.(data.Caller).Call()
	as.String("hello", c1())
}

func TestPromiseFailure(t *testing.T) {
	as := assert.New(t)
	p1 := stdlib.NewPromise(func(_ ...data.Value) data.Value {
		panic(fmt.Errorf("'splosion!"))
	})
	c1 := p1.(data.Caller).Call()
	defer as.ExpectPanic("'splosion!")
	c1()
}
