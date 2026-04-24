// Port of com/google/javascript/jscomp/parsing/parser/trees/OptionalMemberLookupExpressionTree.java
// @ closure-compiler v20260415

package trees

// OptionalMemberLookupExpressionTree represents a computed member access within an optional chain,
// e.g. a?.[expr] or a?.foo[expr] (continuation of chain).
// IsStartOfOptionalChain is true when this node is the ?. site itself.
type OptionalMemberLookupExpressionTree struct {
	ParseTree
	Operand                *ParseTree
	MemberExpression       *ParseTree
	IsStartOfOptionalChain bool
}

// NewOptionalMemberLookupExpressionTree constructs an OptionalMemberLookupExpressionTree.
func NewOptionalMemberLookupExpressionTree(location SourceRange, operand *ParseTree, memberExpression *ParseTree, isStartOfOptionalChain bool) *OptionalMemberLookupExpressionTree {
	return &OptionalMemberLookupExpressionTree{
		ParseTree:              ParseTree{Type: ParseTreeTypeOptChainMemberLookupExpression, Location: location},
		Operand:                operand,
		MemberExpression:       memberExpression,
		IsStartOfOptionalChain: isStartOfOptionalChain,
	}
}
