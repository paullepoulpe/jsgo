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
// Port of com/google/javascript/jscomp/ConvertChunksToESModules.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_convert_chunks_to_es_modules_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_convert_chunks_to_es_modules_2 j2gort.Ext_java_util_function_Function
var lambdaSingleton_convert_chunks_to_es_modules_3 j2gort.Ext_java_util_function_Function
var lambdaSingleton_convert_chunks_to_es_modules_4 j2gort.Ext_java_util_function_Function
var lambdaSingleton_convert_chunks_to_es_modules_5 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_convert_chunks_to_es_modules_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_20 any) any {
		k := j2gort.CastPtr[JSChunk](lp_20)
		_ = k
		return j2gort.NewExt_java_util_LinkedHashSet()
	})
	lambdaSingleton_convert_chunks_to_es_modules_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_25 any) any {
		k := j2gort.CastPtr[JSChunk](lp_25)
		_ = k
		return j2gort.NewExt_java_util_LinkedHashMap()
	})
	lambdaSingleton_convert_chunks_to_es_modules_3 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_30 any) any {
		k := j2gort.CastPtr[JSChunk](lp_30)
		_ = k
		return j2gort.NewExt_java_util_LinkedHashSet()
	})
	lambdaSingleton_convert_chunks_to_es_modules_4 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_35 any) any {
		k := j2gort.CastPtr[JSChunk](lp_35)
		_ = k
		return j2gort.NewExt_java_util_LinkedHashMap()
	})
	lambdaSingleton_convert_chunks_to_es_modules_5 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_42 any) any {
		k := j2gort.CastPtr[JSChunk](lp_42)
		_ = k
		return j2gort.NewExt_java_util_LinkedHashSet()
	})
}

// Java class: com.google.javascript.jscomp.ConvertChunksToESModules
type ConvertChunksToESModules struct {
	compiler               AbstractCompiler                                      `gson:"compiler"`
	crossChunkExports      j2gort.OrderedMap[*JSChunk, j2gort.Ext_java_util_Set] `gson:"crossChunkExports"`
	crossChunkImports      j2gort.OrderedMap[*JSChunk, j2gort.Ext_java_util_Map] `gson:"crossChunkImports"`
	dynamicImportCallbacks []pkg_rhino.Node                                      `gson:"dynamicImportCallbacks"`
	idHash_                atomic.Int32
}

func (this_ *ConvertChunksToESModules) BindSelf_(s any) {
	_ = s
}
func (this_ *ConvertChunksToESModules) IsJavaNil_() bool { return this_ == nil }
func (*ConvertChunksToESModules) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConvertChunksToESModules", "ConvertChunksToESModules"
}
func (this_ *ConvertChunksToESModules) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConvertChunksToESModules) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConvertChunksToESModules", "com.google.javascript.jscomp.ConvertChunksToESModules", func(v_ any) bool {
		return j2gort.InstanceOf[*ConvertChunksToESModules](v_)
	})
}
func (*ConvertChunksToESModules) Implements_jscomp_CompilerPass() {}

func StaticGet_ConvertChunksToESModules_DYNAMIC_IMPORT_CALLBACK_FN() string {
	return "jscomp$DynamicImportCallback"
}

var static_ConvertChunksToESModules_ASSIGNMENT_TO_IMPORT *DiagnosticType

func StaticGet_ConvertChunksToESModules_ASSIGNMENT_TO_IMPORT() *DiagnosticType {
	clinit_ConvertChunksToESModules()
	return static_ConvertChunksToESModules_ASSIGNMENT_TO_IMPORT
}
func StaticSet_ConvertChunksToESModules_ASSIGNMENT_TO_IMPORT(v *DiagnosticType) {
	clinit_ConvertChunksToESModules()
	static_ConvertChunksToESModules_ASSIGNMENT_TO_IMPORT = v
}

var static_ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH *DiagnosticType

func StaticGet_ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH() *DiagnosticType {
	clinit_ConvertChunksToESModules()
	return static_ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH
}
func StaticSet_ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH(v *DiagnosticType) {
	clinit_ConvertChunksToESModules()
	static_ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH = v
}

var static_ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK *DiagnosticType

func StaticGet_ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK() *DiagnosticType {
	clinit_ConvertChunksToESModules()
	return static_ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK
}
func StaticSet_ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK(v *DiagnosticType) {
	clinit_ConvertChunksToESModules()
	static_ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK = v
}

var clinitDone_ConvertChunksToESModules uint32
var clinitState_ConvertChunksToESModules j2gort.ClinitState

func clinit_ConvertChunksToESModules() {
	if atomic.LoadUint32(&clinitDone_ConvertChunksToESModules) != 0 {
		return
	}
	clinitSlow_ConvertChunksToESModules()
}

func clinitSlow_ConvertChunksToESModules() {
	if !j2gort.EnterClinit(&clinitState_ConvertChunksToESModules, &clinitDone_ConvertChunksToESModules) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ConvertChunksToESModules, &clinitDone_ConvertChunksToESModules)
	static_ConvertChunksToESModules_ASSIGNMENT_TO_IMPORT = DiagnosticType_Error("JSC_IMPORT_ASSIGN", "Imported symbol \"{0}\" in chunk \"{1}\" cannot be assigned (defined in \"{2}\")")
	static_ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH = DiagnosticType_Error("JSC_UNABLE_TO_COMPUTE_RELATIVE_PATH", "Unable to compute relative import path from \"{0}\" to \"{1}\"")
	static_ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK = DiagnosticType_Error("JSC_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK", "Dynamic import callback encountered wih an invalid format.{0}")
}

func (this_ *ConvertChunksToESModules) initFields_ConvertChunksToESModules() {
	_ = this_
	this_.crossChunkExports = j2gort.NewOrderedMap[*JSChunk, j2gort.Ext_java_util_Set]()
	this_.crossChunkImports = j2gort.NewOrderedMap[*JSChunk, j2gort.Ext_java_util_Map]()
	this_.dynamicImportCallbacks = []pkg_rhino.Node{}
}

// Java constructor: ConvertChunksToESModules(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewConvertChunksToESModules(compiler AbstractCompiler) *ConvertChunksToESModules {
	clinit_ConvertChunksToESModules()
	this_ := &ConvertChunksToESModules{}
	InitConvertChunksToESModules(this_, compiler)
	return this_
}

func InitConvertChunksToESModules(this_ *ConvertChunksToESModules, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_ConvertChunksToESModules()
	this_.compiler = compiler
}

// Java method: ConvertChunksToESModules.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ConvertChunksToESModules) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, NewConvertChunksToESModules_FindCrossChunkReferences(this_))
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.compiler.GetChunkGraph().GetAllChunks()))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			chunk := j2gort.CastPtr[JSChunk](rawElem_2)
			_ = chunk
			if (!(this_.crossChunkExports.ContainsKey(chunk)) && !(this_.crossChunkImports.ContainsKey(chunk))) && !(chunk.GetInputs().IsEmpty()) {
				this_.crossChunkExports.Put(chunk, j2gort.NewExt_java_util_LinkedHashSet())
			}
		}
	}
	this_.priv_ConvertChunksToESModules_convertChunkSourcesToModules_2226()
	this_.priv_ConvertChunksToESModules_addExportStatements_2227()
	this_.priv_ConvertChunksToESModules_addImportStatements_2228()
	this_.priv_ConvertChunksToESModules_rewriteDynamicImportCallbacks_2229()
	NodeUtil_AddFeatureToAllScripts(root, pkg_parser.StaticGet_FeatureSet_Feature_MODULES(), this_.compiler)
}

// Java method: ConvertChunksToESModules.convertChunkSourcesToModules()
func (this_ *ConvertChunksToESModules) priv_ConvertChunksToESModules_convertChunkSourcesToModules_2226() {
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.compiler.GetChunkGraph().GetAllChunks()))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			chunk := j2gort.CastPtr[JSChunk](rawElem_4)
			_ = chunk
			if chunk.GetInputs().IsEmpty() {
				continue
			}
			var firstInput *CompilerInput = nil
			_ = firstInput
			{
				cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunk.GetInputs()))
				for cur_5.HasNext() {
					rawElem_6 := cur_5.Next()
					input := j2gort.CastPtr[CompilerInput](rawElem_6)
					_ = input
					var astRoot pkg_rhino.Node = input.GetAstRoot(this_.compiler)
					_ = astRoot
					var scriptFeatures *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(astRoot)
					_ = scriptFeatures
					if !(!(scriptFeatures.Contains_1_FeatureSet(pkg_parser.StaticGet_FeatureSet_ES2015_MODULES()))) {
						j2gort.CheckState(false)
					}
					if firstInput == nil {
						firstInput = input
						scriptFeatures = scriptFeatures.Union(pkg_parser.StaticGet_FeatureSet_ES2015_MODULES())
						astRoot.SetFeatureSet(j2gort.PtrToAny(scriptFeatures))
						var moduleBody pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_MODULE_BODY())
						_ = moduleBody
						moduleBody.Srcref(astRoot)
						moduleBody.AddChildrenToFront(astRoot.RemoveChildren())
						astRoot.AddChildToFront(moduleBody)
						this_.compiler.ReportChangeToEnclosingScope(moduleBody)
					} else {
						var firstInputAstRoot pkg_rhino.Node = firstInput.GetAstRoot(this_.compiler)
						_ = firstInputAstRoot
						var firstInputScriptFeatures *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(firstInputAstRoot)
						_ = firstInputScriptFeatures
						var combinedFeatureSet *pkg_parser.FeatureSet = firstInputScriptFeatures.Union(NodeUtil_GetFeatureSetOfScript(astRoot))
						_ = combinedFeatureSet
						astRoot.SetFeatureSet(j2gort.PtrToAny(combinedFeatureSet))
						var moduleBody_2 pkg_rhino.Node = firstInputAstRoot.GetFirstChild()
						_ = moduleBody_2
						if !((moduleBody_2 != 0) && moduleBody_2.IsModuleBody()) {
							j2gort.CheckState(false)
						}
						moduleBody_2.AddChildrenToBack(astRoot.RemoveChildren())
						this_.compiler.ReportChangeToEnclosingScope(firstInputAstRoot)
						this_.compiler.ReportChangeToChangeScope(astRoot)
					}
				}
			}
		}
	}
}

// Java method: ConvertChunksToESModules.addExportStatements()
func (this_ *ConvertChunksToESModules) priv_ConvertChunksToESModules_addExportStatements_2227() {
	{
		cur_7 := this_.crossChunkExports.IterEntries()
		for cur_7.NextEntry() {
			var firstInput *CompilerInput = cur_7.Key().GetFirst()
			_ = firstInput
			var moduleBody pkg_rhino.Node = firstInput.GetAstRoot(this_.compiler).GetFirstChild()
			_ = moduleBody
			if !((moduleBody != 0) && moduleBody.IsModuleBody()) {
				j2gort.CheckState(false)
			}
			var exportSpecs pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_EXPORT_SPECS())
			_ = exportSpecs
			{
				cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(cur_7.Value()))
				for cur_8.HasNext() {
					name := cur_8.NextStr()
					_ = name
					var exportSpec pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_EXPORT_SPEC())
					_ = exportSpec
					exportSpec.AddChildToFront(pkg_rhino.IR_Name(name))
					exportSpec.AddChildToFront(pkg_rhino.IR_Name(name))
					exportSpec.PutIntProp(pkg_rhino.StaticGet_Node_IS_SHORTHAND_PROPERTY(), int32(1))
					exportSpecs.AddChildToBack(exportSpec)
				}
			}
			var importsByChunk j2gort.Ext_java_util_Map = this_.crossChunkImports.Get(cur_7.Key())
			_ = importsByChunk
			if (exportSpecs.HasChildren() || (importsByChunk == nil)) || importsByChunk.IsEmpty() {
				var export pkg_rhino.Node = pkg_rhino.IR_Export(exportSpecs).SrcrefTree(moduleBody)
				_ = export
				moduleBody.AddChildToBack(export)
				this_.compiler.ReportChangeToEnclosingScope(moduleBody)
			}
		}
	}
}

// Java static method: ConvertChunksToESModules.getChunkName(com.google.javascript.jscomp.JSChunk chunk)
func ConvertChunksToESModules_GetChunkName(chunk *JSChunk) string {
	clinit_ConvertChunksToESModules()
	_ = chunk
	return (j2gort.StrValString(chunk.GetName()) + ".js")
}

// Java method: ConvertChunksToESModules.addImportStatements()
func (this_ *ConvertChunksToESModules) priv_ConvertChunksToESModules_addImportStatements_2228() {
	{
		cur_9 := this_.crossChunkImports.IterEntries()
		for cur_9.NextEntry() {
			var importStatements []pkg_rhino.Node = []pkg_rhino.Node{}
			_ = importStatements
			var importingChunk *JSChunk = cur_9.Key()
			_ = importingChunk
			var firstInput *CompilerInput = importingChunk.GetFirst()
			_ = firstInput
			var moduleBody pkg_rhino.Node = firstInput.GetAstRoot(this_.compiler).GetFirstChild()
			_ = moduleBody
			if !((moduleBody != 0) && moduleBody.IsModuleBody()) {
				j2gort.CheckState(false)
			}
			{
				cur_10 := j2gort.NewMapEntryCursor(cur_9.Value())
				for cur_10.NextEntry() {
					var importSpecs pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_IMPORT_SPECS())
					_ = importSpecs
					{
						cur_11 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](cur_10.GetValue())))
						for cur_11.HasNext() {
							name := cur_11.NextStr()
							_ = name
							var importSpec pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_IMPORT_SPEC())
							_ = importSpec
							importSpec.AddChildToFront(pkg_rhino.IR_Name(name))
							importSpec.AddChildToFront(pkg_rhino.IR_Name(name))
							importSpec.PutIntProp(pkg_rhino.StaticGet_Node_IS_SHORTHAND_PROPERTY(), int32(1))
							importSpecs.AddChildToBack(importSpec)
						}
					}
					var importStatement pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_IMPORT())
					_ = importStatement
					var exportingChunk *JSChunk = j2gort.CastPtr[JSChunk](cur_10.GetKey())
					_ = exportingChunk
					var importPath string = ConvertChunksToESModules_GetChunkName(exportingChunk)
					_ = importPath
					func() {
						defer func() {
							if rec_12 := recover(); rec_12 != nil {
								jex_13 := javalang.RecoverException(rec_12)
								_ = jex_13
								if j2gort.InstanceOf[javalang.IIllegalArgumentException](jex_13) {
									e := j2gort.AsIface[javalang.IIllegalArgumentException](jex_13)
									_ = e
									this_.compiler.Report_1(JSError_Make_3(moduleBody, StaticGet_ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH(), []string{ConvertChunksToESModules_GetChunkName(importingChunk), ConvertChunksToESModules_GetChunkName(exportingChunk)}))
								} else {
									panic(rec_12)
								}
							}
						}()
						importPath = ModuleLoader_RelativePathFrom(ConvertChunksToESModules_GetChunkName(importingChunk), ConvertChunksToESModules_GetChunkName(exportingChunk))
					}()
					importStatement.AddChildToFront(pkg_rhino.IR_String(importPath))
					if importSpecs.HasChildren() {
						importStatement.AddChildToFront(importSpecs)
					} else {
						importStatement.AddChildToFront(pkg_rhino.IR_Empty())
					}
					importStatement.AddChildToFront(pkg_rhino.IR_Empty())
					importStatement.SrcrefTree(moduleBody)
					importStatements = j2gort.SliceInsertAt(importStatements, int32(0), importStatement)
				}
			}
			for _, importStatement_2 := range importStatements {
				_ = importStatement_2
				moduleBody.AddChildToFront(importStatement_2)
			}

			this_.compiler.ReportChangeToEnclosingScope(moduleBody)
		}
	}
}

// Java static method: ConvertChunksToESModules.getDynamicImportCallbackModuleNamespace(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node call)
func ConvertChunksToESModules_GetDynamicImportCallbackModuleNamespace(compiler AbstractCompiler, call pkg_rhino.Node) pkg_rhino.Node {
	clinit_ConvertChunksToESModules()
	_ = compiler
	_ = call
	if !(call.IsCall()) {
		j2gort.CheckState(false)
	}
	var callbackFn pkg_rhino.Node = NodeUtil_GetArgumentForCallOrNew(call, int32(0))
	_ = callbackFn
	if ((callbackFn == 0) || !(callbackFn.IsFunction())) || NodeUtil_GetFunctionParameters(callbackFn).HasChildren() {
		compiler.Report_1(JSError_Make_3(call, StaticGet_ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK(), []string{" Unable to find valid callback function."}))
		return pkg_rhino.Node(0)
	}
	var callbackBody pkg_rhino.Node = NodeUtil_GetFunctionBody(callbackFn)
	_ = callbackBody
	if callbackBody.IsName() {
		return callbackBody
	} else if (((callbackBody.IsBlock() && callbackBody.HasOneChild()) && callbackBody.GetFirstChild().IsReturn()) && callbackBody.GetFirstChild().HasOneChild()) && callbackBody.GetFirstFirstChild().IsName() {
		return callbackBody.GetFirstFirstChild()
	}
	compiler.Report_1(JSError_Make_3(call, StaticGet_ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK(), []string{" Unable to find valid namespace reference."}))
	return pkg_rhino.Node(0)
}

// Java method: ConvertChunksToESModules.rewriteDynamicImportCallbacks()
func (this_ *ConvertChunksToESModules) priv_ConvertChunksToESModules_rewriteDynamicImportCallbacks_2229() {
	var astFactory *AstFactory = this_.compiler.CreateAstFactory()
	_ = astFactory
	for _, dynamicImportCallback := range this_.dynamicImportCallbacks {
		_ = dynamicImportCallback
		if !(dynamicImportCallback.IsCall()) {
			j2gort.CheckState(false)
		}
		var moduleNamespace pkg_rhino.Node = ConvertChunksToESModules_GetDynamicImportCallbackModuleNamespace(this_.compiler, dynamicImportCallback)
		_ = moduleNamespace
		if moduleNamespace == 0 {
			continue
		}
		var callbackFn pkg_rhino.Node = NodeUtil_GetArgumentForCallOrNew(dynamicImportCallback, int32(0))
		_ = callbackFn
		var callbackParamList pkg_rhino.Node = NodeUtil_GetFunctionParameters(callbackFn)
		_ = callbackParamList
		var importNamespaceParam pkg_rhino.Node = astFactory.CreateNameWithUnknownType("$").Srcref(moduleNamespace)
		_ = importNamespaceParam
		callbackParamList.AddChildToFront(importNamespaceParam)
		this_.compiler.ReportChangeToEnclosingScope(importNamespaceParam)
		var namespaceGetprop pkg_rhino.Node = astFactory.CreateGetPropWithUnknownType(astFactory.CreateNameWithUnknownType("$").Srcref(moduleNamespace), moduleNamespace.GetString())
		_ = namespaceGetprop
		moduleNamespace.ReplaceWith(namespaceGetprop)
		this_.compiler.ReportChangeToEnclosingScope(namespaceGetprop)
		var innerCallback pkg_rhino.Node = NodeUtil_GetArgumentForCallOrNew(dynamicImportCallback, int32(0)).Detach()
		_ = innerCallback
		dynamicImportCallback.ReplaceWith(innerCallback)
		this_.compiler.ReportChangeToEnclosingScope(innerCallback)
	}

}

// Java static method: ConvertChunksToESModules.isDynamicImportCallback(com.google.javascript.rhino.Node call)
func ConvertChunksToESModules_IsDynamicImportCallback(call pkg_rhino.Node) bool {
	clinit_ConvertChunksToESModules()
	_ = call
	if !(call.IsCall()) {
		return false
	}
	if NodeUtil_IsCallTo_2_Node_string(call, "jscomp$DynamicImportCallback") {
		return true
	}
	return false
}

// Java method: ConvertChunksToESModules.visitCallAndTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *ConvertChunksToESModules) priv_ConvertChunksToESModules_visitCallAndTraverse_2230(t *NodeTraversal, call pkg_rhino.Node) bool {
	_ = t
	_ = call
	if !(call.IsCall()) {
		j2gort.CheckState(false)
	}
	if !(ConvertChunksToESModules_IsDynamicImportCallback(call)) {
		return true
	}
	var moduleNamespace pkg_rhino.Node = ConvertChunksToESModules_GetDynamicImportCallbackModuleNamespace(this_.compiler, call)
	_ = moduleNamespace
	if moduleNamespace == 0 {
		return true
	}
	var isValidModuleNamespace bool = this_.VisitName(t, moduleNamespace, StaticGet_ConvertChunksToESModules_ImportType_DYNAMIC())
	_ = isValidModuleNamespace
	if isValidModuleNamespace {
		this_.dynamicImportCallbacks = append(this_.dynamicImportCallbacks, call)
	} else {
		this_.compiler.Report_1(JSError_Make_3(call, StaticGet_ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK(), []string{" Unable to find valid namespace reference."}))
	}
	return false
}

// Java method: ConvertChunksToESModules.visitName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node nameNode, com.google.javascript.jscomp.ConvertChunksToESModules.ImportType importType)
func (this_ *ConvertChunksToESModules) VisitName(t *NodeTraversal, nameNode pkg_rhino.Node, importType ConvertChunksToESModules_ImportType) bool {
	_ = t
	_ = nameNode
	_ = importType
	if !(nameNode.IsName()) {
		j2gort.CheckState(false)
	}
	var name string = nameNode.GetString()
	_ = name
	if j2gort.CheckNotNullString(name) == "" {
		return false
	}
	var s *Scope = t.GetScope()
	_ = s
	var v *Var = j2gort.CastPtr[Var](s.GetVar(name))
	_ = v
	if (v == nil) || !(v.IsGlobal()) {
		return false
	}
	var input *CompilerInput = v.GetInput()
	_ = input
	if input == nil {
		return false
	}
	var definingChunk *JSChunk = input.GetChunk()
	_ = definingChunk
	var referencingChunk *JSChunk = t.GetChunk()
	_ = referencingChunk
	if definingChunk != referencingChunk {
		if NodeUtil_IsLhsOfAssign(nameNode) {
			t.Report_3(nameNode, StaticGet_ConvertChunksToESModules_ASSIGNMENT_TO_IMPORT(), []string{nameNode.GetString(), ConvertChunksToESModules_GetChunkName(referencingChunk), ConvertChunksToESModules_GetChunkName(definingChunk)})
		}
		var namesToExport j2gort.Ext_java_util_Set = j2gort.OrderedMapComputeIfAbsentFn(&this_.crossChunkExports, definingChunk, nil, lambdaSingleton_convert_chunks_to_es_modules_1)

		_ = namesToExport
		j2gort.SetAddStr(namesToExport, name)
		var namesToImportByChunk j2gort.Ext_java_util_Map = j2gort.OrderedMapComputeIfAbsentFn(&this_.crossChunkImports, referencingChunk, nil, lambdaSingleton_convert_chunks_to_es_modules_2)

		_ = namesToImportByChunk
		if importType == StaticGet_ConvertChunksToESModules_ImportType_STATIC() {
			var importsForChunk j2gort.Ext_java_util_Set
		iifeSwitch_1:
			switch {
			default:
				m_26 := namesToImportByChunk
				k_27 := j2gort.PtrToAny(definingChunk)
				if v_28, ok_29 := j2gort.MapTryGetComputeAny(m_26, k_27); ok_29 {
					importsForChunk = j2gort.AsIface[j2gort.Ext_java_util_Set](v_28)
					break iifeSwitch_1
				}
				importsForChunk = j2gort.AsIface[j2gort.Ext_java_util_Set](m_26.ComputeIfAbsent(k_27, lambdaSingleton_convert_chunks_to_es_modules_3))
				break iifeSwitch_1
			}

			_ = importsForChunk
			j2gort.SetAddStr(importsForChunk, name)
		}
	}
	return true
}

// Java enum: com.google.javascript.jscomp.ConvertChunksToESModules.ImportType
type ConvertChunksToESModules_ImportType uint8

const (
	enum_ConvertChunksToESModules_ImportType_STATIC  ConvertChunksToESModules_ImportType = 1
	enum_ConvertChunksToESModules_ImportType_DYNAMIC ConvertChunksToESModules_ImportType = 2
)

var enumBaseTable_ConvertChunksToESModules_ImportType = [...]j2gort.EnumBase{
	{Name_: "STATIC", Ordinal_: 0},
	{Name_: "DYNAMIC", Ordinal_: 1},
}
var enumValues_ConvertChunksToESModules_ImportType = []ConvertChunksToESModules_ImportType{1, 2}

func StaticGet_ConvertChunksToESModules_ImportType_STATIC() ConvertChunksToESModules_ImportType {
	return enum_ConvertChunksToESModules_ImportType_STATIC
}
func StaticGet_ConvertChunksToESModules_ImportType_DYNAMIC() ConvertChunksToESModules_ImportType {
	return enum_ConvertChunksToESModules_ImportType_DYNAMIC
}

func ConvertChunksToESModules_ImportType_Values() []ConvertChunksToESModules_ImportType {
	return j2gort.CloneSlice(enumValues_ConvertChunksToESModules_ImportType)
}

func ConvertChunksToESModules_ImportType_ValueOf(name string) ConvertChunksToESModules_ImportType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ConvertChunksToESModules_ImportType {
		if enumBaseTable_ConvertChunksToESModules_ImportType[i_].Name_ == name {
			return ConvertChunksToESModules_ImportType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ConvertChunksToESModules.ImportType." + name))
}

func (this_ ConvertChunksToESModules_ImportType) IsJavaNil_() bool                              { return this_ == 0 }
func (ConvertChunksToESModules_ImportType) HasIdentityEquals_()                                 {}
func (ConvertChunksToESModules_ImportType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ConvertChunksToESModules_ImportType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ConvertChunksToESModules_ImportType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ConvertChunksToESModules_ImportType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConvertChunksToESModules$ImportType", "ImportType"
}
func (this_ ConvertChunksToESModules_ImportType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ConvertChunksToESModules_ImportType[this_-1]
}
func (this_ ConvertChunksToESModules_ImportType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ConvertChunksToESModules_ImportType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ConvertChunksToESModules_ImportType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ConvertChunksToESModules_ImportType[this_-1].Name_
}
func (this_ ConvertChunksToESModules_ImportType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ConvertChunksToESModules_ImportType[this_-1].Name_
}
func (this_ ConvertChunksToESModules_ImportType) Equals(other_ any) bool {
	o_, ok_ := other_.(ConvertChunksToESModules_ImportType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ConvertChunksToESModules_ImportType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ConvertChunksToESModules_ImportType[this_-1].HashCode()
}
func (this_ ConvertChunksToESModules_ImportType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ConvertChunksToESModules_ImportType[this_-1].CompareTo(other_)
}
func (this_ ConvertChunksToESModules_ImportType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ConvertChunksToESModules_ImportType[this_-1].Name_)
}
func init() {
	vals_ := []any{ConvertChunksToESModules_ImportType(1), ConvertChunksToESModules_ImportType(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ConvertChunksToESModules.ImportType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ConvertChunksToESModules$ImportType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConvertChunksToESModules.ImportType", "com.google.javascript.jscomp.ConvertChunksToESModules$ImportType", func(v_ any) bool {
		return j2gort.InstanceOf[ConvertChunksToESModules_ImportType](v_)
	})
}

// Java class: com.google.javascript.jscomp.ConvertChunksToESModules.FindCrossChunkReferences
type ConvertChunksToESModules_FindCrossChunkReferences struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer_ *ConvertChunksToESModules
}

func (this_ *ConvertChunksToESModules_FindCrossChunkReferences) BindSelf_(s any) {
	_ = s
}
func (this_ *ConvertChunksToESModules_FindCrossChunkReferences) IsJavaNil_() bool {
	return this_ == nil
}
func (*ConvertChunksToESModules_FindCrossChunkReferences) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConvertChunksToESModules$FindCrossChunkReferences", "FindCrossChunkReferences"
}
func (this_ *ConvertChunksToESModules_FindCrossChunkReferences) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*ConvertChunksToESModules_FindCrossChunkReferences) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConvertChunksToESModules.FindCrossChunkReferences", "com.google.javascript.jscomp.ConvertChunksToESModules$FindCrossChunkReferences", func(v_ any) bool {
		return j2gort.InstanceOf[*ConvertChunksToESModules_FindCrossChunkReferences](v_)
	})
}
func (this_ *ConvertChunksToESModules_FindCrossChunkReferences) SetOuter_(o *ConvertChunksToESModules) {
	this_.outer_ = o
}
func (this_ *ConvertChunksToESModules_FindCrossChunkReferences) GetOuter_() *ConvertChunksToESModules {
	return this_.outer_
}
func (*ConvertChunksToESModules_FindCrossChunkReferences) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {
}
func (*ConvertChunksToESModules_FindCrossChunkReferences) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *ConvertChunksToESModules_FindCrossChunkReferences) initFields_ConvertChunksToESModules_FindCrossChunkReferences() {
	_ = this_
}

// Java constructor: FindCrossChunkReferences()
func NewConvertChunksToESModules_FindCrossChunkReferences(outer_ *ConvertChunksToESModules) *ConvertChunksToESModules_FindCrossChunkReferences {
	this_ := &ConvertChunksToESModules_FindCrossChunkReferences{}
	this_.outer_ = outer_
	InitConvertChunksToESModules_FindCrossChunkReferences(this_)
	return this_
}

func InitConvertChunksToESModules_FindCrossChunkReferences(this_ *ConvertChunksToESModules_FindCrossChunkReferences) {
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_ConvertChunksToESModules_FindCrossChunkReferences()
}

// Java method: FindCrossChunkReferences.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ConvertChunksToESModules_FindCrossChunkReferences) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		this_.VisitScript(t, n)
		return true
	} else if n.IsCall() {
		return this_.outer_.priv_ConvertChunksToESModules_visitCallAndTraverse_2230(t, n)
	} else if n.IsName() {
		this_.outer_.VisitName(t, n, StaticGet_ConvertChunksToESModules_ImportType_STATIC())
		return true
	}
	return true
}

// Java method: FindCrossChunkReferences.visitScript(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node script)
func (this_ *ConvertChunksToESModules_FindCrossChunkReferences) VisitScript(t *NodeTraversal, script pkg_rhino.Node) {
	_ = t
	_ = script
	if !(script.IsScript()) {
		j2gort.CheckState(false)
	}
	var chunk *JSChunk = t.GetChunk()
	_ = chunk
	var chunkDependencies []*JSChunk = j2gort.ImmListToSliceNullable[*JSChunk](chunk.GetDependencies())
	_ = chunkDependencies
	this_.outer_.crossChunkExports.PutIfAbsent(chunk, j2gort.NewExt_java_util_LinkedHashSet())
	if len(chunkDependencies) != 0 {
		var namesToImportByChunk j2gort.Ext_java_util_Map = j2gort.OrderedMapComputeIfAbsentFn(&this_.outer_.crossChunkImports, chunk, nil, lambdaSingleton_convert_chunks_to_es_modules_4)

		_ = namesToImportByChunk
		{

			for _, dependency := range chunkDependencies {

				_ = dependency
				func() {
					m_38 := namesToImportByChunk
					k_39 := j2gort.PtrToAny(dependency)
					if _, ok_41 := j2gort.MapTryGetComputeAny(m_38, k_39); !ok_41 {
						m_38.ComputeIfAbsent(k_39, lambdaSingleton_convert_chunks_to_es_modules_5)
					}
				}()
			}
		}
	}
}
func (*ConvertChunksToESModules_FindCrossChunkReferences) NoOpVisit_()                 {}
func (*ConvertChunksToESModules_FindCrossChunkReferences) CanReleaseSyntacticScopes_() {}
