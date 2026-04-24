// Port of com/google/javascript/jscomp/parsing/parser/trees/SetAccessorTree.java
// @ closure-compiler v20260415

package trees

// SetAccessorTree represents a setter accessor in an object literal or class,
// e.g. set foo(v) { ... }.
type SetAccessorTree struct {
	ParseTree
	PropertyName  *Token
	Parameter     *FormalParameterListTree
	IsClassMember bool
	IsStatic      bool
	Body          *BlockTree
}

// NewSetAccessorTree constructs a SetAccessorTree.
// Go-only: language-adapter
func NewSetAccessorTree(location SourceRange, propertyName *Token, isClassMember bool, isStatic bool, parameter *FormalParameterListTree, body *BlockTree) *SetAccessorTree {
	return &SetAccessorTree{
		ParseTree:     ParseTree{Type: ParseTreeTypeSetAccessor, Location: location},
		PropertyName:  propertyName,
		IsClassMember: isClassMember,
		IsStatic:      isStatic,
		Parameter:     parameter,
		Body:          body,
	}
}
