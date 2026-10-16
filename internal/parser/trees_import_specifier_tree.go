// Port of com/google/javascript/jscomp/parsing/parser/trees/ImportSpecifierTree.java @ closure-compiler v20260415

package parser

// ImportSpecifierTree is the parse tree for ParseTreeType IMPORT_SPECIFIER.
type ImportSpecifierTree struct {
	ParseTreeBase

	ImportedName    *IdentifierToken
	DestinationName *IdentifierToken
}

// NewImportSpecifierTree constructs a ImportSpecifierTree.
//
// Java: ImportSpecifierTree(SourceRange, IdentifierToken, IdentifierToken)
func NewImportSpecifierTree(location *SourceRange, importedName *IdentifierToken, destinationName *IdentifierToken) *ImportSpecifierTree {
	t := &ImportSpecifierTree{}
	t.InitParseTreeBase(t, ParseTreeType_IMPORT_SPECIFIER, location)
	t.ImportedName = importedName
	t.DestinationName = destinationName
	return t
}
