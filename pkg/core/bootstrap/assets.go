package bootstrap

import (
	"fmt"
	"os"

	lang "github.com/kode4food/ale/internal/lang/env"
	"github.com/kode4food/ale/internal/stream"
	"github.com/kode4food/ale/pkg/core/source"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/eval"
	"github.com/kode4food/ale/pkg/read"
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
	if err := populateFileSystem(ns); err != nil {
		panic(err)
	}

	src, err := source.Assets.ReadFile(InitFile)
	if err != nil {
		panic(err)
	}

	seq := read.FromString(data.String(src))
	if _, err := eval.Block(ns, seq); err != nil {
		panic(err)
	}
}

func populateFileSystem(ns env.Namespace) error {
	e, err := ns.Private(lang.FS)
	if err != nil {
		return fmt.Errorf("failed to get filesystem: %w", err)
	}
	if err = e.Bind(stream.WrapFileSystem(source.Assets)); err != nil {
		return fmt.Errorf("failed to bind filesystem: %w", err)
	}
	return nil
}
