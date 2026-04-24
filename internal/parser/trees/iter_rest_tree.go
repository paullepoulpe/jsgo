// Port of com/google/javascript/jscomp/parsing/parser/trees/IterRestTree.java
// @ closure-compiler v20260415

package trees

// IterRestTree represents an iterable rest element (...target).
type IterRestTree struct {
	ParseTree
	AssignmentTarget *ParseTree
}

func NewIterRestTree(location SourceRange, assignmentTarget *ParseTree) *IterRestTree {
	return &IterRestTree{
		ParseTree:        NewParseTree(ParseTreeTypeIterRest, location),
		AssignmentTarget: assignmentTarget,
	}
}
