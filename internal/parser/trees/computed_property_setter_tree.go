// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertySetterTree.java
// @ closure-compiler v20260415

package trees

// ComputedPropertySetterTree represents a computed property setter in an object literal or class.
type ComputedPropertySetterTree struct {
	ParseTree
	Property  *ParseTree
	Parameter *FormalParameterListTree
	IsStatic  bool
	Body      *BlockTree
}

func NewComputedPropertySetterTree(
	location SourceRange,
	property *ParseTree,
	isStatic bool,
	parameter *FormalParameterListTree,
	body *BlockTree,
) *ComputedPropertySetterTree {
	return &ComputedPropertySetterTree{
		ParseTree: NewParseTree(ParseTreeTypeComputedPropertySetter, location),
		Property:  property,
		IsStatic:  isStatic,
		Parameter: parameter,
		Body:      body,
	}
}
