package asm

import (
	"errors"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/compiler"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/lang/params"
	"github.com/kode4food/ale/internal/runtime/isa"
)

// Error messages
const (
	// ErrUnexpectedParameter is raised when an encoder parameter is not found.
	// These are declared using the special* built-in
	ErrUnexpectedParameter = "unexpected parameter name: %s"
)

func MakeSpecial(pc *params.ParamCases) EmitBuilder {
	return func(p *Parser) (Emit, error) {
		cases := pc.Cases
		ap := make([]*Parser, len(cases))
		emitters := make([]Emit, len(cases))
		for i, c := range cases {
			ap[i] = p.withParams(c.Params)
			e, err := ap[i].sequence(c.Body)
			if err != nil {
				return nil, err
			}
			emitters[i] = e
		}

		ac := pc.MakeArityChecker()
		fetchers := pc.MakeArgFetchers()

		fn := func(e encoder.Encoder, args ...ale.Value) error {
			if err := ac(len(args)); err != nil {
				return err
			}
			for i, f := range fetchers {
				if a, ok := f(args); ok {
					ae := ap[i].wrapEncoder(e, a...)
					return emitters[i](ae)
				}
			}
			return errors.New(params.ErrNoMatchingParamPattern)
		}

		return func(e *Encoder) error {
			e.Emit(isa.Const, e.AddConstant(compiler.Call(fn)))
			return nil
		}, nil
	}
}
