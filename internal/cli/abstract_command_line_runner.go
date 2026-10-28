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
	"archive/zip"
	"bytes"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
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
	// Java: AbstractCommandLineRunner.setExitCodeReceiver(Function<Integer, Void>)
	SetExitCodeReceiver(newExitCodeReceiver func(int))
	// Java: AbstractCommandLineRunner.isInTestMode()
	IsInTestMode() bool
	// Java: AbstractCommandLineRunner.getCommandLineConfig()
	GetCommandLineConfig() *AbstractCommandLineRunner_CommandLineConfig
	// Java: AbstractCommandLineRunner.createCompiler()
	CreateCompiler() *jscomp.Compiler
	// Java: AbstractCommandLineRunner.prepForBundleAndAppendTo(Appendable, CompilerInput, String)
	PrepForBundleAndAppendTo(out jutil.Appendable, input *jscomp.CompilerInput, content string)
	// Java: AbstractCommandLineRunner.appendRuntimeTo(Appendable)
	AppendRuntimeTo(out jutil.Appendable)
	// Java: AbstractCommandLineRunner.createOptions()
	CreateOptions() *jscomp.CompilerOptions
	// Java: AbstractCommandLineRunner.getDiagnosticGroups()
	GetDiagnosticGroups() *jscomp.DiagnosticGroups
	// Java: AbstractCommandLineRunner.addAllowlistWarningsGuard(CompilerOptions, File)
	AddAllowlistWarningsGuard(options *jscomp.CompilerOptions, allowlistFile *jutil.File)
	// Java: AbstractCommandLineRunner.setRunOptions(CompilerOptions)
	SetRunOptions(options *jscomp.CompilerOptions)
	// Java: AbstractCommandLineRunner.getVersionText()
	GetVersionText() string
	// Java: AbstractCommandLineRunner.parseJsonFilesFromInputStream()
	ParseJsonFilesFromInputStream() []*AbstractCommandLineRunner_JsonFileSpec
	// Java: AbstractCommandLineRunner.createInputs(List<FlagEntry<JsSourceType>>, List<JsonFileSpec>, boolean, List<JsChunkSpec>)
	CreateInputsWithJsonFilesAndAllowStdIn(
		files []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType],
		jsonFiles []*AbstractCommandLineRunner_JsonFileSpec,
		allowStdIn bool,
		jsChunkSpecs []*AbstractCommandLineRunner_JsChunkSpec) []*jscomp.SourceFile
	// Java: AbstractCommandLineRunner.checkChunkName(String)
	CheckChunkName(name string)
	// Java: AbstractCommandLineRunner.getChunkOutputFileName(JSChunk)
	GetChunkOutputFileName(m *jscomp.JSChunk) string
	// Java: AbstractCommandLineRunner.getChunkConformanceFileName(JSChunk)
	GetChunkConformanceFileName(m *jscomp.JSChunk) *string
	// Java: AbstractCommandLineRunner.writeChunkOutput(String, Appendable, LicenseTracker, JSChunk)
	WriteChunkOutput(fileName string, out jutil.Appendable, lt jscomp.CodePrinter_LicenseTracker, m *jscomp.JSChunk)
	// Java: AbstractCommandLineRunner.writeOutput(Appendable, Compiler, LicenseTracker, JSChunk, String, String, Function<String, String>, String)
	WriteOutput(
		out jutil.Appendable,
		compiler *jscomp.Compiler,
		licenseTracker jscomp.CodePrinter_LicenseTracker,
		chunk *jscomp.JSChunk,
		wrapper string,
		codePlaceholder string,
		escaper func(string) string,
		filename string)
	// Java: AbstractCommandLineRunner.writeOutput(Appendable, Compiler, String, String, String, Function<String, String>, String)
	WriteOutputOfCode(
		out jutil.Appendable,
		compiler *jscomp.Compiler,
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
	ProcessResults(result *jscomp.Result, chunks []*jscomp.JSChunk, options *jscomp.CompilerOptions) int
	// Java: AbstractCommandLineRunner.getJavascriptEscaper()
	GetJavascriptEscaper() func(string) string
	// Java: AbstractCommandLineRunner.outputSingleBinary(B)
	OutputSingleBinary(options *jscomp.CompilerOptions)
	// Java: AbstractCommandLineRunner.createJsonFile(B, String, Function<String, String>)
	CreateJsonFile(options *jscomp.CompilerOptions, outputMarker string, escaper func(string) string) *AbstractCommandLineRunner_JsonFileSpec
	// Java: AbstractCommandLineRunner.outputJsonStream()
	OutputJsonStream()
	// Java: AbstractCommandLineRunner.createExterns(CompilerOptions)
	CreateExterns(options *jscomp.CompilerOptions) []*jscomp.SourceFile
	// Java: AbstractCommandLineRunner.shouldGenerateMapPerChunk(B)
	ShouldGenerateMapPerChunk(options *jscomp.CompilerOptions) bool
	// Java: AbstractCommandLineRunner.expandSourceMapPath(B, JSChunk)
	ExpandSourceMapPath(options *jscomp.CompilerOptions, forChunk *jscomp.JSChunk) *string
	// Java: AbstractCommandLineRunner.filenameToOutputStream(String)
	FilenameToOutputStream(fileName *string) io.WriteCloser
	// Java: AbstractCommandLineRunner.createWriter(OutputStream, Charset)
	CreateWriter(stream io.WriteCloser, charset *jutil.Charset) AbstractCommandLineRunner_Writer
	// Java: AbstractCommandLineRunner.printChunkGraphJsonTo(Appendable)
	PrintChunkGraphJsonTo(out jutil.Appendable)
	// Java: AbstractCommandLineRunner.printChunkGraphManifestOrBundleTo(JSChunkGraph, Appendable, boolean)
	PrintChunkGraphManifestOrBundleTo(graph *jscomp.JSChunkGraph, out jutil.Appendable, isManifest bool)
	// Java: AbstractCommandLineRunner.printManifestTo(JSChunk, Appendable)
	PrintManifestTo(chunk *jscomp.JSChunk, out jutil.Appendable)
	// Java: AbstractCommandLineRunner.printBundleTo(JSChunk, Appendable)
	PrintBundleTo(chunk *jscomp.JSChunk, out jutil.Appendable)
}

// AbstractCommandLineRunnerBase holds AbstractCommandLineRunner's fields and
// base implementations (PORTING.md §4).
type AbstractCommandLineRunnerBase struct {
	self AbstractCommandLineRunner

	config *AbstractCommandLineRunner_CommandLineConfig

	in io.Reader

	defaultJsOutput *jutil.PrintStream

	err *jutil.PrintStream

	compiler *jscomp.Compiler

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

	externsSupplierForTesting func() []*jscomp.SourceFile // nil is null

	inputsSupplierForTesting func() []*jscomp.SourceFile // nil is null

	chunksSupplierForTesting func() []*jscomp.JSChunk // nil is null

	exitCodeReceiver func(int)

	rootRelativePathsMap *jutil.LinkedHashMap[string, string] // nil is null

	parsedChunkWrappers *jutil.LinkedHashMap[string, string] // nil is null

	parsedChunkOutputFiles *jutil.LinkedHashMap[string, string] // nil is null

	parsedChunkConformanceFiles *jutil.LinkedHashMap[string, string] // nil is null

	filesToStreamOut []*AbstractCommandLineRunner_JsonFileSpec
}

var (
	// Java: AbstractCommandLineRunner.OUTPUT_SAME_AS_INPUT_ERROR
	AbstractCommandLineRunner_OUTPUT_SAME_AS_INPUT_ERROR = jscomp.DiagnosticType_Error(
		"JSC_OUTPUT_SAME_AS_INPUT_ERROR", "Bad output file (already listed as input file): {0}")

	// Java: AbstractCommandLineRunner.COULD_NOT_SERIALIZE_AST
	AbstractCommandLineRunner_COULD_NOT_SERIALIZE_AST = jscomp.DiagnosticType_Error("JSC_COULD_NOT_SERIALIZE_AST", "Could not serialize ast to: {0}")

	// Java: AbstractCommandLineRunner.COULD_NOT_DESERIALIZE_AST
	AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST = jscomp.DiagnosticType_Error("JSC_COULD_NOT_DESERIALIZE_AST", "Could not deserialize ast from: {0}")

	// Java: AbstractCommandLineRunner.NO_TREE_GENERATED_ERROR
	AbstractCommandLineRunner_NO_TREE_GENERATED_ERROR = jscomp.DiagnosticType_Error(
		"JSC_NO_TREE_GENERATED_ERROR", "Code contains errors. No tree was generated.")
	// Java: AbstractCommandLineRunner.INVALID_CHUNK_SOURCEMAP_PATTERN
	AbstractCommandLineRunner_INVALID_CHUNK_SOURCEMAP_PATTERN = jscomp.DiagnosticType_Error(
		"JSC_INVALID_CHUNK_SOURCEMAP_PATTERN",
		"When using --chunk or --module flags, the --create_source_map flag must contain "+
			"%outname% in the value.")

	// Java: AbstractCommandLineRunner.EXPECTED_DIAGNOSTIC_NOT_FOUND
	AbstractCommandLineRunner_EXPECTED_DIAGNOSTIC_NOT_FOUND = jscomp.DiagnosticType_Error(
		"JSC_EXPECTED_DIAGNOSTIC_NOT_FOUND", "Expected diagnostic not found: {0}")

	// Java: AbstractCommandLineRunner.AMBIGUOUS_EXPECTATION
	AbstractCommandLineRunner_AMBIGUOUS_EXPECTATION = jscomp.DiagnosticType_Error(
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
	externsSupplier func() []*jscomp.SourceFile,
	inputsSupplier func() []*jscomp.SourceFile,
	chunksSupplier func() []*jscomp.JSChunk,
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
	return r.config.jsonStreamMode == jscomp.CompilerOptions_JsonStreamMode_OUT ||
		r.config.jsonStreamMode == jscomp.CompilerOptions_JsonStreamMode_BOTH
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
func (r *AbstractCommandLineRunnerBase) GetDiagnosticGroups() *jscomp.DiagnosticGroups {
	if r.compiler == nil {
		return jscomp.NewDiagnosticGroups()
	}
	return r.compiler.GetDiagnosticGroups()
}

// Java: AbstractCommandLineRunner.setWarningGuardOptions(CompilerOptions, ArrayList<FlagEntry<CheckLevel>>, DiagnosticGroups)
func AbstractCommandLineRunner_SetWarningGuardOptions(
	options *jscomp.CompilerOptions,
	warningGuards []*AbstractCommandLineRunner_FlagEntry[jscomp.CheckLevel],
	diagnosticGroups *jscomp.DiagnosticGroups) {
	if warningGuards != nil {
		groupNames := jscomp.DiagnosticGroups_GetRegisteredGroups() // Java: .keySet()
		for _, entry := range warningGuards {
			if "*" == entry.value {
				for groupName := range groupNames.Keys() {
					if !jscomp.DiagnosticGroups_wildcardExcludedGroups.Contains(groupName) {
						diagnosticGroups.SetWarningLevel(options, groupName, entry.flag)
					}
				}
			} else if groupNames.ContainsKey(entry.value) {
				diagnosticGroups.SetWarningLevel(options, entry.value, entry.flag)
			} else {
				panic(jscomp.NewFlagUsageException("Unknown diagnostic group: '" + entry.value + "'"))
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
func (r *AbstractCommandLineRunnerBase) SetRunOptions(options *jscomp.CompilerOptions) {
	config := r.config
	diagnosticGroups := r.self.GetDiagnosticGroups()

	if config.ShouldSaveAfterStage1() || config.ShouldContinueCompilation() {
		if options.IsChecksOnly() {
			panic(jscomp.NewFlagUsageException(
				"checks_only mode is incompatible with multi-stage compilation"))
		}
		if options.GetExternExportsPath() != nil {
			panic(jscomp.NewFlagUsageException(
				"generating externs from exports is incompatible with multi-stage compilation"))
		}
	}

	AbstractCommandLineRunner_SetWarningGuardOptions(options, config.warningGuards, diagnosticGroups)

	if config.warningsAllowFile != "" {
		r.self.AddAllowlistWarningsGuard(options, jutil.NewFile(config.warningsAllowFile))
	}

	if len(config.hideWarningsFor) != 0 {
		options.AddWarningsGuard(
			jscomp.NewShowByPathWarningsGuardWithPathsAndShowType(
				append([]string{}, config.hideWarningsFor...),
				jscomp.ShowByPathWarningsGuard_ShowType_EXCLUDE))
	}

	define := append([]string{}, config.define...)
	if config.browserFeaturesetYear != 0 {
		jutil.Try(func() {
			options.SetBrowserFeaturesetYear(config.browserFeaturesetYear)
		}, func(e jutil.IIllegalStateException) {
			panic(jscomp.NewFlagUsageException(e.GetMessage()))
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
			panic(jscomp.NewFlagUsageException(
				"skip_normal_outputs and js_output_file cannot be used together."))
		}
	}

	if config.skipNormalOutputs && config.printAst {
		panic(jscomp.NewFlagUsageException("skip_normal_outputs and print_ast cannot be used together."))
	}

	if config.skipNormalOutputs && config.printTree {
		panic(jscomp.NewFlagUsageException("skip_normal_outputs and print_tree cannot be used together."))
	}

	if config.skipNormalOutputs && config.printTreeJson {
		panic(jscomp.NewFlagUsageException(
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

	inputSourceMaps := jutil.NewImmutableMapBuilder[string, *jscomp.SourceMapInput]()
	for key, value := range config.sourceMapInputFiles.All() {
		sourceMap := jscomp.SourceFile_Builder_().WithKind(rhino.StaticSourceFile_SourceKind_NON_CODE).WithPath(value).Build()
		inputSourceMaps.Put(key, jscomp.NewSourceMapInput(sourceMap))
	}
	options.SetInputSourceMaps(inputSourceMaps.BuildOrThrow())

	if config.variableMapInputFile != "" {
		options.SetInputVariableMap(jscomp.VariableMap_Load(config.variableMapInputFile))
	}

	if config.propertyMapInputFile != "" {
		options.SetInputPropertyMap(jscomp.VariableMap_Load(config.propertyMapInputFile))
	}

	if len(config.outputManifests) != 0 {
		uniqueNames := jutil.NewLinkedHashSet[string]()
		for _, filename := range config.outputManifests {
			if !uniqueNames.Add(filename) {
				panic(jscomp.NewFlagUsageException(
					"output_manifest flags specify duplicate file names: " + filename))
			}
		}
	}

	if len(config.outputBundles) != 0 {
		uniqueNames := jutil.NewLinkedHashSet[string]()
		for _, filename := range config.outputBundles {
			if !uniqueNames.Add(filename) {
				panic(jscomp.NewFlagUsageException(
					"output_bundle flags specify duplicate file names: " + filename))
			}
		}
	}

	options.SetProcessCommonJSModules(config.processCommonJSModules)
	options.SetModuleRoots(config.moduleRoots)
	options.SetAngularPass(config.angularPass)

	if config.jsonWarningsFile != "" {
		options.AddReportGenerator(
			jscomp.NewJsonErrorReportGenerator(jutil.NewPrintStreamFromFile(jutil.NewFile(config.jsonWarningsFile)), r.compiler))
	}

	if config.errorFormat == AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_JSON {
		errorGenerator := jscomp.NewJsonErrorReportGenerator(r.GetErrorPrintStream(), r.compiler)
		r.compiler.SetErrorManager(jscomp.NewSortingErrorManager([]jscomp.SortingErrorManager_ErrorReportGenerator{errorGenerator}))
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
func (r *AbstractCommandLineRunnerBase) GetCompiler() *jscomp.Compiler {
	return r.compiler
}

// AbstractCommandLineRunner_GetBuiltinExterns returns a mutable list.
//
// Java: AbstractCommandLineRunner.getBuiltinExterns(CompilerOptions.Environment)
func AbstractCommandLineRunner_GetBuiltinExterns(env jscomp.CompilerOptions_Environment) []*jscomp.SourceFile {
	input := abstractCommandLineRunner_getExternsInput()
	zipReader, err := zip.NewReader(bytes.NewReader(input), int64(len(input)))
	if err != nil {
		panic(jutil.NewZipException(err.Error()))
	}
	envPrefix := jutil.Ascii_ToLowerCase(env.String()) + "/"
	mapFromExternsZip := jutil.NewLinkedHashMap[string, *jscomp.SourceFile]()
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

		mapFromExternsZip.Put(
			filename,
			// Give the files an odd prefix, so that they do not conflict
			// with the user's files.
			//
			// Go-only: performance — the content is read from the zip
			// (embedded in the binary) when needed, and again if needed
			// after parsing dropped it (Java: withContentFromStream, which
			// keeps it for the whole compilation).
			jscomp.SourceFile_Builder_().
				WithPath("externs.zip//"+filename).
				WithReloadableContent(func() string {
					entryStream, err := entry.Open()
					if err != nil {
						panic(jutil.NewZipException(err.Error()))
					}
					defer entryStream.Close()
					return jutil.CharStreams_ToString(jutil.NewInputStreamReader(entryStream, jutil.StandardCharsets_UTF_8))
				}).
				Build())
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
		}, func(e jscomp.IFlagUsageException) {
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
		jsonFiles = append(jsonFiles, parseJsonFileSpec(reader))
	}
	reader.EndArray()
	return jsonFiles
}

func parseJsonFileSpec(reader *jutil.JsonReader) (result *AbstractCommandLineRunner_JsonFileSpec) {
	isEmpty := true
	oldLenient := reader.IsLenient()
	reader.SetLenient(true)
	defer reader.SetLenient(oldLenient)
	jutil.Try(func() {
		tok := reader.Peek()
		isEmpty = false
		if tok == jutil.JsonToken_NULL {
			reader.NextNull()
			return
		}
		spec := newAbstractCommandLineRunner_JsonFileSpecDefault()
		reader.BeginObject()
		for reader.HasNext() {
			switch reader.NextName() {
			case "src":
				spec.Src = readJsonNullableString(reader)
			case "path":
				spec.Path = readJsonNullableString(reader)
			case "source_map", "sourceMap":
				spec.SourceMap = readJsonNullableString(reader)
			case "webpack_id", "webpackId":
				spec.WebpackId = readJsonNullableString(reader)
			default:
				reader.SkipValue()
			}
		}
		reader.EndObject()
		result = spec
	}, func(t jutil.Throwable) {
		switch e := t.(type) {
		case jutil.IEOFException:
			if isEmpty {
				result = nil
				return
			}
			panic(jutil.NewJsonSyntaxExceptionWithCause(e))
		case jutil.IIllegalStateException:
			panic(jutil.NewJsonSyntaxExceptionWithCause(e))
		case jutil.IIOException:
			panic(jutil.NewJsonSyntaxExceptionWithCause(e))
		default:
			panic(t)
		}
	})
	return result
}

func readJsonNullableString(reader *jutil.JsonReader) *string {
	switch reader.Peek() {
	case jutil.JsonToken_NULL:
		reader.NextNull()
		return nil
	case jutil.JsonToken_BOOLEAN:
		if reader.NextBoolean() {
			s := "true"
			return &s
		}
		s := "false"
		return &s
	default:
		s := reader.NextString()
		return &s
	}
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
	jsChunkSpecs []*AbstractCommandLineRunner_JsChunkSpec) []*jscomp.SourceFile {
	return r.self.CreateInputsWithJsonFilesAndAllowStdIn(files /* jsonFiles= */, nil, allowStdIn, jsChunkSpecs)
}

// CreateInputsWithJsonFiles creates inputs from a list of source files and
// json files.
//
// Java: AbstractCommandLineRunner.createInputs(List<FlagEntry<JsSourceType>>, List<JsonFileSpec>, List<JsChunkSpec>)
func (r *AbstractCommandLineRunnerBase) CreateInputsWithJsonFiles(
	files []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType],
	jsonFiles []*AbstractCommandLineRunner_JsonFileSpec,
	jsChunkSpecs []*AbstractCommandLineRunner_JsChunkSpec) []*jscomp.SourceFile {
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
	jsChunkSpecs []*AbstractCommandLineRunner_JsChunkSpec) []*jscomp.SourceFile {
	inputs := make([]*jscomp.SourceFile, 0, len(files))
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
				panic(jscomp.NewFlagUsageException("Can't use TypedASTs with --zip."))
			}

			if "-" != filename {
				newFiles := jscomp.SourceFile_FromZipFile(filename, r.inputCharset)
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
			inputs = append(inputs,
				jscomp.SourceFile_Builder_().
					WithKind(kind).
					WithCharset(r.inputCharset).
					WithPath(filename).
					Build())
		} else {
			if r.config.typedAstListInputFilename != nil {
				panic(jscomp.NewFlagUsageException("Can't use TypedASTs with stdin."))
			}
			if !r.config.defaultToStdin {
				panic(jscomp.NewFlagUsageException("Can't specify stdin."))
			}
			if usingStdin {
				panic(jscomp.NewFlagUsageException("Can't specify stdin twice."))
			}

			if len(r.config.outputManifests) != 0 {
				panic(jscomp.NewFlagUsageException(
					"Manifest files cannot be generated when the input is from stdin."))
			}
			if len(r.config.outputBundles) != 0 {
				panic(jscomp.NewFlagUsageException(
					"Bundle files cannot be generated when the input is from stdin."))
			}

			r.err.Println(AbstractCommandLineRunner_WAITING_FOR_INPUT_WARNING)
			inputs = append(inputs,
				jscomp.SourceFile_Builder_().
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
			inputs = append(inputs, jscomp.SourceFile_FromCode(*jsonFile.GetPath(), *jsonFile.GetSrc()))
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
	hasChunkSpecs bool) []*jscomp.JSError {
	var errors []*jscomp.JSError

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
				panic(jscomp.NewFlagUsageException("--ijs is incompatible with --chunk or --module."))
			}
			absoluteName := file.value
			if !strings.HasSuffix(absoluteName, ".i.js") {
				errors = append(errors, jscomp.JSError_Make(jscomp.IjsErrors_BAD_IJS_FILE_NAME, absoluteName))
				i++
				continue
			}
			relativeName := r.GetModuleRootRelativeName(absoluteName, moduleRoots)
			relativeNonIjsName := relativeName[:len(relativeName)-len(".i.js")]
			if relativeToAbsoluteName.ContainsKey(relativeNonIjsName) {
				errors = append(errors,
					jscomp.JSError_Make(
						jscomp.IjsErrors_CONFLICTING_IJS_FILE,
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
	moduleRoots []string) []*jscomp.SourceFile {
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

	var result []*jscomp.SourceFile
	jutil.Try(func() {
		if jsonFiles != nil {
			result = r.CreateInputsWithJsonFiles(files, jsonFiles, jsChunkSpecs)
		} else {
			result = r.CreateInputs(files, true, jsChunkSpecs)
		}
	}, func(e jscomp.IFlagUsageException) {
		panic(jscomp.NewFlagUsageException("Bad --js flag. " + e.GetMessage()))
	})
	return result
}

// CreateExternInputs creates JS extern inputs from a list of files.
//
// Java: AbstractCommandLineRunner.createExternInputs(List<String>)
func (r *AbstractCommandLineRunnerBase) CreateExternInputs(files []string) []*jscomp.SourceFile {
	var externFiles []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]
	for _, file := range files {
		externFiles = append(externFiles, NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_EXTERN, file))
	}
	var result []*jscomp.SourceFile
	jutil.Try(func() {
		result = r.CreateInputs(externFiles, false, []*AbstractCommandLineRunner_JsChunkSpec{})
	}, func(e jscomp.IFlagUsageException) {
		panic(jscomp.NewFlagUsageException("Bad --externs flag. " + e.GetMessage()))
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
func AbstractCommandLineRunner_CreateJsChunks(specs []*AbstractCommandLineRunner_JsChunkSpec, inputs []*jscomp.CompilerInput) []*jscomp.JSChunk {
	jutil.CheckState(specs != nil)
	jutil.CheckState(len(specs) != 0)
	jutil.CheckState(inputs != nil)

	result := make([]*jscomp.JSChunk, 0, len(specs))
	chunksByName := make(map[string]*jscomp.JSChunk, len(specs))
	numJsFilesExpected := 0
	minJsFilesRequired := 0
	for _, spec := range specs {
		if _, exists := chunksByName[spec.name]; exists {
			panic(jscomp.NewFlagUsageException("Duplicate chunk name: " + spec.name))
		}
		chunk := jscomp.NewJSChunk(spec.name)

		for _, dep := range spec.deps {
			other := chunksByName[dep]
			if other == nil {
				panic(jscomp.NewFlagUsageException(
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

		chunksByName[spec.name] = chunk
		result = append(result, chunk)
	}

	totalNumJsFiles := len(inputs)

	if numJsFilesExpected >= 0 || minJsFilesRequired > totalNumJsFiles {
		if minJsFilesRequired > totalNumJsFiles {
			numJsFilesExpected = minJsFilesRequired
		}

		if numJsFilesExpected > totalNumJsFiles {
			panic(jscomp.NewFlagUsageException(
				"Not enough JS files specified. Expected " +
					strconv.Itoa(numJsFilesExpected) +
					" but found " +
					strconv.Itoa(totalNumJsFiles)))
		} else if numJsFilesExpected < totalNumJsFiles {
			panic(jscomp.NewFlagUsageException(
				"Too many JS files specified. Expected " +
					strconv.Itoa(numJsFilesExpected) +
					" but found " +
					strconv.Itoa(totalNumJsFiles)))
		}
	}

	// Allocate source files to chunks in reverse order so that the first chunk
	// (presumably the base chunk) can have a size of 'auto' (-1).
	numJsFilesLeft := totalNumJsFiles
	for i := len(specs) - 1; i >= 0; i-- {
		numJsFiles := specs[i].numJsFiles
		chunk := result[i]

		// Check if the first chunk specified 'auto' for the number of files
		if i == 0 && numJsFiles == -1 {
			numJsFiles = numJsFilesLeft
		}

		chunkFiles := inputs[numJsFilesLeft-numJsFiles : numJsFilesLeft]
		for _, input := range chunkFiles {
			chunk.Add(input)
		}
		numJsFilesLeft -= numJsFiles
	}

	return result
}

// CheckChunkName validates the chunk name. Can be overridden by subclasses.
//
// Java: AbstractCommandLineRunner.checkChunkName(String)
func (r *AbstractCommandLineRunnerBase) CheckChunkName(name string) {
	if !rhino.TokenStream_IsJSIdentifier(name) {
		panic(jscomp.NewFlagUsageException("Invalid chunk name: '" + name + "'"))
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

	jscomp.Compiler_SetLoggingLevel(jutil.Level_Parse(config.loggingLevel))

	r.compiler = r.self.CreateCompiler()
	if len(config.expectedDiagnostics) != 0 {
		r.compiler.SetErrorManager(
			NewVerifyingErrorManager(r.compiler.GetErrorManager(), config.expectedDiagnostics))
	}
	options := r.self.CreateOptions()
	r.self.SetRunOptions(options)

	typedAstListInputFilename := config.typedAstListInputFilename

	externs := r.self.CreateExterns(options)
	var chunks []*jscomp.JSChunk

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

	if config.jsonStreamMode == jscomp.CompilerOptions_JsonStreamMode_IN ||
		config.jsonStreamMode == jscomp.CompilerOptions_JsonStreamMode_BOTH {
		jsonFiles = r.self.ParseJsonFilesFromInputStream()
		if jsonFiles == nil {
			jsonFiles = []*AbstractCommandLineRunner_JsonFileSpec{}
		}

		inputSourceMaps := jutil.NewImmutableMapBuilder[string, *jscomp.SourceMapInput]()
		inputPathByWebpackId := jutil.NewImmutableMapBuilder[string, string]()

		foundJsonInputSourceMap := false
		for _, jsonFile := range jsonFiles {
			if jsonFile.GetSourceMap() != nil && jutil.StringLength(*jsonFile.GetSourceMap()) > 0 {
				sourceMapPath := jutil.ToJavaString(jsonFile.GetPath()) + ".map"
				sourceMap := jscomp.SourceFile_FromCode(sourceMapPath, *jsonFile.GetSourceMap())
				inputSourceMaps.Put(*jsonFile.GetPath(), jscomp.NewSourceMapInput(sourceMap))
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
			if jscomp.CompilerOptions_JsonStreamMode_IN == config.jsonStreamMode ||
				jscomp.CompilerOptions_JsonStreamMode_NONE == config.jsonStreamMode {
				for _, m := range jsChunkSpecs {
					r.self.CheckChunkName(m.GetName())
				}
			}
			inputs := make([]*jscomp.CompilerInput, 0, len(sources))
			for _, source := range sources {
				inputs = append(inputs, jscomp.NewCompilerInput(source))
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

	var result *jscomp.Result
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
				chunks = append([]*jscomp.JSChunk{}, r.compiler.GetChunks()...)
			}
		}()
		result = r.compiler.GetResult()
	}

	if createCommonJsModules {
		// For CommonJS modules, construct chunks from actual inputs.
		chunks = append([]*jscomp.JSChunk{}, r.compiler.GetChunks()...)
		for _, c := range chunks {
			outputFileNames = append(outputFileNames, r.self.GetChunkOutputFileName(c))
		}
	}

	for _, outputFileName := range outputFileNames {
		if r.compiler.GetSourceFileByName(outputFileName) != nil {
			r.compiler.Report(jscomp.JSError_Make(AbstractCommandLineRunner_OUTPUT_SAME_AS_INPUT_ERROR, outputFileName))
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
		r.compiler.Stage2Passes(config.stage2SegmentOfCompilationToRun)
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
	externs []*jscomp.SourceFile,
	sources []*jscomp.SourceFile,
	options *jscomp.CompilerOptions,
	filename string) {
	jutil.Try(func() {
		typedAstListStream := jutil.NewGZIPInputStreamWithSize(jutil.NewFileInputStream(filename), abstractCommandLineRunner_GZIPPED_TYPEDAST_BUFFER_SIZE)
		defer typedAstListStream.Close()
		r.compiler.InitWithTypedAstFilesystem(externs, sources, options, typedAstListStream)
	}, func(e jutil.IIOException) {
		r.compiler.Report(jscomp.JSError_Make(AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST, filename))
	})
}

// Java: AbstractCommandLineRunner.initChunksWithTypedAstFilesystem(List<SourceFile>, List<JSChunk>, CompilerOptions, String)
func (r *AbstractCommandLineRunnerBase) InitChunksWithTypedAstFilesystem(
	externs []*jscomp.SourceFile, chunks []*jscomp.JSChunk, options *jscomp.CompilerOptions, filename string) {
	jutil.Try(func() {
		typedAstListStream := jutil.NewGZIPInputStreamWithSize(jutil.NewFileInputStream(filename), abstractCommandLineRunner_GZIPPED_TYPEDAST_BUFFER_SIZE)
		defer typedAstListStream.Close()
		r.compiler.InitChunksWithTypedAstFilesystem(externs, chunks, options, typedAstListStream)
	}, func(e jutil.IIOException) {
		r.compiler.Report(jscomp.JSError_Make(AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST, filename))
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
		r.compiler.Report(jscomp.JSError_Make(AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST, filename))
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
		r.compiler.Report(jscomp.JSError_Make(AbstractCommandLineRunner_COULD_NOT_SERIALIZE_AST, filename))
	})
}

// ProcessResults processes the results of the compile job, and returns an
// error code.
//
// Java: AbstractCommandLineRunner.processResults(Result, List<JSChunk>, B)
func (r *AbstractCommandLineRunnerBase) ProcessResults(result *jscomp.Result, chunks []*jscomp.JSChunk, options *jscomp.CompilerOptions) int {
	config := r.config
	if config.printAst {
		if r.compiler.GetRoot().IsNil() {
			return 1
		} else {
			jsOutput := r.CreateDefaultOutput()
			cfg := r.compiler.ComputeCFG()
			jscomp.DotFormatter_AppendDot(r.compiler.GetRoot().GetLastChild(), cfg, jsOutput)
			jsOutput.AppendCharacter('\n')
			abstractCommandLineRunner_closeAppendable(jsOutput)
			return 0
		}
	}

	if config.printTree {
		if r.compiler.GetRoot().IsNil() {
			r.compiler.Report(jscomp.JSError_Make(AbstractCommandLineRunner_NO_TREE_GENERATED_ERROR))
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
		if r.compiler.GetRoot().IsNil() {
			r.compiler.Report(jscomp.JSError_Make(AbstractCommandLineRunner_NO_TREE_GENERATED_ERROR))
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
	} else if options.GetOutputJs() != jscomp.CompilerOptions_OutputJs_NONE && result.Success {
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
				r.compiler.Report(jscomp.JSError_Make(error))
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
				if result.ExternExport != nil {
					externExport = *result.ExternExport
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
func abstractCommandLineRunner_getExitStatusForResult(result *jscomp.Result) int {
	// return 0 if no errors, the error count otherwise
	return min(len(result.Errors), 0x7f)
}

// GetInputCharset queries the flag for the input charset, and returns a
// Charset object representing the selection. It throws FlagUsageException
// if flag is not a valid Charset name.
//
// Java: AbstractCommandLineRunner.getInputCharset()
func (r *AbstractCommandLineRunnerBase) GetInputCharset() *jutil.Charset {
	if r.config.charset != "" {
		if !jutil.Charset_IsSupported(r.config.charset) {
			panic(jscomp.NewFlagUsageException(r.config.charset + " is not a valid charset name."))
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
			panic(jscomp.NewFlagUsageException(r.config.charset + " is not a valid charset name."))
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
			panic(jscomp.NewFlagUsageException(r.config.charset + " is not a valid charset name."))
		}
		return jutil.Charset_ForName(r.config.charset)
	}
	return jutil.StandardCharsets_UTF_8
}

// Java: AbstractCommandLineRunner.createExterns(CompilerOptions)
func (r *AbstractCommandLineRunnerBase) CreateExterns(options *jscomp.CompilerOptions) []*jscomp.SourceFile {
	if r.self.IsInTestMode() {
		return r.externsSupplierForTesting()
	}
	return r.CreateExternInputs(r.config.externs)
}

// AbstractCommandLineRunner_CreateDefineReplacements creates a map of
// constant names to constant values from a textual description of the map.
//
// definitions is a list of overriding definitions for defines in the form
// <name>[=<val>], where <val> is a number, boolean, or single-quoted string
// without single quotes.
//
// Java: AbstractCommandLineRunner.createDefineReplacements(List<String>, CompilerOptions)
func AbstractCommandLineRunner_CreateDefineReplacements(definitions []string, options *jscomp.CompilerOptions) {
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
