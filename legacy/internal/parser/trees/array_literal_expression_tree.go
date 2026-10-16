// Port of com/google/javascript/jscomp/parsing/parser/trees/ArrayLiteralExpressionTree.java
// @ closure-compiler v20260415

package trees

// ArrayLiteralExpressionTree represents an array literal expression,
// e.g. [a, b, c].
type ArrayLiteralExpressionTree struct {
	ParseTree
	Elements         []*ParseTree
	HasTrailingComma bool
}

// NewArrayLiteralExpressionTree constructs an ArrayLiteralExpressionTree.
// Go-only: language-adapter
func NewArrayLiteralExpressionTree(location SourceRange, elements []*ParseTree, hasTrailingComma bool) *ArrayLiteralExpressionTree {
	return &ArrayLiteralExpressionTree{
		ParseTree:        ParseTree{Type: ParseTreeTypeArrayLiteralExpression, Location: location},
		Elements:         elements,
		HasTrailingComma: hasTrailingComma,
	}
}
