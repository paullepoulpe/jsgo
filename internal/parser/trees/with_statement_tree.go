// Port of com/google/javascript/jscomp/parsing/parser/trees/WithStatementTree.java
// @ closure-compiler v20260415

package trees

// WithStatementTree represents a with statement.
type WithStatementTree struct {
	ParseTree
	Expression *ParseTree
	Body       *ParseTree
}

func NewWithStatementTree(
	location SourceRange,
	expression *ParseTree,
	body *ParseTree,
) *WithStatementTree {
	return &WithStatementTree{
		ParseTree:  NewParseTree(ParseTreeTypeWithStatement, location),
		Expression: expression,
		Body:       body,
	}
}
