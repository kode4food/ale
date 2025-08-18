package bootstrap

import (
	"fmt"
	"os"

	"github.com/kode4food/ale/core/source"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/read"
)

const InitFile = "bootstrap.ale"

func (b *bootstrap) populateAssets() {
	defer func() {
		if rec := recover(); rec != nil {
			_, _ = fmt.Fprint(os.Stderr, "\nBootstrap Error\n\n")
			if ev, ok := rec.(error); ok {
				msg := ev.Error()
				_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", msg)
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", rec)
			}
			os.Exit(-1)
		}
	}()

	ns := b.environment.GetRoot()
	if err := BindFileSystem(ns, source.Assets); err != nil {
		panic(err)
	}

	src, err := source.Assets.ReadFile(InitFile)
	if err != nil {
		panic(err)
	}

	seq := read.MustFromString(ns, data.String(src))
	if _, err := eval.Block(ns, seq); err != nil {
		panic(err)
	}
}
