// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableDeclarationTree.java
// @ closure-compiler v20260415

package trees

type VariableDeclarationTree struct {
	ParseTree
	// An IdentifierExpressionTree, ArrayPatternTree or ObjectLiteralTree
	Lvalue      *ParseTree
	Initializer *ParseTree
}

// Go-only: language-adapter
func NewVariableDeclarationTree(location SourceRange, lvalue *ParseTree, initializer *ParseTree) *VariableDeclarationTree {
	return &VariableDeclarationTree{
		ParseTree:   ParseTree{Type: ParseTreeTypeVariableDeclaration, Location: location},
		Lvalue:      lvalue,
		Initializer: initializer,
	}
}
