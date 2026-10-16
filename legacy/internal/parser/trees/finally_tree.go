// Port of com/google/javascript/jscomp/parsing/parser/trees/FinallyTree.java
// @ closure-compiler v20260415

package trees

type FinallyTree struct {
	ParseTree
	Block *ParseTree
}

// Go-only: language-adapter
func NewFinallyTree(location SourceRange, block *ParseTree) *FinallyTree {
	return &FinallyTree{
		ParseTree: ParseTree{Type: ParseTreeTypeFinally, Location: location},
		Block:     block,
	}
}
