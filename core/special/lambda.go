package special

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/lang/params"
	"github.com/kode4food/ale/internal/runtime/isa"
)

// Lambda encodes a lambda
func Lambda(e encoder.Encoder, args ...ale.Value) error {
	pc, err := params.ParseCases(data.Vector(args))
	if err != nil {
		return err
	}
	fn, err := generate.Procedure(e, func(pe encoder.Encoder) error {
		return makeLambda(pe, pc)
	})
	if err != nil {
		return err
	}
	fn.ArityChecker = pc.MakeArityChecker()
	return nil
}

func makeLambda(e encoder.Encoder, pc *params.ParamCases) error {
	if len(pc.Cases) == 0 {
		e.Emit(isa.RetNull)
		return nil
	}
	return encodeCases(e, pc.Cases)

}

func encodeCases(e encoder.Encoder, cases []*params.ParamCase) error {
	switch len(cases) {
	case 0:
		noMatch := data.String(params.ErrNoMatchingParamPattern)
		if err := generate.Literal(e, noMatch); err != nil {
			return err
		}
		e.Emit(isa.Panic)
		return nil
	case 1:
		if c := cases[0]; c.Rest && len(c.Params) == 1 {
			return encodeConsequent(e, c)
		}
		fallthrough
	default:
		c := cases[0]
		return generate.Branch(e,
			func(encoder.Encoder) error { return encodePredicate(e, c) },
			func(encoder.Encoder) error { return encodeConsequent(e, c) },
			func(encoder.Encoder) error { return encodeCases(e, cases[1:]) },
		)
	}
}

func encodePredicate(e encoder.Encoder, c *params.ParamCase) error {
	e.Emit(isa.ArgsLen)
	cl := len(c.Params)
	if c.Rest {
		if err := generate.Literal(e, data.Integer(cl-1)); err != nil {
			return err
		}
		e.Emit(isa.NumGte)
		return nil
	}
	if err := generate.Literal(e, data.Integer(cl)); err != nil {
		return err
	}
	e.Emit(isa.NumEq)
	return nil
}

func encodeConsequent(e encoder.Encoder, c *params.ParamCase) error {
	e.PushParams(c.Params, c.Rest)
	e.PushLocals()
	if err := generate.Block(e, c.Body); err != nil {
		return err
	}
	e.Emit(isa.Return)
	if err := e.PopLocals(); err != nil {
		return err
	}
	e.PopParams()
	return nil
}
