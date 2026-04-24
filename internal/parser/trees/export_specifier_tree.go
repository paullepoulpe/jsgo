// Port of com/google/javascript/jscomp/parsing/parser/trees/ExportSpecifierTree.java
// @ closure-compiler v20260415

package trees

// ExportSpecifierTree represents a single specifier in an export clause.
type ExportSpecifierTree struct {
	ParseTree
	ImportedName    *IdentifierToken
	DestinationName *IdentifierToken
}

func NewExportSpecifierTree(
	location SourceRange,
	importedName *IdentifierToken,
	destinationName *IdentifierToken,
) *ExportSpecifierTree {
	return &ExportSpecifierTree{
		ParseTree:       NewParseTree(ParseTreeTypeExportSpecifier, location),
		ImportedName:    importedName,
		DestinationName: destinationName,
	}
}
