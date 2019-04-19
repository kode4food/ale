package isa_test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/assert"
	"gitlab.com/kode4food/ale/runtime/isa"
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
	as.Integer(99, api.Integer(out1[0]))
	as.Integer(5, api.Integer(out1[1]))
	as.Integer(37, api.Integer(out1[2]))
}

func TestInstructions(t *testing.T) {
	as := assert.New(t)

	i1 := isa.New(isa.CondJump, isa.Offset(27).Word())
	as.String("CondJump(27)", i1)

	defer func() {
		if rec := recover(); rec != nil {
			err := fmt.Sprintf(isa.BadInstructionArgs, "CondJump")
			as.String(err, rec)
		} else {
			as.Fail("proper error not raised")
		}
	}()

	isa.New(isa.CondJump, isa.Word(12), isa.Word(32))
}
