package builtin_test

import (
	"bytes"
	"testing"

	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/core/internal/builtin"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/stream"
)

const stdoutName = "*out*"

func bindWrite(w stream.Writer) data.Function {
	return data.Applicative(func(args ...data.Value) data.Value {
		for _, af := range args {
			w.Write(af)
		}
		return nil
	})
}

func testOutput(t *testing.T, src, expected string) {
	as := assert.New(t)

	buf := bytes.NewBufferString("")
	w := stream.NewWriter(buf, stream.StrOutput)

	e := env.NewEnvironment()
	ns := e.GetRoot()
	ns.Declare(stdoutName).Bind(O(
		C(builtin.WriterKey, w),
		C(builtin.WriteKey, bindWrite(w)),
	))
	bootstrap.Into(e)

	anon := e.GetAnonymous()
	eval.String(anon, S(src))

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
