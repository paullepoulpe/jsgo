// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableDeclarationTree.java @ closure-compiler v20260415

package parser

// VariableDeclarationTree is the parse tree for ParseTreeType VARIABLE_DECLARATION.
type VariableDeclarationTree struct {
	ParseTreeBase

	// An IdentifierExpressionTree, ArrayPatternTree or ObjectLiteralTree
	Lvalue      ParseTree
	Initializer ParseTree
}

// NewVariableDeclarationTree constructs a VariableDeclarationTree.
//
// Java: VariableDeclarationTree(SourceRange, ParseTree, ParseTree)
func NewVariableDeclarationTree(location *SourceRange, lvalue ParseTree, initializer ParseTree) *VariableDeclarationTree {
	t := &VariableDeclarationTree{}
	t.InitParseTreeBase(t, ParseTreeType_VARIABLE_DECLARATION, location)
	t.Lvalue = lvalue
	t.Initializer = initializer
	return t
}
