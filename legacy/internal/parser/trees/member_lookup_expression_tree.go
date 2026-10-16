// Port of com/google/javascript/jscomp/parsing/parser/trees/MemberLookupExpressionTree.java
// @ closure-compiler v20260415

package trees

// MemberLookupExpressionTree represents a computed member access expression,
// e.g. obj[expr].
type MemberLookupExpressionTree struct {
	ParseTree
	Operand          *ParseTree
	MemberExpression *ParseTree
}

// NewMemberLookupExpressionTree constructs a MemberLookupExpressionTree.
// Go-only: language-adapter
func NewMemberLookupExpressionTree(location SourceRange, operand *ParseTree, memberExpression *ParseTree) *MemberLookupExpressionTree {
	return &MemberLookupExpressionTree{
		ParseTree:        ParseTree{Type: ParseTreeTypeMemberLookupExpression, Location: location},
		Operand:          operand,
		MemberExpression: memberExpression,
	}
}
