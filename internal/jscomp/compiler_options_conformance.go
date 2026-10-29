// Port of com/google/javascript/jscomp/CompilerOptions.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"fmt"
	"io"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
)

// Java: CompilerOptions.setPrettyPrint(boolean)
func (o *CompilerOptions) SetPrettyPrint(prettyPrint bool) {
	o.prettyPrint = prettyPrint
}

// Java: CompilerOptions.isPrettyPrint()
func (o *CompilerOptions) IsPrettyPrint() bool {
	return o.prettyPrint
}

// Java: CompilerOptions.setLineBreak(boolean)
func (o *CompilerOptions) SetLineBreak(lineBreak bool) {
	o.lineBreak = lineBreak
}

// Java: CompilerOptions.shouldAddLineBreak()
func (o *CompilerOptions) ShouldAddLineBreak() bool {
	return o.lineBreak
}

// Java: CompilerOptions.setPrintInputDelimiter(boolean)
func (o *CompilerOptions) SetPrintInputDelimiter(printInputDelimiter bool) {
	o.printInputDelimiter = printInputDelimiter
}

// Java: CompilerOptions.shouldPrintInputDelimiter()
func (o *CompilerOptions) ShouldPrintInputDelimiter() bool {
	return o.printInputDelimiter
}

// Java: CompilerOptions.setInputDelimiter(String)
func (o *CompilerOptions) SetInputDelimiter(inputDelimiter string) {
	o.inputDelimiter = inputDelimiter
}

// Java: CompilerOptions.getInputDelimiter()
func (o *CompilerOptions) GetInputDelimiter() string {
	return o.inputDelimiter
}

// Java: CompilerOptions.setDebugLogDirectory(Path)
func (o *CompilerOptions) SetDebugLogDirectory(dir *jutil.Path) {
	o.debugLogDirectory = dir
}

// Java: CompilerOptions.getDebugLogDirectory()
func (o *CompilerOptions) GetDebugLogDirectory() *jutil.Path {
	return o.debugLogDirectory
}

// Java: CompilerOptions.setDebugLogFilter(String)
func (o *CompilerOptions) SetDebugLogFilter(filter *string) {
	o.debugLogFilter = filter
}

// Java: CompilerOptions.getDebugLogFilter()
func (o *CompilerOptions) GetDebugLogFilter() *string {
	return o.debugLogFilter
}

// Java: CompilerOptions.shouldSerializeExtraDebugInfo()
func (o *CompilerOptions) ShouldSerializeExtraDebugInfo() bool {
	// NOTE: this is tied to debug logging (via checking getDebugLogDirectory()) for convenience,
	// and it's possible to split up the two options if this begins to trigger memory issues.
	return o.serializeExtraDebugInfo || o.GetDebugLogDirectory() != nil
}

// Java: CompilerOptions.getStateCompressionWrapper()
func (o *CompilerOptions) GetStateCompressionWrapper() func(io.Writer) io.Writer {
	return o.stateCompressionWrapper
}

// Java: CompilerOptions.getStateDecompressionWrapper()
func (o *CompilerOptions) GetStateDecompressionWrapper() func(io.Reader) io.Reader {
	return o.stateDecompressionWrapper
}

// Java: CompilerOptions.setQuoteKeywordProperties(boolean)
func (o *CompilerOptions) SetQuoteKeywordProperties(quoteKeywordProperties bool) {
	o.quoteKeywordProperties = quoteKeywordProperties
}

// Java: CompilerOptions.shouldQuoteKeywordProperties()
func (o *CompilerOptions) ShouldQuoteKeywordProperties() bool {
	// Never quote properties in .i.js files
	if o.incrementalCheckMode == CompilerOptions_IncrementalCheckMode_GENERATE_IJS {
		return false
	}
	return o.quoteKeywordProperties || parser.FeatureSet_ES3.Contains(o.GetOutputFeatureSet())
}

// Java: CompilerOptions.setErrorFormat(ErrorFormat)
func (o *CompilerOptions) SetErrorFormat(errorFormat ErrorFormat) {
	o.errorFormat = errorFormat
}

// Java: CompilerOptions.getErrorFormat()
func (o *CompilerOptions) GetErrorFormat() ErrorFormat {
	return o.errorFormat
}

// Java: CompilerOptions.setLineLengthThreshold(int)
func (o *CompilerOptions) SetLineLengthThreshold(lineLengthThreshold int) {
	o.lineLengthThreshold = lineLengthThreshold
}

// Java: CompilerOptions.getLineLengthThreshold()
func (o *CompilerOptions) GetLineLengthThreshold() int {
	return o.lineLengthThreshold
}

// Java: CompilerOptions.setUseOriginalNamesInOutput(boolean)
func (o *CompilerOptions) SetUseOriginalNamesInOutput(useOriginalNamesInOutput bool) {
	o.useOriginalNamesInOutput = useOriginalNamesInOutput
}

// Java: CompilerOptions.getUseOriginalNamesInOutput()
func (o *CompilerOptions) GetUseOriginalNamesInOutput() bool {
	return o.useOriginalNamesInOutput
}

// Java: CompilerOptions.setExternExportsPath(String)
func (o *CompilerOptions) SetExternExportsPath(externExportsPath *string) {
	o.externExportsPath = externExportsPath
}

// Java: CompilerOptions.getExternExportsPath()
func (o *CompilerOptions) GetExternExportsPath() *string {
	return o.externExportsPath
}

// Java: CompilerOptions.setSourceMapOutputPath(String)
func (o *CompilerOptions) SetSourceMapOutputPath(sourceMapOutputPath *string) {
	o.sourceMapOutputPath = sourceMapOutputPath
}

// Java: CompilerOptions.setApplyInputSourceMaps(boolean)
func (o *CompilerOptions) SetApplyInputSourceMaps(applyInputSourceMaps bool) {
	o.applyInputSourceMaps = applyInputSourceMaps
}

// Java: CompilerOptions.getApplyInputSourceMaps()
func (o *CompilerOptions) GetApplyInputSourceMaps() bool {
	return o.applyInputSourceMaps
}

// Java: CompilerOptions.getResolveSourceMapAnnotations()
func (o *CompilerOptions) GetResolveSourceMapAnnotations() bool {
	return o.resolveSourceMapAnnotations
}

// Java: CompilerOptions.setSourceMapIncludeSourcesContent(boolean)
func (o *CompilerOptions) SetSourceMapIncludeSourcesContent(sourceMapIncludeSourcesContent bool) {
	o.sourceMapIncludeSourcesContent = sourceMapIncludeSourcesContent
}

// Java: CompilerOptions.getSourceMapIncludeSourcesContent()
func (o *CompilerOptions) GetSourceMapIncludeSourcesContent() bool {
	return o.sourceMapIncludeSourcesContent
}

// Java: CompilerOptions.setParseInlineSourceMaps(boolean)
func (o *CompilerOptions) SetParseInlineSourceMaps(parseInlineSourceMaps bool) {
	o.parseInlineSourceMaps = parseInlineSourceMaps
}

// Java: CompilerOptions.getParseInlineSourceMaps()
func (o *CompilerOptions) GetParseInlineSourceMaps() bool {
	return o.parseInlineSourceMaps
}

// Java: CompilerOptions.setSourceMapDetailLevel(SourceMap.DetailLevel)
func (o *CompilerOptions) SetSourceMapDetailLevel(sourceMapDetailLevel SourceMap_DetailLevel) {
	o.sourceMapDetailLevel = sourceMapDetailLevel
}

// Java: CompilerOptions.getSourceMapDetailLevel()
func (o *CompilerOptions) GetSourceMapDetailLevel() SourceMap_DetailLevel {
	return o.sourceMapDetailLevel
}

// Java: CompilerOptions.setSourceMapFormat(SourceMap.Format)
func (o *CompilerOptions) SetSourceMapFormat(sourceMapFormat SourceMap_Format) {
	o.sourceMapFormat = sourceMapFormat
}

// Java: CompilerOptions.getSourceMapFormat()
func (o *CompilerOptions) GetSourceMapFormat() SourceMap_Format {
	return o.sourceMapFormat
}

// Java: CompilerOptions.setSourceMapLocationMappings(List<? extends SourceMap.LocationMapping>)
func (o *CompilerOptions) SetSourceMapLocationMappings(sourceMapLocationMappings []SourceMap_LocationMapping) {
	o.sourceMapLocationMappings = append([]SourceMap_LocationMapping{}, sourceMapLocationMappings...)
}

// Java: CompilerOptions.getSourceMapLocationMappings()
func (o *CompilerOptions) GetSourceMapLocationMappings() []SourceMap_LocationMapping {
	return o.sourceMapLocationMappings
}

// Rewrites CommonJS modules so that modules can be concatenated together, by renaming all globals
// to avoid conflicting with other modules.
//
// Java: CompilerOptions.setProcessCommonJSModules(boolean)
func (o *CompilerOptions) SetProcessCommonJSModules(processCommonJSModules bool) {
	o.processCommonJSModules = processCommonJSModules
}

// Java: CompilerOptions.getProcessCommonJSModules()
func (o *CompilerOptions) GetProcessCommonJSModules() bool {
	return o.processCommonJSModules
}

// Java: CompilerOptions.setEs6ModuleTranspilation(Es6ModuleTranspilation)
func (o *CompilerOptions) SetEs6ModuleTranspilation(value CompilerOptions_Es6ModuleTranspilation) {
	o.es6ModuleTranspilation = value
}

// Java: CompilerOptions.getEs6ModuleTranspilation()
func (o *CompilerOptions) GetEs6ModuleTranspilation() CompilerOptions_Es6ModuleTranspilation {
	return o.es6ModuleTranspilation
}

// Sets a path prefix for CommonJS modules (maps to #setModuleRoots(List)).
//
// Java: CompilerOptions.setCommonJSModulePathPrefix(String)
func (o *CompilerOptions) SetCommonJSModulePathPrefix(commonJSModulePathPrefix string) {
	o.SetModuleRoots([]string{commonJSModulePathPrefix})
}

// Sets the module roots.
//
// Java: CompilerOptions.setModuleRoots(List<String>)
func (o *CompilerOptions) SetModuleRoots(moduleRoots []string) {
	o.moduleRoots = append([]string{}, moduleRoots...)
}

// Java: CompilerOptions.getModuleRoots()
func (o *CompilerOptions) GetModuleRoots() []string {
	return o.moduleRoots
}

// Sets whether to rewrite polyfills.
//
// Java: CompilerOptions.setRewritePolyfills(boolean)
func (o *CompilerOptions) SetRewritePolyfills(rewritePolyfills bool) {
	o.rewritePolyfills = rewritePolyfills
}

// Java: CompilerOptions.getRewritePolyfills()
func (o *CompilerOptions) GetRewritePolyfills() bool {
	return o.rewritePolyfills
}

// Sets whether to isolate polyfills from the global scope.
//
// Java: CompilerOptions.setIsolatePolyfills(boolean)
func (o *CompilerOptions) SetIsolatePolyfills(isolatePolyfills bool) {
	o.isolatePolyfills = isolatePolyfills
	if o.isolatePolyfills {
		o.SetDefineToBooleanLiteral("$jscomp.ISOLATE_POLYFILLS", isolatePolyfills)
	}
}

// Java: CompilerOptions.getIsolatePolyfills()
func (o *CompilerOptions) GetIsolatePolyfills() bool {
	return o.isolatePolyfills
}

// Java: CompilerOptions.setInjectPolyfillsNewerThan(LanguageMode)
func (o *CompilerOptions) SetInjectPolyfillsNewerThan(injectPolyfillsNewerThan *CompilerOptions_LanguageMode) {
	o.injectPolyfillsNewerThan = injectPolyfillsNewerThan
}

// Java: CompilerOptions.getInjectPolyfillsNewerThan()
func (o *CompilerOptions) GetInjectPolyfillsNewerThan() *CompilerOptions_LanguageMode {
	return o.injectPolyfillsNewerThan
}

// Sets whether to transpile async functions and generators for AsyncContext.
//
// Java: CompilerOptions.setInstrumentAsyncContext(boolean)
func (o *CompilerOptions) SetInstrumentAsyncContext(instrumentAsyncContext bool) {
	o.instrumentAsyncContext = instrumentAsyncContext
	o.SetDefineToBooleanLiteral("$jscomp.INSTRUMENT_ASYNC_CONTEXT", instrumentAsyncContext)
}

// Java: CompilerOptions.getInstrumentAsyncContext()
func (o *CompilerOptions) GetInstrumentAsyncContext() bool {
	return o.instrumentAsyncContext
}

// Java: CompilerOptions.setEs6SubclassTranspilation(Es6SubclassTranspilation)
func (o *CompilerOptions) SetEs6SubclassTranspilation(es6SubclassTranspilation CompilerOptions_Es6SubclassTranspilation) {
	o.es6SubclassTranspilation = es6SubclassTranspilation
}

// Java: CompilerOptions.getEs6SubclassTranspilation()
func (o *CompilerOptions) GetEs6SubclassTranspilation() CompilerOptions_Es6SubclassTranspilation {
	return o.es6SubclassTranspilation
}

// Sets list of libraries to always inject, even if not needed.
//
// Java: CompilerOptions.setForceLibraryInjection(Iterable<String>)
func (o *CompilerOptions) SetForceLibraryInjection(libraries []string) {
	o.forceLibraryInjection = append([]string{}, libraries...)
}

// Java: CompilerOptions.getForceLibraryInjectionList()
func (o *CompilerOptions) GetForceLibraryInjectionList() []string {
	return o.forceLibraryInjection
}

// Sets the set of libraries to never inject, even if required.
//
// Java: CompilerOptions.setPreventLibraryInjection(boolean)
func (o *CompilerOptions) SetPreventLibraryInjection(preventLibraryInjection bool) {
	if preventLibraryInjection {
		o.runtimeLibraryMode = RuntimeJsLibManager_RuntimeLibraryMode_NO_OP
	} else {
		o.runtimeLibraryMode = RuntimeJsLibManager_RuntimeLibraryMode_INJECT
	}
}

// Java: CompilerOptions.setRuntimeLibraryMode(RuntimeLibraryMode)
func (o *CompilerOptions) SetRuntimeLibraryMode(runtimeLibraryMode RuntimeJsLibManager_RuntimeLibraryMode) {
	o.runtimeLibraryMode = runtimeLibraryMode
}

// Java: CompilerOptions.getRuntimeLibraryMode()
func (o *CompilerOptions) GetRuntimeLibraryMode() RuntimeJsLibManager_RuntimeLibraryMode {
	return o.runtimeLibraryMode
}

// Java: CompilerOptions.getUnusedImportsToRemove()
func (o *CompilerOptions) GetUnusedImportsToRemove() *jutil.LinkedHashSet[string] {
	return o.unusedImportsToRemove
}

// Java: CompilerOptions.setInstrumentForCoverageOption(InstrumentOption)
func (o *CompilerOptions) SetInstrumentForCoverageOption(instrumentForCoverageOption CompilerOptions_InstrumentOption) {
	o.instrumentForCoverageOption = instrumentForCoverageOption
}

// Java: CompilerOptions.getInstrumentForCoverageOption()
func (o *CompilerOptions) GetInstrumentForCoverageOption() CompilerOptions_InstrumentOption {
	return o.instrumentForCoverageOption
}

// Sets the name for the global array which is used by PRODUCTION instrumentation. The array is
// declared during the instrumentation pass with the name provided through this setter.
//
// Java: CompilerOptions.setProductionInstrumentationArrayName(String)
func (o *CompilerOptions) SetProductionInstrumentationArrayName(productionInstrumentationArrayName string) {
	o.productionInstrumentationArrayName = productionInstrumentationArrayName
}

// Java: CompilerOptions.getProductionInstrumentationArrayName()
func (o *CompilerOptions) GetProductionInstrumentationArrayName() string {
	return o.productionInstrumentationArrayName
}

// Java: CompilerOptions.getConformanceConfigs()
func (o *CompilerOptions) GetConformanceConfigs() []*ConformanceConfig {
	return o.conformanceConfigs
}

// Both enable and configure conformance checks, if non-null.
//
// @deprecated See go/binary-level-conformance-deprecated.
//
// Java: CompilerOptions.setConformanceConfigs(List<ConformanceConfig>)
func (o *CompilerOptions) SetConformanceConfigs(configs []*ConformanceConfig) {
	for _, c := range configs {
		jutil.CheckNotNull(c) // ImmutableList.copyOf rejects null elements.
	}
	o.conformanceConfigs = append([]*ConformanceConfig{}, configs...)
}

// Whether the output should contain a 'use strict' directive.
//
// Java: CompilerOptions.shouldEmitUseStrict()
func (o *CompilerOptions) ShouldEmitUseStrict() bool {
	// Fall back to the language in's strictness if there is no output language explicitly set
	// for backwards compatibility.
	if o.emitUseStrict != base.Tri_UNKNOWN {
		return o.emitUseStrict == base.Tri_TRUE
	}
	if o.languageOutIsDefaultStrict != base.Tri_UNKNOWN {
		return o.languageOutIsDefaultStrict == base.Tri_TRUE
	}
	return o.languageIn.IsDefaultStrict()
}

// Java: CompilerOptions.setEmitUseStrict(boolean)
func (o *CompilerOptions) SetEmitUseStrict(emitUseStrict bool) *CompilerOptions {
	o.emitUseStrict = base.Tri_ForBoolean(emitUseStrict)
	return o
}

// Java: CompilerOptions.getModuleResolutionMode()
func (o *CompilerOptions) GetModuleResolutionMode() ModuleLoader_ResolutionMode {
	return o.moduleResolutionMode
}

// Java: CompilerOptions.setModuleResolutionMode(ResolutionMode)
func (o *CompilerOptions) SetModuleResolutionMode(moduleResolutionMode ModuleLoader_ResolutionMode) {
	o.moduleResolutionMode = moduleResolutionMode
}

// Java: CompilerOptions.getBrowserResolverPrefixReplacements()
func (o *CompilerOptions) GetBrowserResolverPrefixReplacements() *jutil.LinkedHashMap[string, string] {
	return o.browserResolverPrefixReplacements
}

// Java: CompilerOptions.setBrowserResolverPrefixReplacements(ImmutableMap<String, String>)
func (o *CompilerOptions) SetBrowserResolverPrefixReplacements(browserResolverPrefixReplacements *jutil.LinkedHashMap[string, string]) {
	o.browserResolverPrefixReplacements = browserResolverPrefixReplacements
}

// Java: CompilerOptions.setPathEscaper(ModuleLoader.PathEscaper)
func (o *CompilerOptions) SetPathEscaper(pathEscaper ModuleLoader_PathEscaper) {
	o.pathEscaper = pathEscaper
}

// Java: CompilerOptions.getPathEscaper()
func (o *CompilerOptions) GetPathEscaper() ModuleLoader_PathEscaper {
	return o.pathEscaper
}

// Java: CompilerOptions.getPackageJsonEntryNames()
func (o *CompilerOptions) GetPackageJsonEntryNames() []string {
	return o.packageJsonEntryNames
}

// Java: CompilerOptions.setPackageJsonEntryNames(List<String>)
func (o *CompilerOptions) SetPackageJsonEntryNames(names []string) {
	o.packageJsonEntryNames = names
}

// Java: CompilerOptions.setUseSizeHeuristicToStopOptimizationLoop(boolean)
func (o *CompilerOptions) SetUseSizeHeuristicToStopOptimizationLoop(mayStopEarly bool) {
	o.useSizeHeuristicToStopOptimizationLoop = mayStopEarly
}

// Java: CompilerOptions.shouldUseSizeHeuristicToStopOptimizationLoop()
func (o *CompilerOptions) ShouldUseSizeHeuristicToStopOptimizationLoop() bool {
	return o.useSizeHeuristicToStopOptimizationLoop
}

// Java: CompilerOptions.setMaxOptimizationLoopIterations(int)
func (o *CompilerOptions) SetMaxOptimizationLoopIterations(maxIterations int) {
	o.optimizationLoopMaxIterations = maxIterations
}

// Java: CompilerOptions.getMaxOptimizationLoopIterations()
func (o *CompilerOptions) GetMaxOptimizationLoopIterations() int {
	return o.optimizationLoopMaxIterations
}

// Java: CompilerOptions.getChunkOutputType()
func (o *CompilerOptions) GetChunkOutputType() CompilerOptions_ChunkOutputType {
	return o.chunkOutputType
}

// Java: CompilerOptions.setChunkOutputType(ChunkOutputType)
func (o *CompilerOptions) SetChunkOutputType(chunkOutputType CompilerOptions_ChunkOutputType) {
	o.chunkOutputType = chunkOutputType
}

// Java: CompilerOptions.getStrictMessageReplacement()
func (o *CompilerOptions) GetStrictMessageReplacement() bool {
	return o.strictMessageReplacement
}

// Java: CompilerOptions.toString()
func (o *CompilerOptions) String() string {
	h := jutil.MoreObjects_ToStringHelper("CompilerOptions").OmitNullValues()
	h.Add("aliasStringsMode", compilerOptions_valueOf(o.GetAliasStringsMode()))
	h.AddBool("ambiguateProperties", o.ambiguateProperties)
	h.AddBool("angularPass", o.angularPass)
	h.AddBool("assumeGettersArePure", o.assumeGettersArePure)
	h.AddBool("assumePropertiesAreStaticallyAnalyzable", o.assumePropertiesAreStaticallyAnalyzable)
	h.AddBool("assumeStrictThis", o.AssumeStrictThis())
	h.Add("browserResolverPrefixReplacements", compilerOptions_valueOf(o.browserResolverPrefixReplacements))
	h.AddBool("checkDeterminism", o.GetCheckDeterminism())
	h.AddBool("checkSuspiciousCode", o.checkSuspiciousCode)
	h.AddBool("checkSymbols", o.checkSymbols)
	h.AddBool("checkTypes", o.checkTypes)
	h.AddBool("checksOnly", o.checksOnly)
	h.Add("chunksToPrintAfterEachPassRegexList", compilerOptions_valueOf(o.chunksToPrintAfterEachPassRegexList))
	h.AddBool("closurePass", o.closurePass)
	h.AddBool("coalesceVariableNames", o.coalesceVariableNames)
	h.Add("codingConvention", compilerOptions_valueOf(o.GetCodingConvention()))
	h.AddBool("collapseAnonymousFunctions", o.collapseAnonymousFunctions)
	h.AddBool("collapseObjectLiterals", o.collapseObjectLiterals)
	h.Add("collapseProperties", compilerOptions_valueOf(o.collapsePropertiesLevel))
	h.AddBool("collapseVariableDeclarations", o.collapseVariableDeclarations)
	h.AddBool("colorizeErrorOutput", o.ShouldColorizeErrorOutput())
	h.AddBool("computeFunctionSideEffects", o.computeFunctionSideEffects)
	h.Add("conformanceConfigs", compilerOptions_valueOf(o.GetConformanceConfigs()))
	h.Add("conformanceRemoveRegexFromPath", compilerOptions_valueOf(compilerOptions_optionalOf(o.conformanceRemoveRegexFromPath)))
	h.Add("conformanceReportingMode", compilerOptions_valueOf(o.conformanceReportingMode))
	h.AddBool("continueAfterErrors", o.CanContinueAfterErrors())
	h.AddBool("convertToDottedProperties", o.convertToDottedProperties)
	h.AddBool("crossChunkCodeMotion", o.crossChunkCodeMotion)
	h.AddBool("crossChunkCodeMotionNoStubMethods", o.crossChunkCodeMotionNoStubMethods)
	h.AddBool("crossChunkMethodMotion", o.crossChunkMethodMotion)
	h.Add("cssRenamingMap", compilerOptions_valueOf(o.cssRenamingMap))
	h.Add("cssRenamingSkiplist", compilerOptions_valueOf(o.cssRenamingSkiplist))
	h.Add("customPasses", compilerOptions_valueOf(o.customPasses))
	h.AddBool("deadAssignmentElimination", o.deadAssignmentElimination)
	h.Add("debugLogDirectory", compilerOptions_valueOf(o.debugLogDirectory))
	h.Add("defineReplacements", compilerOptions_valueOf(o.GetDefineReplacements()))
	h.Add("dependencyOptions", compilerOptions_valueOf(o.GetDependencyOptions()))
	h.Add("devMode", compilerOptions_valueOf(o.devMode))
	h.AddBool("devirtualizeMethods", o.devirtualizeMethods)
	h.AddBool("disambiguateProperties", o.disambiguateProperties)
	h.Add("emitUseStrict", compilerOptions_valueOf(compilerOptions_optionalOfTri(o.emitUseStrict)))
	h.AddBool("enableModuleRewriting", o.enableModuleRewriting)
	h.Add("environment", compilerOptions_valueOf(o.GetEnvironment()))
	h.Add("errorFormat", compilerOptions_valueOf(o.errorFormat))
	h.Add("errorHandler", compilerOptions_valueOf(o.errorHandler))
	h.Add("es6ModuleTranspilation", compilerOptions_valueOf(o.es6ModuleTranspilation))
	h.Add("es6SubclassTranspilation", compilerOptions_valueOf(o.es6SubclassTranspilation))
	h.AddBool("exportLocalPropertyDefinitions", o.exportLocalPropertyDefinitions)
	h.AddBool("exportTestFunctions", o.exportTestFunctions)
	h.Add("externExportsPath", compilerOptions_valueOf(o.externExportsPath))
	h.Add("extraAnnotationNames", compilerOptions_valueOf(o.extraAnnotationNames))
	h.Add("extractPrototypeMemberDeclarations", compilerOptions_valueOf(o.extractPrototypeMemberDeclarations))
	h.Add("filesToPrintAfterEachPassRegexList", compilerOptions_valueOf(o.filesToPrintAfterEachPassRegexList))
	h.AddBool("flowSensitiveInlineVariables", o.flowSensitiveInlineVariables)
	h.AddBool("foldConstants", o.foldConstants)
	h.Add("forceLibraryInjection", compilerOptions_valueOf(o.forceLibraryInjection))
	h.AddBool("gatherCssNames", o.gatherCssNames)
	h.AddBool("generateExports", o.generateExports)
	h.AddBool("generatePseudoNames", o.generatePseudoNames)
	h.AddBool("generateTypedExterns", o.ShouldGenerateTypedExterns())
	h.Add("idGenerators", compilerOptions_valueOf(o.idGenerators))
	h.Add("idGeneratorsMapSerialized", compilerOptions_valueOf(o.idGeneratorsMapSerialized))
	h.Add("incrementalCheckMode", compilerOptions_valueOf(o.incrementalCheckMode))
	h.AddBool("inferConsts", o.inferConsts)
	h.AddBool("inferTypes", o.inferTypes)
	h.AddBool("inlineConstantVars", o.inlineConstantVars)
	h.Add("inlineFunctionsLevel", compilerOptions_valueOf(o.inlineFunctionsLevel))
	h.AddBool("inlineGetters", o.inlineGetters)
	h.AddBool("inlineLocalVariables", o.inlineLocalVariables)
	h.AddBool("inlineProperties", o.inlineProperties)
	h.AddBool("inlineVariables", o.inlineVariables)
	h.Add("inputDelimiter", compilerOptions_valueOf(o.inputDelimiter))
	h.Add("inputPropertyMap", compilerOptions_valueOf(o.inputPropertyMap))
	h.Add("inputSourceMaps", compilerOptions_valueOf(o.inputSourceMaps))
	h.Add("inputVariableMap", compilerOptions_valueOf(o.inputVariableMap))
	h.AddBool("instrumentForCoverageOnly", o.instrumentForCoverageOnly)
	h.Add("instrumentForCoverageOption", compilerOptions_valueOf(o.instrumentForCoverageOption.String()))
	h.AddBool("isolatePolyfills", o.isolatePolyfills)
	h.AddBool("j2clMinifierEnabled", o.j2clMinifierEnabled)
	h.Add("j2clMinifierPruningManifest", compilerOptions_valueOf(o.j2clMinifierPruningManifest))
	h.Add("j2clPassMode", compilerOptions_valueOf(o.j2clPassMode))
	h.AddBool("labelRenaming", o.labelRenaming)
	h.Add("languageIn", compilerOptions_valueOf(o.GetLanguageIn()))
	h.Add("languageOutIsDefaultStrict", compilerOptions_valueOf(compilerOptions_optionalOfTri(o.languageOutIsDefaultStrict)))
	h.AddBool("lineBreak", o.lineBreak)
	h.Add("lineLengthThreshold", o.lineLengthThreshold) // add(String, int): never omitted
	h.Add("locale", compilerOptions_valueOf(o.locale))
	h.AddBool("markAsCompiled", o.markAsCompiled)
	h.Add("maxFunctionSizeAfterInlining", o.maxFunctionSizeAfterInlining) // add(String, int): never omitted
	h.Add("messageBundle", compilerOptions_valueOf(o.messageBundle))
	h.Add("moduleRoots", compilerOptions_valueOf(o.moduleRoots))
	h.Add("nameGenerator", compilerOptions_valueOf(o.nameGenerator))
	h.Add("numParallelThreads", o.numParallelThreads) // add(String, int): never omitted
	h.AddBool("optimizeCalls", o.optimizeCalls)
	h.AddBool("optimizeESClassConstructors", o.optimizeESClassConstructors)
	h.Add(
		"optimizeLocalAccessForGlobalSymbolNamespace",
		compilerOptions_valueOf(o.optimizeLocalAccessForGlobalSymbolNamespace))
	h.Add("outputCharset", compilerOptions_valueOf(o.outputCharset))
	h.Add("outputFeatureSet", compilerOptions_valueOf(compilerOptions_optionalOf(o.outputFeatureSet)))
	h.Add("outputJs", compilerOptions_valueOf(o.outputJs))
	h.AddBool("outputJsStringUsage", o.outputJsStringUsage)
	h.AddBool("parentChunkCanSeeSymbolsDeclaredInChildren", o.parentChunkCanSeeSymbolsDeclaredInChildren)
	h.Add("parseJsDocDocumentation", compilerOptions_valueOf(o.IsParseJsDocDocumentation()))
	h.Add("pathEscaper", compilerOptions_valueOf(o.pathEscaper))
	h.AddBool("polymerPass", o.polymerPass)
	h.AddBool("preferSingleQuotes", o.preferSingleQuotes)
	h.AddBool("preferStableNames", o.preferStableNames)
	h.AddBool("preserveDetailedSourceInfo", o.PreservesDetailedSourceInfo())
	h.AddBool("preserveGoogProvidesAndRequires", o.preserveClosurePrimitives)
	h.AddBool("preserveNonJSDocComments", o.GetPreserveNonJSDocComments())
	h.AddBool("preserveTypeAnnotations", o.preserveTypeAnnotations)
	h.AddBool("prettyPrint", o.prettyPrint)
	h.AddBool("printConfig", o.printConfig)
	h.AddBool("printInputDelimiter", o.printInputDelimiter)
	h.AddBool("printSourceAfterEachPass", o.printSourceAfterEachPass)
	h.AddBool("processCommonJSModules", o.processCommonJSModules)
	h.Add("productionInstrumentationArrayName", compilerOptions_valueOf(o.productionInstrumentationArrayName))
	h.Add("propertiesThatMustDisambiguate", compilerOptions_valueOf(o.propertiesThatMustDisambiguate))
	h.Add("propertyRenaming", compilerOptions_valueOf(o.propertyRenaming))
	h.AddBool("propertyRenamingOnlyCompilationMode", o.propertyRenamingOnlyCompilationMode)
	h.AddBool("protectHiddenSideEffects", o.protectHiddenSideEffects)
	h.Add("qnameUsesToPrintAfterEachPassRegexList", compilerOptions_valueOf(o.qnameUsesToPrintAfterEachPassList))
	h.AddBool("quoteKeywordProperties", o.quoteKeywordProperties)
	h.AddBool("removeAbstractMethods", o.removeAbstractMethods)
	h.AddBool("removeClosureAsserts", o.removeClosureAsserts)
	h.AddBool("removeJ2clAsserts", o.removeJ2clAsserts)
	h.AddBool("removeUnusedClassProperties", o.removeUnusedClassProperties)
	h.AddBool("removeUnusedLocalVars", o.removeUnusedLocalVars)
	h.AddBool("removeUnusedPrototypeProperties", o.removeUnusedPrototypeProperties)
	h.AddBool("removeUnusedVars", o.removeUnusedVars)
	h.Add("renamePrefix", compilerOptions_valueOf(o.renamePrefix))
	h.Add("renamePrefixNamespace", compilerOptions_valueOf(o.renamePrefixNamespace))
	h.AddBool("renamePrefixNamespaceAssumeCrossChunkNames", o.renamePrefixNamespaceAssumeCrossChunkNames)
	h.AddBool("replaceIdGenerators", o.replaceIdGenerators)
	h.AddBool("replaceMessagesWithChromeI18n", o.replaceMessagesWithChromeI18n)
	h.Add("replaceStringsFunctionDescriptions", compilerOptions_valueOf(o.replaceStringsFunctionDescriptions))
	h.Add("replaceStringsPlaceholderToken", compilerOptions_valueOf(o.replaceStringsPlaceholderToken))
	h.AddBool("reserveRawExports", o.reserveRawExports)
	h.AddBool("rewriteFunctionExpressions", o.rewriteFunctionExpressions)
	h.AddBool("rewriteGlobalDeclarationsForTryCatchWrapping", o.rewriteGlobalDeclarationsForTryCatchWrapping)
	h.AddBool("rewriteModulesBeforeTypechecking", o.rewriteModulesBeforeTypechecking)
	h.AddBool("rewritePolyfills", o.rewritePolyfills)
	h.Add("runtimeLibraryMode", compilerOptions_valueOf(o.runtimeLibraryMode))
	h.AddBool("skipNonTranspilationPasses", o.skipNonTranspilationPasses)
	h.AddBool("smartNameRemoval", o.smartNameRemoval)
	h.Add("sourceMapDetailLevel", compilerOptions_valueOf(o.sourceMapDetailLevel))
	h.Add("sourceMapFormat", compilerOptions_valueOf(o.sourceMapFormat))
	h.Add("sourceMapLocationMappings", compilerOptions_valueOf(o.sourceMapLocationMappings))
	h.Add("sourceMapOutputPath", compilerOptions_valueOf(o.sourceMapOutputPath))
	h.AddBool("strictMessageReplacement", o.strictMessageReplacement)
	h.Add("stripNamePrefixes", compilerOptions_valueOf(o.stripNamePrefixes))
	h.Add("stripNameSuffixes", compilerOptions_valueOf(o.stripNameSuffixes))
	h.Add("stripTypes", compilerOptions_valueOf(o.stripTypes))
	h.Add("summaryDetailLevel", o.summaryDetailLevel) // add(String, int): never omitted
	h.Add("syntheticBlockEndMarker", compilerOptions_valueOf(o.syntheticBlockEndMarker))
	h.Add("syntheticBlockStartMarker", compilerOptions_valueOf(o.syntheticBlockStartMarker))
	h.Add("tcProjectId", compilerOptions_valueOf(o.tcProjectId))
	h.Add("tracer", compilerOptions_valueOf(o.tracer))
	h.AddBool("trustedStrings", o.trustedStrings)
	h.Add("tweakProcessing", compilerOptions_valueOf(o.GetTweakProcessing()))
	h.Add("unusedImportsToRemove", compilerOptions_valueOf(o.unusedImportsToRemove))
	h.AddBool("useTypesForLocalOptimization", o.useTypesForLocalOptimization)
	h.Add("variableRenaming", compilerOptions_valueOf(o.variableRenaming))
	h.Add("warningsGuard", compilerOptions_valueOf(o.GetWarningsGuard()))
	h.AddBool("wrapGoogModulesForWhitespaceOnly", o.wrapGoogModulesForWhitespaceOnly)
	return h.String()
}

// Java: CompilerOptions.expectStrictModeInput()
func (o *CompilerOptions) ExpectStrictModeInput() bool {
	if o.isStrictModeInput != base.Tri_UNKNOWN {
		return o.isStrictModeInput == base.Tri_TRUE
	}
	return o.GetLanguageIn().IsDefaultStrict()
}

// Java: CompilerOptions.setStrictModeInput(boolean)
func (o *CompilerOptions) SetStrictModeInput(isStrictModeInput bool) *CompilerOptions {
	o.isStrictModeInput = base.Tri_ForBoolean(isStrictModeInput)
	return o
}

// Java: CompilerOptions.getPropertyReservedNamingFirstChars()
func (o *CompilerOptions) GetPropertyReservedNamingFirstChars() ASCIICharSet {
	if o.polymerPass {
		return compilerOptions_POLYMER_PROPERTY_RESERVED_FIRST_CHARS
	} else if o.angularPass {
		return compilerOptions_ANGULAR_PROPERTY_RESERVED_FIRST_CHARS
	}
	return ASCIICharSet{}
}

// Java: CompilerOptions.getPropertyReservedNamingNonFirstChars()
func (o *CompilerOptions) GetPropertyReservedNamingNonFirstChars() ASCIICharSet {
	if o.polymerPass {
		return compilerOptions_POLYMER_PROPERTY_RESERVED_NON_FIRST_CHARS
	}
	return ASCIICharSet{}
}

// Java: CompilerOptions.shouldOptimize()
func (o *CompilerOptions) ShouldOptimize() bool {
	return !o.skipNonTranspilationPasses &&
		!o.checksOnly &&
		!o.ShouldGenerateTypedExterns() &&
		!o.instrumentForCoverageOnly
}

// compilerOptions_valueOf is the value ToStringHelper.add(String, Object)
// receives: nil for a Java null, otherwise its String.valueOf.
func compilerOptions_valueOf(value any) any {
	if jutil.IsNil(value) {
		return nil
	}
	return compilerOptions_toJavaString(value)
}

// compilerOptions_optional is a present or absent Guava Optional.
type compilerOptions_optional struct {
	value any // nil when absent
}

// compilerOptions_optionalOf wraps a nil-able pointer as an Optional.
func compilerOptions_optionalOf[T any](p *T) compilerOptions_optional {
	if p == nil {
		return compilerOptions_optional{}
	}
	return compilerOptions_optional{value: p}
}

func compilerOptions_optionalOfTri(t base.Tri) compilerOptions_optional {
	switch t {
	case base.Tri_TRUE:
		return compilerOptions_optional{value: true}
	case base.Tri_FALSE:
		return compilerOptions_optional{value: false}
	default:
		return compilerOptions_optional{}
	}
}

func compilerOptions_formatSlice[T any](s []T) string {
	parts := make([]string, len(s))
	for i, elem := range s {
		parts[i] = compilerOptions_toJavaString(elem)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// compilerOptions_toJavaString is String.valueOf(Object) for the option
// values.
//
// Go-only: language-adapter — slices are printed like Java lists, and a
// value whose Java class does not override toString() prints as
// Object.toString() (class name, '@', identity hash).
func compilerOptions_toJavaString(value any) string {
	switch v := value.(type) {
	case compilerOptions_optional:
		if v.value == nil {
			return "Optional.absent()"
		}
		return "Optional.of(" + compilerOptions_toJavaString(v.value) + ")"
	case *string:
		return *v
	case string, bool, int, int32, int64, float64:
		return jutil.ToJavaString(v)
	case []string:
		return "[" + strings.Join(v, ", ") + "]"
	case []*ConformanceConfig:
		return compilerOptions_formatSlice(v)
	case []SourceMap_LocationMapping:
		return compilerOptions_formatSlice(v)
	case fmt.Stringer:
		return v.String()
	}
	typeName := fmt.Sprintf("%T", value)
	typeName = strings.TrimPrefix(typeName, "*")
	if dot := strings.LastIndexByte(typeName, '.'); dot >= 0 {
		typeName = typeName[dot+1:]
	}
	return jutil.ObjectToString(
		"com.google.javascript.jscomp."+strings.ReplaceAll(typeName, "_", "$"), value)
}
