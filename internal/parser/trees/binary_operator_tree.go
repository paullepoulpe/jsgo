// Port of com/google/javascript/jscomp/parsing/parser/trees/BinaryOperatorTree.java
// @ closure-compiler v20260415

package trees

// BinaryOperatorTree represents a binary operator expression, e.g. a + b.
type BinaryOperatorTree struct {
	ParseTree
	Left     *ParseTree
	Operator *Token
	Right    *ParseTree
}

// NewBinaryOperatorTree constructs a BinaryOperatorTree.
func NewBinaryOperatorTree(location SourceRange, left *ParseTree, operator *Token, right *ParseTree) *BinaryOperatorTree {
	return &BinaryOperatorTree{
		ParseTree: ParseTree{Type: ParseTreeTypeBinaryOperator, Location: location},
		Left:      left,
		Operator:  operator,
		Right:     right,
	}
}
