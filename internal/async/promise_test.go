package async_test

import (
	"errors"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/async"
)

func TestPromiseCaller(t *testing.T) {
	as := assert.New(t)
	p1 := async.NewPromise(data.Applicative(func(_ ...data.Value) data.Value {
		return S("hello")
	}, 0))
	c1 := p1.(data.Function)
	as.String("hello", c1.Call())
}

func TestPromiseFailure(t *testing.T) {
	as := assert.New(t)
	p1 := async.NewPromise(data.Applicative(func(_ ...data.Value) data.Value {
		panic(errors.New("explosion"))
	}, 0))
	c1 := p1.(data.Function)
	defer as.ExpectPanic("explosion")
	c1.Call()
}
