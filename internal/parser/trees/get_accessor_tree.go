// Port of com/google/javascript/jscomp/parsing/parser/trees/GetAccessorTree.java
// @ closure-compiler v20260415

package trees

// GetAccessorTree represents a getter property in an object literal or class.
type GetAccessorTree struct {
	ParseTree
	PropertyName  *Token
	IsClassMember bool
	IsStatic      bool
	Body          *BlockTree
}

func NewGetAccessorTree(
	location SourceRange,
	propertyName *Token,
	isClassMember bool,
	isStatic bool,
	body *BlockTree,
) *GetAccessorTree {
	return &GetAccessorTree{
		ParseTree:     NewParseTree(ParseTreeTypeGetAccessor, location),
		PropertyName:  propertyName,
		IsClassMember: isClassMember,
		IsStatic:      isStatic,
		Body:          body,
	}
}
