// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableDeclarationTree.java
// @ closure-compiler v20260415

package trees

// VariableDeclarationTree represents a single variable declaration
// (an IdentifierExpression, ArrayPattern, or ObjectLiteral lvalue).
type VariableDeclarationTree struct {
	ParseTree
	Lvalue      *ParseTree
	// Initializer is nil when the declaration has no initializer.
	Initializer *ParseTree
}

func NewVariableDeclarationTree(
	location SourceRange,
	lvalue *ParseTree,
	initializer *ParseTree,
) *VariableDeclarationTree {
	return &VariableDeclarationTree{
		ParseTree:   NewParseTree(ParseTreeTypeVariableDeclaration, location),
		Lvalue:      lvalue,
		Initializer: initializer,
	}
}
