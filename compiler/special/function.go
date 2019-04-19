package special

import (
	"fmt"

	"gitlab.com/kode4food/ale/compiler/arity"
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/compiler/generate"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/macro"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
	"gitlab.com/kode4food/ale/runtime/vm"
)

type (
	funcEncoder struct {
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
	restMarker  = data.Name("&")
)

// Fn encodes a lambda
func Fn(e encoder.Type, args ...data.Value) {
	name, vars := parseFunction(args)
	fe := makeFunctionEncoder(e, name, vars)
	arityChecker := fe.makeArityChecker()
	fe.encodeCall()
	generate.Literal(e, data.Call(func(args ...data.Value) data.Value {
		return &data.Function{
			Call:         args[0].(data.Call),
			Convention:   data.ApplicativeCall,
			ArityChecker: arityChecker,
		}
	}))
	e.Emit(isa.Call1)
}

// DefMacro encodes and registers a macro
func DefMacro(e encoder.Type, args ...data.Value) {
	name, vars := parseNamedFunction(args)
	fe := makeFunctionEncoder(e, name, vars)
	arityChecker := fe.makeArityChecker()
	fe.encodeCall()
	generate.Literal(e, data.Call(func(args ...data.Value) data.Value {
		body := args[0].(data.Call)
		wrapper := func(_ namespace.Type, args ...data.Value) data.Value {
			if err := arityChecker(len(args)); err != nil {
				panic(err)
			}
			return body(args...)
		}
		return macro.Call(wrapper)
	}))
	e.Emit(isa.Call1)
	generate.Literal(e, fe.Name())
	e.Emit(isa.Bind)
	generate.Literal(e, fe.Name())
}

func makeFunctionEncoder(e encoder.Type, n data.Name, v variants) *funcEncoder {
	child := makeChildEncoder(e, n)
	res := &funcEncoder{
		Type:     child,
		variants: v,
	}
	res.PushArgs(data.Names{allArgsName}, true)
	return res
}

func makeChildEncoder(e encoder.Type, n data.Name) encoder.Type {
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
		generate.Symbol(e, data.NewLocalSymbol(name))
	}
	e.Emit(isa.Const, idx)
	e.Emit(isa.Call, isa.Count(nl))
}

func (fe *funcEncoder) makeClosure() data.Call {
	if len(fe.variants) == 0 {
		fe.Emit(isa.RetNil)
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
		generate.Literal(fe, data.String("no matching argument pattern"))
		fe.Emit(isa.Panic)
		return
	}

	v := vars[0]
	generate.Branch(fe,
		func() { fe.makeCond(v) },
		func() { fe.makeThen(v) },
		func() { fe.makeVariants(vars[1:]) },
	)
}

func (fe *funcEncoder) makeArityChecker() data.ArityChecker {
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
	fe.Emit(isa.ArgLen)
	al := len(v.args)
	if v.rest {
		generate.Literal(fe, data.Integer(al-1))
		fe.Emit(isa.Gte)
		return
	}
	generate.Literal(fe, data.Integer(al))
	fe.Emit(isa.Eq)
}

func (fe *funcEncoder) makeThen(v *variant) {
	body := v.body
	if body.IsEmpty() {
		fe.Emit(isa.RetNil)
		return
	}

	fe.PushArgs(v.args, v.rest)
	fe.PushLocals()
	generate.Block(fe, v.body)
	fe.Emit(isa.Return)
	fe.PopLocals()
	fe.PopArgs()
}

func parseNamedFunction(args data.Vector) (data.Name, variants) {
	name := args[0].(data.LocalSymbol).Name()
	vars := parseFunctionVariants(args[1:])
	return name, vars
}

func parseFunction(args data.Vector) (data.Name, variants) {
	name, r := parseOptionalName(args)
	vars := parseFunctionVariants(r)
	return name, vars
}

func parseOptionalName(args data.Vector) (data.Name, data.Vector) {
	if s, ok := args[0].(data.Symbol); ok {
		ls := s.(data.LocalSymbol)
		return ls.Name(), args[1:]
	}
	return "", args
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

func parseFunctionVariants(s data.Sequence) variants {
	if _, ok := s.First().(data.Vector); ok {
		v := parseFunctionVariant(s)
		return variants{v}
	}
	var res variants
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		v := parseFunctionVariant(f.(*data.List))
		res = append(res, v)
	}
	return res
}

func parseFunctionVariant(s data.Sequence) *variant {
	f, body, _ := s.Split()
	argNames, restArg := parseArgNames(f.(data.Vector))
	return &variant{
		args: argNames,
		rest: restArg,
		body: body,
	}
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
		n := f.(data.Symbol).Name()
		if n != restMarker && r.IsEmpty() {
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
