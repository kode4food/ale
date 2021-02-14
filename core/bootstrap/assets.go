package bootstrap

import (
	"fmt"
	"os"
	"strings"

	"github.com/kode4food/ale/core/assets"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/eval"
)

const prefix = "core/"

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

	ns := b.environment.GetRoot()
	for _, filename = range assets.AssetNames() {
		if !strings.HasPrefix(filename, prefix) {
			continue
		}
		src := data.String(assets.MustGet(filename))
		eval.String(ns, src)
	}
}
