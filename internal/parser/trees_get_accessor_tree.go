// Port of com/google/javascript/jscomp/parsing/parser/trees/GetAccessorTree.java @ closure-compiler v20260415

package parser

// GetAccessorTree is the parse tree for ParseTreeType GET_ACCESSOR.
type GetAccessorTree struct {
	ParseTreeBase

	PropertyName  Token
	IsClassMember bool
	IsStatic      bool
	Body          *BlockTree
}

// NewGetAccessorTree constructs a GetAccessorTree.
//
// Java: GetAccessorTree(SourceRange, Token, boolean, boolean, BlockTree)
func NewGetAccessorTree(location *SourceRange, propertyName Token, isClassMember bool, isStatic bool, body *BlockTree) *GetAccessorTree {
	t := &GetAccessorTree{}
	t.InitParseTreeBase(t, ParseTreeType_GET_ACCESSOR, location)
	t.PropertyName = propertyName
	t.IsClassMember = isClassMember
	t.IsStatic = isStatic
	t.Body = body
	return t
}
