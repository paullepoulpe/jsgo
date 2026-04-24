// Port of com/google/javascript/jscomp/parsing/parser/trees/LabelledStatementTree.java
// @ closure-compiler v20260415

package trees

// LabelledStatementTree represents a labelled statement.
type LabelledStatementTree struct {
	ParseTree
	Name      *IdentifierToken
	Statement *ParseTree
}

func NewLabelledStatementTree(
	location SourceRange,
	name *IdentifierToken,
	statement *ParseTree,
) *LabelledStatementTree {
	return &LabelledStatementTree{
		ParseTree: NewParseTree(ParseTreeTypeLabelledStatement, location),
		Name:      name,
		Statement: statement,
	}
}
