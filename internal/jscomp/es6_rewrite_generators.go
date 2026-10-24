// Port of com/google/javascript/jscomp/Es6RewriteGenerators.java @ closure-compiler v20261015

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

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteGenerators converts ES6 generator functions to valid ES3 code.
// This pass runs after all ES6 features except for yield and generators have
// been transpiled.
//
// Genertor transpilation pass uses two sets of node properties:
//
//   - generatorMarker property - to indicate that subtee contains YIELD nodes;
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

	nullableStringType *colors.Color
	// Go-only: language-adapter — Suppliers.memoize(...) is the lazily set
	// generatorContextTypeValue, read through generatorContextType().
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
	e := &Es6RewriteGenerators{}
	jutil.CheckNotNull(compiler)
	e.compiler = compiler
	e.astFactory = compiler.CreateAstFactory()
	e.namespace = compiler.GetTranspilationNamespace()

	if compiler.HasOptimizationColors() {
		// typechecking has run, so we must preserve and propagate type information
		e.nullableStringType =
			colors.Color_CreateUnion(jutil.NewObjectLinkedHashSet(colors.StandardColors_NULL_OR_VOID, colors.StandardColors_STRING))
	} else {
		e.nullableStringType = nil
	}

	runtimeJsLibManager := compiler.GetRuntimeJsLibManager()
	e.generatorContext = runtimeJsLibManager.GetJsLibField("$jscomp.generator.Context")
	e.jscompAsyncExecuteFunction =
		runtimeJsLibManager.GetJsLibField("$jscomp.asyncExecutePromiseGeneratorFunction")
	e.jscompAsyncExecuteProgram =
		runtimeJsLibManager.GetJsLibField("$jscomp.asyncExecutePromiseGeneratorProgram")
	e.jscompCreateGenerator =
		runtimeJsLibManager.GetJsLibField("$jscomp.generator.createGenerator")

	e.generatorContextTypeValue = nil
	e.propertyIteratorType = AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_FUNCTION_TYPE, colors.StandardColors_TOP_OBJECT)
	e.uniqueIdSupplier = compiler.GetUniqueIdSupplier()
	return e
}

// generatorContextType is generatorContextType.get() of the memoizing
// supplier.
//
// Go-only: language-adapter — Suppliers.memoize.
func (e *Es6RewriteGenerators) generatorContextType() AstFactory_Type {
	if e.generatorContextTypeValue == nil {
		e.generatorContextTypeValue =
			AstFactory_Type_(
				e.astFactory.CreateNewNode(
					e.astFactory.CreateQNameOfField(e.namespace, e.generatorContext)))
	}
	return e.generatorContextTypeValue
}

// Java: Es6RewriteGenerators.process(Node, Node)
func (e *Es6RewriteGenerators) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckState(e.compiler.GetLifeCycleStage().IsNormalized())
	TranspilationPasses_ProcessTranspile(
		e.compiler, root, Es6RewriteGenerators_transpiledFeatures, &Es6RewriteGenerators_GeneratorFunctionsTranspiler{outer: e})
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, Es6RewriteGenerators_transpiledFeatures)
}

// Es6RewriteGenerators_YieldExposer exposes expression with yield inside to
// an equivalent expression in which yield is of the form:
//
//	var name = yield expr;
//
// For example, changes the following code:
//
//	{ return x || yield y; }
//
// into:
//
//	{
//	  var temp$$0;
//	  if (temp$$0 = x); else temp$$0 = yield y;
//	  return temp$$0;
//	}
//
// Expression should always be inside a block, so that other statements could
// be added at need.
//
// Uses the ExpressionDecomposer class.
type Es6RewriteGenerators_YieldExposer struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer *Es6RewriteGenerators

	decomposer *ExpressionDecomposer
}

// Java: YieldExposer()
func NewEs6RewriteGenerators_YieldExposer(outer *Es6RewriteGenerators) *Es6RewriteGenerators_YieldExposer {
	y := &Es6RewriteGenerators_YieldExposer{outer: outer}
	y.decomposer = outer.compiler.CreateDefaultExpressionDecomposer()
	return y
}

// Java: YieldExposer.shouldTraverse(NodeTraversal, Node, Node)
func (y *Es6RewriteGenerators_YieldExposer) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	n.SetGeneratorMarker(false)
	if n.IsFunction() {
		return false
	}
	if n.IsYield() {
		y.VisitYield(n)
		return false
	}
	return true
}

// Java: YieldExposer.visitYield(Node)
func (y *Es6RewriteGenerators_YieldExposer) VisitYield(n rhino.Node) {
	if n.GetParent().IsExprResult() {
		return
	}
	if y.decomposer.CanExposeExpression(n) !=
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE {
		y.decomposer.MaybeExposeExpression(n)
	} else {
		link :=
			"https://github.com/google/closure-compiler/wiki/FAQ" +
				"#i-get-an-undecomposable-expression-error-for-my-yield-or-await-expression" +
				"-what-do-i-do"
		suggestion := "Please rewrite the yield or await as a separate statement."
		message := "Undecomposable expression: " + suggestion + "\nSee " + link
		y.outer.compiler.Report(JSError_MakeWithNode(n, TranspilationUtil_CANNOT_CONVERT, message))
	}
}

// Es6RewriteGenerators_GeneratorFunctionsTranspiler finds generator functions
// and performs ES6 -> ES3 trnspilation
type Es6RewriteGenerators_GeneratorFunctionsTranspiler struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *Es6RewriteGenerators
}

// Java: GeneratorFunctionsTranspiler.visit(NodeTraversal, Node, Node)
func (g *Es6RewriteGenerators_GeneratorFunctionsTranspiler) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsGeneratorFunction() {
		NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler(g.outer, n, g.outer.uniqueIdSupplier.GetUniqueId(t.GetInput())).
			Transpile()
	}
}

// Es6RewriteGenerators_SingleGeneratorFunctionTranspiler transpiles a single
// generator function into a state machine program.
type Es6RewriteGenerators_SingleGeneratorFunctionTranspiler struct {
	outer *Es6RewriteGenerators

	uniqueId string

	// The transpilation context for the state machine program.
	context *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext

	// The body of original generator function that should be transpiled
	originalGeneratorBody rhino.Node

	// Counter to generate unique variable names for for-in loops.
	//
	// If there are multiple for-in loops, we'll need a separate variable for each of them.
	forInCounter int

	// The body of a replacement function.
	//
	// It's a block node that hoists local variables of a generator program and returns an actual
	// generator object created from that program:
	//
	//	{
	//	  var a;
	//	  var b;
	//	  ...
	//	  return createGenerator(function ($jscomp$generator$context) { ... });
	//	}
	//
	// The assumption is that the hoist block always ends with a return statement, and all local
	// variables are added before this "return" statement.
	newGeneratorHoistBlock rhino.Node
}

// Java: SingleGeneratorFunctionTranspiler(Node, String)
func NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler(outer *Es6RewriteGenerators, genFunc rhino.Node, uniqueId string) *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler {
	s := &Es6RewriteGenerators_SingleGeneratorFunctionTranspiler{outer: outer}
	s.forInCounter = 0
	s.originalGeneratorBody = genFunc.GetLastChild()
	s.uniqueId = uniqueId
	s.context = NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext(s)
	return s
}

// HoistVarNode hoists a var node inside newGeneratorHoistBlock.
//
// The last statement in the block is expected to be the call to the
// generator execution function.
//
// Java: SingleGeneratorFunctionTranspiler.hoistVarNode(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) HoistVarNode(node rhino.Node) {
	jutil.CheckStateMsg1(node.IsVar(), "%s", node)
	node.InsertBefore(s.newGeneratorHoistBlock.GetLastChild())
}

// HoistFunctionDeclarationNode hoists a function declaration statement inside
// newGeneratorHoistBlock.
//
// In order to maintain the requirements of normalization, function
// declaration statements must always be at the beginning of the function
// body.
//
// Java: SingleGeneratorFunctionTranspiler.hoistFunctionDeclarationNode(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) HoistFunctionDeclarationNode(node rhino.Node) {
	jutil.CheckStateMsg1(node.IsFunction(), "%s", node)
	s.newGeneratorHoistBlock.AddChildToFront(node)
}

// IsTranspiledAsyncFunction detects whether the generator function was
// generated by async function transpilation:
//
//	function() {
//	  ...
//	  return $jscomp.asyncExecutePromiseGeneratorFunction(function* genFunc() {...});
//	}
//
// Java: SingleGeneratorFunctionTranspiler.isTranspiledAsyncFunction(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) IsTranspiledAsyncFunction(generatorFunction rhino.Node) bool {
	if generatorFunction.GetParent().IsCall() && !generatorFunction.GetPrevious().IsNil() {
		callTarget := generatorFunction.GetParent().GetFirstChild()
		if generatorFunction.GetPrevious() == callTarget &&
			generatorFunction.GetNext().IsNil() &&
			s.outer.jscompAsyncExecuteFunction.Matches(callTarget) {
			jutil.CheckState(generatorFunction.GetGrandparent().IsReturn())
			jutil.CheckState(generatorFunction.GetGrandparent().GetNext().IsNil())
			return true
		}
	}
	return false
}

// Java: SingleGeneratorFunctionTranspiler.transpile()
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) Transpile() {
	astFactory := s.outer.astFactory
	generatorFunction := s.originalGeneratorBody.GetParent()
	jutil.CheckState(generatorFunction.IsGeneratorFunction())
	generatorFunction.PutBooleanProp(rhino.Node_GENERATOR_FN, false)

	// A "program" function:
	//   function ($jscomp$generator$context) {
	//   }
	var program rhino.Node
	// function(!Context<YIELD_TYPE>): (void|{value: YIELD_TYPE})
	programType := colors.StandardColors_TOP_OBJECT
	generatorBody := astFactory.CreateBlock()

	var changeScopeNode rhino.Node
	if s.IsTranspiledAsyncFunction(generatorFunction) {
		// Our generatorFunction is a transpiled async function

		// $jscomp.asyncExecutePromiseGeneratorFunction
		callTarget := generatorFunction.GetPrevious()
		jutil.CheckStateMsg1(callTarget.IsGetProp() || callTarget.IsName(), "%s", callTarget)

		// Use original async function as a hoist block for local generator variables:
		// generator function -> call -> return -> async function body
		s.newGeneratorHoistBlock = generatorFunction.GetGrandparent().GetParent()
		jutil.CheckStateMsg1(s.newGeneratorHoistBlock.IsBlock(), "%s", s.newGeneratorHoistBlock)
		changeScopeNode = NodeUtil_GetEnclosingFunction(s.newGeneratorHoistBlock)
		jutil.CheckStateMsg1(changeScopeNode.IsFunction(), "%s", changeScopeNode)

		// asyncExecutePromiseGeneratorFunction   =>   asyncExecutePromiseGeneratorProgram
		callTarget.ReplaceWith(
			astFactory.CreateQNameOfField(s.outer.namespace, s.outer.jscompAsyncExecuteProgram).SrcrefTree(callTarget))

		program = s.originalGeneratorBody.GetParent()
		// function *() {...}   =>   function *(context) {}
		s.originalGeneratorBody.
			GetPrevious().
			AddChildToBack(s.context.GetJsContextNameNode(s.originalGeneratorBody))
		s.originalGeneratorBody.ReplaceWith(generatorBody)
	} else {
		changeScopeNode = generatorFunction
		genFuncName := generatorFunction.GetFirstChild()
		jutil.CheckState(genFuncName.IsName())
		// The transpiled function needs to be able to refer to itself, so make sure it has a name.
		if genFuncName.GetString() == "" {
			genFuncName.SetString(s.context.GetScopedName(Es6RewriteGenerators_GENERATOR_FUNCTION))
			if astFactory.IsAddingColors() {
				// The name of the function is a variable with the same type as the function expression
				// itself.
				genFuncName.SetColor(generatorFunction.GetColor())
			}
			// Function expression name nodes are always marked constant.
			genFuncName.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		}

		// Prepare a "program" function:
		//   function ($jscomp$generator$context) {
		//   }
		program =
			astFactory.CreateFunction(
				"",
				rhino.IR_ParamList(s.context.GetJsContextNameNode(s.originalGeneratorBody)),
				generatorBody,
				AstFactory_Type_OfColor(programType))

		// $jscomp.generator.createGenerator
		createGenerator := astFactory.CreateQNameOfField(s.outer.namespace, s.outer.jscompCreateGenerator)
		// Replace original generator function body with:
		// return $jscomp.generator.createGenerator(<origGenerator>, <program function>);
		s.newGeneratorHoistBlock =
			astFactory.CreateBlock(
				astFactory.
					CreateReturn(
						// TODO(b/142881197): we can't give a more accurate type right now.
						astFactory.CreateCallWithUnknownType(
							createGenerator, genFuncName.CloneNode(), program)).
					SrcrefTree(s.originalGeneratorBody))
		s.originalGeneratorBody.ReplaceWith(s.newGeneratorHoistBlock)
	}

	// New scopes and any changes to scopes should be reported individually.
	s.outer.compiler.ReportChangeToChangeScope(program)

	NodeTraversal_Traverse(s.outer.compiler, s.originalGeneratorBody, &Es6RewriteGenerators_YieldNodeMarker{})

	// Test if end of generator function is reachable
	shouldAddFinalJump := !s.IsEndOfBlockUnreachable(s.originalGeneratorBody)

	// Transpile statements from original generator function
	for s.originalGeneratorBody.HasChildren() {
		s.TranspileStatement(s.originalGeneratorBody.RemoveFirstChild())
	}

	// Ensure that the state machine program ends
	finalBlock := astFactory.CreateBlock()
	if shouldAddFinalJump {
		finalBlock.AddChildToBack(
			s.context.CallContextMethodResult(
				s.originalGeneratorBody, "jumpToEnd", AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID)))
	}
	s.context.currentCase.JumpTo(s.context.programEndCase, finalBlock)
	s.context.currentCase.mayFallThrough = true

	s.context.FinalizeTransformation(generatorBody)
	s.context.CheckStateIsEmpty()

	s.outer.compiler.ReportChangeToChangeScope(changeScopeNode)
}

// TranspileStatement: see TranspileStatementWithCases.
//
// Java: SingleGeneratorFunctionTranspiler.transpileStatement(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileStatement(statement rhino.Node) {
	s.TranspileStatementWithCases(statement, nil, nil)
}

// TranspileStatementWithCases transpiles a detached node and adds transpiled
// version of it to the currentCase of the context.
//
// Java: SingleGeneratorFunctionTranspiler.transpileStatement(Node, TranspilationContext.Case, TranspilationContext.Case)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileStatementWithCases(
	statement rhino.Node,
	breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case,
	continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	jutil.CheckState(rhino.IR_MayBeStatement(statement))
	jutil.CheckState(statement.GetParent().IsNil())

	if !statement.IsGeneratorMarker() {
		s.TranspileUnmarkedNode(statement)
		return
	}
	switch statement.GetToken() {
	case rhino.Token_LABEL:
		s.TranspileLabel(statement)
	case rhino.Token_BLOCK:
		s.TranspileBlock(statement)
	case rhino.Token_EXPR_RESULT:
		s.TranspileExpressionResult(statement)
	case rhino.Token_VAR:
		s.TranspileVar(statement)
	case rhino.Token_RETURN:
		s.TranspileReturn(statement)
	case rhino.Token_THROW:
		s.TranspileThrow(statement)
	case rhino.Token_IF:
		s.TranspileIf(statement, breakCase)
	case rhino.Token_FOR:
		s.TranspileFor(statement, breakCase, continueCase)
	case rhino.Token_FOR_IN:
		s.TranspileForIn(statement, breakCase, continueCase)
	case rhino.Token_DO:
		s.TranspileDo(statement, breakCase, continueCase)
	case rhino.Token_TRY:
		s.TranspileTry(statement, breakCase)
	case rhino.Token_SWITCH:
		s.TranspileSwitch(statement, breakCase)
	default:
		// NOTE: There is no WHILE case, becasue this pass runs after normalization,
		// which converts all while loops to for loops.
		panic(jutil.NewIllegalStateExceptionMsg("Unsupported token: " + statement.GetToken().String()))
	}
}

// TranspileUnmarkedNode transpiles code that doesn't contain `yield`s.
//
// Java: SingleGeneratorFunctionTranspiler.transpileUnmarkedNode(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileUnmarkedNode(n rhino.Node) {
	jutil.CheckState(!n.IsGeneratorMarker())
	if n.IsFunction() {
		// An inner function should be created only once, when the generator function is first
		// called, not every time the inner callback we create is called, so we will hoist
		// function declaration statements into the outer scope along with the other variable
		// declarations.
		// TODO(bradfordcsmith): Ideally we should probably also hoist out function expressions,
		// but there's no strong need to do this now and it would be tricky to get right.
		functionName := n.GetFirstChild().GetString()
		// Make sure there are no "function (...) {...}" statements (note that
		// "function *(...) {...}" becomes "function $jscomp$generator$function(...) {...}" as
		// inner generator functions are transpiled first).
		jutil.CheckState(functionName != "" && !strings.HasPrefix(functionName, Es6RewriteGenerators_GENERATOR_FUNCTION))
		s.HoistFunctionDeclarationNode(n)
		return
	}
	if n.IsBlock() || n.IsAddedBlock() {
		s.context.TranspileUnmarkedBlock(n)
	} else {
		s.context.TranspileUnmarkedBlock(rhino.IR_BlockOfNode(n))
	}
}

// TranspileLabel transpiles a label with marked statement.
//
// Java: SingleGeneratorFunctionTranspiler.transpileLabel(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileLabel(n rhino.Node) {
	// Collect all labels names in "a: b: c: {}" statement
	labelNames := []rhino.Node{}
	for n.IsLabel() {
		labelNames = append(labelNames, n.RemoveFirstChild())
		n = n.RemoveFirstChild()
	}

	// Push label names and continue transpilation
	var continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case
	if NodeUtil_IsLoopStructure(n) {
		continueCase = s.context.CreateCase()
	}
	breakCase := s.context.CreateCase()
	s.context.PushLabels(labelNames, breakCase, continueCase)
	s.TranspileStatementWithCases(n, breakCase, continueCase)
	s.context.PopLabels(labelNames)

	// Switch to endCase if it's not yet active.
	if breakCase != s.context.currentCase {
		s.context.SwitchCaseTo(breakCase)
	}
}

// TranspileBlock transpiles a block.
//
// Java: SingleGeneratorFunctionTranspiler.transpileBlock(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileBlock(n rhino.Node) {
	for n.HasChildren() {
		s.TranspileStatement(n.RemoveFirstChild())
	}
}

// TranspileExpressionResult transpiles marked expression result statement.
//
// Java: SingleGeneratorFunctionTranspiler.transpileExpressionResult(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileExpressionResult(n rhino.Node) {
	exposedExpression := s.ExposeYieldAndTranspileRest(n.RemoveFirstChild())
	decomposed := s.TranspileYields(exposedExpression)

	// Tanspile "a = yield;" into "a = $context.yieldResult;"
	// But don't transpile "yield;" into "$context.yieldResult;"
	// As it influences the collapsing of empty case sections.
	if !exposedExpression.IsYield() {
		n.AddChildToFront(s.PrepareNodeForWrite(decomposed))
		n.SetGeneratorMarker(false)
		s.context.WriteGeneratedNode(n)
	}
}

// TranspileVar transpiles marked "var" statement.
//
// Java: SingleGeneratorFunctionTranspiler.transpileVar(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileVar(n rhino.Node) {
	n.SetGeneratorMarker(false)
	newVars := n.CloneNode()
	for n.HasChildren() {
		var var_ rhino.Node
		// Just collect all unmarked vars and transpile them together.
		for {
			var_ = n.RemoveFirstChild()
			if !(!var_.IsNil() && !var_.IsGeneratorMarker()) {
				break
			}
			newVars.AddChildToBack(var_)
		}
		if newVars.HasChildren() {
			s.TranspileUnmarkedNode(newVars)
			newVars = n.CloneNode()
		}

		// Transpile marked var
		if !var_.IsNil() {
			jutil.CheckState(var_.IsGeneratorMarker())
			var_.AddChildToFront(s.MaybeDecomposeExpression(var_.RemoveFirstChild()))
			var_.SetGeneratorMarker(false)
			newVars.AddChildToBack(var_)
		}
	}

	// Flush the vars if not empty
	if newVars.HasChildren() {
		s.TranspileUnmarkedNode(newVars)
	}
}

// TranspileReturn transpiles marked "return" statement.
//
// Java: SingleGeneratorFunctionTranspiler.transpileReturn(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileReturn(n rhino.Node) {
	n.AddChildToFront(
		s.context.ReturnExpression(
			n, s.PrepareNodeForWrite(s.MaybeDecomposeExpression(n.RemoveFirstChild()))))
	s.context.WriteGeneratedNode(n)
	s.context.currentCase.mayFallThrough = false
}

// TranspileThrow transpiles marked "throw" statement.
//
// Java: SingleGeneratorFunctionTranspiler.transpileThrow(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileThrow(n rhino.Node) {
	n.AddChildToFront(s.PrepareNodeForWrite(s.MaybeDecomposeExpression(n.RemoveFirstChild())))
	s.context.WriteGeneratedNode(n)
	s.context.currentCase.mayFallThrough = false
}

// ExposeYieldAndTranspileRest exposes YIELD operator so it's free of side
// effects transpiling some code on the way.
//
// Java: SingleGeneratorFunctionTranspiler.exposeYieldAndTranspileRest(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) ExposeYieldAndTranspileRest(n rhino.Node) rhino.Node {
	jutil.CheckState(n.IsGeneratorMarker())
	if n.IsYield() {
		return n
	}

	// Assuming the initial node is "a + (a = b) + (b = yield) + a".

	// YieldExposer may break n up into multiple statements.
	// Place n into a temporary block to hold those statements:
	// {
	//   var JSCompiler_temp_const$jscomp$0 = a + (a = b);
	//   return JSCompiler_temp_const$jscomp$0 + (b = yield) + a;
	// }
	// Need to put expression nodes into return node so that they always stay expression nodes
	// If expression put into expression result YieldExposer may turn it into an "if" statement.
	isExpression := rhino.IR_MayBeExpression(n)
	var block rhino.Node
	if isExpression {
		block = s.outer.astFactory.CreateBlock(s.outer.astFactory.CreateReturn(n))
	} else {
		block = s.outer.astFactory.CreateBlock(n)
	}
	NodeTraversal_Traverse(s.outer.compiler, n, NewEs6RewriteGenerators_YieldExposer(s.outer))
	// Make sure newly created statements are correctly marked for recursive transpileStatement()
	// calls.
	NodeTraversal_Traverse(s.outer.compiler, block, &Es6RewriteGenerators_YieldNodeMarker{})

	// The last child of decomposed block free of side effects.
	decomposed := block.GetLastChild().Detach()
	s.TranspileStatement(block)
	if isExpression {
		return decomposed.RemoveFirstChild()
	}
	return decomposed
}

// MaybeDecomposeExpression converts an expression node containing YIELD into
// an unmarked analogue.
//
// Java: SingleGeneratorFunctionTranspiler.maybeDecomposeExpression(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) MaybeDecomposeExpression(n rhino.Node) rhino.Node {
	if n.IsNil() || !n.IsGeneratorMarker() {
		return n
	}
	return s.TranspileYields(s.ExposeYieldAndTranspileRest(n))
}

// PrepareNodeForWrite makes unmarked node containing arbitrary code suitable
// to write using the TranspilationContext.writeGeneratedNode method.
//
// Java: SingleGeneratorFunctionTranspiler.prepareNodeForWrite(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) PrepareNodeForWrite(n rhino.Node) rhino.Node {
	if n.IsNil() {
		return rhino.NilNode
	}

	// Need to wrap a node so, it can be replaced in the tree with some other node if necessary.
	var wrapper rhino.Node
	if rhino.IR_MayBeStatement(n) {
		wrapper = s.outer.astFactory.CreateBlock(n)
	} else {
		wrapper = s.outer.astFactory.ExprResult(n)
	}
	NodeTraversal_Traverse(s.outer.compiler, wrapper, NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler(s.context))
	jutil.CheckState(wrapper.HasOneChild())
	return wrapper.RemoveFirstChild()
}

// TranspileYields converts node with YIELD into
// $jscomp$generator$context.yieldResult.
//
// Java: SingleGeneratorFunctionTranspiler.transpileYields(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileYields(n rhino.Node) rhino.Node {
	if !n.IsGeneratorMarker() {
		// In some cases exposing yield causes it to disapear from the resulting statement.
		// I.e. the following node: "0 || yield;" becomes:
		// {
		//   var JSCompiler_temp$jscomp$0;
		//   if (JSCompiler_temp$jscomp$0 = 0); else JSCompiler_temp$jscomp$0 = yield;
		// }
		// JSCompiler_temp$jscomp$0; // This is our resulting statement.
		return n
	}
	jumpToSection := s.context.CreateCase()
	yieldNode := s.FindYield(n)
	yieldExpression :=
		s.PrepareNodeForWrite(s.MaybeDecomposeExpression(yieldNode.RemoveFirstChild()))
	if yieldNode.IsYieldAll() {
		s.context.YieldAll(yieldExpression, jumpToSection, yieldNode)
	} else {
		s.context.YieldValue(yieldExpression, jumpToSection, yieldNode)
	}
	s.context.SwitchCaseTo(jumpToSection)
	yieldResult := s.context.YieldResult(yieldNode)
	if yieldNode == n {
		return yieldResult
	}
	// Replace YIELD with $context.yeildResult
	yieldNode.ReplaceWith(yieldResult)
	// Remove generator markings from subtree
	for yieldResult != n {
		yieldResult = yieldResult.GetParent()
		yieldResult.SetGeneratorMarker(false)
	}
	return n
}

// TranspileIf transpiles marked "if" stetement.
//
// Java: SingleGeneratorFunctionTranspiler.transpileIf(Node, TranspilationContext.Case)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileIf(n rhino.Node, breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	astFactory := s.outer.astFactory
	// Decompose condition first
	condition := s.MaybeDecomposeExpression(n.RemoveFirstChild())
	ifBlock := n.GetFirstChild()
	elseBlock := ifBlock.GetNext()

	// No transpilation is needed
	if !ifBlock.IsGeneratorMarker() && (elseBlock.IsNil() || !elseBlock.IsGeneratorMarker()) {
		n.AddChildToFront(condition)
		n.SetGeneratorMarker(false)
		s.TranspileUnmarkedNode(n)
		return
	}

	ifBlock.Detach()
	if elseBlock.IsNil() {
		// No "else" block, just create an empty one as will need it anyway.
		elseBlock = astFactory.CreateBlock().Srcref(n)
	} else {
		elseBlock.Detach()
	}

	// Only "else" block is unmarked, swap "if" and "else" blocks and negate the condition.
	if ifBlock.IsGeneratorMarker() && !elseBlock.IsGeneratorMarker() {
		condition = astFactory.CreateNot(condition).Srcref(condition)
		tmpNode := ifBlock
		ifBlock = elseBlock
		elseBlock = tmpNode
	}

	// Unmarked "if" block (marked "else")
	if !ifBlock.IsGeneratorMarker() {
		endCase := s.context.MaybeCreateCase(breakCase)
		jumpToBlock := s.context.CreateJumpToBlock(endCase /* allowEmbedding= */, false, ifBlock)
		for jumpToBlock.HasChildren() {
			jumpToNode := jumpToBlock.RemoveFirstChild()
			jumpToNode.SetGeneratorSafe(true)
			ifBlock.AddChildToBack(jumpToNode)
		}
		s.TranspileUnmarkedNode(astFactory.CreateIf(condition, ifBlock).Srcref(n))
		s.TranspileStatement(elseBlock)
		s.context.SwitchCaseTo(endCase)
		return
	}

	ifCase := s.context.CreateCase()
	endCase := s.context.MaybeCreateCase(breakCase)

	// "if" and "else" blocks marked
	condition = s.PrepareNodeForWrite(condition)
	newIfBlock := s.context.CreateJumpToBlock(ifCase /* allowEmbedding= */, true, n)
	s.context.WriteGeneratedNode(
		astFactory.CreateIf(s.PrepareNodeForWrite(condition), newIfBlock).Srcref(n))
	s.TranspileStatement(elseBlock)
	s.context.WriteJumpTo(endCase, elseBlock)
	s.context.SwitchCaseTo(ifCase)
	s.TranspileStatement(ifBlock)
	s.context.SwitchCaseTo(endCase)
}

// TranspileFor transpiles marked "for" statement.
//
// Java: SingleGeneratorFunctionTranspiler.transpileFor(Node, TranspilationContext.Case, TranspilationContext.Case)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileFor(
	n rhino.Node,
	breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case,
	continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	astFactory := s.outer.astFactory
	// Decompose init first
	init := s.MaybeDecomposeExpression(n.RemoveFirstChild())
	condition := n.GetFirstChild()
	increment := condition.GetNext()
	body := increment.GetNext()

	// No transpilation is needed
	if !condition.IsGeneratorMarker() &&
		!increment.IsGeneratorMarker() &&
		!body.IsGeneratorMarker() {
		n.AddChildToFront(init)
		n.SetGeneratorMarker(false)
		s.TranspileUnmarkedNode(n)
		return
	}

	// Move init expression out of for loop.
	if !init.IsEmpty() {
		if rhino.IR_MayBeExpression(init) {
			// Convert expression into expression result.
			init = astFactory.ExprResult(init).Srcref(init)
		}
		s.TranspileUnmarkedNode(init)
	}

	startCase := s.context.CreateCase()
	incrementCase := s.context.MaybeCreateCase(continueCase)
	endCase := s.context.MaybeCreateCase(breakCase)

	s.context.SwitchCaseTo(startCase)

	// Transpile condition expression
	if !condition.IsEmpty() {
		condition = s.PrepareNodeForWrite(s.MaybeDecomposeExpression(condition.Detach()))
		s.context.WriteGeneratedNode(
			astFactory.
				CreateIf(
					astFactory.CreateNot(condition).Srcref(condition),
					s.context.CreateJumpToBlock(endCase /* allowEmbedding= */, true, n)).
				Srcref(n))
	}

	// Transpile "for" body
	s.context.PushBreakContinueContext(endCase, incrementCase)
	s.TranspileStatement(body.Detach())
	s.context.PopBreakContinueContext()

	// Transpile increment expression
	s.context.SwitchCaseTo(incrementCase)
	if !increment.IsEmpty() {
		increment = s.MaybeDecomposeExpression(increment.Detach())
		s.TranspileUnmarkedNode(astFactory.ExprResult(increment).Srcref(increment))
	}
	s.context.WriteJumpTo(startCase, n)

	s.context.SwitchCaseTo(endCase)
}

// TranspileForIn transpiles "for in" statement by converting it into "for".
//
// `for (var i in expr) {}` will be converted into
// `for (var i, $for$in = $context.forIn(expr); i = $for$in.getNext(); ) {}`
//
// Java: SingleGeneratorFunctionTranspiler.transpileForIn(Node, TranspilationContext.Case, TranspilationContext.Case)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileForIn(
	n rhino.Node,
	breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case,
	continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	astFactory := s.outer.astFactory
	// Decompose condition first
	detachedExpr := s.MaybeDecomposeExpression(n.GetSecondChild().Detach())
	target := n.GetFirstChild()
	body := n.GetSecondChild()

	// No transpilation is needed
	if !target.IsGeneratorMarker() && !body.IsGeneratorMarker() {
		detachedExpr.InsertAfter(target)
		n.SetGeneratorMarker(false)
		s.TranspileUnmarkedNode(n)
		return
	}

	// Prepare a new init statement
	var init rhino.Node
	if target.Detach().IsVar() {
		// "var i in x"    =>   "var i"
		jutil.CheckState(!target.IsGeneratorMarker())
		init = target
		jutil.CheckState(!init.GetFirstChild().HasChildren())
		target = init.GetFirstChild().CloneNode()
	} else {
		// "i in x"   =>    "var"
		init = rhino.NewNode(rhino.Token_VAR).Srcref(target)
	}

	// "$for$in"
	child := s.context.CallContextMethod(target, "forIn", s.outer.propertyIteratorType, detachedExpr)
	forInCounter := s.forInCounter
	s.forInCounter++
	forIn :=
		astFactory.
			CreateName(
				s.context.GetScopedName(Es6RewriteGenerators_GENERATOR_FORIN_PREFIX)+"$"+strconv.Itoa(forInCounter), AstFactory_Type_(child)).
			Srcref(target)
	// "$context.forIn(x)"
	forIn.AddChildToFront(child)
	// "var ..., $for$in = $context.forIn(expr)"
	init.AddChildToBack(forIn)

	// "$for$in.getNext()"
	forInGetNext :=
		astFactory.
			CreateGetProp(forIn.CloneNode(), "getNext", AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT)).
			Srcref(detachedExpr)

	// "(i = $for$in.getNext()) != null"
	forCond :=
		astFactory.
			CreateNe(
				astFactory.
					CreateAssign(
						target.SetColor(s.outer.nullableStringType),
						astFactory.
							CreateCall(forInGetNext, AstFactory_Type_OfColor(s.outer.nullableStringType)).
							Srcref(detachedExpr)).
					Srcref(detachedExpr),
				astFactory.CreateNull().Srcref(forIn)).
			Srcref(detachedExpr)
	forCond.SetGeneratorMarker(target.IsGeneratorMarker())

	// Prepare "for" statement.
	// "for (var i, $for$in = $context.forIn(expr); (i = $for$in.getNext()) != null; ) {}"
	forNode := rhino.IR_ForNode(init, forCond, rhino.IR_Empty().Srcref(n), body.Detach()).Srcref(n)

	// Transpile "for" instead of "for in".
	s.TranspileFor(forNode, breakCase, continueCase)
}

// TranspileDo transpiles "do while" statement.
//
// Java: SingleGeneratorFunctionTranspiler.transpileDo(Node, TranspilationContext.Case, TranspilationContext.Case)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileDo(
	n rhino.Node,
	breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case,
	continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	startCase := s.context.CreateCase()
	breakCase = s.context.MaybeCreateCase(breakCase)
	continueCase = s.context.MaybeCreateCase(continueCase)

	s.context.SwitchCaseTo(startCase)

	// Transpile body
	body := n.RemoveFirstChild()
	s.context.PushBreakContinueContext(breakCase, continueCase)
	s.TranspileStatement(body)
	s.context.PopBreakContinueContext()

	// Transpile condition
	s.context.SwitchCaseTo(continueCase)
	condition := s.PrepareNodeForWrite(s.MaybeDecomposeExpression(n.RemoveFirstChild()))
	s.context.WriteGeneratedNode(
		s.outer.astFactory.
			CreateIf(
				condition, s.context.CreateJumpToBlock(startCase /* allowEmbedding= */, false, n)).
			Srcref(n))
	s.context.SwitchCaseTo(breakCase)
}

// TranspileTry transpiles "try" statement
//
// Java: SingleGeneratorFunctionTranspiler.transpileTry(Node, TranspilationContext.Case)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileTry(n rhino.Node, breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	tryBlock := n.RemoveFirstChild()
	catchBlock := n.RemoveFirstChild()
	finallyBlock := n.RemoveFirstChild()

	var catchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case
	if catchBlock.HasChildren() {
		catchCase = s.context.CreateCase()
	}
	var finallyCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case
	if !finallyBlock.IsNil() {
		finallyCase = s.context.CreateCase()
	}
	endCase := s.context.MaybeCreateCase(breakCase)

	// Transpile "try" block
	s.context.EnterTryBlock(catchCase, finallyCase, tryBlock)
	s.TranspileStatement(tryBlock)

	if finallyBlock.IsNil() {
		s.context.LeaveTryBlock(catchCase, endCase, tryBlock)
	} else {
		// Transpile "finally" block
		s.context.SwitchCaseTo(finallyCase)
		s.context.EnterFinallyBlock(catchCase, finallyCase, finallyBlock)
		s.TranspileStatement(finallyBlock)
		s.context.LeaveFinallyBlock(endCase, finallyBlock)
	}

	// Transpile "catch" block
	if catchBlock.HasChildren() {
		jutil.CheckState(catchBlock.GetFirstChild().IsCatch())

		s.context.SwitchCaseTo(catchCase)
		exceptionName := catchBlock.GetFirstFirstChild().Detach()
		s.context.EnterCatchBlock(finallyCase, exceptionName)

		catchBody := catchBlock.GetFirstFirstChild().Detach()
		jutil.CheckState(catchBody.IsBlock())
		s.TranspileStatement(catchBody)
		s.context.LeaveCatchBlock(finallyCase, catchBody)
	}

	s.context.SwitchCaseTo(endCase)
}

// TranspileSwitch transpiles "switch" statement.
//
// Java: SingleGeneratorFunctionTranspiler.transpileSwitch(Node, TranspilationContext.Case)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) TranspileSwitch(n rhino.Node, breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	// Transpile condition first
	n.AddChildToFront(s.MaybeDecomposeExpression(n.RemoveFirstChild()))

	// Are all "switch" cases unmarked?
	hasGeneratorMarker := false
	switchBody := n.GetSecondChild()
	for caseSection := switchBody.GetFirstChild(); !caseSection.IsNil(); caseSection = caseSection.GetNext() {
		if caseSection.IsGeneratorMarker() {
			hasGeneratorMarker = true
			break
		}
	}
	// No transpilation is needed
	if !hasGeneratorMarker {
		n.SetGeneratorMarker(false)
		s.TranspileUnmarkedNode(n)
		return
	}

	// Stores a detached body of a case statement and a case section assosiated with it.
	type SwitchCase struct {
		generatedCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case
		body          rhino.Node
	}

	// TODO(skill): Don't move all case sections.
	detachedCases := []*SwitchCase{}

	// We don't have to transpile unmarked cases at the beginning of "switch".
	canSkipUnmarkedCases := true
	for caseSection := switchBody.GetFirstChild(); !caseSection.IsNil(); caseSection = caseSection.GetNext() {
		if !caseSection.IsDefaultCase() && caseSection.GetFirstChild().IsGeneratorMarker() {
			// Following example is possible to transpile, but it's not trivial.
			// switch (cond) {
			//   case yield "test": break;
			//   case 5 + yield:  break;
			// }
			s.outer.compiler.Report(
				JSError_MakeWithNode(
					n, TranspilationUtil_CANNOT_CONVERT_YET, "Case statements that contain yields"))
			return
		}
		body := caseSection.GetLastChild()

		if !body.HasChildren() || (canSkipUnmarkedCases && !body.IsGeneratorMarker()) {
			// Can skip empty or unmarked case.
			continue
		}

		// Check whether we can start skipping unmarked cases again
		canSkipUnmarkedCases = s.IsEndOfBlockUnreachable(body)

		// Move case's body under a global switch statement...

		// Allocate a new case
		generatedCase := s.context.CreateCase()
		generatedCase.caseBlock.Srcref(body)

		// Replace old body with a jump instruction.
		newBody := rhino.IR_BlockOfNode(s.context.CreateJumpToNode(generatedCase, body))
		newBody.SetIsAddedBlock(true)
		newBody.SetGeneratorSafe(true) // make sure we don't transpile generated "jump" instruction
		body.ReplaceWith(newBody)

		// Remember the body and the case under which the body will be moved.
		detachedCases = append(detachedCases, &SwitchCase{generatedCase: generatedCase, body: body})

		caseSection.SetGeneratorMarker(false)
	}

	endCase := s.context.MaybeCreateCase(breakCase)

	// Transpile the barebone of original "switch" statement
	n.SetGeneratorMarker(false)
	switchBody.SetGeneratorMarker(false)
	s.TranspileUnmarkedNode(n)
	s.context.WriteJumpTo(endCase, n) // TODO(skill): do not always add this.

	// Transpile all detached case bodies
	s.context.PushBreakContext(endCase)
	for _, detachedCase := range detachedCases {
		generatedCase := detachedCase.generatedCase
		s.context.SwitchCaseTo(generatedCase)
		s.TranspileStatement(detachedCase.body)
	}
	s.context.PopBreakContext()

	s.context.SwitchCaseTo(endCase)
}

// FindYield finds the only YIELD node in a tree.
//
// Java: SingleGeneratorFunctionTranspiler.findYield(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) FindYield(n rhino.Node) rhino.Node {
	yieldFinder := &Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder{}
	NodeTraversal_Traverse(s.outer.compiler, n, yieldFinder)
	return yieldFinder.GetYieldNode()
}

// Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder finds
// the only YIELD node in a tree.
type Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder struct {
	NodeTraversal_AbstractPreOrderCallbackBase

	yieldNode rhino.Node
}

// Java: YieldFinder.shouldTraverse(NodeTraversal, Node, Node)
func (y *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsFunction() {
		return false
	}
	if n.IsYield() {
		jutil.CheckState(y.yieldNode.IsNil())
		y.yieldNode = n
		return false
	}
	return true
}

// Java: YieldFinder.getYieldNode()
func (y *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_YieldFinder) GetYieldNode() rhino.Node {
	rhino.CheckNotNull(y.yieldNode)
	return y.yieldNode
}

// IsEndOfBlockUnreachable returns whether any statements added to the end of
// the block would be unreachable.
//
// It's OK for this method to return false-negatives.
//
// Java: SingleGeneratorFunctionTranspiler.isEndOfBlockUnreachable(Node)
func (s *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) IsEndOfBlockUnreachable(block rhino.Node) bool {
	jutil.CheckState(block.IsBlock())
	if !block.HasChildren() {
		return false
	}
	switch block.GetLastChild().GetToken() {
	case rhino.Token_BLOCK:
		return s.IsEndOfBlockUnreachable(block.GetLastChild())
	case rhino.Token_RETURN, rhino.Token_THROW, rhino.Token_CONTINUE, rhino.Token_BREAK:
		return true
	default:
		return false
	}
}

// Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext
// is the state machine context that is used during generator function
// transpilation.
type Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext struct {
	outer *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler

	namedLabels *jutil.LinkedHashMap[string, *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases]
	// Go-only: language-adapter — each ArrayDeque used as a stack is a slice
	// whose last element is its first (push/pop/addFirst/removeFirst/getFirst).
	breakCases    []*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case
	continueCases []*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case

	catchCases   []*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase
	finallyCases []*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case
	catchNames   *jutil.LinkedHashSet[string]

	// All "case" sections that will be added to generator program.
	allCases []*Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case

	// All "break" nodes that exit from the generator primary switch statement
	switchBreaks []rhino.Node

	// A virtual case that indicates end of program
	programEndCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case

	// Most recently assigned id.
	caseIdCounter int

	// Points to the switch case that is being populated with transpiled instructions from the
	// original generator function that is being transpiled.
	currentCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case

	// A counter for the number of finally blocks we are currently inside.
	// This value is used for two purposes:
	// 1. At COMPILE-TIME, to determine if a break/continue is inside a finally block.
	// 2. At RUNTIME, its value is emitted into the generated code to manage the
	//    exception-handling stack.
	nestedFinallyBlockCount int
	thisReferenceFound      bool
	argumentsReferenceFound bool
}

// Java: TranspilationContext()
func NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext(outer *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler) *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext {
	c := &Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext{outer: outer}
	c.namedLabels = jutil.NewLinkedHashMap[string, *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases]()
	c.catchNames = jutil.NewLinkedHashSet[string]()
	c.nestedFinallyBlockCount = 0
	c.programEndCase = NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case(c)
	jutil.CheckState(c.programEndCase.id == 0)
	c.currentCase = NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case(c)
	jutil.CheckState(c.currentCase.id == 1)
	c.allCases = append(c.allCases, c.currentCase)
	return c
}

// OptimizeCaseIds removes unnecessary cases.
//
// This optimization is needed to reduce number of switch cases, which is used
// then to generate even shorter state machine programs.
//
// Java: TranspilationContext.optimizeCaseIds()
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) OptimizeCaseIds() {
	if len(c.allCases) == 0 {
		jutil.CheckStateMsg1(false, "%s", jutil.CollectionToString(c.allCases))
	}

	// Shortcut jump chains:
	//   case 100:
	//     $context.yieldValue("something", 101);
	//     break;
	//   case 101:
	//     $context.jumpTo(102);
	//     break;
	//   case 102:
	//     $context.jumpTo(200);
	//     break;
	// becomes:
	//   case 100:
	//     $context.yieldValue("something", 200);
	//     break;
	//   case 101:
	//     $context.jumpTo(102);
	//     break;
	//   case 102:
	//     $context.jumpTo(200);
	//     break;
	for _, currentCase := range c.allCases {
		if currentCase.jumpTo != nil {
			// Flatten jumps chains:
			// 1 -> 2
			// 2 -> 8
			// 8 -> 300
			// to:
			// 1 -> 300
			// 2 -> 300
			// 8 -> 300
			for currentCase.jumpTo.jumpTo != nil {
				currentCase.jumpTo = currentCase.jumpTo.jumpTo
			}

			if !currentCase.embedInto.IsNil() && len(currentCase.references) == 1 {
				currentCase.jumpTo.embedInto = currentCase.embedInto
			}
			currentCase.embedInto = rhino.NilNode

			// Update references to jump to the final case in the chain
			for _, reference := range currentCase.references {
				reference.SetDouble(float64(currentCase.jumpTo.id))
			}
			currentCase.jumpTo.references = append(currentCase.jumpTo.references, currentCase.references...)
			currentCase.references = currentCase.references[:0]
		}
	}

	// Merge cases without any references with the previous case:
	//   case 100:
	//     doSomething();
	//   case 101:
	//     doSomethingElse();
	//     break;
	//   case 102:
	//     doEvenMore();
	// becomes:
	//   case 100:
	//     doSomething();
	//     doSomethingElse();
	//     break;
	//
	// Go-only: language-adapter — the Iterator with it.remove() is a filter
	// that keeps the cases not removed, in order.
	kept := c.allCases[:1]
	prevCase := c.allCases[0]
	jutil.CheckState(prevCase.id == 1)
	for _, currentCase := range c.allCases[1:] {
		if len(currentCase.references) == 0 {
			// No jump references, just append the body to a previous case if needed.
			jutil.CheckState(currentCase.embedInto.IsNil())
			if prevCase.mayFallThrough {
				prevCase.caseBlock.AddChildrenToBack(currentCase.caseBlock.RemoveChildren())
				prevCase.mayFallThrough = currentCase.mayFallThrough
			}
			continue // it.remove()
		}
		if !currentCase.embedInto.IsNil() {
			jutil.CheckState(currentCase.jumpTo == nil)
			// Cases can be embedded only if they are referenced once and don't fall through.
			if len(currentCase.references) == 1 && !currentCase.mayFallThrough {
				currentCase.embedInto.ReplaceWith(currentCase.caseBlock)
				continue // it.remove()
			}
		}
		if prevCase.jumpTo == currentCase {
			// Merging "case 1:" with the following case. The standard merging cannot be used
			// as "case 1:" is an entry point and it cannot be renamed.
			//   case 1:
			//   case 2:
			//     doSomethingElse();
			//     break;
			//   case 102:
			//     $context.jumpTo(2);
			//     break;
			// becomes:
			//   case 1:
			//     doSomethingElse();
			//     break;
			//   case 102:
			//     $context.jumpTo(1);
			//     break;
			jutil.CheckState(prevCase.mayFallThrough)
			jutil.CheckState(!prevCase.caseBlock.HasChildren())
			jutil.CheckState(currentCase.jumpTo == nil)

			prevCase.caseBlock.AddChildrenToBack(currentCase.caseBlock.RemoveChildren())
			prevCase.mayFallThrough = currentCase.mayFallThrough
			for _, reference := range currentCase.references {
				reference.SetDouble(float64(prevCase.id))
			}
			prevCase.jumpTo = currentCase.jumpTo
			prevCase.references = append(prevCase.references, currentCase.references...)
			continue // it.remove()
		}
		kept = append(kept, currentCase)
		prevCase = currentCase
	}
	clear(c.allCases[len(kept):])
	c.allCases = kept
}

// EliminateSwitchBreaks replaces "...; break;" with "return ...;".
//
// Java: TranspilationContext.eliminateSwitchBreaks()
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) EliminateSwitchBreaks() {
	for _, breakNode := range c.switchBreaks {
		prevStatement := breakNode.GetPrevious()
		jutil.CheckState(!prevStatement.IsNil())
		jutil.CheckState(prevStatement.IsExprResult())
		prevStatement.ReplaceWith(rhino.IR_ReturnNodeWithExpr(prevStatement.RemoveFirstChild()))
		breakNode.Detach()
	}
	c.switchBreaks = nil
}

// allCasesRemoveFirst is allCases.remove(0).
//
// Go-only: language-adapter — ArrayList.remove(0) on the slice.
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) allCasesRemoveFirst() *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case {
	first := c.allCases[0]
	c.allCases = c.allCases[1:]
	return first
}

// FinalizeTransformation finalizes transpilation by dumping all generated
// "case" nodes.
//
// Java: TranspilationContext.finalizeTransformation(Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) FinalizeTransformation(generatorBody rhino.Node) {
	astFactory := c.outer.outer.astFactory
	c.OptimizeCaseIds()

	// If number of cases is small we render them without using "switch"
	//   switch ($context.getNextAddressJsc()) {
	//     case 1: a();
	//     case 2: b();
	//     case 3: c();
	//   }
	// are rendered as:
	//   if ($context.getNextAddressJsc() == 1) a();
	//   if ($context.getNextAddressJsc() != 3) b();
	//   c();
	if len(c.allCases) == 2 || len(c.allCases) == 3 {
		generatorBody.AddChildToBack(
			astFactory.
				CreateIf(
					astFactory.CreateEq(
						c.CallContextMethod(
							generatorBody,
							"getNextAddressJsc",
							AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_NUMBER_TYPE, colors.StandardColors_NUMBER)),
						astFactory.CreateNumber(1)),
					c.allCasesRemoveFirst().caseBlock).
				SrcrefTreeIfMissing(generatorBody))
	}

	// If number of cases is small we render them without using "switch"
	//   switch ($context.getNextAddressJsc()) {
	//     case 1: a();
	//     case 2: b();
	//   }
	// are rendered as:
	//   if ($context.getNextAddressJsc() == 1) a();
	//   b();
	if len(c.allCases) == 2 {
		generatorBody.AddChildToBack(
			astFactory.
				CreateIf(
					astFactory.CreateNe(
						c.CallContextMethod(
							generatorBody,
							"getNextAddressJsc",
							AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_NUMBER_TYPE, colors.StandardColors_NUMBER)),
						astFactory.CreateNumber(float64(c.allCases[1].id))),
					c.allCasesRemoveFirst().caseBlock).
				SrcrefTreeIfMissing(generatorBody))
	}

	// If number of cases is small we render them without using "switch"
	//   switch ($context.getNextAddressJsc()) {
	//     case 1: a();
	//   }
	// are rendered as:
	//   a();
	if len(c.allCases) == 1 {
		generatorBody.AddChildrenToBack(c.allCasesRemoveFirst().caseBlock.RemoveChildren())
		c.EliminateSwitchBreaks()
		return
	}

	//  switch ($jscomp$generator$context.getNextAddressJsc()) {}
	switchNode :=
		rhino.IR_SwitchNode(
			c.CallContextMethod(
				generatorBody,
				"getNextAddressJsc",
				AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_NUMBER_TYPE, colors.StandardColors_NUMBER))).
			Srcref(generatorBody)
	generatorBody.AddChildToBack(switchNode)
	switchBody := switchNode.GetSecondChild().Srcref(generatorBody)

	// Populate "switch" statement with "case"s.
	for _, currentCase := range c.allCases {
		switchBody.AddChildToBack(currentCase.CreateCaseNode())
	}
	c.allCases = nil
}

// CheckStateIsEmpty ensures that the context has an empty state.
//
// Java: TranspilationContext.checkStateIsEmpty()
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) CheckStateIsEmpty() {
	jutil.CheckState(c.namedLabels.IsEmpty())
	jutil.CheckState(len(c.breakCases) == 0)
	jutil.CheckState(len(c.continueCases) == 0)
	jutil.CheckState(len(c.catchCases) == 0)
	jutil.CheckState(len(c.finallyCases) == 0)
	jutil.CheckState(c.nestedFinallyBlockCount == 0)
	jutil.CheckState(len(c.allCases) == 0)
}

// TranspileUnmarkedBlock adds a block of original code to the end of the
// current case.
//
// Java: TranspilationContext.transpileUnmarkedBlock(Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) TranspileUnmarkedBlock(block rhino.Node) {
	if block.HasChildren() {
		NodeTraversal_Traverse(c.outer.outer.compiler, block, NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler(c))
		for block.HasChildren() {
			c.WriteGeneratedNode(block.RemoveFirstChild())
		}
	}
}

// WriteGeneratedNode adds a new generated node to the end of the current
// case.
//
// Java: TranspilationContext.writeGeneratedNode(Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) WriteGeneratedNode(n rhino.Node) {
	c.currentCase.AddNode(n)
}

// WriteGeneratedNodeAndBreak adds a new generated node to the end of the
// current case and finializes it.
//
// Java: TranspilationContext.writeGeneratedNodeAndBreak(Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) WriteGeneratedNodeAndBreak(n rhino.Node) {
	c.WriteGeneratedNode(n)
	c.WriteGeneratedNode(c.CreateBreakNodeFor(n))
	c.currentCase.mayFallThrough = false
}

// CreateCase creates a new detached case statement.
//
// Java: TranspilationContext.createCase()
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) CreateCase() *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case {
	return NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case(c)
}

// MaybeCreateCase returns a passed case object or creates a new one if it's
// null.
//
// Java: TranspilationContext.maybeCreateCase(Case)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) MaybeCreateCase(other *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case {
	if other != nil {
		return other
	}
	return c.CreateCase()
}

// GetJsContextNameNode returns the name node of context parameter passed to
// the program.
//
// Java: TranspilationContext.getJsContextNameNode(Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) GetJsContextNameNode(sourceNode rhino.Node) rhino.Node {
	return c.outer.outer.astFactory.
		CreateName(c.GetScopedName(Es6RewriteGenerators_GENERATOR_CONTEXT), c.outer.outer.generatorContextType()).
		Srcref(sourceNode)
}

// GetScopedName returns unique name in the current context.
//
// Java: TranspilationContext.getScopedName(String)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) GetScopedName(name string) string {
	return name + "$" + c.outer.uniqueId
}

// GetContextField creates node that access a specified field of the current
// context.
//
// Java: TranspilationContext.getContextField(Node, String)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) GetContextField(sourceNode rhino.Node, fieldName string) rhino.Node {
	return c.outer.outer.astFactory.
		CreateGetPropWithUnknownType(c.GetJsContextNameNode(sourceNode), fieldName).
		Srcref(sourceNode)
}

// CallContextMethod creates node that make a call to a context function.
//
// Java: TranspilationContext.callContextMethod(Node, String, AstFactory.Type, Node...)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) CallContextMethod(
	sourceNode rhino.Node, methodName string, type_ AstFactory_Type, args ...rhino.Node) rhino.Node {
	contextField := c.GetContextField(sourceNode, methodName)
	return c.outer.outer.astFactory.CreateCall(contextField, type_, args...).Srcref(sourceNode)
}

// CallContextMethodResult creates node that make a call to a context
// function.
//
// Java: TranspilationContext.callContextMethodResult(Node, String, AstFactory.Type, Node...)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) CallContextMethodResult(
	sourceNode rhino.Node, methodName string, type_ AstFactory_Type, args ...rhino.Node) rhino.Node {
	return c.outer.outer.astFactory.
		ExprResult(c.CallContextMethod(sourceNode, methodName, type_, args...)).
		Srcref(sourceNode)
}

// ReturnContextMethod creates node that returns the result of a call to a
// context function.
//
// Java: TranspilationContext.returnContextMethod(Node, String, AstFactory.Type, Node...)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) ReturnContextMethod(
	sourceNode rhino.Node, methodName string, type_ AstFactory_Type, args ...rhino.Node) rhino.Node {
	return c.outer.outer.astFactory.
		CreateReturn(c.CallContextMethod(sourceNode, methodName, type_, args...)).
		Srcref(sourceNode)
}

// CreateBreakNodeFor creates a "break;" statement that will follow preBreak
// node.
//
// This is used to be able to generatate a state machine program outside of
// "swtich" statement so:
//
//	$context.jumpTo(5);
//	break;
//
// could be converted into:
//
//	return $context.jumpTo(5);
//
// Java: TranspilationContext.createBreakNodeFor(Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) CreateBreakNodeFor(preBreak rhino.Node) rhino.Node {
	breakNode := rhino.IR_BreakNode().Srcref(preBreak)
	c.switchBreaks = append(c.switchBreaks, breakNode)
	return breakNode
}

// CreateJumpToNode returns a node that instructs a state machine program to
// jump to a selected case section.
//
// Java: TranspilationContext.createJumpToNode(Case, Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) CreateJumpToNode(section *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode rhino.Node) rhino.Node {
	return c.ReturnContextMethod(
		sourceNode, "jumpTo", AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID), section.GetNumber(sourceNode))
}

// WriteJumpTo instructs a state machine program to jump to a selected case
// section.
//
// Java: TranspilationContext.writeJumpTo(Case, Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) WriteJumpTo(section *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode rhino.Node) {
	c.currentCase.JumpTo(
		section, c.CreateJumpToBlock(section /* allowEmbedding= */, false, sourceNode))
}

// CreateJumpToBlock creates a block node that contains a jump instruction.
//
// allowEmbedding tells whether the code from the target section can be
// embedded into jump block.
//
// Java: TranspilationContext.createJumpToBlock(Case, boolean, Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) CreateJumpToBlock(section *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, allowEmbedding bool, sourceNode rhino.Node) rhino.Node {
	jutil.CheckState(section.embedInto.IsNil())
	jumpBlock :=
		c.outer.outer.astFactory.
			CreateBlock(
				c.CallContextMethodResult(
					sourceNode,
					"jumpTo",
					AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID),
					section.GetNumber(sourceNode)),
				c.CreateBreakNodeFor(sourceNode)).
			Srcref(sourceNode)
	if allowEmbedding {
		section.embedInto = jumpBlock
	}
	return jumpBlock
}

// ReplaceBreakContinueWithJump converts "break" and "continue" statements
// into state machine jumps.
//
// Java: TranspilationContext.replaceBreakContinueWithJump(Node, Case, int)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) ReplaceBreakContinueWithJump(sourceNode rhino.Node, section *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, breakSuppressors int) {
	var jumpMethod string
	if c.nestedFinallyBlockCount > 0 {
		// If we are in a finally block, we need to use jumpThroughFinallyBlocks to ensure that
		// the finally block is correctly exited.
		jumpMethod = "jumpThroughFinallyBlocks"
	} else if len(c.finallyCases) == 0 || c.finallyCasesGetFirst().id < section.id {
		// There are no finally blocks that should be exectuted pior to jumping
		jumpMethod = "jumpTo"
	} else {
		// There are some finally blocks that should be exectuted before we can break/continue.
		jutil.CheckState(c.finallyCasesGetFirst().id != section.id)
		jumpMethod = "jumpThroughFinallyBlocks"
	}
	if breakSuppressors == 0 {
		// continue;  =>  $context.jumpTo(x); break;
		c.CallContextMethodResult(
			sourceNode,
			jumpMethod,
			AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID),
			section.GetNumber(sourceNode)).
			InsertBefore(sourceNode)
		sourceNode.ReplaceWith(c.CreateBreakNodeFor(sourceNode))
	} else {
		// "break;" inside a loop or swtich statement:
		// for (...) {
		//   break l1;
		// }
		// becomes:
		// for (...) {                  // loop doesn't allow to use "break" to advance to the
		//   return $context.jumpTo(x); // next address, so "return" is used instead.
		// }
		sourceNode.ReplaceWith(
			c.ReturnContextMethod(
				sourceNode,
				jumpMethod,
				AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID),
				section.GetNumber(sourceNode)))
	}
}

// YieldValue instructs a state machine program to yield a value and then
// jump to a selected case section.
//
// Java: TranspilationContext.yieldValue(Node, TranspilationContext.Case, Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) YieldValue(
	expression rhino.Node, jumpToSection *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode rhino.Node) {
	args := []rhino.Node{}
	if expression.IsNil() {
		args = append(args, c.outer.outer.astFactory.CreateUndefinedValue().SrcrefTree(sourceNode))
	} else {
		args = append(args, expression)
	}
	args = append(args, jumpToSection.GetNumber(sourceNode))
	c.outer.context.WriteGeneratedNode(
		c.ReturnContextMethod(
			sourceNode, "yield", AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN), args...))
	c.outer.context.currentCase.mayFallThrough = false
}

// YieldAll instructs a state machine program to yield all values and then
// jump to a selected case section.
//
// Java: TranspilationContext.yieldAll(Node, TranspilationContext.Case, Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) YieldAll(expression rhino.Node, jumpToSection *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode rhino.Node) {
	c.WriteGeneratedNode(
		c.ReturnContextMethod(
			sourceNode,
			"yieldAll",
			AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN),
			expression,
			jumpToSection.GetNumber(sourceNode)))
	c.outer.context.currentCase.mayFallThrough = false
}

// ReturnExpression instructs a state machine program to return a given
// expression.
//
// Java: TranspilationContext.returnExpression(Node, Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) ReturnExpression(sourceNode rhino.Node, expression rhino.Node) rhino.Node {
	if expression.IsNil() {
		return c.CallContextMethod(sourceNode, "return", AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID))
	}
	return c.CallContextMethod(
		sourceNode, "return", AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID), expression)
}

// YieldResult instructs a state machine program to consume a yield result
// after yielding.
//
// Java: TranspilationContext.yieldResult(Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) YieldResult(sourceNode rhino.Node) rhino.Node {
	return c.CallContextMethod(sourceNode, "getYieldResultJsc", AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN))
}

// catchCasesGetFirst is ArrayDeque.getFirst(), or NoSuchElementException.
//
// Go-only: language-adapter — the ArrayDeque operation on the slice stack.
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) catchCasesGetFirst() *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase {
	if len(c.catchCases) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	return c.catchCases[len(c.catchCases)-1]
}

// catchCasesRemoveFirst is ArrayDeque.removeFirst(), or
// NoSuchElementException.
//
// Go-only: language-adapter — the ArrayDeque operation on the slice stack.
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) catchCasesRemoveFirst() *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase {
	first := c.catchCasesGetFirst()
	c.catchCases = c.catchCases[:len(c.catchCases)-1]
	return first
}

// finallyCasesGetFirst is ArrayDeque.getFirst(), or NoSuchElementException.
//
// Go-only: language-adapter — the ArrayDeque operation on the slice stack.
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) finallyCasesGetFirst() *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case {
	if len(c.finallyCases) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	return c.finallyCases[len(c.finallyCases)-1]
}

// finallyCasesRemoveFirst is ArrayDeque.removeFirst(), or
// NoSuchElementException.
//
// Go-only: language-adapter — the ArrayDeque operation on the slice stack.
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) finallyCasesRemoveFirst() *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case {
	first := c.finallyCasesGetFirst()
	c.finallyCases = c.finallyCases[:len(c.finallyCases)-1]
	return first
}

// AddCatchFinallyCases adds references to catch and finally blocks to the
// transpilation context.
//
// Java: TranspilationContext.addCatchFinallyCases(Case, Case)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) AddCatchFinallyCases(catchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, finallyCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	if finallyCase != nil {
		if len(c.catchCases) != 0 {
			c.catchCasesGetFirst().finallyBlocks++
		}
		c.finallyCases = append(c.finallyCases, finallyCase)
	}
	if catchCase != nil {
		c.catchCases = append(c.catchCases, NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase(catchCase))
	}
}

// GetNextCatchCase returns the case section of the next catch block that is
// not hidden by finally blocks.
//
// Java: TranspilationContext.getNextCatchCase()
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) GetNextCatchCase() *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case {
	for i := len(c.catchCases) - 1; i >= 0; i-- {
		catchCase := c.catchCases[i]
		if catchCase.finallyBlocks == 0 {
			return catchCase.catchCase
		}
		break
	}
	return nil
}

// GetNextFinallyCase returns the case section of the next finally block.
//
// Java: TranspilationContext.getNextFinallyCase()
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) GetNextFinallyCase() *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case {
	if len(c.finallyCases) == 0 {
		return nil
	}
	return c.finallyCasesGetFirst()
}

// RemoveCatchFinallyCases removes references to catch and finally blocks from
// the transpilation context.
//
// Java: TranspilationContext.removeCatchFinallyCases(Case, Case)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) RemoveCatchFinallyCases(catchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, finallyCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	if catchCase != nil {
		lastCatch := c.catchCasesRemoveFirst()
		jutil.CheckState(lastCatch.finallyBlocks == 0)
		jutil.CheckState(lastCatch.catchCase == catchCase)
	}
	if finallyCase != nil {
		if len(c.catchCases) != 0 {
			c.catchCasesGetFirst().finallyBlocks--
			finallyBlocks := c.catchCasesGetFirst().finallyBlocks
			jutil.CheckState(finallyBlocks >= 0)
		}
		lastFinally := c.finallyCasesRemoveFirst()
		jutil.CheckState(lastFinally == finallyCase)
	}
}

// EnterTryBlock writes a statement Node that should be placed at the
// beginning of try block.
//
// Java: TranspilationContext.enterTryBlock(Case, Case, Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) EnterTryBlock(catchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, finallyCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode rhino.Node) {
	c.AddCatchFinallyCases(catchCase, finallyCase)

	var methodName string
	args := []rhino.Node{}
	if catchCase == nil {
		methodName = "setFinallyBlock"
		args = append(args, finallyCase.GetNumber(sourceNode))
	} else {
		methodName = "setCatchFinallyBlocks"
		args = append(args, catchCase.GetNumber(sourceNode))
		if finallyCase != nil {
			args = append(args, finallyCase.GetNumber(sourceNode))
		}
	}
	c.WriteGeneratedNode(
		c.CallContextMethodResult(
			sourceNode,
			methodName,
			AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID),
			args...))
}

// LeaveTryBlock writes a statements that should be placed at the end of try
// block if finally block is not present.
//
// Java: TranspilationContext.leaveTryBlock(Case, Case, Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) LeaveTryBlock(catchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, endCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode rhino.Node) {
	c.RemoveCatchFinallyCases(catchCase, nil)
	args := []rhino.Node{}
	args = append(args, endCase.GetNumber(sourceNode))
	// Find the next catch block that is not hidden by any finally blocks.
	nextCatchCase := c.GetNextCatchCase()
	if nextCatchCase != nil {
		args = append(args, nextCatchCase.GetNumber(sourceNode))
	}
	c.WriteGeneratedNodeAndBreak(
		c.CallContextMethodResult(
			sourceNode,
			"leaveTryBlock",
			AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID),
			args...))
}

// EnterCatchBlock writes a statement Node that should be placed at the
// beginning of catch block.
//
// Java: TranspilationContext.enterCatchBlock(Case, Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) EnterCatchBlock(finallyCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, exceptionName rhino.Node) {
	astFactory := c.outer.outer.astFactory
	jutil.CheckState(exceptionName.IsName())
	c.AddCatchFinallyCases(nil, finallyCase)

	// Find the next catch block that is not hidden by any finally blocks.
	nextCatchCase := c.GetNextCatchCase()

	if c.catchNames.Add(exceptionName.GetString()) {
		c.outer.HoistVarNode(rhino.IR_Var(exceptionName.CloneNode()).Srcref(exceptionName))
	}

	args := []rhino.Node{}
	if nextCatchCase != nil {
		args = append(args, nextCatchCase.GetNumber(exceptionName))
	}

	enterCatchBlockCall :=
		c.CallContextMethod(
			exceptionName,
			"enterCatchBlock",
			AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN),
			args...)
	exceptionName.SetColor(enterCatchBlockCall.GetColor())
	c.WriteGeneratedNode(
		astFactory.
			ExprResult(
				astFactory.
					CreateAssign(exceptionName, enterCatchBlockCall).
					Srcref(exceptionName)).
			Srcref(exceptionName))
}

// LeaveCatchBlock writes a statement to jump to the finally block if it's
// present.
//
// Java: TranspilationContext.leaveCatchBlock(Case, Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) LeaveCatchBlock(finallyCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode rhino.Node) {
	if finallyCase != nil {
		c.RemoveCatchFinallyCases(nil, finallyCase)
		c.WriteJumpTo(finallyCase, sourceNode)
	}
}

// EnterFinallyBlock writes a Node that should be placed at the beginning of
// finally block.
//
// Java: TranspilationContext.enterFinallyBlock(Case, Case, Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) EnterFinallyBlock(
	catchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, finallyCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode rhino.Node) {
	astFactory := c.outer.outer.astFactory
	c.RemoveCatchFinallyCases(catchCase, finallyCase)

	nextCatchCase := c.GetNextCatchCase()
	nextFinallyCase := c.GetNextFinallyCase()

	args := []rhino.Node{}
	if c.nestedFinallyBlockCount == 0 {
		if nextCatchCase != nil || nextFinallyCase != nil {
			if nextCatchCase == nil {
				args = append(args, astFactory.CreateNumber(0).Srcref(sourceNode))
			} else {
				args = append(args, nextCatchCase.GetNumber(sourceNode))
			}
			if nextFinallyCase != nil {
				args = append(args, nextFinallyCase.GetNumber(sourceNode))
			}
		}
	} else {
		if nextCatchCase == nil {
			args = append(args, astFactory.CreateNumber(0).Srcref(sourceNode))
		} else {
			args = append(args, nextCatchCase.GetNumber(sourceNode))
		}
		if nextFinallyCase == nil {
			args = append(args, astFactory.CreateNumber(0).Srcref(sourceNode))
		} else {
			args = append(args, nextFinallyCase.GetNumber(sourceNode))
		}
		args = append(args, astFactory.CreateNumber(float64(c.nestedFinallyBlockCount)).Srcref(sourceNode))
	}

	c.WriteGeneratedNode(
		c.CallContextMethodResult(
			sourceNode,
			"enterFinallyBlock",
			AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID),
			args...))

	c.nestedFinallyBlockCount++
}

// LeaveFinallyBlock writes a Node that should be placed at the end of
// finally block.
//
// Java: TranspilationContext.leaveFinallyBlock(Case, Node)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) LeaveFinallyBlock(endCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, sourceNode rhino.Node) {
	args := []rhino.Node{}
	args = append(args, endCase.GetNumber(sourceNode))
	c.nestedFinallyBlockCount--
	if c.nestedFinallyBlockCount != 0 {
		args = append(args, c.outer.outer.astFactory.CreateNumber(float64(c.nestedFinallyBlockCount)).Srcref(sourceNode))
	}

	c.WriteGeneratedNodeAndBreak(
		c.CallContextMethodResult(
			sourceNode,
			"leaveFinallyBlock",
			AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID),
			args...))
}

// SwitchCaseTo changes the currentCase to a new one.
//
// Java: TranspilationContext.switchCaseTo(Case)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) SwitchCaseTo(caseSection *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	c.currentCase.WillFollowBy(caseSection)
	c.allCases = append(c.allCases, caseSection)
	c.currentCase = caseSection
}

// PushLabels adds a named labels to the context.
//
// Java: TranspilationContext.pushLabels(ArrayList<Node>, Case, Case)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) PushLabels(labelNames []rhino.Node, breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	for _, labelName := range labelNames {
		jutil.CheckState(labelName.IsLabelName())
		c.namedLabels.Put(labelName.GetString(), NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases(breakCase, continueCase))
	}
}

// PopLabels removes the named labels from the context.
//
// Java: TranspilationContext.popLabels(ArrayList<Node>)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) PopLabels(labelNames []rhino.Node) {
	for _, labelName := range labelNames {
		jutil.CheckState(labelName.IsLabelName())
		c.namedLabels.Remove(labelName.GetString())
	}
}

// PushBreakContext adds "break" jump point to the context
//
// Java: TranspilationContext.pushBreakContext(Case)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) PushBreakContext(breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	c.breakCases = append(c.breakCases, breakCase)
}

// PushBreakContinueContext adds "break" and "continue" jump points to the
// context
//
// Java: TranspilationContext.pushBreakContinueContext(Case, Case)
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) PushBreakContinueContext(breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	c.PushBreakContext(breakCase)
	c.continueCases = append(c.continueCases, continueCase)
}

// PopBreakContext removes "break" jump point from the context, restoring the
// previous one
//
// Java: TranspilationContext.popBreakContext()
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) PopBreakContext() {
	if len(c.breakCases) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
	}
	c.breakCases = c.breakCases[:len(c.breakCases)-1]
}

// PopBreakContinueContext removes "break" and "continue" jump points from the
// context, restoring the previous ones.
//
// Java: TranspilationContext.popBreakContinueContext()
func (c *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) PopBreakContinueContext() {
	c.PopBreakContext()
	if len(c.continueCases) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
	}
	c.continueCases = c.continueCases[:len(c.continueCases)-1]
}

// Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case
// is a case section in a switch block of generator program.
type Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case struct {
	outer *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext

	id        int
	caseBlock rhino.Node

	// Records number of times the section was referenced.
	//
	// It's used to drop unreferenced sections.
	references []rhino.Node

	// Indicates that this case is a simple jump or a fall-though case. Points to the target
	// case.
	jumpTo *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case

	// Indicates that the body of this case could potentially be embedded into another block
	// node.
	//
	// Usually "if (a) {b();} else { c(); }" is transpiled into:
	//
	//	  if (a) { goto labelIf; }
	//	  c();
	//	  goto labelEnd;
	//	labelIf:
	//	  b();
	//	labelEnd:
	//
	// but "labelIf: b();" can be inlined to get shorter output:
	//
	//	  if (a) { b(); goto labelEnd; }
	//	  c();
	//	labelEnd:
	//
	// In this example "labelIf" case can be embedded into "{ goto labelIf; }" block.
	embedInto rhino.Node

	// Tells whether this case might fall-through.
	mayFallThrough bool
}

// NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case
// creates a new empty case section and assings a new id.
//
// Java: Case()
func NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case(outer *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case {
	cs := &Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case{outer: outer}
	cs.references = nil
	cs.mayFallThrough = true
	cs.id = outer.caseIdCounter
	outer.caseIdCounter++
	cs.caseBlock = outer.outer.outer.astFactory.CreateBlock().Srcref(outer.outer.originalGeneratorBody)
	return cs
}

// Java: Case.createCaseNode()
func (cs *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) CreateCaseNode() rhino.Node {
	return rhino.IR_CaseNode(cs.outer.outer.outer.astFactory.CreateNumber(float64(cs.id)).Srcref(cs.caseBlock), cs.caseBlock).
		Srcref(cs.caseBlock)
}

// GetNumber returns the number node of the case section and increments a
// reference counter.
//
// Java: Case.getNumber(Node)
func (cs *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) GetNumber(sourceNode rhino.Node) rhino.Node {
	if cs.jumpTo != nil {
		return cs.jumpTo.GetNumber(sourceNode)
	}
	node := cs.outer.outer.outer.astFactory.CreateNumber(float64(cs.id)).Srcref(sourceNode)
	cs.references = append(cs.references, node)
	return node
}

// JumpTo finalizes the case section with a jump instruction.
//
// AddNode cannot be invoked after this method is called.
//
// Java: Case.jumpTo(Case, Node)
func (cs *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) JumpTo(other *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, jumpBlock rhino.Node) {
	jutil.CheckState(jumpBlock.IsBlock())
	jutil.CheckState(cs.jumpTo == nil)
	cs.WillFollowBy(other)
	cs.caseBlock.AddChildrenToBack(jumpBlock.RemoveChildren())
	cs.mayFallThrough = false
}

// WillFollowBy informs which other case will be executed after this one.
//
// It's used to detect and then eliminate case statements that are used as
// simple jump hops:
//
//	case 100:
//	  $context.jumpTo(200);
//	  break;
//
// or
//
//	case 300:
//
// Java: Case.willFollowBy(Case)
func (cs *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) WillFollowBy(other *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) {
	if cs.jumpTo == nil && !cs.caseBlock.HasChildren() {
		jutil.CheckState(other.jumpTo == nil)
		cs.jumpTo = other
	}
}

// AddNode adds a new node to the end of the case block.
//
// Java: Case.addNode(Node)
func (cs *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) AddNode(n rhino.Node) {
	jutil.CheckState(cs.jumpTo == nil)
	jutil.CheckState(rhino.IR_MayBeStatement(n))
	cs.caseBlock.AddChildToBack(n)
}

// Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler
// adjusts YIELD-free nodes to run correctly inside a state machine program.
//
// The following transformations are performed:
//
//   - moving `var` into hoist scope;
//   - transpiling `return` statements;
//   - transpiling `break` and `continue` statements;
//   - transpiling references to `this` and `arguments`.
type Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler struct {
	outer *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext

	// Count the number of enclosing statements that a bare break could address.
	// A value > 0 means that a bare break statement we encounter can be left unmodified,
	// since it addresses a statement within the node we are transpiling.
	breakSuppressors int
	// Same as breakSuppressors, but for bare continue statements.
	continueSuppressors int
}

// Java: UnmarkedNodeTranspiler()
func NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler(outer *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext) *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler {
	return &Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler{outer: outer}
}

// Java: UnmarkedNodeTranspiler.shouldTraverse(NodeTraversal, Node, Node)
func (u *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsGeneratorSafe() {
		// Skip nodes that were generated by the compiler.
		n.SetGeneratorSafe(false)
		return false
	}

	jutil.CheckState(!n.IsGeneratorMarker())
	jutil.CheckStateMsg(!n.IsSuper(), "Reference to SUPER is not supported")

	if NodeUtil_IsLoopStructure(n) {
		u.continueSuppressors++
		u.breakSuppressors++
	} else if n.IsSwitch() {
		u.breakSuppressors++
	}

	if n.IsBreak() || n.IsContinue() {
		if n.HasChildren() {
			u.VisitNamedBreakContinue(n)
		} else {
			u.VisitBreakContinue(n)
		}
		return false
	}

	return !n.IsFunction()
}

// Java: UnmarkedNodeTranspiler.visit(NodeTraversal, Node, Node)
func (u *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if NodeUtil_IsLoopStructure(n) {
		u.continueSuppressors--
		u.breakSuppressors--
	} else if n.IsSwitch() {
		u.breakSuppressors--
	} else if n.IsThis() {
		u.VisitThis(n)
	} else if n.IsReturn() {
		u.VisitReturn(n)
	} else if n.IsName() && n.GetString() == "arguments" {
		u.VisitArguments(n)
	} else if n.IsVar() {
		if parent.IsVanillaFor() {
			u.VisitVanillaForLoopVar(n)
		} else if parent.IsForIn() {
			u.VisitForInLoopVar(n)
		} else {
			// NOTE: for-of loops are transpiled away before this pass
			u.VisitVarStatement(n)
		}
	} // else no changes need to be made
}

// VisitReturn adjusts return statements.
//
// Java: UnmarkedNodeTranspiler.visitReturn(Node)
func (u *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) VisitReturn(n rhino.Node) {
	// return ...;   =>   return $context.return(...);
	n.AddChildToFront(u.outer.ReturnExpression(n, n.RemoveFirstChild()))
}

// VisitNamedBreakContinue converts labeled `break` or `continue` statement
// into a jump.
//
// Java: UnmarkedNodeTranspiler.visitNamedBreakContinue(Node)
func (u *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) VisitNamedBreakContinue(n rhino.Node) {
	jutil.CheckState(n.GetFirstChild().IsLabelName())
	cases := u.outer.namedLabels.Get(n.GetFirstChild().GetString())
	if cases != nil {
		var caseSection *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case
		if n.IsBreak() {
			caseSection = cases.breakCase
		} else {
			caseSection = cases.continueCase
		}
		u.outer.outer.context.ReplaceBreakContinueWithJump(n, caseSection, u.breakSuppressors)
	}
}

// VisitBreakContinue converts `break` or `continue` statement into a jump.
//
// Java: UnmarkedNodeTranspiler.visitBreakContinue(Node)
func (u *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) VisitBreakContinue(n rhino.Node) {
	var caseSection *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case = nil
	if n.IsBreak() && u.breakSuppressors == 0 {
		if len(u.outer.breakCases) == 0 {
			panic(jutil.NewNoSuchElementException()) // ArrayDeque.getFirst()
		}
		caseSection = u.outer.breakCases[len(u.outer.breakCases)-1]
	}
	if n.IsContinue() && u.continueSuppressors == 0 {
		if len(u.outer.continueCases) == 0 {
			panic(jutil.NewNoSuchElementException()) // ArrayDeque.getFirst()
		}
		caseSection = u.outer.continueCases[len(u.outer.continueCases)-1]
	}
	if caseSection != nil {
		u.outer.outer.context.ReplaceBreakContinueWithJump(n, caseSection, u.breakSuppressors)
	}
}

// VisitThis replaces reference to `this` with `$jscomp$generator$this`.
//
// Java: UnmarkedNodeTranspiler.visitThis(Node)
func (u *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) VisitThis(n rhino.Node) {
	s := u.outer.outer
	newThis := s.outer.astFactory.CreateName(s.context.GetScopedName(Es6RewriteGenerators_GENERATOR_THIS), AstFactory_Type_(n))
	n.ReplaceWith(newThis)
	if !u.outer.thisReferenceFound {
		var_ := rhino.IR_VarWithValue(newThis.CloneNode().Srcref(n), n).Srcref(s.newGeneratorHoistBlock)
		s.HoistVarNode(var_)
		u.outer.thisReferenceFound = true
	}
}

// VisitArguments replaces reference to `arguments` with
// `$jscomp$generator$arguments`.
//
// Java: UnmarkedNodeTranspiler.visitArguments(Node)
func (u *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) VisitArguments(n rhino.Node) {
	s := u.outer.outer
	newArguments :=
		s.outer.astFactory.CreateName(s.context.GetScopedName(Es6RewriteGenerators_GENERATOR_ARGUMENTS), AstFactory_Type_(n)).Srcref(n)
	n.ReplaceWith(newArguments)
	if !u.outer.argumentsReferenceFound {
		var_ := rhino.IR_VarWithValue(newArguments.CloneNode(), n).Srcref(s.newGeneratorHoistBlock)
		s.HoistVarNode(var_)
		u.outer.argumentsReferenceFound = true
	}
}

// VisitVarStatement hoists `var` statements into the closure containing the
// generator to preserve their state across multiple invocation of state
// machine program.
//
//	var a = "test", b = i + 5;
//
// is transpiled to:
//
//	var a, b;
//	a = "test", b = i + 5;
//
// Java: UnmarkedNodeTranspiler.visitVarStatement(Node)
func (u *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) VisitVarStatement(varStatement rhino.Node) {
	commaExpression := u.ExtractAssignmentsToCommaExpression(varStatement)
	if commaExpression.IsNil() {
		varStatement.Detach()
	} else {
		varStatement.ReplaceWith(u.outer.outer.outer.astFactory.ExprResult(commaExpression))
	}
	// Move declaration without initial values to just before the program method definition.
	u.outer.outer.HoistVarNode(varStatement)
}

// VisitVanillaForLoopVar hoists `var` declarations in vanilla for loops into
// the closure containing the generator to preserve their state across
// multiple invocation of state machine program.
//
//	for (var a = "test", b = i + 5; ... ; )
//
// is transpiled to:
//
//	var a, b;
//	for (a = "test", b = i + 5; ...; )
//
// Java: UnmarkedNodeTranspiler.visitVanillaForLoopVar(Node)
func (u *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) VisitVanillaForLoopVar(varDeclaration rhino.Node) {
	commaExpression := u.ExtractAssignmentsToCommaExpression(varDeclaration)
	if commaExpression.IsNil() {
		// `for (var x; ` becomes `for (; `
		varDeclaration.ReplaceWith(rhino.IR_Empty())
	} else {
		// `for (var i = 0, j = 0; `... becomes `for (i = 0, j = 0; `...
		varDeclaration.ReplaceWith(commaExpression)
	}
	// Move declaration without initial values to just before the program method definition.
	u.outer.outer.HoistVarNode(varDeclaration)
}

// VisitForInLoopVar hoists `var` declarations in for-in loops into the
// closure containing the generator to preserve their state across multiple
// invocation of state machine program.
//
//	for (var a in obj)
//
// is transpiled to:
//
//	var a;
//	for (a in obj))
//
// Java: UnmarkedNodeTranspiler.visitForInLoopVar(Node)
func (u *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) VisitForInLoopVar(varDeclaration rhino.Node) {
	// `for (var varName in ` ...
	varName := varDeclaration.GetOnlyChild()
	jutil.CheckStateMsg1(!varName.HasChildren(), "%s", varName)
	clonedVarName := varName.CloneNode().SetJSDocInfo(nil)
	// becomes `for (varName in ` ...
	varDeclaration.ReplaceWith(clonedVarName)
	// Move declaration without initial values to just before the program method definition.
	u.outer.outer.HoistVarNode(varDeclaration)
}

// ExtractAssignmentsToCommaExpression removes all initializers from a var
// declaration and returns them as a single expression of comma-separated
// assignments or null if there aren't any initializers.
//
// varDeclaration is the VAR node; returns null or an expression node (e.g.
// `varName1 = 1, varName2 = y`).
//
// Java: UnmarkedNodeTranspiler.extractAssignmentsToCommaExpression(Node)
func (u *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_UnmarkedNodeTranspiler) ExtractAssignmentsToCommaExpression(varDeclaration rhino.Node) rhino.Node {
	astFactory := u.outer.outer.outer.astFactory
	assignments := []rhino.Node{}
	for varName := varDeclaration.GetFirstChild(); !varName.IsNil(); varName = varName.GetNext() {
		if varName.HasChildren() {
			copiedVarName := varName.CloneNode().SetJSDocInfo(nil)
			assign :=
				astFactory.
					CreateAssign(copiedVarName, varName.RemoveFirstChild()).
					Srcref(varName)
			assignments = append(assignments, assign)
		}
	}
	var commaExpression rhino.Node = rhino.NilNode
	for _, assignment := range assignments {
		if commaExpression.IsNil() {
			commaExpression = assignment
		} else {
			commaExpression = astFactory.CreateComma(commaExpression, assignment).Srcref(assignment)
		}
	}
	return commaExpression
}

// Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase
// reprasents a catch case that is used by try/catch transpilation
type Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase struct {
	catchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case

	// Number of finally blocks that should be executed before exception can be handled by this
	// catch case.
	finallyBlocks int
}

// Java: CatchCase(Case)
func NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase(catchCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase {
	cc := &Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_CatchCase{}
	cc.catchCase = catchCase
	return cc
}

// Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases
// stores "break" and "continue" case sections assosiated with a label.
type Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases struct {
	breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case

	continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case
}

// Java: LabelCases(Case, Case)
func NewEs6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases(breakCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case, continueCase *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_Case) *Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases {
	lc := &Es6RewriteGenerators_SingleGeneratorFunctionTranspiler_TranspilationContext_LabelCases{}
	lc.breakCase = breakCase
	lc.continueCase = continueCase
	return lc
}

// Es6RewriteGenerators_YieldNodeMarker marks "yield" nodes and propagates
// this information up through the tree
type Es6RewriteGenerators_YieldNodeMarker struct{}

// Java: YieldNodeMarker.shouldTraverse(NodeTraversal, Node, Node)
func (y *Es6RewriteGenerators_YieldNodeMarker) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return !n.IsFunction()
}

// Java: YieldNodeMarker.visit(NodeTraversal, Node, Node)
func (y *Es6RewriteGenerators_YieldNodeMarker) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsYield() {
		n.SetGeneratorMarker(true)
	}

	// This class is used on a tree that is detached from the main AST, so this will not end up
	// marking the parent of the node used to start the traversal.
	if !parent.IsNil() && n.IsGeneratorMarker() {
		parent.SetGeneratorMarker(true)
	}
}
