// Port of com/google/javascript/jscomp/parsing/parser/trees/ComprehensionForTree.java @ closure-compiler v20260415

package parser

// ComprehensionForTree is the parse tree for ParseTreeType COMPREHENSION_FOR.
type ComprehensionForTree struct {
	ParseTreeBase

	Initializer ParseTree
	Collection  ParseTree
}

// NewComprehensionForTree constructs a ComprehensionForTree.
//
// Java: ComprehensionForTree(SourceRange, ParseTree, ParseTree)
func NewComprehensionForTree(location *SourceRange, initializer ParseTree, collection ParseTree) *ComprehensionForTree {
	t := &ComprehensionForTree{}
	t.InitParseTreeBase(t, ParseTreeType_COMPREHENSION_FOR, location)
	t.Initializer = initializer
	t.Collection = collection
	return t
}
