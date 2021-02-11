package builtin_test

import (
	"testing"

	"github.com/kode4food/ale/core/internal/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestGo(t *testing.T) {
	as := assert.New(t)
	done := make(chan bool, 0)

	var called bool
	fn := data.Applicative(func(args ...data.Value) data.Value {
		res := builtin.Str.Call(args...)
		as.String("helloworld", res)
		called = true
		done <- true
		return data.Nil
	})
	builtin.Go.Call(fn, S("hello"), S("world"))
	<-done
	as.True(called)
}

func TestChan(t *testing.T) {
	as := assert.New(t)

	ch := builtin.Chan.Call(I(0)).(data.Mapped)
	emit, ok1 := ch.Get(builtin.EmitKey)
	closeChan, ok2 := ch.Get(builtin.CloseKey)
	seq, ok3 := ch.Get(builtin.SequenceKey)
	as.True(ok1)
	as.True(ok2)
	as.True(ok3)

	go func() {
		emit.(data.Function).Call(S("hello"))
		closeChan.(data.Function).Call()
	}()

	f, r, ok := seq.(data.Sequence).Split()
	as.String("hello", f)
	as.True(r.IsEmpty())
	as.True(ok)
}

func makeWrapperFunc(v data.Value) data.Function {
	return data.Applicative(func(_ ...data.Value) data.Value {
		return v
	})
}
func TestPromise(t *testing.T) {
	as := assert.New(t)

	p1 := builtin.Promise.Call(makeWrapperFunc(S("with initial")))
	as.True(builtin.IsPromise.Call(p1))
	as.False(builtin.IsResolved.Call(p1))
	res := p1.(data.Function).Call()
	as.True(builtin.IsResolved.Call(p1))
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
