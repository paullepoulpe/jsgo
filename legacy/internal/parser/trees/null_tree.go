// Port of com/google/javascript/jscomp/parsing/parser/trees/NullTree.java
// @ closure-compiler v20260415

package trees

// NullTree represents a null parse tree element.
type NullTree struct {
	ParseTree
}

// NewNullTree constructs a NullTree.
// Go-only: language-adapter
func NewNullTree(location SourceRange) *NullTree {
	return &NullTree{
		ParseTree: ParseTree{Type: ParseTreeTypeNull, Location: location},
	}
}
