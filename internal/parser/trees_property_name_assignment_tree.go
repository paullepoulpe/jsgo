// Port of com/google/javascript/jscomp/parsing/parser/trees/PropertyNameAssignmentTree.java @ closure-compiler v20260415

package parser

// PropertyNameAssignmentTree is the parse tree for ParseTreeType PROPERTY_NAME_ASSIGNMENT.
type PropertyNameAssignmentTree struct {
	ParseTreeBase

	Name Token
	// May be null in ES6 and above.
	Value ParseTree
}

// NewPropertyNameAssignmentTree constructs a PropertyNameAssignmentTree.
//
// Java: PropertyNameAssignmentTree(SourceRange, Token, ParseTree)
func NewPropertyNameAssignmentTree(location *SourceRange, name Token, value ParseTree) *PropertyNameAssignmentTree {
	t := &PropertyNameAssignmentTree{}
	t.InitParseTreeBase(t, ParseTreeType_PROPERTY_NAME_ASSIGNMENT, location)
	t.Name = name
	t.Value = value
	return t
}
