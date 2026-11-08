// Port of com/google/javascript/jscomp/AbstractCommandLineRunner.java @ closure-compiler v20261104
// Port of com/google/javascript/jscomp/VerifyingErrorManager.java @ closure-compiler v20261104

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
// CompileMetricsRecorder records lifecycle and result metrics for a CLI
// compilation action.
//
// Java: interface CompileMetricsRecorderInterface
type CompileMetricsRecorder interface {
	RecordActionStart()
	RecordActionName(actionName string)
	RecordStartState(compiler jscomp.AbstractCompiler)
	RecordResultMetrics(compiler jscomp.AbstractCompiler, result *jscomp.Result)
}

// AbstractCommandLineRunner defines the virtual hook methods on
// AbstractCommandLineRunnerBase that may be overridden by CommandLineRunner or
// custom runners.
//
// Java: abstract class AbstractCommandLineRunner<A extends Compiler, B extends CompilerOptions>
type AbstractCommandLineRunner interface {
	CreateCompiler() *jscomp.Compiler
	CreateOptions() *jscomp.CompilerOptions
	CreateExterns(options *jscomp.CompilerOptions) []*jscomp.SourceFile
	PrepForBundleAndAppendTo(out jutil.Appendable, input *jscomp.CompilerInput, content string)
	AppendRuntimeTo(out jutil.Appendable)
	AddAllowlistWarningsGuard(options *jscomp.CompilerOptions, allowlistFile *jutil.File)
	CheckChunkName(name string)
	GetVersionText() string
	DoRun() int
	SetRunOptions(options *jscomp.CompilerOptions)
	ShouldGenerateMapPerChunk(options *jscomp.CompilerOptions) bool
	WriteChunkOutput(fileName string, out jutil.Appendable, lt jscomp.CodePrinter_LicenseTracker, m *jscomp.JSChunk)
	WriteOutput(
		out jutil.Appendable,
		compiler *jscomp.Compiler,
		licenseTracker jscomp.CodePrinter_LicenseTracker,
		chunk *jscomp.JSChunk,
		wrapper string,
		codePlaceholder string,
		escaper func(string) string,
		filename string,
	)
}

type runnerHooks = AbstractCommandLineRunner

// AbstractCommandLineRunnerBase holds AbstractCommandLineRunner's fields and
// base implementations for embedding in custom runners.
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

	sourceFileCache sourceFileLookupFunc

	metricsRecorder CompileMetricsRecorder
}

type abstractCommandLineRunnerBase = AbstractCommandLineRunnerBase

var (
	// Java: AbstractCommandLineRunner.OUTPUT_SAME_AS_INPUT_ERROR
	AbstractCommandLineRunner_OUTPUT_SAME_AS_INPUT_ERROR = jscomp.DiagnosticType_Error(
		"JSC_OUTPUT_SAME_AS_INPUT_ERROR", "Bad output file (already listed as input file): {0}")

	// Java: AbstractCommandLineRunner.COULD_NOT_SERIALIZE_AST
	couldNotSerializeAstError = jscomp.DiagnosticType_Error(
		"JSC_COULD_NOT_SERIALIZE_AST", "Could not serialize ast to: {0}")

	// Java: AbstractCommandLineRunner.COULD_NOT_DESERIALIZE_AST
	couldNotDeserializeAstError = jscomp.DiagnosticType_Error(
		"JSC_COULD_NOT_DESERIALIZE_AST", "Could not deserialize ast from: {0}")

	// Java: AbstractCommandLineRunner.NO_TREE_GENERATED_ERROR
	noTreeGeneratedError = jscomp.DiagnosticType_Error(
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
const waitingForInputWarning = "The compiler is waiting for input via stdin."

// Use an 8MiB buffer since the concatenated TypedAst file can be very large.
//
// Java: AbstractCommandLineRunner.GZIPPED_TYPEDAST_BUFFER_SIZE
const gzippedTypedAstBufferSize = 8 * 1024 * 1024

// Java: AbstractCommandLineRunner.OUTPUT_MARKER
const defaultOutputMarker = "%output%"

// Java: AbstractCommandLineRunner.OUTPUT_MARKER_JS_STRING
const outputMarkerJsString = "%output|jsstring%"

// InitAbstractCommandLineRunnerBase initializes the common runner base fields
// with the provided streams and virtual hook target.
//
// Java: AbstractCommandLineRunner(InputStream, PrintStream, PrintStream)
func (r *abstractCommandLineRunnerBase) InitAbstractCommandLineRunnerBase(
	self AbstractCommandLineRunner, in io.Reader, out *jutil.PrintStream, err *jutil.PrintStream) {
	r.initRunnerBase(self, in, out, err)
}

func (r *abstractCommandLineRunnerBase) initRunnerBase(
	self runnerHooks, in io.Reader, out *jutil.PrintStream, err *jutil.PrintStream) {
	r.self = self
	r.exitCodeReceiver = AbstractCommandLineRunner_SystemExitCodeReceiver_INSTANCE.Apply
	r.config = NewAbstractCommandLineRunner_CommandLineConfig()
	r.in = jutil.CheckNotNull(in)
	r.defaultJsOutput = jutil.CheckNotNullPtr(out)
	r.err = jutil.CheckNotNullPtr(err)
}

// GetErrorPrintStream returns the runner's error output PrintStream.
//
// Java: AbstractCommandLineRunner.getErrorPrintStream()
func (r *abstractCommandLineRunnerBase) GetErrorPrintStream() *jutil.PrintStream {
	return r.err
}

// SetCompiler overrides the underlying Compiler instance on the runner.
func (r *abstractCommandLineRunnerBase) SetCompiler(compiler *jscomp.Compiler) {
	r.compiler = compiler
}

// SetCompileMetricsRecorder installs an optional CompileMetricsRecorder for
// DoRun lifecycle callbacks.
func (r *abstractCommandLineRunnerBase) SetCompileMetricsRecorder(recorder CompileMetricsRecorder) {
	r.metricsRecorder = recorder
}

// GetCompileMetricsRecorder returns the configured CompileMetricsRecorder, or nil.
//
// Java: AbstractCommandLineRunner.getCompileMetricsRecorder()
func (r *abstractCommandLineRunnerBase) GetCompileMetricsRecorder() CompileMetricsRecorder {
	return r.metricsRecorder
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
func (r *abstractCommandLineRunnerBase) EnableTestMode(
	externsSupplier func() []*jscomp.SourceFile,
	inputsSupplier func() []*jscomp.SourceFile,
	chunksSupplier func() []*jscomp.JSChunk,
	exitCodeReceiver func(int)) {
	if (inputsSupplier == nil) == (chunksSupplier == nil) {
		panic(jutil.NewIllegalArgumentException())
	}
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
func (r *abstractCommandLineRunnerBase) SetExitCodeReceiver(newExitCodeReceiver func(int)) {
	if newExitCodeReceiver == nil {
		panic(jutil.NewNullPointerException())
	}
	r.exitCodeReceiver = newExitCodeReceiver
}

// isOutputInJson returns whether output should be a JSON stream.
//
// Java: AbstractCommandLineRunner.isOutputInJson()
func (r *abstractCommandLineRunnerBase) isOutputInJson() bool {
	return r.config.jsonStreamMode == jscomp.CompilerOptions_JsonStreamMode_OUT ||
		r.config.jsonStreamMode == jscomp.CompilerOptions_JsonStreamMode_BOTH
}

// GetCommandLineConfig gets the command line config, so that it can be
// initialized.
//
// Java: AbstractCommandLineRunner.getCommandLineConfig()
func (r *abstractCommandLineRunnerBase) GetCommandLineConfig() *AbstractCommandLineRunner_CommandLineConfig {
	return r.config
}

// getDiagnosticGroups returns the warning classes that are available from
// the command-line.
//
// Java: AbstractCommandLineRunner.getDiagnosticGroups()
func (r *abstractCommandLineRunnerBase) getDiagnosticGroups() *jscomp.DiagnosticGroups {
	if r.compiler == nil {
		return jscomp.NewDiagnosticGroups()
	}
	return r.compiler.GetDiagnosticGroups()
}

func checkNoDuplicateFlagFilenames(filenames []string, flagName string) {
	if len(filenames) == 0 {
		return
	}
	uniqueNames := make(map[string]struct{}, len(filenames))
	for _, filename := range filenames {
		if _, exists := uniqueNames[filename]; exists {
			panic(jscomp.NewFlagUsageException(
				flagName + " flags specify duplicate file names: " + filename))
		}
		uniqueNames[filename] = struct{}{}
	}
}

// SetRunOptions sets options based on the configurations set flags API.
//
// Java: AbstractCommandLineRunner.setRunOptions(CompilerOptions)
func (r *abstractCommandLineRunnerBase) SetRunOptions(options *jscomp.CompilerOptions) {
	r.setRunOptions(options)
}

func (r *abstractCommandLineRunnerBase) setRunOptions(options *jscomp.CompilerOptions) {
	config := r.config
	diagnosticGroups := r.getDiagnosticGroups()

	if config.ShouldSaveAfterStage1() || config.ShouldContinueCompilation() {
		if options.IsChecksOnly() {
			panic(jscomp.NewFlagUsageException(
				"checks_only mode is incompatible with multi-stage compilation"))
		}
		if options.HasExternExportsPath() {
			panic(jscomp.NewFlagUsageException(
				"generating externs from exports is incompatible with multi-stage compilation"))
		}
	}

	setWarningGuardOptions(options, config.warningGuards, diagnosticGroups)

	if config.warningsAllowFile != "" {
		r.self.AddAllowlistWarningsGuard(options, jutil.NewFile(config.warningsAllowFile))
	}

	if len(config.hideWarningsFor) != 0 {
		options.AddWarningsGuard(
			jscomp.NewShowByPathWarningsGuardWithPathsAndShowType(
				slices.Clone(config.hideWarningsFor),
				jscomp.ShowByPathWarningsGuard_ShowType_EXCLUDE))
	}

	define := slices.Clone(config.define)
	if config.browserFeaturesetYear != 0 {
		if err := setBrowserFeaturesetYear(options, config.browserFeaturesetYear); err != nil {
			panic(err)
		}
	}

	createDefineReplacements(define, options)

	options.SetTweakProcessing(config.tweakProcessing)

	// TODO(tjgq): Unconditionally set the options.
	if config.dependencyOptions != nil {
		options.SetDependencyOptions(config.dependencyOptions)
	}

	options.SetDevMode(config.jscompDevMode)
	options.SetCodingConvention(config.codingConvention)
	options.SetSummaryDetailLevel(config.summaryDetailLevel)
	options.SetTrustedStrings(true)

	r.legacyOutputCharset = r.getLegacyOutputCharset()
	options.SetOutputCharset(r.legacyOutputCharset)
	r.outputCharset2 = r.getOutputCharset2()
	r.inputCharset = r.getInputCharset()

	if config.jsOutputFile != "" && config.skipNormalOutputs {
		panic(jscomp.NewFlagUsageException(
			"skip_normal_outputs and js_output_file cannot be used together."))
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
		options.SetSourceMapOutputPath(config.createSourceMap)
	} else if r.isOutputInJson() {
		options.SetSourceMapOutputPath("%outname%")
	}
	options.SetSourceMapDetailLevel(config.sourceMapDetailLevel)
	options.SetSourceMapFormat(config.sourceMapFormat)
	options.SetSourceMapLocationMappings(config.sourceMapLocationMappings)
	options.SetParseInlineSourceMaps(config.parseInlineSourceMaps)
	options.SetApplyInputSourceMaps(config.applyInputSourceMaps)

	if len(config.sourceMapInputFiles) == 0 {
		options.SetInputSourceMaps(jutil.NewOrderedMap[string, *jscomp.SourceMapInput]())
	} else {
		entries := make([]kvEntry[*jscomp.SourceMapInput], len(config.sourceMapInputFiles))
		for i, entry := range config.sourceMapInputFiles {
			sourceMap := r.newSourceFile(entry.Value, nil, rhino.StaticSourceFile_SourceKind_NON_CODE)
			entries[i] = kvEntry[*jscomp.SourceMapInput]{key: entry.Key, value: jscomp.NewSourceMapInput(sourceMap)}
		}
		options.SetInputSourceMaps(buildUniqueOrderedMap(entries))
	}

	if config.variableMapInputFile != "" {
		options.SetInputVariableMap(jscomp.VariableMap_Load(config.variableMapInputFile))
	}

	if config.propertyMapInputFile != "" {
		options.SetInputPropertyMap(jscomp.VariableMap_Load(config.propertyMapInputFile))
	}

	checkNoDuplicateFlagFilenames(config.outputManifests, "output_manifest")
	checkNoDuplicateFlagFilenames(config.outputBundles, "output_bundle")

	options.SetProcessCommonJSModules(config.processCommonJSModules)
	options.SetModuleRoots(config.moduleRoots)
	options.SetAngularPass(config.angularPass)

	if config.jsonWarningsFile != "" {
		options.AddReportGenerator(
			jscomp.NewJsonErrorReportGenerator(jutil.NewPrintStream(r.filenameToOutputStream(config.jsonWarningsFile)), r.compiler))
	}

	if config.errorFormat == AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_JSON {
		errorGenerator := jscomp.NewJsonErrorReportGenerator(r.err, r.compiler)
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
func (r *abstractCommandLineRunnerBase) GetCompiler() *jscomp.Compiler {
	return r.compiler
}

// Run runs the Compiler and calls System.exit() with the exit status of the
// compiler.
//
// Java: AbstractCommandLineRunner.run()
func (r *abstractCommandLineRunnerBase) Run() {
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

func (r *abstractCommandLineRunnerBase) tryDoRun() (result int, thrown jutil.Throwable) {
	defer func() {
		if rec := recover(); rec != nil {
			thrown = jutil.AsThrowable(rec)
		}
	}()
	return r.self.DoRun(), nil
}

// DoRun parses command-line arguments and runs the compiler. It returns the
// system exit status.
//
// Java: AbstractCommandLineRunner.doRun()
func (r *abstractCommandLineRunnerBase) DoRun() int {
	metricsRecorder := r.GetCompileMetricsRecorder()
	if metricsRecorder != nil {
		metricsRecorder.RecordActionStart()
	}
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
			newVerifyingErrorManager(r.compiler.GetErrorManager(), config.expectedDiagnostics))
	}
	options := r.self.CreateOptions()
	r.self.SetRunOptions(options)

	typedAstListInputFilename, hasTypedAstListInput := config.typedAstListInputFilename.Get()

	externs := r.self.CreateExterns(options)
	var chunks []*jscomp.JSChunk

	r.rootRelativePathsMap = r.constructRootRelativePathsMap()

	var outputFileNames []string
	if config.jsOutputFile != "" {
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
		jsonFiles = r.parseJsonFilesFromInputStream()
		if jsonFiles == nil {
			jsonFiles = []*AbstractCommandLineRunner_JsonFileSpec{}
		}

		smEntries, webpackEntries := collectJsonStreamMaps(jsonFiles)

		if len(smEntries) > 0 {
			for k, v := range options.GetInputSourceMaps().All() {
				if v == nil {
					panic(jutil.NewNullPointerExceptionMsg("null value in entry: " + k + "=null"))
				}
				smEntries = append(smEntries, kvEntry[*jscomp.SourceMapInput]{key: k, value: v})
			}
			options.SetInputSourceMaps(buildUniqueOrderedMap(smEntries))
		}

		r.compiler.InitWebpackMap(buildUniqueOrderedMap(webpackEntries))
	} else {
		r.compiler.InitWebpackMap(&jutil.OrderedMap[string, string]{})
	}

	options.SetDoLateLocalization(config.ShouldDoLateLocalization())
	options.SetAlwaysGatherSourceMapInfo(config.ShouldAlwaysGatherSourceMapInfo())

	r.compiler.InitOptions(options)

	sources := r.createSourceInputs(jsChunkSpecs, config.mixedJsSources, jsonFiles, config.moduleRoots)

	if len(jsChunkSpecs) != 0 {
		if r.testMode {
			chunks = r.chunksSupplierForTesting()
		} else {
			if config.jsonStreamMode == jscomp.CompilerOptions_JsonStreamMode_IN ||
				config.jsonStreamMode == jscomp.CompilerOptions_JsonStreamMode_NONE {
				for _, m := range jsChunkSpecs {
					r.self.CheckChunkName(m.GetName())
				}
			}
			inputs := make([]*jscomp.CompilerInput, 0, len(sources))
			for _, source := range sources {
				inputs = append(inputs, jscomp.NewCompilerInput(source))
			}
			chunks = createJsChunks(jsChunkSpecs, inputs)
		}
		for _, m := range chunks {
			outputFileNames = append(outputFileNames, r.GetChunkOutputFileName(m))
		}

		if hasTypedAstListInput {
			r.initChunksWithTypedAstFilesystem(externs, chunks, options, typedAstListInputFilename)
		} else {
			r.compiler.InitChunks(externs, chunks, options)
		}
	} else {
		if hasTypedAstListInput {
			r.initWithTypedAstFilesystem(externs, sources, options, typedAstListInputFilename)
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
		if metricsRecorder != nil {
			metricsRecorder.RecordActionName("skip normal outputs")
		}
		// TODO(bradfordcsmith): Should we be ignoring possible init/initChunks() errors here?
		r.compiler.OrderInputsWithLargeStack()
	} else if r.compiler.HasErrors() {
		if metricsRecorder != nil {
			metricsRecorder.RecordActionName("initialization errors occurred")
		}
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
			r.performCompilation(metricsRecorder)

			// If we're finished with compilation (i.e. we're not saving state), /and/ the compiler was
			// restored from a previous state, then we need to re-initialize the set of chunks.
			// TODO(lharker): figure out if this is still needed.
			if chunks != nil &&
				config.restoredCompilationStage != -1 &&
				config.saveAfterCompilationStage == -1 {
				chunks = slices.Clone(r.compiler.GetChunks())
			}
		}()
		result = r.compiler.GetResult()
	}

	if createCommonJsModules {
		// For CommonJS modules, construct chunks from actual inputs.
		chunks = slices.Clone(r.compiler.GetChunks())
		for _, c := range chunks {
			outputFileNames = append(outputFileNames, r.GetChunkOutputFileName(c))
		}
	}

	for _, outputFileName := range outputFileNames {
		if r.compiler.GetSourceFileByName(outputFileName) != nil {
			r.compiler.Report(jscomp.JSError_Make(AbstractCommandLineRunner_OUTPUT_SAME_AS_INPUT_ERROR, outputFileName))
			return 1
		}
	}

	// We won't want to process results for cases where compilation is only partially done.
	var exitStatus int
	if !config.saveCompilationStateToFilename.IsPresent() {
		exitStatus = r.processResults(result, chunks, options)
	} else {
		exitStatus = exitStatusForResult(result)
	}
	if metricsRecorder != nil {
		metricsRecorder.RecordResultMetrics(r.compiler, result)
	}
	return exitStatus
}

// Java: AbstractCommandLineRunner.performCompilation(CompileMetricsRecorderInterface)
func (r *abstractCommandLineRunnerBase) performCompilation(metricsRecorder CompileMetricsRecorder) {
	// Parse, restore from a save file, or initialize from a TypedAST list.
	r.initializeStateBeforeCompilation()

	// Do the interesting work - checks, optimizations, etc.
	r.runCompilerPasses(metricsRecorder)

	// Save state unless there was an exception or error - we won't go on to the next stage if there
	// was an error, and saving state itself might throw an exception if the compiler is in a bad
	// state.
	if !r.compiler.HasErrors() {
		r.saveState()
	}
	// Perform post-compilation tasks even if compiler.hasErrors()
	r.compiler.PerformPostCompilationTasks()
}

// Java: AbstractCommandLineRunner.runCompilerPasses(CompileMetricsRecorderInterface)
func (r *abstractCommandLineRunnerBase) runCompilerPasses(metricsRecorder CompileMetricsRecorder) {
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

	if metricsRecorder != nil {
		metricsRecorder.RecordActionName(actionMetricsName)
	}
	if r.compiler.HasErrors() {
		return
	}
	if metricsRecorder != nil {
		metricsRecorder.RecordStartState(r.compiler)
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

// initializeStateBeforeCompilation is called at the beginning of
// compilation to initialize the compiler state.
//
// Compiler state should be initialized no matter what the compilation stage
// is, but this method handles the different ways it might be initialized -
// whether from parsing actual JS files, from reading library-level
// TypedASTs, or from restoring a previous compilation state.
//
// Java: AbstractCommandLineRunner.initializeStateBeforeCompilation()
func (r *abstractCommandLineRunnerBase) initializeStateBeforeCompilation() {
	if r.config.restoredCompilationStage != -1 {
		fn, ok := r.config.continueSavedCompilationFileName.Get()
		if !ok {
			panic(jutil.NewNullPointerException())
		}
		r.restoreState(fn)
	} else if !r.config.typedAstListInputFilename.IsPresent() {
		r.compiler.ParseForCompilation()
	}
}

// processResults processes the results of the compile job, and returns an
// error code.
//
// Java: AbstractCommandLineRunner.processResults(Result, List<JSChunk>, B)
func (r *abstractCommandLineRunnerBase) processResults(result *jscomp.Result, chunks []*jscomp.JSChunk, options *jscomp.CompilerOptions) int {
	config := r.config
	if config.printAst {
		if r.compiler.GetRoot().IsNil() {
			return 1
		}
		jsOutput := r.createDefaultOutput()
		cfg := r.compiler.ComputeCFG()
		jscomp.DotFormatter_AppendDot(r.compiler.GetRoot().GetLastChild(), cfg, jsOutput)
		jsOutput.AppendCharacter('\n')
		closeAppendable(jsOutput)
		return 0
	}

	if config.printTree {
		if r.compiler.GetRoot().IsNil() {
			r.compiler.Report(jscomp.JSError_Make(noTreeGeneratedError))
			return 1
		}
		jsOutput := r.createDefaultOutput()
		r.compiler.GetJsRoot().AppendStringTree(jsOutput)
		jsOutput.AppendCharSequence("\n")
		closeAppendable(jsOutput)
		return 0
	}

	if config.printTreeJson {
		if r.compiler.GetRoot().IsNil() {
			r.compiler.Report(jscomp.JSError_Make(noTreeGeneratedError))
			return 1
		}
		jsOutput := r.createDefaultOutput()
		r.compiler.GetJsRoot().AppendJsonTree(jsOutput)
		closeAppendable(jsOutput)
		return 0
	}

	if config.skipNormalOutputs {
		// Output the manifest and bundle files if requested.
		r.outputManifest()
		r.outputBundle()
		r.outputChunkGraphJson()
		return 0
	}

	if options.GetOutputJs() != jscomp.CompilerOptions_OutputJs_NONE && result.Success {
		r.outputChunkGraphJson()
		if chunks == nil {
			r.outputSingleBinary(options)

			// Output the source map if requested.
			// If output files are being written to stdout as a JSON string,
			// outputSingleBinary will have added the sourcemap to the output file
			if !r.isOutputInJson() {
				r.outputSourceMap(options, config.jsOutputFile)
			}
		} else if diagErr := r.outputChunkBinaryAndSourceMaps(r.compiler.GetChunkGraph(), options); diagErr != nil {
			r.compiler.Report(jscomp.JSError_Make(diagErr))
			return 1
		}

		// Output the externs if required.
		if options.HasExternExportsPath() {
			func() {
				eeOut := r.openExternExportsStream(options, config.jsOutputFile)
				defer closeBufferedWriter(eeOut)
				// Java: Writer.append(null) appends "null".
				externExport := "null"
				if ee, ok := result.ExternExportString(); ok {
					externExport = ee
				}
				eeOut.AppendCharSequence(externExport)
			}()
		}

		// Output the variable and property name maps if requested.
		r.outputNameMaps()

		// Output the ReplaceStrings map if requested
		r.outputStringMap()

		// Output the manifest and bundle files if requested.
		r.outputManifest()
		r.outputBundle()

		// Output the production instrumentation param mapping if requested.
		r.outputInstrumentationMapping()

		if r.isOutputInJson() {
			r.outputJsonStream()
		}
	}

	return exitStatusForResult(result)
}

// Java: AbstractCommandLineRunner.getExitStatusForResult(Result)
func exitStatusForResult(result *jscomp.Result) int {
	// return 0 if no errors, the error count otherwise
	return min(len(result.Errors), 0x7f)
}

type kvEntry[V any] struct {
	key   string
	value V
}

// checkUniqueKeys validates that entries contain no duplicate keys, panicking
// with Guava ImmutableMap.Builder.buildOrThrow's exact error message on conflict.
func checkUniqueKeys[E any](entries []E, keyOf func(E) string, valOf func(E) string) {
	n := len(entries)
	if n <= 1 {
		return
	}
	if n > 8 {
		seen := make(map[string]struct{}, n)
		hasDup := false
		for _, e := range entries {
			k := keyOf(e)
			if _, ok := seen[k]; ok {
				hasDup = true
				break
			}
			seen[k] = struct{}{}
		}
		if !hasDup {
			return
		}
	}
	for i := n - 1; i >= 0; i-- {
		k := keyOf(entries[i])
		for j := i + 1; j < n; j++ {
			if keyOf(entries[j]) == k {
				panic(jutil.NewIllegalArgumentExceptionMsg("Multiple entries with same key: " +
					k + "=" + valOf(entries[j]) + " and " +
					k + "=" + valOf(entries[i])))
			}
		}
	}
}

func buildUniqueOrderedMap[V any](entries []kvEntry[V]) *jutil.OrderedMap[string, V] {
	n := len(entries)
	out := jutil.NewOrderedMapWithCapacity[string, V](n)
	if n == 0 {
		return out
	}
	checkUniqueKeys(entries,
		func(e kvEntry[V]) string { return e.key },
		func(e kvEntry[V]) string { return jutil.ToJavaString(e.value) })
	for _, e := range entries {
		out.PutNew(e.key, e.value)
	}
	return out
}

// verifyingErrorManager is an ErrorManager that verifies expected
// diagnostics and reports missing or unexpected ones.
//
// It wraps a delegate ErrorManager. Expected diagnostics (passed as regexes)
// are swallowed if matched. Unexpected diagnostics are passed to the
// delegate. At the end of compilation (when generateReport is called), any
// unmatched expectations are reported as errors to the delegate.
type verifyingErrorManager struct {
	*jscomp.ThreadSafeDelegatingErrorManager
	mu                    sync.Mutex
	unmatchedExpectations []*jutil.Pattern
	formatter             *jscomp.LightweightMessageFormatter
}

var _ jscomp.ErrorManager = (*verifyingErrorManager)(nil)

// Java: VerifyingErrorManager(ErrorManager, List<String>)
func newVerifyingErrorManager(delegate jscomp.ErrorManager, expectedDiagnostics []string) *verifyingErrorManager {
	unmatched := make([]*jutil.Pattern, len(expectedDiagnostics))
	for i, exp := range expectedDiagnostics {
		unmatched[i] = jutil.Pattern_Compile(exp)
	}
	return &verifyingErrorManager{
		ThreadSafeDelegatingErrorManager: jscomp.NewThreadSafeDelegatingErrorManager(delegate),
		unmatchedExpectations:            unmatched,
		formatter:                        jscomp.LightweightMessageFormatter_WithoutSource(),
	}
}

// Java: VerifyingErrorManager.report(CheckLevel, JSError)
func (m *verifyingErrorManager) Report(level jscomp.CheckLevel, jsErr *jscomp.JSError) {
	m.mu.Lock()
	defer m.mu.Unlock()
	formattedError := m.formatErrorForMatching(level, jsErr)

	var matchedPattern *jutil.Pattern
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
				jscomp.CheckLevel_ERROR,
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
		m.ThreadSafeDelegatingErrorManager.Report(level, jsErr)
	}
}

// Java: VerifyingErrorManager.generateReport()
func (m *verifyingErrorManager) GenerateReport() {
	m.mu.Lock()
	m.reportMissingExpectationsLocked()
	m.mu.Unlock()
	m.ThreadSafeDelegatingErrorManager.GenerateReport()
}

func (m *verifyingErrorManager) reportMissingExpectationsLocked() {
	for _, p := range m.unmatchedExpectations {
		m.ThreadSafeDelegatingErrorManager.Report(
			jscomp.CheckLevel_ERROR,
			jscomp.JSError_Make(AbstractCommandLineRunner_EXPECTED_DIAGNOSTIC_NOT_FOUND, p.String()))
	}
	m.unmatchedExpectations = m.unmatchedExpectations[:0]
}

// Java: VerifyingErrorManager.formatErrorForMatching(CheckLevel, JSError)
func (m *verifyingErrorManager) formatErrorForMatching(level jscomp.CheckLevel, jsErr *jscomp.JSError) string {
	if level == jscomp.CheckLevel_ERROR {
		return m.formatter.FormatError(jsErr)
	}
	return m.formatter.FormatWarning(jsErr)
}
