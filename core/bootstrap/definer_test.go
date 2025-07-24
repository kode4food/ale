package bootstrap_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/internal/assert"
)

func TestDefiner(t *testing.T) {
	as := assert.New(t)

	e := bootstrap.DevNullEnvironment()
	ns := e.GetRoot()
	res, err := eval.String(ns, "read")
	if as.NoError(err) {
		_, ok := res.(data.Procedure)
		as.True(ok)
	}

	_, err = eval.String(ns, "(def-builtin read)")
	as.ExpectError(fmt.Errorf(env.ErrNameAlreadyBound, "read"), err)

	_, err = eval.String(ns, "(def-builtin nope)")
	as.ExpectError(fmt.Errorf(bootstrap.ErrBuiltInNotFound, "nope"), err)

	_, err = eval.String(ns, "(def-builtin too many)")
	as.ExpectError(fmt.Errorf(data.ErrFixedArity, 1, 2), err)
}
