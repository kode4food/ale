package special

import (
	"fmt"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

type (
	lambdaEncoder struct {
		encoder.Encoder
		cases lambdaCases
	}

	lambdaCase struct {
		args data.Names
		rest bool
		body data.Sequence
	}

	lambdaCases []*lambdaCase
)

// Error messages
const (
	ErrUnexpectedLambdaSyntax   = "unexpected lambda syntax: %s"
	ErrUnexpectedArgumentSyntax = "unexpected argument syntax: %s"
)

// Lambda encodes a lambda
func Lambda(e encoder.Encoder, args ...data.Value) {
	var le *lambdaEncoder
	vars := parseLambda(data.NewVector(args...))
	fn := generate.Lambda(e, func(c encoder.Encoder) {
		le = makeLambda(c, vars)
		le.encode()
	})
	fn.ArityChecker = le.makeArityChecker()
}

func makeLambda(e encoder.Encoder, v lambdaCases) *lambdaEncoder {
	res := &lambdaEncoder{
		Encoder: e,
		cases:   v,
	}
	return res
}

func (le *lambdaEncoder) encode() {
	if len(le.cases) == 0 {
		le.Emit(isa.RetNil)
		return
	}
	le.encodeCases(le.cases)
}

func (le *lambdaEncoder) encodeCases(cases lambdaCases) {
	if len(cases) == 0 {
		generate.Literal(le, data.String("no matching argument pattern"))
		le.Emit(isa.Panic)
		return
	}

	c := cases[0]
	generate.Branch(le,
		func(encoder.Encoder) { le.predicate(c) },
		func(encoder.Encoder) { le.consequent(c) },
		func(encoder.Encoder) { le.encodeCases(cases[1:]) },
	)
}

func (le *lambdaEncoder) makeArityChecker() data.ArityChecker {
	v0 := le.cases[0]
	lower, upper := v0.arityRange()
	for _, s := range le.cases[1:] {
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

func (le *lambdaEncoder) predicate(c *lambdaCase) {
	le.Emit(isa.ArgLen)
	al := len(c.args)
	if c.rest {
		generate.Literal(le, data.Integer(al-1))
		le.Emit(isa.Gte)
		return
	}
	generate.Literal(le, data.Integer(al))
	le.Emit(isa.Eq)
}

func (le *lambdaEncoder) consequent(c *lambdaCase) {
	body := c.body
	if body.IsEmpty() {
		le.Emit(isa.RetNil)
		return
	}

	le.PushArgs(c.args, c.rest)
	le.PushLocals()
	generate.Block(le, c.body)
	le.Emit(isa.Return)
	le.PopLocals()
	le.PopArgs()
}

func parseLambda(s data.Vector) lambdaCases {
	f := s.First()
	switch f.(type) {
	case data.List, data.Cons, data.LocalSymbol:
		c := parseLambdaCase(s)
		return lambdaCases{c}
	case data.Vector:
		var res lambdaCases
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			c := parseLambdaCase(f.(data.Vector))
			res = append(res, c)
		}
		return res
	default:
		panic(fmt.Errorf(ErrUnexpectedLambdaSyntax, f))
	}
}

func parseLambdaCase(s data.Sequence) *lambdaCase {
	f, body, _ := s.Split()
	argNames, restArg := parseArgBindings(f)
	return &lambdaCase{
		args: argNames,
		rest: restArg,
		body: body,
	}
}

func (c *lambdaCase) fixedArgs() data.Names {
	if c.rest {
		return c.args[0 : len(c.args)-1]
	}
	return c.args
}

func (c *lambdaCase) restArg() (data.Name, bool) {
	if c.rest {
		return c.args[len(c.args)-1], true
	}
	return "", false
}

func (c *lambdaCase) arityRange() (int, int) {
	fl := len(c.fixedArgs())
	if _, ok := c.restArg(); ok {
		return fl, data.OrMore
	}
	return fl, fl
}

func parseArgBindings(v data.Value) (data.Names, bool) {
	switch v := v.(type) {
	case data.LocalSymbol:
		return data.Names{v.Name()}, true
	case data.List:
		return parseListArgNames(v), false
	case data.Cons:
		return parseConsArgNames(v), true
	default:
		panic(fmt.Errorf(ErrUnexpectedArgumentSyntax, v))
	}
}

func parseListArgNames(l data.List) data.Names {
	var an data.Names
	for f, r, ok := l.Split(); ok; f, r, ok = r.Split() {
		n := f.(data.LocalSymbol).Name()
		an = append(an, n)
	}
	return an
}

func parseConsArgNames(c data.Cons) data.Names {
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
