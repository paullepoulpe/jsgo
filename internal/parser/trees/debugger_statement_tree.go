// Port of com/google/javascript/jscomp/parsing/parser/trees/DebuggerStatementTree.java
// @ closure-compiler v20260415

package trees

// DebuggerStatementTree represents a debugger statement.
type DebuggerStatementTree struct {
	ParseTree
}

func NewDebuggerStatementTree(location SourceRange) *DebuggerStatementTree {
	return &DebuggerStatementTree{
		ParseTree: NewParseTree(ParseTreeTypeDebuggerStatement, location),
	}
}
