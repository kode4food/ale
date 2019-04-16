package bootstrap

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/eval"
	"gitlab.com/kode4food/ale/internal/assets"
)

const prefix = "internal/bootstrap/"

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

	ns := b.manager.GetRoot()
	for _, filename = range assets.AssetNames() {
		if !strings.HasPrefix(filename, prefix) {
			continue
		}
		src := api.String(assets.MustGet(filename))
		eval.String(ns, src)
	}
}
