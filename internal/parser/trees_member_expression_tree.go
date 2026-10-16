// Port of com/google/javascript/jscomp/parsing/parser/trees/MemberExpressionTree.java @ closure-compiler v20260415

package parser

// MemberExpressionTree is the parse tree for ParseTreeType MEMBER_EXPRESSION.
type MemberExpressionTree struct {
	ParseTreeBase

	Operand    ParseTree
	MemberName *IdentifierToken
}

// NewMemberExpressionTree constructs a MemberExpressionTree.
//
// Java: MemberExpressionTree(SourceRange, ParseTree, IdentifierToken)
func NewMemberExpressionTree(location *SourceRange, operand ParseTree, memberName *IdentifierToken) *MemberExpressionTree {
	t := &MemberExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_MEMBER_EXPRESSION, location)
	t.Operand = operand
	t.MemberName = memberName
	return t
}
