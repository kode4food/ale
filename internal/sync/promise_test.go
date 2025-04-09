package sync_test

import (
	"errors"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/sync"
	"github.com/kode4food/ale/pkg/data"
)

func TestPromiseCaller(t *testing.T) {
	as := assert.New(t)
	p1 := sync.NewPromise(data.MakeProcedure(func(...data.Value) data.Value {
		return S("hello")
	}, 0))
	as.String("hello", p1.Call())
}

func TestPromiseFailure(t *testing.T) {
	as := assert.New(t)
	p1 := sync.NewPromise(data.MakeProcedure(func(...data.Value) data.Value {
		panic(errors.New("explosion"))
	}, 0))
	defer as.ExpectPanic(errors.New("explosion"))
	p1.Call()
}

func TestPromiseEval(t *testing.T) {
	as := assert.New(t)

	as.MustEvalTo(`
		(define p (delay "hello"))
		(let* ([p?  (promise? p)]
			   [r1? (resolved? p)]
 			   [f   (p)]
  			   [r2? (resolved? p)])
		  [p? r1? f r2?])
	`, V(data.True, data.False, S("hello"), data.True))
}
