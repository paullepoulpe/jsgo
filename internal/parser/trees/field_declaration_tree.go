// Port of com/google/javascript/jscomp/parsing/parser/trees/FieldDeclarationTree.java
// @ closure-compiler v20260415

package trees

// FieldDeclarationTree represents a class field declaration
// (`name = expr;` or `name;` or `static name = expr;` inside a class body).
type FieldDeclarationTree struct {
	ParseTree
	Name        *IdentifierToken
	IsStatic    bool
	Initializer *ParseTree // nil when the field has no initialiser
}

// NewFieldDeclarationTree constructs a FieldDeclarationTree.
func NewFieldDeclarationTree(location SourceRange, name *IdentifierToken, isStatic bool, initializer *ParseTree) *FieldDeclarationTree {
	return &FieldDeclarationTree{
		ParseTree:   ParseTree{Type: ParseTreeTypeFieldDeclaration, Location: location},
		Name:        name,
		IsStatic:    isStatic,
		Initializer: initializer,
	}
}
