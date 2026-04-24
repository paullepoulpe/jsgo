// Port of com/google/javascript/jscomp/parsing/parser/trees/OptChainCallExpressionTree.java
// @ closure-compiler v20260415

package trees

// OptChainCallExpressionTree represents a call expression within an optional chain.
// IsStartOfOptionalChain indicates whether this node begins the optional chain.
type OptChainCallExpressionTree struct {
	ParseTree
	Operand                *ParseTree
	Arguments              *ArgumentListTree
	IsStartOfOptionalChain bool
	HasTrailingComma       bool
}

func NewOptChainCallExpressionTree(
	location SourceRange,
	operand *ParseTree,
	arguments *ArgumentListTree,
	isStartOfOptionalChain bool,
	hasTrailingComma bool,
) *OptChainCallExpressionTree {
	return &OptChainCallExpressionTree{
		ParseTree:              NewParseTree(ParseTreeTypeOptChainCallExpression, location),
		Operand:                operand,
		Arguments:              arguments,
		IsStartOfOptionalChain: isStartOfOptionalChain,
		HasTrailingComma:       hasTrailingComma,
	}
}
