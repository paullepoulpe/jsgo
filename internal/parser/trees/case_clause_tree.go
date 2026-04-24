// Port of com/google/javascript/jscomp/parsing/parser/trees/CaseClauseTree.java
// @ closure-compiler v20260415

package trees

type CaseClauseTree struct {
	ParseTree
	Expression *ParseTree
	Statements []*ParseTree
}

func NewCaseClauseTree(location SourceRange, expression *ParseTree, statements []*ParseTree) *CaseClauseTree {
	return &CaseClauseTree{
		ParseTree:  ParseTree{Type: ParseTreeTypeCaseClause, Location: location},
		Expression: expression,
		Statements: statements,
	}
}
