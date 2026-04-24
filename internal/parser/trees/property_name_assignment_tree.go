// Port of com/google/javascript/jscomp/parsing/parser/trees/PropertyNameAssignmentTree.java
// @ closure-compiler v20260415

package trees

// PropertyNameAssignmentTree represents a property name/value pair in an object literal.
// Value may be nil in ES6 and above (shorthand property).
type PropertyNameAssignmentTree struct {
	ParseTree
	Name  *Token
	Value *ParseTree
}

// NewPropertyNameAssignmentTree constructs a PropertyNameAssignmentTree.
// Go-only: language-adapter
func NewPropertyNameAssignmentTree(location SourceRange, name *Token, value *ParseTree) *PropertyNameAssignmentTree {
	return &PropertyNameAssignmentTree{
		ParseTree: ParseTree{Type: ParseTreeTypePropertyNameAssignment, Location: location},
		Name:      name,
		Value:     value,
	}
}
