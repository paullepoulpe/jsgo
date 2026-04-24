// Port of com/google/javascript/jscomp/parsing/parser/trees/CatchTree.java
// @ closure-compiler v20260415

package trees

type CatchTree struct {
	ParseTree
	Exception *ParseTree
	CatchBody *ParseTree
}

// Go-only: language-adapter
func NewCatchTree(location SourceRange, exception *ParseTree, catchBody *ParseTree) *CatchTree {
	return &CatchTree{
		ParseTree: ParseTree{Type: ParseTreeTypeCatch, Location: location},
		Exception: exception,
		CatchBody: catchBody,
	}
}
