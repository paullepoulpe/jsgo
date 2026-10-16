// Port of com/google/javascript/jscomp/parsing/parser/trees/IterSpreadTree.java @ closure-compiler v20260415

package parser

// IterSpreadTree is the parse tree for ParseTreeType ITER_SPREAD.
type IterSpreadTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewIterSpreadTree constructs a IterSpreadTree.
//
// Java: IterSpreadTree(SourceRange, ParseTree)
func NewIterSpreadTree(location *SourceRange, expression ParseTree) *IterSpreadTree {
	t := &IterSpreadTree{}
	t.InitParseTreeBase(t, ParseTreeType_ITER_SPREAD, location)
	t.Expression = expression
	return t
}
