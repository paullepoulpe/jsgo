// Port of com/google/javascript/jscomp/Es6RewriteGenerators.java @ closure-compiler v20261028

/*
 * Copyright 2014 The Closure Compiler Authors.
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

package jscomp

import (
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteGenerators converts ES6 generator functions to valid ES3 code.
// This pass runs after all ES6 features except for yield and generators have
// been transpiled.
//
// Generator transpilation pass uses two sets of node properties:
//
//   - generatorMarker property - to indicate that subtree contains YIELD nodes;
//   - generatorSafe property - the node is known to require no further
//     modifications to work in the transpiled form of the generator body.
//
// The conversion is done in the following steps:
//
//   - Find a generator function: `function *() {}`
//   - Replace its original body with a template
//   - Mark all nodes in original body that contain any YIELD nodes
//   - Transpile every statement of the original body into replaced template
//   - unmarked nodes may be copied into the template with a trivial
//     transpilation of "this", "break", "continue", "return" and "arguments"
//     keywords.
//   - marked nodes must be broken up into multiple states to support the
//     yields they contain.
//
// Es6RewriteGenerators depends on InjectTranspilationRuntimeLibraries to
// inject `generator_engine.js` template.
type Es6RewriteGenerators struct {
	jscompAsyncExecuteFunction RuntimeJsLibManager_JsLibField
	jscompAsyncExecuteProgram  RuntimeJsLibManager_JsLibField
	generatorContext           RuntimeJsLibManager_JsLibField
	jscompCreateGenerator      RuntimeJsLibManager_JsLibField

	compiler   AbstractCompiler
	namespace  rhino.StaticScope
	astFactory *AstFactory

	nullableStringType        *colors.Color
	generatorContextTypeValue AstFactory_Type
	propertyIteratorType      AstFactory_Type
	uniqueIdSupplier          *UniqueIdSupplier
}

const (
	// Java: Es6RewriteGenerators.GENERATOR_FUNCTION
	Es6RewriteGenerators_GENERATOR_FUNCTION = "$jscomp$generator$function"
	// Java: Es6RewriteGenerators.GENERATOR_CONTEXT
	Es6RewriteGenerators_GENERATOR_CONTEXT = "$jscomp$generator$context"
	// Java: Es6RewriteGenerators.GENERATOR_ARGUMENTS
	Es6RewriteGenerators_GENERATOR_ARGUMENTS = "$jscomp$generator$arguments"
	// Java: Es6RewriteGenerators.GENERATOR_THIS
	Es6RewriteGenerators_GENERATOR_THIS = "$jscomp$generator$this"
	// Java: Es6RewriteGenerators.GENERATOR_FORIN_PREFIX
	Es6RewriteGenerators_GENERATOR_FORIN_PREFIX = "$jscomp$generator$forin"
)

// Java: Es6RewriteGenerators.transpiledFeatures
var Es6RewriteGenerators_transpiledFeatures = parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_GENERATORS)

// Java: Es6RewriteGenerators(AbstractCompiler)
func NewEs6RewriteGenerators(compiler AbstractCompiler) *Es6RewriteGenerators {
	if compiler == nil {
		panic(jutil.NewNullPointerException())
	}
	e := &Es6RewriteGenerators{
		compiler:             compiler,
		astFactory:           compiler.CreateAstFactory(),
		namespace:            compiler.GetTranspilationNamespace(),
		propertyIteratorType: AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_FUNCTION_TYPE, colors.StandardColors_TOP_OBJECT),
		uniqueIdSupplier:     compiler.GetUniqueIdSupplier(),
	}

	if compiler.HasOptimizationColors() {
		// typechecking has run, so we must preserve and propagate type information
		e.nullableStringType =
			colors.Color_CreateUnion2(colors.StandardColors_NULL_OR_VOID, colors.StandardColors_STRING)
	}

	runtimeJsLibManager := compiler.GetRuntimeJsLibManager()
	e.generatorContext = runtimeJsLibManager.GetJsLibField("$jscomp.generator.Context")
	e.jscompAsyncExecuteFunction =
		runtimeJsLibManager.GetJsLibField("$jscomp.asyncExecutePromiseGeneratorFunction")
	e.jscompAsyncExecuteProgram =
		runtimeJsLibManager.GetJsLibField("$jscomp.asyncExecutePromiseGeneratorProgram")
	e.jscompCreateGenerator =
		runtimeJsLibManager.GetJsLibField("$jscomp.generator.createGenerator")
	return e
}

func (e *Es6RewriteGenerators) generatorContextType() AstFactory_Type {
	if e.generatorContextTypeValue.IsZero() {
		e.generatorContextTypeValue =
			AstFactory_Type_(
				e.astFactory.CreateNewNode(
					e.astFactory.CreateQNameOfField(e.namespace, e.generatorContext)))
	}
	return e.generatorContextTypeValue
}

// Java: Es6RewriteGenerators.process(Node, Node)
func (e *Es6RewriteGenerators) Process(externs rhino.Node, root rhino.Node) {
	if !e.compiler.GetLifeCycleStage().IsNormalized() {
		panic(jutil.NewIllegalStateException())
	}
	TranspilationPasses_ProcessTranspile(
		e.compiler, root, Es6RewriteGenerators_transpiledFeatures, &generatorFunctionsTranspiler{pass: e})
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, Es6RewriteGenerators_transpiledFeatures)
}

// generatorYieldExposer exposes an expression containing `yield` into an
// equivalent expression where `yield` is in statement/initializer position,
// using ExpressionDecomposer.
type generatorYieldExposer struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	pass       *Es6RewriteGenerators
	decomposer *ExpressionDecomposer
}

func newGeneratorYieldExposer(pass *Es6RewriteGenerators) *generatorYieldExposer {
	return &generatorYieldExposer{
		pass:       pass,
		decomposer: pass.compiler.CreateDefaultExpressionDecomposer(),
	}
}

func (y *generatorYieldExposer) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	n.SetGeneratorMarker(false)
	if _, ok := ast.AsFunction(n); ok {
		return false
	}
	if _, ok := ast.AsYield(n); ok {
		y.visitYield(n)
		return false
	}
	return true
}

func (y *generatorYieldExposer) visitYield(n rhino.Node) {
	if _, ok := ast.AsExprStmt(n.GetParent()); ok {
		return
	}
	if y.decomposer.CanExposeExpression(n) != ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE {
		y.decomposer.MaybeExposeExpression(n)
	} else {
		link :=
			"https://github.com/google/closure-compiler/wiki/FAQ" +
				"#i-get-an-undecomposable-expression-error-for-my-yield-or-await-expression" +
				"-what-do-i-do"
		suggestion := "Please rewrite the yield or await as a separate statement."
		message := "Undecomposable expression: " + suggestion + "\nSee " + link
		y.pass.compiler.Report(JSError_MakeWithNode(n, TranspilationUtil_CANNOT_CONVERT, message))
	}
}

// generatorFunctionsTranspiler finds generator functions and performs ES6 -> ES3 transpilation.
type generatorFunctionsTranspiler struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	pass *Es6RewriteGenerators
}

var generatorFunctionsTranspilerVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_FUNCTION)
	return s
}()

func (*generatorFunctionsTranspiler) visitTokens() *ast.TokenSet {
	return &generatorFunctionsTranspilerVisitTokens
}

var _ nodeTraversal_visitTokens = (*generatorFunctionsTranspiler)(nil)

func (g *generatorFunctionsTranspiler) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if fn, ok := ast.AsFunction(n); ok && fn.IsGenerator() {
		newGeneratorFunctionTranspiler(g.pass, fn, g.pass.uniqueIdSupplier.GetUniqueId(t.GetInput())).
			transpile()
	}
}

// generatorFunctionTranspiler transpiles a single generator function into a
// state machine program.
type generatorFunctionTranspiler struct {
	pass     *Es6RewriteGenerators
	uniqueId string

	// The transpilation context for the state machine program.
	ctx *generatorTranspilationContext

	// The body of original generator function that should be transpiled.
	originalGeneratorBody rhino.Node

	// Counter to generate unique variable names for for-in loops.
	forInCounter int

	// The body of a replacement function: a block node that hoists local variables of a generator
	// program and ends with `return createGenerator(function ($jscomp$generator$context) { ... });`.
	newGeneratorHoistBlock rhino.Node
}

func newGeneratorFunctionTranspiler(pass *Es6RewriteGenerators, genFunc ast.Function, uniqueId string) *generatorFunctionTranspiler {
	s := &generatorFunctionTranspiler{
		pass:                  pass,
		uniqueId:              uniqueId,
		originalGeneratorBody: genFunc.Body().Node(),
	}
	s.ctx = newGeneratorTranspilationContext(pass, s)
	return s
}

// hoistVarNode hoists a var node inside newGeneratorHoistBlock before the final return statement.
func (s *generatorFunctionTranspiler) hoistVarNode(node rhino.Node) {
	ast.MustNameDeclaration(node)
	lastStmt, _ := ast.MustBlock(s.newGeneratorHoistBlock).LastStatement()
	node.InsertBefore(lastStmt.Node())
}

// hoistFunctionDeclarationNode hoists a function declaration statement to the beginning of
// newGeneratorHoistBlock to maintain normalization invariants.
func (s *generatorFunctionTranspiler) hoistFunctionDeclarationNode(node rhino.Node) {
	ast.MustFunction(node)
	s.newGeneratorHoistBlock.AddChildToFront(node)
}

// isTranspiledAsyncFunction detects whether the generator function was generated by
// RewriteAsyncFunctions (`return $jscomp.asyncExecutePromiseGeneratorFunction(function* () {...});`).
func (s *generatorFunctionTranspiler) isTranspiledAsyncFunction(generatorFunction rhino.Node) bool {
	if call, ok := ast.AsCall(generatorFunction.GetParent()); ok {
		callTarget := call.Callee().Node()
		if generatorFunction.GetPrevious() == callTarget &&
			generatorFunction.GetNext().IsNil() &&
			s.pass.jscompAsyncExecuteFunction.Matches(callTarget) {
			grandparent := generatorFunction.GetGrandparent()
			_, isRet := ast.AsReturn(grandparent)
			return isRet && grandparent.GetNext().IsNil()
		}
	}
	return false
}

func (s *generatorFunctionTranspiler) transpile() {
	astFactory := s.pass.astFactory
	genFn := ast.MustFunction(s.originalGeneratorBody.GetParent())
	generatorFunction := genFn.Node()
	genFn.SetGenerator(false)

	// A "program" function: `function ($jscomp$generator$context) {}`
	var program rhino.Node
	programType := colors.StandardColors_TOP_OBJECT
	generatorBody := astFactory.CreateBlock()

	var changeScopeNode rhino.Node
	if s.isTranspiledAsyncFunction(generatorFunction) {
		callTarget := generatorFunction.GetPrevious()

		// Use original async function as a hoist block for local generator variables:
		// generator function -> call -> return -> async function body
		s.newGeneratorHoistBlock = ast.MustBlock(generatorFunction.GetGrandparent().GetParent()).Node()
		changeScopeNode = ast.MustFunction(NodeUtil_GetEnclosingFunction(s.newGeneratorHoistBlock)).Node()

		// asyncExecutePromiseGeneratorFunction   =>   asyncExecutePromiseGeneratorProgram
		callTarget.ReplaceWith(
			astFactory.CreateQNameOfField(s.pass.namespace, s.pass.jscompAsyncExecuteProgram).SrcrefTree(callTarget))

		program = generatorFunction
		// function *() {...}   =>   function *(context) {}
		genFn.Params().Node().AddChildToBack(s.ctx.getJsContextNameNode(s.originalGeneratorBody))
		s.originalGeneratorBody.ReplaceWith(generatorBody)
	} else {
		changeScopeNode = generatorFunction
		genFuncName := genFn.NameNode()
		// The transpiled function needs to be able to refer to itself, so make sure it has a name.
		if genFuncName.Ident() == "" {
			genFuncName.SetIdent(s.ctx.getScopedName(Es6RewriteGenerators_GENERATOR_FUNCTION))
			if astFactory.IsAddingColors() {
				// The name of the function is a variable with the same type as the function expression itself.
				genFuncName.Node().SetColor(generatorFunction.GetColor())
			}
			// Function expression name nodes are always marked constant.
			genFuncName.SetConstantName(true)
		}

		program =
			astFactory.CreateFunction(
				"",
				rhino.IR_ParamList(s.ctx.getJsContextNameNode(s.originalGeneratorBody)),
				generatorBody,
				AstFactory_Type_OfColor(programType))

		// $jscomp.generator.createGenerator
		createGenerator := astFactory.CreateQNameOfField(s.pass.namespace, s.pass.jscompCreateGenerator)
		// Replace original generator function body with:
		// return $jscomp.generator.createGenerator(<origGenerator>, <program function>);
		s.newGeneratorHoistBlock =
			astFactory.CreateBlock(
				astFactory.
					CreateReturn(
						// TODO(b/142881197): we can't give a more accurate type right now.
						astFactory.CreateCallWithUnknownType(
							createGenerator, genFuncName.Node().CloneNode(), program)).
					SrcrefTree(s.originalGeneratorBody))
		s.originalGeneratorBody.ReplaceWith(s.newGeneratorHoistBlock)
	}

	// New scopes and any changes to scopes should be reported individually.
	s.pass.compiler.ReportChangeToChangeScope(program)

	NodeTraversal_Traverse(s.pass.compiler, s.originalGeneratorBody, &generatorYieldNodeMarker{})

	// Test if end of generator function is reachable
	shouldAddFinalJump := !s.isEndOfBlockUnreachable(s.originalGeneratorBody)

	// Transpile statements from original generator function
	origBody := ast.MustBlock(s.originalGeneratorBody)
	for stmt, ok := origBody.FirstStatement(); ok; stmt, ok = origBody.FirstStatement() {
		s.transpileStatement(stmt.Detach().Node())
	}

	// Ensure that the state machine program ends
	finalBlock := astFactory.CreateBlock()
	if shouldAddFinalJump {
		finalBlock.AddChildToBack(
			s.ctx.callContextMethodResult(
				s.originalGeneratorBody, "jumpToEnd", AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID)))
	}
	s.ctx.currentCase.jumpToCase(s.ctx.programEndCase, finalBlock)
	s.ctx.currentCase.mayFallThrough = true

	s.ctx.finalizeTransformation(generatorBody)
	s.ctx.checkStateIsEmpty()

	s.pass.compiler.ReportChangeToChangeScope(changeScopeNode)
}

func (s *generatorFunctionTranspiler) transpileStatement(statement rhino.Node) {
	s.transpileStatementWithCases(statement, nil, nil)
}

// transpileStatementWithCases transpiles a detached statement node and appends
// the transpiled instructions to ctx.currentCase.
func (s *generatorFunctionTranspiler) transpileStatementWithCases(
	statement rhino.Node,
	breakCase *generatorStateCase,
	continueCase *generatorStateCase,
) {
	if !rhino.IR_MayBeStatement(statement) || !statement.GetParent().IsNil() {
		panic(jutil.NewIllegalStateException())
	}

	if !statement.IsGeneratorMarker() {
		s.transpileUnmarkedNode(statement)
		return
	}
	switch statement.GetToken() {
	case rhino.Token_LABEL:
		s.transpileLabel(statement)
	case rhino.Token_BLOCK:
		s.transpileBlock(statement)
	case rhino.Token_EXPR_RESULT:
		s.transpileExpressionResult(statement)
	case rhino.Token_VAR:
		s.transpileVar(statement)
	case rhino.Token_RETURN:
		s.transpileReturn(statement)
	case rhino.Token_THROW:
		s.transpileThrow(statement)
	case rhino.Token_IF:
		s.transpileIf(statement, breakCase)
	case rhino.Token_FOR:
		s.transpileFor(statement, breakCase, continueCase)
	case rhino.Token_FOR_IN:
		s.transpileForIn(statement, breakCase, continueCase)
	case rhino.Token_DO:
		s.transpileDo(statement, breakCase, continueCase)
	case rhino.Token_TRY:
		s.transpileTry(statement, breakCase)
	case rhino.Token_SWITCH:
		s.transpileSwitch(statement, breakCase)
	default:
		// NOTE: There is no WHILE case because this pass runs after normalization,
		// which converts all while loops to for loops.
		panic(jutil.NewIllegalStateExceptionMsg("Unsupported token: " + statement.GetToken().String()))
	}
}

// transpileUnmarkedNode transpiles code that doesn't contain `yield`s.
func (s *generatorFunctionTranspiler) transpileUnmarkedNode(n rhino.Node) {
	if fn, ok := ast.AsFunction(n); ok {
		// An inner function should be created only once, when the generator function is first
		// called, not every time the inner callback we create is called, so hoist function
		// declaration statements into the outer scope along with variable declarations.
		functionName := fn.NameNode().Ident()
		if functionName == "" || strings.HasPrefix(functionName, Es6RewriteGenerators_GENERATOR_FUNCTION) {
			panic(jutil.NewIllegalStateException())
		}
		s.hoistFunctionDeclarationNode(n)
		return
	}
	if n.IsBlock() || n.IsAddedBlock() {
		s.ctx.transpileUnmarkedBlock(n)
	} else {
		s.ctx.transpileUnmarkedBlock(rhino.IR_BlockOfNode(n))
	}
}

func (s *generatorFunctionTranspiler) transpileLabel(n rhino.Node) {
	// Collect all label names in "a: b: c: {}" statement
	var labelNames []rhino.Node
	for label, ok := ast.AsLabel(n); ok; label, ok = ast.AsLabel(n) {
		name, body := label.Name().Node(), label.Body().Node()
		labelNames = append(labelNames, name.Detach())
		n = body.Detach()
	}

	var continueCase *generatorStateCase
	if NodeUtil_IsLoopStructure(n) {
		continueCase = s.ctx.createCase()
	}
	breakCase := s.ctx.createCase()
	s.ctx.pushLabels(labelNames, breakCase, continueCase)
	s.transpileStatementWithCases(n, breakCase, continueCase)
	s.ctx.popLabels(labelNames)

	if breakCase != s.ctx.currentCase {
		s.ctx.switchCaseTo(breakCase)
	}
}

func (s *generatorFunctionTranspiler) transpileBlock(n rhino.Node) {
	blk := ast.MustBlock(n)
	for stmt, ok := blk.FirstStatement(); ok; stmt, ok = blk.FirstStatement() {
		s.transpileStatement(stmt.Detach().Node())
	}
}

func (s *generatorFunctionTranspiler) transpileExpressionResult(n rhino.Node) {
	exprStmt := ast.MustExprStmt(n)
	exposedExpression := s.exposeYieldAndTranspileRest(exprStmt.Expr().Detach().Node())
	decomposed := s.transpileYields(exposedExpression)

	// Transpile "a = yield;" into "a = $context.yieldResult;"
	// But don't transpile "yield;" into "$context.yieldResult;"
	// as it influences the collapsing of empty case sections.
	if !exposedExpression.IsYield() {
		n.AddChildToFront(s.prepareNodeForWrite(decomposed))
		n.SetGeneratorMarker(false)
		s.ctx.writeGeneratedNode(n)
	}
}

func (s *generatorFunctionTranspiler) transpileVar(n rhino.Node) {
	n.SetGeneratorMarker(false)
	varDecl := ast.MustNameDeclaration(n)
	newVars := n.CloneNode()
	for decl, ok := varDecl.FirstDeclarator(); ok; {
		for ok && !decl.Node().IsGeneratorMarker() {
			next, nextOk := decl.Next()
			newVars.AddChildToBack(decl.Detach().Node())
			decl, ok = next, nextOk
		}
		if newVars.HasChildren() {
			s.transpileUnmarkedNode(newVars)
			newVars = n.CloneNode()
		}

		if ok {
			next, nextOk := decl.Next()
			nameDecl := decl.Detach().MustName()
			if initExpr, hasInit := nameDecl.Init(); hasInit {
				nameDecl.Node().AddChildToFront(s.maybeDecomposeExpression(initExpr.Detach().Node()))
			}
			nameDecl.Node().SetGeneratorMarker(false)
			newVars.AddChildToBack(nameDecl.Node())
			decl, ok = next, nextOk
		}
	}

	if newVars.HasChildren() {
		s.transpileUnmarkedNode(newVars)
	}
}

func (s *generatorFunctionTranspiler) transpileReturn(n rhino.Node) {
	retStmt := ast.MustReturn(n)
	var retVal rhino.Node
	if val, hasVal := retStmt.Value(); hasVal {
		retVal = val.Detach().Node()
	}
	n.AddChildToFront(
		s.ctx.returnExpression(
			n, s.prepareNodeForWrite(s.maybeDecomposeExpression(retVal))))
	s.ctx.writeGeneratedNode(n)
	s.ctx.currentCase.mayFallThrough = false
}

func (s *generatorFunctionTranspiler) transpileThrow(n rhino.Node) {
	throwStmt := ast.MustThrow(n)
	n.AddChildToFront(s.prepareNodeForWrite(s.maybeDecomposeExpression(throwStmt.Value().Detach().Node())))
	s.ctx.writeGeneratedNode(n)
	s.ctx.currentCase.mayFallThrough = false
}

// exposeYieldAndTranspileRest exposes YIELD operator so it's free of side
// effects, transpiling decomposed prefix statements along the way.
func (s *generatorFunctionTranspiler) exposeYieldAndTranspileRest(n rhino.Node) rhino.Node {
	if n.IsYield() {
		return n
	}

	isExpression := rhino.IR_MayBeExpression(n)
	var block rhino.Node
	if isExpression {
		block = s.pass.astFactory.CreateBlock(s.pass.astFactory.CreateReturn(n))
	} else {
		block = s.pass.astFactory.CreateBlock(n)
	}
	NodeTraversal_Traverse(s.pass.compiler, n, newGeneratorYieldExposer(s.pass))
	NodeTraversal_Traverse(s.pass.compiler, block, &generatorYieldNodeMarker{})

	lastStmt, _ := ast.MustBlock(block).LastStatement()
	decomposed := lastStmt.Detach().Node()
	s.transpileStatement(block)
	if isExpression {
		val, hasVal := ast.MustReturn(decomposed).Value()
		if !hasVal {
			return rhino.NilNode
		}
		return val.Detach().Node()
	}
	return decomposed
}

func (s *generatorFunctionTranspiler) maybeDecomposeExpression(n rhino.Node) rhino.Node {
	if n.IsNil() || !n.IsGeneratorMarker() {
		return n
	}
	return s.transpileYields(s.exposeYieldAndTranspileRest(n))
}

func (s *generatorFunctionTranspiler) prepareNodeForWrite(n rhino.Node) rhino.Node {
	if n.IsNil() {
		return rhino.NilNode
	}

	var wrapper rhino.Node
	if rhino.IR_MayBeStatement(n) {
		wrapper = s.pass.astFactory.CreateBlock(n)
	} else {
		wrapper = s.pass.astFactory.ExprResult(n)
	}
	NodeTraversal_Traverse(s.pass.compiler, wrapper, newGeneratorUnmarkedNodeTranspiler(s.pass, s, s.ctx))
	return wrapper.RemoveFirstChild()
}

func (s *generatorFunctionTranspiler) transpileYields(n rhino.Node) rhino.Node {
	if !n.IsGeneratorMarker() {
		return n
	}
	jumpToSection := s.ctx.createCase()
	yieldNode := s.findYield(n)
	yld := ast.MustYield(yieldNode)
	var rawYieldVal rhino.Node
	if val, hasVal := yld.Value(); hasVal {
		rawYieldVal = val.Detach().Node()
	}
	yieldExpression :=
		s.prepareNodeForWrite(s.maybeDecomposeExpression(rawYieldVal))
	if yld.IsYieldAll() {
		s.ctx.yieldAll(yieldExpression, jumpToSection, yieldNode)
	} else {
		s.ctx.yieldValue(yieldExpression, jumpToSection, yieldNode)
	}
	s.ctx.switchCaseTo(jumpToSection)
	yieldResult := s.ctx.yieldResult(yieldNode)
	if yieldNode == n {
		return yieldResult
	}
	yieldNode.ReplaceWith(yieldResult)
	for yieldResult != n {
		yieldResult = yieldResult.GetParent()
		yieldResult.SetGeneratorMarker(false)
	}
	return n
}

func (s *generatorFunctionTranspiler) transpileIf(n rhino.Node, breakCase *generatorStateCase) {
	astFactory := s.pass.astFactory
	ifStmt := ast.MustIf(n)
	ifBlock := ifStmt.Then().Node()
	elsePart, hasElse := ifStmt.Else()
	elseBlock := elsePart.Node()
	condition := s.maybeDecomposeExpression(ifStmt.Cond().Node().Detach())

	if !ifBlock.IsGeneratorMarker() && (!hasElse || !elseBlock.IsGeneratorMarker()) {
		n.AddChildToFront(condition)
		n.SetGeneratorMarker(false)
		s.transpileUnmarkedNode(n)
		return
	}

	ifBlock.Detach()
	if !hasElse {
		elseBlock = astFactory.CreateBlock().Srcref(n)
	} else {
		elseBlock.Detach()
	}

	// Only "else" block is unmarked: swap "if" and "else" blocks and negate the condition.
	if ifBlock.IsGeneratorMarker() && !elseBlock.IsGeneratorMarker() {
		condition = astFactory.CreateNot(condition).Srcref(condition)
		ifBlock, elseBlock = elseBlock, ifBlock
	}

	// Unmarked "if" block (marked "else")
	if !ifBlock.IsGeneratorMarker() {
		endCase := s.ctx.maybeCreateCase(breakCase)
		jumpToBlock := ast.MustBlock(s.ctx.createJumpToBlock(endCase /* allowEmbedding= */, false, ifBlock))
		for stmt, ok := jumpToBlock.FirstStatement(); ok; stmt, ok = jumpToBlock.FirstStatement() {
			jumpToNode := stmt.Detach().Node()
			jumpToNode.SetGeneratorSafe(true)
			ifBlock.AddChildToBack(jumpToNode)
		}
		s.transpileUnmarkedNode(astFactory.CreateIf(condition, ifBlock).Srcref(n))
		s.transpileStatement(elseBlock)
		s.ctx.switchCaseTo(endCase)
		return
	}

	ifCase := s.ctx.createCase()
	endCase := s.ctx.maybeCreateCase(breakCase)

	// Both "if" and "else" blocks marked
	condition = s.prepareNodeForWrite(condition)
	newIfBlock := s.ctx.createJumpToBlock(ifCase /* allowEmbedding= */, true, n)
	s.ctx.writeGeneratedNode(
		astFactory.CreateIf(condition, newIfBlock).Srcref(n))
	s.transpileStatement(elseBlock)
	s.ctx.writeJumpTo(endCase, elseBlock)
	s.ctx.switchCaseTo(ifCase)
	s.transpileStatement(ifBlock)
	s.ctx.switchCaseTo(endCase)
}

func (s *generatorFunctionTranspiler) transpileFor(
	n rhino.Node,
	breakCase *generatorStateCase,
	continueCase *generatorStateCase,
) {
	astFactory := s.pass.astFactory
	forStmt := ast.MustFor(n)
	initPart, hasInit := forStmt.Init()
	condPart, hasCond := forStmt.Cond()
	condition := condPart.Node()
	updatePart, hasIncrement := forStmt.Update()
	increment := updatePart.Node()
	body := forStmt.Body().Node()
	var init rhino.Node
	if hasInit {
		init = s.maybeDecomposeExpression(initPart.Detach().Node())
	}

	if (!hasCond || !condition.IsGeneratorMarker()) &&
		(!hasIncrement || !increment.IsGeneratorMarker()) &&
		!body.IsGeneratorMarker() {
		if hasInit {
			n.AddChildToFront(init)
		}
		n.SetGeneratorMarker(false)
		s.transpileUnmarkedNode(n)
		return
	}

	if hasInit {
		if rhino.IR_MayBeExpression(init) {
			init = astFactory.ExprResult(init).Srcref(init)
		}
		s.transpileUnmarkedNode(init)
	}

	startCase := s.ctx.createCase()
	incrementCase := s.ctx.maybeCreateCase(continueCase)
	endCase := s.ctx.maybeCreateCase(breakCase)

	s.ctx.switchCaseTo(startCase)

	if hasCond {
		condition = s.prepareNodeForWrite(s.maybeDecomposeExpression(condition.Detach()))
		s.ctx.writeGeneratedNode(
			astFactory.
				CreateIf(
					astFactory.CreateNot(condition).Srcref(condition),
					s.ctx.createJumpToBlock(endCase /* allowEmbedding= */, true, n)).
				Srcref(n))
	}

	s.ctx.pushBreakContinueContext(endCase, incrementCase)
	s.transpileStatement(body.Detach())
	s.ctx.popBreakContinueContext()

	s.ctx.switchCaseTo(incrementCase)
	if hasIncrement {
		increment = s.maybeDecomposeExpression(increment.Detach())
		s.transpileUnmarkedNode(astFactory.ExprResult(increment).Srcref(increment))
	}
	s.ctx.writeJumpTo(startCase, n)

	s.ctx.switchCaseTo(endCase)
}

func (s *generatorFunctionTranspiler) transpileForIn(
	n rhino.Node,
	breakCase *generatorStateCase,
	continueCase *generatorStateCase,
) {
	astFactory := s.pass.astFactory
	forInStmt := ast.MustForIn(n)
	target := forInStmt.Target().Node()
	expr := forInStmt.Object().Node()
	body := forInStmt.Body().Node()
	detachedExpr := s.maybeDecomposeExpression(expr.Detach())

	if !target.IsGeneratorMarker() && !body.IsGeneratorMarker() {
		detachedExpr.InsertAfter(target)
		n.SetGeneratorMarker(false)
		s.transpileUnmarkedNode(n)
		return
	}

	var init rhino.Node
	if varDecl, isVar := ast.AsNameDeclarationOf(target.Detach(), rhino.Token_VAR); isVar {
		if target.IsGeneratorMarker() {
			panic(jutil.NewIllegalStateException())
		}
		init = target
		nameDecl := varDecl.MustOnlyDeclarator().MustName()
		if _, hasInit := nameDecl.Init(); hasInit {
			panic(jutil.NewIllegalStateException())
		}
		target = nameDecl.Node().CloneNode()
	} else {
		init = rhino.NewNode(rhino.Token_VAR).Srcref(target)
	}

	child := s.ctx.callContextMethod(target, "forIn", s.pass.propertyIteratorType, detachedExpr)
	forInCounter := s.forInCounter
	s.forInCounter++
	forIn :=
		astFactory.
			CreateName(
				s.ctx.getScopedName(Es6RewriteGenerators_GENERATOR_FORIN_PREFIX)+"$"+strconv.Itoa(forInCounter), AstFactory_Type_(child)).
			Srcref(target)
	forIn.AddChildToFront(child)
	init.AddChildToBack(forIn)

	forInGetNext :=
		astFactory.
			CreateGetProp(forIn.CloneNode(), "getNext", AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT)).
			Srcref(detachedExpr)

	forCond :=
		astFactory.
			CreateNe(
				astFactory.
					CreateAssign(
						target.SetColor(s.pass.nullableStringType),
						astFactory.
							CreateCall(forInGetNext, AstFactory_Type_OfColor(s.pass.nullableStringType)).
							Srcref(detachedExpr)).
					Srcref(detachedExpr),
				astFactory.CreateNull().Srcref(forIn)).
			Srcref(detachedExpr)
	forCond.SetGeneratorMarker(target.IsGeneratorMarker())

	forNode := rhino.IR_ForNode(init, forCond, rhino.IR_Empty().Srcref(n), body.Detach()).Srcref(n)
	s.transpileFor(forNode, breakCase, continueCase)
}

func (s *generatorFunctionTranspiler) transpileDo(
	n rhino.Node,
	breakCase *generatorStateCase,
	continueCase *generatorStateCase,
) {
	doStmt := ast.MustDo(n)
	body := doStmt.Body().Node()
	condNode := doStmt.Cond().Node()
	body.Detach()
	condNode.Detach()

	startCase := s.ctx.createCase()
	breakCase = s.ctx.maybeCreateCase(breakCase)
	continueCase = s.ctx.maybeCreateCase(continueCase)

	s.ctx.switchCaseTo(startCase)

	s.ctx.pushBreakContinueContext(breakCase, continueCase)
	s.transpileStatement(body)
	s.ctx.popBreakContinueContext()

	s.ctx.switchCaseTo(continueCase)
	condition := s.prepareNodeForWrite(s.maybeDecomposeExpression(condNode))
	s.ctx.writeGeneratedNode(
		s.pass.astFactory.
			CreateIf(
				condition, s.ctx.createJumpToBlock(startCase /* allowEmbedding= */, false, n)).
			Srcref(n))
	s.ctx.switchCaseTo(breakCase)
}

func (s *generatorFunctionTranspiler) transpileTry(n rhino.Node, breakCase *generatorStateCase) {
	tryStmt := ast.MustTry(n)
	tryBlock := tryStmt.Body().Node()
	catchNode, hasCatch := tryStmt.Catch()
	finallyPart, hasFinally := tryStmt.Finally()
	finallyBlock := finallyPart.Node()
	tryBlock.Detach()
	if hasCatch {
		catchNode.Detach()
	}
	if hasFinally {
		finallyBlock.Detach()
	}

	var catchCase *generatorStateCase
	if hasCatch {
		catchCase = s.ctx.createCase()
	}
	var finallyCase *generatorStateCase
	if hasFinally {
		finallyCase = s.ctx.createCase()
	}
	endCase := s.ctx.maybeCreateCase(breakCase)

	s.ctx.enterTryBlock(catchCase, finallyCase, tryBlock)
	s.transpileStatement(tryBlock)

	if !hasFinally {
		s.ctx.leaveTryBlock(catchCase, endCase, tryBlock)
	} else {
		s.ctx.switchCaseTo(finallyCase)
		s.ctx.enterFinallyBlock(catchCase, finallyCase, finallyBlock)
		s.transpileStatement(finallyBlock)
		s.ctx.leaveFinallyBlock(endCase, finallyBlock)
	}

	if hasCatch {
		s.ctx.switchCaseTo(catchCase)
		catchBody := catchNode.Body().Node().Detach()
		param, hasParam := catchNode.Param()
		if !hasParam {
			panic(jutil.NewIllegalStateExceptionMsg("Expected catch parameter to be normalized to a NAME"))
		}
		nameParam, ok := param.AsName()
		if !ok {
			panic(jutil.NewIllegalStateExceptionMsg("Expected catch parameter to be normalized to a NAME, found: " + param.Token().String()))
		}
		exceptionName := nameParam.Detach().Node()
		s.ctx.enterCatchBlock(finallyCase, exceptionName)

		s.transpileStatement(catchBody)
		s.ctx.leaveCatchBlock(finallyCase, catchBody)
	}

	s.ctx.switchCaseTo(endCase)
}

func (s *generatorFunctionTranspiler) transpileSwitch(n rhino.Node, breakCase *generatorStateCase) {
	sw := ast.MustSwitch(n)
	switchBody := sw.Body()
	disc := sw.Discriminant().Node()
	n.AddChildToFront(s.maybeDecomposeExpression(disc.Detach()))

	hasGeneratorMarker := false
	for caseSection, ok := switchBody.FirstCase(); ok; caseSection, ok = caseSection.Next() {
		if caseSection.Node().IsGeneratorMarker() {
			hasGeneratorMarker = true
			break
		}
	}
	if !hasGeneratorMarker {
		n.SetGeneratorMarker(false)
		s.transpileUnmarkedNode(n)
		return
	}

	type detachedSwitchCase struct {
		generatedCase *generatorStateCase
		body          rhino.Node
	}
	var detachedCases []detachedSwitchCase

	// We don't have to transpile unmarked cases at the beginning of "switch".
	canSkipUnmarkedCases := true
	for caseSection, ok := switchBody.FirstCase(); ok; caseSection, ok = caseSection.Next() {
		if c, ok := caseSection.AsCase(); ok && c.Value().Node().IsGeneratorMarker() {
			s.pass.compiler.Report(
				JSError_MakeWithNode(
					n, TranspilationUtil_CANNOT_CONVERT_YET, "Case statements that contain yields"))
			return
		}
		body := caseSection.Body().Node()

		if !body.HasChildren() || (canSkipUnmarkedCases && !body.IsGeneratorMarker()) {
			continue
		}

		canSkipUnmarkedCases = s.isEndOfBlockUnreachable(body)

		generatedCase := s.ctx.createCase()
		generatedCase.caseBlock.Srcref(body)

		newBody := rhino.IR_BlockOfNode(s.ctx.createJumpToNode(generatedCase, body))
		newBody.SetIsAddedBlock(true)
		newBody.SetGeneratorSafe(true)
		body.ReplaceWith(newBody)

		detachedCases = append(detachedCases, detachedSwitchCase{generatedCase: generatedCase, body: body})
		caseSection.Node().SetGeneratorMarker(false)
	}

	endCase := s.ctx.maybeCreateCase(breakCase)

	n.SetGeneratorMarker(false)
	switchBody.Node().SetGeneratorMarker(false)
	s.transpileUnmarkedNode(n)
	s.ctx.writeJumpTo(endCase, n)

	s.ctx.pushBreakContext(endCase)
	for _, detachedCase := range detachedCases {
		s.ctx.switchCaseTo(detachedCase.generatedCase)
		s.transpileStatement(detachedCase.body)
	}
	s.ctx.popBreakContext()

	s.ctx.switchCaseTo(endCase)
}

func (s *generatorFunctionTranspiler) findYield(n rhino.Node) rhino.Node {
	yieldFinder := &generatorYieldFinder{}
	NodeTraversal_Traverse(s.pass.compiler, n, yieldFinder)
	rhino.CheckNotNull(yieldFinder.yieldNode)
	return yieldFinder.yieldNode
}

type generatorYieldFinder struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	yieldNode rhino.Node
}

func (y *generatorYieldFinder) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if _, ok := ast.AsFunction(n); ok {
		return false
	}
	if _, ok := ast.AsYield(n); ok {
		if !y.yieldNode.IsNil() {
			panic(jutil.NewIllegalStateException())
		}
		y.yieldNode = n
		return false
	}
	return true
}

func (s *generatorFunctionTranspiler) isEndOfBlockUnreachable(block rhino.Node) bool {
	lastStmt, ok := ast.MustBlock(block).LastStatement()
	if !ok {
		return false
	}
	if innerBlock, isBlk := lastStmt.AsBlock(); isBlk {
		return s.isEndOfBlockUnreachable(innerBlock.Node())
	}
	switch lastStmt.Token() {
	case rhino.Token_RETURN, rhino.Token_THROW, rhino.Token_CONTINUE, rhino.Token_BREAK:
		return true
	default:
		return false
	}
}

// generatorYieldNodeMarker marks "yield" nodes and propagates this information
// up through the detached subtree.
type generatorYieldNodeMarker struct{}

func (y *generatorYieldNodeMarker) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	_, isFn := ast.AsFunction(n)
	return !isFn
}

func (y *generatorYieldNodeMarker) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if _, ok := ast.AsYield(n); ok {
		n.SetGeneratorMarker(true)
	}
	if !parent.IsNil() && n.IsGeneratorMarker() {
		parent.SetGeneratorMarker(true)
	}
}
