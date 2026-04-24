// Port of com/google/javascript/jscomp/parsing/parser/trees/MissingPrimaryExpressionTree.java
// @ closure-compiler v20260415

package trees

// MissingPrimaryExpressionTree is inserted when a primary expression is expected
// but not found (error recovery).
type MissingPrimaryExpressionTree struct {
	ParseTree
}

// NewMissingPrimaryExpressionTree constructs a MissingPrimaryExpressionTree.
func NewMissingPrimaryExpressionTree(location SourceRange) *MissingPrimaryExpressionTree {
	return &MissingPrimaryExpressionTree{
		ParseTree: ParseTree{Type: ParseTreeTypeMissingPrimaryExpression, Location: location},
	}
}
