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
// Port of com/google/javascript/jscomp/deps/DepsGenerator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_deps_generator_1 j2gort.Ext_com_google_common_base_Function

func init() {
	lambdaSingleton_deps_generator_1 = j2gort.NewLambda_Ext_com_google_common_base_Function(func(mrp0_ any) any { return j2gort.BoxString(j2gort.AsIface[DependencyInfo](mrp0_).GetName()) })
}

// Java class: com.google.javascript.jscomp.deps.DepsGenerator
type DepsGenerator struct {
	srcs           j2gort.Ext_java_util_Collection `gson:"srcs"`
	deps           j2gort.Ext_java_util_Collection `gson:"deps"`
	closurePathAbs string                          `gson:"closurePathAbs"`
	mergeStrategy  DepsGenerator_InclusionStrategy `gson:"mergeStrategy"`
	loader         *ModuleLoader                   `gson:"loader"`
	errorManager   ErrorManager                    `gson:"errorManager"`
	idHash_        atomic.Int32
}

func (this_ *DepsGenerator) BindSelf_(s any) {
	_ = s
}
func (this_ *DepsGenerator) IsJavaNil_() bool { return this_ == nil }
func (*DepsGenerator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.DepsGenerator", "DepsGenerator"
}
func (this_ *DepsGenerator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DepsGenerator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.DepsGenerator", "com.google.javascript.jscomp.deps.DepsGenerator", func(v_ any) bool {
		return j2gort.InstanceOf[*DepsGenerator](v_)
	})
}

var static_DepsGenerator_logger *j2gort.Ext_java_util_logging_Logger

func StaticGet_DepsGenerator_Logger() *j2gort.Ext_java_util_logging_Logger {
	clinit_DepsGenerator()
	return static_DepsGenerator_logger
}
func StaticSet_DepsGenerator_Logger(v *j2gort.Ext_java_util_logging_Logger) {
	clinit_DepsGenerator()
	static_DepsGenerator_logger = v
}

var static_DepsGenerator_ES6_IMPORT_FOR_NON_ES6_MODULE *DiagnosticType

func StaticGet_DepsGenerator_ES6_IMPORT_FOR_NON_ES6_MODULE() *DiagnosticType {
	clinit_DepsGenerator()
	return static_DepsGenerator_ES6_IMPORT_FOR_NON_ES6_MODULE
}
func StaticSet_DepsGenerator_ES6_IMPORT_FOR_NON_ES6_MODULE(v *DiagnosticType) {
	clinit_DepsGenerator()
	static_DepsGenerator_ES6_IMPORT_FOR_NON_ES6_MODULE = v
}

var static_DepsGenerator_UNKNOWN_PATH_IMPORT *DiagnosticType

func StaticGet_DepsGenerator_UNKNOWN_PATH_IMPORT() *DiagnosticType {
	clinit_DepsGenerator()
	return static_DepsGenerator_UNKNOWN_PATH_IMPORT
}
func StaticSet_DepsGenerator_UNKNOWN_PATH_IMPORT(v *DiagnosticType) {
	clinit_DepsGenerator()
	static_DepsGenerator_UNKNOWN_PATH_IMPORT = v
}

var static_DepsGenerator_SAME_FILE_WARNING *DiagnosticType

func StaticGet_DepsGenerator_SAME_FILE_WARNING() *DiagnosticType {
	clinit_DepsGenerator()
	return static_DepsGenerator_SAME_FILE_WARNING
}
func StaticSet_DepsGenerator_SAME_FILE_WARNING(v *DiagnosticType) {
	clinit_DepsGenerator()
	static_DepsGenerator_SAME_FILE_WARNING = v
}

var static_DepsGenerator_NEVER_PROVIDED_ERROR *DiagnosticType

func StaticGet_DepsGenerator_NEVER_PROVIDED_ERROR() *DiagnosticType {
	clinit_DepsGenerator()
	return static_DepsGenerator_NEVER_PROVIDED_ERROR
}
func StaticSet_DepsGenerator_NEVER_PROVIDED_ERROR(v *DiagnosticType) {
	clinit_DepsGenerator()
	static_DepsGenerator_NEVER_PROVIDED_ERROR = v
}

var static_DepsGenerator_DUPE_PROVIDES_WARNING *DiagnosticType

func StaticGet_DepsGenerator_DUPE_PROVIDES_WARNING() *DiagnosticType {
	clinit_DepsGenerator()
	return static_DepsGenerator_DUPE_PROVIDES_WARNING
}
func StaticSet_DepsGenerator_DUPE_PROVIDES_WARNING(v *DiagnosticType) {
	clinit_DepsGenerator()
	static_DepsGenerator_DUPE_PROVIDES_WARNING = v
}

var static_DepsGenerator_MULTIPLE_PROVIDES_ERROR *DiagnosticType

func StaticGet_DepsGenerator_MULTIPLE_PROVIDES_ERROR() *DiagnosticType {
	clinit_DepsGenerator()
	return static_DepsGenerator_MULTIPLE_PROVIDES_ERROR
}
func StaticSet_DepsGenerator_MULTIPLE_PROVIDES_ERROR(v *DiagnosticType) {
	clinit_DepsGenerator()
	static_DepsGenerator_MULTIPLE_PROVIDES_ERROR = v
}

var static_DepsGenerator_NO_DEPS_WARNING *DiagnosticType

func StaticGet_DepsGenerator_NO_DEPS_WARNING() *DiagnosticType {
	clinit_DepsGenerator()
	return static_DepsGenerator_NO_DEPS_WARNING
}
func StaticSet_DepsGenerator_NO_DEPS_WARNING(v *DiagnosticType) {
	clinit_DepsGenerator()
	static_DepsGenerator_NO_DEPS_WARNING = v
}

var clinitDone_DepsGenerator uint32
var clinitState_DepsGenerator j2gort.ClinitState

func clinit_DepsGenerator() {
	if atomic.LoadUint32(&clinitDone_DepsGenerator) != 0 {
		return
	}
	clinitSlow_DepsGenerator()
}

func clinitSlow_DepsGenerator() {
	if !j2gort.EnterClinit(&clinitState_DepsGenerator, &clinitDone_DepsGenerator) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_DepsGenerator, &clinitDone_DepsGenerator)
	static_DepsGenerator_logger = j2gort.Ext_java_util_logging_Logger_GetLogger_1(j2gort.CheckNotNullPtr(j2gort.ClassLiteral("com.google.javascript.jscomp.deps.DepsGenerator", "DepsGenerator")).GetName())
	static_DepsGenerator_ES6_IMPORT_FOR_NON_ES6_MODULE = DiagnosticType_Warning("DEPS_ES6_IMPORT_FOR_NON_ES6_MODULE", "Cannot import file \"{0}\" because it is not an ES6 module.")
	static_DepsGenerator_UNKNOWN_PATH_IMPORT = DiagnosticType_Warning("DEPS_UNKNOWN_PATH_IMPORT", "Could not find file \"{0}\".")
	static_DepsGenerator_SAME_FILE_WARNING = DiagnosticType_Warning("DEPS_SAME_FILE", "Namespace \"{0}\" is both required and provided in the same file.")
	static_DepsGenerator_NEVER_PROVIDED_ERROR = DiagnosticType_Error("DEPS_NEVER_PROVIDED", "Namespace \"{0}\" is required but never provided.\nYou need to pass a library that has it in srcs or exports to your target''s deps.")
	static_DepsGenerator_DUPE_PROVIDES_WARNING = DiagnosticType_Warning("DEPS_DUPE_PROVIDES", "Multiple calls to goog.provide(\"{0}\")")
	static_DepsGenerator_MULTIPLE_PROVIDES_ERROR = DiagnosticType_Error("DEPS_DUPE_PROVIDES", "Namespace \"{0}\" is already provided in other file {1}")
	static_DepsGenerator_NO_DEPS_WARNING = DiagnosticType_Warning("DEPS_NO_DEPS", "No dependencies found in file")
}

func (this_ *DepsGenerator) initFields_DepsGenerator() {
	_ = this_
	this_.closurePathAbs = j2gort.NullString
}

// Java constructor: DepsGenerator(java.util.Collection<com.google.javascript.jscomp.SourceFile> deps, java.util.Collection<com.google.javascript.jscomp.SourceFile> srcs, com.google.javascript.jscomp.deps.DepsGenerator.InclusionStrategy mergeStrategy, java.lang.String closurePathAbs, com.google.javascript.jscomp.ErrorManager errorManager, com.google.javascript.jscomp.deps.ModuleLoader loader)
func NewDepsGenerator(deps j2gort.Ext_java_util_Collection, srcs j2gort.Ext_java_util_Collection, mergeStrategy DepsGenerator_InclusionStrategy, closurePathAbs string, errorManager ErrorManager, loader *ModuleLoader) *DepsGenerator {
	clinit_DepsGenerator()
	this_ := &DepsGenerator{}
	InitDepsGenerator(this_, deps, srcs, mergeStrategy, closurePathAbs, errorManager, loader)
	return this_
}

func InitDepsGenerator(this_ *DepsGenerator, deps j2gort.Ext_java_util_Collection, srcs j2gort.Ext_java_util_Collection, mergeStrategy DepsGenerator_InclusionStrategy, closurePathAbs string, errorManager ErrorManager, loader *ModuleLoader) {
	_ = deps
	_ = srcs
	_ = mergeStrategy
	_ = closurePathAbs
	_ = errorManager
	_ = loader
	this_.initFields_DepsGenerator()
	this_.deps = deps
	this_.srcs = srcs
	this_.mergeStrategy = mergeStrategy
	this_.closurePathAbs = closurePathAbs
	this_.errorManager = errorManager
	this_.loader = loader
}

// Java method: DepsGenerator.computeDependencyCalls()
func (this_ *DepsGenerator) ComputeDependencyCalls() string {
	var depsFiles j2gort.Ext_java_util_Map = this_.priv_DepsGenerator_parseDepsFiles_6097()
	_ = depsFiles
	if StaticGet_DepsGenerator_Logger().IsLoggable(j2gort.Ext_java_util_logging_Level_F_FINE) {
		StaticGet_DepsGenerator_Logger().Fine_1_string(("preparsedFiles: " + j2gort.StrVal(depsFiles)))
	}
	var jsFiles j2gort.Ext_java_util_Map = this_.priv_DepsGenerator_parseSources_6099(j2gort.AsIface[j2gort.Ext_java_util_Set](depsFiles.KeySet_0()))
	_ = jsFiles
	if this_.errorManager.GetErrorCount() > int32(0) {
		return j2gort.NullString
	}
	this_.CleanUpDuplicatedFiles(depsFiles, jsFiles)
	jsFiles = this_.priv_DepsGenerator_removeMungedSymbols_6086(depsFiles, jsFiles)
	this_.priv_DepsGenerator_validateDependencies_6087(j2gort.AsIface[j2gort.Ext_java_util_Collection](depsFiles.Values()), j2gort.AsIface[j2gort.Ext_java_util_Collection](jsFiles.Values()))
	if this_.errorManager.GetErrorCount() > int32(0) {
		return j2gort.NullString
	}
	var output *j2gort.Ext_java_io_ByteArrayOutputStream = j2gort.NewExt_java_io_ByteArrayOutputStream()
	_ = output
	this_.priv_DepsGenerator_writeDepsContent_6100(depsFiles, jsFiles, j2gort.NewExt_java_io_PrintStream_1_Ext_java_io_OutputStream(output))
	return func() string {
		b_1 := output.ToByteArray()
		return j2gort.NewStringFromBytes(b_1, 0, int32(len(b_1)), j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8)
	}()
}

// Java method: DepsGenerator.cleanUpDuplicatedFiles(java.util.Map<java.lang.String,com.google.javascript.jscomp.deps.DependencyInfo> depsFiles, java.util.Map<java.lang.String,com.google.javascript.jscomp.deps.DependencyInfo> jsFiles)
func (this_ *DepsGenerator) CleanUpDuplicatedFiles(depsFiles j2gort.Ext_java_util_Map, jsFiles j2gort.Ext_java_util_Map) {
	_ = depsFiles
	_ = jsFiles
	var depsPathsCopy j2gort.LocalOrderedSet[string] = j2gort.LocalOrderedSetFromStrCol(j2gort.AsIface[j2gort.Ext_java_util_Set](depsFiles.KeySet_0()))
	_ = depsPathsCopy
	{
		cur_2 := depsPathsCopy.Iter()
		for cur_2.HasNext() {
			path := cur_2.
				Next()
			_ = path
			if this_.mergeStrategy != StaticGet_DepsGenerator_InclusionStrategy_WHEN_IN_SRCS() {
				j2gort.MapRemoveStr(jsFiles, path)
			}
		}
	}
	{
		cur_3 := j2gort.NewMapKeyStrCursor(jsFiles)
		for cur_3.HasNext() {
			path_2 := cur_3.NextStr()
			_ = path_2
			j2gort.MapRemoveStr(depsFiles, path_2)
		}
	}
}

// Java method: DepsGenerator.removeMungedSymbols(java.util.Map<java.lang.String,com.google.javascript.jscomp.deps.DependencyInfo> depFiles, java.util.Map<java.lang.String,com.google.javascript.jscomp.deps.DependencyInfo> jsFiles)
func (this_ *DepsGenerator) priv_DepsGenerator_removeMungedSymbols_6086(depFiles j2gort.Ext_java_util_Map, jsFiles j2gort.Ext_java_util_Map) j2gort.Ext_java_util_Map {
	_ = depFiles
	_ = jsFiles
	var newJsFiles j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = newJsFiles
	var providesMap j2gort.OrderedMap[string, DependencyInfo] = j2gort.NewOrderedMap[string, DependencyInfo]()
	_ = providesMap
	this_.priv_DepsGenerator_addToProvideMap_6094(j2gort.AsIface[j2gort.Ext_java_util_Collection](depFiles.Values()), &providesMap, true)
	this_.priv_DepsGenerator_addToProvideMap_6094(j2gort.AsIface[j2gort.Ext_java_util_Collection](jsFiles.Values()), &providesMap, false)
	{
		cur_4 := j2gort.NewMapValCursor(jsFiles)
		for cur_4.HasNext() {
			rawElem_5 := cur_4.Next()
			dependencyInfo := j2gort.AsIface[DependencyInfo](rawElem_5)
			_ = dependencyInfo
			var newRequires j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
			_ = newRequires
			{
				cur_6 := j2gort.NewIterCursor(j2gort.CheckNotNull(dependencyInfo.GetRequires()))
				for cur_6.HasNext() {
					rawElem_7 := cur_6.Next()
					require := j2gort.CastPtr[DependencyInfo_Require](rawElem_7)
					_ = require
					if require.GetType() == StaticGet_DependencyInfo_Require_Type_ES6_IMPORT() {
						var provider DependencyInfo = providesMap.Get(require.GetSymbol())
						_ = provider
						if provider == nil || provider.IsJavaNil_() {
							this_.priv_DepsGenerator_reportMissingFile_6088(dependencyInfo, require.GetRawText())
						} else {
							newRequires.Add_1(j2gort.PtrToAny(require.WithSymbol(provider.GetPathRelativeToClosureBase())))
						}
					} else {
						newRequires.Add_1(j2gort.PtrToAny(require))
					}
				}
			}
			var provides []string = j2gort.ImmListToSliceNullable[string](dependencyInfo.GetProvides())
			_ = provides
			if dependencyInfo.IsEs6Module() {
				var mungedProvide string = this_.loader.Resolve(dependencyInfo.GetName()).ToModuleName()
				_ = mungedProvide
				var builder []string = []string{}
				_ = builder
				{

					for _, provide := range provides {

						_ = provide
						if !(j2gort.CheckNotNullString(provide) == mungedProvide) && !(func() bool {
							sr_9 := provide
							se_10 := dependencyInfo.GetPathRelativeToClosureBase()
							return j2gort.CheckNotNullString(sr_9) == se_10
						}()) {
							builder = append(builder, j2gort.CheckNotNullString(provide))
						}
					}
				}
				provides = builder
			}
			j2gort.MapPutStr(newJsFiles, dependencyInfo.GetPathRelativeToClosureBase(), SimpleDependencyInfo_S_Builder(dependencyInfo.GetPathRelativeToClosureBase(), dependencyInfo.GetName()).SetProvides_1_Ext_java_util_Collection(j2gort.ImmListFromTypedSlice(provides)).SetRequires_1_Ext_java_util_Collection(newRequires).SetLoadFlags(j2gort.AsIface[j2gort.Ext_java_util_Map](dependencyInfo.GetLoadFlags())).Build())
		}
	}
	return newJsFiles
}

// Java method: DepsGenerator.validateDependencies(java.lang.Iterable<com.google.javascript.jscomp.deps.DependencyInfo> preparsedFileDependencies, java.lang.Iterable<com.google.javascript.jscomp.deps.DependencyInfo> parsedFileDependencies)
func (this_ *DepsGenerator) priv_DepsGenerator_validateDependencies_6087(preparsedFileDependencies j2gort.Ext_java_lang_Iterable, parsedFileDependencies j2gort.Ext_java_lang_Iterable) {
	_ = preparsedFileDependencies
	_ = parsedFileDependencies
	var providesMap j2gort.OrderedMap[string, DependencyInfo] = j2gort.NewOrderedMap[string, DependencyInfo]()
	_ = providesMap
	this_.priv_DepsGenerator_addToProvideMap_6094(preparsedFileDependencies, &providesMap, true)
	this_.priv_DepsGenerator_addToProvideMap_6094(parsedFileDependencies, &providesMap, false)
	{
		cur_11 := j2gort.NewIterCursor(j2gort.CheckNotNull(parsedFileDependencies))
		for cur_11.HasNext() {
			rawElem_12 := cur_11.Next()
			depInfo := j2gort.AsIface[DependencyInfo](rawElem_12)
			_ = depInfo
			{
				cur_13 := j2gort.NewIterCursor(j2gort.CheckNotNull(depInfo.GetRequires()))
				for cur_13.HasNext() {
					rawElem_14 := cur_13.Next()
					require := j2gort.CastPtr[DependencyInfo_Require](rawElem_14)
					_ = require
					var namespace string = require.GetSymbol()
					_ = namespace
					var provider DependencyInfo = providesMap.Get(namespace)
					_ = provider
					if provider == nil || provider.IsJavaNil_() {
						this_.priv_DepsGenerator_reportUndefinedNamespace_6091(namespace, depInfo)
					} else if provider == depInfo || ((provider == nil || provider.IsJavaNil_()) && (depInfo == nil || depInfo.IsJavaNil_())) {
						this_.priv_DepsGenerator_reportSameFile_6090(namespace, depInfo)
					} else {
						var providerIsEs6Module bool = provider.IsEs6Module()
						_ = providerIsEs6Module
						switch (require.GetType()).Ordinal() {
						case int32(1) /* Type.ES6_IMPORT */ :
							if !(providerIsEs6Module) {
								this_.priv_DepsGenerator_reportEs6ImportForNonEs6Module_6089(provider, depInfo)
							}
						case int32(0) /* Type.GOOG_REQUIRE_SYMBOL */, int32(2) /* Type.PARSED_FROM_DEPS */ :
						case int32(3) /* Type.COMMON_JS */, int32(4) /* Type.COMPILER_CHUNK */ :
							panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected import type: "+j2gort.StrVal(require.GetType()))))
						}
					}
				}
			}
		}
	}
}

// Java method: DepsGenerator.reportMissingFile(com.google.javascript.jscomp.deps.DependencyInfo depInfo, java.lang.String path)
func (this_ *DepsGenerator) priv_DepsGenerator_reportMissingFile_6088(depInfo DependencyInfo, path string) {
	_ = depInfo
	_ = path
	this_.errorManager.Report_2(StaticGet_CheckLevel_ERROR(), JSError_Make_5(depInfo.GetName(), int32(-1), int32(-1), StaticGet_DepsGenerator_UNKNOWN_PATH_IMPORT(), []string{path}))
}

// Java method: DepsGenerator.reportEs6ImportForNonEs6Module(com.google.javascript.jscomp.deps.DependencyInfo provider, com.google.javascript.jscomp.deps.DependencyInfo depInfo)
func (this_ *DepsGenerator) priv_DepsGenerator_reportEs6ImportForNonEs6Module_6089(provider DependencyInfo, depInfo DependencyInfo) {
	_ = provider
	_ = depInfo
	this_.errorManager.Report_2(StaticGet_CheckLevel_ERROR(), JSError_Make_5(depInfo.GetName(), int32(-1), int32(-1), StaticGet_DepsGenerator_ES6_IMPORT_FOR_NON_ES6_MODULE(), []string{provider.GetName()}))
}

// Java method: DepsGenerator.reportSameFile(java.lang.String namespace, com.google.javascript.jscomp.deps.DependencyInfo depInfo)
func (this_ *DepsGenerator) priv_DepsGenerator_reportSameFile_6090(namespace string, depInfo DependencyInfo) {
	_ = namespace
	_ = depInfo
	this_.errorManager.Report_2(StaticGet_CheckLevel_WARNING(), JSError_Make_5(depInfo.GetName(), int32(-1), int32(-1), StaticGet_DepsGenerator_SAME_FILE_WARNING(), []string{namespace}))
}

// Java method: DepsGenerator.reportUndefinedNamespace(java.lang.String namespace, com.google.javascript.jscomp.deps.DependencyInfo depInfo)
func (this_ *DepsGenerator) priv_DepsGenerator_reportUndefinedNamespace_6091(namespace string, depInfo DependencyInfo) {
	_ = namespace
	_ = depInfo
	this_.errorManager.Report_2(StaticGet_CheckLevel_ERROR(), JSError_Make_5(depInfo.GetName(), int32(-1), int32(-1), StaticGet_DepsGenerator_NEVER_PROVIDED_ERROR(), []string{namespace}))
}

// Java method: DepsGenerator.reportDuplicateProvide(java.lang.String namespace, com.google.javascript.jscomp.deps.DependencyInfo firstDep, com.google.javascript.jscomp.deps.DependencyInfo secondDep)
func (this_ *DepsGenerator) priv_DepsGenerator_reportDuplicateProvide_6092(namespace string, firstDep DependencyInfo, secondDep DependencyInfo) {
	_ = namespace
	_ = firstDep
	_ = secondDep
	if firstDep == secondDep || ((firstDep == nil || firstDep.IsJavaNil_()) && (secondDep == nil || secondDep.IsJavaNil_())) {
		if !(j2gort.CheckNotNullString(firstDep.GetPathRelativeToClosureBase()) == namespace) {
			this_.errorManager.Report_2(StaticGet_CheckLevel_WARNING(), JSError_Make_5(firstDep.GetName(), int32(-1), int32(-1), StaticGet_DepsGenerator_DUPE_PROVIDES_WARNING(), []string{namespace}))
		}
	} else {
		this_.errorManager.Report_2(StaticGet_CheckLevel_ERROR(), JSError_Make_5(secondDep.GetName(), int32(-1), int32(-1), StaticGet_DepsGenerator_MULTIPLE_PROVIDES_ERROR(), []string{namespace, firstDep.GetName()}))
	}
}

// Java method: DepsGenerator.reportNoDepsInDepsFile(java.lang.String filePath)
func (this_ *DepsGenerator) priv_DepsGenerator_reportNoDepsInDepsFile_6093(filePath string) {
	_ = filePath
	this_.errorManager.Report_2(StaticGet_CheckLevel_WARNING(), JSError_Make_5(filePath, int32(-1), int32(-1), StaticGet_DepsGenerator_NO_DEPS_WARNING(), []string{}))
}

// Java method: DepsGenerator.addToProvideMap(java.lang.Iterable<com.google.javascript.jscomp.deps.DependencyInfo> depInfos, java.util.Map<java.lang.String,com.google.javascript.jscomp.deps.DependencyInfo> providesMap, boolean isFromDepsFile)
func (this_ *DepsGenerator) priv_DepsGenerator_addToProvideMap_6094(depInfos j2gort.Ext_java_lang_Iterable, providesMap *j2gort.OrderedMap[string, DependencyInfo], isFromDepsFile bool) {
	_ = depInfos
	_ = providesMap
	_ = isFromDepsFile
	{
		cur_15 := j2gort.NewIterCursor(j2gort.CheckNotNull(depInfos))
		for cur_15.HasNext() {
			rawElem_16 := cur_15.Next()
			depInfo := j2gort.AsIface[DependencyInfo](rawElem_16)
			_ = depInfo
			var provides []string = j2gort.CopyStrColToSlice(j2gort.AsIface[j2gort.Ext_java_util_Collection](depInfo.GetProvides()))
			_ = provides
			if isFromDepsFile {
				provides = append(provides, this_.loader.Resolve(PathUtil_MakeAbsolute_2(depInfo.GetPathRelativeToClosureBase(), this_.closurePathAbs)).ToModuleName())
			} else {
				if !(depInfo.IsEs6Module()) {
					provides = append(provides, this_.loader.Resolve(depInfo.GetName()).ToModuleName())
				}
			}
			provides = append(provides, depInfo.GetPathRelativeToClosureBase())
			for _, provide := range provides {
				_ = provide
				var prevValue DependencyInfo = providesMap.Put(provide, depInfo)
				_ = prevValue
				if prevValue != nil && !prevValue.IsJavaNil_() {
					this_.priv_DepsGenerator_reportDuplicateProvide_6092(provide, prevValue, depInfo)
				}
			}

		}
	}
}

// Java method: DepsGenerator.createDepsFileParser()
func (this_ *DepsGenerator) CreateDepsFileParser() *DepsFileRegexParser {
	var depsParser *DepsFileRegexParser = NewDepsFileRegexParser_1(this_.errorManager)
	_ = depsParser
	depsParser.SetShortcutMode(true)
	return depsParser
}

// Java method: DepsGenerator.shouldSkipDepsFile(com.google.javascript.jscomp.SourceFile file)
func (this_ *DepsGenerator) ShouldSkipDepsFile(file *SourceFile) bool {
	_ = file
	return false
}

// Java method: DepsGenerator.parseDepsFiles()
func (this_ *DepsGenerator) priv_DepsGenerator_parseDepsFiles_6097() j2gort.Ext_java_util_Map {
	var depsParser *DepsFileRegexParser = this_.CreateDepsFileParser()
	_ = depsParser
	var depsFiles j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = depsFiles
	{
		cur_18 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.deps))
		for cur_18.HasNext() {
			rawElem_19 := cur_18.Next()
			file := j2gort.CastPtr[SourceFile](rawElem_19)
			_ = file
			if !(this_.ShouldSkipDepsFile(file)) {
				var depInfos j2gort.Ext_java_util_List = depsParser.ParseFileReader(file.GetName(), file.GetCodeReader())
				_ = depInfos
				if depInfos.IsEmpty() {
					this_.priv_DepsGenerator_reportNoDepsInDepsFile_6093(file.GetName())
				} else {
					{
						cur_20 := j2gort.NewIterCursor(j2gort.CheckNotNull(depInfos))
						for cur_20.HasNext() {
							rawElem_21 := cur_20.Next()
							info := j2gort.AsIface[DependencyInfo](rawElem_21)
							_ = info
							j2gort.MapPutStr(depsFiles, info.GetPathRelativeToClosureBase(), this_.priv_DepsGenerator_removeRelativePathProvide_6098(info))
						}
					}
				}
			}
		}
	}
	{
		cur_22 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.srcs))
		for cur_22.HasNext() {
			rawElem_23 := cur_22.Next()
			src := j2gort.CastPtr[SourceFile](rawElem_23)
			_ = src
			if !(this_.ShouldSkipDepsFile(src)) {
				var srcInfos j2gort.Ext_java_util_List = depsParser.ParseFileReader(src.GetName(), src.GetCodeReader())
				_ = srcInfos
				{
					cur_24 := j2gort.NewIterCursor(j2gort.CheckNotNull(srcInfos))
					for cur_24.HasNext() {
						rawElem_25 := cur_24.Next()
						info_2 := j2gort.AsIface[DependencyInfo](rawElem_25)
						_ = info_2
						j2gort.MapPutStr(depsFiles, info_2.GetPathRelativeToClosureBase(), this_.priv_DepsGenerator_removeRelativePathProvide_6098(info_2))
					}
				}
			}
		}
	}
	return depsFiles
}

// Java method: DepsGenerator.removeRelativePathProvide(com.google.javascript.jscomp.deps.DependencyInfo info)
func (this_ *DepsGenerator) priv_DepsGenerator_removeRelativePathProvide_6098(info DependencyInfo) DependencyInfo {
	_ = info
	return Upcast_DependencyInfo(SimpleDependencyInfo_Builder_From(info).SetProvides_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_List](info.GetProvides().Stream().Filter(j2gort.NewLambda_Ext_java_util_function_Predicate(func(lp_26 any) bool {
		p := j2gort.CastString(lp_26)
		_ = p
		return !(func() bool {
			sr_27 := p
			se_28 := info.GetPathRelativeToClosureBase()
			return j2gort.CheckNotNullString(sr_27) == se_28
		}())
	})).Collect_1(j2gort.Ext_java_util_stream_Collectors_ToList()))).Build())
}

// Java method: DepsGenerator.parseSources(java.util.Set<java.lang.String> preparsedFiles)
func (this_ *DepsGenerator) priv_DepsGenerator_parseSources_6099(preparsedFiles j2gort.Ext_java_util_Set) j2gort.Ext_java_util_Map {
	_ = preparsedFiles
	var parsedFiles j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = parsedFiles
	var jsParser *JsFileRegexParser = NewJsFileRegexParser(this_.errorManager).SetModuleLoader(this_.loader)
	_ = jsParser
	var compiler *Compiler = NewCompiler()
	_ = compiler
	compiler.Init(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()), j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()), NewCompilerOptions())
	{
		cur_29 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.srcs))
		for cur_29.HasNext() {
			rawElem_30 := cur_29.Next()
			file := j2gort.CastPtr[SourceFile](rawElem_30)
			_ = file
			var closureRelativePath string = PathUtil_MakeRelative(this_.closurePathAbs, PathUtil_MakeAbsolute_1(file.GetName()))
			_ = closureRelativePath
			if StaticGet_DepsGenerator_Logger().IsLoggable(j2gort.Ext_java_util_logging_Level_F_FINE) {
				StaticGet_DepsGenerator_Logger().Fine_1_string(("Closure-relative path: " + j2gort.StrValString(closureRelativePath)))
			}
			if (StaticGet_DepsGenerator_InclusionStrategy_WHEN_IN_SRCS() == this_.mergeStrategy) || !(j2gort.SetContainsStr(preparsedFiles, closureRelativePath)) {
				var depInfo DependencyInfo = jsParser.ParseFile(file.GetName(), closureRelativePath, file.GetCode())
				_ = depInfo
				depInfo = NewLazyParsedDependencyInfo(depInfo, NewCompilerInput_1(file), compiler)
				if depInfo.GetHasExternsAnnotation() {
					continue
				}
				file.ClearCachedSource()
				j2gort.MapPutStr(parsedFiles, closureRelativePath, depInfo)
			}
		}
	}
	return parsedFiles
}

// Java method: DepsGenerator.writeDepsContent(java.util.Map<java.lang.String,com.google.javascript.jscomp.deps.DependencyInfo> depsFiles, java.util.Map<java.lang.String,com.google.javascript.jscomp.deps.DependencyInfo> jsFiles, java.io.PrintStream out)
func (this_ *DepsGenerator) priv_DepsGenerator_writeDepsContent_6100(depsFiles j2gort.Ext_java_util_Map, jsFiles j2gort.Ext_java_util_Map, out *j2gort.Ext_java_io_PrintStream) {
	_ = depsFiles
	_ = jsFiles
	_ = out
	DepsGenerator_WriteDepInfos(out, j2gort.AsIface[j2gort.Ext_java_util_Collection](jsFiles.Values()))
	if this_.mergeStrategy == StaticGet_DepsGenerator_InclusionStrategy_ALWAYS() {
		var infosIndex *j2gort.Ext_com_google_common_collect_ImmutableListMultimap = j2gort.Ext_com_google_common_collect_Multimaps_Index_2_Ext_java_lang_Iterable_Ext_com_google_common_base_Function(j2gort.AsIface[j2gort.Ext_java_util_Collection](depsFiles.Values()), lambdaSingleton_deps_generator_1)
		_ = infosIndex
		{
			cur_31 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](infosIndex.KeySet())))
			for cur_31.HasNext() {
				depsPath := cur_31.NextStr()
				_ = depsPath
				var path string = this_.FormatPathToDepsFile(depsPath)
				_ = path
				out.Println_1_string(("\n// Included from: " + j2gort.StrValString(path)))
				DepsGenerator_WriteDepInfos(out, j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](j2gort.MultimapGetStr(infosIndex, depsPath)))
			}
		}
	}
}

// Java method: DepsGenerator.formatPathToDepsFile(java.lang.String path)
func (this_ *DepsGenerator) FormatPathToDepsFile(path string) string {
	_ = path
	return path
}

// Java static method: DepsGenerator.writeDepInfos(java.io.PrintStream out, java.util.Collection<com.google.javascript.jscomp.deps.DependencyInfo> depInfos)
func DepsGenerator_WriteDepInfos(out *j2gort.Ext_java_io_PrintStream, depInfos j2gort.Ext_java_util_Collection) {
	clinit_DepsGenerator()
	_ = out
	_ = depInfos
	{
		cur_32 := j2gort.NewIterCursor(j2gort.CheckNotNull(depInfos))
		for cur_32.HasNext() {
			rawElem_33 := cur_32.Next()
			depInfo := j2gort.AsIface[DependencyInfo](rawElem_33)
			_ = depInfo
			DependencyInfo_Util_WriteAddDependency(j2gort.AsIface[j2gort.Ext_java_lang_Appendable](out), depInfo)
		}
	}
}

// Java static method: DepsGenerator.createSourceFilesFromPaths(java.util.Collection<java.lang.String> paths)
func DepsGenerator_CreateSourceFilesFromPaths_1_Ext_java_util_Collection(paths j2gort.Ext_java_util_Collection) j2gort.Ext_java_util_List {
	clinit_DepsGenerator()
	_ = paths
	var files []*SourceFile = []*SourceFile{}
	_ = files
	{
		cur_34 := j2gort.NewIterCursor(j2gort.CheckNotNull(paths))
		for cur_34.HasNext() {
			path := cur_34.NextStr()
			_ = path
			files = append(files, SourceFile_FromFile_1(path))
		}
	}
	return j2gort.WrapTypedArrayList(files)
}

// Java static method: DepsGenerator.createSourceFilesFromPaths(java.lang.String[] paths)
func DepsGenerator_CreateSourceFilesFromPaths_1_slice_string(paths []string) j2gort.Ext_java_util_List {
	clinit_DepsGenerator()
	_ = paths
	return DepsGenerator_CreateSourceFilesFromPaths_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_java_util_Arrays_AsList(func() []any {
		arr_35 := paths
		if arr_35 == nil {
			return nil
		}
		out_36 := make([]any, len(arr_35))
		for i_, el_37 := range arr_35 {
			out_36[i_] = j2gort.BoxString(el_37)
		}
		return out_36
	}())))
}

// Java static method: DepsGenerator.createSourceFilesFromZipPaths(java.util.Collection<java.lang.String> paths)
func DepsGenerator_CreateSourceFilesFromZipPaths(paths j2gort.Ext_java_util_Collection) j2gort.Ext_java_util_List {
	clinit_DepsGenerator()
	_ = paths
	var zipSourceFiles j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = zipSourceFiles
	{
		cur_38 := j2gort.NewIterCursor(j2gort.CheckNotNull(paths))
		for cur_38.HasNext() {
			path := cur_38.NextStr()
			_ = path
			zipSourceFiles.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](SourceFile_FromZipFile(path, j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8)))
		}
	}
	return zipSourceFiles
}

// Java enum: com.google.javascript.jscomp.deps.DepsGenerator.InclusionStrategy
type DepsGenerator_InclusionStrategy uint8

const (
	enum_DepsGenerator_InclusionStrategy_ALWAYS           DepsGenerator_InclusionStrategy = 1
	enum_DepsGenerator_InclusionStrategy_WHEN_IN_SRCS     DepsGenerator_InclusionStrategy = 2
	enum_DepsGenerator_InclusionStrategy_DO_NOT_DUPLICATE DepsGenerator_InclusionStrategy = 3
)

var enumBaseTable_DepsGenerator_InclusionStrategy = [...]j2gort.EnumBase{
	{Name_: "ALWAYS", Ordinal_: 0},
	{Name_: "WHEN_IN_SRCS", Ordinal_: 1},
	{Name_: "DO_NOT_DUPLICATE", Ordinal_: 2},
}
var enumValues_DepsGenerator_InclusionStrategy = []DepsGenerator_InclusionStrategy{1, 2, 3}

func StaticGet_DepsGenerator_InclusionStrategy_ALWAYS() DepsGenerator_InclusionStrategy {
	return enum_DepsGenerator_InclusionStrategy_ALWAYS
}
func StaticGet_DepsGenerator_InclusionStrategy_WHEN_IN_SRCS() DepsGenerator_InclusionStrategy {
	return enum_DepsGenerator_InclusionStrategy_WHEN_IN_SRCS
}
func StaticGet_DepsGenerator_InclusionStrategy_DO_NOT_DUPLICATE() DepsGenerator_InclusionStrategy {
	return enum_DepsGenerator_InclusionStrategy_DO_NOT_DUPLICATE
}

func DepsGenerator_InclusionStrategy_Values() []DepsGenerator_InclusionStrategy {
	return j2gort.CloneSlice(enumValues_DepsGenerator_InclusionStrategy)
}

func DepsGenerator_InclusionStrategy_ValueOf(name string) DepsGenerator_InclusionStrategy {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_DepsGenerator_InclusionStrategy {
		if enumBaseTable_DepsGenerator_InclusionStrategy[i_].Name_ == name {
			return DepsGenerator_InclusionStrategy(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.deps.DepsGenerator.InclusionStrategy." + name))
}

func (this_ DepsGenerator_InclusionStrategy) IsJavaNil_() bool                              { return this_ == 0 }
func (DepsGenerator_InclusionStrategy) HasIdentityEquals_()                                 {}
func (DepsGenerator_InclusionStrategy) Implements_j2gort_Ext_java_lang_Enum()               {}
func (DepsGenerator_InclusionStrategy) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ DepsGenerator_InclusionStrategy) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (DepsGenerator_InclusionStrategy) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.DepsGenerator$InclusionStrategy", "InclusionStrategy"
}
func (this_ DepsGenerator_InclusionStrategy) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_DepsGenerator_InclusionStrategy[this_-1]
}
func (this_ DepsGenerator_InclusionStrategy) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ DepsGenerator_InclusionStrategy) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ DepsGenerator_InclusionStrategy) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_DepsGenerator_InclusionStrategy[this_-1].Name_
}
func (this_ DepsGenerator_InclusionStrategy) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_DepsGenerator_InclusionStrategy[this_-1].Name_
}
func (this_ DepsGenerator_InclusionStrategy) Equals(other_ any) bool {
	o_, ok_ := other_.(DepsGenerator_InclusionStrategy)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ DepsGenerator_InclusionStrategy) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_DepsGenerator_InclusionStrategy[this_-1].HashCode()
}
func (this_ DepsGenerator_InclusionStrategy) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_DepsGenerator_InclusionStrategy[this_-1].CompareTo(other_)
}
func (this_ DepsGenerator_InclusionStrategy) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_DepsGenerator_InclusionStrategy[this_-1].Name_)
}
func init() {
	vals_ := []any{DepsGenerator_InclusionStrategy(1), DepsGenerator_InclusionStrategy(2), DepsGenerator_InclusionStrategy(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.DepsGenerator.InclusionStrategy", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.DepsGenerator$InclusionStrategy", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.DepsGenerator.InclusionStrategy", "com.google.javascript.jscomp.deps.DepsGenerator$InclusionStrategy", func(v_ any) bool {
		return j2gort.InstanceOf[DepsGenerator_InclusionStrategy](v_)
	})
}
