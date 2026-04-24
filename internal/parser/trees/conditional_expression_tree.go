// Port of com/google/javascript/jscomp/parsing/parser/trees/ConditionalExpressionTree.java
// @ closure-compiler v20260415

package trees

// ConditionalExpressionTree represents a ternary conditional expression (condition ? left : right).
type ConditionalExpressionTree struct {
	ParseTree
	Condition *ParseTree
	Left      *ParseTree
	Right     *ParseTree
}

func NewConditionalExpressionTree(
	location SourceRange,
	condition *ParseTree,
	left *ParseTree,
	right *ParseTree,
) *ConditionalExpressionTree {
	return &ConditionalExpressionTree{
		ParseTree: NewParseTree(ParseTreeTypeConditionalExpression, location),
		Condition: condition,
		Left:      left,
		Right:     right,
	}
}
