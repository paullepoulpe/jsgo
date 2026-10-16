// Port of com/google/javascript/jscomp/parsing/parser/trees/CaseClauseTree.java @ closure-compiler v20260415

package parser

// CaseClauseTree is the parse tree for ParseTreeType CASE_CLAUSE.
type CaseClauseTree struct {
	ParseTreeBase

	Expression ParseTree
	Statements []ParseTree
}

// NewCaseClauseTree constructs a CaseClauseTree.
//
// Java: CaseClauseTree(SourceRange, ParseTree, ImmutableList)
func NewCaseClauseTree(location *SourceRange, expression ParseTree, statements []ParseTree) *CaseClauseTree {
	t := &CaseClauseTree{}
	t.InitParseTreeBase(t, ParseTreeType_CASE_CLAUSE, location)
	t.Expression = expression
	t.Statements = statements
	return t
}
