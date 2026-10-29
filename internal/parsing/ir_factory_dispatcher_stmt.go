// Port of com/google/javascript/jscomp/parsing/IRFactory.java @ closure-compiler v20261021

/*
 * Copyright 2013 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package parsing

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// IRFactory_TransformDispatcher dispatches each ParseTree to its process
// method.
//
// Java: private class IRFactory.TransformDispatcher (an inner class)
type IRFactory_TransformDispatcher struct {
	// Go-only: language-adapter — the enclosing IRFactory instance of the
	// inner class, embedded so that its members resolve as they do in Java.
	*IRFactory

	// Tracks whether a class is currently in scope. Presently used to detect
	// valid usage of private fields.
	classScope *IRFactory_TransformDispatcher_ScopeTracker

	// Used to detect an IdentifierExpression in the following context:
	// Valid: `class A { #x; static isA(o) { return #x in o; } }`
	// Invalid: `class A { #x; } function isA(o) { return #x in o; } }`
	privateIdLhsOfInScope *IRFactory_TransformDispatcher_ScopeTracker
}

// Java: IRFactory.TransformDispatcher()
func NewIRFactory_TransformDispatcher(outer *IRFactory) *IRFactory_TransformDispatcher {
	d := &IRFactory_TransformDispatcher{IRFactory: outer}
	d.classScope = NewIRFactory_TransformDispatcher_ScopeTracker()
	d.privateIdLhsOfInScope = NewIRFactory_TransformDispatcher_ScopeTracker()
	return d
}

// IRFactory_TransformDispatcher_ScopeTracker tracks whether a certain
// language feature is currently in scope. Use in conjunction with
// try-with-resources to auto-handle decrementing.
//
// Java: private static class IRFactory.TransformDispatcher.ScopeTracker
type IRFactory_TransformDispatcher_ScopeTracker struct {
	usageCount int
}

// Java: IRFactory.TransformDispatcher.ScopeTracker()
func NewIRFactory_TransformDispatcher_ScopeTracker() *IRFactory_TransformDispatcher_ScopeTracker {
	return &IRFactory_TransformDispatcher_ScopeTracker{usageCount: 0}
}

// Java: IRFactory.TransformDispatcher.ScopeTracker.inScope()
func (s *IRFactory_TransformDispatcher_ScopeTracker) InScope() bool {
	return s.usageCount > 0
}

// Java: IRFactory.TransformDispatcher.ScopeTracker.increment()
func (s *IRFactory_TransformDispatcher_ScopeTracker) Increment() irFactory_TransformDispatcher_ScopeTracker_autoDecrement {
	return s.MaybeIncrement(true)
}

// MaybeIncrement increments the usage count if condition holds and returns
// the AutoDecrement that undoes it.
//
// Go-only: language-adapter — returns the anonymous AutoDecrement class by
// value rather than as the AutoDecrement interface, so that it stays on the
// stack (the JVM's escape analysis removes the allocation; Go's would not
// once it is boxed in an interface).
//
// Java: IRFactory.TransformDispatcher.ScopeTracker.maybeIncrement(boolean)
func (s *IRFactory_TransformDispatcher_ScopeTracker) MaybeIncrement(condition bool) irFactory_TransformDispatcher_ScopeTracker_autoDecrement {
	if condition {
		jutil.CheckState(s.usageCount >= 0)
		s.usageCount++
	}
	return irFactory_TransformDispatcher_ScopeTracker_autoDecrement{tracker: s, condition: condition, closed: false}
}

// irFactory_TransformDispatcher_ScopeTracker_autoDecrement is the anonymous
// AutoDecrement class of ScopeTracker.maybeIncrement.
type irFactory_TransformDispatcher_ScopeTracker_autoDecrement struct {
	tracker   *IRFactory_TransformDispatcher_ScopeTracker
	condition bool
	closed    bool
}

// Java: IRFactory.TransformDispatcher.ScopeTracker.maybeIncrement(boolean).new AutoDecrement() {...}.close()
func (a *irFactory_TransformDispatcher_ScopeTracker_autoDecrement) Close() {
	jutil.CheckState(!a.closed)
	if a.condition {
		a.tracker.usageCount--
		jutil.CheckState(a.tracker.usageCount >= 0)
	}
	a.closed = true
}

// IRFactory_TransformDispatcher_IdentifierType indicates the valid values
// usages for an identifier.
//
// Java: private enum IRFactory.TransformDispatcher.IdentifierType
type IRFactory_TransformDispatcher_IdentifierType int32

const (
	// The identifier can never be private.
	IRFactory_TransformDispatcher_IdentifierType_STANDARD IRFactory_TransformDispatcher_IdentifierType = iota
	// The identifier can be private if it is in a class scope.
	IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE
)

var irFactory_TransformDispatcher_IdentifierType_names = [...]string{"STANDARD", "CAN_BE_PRIVATE"}

// String is Enum.name().
func (t IRFactory_TransformDispatcher_IdentifierType) String() string {
	return irFactory_TransformDispatcher_IdentifierType_names[t]
}

// NewStringNodeFromIdentifier creates a StringNode using the value of an
// IdentifierToken. Reports an error if the identifier value is private and
// the identifier type is STANDARD or the identifier is not within the scope
// of a class.
//

// Java: IRFactory.TransformDispatcher.processBlock(BlockTree)
func (d *IRFactory_TransformDispatcher) ProcessBlock(blockNode *parser.BlockTree) rhino.Node {
	node := d.NewNode(rhino.Token_BLOCK)
	for i, child := range blockNode.Statements {
		childNode := d.Transform(child)
		node.AddChildToBack(childNode)
		d.AttachPossibleTrailingComment(childNode, child.GetEnd())
		// Go-only: performance — nothing reads a statement's parse tree once it
		// is transformed: dropping it lets the GC free it while the rest of the
		// file is transformed, instead of holding a large file's whole parse
		// tree next to its whole AST.
		blockNode.Statements[i] = nil
	}

	if endOff := blockNode.GetEnd(); d.hasAnyPendingCommentBeforeOffset(endOff.Offset) {
		lastComment := d.ParseNonJSDocCommentAt(d.linePos(endOff), false)
		d.AddExtraTrailingComment(node, lastComment)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processBreakStatement(BreakStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessBreakStatement(statementNode *parser.BreakStatementTree) rhino.Node {
	node := d.NewNode(rhino.Token_BREAK)
	if statementNode.HasLabel() {
		labelName := d.TransformLabelName(statementNode.Name)
		node.AddChildToBack(labelName)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.transformLabelName(IdentifierToken)
func (d *IRFactory_TransformDispatcher) TransformLabelName(token parser.Token) rhino.Node {
	checkTokenNotNull(token)
	label := d.MaybeAddNonJsDocComment(
		d.NewStringNodeFromIdentifier(rhino.Token_LABEL_NAME, IRFactory_TransformDispatcher_IdentifierType_STANDARD, token),
		token.GetStart())
	d.SetSourceInfoFromToken(label, token)
	return label
}

// Java: IRFactory.TransformDispatcher.processContinueStatement(ContinueStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessContinueStatement(statementNode *parser.ContinueStatementTree) rhino.Node {
	node := d.NewNode(rhino.Token_CONTINUE)
	if statementNode.HasLabel() {
		labelName := d.TransformLabelName(statementNode.Name)
		node.AddChildToBack(labelName)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processDoLoop(DoWhileStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessDoLoop(loopNode *parser.DoWhileStatementTree) rhino.Node {
	return d.NewNodeWith2Children(rhino.Token_DO, d.TransformBlock(loopNode.Body), d.Transform(loopNode.Condition))
}

// Java: IRFactory.TransformDispatcher.processEmptyStatement(EmptyStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessEmptyStatement(unused *parser.EmptyStatementTree) rhino.Node {
	return d.NewNode(rhino.Token_EMPTY)
}

// Java: IRFactory.TransformDispatcher.processExpressionStatement(ExpressionStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessExpressionStatement(statementNode *parser.ExpressionStatementTree) rhino.Node {
	node := d.NewNode(rhino.Token_EXPR_RESULT)
	node.AddChildToBack(d.Transform(statementNode.Expression))
	return node
}

// Java: IRFactory.TransformDispatcher.processForInLoop(ForInStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessForInLoop(loopNode *parser.ForInStatementTree) rhino.Node {
	// TODO(bradfordcsmith): Rename initializer to something more intuitive
	// like "lhs"
	initializer := d.Transform(loopNode.Initializer)
	return d.NewNodeWith3Children(
		rhino.Token_FOR_IN, initializer, d.Transform(loopNode.Collection), d.TransformBlock(loopNode.Body))
}

// Java: IRFactory.TransformDispatcher.processForOf(ForOfStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessForOf(loopNode *parser.ForOfStatementTree) rhino.Node {
	d.MaybeWarnForFeature(loopNode, parser.FeatureSet_Feature_FOR_OF)
	initializer := d.Transform(loopNode.Initializer)
	return d.NewNodeWith3Children(
		rhino.Token_FOR_OF, initializer, d.Transform(loopNode.Collection), d.TransformBlock(loopNode.Body))
}

// Java: IRFactory.TransformDispatcher.processForAwaitOf(ForAwaitOfStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessForAwaitOf(loopNode *parser.ForAwaitOfStatementTree) rhino.Node {
	d.MaybeWarnForFeature(loopNode, parser.FeatureSet_Feature_FOR_AWAIT_OF)
	initializer := d.Transform(loopNode.Initializer)
	return d.NewNodeWith3Children(
		rhino.Token_FOR_AWAIT_OF,
		initializer,
		d.Transform(loopNode.Collection),
		d.TransformBlock(loopNode.Body))
}

// Java: IRFactory.TransformDispatcher.processForLoop(ForStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessForLoop(loopNode *parser.ForStatementTree) rhino.Node {
	node := d.NewNodeWith3Children(
		rhino.Token_FOR,
		d.TransformOrEmpty(loopNode.Initializer, loopNode),
		d.TransformOrEmpty(loopNode.Condition, loopNode),
		d.TransformOrEmpty(loopNode.Increment, loopNode))
	node.AddChildToBack(d.TransformBlock(loopNode.Body))
	return node
}

// Java: IRFactory.TransformDispatcher.transformOrEmpty(ParseTree, ParseTree)
func (d *IRFactory_TransformDispatcher) TransformOrEmpty(tree parser.ParseTree, parent parser.ParseTree) rhino.Node {
	if tree == nil {
		n := d.NewNode(rhino.Token_EMPTY)
		d.SetSourceInfo(n, parent)
		return n
	}
	return d.Transform(tree)
}

// Java: IRFactory.TransformDispatcher.transformOrEmpty(IdentifierToken, ParseTree)
func (d *IRFactory_TransformDispatcher) TransformOrEmptyIdentifier(token parser.Token, parent parser.ParseTree) rhino.Node {
	if token.IsNull() {
		n := d.NewNode(rhino.Token_EMPTY)
		d.SetSourceInfo(n, parent)
		return n
	}
	return d.ProcessNameOfToken(token, rhino.Token_NAME, IRFactory_TransformDispatcher_IdentifierType_STANDARD)
}

// Java: IRFactory.TransformDispatcher.processFunction(FunctionDeclarationTree)
func (d *IRFactory_TransformDispatcher) ProcessFunction(functionTree *parser.FunctionDeclarationTree) rhino.Node {
	isDeclaration := functionTree.Kind == parser.FunctionDeclarationTree_Kind_DECLARATION
	isMember := functionTree.Kind == parser.FunctionDeclarationTree_Kind_MEMBER
	isArrow := functionTree.Kind == parser.FunctionDeclarationTree_Kind_ARROW
	isAsync := functionTree.IsAsync
	isGenerator := functionTree.IsGenerator
	isSignature := functionTree.FunctionBody.GetParseTreeBase().Type == parser.ParseTreeType_EMPTY_STATEMENT

	identifierType := IRFactory_TransformDispatcher_IdentifierType_STANDARD
	if functionTree.IsClassMember {
		identifierType = IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE
	}

	if isGenerator {
		d.MaybeWarnForFeature(functionTree, parser.FeatureSet_Feature_GENERATORS)
	}

	if isMember {
		d.MaybeWarnForFeature(functionTree, parser.FeatureSet_Feature_MEMBER_DECLARATIONS)
	}

	if isArrow {
		d.MaybeWarnForFeature(functionTree, parser.FeatureSet_Feature_ARROW_FUNCTIONS)
	}

	if isAsync {
		d.MaybeWarnForFeature(functionTree, parser.FeatureSet_Feature_ASYNC_FUNCTIONS)
	}

	if isGenerator && isAsync {
		d.MaybeWarnForFeature(functionTree, parser.FeatureSet_Feature_ASYNC_GENERATORS)
	}

	name := functionTree.Name
	var newName rhino.Node
	if name.IsPresent() {
		newName = d.ProcessNameWithInlineCommentsOfToken(name, identifierType)
	} else {
		if isDeclaration || isMember {
			d.errorReporter.Error(
				"unnamed function statement", d.sourceName, d.Lineno(functionTree), d.Charno(functionTree))

			// Return the bare minimum to put the AST in a valid state.
			newName = d.CreateMissingNameNode()
		} else {
			newName = d.NewStringNodeWithToken(rhino.Token_NAME, "")
		}

		// Old Rhino tagged the empty name node with the line number of the
		// declaration.
		d.SetSourceInfo(newName, functionTree)
	}

	node := d.NewNode(rhino.Token_FUNCTION)
	if isMember {
		newName.SetString("")
	}
	node.AddChildToBack(newName)

	node.AddChildToBack(d.Transform(functionTree.FormalParameterList))

	bodyNode := d.Transform(functionTree.FunctionBody)
	if !isArrow && !isSignature && !bodyNode.IsBlock() {
		// When in "keep going" mode the parser tries to parse some constructs
		// the compiler doesn't support, repair it here.
		jutil.CheckState(d.config.RunMode() == Config_RunMode_KEEP_GOING)
		bodyNode = rhino.IR_Block()
	}
	d.ParseDirectives(bodyNode)
	node.AddChildToBack(bodyNode)

	node.SetIsGeneratorFunction(isGenerator)
	node.SetIsArrowFunction(isArrow)
	node.SetIsAsyncFunction(isAsync)
	node.PutBooleanProp(rhino.Node_OPT_ES6_TYPED, functionTree.IsOptional)

	d.AttachPossibleTrailingComment(node, functionTree.GetEnd())

	var result rhino.Node

	if isMember {
		d.SetSourceInfo(node, functionTree)
		member := d.NewStringNodeFromIdentifier(rhino.Token_MEMBER_FUNCTION_DEF, identifierType, name)
		member.AddChildToBack(node)
		member.SetStaticMember(functionTree.IsStatic)
		// The source info should only include the identifier, not the entire
		// function expression
		d.SetSourceInfoFromToken(member, name)
		result = member
	} else {
		result = node
	}

	return result
}

// Java: IRFactory.TransformDispatcher.processIfStatement(IfStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessIfStatement(statementNode *parser.IfStatementTree) rhino.Node {
	node := d.NewNode(rhino.Token_IF)
	node.AddChildToBack(d.Transform(statementNode.Condition))
	node.AddChildToBack(d.TransformBlock(statementNode.IfClause))
	if statementNode.ElseClause != nil {
		node.AddChildToBack(d.TransformBlock(statementNode.ElseClause))
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processDebuggerStatement(DebuggerStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessDebuggerStatement(unused *parser.DebuggerStatementTree) rhino.Node {
	return d.NewNode(rhino.Token_DEBUGGER)
}

// Java: IRFactory.TransformDispatcher.processLabeledStatement(LabelledStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessLabeledStatement(labelTree *parser.LabelledStatementTree) rhino.Node {
	statement := d.Transform(labelTree.Statement)
	if statement.IsFunction() ||
		statement.IsClass() ||
		statement.IsLet() ||
		statement.IsConst() {
		d.errorReporter.Error(
			"Lexical declarations are only allowed at top level or inside a block.",
			d.sourceName,
			d.Lineno(labelTree),
			d.Charno(labelTree))
		return statement // drop the LABEL node so that the resulting AST is valid
	}
	return d.NewNodeWith2Children(rhino.Token_LABEL, d.TransformLabelName(labelTree.Name), statement)
}

// Java: IRFactory.TransformDispatcher.processReturnStatement(ReturnStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessReturnStatement(statementNode *parser.ReturnStatementTree) rhino.Node {
	node := d.NewNode(rhino.Token_RETURN)
	if statementNode.Expression != nil {
		node.AddChildToBack(d.Transform(statementNode.Expression))
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processSwitchCase(CaseClauseTree)
func (d *IRFactory_TransformDispatcher) ProcessSwitchCase(caseNode *parser.CaseClauseTree) rhino.Node {
	expr := caseNode.Expression
	node := d.NewNodeWithChild(rhino.Token_CASE, d.Transform(expr))
	block := d.NewNode(rhino.Token_BLOCK)
	block.SetIsAddedBlock(true)
	d.SetSourceInfo(block, caseNode)
	if caseNode.Statements != nil {
		for _, child := range caseNode.Statements {
			block.AddChildToBack(d.Transform(child))
		}
	}
	node.AddChildToBack(block)
	return node
}

// Java: IRFactory.TransformDispatcher.processSwitchDefault(DefaultClauseTree)
func (d *IRFactory_TransformDispatcher) ProcessSwitchDefault(caseNode *parser.DefaultClauseTree) rhino.Node {
	node := d.NewNode(rhino.Token_DEFAULT_CASE)
	block := d.NewNode(rhino.Token_BLOCK)
	block.SetIsAddedBlock(true)
	d.SetSourceInfo(block, caseNode)
	if caseNode.Statements != nil {
		for _, child := range caseNode.Statements {
			block.AddChildToBack(d.Transform(child))
		}
	}
	node.AddChildToBack(block)
	return node
}

// Java: IRFactory.TransformDispatcher.processSwitchStatement(SwitchStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessSwitchStatement(statementNode *parser.SwitchStatementTree) rhino.Node {
	node := d.NewNodeWithChild(rhino.Token_SWITCH, d.Transform(statementNode.Expression))
	switchBody := d.NewNode(rhino.Token_SWITCH_BODY)
	d.SetSourceInfo(switchBody, statementNode)
	node.AddChildToBack(switchBody)
	for _, child := range statementNode.CaseClauses {
		switchBody.AddChildToBack(d.Transform(child))
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processThrowStatement(ThrowStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessThrowStatement(statementNode *parser.ThrowStatementTree) rhino.Node {
	return d.NewNodeWithChild(rhino.Token_THROW, d.Transform(statementNode.Value))
}

// Java: IRFactory.TransformDispatcher.processTryStatement(TryStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessTryStatement(statementNode *parser.TryStatementTree) rhino.Node {
	node := d.NewNodeWithChild(rhino.Token_TRY, d.TransformBlock(statementNode.Body))
	block := d.NewNode(rhino.Token_BLOCK)
	lineSet := false

	cc := statementNode.CatchBlock
	if cc != nil {
		// Mark the enclosing block at the same line as the first catch
		// clause.
		d.SetSourceInfo(block, cc)
		lineSet = true
		block.AddChildToBack(d.Transform(cc))
	}

	node.AddChildToBack(block)

	finallyBlock := statementNode.FinallyBlock
	if finallyBlock != nil {
		node.AddChildToBack(d.TransformBlock(finallyBlock))
	}

	// If we didn't set the line on the catch clause, then
	// we've got an empty catch clause.  Set its line to be the same
	// as the finally block (to match Old Rhino's behavior.)
	if !lineSet && (finallyBlock != nil) {
		d.SetSourceInfo(block, finallyBlock)
	}

	return node
}

// Java: IRFactory.TransformDispatcher.processCatchClause(CatchTree)
func (d *IRFactory_TransformDispatcher) ProcessCatchClause(clauseNode *parser.CatchTree) rhino.Node {
	if clauseNode.Exception.GetParseTreeBase().Type == parser.ParseTreeType_EMPTY_STATEMENT {
		d.MaybeWarnForFeature(clauseNode, parser.FeatureSet_Feature_OPTIONAL_CATCH_BINDING)
	}

	return d.NewNodeWith2Children(
		rhino.Token_CATCH, d.Transform(clauseNode.Exception), d.TransformBlock(clauseNode.CatchBody))
}

// Java: IRFactory.TransformDispatcher.processFinally(FinallyTree)
func (d *IRFactory_TransformDispatcher) ProcessFinally(finallyNode *parser.FinallyTree) rhino.Node {
	return d.TransformBlock(finallyNode.Block)
}

// Java: IRFactory.TransformDispatcher.processVariableStatement(VariableStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessVariableStatement(stmt *parser.VariableStatementTree) rhino.Node {
	// TODO(moz): Figure out why we still need the special handling
	return d.transformDispatcher.Process(stmt.Declarations)
}

// Java: IRFactory.TransformDispatcher.processVariableDeclarationList(VariableDeclarationListTree)
func (d *IRFactory_TransformDispatcher) ProcessVariableDeclarationList(decl *parser.VariableDeclarationListTree) rhino.Node {
	var declType rhino.Token
	switch decl.DeclarationType {
	case parser.TokenType_CONST:
		d.MaybeWarnForFeature(decl, parser.FeatureSet_Feature_CONST_DECLARATIONS)
		declType = rhino.Token_CONST
	case parser.TokenType_LET:
		d.MaybeWarnForFeature(decl, parser.FeatureSet_Feature_LET_DECLARATIONS)
		declType = rhino.Token_LET
	case parser.TokenType_VAR:
		declType = rhino.Token_VAR
	default:
		panic(jutil.NewIllegalStateException())
	}

	node := d.NewNode(declType)
	for _, child := range decl.Declarations {
		node.AddChildToBack(d.TransformNodeWithInlineComments(child))
	}

	d.AttachPossibleTrailingComment(node, decl.GetEnd())
	return node
}

// Java: IRFactory.TransformDispatcher.processVariableDeclaration(VariableDeclarationTree)
func (d *IRFactory_TransformDispatcher) ProcessVariableDeclaration(decl *parser.VariableDeclarationTree) rhino.Node {
	node := d.TransformNodeWithInlineComments(decl.Lvalue)
	lhs := node
	if node.IsDestructuringPattern() {
		lhs = d.NewNodeWithChild(rhino.Token_DESTRUCTURING_LHS, node)
	}
	if decl.Initializer != nil {
		initializer := d.Transform(decl.Initializer)
		lhs.AddChildToBack(initializer)
		d.SetLength(lhs, decl.Location.Start, decl.Location.End)
	}
	return lhs
}

// Java: IRFactory.TransformDispatcher.processWhileLoop(WhileStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessWhileLoop(stmt *parser.WhileStatementTree) rhino.Node {
	return d.NewNodeWith2Children(rhino.Token_WHILE, d.Transform(stmt.Condition), d.TransformBlock(stmt.Body))
}

// Java: IRFactory.TransformDispatcher.processWithStatement(WithStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessWithStatement(stmt *parser.WithStatementTree) rhino.Node {
	return d.NewNodeWith2Children(rhino.Token_WITH, d.Transform(stmt.Expression), d.TransformBlock(stmt.Body))
}

// irFactory_javaClassOf is Object.getClass().toString() for a parser token:
// "class com.google.javascript.jscomp.parsing.parser.<simple name>".
func irFactory_javaClassOf(token parser.Token) string {
	return "class com.google.javascript.jscomp.parsing.parser." + token.JavaSimpleClassName()
}

// Java: IRFactory.TransformDispatcher.process(ParseTree)
func (d *IRFactory_TransformDispatcher) Process(node parser.ParseTree) rhino.Node {
	switch node.GetParseTreeBase().Type {
	case parser.ParseTreeType_BINARY_OPERATOR:
		return d.ProcessBinaryExpression(node.(*parser.BinaryOperatorTree))
	case parser.ParseTreeType_ARRAY_LITERAL_EXPRESSION:
		return d.ProcessArrayLiteral(node.(*parser.ArrayLiteralExpressionTree))
	case parser.ParseTreeType_TEMPLATE_LITERAL_EXPRESSION:
		return d.ProcessTemplateLiteral(node.(*parser.TemplateLiteralExpressionTree))
	case parser.ParseTreeType_TEMPLATE_LITERAL_PORTION:
		return d.ProcessTemplateLiteralPortion(node.(*parser.TemplateLiteralPortionTree))
	case parser.ParseTreeType_TEMPLATE_SUBSTITUTION:
		return d.ProcessTemplateSubstitution(node.(*parser.TemplateSubstitutionTree))
	case parser.ParseTreeType_UNARY_EXPRESSION:
		return d.ProcessUnaryExpression(node.(*parser.UnaryExpressionTree))
	case parser.ParseTreeType_BLOCK:
		return d.ProcessBlock(node.(*parser.BlockTree))
	case parser.ParseTreeType_BREAK_STATEMENT:
		return d.ProcessBreakStatement(node.(*parser.BreakStatementTree))
	case parser.ParseTreeType_CALL_EXPRESSION:
		return d.ProcessFunctionCall(node.(*parser.CallExpressionTree))
	case parser.ParseTreeType_OPT_CHAIN_CALL_EXPRESSION:
		return d.ProcessOptChainFunctionCall(node.(*parser.OptChainCallExpressionTree))
	case parser.ParseTreeType_CASE_CLAUSE:
		return d.ProcessSwitchCase(node.(*parser.CaseClauseTree))
	case parser.ParseTreeType_DEFAULT_CLAUSE:
		return d.ProcessSwitchDefault(node.(*parser.DefaultClauseTree))
	case parser.ParseTreeType_CATCH:
		return d.ProcessCatchClause(node.(*parser.CatchTree))
	case parser.ParseTreeType_CONTINUE_STATEMENT:
		return d.ProcessContinueStatement(node.(*parser.ContinueStatementTree))
	case parser.ParseTreeType_DO_WHILE_STATEMENT:
		return d.ProcessDoLoop(node.(*parser.DoWhileStatementTree))
	case parser.ParseTreeType_EMPTY_STATEMENT:
		return d.ProcessEmptyStatement(node.(*parser.EmptyStatementTree))
	case parser.ParseTreeType_EXPRESSION_STATEMENT:
		return d.ProcessExpressionStatement(node.(*parser.ExpressionStatementTree))
	case parser.ParseTreeType_DEBUGGER_STATEMENT:
		return d.ProcessDebuggerStatement(node.(*parser.DebuggerStatementTree))
	case parser.ParseTreeType_THIS_EXPRESSION:
		return d.ProcessThisExpression(node.(*parser.ThisExpressionTree))
	case parser.ParseTreeType_FOR_STATEMENT:
		return d.ProcessForLoop(node.(*parser.ForStatementTree))
	case parser.ParseTreeType_FOR_IN_STATEMENT:
		return d.ProcessForInLoop(node.(*parser.ForInStatementTree))
	case parser.ParseTreeType_FUNCTION_DECLARATION:
		return d.ProcessFunction(node.(*parser.FunctionDeclarationTree))
	case parser.ParseTreeType_MEMBER_LOOKUP_EXPRESSION:
		return d.ProcessElementGet(node.(*parser.MemberLookupExpressionTree))
	case parser.ParseTreeType_OPT_CHAIN_MEMBER_LOOKUP_EXPRESSION:
		return d.ProcessOptChainElementGet(node.(*parser.OptionalMemberLookupExpressionTree))
	case parser.ParseTreeType_MEMBER_EXPRESSION:
		return d.ProcessPropertyGet(node.(*parser.MemberExpressionTree))
	case parser.ParseTreeType_OPT_CHAIN_MEMBER_EXPRESSION:
		return d.ProcessOptChainPropertyGet(node.(*parser.OptionalMemberExpressionTree))
	case parser.ParseTreeType_CONDITIONAL_EXPRESSION:
		return d.ProcessConditionalExpression(node.(*parser.ConditionalExpressionTree))
	case parser.ParseTreeType_IF_STATEMENT:
		return d.ProcessIfStatement(node.(*parser.IfStatementTree))
	case parser.ParseTreeType_LABELLED_STATEMENT:
		return d.ProcessLabeledStatement(node.(*parser.LabelledStatementTree))
	case parser.ParseTreeType_PAREN_EXPRESSION:
		return d.ProcessParenthesizedExpression(node.(*parser.ParenExpressionTree))
	case parser.ParseTreeType_IDENTIFIER_EXPRESSION:
		identifierType := IRFactory_TransformDispatcher_IdentifierType_STANDARD
		if d.privateIdLhsOfInScope.InScope() {
			identifierType = IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE
		}
		return d.ProcessName(node.(*parser.IdentifierExpressionTree), identifierType)
	case parser.ParseTreeType_NEW_EXPRESSION:
		return d.ProcessNewExpression(node.(*parser.NewExpressionTree))
	case parser.ParseTreeType_OBJECT_LITERAL_EXPRESSION:
		return d.ProcessObjectLiteral(node.(*parser.ObjectLiteralExpressionTree))
	case parser.ParseTreeType_COMPUTED_PROPERTY_DEFINITION:
		return d.ProcessComputedPropertyDefinition(node.(*parser.ComputedPropertyDefinitionTree))
	case parser.ParseTreeType_COMPUTED_PROPERTY_GETTER:
		return d.ProcessComputedPropertyGetter(node.(*parser.ComputedPropertyGetterTree))
	case parser.ParseTreeType_COMPUTED_PROPERTY_METHOD:
		return d.ProcessComputedPropertyMethod(node.(*parser.ComputedPropertyMethodTree))
	case parser.ParseTreeType_COMPUTED_PROPERTY_SETTER:
		return d.ProcessComputedPropertySetter(node.(*parser.ComputedPropertySetterTree))
	case parser.ParseTreeType_RETURN_STATEMENT:
		return d.ProcessReturnStatement(node.(*parser.ReturnStatementTree))
	case parser.ParseTreeType_UPDATE_EXPRESSION:
		return d.ProcessUpdateExpression(node.(*parser.UpdateExpressionTree))
	case parser.ParseTreeType_PROGRAM:
		return d.ProcessAstRoot(node.(*parser.ProgramTree))
	case parser.ParseTreeType_LITERAL_EXPRESSION:
		// STRING, NUMBER, TRUE, FALSE, NULL, REGEXP
		return d.ProcessLiteralExpression(node.(*parser.LiteralExpressionTree))
	case parser.ParseTreeType_SWITCH_STATEMENT:
		return d.ProcessSwitchStatement(node.(*parser.SwitchStatementTree))
	case parser.ParseTreeType_THROW_STATEMENT:
		return d.ProcessThrowStatement(node.(*parser.ThrowStatementTree))
	case parser.ParseTreeType_TRY_STATEMENT:
		return d.ProcessTryStatement(node.(*parser.TryStatementTree))
	case parser.ParseTreeType_VARIABLE_STATEMENT:
		// var const let
		return d.ProcessVariableStatement(node.(*parser.VariableStatementTree))
	case parser.ParseTreeType_VARIABLE_DECLARATION_LIST:
		return d.ProcessVariableDeclarationList(node.(*parser.VariableDeclarationListTree))
	case parser.ParseTreeType_VARIABLE_DECLARATION:
		return d.ProcessVariableDeclaration(node.(*parser.VariableDeclarationTree))
	case parser.ParseTreeType_WHILE_STATEMENT:
		return d.ProcessWhileLoop(node.(*parser.WhileStatementTree))
	case parser.ParseTreeType_WITH_STATEMENT:
		return d.ProcessWithStatement(node.(*parser.WithStatementTree))
	case parser.ParseTreeType_COMMA_EXPRESSION:
		return d.ProcessCommaExpression(node.(*parser.CommaExpressionTree))
	case parser.ParseTreeType_NULL:
		// this is not the null literal
		return d.ProcessNull(node.(*parser.NullTree))
	case parser.ParseTreeType_FINALLY:
		return d.ProcessFinally(node.(*parser.FinallyTree))
	case parser.ParseTreeType_MISSING_PRIMARY_EXPRESSION:
		return d.ProcessMissingExpression(node.(*parser.MissingPrimaryExpressionTree))
	case parser.ParseTreeType_PROPERTY_NAME_ASSIGNMENT:
		return d.ProcessPropertyNameAssignment(node.(*parser.PropertyNameAssignmentTree))
	case parser.ParseTreeType_GET_ACCESSOR:
		return d.ProcessGetAccessor(node.(*parser.GetAccessorTree))
	case parser.ParseTreeType_SET_ACCESSOR:
		return d.ProcessSetAccessor(node.(*parser.SetAccessorTree))
	case parser.ParseTreeType_FORMAL_PARAMETER_LIST:
		return d.ProcessFormalParameterList(node.(*parser.FormalParameterListTree))
	case parser.ParseTreeType_CLASS_DECLARATION:
		return d.ProcessClassDeclaration(node.(*parser.ClassDeclarationTree))
	case parser.ParseTreeType_SUPER_EXPRESSION:
		return d.ProcessSuper(node.(*parser.SuperExpressionTree))
	case parser.ParseTreeType_NEW_TARGET_EXPRESSION:
		return d.ProcessNewTarget(node.(*parser.NewTargetExpressionTree))
	case parser.ParseTreeType_YIELD_EXPRESSION:
		return d.ProcessYield(node.(*parser.YieldExpressionTree))
	case parser.ParseTreeType_AWAIT_EXPRESSION:
		return d.ProcessAwait(node.(*parser.AwaitExpressionTree))
	case parser.ParseTreeType_FOR_OF_STATEMENT:
		return d.ProcessForOf(node.(*parser.ForOfStatementTree))
	case parser.ParseTreeType_FOR_AWAIT_OF_STATEMENT:
		return d.ProcessForAwaitOf(node.(*parser.ForAwaitOfStatementTree))
	case parser.ParseTreeType_EXPORT_DECLARATION:
		return d.ProcessExportDecl(node.(*parser.ExportDeclarationTree))
	case parser.ParseTreeType_EXPORT_SPECIFIER:
		return d.ProcessExportSpec(node.(*parser.ExportSpecifierTree))
	case parser.ParseTreeType_IMPORT_DECLARATION:
		return d.ProcessImportDecl(node.(*parser.ImportDeclarationTree))
	case parser.ParseTreeType_IMPORT_SPECIFIER:
		return d.ProcessImportSpec(node.(*parser.ImportSpecifierTree))
	case parser.ParseTreeType_DYNAMIC_IMPORT_EXPRESSION:
		return d.ProcessDynamicImport(node.(*parser.DynamicImportTree))
	case parser.ParseTreeType_IMPORT_META_EXPRESSION:
		return d.ProcessImportMeta(node.(*parser.ImportMetaExpressionTree))
	case parser.ParseTreeType_ARRAY_PATTERN:
		return d.ProcessArrayPattern(node.(*parser.ArrayPatternTree))
	case parser.ParseTreeType_OBJECT_PATTERN:
		return d.ProcessObjectPattern(node.(*parser.ObjectPatternTree))
	case parser.ParseTreeType_DEFAULT_PARAMETER:
		return d.ProcessDefaultParameter(node.(*parser.DefaultParameterTree))
	case parser.ParseTreeType_ITER_REST:
		return d.ProcessIterRest(node.(*parser.IterRestTree))
	case parser.ParseTreeType_ITER_SPREAD:
		return d.ProcessIterSpread(node.(*parser.IterSpreadTree))
		// ES2019
	case parser.ParseTreeType_OBJECT_REST:
		return d.ProcessObjectPatternElement(node.(*parser.ObjectRestTree))
	case parser.ParseTreeType_OBJECT_SPREAD:
		return d.ProcessObjectSpread(node.(*parser.ObjectSpreadTree))
		// ES2022
	case parser.ParseTreeType_FIELD_DECLARATION:
		return d.ProcessField(node.(*parser.FieldDeclarationTree))
	case parser.ParseTreeType_COMPUTED_PROPERTY_FIELD:
		return d.ProcessComputedPropertyField(node.(*parser.ComputedPropertyFieldTree))
	case parser.ParseTreeType_ARGUMENT_LIST:
		// TODO(johnlenz): handle these or remove parser support
	default:
	}
	return d.ProcessIllegalToken(node)
}
