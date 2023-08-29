package special

import (
	"fmt"

	"github.com/kode4food/ale/data"
)

type (
	paramCase struct {
		params data.Names
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
	f := s.First()
	switch f.(type) {
	case data.List, data.Cons, data.LocalSymbol:
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

func (c *paramCase) fixedArgs() data.Names {
	if c.rest {
		return c.params[0 : len(c.params)-1]
	}
	return c.params
}

func (c *paramCase) restArg() (data.Name, bool) {
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
			return append(args[0:cl-2], args[cl-1:]...), true
		}
	}
	return func(args data.Values) (data.Values, bool) {
		return args, cl == len(args)
	}
}

func parseParamNames(v data.Value) (data.Names, bool) {
	switch v := v.(type) {
	case data.LocalSymbol:
		return data.Names{v.Name()}, true
	case data.List:
		return parseListParamNames(v), false
	case data.Cons:
		return parseConsParamNames(v), true
	default:
		panic(fmt.Errorf(ErrUnexpectedParamSyntax, v))
	}
}

func parseListParamNames(l data.List) data.Names {
	var an data.Names
	for f, r, ok := l.Split(); ok; f, r, ok = r.Split() {
		n := f.(data.LocalSymbol).Name()
		an = append(an, n)
	}
	return an
}

func parseConsParamNames(c data.Cons) data.Names {
	var an data.Names
	next := c
	for {
		an = append(an, next.Car().(data.LocalSymbol).Name())

		cdr := next.Cdr()
		if nc, ok := cdr.(data.Cons); ok {
			next = nc
			continue
		}

		an = append(an, cdr.(data.LocalSymbol).Name())
		return an
	}
}
