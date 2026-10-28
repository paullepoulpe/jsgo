// Port of com/google/javascript/jscomp/deps/DepsGenerator.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// DepsGenerator generates deps.js files by scanning JavaScript files for
// calls to goog.provide(), goog.require() and goog.addDependency().
type DepsGenerator struct {
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

func (s DepsGenerator_InclusionStrategy) String() string {
	return depsGenerator_InclusionStrategy_names[s]
}

var depsGenerator_logger = jutil.Logger_GetLogger("com.google.javascript.jscomp.deps.DepsGenerator")

var (
	DepsGenerator_ES6_IMPORT_FOR_NON_ES6_MODULE = DiagnosticType_Warning(
		"DEPS_ES6_IMPORT_FOR_NON_ES6_MODULE",
		"Cannot import file \"{0}\" because it is not an ES6 module.")
	DepsGenerator_UNKNOWN_PATH_IMPORT = DiagnosticType_Warning(
		"DEPS_UNKNOWN_PATH_IMPORT",
		"Could not find file \"{0}\".")
	DepsGenerator_SAME_FILE_WARNING = DiagnosticType_Warning(
		"DEPS_SAME_FILE",
		"Namespace \"{0}\" is both required and provided in the same file.")
	DepsGenerator_NEVER_PROVIDED_ERROR = DiagnosticType_Error(
		"DEPS_NEVER_PROVIDED",
		"Namespace \"{0}\" is required but never provided.\nYou "+
			"need to pass a library that has it in srcs or exports to your target''s deps.")
	DepsGenerator_DUPE_PROVIDES_WARNING = DiagnosticType_Warning(
		"DEPS_DUPE_PROVIDES",
		"Multiple calls to goog.provide(\"{0}\")")
	DepsGenerator_MULTIPLE_PROVIDES_ERROR = DiagnosticType_Error(
		"DEPS_DUPE_PROVIDES",
		"Namespace \"{0}\" is already provided in other file {1}")
	DepsGenerator_NO_DEPS_WARNING = DiagnosticType_Warning(
		"DEPS_NO_DEPS",
		"No dependencies found in file")
)

type orderedDepMap struct {
	order  []string
	byPath map[string]DependencyInfo
}

func newOrderedDepMap() *orderedDepMap {
	return &orderedDepMap{byPath: make(map[string]DependencyInfo)}
}

func (m *orderedDepMap) put(path string, info DependencyInfo) {
	if _, exists := m.byPath[path]; !exists {
		m.order = append(m.order, path)
	}
	m.byPath[path] = info
}

func (m *orderedDepMap) contains(path string) bool {
	_, ok := m.byPath[path]
	return ok
}

func (m *orderedDepMap) remove(path string) {
	delete(m.byPath, path)
}

func (m *orderedDepMap) values() []DependencyInfo {
	out := make([]DependencyInfo, 0, len(m.byPath))
	for _, path := range m.order {
		if info, ok := m.byPath[path]; ok {
			out = append(out, info)
		}
	}
	return out
}

// Java: DepsGenerator(Collection<SourceFile>, Collection<SourceFile>, InclusionStrategy, String, ErrorManager, ModuleLoader)
func NewDepsGenerator(
	deps []*SourceFile,
	srcs []*SourceFile,
	mergeStrategy DepsGenerator_InclusionStrategy,
	closurePathAbs string,
	errorManager ErrorManager,
	loader *ModuleLoader) *DepsGenerator {
	return &DepsGenerator{
		deps:           deps,
		srcs:           srcs,
		mergeStrategy:  mergeStrategy,
		closurePathAbs: closurePathAbs,
		errorManager:   errorManager,
		loader:         loader,
	}
}

// ComputeDependencyCalls performs the parsing inputs and writing of outputs.
// Returns a String of goog.addDependency calls that will build the dependency
// graph, or nil if there was an error.
//
// Java: DepsGenerator.computeDependencyCalls()
func (g *DepsGenerator) ComputeDependencyCalls() *string {
	depsFiles := g.parseDepsFiles()
	jsFiles := g.parseSources(depsFiles)
	if g.errorManager.GetErrorCount() > 0 {
		return nil
	}

	g.cleanUpDuplicatedFiles(depsFiles, jsFiles)
	jsValues := g.removeMungedSymbols(depsFiles.values(), jsFiles.values())
	depsValues := depsFiles.values()

	g.validateDependencies(depsValues, jsValues)
	if g.errorManager.GetErrorCount() > 0 {
		return nil
	}

	var output bytes.Buffer
	g.writeDepsContent(depsValues, jsValues, jutil.NewPrintStream(&output))
	s := jutil.NewStringFromUTF8Bytes(output.Bytes())
	return &s
}

func (g *DepsGenerator) cleanUpDuplicatedFiles(depsFiles, jsFiles *orderedDepMap) {
	if g.mergeStrategy != DepsGenerator_InclusionStrategy_WHEN_IN_SRCS {
		for _, path := range depsFiles.order {
			if depsFiles.contains(path) {
				jsFiles.remove(path)
			}
		}
	}
	for _, path := range jsFiles.order {
		if jsFiles.contains(path) {
			depsFiles.remove(path)
		}
	}
}

func (g *DepsGenerator) removeMungedSymbols(depValues, jsValues []DependencyInfo) []DependencyInfo {
	providesMap := make(map[string]DependencyInfo, len(depValues)+len(jsValues))
	g.addToProvideMap(depValues, providesMap, true)
	g.addToProvideMap(jsValues, providesMap, false)

	out := make([]DependencyInfo, 0, len(jsValues))
	for _, dependencyInfo := range jsValues {
		origReqs := dependencyInfo.GetRequires()
		newRequires := make([]*DependencyInfo_Require, 0, len(origReqs))
		for _, require := range origReqs {
			if require.GetType() == DependencyInfo_Require_Type_ES6_IMPORT {
				provider := providesMap[require.GetSymbol()]
				if provider == nil {
					g.reportMissingFile(dependencyInfo, require.GetRawText())
				} else {
					newRequires = append(newRequires, require.WithSymbol(provider.GetPathRelativeToClosureBase()))
				}
			} else {
				newRequires = append(newRequires, require)
			}
		}

		provides := dependencyInfo.GetProvides()
		if dependencyInfo.IsEs6Module() {
			mungedProvide := g.loader.Resolve(dependencyInfo.GetName()).ToModuleName()
			relPath := dependencyInfo.GetPathRelativeToClosureBase()
			filtered := make([]string, 0, len(provides))
			for _, provide := range provides {
				if provide != mungedProvide && provide != relPath {
					filtered = append(filtered, provide)
				}
			}
			provides = filtered
		}

		out = append(out, newSimpleDependencyInfoDirect(
			dependencyInfo.GetPathRelativeToClosureBase(),
			dependencyInfo.GetName(),
			provides,
			newRequires,
			nil,
			dependencyInfo.GetLoadFlags(),
			false,
			false,
		))
	}
	return out
}

func (g *DepsGenerator) validateDependencies(preparsedFileDependencies, parsedFileDependencies []DependencyInfo) {
	providesMap := make(map[string]DependencyInfo, len(preparsedFileDependencies)+len(parsedFileDependencies))
	g.addToProvideMap(preparsedFileDependencies, providesMap, true)
	g.addToProvideMap(parsedFileDependencies, providesMap, false)
	for _, depInfo := range parsedFileDependencies {
		for _, require := range depInfo.GetRequires() {
			namespace := require.GetSymbol()
			provider := providesMap[namespace]
			if provider == nil {
				g.reportUndefinedNamespace(namespace, depInfo)
			} else if provider == depInfo {
				g.reportSameFile(namespace, depInfo)
			} else {
				switch require.GetType() {
				case DependencyInfo_Require_Type_ES6_IMPORT:
					if !provider.IsEs6Module() {
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

func (g *DepsGenerator) reportMissingFile(depInfo DependencyInfo, path string) {
	name := depInfo.GetName()
	g.errorManager.Report(
		CheckLevel_ERROR.Ptr(), JSError_MakeWithSourceLocation(&name, -1, -1, DepsGenerator_UNKNOWN_PATH_IMPORT, path))
}

func (g *DepsGenerator) reportEs6ImportForNonEs6Module(provider DependencyInfo, depInfo DependencyInfo) {
	name := depInfo.GetName()
	g.errorManager.Report(
		CheckLevel_ERROR.Ptr(),
		JSError_MakeWithSourceLocation(&name, -1, -1, DepsGenerator_ES6_IMPORT_FOR_NON_ES6_MODULE, provider.GetName()))
}

func (g *DepsGenerator) reportSameFile(namespace string, depInfo DependencyInfo) {
	name := depInfo.GetName()
	g.errorManager.Report(CheckLevel_WARNING.Ptr(),
		JSError_MakeWithSourceLocation(&name, -1, -1, DepsGenerator_SAME_FILE_WARNING, namespace))
}

func (g *DepsGenerator) reportUndefinedNamespace(namespace string, depInfo DependencyInfo) {
	name := depInfo.GetName()
	g.errorManager.Report(CheckLevel_ERROR.Ptr(),
		JSError_MakeWithSourceLocation(&name, -1, -1, DepsGenerator_NEVER_PROVIDED_ERROR, namespace))
}

func (g *DepsGenerator) reportDuplicateProvide(namespace string, firstDep, secondDep DependencyInfo) {
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

func (g *DepsGenerator) reportNoDepsInDepsFile(filePath string) {
	g.errorManager.Report(CheckLevel_WARNING.Ptr(),
		JSError_MakeWithSourceLocation(&filePath, -1, -1, DepsGenerator_NO_DEPS_WARNING))
}

func (g *DepsGenerator) addProvideSymbol(provide string, depInfo DependencyInfo, providesMap map[string]DependencyInfo) {
	prevValue := providesMap[provide]
	providesMap[provide] = depInfo
	if prevValue != nil {
		g.reportDuplicateProvide(provide, prevValue, depInfo)
	}
}

func (g *DepsGenerator) addToProvideMap(depInfos []DependencyInfo, providesMap map[string]DependencyInfo, isFromDepsFile bool) {
	for _, depInfo := range depInfos {
		for _, provide := range depInfo.GetProvides() {
			g.addProvideSymbol(provide, depInfo, providesMap)
		}
		if isFromDepsFile {
			munged := g.loader.Resolve(PathUtil_MakeAbsoluteWithRoot(depInfo.GetPathRelativeToClosureBase(), g.closurePathAbs)).ToModuleName()
			g.addProvideSymbol(munged, depInfo, providesMap)
		} else if !depInfo.IsEs6Module() {
			munged := g.loader.Resolve(depInfo.GetName()).ToModuleName()
			g.addProvideSymbol(munged, depInfo, providesMap)
		}
		g.addProvideSymbol(depInfo.GetPathRelativeToClosureBase(), depInfo, providesMap)
	}
}

func (g *DepsGenerator) parseDepsFiles() *orderedDepMap {
	depsParser := NewDepsFileRegexParser(g.errorManager)
	depsParser.SetShortcutMode(true)
	depsFiles := newOrderedDepMap()
	for _, file := range g.deps {
		depInfos := depsParser.ParseFileReader(file.GetName(), file.GetCodeReader())
		if len(depInfos) == 0 {
			g.reportNoDepsInDepsFile(file.GetName())
		} else {
			for _, info := range depInfos {
				depsFiles.put(info.GetPathRelativeToClosureBase(), depsGenerator_removeRelativePathProvide(info))
			}
		}
	}
	for _, src := range g.srcs {
		for _, info := range depsParser.ParseFileReader(src.GetName(), src.GetCodeReader()) {
			depsFiles.put(info.GetPathRelativeToClosureBase(), depsGenerator_removeRelativePathProvide(info))
		}
	}
	return depsFiles
}

func depsGenerator_removeRelativePathProvide(info DependencyInfo) DependencyInfo {
	relPath := info.GetPathRelativeToClosureBase()
	origProvides := info.GetProvides()
	provides := make([]string, 0, len(origProvides))
	for _, p := range origProvides {
		if p != relPath {
			provides = append(provides, p)
		}
	}
	return newSimpleDependencyInfoDirect(
		relPath,
		info.GetName(),
		provides,
		info.GetRequires(),
		info.GetTypeRequires(),
		info.GetLoadFlags(),
		info.GetHasExternsAnnotation(),
		info.GetHasNoCompileAnnotation(),
	)
}

func (g *DepsGenerator) parseSources(preparsedFiles *orderedDepMap) *orderedDepMap {
	parsedFiles := newOrderedDepMap()
	jsParser := NewJsFileRegexParser(g.errorManager).SetModuleLoader(g.loader)
	compiler := NewCompiler()
	compiler.Init([]*SourceFile{}, []*SourceFile{}, NewCompilerOptions())

	for _, file := range g.srcs {
		closureRelativePath := PathUtil_MakeRelative(g.closurePathAbs, PathUtil_MakeAbsolute(file.GetName()))
		if depsGenerator_logger.IsLoggable(jutil.Level_FINE) {
			depsGenerator_logger.Fine("Closure-relative path: " + closureRelativePath)
		}
		if g.mergeStrategy == DepsGenerator_InclusionStrategy_WHEN_IN_SRCS || !preparsedFiles.contains(closureRelativePath) {
			var depInfo DependencyInfo = jsParser.ParseFile(file.GetName(), closureRelativePath, file.GetCode())
			depInfo = NewLazyParsedDependencyInfo(depInfo, NewCompilerInput(file), compiler)
			if depInfo.GetHasExternsAnnotation() {
				continue
			}
			file.ClearCachedSource()
			parsedFiles.put(closureRelativePath, depInfo)
		}
	}
	return parsedFiles
}

func (g *DepsGenerator) writeDepsContent(depsValues, jsValues []DependencyInfo, out *jutil.PrintStream) {
	depsGenerator_writeDepInfos(out, jsValues)
	if g.mergeStrategy == DepsGenerator_InclusionStrategy_ALWAYS {
		var depsPaths []string
		infosByDepsPath := make(map[string][]DependencyInfo)
		for _, info := range depsValues {
			name := info.GetName()
			if _, seen := infosByDepsPath[name]; !seen {
				depsPaths = append(depsPaths, name)
			}
			infosByDepsPath[name] = append(infosByDepsPath[name], info)
		}
		for _, depsPath := range depsPaths {
			out.Println("\n// Included from: " + depsPath)
			depsGenerator_writeDepInfos(out, infosByDepsPath[depsPath])
		}
	}
}

func depsGenerator_writeDepInfos(out *jutil.PrintStream, depInfos []DependencyInfo) {
	for _, depInfo := range depInfos {
		DependencyInfo_Util_WriteAddDependency(out, depInfo)
	}
}
