// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableStatementTree.java
// @ closure-compiler v20260415

package trees

// VariableStatementTree represents a variable statement (var/let/const declaration).
type VariableStatementTree struct {
	ParseTree
	Declarations *VariableDeclarationListTree
}

func NewVariableStatementTree(
	location SourceRange,
	declarations *VariableDeclarationListTree,
) *VariableStatementTree {
	return &VariableStatementTree{
		ParseTree:    NewParseTree(ParseTreeTypeVariableStatement, location),
		Declarations: declarations,
	}
}
