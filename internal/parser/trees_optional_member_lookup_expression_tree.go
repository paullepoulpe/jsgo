// Port of com/google/javascript/jscomp/parsing/parser/trees/OptionalMemberLookupExpressionTree.java @ closure-compiler v20260415

package parser

// Creates a OptionalMemberExpressionTree that represents a getprop expression within an optional
// chain. Has an isOptionalChainStart field to indicate whether it is the start of an optional
// chain.
type OptionalMemberLookupExpressionTree struct {
	ParseTreeBase

	Operand                ParseTree
	MemberExpression       ParseTree
	IsStartOfOptionalChain bool
}

// NewOptionalMemberLookupExpressionTree constructs a OptionalMemberLookupExpressionTree.
//
// Java: OptionalMemberLookupExpressionTree(SourceRange, ParseTree, ParseTree, boolean)
func NewOptionalMemberLookupExpressionTree(location *SourceRange, operand ParseTree, memberExpression ParseTree, isStartOfOptionalChain bool) *OptionalMemberLookupExpressionTree {
	t := &OptionalMemberLookupExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_OPT_CHAIN_MEMBER_LOOKUP_EXPRESSION, location)
	t.Operand = operand
	t.MemberExpression = memberExpression
	t.IsStartOfOptionalChain = isStartOfOptionalChain
	return t
}
