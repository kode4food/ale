package eval

import (
	"fmt"

	lang "github.com/kode4food/ale/internal/lang/env"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/internal/stream"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/read"
)

const (
	ErrExpectedPath       = "expected string path, got: %s"
	ErrExpectedFileSystem = "expected file system, got: %s"
)

type Include data.Sequence

func processInclude(ns env.Namespace, v data.Value) (Include, error) {
	l, ok := v.(*data.List)
	if !ok {
		return nil, nil
	}
	f, r, ok := l.Split()
	if !ok || !lang.Include.Equal(f) {
		return nil, nil
	}
	args := sequence.ToVector(r)
	return readInclude(ns, args...)
}

func readInclude(ns env.Namespace, args ...data.Value) (Include, error) {
	if err := data.CheckFixedArity(1, len(args)); err != nil {
		return nil, err
	}
	path, ok := args[0].(data.String)
	if !ok {
		return nil, fmt.Errorf(ErrExpectedPath, args[0])
	}
	c, err := fetchOpenCall(ns)
	if err != nil {
		return nil, err
	}
	res := c.Call(path, stream.ReadAll)
	return read.FromString(res.(data.String)), nil
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
