// Port of com/google/javascript/jscomp/parsing/parser/trees/DoWhileStatementTree.java
// @ closure-compiler v20260415

package trees

// DoWhileStatementTree represents a do-while loop.
type DoWhileStatementTree struct {
	ParseTree
	Body      *ParseTree
	Condition *ParseTree
}

func NewDoWhileStatementTree(
	location SourceRange,
	body *ParseTree,
	condition *ParseTree,
) *DoWhileStatementTree {
	return &DoWhileStatementTree{
		ParseTree: NewParseTree(ParseTreeTypeDoWhileStatement, location),
		Body:      body,
		Condition: condition,
	}
}
