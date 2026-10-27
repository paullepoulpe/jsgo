// Port of com/google/javascript/jscomp/transpile/BaseTranspiler.java @ closure-compiler v20261021

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

package transpile

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/bundle"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
)

// BaseTranspiler is the basic Transpiler implementation for outputting ES5
// code.
//
// Java: final class transpile.BaseTranspiler implements Transpiler
type BaseTranspiler struct {
	compilerSupplier   BaseTranspiler_CompilerSupplier
	runtimeLibraryName string
}

var _ Transpiler = (*BaseTranspiler)(nil)

// Java: BaseTranspiler(CompilerSupplier, String)
func NewBaseTranspiler(compilerSupplier BaseTranspiler_CompilerSupplier, runtimeLibraryName string) *BaseTranspiler {
	t := &BaseTranspiler{}
	t.compilerSupplier = jutil.CheckNotNull(compilerSupplier)
	t.runtimeLibraryName = runtimeLibraryName
	return t
}

// Java: BaseTranspiler.transpile(URI, String)
func (t *BaseTranspiler) Transpile(path *jutil.URI, code string) *TranspileResult {
	result := t.compilerSupplier.Compile(path, code)
	if !result.transpiled {
		return NewTranspileResult(path, code, code, "")
	}
	return NewTranspileResult(path, code, result.source, result.sourceMap)
}

// Java: BaseTranspiler.runtime()
func (t *BaseTranspiler) Runtime() string {
	var sb strings.Builder
	if t.runtimeLibraryName != "" {
		sb.WriteString(t.compilerSupplier.Runtime(t.runtimeLibraryName))
	}
	sb.WriteString(t.compilerSupplier.Runtime("modules"))
	return sb.String()
}

// Java: BaseTranspiler.LATEST_TRANSPILER
var BaseTranspiler_LATEST_TRANSPILER = BaseTranspiler_ToWithRuntime(parser.FeatureSet_Latest(), "")

// Java: BaseTranspiler.ES5_TRANSPILER
var BaseTranspiler_ES5_TRANSPILER = BaseTranspiler_To(jscomp.CompilerOptions_LanguageMode_ECMASCRIPT5.ToFeatureSet())

// Java: BaseTranspiler.to(FeatureSet, String)
func BaseTranspiler_ToWithRuntime(featureSet *parser.FeatureSet, runtime string) *BaseTranspiler {
	return NewBaseTranspiler(NewBaseTranspiler_CompilerSupplierWithOutputFeatureSet(featureSet), runtime)
}

// Java: BaseTranspiler.to(FeatureSet)
func BaseTranspiler_To(featureSet *parser.FeatureSet) *BaseTranspiler {
	return BaseTranspiler_ToWithRuntime(featureSet, "es6_runtime")
}

// BaseTranspiler_CompilerSupplier wraps the Compiler into a more relevant
// interface, making it easy to test the Transpiler without depending on
// implementation details of the Compiler itself. Also works around the fact
// that the Compiler is not thread-safe (since we may do multiple transpiles
// concurrently), so we supply a fresh instance each time when we're in
// single-file mode.
//
// Java: static class BaseTranspiler.CompilerSupplier
type BaseTranspiler_CompilerSupplier interface {
	// Java: CompilerSupplier.compile(URI, String)
	Compile(path *jutil.URI, code string) *BaseTranspiler_CompileResult
	// Java: CompilerSupplier.runtime(String)
	Runtime(library string) string
	// Java: CompilerSupplier.compiler()
	Compiler() *jscomp.Compiler
	// Java: CompilerSupplier.options()
	Options() *jscomp.CompilerOptions
	// Java: CompilerSupplier.setOptions(CompilerOptions)
	SetOptions(options *jscomp.CompilerOptions)
	// Java: CompilerSupplier.createTrivialExterns()
	CreateTrivialExterns() *jscomp.SourceFile
	// Java: CompilerSupplier.createEmptySource()
	CreateEmptySource() *jscomp.SourceFile
}

// BaseTranspiler_CompilerSupplierBase holds the fields of CompilerSupplier
// and its implementation; subclasses embed it.
type BaseTranspiler_CompilerSupplierBase struct {
	self               BaseTranspiler_CompilerSupplier
	moduleResolution   jscomp.ModuleLoader_ResolutionMode
	moduleRoots        []string
	prefixReplacements *jutil.LinkedHashMap[string, string]
	outputFeatureSet   *parser.FeatureSet
}

// Java: CompilerSupplier()
func NewBaseTranspiler_CompilerSupplier() *BaseTranspiler_CompilerSupplierBase {
	return NewBaseTranspiler_CompilerSupplierWithOutputFeatureSet(jscomp.CompilerOptions_LanguageMode_ECMASCRIPT5.ToFeatureSet())
}

// Java: CompilerSupplier(FeatureSet)
func NewBaseTranspiler_CompilerSupplierWithOutputFeatureSet(outputFeatureSet *parser.FeatureSet) *BaseTranspiler_CompilerSupplierBase {
	// Use the default resolution mode
	return NewBaseTranspiler_CompilerSupplierWithModuleResolution(
		outputFeatureSet,
		jscomp.NewCompilerOptions().GetModuleResolutionMode(),
		[]string{},
		jutil.NewLinkedHashMap[string, string]())
}

// NewBaseTranspiler_CompilerSupplierWithModuleResolution accepts commonly
// overridden options for ES6 modules to avoid needed to subclass.
//
// moduleResolution: module resolution for resolving import paths;
// prefixReplacements: prefix replacements for when moduleResolution is
// ModuleLoader.ResolutionMode.BROWSER_WITH_TRANSFORMED_PREFIXES.
//
// Java: CompilerSupplier(FeatureSet, ResolutionMode, ImmutableList<String>, ImmutableMap<String, String>)
func NewBaseTranspiler_CompilerSupplierWithModuleResolution(
	outputFeatureSet *parser.FeatureSet,
	moduleResolution jscomp.ModuleLoader_ResolutionMode,
	moduleRoots []string,
	prefixReplacements *jutil.LinkedHashMap[string, string]) *BaseTranspiler_CompilerSupplierBase {
	s := &BaseTranspiler_CompilerSupplierBase{}
	s.InitBaseTranspiler_CompilerSupplierBase(s, outputFeatureSet, moduleResolution, moduleRoots, prefixReplacements)
	return s
}

// InitBaseTranspiler_CompilerSupplierBase is the constructor as subclasses
// call it; self is the subclass instance.
//
// Java: CompilerSupplier(FeatureSet, ResolutionMode, ImmutableList<String>, ImmutableMap<String, String>)
func (s *BaseTranspiler_CompilerSupplierBase) InitBaseTranspiler_CompilerSupplierBase(
	self BaseTranspiler_CompilerSupplier,
	outputFeatureSet *parser.FeatureSet,
	moduleResolution jscomp.ModuleLoader_ResolutionMode,
	moduleRoots []string,
	prefixReplacements *jutil.LinkedHashMap[string, string]) {
	s.self = self
	s.outputFeatureSet = outputFeatureSet
	s.moduleResolution = moduleResolution
	s.moduleRoots = moduleRoots
	s.prefixReplacements = prefixReplacements
}

// Java: CompilerSupplier.compile(URI, String)
func (s *BaseTranspiler_CompilerSupplierBase) Compile(path *jutil.URI, code string) *BaseTranspiler_CompileResult {
	compiler := s.self.Compiler()
	result :=
		compiler.CompileOfSourceFiles(
			s.self.CreateTrivialExterns(), jscomp.SourceFile_FromCode(path.String(), code), s.self.Options())
	source := compiler.ToSource()
	var sourceMap strings.Builder
	if result.SourceMap != nil {
		name := path.String()
		result.SourceMap.AppendTo(&sourceMap, &name)
	}
	transpiled := result.TranspiledFiles
	if len(result.Errors) != 0 {
		panic(bundle.NewTranspilationException(compiler, result.Errors, result.Warnings))
	}
	sm := ""
	if transpiled {
		sm = sourceMap.String()
	}
	return NewBaseTranspiler_CompileResult(source, transpiled, sm)
}

// Java: CompilerSupplier.runtime(String)
func (s *BaseTranspiler_CompilerSupplierBase) Runtime(library string) string {
	compiler := s.self.Compiler()
	options := s.self.Options()
	options.SetForceLibraryInjection([]string{library})
	compiler.CompileOfSourceFiles(s.self.CreateTrivialExterns(), s.self.CreateEmptySource(), options)
	return compiler.ToSource()
}

// Java: CompilerSupplier.compiler()
func (s *BaseTranspiler_CompilerSupplierBase) Compiler() *jscomp.Compiler {
	return jscomp.NewCompiler()
}

// Java: CompilerSupplier.options()
func (s *BaseTranspiler_CompilerSupplierBase) Options() *jscomp.CompilerOptions {
	options := jscomp.NewCompilerOptions()
	s.self.SetOptions(options)
	return options
}

// Java: CompilerSupplier.setOptions(CompilerOptions)
func (s *BaseTranspiler_CompilerSupplierBase) SetOptions(options *jscomp.CompilerOptions) {
	options.SetLanguageIn(jscomp.CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.LegacySetOutputFeatureSet(s.outputFeatureSet.WithoutFeatures(parser.FeatureSet_Feature_MODULES))
	options.SetEmitUseStrict(false)
	options.SetQuoteKeywordProperties(true)
	options.SetSkipNonTranspilationPasses(true)
	options.SetVariableRenaming(jscomp.VariableRenamingPolicy_OFF)
	options.SetPropertyRenaming(jscomp.PropertyRenamingPolicy_OFF)
	options.SetWrapGoogModulesForWhitespaceOnly(false)
	options.SetPrettyPrint(true)
	options.SetWarningLevel(jscomp.DiagnosticGroups_NON_STANDARD_JSDOC, jscomp.CheckLevel_OFF)
	options.SetEs6ModuleTranspilation(jscomp.CompilerOptions_Es6ModuleTranspilation_TO_COMMON_JS_LIKE_MODULES)
	options.SetModuleResolutionMode(s.moduleResolution)
	options.SetModuleRoots(s.moduleRoots)
	options.SetBrowserResolverPrefixReplacements(s.prefixReplacements)

	// Transpiler often used in isolation, so references to other files will
	// never exist.
	options.SetWarningLevel(jscomp.DiagnosticGroups_MODULE_LOAD, jscomp.CheckLevel_OFF)

	// Don't escape module paths when bundling in the event paths are URLs.
	options.SetPathEscaper(jscomp.ModuleLoader_PathEscaper_CANONICALIZE_ONLY)
	options.SetParseInlineSourceMaps(true)
	options.SetApplyInputSourceMaps(true)
	devNull := "/dev/null"
	options.SetSourceMapOutputPath(&devNull)
	options.SetSourceMapIncludeSourcesContent(true)
	// Make sourcemaps use absolute paths, so that the path is not duplicated
	// if a build tool adds a sourceurl. Exception: if the location has a
	// scheme (like http:) then leave the path intact. This makes this usable
	// from web servers.
	options.SetSourceMapLocationMappings(
		[]jscomp.SourceMap_LocationMapping{baseTranspiler_absoluteLocationMapping{}})
}

// baseTranspiler_absoluteLocationMapping is the location mapping lambda of
// CompilerSupplier.setOptions.
type baseTranspiler_absoluteLocationMapping struct{}

// Java: SourceMap.LocationMapping.map(String)
func (baseTranspiler_absoluteLocationMapping) Map(location string) *string {
	hasScheme := false
	jutil.Try(func() {
		if jutil.NewURI(location).GetScheme() != nil {
			hasScheme = true
		}
	}, func(e jutil.IURISyntaxException) {
		// Swallow, return the absolute version below.
	})
	if hasScheme {
		return &location
	}
	return jscomp.NewSourceMap_PrefixLocationMapping("", "/").Map(location)
}

// Java: CompilerSupplier.createTrivialExterns()
func (s *BaseTranspiler_CompilerSupplierBase) CreateTrivialExterns() *jscomp.SourceFile {
	return jscomp.SourceFile_FromCode("externs.js", "function Symbol() {}")
}

// Java: CompilerSupplier.createEmptySource()
func (s *BaseTranspiler_CompilerSupplierBase) CreateEmptySource() *jscomp.SourceFile {
	return jscomp.SourceFile_FromCode("empty.js", "")
}

// BaseTranspiler_CompileResult is the source together with the additional
// compilation results.
//
// Java: static class BaseTranspiler.CompileResult
type BaseTranspiler_CompileResult struct {
	source     string
	transpiled bool
	sourceMap  string
}

// Java: CompileResult(String, boolean, String)
func NewBaseTranspiler_CompileResult(source string, transpiled bool, sourceMap string) *BaseTranspiler_CompileResult {
	r := &BaseTranspiler_CompileResult{}
	r.source = source
	r.transpiled = transpiled
	r.sourceMap = sourceMap
	return r
}
