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
// Port of com/google/javascript/jscomp/deps/ModuleResolver.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.deps.ModuleResolver
type ModuleResolver interface {
	Implements_jscomp_ModuleResolver()
	IsJavaNil_() bool
	GetModuleResolverBase() *ModuleResolverBase
	GetPackageJsonMainEntries() j2gort.Ext_java_util_Map
	ResolveJsModule(string, string, string, int32, int32) string
	ResolveJsModuleSilently(string, string) string
	ResolveModuleAsPath(string, string) string
	Locate(string, string) string
	CanonicalizePath(string, string) string
	SetErrorHandler(ErrorHandler)
}

func ToIface_ModuleResolver[T interface {
	ModuleResolver
	comparable
}](v T) ModuleResolver {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ModuleResolver[T ModuleResolver](v T) ModuleResolver {
	var vi ModuleResolver = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.deps.ModuleResolver
type ModuleResolverBase struct {
	modulePaths     j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"modulePaths"`
	moduleRootPaths []string                                          `gson:"moduleRootPaths"`
	errorHandler    ErrorHandler                                      `gson:"errorHandler"`
	pathEscaper     *ModuleLoader_PathEscaper                         `gson:"pathEscaper"`
	idHash_         atomic.Int32
}

func (this_ *ModuleResolverBase) BindSelf_(s any) {
	_ = s
}
func (this_ *ModuleResolverBase) IsJavaNil_() bool { return this_ == nil }
func (*ModuleResolverBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.ModuleResolver", "ModuleResolver"
}
func (this_ *ModuleResolverBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.ModuleResolver", "com.google.javascript.jscomp.deps.ModuleResolver", func(v_ any) bool {
		return j2gort.InstanceOf[ModuleResolver](v_)
	})
}
func (b *ModuleResolverBase) GetModuleResolverBase() *ModuleResolverBase { return b }
func (*ModuleResolverBase) Implements_jscomp_ModuleResolver()            {}

func (this_ *ModuleResolverBase) initFields_ModuleResolver() {
	_ = this_
}

func InitModuleResolver(this_ *ModuleResolverBase, modulePaths j2gort.Ext_com_google_common_collect_ImmutableSet, moduleRootPaths j2gort.Ext_com_google_common_collect_ImmutableList, errorHandler ErrorHandler, pathEscaper *ModuleLoader_PathEscaper) {
	_ = modulePaths
	_ = moduleRootPaths
	_ = errorHandler
	_ = pathEscaper
	this_.initFields_ModuleResolver()
	this_.modulePaths = modulePaths
	this_.moduleRootPaths = j2gort.ImmListToSliceNullable[string](moduleRootPaths)
	this_.errorHandler = errorHandler
	this_.pathEscaper = pathEscaper
}

// Java method: ModuleResolver.getPackageJsonMainEntries()
func (this_ *ModuleResolverBase) GetPackageJsonMainEntries() j2gort.Ext_java_util_Map {
	return j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0())
}

// Java method: ModuleResolver.resolveJsModule(java.lang.String scriptAddress, java.lang.String moduleAddress, java.lang.String sourcename, int lineno, int colno)
func (this_ *ModuleResolverBase) ResolveJsModule(scriptAddress string, moduleAddress string, sourcename string, lineno int32, colno int32) string {
	_ = scriptAddress
	_ = moduleAddress
	_ = sourcename
	_ = lineno
	_ = colno
	panic("abstract or native method invoked")
}

// Java method: ModuleResolver.resolveJsModuleSilently(java.lang.String scriptAddress, java.lang.String moduleAddress)
func (this_ *ModuleResolverBase) ResolveJsModuleSilently(scriptAddress string, moduleAddress string) string {
	_ = scriptAddress
	_ = moduleAddress
	panic("abstract or native method invoked")
}

// Java method: ModuleResolver.resolveModuleAsPath(java.lang.String scriptAddress, java.lang.String moduleAddress)
func (this_ *ModuleResolverBase) ResolveModuleAsPath(scriptAddress string, moduleAddress string) string {
	_ = scriptAddress
	_ = moduleAddress
	if !(javalang.StringEndsWith(j2gort.CheckNotNullString(moduleAddress), ".js")) {
		moduleAddress = (j2gort.StrValString(moduleAddress) + ".js")
	}
	var path string = this_.pathEscaper.Escape(moduleAddress)
	_ = path
	if ModuleLoader_IsRelativeIdentifier(moduleAddress) {
		var ourPath string = scriptAddress
		_ = ourPath
		var lastIndex int32 = j2gort.JStrLastIndexOf(ourPath, "/")
		_ = lastIndex
		path = ModuleNames_CanonicalizePath((j2gort.StrValString(j2gort.JStrSub(ourPath, int32(0), (lastIndex+j2gort.JStrLen("/")))) + j2gort.StrValString(path)))
	}
	return ModuleLoader_Normalize(path, j2gort.ImmListFromTypedSlice(this_.moduleRootPaths))
}

// Java method: ModuleResolver.locate(java.lang.String scriptAddress, java.lang.String name)
func (this_ *ModuleResolverBase) Locate(scriptAddress string, name string) string {
	_ = scriptAddress
	_ = name
	var canonicalizedPath string = this_.CanonicalizePath(scriptAddress, name)
	_ = canonicalizedPath
	var normalizedPath string = canonicalizedPath
	_ = normalizedPath
	if ModuleLoader_IsAmbiguousIdentifier(canonicalizedPath) {
		normalizedPath = ("/" + j2gort.StrValString(canonicalizedPath))
	}
	if j2gort.SetContainsStr(this_.modulePaths, normalizedPath) {
		return canonicalizedPath
	}
	{

		for _, rootPath := range this_.moduleRootPaths {

			_ = rootPath
			var modulePath string = (j2gort.StrValString(rootPath) + j2gort.StrValString(normalizedPath))
			_ = modulePath
			if j2gort.SetContainsStr(this_.modulePaths, modulePath) {
				return canonicalizedPath
			}
		}
	}
	return j2gort.NullString
}

// Java method: ModuleResolver.canonicalizePath(java.lang.String scriptAddress, java.lang.String moduleAddress)
func (this_ *ModuleResolverBase) CanonicalizePath(scriptAddress string, moduleAddress string) string {
	_ = scriptAddress
	_ = moduleAddress
	var path string = this_.pathEscaper.Escape(moduleAddress)
	_ = path
	if ModuleLoader_IsRelativeIdentifier(moduleAddress) {
		var ourPath string = scriptAddress
		_ = ourPath
		var lastIndex int32 = j2gort.JStrLastIndexOf(ourPath, "/")
		_ = lastIndex
		path = ModuleNames_CanonicalizePath((j2gort.StrValString(j2gort.JStrSub(ourPath, int32(0), (lastIndex+j2gort.JStrLen("/")))) + j2gort.StrValString(path)))
	}
	return path
}

// Java method: ModuleResolver.setErrorHandler(com.google.javascript.jscomp.ErrorHandler errorHandler)
func (this_ *ModuleResolverBase) SetErrorHandler(errorHandler ErrorHandler) {
	_ = errorHandler
	this_.errorHandler = errorHandler
}
