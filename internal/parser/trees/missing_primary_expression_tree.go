// Port of com/google/javascript/jscomp/parsing/parser/trees/MissingPrimaryExpressionTree.java
// @ closure-compiler v20260415

package trees

// MissingPrimaryExpressionTree is a sentinel node for parse errors.
type MissingPrimaryExpressionTree struct {
	ParseTree
}

func NewMissingPrimaryExpressionTree(location SourceRange) *MissingPrimaryExpressionTree {
	return &MissingPrimaryExpressionTree{
		ParseTree: NewParseTree(ParseTreeTypeMissingPrimaryExpression, location),
	}
}
