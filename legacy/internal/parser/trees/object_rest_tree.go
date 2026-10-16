// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectRestTree.java
// @ closure-compiler v20260415

package trees

// ObjectRestTree represents the rest element in object destructuring, e.g. { ...rest }.
type ObjectRestTree struct {
	ParseTree
	AssignmentTarget *ParseTree
}

// NewObjectRestTree constructs an ObjectRestTree.
// Go-only: language-adapter
func NewObjectRestTree(location SourceRange, assignmentTarget *ParseTree) *ObjectRestTree {
	return &ObjectRestTree{
		ParseTree:        ParseTree{Type: ParseTreeTypeObjectRest, Location: location},
		AssignmentTarget: assignmentTarget,
	}
}
