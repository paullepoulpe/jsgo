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
// Port of com/google/javascript/jscomp/deps/WebpackModuleResolver.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.deps.WebpackModuleResolver
type WebpackModuleResolver struct {
	NodeModuleResolverBase
	modulesById j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"modulesById"`
}

func (this_ *WebpackModuleResolver) BindSelf_(s any) {
	_ = s
}
func (this_ *WebpackModuleResolver) IsJavaNil_() bool { return this_ == nil }
func (*WebpackModuleResolver) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.WebpackModuleResolver", "WebpackModuleResolver"
}
func (this_ *WebpackModuleResolver) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeModuleResolverBase.IdentityHash_()
}
func (*WebpackModuleResolver) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.WebpackModuleResolver", "com.google.javascript.jscomp.deps.WebpackModuleResolver", func(v_ any) bool {
		return j2gort.InstanceOf[*WebpackModuleResolver](v_)
	})
}
func (*WebpackModuleResolver) Implements_jscomp_NodeModuleResolver() {}
func (*WebpackModuleResolver) Implements_jscomp_ModuleResolver()     {}

func (this_ *WebpackModuleResolver) initFields_WebpackModuleResolver() {
	_ = this_
}

// Java constructor: WebpackModuleResolver(com.google.common.collect.ImmutableSet<java.lang.String> modulePaths, com.google.common.collect.ImmutableList<java.lang.String> moduleRootPaths, java.util.Map<java.lang.String,java.lang.String> modulesById, com.google.javascript.jscomp.ErrorHandler errorHandler, com.google.javascript.jscomp.deps.ModuleLoader.PathEscaper pathEscaper)
func NewWebpackModuleResolver(modulePaths j2gort.Ext_com_google_common_collect_ImmutableSet, moduleRootPaths j2gort.Ext_com_google_common_collect_ImmutableList, modulesById *j2gort.OrderedMap[string, string], errorHandler ErrorHandler, pathEscaper *ModuleLoader_PathEscaper) *WebpackModuleResolver {
	this_ := &WebpackModuleResolver{}
	InitWebpackModuleResolver(this_, modulePaths, moduleRootPaths, modulesById, errorHandler, pathEscaper)
	return this_
}

func InitWebpackModuleResolver(this_ *WebpackModuleResolver, modulePaths j2gort.Ext_com_google_common_collect_ImmutableSet, moduleRootPaths j2gort.Ext_com_google_common_collect_ImmutableList, modulesById *j2gort.OrderedMap[string, string], errorHandler ErrorHandler, pathEscaper *ModuleLoader_PathEscaper) {
	_ = modulePaths
	_ = moduleRootPaths
	_ = modulesById
	_ = errorHandler
	_ = pathEscaper
	InitNodeModuleResolver(&this_.NodeModuleResolverBase, modulePaths, moduleRootPaths, nil, errorHandler, pathEscaper)
	this_.initFields_WebpackModuleResolver()
	this_.modulesById = modulesById.ToImmutableMap()
}

// Java method: WebpackModuleResolver.resolveJsModule(java.lang.String scriptAddress, java.lang.String moduleAddress, java.lang.String sourcename, int lineno, int colno)
func (this_ *WebpackModuleResolver) ResolveJsModule(scriptAddress string, moduleAddress string, sourcename string, lineno int32, colno int32) string {
	_ = scriptAddress
	_ = moduleAddress
	_ = sourcename
	_ = lineno
	_ = colno
	var loadAddress string = j2gort.CastString(j2gort.MapGetStr(this_.modulesById, moduleAddress))
	_ = loadAddress
	if loadAddress == j2gort.NullString {
		return this_.NodeModuleResolverBase.ResolveJsModule(scriptAddress, moduleAddress, sourcename, lineno, colno)
	}
	return loadAddress
}

// Java method: WebpackModuleResolver.resolveJsModuleSilently(java.lang.String scriptAddress, java.lang.String moduleAddress)
func (this_ *WebpackModuleResolver) ResolveJsModuleSilently(scriptAddress string, moduleAddress string) string {
	_ = scriptAddress
	_ = moduleAddress
	var loadAddress string = j2gort.CastString(j2gort.MapGetStr(this_.modulesById, moduleAddress))
	_ = loadAddress
	if loadAddress == j2gort.NullString {
		return this_.NodeModuleResolverBase.ResolveJsModuleSilently(scriptAddress, moduleAddress)
	}
	return loadAddress
}

// Java class: com.google.javascript.jscomp.deps.WebpackModuleResolver.Factory
type WebpackModuleResolver_Factory struct {
	lookupMap j2gort.Ext_java_util_Map `gson:"lookupMap"`
	idHash_   atomic.Int32
}

func (this_ *WebpackModuleResolver_Factory) BindSelf_(s any) {
	_ = s
}
func (this_ *WebpackModuleResolver_Factory) IsJavaNil_() bool { return this_ == nil }
func (*WebpackModuleResolver_Factory) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.WebpackModuleResolver$Factory", "Factory"
}
func (this_ *WebpackModuleResolver_Factory) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*WebpackModuleResolver_Factory) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.WebpackModuleResolver.Factory", "com.google.javascript.jscomp.deps.WebpackModuleResolver$Factory", func(v_ any) bool {
		return j2gort.InstanceOf[*WebpackModuleResolver_Factory](v_)
	})
}
func (*WebpackModuleResolver_Factory) Implements_jscomp_ModuleLoader_ModuleResolverFactory() {}

func (this_ *WebpackModuleResolver_Factory) initFields_WebpackModuleResolver_Factory() {
	_ = this_
}

// Java constructor: Factory(java.util.Map<java.lang.String,java.lang.String> lookupMap)
func NewWebpackModuleResolver_Factory(lookupMap j2gort.Ext_java_util_Map) *WebpackModuleResolver_Factory {
	this_ := &WebpackModuleResolver_Factory{}
	InitWebpackModuleResolver_Factory(this_, lookupMap)
	return this_
}

func InitWebpackModuleResolver_Factory(this_ *WebpackModuleResolver_Factory, lookupMap j2gort.Ext_java_util_Map) {
	_ = lookupMap
	this_.initFields_WebpackModuleResolver_Factory()
	this_.lookupMap = lookupMap
}

// Java method: Factory.create(com.google.common.collect.ImmutableSet<java.lang.String> modulePaths, com.google.common.collect.ImmutableList<java.lang.String> moduleRootPaths, com.google.javascript.jscomp.ErrorHandler errorHandler, com.google.javascript.jscomp.deps.ModuleLoader.PathEscaper pathEscaper)
func (this_ *WebpackModuleResolver_Factory) Create(modulePaths j2gort.Ext_com_google_common_collect_ImmutableSet, moduleRootPaths j2gort.Ext_com_google_common_collect_ImmutableList, errorHandler ErrorHandler, pathEscaper *ModuleLoader_PathEscaper) ModuleResolver {
	_ = modulePaths
	_ = moduleRootPaths
	_ = errorHandler
	_ = pathEscaper
	var normalizedPathsById j2gort.OrderedMap[string, string] = j2gort.NewOrderedMap[string, string]()
	_ = normalizedPathsById
	{
		cur_1 := j2gort.NewMapEntryCursor(this_.lookupMap)
		for cur_1.NextEntry() {
			var canonicalizedPath string = ModuleLoader_Normalize(pathEscaper.Escape(j2gort.CastString(cur_1.GetValue())), j2gort.AsIface[j2gort.Ext_java_lang_Iterable](moduleRootPaths))
			_ = canonicalizedPath
			if ModuleLoader_IsAmbiguousIdentifier(canonicalizedPath) {
				canonicalizedPath = ("/" + j2gort.StrValString(canonicalizedPath))
			}
			normalizedPathsById.Put(cur_1.KeyStr(), canonicalizedPath)
		}
	}
	return NewWebpackModuleResolver(modulePaths, moduleRootPaths, &normalizedPathsById, errorHandler, pathEscaper)
}
