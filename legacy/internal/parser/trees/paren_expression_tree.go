// Port of com/google/javascript/jscomp/parsing/parser/trees/ParenExpressionTree.java
// @ closure-compiler v20260415

package trees

// ParenExpressionTree represents a parenthesized expression, e.g. (expr).
type ParenExpressionTree struct {
	ParseTree
	Expression *ParseTree
}

// NewParenExpressionTree constructs a ParenExpressionTree.
// Go-only: language-adapter
func NewParenExpressionTree(location SourceRange, expression *ParseTree) *ParenExpressionTree {
	return &ParenExpressionTree{
		ParseTree:  ParseTree{Type: ParseTreeTypeParenExpression, Location: location},
		Expression: expression,
	}
}
