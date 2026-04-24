// Port of com/google/javascript/jscomp/parsing/parser/trees/SuperExpressionTree.java
// @ closure-compiler v20260415

package trees

// SuperExpressionTree represents the super keyword expression.
type SuperExpressionTree struct {
	ParseTree
}

// NewSuperExpressionTree constructs a SuperExpressionTree.
// Go-only: language-adapter
func NewSuperExpressionTree(location SourceRange) *SuperExpressionTree {
	return &SuperExpressionTree{
		ParseTree: ParseTree{Type: ParseTreeTypeSuperExpression, Location: location},
	}
}
