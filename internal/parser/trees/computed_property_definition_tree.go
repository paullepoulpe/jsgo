// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyDefinitionTree.java
// @ closure-compiler v20260415

package trees

// ComputedPropertyDefinitionTree represents a computed property in an object literal.
type ComputedPropertyDefinitionTree struct {
	ParseTree
	Property *ParseTree
	Value    *ParseTree
}

func NewComputedPropertyDefinitionTree(
	location SourceRange,
	property *ParseTree,
	value *ParseTree,
) *ComputedPropertyDefinitionTree {
	return &ComputedPropertyDefinitionTree{
		ParseTree: NewParseTree(ParseTreeTypeComputedPropertyDefinition, location),
		Property:  property,
		Value:     value,
	}
}
