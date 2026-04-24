// Port of com/google/javascript/jscomp/parsing/parser/trees/SwitchStatementTree.java
// @ closure-compiler v20260415

package trees

// SwitchStatementTree represents a switch statement.
type SwitchStatementTree struct {
	ParseTree
	Expression  *ParseTree
	CaseClauses []*ParseTree
}

func NewSwitchStatementTree(
	location SourceRange,
	expression *ParseTree,
	caseClauses []*ParseTree,
) *SwitchStatementTree {
	return &SwitchStatementTree{
		ParseTree:   NewParseTree(ParseTreeTypeSwitchStatement, location),
		Expression:  expression,
		CaseClauses: caseClauses,
	}
}
