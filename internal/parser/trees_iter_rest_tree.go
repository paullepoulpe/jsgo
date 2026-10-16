// Port of com/google/javascript/jscomp/parsing/parser/trees/IterRestTree.java @ closure-compiler v20260415

package parser

// IterRestTree is the parse tree for ParseTreeType ITER_REST.
type IterRestTree struct {
	ParseTreeBase

	AssignmentTarget ParseTree
}

// NewIterRestTree constructs a IterRestTree.
//
// Java: IterRestTree(SourceRange, ParseTree)
func NewIterRestTree(location *SourceRange, assignmentTarget ParseTree) *IterRestTree {
	t := &IterRestTree{}
	t.InitParseTreeBase(t, ParseTreeType_ITER_REST, location)
	t.AssignmentTarget = assignmentTarget
	return t
}
