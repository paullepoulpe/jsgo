// Port of com/google/javascript/jscomp/parsing/parser/trees/BlockTree.java
// @ closure-compiler v20260415

package trees

// BlockTree represents a block statement, e.g. { stmt1; stmt2; }.
type BlockTree struct {
	ParseTree
	Statements []*ParseTree
}

// NewBlockTree constructs a BlockTree.
// Go-only: language-adapter
func NewBlockTree(location SourceRange, statements []*ParseTree) *BlockTree {
	return &BlockTree{
		ParseTree:  ParseTree{Type: ParseTreeTypeBlock, Location: location},
		Statements: statements,
	}
}
