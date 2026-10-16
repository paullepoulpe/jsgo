// Port of com/google/javascript/jscomp/parsing/parser/trees/CommaExpressionTree.java @ closure-compiler v20260415

package parser

// CommaExpressionTree is the parse tree for ParseTreeType COMMA_EXPRESSION.
type CommaExpressionTree struct {
	ParseTreeBase

	Expressions []ParseTree
}

// NewCommaExpressionTree constructs a CommaExpressionTree.
//
// Java: CommaExpressionTree(SourceRange, ImmutableList)
func NewCommaExpressionTree(location *SourceRange, expressions []ParseTree) *CommaExpressionTree {
	t := &CommaExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_COMMA_EXPRESSION, location)
	t.Expressions = expressions
	return t
}
