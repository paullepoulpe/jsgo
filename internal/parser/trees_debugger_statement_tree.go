// Port of com/google/javascript/jscomp/parsing/parser/trees/DebuggerStatementTree.java @ closure-compiler v20260415

package parser

// DebuggerStatementTree is the parse tree for ParseTreeType DEBUGGER_STATEMENT.
type DebuggerStatementTree struct {
	ParseTreeBase
}

// NewDebuggerStatementTree constructs a DebuggerStatementTree.
//
// Java: DebuggerStatementTree(SourceRange)
func NewDebuggerStatementTree(location *SourceRange) *DebuggerStatementTree {
	t := &DebuggerStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_DEBUGGER_STATEMENT, location)
	return t
}
