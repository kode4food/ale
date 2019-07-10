package special

import (
	"fmt"

	"gitlab.com/kode4food/ale/compiler/arity"
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/compiler/generate"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/util"
	"gitlab.com/kode4food/ale/runtime/isa"
	"gitlab.com/kode4food/ale/runtime/vm"
)

type (
	lambdaEncoder struct {
		encoder.Type
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
	UnexpectedLambdaSyntax = "unexpected lambda syntax: %s"
)

const allArgsName = data.Name("*args*")

// Lambda encodes a lambda
func Lambda(e encoder.Type, args ...data.Value) {
	vars := parseLambda(args)
	le := makeLambdaEncoder(e, vars)
	le.encodeCall()
}

func makeLambdaEncoder(e encoder.Type, v lambdaCases) *lambdaEncoder {
	child := e.Child()
	res := &lambdaEncoder{
		Type:  child,
		cases: v,
	}
	res.PushArgs(data.Names{allArgsName}, true)
	return res
}

func (le *lambdaEncoder) encodeCall() {
	e := le.Parent()
	fn := le.makeLambda().Caller()

	cells := le.Closure()
	nl := len(cells)
	if nl == 0 {
		// nothing needed to be captured from local variables,
		// so just pass the newly instantiated closure through
		generate.Literal(e, fn())
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
		le.Emit(isa.RetNull)
	} else {
		le.makeVariants(le.cases)
	}
	res := vm.LambdaFromEncoder(le)
	res.ArityChecker = le.makeArityChecker()
	return res
}

func (le *lambdaEncoder) makeVariants(vars lambdaCases) {
	if len(vars) == 0 {
		generate.Literal(le, data.String("no matching argument pattern"))
		le.Emit(isa.Panic)
		return
	}

	v := vars[0]
	generate.Branch(le,
		func() { le.makeCond(v) },
		func() { le.makeThen(v) },
		func() { le.makeVariants(vars[1:]) },
	)
}

func (le *lambdaEncoder) makeArityChecker() data.ArityChecker {
	v0 := le.cases[0]
	lower, upper := v0.arityRange()
	for _, s := range le.cases[1:] {
		l, u := s.arityRange()
		lower = util.IntMin(l, lower)
		if u == -1 || upper == -1 {
			upper = -1
			continue
		}
		upper = util.IntMax(u, upper)
	}
	return arity.MakeChecker(lower, upper)
}

func (le *lambdaEncoder) makeCond(v *lambdaCase) {
	le.Emit(isa.ArgLen)
	al := len(v.args)
	if v.rest {
		generate.Literal(le, data.Integer(al-1))
		le.Emit(isa.Gte)
		return
	}
	generate.Literal(le, data.Integer(al))
	le.Emit(isa.Eq)
}

func (le *lambdaEncoder) makeThen(v *lambdaCase) {
	body := v.body
	if body.IsEmpty() {
		le.Emit(isa.RetNull)
		return
	}

	le.PushArgs(v.args, v.rest)
	le.PushLocals()
	generate.Block(le, v.body)
	le.Emit(isa.Return)
	le.PopLocals()
	le.PopArgs()
}

func parseLambda(s data.Vector) lambdaCases {
	f := s.First()
	switch f.(type) {
	case data.List, *data.Cons, data.LocalSymbol:
		v := parseLambdaCase(s)
		return lambdaCases{v}
	case data.Vector:
		var res lambdaCases
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			v := parseLambdaCase(f.(data.Vector))
			res = append(res, v)
		}
		return res
	default:
		panic(fmt.Errorf(UnexpectedLambdaSyntax, f))
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

func (v *lambdaCase) fixedArgs() data.Names {
	if v.rest {
		return v.args[0 : len(v.args)-1]
	}
	return v.args
}

func (v *lambdaCase) restArg() (data.Name, bool) {
	if v.rest {
		return v.args[len(v.args)-1], true
	}
	return "", false
}

func (v *lambdaCase) arityRange() (int, int) {
	fl := len(v.fixedArgs())
	if _, ok := v.restArg(); ok {
		return fl, -1
	}
	return fl, fl
}

func parseArgBindings(v data.Value) (data.Names, bool) {
	switch typed := v.(type) {
	case data.LocalSymbol:
		return data.Names{typed.Name()}, true
	case *data.Cons:
		return parseConsArgNames(typed), true
	case data.List:
		return parseListArgNames(typed), false
	default:
		panic("what the shit?")
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

func parseConsArgNames(c *data.Cons) data.Names {
	var an data.Names
	next := c
	for {
		an = append(an, next.Car().(data.LocalSymbol).Name())

		cdr := next.Cdr()
		if nc, ok := cdr.(*data.Cons); ok {
			next = nc
			continue
		}

		an = append(an, cdr.(data.LocalSymbol).Name())
		return an
	}
}
