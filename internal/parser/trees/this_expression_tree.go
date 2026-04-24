// Port of com/google/javascript/jscomp/parsing/parser/trees/ThisExpressionTree.java
// @ closure-compiler v20260415

package trees

// ThisExpressionTree represents a this expression.
type ThisExpressionTree struct {
	ParseTree
}

func NewThisExpressionTree(location SourceRange) *ThisExpressionTree {
	return &ThisExpressionTree{
		ParseTree: NewParseTree(ParseTreeTypeThisExpression, location),
	}
}
