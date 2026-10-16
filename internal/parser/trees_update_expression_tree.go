// Port of com/google/javascript/jscomp/parsing/parser/trees/UpdateExpressionTree.java @ closure-compiler v20260415

package parser

// Represents UpdateExpression productions from the spec.
//
// UpdateExpression :=
//
//	{ ++ | -- } UnaryExpression
//	LeftHandSideExpression [no LineTerminator here] { ++ | -- }
type UpdateExpressionTree struct {
	ParseTreeBase

	Operator         Token
	OperatorPosition UpdateExpressionTree_OperatorPosition
	Operand          ParseTree
}

// UpdateExpressionTree_OperatorPosition is UpdateExpressionTree.OperatorPosition:
// position of the operator relative to the operand.
type UpdateExpressionTree_OperatorPosition int32

const (
	UpdateExpressionTree_OperatorPosition_PREFIX UpdateExpressionTree_OperatorPosition = iota
	UpdateExpressionTree_OperatorPosition_POSTFIX
)

var updateExpressionTree_OperatorPosition_names = [...]string{"PREFIX", "POSTFIX"}

// UpdateExpressionTree_OperatorPosition_values is OperatorPosition.values().
func UpdateExpressionTree_OperatorPosition_values() []UpdateExpressionTree_OperatorPosition {
	return []UpdateExpressionTree_OperatorPosition{UpdateExpressionTree_OperatorPosition_PREFIX, UpdateExpressionTree_OperatorPosition_POSTFIX}
}

// Name is Enum.name().
func (p UpdateExpressionTree_OperatorPosition) Name() string {
	return updateExpressionTree_OperatorPosition_names[p]
}

// String is Enum.toString().
func (p UpdateExpressionTree_OperatorPosition) String() string { return p.Name() }

// NewUpdateExpressionTree constructs a UpdateExpressionTree.
//
// Java: UpdateExpressionTree(SourceRange, Token, OperatorPosition, ParseTree)
func NewUpdateExpressionTree(location *SourceRange, operator Token, operatorPosition UpdateExpressionTree_OperatorPosition, operand ParseTree) *UpdateExpressionTree {
	t := &UpdateExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_UPDATE_EXPRESSION, location)
	t.Operator = operator
	t.OperatorPosition = operatorPosition
	t.Operand = operand
	return t
}

// UpdateExpressionTree_Prefix creates a prefix update expression.
//
// Java: UpdateExpressionTree.prefix(SourceRange, Token, ParseTree)
func UpdateExpressionTree_Prefix(location *SourceRange, operator Token, operand ParseTree) *UpdateExpressionTree {
	return NewUpdateExpressionTree(location, operator, UpdateExpressionTree_OperatorPosition_PREFIX, operand)
}

// UpdateExpressionTree_Postfix creates a postfix update expression.
//
// Java: UpdateExpressionTree.postfix(SourceRange, Token, ParseTree)
func UpdateExpressionTree_Postfix(location *SourceRange, operator Token, operand ParseTree) *UpdateExpressionTree {
	return NewUpdateExpressionTree(location, operator, UpdateExpressionTree_OperatorPosition_POSTFIX, operand)
}
