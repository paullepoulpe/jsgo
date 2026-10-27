// Port of com/google/javascript/jscomp/Compiler.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"io"
	"slices"
	"strings"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/resources"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// Java: Compiler.getChangeTracker()
func (c *Compiler) GetChangeTracker() *ChangeTracker {
	return c.changeTracker
}

// Java: Compiler.getExternsRoot()
func (c *Compiler) GetExternsRoot() rhino.Node {
	return c.externsRoot
}

// Java: Compiler.getJsRoot()
func (c *Compiler) GetJsRoot() rhino.Node {
	return c.jsRoot
}

// Java: Compiler.hasScopeChanged(Node)
func (c *Compiler) HasScopeChanged(n rhino.Node) bool {
	if c.phaseOptimizer == nil {
		return true
	}
	return c.phaseOptimizer.HasScopeChanged(n)
}

// Java: Compiler.reportChangeToChangeScope(Node)
func (c *Compiler) ReportChangeToChangeScope(changeScopeRoot rhino.Node) {
	c.changeTracker.ReportChangeToChangeScope(changeScopeRoot)
}

// Java: Compiler.reportFunctionDeleted(Node)
func (c *Compiler) ReportFunctionDeleted(n rhino.Node) {
	c.changeTracker.ReportFunctionDeleted(n)
}

// Java: Compiler.reportDisambiguatePropertiesSummary(Supplier<String>)
func (c *Compiler) ReportDisambiguatePropertiesSummary(summarySupplier func() string) {
	if c.tracker != nil {
		c.tracker.SetDisambiguatePropertiesSummary(summarySupplier())
	}
}

// Java: Compiler.reportAmbiguatePropertiesSummary(Supplier<String>)
func (c *Compiler) ReportAmbiguatePropertiesSummary(summarySupplier func() string) {
	if c.tracker != nil {
		c.tracker.SetAmbiguatePropertiesSummary(summarySupplier())
	}
}

// Java: Compiler.reportChangeToEnclosingScope(Node)
func (c *Compiler) ReportChangeToEnclosingScope(n rhino.Node) {
	c.changeTracker.ReportChangeToEnclosingScope(n)
}

// Java: Compiler.getCodingConvention()
func (c *Compiler) GetCodingConvention() CodingConvention {
	convention := c.options.GetCodingConvention()
	if convention == nil {
		convention = c.defaultCodingConvention
	}
	return convention
}

// Java: Compiler.getParserConfigLanguageMode(CompilerOptions.LanguageMode)
func (c *Compiler) getParserConfigLanguageMode(languageMode CompilerOptions_LanguageMode) parsing.Config_LanguageMode {
	switch languageMode {
	case CompilerOptions_LanguageMode_ECMASCRIPT3:
		return parsing.Config_LanguageMode_ECMASCRIPT3
	case CompilerOptions_LanguageMode_ECMASCRIPT5, CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT:
		return parsing.Config_LanguageMode_ECMASCRIPT5
	case CompilerOptions_LanguageMode_ECMASCRIPT_2015:
		return parsing.Config_LanguageMode_ECMASCRIPT_2015
	case CompilerOptions_LanguageMode_ECMASCRIPT_2016:
		return parsing.Config_LanguageMode_ECMASCRIPT_2016
	case CompilerOptions_LanguageMode_ECMASCRIPT_2017:
		return parsing.Config_LanguageMode_ECMASCRIPT_2017
	case CompilerOptions_LanguageMode_ECMASCRIPT_2018:
		return parsing.Config_LanguageMode_ECMASCRIPT_2018
	case CompilerOptions_LanguageMode_ECMASCRIPT_2019:
		return parsing.Config_LanguageMode_ECMASCRIPT_2019
	case CompilerOptions_LanguageMode_ECMASCRIPT_2020:
		return parsing.Config_LanguageMode_ECMASCRIPT_2020
	case CompilerOptions_LanguageMode_ECMASCRIPT_2021:
		return parsing.Config_LanguageMode_ECMASCRIPT_2021
	case CompilerOptions_LanguageMode_ECMASCRIPT_NEXT:
		return parsing.Config_LanguageMode_ES_NEXT
	case CompilerOptions_LanguageMode_UNSTABLE:
		return parsing.Config_LanguageMode_UNSTABLE
	case CompilerOptions_LanguageMode_UNSUPPORTED:
		return parsing.Config_LanguageMode_UNSUPPORTED
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected language mode: " + c.options.GetLanguageIn().String()))
	}
}

// Java: Compiler.getParserConfig(ConfigContext)
func (c *Compiler) GetParserConfig(context AbstractCompiler_ConfigContext) *parsing.Config {
	if c.parserConfig == nil || c.externsParserConfig == nil {
		func() {
			c.monitor.Lock()
			defer c.monitor.Unlock()
			if c.parserConfig == nil {
				configLanguageMode := c.getParserConfigLanguageMode(c.options.GetLanguageIn())
				strictMode := parsing.Config_StrictMode_SLOPPY
				if c.options.ExpectStrictModeInput() {
					strictMode = parsing.Config_StrictMode_STRICT
				}
				c.parserConfig = c.CreateConfig(configLanguageMode, strictMode)
				// Externs must always be parsed with at least ES5 language mode.
				if configLanguageMode == parsing.Config_LanguageMode_ECMASCRIPT3 {
					c.externsParserConfig = c.CreateConfig(parsing.Config_LanguageMode_ECMASCRIPT5, strictMode)
				} else {
					c.externsParserConfig = c.parserConfig
				}
			}
		}()
	}
	if context == AbstractCompiler_ConfigContext_EXTERNS {
		return c.externsParserConfig
	}
	return c.parserConfig
}

// Java: Compiler.createConfig(LanguageMode, StrictMode)
func (c *Compiler) CreateConfig(mode parsing.Config_LanguageMode, strictMode parsing.Config_StrictMode) *parsing.Config {
	runMode := parsing.Config_RunMode_STOP_AFTER_ERROR
	if c.options.CanContinueAfterErrors() {
		runMode = parsing.Config_RunMode_KEEP_GOING
	}
	return parsing.ParserRunner_CreateConfigWithModes(
		mode,
		c.options.IsParseJsDocDocumentation(),
		runMode,
		c.options.GetExtraAnnotationNames(),
		c.options.GetParseInlineSourceMaps(),
		strictMode)
}

// ------------------------------------------------------------------------
// Error reporting
// ------------------------------------------------------------------------

// GetDiagnosticGroups: the warning classes that are available from the
// command-line, and are suppressible by the @suppress annotation.
//
// Java: Compiler.getDiagnosticGroups()
func (c *Compiler) GetDiagnosticGroups() *DiagnosticGroups {
	return NewDiagnosticGroups()
}

// Java: Compiler.report(JSError)
func (c *Compiler) Report(error *JSError) {
	level := error.DefaultLevel().Ptr()
	if c.warningsGuard != nil {
		newLevel := c.warningsGuard.Level(error)
		if newLevel != nil {
			level = newLevel
		}
	}

	if level.IsOn() {
		c.InitCompilerOptionsIfTesting()
		if c.GetOptions().GetErrorHandler() != nil {
			c.GetOptions().GetErrorHandler().Report(level, error)
		}
		c.errorManager.Report(level, error)
	}
}

// ReportWithLevel is report(CheckLevel, JSError), Compiler's ErrorHandler
// method (see AsErrorHandler).
//
// Java: Compiler.report(CheckLevel, JSError)
func (c *Compiler) ReportWithLevel(ignoredLevel *CheckLevel, error *JSError) {
	c.Report(error)
}

// Java: Compiler.getErrorLevel(JSError)
func (c *Compiler) GetErrorLevel(error *JSError) *CheckLevel {
	jutil.CheckNotNull(c.options)
	return c.warningsGuard.Level(error)
}

// ThrowInternalError reports an internal error.
//
// Java: Compiler.throwInternalError(String, Throwable)
func (c *Compiler) ThrowInternalError(message string, cause jutil.Throwable) {
	panic(jutil.NewRuntimeExceptionWithMessageAndCause(
		"INTERNAL COMPILER ERROR.\nPlease report this problem.\n\n"+message, cause))
}

// GetErrorCount gets the number of errors.
//
// Java: Compiler.getErrorCount()
func (c *Compiler) GetErrorCount() int {
	return c.errorManager.GetErrorCount()
}

// GetWarningCount gets the number of warnings.
//
// Java: Compiler.getWarningCount()
func (c *Compiler) GetWarningCount() int {
	return c.errorManager.GetWarningCount()
}

// Java: Compiler.hasHaltingErrors()
func (c *Compiler) HasHaltingErrors() bool {
	return !c.GetOptions().CanContinueAfterErrors() && c.errorManager.HasHaltingErrors()
}

// HasErrors consults the ErrorManager to see if we've encountered errors
// that should halt compilation.
//
// If CompilerOptions#canContinueAfterErrors is true, this function always
// returns false without consulting the error manager. The error manager
// will continue to be told about new errors and warnings, but the compiler
// will complete compilation of all inputs.
//
// Java: Compiler.hasErrors()
func (c *Compiler) HasErrors() bool {
	return c.HasHaltingErrors()
}

// Java: Compiler.getSourceFileByName(String)
func (c *Compiler) GetSourceFileByName(sourceName string) *SourceFile {
	// Here we assume that the source name is the input name, this
	// is true of JavaScript parsed from source.
	input := c.inputsById.Get(rhino.NewInputId(sourceName))
	if input != nil {
		return input.GetSourceFile()
	}
	// Alternatively, the sourceName might have been reverse-mapped by
	// an input source-map, so let's look in our sourcemap original sources.
	c.sourceMapMu.RLock()
	sf := c.sourceMapOriginalSources.Get(sourceName)
	c.sourceMapMu.RUnlock()
	return sf
}

// getSourceFileByNullableName is getSourceFileByName for a nullable name.
//
// Go-only: language-adapter — the Java null check of getSourceFileByName,
// for the callers whose name may be null.
func (c *Compiler) getSourceFileByNullableName(sourceName *string) *SourceFile {
	if sourceName != nil {
		return c.GetSourceFileByName(*sourceName)
	}
	return nil
}

// Java: Compiler.getSourceFileContentByName(String)
func (c *Compiler) GetSourceFileContentByName(sourceName string) *string {
	file := c.GetSourceFileByName(sourceName)
	if file.IsStubSourceFileForAlreadyProvidedInput() {
		return nil
	}
	jutil.CheckNotNull(file)
	var result *string
	jutil.Try(func() {
		code := file.GetCode()
		result = &code
	}, func(e jutil.IIOException) {
		result = nil
	})
	return result
}

// Java: Compiler.addInputSourceMap(String, SourceMapInput)
func (c *Compiler) AddInputSourceMap(sourceFileName string, inputSourceMap *SourceMapInput) {
	c.inputSourceMaps.Put(sourceFileName, inputSourceMap)
	if c.options.GetSourceMapIncludeSourcesContent() && c.sourceMap != nil {
		c.addSourceMapSourceFiles(inputSourceMap)
	}
}

// GetBase64SourceMapContents returns the <encoded_source_map> from a `//#
// sourceMappingURL=data:application/json;base64,<encoded_source_map>`
// comment. This sourceMappingURL comment is a JS file's inline source map,
// embedded into the input JS file with a base64-encoded "data url" stored
// in `//# sourceMappingURL=` comment.
//
// Java: Compiler.getBase64SourceMapContents(String)
func (c *Compiler) GetBase64SourceMapContents(sourceFileName string) *string {
	sourceMapInput := c.inputSourceMaps.GetOrDefault(sourceFileName, nil)
	// This sourceMappingURL ends with a ".inline.map", and we want to get its base4-encoded
	// contents (i.e. "data:application/json;base64,<encoded_source_map>").
	var sourceMappingURL *string
	if sourceMapInput != nil {
		originalPath := sourceMapInput.GetOriginalPath()
		sourceMappingURL = &originalPath
	}
	if sourceMappingURL != nil && strings.HasSuffix(*sourceMappingURL, ".inline.map") {
		if sf := sourceMapInput.sourceFile; sf != nil && sf.code.Load() == nil {
			if l, ok := sf.loader.(*sourceFile_CodeLoader_InlineBase64); ok && l.canonicalUTF8 {
				raw := l.rawBase64
				return &raw
			}
		}
		unencoded := sourceMapInput.GetRawSourceMapContents() // This is the sourcemap.
		encoded := jutil.BaseEncoding_Base64().Encode([]byte(*unencoded))
		return &encoded
	}
	return nil
}

// addSourceMapSourceFiles adds file name to content mappings for all
// sources found in a source map. This is used to populate sourcesContent
// array in the output source map even for sources embedded in the input
// source map.
//
// Java: Compiler.addSourceMapSourceFiles(SourceMapInput)
func (c *Compiler) addSourceMapSourceFiles(inputSourceMap *SourceMapInput) {
	// synchronized annotation guards concurrent access to sourceMap during parsing.
	c.monitor.Lock()
	defer c.monitor.Unlock()
	consumer := inputSourceMap.GetSourceMap(c.errorManager)
	if consumer == nil {
		return
	}
	sourcesContent := consumer.GetOriginalSourcesContent()
	if sourcesContent == nil {
		return
	}
	content := sourcesContent
	sources := consumer.GetOriginalSources()
	i := 0
	for i < len(sources) && i < len(content) {
		code := content[i]
		source := SourceMapResolver_GetRelativePath(inputSourceMap.GetOriginalPath(), *sources[i])
		i++
		if source != nil {
			c.sourceMap.AddSourceFile(source.GetName(), code)
		}
	}
	if i < len(sources) || i < len(content) {
		panic(jutil.NewRuntimeException(
			"Source map's \"sources\" and \"sourcesContent\" lengths do not match."))
	}
}

// Java: Compiler.getSourceMapping(String, int, int)
func (c *Compiler) GetSourceMapping(
	sourceName *string, lineNumber int, columnNumber int) *sourcemap.Mapping_OriginalMapping {
	if sourceName == nil {
		return nil
	}
	result, _, hasInputMap, isDefaultEmpty := c.getSourceMappingValue(*sourceName, lineNumber, columnNumber)
	if !hasInputMap {
		return nil
	}
	if isDefaultEmpty {
		return sourcemap.Mapping_OriginalMapping_GetDefaultInstance()
	}
	return &result
}

type compiler_lockedErrorManager struct {
	ErrorManager
	mu *sync.RWMutex
}

func (m compiler_lockedErrorManager) Report(level *CheckLevel, err *JSError) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ErrorManager.Report(level, err)
}

func (c *Compiler) newWorkerSourceMap() *SourceMap {
	sm := c.options.GetSourceMapFormat().GetInstance()
	sm.SetPrefixMappings(c.options.GetSourceMapLocationMappings())
	if c.sourceMap != nil {
		sm.SetSourceFileMapping(c.sourceMap.mapping)
		if sm.generatorV3 != nil && c.sourceMap.generatorV3 != nil {
			sm.generatorV3.ShareConfigFrom(c.sourceMap.generatorV3)
		}
	}
	return sm
}

func (c *Compiler) getSourceMappingValue(
	sourceName string, lineNumber int, columnNumber int) (sourcemap.Mapping_OriginalMapping, *string, bool, bool) {
	return c.getSourceMappingValueCached(sourceName, lineNumber, columnNumber, &c.resolvedSourceMap, nil)
}

func (c *Compiler) getSourceMappingValueCached(
	sourceName string,
	lineNumber int,
	columnNumber int,
	resolvedCache *compiler_ResolvedSourceMap,
	lineCache *sourcemap.SourceMapLineCache,
) (sourcemap.Mapping_OriginalMapping, *string, bool, bool) {
	if c.inputSourceMaps.IsEmpty() {
		return sourcemap.Mapping_OriginalMapping{}, nil, false, false
	}
	sourceMap := c.inputSourceMaps.Get(sourceName)
	if sourceMap == nil {
		return sourcemap.Mapping_OriginalMapping{}, nil, false, false
	}

	// JSCompiler uses 1-indexing for lineNumber and 0-indexing for columnNumber.
	// Sourcemaps use 1-indexing for both.
	var consumer *sourcemap.SourceMapConsumerV3
	if sourceMap.ready.Load() {
		consumer = sourceMap.parsedSourceMap
	} else {
		consumer = sourceMap.GetSourceMap(compiler_lockedErrorManager{ErrorManager: c.errorManager, mu: &c.sourceMapMu})
	}
	if consumer == nil {
		return sourcemap.Mapping_OriginalMapping{}, nil, false, false
	}
	if !c.options.GetSourceMapIncludeSourcesContent() {
		// Go-only: performance — only addSourceMapSourceFiles reads the
		// sources' contents, and only with that option.
		sourceMap.dropSourcesContentOnce(consumer)
	}
	result, identifierPtr, ok := consumer.GetMappingForLineValueCached(lineNumber, columnNumber+1, lineCache)
	if !ok {
		return sourcemap.Mapping_OriginalMapping{}, nil, true, true
	}

	// First check to see if the original file was loaded from an input source map.
	sourceMapOriginalPath := sourceMap.GetOriginalPath()
	resultOriginalPath := result.GetOriginalFile()
	var relativePath string

	// Resolving the paths to a source file is expensive, so check the cache first.
	if resolvedCache.hasCached &&
		sourceMapOriginalPath == resolvedCache.originalPath &&
		resultOriginalPath == resolvedCache.sourceMapPath {
		relativePath = resolvedCache.relativePath
	} else {
		relativePath = compiler_resolveSibling(sourceMapOriginalPath, resultOriginalPath)
		source := c.GetSourceFileByName(relativePath)
		if source == nil && resultOriginalPath != "" {
			source = SourceMapResolver_GetRelativePath(
				sourceMap.GetOriginalPath(), result.GetOriginalFile())
			if source != nil {
				c.sourceMapMu.Lock()
				c.sourceMapOriginalSources.PutIfAbsent(relativePath, source)
				c.sourceMapMu.Unlock()
			}
		}

		// Cache this resolved source for the next caller.
		resolvedCache.hasCached = true
		resolvedCache.originalPath = sourceMapOriginalPath
		resolvedCache.sourceMapPath = resultOriginalPath
		resolvedCache.relativePath = relativePath
	}

	result.SetOriginalFileAndColumnInPlace(relativePath, result.GetColumnPosition()-1)
	return result, identifierPtr, true, false
}

// Java: Compiler.getSourceLine(String, int)
func (c *Compiler) GetSourceLine(sourceName *string, lineNumber int) *string {
	if lineNumber < 1 {
		return nil
	}
	input := c.getSourceFileByNullableName(sourceName)
	if input != nil && !input.IsStubSourceFileForAlreadyProvidedInput() {
		return input.GetLine(lineNumber)
	}
	return nil
}

// Java: Compiler.getSourceLines(String, int, int)
func (c *Compiler) GetSourceLines(sourceName *string, lineNumber int, length int) Region {
	if lineNumber < 1 {
		return nil
	}
	input := c.getSourceFileByNullableName(sourceName)
	if input != nil && !input.IsStubSourceFileForAlreadyProvidedInput() {
		return input.GetLines(lineNumber, length)
	}
	return nil
}

// Java: Compiler.getSourceRegion(String, int)
func (c *Compiler) GetSourceRegion(sourceName *string, lineNumber int) Region {
	if lineNumber < 1 {
		return nil
	}
	input := c.getSourceFileByNullableName(sourceName)
	if input != nil && !input.IsStubSourceFileForAlreadyProvidedInput() {
		return input.GetRegion(lineNumber)
	}
	return nil
}

// ------------------------------------------------------------------------
// Package-private and Protected helpers
// ------------------------------------------------------------------------

// Java: Compiler.getNodeForCodeInsertion(JSChunk)
func (c *Compiler) GetNodeForCodeInsertion(chunk *JSChunk) rhino.Node {
	if c.inputsById.ContainsKey(compiler_SYNTHETIC_CODE_INPUT_ID) {
		return c.inputsById.Get(compiler_SYNTHETIC_CODE_INPUT_ID).GetAstRoot(c)
	}
	if chunk == nil {
		allInputs := c.chunkGraph.GetAllInputs()
		if c.chunkGraph == nil || len(allInputs) == 0 {
			panic(jutil.NewIllegalStateExceptionMsg("No inputs"))
		}
		firstInput := allInputs[0]
		return compiler_checkNotModule(firstInput.GetAstRoot(c), "Cannot insert code into a module")
	}

	chunkInputs := chunk.GetInputs()
	if len(chunkInputs) != 0 {
		return compiler_checkNotModule(
			chunkInputs[0].GetAstRoot(c), "Cannot insert code into a module")
	}
	panic(jutil.NewIllegalStateExceptionMsg("Root chunk has no inputs"))
}

// Java: Compiler.getSourceMap()
func (c *Compiler) GetSourceMap() *SourceMap {
	return c.sourceMap
}

// Java: Compiler.setVariableMap(VariableMap)
func (c *Compiler) SetVariableMap(variableMap *VariableMap) {
	c.variableMap = variableMap
}

// Java: Compiler.getVariableMap()
func (c *Compiler) GetVariableMap() *VariableMap {
	return c.variableMap
}

// Java: Compiler.setPropertyMap(VariableMap)
func (c *Compiler) SetPropertyMap(propertyMap *VariableMap) {
	c.propertyMap = propertyMap
}

// Java: Compiler.getPropertyMap()
func (c *Compiler) GetPropertyMap() *VariableMap {
	return c.propertyMap
}

// Java: Compiler.setStringMap(VariableMap)
func (c *Compiler) SetStringMap(stringMap *VariableMap) {
	c.stringMap = stringMap
}

// Java: Compiler.setCssNames(Set<String>)
func (c *Compiler) SetCssNames(cssNames *jutil.LinkedHashSet[string]) {
	c.cssNames = cssNames
}

// Java: Compiler.setIdGeneratorMap(String)
func (c *Compiler) SetIdGeneratorMap(serializedIdMappings string) {
	c.idGeneratorMap = &serializedIdMappings
}

// Java: Compiler.getTranspiledFiles()
func (c *Compiler) GetTranspiledFiles() bool {
	return c.transpiledFiles
}

// Java: Compiler.setTranspiledFiles(boolean)
func (c *Compiler) SetTranspiledFiles(transpiledFiles bool) {
	c.transpiledFiles = transpiledFiles
}

// Java: Compiler.getCrossChunkIdGenerator()
func (c *Compiler) GetCrossChunkIdGenerator() *IdGenerator {
	return c.crossChunkIdGenerator
}

// Java: Compiler.setAnonymousFunctionNameMap(VariableMap)
func (c *Compiler) SetAnonymousFunctionNameMap(functionMap *VariableMap) {
	// NOTE: remove this method
}

// Java: Compiler.getStringMap()
func (c *Compiler) GetStringMap() *VariableMap {
	return c.stringMap
}

// Java: Compiler.setInstrumentationMapping(VariableMap)
func (c *Compiler) SetInstrumentationMapping(instrumentationMapping *VariableMap) {
	c.instrumentationMapping = instrumentationMapping
}

// Java: Compiler.getInstrumentationMapping()
func (c *Compiler) GetInstrumentationMapping() *VariableMap {
	return c.instrumentationMapping
}

// Java: Compiler.addExportedNames(Set<String>)
func (c *Compiler) AddExportedNames(exportedNames *jutil.LinkedHashSet[string]) {
	c.exportedNames.AddAll(exportedNames.ToSlice()...)
}

// Java: Compiler.getExportedNames()
func (c *Compiler) GetExportedNames() *jutil.LinkedHashSet[string] {
	return c.exportedNames
}

// Java: Compiler.setDefineNames(Collection<String>)
func (c *Compiler) SetDefineNames(defineNames []string) {
	c.defineNames = jutil.NewLinkedHashSet[string](defineNames...)
}

// Java: Compiler.getDefineNames()
func (c *Compiler) GetDefineNames() *jutil.LinkedHashSet[string] {
	return c.defineNames
}

// Java: Compiler.getOptions()
func (c *Compiler) GetOptions() *CompilerOptions {
	return c.options
}

// Compiler_SetLoggingLevel sets the logging level for the
// com.google.javascript.jscomp package.
//
// Java: Compiler.setLoggingLevel(Level)
func Compiler_SetLoggingLevel(level *jutil.Level) {
	Compiler_logger.SetLevel(level)
}

// GetAstDotGraph gets the DOT graph of the AST generated at the end of
// compilation.
//
// Java: Compiler.getAstDotGraph()
func (c *Compiler) GetAstDotGraph() string {
	if !c.jsRoot.IsNil() {
		cfg := ControlFlowAnalysis_Builder_().
			SetCompiler(c).
			SetCfgRoot(c.jsRoot).
			SetTraverseFunctions(true).
			ComputeCfg()
		return DotFormatter_ToDotWithCfg(c.jsRoot, cfg)
	} else {
		return ""
	}
}

// Java: Compiler.getErrorManager()
func (c *Compiler) GetErrorManager() ErrorManager {
	if c.options == nil {
		c.InitOptions(NewCompilerOptions())
	}
	return c.errorManager
}

// Java: Compiler.getInputsInOrder()
func (c *Compiler) GetInputsInOrder() []*CompilerInput {
	return c.chunkGraph.GetAllInputs()
}

// Java: Compiler.getNumberOfInputs()
func (c *Compiler) GetNumberOfInputs() int {
	// In some testing cases inputs will be null, but obviously there must be at least one input.
	// The intended use of this method is to allow passes to estimate how much memory they will
	// need for data structures, so it's not necessary that the returned value be exactly right
	// in the corner cases where inputs ends up being null.
	if c.chunkGraph != nil {
		return c.chunkGraph.GetInputCount()
	}
	return 1
}

// GetInputsById returns an unmodifiable view of the compiler inputs
// indexed by id.
//
// Java: Compiler.getInputsById()
func (c *Compiler) GetInputsById() *jutil.JavaHashMap[*rhino.InputId, *CompilerInput] {
	return c.inputsById
}

// GetExternsInOrder gets the externs in the order in which they are being
// processed.
//
// Java: Compiler.getExternsInOrder()
func (c *Compiler) GetExternsInOrder() []*CompilerInput {
	return c.externs
}

// Java: Compiler.getInputsForTesting()
func (c *Compiler) GetInputsForTesting() []*CompilerInput {
	if c.chunkGraph != nil {
		return append([]*CompilerInput(nil), c.chunkGraph.GetAllInputs()...)
	}
	return nil
}

// Java: Compiler.getExternsForTesting()
func (c *Compiler) GetExternsForTesting() []*CompilerInput {
	return c.externs
}

// Java: Compiler.hasRegExpGlobalReferences()
func (c *Compiler) HasRegExpGlobalReferences() bool {
	return c.hasRegExpGlobalReferences
}

// Java: Compiler.setHasRegExpGlobalReferences(boolean)
func (c *Compiler) SetHasRegExpGlobalReferences(references bool) {
	c.hasRegExpGlobalReferences = references
}

// SetRunJ2clPasses sets whether J2CL passes should run
//
// Java: Compiler.setRunJ2clPasses(boolean)
func (c *Compiler) SetRunJ2clPasses(runJ2clPasses bool) {
	c.runJ2clPasses = runJ2clPasses
}

// RunJ2clPasses reports whether J2CL passes should run
//
// Java: Compiler.runJ2clPasses()
func (c *Compiler) RunJ2clPasses() bool {
	return c.runJ2clPasses
}

// Java: Compiler.getSynthesizedExternsInput()
func (c *Compiler) GetSynthesizedExternsInput() *CompilerInput {
	if c.syntheticExternsInput != nil {
		return c.syntheticExternsInput
	}

	input := NewCompilerInputWithIsExtern(c.SYNTHETIC_EXTERNS_FILE /* isExtern= */, true)
	root := rhino.CheckNotNull(input.GetAstRoot(c))
	c.putCompilerInput(input)
	c.syntheticExternsInput = input
	c.externsRoot.AddChildToFront(root)
	c.externs = slices.Insert(c.externs, 0, input)
	c.scriptNodeByFilename.Put(input.GetSourceFile().GetName(), root)

	return input
}

// Java: Compiler.getSynthesizedExternsRoot()
func (c *Compiler) getSynthesizedExternsRoot() rhino.Node {
	externsInput := c.GetSynthesizedExternsInput()
	return externsInput.GetAstRoot(c)
}

// Java: Compiler.getSynthesizedTypeSummaryInput()
func (c *Compiler) GetSynthesizedTypeSummaryInput() *CompilerInput {
	if c.syntheticTypeSummaryInput != nil {
		return c.syntheticTypeSummaryInput
	}

	input := NewCompilerInputWithIsExtern(c.syntheticTypeSummaryFile /* isExtern= */, true)
	root := rhino.CheckNotNull(input.GetAstRoot(c))
	c.putCompilerInput(input)
	c.syntheticTypeSummaryInput = input
	c.externsRoot.AddChildToFront(root)
	c.externs = slices.Insert(c.externs, 0, input)
	c.scriptNodeByFilename.Put(input.GetSourceFile().GetName(), root)
	builder := rhino.JSDocInfo_Builder_()
	builder.RecordTypeSummary()
	root.SetJSDocInfo(builder.Build())

	return input
}

// Java: Compiler.getSyntheticCodeInputId()
func (c *Compiler) GetSyntheticCodeInputId() *rhino.InputId {
	return compiler_SYNTHETIC_CODE_INPUT_ID
}

// Java: Compiler.initializeSyntheticCodeInput()
func (c *Compiler) InitializeSyntheticCodeInput() {
	jutil.CheckStateMsg(
		!c.inputsById.ContainsKey(compiler_SYNTHETIC_CODE_INPUT_ID),
		"Already initialized synthetic input")
	ast := NewCompilerInput(SourceFile_FromCode(compiler_SYNTHETIC_CODE_INPUT_ID.GetIdName(), ""))
	if c.inputsById.ContainsKey(ast.GetInputId()) {
		panic(jutil.NewIllegalStateExceptionMsg("Conflicting synthetic id name"))
	}
	c.jsRoot.AddChildToFront(rhino.CheckNotNull(ast.GetAstRoot(c)))

	var firstChunk *JSChunk
	if chunks := c.GetChunks(); len(chunks) > 0 {
		firstChunk = chunks[0]
	}
	if firstChunk.GetName() == JSChunk_STRONG_CHUNK_NAME {
		firstChunk.Add(ast)
	}
	ast.SetChunk(firstChunk)
	c.putCompilerInput(ast)

	c.commentsPerFile.Put(compiler_SYNTHETIC_CODE_INPUT_ID.GetIdName(), []*parser.Comment{})
	c.ReportChangeToChangeScope(ast.GetAstRoot(c))
}

// Java: Compiler.removeSyntheticCodeInput()
func (c *Compiler) RemoveSyntheticCodeInput() {
	c.removeSyntheticCodeInputMerging( /* mergeContentIntoFirstInput= */ false)
}

// Java: Compiler.mergeSyntheticCodeInput()
func (c *Compiler) MergeSyntheticCodeInput() {
	c.removeSyntheticCodeInputMerging( /* mergeContentIntoFirstInput= */ true)
}

// removeSyntheticCodeInputMerging deletes the synthesized code input and
// optionally moves all its subtree contents into the first non-synthetic
// input
//
// Java: Compiler.removeSyntheticCodeInput(boolean)
func (c *Compiler) removeSyntheticCodeInputMerging(mergeContentIntoFirstInput bool) {
	jutil.CheckStateMsg(
		c.inputsById.ContainsKey(compiler_SYNTHETIC_CODE_INPUT_ID),
		"Never initialized the synthetic input")
	input := c.inputsById.Get(compiler_SYNTHETIC_CODE_INPUT_ID)
	astRoot := input.GetAstRoot(c)
	jutil.CheckState(astRoot.IsFirstChildOf(c.jsRoot))
	jutil.CheckState(compiler_SYNTHETIC_CODE_INPUT_ID.Equals(input.GetInputId()))

	if mergeContentIntoFirstInput && astRoot.HasChildren() {
		// If we've inserted anything into the synthetic AST, move it into the top of the next script.
		next := astRoot.GetNext()
		rhino.CheckNotNullMsg(next, "Must provide at least one source")
		compiler_checkNotModule(
			next, "Cannot remove synthetic code input until modules are rewritten: %s", next)
		next.AddChildrenToFront(astRoot.RemoveChildren())
		c.ReportChangeToChangeScope(next)
	}

	astRoot.Detach()

	// bookkeeping to mark scopes and nodes as deleted
	c.ReportChangeToChangeScope(astRoot)
	astRoot.SetDeleted(true)
	NodeUtil_MarkFunctionsDeleted(astRoot, c)

	input.GetChunk().Remove(input)
	c.inputsById.Remove(input.GetInputId())
}

// Java: Compiler.getLastPassName()
func (c *Compiler) GetLastPassName() *string {
	return c.lastPassName
}

// Java: Compiler.setExternProperties(ImmutableSet<String>)
func (c *Compiler) SetExternProperties(externProperties *jutil.LinkedHashSet[string]) {
	c.externProperties = externProperties
}

// Java: Compiler.getExternProperties()
func (c *Compiler) GetExternProperties() *jutil.LinkedHashSet[string] {
	return c.externProperties
}

// Java: Compiler.getAccessorSummary()
func (c *Compiler) GetAccessorSummary() *AccessorSummary {
	return c.accessorSummary
}

// Java: Compiler.setAccessorSummary(AccessorSummary)
func (c *Compiler) SetAccessorSummary(summary *AccessorSummary) {
	c.accessorSummary = jutil.CheckNotNull(summary)
}

// loadResourceContents reads runtime library .js files into memory in AST
// format
//
// Java: Compiler.loadResourceContents(String, String)
func (c *Compiler) loadResourceContents(resourceName string, path string) rhino.Node {
	if c.GetLifeCycleStage().HasColorAndSimplifiedJSDoc() {
		jutil.CheckNotNullMsg(
			c.runtimeLibraryTypedAsts,
			"Must call initRuntimeLibraryTypedAsts before calling ensureLibraryInjected during"+
				" optimizations")

		typedAstSupplier := jutil.CheckNotNullMsg1(
			c.runtimeLibraryTypedAsts.Get(path),
			strings.Join([]string{
				"Missing precompiled .typedast for '%s'. If this file is newly added, ",
				"you may need to regenerate runtime_libs.typedast.textproto",
				""}, ""),
			path)
		return typedAstSupplier()
	}
	jutil.CheckStateMsg(
		!c.HasTypeCheckingRun(),
		"runtime library injected after type checking but before optimization colors")
	originalCode := resources.ResourceLoader_LoadTextResource("com.google.javascript.jscomp.Compiler", "js/"+resourceName+".js")

	source := SourceFile_FromCode(path, originalCode)
	c.addFilesToSourceMap([]*SourceFile{source})
	return c.parseCodeHelper(source)
}

// Java: Compiler.getNodeForRuntimeCodeInsertion()
func (c *Compiler) getNodeForRuntimeCodeInsertion() rhino.Node {
	// For stage 2 optimizing builds, insert runtime libraries as actual code in the AST.
	// For library builds that outputs a TypedAST, the runtime libraries are injected to a synthetic
	// @typeSummary (.i.js) node so that the types are available.
	if c.options.GetTypedAstOutputFile() == nil {
		return c.GetNodeForCodeInsertion(nil)
	}
	return c.GetSynthesizedTypeSummaryInput().GetAstRoot(c)
}

// Java: Compiler.addComments(String, List<Comment>)
func (c *Compiler) AddComments(filename string, comments []*parser.Comment) {
	if !c.GetOptions().PreservesDetailedSourceInfo() {
		panic(jutil.NewUnsupportedOperationExceptionMsg("addComments may only be called in IDE mode."))
	}
	c.commentsPerFile.Put(filename, comments)
}

// Java: Compiler.getComments(String)
func (c *Compiler) GetComments(filename string) []*parser.Comment {
	if !c.GetOptions().PreservesDetailedSourceInfo() {
		panic(jutil.NewUnsupportedOperationExceptionMsg("getComments may only be called in IDE mode."))
	}
	return c.commentsPerFile.Get(filename)
}

// Java: Compiler.getModuleLoader()
func (c *Compiler) GetModuleLoader() *ModuleLoader {
	return c.moduleLoader
}

// Java: Compiler.addFilesToSourceMap(Iterable<SourceFile>)
func (c *Compiler) addFilesToSourceMap(files []*SourceFile) {
	// synchronized annotation guards concurrent access to sourceMap during parsing.
	c.monitor.Lock()
	defer c.monitor.Unlock()
	if c.GetOptions().GetSourceMapIncludeSourcesContent() && c.GetSourceMap() != nil {
		for _, file := range files {
			jutil.Try(func() {
				c.GetSourceMap().AddSourceFile(file.GetName(), jutil.StringPtr(file.GetCode()))
			}, func(e jutil.IIOException) {
				panic(jutil.NewRuntimeExceptionWithMessageAndCause("Cannot read code of a source map's source file.", e))
			})
		}
	}
}

// Java: Compiler.initWebpackMap(ImmutableMap<String, String>)
func (c *Compiler) InitWebpackMap(inputPathByWebpackId *jutil.LinkedHashMap[string, string]) {
	c.inputPathByWebpackId = inputPathByWebpackId
}

// SetDebugMessage sets an additional human-readable message appended to
// debugging contexts, like crash reports & source-after-each-pass logging.
//
// This is used by NestedCompilerRunner to distinguish between the main
// compiler & nested compiler.
//
// Java: Compiler.setDebugMessage(String)
func (c *Compiler) SetDebugMessage(debugMessage *string) {
	c.compilerExecutor.SetDebugMessage(debugMessage)
	c.debugMessage = debugMessage
}

// Java: Compiler.createCompilerExecutor()
func (c *Compiler) CreateCompilerExecutor() *CompilerExecutor {
	return NewCompilerExecutor()
}

// Java: Compiler.getCompilerExecutor()
func (c *Compiler) GetCompilerExecutor() *CompilerExecutor {
	return c.compilerExecutor
}

// Serializable state of the compiler specific to multistage binary builds
//
// Only contains state that does not make sense in 'multilevel' binary builds
// (where library-level TypedASTs are the input). Such state belongs in the
// jscomp.TypedAst proto.

// Java: Compiler.toProto(Feature)
func compiler_toProto(feature parser.FeatureSet_Feature) FeatureProto {
	var result FeatureProto
	jutil.Try(func() {
		result = FeatureProto_ValueOfName("FEATURE_" + feature.Name())
	}, func(e jutil.IIllegalArgumentException) {
		result = FeatureProto_FEATURE_UNKNOWN
	})
	return result
}

// Java: Compiler.getModuleMetadataMap()
func (c *Compiler) GetModuleMetadataMap() *ModuleMetadataMap {
	return c.moduleMetadataMap
}

// Java: Compiler.setModuleMetadataMap(ModuleMetadataMap)
func (c *Compiler) SetModuleMetadataMap(moduleMetadataMap *ModuleMetadataMap) {
	c.moduleMetadataMap = moduleMetadataMap
}

// Java: Compiler.getModuleMap()
func (c *Compiler) GetModuleMap() *ModuleMap {
	return c.moduleMap
}

// Java: Compiler.setModuleMap(ModuleMap)
func (c *Compiler) SetModuleMap(moduleMap *ModuleMap) {
	c.moduleMap = moduleMap
}

// compiler_resolveSibling is a simplistic implementation of the
// java.nio.file.Path resolveSibling method that always uses "/" regardless
// of platform. This is required for consistency of path handling when
// using input source maps, which are expected to always use "/".
//
// If `toPath` is already absolute, just return it.
//
// Otherwise, strip the filename off of `fromPath` and logically replaces
// it with `toPath`. Then normalize the path by resolving and removing any
// "." or ".." elements.
//
// fromPath - must be a file (not directory); toPath - must be a file (not
// directory). Returns `toPath` adjusted so that it has the same parent
// directory as `fromPath`.
//
// Java: Compiler.resolveSibling(String, String)
func compiler_resolveSibling(fromPath string, toPath string) string {
	// If the destination is an absolute path, nothing to do.
	if strings.HasPrefix(toPath, "/") {
		return toPath
	}

	fromPathParts := jutil.StringSplitLiteral(fromPath, "/")
	toPathParts := jutil.StringSplitLiteral(toPath, "/")
	if len(fromPathParts) != 0 {
		fromPathParts = fromPathParts[:len(fromPathParts)-1]
	}

	for len(fromPathParts) != 0 && len(toPathParts) != 0 {
		if toPathParts[0] == "." {
			toPathParts = toPathParts[1:]
		} else if toPathParts[0] == ".." {
			toPathParts = toPathParts[1:]
			fromPathParts = fromPathParts[:len(fromPathParts)-1]
		} else {
			break
		}
	}

	fromPathParts = append(fromPathParts, toPathParts...)
	return strings.Join(fromPathParts, "/")
}

// Java: Compiler.resetAndIntitializeSourceMap()
func (c *Compiler) ResetAndIntitializeSourceMap() {
	if c.sourceMap == nil {
		return
	}
	c.sourceMap.Reset()
	if c.options.GetSourceMapIncludeSourcesContent() {
		if c.options.GetApplyInputSourceMaps() {
			// Add any input source map content files to the source map as potential sources
			for inputSourceMap := range c.inputSourceMaps.Values() {
				c.addSourceMapSourceFiles(inputSourceMap)
			}
		}

		// Add all the compilation sources to the source map as potential sources
		allChunks := c.GetChunks()
		if allChunks != nil {
			sourceFiles := []*SourceFile{}
			for _, chunk := range allChunks {
				for _, input := range chunk.GetInputs() {
					sourceFile := input.GetSourceFile()
					if sourceFile.IsStubSourceFileForAlreadyProvidedInput() {
						// do not add stub source files to the source map
						continue
					}
					sourceFiles = append(sourceFiles, input.GetSourceFile())
				}
			}
			c.addFilesToSourceMap(sourceFiles)
		}
	}
}

// Java: Compiler.checkNotModule(Node, String, Object...)
func compiler_checkNotModule(script rhino.Node, msg string, args ...any) rhino.Node {
	jutil.CheckArgumentMsg1(script.IsScript(), "%s", script)
	if !script.HasOneChild() {
		return script
	}
	jutil.CheckStateMsg(!script.GetFirstChild().IsModuleBody(), msg, args...)
	return script
}

// Java: Compiler.maybeLogPruningAnalysis(DependencyManagementResult)
func (c *Compiler) maybeLogPruningAnalysis(result *JSChunkGraph_DependencyManagementResult) {
	if c.tracker == nil || !c.options.GetTracerMode().DoPruningAnalysis() {
		return
	}
	analysis := PruningAnalysis_Create(result.Sorter(), result.EntryPointInputs().ToSlice())
	pruningResult := analysis.Analyze()
	c.tracker.SetPruningAnalysisSummary(
		"\nTransitive dependency count per entry point:\n" +
			jutil.ToJavaString(pruningResult.EntryPointDependencyCount()) +
			"\n\nTop bottleneck dependencies:\n" +
			jutil.ToJavaString(pruningResult.BottleneckBlame()))
}

// Compiler_NonClosingOutputStream is a FilterOutputStream whose close() only
// flushes.
//
// Java: Compiler.NonClosingOutputStream
type Compiler_NonClosingOutputStream struct {
	out io.Writer
}

// Java: Compiler.NonClosingOutputStream(OutputStream)
func NewCompiler_NonClosingOutputStream(out io.Writer) *Compiler_NonClosingOutputStream {
	return &Compiler_NonClosingOutputStream{out: out}
}

// Write is FilterOutputStream.write(byte[], int, int).
//
// Go-only: language-adapter — io.Writer for OutputStream.
func (s *Compiler_NonClosingOutputStream) Write(b []byte) (int, error) {
	return s.out.Write(b)
}

// Flush is FilterOutputStream.flush(): it flushes out.
//
// Go-only: language-adapter — an io.Writer flushes if it has a Flush method.
func (s *Compiler_NonClosingOutputStream) Flush() error {
	if f, ok := s.out.(interface{ Flush() error }); ok {
		return f.Flush()
	}
	return nil
}

// Java: Compiler.NonClosingOutputStream.close()
func (s *Compiler_NonClosingOutputStream) Close() error {
	return s.Flush()
}

// Compiler_NonClosingInputStream is a FilterInputStream whose close() does not
// close the underlying stream.
//
// Java: Compiler.NonClosingInputStream
type Compiler_NonClosingInputStream struct {
	in io.Reader
}

// Java: Compiler.NonClosingInputStream(InputStream)
func NewCompiler_NonClosingInputStream(in io.Reader) *Compiler_NonClosingInputStream {
	return &Compiler_NonClosingInputStream{in: in}
}

// Read is FilterInputStream.read(byte[], int, int).
//
// Go-only: language-adapter — io.Reader for InputStream.
func (s *Compiler_NonClosingInputStream) Read(b []byte) (int, error) {
	return s.in.Read(b)
}

// Java: Compiler.NonClosingInputStream.close()
func (s *Compiler_NonClosingInputStream) Close() error {
	// Do not close the underlying stream
	return nil
}
