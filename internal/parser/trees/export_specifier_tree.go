// Port of com/google/javascript/jscomp/parsing/parser/trees/ExportSpecifierTree.java
// @ closure-compiler v20260415

package trees

type ExportSpecifierTree struct {
	ImportedName    *IdentifierToken
	DestinationName *IdentifierToken
}

func NewExportSpecifierTree(location SourceRange, importedName *IdentifierToken, destinationName *IdentifierToken) *ExportSpecifierTree {
	return &ExportSpecifierTree{
		ImportedName:    importedName,
		DestinationName: destinationName,
	}
}
