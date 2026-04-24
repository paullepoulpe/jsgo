// Port of com/google/javascript/jscomp/parsing/parser/trees/CaseClauseTree.java
// @ closure-compiler v20260415

package trees

// CaseClauseTree represents a case clause in a switch statement.
type CaseClauseTree struct {
	ParseTree
	Expression *ParseTree
	Statements []*ParseTree
}

func NewCaseClauseTree(
	location SourceRange,
	expression *ParseTree,
	statements []*ParseTree,
) *CaseClauseTree {
	return &CaseClauseTree{
		ParseTree:  NewParseTree(ParseTreeTypeCaseClause, location),
		Expression: expression,
		Statements: statements,
	}
}
