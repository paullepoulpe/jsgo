// Port of com/google/javascript/jscomp/parsing/parser/trees/SwitchStatementTree.java @ closure-compiler v20260415

package parser

// SwitchStatementTree is the parse tree for ParseTreeType SWITCH_STATEMENT.
type SwitchStatementTree struct {
	ParseTreeBase

	Expression  ParseTree
	CaseClauses []ParseTree
}

// NewSwitchStatementTree constructs a SwitchStatementTree.
//
// Java: SwitchStatementTree(SourceRange, ParseTree, ImmutableList)
func NewSwitchStatementTree(location *SourceRange, expression ParseTree, caseClauses []ParseTree) *SwitchStatementTree {
	t := &SwitchStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_SWITCH_STATEMENT, location)
	t.Expression = expression
	t.CaseClauses = caseClauses
	return t
}
