// Port of com/google/javascript/jscomp/parsing/parser/trees/ForStatementTree.java
// @ closure-compiler v20260415

package trees

// ForStatementTree represents a traditional for loop.
type ForStatementTree struct {
	ParseTree
	Initializer *ParseTree
	Condition   *ParseTree
	Increment   *ParseTree
	Body        *ParseTree
}

func NewForStatementTree(
	location SourceRange,
	variables *ParseTree,
	condition *ParseTree,
	increment *ParseTree,
	body *ParseTree,
) *ForStatementTree {
	return &ForStatementTree{
		ParseTree:   NewParseTree(ParseTreeTypeForStatement, location),
		Initializer: variables,
		Condition:   condition,
		Increment:   increment,
		Body:        body,
	}
}
