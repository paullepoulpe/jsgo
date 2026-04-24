// Port of com/google/javascript/jscomp/parsing/parser/trees/ParenExpressionTree.java
// @ closure-compiler v20260415

package trees

// ParenExpressionTree represents a parenthesised expression.
type ParenExpressionTree struct {
	ParseTree
	Expression *ParseTree
}

func NewParenExpressionTree(location SourceRange, expression *ParseTree) *ParenExpressionTree {
	return &ParenExpressionTree{
		ParseTree:  NewParseTree(ParseTreeTypeParenExpression, location),
		Expression: expression,
	}
}
