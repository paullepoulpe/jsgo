// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectLiteralExpressionTree.java @ closure-compiler v20260415

package parser

// ObjectLiteralExpressionTree is the parse tree for ParseTreeType OBJECT_LITERAL_EXPRESSION.
type ObjectLiteralExpressionTree struct {
	ParseTreeBase

	PropertyNameAndValues []ParseTree
	HasTrailingComma      bool
}

// NewObjectLiteralExpressionTree constructs a ObjectLiteralExpressionTree.
//
// Java: ObjectLiteralExpressionTree(SourceRange, ImmutableList, boolean)
func NewObjectLiteralExpressionTree(location *SourceRange, propertyNameAndValues []ParseTree, hasTrailingComma bool) *ObjectLiteralExpressionTree {
	t := &ObjectLiteralExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_OBJECT_LITERAL_EXPRESSION, location)
	t.PropertyNameAndValues = propertyNameAndValues
	t.HasTrailingComma = hasTrailingComma
	return t
}
