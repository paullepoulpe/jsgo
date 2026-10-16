// Port of com/google/javascript/jscomp/parsing/parser/trees/CallExpressionTree.java @ closure-compiler v20260415

package parser

// CallExpressionTree is the parse tree for ParseTreeType CALL_EXPRESSION.
type CallExpressionTree struct {
	ParseTreeBase

	Operand   ParseTree
	Arguments *ArgumentListTree
}

// NewCallExpressionTree constructs a CallExpressionTree.
//
// Java: CallExpressionTree(SourceRange, ParseTree, ArgumentListTree)
func NewCallExpressionTree(location *SourceRange, operand ParseTree, arguments *ArgumentListTree) *CallExpressionTree {
	t := &CallExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_CALL_EXPRESSION, location)
	t.Operand = operand
	t.Arguments = arguments
	return t
}
