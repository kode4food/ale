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
	p1 := async.NewPromise(func(_ ...data.Value) data.Value {
		return S("hello")
	})
	c1 := p1.(data.Caller).Call()
	as.String("hello", c1())
}

func TestPromiseFailure(t *testing.T) {
	as := assert.New(t)
	p1 := async.NewPromise(func(_ ...data.Value) data.Value {
		panic(errors.New("explosion"))
	})
	c1 := p1.(data.Caller).Call()
	defer as.ExpectPanic("explosion")
	c1()
}
