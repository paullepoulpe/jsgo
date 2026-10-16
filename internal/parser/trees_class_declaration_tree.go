// Port of com/google/javascript/jscomp/parsing/parser/trees/ClassDeclarationTree.java @ closure-compiler v20260415

package parser

// ClassDeclarationTree is the parse tree for ParseTreeType CLASS_DECLARATION.
type ClassDeclarationTree struct {
	ParseTreeBase

	Name       *IdentifierToken
	SuperClass ParseTree
	Elements   []ParseTree
}

// NewClassDeclarationTree constructs a ClassDeclarationTree.
//
// Java: ClassDeclarationTree(SourceRange, IdentifierToken, ParseTree, ImmutableList)
func NewClassDeclarationTree(location *SourceRange, name *IdentifierToken, superClass ParseTree, elements []ParseTree) *ClassDeclarationTree {
	t := &ClassDeclarationTree{}
	t.InitParseTreeBase(t, ParseTreeType_CLASS_DECLARATION, location)
	t.Name = name
	t.SuperClass = superClass
	t.Elements = elements
	return t
}
