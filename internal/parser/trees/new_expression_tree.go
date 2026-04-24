// Port of com/google/javascript/jscomp/parsing/parser/trees/NewExpressionTree.java
// @ closure-compiler v20260415

package trees

// NewExpressionTree represents a new expression.
type NewExpressionTree struct {
	ParseTree
	Operand          *ParseTree
	Arguments        *ArgumentListTree
	HasTrailingComma bool
}

func NewNewExpressionTree(
	location SourceRange,
	operand *ParseTree,
	arguments *ArgumentListTree,
	hasTrailingComma bool,
) *NewExpressionTree {
	return &NewExpressionTree{
		ParseTree:        NewParseTree(ParseTreeTypeNewExpression, location),
		Operand:          operand,
		Arguments:        arguments,
		HasTrailingComma: hasTrailingComma,
	}
}
