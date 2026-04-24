// Port of com/google/javascript/jscomp/parsing/parser/trees/MemberExpressionTree.java
// @ closure-compiler v20260415

package trees

// MemberExpressionTree represents a member access expression, e.g. obj.foo.
type MemberExpressionTree struct {
	ParseTree
	Operand    *ParseTree
	MemberName *IdentifierToken
}

// NewMemberExpressionTree constructs a MemberExpressionTree.
// Go-only: language-adapter
func NewMemberExpressionTree(location SourceRange, operand *ParseTree, memberName *IdentifierToken) *MemberExpressionTree {
	return &MemberExpressionTree{
		ParseTree:  ParseTree{Type: ParseTreeTypeMemberExpression, Location: location},
		Operand:    operand,
		MemberName: memberName,
	}
}
