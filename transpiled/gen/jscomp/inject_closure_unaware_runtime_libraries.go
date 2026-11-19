/*
 * Copyright 2025 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/InjectClosureUnawareRuntimeLibraries.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_inject_closure_unaware_runtime_libraries_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_inject_closure_unaware_runtime_libraries_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.BoxString(j2gort.AsIface[RuntimeJsLibManager_InjectedJsLibField](mrp0_).QualifiedName())
	})
}

// Java class: com.google.javascript.jscomp.InjectClosureUnawareRuntimeLibraries
type InjectClosureUnawareRuntimeLibraries struct {
	mainCompiler        AbstractCompiler          `gson:"mainCompiler"`
	shadowCompiler      AbstractCompiler          `gson:"shadowCompiler"`
	closureUnawareCalls j2gort.Ext_java_util_List `gson:"closureUnawareCalls"`
	idHash_             atomic.Int32
}

func (this_ *InjectClosureUnawareRuntimeLibraries) BindSelf_(s any) {
	_ = s
}
func (this_ *InjectClosureUnawareRuntimeLibraries) IsJavaNil_() bool { return this_ == nil }
func (*InjectClosureUnawareRuntimeLibraries) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InjectClosureUnawareRuntimeLibraries", "InjectClosureUnawareRuntimeLibraries"
}
func (this_ *InjectClosureUnawareRuntimeLibraries) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InjectClosureUnawareRuntimeLibraries) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InjectClosureUnawareRuntimeLibraries", "com.google.javascript.jscomp.InjectClosureUnawareRuntimeLibraries", func(v_ any) bool {
		return j2gort.InstanceOf[*InjectClosureUnawareRuntimeLibraries](v_)
	})
}
func (*InjectClosureUnawareRuntimeLibraries) Implements_jscomp_CompilerPass() {}

func (this_ *InjectClosureUnawareRuntimeLibraries) initFields_InjectClosureUnawareRuntimeLibraries() {
	_ = this_
}

// Java constructor: InjectClosureUnawareRuntimeLibraries(com.google.javascript.jscomp.AbstractCompiler mainCompiler, com.google.javascript.jscomp.AbstractCompiler shadowCompiler, java.util.List<com.google.javascript.jscomp.InjectClosureUnawareRuntimeLibraries.ClosureUnawareCallSite> closureUnawareCalls)
func NewInjectClosureUnawareRuntimeLibraries(mainCompiler AbstractCompiler, shadowCompiler AbstractCompiler, closureUnawareCalls j2gort.Ext_java_util_List) *InjectClosureUnawareRuntimeLibraries {
	this_ := &InjectClosureUnawareRuntimeLibraries{}
	InitInjectClosureUnawareRuntimeLibraries(this_, mainCompiler, shadowCompiler, closureUnawareCalls)
	return this_
}

func InitInjectClosureUnawareRuntimeLibraries(this_ *InjectClosureUnawareRuntimeLibraries, mainCompiler AbstractCompiler, shadowCompiler AbstractCompiler, closureUnawareCalls j2gort.Ext_java_util_List) {
	_ = mainCompiler
	_ = shadowCompiler
	_ = closureUnawareCalls
	this_.initFields_InjectClosureUnawareRuntimeLibraries()
	this_.mainCompiler = mainCompiler
	this_.shadowCompiler = shadowCompiler
	this_.closureUnawareCalls = closureUnawareCalls
}

// Java method: InjectClosureUnawareRuntimeLibraries.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InjectClosureUnawareRuntimeLibraries) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.priv_InjectClosureUnawareRuntimeLibraries_injectLibrariesIntoMainAst_3201()
	this_.priv_InjectClosureUnawareRuntimeLibraries_injectFieldsIntoCallSites_3202()
}

// Java method: InjectClosureUnawareRuntimeLibraries.injectLibrariesIntoMainAst()
func (this_ *InjectClosureUnawareRuntimeLibraries) priv_InjectClosureUnawareRuntimeLibraries_injectLibrariesIntoMainAst_3201() {
	var mainManager *RuntimeJsLibManager = this_.mainCompiler.GetRuntimeJsLibManager()
	_ = mainManager
	var shadowManager *RuntimeJsLibManager = this_.shadowCompiler.GetRuntimeJsLibManager()
	_ = shadowManager
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(shadowManager.GetInjectedLibraries()))
		for cur_1.HasNext() {
			runtimeLib := cur_1.NextStr()
			_ = runtimeLib
			mainManager.EnsureLibraryInjected(runtimeLib, false)
		}
	}
}

// Java method: InjectClosureUnawareRuntimeLibraries.injectFieldsIntoCallSites()
func (this_ *InjectClosureUnawareRuntimeLibraries) priv_InjectClosureUnawareRuntimeLibraries_injectFieldsIntoCallSites_3202() {
	var fieldsByName j2gort.Ext_com_google_common_collect_ImmutableMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](this_.shadowCompiler.GetRuntimeJsLibManager().GetExternedFields().Stream().Collect_1(j2gort.Ext_com_google_common_collect_ImmutableMap_ToImmutableMap_2(lambdaSingleton_inject_closure_unaware_runtime_libraries_1, j2gort.Ext_java_util_function_Function_Identity())))
	_ = fieldsByName
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.closureUnawareCalls))
		for cur_2.HasNext() {
			rawElem_3 := cur_2.Next()
			callSite := j2gort.CastPtr[InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite](rawElem_3)
			_ = callSite
			var toInject j2gort.OrderedSet[RuntimeJsLibManager_ExternedField] = ordSet_InjectClosureUnawareRuntimeLibraries_GatherRuntimeFields(callSite.ClosureUnawareFunction(), fieldsByName)
			_ = toInject
			this_.priv_InjectClosureUnawareRuntimeLibraries_injectAll_3203(callSite, &toInject)
		}
	}
}

// Java method: InjectClosureUnawareRuntimeLibraries.injectAll(com.google.javascript.jscomp.InjectClosureUnawareRuntimeLibraries.ClosureUnawareCallSite callSite, java.util.Set<com.google.javascript.jscomp.js.RuntimeJsLibManager.ExternedField> toInject)
func (this_ *InjectClosureUnawareRuntimeLibraries) priv_InjectClosureUnawareRuntimeLibraries_injectAll_3203(callSite *InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite, toInject *j2gort.OrderedSet[RuntimeJsLibManager_ExternedField]) {
	_ = callSite
	_ = toInject
	if j2gort.CheckNotNullPtr(toInject).IsEmpty() {
		return
	}
	var parameterParent pkg_rhino.Node = NodeUtil_GetFunctionParameters(callSite.ClosureUnawareFunction())
	_ = parameterParent
	var argumentParent pkg_rhino.Node = callSite.CallNode()
	_ = argumentParent
	if parameterParent.HasChildren() && parameterParent.GetLastChild().IsRest() {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "rest parameters not allowed in closureUnaware function parameters, found "+j2gort.StrVal(parameterParent.GetLastChild()))))
	}
	{
		cur_4 := toInject.Iter()
		for cur_4.HasNext() {

			field := cur_4.
				Next()
			_ = field
			var parameter pkg_rhino.Node = pkg_rhino.IR_Name(field.QualifiedName()).Srcref(parameterParent)
			_ = parameter
			var argument pkg_rhino.Node = NodeUtil_NewQName_2_AbstractCompiler_string(this_.mainCompiler, field.UncompiledName()).SrcrefTree(argumentParent)
			_ = argument
			parameterParent.AddChildToBack(parameter)
			this_.shadowCompiler.ReportChangeToEnclosingScope(parameterParent)
			argumentParent.AddChildToBack(argument)
			this_.mainCompiler.ReportChangeToEnclosingScope(argumentParent)
		}
	}
}

// Java static method: InjectClosureUnawareRuntimeLibraries.gatherRuntimeFields(com.google.javascript.rhino.Node root, com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.js.RuntimeJsLibManager.ExternedField> allKnownFields)
func ordSet_InjectClosureUnawareRuntimeLibraries_GatherRuntimeFields(root pkg_rhino.Node, allKnownFields j2gort.Ext_com_google_common_collect_ImmutableMap) j2gort.OrderedSet[RuntimeJsLibManager_ExternedField] {
	_ = root
	_ = allKnownFields
	var seen j2gort.OrderedSet[RuntimeJsLibManager_ExternedField] = j2gort.NewOrderedSet[RuntimeJsLibManager_ExternedField]()
	_ = seen
	var visitor NodeUtil_Visitor = NewLambda_NodeUtil_Visitor(func(n pkg_rhino.Node) {
		_ = n
		if !(n.IsName()) {
			return
		}
		var field RuntimeJsLibManager_ExternedField = j2gort.AsIface[RuntimeJsLibManager_ExternedField](j2gort.MapGetStr(allKnownFields, n.GetString()))
		_ = field
		if field != nil && !field.IsJavaNil_() {
			seen.Add(field)
		}
	})
	_ = visitor
	NodeUtil_VisitPreOrder_2(root, visitor)
	return seen
}

// Java record: com.google.javascript.jscomp.InjectClosureUnawareRuntimeLibraries.ClosureUnawareCallSite
type InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite struct {
	callNode               pkg_rhino.Node `gson:"callNode"`
	closureUnawareFunction pkg_rhino.Node `gson:"closureUnawareFunction"`
}

func (this_ *InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite) BindSelf_(s any) {
	_ = s
}
func (this_ *InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite) IsJavaNil_() bool {
	return this_ == nil
}
func (*InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InjectClosureUnawareRuntimeLibraries$ClosureUnawareCallSite", "ClosureUnawareCallSite"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InjectClosureUnawareRuntimeLibraries.ClosureUnawareCallSite", "com.google.javascript.jscomp.InjectClosureUnawareRuntimeLibraries$ClosureUnawareCallSite", func(v_ any) bool {
		return j2gort.InstanceOf[*InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite](v_)
	})
}
func (*InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite) Implements_j2gort_Ext_java_lang_Record() {
}

func (this_ *InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite) initFields_InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite() {
	_ = this_
}

// Java constructor: ClosureUnawareCallSite(com.google.javascript.rhino.Node callNode, com.google.javascript.rhino.Node closureUnawareFunction)
func NewInjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite(callNode pkg_rhino.Node, closureUnawareFunction pkg_rhino.Node) *InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite {
	this_ := &InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite{}
	InitInjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite(this_, callNode, closureUnawareFunction)
	return this_
}

func InitInjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite(this_ *InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite, callNode pkg_rhino.Node, closureUnawareFunction pkg_rhino.Node) {
	_ = callNode
	_ = closureUnawareFunction
	this_.initFields_InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite()
	if !(callNode.IsCall()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(callNode))
	}
	if !(closureUnawareFunction.IsFunction()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(closureUnawareFunction))
	}
	this_.callNode = callNode
	this_.closureUnawareFunction = closureUnawareFunction
}

func (this_ *InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite) CallNode() pkg_rhino.Node {
	return this_.callNode
}
func (this_ *InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite) ClosureUnawareFunction() pkg_rhino.Node {
	return this_.closureUnawareFunction
}
func (this_ *InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite) Equals(other_ any) bool {
	o_, ok_ := other_.(*InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.callNode != o_.callNode {
		return false
	}
	if this_.closureUnawareFunction != o_.closureUnawareFunction {
		return false
	}
	return true
}
func (this_ *InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite) HashCode() int32 {
	var h_ int32 = 0
	h_ = 31*h_ + int32(this_.callNode)
	h_ = 31*h_ + int32(this_.closureUnawareFunction)
	return h_
}
func (this_ *InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite) String() string {
	_ = fmt.Sprint
	return "ClosureUnawareCallSite[" + "callNode=" + j2gort.StrVal(j2gort.BoxNode(this_.callNode)) + ", " + "closureUnawareFunction=" + j2gort.StrVal(j2gort.BoxNode(this_.closureUnawareFunction)) + "]"
}
func InjectClosureUnawareRuntimeLibraries_GatherRuntimeFields(root pkg_rhino.Node, allKnownFields j2gort.Ext_com_google_common_collect_ImmutableMap) j2gort.Ext_java_util_Set {
	return ordSet_InjectClosureUnawareRuntimeLibraries_GatherRuntimeFields(root, allKnownFields).ToJavaSet()
}
