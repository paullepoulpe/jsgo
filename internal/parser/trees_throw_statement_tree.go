// Port of com/google/javascript/jscomp/parsing/parser/trees/ThrowStatementTree.java @ closure-compiler v20260415

package parser

// ThrowStatementTree is the parse tree for ParseTreeType THROW_STATEMENT.
type ThrowStatementTree struct {
	ParseTreeBase

	Value ParseTree
}

// NewThrowStatementTree constructs a ThrowStatementTree.
//
// Java: ThrowStatementTree(SourceRange, ParseTree)
func NewThrowStatementTree(location *SourceRange, value ParseTree) *ThrowStatementTree {
	t := &ThrowStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_THROW_STATEMENT, location)
	t.Value = value
	return t
}
