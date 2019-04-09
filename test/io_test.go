package test

import (
	"bytes"
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/eval"
	"gitlab.com/kode4food/ale/internal/assert"
	"gitlab.com/kode4food/ale/internal/bootstrap"
	"gitlab.com/kode4food/ale/internal/builtin"
	"gitlab.com/kode4food/ale/internal/namespace"
	"gitlab.com/kode4food/ale/stdlib"
)

const stdoutName = "*out*"

func bindWrite(w stdlib.Writer) api.Call {
	return func(args ...api.Value) api.Value {
		for _, af := range args {
			w.Write(af)
		}
		return nil
	}
}

func testOutput(t *testing.T, src, expected string) {
	as := assert.New(t)

	buf := bytes.NewBufferString("")
	w := stdlib.NewWriter(buf, stdlib.StrOutput)

	manager := namespace.NewManager()
	ns := manager.GetRoot()
	ns.Bind(stdoutName, api.Object{
		builtin.WriterKey: w,
		builtin.WriteKey:  bindWrite(w),
	})
	bootstrap.Into(manager)

	anon := manager.GetAnonymous()
	eval.String(anon, api.String(src))

	as.String(expected, buf.String())
}

func TestIO(t *testing.T) {
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
