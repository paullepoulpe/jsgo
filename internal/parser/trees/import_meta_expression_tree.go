// Port of com/google/javascript/jscomp/parsing/parser/trees/ImportMetaExpressionTree.java
// @ closure-compiler v20260415

package trees

// ImportMetaExpressionTree represents the import.meta expression.
type ImportMetaExpressionTree struct {
	ParseTree
}

func NewImportMetaExpressionTree(location SourceRange) *ImportMetaExpressionTree {
	return &ImportMetaExpressionTree{
		ParseTree: NewParseTree(ParseTreeTypeImportMetaExpression, location),
	}
}
