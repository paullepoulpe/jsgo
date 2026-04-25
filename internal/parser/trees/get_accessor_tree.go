// Port of com/google/javascript/jscomp/parsing/parser/trees/GetAccessorTree.java
// @ closure-compiler v20260415

package trees

// GetAccessorTree represents a getter accessor in an object literal or class,
// e.g. get foo() { ... }.
type GetAccessorTree struct {
	ParseTree
	PropertyName  *Token
	IsClassMember bool
	IsStatic      bool
	Body          *BlockTree
}

// NewGetAccessorTree constructs a GetAccessorTree.
// Go-only: language-adapter
func NewGetAccessorTree(location SourceRange, propertyName *Token, isClassMember bool, isStatic bool, body *BlockTree) *GetAccessorTree {
	return &GetAccessorTree{
		ParseTree:     ParseTree{Type: ParseTreeTypeGetAccessor, Location: location},
		PropertyName:  propertyName,
		IsClassMember: isClassMember,
		IsStatic:      isStatic,
		Body:          body,
	}
}
