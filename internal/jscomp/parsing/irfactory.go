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

// Go-only: Go constructor — Java uses the static IRFactory.transformTree(...) factory method.
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
	node := f.process(tree)
	if node != nil {
		f.setSourceInfo(node, tree)
	}
	return node
}

// process routes a ParseTree to the appropriate transform* method based on
// its type. Mirrors the TransformDispatcher.process(ParseTree) switch in Java.
func (f *IRFactory) process(tree *trees.ParseTree) *rhino.Node {
	switch tree.Type {
	// ---- programme / block -----------------------------------------
	case trees.ParseTreeTypeProgram:
		return f.processAstRoot((*trees.ProgramTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeBlock:
		return f.processBlock((*trees.BlockTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeExpressionStatement:
		return f.processExpressionStatement((*trees.ExpressionStatementTree)(unsafe.Pointer(tree)))

	// ---- basic statements ------------------------------------------
	case trees.ParseTreeTypeReturnStatement:
		return f.processReturnStatement((*trees.ReturnStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeThrowStatement:
		return f.processThrowStatement((*trees.ThrowStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeBreakStatement:
		return f.processBreakStatement((*trees.BreakStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeContinueStatement:
		return f.processContinueStatement((*trees.ContinueStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeIfStatement:
		return f.processIfStatement((*trees.IfStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeWhileStatement:
		return f.processWhileLoop((*trees.WhileStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeDoWhileStatement:
		return f.processDoLoop((*trees.DoWhileStatementTree)(unsafe.Pointer(tree)))

	// ---- loops + try/switch ----------------------------------------
	case trees.ParseTreeTypeForStatement:
		return f.processForLoop((*trees.ForStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeForInStatement:
		return f.processForInLoop((*trees.ForInStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeForOfStatement:
		return f.processForOf((*trees.ForOfStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeForAwaitOfStatement:
		return f.processForAwaitOf((*trees.ForAwaitOfStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeSwitchStatement:
		return f.processSwitchStatement((*trees.SwitchStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeCaseClause:
		return f.processSwitchCase((*trees.CaseClauseTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeDefaultClause:
		return f.processSwitchDefault((*trees.DefaultClauseTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeTryStatement:
		return f.processTryStatement((*trees.TryStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeCatch:
		return f.processCatchClause((*trees.CatchTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeFinally:
		return f.processFinally((*trees.FinallyTree)(unsafe.Pointer(tree)))

	// ---- declarations + functions + classes ------------------------
	case trees.ParseTreeTypeFunctionDeclaration:
		return f.processFunction((*trees.FunctionDeclarationTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeFormalParameterList:
		return f.processFormalParameterList((*trees.FormalParameterListTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeDefaultParameter:
		return f.processDefaultParameter((*trees.DefaultParameterTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeClassDeclaration:
		return f.processClassDeclaration((*trees.ClassDeclarationTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeFieldDeclaration:
		return f.processField((*trees.FieldDeclarationTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeVariableStatement:
		return f.processVariableStatement((*trees.VariableStatementTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeVariableDeclarationList:
		return f.processVariableDeclarationList((*trees.VariableDeclarationListTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeVariableDeclaration:
		return f.processVariableDeclaration((*trees.VariableDeclarationTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeImportDeclaration:
		return f.processImportDecl((*trees.ImportDeclarationTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeImportSpecifier:
		return f.processImportSpec((*trees.ImportSpecifierTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeExportDeclaration:
		return f.processExportDecl((*trees.ExportDeclarationTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeExportSpecifier:
		return f.processExportSpec((*trees.ExportSpecifierTree)(unsafe.Pointer(tree)))

	// ---- primitive expressions -------------------------------------
	case trees.ParseTreeTypeIdentifierExpression:
		return f.processName((*trees.IdentifierExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeLiteralExpression:
		return f.processLiteralExpression((*trees.LiteralExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeThisExpression:
		return f.processThisExpression()
	case trees.ParseTreeTypeSuperExpression:
		return f.processSuper()

	// ---- unary / binary / update -----------------------------------
	case trees.ParseTreeTypeUnaryExpression:
		return f.processUnaryExpression((*trees.UnaryExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeUpdateExpression:
		return f.processUpdateExpression((*trees.UpdateExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeBinaryOperator:
		return f.processBinaryExpression((*trees.BinaryOperatorTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeConditionalExpression:
		return f.processConditionalExpression((*trees.ConditionalExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeCommaExpression:
		return f.processCommaExpression((*trees.CommaExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeParenExpression:
		return f.processParenthesizedExpression((*trees.ParenExpressionTree)(unsafe.Pointer(tree)))

	// ---- calls / member access / new -------------------------------
	case trees.ParseTreeTypeCallExpression:
		return f.processFunctionCall((*trees.CallExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeNewExpression:
		return f.processNewExpression((*trees.NewExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeMemberExpression:
		return f.processPropertyGet((*trees.MemberExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeMemberLookupExpression:
		return f.processElementGet((*trees.MemberLookupExpressionTree)(unsafe.Pointer(tree)))

	// ---- structured literals ---------------------------------------
	case trees.ParseTreeTypeArrayLiteralExpression:
		return f.processArrayLiteral((*trees.ArrayLiteralExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeObjectLiteralExpression:
		return f.processObjectLiteral((*trees.ObjectLiteralExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypePropertyNameAssignment:
		return f.processPropertyNameAssignment((*trees.PropertyNameAssignmentTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeGetAccessor:
		return f.processGetAccessor((*trees.GetAccessorTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeSetAccessor:
		return f.processSetAccessor((*trees.SetAccessorTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeTemplateLiteralExpression:
		return f.processTemplateLiteral((*trees.TemplateLiteralExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeTemplateLiteralPortion:
		return f.processTemplateLiteralPortion((*trees.TemplateLiteralPortionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeTemplateSubstitution:
		return f.processTemplateSubstitution((*trees.TemplateSubstitutionTree)(unsafe.Pointer(tree)))

	// ---- await / dynamic-import / debugger ------------------------
	case trees.ParseTreeTypeAwaitExpression:
		return f.transformAwaitExpression((*trees.AwaitExpressionTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeDynamicImportExpression:
		return f.transformDynamicImport((*trees.DynamicImportTree)(unsafe.Pointer(tree)))
	case trees.ParseTreeTypeDebuggerStatement:
		return f.newNode(rhino.DEBUGGER)

	// ---- null / empty placeholder ----------------------------------
	case trees.ParseTreeTypeNull:
		return rhino.New(rhino.EMPTY)

	default:
		return rhino.New(rhino.EMPTY)
	}
}

// transformAwaitExpression mirrors IRFactory.processAwait(AwaitExpressionTree).
func (f *IRFactory) transformAwaitExpression(tree *trees.AwaitExpressionTree) *rhino.Node {
	await := f.newNode(rhino.AWAIT)
	await.AddChildToBack(f.Transform(tree.Expression))
	return await
}

// transformDynamicImport mirrors IRFactory.processDynamicImport(DynamicImportTree).
func (f *IRFactory) transformDynamicImport(tree *trees.DynamicImportTree) *rhino.Node {
	n := f.newNode(rhino.DYNAMIC_IMPORT)
	n.AddChildToBack(f.Transform(tree.Argument))
	return n
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

// processAstRoot mirrors TransformDispatcher.processAstRoot(ProgramTree).
func (f *IRFactory) processAstRoot(tree *trees.ProgramTree) *rhino.Node {
	script := f.newNode(rhino.SCRIPT)
	for _, child := range tree.SourceElements {
		script.AddChildToBack(f.Transform(child))
	}
	return script
}

// processBlock mirrors TransformDispatcher.processBlock(BlockTree).
func (f *IRFactory) processBlock(tree *trees.BlockTree) *rhino.Node {
	node := f.newNode(rhino.BLOCK)
	for _, child := range tree.Statements {
		node.AddChildToBack(f.Transform(child))
	}
	return node
}

// Go-only: dispatches and wraps single-statement bodies in a synthetic BLOCK.
// Java inlines this wrapping at each call site via processBlock(ParseTree node).
func (f *IRFactory) processBlockTree(tree *trees.ParseTree) *rhino.Node {
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

// Go-only: helper to wrap a bare expression ParseTree in EXPR_RESULT.
// Java's processExpressionStatement always has an ExpressionStatementTree wrapper.
func (f *IRFactory) transformExpressionStatementExpr(expr *trees.ParseTree) *rhino.Node {
	node := f.newNode(rhino.EXPR_RESULT)
	node.AddChildToBack(f.Transform(expr))
	return node
}

// processExpressionStatement mirrors
// TransformDispatcher.processExpressionStatement(ExpressionStatementTree).
func (f *IRFactory) processExpressionStatement(tree *trees.ExpressionStatementTree) *rhino.Node {
	node := f.newNode(rhino.EXPR_RESULT)
	node.AddChildToBack(f.Transform(tree.Expression))
	return node
}

// ---- basic statements ----------------------------------------------

// processReturnStatement mirrors
// TransformDispatcher.processReturnStatement(ReturnStatementTree).
func (f *IRFactory) processReturnStatement(tree *trees.ReturnStatementTree) *rhino.Node {
	node := f.newNode(rhino.RETURN)
	if tree.Expression != nil {
		node.AddChildToBack(f.Transform(tree.Expression))
	}
	return node
}

// processThrowStatement mirrors
// TransformDispatcher.processThrowStatement(ThrowStatementTree).
func (f *IRFactory) processThrowStatement(tree *trees.ThrowStatementTree) *rhino.Node {
	node := f.newNode(rhino.THROW)
	node.AddChildToBack(f.Transform(tree.Value))
	return node
}

// processBreakStatement mirrors
// TransformDispatcher.processBreakStatement(BreakStatementTree).
func (f *IRFactory) processBreakStatement(tree *trees.BreakStatementTree) *rhino.Node {
	node := f.newNode(rhino.BREAK)
	if tree.HasLabel() {
		labelName := f.transformLabelName(tree.Name)
		node.AddChildToBack(labelName)
	}
	return node
}

// processContinueStatement mirrors
// TransformDispatcher.processContinueStatement(ContinueStatementTree).
func (f *IRFactory) processContinueStatement(tree *trees.ContinueStatementTree) *rhino.Node {
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

// processIfStatement mirrors
// TransformDispatcher.processIfStatement(IfStatementTree).
func (f *IRFactory) processIfStatement(tree *trees.IfStatementTree) *rhino.Node {
	node := f.newNode(rhino.IF)
	node.AddChildToBack(f.Transform(tree.Condition))
	node.AddChildToBack(f.processBlockTree(tree.IfClause))
	if tree.ElseClause != nil {
		node.AddChildToBack(f.processBlockTree(tree.ElseClause))
	}
	return node
}

// processWhileLoop mirrors
// TransformDispatcher.processWhileLoop(WhileStatementTree).
func (f *IRFactory) processWhileLoop(tree *trees.WhileStatementTree) *rhino.Node {
	node := f.newNode(rhino.WHILE)
	node.AddChildToBack(f.Transform(tree.Condition))
	node.AddChildToBack(f.processBlockTree(tree.Body))
	return node
}

// processDoLoop mirrors
// TransformDispatcher.processDoLoop(DoWhileStatementTree).
func (f *IRFactory) processDoLoop(tree *trees.DoWhileStatementTree) *rhino.Node {
	node := f.newNode(rhino.DO)
	node.AddChildToBack(f.processBlockTree(tree.Body))
	node.AddChildToBack(f.Transform(tree.Condition))
	return node
}

// ---- loops + try/switch -------------------------------------------

// processForLoop mirrors
// TransformDispatcher.processForLoop(ForStatementTree).
func (f *IRFactory) processForLoop(tree *trees.ForStatementTree) *rhino.Node {
	node := f.newNode(rhino.FOR)
	node.AddChildToBack(f.transformOrEmpty(tree.Initializer, &tree.ParseTree))
	node.AddChildToBack(f.transformOrEmpty(tree.Condition, &tree.ParseTree))
	node.AddChildToBack(f.transformOrEmpty(tree.Increment, &tree.ParseTree))
	node.AddChildToBack(f.processBlockTree(tree.Body))
	return node
}

// processForInLoop mirrors
// TransformDispatcher.processForInLoop(ForInStatementTree).
func (f *IRFactory) processForInLoop(tree *trees.ForInStatementTree) *rhino.Node {
	node := f.newNode(rhino.FOR_IN)
	node.AddChildToBack(f.Transform(tree.Initializer))
	node.AddChildToBack(f.Transform(tree.Collection))
	node.AddChildToBack(f.processBlockTree(tree.Body))
	return node
}

// processForOf mirrors
// TransformDispatcher.processForOf(ForOfStatementTree).
func (f *IRFactory) processForOf(tree *trees.ForOfStatementTree) *rhino.Node {
	node := f.newNode(rhino.FOR_OF)
	node.AddChildToBack(f.Transform(tree.Initializer))
	node.AddChildToBack(f.Transform(tree.Collection))
	node.AddChildToBack(f.processBlockTree(tree.Body))
	return node
}

// processForAwaitOf mirrors
// TransformDispatcher.processForAwaitOf(ForAwaitOfStatementTree).
func (f *IRFactory) processForAwaitOf(tree *trees.ForAwaitOfStatementTree) *rhino.Node {
	node := f.newNode(rhino.FOR_AWAIT_OF)
	node.AddChildToBack(f.Transform(tree.Initializer))
	node.AddChildToBack(f.Transform(tree.Collection))
	node.AddChildToBack(f.processBlockTree(tree.Body))
	return node
}

// processSwitchStatement mirrors
// TransformDispatcher.processSwitchStatement(SwitchStatementTree).
func (f *IRFactory) processSwitchStatement(tree *trees.SwitchStatementTree) *rhino.Node {
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

// processSwitchCase mirrors
// TransformDispatcher.processSwitchCase(CaseClauseTree).
func (f *IRFactory) processSwitchCase(tree *trees.CaseClauseTree) *rhino.Node {
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

// processSwitchDefault mirrors
// TransformDispatcher.processSwitchDefault(DefaultClauseTree).
func (f *IRFactory) processSwitchDefault(tree *trees.DefaultClauseTree) *rhino.Node {
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

// processTryStatement mirrors
// TransformDispatcher.processTryStatement(TryStatementTree).
func (f *IRFactory) processTryStatement(tree *trees.TryStatementTree) *rhino.Node {
	node := f.newNode(rhino.TRY)
	node.AddChildToBack(f.processBlockTree(tree.Body))
	block := f.newNode(rhino.BLOCK)
	lineSet := false
	if tree.CatchBlock != nil {
		f.setSourceInfo(block, tree.CatchBlock)
		lineSet = true
		block.AddChildToBack(f.Transform(tree.CatchBlock))
	}
	node.AddChildToBack(block)
	if tree.FinallyBlock != nil {
		node.AddChildToBack(f.processBlockTree(tree.FinallyBlock))
	}
	// Mirror Java: if no catch clause, set the block's position from finally.
	if !lineSet && tree.FinallyBlock != nil {
		f.setSourceInfo(block, tree.FinallyBlock)
	}
	return node
}

// processCatchClause mirrors
// TransformDispatcher.processCatchClause(CatchTree).
func (f *IRFactory) processCatchClause(tree *trees.CatchTree) *rhino.Node {
	node := f.newNode(rhino.CATCH)
	node.AddChildToBack(f.Transform(tree.Exception))
	node.AddChildToBack(f.processBlockTree(tree.CatchBody))
	return node
}

// processFinally mirrors
// TransformDispatcher.processFinally(FinallyTree).
// Java's processFinally calls processBlock on the block field.
func (f *IRFactory) processFinally(tree *trees.FinallyTree) *rhino.Node {
	return f.processBlockTree(tree.Block)
}

// ---- declarations + functions + classes ---------------------------

// processFunction mirrors
// TransformDispatcher.processFunction(FunctionDeclarationTree).
func (f *IRFactory) processFunction(tree *trees.FunctionDeclarationTree) *rhino.Node {
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

// processFormalParameterList mirrors
// TransformDispatcher.processFormalParameterList(FormalParameterListTree).
func (f *IRFactory) processFormalParameterList(tree *trees.FormalParameterListTree) *rhino.Node {
	params := f.newNode(rhino.PARAM_LIST)
	for _, param := range tree.Parameters {
		params.AddChildToBack(f.Transform(param))
	}
	return params
}

// processDefaultParameter mirrors
// TransformDispatcher.processDefaultParameter(DefaultParameterTree).
func (f *IRFactory) processDefaultParameter(tree *trees.DefaultParameterTree) *rhino.Node {
	targetNode := f.Transform(tree.Lhs)
	defaultValue := f.Transform(tree.DefaultValue)
	node := f.newNode(rhino.DEFAULT_VALUE)
	node.AddChildToBack(targetNode)
	node.AddChildToBack(defaultValue)
	f.setSourceInfo(node, &tree.ParseTree)
	return node
}

// processClassDeclaration mirrors
// TransformDispatcher.processClassDeclaration(ClassDeclarationTree).
func (f *IRFactory) processClassDeclaration(tree *trees.ClassDeclarationTree) *rhino.Node {
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

// processField mirrors
// TransformDispatcher.processField(FieldDeclarationTree).
func (f *IRFactory) processField(tree *trees.FieldDeclarationTree) *rhino.Node {
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

// processVariableStatement mirrors
// TransformDispatcher.processVariableStatement(VariableStatementTree).
// Java delegates directly to the inner VariableDeclarationListTree.
func (f *IRFactory) processVariableStatement(tree *trees.VariableStatementTree) *rhino.Node {
	if tree.Declarations != nil {
		return f.Transform(&tree.Declarations.ParseTree)
	}
	return f.newNode(rhino.EMPTY)
}

// processVariableDeclarationList mirrors
// TransformDispatcher.processVariableDeclarationList(VariableDeclarationListTree).
func (f *IRFactory) processVariableDeclarationList(tree *trees.VariableDeclarationListTree) *rhino.Node {
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

// processVariableDeclaration mirrors
// TransformDispatcher.processVariableDeclaration(VariableDeclarationTree).
func (f *IRFactory) processVariableDeclaration(tree *trees.VariableDeclarationTree) *rhino.Node {
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

// processImportDecl mirrors
// TransformDispatcher.processImportDecl(ImportDeclarationTree).
func (f *IRFactory) processImportDecl(tree *trees.ImportDeclarationTree) *rhino.Node {
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

// processImportSpec mirrors
// TransformDispatcher.processImportSpec(ImportSpecifierTree).
func (f *IRFactory) processImportSpec(tree *trees.ImportSpecifierTree) *rhino.Node {
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

// processExportDecl mirrors
// TransformDispatcher.processExportDecl(ExportDeclarationTree).
func (f *IRFactory) processExportDecl(tree *trees.ExportDeclarationTree) *rhino.Node {
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

// processExportSpec mirrors
// TransformDispatcher.processExportSpec(ExportSpecifierTree).
func (f *IRFactory) processExportSpec(tree *trees.ExportSpecifierTree) *rhino.Node {
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

// processName mirrors
// TransformDispatcher.processName(IdentifierExpressionTree, IdentifierType).
func (f *IRFactory) processName(tree *trees.IdentifierExpressionTree) *rhino.Node {
	if tree.IdentifierToken == nil {
		return f.newStringNode(rhino.NAME, "__missing_name__")
	}
	return f.newStringNode(rhino.NAME, tree.IdentifierToken.Value)
}

// processLiteralExpression processes on the literal token type.
// Mirrors the per-type literal processing in TransformDispatcher.
func (f *IRFactory) processLiteralExpression(tree *trees.LiteralExpressionTree) *rhino.Node {
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
		return f.processRegExpLiteral("")
	default:
		return f.newNode(rhino.EMPTY)
	}
}

// processThisExpression mirrors
// TransformDispatcher.processThisExpression(ThisExpressionTree).
func (f *IRFactory) processThisExpression() *rhino.Node {
	return f.newNode(rhino.THIS)
}

// processSuper mirrors Java super expression → Token.SUPER.
func (f *IRFactory) processSuper() *rhino.Node {
	return f.newNode(rhino.SUPER)
}

// ---- unary / binary / update ---------------------------------------

// processUnaryExpression mirrors
// TransformDispatcher.processUnaryExpression(UnaryExpressionTree).
func (f *IRFactory) processUnaryExpression(tree *trees.UnaryExpressionTree) *rhino.Node {
	if tree.Operator == nil {
		return f.newNode(rhino.EMPTY)
	}
	tok := transformUnaryTokenType(tree.Operator.Type)
	operand := f.Transform(tree.Operand)
	node := f.newNode(tok)
	node.AddChildToBack(operand)
	return node
}

// processUpdateExpression mirrors
// TransformDispatcher.processUpdateExpression(UpdateExpressionTree).
func (f *IRFactory) processUpdateExpression(tree *trees.UpdateExpressionTree) *rhino.Node {
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

// processBinaryExpression mirrors
// TransformDispatcher.processBinaryExpression(BinaryOperatorTree).
func (f *IRFactory) processBinaryExpression(tree *trees.BinaryOperatorTree) *rhino.Node {
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

// processConditionalExpression mirrors
// TransformDispatcher.processConditionalExpression(ConditionalExpressionTree).
func (f *IRFactory) processConditionalExpression(tree *trees.ConditionalExpressionTree) *rhino.Node {
	node := f.newNode(rhino.HOOK)
	node.AddChildToBack(f.Transform(tree.Condition))
	node.AddChildToBack(f.Transform(tree.Left))
	node.AddChildToBack(f.Transform(tree.Right))
	return node
}

// processCommaExpression mirrors
// TransformDispatcher.processCommaExpression(CommaExpressionTree).
func (f *IRFactory) processCommaExpression(tree *trees.CommaExpressionTree) *rhino.Node {
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

// processParenthesizedExpression mirrors
// TransformDispatcher.processParenthesizedExpression(ParenExpressionTree).
func (f *IRFactory) processParenthesizedExpression(tree *trees.ParenExpressionTree) *rhino.Node {
	expr := f.Transform(tree.Expression)
	if expr != nil {
		expr.SetIsParenthesized(true)
	}
	return expr
}

// ---- calls / member access / new -----------------------------------

// processFunctionCall mirrors
// TransformDispatcher.processFunctionCall(CallExpressionTree).
func (f *IRFactory) processFunctionCall(tree *trees.CallExpressionTree) *rhino.Node {
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

// processNewExpression mirrors
// TransformDispatcher.processNewExpression(NewExpressionTree).
func (f *IRFactory) processNewExpression(tree *trees.NewExpressionTree) *rhino.Node {
	node := f.newNode(rhino.NEW)
	node.AddChildToBack(f.Transform(tree.Operand))
	if tree.Arguments != nil {
		for _, arg := range tree.Arguments.Arguments {
			node.AddChildToBack(f.Transform(arg))
		}
	}
	return node
}

// processPropertyGet mirrors
// TransformDispatcher.processPropertyGet(MemberExpressionTree).
func (f *IRFactory) processPropertyGet(tree *trees.MemberExpressionTree) *rhino.Node {
	leftChild := f.Transform(tree.Operand)
	if tree.MemberName == nil {
		return leftChild
	}
	getProp := f.newStringNode(rhino.GETPROP, tree.MemberName.Value)
	getProp.AddChildToBack(leftChild)
	return getProp
}

// processElementGet mirrors
// TransformDispatcher.processElementGet(MemberLookupExpressionTree).
func (f *IRFactory) processElementGet(tree *trees.MemberLookupExpressionTree) *rhino.Node {
	node := f.newNode(rhino.GETELEM)
	node.AddChildToBack(f.Transform(tree.Operand))
	node.AddChildToBack(f.Transform(tree.MemberExpression))
	return node
}

// ---- structured literals -------------------------------------------

// processArrayLiteral mirrors
// TransformDispatcher.processArrayLiteral(ArrayLiteralExpressionTree).
func (f *IRFactory) processArrayLiteral(tree *trees.ArrayLiteralExpressionTree) *rhino.Node {
	node := f.newNode(rhino.ARRAYLIT)
	for _, el := range tree.Elements {
		node.AddChildToBack(f.Transform(el))
	}
	return node
}

// processObjectLiteral mirrors
// TransformDispatcher.processObjectLiteral(ObjectLiteralExpressionTree).
func (f *IRFactory) processObjectLiteral(tree *trees.ObjectLiteralExpressionTree) *rhino.Node {
	node := f.newNode(rhino.OBJECTLIT)
	for _, el := range tree.PropertyNameAndValues {
		key := f.Transform(el)
		if key != nil {
			node.AddChildToBack(key)
		}
	}
	return node
}

// processPropertyNameAssignment mirrors
// TransformDispatcher.processPropertyNameAssignment(PropertyNameAssignmentTree).
func (f *IRFactory) processPropertyNameAssignment(tree *trees.PropertyNameAssignmentTree) *rhino.Node {
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

// processGetAccessor mirrors
// TransformDispatcher.processGetAccessor(GetAccessorTree).
func (f *IRFactory) processGetAccessor(tree *trees.GetAccessorTree) *rhino.Node {
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

// processSetAccessor mirrors
// TransformDispatcher.processSetAccessor(SetAccessorTree).
func (f *IRFactory) processSetAccessor(tree *trees.SetAccessorTree) *rhino.Node {
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

// processTemplateLiteral mirrors
// TransformDispatcher.processTemplateLiteral(TemplateLiteralExpressionTree).
func (f *IRFactory) processTemplateLiteral(tree *trees.TemplateLiteralExpressionTree) *rhino.Node {
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

// processTemplateLiteralPortion mirrors
// TransformDispatcher.processTemplateLiteralPortion(TemplateLiteralPortionTree).
// Note: trees.Token has no Value field in the stub; cooked/raw are empty stubs.
func (f *IRFactory) processTemplateLiteralPortion(tree *trees.TemplateLiteralPortionTree) *rhino.Node {
	return rhino.NewTemplateLitString("", "")
}

// processTemplateSubstitution mirrors
// TransformDispatcher.processTemplateSubstitution(TemplateSubstitutionTree).
func (f *IRFactory) processTemplateSubstitution(tree *trees.TemplateSubstitutionTree) *rhino.Node {
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

// processRegExpLiteral builds a REGEXP node from the raw literal string.
// Mirrors TransformDispatcher.processRegExpLiteral.
func (f *IRFactory) processRegExpLiteral(raw string) *rhino.Node {
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

// Go-only: token-based overload of setSourceInfo; Java's setSourceInfo takes a ParseTree.
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

// Go-only: Go-level numeric literal parser; Java uses removeNumericSeparators + Double.parseDouble.
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

// Go-only: converts float64 to JS number string; Java uses DToA.numberToString.
func floatToString(v float64) string {
	return fmt.Sprintf("%g", v)
}
