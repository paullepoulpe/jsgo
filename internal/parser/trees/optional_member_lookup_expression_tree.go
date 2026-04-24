// Port of com/google/javascript/jscomp/parsing/parser/trees/OptionalMemberLookupExpressionTree.java
// @ closure-compiler v20260415

package trees

// OptionalMemberLookupExpressionTree represents a computed member expression within an optional
// chain (operand?.[memberExpression]). IsStartOfOptionalChain indicates whether this begins the
// chain.
type OptionalMemberLookupExpressionTree struct {
	ParseTree
	Operand                *ParseTree
	MemberExpression       *ParseTree
	IsStartOfOptionalChain bool
}

func NewOptionalMemberLookupExpressionTree(
	location SourceRange,
	operand *ParseTree,
	memberExpression *ParseTree,
	isStartOfOptionalChain bool,
) *OptionalMemberLookupExpressionTree {
	return &OptionalMemberLookupExpressionTree{
		ParseTree:              NewParseTree(ParseTreeTypeOptChainMemberLookupExpression, location),
		Operand:                operand,
		MemberExpression:       memberExpression,
		IsStartOfOptionalChain: isStartOfOptionalChain,
	}
}
