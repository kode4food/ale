package internal

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/kode4food/ale/pkg/data"
)

type (
	ParamCase struct {
		Signature data.Value
		Body      data.Sequence
		Params    data.Locals
		Rest      bool
	}

	ParamCases struct {
		Cases   []*ParamCase
		Fixed   []uint8
		HasRest bool
		LowRest int
	}

	ArgFetcher func(data.Vector) (data.Vector, bool)
)

const (
	// ErrNoCasesDefined is raised when a call to Lambda doesn't include any
	// parameter cases definitions.
	ErrNoCasesDefined = "no parameter cases defined"

	// ErrUnexpectedCaseSyntax is raised when a call to Lambda doesn't include
	// a proper parameter case initializer. If it first encounters a Vector,
	// the parsing will assume multiple parameter cases, otherwise it will
	// assume a single parameter case
	ErrUnexpectedCaseSyntax = "unexpected case syntax: %s"

	// ErrUnexpectedParamSyntax is raised when a Lambda parameter case is
	// represented by an unexpected syntax. Valid syntax representations are
	// data.List, data.Cons, or data.Local
	ErrUnexpectedParamSyntax = "unexpected parameter syntax: %s"

	// ErrNoCaseBodyDefined is raised when a Lambda parameter case defines its
	// parameters, but not an associated body to evaluate
	ErrNoCaseBodyDefined = "no parameter case body: %s"

	// ErrUnreachableCase is raised when a Lambda parameter is defined that
	// would otherwise be impossible to reach given previous definitions
	ErrUnreachableCase = "unreachable parameter case: %s"

	// ErrUnmatchedCase is raised when a Lambda is called and the number of
	// arguments provided doesn't match any of the declared parameter cases
	ErrUnmatchedCase = "got %d arguments, expected %s"

	// ErrNoMatchingParamPattern is raised when none of the parameter patterns
	// for a Lambda were capable of being matched
	ErrNoMatchingParamPattern = "no matching parameter pattern"
)

const arityBits = 8

func ParseParamCases(s data.Sequence) (*ParamCases, error) {
	if s.IsEmpty() {
		return nil, errors.New(ErrNoCasesDefined)
	}
	res := &ParamCases{}
	f := s.Car()
	switch f.(type) {
	case *data.List, *data.Cons, data.Local:
		c, err := parseParamCase(s)
		if err != nil {
			return nil, err
		}
		if err := res.addParamCase(c); err != nil {
			return nil, err
		}
		return res, nil
	case data.Vector:
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			c, err := parseParamCase(f.(data.Vector))
			if err != nil {
				return nil, err
			}
			if err := res.addParamCase(c); err != nil {
				return nil, err
			}
		}
		return res, nil
	default:
		return nil, fmt.Errorf(ErrUnexpectedCaseSyntax, f)
	}
}

func (pc *ParamCases) MakeArgFetchers() []ArgFetcher {
	res := make([]ArgFetcher, len(pc.Cases))
	for i, c := range pc.Cases {
		res[i] = c.makeArgFetcher()
	}
	return res
}

func (pc *ParamCases) MakeArityChecker() data.ArityChecker {
	if pc.HasRest {
		return pc.makeRestChecker()
	}
	return pc.makeFixedChecker()
}

func (pc *ParamCases) addParamCase(added *ParamCase) error {
	a, ar := added.getArity()
	if !pc.isReachable(a, ar) {
		return fmt.Errorf(ErrUnreachableCase, added.Signature)
	}
	pc.Cases = append(pc.Cases, added)
	if ar {
		pc.addRest(a)
	} else {
		pc.addFixed(a)
	}
	return nil
}

func (pc *ParamCases) isReachable(i int, isRest bool) bool {
	if len(pc.Cases) == 0 {
		return true
	}
	if pc.HasRest {
		return i < pc.LowRest
	}
	if isRest {
		return true
	}
	index, offset := i/arityBits, i%arityBits
	if index < len(pc.Fixed) {
		return (pc.Fixed[index] & (1 << offset)) == 0
	}
	return true
}

func (pc *ParamCases) makeFixedChecker() data.ArityChecker {
	fixed := pc.Fixed
	signatures := pc.signatures()
	return func(i int) error {
		index, offset := i/arityBits, i%arityBits
		if index >= len(fixed) || fixed[index]&(1<<offset) == 0 {
			return fmt.Errorf(ErrUnmatchedCase, i, signatures)
		}
		return nil
	}
}

func (pc *ParamCases) makeRestChecker() data.ArityChecker {
	lowRest := pc.LowRest
	fixedChecker := pc.makeFixedChecker()
	return func(i int) error {
		if i >= lowRest {
			return nil
		}
		return fixedChecker(i)
	}
}

func (pc *ParamCases) addFixed(i int) {
	index, offset := i/arityBits, i%arityBits
	for len(pc.Fixed) <= index {
		pc.Fixed = append(pc.Fixed, 0)
	}
	pc.Fixed[index] |= 1 << offset
}

func (pc *ParamCases) addRest(i int) {
	if pc.HasRest {
		pc.LowRest = min(pc.LowRest, i)
	} else {
		pc.LowRest = i
		pc.HasRest = true
	}
}

func (pc *ParamCases) signatures() string {
	var res []string
	for _, r := range pc.fixedRanges() {
		if pc.HasRest && r[1] >= pc.LowRest-1 {
			res = append(res, formatOrMore(r[0]))
			return strings.Join(res, ", ")
		}
		res = append(res, formatRange(r))
	}

	if pc.HasRest {
		res = append(res, formatOrMore(pc.LowRest))
	}

	return strings.Join(res, ", ")
}

func (pc *ParamCases) fixedRanges() [][2]int {
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

func (pc *ParamCases) fixedSet() []int {
	var res []int
	for i := range len(pc.Fixed) * arityBits {
		index, offset := i/arityBits, i%arityBits
		if pc.Fixed[index]&(1<<offset) != 0 {
			res = append(res, i)
		}
	}
	return res
}

func parseParamCase(s data.Sequence) (*ParamCase, error) {
	f, body, _ := s.Split()
	argNames, restArg, err := parseParamNames(f)
	if err != nil {
		return nil, err
	}
	if body.IsEmpty() {
		return nil, fmt.Errorf(ErrNoCaseBodyDefined, f)
	}
	return &ParamCase{
		Signature: f,
		Params:    argNames,
		Rest:      restArg,
		Body:      body,
	}, nil
}

func (c *ParamCase) fixedArgs() data.Locals {
	if c.Rest {
		return c.Params[0 : len(c.Params)-1]
	}
	return c.Params
}

func (c *ParamCase) restArg() (data.Local, bool) {
	if c.Rest {
		return c.Params[len(c.Params)-1], true
	}
	return "", false
}

func (c *ParamCase) getArity() (int, bool) {
	fl := len(c.fixedArgs())
	if _, ok := c.restArg(); ok {
		return fl, true
	}
	return fl, false
}

func (c *ParamCase) makeArgFetcher() ArgFetcher {
	cl := len(c.Params)
	if c.Rest {
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

func parseParamNames(v data.Value) (data.Locals, bool, error) {
	switch v := v.(type) {
	case data.Local:
		return data.Locals{v}, true, nil
	case *data.List:
		return parseListParamNames(v), false, nil
	case *data.Cons:
		return parseConsParamNames(v), true, nil
	default:
		return nil, false, fmt.Errorf(ErrUnexpectedParamSyntax, v)
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
