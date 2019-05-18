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
			i.Set(isa.Instructions{})
			addReturnToNode(b)
		}
	}
}

func addReturnToNode(n visitor.Node) {
	switch typed := n.(type) {
	case visitor.Branches:
		addReturnToBranches(typed)
	case visitor.Instructions:
		addReturnToInstructions(typed)
	}
}

func addReturnToBranches(b visitor.Branches) {
	if addReturnToEpilogue(b.Epilogue()) {
		return
	}
	addReturnToNode(b.ThenBranch())
	addReturnToNode(b.ElseBranch())
}

func addReturnToEpilogue(n visitor.Node) bool {
	if i, ok := n.(visitor.Instructions); ok {
		if len(i.Code()) > 0 {
			addReturnToInstructions(i)
			return true
		}
		return false
	}
	addReturnToBranches(n.(visitor.Branches))
	return true
}

func addReturnToInstructions(i visitor.Instructions) {
	code := i.Code()
	code = append(code, isa.New(isa.Return))
	i.Set(code)
}
