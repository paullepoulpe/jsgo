// Port of com/google/javascript/jscomp/parsing/parser/trees/ConditionalExpressionTree.java
// @ closure-compiler v20260415

package trees

// ConditionalExpressionTree represents a ternary expression, e.g. a ? b : c.
type ConditionalExpressionTree struct {
	ParseTree
	Condition *ParseTree
	Left      *ParseTree
	Right     *ParseTree
}

// NewConditionalExpressionTree constructs a ConditionalExpressionTree.
// Go-only: language-adapter
func NewConditionalExpressionTree(location SourceRange, condition *ParseTree, left *ParseTree, right *ParseTree) *ConditionalExpressionTree {
	return &ConditionalExpressionTree{
		ParseTree: ParseTree{Type: ParseTreeTypeConditionalExpression, Location: location},
		Condition: condition,
		Left:      left,
		Right:     right,
	}
}
