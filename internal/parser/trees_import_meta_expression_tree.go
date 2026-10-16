// Port of com/google/javascript/jscomp/parsing/parser/trees/ImportMetaExpressionTree.java @ closure-compiler v20260415

package parser

// Represents 'new.target'.
type ImportMetaExpressionTree struct {
	ParseTreeBase
}

// NewImportMetaExpressionTree constructs a ImportMetaExpressionTree.
//
// Java: ImportMetaExpressionTree(SourceRange)
func NewImportMetaExpressionTree(location *SourceRange) *ImportMetaExpressionTree {
	t := &ImportMetaExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_IMPORT_META_EXPRESSION, location)
	return t
}
