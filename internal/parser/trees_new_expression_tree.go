// Port of com/google/javascript/jscomp/parsing/parser/trees/NewExpressionTree.java @ closure-compiler v20260415

package parser

// NewExpressionTree is the parse tree for ParseTreeType NEW_EXPRESSION.
type NewExpressionTree struct {
	ParseTreeBase

	Operand          ParseTree
	Arguments        *ArgumentListTree
	HasTrailingComma bool
}

// NewNewExpressionTree constructs a NewExpressionTree.
//
// Java: NewExpressionTree(SourceRange, ParseTree, ArgumentListTree, boolean)
func NewNewExpressionTree(location *SourceRange, operand ParseTree, arguments *ArgumentListTree, hasTrailingComma bool) *NewExpressionTree {
	t := &NewExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_NEW_EXPRESSION, location)
	t.Operand = operand
	t.Arguments = arguments
	t.HasTrailingComma = hasTrailingComma
	return t
}
