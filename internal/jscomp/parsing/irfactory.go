// Port of com/google/javascript/jscomp/parsing/IRFactory.java
// @ closure-compiler v20260415
//
// IRFactory transforms ParseTree nodes (from internal/parser/trees) into
// rhino IR Nodes (from internal/rhino). Only the subset of transforms
// needed for the current corpus tier is implemented; stub methods are
// marked with a TODO comment.

package parsing

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/parser/trees"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// IRFactory transforms the external parse-tree AST (trees.ParseTree) into
// the internal IR AST (rhino.Node). Mirrors
// com.google.javascript.jscomp.parsing.IRFactory.
type IRFactory struct {
	sourceFile    rhino.StaticSourceFile
	errorReporter errorReporter
}

// errorReporter is a minimal error-reporting interface used by IRFactory.
// Mirrors com.google.javascript.rhino.ErrorReporter.
type errorReporter interface {
	// Error reports a hard error at the given source position.
	Error(msg, sourceName string, lineno, charno int)
	// Warning reports a warning at the given source position.
	Warning(msg, sourceName string, lineno, charno int)
}

// NewIRFactory constructs an IRFactory for the given source file.
func NewIRFactory(sourceFile rhino.StaticSourceFile, er errorReporter) *IRFactory {
	return &IRFactory{
		sourceFile:    sourceFile,
		errorReporter: er,
	}
}

// Transform converts a ParseTree node into a rhino IR Node.
// Mirrors IRFactory.transform(ParseTree) in Java.
func (f *IRFactory) Transform(tree *trees.ParseTree) *rhino.Node {
	if tree == nil {
		return nil
	}
	node := f.dispatch(tree)
	if node != nil {
		f.setSourceInfo(node, tree)
	}
	return node
}

// dispatch routes a ParseTree to the appropriate transform* method based on
// its type. Mirrors the TransformDispatcher.process(ParseTree) switch in Java.
func (f *IRFactory) dispatch(tree *trees.ParseTree) *rhino.Node {
	switch tree.Type {
	// ---- programme / block -----------------------------------------
	case trees.ParseTreeTypeProgram:
		return f.transformProgram((*trees.ProgramTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeBlock:
		return f.transformBlock((*trees.BlockTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeExpressionStatement:
		return f.transformExpressionStatement((*trees.ExpressionStatementTree)(unsafe.Pointer(tree)))

	// ---- basic statements ------------------------------------------
	case trees.ParseTreeTypeReturnStatement:
		return f.transformReturnStatement((*trees.ReturnStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeThrowStatement:
		return f.transformThrowStatement((*trees.ThrowStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeBreakStatement:
		return f.transformBreakStatement((*trees.BreakStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeContinueStatement:
		return f.transformContinueStatement((*trees.ContinueStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeIfStatement:
		return f.transformIfStatement((*trees.IfStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeWhileStatement:
		return f.transformWhileStatement((*trees.WhileStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeDoWhileStatement:
		return f.transformDoWhileStatement((*trees.DoWhileStatementTree)(unsafe.Pointer(tree)))

	// ---- loops + try/switch ----------------------------------------
	case trees.ParseTreeTypeForStatement:
		return f.transformForStatement((*trees.ForStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeForInStatement:
		return f.transformForInStatement((*trees.ForInStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeForOfStatement:
		return f.transformForOfStatement((*trees.ForOfStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeForAwaitOfStatement:
		return f.transformForAwaitOfStatement((*trees.ForAwaitOfStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeSwitchStatement:
		return f.transformSwitchStatement((*trees.SwitchStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeCaseClause:
		return f.transformCaseClause((*trees.CaseClauseTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeDefaultClause:
		return f.transformDefaultClause((*trees.DefaultClauseTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeTryStatement:
		return f.transformTryStatement((*trees.TryStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeCatch:
		return f.transformCatch((*trees.CatchTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeFinally:
		return f.transformFinally((*trees.FinallyTree)(unsafe.Pointer(tree)))

	// ---- declarations + functions + classes ------------------------
	case trees.ParseTreeTypeFunctionDeclaration:
		return f.transformFunctionDeclaration((*trees.FunctionDeclarationTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeFormalParameterList:
		return f.transformFormalParameterList((*trees.FormalParameterListTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeDefaultParameter:
		return f.transformDefaultParameter((*trees.DefaultParameterTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeClassDeclaration:
		return f.transformClassDeclaration((*trees.ClassDeclarationTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeFieldDeclaration:
		return f.transformFieldDeclaration((*trees.FieldDeclarationTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeVariableStatement:
		return f.transformVariableStatement((*trees.VariableStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeVariableDeclarationList:
		return f.transformVariableDeclarationList((*trees.VariableDeclarationListTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeVariableDeclaration:
		return f.transformVariableDeclaration((*trees.VariableDeclarationTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeImportDeclaration:
		return f.transformImportDeclaration((*trees.ImportDeclarationTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeImportSpecifier:
		return f.transformImportSpecifier((*trees.ImportSpecifierTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeExportDeclaration:
		return f.transformExportDeclaration((*trees.ExportDeclarationTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeExportSpecifier:
		return f.transformExportSpecifier((*trees.ExportSpecifierTree)(unsafe.Pointer(tree)))

	// ---- primitive expressions -------------------------------------
	case trees.ParseTreeTypeIdentifierExpression:
		return f.transformIdentifierExpression((*trees.IdentifierExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeLiteralExpression:
		return f.transformLiteralExpression((*trees.LiteralExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeThisExpression:
		return f.transformThisExpression()
	case trees.ParseTreeTypeSuperExpression:
		return f.transformSuperExpression()

	// ---- unary / binary / update -----------------------------------
	case trees.ParseTreeTypeUnaryExpression:
		return f.transformUnaryExpression((*trees.UnaryExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeUpdateExpression:
		return f.transformUpdateExpression((*trees.UpdateExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeBinaryOperator:
		return f.transformBinaryOperator((*trees.BinaryOperatorTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeConditionalExpression:
		return f.transformConditionalExpression((*trees.ConditionalExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeCommaExpression:
		return f.transformCommaExpression((*trees.CommaExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeParenExpression:
		return f.transformParenExpression((*trees.ParenExpressionTree)(unsafe.Pointer(tree)))

	// ---- calls / member access / new -------------------------------
	case trees.ParseTreeTypeCallExpression:
		return f.transformCallExpression((*trees.CallExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeNewExpression:
		return f.transformNewExpression((*trees.NewExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeMemberExpression:
		return f.transformMemberExpression((*trees.MemberExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeMemberLookupExpression:
		return f.transformMemberLookupExpression((*trees.MemberLookupExpressionTree)(unsafe.Pointer(tree)))

	// ---- structured literals ---------------------------------------
	case trees.ParseTreeTypeArrayLiteralExpression:
		return f.transformArrayLiteralExpression((*trees.ArrayLiteralExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeObjectLiteralExpression:
		return f.transformObjectLiteralExpression((*trees.ObjectLiteralExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypePropertyNameAssignment:
		return f.transformPropertyNameAssignment((*trees.PropertyNameAssignmentTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeGetAccessor:
		return f.transformGetAccessor((*trees.GetAccessorTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeSetAccessor:
		return f.transformSetAccessor((*trees.SetAccessorTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeTemplateLiteralExpression:
		return f.transformTemplateLiteralExpression((*trees.TemplateLiteralExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeTemplateLiteralPortion:
		return f.transformTemplateLiteralPortion((*trees.TemplateLiteralPortionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeTemplateSubstitution:
		return f.transformTemplateSubstitution((*trees.TemplateSubstitutionTree)(unsafe.Pointer(tree)))

	// ---- null / empty placeholder ----------------------------------
	case trees.ParseTreeTypeNull:
		return rhino.New(rhino.EMPTY)

	default:
		return rhino.New(rhino.EMPTY)
	}
}

// setSourceInfo copies source-position information from tree to node.
// Mirrors IRFactory.setSourceInfo(Node, ParseTree) in Java.
func (f *IRFactory) setSourceInfo(node *rhino.Node, tree *trees.ParseTree) {
	if node.Lineno() == 0 {
		// Location lines in ParseTree are 0-based; rhino expects 1-based.
		node.SetPos(tree.Location.Start.Line+1, tree.Location.Start.Column)
	}
}

// newNode creates a new rhino Node with the given token. Mirrors
// IRFactory.newNode(Token) in Java.
func (f *IRFactory) newNode(tok rhino.Token) *rhino.Node {
	return rhino.New(tok)
}

// newStringNode creates a new rhino string Node. Mirrors
// IRFactory.newStringNode(Token, String) in Java.
func (f *IRFactory) newStringNode(tok rhino.Token, s string) *rhino.Node {
	return rhino.NewString(tok, s)
}

// newNumberNode creates a new rhino number Node. Mirrors
// IRFactory.newNumberNode(double) in Java.
func (f *IRFactory) newNumberNode(v float64) *rhino.Node {
	return rhino.NewNumber(v)
}

// ---- programme / block ---------------------------------------------

// transformProgram mirrors TransformDispatcher.processAstRoot(ProgramTree).
func (f *IRFactory) transformProgram(tree *trees.ProgramTree) *rhino.Node {
	script := f.newNode(rhino.SCRIPT)
	for _, child := range tree.SourceElements {
		script.AddChildToBack(f.Transform(child))
	}
	return script
}

// transformBlock mirrors TransformDispatcher.processBlock(BlockTree).
func (f *IRFactory) transformBlock(tree *trees.BlockTree) *rhino.Node {
	node := f.newNode(rhino.BLOCK)
	for _, child := range tree.Statements {
		node.AddChildToBack(f.Transform(child))
	}
	return node
}

// transformBlockTree transforms any ParseTree as a block, matching Java's
// transformBlock(ParseTree) which dispatches and wraps non-block results.
// When the tree is already a BLOCK the result is used directly; when it is
// anything else (e.g. a single statement) it is wrapped in a synthetic BLOCK.
func (f *IRFactory) transformBlockTree(tree *trees.ParseTree) *rhino.Node {
	if tree == nil {
		return f.newNode(rhino.BLOCK)
	}
	result := f.Transform(tree)
	if result == nil {
		return f.newNode(rhino.BLOCK)
	}
	if result.Token() == rhino.BLOCK {
		return result
	}
	// Wrap a single statement in a synthetic block (mirrors Java behaviour for
	// braceless if/while/for bodies).
	block := f.newNode(rhino.BLOCK)
	block.PutBooleanProp(rhino.PropAddedBlock, true)
	f.setSourceInfo(block, tree)
	block.AddChildToBack(result)
	return block
}

// transformOrEmpty transforms tree if non-nil, else returns an EMPTY node
// positioned at parent. Mirrors IRFactory.TransformDispatcher.transformOrEmpty.
func (f *IRFactory) transformOrEmpty(tree *trees.ParseTree, parent *trees.ParseTree) *rhino.Node {
	if tree == nil {
		n := f.newNode(rhino.EMPTY)
		if parent != nil {
			f.setSourceInfo(n, parent)
		}
		return n
	}
	return f.Transform(tree)
}

// transformExpressionStatementExpr mirrors
// TransformDispatcher.processExpressionStatement, given the inner expression
// tree. Used when the expression-statement's child is dispatched directly.
func (f *IRFactory) transformExpressionStatementExpr(expr *trees.ParseTree) *rhino.Node {
	node := f.newNode(rhino.EXPR_RESULT)
	node.AddChildToBack(f.Transform(expr))
	return node
}

// transformExpressionStatement mirrors
// TransformDispatcher.processExpressionStatement(ExpressionStatementTree).
func (f *IRFactory) transformExpressionStatement(tree *trees.ExpressionStatementTree) *rhino.Node {
	node := f.newNode(rhino.EXPR_RESULT)
	node.AddChildToBack(f.Transform(tree.Expression))
	return node
}

// ---- basic statements ----------------------------------------------

// transformReturnStatement mirrors
// TransformDispatcher.processReturnStatement(ReturnStatementTree).
func (f *IRFactory) transformReturnStatement(tree *trees.ReturnStatementTree) *rhino.Node {
	node := f.newNode(rhino.RETURN)
	if tree.Expression != nil {
		node.AddChildToBack(f.Transform(tree.Expression))
	}
	return node
}

// transformThrowStatement mirrors
// TransformDispatcher.processThrowStatement(ThrowStatementTree).
func (f *IRFactory) transformThrowStatement(tree *trees.ThrowStatementTree) *rhino.Node {
	node := f.newNode(rhino.THROW)
	node.AddChildToBack(f.Transform(tree.Value))
	return node
}

// transformBreakStatement mirrors
// TransformDispatcher.processBreakStatement(BreakStatementTree).
func (f *IRFactory) transformBreakStatement(tree *trees.BreakStatementTree) *rhino.Node {
	node := f.newNode(rhino.BREAK)
	if tree.HasLabel() {
		labelName := f.transformLabelName(tree.Name)
		node.AddChildToBack(labelName)
	}
	return node
}

// transformContinueStatement mirrors
// TransformDispatcher.processContinueStatement(ContinueStatementTree).
func (f *IRFactory) transformContinueStatement(tree *trees.ContinueStatementTree) *rhino.Node {
	node := f.newNode(rhino.CONTINUE)
	if tree.HasLabel() {
		labelName := f.transformLabelName(tree.Name)
		node.AddChildToBack(labelName)
	}
	return node
}

// transformLabelName builds a LABEL_NAME node from an IdentifierToken.
// Mirrors IRFactory.TransformDispatcher.transformLabelName(IdentifierToken).
func (f *IRFactory) transformLabelName(tok *trees.IdentifierToken) *rhino.Node {
	label := f.newStringNode(rhino.LABEL_NAME, tok.Value)
	label.SetPos(tok.Location.Start.Line+1, tok.Location.Start.Column)
	return label
}

// transformIfStatement mirrors
// TransformDispatcher.processIfStatement(IfStatementTree).
func (f *IRFactory) transformIfStatement(tree *trees.IfStatementTree) *rhino.Node {
	node := f.newNode(rhino.IF)
	node.AddChildToBack(f.Transform(tree.Condition))
	node.AddChildToBack(f.transformBlockTree(tree.IfClause))
	if tree.ElseClause != nil {
		node.AddChildToBack(f.transformBlockTree(tree.ElseClause))
	}
	return node
}

// transformWhileStatement mirrors
// TransformDispatcher.processWhileLoop(WhileStatementTree).
func (f *IRFactory) transformWhileStatement(tree *trees.WhileStatementTree) *rhino.Node {
	node := f.newNode(rhino.WHILE)
	node.AddChildToBack(f.Transform(tree.Condition))
	node.AddChildToBack(f.transformBlockTree(tree.Body))
	return node
}

// transformDoWhileStatement mirrors
// TransformDispatcher.processDoLoop(DoWhileStatementTree).
func (f *IRFactory) transformDoWhileStatement(tree *trees.DoWhileStatementTree) *rhino.Node {
	node := f.newNode(rhino.DO)
	node.AddChildToBack(f.transformBlockTree(tree.Body))
	node.AddChildToBack(f.Transform(tree.Condition))
	return node
}

// ---- loops + try/switch -------------------------------------------

// transformForStatement mirrors
// TransformDispatcher.processForLoop(ForStatementTree).
func (f *IRFactory) transformForStatement(tree *trees.ForStatementTree) *rhino.Node {
	node := f.newNode(rhino.FOR)
	node.AddChildToBack(f.transformOrEmpty(tree.Initializer, &tree.ParseTree))
	node.AddChildToBack(f.transformOrEmpty(tree.Condition, &tree.ParseTree))
	node.AddChildToBack(f.transformOrEmpty(tree.Increment, &tree.ParseTree))
	node.AddChildToBack(f.transformBlockTree(tree.Body))
	return node
}

// transformForInStatement mirrors
// TransformDispatcher.processForInLoop(ForInStatementTree).
func (f *IRFactory) transformForInStatement(tree *trees.ForInStatementTree) *rhino.Node {
	node := f.newNode(rhino.FOR_IN)
	node.AddChildToBack(f.Transform(tree.Initializer))
	node.AddChildToBack(f.Transform(tree.Collection))
	node.AddChildToBack(f.transformBlockTree(tree.Body))
	return node
}

// transformForOfStatement mirrors
// TransformDispatcher.processForOf(ForOfStatementTree).
func (f *IRFactory) transformForOfStatement(tree *trees.ForOfStatementTree) *rhino.Node {
	node := f.newNode(rhino.FOR_OF)
	node.AddChildToBack(f.Transform(tree.Initializer))
	node.AddChildToBack(f.Transform(tree.Collection))
	node.AddChildToBack(f.transformBlockTree(tree.Body))
	return node
}

// transformForAwaitOfStatement mirrors
// TransformDispatcher.processForAwaitOf(ForAwaitOfStatementTree).
func (f *IRFactory) transformForAwaitOfStatement(tree *trees.ForAwaitOfStatementTree) *rhino.Node {
	node := f.newNode(rhino.FOR_AWAIT_OF)
	node.AddChildToBack(f.Transform(tree.Initializer))
	node.AddChildToBack(f.Transform(tree.Collection))
	node.AddChildToBack(f.transformBlockTree(tree.Body))
	return node
}

// transformSwitchStatement mirrors
// TransformDispatcher.processSwitchStatement(SwitchStatementTree).
func (f *IRFactory) transformSwitchStatement(tree *trees.SwitchStatementTree) *rhino.Node {
	node := f.newNode(rhino.SWITCH)
	node.AddChildToBack(f.Transform(tree.Expression))
	switchBody := f.newNode(rhino.SWITCH_BODY)
	f.setSourceInfo(switchBody, &tree.ParseTree)
	node.AddChildToBack(switchBody)
	for _, child := range tree.CaseClauses {
		switchBody.AddChildToBack(f.Transform(child))
	}
	return node
}

// transformCaseClause mirrors
// TransformDispatcher.processSwitchCase(CaseClauseTree).
func (f *IRFactory) transformCaseClause(tree *trees.CaseClauseTree) *rhino.Node {
	node := f.newNode(rhino.CASE)
	node.AddChildToBack(f.Transform(tree.Expression))
	block := f.newNode(rhino.BLOCK)
	block.PutBooleanProp(rhino.PropAddedBlock, true)
	f.setSourceInfo(block, &tree.ParseTree)
	for _, stmt := range tree.Statements {
		block.AddChildToBack(f.Transform(stmt))
	}
	node.AddChildToBack(block)
	return node
}

// transformDefaultClause mirrors
// TransformDispatcher.processSwitchDefault(DefaultClauseTree).
func (f *IRFactory) transformDefaultClause(tree *trees.DefaultClauseTree) *rhino.Node {
	node := f.newNode(rhino.DEFAULT_CASE)
	block := f.newNode(rhino.BLOCK)
	block.PutBooleanProp(rhino.PropAddedBlock, true)
	f.setSourceInfo(block, &tree.ParseTree)
	for _, stmt := range tree.Statements {
		block.AddChildToBack(f.Transform(stmt))
	}
	node.AddChildToBack(block)
	return node
}

// transformTryStatement mirrors
// TransformDispatcher.processTryStatement(TryStatementTree).
func (f *IRFactory) transformTryStatement(tree *trees.TryStatementTree) *rhino.Node {
	node := f.newNode(rhino.TRY)
	node.AddChildToBack(f.transformBlockTree(tree.Body))
	block := f.newNode(rhino.BLOCK)
	lineSet := false
	if tree.CatchBlock != nil {
		f.setSourceInfo(block, tree.CatchBlock)
		lineSet = true
		block.AddChildToBack(f.Transform(tree.CatchBlock))
	}
	node.AddChildToBack(block)
	if tree.FinallyBlock != nil {
		node.AddChildToBack(f.transformBlockTree(tree.FinallyBlock))
	}
	// Mirror Java: if no catch clause, set the block's position from finally.
	if !lineSet && tree.FinallyBlock != nil {
		f.setSourceInfo(block, tree.FinallyBlock)
	}
	return node
}

// transformCatch mirrors
// TransformDispatcher.processCatchClause(CatchTree).
func (f *IRFactory) transformCatch(tree *trees.CatchTree) *rhino.Node {
	node := f.newNode(rhino.CATCH)
	node.AddChildToBack(f.Transform(tree.Exception))
	node.AddChildToBack(f.transformBlockTree(tree.CatchBody))
	return node
}

// transformFinally mirrors
// TransformDispatcher.processFinally(FinallyTree).
// Java's processFinally calls transformBlock on the block field.
func (f *IRFactory) transformFinally(tree *trees.FinallyTree) *rhino.Node {
	return f.transformBlockTree(tree.Block)
}

// ---- declarations + functions + classes ---------------------------

// transformFunctionDeclaration mirrors
// TransformDispatcher.processFunction(FunctionDeclarationTree).
func (f *IRFactory) transformFunctionDeclaration(tree *trees.FunctionDeclarationTree) *rhino.Node {
	isDeclaration := tree.Kind == trees.FunctionDeclarationKindDeclaration
	isMember := tree.Kind == trees.FunctionDeclarationKindMember
	isArrow := tree.Kind == trees.FunctionDeclarationKindArrow
	isAsync := tree.IsAsync
	isGenerator := tree.IsGenerator

	var newName *rhino.Node
	if tree.Name != nil {
		newName = f.newStringNode(rhino.NAME, tree.Name.Value)
		f.setSourceInfoFromToken(newName, &tree.Name.Token)
	} else {
		if isDeclaration || isMember {
			newName = f.newStringNode(rhino.NAME, "__missing_name__")
		} else {
			newName = f.newStringNode(rhino.NAME, "")
		}
		f.setSourceInfo(newName, &tree.ParseTree)
	}

	node := f.newNode(rhino.FUNCTION)
	if isMember {
		newName.SetString("")
	}
	node.AddChildToBack(newName)

	if tree.FormalParameterList != nil {
		node.AddChildToBack(f.Transform(&tree.FormalParameterList.ParseTree))
	} else {
		node.AddChildToBack(f.newNode(rhino.PARAM_LIST))
	}

	if tree.FunctionBody != nil {
		node.AddChildToBack(f.Transform(tree.FunctionBody))
	} else {
		node.AddChildToBack(f.newNode(rhino.BLOCK))
	}

	node.SetIsGeneratorFunction(isGenerator)
	node.SetIsArrowFunction(isArrow)
	node.SetIsAsyncFunction(isAsync)
	node.PutBooleanProp(rhino.PropOptEs6Typed, tree.IsOptional)

	if isMember {
		f.setSourceInfo(node, &tree.ParseTree)
		memberName := f.newStringNode(rhino.MEMBER_FUNCTION_DEF, "")
		if tree.Name != nil {
			memberName.SetString_(tree.Name.Value)
			f.setSourceInfoFromToken(memberName, &tree.Name.Token)
		}
		memberName.AddChildToBack(node)
		memberName.SetStaticMember(tree.IsStatic)
		return memberName
	}
	return node
}

// transformFormalParameterList mirrors
// TransformDispatcher.processFormalParameterList(FormalParameterListTree).
func (f *IRFactory) transformFormalParameterList(tree *trees.FormalParameterListTree) *rhino.Node {
	params := f.newNode(rhino.PARAM_LIST)
	for _, param := range tree.Parameters {
		params.AddChildToBack(f.Transform(param))
	}
	return params
}

// transformDefaultParameter mirrors
// TransformDispatcher.processDefaultParameter(DefaultParameterTree).
func (f *IRFactory) transformDefaultParameter(tree *trees.DefaultParameterTree) *rhino.Node {
	targetNode := f.Transform(tree.Lhs)
	defaultValue := f.Transform(tree.DefaultValue)
	node := f.newNode(rhino.DEFAULT_VALUE)
	node.AddChildToBack(targetNode)
	node.AddChildToBack(defaultValue)
	f.setSourceInfo(node, &tree.ParseTree)
	return node
}

// transformClassDeclaration mirrors
// TransformDispatcher.processClassDeclaration(ClassDeclarationTree).
func (f *IRFactory) transformClassDeclaration(tree *trees.ClassDeclarationTree) *rhino.Node {
	var name *rhino.Node
	if tree.Name != nil {
		name = f.newStringNode(rhino.NAME, tree.Name.Value)
		f.setSourceInfoFromToken(name, &tree.Name.Token)
	} else {
		name = f.newNode(rhino.EMPTY)
		f.setSourceInfo(name, &tree.ParseTree)
	}

	var superClass *rhino.Node
	if tree.SuperClass != nil {
		superClass = f.Transform(tree.SuperClass)
	} else {
		superClass = f.newNode(rhino.EMPTY)
		f.setSourceInfo(superClass, &tree.ParseTree)
	}

	classMembers := f.newNode(rhino.CLASS_MEMBERS)
	f.setSourceInfo(classMembers, &tree.ParseTree)
	for _, child := range tree.Elements {
		classMembers.AddChildToBack(f.Transform(child))
	}

	node := f.newNode(rhino.CLASS)
	node.AddChildToBack(name)
	node.AddChildToBack(superClass)
	node.AddChildToBack(classMembers)
	return node
}

// transformFieldDeclaration mirrors
// TransformDispatcher.processField(FieldDeclarationTree).
func (f *IRFactory) transformFieldDeclaration(tree *trees.FieldDeclarationTree) *rhino.Node {
	var name string
	if tree.Name != nil {
		name = tree.Name.Value
	}
	node := f.newStringNode(rhino.MEMBER_FIELD_DEF, name)
	if tree.Initializer != nil {
		node.AddChildToBack(f.Transform(tree.Initializer))
	}
	node.PutBooleanProp(rhino.PropStaticMember, tree.IsStatic)
	return node
}

// transformVariableStatement mirrors
// TransformDispatcher.processVariableStatement(VariableStatementTree).
// Java delegates directly to the inner VariableDeclarationListTree.
func (f *IRFactory) transformVariableStatement(tree *trees.VariableStatementTree) *rhino.Node {
	if tree.Declarations != nil {
		return f.Transform(&tree.Declarations.ParseTree)
	}
	return f.newNode(rhino.EMPTY)
}

// transformVariableDeclarationList mirrors
// TransformDispatcher.processVariableDeclarationList(VariableDeclarationListTree).
func (f *IRFactory) transformVariableDeclarationList(tree *trees.VariableDeclarationListTree) *rhino.Node {
	var tok rhino.Token
	switch tree.DeclarationType {
	case trees.TokenTypeConst:
		tok = rhino.CONST
	case trees.TokenTypeLet:
		tok = rhino.LET
	default:
		tok = rhino.VAR
	}
	node := f.newNode(tok)
	for _, child := range tree.Declarations {
		node.AddChildToBack(f.Transform(&child.ParseTree))
	}
	return node
}

// transformVariableDeclaration mirrors
// TransformDispatcher.processVariableDeclaration(VariableDeclarationTree).
func (f *IRFactory) transformVariableDeclaration(tree *trees.VariableDeclarationTree) *rhino.Node {
	lhsNode := f.Transform(tree.Lvalue)
	var lhs *rhino.Node
	if lhsNode != nil && lhsNode.IsDestructuringPattern() {
		lhs = f.newNode(rhino.DESTRUCTURING_LHS)
		lhs.AddChildToBack(lhsNode)
	} else {
		lhs = lhsNode
	}
	if tree.Initializer != nil {
		lhs.AddChildToBack(f.Transform(tree.Initializer))
	}
	return lhs
}

// transformImportDeclaration mirrors
// TransformDispatcher.processImportDecl(ImportDeclarationTree).
func (f *IRFactory) transformImportDeclaration(tree *trees.ImportDeclarationTree) *rhino.Node {
	// First child: default binding or EMPTY.
	var firstChild *rhino.Node
	if tree.DefaultBindingIdentifier != nil {
		firstChild = f.newStringNode(rhino.NAME, tree.DefaultBindingIdentifier.Value)
		f.setSourceInfoFromToken(firstChild, &tree.DefaultBindingIdentifier.Token)
	} else {
		firstChild = f.newNode(rhino.EMPTY)
		f.setSourceInfo(firstChild, &tree.ParseTree)
	}

	// Second child: namespace import or import-spec list.
	var secondChild *rhino.Node
	if tree.NameSpaceImportIdentifier != nil {
		secondChild = f.newStringNode(rhino.IMPORT_STAR, tree.NameSpaceImportIdentifier.Value)
		f.setSourceInfoFromToken(secondChild, &tree.NameSpaceImportIdentifier.Token)
	} else if len(tree.ImportSpecifierList) > 0 {
		secondChild = f.newNode(rhino.IMPORT_SPECS)
		f.setSourceInfo(secondChild, &tree.ParseTree)
		for _, spec := range tree.ImportSpecifierList {
			secondChild.AddChildToBack(f.Transform(spec))
		}
	} else {
		secondChild = f.newNode(rhino.EMPTY)
		f.setSourceInfo(secondChild, &tree.ParseTree)
	}

	// Third child: module specifier string.
	// Note: LiteralToken is a stub and carries no value field; the string
	// value is unavailable until the token layer is fully ported.
	thirdChild := f.newStringNode(rhino.STRINGLIT, "")

	node := f.newNode(rhino.IMPORT)
	node.AddChildToBack(firstChild)
	node.AddChildToBack(secondChild)
	node.AddChildToBack(thirdChild)
	return node
}

// transformImportSpecifier mirrors
// TransformDispatcher.processImportSpec(ImportSpecifierTree).
func (f *IRFactory) transformImportSpecifier(tree *trees.ImportSpecifierTree) *rhino.Node {
	var importedName *rhino.Node
	if tree.ImportedName != nil {
		importedName = f.newStringNode(rhino.NAME, tree.ImportedName.Value)
	} else {
		importedName = f.newStringNode(rhino.NAME, "")
	}
	importSpec := f.newNode(rhino.IMPORT_SPEC)
	importSpec.AddChildToBack(importedName)
	if tree.DestinationName == nil {
		importSpec.SetShorthandProperty(true)
		importSpec.AddChildToBack(importedName.CloneTree())
	} else {
		destName := f.newStringNode(rhino.NAME, tree.DestinationName.Value)
		importSpec.AddChildToBack(destName)
	}
	return importSpec
}

// transformExportDeclaration mirrors
// TransformDispatcher.processExportDecl(ExportDeclarationTree).
func (f *IRFactory) transformExportDeclaration(tree *trees.ExportDeclarationTree) *rhino.Node {
	var decls *rhino.Node
	if tree.IsExportAll {
		// export * from 'mod' — no declaration or spec list.
		decls = f.newNode(rhino.EMPTY)
	} else if tree.Declaration != nil {
		decls = f.Transform(tree.Declaration)
	} else if len(tree.ExportSpecifierList) > 0 {
		decls = f.newNode(rhino.EXPORT_SPECS)
		for _, spec := range tree.ExportSpecifierList {
			decls.AddChildToBack(f.Transform(spec))
		}
	} else {
		decls = f.newNode(rhino.EMPTY)
	}
	if decls == nil {
		decls = f.newNode(rhino.EMPTY)
	}
	f.setSourceInfo(decls, &tree.ParseTree)

	export := f.newNode(rhino.EXPORT)
	export.AddChildToBack(decls)
	if tree.From != nil {
		// LiteralToken carries no value field in the current stub.
		from := f.newStringNode(rhino.STRINGLIT, "")
		export.AddChildToBack(from)
	}
	export.PutBooleanProp(rhino.PropExportAll, tree.IsExportAll)
	export.PutBooleanProp(rhino.PropExportDefault, tree.IsDefault)
	return export
}

// transformExportSpecifier mirrors
// TransformDispatcher.processExportSpec(ExportSpecifierTree).
func (f *IRFactory) transformExportSpecifier(tree *trees.ExportSpecifierTree) *rhino.Node {
	var importedName *rhino.Node
	if tree.ImportedName != nil {
		importedName = f.newStringNode(rhino.NAME, tree.ImportedName.Value)
	} else {
		importedName = f.newStringNode(rhino.NAME, "")
	}
	exportSpec := f.newNode(rhino.EXPORT_SPEC)
	exportSpec.AddChildToBack(importedName)
	if tree.DestinationName == nil {
		exportSpec.SetShorthandProperty(true)
		exportSpec.AddChildToBack(importedName.CloneTree())
	} else {
		destName := f.newStringNode(rhino.NAME, tree.DestinationName.Value)
		exportSpec.AddChildToBack(destName)
	}
	return exportSpec
}

// ---- primitive expressions -----------------------------------------

// transformIdentifierExpression mirrors
// TransformDispatcher.processName(IdentifierExpressionTree, IdentifierType).
func (f *IRFactory) transformIdentifierExpression(tree *trees.IdentifierExpressionTree) *rhino.Node {
	if tree.IdentifierToken == nil {
		return f.newStringNode(rhino.NAME, "__missing_name__")
	}
	return f.newStringNode(rhino.NAME, tree.IdentifierToken.Value)
}

// transformLiteralExpression dispatches on the literal token type.
// Mirrors the per-type literal processing in TransformDispatcher.
func (f *IRFactory) transformLiteralExpression(tree *trees.LiteralExpressionTree) *rhino.Node {
	if tree.LiteralToken == nil {
		return f.newNode(rhino.EMPTY)
	}
	tok := tree.LiteralToken
	switch tok.Type {
	case trees.TokenTypeTrue:
		return f.newNode(rhino.TRUE)
	case trees.TokenTypeFalse:
		return f.newNode(rhino.FALSE)
	case trees.TokenTypeNull:
		return f.newNode(rhino.NULL)
	case trees.TokenTypeNumber:
		// trees.Token has no Value field; number value is not accessible via the
		// stub. Return 0 as a placeholder until the full token is wired.
		return f.newNumberNode(0)
	case trees.TokenTypeString:
		return f.newStringNode(rhino.STRINGLIT, "")
	case trees.TokenTypeRegularExpression:
		return f.transformRegexpLiteral("")
	default:
		return f.newNode(rhino.EMPTY)
	}
}

// transformThisExpression mirrors
// TransformDispatcher.processThisExpression(ThisExpressionTree).
func (f *IRFactory) transformThisExpression() *rhino.Node {
	return f.newNode(rhino.THIS)
}

// transformSuperExpression mirrors Java super expression → Token.SUPER.
func (f *IRFactory) transformSuperExpression() *rhino.Node {
	return f.newNode(rhino.SUPER)
}

// ---- unary / binary / update ---------------------------------------

// transformUnaryExpression mirrors
// TransformDispatcher.processUnaryExpression(UnaryExpressionTree).
func (f *IRFactory) transformUnaryExpression(tree *trees.UnaryExpressionTree) *rhino.Node {
	if tree.Operator == nil {
		return f.newNode(rhino.EMPTY)
	}
	tok := transformUnaryTokenType(tree.Operator.Type)
	operand := f.Transform(tree.Operand)
	node := f.newNode(tok)
	node.AddChildToBack(operand)
	return node
}

// transformUpdateExpression mirrors
// TransformDispatcher.processUpdateExpression(UpdateExpressionTree).
func (f *IRFactory) transformUpdateExpression(tree *trees.UpdateExpressionTree) *rhino.Node {
	if tree.Operator == nil {
		return f.newNode(rhino.EMPTY)
	}
	tok := transformUpdateTokenType(tree.Operator.Type)
	operand := f.Transform(tree.Operand)
	node := f.newNode(tok)
	node.AddChildToBack(operand)
	postfix := tree.OperatorPosition == trees.OperatorPositionPostfix
	node.PutBooleanProp(rhino.PropIncrDecr, postfix)
	return node
}

// transformBinaryOperator mirrors
// TransformDispatcher.processBinaryExpression(BinaryOperatorTree).
func (f *IRFactory) transformBinaryOperator(tree *trees.BinaryOperatorTree) *rhino.Node {
	if tree.Operator == nil {
		return f.newNode(rhino.EMPTY)
	}
	tok := transformBinaryTokenType(tree.Operator.Type)
	left := f.Transform(tree.Left)
	right := f.Transform(tree.Right)
	node := f.newNode(tok)
	node.AddChildToBack(left)
	node.AddChildToBack(right)
	return node
}

// transformConditionalExpression mirrors
// TransformDispatcher.processConditionalExpression(ConditionalExpressionTree).
func (f *IRFactory) transformConditionalExpression(tree *trees.ConditionalExpressionTree) *rhino.Node {
	node := f.newNode(rhino.HOOK)
	node.AddChildToBack(f.Transform(tree.Condition))
	node.AddChildToBack(f.Transform(tree.Left))
	node.AddChildToBack(f.Transform(tree.Right))
	return node
}

// transformCommaExpression mirrors
// TransformDispatcher.processCommaExpression(CommaExpressionTree).
func (f *IRFactory) transformCommaExpression(tree *trees.CommaExpressionTree) *rhino.Node {
	if len(tree.Expressions) == 0 {
		return f.newNode(rhino.EMPTY)
	}
	// Java builds a left-associative chain of COMMA nodes.
	root := f.newNode(rhino.COMMA)
	for i, expr := range tree.Expressions {
		if i < 2 {
			root.AddChildToBack(f.Transform(expr))
		} else {
			parent := f.newNode(rhino.COMMA)
			parent.AddChildToBack(root)
			parent.AddChildToBack(f.Transform(expr))
			root = parent
		}
	}
	return root
}

// transformParenExpression mirrors
// TransformDispatcher.processParenthesizedExpression(ParenExpressionTree).
func (f *IRFactory) transformParenExpression(tree *trees.ParenExpressionTree) *rhino.Node {
	expr := f.Transform(tree.Expression)
	if expr != nil {
		expr.SetIsParenthesized(true)
	}
	return expr
}

// ---- calls / member access / new -----------------------------------

// transformCallExpression mirrors
// TransformDispatcher.processFunctionCall(CallExpressionTree).
func (f *IRFactory) transformCallExpression(tree *trees.CallExpressionTree) *rhino.Node {
	node := f.newNode(rhino.CALL)
	node.AddChildToBack(f.Transform(tree.Operand))
	if tree.Arguments != nil {
		for _, arg := range tree.Arguments.Arguments {
			node.AddChildToBack(f.Transform(arg))
		}
	}
	// Free-call annotation: if callee is not a GETPROP/GETELEM, mark FREE_CALL.
	callee := node.FirstChild()
	if callee != nil && !callee.IsGetProp() && !callee.IsGetElem() {
		node.PutBooleanProp(rhino.PropFreeCall, true)
	}
	return node
}

// transformNewExpression mirrors
// TransformDispatcher.processNewExpression(NewExpressionTree).
func (f *IRFactory) transformNewExpression(tree *trees.NewExpressionTree) *rhino.Node {
	node := f.newNode(rhino.NEW)
	node.AddChildToBack(f.Transform(tree.Operand))
	if tree.Arguments != nil {
		for _, arg := range tree.Arguments.Arguments {
			node.AddChildToBack(f.Transform(arg))
		}
	}
	return node
}

// transformMemberExpression mirrors
// TransformDispatcher.processPropertyGet(MemberExpressionTree).
func (f *IRFactory) transformMemberExpression(tree *trees.MemberExpressionTree) *rhino.Node {
	leftChild := f.Transform(tree.Operand)
	if tree.MemberName == nil {
		return leftChild
	}
	getProp := f.newStringNode(rhino.GETPROP, tree.MemberName.Value)
	getProp.AddChildToBack(leftChild)
	return getProp
}

// transformMemberLookupExpression mirrors
// TransformDispatcher.processElementGet(MemberLookupExpressionTree).
func (f *IRFactory) transformMemberLookupExpression(tree *trees.MemberLookupExpressionTree) *rhino.Node {
	node := f.newNode(rhino.GETELEM)
	node.AddChildToBack(f.Transform(tree.Operand))
	node.AddChildToBack(f.Transform(tree.MemberExpression))
	return node
}

// ---- structured literals -------------------------------------------

// transformArrayLiteralExpression mirrors
// TransformDispatcher.processArrayLiteral(ArrayLiteralExpressionTree).
func (f *IRFactory) transformArrayLiteralExpression(tree *trees.ArrayLiteralExpressionTree) *rhino.Node {
	node := f.newNode(rhino.ARRAYLIT)
	for _, el := range tree.Elements {
		node.AddChildToBack(f.Transform(el))
	}
	return node
}

// transformObjectLiteralExpression mirrors
// TransformDispatcher.processObjectLiteral(ObjectLiteralExpressionTree).
func (f *IRFactory) transformObjectLiteralExpression(tree *trees.ObjectLiteralExpressionTree) *rhino.Node {
	node := f.newNode(rhino.OBJECTLIT)
	for _, el := range tree.PropertyNameAndValues {
		key := f.Transform(el)
		if key != nil {
			node.AddChildToBack(key)
		}
	}
	return node
}

// transformPropertyNameAssignment mirrors
// TransformDispatcher.processPropertyNameAssignment(PropertyNameAssignmentTree).
func (f *IRFactory) transformPropertyNameAssignment(tree *trees.PropertyNameAssignmentTree) *rhino.Node {
	key := f.processObjectLitKey(tree.Name, rhino.STRING_KEY)
	if tree.Value != nil {
		key.AddChildToFront(f.Transform(tree.Value))
	} else if key.Token() == rhino.STRING_KEY {
		// Shorthand property: value is a NAME node with the same string.
		val := f.newStringNode(rhino.NAME, key.String_())
		key.AddChildToFront(val)
		key.PutBooleanProp(rhino.PropShorthandProperty, true)
	}
	return key
}

// transformGetAccessor mirrors
// TransformDispatcher.processGetAccessor(GetAccessorTree).
func (f *IRFactory) transformGetAccessor(tree *trees.GetAccessorTree) *rhino.Node {
	key := f.processObjectLitKey(tree.PropertyName, rhino.GETTER_DEF)
	body := f.Transform(&tree.Body.ParseTree)
	dummyName := f.newStringNode(rhino.NAME, "")
	f.setSourceInfo(dummyName, &tree.Body.ParseTree)
	paramList := f.newNode(rhino.PARAM_LIST)
	f.setSourceInfo(paramList, &tree.Body.ParseTree)
	fn := f.newNode(rhino.FUNCTION)
	fn.AddChildToBack(dummyName)
	fn.AddChildToBack(paramList)
	fn.AddChildToBack(body)
	f.setSourceInfo(fn, &tree.Body.ParseTree)
	key.AddChildToFront(fn)
	if tree.IsStatic {
		key.PutBooleanProp(rhino.PropStaticMember, true)
	}
	return key
}

// transformSetAccessor mirrors
// TransformDispatcher.processSetAccessor(SetAccessorTree).
func (f *IRFactory) transformSetAccessor(tree *trees.SetAccessorTree) *rhino.Node {
	key := f.processObjectLitKey(tree.PropertyName, rhino.SETTER_DEF)
	// Build a minimal PARAM_LIST from the formal parameter.
	paramList := f.newNode(rhino.PARAM_LIST)
	if tree.Parameter != nil {
		for _, p := range tree.Parameter.Parameters {
			paramList.AddChildToBack(f.Transform(p))
		}
		f.setSourceInfo(paramList, &tree.Parameter.ParseTree)
	}
	body := f.Transform(&tree.Body.ParseTree)
	dummyName := f.newStringNode(rhino.NAME, "")
	if tree.PropertyName != nil {
		f.setSourceInfoFromToken(dummyName, tree.PropertyName)
	}
	fn := f.newNode(rhino.FUNCTION)
	fn.AddChildToBack(dummyName)
	fn.AddChildToBack(paramList)
	fn.AddChildToBack(body)
	f.setSourceInfo(fn, &tree.Body.ParseTree)
	key.AddChildToFront(fn)
	if tree.IsStatic {
		key.PutBooleanProp(rhino.PropStaticMember, true)
	}
	return key
}

// transformTemplateLiteralExpression mirrors
// TransformDispatcher.processTemplateLiteral(TemplateLiteralExpressionTree).
func (f *IRFactory) transformTemplateLiteralExpression(tree *trees.TemplateLiteralExpressionTree) *rhino.Node {
	templateLitNode := f.newNode(rhino.TEMPLATELIT)
	f.setSourceInfo(templateLitNode, &tree.ParseTree)
	var node *rhino.Node
	if tree.Operand == nil {
		node = templateLitNode
	} else {
		node = f.newNode(rhino.TAGGED_TEMPLATELIT)
		node.AddChildToBack(f.Transform(tree.Operand))
		node.AddChildToBack(templateLitNode)
	}
	for _, child := range tree.Elements {
		templateLitNode.AddChildToBack(f.Transform(child))
	}
	if node.Token() == rhino.TAGGED_TEMPLATELIT {
		// Free-call annotation for tagged template literals.
		callee := node.FirstChild()
		if callee != nil && !callee.IsGetProp() && !callee.IsGetElem() {
			node.PutBooleanProp(rhino.PropFreeCall, true)
		}
	}
	return node
}

// transformTemplateLiteralPortion mirrors
// TransformDispatcher.processTemplateLiteralPortion(TemplateLiteralPortionTree).
// Note: trees.Token has no Value field in the stub; cooked/raw are empty stubs.
func (f *IRFactory) transformTemplateLiteralPortion(tree *trees.TemplateLiteralPortionTree) *rhino.Node {
	return rhino.NewTemplateLitString("", "")
}

// transformTemplateSubstitution mirrors
// TransformDispatcher.processTemplateSubstitution(TemplateSubstitutionTree).
func (f *IRFactory) transformTemplateSubstitution(tree *trees.TemplateSubstitutionTree) *rhino.Node {
	node := f.newNode(rhino.TEMPLATELIT_SUB)
	node.AddChildToBack(f.Transform(tree.Expression))
	return node
}

// ---- helpers -------------------------------------------------------

// processObjectLitKey creates a STRING_KEY / GETTER_DEF / SETTER_DEF node
// from the property token. Mirrors IRFactory.TransformDispatcher.processObjectLitKey.
// Note: trees.Token is a stub without a Value field; the property name
// is not accessible through the current trees.Token type.
func (f *IRFactory) processObjectLitKey(tok *trees.Token, output rhino.Token) *rhino.Node {
	if tok == nil {
		return f.newStringNode(output, "__missing_key__")
	}
	switch tok.Type {
	case trees.TokenTypeIdentifier:
		// Identifier key: name carried by IdentifierToken.Value (not on trees.Token stub).
		node := f.newStringNode(output, "")
		f.setSourceInfoFromToken(node, tok)
		return node
	case trees.TokenTypeNumber:
		node := f.newStringNode(output, "0")
		node.PutBooleanProp(rhino.PropQuotedStringKey, true)
		f.setSourceInfoFromToken(node, tok)
		return node
	default:
		node := f.newStringNode(output, "")
		node.PutBooleanProp(rhino.PropQuotedStringKey, true)
		f.setSourceInfoFromToken(node, tok)
		return node
	}
}

// transformRegexpLiteral builds a REGEXP node from the raw literal string.
// Mirrors TransformDispatcher.processRegExpLiteral.
func (f *IRFactory) transformRegexpLiteral(raw string) *rhino.Node {
	lastSlash := len(raw) - 1
	for lastSlash > 0 && raw[lastSlash] != '/' {
		lastSlash--
	}
	pattern := raw[1:lastSlash]
	flags := raw[lastSlash+1:]
	patternNode := f.newStringNode(rhino.STRINGLIT, pattern)
	node := f.newNode(rhino.REGEXP)
	node.AddChildToBack(patternNode)
	if flags != "" {
		flagsNode := f.newStringNode(rhino.STRINGLIT, flags)
		node.AddChildToBack(flagsNode)
	}
	return node
}

// setSourceInfoFromToken copies position from a trees.Token to a rhino Node.
func (f *IRFactory) setSourceInfoFromToken(node *rhino.Node, tok *trees.Token) {
	if node.Lineno() == 0 {
		node.SetPos(tok.Location.Start.Line+1, tok.Location.Start.Column)
	}
}

// ---- token-type mapping functions ---------------------------------

// transformUnaryTokenType maps a parse-tree TokenType to the rhino unary
// operator token. Mirrors IRFactory.transformUnaryTokenType(TokenType).
func transformUnaryTokenType(t trees.TokenType) rhino.Token {
	switch t {
	case trees.TokenTypeBang:
		return rhino.NOT
	case trees.TokenTypeTilde:
		return rhino.BITNOT
	case trees.TokenTypePlus:
		return rhino.POS
	case trees.TokenTypeMinus:
		return rhino.NEG
	case trees.TokenTypeDelete:
		return rhino.DELPROP
	case trees.TokenTypeTypeof:
		return rhino.TYPEOF
	case trees.TokenTypeVoid:
		return rhino.VOID
	default:
		return rhino.EMPTY
	}
}

// transformUpdateTokenType maps a parse-tree TokenType to the rhino
// increment/decrement token.
func transformUpdateTokenType(t trees.TokenType) rhino.Token {
	switch t {
	case trees.TokenTypePlusPlus:
		return rhino.INC
	case trees.TokenTypeMinusMinus:
		return rhino.DEC
	default:
		return rhino.EMPTY
	}
}

// transformBinaryTokenType maps a parse-tree TokenType to the rhino binary
// operator token. Mirrors IRFactory.transformBinaryTokenType(TokenType).
func transformBinaryTokenType(t trees.TokenType) rhino.Token {
	switch t {
	case trees.TokenTypeBar:
		return rhino.BITOR
	case trees.TokenTypeCaret:
		return rhino.BITXOR
	case trees.TokenTypeAmpersand:
		return rhino.BITAND
	case trees.TokenTypeEqualEqual:
		return rhino.EQ
	case trees.TokenTypeNotEqual:
		return rhino.NE
	case trees.TokenTypeOpenAngle:
		return rhino.LT
	case trees.TokenTypeLessEqual:
		return rhino.LE
	case trees.TokenTypeCloseAngle:
		return rhino.GT
	case trees.TokenTypeGreaterEqual:
		return rhino.GE
	case trees.TokenTypeLeftShift:
		return rhino.LSH
	case trees.TokenTypeRightShift:
		return rhino.RSH
	case trees.TokenTypeUnsignedRightShift:
		return rhino.URSH
	case trees.TokenTypePlus:
		return rhino.ADD
	case trees.TokenTypeMinus:
		return rhino.SUB
	case trees.TokenTypeStar:
		return rhino.MUL
	case trees.TokenTypeSlash:
		return rhino.DIV
	case trees.TokenTypePercent:
		return rhino.MOD
	case trees.TokenTypeStarStar:
		return rhino.EXPONENT
	case trees.TokenTypeEqualEqualEqual:
		return rhino.SHEQ
	case trees.TokenTypeNotEqualEqual:
		return rhino.SHNE
	case trees.TokenTypeIn:
		return rhino.IN
	case trees.TokenTypeInstanceof:
		return rhino.INSTANCEOF
	case trees.TokenTypeComma:
		return rhino.COMMA
	case trees.TokenTypeEqual:
		return rhino.ASSIGN
	case trees.TokenTypeBarEqual:
		return rhino.ASSIGN_BITOR
	case trees.TokenTypeCaretEqual:
		return rhino.ASSIGN_BITXOR
	case trees.TokenTypeAmpersandEqual:
		return rhino.ASSIGN_BITAND
	case trees.TokenTypeLeftShiftEqual:
		return rhino.ASSIGN_LSH
	case trees.TokenTypeRightShiftEqual:
		return rhino.ASSIGN_RSH
	case trees.TokenTypeUnsignedRightShiftEqual:
		return rhino.ASSIGN_URSH
	case trees.TokenTypePlusEqual:
		return rhino.ASSIGN_ADD
	case trees.TokenTypeMinusEqual:
		return rhino.ASSIGN_SUB
	case trees.TokenTypeStarEqual:
		return rhino.ASSIGN_MUL
	case trees.TokenTypeStarStarEqual:
		return rhino.ASSIGN_EXPONENT
	case trees.TokenTypeSlashEqual:
		return rhino.ASSIGN_DIV
	case trees.TokenTypePercentEqual:
		return rhino.ASSIGN_MOD
	case trees.TokenTypeOr:
		return rhino.OR
	case trees.TokenTypeAnd:
		return rhino.AND
	case trees.TokenTypeQuestionQuestion:
		return rhino.COALESCE
	case trees.TokenTypeOrEqual:
		return rhino.ASSIGN_OR
	case trees.TokenTypeAndEqual:
		return rhino.ASSIGN_AND
	case trees.TokenTypeQuestionQuestionEqual:
		return rhino.ASSIGN_COALESCE
	default:
		return rhino.EMPTY
	}
}

// ---- string/number utilities ---------------------------------------

// normalizeString does minimal string literal normalization (strip quotes).
// A full implementation would handle escape sequences; this stub is
// sufficient for the current corpus tier.
func normalizeString(raw string) string {
	if len(raw) >= 2 {
		q := raw[0]
		if (q == '"' || q == '\'' || q == '`') && raw[len(raw)-1] == q {
			return raw[1 : len(raw)-1]
		}
	}
	return raw
}

// parseFloat converts a JS number literal string to float64.
// A full implementation would handle hex/octal/binary; this handles the
// common decimal case.
func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

// floatToString converts a float64 to a JS-style string (used for numeric
// property keys). Mirrors DToA.numberToString.
func floatToString(v float64) string {
	return fmt.Sprintf("%g", v)
}
