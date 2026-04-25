// Port of com/google/javascript/jscomp/parsing/parser/trees/NewExpressionTree.java
// @ closure-compiler v20260415

package trees

// NewExpressionTree represents a new expression, e.g. new Foo(a, b).
type NewExpressionTree struct {
	ParseTree
	Operand          *ParseTree
	Arguments        *ArgumentListTree
	HasTrailingComma bool
}

// NewNewExpressionTree constructs a NewExpressionTree.
// Go-only: language-adapter
func NewNewExpressionTree(location SourceRange, operand *ParseTree, arguments *ArgumentListTree, hasTrailingComma bool) *NewExpressionTree {
	return &NewExpressionTree{
		ParseTree:        ParseTree{Type: ParseTreeTypeNewExpression, Location: location},
		Operand:          operand,
		Arguments:        arguments,
		HasTrailingComma: hasTrailingComma,
	}
}
