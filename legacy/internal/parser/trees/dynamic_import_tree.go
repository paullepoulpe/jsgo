// Port of com/google/javascript/jscomp/parsing/parser/trees/DynamicImportTree.java
// @ closure-compiler v20260415

package trees

// DynamicImportTree represents a dynamic import expression,
// e.g. import(myModuleName) in const modulePromise = import(myModuleName).
type DynamicImportTree struct {
	ParseTree
	Argument *ParseTree
}

// NewDynamicImportTree constructs a DynamicImportTree.
// Go-only: language-adapter
func NewDynamicImportTree(location SourceRange, argument *ParseTree) *DynamicImportTree {
	return &DynamicImportTree{
		ParseTree: ParseTree{Type: ParseTreeTypeDynamicImportExpression, Location: location},
		Argument:  argument,
	}
}
