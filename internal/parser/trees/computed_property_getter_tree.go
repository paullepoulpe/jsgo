// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyGetterTree.java
// @ closure-compiler v20260415

package trees

// ComputedPropertyGetterTree represents a computed property getter in an object literal or class.
type ComputedPropertyGetterTree struct {
	ParseTree
	Property *ParseTree
	IsStatic bool
	Body     *BlockTree
}

func NewComputedPropertyGetterTree(
	location SourceRange,
	property *ParseTree,
	isStatic bool,
	body *BlockTree,
) *ComputedPropertyGetterTree {
	return &ComputedPropertyGetterTree{
		ParseTree: NewParseTree(ParseTreeTypeComputedPropertyGetter, location),
		Property:  property,
		IsStatic:  isStatic,
		Body:      body,
	}
}
