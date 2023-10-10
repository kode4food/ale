package special

import (
	"fmt"

	"github.com/kode4food/ale/data"
)

type (
	paramCase struct {
		params data.Locals
		rest   bool
		body   data.Sequence
	}

	paramCases []*paramCase

	argFetcher func(data.Values) (data.Values, bool)
)

// Error messages
const (
	ErrUnexpectedParamSyntax = "unexpected parameter syntax: %s"
)

func parseParamCases(s data.Sequence) paramCases {
	f := s.Car()
	switch f.(type) {
	case *data.List, *data.Cons, data.Local:
		c := parseParamCase(s)
		return paramCases{c}
	case data.Vector:
		var res paramCases
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			c := parseParamCase(f.(data.Vector))
			res = append(res, c)
		}
		return res
	default:
		panic(fmt.Errorf(ErrUnexpectedCaseSyntax, f))
	}
}

func (lc paramCases) makeArityChecker() data.ArityChecker {
	v0 := lc[0]
	lower, upper := v0.arityRange()
	for _, s := range lc[1:] {
		l, u := s.arityRange()
		lower = min(l, lower)
		if u == data.OrMore || upper == data.OrMore {
			upper = data.OrMore
			continue
		}
		upper = max(u, upper)
	}
	return data.MakeChecker(lower, upper)
}

func (lc paramCases) makeFetchers() []argFetcher {
	res := make([]argFetcher, len(lc))
	for i, c := range lc {
		res[i] = c.makeFetcher()
	}
	return res
}

func parseParamCase(s data.Sequence) *paramCase {
	f, body, _ := s.Split()
	argNames, restArg := parseParamNames(f)
	return &paramCase{
		params: argNames,
		rest:   restArg,
		body:   body,
	}
}

func (c *paramCase) fixedArgs() data.Locals {
	if c.rest {
		return c.params[0 : len(c.params)-1]
	}
	return c.params
}

func (c *paramCase) restArg() (data.Local, bool) {
	if c.rest {
		return c.params[len(c.params)-1], true
	}
	return "", false
}

func (c *paramCase) arityRange() (int, int) {
	fl := len(c.fixedArgs())
	if _, ok := c.restArg(); ok {
		return fl, data.OrMore
	}
	return fl, fl
}

func (c *paramCase) makeFetcher() argFetcher {
	cl := len(c.params)
	if c.rest {
		return func(args data.Values) (data.Values, bool) {
			if len(args) < cl {
				return args, false
			}
			res := make(data.Values, cl)
			copy(res, args[0:cl-1])
			res[cl-1] = data.NewVector(args[cl-1:]...)
			return res, true
		}
	}
	return func(args data.Values) (data.Values, bool) {
		return args, cl == len(args)
	}
}

func parseParamNames(v data.Value) (data.Locals, bool) {
	switch v := v.(type) {
	case data.Local:
		return data.Locals{v}, true
	case *data.List:
		return parseListParamNames(v), false
	case *data.Cons:
		return parseConsParamNames(v), true
	default:
		panic(fmt.Errorf(ErrUnexpectedParamSyntax, v))
	}
}

func parseListParamNames(l *data.List) data.Locals {
	var an data.Locals
	for f, r, ok := l.Split(); ok; f, r, ok = r.Split() {
		n := f.(data.Local)
		an = append(an, n)
	}
	return an
}

func parseConsParamNames(c *data.Cons) data.Locals {
	var an data.Locals
	next := c
	for {
		an = append(an, next.Car().(data.Local))

		cdr := next.Cdr()
		if nc, ok := cdr.(*data.Cons); ok {
			next = nc
			continue
		}

		an = append(an, cdr.(data.Local))
		return an
	}
}
