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
// Port of com/google/javascript/jscomp/FunctionArgumentInjector.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.FunctionArgumentInjector
type FunctionArgumentInjector struct {
	astAnalyzer          *AstAnalyzer                                        `gson:"astAnalyzer"`
	reusableCollector    *FunctionArgumentInjector_ReferencedAfterSideEffect `gson:"reusableCollector"`
	reusableParamArgList *FunctionArgumentInjector_ParamArgList              `gson:"reusableParamArgList"`
	idHash_              atomic.Int32
	collectorInUse       bool `gson:"collectorInUse"`
	paramArgListInUse    bool `gson:"paramArgListInUse"`
}

func (this_ *FunctionArgumentInjector) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionArgumentInjector) IsJavaNil_() bool { return this_ == nil }
func (*FunctionArgumentInjector) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionArgumentInjector", "FunctionArgumentInjector"
}
func (this_ *FunctionArgumentInjector) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FunctionArgumentInjector) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionArgumentInjector", "com.google.javascript.jscomp.FunctionArgumentInjector", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionArgumentInjector](v_)
	})
}

func StaticGet_FunctionArgumentInjector_THIS_MARKER() string { return "this" }

var static_FunctionArgumentInjector_THIS_MARKER_NODE pkg_rhino.Node

func StaticGet_FunctionArgumentInjector_THIS_MARKER_NODE() pkg_rhino.Node {
	clinit_FunctionArgumentInjector()
	return static_FunctionArgumentInjector_THIS_MARKER_NODE
}
func StaticSet_FunctionArgumentInjector_THIS_MARKER_NODE(v pkg_rhino.Node) {
	clinit_FunctionArgumentInjector()
	static_FunctionArgumentInjector_THIS_MARKER_NODE = v
}
func StaticGet_FunctionArgumentInjector_REST_MARKER() string { return "rest param" }

var clinitDone_FunctionArgumentInjector uint32
var clinitState_FunctionArgumentInjector j2gort.ClinitState

func clinit_FunctionArgumentInjector() {
	if atomic.LoadUint32(&clinitDone_FunctionArgumentInjector) != 0 {
		return
	}
	clinitSlow_FunctionArgumentInjector()
}

func clinitSlow_FunctionArgumentInjector() {
	if !j2gort.EnterClinit(&clinitState_FunctionArgumentInjector, &clinitDone_FunctionArgumentInjector) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_FunctionArgumentInjector, &clinitDone_FunctionArgumentInjector)
	static_FunctionArgumentInjector_THIS_MARKER_NODE = pkg_rhino.IR_Name("this")
}

func (this_ *FunctionArgumentInjector) initFields_FunctionArgumentInjector() {
	_ = this_
}

// Java constructor: FunctionArgumentInjector(com.google.javascript.jscomp.AstAnalyzer astAnalyzer)
func NewFunctionArgumentInjector(astAnalyzer *AstAnalyzer) *FunctionArgumentInjector {
	clinit_FunctionArgumentInjector()
	this_ := &FunctionArgumentInjector{}
	InitFunctionArgumentInjector(this_, astAnalyzer)
	return this_
}

func InitFunctionArgumentInjector(this_ *FunctionArgumentInjector, astAnalyzer *AstAnalyzer) {
	_ = astAnalyzer
	this_.initFields_FunctionArgumentInjector()
	this_.astAnalyzer = astAnalyzer
}

// Java method: FunctionArgumentInjector.inject(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent, com.google.javascript.jscomp.FunctionArgumentInjector.ParamArgList replacements)
func (this_ *FunctionArgumentInjector) Inject_4_AbstractCompiler_Node_Node_FunctionArgumentInjector_ParamArgList(compiler AbstractCompiler, node pkg_rhino.Node, parent pkg_rhino.Node, replacements *FunctionArgumentInjector_ParamArgList) pkg_rhino.Node {
	_ = compiler
	_ = node
	_ = parent
	_ = replacements
	return this_.priv_FunctionArgumentInjector_inject_2898(compiler, node, parent, replacements, true)
}

// Java method: FunctionArgumentInjector.inject(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent, com.google.javascript.jscomp.FunctionArgumentInjector.ParamArgList replacements, boolean replaceThis)
func (this_ *FunctionArgumentInjector) priv_FunctionArgumentInjector_inject_2898(compiler AbstractCompiler, node pkg_rhino.Node, parent pkg_rhino.Node, replacements *FunctionArgumentInjector_ParamArgList, replaceThis bool) pkg_rhino.Node {
	_ = compiler
	_ = node
	_ = parent
	_ = replacements
	_ = replaceThis
	if node.IsName() {
		var replacementTemplate pkg_rhino.Node = replacements.GetArg(node.GetString())
		_ = replacementTemplate
		if replacementTemplate != 0 {
			if !(!((parent.IsFunction() || parent.IsVar()) || parent.IsCatch())) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
			}
			var replacement pkg_rhino.Node = replacementTemplate.CloneTree_0()
			_ = replacement
			node.ReplaceWith(replacement)
			return replacement
		}
	} else if replaceThis && node.IsThis() {
		var replacementTemplate_2 pkg_rhino.Node = replacements.GetThisArg()
		_ = replacementTemplate_2
		j2gort.CheckNotNullNode(replacementTemplate_2)
		if !(replacementTemplate_2.IsThis()) {
			var replacement_2 pkg_rhino.Node = replacementTemplate_2.CloneTree_0()
			_ = replacement_2
			node.ReplaceWith(replacement_2)
			if compiler.GetAstAnalyzer().MayHaveSideEffects(replacementTemplate_2) {
				replacements.ClearThisArg()
			}
			return replacement_2
		}
	} else if ((node.IsFunction() && !(node.IsArrowFunction())) || node.IsMemberFieldDef()) || NodeUtil_IsClassStaticBlock(node) {
		replaceThis = false
	}
	if node.IsComputedFieldDef() {
		var key pkg_rhino.Node = node.GetFirstChild()
		_ = key
		if key != 0 {
			key = this_.priv_FunctionArgumentInjector_inject_2898(compiler, key, node, replacements, replaceThis)
			var value pkg_rhino.Node = key.GetNext()
			_ = value
			if value != 0 {
				value = this_.priv_FunctionArgumentInjector_inject_2898(compiler, value, node, replacements, false)
			}
		}
		return node
	}
	{
		var c pkg_rhino.Node = node.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			c = this_.priv_FunctionArgumentInjector_inject_2898(compiler, c, node, replacements, replaceThis)
		}
	}
	return node
}

// Java method: FunctionArgumentInjector.inject(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent, java.util.Map<java.lang.String,com.google.javascript.rhino.Node> replacements)
func (this_ *FunctionArgumentInjector) Inject_4_AbstractCompiler_Node_Node_Ext_java_util_Map(compiler AbstractCompiler, node pkg_rhino.Node, parent pkg_rhino.Node, replacements j2gort.Ext_java_util_Map) pkg_rhino.Node {
	_ = compiler
	_ = node
	_ = parent
	_ = replacements
	return this_.priv_FunctionArgumentInjector_inject_2900(compiler, node, parent, replacements, true)
}

// Java method: FunctionArgumentInjector.inject(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent, java.util.Map<java.lang.String,com.google.javascript.rhino.Node> replacements, boolean replaceThis)
func (this_ *FunctionArgumentInjector) priv_FunctionArgumentInjector_inject_2900(compiler AbstractCompiler, node pkg_rhino.Node, parent pkg_rhino.Node, replacements j2gort.Ext_java_util_Map, replaceThis bool) pkg_rhino.Node {
	_ = compiler
	_ = node
	_ = parent
	_ = replacements
	_ = replaceThis
	if node.IsName() {
		var replacementTemplate pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](j2gort.MapGetStr(replacements, node.GetString()))
		_ = replacementTemplate
		if replacementTemplate != 0 {
			if !(!((parent.IsFunction() || parent.IsVar()) || parent.IsCatch())) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
			}
			var replacement pkg_rhino.Node = replacementTemplate.CloneTree_0()
			_ = replacement
			node.ReplaceWith(replacement)
			return replacement
		}
	} else if replaceThis && node.IsThis() {
		var replacementTemplate_2 pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](j2gort.MapGetStr(replacements, "this"))
		_ = replacementTemplate_2
		j2gort.CheckNotNullNode(replacementTemplate_2)
		if !(replacementTemplate_2.IsThis()) {
			var replacement_2 pkg_rhino.Node = replacementTemplate_2.CloneTree_0()
			_ = replacement_2
			node.ReplaceWith(replacement_2)
			if compiler.GetAstAnalyzer().MayHaveSideEffects(replacementTemplate_2) {
				j2gort.MapRemoveStr(replacements, "this")
			}
			return replacement_2
		}
	} else if ((node.IsFunction() && !(node.IsArrowFunction())) || node.IsMemberFieldDef()) || NodeUtil_IsClassStaticBlock(node) {
		replaceThis = false
	}
	if node.IsComputedFieldDef() {
		var key pkg_rhino.Node = node.GetFirstChild()
		_ = key
		if key != 0 {
			key = this_.priv_FunctionArgumentInjector_inject_2900(compiler, key, node, replacements, replaceThis)
			var value pkg_rhino.Node = key.GetNext()
			_ = value
			if value != 0 {
				value = this_.priv_FunctionArgumentInjector_inject_2900(compiler, value, node, replacements, false)
			}
		}
		return node
	}
	{
		var c pkg_rhino.Node = node.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			c = this_.priv_FunctionArgumentInjector_inject_2900(compiler, c, node, replacements, replaceThis)
		}
	}
	return node
}

// Java method: FunctionArgumentInjector.getFunctionCallParameterMap(com.google.javascript.rhino.Node fnNode, com.google.javascript.rhino.Node callNode, java.util.function.Supplier<java.lang.String> safeNameIdSupplier)
func (this_ *FunctionArgumentInjector) GetFunctionCallParameterMap(fnNode pkg_rhino.Node, callNode pkg_rhino.Node, safeNameIdSupplier j2gort.Ext_java_util_function_Supplier) *FunctionArgumentInjector_ParamArgList {
	_ = fnNode
	_ = callNode
	_ = safeNameIdSupplier
	var val_1 FunctionArgumentInjector_ParamArgList
	InitFunctionArgumentInjector_ParamArgList(&val_1)
	var argMap *FunctionArgumentInjector_ParamArgList = &val_1
	_ = argMap
	this_.priv_FunctionArgumentInjector_populateFunctionCallParameterMap_2904(argMap, fnNode, callNode, safeNameIdSupplier)
	return argMap
}

// Java method: FunctionArgumentInjector.acquireFunctionCallParameterMap(com.google.javascript.rhino.Node fnNode, com.google.javascript.rhino.Node callNode, java.util.function.Supplier<java.lang.String> safeNameIdSupplier)
func (this_ *FunctionArgumentInjector) AcquireFunctionCallParameterMap(fnNode pkg_rhino.Node, callNode pkg_rhino.Node, safeNameIdSupplier j2gort.Ext_java_util_function_Supplier) *FunctionArgumentInjector_ParamArgList {
	_ = fnNode
	_ = callNode
	_ = safeNameIdSupplier
	var argMap *FunctionArgumentInjector_ParamArgList = nil
	_ = argMap
	if !(this_.paramArgListInUse) {
		argMap = this_.reusableParamArgList
		if argMap == nil {
			argMap = NewFunctionArgumentInjector_ParamArgList()
			this_.reusableParamArgList = argMap
		}
		this_.paramArgListInUse = true
	} else {
		argMap = NewFunctionArgumentInjector_ParamArgList()
	}
	this_.priv_FunctionArgumentInjector_populateFunctionCallParameterMap_2904(argMap, fnNode, callNode, safeNameIdSupplier)
	return argMap
}

// Java method: FunctionArgumentInjector.releaseParamArgList(com.google.javascript.jscomp.FunctionArgumentInjector.ParamArgList argMap)
func (this_ *FunctionArgumentInjector) ReleaseParamArgList(argMap *FunctionArgumentInjector_ParamArgList) {
	_ = argMap
	argMap.Reset()
	if argMap == this_.reusableParamArgList {
		this_.paramArgListInUse = false
	}
}

// Java method: FunctionArgumentInjector.populateFunctionCallParameterMap(com.google.javascript.jscomp.FunctionArgumentInjector.ParamArgList argMap, com.google.javascript.rhino.Node fnNode, com.google.javascript.rhino.Node callNode, java.util.function.Supplier<java.lang.String> safeNameIdSupplier)
func (this_ *FunctionArgumentInjector) priv_FunctionArgumentInjector_populateFunctionCallParameterMap_2904(argMap *FunctionArgumentInjector_ParamArgList, fnNode pkg_rhino.Node, callNode pkg_rhino.Node, safeNameIdSupplier j2gort.Ext_java_util_function_Supplier) {
	_ = argMap
	_ = fnNode
	_ = callNode
	_ = safeNameIdSupplier
	j2gort.CheckNotNullNode(fnNode)
	var firstParam pkg_rhino.Node = NodeUtil_GetFunctionParameters(fnNode).GetFirstChild()
	_ = firstParam
	var cArg pkg_rhino.Node = callNode.GetSecondChild()
	_ = cArg
	if (cArg != 0) && NodeUtil_IsFunctionObjectCall(callNode) {
		if !(fnNode.IsArrowFunction()) {
			argMap.PutFirst("this", StaticGet_FunctionArgumentInjector_THIS_MARKER_NODE(), cArg)
		} else {
			argMap.PutFirst("this", StaticGet_FunctionArgumentInjector_THIS_MARKER_NODE(), pkg_rhino.IR_ThisNode())
		}
		cArg = cArg.GetNext()
	} else {
		if !(!(NodeUtil_IsFunctionObjectApply(callNode))) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(callNode))
		}
		argMap.PutFirst("this", StaticGet_FunctionArgumentInjector_THIS_MARKER_NODE(), NodeUtil_NewUndefinedNode(callNode))
	}
	{
		var fnParam pkg_rhino.Node = firstParam
		_ = fnParam
		for ; fnParam != 0; fnParam = fnParam.GetNext() {
			if cArg != 0 {
				if fnParam.IsRest() {
					var restName pkg_rhino.Node = fnParam.GetOnlyChild()
					_ = restName
					if !(restName.IsName()) {
						j2gort.CheckStateMsg(false, j2gort.BoxNode(restName))
					}
					var array pkg_rhino.Node = pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{})
					_ = array
					array.SrcrefTreeIfMissing(cArg)
					for cArg != 0 {
						array.AddChildToBack(cArg.CloneTree_0())
						cArg = cArg.GetNext()
					}
					argMap.Put(restName.GetString(), restName, array)
					return
				} else {
					if !(fnParam.IsName()) {
						j2gort.CheckStateMsg(false, j2gort.BoxNode(fnParam))
					}
					argMap.Put(fnParam.GetString(), fnParam, cArg)
				}
				cArg = cArg.GetNext()
			} else {
				if fnParam.IsRest() {
					var restName_2 pkg_rhino.Node = fnParam.GetOnlyChild()
					_ = restName_2
					if !(restName_2.IsName()) {
						j2gort.CheckStateMsg(false, j2gort.BoxNode(fnParam))
					}
					var array_2 pkg_rhino.Node = pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{}).Srcref(fnParam)
					_ = array_2
					argMap.Put(restName_2.GetString(), restName_2, array_2)
				} else {
					if !(fnParam.IsName()) {
						j2gort.CheckStateMsg(false, j2gort.BoxNode(fnParam))
					}
					argMap.Put(fnParam.GetString(), fnParam, NodeUtil_NewUndefinedNode(callNode))
				}
			}
		}
	}
	for cArg != 0 {
		var uniquePlaceholder string = FunctionArgumentInjector_GetUniqueAnonymousParameterName(safeNameIdSupplier)
		_ = uniquePlaceholder
		argMap.Put(uniquePlaceholder, pkg_rhino.IR_Name(uniquePlaceholder), cArg)
		cArg = cArg.GetNext()
	}
}

// Java static method: FunctionArgumentInjector.getUniqueAnonymousParameterName(java.util.function.Supplier<java.lang.String> safeNameIdSupplier)
func FunctionArgumentInjector_GetUniqueAnonymousParameterName(safeNameIdSupplier j2gort.Ext_java_util_function_Supplier) string {
	clinit_FunctionArgumentInjector()
	_ = safeNameIdSupplier
	return ("JSCompiler_inline_anon_param_" + j2gort.StrValString(j2gort.CastString(safeNameIdSupplier.Get())))
}

// Java method: FunctionArgumentInjector.findModifiedParameters(com.google.javascript.rhino.Node fnNode)
func (this_ *FunctionArgumentInjector) FindModifiedParameters(fnNode pkg_rhino.Node) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = fnNode
	var firstParam pkg_rhino.Node = NodeUtil_GetFunctionParameters(fnNode).GetFirstChild()
	_ = firstParam
	if firstParam == 0 {
		return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	}
	var p0 string = j2gort.NullString
	_ = p0
	var p1 string = j2gort.NullString
	_ = p1
	var p2 string = j2gort.NullString
	_ = p2
	var p3 string = j2gort.NullString
	_ = p3
	var moreParams []string = nil
	_ = moreParams
	var paramCount int32 = int32(0)
	_ = paramCount
	{
		var n pkg_rhino.Node = firstParam
		_ = n
		for ; n != 0; n = n.GetNext() {
			if n.IsRest() {
				continue
			} else if (n.IsDefaultValue() || n.IsObjectPattern()) || n.IsArrayPattern() {
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Not supported: "+j2gort.StrVal(n))))
			} else {
				var s string = n.GetString()
				_ = s
				switch paramCount {
				case int32(0):
					p0 = s
				case int32(1):
					p1 = s
				case int32(2):
					p2 = s
				case int32(3):
					p3 = s
				default:
					if moreParams == nil {
						moreParams = j2gort.NewStringArray(int32(8))
						{
							a_2 := moreParams
							i_3 := int32(0)
							a_2[i_3] = p0
						}
						{
							a_4 := moreParams
							i_5 := int32(1)
							a_4[i_5] = p1
						}
						{
							a_6 := moreParams
							i_7 := int32(2)
							a_6[i_7] = p2
						}
						{
							a_8 := moreParams
							i_9 := int32(3)
							a_8[i_9] = p3
						}
					} else if paramCount == int32(len(moreParams)) {
						{

							arrSrc_10 := moreParams
							arrLen_11 := (paramCount * int32(2))
							j2gort.CheckNotNullSlice(arrSrc_10)
							if arrLen_11 < 0 {
								panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
							}
							arrOut_12 := j2gort.NewStringArray(arrLen_11)
							copy(arrOut_12, arrSrc_10)
							moreParams = arrOut_12
						}

					}
					{
						a_13 := moreParams
						i_14 := paramCount
						a_13[i_14] = s
					}
				}
				paramCount = paramCount + 1
			}
		}
	}
	if paramCount == int32(0) {
		return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	}
	var unsafe *j2gort.OrderedSet[string] = ordSet_FunctionArgumentInjector_CollectModifiedParameters(fnNode.GetLastChild(), p0, p1, p2, p3, moreParams, paramCount, false, nil)
	_ = unsafe
	if unsafe != nil {
		return unsafe.ToImmutableSet()
	} else {
		return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	}
}

// Java static method: FunctionArgumentInjector.matchesParameterName(java.lang.String name, java.lang.@org.jspecify.annotations.Nullable String p0, java.lang.@org.jspecify.annotations.Nullable String p1, java.lang.@org.jspecify.annotations.Nullable String p2, java.lang.@org.jspecify.annotations.Nullable String p3, java.lang.String @org.jspecify.annotations.Nullable [] moreParams, int paramCount)
func FunctionArgumentInjector_MatchesParameterName(name string, p0 string, p1 string, p2 string, p3 string, moreParams []string, paramCount int32) bool {
	clinit_FunctionArgumentInjector()
	_ = name
	_ = p0
	_ = p1
	_ = p2
	_ = p3
	_ = moreParams
	_ = paramCount
	if moreParams != nil {
		{
			var i int32 = int32(0)
			_ = i
			for ; i < paramCount; i = i + 1 {
				if j2gort.CheckNotNullString(moreParams[i]) == name {
					return true
				}
			}
		}
		return false
	}
	return ((((j2gort.CheckNotNullString(name) == p0) || ((p1 != j2gort.NullString) && (j2gort.CheckNotNullString(name) == p1))) || ((p2 != j2gort.NullString) && (j2gort.CheckNotNullString(name) == p2))) || ((p3 != j2gort.NullString) && (j2gort.CheckNotNullString(name) == p3)))
}

// Java static method: FunctionArgumentInjector.collectModifiedParameters(com.google.javascript.rhino.Node n, java.lang.@org.jspecify.annotations.Nullable String p0, java.lang.@org.jspecify.annotations.Nullable String p1, java.lang.@org.jspecify.annotations.Nullable String p2, java.lang.@org.jspecify.annotations.Nullable String p3, java.lang.String @org.jspecify.annotations.Nullable [] moreParams, int paramCount, boolean inInnerClosure, java.util.@org.jspecify.annotations.Nullable LinkedHashSet<java.lang.String> unsafe)
func ordSet_FunctionArgumentInjector_CollectModifiedParameters(n pkg_rhino.Node, p0 string, p1 string, p2 string, p3 string, moreParams []string, paramCount int32, inInnerClosure bool, unsafe *j2gort.OrderedSet[string]) *j2gort.OrderedSet[string] {
	clinit_FunctionArgumentInjector()
	_ = n
	_ = p0
	_ = p1
	_ = p2
	_ = p3
	_ = moreParams
	_ = paramCount
	_ = inInnerClosure
	_ = unsafe
	if NodeUtil_IsFunctionOrNonStaticClassFieldInitializer(n) {
		inInnerClosure = true
	}
	if n.IsName() {
		var name string = n.GetString()
		_ = name
		if FunctionArgumentInjector_MatchesParameterName(name, p0, p1, p2, p3, moreParams, paramCount) && (inInnerClosure || FunctionArgumentInjector_CanNameValueChange(n)) {
			if unsafe == nil {
				unsafe = &j2gort.OrderedSet[string]{}
			}
			unsafe.Add(name)
		}
	}
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			unsafe = ordSet_FunctionArgumentInjector_CollectModifiedParameters(c, p0, p1, p2, p3, moreParams, paramCount, inInnerClosure, unsafe)
		}
	}
	return unsafe
}

// Java static method: FunctionArgumentInjector.canNameValueChange(com.google.javascript.rhino.Node n)
func FunctionArgumentInjector_CanNameValueChange(n pkg_rhino.Node) bool {
	clinit_FunctionArgumentInjector()
	_ = n
	if !(NodeUtil_IsLValue(n)) {
		return false
	}
	if NodeUtil_IsDeclarationLValue(n) {
		var declaring pkg_rhino.Node = NodeUtil_GetDeclaringParent(n)
		_ = declaring
		return (!(declaring.IsConst()) && !(declaring.IsLet()))
	}
	return true
}

// Java method: FunctionArgumentInjector.gatherCallArgumentsNeedingTemps(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node fnNode, com.google.javascript.jscomp.FunctionArgumentInjector.ParamArgList argMap, com.google.common.collect.ImmutableSet<java.lang.String> modifiedParameters, com.google.javascript.jscomp.CodingConvention convention)
func (this_ *FunctionArgumentInjector) GatherCallArgumentsNeedingTemps(compiler AbstractCompiler, fnNode pkg_rhino.Node, argMap *FunctionArgumentInjector_ParamArgList, modifiedParameters j2gort.Ext_com_google_common_collect_ImmutableSet, convention CodingConvention) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = compiler
	_ = fnNode
	_ = argMap
	_ = modifiedParameters
	_ = convention
	if !(fnNode.IsFunction()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(fnNode))
	}
	if argMap.IsEmpty() {
		return modifiedParameters
	}
	var block pkg_rhino.Node = fnNode.GetLastChild()
	_ = block
	var requiresTempsUpToThisParameterName string = ""
	_ = requiresTempsUpToThisParameterName
	var argCount int32 = argMap.Size()
	_ = argCount
	var thisArg pkg_rhino.Node = argMap.GetThisArg()
	_ = thisArg
	var isTrivialBody bool = (!(block.HasChildren()) || (block.HasOneChild() && !(this_.BodyMayHaveConditionalCode(block.GetLastChild()))))
	_ = isTrivialBody
	var implicitThis bool = (NodeUtil_IsUndefined(thisArg) || (fnNode.IsArrowFunction() && thisArg.IsThis()))
	_ = implicitThis
	var hasMinimalParameters bool = (implicitThis && (argCount <= int32(2)))
	_ = hasMinimalParameters
	var collector *FunctionArgumentInjector_ReferencedAfterSideEffect = this_.priv_FunctionArgumentInjector_obtainCollector_2909(argMap, block)
	_ = collector
	func() {
		defer func() {
			this_.priv_FunctionArgumentInjector_releaseCollector_2910(collector)
		}()
		{
			var i int32 = int32(0)
			_ = i
			for ; i < argCount; i = i + 1 {
				var parameterName string = argMap.GetNameAt(i)
				_ = parameterName
				var cArg pkg_rhino.Node = argMap.GetArgAt(i)
				_ = cArg
				if !(modifiedParameters.IsEmpty()) && j2gort.SetContainsStr(modifiedParameters, parameterName) {
					requiresTempsUpToThisParameterName = parameterName
					continue
				}
				var requiresTemporary bool = false
				_ = requiresTemporary
				var references int32 = NodeUtil_GetNameReferenceCount(block, parameterName)
				_ = references
				var argSideEffects bool = compiler.GetAstAnalyzer().MayHaveSideEffects(cArg)
				_ = argSideEffects
				if !(argSideEffects) && (references == int32(0)) {
					requiresTemporary = false
				} else if (((isTrivialBody && hasMinimalParameters) && (references == int32(1))) && !(NodeUtil_CanBeSideEffected_1(cArg) && collector.IsParameterReferencedAfterSideEffect(parameterName))) && (!(argSideEffects) || !(collector.HasNamesReferencedBeforeParameter(parameterName))) {
					requiresTemporary = false
				} else if compiler.GetAstAnalyzer().MayEffectMutableState(cArg) && (references > int32(0)) {
					requiresTemporary = true
				} else if argSideEffects {
					requiresTemporary = true
				} else if NodeUtil_CanBeSideEffected_1(cArg) && collector.IsParameterReferencedAfterSideEffect(parameterName) {
					requiresTemporary = true
				} else if references > int32(1) {
					switch cArg.TokenSwitchOrd_() {
					case int32(32) /* Token.NAME */ :
						var name string = cArg.GetString()
						_ = name
						requiresTemporary = (convention.IsExported_2(name, true) || convention.IsExported_2(name, false))
					case int32(37) /* Token.THIS */ :
						requiresTemporary = false
					case int32(35) /* Token.STRINGLIT */ :
						requiresTemporary = (j2gort.JStrLen(cArg.GetString()) >= int32(2))
					default:
						requiresTemporary = !(NodeUtil_IsImmutableValue(cArg))
					}
				}
				if requiresTemporary {
					requiresTempsUpToThisParameterName = parameterName
				}
			}
		}
	}()
	if j2gort.CheckNotNullString(requiresTempsUpToThisParameterName) == "" {
		return modifiedParameters
	}
	var namesNeedingTemps *j2gort.OrderedSet[string]
	_ = namesNeedingTemps
	if modifiedParameters.IsEmpty() {
		namesNeedingTemps = &j2gort.OrderedSet[string]{}
	} else {
		namesNeedingTemps = j2gort.OrderedSetPtrFromStrCol(j2gort.AsIface[j2gort.Ext_java_util_Collection](modifiedParameters))
	}
	{
		var i_2 int32 = int32(0)
		_ = i_2
		for ; i_2 < argCount; i_2 = i_2 + 1 {
			var parameterName_2 string = argMap.GetNameAt(i_2)
			_ = parameterName_2
			if (i_2 == int32(0)) && implicitThis {
				continue
			}
			namesNeedingTemps.Add(parameterName_2)
			if j2gort.CheckNotNullString(parameterName_2) == requiresTempsUpToThisParameterName {
				break
			}
		}
	}
	return j2gort.CheckNotNullPtr(namesNeedingTemps).ToImmutableSet()
}

// Java method: FunctionArgumentInjector.bodyMayHaveConditionalCode(com.google.javascript.rhino.Node n)
func (this_ *FunctionArgumentInjector) BodyMayHaveConditionalCode(n pkg_rhino.Node) bool {
	_ = n
	if !(n.IsReturn()) && !(n.IsExprResult()) {
		return true
	}
	return this_.MayHaveConditionalCode(n)
}

// Java method: FunctionArgumentInjector.mayHaveConditionalCode(com.google.javascript.rhino.Node n)
func (this_ *FunctionArgumentInjector) MayHaveConditionalCode(n pkg_rhino.Node) bool {
	_ = n
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			switch c.TokenSwitchOrd_() {
			case int32(73) /* Token.FUNCTION */, int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(67) /* Token.HOOK */, int32(70) /* Token.COALESCE */, int32(31) /* Token.OPTCHAIN_CALL */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
				return true
			default:
				if NodeUtil_IsNonStaticClassFieldInitializer(c) {
					return true
				}
			}
			if this_.MayHaveConditionalCode(c) {
				return true
			}
		}
	}
	return false
}

// Java method: FunctionArgumentInjector.obtainCollector(com.google.javascript.jscomp.FunctionArgumentInjector.ParamArgList parameters, com.google.javascript.rhino.Node root)
func (this_ *FunctionArgumentInjector) priv_FunctionArgumentInjector_obtainCollector_2909(parameters *FunctionArgumentInjector_ParamArgList, root pkg_rhino.Node) *FunctionArgumentInjector_ReferencedAfterSideEffect {
	_ = parameters
	_ = root
	var extraLocals *j2gort.Ext_java_util_LinkedHashSet = FunctionArgumentInjector_GatherLocalNames(root, nil)
	_ = extraLocals
	var collector *FunctionArgumentInjector_ReferencedAfterSideEffect = nil
	_ = collector
	if !(this_.collectorInUse) {
		collector = this_.reusableCollector
		if collector == nil {
			collector = NewFunctionArgumentInjector_ReferencedAfterSideEffect(this_)
			this_.reusableCollector = collector
		}
		this_.collectorInUse = true
	} else {
		collector = NewFunctionArgumentInjector_ReferencedAfterSideEffect(this_)
	}
	collector.Reset(parameters, extraLocals)
	collector.Traverse(root)
	return collector
}

// Java method: FunctionArgumentInjector.releaseCollector(com.google.javascript.jscomp.FunctionArgumentInjector.ReferencedAfterSideEffect collector)
func (this_ *FunctionArgumentInjector) priv_FunctionArgumentInjector_releaseCollector_2910(collector *FunctionArgumentInjector_ReferencedAfterSideEffect) {
	_ = collector
	collector.ClearRefs()
	if collector == this_.reusableCollector {
		this_.collectorInUse = false
	}
}

// Java static method: FunctionArgumentInjector.gatherLocalNames(com.google.javascript.rhino.Node n, java.util.@org.jspecify.annotations.Nullable LinkedHashSet<java.lang.String> names)
func FunctionArgumentInjector_GatherLocalNames(n pkg_rhino.Node, names *j2gort.Ext_java_util_LinkedHashSet) *j2gort.Ext_java_util_LinkedHashSet {
	clinit_FunctionArgumentInjector()
	_ = n
	_ = names
	if n.IsFunction() {
		if NodeUtil_IsFunctionDeclaration(n) {
			if names == nil {
				names = j2gort.NewExt_java_util_LinkedHashSet()
			}
			j2gort.SetAddStr(names, n.GetFirstChild().GetString())
		}
		return names
	} else if n.IsName() {
		switch n.GetParent().TokenSwitchOrd_() {
		case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(86) /* Token.CATCH */ :
			if names == nil {
				names = j2gort.NewExt_java_util_LinkedHashSet()
			}
			j2gort.SetAddStr(names, n.GetString())
		default:
		}
	}
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			names = FunctionArgumentInjector_GatherLocalNames(c, names)
		}
	}
	return names
}

// Java class: com.google.javascript.jscomp.FunctionArgumentInjector.ReferencedAfterSideEffect
type FunctionArgumentInjector_ReferencedAfterSideEffect struct {
	outer_                                       *FunctionArgumentInjector
	parameters                                   *FunctionArgumentInjector_ParamArgList `gson:"parameters"`
	extraLocals                                  *j2gort.Ext_java_util_LinkedHashSet    `gson:"extraLocals"`
	parametersReferenced                         *j2gort.OrderedSet[string]             `gson:"parametersReferenced"`
	parametersWithNamesReferencedBeforeParameter *j2gort.OrderedSet[string]             `gson:"parametersWithNamesReferencedBeforeParameter"`
	loopsEntered                                 int32                                  `gson:"loopsEntered"`
	idHash_                                      atomic.Int32
	sideEffectSeen                               bool `gson:"sideEffectSeen"`
	nameNodeHasBeenSeen                          bool `gson:"nameNodeHasBeenSeen"`
}

func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) IsJavaNil_() bool {
	return this_ == nil
}
func (*FunctionArgumentInjector_ReferencedAfterSideEffect) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionArgumentInjector$ReferencedAfterSideEffect", "ReferencedAfterSideEffect"
}
func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FunctionArgumentInjector_ReferencedAfterSideEffect) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionArgumentInjector.ReferencedAfterSideEffect", "com.google.javascript.jscomp.FunctionArgumentInjector$ReferencedAfterSideEffect", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionArgumentInjector_ReferencedAfterSideEffect](v_)
	})
}
func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) SetOuter_(o *FunctionArgumentInjector) {
	this_.outer_ = o
}
func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) GetOuter_() *FunctionArgumentInjector {
	return this_.outer_
}

func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) initFields_FunctionArgumentInjector_ReferencedAfterSideEffect() {
	_ = this_

}

// Java constructor: ReferencedAfterSideEffect()
func NewFunctionArgumentInjector_ReferencedAfterSideEffect(outer_ *FunctionArgumentInjector) *FunctionArgumentInjector_ReferencedAfterSideEffect {
	this_ := &FunctionArgumentInjector_ReferencedAfterSideEffect{}
	this_.outer_ = outer_
	InitFunctionArgumentInjector_ReferencedAfterSideEffect(this_)
	return this_
}

func InitFunctionArgumentInjector_ReferencedAfterSideEffect(this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) {
	this_.initFields_FunctionArgumentInjector_ReferencedAfterSideEffect()
}

// Java method: ReferencedAfterSideEffect.reset(com.google.javascript.jscomp.FunctionArgumentInjector.ParamArgList parameters, java.util.@org.jspecify.annotations.Nullable LinkedHashSet<java.lang.String> extraLocals)
func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) Reset(parameters *FunctionArgumentInjector_ParamArgList, extraLocals *j2gort.Ext_java_util_LinkedHashSet) {
	_ = parameters
	_ = extraLocals
	this_.parameters = parameters
	this_.extraLocals = extraLocals
	this_.sideEffectSeen = false
	if this_.parametersReferenced != nil {
		this_.parametersReferenced.Clear()
	}
	this_.nameNodeHasBeenSeen = false
	if this_.parametersWithNamesReferencedBeforeParameter != nil {
		this_.parametersWithNamesReferencedBeforeParameter.Clear()
	}
	this_.loopsEntered = int32(0)
}

// Java method: ReferencedAfterSideEffect.clearRefs()
func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) ClearRefs() {
	this_.parameters = nil
	this_.extraLocals = nil
}

// Java method: ReferencedAfterSideEffect.isParameterReferencedAfterSideEffect(java.lang.String name)
func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) IsParameterReferencedAfterSideEffect(name string) bool {
	_ = name
	return ((this_.parametersReferenced != nil) && this_.parametersReferenced.Contains(name))
}

// Java method: ReferencedAfterSideEffect.hasNamesReferencedBeforeParameter(java.lang.String name)
func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) HasNamesReferencedBeforeParameter(name string) bool {
	_ = name
	return ((this_.parametersWithNamesReferencedBeforeParameter != nil) && this_.parametersWithNamesReferencedBeforeParameter.Contains(name))
}

// Java method: ReferencedAfterSideEffect.traverse(com.google.javascript.rhino.Node node)
func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) Traverse(node pkg_rhino.Node) {
	_ = node
	if this_.priv_ReferencedAfterSideEffect_shouldTraverseChildren_2916(node) {
		{
			var c pkg_rhino.Node = node.GetFirstChild()
			_ = c
			for c != 0 {
				var next pkg_rhino.Node = c.GetNext()
				_ = next
				this_.Traverse(c)
				c = next
			}
		}
	}
	this_.Visit(node)
}

// Java method: ReferencedAfterSideEffect.shouldTraverseChildren(com.google.javascript.rhino.Node node)
func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) priv_ReferencedAfterSideEffect_shouldTraverseChildren_2916(node pkg_rhino.Node) bool {
	_ = node
	if NodeUtil_IsLoopStructure(node) {
		this_.loopsEntered = this_.loopsEntered + 1
	}
	return !((this_.sideEffectSeen && (this_.parametersReferenced != nil)) && (this_.parameters.Size() == this_.parametersReferenced.Size()))
}

// Java method: ReferencedAfterSideEffect.inLoop()
func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) InLoop() bool {
	return (this_.loopsEntered != int32(0))
}

// Java method: ReferencedAfterSideEffect.recordReferenced(java.lang.String name)
func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) priv_ReferencedAfterSideEffect_recordReferenced_2918(name string) {
	_ = name
	if this_.parametersReferenced == nil {
		this_.parametersReferenced = &j2gort.OrderedSet[string]{}
	}
	this_.parametersReferenced.Add(name)
}

// Java method: ReferencedAfterSideEffect.visit(com.google.javascript.rhino.Node n)
func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) Visit(n pkg_rhino.Node) {
	_ = n
	if n.IsName() {
		var name string = n.GetString()
		_ = name
		if this_.nameNodeHasBeenSeen && this_.parameters.ContainsName(name) {
			if this_.parametersWithNamesReferencedBeforeParameter == nil {
				this_.parametersWithNamesReferencedBeforeParameter = &j2gort.OrderedSet[string]{}
			}
			this_.parametersWithNamesReferencedBeforeParameter.Add(name)
		}
		this_.nameNodeHasBeenSeen = true
	}
	if NodeUtil_IsLoopStructure(n) {
		this_.loopsEntered = this_.loopsEntered - 1
		if (!(this_.InLoop()) && !(this_.sideEffectSeen)) && (this_.parametersReferenced != nil) {
			this_.parametersReferenced.Clear()
		}
	}
	if !(this_.sideEffectSeen) && this_.priv_ReferencedAfterSideEffect_hasNonLocalSideEffect_2920(n) {
		this_.sideEffectSeen = true
	}
	if this_.InLoop() || this_.sideEffectSeen {
		if n.IsName() {
			var name_2 string = n.GetString()
			_ = name_2
			if this_.parameters.ContainsName(name_2) {
				this_.priv_ReferencedAfterSideEffect_recordReferenced_2918(name_2)
			}
		} else if n.IsThis() {
			this_.priv_ReferencedAfterSideEffect_recordReferenced_2918("this")
		}
	}
}

// Java method: ReferencedAfterSideEffect.hasNonLocalSideEffect(com.google.javascript.rhino.Node n)
func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) priv_ReferencedAfterSideEffect_hasNonLocalSideEffect_2920(n pkg_rhino.Node) bool {
	_ = n
	var sideEffect bool = false
	_ = sideEffect
	var type_ pkg_rhino.Token = n.GetToken()
	_ = type_
	if (NodeUtil_IsAssignmentOp(n) || (type_ == pkg_rhino.StaticGet_Token_INC())) || (type_ == pkg_rhino.StaticGet_Token_DEC()) {
		var lhs pkg_rhino.Node = n.GetFirstChild()
		_ = lhs
		if !(this_.priv_ReferencedAfterSideEffect_isLocalName_2921(lhs)) {
			sideEffect = true
		}
	} else if type_ == pkg_rhino.StaticGet_Token_CALL() {
		sideEffect = this_.outer_.astAnalyzer.FunctionCallHasSideEffects(n)
	} else if type_ == pkg_rhino.StaticGet_Token_NEW() {
		sideEffect = this_.outer_.astAnalyzer.ConstructorCallHasSideEffects(n)
	} else if type_ == pkg_rhino.StaticGet_Token_DELPROP() {
		sideEffect = true
	}
	return sideEffect
}

// Java method: ReferencedAfterSideEffect.isLocalName(com.google.javascript.rhino.Node node)
func (this_ *FunctionArgumentInjector_ReferencedAfterSideEffect) priv_ReferencedAfterSideEffect_isLocalName_2921(node pkg_rhino.Node) bool {
	_ = node
	if node.IsName() {
		var name string = node.GetString()
		_ = name
		return (this_.parameters.ContainsName(name) || ((this_.extraLocals != nil) && j2gort.SetContainsStr(this_.extraLocals, name)))
	}
	return false
}

// Java class: com.google.javascript.jscomp.FunctionArgumentInjector.ParamArgList
type FunctionArgumentInjector_ParamArgList struct {
	n0             string           `gson:"n0"`
	n1             string           `gson:"n1"`
	n2             string           `gson:"n2"`
	n3             string           `gson:"n3"`
	n4             string           `gson:"n4"`
	restNames      []string         `gson:"restNames"`
	restParamNodes []pkg_rhino.Node `gson:"restParamNodes"`
	restArgs       []pkg_rhino.Node `gson:"restArgs"`
	p0             pkg_rhino.Node   `gson:"p0"`
	p1             pkg_rhino.Node   `gson:"p1"`
	p2             pkg_rhino.Node   `gson:"p2"`
	p3             pkg_rhino.Node   `gson:"p3"`
	p4             pkg_rhino.Node   `gson:"p4"`
	a0             pkg_rhino.Node   `gson:"a0"`
	a1             pkg_rhino.Node   `gson:"a1"`
	a2             pkg_rhino.Node   `gson:"a2"`
	a3             pkg_rhino.Node   `gson:"a3"`
	a4             pkg_rhino.Node   `gson:"a4"`
	size           int32            `gson:"size"`
	idHash_        atomic.Int32
}

func (this_ *FunctionArgumentInjector_ParamArgList) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionArgumentInjector_ParamArgList) IsJavaNil_() bool { return this_ == nil }
func (*FunctionArgumentInjector_ParamArgList) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionArgumentInjector$ParamArgList", "ParamArgList"
}
func (this_ *FunctionArgumentInjector_ParamArgList) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FunctionArgumentInjector_ParamArgList) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionArgumentInjector.ParamArgList", "com.google.javascript.jscomp.FunctionArgumentInjector$ParamArgList", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionArgumentInjector_ParamArgList](v_)
	})
}

func (this_ *FunctionArgumentInjector_ParamArgList) initFields_FunctionArgumentInjector_ParamArgList() {
	_ = this_
	this_.n0 = j2gort.NullString
	this_.n1 = j2gort.NullString
	this_.n2 = j2gort.NullString
	this_.n3 = j2gort.NullString
	this_.n4 = j2gort.NullString
}

// Java constructor: ParamArgList()
func NewFunctionArgumentInjector_ParamArgList() *FunctionArgumentInjector_ParamArgList {
	this_ := &FunctionArgumentInjector_ParamArgList{}
	InitFunctionArgumentInjector_ParamArgList(this_)
	return this_
}

func InitFunctionArgumentInjector_ParamArgList(this_ *FunctionArgumentInjector_ParamArgList) {
	this_.initFields_FunctionArgumentInjector_ParamArgList()
}

// Java method: ParamArgList.reset()
func (this_ *FunctionArgumentInjector_ParamArgList) Reset() {
	var n int32 = this_.size
	_ = n
	if n == int32(0) {
		return
	}
	this_.n0 = j2gort.NullString
	this_.p0 = pkg_rhino.Node(0)
	this_.a0 = pkg_rhino.Node(0)
	if n > int32(1) {
		this_.n1 = j2gort.NullString
		this_.p1 = pkg_rhino.Node(0)
		this_.a1 = pkg_rhino.Node(0)
	}
	if n > int32(2) {
		this_.n2 = j2gort.NullString
		this_.p2 = pkg_rhino.Node(0)
		this_.a2 = pkg_rhino.Node(0)
	}
	if n > int32(3) {
		this_.n3 = j2gort.NullString
		this_.p3 = pkg_rhino.Node(0)
		this_.a3 = pkg_rhino.Node(0)
	}
	if n > int32(4) {
		this_.n4 = j2gort.NullString
		this_.p4 = pkg_rhino.Node(0)
		this_.a4 = pkg_rhino.Node(0)
	}
	if (n > int32(5)) && (this_.restNames != nil) {
		var extra int32 = (n - int32(5))
		_ = extra
		j2gort.ArraysFillRange(this_.restNames, int32(0), extra, j2gort.NullString)
		j2gort.ArraysFillRange(this_.restParamNodes, int32(0), extra, pkg_rhino.Node(0))
		j2gort.ArraysFillRange(this_.restArgs, int32(0), extra, pkg_rhino.Node(0))
	}
	this_.size = int32(0)
}

// Java method: ParamArgList.putFirst(java.lang.String name, com.google.javascript.rhino.Node paramNode, com.google.javascript.rhino.Node arg)
func (this_ *FunctionArgumentInjector_ParamArgList) PutFirst(name string, paramNode pkg_rhino.Node, arg pkg_rhino.Node) {
	_ = name
	_ = paramNode
	_ = arg
	j2gort.CheckNotNullString(name)
	j2gort.CheckNotNullNode(paramNode)
	j2gort.CheckNotNullNode(arg)
	this_.n0 = name
	this_.p0 = paramNode
	this_.a0 = arg
	this_.size = int32(1)
}

// Java method: ParamArgList.put(java.lang.String name, com.google.javascript.rhino.Node paramNode, com.google.javascript.rhino.Node arg)
func (this_ *FunctionArgumentInjector_ParamArgList) Put(name string, paramNode pkg_rhino.Node, arg pkg_rhino.Node) {
	_ = name
	_ = paramNode
	_ = arg
	j2gort.CheckNotNullString(name)
	j2gort.CheckNotNullNode(paramNode)
	j2gort.CheckNotNullNode(arg)
	if this_.ContainsName(name) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Multiple entries with same key: "+j2gort.StrValString(name))))
	}
	var n int32 = this_.size
	_ = n
	switch n {
	case int32(0):
		this_.n0 = name
		this_.p0 = paramNode
		this_.a0 = arg
	case int32(1):
		this_.n1 = name
		this_.p1 = paramNode
		this_.a1 = arg
	case int32(2):
		this_.n2 = name
		this_.p2 = paramNode
		this_.a2 = arg
	case int32(3):
		this_.n3 = name
		this_.p3 = paramNode
		this_.a3 = arg
	case int32(4):
		this_.n4 = name
		this_.p4 = paramNode
		this_.a4 = arg
	default:
		var restIdx int32 = (n - int32(5))
		_ = restIdx
		if this_.restNames == nil {
			this_.restNames = j2gort.NewStringArray(int32(4))
			this_.restParamNodes = make([]pkg_rhino.Node, int32(4))
			this_.restArgs = make([]pkg_rhino.Node, int32(4))
		} else if restIdx == int32(len(this_.restNames)) {
			var nextCap int32 = (restIdx * int32(2))
			_ = nextCap
			this_.restNames = func() []string {
				arrSrc_15 := this_.restNames
				arrLen_16 := nextCap
				j2gort.CheckNotNullSlice(arrSrc_15)
				if arrLen_16 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_17 := j2gort.NewStringArray(arrLen_16)
				copy(arrOut_17, arrSrc_15)
				return arrOut_17
			}()
			this_.restParamNodes = func() []pkg_rhino.Node {
				arrSrc_18 := this_.restParamNodes
				arrLen_19 := nextCap
				j2gort.CheckNotNullSlice(arrSrc_18)
				if arrLen_19 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_20 := make([]pkg_rhino.Node, arrLen_19)
				copy(arrOut_20, arrSrc_18)
				return arrOut_20
			}()
			this_.restArgs = func() []pkg_rhino.Node {
				arrSrc_21 := this_.restArgs
				arrLen_22 := nextCap
				j2gort.CheckNotNullSlice(arrSrc_21)
				if arrLen_22 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_23 := make([]pkg_rhino.Node, arrLen_22)
				copy(arrOut_23, arrSrc_21)
				return arrOut_23
			}()
		}
		{
			a_24 := this_.restNames
			i_25 := restIdx
			a_24[i_25] = name
		}
		{
			a_26 := this_.restParamNodes
			i_27 := restIdx
			a_26[i_27] = paramNode
		}
		{
			a_28 := this_.restArgs
			i_29 := restIdx
			a_28[i_29] = arg
		}
	}
	this_.size = (n + int32(1))
}

// Java method: ParamArgList.isEmpty()
func (this_ *FunctionArgumentInjector_ParamArgList) IsEmpty() bool {
	return (this_.size == int32(0))
}

// Java method: ParamArgList.size()
func (this_ *FunctionArgumentInjector_ParamArgList) Size() int32 {
	return this_.size
}

// Java method: ParamArgList.getThisArg()
func (this_ *FunctionArgumentInjector_ParamArgList) GetThisArg() pkg_rhino.Node {
	return this_.a0
}

// Java method: ParamArgList.clearThisArg()
func (this_ *FunctionArgumentInjector_ParamArgList) ClearThisArg() {
	this_.a0 = pkg_rhino.Node(0)
}

// Java method: ParamArgList.getNameAt(int index)
func (this_ *FunctionArgumentInjector_ParamArgList) GetNameAt(index int32) string {
	_ = index
	switch index {
	case int32(0):
		return this_.n0
	case int32(1):
		return this_.n1
	case int32(2):
		return this_.n2
	case int32(3):
		return this_.n3
	case int32(4):
		return this_.n4
	default:
		return this_.restNames[(index - int32(5))]
	}
	panic("unreachable switch expression")
}

// Java method: ParamArgList.getParamNodeAt(int index)
func (this_ *FunctionArgumentInjector_ParamArgList) GetParamNodeAt(index int32) pkg_rhino.Node {
	_ = index
	switch index {
	case int32(0):
		return this_.p0
	case int32(1):
		return this_.p1
	case int32(2):
		return this_.p2
	case int32(3):
		return this_.p3
	case int32(4):
		return this_.p4
	default:
		return this_.restParamNodes[(index - int32(5))]
	}
	panic("unreachable switch expression")
}

// Java method: ParamArgList.getArgAt(int index)
func (this_ *FunctionArgumentInjector_ParamArgList) GetArgAt(index int32) pkg_rhino.Node {
	_ = index
	switch index {
	case int32(0):
		return this_.a0
	case int32(1):
		return this_.a1
	case int32(2):
		return this_.a2
	case int32(3):
		return this_.a3
	case int32(4):
		return this_.a4
	default:
		return this_.restArgs[(index - int32(5))]
	}
	panic("unreachable switch expression")
}

// Java method: ParamArgList.setArgAt(int index, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node arg)
func (this_ *FunctionArgumentInjector_ParamArgList) SetArgAt(index int32, arg pkg_rhino.Node) {
	_ = index
	_ = arg
	switch index {
	case int32(0):
		this_.a0 = arg
	case int32(1):
		this_.a1 = arg
	case int32(2):
		this_.a2 = arg
	case int32(3):
		this_.a3 = arg
	case int32(4):
		this_.a4 = arg
	default:
		{
			a_30 := this_.restArgs
			i_31 := (index - int32(5))
			a_30[i_31] = arg
		}
	}
}

// Java method: ParamArgList.containsName(java.lang.String name)
func (this_ *FunctionArgumentInjector_ParamArgList) ContainsName(name string) bool {
	_ = name
	var n int32 = this_.size
	_ = n
	if n == int32(0) {
		return false
	}
	{
		var iifeCond_1 bool
		{

			sr_32 := name
			se_33 := this_.n0
			iifeCond_1 = j2gort.CheckNotNullString(sr_32) == se_33
		}
		if iifeCond_1 {
			return true
		}
	}
	if (n > int32(1)) && func() bool { sr_34 := name; se_35 := this_.n1; return j2gort.CheckNotNullString(sr_34) == se_35 }() {
		return true
	}
	if (n > int32(2)) && func() bool { sr_36 := name; se_37 := this_.n2; return j2gort.CheckNotNullString(sr_36) == se_37 }() {
		return true
	}
	if (n > int32(3)) && func() bool { sr_38 := name; se_39 := this_.n3; return j2gort.CheckNotNullString(sr_38) == se_39 }() {
		return true
	}
	if (n > int32(4)) && func() bool { sr_40 := name; se_41 := this_.n4; return j2gort.CheckNotNullString(sr_40) == se_41 }() {
		return true
	}
	if n > int32(5) {
		var rest []string = this_.restNames
		_ = rest
		{
			var i int32 = int32(0)
			_ = i
			var extra int32 = (n - int32(5))
			_ = extra
			for ; i < extra; i = i + 1 {
				{
					var iifeCond_1 bool
					{

						sr_42 := name
						se_43 := rest[i]
						iifeCond_1 = j2gort.CheckNotNullString(sr_42) == se_43
					}
					if iifeCond_1 {
						return true
					}
				}
			}
		}
	}
	return false
}

// Java method: ParamArgList.getArg(java.lang.String name)
func (this_ *FunctionArgumentInjector_ParamArgList) GetArg(name string) pkg_rhino.Node {
	_ = name
	var n int32 = this_.size
	_ = n
	if n == int32(0) {
		return pkg_rhino.Node(0)
	}
	{
		var iifeCond_1 bool
		{

			sr_44 := name
			se_45 := this_.n0
			iifeCond_1 = j2gort.CheckNotNullString(sr_44) == se_45
		}
		if iifeCond_1 {
			return this_.a0
		}
	}
	if (n > int32(1)) && func() bool { sr_46 := name; se_47 := this_.n1; return j2gort.CheckNotNullString(sr_46) == se_47 }() {
		return this_.a1
	}
	if (n > int32(2)) && func() bool { sr_48 := name; se_49 := this_.n2; return j2gort.CheckNotNullString(sr_48) == se_49 }() {
		return this_.a2
	}
	if (n > int32(3)) && func() bool { sr_50 := name; se_51 := this_.n3; return j2gort.CheckNotNullString(sr_50) == se_51 }() {
		return this_.a3
	}
	if (n > int32(4)) && func() bool { sr_52 := name; se_53 := this_.n4; return j2gort.CheckNotNullString(sr_52) == se_53 }() {
		return this_.a4
	}
	if n > int32(5) {
		var rest []string = this_.restNames
		_ = rest
		{
			var i int32 = int32(0)
			_ = i
			var extra int32 = (n - int32(5))
			_ = extra
			for ; i < extra; i = i + 1 {
				{
					var iifeCond_1 bool
					{

						sr_54 := name
						se_55 := rest[i]
						iifeCond_1 = j2gort.CheckNotNullString(sr_54) == se_55
					}
					if iifeCond_1 {
						return this_.restArgs[i]
					}
				}
			}
		}
	}
	return pkg_rhino.Node(0)
}
func FunctionArgumentInjector_CollectModifiedParameters(n pkg_rhino.Node, p0 string, p1 string, p2 string, p3 string, moreParams []string, paramCount int32, inInnerClosure bool, unsafe *j2gort.Ext_java_util_LinkedHashSet) *j2gort.Ext_java_util_LinkedHashSet {
	clinit_FunctionArgumentInjector()
	_ = n
	_ = p0
	_ = p1
	_ = p2
	_ = p3
	_ = moreParams
	_ = paramCount
	_ = inInnerClosure
	_ = unsafe
	if NodeUtil_IsFunctionOrNonStaticClassFieldInitializer(n) {
		inInnerClosure = true
	}
	if n.IsName() {
		var name string = n.GetString()
		_ = name
		if FunctionArgumentInjector_MatchesParameterName(name, p0, p1, p2, p3, moreParams, paramCount) && (inInnerClosure || FunctionArgumentInjector_CanNameValueChange(n)) {
			if unsafe == nil {
				unsafe = j2gort.NewExt_java_util_LinkedHashSet()
			}
			j2gort.SetAddStr(unsafe, name)
		}
	}
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			unsafe = FunctionArgumentInjector_CollectModifiedParameters(c, p0, p1, p2, p3, moreParams, paramCount, inInnerClosure, unsafe)
		}
	}
	return unsafe
}
