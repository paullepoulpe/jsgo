// Port of com/google/javascript/jscomp/parsing/parser/trees/WhileStatementTree.java
// @ closure-compiler v20260415

package trees

// WhileStatementTree represents a while loop.
type WhileStatementTree struct {
	ParseTree
	Condition *ParseTree
	Body      *ParseTree
}

func NewWhileStatementTree(
	location SourceRange,
	condition *ParseTree,
	body *ParseTree,
) *WhileStatementTree {
	return &WhileStatementTree{
		ParseTree: NewParseTree(ParseTreeTypeWhileStatement, location),
		Condition: condition,
		Body:      body,
	}
}
