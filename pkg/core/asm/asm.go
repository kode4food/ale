package asm

import (
	"sync"

	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/pkg/data"
)

var (
	callsOnce sync.Once
	calls     namedAsmParsers
)

// Asm provides indirect access to the Encoder's methods and generators
func Asm(e encoder.Encoder, args ...data.Value) {
	c := getCalls()
	p := makeAsmParser(c)
	emit, err := p.parse(data.Vector(args))
	if err != nil {
		panic(err)
	}
	ae := p.wrapEncoder(e)
	if err := emit(ae); err != nil {
		panic(err)
	}
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
