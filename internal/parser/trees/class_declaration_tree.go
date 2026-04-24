// Port of com/google/javascript/jscomp/parsing/parser/trees/ClassDeclarationTree.java
// @ closure-compiler v20260415

package trees

type ClassDeclarationTree struct {
	ParseTree
	Name       *IdentifierToken
	SuperClass *ParseTree
	Elements   []*ParseTree
}

// Go-only: language-adapter
func NewClassDeclarationTree(location SourceRange, name *IdentifierToken, superClass *ParseTree, elements []*ParseTree) *ClassDeclarationTree {
	return &ClassDeclarationTree{
		ParseTree:  ParseTree{Type: ParseTreeTypeClassDeclaration, Location: location},
		Name:       name,
		SuperClass: superClass,
		Elements:   elements,
	}
}
