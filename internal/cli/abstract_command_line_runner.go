// Port of com/google/javascript/jscomp/AbstractCommandLineRunner.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/VerifyingErrorManager.java @ closure-compiler v20261021

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

/*
 * Copyright 2026 The Closure Compiler Authors.
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
	"io"
	"slices"
	"strconv"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
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
	ExpandSourceMapPath(options *jscomp.CompilerOptions, forChunk *jscomp.JSChunk) string
	// Java: AbstractCommandLineRunner.filenameToOutputStream(String)
	FilenameToOutputStream(fileName string) io.WriteCloser
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

	rootRelativePathsMap map[string]string // nil is null

	parsedChunkWrappers map[string]string // nil is null

	parsedChunkOutputFiles map[string]string // nil is null

	parsedChunkConformanceFiles map[string]string // nil is null

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
		if err := setBrowserFeaturesetYear(options, config.browserFeaturesetYear); err != nil {
			panic(err)
		}
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

	if config.jsOutputFile != "" {
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

	if config.createSourceMap != "" {
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

	if len(config.sourceMapInputFiles) == 0 {
		options.SetInputSourceMaps(jutil.NewLinkedHashMap[string, *jscomp.SourceMapInput]())
	} else {
		entries := make([]kvEntry[*jscomp.SourceMapInput], len(config.sourceMapInputFiles))
		for i, entry := range config.sourceMapInputFiles {
			sourceMap := jscomp.SourceFile_Builder_().WithKind(rhino.StaticSourceFile_SourceKind_NON_CODE).WithPath(entry.Value).Build()
			entries[i] = kvEntry[*jscomp.SourceMapInput]{key: entry.Key, value: jscomp.NewSourceMapInput(sourceMap)}
		}
		options.SetInputSourceMaps(buildUniqueLinkedMap(entries))
	}

	if config.variableMapInputFile != "" {
		options.SetInputVariableMap(jscomp.VariableMap_Load(config.variableMapInputFile))
	}

	if config.propertyMapInputFile != "" {
		options.SetInputPropertyMap(jscomp.VariableMap_Load(config.propertyMapInputFile))
	}

	if len(config.outputManifests) != 0 {
		uniqueNames := make(map[string]struct{}, len(config.outputManifests))
		for _, filename := range config.outputManifests {
			if _, exists := uniqueNames[filename]; exists {
				panic(jscomp.NewFlagUsageException(
					"output_manifest flags specify duplicate file names: " + filename))
			}
			uniqueNames[filename] = struct{}{}
		}
	}

	if len(config.outputBundles) != 0 {
		uniqueNames := make(map[string]struct{}, len(config.outputBundles))
		for _, filename := range config.outputBundles {
			if _, exists := uniqueNames[filename]; exists {
				panic(jscomp.NewFlagUsageException(
					"output_bundle flags specify duplicate file names: " + filename))
			}
			uniqueNames[filename] = struct{}{}
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

func setBrowserFeaturesetYear(options *jscomp.CompilerOptions, year int) *jscomp.FlagUsageException {
	if year != 2012 && (year < 2018 || year > 2026) {
		return jscomp.NewFlagUsageException(
			"Illegal browser_featureset_year=" + strconv.Itoa(year) + ". We support values 2012, or 2018..2026 only")
	}
	options.SetBrowserFeaturesetYear(year)
	return nil
}

// Java: AbstractCommandLineRunner.getCompiler()
func (r *AbstractCommandLineRunnerBase) GetCompiler() *jscomp.Compiler {
	return r.compiler
}

// Run runs the Compiler and calls System.exit() with the exit status of the
// compiler.
//
// Java: AbstractCommandLineRunner.run()
func (r *AbstractCommandLineRunnerBase) Run() {
	result, err := r.tryDoRun()
	if err != nil {
		if e, ok := err.(jscomp.IFlagUsageException); ok {
			r.err.Println(e.GetMessage())
			result = -1
		} else {
			jutil.Throwable_PrintStackTrace(err, r.err)
			result = -2
		}
	}

	r.exitCodeReceiver(result)
}

func (r *AbstractCommandLineRunnerBase) tryDoRun() (result int, thrown jutil.Throwable) {
	defer func() {
		if rec := recover(); rec != nil {
			thrown = jutil.AsThrowable(rec)
		}
	}()
	return r.self.DoRun(), nil
}

// GetErrorPrintStream returns the PrintStream for writing errors associated
// with this AbstractCommandLineRunner.
//
// Java: AbstractCommandLineRunner.getErrorPrintStream()
func (r *AbstractCommandLineRunnerBase) GetErrorPrintStream() *jutil.PrintStream {
	return r.err
}

// DoRun parses command-line arguments and runs the compiler. It returns the
// system exit status.
//
// Java: AbstractCommandLineRunner.doRun()
func (r *AbstractCommandLineRunnerBase) DoRun() int {
	config := r.config
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
	if len(config.chunk) > 0 {
		jsChunkSpecs = make([]*AbstractCommandLineRunner_JsChunkSpec, len(config.chunk))
		for i, c := range config.chunk {
			jsChunkSpecs[i] = AbstractCommandLineRunner_JsChunkSpec_Create(c, i == 0)
		}
	}
	var jsonFiles []*AbstractCommandLineRunner_JsonFileSpec

	if config.jsonStreamMode == jscomp.CompilerOptions_JsonStreamMode_IN ||
		config.jsonStreamMode == jscomp.CompilerOptions_JsonStreamMode_BOTH {
		jsonFiles = r.self.ParseJsonFilesFromInputStream()
		if jsonFiles == nil {
			jsonFiles = []*AbstractCommandLineRunner_JsonFileSpec{}
		}

		var smEntries []kvEntry[*jscomp.SourceMapInput]
		var webpackEntries []kvEntry[string]

		for _, jsonFile := range jsonFiles {
			if sm := jsonFile.GetSourceMap(); sm != nil && *sm != "" {
				pathStr := "null"
				if p := jsonFile.GetPath(); p != nil {
					pathStr = *p
				}
				sourceMapPath := pathStr + ".map"
				sourceMap := jscomp.SourceFile_FromCode(sourceMapPath, *sm)
				smEntries = append(smEntries, kvEntry[*jscomp.SourceMapInput]{
					key:   *jsonFile.GetPath(),
					value: jscomp.NewSourceMapInput(sourceMap),
				})
			}
			if wid := jsonFile.GetWebpackId(); wid != nil {
				webpackEntries = append(webpackEntries, kvEntry[string]{
					key:   *wid,
					value: *jsonFile.GetPath(),
				})
			}
		}

		if len(smEntries) > 0 {
			for k, v := range options.GetInputSourceMaps().All() {
				smEntries = append(smEntries, kvEntry[*jscomp.SourceMapInput]{key: k, value: v})
			}
			options.SetInputSourceMaps(buildUniqueLinkedMap(smEntries))
		}

		r.compiler.InitWebpackMap(buildUniqueLinkedMap(webpackEntries))
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
		// TODO(bradfordcsmith): Should we be ignoring possible init/initChunks() errors here?
		r.compiler.OrderInputsWithLargeStack()
		result = nil
	} else if r.compiler.HasErrors() {
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
			r.PerformCompilation()

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
	return exitStatus
}

// Java: AbstractCommandLineRunner.performCompilation()
func (r *AbstractCommandLineRunnerBase) PerformCompilation() {
	// Parse, restore from a save file, or initialize from a TypedAST list.
	r.InitializeStateBeforeCompilation()

	// Do the interesting work - checks, optimizations, etc.
	r.RunCompilerPasses()

	// Save state unless there was an exception or error - we won't go on to the next stage if there
	// was an error, and saving state itself might throw an exception if the compiler is in a bad
	// state.
	if !r.compiler.HasErrors() {
		r.SaveState()
	}
	// Perform post-compilation tasks even if compiler.hasErrors()
	r.compiler.PerformPostCompilationTasks()
}

// Java: AbstractCommandLineRunner.runCompilerPasses()
func (r *AbstractCommandLineRunnerBase) RunCompilerPasses() {
	config := r.config
	runStage1 := false
	runStage2 := false
	runStage3 := false
	instrumentForCoverage := false
	if r.compiler.GetOptions().GetInstrumentForCoverageOnly() {
		instrumentForCoverage = true
	} else if config.ShouldSaveAfterStage1() {
		runStage1 = true
	} else if config.ShouldRestoreTypedAstsPerformStage2AndSave() {
		runStage2 = true
	} else if config.ShouldRestoreTypedAstsPerformStages2And3() {
		runStage2 = true
		runStage3 = true
	} else if config.ShouldRestoreAndPerformStage2AndSave() {
		runStage2 = true
	} else if config.ShouldRestoreAndPerformStages2And3() {
		runStage2 = true
		runStage3 = true
	} else if config.ShouldRestoreAndPerformStage3() {
		runStage3 = true
	} else {
		// This is the code path taken when "building" a library by just checking it for errors
		// and generating an .ijs file and also when doing a full compilation.
		runStage1 = true
		runStage2 = true
		runStage3 = true
	}

	if r.compiler.HasErrors() {
		return
	}
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

type kvEntry[V any] struct {
	key   string
	value V
}

func buildUniqueLinkedMap[V any](entries []kvEntry[V]) *jutil.LinkedHashMap[string, V] {
	n := len(entries)
	out := jutil.NewLinkedHashMapWithCapacity[string, V](n)
	if n == 0 {
		return out
	}
	for _, e := range entries {
		if jutil.IsNil(e.value) {
			panic(jutil.NewNullPointerExceptionMsg("null value in entry: " + e.key + "=null"))
		}
	}
	hasDup := false
	if n > 8 {
		seen := make(map[string]struct{}, n)
		for _, e := range entries {
			if _, ok := seen[e.key]; ok {
				hasDup = true
				break
			}
			seen[e.key] = struct{}{}
		}
	} else {
		hasDup = true
	}
	if hasDup {
		for i := n - 1; i >= 0; i-- {
			k := entries[i].key
			for j := i + 1; j < n; j++ {
				if entries[j].key == k {
					panic(jutil.NewIllegalArgumentExceptionMsg("Multiple entries with same key: " +
						k + "=" + jutil.ToJavaString(entries[j].value) + " and " +
						k + "=" + jutil.ToJavaString(entries[i].value)))
				}
			}
		}
	}
	for _, e := range entries {
		out.Put(e.key, e.value)
	}
	return out
}

// VerifyingErrorManager is an ErrorManager that verifies expected
// diagnostics and reports missing or unexpected ones.
//
// It wraps a delegate ErrorManager. Expected diagnostics (passed as regexes)
// are swallowed if matched. Unexpected diagnostics are passed to the
// delegate. At the end of compilation (when generateReport is called), any
// unmatched expectations are reported as errors to the delegate.
type VerifyingErrorManager struct {
	*jscomp.ThreadSafeDelegatingErrorManager
	mu                    sync.Mutex
	unmatchedExpectations []*jutil.Pattern
	formatter             *jscomp.LightweightMessageFormatter
}

var _ jscomp.ErrorManager = (*VerifyingErrorManager)(nil)

// Java: VerifyingErrorManager(ErrorManager, List<String>)
func NewVerifyingErrorManager(delegate jscomp.ErrorManager, expectedDiagnostics []string) *VerifyingErrorManager {
	m := &VerifyingErrorManager{
		ThreadSafeDelegatingErrorManager: jscomp.NewThreadSafeDelegatingErrorManager(delegate),
	}
	m.formatter = jscomp.LightweightMessageFormatter_WithoutSource()
	m.unmatchedExpectations = make([]*jutil.Pattern, 0, len(expectedDiagnostics))
	for _, exp := range expectedDiagnostics {
		m.unmatchedExpectations = append(m.unmatchedExpectations, jutil.Pattern_Compile(exp))
	}
	return m
}

// Java: VerifyingErrorManager.report(CheckLevel, JSError)
func (m *VerifyingErrorManager) Report(level *jscomp.CheckLevel, error *jscomp.JSError) {
	m.mu.Lock()
	defer m.mu.Unlock()
	formattedError := m.FormatErrorForMatching(level, error)

	var matchedPattern *jutil.Pattern
	// Go-only: language-adapter — the Iterator with remove() is an index
	// loop over the slice; i is not advanced past a removed element.
	for i := 0; i < len(m.unmatchedExpectations); {
		p := m.unmatchedExpectations[i]
		if p.Matcher(formattedError).Find() {
			if matchedPattern == nil {
				matchedPattern = p
				m.unmatchedExpectations = slices.Delete(m.unmatchedExpectations, i, i+1)
				continue
			}

			if matchedPattern.String() == p.String() {
				// Allow duplicates that are identical.
				i++
				continue
			}

			m.ThreadSafeDelegatingErrorManager.Report(
				jscomp.CheckLevel_ERROR.Ptr(),
				jscomp.JSError_Make(
					AbstractCommandLineRunner_AMBIGUOUS_EXPECTATION,
					formattedError,
					matchedPattern.String(),
					p.String()))
			m.unmatchedExpectations = slices.Delete(m.unmatchedExpectations, i, i+1)
			continue
		}
		i++
	}

	if matchedPattern == nil {
		m.ThreadSafeDelegatingErrorManager.Report(level, error)
	}
}

// Java: VerifyingErrorManager.generateReport()
func (m *VerifyingErrorManager) GenerateReport() {
	m.mu.Lock()
	m.reportMissingExpectationsLocked()
	m.mu.Unlock()
	m.ThreadSafeDelegatingErrorManager.GenerateReport()
}

// Java: VerifyingErrorManager.reportMissingExpectations()
func (m *VerifyingErrorManager) ReportMissingExpectations() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.reportMissingExpectationsLocked()
}

func (m *VerifyingErrorManager) reportMissingExpectationsLocked() {
	for _, p := range m.unmatchedExpectations {
		m.ThreadSafeDelegatingErrorManager.Report(
			jscomp.CheckLevel_ERROR.Ptr(),
			jscomp.JSError_Make(AbstractCommandLineRunner_EXPECTED_DIAGNOSTIC_NOT_FOUND, p.String()))
	}
	m.unmatchedExpectations = m.unmatchedExpectations[:0]
}

// Java: VerifyingErrorManager.formatErrorForMatching(CheckLevel, JSError)
func (m *VerifyingErrorManager) FormatErrorForMatching(level *jscomp.CheckLevel, error *jscomp.JSError) string {
	if level != nil && *level == jscomp.CheckLevel_ERROR {
		return m.formatter.FormatError(error)
	}
	return m.formatter.FormatWarning(error)
}
