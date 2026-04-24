// Port of com/google/javascript/jscomp/parsing/parser/trees/CatchTree.java
// @ closure-compiler v20260415

package trees

// CatchTree represents a catch clause in a try statement.
type CatchTree struct {
	ParseTree
	Exception *ParseTree
	CatchBody *ParseTree
}

func NewCatchTree(location SourceRange, exception *ParseTree, catchBody *ParseTree) *CatchTree {
	return &CatchTree{
		ParseTree: NewParseTree(ParseTreeTypeCatch, location),
		Exception: exception,
		CatchBody: catchBody,
	}
}
