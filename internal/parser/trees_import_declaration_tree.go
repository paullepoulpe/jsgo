// Port of com/google/javascript/jscomp/parsing/parser/trees/ImportDeclarationTree.java @ closure-compiler v20260415

package parser

// ImportDeclarationTree is the parse tree for ParseTreeType IMPORT_DECLARATION.
type ImportDeclarationTree struct {
	ParseTreeBase

	// The identifier for the "default" import from the module.
	DefaultBindingIdentifier *IdentifierToken
	// The list of names imported from the module when using the
	// `import {...} from 'specifier'` form.
	ImportSpecifierList []ParseTree
	// The namespace into which imported names are put when using the
	// `import * as nameSpaceImportIdentifier from 'specifier'` form.
	NameSpaceImportIdentifier *IdentifierToken
	// The string identifying the module.
	ModuleSpecifier LiteralToken
}

// NewImportDeclarationTree constructs a ImportDeclarationTree.
//
// Java: ImportDeclarationTree(SourceRange, IdentifierToken, ImmutableList, IdentifierToken, LiteralToken)
func NewImportDeclarationTree(location *SourceRange, defaultBindingIdentifier *IdentifierToken, importSpecifierList []ParseTree, nameSpaceImportIdentifier *IdentifierToken, moduleSpecifier LiteralToken) *ImportDeclarationTree {
	t := &ImportDeclarationTree{}
	t.InitParseTreeBase(t, ParseTreeType_IMPORT_DECLARATION, location)
	t.DefaultBindingIdentifier = defaultBindingIdentifier
	t.ImportSpecifierList = importSpecifierList
	t.NameSpaceImportIdentifier = nameSpaceImportIdentifier
	t.ModuleSpecifier = moduleSpecifier
	return t
}
