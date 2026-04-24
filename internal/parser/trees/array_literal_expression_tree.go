// Port of com/google/javascript/jscomp/parsing/parser/trees/ArrayLiteralExpressionTree.java
// @ closure-compiler v20260415

package trees

// ArrayLiteralExpressionTree represents an array literal expression.
type ArrayLiteralExpressionTree struct {
	ParseTree
	Elements         []*ParseTree
	HasTrailingComma bool
}

func NewArrayLiteralExpressionTree(
	location SourceRange,
	elements []*ParseTree,
	hasTrailingComma bool,
) *ArrayLiteralExpressionTree {
	return &ArrayLiteralExpressionTree{
		ParseTree:        NewParseTree(ParseTreeTypeArrayLiteralExpression, location),
		Elements:         elements,
		HasTrailingComma: hasTrailingComma,
	}
}
