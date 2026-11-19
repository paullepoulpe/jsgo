/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/PolymerBehaviorExtractor.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.PolymerBehaviorExtractor
type PolymerBehaviorExtractor struct {
	compiler              AbstractCompiler                                                               `gson:"compiler"`
	globalNames           *GlobalNamespace                                                               `gson:"globalNames"`
	moduleMetadataMap     *ModuleMetadataMap                                                             `gson:"moduleMetadataMap"`
	moduleMap             *ModuleMap                                                                     `gson:"moduleMap"`
	resolveMemoized       j2gort.Ext_com_google_common_collect_Table                                     `gson:"resolveMemoized"`
	globalResolveMemoized j2gort.OrderedMap[string, *PolymerBehaviorExtractor_ResolveBehaviorNameResult] `gson:"globalResolveMemoized"`
	idHash_               atomic.Int32
}

func (this_ *PolymerBehaviorExtractor) BindSelf_(s any) {
	_ = s
}
func (this_ *PolymerBehaviorExtractor) IsJavaNil_() bool { return this_ == nil }
func (*PolymerBehaviorExtractor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolymerBehaviorExtractor", "PolymerBehaviorExtractor"
}
func (this_ *PolymerBehaviorExtractor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PolymerBehaviorExtractor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolymerBehaviorExtractor", "com.google.javascript.jscomp.PolymerBehaviorExtractor", func(v_ any) bool {
		return j2gort.InstanceOf[*PolymerBehaviorExtractor](v_)
	})
}

var static_PolymerBehaviorExtractor_BEHAVIOR_NAMES_NOT_TO_COPY j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_PolymerBehaviorExtractor_BEHAVIOR_NAMES_NOT_TO_COPY() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_PolymerBehaviorExtractor()
	return static_PolymerBehaviorExtractor_BEHAVIOR_NAMES_NOT_TO_COPY
}
func StaticSet_PolymerBehaviorExtractor_BEHAVIOR_NAMES_NOT_TO_COPY(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_PolymerBehaviorExtractor()
	static_PolymerBehaviorExtractor_BEHAVIOR_NAMES_NOT_TO_COPY = v
}
func StaticGet_PolymerBehaviorExtractor_GOOG_MODULE_EXPORTS() string { return "exports" }

var static_PolymerBehaviorExtractor_FAILED_RESOLVE_RESULT *PolymerBehaviorExtractor_ResolveBehaviorNameResult

func StaticGet_PolymerBehaviorExtractor_FAILED_RESOLVE_RESULT() *PolymerBehaviorExtractor_ResolveBehaviorNameResult {
	clinit_PolymerBehaviorExtractor()
	return static_PolymerBehaviorExtractor_FAILED_RESOLVE_RESULT
}
func StaticSet_PolymerBehaviorExtractor_FAILED_RESOLVE_RESULT(v *PolymerBehaviorExtractor_ResolveBehaviorNameResult) {
	clinit_PolymerBehaviorExtractor()
	static_PolymerBehaviorExtractor_FAILED_RESOLVE_RESULT = v
}

var static_PolymerBehaviorExtractor_GOOG_MODULE_GET pkg_rhino.QualifiedName

func StaticGet_PolymerBehaviorExtractor_GOOG_MODULE_GET() pkg_rhino.QualifiedName {
	clinit_PolymerBehaviorExtractor()
	return static_PolymerBehaviorExtractor_GOOG_MODULE_GET
}
func StaticSet_PolymerBehaviorExtractor_GOOG_MODULE_GET(v pkg_rhino.QualifiedName) {
	clinit_PolymerBehaviorExtractor()
	static_PolymerBehaviorExtractor_GOOG_MODULE_GET = v
}

var clinitDone_PolymerBehaviorExtractor uint32
var clinitState_PolymerBehaviorExtractor j2gort.ClinitState

func clinit_PolymerBehaviorExtractor() {
	if atomic.LoadUint32(&clinitDone_PolymerBehaviorExtractor) != 0 {
		return
	}
	clinitSlow_PolymerBehaviorExtractor()
}

func clinitSlow_PolymerBehaviorExtractor() {
	if !j2gort.EnterClinit(&clinitState_PolymerBehaviorExtractor, &clinitDone_PolymerBehaviorExtractor) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PolymerBehaviorExtractor, &clinitDone_PolymerBehaviorExtractor)
	static_PolymerBehaviorExtractor_BEHAVIOR_NAMES_NOT_TO_COPY = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("created", "attached", "detached", "attributeChanged", "configure", "ready", []any{"properties", "listeners", "observers", "hostAttributes"})
	static_PolymerBehaviorExtractor_FAILED_RESOLVE_RESULT = NewPolymerBehaviorExtractor_ResolveBehaviorNameResult(pkg_rhino.Node(0), false, nil)
	static_PolymerBehaviorExtractor_GOOG_MODULE_GET = pkg_rhino.QualifiedName_Of("goog.module.get")
}

func (this_ *PolymerBehaviorExtractor) initFields_PolymerBehaviorExtractor() {
	_ = this_
	this_.resolveMemoized = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Table](j2gort.Ext_com_google_common_collect_HashBasedTable_Create_0())
	this_.globalResolveMemoized = j2gort.NewOrderedMap[string, *PolymerBehaviorExtractor_ResolveBehaviorNameResult](

	// Java constructor: PolymerBehaviorExtractor(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.GlobalNamespace globalNames, com.google.javascript.jscomp.modules.ModuleMetadataMap moduleMetadataMap, com.google.javascript.jscomp.modules.ModuleMap moduleMap)
	)
}

func NewPolymerBehaviorExtractor(compiler AbstractCompiler, globalNames *GlobalNamespace, moduleMetadataMap *ModuleMetadataMap, moduleMap *ModuleMap) *PolymerBehaviorExtractor {
	clinit_PolymerBehaviorExtractor()
	this_ := &PolymerBehaviorExtractor{}
	InitPolymerBehaviorExtractor(this_, compiler, globalNames, moduleMetadataMap, moduleMap)
	return this_
}

func InitPolymerBehaviorExtractor(this_ *PolymerBehaviorExtractor, compiler AbstractCompiler, globalNames *GlobalNamespace, moduleMetadataMap *ModuleMetadataMap, moduleMap *ModuleMap) {
	_ = compiler
	_ = globalNames
	_ = moduleMetadataMap
	_ = moduleMap
	this_.initFields_PolymerBehaviorExtractor()
	this_.compiler = compiler
	this_.globalNames = globalNames
	this_.moduleMetadataMap = moduleMetadataMap
	this_.moduleMap = moduleMap
}

// Java method: PolymerBehaviorExtractor.extractBehaviors(com.google.javascript.rhino.Node behaviorArray, com.google.javascript.jscomp.modules.ModuleMetadataMap.@org.jspecify.annotations.Nullable ModuleMetadata moduleMetadata)
func (this_ *PolymerBehaviorExtractor) ExtractBehaviors(behaviorArray pkg_rhino.Node, moduleMetadata ModuleMetadataMap_ModuleMetadata) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = behaviorArray
	_ = moduleMetadata
	if behaviorArray == 0 {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	if !(behaviorArray.IsArrayLit()) {
		this_.compiler.Report_1(JSError_Make_3(behaviorArray, StaticGet_PolymerPassErrors_POLYMER_INVALID_BEHAVIOR_ARRAY(), []string{}))
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	var behaviors []*PolymerBehaviorExtractor_BehaviorDefinition = []*PolymerBehaviorExtractor_BehaviorDefinition{}
	_ = behaviors
	{
		var behaviorName pkg_rhino.Node = behaviorArray.GetFirstChild()
		_ = behaviorName
		for ; behaviorName != 0; behaviorName = behaviorName.GetNext() {
			if behaviorName.IsObjectLit() {
				PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(behaviorName, this_.compiler)
				PolymerPassStaticUtils_QuoteListenerAndHostAttributeKeys(behaviorName, this_.compiler)
				PolymerPassStaticUtils_ProtectObserverAndPropertyFunctionKeys(behaviorName)
				if NodeUtil_GetFirstPropMatchingKey(behaviorName, "is") != 0 {
					this_.compiler.Report_1(JSError_Make_3(behaviorName, StaticGet_PolymerPassErrors_POLYMER_INVALID_BEHAVIOR(), []string{}))
				}
				var behaviorModule pkg_rhino.Node = NodeUtil_GetEnclosingModuleIfPresent(behaviorName)
				_ = behaviorModule
				behaviors = append(behaviors, NewPolymerBehaviorExtractor_BehaviorDefinition(j2gort.AsIface[j2gort.Ext_java_util_List](PolymerPassStaticUtils_ExtractProperties(behaviorName, StaticGet_PolymerClassDefinition_DefinitionType_ObjectLiteral(), this_.compiler, pkg_rhino.Node(0))), j2gort.AsIface[j2gort.Ext_java_util_List](PolymerBehaviorExtractor_GetBehaviorFunctionsToCopy(behaviorName)), j2gort.AsIface[j2gort.Ext_java_util_List](PolymerBehaviorExtractor_GetNonPropertyMembersToCopy(behaviorName)), NodeUtil_GetEnclosingScopeRoot(behaviorName).IsRoot(), (j2gort.CastPtr[pkg_parser.FeatureSet](NodeUtil_GetEnclosingScript(behaviorName).GetFeatureSet())), behaviorModule))
				continue
			}
			var resolveResult *PolymerBehaviorExtractor_ResolveBehaviorNameResult = nil
			_ = resolveResult
			if PolymerBehaviorExtractor_IsGoogModuleGetCall(behaviorName) {
				resolveResult = this_.priv_PolymerBehaviorExtractor_resolveGoogModuleGet_4366(behaviorName.GetSecondChild().GetString())
			} else {
				resolveResult = this_.priv_PolymerBehaviorExtractor_resolveBehaviorName_4363(PolymerBehaviorExtractor_GetQualifiedNameThroughCast(behaviorName), moduleMetadata)
			}
			{
				var iifeCond_1 bool
				{

					rp_1 := resolveResult
					ap_2 := StaticGet_PolymerBehaviorExtractor_FAILED_RESOLVE_RESULT()
					iifeCond_1 = j2gort.CheckNotNullPtr(rp_1) == ap_2
				}
				if iifeCond_1 {
					this_.compiler.Report_1(JSError_Make_3(behaviorName, StaticGet_PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR(), []string{}))
					continue
				}
			}
			var behaviorValue pkg_rhino.Node = resolveResult.node
			_ = behaviorValue
			if behaviorValue.IsArrayLit() {
				behaviors = j2gort.AppendImmCollection[*PolymerBehaviorExtractor_BehaviorDefinition](behaviors, j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.ExtractBehaviors(behaviorValue, resolveResult.moduleMetadata)))
			} else if behaviorValue.IsObjectLit() {
				PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(behaviorValue, this_.compiler)
				PolymerPassStaticUtils_QuoteListenerAndHostAttributeKeys(behaviorValue, this_.compiler)
				if NodeUtil_GetFirstPropMatchingKey(behaviorValue, "is") != 0 {
					this_.compiler.Report_1(JSError_Make_3(behaviorValue, StaticGet_PolymerPassErrors_POLYMER_INVALID_BEHAVIOR(), []string{}))
				}
				var behaviorModule_2 pkg_rhino.Node = NodeUtil_GetEnclosingModuleIfPresent(behaviorValue)
				_ = behaviorModule_2
				behaviors = append(behaviors, NewPolymerBehaviorExtractor_BehaviorDefinition(j2gort.AsIface[j2gort.Ext_java_util_List](PolymerPassStaticUtils_ExtractProperties(behaviorValue, StaticGet_PolymerClassDefinition_DefinitionType_ObjectLiteral(), this_.compiler, pkg_rhino.Node(0))), j2gort.AsIface[j2gort.Ext_java_util_List](PolymerBehaviorExtractor_GetBehaviorFunctionsToCopy(behaviorValue)), j2gort.AsIface[j2gort.Ext_java_util_List](PolymerBehaviorExtractor_GetNonPropertyMembersToCopy(behaviorValue)), resolveResult.isGlobalDeclaration, (j2gort.CastPtr[pkg_parser.FeatureSet](NodeUtil_GetEnclosingScript(behaviorValue).GetFeatureSet())), behaviorModule_2))
			} else {
				this_.compiler.Report_1(JSError_Make_3(behaviorName, StaticGet_PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR(), []string{}))
			}
		}
	}
	return j2gort.ImmListFromTypedSlice(behaviors)
}

// Java method: PolymerBehaviorExtractor.resolveBehaviorName(java.lang.@org.jspecify.annotations.Nullable String name, com.google.javascript.jscomp.modules.ModuleMetadataMap.@org.jspecify.annotations.Nullable ModuleMetadata moduleMetadata)
func (this_ *PolymerBehaviorExtractor) priv_PolymerBehaviorExtractor_resolveBehaviorName_4363(name string, moduleMetadata ModuleMetadataMap_ModuleMetadata) *PolymerBehaviorExtractor_ResolveBehaviorNameResult {
	_ = name
	_ = moduleMetadata
	if name == j2gort.NullString {
		return StaticGet_PolymerBehaviorExtractor_FAILED_RESOLVE_RESULT()
	}
	var memoized *PolymerBehaviorExtractor_ResolveBehaviorNameResult
	_ = memoized
	if moduleMetadata != nil && !moduleMetadata.IsJavaNil_() {
		memoized = j2gort.CastPtr[PolymerBehaviorExtractor_ResolveBehaviorNameResult](this_.resolveMemoized.Get(j2gort.BoxString(name), moduleMetadata))
	} else {
		memoized = this_.globalResolveMemoized.Get(name)
	}
	if memoized == nil {
		memoized = j2gort.CheckNotNullPtr(this_.priv_PolymerBehaviorExtractor_resolveBehaviorNameInternal_4364(name, moduleMetadata))
		if moduleMetadata != nil && !moduleMetadata.IsJavaNil_() {
			this_.resolveMemoized.Put(j2gort.BoxString(name), moduleMetadata, j2gort.PtrToAny(memoized))
		} else {
			this_.globalResolveMemoized.Put(name, memoized)
		}
	}
	return memoized
}

// Java method: PolymerBehaviorExtractor.resolveBehaviorNameInternal(java.lang.String name, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata moduleMetadata)
func (this_ *PolymerBehaviorExtractor) priv_PolymerBehaviorExtractor_resolveBehaviorNameInternal_4364(name string, moduleMetadata ModuleMetadataMap_ModuleMetadata) *PolymerBehaviorExtractor_ResolveBehaviorNameResult {
	_ = name
	_ = moduleMetadata
	var result *PolymerBehaviorExtractor_ResolveBehaviorNameResult = this_.priv_PolymerBehaviorExtractor_getNameIfModuleImport_4367(name, moduleMetadata)
	_ = result
	if result != nil {
		return result
	}
	var legacyResolve *PolymerBehaviorExtractor_ResolveBehaviorNameResult = this_.priv_PolymerBehaviorExtractor_resolveReferenceToLegacyGoogModule_4365(name)
	_ = legacyResolve
	if legacyResolve != nil {
		return legacyResolve
	}
	var moduleLevelName *GlobalNamespace_Name
	_ = moduleLevelName
	if moduleMetadata != nil && !moduleMetadata.IsJavaNil_() {
		moduleLevelName = this_.globalNames.GetNameFromModule(moduleMetadata, name)
	} else {
		moduleLevelName = nil
	}
	var globalName *GlobalNamespace_Name
	_ = globalName
	if moduleLevelName == nil {
		globalName = j2gort.CastPtr[GlobalNamespace_Name](this_.globalNames.GetSlot(name))
	} else {
		globalName = moduleLevelName
	}
	if globalName == nil {
		return StaticGet_PolymerBehaviorExtractor_FAILED_RESOLVE_RESULT()
	}
	var isGlobalDeclaration bool = (moduleLevelName == nil)
	_ = isGlobalDeclaration
	var declarationRef *GlobalNamespace_Ref = j2gort.CastPtr[GlobalNamespace_Ref](globalName.GetDeclaration())
	_ = declarationRef
	if declarationRef == nil {
		{
			cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(globalName.GetRefs()))
			for cur_3.HasNext() {
				rawElem_4 := cur_3.Next()
				ref := j2gort.CastPtr[GlobalNamespace_Ref](rawElem_4)
				_ = ref
				if ref.IsSet() {
					isGlobalDeclaration = false
					declarationRef = ref
					break
				}
			}
		}
	}
	if declarationRef == nil {
		return StaticGet_PolymerBehaviorExtractor_FAILED_RESOLVE_RESULT()
	}
	var declarationNode pkg_rhino.Node = declarationRef.GetNode()
	_ = declarationNode
	if declarationNode == 0 {
		return StaticGet_PolymerBehaviorExtractor_FAILED_RESOLVE_RESULT()
	}
	var rValue pkg_rhino.Node = NodeUtil_GetRValueOfLValue(declarationNode)
	_ = rValue
	if rValue == 0 {
		return StaticGet_PolymerBehaviorExtractor_FAILED_RESOLVE_RESULT()
	}
	if rValue.IsQualifiedName() {
		var declarationScope *Scope = j2gort.CastPtr[Scope](declarationRef.scope.GetClosestHoistScope())
		_ = declarationScope
		var m *Module = ModuleImportResolver_GetModuleFromScopeRoot(this_.compiler.GetModuleMap(), j2gort.AsIface[CompilerInputProvider](this_.compiler), declarationScope.GetRootNode_0())
		_ = m
		return this_.priv_PolymerBehaviorExtractor_resolveBehaviorName_4363(PolymerBehaviorExtractor_GetQualifiedNameThroughCast(rValue), func() ModuleMetadataMap_ModuleMetadata {
			if m != nil {
				return m.Metadata()
			}
			return nil
		}())
	}
	var behaviorInfo *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(declarationNode)
	_ = behaviorInfo
	if (behaviorInfo == nil) || !(behaviorInfo.IsPolymerBehavior()) {
		this_.compiler.Report_1(JSError_Make_3(declarationNode, StaticGet_PolymerPassErrors_POLYMER_UNANNOTATED_BEHAVIOR(), []string{}))
	}
	return NewPolymerBehaviorExtractor_ResolveBehaviorNameResult(rValue, isGlobalDeclaration, moduleMetadata)
}

// Java method: PolymerBehaviorExtractor.resolveReferenceToLegacyGoogModule(java.lang.String name)
func (this_ *PolymerBehaviorExtractor) priv_PolymerBehaviorExtractor_resolveReferenceToLegacyGoogModule_4365(name string) *PolymerBehaviorExtractor_ResolveBehaviorNameResult {
	_ = name
	var dot int32 = j2gort.JStrLen(name)
	_ = dot
	for dot >= int32(0) {
		var subNamespace string = j2gort.JStrSub(name, int32(0), dot)
		_ = subNamespace
		var metadata ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.moduleMetadataMap.GetModulesByGoogNamespace(), subNamespace))
		_ = metadata
		if (metadata == nil || metadata.IsJavaNil_()) || !(metadata.IsLegacyGoogModule()) {
			dot = j2gort.JStrLastIndexOfCpFrom(name, int32(uint16(46)), (dot - int32(1)))
			continue
		}
		var rest string
		_ = rest
		if dot == j2gort.JStrLen(name) {
			rest = ""
		} else {
			rest = j2gort.JStrSubFrom(name, dot)
		}
		var result *PolymerBehaviorExtractor_ResolveBehaviorNameResult = this_.priv_PolymerBehaviorExtractor_resolveBehaviorName_4363(("exports" + j2gort.StrValString(rest)), metadata)
		_ = result
		{
			var iifeCond_1 bool
			{

				rp_5 := result
				ap_6 := StaticGet_PolymerBehaviorExtractor_FAILED_RESOLVE_RESULT()
				iifeCond_1 = j2gort.CheckNotNullPtr(rp_5) == ap_6
			}
			if iifeCond_1 {
				return nil
			} else {
				return result
			}
		}
	}
	return nil
}

// Java static method: PolymerBehaviorExtractor.isGoogModuleGetCall(com.google.javascript.rhino.Node callNode)
func PolymerBehaviorExtractor_IsGoogModuleGetCall(callNode pkg_rhino.Node) bool {
	clinit_PolymerBehaviorExtractor()
	_ = callNode
	if !(callNode.IsCall()) {
		return false
	}
	return ((StaticGet_PolymerBehaviorExtractor_GOOG_MODULE_GET().Matches(callNode.GetFirstChild()) && callNode.HasTwoChildren()) && callNode.GetSecondChild().IsStringLit())
}

// Java method: PolymerBehaviorExtractor.resolveGoogModuleGet(java.lang.String moduleNamespace)
func (this_ *PolymerBehaviorExtractor) priv_PolymerBehaviorExtractor_resolveGoogModuleGet_4366(moduleNamespace string) *PolymerBehaviorExtractor_ResolveBehaviorNameResult {
	_ = moduleNamespace
	var closureModule ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.moduleMetadataMap.GetModulesByGoogNamespace(), moduleNamespace))
	_ = closureModule
	if closureModule == nil || closureModule.IsJavaNil_() {
		return StaticGet_PolymerBehaviorExtractor_FAILED_RESOLVE_RESULT()
	} else if closureModule.IsGoogProvide() {
		return this_.priv_PolymerBehaviorExtractor_resolveBehaviorName_4363(moduleNamespace, nil)
	}
	if !(closureModule.IsGoogModule()) {
		j2gort.CheckStateMsg(false, closureModule)
	}
	return this_.priv_PolymerBehaviorExtractor_resolveBehaviorName_4363("exports", closureModule)
}

// Java method: PolymerBehaviorExtractor.getNameIfModuleImport(java.lang.String name, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata)
func (this_ *PolymerBehaviorExtractor) priv_PolymerBehaviorExtractor_getNameIfModuleImport_4367(name string, metadata ModuleMetadataMap_ModuleMetadata) *PolymerBehaviorExtractor_ResolveBehaviorNameResult {
	_ = name
	_ = metadata
	if (metadata == nil || metadata.IsJavaNil_()) || (!(metadata.IsEs6Module()) && !(metadata.IsGoogModule())) {
		return nil
	}
	var module *Module
	_ = module
	if metadata.IsGoogModule() {
		module = this_.moduleMap.GetClosureModule(j2gort.ListGetStr(metadata.GoogNamespaces().AsList(), int32(0)))
	} else {
		module = this_.moduleMap.GetModule_1_ModuleLoader_ModulePath(metadata.Path())
	}
	j2gort.CheckNotNullPtr(module)
	var dot int32 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
	_ = dot
	var root string
	_ = root
	if dot == int32(-1) {
		root = name
	} else {
		root = j2gort.JStrSub(name, int32(0), dot)
	}
	var b *Binding = j2gort.CastPtr[Binding](j2gort.MapGetStr(module.BoundNames(), root))
	_ = b
	if (b == nil) || !(b.IsSomeImport()) {
		return nil
	}
	var rest string
	_ = rest
	if dot == int32(-1) {
		rest = ""
	} else {
		rest = j2gort.JStrSubFrom(name, dot)
	}
	if b.IsModuleNamespace() {
		return this_.priv_PolymerBehaviorExtractor_resolveModuleNamespaceBinding_4368(b, rest)
	}
	var importMetadata ModuleMetadataMap_ModuleMetadata = b.OriginatingExport().ModuleMetadata()
	_ = importMetadata
	var originatingName string = j2gort.NullString
	_ = originatingName
	if importMetadata.IsEs6Module() {
		originatingName = (j2gort.StrValString(b.OriginatingExport().LocalName()) + j2gort.StrValString(rest))
	} else if importMetadata.IsGoogModule() {
		originatingName = ("exports" + "." + j2gort.StrValString(b.OriginatingExport().ExportName()) + j2gort.StrValString(rest))
	} else {
		if !(importMetadata.IsGoogProvide()) {
			j2gort.CheckStateMsg(false, importMetadata)
		}
		originatingName = (j2gort.StrValString(b.ClosureNamespace()) + "." + j2gort.StrValString(b.OriginatingExport().ExportName()) + j2gort.StrValString(rest))
	}
	return this_.priv_PolymerBehaviorExtractor_resolveBehaviorName_4363(originatingName, importMetadata)
}

// Java method: PolymerBehaviorExtractor.resolveModuleNamespaceBinding(com.google.javascript.jscomp.modules.Binding b, java.lang.String rest)
func (this_ *PolymerBehaviorExtractor) priv_PolymerBehaviorExtractor_resolveModuleNamespaceBinding_4368(b *Binding, rest string) *PolymerBehaviorExtractor_ResolveBehaviorNameResult {
	_ = b
	_ = rest
	if b.Metadata().IsGoogModule() {
		return this_.priv_PolymerBehaviorExtractor_resolveBehaviorName_4363(("exports" + j2gort.StrValString(rest)), b.Metadata())
	} else if b.Metadata().IsGoogProvide() {
		return this_.priv_PolymerBehaviorExtractor_resolveBehaviorName_4363((j2gort.StrValString(b.ClosureNamespace()) + j2gort.StrValString(rest)), b.Metadata())
	}
	if !(b.Metadata().IsEs6Module()) {
		j2gort.CheckState(false)
	}
	if j2gort.CheckNotNullString(rest) == "" {
		return StaticGet_PolymerBehaviorExtractor_FAILED_RESOLVE_RESULT()
	}
	rest = j2gort.JStrSubFrom(rest, int32(1))
	var dot int32 = j2gort.JStrIndexOfCp(rest, int32(uint16(46)))
	_ = dot
	var exportedName string
	_ = exportedName
	if dot == int32(-1) {
		exportedName = rest
	} else {
		exportedName = j2gort.JStrSub(rest, int32(0), dot)
	}
	var originalModule *Module = this_.moduleMap.GetModule_1_ModuleLoader_ModulePath(b.Metadata().Path())
	_ = originalModule
	var exportBinding *Binding = j2gort.CastPtr[Binding](j2gort.MapGetStr(originalModule.Namespace(), exportedName))
	_ = exportBinding
	if (exportBinding == nil) || !(exportBinding.IsCreatedByEsExport()) {
		return StaticGet_PolymerBehaviorExtractor_FAILED_RESOLVE_RESULT()
	}
	return this_.priv_PolymerBehaviorExtractor_resolveBehaviorName_4363((j2gort.StrValString(exportBinding.OriginatingExport().LocalName()) + j2gort.StrValString(func() string {
		if dot == int32(-1) {
			return ""
		}
		return j2gort.JStrSubFrom(rest, dot)
	}())), b.Metadata())
}

// Java static method: PolymerBehaviorExtractor.getBehaviorFunctionsToCopy(com.google.javascript.rhino.Node behaviorObjLit)
func PolymerBehaviorExtractor_GetBehaviorFunctionsToCopy(behaviorObjLit pkg_rhino.Node) j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_PolymerBehaviorExtractor()
	_ = behaviorObjLit
	if !(behaviorObjLit.IsObjectLit()) {
		j2gort.CheckState(false)
	}
	var functionsToCopy []*PolymerPass_MemberDefinition = []*PolymerPass_MemberDefinition{}
	_ = functionsToCopy
	var enclosingModule pkg_rhino.Node = NodeUtil_GetEnclosingModuleIfPresent(behaviorObjLit)
	_ = enclosingModule
	{
		var keyNode pkg_rhino.Node = behaviorObjLit.GetFirstChild()
		_ = keyNode
		for ; keyNode != 0; keyNode = keyNode.GetNext() {
			var isFunctionDefinition bool = ((keyNode.IsStringKey() && keyNode.GetFirstChild().IsFunction()) || keyNode.IsMemberFunctionDef())
			_ = isFunctionDefinition
			if isFunctionDefinition && !(j2gort.SetContainsStr(StaticGet_PolymerBehaviorExtractor_BEHAVIOR_NAMES_NOT_TO_COPY(), keyNode.GetString())) {
				functionsToCopy = append(functionsToCopy, NewPolymerPass_MemberDefinition(NodeUtil_GetBestJSDocInfo(keyNode), keyNode, keyNode.GetFirstChild(), enclosingModule))
			}
		}
	}
	return j2gort.ImmListFromTypedSlice(functionsToCopy)
}

// Java static method: PolymerBehaviorExtractor.getQualifiedNameThroughCast(com.google.javascript.rhino.Node node)
func PolymerBehaviorExtractor_GetQualifiedNameThroughCast(node pkg_rhino.Node) string {
	clinit_PolymerBehaviorExtractor()
	_ = node
	if node.IsName() {
		var name string = node.GetString()
		_ = name
		if j2gort.CheckNotNullString(name) == "" {
			return j2gort.NullString
		} else {
			return name
		}
	} else if node.IsGetProp() {
		var left string = PolymerBehaviorExtractor_GetQualifiedNameThroughCast(node.GetFirstChild())
		_ = left
		if left == j2gort.NullString {
			return j2gort.NullString
		}
		var right string = node.GetString()
		_ = right
		return (j2gort.StrValString(left) + "." + j2gort.StrValString(right))
	} else if node.IsCast() {
		return PolymerBehaviorExtractor_GetQualifiedNameThroughCast(node.GetFirstChild())
	}
	return j2gort.NullString
}

// Java static method: PolymerBehaviorExtractor.getNonPropertyMembersToCopy(com.google.javascript.rhino.Node behaviorObjLit)
func PolymerBehaviorExtractor_GetNonPropertyMembersToCopy(behaviorObjLit pkg_rhino.Node) j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_PolymerBehaviorExtractor()
	_ = behaviorObjLit
	if !(behaviorObjLit.IsObjectLit()) {
		j2gort.CheckState(false)
	}
	var membersToCopy []*PolymerPass_MemberDefinition = []*PolymerPass_MemberDefinition{}
	_ = membersToCopy
	var enclosingModule pkg_rhino.Node = NodeUtil_GetEnclosingModuleIfPresent(behaviorObjLit)
	_ = enclosingModule
	{
		var keyNode pkg_rhino.Node = behaviorObjLit.GetFirstChild()
		_ = keyNode
		for ; keyNode != 0; keyNode = keyNode.GetNext() {
			var isNonFunctionMember bool = (keyNode.IsGetterDef() || (keyNode.IsStringKey() && !(keyNode.GetFirstChild().IsFunction())))
			_ = isNonFunctionMember
			if isNonFunctionMember && !(j2gort.SetContainsStr(StaticGet_PolymerBehaviorExtractor_BEHAVIOR_NAMES_NOT_TO_COPY(), keyNode.GetString())) {
				membersToCopy = append(membersToCopy, NewPolymerPass_MemberDefinition(NodeUtil_GetBestJSDocInfo(keyNode), keyNode, keyNode.GetFirstChild(), enclosingModule))
			}
		}
	}
	return j2gort.ImmListFromTypedSlice(membersToCopy)
}

// Java class: com.google.javascript.jscomp.PolymerBehaviorExtractor.ResolveBehaviorNameResult
type PolymerBehaviorExtractor_ResolveBehaviorNameResult struct {
	moduleMetadata      ModuleMetadataMap_ModuleMetadata `gson:"moduleMetadata"`
	node                pkg_rhino.Node                   `gson:"node"`
	idHash_             atomic.Int32
	isGlobalDeclaration bool `gson:"isGlobalDeclaration"`
}

func (this_ *PolymerBehaviorExtractor_ResolveBehaviorNameResult) BindSelf_(s any) {
	_ = s
}
func (this_ *PolymerBehaviorExtractor_ResolveBehaviorNameResult) IsJavaNil_() bool {
	return this_ == nil
}
func (*PolymerBehaviorExtractor_ResolveBehaviorNameResult) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolymerBehaviorExtractor$ResolveBehaviorNameResult", "ResolveBehaviorNameResult"
}
func (this_ *PolymerBehaviorExtractor_ResolveBehaviorNameResult) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PolymerBehaviorExtractor_ResolveBehaviorNameResult) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolymerBehaviorExtractor.ResolveBehaviorNameResult", "com.google.javascript.jscomp.PolymerBehaviorExtractor$ResolveBehaviorNameResult", func(v_ any) bool {
		return j2gort.InstanceOf[*PolymerBehaviorExtractor_ResolveBehaviorNameResult](v_)
	})
}

func (this_ *PolymerBehaviorExtractor_ResolveBehaviorNameResult) initFields_PolymerBehaviorExtractor_ResolveBehaviorNameResult() {
	_ = this_
}

// Java constructor: ResolveBehaviorNameResult(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node node, boolean isGlobalDeclaration, com.google.javascript.jscomp.modules.ModuleMetadataMap.@org.jspecify.annotations.Nullable ModuleMetadata moduleMetadata)
func NewPolymerBehaviorExtractor_ResolveBehaviorNameResult(node pkg_rhino.Node, isGlobalDeclaration bool, moduleMetadata ModuleMetadataMap_ModuleMetadata) *PolymerBehaviorExtractor_ResolveBehaviorNameResult {
	this_ := &PolymerBehaviorExtractor_ResolveBehaviorNameResult{}
	InitPolymerBehaviorExtractor_ResolveBehaviorNameResult(this_, node, isGlobalDeclaration, moduleMetadata)
	return this_
}

func InitPolymerBehaviorExtractor_ResolveBehaviorNameResult(this_ *PolymerBehaviorExtractor_ResolveBehaviorNameResult, node pkg_rhino.Node, isGlobalDeclaration bool, moduleMetadata ModuleMetadataMap_ModuleMetadata) {
	_ = node
	_ = isGlobalDeclaration
	_ = moduleMetadata
	this_.initFields_PolymerBehaviorExtractor_ResolveBehaviorNameResult()
	this_.node = node
	this_.isGlobalDeclaration = isGlobalDeclaration
	this_.moduleMetadata = moduleMetadata
}

// Java class: com.google.javascript.jscomp.PolymerBehaviorExtractor.BehaviorDefinition
type PolymerBehaviorExtractor_BehaviorDefinition struct {
	props                    j2gort.Ext_java_util_List `gson:"props"`
	functionsToCopy          j2gort.Ext_java_util_List `gson:"functionsToCopy"`
	nonPropertyMembersToCopy j2gort.Ext_java_util_List `gson:"nonPropertyMembersToCopy"`
	features                 *pkg_parser.FeatureSet    `gson:"features"`
	lazyModuleLocalNames     j2gort.Ext_java_util_Set  `gson:"lazyModuleLocalNames"`
	behaviorModule           pkg_rhino.Node            `gson:"behaviorModule"`
	idHash_                  atomic.Int32
	isGlobalDeclaration      bool `gson:"isGlobalDeclaration"`
}

func (this_ *PolymerBehaviorExtractor_BehaviorDefinition) BindSelf_(s any) {
	_ = s
}
func (this_ *PolymerBehaviorExtractor_BehaviorDefinition) IsJavaNil_() bool { return this_ == nil }
func (*PolymerBehaviorExtractor_BehaviorDefinition) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolymerBehaviorExtractor$BehaviorDefinition", "BehaviorDefinition"
}
func (this_ *PolymerBehaviorExtractor_BehaviorDefinition) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PolymerBehaviorExtractor_BehaviorDefinition) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolymerBehaviorExtractor.BehaviorDefinition", "com.google.javascript.jscomp.PolymerBehaviorExtractor$BehaviorDefinition", func(v_ any) bool {
		return j2gort.InstanceOf[*PolymerBehaviorExtractor_BehaviorDefinition](v_)
	})
}

func (this_ *PolymerBehaviorExtractor_BehaviorDefinition) initFields_PolymerBehaviorExtractor_BehaviorDefinition() {
	_ = this_
}

// Java constructor: BehaviorDefinition(java.util.List<com.google.javascript.jscomp.PolymerPass.MemberDefinition> props, java.util.List<com.google.javascript.jscomp.PolymerPass.MemberDefinition> functionsToCopy, java.util.List<com.google.javascript.jscomp.PolymerPass.MemberDefinition> nonPropertyMembersToCopy, boolean isGlobalDeclaration, com.google.javascript.jscomp.parsing.parser.FeatureSet features, com.google.javascript.rhino.Node behaviorModule)
func NewPolymerBehaviorExtractor_BehaviorDefinition(props j2gort.Ext_java_util_List, functionsToCopy j2gort.Ext_java_util_List, nonPropertyMembersToCopy j2gort.Ext_java_util_List, isGlobalDeclaration bool, features *pkg_parser.FeatureSet, behaviorModule pkg_rhino.Node) *PolymerBehaviorExtractor_BehaviorDefinition {
	this_ := &PolymerBehaviorExtractor_BehaviorDefinition{}
	InitPolymerBehaviorExtractor_BehaviorDefinition(this_, props, functionsToCopy, nonPropertyMembersToCopy, isGlobalDeclaration, features, behaviorModule)
	return this_
}

func InitPolymerBehaviorExtractor_BehaviorDefinition(this_ *PolymerBehaviorExtractor_BehaviorDefinition, props j2gort.Ext_java_util_List, functionsToCopy j2gort.Ext_java_util_List, nonPropertyMembersToCopy j2gort.Ext_java_util_List, isGlobalDeclaration bool, features *pkg_parser.FeatureSet, behaviorModule pkg_rhino.Node) {
	_ = props
	_ = functionsToCopy
	_ = nonPropertyMembersToCopy
	_ = isGlobalDeclaration
	_ = features
	_ = behaviorModule
	this_.initFields_PolymerBehaviorExtractor_BehaviorDefinition()
	this_.props = props
	this_.functionsToCopy = functionsToCopy
	this_.nonPropertyMembersToCopy = nonPropertyMembersToCopy
	this_.isGlobalDeclaration = isGlobalDeclaration
	this_.features = features
	this_.behaviorModule = behaviorModule
}

// Java method: BehaviorDefinition.getModuleLocalNames(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PolymerBehaviorExtractor_BehaviorDefinition) GetModuleLocalNames(compiler AbstractCompiler) j2gort.Ext_java_util_Set {
	_ = compiler
	if this_.lazyModuleLocalNames == nil {
		this_.lazyModuleLocalNames = this_.priv_BehaviorDefinition_accumulateModuleLocalVars_4370(compiler)
	}
	return this_.lazyModuleLocalNames
}

// Java method: BehaviorDefinition.accumulateModuleLocalVars(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PolymerBehaviorExtractor_BehaviorDefinition) priv_BehaviorDefinition_accumulateModuleLocalVars_4370(compiler AbstractCompiler) j2gort.Ext_java_util_Set {
	_ = compiler
	clinit_SyntacticScopeCreator()
	var val_7 SyntacticScopeCreator
	InitSyntacticScopeCreator_1(&val_7, compiler)
	var scopeCreator *SyntacticScopeCreator = &val_7
	_ = scopeCreator
	var globalScope *Scope = Scope_CreateGlobalScope(this_.behaviorModule.GetParent())
	_ = globalScope
	return NodeUtil_GetAllVarNamesDeclaredInModule(this_.behaviorModule, compiler, scopeCreator, globalScope)
}
