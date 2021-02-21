package bootstrap

import (
	"fmt"
	"os"

	"github.com/kode4food/ale/core"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/internal/do"
	"github.com/kode4food/ale/read"
)

var (
	readAssetsOnce = do.Once()

	assets []data.Sequence
)

func (b *bootstrap) assets() {
	var filename string

	defer func() {
		if rec := recover(); rec != nil {
			fmt.Fprint(os.Stderr, "\nBootstrap Error\n\n")
			if ev, ok := rec.(error); ok {
				msg := ev.Error()
				fmt.Fprintf(os.Stderr, "  %s: %s\n\n", filename, msg)
			} else {
				fmt.Fprintf(os.Stderr, "  %s: %s\n\n", filename, rec)
			}
			os.Exit(-1)
		}
	}()

	readAssetsOnce(func() {
		names := core.Names()
		assets = make([]data.Sequence, len(names))
		for i, filename := range names {
			src, _ := core.Get(filename)
			assets[i] = read.FromString(data.String(src))
		}
	})

	ns := b.environment.GetRoot()
	for _, s := range assets {
		eval.Block(ns, s)
	}
}
