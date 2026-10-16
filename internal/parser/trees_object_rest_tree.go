// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectRestTree.java @ closure-compiler v20260415

package parser

// ObjectRestTree is the parse tree for ParseTreeType OBJECT_REST.
type ObjectRestTree struct {
	ParseTreeBase

	AssignmentTarget ParseTree
}

// NewObjectRestTree constructs a ObjectRestTree.
//
// Java: ObjectRestTree(SourceRange, ParseTree)
func NewObjectRestTree(location *SourceRange, assignmentTarget ParseTree) *ObjectRestTree {
	t := &ObjectRestTree{}
	t.InitParseTreeBase(t, ParseTreeType_OBJECT_REST, location)
	t.AssignmentTarget = assignmentTarget
	return t
}
