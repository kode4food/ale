package internal_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kode4food/ale/cmd/ale/internal"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/internal/assert"
)

func TestBuiltIns(t *testing.T) {
	as := assert.New(t)
	r := internal.NewREPL()
	ns := r.GetNS()

	for _, name := range []data.Local{"help", "debug", "cls"} {
		v, err := env.ResolveValue(ns, name)
		if as.NoError(err) {
			p, ok := v.(data.Procedure)
			as.True(ok)
			if ok {
				as.NotPanics(func() {
					_ = p.Call()
				})
			}
		}
	}
}

func TestEvaluateFile(t *testing.T) {
	as := assert.New(t)
	dir := t.TempDir()
	fileName := filepath.Join(dir, "evaluate.ale")
	as.NoError(os.WriteFile(fileName, []byte("(+ 1 2)\n"), 0o600))
	as.NotPanics(func() {
		internal.EvaluateFile(fileName)
	})
}

func TestEvaluateStdIn(t *testing.T) {
	as := assert.New(t)
	as.NotPanics(func() {
		withStdInFile(t, "(+ 2 3)\n", func() {
			internal.EvaluateStdIn()
		})
	})
}

func TestDocBuiltIn(t *testing.T) {
	as := assert.New(t)
	ns := internal.NewREPL().GetNS()

	as.NotPanics(func() {
		_, err := eval.String(ns, "(doc)")
		as.NoError(err)
	})

	as.NotPanics(func() {
		_, err := eval.String(ns, "(doc doc)")
		as.NoError(err)
	})

	as.NotPanics(func() {
		_, err := eval.String(ns, "(doc +)")
		as.NoError(err)
	})
}

func withStdInFile(t *testing.T, src string, fn func()) {
	t.Helper()
	fileName := filepath.Join(t.TempDir(), "stdin.ale")
	if err := os.WriteFile(fileName, []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}

	oldIn := os.Stdin
	f, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Stdin = oldIn
		_ = f.Close()
	}()

	os.Stdin = f
	fn()
}
