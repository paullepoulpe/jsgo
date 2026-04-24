// Port of com/google/javascript/jscomp/parsing/parser/trees/YieldExpressionTree.java
// @ closure-compiler v20260415

package trees

// YieldExpressionTree represents a yield expression in a generator function,
// e.g. yield expr or yield* expr.
type YieldExpressionTree struct {
	ParseTree
	Expression *ParseTree
	IsYieldAll bool
}

// NewYieldExpressionTree constructs a YieldExpressionTree.
func NewYieldExpressionTree(location SourceRange, isYieldAll bool, expression *ParseTree) *YieldExpressionTree {
	return &YieldExpressionTree{
		ParseTree:  ParseTree{Type: ParseTreeTypeYieldExpression, Location: location},
		IsYieldAll: isYieldAll,
		Expression: expression,
	}
}
