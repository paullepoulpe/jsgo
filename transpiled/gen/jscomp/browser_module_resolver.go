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
// Port of com/google/javascript/jscomp/deps/BrowserModuleResolver.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

var lambdaSingleton_browser_module_resolver_1 ModuleLoader_ModuleResolverFactory

func init() {
	lambdaSingleton_browser_module_resolver_1 = NewLambda_ModuleLoader_ModuleResolverFactory(func(mrp0_ j2gort.Ext_com_google_common_collect_ImmutableSet, mrp1_ j2gort.Ext_com_google_common_collect_ImmutableList, mrp2_ ErrorHandler, mrp3_ *ModuleLoader_PathEscaper) ModuleResolver {
		return NewBrowserModuleResolver(mrp0_, mrp1_, mrp2_, mrp3_)
	})
}

// Java class: com.google.javascript.jscomp.deps.BrowserModuleResolver
type BrowserModuleResolver struct {
	ModuleResolverBase
}

func (this_ *BrowserModuleResolver) BindSelf_(s any) {
	_ = s
}
func (this_ *BrowserModuleResolver) IsJavaNil_() bool { return this_ == nil }
func (*BrowserModuleResolver) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.BrowserModuleResolver", "BrowserModuleResolver"
}
func (this_ *BrowserModuleResolver) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ModuleResolverBase.IdentityHash_()
}
func (*BrowserModuleResolver) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.BrowserModuleResolver", "com.google.javascript.jscomp.deps.BrowserModuleResolver", func(v_ any) bool {
		return j2gort.InstanceOf[*BrowserModuleResolver](v_)
	})
}
func (*BrowserModuleResolver) Implements_jscomp_ModuleResolver() {}

var static_BrowserModuleResolver_FACTORY ModuleLoader_ModuleResolverFactory

func StaticGet_BrowserModuleResolver_FACTORY() ModuleLoader_ModuleResolverFactory {
	clinit_BrowserModuleResolver()
	return static_BrowserModuleResolver_FACTORY
}
func StaticSet_BrowserModuleResolver_FACTORY(v ModuleLoader_ModuleResolverFactory) {
	clinit_BrowserModuleResolver()
	static_BrowserModuleResolver_FACTORY = v
}

var clinitDone_BrowserModuleResolver uint32
var clinitState_BrowserModuleResolver j2gort.ClinitState

func clinit_BrowserModuleResolver() {
	if atomic.LoadUint32(&clinitDone_BrowserModuleResolver) != 0 {
		return
	}
	clinitSlow_BrowserModuleResolver()
}

func clinitSlow_BrowserModuleResolver() {
	if !j2gort.EnterClinit(&clinitState_BrowserModuleResolver, &clinitDone_BrowserModuleResolver) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_BrowserModuleResolver, &clinitDone_BrowserModuleResolver)
	static_BrowserModuleResolver_FACTORY = lambdaSingleton_browser_module_resolver_1
}

func (this_ *BrowserModuleResolver) initFields_BrowserModuleResolver() {
	_ = this_
}

// Java constructor: BrowserModuleResolver(com.google.common.collect.ImmutableSet<java.lang.String> modulePaths, com.google.common.collect.ImmutableList<java.lang.String> moduleRootPaths, com.google.javascript.jscomp.ErrorHandler errorHandler, com.google.javascript.jscomp.deps.ModuleLoader.PathEscaper pathEscaper)
func NewBrowserModuleResolver(modulePaths j2gort.Ext_com_google_common_collect_ImmutableSet, moduleRootPaths j2gort.Ext_com_google_common_collect_ImmutableList, errorHandler ErrorHandler, pathEscaper *ModuleLoader_PathEscaper) *BrowserModuleResolver {
	this_ := &BrowserModuleResolver{}
	InitBrowserModuleResolver(this_, modulePaths, moduleRootPaths, errorHandler, pathEscaper)
	return this_
}

func InitBrowserModuleResolver(this_ *BrowserModuleResolver, modulePaths j2gort.Ext_com_google_common_collect_ImmutableSet, moduleRootPaths j2gort.Ext_com_google_common_collect_ImmutableList, errorHandler ErrorHandler, pathEscaper *ModuleLoader_PathEscaper) {
	_ = modulePaths
	_ = moduleRootPaths
	_ = errorHandler
	_ = pathEscaper
	InitModuleResolver(&this_.ModuleResolverBase, modulePaths, moduleRootPaths, errorHandler, pathEscaper)
	this_.initFields_BrowserModuleResolver()
}

// Java method: BrowserModuleResolver.resolveJsModule(java.lang.String scriptAddress, java.lang.String moduleAddress, java.lang.String sourcename, int lineno, int colno)
func (this_ *BrowserModuleResolver) ResolveJsModule(scriptAddress string, moduleAddress string, sourcename string, lineno int32, colno int32) string {
	_ = scriptAddress
	_ = moduleAddress
	_ = sourcename
	_ = lineno
	_ = colno
	if ModuleLoader_IsAmbiguousIdentifier(moduleAddress) {
		this_.GetModuleResolverBase().errorHandler.Report_2(StaticGet_CheckLevel_WARNING(), JSError_Make_5(sourcename, lineno, colno, StaticGet_ModuleLoader_INVALID_MODULE_PATH(), []string{moduleAddress, j2gort.CheckNotNullEnum(StaticGet_ModuleLoader_ResolutionMode_BROWSER()).String()}))
		return j2gort.NullString
	}
	var loadAddress string = this_.Locate(scriptAddress, moduleAddress)
	_ = loadAddress
	if loadAddress == j2gort.NullString {
		this_.GetModuleResolverBase().errorHandler.Report_2(StaticGet_CheckLevel_WARNING(), JSError_Make_5(sourcename, lineno, colno, StaticGet_ModuleLoader_LOAD_WARNING(), []string{moduleAddress}))
	}
	return loadAddress
}

// Java method: BrowserModuleResolver.resolveJsModuleSilently(java.lang.String scriptAddress, java.lang.String moduleAddress)
func (this_ *BrowserModuleResolver) ResolveJsModuleSilently(scriptAddress string, moduleAddress string) string {
	_ = scriptAddress
	_ = moduleAddress
	if ModuleLoader_IsAmbiguousIdentifier(moduleAddress) {
		return j2gort.NullString
	} else {
		return this_.Locate(scriptAddress, moduleAddress)
	}
}
