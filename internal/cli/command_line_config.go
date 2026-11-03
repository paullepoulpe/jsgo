// Port of com/google/javascript/jscomp/AbstractCommandLineRunner.java @ closure-compiler v20261028

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

	"github.com/paullepoulpe/jsgo/internal/base"
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

	// When present specifies a file containing saved compiler state to
	// restore and continue compiling.
	continueSavedCompilationFileName base.OptString

	// When > 0 indicates the stage at which compilation stopped for the
	// compilation state that is being restored from a save file.
	restoredCompilationStage int

	// When > 0 indicates the stage at which compilation should stop and the
	// compilation state be stored into a save file.
	saveAfterCompilationStage int

	stage2SegmentOfCompilationToRun jscomp.CompilerOptions_SegmentOfCompilationToRun

	typedAstListInputFilename base.OptString

	saveCompilationStateToFilename base.OptString

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

	outputChunkDependencies string

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
		jscompDevMode:                   jscomp.CompilerOptions_DevMode_OFF,
		loggingLevel:                    jutil.Level_WARNING.GetName(),
		externs:                         []string{},
		mixedJsSources:                  []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]{},
		restoredCompilationStage:        -1,
		saveAfterCompilationStage:       -1,
		stage2SegmentOfCompilationToRun: jscomp.CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS,
		chunk:                           []string{},
		sourceMapInputFiles:             []PipeFlagEntry{},
		expectedDiagnostics:             []string{},
		codingConvention:                jscomp.CodingConventions_GetDefault(),
		summaryDetailLevel:              1,
		chunkWrapper:                    []string{},
		chunkOutputFiles:                []string{},
		chunkConformanceFiles:           []string{},
		sourceMapDetailLevel:            jscomp.SourceMap_DetailLevel_ALL,
		sourceMapFormat:                 jscomp.SourceMap_Format_DEFAULT,
		sourceMapLocationMappings:       []jscomp.SourceMap_LocationMapping{},
		warningGuards:                   []*AbstractCommandLineRunner_FlagEntry[jscomp.CheckLevel]{},
		define:                          []string{},
		tweakProcessing:                 jscomp.CompilerOptions_TweakProcessing_OFF,
		outputManifests:                 []string{},
		outputBundles:                   []string{},
		manifestMaps:                    []string{},
		moduleRoots:                     []string{jscomp.ModuleLoader_DEFAULT_FILENAME_PREFIX},
		hideWarningsFor:                 []string{},
		jsonStreamMode:                  jscomp.CompilerOptions_JsonStreamMode_NONE,
		errorFormat:                     AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_STANDARD,
		includeTrailingNewline:          true,
	}
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

// SetContinueSavedCompilationFileName sets the compiler to resume a saved
// compilation state from a file.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setContinueSavedCompilationFileName(String, int)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetContinueSavedCompilationFileName(
	fileName *string, restoredCompilationStage int) *AbstractCommandLineRunner_CommandLineConfig {
	return c.setContinueSavedCompilationOpt(base.OptStringFromPtr(fileName), restoredCompilationStage)
}

func (c *AbstractCommandLineRunner_CommandLineConfig) setContinueSavedCompilationOpt(
	fileName base.OptString, restoredCompilationStage int) *AbstractCommandLineRunner_CommandLineConfig {
	if fileName.IsPresent() {
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
	if c.saveCompilationStateToFilename.IsPresent() && c.saveAfterCompilationStage == 1 {
		jutil.CheckStateMsg1(
			c.restoredCompilationStage < 0,
			"cannot perform stage 1 after restoring from stage %s",
			c.restoredCompilationStage)
		jutil.CheckStateMsg(
			!c.continueSavedCompilationFileName.IsPresent(),
			"cannot restore a saved compilation and also save after stage 1")
		return true
	}
	return false
}

// Java: AbstractCommandLineRunner.CommandLineConfig.getContinueSavedCompilationFileName()
func (c *AbstractCommandLineRunner_CommandLineConfig) GetContinueSavedCompilationFileName() *string {
	return c.continueSavedCompilationFileName.Ptr()
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
	if c.continueSavedCompilationFileName.IsPresent() {
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
	return c.typedAstListInputFilename.IsPresent() &&
		// we want to stop and save after optimizations
		c.saveAfterCompilationStage == 2
}

// Java: AbstractCommandLineRunner.CommandLineConfig.shouldRestoreTypedAstsPerformStages2And3()
func (c *AbstractCommandLineRunner_CommandLineConfig) ShouldRestoreTypedAstsPerformStages2And3() bool {
	// We have a typed ast input list to parse
	return c.typedAstListInputFilename.IsPresent() &&
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

// SetSaveCompilationStateToFilename sets the compiler to perform the first
// phase and save the intermediate result to a file.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setSaveCompilationStateToFilename(String, int)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSaveCompilationStateToFilename(
	fileName *string, saveAfterCompilationStage int) *AbstractCommandLineRunner_CommandLineConfig {
	return c.setSaveCompilationStateOpt(base.OptStringFromPtr(fileName), saveAfterCompilationStage)
}

func (c *AbstractCommandLineRunner_CommandLineConfig) setSaveCompilationStateOpt(
	fileName base.OptString, saveAfterCompilationStage int) *AbstractCommandLineRunner_CommandLineConfig {
	c.saveCompilationStateToFilename = fileName
	c.saveAfterCompilationStage = saveAfterCompilationStage
	return c
}

// Java: AbstractCommandLineRunner.CommandLineConfig.getSaveCompilationStateToFilename()
func (c *AbstractCommandLineRunner_CommandLineConfig) GetSaveCompilationStateToFilename() *string {
	return c.saveCompilationStateToFilename.Ptr()
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

// Java: AbstractCommandLineRunner.CommandLineConfig.setSourceMapInputFiles(Map<String, String>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSourceMapInputFiles(sourceMapInputFiles []PipeFlagEntry) *AbstractCommandLineRunner_CommandLineConfig {
	c.sourceMapInputFiles = slices.Clone(sourceMapInputFiles)
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

// SetChunkWrapper sets an output wrapper for a JavaScript module (optional).
// See the flag description for formatting requirements.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setChunkWrapper(List<String>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetChunkWrapper(chunkWrapper []string) *AbstractCommandLineRunner_CommandLineConfig {
	c.chunkWrapper = append(c.chunkWrapper[:0:0], chunkWrapper...)
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

// SetOutputManifest sets whether to print output manifest files. Filter out
// empty file names.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setOutputManifest(List<String>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetOutputManifest(outputManifests []string) *AbstractCommandLineRunner_CommandLineConfig {
	n := 0
	for _, manifestName := range outputManifests {
		if manifestName != "" {
			n++
		}
	}
	if n == 0 {
		c.outputManifests = []string{}
		return c
	}
	manifests := make([]string, 0, n)
	for _, manifestName := range outputManifests {
		if manifestName != "" {
			manifests = append(manifests, manifestName)
		}
	}
	c.outputManifests = manifests
	return c
}
