// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableDeclarationListTree.java
// @ closure-compiler v20260415

package trees

type VariableDeclarationListTree struct {
	ParseTree
	DeclarationType TokenType // const, var, let
	Declarations    []*VariableDeclarationTree
}

func NewVariableDeclarationListTree(location SourceRange, declarationType TokenType, declarations []*VariableDeclarationTree) *VariableDeclarationListTree {
	return &VariableDeclarationListTree{
		ParseTree:       ParseTree{Type: ParseTreeTypeVariableDeclarationList, Location: location},
		DeclarationType: declarationType,
		Declarations:    declarations,
	}
}
