// Port of com/google/javascript/jscomp/parsing/parser/trees/YieldExpressionTree.java
// @ closure-compiler v20260415

package trees

// YieldExpressionTree represents a yield expression.
// IsYieldAll is true for yield* (delegate) expressions.
type YieldExpressionTree struct {
	ParseTree
	IsYieldAll bool
	Expression *ParseTree
}

func NewYieldExpressionTree(
	location SourceRange,
	isYieldAll bool,
	expression *ParseTree,
) *YieldExpressionTree {
	return &YieldExpressionTree{
		ParseTree:  NewParseTree(ParseTreeTypeYieldExpression, location),
		IsYieldAll: isYieldAll,
		Expression: expression,
	}
}
