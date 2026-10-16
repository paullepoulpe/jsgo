// Port of com/google/javascript/jscomp/parsing/parser/trees/CatchTree.java @ closure-compiler v20260415

package parser

// CatchTree is the parse tree for ParseTreeType CATCH.
type CatchTree struct {
	ParseTreeBase

	Exception ParseTree
	CatchBody ParseTree
}

// NewCatchTree constructs a CatchTree.
//
// Java: CatchTree(SourceRange, ParseTree, ParseTree)
func NewCatchTree(location *SourceRange, exception ParseTree, catchBody ParseTree) *CatchTree {
	t := &CatchTree{}
	t.InitParseTreeBase(t, ParseTreeType_CATCH, location)
	t.Exception = exception
	t.CatchBody = catchBody
	return t
}
