// Port of com/google/javascript/jscomp/parsing/parser/trees/IdentifierExpressionTree.java
// @ closure-compiler v20260415

package trees

// IdentifierExpressionTree represents an identifier reference expression, e.g. foo.
type IdentifierExpressionTree struct {
	ParseTree
	IdentifierToken *IdentifierToken
}

// NewIdentifierExpressionTree constructs an IdentifierExpressionTree.
func NewIdentifierExpressionTree(location SourceRange, identifierToken *IdentifierToken) *IdentifierExpressionTree {
	return &IdentifierExpressionTree{
		ParseTree:       ParseTree{Type: ParseTreeTypeIdentifierExpression, Location: location},
		IdentifierToken: identifierToken,
	}
}
