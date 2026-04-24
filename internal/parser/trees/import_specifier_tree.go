// Port of com/google/javascript/jscomp/parsing/parser/trees/ImportSpecifierTree.java
// @ closure-compiler v20260415

package trees

type ImportSpecifierTree struct {
	ParseTree
	ImportedName    *IdentifierToken
	DestinationName *IdentifierToken
}

func NewImportSpecifierTree(location SourceRange, importedName *IdentifierToken, destinationName *IdentifierToken) *ImportSpecifierTree {
	return &ImportSpecifierTree{
		ParseTree:       ParseTree{Type: ParseTreeTypeImportSpecifier, Location: location},
		ImportedName:    importedName,
		DestinationName: destinationName,
	}
}
