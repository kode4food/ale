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

type asset struct {
	name  string
	block data.Sequence
}

var (
	readAssetsOnce = do.Once()

	assets []asset
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
		for _, filename = range core.Names() {
			src, _ := core.Get(filename)
			assets = append(assets, asset{
				name:  filename,
				block: read.FromString(data.String(src)),
			})
		}
	})

	ns := b.environment.GetRoot()
	for _, a := range assets {
		filename = a.name
		eval.Block(ns, a.block)
	}
}
