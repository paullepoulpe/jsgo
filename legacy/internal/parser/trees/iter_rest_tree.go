// Port of com/google/javascript/jscomp/parsing/parser/trees/IterRestTree.java
// @ closure-compiler v20260415

package trees

// IterRestTree represents an iterable rest element in array destructuring,
// e.g. ...rest in [a, ...rest].
type IterRestTree struct {
	ParseTree
	AssignmentTarget *ParseTree
}

// NewIterRestTree constructs an IterRestTree.
// Go-only: language-adapter
func NewIterRestTree(location SourceRange, assignmentTarget *ParseTree) *IterRestTree {
	return &IterRestTree{
		ParseTree:        ParseTree{Type: ParseTreeTypeIterRest, Location: location},
		AssignmentTarget: assignmentTarget,
	}
}
