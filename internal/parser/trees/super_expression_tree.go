// Port of com/google/javascript/jscomp/parsing/parser/trees/SuperExpressionTree.java
// @ closure-compiler v20260415

package trees

// SuperExpressionTree represents a super expression.
type SuperExpressionTree struct {
	ParseTree
}

func NewSuperExpressionTree(location SourceRange) *SuperExpressionTree {
	return &SuperExpressionTree{
		ParseTree: NewParseTree(ParseTreeTypeSuperExpression, location),
	}
}
