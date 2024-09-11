package asm

import (
	"errors"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/special"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/pkg/core/internal"
	"github.com/kode4food/ale/pkg/data"
)

func (p *asmParser) specialCall(forms data.Sequence) (asmEmit, error) {
	pc := internal.ParseParamCases(forms)
	cases := pc.Cases
	ap := make([]*asmParser, len(cases))
	emitters := make([]asmEmit, len(cases))
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

	fn := func(e encoder.Encoder, args ...data.Value) {
		if err := ac(len(args)); err != nil {
			panic(err)
		}
		for i, f := range fetchers {
			if a, ok := f(args); ok {
				ae := ap[i].wrapEncoder(e, a...)
				if err := emitters[i](ae); err != nil {
					panic(err)
				}
				return
			}
		}
		panic(errors.New(internal.ErrNoMatchingParamPattern))
	}

	return func(e *asmEncoder) error {
		e.Emit(isa.Const, e.AddConstant(special.Call(fn)))
		return nil
	}, nil
}
