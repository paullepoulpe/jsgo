// Port of com/google/javascript/jscomp/parsing/parser/trees/OptionalMemberExpressionTree.java
// @ closure-compiler v20260415

package trees

// OptionalMemberExpressionTree represents a getprop expression within an optional chain
// (operand?.memberName). IsStartOfOptionalChain indicates whether this node begins the chain.
type OptionalMemberExpressionTree struct {
	ParseTree
	Operand                *ParseTree
	MemberName             *IdentifierToken
	IsStartOfOptionalChain bool
}

func NewOptionalMemberExpressionTree(
	location SourceRange,
	operand *ParseTree,
	memberName *IdentifierToken,
	isStartOfOptionalChain bool,
) *OptionalMemberExpressionTree {
	return &OptionalMemberExpressionTree{
		ParseTree:              NewParseTree(ParseTreeTypeOptChainMemberExpression, location),
		Operand:                operand,
		MemberName:             memberName,
		IsStartOfOptionalChain: isStartOfOptionalChain,
	}
}
