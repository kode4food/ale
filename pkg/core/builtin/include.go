package builtin

import (
	"fmt"

	lang "github.com/kode4food/ale/internal/lang/env"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/internal/stream"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/eval"
	"github.com/kode4food/ale/pkg/read"
)

const (
	ErrExpectedPath       = "expected string path, got: %s"
	ErrExpectedFileSystem = "expected file system, got: %s"
)

func Include(ns env.Namespace, args ...data.Value) data.Value {
	data.MustCheckMinimumArity(1, len(args))
	c, err := fetchOpenCall(ns)
	if err != nil {
		panic(err)
	}
	var res []data.Sequence
	for _, arg := range args {
		path, ok := arg.(data.String)
		if !ok {
			panic(fmt.Errorf(ErrExpectedPath, arg))
		}
		str := c.Call(path, stream.ReadAll).(data.String)
		seq := read.FromString(str)
		res = append(res, seq)
	}
	return eval.Include(sequence.Concat(res...))
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
