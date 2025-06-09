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
	ErrExpectedPath    = "expected string path, got: %s"
	ErrOpenUnsupported = "file system does not support opening files: %s"
	ErrExpectedString  = "expected string result from open, got: %s"
	ErrExpectedObject  = "expected file system to be an object, got: %s"
	ErrExpectedCaller  = "expected open to be proc, got: %s"
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
	str, ok := res.(data.String)
	if !ok {
		return nil, fmt.Errorf(ErrExpectedString, res)
	}
	return read.FromString(str), nil
}

func fetchOpenCall(ns env.Namespace) (data.Caller, error) {
	e, _, err := ns.Resolve(lang.FS)
	if err != nil {
		return nil, err
	}
	v, err := e.Value()
	if err != nil {
		return nil, err
	}
	fs, ok := v.(*data.Object)
	if !ok {
		return nil, fmt.Errorf(ErrExpectedObject, v)
	}
	v, ok = fs.Get(stream.OpenKey)
	if !ok {
		return nil, fmt.Errorf(ErrOpenUnsupported, fs)
	}
	c, ok := v.(data.Caller)
	if !ok {
		return nil, fmt.Errorf(ErrExpectedCaller, v)
	}
	return c, nil
}
