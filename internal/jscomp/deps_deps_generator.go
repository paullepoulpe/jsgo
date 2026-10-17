// Port of com/google/javascript/jscomp/deps/DepsGenerator.java @ closure-compiler v20260415

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"bytes"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// DepsGenerator generates deps.js files by scanning JavaScript files for
// calls to goog.provide(), goog.require() and goog.addDependency().
//
// Java: class deps.DepsGenerator
type DepsGenerator struct {
	// See the Flags in MakeJsDeps for descriptions of these.
	srcs           []*SourceFile
	deps           []*SourceFile
	closurePathAbs string
	mergeStrategy  DepsGenerator_InclusionStrategy
	loader         *ModuleLoader
	errorManager   ErrorManager
}

// Java: enum DepsGenerator.InclusionStrategy
type DepsGenerator_InclusionStrategy int32

const (
	DepsGenerator_InclusionStrategy_ALWAYS DepsGenerator_InclusionStrategy = iota
	DepsGenerator_InclusionStrategy_WHEN_IN_SRCS
	DepsGenerator_InclusionStrategy_DO_NOT_DUPLICATE
)

var depsGenerator_InclusionStrategy_names = [...]string{"ALWAYS", "WHEN_IN_SRCS", "DO_NOT_DUPLICATE"}

// Java: DepsGenerator.InclusionStrategy.values()
func DepsGenerator_InclusionStrategy_values() []DepsGenerator_InclusionStrategy {
	return []DepsGenerator_InclusionStrategy{
		DepsGenerator_InclusionStrategy_ALWAYS,
		DepsGenerator_InclusionStrategy_WHEN_IN_SRCS,
		DepsGenerator_InclusionStrategy_DO_NOT_DUPLICATE,
	}
}

// String is InclusionStrategy.name().
func (s DepsGenerator_InclusionStrategy) String() string {
	return depsGenerator_InclusionStrategy_names[s]
}

var depsGenerator_logger = jutil.Logger_GetLogger("com.google.javascript.jscomp.deps.DepsGenerator")

// Java: DepsGenerator.ES6_IMPORT_FOR_NON_ES6_MODULE
var DepsGenerator_ES6_IMPORT_FOR_NON_ES6_MODULE = DiagnosticType_Warning(
	"DEPS_ES6_IMPORT_FOR_NON_ES6_MODULE",
	"Cannot import file \"{0}\" because it is not an ES6 module.")

// Java: DepsGenerator.UNKNOWN_PATH_IMPORT
var DepsGenerator_UNKNOWN_PATH_IMPORT = DiagnosticType_Warning("DEPS_UNKNOWN_PATH_IMPORT", "Could not find file \"{0}\".")

// Java: DepsGenerator.SAME_FILE_WARNING
var DepsGenerator_SAME_FILE_WARNING = DiagnosticType_Warning(
	"DEPS_SAME_FILE",
	"Namespace \"{0}\" is both required and provided in the same file.")

// Java: DepsGenerator.NEVER_PROVIDED_ERROR
var DepsGenerator_NEVER_PROVIDED_ERROR = DiagnosticType_Error(
	"DEPS_NEVER_PROVIDED",
	"Namespace \"{0}\" is required but never provided.\nYou "+
		"need to pass a library that has it in srcs or exports to your target''s deps.")

// Java: DepsGenerator.DUPE_PROVIDES_WARNING
var DepsGenerator_DUPE_PROVIDES_WARNING = DiagnosticType_Warning(
	"DEPS_DUPE_PROVIDES",
	"Multiple calls to goog.provide(\"{0}\")")

// Java: DepsGenerator.MULTIPLE_PROVIDES_ERROR
var DepsGenerator_MULTIPLE_PROVIDES_ERROR = DiagnosticType_Error(
	"DEPS_DUPE_PROVIDES",
	"Namespace \"{0}\" is already provided in other file {1}")

// Java: DepsGenerator.NO_DEPS_WARNING
var DepsGenerator_NO_DEPS_WARNING = DiagnosticType_Warning(
	"DEPS_NO_DEPS",
	"No dependencies found in file")

// NewDepsGenerator creates a new DepsGenerator.
//
// Java: DepsGenerator(Collection<SourceFile>, Collection<SourceFile>, InclusionStrategy, String, ErrorManager, ModuleLoader)
func NewDepsGenerator(
	deps []*SourceFile,
	srcs []*SourceFile,
	mergeStrategy DepsGenerator_InclusionStrategy,
	closurePathAbs string,
	errorManager ErrorManager,
	loader *ModuleLoader) *DepsGenerator {
	g := &DepsGenerator{}
	g.deps = deps
	g.srcs = srcs
	g.mergeStrategy = mergeStrategy
	g.closurePathAbs = closurePathAbs
	g.errorManager = errorManager
	g.loader = loader
	return g
}

// GetErrorManager is the package-private field errorManager, read by
// tests.
func (g *DepsGenerator) GetErrorManager() ErrorManager { return g.errorManager }

// ComputeDependencyCalls performs the parsing inputs and writing of outputs.
//
// It throws IOException upon an IO error. It returns a String of
// goog.addDependency calls that will build the dependency graph, or nil if
// there was an error.
//
// Java: DepsGenerator.computeDependencyCalls()
func (g *DepsGenerator) ComputeDependencyCalls() *string {
	// Build a map of closure-relative path -> DepInfo.
	depsFiles := g.parseDepsFiles()
	if depsGenerator_logger.IsLoggable(jutil.Level_FINE) {
		depsGenerator_logger.Fine("preparsedFiles: " + jutil.ToJavaString(depsFiles))
	}
	// Find all goog.provides & goog.requires in src files
	jsFiles := g.parseSources(jutil.NewLinkedHashSet(depsFiles.KeySlice()...))

	// Check if there were any parse errors.
	if g.errorManager.GetErrorCount() > 0 {
		return nil
	}

	g.CleanUpDuplicatedFiles(depsFiles, jsFiles)

	jsFiles = g.removeMungedSymbols(depsFiles, jsFiles)

	// Check for missing provides or other semantic inconsistencies.
	g.validateDependencies(slices.Collect(depsFiles.Values()), slices.Collect(jsFiles.Values()))

	if g.errorManager.GetErrorCount() > 0 {
		return nil
	}

	var output bytes.Buffer
	g.writeDepsContent(depsFiles, jsFiles, jutil.NewPrintStream(&output))
	s := jutil.NewStringFromUTF8Bytes(output.Bytes())
	return &s
}

// CleanUpDuplicatedFiles removes duplicated depsInfo from jsFiles if this
// info already present in some of the parsed deps.js. depsFiles is the
// DepsInfo from deps.js dependencies, jsFiles the DepsInfo from some of
// jsSources.
//
// Java: DepsGenerator.cleanUpDuplicatedFiles(Map<String, DependencyInfo>, Map<String, DependencyInfo>)
func (g *DepsGenerator) CleanUpDuplicatedFiles(depsFiles *jutil.LinkedHashMap[string, DependencyInfo],
	jsFiles *jutil.LinkedHashMap[string, DependencyInfo]) {
	// Go-only: new HashSet<>(depsFiles.keySet()) is iterated only to remove
	// its elements from another map, so its order is not observable.
	depsPathsCopy := depsFiles.KeySlice()
	for _, path := range depsPathsCopy {
		if g.mergeStrategy != DepsGenerator_InclusionStrategy_WHEN_IN_SRCS {
			jsFiles.Remove(path)
		}
	}

	for path := range jsFiles.Keys() {
		// If a generated file appears in both the jsFiles and in depsFiles,
		// then remove it from depsFiles in order to get the full path the
		// generated file.
		depsFiles.Remove(path)
	}
}

// removeMungedSymbols removes munged symbols in requires and provides. These
// munged symbols are from ES6 modules and are generated by
// ModulePath.toModuleName(). We do not wish to write these munged symbols to
// the dependency file.
//
//   - Makes any require'd munged symbol the require'd file's relative path to
//     Closure.
//   - Removes munged symbols from the provides list.
//
// Java: DepsGenerator.removeMungedSymbols(Map<String, DependencyInfo>, Map<String, DependencyInfo>)
func (g *DepsGenerator) removeMungedSymbols(
	depFiles *jutil.LinkedHashMap[string, DependencyInfo], jsFiles *jutil.LinkedHashMap[string, DependencyInfo]) *jutil.LinkedHashMap[string, DependencyInfo] {
	newJsFiles := jutil.NewLinkedHashMap[string, DependencyInfo]()
	providesMap := jutil.NewLinkedHashMap[string, DependencyInfo]()
	g.addToProvideMap(slices.Collect(depFiles.Values()), providesMap, true)
	g.addToProvideMap(slices.Collect(jsFiles.Values()), providesMap, false)

	for dependencyInfo := range jsFiles.Values() {
		newRequires := []*DependencyInfo_Require{}
		for _, require := range dependencyInfo.GetRequires() {
			if require.GetType() == DependencyInfo_Require_Type_ES6_IMPORT {
				// Symbols are unique per file and have nothing to do with paths
				// so map lookups are safe here.
				provider := providesMap.Get(require.GetSymbol())
				if provider == nil {
					g.reportMissingFile(dependencyInfo, require.GetRawText())
				} else {
					// If this is an ES6 module then set the symbol to be its
					// relative path to Closure. ES6 modules in a dependency file
					// do not "provide" anything. Requires can match a provided
					// symbol or a relative path to Closure.
					newRequires = append(newRequires, require.WithSymbol(provider.GetPathRelativeToClosureBase()))
				}
			} else {
				// Require is by symbol already so no need to change it.
				newRequires = append(newRequires, require)
			}
		}

		provides := dependencyInfo.GetProvides()

		if dependencyInfo.IsEs6Module() {
			mungedProvide := g.loader.Resolve(dependencyInfo.GetName()).ToModuleName()
			// Filter out the munged symbol.
			// Note that at the moment ES6 modules should not have any other
			// provides! In the future we may have additional mechanisms to add
			// goog symbols. But for not nothing is officially supported.
			builder := []string{}
			for _, provide := range provides {
				if provide != mungedProvide &&
					provide != dependencyInfo.GetPathRelativeToClosureBase() {
					builder = append(builder, provide)
				}
			}
			provides = builder
		}

		newJsFiles.Put(
			dependencyInfo.GetPathRelativeToClosureBase(),
			SimpleDependencyInfo_Builder_(
				dependencyInfo.GetPathRelativeToClosureBase(), dependencyInfo.GetName()).
				SetProvides(provides).
				SetRequires(newRequires).
				SetLoadFlags(dependencyInfo.GetLoadFlags()).
				Build())
	}
	return newJsFiles
}

// validateDependencies reports if there are any dependency problems with the
// given dependency information. Reported problems include:
//   - A namespace being provided more than once
//   - A namespace being required multiple times from within one file
//   - A namespace being provided and required in the same file
//   - A namespace being required that is never provided
//
// preparsedFileDependencies is the dependency information from existing
// deps.js files, parsedFileDependencies the dependency information from
// parsed .js files.
//
// Java: DepsGenerator.validateDependencies(Iterable<DependencyInfo>, Iterable<DependencyInfo>)
func (g *DepsGenerator) validateDependencies(preparsedFileDependencies []DependencyInfo,
	parsedFileDependencies []DependencyInfo) {
	// Create a map of namespace -> file providing it.
	// Also report any duplicate provides.
	providesMap := jutil.NewLinkedHashMap[string, DependencyInfo]()
	g.addToProvideMap(preparsedFileDependencies, providesMap, true)
	g.addToProvideMap(parsedFileDependencies, providesMap, false)
	// For each require in the parsed sources:
	for _, depInfo := range parsedFileDependencies {
		for _, require := range depInfo.GetRequires() {
			namespace := require.GetSymbol()
			// Check for missing provides.
			provider := providesMap.Get(namespace)
			if provider == nil {
				g.reportUndefinedNamespace(namespace, depInfo)
			} else if provider == depInfo {
				g.reportSameFile(namespace, depInfo)
			} else {
				providerIsEs6Module := provider.IsEs6Module()

				switch require.GetType() {
				case DependencyInfo_Require_Type_ES6_IMPORT:
					if !providerIsEs6Module {
						g.reportEs6ImportForNonEs6Module(provider, depInfo)
					}
				case DependencyInfo_Require_Type_GOOG_REQUIRE_SYMBOL, DependencyInfo_Require_Type_PARSED_FROM_DEPS:
				case DependencyInfo_Require_Type_COMMON_JS, DependencyInfo_Require_Type_COMPILER_CHUNK:
					panic(jutil.NewIllegalStateExceptionMsg("Unexpected import type: " + require.GetType().String()))
				}
			}
		}
	}
}

// Java: DepsGenerator.reportMissingFile(DependencyInfo, String)
func (g *DepsGenerator) reportMissingFile(depInfo DependencyInfo, path string) {
	name := depInfo.GetName()
	g.errorManager.Report(
		CheckLevel_ERROR.Ptr(), JSError_MakeWithSourceLocation(&name, -1, -1, DepsGenerator_UNKNOWN_PATH_IMPORT, path))
}

// Java: DepsGenerator.reportEs6ImportForNonEs6Module(DependencyInfo, DependencyInfo)
func (g *DepsGenerator) reportEs6ImportForNonEs6Module(provider DependencyInfo, depInfo DependencyInfo) {
	name := depInfo.GetName()
	g.errorManager.Report(
		CheckLevel_ERROR.Ptr(),
		JSError_MakeWithSourceLocation(&name, -1, -1, DepsGenerator_ES6_IMPORT_FOR_NON_ES6_MODULE, provider.GetName()))
}

// Java: DepsGenerator.reportSameFile(String, DependencyInfo)
func (g *DepsGenerator) reportSameFile(namespace string, depInfo DependencyInfo) {
	name := depInfo.GetName()
	g.errorManager.Report(CheckLevel_WARNING.Ptr(),
		JSError_MakeWithSourceLocation(&name, -1, -1,
			DepsGenerator_SAME_FILE_WARNING, namespace))
}

// Java: DepsGenerator.reportUndefinedNamespace(String, DependencyInfo)
func (g *DepsGenerator) reportUndefinedNamespace(
	namespace string, depInfo DependencyInfo) {
	name := depInfo.GetName()
	g.errorManager.Report(CheckLevel_ERROR.Ptr(),
		JSError_MakeWithSourceLocation(&name, -1, -1,
			DepsGenerator_NEVER_PROVIDED_ERROR, namespace))
}

// Java: DepsGenerator.reportDuplicateProvide(String, DependencyInfo, DependencyInfo)
func (g *DepsGenerator) reportDuplicateProvide(namespace string, firstDep DependencyInfo,
	secondDep DependencyInfo) {
	if firstDep == secondDep {
		if firstDep.GetPathRelativeToClosureBase() != namespace {
			name := firstDep.GetName()
			g.errorManager.Report(
				CheckLevel_WARNING.Ptr(),
				JSError_MakeWithSourceLocation(&name, -1, -1, DepsGenerator_DUPE_PROVIDES_WARNING, namespace))
		}
	} else {
		name := secondDep.GetName()
		g.errorManager.Report(CheckLevel_ERROR.Ptr(),
			JSError_MakeWithSourceLocation(&name, -1, -1,
				DepsGenerator_MULTIPLE_PROVIDES_ERROR, namespace, firstDep.GetName()))
	}
}

// Java: DepsGenerator.reportNoDepsInDepsFile(String)
func (g *DepsGenerator) reportNoDepsInDepsFile(filePath string) {
	g.errorManager.Report(CheckLevel_WARNING.Ptr(),
		JSError_MakeWithSourceLocation(&filePath, -1, -1, DepsGenerator_NO_DEPS_WARNING))
}

// addToProvideMap adds the given DependencyInfos to the given providesMap.
// Also checks for and reports duplicate provides.
//
// Java: DepsGenerator.addToProvideMap(Iterable<DependencyInfo>, Map<String, DependencyInfo>, boolean)
func (g *DepsGenerator) addToProvideMap(
	depInfos []DependencyInfo,
	providesMap *jutil.LinkedHashMap[string, DependencyInfo],
	isFromDepsFile bool) {
	for _, depInfo := range depInfos {
		provides := append([]string{}, depInfo.GetProvides()...)

		// Add a munged symbol to the provides map so that lookups by path
		// requires work as intended.
		if isFromDepsFile {
			// Don't add the dependency file itself but every file it says
			// exists instead.
			provides = append(provides,
				g.loader.
					Resolve(
						PathUtil_MakeAbsoluteWithRoot(depInfo.GetPathRelativeToClosureBase(), g.closurePathAbs)).
					ToModuleName())
		} else {
			// ES6 modules already provide these munged symbols.
			if !depInfo.IsEs6Module() {
				provides = append(provides, g.loader.Resolve(depInfo.GetName()).ToModuleName())
			}
		}

		// Also add the relative closure path as a provide. At some point we'll
		// swap out the munged symbols for these relative paths. So looks ups
		// by either need to work.
		provides = append(provides, depInfo.GetPathRelativeToClosureBase())

		for _, provide := range provides {
			prevValue := providesMap.Put(provide, depInfo)
			// Check for duplicate provides.
			if prevValue != nil {
				g.reportDuplicateProvide(provide, prevValue, depInfo)
			}
		}
	}
}

// Java: DepsGenerator.createDepsFileParser()
func (g *DepsGenerator) CreateDepsFileParser() *DepsFileRegexParser {
	depsParser := NewDepsFileRegexParser(g.errorManager)
	depsParser.SetShortcutMode(true)
	return depsParser
}

// ShouldSkipDepsFile returns whether we should ignore dependency info in the
// given deps file.
//
// Java: DepsGenerator.shouldSkipDepsFile(SourceFile)
func (g *DepsGenerator) ShouldSkipDepsFile(file *SourceFile) bool {
	return false
}

// parseDepsFiles parses all deps.js files in the deps list and creates a
// map of closure-relative path -> DependencyInfo. It throws IOException.
//
// Java: DepsGenerator.parseDepsFiles()
func (g *DepsGenerator) parseDepsFiles() *jutil.LinkedHashMap[string, DependencyInfo] {
	depsParser := g.CreateDepsFileParser()
	depsFiles := jutil.NewLinkedHashMap[string, DependencyInfo]()
	for _, file := range g.deps {
		if !g.ShouldSkipDepsFile(file) {
			depInfos := depsParser.ParseFileReader(
				file.GetName(), file.GetCodeReader())
			if len(depInfos) == 0 {
				g.reportNoDepsInDepsFile(file.GetName())
			} else {
				for _, info := range depInfos {
					depsFiles.Put(info.GetPathRelativeToClosureBase(), g.removeRelativePathProvide(info))
				}
			}
		}
	}

	// If a deps file also appears in srcs, our build tools will move it
	// into srcs.  So we need to scan all the src files for addDependency
	// calls as well.
	for _, src := range g.srcs {
		if !g.ShouldSkipDepsFile(src) {
			srcInfos :=
				depsParser.ParseFileReader(src.GetName(), src.GetCodeReader())
			for _, info := range srcInfos {
				depsFiles.Put(info.GetPathRelativeToClosureBase(), g.removeRelativePathProvide(info))
			}
		}
	}

	return depsFiles
}

// Java: DepsGenerator.removeRelativePathProvide(DependencyInfo)
func (g *DepsGenerator) removeRelativePathProvide(info DependencyInfo) DependencyInfo {
	// DepsFileRegexParser adds an ES6 module's relative path to closure as a
	// provide so that the resulting depgraph is valid. But we don't want to
	// write this "fake" provide back out, so remove it here.
	provides := []string{}
	for _, p := range info.GetProvides() {
		if p != info.GetPathRelativeToClosureBase() {
			provides = append(provides, p)
		}
	}
	return SimpleDependencyInfo_Builder_From(info).
		SetProvides(provides).
		Build()
}

// parseSources parses all source files for dependency information.
// preparsedFiles is a set of closure-relative paths: files in this set are
// not parsed if they are encountered in srcs. It returns a map of
// closure-relative paths -> DependencyInfo for the newly parsed files. It
// throws IOException upon an IO error.
//
// Java: DepsGenerator.parseSources(Set<String>)
func (g *DepsGenerator) parseSources(
	preparsedFiles *jutil.LinkedHashSet[string]) *jutil.LinkedHashMap[string, DependencyInfo] {
	parsedFiles := jutil.NewLinkedHashMap[string, DependencyInfo]()
	jsParser := NewJsFileRegexParser(g.errorManager).SetModuleLoader(g.loader)
	compiler := NewCompiler()
	compiler.Init([]*SourceFile{}, []*SourceFile{}, NewCompilerOptions())

	for _, file := range g.srcs {
		closureRelativePath :=
			PathUtil_MakeRelative(
				g.closurePathAbs, PathUtil_MakeAbsolute(file.GetName()))
		if depsGenerator_logger.IsLoggable(jutil.Level_FINE) {
			depsGenerator_logger.Fine("Closure-relative path: " + closureRelativePath)
		}
		if DepsGenerator_InclusionStrategy_WHEN_IN_SRCS == g.mergeStrategy ||
			!preparsedFiles.Contains(closureRelativePath) {
			var depInfo DependencyInfo = jsParser.ParseFile(
				file.GetName(), closureRelativePath,
				file.GetCode())
			depInfo = NewLazyParsedDependencyInfo(depInfo, NewCompilerInput(file), compiler)

			// Skip externs files, which should never be loaded.
			if depInfo.GetHasExternsAnnotation() {
				continue
			}

			// Kick the source out of memory.
			file.ClearCachedSource()
			parsedFiles.Put(closureRelativePath, depInfo)
		}
	}

	return parsedFiles
}

// writeDepsContent creates the content to put into the output deps.js file.
// If mergeDeps is true, then all of the dependency information in the
// providedDeps will be included in the output. It throws IOException upon
// an IO error.
//
// Java: DepsGenerator.writeDepsContent(Map<String, DependencyInfo>, Map<String, DependencyInfo>, PrintStream)
func (g *DepsGenerator) writeDepsContent(depsFiles *jutil.LinkedHashMap[string, DependencyInfo],
	jsFiles *jutil.LinkedHashMap[string, DependencyInfo], out *jutil.PrintStream) {
	// Print all dependencies extracted from srcs.
	depsGenerator_writeDepInfos(out, slices.Collect(jsFiles.Values()))

	// Print all dependencies extracted from deps.
	if g.mergeStrategy == DepsGenerator_InclusionStrategy_ALWAYS {
		// This multimap is just for splitting DepsInfo objects by
		// it's definition deps.js file
		//
		// Multimaps.index: keys in order of first appearance, each with its
		// values in order.
		infosIndex := jutil.NewLinkedHashMap[string, []DependencyInfo]()
		for info := range depsFiles.Values() {
			infosIndex.Put(info.GetName(), append(infosIndex.Get(info.GetName()), info))
		}

		for depsPath := range infosIndex.Keys() {
			path := g.FormatPathToDepsFile(depsPath)
			out.Println("\n// Included from: " + path)
			depsGenerator_writeDepInfos(out, infosIndex.Get(depsPath))
		}
	}
}

// FormatPathToDepsFile formats the deps file path so that it can be
// included in the output file.
//
// Java: DepsGenerator.formatPathToDepsFile(String)
func (g *DepsGenerator) FormatPathToDepsFile(path string) string {
	return path
}

// depsGenerator_writeDepInfos writes goog.addDependency() lines for each
// DependencyInfo in depInfos. It throws IOException.
//
// Java: DepsGenerator.writeDepInfos(PrintStream, Collection<DependencyInfo>)
func depsGenerator_writeDepInfos(out *jutil.PrintStream, depInfos []DependencyInfo) {
	// Print dependencies.
	// Lines look like this:
	// goog.addDependency('../../path/to/file.js', ['goog.Delay'],
	//     ['goog.Disposable', 'goog.Timer']);
	for _, depInfo := range depInfos {
		DependencyInfo_Util_WriteAddDependency(out, depInfo)
	}
}

// Java: DepsGenerator.createSourceFilesFromPaths(Collection<String>)
func DepsGenerator_CreateSourceFilesFromPaths(
	paths []string) []*SourceFile {
	files := []*SourceFile{}
	for _, path := range paths {
		files = append(files, SourceFile_FromFile(path))
	}
	return files
}

// Java: DepsGenerator.createSourceFilesFromPaths(String...)
func DepsGenerator_CreateSourceFilesFromPathsArray(paths ...string) []*SourceFile {
	return DepsGenerator_CreateSourceFilesFromPaths(append([]string{}, paths...))
}

// DepsGenerator_CreateSourceFilesFromZipPaths throws IOException.
//
// Java: DepsGenerator.createSourceFilesFromZipPaths(Collection<String>)
func DepsGenerator_CreateSourceFilesFromZipPaths(
	paths []string) []*SourceFile {
	zipSourceFiles := []*SourceFile{}
	for _, path := range paths {
		zipSourceFiles = append(zipSourceFiles, SourceFile_FromZipFile(path, jutil.StandardCharsets_UTF_8)...)
	}
	return zipSourceFiles
}
