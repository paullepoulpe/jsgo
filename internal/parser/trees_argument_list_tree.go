// Port of com/google/javascript/jscomp/parsing/parser/trees/ArgumentListTree.java @ closure-compiler v20260415

package parser

// ArgumentListTree is the parse tree for ParseTreeType ARGUMENT_LIST.
type ArgumentListTree struct {
	ParseTreeBase

	Arguments        []ParseTree
	HasTrailingComma bool
	CommaPositions   []*SourcePosition
}

// NewArgumentListTree constructs a ArgumentListTree.
//
// Java: ArgumentListTree(SourceRange, ImmutableList, boolean, ImmutableList)
func NewArgumentListTree(location *SourceRange, arguments []ParseTree, hasTrailingComma bool, commaPositions []*SourcePosition) *ArgumentListTree {
	t := &ArgumentListTree{}
	t.InitParseTreeBase(t, ParseTreeType_ARGUMENT_LIST, location)
	t.Arguments = arguments
	t.HasTrailingComma = hasTrailingComma
	t.CommaPositions = commaPositions
	return t
}
