// Port of com/google/javascript/jscomp/parsing/parser/trees/IfStatementTree.java @ closure-compiler v20260415

package parser

// IfStatementTree is the parse tree for ParseTreeType IF_STATEMENT.
type IfStatementTree struct {
	ParseTreeBase

	Condition  ParseTree
	IfClause   ParseTree
	ElseClause ParseTree
}

// NewIfStatementTree constructs a IfStatementTree.
//
// Java: IfStatementTree(SourceRange, ParseTree, ParseTree, ParseTree)
func NewIfStatementTree(location *SourceRange, condition ParseTree, ifClause ParseTree, elseClause ParseTree) *IfStatementTree {
	t := &IfStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_IF_STATEMENT, location)
	t.Condition = condition
	t.IfClause = ifClause
	t.ElseClause = elseClause
	return t
}
