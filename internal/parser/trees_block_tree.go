// Port of com/google/javascript/jscomp/parsing/parser/trees/BlockTree.java @ closure-compiler v20260415

package parser

// BlockTree is the parse tree for ParseTreeType BLOCK.
type BlockTree struct {
	ParseTreeBase

	Statements []ParseTree
}

// NewBlockTree constructs a BlockTree.
//
// Java: BlockTree(SourceRange, ImmutableList)
func NewBlockTree(location *SourceRange, statements []ParseTree) *BlockTree {
	t := &BlockTree{}
	t.InitParseTreeBase(t, ParseTreeType_BLOCK, location)
	t.Statements = statements
	return t
}
