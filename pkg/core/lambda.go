package core

import (
	"github.com/kode4food/ale/pkg/compiler/encoder"
	"github.com/kode4food/ale/pkg/compiler/generate"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/runtime/isa"
)

type lambdaEncoder struct {
	encoder.Encoder
	cases *paramCases
}

// ErrNoMatchingParamPattern is raised when none of the parameter patterns for
// a Lambda were capable of being matched
const ErrNoMatchingParamPattern = "no matching parameter pattern"

// Lambda encodes a lambda
func Lambda(e encoder.Encoder, args ...data.Value) {
	var le *lambdaEncoder
	pc := parseParamCases(data.Vector(args))
	fn := generate.Procedure(e, func(c encoder.Encoder) {
		le = makeLambda(c, pc)
		le.encode()
	})
	fn.ArityChecker = pc.makeChecker()
}

func makeLambda(e encoder.Encoder, v *paramCases) *lambdaEncoder {
	res := &lambdaEncoder{
		Encoder: e,
		cases:   v,
	}
	return res
}

func (le *lambdaEncoder) encode() {
	cases := le.cases.Cases()
	if len(cases) == 0 {
		le.Emit(isa.RetNull)
		return
	}
	le.encodeCases(cases)
}

func (le *lambdaEncoder) encodeCases(cases []*paramCase) {
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
	cl := len(c.params)
	if c.rest {
		generate.Literal(le, data.Integer(cl-1))
		le.Emit(isa.NumGte)
		return
	}
	generate.Literal(le, data.Integer(cl))
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
