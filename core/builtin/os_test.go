package builtin_test

import (
	"testing"
	"time"

	"github.com/kode4food/ale/core/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestEnv(t *testing.T) {
	as := assert.New(t)

	env := builtin.Env().(*data.Object)
	as.NotNil(env)
	as.False(env.IsEmpty())
	p := as.MustGet(env, K("PATH")).(data.String)
	as.True(len(p) > 0)
}

func TestArgs(t *testing.T) {
	as := assert.New(t)

	args := builtin.Args().(data.Vector)
	as.NotNil(args)
	as.False(args.IsEmpty())
	as.Contains("test", args[0])
}

func TestCurrentTime(t *testing.T) {
	as := assert.New(t)

	t1 := time.Now().UnixNano()
	t2 := int64(builtin.CurrentTime.Call().(data.Integer))
	as.Equal(t1-(t1%1000000), t2-(t2%1000000))
}
