package optimize

import (
	"gitlab.com/kode4food/ale/compiler/ir/visitor"
	"gitlab.com/kode4food/ale/runtime/isa"
)

type returnRoller struct{}

func rollReturns(root visitor.Node) visitor.Node {
	r := &returnRoller{}
	visitor.DepthFirst(root, r)
	return root
}

func (*returnRoller) EnterRoot(visitor.Node)            {}
func (*returnRoller) ExitRoot(visitor.Node)             {}
func (*returnRoller) EnterBranches(visitor.Branches)    {}
func (*returnRoller) Instructions(visitor.Instructions) {}

func (*returnRoller) ExitBranches(b visitor.Branches) {
	if i, ok := b.Epilogue().(visitor.Instructions); ok {
		code := i.Code()
		if len(code) == 1 && code[0].Opcode == isa.Return {
			if addReturnToBranches(b) {
				i.Set(isa.Instructions{})
			}
		}
	}
}

func addReturnToBranches(b visitor.Branches) bool {
	if ti, ok := b.ThenBranch().(visitor.Instructions); ok {
		if ei, ok := b.ElseBranch().(visitor.Instructions); ok {
			addReturnToInstructions(ti)
			addReturnToInstructions(ei)
			return true
		}
	}
	return false
}

func addReturnToInstructions(i visitor.Instructions) {
	code := i.Code()
	code = append(code, isa.New(isa.Return))
	i.Set(code)
}
