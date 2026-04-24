// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectLiteralExpressionTree.java
// @ closure-compiler v20260415

package trees

// ObjectLiteralExpressionTree represents an object literal expression.
type ObjectLiteralExpressionTree struct {
	ParseTree
	PropertyNameAndValues []*ParseTree
	HasTrailingComma      bool
}

func NewObjectLiteralExpressionTree(
	location SourceRange,
	propertyNameAndValues []*ParseTree,
	hasTrailingComma bool,
) *ObjectLiteralExpressionTree {
	return &ObjectLiteralExpressionTree{
		ParseTree:             NewParseTree(ParseTreeTypeObjectLiteralExpression, location),
		PropertyNameAndValues: propertyNameAndValues,
		HasTrailingComma:      hasTrailingComma,
	}
}
