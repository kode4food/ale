package isa_test

import (
	"testing"

	"gitlab.com/kode4food/ale/internal/assert"
	"gitlab.com/kode4food/ale/runtime/isa"
)

func TestEffects(t *testing.T) {
	as := assert.New(t)

	e1 := isa.MustGetEffect(isa.CondJump)
	as.Number(2, e1.Size)

	defer func() {
		rec := recover()
		err := "effect not declared for opcode: Opcode(5000)"
		as.String(err, rec)
	}()

	isa.MustGetEffect(isa.Opcode(5000))
}
