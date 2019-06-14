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
		variants variants
	}

	variant struct {
		args data.Names
		rest bool
		body data.Sequence
	}

	variants []*variant
)

// Error messages
const (
	InvalidRestArgument = "rest-argument not well-formed: %s"
)

const (
	allArgsName = data.Name("*args*")
	restMarker  = data.Name(".")
)

// Lambda encodes a lambda
func Lambda(e encoder.Type, args ...data.Value) {
	vars := parseLambda(args)
	le := makeLambdaEncoder(e, vars)
	le.encodeCall()
}

func makeLambdaEncoder(e encoder.Type, v variants) *lambdaEncoder {
	child := e.Child()
	res := &lambdaEncoder{
		Type:     child,
		variants: v,
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
	if len(le.variants) == 0 {
		le.Emit(isa.RetNull)
	} else {
		le.makeVariants(le.variants)
	}
	res := vm.LambdaFromEncoder(le)
	res.ArityChecker = le.makeArityChecker()
	return res
}

func (le *lambdaEncoder) makeVariants(vars variants) {
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
	v0 := le.variants[0]
	lower, upper := v0.arityRange()
	for _, s := range le.variants[1:] {
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

func (le *lambdaEncoder) makeCond(v *variant) {
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

func (le *lambdaEncoder) makeThen(v *variant) {
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

func parseLambda(s data.Vector) variants {
	switch s.First().(type) {
	case data.Vector, data.LocalSymbol:
		v := parseLambdaVariant(s)
		return variants{v}
	default:
		var res variants
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			v := parseLambdaVariant(f.(data.List))
			res = append(res, v)
		}
		return res
	}
}

func parseLambdaVariant(s data.Sequence) *variant {
	f, body, _ := s.Split()
	argNames, restArg := parseArgBindings(f)
	return &variant{
		args: argNames,
		rest: restArg,
		body: body,
	}
}

func (v *variant) fixedArgs() data.Names {
	if v.rest {
		return v.args[0 : len(v.args)-1]
	}
	return v.args
}

func (v *variant) restArg() (data.Name, bool) {
	if v.rest {
		return v.args[len(v.args)-1], true
	}
	return "", false
}

func (v *variant) arityRange() (int, int) {
	fl := len(v.fixedArgs())
	if _, ok := v.restArg(); ok {
		return fl, -1
	}
	return fl, fl
}

func parseArgBindings(v data.Value) (data.Names, bool) {
	if l, ok := v.(data.LocalSymbol); ok {
		s := data.NewVector(data.NewLocalSymbol(restMarker), l)
		return parseArgNames(s)
	}
	return parseArgNames(v.(data.Vector))
}

func parseArgNames(s data.Sequence) (data.Names, bool) {
	var an data.Names
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		n := f.(data.LocalSymbol).Name()
		if n == restMarker {
			rn := parseRestArg(r)
			return append(an, rn), true
		}
		an = append(an, n)
	}
	return an, false
}

func parseRestArg(s data.Sequence) data.Name {
	if f, r, ok := s.Split(); ok {
		n := f.(data.LocalSymbol).Name()
		if n != restMarker && r.IsEmpty() {
			return n
		}
	}
	panic(fmt.Errorf(InvalidRestArgument, s))
}
