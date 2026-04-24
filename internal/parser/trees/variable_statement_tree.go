// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableStatementTree.java
// @ closure-compiler v20260415

package trees

type VariableStatementTree struct {
	Declarations *VariableDeclarationListTree
}

func NewVariableStatementTree(location SourceRange, declarations *VariableDeclarationListTree) *VariableStatementTree {
	return &VariableStatementTree{
		Declarations: declarations,
	}
}
