package parse

import (
	"fmt"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	lang "github.com/kode4food/ale/internal/lang/env"
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/internal/stream"
)

const (
	ErrExpectedPathString = "expected include path as string, got: %s"
	ErrExpectedFileSystem = "expected file system, got: %s"
)

func (p *parser) processInclude(v ale.Value) (ale.Value, error) {
	path, ok, err := parseInclude(v)
	if err != nil {
		return data.Null, err
	}
	if !ok {
		return v, nil
	}
	inc, err := p.readInclude(path)
	if err != nil {
		return data.Null, err
	}
	p.seq = sequence.Concat(lex.StripWhitespace(inc), p.seq)
	t, err := p.nextToken()
	if err != nil {
		return data.Null, p.maybeWrap(err)
	}
	return p.value(t)
}

func parseInclude(v ale.Value) (string, bool, error) {
	l, ok := v.(*data.List)
	if !ok || l.IsEmpty() {
		return "", false, nil
	}
	f, r, _ := l.Split()
	if _, ok := f.(*lex.Include); !ok {
		return "", false, nil
	}

	args := sequence.ToVector(r)
	if err := data.CheckFixedArity(1, len(args)); err != nil {
		return "", false, err
	}
	if path, ok := args[0].(data.String); ok {
		return string(path), true, nil
	}
	return "", false, fmt.Errorf(ErrExpectedPathString, args[0])
}

func (p *parser) readInclude(path string) (data.Sequence, error) {
	c, err := fetchOpenCall(p.ns)
	if err != nil {
		return data.Null, err
	}
	str := c.Call(data.String(path), stream.ReadString).(data.String)
	return p.tokenize(str)
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
