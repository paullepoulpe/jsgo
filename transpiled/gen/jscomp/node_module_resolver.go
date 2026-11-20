/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/deps/NodeModuleResolver.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_node_module_resolver_1 j2gort.Ext_java_util_Comparator

func init() {
	lambdaSingleton_node_module_resolver_1 = j2gort.NewLambda_Ext_java_util_Comparator(func(lp_1 any, lp_2 any) int32 {
		a := j2gort.CastString(lp_1)
		_ = a
		b := j2gort.CastString(lp_2)
		_ = b
		var comparison int32 = j2gort.CompareValues(j2gort.ValueOfInt(j2gort.JStrLen(b)), j2gort.ValueOfInt(j2gort.JStrLen(a)))
		_ = comparison
		if comparison != int32(0) {
			return comparison
		}
		return j2gort.JStrCompareTo(a, b)
	})
}

// Java polymorphic class interface: com.google.javascript.jscomp.deps.NodeModuleResolver
type NodeModuleResolver interface {
	Implements_jscomp_NodeModuleResolver()
	IsJavaNil_() bool
	GetNodeModuleResolverBase() *NodeModuleResolverBase
	ModuleResolver
	ResolveJsModuleFile(string, string) string
}

func ToIface_NodeModuleResolver[T interface {
	NodeModuleResolver
	comparable
}](v T) NodeModuleResolver {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_NodeModuleResolver[T NodeModuleResolver](v T) NodeModuleResolver {
	var vi NodeModuleResolver = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.deps.NodeModuleResolver
type NodeModuleResolverBase struct {
	ModuleResolverBase
	packageJsonMainEntries j2gort.Ext_com_google_common_collect_ImmutableMap       `gson:"packageJsonMainEntries"`
	nodeModulesFolders     j2gort.Ext_com_google_common_collect_ImmutableSortedSet `gson:"nodeModulesFolders"`
}

func (this_ *NodeModuleResolverBase) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeModuleResolverBase) IsJavaNil_() bool { return this_ == nil }
func (*NodeModuleResolverBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.NodeModuleResolver", "NodeModuleResolver"
}
func (this_ *NodeModuleResolverBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ModuleResolverBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.NodeModuleResolver", "com.google.javascript.jscomp.deps.NodeModuleResolver", func(v_ any) bool {
		return j2gort.InstanceOf[NodeModuleResolver](v_)
	})
}
func (b *NodeModuleResolverBase) GetNodeModuleResolverBase() *NodeModuleResolverBase { return b }
func (*NodeModuleResolverBase) Implements_jscomp_NodeModuleResolver()                {}
func (*NodeModuleResolverBase) Implements_jscomp_ModuleResolver()                    {}

var static_NodeModuleResolver_FILE_EXTENSIONS_TO_SEARCH []string

func StaticGet_NodeModuleResolver_FILE_EXTENSIONS_TO_SEARCH() []string {
	clinit_NodeModuleResolver()
	return static_NodeModuleResolver_FILE_EXTENSIONS_TO_SEARCH
}
func StaticSet_NodeModuleResolver_FILE_EXTENSIONS_TO_SEARCH(v []string) {
	clinit_NodeModuleResolver()
	static_NodeModuleResolver_FILE_EXTENSIONS_TO_SEARCH = v
}

var static_NodeModuleResolver_FILES_TO_SEARCH []string

func StaticGet_NodeModuleResolver_FILES_TO_SEARCH() []string {
	clinit_NodeModuleResolver()
	return static_NodeModuleResolver_FILES_TO_SEARCH
}
func StaticSet_NodeModuleResolver_FILES_TO_SEARCH(v []string) {
	clinit_NodeModuleResolver()
	static_NodeModuleResolver_FILES_TO_SEARCH = v
}

var clinitDone_NodeModuleResolver uint32
var clinitState_NodeModuleResolver j2gort.ClinitState

func clinit_NodeModuleResolver() {
	if atomic.LoadUint32(&clinitDone_NodeModuleResolver) != 0 {
		return
	}
	clinitSlow_NodeModuleResolver()
}

func clinitSlow_NodeModuleResolver() {
	if !j2gort.EnterClinit(&clinitState_NodeModuleResolver, &clinitDone_NodeModuleResolver) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_NodeModuleResolver, &clinitDone_NodeModuleResolver)
	static_NodeModuleResolver_FILE_EXTENSIONS_TO_SEARCH = []string{"", ".js", ".json", ".i.js", ".js.i.js"}
	static_NodeModuleResolver_FILES_TO_SEARCH = []string{("/" + "package.json"), ("/" + "index.js"), ("/" + "index.json")}
}

func (this_ *NodeModuleResolverBase) initFields_NodeModuleResolver() {
	_ = this_
}

// Java constructor: NodeModuleResolver(com.google.common.collect.ImmutableSet<java.lang.String> modulePaths, com.google.common.collect.ImmutableList<java.lang.String> moduleRootPaths, java.util.Map<java.lang.String,java.lang.String> packageJsonMainEntries, com.google.javascript.jscomp.ErrorHandler errorHandler, com.google.javascript.jscomp.deps.ModuleLoader.PathEscaper pathEscaper)
func NewNodeModuleResolver(modulePaths j2gort.Ext_com_google_common_collect_ImmutableSet, moduleRootPaths j2gort.Ext_com_google_common_collect_ImmutableList, packageJsonMainEntries j2gort.Ext_java_util_Map, errorHandler ErrorHandler, pathEscaper *ModuleLoader_PathEscaper) NodeModuleResolver {
	clinit_NodeModuleResolver()
	this_ := &NodeModuleResolverBase{}
	InitNodeModuleResolver(this_, modulePaths, moduleRootPaths, packageJsonMainEntries, errorHandler, pathEscaper)
	return this_
}

func InitNodeModuleResolver(this_ *NodeModuleResolverBase, modulePaths j2gort.Ext_com_google_common_collect_ImmutableSet, moduleRootPaths j2gort.Ext_com_google_common_collect_ImmutableList, packageJsonMainEntries j2gort.Ext_java_util_Map, errorHandler ErrorHandler, pathEscaper *ModuleLoader_PathEscaper) {
	_ = modulePaths
	_ = moduleRootPaths
	_ = packageJsonMainEntries
	_ = errorHandler
	_ = pathEscaper
	InitModuleResolver(&this_.ModuleResolverBase, modulePaths, moduleRootPaths, errorHandler, pathEscaper)
	this_.initFields_NodeModuleResolver()
	this_.nodeModulesFolders = NodeModuleResolver_BuildNodeModulesFoldersRegistry(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](modulePaths), j2gort.AsIface[j2gort.Ext_java_lang_Iterable](moduleRootPaths))
	if packageJsonMainEntries == nil {
		this_.packageJsonMainEntries = j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()
	} else {
		this_.packageJsonMainEntries = this_.priv_NodeModuleResolver_buildPackageJsonMainEntries_6171(packageJsonMainEntries)
	}
}

// Java static method: NodeModuleResolver.buildNodeModulesFoldersRegistry(java.lang.Iterable<java.lang.String> modulePaths, java.lang.Iterable<java.lang.String> moduleRootPaths)
func NodeModuleResolver_BuildNodeModulesFoldersRegistry(modulePaths j2gort.Ext_java_lang_Iterable, moduleRootPaths j2gort.Ext_java_lang_Iterable) j2gort.Ext_com_google_common_collect_ImmutableSortedSet {
	clinit_NodeModuleResolver()
	_ = modulePaths
	_ = moduleRootPaths
	var registry j2gort.Ext_java_util_SortedSet = j2gort.NewExt_java_util_TreeSet_1_Ext_java_util_Comparator(lambdaSingleton_node_module_resolver_1)
	_ = registry
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(modulePaths))
		for cur_3.HasNext() {
			modulePath := cur_3.NextStr()
			_ = modulePath
			{
				cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(moduleRootPaths))
				for cur_4.HasNext() {
					moduleRootPath := cur_4.NextStr()
					_ = moduleRootPath
					if javalang.StringStartsWith(j2gort.CheckNotNullString(modulePath), j2gort.CheckNotNullString(moduleRootPath)) {
						modulePath = j2gort.JStrSubFrom(modulePath, j2gort.JStrLen(moduleRootPath))
						break
					}
				}
			}
			var nodeModulesDirs []string = j2gort.JStrSplit(modulePath, "/node_modules/")
			_ = nodeModulesDirs
			var parentPath string = ""
			_ = parentPath
			{
				var i int32 = int32(0)
				_ = i
				for ; i < (int32(len(nodeModulesDirs)) - int32(1)); i = i + 1 {
					if (i + int32(1)) < int32(len(nodeModulesDirs)) {
						parentPath = (j2gort.StrValString(parentPath) + j2gort.StrValString(nodeModulesDirs[i]) + "/")
					}
					j2gort.SetAddStr(registry, parentPath)
					parentPath = (j2gort.StrValString(parentPath) + "node_modules/")
				}
			}
		}
	}
	return j2gort.Ext_com_google_common_collect_ImmutableSortedSet_CopyOfSorted(registry)
}

// Java method: NodeModuleResolver.buildPackageJsonMainEntries(java.util.Map<java.lang.String,java.lang.String> packageJsonMainEntries)
func (this_ *NodeModuleResolverBase) priv_NodeModuleResolver_buildPackageJsonMainEntries_6171(packageJsonMainEntries j2gort.Ext_java_util_Map) j2gort.Ext_com_google_common_collect_ImmutableMap {
	_ = packageJsonMainEntries
	var builder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = builder
	{
		cur_5 := j2gort.NewMapEntryCursor(packageJsonMainEntries)
		for cur_5.NextEntry() {
			var entryKey string = cur_5.KeyStr()
			_ = entryKey
			if ModuleLoader_IsAmbiguousIdentifier(entryKey) {
				entryKey = ("/" + j2gort.StrValString(entryKey))
			}
			builder.Put_2(j2gort.BoxString(entryKey), j2gort.BoxString(j2gort.CastString(cur_5.GetValue())))
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](builder.BuildOrThrow())
}

// Java method: NodeModuleResolver.getPackageJsonMainEntries()
func (this_ *NodeModuleResolverBase) GetPackageJsonMainEntries() j2gort.Ext_java_util_Map {
	return j2gort.AsIface[j2gort.Ext_java_util_Map](this_.packageJsonMainEntries)
}

// Java method: NodeModuleResolver.resolveJsModule(java.lang.String scriptAddress, java.lang.String moduleAddress, java.lang.String sourcename, int lineno, int colno)
func (this_ *NodeModuleResolverBase) ResolveJsModule(scriptAddress string, moduleAddress string, sourcename string, lineno int32, colno int32) string {
	_ = scriptAddress
	_ = moduleAddress
	_ = sourcename
	_ = lineno
	_ = colno
	var loadAddress string = this_.priv_NodeModuleResolver_internalResolveJsModule_6172(scriptAddress, moduleAddress)
	_ = loadAddress
	if loadAddress == j2gort.NullString {
		this_.GetModuleResolverBase().errorHandler.Report_2(StaticGet_CheckLevel_WARNING(), JSError_Make_5(sourcename, lineno, colno, StaticGet_ModuleLoader_LOAD_WARNING(), []string{moduleAddress}))
	}
	return loadAddress
}

// Java method: NodeModuleResolver.internalResolveJsModule(java.lang.String scriptAddress, java.lang.String moduleAddress)
func (this_ *NodeModuleResolverBase) priv_NodeModuleResolver_internalResolveJsModule_6172(scriptAddress string, moduleAddress string) string {
	_ = scriptAddress
	_ = moduleAddress
	var loadAddress string = j2gort.NullString
	_ = loadAddress
	if ModuleLoader_IsAbsoluteIdentifier(moduleAddress) || ModuleLoader_IsRelativeIdentifier(moduleAddress) {
		loadAddress = this_.priv_NodeModuleResolver_resolveJsModuleNodeFileOrDirectory_6174(scriptAddress, moduleAddress)
	} else {
		loadAddress = this_.priv_NodeModuleResolver_resolveJsModuleFromRegistry_6176(scriptAddress, moduleAddress)
	}
	return loadAddress
}

// Java method: NodeModuleResolver.resolveJsModuleSilently(java.lang.String scriptAddress, java.lang.String moduleAddress)
func (this_ *NodeModuleResolverBase) ResolveJsModuleSilently(scriptAddress string, moduleAddress string) string {
	_ = scriptAddress
	_ = moduleAddress
	return this_.priv_NodeModuleResolver_internalResolveJsModule_6172(scriptAddress, moduleAddress)
}

// Java method: NodeModuleResolver.resolveJsModuleFile(java.lang.String scriptAddress, java.lang.String moduleAddress)
func (this_ *NodeModuleResolverBase) ResolveJsModuleFile(scriptAddress string, moduleAddress string) string {
	_ = scriptAddress
	_ = moduleAddress
	for _, el_6 := range j2gort.CheckNotNullSlice(StaticGet_NodeModuleResolver_FILE_EXTENSIONS_TO_SEARCH()) {
		extension := el_6
		_ = extension
		var moduleAddressCandidate string = (j2gort.StrValString(moduleAddress) + j2gort.StrValString(extension))
		_ = moduleAddressCandidate
		var canonicalizedCandidatePath string = this_.CanonicalizePath(scriptAddress, moduleAddressCandidate)
		_ = canonicalizedCandidatePath
		if j2gort.MapContainsKeyStr(this_.packageJsonMainEntries, canonicalizedCandidatePath) {
			moduleAddressCandidate = j2gort.CastString(j2gort.MapGetStr(this_.packageJsonMainEntries, canonicalizedCandidatePath))
			if "$jscomp$browser$skiplisted" == moduleAddressCandidate {
				return j2gort.NullString
			}
		}
		var loadAddress string = this_.Locate(scriptAddress, moduleAddressCandidate)
		_ = loadAddress
		if loadAddress != j2gort.NullString {
			return loadAddress
		}
	}
	return j2gort.NullString
}

// Java method: NodeModuleResolver.resolveJsModuleNodeFileOrDirectory(java.lang.String scriptAddress, java.lang.String moduleAddress)
func (this_ *NodeModuleResolverBase) priv_NodeModuleResolver_resolveJsModuleNodeFileOrDirectory_6174(scriptAddress string, moduleAddress string) string {
	_ = scriptAddress
	_ = moduleAddress
	var loadAddress string = j2gort.NullString
	_ = loadAddress
	loadAddress = this_.ResolveJsModuleFile(scriptAddress, moduleAddress)
	if loadAddress == j2gort.NullString {
		loadAddress = this_.priv_NodeModuleResolver_resolveJsModuleNodeDirectory_6175(scriptAddress, moduleAddress)
	}
	return loadAddress
}

// Java method: NodeModuleResolver.resolveJsModuleNodeDirectory(java.lang.String scriptAddress, java.lang.String moduleAddress)
func (this_ *NodeModuleResolverBase) priv_NodeModuleResolver_resolveJsModuleNodeDirectory_6175(scriptAddress string, moduleAddress string) string {
	_ = scriptAddress
	_ = moduleAddress
	if javalang.StringEndsWith(j2gort.CheckNotNullString(moduleAddress), "/") {
		moduleAddress = j2gort.JStrSub(moduleAddress, int32(0), (j2gort.JStrLen(moduleAddress) - int32(1)))
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(StaticGet_NodeModuleResolver_FILES_TO_SEARCH())); i = i + 1 {
			var loadAddress string = this_.Locate(scriptAddress, (j2gort.StrValString(moduleAddress) + j2gort.StrValString(StaticGet_NodeModuleResolver_FILES_TO_SEARCH()[i])))
			_ = loadAddress
			if loadAddress != j2gort.NullString {
				{
					var iifeCond_1 bool
					{

						sr_7 := StaticGet_NodeModuleResolver_FILES_TO_SEARCH()[i]
						se_8 := ("/" + "package.json")
						iifeCond_1 = j2gort.CheckNotNullString(sr_7) == se_8
					}
					if iifeCond_1 {
						if j2gort.MapContainsKeyStr(this_.packageJsonMainEntries, loadAddress) {
							return this_.ResolveJsModuleFile(scriptAddress, j2gort.CastString(j2gort.MapGetStr(this_.packageJsonMainEntries, loadAddress)))
						}
					} else {
						return loadAddress
					}
				}
			}
		}
	}
	return j2gort.NullString
}

// Java method: NodeModuleResolver.resolveJsModuleFromRegistry(java.lang.String scriptAddress, java.lang.String moduleAddress)
func (this_ *NodeModuleResolverBase) priv_NodeModuleResolver_resolveJsModuleFromRegistry_6176(scriptAddress string, moduleAddress string) string {
	_ = scriptAddress
	_ = moduleAddress
	var normalizedScriptAddress string = (j2gort.StrValString(func() string {
		if ModuleLoader_IsAmbiguousIdentifier(scriptAddress) {
			return "/"
		}
		return ""
	}()) + j2gort.StrValString(scriptAddress))
	_ = normalizedScriptAddress
	{
		cur_9 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.nodeModulesFolders))
		for cur_9.HasNext() {
			nodeModulesFolder := cur_9.NextStr()
			_ = nodeModulesFolder
			if !(javalang.StringStartsWith(j2gort.CheckNotNullString(normalizedScriptAddress), j2gort.CheckNotNullString(nodeModulesFolder))) {
				continue
			}
			var fullModulePath string = (j2gort.StrValString(nodeModulesFolder) + "node_modules/" + j2gort.StrValString(moduleAddress))
			_ = fullModulePath
			var loadAddress string = this_.ResolveJsModuleFile(scriptAddress, fullModulePath)
			_ = loadAddress
			if loadAddress == j2gort.NullString {
				loadAddress = this_.priv_NodeModuleResolver_resolveJsModuleNodeDirectory_6175(scriptAddress, fullModulePath)
			}
			if loadAddress != j2gort.NullString {
				return loadAddress
			}
		}
	}
	return j2gort.NullString
}

// Java class: com.google.javascript.jscomp.deps.NodeModuleResolver.Factory
type NodeModuleResolver_Factory struct {
	packageJsonMainEntries j2gort.Ext_java_util_Map `gson:"packageJsonMainEntries"`
	idHash_                atomic.Int32
}

func (this_ *NodeModuleResolver_Factory) BindSelf_(s any) {
	_ = s
}
func (this_ *NodeModuleResolver_Factory) IsJavaNil_() bool { return this_ == nil }
func (*NodeModuleResolver_Factory) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.NodeModuleResolver$Factory", "Factory"
}
func (this_ *NodeModuleResolver_Factory) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NodeModuleResolver_Factory) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.NodeModuleResolver.Factory", "com.google.javascript.jscomp.deps.NodeModuleResolver$Factory", func(v_ any) bool {
		return j2gort.InstanceOf[*NodeModuleResolver_Factory](v_)
	})
}
func (*NodeModuleResolver_Factory) Implements_jscomp_ModuleLoader_ModuleResolverFactory() {}

func (this_ *NodeModuleResolver_Factory) initFields_NodeModuleResolver_Factory() {
	_ = this_
}

// Java constructor: Factory()
func NewNodeModuleResolver_Factory() *NodeModuleResolver_Factory {
	this_ := &NodeModuleResolver_Factory{}
	InitNodeModuleResolver_Factory(this_)
	return this_
}

func InitNodeModuleResolver_Factory(this_ *NodeModuleResolver_Factory) {
	InitNodeModuleResolver_Factory_1(this_, nil)
}

// Java constructor: Factory(java.util.@org.jspecify.annotations.Nullable Map<java.lang.String,java.lang.String> packageJsonMainEntries)
func NewNodeModuleResolver_Factory_1(packageJsonMainEntries j2gort.Ext_java_util_Map) *NodeModuleResolver_Factory {
	this_ := &NodeModuleResolver_Factory{}
	InitNodeModuleResolver_Factory_1(this_, packageJsonMainEntries)
	return this_
}

func InitNodeModuleResolver_Factory_1(this_ *NodeModuleResolver_Factory, packageJsonMainEntries j2gort.Ext_java_util_Map) {
	_ = packageJsonMainEntries
	this_.initFields_NodeModuleResolver_Factory()
	this_.packageJsonMainEntries = packageJsonMainEntries
}

// Java method: Factory.create(com.google.common.collect.ImmutableSet<java.lang.String> modulePaths, com.google.common.collect.ImmutableList<java.lang.String> moduleRootPaths, com.google.javascript.jscomp.ErrorHandler errorHandler, com.google.javascript.jscomp.deps.ModuleLoader.PathEscaper pathEscaper)
func (this_ *NodeModuleResolver_Factory) Create(modulePaths j2gort.Ext_com_google_common_collect_ImmutableSet, moduleRootPaths j2gort.Ext_com_google_common_collect_ImmutableList, errorHandler ErrorHandler, pathEscaper *ModuleLoader_PathEscaper) ModuleResolver {
	_ = modulePaths
	_ = moduleRootPaths
	_ = errorHandler
	_ = pathEscaper
	return Upcast_ModuleResolver(NewNodeModuleResolver(modulePaths, moduleRootPaths, this_.packageJsonMainEntries, errorHandler, pathEscaper))
}
