// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableDeclarationListTree.java @ closure-compiler v20260415

package parser

// VariableDeclarationListTree is the parse tree for ParseTreeType VARIABLE_DECLARATION_LIST.
type VariableDeclarationListTree struct {
	ParseTreeBase

	DeclarationType TokenType // const, var, let
	Declarations    []*VariableDeclarationTree
}

// NewVariableDeclarationListTree constructs a VariableDeclarationListTree.
//
// Java: VariableDeclarationListTree(SourceRange, TokenType, ImmutableList)
func NewVariableDeclarationListTree(location *SourceRange, declarationType TokenType, declarations []*VariableDeclarationTree) *VariableDeclarationListTree {
	t := &VariableDeclarationListTree{}
	t.InitParseTreeBase(t, ParseTreeType_VARIABLE_DECLARATION_LIST, location)
	t.DeclarationType = declarationType
	t.Declarations = declarations
	return t
}
