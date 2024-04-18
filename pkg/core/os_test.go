package core_test

import (
	"testing"
	"time"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/pkg/core"
	"github.com/kode4food/ale/pkg/data"
)

func TestEnv(t *testing.T) {
	as := assert.New(t)

	env := core.Env().(*data.Object)
	as.NotNil(env)
	as.False(env.IsEmpty())
	p := as.MustGet(env, K("PATH")).(data.String)
	as.True(len(p) > 0)
}

func TestArgs(t *testing.T) {
	as := assert.New(t)

	args := core.Args().(data.Vector)
	as.NotNil(args)
	as.False(args.IsEmpty())
	as.Contains("test", args[0])
}

func TestCurrentTime(t *testing.T) {
	as := assert.New(t)

	t1 := time.Now().UnixNano()
	t2 := int64(core.CurrentTime.Call().(data.Integer))
	as.Equal(t1-(t1%1000000), t2-(t2%1000000))
}
