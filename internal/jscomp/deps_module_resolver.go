// Port of com/google/javascript/jscomp/deps/ModuleResolver.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/deps/BrowserModuleResolver.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/deps/WebpackModuleResolver.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/deps/BrowserWithTransformedPrefixesModuleResolver.java @ closure-compiler v20261021

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

import (
	"cmp"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ModuleResolver is the base class for algorithms that resolve JavaScript
// module references to input files.
type ModuleResolver interface {
	// Java: ModuleResolver.getPackageJsonMainEntries()
	GetPackageJsonMainEntries() *jutil.LinkedHashMap[string, string]

	// ResolveJsModule resolves a JS module reference, returning (path, true) if found.
	//
	// Java: ModuleResolver.resolveJsModule(String, String, String, int, int)
	ResolveJsModule(scriptAddress string, moduleAddress string, sourcename string, lineno int, colno int) (string, bool)

	// ResolveJsModuleSilently resolves a JS module reference silently, returning (path, true) if found.
	//
	// Java: ModuleResolver.resolveJsModuleSilently(String, String)
	ResolveJsModuleSilently(scriptAddress string, moduleAddress string) (string, bool)

	// Java: ModuleResolver.resolveModuleAsPath(String, String)
	ResolveModuleAsPath(scriptAddress string, moduleAddress string) string

	// Java: ModuleResolver.setErrorHandler(ErrorHandler)
	SetErrorHandler(errorHandler ErrorHandler)
}

func isASCIIPath(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

func comparePathsLongestFirst(a, b string) int {
	if isASCIIPath(a) && isASCIIPath(b) {
		if len(a) != len(b) {
			return cmp.Compare(len(b), len(a))
		}
		return strings.Compare(a, b)
	}
	if c := cmp.Compare(jutil.StringLength(b), jutil.StringLength(a)); c != 0 {
		return c
	}
	return jutil.StringCompareTo(a, b)
}

func comparePathCodeUnits(a, b string) int {
	if isASCIIPath(a) && isASCIIPath(b) {
		return strings.Compare(a, b)
	}
	return jutil.StringCompareTo(a, b)
}

var emptyPackageJsonMainEntries = jutil.NewLinkedHashMap[string, string]()

// moduleResolverCommon holds the shared fields and path helpers for concrete
// ModuleResolver implementations without Base/Init fake-inheritance boilerplate.
type moduleResolverCommon struct {
	// The set of all known input module URIs (including trailing .js), after normalization.
	modulePaths map[string]struct{}

	// Root URIs to match module roots against.
	moduleRootPaths []string

	errorHandler ErrorHandler
	pathEscaper  ModuleLoader_PathEscaper
}

func newModuleResolverCommon(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) moduleResolverCommon {
	var mp map[string]struct{}
	if modulePaths != nil && modulePaths.Size() > 0 {
		mp = make(map[string]struct{}, modulePaths.Size())
		for p := range modulePaths.All() {
			mp[p] = struct{}{}
		}
	}
	return moduleResolverCommon{
		modulePaths:     mp,
		moduleRootPaths: moduleRootPaths,
		errorHandler:    errorHandler,
		pathEscaper:     pathEscaper,
	}
}

// Java: ModuleResolver.getPackageJsonMainEntries()
func (r *moduleResolverCommon) GetPackageJsonMainEntries() *jutil.LinkedHashMap[string, string] {
	return emptyPackageJsonMainEntries
}

// Java: ModuleResolver.resolveModuleAsPath(String, String)
func (r *moduleResolverCommon) ResolveModuleAsPath(scriptAddress string, moduleAddress string) string {
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

// Locate locates the module with the given name, returning ("", false) if there
// is no JS file in the expected location.
//
// Java: ModuleResolver.locate(String, String)
func (r *moduleResolverCommon) Locate(scriptAddress string, name string) (string, bool) {
	canonicalizedPath := r.CanonicalizePath(scriptAddress, name)

	normalizedPath := canonicalizedPath
	if ModuleLoader_IsAmbiguousIdentifier(canonicalizedPath) {
		normalizedPath = ModuleLoader_MODULE_SLASH + canonicalizedPath
	}

	// First check to see if the module is known with its provided path
	if _, ok := r.modulePaths[normalizedPath]; ok {
		return canonicalizedPath, true
	}

	// Check for the module beneath each of the module roots
	for _, rootPath := range r.moduleRootPaths {
		modulePath := rootPath + normalizedPath

		// Since there might be code relying on whether the path has a leading slash or not,
		// honor the state it was provided in.
		if _, ok := r.modulePaths[modulePath]; ok {
			return canonicalizedPath, true
		}
	}

	return "", false
}

// CanonicalizePath normalizes a module path reference. Includes escaping
// special characters and converting relative paths to absolute references.
//
// Java: ModuleResolver.canonicalizePath(String, String)
func (r *moduleResolverCommon) CanonicalizePath(scriptAddress string, moduleAddress string) string {
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
func (r *moduleResolverCommon) SetErrorHandler(errorHandler ErrorHandler) {
	r.errorHandler = errorHandler
}

// BrowserModuleResolver is the resolution algorithm for Browsers.
//
// Only unambiguous paths are supported. Paths must specify a file
// extension.
type BrowserModuleResolver struct {
	moduleResolverCommon
}

var _ ModuleResolver = (*BrowserModuleResolver)(nil)

// Java: BrowserModuleResolver.FACTORY
var BrowserModuleResolver_FACTORY ModuleLoader_ModuleResolverFactory = ModuleLoader_ModuleResolverFactoryFunc(
	func(modulePaths *jutil.LinkedHashSet[string], moduleRootPaths []string, errorHandler ErrorHandler, pathEscaper ModuleLoader_PathEscaper) ModuleResolver {
		return NewBrowserModuleResolver(modulePaths, moduleRootPaths, errorHandler, pathEscaper)
	})

// Java: BrowserModuleResolver(ImmutableSet<String>, ImmutableList<String>, ErrorHandler, PathEscaper)
func NewBrowserModuleResolver(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) *BrowserModuleResolver {
	return &BrowserModuleResolver{
		moduleResolverCommon: newModuleResolverCommon(modulePaths, moduleRootPaths, errorHandler, pathEscaper),
	}
}

// Java: BrowserModuleResolver.resolveJsModule(String, String, String, int, int)
func (r *BrowserModuleResolver) ResolveJsModule(
	scriptAddress string, moduleAddress string, sourcename string, lineno int, colno int) (string, bool) {

	if ModuleLoader_IsAmbiguousIdentifier(moduleAddress) {
		r.errorHandler.Report(
			CheckLevel_WARNING.Ptr(),
			JSError_MakeWithSourceLocation(
				&sourcename,
				lineno,
				colno,
				ModuleLoader_INVALID_MODULE_PATH,
				moduleAddress,
				ModuleLoader_ResolutionMode_BROWSER.String()))
		return "", false
	}

	loadAddress, ok := r.Locate(scriptAddress, moduleAddress)
	if !ok {
		r.errorHandler.Report(
			CheckLevel_WARNING.Ptr(),
			JSError_MakeWithSourceLocation(&sourcename, lineno, colno, ModuleLoader_LOAD_WARNING, moduleAddress))
	}
	return loadAddress, ok
}

// Java: BrowserModuleResolver.resolveJsModuleSilently(String, String)
func (r *BrowserModuleResolver) ResolveJsModuleSilently(scriptAddress string, moduleAddress string) (string, bool) {
	if ModuleLoader_IsAmbiguousIdentifier(moduleAddress) {
		return "", false
	}
	return r.Locate(scriptAddress, moduleAddress)
}

// WebpackModuleResolver is the resolution algorithm for Webpack. Modules are
// located by a map of numeric ids to module paths.
type WebpackModuleResolver struct {
	NodeModuleResolver
	modulesById map[string]string
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
	normalizedPathsById := make(map[string]string, f.lookupMap.Size())
	for key, value := range f.lookupMap.All() {
		canonicalizedPath := ModuleLoader_Normalize(pathEscaper.Escape(value), moduleRootPaths)
		if ModuleLoader_IsAmbiguousIdentifier(canonicalizedPath) {
			canonicalizedPath = ModuleLoader_MODULE_SLASH + canonicalizedPath
		}
		normalizedPathsById[key] = canonicalizedPath
	}
	return &WebpackModuleResolver{
		NodeModuleResolver: *NewNodeModuleResolver(
			modulePaths,
			moduleRootPaths,
			nil,
			errorHandler,
			pathEscaper),
		modulesById: normalizedPathsById,
	}
}

// Java: WebpackModuleResolver.resolveJsModule(String, String, String, int, int)
func (r *WebpackModuleResolver) ResolveJsModule(
	scriptAddress string, moduleAddress string, sourcename string, lineno int, colno int) (string, bool) {
	if loadAddress, ok := r.modulesById[moduleAddress]; ok {
		return loadAddress, true
	}
	return r.NodeModuleResolver.ResolveJsModule(scriptAddress, moduleAddress, sourcename, lineno, colno)
}

// Java: WebpackModuleResolver.resolveJsModuleSilently(String, String)
func (r *WebpackModuleResolver) ResolveJsModuleSilently(scriptAddress string, moduleAddress string) (string, bool) {
	if loadAddress, ok := r.modulesById[moduleAddress]; ok {
		return loadAddress, true
	}
	return r.NodeModuleResolver.ResolveJsModuleSilently(scriptAddress, moduleAddress)
}

type browserPrefixReplacement struct {
	prefix      string
	replacement string
}

// BrowserWithTransformedPrefixesModuleResolver is a limited superset of the
// BrowserModuleResolver that allows for replacing some path prefixes before
// resolving.
type BrowserWithTransformedPrefixesModuleResolver struct {
	moduleResolverCommon
	prefixReplacements []browserPrefixReplacement
	expectedPrefixes   string
}

var _ ModuleResolver = (*BrowserWithTransformedPrefixesModuleResolver)(nil)

var BrowserWithTransformedPrefixesModuleResolver_INVALID_AMBIGUOUS_PATH = DiagnosticType_Error(
	"JSC_INVALID_AMBIGUOUS_PATH",
	"The path \"{0}\" is invalid. Expected any of the following prefixes for non-relative "+
		"paths: {1}.")

// BrowserWithTransformedPrefixesModuleResolver_Factory is the factory for
// BrowserWithTransformedPrefixesModuleResolver.
type BrowserWithTransformedPrefixesModuleResolver_Factory struct {
	prefixReplacements *jutil.LinkedHashMap[string, string]
}

var _ ModuleLoader_ModuleResolverFactory = (*BrowserWithTransformedPrefixesModuleResolver_Factory)(nil)

// Java: BrowserWithTransformedPrefixesModuleResolver.Factory(ImmutableMap<String, String>)
func NewBrowserWithTransformedPrefixesModuleResolver_Factory(prefixReplacements *jutil.LinkedHashMap[string, string]) *BrowserWithTransformedPrefixesModuleResolver_Factory {
	return &BrowserWithTransformedPrefixesModuleResolver_Factory{prefixReplacements: prefixReplacements}
}

// Java: BrowserWithTransformedPrefixesModuleResolver.Factory.create(ImmutableSet<String>, ImmutableList<String>, ErrorHandler, PathEscaper)
func (f *BrowserWithTransformedPrefixesModuleResolver_Factory) Create(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) ModuleResolver {
	return NewBrowserWithTransformedPrefixesModuleResolver(
		modulePaths, moduleRootPaths, errorHandler, pathEscaper, f.prefixReplacements)
}

// Java: BrowserWithTransformedPrefixesModuleResolver(ImmutableSet<String>, ImmutableList<String>, ErrorHandler, PathEscaper, ImmutableMap<String, String>)
func NewBrowserWithTransformedPrefixesModuleResolver(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper,
	prefixReplacements *jutil.LinkedHashMap[string, string]) *BrowserWithTransformedPrefixesModuleResolver {
	r := &BrowserWithTransformedPrefixesModuleResolver{
		moduleResolverCommon: newModuleResolverCommon(modulePaths, moduleRootPaths, errorHandler, pathEscaper),
	}
	n := prefixReplacements.Size()
	replacements := make([]browserPrefixReplacement, 0, n)
	prefixes := make([]string, 0, n)
	for k, v := range prefixReplacements.All() {
		replacements = append(replacements, browserPrefixReplacement{prefix: k, replacement: v})
		prefixes = append(prefixes, k)
	}
	slices.SortFunc(replacements, func(a, b browserPrefixReplacement) int {
		return comparePathsLongestFirst(a.prefix, b.prefix)
	})
	r.prefixReplacements = replacements
	slices.SortFunc(prefixes, comparePathCodeUnits)
	r.expectedPrefixes = strings.Join(prefixes, ", ")
	return r
}

func (r *BrowserWithTransformedPrefixesModuleResolver) replacePathPrefixes(moduleAddress string) (string, bool) {
	for _, pr := range r.prefixReplacements {
		if strings.HasPrefix(moduleAddress, pr.prefix) {
			return pr.replacement + moduleAddress[len(pr.prefix):], true
		}
	}
	return moduleAddress, false
}

// Java: BrowserWithTransformedPrefixesModuleResolver.resolveJsModule(String, String, String, int, int)
func (r *BrowserWithTransformedPrefixesModuleResolver) ResolveJsModule(
	scriptAddress string, moduleAddress string, sourcename string, lineno int, colno int) (string, bool) {
	transformedAddress, replaced := r.replacePathPrefixes(moduleAddress)
	//Note: in Java, `replaced = !transformedAddress.equals(moduleAddress)` (so replacing "" -> "" or "x" -> "x" has replaced == false!)
	replaced = replaced && transformedAddress != moduleAddress
	if !replaced && ModuleLoader_IsAmbiguousIdentifier(transformedAddress) {
		r.errorHandler.Report(
			CheckLevel_WARNING.Ptr(),
			JSError_MakeWithSourceLocation(
				&sourcename,
				lineno,
				colno,
				BrowserWithTransformedPrefixesModuleResolver_INVALID_AMBIGUOUS_PATH,
				transformedAddress,
				r.expectedPrefixes))
		return "", false
	}
	return r.Locate(scriptAddress, transformedAddress)
}

// Java: BrowserWithTransformedPrefixesModuleResolver.resolveJsModuleSilently(String, String)
func (r *BrowserWithTransformedPrefixesModuleResolver) ResolveJsModuleSilently(scriptAddress string, moduleAddress string) (string, bool) {
	transformedAddress, replaced := r.replacePathPrefixes(moduleAddress)
	replaced = replaced && transformedAddress != moduleAddress
	if !replaced && ModuleLoader_IsAmbiguousIdentifier(transformedAddress) {
		return "", false
	}
	return r.Locate(scriptAddress, transformedAddress)
}

// Java: BrowserWithTransformedPrefixesModuleResolver.resolveModuleAsPath(String, String)
func (r *BrowserWithTransformedPrefixesModuleResolver) ResolveModuleAsPath(scriptAddress string, moduleAddress string) string {
	if ModuleLoader_IsRelativeIdentifier(moduleAddress) {
		return r.moduleResolverCommon.ResolveModuleAsPath(scriptAddress, moduleAddress)
	}
	transformedAddress, _ := r.replacePathPrefixes(moduleAddress)
	return ModuleLoader_Normalize(transformedAddress, r.moduleRootPaths)
}
