// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectRestTree.java
// @ closure-compiler v20260415

package trees

// ObjectRestTree represents an object rest element (...target).
type ObjectRestTree struct {
	ParseTree
	AssignmentTarget *ParseTree
}

func NewObjectRestTree(location SourceRange, assignmentTarget *ParseTree) *ObjectRestTree {
	return &ObjectRestTree{
		ParseTree:        NewParseTree(ParseTreeTypeObjectRest, location),
		AssignmentTarget: assignmentTarget,
	}
}
