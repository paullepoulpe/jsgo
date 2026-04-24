// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableDeclarationListTree.java
// @ closure-compiler v20260415

package trees

type VariableDeclarationListTree struct {
	DeclarationType TokenType // const, var, let
	Declarations    []*VariableDeclarationTree
}

func NewVariableDeclarationListTree(location SourceRange, declarationType TokenType, declarations []*VariableDeclarationTree) *VariableDeclarationListTree {
	return &VariableDeclarationListTree{
		DeclarationType: declarationType,
		Declarations:    declarations,
	}
}
