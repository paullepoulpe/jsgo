// Port of com/google/javascript/jscomp/FunctionInjector.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: the call sites and functions this file reads (the callee, the
// this-argument and the arguments of a call, the assignment or var around
// it, the name, params and body of a function) go through the typed view
// (internal/ast); the subtree walks stay on rhino.Node.

// FunctionInjector is a set of utility functions that replaces CALL with a
// specified FUNCTION body, replacing and aliasing function parameters as
// necessary.
type FunctionInjector struct {
	compiler               AbstractCompiler
	allowDecomposition     bool
	knownConstantFunctions map[string]bool
	assumeStrictThis       bool
	safeNameIdSupplier     func() string
	throwawayNameSupplier  func() string

	functionArgumentInjector *FunctionArgumentInjector
	mutator                  FunctionToBlockMutator
	paramReplacements        inlineIDNodeMap
	decomposer               *ExpressionDecomposer

	// Combined cache of function node to whether it deeply contains an `eval`
	// call (evalState: 0 = unknown, 1 = false, 2 = true) and any inner function.
	fnInfoCache map[rhino.Node]functionInjectorFnInfo
}

type functionInjectorFnInfo struct {
	inner     rhino.Node
	evalState uint8
}

// Sentinel value indicating that the key contains no functions.
var FunctionInjector_NO_FUNCTIONS = rhino.NewNode(rhino.Token_FUNCTION)

// Sentinel value indicating that the key contains multiple distinct
// functions.
var FunctionInjector_MULTIPLE_FUNCTIONS = rhino.NewNode(rhino.Token_FUNCTION)

// Java: FunctionInjector(Builder)
func NewFunctionInjector(builder *FunctionInjector_Builder) *FunctionInjector {
	f := &FunctionInjector{}
	// The anonymous Supplier<String> throwawayNameSupplier, with its nextId
	// field.
	nextId := 0
	f.throwawayNameSupplier = func() string {
		id := nextId
		nextId++
		return strconv.Itoa(id)
	}
	f.compiler = jutil.CheckNotNull(builder.compiler)
	f.safeNameIdSupplier = jutil.CheckNotNull(builder.safeNameIdSupplier)
	f.assumeStrictThis = builder.assumeStrictThis
	f.allowDecomposition = builder.allowDecomposition
	f.functionArgumentInjector = jutil.CheckNotNull(builder.functionArgumentInjector)
	f.mutator = FunctionToBlockMutator{
		compiler:                 f.compiler,
		safeNameIdSupplier:       f.safeNameIdSupplier,
		functionArgumentInjector: f.functionArgumentInjector,
	}
	return f
}

// FunctionInjector_Builder is FunctionInjector.Builder.
type FunctionInjector_Builder struct {
	compiler                 AbstractCompiler
	safeNameIdSupplier       func() string
	assumeStrictThis         bool
	allowDecomposition       bool
	functionArgumentInjector *FunctionArgumentInjector
}

// Java: FunctionInjector.Builder(AbstractCompiler)
func NewFunctionInjector_Builder(compiler AbstractCompiler) *FunctionInjector_Builder {
	b := &FunctionInjector_Builder{
		safeNameIdSupplier:       nil,
		assumeStrictThis:         true,
		allowDecomposition:       true,
		functionArgumentInjector: nil,
	}
	b.compiler = jutil.CheckNotNull(compiler)
	return b
}

// SafeNameIdSupplier provides the name supplier to use for injection.
//
// If this method is not called, `compiler.getUniqueNameIdSupplier()` will be
// used.
//
// Java: FunctionInjector.Builder.safeNameIdSupplier(Supplier<String>)
func (b *FunctionInjector_Builder) SafeNameIdSupplier(safeNameIdSupplier func() string) *FunctionInjector_Builder {
	b.safeNameIdSupplier = jutil.CheckNotNull(safeNameIdSupplier)
	return b
}

// AllowDecomposition allows decomposition of expressions.
//
// Default is `true`.
//
// Java: FunctionInjector.Builder.allowDecomposition(boolean)
func (b *FunctionInjector_Builder) AllowDecomposition(allowDecomposition bool) *FunctionInjector_Builder {
	b.allowDecomposition = allowDecomposition
	return b
}

// Java: FunctionInjector.Builder.assumeStrictThis(boolean)
func (b *FunctionInjector_Builder) AssumeStrictThis(assumeStrictThis bool) *FunctionInjector_Builder {
	b.assumeStrictThis = assumeStrictThis
	return b
}

// FunctionArgumentInjector specifies the `FunctionArgumentInjector` to be
// used.
//
// Default is for the builder to create this. This method exists for testing
// purposes.
//
// Java: FunctionInjector.Builder.functionArgumentInjector(FunctionArgumentInjector)
func (b *FunctionInjector_Builder) FunctionArgumentInjector(functionArgumentInjector *FunctionArgumentInjector) *FunctionInjector_Builder {
	b.functionArgumentInjector = jutil.CheckNotNull(functionArgumentInjector)
	return b
}

// Java: FunctionInjector.Builder.build()
func (b *FunctionInjector_Builder) Build() *FunctionInjector {
	if b.safeNameIdSupplier == nil {
		b.safeNameIdSupplier = b.compiler.GetUniqueNameIdSupplier()
	}
	if b.functionArgumentInjector == nil {
		b.functionArgumentInjector = NewFunctionArgumentInjector(jutil.CheckNotNull(b.compiler.GetAstAnalyzer()))
	}
	return NewFunctionInjector(b)
}

// FunctionInjector_InliningMode is the type of inlining to perform.
type FunctionInjector_InliningMode int32

const (
	// Directly replace the call expression. Only functions of meeting strict
	// preconditions can be inlined.
	FunctionInjector_InliningMode_DIRECT FunctionInjector_InliningMode = iota

	// Replaces the call expression with a block of statements. Conditions on
	// the function are looser in mode, but stricter on the call site.
	FunctionInjector_InliningMode_BLOCK
)

var functionInjector_InliningMode_names = [...]string{"DIRECT", "BLOCK"}

// Java: FunctionInjector.InliningMode.name()
func (m FunctionInjector_InliningMode) String() string { return functionInjector_InliningMode_names[m] }

// Java: FunctionInjector.InliningMode.ordinal()
func (m FunctionInjector_InliningMode) Ordinal() int { return int(m) }

// FunctionInjector_Reference holds a reference to the call node of a
// function call.
type FunctionInjector_Reference struct {
	callNode              rhino.Node
	scope                 *Scope
	t                     *NodeTraversal
	chunk                 *JSChunk
	mode                  FunctionInjector_InliningMode
	requiresDecomposition bool
	inlined               bool
}

// NewFunctionInjector_Reference constructs a FunctionInjector_Reference.
//
// Java: FunctionInjector.Reference(Node, Scope, JSChunk, InliningMode)
func NewFunctionInjector_Reference(callNode rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) *FunctionInjector_Reference {
	r := &FunctionInjector_Reference{}
	r.init(callNode, scope, chunk, mode, false)
	return r
}

func (r *FunctionInjector_Reference) init(callNode rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode, decompose bool) {
	r.callNode = callNode
	r.scope = scope
	r.t = nil
	r.chunk = chunk
	r.mode = mode
	r.requiresDecomposition = decompose
	r.inlined = false
}

// SetRequiresDecomposition records whether the call site must be moved or
// decomposed before block inlining.
func (r *FunctionInjector_Reference) SetRequiresDecomposition(val bool) {
	r.requiresDecomposition = val
}

func (r *FunctionInjector_Reference) inGlobalScope() bool {
	if r.scope != nil {
		return r.scope.IsGlobal()
	}
	if r.t != nil {
		return r.t.InGlobalScope()
	}
	return false
}

func (r *FunctionInjector_Reference) closestHoistScopeRoot() rhino.Node {
	if r.scope != nil {
		hs := r.scope.GetClosestHoistScope()
		if hs == nil || hs.IsGlobal() {
			return rhino.NilNode
		}
		return hs.GetRootNode()
	}
	if r.t != nil {
		root := r.t.GetClosestHoistScopeRoot()
		if root.IsNil() || root.IsRoot() || root.IsScript() {
			return rhino.NilNode
		}
		return root
	}
	return rhino.NilNode
}

func (r *FunctionInjector_Reference) lookupScope() *Scope {
	if r.scope != nil {
		return r.scope
	}
	if r.t != nil {
		return r.t.resolveLookupScope(len(r.t.scopes) - 1)
	}
	return nil
}

// Java: FunctionInjector.Reference.toString()
func (r *FunctionInjector_Reference) String() string {
	return "Reference @ " + jutil.ToJavaString(r.callNode)
}

// Inline inlines a function into the call site.
//
// Java: FunctionInjector.inline(Reference, String, Node)
func (f *FunctionInjector) Inline(ref *FunctionInjector_Reference, fnName string, fnNode rhino.Node) rhino.Node {
	return f.inlineWithSummary(ref, fnName, fnNode, nil)
}

func (f *FunctionInjector) inlineWithSummary(
	ref *FunctionInjector_Reference, fnName string, fnNode rhino.Node, summary *calleeSummary) rhino.Node {
	jutil.CheckState(f.compiler.GetLifeCycleStage().IsNormalized())
	return f.internalInlineWithSummary(ref, fnName, fnNode, summary)
}

// UnsafeInline inlines a function into the call site. Note that this unsafe
// version doesn't verify if the AST is normalized. You should use inline
// instead, unless you are 100% certain that the bit of code you're inlining
// is safe without being normalized first.
//
// Java: FunctionInjector.unsafeInline(Reference, String, Node)
func (f *FunctionInjector) UnsafeInline(ref *FunctionInjector_Reference, fnName string, fnNode rhino.Node) rhino.Node {
	return f.InternalInline(ref, fnName, fnNode)
}

// Java: FunctionInjector.internalInline(Reference, String, Node)
func (f *FunctionInjector) InternalInline(ref *FunctionInjector_Reference, fnName string, fnNode rhino.Node) rhino.Node {
	return f.internalInlineWithSummary(ref, fnName, fnNode, nil)
}

func (f *FunctionInjector) internalInlineWithSummary(
	ref *FunctionInjector_Reference, fnName string, fnNode rhino.Node, summary *calleeSummary) rhino.Node {
	var result rhino.Node
	if ref.mode == FunctionInjector_InliningMode_DIRECT {
		result = f.InlineReturnValue(ref, fnNode)
	} else {
		result = f.inlineFunctionWithSummary(ref, fnNode, fnName, summary)
	}
	f.compiler.ReportChangeToEnclosingScope(result)
	return result
}

// InlineReturnValue inlines a function that fulfills the requirements of
// canInlineReferenceDirectly into the call site, replacing only the CALL
// node.
//
// Java: FunctionInjector.inlineReturnValue(Reference, Node)
func (f *FunctionInjector) InlineReturnValue(ref *FunctionInjector_Reference, fnNode rhino.Node) rhino.Node {
	callNode := ref.callNode
	body := ast.MustFunction(fnNode).Body()

	// NOTE: As the normalize pass guarantees globals aren't being
	// shadowed and an expression can't introduce new names, there is
	// no need to check for conflicts.

	// Create an paramName -> argument value map, checking for side effects.
	paramReplacements := &f.paramReplacements
	f.functionArgumentInjector.fillCallParamNodeMap(
		fnNode, callNode, f.safeNameIdSupplier, paramReplacements)

	var newExpression rhino.Node
	if block, ok := body.AsBlock(); ok && !block.HasStatements() {
		newExpression = NodeUtil_NewUndefinedNode(block.Node())
	} else if !body.Node().HasChildren() {
		newExpression = NodeUtil_NewUndefinedNode(body.Node())
	} else {
		firstStmt := body.MustBlock().MustFirstStatement()
		returnNode, ok := firstStmt.AsReturn()
		jutil.CheckArgumentMsg1(ok, "%s", firstStmt.Node())

		// Clone the return node first.
		safeReturnNode := ast.MustReturn(returnNode.Node().CloneTree())
		inlineResult := f.functionArgumentInjector.injectByID(nil, safeReturnNode.Node(), rhino.NilNode, paramReplacements, true)
		jutil.CheckArgument(safeReturnNode.Is(inlineResult))
		value, ok := safeReturnNode.Value()
		jutil.CheckState(ok) // canInlineReferenceDirectly: the return has a value
		newExpression = value.Detach().Node()
		NodeUtil_MarkNewScopesChanged(newExpression, f.compiler)
	}

	// If the call site had a cast ensure it's persisted to the new expression that replaces it.
	typeBeforeCast := callNode.GetJSTypeBeforeCast()
	if typeBeforeCast != nil {
		newExpression.SetJSTypeBeforeCast(typeBeforeCast)
		newExpression.SetJSType(callNode.GetJSType())
	}
	// If the new expression has no color or the UNKNOWN color, attach the color the call node. It
	// may be more accurate if the call node was in a cast (we don't track information about casts,
	// though)
	if callNode.GetColor() != nil && callNode.IsColorFromTypeCast() {
		newExpression.SetColor(callNode.GetColor())
		newExpression.SetColorFromTypeCast()
	}
	callNode.ReplaceWith(newExpression)
	NodeUtil_MarkFunctionsDeleted(callNode, f.compiler)
	return newExpression
}

// FunctionInjector_CallSiteType is the enum of supported call site types.
type FunctionInjector_CallSiteType int32

const (
	// Used for a call site for which there does not exist a method to inline
	// it.
	FunctionInjector_CallSiteType_UNSUPPORTED FunctionInjector_CallSiteType = iota

	// A call as a statement. For example: "foo();". EXPR_RESULT CALL
	FunctionInjector_CallSiteType_SIMPLE_CALL

	// An assignment, where the result of the call is assigned to a simple
	// name. For example: "a = foo();". EXPR_RESULT NAME A CALL FOO
	FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT

	// An var declaration and initialization, where the result of the call is
	// assigned to the declared name name. For example: "var a = foo();". VAR
	// NAME A CALL FOO
	FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT

	// An arbitrary expression, the root of which is a EXPR_RESULT, IF,
	// RETURN, SWITCH or VAR. The call must be the first side-effect in the
	// expression.
	//
	// Examples include: "if (foo()) {..." "return foo();" "var a = 1 +
	// foo();" "a = 1 + foo()" "foo() ? 1:0" "foo() && x"
	FunctionInjector_CallSiteType_EXPRESSION

	// An arbitrary expression, the root of which is a EXPR_RESULT, IF,
	// RETURN, SWITCH or VAR. Where the call is not the first side-effect in
	// the expression.
	FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION
)

var functionInjector_CallSiteType_names = [...]string{
	"UNSUPPORTED",
	"SIMPLE_CALL",
	"SIMPLE_ASSIGNMENT",
	"VAR_DECL_SIMPLE_ASSIGNMENT",
	"EXPRESSION",
	"DECOMPOSABLE_EXPRESSION",
}

// Java: FunctionInjector.CallSiteType.name()
func (c FunctionInjector_CallSiteType) String() string { return functionInjector_CallSiteType_names[c] }

// Java: FunctionInjector.CallSiteType.ordinal()
func (c FunctionInjector_CallSiteType) Ordinal() int { return int(c) }

// Prepare is the abstract enum method, with each constant's body.
//
// Java: FunctionInjector.CallSiteType.prepare(FunctionInjector, Reference)
func (c FunctionInjector_CallSiteType) Prepare(injector *FunctionInjector, ref *FunctionInjector_Reference) {
	switch c {
	case FunctionInjector_CallSiteType_UNSUPPORTED:
		panic(jutil.NewIllegalStateExceptionMsg("unexpected: " + ref.String()))
	case FunctionInjector_CallSiteType_SIMPLE_CALL:
		// Nothing to do.
	case FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT:
		// Nothing to do.
	case FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT:
		// Nothing to do.
	case FunctionInjector_CallSiteType_EXPRESSION:
		callNode := ref.callNode
		injector.GetDecomposer(ref.scope).MoveExpression(callNode)

		// Reclassify after move
		callSiteType := injector.ClassifyCallSite(ref)
		jutil.CheckState(c != callSiteType)
		callSiteType.Prepare(injector, ref)
	case FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION:
		callNode := ref.callNode
		injector.GetDecomposer(ref.scope).MaybeExposeExpression(callNode)

		// Reclassify after decomposition
		callSiteType := injector.ClassifyCallSite(ref)
		jutil.CheckState(c != callSiteType)
		callSiteType.Prepare(injector, ref)
	}
}

// ClassifyCallSite determines which, if any, of the supported types the call
// site is.
//
// Constant vars are treated differently so that we don't break their
// const-ness when we decompose the expression. Once the CONSTANT_VAR
// annotation is used everywhere instead of coding conventions, we should
// just teach this pass how to remove the annotation.
//
// Java: FunctionInjector.classifyCallSite(Reference)
func (f *FunctionInjector) ClassifyCallSite(ref *FunctionInjector_Reference) FunctionInjector_CallSiteType {
	return f.classifyCallSite(ref, true)
}

func (f *FunctionInjector) classifyCallSite(ref *FunctionInjector_Reference, allowDecomposition bool) FunctionInjector_CallSiteType {
	callNode := ref.callNode

	// Verify the call site:
	if NodeUtil_IsExprCall(callNode.GetParent()) {
		// This is a simple call. Example: "foo();".
		return FunctionInjector_CallSiteType_SIMPLE_CALL
	} else if _, ok := functionInjector_simpleAssignTarget(callNode); ok {
		// This is a simple assignment.  Example: "x = foo();"
		return FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT
	} else if _, ok := functionInjector_simpleVarDecl(callNode); ok {
		// This is a var declaration.  Example: "var x = foo();"
		// TODO(johnlenz): Should we be checking for constants on the
		// left-hand-side of the assignments and handling them as EXPRESSION?
		return FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT
	} else if allowDecomposition {
		if ExpressionDecomposer_FindExpressionRoot(callNode).IsNil() {
			return FunctionInjector_CallSiteType_UNSUPPORTED
		}
		decomposer := f.GetDecomposer(ref.lookupScope())
		switch decomposer.CanExposeExpression(callNode) {
		case ExpressionDecomposer_DecompositionType_MOVABLE:
			return FunctionInjector_CallSiteType_EXPRESSION
		case ExpressionDecomposer_DecompositionType_DECOMPOSABLE:
			return FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION
		case ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE:
		}
	}

	return FunctionInjector_CallSiteType_UNSUPPORTED
}

// functionInjector_simpleAssignTarget returns the name that the statement
// x = call() assigns the call to; ok is false for any other call site.
func functionInjector_simpleAssignTarget(callNode rhino.Node) (ast.Name, bool) {
	if ast.RoleOf(callNode) != ast.RoleAssignValue {
		return ast.Name{}, false
	}
	assign := ast.MustAssign(callNode.GetParent())
	if ast.RoleOf(assign.Node()) != ast.RoleExprStmtExpr {
		return ast.Name{}, false
	}
	name, ok := assign.Target().AsName()
	// TODO(nicksantos): Remove this once everyone is using
	// the CONSTANT_VAR annotation. We know how to remove that.
	return name, ok && !NodeUtil_IsConstantName(name.Node())
}

// functionInjector_simpleVarDecl returns the name that var x = call()
// declares, the only one of its var; ok is false for any other call site.
func functionInjector_simpleVarDecl(callNode rhino.Node) (ast.Name, bool) {
	if ast.RoleOf(callNode) != ast.RoleNameInit {
		return ast.Name{}, false
	}
	name := ast.MustName(callNode.GetParent())
	// TODO(nicksantos): Remove this once everyone is using the CONSTANT_VAR annotation.
	if NodeUtil_IsConstantName(name.Node()) {
		return ast.Name{}, false
	}
	// Note: not let or const. See InlineFunctionsTest.testInlineFunctions35
	decl, ok := ast.AsNameDeclarationOf(name.Node().GetParent(), rhino.Token_VAR)
	return name, ok && decl.NumDeclarators() == 1
}

// Java: FunctionInjector.getDecomposer(Scope)
func (f *FunctionInjector) GetDecomposer(scope *Scope) *ExpressionDecomposer {
	if f.decomposer == nil {
		f.decomposer = f.compiler.CreateExpressionDecomposer(f.safeNameIdSupplier, f.knownConstantFunctions, scope)
	} else {
		f.decomposer.scope = scope
		f.decomposer.knownConstantFunctions = f.knownConstantFunctions
	}
	return f.decomposer
}

// MaybePrepareCall rewrites, if required, the statement containing the call
// expression.
//
// See ExpressionDecomposer#canExposeExpression.
//
// Java: FunctionInjector.maybePrepareCall(Reference)
func (f *FunctionInjector) MaybePrepareCall(ref *FunctionInjector_Reference) {
	callSiteType := f.ClassifyCallSite(ref)
	callSiteType.Prepare(f, ref)
}

// InlineFunction inlines a function which fulfills the requirements of
// canInlineReferenceAsStatementBlock into the call site, replacing the
// parent expression.
//
// Java: FunctionInjector.inlineFunction(Reference, Node, String)
func (f *FunctionInjector) InlineFunction(ref *FunctionInjector_Reference, fnNode rhino.Node, fnName string) rhino.Node {
	return f.inlineFunctionWithSummary(ref, fnNode, fnName, nil)
}

func (f *FunctionInjector) inlineFunctionWithSummary(
	ref *FunctionInjector_Reference, fnNode rhino.Node, fnName string, summary *calleeSummary) rhino.Node {
	callNode := ref.callNode
	parent := callNode.GetParent()
	grandParent := parent.GetParent()

	// TODO(johnlenz): Consider storing the callSite classification in the
	// reference object and passing it in here.
	callSiteType := f.ClassifyCallSite(ref)
	jutil.CheckArgument(callSiteType != FunctionInjector_CallSiteType_UNSUPPORTED)

	// Store the name for the result. This will be used to
	// replace "return expr" statements with "resultName = expr"
	// to replace
	var resultName string
	needsDefaultReturnResult := true
	switch callSiteType {
	case FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT:
		target, _ := functionInjector_simpleAssignTarget(callNode)
		resultName = target.Ident()
		removeConstantVarAnnotationByID(ref.scope, target.Node().GetStringID())
	case FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT:
		resultName = ast.MustName(parent).Ident()
		removeConstantVarAnnotationByID(ref.scope, parent.GetStringID())
	case FunctionInjector_CallSiteType_SIMPLE_CALL:
		resultName = "" // "foo()" doesn't need a result.
		needsDefaultReturnResult = false
	case FunctionInjector_CallSiteType_EXPRESSION:
		panic(jutil.NewIllegalStateExceptionMsg("Movable expressions must be moved before inlining."))
	case FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Decomposable expressions must be decomposed before inlining."))
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected call site type."))
	}

	isCallInLoop := NodeUtil_IsWithinLoop(callNode)
	f.mutator.allowMoveCallArgs = true
	newBlock := f.mutator.mutateWithSummary(
		fnName, fnNode, callNode, resultName, needsDefaultReturnResult, isCallInLoop, summary)
	f.mutator.allowMoveCallArgs = false
	NodeUtil_MarkNewScopesChanged(newBlock, f.compiler)

	// TODO(nicksantos): Create a common mutation function that
	// can replace either a VAR name assignment, assignment expression or
	// a EXPR_RESULT.
	switch callSiteType {
	case FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT:
		// Remove the call from the name node.
		name := ast.MustName(parent)
		value, _ := name.Init()
		name.RemoveInit()
		NodeUtil_MarkFunctionsDeleted(value.Node(), f.compiler)
		// Add the call, after the VAR.
		newBlock.InsertAfter(grandParent)
	case FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT:
		// The assignment is now part of the inline function so
		// replace it completely.
		jutil.CheckState(grandParent.IsExprResult())
		grandParent.ReplaceWith(newBlock)
		NodeUtil_MarkFunctionsDeleted(grandParent, f.compiler)
	case FunctionInjector_CallSiteType_SIMPLE_CALL:
		// If nothing is looking at the result just replace the call.
		jutil.CheckState(parent.IsExprResult())
		parent.ReplaceWith(newBlock)
		NodeUtil_MarkFunctionsDeleted(parent, f.compiler)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected call site type."))
	}

	return newBlock
}

func removeConstantVarAnnotationByID(scope *Scope, id uint32) {
	var_ := scope.GetVarByID(id)
	var nameNode rhino.Node
	if var_ != nil {
		nameNode = var_.GetNameNode()
	}
	if nameNode.IsNil() {
		return
	}

	if nameNode.IsDeclaredConstantVar() {
		nameNode.SetDeclaredConstantVar(false)
	}
}

// SetKnownConstantFunctions stores the names of known constants to be used
// when classifying call-sites in expressions.
//
// Java: FunctionInjector.setKnownConstantFunctions(ImmutableSet<String>)
func (f *FunctionInjector) SetKnownConstantFunctions(knownConstantFunctions map[string]bool) {
	// This is only expected to be set once. The same set should be used
	// when evaluating call-sites and inlining calls.
	jutil.CheckState(len(f.knownConstantFunctions) == 0)
	f.knownConstantFunctions = knownConstantFunctions
}
