// Port of com/google/javascript/jscomp/parsing/parser/trees/UpdateExpressionTree.java
// @ closure-compiler v20260415

package trees

// OperatorPosition records whether a ++ / -- operator is prefix or postfix.
type OperatorPosition int

const (
	OperatorPositionPrefix  OperatorPosition = iota
	OperatorPositionPostfix
)

// UpdateExpressionTree represents an update expression (++ / --).
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

func NewUpdateExpressionTree(
	location SourceRange,
	operator *Token,
	operatorPosition OperatorPosition,
	operand *ParseTree,
) *UpdateExpressionTree {
	return &UpdateExpressionTree{
		ParseTree:        NewParseTree(ParseTreeTypeUpdateExpression, location),
		Operator:         operator,
		OperatorPosition: operatorPosition,
		Operand:          operand,
	}
}

func UpdateExpressionPrefix(location SourceRange, operator *Token, operand *ParseTree) *UpdateExpressionTree {
	return NewUpdateExpressionTree(location, operator, OperatorPositionPrefix, operand)
}

func UpdateExpressionPostfix(location SourceRange, operator *Token, operand *ParseTree) *UpdateExpressionTree {
	return NewUpdateExpressionTree(location, operator, OperatorPositionPostfix, operand)
}
