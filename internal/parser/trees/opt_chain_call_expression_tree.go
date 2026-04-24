// Port of com/google/javascript/jscomp/parsing/parser/trees/OptChainCallExpressionTree.java
// @ closure-compiler v20260415

package trees

// OptChainCallExpressionTree represents a call expression within an optional chain,
// e.g. a?.() or a?.b().
// IsStartOfOptionalChain is true when this node is the ?. site itself.
type OptChainCallExpressionTree struct {
	ParseTree
	Operand                *ParseTree
	Arguments              *ArgumentListTree
	IsStartOfOptionalChain bool
	HasTrailingComma       bool
}

// NewOptChainCallExpressionTree constructs an OptChainCallExpressionTree.
func NewOptChainCallExpressionTree(location SourceRange, operand *ParseTree, arguments *ArgumentListTree, isStartOfOptionalChain bool, hasTrailingComma bool) *OptChainCallExpressionTree {
	return &OptChainCallExpressionTree{
		ParseTree:              ParseTree{Type: ParseTreeTypeOptChainCallExpression, Location: location},
		Operand:                operand,
		Arguments:              arguments,
		IsStartOfOptionalChain: isStartOfOptionalChain,
		HasTrailingComma:       hasTrailingComma,
	}
}
