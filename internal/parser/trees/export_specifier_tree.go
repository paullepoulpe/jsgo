// Port of com/google/javascript/jscomp/parsing/parser/trees/ExportSpecifierTree.java
// @ closure-compiler v20260415

package trees

type ExportSpecifierTree struct {
	ParseTree
	ImportedName    *IdentifierToken
	DestinationName *IdentifierToken
}

func NewExportSpecifierTree(location SourceRange, importedName *IdentifierToken, destinationName *IdentifierToken) *ExportSpecifierTree {
	return &ExportSpecifierTree{
		ParseTree:       ParseTree{Type: ParseTreeTypeExportSpecifier, Location: location},
		ImportedName:    importedName,
		DestinationName: destinationName,
	}
}
