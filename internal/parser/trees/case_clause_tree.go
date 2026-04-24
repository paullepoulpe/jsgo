// Port of com/google/javascript/jscomp/parsing/parser/trees/CaseClauseTree.java
// @ closure-compiler v20260415

package trees

type CaseClauseTree struct {
	Expression ParseTree
	Statements []ParseTree
}

func NewCaseClauseTree(location SourceRange, expression ParseTree, statements []ParseTree) *CaseClauseTree {
	return &CaseClauseTree{
		Expression: expression,
		Statements: statements,
	}
}
