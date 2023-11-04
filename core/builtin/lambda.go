package builtin

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/compiler/generate"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/runtime/isa"
)

type lambdaEncoder struct {
	encoder.Encoder
	cases paramCases
}

// Error messages
const (
	ErrUnexpectedCaseSyntax   = "unexpected case syntax: %s"
	ErrNoMatchingParamPattern = "no matching parameter pattern"
)

// Lambda encodes a lambda
func Lambda(e encoder.Encoder, args ...data.Value) {
	var le *lambdaEncoder
	cases := parseParamCases(data.NewVector(args...))
	fn := generate.Procedure(e, func(c encoder.Encoder) {
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
		le.Emit(isa.RetNull)
		return
	}
	le.encodeCases(le.cases)
}

func (le *lambdaEncoder) encodeCases(cases paramCases) {
	switch len(cases) {
	case 0:
		generate.Literal(le, data.String(ErrNoMatchingParamPattern))
		le.Emit(isa.Panic)
		return
	case 1:
		if c := cases[0]; c.rest && len(c.params) == 1 {
			le.consequent(c)
			return
		}
		fallthrough
	default:
		c := cases[0]
		generate.Branch(le,
			func(encoder.Encoder) { le.predicate(c) },
			func(encoder.Encoder) { le.consequent(c) },
			func(encoder.Encoder) { le.encodeCases(cases[1:]) },
		)
	}
}

func (le *lambdaEncoder) predicate(c *paramCase) {
	le.Emit(isa.ArgLen)
	al := len(c.params)
	if c.rest {
		generate.Literal(le, data.Integer(al-1))
		le.Emit(isa.NumGte)
		return
	}
	generate.Literal(le, data.Integer(al))
	le.Emit(isa.NumEq)
}

func (le *lambdaEncoder) consequent(c *paramCase) {
	body := c.body
	if body.IsEmpty() {
		le.Emit(isa.RetNull)
		return
	}

	le.PushParams(c.params, c.rest)
	le.PushLocals()
	generate.Block(le, c.body)
	le.Emit(isa.Return)
	le.PopLocals()
	le.PopParams()
}
