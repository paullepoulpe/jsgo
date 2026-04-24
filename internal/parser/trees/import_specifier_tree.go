// Port of com/google/javascript/jscomp/parsing/parser/trees/ImportSpecifierTree.java
// @ closure-compiler v20260415

package trees

// ImportSpecifierTree represents a single specifier in an import clause.
type ImportSpecifierTree struct {
	ParseTree
	ImportedName    *IdentifierToken
	DestinationName *IdentifierToken
}

func NewImportSpecifierTree(
	location SourceRange,
	importedName *IdentifierToken,
	destinationName *IdentifierToken,
) *ImportSpecifierTree {
	return &ImportSpecifierTree{
		ParseTree:       NewParseTree(ParseTreeTypeImportSpecifier, location),
		ImportedName:    importedName,
		DestinationName: destinationName,
	}
}
