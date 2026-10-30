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

// IRFactory_ScopeTracker tracks whether a certain language feature is
// currently in scope. Use in conjunction with try-with-resources to
// auto-handle decrementing.
//
// Java: private static class IRFactory.TransformDispatcher.ScopeTracker
type IRFactory_ScopeTracker struct {
	usageCount int
}

// Java: IRFactory.TransformDispatcher.ScopeTracker.inScope()
func (s *IRFactory_ScopeTracker) InScope() bool {
	return s.usageCount > 0
}

// Java: IRFactory.TransformDispatcher.ScopeTracker.increment()
func (s *IRFactory_ScopeTracker) Increment() irFactory_ScopeTracker_autoDecrement {
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
func (s *IRFactory_ScopeTracker) MaybeIncrement(condition bool) irFactory_ScopeTracker_autoDecrement {
	if condition {
		jutil.CheckState(s.usageCount >= 0)
		s.usageCount++
	}
	return irFactory_ScopeTracker_autoDecrement{tracker: s, condition: condition, closed: false}
}

// irFactory_ScopeTracker_autoDecrement is the anonymous AutoDecrement class
// of ScopeTracker.maybeIncrement.
type irFactory_ScopeTracker_autoDecrement struct {
	tracker   *IRFactory_ScopeTracker
	condition bool
	closed    bool
}

// Java: IRFactory.TransformDispatcher.ScopeTracker.maybeIncrement(boolean).new AutoDecrement() {...}.close()
func (a *irFactory_ScopeTracker_autoDecrement) Close() {
	jutil.CheckState(!a.closed)
	if a.condition {
		a.tracker.usageCount--
		jutil.CheckState(a.tracker.usageCount >= 0)
	}
	a.closed = true
}

// IRFactory_IdentifierType indicates the valid values usages for an
// identifier.
//
// Java: private enum IRFactory.TransformDispatcher.IdentifierType
type IRFactory_IdentifierType int32

const (
	// The identifier can never be private.
	IRFactory_IdentifierType_STANDARD IRFactory_IdentifierType = iota
	// The identifier can be private if it is in a class scope.
	IRFactory_IdentifierType_CAN_BE_PRIVATE
)

var irFactory_IdentifierType_names = [...]string{"STANDARD", "CAN_BE_PRIVATE"}

// String is Enum.name().
func (t IRFactory_IdentifierType) String() string {
	return irFactory_IdentifierType_names[t]
}

// NewStringNodeFromIdentifier creates a StringNode using the value of an
// IdentifierToken. Reports an error if the identifier value is private and
// the identifier type is STANDARD or the identifier is not within the scope
// of a class.
//

// Java: IRFactory.TransformDispatcher.processBlock(BlockTree)
func (f *IRFactory) ProcessBlock(blockNode *parser.BlockTree) rhino.Node {
	node := f.NewNode(rhino.Token_BLOCK)
	for i, child := range blockNode.Statements {
		childNode := f.Transform(child)
		node.AddChildToBack(childNode)
		f.AttachPossibleTrailingComment(childNode, child.GetEnd())
		// Go-only: performance — nothing reads a statement's parse tree once it
		// is transformed: dropping it lets the GC free it while the rest of the
		// file is transformed, instead of holding a large file's whole parse
		// tree next to its whole AST.
		blockNode.Statements[i] = nil
	}

	if endOff := blockNode.GetEnd(); f.hasAnyPendingCommentBeforeOffset(endOff.Offset) {
		lastComment := f.ParseNonJSDocCommentAt(f.linePos(endOff), false)
		f.AddExtraTrailingComment(node, lastComment)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processBreakStatement(BreakStatementTree)
func (f *IRFactory) ProcessBreakStatement(statementNode *parser.BreakStatementTree) rhino.Node {
	node := f.NewNode(rhino.Token_BREAK)
	if statementNode.HasLabel() {
		labelName := f.TransformLabelName(statementNode.Name)
		node.AddChildToBack(labelName)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.transformLabelName(IdentifierToken)
func (f *IRFactory) TransformLabelName(token parser.Token) rhino.Node {
	checkTokenNotNull(token)
	label := f.MaybeAddNonJsDocComment(
		f.NewStringNodeFromIdentifier(rhino.Token_LABEL_NAME, IRFactory_IdentifierType_STANDARD, token),
		token.GetStart())
	f.SetSourceInfoFromToken(label, token)
	return label
}

// Java: IRFactory.TransformDispatcher.processContinueStatement(ContinueStatementTree)
func (f *IRFactory) ProcessContinueStatement(statementNode *parser.ContinueStatementTree) rhino.Node {
	node := f.NewNode(rhino.Token_CONTINUE)
	if statementNode.HasLabel() {
		labelName := f.TransformLabelName(statementNode.Name)
		node.AddChildToBack(labelName)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processDoLoop(DoWhileStatementTree)
func (f *IRFactory) ProcessDoLoop(loopNode *parser.DoWhileStatementTree) rhino.Node {
	return f.NewNodeWith2Children(rhino.Token_DO, f.TransformBlock(loopNode.Body), f.Transform(loopNode.Condition))
}

// Java: IRFactory.TransformDispatcher.processEmptyStatement(EmptyStatementTree)
func (f *IRFactory) ProcessEmptyStatement(unused *parser.EmptyStatementTree) rhino.Node {
	return f.NewNode(rhino.Token_EMPTY)
}

// Java: IRFactory.TransformDispatcher.processExpressionStatement(ExpressionStatementTree)
func (f *IRFactory) ProcessExpressionStatement(statementNode *parser.ExpressionStatementTree) rhino.Node {
	node := f.NewNode(rhino.Token_EXPR_RESULT)
	node.AddChildToBack(f.Transform(statementNode.Expression))
	return node
}

// Java: IRFactory.TransformDispatcher.processForInLoop(ForInStatementTree)
func (f *IRFactory) ProcessForInLoop(loopNode *parser.ForInStatementTree) rhino.Node {
	// TODO(bradfordcsmith): Rename initializer to something more intuitive
	// like "lhs"
	initializer := f.Transform(loopNode.Initializer)
	return f.NewNodeWith3Children(
		rhino.Token_FOR_IN, initializer, f.Transform(loopNode.Collection), f.TransformBlock(loopNode.Body))
}

// Java: IRFactory.TransformDispatcher.processForOf(ForOfStatementTree)
func (f *IRFactory) ProcessForOf(loopNode *parser.ForOfStatementTree) rhino.Node {
	f.MaybeWarnForFeature(loopNode, parser.FeatureSet_Feature_FOR_OF)
	initializer := f.Transform(loopNode.Initializer)
	return f.NewNodeWith3Children(
		rhino.Token_FOR_OF, initializer, f.Transform(loopNode.Collection), f.TransformBlock(loopNode.Body))
}

// Java: IRFactory.TransformDispatcher.processForAwaitOf(ForAwaitOfStatementTree)
func (f *IRFactory) ProcessForAwaitOf(loopNode *parser.ForAwaitOfStatementTree) rhino.Node {
	f.MaybeWarnForFeature(loopNode, parser.FeatureSet_Feature_FOR_AWAIT_OF)
	initializer := f.Transform(loopNode.Initializer)
	return f.NewNodeWith3Children(
		rhino.Token_FOR_AWAIT_OF,
		initializer,
		f.Transform(loopNode.Collection),
		f.TransformBlock(loopNode.Body))
}

// Java: IRFactory.TransformDispatcher.processForLoop(ForStatementTree)
func (f *IRFactory) ProcessForLoop(loopNode *parser.ForStatementTree) rhino.Node {
	node := f.NewNodeWith3Children(
		rhino.Token_FOR,
		f.TransformOrEmpty(loopNode.Initializer, loopNode),
		f.TransformOrEmpty(loopNode.Condition, loopNode),
		f.TransformOrEmpty(loopNode.Increment, loopNode))
	node.AddChildToBack(f.TransformBlock(loopNode.Body))
	return node
}

// Java: IRFactory.TransformDispatcher.transformOrEmpty(ParseTree, ParseTree)
func (f *IRFactory) TransformOrEmpty(tree parser.ParseTree, parent parser.ParseTree) rhino.Node {
	if tree == nil {
		n := f.NewNode(rhino.Token_EMPTY)
		f.SetSourceInfo(n, parent)
		return n
	}
	return f.Transform(tree)
}

// Java: IRFactory.TransformDispatcher.transformOrEmpty(IdentifierToken, ParseTree)
func (f *IRFactory) TransformOrEmptyIdentifier(token parser.Token, parent parser.ParseTree) rhino.Node {
	if token.IsNull() {
		n := f.NewNode(rhino.Token_EMPTY)
		f.SetSourceInfo(n, parent)
		return n
	}
	return f.ProcessNameOfToken(token, rhino.Token_NAME, IRFactory_IdentifierType_STANDARD)
}

// Java: IRFactory.TransformDispatcher.processFunction(FunctionDeclarationTree)
func (f *IRFactory) ProcessFunction(functionTree *parser.FunctionDeclarationTree) rhino.Node {
	isDeclaration := functionTree.Kind == parser.FunctionDeclarationTree_Kind_DECLARATION
	isMember := functionTree.Kind == parser.FunctionDeclarationTree_Kind_MEMBER
	isArrow := functionTree.Kind == parser.FunctionDeclarationTree_Kind_ARROW
	isAsync := functionTree.IsAsync
	isGenerator := functionTree.IsGenerator
	isSignature := functionTree.FunctionBody.GetParseTreeBase().Type == parser.ParseTreeType_EMPTY_STATEMENT

	identifierType := IRFactory_IdentifierType_STANDARD
	if functionTree.IsClassMember {
		identifierType = IRFactory_IdentifierType_CAN_BE_PRIVATE
	}

	if isGenerator {
		f.MaybeWarnForFeature(functionTree, parser.FeatureSet_Feature_GENERATORS)
	}

	if isMember {
		f.MaybeWarnForFeature(functionTree, parser.FeatureSet_Feature_MEMBER_DECLARATIONS)
	}

	if isArrow {
		f.MaybeWarnForFeature(functionTree, parser.FeatureSet_Feature_ARROW_FUNCTIONS)
	}

	if isAsync {
		f.MaybeWarnForFeature(functionTree, parser.FeatureSet_Feature_ASYNC_FUNCTIONS)
	}

	if isGenerator && isAsync {
		f.MaybeWarnForFeature(functionTree, parser.FeatureSet_Feature_ASYNC_GENERATORS)
	}

	name := functionTree.Name
	var newName rhino.Node
	if name.IsPresent() {
		newName = f.ProcessNameWithInlineCommentsOfToken(name, identifierType)
	} else {
		if isDeclaration || isMember {
			f.errorReporter.Error(
				"unnamed function statement", f.sourceName, f.Lineno(functionTree), f.Charno(functionTree))

			// Return the bare minimum to put the AST in a valid state.
			newName = f.CreateMissingNameNode()
		} else {
			newName = f.NewStringNodeWithToken(rhino.Token_NAME, "")
		}

		// Old Rhino tagged the empty name node with the line number of the
		// declaration.
		f.SetSourceInfo(newName, functionTree)
	}

	node := f.NewNode(rhino.Token_FUNCTION)
	if isMember {
		newName.SetString("")
	}
	node.AddChildToBack(newName)

	node.AddChildToBack(f.Transform(functionTree.FormalParameterList))

	bodyNode := f.Transform(functionTree.FunctionBody)
	if !isArrow && !isSignature && !bodyNode.IsBlock() {
		// When in "keep going" mode the parser tries to parse some constructs
		// the compiler doesn't support, repair it here.
		jutil.CheckState(f.config.RunMode() == Config_RunMode_KEEP_GOING)
		bodyNode = rhino.IR_Block()
	}
	f.ParseDirectives(bodyNode)
	node.AddChildToBack(bodyNode)

	node.SetIsGeneratorFunction(isGenerator)
	node.SetIsArrowFunction(isArrow)
	node.SetIsAsyncFunction(isAsync)
	node.PutBooleanProp(rhino.Node_OPT_ES6_TYPED, functionTree.IsOptional)

	f.AttachPossibleTrailingComment(node, functionTree.GetEnd())

	var result rhino.Node

	if isMember {
		f.SetSourceInfo(node, functionTree)
		member := f.NewStringNodeFromIdentifier(rhino.Token_MEMBER_FUNCTION_DEF, identifierType, name)
		member.AddChildToBack(node)
		member.SetStaticMember(functionTree.IsStatic)
		// The source info should only include the identifier, not the entire
		// function expression
		f.SetSourceInfoFromToken(member, name)
		result = member
	} else {
		result = node
	}

	return result
}

// Java: IRFactory.TransformDispatcher.processIfStatement(IfStatementTree)
func (f *IRFactory) ProcessIfStatement(statementNode *parser.IfStatementTree) rhino.Node {
	node := f.NewNode(rhino.Token_IF)
	node.AddChildToBack(f.Transform(statementNode.Condition))
	node.AddChildToBack(f.TransformBlock(statementNode.IfClause))
	if statementNode.ElseClause != nil {
		node.AddChildToBack(f.TransformBlock(statementNode.ElseClause))
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processDebuggerStatement(DebuggerStatementTree)
func (f *IRFactory) ProcessDebuggerStatement(unused *parser.DebuggerStatementTree) rhino.Node {
	return f.NewNode(rhino.Token_DEBUGGER)
}

// Java: IRFactory.TransformDispatcher.processLabeledStatement(LabelledStatementTree)
func (f *IRFactory) ProcessLabeledStatement(labelTree *parser.LabelledStatementTree) rhino.Node {
	statement := f.Transform(labelTree.Statement)
	if statement.IsFunction() ||
		statement.IsClass() ||
		statement.IsLet() ||
		statement.IsConst() {
		f.errorReporter.Error(
			"Lexical declarations are only allowed at top level or inside a block.",
			f.sourceName,
			f.Lineno(labelTree),
			f.Charno(labelTree))
		return statement // drop the LABEL node so that the resulting AST is valid
	}
	return f.NewNodeWith2Children(rhino.Token_LABEL, f.TransformLabelName(labelTree.Name), statement)
}

// Java: IRFactory.TransformDispatcher.processReturnStatement(ReturnStatementTree)
func (f *IRFactory) ProcessReturnStatement(statementNode *parser.ReturnStatementTree) rhino.Node {
	node := f.NewNode(rhino.Token_RETURN)
	if statementNode.Expression != nil {
		node.AddChildToBack(f.Transform(statementNode.Expression))
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processSwitchCase(CaseClauseTree)
func (f *IRFactory) ProcessSwitchCase(caseNode *parser.CaseClauseTree) rhino.Node {
	expr := caseNode.Expression
	node := f.NewNodeWithChild(rhino.Token_CASE, f.Transform(expr))
	block := f.NewNode(rhino.Token_BLOCK)
	block.SetIsAddedBlock(true)
	f.SetSourceInfo(block, caseNode)
	if caseNode.Statements != nil {
		for _, child := range caseNode.Statements {
			block.AddChildToBack(f.Transform(child))
		}
	}
	node.AddChildToBack(block)
	return node
}

// Java: IRFactory.TransformDispatcher.processSwitchDefault(DefaultClauseTree)
func (f *IRFactory) ProcessSwitchDefault(caseNode *parser.DefaultClauseTree) rhino.Node {
	node := f.NewNode(rhino.Token_DEFAULT_CASE)
	block := f.NewNode(rhino.Token_BLOCK)
	block.SetIsAddedBlock(true)
	f.SetSourceInfo(block, caseNode)
	if caseNode.Statements != nil {
		for _, child := range caseNode.Statements {
			block.AddChildToBack(f.Transform(child))
		}
	}
	node.AddChildToBack(block)
	return node
}

// Java: IRFactory.TransformDispatcher.processSwitchStatement(SwitchStatementTree)
func (f *IRFactory) ProcessSwitchStatement(statementNode *parser.SwitchStatementTree) rhino.Node {
	node := f.NewNodeWithChild(rhino.Token_SWITCH, f.Transform(statementNode.Expression))
	switchBody := f.NewNode(rhino.Token_SWITCH_BODY)
	f.SetSourceInfo(switchBody, statementNode)
	node.AddChildToBack(switchBody)
	for _, child := range statementNode.CaseClauses {
		switchBody.AddChildToBack(f.Transform(child))
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processThrowStatement(ThrowStatementTree)
func (f *IRFactory) ProcessThrowStatement(statementNode *parser.ThrowStatementTree) rhino.Node {
	return f.NewNodeWithChild(rhino.Token_THROW, f.Transform(statementNode.Value))
}

// Java: IRFactory.TransformDispatcher.processTryStatement(TryStatementTree)
func (f *IRFactory) ProcessTryStatement(statementNode *parser.TryStatementTree) rhino.Node {
	node := f.NewNodeWithChild(rhino.Token_TRY, f.TransformBlock(statementNode.Body))
	block := f.NewNode(rhino.Token_BLOCK)
	lineSet := false

	cc := statementNode.CatchBlock
	if cc != nil {
		// Mark the enclosing block at the same line as the first catch
		// clause.
		f.SetSourceInfo(block, cc)
		lineSet = true
		block.AddChildToBack(f.Transform(cc))
	}

	node.AddChildToBack(block)

	finallyBlock := statementNode.FinallyBlock
	if finallyBlock != nil {
		node.AddChildToBack(f.TransformBlock(finallyBlock))
	}

	// If we didn't set the line on the catch clause, then
	// we've got an empty catch clause.  Set its line to be the same
	// as the finally block (to match Old Rhino's behavior.)
	if !lineSet && (finallyBlock != nil) {
		f.SetSourceInfo(block, finallyBlock)
	}

	return node
}

// Java: IRFactory.TransformDispatcher.processCatchClause(CatchTree)
func (f *IRFactory) ProcessCatchClause(clauseNode *parser.CatchTree) rhino.Node {
	if clauseNode.Exception.GetParseTreeBase().Type == parser.ParseTreeType_EMPTY_STATEMENT {
		f.MaybeWarnForFeature(clauseNode, parser.FeatureSet_Feature_OPTIONAL_CATCH_BINDING)
	}

	return f.NewNodeWith2Children(
		rhino.Token_CATCH, f.Transform(clauseNode.Exception), f.TransformBlock(clauseNode.CatchBody))
}

// Java: IRFactory.TransformDispatcher.processFinally(FinallyTree)
func (f *IRFactory) ProcessFinally(finallyNode *parser.FinallyTree) rhino.Node {
	return f.TransformBlock(finallyNode.Block)
}

// Java: IRFactory.TransformDispatcher.processVariableStatement(VariableStatementTree)
func (f *IRFactory) ProcessVariableStatement(stmt *parser.VariableStatementTree) rhino.Node {
	// TODO(moz): Figure out why we still need the special handling
	return f.Process(stmt.Declarations)
}

// Java: IRFactory.TransformDispatcher.processVariableDeclarationList(VariableDeclarationListTree)
func (f *IRFactory) ProcessVariableDeclarationList(decl *parser.VariableDeclarationListTree) rhino.Node {
	var declType rhino.Token
	switch decl.DeclarationType {
	case parser.TokenType_CONST:
		f.MaybeWarnForFeature(decl, parser.FeatureSet_Feature_CONST_DECLARATIONS)
		declType = rhino.Token_CONST
	case parser.TokenType_LET:
		f.MaybeWarnForFeature(decl, parser.FeatureSet_Feature_LET_DECLARATIONS)
		declType = rhino.Token_LET
	case parser.TokenType_VAR:
		declType = rhino.Token_VAR
	default:
		panic(jutil.NewIllegalStateException())
	}

	node := f.NewNode(declType)
	for _, child := range decl.Declarations {
		node.AddChildToBack(f.TransformNodeWithInlineComments(child))
	}

	f.AttachPossibleTrailingComment(node, decl.GetEnd())
	return node
}

// Java: IRFactory.TransformDispatcher.processVariableDeclaration(VariableDeclarationTree)
func (f *IRFactory) ProcessVariableDeclaration(decl *parser.VariableDeclarationTree) rhino.Node {
	node := f.TransformNodeWithInlineComments(decl.Lvalue)
	lhs := node
	if node.IsDestructuringPattern() {
		lhs = f.NewNodeWithChild(rhino.Token_DESTRUCTURING_LHS, node)
	}
	if decl.Initializer != nil {
		initializer := f.Transform(decl.Initializer)
		lhs.AddChildToBack(initializer)
		f.SetLength(lhs, decl.Location.Start, decl.Location.End)
	}
	return lhs
}

// Java: IRFactory.TransformDispatcher.processWhileLoop(WhileStatementTree)
func (f *IRFactory) ProcessWhileLoop(stmt *parser.WhileStatementTree) rhino.Node {
	return f.NewNodeWith2Children(rhino.Token_WHILE, f.Transform(stmt.Condition), f.TransformBlock(stmt.Body))
}

// Java: IRFactory.TransformDispatcher.processWithStatement(WithStatementTree)
func (f *IRFactory) ProcessWithStatement(stmt *parser.WithStatementTree) rhino.Node {
	return f.NewNodeWith2Children(rhino.Token_WITH, f.Transform(stmt.Expression), f.TransformBlock(stmt.Body))
}

// irFactory_javaClassOf is Object.getClass().toString() for a parser token:
// "class com.google.javascript.jscomp.parsing.parser.<simple name>".
func irFactory_javaClassOf(token parser.Token) string {
	return "class com.google.javascript.jscomp.parsing.parser." + token.JavaSimpleClassName()
}

// Java: IRFactory.TransformDispatcher.process(ParseTree)
func (f *IRFactory) Process(node parser.ParseTree) rhino.Node {
	switch n := node.(type) {
	case *parser.BinaryOperatorTree:
		return f.ProcessBinaryExpression(n)
	case *parser.ArrayLiteralExpressionTree:
		return f.ProcessArrayLiteral(n)
	case *parser.TemplateLiteralExpressionTree:
		return f.ProcessTemplateLiteral(n)
	case *parser.TemplateLiteralPortionTree:
		return f.ProcessTemplateLiteralPortion(n)
	case *parser.TemplateSubstitutionTree:
		return f.ProcessTemplateSubstitution(n)
	case *parser.UnaryExpressionTree:
		return f.ProcessUnaryExpression(n)
	case *parser.BlockTree:
		return f.ProcessBlock(n)
	case *parser.BreakStatementTree:
		return f.ProcessBreakStatement(n)
	case *parser.CallExpressionTree:
		return f.ProcessFunctionCall(n)
	case *parser.OptChainCallExpressionTree:
		return f.ProcessOptChainFunctionCall(n)
	case *parser.CaseClauseTree:
		return f.ProcessSwitchCase(n)
	case *parser.DefaultClauseTree:
		return f.ProcessSwitchDefault(n)
	case *parser.CatchTree:
		return f.ProcessCatchClause(n)
	case *parser.ContinueStatementTree:
		return f.ProcessContinueStatement(n)
	case *parser.DoWhileStatementTree:
		return f.ProcessDoLoop(n)
	case *parser.EmptyStatementTree:
		return f.ProcessEmptyStatement(n)
	case *parser.ExpressionStatementTree:
		return f.ProcessExpressionStatement(n)
	case *parser.DebuggerStatementTree:
		return f.ProcessDebuggerStatement(n)
	case *parser.ThisExpressionTree:
		return f.ProcessThisExpression(n)
	case *parser.ForStatementTree:
		return f.ProcessForLoop(n)
	case *parser.ForInStatementTree:
		return f.ProcessForInLoop(n)
	case *parser.FunctionDeclarationTree:
		return f.ProcessFunction(n)
	case *parser.MemberLookupExpressionTree:
		return f.ProcessElementGet(n)
	case *parser.OptionalMemberLookupExpressionTree:
		return f.ProcessOptChainElementGet(n)
	case *parser.MemberExpressionTree:
		return f.ProcessPropertyGet(n)
	case *parser.OptionalMemberExpressionTree:
		return f.ProcessOptChainPropertyGet(n)
	case *parser.ConditionalExpressionTree:
		return f.ProcessConditionalExpression(n)
	case *parser.IfStatementTree:
		return f.ProcessIfStatement(n)
	case *parser.LabelledStatementTree:
		return f.ProcessLabeledStatement(n)
	case *parser.ParenExpressionTree:
		return f.ProcessParenthesizedExpression(n)
	case *parser.IdentifierExpressionTree:
		identifierType := IRFactory_IdentifierType_STANDARD
		if f.privateIdLhsOfInScope.InScope() {
			identifierType = IRFactory_IdentifierType_CAN_BE_PRIVATE
		}
		return f.ProcessName(n, identifierType)
	case *parser.NewExpressionTree:
		return f.ProcessNewExpression(n)
	case *parser.ObjectLiteralExpressionTree:
		return f.ProcessObjectLiteral(n)
	case *parser.ComputedPropertyDefinitionTree:
		return f.ProcessComputedPropertyDefinition(n)
	case *parser.ComputedPropertyGetterTree:
		return f.ProcessComputedPropertyGetter(n)
	case *parser.ComputedPropertyMethodTree:
		return f.ProcessComputedPropertyMethod(n)
	case *parser.ComputedPropertySetterTree:
		return f.ProcessComputedPropertySetter(n)
	case *parser.ReturnStatementTree:
		return f.ProcessReturnStatement(n)
	case *parser.UpdateExpressionTree:
		return f.ProcessUpdateExpression(n)
	case *parser.ProgramTree:
		return f.ProcessAstRoot(n)
	case *parser.LiteralExpressionTree:
		// STRING, NUMBER, TRUE, FALSE, NULL, REGEXP
		return f.ProcessLiteralExpression(n)
	case *parser.SwitchStatementTree:
		return f.ProcessSwitchStatement(n)
	case *parser.ThrowStatementTree:
		return f.ProcessThrowStatement(n)
	case *parser.TryStatementTree:
		return f.ProcessTryStatement(n)
	case *parser.VariableStatementTree:
		// var const let
		return f.ProcessVariableStatement(n)
	case *parser.VariableDeclarationListTree:
		return f.ProcessVariableDeclarationList(n)
	case *parser.VariableDeclarationTree:
		return f.ProcessVariableDeclaration(n)
	case *parser.WhileStatementTree:
		return f.ProcessWhileLoop(n)
	case *parser.WithStatementTree:
		return f.ProcessWithStatement(n)
	case *parser.CommaExpressionTree:
		return f.ProcessCommaExpression(n)
	case *parser.NullTree:
		// this is not the null literal
		return f.ProcessNull(n)
	case *parser.FinallyTree:
		return f.ProcessFinally(n)
	case *parser.MissingPrimaryExpressionTree:
		return f.ProcessMissingExpression(n)
	case *parser.PropertyNameAssignmentTree:
		return f.ProcessPropertyNameAssignment(n)
	case *parser.GetAccessorTree:
		return f.ProcessGetAccessor(n)
	case *parser.SetAccessorTree:
		return f.ProcessSetAccessor(n)
	case *parser.FormalParameterListTree:
		return f.ProcessFormalParameterList(n)
	case *parser.ClassDeclarationTree:
		return f.ProcessClassDeclaration(n)
	case *parser.SuperExpressionTree:
		return f.ProcessSuper(n)
	case *parser.NewTargetExpressionTree:
		return f.ProcessNewTarget(n)
	case *parser.YieldExpressionTree:
		return f.ProcessYield(n)
	case *parser.AwaitExpressionTree:
		return f.ProcessAwait(n)
	case *parser.ForOfStatementTree:
		return f.ProcessForOf(n)
	case *parser.ForAwaitOfStatementTree:
		return f.ProcessForAwaitOf(n)
	case *parser.ExportDeclarationTree:
		return f.ProcessExportDecl(n)
	case *parser.ExportSpecifierTree:
		return f.ProcessExportSpec(n)
	case *parser.ImportDeclarationTree:
		return f.ProcessImportDecl(n)
	case *parser.ImportSpecifierTree:
		return f.ProcessImportSpec(n)
	case *parser.DynamicImportTree:
		return f.ProcessDynamicImport(n)
	case *parser.ImportMetaExpressionTree:
		return f.ProcessImportMeta(n)
	case *parser.ArrayPatternTree:
		return f.ProcessArrayPattern(n)
	case *parser.ObjectPatternTree:
		return f.ProcessObjectPattern(n)
	case *parser.DefaultParameterTree:
		return f.ProcessDefaultParameter(n)
	case *parser.IterRestTree:
		return f.ProcessIterRest(n)
	case *parser.IterSpreadTree:
		return f.ProcessIterSpread(n)
		// ES2019
	case *parser.ObjectRestTree:
		return f.ProcessObjectPatternElement(n)
	case *parser.ObjectSpreadTree:
		return f.ProcessObjectSpread(n)
		// ES2022
	case *parser.FieldDeclarationTree:
		return f.ProcessField(n)
	case *parser.ComputedPropertyFieldTree:
		return f.ProcessComputedPropertyField(n)
	case *parser.ArgumentListTree:
		// TODO(johnlenz): handle these or remove parser support
	default:
	}
	return f.ProcessIllegalToken(node)
}
