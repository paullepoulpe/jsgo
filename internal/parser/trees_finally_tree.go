// Port of com/google/javascript/jscomp/parsing/parser/trees/FinallyTree.java @ closure-compiler v20260415

package parser

// FinallyTree is the parse tree for ParseTreeType FINALLY.
type FinallyTree struct {
	ParseTreeBase

	Block ParseTree
}

// NewFinallyTree constructs a FinallyTree.
//
// Java: FinallyTree(SourceRange, ParseTree)
func NewFinallyTree(location *SourceRange, block ParseTree) *FinallyTree {
	t := &FinallyTree{}
	t.InitParseTreeBase(t, ParseTreeType_FINALLY, location)
	t.Block = block
	return t
}
