// Port of com/google/javascript/jscomp/deps/ModuleResolver.java @ closure-compiler v20261021

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

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ModuleResolver is the base class for algorithms that resolve JavaScript
// module references to input files.
type ModuleResolver interface {
	// Java: ModuleResolver.getPackageJsonMainEntries()
	GetPackageJsonMainEntries() *jutil.LinkedHashMap[string, string]

	// ResolveJsModule may return nil.
	//
	// Java: ModuleResolver.resolveJsModule(String, String, String, int, int)
	ResolveJsModule(scriptAddress string, moduleAddress string, sourcename string, lineno int, colno int) *string

	// ResolveJsModuleSilently may return nil.
	//
	// Java: ModuleResolver.resolveJsModuleSilently(String, String)
	ResolveJsModuleSilently(scriptAddress string, moduleAddress string) *string

	// Java: ModuleResolver.resolveModuleAsPath(String, String)
	ResolveModuleAsPath(scriptAddress string, moduleAddress string) string

	// Java: ModuleResolver.setErrorHandler(ErrorHandler)
	SetErrorHandler(errorHandler ErrorHandler)

	GetModuleResolverBase() *ModuleResolverBase
}

// ModuleResolverBase holds the fields and base implementations of
// ModuleResolver; subclasses embed it.
type ModuleResolverBase struct {
	// The set of all known input module URIs (including trailing .js), after normalization.
	modulePaths *jutil.LinkedHashSet[string]

	// Root URIs to match module roots against.
	moduleRootPaths []string

	errorHandler ErrorHandler
	pathEscaper  ModuleLoader_PathEscaper
}

// InitModuleResolverBase is the Go form of super(...).
//
// Java: ModuleResolver(ImmutableSet<String>, ImmutableList<String>, ErrorHandler, PathEscaper)
func (r *ModuleResolverBase) InitModuleResolverBase(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) {
	r.modulePaths = modulePaths
	r.moduleRootPaths = moduleRootPaths
	r.errorHandler = errorHandler
	r.pathEscaper = pathEscaper
}

// GetModuleResolverBase gives access to the base fields.
func (r *ModuleResolverBase) GetModuleResolverBase() *ModuleResolverBase { return r }

// Java: ModuleResolver.getPackageJsonMainEntries()
func (r *ModuleResolverBase) GetPackageJsonMainEntries() *jutil.LinkedHashMap[string, string] {
	return jutil.NewLinkedHashMap[string, string]()
}

// Java: ModuleResolver.resolveModuleAsPath(String, String)
func (r *ModuleResolverBase) ResolveModuleAsPath(scriptAddress string, moduleAddress string) string {
	if !strings.HasSuffix(moduleAddress, ".js") {
		moduleAddress += ".js"
	}
	path := r.pathEscaper.Escape(moduleAddress)
	if ModuleLoader_IsRelativeIdentifier(moduleAddress) {
		ourPath := scriptAddress
		lastIndex := strings.LastIndex(ourPath, ModuleLoader_MODULE_SLASH)
		path =
			ModuleNames_CanonicalizePath(
				ourPath[0:lastIndex+len(ModuleLoader_MODULE_SLASH)] + path)
	}
	return ModuleLoader_Normalize(path, r.moduleRootPaths)
}

// Locate locates the module with the given name, but returns nil if there
// is no JS file in the expected location.
//
// Java: ModuleResolver.locate(String, String)
func (r *ModuleResolverBase) Locate(scriptAddress string, name string) *string {
	canonicalizedPath := r.CanonicalizePath(scriptAddress, name)

	normalizedPath := canonicalizedPath
	if ModuleLoader_IsAmbiguousIdentifier(canonicalizedPath) {
		normalizedPath = ModuleLoader_MODULE_SLASH + canonicalizedPath
	}

	// First check to see if the module is known with it's provided path
	if r.modulePaths.Contains(normalizedPath) {
		return &canonicalizedPath
	}

	// Check for the module beneath each of the module roots
	for _, rootPath := range r.moduleRootPaths {
		modulePath := rootPath + normalizedPath

		// Since there might be code that relying on whether the path has a leading slash or not,
		// honor the state it was provided in. In an ideal world this would always be normalized
		// to contain a leading slash.
		if r.modulePaths.Contains(modulePath) {
			return &canonicalizedPath
		}
	}

	return nil
}

// CanonicalizePath normalizes a module path reference. Includes escaping
// special characters and converting relative paths to absolute references.
//
// Java: ModuleResolver.canonicalizePath(String, String)
func (r *ModuleResolverBase) CanonicalizePath(scriptAddress string, moduleAddress string) string {
	path := r.pathEscaper.Escape(moduleAddress)
	if ModuleLoader_IsRelativeIdentifier(moduleAddress) {
		ourPath := scriptAddress
		lastIndex := strings.LastIndex(ourPath, ModuleLoader_MODULE_SLASH)
		path =
			ModuleNames_CanonicalizePath(
				ourPath[0:lastIndex+len(ModuleLoader_MODULE_SLASH)] + path)
	}
	return path
}

// Java: ModuleResolver.setErrorHandler(ErrorHandler)
func (r *ModuleResolverBase) SetErrorHandler(errorHandler ErrorHandler) {
	r.errorHandler = errorHandler
}
