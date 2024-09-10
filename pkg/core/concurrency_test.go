package core_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/stream"
	"github.com/kode4food/ale/pkg/core"
	"github.com/kode4food/ale/pkg/data"
)

func TestGo(t *testing.T) {
	as := assert.New(t)
	done := make(chan bool)

	var called bool
	fn := data.MakeProcedure(func(args ...data.Value) data.Value {
		res := core.Str.Call(args...)
		as.String("helloworld", res)
		called = true
		done <- true
		return data.Null
	})
	core.Go.Call(fn, S("hello"), S("world"))
	<-done
	as.True(called)
}

func TestChan(t *testing.T) {
	as := assert.New(t)

	ch := core.Chan.Call(I(0)).(data.Mapped)
	emit, ok1 := ch.Get(stream.EmitKey)
	closeChan, ok2 := ch.Get(stream.CloseKey)
	seq, ok3 := ch.Get(stream.SequenceKey)
	as.True(ok1)
	as.True(ok2)
	as.True(ok3)

	go func() {
		emit.(data.Procedure).Call(S("hello"))
		closeChan.(data.Procedure).Call()
	}()

	f, r, ok := seq.(data.Sequence).Split()
	as.String("hello", f)
	as.True(r.IsEmpty())
	as.True(ok)
}

func makeWrapperFunc(v data.Value) data.Procedure {
	return data.MakeProcedure(func(...data.Value) data.Value {
		return v
	})
}

func TestPromise(t *testing.T) {
	as := assert.New(t)

	p1 := core.Promise.Call(makeWrapperFunc(S("with initial")))
	as.True(getPredicate(core.PromiseKey).Call(p1))
	as.False(getPredicate(core.ResolvedKey).Call(p1))
	res := p1.(data.Procedure).Call()
	as.True(getPredicate(core.ResolvedKey).Call(p1))
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

	as.EvalTo(`
		(define p3 (delay "nope"))
		[(promise-forced? p3)
         (promise-forced? "string")]
	`, V(data.False, data.True))
}

func TestFutureEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(define p (future "hello"))
		(p)
	`, S("hello"))
}
