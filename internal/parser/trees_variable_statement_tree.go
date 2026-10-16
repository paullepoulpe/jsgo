// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableStatementTree.java @ closure-compiler v20260415

package parser

// VariableStatementTree is the parse tree for ParseTreeType VARIABLE_STATEMENT.
type VariableStatementTree struct {
	ParseTreeBase

	Declarations *VariableDeclarationListTree
}

// NewVariableStatementTree constructs a VariableStatementTree.
//
// Java: VariableStatementTree(SourceRange, VariableDeclarationListTree)
func NewVariableStatementTree(location *SourceRange, declarations *VariableDeclarationListTree) *VariableStatementTree {
	t := &VariableStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_VARIABLE_STATEMENT, location)
	t.Declarations = declarations
	return t
}
