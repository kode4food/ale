package builtin

import (
	"fmt"
	"slices"

	"github.com/kode4food/ale/data"
)

type (
	paramCase struct {
		signature data.Value
		params    data.Locals
		rest      bool
		body      data.Sequence
	}

	paramCases []*paramCase

	argFetcher func(data.Vector) (data.Vector, bool)
)

const (
	// ErrUnexpectedCaseSyntax is raised when a call to Lambda doesn't include
	// a proper parameter case initializer. If it first encounters a Vector,
	// the parsing will assume multiple parameter cases, otherwise it will
	// assume a single parameter case
	ErrUnexpectedCaseSyntax = "unexpected case syntax: %s"

	// ErrUnexpectedParamSyntax is raised when a Lambda parameter case is
	// represented by an unexpected syntax. Valid syntax representations are
	// data.List, data.Cons, or data.Local
	ErrUnexpectedParamSyntax = "unexpected parameter syntax: %s"

	// ErrUnreachableCase is raised when a Lambda parameter is defined that
	// would otherwise be impossible to reach given previous definitions
	ErrUnreachableCase = "unreachable parameter case: %s, matched by: %s"
)

func parseParamCases(s data.Sequence) paramCases {
	if s.IsEmpty() {
		return paramCases{}
	}
	f := s.Car()
	switch f.(type) {
	case *data.List, *data.Cons, data.Local:
		c := parseParamCase(s)
		return paramCases{c}
	case data.Vector:
		var res paramCases
		var err error
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			c := parseParamCase(f.(data.Vector))
			res, err = res.addParamCase(c)
			if err != nil {
				panic(err)
			}
		}
		return res
	default:
		panic(fmt.Errorf(ErrUnexpectedCaseSyntax, f))
	}
}

func (pc paramCases) addParamCase(added *paramCase) (paramCases, error) {
	addedLow, addedHigh := added.arityRange()
	for _, orig := range pc {
		origLow, origHigh := orig.arityRange()
		if isUnreachable(origLow, origHigh, addedLow, addedHigh) {
			return pc, fmt.Errorf(
				ErrUnreachableCase, added.signature, orig.signature,
			)
		}
	}
	return append(pc, added), nil
}

func isUnreachable(origLow, origHigh, addedLow, addedHigh int) bool {
	if origHigh == data.OrMore {
		return addedLow >= origLow
	}
	return addedHigh != data.OrMore && addedLow == origLow
}

func (pc paramCases) makeArityChecker() data.ArityChecker {
	switch len(pc) {
	case 0:
		return data.MakeChecker(0)
	case 1:
		l, h := pc[0].arityRange()
		return data.MakeChecker(l, h)
	default:
		lower, upper := pc[0].arityRange()
		for _, s := range pc[1:] {
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
}

func (pc paramCases) makeFetchers() []argFetcher {
	res := make([]argFetcher, len(pc))
	for i, c := range pc {
		res[i] = c.makeFetcher()
	}
	return res
}

func parseParamCase(s data.Sequence) *paramCase {
	f, body, _ := s.Split()
	argNames, restArg := parseParamNames(f)
	return &paramCase{
		signature: f,
		params:    argNames,
		rest:      restArg,
		body:      body,
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
		return func(args data.Vector) (data.Vector, bool) {
			if len(args) < cl-1 {
				return args, false
			}
			res := append(slices.Clone(args[0:cl-1]), args[cl-1:])
			return res, true
		}
	}
	return func(args data.Vector) (data.Vector, bool) {
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
