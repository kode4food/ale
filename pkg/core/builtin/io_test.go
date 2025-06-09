package builtin_test

import (
	"bytes"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	lang "github.com/kode4food/ale/internal/lang/env"
	"github.com/kode4food/ale/internal/stream"
	"github.com/kode4food/ale/pkg/core/bootstrap"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/eval"
)

func testOutput(t *testing.T, src, expected string) {
	as := assert.New(t)

	buf := bytes.NewBufferString("")
	w, _ := stream.NewWriter(buf, stream.StrOutput).Get(stream.WriteKey)

	e := env.NewEnvironment()
	ns := e.GetRoot()
	as.NoError(env.BindPublic(ns, lang.Out, O(
		C(stream.WriteKey, w),
	)))
	bootstrap.Into(e)

	anon := e.GetAnonymous()
	as.Nil(eval.String(anon, S(src)))

	as.String(expected, buf.String())
}

func TestIOEval(t *testing.T) {
	testOutput(t, `
		(println "hello" "there")
	`, "hello there\n")

	testOutput(t, `
		(print "hello" "there")
	`, "hello there")

	testOutput(t, `
		(print "hello" 99)
	`, "hello 99")

	testOutput(t, `
		(prn "hello" "there")
	`, "\"hello\" \"there\"\n")

	testOutput(t, `
		(pr "hello" "there")
	`, "\"hello\" \"there\"")

	testOutput(t, `
		(pr "hello" 99)
	`, "\"hello\" 99")
}
