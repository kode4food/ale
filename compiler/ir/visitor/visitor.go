package visitor

// Visitor is an interface that is called back upon visiting
type Visitor interface {
	EnterRoot(Node)
	ExitRoot(Node)
	EnterBranches(Branches)
	ExitBranches(Branches)
	Instructions(Instructions)
}

// Error messages
const (
	errUnexpectedNodeType = "unexpected node type"
)

// DepthFirst performs a depth-first visitation
func DepthFirst(root Node, visitor Visitor) {
	visitor.EnterRoot(root)
	depthFirst(root, visitor)
	visitor.ExitRoot(root)
}

func depthFirst(node Node, visitor Visitor) {
	switch node := node.(type) {
	case Instructions:
		visitor.Instructions(node)
	case Branches:
		visitor.EnterBranches(node)
		depthFirst(node.Epilogue(), visitor)
		depthFirst(node.ThenBranch(), visitor)
		depthFirst(node.ElseBranch(), visitor)
		depthFirst(node.Prologue(), visitor)
		visitor.ExitBranches(node)
	default:
		// Programmer error
		panic(errUnexpectedNodeType)
	}
}
