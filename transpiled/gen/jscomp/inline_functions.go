/*
 * Copyright 2005 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/InlineFunctions.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_inline_functions_1 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_inline_functions_2 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_inline_functions_3 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_inline_functions_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool {
		return NodeUtil_IsFunctionOrNonStaticClassFieldInitializer(pkg_rhino.Node(rawNode_))
	})
	lambdaSingleton_inline_functions_2 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_1 uint32) bool {
		n := pkg_rhino.Node(rawNode_1)
		_ = n
		return (n.IsLet() || n.IsConst())
	})
	lambdaSingleton_inline_functions_3 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_6 uint32) bool {
		input := pkg_rhino.Node(rawNode_6)
		_ = input
		return (input.IsDefaultValue() || input.IsDestructuringPattern())
	})
}

// Java class: com.google.javascript.jscomp.InlineFunctions
type InlineFunctions struct {
	fns                                     j2gort.OrderedMap[string, *InlineFunctions_FunctionState] `gson:"fns"`
	anonFns                                 j2gort.OrderedMap[pkg_rhino.Node, string]                 `gson:"anonFns"`
	compiler                                AbstractCompiler                                          `gson:"compiler"`
	injector                                *FunctionInjector                                         `gson:"injector"`
	functionArgumentInjector                *FunctionArgumentInjector                                 `gson:"functionArgumentInjector"`
	maxSizeAfterInlining                    int32                                                     `gson:"maxSizeAfterInlining"`
	idHash_                                 atomic.Int32
	reach                                   CompilerOptions_Reach `gson:"reach"`
	enforceMaxSizeAfterInlining             bool                  `gson:"enforceMaxSizeAfterInlining"`
	slab_InlineFunctions_Reference          j2gort.PooledSlab[InlineFunctions_Reference]
	slab_InlineFunctions_FunctionVar        j2gort.PooledSlab[InlineFunctions_FunctionVar]
	slab_InlineFunctions_NamedFunction      j2gort.PooledSlab[InlineFunctions_NamedFunction]
	slab_InlineFunctions_FunctionExpression j2gort.PooledSlab[InlineFunctions_FunctionExpression]
}

func (this_ *InlineFunctions) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineFunctions) IsJavaNil_() bool { return this_ == nil }
func (*InlineFunctions) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineFunctions", "InlineFunctions"
}
func (this_ *InlineFunctions) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineFunctions) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineFunctions", "com.google.javascript.jscomp.InlineFunctions", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineFunctions](v_)
	})
}
func (*InlineFunctions) Implements_jscomp_CompilerPass() {}

var static_InlineFunctions_MISSED_REQUIRED_INLINING *DiagnosticType

func StaticGet_InlineFunctions_MISSED_REQUIRED_INLINING() *DiagnosticType {
	clinit_InlineFunctions()
	return static_InlineFunctions_MISSED_REQUIRED_INLINING
}
func StaticSet_InlineFunctions_MISSED_REQUIRED_INLINING(v *DiagnosticType) {
	clinit_InlineFunctions()
	static_InlineFunctions_MISSED_REQUIRED_INLINING = v
}

var clinitDone_InlineFunctions uint32
var clinitState_InlineFunctions j2gort.ClinitState

func clinit_InlineFunctions() {
	if atomic.LoadUint32(&clinitDone_InlineFunctions) != 0 {
		return
	}
	clinitSlow_InlineFunctions()
}

func clinitSlow_InlineFunctions() {
	if !j2gort.EnterClinit(&clinitState_InlineFunctions, &clinitDone_InlineFunctions) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_InlineFunctions, &clinitDone_InlineFunctions)
	static_InlineFunctions_MISSED_REQUIRED_INLINING = DiagnosticType_Warning("JSC_MISSED_REQUIRED_INLINING", "function {1} annotated @requireInlining could not be inlined here")
}

func (this_ *InlineFunctions) initFields_InlineFunctions() {
	_ = this_
	this_.fns = j2gort.AcquireOrderedMapCap[string, *InlineFunctions_FunctionState](int32(1024))
	this_.anonFns = j2gort.AcquireOrderedMap[pkg_rhino.Node, string]()
}

// Java constructor: InlineFunctions(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.function.Supplier<java.lang.String> safeNameIdSupplier, com.google.javascript.jscomp.CompilerOptions.Reach reach, boolean assumeStrictThis, int maxSizeAfterInlining)
func NewInlineFunctions(compiler AbstractCompiler, safeNameIdSupplier j2gort.Ext_java_util_function_Supplier, reach CompilerOptions_Reach, assumeStrictThis bool, maxSizeAfterInlining int32) *InlineFunctions {
	clinit_InlineFunctions()
	this_ := &InlineFunctions{}
	InitInlineFunctions(this_, compiler, safeNameIdSupplier, reach, assumeStrictThis, maxSizeAfterInlining)
	return this_
}

func InitInlineFunctions(this_ *InlineFunctions, compiler AbstractCompiler, safeNameIdSupplier j2gort.Ext_java_util_function_Supplier, reach CompilerOptions_Reach, assumeStrictThis bool, maxSizeAfterInlining int32) {
	_ = compiler
	_ = safeNameIdSupplier
	_ = reach
	_ = assumeStrictThis
	_ = maxSizeAfterInlining
	this_.initFields_InlineFunctions()
	if !(compiler != nil && !compiler.IsJavaNil_()) {
		j2gort.CheckArgument(false)
	}
	if !(safeNameIdSupplier != nil) {
		j2gort.CheckArgument(false)
	}
	if !(reach != StaticGet_CompilerOptions_Reach_NONE()) {
		j2gort.CheckArgument(false)
	}
	this_.compiler = compiler
	this_.reach = reach
	this_.maxSizeAfterInlining = maxSizeAfterInlining
	this_.enforceMaxSizeAfterInlining = (maxSizeAfterInlining != int32(-1))
	this_.functionArgumentInjector = NewFunctionArgumentInjector(compiler.GetAstAnalyzer())
	this_.injector = NewFunctionInjector_Builder(compiler).SafeNameIdSupplier(safeNameIdSupplier).AssumeStrictThis(assumeStrictThis).FunctionArgumentInjector(this_.functionArgumentInjector).Build()
}

// Java method: InlineFunctions.getOrCreateFunctionState(java.lang.String fnName)
func (this_ *InlineFunctions) priv_InlineFunctions_getOrCreateFunctionState_3269(fnName string) *InlineFunctions_FunctionState {
	_ = fnName
	var fs *InlineFunctions_FunctionState = this_.fns.Get(fnName)
	_ = fs
	if fs == nil {
		fs = NewInlineFunctions_FunctionState()
		this_.fns.Put(fnName, fs)
	}
	return fs
}

// Java method: InlineFunctions.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InlineFunctions) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer this_.ReleaseSlabs_()
	if !(this_.compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckState(false)
	}
	NodeTraversal_Traverse(this_.compiler, root, NewInlineFunctions_FindCandidateFunctions(this_))
	if this_.fns.IsEmpty() {
		return
	}
	sc_ := NewSyntacticScopeCreator_1(this_.compiler)
	defer func() {
		for cur_ := this_.fns.IterValues(); cur_.HasNext(); {
			fs_ := cur_.Next()
			if fs_ != nil && fs_.references != nil {
				j2gort.ReleaseOrderedMap(fs_.references)
				fs_.references = nil
			}
		}
		this_.fns.ReleaseToPool()
		this_.anonFns.ReleaseToPool()
		sc_.ReleaseSlabs_()
	}()
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(NewInlineFunctions_FindCandidatesReferences(this_))
		(&ntb_).SetScopeCreator(sc_)
		(&ntb_).Traverse(root)
	}

	this_.priv_InlineFunctions_trimCandidatesNotMeetingMinimumRequirements_3275()
	if this_.fns.IsEmpty() {
		return
	}
	var fnNames j2gort.Ext_com_google_common_collect_ImmutableSet = this_.fns.KeysToImmutableSet()
	_ = fnNames
	this_.injector.SetKnownConstantFunctions(fnNames)
	this_.priv_InlineFunctions_trimCandidatesUsingOnCost_3276()
	if this_.fns.IsEmpty() {
		return
	}
	this_.priv_InlineFunctions_resolveInlineConflicts_3279()
	this_.priv_InlineFunctions_decomposeExpressions_3282()
	NodeTraversal_Traverse(this_.compiler, root, Upcast_NodeTraversal_Callback(NewInlineFunctions_CallVisitor(this_, NewInlineFunctions_Inline(this_.injector))))
	this_.priv_InlineFunctions_removeInlinedFunctions_3283()
}

// Java static method: InlineFunctions.isAlwaysInlinable(com.google.javascript.rhino.Node fn)
func InlineFunctions_IsAlwaysInlinable(fn pkg_rhino.Node) bool {
	clinit_InlineFunctions()
	_ = fn
	if !(fn.IsFunction()) {
		j2gort.CheckArgument(false)
	}
	var body pkg_rhino.Node = NodeUtil_GetFunctionBody(fn)
	_ = body
	return ((!(body.HasChildren())) || (body.HasOneChild() && body.GetFirstChild().IsReturn()))
}

// Java method: InlineFunctions.targetSizeAfterInlineExceedsLimit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.InlineFunctions.FunctionState functionState)
func (this_ *InlineFunctions) priv_InlineFunctions_targetSizeAfterInlineExceedsLimit_3270(t *NodeTraversal, functionState *InlineFunctions_FunctionState) bool {
	_ = t
	_ = functionState
	var containingFunction pkg_rhino.Node = t.GetEnclosingFunction()
	_ = containingFunction
	if containingFunction == 0 {
		return false
	}
	var inlinedFun pkg_rhino.Node = functionState.GetFn().GetFunctionNode()
	_ = inlinedFun
	if InlineFunctions_IsAlwaysInlinable(inlinedFun) {
		return false
	}
	var inlinedFunSize int32 = NodeUtil_CountAstSizeUpToLimit(NodeUtil_GetFunctionBody(inlinedFun), this_.maxSizeAfterInlining)
	_ = inlinedFunSize
	var targetFunSize int32 = NodeUtil_CountAstSizeUpToLimit(containingFunction, this_.maxSizeAfterInlining)
	_ = targetFunSize
	return ((inlinedFunSize + targetFunSize) > this_.maxSizeAfterInlining)
}

// Java method: InlineFunctions.maybeAddFunction(com.google.javascript.jscomp.InlineFunctions.Function fn, com.google.javascript.jscomp.JSChunk chunk)
func (this_ *InlineFunctions) priv_InlineFunctions_maybeAddFunction_3271(fn InlineFunctions_Function, chunk *JSChunk) {
	_ = fn
	_ = chunk
	var name string = fn.GetName()
	_ = name
	var functionState *InlineFunctions_FunctionState = this_.priv_InlineFunctions_getOrCreateFunctionState_3269(name)
	_ = functionState
	this_.priv_InlineFunctions_updateFunctionStateForInlining_3272(fn, chunk, functionState)
	if !(!(InlineFunctions_HasRequireInliningAnnotation(fn.GetFunctionNode()) && !(functionState.CanInline()))) {
		j2gort.CheckState(false)
	}
}

// Java method: InlineFunctions.updateFunctionStateForInlining(com.google.javascript.jscomp.InlineFunctions.Function fn, com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.InlineFunctions.FunctionState functionState)
func (this_ *InlineFunctions) priv_InlineFunctions_updateFunctionStateForInlining_3272(fn InlineFunctions_Function, chunk *JSChunk, functionState *InlineFunctions_FunctionState) {
	_ = fn
	_ = chunk
	_ = functionState
	if functionState.HasExistingFunctionDefinition() {
		functionState.DisallowInlining_2(this_.compiler, StaticGet_InlineFunctions_DisallowInliningReason_MULTIPLE_FUNCTION_DEFINITIONS())
		return
	}
	var fnNode pkg_rhino.Node = fn.GetFunctionNode()
	_ = fnNode
	if this_.priv_InlineFunctions_hasNoInlineAnnotation_3273(fnNode) {
		functionState.DisallowInlining_2(this_.compiler, StaticGet_InlineFunctions_DisallowInliningReason_NO_INLINE_ANNOTATION())
		return
	}
	if (this_.enforceMaxSizeAfterInlining && !(InlineFunctions_IsAlwaysInlinable(fnNode))) && (this_.maxSizeAfterInlining <= NodeUtil_CountAstSizeUpToLimit(fnNode, this_.maxSizeAfterInlining)) {
		functionState.DisallowInlining_2(this_.compiler, StaticGet_InlineFunctions_DisallowInliningReason_TARGET_SIZE_EXCEEDS_LIMIT())
		return
	}
	if functionState.CanInline() {
		functionState.SetFn(fn)
		if FunctionInjector_IsDirectCallNodeReplacementPossible(fn.GetFunctionNode()) {
			functionState.InlineDirectly(true)
		}
		if InlineFunctions_HasNonInlinableParam(NodeUtil_GetFunctionParameters(fnNode)) {
			functionState.DisallowInlining_2(this_.compiler, StaticGet_InlineFunctions_DisallowInliningReason_NON_INLINABLE_PARAM())
		}
		if !(this_.priv_InlineFunctions_isCandidateFunction_3274(fn)) {
			functionState.DisallowInlining_2(this_.compiler, StaticGet_InlineFunctions_DisallowInliningReason_NOT_CANDIDATE_FUNCTION())
		}
		if functionState.CanInline() {
			functionState.SetChunk(chunk)
			var namesToAlias j2gort.Ext_java_util_Set = j2gort.AsIface[j2gort.Ext_java_util_Set](this_.functionArgumentInjector.FindModifiedParameters(fnNode))
			_ = namesToAlias
			if !(namesToAlias.IsEmpty()) {
				functionState.InlineDirectly(false)
				functionState.SetNamesToAlias(namesToAlias)
			}
			var block pkg_rhino.Node = NodeUtil_GetFunctionBody(fnNode)
			_ = block
			if NodeUtil_ReferencesEnclosingReceiver(block) {
				functionState.SetReferencesThis(true)
			}
			if NodeUtil_Has(block, lambdaSingleton_inline_functions_1, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue()) {
				functionState.SetHasInnerClosures(true)
				if InlineFunctions_HasLocalNames(fnNode) {
					functionState.DisallowInlining_2(this_.compiler, StaticGet_InlineFunctions_DisallowInliningReason_HAS_LOCAL_NAMES())
				}
			}
		}
		if fnNode.GetGrandparent().IsVar() {
			var block_2 pkg_rhino.Node = functionState.GetFn().GetDeclaringBlock()
			_ = block_2
			if (block_2.IsBlock() && !(block_2.GetParent().IsFunction())) && NodeUtil_Has(block_2, lambdaSingleton_inline_functions_2, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue()) {
				functionState.DisallowInlining_2(this_.compiler, StaticGet_InlineFunctions_DisallowInliningReason_CAPTURED_VARIABLE_NOT_IN_SCOPE())
			}
		}
		if fnNode.IsGeneratorFunction() {
			functionState.DisallowInlining_2(this_.compiler, StaticGet_InlineFunctions_DisallowInliningReason_GENERATOR_FUNCTION())
		}
		if fnNode.IsAsyncFunction() {
			functionState.DisallowInlining_2(this_.compiler, StaticGet_InlineFunctions_DisallowInliningReason_ASYNC_FUNCTION())
		}
	}
}

// Java method: InlineFunctions.hasNoInlineAnnotation(com.google.javascript.rhino.Node fnNode)
func (this_ *InlineFunctions) priv_InlineFunctions_hasNoInlineAnnotation_3273(fnNode pkg_rhino.Node) bool {
	_ = fnNode
	var jsDocInfo *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(fnNode)
	_ = jsDocInfo
	return ((jsDocInfo != nil) && jsDocInfo.IsNoInline())
}

// Java static method: InlineFunctions.hasEncourageInliningAnnotation(com.google.javascript.rhino.Node fnNode)
func InlineFunctions_HasEncourageInliningAnnotation(fnNode pkg_rhino.Node) bool {
	clinit_InlineFunctions()
	_ = fnNode
	var jsDocInfo *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(fnNode)
	_ = jsDocInfo
	return ((jsDocInfo != nil) && jsDocInfo.IsEncourageInlining())
}

// Java static method: InlineFunctions.hasRequireInliningAnnotation(com.google.javascript.rhino.Node fnNode)
func InlineFunctions_HasRequireInliningAnnotation(fnNode pkg_rhino.Node) bool {
	clinit_InlineFunctions()
	_ = fnNode
	var jsDocInfo *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(fnNode)
	_ = jsDocInfo
	return ((jsDocInfo != nil) && jsDocInfo.IsRequireInlining())
}

// Java static method: InlineFunctions.hasLocalNames(com.google.javascript.rhino.Node fnNode)
func InlineFunctions_HasLocalNames(fnNode pkg_rhino.Node) bool {
	clinit_InlineFunctions()
	_ = fnNode
	var block pkg_rhino.Node = NodeUtil_GetFunctionBody(fnNode)
	_ = block
	return (NodeUtil_GetFunctionParameters(fnNode).HasChildren() || NodeUtil_Has(block, NewNodeUtil_MatchDeclaration(), NewNodeUtil_MatchShallowStatement()))
}

// Java method: InlineFunctions.isCandidateFunction(com.google.javascript.jscomp.InlineFunctions.Function fn)
func (this_ *InlineFunctions) priv_InlineFunctions_isCandidateFunction_3274(fn InlineFunctions_Function) bool {
	_ = fn
	var fnName string = fn.GetName()
	_ = fnName
	if this_.compiler.GetCodingConvention().IsExported_1(fnName) {
		return false
	}
	if this_.compiler.GetCodingConvention().IsPropertyRenameFunction(fn.GetNameNode()) {
		return false
	}
	var fnNode pkg_rhino.Node = fn.GetFunctionNode()
	_ = fnNode
	return this_.injector.DoesFunctionMeetMinimumRequirements(fnName, fnNode)
}

// Java static method: InlineFunctions.isCandidateUsage(com.google.javascript.rhino.Node name)
func InlineFunctions_IsCandidateUsage(name pkg_rhino.Node) bool {
	clinit_InlineFunctions()
	_ = name
	var parent pkg_rhino.Node = name.GetParent()
	_ = parent
	if !(name.IsName()) {
		j2gort.CheckState(false)
	}
	if NodeUtil_IsNameDeclaration(parent) || parent.IsFunction() {
		return true
	}
	if NodeUtil_IsNormalOrOptChainCall(parent) && (parent.GetFirstChild() == name) {
		return true
	}
	if (((parent.IsGetElem() && (name == parent.GetFirstChild())) && parent.GetSecondChild().IsStringLit()) && (j2gort.CheckNotNullString(parent.GetSecondChild().GetString()) == "call")) || (parent.IsGetProp() && (j2gort.CheckNotNullString(parent.GetString()) == "call")) {
		var grandparent pkg_rhino.Node = name.GetAncestor(int32(2))
		_ = grandparent
		if grandparent.IsCall() && (grandparent.GetFirstChild() == parent) {
			return true
		}
	}
	return false
}

// Java method: InlineFunctions.trimCandidatesNotMeetingMinimumRequirements()
func (this_ *InlineFunctions) priv_InlineFunctions_trimCandidatesNotMeetingMinimumRequirements_3275() {
	var toRemove []string = nil
	_ = toRemove
	{
		cur_2 := this_.fns.IterEntries()
		for cur_2.NextEntry() {
			var functionState *InlineFunctions_FunctionState = cur_2.Value()
			_ = functionState
			if !(functionState.HasExistingFunctionDefinition()) || !(functionState.CanInline()) {
				if toRemove == nil {
					toRemove = []string{}
				}
				toRemove = append(toRemove, cur_2.Key())
			}
		}
	}
	if toRemove != nil {
		{
			var i int32 = int32(0)
			_ = i
			var size int32 = int32(len(toRemove))
			_ = size
			for ; i < size; i = i + 1 {
				if fs_ := this_.fns.Remove(toRemove[i]); fs_ != nil && fs_.references != nil {
					j2gort.ReleaseOrderedMap(fs_.references)
					fs_.references = nil
				}

			}
		}
	}
}

// Java method: InlineFunctions.trimCandidatesUsingOnCost()
func (this_ *InlineFunctions) priv_InlineFunctions_trimCandidatesUsingOnCost_3276() {
	var toRemove []string = nil
	_ = toRemove
	{
		cur_3 := this_.fns.IterEntries()
		for cur_3.NextEntry() {
			var functionState *InlineFunctions_FunctionState = cur_3.Value()
			_ = functionState
			if functionState.HasReferences() {
				var lowersCost bool = this_.priv_InlineFunctions_minimizeCost_3277(functionState)
				_ = lowersCost
				if !(lowersCost) {
					if toRemove == nil {
						toRemove = []string{}
					}
					toRemove = append(toRemove, cur_3.Key())
				}
			} else if !(functionState.CanRemove()) {
				if toRemove == nil {
					toRemove = []string{}
				}
				toRemove = append(toRemove, cur_3.Key())
			}
		}
	}
	if toRemove != nil {
		{
			var i int32 = int32(0)
			_ = i
			var size int32 = int32(len(toRemove))
			_ = size
			for ; i < size; i = i + 1 {
				if fs_ := this_.fns.Remove(toRemove[i]); fs_ != nil && fs_.references != nil {
					j2gort.ReleaseOrderedMap(fs_.references)
					fs_.references = nil
				}

			}
		}
	}
}

// Java method: InlineFunctions.minimizeCost(com.google.javascript.jscomp.InlineFunctions.FunctionState functionState)
func (this_ *InlineFunctions) priv_InlineFunctions_minimizeCost_3277(functionState *InlineFunctions_FunctionState) bool {
	_ = functionState
	if !(this_.priv_InlineFunctions_inliningLowersCost_3278(functionState)) {
		if functionState.HasBlockInliningReferences() {
			functionState.SetRemove_4(false, this_.compiler, StaticGet_InlineFunctions_CannotRemoveReason_BLOCK_INLINING_REFERENCE(), StaticGet_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES())
			functionState.RemoveBlockInliningReferences()
			if !(functionState.HasReferences()) || !(this_.priv_InlineFunctions_inliningLowersCost_3278(functionState)) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

// Java method: InlineFunctions.inliningLowersCost(com.google.javascript.jscomp.InlineFunctions.FunctionState functionState)
func (this_ *InlineFunctions) priv_InlineFunctions_inliningLowersCost_3278(functionState *InlineFunctions_FunctionState) bool {
	_ = functionState
	if functionState.RequireInlining() || functionState.EncourageInlining() {
		return true
	}
	if functionState.HasInlineReferences() {
		return this_.injector.InliningLowersCostInline(functionState.GetChunk(), functionState.GetFn().GetFunctionNode(), ToIface_FunctionInjector_Reference(functionState.InlineRef0()), ToIface_FunctionInjector_Reference(functionState.InlineRef1()), functionState.GetNamesToAliasReadOnly(), functionState.CanRemove(), functionState.GetReferencesThis())
	}
	return this_.injector.InliningLowersCost(functionState.GetChunk(), functionState.GetFn().GetFunctionNode(), functionState.GetReferences(), functionState.GetNamesToAliasReadOnly(), functionState.CanRemove(), functionState.GetReferencesThis())
}

// Java method: InlineFunctions.resolveInlineConflicts()
func (this_ *InlineFunctions) priv_InlineFunctions_resolveInlineConflicts_3279() {
	{
		cur_4 := this_.fns.IterValues()
		for cur_4.HasNext() {

			functionState := cur_4.
				Next()
			_ = functionState
			this_.priv_InlineFunctions_resolveInlineConflictsForFunction_3280(functionState)
		}
	}
}

// Java static method: InlineFunctions.hasNonInlinableParam(com.google.javascript.rhino.Node node)
func InlineFunctions_HasNonInlinableParam(node pkg_rhino.Node) bool {
	clinit_InlineFunctions()
	_ = node
	j2gort.CheckNotNullNode(node)
	var pred j2gort.Ext_com_google_common_base_Predicate = lambdaSingleton_inline_functions_3
	_ = pred
	return NodeUtil_Has(node, pred, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue())
}

// Java method: InlineFunctions.resolveInlineConflictsForFunction(com.google.javascript.jscomp.InlineFunctions.FunctionState functionState)
func (this_ *InlineFunctions) priv_InlineFunctions_resolveInlineConflictsForFunction_3280(functionState *InlineFunctions_FunctionState) {
	_ = functionState
	if !(functionState.HasReferences()) || !(functionState.CanInline()) {
		return
	}
	var fnNode pkg_rhino.Node = functionState.GetFn().GetFunctionNode()
	_ = fnNode
	var names j2gort.OrderedSet[string] = this_.priv_InlineFunctions_findCalledFunctions_3281(fnNode)
	_ = names
	if !(names.IsEmpty()) {
		{
			cur_7 := names.Iter()
			for cur_7.HasNext() {
				name := cur_7.
					Next()
				_ = name
				var fsCalled *InlineFunctions_FunctionState = this_.fns.Get(name)
				_ = fsCalled
				if (fsCalled != nil) && fsCalled.CanRemove() {
					fsCalled.SetRemove_4(false, this_.compiler, StaticGet_InlineFunctions_CannotRemoveReason_CALLED_BY_INLINED_FUNCTION(), StaticGet_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO())
					if (!(this_.priv_InlineFunctions_minimizeCost_3277(fsCalled)) && !(fsCalled.RequireInlining())) && !(fsCalled.EncourageInlining()) {
						fsCalled.DisallowInlining_2(this_.compiler, StaticGet_InlineFunctions_DisallowInliningReason_COST())
					}
				}
			}
		}
		functionState.SetSafeFnNode(functionState.GetFn().GetFunctionNode().CloneTree_0())
	}
}

// Java method: InlineFunctions.findCalledFunctions(com.google.javascript.rhino.Node node)
func (this_ *InlineFunctions) priv_InlineFunctions_findCalledFunctions_3281(node pkg_rhino.Node) j2gort.OrderedSet[string] {
	_ = node
	var changed j2gort.OrderedSet[string] = j2gort.NewOrderedSet[string]()
	_ = changed
	ordSet_InlineFunctions_FindCalledFunctions(NodeUtil_GetFunctionBody(node), &changed)
	return changed
}

// Java static method: InlineFunctions.findCalledFunctions(com.google.javascript.rhino.Node node, java.util.Set<java.lang.String> changed)
func ordSet_InlineFunctions_FindCalledFunctions(node pkg_rhino.Node, changed *j2gort.OrderedSet[string]) {
	clinit_InlineFunctions()
	_ = node
	_ = changed
	if !(changed != nil) {
		j2gort.CheckArgument(false)
	}
	if node.IsName() && InlineFunctions_IsCandidateUsage(node) {
		changed.Add(node.GetString())
	}
	{
		var c pkg_rhino.Node = node.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			ordSet_InlineFunctions_FindCalledFunctions(c, changed)
		}
	}
}

// Java method: InlineFunctions.decomposeExpressions()
func (this_ *InlineFunctions) priv_InlineFunctions_decomposeExpressions_3282() {
	{
		cur_8 := this_.fns.IterValues()
		for cur_8.HasNext() {

			functionState := cur_8.
				Next()
			_ = functionState
			if functionState.CanInline() {
				if functionState.HasInlineReferences() {
					var r0 *InlineFunctions_Reference = functionState.InlineRef0()
					_ = r0
					if r0 != nil {
						if r0.requiresDecomposition {
							this_.injector.MaybePrepareCall(ToIface_FunctionInjector_Reference(r0))
						}
						var r1 *InlineFunctions_Reference = functionState.InlineRef1()
						_ = r1
						if (r1 != nil) && r1.requiresDecomposition {
							this_.injector.MaybePrepareCall(ToIface_FunctionInjector_Reference(r1))
						}
					}
				} else {
					{
						cur_10 := j2gort.NewIterCursor(j2gort.CheckNotNull(functionState.GetReferences()))
						for cur_10.HasNext() {
							rawElem_11 := cur_10.Next()
							ref := j2gort.CastPtr[InlineFunctions_Reference](rawElem_11)
							_ = ref
							if ref.requiresDecomposition {
								this_.injector.MaybePrepareCall(ToIface_FunctionInjector_Reference(ref))
							}
						}
					}
				}
			}
		}
	}
}

// Java method: InlineFunctions.removeInlinedFunctions()
func (this_ *InlineFunctions) priv_InlineFunctions_removeInlinedFunctions_3283() {
	{
		cur_12 := this_.fns.IterEntries()
		for cur_12.NextEntry() {
			var name string = cur_12.Key()
			_ = name
			var functionState *InlineFunctions_FunctionState = cur_12.Value()
			_ = functionState
			if functionState.CanRemove() {
				var fn InlineFunctions_Function = functionState.GetFn()
				_ = fn
				if !(functionState.CanInline()) {
					j2gort.CheckState(false)
				}
				if !(fn != nil && !fn.IsJavaNil_()) {
					j2gort.CheckState(false)
				}
				this_.priv_InlineFunctions_verifyAllReferencesInlined_3284(name, functionState)
				fn.Remove()
				NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(fn.GetFunctionNode(), this_.compiler)
			}
		}
	}
}

// Java method: InlineFunctions.verifyAllReferencesInlined(java.lang.String name, com.google.javascript.jscomp.InlineFunctions.FunctionState functionState)
func (this_ *InlineFunctions) priv_InlineFunctions_verifyAllReferencesInlined_3284(name string, functionState *InlineFunctions_FunctionState) {
	_ = name
	_ = functionState
	if functionState.HasInlineReferences() {
		var r0 *InlineFunctions_Reference = functionState.InlineRef0()
		_ = r0
		if r0 != nil {
			InlineFunctions_VerifyReferenceInlined(name, r0)
			var r1 *InlineFunctions_Reference = functionState.InlineRef1()
			_ = r1
			if r1 != nil {
				InlineFunctions_VerifyReferenceInlined(name, r1)
			}
		}
		return
	}
	{
		cur_13 := j2gort.NewIterCursor(j2gort.CheckNotNull(functionState.GetReferences()))
		for cur_13.HasNext() {
			rawElem_14 := cur_13.Next()
			ref := j2gort.CastPtr[InlineFunctions_Reference](rawElem_14)
			_ = ref
			InlineFunctions_VerifyReferenceInlined(name, ref)
		}
	}
}

// Java static method: InlineFunctions.verifyReferenceInlined(java.lang.String name, com.google.javascript.jscomp.InlineFunctions.Reference ref)
func InlineFunctions_VerifyReferenceInlined(name string, ref *InlineFunctions_Reference) {
	clinit_InlineFunctions()
	_ = name
	_ = ref
	if !(ref.inlined) {
		var parent pkg_rhino.Node = ref.GetFunctionInjector_ReferenceBase().callNode.GetParent()
		_ = parent
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Call site missed (" + j2gort.StrValString(name) + ").\n call: " + j2gort.StrValString(ref.GetFunctionInjector_ReferenceBase().callNode.ToStringTree()) + "\n parent:  " + j2gort.StrValString(func() string {
			if parent == 0 {
				return "null"
			}
			return parent.ToStringTree()
		}()))))))
	}
}

// Java static method: InlineFunctions.checkOptimizationLevelSupportsRequireInlining(com.google.javascript.jscomp.CompilerOptions options)
func InlineFunctions_CheckOptimizationLevelSupportsRequireInlining(options *CompilerOptions) *InlineFunctions_ShouldRequireInlining {
	clinit_InlineFunctions()
	_ = options
	var problems j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = problems
	if !(options.ShouldOptimize()) {
		problems.Add_1("non-optimized code")
	}
	if !(options.ShouldInlineProperties()) {
		problems.Add_1("missing shouldInlineProperties")
	}
	if !(options.ShouldInlineConstantVars()) {
		problems.Add_1("missing inlineConstantVars")
	}
	if !(options.GetSmartNameRemoval()) {
		problems.Add_1("missing smartNameRemoval")
	}
	if !(options.GetInlineFunctionsLevel().IsOn()) {
		problems.Add_1("missing inlineFunctions")
	}
	if !(options.GetInlineFunctionsLevel().IncludesGlobals()) {
		problems.Add_1("inlineFunctions does not include globals")
	}
	if options.GetInstrumentForCoverageOption() != StaticGet_CompilerOptions_InstrumentOption_NONE() {
		problems.Add_1("instrumentForCoverageOption enabled")
	}
	{
		var iifeCond_1 bool
		{

			re_15 := options.GetPropertyCollapseLevel()
			ae_16 := StaticGet_CompilerOptions_PropertyCollapseLevel_ALL()
			if re_15 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_15 == ae_16
		}
		if !(iifeCond_1) {
			problems.Add_1("missing property collapsing")
		}
	}
	if !(options.ShouldRemoveClosureAsserts()) {
		problems.Add_1("missing removeClosureAsserts")
	}
	if !(options.ShouldDevirtualizeMethods()) {
		problems.Add_1("missing devirtualizeMethods")
	}
	if options.GetMaxOptimizationLoopIterations() == int32(1) {
		problems.Add_1("optimizationLoopMaxIterations should be >= 2, but is: " + j2gort.StrValInt32(options.GetMaxOptimizationLoopIterations()))
	}
	if problems.IsEmpty() {
		return NewInlineFunctions_ShouldRequireInlining(true, "")
	}
	return NewInlineFunctions_ShouldRequireInlining(false, j2gort.Ext_java_lang_String_Format_2("[%s]", []any{j2gort.BoxString(j2gort.Ext_java_lang_String_Join_2_any_Ext_java_lang_Iterable(", ", problems))}))
}

// Java class: com.google.javascript.jscomp.InlineFunctions.FindCandidateFunctions
type InlineFunctions_FindCandidateFunctions struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_    *InlineFunctions
	callsSeen int32 `gson:"callsSeen"`
}

func (this_ *InlineFunctions_FindCandidateFunctions) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineFunctions_FindCandidateFunctions) IsJavaNil_() bool { return this_ == nil }
func (*InlineFunctions_FindCandidateFunctions) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineFunctions$FindCandidateFunctions", "FindCandidateFunctions"
}
func (this_ *InlineFunctions_FindCandidateFunctions) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*InlineFunctions_FindCandidateFunctions) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineFunctions.FindCandidateFunctions", "com.google.javascript.jscomp.InlineFunctions$FindCandidateFunctions", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineFunctions_FindCandidateFunctions](v_)
	})
}
func (this_ *InlineFunctions_FindCandidateFunctions) SetOuter_(o *InlineFunctions) { this_.outer_ = o }
func (this_ *InlineFunctions_FindCandidateFunctions) GetOuter_() *InlineFunctions {
	return this_.outer_
}
func (*InlineFunctions_FindCandidateFunctions) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*InlineFunctions_FindCandidateFunctions) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *InlineFunctions_FindCandidateFunctions) initFields_InlineFunctions_FindCandidateFunctions() {
	_ = this_

}

// Java constructor: FindCandidateFunctions()
func NewInlineFunctions_FindCandidateFunctions(outer_ *InlineFunctions) *InlineFunctions_FindCandidateFunctions {
	this_ := &InlineFunctions_FindCandidateFunctions{}
	this_.outer_ = outer_
	InitInlineFunctions_FindCandidateFunctions(this_)
	return this_
}

func InitInlineFunctions_FindCandidateFunctions(this_ *InlineFunctions_FindCandidateFunctions) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_InlineFunctions_FindCandidateFunctions()
}

// Java method: FindCandidateFunctions.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InlineFunctions_FindCandidateFunctions) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if this_.outer_.reach.IncludesGlobals() || !(t.InGlobalHoistScope()) {
		this_.FindNamedFunctions(t, n, parent)
		this_.FindFunctionExpressions(t, n)
	}
}

// Java method: FindCandidateFunctions.findNamedFunctions(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InlineFunctions_FindCandidateFunctions) FindNamedFunctions(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(NodeUtil_IsStatement(n)) {
		return
	}
	switch n.TokenSwitchOrd_() {
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		if !(n.HasOneChild()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
		}
		var nameNode pkg_rhino.Node = n.GetFirstChild()
		_ = nameNode
		if (nameNode.IsName() && nameNode.HasChildren()) && nameNode.GetFirstChild().IsFunction() {
			this_.outer_.priv_InlineFunctions_maybeAddFunction_3271(ToIface_InlineFunctions_Function(allocPooled_InlineFunctions_FunctionVar(&this_.outer_.slab_InlineFunctions_FunctionVar, this_.outer_, n)), t.GetChunk())
		}
	case int32(73) /* Token.FUNCTION */ :
		if !(NodeUtil_IsStatementBlock(parent) || parent.IsLabel()) {
			j2gort.CheckState(false)
		}
		if NodeUtil_IsFunctionDeclaration(n) {
			var fn InlineFunctions_Function = ToIface_InlineFunctions_Function(allocPooled_InlineFunctions_NamedFunction(&this_.outer_.slab_InlineFunctions_NamedFunction, this_.outer_, n))
			_ = fn
			this_.outer_.priv_InlineFunctions_maybeAddFunction_3271(fn, t.GetChunk())
		}
	default:
	}
}

// Java method: FindCandidateFunctions.findFunctionExpressions(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *InlineFunctions_FindCandidateFunctions) FindFunctionExpressions(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(31) /* Token.OPTCHAIN_CALL */, int32(28) /* Token.CALL */ :
		var fnNode pkg_rhino.Node = pkg_rhino.Node(0)
		_ = fnNode
		if n.GetFirstChild().IsFunction() {
			fnNode = n.GetFirstChild()
		} else if NodeUtil_IsFunctionObjectCall(n) {
			var fnIdentifyingNode pkg_rhino.Node = n.GetFirstFirstChild()
			_ = fnIdentifyingNode
			if fnIdentifyingNode.IsFunction() {
				fnNode = fnIdentifyingNode
			}
		}
		if fnNode != 0 {
			var fn InlineFunctions_Function = ToIface_InlineFunctions_Function(allocPooled_InlineFunctions_FunctionExpression(&this_.outer_.slab_InlineFunctions_FunctionExpression, fnNode, func() int32 { old_17 := this_.callsSeen; nw_18 := old_17 + 1; this_.callsSeen = nw_18; return old_17 }()))
			_ = fn
			this_.outer_.priv_InlineFunctions_maybeAddFunction_3271(fn, t.GetChunk())
			this_.outer_.anonFns.Put(fnNode, fn.GetName())
		}
	default:
	}
}

// Java interface: com.google.javascript.jscomp.InlineFunctions.CallVisitorCallback
type InlineFunctions_CallVisitorCallback interface {
	Implements_jscomp_InlineFunctions_CallVisitorCallback()
	IsJavaNil_() bool
	VisitCallSite(*NodeTraversal, pkg_rhino.Node, *InlineFunctions_FunctionState)
}

func ToIface_InlineFunctions_CallVisitorCallback[T interface {
	InlineFunctions_CallVisitorCallback
	comparable
}](v T) InlineFunctions_CallVisitorCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_InlineFunctions_CallVisitorCallback[T InlineFunctions_CallVisitorCallback](v T) InlineFunctions_CallVisitorCallback {
	var vi InlineFunctions_CallVisitorCallback = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineFunctions.CallVisitorCallback", "com.google.javascript.jscomp.InlineFunctions$CallVisitorCallback", func(v_ any) bool {
		return j2gort.InstanceOf[InlineFunctions_CallVisitorCallback](v_)
	})
}

// Java polymorphic class interface: com.google.javascript.jscomp.InlineFunctions.CallVisitor
type InlineFunctions_CallVisitor interface {
	Implements_jscomp_InlineFunctions_CallVisitor()
	IsJavaNil_() bool
	GetInlineFunctions_CallVisitorBase() *InlineFunctions_CallVisitorBase
	NodeTraversal_AbstractPostOrderCallback
}

func ToIface_InlineFunctions_CallVisitor[T interface {
	InlineFunctions_CallVisitor
	comparable
}](v T) InlineFunctions_CallVisitor {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_InlineFunctions_CallVisitor[T InlineFunctions_CallVisitor](v T) InlineFunctions_CallVisitor {
	var vi InlineFunctions_CallVisitor = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.InlineFunctions.CallVisitor
type InlineFunctions_CallVisitorBase struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	callback InlineFunctions_CallVisitorCallback `gson:"callback"`
	owner    *InlineFunctions                    `gson:"owner"`
}

func (this_ *InlineFunctions_CallVisitorBase) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineFunctions_CallVisitorBase) IsJavaNil_() bool { return this_ == nil }
func (*InlineFunctions_CallVisitorBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineFunctions$CallVisitor", "CallVisitor"
}
func (this_ *InlineFunctions_CallVisitorBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineFunctions.CallVisitor", "com.google.javascript.jscomp.InlineFunctions$CallVisitor", func(v_ any) bool {
		return j2gort.InstanceOf[InlineFunctions_CallVisitor](v_)
	})
}
func (b *InlineFunctions_CallVisitorBase) GetInlineFunctions_CallVisitorBase() *InlineFunctions_CallVisitorBase {
	return b
}
func (*InlineFunctions_CallVisitorBase) Implements_jscomp_InlineFunctions_CallVisitor()             {}
func (*InlineFunctions_CallVisitorBase) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*InlineFunctions_CallVisitorBase) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *InlineFunctions_CallVisitorBase) initFields_InlineFunctions_CallVisitor() {
	_ = this_
}

// Java constructor: CallVisitor(com.google.javascript.jscomp.InlineFunctions owner, com.google.javascript.jscomp.InlineFunctions.@org.jspecify.annotations.Nullable CallVisitorCallback callback)
func NewInlineFunctions_CallVisitor(owner *InlineFunctions, callback InlineFunctions_CallVisitorCallback) InlineFunctions_CallVisitor {
	this_ := &InlineFunctions_CallVisitorBase{}
	InitInlineFunctions_CallVisitor(this_, owner, callback)
	return this_
}

func InitInlineFunctions_CallVisitor(this_ *InlineFunctions_CallVisitorBase, owner *InlineFunctions, callback InlineFunctions_CallVisitorCallback) {
	_ = owner
	_ = callback
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_InlineFunctions_CallVisitor()
	this_.owner = owner
	this_.callback = callback
}

// Java method: CallVisitor.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InlineFunctions_CallVisitorBase) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(31) /* Token.OPTCHAIN_CALL */, int32(28) /* Token.CALL */ :
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		var name string = j2gort.NullString
		_ = name
		if child.IsName() {
			name = child.GetString()
		} else if child.IsFunction() {
			name = j2gort.OrderedMapGetStr(this_.owner.anonFns, child)
		} else if NodeUtil_IsFunctionObjectCall(n) {
			if !(NodeUtil_IsNormalOrOptChainGet(child)) {
				j2gort.CheckState(false)
			}
			var fnIdentifyingNode pkg_rhino.Node = child.GetFirstChild()
			_ = fnIdentifyingNode
			if fnIdentifyingNode.IsName() {
				name = fnIdentifyingNode.GetString()
			} else if fnIdentifyingNode.IsFunction() {
				name = j2gort.OrderedMapGetStr(this_.owner.anonFns, fnIdentifyingNode)
			}
		}
		if name != j2gort.NullString {
			var functionState *InlineFunctions_FunctionState = this_.owner.fns.Get(name)
			_ = functionState
			if functionState != nil {
				this_.callback.VisitCallSite(t, n, functionState)
			}
		}
	default:
	}
}

// Java class: com.google.javascript.jscomp.InlineFunctions.FindCandidatesReferences
type InlineFunctions_FindCandidatesReferences struct {
	InlineFunctions_CallVisitorBase
	outer_ *InlineFunctions
}

func (this_ *InlineFunctions_FindCandidatesReferences) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineFunctions_FindCandidatesReferences) IsJavaNil_() bool { return this_ == nil }
func (*InlineFunctions_FindCandidatesReferences) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineFunctions$FindCandidatesReferences", "FindCandidatesReferences"
}
func (this_ *InlineFunctions_FindCandidatesReferences) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.InlineFunctions_CallVisitorBase.IdentityHash_()
}
func (*InlineFunctions_FindCandidatesReferences) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineFunctions.FindCandidatesReferences", "com.google.javascript.jscomp.InlineFunctions$FindCandidatesReferences", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineFunctions_FindCandidatesReferences](v_)
	})
}
func (this_ *InlineFunctions_FindCandidatesReferences) SetOuter_(o *InlineFunctions) {
	this_.outer_ = o
}
func (this_ *InlineFunctions_FindCandidatesReferences) GetOuter_() *InlineFunctions {
	return this_.outer_
}
func (*InlineFunctions_FindCandidatesReferences) Implements_jscomp_InlineFunctions_CallVisitor() {}
func (*InlineFunctions_FindCandidatesReferences) Implements_jscomp_InlineFunctions_CallVisitorCallback() {
}
func (*InlineFunctions_FindCandidatesReferences) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*InlineFunctions_FindCandidatesReferences) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *InlineFunctions_FindCandidatesReferences) initFields_InlineFunctions_FindCandidatesReferences() {
	_ = this_
}

// Java constructor: FindCandidatesReferences()
func NewInlineFunctions_FindCandidatesReferences(outer_ *InlineFunctions) *InlineFunctions_FindCandidatesReferences {
	this_ := &InlineFunctions_FindCandidatesReferences{}
	this_.outer_ = outer_
	InitInlineFunctions_FindCandidatesReferences(this_)
	return this_
}

func InitInlineFunctions_FindCandidatesReferences(this_ *InlineFunctions_FindCandidatesReferences) {
	InitInlineFunctions_CallVisitor(&this_.InlineFunctions_CallVisitorBase, this_.outer_, nil)
	this_.initFields_InlineFunctions_FindCandidatesReferences()
	this_.GetInlineFunctions_CallVisitorBase().callback = ToIface_InlineFunctions_CallVisitorCallback(this_)
}

// Java method: FindCandidatesReferences.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InlineFunctions_FindCandidatesReferences) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	this_.InlineFunctions_CallVisitorBase.Visit_3(t, n, parent)
	if n.IsName() {
		this_.priv_FindCandidatesReferences_checkNameUsage_3290(n, parent)
	}
}

// Java method: FindCandidatesReferences.visitCallSite(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node callNode, com.google.javascript.jscomp.InlineFunctions.FunctionState functionState)
func (this_ *InlineFunctions_FindCandidatesReferences) VisitCallSite(t *NodeTraversal, callNode pkg_rhino.Node, functionState *InlineFunctions_FunctionState) {
	_ = t
	_ = callNode
	_ = functionState
	this_.MaybeAddReference(t, functionState, callNode, t.GetChunk())
}

// Java method: FindCandidatesReferences.maybeAddReference(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.InlineFunctions.FunctionState functionState, com.google.javascript.rhino.Node callNode, com.google.javascript.jscomp.JSChunk chunk)
func (this_ *InlineFunctions_FindCandidatesReferences) MaybeAddReference(t *NodeTraversal, functionState *InlineFunctions_FunctionState, callNode pkg_rhino.Node, chunk *JSChunk) {
	_ = t
	_ = functionState
	_ = callNode
	_ = chunk
	if !(functionState.CanInline()) {
		return
	}
	var mode FunctionInjector_InliningMode
	_ = mode
	if functionState.CanInlineDirectly() {
		mode = StaticGet_FunctionInjector_InliningMode_DIRECT()
	} else {
		mode = StaticGet_FunctionInjector_InliningMode_BLOCK()
	}
	var referenceAdded bool = this_.priv_FindCandidatesReferences_maybeAddReferenceUsingMode_3289(t, functionState, callNode, chunk, mode)
	_ = referenceAdded
	if !(referenceAdded) && (mode == StaticGet_FunctionInjector_InliningMode_DIRECT()) {
		mode = StaticGet_FunctionInjector_InliningMode_BLOCK()
		referenceAdded = this_.priv_FindCandidatesReferences_maybeAddReferenceUsingMode_3289(t, functionState, callNode, chunk, mode)
	}
	if !(referenceAdded) {
		functionState.SetRemove_5(false, this_.outer_.compiler, StaticGet_InlineFunctions_CannotRemoveReason_UNINLINABLE_REFERENCE(), StaticGet_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO(), callNode)
	}
}

// Java method: FindCandidatesReferences.maybeAddReferenceUsingMode(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.InlineFunctions.FunctionState functionState, com.google.javascript.rhino.Node callNode, com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.FunctionInjector.InliningMode mode)
func (this_ *InlineFunctions_FindCandidatesReferences) priv_FindCandidatesReferences_maybeAddReferenceUsingMode_3289(t *NodeTraversal, functionState *InlineFunctions_FunctionState, callNode pkg_rhino.Node, chunk *JSChunk, mode FunctionInjector_InliningMode) bool {
	_ = t
	_ = functionState
	_ = callNode
	_ = chunk
	_ = mode
	if this_.outer_.enforceMaxSizeAfterInlining && this_.outer_.priv_InlineFunctions_targetSizeAfterInlineExceedsLimit_3270(t, functionState) {
		return false
	}
	var candidate *InlineFunctions_Reference = allocPooled_InlineFunctions_Reference(&this_.outer_.slab_InlineFunctions_Reference, callNode, t.GetScope(), chunk, mode)
	_ = candidate
	var result FunctionInjector_CanInlineResult = this_.outer_.injector.CanInlineReferenceToFunction(ToIface_FunctionInjector_Reference(candidate), functionState.GetFn().GetFunctionNode(), functionState.GetNamesToAlias(), functionState.GetReferencesThis(), functionState.HasInnerClosures())
	_ = result
	if result != StaticGet_FunctionInjector_CanInlineResult_NO() {
		candidate.SetRequiresDecomposition((result == StaticGet_FunctionInjector_CanInlineResult_AFTER_PREPARATION()))
		functionState.AddReference(candidate)
		return true
	}
	return false
}

// Java method: FindCandidatesReferences.checkNameUsage(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InlineFunctions_FindCandidatesReferences) priv_FindCandidatesReferences_checkNameUsage_3290(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	if !(n.IsName()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	if InlineFunctions_IsCandidateUsage(n) {
		return
	}
	var name string = n.GetString()
	_ = name
	var functionState *InlineFunctions_FunctionState = this_.GetInlineFunctions_CallVisitorBase().owner.fns.Get(name)
	_ = functionState
	if functionState == nil {
		return
	}
	if NodeUtil_IsLValue(n) {
		functionState.DisallowInlining_2(this_.outer_.compiler, StaticGet_InlineFunctions_DisallowInliningReason_REASSIGNED())
	} else if (((parent.IsAssign() && (parent.GetJSDocInfo() != nil)) && parent.GetJSDocInfo().IsConstant()) && (parent.GetSecondChild() == n)) || (((parent.GetParent() != 0) && parent.GetParent().IsConst()) && (parent.GetFirstChild() == n)) {
		functionState.SetRemove_4(false, this_.outer_.compiler, StaticGet_InlineFunctions_CannotRemoveReason_UNREMOVABLE_REFERENCE(), StaticGet_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO())
	} else {
		functionState.SetRemove_5(false, this_.outer_.compiler, StaticGet_InlineFunctions_CannotRemoveReason_UNREMOVABLE_REFERENCE(), StaticGet_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES(), parent)
	}
}

// Java class: com.google.javascript.jscomp.InlineFunctions.Inline
type InlineFunctions_Inline struct {
	injector *FunctionInjector `gson:"injector"`
	idHash_  atomic.Int32
}

func (this_ *InlineFunctions_Inline) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineFunctions_Inline) IsJavaNil_() bool { return this_ == nil }
func (*InlineFunctions_Inline) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineFunctions$Inline", "Inline"
}
func (this_ *InlineFunctions_Inline) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineFunctions_Inline) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineFunctions.Inline", "com.google.javascript.jscomp.InlineFunctions$Inline", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineFunctions_Inline](v_)
	})
}
func (*InlineFunctions_Inline) Implements_jscomp_InlineFunctions_CallVisitorCallback() {}

func (this_ *InlineFunctions_Inline) initFields_InlineFunctions_Inline() {
	_ = this_
}

// Java constructor: Inline(com.google.javascript.jscomp.FunctionInjector injector)
func NewInlineFunctions_Inline(injector *FunctionInjector) *InlineFunctions_Inline {
	this_ := &InlineFunctions_Inline{}
	InitInlineFunctions_Inline(this_, injector)
	return this_
}

func InitInlineFunctions_Inline(this_ *InlineFunctions_Inline, injector *FunctionInjector) {
	_ = injector
	this_.initFields_InlineFunctions_Inline()
	this_.injector = injector
}

// Java method: Inline.visitCallSite(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node callNode, com.google.javascript.jscomp.InlineFunctions.FunctionState functionState)
func (this_ *InlineFunctions_Inline) VisitCallSite(t *NodeTraversal, callNode pkg_rhino.Node, functionState *InlineFunctions_FunctionState) {
	_ = t
	_ = callNode
	_ = functionState
	if !(functionState.HasExistingFunctionDefinition()) {
		j2gort.CheckState(false)
	}
	if functionState.CanInline() {
		var ref *InlineFunctions_Reference = functionState.GetReference(callNode)
		_ = ref
		if ref != nil {
			this_.priv_Inline_inlineFunction_3291(t, ref, functionState)
			{
				r_19 := ref
				r_19.inlined = true
			}
		}
	}
}

// Java method: Inline.inlineFunction(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.InlineFunctions.Reference ref, com.google.javascript.jscomp.InlineFunctions.FunctionState functionState)
func (this_ *InlineFunctions_Inline) priv_Inline_inlineFunction_3291(t *NodeTraversal, ref *InlineFunctions_Reference, functionState *InlineFunctions_FunctionState) {
	_ = t
	_ = ref
	_ = functionState
	var fn InlineFunctions_Function = functionState.GetFn()
	_ = fn
	var fnName string = fn.GetName()
	_ = fnName
	var fnNode pkg_rhino.Node = functionState.GetSafeFnNode()
	_ = fnNode
	var newExpr pkg_rhino.Node = this_.injector.Inline(ToIface_FunctionInjector_Reference(ref), fnName, fnNode)
	_ = newExpr
	{
		var iifeCond_1 bool
		{

			rn_20 := newExpr
			an_21 := ref.GetFunctionInjector_ReferenceBase().callNode
			iifeCond_1 = j2gort.CheckNotNullNode(rn_20) == an_21
		}
		if !(iifeCond_1) {
			t.GetCompiler().ReportChangeToEnclosingScope(newExpr)
		}
	}
}

// Java class: com.google.javascript.jscomp.InlineFunctions.FunctionState
type InlineFunctions_FunctionState struct {
	fn                                        InlineFunctions_Function                                       `gson:"fn"`
	ref0                                      *InlineFunctions_Reference                                     `gson:"ref0"`
	ref1                                      *InlineFunctions_Reference                                     `gson:"ref1"`
	references                                *j2gort.OrderedMap[pkg_rhino.Node, *InlineFunctions_Reference] `gson:"references"`
	chunk                                     *JSChunk                                                       `gson:"chunk"`
	namesToAlias                              j2gort.Ext_java_util_Set                                       `gson:"namesToAlias"`
	safeFnNode                                pkg_rhino.Node                                                 `gson:"safeFnNode"`
	idHash_                                   atomic.Int32
	inline                                    bool `gson:"inline"`
	remove                                    bool `gson:"remove"`
	inlineDirectly                            bool `gson:"inlineDirectly"`
	referencesThis                            bool `gson:"referencesThis"`
	hasInnerClosures                          bool `gson:"hasInnerClosures"`
	hasRequireInliningAnnotation              bool `gson:"hasRequireInliningAnnotation"`
	hasRequireInliningAnnotationInitialized   bool `gson:"hasRequireInliningAnnotationInitialized"`
	hasEncourageInliningAnnotation            bool `gson:"hasEncourageInliningAnnotation"`
	hasEncourageInliningAnnotationInitialized bool `gson:"hasEncourageInliningAnnotationInitialized"`
}

func (this_ *InlineFunctions_FunctionState) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineFunctions_FunctionState) IsJavaNil_() bool { return this_ == nil }
func (*InlineFunctions_FunctionState) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineFunctions$FunctionState", "FunctionState"
}
func (this_ *InlineFunctions_FunctionState) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineFunctions_FunctionState) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineFunctions.FunctionState", "com.google.javascript.jscomp.InlineFunctions$FunctionState", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineFunctions_FunctionState](v_)
	})
}

func (this_ *InlineFunctions_FunctionState) initFields_InlineFunctions_FunctionState() {
	_ = this_

	this_.inline = true
	this_.remove = true

	if this_.references != nil {
		j2gort.ReleaseOrderedMap(this_.references)
		this_.references = nil
	}

	this_.chunk = nil
	this_.namesToAlias = nil
	this_.hasRequireInliningAnnotation = false
	this_.hasRequireInliningAnnotationInitialized = false
	this_.hasEncourageInliningAnnotation = false
	this_.hasEncourageInliningAnnotationInitialized = false
}

// Java constructor: FunctionState()
func NewInlineFunctions_FunctionState() *InlineFunctions_FunctionState {
	this_ := &InlineFunctions_FunctionState{}
	InitInlineFunctions_FunctionState(this_)
	return this_
}

func InitInlineFunctions_FunctionState(this_ *InlineFunctions_FunctionState) {
	this_.initFields_InlineFunctions_FunctionState()
}

// Java method: FunctionState.hasExistingFunctionDefinition()
func (this_ *InlineFunctions_FunctionState) HasExistingFunctionDefinition() bool {
	return (this_.fn != nil && !this_.fn.IsJavaNil_())
}

// Java method: FunctionState.setReferencesThis(boolean referencesThis)
func (this_ *InlineFunctions_FunctionState) SetReferencesThis(referencesThis bool) {
	_ = referencesThis
	this_.referencesThis = referencesThis
}

// Java method: FunctionState.getReferencesThis()
func (this_ *InlineFunctions_FunctionState) GetReferencesThis() bool {
	return this_.referencesThis
}

// Java method: FunctionState.setHasInnerClosures(boolean hasInnerClosures)
func (this_ *InlineFunctions_FunctionState) SetHasInnerClosures(hasInnerClosures bool) {
	_ = hasInnerClosures
	this_.hasInnerClosures = hasInnerClosures
}

// Java method: FunctionState.hasInnerClosures()
func (this_ *InlineFunctions_FunctionState) HasInnerClosures() bool {
	return this_.hasInnerClosures
}

// Java method: FunctionState.removeBlockInliningReferences()
func (this_ *InlineFunctions_FunctionState) RemoveBlockInliningReferences() {
	if this_.references != nil {
		var toRemove []pkg_rhino.Node = nil
		_ = toRemove
		{
			cur_22 := this_.references.IterEntries()
			for cur_22.NextEntry() {
				if cur_22.Value().GetFunctionInjector_ReferenceBase().mode == StaticGet_FunctionInjector_InliningMode_BLOCK() {
					if toRemove == nil {
						toRemove = []pkg_rhino.Node{}
					}
					toRemove = append(toRemove, cur_22.Key())
				}
			}
		}
		if toRemove != nil {
			{
				var i int32 = int32(0)
				_ = i
				var size int32 = int32(len(toRemove))
				_ = size
				for ; i < size; i = i + 1 {
					this_.references.Remove(toRemove[i])
				}
			}
		}
		return
	}
	if (this_.ref1 != nil) && (this_.ref1.GetFunctionInjector_ReferenceBase().mode == StaticGet_FunctionInjector_InliningMode_BLOCK()) {
		this_.ref1 = nil
	}
	if (this_.ref0 != nil) && (this_.ref0.GetFunctionInjector_ReferenceBase().mode == StaticGet_FunctionInjector_InliningMode_BLOCK()) {
		this_.ref0 = this_.ref1
		this_.ref1 = nil
	}
}

// Java method: FunctionState.hasBlockInliningReferences()
func (this_ *InlineFunctions_FunctionState) HasBlockInliningReferences() bool {
	if this_.references != nil {
		{
			cur_23 := this_.references.IterValues()
			for cur_23.HasNext() {

				r := cur_23.
					Next()
				_ = r
				if r.GetFunctionInjector_ReferenceBase().mode == StaticGet_FunctionInjector_InliningMode_BLOCK() {
					return true
				}
			}
		}
		return false
	}
	return (((this_.ref0 != nil) && (this_.ref0.GetFunctionInjector_ReferenceBase().mode == StaticGet_FunctionInjector_InliningMode_BLOCK())) || ((this_.ref1 != nil) && (this_.ref1.GetFunctionInjector_ReferenceBase().mode == StaticGet_FunctionInjector_InliningMode_BLOCK())))
}

// Java method: FunctionState.getFn()
func (this_ *InlineFunctions_FunctionState) GetFn() InlineFunctions_Function {
	return this_.fn
}

// Java method: FunctionState.setFn(com.google.javascript.jscomp.InlineFunctions.Function fn)
func (this_ *InlineFunctions_FunctionState) SetFn(fn InlineFunctions_Function) {
	_ = fn
	if this_.fn != nil && !this_.fn.IsJavaNil_() {
		j2gort.CheckState(false)
	}
	this_.fn = fn
}

// Java method: FunctionState.getSafeFnNode()
func (this_ *InlineFunctions_FunctionState) GetSafeFnNode() pkg_rhino.Node {
	if this_.safeFnNode != 0 {
		return this_.safeFnNode
	} else {
		return this_.fn.GetFunctionNode()
	}
}

// Java method: FunctionState.setSafeFnNode(com.google.javascript.rhino.Node safeFnNode)
func (this_ *InlineFunctions_FunctionState) SetSafeFnNode(safeFnNode pkg_rhino.Node) {
	_ = safeFnNode
	this_.safeFnNode = safeFnNode
}

// Java method: FunctionState.canInline()
func (this_ *InlineFunctions_FunctionState) CanInline() bool {
	return this_.inline
}

// Java method: FunctionState.encourageInlining()
func (this_ *InlineFunctions_FunctionState) EncourageInlining() bool {
	if this_.fn == nil || this_.fn.IsJavaNil_() {
		return false
	}
	if !(this_.hasEncourageInliningAnnotationInitialized) {
		this_.hasEncourageInliningAnnotationInitialized = true
		this_.hasEncourageInliningAnnotation = InlineFunctions_HasEncourageInliningAnnotation(this_.fn.GetFunctionNode())
	}
	return this_.hasEncourageInliningAnnotation
}

// Java method: FunctionState.requireInlining()
func (this_ *InlineFunctions_FunctionState) RequireInlining() bool {
	if this_.fn == nil || this_.fn.IsJavaNil_() {
		return false
	}
	if !(this_.hasRequireInliningAnnotationInitialized) {
		this_.hasRequireInliningAnnotationInitialized = true
		this_.hasRequireInliningAnnotation = InlineFunctions_HasRequireInliningAnnotation(this_.fn.GetFunctionNode())
	}
	return this_.hasRequireInliningAnnotation
}

// Java method: FunctionState.disallowInlining(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.InlineFunctions.DisallowInliningReason reason)
func (this_ *InlineFunctions_FunctionState) DisallowInlining_2(compiler AbstractCompiler, reason InlineFunctions_DisallowInliningReason) {
	_ = compiler
	_ = reason
	this_.DisallowInlining_3(compiler, reason, StaticGet_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES())
}

// Java method: FunctionState.disallowInlining(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.InlineFunctions.DisallowInliningReason reason, com.google.javascript.jscomp.InlineFunctions.ShouldWarnWhenRequireInliningCannotInline shouldWarnWhenRequireInliningCannotInline)
func (this_ *InlineFunctions_FunctionState) DisallowInlining_3(compiler AbstractCompiler, reason InlineFunctions_DisallowInliningReason, shouldWarnWhenRequireInliningCannotInline InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) {
	_ = compiler
	_ = reason
	_ = shouldWarnWhenRequireInliningCannotInline
	if (this_.RequireInlining() && compiler.IsDebugLoggingEnabled()) && (shouldWarnWhenRequireInliningCannotInline == StaticGet_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES()) {
		compiler.Report_1(JSError_Make_3(this_.fn.GetFunctionNode(), StaticGet_InlineFunctions_MISSED_REQUIRED_INLINING(), []string{(reason).Name()}))
	}
	this_.inline = false
	this_.ref0 = nil
	this_.ref1 = nil
	if this_.references != nil {
		j2gort.ReleaseOrderedMap(this_.references)
		this_.references = nil
	}

	this_.remove = false
}

// Java method: FunctionState.canRemove()
func (this_ *InlineFunctions_FunctionState) CanRemove() bool {
	return this_.remove
}

// Java method: FunctionState.setRemove(boolean remove, com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.InlineFunctions.CannotRemoveReason reason, com.google.javascript.jscomp.InlineFunctions.ShouldWarnWhenRequireInliningCannotInline shouldWarnIfRequireInlining)
func (this_ *InlineFunctions_FunctionState) SetRemove_4(remove bool, compiler AbstractCompiler, reason InlineFunctions_CannotRemoveReason, shouldWarnIfRequireInlining InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) {
	_ = remove
	_ = compiler
	_ = reason
	_ = shouldWarnIfRequireInlining
	this_.SetRemove_5(remove, compiler, reason, shouldWarnIfRequireInlining, pkg_rhino.Node(0))
}

// Java method: FunctionState.setRemove(boolean remove, com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.InlineFunctions.CannotRemoveReason reason, com.google.javascript.jscomp.InlineFunctions.ShouldWarnWhenRequireInliningCannotInline shouldWarnIfRequireInlining, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node contextNode)
func (this_ *InlineFunctions_FunctionState) SetRemove_5(remove bool, compiler AbstractCompiler, reason InlineFunctions_CannotRemoveReason, shouldWarnIfRequireInlining InlineFunctions_ShouldWarnWhenRequireInliningCannotInline, contextNode pkg_rhino.Node) {
	_ = remove
	_ = compiler
	_ = reason
	_ = shouldWarnIfRequireInlining
	_ = contextNode
	if ((!(remove) && (shouldWarnIfRequireInlining != StaticGet_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO())) && this_.RequireInlining()) && compiler.IsDebugLoggingEnabled() {
		compiler.Report_1(JSError_Make_3(this_.fn.GetFunctionNode(), StaticGet_InlineFunctions_MISSED_REQUIRED_INLINING(), []string{(reason).Name(), func() string {
			if contextNode == 0 {
				return "unknown"
			}
			return contextNode.String()
		}()}))
	}
	this_.remove = remove
}

// Java method: FunctionState.canInlineDirectly()
func (this_ *InlineFunctions_FunctionState) CanInlineDirectly() bool {
	return this_.inlineDirectly
}

// Java method: FunctionState.inlineDirectly(boolean directReplacement)
func (this_ *InlineFunctions_FunctionState) InlineDirectly(directReplacement bool) {
	_ = directReplacement
	this_.inlineDirectly = directReplacement
}

// Java method: FunctionState.hasReferences()
func (this_ *InlineFunctions_FunctionState) HasReferences() bool {
	return ((this_.ref0 != nil) || ((this_.references != nil) && !(this_.references.IsEmpty())))
}

// Java method: FunctionState.addReference(com.google.javascript.jscomp.InlineFunctions.Reference ref)
func (this_ *InlineFunctions_FunctionState) AddReference(ref *InlineFunctions_Reference) {
	_ = ref
	if this_.references != nil {
		this_.references.Put(ref.GetFunctionInjector_ReferenceBase().callNode, ref)
		return
	}
	if (this_.ref0 == nil) || (this_.ref0.GetFunctionInjector_ReferenceBase().callNode == ref.GetFunctionInjector_ReferenceBase().callNode) {
		this_.ref0 = ref
		return
	}
	if (this_.ref1 == nil) || (this_.ref1.GetFunctionInjector_ReferenceBase().callNode == ref.GetFunctionInjector_ReferenceBase().callNode) {
		this_.ref1 = ref
		return
	}
	var r0 *InlineFunctions_Reference = this_.ref0
	_ = r0
	var r1 *InlineFunctions_Reference = this_.ref1
	_ = r1
	this_.ref0 = nil
	this_.ref1 = nil
	this_.references = j2gort.AcquireOrderedMapPtrCap[pkg_rhino.Node, *InlineFunctions_Reference](int32(8))
	this_.references.Put(r0.GetFunctionInjector_ReferenceBase().callNode, r0)
	this_.references.Put(r1.GetFunctionInjector_ReferenceBase().callNode, r1)
	this_.references.Put(ref.GetFunctionInjector_ReferenceBase().callNode, ref)
}

// Java method: FunctionState.hasInlineReferences()
func (this_ *InlineFunctions_FunctionState) HasInlineReferences() bool {
	return (this_.references == nil)
}

// Java method: FunctionState.inlineRef0()
func (this_ *InlineFunctions_FunctionState) InlineRef0() *InlineFunctions_Reference {
	return this_.ref0
}

// Java method: FunctionState.inlineRef1()
func (this_ *InlineFunctions_FunctionState) InlineRef1() *InlineFunctions_Reference {
	return this_.ref1
}

// Java method: FunctionState.getReferences()
func (this_ *InlineFunctions_FunctionState) GetReferences() j2gort.Ext_java_util_Collection {
	if this_.references != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.references.LiveValues())
	}
	if this_.ref0 == nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
	if this_.ref1 != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.PtrToAny(this_.ref0), j2gort.PtrToAny(this_.ref1)))
	} else {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.ImmListOfPtr(this_.ref0))
	}
}

// Java method: FunctionState.getReference(com.google.javascript.rhino.Node n)
func (this_ *InlineFunctions_FunctionState) GetReference(n pkg_rhino.Node) *InlineFunctions_Reference {
	_ = n
	if this_.references != nil {
		return this_.references.Get(n)
	}
	if this_.ref0 != nil {
		if this_.ref0.GetFunctionInjector_ReferenceBase().callNode == n {
			return this_.ref0
		}
		if (this_.ref1 != nil) && (this_.ref1.GetFunctionInjector_ReferenceBase().callNode == n) {
			return this_.ref1
		}
	}
	return nil
}

// Java method: FunctionState.getNamesToAliasReadOnly()
func (this_ *InlineFunctions_FunctionState) GetNamesToAliasReadOnly() j2gort.Ext_java_util_Set {
	if this_.namesToAlias != nil {
		return this_.namesToAlias
	} else {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
	}
}

// Java method: FunctionState.getNamesToAlias()
func (this_ *InlineFunctions_FunctionState) GetNamesToAlias() j2gort.Ext_com_google_common_collect_ImmutableSet {
	if this_.namesToAlias == nil {
		return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	}
	return j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.namesToAlias))
}

// Java method: FunctionState.setNamesToAlias(java.util.Set<java.lang.String> names)
func (this_ *InlineFunctions_FunctionState) SetNamesToAlias(names j2gort.Ext_java_util_Set) {
	_ = names
	this_.namesToAlias = names
}

// Java method: FunctionState.setChunk(com.google.javascript.jscomp.JSChunk chunk)
func (this_ *InlineFunctions_FunctionState) SetChunk(chunk *JSChunk) {
	_ = chunk
	this_.chunk = chunk
}

// Java method: FunctionState.getChunk()
func (this_ *InlineFunctions_FunctionState) GetChunk() *JSChunk {
	return this_.chunk
}

// Java interface: com.google.javascript.jscomp.InlineFunctions.Function
type InlineFunctions_Function interface {
	Implements_jscomp_InlineFunctions_Function()
	IsJavaNil_() bool
	GetName() string
	GetNameNode() pkg_rhino.Node
	GetFunctionNode() pkg_rhino.Node
	Remove()
	GetDeclaringBlock() pkg_rhino.Node
}

func ToIface_InlineFunctions_Function[T interface {
	InlineFunctions_Function
	comparable
}](v T) InlineFunctions_Function {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_InlineFunctions_Function[T InlineFunctions_Function](v T) InlineFunctions_Function {
	var vi InlineFunctions_Function = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineFunctions.Function", "com.google.javascript.jscomp.InlineFunctions$Function", func(v_ any) bool {
		return j2gort.InstanceOf[InlineFunctions_Function](v_)
	})
}

// Java class: com.google.javascript.jscomp.InlineFunctions.NamedFunction
type InlineFunctions_NamedFunction struct {
	outer_  *InlineFunctions
	fn      pkg_rhino.Node `gson:"fn"`
	idHash_ atomic.Int32
}

func (this_ *InlineFunctions_NamedFunction) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineFunctions_NamedFunction) IsJavaNil_() bool { return this_ == nil }
func (*InlineFunctions_NamedFunction) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineFunctions$NamedFunction", "NamedFunction"
}
func (this_ *InlineFunctions_NamedFunction) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineFunctions_NamedFunction) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineFunctions.NamedFunction", "com.google.javascript.jscomp.InlineFunctions$NamedFunction", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineFunctions_NamedFunction](v_)
	})
}
func (this_ *InlineFunctions_NamedFunction) SetOuter_(o *InlineFunctions)          { this_.outer_ = o }
func (this_ *InlineFunctions_NamedFunction) GetOuter_() *InlineFunctions           { return this_.outer_ }
func (*InlineFunctions_NamedFunction) Implements_jscomp_InlineFunctions_Function() {}

func (this_ *InlineFunctions_NamedFunction) initFields_InlineFunctions_NamedFunction() {
	_ = this_
}

// Java constructor: NamedFunction(com.google.javascript.rhino.Node fn)
func NewInlineFunctions_NamedFunction(outer_ *InlineFunctions, fn pkg_rhino.Node) *InlineFunctions_NamedFunction {
	this_ := &InlineFunctions_NamedFunction{}
	this_.outer_ = outer_
	InitInlineFunctions_NamedFunction(this_, fn)
	return this_
}
func allocPooled_InlineFunctions_NamedFunction(slab_ *j2gort.PooledSlab[InlineFunctions_NamedFunction], outer_ *InlineFunctions, fn pkg_rhino.Node) *InlineFunctions_NamedFunction {
	this_ := slab_.Alloc()
	this_.outer_ = outer_
	InitInlineFunctions_NamedFunction(this_, fn)
	return this_
}

func InitInlineFunctions_NamedFunction(this_ *InlineFunctions_NamedFunction, fn pkg_rhino.Node) {
	_ = fn
	this_.initFields_InlineFunctions_NamedFunction()
	this_.fn = fn
}

// Java method: NamedFunction.getName()
func (this_ *InlineFunctions_NamedFunction) GetName() string {
	return this_.fn.GetFirstChild().GetString()
}

// Java method: NamedFunction.getNameNode()
func (this_ *InlineFunctions_NamedFunction) GetNameNode() pkg_rhino.Node {
	return this_.fn.GetFirstChild()
}

// Java method: NamedFunction.getFunctionNode()
func (this_ *InlineFunctions_NamedFunction) GetFunctionNode() pkg_rhino.Node {
	return this_.fn
}

// Java method: NamedFunction.remove()
func (this_ *InlineFunctions_NamedFunction) Remove() {
	this_.outer_.compiler.ReportChangeToEnclosingScope(this_.fn)
	NodeUtil_RemoveChild(this_.fn.GetParent(), this_.fn)
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(this_.fn, this_.outer_.compiler)
}

// Java method: NamedFunction.getDeclaringBlock()
func (this_ *InlineFunctions_NamedFunction) GetDeclaringBlock() pkg_rhino.Node {
	return this_.fn.GetParent()
}

// Java class: com.google.javascript.jscomp.InlineFunctions.FunctionVar
type InlineFunctions_FunctionVar struct {
	outer_  *InlineFunctions
	var_    pkg_rhino.Node `gson:"var"`
	idHash_ atomic.Int32
}

func (this_ *InlineFunctions_FunctionVar) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineFunctions_FunctionVar) IsJavaNil_() bool { return this_ == nil }
func (*InlineFunctions_FunctionVar) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineFunctions$FunctionVar", "FunctionVar"
}
func (this_ *InlineFunctions_FunctionVar) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineFunctions_FunctionVar) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineFunctions.FunctionVar", "com.google.javascript.jscomp.InlineFunctions$FunctionVar", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineFunctions_FunctionVar](v_)
	})
}
func (this_ *InlineFunctions_FunctionVar) SetOuter_(o *InlineFunctions)          { this_.outer_ = o }
func (this_ *InlineFunctions_FunctionVar) GetOuter_() *InlineFunctions           { return this_.outer_ }
func (*InlineFunctions_FunctionVar) Implements_jscomp_InlineFunctions_Function() {}

func (this_ *InlineFunctions_FunctionVar) initFields_InlineFunctions_FunctionVar() {
	_ = this_
}

// Java constructor: FunctionVar(com.google.javascript.rhino.Node var)
func NewInlineFunctions_FunctionVar(outer_ *InlineFunctions, var_ pkg_rhino.Node) *InlineFunctions_FunctionVar {
	this_ := &InlineFunctions_FunctionVar{}
	this_.outer_ = outer_
	InitInlineFunctions_FunctionVar(this_, var_)
	return this_
}
func allocPooled_InlineFunctions_FunctionVar(slab_ *j2gort.PooledSlab[InlineFunctions_FunctionVar], outer_ *InlineFunctions, var_ pkg_rhino.Node) *InlineFunctions_FunctionVar {
	this_ := slab_.Alloc()
	this_.outer_ = outer_
	InitInlineFunctions_FunctionVar(this_, var_)
	return this_
}

func InitInlineFunctions_FunctionVar(this_ *InlineFunctions_FunctionVar, var_ pkg_rhino.Node) {
	_ = var_
	this_.initFields_InlineFunctions_FunctionVar()
	this_.var_ = var_
}

// Java method: FunctionVar.getName()
func (this_ *InlineFunctions_FunctionVar) GetName() string {
	return this_.var_.GetFirstChild().GetString()
}

// Java method: FunctionVar.getNameNode()
func (this_ *InlineFunctions_FunctionVar) GetNameNode() pkg_rhino.Node {
	return this_.var_.GetFirstChild()
}

// Java method: FunctionVar.getFunctionNode()
func (this_ *InlineFunctions_FunctionVar) GetFunctionNode() pkg_rhino.Node {
	return this_.var_.GetFirstFirstChild()
}

// Java method: FunctionVar.remove()
func (this_ *InlineFunctions_FunctionVar) Remove() {
	this_.outer_.compiler.ReportChangeToEnclosingScope(this_.var_)
	NodeUtil_RemoveChild(this_.var_.GetParent(), this_.var_)
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(this_.var_, this_.outer_.compiler)
}

// Java method: FunctionVar.getDeclaringBlock()
func (this_ *InlineFunctions_FunctionVar) GetDeclaringBlock() pkg_rhino.Node {
	return this_.var_.GetParent()
}

// Java class: com.google.javascript.jscomp.InlineFunctions.FunctionExpression
type InlineFunctions_FunctionExpression struct {
	fakeName     string         `gson:"fakeName"`
	fn           pkg_rhino.Node `gson:"fn"`
	fakeNameNode pkg_rhino.Node `gson:"fakeNameNode"`
	idHash_      atomic.Int32
}

func (this_ *InlineFunctions_FunctionExpression) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineFunctions_FunctionExpression) IsJavaNil_() bool { return this_ == nil }
func (*InlineFunctions_FunctionExpression) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineFunctions$FunctionExpression", "FunctionExpression"
}
func (this_ *InlineFunctions_FunctionExpression) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineFunctions_FunctionExpression) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineFunctions.FunctionExpression", "com.google.javascript.jscomp.InlineFunctions$FunctionExpression", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineFunctions_FunctionExpression](v_)
	})
}
func (*InlineFunctions_FunctionExpression) Implements_jscomp_InlineFunctions_Function() {}

func (this_ *InlineFunctions_FunctionExpression) initFields_InlineFunctions_FunctionExpression() {
	_ = this_
	this_.fakeName = j2gort.NullString
}

// Java constructor: FunctionExpression(com.google.javascript.rhino.Node fn, int index)
func NewInlineFunctions_FunctionExpression(fn pkg_rhino.Node, index int32) *InlineFunctions_FunctionExpression {
	this_ := &InlineFunctions_FunctionExpression{}
	InitInlineFunctions_FunctionExpression(this_, fn, index)
	return this_
}
func allocPooled_InlineFunctions_FunctionExpression(slab_ *j2gort.PooledSlab[InlineFunctions_FunctionExpression], fn pkg_rhino.Node, index int32) *InlineFunctions_FunctionExpression {
	this_ := slab_.Alloc()
	InitInlineFunctions_FunctionExpression(this_, fn, index)
	return this_
}

func InitInlineFunctions_FunctionExpression(this_ *InlineFunctions_FunctionExpression, fn pkg_rhino.Node, index int32) {
	_ = fn
	_ = index
	this_.initFields_InlineFunctions_FunctionExpression()
	this_.fn = fn
	this_.fakeName = strconv.FormatInt(int64(index), 10)
	this_.fakeNameNode = pkg_rhino.IR_Name(this_.fakeName)
}

// Java method: FunctionExpression.getName()
func (this_ *InlineFunctions_FunctionExpression) GetName() string {
	return this_.fakeName
}

// Java method: FunctionExpression.getNameNode()
func (this_ *InlineFunctions_FunctionExpression) GetNameNode() pkg_rhino.Node {
	return this_.fakeNameNode
}

// Java method: FunctionExpression.getFunctionNode()
func (this_ *InlineFunctions_FunctionExpression) GetFunctionNode() pkg_rhino.Node {
	return this_.fn
}

// Java method: FunctionExpression.remove()
func (this_ *InlineFunctions_FunctionExpression) Remove() {
}

// Java method: FunctionExpression.getDeclaringBlock()
func (this_ *InlineFunctions_FunctionExpression) GetDeclaringBlock() pkg_rhino.Node {
	return pkg_rhino.Node(0)
}

// Java class: com.google.javascript.jscomp.InlineFunctions.Reference
type InlineFunctions_Reference struct {
	FunctionInjector_ReferenceBase
	requiresDecomposition bool `gson:"requiresDecomposition"`
	inlined               bool `gson:"inlined"`
}

func (this_ *InlineFunctions_Reference) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineFunctions_Reference) IsJavaNil_() bool { return this_ == nil }
func (*InlineFunctions_Reference) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineFunctions$Reference", "Reference"
}
func (this_ *InlineFunctions_Reference) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.FunctionInjector_ReferenceBase.IdentityHash_()
}
func (*InlineFunctions_Reference) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineFunctions.Reference", "com.google.javascript.jscomp.InlineFunctions$Reference", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineFunctions_Reference](v_)
	})
}
func (*InlineFunctions_Reference) Implements_jscomp_FunctionInjector_Reference() {}

func (this_ *InlineFunctions_Reference) initFields_InlineFunctions_Reference() {
	_ = this_

}

// Java constructor: Reference(com.google.javascript.rhino.Node callNode, com.google.javascript.jscomp.Scope scope, com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.FunctionInjector.InliningMode mode)
func NewInlineFunctions_Reference(callNode pkg_rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) *InlineFunctions_Reference {
	this_ := &InlineFunctions_Reference{}
	InitInlineFunctions_Reference(this_, callNode, scope, chunk, mode)
	return this_
}
func allocPooled_InlineFunctions_Reference(slab_ *j2gort.PooledSlab[InlineFunctions_Reference], callNode pkg_rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) *InlineFunctions_Reference {
	this_ := slab_.Alloc()
	InitInlineFunctions_Reference(this_, callNode, scope, chunk, mode)
	return this_
}

func InitInlineFunctions_Reference(this_ *InlineFunctions_Reference, callNode pkg_rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) {
	_ = callNode
	_ = scope
	_ = chunk
	_ = mode
	InitFunctionInjector_Reference(&this_.FunctionInjector_ReferenceBase, callNode, scope, chunk, mode)
	this_.initFields_InlineFunctions_Reference()
}

// Java method: Reference.setRequiresDecomposition(boolean newVal)
func (this_ *InlineFunctions_Reference) SetRequiresDecomposition(newVal bool) {
	_ = newVal
	this_.requiresDecomposition = newVal
}

// Java enum: com.google.javascript.jscomp.InlineFunctions.ShouldWarnWhenRequireInliningCannotInline
type InlineFunctions_ShouldWarnWhenRequireInliningCannotInline uint8

const (
	enum_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES InlineFunctions_ShouldWarnWhenRequireInliningCannotInline = 1
	enum_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO  InlineFunctions_ShouldWarnWhenRequireInliningCannotInline = 2
)

var enumBaseTable_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline = [...]j2gort.EnumBase{
	{Name_: "YES", Ordinal_: 0},
	{Name_: "NO", Ordinal_: 1},
}
var enumValues_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline = []InlineFunctions_ShouldWarnWhenRequireInliningCannotInline{1, 2}

func StaticGet_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES() InlineFunctions_ShouldWarnWhenRequireInliningCannotInline {
	return enum_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES
}
func StaticGet_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO() InlineFunctions_ShouldWarnWhenRequireInliningCannotInline {
	return enum_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO
}

func InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_Values() []InlineFunctions_ShouldWarnWhenRequireInliningCannotInline {
	return j2gort.CloneSlice(enumValues_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline)
}

func InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_ValueOf(name string) InlineFunctions_ShouldWarnWhenRequireInliningCannotInline {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline {
		if enumBaseTable_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline[i_].Name_ == name {
			return InlineFunctions_ShouldWarnWhenRequireInliningCannotInline(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.InlineFunctions.ShouldWarnWhenRequireInliningCannotInline." + name))
}

func (this_ InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) IsJavaNil_() bool {
	return this_ == 0
}
func (InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) HasIdentityEquals_() {}
func (InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) Implements_j2gort_Ext_java_lang_Enum() {
}
func (InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineFunctions$ShouldWarnWhenRequireInliningCannotInline", "ShouldWarnWhenRequireInliningCannotInline"
}
func (this_ InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline[this_-1]
}
func (this_ InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) EnumOrdinal_() int32 {
	return int32(this_) - 1
}
func (this_ InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline[this_-1].Name_
}
func (this_ InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline[this_-1].Name_
}
func (this_ InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) Equals(other_ any) bool {
	o_, ok_ := other_.(InlineFunctions_ShouldWarnWhenRequireInliningCannotInline)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline[this_-1].HashCode()
}
func (this_ InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline[this_-1].CompareTo(other_)
}
func (this_ InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_InlineFunctions_ShouldWarnWhenRequireInliningCannotInline[this_-1].Name_)
}
func init() {
	vals_ := []any{InlineFunctions_ShouldWarnWhenRequireInliningCannotInline(1), InlineFunctions_ShouldWarnWhenRequireInliningCannotInline(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.InlineFunctions.ShouldWarnWhenRequireInliningCannotInline", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.InlineFunctions$ShouldWarnWhenRequireInliningCannotInline", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineFunctions.ShouldWarnWhenRequireInliningCannotInline", "com.google.javascript.jscomp.InlineFunctions$ShouldWarnWhenRequireInliningCannotInline", func(v_ any) bool {
		return j2gort.InstanceOf[InlineFunctions_ShouldWarnWhenRequireInliningCannotInline](v_)
	})
}

// Java enum: com.google.javascript.jscomp.InlineFunctions.CannotRemoveReason
type InlineFunctions_CannotRemoveReason uint8

const (
	enum_InlineFunctions_CannotRemoveReason_BLOCK_INLINING_REFERENCE   InlineFunctions_CannotRemoveReason = 1
	enum_InlineFunctions_CannotRemoveReason_CALLED_BY_INLINED_FUNCTION InlineFunctions_CannotRemoveReason = 2
	enum_InlineFunctions_CannotRemoveReason_UNINLINABLE_REFERENCE      InlineFunctions_CannotRemoveReason = 3
	enum_InlineFunctions_CannotRemoveReason_UNREMOVABLE_REFERENCE      InlineFunctions_CannotRemoveReason = 4
)

var enumBaseTable_InlineFunctions_CannotRemoveReason = [...]j2gort.EnumBase{
	{Name_: "BLOCK_INLINING_REFERENCE", Ordinal_: 0},
	{Name_: "CALLED_BY_INLINED_FUNCTION", Ordinal_: 1},
	{Name_: "UNINLINABLE_REFERENCE", Ordinal_: 2},
	{Name_: "UNREMOVABLE_REFERENCE", Ordinal_: 3},
}
var enumValues_InlineFunctions_CannotRemoveReason = []InlineFunctions_CannotRemoveReason{1, 2, 3, 4}

func StaticGet_InlineFunctions_CannotRemoveReason_BLOCK_INLINING_REFERENCE() InlineFunctions_CannotRemoveReason {
	return enum_InlineFunctions_CannotRemoveReason_BLOCK_INLINING_REFERENCE
}
func StaticGet_InlineFunctions_CannotRemoveReason_CALLED_BY_INLINED_FUNCTION() InlineFunctions_CannotRemoveReason {
	return enum_InlineFunctions_CannotRemoveReason_CALLED_BY_INLINED_FUNCTION
}
func StaticGet_InlineFunctions_CannotRemoveReason_UNINLINABLE_REFERENCE() InlineFunctions_CannotRemoveReason {
	return enum_InlineFunctions_CannotRemoveReason_UNINLINABLE_REFERENCE
}
func StaticGet_InlineFunctions_CannotRemoveReason_UNREMOVABLE_REFERENCE() InlineFunctions_CannotRemoveReason {
	return enum_InlineFunctions_CannotRemoveReason_UNREMOVABLE_REFERENCE
}

func InlineFunctions_CannotRemoveReason_Values() []InlineFunctions_CannotRemoveReason {
	return j2gort.CloneSlice(enumValues_InlineFunctions_CannotRemoveReason)
}

func InlineFunctions_CannotRemoveReason_ValueOf(name string) InlineFunctions_CannotRemoveReason {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_InlineFunctions_CannotRemoveReason {
		if enumBaseTable_InlineFunctions_CannotRemoveReason[i_].Name_ == name {
			return InlineFunctions_CannotRemoveReason(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.InlineFunctions.CannotRemoveReason." + name))
}

func (this_ InlineFunctions_CannotRemoveReason) IsJavaNil_() bool                              { return this_ == 0 }
func (InlineFunctions_CannotRemoveReason) HasIdentityEquals_()                                 {}
func (InlineFunctions_CannotRemoveReason) Implements_j2gort_Ext_java_lang_Enum()               {}
func (InlineFunctions_CannotRemoveReason) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ InlineFunctions_CannotRemoveReason) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (InlineFunctions_CannotRemoveReason) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineFunctions$CannotRemoveReason", "CannotRemoveReason"
}
func (this_ InlineFunctions_CannotRemoveReason) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_InlineFunctions_CannotRemoveReason[this_-1]
}
func (this_ InlineFunctions_CannotRemoveReason) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ InlineFunctions_CannotRemoveReason) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ InlineFunctions_CannotRemoveReason) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_InlineFunctions_CannotRemoveReason[this_-1].Name_
}
func (this_ InlineFunctions_CannotRemoveReason) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_InlineFunctions_CannotRemoveReason[this_-1].Name_
}
func (this_ InlineFunctions_CannotRemoveReason) Equals(other_ any) bool {
	o_, ok_ := other_.(InlineFunctions_CannotRemoveReason)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ InlineFunctions_CannotRemoveReason) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_InlineFunctions_CannotRemoveReason[this_-1].HashCode()
}
func (this_ InlineFunctions_CannotRemoveReason) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_InlineFunctions_CannotRemoveReason[this_-1].CompareTo(other_)
}
func (this_ InlineFunctions_CannotRemoveReason) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_InlineFunctions_CannotRemoveReason[this_-1].Name_)
}
func init() {
	vals_ := []any{InlineFunctions_CannotRemoveReason(1), InlineFunctions_CannotRemoveReason(2), InlineFunctions_CannotRemoveReason(3), InlineFunctions_CannotRemoveReason(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.InlineFunctions.CannotRemoveReason", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.InlineFunctions$CannotRemoveReason", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineFunctions.CannotRemoveReason", "com.google.javascript.jscomp.InlineFunctions$CannotRemoveReason", func(v_ any) bool {
		return j2gort.InstanceOf[InlineFunctions_CannotRemoveReason](v_)
	})
}

// Java enum: com.google.javascript.jscomp.InlineFunctions.DisallowInliningReason
type InlineFunctions_DisallowInliningReason uint8

const (
	enum_InlineFunctions_DisallowInliningReason_ASYNC_FUNCTION                 InlineFunctions_DisallowInliningReason = 1
	enum_InlineFunctions_DisallowInliningReason_CAPTURED_VARIABLE_NOT_IN_SCOPE InlineFunctions_DisallowInliningReason = 2
	enum_InlineFunctions_DisallowInliningReason_COST                           InlineFunctions_DisallowInliningReason = 3
	enum_InlineFunctions_DisallowInliningReason_GENERATOR_FUNCTION             InlineFunctions_DisallowInliningReason = 4
	enum_InlineFunctions_DisallowInliningReason_HAS_LOCAL_NAMES                InlineFunctions_DisallowInliningReason = 5
	enum_InlineFunctions_DisallowInliningReason_MULTIPLE_FUNCTION_DEFINITIONS  InlineFunctions_DisallowInliningReason = 6
	enum_InlineFunctions_DisallowInliningReason_NON_INLINABLE_PARAM            InlineFunctions_DisallowInliningReason = 7
	enum_InlineFunctions_DisallowInliningReason_NOT_CANDIDATE_FUNCTION         InlineFunctions_DisallowInliningReason = 8
	enum_InlineFunctions_DisallowInliningReason_NO_INLINE_ANNOTATION           InlineFunctions_DisallowInliningReason = 9
	enum_InlineFunctions_DisallowInliningReason_REASSIGNED                     InlineFunctions_DisallowInliningReason = 10
	enum_InlineFunctions_DisallowInliningReason_TARGET_SIZE_EXCEEDS_LIMIT      InlineFunctions_DisallowInliningReason = 11
)

var enumBaseTable_InlineFunctions_DisallowInliningReason = [...]j2gort.EnumBase{
	{Name_: "ASYNC_FUNCTION", Ordinal_: 0},
	{Name_: "CAPTURED_VARIABLE_NOT_IN_SCOPE", Ordinal_: 1},
	{Name_: "COST", Ordinal_: 2},
	{Name_: "GENERATOR_FUNCTION", Ordinal_: 3},
	{Name_: "HAS_LOCAL_NAMES", Ordinal_: 4},
	{Name_: "MULTIPLE_FUNCTION_DEFINITIONS", Ordinal_: 5},
	{Name_: "NON_INLINABLE_PARAM", Ordinal_: 6},
	{Name_: "NOT_CANDIDATE_FUNCTION", Ordinal_: 7},
	{Name_: "NO_INLINE_ANNOTATION", Ordinal_: 8},
	{Name_: "REASSIGNED", Ordinal_: 9},
	{Name_: "TARGET_SIZE_EXCEEDS_LIMIT", Ordinal_: 10},
}
var enumValues_InlineFunctions_DisallowInliningReason = []InlineFunctions_DisallowInliningReason{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

func StaticGet_InlineFunctions_DisallowInliningReason_ASYNC_FUNCTION() InlineFunctions_DisallowInliningReason {
	return enum_InlineFunctions_DisallowInliningReason_ASYNC_FUNCTION
}
func StaticGet_InlineFunctions_DisallowInliningReason_CAPTURED_VARIABLE_NOT_IN_SCOPE() InlineFunctions_DisallowInliningReason {
	return enum_InlineFunctions_DisallowInliningReason_CAPTURED_VARIABLE_NOT_IN_SCOPE
}
func StaticGet_InlineFunctions_DisallowInliningReason_COST() InlineFunctions_DisallowInliningReason {
	return enum_InlineFunctions_DisallowInliningReason_COST
}
func StaticGet_InlineFunctions_DisallowInliningReason_GENERATOR_FUNCTION() InlineFunctions_DisallowInliningReason {
	return enum_InlineFunctions_DisallowInliningReason_GENERATOR_FUNCTION
}
func StaticGet_InlineFunctions_DisallowInliningReason_HAS_LOCAL_NAMES() InlineFunctions_DisallowInliningReason {
	return enum_InlineFunctions_DisallowInliningReason_HAS_LOCAL_NAMES
}
func StaticGet_InlineFunctions_DisallowInliningReason_MULTIPLE_FUNCTION_DEFINITIONS() InlineFunctions_DisallowInliningReason {
	return enum_InlineFunctions_DisallowInliningReason_MULTIPLE_FUNCTION_DEFINITIONS
}
func StaticGet_InlineFunctions_DisallowInliningReason_NON_INLINABLE_PARAM() InlineFunctions_DisallowInliningReason {
	return enum_InlineFunctions_DisallowInliningReason_NON_INLINABLE_PARAM
}
func StaticGet_InlineFunctions_DisallowInliningReason_NOT_CANDIDATE_FUNCTION() InlineFunctions_DisallowInliningReason {
	return enum_InlineFunctions_DisallowInliningReason_NOT_CANDIDATE_FUNCTION
}
func StaticGet_InlineFunctions_DisallowInliningReason_NO_INLINE_ANNOTATION() InlineFunctions_DisallowInliningReason {
	return enum_InlineFunctions_DisallowInliningReason_NO_INLINE_ANNOTATION
}
func StaticGet_InlineFunctions_DisallowInliningReason_REASSIGNED() InlineFunctions_DisallowInliningReason {
	return enum_InlineFunctions_DisallowInliningReason_REASSIGNED
}
func StaticGet_InlineFunctions_DisallowInliningReason_TARGET_SIZE_EXCEEDS_LIMIT() InlineFunctions_DisallowInliningReason {
	return enum_InlineFunctions_DisallowInliningReason_TARGET_SIZE_EXCEEDS_LIMIT
}

func InlineFunctions_DisallowInliningReason_Values() []InlineFunctions_DisallowInliningReason {
	return j2gort.CloneSlice(enumValues_InlineFunctions_DisallowInliningReason)
}

func InlineFunctions_DisallowInliningReason_ValueOf(name string) InlineFunctions_DisallowInliningReason {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_InlineFunctions_DisallowInliningReason {
		if enumBaseTable_InlineFunctions_DisallowInliningReason[i_].Name_ == name {
			return InlineFunctions_DisallowInliningReason(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.InlineFunctions.DisallowInliningReason." + name))
}

func (this_ InlineFunctions_DisallowInliningReason) IsJavaNil_() bool                              { return this_ == 0 }
func (InlineFunctions_DisallowInliningReason) HasIdentityEquals_()                                 {}
func (InlineFunctions_DisallowInliningReason) Implements_j2gort_Ext_java_lang_Enum()               {}
func (InlineFunctions_DisallowInliningReason) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ InlineFunctions_DisallowInliningReason) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (InlineFunctions_DisallowInliningReason) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineFunctions$DisallowInliningReason", "DisallowInliningReason"
}
func (this_ InlineFunctions_DisallowInliningReason) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_InlineFunctions_DisallowInliningReason[this_-1]
}
func (this_ InlineFunctions_DisallowInliningReason) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ InlineFunctions_DisallowInliningReason) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ InlineFunctions_DisallowInliningReason) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_InlineFunctions_DisallowInliningReason[this_-1].Name_
}
func (this_ InlineFunctions_DisallowInliningReason) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_InlineFunctions_DisallowInliningReason[this_-1].Name_
}
func (this_ InlineFunctions_DisallowInliningReason) Equals(other_ any) bool {
	o_, ok_ := other_.(InlineFunctions_DisallowInliningReason)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ InlineFunctions_DisallowInliningReason) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_InlineFunctions_DisallowInliningReason[this_-1].HashCode()
}
func (this_ InlineFunctions_DisallowInliningReason) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_InlineFunctions_DisallowInliningReason[this_-1].CompareTo(other_)
}
func (this_ InlineFunctions_DisallowInliningReason) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_InlineFunctions_DisallowInliningReason[this_-1].Name_)
}
func init() {
	vals_ := []any{InlineFunctions_DisallowInliningReason(1), InlineFunctions_DisallowInliningReason(2), InlineFunctions_DisallowInliningReason(3), InlineFunctions_DisallowInliningReason(4), InlineFunctions_DisallowInliningReason(5), InlineFunctions_DisallowInliningReason(6), InlineFunctions_DisallowInliningReason(7), InlineFunctions_DisallowInliningReason(8), InlineFunctions_DisallowInliningReason(9), InlineFunctions_DisallowInliningReason(10), InlineFunctions_DisallowInliningReason(11)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.InlineFunctions.DisallowInliningReason", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.InlineFunctions$DisallowInliningReason", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineFunctions.DisallowInliningReason", "com.google.javascript.jscomp.InlineFunctions$DisallowInliningReason", func(v_ any) bool {
		return j2gort.InstanceOf[InlineFunctions_DisallowInliningReason](v_)
	})
}

// Java record: com.google.javascript.jscomp.InlineFunctions.ShouldRequireInlining
type InlineFunctions_ShouldRequireInlining struct {
	optMessage     string `gson:"optMessage"`
	shouldValidate bool   `gson:"shouldValidate"`
}

func (this_ *InlineFunctions_ShouldRequireInlining) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineFunctions_ShouldRequireInlining) IsJavaNil_() bool { return this_ == nil }
func (*InlineFunctions_ShouldRequireInlining) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineFunctions$ShouldRequireInlining", "ShouldRequireInlining"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineFunctions.ShouldRequireInlining", "com.google.javascript.jscomp.InlineFunctions$ShouldRequireInlining", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineFunctions_ShouldRequireInlining](v_)
	})
}
func (*InlineFunctions_ShouldRequireInlining) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *InlineFunctions_ShouldRequireInlining) initFields_InlineFunctions_ShouldRequireInlining() {
	_ = this_
}

// Java constructor: ShouldRequireInlining(boolean shouldValidate, java.lang.String optMessage)
func NewInlineFunctions_ShouldRequireInlining(shouldValidate bool, optMessage string) *InlineFunctions_ShouldRequireInlining {
	this_ := &InlineFunctions_ShouldRequireInlining{}
	InitInlineFunctions_ShouldRequireInlining(this_, shouldValidate, optMessage)
	return this_
}

func InitInlineFunctions_ShouldRequireInlining(this_ *InlineFunctions_ShouldRequireInlining, shouldValidate bool, optMessage string) {
	_ = shouldValidate
	_ = optMessage
	this_.initFields_InlineFunctions_ShouldRequireInlining()
	this_.shouldValidate = shouldValidate
	this_.optMessage = optMessage
}

func (this_ *InlineFunctions_ShouldRequireInlining) ShouldValidate() bool {
	return this_.shouldValidate
}
func (this_ *InlineFunctions_ShouldRequireInlining) OptMessage() string { return this_.optMessage }
func (this_ *InlineFunctions_ShouldRequireInlining) Equals(other_ any) bool {
	o_, ok_ := other_.(*InlineFunctions_ShouldRequireInlining)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.shouldValidate != o_.shouldValidate {
		return false
	}
	if this_.optMessage != o_.optMessage {
		return false
	}
	return true
}
func (this_ *InlineFunctions_ShouldRequireInlining) HashCode() int32 {
	var h_ int32 = 0
	if this_.shouldValidate {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	{
		hs_25 := this_.optMessage
		h_ = 31 * h_
		if hs_25 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_25)
		}
	}
	return h_
}
func (this_ *InlineFunctions_ShouldRequireInlining) String() string {
	_ = fmt.Sprint
	return "ShouldRequireInlining[" + "shouldValidate=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.shouldValidate)) + ", " + "optMessage=" + j2gort.StrVal(j2gort.BoxString(this_.optMessage)) + "]"
}
func (this_ *InlineFunctions) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_InlineFunctions_Reference.Release()
	this_.slab_InlineFunctions_FunctionVar.Release()
	this_.slab_InlineFunctions_NamedFunction.Release()
	this_.slab_InlineFunctions_FunctionExpression.Release()
}
func (*InlineFunctions_FindCandidateFunctions) AlwaysShouldTraverse_()      {}
func (*InlineFunctions_FindCandidateFunctions) CanReleaseSyntacticScopes_() {}
func (*InlineFunctions_FindCandidatesReferences) AlwaysShouldTraverse_()    {}
func InlineFunctions_FindCalledFunctions(node pkg_rhino.Node, changed j2gort.Ext_java_util_Set) {
	clinit_InlineFunctions()
	_ = node
	_ = changed
	if !(changed != nil) {
		j2gort.CheckArgument(false)
	}
	if node.IsName() && InlineFunctions_IsCandidateUsage(node) {
		j2gort.SetAddStr(changed, node.GetString())
	}
	{
		var c pkg_rhino.Node = node.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			InlineFunctions_FindCalledFunctions(c, changed)
		}
	}
}
