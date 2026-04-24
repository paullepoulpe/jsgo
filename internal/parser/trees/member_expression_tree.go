// Port of com/google/javascript/jscomp/parsing/parser/trees/MemberExpressionTree.java
// @ closure-compiler v20260415

package trees

// MemberExpressionTree represents a member expression (operand.memberName).
type MemberExpressionTree struct {
	ParseTree
	Operand    *ParseTree
	MemberName *IdentifierToken
}

func NewMemberExpressionTree(
	location SourceRange,
	operand *ParseTree,
	memberName *IdentifierToken,
) *MemberExpressionTree {
	return &MemberExpressionTree{
		ParseTree:  NewParseTree(ParseTreeTypeMemberExpression, location),
		Operand:    operand,
		MemberName: memberName,
	}
}
