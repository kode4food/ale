package special

import (
	"fmt"

	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/util"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/ale/runtime/vm"
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
	ErrUnexpectedLambdaSyntax = "unexpected lambda syntax: %s"
)

const allArgsName = data.Name("*args*")

// Lambda encodes a lambda
func Lambda(e encoder.Encoder, args ...data.Value) {
	vars := parseLambda(args)
	le := makeLambdaEncoder(e, vars)
	le.encodeCall()
}

func makeLambdaEncoder(e encoder.Encoder, v lambdaCases) *lambdaEncoder {
	child := e.Child()
	res := &lambdaEncoder{
		Encoder: child,
		cases:   v,
	}
	res.PushArgs(data.Names{allArgsName}, true)
	return res
}

func (le *lambdaEncoder) encodeCall() {
	e := le.Parent()
	fn := le.makeLambda()

	cells := le.Closure()
	nl := len(cells)
	if nl == 0 {
		// nothing needed to be captured from local variables,
		// so just pass the newly instantiated closure through
		generate.Literal(e, fn.Call())
		return
	}

	for i := nl - 1; i >= 0; i-- {
		name := cells[i].Name
		generate.Symbol(e, data.NewLocalSymbol(name))
	}
	e.Emit(isa.Const, e.AddConstant(fn))
	e.Emit(isa.Call, isa.Count(nl))
}

func (le *lambdaEncoder) makeLambda() *vm.Lambda {
	if len(le.cases) == 0 {
		le.Emit(isa.RetNil)
	} else {
		le.makeLambdaCases(le.cases)
	}
	res := vm.LambdaFromEncoder(le)
	res.ArityChecker = le.makeArityChecker()
	return res
}

func (le *lambdaEncoder) makeLambdaCases(cases lambdaCases) {
	if len(cases) == 0 {
		generate.Literal(le, data.String("no matching argument pattern"))
		le.Emit(isa.Panic)
		return
	}

	c := cases[0]
	generate.Branch(le,
		func() { le.makePredicate(c) },
		func() { le.makeConsequent(c) },
		func() { le.makeLambdaCases(cases[1:]) },
	)
}

func (le *lambdaEncoder) makeArityChecker() data.ArityChecker {
	v0 := le.cases[0]
	lower, upper := v0.arityRange()
	for _, s := range le.cases[1:] {
		l, u := s.arityRange()
		lower = util.IntMin(l, lower)
		if u == data.OrMore || upper == data.OrMore {
			upper = data.OrMore
			continue
		}
		upper = util.IntMax(u, upper)
	}
	return data.MakeChecker(lower, upper)
}

func (le *lambdaEncoder) makePredicate(c *lambdaCase) {
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

func (le *lambdaEncoder) makeConsequent(c *lambdaCase) {
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
		panic(fmt.Errorf(ErrUnexpectedLambdaSyntax, v))
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
