// Port of com/google/javascript/jscomp/parsing/parser/trees/CallExpressionTree.java
// @ closure-compiler v20260415

package trees

// CallExpressionTree represents a function call expression, e.g. f(a, b).
type CallExpressionTree struct {
	ParseTree
	Operand   *ParseTree
	Arguments *ArgumentListTree
}

// NewCallExpressionTree constructs a CallExpressionTree.
// Go-only: language-adapter
func NewCallExpressionTree(location SourceRange, operand *ParseTree, arguments *ArgumentListTree) *CallExpressionTree {
	return &CallExpressionTree{
		ParseTree: ParseTree{Type: ParseTreeTypeCallExpression, Location: location},
		Operand:   operand,
		Arguments: arguments,
	}
}
