// Port of com/google/javascript/jscomp/parsing/parser/trees/OptChainCallExpressionTree.java @ closure-compiler v20260415

package parser

// Creates a OptChainCallExpressionTree that represents a call expression within an optional chain.
// Has an isOptionalChainStart field to indicate whether it is the start of an optional chain.
type OptChainCallExpressionTree struct {
	ParseTreeBase

	Operand                ParseTree
	Arguments              *ArgumentListTree
	IsStartOfOptionalChain bool
	HasTrailingComma       bool
}

// NewOptChainCallExpressionTree constructs a OptChainCallExpressionTree.
//
// Java: OptChainCallExpressionTree(SourceRange, ParseTree, ArgumentListTree, boolean, boolean)
func NewOptChainCallExpressionTree(location *SourceRange, operand ParseTree, arguments *ArgumentListTree, isStartOfOptionalChain bool, hasTrailingComma bool) *OptChainCallExpressionTree {
	t := &OptChainCallExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_OPT_CHAIN_CALL_EXPRESSION, location)
	t.Operand = operand
	t.Arguments = arguments
	t.IsStartOfOptionalChain = isStartOfOptionalChain
	t.HasTrailingComma = hasTrailingComma
	return t
}
