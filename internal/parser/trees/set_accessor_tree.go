// Port of com/google/javascript/jscomp/parsing/parser/trees/SetAccessorTree.java
// @ closure-compiler v20260415

package trees

// SetAccessorTree represents a setter property in an object literal or class.
type SetAccessorTree struct {
	ParseTree
	PropertyName  *Token
	Parameter     *FormalParameterListTree
	IsClassMember bool
	IsStatic      bool
	Body          *BlockTree
}

func NewSetAccessorTree(
	location SourceRange,
	propertyName *Token,
	isClassMember bool,
	isStatic bool,
	parameter *FormalParameterListTree,
	body *BlockTree,
) *SetAccessorTree {
	return &SetAccessorTree{
		ParseTree:     NewParseTree(ParseTreeTypeSetAccessor, location),
		PropertyName:  propertyName,
		IsClassMember: isClassMember,
		IsStatic:      isStatic,
		Parameter:     parameter,
		Body:          body,
	}
}
