// Port of com/google/javascript/jscomp/parsing/parser/trees/BinaryOperatorTree.java
// @ closure-compiler v20260415

package trees

// BinaryOperatorTree represents a binary operator expression.
type BinaryOperatorTree struct {
	ParseTree
	Left     *ParseTree
	Operator *Token
	Right    *ParseTree
}

func NewBinaryOperatorTree(
	location SourceRange,
	left *ParseTree,
	operator *Token,
	right *ParseTree,
) *BinaryOperatorTree {
	return &BinaryOperatorTree{
		ParseTree: NewParseTree(ParseTreeTypeBinaryOperator, location),
		Left:      left,
		Operator:  operator,
		Right:     right,
	}
}
