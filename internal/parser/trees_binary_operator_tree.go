// Port of com/google/javascript/jscomp/parsing/parser/trees/BinaryOperatorTree.java @ closure-compiler v20260415

package parser

// BinaryOperatorTree is the parse tree for ParseTreeType BINARY_OPERATOR.
type BinaryOperatorTree struct {
	ParseTreeBase

	Left     ParseTree
	Operator Token
	Right    ParseTree
}

// NewBinaryOperatorTree constructs a BinaryOperatorTree.
//
// Java: BinaryOperatorTree(SourceRange, ParseTree, Token, ParseTree)
func NewBinaryOperatorTree(location *SourceRange, left ParseTree, operator Token, right ParseTree) *BinaryOperatorTree {
	t := &BinaryOperatorTree{}
	t.InitParseTreeBase(t, ParseTreeType_BINARY_OPERATOR, location)
	t.Left = left
	t.Operator = operator
	t.Right = right
	return t
}
