// Port of com/google/javascript/jscomp/parsing/parser/trees/LiteralExpressionTree.java
// @ closure-compiler v20260415

package trees

// LiteralExpressionTree represents a literal expression (string, number, regexp, etc.).
type LiteralExpressionTree struct {
	ParseTree
	LiteralToken *Token
}

func NewLiteralExpressionTree(location SourceRange, literalToken *Token) *LiteralExpressionTree {
	return &LiteralExpressionTree{
		ParseTree:    NewParseTree(ParseTreeTypeLiteralExpression, location),
		LiteralToken: literalToken,
	}
}
