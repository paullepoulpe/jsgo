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
// Port of com/google/javascript/jscomp/TranspileAndOptimizeClosureUnaware.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

var lambdaSingleton_transpile_and_optimize_closure_unaware_1 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_transpile_and_optimize_closure_unaware_2 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_transpile_and_optimize_closure_unaware_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return pkg_rhino.Node(rawNode_).IsFunction() })
	lambdaSingleton_transpile_and_optimize_closure_unaware_2 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_14 uint32) bool {
		childToVisit := pkg_rhino.Node(rawNode_14)
		_ = childToVisit
		return true
	})
}

// Java class: com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware
type TranspileAndOptimizeClosureUnaware struct {
	original AbstractCompiler `gson:"original"`
	idHash_  atomic.Int32
	mode     NestedCompilerRunner_Mode `gson:"mode"`
}

func (this_ *TranspileAndOptimizeClosureUnaware) BindSelf_(s any) {
	_ = s
}
func (this_ *TranspileAndOptimizeClosureUnaware) IsJavaNil_() bool { return this_ == nil }
func (*TranspileAndOptimizeClosureUnaware) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware", "TranspileAndOptimizeClosureUnaware"
}
func (this_ *TranspileAndOptimizeClosureUnaware) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TranspileAndOptimizeClosureUnaware) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware", "com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware", func(v_ any) bool {
		return j2gort.InstanceOf[*TranspileAndOptimizeClosureUnaware](v_)
	})
}
func (*TranspileAndOptimizeClosureUnaware) Implements_jscomp_CompilerPass() {}

func (this_ *TranspileAndOptimizeClosureUnaware) initFields_TranspileAndOptimizeClosureUnaware() {
	_ = this_
}

// Java constructor: TranspileAndOptimizeClosureUnaware(com.google.javascript.jscomp.AbstractCompiler original, com.google.javascript.jscomp.NestedCompilerRunner.Mode mode)
func NewTranspileAndOptimizeClosureUnaware(original AbstractCompiler, mode NestedCompilerRunner_Mode) *TranspileAndOptimizeClosureUnaware {
	this_ := &TranspileAndOptimizeClosureUnaware{}
	InitTranspileAndOptimizeClosureUnaware(this_, original, mode)
	return this_
}

func InitTranspileAndOptimizeClosureUnaware(this_ *TranspileAndOptimizeClosureUnaware, original AbstractCompiler, mode NestedCompilerRunner_Mode) {
	_ = original
	_ = mode
	this_.initFields_TranspileAndOptimizeClosureUnaware()
	this_.original = original
	this_.mode = mode
}

// Java method: TranspileAndOptimizeClosureUnaware.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *TranspileAndOptimizeClosureUnaware) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var val_1 TranspileAndOptimizeClosureUnaware_CollectShadowAsts
	InitTranspileAndOptimizeClosureUnaware_CollectShadowAsts(&val_1)
	var collector *TranspileAndOptimizeClosureUnaware_CollectShadowAsts = &val_1
	_ = collector
	NodeTraversal_Traverse(this_.original, root, collector)
	if collector.shadowAsts.IsEmpty() {
		return
	}
	this_.priv_TranspileAndOptimizeClosureUnaware_setScriptFeaturesets_5377(j2gort.AsIface[j2gort.Ext_java_util_Collection](collector.shadowAsts.Values()))
	var shadowOptions *CompilerOptions = ClosureUnawareOptions_Convert(this_.original.GetOptions())
	_ = shadowOptions
	var shadowPassConfig *TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig = NewTranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig(shadowOptions, this_.priv_TranspileAndOptimizeClosureUnaware_createInjectRuntimeLibrariesPass_5378(j2gort.AsIface[j2gort.Ext_java_util_Collection](collector.shadowAsts.Values())))
	_ = shadowPassConfig
	var shadowCompiler *NestedCompilerRunner = NestedCompilerRunner_Create(this_.original, shadowOptions, this_.mode, shadowPassConfig)
	_ = shadowCompiler
	this_.priv_TranspileAndOptimizeClosureUnaware_initShadowInputs_5375(shadowCompiler, collector.shadowAsts)
	shadowCompiler.Compile()
	this_.priv_TranspileAndOptimizeClosureUnaware_reattachShadowNodes_5376(collector.shadowAsts)
}

// Java method: TranspileAndOptimizeClosureUnaware.initShadowInputs(com.google.javascript.jscomp.NestedCompilerRunner shadowCompiler, com.google.common.collect.LinkedHashMultimap<com.google.javascript.jscomp.SourceFile,com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware.ShadowAst> perFileInputs)
func (this_ *TranspileAndOptimizeClosureUnaware) priv_TranspileAndOptimizeClosureUnaware_initShadowInputs_5375(shadowCompiler *NestedCompilerRunner, perFileInputs *j2gort.Ext_com_google_common_collect_LinkedHashMultimap) {
	_ = shadowCompiler
	_ = perFileInputs
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](perFileInputs.KeySet())))
		for cur_2.HasNext() {
			rawElem_3 := cur_2.Next()
			sourceFile := j2gort.CastPtr[SourceFile](rawElem_3)
			_ = sourceFile
			var indexInFile int32 = int32(0)
			_ = indexInFile
			{
				cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](perFileInputs.Get(j2gort.PtrToAny(sourceFile)))))
				for cur_4.HasNext() {
					rawElem_5 := cur_4.Next()
					shadowAst := j2gort.CastPtr[TranspileAndOptimizeClosureUnaware_ShadowAst](rawElem_5)
					_ = shadowAst
					var uniqueName string = (j2gort.StrValString(sourceFile.GetName()) + ".shadow" + j2gort.StrValInt32(func() int32 { old_6 := indexInFile; nw_7 := old_6 + 1; indexInFile = nw_7; return old_6 }()))
					_ = uniqueName
					shadowCompiler.AddScript(shadowAst.Script().Detach(), uniqueName)
				}
			}
		}
	}
}

// Java method: TranspileAndOptimizeClosureUnaware.reattachShadowNodes(com.google.common.collect.LinkedHashMultimap<com.google.javascript.jscomp.SourceFile,com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware.ShadowAst> inputs)
func (this_ *TranspileAndOptimizeClosureUnaware) priv_TranspileAndOptimizeClosureUnaware_reattachShadowNodes_5376(inputs *j2gort.Ext_com_google_common_collect_LinkedHashMultimap) {
	_ = inputs
	{
		cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Collection](inputs.Values())))
		for cur_8.HasNext() {
			rawElem_9 := cur_8.Next()
			shadowAst := j2gort.CastPtr[TranspileAndOptimizeClosureUnaware_ShadowAst](rawElem_9)
			_ = shadowAst
			var script pkg_rhino.Node = shadowAst.Script()
			_ = script
			if !(script.HasOneChild()) {
				j2gort.CheckStateMsg(false, j2gort.BoxString(script.ToStringTree()))
			}
			shadowAst.OriginalRoot().AddChildToFront(script.Detach())
			script.SetInputId(nil)
			script.SetStaticSourceFile(nil)
		}
	}
}

// Java method: TranspileAndOptimizeClosureUnaware.setScriptFeaturesets(java.util.Collection<com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware.ShadowAst> inputs)
func (this_ *TranspileAndOptimizeClosureUnaware) priv_TranspileAndOptimizeClosureUnaware_setScriptFeaturesets_5377(inputs j2gort.Ext_java_util_Collection) {
	_ = inputs
	{
		cur_10 := j2gort.NewIterCursor(j2gort.CheckNotNull(inputs))
		for cur_10.HasNext() {
			rawElem_11 := cur_10.Next()
			shadowAst := j2gort.CastPtr[TranspileAndOptimizeClosureUnaware_ShadowAst](rawElem_11)
			_ = shadowAst
			var shadowScript pkg_rhino.Node = shadowAst.Script()
			_ = shadowScript
			var collector *pkg_parsing.FeatureCollector = pkg_parsing.NewFeatureCollector()
			_ = collector
			collector.VisitScript(shadowScript)
			shadowScript.SetFeatureSet(j2gort.PtrToAny(collector.AllFeatures()))
		}
	}
}

// Java method: TranspileAndOptimizeClosureUnaware.createInjectRuntimeLibrariesPass(java.util.Collection<com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware.ShadowAst> shadowAsts)
func (this_ *TranspileAndOptimizeClosureUnaware) priv_TranspileAndOptimizeClosureUnaware_createInjectRuntimeLibrariesPass_5378(shadowAsts j2gort.Ext_java_util_Collection) PassFactory {
	_ = shadowAsts
	var callSites *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	_ = callSites
	{
		cur_12 := j2gort.NewIterCursor(j2gort.CheckNotNull(shadowAsts))
		for cur_12.HasNext() {
			rawElem_13 := cur_12.Next()
			shadowAst := j2gort.CastPtr[TranspileAndOptimizeClosureUnaware_ShadowAst](rawElem_13)
			_ = shadowAst
			var closureUnawareFunction pkg_rhino.Node = NodeUtil_FindPreorder(shadowAst.Script(), lambdaSingleton_transpile_and_optimize_closure_unaware_1, lambdaSingleton_transpile_and_optimize_closure_unaware_2)
			_ = closureUnawareFunction
			callSites.AddVoid_(NewInjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite(shadowAst.CallNode(), closureUnawareFunction))
		}
	}
	return PassFactory_S_Builder().SetName("injectClosureUnawareRuntimeLibs").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_15 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_15)
		_ = compiler
		return NewInjectClosureUnawareRuntimeLibraries(this_.original, compiler, j2gort.AsIface[j2gort.Ext_java_util_List](callSites.BuildImmList_()))
	})).Build()
}

// Java record: com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware.ShadowAst
type TranspileAndOptimizeClosureUnaware_ShadowAst struct {
	originalRoot pkg_rhino.Node `gson:"originalRoot"`
	callNode     pkg_rhino.Node `gson:"callNode"`
	script       pkg_rhino.Node `gson:"script"`
}

func (this_ *TranspileAndOptimizeClosureUnaware_ShadowAst) BindSelf_(s any) {
	_ = s
}
func (this_ *TranspileAndOptimizeClosureUnaware_ShadowAst) IsJavaNil_() bool { return this_ == nil }
func (*TranspileAndOptimizeClosureUnaware_ShadowAst) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware$ShadowAst", "ShadowAst"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware.ShadowAst", "com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware$ShadowAst", func(v_ any) bool {
		return j2gort.InstanceOf[*TranspileAndOptimizeClosureUnaware_ShadowAst](v_)
	})
}
func (*TranspileAndOptimizeClosureUnaware_ShadowAst) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *TranspileAndOptimizeClosureUnaware_ShadowAst) initFields_TranspileAndOptimizeClosureUnaware_ShadowAst() {
	_ = this_
}

// Java constructor: ShadowAst(com.google.javascript.rhino.Node originalRoot, com.google.javascript.rhino.Node callNode, com.google.javascript.rhino.Node script)
func NewTranspileAndOptimizeClosureUnaware_ShadowAst(originalRoot pkg_rhino.Node, callNode pkg_rhino.Node, script pkg_rhino.Node) *TranspileAndOptimizeClosureUnaware_ShadowAst {
	this_ := &TranspileAndOptimizeClosureUnaware_ShadowAst{}
	InitTranspileAndOptimizeClosureUnaware_ShadowAst(this_, originalRoot, callNode, script)
	return this_
}

func InitTranspileAndOptimizeClosureUnaware_ShadowAst(this_ *TranspileAndOptimizeClosureUnaware_ShadowAst, originalRoot pkg_rhino.Node, callNode pkg_rhino.Node, script pkg_rhino.Node) {
	_ = originalRoot
	_ = callNode
	_ = script
	this_.initFields_TranspileAndOptimizeClosureUnaware_ShadowAst()
	this_.originalRoot = originalRoot
	this_.callNode = callNode
	this_.script = script
}

func (this_ *TranspileAndOptimizeClosureUnaware_ShadowAst) OriginalRoot() pkg_rhino.Node {
	return this_.originalRoot
}
func (this_ *TranspileAndOptimizeClosureUnaware_ShadowAst) CallNode() pkg_rhino.Node {
	return this_.callNode
}
func (this_ *TranspileAndOptimizeClosureUnaware_ShadowAst) Script() pkg_rhino.Node {
	return this_.script
}
func (this_ *TranspileAndOptimizeClosureUnaware_ShadowAst) Equals(other_ any) bool {
	o_, ok_ := other_.(*TranspileAndOptimizeClosureUnaware_ShadowAst)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.originalRoot != o_.originalRoot {
		return false
	}
	if this_.callNode != o_.callNode {
		return false
	}
	if this_.script != o_.script {
		return false
	}
	return true
}
func (this_ *TranspileAndOptimizeClosureUnaware_ShadowAst) HashCode() int32 {
	var h_ int32 = 0
	h_ = 31*h_ + int32(this_.originalRoot)
	h_ = 31*h_ + int32(this_.callNode)
	h_ = 31*h_ + int32(this_.script)
	return h_
}
func (this_ *TranspileAndOptimizeClosureUnaware_ShadowAst) String() string {
	_ = fmt.Sprint
	return "ShadowAst[" + "originalRoot=" + j2gort.StrVal(j2gort.BoxNode(this_.originalRoot)) + ", " + "callNode=" + j2gort.StrVal(j2gort.BoxNode(this_.callNode)) + ", " + "script=" + j2gort.StrVal(j2gort.BoxNode(this_.script)) + "]"
}

// Java class: com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware.CollectShadowAsts
type TranspileAndOptimizeClosureUnaware_CollectShadowAsts struct {
	shadowAsts *j2gort.Ext_com_google_common_collect_LinkedHashMultimap `gson:"shadowAsts"`
	idHash_    atomic.Int32
}

func (this_ *TranspileAndOptimizeClosureUnaware_CollectShadowAsts) BindSelf_(s any) {
	_ = s
}
func (this_ *TranspileAndOptimizeClosureUnaware_CollectShadowAsts) IsJavaNil_() bool {
	return this_ == nil
}
func (*TranspileAndOptimizeClosureUnaware_CollectShadowAsts) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware$CollectShadowAsts", "CollectShadowAsts"
}
func (this_ *TranspileAndOptimizeClosureUnaware_CollectShadowAsts) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TranspileAndOptimizeClosureUnaware_CollectShadowAsts) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware.CollectShadowAsts", "com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware$CollectShadowAsts", func(v_ any) bool {
		return j2gort.InstanceOf[*TranspileAndOptimizeClosureUnaware_CollectShadowAsts](v_)
	})
}
func (*TranspileAndOptimizeClosureUnaware_CollectShadowAsts) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *TranspileAndOptimizeClosureUnaware_CollectShadowAsts) initFields_TranspileAndOptimizeClosureUnaware_CollectShadowAsts() {
	_ = this_
	this_.shadowAsts = j2gort.Ext_com_google_common_collect_LinkedHashMultimap_Create_0()
}

// Java constructor: CollectShadowAsts()
func NewTranspileAndOptimizeClosureUnaware_CollectShadowAsts() *TranspileAndOptimizeClosureUnaware_CollectShadowAsts {
	this_ := &TranspileAndOptimizeClosureUnaware_CollectShadowAsts{}
	InitTranspileAndOptimizeClosureUnaware_CollectShadowAsts(this_)
	return this_
}

func InitTranspileAndOptimizeClosureUnaware_CollectShadowAsts(this_ *TranspileAndOptimizeClosureUnaware_CollectShadowAsts) {
	this_.initFields_TranspileAndOptimizeClosureUnaware_CollectShadowAsts()
}

// Java method: CollectShadowAsts.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TranspileAndOptimizeClosureUnaware_CollectShadowAsts) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if !(n.IsScript()) {
		return true
	}
	if !(n.IsClosureUnawareCode()) {
		return false
	}
	var fileJsDoc *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = fileJsDoc
	if (fileJsDoc == nil) || (fileJsDoc.GetPerFileClosureUnawareMode() == 0) {
		return true
	}
	var minificationLevel pkg_rhino.JSDocInfo_PerFileClosureUnawareMode = fileJsDoc.GetPerFileClosureUnawareMode()
	_ = minificationLevel
	switch (minificationLevel).Ordinal() {
	case int32(0) /* PerFileClosureUnawareMode.UNSPECIFIED */, int32(1) /* PerFileClosureUnawareMode.SIMPLE */ :
		return true
	case int32(2) /* PerFileClosureUnawareMode.WHITESPACE */ :
		return false
	}
	panic("unreachable switch expression")
}

// Java method: CollectShadowAsts.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TranspileAndOptimizeClosureUnaware_CollectShadowAsts) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var shadow pkg_rhino.Node = n.GetClosureUnawareShadow()
	_ = shadow
	if shadow == 0 {
		return
	}
	var callNode pkg_rhino.Node
	_ = callNode
	if n.GetParent().IsCall() {
		callNode = n.GetParent()
	} else {
		callNode = n.GetGrandparent()
	}
	if !(callNode.IsCall()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(callNode))
	}
	this_.shadowAsts.Put(j2gort.PtrToAny(t.GetInput().GetSourceFile()), NewTranspileAndOptimizeClosureUnaware_ShadowAst(shadow, callNode, shadow.GetOnlyChild()))
}

// Java class: com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware.RuntimeLibInjectionPassConfig
type TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig struct {
	PassConfig_PassConfigDelegateBase
	injectRuntimeLibs PassFactory `gson:"injectRuntimeLibs"`
}

func (this_ *TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig) BindSelf_(s any) {
	_ = s
	this_.PassConfig_PassConfigDelegateBase.BindSelf_(s)
}
func (this_ *TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig) IsJavaNil_() bool {
	return this_ == nil
}
func (*TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware$RuntimeLibInjectionPassConfig", "RuntimeLibInjectionPassConfig"
}
func (this_ *TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PassConfig_PassConfigDelegateBase.IdentityHash_()
}
func (*TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware.RuntimeLibInjectionPassConfig", "com.google.javascript.jscomp.TranspileAndOptimizeClosureUnaware$RuntimeLibInjectionPassConfig", func(v_ any) bool {
		return j2gort.InstanceOf[*TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig](v_)
	})
}
func (*TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig) Implements_jscomp_PassConfig_PassConfigDelegate() {
}
func (*TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig) Implements_jscomp_PassConfig() {
}

func (this_ *TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig) initFields_TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig() {
	_ = this_
}

// Java constructor: RuntimeLibInjectionPassConfig(com.google.javascript.jscomp.CompilerOptions options, com.google.javascript.jscomp.PassFactory injectRuntimeLibs)
func NewTranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig(options *CompilerOptions, injectRuntimeLibs PassFactory) *TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig {
	this_ := &TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig{}
	this_.BindSelf_(this_)
	InitTranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig(this_, options, injectRuntimeLibs)
	return this_
}

func InitTranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig(this_ *TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig, options *CompilerOptions, injectRuntimeLibs PassFactory) {
	_ = options
	_ = injectRuntimeLibs
	InitPassConfig_PassConfigDelegate(&this_.PassConfig_PassConfigDelegateBase, NewDefaultPassConfig(options))
	this_.initFields_TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig()
	this_.injectRuntimeLibs = injectRuntimeLibs
}

// Java method: RuntimeLibInjectionPassConfig.getTranspileOnlyPasses()
func (this_ *TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig) GetTranspileOnlyPasses() *PassListBuilder {
	var passes *PassListBuilder = this_.PassConfig_PassConfigDelegateBase.GetTranspileOnlyPasses()
	_ = passes
	passes.AddAfter(this_.injectRuntimeLibs, "markUnnormalized")
	return passes
}

// Java method: RuntimeLibInjectionPassConfig.getFinalizations()
func (this_ *TranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig) GetFinalizations() *PassListBuilder {
	var passes *PassListBuilder = this_.PassConfig_PassConfigDelegateBase.GetFinalizations()
	_ = passes
	passes.AddAfter(this_.injectRuntimeLibs, "markUnnormalized")
	return passes
}
func (*TranspileAndOptimizeClosureUnaware_CollectShadowAsts) CanReleaseSyntacticScopes_() {}
