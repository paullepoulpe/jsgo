// Port of com/google/javascript/jscomp/deps/ModuleResolver.java @ closure-compiler v20261104
// Port of com/google/javascript/jscomp/deps/BrowserModuleResolver.java @ closure-compiler v20261104
// Port of com/google/javascript/jscomp/deps/WebpackModuleResolver.java @ closure-compiler v20261104
// Port of com/google/javascript/jscomp/deps/BrowserWithTransformedPrefixesModuleResolver.java @ closure-compiler v20261104

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
	"iter"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ModuleResolver is the base class for algorithms that resolve JavaScript
// module references to input files.
type ModuleResolver interface {
	GetPackageJsonMainEntries() *OrderedStringMap[string]

	// ResolveJsModule resolves a JS module reference, returning (path, true) if found.
	ResolveJsModule(scriptAddress string, moduleAddress string, sourcename string, lineno int, colno int) (string, bool)

	// ResolveJsModuleSilently resolves a JS module reference silently, returning (path, true) if found.
	ResolveJsModuleSilently(scriptAddress string, moduleAddress string) (string, bool)

	ResolveModuleAsPath(scriptAddress string, moduleAddress string) string

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

var emptyPackageJsonMainEntries = &OrderedStringMap[string]{}

// moduleResolverCommon holds the shared fields and path helpers for concrete
// ModuleResolver implementations without Base/Init fake-inheritance boilerplate.
type moduleResolverCommon struct {
	// The set of all known input module URIs (including trailing .js), after normalization.
	modulePaths *OrderedSet[string]

	// Root URIs to match module roots against.
	moduleRootPaths []string

	errorHandler ErrorHandler
	pathEscaper  ModuleLoader_PathEscaper
}

func newModuleResolverCommon(
	modulePaths *OrderedSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) moduleResolverCommon {
	return moduleResolverCommon{
		modulePaths:     modulePaths,
		moduleRootPaths: moduleRootPaths,
		errorHandler:    errorHandler,
		pathEscaper:     pathEscaper,
	}
}

func (r *moduleResolverCommon) GetPackageJsonMainEntries() *OrderedStringMap[string] {
	return emptyPackageJsonMainEntries
}

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
func (r *moduleResolverCommon) Locate(scriptAddress string, name string) (string, bool) {
	canonicalizedPath := r.CanonicalizePath(scriptAddress, name)

	normalizedPath := canonicalizedPath
	if ModuleLoader_IsAmbiguousIdentifier(canonicalizedPath) {
		normalizedPath = ModuleLoader_MODULE_SLASH + canonicalizedPath
	}

	// First check to see if the module is known with its provided path
	if r.modulePaths.Has(normalizedPath) {
		return canonicalizedPath, true
	}

	// Check for the module beneath each of the module roots
	for _, rootPath := range r.moduleRootPaths {
		modulePath := rootPath + normalizedPath

		// Since there might be code relying on whether the path has a leading slash or not,
		// honor the state it was provided in.
		if r.modulePaths.Has(modulePath) {
			return canonicalizedPath, true
		}
	}

	return "", false
}

// CanonicalizePath normalizes a module path reference. Includes escaping
// special characters and converting relative paths to absolute references.
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

var BrowserModuleResolver_FACTORY ModuleLoader_ModuleResolverFactory = ModuleLoader_ModuleResolverFactoryFunc(
	func(modulePaths *OrderedSet[string], moduleRootPaths []string, errorHandler ErrorHandler, pathEscaper ModuleLoader_PathEscaper) ModuleResolver {
		return NewBrowserModuleResolver(modulePaths, moduleRootPaths, errorHandler, pathEscaper)
	})

func NewBrowserModuleResolver(
	modulePaths *OrderedSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) *BrowserModuleResolver {
	return &BrowserModuleResolver{
		moduleResolverCommon: newModuleResolverCommon(modulePaths, moduleRootPaths, errorHandler, pathEscaper),
	}
}

func (r *BrowserModuleResolver) ResolveJsModule(
	scriptAddress string, moduleAddress string, sourcename string, lineno int, colno int) (string, bool) {

	if ModuleLoader_IsAmbiguousIdentifier(moduleAddress) {
		r.errorHandler.Report(
			CheckLevel_WARNING,
			JSError_MakeAtSourceLocation(
				sourcename,
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
			CheckLevel_WARNING,
			JSError_MakeAtSourceLocation(sourcename, lineno, colno, ModuleLoader_LOAD_WARNING, moduleAddress))
	}
	return loadAddress, ok
}

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

type orderedStringPairs interface {
	comparable
	Size() int
	All() iter.Seq2[string, string]
}

func cloneOrderedStringMap[M orderedStringPairs](src M) *OrderedStringMap[string] {
	var zero M
	if src == zero {
		return nil
	}
	if om, ok := any(src).(*OrderedStringMap[string]); ok {
		c := om.Clone()
		return &c
	}
	dst := NewOrderedStringMapWithCapacity[string](src.Size())
	for k, v := range src.All() {
		dst.Put(k, v)
	}
	return dst
}

// WebpackModuleResolver_Factory uses a lookup map provided by webpack to
// locate modules from a numeric id used during import.
type WebpackModuleResolver_Factory struct {
	lookupMap *OrderedStringMap[string]
}

var _ ModuleLoader_ModuleResolverFactory = (*WebpackModuleResolver_Factory)(nil)

func NewWebpackModuleResolver_Factory[M orderedStringPairs](lookupMap M) *WebpackModuleResolver_Factory {
	return &WebpackModuleResolver_Factory{lookupMap: cloneOrderedStringMap(lookupMap)}
}

func (f *WebpackModuleResolver_Factory) Create(
	modulePaths *OrderedSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) ModuleResolver {
	jutil.CheckNotNull(f.lookupMap)
	normalizedPathsById := make(map[string]string, f.lookupMap.Size())
	for _, e := range f.lookupMap.Entries() {
		canonicalizedPath := ModuleLoader_Normalize(pathEscaper.Escape(e.Value), moduleRootPaths)
		if ModuleLoader_IsAmbiguousIdentifier(canonicalizedPath) {
			canonicalizedPath = ModuleLoader_MODULE_SLASH + canonicalizedPath
		}
		normalizedPathsById[e.Key] = canonicalizedPath
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

func (r *WebpackModuleResolver) ResolveJsModule(
	scriptAddress string, moduleAddress string, sourcename string, lineno int, colno int) (string, bool) {
	if loadAddress, ok := r.modulesById[moduleAddress]; ok {
		return loadAddress, true
	}
	return r.NodeModuleResolver.ResolveJsModule(scriptAddress, moduleAddress, sourcename, lineno, colno)
}

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
	prefixReplacements *OrderedStringMap[string]
}

var _ ModuleLoader_ModuleResolverFactory = (*BrowserWithTransformedPrefixesModuleResolver_Factory)(nil)

func NewBrowserWithTransformedPrefixesModuleResolver_Factory[M orderedStringPairs](prefixReplacements M) *BrowserWithTransformedPrefixesModuleResolver_Factory {
	return &BrowserWithTransformedPrefixesModuleResolver_Factory{
		prefixReplacements: cloneOrderedStringMap(prefixReplacements),
	}
}

func (f *BrowserWithTransformedPrefixesModuleResolver_Factory) Create(
	modulePaths *OrderedSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) ModuleResolver {
	return NewBrowserWithTransformedPrefixesModuleResolver(
		modulePaths, moduleRootPaths, errorHandler, pathEscaper, f.prefixReplacements)
}

func NewBrowserWithTransformedPrefixesModuleResolver(
	modulePaths *OrderedSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper,
	prefixReplacements *OrderedStringMap[string]) *BrowserWithTransformedPrefixesModuleResolver {
	jutil.CheckNotNull(prefixReplacements)
	r := &BrowserWithTransformedPrefixesModuleResolver{
		moduleResolverCommon: newModuleResolverCommon(modulePaths, moduleRootPaths, errorHandler, pathEscaper),
	}
	entries := prefixReplacements.Entries()
	replacements := make([]browserPrefixReplacement, len(entries))
	prefixes := make([]string, len(entries))
	for i, e := range entries {
		replacements[i] = browserPrefixReplacement{prefix: e.Key, replacement: e.Value}
		prefixes[i] = e.Key
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

func (r *BrowserWithTransformedPrefixesModuleResolver) ResolveJsModule(
	scriptAddress string, moduleAddress string, sourcename string, lineno int, colno int) (string, bool) {
	transformedAddress, replaced := r.replacePathPrefixes(moduleAddress)
	//Note: in Java, `replaced = !transformedAddress.equals(moduleAddress)` (so replacing "" -> "" or "x" -> "x" has replaced == false!)
	replaced = replaced && transformedAddress != moduleAddress
	if !replaced && ModuleLoader_IsAmbiguousIdentifier(transformedAddress) {
		r.errorHandler.Report(
			CheckLevel_WARNING,
			JSError_MakeAtSourceLocation(
				sourcename,
				lineno,
				colno,
				BrowserWithTransformedPrefixesModuleResolver_INVALID_AMBIGUOUS_PATH,
				transformedAddress,
				r.expectedPrefixes))
		return "", false
	}
	return r.Locate(scriptAddress, transformedAddress)
}

func (r *BrowserWithTransformedPrefixesModuleResolver) ResolveJsModuleSilently(scriptAddress string, moduleAddress string) (string, bool) {
	transformedAddress, replaced := r.replacePathPrefixes(moduleAddress)
	replaced = replaced && transformedAddress != moduleAddress
	if !replaced && ModuleLoader_IsAmbiguousIdentifier(transformedAddress) {
		return "", false
	}
	return r.Locate(scriptAddress, transformedAddress)
}

func (r *BrowserWithTransformedPrefixesModuleResolver) ResolveModuleAsPath(scriptAddress string, moduleAddress string) string {
	if ModuleLoader_IsRelativeIdentifier(moduleAddress) {
		return r.moduleResolverCommon.ResolveModuleAsPath(scriptAddress, moduleAddress)
	}
	transformedAddress, _ := r.replacePathPrefixes(moduleAddress)
	return ModuleLoader_Normalize(transformedAddress, r.moduleRootPaths)
}
