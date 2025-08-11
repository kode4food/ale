package visitor

import "github.com/kode4food/ale/internal/debug"

// Visitor is an interface called back upon visiting
type Visitor interface {
	// EnterRoot is called when entering the root node
	EnterRoot(Node)

	// ExitRoot is called when exiting the root node
	ExitRoot(Node)

	// EnterBranches is called when entering a branches node
	EnterBranches(Branches)

	// ExitBranches is called when exiting a branches node
	ExitBranches(Branches)

	// Instructions is called when visiting an instructions node
	Instructions(Instructions)
}

// Visit performs a branched visitation
func Visit(root Node, visitor Visitor) {
	visitor.EnterRoot(root)
	visit(root, visitor)
	visitor.ExitRoot(root)
}

func visit(node Node, visitor Visitor) {
	switch node := node.(type) {
	case Instructions:
		visitor.Instructions(node)
	case Branches:
		visitor.EnterBranches(node)
		visit(node.Prologue(), visitor)
		visit(node.ElseBranch(), visitor)
		visit(node.ThenBranch(), visitor)
		visit(node.Epilogue(), visitor)
		visitor.ExitBranches(node)
	default:
		panic(debug.ProgrammerError("unexpected node type"))
	}
}
