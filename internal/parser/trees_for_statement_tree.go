// Port of com/google/javascript/jscomp/parsing/parser/trees/ForStatementTree.java @ closure-compiler v20260415

package parser

// ForStatementTree is the parse tree for ParseTreeType FOR_STATEMENT.
type ForStatementTree struct {
	ParseTreeBase

	Initializer ParseTree
	Condition   ParseTree
	Increment   ParseTree
	Body        ParseTree
}

// NewForStatementTree constructs a ForStatementTree.
//
// Java: ForStatementTree(SourceRange, ParseTree, ParseTree, ParseTree, ParseTree)
func NewForStatementTree(location *SourceRange, variables ParseTree, condition ParseTree, increment ParseTree, body ParseTree) *ForStatementTree {
	t := &ForStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_FOR_STATEMENT, location)
	t.Initializer = variables
	t.Condition = condition
	t.Increment = increment
	t.Body = body
	return t
}
