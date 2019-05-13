package optimize

import "gitlab.com/kode4food/ale/compiler/internal/visitor"

type tailCallsVisitor struct{}

func tailCalls(root visitor.Node) visitor.Node {
	visitor.DepthFirst(root, &tailCallsVisitor{})
	return root
}

func (t *tailCallsVisitor) EnterRoot(_ visitor.Node)            {}
func (t *tailCallsVisitor) ExitRoot(_ visitor.Node)             {}
func (t *tailCallsVisitor) EnterBranches(_ visitor.Branches)    {}
func (t *tailCallsVisitor) ExitBranches(_ visitor.Branches)     {}
func (t *tailCallsVisitor) Instructions(_ visitor.Instructions) {}
