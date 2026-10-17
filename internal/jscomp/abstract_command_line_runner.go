// Port of com/google/javascript/jscomp/AbstractCommandLineRunner.java @ closure-compiler v20260415

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
	"archive/zip"
	"bytes"
	"io"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/resources"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AbstractCommandLineRunner implementations translate flags into Java API
// calls on the Compiler. AbstractCompiler contains common flags and logic to
// make that happen.
//
// This class may be extended and used to create other Java classes that
// behave the same as running the Compiler from the command line. Example:
//
//	class MyCommandLineRunner extends
//	    AbstractCommandLineRunner<MyCompiler, MyOptions> {
//	  MyCommandLineRunner(String[] args) {
//	    super(args);
//	  }
//
//	  @Override
//	  protected MyOptions createOptions() {
//	    MyOptions options = new MyOptions();
//	    CompilerFlagTranslator.setOptionsFromFlags(options);
//	    addMyCrazyCompilerPassThatOutputsAnExtraFile(options);
//	    return options;
//	  }
//
//	  @Override
//	  protected MyCompiler createCompiler() {
//	    return new MyCompiler();
//	  }
//
//	  public static void main(String[] args) {
//	    (new MyCommandLineRunner(args)).run();
//	  }
//	}
//
// Go-only: language-adapter — the type parameters <A extends Compiler, B
// extends CompilerOptions> are *Compiler and *CompilerOptions: their only
// instantiation is CommandLineRunner's, and Go cannot bound a type parameter
// by the subclasses of a struct type.
//
// Java: abstract class AbstractCommandLineRunner<A extends Compiler, B extends CompilerOptions>
type AbstractCommandLineRunner interface {
	GetAbstractCommandLineRunnerBase() *AbstractCommandLineRunnerBase

	// Java: AbstractCommandLineRunner.setExitCodeReceiver(Function<Integer, Void>)
	SetExitCodeReceiver(newExitCodeReceiver func(int))
	// Java: AbstractCommandLineRunner.isInTestMode()
	IsInTestMode() bool
	// Java: AbstractCommandLineRunner.getCommandLineConfig()
	GetCommandLineConfig() *AbstractCommandLineRunner_CommandLineConfig
	// Java: AbstractCommandLineRunner.createCompiler()
	CreateCompiler() *Compiler
	// Java: AbstractCommandLineRunner.prepForBundleAndAppendTo(Appendable, CompilerInput, String)
	PrepForBundleAndAppendTo(out jutil.Appendable, input *CompilerInput, content string)
	// Java: AbstractCommandLineRunner.appendRuntimeTo(Appendable)
	AppendRuntimeTo(out jutil.Appendable)
	// Java: AbstractCommandLineRunner.createOptions()
	CreateOptions() *CompilerOptions
	// Java: AbstractCommandLineRunner.getDiagnosticGroups()
	GetDiagnosticGroups() *DiagnosticGroups
	// Java: AbstractCommandLineRunner.addAllowlistWarningsGuard(CompilerOptions, File)
	AddAllowlistWarningsGuard(options *CompilerOptions, allowlistFile *jutil.File)
	// Java: AbstractCommandLineRunner.setRunOptions(CompilerOptions)
	SetRunOptions(options *CompilerOptions)
	// Java: AbstractCommandLineRunner.getVersionText()
	GetVersionText() string
	// Java: AbstractCommandLineRunner.parseJsonFilesFromInputStream()
	ParseJsonFilesFromInputStream() []*AbstractCommandLineRunner_JsonFileSpec
	// Java: AbstractCommandLineRunner.createInputs(List<FlagEntry<JsSourceType>>, List<JsonFileSpec>, boolean, List<JsChunkSpec>)
	CreateInputsWithJsonFilesAndAllowStdIn(
		files []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType],
		jsonFiles []*AbstractCommandLineRunner_JsonFileSpec,
		allowStdIn bool,
		jsChunkSpecs []*AbstractCommandLineRunner_JsChunkSpec) []*SourceFile
	// Java: AbstractCommandLineRunner.checkChunkName(String)
	CheckChunkName(name string)
	// Java: AbstractCommandLineRunner.getChunkOutputFileName(JSChunk)
	GetChunkOutputFileName(m *JSChunk) string
	// Java: AbstractCommandLineRunner.getChunkConformanceFileName(JSChunk)
	GetChunkConformanceFileName(m *JSChunk) *string
	// Java: AbstractCommandLineRunner.writeChunkOutput(String, Appendable, LicenseTracker, JSChunk)
	WriteChunkOutput(fileName string, out jutil.Appendable, lt CodePrinter_LicenseTracker, m *JSChunk)
	// Java: AbstractCommandLineRunner.writeOutput(Appendable, Compiler, LicenseTracker, JSChunk, String, String, Function<String, String>, String)
	WriteOutput(
		out jutil.Appendable,
		compiler *Compiler,
		licenseTracker CodePrinter_LicenseTracker,
		chunk *JSChunk,
		wrapper string,
		codePlaceholder string,
		escaper func(string) string,
		filename string)
	// Java: AbstractCommandLineRunner.writeOutput(Appendable, Compiler, String, String, String, Function<String, String>, String)
	WriteOutputOfCode(
		out jutil.Appendable,
		compiler *Compiler,
		code string,
		wrapper string,
		codePlaceholder string,
		escaper func(string) string,
		filename string)
	// Java: AbstractCommandLineRunner.doRun()
	DoRun() int
	// Java: AbstractCommandLineRunner.getCompileMetricsRecorder()
	GetCompileMetricsRecorder() CompileMetricsRecorderInterface
	// Java: AbstractCommandLineRunner.processResults(Result, List<JSChunk>, B)
	ProcessResults(result *Result, chunks []*JSChunk, options *CompilerOptions) int
	// Java: AbstractCommandLineRunner.getJavascriptEscaper()
	GetJavascriptEscaper() func(string) string
	// Java: AbstractCommandLineRunner.outputSingleBinary(B)
	OutputSingleBinary(options *CompilerOptions)
	// Java: AbstractCommandLineRunner.createJsonFile(B, String, Function<String, String>)
	CreateJsonFile(options *CompilerOptions, outputMarker string, escaper func(string) string) *AbstractCommandLineRunner_JsonFileSpec
	// Java: AbstractCommandLineRunner.outputJsonStream()
	OutputJsonStream()
	// Java: AbstractCommandLineRunner.createExterns(CompilerOptions)
	CreateExterns(options *CompilerOptions) []*SourceFile
	// Java: AbstractCommandLineRunner.shouldGenerateMapPerChunk(B)
	ShouldGenerateMapPerChunk(options *CompilerOptions) bool
	// Java: AbstractCommandLineRunner.expandSourceMapPath(B, JSChunk)
	ExpandSourceMapPath(options *CompilerOptions, forChunk *JSChunk) *string
	// Java: AbstractCommandLineRunner.filenameToOutputStream(String)
	FilenameToOutputStream(fileName *string) io.WriteCloser
	// Java: AbstractCommandLineRunner.createWriter(OutputStream, Charset)
	CreateWriter(stream io.WriteCloser, charset *jutil.Charset) AbstractCommandLineRunner_Writer
	// Java: AbstractCommandLineRunner.printChunkGraphJsonTo(Appendable)
	PrintChunkGraphJsonTo(out jutil.Appendable)
	// Java: AbstractCommandLineRunner.printChunkGraphManifestOrBundleTo(JSChunkGraph, Appendable, boolean)
	PrintChunkGraphManifestOrBundleTo(graph *JSChunkGraph, out jutil.Appendable, isManifest bool)
	// Java: AbstractCommandLineRunner.printManifestTo(JSChunk, Appendable)
	PrintManifestTo(chunk *JSChunk, out jutil.Appendable)
	// Java: AbstractCommandLineRunner.printBundleTo(JSChunk, Appendable)
	PrintBundleTo(chunk *JSChunk, out jutil.Appendable)
}

// AbstractCommandLineRunnerBase holds AbstractCommandLineRunner's fields and
// base implementations (PORTING.md §4).
type AbstractCommandLineRunnerBase struct {
	self AbstractCommandLineRunner

	config *AbstractCommandLineRunner_CommandLineConfig

	in io.Reader

	defaultJsOutput *jutil.PrintStream

	err *jutil.PrintStream

	compiler *Compiler

	inputCharset *jutil.Charset

	// NOTE(nicksantos): JSCompiler has always used ASCII as the default
	// output charset. This was done to solve legacy problems with
	// bad proxies, etc. We are not sure if these issues are still problems,
	// and changing the encoding would require a big obnoxious migration plan.
	//
	// New outputs should use outputCharset2, which is how I would have
	// designed this if I had a time machine.
	outputCharset2 *jutil.Charset

	legacyOutputCharset *jutil.Charset

	testMode bool

	externsSupplierForTesting func() []*SourceFile // nil is null

	inputsSupplierForTesting func() []*SourceFile // nil is null

	chunksSupplierForTesting func() []*JSChunk // nil is null

	exitCodeReceiver func(int)

	rootRelativePathsMap *jutil.LinkedHashMap[string, string] // nil is null

	parsedChunkWrappers *jutil.LinkedHashMap[string, string] // nil is null

	parsedChunkOutputFiles *jutil.LinkedHashMap[string, string] // nil is null

	parsedChunkConformanceFiles *jutil.LinkedHashMap[string, string] // nil is null

	gson *jutil.Gson

	filesToStreamOut []*AbstractCommandLineRunner_JsonFileSpec
}

// GetAbstractCommandLineRunnerBase gives access to the base fields through
// the interface (PORTING.md §4).
func (r *AbstractCommandLineRunnerBase) GetAbstractCommandLineRunnerBase() *AbstractCommandLineRunnerBase {
	return r
}

var (
	// Java: AbstractCommandLineRunner.OUTPUT_SAME_AS_INPUT_ERROR
	AbstractCommandLineRunner_OUTPUT_SAME_AS_INPUT_ERROR = DiagnosticType_Error(
		"JSC_OUTPUT_SAME_AS_INPUT_ERROR", "Bad output file (already listed as input file): {0}")

	// Java: AbstractCommandLineRunner.COULD_NOT_SERIALIZE_AST
	AbstractCommandLineRunner_COULD_NOT_SERIALIZE_AST = DiagnosticType_Error("JSC_COULD_NOT_SERIALIZE_AST", "Could not serialize ast to: {0}")

	// Java: AbstractCommandLineRunner.COULD_NOT_DESERIALIZE_AST
	AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST = DiagnosticType_Error("JSC_COULD_NOT_DESERIALIZE_AST", "Could not deserialize ast from: {0}")

	// Java: AbstractCommandLineRunner.NO_TREE_GENERATED_ERROR
	AbstractCommandLineRunner_NO_TREE_GENERATED_ERROR = DiagnosticType_Error(
		"JSC_NO_TREE_GENERATED_ERROR", "Code contains errors. No tree was generated.")
	// Java: AbstractCommandLineRunner.INVALID_CHUNK_SOURCEMAP_PATTERN
	AbstractCommandLineRunner_INVALID_CHUNK_SOURCEMAP_PATTERN = DiagnosticType_Error(
		"JSC_INVALID_CHUNK_SOURCEMAP_PATTERN",
		"When using --chunk or --module flags, the --create_source_map flag must contain "+
			"%outname% in the value.")

	// Java: AbstractCommandLineRunner.EXPECTED_DIAGNOSTIC_NOT_FOUND
	AbstractCommandLineRunner_EXPECTED_DIAGNOSTIC_NOT_FOUND = DiagnosticType_Error(
		"JSC_EXPECTED_DIAGNOSTIC_NOT_FOUND", "Expected diagnostic not found: {0}")

	// Java: AbstractCommandLineRunner.AMBIGUOUS_EXPECTATION
	AbstractCommandLineRunner_AMBIGUOUS_EXPECTATION = DiagnosticType_Error(
		"JSC_AMBIGUOUS_EXPECTATION",
		"Multiple expected diagnostics matched the error: \"{0}\". Matches: \"{1}\"")
)

// Java: AbstractCommandLineRunner.WAITING_FOR_INPUT_WARNING
const AbstractCommandLineRunner_WAITING_FOR_INPUT_WARNING = "The compiler is waiting for input via stdin."

// Use an 8MiB buffer since the concatenated TypedAst file can be very large.
//
// Java: AbstractCommandLineRunner.GZIPPED_TYPEDAST_BUFFER_SIZE
const abstractCommandLineRunner_GZIPPED_TYPEDAST_BUFFER_SIZE = 8 * 1024 * 1024

// Java: AbstractCommandLineRunner.OUTPUT_MARKER
const AbstractCommandLineRunner_OUTPUT_MARKER = "%output%"

// Java: AbstractCommandLineRunner.OUTPUT_MARKER_JS_STRING
const abstractCommandLineRunner_OUTPUT_MARKER_JS_STRING = "%output|jsstring%"

// InitAbstractCommandLineRunnerBase is the Go form of the constructor
// AbstractCommandLineRunner(InputStream, PrintStream, PrintStream), called
// by subclasses (PORTING.md §4). The Java no-argument and (PrintStream,
// PrintStream) constructors pass System.in, System.out and System.err for the
// missing streams; callers pass jutil.System_In(), jutil.System_Out() and
// jutil.System_Err().
//
// Java: AbstractCommandLineRunner(InputStream, PrintStream, PrintStream)
func (r *AbstractCommandLineRunnerBase) InitAbstractCommandLineRunnerBase(
	self AbstractCommandLineRunner, in io.Reader, out *jutil.PrintStream, err *jutil.PrintStream) {
	r.self = self
	r.exitCodeReceiver = AbstractCommandLineRunner_SystemExitCodeReceiver_INSTANCE.Apply
	r.config = NewAbstractCommandLineRunner_CommandLineConfig()
	r.in = jutil.CheckNotNull(in)
	r.defaultJsOutput = jutil.CheckNotNull(out)
	r.err = jutil.CheckNotNull(err)
	r.gson = jutil.NewGson()
}

// EnableTestMode puts the command line runner into test mode. In test mode,
// all outputs will be blackholed.
//
// externsSupplier is a provider for externs; inputsSupplier a provider for
// source inputs; chunksSupplier a provider for chunks. Only one of
// inputsSupplier and chunksSupplier may be non-null. exitCodeReceiver is a
// receiver for the status code that would have been passed to System.exit
// in non-test mode.
//
// Java: AbstractCommandLineRunner.enableTestMode(Supplier<List<SourceFile>>, Supplier<List<SourceFile>>, Supplier<List<JSChunk>>, Function<Integer, Void>)
func (r *AbstractCommandLineRunnerBase) EnableTestMode(
	externsSupplier func() []*SourceFile,
	inputsSupplier func() []*SourceFile,
	chunksSupplier func() []*JSChunk,
	exitCodeReceiver func(int)) {
	jutil.CheckArgument((inputsSupplier == nil) != (chunksSupplier == nil))
	r.testMode = true
	r.externsSupplierForTesting = externsSupplier
	r.inputsSupplierForTesting = inputsSupplier
	r.chunksSupplierForTesting = chunksSupplier
	r.exitCodeReceiver = exitCodeReceiver
}

// SetExitCodeReceiver sets the receiver of a non-zero integer to indicate a
// problem during execution or 0i to indicate success.
//
// Java: AbstractCommandLineRunner.setExitCodeReceiver(Function<Integer, Void>)
func (r *AbstractCommandLineRunnerBase) SetExitCodeReceiver(newExitCodeReceiver func(int)) {
	r.exitCodeReceiver = jutil.CheckNotNull(newExitCodeReceiver)
}

// IsInTestMode returns whether we're in test mode.
//
// Java: AbstractCommandLineRunner.isInTestMode()
func (r *AbstractCommandLineRunnerBase) IsInTestMode() bool {
	return r.testMode
}

// IsOutputInJson returns whether output should be a JSON stream.
//
// Java: AbstractCommandLineRunner.isOutputInJson()
func (r *AbstractCommandLineRunnerBase) IsOutputInJson() bool {
	return r.config.jsonStreamMode == CompilerOptions_JsonStreamMode_OUT ||
		r.config.jsonStreamMode == CompilerOptions_JsonStreamMode_BOTH
}

// GetCommandLineConfig gets the command line config, so that it can be
// initialized.
//
// Java: AbstractCommandLineRunner.getCommandLineConfig()
func (r *AbstractCommandLineRunnerBase) GetCommandLineConfig() *AbstractCommandLineRunner_CommandLineConfig {
	return r.config
}

// GetDiagnosticGroups returns the warning classes that are available from
// the command-line.
//
// Java: AbstractCommandLineRunner.getDiagnosticGroups()
func (r *AbstractCommandLineRunnerBase) GetDiagnosticGroups() *DiagnosticGroups {
	if r.compiler == nil {
		return NewDiagnosticGroups()
	}
	return r.compiler.GetDiagnosticGroups()
}

// Java: AbstractCommandLineRunner.setWarningGuardOptions(CompilerOptions, ArrayList<FlagEntry<CheckLevel>>, DiagnosticGroups)
func AbstractCommandLineRunner_SetWarningGuardOptions(
	options *CompilerOptions,
	warningGuards []*AbstractCommandLineRunner_FlagEntry[CheckLevel],
	diagnosticGroups *DiagnosticGroups) {
	if warningGuards != nil {
		groupNames := DiagnosticGroups_GetRegisteredGroups() // Java: .keySet()
		for _, entry := range warningGuards {
			if "*" == entry.value {
				for groupName := range groupNames.Keys() {
					if !DiagnosticGroups_wildcardExcludedGroups.Contains(groupName) {
						diagnosticGroups.SetWarningLevel(options, groupName, entry.flag)
					}
				}
			} else if groupNames.ContainsKey(entry.value) {
				diagnosticGroups.SetWarningLevel(options, entry.value, entry.flag)
			} else {
				panic(NewFlagUsageException("Unknown diagnostic group: '" + entry.value + "'"))
			}
		}
	}
}

// SetRunOptions sets options based on the configurations set flags API.
// Called during the run() run() method. If you want to ignore the flags
// API, or interpret flags your own way, then you should override this
// method.
//
// Java: AbstractCommandLineRunner.setRunOptions(CompilerOptions)
func (r *AbstractCommandLineRunnerBase) SetRunOptions(options *CompilerOptions) {
	config := r.config
	diagnosticGroups := r.self.GetDiagnosticGroups()

	if config.ShouldSaveAfterStage1() || config.ShouldContinueCompilation() {
		if options.IsChecksOnly() {
			panic(NewFlagUsageException(
				"checks_only mode is incompatible with multi-stage compilation"))
		}
		if options.GetExternExportsPath() != nil {
			panic(NewFlagUsageException(
				"generating externs from exports is incompatible with multi-stage compilation"))
		}
	}

	AbstractCommandLineRunner_SetWarningGuardOptions(options, config.warningGuards, diagnosticGroups)

	if config.warningsAllowFile != "" {
		r.self.AddAllowlistWarningsGuard(options, jutil.NewFile(config.warningsAllowFile))
	}

	if len(config.hideWarningsFor) != 0 {
		options.AddWarningsGuard(
			NewShowByPathWarningsGuardWithPathsAndShowType(
				append([]string{}, config.hideWarningsFor...),
				ShowByPathWarningsGuard_ShowType_EXCLUDE))
	}

	define := append([]string{}, config.define...)
	if config.browserFeaturesetYear != 0 {
		jutil.Try(func() {
			options.SetBrowserFeaturesetYear(config.browserFeaturesetYear)
		}, func(e jutil.IIllegalStateException) {
			panic(NewFlagUsageException(e.GetMessage()))
		})
	}

	AbstractCommandLineRunner_CreateDefineReplacements(define, options)

	options.SetTweakProcessing(config.tweakProcessing)

	// TODO(tjgq): Unconditionally set the options.
	if config.dependencyOptions != nil {
		options.SetDependencyOptions(config.dependencyOptions)
	}

	options.SetDevMode(config.jscompDevMode)
	options.SetCodingConvention(config.codingConvention)
	options.SetSummaryDetailLevel(config.summaryDetailLevel)
	options.SetTrustedStrings(true)

	r.legacyOutputCharset = r.GetLegacyOutputCharset()
	options.SetOutputCharset(r.legacyOutputCharset)
	r.outputCharset2 = r.GetOutputCharset2()
	r.inputCharset = r.GetInputCharset()

	if jutil.StringLength(config.jsOutputFile) > 0 {
		if config.skipNormalOutputs {
			panic(NewFlagUsageException(
				"skip_normal_outputs and js_output_file cannot be used together."))
		}
	}

	if config.skipNormalOutputs && config.printAst {
		panic(NewFlagUsageException("skip_normal_outputs and print_ast cannot be used together."))
	}

	if config.skipNormalOutputs && config.printTree {
		panic(NewFlagUsageException("skip_normal_outputs and print_tree cannot be used together."))
	}

	if config.skipNormalOutputs && config.printTreeJson {
		panic(NewFlagUsageException(
			"skip_normal_outputs and print_tree_json_path cannot be used together."))
	}

	if jutil.StringLength(config.createSourceMap) > 0 {
		options.SetSourceMapOutputPath(&config.createSourceMap)
	} else if r.IsOutputInJson() {
		outname := "%outname%"
		options.SetSourceMapOutputPath(&outname)
	}
	options.SetSourceMapDetailLevel(config.sourceMapDetailLevel)
	options.SetSourceMapFormat(config.sourceMapFormat)
	options.SetSourceMapLocationMappings(config.sourceMapLocationMappings)
	options.SetParseInlineSourceMaps(config.parseInlineSourceMaps)
	options.SetApplyInputSourceMaps(config.applyInputSourceMaps)

	inputSourceMaps := jutil.NewImmutableMapBuilder[string, *SourceMapInput]()
	for key, value := range config.sourceMapInputFiles.All() {
		sourceMap := SourceFile_Builder_().WithKind(rhino.StaticSourceFile_SourceKind_NON_CODE).WithPath(value).Build()
		inputSourceMaps.Put(key, NewSourceMapInput(sourceMap))
	}
	options.SetInputSourceMaps(inputSourceMaps.BuildOrThrow())

	if config.variableMapInputFile != "" {
		options.SetInputVariableMap(VariableMap_Load(config.variableMapInputFile))
	}

	if config.propertyMapInputFile != "" {
		options.SetInputPropertyMap(VariableMap_Load(config.propertyMapInputFile))
	}

	if len(config.outputManifests) != 0 {
		uniqueNames := jutil.NewLinkedHashSet[string]()
		for _, filename := range config.outputManifests {
			if !uniqueNames.Add(filename) {
				panic(NewFlagUsageException(
					"output_manifest flags specify duplicate file names: " + filename))
			}
		}
	}

	if len(config.outputBundles) != 0 {
		uniqueNames := jutil.NewLinkedHashSet[string]()
		for _, filename := range config.outputBundles {
			if !uniqueNames.Add(filename) {
				panic(NewFlagUsageException(
					"output_bundle flags specify duplicate file names: " + filename))
			}
		}
	}

	options.SetProcessCommonJSModules(config.processCommonJSModules)
	options.SetModuleRoots(config.moduleRoots)
	options.SetAngularPass(config.angularPass)

	if config.jsonWarningsFile != "" {
		options.AddReportGenerator(
			NewJsonErrorReportGenerator(jutil.NewPrintStreamFromFile(jutil.NewFile(config.jsonWarningsFile)), r.compiler))
	}

	if config.errorFormat == AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_JSON {
		errorGenerator := NewJsonErrorReportGenerator(r.GetErrorPrintStream(), r.compiler)
		r.compiler.SetErrorManager(NewSortingErrorManager([]SortingErrorManager_ErrorReportGenerator{errorGenerator}))
	}
	if config.printTree || config.printTreeJson {
		options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	}
	if config.skipNormalOutputs {
		// If skipping normal outputs, it's unnecessary to do a full AST parse of each input file.
		// The regex parser may still run if ordering/pruning inputs.
		r.compiler.SetPreferRegexParser(true)
	}
}

// Java: AbstractCommandLineRunner.getCompiler()
func (r *AbstractCommandLineRunnerBase) GetCompiler() *Compiler {
	return r.compiler
}

// AbstractCommandLineRunner_GetBuiltinExterns returns a mutable list.
//
// Java: AbstractCommandLineRunner.getBuiltinExterns(CompilerOptions.Environment)
func AbstractCommandLineRunner_GetBuiltinExterns(env CompilerOptions_Environment) []*SourceFile {
	input := abstractCommandLineRunner_getExternsInput()
	zipReader, err := zip.NewReader(bytes.NewReader(input), int64(len(input)))
	if err != nil {
		panic(jutil.NewZipException(err.Error()))
	}
	envPrefix := jutil.Ascii_ToLowerCase(env.String()) + "/"
	mapFromExternsZip := jutil.NewLinkedHashMap[string, *SourceFile]()
	for _, entry := range zipReader.File {
		filename := entry.Name

		// Always load externs in the root folder.
		// If the non-core-JS externs are organized in subfolders, only load
		// the ones in a subfolder matching the specified environment. Strip the subfolder.
		if strings.Contains(filename, "/") {
			if !strings.HasPrefix(filename, envPrefix) {
				continue
			}
			filename = filename[len(envPrefix):] // remove envPrefix, including '/'
		}

		entryStream, err := entry.Open()
		if err != nil {
			panic(jutil.NewZipException(err.Error()))
		}
		mapFromExternsZip.Put(
			filename,
			// Give the files an odd prefix, so that they do not conflict
			// with the user's files.
			SourceFile_Builder_().
				WithPath("externs.zip//"+filename).
				WithContentFromStream(entryStream).
				Build())
		entryStream.Close()
	}
	return DefaultExterns_PrepareExterns(env, mapFromExternsZip)
}

// Java: AbstractCommandLineRunner.getExternsInput()
func abstractCommandLineRunner_getExternsInput() []byte {
	input := resources.ResourceLoader_GetResourceAsStream("com.google.javascript.jscomp.AbstractCommandLineRunner", "/externs.zip")
	if input == nil {
		// In some environments, the externs.zip is relative to this class.
		input = resources.ResourceLoader_GetResourceAsStream("com.google.javascript.jscomp.AbstractCommandLineRunner", "externs.zip")
	}
	jutil.CheckNotNull(input)
	return input
}

// Run runs the Compiler and calls System.exit() with the exit status of the
// compiler.
//
// Java: AbstractCommandLineRunner.run()
func (r *AbstractCommandLineRunnerBase) Run() {
	var result int
	jutil.Try(func() {
		jutil.Try(func() {
			result = r.self.DoRun()
		}, func(e IFlagUsageException) {
			r.err.Println(e.GetMessage())
			result = -1
		})
	}, func(t jutil.Throwable) {
		jutil.Throwable_PrintStackTrace(t, r.err)
		result = -2
	})

	r.exitCodeReceiver(result)
}

// GetErrorPrintStream returns the PrintStream for writing errors associated
// with this AbstractCommandLineRunner.
//
// Java: AbstractCommandLineRunner.getErrorPrintStream()
func (r *AbstractCommandLineRunnerBase) GetErrorPrintStream() *jutil.PrintStream {
	return r.err
}

// Java: AbstractCommandLineRunner.parseJsonFilesFromInputStream()
func (r *AbstractCommandLineRunnerBase) ParseJsonFilesFromInputStream() []*AbstractCommandLineRunner_JsonFileSpec {
	var jsonFiles []*AbstractCommandLineRunner_JsonFileSpec
	reader := jutil.NewJsonReader(jutil.NewInputStreamReader(r.in, r.inputCharset))
	defer reader.Close()
	reader.BeginArray()
	for reader.HasNext() {
		var jsonFile *AbstractCommandLineRunner_JsonFileSpec
		if v := r.gson.FromJsonReaderReflective(reader, func() any { return newAbstractCommandLineRunner_JsonFileSpecDefault() }); v != nil {
			jsonFile = v.(*AbstractCommandLineRunner_JsonFileSpec)
		}
		jsonFiles = append(jsonFiles, jsonFile)
	}
	reader.EndArray()
	return jsonFiles
}

// CreateInputs creates inputs from a list of files.
//
// files is a list of flag entries indicates js and zip file names;
// allowStdIn whether '-' is allowed appear as a filename to represent stdin
// (if true, '-' is only allowed to appear once); jsChunkSpecs a list chunk
// specs. It returns an array of inputs.
//
// Java: AbstractCommandLineRunner.createInputs(List<FlagEntry<JsSourceType>>, boolean, List<JsChunkSpec>)
func (r *AbstractCommandLineRunnerBase) CreateInputs(
	files []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType],
	allowStdIn bool,
	jsChunkSpecs []*AbstractCommandLineRunner_JsChunkSpec) []*SourceFile {
	return r.self.CreateInputsWithJsonFilesAndAllowStdIn(files /* jsonFiles= */, nil, allowStdIn, jsChunkSpecs)
}

// CreateInputsWithJsonFiles creates inputs from a list of source files and
// json files.
//
// Java: AbstractCommandLineRunner.createInputs(List<FlagEntry<JsSourceType>>, List<JsonFileSpec>, List<JsChunkSpec>)
func (r *AbstractCommandLineRunnerBase) CreateInputsWithJsonFiles(
	files []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType],
	jsonFiles []*AbstractCommandLineRunner_JsonFileSpec,
	jsChunkSpecs []*AbstractCommandLineRunner_JsChunkSpec) []*SourceFile {
	return r.self.CreateInputsWithJsonFilesAndAllowStdIn(files, jsonFiles /* allowStdIn= */, false, jsChunkSpecs)
}

// CreateInputsWithJsonFilesAndAllowStdIn creates inputs from a list of
// source files, zips and json files.
//
// Can be overridden by subclasses who want to pull files from different
// places.
//
// Go-only: language-adapter — a nil jsonFiles is Java's null list; an empty
// one is passed as a non-nil empty slice.
//
// Java: AbstractCommandLineRunner.createInputs(List<FlagEntry<JsSourceType>>, List<JsonFileSpec>, boolean, List<JsChunkSpec>)
func (r *AbstractCommandLineRunnerBase) CreateInputsWithJsonFilesAndAllowStdIn(
	files []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType],
	jsonFiles []*AbstractCommandLineRunner_JsonFileSpec,
	allowStdIn bool,
	jsChunkSpecs []*AbstractCommandLineRunner_JsChunkSpec) []*SourceFile {
	inputs := make([]*SourceFile, 0, len(files))
	usingStdin := false
	jsChunkIndex := 0
	var jsChunkSpec *AbstractCommandLineRunner_JsChunkSpec
	if len(jsChunkSpecs) > 0 {
		jsChunkSpec = jsChunkSpecs[0]
	}
	cumulatedInputFilesExpected := math.MaxInt32
	if jsChunkSpec != nil {
		cumulatedInputFilesExpected = jsChunkSpec.GetNumInputs()
	}
	for i := 0; i < len(files); i++ {
		file := files[i]
		filename := file.value
		if file.flag == AbstractCommandLineRunner_JsSourceType_JS_ZIP {
			if r.config.typedAstListInputFilename != nil {
				panic(NewFlagUsageException("Can't use TypedASTs with --zip."))
			}

			if "-" != filename {
				newFiles := SourceFile_FromZipFile(filename, r.inputCharset)
				// Update the manifest maps for new zip entries.
				if r.rootRelativePathsMap.ContainsKey(filename) {
					rootFilename := r.rootRelativePathsMap.Get(filename)
					for _, zipEntry := range newFiles {
						zipEntryName := zipEntry.GetName()
						jutil.CheckState(strings.Contains(zipEntryName, filename))
						zipmap := strings.ReplaceAll(zipEntryName, filename, rootFilename)
						r.rootRelativePathsMap.Put(zipEntryName, zipmap)
					}
				}

				inputs = append(inputs, newFiles...)
				if jsChunkSpec != nil {
					jsChunkSpec.numJsFiles += len(newFiles) - 1
				}
			}
		} else if "-" != filename {
			kind := rhino.StaticSourceFile_SourceKind_STRONG
			if file.flag == AbstractCommandLineRunner_JsSourceType_WEAKDEP {
				kind = rhino.StaticSourceFile_SourceKind_WEAK
			}
			newFile := SourceFile_Builder_().
				WithPath(filename).
				WithCharset(r.inputCharset).
				WithKind(kind).
				Build()
			inputs = append(inputs, newFile)
		} else {
			if r.config.typedAstListInputFilename != nil {
				panic(NewFlagUsageException("Can't use TypedASTs with stdin."))
			}
			if !r.config.defaultToStdin {
				panic(NewFlagUsageException("Can't specify stdin."))
			}
			if usingStdin {
				panic(NewFlagUsageException("Can't specify stdin twice."))
			}

			if len(r.config.outputManifests) != 0 {
				panic(NewFlagUsageException(
					"Manifest files cannot be generated when the input is from stdin."))
			}
			if len(r.config.outputBundles) != 0 {
				panic(NewFlagUsageException(
					"Bundle files cannot be generated when the input is from stdin."))
			}

			r.err.Println(AbstractCommandLineRunner_WAITING_FOR_INPUT_WARNING)
			inputs = append(inputs,
				SourceFile_Builder_().
					WithPath("stdin").
					WithContentFromStream(r.in).
					WithCharset(r.inputCharset).
					Build())
			usingStdin = true
		}
		if i >= cumulatedInputFilesExpected-1 {
			jsChunkIndex++
			if jsChunkIndex < len(jsChunkSpecs) {
				jsChunkSpec = jsChunkSpecs[jsChunkIndex]
				cumulatedInputFilesExpected += jsChunkSpec.GetNumInputs()
			}
		}
	}
	if jsonFiles != nil {
		for _, jsonFile := range jsonFiles {
			inputs = append(inputs, SourceFile_FromCode(*jsonFile.GetPath(), *jsonFile.GetSrc()))
		}
	}

	return inputs
}

// DeduplicateIjsFiles removes input --ijs files whose basename matches an
// input --js or --weakdep file. It throws FlagUsageException if there are
// both input --ijs files and chunk specs.
//
// Go-only: language-adapter — the Java method removes entries from the
// caller's list through its iterator; files is a pointer to the caller's
// slice.
//
// Java: AbstractCommandLineRunner.deduplicateIjsFiles(List<FlagEntry<JsSourceType>>, List<String>, boolean)
func (r *AbstractCommandLineRunnerBase) DeduplicateIjsFiles(
	files *[]*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType],
	moduleRoots []string,
	hasChunkSpecs bool) []*JSError {
	var errors []*JSError

	// First pass: collect the (module root relative) names of --js and --weakdep files.
	relativeToAbsoluteName := jutil.NewLinkedHashMap[string, string]()
	for _, file := range *files {
		// TODO(tjgq): Handle zip files.
		if file.flag == AbstractCommandLineRunner_JsSourceType_JS || file.flag == AbstractCommandLineRunner_JsSourceType_WEAKDEP {
			absoluteName := file.value
			relativeName := r.GetModuleRootRelativeName(absoluteName, moduleRoots)
			relativeToAbsoluteName.Put(relativeName, absoluteName)
		}
	}

	// Second pass: drop --ijs files whose (module root relative) name matches a --js or --weakdep
	// file.
	// Go-only: language-adapter — Iterator.remove() as an index loop.
	for i := 0; i < len(*files); {
		file := (*files)[i]
		if file.flag == AbstractCommandLineRunner_JsSourceType_IJS {
			if hasChunkSpecs {
				panic(NewFlagUsageException("--ijs is incompatible with --chunk or --module."))
			}
			absoluteName := file.value
			if !strings.HasSuffix(absoluteName, ".i.js") {
				errors = append(errors, JSError_Make(IjsErrors_BAD_IJS_FILE_NAME, absoluteName))
				i++
				continue
			}
			relativeName := r.GetModuleRootRelativeName(absoluteName, moduleRoots)
			relativeNonIjsName := relativeName[:len(relativeName)-len(".i.js")]
			if relativeToAbsoluteName.ContainsKey(relativeNonIjsName) {
				errors = append(errors,
					JSError_Make(
						IjsErrors_CONFLICTING_IJS_FILE,
						relativeToAbsoluteName.Get(relativeNonIjsName),
						absoluteName))
				*files = append((*files)[:i:i], (*files)[i+1:]...)
				continue
			}
		}
		i++
	}

	return errors
}

// Java: AbstractCommandLineRunner.getModuleRootRelativeName(String, List<String>)
func (r *AbstractCommandLineRunnerBase) GetModuleRootRelativeName(filename string, moduleRoots []string) string {
	for _, moduleRoot := range moduleRoots {
		if strings.HasPrefix(filename, moduleRoot+"/") {
			return filename[len(moduleRoot)+1:]
		}
	}
	return filename
}

// CreateSourceInputs creates JS source code inputs from a list of files.
//
// Java: AbstractCommandLineRunner.createSourceInputs(List<JsChunkSpec>, List<FlagEntry<JsSourceType>>, List<JsonFileSpec>, List<String>)
func (r *AbstractCommandLineRunnerBase) CreateSourceInputs(
	jsChunkSpecs []*AbstractCommandLineRunner_JsChunkSpec,
	files []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType],
	jsonFiles []*AbstractCommandLineRunner_JsonFileSpec,
	moduleRoots []string) []*SourceFile {
	if r.self.IsInTestMode() {
		if r.inputsSupplierForTesting != nil {
			return r.inputsSupplierForTesting()
		}
		return nil
	}
	if len(files) == 0 && jsonFiles == nil && r.config.defaultToStdin {
		// Request to read from stdin.
		files = []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]{
			NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS, "-")}
	}

	for _, error := range r.DeduplicateIjsFiles(&files, moduleRoots, len(jsChunkSpecs) != 0) {
		r.compiler.Report(error)
	}

	var result []*SourceFile
	jutil.Try(func() {
		if jsonFiles != nil {
			result = r.CreateInputsWithJsonFiles(files, jsonFiles, jsChunkSpecs)
		} else {
			result = r.CreateInputs(files, true, jsChunkSpecs)
		}
	}, func(e IFlagUsageException) {
		panic(NewFlagUsageException("Bad --js flag. " + e.GetMessage()))
	})
	return result
}

// CreateExternInputs creates JS extern inputs from a list of files.
//
// Java: AbstractCommandLineRunner.createExternInputs(List<String>)
func (r *AbstractCommandLineRunnerBase) CreateExternInputs(files []string) []*SourceFile {
	var externFiles []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]
	for _, file := range files {
		externFiles = append(externFiles, NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_EXTERN, file))
	}
	var result []*SourceFile
	jutil.Try(func() {
		result = r.CreateInputs(externFiles, false, []*AbstractCommandLineRunner_JsChunkSpec{})
	}, func(e IFlagUsageException) {
		panic(NewFlagUsageException("Bad --externs flag. " + e.GetMessage()))
	})
	return result
}

// AbstractCommandLineRunner_CreateJsChunks creates chunk objects from a
// list of chunk specifications.
//
// specs is a list of chunk specifications, not null or empty; inputs a list
// of JS file paths, not null. It returns an array of chunk objects.
//
// Java: AbstractCommandLineRunner.createJsChunks(List<JsChunkSpec>, List<CompilerInput>)
func AbstractCommandLineRunner_CreateJsChunks(specs []*AbstractCommandLineRunner_JsChunkSpec, inputs []*CompilerInput) []*JSChunk {
	jutil.CheckState(specs != nil)
	jutil.CheckState(len(specs) != 0)
	jutil.CheckState(inputs != nil)

	chunkNames := make([]string, 0, len(specs))
	chunksByName := jutil.NewLinkedHashMap[string, *JSChunk]()
	chunksFileCountMap := jutil.NewLinkedHashMap[string, int]()
	numJsFilesExpected := 0
	minJsFilesRequired := 0
	for _, spec := range specs {
		if chunksByName.ContainsKey(spec.name) {
			panic(NewFlagUsageException("Duplicate chunk name: " + spec.name))
		}
		chunk := NewJSChunk(spec.name)

		for _, dep := range spec.deps {
			other := chunksByName.Get(dep)
			if other == nil {
				panic(NewFlagUsageException(
					"Chunk '" +
						spec.name +
						"' depends on unknown chunk '" +
						dep +
						"'. Be sure to list chunks in dependency order."))
			}
			chunk.AddDependency(other)
		}

		// We will allow chunks of zero input.
		if spec.numJsFiles < 0 {
			numJsFilesExpected = -1
		} else {
			minJsFilesRequired += spec.numJsFiles
		}

		if numJsFilesExpected >= 0 {
			numJsFilesExpected += spec.numJsFiles
		}

		// Add chunks in reverse order so that source files are allocated to
		// chunks in reverse order. This allows the first chunk
		// (presumably the base chunk) to have a size of 'auto'
		chunkNames = append([]string{spec.name}, chunkNames...)
		chunksFileCountMap.Put(spec.name, spec.numJsFiles)
		chunksByName.Put(spec.name, chunk)
	}

	totalNumJsFiles := len(inputs)

	if numJsFilesExpected >= 0 || minJsFilesRequired > totalNumJsFiles {
		if minJsFilesRequired > totalNumJsFiles {
			numJsFilesExpected = minJsFilesRequired
		}

		if numJsFilesExpected > totalNumJsFiles {
			panic(NewFlagUsageException(
				"Not enough JS files specified. Expected " +
					strconv.Itoa(numJsFilesExpected) +
					" but found " +
					strconv.Itoa(totalNumJsFiles)))
		} else if numJsFilesExpected < totalNumJsFiles {
			panic(NewFlagUsageException(
				"Too many JS files specified. Expected " +
					strconv.Itoa(numJsFilesExpected) +
					" but found " +
					strconv.Itoa(totalNumJsFiles)))
		}
	}

	numJsFilesLeft := totalNumJsFiles
	chunkIndex := 0
	for _, chunkName := range chunkNames {
		// Parse chunk inputs.
		numJsFiles := chunksFileCountMap.Get(chunkName)
		chunk := chunksByName.Get(chunkName)

		// Check if the first chunk specified 'auto' for the number of files
		if chunkIndex == len(chunkNames)-1 && numJsFiles == -1 {
			numJsFiles = numJsFilesLeft
		}

		chunkFiles := inputs[numJsFilesLeft-numJsFiles : numJsFilesLeft]
		for _, input := range chunkFiles {
			chunk.Add(input)
		}
		numJsFilesLeft -= numJsFiles
		chunkIndex++
	}

	var result []*JSChunk
	for chunk := range chunksByName.Values() {
		result = append(result, chunk)
	}
	return result
}

// CheckChunkName validates the chunk name. Can be overridden by subclasses.
//
// Java: AbstractCommandLineRunner.checkChunkName(String)
func (r *AbstractCommandLineRunnerBase) CheckChunkName(name string) {
	if !rhino.TokenStream_IsJSIdentifier(name) {
		panic(NewFlagUsageException("Invalid chunk name: '" + name + "'"))
	}
}

// AbstractCommandLineRunner_ParseChunkWrappers parses chunk wrapper
// specifications.
//
// specs is a list of chunk wrapper specifications, not null. The spec
// format is: name:wrapper. Wrappers. chunks are the JS chunks whose
// wrappers are specified. It returns a map from chunk name to chunk wrapper.
// Chunks with no wrapper will have the empty string as their value in this
// map.
//
// Java: AbstractCommandLineRunner.parseChunkWrappers(List<String>, Iterable<JSChunk>)
func AbstractCommandLineRunner_ParseChunkWrappers(specs []string, chunks []*JSChunk) *jutil.LinkedHashMap[string, string] {
	jutil.CheckState(specs != nil)

	wrappers := jutil.NewLinkedHashMap[string, string]()

	// Prepopulate the map with chunk names.
	for _, c := range chunks {
		wrappers.Put(c.GetName(), "")
	}

	for _, spec := range specs {
		// Format is "<name>:<wrapper>".
		pos := jutil.StringIndexOfChar(spec, ':')
		if pos == -1 {
			panic(NewFlagUsageException(
				"Expected chunk wrapper to have " + "<name>:<wrapper> format: " + spec))
		}

		// Parse chunk name.
		name := jutil.StringSubstring(spec, 0, pos)
		if !wrappers.ContainsKey(name) {
			panic(NewFlagUsageException("Unknown chunk: '" + name + "'"))
		}
		wrapper := jutil.StringSubstringFrom(spec, pos+1)
		// Support for %n% and %output%
		wrapper = strings.ReplaceAll(strings.ReplaceAll(wrapper, "%output%", "%s"), "%n%", "\n")
		if !strings.Contains(wrapper, "%s") {
			panic(NewFlagUsageException("No %s placeholder in chunk wrapper: '" + wrapper + "'"))
		}

		wrappers.Put(name, wrapper)
	}
	return wrappers
}

// abstractCommandLineRunner_parseChunkOutputFiles parses chunk output name
// specifications.
//
// specs is a list of chunk output name specifications, not null. The spec
// format is: name:output_file. It returns a map from chunk name to chunk
// output file name, if declared. Chunks with no predeclared output file name
// will have no entry in this map.
//
// Java: AbstractCommandLineRunner.parseChunkOutputFiles(List<String>)
func abstractCommandLineRunner_parseChunkOutputFiles(specs []string) *jutil.LinkedHashMap[string, string] {
	jutil.CheckArgument(specs != nil)

	outputFilesBuilder := jutil.NewImmutableMapBuilder[string, string]()

	for _, spec := range specs {
		// Format is "<name>:<output_file>".
		pos := jutil.StringIndexOfChar(spec, ':')
		if pos == -1 {
			panic(NewFlagUsageException(
				"Expected chunk_output_file to have " + "<name>:<output_file> format: " + spec))
		}

		name := jutil.StringSubstring(spec, 0, pos)
		filename := jutil.StringSubstringFrom(spec, pos+1)
		outputFilesBuilder.Put(name, filename)
	}
	return outputFilesBuilder.BuildOrThrow()
}

// GetChunkOutputFileName returns the output file name for a chunk.
//
// For chunks with predeclared output file names specified using
// --chunk_output_file, the output file name is
// <outputPathPrefix>/<output_file>.
//
// Otherwise, the output file name is <outputPathPrefix>/<chunkName>.js
//
// Java: AbstractCommandLineRunner.getChunkOutputFileName(JSChunk)
func (r *AbstractCommandLineRunnerBase) GetChunkOutputFileName(m *JSChunk) string {
	if r.parsedChunkOutputFiles == nil {
		r.parsedChunkOutputFiles = abstractCommandLineRunner_parseChunkOutputFiles(r.config.chunkOutputFiles)
	}
	if outputFile, ok := r.parsedChunkOutputFiles.Lookup(m.GetName()); ok {
		return r.config.chunkOutputPathPrefix + outputFile
	}
	return r.config.chunkOutputPathPrefix + m.GetName() + ".js"
}

// GetChunkConformanceFileName returns the conformance file name for a
// chunk.
//
// Java: AbstractCommandLineRunner.getChunkConformanceFileName(JSChunk)
func (r *AbstractCommandLineRunnerBase) GetChunkConformanceFileName(m *JSChunk) *string {
	if r.parsedChunkConformanceFiles == nil {
		r.parsedChunkConformanceFiles = abstractCommandLineRunner_parseChunkOutputFiles(r.config.chunkConformanceFiles)
	}
	if name, ok := r.parsedChunkConformanceFiles.Lookup(m.GetName()); ok {
		return &name
	}
	return nil
}

// Java: AbstractCommandLineRunner.writeChunkOutput(String, Appendable, LicenseTracker, JSChunk)
func (r *AbstractCommandLineRunnerBase) WriteChunkOutput(fileName string, out jutil.Appendable, lt CodePrinter_LicenseTracker, m *JSChunk) {
	if r.parsedChunkWrappers == nil {
		r.parsedChunkWrappers = AbstractCommandLineRunner_ParseChunkWrappers(
			r.config.chunkWrapper, append([]*JSChunk{}, r.compiler.GetChunkGraph().GetAllChunks()...))
	}

	if !r.IsOutputInJson() {
		abstractCommandLineRunner_maybeCreateDirsForPath(&fileName)
	}
	baseName := jutil.NewFile(fileName).GetName()
	r.self.WriteOutput(
		out,
		r.compiler,
		lt,
		m,
		strings.ReplaceAll(r.parsedChunkWrappers.Get(m.GetName()), "%basename%", baseName),
		"%s",
		nil,
		fileName)
}

// WriteOutput writes code to an output stream, optionally wrapping it in an
// arbitrary wrapper that contains a placeholder where the code should be
// inserted.
//
// chunk is which chunk to write. If this is null, write the entire AST.
//
// Java: AbstractCommandLineRunner.writeOutput(Appendable, Compiler, LicenseTracker, JSChunk, String, String, Function<String, String>, String)
func (r *AbstractCommandLineRunnerBase) WriteOutput(
	out jutil.Appendable,
	compiler *Compiler,
	licenseTracker CodePrinter_LicenseTracker,
	chunk *JSChunk,
	wrapper string,
	codePlaceholder string,
	escaper func(string) string,
	filename string) {
	if compiler.GetOptions().GetOutputJs() == CompilerOptions_OutputJs_SENTINEL {
		out.AppendCharSequence("// No JS output because the compiler was run in checks-only mode.\n")
		return
	}
	jutil.CheckState(compiler.GetOptions().GetOutputJs() == CompilerOptions_OutputJs_NORMAL)

	var code string
	if chunk == nil {
		code = compiler.ToSource()
	} else {
		code = compiler.ToSourceOfLicenseTrackerAndJSChunk(licenseTracker, chunk)
	}
	r.self.WriteOutputOfCode(out, compiler, code, wrapper, codePlaceholder, escaper, filename)
}

// WriteOutputOfCode writes code to an output stream, optionally wrapping it
// in an arbitrary wrapper that contains a placeholder where the code should
// be inserted.
//
// Java: AbstractCommandLineRunner.writeOutput(Appendable, Compiler, String, String, String, Function<String, String>, String)
func (r *AbstractCommandLineRunnerBase) WriteOutputOfCode(
	out jutil.Appendable,
	compiler *Compiler,
	code string,
	wrapper string,
	codePlaceholder string,
	escaper func(string) string,
	filename string) {
	pos := jutil.StringIndexOf(wrapper, codePlaceholder)
	if pos != -1 {
		prefix := ""

		if pos > 0 {
			prefix = jutil.StringSubstring(wrapper, 0, pos)
			out.AppendCharSequence(prefix)
		}

		if escaper == nil {
			out.AppendCharSequence(code)
		} else {
			out.AppendCharSequence(escaper(code))
		}

		suffixStart := pos + jutil.StringLength(codePlaceholder)
		if suffixStart != jutil.StringLength(wrapper) {
			// Something after placeholder?
			out.AppendCharSequence(jutil.StringSubstringFrom(wrapper, suffixStart))
		}
		if r.self.GetCommandLineConfig().includeTrailingNewline {
			out.AppendCharacter('\n')
		}

		// If we have a source map, adjust its offsets to match
		// the code WITHIN the wrapper.
		if compiler != nil && compiler.GetSourceMap() != nil {
			compiler.GetSourceMap().SetWrapperPrefix(prefix)
		}
	} else {
		out.AppendCharSequence(code)
		if r.self.GetCommandLineConfig().includeTrailingNewline {
			out.AppendCharacter('\n')
		}
	}
}

// abstractCommandLineRunner_maybeCreateDirsForPath creates any directories
// necessary to write a file that will have a given path prefix.
//
// Java: AbstractCommandLineRunner.maybeCreateDirsForPath(String)
func abstractCommandLineRunner_maybeCreateDirsForPath(pathPrefix *string) {
	if !jutil.Strings_IsNullOrEmpty(pathPrefix) {
		parent := jutil.NewFile(*pathPrefix).GetParentFile()
		if parent != nil {
			parent.Mkdirs()
		}
	}
}

// Java: AbstractCommandLineRunner.createDefaultOutput()
func (r *AbstractCommandLineRunnerBase) CreateDefaultOutput() AbstractCommandLineRunner_Writer {
	writeOutputToFile := r.config.jsOutputFile != ""
	if writeOutputToFile {
		return r.FileNameToLegacyOutputWriter(&r.config.jsOutputFile)
	}
	return r.StreamToLegacyOutputWriter(abstractCommandLineRunner_printStreamAsOutputStream{r.defaultJsOutput})
}

// Java: AbstractCommandLineRunner.closeAppendable(Appendable)
func abstractCommandLineRunner_closeAppendable(output jutil.Appendable) {
	if flushable, ok := output.(interface{ Flush() }); ok {
		flushable.Flush()
	}
	if closeable, ok := output.(interface{ Close() }); ok {
		closeable.Close()
	}
}

// DoRun parses command-line arguments and runs the compiler. It returns the
// system exit status.
//
// Java: AbstractCommandLineRunner.doRun()
func (r *AbstractCommandLineRunnerBase) DoRun() int {
	config := r.config
	metricsRecorder := r.self.GetCompileMetricsRecorder()
	metricsRecorder.RecordActionStart()
	if config.printVersion {
		r.defaultJsOutput.Println(r.self.GetVersionText())
		r.defaultJsOutput.Flush()
		return 0
	}

	Compiler_SetLoggingLevel(jutil.Level_Parse(config.loggingLevel))

	r.compiler = r.self.CreateCompiler()
	if len(config.expectedDiagnostics) != 0 {
		r.compiler.SetErrorManager(
			NewVerifyingErrorManager(r.compiler.GetErrorManager(), config.expectedDiagnostics))
	}
	options := r.self.CreateOptions()
	r.self.SetRunOptions(options)

	typedAstListInputFilename := config.typedAstListInputFilename

	externs := r.self.CreateExterns(options)
	var chunks []*JSChunk

	r.rootRelativePathsMap = r.ConstructRootRelativePathsMap()

	writeOutputToFile := config.jsOutputFile != ""
	var outputFileNames []string
	if writeOutputToFile {
		outputFileNames = append(outputFileNames, config.jsOutputFile)
	}

	createCommonJsModules := false
	if options.GetProcessCommonJSModules() &&
		(len(config.chunk) == 1 && config.chunk[0] == "auto") {
		createCommonJsModules = true
		config.chunk = config.chunk[1:]
	}

	var jsChunkSpecs []*AbstractCommandLineRunner_JsChunkSpec
	for i := 0; i < len(config.chunk); i++ {
		jsChunkSpecs = append(jsChunkSpecs, AbstractCommandLineRunner_JsChunkSpec_Create(config.chunk[i], i == 0))
	}
	var jsonFiles []*AbstractCommandLineRunner_JsonFileSpec

	if config.jsonStreamMode == CompilerOptions_JsonStreamMode_IN ||
		config.jsonStreamMode == CompilerOptions_JsonStreamMode_BOTH {
		jsonFiles = r.self.ParseJsonFilesFromInputStream()
		if jsonFiles == nil {
			jsonFiles = []*AbstractCommandLineRunner_JsonFileSpec{}
		}

		inputSourceMaps := jutil.NewImmutableMapBuilder[string, *SourceMapInput]()
		inputPathByWebpackId := jutil.NewImmutableMapBuilder[string, string]()

		foundJsonInputSourceMap := false
		for _, jsonFile := range jsonFiles {
			if jsonFile.GetSourceMap() != nil && jutil.StringLength(*jsonFile.GetSourceMap()) > 0 {
				sourceMapPath := jutil.ToJavaString(jsonFile.GetPath()) + ".map"
				sourceMap := SourceFile_FromCode(sourceMapPath, *jsonFile.GetSourceMap())
				inputSourceMaps.Put(*jsonFile.GetPath(), NewSourceMapInput(sourceMap))
				foundJsonInputSourceMap = true
			}
			if jsonFile.GetWebpackId() != nil {
				inputPathByWebpackId.Put(*jsonFile.GetWebpackId(), *jsonFile.GetPath())
			}
		}

		if foundJsonInputSourceMap {
			inputSourceMaps.PutAll(options.GetInputSourceMaps())
			options.SetInputSourceMaps(inputSourceMaps.BuildOrThrow())
		}

		r.compiler.InitWebpackMap(inputPathByWebpackId.BuildOrThrow())
	} else {
		emptyMap := jutil.NewLinkedHashMap[string, string]()
		r.compiler.InitWebpackMap(emptyMap)
	}

	options.SetDoLateLocalization(config.ShouldDoLateLocalization())
	options.SetAlwaysGatherSourceMapInfo(config.ShouldAlwaysGatherSourceMapInfo())

	r.compiler.InitOptions(options)

	sources := r.CreateSourceInputs(jsChunkSpecs, config.mixedJsSources, jsonFiles, config.moduleRoots)

	if len(jsChunkSpecs) != 0 {
		if r.self.IsInTestMode() {
			chunks = r.chunksSupplierForTesting()
		} else {
			if CompilerOptions_JsonStreamMode_IN == config.jsonStreamMode ||
				CompilerOptions_JsonStreamMode_NONE == config.jsonStreamMode {
				for _, m := range jsChunkSpecs {
					r.self.CheckChunkName(m.GetName())
				}
			}
			inputs := make([]*CompilerInput, 0, len(sources))
			for _, source := range sources {
				inputs = append(inputs, NewCompilerInputWithIsExtern(source /* isExtern= */, false))
			}
			chunks = AbstractCommandLineRunner_CreateJsChunks(jsChunkSpecs, inputs)
		}
		for _, m := range chunks {
			outputFileNames = append(outputFileNames, r.self.GetChunkOutputFileName(m))
		}

		if typedAstListInputFilename != nil {
			r.InitChunksWithTypedAstFilesystem(externs, chunks, options, *typedAstListInputFilename)
		} else {
			r.compiler.InitChunks(externs, chunks, options)
		}
	} else {
		if typedAstListInputFilename != nil {
			r.InitWithTypedAstFilesystem(externs, sources, options, *typedAstListInputFilename)
		} else {
			r.compiler.Init(externs, sources, options)
		}
	}

	// Release temporary data structures now that the compiler has
	// been initialized
	jsChunkSpecs = nil
	jsonFiles = nil
	externs = nil
	sources = nil

	var result *Result
	if config.skipNormalOutputs {
		metricsRecorder.RecordActionName("skip normal outputs")
		// TODO(bradfordcsmith): Should we be ignoring possible init/initChunks() errors here?
		r.compiler.OrderInputsWithLargeStack()
		result = nil
	} else if r.compiler.HasErrors() {
		metricsRecorder.RecordActionName("initialization errors occurred")
		// init() or initChunks() encountered an error.
		r.compiler.GenerateReport()
		result = r.compiler.GetResult()
	} else {
		func() {
			defer func() {
				// Make sure we generate a report of errors and warnings even if the compiler throws an
				// exception somewhere.
				r.compiler.GenerateReport()
			}()
			// This is the common case - we're actually compiling something.
			r.PerformCompilation(metricsRecorder)

			// If we're finished with compilation (i.e. we're not saving state), /and/ the compiler was
			// restored from a previous state, then we need to re-initialize the set of chunks.
			// TODO(lharker): figure out if this is still needed.
			refresh := chunks != nil &&
				config.restoredCompilationStage != -1 &&
				config.saveAfterCompilationStage == -1
			if refresh {
				chunks = append([]*JSChunk{}, r.compiler.GetChunks()...)
			}
		}()
		result = r.compiler.GetResult()
	}

	if createCommonJsModules {
		// For CommonJS modules, construct chunks from actual inputs.
		chunks = append([]*JSChunk{}, r.compiler.GetChunks()...)
		for _, c := range chunks {
			outputFileNames = append(outputFileNames, r.self.GetChunkOutputFileName(c))
		}
	}

	for _, outputFileName := range outputFileNames {
		if r.compiler.GetSourceFileByName(outputFileName) != nil {
			r.compiler.Report(JSError_Make(AbstractCommandLineRunner_OUTPUT_SAME_AS_INPUT_ERROR, outputFileName))
			return 1
		}
	}

	// We won't want to process results for cases where compilation is only partially done.
	shouldProcessResults := config.GetSaveCompilationStateToFilename() == nil
	var exitStatus int
	if shouldProcessResults {
		exitStatus = r.self.ProcessResults(result, chunks, options)
	} else {
		exitStatus = abstractCommandLineRunner_getExitStatusForResult(result)
	}
	metricsRecorder.RecordResultMetrics(r.compiler, result)
	return exitStatus
}

// Java: AbstractCommandLineRunner.performCompilation(CompileMetricsRecorderInterface)
func (r *AbstractCommandLineRunnerBase) PerformCompilation(metricsRecorder CompileMetricsRecorderInterface) {
	// Parse, restore from a save file, or initialize from a TypedAST list.
	r.InitializeStateBeforeCompilation()

	// Do the interesting work - checks, optimizations, etc.
	r.RunCompilerPasses(metricsRecorder)

	// Save state unless there was an exception or error - we won't go on to the next stage if there
	// was an error, and saving state itself might throw an exception if the compiler is in a bad
	// state.
	if !r.compiler.HasErrors() {
		r.SaveState()
	}
	// Perform post-compilation tasks even if compiler.hasErrors()
	r.compiler.PerformPostCompilationTasks()
}

// Java: AbstractCommandLineRunner.runCompilerPasses(CompileMetricsRecorderInterface)
func (r *AbstractCommandLineRunnerBase) RunCompilerPasses(metricsRecorder CompileMetricsRecorderInterface) {
	config := r.config
	runStage1 := false
	runStage2 := false
	runStage3 := false
	instrumentForCoverage := false
	var actionMetricsName string
	if r.compiler.GetOptions().GetInstrumentForCoverageOnly() {
		actionMetricsName = "instrument for coverage"
		instrumentForCoverage = true
	} else if config.ShouldSaveAfterStage1() {
		actionMetricsName = "stage 1"
		runStage1 = true
	} else if config.ShouldRestoreTypedAstsPerformStage2AndSave() {
		actionMetricsName = "parse & optimize"
		runStage2 = true
	} else if config.ShouldRestoreTypedAstsPerformStages2And3() {
		actionMetricsName = "skip-checks compile"
		runStage2 = true
		runStage3 = true
	} else if config.ShouldRestoreAndPerformStage2AndSave() {
		actionMetricsName = "stage 2/3"
		runStage2 = true
	} else if config.ShouldRestoreAndPerformStages2And3() {
		// From the outside this looks like the second stage of a 2-stage compile.
		actionMetricsName = "stage 2/2"
		runStage2 = true
		runStage3 = true
	} else if config.ShouldRestoreAndPerformStage3() {
		actionMetricsName = "stage 3/3"
		runStage3 = true
	} else {
		// This is the code path taken when "building" a library by just checking it for errors
		// and generating an .ijs file and also when doing a full compilation.
		if r.compiler.GetOptions().IsChecksOnly() {
			actionMetricsName = "checks-only"
		} else {
			actionMetricsName = "full compile"
		}
		runStage1 = true
		runStage2 = true
		runStage3 = true
	}

	metricsRecorder.RecordActionName(actionMetricsName)
	if r.compiler.HasErrors() {
		return
	}
	metricsRecorder.RecordStartState(r.compiler)
	if runStage1 {
		r.compiler.Stage1Passes()
		if r.compiler.HasErrors() {
			return
		}
	}
	if runStage2 {
		r.compiler.Stage2Passes(CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS)
		if r.compiler.HasErrors() {
			return
		}
	}
	if runStage3 {
		r.compiler.Stage3Passes()
		if r.compiler.HasErrors() {
			return
		}
	}
	if instrumentForCoverage {
		r.compiler.InstrumentForCoverage()
	}
}

// GetCompileMetricsRecorder: child classes should override this if they
// want to actually record metrics about the compilation.
//
// Java: AbstractCommandLineRunner.getCompileMetricsRecorder()
func (r *AbstractCommandLineRunnerBase) GetCompileMetricsRecorder() CompileMetricsRecorderInterface {
	return NewDummyCompileMetricsRecorder()
}

// Java: AbstractCommandLineRunner.initWithTypedAstFilesystem(List<SourceFile>, List<SourceFile>, CompilerOptions, String)
func (r *AbstractCommandLineRunnerBase) InitWithTypedAstFilesystem(
	externs []*SourceFile,
	sources []*SourceFile,
	options *CompilerOptions,
	filename string) {
	jutil.Try(func() {
		typedAstListStream := jutil.NewGZIPInputStreamWithSize(jutil.NewFileInputStream(filename), abstractCommandLineRunner_GZIPPED_TYPEDAST_BUFFER_SIZE)
		defer typedAstListStream.Close()
		r.compiler.InitWithTypedAstFilesystem(externs, sources, options, typedAstListStream)
	}, func(e jutil.IIOException) {
		r.compiler.Report(JSError_Make(AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST, filename))
	})
}

// Java: AbstractCommandLineRunner.initChunksWithTypedAstFilesystem(List<SourceFile>, List<JSChunk>, CompilerOptions, String)
func (r *AbstractCommandLineRunnerBase) InitChunksWithTypedAstFilesystem(
	externs []*SourceFile, chunks []*JSChunk, options *CompilerOptions, filename string) {
	jutil.Try(func() {
		typedAstListStream := jutil.NewGZIPInputStreamWithSize(jutil.NewFileInputStream(filename), abstractCommandLineRunner_GZIPPED_TYPEDAST_BUFFER_SIZE)
		defer typedAstListStream.Close()
		r.compiler.InitChunksWithTypedAstFilesystem(externs, chunks, options, typedAstListStream)
	}, func(e jutil.IIOException) {
		r.compiler.Report(JSError_Make(AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST, filename))
	})
}

// InitializeStateBeforeCompilation is called at the beginning of
// compilation to initialize the compiler state.
//
// Compiler state should be initialized no matter what the compilation stage
// is, but this method handles the different ways it might be initialized -
// whether from parsing actual JS files, from reading library-level
// TypedASTs, or from restoring a previous compilation state.
//
// Java: AbstractCommandLineRunner.initializeStateBeforeCompilation()
func (r *AbstractCommandLineRunnerBase) InitializeStateBeforeCompilation() {
	if r.config.restoredCompilationStage != -1 {
		r.RestoreState(*r.config.GetContinueSavedCompilationFileName())
	} else if r.config.typedAstListInputFilename != nil {
		// we did this elsewhere
	} else {
		// parsing!
		r.compiler.ParseForCompilation()
	}
}

// Java: AbstractCommandLineRunner.restoreState(String)
func (r *AbstractCommandLineRunnerBase) RestoreState(filename string) {
	jutil.Try(func() {
		serializedInputStream := jutil.NewFileInputStream(filename)
		defer serializedInputStream.Close()
		r.compiler.RestoreState(serializedInputStream)
	}, func(e jutil.Throwable) {
		_, isIO := e.(jutil.IIOException)
		_, isCNF := e.(jutil.IClassNotFoundException)
		if !isIO && !isCNF {
			panic(e)
		}
		r.compiler.Report(JSError_Make(AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST, filename))
	})
}

// SaveState is called at the end of compilation to save the compiler state
// if applicable.
//
// Java: AbstractCommandLineRunner.saveState()
func (r *AbstractCommandLineRunnerBase) SaveState() {
	if r.config.GetSaveCompilationStateToFilename() == nil {
		// nothing to save to.
		return
	}
	filename := *r.config.GetSaveCompilationStateToFilename()
	jutil.Try(func() {
		serializedOutputStream := jutil.NewFileOutputStream(filename)
		defer serializedOutputStream.Close()
		r.compiler.SaveState(serializedOutputStream)
	}, func(e jutil.IIOException) {
		r.compiler.Report(JSError_Make(AbstractCommandLineRunner_COULD_NOT_SERIALIZE_AST, filename))
	})
}

// ProcessResults processes the results of the compile job, and returns an
// error code.
//
// Java: AbstractCommandLineRunner.processResults(Result, List<JSChunk>, B)
func (r *AbstractCommandLineRunnerBase) ProcessResults(result *Result, chunks []*JSChunk, options *CompilerOptions) int {
	config := r.config
	if config.printAst {
		if r.compiler.GetRoot() == nil {
			return 1
		} else {
			jsOutput := r.CreateDefaultOutput()
			cfg := r.compiler.ComputeCFG()
			DotFormatter_AppendDot(r.compiler.GetRoot().GetLastChild(), cfg, jsOutput)
			jsOutput.AppendCharacter('\n')
			abstractCommandLineRunner_closeAppendable(jsOutput)
			return 0
		}
	}

	if config.printTree {
		if r.compiler.GetRoot() == nil {
			r.compiler.Report(JSError_Make(AbstractCommandLineRunner_NO_TREE_GENERATED_ERROR))
			return 1
		} else {
			jsOutput := r.CreateDefaultOutput()
			r.compiler.GetJsRoot().AppendStringTree(jsOutput)
			jsOutput.AppendCharSequence("\n")
			abstractCommandLineRunner_closeAppendable(jsOutput)
			return 0
		}
	}

	if config.printTreeJson {
		if r.compiler.GetRoot() == nil {
			r.compiler.Report(JSError_Make(AbstractCommandLineRunner_NO_TREE_GENERATED_ERROR))
			return 1
		} else {
			jsOutput := r.CreateDefaultOutput()
			r.compiler.GetJsRoot().AppendJsonTree(jsOutput)
			abstractCommandLineRunner_closeAppendable(jsOutput)
			return 0
		}
	}

	if config.skipNormalOutputs {
		// Output the manifest and bundle files if requested.
		r.OutputManifest()
		r.OutputBundle()
		r.OutputChunkGraphJson()
		return 0
	} else if options.GetOutputJs() != CompilerOptions_OutputJs_NONE && result.success {
		r.OutputChunkGraphJson()
		if chunks == nil {
			r.self.OutputSingleBinary(options)

			// Output the source map if requested.
			// If output files are being written to stdout as a JSON string,
			// outputSingleBinary will have added the sourcemap to the output file
			if !r.IsOutputInJson() {
				r.OutputSourceMap(options, config.jsOutputFile)
			}
		} else {
			error := r.OutputChunkBinaryAndSourceMaps(r.compiler.GetChunkGraph(), options)
			if error != nil {
				r.compiler.Report(JSError_Make(error))
				return 1
			}
		}

		// Output the externs if required.
		if options.GetExternExportsPath() != nil {
			func() {
				eeOut := r.OpenExternExportsStream(options, config.jsOutputFile)
				defer abstractCommandLineRunner_closeWriter(eeOut)
				// Java: Writer.append(null) appends "null".
				externExport := "null"
				if result.externExport != nil {
					externExport = *result.externExport
				}
				eeOut.AppendCharSequence(externExport)
			}()
		}

		// Output the variable and property name maps if requested.
		r.OutputNameMaps()

		// Output the ReplaceStrings map if requested
		r.OutputStringMap()

		// Output the manifest and bundle files if requested.
		r.OutputManifest()
		r.OutputBundle()

		// Output the production instrumentation param mapping if requested.
		r.OutputInstrumentationMapping()

		if r.IsOutputInJson() {
			r.self.OutputJsonStream()
		}
	}

	return abstractCommandLineRunner_getExitStatusForResult(result)
}

// Java: AbstractCommandLineRunner.getExitStatusForResult(Result)
func abstractCommandLineRunner_getExitStatusForResult(result *Result) int {
	// return 0 if no errors, the error count otherwise
	return min(len(result.errors), 0x7f)
}

// Java: AbstractCommandLineRunner.getJavascriptEscaper()
func (r *AbstractCommandLineRunnerBase) GetJavascriptEscaper() func(string) string {
	return jutil.Escaper_AsFunction(SourceCodeEscapers_JavascriptEscaper())
}

// Java: AbstractCommandLineRunner.outputSingleBinary(B)
func (r *AbstractCommandLineRunnerBase) OutputSingleBinary(options *CompilerOptions) {
	config := r.config
	var escaper func(string) string
	marker := AbstractCommandLineRunner_OUTPUT_MARKER
	if strings.Contains(config.outputWrapper, abstractCommandLineRunner_OUTPUT_MARKER_JS_STRING) {
		marker = abstractCommandLineRunner_OUTPUT_MARKER_JS_STRING
		escaper = r.self.GetJavascriptEscaper()
	}

	if r.IsOutputInJson() {
		r.filesToStreamOut = append(r.filesToStreamOut, r.self.CreateJsonFile(options, marker, escaper))
	} else {
		if config.jsOutputFile != "" {
			abstractCommandLineRunner_maybeCreateDirsForPath(&config.jsOutputFile)
		}

		jsOutput := r.CreateDefaultOutput()
		r.self.WriteOutput(
			jsOutput,
			r.compiler,
			// So long as the JSChunk arg is null the compiler will write all sources to jsOutput
			// Use single-binary license tracking to dedupe licenses among all the inputs
			newCompiler_SingleBinaryLicenseTracker(r.compiler),
			nil,
			config.outputWrapper,
			marker,
			escaper,
			config.jsOutputFile)
		abstractCommandLineRunner_closeAppendable(jsOutput)
	}
}

// CreateJsonFile saves the compiler output to a JsonFileSpec to be later
// written to stdout.
//
// Java: AbstractCommandLineRunner.createJsonFile(B, String, Function<String, String>)
func (r *AbstractCommandLineRunnerBase) CreateJsonFile(options *CompilerOptions, outputMarker string, escaper func(string) string) *AbstractCommandLineRunner_JsonFileSpec {
	config := r.config
	jsOutput := jutil.NewStringBuilder()
	r.self.WriteOutput(
		jsOutput,
		r.compiler,
		// So long as the JSChunk arg is null the compiler will write all sources to jsOutput
		// Use single-binary license tracking to dedupe licenses among all the inputs
		newCompiler_SingleBinaryLicenseTracker(r.compiler),
		nil,
		config.outputWrapper,
		outputMarker,
		escaper,
		config.jsOutputFile)

	path := config.jsOutputFile
	if jutil.Strings_IsNullOrEmpty(&config.jsOutputFile) {
		path = "compiled.js"
	}
	jsonOutput := NewAbstractCommandLineRunner_JsonFileSpec(jsOutput.String(), path)

	if options.ShouldGatherSourceMapInfo() {
		sourcemap := &strings.Builder{} // Go-only: language-adapter — StringBuilder as SourceMap.appendTo's Appendable
		r.compiler.GetSourceMap().AppendTo(sourcemap, jsonOutput.GetPath())
		jsonOutput.SetSourceMap(sourcemap.String())
	}

	return jsonOutput
}

// Java: AbstractCommandLineRunner.outputJsonStream()
func (r *AbstractCommandLineRunnerBase) OutputJsonStream() {
	jsonWriter := jutil.NewJsonWriter(jutil.NewBufferedWriter(jutil.NewOutputStreamWriterWithCharset(
		abstractCommandLineRunner_printStreamAsOutputStream{r.defaultJsOutput}, jutil.StandardCharsets_UTF_8)))
	defer jsonWriter.Close()
	gsonOut := jutil.NewGsonBuilder().DisableHtmlEscaping().SetPrettyPrinting().Create()
	// Java: the TypeToken<List<JsChunkSpec>> is not used to write the elements:
	// Gson's runtime type wrapper writes each JsonFileSpec with its own
	// reflective adapter.
	gsonOut.ToJsonWithWriter(r.filesToStreamOut, jsonWriter)
}

// Java: AbstractCommandLineRunner.outputChunkBinaryAndSourceMaps(JSChunkGraph, B)
func (r *AbstractCommandLineRunnerBase) OutputChunkBinaryAndSourceMaps(
	chunkGraph *JSChunkGraph, options *CompilerOptions) *DiagnosticType {
	config := r.config
	chunks := chunkGraph.GetAllChunks()
	r.parsedChunkWrappers = AbstractCommandLineRunner_ParseChunkWrappers(config.chunkWrapper, chunks)
	if !r.IsOutputInJson() {
		// make sure the method generates all dirs up to the latest /
		dummy := config.chunkOutputPathPrefix + "dummy"
		abstractCommandLineRunner_maybeCreateDirsForPath(&dummy)
	}

	// If the source map path is in fact a pattern for each
	// chunk, create a stream per-chunk. Otherwise, create
	// a single source map.
	var mapFileOut AbstractCommandLineRunner_Writer

	// When the json_streams flag is specified, sourcemaps are always generated
	// per chunk
	if !(r.self.ShouldGenerateMapPerChunk(options) ||
		!options.ShouldGatherSourceMapInfo() ||
		config.jsonStreamMode == CompilerOptions_JsonStreamMode_OUT ||
		config.jsonStreamMode == CompilerOptions_JsonStreamMode_BOTH) {
		// warn that this is not supported
		return AbstractCommandLineRunner_INVALID_CHUNK_SOURCEMAP_PATTERN
	}

	mlicenseTracker := NewCompiler_ChunkGraphAwareLicenseTracker(r.compiler)
	for _, m := range chunks {
		if m.GetName() == JSChunk_WEAK_CHUNK_NAME {
			// Skip the weak chunk, which is always empty.
			continue
		}
		if r.IsOutputInJson() {
			r.filesToStreamOut = append(r.filesToStreamOut, r.CreateJsonFileFromChunk(m))
		} else {
			if r.self.ShouldGenerateMapPerChunk(options) {
				mapFileOut = r.FileNameToOutputWriter2(r.self.ExpandSourceMapPath(options, m))
			}

			chunkFilename := r.self.GetChunkOutputFileName(m)
			abstractCommandLineRunner_maybeCreateDirsForPath(&chunkFilename)
			func() {
				writer := r.FileNameToLegacyOutputWriter(&chunkFilename)
				defer abstractCommandLineRunner_closeWriter(writer)
				if options.ShouldGatherSourceMapInfo() {
					r.compiler.ResetAndIntitializeSourceMap()
				}
				mlicenseTracker.SetCurrentChunkContext(m)
				r.self.WriteChunkOutput(chunkFilename, writer, mlicenseTracker, m)
				if options.ShouldGatherSourceMapInfo() {
					r.compiler.GetSourceMap().AppendTo(mapFileOut, &chunkFilename)
				}
			}()

			if r.self.ShouldGenerateMapPerChunk(options) && mapFileOut != nil {
				mapFileOut.Close()
				mapFileOut = nil
			}
		}
	}

	if mapFileOut != nil {
		mapFileOut.Close()
	}
	return nil
}

// CreateJsonFileFromChunk converts the given output chunk to a JSONFileSpec
// with associated sourcemap.
//
// Java: AbstractCommandLineRunner.createJsonFileFromChunk(JSChunk)
func (r *AbstractCommandLineRunnerBase) CreateJsonFileFromChunk(chunk *JSChunk) *AbstractCommandLineRunner_JsonFileSpec {
	r.compiler.ResetAndIntitializeSourceMap()

	filename := r.self.GetChunkOutputFileName(chunk)
	output := jutil.NewStringBuilder()
	r.self.WriteChunkOutput(filename, output, NewCompiler_ScriptNodeLicensesOnlyTracker(r.compiler), chunk)

	jsonFile := NewAbstractCommandLineRunner_JsonFileSpec(output.String(), filename)

	chunkSourceMap := &strings.Builder{} // Go-only: language-adapter — StringBuilder as SourceMap.appendTo's Appendable

	chunkOutputFileName := r.self.GetChunkOutputFileName(chunk)
	r.compiler.GetSourceMap().AppendTo(chunkSourceMap, &chunkOutputFileName)

	jsonFile.SetSourceMap(chunkSourceMap.String())

	return jsonFile
}

// GetInputCharset queries the flag for the input charset, and returns a
// Charset object representing the selection. It throws FlagUsageException
// if flag is not a valid Charset name.
//
// Java: AbstractCommandLineRunner.getInputCharset()
func (r *AbstractCommandLineRunnerBase) GetInputCharset() *jutil.Charset {
	if r.config.charset != "" {
		if !jutil.Charset_IsSupported(r.config.charset) {
			panic(NewFlagUsageException(r.config.charset + " is not a valid charset name."))
		}
		return jutil.Charset_ForName(r.config.charset)
	}
	return jutil.StandardCharsets_UTF_8
}

// GetLegacyOutputCharset queries the flag for the output charset.
//
// Let the outputCharset be the same as the input charset... except if we're
// reading in UTF-8 by default. By tradition, we've always output ASCII to
// avoid various hiccups with different browsers, proxies and firewalls.
//
// It returns the charset to use when writing outputs (guaranteed to be a
// supported charset) and throws FlagUsageException if flag is not a valid
// Charset name.
//
// Java: AbstractCommandLineRunner.getLegacyOutputCharset()
func (r *AbstractCommandLineRunnerBase) GetLegacyOutputCharset() *jutil.Charset {
	if r.config.charset != "" {
		if !jutil.Charset_IsSupported(r.config.charset) {
			panic(NewFlagUsageException(r.config.charset + " is not a valid charset name."))
		}
		return jutil.Charset_ForName(r.config.charset)
	}
	return jutil.StandardCharsets_US_ASCII
}

// GetOutputCharset2 queries the flag for the output charset. Defaults to
// UTF-8. It throws FlagUsageException if flag is not a valid Charset name.
//
// Java: AbstractCommandLineRunner.getOutputCharset2()
func (r *AbstractCommandLineRunnerBase) GetOutputCharset2() *jutil.Charset {
	if r.config.charset != "" {
		if !jutil.Charset_IsSupported(r.config.charset) {
			panic(NewFlagUsageException(r.config.charset + " is not a valid charset name."))
		}
		return jutil.Charset_ForName(r.config.charset)
	}
	return jutil.StandardCharsets_UTF_8
}

// Java: AbstractCommandLineRunner.createExterns(CompilerOptions)
func (r *AbstractCommandLineRunnerBase) CreateExterns(options *CompilerOptions) []*SourceFile {
	if r.self.IsInTestMode() {
		return r.externsSupplierForTesting()
	}
	return r.CreateExternInputs(r.config.externs)
}

// ShouldGenerateMapPerChunk returns true if and only if a source map file
// should be generated for each chunk, as opposed to one unified map. This is
// specified by having the source map pattern include the %outname%
// variable.
//
// Java: AbstractCommandLineRunner.shouldGenerateMapPerChunk(B)
func (r *AbstractCommandLineRunnerBase) ShouldGenerateMapPerChunk(options *CompilerOptions) bool {
	return options.ShouldGatherSourceMapInfo() &&
		strings.Contains(*options.GetSourceMapOutputPath(), "%outname%")
}

// OpenExternExportsStream returns a stream for outputting the generated
// externs file.
//
// options are the options to the Compiler; path is the path of the
// generated JS source file. It returns the stream or null if no extern-ed
// exports are being generated.
//
// Java: AbstractCommandLineRunner.openExternExportsStream(B, String)
func (r *AbstractCommandLineRunnerBase) OpenExternExportsStream(options *CompilerOptions, path string) AbstractCommandLineRunner_Writer {
	externExportsPath := options.GetExternExportsPath()
	if externExportsPath == nil {
		return nil
	}

	exPath := *externExportsPath

	if !strings.Contains(exPath, jutil.File_separator) {
		outputFile := jutil.NewFile(path)
		parent := "null" // Java: a null parent concatenates as "null".
		if p := outputFile.GetParent(); p != nil {
			parent = *p
		}
		exPath = parent + string(rune(jutil.File_separatorChar)) + exPath
	}

	return r.FileNameToOutputWriter2(&exPath)
}

// ExpandCommandLinePath expands a file path specified on the command-line.
//
// Most file paths on the command-line allow an %outname% placeholder. The
// placeholder will expand to a different value depending on the current
// output mode. There are three scenarios:
//
// 1) Single JS output, single extra output: sub in jsOutputPath. 2) Multiple
// JS output, single extra output: sub in the base chunk name. 3) Multiple JS
// output, multiple extra output: sub in the module output file.
//
// Passing a JSChunk to this function automatically triggers case #3.
// Otherwise, we'll use strategy #1 or #2 based on the current output mode.
//
// Java: AbstractCommandLineRunner.expandCommandLinePath(String, JSChunk)
func (r *AbstractCommandLineRunnerBase) ExpandCommandLinePath(path string, forChunk *JSChunk) string {
	var sub string
	if forChunk != nil {
		sub = r.self.GetChunkOutputFileName(forChunk)
	} else if len(r.config.chunk) != 0 {
		sub = r.config.chunkOutputPathPrefix
	} else {
		sub = r.config.jsOutputFile
	}
	return strings.ReplaceAll(path, "%outname%", sub)
}

// ExpandSourceMapPath is the expansion function for source map.
//
// Java: AbstractCommandLineRunner.expandSourceMapPath(B, JSChunk)
func (r *AbstractCommandLineRunnerBase) ExpandSourceMapPath(options *CompilerOptions, forChunk *JSChunk) *string {
	if !options.ShouldGatherSourceMapInfo() {
		return nil
	}
	expanded := r.ExpandCommandLinePath(*options.GetSourceMapOutputPath(), forChunk)
	return &expanded
}

// FileNameToLegacyOutputWriter converts a file name into a Writer taking in
// account the output charset. Returns null if the file name is null.
//
// Java: AbstractCommandLineRunner.fileNameToLegacyOutputWriter(String)
func (r *AbstractCommandLineRunnerBase) FileNameToLegacyOutputWriter(fileName *string) AbstractCommandLineRunner_Writer {
	if fileName == nil {
		return nil
	}
	if r.self.IsInTestMode() {
		return jutil.NewStringWriter()
	}

	return r.StreamToLegacyOutputWriter(r.self.FilenameToOutputStream(fileName))
}

// FileNameToOutputWriter2 converts a file name into a Writer taking in
// account the output charset. Returns null if the file name is null.
//
// Java: AbstractCommandLineRunner.fileNameToOutputWriter2(String)
func (r *AbstractCommandLineRunnerBase) FileNameToOutputWriter2(fileName *string) AbstractCommandLineRunner_Writer {
	if fileName == nil {
		return nil
	}
	if r.self.IsInTestMode() {
		return jutil.NewStringWriter()
	}

	return r.StreamToOutputWriter2(r.self.FilenameToOutputStream(fileName))
}

// FilenameToOutputStream converts a file name into a Outputstream. Returns
// null if the file name is null.
//
// Java: AbstractCommandLineRunner.filenameToOutputStream(String)
func (r *AbstractCommandLineRunnerBase) FilenameToOutputStream(fileName *string) io.WriteCloser {
	if fileName == nil {
		return nil
	}
	return jutil.NewFileOutputStream(*fileName)
}

// StreamToLegacyOutputWriter creates a writer with the legacy output
// charset.
//
// Java: AbstractCommandLineRunner.streamToLegacyOutputWriter(OutputStream)
func (r *AbstractCommandLineRunnerBase) StreamToLegacyOutputWriter(stream io.WriteCloser) AbstractCommandLineRunner_Writer {
	if r.legacyOutputCharset == nil {
		return r.self.CreateWriter(stream, jutil.StandardCharsets_UTF_8)
	}
	return r.self.CreateWriter(stream, r.legacyOutputCharset)
}

// StreamToOutputWriter2 creates a writer with the newer output charset.
//
// Java: AbstractCommandLineRunner.streamToOutputWriter2(OutputStream)
func (r *AbstractCommandLineRunnerBase) StreamToOutputWriter2(stream io.WriteCloser) AbstractCommandLineRunner_Writer {
	if r.outputCharset2 == nil {
		return r.self.CreateWriter(stream, jutil.StandardCharsets_UTF_8)
	}
	return r.self.CreateWriter(stream, r.outputCharset2)
}

// CreateWriter creates a buffered Writer that writes to the given stream
// using the given encoding.
//
// Java: AbstractCommandLineRunner.createWriter(OutputStream, Charset)
func (r *AbstractCommandLineRunnerBase) CreateWriter(stream io.WriteCloser, charset *jutil.Charset) AbstractCommandLineRunner_Writer {
	return jutil.NewBufferedWriter(jutil.NewOutputStreamWriterWithCharset(stream, charset))
}

// OutputSourceMap outputs the source map found in the compiler to the
// proper path if one exists.
//
// Java: AbstractCommandLineRunner.outputSourceMap(B, String)
func (r *AbstractCommandLineRunnerBase) OutputSourceMap(options *CompilerOptions, associatedName string) {
	if !options.ShouldGatherSourceMapInfo() ||
		*options.GetSourceMapOutputPath() == "/dev/null" {
		return
	}

	outName := r.self.ExpandSourceMapPath(options, nil)
	abstractCommandLineRunner_maybeCreateDirsForPath(outName)
	out := r.FileNameToOutputWriter2(outName)
	defer abstractCommandLineRunner_closeWriter(out)
	r.compiler.GetSourceMap().AppendTo(out, &associatedName)
}

// GetMapPath returns the path at which to output map file(s) based on the
// path at which the JS binary will be placed. It returns the path in which
// to place the generated map file(s).
//
// Java: AbstractCommandLineRunner.getMapPath(String)
func (r *AbstractCommandLineRunnerBase) GetMapPath(outputFile string) string {
	basePath := ""

	if outputFile == "" {
		// If we have a js_module_binary rule, output the maps
		// at chunkname_props_map.out, etc.
		if r.config.chunkOutputPathPrefix != "" {
			basePath = r.config.chunkOutputPathPrefix
		} else {
			basePath = "jscompiler"
		}
	} else {
		// Get the path of the output file.
		file := jutil.NewFile(outputFile)

		outputFileName := file.GetName()

		// Strip the .js from the name.
		if strings.HasSuffix(outputFileName, ".js") {
			outputFileName = outputFileName[:len(outputFileName)-3]
		}

		fileParent := file.GetParent()
		if fileParent == nil {
			basePath = outputFileName
		} else {
			basePath = *file.GetParent() + string(rune(jutil.File_separatorChar)) + outputFileName
		}
	}

	return basePath
}

// OutputNameMaps outputs the variable and property name maps for the
// specified compiler if the proper FLAGS are set.
//
// Java: AbstractCommandLineRunner.outputNameMaps()
func (r *AbstractCommandLineRunnerBase) OutputNameMaps() {
	config := r.config

	var propertyMapOutputPath *string
	var variableMapOutputPath *string

	// Check the create_name_map_files FLAG.
	if config.createNameMapFiles {
		basePath := r.GetMapPath(config.jsOutputFile)

		propertyMapOutputPath = jutil.StringPtr(basePath + "_props_map.out")
		variableMapOutputPath = jutil.StringPtr(basePath + "_vars_map.out")
	}

	// Check the individual FLAGS.
	if config.variableMapOutputFile != "" {
		if variableMapOutputPath != nil {
			panic(NewFlagUsageException(
				"The flags variable_map_output_file and " +
					"create_name_map_files cannot both be used simultaneously."))
		}

		variableMapOutputPath = &config.variableMapOutputFile
	}

	if config.propertyMapOutputFile != "" {
		if propertyMapOutputPath != nil {
			panic(NewFlagUsageException(
				"The flags property_map_output_file and " +
					"create_name_map_files cannot both be used simultaneously."))
		}

		propertyMapOutputPath = &config.propertyMapOutputFile
	}

	// Output the maps.
	if variableMapOutputPath != nil && r.compiler.GetVariableMap() != nil {
		r.compiler.GetVariableMap().Save(*variableMapOutputPath)
	}

	if propertyMapOutputPath != nil && r.compiler.GetPropertyMap() != nil {
		r.compiler.GetPropertyMap().Save(*propertyMapOutputPath)
	}
}

// OutputStringMap outputs the string map generated by the ReplaceStrings
// pass if an output path exists.
//
// Java: AbstractCommandLineRunner.outputStringMap()
func (r *AbstractCommandLineRunnerBase) OutputStringMap() {
	config := r.config
	if config.stringMapOutputPath != "" {
		if r.compiler.GetStringMap() == nil {
			// Ensure an empty file is created if there is no string map.
			// This avoids confusing some build tools that expect to see the file, even if it is empty.
			if !jutil.NewFile(config.stringMapOutputPath).CreateNewFile() {
				panic(jutil.NewIOException("Could not create file: " + config.stringMapOutputPath))
			}
		} else {
			r.compiler.GetStringMap().Save(config.stringMapOutputPath)
		}
	}
}

// AbstractCommandLineRunner_CreateDefineReplacements creates a map of
// constant names to constant values from a textual description of the map.
//
// definitions is a list of overriding definitions for defines in the form
// <name>[=<val>], where <val> is a number, boolean, or single-quoted string
// without single quotes.
//
// Java: AbstractCommandLineRunner.createDefineReplacements(List<String>, CompilerOptions)
func AbstractCommandLineRunner_CreateDefineReplacements(definitions []string, options *CompilerOptions) {
	// Parse the definitions
	for _, override := range definitions {
		assignment := jutil.StringSplitLiteralLimit(override, "=", 2)
		defName := assignment[0]

		if jutil.StringLength(defName) > 0 {
			defValue := "true"
			if len(assignment) != 1 {
				defValue = assignment[1]
			}

			isTrue := defValue == "true"
			isFalse := defValue == "false"
			if isTrue || isFalse {
				options.SetDefineToBooleanLiteral(defName, isTrue)
				continue
			} else if jutil.StringLength(defValue) > 1 &&
				((jutil.StringCharAt(defValue, 0) == '\'' && jutil.StringCharAt(defValue, jutil.StringLength(defValue)-1) == '\'') ||
					(jutil.StringCharAt(defValue, 0) == '"' &&
						jutil.StringCharAt(defValue, jutil.StringLength(defValue)-1) == '"')) {
				// If the value starts and ends with a single quote,
				// we assume that it's a string.
				maybeStringVal := jutil.StringSubstring(defValue, 1, jutil.StringLength(defValue)-1)
				if jutil.StringIndexOfChar(maybeStringVal, jutil.StringCharAt(defValue, 0)) == -1 {
					options.SetDefineToStringLiteral(defName, maybeStringVal)
					continue
				}
			} else {
				parsed := false
				jutil.Try(func() {
					value := jutil.DoubleParseDouble(defValue)
					options.SetDefineToDoubleLiteral(defName, value)
					parsed = true
				}, func(e jutil.INumberFormatException) {
					// do nothing, it will be caught at the end
				})
				if parsed {
					continue
				}

				if jutil.StringLength(defValue) > 0 {
					options.SetDefineToStringLiteral(defName, defValue)
					continue
				}
			}
		}

		panic(jutil.NewRuntimeException("--define flag syntax invalid: " + override))
	}
}

// ShouldGenerateOutputPerChunk returns true if and only if a manifest or
// bundle should be generated for each chunk, as opposed to one unified
// manifest.
//
// Java: AbstractCommandLineRunner.shouldGenerateOutputPerChunk(String)
func (r *AbstractCommandLineRunnerBase) ShouldGenerateOutputPerChunk(output string) bool {
	return len(r.config.chunk) != 0 && strings.Contains(output, "%outname%")
}

// Java: AbstractCommandLineRunner.outputManifest()
func (r *AbstractCommandLineRunnerBase) OutputManifest() {
	r.OutputManifestOrBundle(r.config.outputManifests, true)
}

// Java: AbstractCommandLineRunner.outputBundle()
func (r *AbstractCommandLineRunnerBase) OutputBundle() {
	r.OutputManifestOrBundle(r.config.outputBundles, false)
}

// Java: AbstractCommandLineRunner.outputInstrumentationMapping()
func (r *AbstractCommandLineRunnerBase) OutputInstrumentationMapping() {
	if !jutil.Strings_IsNullOrEmpty(&r.config.instrumentationMappingFile) {
		path := r.ExpandCommandLinePath(r.config.instrumentationMappingFile /* forChunk= */, nil)
		r.compiler.GetInstrumentationMapping().Save(path)
	}
}

// OutputManifestOrBundle writes the manifest or bundle of all compiler
// input files that were included as controlled by --dependency_mode, if
// requested.
//
// Java: AbstractCommandLineRunner.outputManifestOrBundle(List<String>, boolean)
func (r *AbstractCommandLineRunnerBase) OutputManifestOrBundle(outputFiles []string, isManifest bool) {
	if len(outputFiles) == 0 {
		return
	}

	for _, output := range outputFiles {
		if output == "" {
			continue
		}

		if r.ShouldGenerateOutputPerChunk(output) {
			// Generate per-chunk manifests or bundles.
			chunks := r.compiler.GetChunkGraph().GetAllChunks()
			for _, chunk := range chunks {
				func() {
					path := r.ExpandCommandLinePath(output, chunk)
					out := r.FileNameToOutputWriter2(&path)
					defer abstractCommandLineRunner_closeWriter(out)
					if isManifest {
						r.self.PrintManifestTo(chunk, out)
					} else {
						r.self.PrintBundleTo(chunk, out)
					}
				}()
			}
		} else {
			// Generate a single file manifest or bundle.
			func() {
				path := r.ExpandCommandLinePath(output, nil)
				out := r.FileNameToOutputWriter2(&path)
				defer abstractCommandLineRunner_closeWriter(out)
				if len(r.config.chunk) == 0 {
					// For a single-chunk compilation, generate a single headerless manifest or bundle
					// containing only the strong files.
					chunk := r.compiler.GetChunkGraph().GetChunkByName(JSChunk_STRONG_CHUNK_NAME)
					if isManifest {
						r.self.PrintManifestTo(chunk, out)
					} else {
						r.self.PrintBundleTo(chunk, out)
					}
				} else {
					// For a multi-chunk compilation, generate a single manifest file with chunk headers.
					r.self.PrintChunkGraphManifestOrBundleTo(r.compiler.GetChunkGraph(), out, isManifest)
				}
			}()
		}
	}
}

// OutputChunkGraphJson creates a file containing the current chunk graph in
// JSON serialization.
//
// Java: AbstractCommandLineRunner.outputChunkGraphJson()
func (r *AbstractCommandLineRunnerBase) OutputChunkGraphJson() {
	if r.config.outputChunkDependencies != nil && jutil.StringLength(*r.config.outputChunkDependencies) != 0 {
		out := r.FileNameToOutputWriter2(r.config.outputChunkDependencies)
		defer abstractCommandLineRunner_closeWriter(out)
		r.self.PrintChunkGraphJsonTo(out)
	}
}

// PrintChunkGraphJsonTo prints the current chunk graph as JSON.
//
// Java: AbstractCommandLineRunner.printChunkGraphJsonTo(Appendable)
func (r *AbstractCommandLineRunnerBase) PrintChunkGraphJsonTo(out jutil.Appendable) {
	out.AppendCharSequence(r.compiler.GetChunkGraph().ToJson().String())
}

// PrintChunkGraphManifestOrBundleTo prints a set of chunks to the manifest
// or bundle file.
//
// Java: AbstractCommandLineRunner.printChunkGraphManifestOrBundleTo(JSChunkGraph, Appendable, boolean)
func (r *AbstractCommandLineRunnerBase) PrintChunkGraphManifestOrBundleTo(graph *JSChunkGraph, out jutil.Appendable, isManifest bool) {
	requiresNewline := false
	for _, chunk := range graph.GetAllChunks() {
		if !isManifest && chunk.IsWeak() {
			// Skip the weak chunk on a multi-chunk bundle, but not a multi-chunk manifest.
			continue
		}

		if requiresNewline {
			out.AppendCharSequence("\n")
		}

		if isManifest {
			// See CommandLineRunnerTest to see what the format of this
			// manifest looks like.
			dependencies := strings.Join(chunk.GetSortedDependencyNames(), ",")
			deps := ""
			if dependencies != "" {
				deps = ":" + dependencies
			}
			out.AppendCharSequence("{" + chunk.GetName() + deps + "}\n")
			r.self.PrintManifestTo(chunk, out)
		} else {
			r.self.PrintBundleTo(chunk, out)
		}
		requiresNewline = true
	}
}

// PrintManifestTo prints a list of input names (using root-relative paths),
// delimited by newlines, to the manifest file.
//
// Java: AbstractCommandLineRunner.printManifestTo(JSChunk, Appendable)
func (r *AbstractCommandLineRunnerBase) PrintManifestTo(chunk *JSChunk, out jutil.Appendable) {
	for _, input := range chunk.GetInputs() {
		name := input.GetName()

		// Ignore fill files created by the compiler to facilitate cross-chunk code motion.
		if AbstractCompiler_IsFillFileName(name) {
			continue
		}

		displayName := name
		if rootRelativePath, ok := r.rootRelativePathsMap.Lookup(name); ok {
			displayName = rootRelativePath
		}
		out.AppendCharSequence(displayName)
		out.AppendCharSequence("\n")
	}
}

// PrintBundleTo prints all the input contents, starting with a comment that
// specifies the input file name (using root-relative paths) before each
// file.
//
// Java: AbstractCommandLineRunner.printBundleTo(JSChunk, Appendable)
func (r *AbstractCommandLineRunnerBase) PrintBundleTo(chunk *JSChunk, out jutil.Appendable) {
	inputs := chunk.GetInputs()
	options := r.compiler.GetOptions()
	// Prebuild ASTs before they're needed in getLoadFlags, for performance and because
	// StackOverflowErrors can be hit if not prebuilt.
	if options.GetNumParallelThreads() > 1 {
		NewPrebuildAst(r.compiler, r.compiler.GetOptions().GetNumParallelThreads()).Prebuild(slices.Values(inputs))
	}
	if options.GetRuntimeLibraryMode() == RuntimeJsLibManager_RuntimeLibraryMode_INJECT {
		// ES6 modules will need a runtime in a bundle. Skip appending this runtime if there are no
		// ES6 modules to cut down on size.
		for _, input := range inputs {
			if input.IsEs6Module() {
				r.self.AppendRuntimeTo(out)
				break
			}
		}
	}

	for _, input := range inputs {
		name := input.GetName()
		code := input.GetSourceFile().GetCode()

		// Ignore empty fill files created by the compiler to facilitate cross-chunk code motion.
		// Note that non-empty fill files (ones whose code has actually been moved into) are still
		// emitted. In particular, this ensures that if there are no (real) inputs the bundle will be
		// empty.
		if AbstractCompiler_IsFillFileName(name) && code == "" {
			continue
		}

		displayName := input.GetName()
		if rootRelativePath, ok := r.rootRelativePathsMap.Lookup(name); ok {
			displayName = rootRelativePath
		}
		out.AppendCharSequence("//")
		out.AppendCharSequence(displayName)
		out.AppendCharSequence("\n")

		r.self.PrepForBundleAndAppendTo(out, input, code)

		out.AppendCharSequence("\n")
	}
}

// ConstructRootRelativePathsMap constructs and returns the input root path
// map. The key is the exec path of each input file, and the value is the
// corresponding root relative path.
//
// Java: AbstractCommandLineRunner.constructRootRelativePathsMap()
func (r *AbstractCommandLineRunnerBase) ConstructRootRelativePathsMap() *jutil.LinkedHashMap[string, string] {
	rootRelativePathsMap := jutil.NewLinkedHashMap[string, string]()
	for _, mapString := range r.config.manifestMaps {
		colonIndex := jutil.StringIndexOfChar(mapString, ':')
		jutil.CheckState(colonIndex > 0)
		execPath := jutil.StringSubstring(mapString, 0, colonIndex)
		rootRelativePath := jutil.StringSubstringFrom(mapString, colonIndex+1)
		jutil.CheckState(jutil.StringIndexOfChar(rootRelativePath, ':') == -1)
		rootRelativePathsMap.Put(execPath, rootRelativePath)
	}
	return rootRelativePathsMap
}

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

	jscompDevMode CompilerOptions_DevMode

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

	typedAstListInputFilename *string

	saveCompilationStateToFilename *string

	chunk []string

	sourceMapInputFiles *jutil.LinkedHashMap[string, string]

	parseInlineSourceMaps bool

	expectedDiagnostics []string

	variableMapInputFile string

	propertyMapInputFile string

	variableMapOutputFile string

	createNameMapFiles bool

	propertyMapOutputFile string

	stringMapOutputPath string

	instrumentationMappingFile string

	codingConvention CodingConvention

	summaryDetailLevel int

	outputWrapper string

	chunkWrapper []string

	chunkOutputPathPrefix string

	chunkOutputFiles []string

	chunkConformanceFiles []string

	createSourceMap string

	sourceMapDetailLevel SourceMap_DetailLevel

	sourceMapFormat SourceMap_Format

	sourceMapLocationMappings []SourceMap_LocationMapping

	applyInputSourceMaps bool

	warningGuards []*AbstractCommandLineRunner_FlagEntry[CheckLevel]

	define []string

	browserFeaturesetYear int

	tweakProcessing CompilerOptions_TweakProcessing

	charset string

	dependencyOptions *DependencyOptions // nil is null

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

	jsonStreamMode CompilerOptions_JsonStreamMode

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
		printTree:                  false,
		printTreeJson:              false,
		printAst:                   false,
		jscompDevMode:              CompilerOptions_DevMode_OFF,
		loggingLevel:               jutil.Level_WARNING.GetName(),
		externs:                    []string{},
		mixedJsSources:             []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]{},
		defaultToStdin:             false,
		jsOutputFile:               "",
		restoredCompilationStage:   -1,
		saveAfterCompilationStage:  -1,
		chunk:                      []string{},
		sourceMapInputFiles:        jutil.NewLinkedHashMap[string, string](),
		parseInlineSourceMaps:      false,
		expectedDiagnostics:        []string{},
		variableMapInputFile:       "",
		propertyMapInputFile:       "",
		variableMapOutputFile:      "",
		createNameMapFiles:         false,
		propertyMapOutputFile:      "",
		stringMapOutputPath:        "",
		instrumentationMappingFile: "",
		codingConvention:           CodingConventions_GetDefault(),
		summaryDetailLevel:         1,
		outputWrapper:              "",
		chunkWrapper:               []string{},
		chunkOutputPathPrefix:      "",
		chunkOutputFiles:           []string{},
		chunkConformanceFiles:      []string{},
		createSourceMap:            "",
		sourceMapDetailLevel:       SourceMap_DetailLevel_ALL,
		sourceMapFormat:            SourceMap_Format_DEFAULT,
		sourceMapLocationMappings:  []SourceMap_LocationMapping{},
		applyInputSourceMaps:       false,
		warningGuards:              []*AbstractCommandLineRunner_FlagEntry[CheckLevel]{},
		define:                     []string{},
		browserFeaturesetYear:      0,
		tweakProcessing:            CompilerOptions_TweakProcessing_OFF,
		charset:                    "",
		outputManifests:            []string{},
		outputBundles:              []string{},
		skipNormalOutputs:          false,
		manifestMaps:               []string{},
		processCommonJSModules:     false,
		moduleRoots:                []string{ModuleLoader_DEFAULT_FILENAME_PREFIX},
		warningsAllowFile:          "",
		hideWarningsFor:            []string{},
		angularPass:                false,
		jsonStreamMode:             CompilerOptions_JsonStreamMode_NONE,
		errorFormat:                AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_STANDARD,
		jsonWarningsFile:           "",
		includeTrailingNewline:     true,
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
func (c *AbstractCommandLineRunner_CommandLineConfig) SetJscompDevMode(jscompDevMode CompilerOptions_DevMode) *AbstractCommandLineRunner_CommandLineConfig {
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
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSourceMapInputFiles(sourceMapInputFiles *jutil.LinkedHashMap[string, string]) *AbstractCommandLineRunner_CommandLineConfig {
	c.sourceMapInputFiles = sourceMapInputFiles
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
func (c *AbstractCommandLineRunner_CommandLineConfig) SetCodingConvention(codingConvention CodingConvention) *AbstractCommandLineRunner_CommandLineConfig {
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
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSourceMapDetailLevel(level SourceMap_DetailLevel) *AbstractCommandLineRunner_CommandLineConfig {
	c.sourceMapDetailLevel = level
	return c
}

// SetSourceMapFormat sets the source map format to use, if generated.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setSourceMapFormat(SourceMap.Format)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSourceMapFormat(format SourceMap_Format) *AbstractCommandLineRunner_CommandLineConfig {
	c.sourceMapFormat = format
	return c
}

// SetSourceMapLocationMappings sets the source map location mappings to
// use, if generated.
//
// Java: AbstractCommandLineRunner.CommandLineConfig.setSourceMapLocationMappings(List<SourceMap.LocationMapping>)
func (c *AbstractCommandLineRunner_CommandLineConfig) SetSourceMapLocationMappings(
	locationMappings []SourceMap_LocationMapping) *AbstractCommandLineRunner_CommandLineConfig {

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
func (c *AbstractCommandLineRunner_CommandLineConfig) SetWarningGuards(warningGuards []*AbstractCommandLineRunner_FlagEntry[CheckLevel]) *AbstractCommandLineRunner_CommandLineConfig {
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
func (c *AbstractCommandLineRunner_CommandLineConfig) SetTweakProcessing(tweakProcessing CompilerOptions_TweakProcessing) *AbstractCommandLineRunner_CommandLineConfig {
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
func (c *AbstractCommandLineRunner_CommandLineConfig) SetDependencyOptions(dependencyOptions *DependencyOptions) *AbstractCommandLineRunner_CommandLineConfig {
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
func (c *AbstractCommandLineRunner_CommandLineConfig) SetJsonStreamMode(mode CompilerOptions_JsonStreamMode) *AbstractCommandLineRunner_CommandLineConfig {
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
// Go-only: language-adapter — the `gson` tags are the @SerializedName
// values and `gsonalt` their alternates (jutil Gson's reflective adapter);
// a nil field is Java's null.
//
// Java: AbstractCommandLineRunner.JsonFileSpec
type AbstractCommandLineRunner_JsonFileSpec struct {
	src  *string
	path *string

	sourceMap *string `gson:"source_map" gsonalt:"sourceMap"`

	webpackId *string `gson:"webpack_id" gsonalt:"webpackId"`
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
		src:       src,
		path:      path,
		sourceMap: sourceMap,
		webpackId: webpackId,
	}
}

// Java: AbstractCommandLineRunner.JsonFileSpec.getSrc()
func (s *AbstractCommandLineRunner_JsonFileSpec) GetSrc() *string {
	return s.src
}

// Java: AbstractCommandLineRunner.JsonFileSpec.getPath()
func (s *AbstractCommandLineRunner_JsonFileSpec) GetPath() *string {
	return s.path
}

// Java: AbstractCommandLineRunner.JsonFileSpec.getSourceMap()
func (s *AbstractCommandLineRunner_JsonFileSpec) GetSourceMap() *string {
	return s.sourceMap
}

// Java: AbstractCommandLineRunner.JsonFileSpec.getWebpackId()
func (s *AbstractCommandLineRunner_JsonFileSpec) GetWebpackId() *string {
	return s.webpackId
}

// Java: AbstractCommandLineRunner.JsonFileSpec.setSourceMap(String)
func (s *AbstractCommandLineRunner_JsonFileSpec) SetSourceMap(map_ string) {
	s.sourceMap = &map_
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

// AbstractCommandLineRunner_JsSourceType_values is JsSourceType.values().
func AbstractCommandLineRunner_JsSourceType_values() []AbstractCommandLineRunner_JsSourceType {
	return []AbstractCommandLineRunner_JsSourceType{
		AbstractCommandLineRunner_JsSourceType_EXTERN,
		AbstractCommandLineRunner_JsSourceType_JS,
		AbstractCommandLineRunner_JsSourceType_JS_ZIP,
		AbstractCommandLineRunner_JsSourceType_WEAKDEP,
		AbstractCommandLineRunner_JsSourceType_IJS,
	}
}

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
		panic(NewFlagUsageException(
			"Expected 2-4 colon-delimited parts in " + "chunk spec: " + specString))
	}

	// Parse chunk name.
	name := parts[0]

	// Parse chunk dependencies.
	deps := []string{}
	if len(parts) > 2 && jutil.StringLength(parts[2]) > 0 {
		deps = jutil.StringSplitLiteral(parts[2], ",")
	}

	// Parse chunk inputs.
	numInputs := -1
	jutil.Try(func() {
		numInputs = int(jutil.IntegerParseInt(parts[1], 10))
	}, func(ignored jutil.INumberFormatException) {
		numInputs = -1
	})

	// We will allow chunks of zero input.
	if numInputs < 0 {
		// A size of 'auto' is only allowed on the base chunk if
		// and it must also be the first chunk
		if len(parts) == 2 && "auto" == parts[1] {
			if !isFirstChunk {
				panic(NewFlagUsageException(
					"Invalid JS file count '" +
						parts[1] +
						"' for chunk: " +
						name +
						". Only the first chunk may specify " +
						"a size of 'auto' and it must have no dependencies."))
			}
		} else {
			panic(NewFlagUsageException(
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

// AbstractCommandLineRunner_Writer is a java.io.Writer as this class uses
// it: Writer.write/flush/close, the Appendable methods it is passed through,
// and io.StringWriter for SourceMap.appendTo.
//
// Go-only: language-adapter — java.io.Writer implements Appendable.
type AbstractCommandLineRunner_Writer interface {
	jutil.Writer
	jutil.Appendable
	io.StringWriter
}

// abstractCommandLineRunner_closeWriter is the close() of a
// try-with-resources statement, which skips a null resource.
//
// Go-only: language-adapter — try-with-resources.
func abstractCommandLineRunner_closeWriter(w AbstractCommandLineRunner_Writer) {
	if w != nil {
		w.Close()
	}
}

// abstractCommandLineRunner_printStreamAsOutputStream passes a PrintStream
// where Java passes it as an OutputStream: bytes go through
// PrintStream.write(byte[]) unchanged, and close() closes the PrintStream.
//
// Go-only: language-adapter — PrintStream is an OutputStream.
type abstractCommandLineRunner_printStreamAsOutputStream struct {
	ps *jutil.PrintStream
}

func (s abstractCommandLineRunner_printStreamAsOutputStream) Write(b []byte) (int, error) {
	return s.ps.Write(b)
}

func (s abstractCommandLineRunner_printStreamAsOutputStream) Close() error {
	s.ps.Close()
	return nil
}
