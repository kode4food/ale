package visitor

// Type is an interface that is called back upon visiting
type Type interface {
	EnterRoot(Node)
	ExitRoot(Node)
	EnterBranches(Branches)
	ExitBranches(Branches)
	Instructions(Instructions)
}

// DepthFirst performs a depth-first visitation
func DepthFirst(root Node, visitor Type) {
	visitor.EnterRoot(root)
	depthFirst(root, visitor)
	visitor.ExitRoot(root)
}

func depthFirst(node Node, visitor Type) {
	switch typed := node.(type) {
	case Instructions:
		visitor.Instructions(typed)
	case Branches:
		visitor.EnterBranches(typed)
		depthFirst(typed.Epilogue(), visitor)
		depthFirst(typed.ThenBranch(), visitor)
		depthFirst(typed.ElseBranch(), visitor)
		depthFirst(typed.Prologue(), visitor)
		visitor.ExitBranches(typed)
	default:
		panic("unexpected node type")
	}
}
