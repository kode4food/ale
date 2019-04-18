package builtin_test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/bootstrap/builtin"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestGo(t *testing.T) {
	as := assert.New(t)
	done := make(chan bool, 0)

	var called bool
	fn := api.Call(func(args ...api.Value) api.Value {
		res := builtin.Str(args...)
		as.String("helloworld", res)
		called = true
		done <- true
		return api.Nil
	})
	builtin.Go(fn, S("hello"), S("world"))
	<-done
	as.True(called)
}
