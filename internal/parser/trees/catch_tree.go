// Port of com/google/javascript/jscomp/parsing/parser/trees/CatchTree.java
// @ closure-compiler v20260415

package trees

type CatchTree struct {
	Exception ParseTree
	CatchBody ParseTree
}

func NewCatchTree(location SourceRange, exception ParseTree, catchBody ParseTree) *CatchTree {
	return &CatchTree{
		Exception: exception,
		CatchBody: catchBody,
	}
}
