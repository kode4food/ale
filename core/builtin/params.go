package builtin

import (
	"fmt"
	"slices"
	"strings"

	"github.com/kode4food/ale/data"
)

type (
	paramCase struct {
		signature data.Value
		params    data.Locals
		rest      bool
		body      data.Sequence
	}

	paramCases struct {
		cases   []*paramCase
		fixed   []uint8
		hasRest bool
		lowRest int
	}

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
	ErrUnreachableCase = "unreachable parameter case: %s"

	// ErrUnmatchedCase is raised when a Lambda is called and the number of
	// arguments provided doesn't match any of the declared parameter cases
	ErrUnmatchedCase = "got %d arguments, expected %s"
)

const arityBits = 8

func parseParamCases(s data.Sequence) *paramCases {
	res := &paramCases{}
	if s.IsEmpty() {
		return res
	}
	f := s.Car()
	switch f.(type) {
	case *data.List, *data.Cons, data.Local:
		c := parseParamCase(s)
		if err := res.addParamCase(c); err != nil {
			panic(err)
		}
		return res
	case data.Vector:
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			c := parseParamCase(f.(data.Vector))
			if err := res.addParamCase(c); err != nil {
				panic(err)
			}
		}
		return res
	default:
		panic(fmt.Errorf(ErrUnexpectedCaseSyntax, f))
	}
}

func (pc *paramCases) Cases() []*paramCase {
	return pc.cases
}

func (pc *paramCases) addParamCase(added *paramCase) error {
	a, ar := added.getArity()
	if !pc.isReachable(a, ar) {
		return fmt.Errorf(ErrUnreachableCase, added.signature)
	}
	pc.cases = append(pc.cases, added)
	if ar {
		pc.addRest(a)
	} else {
		pc.addFixed(a)
	}
	return nil
}

func (pc *paramCases) isReachable(i int, isRest bool) bool {
	if len(pc.cases) == 0 {
		return true
	}
	if pc.hasRest {
		return i < pc.lowRest
	}
	if isRest {
		return true
	}
	index, offset := i/arityBits, i%arityBits
	if index < len(pc.fixed) {
		return (pc.fixed[index] & (1 << offset)) == 0
	}
	return true
}

func (pc *paramCases) makeFetchers() []argFetcher {
	res := make([]argFetcher, len(pc.cases))
	for i, c := range pc.cases {
		res[i] = c.makeFetcher()
	}
	return res
}

func (pc *paramCases) makeChecker() data.ArityChecker {
	if pc.hasRest {
		return pc.makeRestChecker()
	}
	return pc.makeFixedChecker()
}

func (pc *paramCases) makeFixedChecker() data.ArityChecker {
	fixed := pc.fixed
	signatures := pc.signatures()
	return func(i int) error {
		index, offset := i/arityBits, i%arityBits
		if index >= len(fixed) || fixed[index]&(1<<offset) == 0 {
			return fmt.Errorf(ErrUnmatchedCase, i, signatures)
		}
		return nil
	}
}

func (pc *paramCases) makeRestChecker() data.ArityChecker {
	lowRest := pc.lowRest
	fixedChecker := pc.makeFixedChecker()
	return func(i int) error {
		if i >= lowRest {
			return nil
		}
		return fixedChecker(i)
	}
}

func (pc *paramCases) addFixed(i int) {
	index, offset := i/arityBits, i%arityBits
	for len(pc.fixed) <= index {
		pc.fixed = append(pc.fixed, 0)
	}
	pc.fixed[index] |= 1 << offset
}

func (pc *paramCases) addRest(i int) {
	if pc.hasRest {
		pc.lowRest = min(pc.lowRest, i)
	} else {
		pc.lowRest = i
		pc.hasRest = true
	}
}

func (pc *paramCases) signatures() string {
	var res []string
	for _, r := range pc.fixedRanges() {
		if pc.hasRest && r[1] >= pc.lowRest-1 {
			res = append(res, formatOrMore(r[0]))
			return strings.Join(res, ", ")
		}
		res = append(res, formatRange(r))
	}

	if pc.hasRest {
		res = append(res, formatOrMore(pc.lowRest))
	}

	return strings.Join(res, ", ")
}

func (pc *paramCases) fixedRanges() [][2]int {
	fixed := pc.fixedSet()
	if len(fixed) == 0 {
		return [][2]int{}
	}

	var res [][2]int
	start := fixed[0]
	for i := 1; i < len(fixed); i++ {
		if fixed[i] != fixed[i-1]+1 {
			res = append(res, [2]int{start, fixed[i-1]})
			start = fixed[i]
		}
	}
	res = append(res, [2]int{start, fixed[len(fixed)-1]})
	return res
}

func (pc *paramCases) fixedSet() []int {
	var res []int
	for i := 0; i < len(pc.fixed)*arityBits; i++ {
		index, offset := i/arityBits, i%arityBits
		if pc.fixed[index]&(1<<offset) != 0 {
			res = append(res, i)
		}
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

func (c *paramCase) getArity() (int, bool) {
	fl := len(c.fixedArgs())
	if _, ok := c.restArg(); ok {
		return fl, true
	}
	return fl, false
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

func formatRange(r [2]int) string {
	if r[0] == r[1] {
		return fmt.Sprintf("%d", r[0])
	} else {
		return fmt.Sprintf("%d-%d", r[0], r[1])
	}
}

func formatOrMore(i int) string {
	return fmt.Sprintf("%d or more", i)
}
