// Port of com/google/javascript/jscomp/parsing/parser/trees/BlockTree.java
// @ closure-compiler v20260415

package trees

type BlockTree struct {
	Statements []ParseTree
}

func NewBlockTree(location SourceRange, statements []ParseTree) *BlockTree {
	return &BlockTree{
		Statements: statements,
	}
}
