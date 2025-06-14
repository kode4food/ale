package parse

import (
	"fmt"

	lang "github.com/kode4food/ale/internal/lang/env"
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/internal/stream"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

const (
	ErrExpectedPath       = "expected include path, got: %s"
	ErrExpectedFileSystem = "expected file system, got: %s"
)

func (r *parser) processInclude(v data.Value, err error) (data.Value, error) {
	if err != nil {
		return v, err
	}
	path, ok, err := parseInclude(v)
	if err != nil {
		return data.Null, err
	}
	if !ok {
		return v, nil
	}
	inc, err := r.readInclude(path)
	if err != nil {
		return data.Null, err
	}
	r.seq = sequence.Concat(lex.StripWhitespace(inc), r.seq)
	t, err := r.nextToken()
	if err != nil {
		return data.Null, r.maybeWrap(err)
	}
	return r.value(t)
}

func parseInclude(v data.Value) (string, bool, error) {
	l, ok := v.(*data.List)
	if !ok || l.IsEmpty() {
		return "", false, nil
	}
	f, r, _ := l.Split()
	if _, ok := f.(*lex.Include); !ok {
		return "", false, nil
	}

	args := sequence.ToVector(r)
	if len(args) != 1 {
		return "", false, fmt.Errorf(ErrExpectedPath, data.Null)
	}
	if path, ok := args[0].(data.String); ok {
		return string(path), true, nil
	}
	return "", false, fmt.Errorf(ErrExpectedPath, args[0])
}

func (r *parser) readInclude(path string) (data.Sequence, error) {
	c, err := fetchOpenCall(r.ns)
	if err != nil {
		panic(err)
	}
	str := c.Call(data.String(path), stream.ReadAll).(data.String)
	return r.tokenize(str)
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
