// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertySetterTree.java
// @ closure-compiler v20260415

package trees

// ComputedPropertySetterTree represents a computed setter in an object literal or class,
// e.g. { set [expr](v) {} }.
type ComputedPropertySetterTree struct {
	ParseTree
	Property  *ParseTree
	IsStatic  bool
	Parameter *FormalParameterListTree
	Body      *BlockTree
}

// NewComputedPropertySetterTree constructs a ComputedPropertySetterTree.
func NewComputedPropertySetterTree(location SourceRange, property *ParseTree, isStatic bool, parameter *FormalParameterListTree, body *BlockTree) *ComputedPropertySetterTree {
	return &ComputedPropertySetterTree{
		ParseTree: ParseTree{Type: ParseTreeTypeComputedPropertySetter, Location: location},
		Property:  property,
		IsStatic:  isStatic,
		Parameter: parameter,
		Body:      body,
	}
}
