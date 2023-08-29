package special

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

type (
	lambdaEncoder struct {
		encoder.Encoder
		cases paramCases
	}
)

// Error messages
const (
	ErrUnexpectedCaseSyntax      = "unexpected case syntax: %s"
	ErrNoMatchingArgumentPattern = "no matching argument pattern"
)

// Lambda encodes a lambda
func Lambda(e encoder.Encoder, args ...data.Value) {
	var le *lambdaEncoder
	cases := parseParamCases(data.NewVector(args...))
	fn := generate.Lambda(e, func(c encoder.Encoder) {
		le = makeLambda(c, cases)
		le.encode()
	})
	fn.ArityChecker = cases.makeArityChecker()
}

func makeLambda(e encoder.Encoder, v paramCases) *lambdaEncoder {
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

func (le *lambdaEncoder) encodeCases(cases paramCases) {
	if len(cases) == 0 {
		generate.Literal(le, data.String(ErrNoMatchingArgumentPattern))
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

func (le *lambdaEncoder) predicate(c *paramCase) {
	le.Emit(isa.ArgLen)
	al := len(c.params)
	if c.rest {
		generate.Literal(le, data.Integer(al-1))
		le.Emit(isa.Gte)
		return
	}
	generate.Literal(le, data.Integer(al))
	le.Emit(isa.Eq)
}

func (le *lambdaEncoder) consequent(c *paramCase) {
	body := c.body
	if body.IsEmpty() {
		le.Emit(isa.RetNil)
		return
	}

	le.PushArgs(c.params, c.rest)
	le.PushLocals()
	generate.Block(le, c.body)
	le.Emit(isa.Return)
	le.PopLocals()
	le.PopArgs()
}
