// Port of com/google/javascript/jscomp/parsing/parser/trees/FinallyTree.java
// @ closure-compiler v20260415

package trees

type FinallyTree struct {
	Block ParseTree
}

func NewFinallyTree(location SourceRange, block ParseTree) *FinallyTree {
	return &FinallyTree{
		Block: block,
	}
}
