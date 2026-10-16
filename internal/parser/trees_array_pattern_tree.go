// Port of com/google/javascript/jscomp/parsing/parser/trees/ArrayPatternTree.java @ closure-compiler v20260415

package parser

// ArrayPatternTree is the parse tree for ParseTreeType ARRAY_PATTERN.
type ArrayPatternTree struct {
	ParseTreeBase

	Elements []ParseTree
}

// NewArrayPatternTree constructs a ArrayPatternTree.
//
// Java: ArrayPatternTree(SourceRange, ImmutableList)
func NewArrayPatternTree(location *SourceRange, elements []ParseTree) *ArrayPatternTree {
	t := &ArrayPatternTree{}
	t.InitParseTreeBase(t, ParseTreeType_ARRAY_PATTERN, location)
	t.Elements = elements
	return t
}
