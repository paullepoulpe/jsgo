// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableDeclarationListTree.java
// @ closure-compiler v20260415

package trees

// VariableDeclarationListTree represents a variable declaration list (var/let/const).
type VariableDeclarationListTree struct {
	ParseTree
	// DeclarationType is the keyword token type: var, let, or const.
	DeclarationType TokenType
	Declarations    []*VariableDeclarationTree
}

func NewVariableDeclarationListTree(
	location SourceRange,
	declarationType TokenType,
	declarations []*VariableDeclarationTree,
) *VariableDeclarationListTree {
	return &VariableDeclarationListTree{
		ParseTree:       NewParseTree(ParseTreeTypeVariableDeclarationList, location),
		DeclarationType: declarationType,
		Declarations:    declarations,
	}
}
