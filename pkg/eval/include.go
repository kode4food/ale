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
	ErrExpectedPath    = "expected a string path, got: %s"
	ErrOpenUnsupported = "file system does not support opening files: %s"
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
	data.MustCheckFixedArity(1, len(args))
	path, ok := args[0].(data.String)
	if !ok {
		panic(fmt.Errorf(ErrExpectedPath, args[0]))
	}
	c := mustFetchOpenCall(ns)
	res := c.Call(path, stream.ReadAll).(data.String)
	return read.FromString(res), nil
}

func mustFetchOpenCall(ns env.Namespace) data.Caller {
	fs := env.MustResolveValue(ns, lang.FS).(*data.Object)
	v, ok := fs.Get(stream.OpenKey)
	if !ok {
		panic(fmt.Errorf(ErrOpenUnsupported, fs))
	}
	return v.(data.Caller)
}
