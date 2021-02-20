package optimize

import (
	"github.com/kode4food/ale/compiler/ir/visitor"
	"github.com/kode4food/ale/runtime/isa"
)

type returnSplitter struct{}

func splitReturns(root visitor.Node) visitor.Node {
	r := new(returnSplitter)
	visitor.DepthFirst(root, r)
	return root
}

func (*returnSplitter) EnterRoot(visitor.Node)            {}
func (*returnSplitter) ExitRoot(visitor.Node)             {}
func (*returnSplitter) EnterBranches(visitor.Branches)    {}
func (*returnSplitter) Instructions(visitor.Instructions) {}

func (*returnSplitter) ExitBranches(b visitor.Branches) {
	if i, ok := b.Epilogue().(visitor.Instructions); ok {
		code := i.Code()
		if len(code) == 1 && code[0].Opcode == isa.Return {
			i.Set(isa.Instructions{})
			addReturnToNode(b)
		}
	}
}

func addReturnToNode(n visitor.Node) {
	switch n := n.(type) {
	case visitor.Branches:
		addReturnToBranches(n)
	case visitor.Instructions:
		addReturnToInstructions(n)
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
