// Port of com/google/javascript/jscomp/parsing/parser/trees/ImportDeclarationTree.java
// @ closure-compiler v20260415

package trees

type ImportDeclarationTree struct {
	// The identifier for the "default" import from the module.
	DefaultBindingIdentifier *IdentifierToken
	// The list of names imported from the module when using the
	// `import {...} from 'specifier'` form.
	ImportSpecifierList []ParseTree
	// The namespace into which imported names are put when using the
	// `import * as nameSpaceImportIdentifier from 'specifier'` form.
	NameSpaceImportIdentifier *IdentifierToken
	// The string identifying the module.
	ModuleSpecifier *LiteralToken
}

func NewImportDeclarationTree(location SourceRange, defaultBindingIdentifier *IdentifierToken, importSpecifierList []ParseTree, nameSpaceImportIdentifier *IdentifierToken, moduleSpecifier *LiteralToken) *ImportDeclarationTree {
	return &ImportDeclarationTree{
		DefaultBindingIdentifier:  defaultBindingIdentifier,
		ImportSpecifierList:       importSpecifierList,
		NameSpaceImportIdentifier: nameSpaceImportIdentifier,
		ModuleSpecifier:           moduleSpecifier,
	}
}
