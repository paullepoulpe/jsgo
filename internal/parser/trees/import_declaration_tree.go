// Port of com/google/javascript/jscomp/parsing/parser/trees/ImportDeclarationTree.java
// @ closure-compiler v20260415

package trees

// ImportDeclarationTree represents an import declaration.
type ImportDeclarationTree struct {
	ParseTree
	// DefaultBindingIdentifier is the identifier for the "default" import; nil when absent.
	DefaultBindingIdentifier *IdentifierToken
	// ImportSpecifierList holds the named imports; nil/empty when absent.
	ImportSpecifierList []*ParseTree
	// NameSpaceImportIdentifier is the namespace import identifier (* as foo); nil when absent.
	NameSpaceImportIdentifier *IdentifierToken
	// ModuleSpecifier is the string module specifier.
	ModuleSpecifier *LiteralToken
}

func NewImportDeclarationTree(
	location SourceRange,
	defaultBindingIdentifier *IdentifierToken,
	importSpecifierList []*ParseTree,
	nameSpaceImportIdentifier *IdentifierToken,
	moduleSpecifier *LiteralToken,
) *ImportDeclarationTree {
	return &ImportDeclarationTree{
		ParseTree:                 NewParseTree(ParseTreeTypeImportDeclaration, location),
		DefaultBindingIdentifier:  defaultBindingIdentifier,
		ImportSpecifierList:       importSpecifierList,
		NameSpaceImportIdentifier: nameSpaceImportIdentifier,
		ModuleSpecifier:           moduleSpecifier,
	}
}
