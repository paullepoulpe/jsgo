// Port of com/google/javascript/jscomp/parsing/parser/trees/NullTree.java
// @ closure-compiler v20260415

package trees

// NullTree represents a null parse-tree element (an absent optional node).
type NullTree struct {
	ParseTree
}

func NewNullTree(location SourceRange) *NullTree {
	return &NullTree{
		ParseTree: NewParseTree(ParseTreeTypeNull, location),
	}
}
