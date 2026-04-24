// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableDeclarationTree.java
// @ closure-compiler v20260415

package trees

type VariableDeclarationTree struct {
	// An IdentifierExpressionTree, ArrayPatternTree or ObjectLiteralTree
	Lvalue      ParseTree
	Initializer ParseTree
}

func NewVariableDeclarationTree(location SourceRange, lvalue ParseTree, initializer ParseTree) *VariableDeclarationTree {
	return &VariableDeclarationTree{
		Lvalue:      lvalue,
		Initializer: initializer,
	}
}
