// Port of com/google/javascript/jscomp/deps/WebpackModuleResolver.java @ closure-compiler v20261015

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

package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

// WebpackModuleResolver is the resolution algorithm for Webpack. Modules are
// located by a map of numeric ids to module paths.
//
// As the compiler normally locates modules by path string, webpack numeric
// ids are converted to strings.
type WebpackModuleResolver struct {
	NodeModuleResolver
	modulesById *jutil.LinkedHashMap[string, string]
}

var _ ModuleResolver = (*WebpackModuleResolver)(nil)

// WebpackModuleResolver_Factory uses a lookup map provided by webpack to
// locate modules from a numeric id used during import.
type WebpackModuleResolver_Factory struct {
	lookupMap *jutil.LinkedHashMap[string, string]
}

var _ ModuleLoader_ModuleResolverFactory = (*WebpackModuleResolver_Factory)(nil)

// Java: WebpackModuleResolver.Factory(Map<String, String>)
func NewWebpackModuleResolver_Factory(lookupMap *jutil.LinkedHashMap[string, string]) *WebpackModuleResolver_Factory {
	return &WebpackModuleResolver_Factory{lookupMap: lookupMap}
}

// Java: WebpackModuleResolver.Factory.create(ImmutableSet<String>, ImmutableList<String>, ErrorHandler, PathEscaper)
func (f *WebpackModuleResolver_Factory) Create(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) ModuleResolver {
	normalizedPathsById := jutil.NewLinkedHashMap[string, string]()
	for key, value := range f.lookupMap.All() {
		canonicalizedPath :=
			ModuleLoader_Normalize(pathEscaper.Escape(value), moduleRootPaths)
		if ModuleLoader_IsAmbiguousIdentifier(canonicalizedPath) {
			canonicalizedPath = ModuleLoader_MODULE_SLASH + canonicalizedPath
		}
		normalizedPathsById.Put(key, canonicalizedPath)
	}
	return NewWebpackModuleResolver(
		modulePaths, moduleRootPaths, normalizedPathsById, errorHandler, pathEscaper)
}

// Java: WebpackModuleResolver(ImmutableSet<String>, ImmutableList<String>, Map<String, String>, ErrorHandler, PathEscaper)
func NewWebpackModuleResolver(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	modulesById *jutil.LinkedHashMap[string, string],
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) *WebpackModuleResolver {
	r := &WebpackModuleResolver{
		NodeModuleResolver: *NewNodeModuleResolver(
			modulePaths,
			moduleRootPaths,
			/* packageJsonMainEntries= */ nil,
			errorHandler,
			pathEscaper),
	}
	r.modulesById = jutil.NewLinkedHashMapFrom(modulesById)
	return r
}

// ResolveJsModule may return nil.
//
// Java: WebpackModuleResolver.resolveJsModule(String, String, String, int, int)
func (r *WebpackModuleResolver) ResolveJsModule(
	scriptAddress string, moduleAddress string, sourcename string, lineno int, colno int) *string {
	loadAddress, ok := r.modulesById.Lookup(moduleAddress)
	if !ok {
		// Module paths may still be used in type nodes so we need to fall back
		// to node module resolution for those.
		return r.NodeModuleResolver.ResolveJsModule(scriptAddress, moduleAddress, sourcename, lineno, colno)
	}
	return &loadAddress
}

// ResolveJsModuleSilently may return nil.
//
// Java: WebpackModuleResolver.resolveJsModuleSilently(String, String)
func (r *WebpackModuleResolver) ResolveJsModuleSilently(scriptAddress string, moduleAddress string) *string {
	loadAddress, ok := r.modulesById.Lookup(moduleAddress)
	if !ok {
		// Module paths may still be used in type nodes so we need to fall back
		// to node module resolution for those.
		return r.NodeModuleResolver.ResolveJsModuleSilently(scriptAddress, moduleAddress)
	}
	return &loadAddress
}
