// Port of com/google/javascript/jscomp/parsing/parser/trees/LiteralExpressionTree.java
// @ closure-compiler v20260415

package trees

// LiteralExpressionTree represents a literal value expression,
// e.g. 42, "hello", true, null.
type LiteralExpressionTree struct {
	ParseTree
	LiteralToken *Token
}

// NewLiteralExpressionTree constructs a LiteralExpressionTree.
func NewLiteralExpressionTree(location SourceRange, literalToken *Token) *LiteralExpressionTree {
	return &LiteralExpressionTree{
		ParseTree:    ParseTree{Type: ParseTreeTypeLiteralExpression, Location: location},
		LiteralToken: literalToken,
	}
}
