package builtin_test

import (
	"testing"

	"gitlab.com/kode4food/ale/bootstrap/builtin"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestGo(t *testing.T) {
	as := assert.New(t)
	done := make(chan bool, 0)

	var called bool
	fn := data.Call(func(args ...data.Value) data.Value {
		res := builtin.Str(args...)
		as.String("helloworld", res)
		called = true
		done <- true
		return data.Nil
	})
	builtin.Go(fn, S("hello"), S("world"))
	<-done
	as.True(called)
}

func TestChan(t *testing.T) {
	as := assert.New(t)

	ch := builtin.Chan().(data.Mapped)
	emit, ok1 := ch.Get(builtin.EmitKey)
	closeChan, ok2 := ch.Get(builtin.CloseKey)
	seq, ok3 := ch.Get(builtin.SequenceKey)
	as.True(ok1)
	as.True(ok2)
	as.True(ok3)

	go func() {
		emit.(*data.Function).Call(S("hello"))
		closeChan.(*data.Function).Call()
	}()

	f, r, ok := seq.(data.Sequence).Split()
	as.String("hello", f)
	as.True(r.IsEmpty())
	as.True(ok)
}

func TestPromise(t *testing.T) {
	as := assert.New(t)

	p1 := builtin.Promise(S("with initial"))
	as.True(builtin.IsPromise(p1))
	res := getCall(p1)()
	as.String("with initial", res)

	p2 := builtin.Promise()
	go func() {
		getCall(p2)(S("no initial"))
	}()
	res = getCall(p2)()
	as.String("no initial", res)
	as.False(builtin.IsPromise(res))

	defer as.ExpectPanic("can't deliver a promise twice")
	getCall(p1)(S("new value"))
}
