// Port of com/google/javascript/jscomp/parsing/parser/trees/ClassDeclarationTree.java
// @ closure-compiler v20260415

package trees

// ClassDeclarationTree represents a class declaration or expression.
type ClassDeclarationTree struct {
	ParseTree
	// Name is nil for anonymous class expressions.
	Name       *IdentifierToken
	SuperClass *ParseTree
	Elements   []*ParseTree
}

func NewClassDeclarationTree(
	location SourceRange,
	name *IdentifierToken,
	superClass *ParseTree,
	elements []*ParseTree,
) *ClassDeclarationTree {
	return &ClassDeclarationTree{
		ParseTree:  NewParseTree(ParseTreeTypeClassDeclaration, location),
		Name:       name,
		SuperClass: superClass,
		Elements:   elements,
	}
}
