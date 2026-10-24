// Port of com/google/javascript/jscomp/deps/ModuleLoader.java @ closure-compiler v20261015

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

package jscomp

import (
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ModuleLoader provides compile-time locate semantics for ES6 and CommonJS
// modules.
//
// See https://tc39.github.io/ecma262/#sec-module-semantics and
// http://wiki.commonjs.org/wiki/Modules/1.1.
type ModuleLoader struct {
	errorHandler ErrorHandler

	// Root URIs to match module roots against.
	moduleRootPaths []string
	// The set of all known input module URIs (including trailing .js), after normalization.
	// Used to canonicalize paths before resolution.
	pathResolver ModuleLoader_PathResolver

	pathEscaper ModuleLoader_PathEscaper

	moduleResolver ModuleResolver
}

var ModuleLoader_MODULE_CONFLICT = DiagnosticType_Warning(
	"JSC_MODULE_CONFLICT",
	"File cannot be a combination of goog.provide, goog.module, and/or ES6 module: {0}")

// ModuleLoader_MODULE_SLASH: according to the spec, the forward slash should
// be the delimiter on all platforms.
const ModuleLoader_MODULE_SLASH = ModuleNames_MODULE_SLASH

// ModuleLoader_DEFAULT_FILENAME_PREFIX is the default module root, the
// current directory.
const ModuleLoader_DEFAULT_FILENAME_PREFIX = "." + ModuleLoader_MODULE_SLASH

const ModuleLoader_JSC_BROWSER_SKIPLISTED_MARKER = "$jscomp$browser$skiplisted"

var ModuleLoader_LOAD_WARNING = DiagnosticType_Error("JSC_JS_MODULE_LOAD_WARNING", "Failed to load module \"{0}\"")

var ModuleLoader_INVALID_MODULE_PATH = DiagnosticType_Error(
	"JSC_INVALID_MODULE_PATH", "Invalid module path \"{0}\" for resolution mode \"{1}\"")

// Java: ModuleLoader.builder()
func ModuleLoader_Builder_() *ModuleLoader_Builder {
	return newModuleLoader_Builder()
}

// ModuleLoader_Builder is the Builder.
//
// Go-only: language-adapter — moduleRootsSet/inputsSet record whether the
// Java Iterables were set (non-null), since a nil Go slice is also empty.
type ModuleLoader_Builder struct {
	errorHandler   ErrorHandler
	moduleRoots    []string
	moduleRootsSet bool
	inputs         []DependencyInfo
	inputsSet      bool
	factory        ModuleLoader_ModuleResolverFactory
	pathResolver   ModuleLoader_PathResolver
	pathEscaper    ModuleLoader_PathEscaper
}

// Java: ModuleLoader.Builder()
func newModuleLoader_Builder() *ModuleLoader_Builder {
	return &ModuleLoader_Builder{
		errorHandler: moduleLoader_NOOP_ERROR_HANDER,
		pathResolver: ModuleLoader_PathResolver_RELATIVE,
		pathEscaper:  ModuleLoader_PathEscaper_ESCAPE,
	}
}

// Java: ModuleLoader.Builder.setErrorHandler(ErrorHandler)
func (b *ModuleLoader_Builder) SetErrorHandler(x ErrorHandler) *ModuleLoader_Builder {
	if x != nil {
		b.errorHandler = x
	}
	return b
}

// SetModuleRoots sets the path prefixes to strip from module paths.
//
// Java: ModuleLoader.Builder.setModuleRoots(Iterable<String>)
func (b *ModuleLoader_Builder) SetModuleRoots(x []string) *ModuleLoader_Builder {
	b.moduleRoots = x
	b.moduleRootsSet = true
	return b
}

// SetInputs sets all inputs to the compilation process.
//
// Used to ensure that resolved paths references a valid input.
//
// Java: ModuleLoader.Builder.setInputs(Iterable<? extends DependencyInfo>)
func (b *ModuleLoader_Builder) SetInputs(x []DependencyInfo) *ModuleLoader_Builder {
	b.inputs = x
	b.inputsSet = true
	return b
}

// SetFactory creates a module resolver, which determines how module
// identifiers are resolved.
//
// Java: ModuleLoader.Builder.setFactory(ModuleResolverFactory)
func (b *ModuleLoader_Builder) SetFactory(x ModuleLoader_ModuleResolverFactory) *ModuleLoader_Builder {
	b.factory = x
	return b
}

// SetPathResolver determines how to sanitize paths before resolving.
//
// Java: ModuleLoader.Builder.setPathResolver(PathResolver)
func (b *ModuleLoader_Builder) SetPathResolver(x ModuleLoader_PathResolver) *ModuleLoader_Builder {
	b.pathResolver = x
	return b
}

// SetPathEscaper determines if / how paths should be escaped.
//
// Java: ModuleLoader.Builder.setPathEscaper(PathEscaper)
func (b *ModuleLoader_Builder) SetPathEscaper(x ModuleLoader_PathEscaper) *ModuleLoader_Builder {
	b.pathEscaper = x
	return b
}

// Java: ModuleLoader.Builder.build()
func (b *ModuleLoader_Builder) Build() *ModuleLoader {
	return newModuleLoader(b)
}

// Java: ModuleLoader(Builder)
func newModuleLoader(builder *ModuleLoader_Builder) *ModuleLoader {
	l := &ModuleLoader{}
	l.errorHandler = jutil.CheckNotNull(builder.errorHandler)
	l.pathResolver = builder.pathResolver
	l.pathEscaper = builder.pathEscaper
	if !builder.moduleRootsSet {
		panic(jutil.NewNullPointerException())
	}
	l.moduleRootPaths = l.CreateRootPaths(builder.moduleRoots)

	if !builder.inputsSet {
		panic(jutil.NewNullPointerException())
	}
	modulePaths := l.ResolvePaths(builder.inputs)
	l.moduleResolver =
		builder.factory.Create(
			modulePaths, l.moduleRootPaths, l.errorHandler, l.pathEscaper)
	return l
}

// Java: ModuleLoader.getPackageJsonMainEntries()
func (l *ModuleLoader) GetPackageJsonMainEntries() *jutil.LinkedHashMap[string, string] {
	return l.moduleResolver.GetPackageJsonMainEntries()
}

// ModuleLoader_ModulePath is a path to a module. Provides access to the
// module's closurized name and a way to resolve relative paths.
type ModuleLoader_ModulePath struct {
	outer *ModuleLoader
	path  string
}

// Java: ModuleLoader.ModulePath(String)
func (l *ModuleLoader) newModulePath(path string) *ModuleLoader_ModulePath {
	return &ModuleLoader_ModulePath{outer: l, path: path}
}

// Java: ModuleLoader.ModulePath.toString()
func (p *ModuleLoader_ModulePath) String() string {
	return p.path
}

// EqualsIgnoreLeadingSlash determines if this path is the same as another
// path, ignoring any potential leading slashes on both.
//
// Java: ModuleLoader.ModulePath.equalsIgnoreLeadingSlash(ModulePath)
func (p *ModuleLoader_ModulePath) EqualsIgnoreLeadingSlash(other *ModuleLoader_ModulePath) bool {
	return other != nil && p.ToModuleName() == other.ToModuleName()
}

// ToModuleName turns a filename into a JS identifier that is used for
// moduleNames in rewritten code. Removes leading /, replaces / with $,
// removes trailing .js and replaces - with _. All moduleNames get a
// "module$" prefix.
//
// Java: ModuleLoader.ModulePath.toModuleName()
func (p *ModuleLoader_ModulePath) ToModuleName() string {
	return ModuleNames_ToModuleName(p.path)
}

// ResolveJsModule finds a JS module requireName. See
// https://nodejs.org/api/modules.html#modules_all_together
//
// It returns the normalized module URI, or nil if not found.
//
// Java: ModuleLoader.ModulePath.resolveJsModule(String, String, int, int)
func (p *ModuleLoader_ModulePath) ResolveJsModule(
	moduleAddress string, sourcename string, lineno int, colno int) *ModuleLoader_ModulePath {
	loadAddress :=
		p.outer.moduleResolver.ResolveJsModule(p.path, moduleAddress, sourcename, lineno, colno)

	if loadAddress != nil {
		return p.outer.newModulePath(*loadAddress)
	}

	return nil
}

// ResolveJsModuleSilently resolves a JS module for dynamic imports without
// reporting errors for unresolved addresses.
//
// It returns the resolved module path, or nil if unresolvable (or
// ambiguous).
//
// Java: ModuleLoader.ModulePath.resolveJsModuleSilently(String)
func (p *ModuleLoader_ModulePath) ResolveJsModuleSilently(moduleAddress string) *ModuleLoader_ModulePath {
	loadAddress := p.outer.moduleResolver.ResolveJsModuleSilently(p.path, moduleAddress)

	if loadAddress != nil {
		return p.outer.newModulePath(*loadAddress)
	}

	return nil
}

// ResolveModuleAsPath treats the module address as a path and returns the
// name of that module. Does not verify that there is actually a JS file at
// the provided URI.
//
// Primarily used for per-file ES6 module transpilation.
//
// Java: ModuleLoader.ModulePath.resolveModuleAsPath(String)
func (p *ModuleLoader_ModulePath) ResolveModuleAsPath(moduleAddress string) *ModuleLoader_ModulePath {
	return p.outer.newModulePath(p.outer.moduleResolver.ResolveModuleAsPath(p.path, moduleAddress))
}

// Resolve resolves a path into a ModulePath.
//
// Java: ModuleLoader.resolve(String)
func (l *ModuleLoader) Resolve(path string) *ModuleLoader_ModulePath {
	return l.newModulePath(ModuleLoader_Normalize(l.pathEscaper.Escape(l.pathResolver.Apply(path)), l.moduleRootPaths))
}

// ModuleLoader_IsRelativeIdentifier: whether this is relative to the
// current file, or a top-level identifier.
//
// Java: ModuleLoader.isRelativeIdentifier(String)
func ModuleLoader_IsRelativeIdentifier(name string) bool {
	return strings.HasPrefix(name, "."+ModuleLoader_MODULE_SLASH) ||
		strings.HasPrefix(name, ".."+ModuleLoader_MODULE_SLASH) ||
		name == "." ||
		name == ".."
}

// ModuleLoader_IsAbsoluteIdentifier: whether this is absolute to the
// compilation.
//
// Java: ModuleLoader.isAbsoluteIdentifier(String)
func ModuleLoader_IsAbsoluteIdentifier(name string) bool {
	return strings.HasPrefix(name, ModuleLoader_MODULE_SLASH)
}

// ModuleLoader_IsAmbiguousIdentifier: whether this is neither absolute or
// relative.
//
// Java: ModuleLoader.isAmbiguousIdentifier(String)
func ModuleLoader_IsAmbiguousIdentifier(name string) bool {
	return !ModuleLoader_IsAbsoluteIdentifier(name) && !ModuleLoader_IsRelativeIdentifier(name)
}

// ModuleLoader_IsPathIdentifier: whether name is a path-based identifier
// (has a '/' character).
//
// Java: ModuleLoader.isPathIdentifier(String)
func ModuleLoader_IsPathIdentifier(name string) bool {
	return strings.Contains(name, ModuleLoader_MODULE_SLASH)
}

// CreateRootPaths normalizes the given root paths, which are path prefixes
// to be removed from a module path when resolved.
//
// Java: ModuleLoader.createRootPaths(Iterable<String>)
func (l *ModuleLoader) CreateRootPaths(roots []string) []string {
	var paths []string
	for _, root := range roots {
		n := l.pathEscaper.Escape(l.pathResolver.Apply(root))
		if ModuleLoader_IsAmbiguousIdentifier(n) {
			n = ModuleLoader_MODULE_SLASH + n
		}
		paths = append(paths, n)
	}
	// toImmutableSortedSet(BEST_MATCH_PATH_ORDERING).asList(): sorted, without duplicates.
	slices.SortStableFunc(paths, moduleLoader_BEST_MATCH_PATH_ORDERING)
	return slices.CompactFunc(paths, func(a, b string) bool { return moduleLoader_BEST_MATCH_PATH_ORDERING(a, b) == 0 })
}

// ResolvePaths returns the list of normalized modules, which always have a
// leading slash.
//
// Java: ModuleLoader.resolvePaths(Iterable<? extends DependencyInfo>)
func (l *ModuleLoader) ResolvePaths(inputs []DependencyInfo) *jutil.LinkedHashSet[string] {
	var paths []string
	for _, input := range inputs {
		p := l.pathEscaper.Escape(l.pathResolver.Apply(input.GetName()))
		n := ModuleLoader_Normalize(p, l.moduleRootPaths)
		if ModuleLoader_IsAmbiguousIdentifier(n) {
			n = ModuleLoader_MODULE_SLASH + n
		}
		paths = append(paths, n)
	}
	slices.SortStableFunc(paths, moduleLoader_BEST_MATCH_PATH_ORDERING) // Sort so that this errors are easier to read.
	dupeModulePaths := jutil.ImmutableMultiset_Of(paths...)

	jutil.CheckStateMsg1(
		dupeModulePaths.Size() == dupeModulePaths.ElementSet().Size(),
		"Duplicate module paths after resolving: %s",
		dupeModulePaths)

	return dupeModulePaths.ElementSet()
}

// ModuleLoader_Normalize normalizes the name and resolves it against the
// module roots.
//
// Java: ModuleLoader.normalize(String, Iterable<String>)
func ModuleLoader_Normalize(path string, moduleRootPaths []string) string {
	normalizedPath := path
	if ModuleLoader_IsAmbiguousIdentifier(normalizedPath) {
		normalizedPath = ModuleLoader_MODULE_SLASH + normalizedPath
	}

	// Find a moduleRoot that this URI is under. If none, use as is.
	for _, moduleRoot := range moduleRootPaths {
		if strings.HasPrefix(normalizedPath, moduleRoot) {
			// Make sure that e.g. path "foobar/test.js" is not matched by module "foo", by checking for
			// a leading slash.
			trailing := normalizedPath[len(moduleRoot):]
			if strings.HasPrefix(trailing, ModuleLoader_MODULE_SLASH) {
				return trailing[len(ModuleLoader_MODULE_SLASH):]
			}
		}
	}
	// Not underneath any of the roots.
	return path
}

// Java: ModuleLoader.relativePathFrom(String, String)
func ModuleLoader_RelativePathFrom(fromUriPath string, toUriPath string) string {
	fromPath := jutil.Path_Of(fromUriPath)
	toPath := jutil.Path_Of(toUriPath)
	fromFolder := fromPath.GetParent()

	// if the from URIs are simple names without paths, they are in the same folder
	// example: m0.js
	if fromFolder == nil && toPath.GetParent() == nil {
		return "./" + toUriPath
	} else if fromFolder == nil &&
		(strings.HasPrefix(toUriPath, ".") ||
			strings.HasPrefix(toPath.String(), "/") ||
			strings.HasPrefix(toPath.String(), "\\")) {
		return toUriPath
	} else if fromFolder == nil {
		panic(jutil.NewIllegalArgumentExceptionMsg("Relative path between URIs cannot be calculated"))
	}

	calculatedPath := fromFolder.Relativize(toPath).String()
	if strings.HasPrefix(calculatedPath, ".") || strings.HasPrefix(calculatedPath, "/") {
		return calculatedPath
	}
	return "./" + calculatedPath
}

// Java: ModuleLoader.setErrorHandler(ErrorHandler)
func (l *ModuleLoader) SetErrorHandler(errorHandler ErrorHandler) {
	if errorHandler == nil {
		l.errorHandler = moduleLoader_NOOP_ERROR_HANDER
	} else {
		l.errorHandler = errorHandler
	}
	l.moduleResolver.SetErrorHandler(l.errorHandler)
}

// Java: ModuleLoader.getErrorHandler()
func (l *ModuleLoader) GetErrorHandler() ErrorHandler {
	return l.errorHandler
}

// ModuleLoader_PathEscaper indicates whether to escape characters in paths.
//
// Java: ModuleLoader.PathEscaper
type ModuleLoader_PathEscaper int32

const (
	// ModuleLoader_PathEscaper_ESCAPE escapes characters in paths according
	// to ModuleNames.escapePath(String) and then canonicalizes it.
	ModuleLoader_PathEscaper_ESCAPE ModuleLoader_PathEscaper = iota

	// ModuleLoader_PathEscaper_CANONICALIZE_ONLY does not escaped characters
	// in paths, but does canonicalize it according to
	// ModuleNames.canonicalizePath(String).
	ModuleLoader_PathEscaper_CANONICALIZE_ONLY
)

var moduleLoader_PathEscaper_names = [...]string{"ESCAPE", "CANONICALIZE_ONLY"}

// Java: ModuleLoader.PathEscaper.values()
func ModuleLoader_PathEscaper_values() []ModuleLoader_PathEscaper {
	return []ModuleLoader_PathEscaper{ModuleLoader_PathEscaper_ESCAPE, ModuleLoader_PathEscaper_CANONICALIZE_ONLY}
}

// String is Enum.name().
func (e ModuleLoader_PathEscaper) String() string { return moduleLoader_PathEscaper_names[e] }

// Java: ModuleLoader.PathEscaper.escape(String)
func (e ModuleLoader_PathEscaper) Escape(path string) string {
	switch e {
	case ModuleLoader_PathEscaper_ESCAPE:
		return ModuleNames_EscapePath(path)
	default: // CANONICALIZE_ONLY
		return ModuleNames_CanonicalizePath(path)
	}
}

// ModuleLoader_PathResolver is an enum indicating whether to absolutize
// paths.
//
// Java: ModuleLoader.PathResolver
type ModuleLoader_PathResolver int32

const (
	ModuleLoader_PathResolver_RELATIVE ModuleLoader_PathResolver = iota
	ModuleLoader_PathResolver_ABSOLUTE
)

var moduleLoader_PathResolver_names = [...]string{"RELATIVE", "ABSOLUTE"}

// Java: ModuleLoader.PathResolver.values()
func ModuleLoader_PathResolver_values() []ModuleLoader_PathResolver {
	return []ModuleLoader_PathResolver{ModuleLoader_PathResolver_RELATIVE, ModuleLoader_PathResolver_ABSOLUTE}
}

// String is Enum.name().
func (e ModuleLoader_PathResolver) String() string { return moduleLoader_PathResolver_names[e] }

// Java: ModuleLoader.PathResolver.apply(String)
func (e ModuleLoader_PathResolver) Apply(path string) string {
	switch e {
	case ModuleLoader_PathResolver_RELATIVE:
		return path
	default: // ABSOLUTE
		return jutil.Path_Of(path).ToAbsolutePath().String()
	}
}

// ModuleLoader_ModuleResolverFactory is used to specify what algorithm to
// use to locate non path-based modules.
//
// Java: ModuleLoader.ModuleResolverFactory
type ModuleLoader_ModuleResolverFactory interface {
	// Java: ModuleLoader.ModuleResolverFactory.create(ImmutableSet<String>, ImmutableList<String>, ErrorHandler, PathEscaper)
	Create(
		modulePaths *jutil.LinkedHashSet[string],
		moduleRootPaths []string,
		errorHandler ErrorHandler,
		pathEscaper ModuleLoader_PathEscaper) ModuleResolver
}

// ModuleLoader_ModuleResolverFactoryFunc adapts a function (a Java lambda or
// method reference) to ModuleResolverFactory.
//
// Go-only: language-adapter — @FunctionalInterface.
type ModuleLoader_ModuleResolverFactoryFunc func(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) ModuleResolver

// Create calls f.
func (f ModuleLoader_ModuleResolverFactoryFunc) Create(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) ModuleResolver {
	return f(modulePaths, moduleRootPaths, errorHandler, pathEscaper)
}

// Sort longest length to shortest so that paths are applied most specific to least.
//
// Java: Comparator.comparingInt(String::length).reversed().thenComparing(naturalOrder())
func moduleLoader_BEST_MATCH_PATH_ORDERING(a, b string) int {
	la, lb := jutil.StringLength(a), jutil.StringLength(b)
	if la != lb {
		if lb < la {
			return -1
		}
		return 1
	}
	return jutil.StringCompareTo(a, b)
}

// moduleLoader_noopErrorHandler is the lambda (CheckLevel level, JSError error) -> {}.
type moduleLoader_noopErrorHandler struct{}

func (moduleLoader_noopErrorHandler) Report(level *CheckLevel, error *JSError) {}

var moduleLoader_NOOP_ERROR_HANDER ErrorHandler = moduleLoader_noopErrorHandler{}

// ModuleLoader_EMPTY is a trivial module loader with no roots.
var ModuleLoader_EMPTY = moduleLoader_NewEmpty()

// moduleLoader_NewEmpty returns a new trivial module loader with no roots.
//
// Go-only: a Compiler uses its own instead of the shared ModuleLoader_EMPTY
// (Java's Compiler uses ModuleLoader.EMPTY): the Compiler sets its error
// handler on its loader, and on a shared one that would send a later
// compiler's module errors to an earlier one.
func moduleLoader_NewEmpty() *ModuleLoader {
	return ModuleLoader_Builder_().
		SetModuleRoots([]string{}).
		SetInputs([]DependencyInfo{}).
		SetFactory(BrowserModuleResolver_FACTORY).
		Build()
}

// ModuleLoader_ResolutionMode lists the standard path base resolution
// algorithms that are accepted as a command line flag.
//
// Java: ModuleLoader.ResolutionMode
type ModuleLoader_ResolutionMode int32

const (
	// ModuleLoader_ResolutionMode_BROWSER mimics the behavior of MS Edge.
	//
	// Modules must begin with a "." or "/" character. Modules must include
	// the file extension MS Edge was the only browser to define a module
	// resolution behavior at the time of this writing.
	ModuleLoader_ResolutionMode_BROWSER ModuleLoader_ResolutionMode = iota

	// ModuleLoader_ResolutionMode_BROWSER_WITH_TRANSFORMED_PREFIXES is a
	// limited superset of BROWSER that transforms some path prefixes.
	//
	// For example one could configure this so that "@root/" is replaced with
	// "/my/path/to/project/" within import paths.
	ModuleLoader_ResolutionMode_BROWSER_WITH_TRANSFORMED_PREFIXES

	// ModuleLoader_ResolutionMode_NODE uses the node module resolution
	// algorithm.
	//
	// Modules which do not begin with a "." or "/" character are looked up
	// from the appropriate node_modules folder. Includes the ability to
	// require directories and JSON files. Exact match, then ".js", then
	// ".json" file extensions are searched.
	ModuleLoader_ResolutionMode_NODE

	// ModuleLoader_ResolutionMode_WEBPACK uses a lookup map provided by
	// webpack to locate modules from a numeric id used during import.
	ModuleLoader_ResolutionMode_WEBPACK
)

var moduleLoader_ResolutionMode_names = [...]string{"BROWSER", "BROWSER_WITH_TRANSFORMED_PREFIXES", "NODE", "WEBPACK"}

// Java: ModuleLoader.ResolutionMode.values()
func ModuleLoader_ResolutionMode_values() []ModuleLoader_ResolutionMode {
	return []ModuleLoader_ResolutionMode{
		ModuleLoader_ResolutionMode_BROWSER, ModuleLoader_ResolutionMode_BROWSER_WITH_TRANSFORMED_PREFIXES,
		ModuleLoader_ResolutionMode_NODE, ModuleLoader_ResolutionMode_WEBPACK,
	}
}

// String is Enum.name().
func (e ModuleLoader_ResolutionMode) String() string { return moduleLoader_ResolutionMode_names[e] }
