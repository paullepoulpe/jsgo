// Port of com/google/javascript/jscomp/parsing/parser/trees/CallExpressionTree.java
// @ closure-compiler v20260415

package trees

// CallExpressionTree represents a function call expression.
type CallExpressionTree struct {
	ParseTree
	Operand   *ParseTree
	Arguments *ArgumentListTree
}

func NewCallExpressionTree(
	location SourceRange,
	operand *ParseTree,
	arguments *ArgumentListTree,
) *CallExpressionTree {
	return &CallExpressionTree{
		ParseTree: NewParseTree(ParseTreeTypeCallExpression, location),
		Operand:   operand,
		Arguments: arguments,
	}
}
