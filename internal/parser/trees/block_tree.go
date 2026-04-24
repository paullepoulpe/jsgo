// Port of com/google/javascript/jscomp/parsing/parser/trees/BlockTree.java
// @ closure-compiler v20260415

package trees

// BlockTree represents a block statement.
type BlockTree struct {
	ParseTree
	Statements []*ParseTree
}

func NewBlockTree(location SourceRange, statements []*ParseTree) *BlockTree {
	return &BlockTree{
		ParseTree:  NewParseTree(ParseTreeTypeBlock, location),
		Statements: statements,
	}
}
