// Port of com/google/javascript/jscomp/parsing/parser/trees/NullTree.java @ closure-compiler v20260415

package parser

// Represents a null parse tree element.
type NullTree struct {
	ParseTreeBase
}

// NewNullTree constructs a NullTree.
//
// Java: NullTree(SourceRange)
func NewNullTree(location *SourceRange) *NullTree {
	t := &NullTree{}
	t.InitParseTreeBase(t, ParseTreeType_NULL, location)
	return t
}
