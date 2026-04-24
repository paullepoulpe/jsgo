// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyDefinitionTree.java
// @ closure-compiler v20260415

package trees

// ComputedPropertyDefinitionTree represents a computed property in an object literal,
// e.g. { [expr]: value }.
type ComputedPropertyDefinitionTree struct {
	ParseTree
	Property *ParseTree
	Value    *ParseTree
}

// NewComputedPropertyDefinitionTree constructs a ComputedPropertyDefinitionTree.
func NewComputedPropertyDefinitionTree(location SourceRange, property *ParseTree, value *ParseTree) *ComputedPropertyDefinitionTree {
	return &ComputedPropertyDefinitionTree{
		ParseTree: ParseTree{Type: ParseTreeTypeComputedPropertyDefinition, Location: location},
		Property:  property,
		Value:     value,
	}
}
