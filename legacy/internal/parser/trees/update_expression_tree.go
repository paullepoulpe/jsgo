// Port of com/google/javascript/jscomp/parsing/parser/trees/UpdateExpressionTree.java
// @ closure-compiler v20260415

package trees

// OperatorPosition indicates whether an update operator is prefix or postfix.
type OperatorPosition int

const (
	// OperatorPositionPrefix indicates the operator precedes the operand, e.g. ++x.
	OperatorPositionPrefix OperatorPosition = iota
	// OperatorPositionPostfix indicates the operator follows the operand, e.g. x++.
	OperatorPositionPostfix
)

// UpdateExpressionTree represents UpdateExpression productions from the spec:
//
//	UpdateExpression :=
//	    { ++ | -- } UnaryExpression
//	    LeftHandSideExpression [no LineTerminator here] { ++ | -- }
type UpdateExpressionTree struct {
	ParseTree
	Operator         *Token
	OperatorPosition OperatorPosition
	Operand          *ParseTree
}

// NewUpdateExpressionTree constructs an UpdateExpressionTree.
// Go-only: language-adapter
func NewUpdateExpressionTree(location SourceRange, operator *Token, operatorPosition OperatorPosition, operand *ParseTree) *UpdateExpressionTree {
	return &UpdateExpressionTree{
		ParseTree:        ParseTree{Type: ParseTreeTypeUpdateExpression, Location: location},
		Operator:         operator,
		OperatorPosition: operatorPosition,
		Operand:          operand,
	}
}

// NewPrefixUpdateExpressionTree constructs a prefix update expression, e.g. ++x.
// Go-only: language-adapter
func NewPrefixUpdateExpressionTree(location SourceRange, operator *Token, operand *ParseTree) *UpdateExpressionTree {
	return NewUpdateExpressionTree(location, operator, OperatorPositionPrefix, operand)
}

// NewPostfixUpdateExpressionTree constructs a postfix update expression, e.g. x++.
// Go-only: language-adapter
func NewPostfixUpdateExpressionTree(location SourceRange, operator *Token, operand *ParseTree) *UpdateExpressionTree {
	return NewUpdateExpressionTree(location, operator, OperatorPositionPostfix, operand)
}
