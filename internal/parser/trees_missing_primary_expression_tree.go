// Port of com/google/javascript/jscomp/parsing/parser/trees/MissingPrimaryExpressionTree.java @ closure-compiler v20260415

package parser

// MissingPrimaryExpressionTree is the parse tree for ParseTreeType MISSING_PRIMARY_EXPRESSION.
type MissingPrimaryExpressionTree struct {
	ParseTreeBase
}

// NewMissingPrimaryExpressionTree constructs a MissingPrimaryExpressionTree.
//
// Java: MissingPrimaryExpressionTree(SourceRange)
func NewMissingPrimaryExpressionTree(location *SourceRange) *MissingPrimaryExpressionTree {
	t := &MissingPrimaryExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_MISSING_PRIMARY_EXPRESSION, location)
	return t
}
