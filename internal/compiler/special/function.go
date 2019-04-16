package special

import (
	"fmt"

	"gitlab.com/kode4food/ale/internal/macro"
	"gitlab.com/kode4food/ale/internal/namespace"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/compiler/arity"
	"gitlab.com/kode4food/ale/internal/compiler/build"
	"gitlab.com/kode4food/ale/internal/compiler/encoder"
	"gitlab.com/kode4food/ale/internal/compiler/generate"
	"gitlab.com/kode4food/ale/internal/runtime/isa"
	"gitlab.com/kode4food/ale/internal/runtime/vm"
)

type (
	funcEncoder struct {
		encoder.Type
		variants variants
	}

	variant struct {
		args api.Names
		rest bool
		body api.Sequence
	}

	variants []*variant
)

// Error messages
const (
	InvalidRestArgument = "rest-argument not well-formed: %s"
)

const (
	allArgsName = api.Name("*args*")
	restMarker  = api.Name("&")
)

// Fn encodes a lambda
func Fn(e encoder.Type, args ...api.Value) {
	name, vars := parseFunction(args)
	fe := makeFunctionEncoder(e, name, vars)
	arityChecker := fe.makeArityChecker()
	fe.encodeCall()
	generate.Literal(e, api.Call(func(args ...api.Value) api.Value {
		return &api.Function{
			Call:         args[0].(api.Call),
			Convention:   api.ApplicativeCall,
			ArityChecker: arityChecker,
		}
	}))
	e.Append(isa.Call1)
}

// DefMacro encodes and registers a macro
func DefMacro(e encoder.Type, args ...api.Value) {
	name, vars := parseNamedFunction(args)
	fe := makeFunctionEncoder(e, name, vars)
	arityChecker := fe.makeArityChecker()
	fe.encodeCall()
	generate.Literal(e, api.Call(func(args ...api.Value) api.Value {
		body := args[0].(api.Call)
		wrapper := func(_ namespace.Type, args ...api.Value) api.Value {
			if err := arityChecker(len(args)); err != nil {
				panic(err)
			}
			return body(args...)
		}
		return macro.Call(wrapper)
	}))
	e.Append(isa.Call1)
	generate.Literal(e, fe.Name())
	e.Append(isa.Bind)
	generate.Literal(e, fe.Name())
}

func makeFunctionEncoder(e encoder.Type, n api.Name, v variants) *funcEncoder {
	child := makeChildEncoder(e, n)
	res := &funcEncoder{
		Type:     child,
		variants: v,
	}
	res.PushArgs(api.Names{allArgsName}, true)
	return res
}

func makeChildEncoder(e encoder.Type, n api.Name) encoder.Type {
	if n != "" {
		return e.NamedChild(n)
	}
	return e.Child()
}

func (fe *funcEncoder) encodeCall() {
	e := fe.Parent()
	fn := fe.makeClosure()
	names := fe.Closure()
	nl := len(names)
	if nl == 0 {
		generate.Literal(e, fn())
		return
	}

	idx := e.AddConstant(fn)
	for i := nl - 1; i >= 0; i-- {
		name := names[i]
		generate.Symbol(e, api.NewLocalSymbol(name))
	}
	e.Append(isa.Const, idx)
	e.Append(isa.Call, isa.Count(nl))
}

func (fe *funcEncoder) makeClosure() api.Call {
	if len(fe.variants) == 0 {
		fe.Append(isa.ReturnNil)
	} else {
		fe.makeVariants(fe.variants)
	}
	return vm.NewClosure(&vm.Config{
		Globals:    fe.Globals(),
		Code:       fe.Code(),
		Constants:  fe.Constants(),
		StackSize:  fe.StackSize(),
		LocalCount: fe.LocalCount(),
	})
}

func (fe *funcEncoder) makeVariants(vars variants) {
	if len(vars) == 0 {
		generate.Literal(fe, api.String("no matching argument pattern"))
		fe.Append(isa.Panic)
		return
	}

	v := vars[0]
	build.Cond(fe,
		func() { fe.makeCond(v) },
		func() { fe.makeThen(v) },
		func() { fe.makeVariants(vars[1:]) },
	)
}

func (fe *funcEncoder) makeArityChecker() api.ArityChecker {
	v0 := fe.variants[0]
	lower, upper := v0.arityRange()
	for _, s := range fe.variants[1:] {
		l, u := s.arityRange()
		lower = min(l, lower)
		if u == -1 || upper == -1 {
			upper = -1
			continue
		}
		upper = max(u, upper)
	}
	return arity.MakeChecker(lower, upper)
}

func (fe *funcEncoder) makeCond(v *variant) {
	fe.Append(isa.ArgLen)
	al := len(v.args)
	if v.rest {
		generate.Literal(fe, api.Integer(al-1))
		fe.Append(isa.Gte)
		return
	}
	generate.Literal(fe, api.Integer(al))
	fe.Append(isa.Eq)
}

func (fe *funcEncoder) makeThen(v *variant) {
	body := v.body
	if !body.IsSequence() {
		fe.Append(isa.ReturnNil)
		return
	}

	fe.PushArgs(v.args, v.rest)
	fe.PushLocals()
	generate.Block(fe, v.body)
	fe.Append(isa.Return)
	fe.PopLocals()
	fe.PopArgs()
}

func parseNamedFunction(args api.Vector) (api.Name, variants) {
	name := args[0].(api.LocalSymbol).Name()
	vars := parseFunctionVariants(args[1:])
	return name, vars
}

func parseFunction(args api.Vector) (api.Name, variants) {
	name, r := parseOptionalName(args)
	vars := parseFunctionVariants(r)
	return name, vars
}

func parseOptionalName(args api.Vector) (api.Name, api.Vector) {
	if s, ok := args[0].(api.Symbol); ok {
		ls := s.(api.LocalSymbol)
		return ls.Name(), args[1:]
	}
	return "", args
}

func (v *variant) fixedArgs() api.Names {
	if v.rest {
		return v.args[0 : len(v.args)-1]
	}
	return v.args
}

func (v *variant) restArg() (api.Name, bool) {
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

func parseFunctionVariants(s api.Sequence) variants {
	if _, ok := s.First().(api.Vector); ok {
		v := parseFunctionVariant(s)
		return variants{v}
	}
	var res variants
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		v := parseFunctionVariant(f.(*api.List))
		res = append(res, v)
	}
	return res
}

func parseFunctionVariant(s api.Sequence) *variant {
	f, body, _ := s.Split()
	argNames, restArg := parseArgNames(f.(api.Vector))
	return &variant{
		args: argNames,
		rest: restArg,
		body: body,
	}
}

func parseArgNames(s api.Sequence) (api.Names, bool) {
	var an api.Names
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		n := f.(api.LocalSymbol).Name()
		if n == restMarker {
			rn := parseRestArg(r)
			return append(an, rn), true
		}
		an = append(an, n)
	}
	return an, false
}

func parseRestArg(s api.Sequence) api.Name {
	if f, r, ok := s.Split(); ok {
		n := f.(api.Symbol).Name()
		if n != restMarker && !r.IsSequence() {
			return n
		}
	}
	panic(fmt.Errorf(InvalidRestArgument, s))
}

func min(left, right int) int {
	if left < right {
		return left
	}
	return right
}

func max(left, right int) int {
	if left > right {
		return left
	}
	return right
}
