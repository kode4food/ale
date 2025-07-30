package optimize

import (
	"github.com/kode4food/ale/internal/compiler/ir/visitor"
	"github.com/kode4food/ale/internal/runtime/isa"
)

func paramBranchFor(c isa.Instructions, argc isa.Operand) isa.Instructions {
	b := &visitor.BranchScanner{
		Then:     visitor.All,
		Epilogue: visitor.All,
	}
	b.Else = b.Scan

	if b, ok := b.Scan(c).(visitor.Branches); ok {
		if bc := getParamBranch(b, argc); bc != nil {
			return bc
		}
	}
	return c
}

func getParamBranch(b visitor.Branches, argc isa.Operand) isa.Instructions {
	if len(b.Epilogue().Code()) != 0 {
		// compiled procedures don't include epilogues in the arity branching
		// logic, so if any node along the path has an epilogue, then we can't
		// inline the 'then' branch
		return nil
	}
	oc, op, ok := isParamCase(b)
	if !ok {
		return nil
	}
	if oc == isa.NumEq && argc == op || oc == isa.NumGte && argc >= op {
		return b.ThenBranch().Code()
	}
	if eb, ok := b.ElseBranch().(visitor.Branches); ok {
		return getParamBranch(eb, argc)
	}
	return nil
}

func isParamCase(b visitor.Branches) (isa.Opcode, isa.Operand, bool) {
	p := b.Prologue().Code()
	if len(p) != 4 {
		return isa.NoOp, 0, false
	}
	if p[0].Opcode() != isa.ArgsLen || p[3].Opcode() != isa.CondJump {
		return isa.NoOp, 0, false
	}
	if p[1].Opcode() != isa.PosInt {
		return isa.NoOp, 0, false
	}
	oc := p[2].Opcode()
	op := p[1].Operand()
	return oc, op, oc == isa.NumEq || oc == isa.NumGte
}
