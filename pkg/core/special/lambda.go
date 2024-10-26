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
func Lambda(e encoder.Encoder, args ...data.Value) error {
	var le *lambdaEncoder
	pc := internal.MustParseParamCases(data.Vector(args))
	fn, err := generate.Procedure(e, func(c encoder.Encoder) error {
		le = makeLambda(c, pc)
		return le.encode()
	})
	if err != nil {
		return err
	}
	fn.ArityChecker = pc.MakeArityChecker()
	return nil
}

func makeLambda(e encoder.Encoder, v *internal.ParamCases) *lambdaEncoder {
	res := &lambdaEncoder{
		Encoder: e,
		cases:   v,
	}
	return res
}

func (le *lambdaEncoder) encode() error {
	cases := le.cases.Cases
	if len(cases) == 0 {
		le.Emit(isa.RetNull)
		return nil
	}
	return le.encodeCases(cases)
}

func (le *lambdaEncoder) encodeCases(cases []*internal.ParamCase) error {
	switch len(cases) {
	case 0:
		noMatch := data.String(internal.ErrNoMatchingParamPattern)
		if err := generate.Literal(le, noMatch); err != nil {
			return err
		}
		le.Emit(isa.Panic)
		return nil
	case 1:
		if c := cases[0]; c.Rest && len(c.Params) == 1 {
			return le.consequent(c)
		}
		fallthrough
	default:
		c := cases[0]
		return generate.Branch(le,
			func(encoder.Encoder) error { return le.predicate(c) },
			func(encoder.Encoder) error { return le.consequent(c) },
			func(encoder.Encoder) error { return le.encodeCases(cases[1:]) },
		)
	}
}

func (le *lambdaEncoder) predicate(c *internal.ParamCase) error {
	le.Emit(isa.ArgLen)
	cl := len(c.Params)
	if c.Rest {
		if err := generate.Literal(le, data.Integer(cl-1)); err != nil {
			return err
		}
		le.Emit(isa.NumGte)
		return nil
	}
	if err := generate.Literal(le, data.Integer(cl)); err != nil {
		return err
	}
	le.Emit(isa.NumEq)
	return nil
}

func (le *lambdaEncoder) consequent(c *internal.ParamCase) error {
	le.PushParams(c.Params, c.Rest)
	le.PushLocals()
	if err := generate.Block(le, c.Body); err != nil {
		return err
	}
	le.Emit(isa.Return)
	if err := le.PopLocals(); err != nil {
		return err
	}
	le.PopParams()
	return nil
}
