// Port of com/google/javascript/jscomp/FunctionInjector.java @ closure-compiler v20261015

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// FunctionInjector is a set of utility functions that replaces CALL with a
// specified FUNCTION body, replacing and aliasing function parameters as
// necessary.
type FunctionInjector struct {
	compiler               AbstractCompiler
	allowDecomposition     bool
	knownConstantFunctions *jutil.LinkedHashSet[string]
	assumeStrictThis       bool
	assumeMinimumCapture   bool
	safeNameIdSupplier     func() string
	throwawayNameSupplier  func() string

	functionArgumentInjector *FunctionArgumentInjector
	mutator                  FunctionToBlockMutator
	paramReplacementsScratch *jutil.LinkedHashMap[string, *rhino.Node]
	decomposer               *ExpressionDecomposer

	// Combined cache of function node to whether it deeply contains an `eval`
	// call (evalState: 0 = unknown, 1 = false, 2 = true) and any inner function.
	fnInfoCache map[*rhino.Node]functionInjectorFnInfo
}

type functionInjectorFnInfo struct {
	inner     *rhino.Node
	evalState uint8
}

// Sentinel value indicating that the key contains no functions.
var FunctionInjector_NO_FUNCTIONS = rhino.NewNode(rhino.Token_FUNCTION)

// Sentinel value indicating that the key contains multiple distinct
// functions.
var FunctionInjector_MULTIPLE_FUNCTIONS = rhino.NewNode(rhino.Token_FUNCTION)

// Java: FunctionInjector(Builder)
func NewFunctionInjector(builder *FunctionInjector_Builder) *FunctionInjector {
	f := &FunctionInjector{
		knownConstantFunctions:   jutil.NewLinkedHashSet[string](),
		paramReplacementsScratch: jutil.NewLinkedHashMap[string, *rhino.Node](),
	}
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
	f.assumeMinimumCapture = builder.assumeMinimumCapture
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
	assumeMinimumCapture     bool
	allowDecomposition       bool
	functionArgumentInjector *FunctionArgumentInjector
}

// Java: FunctionInjector.Builder(AbstractCompiler)
func NewFunctionInjector_Builder(compiler AbstractCompiler) *FunctionInjector_Builder {
	b := &FunctionInjector_Builder{
		safeNameIdSupplier:       nil,
		assumeStrictThis:         true,
		assumeMinimumCapture:     true,
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

// Java: FunctionInjector.Builder.assumeMinimumCapture(boolean)
func (b *FunctionInjector_Builder) AssumeMinimumCapture(assumeMinimumCapture bool) *FunctionInjector_Builder {
	b.assumeMinimumCapture = assumeMinimumCapture
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

// Java: FunctionInjector.InliningMode.values()
func FunctionInjector_InliningMode_values() []FunctionInjector_InliningMode {
	return []FunctionInjector_InliningMode{FunctionInjector_InliningMode_DIRECT, FunctionInjector_InliningMode_BLOCK}
}

// Java: FunctionInjector.InliningMode.name()
func (m FunctionInjector_InliningMode) String() string { return functionInjector_InliningMode_names[m] }

// Java: FunctionInjector.InliningMode.ordinal()
func (m FunctionInjector_InliningMode) Ordinal() int { return int(m) }

// FunctionInjector_Reference holds a reference to the call node of a
// function call.
//
// It has a subclass (InlineFunctions.Reference), so it is an interface
// (PORTING.md §4); FunctionInjector_ReferenceBase holds its fields.
type FunctionInjector_Reference interface {
	// Java: FunctionInjector.Reference.toString()
	String() string

	GetFunctionInjector_ReferenceBase() *FunctionInjector_ReferenceBase
}

// FunctionInjector_ReferenceBase holds FunctionInjector.Reference's fields.
type FunctionInjector_ReferenceBase struct {
	callNode *rhino.Node
	scope    *Scope
	chunk    *JSChunk
	mode     FunctionInjector_InliningMode
}

// NewFunctionInjector_Reference constructs a FunctionInjector.Reference.
//
// Java: FunctionInjector.Reference(Node, Scope, JSChunk, InliningMode)
func NewFunctionInjector_Reference(callNode *rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) *FunctionInjector_ReferenceBase {
	r := &FunctionInjector_ReferenceBase{}
	r.InitFunctionInjector_ReferenceBase(callNode, scope, chunk, mode)
	return r
}

// InitFunctionInjector_ReferenceBase is the Reference(Node, Scope, JSChunk,
// InliningMode) constructor body, called by subclasses (super(...)).
//
// Java: FunctionInjector.Reference(Node, Scope, JSChunk, InliningMode)
func (r *FunctionInjector_ReferenceBase) InitFunctionInjector_ReferenceBase(callNode *rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) {
	r.callNode = callNode
	r.scope = scope
	r.chunk = chunk
	r.mode = mode
}

// GetFunctionInjector_ReferenceBase gives access to the base through the
// interface.
func (r *FunctionInjector_ReferenceBase) GetFunctionInjector_ReferenceBase() *FunctionInjector_ReferenceBase {
	return r
}

// Java: FunctionInjector.Reference.toString()
func (r *FunctionInjector_ReferenceBase) String() string {
	return "Reference @ " + jutil.ToJavaString(r.callNode)
}

// In order to estimate the cost of lining, we make the assumption that
// Identifiers are reduced 2 characters. For the call arguments, the
// important thing is that the cost is assumed to be the same in the call and
// the function, so the actual length doesn't matter in most cases.
const FunctionInjector_NAME_COST_ESTIMATE = InlineCostEstimator_ESTIMATED_IDENTIFIER_COST

// The cost of a argument separator (a comma).
const FunctionInjector_COMMA_COST = 1

// The cost of the parentheses needed to make a call.
const FunctionInjector_PAREN_COST = 2

func functionInjector_blocksMinimumRequirements(n *rhino.Node, fnName, fnRecursionName string, checkArguments bool) bool {
	switch n.GetToken() {
	case rhino.Token_NAME:
		s := n.GetString()
		if (checkArguments && scopeNameEq(s, "arguments")) ||
			scopeNameEq(s, "eval") ||
			(fnName != "" && scopeNameEq(s, fnName)) ||
			(fnRecursionName != "" && scopeNameEq(s, fnRecursionName)) {
			return true
		}
	case rhino.Token_SUPER, rhino.Token_NEW_TARGET:
		return true
	case rhino.Token_FUNCTION:
		if !n.IsArrowFunction() {
			checkArguments = false
		}
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		if c.GetFirstChild() == nil {
			switch c.GetToken() {
			case rhino.Token_NAME:
				s := c.GetString()
				if (checkArguments && scopeNameEq(s, "arguments")) ||
					scopeNameEq(s, "eval") ||
					(fnName != "" && scopeNameEq(s, fnName)) ||
					(fnRecursionName != "" && scopeNameEq(s, fnRecursionName)) {
					return true
				}
			case rhino.Token_SUPER, rhino.Token_NEW_TARGET:
				return true
			}
			continue
		}
		if functionInjector_blocksMinimumRequirements(c, fnName, fnRecursionName, checkArguments) {
			return true
		}
	}
	return false
}

// DoesFunctionMeetMinimumRequirements: fnName is the name of this function.
// This either the name of the variable to which the function is assigned or
// the name from the FUNCTION node. fnNode is the FUNCTION node of the
// function to inspect. Returns whether the function node meets the minimum
// requirements for inlining.
//
// Java: FunctionInjector.doesFunctionMeetMinimumRequirements(String, Node)
func (f *FunctionInjector) DoesFunctionMeetMinimumRequirements(fnName string, fnNode *rhino.Node) bool {
	block := NodeUtil_GetFunctionBody(fnNode)

	// Basic restrictions on functions that can be inlined:
	// 1) It contains a reference to itself.
	// 2) It uses its parameters indirectly using "arguments" (it isn't
	//    handled yet.
	// 3) It references "eval". Inline a function containing eval can have
	//    large performance implications.

	fnRecursionName := fnNode.GetFirstChild().GetString()
	// Java checkState(fnRecursionName != null): a Go string is never null.
	if scopeNameEq(fnRecursionName, fnName) {
		fnRecursionName = ""
	}

	return !functionInjector_blocksMinimumRequirements(block, fnName, fnRecursionName, true)
}

// CanInlineReferenceToFunction: fnNode is the function to evaluate for
// inlining; needAliases a set of function parameter names that can not be
// used without aliasing (returned by getUnsafeParameterNames());
// referencesThis whether fnNode contains references to its this object;
// containsClosures whether fnNode contains inner functions or non-static
// class field initializers. Returns whether the inlining can occur.
//
// Java: FunctionInjector.canInlineReferenceToFunction(Reference, Node, ImmutableSet<String>, boolean, boolean)
func (f *FunctionInjector) CanInlineReferenceToFunction(
	ref FunctionInjector_Reference,
	fnNode *rhino.Node,
	needAliases *jutil.LinkedHashSet[string],
	referencesThis bool,
	containsClosures bool) FunctionInjector_CanInlineResult {
	// TODO(johnlenz): This function takes too many parameter, without
	// context.  Modify the API to take a structure describing the function.

	// Allow direct function calls or "fn.call" style calls.
	callNode := ref.GetFunctionInjector_ReferenceBase().callNode
	if !f.IsSupportedCallType(callNode, fnNode) {
		return FunctionInjector_CanInlineResult_NO
	}

	if FunctionInjector_HasSpreadCallArgument(callNode) {
		return FunctionInjector_CanInlineResult_NO
	}

	// Limit where functions that contain functions can be inline.  Introducing
	// an inner function into another function can capture a variable and cause
	// a memory leak.  This isn't a problem in the global scope as those values
	// last until explicitly cleared.
	if containsClosures {
		if !f.assumeMinimumCapture && !ref.GetFunctionInjector_ReferenceBase().scope.IsGlobal() {
			// TODO(johnlenz): Allow inlining into any scope without local names or inner functions.
			return FunctionInjector_CanInlineResult_NO
		} else if NodeUtil_IsWithinLoop(callNode) {
			// An inner closure maybe relying on a local value holding a value for a
			// single iteration through a loop.
			return FunctionInjector_CanInlineResult_NO
		}
	}

	// TODO(johnlenz): Add support for 'apply'
	if referencesThis && !NodeUtil_IsFunctionObjectCall(callNode) {
		// TODO(johnlenz): Allow 'this' references to be replaced with a
		// global 'this' object.
		return FunctionInjector_CanInlineResult_NO
	}

	if ref.GetFunctionInjector_ReferenceBase().mode == FunctionInjector_InliningMode_DIRECT {
		return f.CanInlineReferenceDirectly(ref, fnNode, needAliases)
	} else {
		return f.CanInlineReferenceAsStatementBlock(ref, fnNode, needAliases)
	}
}

// IsSupportedCallType: only ".call" calls and direct calls to functions are
// supported.
//
// callNode is the call evaluate. Returns whether the call is of a type that
// is supported.
//
// Java: FunctionInjector.isSupportedCallType(Node, Node)
func (f *FunctionInjector) IsSupportedCallType(callNode *rhino.Node, fnNode *rhino.Node) bool {
	if !callNode.GetFirstChild().IsName() {
		if NodeUtil_IsFunctionObjectCall(callNode) {
			if !f.assumeStrictThis && !fnNode.IsArrowFunction() {
				thisValue := callNode.GetSecondChild()
				if thisValue == nil || !thisValue.IsThis() {
					return false
				}
			}
		} else if NodeUtil_IsFunctionObjectApply(callNode) {
			return false
		}
	}

	return true
}

// Java: FunctionInjector.hasSpreadCallArgument(Node)
func FunctionInjector_HasSpreadCallArgument(callNode *rhino.Node) bool {
	jutil.CheckArgumentMsg1(NodeUtil_IsNormalOrOptChainCall(callNode), "%s", callNode)
	for arg := callNode.GetSecondChild(); arg != nil; arg = arg.GetNext() {
		if arg.IsSpread() {
			return true
		}
	}
	return false
}

// Inline inlines a function into the call site.
//
// Java: FunctionInjector.inline(Reference, String, Node)
func (f *FunctionInjector) Inline(ref FunctionInjector_Reference, fnName string, fnNode *rhino.Node) *rhino.Node {
	jutil.CheckState(f.compiler.GetLifeCycleStage().IsNormalized())
	return f.InternalInline(ref, fnName, fnNode)
}

// UnsafeInline inlines a function into the call site. Note that this unsafe
// version doesn't verify if the AST is normalized. You should use inline
// instead, unless you are 100% certain that the bit of code you're inlining
// is safe without being normalized first.
//
// Java: FunctionInjector.unsafeInline(Reference, String, Node)
func (f *FunctionInjector) UnsafeInline(ref FunctionInjector_Reference, fnName string, fnNode *rhino.Node) *rhino.Node {
	return f.InternalInline(ref, fnName, fnNode)
}

// Java: FunctionInjector.internalInline(Reference, String, Node)
func (f *FunctionInjector) InternalInline(ref FunctionInjector_Reference, fnName string, fnNode *rhino.Node) *rhino.Node {
	var result *rhino.Node
	if ref.GetFunctionInjector_ReferenceBase().mode == FunctionInjector_InliningMode_DIRECT {
		result = f.InlineReturnValue(ref, fnNode)
	} else {
		result = f.InlineFunction(ref, fnNode, fnName)
	}
	f.compiler.ReportChangeToEnclosingScope(result)
	return result
}

// InlineReturnValue inlines a function that fulfills the requirements of
// canInlineReferenceDirectly into the call site, replacing only the CALL
// node.
//
// Java: FunctionInjector.inlineReturnValue(Reference, Node)
func (f *FunctionInjector) InlineReturnValue(ref FunctionInjector_Reference, fnNode *rhino.Node) *rhino.Node {
	callNode := ref.GetFunctionInjector_ReferenceBase().callNode
	block := fnNode.GetLastChild()

	// NOTE: As the normalize pass guarantees globals aren't being
	// shadowed and an expression can't introduce new names, there is
	// no need to check for conflicts.

	// Create an paramName -> argument value map, checking for side effects.
	paramToArgMap := f.functionArgumentInjector.GetFunctionCallParameterMap(
		fnNode, callNode, f.safeNameIdSupplier)
	paramReplacements := f.paramReplacementsScratch
	paramReplacements.Clear()
	for key, value := range paramToArgMap.All() {
		paramReplacements.Put(key, value.Arg())
	}

	var newExpression *rhino.Node
	if !block.HasChildren() {
		srcLocation := block
		newExpression = NodeUtil_NewUndefinedNode(srcLocation)
	} else {
		returnNode := block.GetFirstChild()
		jutil.CheckArgumentMsg1(returnNode.IsReturn(), "%s", returnNode)

		// Clone the return node first.
		safeReturnNode := returnNode.CloneTree()
		inlineResult := f.functionArgumentInjector.Inject(nil, safeReturnNode, nil, paramReplacements)
		jutil.CheckArgument(safeReturnNode == inlineResult)
		newExpression = safeReturnNode.RemoveFirstChild()
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

// Java: FunctionInjector.CallSiteType.values()
func FunctionInjector_CallSiteType_values() []FunctionInjector_CallSiteType {
	return []FunctionInjector_CallSiteType{
		FunctionInjector_CallSiteType_UNSUPPORTED,
		FunctionInjector_CallSiteType_SIMPLE_CALL,
		FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT,
		FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT,
		FunctionInjector_CallSiteType_EXPRESSION,
		FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION,
	}
}

// Java: FunctionInjector.CallSiteType.name()
func (c FunctionInjector_CallSiteType) String() string { return functionInjector_CallSiteType_names[c] }

// Java: FunctionInjector.CallSiteType.ordinal()
func (c FunctionInjector_CallSiteType) Ordinal() int { return int(c) }

// Prepare is the abstract enum method, with each constant's body.
//
// Java: FunctionInjector.CallSiteType.prepare(FunctionInjector, Reference)
func (c FunctionInjector_CallSiteType) Prepare(injector *FunctionInjector, ref FunctionInjector_Reference) {
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
		callNode := ref.GetFunctionInjector_ReferenceBase().callNode
		injector.GetDecomposer(ref.GetFunctionInjector_ReferenceBase().scope).MoveExpression(callNode)

		// Reclassify after move
		callSiteType := injector.ClassifyCallSite(ref)
		jutil.CheckState(c != callSiteType)
		callSiteType.Prepare(injector, ref)
	case FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION:
		callNode := ref.GetFunctionInjector_ReferenceBase().callNode
		injector.GetDecomposer(ref.GetFunctionInjector_ReferenceBase().scope).MaybeExposeExpression(callNode)

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
func (f *FunctionInjector) ClassifyCallSite(ref FunctionInjector_Reference) FunctionInjector_CallSiteType {
	return f.classifyCallSite(ref, true)
}

func (f *FunctionInjector) classifyCallSite(ref FunctionInjector_Reference, allowDecomposition bool) FunctionInjector_CallSiteType {
	callNode := ref.GetFunctionInjector_ReferenceBase().callNode
	parent := callNode.GetParent()
	grandParent := parent.GetParent()

	// Verify the call site:
	if NodeUtil_IsExprCall(parent) {
		// This is a simple call. Example: "foo();".
		return FunctionInjector_CallSiteType_SIMPLE_CALL
	} else if NodeUtil_IsExprAssign(grandParent) &&
		!NodeUtil_IsNameDeclOrSimpleAssignLhs(callNode, parent) &&
		parent.GetFirstChild().IsName() &&
		// TODO(nicksantos): Remove this once everyone is using
		// the CONSTANT_VAR annotation. We know how to remove that.
		!NodeUtil_IsConstantName(parent.GetFirstChild()) {
		// This is a simple assignment.  Example: "x = foo();"
		return FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT
	} else if parent.IsName() &&
		// TODO(nicksantos): Remove this once everyone is using the CONSTANT_VAR annotation.
		!NodeUtil_IsConstantName(parent) &&
		// Note: not let or const. See InlineFunctionsTest.testInlineFunctions35
		grandParent.IsVar() &&
		grandParent.HasOneChild() {
		// This is a var declaration.  Example: "var x = foo();"
		// TODO(johnlenz): Should we be checking for constants on the
		// left-hand-side of the assignments and handling them as EXPRESSION?
		return FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT
	} else if allowDecomposition {
		decomposer := f.GetDecomposer(ref.GetFunctionInjector_ReferenceBase().scope)
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
func (f *FunctionInjector) MaybePrepareCall(ref FunctionInjector_Reference) {
	callSiteType := f.ClassifyCallSite(ref)
	callSiteType.Prepare(f, ref)
}

// InlineFunction inlines a function which fulfills the requirements of
// canInlineReferenceAsStatementBlock into the call site, replacing the
// parent expression.
//
// Java: FunctionInjector.inlineFunction(Reference, Node, String)
func (f *FunctionInjector) InlineFunction(ref FunctionInjector_Reference, fnNode *rhino.Node, fnName string) *rhino.Node {
	callNode := ref.GetFunctionInjector_ReferenceBase().callNode
	parent := callNode.GetParent()
	grandParent := parent.GetParent()

	// TODO(johnlenz): Consider storing the callSite classification in the
	// reference object and passing it in here.
	callSiteType := f.ClassifyCallSite(ref)
	jutil.CheckArgument(callSiteType != FunctionInjector_CallSiteType_UNSUPPORTED)

	// Store the name for the result. This will be used to
	// replace "return expr" statements with "resultName = expr"
	// to replace
	var resultName *string
	var resultNameVal string
	needsDefaultReturnResult := true
	switch callSiteType {
	case FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT:
		resultNameVal = parent.GetFirstChild().GetString()
		resultName = &resultNameVal
		FunctionInjector_RemoveConstantVarAnnotation(ref.GetFunctionInjector_ReferenceBase().scope, resultNameVal)
	case FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT:
		resultNameVal = parent.GetString()
		resultName = &resultNameVal
		FunctionInjector_RemoveConstantVarAnnotation(ref.GetFunctionInjector_ReferenceBase().scope, resultNameVal)
	case FunctionInjector_CallSiteType_SIMPLE_CALL:
		resultName = nil // "foo()" doesn't need a result.
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
	newBlock := f.mutator.Mutate(
		fnName, fnNode, callNode, resultName, needsDefaultReturnResult, isCallInLoop)
	f.mutator.allowMoveCallArgs = false
	NodeUtil_MarkNewScopesChanged(newBlock, f.compiler)

	// TODO(nicksantos): Create a common mutation function that
	// can replace either a VAR name assignment, assignment expression or
	// a EXPR_RESULT.
	switch callSiteType {
	case FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT:
		// Remove the call from the name node.
		firstChild := parent.RemoveFirstChild()
		NodeUtil_MarkFunctionsDeleted(firstChild, f.compiler)
		jutil.CheckState(!parent.HasChildren())
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

// Java: FunctionInjector.removeConstantVarAnnotation(Scope, String)
func FunctionInjector_RemoveConstantVarAnnotation(scope *Scope, name string) {
	var_ := scope.GetVar(name)
	var nameNode *rhino.Node
	if var_ != nil {
		nameNode = var_.GetNameNode()
	}
	if nameNode == nil {
		return
	}

	if nameNode.IsDeclaredConstantVar() {
		nameNode.SetDeclaredConstantVar(false)
	}
}

// FunctionInjector_IsDirectCallNodeReplacementPossible checks if the given
// function matches the criteria for an inlinable function, and if so, adds
// it to our set of inlinable functions.
//
// Java: FunctionInjector.isDirectCallNodeReplacementPossible(Node)
func FunctionInjector_IsDirectCallNodeReplacementPossible(fnNode *rhino.Node) bool {
	// Only inline single-statement functions
	block := NodeUtil_GetFunctionBody(fnNode)

	// Check if this function is suitable for direct replacement of a CALL node:
	// a function that consists of single return that returns an expression.
	if !block.HasChildren() {
		// special case empty functions.
		return true
	} else if block.HasOneChild() {
		// Only inline functions that return something.
		if block.GetFirstChild().IsReturn() && block.GetFirstFirstChild() != nil {
			return true
		}
	}

	return false
}

// FunctionInjector_CanInlineResult is the enum FunctionInjector.CanInlineResult.
type FunctionInjector_CanInlineResult int32

const (
	FunctionInjector_CanInlineResult_YES FunctionInjector_CanInlineResult = iota
	FunctionInjector_CanInlineResult_AFTER_PREPARATION
	FunctionInjector_CanInlineResult_NO
)

var functionInjector_CanInlineResult_names = [...]string{"YES", "AFTER_PREPARATION", "NO"}

// Java: FunctionInjector.CanInlineResult.values()
func FunctionInjector_CanInlineResult_values() []FunctionInjector_CanInlineResult {
	return []FunctionInjector_CanInlineResult{
		FunctionInjector_CanInlineResult_YES,
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		FunctionInjector_CanInlineResult_NO,
	}
}

// Java: FunctionInjector.CanInlineResult.name()
func (r FunctionInjector_CanInlineResult) String() string {
	return functionInjector_CanInlineResult_names[r]
}

// Java: FunctionInjector.CanInlineResult.ordinal()
func (r FunctionInjector_CanInlineResult) Ordinal() int { return int(r) }

// CanInlineReferenceAsStatementBlock determines whether a function can be
// inlined at a particular call site. There are several criteria that the
// function and reference must hold in order for the functions to be
// inlined: - It must be a simple call, or assignment, or var
// initialization.
//
//	f();
//	a = foo();
//	var a = foo();
//
// Java: FunctionInjector.canInlineReferenceAsStatementBlock(Reference, Node, ImmutableSet<String>)
func (f *FunctionInjector) CanInlineReferenceAsStatementBlock(
	ref FunctionInjector_Reference, fnNode *rhino.Node, namesToAlias *jutil.LinkedHashSet[string]) FunctionInjector_CanInlineResult {
	callSiteType := f.classifyCallSite(ref, f.allowDecomposition)
	if callSiteType == FunctionInjector_CallSiteType_UNSUPPORTED {
		return FunctionInjector_CanInlineResult_NO
	}

	if !f.allowDecomposition &&
		(callSiteType == FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION ||
			callSiteType == FunctionInjector_CallSiteType_EXPRESSION) {
		return FunctionInjector_CanInlineResult_NO
	}

	if !f.CallMeetsBlockInliningRequirements(ref, fnNode, namesToAlias) {
		return FunctionInjector_CanInlineResult_NO
	}

	if callSiteType == FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION ||
		callSiteType == FunctionInjector_CallSiteType_EXPRESSION {
		return FunctionInjector_CanInlineResult_AFTER_PREPARATION
	} else {
		return FunctionInjector_CanInlineResult_YES
	}
}

func functionInjector_isEvalName(n *rhino.Node) bool {
	return n.IsName() && n.GetString() == "eval"
}

func functionInjector_isShallowStatement(n *rhino.Node) bool {
	return n.IsRoot() || n.IsBlock() || (!n.IsFunction() && NodeUtil_IsShallowStatementTree(n.GetParent()))
}

func functionInjector_hasShallowDeclaration(n *rhino.Node) bool {
	if NodeUtil_IsDeclaration(n) {
		return true
	}
	if functionInjector_isShallowStatement(n) {
		for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
			if functionInjector_hasShallowDeclaration(c) {
				return true
			}
		}
	}
	return false
}

func functionInjector_countShallowReturns(n *rhino.Node) int {
	count := 0
	if n.IsReturn() {
		count++
	}
	if functionInjector_isShallowStatement(n) {
		for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
			count += functionInjector_countShallowReturns(c)
		}
	}
	return count
}

func functionInjector_referencesEval(n *rhino.Node, rootFn *rhino.Node) bool {
	if n.IsName() && scopeNameEq(n.GetString(), "eval") {
		return true
	}
	if n.IsFunction() && n != rootFn {
		return false
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		if c.GetFirstChild() == nil {
			if c.IsName() && scopeNameEq(c.GetString(), "eval") {
				return true
			}
			continue
		}
		if functionInjector_referencesEval(c, rootFn) {
			return true
		}
	}
	return false
}

// ReferencesEval returns whether or not fn includes a call to `eval` in its
// scope.
//
// Results are cached to make subsequent calls faster.
//
// Java: FunctionInjector.referencesEval(Node)
func (f *FunctionInjector) ReferencesEval(fn *rhino.Node) bool {
	jutil.CheckState(fn.IsFunction())

	info := f.fnInfoCache[fn]
	if info.evalState != 0 {
		return info.evalState == 2
	}

	result := functionInjector_referencesEval(fn, fn)
	if f.fnInfoCache == nil {
		f.fnInfoCache = make(map[*rhino.Node]functionInjectorFnInfo, 64)
	}
	if result {
		info.evalState = 2
	} else {
		info.evalState = 1
	}
	f.fnInfoCache[fn] = info
	return result
}

func functionInjector_findInnerFunctions(n *rhino.Node, containerFn *rhino.Node, first **rhino.Node, count *int) {
	if n != containerFn && NodeUtil_IsFunctionOrNonStaticClassFieldInitializer(n) {
		if *count == 0 {
			*first = n
		}
		*count++
		if *count > 1 {
			return
		}
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		if c.GetFirstChild() == nil {
			if c.IsMemberFieldDef() && !c.IsStaticMember() {
				if *count == 0 {
					*first = c
				}
				*count++
				if *count > 1 {
					return
				}
			}
			continue
		}
		functionInjector_findInnerFunctions(c, containerFn, first, count)
		if *count > 1 {
			return
		}
	}
}

// InnerFunctionOf returns any inner function or non-static class field
// initializer of containerFn.
//
// Non-static class field initializers are treated like inner functions
// because they are evaluated later, once per instantiation, and may capture
// locals.
//
// If there are no inner functions, or multilple inner functions, the
// sentinel values NO_FUNCTIONS, MULTIPLE_FUNCTIONS are returned
// respectively.
//
// Java: FunctionInjector.innerFunctionOf(Node)
func (f *FunctionInjector) InnerFunctionOf(containerFn *rhino.Node) *rhino.Node {
	jutil.CheckState(containerFn.IsFunction())

	info := f.fnInfoCache[containerFn]
	if info.inner != nil {
		return info.inner
	}

	var firstInner *rhino.Node
	count := 0
	functionInjector_findInnerFunctions(containerFn, containerFn, &firstInner, &count)

	var cached *rhino.Node
	switch count {
	case 0:
		cached = FunctionInjector_NO_FUNCTIONS
	case 1:
		cached = firstInner
	default:
		cached = FunctionInjector_MULTIPLE_FUNCTIONS
	}

	if f.fnInfoCache == nil {
		f.fnInfoCache = make(map[*rhino.Node]functionInjectorFnInfo, 64)
	}
	info.inner = cached
	f.fnInfoCache[containerFn] = info
	return cached
}

// CallMeetsBlockInliningRequirements determines whether a function can be
// inlined at a particular call site. - Don't inline if the calling function
// contains an inner function and inlining would introduce new globals.
//
// Java: FunctionInjector.callMeetsBlockInliningRequirements(Reference, Node, ImmutableSet<String>)
func (f *FunctionInjector) CallMeetsBlockInliningRequirements(
	callRef FunctionInjector_Reference, calleeFn *rhino.Node, namesToAlias *jutil.LinkedHashSet[string]) bool {
	// Note: functions that contain function definitions are filtered out
	// in isCandidateFunction.

	// TODO(johnlenz): Determining if the called function contains VARs
	// or if the caller contains inner functions accounts for 20% of the
	// run-time cost of this pass.
	//
	// Note that as of 2019-10-11, "eval" and function checking are cached so this may be less
	// of an issue.

	forbidTemps := false
	callRefScope := callRef.GetFunctionInjector_ReferenceBase().scope
	if !callRefScope.GetClosestHoistScope().IsGlobal() {
		callerFn := callRefScope.GetClosestHoistScope().GetRootNode().GetParent()

		// Don't allow any new vars into a scope that contains eval or one
		// that contains functions (excluding the function being inlined).
		if f.ReferencesEval(callerFn) {
			forbidTemps = true
		} else if !f.assumeMinimumCapture {
			innerFn := f.InnerFunctionOf(callerFn)
			calleeIsOnlyInnerFn := innerFn == FunctionInjector_NO_FUNCTIONS || innerFn == calleeFn
			forbidTemps = !calleeIsOnlyInnerFn
		}
	}

	// If the caller contains functions or evals, verify we aren't adding any
	// additional VAR declarations because aliasing is needed.
	if forbidTemps {
		if !namesToAlias.IsEmpty() {
			return false
		}

		// Don't inline functions with var declarations into a scope with inner
		// functions as the new vars would leak into the inner function and
		// cause memory leaks.
		calleeContainsVars := functionInjector_hasShallowDeclaration(NodeUtil_GetFunctionBody(calleeFn))
		if calleeContainsVars {
			return false
		}

		args := f.functionArgumentInjector.GetFunctionCallParameterMapScratch(
			calleeFn, callRef.GetFunctionInjector_ReferenceBase().callNode, f.safeNameIdSupplier)
		hasArgs := !args.IsEmpty()
		if hasArgs {
			if f.functionArgumentInjector.HasCallArgumentsNeedingTemps(
				f.compiler, calleeFn, args, namesToAlias, f.compiler.GetCodingConvention()) {
				return false
			}
		}
	}

	return true
}

// CanInlineReferenceDirectly determines whether a function can be inlined at
// a particular call site. There are several criteria that the function and
// reference must hold in order for the functions to be inlined: 1) If a
// call's arguments have side effects, the corresponding argument in the
// function must only be referenced once. For instance, this will not be
// inlined:
//
//	function foo(a) { return a + a }
//	x = foo(i++);
//
// Java: FunctionInjector.canInlineReferenceDirectly(Reference, Node, ImmutableSet<String>)
func (f *FunctionInjector) CanInlineReferenceDirectly(
	ref FunctionInjector_Reference, fnNode *rhino.Node, namesToAlias *jutil.LinkedHashSet[string]) FunctionInjector_CanInlineResult {
	if !FunctionInjector_IsDirectCallNodeReplacementPossible(fnNode) {
		return FunctionInjector_CanInlineResult_NO
	}

	// CALL NODE: [ NAME, ARG1, ARG2, ... ]
	callNode := ref.GetFunctionInjector_ReferenceBase().callNode
	cArg := callNode.GetSecondChild()

	// Functions called via 'call' and 'apply' have a this-object as
	// the first parameter, but this is not part of the called function's
	// parameter list.
	if !callNode.GetFirstChild().IsName() {
		if NodeUtil_IsFunctionObjectCall(callNode) {
			// TODO(johnlenz): Support replace this with a value.
			if !fnNode.IsArrowFunction() && (cArg == nil || !cArg.IsThis()) {
				return FunctionInjector_CanInlineResult_NO
			}
			if cArg != nil {
				cArg = cArg.GetNext()
			}
		} else {
			// ".apply" call should be filtered before this.
			jutil.CheckState(!NodeUtil_IsFunctionObjectApply(callNode))
		}
	}

	if !namesToAlias.IsEmpty() {
		return FunctionInjector_CanInlineResult_NO
	}

	args := f.functionArgumentInjector.GetFunctionCallParameterMapScratch(
		fnNode, callNode, f.throwawayNameSupplier)
	hasArgs := !args.IsEmpty()
	if hasArgs {
		if f.functionArgumentInjector.HasCallArgumentsNeedingTemps(
			f.compiler, fnNode, args, namesToAlias, f.compiler.GetCodingConvention()) {
			return FunctionInjector_CanInlineResult_NO
		}
	}

	return FunctionInjector_CanInlineResult_YES
}

// InliningLowersCost determines if inlining the function is likely to reduce
// the code size.
//
// Go-only: language-adapter — `Collection<? extends Reference> refs` is a
// slice of the interface.
//
// Java: FunctionInjector.inliningLowersCost(JSChunk, Node, Collection<? extends Reference>, Set<String>, boolean, boolean)
func (f *FunctionInjector) InliningLowersCost(
	fnChunk *JSChunk,
	fnNode *rhino.Node,
	refs []FunctionInjector_Reference,
	namesToAlias *jutil.LinkedHashSet[string],
	isRemovable bool,
	referencesThis bool) bool {
	referenceCount := len(refs)
	if referenceCount == 0 {
		return true
	}

	referencesUsingBlockInlining := 0

	checkChunks := isRemovable && fnChunk != nil
	chunkGraph := f.compiler.GetChunkGraph()

	for _, ref := range refs {
		refBase := ref.GetFunctionInjector_ReferenceBase()
		if refBase.mode == FunctionInjector_InliningMode_BLOCK {
			referencesUsingBlockInlining++
		}

		// Check if any of the references cross the chunk boundaries.
		if checkChunks && refBase.chunk != nil {
			if refBase.chunk != fnChunk && !chunkGraph.DependsOn(refBase.chunk, fnChunk) {
				// Calculate the cost as if the function were non-removable,
				// if it still lowers the cost inline it.
				isRemovable = false
				checkChunks = false // no need to check additional chunks.
			}
		}
	}

	referencesUsingDirectInlining := referenceCount - referencesUsingBlockInlining

	// Don't bother calculating the cost of function for simple functions where
	// possible.
	// However, when inlining a complex function, even a single reference may be
	// larger than the original function if there are many returns (resulting
	// in additional assignments) or many parameters that need to be aliased
	// so use the cost estimating.
	if referenceCount == 1 && isRemovable && referencesUsingDirectInlining == 1 {
		return true
	}

	callCost := FunctionInjector_EstimateCallCost(fnNode, referencesThis)
	overallCallCost := callCost * referenceCount

	costDeltaDirect := FunctionInjector_InlineCostDelta(fnNode, namesToAlias, FunctionInjector_InliningMode_DIRECT)
	costDeltaBlock := FunctionInjector_InlineCostDelta(fnNode, namesToAlias, FunctionInjector_InliningMode_BLOCK)

	return FunctionInjector_DoesLowerCost(
		fnNode,
		overallCallCost,
		referencesUsingDirectInlining,
		costDeltaDirect,
		referencesUsingBlockInlining,
		costDeltaBlock,
		isRemovable)
}

// FunctionInjector_DoesLowerCost returns whether inlining will lower cost.
//
// Java: FunctionInjector.doesLowerCost(Node, int, int, int, int, int, boolean)
func FunctionInjector_DoesLowerCost(
	fnNode *rhino.Node,
	callCost int,
	directInlines int,
	costDeltaDirect int,
	blockInlines int,
	costDeltaBlock int,
	removable bool) bool {

	// Determine the threshold value for this inequality:
	//     inline_cost < call_cost
	// But solve it for the function declaration size so the size of it
	// is only calculated once and terminated early if possible.

	removableCount := 0
	if removable {
		removableCount = 1
	}
	fnInstanceCount := directInlines + blockInlines - removableCount
	// Prevent division by zero.
	if fnInstanceCount == 0 {
		// Special case single reference function that are being block inlined:
		// If the cost of the inline is greater than the function definition size,
		// don't inline.
		return blockInlines <= 0 || costDeltaBlock <= 0
	}

	costDelta := (directInlines * -costDeltaDirect) + (blockInlines * -costDeltaBlock)
	threshold := (callCost + costDelta) / fnInstanceCount

	return InlineCostEstimator_GetCostWithThreshold(fnNode, threshold+1) <= threshold
}

// FunctionInjector_EstimateCallCost gets an estimate of the cost in
// characters of making the function call: the sum of the identifiers and the
// separators.
//
// Java: FunctionInjector.estimateCallCost(Node, boolean)
func FunctionInjector_EstimateCallCost(fnNode *rhino.Node, referencesThis bool) int {
	argsNode := NodeUtil_GetFunctionParameters(fnNode)
	numArgs := argsNode.GetChildCount()

	callCost := FunctionInjector_NAME_COST_ESTIMATE + FunctionInjector_PAREN_COST
	if numArgs > 0 {
		callCost += (numArgs * FunctionInjector_NAME_COST_ESTIMATE) + ((numArgs - 1) * FunctionInjector_COMMA_COST)
	}

	if referencesThis {
		// TODO(johnlenz): Update this if we start supporting inlining
		// other functions that reference this.
		// The only functions that reference this that are currently inlined
		// are those that are called via ".call" with an explicit "this".
		callCost += 5 + 5 // ".call" + "this,"
	}

	return callCost
}

// FunctionInjector_InlineCostDelta returns the difference between the
// function definition cost and inline cost.
//
// Java: FunctionInjector.inlineCostDelta(Node, Set<String>, InliningMode)
func FunctionInjector_InlineCostDelta(fnNode *rhino.Node, namesToAlias *jutil.LinkedHashSet[string], mode FunctionInjector_InliningMode) int {
	// The part of the function that is never inlined:
	//    "function xx(xx,xx){}" (15 + (param count * 3) -1;
	paramCount := NodeUtil_GetFunctionParameters(fnNode).GetChildCount()
	commaCount := 0
	if paramCount > 1 {
		commaCount = paramCount - 1
	}
	costDeltaFunctionOverhead :=
		15 + commaCount + (paramCount * InlineCostEstimator_ESTIMATED_IDENTIFIER_COST)

	block := fnNode.GetLastChild()
	if !block.HasChildren() {
		// Assume the inline cost is zero for empty functions.
		return -costDeltaFunctionOverhead
	}

	if mode == FunctionInjector_InliningMode_DIRECT {
		// The part of the function that is inlined using direct inlining:
		//    "return " (7)
		return -(costDeltaFunctionOverhead + 7)
	} else {
		aliasCount := namesToAlias.Size()

		// Originally, we estimated purely base on the function code size, relying
		// on later optimizations. But that did not produce good results, so here
		// we try to estimate the something closer to the actual inlined coded.

		// NOTE 1: Result overhead is only if there is an assignment, but
		// getting that information would require some refactoring.
		// NOTE 2: The aliasing overhead is currently an under-estimate,
		// as some parameters are aliased because of the parameters used.
		// Perhaps we should just assume all parameters will be aliased?
		const inlineBlockOverhead = 4     // "X:{}"
		const perReturnOverhead = 2       // "return" --> "break X"
		const perReturnResultOverhead = 3 // "XX="
		const perAliasOverhead = 3        // "XX="

		// TODO(johnlenz): Counting the number of returns is relatively expensive.
		//   This information should be determined during the traversal and cached.
		returnCount := functionInjector_countShallowReturns(block)
		resultCount := 0
		if returnCount > 0 {
			resultCount = returnCount - 1
		}
		baseOverhead := 0
		if returnCount > 0 {
			baseOverhead = inlineBlockOverhead
		}

		overhead := baseOverhead +
			returnCount*perReturnOverhead +
			resultCount*perReturnResultOverhead +
			aliasCount*perAliasOverhead

		return overhead - costDeltaFunctionOverhead
	}
}

// SetKnownConstantFunctions stores the names of known constants to be used
// when classifying call-sites in expressions.
//
// Java: FunctionInjector.setKnownConstantFunctions(ImmutableSet<String>)
func (f *FunctionInjector) SetKnownConstantFunctions(knownConstantFunctions *jutil.LinkedHashSet[string]) {
	// This is only expected to be set once. The same set should be used
	// when evaluating call-sites and inlining calls.
	jutil.CheckState(f.knownConstantFunctions.IsEmpty())
	f.knownConstantFunctions = knownConstantFunctions
}
