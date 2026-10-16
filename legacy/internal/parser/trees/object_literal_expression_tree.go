// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectLiteralExpressionTree.java
// @ closure-compiler v20260415

package trees

// ObjectLiteralExpressionTree represents an object literal expression,
// e.g. { a: 1, b: 2 }.
type ObjectLiteralExpressionTree struct {
	ParseTree
	PropertyNameAndValues []*ParseTree
	HasTrailingComma      bool
}

// NewObjectLiteralExpressionTree constructs an ObjectLiteralExpressionTree.
// Go-only: language-adapter
func NewObjectLiteralExpressionTree(location SourceRange, propertyNameAndValues []*ParseTree, hasTrailingComma bool) *ObjectLiteralExpressionTree {
	return &ObjectLiteralExpressionTree{
		ParseTree:             ParseTree{Type: ParseTreeTypeObjectLiteralExpression, Location: location},
		PropertyNameAndValues: propertyNameAndValues,
		HasTrailingComma:      hasTrailingComma,
	}
}
