/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RescopeGlobalSymbolsRewriteCallback.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_rescope_global_symbols_rewrite_callback_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_rescope_global_symbols_rewrite_callback_2 j2gort.Ext_java_util_function_Function
var lambdaSingleton_rescope_global_symbols_rewrite_callback_3 j2gort.Ext_java_util_function_Function
var lambdaSingleton_rescope_global_symbols_rewrite_callback_4 j2gort.Ext_java_util_function_Function
var lambdaSingleton_rescope_global_symbols_rewrite_callback_5 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_rescope_global_symbols_rewrite_callback_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_9 any) any {
		k := j2gort.CastString(lp_9)
		_ = k
		return j2gort.NewExt_java_util_LinkedHashSet()
	})
	lambdaSingleton_rescope_global_symbols_rewrite_callback_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_14 any) any {
		k := j2gort.CastPtr[JSChunk](lp_14)
		_ = k
		return j2gort.NewExt_java_util_LinkedHashSet()
	})
	lambdaSingleton_rescope_global_symbols_rewrite_callback_3 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_19 any) any {
		k := j2gort.CastString(lp_19)
		_ = k
		return j2gort.NewExt_java_util_LinkedHashSet()
	})
	lambdaSingleton_rescope_global_symbols_rewrite_callback_4 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_24 any) any {
		k := j2gort.CastString(lp_24)
		_ = k
		return j2gort.NewExt_java_util_LinkedHashSet()
	})
	lambdaSingleton_rescope_global_symbols_rewrite_callback_5 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_29 any) any {
		k := j2gort.CastString(lp_29)
		_ = k
		return j2gort.NewExt_java_util_LinkedHashSet()
	})
}

// Java class: com.google.javascript.jscomp.RescopeGlobalSymbolsRewriteCallback
type RescopeGlobalSymbolsRewriteCallback struct {
	compiler                                AbstractCompiler                                      `gson:"compiler"`
	symbolInfo                              *RescopeGlobalSymbolsInfo                             `gson:"symbolInfo"`
	renameHelper                            *RescopeGlobalSymbolsRenameHelper                     `gson:"renameHelper"`
	preDeclarations                         []*RescopeGlobalSymbolsRewriteCallback_ChunkGlobal    `gson:"preDeclarations"`
	localAliasesForUnwrappedCrossChunkNames j2gort.OrderedMap[*JSChunk, j2gort.Ext_java_util_Set] `gson:"localAliasesForUnwrappedCrossChunkNames"`
	chunksUsingCrossChunkSymbols            j2gort.OrderedMap[string, j2gort.Ext_java_util_Set]   `gson:"chunksUsingCrossChunkSymbols"`
	chunksUsingWrappedReassignableSymbols   j2gort.OrderedMap[string, j2gort.Ext_java_util_Set]   `gson:"chunksUsingWrappedReassignableSymbols"`
	chunksRequiringWrapperInitialization    j2gort.OrderedMap[string, j2gort.Ext_java_util_Set]   `gson:"chunksRequiringWrapperInitialization"`
	idHash_                                 atomic.Int32
	optimizeLocalAccess                     CompilerOptions_OptimizeLocalAccess `gson:"optimizeLocalAccess"`
}

func (this_ *RescopeGlobalSymbolsRewriteCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *RescopeGlobalSymbolsRewriteCallback) IsJavaNil_() bool { return this_ == nil }
func (*RescopeGlobalSymbolsRewriteCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RescopeGlobalSymbolsRewriteCallback", "RescopeGlobalSymbolsRewriteCallback"
}
func (this_ *RescopeGlobalSymbolsRewriteCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RescopeGlobalSymbolsRewriteCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RescopeGlobalSymbolsRewriteCallback", "com.google.javascript.jscomp.RescopeGlobalSymbolsRewriteCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*RescopeGlobalSymbolsRewriteCallback](v_)
	})
}
func (*RescopeGlobalSymbolsRewriteCallback) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *RescopeGlobalSymbolsRewriteCallback) initFields_RescopeGlobalSymbolsRewriteCallback() {
	_ = this_
	this_.preDeclarations = []*RescopeGlobalSymbolsRewriteCallback_ChunkGlobal{}
	this_.localAliasesForUnwrappedCrossChunkNames = j2gort.NewOrderedMap[*JSChunk, j2gort.Ext_java_util_Set]()
	this_.chunksUsingCrossChunkSymbols = j2gort.NewOrderedMap[string, j2gort.Ext_java_util_Set]()
	this_.chunksUsingWrappedReassignableSymbols = j2gort.NewOrderedMap[string, j2gort.Ext_java_util_Set]()
	this_.chunksRequiringWrapperInitialization = j2gort.NewOrderedMap[string, j2gort.Ext_java_util_Set]()
}

// Java constructor: RescopeGlobalSymbolsRewriteCallback(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.CompilerOptions.OptimizeLocalAccess optimizeLocalAccess, com.google.javascript.jscomp.RescopeGlobalSymbolsInfo symbolInfo, com.google.javascript.jscomp.RescopeGlobalSymbolsRenameHelper renameHelper)
func NewRescopeGlobalSymbolsRewriteCallback(compiler AbstractCompiler, optimizeLocalAccess CompilerOptions_OptimizeLocalAccess, symbolInfo *RescopeGlobalSymbolsInfo, renameHelper *RescopeGlobalSymbolsRenameHelper) *RescopeGlobalSymbolsRewriteCallback {
	this_ := &RescopeGlobalSymbolsRewriteCallback{}
	InitRescopeGlobalSymbolsRewriteCallback(this_, compiler, optimizeLocalAccess, symbolInfo, renameHelper)
	return this_
}

func InitRescopeGlobalSymbolsRewriteCallback(this_ *RescopeGlobalSymbolsRewriteCallback, compiler AbstractCompiler, optimizeLocalAccess CompilerOptions_OptimizeLocalAccess, symbolInfo *RescopeGlobalSymbolsInfo, renameHelper *RescopeGlobalSymbolsRenameHelper) {
	_ = compiler
	_ = optimizeLocalAccess
	_ = symbolInfo
	_ = renameHelper
	this_.initFields_RescopeGlobalSymbolsRewriteCallback()
	this_.compiler = compiler
	this_.optimizeLocalAccess = optimizeLocalAccess
	this_.symbolInfo = symbolInfo
	this_.renameHelper = renameHelper
}

// Java method: RescopeGlobalSymbolsRewriteCallback.isCrossChunkName(java.lang.String name)
func (this_ *RescopeGlobalSymbolsRewriteCallback) priv_RescopeGlobalSymbolsRewriteCallback_isCrossChunkName_4956(name string) bool {
	_ = name
	return j2gort.SetContainsStr(this_.symbolInfo.CrossChunkNames(), name)
}

// Java method: RescopeGlobalSymbolsRewriteCallback.isExternVar(java.lang.String varname, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *RescopeGlobalSymbolsRewriteCallback) priv_RescopeGlobalSymbolsRewriteCallback_isExternVar_4957(varname string, t *NodeTraversal) bool {
	_ = varname
	_ = t
	if j2gort.CheckNotNullString(varname) == "" {
		return false
	}
	var v *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(varname))
	_ = v
	return (((v == nil) || v.IsExtern()) || (j2gort.CastPtr[Scope](v.GetScope()).IsGlobal() && j2gort.SetContainsStr(this_.symbolInfo.ExternNames(), varname)))
}

// Java method: RescopeGlobalSymbolsRewriteCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RescopeGlobalSymbolsRewriteCallback) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if NodeUtil_IsNameDeclaration(n) {
		this_.priv_RescopeGlobalSymbolsRewriteCallback_visitNameDeclaration_4958(t, n)
	}
	return true
}

// Java method: RescopeGlobalSymbolsRewriteCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RescopeGlobalSymbolsRewriteCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsName() {
		this_.priv_RescopeGlobalSymbolsRewriteCallback_visitName_4963(t, n, parent)
	}
}

// Java method: RescopeGlobalSymbolsRewriteCallback.visitNameDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node declaration)
func (this_ *RescopeGlobalSymbolsRewriteCallback) priv_RescopeGlobalSymbolsRewriteCallback_visitNameDeclaration_4958(t *NodeTraversal, declaration pkg_rhino.Node) {
	_ = t
	_ = declaration
	var allLhsNodes j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	_ = allLhsNodes
	var iifeArg_1 j2gort.
		Ext_java_util_function_Consumer
	{

		bound_1 := allLhsNodes
		iifeArg_1 = j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) { bound_1.Add_1(mrp0_) })
	}
	NodeUtil_VisitLhsNodesInNode(declaration, iifeArg_1)
	if allLhsNodes.IsEmpty() {
		return
	}
	var hasImportantName bool = false
	_ = hasImportantName
	var isGlobalDeclaration bool = j2gort.CastPtr[Var](t.GetScope().GetVar(j2gort.CastNode[pkg_rhino.Node](allLhsNodes.Get(int32(0))).GetString())).IsGlobal()
	_ = isGlobalDeclaration
	{
		cur_2 := j2gort.NewIterCursor(allLhsNodes)
		for cur_2.HasNext() {
			lhs := pkg_rhino.Node(cur_2.NextNode())
			_ = lhs
			if !(lhs.IsName()) {
				j2gort.CheckStateFmt(false, "Unexpected lhs node %s, expected NAME", j2gort.BoxNode(lhs))
			}
			if (isGlobalDeclaration && this_.priv_RescopeGlobalSymbolsRewriteCallback_isCrossChunkName_4956(lhs.GetString())) || this_.priv_RescopeGlobalSymbolsRewriteCallback_isExternVar_4957(lhs.GetString(), t) {
				hasImportantName = true
				break
			}
		}
	}
	if hasImportantName {
		this_.priv_RescopeGlobalSymbolsRewriteCallback_rewriteNameDeclaration_4960(t, declaration, allLhsNodes, isGlobalDeclaration)
	}
}

// Java method: RescopeGlobalSymbolsRewriteCallback.globalVarStaysLocalInDefiningChunk(java.lang.String name, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *RescopeGlobalSymbolsRewriteCallback) priv_RescopeGlobalSymbolsRewriteCallback_globalVarStaysLocalInDefiningChunk_4959(name string, t *NodeTraversal) bool {
	_ = name
	_ = t
	if !(j2gort.CastPtr[Var](t.GetScope().GetVar(name)).IsGlobal()) {
		j2gort.CheckStateFmt(false, "Expected global variable: %s", j2gort.BoxString(name))
	}
	if this_.priv_RescopeGlobalSymbolsRewriteCallback_isExternVar_4957(name, t) {
		return false
	}
	if !(this_.priv_RescopeGlobalSymbolsRewriteCallback_isCrossChunkName_4956(name)) {
		return true
	}
	if !(this_.optimizeLocalAccess.ShouldOptimize()) {
		return false
	}
	if this_.optimizeLocalAccess.ShouldWrapReassignableSymbols() && j2gort.SetContainsStr(this_.symbolInfo.WrappedCrossChunkNames(), name) {
		return this_.renameHelper.HasUnwrappedLocalNameInDefiningChunk(name)
	}
	return (!(j2gort.SetContainsStr(this_.symbolInfo.CrossChunkNamesWithWriteFromOtherChunk(), name)) && j2gort.SetContainsStr(this_.symbolInfo.GlobalNamesWithReadInDefiningChunk(), name))
}

// Java method: RescopeGlobalSymbolsRewriteCallback.rewriteNameDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node declaration, java.util.List<com.google.javascript.rhino.Node> allLhsNodes, boolean isGlobalDeclaration)
func (this_ *RescopeGlobalSymbolsRewriteCallback) priv_RescopeGlobalSymbolsRewriteCallback_rewriteNameDeclaration_4960(t *NodeTraversal, declaration pkg_rhino.Node, allLhsNodes j2gort.Ext_java_util_List, isGlobalDeclaration bool) {
	_ = t
	_ = declaration
	_ = allLhsNodes
	_ = isGlobalDeclaration
	if this_.optimizeLocalAccess.ShouldPreserveLocalVarDeclarations() && !(NodeUtil_IsAnyFor(declaration.GetParent())) {
		var changed bool = false
		_ = changed
		if (isGlobalDeclaration && !(declaration.IsVar())) && this_.priv_RescopeGlobalSymbolsRewriteCallback_containsNonLocalOrInitializedVars_4962(t, declaration, isGlobalDeclaration) {
			declaration.SetToken(pkg_rhino.StaticGet_Token_VAR())
			changed = true
		}
		{
			var child pkg_rhino.Node = declaration.GetFirstChild()
			_ = child
			for child != 0 {
				var next pkg_rhino.Node = child.GetNext()
				_ = next
				if child.IsName() {
					var name string = child.GetString()
					_ = name
					if !(this_.priv_RescopeGlobalSymbolsRewriteCallback_globalVarStaysLocalInDefiningChunk_4959(name, t)) {
						if child.HasChildren() {
							var assign pkg_rhino.Node = pkg_rhino.IR_Assign(child.CloneNode_0(), child.RemoveFirstChild())
							_ = assign
							child.ReplaceWith(assign)
							assign.SetJSDocInfo(declaration.GetJSDocInfo())
							changed = true
						} else if this_.priv_RescopeGlobalSymbolsRewriteCallback_isExternVar_4957(name, t) {
							child.Detach()
							changed = true
						}
					}
				} else if child.IsDestructuringLhs() {
					NodeUtil_VisitLhsNodesInNode(child, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_3 any) {
						lhs := j2gort.CastNode[pkg_rhino.Node](lp_3)
						_ = lhs
						this_.priv_RescopeGlobalSymbolsRewriteCallback_maybeAddPreDeclaration_4961(t, lhs, isGlobalDeclaration)
					}))
					var assign_2 pkg_rhino.Node = pkg_rhino.IR_Assign(child.RemoveFirstChild(), child.RemoveFirstChild())
					_ = assign_2
					child.ReplaceWith(assign_2)
					assign_2.SetJSDocInfo(declaration.GetJSDocInfo())
					changed = true
				}
				child = next
			}
		}
		if changed {
			this_.compiler.ReportChangeToEnclosingScope(declaration)
		}
		return
	}
	if !(this_.optimizeLocalAccess.ShouldOptimize()) || this_.priv_RescopeGlobalSymbolsRewriteCallback_containsNonLocalOrInitializedVars_4962(t, declaration, isGlobalDeclaration) {
		{
			cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(allLhsNodes))
			for cur_4.HasNext() {
				lhs := pkg_rhino.Node(cur_4.NextNode())
				_ = lhs
				this_.priv_RescopeGlobalSymbolsRewriteCallback_maybeAddPreDeclaration_4961(t, lhs, isGlobalDeclaration)
			}
		}
	}
	{
		var child_2 pkg_rhino.Node = declaration.GetFirstChild()
		_ = child_2
		for child_2 != 0 {
			var next_2 pkg_rhino.Node = child_2.GetNext()
			_ = next_2
			if child_2.IsName() && child_2.HasChildren() {
				var assign_3 pkg_rhino.Node = pkg_rhino.IR_Assign(child_2.CloneNode_0(), child_2.RemoveFirstChild())
				_ = assign_3
				child_2.ReplaceWith(assign_3)
				assign_3.SetJSDocInfo(declaration.GetJSDocInfo())
			} else if child_2.IsDestructuringLhs() {
				if child_2.HasOneChild() {
					if !(NodeUtil_IsEnhancedFor(declaration.GetParent())) {
						j2gort.CheckStateFmt(false, "DESTRUCTURING_LHS should have two children: %s", j2gort.BoxString(declaration.ToStringTree()))
					}
					child_2.ReplaceWith(child_2.RemoveFirstChild())
				} else {
					var assign_4 pkg_rhino.Node = pkg_rhino.IR_Assign(child_2.RemoveFirstChild(), child_2.RemoveFirstChild())
					_ = assign_4
					child_2.ReplaceWith(assign_4)
					assign_4.SetJSDocInfo(declaration.GetJSDocInfo())
				}
			}
			child_2 = next_2
		}
	}
	this_.compiler.ReportChangeToEnclosingScope(declaration)
}

// Java method: RescopeGlobalSymbolsRewriteCallback.maybeAddPreDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node lhs, boolean isGlobalDeclaration)
func (this_ *RescopeGlobalSymbolsRewriteCallback) priv_RescopeGlobalSymbolsRewriteCallback_maybeAddPreDeclaration_4961(t *NodeTraversal, lhs pkg_rhino.Node, isGlobalDeclaration bool) {
	_ = t
	_ = lhs
	_ = isGlobalDeclaration
	var name string = lhs.GetString()
	_ = name
	if this_.priv_RescopeGlobalSymbolsRewriteCallback_globalVarStaysLocalInDefiningChunk_4959(name, t) {
		var localName string = this_.renameHelper.GetLocalNameInDefiningChunk(name, isGlobalDeclaration)
		_ = localName
		this_.preDeclarations = append(this_.preDeclarations, NewRescopeGlobalSymbolsRewriteCallback_ChunkGlobal(t.GetInput().GetAstRoot(this_.compiler), pkg_rhino.IR_Name(localName).Srcref(lhs)))
	}
}

// Java method: RescopeGlobalSymbolsRewriteCallback.containsNonLocalOrInitializedVars(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node declaration, boolean isGlobalDeclaration)
func (this_ *RescopeGlobalSymbolsRewriteCallback) priv_RescopeGlobalSymbolsRewriteCallback_containsNonLocalOrInitializedVars_4962(t *NodeTraversal, declaration pkg_rhino.Node, isGlobalDeclaration bool) bool {
	_ = t
	_ = declaration
	_ = isGlobalDeclaration
	{
		var child pkg_rhino.Node = declaration.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if (isGlobalDeclaration && child.IsName()) && !(this_.priv_RescopeGlobalSymbolsRewriteCallback_globalVarStaysLocalInDefiningChunk_4959(child.GetString(), t)) {
				return true
			}
			if child.HasChildren() || child.IsDestructuringLhs() {
				return true
			}
		}
	}
	return false
}

// Java method: RescopeGlobalSymbolsRewriteCallback.visitName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RescopeGlobalSymbolsRewriteCallback) priv_RescopeGlobalSymbolsRewriteCallback_visitName_4963(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var name string = n.GetString()
	_ = name
	if parent.IsFunction() && (j2gort.CheckNotNullString(name) == "") {
		return
	}
	if this_.priv_RescopeGlobalSymbolsRewriteCallback_isExternVar_4957(name, t) {
		return
	}
	var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(name))
	_ = var_
	this_.renameHelper.MaybeDisambiguateLocalVar(n, var_)
	if !(var_.IsGlobal() && this_.priv_RescopeGlobalSymbolsRewriteCallback_isCrossChunkName_4956(name)) {
		return
	}
	var currentChunk *JSChunk = t.GetChunk()
	_ = currentChunk
	if this_.optimizeLocalAccess.ShouldWrapReassignableSymbols() {
		j2gort.OrderedMapComputeIfAbsentFn(&this_.chunksUsingCrossChunkSymbols, name, nil, lambdaSingleton_rescope_global_symbols_rewrite_callback_1).
			Add_1(j2gort.PtrToAny(currentChunk))
	}
	if this_.optimizeLocalAccess.ShouldOptimize() && !(j2gort.SetContainsStr(this_.symbolInfo.CrossChunkNamesWithWriteFromOtherChunk(), name)) {
		var definingChunk *JSChunk = var_.GetInput().GetChunk()
		_ = definingChunk
		if currentChunk == definingChunk {
			if this_.priv_RescopeGlobalSymbolsRewriteCallback_globalVarStaysLocalInDefiningChunk_4959(name, t) {
				var localName string = this_.renameHelper.MaybeRenameGlobalInDefiningChunk(n)
				_ = localName
				if NodeUtil_IsLValueWrite(n) {
					this_.priv_RescopeGlobalSymbolsRewriteCallback_addGlobalNamespaceAlias_4965(n, name, localName, currentChunk)
				}
				return
			}
			if !(this_.optimizeLocalAccess.ShouldWrapReassignableSymbols()) || !(j2gort.SetContainsStr(this_.symbolInfo.WrappedCrossChunkNames(), name)) {
				var replacement pkg_rhino.Node = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name(this_.symbolInfo.GlobalSymbolNamespace()), name)
				_ = replacement
				this_.priv_RescopeGlobalSymbolsRewriteCallback_replaceSymbol_4964(n, name, replacement)
				return
			}
		} else if this_.optimizeLocalAccess.ShouldOptimizeAllChunks() && (!(this_.optimizeLocalAccess.ShouldWrapReassignableSymbols()) || !(j2gort.SetContainsStr(this_.symbolInfo.WrappedCrossChunkNames(), name))) {
			if !(j2gort.SetContainsStr(this_.symbolInfo.GlobalNamesWithInnerScopeWriteInDefiningChunk(), name)) && this_.compiler.GetChunkGraph().DependsOn(currentChunk, definingChunk) {
				j2gort.SetAddStr(j2gort.OrderedMapComputeIfAbsentFn(&this_.localAliasesForUnwrappedCrossChunkNames, currentChunk, nil, lambdaSingleton_rescope_global_symbols_rewrite_callback_2), name)
				return
			}
		}
	}
	if !(this_.optimizeLocalAccess.ShouldWrapReassignableSymbols()) {
		var replacement_2 pkg_rhino.Node = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name(this_.symbolInfo.GlobalSymbolNamespace()), name)
		_ = replacement_2
		this_.priv_RescopeGlobalSymbolsRewriteCallback_replaceSymbol_4964(n, name, replacement_2)
		return
	}
	if !(j2gort.SetContainsStr(this_.symbolInfo.WrappedCrossChunkNames(), name)) {
		j2gort.CheckStateFmt(false, "Unexpected unaliased cross-chunk symbol under ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS: %s", j2gort.BoxString(name))
	}
	j2gort.OrderedMapComputeIfAbsentFn(&this_.chunksUsingWrappedReassignableSymbols, name, nil, lambdaSingleton_rescope_global_symbols_rewrite_callback_3).
		Add_1(j2gort.PtrToAny(currentChunk))
	j2gort.OrderedMapComputeIfAbsentFn(&this_.chunksRequiringWrapperInitialization, name, nil, lambdaSingleton_rescope_global_symbols_rewrite_callback_4).
		Add_1(j2gort.PtrToAny(currentChunk))
	var replacement_3 pkg_rhino.Node = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name(name), "_")
	_ = replacement_3
	this_.priv_RescopeGlobalSymbolsRewriteCallback_replaceSymbol_4964(n, name, replacement_3)
}

// Java method: RescopeGlobalSymbolsRewriteCallback.replaceSymbol(com.google.javascript.rhino.Node node, java.lang.String name, com.google.javascript.rhino.Node replacement)
func (this_ *RescopeGlobalSymbolsRewriteCallback) priv_RescopeGlobalSymbolsRewriteCallback_replaceSymbol_4964(node pkg_rhino.Node, name string, replacement pkg_rhino.Node) {
	_ = node
	_ = name
	_ = replacement
	var parent pkg_rhino.Node = node.GetParent()
	_ = parent
	replacement.SrcrefTree(node)
	node.ReplaceWith(replacement)
	this_.compiler.ReportChangeToEnclosingScope(replacement)
	if parent.IsCall() && !(j2gort.SetContainsStr(this_.symbolInfo.MaybeReferencesThis(), name)) {
		parent.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), false)
	}
	this_.compiler.ReportChangeToEnclosingScope(parent)
}

// Java method: RescopeGlobalSymbolsRewriteCallback.addGlobalNamespaceAlias(com.google.javascript.rhino.Node n, java.lang.String name, java.lang.String localName, com.google.javascript.jscomp.JSChunk currentChunk)
func (this_ *RescopeGlobalSymbolsRewriteCallback) priv_RescopeGlobalSymbolsRewriteCallback_addGlobalNamespaceAlias_4965(n pkg_rhino.Node, name string, localName string, currentChunk *JSChunk) {
	_ = n
	_ = name
	_ = localName
	_ = currentChunk
	var stmt pkg_rhino.Node = NodeUtil_GetEnclosingStatement(n)
	_ = stmt
	if stmt == 0 {
		return
	}
	var globalTarget pkg_rhino.Node = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name(this_.symbolInfo.GlobalSymbolNamespace()), name)
	_ = globalTarget
	if this_.optimizeLocalAccess.ShouldWrapReassignableSymbols() && j2gort.SetContainsStr(this_.symbolInfo.WrappedCrossChunkNames(), name) {
		j2gort.OrderedMapComputeIfAbsentFn(&this_.chunksRequiringWrapperInitialization, name, nil, lambdaSingleton_rescope_global_symbols_rewrite_callback_5).
			Add_1(j2gort.PtrToAny(currentChunk))
		globalTarget = pkg_rhino.IR_Getprop_2(globalTarget, "_")
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if parent.IsAssign() && (parent.GetFirstChild() == n) {
		var rhs pkg_rhino.Node = parent.GetSecondChild()
		_ = rhs
		var globalNamespaceAssign pkg_rhino.Node = pkg_rhino.IR_Assign(globalTarget, rhs.Detach())
		_ = globalNamespaceAssign
		globalNamespaceAssign.SrcrefTree(n)
		parent.AddChildToBack(globalNamespaceAssign)
		this_.compiler.ReportChangeToEnclosingScope(parent)
	} else if NodeUtil_IsNameDeclaration(parent) {
		if !(n.HasChildren()) {
			j2gort.CheckStateFmt(false, "Expected initialized declaration for %s", j2gort.BoxNode(n))
		}
		var rhs_2 pkg_rhino.Node = n.RemoveFirstChild()
		_ = rhs_2
		var globalNamespaceAssign_2 pkg_rhino.Node = pkg_rhino.IR_Assign(globalTarget, rhs_2)
		_ = globalNamespaceAssign_2
		globalNamespaceAssign_2.SrcrefTree(n)
		n.AddChildToBack(globalNamespaceAssign_2)
		this_.compiler.ReportChangeToEnclosingScope(n)
	} else if (this_.optimizeLocalAccess.ShouldRewriteCompoundAndIncrementWritesInPlace() && NodeUtil_IsCompoundAssignmentOp(parent)) && (parent.GetFirstChild() == n) {
		var placeholder pkg_rhino.Node = pkg_rhino.IR_Empty()
		_ = placeholder
		parent.ReplaceWith(placeholder)
		var globalNamespaceAssign_3 pkg_rhino.Node = pkg_rhino.IR_Assign(globalTarget.SrcrefTree(n), parent).Srcref(n)
		_ = globalNamespaceAssign_3
		placeholder.ReplaceWith(globalNamespaceAssign_3)
		this_.compiler.ReportChangeToEnclosingScope(globalNamespaceAssign_3)
	} else if this_.optimizeLocalAccess.ShouldRewriteCompoundAndIncrementWritesInPlace() && (parent.IsInc() || parent.IsDec()) {
		if !(parent.GetBooleanProp(pkg_rhino.StaticGet_Node_INCRDECR_PROP())) {
			var placeholder_2 pkg_rhino.Node = pkg_rhino.IR_Empty()
			_ = placeholder_2
			parent.ReplaceWith(placeholder_2)
			var globalNamespaceAssign_4 pkg_rhino.Node = pkg_rhino.IR_Assign(globalTarget.SrcrefTree(n), parent).Srcref(n)
			_ = globalNamespaceAssign_4
			placeholder_2.ReplaceWith(globalNamespaceAssign_4)
			this_.compiler.ReportChangeToEnclosingScope(globalNamespaceAssign_4)
		} else {
			var placeholder_3 pkg_rhino.Node = pkg_rhino.IR_Empty()
			_ = placeholder_3
			parent.ReplaceWith(placeholder_3)
			var globalUpdate pkg_rhino.Node = (func() pkg_rhino.Node {
				if parent.IsInc() {
					return pkg_rhino.IR_Inc(globalTarget, false)
				}
				return pkg_rhino.IR_Dec(globalTarget, false)
			}()).SrcrefTree(n)
			_ = globalUpdate
			var comma pkg_rhino.Node = pkg_rhino.IR_Comma(globalUpdate, parent).Srcref(n)
			_ = comma
			placeholder_3.ReplaceWith(comma)
			this_.compiler.ReportChangeToEnclosingScope(comma)
		}
	} else if !(this_.optimizeLocalAccess.ShouldRewriteCompoundAndIncrementWritesInPlace()) || RescopeGlobalSymbolsRewriteCallback_IsSafeDestructuringWrite(n, stmt) {
		var globalNamespaceAssign_5 pkg_rhino.Node = pkg_rhino.IR_Assign(globalTarget, pkg_rhino.IR_Name(localName))
		_ = globalNamespaceAssign_5
		var aliasStatement pkg_rhino.Node = pkg_rhino.IR_ExprResult(globalNamespaceAssign_5)
		_ = aliasStatement
		aliasStatement.SrcrefTree(n)
		aliasStatement.InsertAfter(stmt)
		this_.compiler.ReportChangeToEnclosingScope(stmt.GetParent())
	} else {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected unsafe write to local alias: "+j2gort.StrVal(n))))
	}
}

// Java static method: RescopeGlobalSymbolsRewriteCallback.isSafeDestructuringWrite(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node stmt)
func RescopeGlobalSymbolsRewriteCallback_IsSafeDestructuringWrite(n pkg_rhino.Node, stmt pkg_rhino.Node) bool {
	_ = n
	_ = stmt
	if !(NodeUtil_IsLhsByDestructuring(n)) || !(NodeUtil_IsStatementBlock(stmt.GetParent())) {
		return false
	}
	var rootTarget pkg_rhino.Node = NodeUtil_GetRootTarget(n)
	_ = rootTarget
	var rootTargetParent pkg_rhino.Node = rootTarget.GetParent()
	_ = rootTargetParent
	if !(rootTargetParent.IsAssign()) || (rootTargetParent.GetParent() != stmt) {
		return false
	}
	return (stmt.IsExprResult() || NodeUtil_IsNameDeclaration(stmt))
}

// Java method: RescopeGlobalSymbolsRewriteCallback.addDeclarations()
func (this_ *RescopeGlobalSymbolsRewriteCallback) AddDeclarations() {
	this_.priv_RescopeGlobalSymbolsRewriteCallback_declareChunkGlobals_4967()
	if this_.optimizeLocalAccess.ShouldOptimizeAllChunks() {
		this_.priv_RescopeGlobalSymbolsRewriteCallback_declareLocalAliasesAndWrappers_4968()
	}
}

// Java method: RescopeGlobalSymbolsRewriteCallback.declareChunkGlobals()
func (this_ *RescopeGlobalSymbolsRewriteCallback) priv_RescopeGlobalSymbolsRewriteCallback_declareChunkGlobals_4967() {
	for _, global := range this_.preDeclarations {
		_ = global
		if global.root.HasChildren() && global.root.GetFirstChild().IsVar() {
			global.root.GetFirstChild().AddChildToBack(global.name)
		} else {
			global.root.AddChildToFront(pkg_rhino.IR_Var_1(global.name).Srcref(global.name))
		}
		this_.compiler.ReportChangeToEnclosingScope(global.root)
	}

}

// Java method: RescopeGlobalSymbolsRewriteCallback.declareLocalAliasesAndWrappers()
func (this_ *RescopeGlobalSymbolsRewriteCallback) priv_RescopeGlobalSymbolsRewriteCallback_declareLocalAliasesAndWrappers_4968() {
	var chunkGraph *JSChunkGraph = this_.compiler.GetChunkGraph()
	_ = chunkGraph
	var chunks j2gort.Ext_java_lang_Iterable = chunkGraph.GetAllChunks()
	_ = chunks
	var defaultRootChunk *JSChunk = j2gort.CastPtr[JSChunk](j2gort.AsIface[j2gort.Ext_java_util_Iterator](chunks.Iterator()).Next())
	_ = defaultRootChunk
	var wrapperAssignmentChunks j2gort.OrderedMap[string, *JSChunk] = j2gort.NewOrderedMap[string, *JSChunk]()
	_ = wrapperAssignmentChunks
	if this_.optimizeLocalAccess.ShouldWrapReassignableSymbols() && !(this_.symbolInfo.WrappedCrossChunkNames().IsEmpty()) {
		{
			cur_32 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.symbolInfo.WrappedCrossChunkNames()))
			for cur_32.HasNext() {
				name := cur_32.NextStr()
				_ = name
				var chunksRequiringWrapper j2gort.Ext_java_util_Set = this_.chunksRequiringWrapperInitialization.GetOrDefault(name, j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
				_ = chunksRequiringWrapper
				var wrapperAssignmentChunk *JSChunk
				_ = wrapperAssignmentChunk
				if chunksRequiringWrapper.IsEmpty() {
					wrapperAssignmentChunk = defaultRootChunk
				} else {
					wrapperAssignmentChunk = chunkGraph.GetDeepestCommonDependencyInclusive_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](chunksRequiringWrapper))
				}
				wrapperAssignmentChunks.Put(name, wrapperAssignmentChunk)
			}
		}
	}
	var cleanupHelperName string = this_.priv_RescopeGlobalSymbolsRewriteCallback_getCleanupHelperPropertyName_4969()
	_ = cleanupHelperName
	var needsCleanupHelper bool = false
	_ = needsCleanupHelper
	{
		cur_33 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunkGraph.GetAllChunks()))
		for cur_33.HasNext() {
			rawElem_34 := cur_33.Next()
			chunk := j2gort.CastPtr[JSChunk](rawElem_34)
			_ = chunk
			if this_.priv_RescopeGlobalSymbolsRewriteCallback_declareLocalAliasesAndWrappersForChunk_4971(chunk, &wrapperAssignmentChunks, cleanupHelperName) {
				needsCleanupHelper = true
			}
		}
	}
	if needsCleanupHelper {
		this_.priv_RescopeGlobalSymbolsRewriteCallback_declareCleanupHelper_4970(cleanupHelperName)
	}
}

// Java method: RescopeGlobalSymbolsRewriteCallback.getCleanupHelperPropertyName()
func (this_ *RescopeGlobalSymbolsRewriteCallback) priv_RescopeGlobalSymbolsRewriteCallback_getCleanupHelperPropertyName_4969() string {
	var name string = "_"
	_ = name
	for this_.chunksUsingCrossChunkSymbols.ContainsKey(name) {
		name = (j2gort.StrValString(name) + "$")
	}
	return name
}

// Java method: RescopeGlobalSymbolsRewriteCallback.declareCleanupHelper(java.lang.String cleanupHelperName)
func (this_ *RescopeGlobalSymbolsRewriteCallback) priv_RescopeGlobalSymbolsRewriteCallback_declareCleanupHelper_4970(cleanupHelperName string) {
	_ = cleanupHelperName
	var script pkg_rhino.Node = this_.compiler.GetNodeForCodeInsertion(nil)
	_ = script
	var globalSymbolNamespace string = this_.symbolInfo.GlobalSymbolNamespace()
	_ = globalSymbolNamespace
	var paramName string
	_ = paramName
	if j2gort.CheckNotNullString(globalSymbolNamespace) == "a" {
		paramName = "b"
	} else {
		paramName = "a"
	}
	var loopVar string
	_ = loopVar
	if !(j2gort.CheckNotNullString(paramName) == "b") && !(j2gort.CheckNotNullString(globalSymbolNamespace) == "b") {
		loopVar = "b"
	} else {
		loopVar = "c"
	}
	var splitCall pkg_rhino.Node = pkg_rhino.IR_Call(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name(paramName), "split"), []pkg_rhino.Node{pkg_rhino.IR_String(",")})
	_ = splitCall
	var deleteStmt pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Delprop(pkg_rhino.IR_Getelem(pkg_rhino.IR_Name(globalSymbolNamespace), pkg_rhino.IR_Name(loopVar))))
	_ = deleteStmt
	var forOf pkg_rhino.Node = pkg_rhino.NewNode_4(pkg_rhino.StaticGet_Token_FOR_OF(), pkg_rhino.IR_Let_1(pkg_rhino.IR_Name(loopVar)), splitCall, pkg_rhino.IR_Block_1_Node(deleteStmt))
	_ = forOf
	var helperFn pkg_rhino.Node = pkg_rhino.IR_Function(pkg_rhino.IR_Name(""), pkg_rhino.IR_ParamList([]pkg_rhino.Node{pkg_rhino.IR_Name(paramName)}), pkg_rhino.IR_Block_1_Node(forOf))
	_ = helperFn
	var helperDecl pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name(globalSymbolNamespace), cleanupHelperName), helperFn))
	_ = helperDecl
	RescopeGlobalSymbolsRewriteCallback_AddDeclaration(script, script.GetFirstChild(), helperDecl)
	NodeUtil_MarkNewScopesChanged_2_Node_AbstractCompiler(helperFn, this_.compiler)
	NodeUtil_AddFeatureToScript(script, pkg_parser.StaticGet_FeatureSet_Feature_FOR_OF(), this_.compiler)
	NodeUtil_AddFeatureToScript(script, pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), this_.compiler)
	this_.compiler.ReportChangeToEnclosingScope(script)
}

// Java method: RescopeGlobalSymbolsRewriteCallback.declareLocalAliasesAndWrappersForChunk(com.google.javascript.jscomp.JSChunk chunk, java.util.Map<java.lang.String,com.google.javascript.jscomp.JSChunk> wrapperAssignmentChunks, java.lang.String cleanupHelperName)
func (this_ *RescopeGlobalSymbolsRewriteCallback) priv_RescopeGlobalSymbolsRewriteCallback_declareLocalAliasesAndWrappersForChunk_4971(chunk *JSChunk, wrapperAssignmentChunks *j2gort.OrderedMap[string, *JSChunk], cleanupHelperName string) bool {
	_ = chunk
	_ = wrapperAssignmentChunks
	_ = cleanupHelperName
	var inputs []*CompilerInput = j2gort.ImmListToSliceNullable[*CompilerInput](chunk.GetInputs())
	_ = inputs
	if len(inputs) == 0 {
		return false
	}
	var script pkg_rhino.Node = inputs[int32(0)].GetAstRoot(this_.compiler)
	_ = script
	var insertionPoint pkg_rhino.Node = script.GetFirstChild()
	_ = insertionPoint
	var changed bool = false
	_ = changed
	var usedCleanupHelper bool = false
	_ = usedCleanupHelper
	var wrappedReassignableLocalAliases j2gort.Ext_java_util_Set = j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
	_ = wrappedReassignableLocalAliases
	if this_.optimizeLocalAccess.ShouldWrapReassignableSymbols() {
		wrappedReassignableLocalAliases = j2gort.NewExt_java_util_LinkedHashSet()
		{
			cur_35 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.symbolInfo.WrappedCrossChunkNames()))
			for cur_35.HasNext() {
				name := cur_35.NextStr()
				_ = name
				var assignmentChunk *JSChunk = wrapperAssignmentChunks.Get(name)
				_ = assignmentChunk
				var usingChunks j2gort.Ext_java_util_Set = this_.chunksUsingWrappedReassignableSymbols.Get(name)
				_ = usingChunks
				if j2gort.CheckNotNullPtr(chunk) == assignmentChunk {
					if (usingChunks != nil) && usingChunks.Contains(j2gort.PtrToAny(chunk)) {
						RescopeGlobalSymbolsRewriteCallback_AddDeclaration(script, insertionPoint, pkg_rhino.IR_Var_2(pkg_rhino.IR_Name(name), pkg_rhino.IR_Assign(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name(this_.symbolInfo.GlobalSymbolNamespace()), name), pkg_rhino.IR_Objectlit([]pkg_rhino.Node{}))))
					} else {
						RescopeGlobalSymbolsRewriteCallback_AddDeclaration(script, insertionPoint, pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name(this_.symbolInfo.GlobalSymbolNamespace()), name), pkg_rhino.IR_Objectlit([]pkg_rhino.Node{}))))
					}
					changed = true
				} else if (usingChunks != nil) && usingChunks.Contains(j2gort.PtrToAny(chunk)) {
					j2gort.SetAddStr(wrappedReassignableLocalAliases, name)
				}
			}
		}
	}
	var allLocalAliases j2gort.Ext_java_lang_Iterable = j2gort.Ext_com_google_common_collect_Iterables_Concat_2(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.localAliasesForUnwrappedCrossChunkNames.GetOrDefault(chunk, j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())), j2gort.AsIface[j2gort.Ext_java_lang_Iterable](wrappedReassignableLocalAliases))
	_ = allLocalAliases
	if !(j2gort.Ext_com_google_common_collect_Iterables_IsEmpty(allLocalAliases)) {
		var aliasesToKeep j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
		_ = aliasesToKeep
		var aliasesToDelete j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
		_ = aliasesToDelete
		{
			cur_36 := j2gort.NewIterCursor(j2gort.CheckNotNull(allLocalAliases))
			for cur_36.HasNext() {
				name_2 := cur_36.NextStr()
				_ = name_2
				if this_.optimizeLocalAccess.ShouldDeleteGlobalsOnLastAccess() && this_.priv_RescopeGlobalSymbolsRewriteCallback_canDeleteFromGlobalNamespace_4972(chunk, name_2) {
					aliasesToDelete.Add_1(j2gort.BoxString(name_2))
				} else {
					aliasesToKeep.Add_1(j2gort.BoxString(name_2))
				}
			}
		}
		var localAliases j2gort.Ext_java_lang_Iterable = j2gort.Ext_com_google_common_collect_Iterables_Concat_2(aliasesToKeep, aliasesToDelete)
		_ = localAliases
		var outputFeatures *pkg_parser.FeatureSet = this_.compiler.GetOptions().GetOutputFeatureSet()
		_ = outputFeatures
		if (outputFeatures.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_DESTRUCTURING()) && outputFeatures.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_FOR_OF())) && outputFeatures.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS()) {
			var objectPattern pkg_rhino.Node = pkg_rhino.IR_ObjectPattern([]pkg_rhino.Node{})
			_ = objectPattern
			{
				cur_37 := j2gort.NewIterCursor(j2gort.CheckNotNull(localAliases))
				for cur_37.HasNext() {
					name_3 := cur_37.NextStr()
					_ = name_3
					var stringKey pkg_rhino.Node = pkg_rhino.IR_StringKey_2(name_3, pkg_rhino.IR_Name(name_3))
					_ = stringKey
					stringKey.SetShorthandProperty(true)
					objectPattern.AddChildToBack(stringKey)
				}
			}
			RescopeGlobalSymbolsRewriteCallback_AddDeclaration(script, insertionPoint, pkg_rhino.IR_Var_2(objectPattern, pkg_rhino.IR_Name(this_.symbolInfo.GlobalSymbolNamespace())))
			NodeUtil_AddFeatureToScript(script, pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_DESTRUCTURING(), this_.compiler)
			if !(aliasesToDelete.IsEmpty()) {
				var cleanupCall pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Call(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name(this_.symbolInfo.GlobalSymbolNamespace()), cleanupHelperName), []pkg_rhino.Node{pkg_rhino.IR_String(j2gort.Ext_com_google_common_base_Joiner_On_1_uint16(uint16(44)).Join_1_Ext_java_lang_Iterable(aliasesToDelete))}))
				_ = cleanupCall
				RescopeGlobalSymbolsRewriteCallback_AddDeclaration(script, insertionPoint, cleanupCall)
				usedCleanupHelper = true
			}
		} else {
			{
				cur_38 := j2gort.NewIterCursor(aliasesToKeep)
				for cur_38.HasNext() {
					name_4 := cur_38.NextStr()
					_ = name_4
					RescopeGlobalSymbolsRewriteCallback_AddDeclaration(script, insertionPoint, pkg_rhino.IR_Var_2(pkg_rhino.IR_Name(name_4), pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name(this_.symbolInfo.GlobalSymbolNamespace()), name_4)))
				}
			}
			{
				cur_39 := j2gort.NewIterCursor(aliasesToDelete)
				for cur_39.HasNext() {
					name_5 := cur_39.NextStr()
					_ = name_5
					RescopeGlobalSymbolsRewriteCallback_AddDeclaration(script, insertionPoint, pkg_rhino.IR_Var_2(pkg_rhino.IR_Name(name_5), pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name(this_.symbolInfo.GlobalSymbolNamespace()), name_5)))
					RescopeGlobalSymbolsRewriteCallback_AddDeclaration(script, insertionPoint, pkg_rhino.IR_ExprResult(pkg_rhino.IR_Delprop(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name(this_.symbolInfo.GlobalSymbolNamespace()), name_5))))
				}
			}
		}
		changed = true
	}
	if changed {
		this_.compiler.ReportChangeToEnclosingScope(script)
	}
	return usedCleanupHelper
}

// Java method: RescopeGlobalSymbolsRewriteCallback.canDeleteFromGlobalNamespace(com.google.javascript.jscomp.JSChunk currentChunk, java.lang.String name)
func (this_ *RescopeGlobalSymbolsRewriteCallback) priv_RescopeGlobalSymbolsRewriteCallback_canDeleteFromGlobalNamespace_4972(currentChunk *JSChunk, name string) bool {
	_ = currentChunk
	_ = name
	if this_.compiler.GetCodingConvention().IsExported_2(name, false) {
		return false
	}
	if this_.renameHelper.HasUnwrappedLocalNameInDefiningChunk(name) && j2gort.SetContainsStr(this_.symbolInfo.GlobalNamesWithInnerScopeWriteInDefiningChunk(), name) {
		return false
	}
	var usingChunks j2gort.Ext_java_util_Set = this_.chunksUsingCrossChunkSymbols.Get(name)
	_ = usingChunks
	if (usingChunks == nil) || !(usingChunks.Contains(j2gort.PtrToAny(currentChunk))) {
		return false
	}
	var chunkGraph *JSChunkGraph = this_.compiler.GetChunkGraph()
	_ = chunkGraph
	{
		cur_40 := j2gort.NewIterCursor(j2gort.CheckNotNull(usingChunks))
		for cur_40.HasNext() {
			rawElem_41 := cur_40.Next()
			otherChunk := j2gort.CastPtr[JSChunk](rawElem_41)
			_ = otherChunk
			if !(j2gort.CheckNotNullPtr(otherChunk) == currentChunk) && !(chunkGraph.DependsOn(currentChunk, otherChunk)) {
				return false
			}
		}
	}
	return true
}

// Java static method: RescopeGlobalSymbolsRewriteCallback.addDeclaration(com.google.javascript.rhino.Node script, com.google.javascript.rhino.Node insertionPoint, com.google.javascript.rhino.Node decl)
func RescopeGlobalSymbolsRewriteCallback_AddDeclaration(script pkg_rhino.Node, insertionPoint pkg_rhino.Node, decl pkg_rhino.Node) {
	_ = script
	_ = insertionPoint
	_ = decl
	decl.SrcrefTree(script)
	if insertionPoint == 0 {
		script.AddChildToBack(decl)
	} else {
		decl.InsertBefore(insertionPoint)
	}
}

// Java class: com.google.javascript.jscomp.RescopeGlobalSymbolsRewriteCallback.ChunkGlobal
type RescopeGlobalSymbolsRewriteCallback_ChunkGlobal struct {
	root    pkg_rhino.Node `gson:"root"`
	name    pkg_rhino.Node `gson:"name"`
	idHash_ atomic.Int32
}

func (this_ *RescopeGlobalSymbolsRewriteCallback_ChunkGlobal) BindSelf_(s any) {
	_ = s
}
func (this_ *RescopeGlobalSymbolsRewriteCallback_ChunkGlobal) IsJavaNil_() bool { return this_ == nil }
func (*RescopeGlobalSymbolsRewriteCallback_ChunkGlobal) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RescopeGlobalSymbolsRewriteCallback$ChunkGlobal", "ChunkGlobal"
}
func (this_ *RescopeGlobalSymbolsRewriteCallback_ChunkGlobal) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RescopeGlobalSymbolsRewriteCallback_ChunkGlobal) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RescopeGlobalSymbolsRewriteCallback.ChunkGlobal", "com.google.javascript.jscomp.RescopeGlobalSymbolsRewriteCallback$ChunkGlobal", func(v_ any) bool {
		return j2gort.InstanceOf[*RescopeGlobalSymbolsRewriteCallback_ChunkGlobal](v_)
	})
}

func (this_ *RescopeGlobalSymbolsRewriteCallback_ChunkGlobal) initFields_RescopeGlobalSymbolsRewriteCallback_ChunkGlobal() {
	_ = this_
}

// Java constructor: ChunkGlobal(com.google.javascript.rhino.Node root, com.google.javascript.rhino.Node name)
func NewRescopeGlobalSymbolsRewriteCallback_ChunkGlobal(root pkg_rhino.Node, name pkg_rhino.Node) *RescopeGlobalSymbolsRewriteCallback_ChunkGlobal {
	this_ := &RescopeGlobalSymbolsRewriteCallback_ChunkGlobal{}
	InitRescopeGlobalSymbolsRewriteCallback_ChunkGlobal(this_, root, name)
	return this_
}

func InitRescopeGlobalSymbolsRewriteCallback_ChunkGlobal(this_ *RescopeGlobalSymbolsRewriteCallback_ChunkGlobal, root pkg_rhino.Node, name pkg_rhino.Node) {
	_ = root
	_ = name
	this_.initFields_RescopeGlobalSymbolsRewriteCallback_ChunkGlobal()
	this_.root = root
	this_.name = name
}
