// Port of com/google/javascript/jscomp/parsing/parser/trees/ExportSpecifierTree.java @ closure-compiler v20260415

package parser

// ExportSpecifierTree is the parse tree for ParseTreeType EXPORT_SPECIFIER.
type ExportSpecifierTree struct {
	ParseTreeBase

	ImportedName    *IdentifierToken
	DestinationName *IdentifierToken
}

// NewExportSpecifierTree constructs a ExportSpecifierTree.
//
// Java: ExportSpecifierTree(SourceRange, IdentifierToken, IdentifierToken)
func NewExportSpecifierTree(location *SourceRange, importedName *IdentifierToken, destinationName *IdentifierToken) *ExportSpecifierTree {
	t := &ExportSpecifierTree{}
	t.InitParseTreeBase(t, ParseTreeType_EXPORT_SPECIFIER, location)
	t.ImportedName = importedName
	t.DestinationName = destinationName
	return t
}
