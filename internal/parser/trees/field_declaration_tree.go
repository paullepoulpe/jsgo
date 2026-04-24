// Port of com/google/javascript/jscomp/parsing/parser/trees/FieldDeclarationTree.java
// @ closure-compiler v20260415

package trees

// FieldDeclarationTree represents a public property field in a class.
type FieldDeclarationTree struct {
	ParseTree
	Name     *IdentifierToken
	IsStatic bool
	// Initializer is nil when the field has no initializer.
	Initializer *ParseTree
}

func NewFieldDeclarationTree(
	location SourceRange,
	name *IdentifierToken,
	isStatic bool,
	initializer *ParseTree,
) *FieldDeclarationTree {
	return &FieldDeclarationTree{
		ParseTree:   NewParseTree(ParseTreeTypeFieldDeclaration, location),
		Name:        name,
		IsStatic:    isStatic,
		Initializer: initializer,
	}
}
