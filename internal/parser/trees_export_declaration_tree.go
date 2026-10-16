// Port of com/google/javascript/jscomp/parsing/parser/trees/ExportDeclarationTree.java @ closure-compiler v20260415

package parser

// ExportDeclarationTree is the parse tree for ParseTreeType EXPORT_DECLARATION.
type ExportDeclarationTree struct {
	ParseTreeBase

	// VariableStatement, Declaration (let, const, function or class),
	// or default with assignmentExpression
	IsDefault           bool
	IsExportAll         bool
	Declaration         ParseTree
	ExportSpecifierList []ParseTree
	From                LiteralToken
}

// NewExportDeclarationTree constructs a ExportDeclarationTree.
//
// Java: ExportDeclarationTree(SourceRange, boolean, boolean, ParseTree, ImmutableList, LiteralToken)
func NewExportDeclarationTree(location *SourceRange, isDefault bool, isExportAll bool, declaration ParseTree, exportSpecifierList []ParseTree, from LiteralToken) *ExportDeclarationTree {
	t := &ExportDeclarationTree{}
	t.InitParseTreeBase(t, ParseTreeType_EXPORT_DECLARATION, location)
	t.IsDefault = isDefault
	t.IsExportAll = isExportAll
	t.Declaration = declaration
	t.ExportSpecifierList = exportSpecifierList
	t.From = from
	return t
}
