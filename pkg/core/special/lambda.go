package special

import (
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/core/internal"
	"github.com/kode4food/ale/pkg/data"
)

type lambdaEncoder struct {
	encoder.Encoder
	cases *internal.ParamCases
}

// Lambda encodes a lambda
func Lambda(e encoder.Encoder, args ...data.Value) {
	var le *lambdaEncoder
	pc := internal.ParseParamCases(data.Vector(args))
	fn := generate.Procedure(e, func(c encoder.Encoder) {
		le = makeLambda(c, pc)
		le.encode()
	})
	fn.ArityChecker = pc.MakeArityChecker()
}

func makeLambda(e encoder.Encoder, v *internal.ParamCases) *lambdaEncoder {
	res := &lambdaEncoder{
		Encoder: e,
		cases:   v,
	}
	return res
}

func (le *lambdaEncoder) encode() {
	cases := le.cases.Cases
	if len(cases) == 0 {
		le.Emit(isa.RetNull)
		return
	}
	le.encodeCases(cases)
}

func (le *lambdaEncoder) encodeCases(cases []*internal.ParamCase) {
	switch len(cases) {
	case 0:
		generate.Literal(le, data.String(internal.ErrNoMatchingParamPattern))
		le.Emit(isa.Panic)
		return
	case 1:
		if c := cases[0]; c.Rest && len(c.Params) == 1 {
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

func (le *lambdaEncoder) predicate(c *internal.ParamCase) {
	le.Emit(isa.ArgLen)
	cl := len(c.Params)
	if c.Rest {
		generate.Literal(le, data.Integer(cl-1))
		le.Emit(isa.NumGte)
		return
	}
	generate.Literal(le, data.Integer(cl))
	le.Emit(isa.NumEq)
}

func (le *lambdaEncoder) consequent(c *internal.ParamCase) {
	body := c.Body
	if body.IsEmpty() {
		le.Emit(isa.RetNull)
		return
	}

	le.PushParams(c.Params, c.Rest)
	le.PushLocals()
	generate.Block(le, c.Body)
	le.Emit(isa.Return)
	le.PopLocals()
	le.PopParams()
}
