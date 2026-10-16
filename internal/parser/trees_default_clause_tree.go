// Port of com/google/javascript/jscomp/parsing/parser/trees/DefaultClauseTree.java @ closure-compiler v20260415

package parser

// DefaultClauseTree is the parse tree for ParseTreeType DEFAULT_CLAUSE.
type DefaultClauseTree struct {
	ParseTreeBase

	Statements []ParseTree
}

// NewDefaultClauseTree constructs a DefaultClauseTree.
//
// Java: DefaultClauseTree(SourceRange, ImmutableList)
func NewDefaultClauseTree(location *SourceRange, statements []ParseTree) *DefaultClauseTree {
	t := &DefaultClauseTree{}
	t.InitParseTreeBase(t, ParseTreeType_DEFAULT_CLAUSE, location)
	t.Statements = statements
	return t
}
