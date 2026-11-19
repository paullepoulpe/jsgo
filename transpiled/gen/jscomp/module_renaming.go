/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ModuleRenaming.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ModuleRenaming
type ModuleRenaming struct {
	idHash_ atomic.Int32
}

func (this_ *ModuleRenaming) BindSelf_(s any) {
	_ = s
}
func (this_ *ModuleRenaming) IsJavaNil_() bool { return this_ == nil }
func (*ModuleRenaming) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ModuleRenaming", "ModuleRenaming"
}
func (this_ *ModuleRenaming) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ModuleRenaming) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ModuleRenaming", "com.google.javascript.jscomp.ModuleRenaming", func(v_ any) bool {
		return j2gort.InstanceOf[*ModuleRenaming](v_)
	})
}

func StaticGet_ModuleRenaming_DEFAULT_EXPORT_VAR_PREFIX() string { return "$jscompDefaultExport" }
func (this_ *ModuleRenaming) initFields_ModuleRenaming() {
	_ = this_
}

// Java constructor: ModuleRenaming()
func NewModuleRenaming() *ModuleRenaming {
	this_ := &ModuleRenaming{}
	InitModuleRenaming(this_)
	return this_
}

func InitModuleRenaming(this_ *ModuleRenaming) {
	this_.initFields_ModuleRenaming()
}

// Java static method: ModuleRenaming.getGlobalNameOfEsModuleLocalVariable(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata moduleMetadata, java.lang.String variableName)
func ModuleRenaming_GetGlobalNameOfEsModuleLocalVariable(moduleMetadata ModuleMetadataMap_ModuleMetadata, variableName string) pkg_rhino.QualifiedName {
	_ = moduleMetadata
	_ = variableName
	return pkg_rhino.QualifiedName_Of((j2gort.StrValString(variableName) + "$$" + j2gort.StrValString(ModuleRenaming_GetGlobalName_2(moduleMetadata, j2gort.NullString).Join_0())))
}

// Java static method: ModuleRenaming.getGlobalNameOfAnonymousDefaultExport(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata moduleMetadata)
func ModuleRenaming_GetGlobalNameOfAnonymousDefaultExport(moduleMetadata ModuleMetadataMap_ModuleMetadata) pkg_rhino.QualifiedName {
	_ = moduleMetadata
	return ModuleRenaming_GetGlobalNameOfEsModuleLocalVariable(moduleMetadata, "$jscompDefaultExport")
}

// Java static method: ModuleRenaming.getGlobalName(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata moduleMetadata, java.lang.@org.jspecify.annotations.Nullable String googNamespace)
func ModuleRenaming_GetGlobalName_2(moduleMetadata ModuleMetadataMap_ModuleMetadata, googNamespace string) pkg_rhino.QualifiedName {
	_ = moduleMetadata
	_ = googNamespace
	return ModuleRenaming_GlobalizedModuleName_Create_3(moduleMetadata, googNamespace, nil).AliasName()
}

// Java static method: ModuleRenaming.getGlobalName(com.google.javascript.jscomp.modules.Export export)
func ModuleRenaming_GetGlobalName_1_Export(export *Export) pkg_rhino.QualifiedName {
	_ = export
	if export.ModuleMetadata().IsEs6Module() {
		{
			var iifeCond_1 bool
			{

				sr_1 := export.LocalName()
				se_2 := "*default*"
				iifeCond_1 = j2gort.CheckNotNullString(sr_1) == se_2
			}
			if iifeCond_1 {
				return ModuleRenaming_GetGlobalNameOfAnonymousDefaultExport(export.ModuleMetadata())
			}
		}
		return ModuleRenaming_GetGlobalNameOfEsModuleLocalVariable(export.ModuleMetadata(), export.LocalName())
	}
	return ModuleRenaming_GetGlobalName_2(export.ModuleMetadata(), export.ClosureNamespace()).Getprop(export.ExportName())
}

// Java static method: ModuleRenaming.getGlobalName(com.google.javascript.jscomp.modules.Binding binding)
func ModuleRenaming_GetGlobalName_1_Binding(binding *Binding) pkg_rhino.QualifiedName {
	_ = binding
	if binding.IsModuleNamespace() {
		return ModuleRenaming_GetGlobalName_2(binding.Metadata(), binding.ClosureNamespace())
	}
	return ModuleRenaming_GetGlobalName_1_Export(binding.OriginatingExport())
}

// Java static method: ModuleRenaming.getNameRootType(java.lang.String qname, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable TypedScope globalTypedScope)
func ModuleRenaming_GetNameRootType(qname string, globalTypedScope *TypedScope) pkg_rhino.JSType {
	_ = qname
	_ = globalTypedScope
	if globalTypedScope == nil {
		return nil
	}
	var root string = NodeUtil_GetRootOfQualifiedName_1_string(qname)
	_ = root
	return j2gort.CheckNotNullPtr(j2gort.CastPtr[TypedVar](globalTypedScope.GetVar(root))).GetType()
}

// Java static method: ModuleRenaming.getGlobalNameForJsDoc(com.google.javascript.jscomp.modules.ModuleMap moduleMap, com.google.javascript.jscomp.modules.Binding binding, java.util.List<java.lang.String> propertyChain)
func ModuleRenaming_GetGlobalNameForJsDoc(moduleMap *ModuleMap, binding *Binding, propertyChain j2gort.Ext_java_util_List) string {
	_ = moduleMap
	_ = binding
	_ = propertyChain
	var prop int32 = int32(0)
	_ = prop
	for (binding.IsModuleNamespace() && binding.Metadata().IsEs6Module()) && (prop < propertyChain.Size()) {
		var propertyName string = j2gort.ListGetStr(propertyChain, prop)
		_ = propertyName
		var m *Module = moduleMap.GetModule_1_ModuleLoader_ModulePath(binding.Metadata().Path())
		_ = m
		if j2gort.MapContainsKeyStr(m.Namespace(), propertyName) {
			binding = j2gort.CastPtr[Binding](j2gort.MapGetStr(m.Namespace(), propertyName))
		} else {
			break
		}
		prop = prop + 1
	}
	var globalName pkg_rhino.QualifiedName = ModuleRenaming_GetGlobalName_1_Binding(binding)
	_ = globalName
	for prop < propertyChain.Size() {
		globalName = globalName.Getprop(j2gort.ListGetStr(propertyChain, func() int32 { old_3 := prop; nw_4 := old_3 + 1; prop = nw_4; return old_3 }()))
	}
	return globalName.Join_0()
}

// Java record: com.google.javascript.jscomp.ModuleRenaming.GlobalizedModuleName
type ModuleRenaming_GlobalizedModuleName struct {
	aliasName    pkg_rhino.QualifiedName `gson:"aliasName"`
	rootNameType pkg_rhino.JSType        `gson:"rootNameType"`
}

func (this_ *ModuleRenaming_GlobalizedModuleName) BindSelf_(s any) {
	_ = s
}
func (this_ *ModuleRenaming_GlobalizedModuleName) IsJavaNil_() bool { return this_ == nil }
func (*ModuleRenaming_GlobalizedModuleName) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ModuleRenaming$GlobalizedModuleName", "GlobalizedModuleName"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ModuleRenaming.GlobalizedModuleName", "com.google.javascript.jscomp.ModuleRenaming$GlobalizedModuleName", func(v_ any) bool {
		return j2gort.InstanceOf[*ModuleRenaming_GlobalizedModuleName](v_)
	})
}
func (*ModuleRenaming_GlobalizedModuleName) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *ModuleRenaming_GlobalizedModuleName) initFields_ModuleRenaming_GlobalizedModuleName() {
	_ = this_
}

// Java constructor: GlobalizedModuleName(com.google.javascript.rhino.QualifiedName aliasName, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType rootNameType)
func NewModuleRenaming_GlobalizedModuleName(aliasName pkg_rhino.QualifiedName, rootNameType pkg_rhino.JSType) *ModuleRenaming_GlobalizedModuleName {
	this_ := &ModuleRenaming_GlobalizedModuleName{}
	InitModuleRenaming_GlobalizedModuleName(this_, aliasName, rootNameType)
	return this_
}

func InitModuleRenaming_GlobalizedModuleName(this_ *ModuleRenaming_GlobalizedModuleName, aliasName pkg_rhino.QualifiedName, rootNameType pkg_rhino.JSType) {
	_ = aliasName
	_ = rootNameType
	this_.initFields_ModuleRenaming_GlobalizedModuleName()
	j2gort.CheckNotNullIface(aliasName)
	this_.aliasName = aliasName
	this_.rootNameType = rootNameType
}

// Java method: GlobalizedModuleName.toQname(com.google.javascript.jscomp.AstFactory astFactory)
func (this_ *ModuleRenaming_GlobalizedModuleName) ToQname(astFactory *AstFactory) pkg_rhino.Node {
	_ = astFactory
	var rootName pkg_rhino.Node = astFactory.CreateName_2_string_AstFactory_Type(this_.AliasName().GetRoot(), AstFactory_S_Type_1_JSType(this_.RootNameType()))
	_ = rootName
	if this_.AliasName().IsSimple() {
		return rootName
	}
	return astFactory.CreateGetPropsWithoutColors(rootName, j2gort.Ext_com_google_common_collect_Iterables_Skip(this_.AliasName().Components(), int32(1)))
}

// Java method: GlobalizedModuleName.getprop(java.lang.String property)
func (this_ *ModuleRenaming_GlobalizedModuleName) Getprop(property string) *ModuleRenaming_GlobalizedModuleName {
	_ = property
	if !(!(j2gort.CheckNotNullString(property) == "") && !(j2gort.JStrIndexOf(property, ".") >= 0)) {
		j2gort.CheckArgument(false)
	}
	return ModuleRenaming_GlobalizedModuleName_Create_2(this_.AliasName().Getprop(property), this_.RootNameType())
}

// Java static method: GlobalizedModuleName.create(com.google.javascript.rhino.QualifiedName aliasName, com.google.javascript.rhino.jstype.JSType rootNameType)
func ModuleRenaming_GlobalizedModuleName_Create_2(aliasName pkg_rhino.QualifiedName, rootNameType pkg_rhino.JSType) *ModuleRenaming_GlobalizedModuleName {
	_ = aliasName
	_ = rootNameType
	return NewModuleRenaming_GlobalizedModuleName(aliasName, rootNameType)
}

// Java static method: GlobalizedModuleName.create(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata moduleMetadata, java.lang.@org.jspecify.annotations.Nullable String googNamespace, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable TypedScope globalTypedScope)
func ModuleRenaming_GlobalizedModuleName_Create_3(moduleMetadata ModuleMetadataMap_ModuleMetadata, googNamespace string, globalTypedScope *TypedScope) *ModuleRenaming_GlobalizedModuleName {
	_ = moduleMetadata
	_ = googNamespace
	_ = globalTypedScope
	if !((googNamespace == j2gort.NullString) || j2gort.SetContainsStr(moduleMetadata.GoogNamespaces(), googNamespace)) {
		j2gort.CheckState(false)
	}
	switch (moduleMetadata.ModuleType()).Ordinal() {
	case int32(1) /* ModuleType.GOOG_MODULE */ :
		var moduleBody pkg_rhino.Node = moduleMetadata.RootNode().GetFirstChild()
		_ = moduleBody
		return NewModuleRenaming_GlobalizedModuleName(pkg_rhino.QualifiedName_Of(ClosureRewriteModule_GetBinaryModuleNamespace(googNamespace)), moduleBody.GetJSType())
	case int32(4) /* ModuleType.GOOG_PROVIDE */, int32(2) /* ModuleType.LEGACY_GOOG_MODULE */ :
		return NewModuleRenaming_GlobalizedModuleName(pkg_rhino.QualifiedName_Of(googNamespace), ModuleRenaming_GetNameRootType(googNamespace, globalTypedScope))
	case int32(0) /* ModuleType.ES6_MODULE */, int32(3) /* ModuleType.COMMON_JS */ :
		return NewModuleRenaming_GlobalizedModuleName(pkg_rhino.QualifiedName_Of(moduleMetadata.Path().ToModuleName()), ModuleRenaming_GetNameRootType(moduleMetadata.Path().ToModuleName(), globalTypedScope))
	case int32(5) /* ModuleType.SCRIPT */ :
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected module type: "+j2gort.StrVal(moduleMetadata.ModuleType()))))
}

func (this_ *ModuleRenaming_GlobalizedModuleName) AliasName() pkg_rhino.QualifiedName {
	return this_.aliasName
}
func (this_ *ModuleRenaming_GlobalizedModuleName) RootNameType() pkg_rhino.JSType {
	return this_.rootNameType
}
func (this_ *ModuleRenaming_GlobalizedModuleName) Equals(other_ any) bool {
	o_, ok_ := other_.(*ModuleRenaming_GlobalizedModuleName)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if !javalang.ObjectsEquals(this_.aliasName, o_.aliasName) {
		return false
	}
	if !javalang.ObjectsEquals(this_.rootNameType, o_.rootNameType) {
		return false
	}
	return true
}
func (this_ *ModuleRenaming_GlobalizedModuleName) HashCode() int32 {
	var h_ int32 = 0
	h_ = 31*h_ + j2gort.ObjectHash(this_.aliasName)
	{
		hi_5 := this_.rootNameType
		h_ = 31 * h_
		if hi_5 != nil && !hi_5.IsJavaNil_() {
			h_ += hi_5.HashCode()
		}
	}
	return h_
}
func (this_ *ModuleRenaming_GlobalizedModuleName) String() string {
	_ = fmt.Sprint
	return "GlobalizedModuleName[" + "aliasName=" + j2gort.StrVal(this_.aliasName) + ", " + "rootNameType=" + j2gort.StrVal(this_.rootNameType) + "]"
}
