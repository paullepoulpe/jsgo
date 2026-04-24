// Port of com/google/javascript/jscomp/parsing/parser/trees/NewTargetExpressionTree.java
// @ closure-compiler v20260415

package trees

// NewTargetExpressionTree represents the new.target expression.
type NewTargetExpressionTree struct {
	ParseTree
}

// NewNewTargetExpressionTree constructs a NewTargetExpressionTree.
func NewNewTargetExpressionTree(location SourceRange) *NewTargetExpressionTree {
	return &NewTargetExpressionTree{
		ParseTree: ParseTree{Type: ParseTreeTypeNewTargetExpression, Location: location},
	}
}
