// Port of com/google/javascript/jscomp/parsing/parser/trees/MemberLookupExpressionTree.java
// @ closure-compiler v20260415

package trees

// MemberLookupExpressionTree represents a computed member expression (operand[memberExpression]).
type MemberLookupExpressionTree struct {
	ParseTree
	Operand          *ParseTree
	MemberExpression *ParseTree
}

func NewMemberLookupExpressionTree(
	location SourceRange,
	operand *ParseTree,
	memberExpression *ParseTree,
) *MemberLookupExpressionTree {
	return &MemberLookupExpressionTree{
		ParseTree:        NewParseTree(ParseTreeTypeMemberLookupExpression, location),
		Operand:          operand,
		MemberExpression: memberExpression,
	}
}
