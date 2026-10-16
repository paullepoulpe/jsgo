// Port of com/google/javascript/jscomp/parsing/parser/trees/LiteralExpressionTree.java @ closure-compiler v20260415

package parser

// LiteralExpressionTree is the parse tree for ParseTreeType LITERAL_EXPRESSION.
type LiteralExpressionTree struct {
	ParseTreeBase

	LiteralToken Token
}

// NewLiteralExpressionTree constructs a LiteralExpressionTree.
//
// Java: LiteralExpressionTree(SourceRange, Token)
func NewLiteralExpressionTree(location *SourceRange, literalToken Token) *LiteralExpressionTree {
	t := &LiteralExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_LITERAL_EXPRESSION, location)
	t.LiteralToken = literalToken
	return t
}
