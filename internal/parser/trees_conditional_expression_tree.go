// Port of com/google/javascript/jscomp/parsing/parser/trees/ConditionalExpressionTree.java @ closure-compiler v20260415

package parser

// ConditionalExpressionTree is the parse tree for ParseTreeType CONDITIONAL_EXPRESSION.
type ConditionalExpressionTree struct {
	ParseTreeBase

	Condition ParseTree
	Left      ParseTree
	Right     ParseTree
}

// NewConditionalExpressionTree constructs a ConditionalExpressionTree.
//
// Java: ConditionalExpressionTree(SourceRange, ParseTree, ParseTree, ParseTree)
func NewConditionalExpressionTree(location *SourceRange, condition ParseTree, left ParseTree, right ParseTree) *ConditionalExpressionTree {
	t := &ConditionalExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_CONDITIONAL_EXPRESSION, location)
	t.Condition = condition
	t.Left = left
	t.Right = right
	return t
}
