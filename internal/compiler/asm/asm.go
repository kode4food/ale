package asm

import (
	"sync"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/debug"
)

type (
	EmitBuilder func(p *Parser) (Emit, error)
	Emit        func(*Encoder) error
)

var (
	callsOnce sync.Once
	calls     namedAsmParsers
)

func MakeAsm(args ...ale.Value) EmitBuilder {
	forms := data.Vector(args)
	return func(p *Parser) (Emit, error) {
		return p.sequence(forms)
	}
}

func Encode(e encoder.Encoder, build EmitBuilder) error {
	c := getCalls()
	p := makeAsmParser(c)
	emit, err := build(p)
	if err != nil {
		return err
	}
	ae := p.wrapEncoder(e)
	return emit(ae)
}

func getCalls() namedAsmParsers {
	callsOnce.Do(func() {
		calls = mergeCalls(
			getInstructionCalls(),
			getDirectiveCalls(),
		)
	})
	return calls
}

func mergeCalls(maps ...namedAsmParsers) namedAsmParsers {
	res := namedAsmParsers{}
	for _, m := range maps {
		for k, v := range m {
			if _, ok := res[k]; ok {
				panic(debug.ProgrammerErrorf("duplicate entry: %s", k))
			}
			res[k] = v
		}
	}
	return res
}
