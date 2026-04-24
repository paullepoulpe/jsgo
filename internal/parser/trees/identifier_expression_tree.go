// Port of com/google/javascript/jscomp/parsing/parser/trees/IdentifierExpressionTree.java
// @ closure-compiler v20260415

package trees

// IdentifierExpressionTree represents an identifier expression.
type IdentifierExpressionTree struct {
	ParseTree
	IdentifierToken *IdentifierToken
}

func NewIdentifierExpressionTree(
	location SourceRange,
	identifierToken *IdentifierToken,
) *IdentifierExpressionTree {
	return &IdentifierExpressionTree{
		ParseTree:       NewParseTree(ParseTreeTypeIdentifierExpression, location),
		IdentifierToken: identifierToken,
	}
}
