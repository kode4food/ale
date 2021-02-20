package bootstrap

import (
	"fmt"
	"os"

	"github.com/kode4food/ale/core"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/eval"
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

	ns := b.environment.GetRoot()
	for _, filename = range core.Names() {
		src, _ := core.Get(filename)
		eval.String(ns, data.String(src))
	}
}
