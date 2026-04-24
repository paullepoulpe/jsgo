// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyMethodTree.java
// @ closure-compiler v20260415

package trees

// ComputedPropertyMethodTree represents a computed property method in an object literal or class.
type ComputedPropertyMethodTree struct {
	ParseTree
	Property *ParseTree
	Method   *ParseTree
}

func NewComputedPropertyMethodTree(
	location SourceRange,
	property *ParseTree,
	method *ParseTree,
) *ComputedPropertyMethodTree {
	return &ComputedPropertyMethodTree{
		ParseTree: NewParseTree(ParseTreeTypeComputedPropertyMethod, location),
		Property:  property,
		Method:    method,
	}
}
