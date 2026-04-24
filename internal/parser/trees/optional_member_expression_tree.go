// Port of com/google/javascript/jscomp/parsing/parser/trees/OptionalMemberExpressionTree.java
// @ closure-compiler v20260415

package trees

// OptionalMemberExpressionTree represents a member access within an optional chain,
// e.g. a?.foo or a?.foo.bar (continuation of chain).
// IsStartOfOptionalChain is true when this node is the ?. site itself.
type OptionalMemberExpressionTree struct {
	ParseTree
	Operand                *ParseTree
	MemberName             *IdentifierToken
	IsStartOfOptionalChain bool
}

// NewOptionalMemberExpressionTree constructs an OptionalMemberExpressionTree.
func NewOptionalMemberExpressionTree(location SourceRange, operand *ParseTree, memberName *IdentifierToken, isStartOfOptionalChain bool) *OptionalMemberExpressionTree {
	return &OptionalMemberExpressionTree{
		ParseTree:              ParseTree{Type: ParseTreeTypeOptChainMemberExpression, Location: location},
		Operand:                operand,
		MemberName:             memberName,
		IsStartOfOptionalChain: isStartOfOptionalChain,
	}
}
