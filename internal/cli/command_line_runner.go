// Port of com/google/javascript/jscomp/CommandLineRunner.java @ closure-compiler v20261021

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
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/args4j"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/transpile"
)

// CommandLineRunner translates flags into Java API calls on the Compiler.
//
// This class may be extended and used to create other Java classes that
// behave the same as running the Compiler from the command line. If you want
// to run the compiler in-process in Java, you should look at this class for
// hints on what API calls to make, but you should not use this class
// directly.
//
// Example:
//
//	class MyCommandLineRunner extends CommandLineRunner {
//	  MyCommandLineRunner(String[] args) {
//	    super(args);
//	  }
//
//	  @Override protected CompilerOptions createOptions() {
//	    CompilerOptions options = super.createOptions();
//	    addMyCrazyCompilerPassThatOutputsAnExtraFile(options);
//	    return options;
//	  }
//
//	  public static void main(String[] args) {
//	    MyCommandLineRunner runner = new MyCommandLineRunner(args);
//	    if (runner.shouldRunCompiler()) {
//	      runner.run();
//	    }
//	    if (runner.hasErrors()) {
//	      System.exit(-1);
//	    }
//	  }
//	}
//
// This class is totally not thread-safe.
type CommandLineRunner struct {
	AbstractCommandLineRunnerBase

	flags *commandLineRunner_Flags

	errors bool

	runCompiler bool

	// Cached error stream to avoid passing it as a parameter to helper
	// functions.
	errorStream *jutil.PrintStream // nil is null

	bundler *transpile.ClosureBundler // nil is null
}

var _ AbstractCommandLineRunner = (*CommandLineRunner)(nil)

// Java: CommandLineRunner.OUTPUT_MARKER
const CommandLineRunner_OUTPUT_MARKER = AbstractCommandLineRunner_OUTPUT_MARKER

// UTF-8 BOM is 0xEF, 0xBB, 0xBF, of which character code is 65279.
//
// Java: CommandLineRunner.UTF8_BOM_CODE
const CommandLineRunner_UTF8_BOM_CODE = 65279

// Allowable chunk name characters that aren't valid in a JS identifier
//
// Java: CommandLineRunner.extraChunkNameChars
var commandLineRunner_extraChunkNameChars = jutil.Pattern_Compile("[-.]+")

// CommandLineRunner_FormattingOption is the set of options that can be used
// with the --formatting flag.
//
// Java: CommandLineRunner.FormattingOption
type CommandLineRunner_FormattingOption int32

const (
	CommandLineRunner_FormattingOption_PRETTY_PRINT CommandLineRunner_FormattingOption = iota
	CommandLineRunner_FormattingOption_PRINT_INPUT_DELIMITER
	CommandLineRunner_FormattingOption_SINGLE_QUOTES
)

var commandLineRunner_FormattingOption_names = [...]string{"PRETTY_PRINT", "PRINT_INPUT_DELIMITER", "SINGLE_QUOTES"}

// CommandLineRunner_FormattingOption_values is FormattingOption.values().
func CommandLineRunner_FormattingOption_values() []CommandLineRunner_FormattingOption {
	return []CommandLineRunner_FormattingOption{
		CommandLineRunner_FormattingOption_PRETTY_PRINT,
		CommandLineRunner_FormattingOption_PRINT_INPUT_DELIMITER,
		CommandLineRunner_FormattingOption_SINGLE_QUOTES,
	}
}

// String is Enum.name().
func (o CommandLineRunner_FormattingOption) String() string {
	return commandLineRunner_FormattingOption_names[o]
}

// Ordinal is Enum.ordinal().
func (o CommandLineRunner_FormattingOption) Ordinal() int { return int(o) }

// Java: CommandLineRunner.FormattingOption.applyToOptions(CompilerOptions)
func (o CommandLineRunner_FormattingOption) ApplyToOptions(options *jscomp.CompilerOptions) {
	switch o {
	case CommandLineRunner_FormattingOption_PRETTY_PRINT:
		options.SetPrettyPrint(true)
	case CommandLineRunner_FormattingOption_PRINT_INPUT_DELIMITER:
		options.SetPrintInputDelimiter(true)
	case CommandLineRunner_FormattingOption_SINGLE_QUOTES:
		options.SetPreferSingleQuotes(true)
	}
}

// NewCommandLineRunner creates a new command-line runner. You should only
// need to call the constructor if you're extending this class. Otherwise,
// the main method should instantiate it.
//
// Java: CommandLineRunner(String[])
func NewCommandLineRunner(args []string) *CommandLineRunner {
	return NewCommandLineRunnerWithStreams(args, jutil.System_In(), jutil.System_Out(), jutil.System_Err())
}

// Java: CommandLineRunner(String[], PrintStream, PrintStream)
func NewCommandLineRunnerWithOutAndErr(args []string, out *jutil.PrintStream, err *jutil.PrintStream) *CommandLineRunner {
	return NewCommandLineRunnerWithStreams(args, jutil.System_In(), out, err)
}

// Java: CommandLineRunner(String[], InputStream, PrintStream, PrintStream)
func NewCommandLineRunnerWithStreams(args []string, in io.Reader, out *jutil.PrintStream, err *jutil.PrintStream) *CommandLineRunner {
	r := &CommandLineRunner{}
	r.InitAbstractCommandLineRunnerBase(r, in, out, err)
	r.flags = newCommandLineRunner_Flags()
	r.errors = false
	r.runCompiler = false
	r.InitConfigFromFlags(args, out, err)
	return r
}

var commandLineRunner_quotesPattern = jutil.Pattern_Compile("^['\"](.*)['\"]$")

// Java: CommandLineRunner.processArgs(String[])
func commandLineRunner_processArgs(args []string) []string {
	// Args4j has a different format that the old command-line parser.
	// So we use some voodoo to get the args into the format that args4j
	// expects.
	processedArgs := make([]string, 0, len(args))

	for _, arg := range args {
		if name, value, ok := commandLineRunner_matchArgPattern(arg); ok {
			processedArgs = append(processedArgs, name)

			if len(value) >= 2 && (value[0] == '\'' || value[0] == '"') {
				if quotesMatcher := commandLineRunner_quotesPattern.Matcher(value); quotesMatcher.Matches() {
					processedArgs = append(processedArgs, *quotesMatcher.Group(1))
					continue
				}
			}
			processedArgs = append(processedArgs, value)
		} else {
			processedArgs = append(processedArgs, arg)
		}
	}

	return processedArgs
}

// commandLineRunner_matchArgPattern is
// Pattern.compile("(--?[a-zA-Z_]+)=(.*)", Pattern.DOTALL).matcher(arg)
// .matches(), returning groups 1 and 2: one or two '-', then one or more
// ASCII letters or '_' (which cannot match '=', so backtracking never
// finds another split), then '=', then anything (DOTALL: line terminators
// included).
//
// Go-only: language-adapter — Pattern.DOTALL, which jutil.Pattern does not
// support.
func commandLineRunner_matchArgPattern(arg string) (name, value string, ok bool) {
	i := 0
	if i < len(arg) && arg[i] == '-' {
		i++
	} else {
		return "", "", false
	}
	if i < len(arg) && arg[i] == '-' {
		i++
	}
	start := i
	for i < len(arg) && (arg[i] >= 'a' && arg[i] <= 'z' || arg[i] >= 'A' && arg[i] <= 'Z' || arg[i] == '_') {
		i++
	}
	if i == start || i >= len(arg) || arg[i] != '=' {
		return "", "", false
	}
	return arg[:i], arg[i+1:], true
}

// Java: CommandLineRunner.reportError(String)
func (r *CommandLineRunner) ReportError(message string) {
	r.errors = true
	r.errorStream.Println(message)
	r.errorStream.Flush()
}

// Java: CommandLineRunner.processFlagFiles()
func (r *CommandLineRunner) ProcessFlagFiles() {
	for _, flagFile := range r.flags.flagFiles {
		jutil.Try(func() {
			r.ProcessFlagFile(flagFile)
		}, func(ioErr jutil.IIOException) {
			r.ReportError("ERROR - " + flagFile + " read error.")
		})
	}
}

// Java: CommandLineRunner.processFlagFile(String)
func (r *CommandLineRunner) ProcessFlagFile(flagFileString string) {
	content := strings.TrimPrefix(
		jutil.Files_ReadStringUTF8(jutil.Path_Of(flagFileString)),
		"\uFEFF",
	)

	var tokens []string
	var token []byte
	quoted := false
	escaped := false

	for i := 0; i < len(content); i++ {
		c := content[i]
		switch c {
		case ' ', '\t', '\n', '\r':
			if quoted {
				token = append(token, c)
			} else if len(token) != 0 {
				tokens = append(tokens, string(token))
				token = token[:0]
			}
		case '"':
			if escaped {
				if quoted {
					token[len(token)-1] = '"'
				} else {
					token = append(token, '"')
				}
			} else {
				quoted = !quoted
			}
		default:
			token = append(token, c)
		}
		escaped = c == '\\'
	}

	if len(token) != 0 {
		tokens = append(tokens, string(token))
	}

	r.flags.flagFiles = []string{}

	tokens = commandLineRunner_processArgs(tokens)

	// Command-line warning levels should override flag file settings,
	// which means they should go last.
	previousGuardLevels := slices.Clone(r.flags.guardLevels)
	previousMixedJsSources := slices.Clone(r.flags.mixedJsSources)
	r.flags.guardLevels = r.flags.guardLevels[:0]
	r.flags.mixedJsSources = r.flags.mixedJsSources[:0]
	r.flags.Parse(tokens)
	r.flags.guardLevels = append(r.flags.guardLevels, previousGuardLevels...)
	r.flags.mixedJsSources = append(r.flags.mixedJsSources, previousMixedJsSources...)

	// Currently we are not supporting this (prevent direct/indirect loops)
	if len(r.flags.flagFiles) != 0 {
		r.ReportError("ERROR - Arguments in the file cannot contain " + "--flagfile option.")
	}
}

// Java: CommandLineRunner.getVersionText()
func (r *CommandLineRunner) GetVersionText() string {
	return strings.Join([]string{
		"Closure Compiler (http://github.com/google/closure-compiler)",
		// CommandLineRunnerVersion refers to class created from CommandLineRunnerVersion.template
		// with its static field COMPILER_VERSION substituted at build time.
		"Version: " + jscomp.CommandLineRunnerVersion_COMPILER_VERSION}, "\n")
}

// Java: CommandLineRunner.initConfigFromFlags(String[], PrintStream, PrintStream)
func (r *CommandLineRunner) InitConfigFromFlags(args []string, out *jutil.PrintStream, err *jutil.PrintStream) {
	flags := r.flags

	r.errorStream = err
	processedArgs := commandLineRunner_processArgs(args)

	flags.guardLevels = flags.guardLevels[:0]
	flags.mixedJsSources = flags.mixedJsSources[:0]

	var mixedSources []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]
	var mappings []jscomp.SourceMap_LocationMapping
	var sourceMapInputs []PipeFlagEntry
	parseInlineSourceMaps := false
	applyInputSourceMaps := false
	jutil.Try(func() {
		flags.Parse(processedArgs)

		r.ProcessFlagFiles()

		mixedSources = flags.GetMixedJsSources()
		mappings = flags.GetSourceMapLocationMappings()
		sourceMapInputs = flags.GetSourceMapInputs()
		parseInlineSourceMaps = flags.parseInlineSourceMaps
		applyInputSourceMaps = flags.applyInputSourceMaps
	}, func(t jutil.Throwable) {
		switch e := t.(type) {
		case args4j.ICmdLineException:
			r.ReportError(e.GetMessage())
		case jutil.IIOException:
			r.ReportError("ERROR - ioException: " + e.Error())
		default:
			panic(t)
		}
	})

	if flags.processCommonJsModules {
		flags.processClosurePrimitives = true
	}

	if flags.browserFeaturesetYear != 0 && !flags.languageOutIsDefault {
		panic(jscomp.NewFlagUsageException(
			"ERROR - both flags `--browser_featureset_year` and `--language_out` specified."))
	}

	// Java: if (flags.outputWrapper == null) flags.outputWrapper = ""; the
	// field is a non-null String here.

	if flags.outputWrapperFile != "" {
		jutil.Try(func() {
			flags.outputWrapper = jutil.Files_AsCharSource(jutil.NewFile(flags.outputWrapperFile), jutil.StandardCharsets_UTF_8).Read()
		}, func(e jutil.IException) {
			r.ReportError("ERROR - invalid output_wrapper_file specified.")
		})
	}

	if flags.outputWrapper != "" &&
		!strings.Contains(flags.outputWrapper, CommandLineRunner_OUTPUT_MARKER) {
		r.ReportError(
			"ERROR - invalid output_wrapper specified. Missing '" +
				CommandLineRunner_OUTPUT_MARKER +
				"'.")
	}

	if flags.outputWrapper != "" && flags.isolationMode != jscomp.CompilerOptions_IsolationMode_NONE {
		r.ReportError("--output_wrapper and --isolation_mode may not be used together.")
	}

	if flags.isolationMode == jscomp.CompilerOptions_IsolationMode_IIFE {
		flags.outputWrapper = "(function(){%output%}).call(this);"
	}

	// Handle --compilation_level=BUNDLE
	bundleFiles := []string{}
	skipNormalOutputs := false
	if flags.compilationLevelParsed != nil && *flags.compilationLevelParsed == jscomp.CompilationLevel_BUNDLE {
		if flags.jsOutputFile == "" {
			r.ReportError("--compilation_level=BUNDLE cannot be used without a --js_output_file.")
		} else {
			bundleFiles = []string{flags.jsOutputFile}
			flags.jsOutputFile = ""
			skipNormalOutputs = true
		}
	}

	var conv jscomp.CodingConvention
	if flags.thirdParty {
		conv = jscomp.CodingConventions_GetDefault()
	} else if flags.chromePass {
		conv = jscomp.NewChromeCodingConvention()
	} else {
		conv = jscomp.NewClosureCodingConvention()
	}

	if !flags.renaming &&
		flags.compilationLevelParsed != nil && *flags.compilationLevelParsed == jscomp.CompilationLevel_ADVANCED_OPTIMIZATIONS {
		r.ReportError("ERROR - renaming cannot be disabled when ADVANCED_OPTIMIZATIONS is used.")
	}

	var dependencyOptions *jscomp.DependencyOptions
	jutil.Try(func() {
		dependencyOptions = jscomp.DependencyOptions_FromFlags(
			flags.dependencyMode, flags.entryPoint, []string{}, nil, false, false)
	}, func(e jscomp.IFlagUsageException) {
		r.ReportError(e.GetMessage())
	})

	if r.errors {
		commandLineRunner_Flags_printShortUsageAfterErrors(r.errorStream)
	} else if flags.displayHelp || flags.helpMarkdown {
		flags.PrintUsage(out)
	} else {
		r.runCompiler = true

		config := r.GetCommandLineConfig()
		outputChunkDependencies := flags.outputChunkDependencies
		config.
			SetPrintVersion(flags.version).
			SetPrintTree(flags.printTree).
			SetPrintTreeJson(flags.printTreeJson).
			SetPrintAst(flags.printAst).
			SetJscompDevMode(flags.jscompDevMode).
			SetLoggingLevel(flags.loggingLevel).
			SetExterns(flags.externs).
			SetMixedJsSources(mixedSources).
			SetDefaultToStdin().
			SetJsOutputFile(flags.jsOutputFile).
			SetChunk(flags.chunk).
			SetVariableMapOutputFile(flags.variableMapOutputFile).
			SetCreateNameMapFiles(flags.createNameMapFiles).
			SetPropertyMapOutputFile(flags.propertyMapOutputFile).
			SetExpectedDiagnostics(flags.expectedDiagnostics).
			SetPropertyMapInputFile(flags.propertyMapInputFile).
			SetVariableMapInputFile(flags.variableMapInputFile).
			SetInstrumentationMappingFile(flags.instrumentationMappingOutputFile).
			SetCodingConvention(conv).
			SetSummaryDetailLevel(flags.summaryDetailLevel).
			SetOutputWrapper(flags.outputWrapper).
			SetChunkWrapper(flags.chunkWrapper).
			SetChunkOutputPathPrefix(flags.chunkOutputPathPrefix).
			SetCreateSourceMap(flags.createSourceMap).
			SetSourceMapFormat(flags.sourceMapFormat).
			SetSourceMapLocationMappings(mappings).
			SetSourceMapInputFiles(sourceMapInputs).
			SetParseInlineSourceMaps(parseInlineSourceMaps).
			SetApplyInputSourceMaps(applyInputSourceMaps).
			SetWarningGuards(flags.guardLevels).
			SetDefine(flags.define).
			SetBrowserFeaturesetYear(int(flags.browserFeaturesetYear)).
			SetCharset(flags.charset).
			SetDependencyOptions(dependencyOptions).
			SetOutputManifest([]string{flags.outputManifest}).
			SetOutputBundle(bundleFiles).
			SetSkipNormalOutputs(skipNormalOutputs).
			SetOutputChunkDependencies(&outputChunkDependencies).
			SetProcessCommonJSModules(flags.processCommonJsModules).
			SetModuleRoots(flags.moduleRoot).
			SetWarningsAllowlistFile(flags.warningsAllowlistFile).
			SetHideWarningsFor(flags.hideWarningsFor).
			SetAngularPass(flags.angularPass).
			SetJsonStreamMode(flags.jsonStreamMode).
			SetErrorFormat(flags.errorFormat)

		segmentOfCompilationToRun := flags.segmentOfCompilationToRun
		switch segmentOfCompilationToRun {
		case jscomp.CompilerOptions_SegmentOfCompilationToRun_CHECKS:
			jutil.CheckStateMsg(
				flags.filenameToRestoreFrom == nil,
				"Cannot restore and run CHECKS segment of compilation")
			config.SetSaveCompilationStateToFilename(flags.filenameToSaveTo, 1)
		case jscomp.CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS:
			config.SetContinueSavedCompilationFileName(flags.filenameToRestoreFrom, 1)
			config.SetSaveCompilationStateToFilename(flags.filenameToSaveTo, 2)
		case jscomp.CompilerOptions_SegmentOfCompilationToRun_FINALIZATIONS:
			jutil.CheckStateMsg(
				flags.filenameToSaveTo == nil,
				"Cannot run FINALIZATIONS segment of compilation and then save the result")
			config.SetContinueSavedCompilationFileName(flags.filenameToRestoreFrom, 2)
		case jscomp.CompilerOptions_SegmentOfCompilationToRun_ENTIRE_COMPILATION:
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Cannot run %s segment of compilation: " + segmentOfCompilationToRun.String()))
		}
	}

	r.errorStream = nil
}

// Java: CommandLineRunner.addAllowlistWarningsGuard(CompilerOptions, File)
func (r *CommandLineRunner) AddAllowlistWarningsGuard(options *jscomp.CompilerOptions, allowlistFile *jutil.File) {
	options.AddWarningsGuard(AllowlistWarningsGuard_FromFile(allowlistFile))
}

// Java: CommandLineRunner.checkChunkName(String)
func (r *CommandLineRunner) CheckChunkName(name string) {
	if !rhino.TokenStream_IsJSIdentifier(commandLineRunner_extraChunkNameChars.Matcher(name).ReplaceAll("_")) {
		panic(jscomp.NewFlagUsageException("Invalid chunk name: '" + name + "'"))
	}
}

// Java: CommandLineRunner.createOptions()
func (r *CommandLineRunner) CreateOptions() *jscomp.CompilerOptions {
	flags := r.flags
	options := jscomp.NewCompilerOptions()

	if flags.languageIn != "" {
		languageMode := jscomp.CompilerOptions_LanguageMode_FromString(&flags.languageIn)
		if languageMode != nil && *languageMode == jscomp.CompilerOptions_LanguageMode_UNSUPPORTED {
			panic(jscomp.NewFlagUsageException(
				"Cannot specify the unsupported set of features for language_in."))
		}
		if languageMode != nil {
			options.SetLanguageIn(*languageMode)
		} else {
			panic(jscomp.NewFlagUsageException(
				"Unknown language `" +
					flags.languageIn +
					"' specified. Expected one of: " +
					jutil.CollectionToString(jscomp.CompilerOptions_LanguageMode_ValidCommandLineNames())))
		}
	}

	languageMode := jscomp.CompilerOptions_LanguageMode_FromString(&flags.languageOut)
	if languageMode != nil && *languageMode == jscomp.CompilerOptions_LanguageMode_UNSUPPORTED {
		panic(jscomp.NewFlagUsageException(
			"Cannot specify the unsupported set of features for language_out."))
	}
	if languageMode != nil {
		options.SetLanguageOut(*languageMode)
	} else {
		panic(jscomp.NewFlagUsageException(
			"Unknown language `" +
				flags.languageOut +
				"' specified. Expected one of: " +
				jutil.CollectionToString(jscomp.CompilerOptions_LanguageMode_ValidCommandLineNames())))
	}

	options.SetCodingConvention(jscomp.NewClosureCodingConvention())

	options.SetExtraAnnotationNames(flags.extraAnnotationName)

	level := *flags.compilationLevelParsed
	level.SetOptionsForCompilationLevel(options)

	if flags.debug {
		level.SetDebugOptionsForCompilationLevel(options)
	}

	options.SetNumParallelThreads(flags.numParallelThreads)
	if flags.fastMode {
		options.SetFastMode(true)
	}

	options.SetEnvironment(flags.environment)

	options.SetChecksOnly(flags.checksOnly)
	if flags.checksOnly {
		options.SetOutputJs(jscomp.CompilerOptions_OutputJs_NONE)
	}

	options.SetIncrementalChecks(flags.incrementalCheckMode)

	options.SetContinueAfterErrors(flags.continueAfterErrors)

	if flags.useTypesForOptimization {
		level.SetTypeBasedOptimizationOptions(options)
	}

	if flags.assumeFunctionWrapper ||
		flags.isolationMode == jscomp.CompilerOptions_IsolationMode_IIFE ||
		flags.chunkOutputType == jscomp.CompilerOptions_ChunkOutputType_ES_MODULES {
		level.SetWrappedOutputOptimizations(options)
		options.SetAssumeGlobalScopeIsIsolated(true)
	}

	if flags.typedAstOutputFile != nil {
		options.SetTypedAstOutputFile(jutil.Path_Of(*flags.typedAstOutputFile))
	}
	options.SetGenerateExports(flags.generateExports)
	options.SetExportLocalPropertyDefinitions(flags.exportLocalPropertyDefinitions)

	wLevel := flags.warningLevel
	wLevel.SetOptionsForWarningLevel(options)
	for _, formattingOption := range flags.formatting {
		formattingOption.ApplyToOptions(options)
	}

	options.SetClosurePass(flags.processClosurePrimitives)

	options.SetAngularPass(flags.angularPass)

	options.SetPolymerVersion(flags.polymerVersion)

	options.SetChromePass(flags.chromePass)

	if flags.j2clPassMode != "" {
		jutil.Try(func() {
			j2clPassMode := commandLineRunner_j2clPassModeValueOf(asciiUpper(flags.j2clPassMode))
			options.SetJ2clPass(j2clPassMode)
		}, func(ex jutil.IIllegalArgumentException) {
			panic(jscomp.NewFlagUsageException(
				"Unknown J2clPassMode `" + flags.j2clPassMode + "' specified."))
		})
	}

	options.SetRemoveJ2clAsserts(flags.removeJ2cLAsserts)

	options.SetRenamePrefix(flags.renamePrefix)

	options.SetRenamePrefixNamespace(flags.renamePrefixNamespace)

	options.SetPreserveTypeAnnotations(flags.preserveTypeAnnotations)

	options.SetPreventLibraryInjection(!flags.injectLibraries)

	if len(flags.forceInjectLibraries) != 0 {
		options.SetForceLibraryInjection(flags.forceInjectLibraries)
	}

	options.SetRewritePolyfills(
		flags.rewritePolyfills &&
			options.GetLanguageIn().ToFeatureSet().Contains(parser.FeatureSet_ES2015))
	options.SetIsolatePolyfills(flags.isolatePolyfills)

	if flags.translationsFile != "" {
		jutil.Try(func() {
			options.SetMessageBundle(
				jscomp.NewXtbMessageBundle(
					jutil.NewFileInputStream(flags.translationsFile), flags.translationsProject))
		}, func(e jutil.IIOException) {
			panic(jutil.NewRuntimeExceptionWithMessageAndCause("Reading XTB file", e))
		})
	} else if jscomp.CompilationLevel_ADVANCED_OPTIMIZATIONS == level {
		// In SIMPLE or WHITESPACE mode, if the user hasn't specified a
		// translations file, they might reasonably try to write their own
		// implementation of goog.getMsg that makes the substitution at
		// run-time.
		//
		// In ADVANCED mode, goog.getMsg is going to be renamed anyway,
		// so we might as well inline it. But shut off the i18n warnings,
		// because the user didn't really ask for i18n.
		options.SetMessageBundle(jscomp.NewEmptyMessageBundle())
		options.SetWarningLevel(jscomp.DiagnosticGroups_MSG_CONVENTIONS, jscomp.CheckLevel_OFF)
	}

	options.SetConformanceConfigs(commandLineRunner_loadConformanceConfigs(flags.conformanceConfigs))

	options.SetPrintSourceAfterEachPass(flags.printSourceAfterEachPass)
	options.SetTracerMode(flags.tracerMode)
	options.SetStrictModeInput(flags.strictModeInput)
	options.SetEmitUseStrict(flags.emitUseStrict)
	options.SetSourceMapIncludeSourcesContent(flags.sourceMapIncludeSourcesContent)
	options.SetModuleResolutionMode(flags.moduleResolutionMode)
	options.SetBrowserResolverPrefixReplacements(
		commandLineRunner_immutableMapCopyOf(flags.browserResolverPrefixReplacements))

	if flags.packageJsonEntryNames != nil {
		jutil.Try(func() {
			packageJsonEntryNames := flags.GetPackageJsonEntryNames()
			options.SetPackageJsonEntryNames(packageJsonEntryNames)
		}, func(e args4j.ICmdLineException) {
			r.ReportError("ERROR - invalid package_json_entry_names format specified.")
		})
	}

	if !flags.renaming {
		options.SetVariableRenaming(jscomp.VariableRenamingPolicy_OFF)
		options.SetPropertyRenaming(jscomp.PropertyRenamingPolicy_OFF)
	}

	if flags.instrumentCodeParsed == jscomp.CompilerOptions_InstrumentOption_PRODUCTION &&
		flags.instrumentationMappingOutputFile == "" {
		panic(jscomp.NewFlagUsageException(
			"Expected --instrument_mapping_report to be set when " +
				"--instrument_for_coverage_option is set to Production"))
	}

	if flags.instrumentationMappingOutputFile != "" &&
		flags.instrumentCodeParsed != jscomp.CompilerOptions_InstrumentOption_PRODUCTION {
		panic(jscomp.NewFlagUsageException(
			"Expected --instrument_for_coverage_option to be passed with PRODUCTION " +
				"when --instrument_mapping_report is set"))
	}

	if flags.productionInstrumentationArrayName == "" &&
		flags.instrumentCodeParsed == jscomp.CompilerOptions_InstrumentOption_PRODUCTION {
		panic(jscomp.NewFlagUsageException(
			"Expected --production_instrumentation_array_name to be set when " +
				"--instrument_for_coverage_option is set to Production"))
	}
	options.SetInstrumentForCoverageOption(flags.instrumentCodeParsed)
	options.SetProductionInstrumentationArrayName(flags.productionInstrumentationArrayName)
	options.SetAllowDynamicImport(flags.allowDynamicImport)
	options.SetDynamicImportAlias(flags.dynamicImportAlias)
	options.SetAssumeStaticInheritanceIsNotUsed(flags.assumeStaticInheritanceIsNotUsed)
	options.SetCrossChunkCodeMotionNoStubMethods(flags.assumeNoPrototypeMethodEnumeration)

	if flags.chunkOutputType == jscomp.CompilerOptions_ChunkOutputType_ES_MODULES {
		if flags.renamePrefixNamespace != nil {
			panic(jscomp.NewFlagUsageException(
				"Expected --rename_prefix_namespace not to be specified when " +
					"--chunk_output_type is set to ES_MODULES."))
		}

		if flags.emitUseStrict {
			panic(jscomp.NewFlagUsageException(
				"Expected --emit_use_strict should not be specified when " +
					"--chunk_output_type is set to ES_MODULES."))
		}

		options.SetChunkOutputType(flags.chunkOutputType)
		options.SetEmitUseStrict(false)

		if level == jscomp.CompilationLevel_ADVANCED_OPTIMIZATIONS {
			options.SetExtractPrototypeMemberDeclarationsOfMode(
				jscomp.CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_CHUNK_TEMP)
		}
	}

	return options
}

// commandLineRunner_j2clPassModeValueOf is J2clPassMode.valueOf(String):
// the constant of that name, or IllegalArgumentException.
//
// Go-only: language-adapter — Enum.valueOf.
func commandLineRunner_j2clPassModeValueOf(name string) jscomp.CompilerOptions_J2clPassMode {
	for _, m := range jscomp.CompilerOptions_J2clPassMode_values() {
		if m.String() == name {
			return m
		}
	}
	panic(jutil.NewIllegalArgumentExceptionMsg(
		"No enum constant com.google.javascript.jscomp.CompilerOptions.J2clPassMode." + name))
}

// commandLineRunner_immutableMapCopyOf is ImmutableMap.copyOf(Map) of a
// LinkedHashMap<String, String>: its entries in insertion order.
//
// Go-only: language-adapter — args4j's MapOptionHandler stores String
// values boxed (*string).
func commandLineRunner_immutableMapCopyOf(m *jutil.LinkedHashMap[string, *string]) *jutil.LinkedHashMap[string, string] {
	result := jutil.NewImmutableMapBuilder[string, string]()
	for k, v := range m.All() {
		if v == nil {
			panic(jutil.NewNullPointerExceptionMsg("null value in entry: " + k + "=null"))
		}
		result.Put(k, *v)
	}
	return result.BuildOrThrow()
}

// Java: CommandLineRunner.createCompiler()
func (r *CommandLineRunner) CreateCompiler() *jscomp.Compiler {
	return jscomp.NewCompilerWithOutStream(r.GetErrorPrintStream())
}

// Java: CommandLineRunner.getBundler()
func (r *CommandLineRunner) GetBundler() *transpile.ClosureBundler {
	if r.bundler != nil {
		return r.bundler
	}

	var moduleRoots []string
	if len(r.flags.moduleRoot) != 0 {
		moduleRoots = slices.Clone(r.flags.moduleRoot)
	} else {
		moduleRoots = []string{jscomp.ModuleLoader_DEFAULT_FILENAME_PREFIX}
	}

	options := r.CreateOptions()
	r.bundler = transpile.NewClosureBundlerWithEs6ModuleTranspiler(
		transpile.Transpiler_NULL,
		transpile.NewBaseTranspiler(
			transpile.NewBaseTranspiler_CompilerSupplierWithModuleResolution(
				jscomp.CompilerOptions_LanguageMode_ECMASCRIPT_NEXT.ToFeatureSet().WithoutFeatures(parser.FeatureSet_Feature_MODULES),
				options.GetModuleResolutionMode(),
				moduleRoots,
				options.GetBrowserResolverPrefixReplacements()),
			/* runtimeLibraryName= */ ""))
	return r.bundler
}

// Java: CommandLineRunner.prepForBundleAndAppendTo(Appendable, CompilerInput, String)
func (r *CommandLineRunner) PrepForBundleAndAppendTo(out jutil.Appendable, input *jscomp.CompilerInput, content string) {
	r.GetBundler().WithPath(input.GetName()).AppendTo(out, input, content)
}

// Java: CommandLineRunner.appendRuntimeTo(Appendable)
func (r *CommandLineRunner) AppendRuntimeTo(out jutil.Appendable) {
	r.GetBundler().AppendRuntimeTo(out)
}

// Java: CommandLineRunner.createExterns(CompilerOptions)
func (r *CommandLineRunner) CreateExterns(options *jscomp.CompilerOptions) []*jscomp.SourceFile {
	externs := r.AbstractCommandLineRunnerBase.CreateExterns(options)
	if r.IsInTestMode() {
		return externs
	}
	builtinExterns := AbstractCommandLineRunner_GetBuiltinExterns(options.GetEnvironment())
	builtinExterns = append(builtinExterns, externs...)
	return builtinExterns
}

// Java: CommandLineRunner.loadConformanceConfigs(List<String>)
func commandLineRunner_loadConformanceConfigs(configPaths []string) []*jscomp.ConformanceConfig {
	configs := []*jscomp.ConformanceConfig{}

	for _, configPath := range configPaths {
		jutil.Try(func() {
			configs = append(configs, commandLineRunner_loadConformanceConfig(configPath))
		}, func(e jutil.IIOException) {
			panic(jutil.NewRuntimeExceptionWithMessageAndCause("Error loading conformance config", e))
		})
	}

	return configs
}

// Java: CommandLineRunner.loadConformanceConfig(String)
func commandLineRunner_loadConformanceConfig(configFile string) *jscomp.ConformanceConfig {
	textProto := jutil.Files_AsCharSource(jutil.NewFile(configFile), jutil.StandardCharsets_UTF_8).Read()

	builder := jscomp.ConformanceConfig_NewBuilder()

	// Stripping the BOM if present.
	textProto = strings.TrimPrefix(textProto, "\uFEFF")

	jutil.Try(func() {
		jscomp.TextFormat_Merge(textProto, builder)
	}, func(e jutil.IException) {
		panic(jutil.NewRuntimeExceptionWithCause(e))
	})
	return builder.Build()
}

// Deprecated: use AbstractCommandLineRunner_GetBuiltinExterns(
// CompilerOptions_Environment_BROWSER).
//
// Java: CommandLineRunner.getDefaultExterns()
func CommandLineRunner_GetDefaultExterns() []*jscomp.SourceFile {
	return AbstractCommandLineRunner_GetBuiltinExterns(jscomp.CompilerOptions_Environment_BROWSER)
}

// CommandLineRunner_FindJsFiles returns all the JavaScript files from the
// set of patterns. The patterns support globs, such as '*.js' for all JS
// files in a directory and '**.js' for all JS files within the directory and
// sub-directories.
//
// Java: CommandLineRunner.findJsFiles(Collection<String>)
func CommandLineRunner_FindJsFiles(patterns []string) []string {
	return commandLineRunner_findJsFilesSorted(patterns, false)
}

// jsInputsCollector tracks normalized-absolute-path -> original-path entries
// while preserving insertion order (or sorting keys alphabetically on output).
type jsInputsCollector struct {
	byAbs map[string]string
	order []string
}

func (c *jsInputsCollector) put(absPath, origPath string) {
	if c.byAbs == nil {
		c.byAbs = make(map[string]string)
	}
	if _, exists := c.byAbs[absPath]; !exists {
		c.order = append(c.order, absPath)
	}
	c.byAbs[absPath] = origPath
}

func (c *jsInputsCollector) remove(absPath string) {
	if _, exists := c.byAbs[absPath]; !exists {
		return
	}
	delete(c.byAbs, absPath)
	if i := slices.Index(c.order, absPath); i >= 0 {
		c.order = slices.Delete(c.order, i, i+1)
	}
}

// commandLineRunner_findJsFilesSorted returns all the JavaScript files from
// the set of patterns; sortAlphabetically is whether the output filenames
// should be in alphabetical order. It returns the list of JS filenames found
// by expanding the patterns.
//
// Java: CommandLineRunner.findJsFiles(Collection<String>, boolean)
func commandLineRunner_findJsFilesSorted(patterns []string, sortAlphabetically bool) []string {
	if len(patterns) == 0 {
		return nil
	}
	// A map from normalized absolute paths to original paths. We need to return original paths to
	// support allowlist files that depend on them.
	var allJsInputs jsInputsCollector
	var excludes map[string]struct{}
	var wd *jutil.Path
	for _, pattern := range patterns {
		if !strings.Contains(pattern, "*") && !strings.HasPrefix(pattern, "!") {
			if jutil.File_IsDirectoryPath(pattern) {
				commandLineRunner_matchPaths(jutil.NewFileWithParent(jutil.NewFile(pattern), "**.js").String(), &allJsInputs, &excludes, &wd)
			} else {
				original := jutil.Path_Of(pattern)
				norm := original.Normalize()
				if !norm.IsAbsolute() && wd == nil {
					wd = jutil.Path_WorkingDir()
				}
				pathStringAbsolute := norm.ToAbsolutePathWithWd(wd).String()
				if _, excluded := excludes[pathStringAbsolute]; !excluded {
					allJsInputs.put(pathStringAbsolute, original.String())
				}
			}
		} else {
			commandLineRunner_matchPaths(pattern, &allJsInputs, &excludes, &wd)
		}
	}

	if len(allJsInputs.order) == 0 {
		return nil
	}
	if sortAlphabetically {
		slices.SortFunc(allJsInputs.order, jutil.StringCompareTo)
	}
	result := make([]string, len(allJsInputs.order))
	for i, absPath := range allJsInputs.order {
		result[i] = allJsInputs.byAbs[absPath]
	}
	return result
}

// Java: CommandLineRunner.matchPaths(String, Map<String, String>, Set<String>)
func commandLineRunner_matchPaths(
	pattern string,
	allJsInputs *jsInputsCollector,
	excludes *map[string]struct{},
	wd **jutil.Path,
) {
	remove := strings.HasPrefix(pattern, "!")
	if remove {
		pattern = pattern[1:]
	}

	separator := jutil.File_separator

	// Split the pattern into two pieces: the globbing part
	// and the non-globbing prefix.
	// Go-only: language-adapter — Guava Splitter.on / Joiner.on with the
	// ASCII separator are strings.Split / strings.Join.
	patternParts := strings.Split(pattern, separator)
	prefix := "."
	for i := 0; i < len(patternParts); i++ {
		if strings.Contains(patternParts[i], "*") {
			if i > 0 {
				prefix = strings.Join(patternParts[:i], separator)
				pattern = strings.Join(patternParts[i:], separator)
			}
			break
		}
	}

	matcher := compileGlobMatcher(prefix + separator + pattern)
	visitFile := func(p *jutil.Path) {
		var norm *jutil.Path
		matched := matcher.MatchString(p.String())
		if !matched {
			norm = p.Normalize()
			matched = matcher.MatchString(norm.String())
		}
		if !matched {
			return
		}
		if norm == nil {
			norm = p.Normalize()
		}
		if !norm.IsAbsolute() && *wd == nil {
			*wd = jutil.Path_WorkingDir()
		}
		pathStringAbsolute := norm.ToAbsolutePathWithWd(*wd).String()
		if remove {
			if *excludes == nil {
				*excludes = make(map[string]struct{})
			}
			(*excludes)[pathStringAbsolute] = struct{}{}
			allJsInputs.remove(pathStringAbsolute)
		} else if _, excluded := (*excludes)[pathStringAbsolute]; !excluded {
			allJsInputs.put(pathStringAbsolute, p.String())
		}
	}

	startPath := jutil.Path_Of(prefix)
	startStr := startPath.String()
	st, err := os.Stat(startStr)
	if err != nil {
		// Broken symlink: visit the link itself as a file; nonexistent start: skip.
		if lst, lerr := os.Lstat(startStr); lerr == nil && !lst.IsDir() {
			visitFile(startPath)
		}
		return
	}
	if !st.IsDir() {
		visitFile(startPath)
		return
	}
	evalRoot, err := filepath.EvalSymlinks(startStr)
	if err != nil {
		return
	}
	commandLineRunner_walkDir(evalRoot, startPath, nil, visitFile)
}

func commandLineRunner_walkDir(
	evalRoot string,
	logicalRoot *jutil.Path,
	ancestors []os.FileInfo,
	visitFile func(*jutil.Path),
) {
	_ = filepath.WalkDir(evalRoot, func(walkPath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || walkPath == evalRoot {
			return nil
		}
		var rel string
		switch evalRoot {
		case ".":
			rel = walkPath
		case "/":
			rel = walkPath[1:]
		default:
			rel = walkPath[len(evalRoot)+1:]
		}
		if d.IsDir() {
			return nil
		}
		logicalPath := logicalRoot.ResolveString(rel)
		if d.Type()&fs.ModeSymlink != 0 {
			targetInfo, err := os.Stat(walkPath)
			if err == nil && targetInfo.IsDir() {
				for _, anc := range ancestors {
					if os.SameFile(anc, targetInfo) {
						return nil
					}
				}
				nextAncestors := slices.Clone(ancestors)
				for cur := filepath.Dir(walkPath); ; {
					dirInfo, derr := os.Stat(cur)
					if derr != nil || os.SameFile(dirInfo, targetInfo) {
						return nil
					}
					nextAncestors = append(nextAncestors, dirInfo)
					if cur == evalRoot {
						break
					}
					parent := filepath.Dir(cur)
					if parent == cur {
						break
					}
					cur = parent
				}
				evalTarget, err := filepath.EvalSymlinks(walkPath)
				if err != nil {
					return nil
				}
				commandLineRunner_walkDir(evalTarget, logicalPath, nextAncestors, visitFile)
				return nil
			}
		}
		visitFile(logicalPath)
		return nil
	})
}

// ShouldRunCompiler returns whether the configuration is valid and
// specifies to run the compiler.
//
// Java: CommandLineRunner.shouldRunCompiler()
func (r *CommandLineRunner) ShouldRunCompiler() bool {
	return r.runCompiler
}

// HasErrors returns whether the configuration has errors.
//
// Java: CommandLineRunner.hasErrors()
func (r *CommandLineRunner) HasErrors() bool {
	return r.errors
}

// Java: CommandLineRunner.phaseLogger
var commandLineRunner_phaseLogger = jutil.Logger_GetLogger("com.google.javascript.jscomp.PhaseOptimizer")

// CommandLineRunner_Main runs the Compiler. Exits cleanly in the event of an
// error.
//
// Java: CommandLineRunner.main(String[])
func CommandLineRunner_Main(args []string) {
	// disable any logging messages that can interfere with standard error reporting
	if commandLineRunner_phaseLogger != nil {
		commandLineRunner_phaseLogger.SetLevel(jutil.Level_OFF)
	}
	runner := NewCommandLineRunner(args)
	if runner.ShouldRunCompiler() {
		runner.Run()
	}
	if runner.HasErrors() {
		jutil.System_Exit(-1)
	}
}

func asciiLower(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			b := []byte(s)
			for ; i < len(b); i++ {
				if b[i] >= 'A' && b[i] <= 'Z' {
					b[i] += 'a' - 'A'
				}
			}
			return string(b)
		}
	}
	return s
}

func asciiUpper(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] >= 'a' && s[i] <= 'z' {
			b := []byte(s)
			for ; i < len(b); i++ {
				if b[i] >= 'a' && b[i] <= 'z' {
					b[i] -= 'a' - 'A'
				}
			}
			return string(b)
		}
	}
	return s
}

func compileGlobMatcher(glob string) *regexp.Regexp {
	var sb strings.Builder
	sb.WriteString(`(?s)^`)
	groupOpen := false
	n := len(glob)
	for i := 0; i < n; i++ {
		switch glob[i] {
		case '\\':
			if i+1 == n {
				panic(jutil.NewPatternSyntaxException("No character to escape", glob, i))
			}
			i++
			_, w := utf8.DecodeRuneInString(glob[i:])
			sb.WriteString(regexp.QuoteMeta(glob[i : i+w]))
			i += w - 1
		case '*':
			if i+1 < n && glob[i+1] == '*' {
				sb.WriteString(`.*`)
				i++
			} else {
				sb.WriteString(`[^/]*`)
			}
		case '?':
			sb.WriteString(`[^/]`)
		case '[':
			i = globBracket(glob, i, &sb)
		case '{':
			if groupOpen {
				panic(jutil.NewPatternSyntaxException("Cannot nest groups", glob, i))
			}
			groupOpen = true
			sb.WriteString(`(?:`)
		case '}':
			if groupOpen {
				groupOpen = false
				sb.WriteString(`)`)
			} else {
				sb.WriteString(`\}`)
			}
		case ',':
			if groupOpen {
				sb.WriteString(`|`)
			} else {
				sb.WriteString(`,`)
			}
		default:
			_, w := utf8.DecodeRuneInString(glob[i:])
			sb.WriteString(regexp.QuoteMeta(glob[i : i+w]))
			i += w - 1
		}
	}
	if groupOpen {
		panic(jutil.NewPatternSyntaxException("Missing '}", glob, n-1))
	}
	sb.WriteString(`$`)
	return regexp.MustCompile(sb.String())
}

func globBracket(glob string, start int, sb *strings.Builder) int {
	n := len(glob)
	i := start + 1
	negate := false
	if i < n && glob[i] == '!' {
		negate = true
		i++
	}
	type rng struct{ lo, hi rune }
	var ranges []rng
	for ; i < n && glob[i] != ']'; i++ {
		loIndex := i
		lo, w := utf8.DecodeRuneInString(glob[i:])
		i += w - 1
		if lo == '/' {
			panic(jutil.NewPatternSyntaxException("Explicit 'name separator' in class", glob, i))
		}
		hi := lo
		if i+2 < n && glob[i+1] == '-' && glob[i+2] != ']' {
			i += 2
			var hw int
			hi, hw = utf8.DecodeRuneInString(glob[i:])
			i += hw - 1
			if hi < lo {
				panic(jutil.NewPatternSyntaxException("Invalid range", glob, loIndex))
			}
		}
		ranges = append(ranges, rng{lo, hi})
	}
	if i >= n {
		panic(jutil.NewPatternSyntaxException("Missing ']", glob, n-1))
	}
	if len(ranges) == 0 {
		panic(jutil.NewPatternSyntaxException("Unclosed character class", glob, i))
	}
	sb.WriteByte('[')
	if negate {
		sb.WriteString(`^/`)
	}
	writeClassRune := func(r rune) {
		switch r {
		case '\\', ']', '[', '^', '-':
			sb.WriteByte('\\')
		}
		sb.WriteRune(r)
	}
	for _, r := range ranges {
		for _, part := range [2]rng{{r.lo, min(r.hi, '/'-1)}, {max(r.lo, '/'+1), r.hi}} {
			if part.lo > part.hi {
				continue
			}
			writeClassRune(part.lo)
			if part.hi != part.lo {
				sb.WriteByte('-')
				writeClassRune(part.hi)
			}
		}
	}
	sb.WriteByte(']')
	return i
}
