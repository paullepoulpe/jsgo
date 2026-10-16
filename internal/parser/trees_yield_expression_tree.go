// Port of com/google/javascript/jscomp/parsing/parser/trees/YieldExpressionTree.java @ closure-compiler v20260415

package parser

// YieldExpressionTree is the parse tree for ParseTreeType YIELD_EXPRESSION.
type YieldExpressionTree struct {
	ParseTreeBase

	Expression ParseTree
	IsYieldAll bool
}

// NewYieldExpressionTree constructs a YieldExpressionTree.
//
// Java: YieldExpressionTree(SourceRange, boolean, ParseTree)
func NewYieldExpressionTree(location *SourceRange, isYieldAll bool, expression ParseTree) *YieldExpressionTree {
	t := &YieldExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_YIELD_EXPRESSION, location)
	t.IsYieldAll = isYieldAll
	t.Expression = expression
	return t
}
