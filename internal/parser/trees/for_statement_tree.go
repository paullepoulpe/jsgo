// Port of com/google/javascript/jscomp/parsing/parser/trees/ForStatementTree.java
// @ closure-compiler v20260415

package trees

type ForStatementTree struct {
	ParseTree
	Initializer *ParseTree
	Condition   *ParseTree
	Increment   *ParseTree
	Body        *ParseTree
}

func NewForStatementTree(location SourceRange, variables *ParseTree, condition *ParseTree, increment *ParseTree, body *ParseTree) *ForStatementTree {
	return &ForStatementTree{
		ParseTree:   ParseTree{Type: ParseTreeTypeForStatement, Location: location},
		Initializer: variables,
		Condition:   condition,
		Increment:   increment,
		Body:        body,
	}
}
