/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/FunctionToBlockMutator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.FunctionToBlockMutator
type FunctionToBlockMutator struct {
	compiler                 AbstractCompiler                       `gson:"compiler"`
	safeNameIdSupplier       j2gort.Ext_java_util_function_Supplier `gson:"safeNameIdSupplier"`
	functionArgumentInjector *FunctionArgumentInjector              `gson:"functionArgumentInjector"`
	idHash_                  atomic.Int32
}

func (this_ *FunctionToBlockMutator) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionToBlockMutator) IsJavaNil_() bool { return this_ == nil }
func (*FunctionToBlockMutator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionToBlockMutator", "FunctionToBlockMutator"
}
func (this_ *FunctionToBlockMutator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FunctionToBlockMutator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionToBlockMutator", "com.google.javascript.jscomp.FunctionToBlockMutator", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionToBlockMutator](v_)
	})
}

func (this_ *FunctionToBlockMutator) initFields_FunctionToBlockMutator() {
	_ = this_
}

// Java constructor: FunctionToBlockMutator(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.function.Supplier<java.lang.String> safeNameIdSupplier)
func NewFunctionToBlockMutator_2(compiler AbstractCompiler, safeNameIdSupplier j2gort.Ext_java_util_function_Supplier) *FunctionToBlockMutator {
	this_ := &FunctionToBlockMutator{}
	InitFunctionToBlockMutator_2(this_, compiler, safeNameIdSupplier)
	return this_
}

func InitFunctionToBlockMutator_2(this_ *FunctionToBlockMutator, compiler AbstractCompiler, safeNameIdSupplier j2gort.Ext_java_util_function_Supplier) {
	_ = compiler
	_ = safeNameIdSupplier
	InitFunctionToBlockMutator_3(this_, compiler, safeNameIdSupplier, NewFunctionArgumentInjector(compiler.GetAstAnalyzer()))
}

// Java constructor: FunctionToBlockMutator(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.function.Supplier<java.lang.String> safeNameIdSupplier, com.google.javascript.jscomp.FunctionArgumentInjector functionArgumentInjector)
func NewFunctionToBlockMutator_3(compiler AbstractCompiler, safeNameIdSupplier j2gort.Ext_java_util_function_Supplier, functionArgumentInjector *FunctionArgumentInjector) *FunctionToBlockMutator {
	this_ := &FunctionToBlockMutator{}
	InitFunctionToBlockMutator_3(this_, compiler, safeNameIdSupplier, functionArgumentInjector)
	return this_
}

func InitFunctionToBlockMutator_3(this_ *FunctionToBlockMutator, compiler AbstractCompiler, safeNameIdSupplier j2gort.Ext_java_util_function_Supplier, functionArgumentInjector *FunctionArgumentInjector) {
	_ = compiler
	_ = safeNameIdSupplier
	_ = functionArgumentInjector
	this_.initFields_FunctionToBlockMutator()
	this_.compiler = j2gort.CheckNotNullIface(compiler)
	this_.safeNameIdSupplier = safeNameIdSupplier
	this_.functionArgumentInjector = j2gort.CheckNotNullPtr(functionArgumentInjector)
}

// Java method: FunctionToBlockMutator.mutate(java.lang.String fnName, com.google.javascript.rhino.Node fnNode, com.google.javascript.rhino.Node callNode, java.lang.String resultName, boolean needsDefaultResult, boolean isCallInLoop)
func (this_ *FunctionToBlockMutator) Mutate(fnName string, fnNode pkg_rhino.Node, callNode pkg_rhino.Node, resultName string, needsDefaultResult bool, isCallInLoop bool) pkg_rhino.Node {
	_ = fnName
	_ = fnNode
	_ = callNode
	_ = resultName
	_ = needsDefaultResult
	_ = isCallInLoop
	return this_.priv_FunctionToBlockMutator_mutateInternal_2959(fnName, fnNode, callNode, resultName, needsDefaultResult, isCallInLoop, true)
}

// Java method: FunctionToBlockMutator.mutateWithoutRenaming(java.lang.String fnName, com.google.javascript.rhino.Node fnNode, com.google.javascript.rhino.Node callNode, java.lang.String resultName, boolean needsDefaultResult, boolean isCallInLoop)
func (this_ *FunctionToBlockMutator) MutateWithoutRenaming(fnName string, fnNode pkg_rhino.Node, callNode pkg_rhino.Node, resultName string, needsDefaultResult bool, isCallInLoop bool) pkg_rhino.Node {
	_ = fnName
	_ = fnNode
	_ = callNode
	_ = resultName
	_ = needsDefaultResult
	_ = isCallInLoop
	return this_.priv_FunctionToBlockMutator_mutateInternal_2959(fnName, fnNode, callNode, resultName, needsDefaultResult, isCallInLoop, false)
}

// Java method: FunctionToBlockMutator.mutateInternal(java.lang.String fnName, com.google.javascript.rhino.Node fnNode, com.google.javascript.rhino.Node callNode, java.lang.String resultName, boolean needsDefaultResult, boolean isCallInLoop, boolean renameLocals)
func (this_ *FunctionToBlockMutator) priv_FunctionToBlockMutator_mutateInternal_2959(fnName string, fnNode pkg_rhino.Node, callNode pkg_rhino.Node, resultName string, needsDefaultResult bool, isCallInLoop bool, renameLocals bool) pkg_rhino.Node {
	_ = fnName
	_ = fnNode
	_ = callNode
	_ = resultName
	_ = needsDefaultResult
	_ = isCallInLoop
	_ = renameLocals
	var newFnNode pkg_rhino.Node = fnNode.CloneTree_0()
	_ = newFnNode
	pkg_rhino.IR_Root([]pkg_rhino.Node{pkg_rhino.IR_Script_1_Node([]pkg_rhino.Node{newFnNode})})
	if renameLocals {
		this_.priv_FunctionToBlockMutator_makeLocalNamesUnique_2960(newFnNode, isCallInLoop)
		FunctionToBlockMutator_RewriteFunctionDeclarations(newFnNode.GetLastChild())
	}
	var modifiedParameters j2gort.Ext_com_google_common_collect_ImmutableSet = this_.functionArgumentInjector.FindModifiedParameters(newFnNode)
	_ = modifiedParameters
	var newBlock pkg_rhino.Node = pkg_rhino.Node(0)
	_ = newBlock
	var args *FunctionArgumentInjector_ParamArgList = this_.functionArgumentInjector.AcquireFunctionCallParameterMap(newFnNode, callNode, this_.safeNameIdSupplier)
	_ = args
	func() {
		defer func() {
			this_.functionArgumentInjector.ReleaseParamArgList(args)
		}()
		var hasArgs bool = !(args.IsEmpty())
		_ = hasArgs
		var allNamesToAlias j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
		_ = allNamesToAlias
		if hasArgs {
			var temps j2gort.Ext_com_google_common_collect_ImmutableSet = this_.functionArgumentInjector.GatherCallArgumentsNeedingTemps(this_.compiler, newFnNode, args, modifiedParameters, this_.compiler.GetCodingConvention())
			_ = temps
			if modifiedParameters.IsEmpty() {
				allNamesToAlias = temps
			} else if temps.IsEmpty() {
				allNamesToAlias = modifiedParameters
			} else {
				allNamesToAlias = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder().AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](modifiedParameters))).AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](temps))).Build())
			}
		}
		newBlock = NodeUtil_GetFunctionBody(newFnNode)
		newBlock.Detach()
		if hasArgs {
			var inlineResult pkg_rhino.Node = this_.priv_FunctionToBlockMutator_aliasAndInlineArguments_2963(newBlock, args, j2gort.AsIface[j2gort.Ext_java_util_Set](allNamesToAlias))
			_ = inlineResult
			if !(newBlock == inlineResult) {
				j2gort.CheckState(false)
			}
		}
	}()
	if isCallInLoop {
		FunctionToBlockMutator_FixUninitializedVarDeclarations_2(newBlock, newBlock)
	}
	var labelName string = this_.priv_FunctionToBlockMutator_getLabelNameForFunction_2961(fnName)
	_ = labelName
	var injectableBlock pkg_rhino.Node = FunctionToBlockMutator_ReplaceReturns(newBlock, resultName, labelName, needsDefaultResult)
	_ = injectableBlock
	if !(injectableBlock != 0) {
		j2gort.CheckState(false)
	}
	return injectableBlock
}

// Java static method: FunctionToBlockMutator.rewriteFunctionDeclarations(com.google.javascript.rhino.Node n)
func FunctionToBlockMutator_RewriteFunctionDeclarations(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	if n.IsFunction() {
		if NodeUtil_IsFunctionDeclaration(n) {
			var fnNameNode pkg_rhino.Node = n.GetFirstChild()
			_ = fnNameNode
			var name pkg_rhino.Node = pkg_rhino.IR_Name(fnNameNode.GetString()).Srcref(fnNameNode)
			_ = name
			var var_ pkg_rhino.Node = pkg_rhino.IR_Var_1(name).Srcref(n)
			_ = var_
			fnNameNode.SetString("")
			n.ReplaceWith(var_)
			name.AddChildToFront(n)
			return var_
		}
		return pkg_rhino.Node(0)
	}
	var functionsToHoist []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = functionsToHoist
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		var next pkg_rhino.Node = pkg_rhino.Node(0)
		_ = next
		for ; c != 0; c = next {
			next = c.GetNext()
			var fnVar pkg_rhino.Node = FunctionToBlockMutator_RewriteFunctionDeclarations(c)
			_ = fnVar
			if fnVar != 0 {
				functionsToHoist = j2gort.SliceInsertAt(functionsToHoist, int32(0), fnVar)
			}
		}
	}
	for _, fnVar_2 := range functionsToHoist {
		_ = fnVar_2
		if n.GetFirstChild() != fnVar_2 {
			n.AddChildToFront(fnVar_2.Detach())
		}
	}

	return pkg_rhino.Node(0)
}

// Java static method: FunctionToBlockMutator.fixUninitializedVarDeclarations(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node containingBlock)
func FunctionToBlockMutator_FixUninitializedVarDeclarations_2(n pkg_rhino.Node, containingBlock pkg_rhino.Node) {
	_ = n
	_ = containingBlock
	FunctionToBlockMutator_FixUninitializedVarDeclarations_3(n, containingBlock, false)
}

// Java static method: FunctionToBlockMutator.fixUninitializedVarDeclarations(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node containingBlock, boolean inLoop)
func FunctionToBlockMutator_FixUninitializedVarDeclarations_3(n pkg_rhino.Node, containingBlock pkg_rhino.Node, inLoop bool) {
	_ = n
	_ = containingBlock
	_ = inLoop
	if n.IsFunction() {
		return
	}
	if NodeUtil_IsLoopStructure(n) {
		var body pkg_rhino.Node = NodeUtil_GetLoopCodeBlock(n)
		_ = body
		if body != 0 {
			FunctionToBlockMutator_FixUninitializedVarDeclarations_3(body, containingBlock, true)
		}
		if n.IsVanillaFor() {
			var init pkg_rhino.Node = n.GetFirstChild()
			_ = init
			if (init != 0) && !(init.IsEmpty()) {
				FunctionToBlockMutator_FixUninitializedVarDeclarations_3(init, containingBlock, true)
			}
		}
		return
	}
	if (inLoop && n.IsVar()) && n.HasOneChild() {
		var name pkg_rhino.Node = n.GetFirstChild()
		_ = name
		var srcLocation pkg_rhino.Node = name
		_ = srcLocation
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		var hoistedVar pkg_rhino.Node = pkg_rhino.IR_Var_1(name.CloneNode_0().SetJSDocInfo(nil)).Srcref(n)
		_ = hoistedVar
		hoistedVar.GetFirstChild().AddChildToBack(NodeUtil_NewUndefinedNode(srcLocation))
		containingBlock.AddChildToFront(hoistedVar)
		if name.HasChildren() {
			var value pkg_rhino.Node = name.RemoveFirstChild()
			_ = value
			var assign pkg_rhino.Node = pkg_rhino.IR_Assign(name.Detach(), value).Srcref(n)
			_ = assign
			if ((parent != 0) && parent.IsVanillaFor()) && (parent.GetFirstChild() == n) {
				n.ReplaceWith(assign)
			} else {
				n.ReplaceWith(NodeUtil_NewExpr(assign))
			}
		} else {
			if ((parent != 0) && parent.IsVanillaFor()) && (parent.GetFirstChild() == n) {
				n.ReplaceWith(pkg_rhino.IR_Empty())
			} else {
				n.Detach()
			}
		}
		return
	}
	if (!(inLoop) && (n.IsVar() || n.IsLet())) && n.HasOneChild() {
		var name_2 pkg_rhino.Node = n.GetFirstChild()
		_ = name_2
		if !(name_2.HasChildren()) {
			var srcLocation_2 pkg_rhino.Node = name_2
			_ = srcLocation_2
			var undefinedNode pkg_rhino.Node = NodeUtil_NewUndefinedNode(srcLocation_2)
			_ = undefinedNode
			name_2.AddChildToBack(undefinedNode)
			var parent_2 pkg_rhino.Node = n.GetParent()
			_ = parent_2
			if ((parent_2 != 0) && parent_2.IsVanillaFor()) && (parent_2.GetFirstChild() == n) {
				n.ReplaceWith(pkg_rhino.IR_Empty())
			} else {
				n.Detach()
			}
			containingBlock.AddChildToFront(n)
		}
		return
	}
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			FunctionToBlockMutator_FixUninitializedVarDeclarations_3(c, containingBlock, inLoop)
		}
	}
}

// Java method: FunctionToBlockMutator.makeLocalNamesUnique(com.google.javascript.rhino.Node fnNode, boolean isCallInLoop)
func (this_ *FunctionToBlockMutator) priv_FunctionToBlockMutator_makeLocalNamesUnique_2960(fnNode pkg_rhino.Node, isCallInLoop bool) {
	_ = fnNode
	_ = isCallInLoop
	var idSupplier j2gort.Ext_java_util_function_Supplier = j2gort.AsIface[j2gort.Ext_java_util_function_Supplier](this_.compiler.GetUniqueNameIdSupplier())
	_ = idSupplier
	NodeTraversal_TraverseScopeRoots(this_.compiler, j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(j2gort.BoxNode(fnNode))), ToIface_NodeTraversal_Callback(MakeDeclaredNamesUnique_S_Builder().WithRenamer(NewMakeDeclaredNamesUnique_InlineRenamer(this_.compiler.GetCodingConvention(), idSupplier, "inline_", isCallInLoop, true, nil)).WithMarkChanges(false).Build()), true)
	{
		var val_2 RenameLabels
		InitRenameLabels_4(&val_2, this_.compiler, NewFunctionToBlockMutator_LabelNameSupplier(idSupplier), false, false)
		(&val_2).Process_2(pkg_rhino.Node(0), fnNode)
	}
}

// Java method: FunctionToBlockMutator.getLabelNameForFunction(java.lang.String fnName)
func (this_ *FunctionToBlockMutator) priv_FunctionToBlockMutator_getLabelNameForFunction_2961(fnName string) string {
	_ = fnName
	var name string
	_ = name
	if j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(fnName) {
		name = "anon"
	} else {
		name = fnName
	}
	return ("JSCompiler_inline_label_" + j2gort.StrValString(name) + "_" + j2gort.StrValString(j2gort.CastString(this_.safeNameIdSupplier.Get())))
}

// Java method: FunctionToBlockMutator.getUniqueThisName()
func (this_ *FunctionToBlockMutator) priv_FunctionToBlockMutator_getUniqueThisName_2962() string {
	return ("JSCompiler_inline_this_" + j2gort.StrValString(j2gort.CastString(this_.safeNameIdSupplier.Get())))
}

// Java method: FunctionToBlockMutator.aliasAndInlineArguments(com.google.javascript.rhino.Node fnTemplateRoot, com.google.javascript.jscomp.FunctionArgumentInjector.ParamArgList paramToArgMap, java.util.Set<java.lang.String> namesToAlias)
func (this_ *FunctionToBlockMutator) priv_FunctionToBlockMutator_aliasAndInlineArguments_2963(fnTemplateRoot pkg_rhino.Node, paramToArgMap *FunctionArgumentInjector_ParamArgList, namesToAlias j2gort.Ext_java_util_Set) pkg_rhino.Node {
	_ = fnTemplateRoot
	_ = paramToArgMap
	_ = namesToAlias
	if (namesToAlias == nil) || namesToAlias.IsEmpty() {
		var result pkg_rhino.Node = this_.functionArgumentInjector.Inject_4_AbstractCompiler_Node_Node_FunctionArgumentInjector_ParamArgList(this_.compiler, fnTemplateRoot, pkg_rhino.Node(0), paramToArgMap)
		_ = result
		if !(result == fnTemplateRoot) {
			j2gort.CheckState(false)
		}
		return result
	} else {
		var entryCount int32 = paramToArgMap.size
		_ = entryCount
		var newAliases []pkg_rhino.Node = make([]pkg_rhino.Node, entryCount)
		_ = newAliases
		var aliasCount int32 = int32(0)
		_ = aliasCount
		{
			var i int32 = int32(0)
			_ = i
			for ; i < entryCount; i = i + 1 {
				var name string = paramToArgMap.GetNameAt(i)
				_ = name
				if !(j2gort.SetContainsStr(namesToAlias, name)) {
					continue
				}
				var originalValue pkg_rhino.Node = paramToArgMap.GetArgAt(i)
				_ = originalValue
				{
					var iifeCond_1 bool
					{

						sr_3 := name
						se_4 := "this"
						iifeCond_1 = j2gort.CheckNotNullString(sr_3) == se_4
					}
					if iifeCond_1 {
						var referencesThis bool = NodeUtil_ReferencesEnclosingReceiver(fnTemplateRoot)
						_ = referencesThis
						var replacement pkg_rhino.Node = originalValue
						_ = replacement
						if !(originalValue.IsThis()) && (referencesThis || this_.compiler.GetAstAnalyzer().MayHaveSideEffects(originalValue)) {
							var newName string = this_.priv_FunctionToBlockMutator_getUniqueThisName_2962()
							_ = newName
							var newValue pkg_rhino.Node = originalValue.CloneTree_0()
							_ = newValue
							var newNode pkg_rhino.Node = NodeUtil_NewVarNode_2_string_Node(newName, newValue).SrcrefTreeIfMissing(newValue)
							_ = newNode
							{
								a_5 := newAliases
								var i_6 int32
								{

									old_7 := aliasCount
									nw_8 := old_7 + 1
									aliasCount = nw_8
									i_6 = old_7
								}
								a_5[i_6] = newNode
							}
							replacement = pkg_rhino.IR_Name(newName).SrcrefTree(newValue)
						}
						paramToArgMap.SetArgAt(i, replacement)
					} else {
						var newName_2 pkg_rhino.Node = paramToArgMap.GetParamNodeAt(i).CloneNode_0()
						_ = newName_2
						var newValue_2 pkg_rhino.Node = originalValue.CloneTree_0()
						_ = newValue_2
						newName_2.Srcref(newValue_2)
						var newNode_2 pkg_rhino.Node = pkg_rhino.IR_Var_2(newName_2, newValue_2).SrcrefTreeIfMissing(newValue_2)
						_ = newNode_2
						{
							a_9 := newAliases
							var i_10 int32
							{

								old_11 := aliasCount
								nw_12 := old_11 + 1
								aliasCount = nw_12
								i_10 = old_11
							}
							a_9[i_10] = newNode_2
						}
						paramToArgMap.SetArgAt(i, pkg_rhino.Node(0))
					}
				}
			}
		}
		var result_2 pkg_rhino.Node = this_.functionArgumentInjector.Inject_4_AbstractCompiler_Node_Node_FunctionArgumentInjector_ParamArgList(this_.compiler, fnTemplateRoot, pkg_rhino.Node(0), paramToArgMap)
		_ = result_2
		if !(result_2 == fnTemplateRoot) {
			j2gort.CheckState(false)
		}
		{
			var j int32 = (aliasCount - int32(1))
			_ = j
			for ; j >= int32(0); j = j - 1 {
				fnTemplateRoot.AddChildToFront(newAliases[j])
			}
		}
		return result_2
	}
}

// Java static method: FunctionToBlockMutator.replaceReturns(com.google.javascript.rhino.Node block, java.lang.String resultName, java.lang.String labelName, boolean resultMustBeSet)
func FunctionToBlockMutator_ReplaceReturns(block pkg_rhino.Node, resultName string, labelName string, resultMustBeSet bool) pkg_rhino.Node {
	_ = block
	_ = resultName
	_ = labelName
	_ = resultMustBeSet
	j2gort.CheckNotNullNode(block)
	j2gort.CheckNotNullString(labelName)
	var root pkg_rhino.Node = block
	_ = root
	var hasReturnAtExit bool = false
	_ = hasReturnAtExit
	var returnCount int32 = NodeUtil_GetNodeTypeReferenceCount(block, pkg_rhino.StaticGet_Token_RETURN(), NewNodeUtil_MatchShallowStatement())
	_ = returnCount
	if returnCount > int32(0) {
		hasReturnAtExit = FunctionToBlockMutator_HasReturnAtExit(block)
		if hasReturnAtExit {
			FunctionToBlockMutator_ConvertLastReturnToStatement(block, resultName)
			returnCount = returnCount - 1
		}
		if returnCount > int32(0) {
			FunctionToBlockMutator_ReplaceReturnWithBreak(block, pkg_rhino.Node(0), resultName, labelName)
			var name pkg_rhino.Node = pkg_rhino.IR_LabelName(labelName).Srcref(block)
			_ = name
			var label pkg_rhino.Node = pkg_rhino.IR_Label(name, block).Srcref(block)
			_ = label
			var newRoot pkg_rhino.Node = pkg_rhino.IR_Block_0().Srcref(block)
			_ = newRoot
			newRoot.AddChildToBack(label)
			root = newRoot
		}
	}
	if (resultMustBeSet && !(hasReturnAtExit)) && (resultName != j2gort.NullString) {
		FunctionToBlockMutator_AddDummyAssignment(block, resultName)
	}
	return root
}

// Java static method: FunctionToBlockMutator.addDummyAssignment(com.google.javascript.rhino.Node node, java.lang.String resultName)
func FunctionToBlockMutator_AddDummyAssignment(node pkg_rhino.Node, resultName string) {
	_ = node
	_ = resultName
	if !(node.IsBlock()) {
		j2gort.CheckArgument(false)
	}
	var srcLocation pkg_rhino.Node = node
	_ = srcLocation
	var retVal pkg_rhino.Node = NodeUtil_NewUndefinedNode(srcLocation)
	_ = retVal
	var resultNode pkg_rhino.Node = FunctionToBlockMutator_CreateAssignStatementNode(resultName, retVal)
	_ = resultNode
	resultNode.SrcrefTreeIfMissing(node)
	node.AddChildToBack(resultNode)
}

// Java static method: FunctionToBlockMutator.convertLastReturnToStatement(com.google.javascript.rhino.Node block, java.lang.String resultName)
func FunctionToBlockMutator_ConvertLastReturnToStatement(block pkg_rhino.Node, resultName string) {
	_ = block
	_ = resultName
	var ret pkg_rhino.Node = block.GetLastChild()
	_ = ret
	if !(ret.IsReturn()) {
		j2gort.CheckArgument(false)
	}
	var resultNode pkg_rhino.Node = FunctionToBlockMutator_GetReplacementReturnStatement(ret, resultName)
	_ = resultNode
	if resultNode == 0 {
		ret.Detach()
	} else {
		resultNode.SrcrefTreeIfMissing(ret)
		ret.ReplaceWith(resultNode)
	}
}

// Java static method: FunctionToBlockMutator.createAssignStatementNode(java.lang.String name, com.google.javascript.rhino.Node expression)
func FunctionToBlockMutator_CreateAssignStatementNode(name string, expression pkg_rhino.Node) pkg_rhino.Node {
	_ = name
	_ = expression
	var nameNode pkg_rhino.Node = pkg_rhino.IR_Name(name)
	_ = nameNode
	var assign pkg_rhino.Node = pkg_rhino.IR_Assign(nameNode, expression)
	_ = assign
	return NodeUtil_NewExpr(assign)
}

// Java static method: FunctionToBlockMutator.getReplacementReturnStatement(com.google.javascript.rhino.Node node, java.lang.String resultName)
func FunctionToBlockMutator_GetReplacementReturnStatement(node pkg_rhino.Node, resultName string) pkg_rhino.Node {
	_ = node
	_ = resultName
	var resultNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = resultNode
	var retVal pkg_rhino.Node = pkg_rhino.Node(0)
	_ = retVal
	if node.HasChildren() {
		retVal = node.GetFirstChild().CloneTree_0()
	}
	if resultName == j2gort.NullString {
		if retVal != 0 {
			resultNode = NodeUtil_NewExpr(retVal)
		}
	} else {
		if retVal == 0 {
			var srcLocation pkg_rhino.Node = node
			_ = srcLocation
			retVal = NodeUtil_NewUndefinedNode(srcLocation)
		}
		resultNode = FunctionToBlockMutator_CreateAssignStatementNode(resultName, retVal)
	}
	return resultNode
}

// Java static method: FunctionToBlockMutator.hasReturnAtExit(com.google.javascript.rhino.Node block)
func FunctionToBlockMutator_HasReturnAtExit(block pkg_rhino.Node) bool {
	_ = block
	return (block.GetLastChild().IsReturn())
}

// Java static method: FunctionToBlockMutator.replaceReturnWithBreak(com.google.javascript.rhino.Node current, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node parent, java.lang.String resultName, java.lang.String labelName)
func FunctionToBlockMutator_ReplaceReturnWithBreak(current pkg_rhino.Node, parent pkg_rhino.Node, resultName string, labelName string) pkg_rhino.Node {
	_ = current
	_ = parent
	_ = resultName
	_ = labelName
	if current.IsFunction() || current.IsExprResult() {
		return current
	}
	if current.IsReturn() {
		if !(NodeUtil_IsStatementBlock(parent)) {
			j2gort.CheckState(false)
		}
		var resultNode pkg_rhino.Node = FunctionToBlockMutator_GetReplacementReturnStatement(current, resultName)
		_ = resultNode
		var breakNode pkg_rhino.Node = pkg_rhino.IR_BreakNode_1(pkg_rhino.IR_LabelName(labelName))
		_ = breakNode
		breakNode.SrcrefTreeIfMissing(current)
		current.ReplaceWith(breakNode)
		if resultNode != 0 {
			resultNode.SrcrefTreeIfMissing(current)
			resultNode.InsertBefore(breakNode)
		}
		current = breakNode
	} else {
		{
			var c pkg_rhino.Node = current.GetFirstChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				c = FunctionToBlockMutator_ReplaceReturnWithBreak(c, current, resultName, labelName)
			}
		}
	}
	return current
}

// Java class: com.google.javascript.jscomp.FunctionToBlockMutator.LabelNameSupplier
type FunctionToBlockMutator_LabelNameSupplier struct {
	idSupplier j2gort.Ext_java_util_function_Supplier `gson:"idSupplier"`
	idHash_    atomic.Int32
}

func (this_ *FunctionToBlockMutator_LabelNameSupplier) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionToBlockMutator_LabelNameSupplier) IsJavaNil_() bool { return this_ == nil }
func (*FunctionToBlockMutator_LabelNameSupplier) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionToBlockMutator$LabelNameSupplier", "LabelNameSupplier"
}
func (this_ *FunctionToBlockMutator_LabelNameSupplier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FunctionToBlockMutator_LabelNameSupplier) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionToBlockMutator.LabelNameSupplier", "com.google.javascript.jscomp.FunctionToBlockMutator$LabelNameSupplier", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionToBlockMutator_LabelNameSupplier](v_)
	})
}
func (*FunctionToBlockMutator_LabelNameSupplier) Implements_j2gort_Ext_java_util_function_Supplier() {
}

func (this_ *FunctionToBlockMutator_LabelNameSupplier) initFields_FunctionToBlockMutator_LabelNameSupplier() {
	_ = this_
}

// Java constructor: LabelNameSupplier(java.util.function.Supplier<java.lang.String> idSupplier)
func NewFunctionToBlockMutator_LabelNameSupplier(idSupplier j2gort.Ext_java_util_function_Supplier) *FunctionToBlockMutator_LabelNameSupplier {
	this_ := &FunctionToBlockMutator_LabelNameSupplier{}
	InitFunctionToBlockMutator_LabelNameSupplier(this_, idSupplier)
	return this_
}

func InitFunctionToBlockMutator_LabelNameSupplier(this_ *FunctionToBlockMutator_LabelNameSupplier, idSupplier j2gort.Ext_java_util_function_Supplier) {
	_ = idSupplier
	this_.initFields_FunctionToBlockMutator_LabelNameSupplier()
	this_.idSupplier = idSupplier
}

// Java method: LabelNameSupplier.get()
func (this_ *FunctionToBlockMutator_LabelNameSupplier) Get() any {
	return "JSCompiler_inline_label_" + j2gort.StrValString(j2gort.CastString(this_.idSupplier.Get()))
}
