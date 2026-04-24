// Port of com/google/javascript/jscomp/parsing/parser/trees/CommaExpressionTree.java
// @ closure-compiler v20260415

package trees

// CommaExpressionTree represents a comma (sequence) expression, e.g. a, b, c.
type CommaExpressionTree struct {
	ParseTree
	Expressions []*ParseTree
}

// NewCommaExpressionTree constructs a CommaExpressionTree.
func NewCommaExpressionTree(location SourceRange, expressions []*ParseTree) *CommaExpressionTree {
	return &CommaExpressionTree{
		ParseTree:   ParseTree{Type: ParseTreeTypeCommaExpression, Location: location},
		Expressions: expressions,
	}
}
