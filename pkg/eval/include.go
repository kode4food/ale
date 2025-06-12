package eval

import (
	"fmt"

	lang "github.com/kode4food/ale/internal/lang/env"
	"github.com/kode4food/ale/internal/stream"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/read"
)

const (
	ErrExpectedPath       = "expected string path, got: %s"
	ErrExpectedFileSystem = "expected file system, got: %s"
)

type include struct {
	forms data.Sequence
}

func Include(ns env.Namespace, args ...data.Value) data.Value {
	data.MustCheckFixedArity(1, len(args))
	path, ok := args[0].(data.String)
	if !ok {
		panic(fmt.Errorf(ErrExpectedPath, args[0]))
	}
	c, err := fetchOpenCall(ns)
	if err != nil {
		panic(err)
	}
	res := c.Call(path, stream.ReadAll)
	return &include{read.FromString(res.(data.String))}
}

func (i *include) Equal(other data.Value) bool {
	return i == other
}

func fetchOpenCall(ns env.Namespace) (data.Procedure, error) {
	e, _, err := ns.Resolve(lang.FS)
	if err != nil {
		return nil, err
	}
	v, err := e.Value()
	if err != nil {
		return nil, err
	}
	fs, ok := v.(*stream.FileSystem)
	if !ok {
		return nil, fmt.Errorf(ErrExpectedFileSystem, v)
	}
	return fs.Open, nil
}
