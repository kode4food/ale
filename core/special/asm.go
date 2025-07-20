package special

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/compiler/asm"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/lang/params"
)

// Asm provides indirect access to the Encoder's methods and generators
func Asm(e encoder.Encoder, args ...ale.Value) error {
	return asm.Encode(e, asm.MakeAsm(args...))
}

// Special emits an encoder function for the provided param cases
func Special(e encoder.Encoder, args ...ale.Value) error {
	pc, err := params.ParseCases(data.Vector(args))
	if err != nil {
		return err
	}
	return asm.Encode(e, asm.MakeSpecial(pc))
}
