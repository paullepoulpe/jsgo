// Port of com/google/javascript/jscomp/parsing/parser/trees/ClassDeclarationTree.java
// @ closure-compiler v20260415

package trees

type ClassDeclarationTree struct {
	Name       *IdentifierToken
	SuperClass ParseTree
	Elements   []ParseTree
}

func NewClassDeclarationTree(location SourceRange, name *IdentifierToken, superClass ParseTree, elements []ParseTree) *ClassDeclarationTree {
	return &ClassDeclarationTree{
		Name:       name,
		SuperClass: superClass,
		Elements:   elements,
	}
}
