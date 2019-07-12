package builtin_test

import (
	"testing"

	"gitlab.com/kode4food/ale/core/builtin"
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
		return data.Null
	})
	builtin.Go(fn, S("hello"), S("world"))
	<-done
	as.True(called)
}

func TestChan(t *testing.T) {
	as := assert.New(t)

	ch := builtin.Chan(data.Integer(0)).(data.Mapped)
	emit, ok1 := ch.Get(builtin.EmitKey)
	closeChan, ok2 := ch.Get(builtin.CloseKey)
	seq, ok3 := ch.Get(builtin.SequenceKey)
	as.True(ok1)
	as.True(ok2)
	as.True(ok3)

	go func() {
		emit.(data.Call)(S("hello"))
		closeChan.(data.Call)()
	}()

	f, r, ok := seq.(data.Sequence).Split()
	as.String("hello", f)
	as.True(r.IsEmpty())
	as.True(ok)
}

func makeWrapperFunc(v data.Value) data.Call {
	return data.Call(func(_ ...data.Value) data.Value {
		return v
	})
}
func TestPromise(t *testing.T) {
	as := assert.New(t)

	p1 := builtin.Promise(makeWrapperFunc(S("with initial")))
	as.True(builtin.IsPromise(p1))
	as.False(builtin.IsResolved(p1))
	res := getCall(p1)()
	as.True(builtin.IsResolved(p1))
	as.String("with initial", res)
}

func TestGenerateEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define g (generate
			(emit 99)
			(emit 100 1000)))
		(apply + g)
	`, F(1199))
}

func TestDelayEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define p1 (delay "blah"))
		(promise? p1)
	`, data.True)

	as.EvalTo(`
		(define p2 (delay "hello"))
		(p2)
	`, S("hello"))
}

func TestFutureEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define p (future "hello"))
		(p)
	`, S("hello"))
}
