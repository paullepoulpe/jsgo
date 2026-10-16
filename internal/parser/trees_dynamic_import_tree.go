// Port of com/google/javascript/jscomp/parsing/parser/trees/DynamicImportTree.java @ closure-compiler v20260415

package parser

// Represents a dynamic import expression.
//
// e.g. `import(myModuleName)` in `const modulePromise = import(myModuleName)`.
type DynamicImportTree struct {
	ParseTreeBase

	Argument ParseTree
}

// NewDynamicImportTree constructs a DynamicImportTree.
//
// Java: DynamicImportTree(SourceRange, ParseTree)
func NewDynamicImportTree(location *SourceRange, argument ParseTree) *DynamicImportTree {
	t := &DynamicImportTree{}
	t.InitParseTreeBase(t, ParseTreeType_DYNAMIC_IMPORT_EXPRESSION, location)
	t.Argument = argument
	return t
}
