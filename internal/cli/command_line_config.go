// Port of com/google/javascript/jscomp/AbstractCommandLineRunner.java @ closure-compiler v20261021

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

package cli

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// AbstractCommandLineRunner_CommandLineConfig holds the configurations for
// the command line configs. Designed for easy building, so that we can
// decouple the flags-parsing library from the actual configuration options.
//
// TODO(tjgq): Investigate whether this class is really needed to mediate
// between the CompilerOptions and runner implementations. An alternative
// would be for the runners to fill in the CompilerOptions directly, but that
// conflicts with the latter's mutability and the desire to reuse the same
// options across multiple compilations.
//
// Java: AbstractCommandLineRunner.CommandLineConfig
type AbstractCommandLineRunner_CommandLineConfig struct {
	printVersion bool

	printTree     bool
	printTreeJson bool

	printAst bool

	jscompDevMode jscomp.CompilerOptions_DevMode

	loggingLevel string

	externs []string

	mixedJsSources []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]

	defaultToStdin bool

	jsOutputFile string

	// When non-null specifies a file containing saved compiler state to
	// restore and continue compiling.
	continueSavedCompilationFileName *string

	// When > 0 indicates the stage at which compilation stopped for the
	// compilation state that is being restored from a save file.
	restoredCompilationStage int

	// When > 0 indicates the stage at which compilation should stop and the
	// compilation state be stored into a save file.
	saveAfterCompilationStage int

	stage2SegmentOfCompilationToRun jscomp.CompilerOptions_SegmentOfCompilationToRun

	typedAstListInputFilename *string

	saveCompilationStateToFilename *string

	chunk []string

	sourceMapInputFiles []PipeFlagEntry

	parseInlineSourceMaps bool

	expectedDiagnostics []string

	variableMapInputFile string

	propertyMapInputFile string

	variableMapOutputFile string

	createNameMapFiles bool

	propertyMapOutputFile string

	stringMapOutputPath string

	instrumentationMappingFile string

	codingConvention jscomp.CodingConvention

	summaryDetailLevel int

	outputWrapper string

	chunkWrapper []string

	chunkOutputPathPrefix string

	chunkOutputFiles []string

	chunkConformanceFiles []string

	createSourceMap string

	sourceMapDetailLevel jscomp.SourceMap_DetailLevel

	sourceMapFormat jscomp.SourceMap_Format

	sourceMapLocationMappings []jscomp.SourceMap_LocationMapping

	applyInputSourceMaps bool

	warningGuards []*AbstractCommandLineRunner_FlagEntry[jscomp.CheckLevel]

	define []string

	browserFeaturesetYear int

	tweakProcessing jscomp.CompilerOptions_TweakProcessing

	charset string

	dependencyOptions *jscomp.DependencyOptions // nil is null

	outputManifests []string

	outputChunkDependencies *string // nil is null

	outputBundles []string

	skipNormalOutputs bool

	manifestMaps []string

	processCommonJSModules bool

	moduleRoots []string

	warningsAllowFile string

	hideWarningsFor []string

	angularPass bool

	jsonStreamMode jscomp.CompilerOptions_JsonStreamMode

	errorFormat AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption

	jsonWarningsFile string

	includeTrailingNewline bool
}

// NewAbstractCommandLineRunner_CommandLineConfig runs the field
// initializers of CommandLineConfig.
//
// Java: AbstractCommandLineRunner.CommandLineConfig()
func NewAbstractCommandLineRunner_CommandLineConfig() *AbstractCommandLineRunner_CommandLineConfig {
	return &AbstractCommandLineRunner_CommandLineConfig{
		printTree:                       false,
		printTreeJson:                   false,
		printAst:                        false,
		jscompDevMode:                   jscomp.CompilerOptions_DevMode_OFF,
		loggingLevel:                    jutil.Level_WARNING.GetName(),
		externs:                         []string{},
		mixedJsSources:                  []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]{},
		defaultToStdin:                  false,
		jsOutputFile:                    "",
		restoredCompilationStage:        -1,
		saveAfterCompilationStage:       -1,
		stage2SegmentOfCompilationToRun: jscomp.CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS,
		chunk:                           []string{},
		sourceMapInputFiles:             []PipeFlagEntry{},
		parseInlineSourceMaps:           false,
		expectedDiagnostics:             []string{},
		variableMapInputFile:            "",
		propertyMapInputFile:            "",
		variableMapOutputFile:           "",
		createNameMapFiles:              false,
		propertyMapOutputFile:           "",
		stringMapOutputPath:             "",
		instrumentationMappingFile:      "",
		codingConvention:                jscomp.CodingConventions_GetDefault(),
		summaryDetailLevel:              1,
		outputWrapper:                   "",
		chunkWrapper:                    []string{},
		chunkOutputPathPrefix:           "",
		chunkOutputFiles:                []string{},
		chunkConformanceFiles:           []string{},
		createSourceMap:                 "",
		sourceMapDetailLevel:            jscomp.SourceMap_DetailLevel_ALL,
		sourceMapFormat:                 jscomp.SourceMap_Format_DEFAULT,
		sourceMapLocationMappings:       []jscomp.SourceMap_LocationMapping{},
		applyInputSourceMaps:            false,
		warningGuards:                   []*AbstractCommandLineRunner_FlagEntry[jscomp.CheckLevel]{},
		define:                          []string{},
		browserFeaturesetYear:           0,
		tweakProcessing:                 jscomp.CompilerOptions_TweakProcessing_OFF,
		charset:                         "",
		outputManifests:                 []string{},
		outputBundles:                   []string{},
		skipNormalOutputs:               false,
		manifestMaps:                    []string{},
		processCommonJSModules:          false,
		moduleRoots:                     []string{jscomp.ModuleLoader_DEFAULT_FILENAME_PREFIX},
		warningsAllowFile:               "",
		hideWarningsFor:                 []string{},
		angularPass:                     false,
		jsonStreamMode:                  jscomp.CompilerOptions_JsonStreamMode_NONE,
		errorFormat:                     AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_STANDARD,
		jsonWarningsFile:                "",
		includeTrailingNewline:          true,
	}
}

// Java: AbstractCommandLineRunner.CommandLineConfig.setPrintVersion(boolean)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetPrintVersion(x bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.printVersion = x
	return c
}

// SetPrintTree prints out the parse tree and exits.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setPrintTree(boolean)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetPrintTree(printTree bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.printTree = printTree
	return c
}

// SetPrintTreeJson prints out the parse tree and exits.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setPrintTreeJson(boolean)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetPrintTreeJson(printTreeJson bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.printTreeJson = printTreeJson
	return c
}

// SetPrintAst prints a dot file describing the internal abstract syntax
// tree and exits.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setPrintAst(boolean)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetPrintAst(printAst bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.printAst = printAst
	return c
}

// SetJscompDevMode turns on extra validity checks.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setJscompDevMode(CompilerOptions.DevMode)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetJscompDevMode(jscompDevMode jscomp.CompilerOptions_DevMode) *AbstractCommandLineRunner_CommandLineConfig {
	c.jscompDevMode = jscompDevMode
	return c
}

// SetLoggingLevel sets the logging level (standard java.util.logging.Level
// values) for Compiler progress. Does not control errors or warnings for the
// JavaScript code under compilation.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setLoggingLevel(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetLoggingLevel(loggingLevel string) *AbstractCommandLineRunner_CommandLineConfig {
	c.loggingLevel = loggingLevel
	return c
}

// SetExterns sets the file containing JavaScript externs. You may specify
// multiple.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setExterns(List<String>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetExterns(externs []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.externs = append(c.externs[:0:0], externs...)
	return c
}

// SetMixedJsSources sets the JavaScript source file names, including .js and
// .zip files. You may specify multiple.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setMixedJsSources(List<FlagEntry<JsSourceType>>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetMixedJsSources(mixedJsSources []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]) *AbstractCommandLineRunner_CommandLineConfig {
	jutil.CheckNotNull(mixedJsSources) // Java: addAll(null) throws NullPointerException.
	c.mixedJsSources = append(c.mixedJsSources[:0:0], mixedJsSources...)
	return c
}

// SetDefaultToStdin sets whether to read a single source file from standard
// input if no input files are explicitly specified.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setDefaultToStdin()
func (c *AbstractCommandLineRunner_CommandLineConfig) SetDefaultToStdin() *AbstractCommandLineRunner_CommandLineConfig {
	c.defaultToStdin = true
	return c
}

// SetJsOutputFile sets the primary output filename. If not specified,
// output is written to stdout.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setJsOutputFile(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetJsOutputFile(jsOutputFile string) *AbstractCommandLineRunner_CommandLineConfig {
	c.jsOutputFile = jsOutputFile
	return c
}

// SetContinueSavedCompilationFileName sets the compiler to resume a saved
// compilation state from a file.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setContinueSavedCompilationFileName(String, int)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetContinueSavedCompilationFileName(
	fileName *string, restoredCompilationStage int) *AbstractCommandLineRunner_CommandLineConfig {
	if fileName != nil {
		jutil.CheckArgumentMsg1(
			restoredCompilationStage > 0 && restoredCompilationStage < 3,
			"invalid compilation stage: %s",
			restoredCompilationStage)
		c.continueSavedCompilationFileName = fileName
		c.restoredCompilationStage = restoredCompilationStage
	}
	return c
}

// Java: AbstractCommandLineRunner.CommandLineConfig.shouldSaveAfterStage1()
func (c *AbstractCommandLineRunner_CommandLineConfig) ShouldSaveAfterStage1() bool {
	if c.saveCompilationStateToFilename != nil && c.saveAfterCompilationStage == 1 {
		jutil.CheckStateMsg1(
			c.restoredCompilationStage < 0,
			"cannot perform stage 1 after restoring from stage %s",
			c.restoredCompilationStage)
		jutil.CheckStateMsg(
			c.continueSavedCompilationFileName == nil,
			"cannot restore a saved compilation and also save after stage 1")
		return true
	}
	return false
}

// Java: AbstractCommandLineRunner.CommandLineConfig.getContinueSavedCompilationFileName()
func (c *AbstractCommandLineRunner_CommandLineConfig) GetContinueSavedCompilationFileName() *string {
	return c.continueSavedCompilationFileName
}

// Java: AbstractCommandLineRunner.CommandLineConfig.shouldRestoreAndPerformStages2And3()
func (c *AbstractCommandLineRunner_CommandLineConfig) ShouldRestoreAndPerformStages2And3() bool {
	// We have a saved compilations state to restore
	return c.ShouldContinueCompilation() &&
		// the restored compilation stopped at stage 1
		c.restoredCompilationStage == 1 &&
		// we want to complete the rest of the compilation stages
		c.saveAfterCompilationStage < 0
}

// Java: AbstractCommandLineRunner.CommandLineConfig.shouldContinueCompilation()
func (c *AbstractCommandLineRunner_CommandLineConfig) ShouldContinueCompilation() bool {
	if c.continueSavedCompilationFileName != nil {
		jutil.CheckStateMsg1(
			c.restoredCompilationStage > 0 && c.restoredCompilationStage < 3,
			"invalid restored compilation stage: %s",
			c.restoredCompilationStage)
		return true
	}
	return false
}

// Java: AbstractCommandLineRunner.CommandLineConfig.shouldRestoreTypedAstsPerformStage2AndSave()
func (c *AbstractCommandLineRunner_CommandLineConfig) ShouldRestoreTypedAstsPerformStage2AndSave() bool {
	// We have a typed ast input list to parse
	return c.typedAstListInputFilename != nil &&
		// we want to stop and save after optimizations
		c.saveAfterCompilationStage == 2
}

// Java: AbstractCommandLineRunner.CommandLineConfig.shouldRestoreTypedAstsPerformStages2And3()
func (c *AbstractCommandLineRunner_CommandLineConfig) ShouldRestoreTypedAstsPerformStages2And3() bool {
	// We have a typed ast input list to parse
	return c.typedAstListInputFilename != nil &&
		// we do not want to stop and save after optimizations
		c.saveAfterCompilationStage == -1
}

// Java: AbstractCommandLineRunner.CommandLineConfig.shouldRestoreAndPerformStage2AndSave()
func (c *AbstractCommandLineRunner_CommandLineConfig) ShouldRestoreAndPerformStage2AndSave() bool {
	// We have a saved compilations state to restore
	return c.ShouldContinueCompilation() &&
		// the restored compilation stopped at stage 1
		c.restoredCompilationStage == 1 &&
		// we want to stop and save after stage 2
		c.saveAfterCompilationStage == 2
}

// Java: AbstractCommandLineRunner.CommandLineConfig.shouldRestoreAndPerformStage3()
func (c *AbstractCommandLineRunner_CommandLineConfig) ShouldRestoreAndPerformStage3() bool {
	if c.ShouldContinueCompilation() && c.restoredCompilationStage == 2 {
		// We have a saved compilations state from stage 2 to restore and continue
		jutil.CheckStateMsg1(
			c.saveAfterCompilationStage < 0,
			"request to save after stage %s is invalid when restoring from stage 2",
			c.saveAfterCompilationStage)
		return true
	}
	return false
}

// Java: AbstractCommandLineRunner.CommandLineConfig.shouldDoLateLocalization()
func (c *AbstractCommandLineRunner_CommandLineConfig) ShouldDoLateLocalization() bool {
	// The point of dividing checks, optimizations, and finalizations into different stages is to
	// localization work for last, so we avoid doing checks and optimizations separately for every
	// locale.
	// If we aren't doing finalizations as a separate final stage, then late localization is just
	// doing more work for a possibly-bigger compiled output (because code gets added after
	// optimizations have already executed).
	return c.ShouldRestoreAndPerformStage2AndSave() ||
		c.ShouldRestoreAndPerformStage3() ||
		c.ShouldRestoreTypedAstsPerformStage2AndSave()
}

// Java: AbstractCommandLineRunner.CommandLineConfig.shouldAlwaysGatherSourceMapInfo()
func (c *AbstractCommandLineRunner_CommandLineConfig) ShouldAlwaysGatherSourceMapInfo() bool {
	// If we're doing a partial compilation that isn't the final stage, we need to always gather
	// source map info in case the final stage requires it.
	return c.ShouldSaveAfterStage1() ||
		c.ShouldRestoreAndPerformStage2AndSave() ||
		c.ShouldRestoreTypedAstsPerformStage2AndSave()
}

// Java: AbstractCommandLineRunner.CommandLineConfig.setTypedAstListInputFilename(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetTypedAstListInputFilename(fileName *string) *AbstractCommandLineRunner_CommandLineConfig {
	c.typedAstListInputFilename = fileName
	return c
}

// SetSaveCompilationStateToFilename sets the compiler to perform the first
// phase and save the intermediate result to a file.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setSaveCompilationStateToFilename(String, int)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSaveCompilationStateToFilename(
	fileName *string, saveAfterCompilationStage int) *AbstractCommandLineRunner_CommandLineConfig {
	c.saveCompilationStateToFilename = fileName
	c.saveAfterCompilationStage = saveAfterCompilationStage
	return c
}

// Java: AbstractCommandLineRunner.CommandLineConfig.getSaveCompilationStateToFilename()
func (c *AbstractCommandLineRunner_CommandLineConfig) GetSaveCompilationStateToFilename() *string {
	return c.saveCompilationStateToFilename
}

// SetStage2SegmentOfCompilationToRun configures which Stage 2 segment
// (OPTIMIZATIONS, OPTIMIZATIONS_FIRST_HALF, or OPTIMIZATIONS_SECOND_HALF)
// RunCompilerPasses executes when Stage 2 runs.
func (c *AbstractCommandLineRunner_CommandLineConfig) SetStage2SegmentOfCompilationToRun(
	segment jscomp.CompilerOptions_SegmentOfCompilationToRun,
) *AbstractCommandLineRunner_CommandLineConfig {
	c.stage2SegmentOfCompilationToRun = segment
	return c
}

// GetStage2SegmentOfCompilationToRun returns the Stage 2 segment configured
// on this CommandLineConfig.
func (c *AbstractCommandLineRunner_CommandLineConfig) GetStage2SegmentOfCompilationToRun() jscomp.CompilerOptions_SegmentOfCompilationToRun {
	return c.stage2SegmentOfCompilationToRun
}

// SetChunk sets the JavaScript chunk specifications. The format is
// <name>:<num-js-files>[:[<dep>,...][:]]]. Chunk names must be unique. Each
// dep is the name of a chunk that this chunk depends on. Chunks must be
// listed in dependency order, and JS source files must be listed in the
// corresponding order. Where --chunk flags occur in relation to --js flags
// is unimportant
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setChunk(List<String>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetChunk(chunk []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.chunk = append(c.chunk[:0:0], chunk...)
	return c
}

// PipeFlagEntry holds an ordered key-value pair parsed from a delimited CLI
// flag (such as --source_map_input or --source_map_location_mapping).
type PipeFlagEntry struct {
	Key   string
	Value string
}

func checkNoDuplicatePipeFlagKeys(entries []PipeFlagEntry) {
	n := len(entries)
	if n <= 1 {
		return
	}
	if n > 8 {
		seen := make(map[string]struct{}, n)
		hasDup := false
		for _, e := range entries {
			if _, ok := seen[e.Key]; ok {
				hasDup = true
				break
			}
			seen[e.Key] = struct{}{}
		}
		if !hasDup {
			return
		}
	}
	for i := n - 1; i >= 0; i-- {
		k := entries[i].Key
		for j := i + 1; j < n; j++ {
			if entries[j].Key == k {
				panic(jutil.NewIllegalArgumentExceptionMsg("Multiple entries with same key: " +
					k + "=" + entries[j].Value + " and " + k + "=" + entries[i].Value))
			}
		}
	}
}

// Java: AbstractCommandLineRunner.CommandLineConfig.setSourceMapInputFiles(Map<String, String>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSourceMapInputFiles(sourceMapInputFiles []PipeFlagEntry) *AbstractCommandLineRunner_CommandLineConfig {
	c.sourceMapInputFiles = slices.Clone(sourceMapInputFiles)
	return c
}

// Java: AbstractCommandLineRunner.CommandLineConfig.setParseInlineSourceMaps(boolean)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetParseInlineSourceMaps(parseInlineSourceMaps bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.parseInlineSourceMaps = parseInlineSourceMaps
	return c
}

// SetExpectedDiagnostics sets the expected diagnostics in the format
// [(line,col)]CODE:regex. CODE is the JSCompiler DiagnosticType key. regex
// matches the description.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setExpectedDiagnostics(List<String>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetExpectedDiagnostics(expectedDiagnostics []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.expectedDiagnostics = append(c.expectedDiagnostics[:0:0], expectedDiagnostics...)
	return c
}

// SetVariableMapInputFile sets the file containing the serialized version of
// the variable renaming map produced by a previous compilation.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setVariableMapInputFile(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetVariableMapInputFile(variableMapInputFile string) *AbstractCommandLineRunner_CommandLineConfig {
	c.variableMapInputFile = variableMapInputFile
	return c
}

// SetPropertyMapInputFile sets the file containing the serialized version
// of the property renaming map produced by a previous compilation.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setPropertyMapInputFile(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetPropertyMapInputFile(propertyMapInputFile string) *AbstractCommandLineRunner_CommandLineConfig {
	c.propertyMapInputFile = propertyMapInputFile
	return c
}

// SetVariableMapOutputFile sets the file where the serialized version of
// the variable renaming map produced should be saved.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setVariableMapOutputFile(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetVariableMapOutputFile(variableMapOutputFile string) *AbstractCommandLineRunner_CommandLineConfig {
	c.variableMapOutputFile = variableMapOutputFile
	return c
}

// SetCreateNameMapFiles: if true, variable renaming and property renaming
// map files will be produced as {binary name}_vars_map.out and {binary
// name}_props_map.out. Note that this flag cannot be used in conjunction
// with either variable_map_output_file or property_map_output_file.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setCreateNameMapFiles(boolean)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetCreateNameMapFiles(createNameMapFiles bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.createNameMapFiles = createNameMapFiles
	return c
}

// SetPropertyMapOutputFile sets the file where the serialized version of
// the property renaming map produced should be saved.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setPropertyMapOutputFile(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetPropertyMapOutputFile(propertyMapOutputFile string) *AbstractCommandLineRunner_CommandLineConfig {
	c.propertyMapOutputFile = propertyMapOutputFile
	return c
}

// SetStringMapOutputFile sets the file where the serialized version of the
// string map produced by the ReplaceStrings pass should be saved.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setStringMapOutputFile(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetStringMapOutputFile(stringMapOutputPath string) *AbstractCommandLineRunner_CommandLineConfig {
	c.stringMapOutputPath = stringMapOutputPath
	return c
}

// Java: AbstractCommandLineRunner.CommandLineConfig.setInstrumentationMappingFile(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetInstrumentationMappingFile(instrumentationMappingFile string) *AbstractCommandLineRunner_CommandLineConfig {
	c.instrumentationMappingFile = instrumentationMappingFile
	return c
}

// SetCodingConvention sets rules and conventions to enforce.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setCodingConvention(CodingConvention)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetCodingConvention(codingConvention jscomp.CodingConvention) *AbstractCommandLineRunner_CommandLineConfig {
	c.codingConvention = codingConvention
	return c
}

// SetSummaryDetailLevel controls how detailed the compilation summary is.
// Values: 0 (never print summary), 1 (print summary only if there are errors
// or warnings), 2 (print summary if type checking is on), 3 (always print
// summary). The default level is 1
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setSummaryDetailLevel(int)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSummaryDetailLevel(summaryDetailLevel int) *AbstractCommandLineRunner_CommandLineConfig {
	c.summaryDetailLevel = summaryDetailLevel
	return c
}

// SetOutputWrapper interpolates output into this string at the place
// denoted by the marker token %output%, or %output|jsstring%.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setOutputWrapper(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetOutputWrapper(outputWrapper string) *AbstractCommandLineRunner_CommandLineConfig {
	c.outputWrapper = outputWrapper
	return c
}

// SetChunkWrapper sets an output wrapper for a JavaScript module (optional).
// See the flag description for formatting requirements.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setChunkWrapper(List<String>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetChunkWrapper(chunkWrapper []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.chunkWrapper = append(c.chunkWrapper[:0:0], chunkWrapper...)
	return c
}

// SetChunkOutputPathPrefix sets the prefix for filenames of compiled
// chunks. <chunk-name>.js will be appended to this prefix. Directories will
// be created as needed. Use with --chunk
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setChunkOutputPathPrefix(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetChunkOutputPathPrefix(chunkOutputPathPrefix string) *AbstractCommandLineRunner_CommandLineConfig {
	c.chunkOutputPathPrefix = chunkOutputPathPrefix
	return c
}

// SetChunkOutputFiles sets the output file name for a JavaScript chunk
// (optional). See the flag description for formatting requirements.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setChunkOutputFiles(List<String>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetChunkOutputFiles(chunkOutputFiles []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.chunkOutputFiles = append(c.chunkOutputFiles[:0:0], chunkOutputFiles...)
	return c
}

// SetChunkConformanceFiles sets the conformance report file name for a
// JavaScript chunk (optional). See the flag description for formatting
// requirements.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setChunkConformanceFiles(List<String>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetChunkConformanceFiles(chunkConformanceFiles []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.chunkConformanceFiles = append(c.chunkConformanceFiles[:0:0], chunkConformanceFiles...)
	return c
}

// SetCreateSourceMap: if specified, a source map file mapping the generated
// source files back to the original source file will be output to the
// specified path. The %outname% placeholder will expand to the name of the
// output file that the source map corresponds to.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setCreateSourceMap(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetCreateSourceMap(createSourceMap string) *AbstractCommandLineRunner_CommandLineConfig {
	c.createSourceMap = createSourceMap
	return c
}

// SetSourceMapDetailLevel sets the detail supplied in the source map file,
// if generated.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setSourceMapDetailLevel(SourceMap.DetailLevel)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSourceMapDetailLevel(level jscomp.SourceMap_DetailLevel) *AbstractCommandLineRunner_CommandLineConfig {
	c.sourceMapDetailLevel = level
	return c
}

// SetSourceMapFormat sets the source map format to use, if generated.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setSourceMapFormat(SourceMap.Format)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSourceMapFormat(format jscomp.SourceMap_Format) *AbstractCommandLineRunner_CommandLineConfig {
	c.sourceMapFormat = format
	return c
}

// SetSourceMapLocationMappings sets the source map location mappings to
// use, if generated.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setSourceMapLocationMappings(List<SourceMap.LocationMapping>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSourceMapLocationMappings(
	locationMappings []jscomp.SourceMap_LocationMapping) *AbstractCommandLineRunner_CommandLineConfig {

	c.sourceMapLocationMappings = slices.Clone(locationMappings)
	return c
}

// SetApplyInputSourceMaps sets whether to apply input source maps to the
// output, i.e. map back to original inputs from input files that have source
// maps applied to them.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setApplyInputSourceMaps(boolean)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetApplyInputSourceMaps(applyInputSourceMaps bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.applyInputSourceMaps = applyInputSourceMaps
	return c
}

// SetWarningGuards adds warning guards.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setWarningGuards(List<FlagEntry<CheckLevel>>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetWarningGuards(warningGuards []*AbstractCommandLineRunner_FlagEntry[jscomp.CheckLevel]) *AbstractCommandLineRunner_CommandLineConfig {
	c.warningGuards = append(c.warningGuards[:0:0], warningGuards...)
	return c
}

// SetDefine overrides the value of a variable annotated @define. The format
// is <name>[=<val>], where <name> is the name of a @define variable and
// <val> is a boolean, number, or a single-quoted string that contains no
// single quotes. If [=<val>] is omitted, the variable is marked true
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setDefine(List<String>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetDefine(define []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.define = append(c.define[:0:0], define...)
	return c
}

// SetBrowserFeaturesetYear indicates target browser's year.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setBrowserFeaturesetYear(Integer)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetBrowserFeaturesetYear(browserFeaturesetYear int) *AbstractCommandLineRunner_CommandLineConfig {
	c.browserFeaturesetYear = browserFeaturesetYear
	return c
}

// SetTweakProcessing sets the kind of processing to do for goog.tweak
// functions.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setTweakProcessing(TweakProcessing)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetTweakProcessing(tweakProcessing jscomp.CompilerOptions_TweakProcessing) *AbstractCommandLineRunner_CommandLineConfig {
	c.tweakProcessing = tweakProcessing
	return c
}

// SetCharset sets the input charset for all files.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setCharset(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetCharset(charset string) *AbstractCommandLineRunner_CommandLineConfig {
	c.charset = charset
	return c
}

// SetDependencyOptions sets the dependency management options.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setDependencyOptions(DependencyOptions)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetDependencyOptions(dependencyOptions *jscomp.DependencyOptions) *AbstractCommandLineRunner_CommandLineConfig {
	c.dependencyOptions = dependencyOptions
	return c
}

// SetOutputManifest sets whether to print output manifest files. Filter out
// empty file names.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setOutputManifest(List<String>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetOutputManifest(outputManifests []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.outputManifests = []string{}
	for _, manifestName := range outputManifests {
		if manifestName != "" {
			c.outputManifests = append(c.outputManifests, manifestName)
		}
	}
	c.outputManifests = slices.Clone(c.outputManifests)
	return c
}

// SetOutputChunkDependencies sets whether a JSON file representing the
// dependencies between modules should be created.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setOutputChunkDependencies(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetOutputChunkDependencies(outputChunkDependencies *string) *AbstractCommandLineRunner_CommandLineConfig {
	c.outputChunkDependencies = outputChunkDependencies
	return c
}

// SetOutputBundle sets whether to print output bundle files.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setOutputBundle(List<String>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetOutputBundle(outputBundles []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.outputBundles = outputBundles
	return c
}

// SetSkipNormalOutputs sets whether the normal outputs of compilation should
// be skipped.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setSkipNormalOutputs(boolean)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSkipNormalOutputs(skipNormalOutputs bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.skipNormalOutputs = skipNormalOutputs
	return c
}

// SetManifestMaps sets the execPath:rootRelativePath mappings.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setManifestMaps(List<String>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetManifestMaps(manifestMaps []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.manifestMaps = manifestMaps
	return c
}

// SetProcessCommonJSModules sets whether to process CommonJS modules,.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setProcessCommonJSModules(boolean)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetProcessCommonJSModules(processCommonJSModules bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.processCommonJSModules = processCommonJSModules
	return c
}

// SetModuleRoots sets the module roots.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setModuleRoots(List<String>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetModuleRoots(jsModuleRoots []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.moduleRoots = jsModuleRoots
	return c
}

// SetWarningsAllowlistFile sets a allowlist file that suppresses warnings.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setWarningsAllowlistFile(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetWarningsAllowlistFile(fileName string) *AbstractCommandLineRunner_CommandLineConfig {
	c.warningsAllowFile = fileName
	return c
}

// SetHideWarningsFor sets the paths for which warnings will be hidden.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setHideWarningsFor(List<String>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetHideWarningsFor(hideWarningsFor []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.hideWarningsFor = hideWarningsFor
	return c
}

// SetAngularPass sets whether to process AngularJS-specific annotations.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setAngularPass(boolean)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetAngularPass(angularPass bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.angularPass = angularPass
	return c
}

// Java: AbstractCommandLineRunner.CommandLineConfig.setJsonStreamMode(JsonStreamMode)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetJsonStreamMode(mode jscomp.CompilerOptions_JsonStreamMode) *AbstractCommandLineRunner_CommandLineConfig {
	c.jsonStreamMode = mode
	return c
}

// AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption is the set
// of options that can be used with the --formatting flag.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.ErrorFormatOption
type AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption int32

const (
	AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_STANDARD AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption = iota
	AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_JSON
)

var abstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_names = [...]string{"STANDARD", "JSON"}

// AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_values is
// ErrorFormatOption.values().
func AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_values() []AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption {
	return []AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption{
		AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_STANDARD,
		AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_JSON,
	}
}

// String is Enum.name().
func (e AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) String() string {
	return abstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_names[e]
}

// Ordinal is Enum.ordinal().
func (e AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) Ordinal() int { return int(e) }

// Java: AbstractCommandLineRunner.CommandLineConfig.setErrorFormat(ErrorFormatOption)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetErrorFormat(errorFormat AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) *AbstractCommandLineRunner_CommandLineConfig {
	c.errorFormat = errorFormat
	return c
}

// Java: AbstractCommandLineRunner.CommandLineConfig.setJsonWarningsFile(String)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetJsonWarningsFile(jsonWarningsFile string) *AbstractCommandLineRunner_CommandLineConfig {
	c.jsonWarningsFile = jsonWarningsFile
	return c
}

// Java: AbstractCommandLineRunner.CommandLineConfig.setIncludeTrailingNewline(boolean)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetIncludeTrailingNewline(includeTrailingNewline bool) *AbstractCommandLineRunner_CommandLineConfig {
	c.includeTrailingNewline = includeTrailingNewline
	return c
}

// AbstractCommandLineRunner_JsonFileSpec is the representation of a source
// file from an encoded json stream input.
//
// Java: AbstractCommandLineRunner.JsonFileSpec
type AbstractCommandLineRunner_JsonFileSpec struct {
	Src       *string
	Path      *string
	SourceMap *string
	WebpackId *string
}

// Graal requires a non-arg constructor for use with GSON
// See https://github.com/oracle/graal/issues/680
//
// Java: AbstractCommandLineRunner.JsonFileSpec()
func newAbstractCommandLineRunner_JsonFileSpecDefault() *AbstractCommandLineRunner_JsonFileSpec {
	return NewAbstractCommandLineRunner_JsonFileSpecWithSourceMapAndWebpackId(nil, nil, nil, nil)
}

// Java: AbstractCommandLineRunner.JsonFileSpec(String, String)
func NewAbstractCommandLineRunner_JsonFileSpec(src string, path string) *AbstractCommandLineRunner_JsonFileSpec {
	return NewAbstractCommandLineRunner_JsonFileSpecWithSourceMapAndWebpackId(&src, &path, nil, nil)
}

// Java: AbstractCommandLineRunner.JsonFileSpec(String, String, String)
func NewAbstractCommandLineRunner_JsonFileSpecWithSourceMap(src string, path string, sourceMap string) *AbstractCommandLineRunner_JsonFileSpec {
	return NewAbstractCommandLineRunner_JsonFileSpecWithSourceMapAndWebpackId(&src, &path, &sourceMap, nil)
}

// Java: AbstractCommandLineRunner.JsonFileSpec(String, String, String, String)
func NewAbstractCommandLineRunner_JsonFileSpecWithSourceMapAndWebpackId(
	src *string,
	path *string,
	sourceMap *string,
	webpackId *string) *AbstractCommandLineRunner_JsonFileSpec {
	return &AbstractCommandLineRunner_JsonFileSpec{
		Src:       src,
		Path:      path,
		SourceMap: sourceMap,
		WebpackId: webpackId,
	}
}

// Java: AbstractCommandLineRunner.JsonFileSpec.getSrc()
func (s *AbstractCommandLineRunner_JsonFileSpec) GetSrc() *string {
	return s.Src
}

// Java: AbstractCommandLineRunner.JsonFileSpec.getPath()
func (s *AbstractCommandLineRunner_JsonFileSpec) GetPath() *string {
	return s.Path
}

// Java: AbstractCommandLineRunner.JsonFileSpec.getSourceMap()
func (s *AbstractCommandLineRunner_JsonFileSpec) GetSourceMap() *string {
	return s.SourceMap
}

// Java: AbstractCommandLineRunner.JsonFileSpec.getWebpackId()
func (s *AbstractCommandLineRunner_JsonFileSpec) GetWebpackId() *string {
	return s.WebpackId
}

// Java: AbstractCommandLineRunner.JsonFileSpec.setSourceMap(String)
func (s *AbstractCommandLineRunner_JsonFileSpec) SetSourceMap(map_ string) {
	s.SourceMap = &map_
}

// AbstractCommandLineRunner_JsSourceType is the flag types for JavaScript
// source files.
//
// Java: AbstractCommandLineRunner.JsSourceType
type AbstractCommandLineRunner_JsSourceType int32

const (
	AbstractCommandLineRunner_JsSourceType_EXTERN AbstractCommandLineRunner_JsSourceType = iota
	AbstractCommandLineRunner_JsSourceType_JS
	AbstractCommandLineRunner_JsSourceType_JS_ZIP
	AbstractCommandLineRunner_JsSourceType_WEAKDEP
	AbstractCommandLineRunner_JsSourceType_IJS
)

var abstractCommandLineRunner_JsSourceType_names = [...]string{"EXTERN", "JS", "JS_ZIP", "WEAKDEP", "IJS"}

// abstractCommandLineRunner_JsSourceType_flagName is the field flagName,
// indexed by ordinal.
//
// Java: AbstractCommandLineRunner.JsSourceType.flagName
var abstractCommandLineRunner_JsSourceType_flagName = [...]string{"extern", "js", "jszip", "weakdep", "ijs"}

// String is Enum.name().
func (t AbstractCommandLineRunner_JsSourceType) String() string {
	return abstractCommandLineRunner_JsSourceType_names[t]
}

// Ordinal is Enum.ordinal().
func (t AbstractCommandLineRunner_JsSourceType) Ordinal() int { return int(t) }

// FlagName reads the field flagName.
func (t AbstractCommandLineRunner_JsSourceType) FlagName() string {
	return abstractCommandLineRunner_JsSourceType_flagName[t]
}

// abstractCommandLineRunner_flag is the constraint on FlagEntry's type
// parameter: Java compares flags with equals() and prints them with
// toString(); the only flags are enums.
type abstractCommandLineRunner_flag interface {
	comparable
	String() string
	Ordinal() int
}

// AbstractCommandLineRunner_FlagEntry is a pair from flag to its value.
//
// Java: AbstractCommandLineRunner.FlagEntry<T>
type AbstractCommandLineRunner_FlagEntry[T abstractCommandLineRunner_flag] struct {
	flag  T
	value string
}

// Java: AbstractCommandLineRunner.FlagEntry(T, String)
func NewAbstractCommandLineRunner_FlagEntry[T abstractCommandLineRunner_flag](flag T, value string) *AbstractCommandLineRunner_FlagEntry[T] {
	return &AbstractCommandLineRunner_FlagEntry[T]{flag: flag, value: value}
}

// Java: AbstractCommandLineRunner.FlagEntry.equals(Object)
func (e *AbstractCommandLineRunner_FlagEntry[T]) Equals(o any) bool {
	if that, ok := o.(*AbstractCommandLineRunner_FlagEntry[T]); ok && that != nil {
		return that.flag == e.flag && that.value == e.value
	}
	return false
}

// HashCode is flag.hashCode() + value.hashCode().
//
// Go-only: language-adapter — a Java enum's hashCode() is its identity hash
// code, which no port can reproduce (and no Closure code depends on); the
// ordinal stands in for it.
//
// Java: AbstractCommandLineRunner.FlagEntry.hashCode()
func (e *AbstractCommandLineRunner_FlagEntry[T]) HashCode() int32 {
	return int32(e.flag.Ordinal()) + jutil.StringHashCode(e.value)
}

// Java: AbstractCommandLineRunner.FlagEntry.toString()
func (e *AbstractCommandLineRunner_FlagEntry[T]) String() string {
	return e.flag.String() + "=" + e.value
}

// Java: AbstractCommandLineRunner.FlagEntry.getFlag()
func (e *AbstractCommandLineRunner_FlagEntry[T]) GetFlag() T {
	return e.flag
}

// Java: AbstractCommandLineRunner.FlagEntry.getValue()
func (e *AbstractCommandLineRunner_FlagEntry[T]) GetValue() string {
	return e.value
}

// AbstractCommandLineRunner_JsChunkSpec represents a specification for a
// serving chunk.
//
// Java: AbstractCommandLineRunner.JsChunkSpec
type AbstractCommandLineRunner_JsChunkSpec struct {
	name string
	// Number of input files, including js and zip files.
	numInputs int
	deps      []string
	// Number of input js files. All zip files should be expanded.
	numJsFiles int
}

// Java: AbstractCommandLineRunner.JsChunkSpec(String, int, ImmutableList<String>)
func newAbstractCommandLineRunner_JsChunkSpec(name string, numInputs int, deps []string) *AbstractCommandLineRunner_JsChunkSpec {
	return &AbstractCommandLineRunner_JsChunkSpec{
		name:       name,
		numInputs:  numInputs,
		deps:       deps,
		numJsFiles: numInputs,
	}
}

// AbstractCommandLineRunner_JsChunkSpec_Create returns a parsed chunk spec.
//
// specString's format is: name:num-js-files[:[dep,...][:]]. Chunk names
// must not contain the ':' character. isFirstChunk is whether the spec is
// for the first chunk.
//
// Java: AbstractCommandLineRunner.JsChunkSpec.create(String, boolean)
func AbstractCommandLineRunner_JsChunkSpec_Create(specString string, isFirstChunk bool) *AbstractCommandLineRunner_JsChunkSpec {
	// Format is "<name>:<num-js-files>[:[<dep>,...][:]]".
	parts := jutil.StringSplitLiteral(specString, ":")
	if len(parts) < 2 || len(parts) > 4 {
		panic(jscomp.NewFlagUsageException(
			"Expected 2-4 colon-delimited parts in " + "chunk spec: " + specString))
	}

	// Parse chunk name.
	name := parts[0]

	// Parse chunk dependencies.
	deps := []string{}
	if len(parts) > 2 && parts[2] != "" {
		deps = jutil.StringSplitLiteral(parts[2], ",")
	}

	// Parse chunk inputs.
	numInputs := -1
	if parsed, ok := jutil.TryIntegerParseInt(parts[1], 10); ok {
		numInputs = int(parsed)
	}

	// We will allow chunks of zero input.
	if numInputs < 0 {
		// A size of 'auto' is only allowed on the base chunk if
		// and it must also be the first chunk
		if len(parts) == 2 && "auto" == parts[1] {
			if !isFirstChunk {
				panic(jscomp.NewFlagUsageException(
					"Invalid JS file count '" +
						parts[1] +
						"' for chunk: " +
						name +
						". Only the first chunk may specify " +
						"a size of 'auto' and it must have no dependencies."))
			}
		} else {
			panic(jscomp.NewFlagUsageException(
				"Invalid JS file count '" + parts[1] + "' for chunk: " + name))
		}
	}

	return newAbstractCommandLineRunner_JsChunkSpec(name, numInputs, slices.Clone(deps))
}

// Java: AbstractCommandLineRunner.JsChunkSpec.getName()
func (s *AbstractCommandLineRunner_JsChunkSpec) GetName() string {
	return s.name
}

// Java: AbstractCommandLineRunner.JsChunkSpec.getNumInputs()
func (s *AbstractCommandLineRunner_JsChunkSpec) GetNumInputs() int {
	return s.numInputs
}

// Java: AbstractCommandLineRunner.JsChunkSpec.getDeps()
func (s *AbstractCommandLineRunner_JsChunkSpec) GetDeps() []string {
	return s.deps
}

// Java: AbstractCommandLineRunner.JsChunkSpec.getNumJsFiles()
func (s *AbstractCommandLineRunner_JsChunkSpec) GetNumJsFiles() int {
	return s.numJsFiles
}

// AbstractCommandLineRunner_SystemExitCodeReceiver passes the exit code to
// System.exit.
//
// Java: AbstractCommandLineRunner.SystemExitCodeReceiver
type AbstractCommandLineRunner_SystemExitCodeReceiver struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// Java: AbstractCommandLineRunner.SystemExitCodeReceiver.INSTANCE
var AbstractCommandLineRunner_SystemExitCodeReceiver_INSTANCE = newAbstractCommandLineRunner_SystemExitCodeReceiver()

// Java: AbstractCommandLineRunner.SystemExitCodeReceiver()
func newAbstractCommandLineRunner_SystemExitCodeReceiver() *AbstractCommandLineRunner_SystemExitCodeReceiver {
	// singleton
	return &AbstractCommandLineRunner_SystemExitCodeReceiver{}
}

// Java: AbstractCommandLineRunner.SystemExitCodeReceiver.apply(Integer)
func (*AbstractCommandLineRunner_SystemExitCodeReceiver) Apply(exitCode int) {
	exitCodeValue := int32(exitCode)
	// Don't spuriously report success.
	// Posix conventions only guarantee that 8b are significant.
	exitCodeByte := int8(exitCodeValue)
	if exitCodeByte == 0 && exitCodeValue != 0 {
		exitCodeByte = -1
	}
	jutil.System_Exit(int(exitCodeByte))
}
