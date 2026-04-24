// Port of com/google/javascript/jscomp/parsing/parser/trees/ExportDeclarationTree.java
// @ closure-compiler v20260415

package trees

type ExportDeclarationTree struct {
	// VariableStatement, Declaration (let, const, function or class),
	// or default with assignmentExpression
	IsDefault          bool
	IsExportAll        bool
	Declaration        ParseTree
	ExportSpecifierList []ParseTree
	From               *LiteralToken
}

func NewExportDeclarationTree(location SourceRange, isDefault bool, isExportAll bool, declaration ParseTree, exportSpecifierList []ParseTree, from *LiteralToken) *ExportDeclarationTree {
	return &ExportDeclarationTree{
		IsDefault:          isDefault,
		IsExportAll:        isExportAll,
		Declaration:        declaration,
		ExportSpecifierList: exportSpecifierList,
		From:               from,
	}
}
