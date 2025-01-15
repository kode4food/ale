package asm

import (
	"sync"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/pkg/data"
)

type emitParser func(*asmParser, data.Sequence) (asmEmit, error)

var (
	callsOnce sync.Once
	calls     namedAsmParsers
)

// Asm provides indirect access to the Encoder's methods and generators
func Asm(e encoder.Encoder, args ...data.Value) error {
	return encodeForm(e, emitAsm, args...)
}

func emitAsm(p *asmParser, forms data.Sequence) (asmEmit, error) {
	return p.sequence(forms)
}

func encodeForm(e encoder.Encoder, fn emitParser, args ...data.Value) error {
	c := getCalls()
	p := makeAsmParser(c)
	emit, err := fn(p, data.Vector(args))
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
				panic(debug.ProgrammerError("duplicate entry: %s", k))
			}
			res[k] = v
		}
	}
	return res
}
