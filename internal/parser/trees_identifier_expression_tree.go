// Port of com/google/javascript/jscomp/parsing/parser/trees/IdentifierExpressionTree.java @ closure-compiler v20260415

package parser

// IdentifierExpressionTree is the parse tree for ParseTreeType IDENTIFIER_EXPRESSION.
type IdentifierExpressionTree struct {
	ParseTreeBase

	IdentifierToken *IdentifierToken
}

// NewIdentifierExpressionTree constructs a IdentifierExpressionTree.
//
// Java: IdentifierExpressionTree(SourceRange, IdentifierToken)
func NewIdentifierExpressionTree(location *SourceRange, identifierToken *IdentifierToken) *IdentifierExpressionTree {
	t := &IdentifierExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_IDENTIFIER_EXPRESSION, location)
	t.IdentifierToken = identifierToken
	return t
}
