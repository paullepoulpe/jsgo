// Port of com/google/javascript/jscomp/parsing/parser/trees/ThisExpressionTree.java
// @ closure-compiler v20260415

package trees

// ThisExpressionTree represents the this keyword expression.
type ThisExpressionTree struct {
	ParseTree
}

// NewThisExpressionTree constructs a ThisExpressionTree.
func NewThisExpressionTree(location SourceRange) *ThisExpressionTree {
	return &ThisExpressionTree{
		ParseTree: ParseTree{Type: ParseTreeTypeThisExpression, Location: location},
	}
}
