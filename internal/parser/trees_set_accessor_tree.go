// Port of com/google/javascript/jscomp/parsing/parser/trees/SetAccessorTree.java @ closure-compiler v20260415

package parser

// SetAccessorTree is the parse tree for ParseTreeType SET_ACCESSOR.
type SetAccessorTree struct {
	ParseTreeBase

	PropertyName  Token
	Parameter     *FormalParameterListTree
	IsClassMember bool
	IsStatic      bool
	Body          *BlockTree
}

// NewSetAccessorTree constructs a SetAccessorTree.
//
// Java: SetAccessorTree(SourceRange, Token, boolean, boolean, FormalParameterListTree, BlockTree)
func NewSetAccessorTree(location *SourceRange, propertyName Token, isClassMember bool, isStatic bool, parameter *FormalParameterListTree, body *BlockTree) *SetAccessorTree {
	t := &SetAccessorTree{}
	t.InitParseTreeBase(t, ParseTreeType_SET_ACCESSOR, location)
	t.PropertyName = propertyName
	t.IsClassMember = isClassMember
	t.IsStatic = isStatic
	t.Parameter = parameter
	t.Body = body
	return t
}
