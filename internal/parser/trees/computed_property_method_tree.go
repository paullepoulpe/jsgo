// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyMethodTree.java
// @ closure-compiler v20260415

package trees

// ComputedPropertyMethodTree represents a computed property method in an object literal or class,
// e.g. { [expr]() {} } or class { [expr]() {} }.
type ComputedPropertyMethodTree struct {
	ParseTree
	Property *ParseTree
	Method   *ParseTree
}

// NewComputedPropertyMethodTree constructs a ComputedPropertyMethodTree.
func NewComputedPropertyMethodTree(location SourceRange, property *ParseTree, method *ParseTree) *ComputedPropertyMethodTree {
	return &ComputedPropertyMethodTree{
		ParseTree: ParseTree{Type: ParseTreeTypeComputedPropertyMethod, Location: location},
		Property:  property,
		Method:    method,
	}
}
