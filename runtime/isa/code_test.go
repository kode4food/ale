package isa_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/runtime/isa"
)

func TestCoders(t *testing.T) {
	as := assert.New(t)

	i1 := isa.Index(99)
	c1 := isa.Count(5)
	o1 := isa.Offset(37)

	in1 := []isa.Coder{i1, c1, o1}
	out1 := make([]isa.Word, len(in1))
	for i, c := range in1 {
		out1[i] = c.Word()
	}
	as.Number(99, data.Integer(out1[0]))
	as.Number(5, data.Integer(out1[1]))
	as.Number(37, data.Integer(out1[2]))
}

func TestInstructions(t *testing.T) {
	as := assert.New(t)

	i1 := isa.New(isa.CondJump, isa.Offset(27).Word())
	as.String("CondJump(27)", i1)

	defer func() {
		if rec := recover(); rec != nil {
			err := fmt.Sprintf(isa.ErrBadInstructionArgs, "CondJump")
			as.String(err, rec)
		} else {
			as.Fail("proper error not raised")
		}
	}()

	isa.New(isa.CondJump, isa.Word(12), isa.Word(32))
}

func TestInstructionString(t *testing.T) {
	as := assert.New(t)
	inst := isa.New(isa.Const, isa.Offset(0).Word())
	as.String(`Const(0)`, inst.String())
}
