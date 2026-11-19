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
// Port of com/google/javascript/jscomp/RescopeGlobalSymbols.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.RescopeGlobalSymbols
type RescopeGlobalSymbols struct {
	compiler              AbstractCompiler `gson:"compiler"`
	globalSymbolNamespace string           `gson:"globalSymbolNamespace"`
	idHash_               atomic.Int32
	addExtern             bool                                `gson:"addExtern"`
	assumeCrossChunkNames bool                                `gson:"assumeCrossChunkNames"`
	optimizeLocalAccess   CompilerOptions_OptimizeLocalAccess `gson:"optimizeLocalAccess"`
}

func (this_ *RescopeGlobalSymbols) BindSelf_(s any) {
	_ = s
}
func (this_ *RescopeGlobalSymbols) IsJavaNil_() bool { return this_ == nil }
func (*RescopeGlobalSymbols) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RescopeGlobalSymbols", "RescopeGlobalSymbols"
}
func (this_ *RescopeGlobalSymbols) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RescopeGlobalSymbols) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RescopeGlobalSymbols", "com.google.javascript.jscomp.RescopeGlobalSymbols", func(v_ any) bool {
		return j2gort.InstanceOf[*RescopeGlobalSymbols](v_)
	})
}
func (*RescopeGlobalSymbols) Implements_jscomp_CompilerPass() {}

func (this_ *RescopeGlobalSymbols) initFields_RescopeGlobalSymbols() {
	_ = this_
	this_.globalSymbolNamespace = j2gort.NullString
}

// Java constructor: RescopeGlobalSymbols(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String globalSymbolNamespace, boolean assumeCrossChunkNames, com.google.javascript.jscomp.CompilerOptions.OptimizeLocalAccess optimizeLocalAccess)
func NewRescopeGlobalSymbols_4(compiler AbstractCompiler, globalSymbolNamespace string, assumeCrossChunkNames bool, optimizeLocalAccess CompilerOptions_OptimizeLocalAccess) *RescopeGlobalSymbols {
	this_ := &RescopeGlobalSymbols{}
	InitRescopeGlobalSymbols_4(this_, compiler, globalSymbolNamespace, assumeCrossChunkNames, optimizeLocalAccess)
	return this_
}

func InitRescopeGlobalSymbols_4(this_ *RescopeGlobalSymbols, compiler AbstractCompiler, globalSymbolNamespace string, assumeCrossChunkNames bool, optimizeLocalAccess CompilerOptions_OptimizeLocalAccess) {
	_ = compiler
	_ = globalSymbolNamespace
	_ = assumeCrossChunkNames
	_ = optimizeLocalAccess
	InitRescopeGlobalSymbols_5(this_, compiler, globalSymbolNamespace, true, assumeCrossChunkNames, optimizeLocalAccess)
}

// Java constructor: RescopeGlobalSymbols(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String globalSymbolNamespace, boolean addExtern, boolean assumeCrossChunkNames, com.google.javascript.jscomp.CompilerOptions.OptimizeLocalAccess optimizeLocalAccess)
func NewRescopeGlobalSymbols_5(compiler AbstractCompiler, globalSymbolNamespace string, addExtern bool, assumeCrossChunkNames bool, optimizeLocalAccess CompilerOptions_OptimizeLocalAccess) *RescopeGlobalSymbols {
	this_ := &RescopeGlobalSymbols{}
	InitRescopeGlobalSymbols_5(this_, compiler, globalSymbolNamespace, addExtern, assumeCrossChunkNames, optimizeLocalAccess)
	return this_
}

func InitRescopeGlobalSymbols_5(this_ *RescopeGlobalSymbols, compiler AbstractCompiler, globalSymbolNamespace string, addExtern bool, assumeCrossChunkNames bool, optimizeLocalAccess CompilerOptions_OptimizeLocalAccess) {
	_ = compiler
	_ = globalSymbolNamespace
	_ = addExtern
	_ = assumeCrossChunkNames
	_ = optimizeLocalAccess
	this_.initFields_RescopeGlobalSymbols()
	this_.compiler = compiler
	this_.globalSymbolNamespace = globalSymbolNamespace
	this_.addExtern = addExtern
	this_.assumeCrossChunkNames = assumeCrossChunkNames
	this_.optimizeLocalAccess = optimizeLocalAccess
}

// Java method: RescopeGlobalSymbols.addExternForGlobalSymbolNamespace()
func (this_ *RescopeGlobalSymbols) priv_RescopeGlobalSymbols_addExternForGlobalSymbolNamespace_4922() {
	var varNode pkg_rhino.Node = pkg_rhino.IR_Var_1(pkg_rhino.IR_Name(this_.globalSymbolNamespace))
	_ = varNode
	var input *CompilerInput = this_.compiler.GetSynthesizedExternsInput()
	_ = input
	input.GetAstRoot(this_.compiler).AddChildToBack(varNode)
	this_.compiler.ReportChangeToEnclosingScope(varNode)
}

// Java method: RescopeGlobalSymbols.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RescopeGlobalSymbols) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if this_.addExtern {
		this_.priv_RescopeGlobalSymbols_addExternForGlobalSymbolNamespace_4922()
	}
	NodeTraversal_Traverse(this_.compiler, root, NewRescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback(this_.compiler))
	var val_1 RescopeGlobalSymbolsRenameHelper_HelperBuilder
	InitRescopeGlobalSymbolsRenameHelper_HelperBuilder(&val_1, this_.compiler, this_.globalSymbolNamespace, this_.optimizeLocalAccess)
	var renameHelperBuilder *RescopeGlobalSymbolsRenameHelper_HelperBuilder = &val_1
	_ = renameHelperBuilder
	var findCrossChunkNames *RescopeGlobalSymbols_FindCrossChunkNamesCallback = NewRescopeGlobalSymbols_FindCrossChunkNamesCallback(this_.compiler, this_.assumeCrossChunkNames, this_.optimizeLocalAccess.ShouldOptimize(), renameHelperBuilder)
	_ = findCrossChunkNames
	var findNamesReferencingThis *RescopeGlobalSymbols_FindNamesReferencingThis = NewRescopeGlobalSymbols_FindNamesReferencingThis()
	_ = findNamesReferencingThis
	CombinedCompilerPass_Traverse(this_.compiler, root, j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(findCrossChunkNames, findNamesReferencingThis)))
	var externNames j2gort.Ext_com_google_common_collect_ImmutableSet = NodeUtil_CollectExternVariableNames(this_.compiler, externs)
	_ = externNames
	var symbolInfo *RescopeGlobalSymbolsInfo = NewRescopeGlobalSymbolsInfo_9(this_.globalSymbolNamespace, j2gort.AsIface[j2gort.Ext_java_util_Set](externNames), findCrossChunkNames.crossChunkNames, findCrossChunkNames.crossChunkNamesWithWriteFromOtherChunk, findCrossChunkNames.crossChunkNamesUsedInNonDependentChunk, findCrossChunkNames.globalNamesWithReadInDefiningChunk, findCrossChunkNames.globalNamesWithInnerScopeWriteInDefiningChunk, findCrossChunkNames.globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk, findNamesReferencingThis.maybeReferencesThis)
	_ = symbolInfo
	var renameHelper *RescopeGlobalSymbolsRenameHelper = renameHelperBuilder.Build(symbolInfo)
	_ = renameHelper
	var val_2 RescopeGlobalSymbolsRewriteCallback
	InitRescopeGlobalSymbolsRewriteCallback(&val_2, this_.compiler, this_.optimizeLocalAccess, symbolInfo, renameHelper)
	var rewriteScope *RescopeGlobalSymbolsRewriteCallback = &val_2
	_ = rewriteScope
	NodeTraversal_Traverse(this_.compiler, root, rewriteScope)
	NodeTraversal_Traverse(this_.compiler, root, NewRescopeGlobalSymbols_RemoveGlobalVarCallback(this_.compiler, this_.optimizeLocalAccess.ShouldPreserveLocalVarDeclarations()))
	rewriteScope.AddDeclarations()
}

// Java class: com.google.javascript.jscomp.RescopeGlobalSymbols.RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback
type RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback struct {
	NodeTraversal_AbstractShallowStatementCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback) IsJavaNil_() bool {
	return this_ == nil
}
func (*RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RescopeGlobalSymbols$RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback", "RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback"
}
func (this_ *RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractShallowStatementCallbackBase.IdentityHash_()
}
func (*RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback) HasIdentityEquals_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RescopeGlobalSymbols.RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback", "com.google.javascript.jscomp.RescopeGlobalSymbols$RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback](v_)
	})
}
func (*RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback) Implements_jscomp_NodeTraversal_AbstractShallowStatementCallback() {
}
func (*RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback) initFields_RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback() {
	_ = this_
}

// Java constructor: RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback(compiler AbstractCompiler) *RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback {
	this_ := &RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback{}
	InitRescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback(this_, compiler)
	return this_
}

func InitRescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback(this_ *RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractShallowStatementCallback(&this_.NodeTraversal_AbstractShallowStatementCallbackBase)
	this_.initFields_RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback()
	this_.compiler = compiler
}

// Java method: RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(t.InGlobalScope()) {
		return
	}
	if !(NodeUtil_IsFunctionDeclaration(n)) && !(NodeUtil_IsClassDeclaration(n)) {
		return
	}
	var nameNode pkg_rhino.Node = NodeUtil_GetNameNode(n)
	_ = nameNode
	var name string = nameNode.GetString()
	_ = name
	if n.IsClass() {
		nameNode.ReplaceWith(pkg_rhino.IR_Empty().Srcref(nameNode))
	} else {
		nameNode.SetString("")
		this_.compiler.ReportChangeToEnclosingScope(nameNode)
	}
	var prev pkg_rhino.Node = n.GetPrevious()
	_ = prev
	n.Detach()
	var var_ pkg_rhino.Node = NodeUtil_NewVarNode_2_string_Node(name, n)
	_ = var_
	if prev == 0 {
		parent.AddChildToFront(var_)
	} else {
		var_.InsertAfter(prev)
	}
	this_.compiler.ReportChangeToEnclosingScope(parent)
}

// Java class: com.google.javascript.jscomp.RescopeGlobalSymbols.FindCrossChunkNamesCallback
type RescopeGlobalSymbols_FindCrossChunkNamesCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler                                               AbstractCompiler                                `gson:"compiler"`
	crossChunkNames                                        j2gort.Ext_java_util_Set                        `gson:"crossChunkNames"`
	crossChunkNamesWithWriteFromOtherChunk                 j2gort.Ext_java_util_Set                        `gson:"crossChunkNamesWithWriteFromOtherChunk"`
	crossChunkNamesUsedInNonDependentChunk                 j2gort.Ext_java_util_Set                        `gson:"crossChunkNamesUsedInNonDependentChunk"`
	globalNamesWithReadInDefiningChunk                     j2gort.Ext_java_util_Set                        `gson:"globalNamesWithReadInDefiningChunk"`
	globalNamesWithInnerScopeWriteInDefiningChunk          j2gort.Ext_java_util_Set                        `gson:"globalNamesWithInnerScopeWriteInDefiningChunk"`
	globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk j2gort.Ext_java_util_Set                        `gson:"globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk"`
	renameHelperBuilder                                    *RescopeGlobalSymbolsRenameHelper_HelperBuilder `gson:"renameHelperBuilder"`
	assumeCrossChunkNames                                  bool                                            `gson:"assumeCrossChunkNames"`
	trackLocalAccessSets                                   bool                                            `gson:"trackLocalAccessSets"`
}

func (this_ *RescopeGlobalSymbols_FindCrossChunkNamesCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *RescopeGlobalSymbols_FindCrossChunkNamesCallback) IsJavaNil_() bool { return this_ == nil }
func (*RescopeGlobalSymbols_FindCrossChunkNamesCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RescopeGlobalSymbols$FindCrossChunkNamesCallback", "FindCrossChunkNamesCallback"
}
func (this_ *RescopeGlobalSymbols_FindCrossChunkNamesCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*RescopeGlobalSymbols_FindCrossChunkNamesCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RescopeGlobalSymbols.FindCrossChunkNamesCallback", "com.google.javascript.jscomp.RescopeGlobalSymbols$FindCrossChunkNamesCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*RescopeGlobalSymbols_FindCrossChunkNamesCallback](v_)
	})
}
func (*RescopeGlobalSymbols_FindCrossChunkNamesCallback) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*RescopeGlobalSymbols_FindCrossChunkNamesCallback) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *RescopeGlobalSymbols_FindCrossChunkNamesCallback) initFields_RescopeGlobalSymbols_FindCrossChunkNamesCallback() {
	_ = this_
	this_.crossChunkNames = j2gort.NewExt_java_util_LinkedHashSet()
	this_.crossChunkNamesWithWriteFromOtherChunk = j2gort.NewExt_java_util_LinkedHashSet()
	this_.crossChunkNamesUsedInNonDependentChunk = j2gort.NewExt_java_util_LinkedHashSet()
	this_.globalNamesWithReadInDefiningChunk = j2gort.NewExt_java_util_LinkedHashSet()
	this_.globalNamesWithInnerScopeWriteInDefiningChunk = j2gort.NewExt_java_util_LinkedHashSet()
	this_.globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk = j2gort.NewExt_java_util_LinkedHashSet()
}

// Java constructor: FindCrossChunkNamesCallback(com.google.javascript.jscomp.AbstractCompiler compiler, boolean assumeCrossChunkNames, boolean trackLocalAccessSets, com.google.javascript.jscomp.RescopeGlobalSymbolsRenameHelper.HelperBuilder renameHelperBuilder)
func NewRescopeGlobalSymbols_FindCrossChunkNamesCallback(compiler AbstractCompiler, assumeCrossChunkNames bool, trackLocalAccessSets bool, renameHelperBuilder *RescopeGlobalSymbolsRenameHelper_HelperBuilder) *RescopeGlobalSymbols_FindCrossChunkNamesCallback {
	this_ := &RescopeGlobalSymbols_FindCrossChunkNamesCallback{}
	InitRescopeGlobalSymbols_FindCrossChunkNamesCallback(this_, compiler, assumeCrossChunkNames, trackLocalAccessSets, renameHelperBuilder)
	return this_
}

func InitRescopeGlobalSymbols_FindCrossChunkNamesCallback(this_ *RescopeGlobalSymbols_FindCrossChunkNamesCallback, compiler AbstractCompiler, assumeCrossChunkNames bool, trackLocalAccessSets bool, renameHelperBuilder *RescopeGlobalSymbolsRenameHelper_HelperBuilder) {
	_ = compiler
	_ = assumeCrossChunkNames
	_ = trackLocalAccessSets
	_ = renameHelperBuilder
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_RescopeGlobalSymbols_FindCrossChunkNamesCallback()
	this_.compiler = compiler
	this_.assumeCrossChunkNames = assumeCrossChunkNames
	this_.trackLocalAccessSets = trackLocalAccessSets
	this_.renameHelperBuilder = renameHelperBuilder
}

// Java method: FindCrossChunkNamesCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RescopeGlobalSymbols_FindCrossChunkNamesCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsName() {
		var name string = n.GetString()
		_ = name
		if (j2gort.CheckNotNullString(name) == "") || (!(this_.trackLocalAccessSets) && j2gort.SetContainsStr(this_.crossChunkNames, name)) {
			return
		}
		this_.renameHelperBuilder.RecordExistingName(name)
		var s *Scope = t.GetScope()
		_ = s
		var v *Var = j2gort.CastPtr[Var](s.GetVar(name))
		_ = v
		if (v == nil) || !(v.IsGlobal()) {
			return
		}
		if this_.assumeCrossChunkNames || this_.compiler.GetCodingConvention().IsExported_2(name, false) {
			j2gort.SetAddStr(this_.crossChunkNames, name)
		}
		var input *CompilerInput = v.GetInput()
		_ = input
		if ((this_.trackLocalAccessSets && !(t.InGlobalScope())) && ((input == nil) || (input.GetChunk() == t.GetChunk()))) && NodeUtil_IsLValueWrite(n) {
			j2gort.SetAddStr(this_.globalNamesWithInnerScopeWriteInDefiningChunk, name)
		}
		if input == nil {
			j2gort.SetAddStr(this_.crossChunkNames, name)
			return
		}
		var chunk *JSChunk = input.GetChunk()
		_ = chunk
		if chunk != t.GetChunk() {
			j2gort.SetAddStr(this_.crossChunkNames, name)
			if this_.trackLocalAccessSets {
				if NodeUtil_IsLValueWrite(n) {
					j2gort.SetAddStr(this_.crossChunkNamesWithWriteFromOtherChunk, name)
				}
				if !(this_.compiler.GetChunkGraph().DependsOn(t.GetChunk(), chunk)) {
					j2gort.SetAddStr(this_.crossChunkNamesUsedInNonDependentChunk, name)
				}
			}
		} else if this_.trackLocalAccessSets {
			if RescopeGlobalSymbols_FindCrossChunkNamesCallback_IsReadReference(n) {
				j2gort.SetAddStr(this_.globalNamesWithReadInDefiningChunk, name)
			}
			if NodeUtil_IsLValueWrite(n) && !(RescopeGlobalSymbols_FindCrossChunkNamesCallback_IsWriteSafeForLocalAliasMirroring(n, parent)) {
				j2gort.SetAddStr(this_.globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk, name)
			}
		}
	}
}

// Java static method: FindCrossChunkNamesCallback.isReadReference(com.google.javascript.rhino.Node n)
func RescopeGlobalSymbols_FindCrossChunkNamesCallback_IsReadReference(n pkg_rhino.Node) bool {
	_ = n
	if !(NodeUtil_IsLValue(n)) {
		return true
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return ((NodeUtil_IsCompoundAssignmentOp(parent) || parent.IsInc()) || parent.IsDec())
}

// Java static method: FindCrossChunkNamesCallback.isWriteSafeForLocalAliasMirroring(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func RescopeGlobalSymbols_FindCrossChunkNamesCallback_IsWriteSafeForLocalAliasMirroring(n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = n
	_ = parent
	if ((NodeUtil_IsAssignmentOp(parent) && (parent.GetFirstChild() == n)) || parent.IsInc()) || parent.IsDec() {
		return true
	}
	if NodeUtil_IsNameDeclaration(parent) {
		return !(NodeUtil_IsEnhancedFor(parent.GetParent()))
	}
	if NodeUtil_IsLhsByDestructuring(n) {
		var rootTarget pkg_rhino.Node = NodeUtil_GetRootTarget(n)
		_ = rootTarget
		var rootTargetParent pkg_rhino.Node = rootTarget.GetParent()
		_ = rootTargetParent
		if rootTargetParent.IsDestructuringLhs() {
			var declaration pkg_rhino.Node = rootTargetParent.GetParent()
			_ = declaration
			return NodeUtil_IsStatementBlock(declaration.GetParent())
		}
		if rootTargetParent.IsAssign() && (rootTargetParent.GetFirstChild() == rootTarget) {
			var exprResult pkg_rhino.Node = rootTargetParent.GetParent()
			_ = exprResult
			return (exprResult.IsExprResult() && NodeUtil_IsStatementBlock(exprResult.GetParent()))
		}
	}
	return false
}

// Java class: com.google.javascript.jscomp.RescopeGlobalSymbols.FindNamesReferencingThis
type RescopeGlobalSymbols_FindNamesReferencingThis struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	maybeReferencesThis j2gort.Ext_java_util_Set `gson:"maybeReferencesThis"`
}

func (this_ *RescopeGlobalSymbols_FindNamesReferencingThis) BindSelf_(s any) {
	_ = s
}
func (this_ *RescopeGlobalSymbols_FindNamesReferencingThis) IsJavaNil_() bool { return this_ == nil }
func (*RescopeGlobalSymbols_FindNamesReferencingThis) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RescopeGlobalSymbols$FindNamesReferencingThis", "FindNamesReferencingThis"
}
func (this_ *RescopeGlobalSymbols_FindNamesReferencingThis) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*RescopeGlobalSymbols_FindNamesReferencingThis) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RescopeGlobalSymbols.FindNamesReferencingThis", "com.google.javascript.jscomp.RescopeGlobalSymbols$FindNamesReferencingThis", func(v_ any) bool {
		return j2gort.InstanceOf[*RescopeGlobalSymbols_FindNamesReferencingThis](v_)
	})
}
func (*RescopeGlobalSymbols_FindNamesReferencingThis) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*RescopeGlobalSymbols_FindNamesReferencingThis) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *RescopeGlobalSymbols_FindNamesReferencingThis) initFields_RescopeGlobalSymbols_FindNamesReferencingThis() {
	_ = this_
	this_.maybeReferencesThis = j2gort.NewExt_java_util_LinkedHashSet()
}

// Java constructor: FindNamesReferencingThis()
func NewRescopeGlobalSymbols_FindNamesReferencingThis() *RescopeGlobalSymbols_FindNamesReferencingThis {
	this_ := &RescopeGlobalSymbols_FindNamesReferencingThis{}
	InitRescopeGlobalSymbols_FindNamesReferencingThis(this_)
	return this_
}

func InitRescopeGlobalSymbols_FindNamesReferencingThis(this_ *RescopeGlobalSymbols_FindNamesReferencingThis) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_RescopeGlobalSymbols_FindNamesReferencingThis()
}

// Java method: FindNamesReferencingThis.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RescopeGlobalSymbols_FindNamesReferencingThis) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsName() {
		var name string = n.GetString()
		_ = name
		if j2gort.CheckNotNullString(name) == "" {
			return
		}
		var value pkg_rhino.Node = pkg_rhino.Node(0)
		_ = value
		if parent.IsAssign() && (n == parent.GetFirstChild()) {
			value = parent.GetLastChild()
		} else if NodeUtil_IsNameDeclaration(parent) {
			value = n.GetFirstChild()
		} else if parent.IsFunction() {
			value = parent
		}
		if (value == 0) && !(NodeUtil_IsLhsByDestructuring(n)) {
			return
		}
		if j2gort.SetContainsStr(this_.maybeReferencesThis, name) {
			return
		}
		var s *Scope = t.GetScope()
		_ = s
		var v *Var = j2gort.CastPtr[Var](s.GetVar(name))
		_ = v
		if (v == nil) || !(v.IsGlobal()) {
			return
		}
		if ((value == 0) || !(value.IsFunction())) || NodeUtil_ReferencesOwnReceiver(value) {
			j2gort.SetAddStr(this_.maybeReferencesThis, name)
		}
	}
}

// Java class: com.google.javascript.jscomp.RescopeGlobalSymbols.RemoveGlobalVarCallback
type RescopeGlobalSymbols_RemoveGlobalVarCallback struct {
	NodeTraversal_AbstractShallowStatementCallbackBase
	compiler                     AbstractCompiler `gson:"compiler"`
	preserveLocalVarDeclarations bool             `gson:"preserveLocalVarDeclarations"`
}

func (this_ *RescopeGlobalSymbols_RemoveGlobalVarCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *RescopeGlobalSymbols_RemoveGlobalVarCallback) IsJavaNil_() bool { return this_ == nil }
func (*RescopeGlobalSymbols_RemoveGlobalVarCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RescopeGlobalSymbols$RemoveGlobalVarCallback", "RemoveGlobalVarCallback"
}
func (this_ *RescopeGlobalSymbols_RemoveGlobalVarCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractShallowStatementCallbackBase.IdentityHash_()
}
func (*RescopeGlobalSymbols_RemoveGlobalVarCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RescopeGlobalSymbols.RemoveGlobalVarCallback", "com.google.javascript.jscomp.RescopeGlobalSymbols$RemoveGlobalVarCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*RescopeGlobalSymbols_RemoveGlobalVarCallback](v_)
	})
}
func (*RescopeGlobalSymbols_RemoveGlobalVarCallback) Implements_jscomp_NodeTraversal_AbstractShallowStatementCallback() {
}
func (*RescopeGlobalSymbols_RemoveGlobalVarCallback) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *RescopeGlobalSymbols_RemoveGlobalVarCallback) initFields_RescopeGlobalSymbols_RemoveGlobalVarCallback() {
	_ = this_
}

// Java constructor: RemoveGlobalVarCallback(com.google.javascript.jscomp.AbstractCompiler compiler, boolean preserveLocalVarDeclarations)
func NewRescopeGlobalSymbols_RemoveGlobalVarCallback(compiler AbstractCompiler, preserveLocalVarDeclarations bool) *RescopeGlobalSymbols_RemoveGlobalVarCallback {
	this_ := &RescopeGlobalSymbols_RemoveGlobalVarCallback{}
	InitRescopeGlobalSymbols_RemoveGlobalVarCallback(this_, compiler, preserveLocalVarDeclarations)
	return this_
}

func InitRescopeGlobalSymbols_RemoveGlobalVarCallback(this_ *RescopeGlobalSymbols_RemoveGlobalVarCallback, compiler AbstractCompiler, preserveLocalVarDeclarations bool) {
	_ = compiler
	_ = preserveLocalVarDeclarations
	InitNodeTraversal_AbstractShallowStatementCallback(&this_.NodeTraversal_AbstractShallowStatementCallbackBase)
	this_.initFields_RescopeGlobalSymbols_RemoveGlobalVarCallback()
	this_.compiler = compiler
	this_.preserveLocalVarDeclarations = preserveLocalVarDeclarations
}

// Java method: RemoveGlobalVarCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RescopeGlobalSymbols_RemoveGlobalVarCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(NodeUtil_IsNameDeclaration(n)) {
		return
	}
	if RescopeGlobalSymbols_RemoveGlobalVarCallback_HasOnlyNameOrDestructuringChildren(n) {
		return
	}
	if NodeUtil_IsAnyFor(parent) && (parent.GetFirstChild() == n) {
		var commas []pkg_rhino.Node = []pkg_rhino.Node{}
		_ = commas
		{
			var c pkg_rhino.Node = n.GetFirstChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				commas = append(commas, c.CloneTree_0())
			}
		}
		if len(commas) != 0 {
			var comma pkg_rhino.Node = this_.priv_RemoveGlobalVarCallback_joinOnComma_4923(&commas, n)
			_ = comma
			comma.InsertBefore(n)
		}
	} else {
		var currentVar pkg_rhino.Node = pkg_rhino.Node(0)
		_ = currentVar
		{
			var c_2 pkg_rhino.Node = n.GetFirstChild()
			_ = c_2
			for ; c_2 != 0; c_2 = c_2.GetNext() {
				if !(!(c_2.IsDestructuringLhs())) {
					j2gort.CheckStateFmt(false, "Unexpected DESTRUCTURING_LHS child: %s", j2gort.BoxNode(c_2))
				}
				if this_.preserveLocalVarDeclarations && c_2.IsName() {
					if currentVar == 0 {
						currentVar = pkg_rhino.IR_Var_1(c_2.CloneTree_0()).Srcref(n)
						currentVar.SetJSDocInfo(n.GetJSDocInfo())
						currentVar.InsertBefore(n)
					} else {
						currentVar.AddChildToBack(c_2.CloneTree_0())
					}
				} else if c_2.IsAssign() {
					if currentVar != 0 {
						NodeUtil_MarkNewScopesChanged_2_Node_AbstractCompiler(currentVar, this_.compiler)
						currentVar = pkg_rhino.Node(0)
					}
					var expr pkg_rhino.Node = pkg_rhino.IR_ExprResult(c_2.CloneTree_0()).Srcref(c_2)
					_ = expr
					NodeUtil_MarkNewScopesChanged_2_Node_AbstractCompiler(expr, this_.compiler)
					expr.InsertBefore(n)
				}
			}
		}
		if currentVar != 0 {
			NodeUtil_MarkNewScopesChanged_2_Node_AbstractCompiler(currentVar, this_.compiler)
		}
	}
	n.Detach()
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(n, this_.compiler)
	this_.compiler.ReportChangeToEnclosingScope(parent)
}

// Java static method: RemoveGlobalVarCallback.hasOnlyNameOrDestructuringChildren(com.google.javascript.rhino.Node n)
func RescopeGlobalSymbols_RemoveGlobalVarCallback_HasOnlyNameOrDestructuringChildren(n pkg_rhino.Node) bool {
	_ = n
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if !(c.IsName()) && !(c.IsDestructuringLhs()) {
				return false
			}
		}
	}
	return n.HasChildren()
}

// Java method: RemoveGlobalVarCallback.joinOnComma(java.util.List<com.google.javascript.rhino.Node> commas, com.google.javascript.rhino.Node source)
func (this_ *RescopeGlobalSymbols_RemoveGlobalVarCallback) priv_RemoveGlobalVarCallback_joinOnComma_4923(commas *[]pkg_rhino.Node, source pkg_rhino.Node) pkg_rhino.Node {
	_ = commas
	_ = source
	var comma pkg_rhino.Node = (*commas)[int32(0)]
	_ = comma
	{
		var i int32 = int32(1)
		_ = i
		for ; i < int32(len(*commas)); i = i + 1 {
			var nextComma pkg_rhino.Node = pkg_rhino.IR_Comma(comma, (*commas)[i])
			_ = nextComma
			nextComma.SrcrefIfMissing(source)
			comma = nextComma
		}
	}
	return comma
}
func (*RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback) CanReleaseSyntacticScopes_() {
}
func (*RescopeGlobalSymbols_FindCrossChunkNamesCallback) AlwaysShouldTraverse_()      {}
func (*RescopeGlobalSymbols_FindCrossChunkNamesCallback) CanReleaseSyntacticScopes_() {}
func (*RescopeGlobalSymbols_FindNamesReferencingThis) AlwaysShouldTraverse_()         {}
func (*RescopeGlobalSymbols_FindNamesReferencingThis) CanReleaseSyntacticScopes_()    {}
func (*RescopeGlobalSymbols_RemoveGlobalVarCallback) CanReleaseSyntacticScopes_()     {}
