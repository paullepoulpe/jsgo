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
