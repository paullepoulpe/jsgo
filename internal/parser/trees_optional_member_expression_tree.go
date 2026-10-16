// Port of com/google/javascript/jscomp/parsing/parser/trees/OptionalMemberExpressionTree.java @ closure-compiler v20260415

package parser

// Creates a OptionalMemberExpressionTree that represents a getprop expression within an optional
// chain. Has an isOptionalChainStart field to indicate whether it is the start of an optional
// chain.
type OptionalMemberExpressionTree struct {
	ParseTreeBase

	Operand                ParseTree
	MemberName             *IdentifierToken
	IsStartOfOptionalChain bool
}

// NewOptionalMemberExpressionTree constructs a OptionalMemberExpressionTree.
//
// Java: OptionalMemberExpressionTree(SourceRange, ParseTree, IdentifierToken, boolean)
func NewOptionalMemberExpressionTree(location *SourceRange, operand ParseTree, memberName *IdentifierToken, isStartOfOptionalChain bool) *OptionalMemberExpressionTree {
	t := &OptionalMemberExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_OPT_CHAIN_MEMBER_EXPRESSION, location)
	t.Operand = operand
	t.MemberName = memberName
	t.IsStartOfOptionalChain = isStartOfOptionalChain
	return t
}
