package isa_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/runtime/isa"
)

func TestEffects(t *testing.T) {
	as := assert.New(t)

	e1 := isa.MustGetEffect(isa.CondJump)
	as.Equal(isa.Labels, e1.Operand)

	defer func() {
		rec := recover()
		err := fmt.Errorf(isa.ErrEffectNotDeclared, isa.Opcode(5000))
		as.Equal(err, rec)
	}()

	isa.MustGetEffect(isa.Opcode(5000))
}
