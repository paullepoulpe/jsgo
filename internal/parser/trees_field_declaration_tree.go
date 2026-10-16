// Port of com/google/javascript/jscomp/parsing/parser/trees/FieldDeclarationTree.java @ closure-compiler v20260415

package parser

// Tree for public property fields
type FieldDeclarationTree struct {
	ParseTreeBase

	Name        *IdentifierToken
	IsStatic    bool
	Initializer ParseTree
}

// NewFieldDeclarationTree constructs a FieldDeclarationTree.
//
// Java: FieldDeclarationTree(SourceRange, IdentifierToken, boolean, ParseTree)
func NewFieldDeclarationTree(location *SourceRange, name *IdentifierToken, isStatic bool, initializer ParseTree) *FieldDeclarationTree {
	t := &FieldDeclarationTree{}
	t.InitParseTreeBase(t, ParseTreeType_FIELD_DECLARATION, location)
	t.Name = name
	t.IsStatic = isStatic
	t.Initializer = initializer
	return t
}
