package optimize

import "gitlab.com/kode4food/ale/compiler/internal/visitor"

type rollUpReturnsVisitor struct{}

func rollUpReturns(root visitor.Node) visitor.Node {
	visitor.DepthFirst(root, &rollUpReturnsVisitor{})
	return root
}

func (r *rollUpReturnsVisitor) EnterRoot(_ visitor.Node)            {}
func (r *rollUpReturnsVisitor) ExitRoot(_ visitor.Node)             {}
func (r *rollUpReturnsVisitor) EnterBranches(_ visitor.Branches)    {}
func (r *rollUpReturnsVisitor) ExitBranches(_ visitor.Branches)     {}
func (r *rollUpReturnsVisitor) Instructions(_ visitor.Instructions) {}
