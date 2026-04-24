// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyGetterTree.java
// @ closure-compiler v20260415

package trees

// ComputedPropertyGetterTree represents a computed getter in an object literal or class,
// e.g. { get [expr]() {} }.
type ComputedPropertyGetterTree struct {
	ParseTree
	Property  *ParseTree
	IsStatic  bool
	Body      *BlockTree
}

// NewComputedPropertyGetterTree constructs a ComputedPropertyGetterTree.
func NewComputedPropertyGetterTree(location SourceRange, property *ParseTree, isStatic bool, body *BlockTree) *ComputedPropertyGetterTree {
	return &ComputedPropertyGetterTree{
		ParseTree: ParseTree{Type: ParseTreeTypeComputedPropertyGetter, Location: location},
		Property:  property,
		IsStatic:  isStatic,
		Body:      body,
	}
}
