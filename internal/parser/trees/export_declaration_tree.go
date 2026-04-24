// Port of com/google/javascript/jscomp/parsing/parser/trees/ExportDeclarationTree.java
// @ closure-compiler v20260415

package trees

// ExportDeclarationTree represents an export declaration.
type ExportDeclarationTree struct {
	ParseTree
	IsDefault           bool
	IsExportAll         bool
	// Declaration is a VariableStatement, class/function Declaration, or default
	// assignmentExpression.
	Declaration         *ParseTree
	ExportSpecifierList []*ParseTree
	// From is nil when there is no module specifier.
	From                *LiteralToken
}

func NewExportDeclarationTree(
	location SourceRange,
	isDefault bool,
	isExportAll bool,
	declaration *ParseTree,
	exportSpecifierList []*ParseTree,
	from *LiteralToken,
) *ExportDeclarationTree {
	return &ExportDeclarationTree{
		ParseTree:           NewParseTree(ParseTreeTypeExportDeclaration, location),
		IsDefault:           isDefault,
		IsExportAll:         isExportAll,
		Declaration:         declaration,
		ExportSpecifierList: exportSpecifierList,
		From:                from,
	}
}
