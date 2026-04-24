// Port of com/google/javascript/jscomp/parsing/parser/trees/NewTargetExpressionTree.java
// @ closure-compiler v20260415

package trees

// NewTargetExpressionTree represents the new.target meta-property.
type NewTargetExpressionTree struct {
	ParseTree
}

func NewNewTargetExpressionTree(location SourceRange) *NewTargetExpressionTree {
	return &NewTargetExpressionTree{
		ParseTree: NewParseTree(ParseTreeTypeNewTargetExpression, location),
	}
}
