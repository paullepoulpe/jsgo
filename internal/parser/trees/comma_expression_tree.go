// Port of com/google/javascript/jscomp/parsing/parser/trees/CommaExpressionTree.java
// @ closure-compiler v20260415

package trees

// CommaExpressionTree represents a comma-separated sequence of expressions.
type CommaExpressionTree struct {
	ParseTree
	Expressions []*ParseTree
}

func NewCommaExpressionTree(location SourceRange, expressions []*ParseTree) *CommaExpressionTree {
	return &CommaExpressionTree{
		ParseTree:   NewParseTree(ParseTreeTypeCommaExpression, location),
		Expressions: expressions,
	}
}
