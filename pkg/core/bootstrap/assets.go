package bootstrap

import (
	"fmt"
	"os"

	"github.com/kode4food/ale/internal/do"
	"github.com/kode4food/ale/pkg/core/source"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/eval"
	"github.com/kode4food/ale/pkg/read"
)

type asset struct {
	block data.Sequence
	name  string
}

var (
	readAssetsOnce = do.Once()

	assets []asset
)

func (b *bootstrap) populateAssets() {
	var filename string

	defer func() {
		if rec := recover(); rec != nil {
			_, _ = fmt.Fprint(os.Stderr, "\nBootstrap Error\n\n")
			if ev, ok := rec.(error); ok {
				msg := ev.Error()
				_, _ = fmt.Fprintf(os.Stderr, "  %s: %s\n\n", filename, msg)
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "  %s: %s\n\n", filename, rec)
			}
			os.Exit(-1)
		}
	}()

	readAssetsOnce(func() {
		for _, filename = range source.Names() {
			src, _ := source.Get(filename)
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
