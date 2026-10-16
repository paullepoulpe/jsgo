// Port of com/google/javascript/jscomp/parsing/parser/trees/MemberLookupExpressionTree.java @ closure-compiler v20260415

package parser

// MemberLookupExpressionTree is the parse tree for ParseTreeType MEMBER_LOOKUP_EXPRESSION.
type MemberLookupExpressionTree struct {
	ParseTreeBase

	Operand          ParseTree
	MemberExpression ParseTree
}

// NewMemberLookupExpressionTree constructs a MemberLookupExpressionTree.
//
// Java: MemberLookupExpressionTree(SourceRange, ParseTree, ParseTree)
func NewMemberLookupExpressionTree(location *SourceRange, operand ParseTree, memberExpression ParseTree) *MemberLookupExpressionTree {
	t := &MemberLookupExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_MEMBER_LOOKUP_EXPRESSION, location)
	t.Operand = operand
	t.MemberExpression = memberExpression
	return t
}
