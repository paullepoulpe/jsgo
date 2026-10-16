// Port of com/google/javascript/jscomp/parsing/parser/trees/LabelledStatementTree.java @ closure-compiler v20260415

package parser

// LabelledStatementTree is the parse tree for ParseTreeType LABELLED_STATEMENT.
type LabelledStatementTree struct {
	ParseTreeBase

	Name      *IdentifierToken
	Statement ParseTree
}

// NewLabelledStatementTree constructs a LabelledStatementTree.
//
// Java: LabelledStatementTree(SourceRange, IdentifierToken, ParseTree)
func NewLabelledStatementTree(location *SourceRange, name *IdentifierToken, statement ParseTree) *LabelledStatementTree {
	t := &LabelledStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_LABELLED_STATEMENT, location)
	t.Name = name
	t.Statement = statement
	return t
}
