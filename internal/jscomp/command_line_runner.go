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

package jscomp

import (
	"io"
	"math"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/args4j"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
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

	bundler *ClosureBundler // nil is null
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

// The args4j classes of the enum types of Flags' fields (EnumOptionHandler
// reads the constants of the field's class).
//
// Go-only: language-adapter — annotations (see package args4j).
var (
	commandLineRunner_CompilerOptions_DevMode_class                                       = args4j.EnumClassOf("com.google.javascript.jscomp.CompilerOptions$DevMode", CompilerOptions_DevMode_values())
	commandLineRunner_CompilerOptions_SegmentOfCompilationToRun_class                     = args4j.EnumClassOf("com.google.javascript.jscomp.CompilerOptions$SegmentOfCompilationToRun", CompilerOptions_SegmentOfCompilationToRun_values())
	commandLineRunner_CompilerOptions_IsolationMode_class                                 = args4j.EnumClassOf("com.google.javascript.jscomp.CompilerOptions$IsolationMode", CompilerOptions_IsolationMode_values())
	commandLineRunner_SourceMap_Format_class                                              = args4j.EnumClassOf("com.google.javascript.jscomp.SourceMap$Format", SourceMap_Format_values())
	commandLineRunner_CompilerOptions_IncrementalCheckMode_class                          = args4j.EnumClassOf("com.google.javascript.jscomp.CompilerOptions$IncrementalCheckMode", CompilerOptions_IncrementalCheckMode_values())
	commandLineRunner_WarningLevel_class                                                  = args4j.EnumClassOf("com.google.javascript.jscomp.WarningLevel", WarningLevel_values())
	commandLineRunner_CompilerOptions_TracerMode_class                                    = args4j.EnumClassOf("com.google.javascript.jscomp.CompilerOptions$TracerMode", CompilerOptions_TracerMode_values())
	commandLineRunner_CompilerOptions_Environment_class                                   = args4j.EnumClassOf("com.google.javascript.jscomp.CompilerOptions$Environment", CompilerOptions_Environment_values())
	commandLineRunner_CompilerOptions_JsonStreamMode_class                                = args4j.EnumClassOf("com.google.javascript.jscomp.CompilerOptions$JsonStreamMode", CompilerOptions_JsonStreamMode_values())
	commandLineRunner_DependencyOptions_DependencyMode_class                              = args4j.EnumClassOf("com.google.javascript.jscomp.DependencyOptions$DependencyMode", DependencyOptions_DependencyMode_values())
	commandLineRunner_ModuleLoader_ResolutionMode_class                                   = args4j.EnumClassOf("com.google.javascript.jscomp.deps.ModuleLoader$ResolutionMode", ModuleLoader_ResolutionMode_values())
	commandLineRunner_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_class = args4j.EnumClassOf(
		"com.google.javascript.jscomp.AbstractCommandLineRunner$CommandLineConfig$ErrorFormatOption",
		AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_values())
	commandLineRunner_CompilerOptions_ChunkOutputType_class = args4j.EnumClassOf("com.google.javascript.jscomp.CompilerOptions$ChunkOutputType", CompilerOptions_ChunkOutputType_values())
	commandLineRunner_FormattingOption_class                = args4j.EnumClassOf("com.google.javascript.jscomp.CommandLineRunner$FormattingOption", CommandLineRunner_FormattingOption_values())
)

// I don't really care about unchecked warnings in this class.
//
// Java: CommandLineRunner.Flags
type commandLineRunner_Flags struct {
	browserFeaturesetYear              int32
	displayHelp                        bool
	printTree                          bool
	printTreeJson                      bool
	printAst                           bool
	emitUseStrict                      bool
	strictModeInput                    bool
	jscompDevMode                      CompilerOptions_DevMode
	loggingLevel                       string
	externs                            []string
	js                                 []string
	unusedJsZip                        []string
	jsOutputFile                       string
	chunk                              []string
	filenameToSaveTo                   *string
	filenameToRestoreFrom              *string
	segmentOfCompilationToRun          *CompilerOptions_SegmentOfCompilationToRun
	variableMapOutputFile              string
	instrumentationMappingOutputFile   string
	createNameMapFiles                 bool
	sourceMapIncludeSourcesContent     bool
	propertyMapOutputFile              string
	thirdParty                         bool
	summaryDetailLevel                 int
	isolationMode                      CompilerOptions_IsolationMode
	outputWrapper                      string
	outputWrapperFile                  string
	chunkWrapper                       []string
	chunkOutputPathPrefix              string
	createSourceMap                    string
	sourceMapFormat                    SourceMap_Format
	sourceMapLocationMapping           []string
	sourceMapInputs                    []string
	parseInlineSourceMaps              *bool
	applyInputSourceMaps               bool
	jscompError                        []string
	jscompWarning                      []string
	jscompOff                          []string
	define                             []string
	charset                            string
	compilationLevel                   string
	numParallelThreads                 int
	fastMode                           bool
	checksOnly                         bool
	incrementalCheckMode               CompilerOptions_IncrementalCheckMode
	continueAfterErrors                bool
	useTypesForOptimization            bool
	assumeFunctionWrapper              bool
	warningLevel                       WarningLevel
	debug                              bool
	typedAstOutputFile                 *string
	generateExports                    bool
	exportLocalPropertyDefinitions     bool
	formatting                         []CommandLineRunner_FormattingOption
	processCommonJsModules             bool
	moduleRoot                         []string
	processClosurePrimitives           bool
	angularPass                        bool
	polymerVersion                     *int32
	chromePass                         bool
	j2clPassMode                       string
	removeJ2cLAsserts                  bool
	outputManifest                     string
	outputChunkDependencies            string
	languageIn                         string
	languageOut                        string
	version                            bool
	translationsFile                   string
	translationsProject                *string
	flagFiles                          []string
	warningsAllowlistFile              string
	hideWarningsFor                    []string
	extraAnnotationName                []string
	tracerMode                         CompilerOptions_TracerMode
	renamePrefix                       *string
	renamePrefixNamespace              *string
	conformanceConfigs                 []string
	environment                        CompilerOptions_Environment
	jsonStreamMode                     CompilerOptions_JsonStreamMode
	preserveTypeAnnotations            bool
	injectLibraries                    bool
	forceInjectLibraries               []string
	dependencyMode                     *DependencyOptions_DependencyMode
	entryPoint                         []string
	rewritePolyfills                   bool
	isolatePolyfills                   bool
	printSourceAfterEachPass           bool
	moduleResolutionMode               ModuleLoader_ResolutionMode
	browserResolverPrefixReplacements  *jutil.LinkedHashMap[string, *string]
	packageJsonEntryNames              *string
	errorFormat                        AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption
	renaming                           bool
	helpMarkdown                       bool
	instrumentForCoverageOption        string
	productionInstrumentationArrayName string
	chunkOutputType                    CompilerOptions_ChunkOutputType
	allowDynamicImport                 bool
	dynamicImportAlias                 *string
	assumeStaticInheritanceIsNotUsed   bool
	assumeNoPrototypeMethodEnumeration bool
	variableMapInputFile               string
	propertyMapInputFile               string
	expectedDiagnostics                []string

	arguments []string

	compilationLevelParsed *CompilationLevel // nil is null

	instrumentCodeParsed CompilerOptions_InstrumentOption

	// languageOutIsDefault is whether languageOut still holds its initial
	// String object.
	//
	// Go-only: language-adapter — initConfigFromFlags compares languageOut
	// with the literal "ECMASCRIPT_NEXT" by reference (!=); only the field's
	// initial value is that interned literal, and every value args4j sets is
	// a String parsed from the arguments, a distinct object.
	languageOutIsDefault bool

	parser *args4j.CmdLineParser

	maxLineLength int
}

// Some clients run a few copies of the compiler through CommandLineRunner
// on parallel threads (thankfully, with the same flags),
// so the access to these lists should be synchronized.
//
// Java: CommandLineRunner.Flags.guardLevels
var commandLineRunner_Flags_guardLevels []*AbstractCommandLineRunner_FlagEntry[CheckLevel]

// Java: CommandLineRunner.Flags.mixedJsSources
var commandLineRunner_Flags_mixedJsSources []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]

// Java: CommandLineRunner.Flags()
func newCommandLineRunner_Flags() *commandLineRunner_Flags {
	f := &commandLineRunner_Flags{
		browserFeaturesetYear: 0,
		displayHelp:           false,
		printTree:             false,
		printTreeJson:         false,
		printAst:              false,
		emitUseStrict:         false,
		strictModeInput:       true,
		jscompDevMode:         CompilerOptions_DevMode_OFF,
		loggingLevel:          jutil.Level_WARNING.GetName(),
		externs:               []string{},
		js:                    []string{},
		unusedJsZip:           nil,
		jsOutputFile:          "",
		chunk:                 []string{},
		filenameToSaveTo:      nil,
		filenameToRestoreFrom: nil,
		segmentOfCompilationToRun: func() *CompilerOptions_SegmentOfCompilationToRun {
			v := CompilerOptions_SegmentOfCompilationToRun_ENTIRE_COMPILATION
			return &v
		}(),
		variableMapOutputFile:              "",
		instrumentationMappingOutputFile:   "",
		createNameMapFiles:                 false,
		sourceMapIncludeSourcesContent:     false,
		propertyMapOutputFile:              "",
		thirdParty:                         false,
		summaryDetailLevel:                 1,
		isolationMode:                      CompilerOptions_IsolationMode_NONE,
		outputWrapper:                      "",
		outputWrapperFile:                  "",
		chunkWrapper:                       []string{},
		chunkOutputPathPrefix:              "./",
		createSourceMap:                    "",
		sourceMapFormat:                    SourceMap_Format_DEFAULT,
		sourceMapLocationMapping:           []string{},
		sourceMapInputs:                    []string{},
		parseInlineSourceMaps:              func() *bool { b := true; return &b }(),
		applyInputSourceMaps:               true,
		jscompError:                        []string{},
		jscompWarning:                      []string{},
		jscompOff:                          []string{},
		define:                             []string{},
		charset:                            "",
		compilationLevel:                   "SIMPLE",
		numParallelThreads:                 1,
		checksOnly:                         false,
		incrementalCheckMode:               CompilerOptions_IncrementalCheckMode_OFF,
		continueAfterErrors:                false,
		useTypesForOptimization:            true,
		assumeFunctionWrapper:              false,
		warningLevel:                       WarningLevel_DEFAULT,
		debug:                              false,
		typedAstOutputFile:                 nil,
		generateExports:                    true,
		exportLocalPropertyDefinitions:     true,
		formatting:                         []CommandLineRunner_FormattingOption{},
		processCommonJsModules:             false,
		moduleRoot:                         []string{},
		processClosurePrimitives:           true,
		angularPass:                        false,
		polymerVersion:                     nil,
		chromePass:                         false,
		j2clPassMode:                       "AUTO",
		removeJ2cLAsserts:                  true,
		outputManifest:                     "",
		outputChunkDependencies:            "",
		languageIn:                         "STABLE",
		languageOut:                        "ECMASCRIPT_NEXT",
		version:                            false,
		translationsFile:                   "",
		translationsProject:                nil,
		flagFiles:                          []string{},
		warningsAllowlistFile:              "",
		hideWarningsFor:                    []string{},
		extraAnnotationName:                []string{},
		tracerMode:                         CompilerOptions_TracerMode_OFF,
		renamePrefix:                       nil,
		renamePrefixNamespace:              nil,
		conformanceConfigs:                 []string{},
		environment:                        CompilerOptions_Environment_BROWSER,
		jsonStreamMode:                     CompilerOptions_JsonStreamMode_NONE,
		preserveTypeAnnotations:            false,
		injectLibraries:                    true,
		forceInjectLibraries:               []string{},
		dependencyMode:                     nil,
		entryPoint:                         []string{},
		rewritePolyfills:                   true,
		isolatePolyfills:                   false,
		printSourceAfterEachPass:           false,
		moduleResolutionMode:               ModuleLoader_ResolutionMode_BROWSER,
		browserResolverPrefixReplacements:  jutil.NewLinkedHashMap[string, *string](),
		packageJsonEntryNames:              nil,
		errorFormat:                        AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_STANDARD,
		renaming:                           true,
		helpMarkdown:                       false,
		instrumentForCoverageOption:        "NONE",
		productionInstrumentationArrayName: "",
		chunkOutputType:                    CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE,
		allowDynamicImport:                 true,
		dynamicImportAlias:                 nil,
		assumeStaticInheritanceIsNotUsed:   true,
		assumeNoPrototypeMethodEnumeration: false,
		variableMapInputFile:               "",
		propertyMapInputFile:               "",
		expectedDiagnostics:                []string{},
		compilationLevelParsed:             nil,
		instrumentCodeParsed:               CompilerOptions_InstrumentOption_NONE,
		languageOutIsDefault:               true,
		maxLineLength:                      80,
	}
	f.parser = args4j.NewCmdLineParser(f)
	return f
}

// Args4jFields lists Flags' annotated fields in declaration order.
//
// Go-only: language-adapter — annotations (see package args4j).
func (f *commandLineRunner_Flags) Args4jFields() []args4j.AnnotatedField {
	const c = "com.google.javascript.jscomp.CommandLineRunner$Flags"
	return []args4j.AnnotatedField{
		{Option: &args4j.Option{Name: "--browser_featureset_year", Usage: "shortcut for defining goog.FEATURESET_YEAR=YYYY. The minimum valid value of the browser year is 2012"},
			Field: args4j.IntegerField(c, "browserFeaturesetYear", &f.browserFeaturesetYear)},
		{Option: &args4j.Option{Name: "--help", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Displays this message on stdout and exit"},
			Field: args4j.BooleanField(c, "displayHelp", &f.displayHelp)},
		{Option: &args4j.Option{Name: "--print_tree", Hidden: true, Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Prints out the parse tree and exits"},
			Field: args4j.BooleanField(c, "printTree", &f.printTree)},
		{Option: &args4j.Option{Name: "--print_tree_json", Hidden: true, Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Prints out the parse tree as json and exits"},
			Field: args4j.BooleanField(c, "printTreeJson", &f.printTreeJson)},
		{Option: &args4j.Option{Name: "--print_ast", Hidden: true, Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Prints a dot file describing the internal abstract syntax tree and exits"},
			Field: args4j.BooleanField(c, "printAst", &f.printAst)},
		{Option: &args4j.Option{Name: "--emit_use_strict", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Start output with \"'use strict';\"."},
			Field: args4j.BooleanField(c, "emitUseStrict", &f.emitUseStrict)},
		{Option: &args4j.Option{Name: "--strict_mode_input", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Assume input sources are to run in strict mode."},
			Field: args4j.BooleanField(c, "strictModeInput", &f.strictModeInput)},
		{Option: &args4j.Option{Name: "--jscomp_dev_mode", Hidden: true},
			Field: args4j.EnumField(c, "jscompDevMode", commandLineRunner_CompilerOptions_DevMode_class, &f.jscompDevMode)},
		{Option: &args4j.Option{Name: "--logging_level", Hidden: true, Usage: "The logging level (standard java.util.logging.Level values) for Compiler progress. Does not control errors or warnings for the JavaScript code under compilation"},
			Field: args4j.StringField(c, "loggingLevel", &f.loggingLevel)},
		{Option: &args4j.Option{Name: "--externs", Usage: "The file containing JavaScript externs. You may specify multiple"},
			Field: args4j.ListField(c, "externs", args4j.Class_String, &f.externs)},
		{Option: &args4j.Option{Name: "--js", Handler: newCommandLineRunner_Flags_JsOptionHandler, Usage: "The JavaScript filename. You may specify multiple. The flag name is optional, because args are interpreted as files by default. You may also use minimatch-style glob patterns. For example, use --js='**.js' --js='!**_test.js' to recursively include all js files that do not end in _test.js"},
			Field: args4j.ListField(c, "js", args4j.Class_String, &f.js)},
		{Option: &args4j.Option{Name: "--jszip", Hidden: true, Handler: newCommandLineRunner_Flags_JsZipOptionHandler, Usage: "The JavaScript zip filename. You may specify multiple."},
			Field: args4j.ListField(c, "unusedJsZip", args4j.Class_String, &f.unusedJsZip)},
		{Option: &args4j.Option{Name: "--js_output_file", Usage: "Primary output filename. If not specified, output is written to stdout"},
			Field: args4j.StringField(c, "jsOutputFile", &f.jsOutputFile)},
		{Option: &args4j.Option{Name: "--chunk", Usage: "A JavaScript chunk specification. The format is <name>:<num-js-files>[:[<dep>,...][:]]]. Chunk names must be unique. Each dep is the name of a chunk that this chunk depends on. Chunks must be listed in dependency order, and JS source files must be listed in the corresponding order. Where --chunk flags occur in relation to --js flags is unimportant. <num-js-files> may be set to 'auto' for the first chunk if it has no dependencies. Provide the value 'auto' to trigger chunk creation from CommonJS modules."},
			Field: args4j.ListField(c, "chunk", args4j.Class_String, &f.chunk)},
		{Option: &args4j.Option{Name: "--filename_to_save_to", Hidden: true, Usage: "Filename to save state so that the compilation can be resumed later."},
			Field: args4j.NullableStringField(c, "filenameToSaveTo", &f.filenameToSaveTo)},
		{Option: &args4j.Option{Name: "--filename_to_restore_from", Hidden: true, Usage: "Filename where a compilation state was previously saved."},
			Field: args4j.NullableStringField(c, "filenameToRestoreFrom", &f.filenameToRestoreFrom)},
		{Option: &args4j.Option{Name: "--segment_of_compilation_to_run", Hidden: true, Usage: "Which segment of the compilation to run."},
			Field: args4j.NullableEnumField(c, "segmentOfCompilationToRun", commandLineRunner_CompilerOptions_SegmentOfCompilationToRun_class, &f.segmentOfCompilationToRun)},
		{Option: &args4j.Option{Name: "--variable_renaming_report", Usage: "File where the serialized version of the variable renaming map produced should be saved"},
			Field: args4j.StringField(c, "variableMapOutputFile", &f.variableMapOutputFile)},
		{Option: &args4j.Option{Name: "--instrument_mapping_report", Usage: "File where the encoded parameters created by Production Instrumentation are mapped to their pre-encoded values. The %outname% placeholder will expand to the name of the output file that the source map corresponds to. Must be used in tandem with --instrument_for_coverage_option=PRODUCTION"},
			Field: args4j.StringField(c, "instrumentationMappingOutputFile", &f.instrumentationMappingOutputFile)},
		{Option: &args4j.Option{Name: "--create_renaming_reports", Hidden: true, Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "If true, variable renaming and property renaming report files will be produced as {binary name}_vars_renaming_report.out and {binary name}_props_renaming_report.out. Note that this flag cannot be used in conjunction with either variable_renaming_report or property_renaming_report"},
			Field: args4j.BooleanField(c, "createNameMapFiles", &f.createNameMapFiles)},
		{Option: &args4j.Option{Name: "--source_map_include_content", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Includes sources content into source map. Greatly increases the size of source maps but offers greater portability"},
			Field: args4j.BooleanField(c, "sourceMapIncludeSourcesContent", &f.sourceMapIncludeSourcesContent)},
		{Option: &args4j.Option{Name: "--property_renaming_report", Usage: "File where the serialized version of the property renaming map produced should be saved"},
			Field: args4j.StringField(c, "propertyMapOutputFile", &f.propertyMapOutputFile)},
		{Option: &args4j.Option{Name: "--third_party", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Check source validity but do not enforce Closure style rules and conventions"},
			Field: args4j.BooleanField(c, "thirdParty", &f.thirdParty)},
		{Option: &args4j.Option{Name: "--summary_detail_level", Hidden: true, Usage: "Controls how detailed the compilation summary is. Values: 0 (never print summary), 1 (print summary only if there are errors or warnings), 2 (print summary if the 'checkTypes' diagnostic  group is enabled, see --jscomp_warning), 3 (always print summary). The default level is 1"},
			Field: args4j.IntField(c, "summaryDetailLevel", &f.summaryDetailLevel)},
		{Option: &args4j.Option{Name: "--isolation_mode", Usage: "If set to IIFE the compiler output will follow the form:\n  (function(){%output%}).call(this);\nOptions: NONE, IIFE"},
			Field: args4j.EnumField(c, "isolationMode", commandLineRunner_CompilerOptions_IsolationMode_class, &f.isolationMode)},
		{Option: &args4j.Option{Name: "--output_wrapper", Usage: "Interpolate output into this string at the place denoted by the marker token %output%. Use marker token %output|jsstring% to do js string escaping on the output. Consider using the --isolation_mode flag instead."},
			Field: args4j.StringField(c, "outputWrapper", &f.outputWrapper)},
		{Option: &args4j.Option{Name: "--output_wrapper_file", Usage: "Loads the specified file and passes the file contents to the --output_wrapper flag, replacing the value if it exists. This is useful if you want special characters like newline in the wrapper."},
			Field: args4j.StringField(c, "outputWrapperFile", &f.outputWrapperFile)},
		{Option: &args4j.Option{Name: "--chunk_wrapper", Usage: "An output wrapper for a JavaScript chunk (optional). The format is <name>:<wrapper>. The chunk name must correspond with a chunk specified using --chunk. The wrapper must contain %s as the code placeholder. Alternately, %output% can be used in place of %s. %n% can be used to represent a newline. The %basename% placeholder can also be used to substitute the base name of the chunk output file."},
			Field: args4j.ListField(c, "chunkWrapper", args4j.Class_String, &f.chunkWrapper)},
		{Option: &args4j.Option{Name: "--chunk_output_path_prefix", Usage: "Prefix for filenames of compiled JS chunks. <chunk-name>.js will be appended to this prefix. Directories will be created as needed. Use with --chunk"},
			Field: args4j.StringField(c, "chunkOutputPathPrefix", &f.chunkOutputPathPrefix)},
		{Option: &args4j.Option{Name: "--create_source_map", Usage: "If specified, a source map file mapping the generated source files back to the original source file will be output to the specified path. The %outname% placeholder will expand to the name of the output file that the source map corresponds to."},
			Field: args4j.StringField(c, "createSourceMap", &f.createSourceMap)},
		{Option: &args4j.Option{Name: "--source_map_format", Hidden: true, Usage: "The source map format to produce. Options are V3 and DEFAULT, which are equivalent."},
			Field: args4j.EnumField(c, "sourceMapFormat", commandLineRunner_SourceMap_Format_class, &f.sourceMapFormat)},
		{Option: &args4j.Option{Name: "--source_map_location_mapping", Usage: "Source map location mapping separated by a '|' (i.e. filesystem-path|webserver-path)"},
			Field: args4j.ListField(c, "sourceMapLocationMapping", args4j.Class_String, &f.sourceMapLocationMapping)},
		{Option: &args4j.Option{Name: "--source_map_input", Usage: "Source map locations for input files, separated by a '|', (i.e. input-file-path|input-source-map)"},
			Field: args4j.ListField(c, "sourceMapInputs", args4j.Class_String, &f.sourceMapInputs)},
		{Option: &args4j.Option{Name: "--parse_inline_source_maps", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Parse inline source maps (//# sourceMappingURL=data:...)"},
			Field: args4j.NullableBooleanField(c, "parseInlineSourceMaps", &f.parseInlineSourceMaps)},
		{Option: &args4j.Option{Name: "--apply_input_source_maps", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Apply input source maps to the output source map, i.e. have the result map back to original inputs.  Input sourcemaps can be located in 2 ways:\n 1) by the//# sourceMappingURL=<url>. \n 2) using the--source_map_location_mapping flag.\nsourceMappingURL=<url> can read both paths and inline Base64 encoded sourcemaps. For inline Base64 encoded sourcemaps, see --parse_inline_source_maps."},
			Field: args4j.BooleanField(c, "applyInputSourceMaps", &f.applyInputSourceMaps)},
		{Option: &args4j.Option{Name: "--jscomp_error", Handler: newCommandLineRunner_Flags_WarningGuardErrorOptionHandler, Usage: "Make the named class of warnings an error. Must be one of the error group items. '*' adds all supported."},
			Field: args4j.ListField(c, "jscompError", args4j.Class_String, &f.jscompError)},
		{Option: &args4j.Option{Name: "--jscomp_warning", Handler: newCommandLineRunner_Flags_WarningGuardWarningOptionHandler, Usage: "Make the named class of warnings a normal warning. Must be one of the error group items. '*' adds all supported."},
			Field: args4j.ListField(c, "jscompWarning", args4j.Class_String, &f.jscompWarning)},
		{Option: &args4j.Option{Name: "--jscomp_off", Handler: newCommandLineRunner_Flags_WarningGuardOffOptionHandler, Usage: "Turn off the named class of warnings. Must be one of the error group items. '*' adds all supported."},
			Field: args4j.ListField(c, "jscompOff", args4j.Class_String, &f.jscompOff)},
		{Option: &args4j.Option{Name: "--define", Aliases: []string{"--D", "-D"}, Usage: "Override the value of a variable annotated @define. The format is <name>[=<val>], where <name> is the name of a @define variable and <val> is a boolean, number, or a single-quoted string that contains no single quotes. If [=<val>] is omitted, the variable is marked true"},
			Field: args4j.ListField(c, "define", args4j.Class_String, &f.define)},
		{Option: &args4j.Option{Name: "--charset", Usage: "Input and output charset for all files. By default, we accept UTF-8 as input and output US_ASCII"},
			Field: args4j.StringField(c, "charset", &f.charset)},
		{Option: &args4j.Option{Name: "--compilation_level", Aliases: []string{"-O"}, Usage: "Specifies the compilation level to use. Options: BUNDLE, WHITESPACE_ONLY, SIMPLE (default), TRANSPILE_ONLY,ADVANCED"},
			Field: args4j.StringField(c, "compilationLevel", &f.compilationLevel)},
		{Option: &args4j.Option{Name: "--num_parallel_threads", Hidden: true, Handler: args4j.NewIntOptionHandler, Usage: "Use multiple threads to parallelize parts of the compilation."},
			Field: args4j.IntField(c, "numParallelThreads", &f.numParallelThreads)},
		{Option: &args4j.Option{Name: "--fast_mode", Hidden: true, Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Enable Go-to-Go deterministic fast mode."},
			Field: args4j.BooleanField(c, "fastMode", &f.fastMode)},
		{Option: &args4j.Option{Name: "--checks_only", Aliases: []string{"--checks-only"}, Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Don't generate output. Run checks, but no optimization passes."},
			Field: args4j.BooleanField(c, "checksOnly", &f.checksOnly)},
		{Option: &args4j.Option{Name: "--incremental_check_mode", Usage: "Generate or check externs-like .i.js files representing individual libraries."},
			Field: args4j.EnumField(c, "incrementalCheckMode", commandLineRunner_CompilerOptions_IncrementalCheckMode_class, &f.incrementalCheckMode)},
		{Option: &args4j.Option{Name: "--continue_after_errors", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Continue trying to compile after an error is encountered."},
			Field: args4j.BooleanField(c, "continueAfterErrors", &f.continueAfterErrors)},
		{Option: &args4j.Option{Name: "--use_types_for_optimization", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Enable or disable the optimizations based on available type information. Inaccurate type annotations may result in incorrect results."},
			Field: args4j.BooleanField(c, "useTypesForOptimization", &f.useTypesForOptimization)},
		{Option: &args4j.Option{Name: "--assume_function_wrapper", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Enable additional optimizations based on the assumption that the output will be wrapped with a function wrapper.  This flag is used to indicate that \"global\" declarations will not actually be global but instead isolated to the compilation unit. This enables additional optimizations."},
			Field: args4j.BooleanField(c, "assumeFunctionWrapper", &f.assumeFunctionWrapper)},
		{Option: &args4j.Option{Name: "--warning_level", Aliases: []string{"-W"}, Usage: "Specifies the warning level to use. Options: QUIET, DEFAULT, VERBOSE"},
			Field: args4j.EnumField(c, "warningLevel", commandLineRunner_WarningLevel_class, &f.warningLevel)},
		{Option: &args4j.Option{Name: "--debug", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Enable debugging options. Property renaming uses long mangled names which can be mapped back to the original name."},
			Field: args4j.BooleanField(c, "debug", &f.debug)},
		{Option: &args4j.Option{Name: "--typed_ast_output_file", Hidden: true, Usage: "Sets file to output in-progress typedAST format. DO NOT USE!"},
			Field: args4j.NullableStringField(c, "typedAstOutputFile", &f.typedAstOutputFile)},
		{Option: &args4j.Option{Name: "--generate_exports", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Generates export code for those marked with @export"},
			Field: args4j.BooleanField(c, "generateExports", &f.generateExports)},
		{Option: &args4j.Option{Name: "--export_local_property_definitions", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Generates export code for local properties marked with @export"},
			Field: args4j.BooleanField(c, "exportLocalPropertyDefinitions", &f.exportLocalPropertyDefinitions)},
		{Option: &args4j.Option{Name: "--formatting", Usage: "Specifies which formatting options, if any, should be applied to the output JS. Options: PRETTY_PRINT, PRINT_INPUT_DELIMITER, SINGLE_QUOTES"},
			Field: args4j.ListField(c, "formatting", commandLineRunner_FormattingOption_class, &f.formatting)},
		{Option: &args4j.Option{Name: "--process_common_js_modules", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Process CommonJS modules to a concatenable form."},
			Field: args4j.BooleanField(c, "processCommonJsModules", &f.processCommonJsModules)},
		{Option: &args4j.Option{Name: "--js_module_root", Usage: "Path prefixes to be removed from ES6 & CommonJS modules."},
			Field: args4j.ListField(c, "moduleRoot", args4j.Class_String, &f.moduleRoot)},
		{Option: &args4j.Option{Name: "--process_closure_primitives", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Processes built-ins from the Closure library, such as goog.require(), goog.provide(), and goog.exportSymbol(). True by default."},
			Field: args4j.BooleanField(c, "processClosurePrimitives", &f.processClosurePrimitives)},
		{Option: &args4j.Option{Name: "--angular_pass", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Generate $inject properties for AngularJS for functions annotated with @ngInject"},
			Field: args4j.BooleanField(c, "angularPass", &f.angularPass)},
		{Option: &args4j.Option{Name: "--polymer_version", Usage: "Which version of Polymer is being used (1 or 2)."},
			Field: args4j.NullableIntegerField(c, "polymerVersion", &f.polymerVersion)},
		{Option: &args4j.Option{Name: "--chrome_pass", Hidden: true, Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Enable Chrome-specific options for handling cr.* functions."},
			Field: args4j.BooleanField(c, "chromePass", &f.chromePass)},
		{Option: &args4j.Option{Name: "--j2cl_pass", Hidden: true, Usage: "Rewrite J2CL output to be compiler-friendly if enabled (ON or AUTO). Options:OFF, ON, AUTO(default)"},
			Field: args4j.StringField(c, "j2clPassMode", &f.j2clPassMode)},
		{Option: &args4j.Option{Name: "--remove_j2cl_asserts", Hidden: true, Usage: "Remove calls to J2CL assertions."},
			Field: args4j.BooleanField(c, "removeJ2cLAsserts", &f.removeJ2cLAsserts)},
		{Option: &args4j.Option{Name: "--output_manifest", Usage: "Prints out a list of all the files in the compilation. If --dependency_mode=PRUNE or PRUNE_LEGACY is specified, this will not include files that got dropped because they were not required. The %outname% placeholder expands to the JS output file. If you're using modularization, using %outname% will create a manifest for each module."},
			Field: args4j.StringField(c, "outputManifest", &f.outputManifest)},
		{Option: &args4j.Option{Name: "--output_chunk_dependencies", Usage: "Prints out a JSON file of dependencies between chunks."},
			Field: args4j.StringField(c, "outputChunkDependencies", &f.outputChunkDependencies)},
		{Option: &args4j.Option{Name: "--language_in", Usage: "Sets the language spec to which input sources should conform. Options: ECMASCRIPT3, ECMASCRIPT5, ECMASCRIPT5_STRICT, ECMASCRIPT_2015, ECMASCRIPT_2016, ECMASCRIPT_2017, ECMASCRIPT_2018, ECMASCRIPT_2019, ECMASCRIPT_2020,ECMASCRIPT_2021, STABLE, ECMASCRIPT_NEXT (latest features supported),UNSTABLE (for testing only)"},
			Field: args4j.StringField(c, "languageIn", &f.languageIn)},
		{Option: &args4j.Option{Name: "--language_out", Usage: "Sets the language spec to which output should conform. Options: ECMASCRIPT3, ECMASCRIPT5, ECMASCRIPT_2015, ECMASCRIPT_2016, ECMASCRIPT_2017, ECMASCRIPT_2018, ECMASCRIPT_2019, ECMASCRIPT_2020, ECMASCRIPT_2021, STABLE, ECMASCRIPT_NEXT (latest features supported)"},
			Field: args4j.NewField(c, "languageOut", args4j.Class_String, nil, false,
				func() any { return f.languageOut },
				func(value any) { f.languageOut, f.languageOutIsDefault = value.(string), false })},
		{Option: &args4j.Option{Name: "--version", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Prints the compiler version to stdout and exit."},
			Field: args4j.BooleanField(c, "version", &f.version)},
		{Option: &args4j.Option{Name: "--translations_file", Hidden: true, Usage: "Source of translated messages. Currently only supports XTB."},
			Field: args4j.StringField(c, "translationsFile", &f.translationsFile)},
		{Option: &args4j.Option{Name: "--translations_project", Hidden: true, Usage: "Scopes all translations to the specified project.When specified, we will use different message ids so that messages in different projects can have different translations."},
			Field: args4j.NullableStringField(c, "translationsProject", &f.translationsProject)},
		{Option: &args4j.Option{Name: "--flagfile", Hidden: true, Usage: "A file (or files) containing additional command-line options."},
			Field: args4j.ListField(c, "flagFiles", args4j.Class_String, &f.flagFiles)},
		{Option: &args4j.Option{Name: "--warnings_allowlist_file", Aliases: []string{"--warnings_whitelist_file"}, Usage: "A file containing warnings to suppress. Each line should be of the form\n<file-name>:<line-number>?  <warning-description>"},
			Field: args4j.StringField(c, "warningsAllowlistFile", &f.warningsAllowlistFile)},
		{Option: &args4j.Option{Name: "--hide_warnings_for", Usage: "If specified, files whose path contains this string will have their warnings hidden. You may specify multiple."},
			Field: args4j.ListField(c, "hideWarningsFor", args4j.Class_String, &f.hideWarningsFor)},
		{Option: &args4j.Option{Name: "--extra_annotation_name", Usage: "A allowlist of tag names in JSDoc. You may specify multiple"},
			Field: args4j.ListField(c, "extraAnnotationName", args4j.Class_String, &f.extraAnnotationName)},
		{Option: &args4j.Option{Name: "--tracer_mode", Hidden: true, Usage: "Shows the duration of each compiler pass and the impact to the compiled output size. Options: ALL, AST_SIZE, RAW_SIZE, TIMING_ONLY, OFF"},
			Field: args4j.EnumField(c, "tracerMode", commandLineRunner_CompilerOptions_TracerMode_class, &f.tracerMode)},
		{Option: &args4j.Option{Name: "--rename_variable_prefix", Usage: "Specifies a prefix that will be prepended to all variables."},
			Field: args4j.NullableStringField(c, "renamePrefix", &f.renamePrefix)},
		{Option: &args4j.Option{Name: "--rename_prefix_namespace", Usage: "Specifies the name of an object that will be used to store all non-extern globals"},
			Field: args4j.NullableStringField(c, "renamePrefixNamespace", &f.renamePrefixNamespace)},
		{Option: &args4j.Option{Name: "--conformance_configs", Usage: "A list of JS Conformance configurations in text protocol buffer format."},
			Field: args4j.ListField(c, "conformanceConfigs", args4j.Class_String, &f.conformanceConfigs)},
		{Option: &args4j.Option{Name: "--env", Usage: "Determines the set of builtin externs to load. Options: BROWSER, CUSTOM. Defaults to BROWSER."},
			Field: args4j.EnumField(c, "environment", commandLineRunner_CompilerOptions_Environment_class, &f.environment)},
		{Option: &args4j.Option{Name: "--json_streams", Usage: "Specifies whether standard input and output streams will be a JSON array of sources. Each source will be an object of the form {path: filename, src: file_contents, source_map: srcmap_contents }. Intended for use by stream-based build systems such as gulpjs. Options: NONE, IN, OUT, BOTH. Defaults to NONE."},
			Field: args4j.EnumField(c, "jsonStreamMode", commandLineRunner_CompilerOptions_JsonStreamMode_class, &f.jsonStreamMode)},
		{Option: &args4j.Option{Name: "--preserve_type_annotations", Hidden: true, Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Preserves type annotations."},
			Field: args4j.BooleanField(c, "preserveTypeAnnotations", &f.preserveTypeAnnotations)},
		{Option: &args4j.Option{Name: "--inject_libraries", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Allow injecting runtime libraries."},
			Field: args4j.BooleanField(c, "injectLibraries", &f.injectLibraries)},
		{Option: &args4j.Option{Name: "--force_inject_library", Usage: "Force injection of named runtime libraries. The format is <name> where <name> is the name of a runtime library. Possible libraries include: base, es6_runtime"},
			Field: args4j.ListField(c, "forceInjectLibraries", args4j.Class_String, &f.forceInjectLibraries)},
		{Option: &args4j.Option{Name: "--dependency_mode", Usage: "Specifies how the compiler should determine the set and order of files for a compilation. Options: NONE the compiler will include all src files in the order listed, SORT_ONLY the compiler will include all source files in dependency order, PRUNE files will only be included if they are transitive dependencies of files listed in the --entry_point flag and then sorted in dependency order, PRUNE_LEGACY same as PRUNE but files that do not goog.provide a namespace and are not modules will be automatically added as --entry_point entries. Defaults to PRUNE_LEGACY if entry points are defined, otherwise to NONE."},
			Field: args4j.NullableEnumField(c, "dependencyMode", commandLineRunner_DependencyOptions_DependencyMode_class, &f.dependencyMode)},
		{Option: &args4j.Option{Name: "--entry_point", Usage: "A file or namespace to use as the starting point for determining which src files to include in the compilation. ES6 and CommonJS modules are specified as file paths (without the extension). Closure-library namespaces are specified with a \"goog:\" prefix. Example: --entry_point=goog:goog.Promise"},
			Field: args4j.ListField(c, "entryPoint", args4j.Class_String, &f.entryPoint)},
		{Option: &args4j.Option{Name: "--rewrite_polyfills", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Injects polyfills for ES2015+ library classes and methods used in source. See also the \"Polyfills\" GitHub Wiki page."},
			Field: args4j.BooleanField(c, "rewritePolyfills", &f.rewritePolyfills)},
		{Option: &args4j.Option{Name: "--isolate_polyfills", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Hides injected polyfills from the global scope and any external code. See the the \"Polyfills\" GitHub Wiki page for details."},
			Field: args4j.BooleanField(c, "isolatePolyfills", &f.isolatePolyfills)},
		{Option: &args4j.Option{Name: "--print_source_after_each_pass", Hidden: true, Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Whether to iteratively print resulting JS source per pass."},
			Field: args4j.BooleanField(c, "printSourceAfterEachPass", &f.printSourceAfterEachPass)},
		{Option: &args4j.Option{Name: "--module_resolution", Usage: "Specifies how the compiler locates modules. BROWSER requires all module imports to begin with a '.' or '/' and have a file extension. NODE uses the node module rules. WEBPACK looks up modules from a special lookup map."},
			Field: args4j.EnumField(c, "moduleResolutionMode", commandLineRunner_ModuleLoader_ResolutionMode_class, &f.moduleResolutionMode)},
		{Option: &args4j.Option{Name: "--browser_resolver_prefix_replacements", Usage: "Prefixes to replace in ES6 import paths before resolving. module_resolution must be BROWSER_WITH_TRANSFORMED_PREFIXES to take effect."},
			Field: args4j.LinkedStringMapField(c, "browserResolverPrefixReplacements", &f.browserResolverPrefixReplacements)},
		{Option: &args4j.Option{Name: "--package_json_entry_names", Usage: "Ordered list of entries to look for in package.json files when processing modules with the NODE module resolution strategy (i.e. esnext:main,browser,main). Defaults to a list with the following entries: \"browser\", \"module\", \"main\"."},
			Field: args4j.NullableStringField(c, "packageJsonEntryNames", &f.packageJsonEntryNames)},
		{Option: &args4j.Option{Name: "--error_format", Usage: "Specifies format for error messages."},
			Field: args4j.EnumField(c, "errorFormat", commandLineRunner_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_class, &f.errorFormat)},
		{Option: &args4j.Option{Name: "--renaming", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Disables variable renaming. Cannot be used with ADVANCED optimizations."},
			Field: args4j.BooleanField(c, "renaming", &f.renaming)},
		{Option: &args4j.Option{Name: "--help_markdown", Hidden: true, Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Prints markdown formatted flag usage"},
			Field: args4j.BooleanField(c, "helpMarkdown", &f.helpMarkdown)},
		{Option: &args4j.Option{Name: "--instrument_for_coverage_option", Usage: "Enable code instrumentation to perform code coverage analysis. Options are:\n 1. NONE (default)\n 2. LINE - Instrument code by line.\n 3. BRANCH - Instrument code by branch.\n 4. PRODUCTION - Function Instrumentation on compiled JS code.\n"},
			Field: args4j.StringField(c, "instrumentForCoverageOption", &f.instrumentForCoverageOption)},
		{Option: &args4j.Option{Name: "--production_instrumentation_array_name", Usage: "Name of the global array used by production instrumentation. The array name should be declared as an extern so it is not renamed by the compiler. A function that parses the global array should also be included. This flag is to be used in tandem with --instrument_code=PRODUCTION"},
			Field: args4j.StringField(c, "productionInstrumentationArrayName", &f.productionInstrumentationArrayName)},
		{Option: &args4j.Option{Name: "--chunk_output_type", Usage: "Indicates what format the compiler should use for output chunks. GLOBAL_NAMESPACE is typically used in conjunction with --rename_prefix_namespace. ES_MODULES outputs chunks as proper modules with 'import' and 'export' statements."},
			Field: args4j.EnumField(c, "chunkOutputType", commandLineRunner_CompilerOptions_ChunkOutputType_class, &f.chunkOutputType)},
		{Option: &args4j.Option{Name: "--allow_dynamic_import", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Indicates that the compiler should allow dynamic import expressions. Dynamic import expressions are not yet fully supported and may lead to broken output code."},
			Field: args4j.BooleanField(c, "allowDynamicImport", &f.allowDynamicImport)},
		{Option: &args4j.Option{Name: "--dynamic_import_alias", Usage: "Instructs the compiler to replace dynamic imports expressions with a function call using the specified name. Allows dynamic import expressions to be externally polyfilled when the output language level does not natively support them. An alias of 'import' is allowed."},
			Field: args4j.NullableStringField(c, "dynamicImportAlias", &f.dynamicImportAlias)},
		{Option: &args4j.Option{Name: "--assume_static_inheritance_is_not_used", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Assume that static (class-side) inheritance is not being used and that static methods will not be referenced via `this` or through subclasses. This enables optimizations that could break code that did those things."},
			Field: args4j.BooleanField(c, "assumeStaticInheritanceIsNotUsed", &f.assumeStaticInheritanceIsNotUsed)},
		{Option: &args4j.Option{Name: "--assume_no_prototype_method_enumeration", Handler: newCommandLineRunner_Flags_BooleanOptionHandler, Usage: "Assume that prototype method enumeration is not being used. This allows the compiler to move a prototype method declaration into a deeper chunk without creating stub functions in a parent chunk."},
			Field: args4j.BooleanField(c, "assumeNoPrototypeMethodEnumeration", &f.assumeNoPrototypeMethodEnumeration)},
		{Option: &args4j.Option{Name: "--variable_map_input_file", Usage: "File containing the serialized version of the variable renaming map produced by a previous compilation"},
			Field: args4j.StringField(c, "variableMapInputFile", &f.variableMapInputFile)},
		{Option: &args4j.Option{Name: "--property_map_input_file", Usage: "File containing the serialized version of the property renaming map produced by a previous compilation"},
			Field: args4j.StringField(c, "propertyMapInputFile", &f.propertyMapInputFile)},
		{Option: &args4j.Option{Name: "--expected_diagnostics", Usage: "Expected diagnostics in the format [(line,col)]CODE:regex. CODE is the JSCompiler DiagnosticType key. regex matches the description."},
			Field: args4j.ListField(c, "expectedDiagnostics", args4j.Class_String, &f.expectedDiagnostics)},
		{Argument: &args4j.Argument{},
			Field: args4j.ListField(c, "arguments", args4j.Class_String, &f.arguments)},
	}
}

// Parse parses the given args list.
//
// Java: CommandLineRunner.Flags.parse(List<String>)
func (f *commandLineRunner_Flags) Parse(args []string) {
	f.parser.ParseArgumentCollection(args)

	compilationLevel := jutil.Ascii_ToUpperCase(f.compilationLevel)
	f.compilationLevelParsed = CompilationLevel_FromString(&compilationLevel)
	if f.compilationLevelParsed == nil {
		panic(args4j.NewCmdLineExceptionWithParser(
			f.parser, "Bad value for --compilation_level: "+f.compilationLevel))
	}

	instrumentForCoverageOption := jutil.Ascii_ToUpperCase(f.instrumentForCoverageOption)
	instrumentCodeParsed := CompilerOptions_InstrumentOption_FromString(&instrumentForCoverageOption)
	if instrumentCodeParsed == nil {
		panic(args4j.NewCmdLineExceptionWithParser(
			f.parser,
			"Bad value for --instrument_for_coverage_option: "+f.instrumentForCoverageOption))
	}
	f.instrumentCodeParsed = *instrumentCodeParsed
}

// commandLineRunner_Flags_categories is an ImmutableMultimap: its keys, in
// insertion order, each with its values.
//
// Go-only: language-adapter — ImmutableMultimap<String, String> as an
// ordered slice of entries (only asMap() is used).
//
// Java: CommandLineRunner.Flags.categories
var commandLineRunner_Flags_categories = []struct {
	key    string
	values []string
}{
	{"Basic Usage", []string{
		"compilation_level",
		"env",
		"externs",
		"js",
		"js_output_file",
		"language_in",
		"language_out",
		"warning_level"}},
	{"Warning and Error Management", []string{
		"conformance_configs",
		"error_format",
		"extra_annotation_name",
		"hide_warnings_for",
		"jscomp_error",
		"jscomp_off",
		"jscomp_warning",
		"strict_mode_input",
		"warnings_allowlist_file"}},
	{"Output", []string{
		"assume_function_wrapper",
		"debug",
		"emit_use_strict",
		"export_local_property_definitions",
		"formatting",
		"generate_exports",
		"isolation_mode",
		"output_wrapper",
		"output_wrapper_file",
		"rename_variable_prefix"}},
	{"Dependency Management", []string{"dependency_mode", "entry_point"}},
	{"JS Modules", []string{
		"dynamic_import_alias",
		"js_module_root",
		"module_resolution",
		"process_common_js_modules",
		"package_json_entry_names"}},
	{"Library and Framework Specific", []string{
		"angular_pass",
		"force_inject_library",
		"inject_libraries",
		"polymer_version",
		"process_closure_primitives",
		"rewrite_polyfills",
		"isolate_polyfills"}},
	{"Code Splitting", []string{
		"chunk",
		"chunk_output_path_prefix",
		"chunk_output_type",
		"chunk_wrapper",
		"rename_prefix_namespace"}},
	{"Reports", []string{
		"create_source_map",
		"output_manifest",
		"output_chunk_dependencies",
		"property_renaming_report",
		"source_map_input",
		"source_map_include_content",
		"source_map_location_mapping",
		"variable_renaming_report"}},
	{"Miscellaneous", []string{
		"assume_static_inheritance_is_not_used",
		"browser_featureset_year",
		"charset",
		"checks_only",
		"define",
		"flagfile",
		"help",
		"json_streams",
		"third_party",
		"use_types_for_optimization",
		"version"}},
}

// Java: CommandLineRunner.Flags.printUsage(PrintStream)
func (f *commandLineRunner_Flags) PrintUsage(ps *jutil.PrintStream) {
	outputStream := jutil.NewOutputStreamWriter(abstractCommandLineRunner_printStreamAsOutputStream{ps})

	isFirst := true
	for _, entry := range commandLineRunner_Flags_categories {
		prefix := "\n\n"
		suffix := ""
		if isFirst {
			isFirst = false
			prefix = ""
		}

		if entry.key == "Warning and Error Management" {
			if f.helpMarkdown {
				suffix = "\n## Available Error Groups\n\n" +
					"  - " +
					strings.ReplaceAll(DiagnosticGroups_DIAGNOSTIC_GROUP_NAMES, ", ", "\n  - ")
			} else {
				suffix = "\n" +
					commandLineRunner_Flags_BOLD_PREFIX +
					"Available Error Groups: " +
					commandLineRunner_Flags_NORMAL_PREFIX +
					DiagnosticGroups_DIAGNOSTIC_GROUP_NAMES
			}
		}

		if f.helpMarkdown {
			// For markdown docs we don't want any line wrapping so we just set a very
			// large line length.
			f.maxLineLength = 5000
			f.parser.SetUsageWidth(f.maxLineLength)
		}

		f.PrintCategoryUsage(entry.key, entry.values, outputStream, &prefix, &suffix)
	}

	ps.Flush()
}

// Java: CommandLineRunner.Flags.BOLD_PREFIX
const commandLineRunner_Flags_BOLD_PREFIX = "\033[1m"

// Java: CommandLineRunner.Flags.NORMAL_PREFIX
const commandLineRunner_Flags_NORMAL_PREFIX = "\033[0m"

// Java: CommandLineRunner.Flags.MARKDOWN_CHARS_TO_ESCAPE
const commandLineRunner_Flags_MARKDOWN_CHARS_TO_ESCAPE = "[-*\\`\\[\\]{}\\(\\)#+\\.!<>]"

// commandLineRunner_Flags_optionName is
// namedOptionDef.name().replaceFirst("^--", "") for a handler whose option
// is a NamedOptionDef, and ok false otherwise.
//
// Go-only: language-adapter — the instanceof pattern of the lambdas below.
func commandLineRunner_Flags_optionName(optionHandler args4j.OptionHandler) (name string, ok bool) {
	namedOptionDef, ok := optionHandler.GetOptionHandlerBase().Option.(*args4j.NamedOptionDef)
	if !ok {
		return "", false
	}
	return jutil.StringReplaceFirst(namedOptionDef.Name(), "^--", ""), true
}

// Java: CommandLineRunner.Flags.printCategoryUsage(String, Collection<String>, OutputStreamWriter, String, String)
func (f *commandLineRunner_Flags) PrintCategoryUsage(
	categoryName string,
	options []string,
	outputStream jutil.Writer,
	prefix *string,
	suffix *string) {

	jutil.Try(func() {
		if prefix != nil {
			f.PrintStringLineWrapped(*prefix, outputStream)
		}

		if f.helpMarkdown {
			outputStream.Write("# " + categoryName + "\n")

			for _, optionName := range options {
				stringWriter := jutil.NewStringWriter()
				f.parser.PrintUsageToWriterWithFilter(
					stringWriter,
					nil,
					args4j.OptionHandlerFilterFunc(func(optionHandler args4j.OptionHandler) bool {
						if name, ok := commandLineRunner_Flags_optionName(optionHandler); ok {
							return !optionHandler.GetOptionHandlerBase().Option.Hidden() &&
								optionName == name
						}
						return false
					}))
				stringWriter.Flush()
				rawOptionUsage := stringWriter.String()
				delimiterIndex := jutil.StringIndexOf(rawOptionUsage, " : ")
				if delimiterIndex > 0 {
					outputStream.Write(
						"\n**" + jutil.StringTrim(jutil.StringSubstring(rawOptionUsage, 0, delimiterIndex)) + "**  \n")

					optionDescription := jutil.StringTrim(jutil.StringReplaceAll(
						jutil.StringSubstringFrom(rawOptionUsage, delimiterIndex+3),
						commandLineRunner_Flags_MARKDOWN_CHARS_TO_ESCAPE, "\\\\$0"))
					outputStream.Write(optionDescription + "\n")
				} else {
					outputStream.Write(jutil.StringReplaceAll(rawOptionUsage, commandLineRunner_Flags_MARKDOWN_CHARS_TO_ESCAPE, "\\\\$0"))
				}
				outputStream.Flush()
			}
		} else {
			outputStream.Write(commandLineRunner_Flags_BOLD_PREFIX + categoryName + ":\n" + commandLineRunner_Flags_NORMAL_PREFIX)
			f.parser.PrintUsageToWriterWithFilter(
				outputStream,
				nil,
				args4j.OptionHandlerFilterFunc(func(optionHandler args4j.OptionHandler) bool {
					if name, ok := commandLineRunner_Flags_optionName(optionHandler); ok {
						return !optionHandler.GetOptionHandlerBase().Option.Hidden() &&
							slices.Contains(options, name)
					}
					return false
				}))
		}

		if suffix != nil {
			f.PrintStringLineWrapped(*suffix, outputStream)
		}
	}, func(e jutil.IIOException) {
		// Ignore.
	})
}

// Java: CommandLineRunner.Flags.WHITESPACE_PATTERN
var commandLineRunner_Flags_WHITESPACE_PATTERN = jutil.Pattern_Compile("\\s")

// Java: CommandLineRunner.Flags.printStringLineWrapped(String, OutputStreamWriter)
func (f *commandLineRunner_Flags) PrintStringLineWrapped(input string, outputStream jutil.Writer) {
	if jutil.StringLength(input) < f.maxLineLength {
		outputStream.Write(input)
		return
	}

	endIndex := f.maxLineLength
	subString := jutil.StringSubstring(input, 0, f.maxLineLength)
	whitespaceMatcher := commandLineRunner_Flags_WHITESPACE_PATTERN.Matcher(subString)
	foundMatch := false
	for whitespaceMatcher.Find() {
		endIndex = whitespaceMatcher.Start()
		foundMatch = true
	}
	outputStream.Write(jutil.StringSubstring(input, 0, endIndex) + "\n")
	rest := endIndex
	if foundMatch {
		rest = endIndex + 1
	}
	f.PrintStringLineWrapped("    "+jutil.StringSubstringFrom(input, rest), outputStream)
}

// Java: CommandLineRunner.Flags.printShortUsageAfterErrors(PrintStream)
func commandLineRunner_Flags_printShortUsageAfterErrors(ps *jutil.PrintStream) {
	ps.Print("Sample usage: ")
	ps.Println(
		"--compilation_level (-O) VAL --externs VAL --js VAL" +
			" --js_output_file VAL" +
			" --warning_level (-W) [QUIET | DEFAULT | VERBOSE]")
	ps.Println("Run with --help for all options and details")
	ps.Flush()
}

// GetMixedJsSources: users may specify JS inputs via the --js flag, as well
// as via additional arguments to the Closure Compiler. For example, it is
// convenient to leverage the additional arguments feature when using the
// Closure Compiler in combination with find and xargs:
//
//	find MY_JS_SRC_DIR -name '*.js' \
//	    | xargs java -jar compiler.jar --dependency_mode=PRUNE_LEGACY
//
// The find command will produce a list of '*.js' source files in the
// MY_JS_SRC_DIR directory while xargs will convert them to a single,
// space-delimited set of arguments that are appended to the java command to
// run the Compiler.
//
// Note that it is important to use the --dependency_mode=PRUNE or
// PRUNE_LEGACY option in this case because the order produced by find is
// unlikely to be sorted correctly with respect to goog.provide() and
// goog.requires().
//
// Java: CommandLineRunner.Flags.getMixedJsSources()
func (f *commandLineRunner_Flags) GetMixedJsSources() []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType] {
	mixedSources := make([]*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType], 0, len(commandLineRunner_Flags_mixedJsSources)+len(f.arguments))
	var excludes *jutil.LinkedHashSet[string]
	allocJsEntry := func(filename string) *AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType] {
		return NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS, filename)
	}
	for _, source := range commandLineRunner_Flags_mixedJsSources {
		val := source.GetValue()
		if strings.HasSuffix(val, ".zip") {
			mixedSources = append(mixedSources, source)
		} else if strings.HasPrefix(val, "!") {
			if excludes == nil {
				excludes = jutil.NewLinkedHashSet[string]()
			}
			for _, filename := range CommandLineRunner_FindJsFiles([]string{val[1:]}) {
				excludes.Add(filename)
				// Java: List.remove(Object) removes the first equal element.
				if i := slices.IndexFunc(mixedSources, func(e *AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]) bool {
					return e != nil && e.flag == AbstractCommandLineRunner_JsSourceType_JS && e.value == filename
				}); i >= 0 {
					mixedSources = slices.Delete(mixedSources, i, i+1)
				}
			}
		} else if !strings.Contains(val, "*") && !jutil.File_IsDirectoryPath(val) {
			filename := jutil.Path_NormalizeString(val)
			if excludes == nil || !excludes.Contains(filename) {
				if source.flag == AbstractCommandLineRunner_JsSourceType_JS && filename == val {
					mixedSources = append(mixedSources, source)
				} else {
					mixedSources = append(mixedSources, allocJsEntry(filename))
				}
			}
		} else {
			for _, filename := range commandLineRunner_findJsFilesSorted([]string{val}, true) {
				if excludes == nil || !excludes.Contains(filename) {
					mixedSources = append(mixedSources, allocJsEntry(filename))
				}
			}
		}
	}
	if len(f.arguments) != 0 {
		fromArguments := CommandLineRunner_FindJsFiles(f.arguments)
		for _, filename := range fromArguments {
			mixedSources = append(mixedSources, allocJsEntry(filename))
		}
	}
	if len(commandLineRunner_Flags_mixedJsSources) != 0 && len(f.arguments) != 0 && len(mixedSources) == 0 {
		panic(args4j.NewCmdLineExceptionWithParser(f.parser, "No inputs matched"))
	}
	return mixedSources
}

// Java: CommandLineRunner.Flags.getSourceMapLocationMappings()
func (f *commandLineRunner_Flags) GetSourceMapLocationMappings() []SourceMap_LocationMapping {
	var locationMappings []SourceMap_LocationMapping

	split := f.SplitPipeParts(f.sourceMapLocationMapping, "--source_map_location_mapping")
	for key, value := range split.All() {
		locationMappings = append(locationMappings,
			NewSourceMap_PrefixLocationMapping(key, value))
	}

	return locationMappings
}

// Java: CommandLineRunner.Flags.getSourceMapInputs()
func (f *commandLineRunner_Flags) GetSourceMapInputs() *jutil.LinkedHashMap[string, string] {
	return f.SplitPipeParts(f.sourceMapInputs, "--source_map_input")
}

// Java: CommandLineRunner.Flags.splitPipeParts(Iterable<String>, String)
func (f *commandLineRunner_Flags) SplitPipeParts(input []string, flagName string) *jutil.LinkedHashMap[string, string] {
	result := jutil.NewImmutableMapBuilder[string, string]()

	// Go-only: language-adapter — Guava Splitter.on('|').limit(2) is
	// strings.SplitN on the ASCII separator.
	for _, inputSourceMap := range input {
		parts := strings.SplitN(inputSourceMap, "|", 2)
		if len(parts) != 2 {
			panic(args4j.NewCmdLineExceptionWithParser(
				f.parser, "Bad value for "+flagName+" (duplicate key): "+jutil.CollectionToString(input)))
		}
		result.Put(parts[0], parts[1])
	}

	return result.BuildOrThrow()
}

// Java: CommandLineRunner.Flags.getPackageJsonEntryNames()
func (f *commandLineRunner_Flags) GetPackageJsonEntryNames() []string {
	// Go-only: language-adapter — Guava Splitter.on(',') is strings.Split on
	// the ASCII separator.
	return strings.Split(*f.packageJsonEntryNames, ",")
}

// commandLineRunner_Flags_BooleanOptionHandler is our own option parser to
// be backwards-compatible. It needs to be public because of the crazy
// reflection that args4j does.
//
// Java: CommandLineRunner.Flags.BooleanOptionHandler
type commandLineRunner_Flags_BooleanOptionHandler struct {
	args4j.OptionHandlerBase
}

// Java: CommandLineRunner.Flags.BooleanOptionHandler.TRUES
var commandLineRunner_Flags_BooleanOptionHandler_TRUES = jutil.NewLinkedHashSet("true", "on", "yes", "1")

// Java: CommandLineRunner.Flags.BooleanOptionHandler.FALSES
var commandLineRunner_Flags_BooleanOptionHandler_FALSES = jutil.NewLinkedHashSet("false", "off", "no", "0")

// Handlers are used reflectively by args4j
//
// Java: CommandLineRunner.Flags.BooleanOptionHandler(CmdLineParser, OptionDef, Setter<? super Boolean>)
func newCommandLineRunner_Flags_BooleanOptionHandler(
	parser *args4j.CmdLineParser, option args4j.OptionDef, setter args4j.Setter) args4j.OptionHandler {
	h := &commandLineRunner_Flags_BooleanOptionHandler{}
	h.InitOptionHandlerBase(h, parser, option, setter)
	return h
}

// Java: CommandLineRunner.Flags.BooleanOptionHandler.parseArguments(Parameters)
func (h *commandLineRunner_Flags_BooleanOptionHandler) ParseArguments(params args4j.Parameters) int {
	var param *string
	jutil.Try(func() {
		p := params.GetParameter(0)
		param = &p
	}, func(e args4j.ICmdLineException) {
		param = nil // to stop linter complaints
	})

	if param == nil {
		h.Setter.AddValue(true)
		return 0
	} else {
		lowerParam := jutil.Ascii_ToLowerCase(*param)
		if commandLineRunner_Flags_BooleanOptionHandler_TRUES.Contains(lowerParam) {
			h.Setter.AddValue(true)
		} else if commandLineRunner_Flags_BooleanOptionHandler_FALSES.Contains(lowerParam) {
			h.Setter.AddValue(false)
		} else {
			h.Setter.AddValue(true)
			return 0
		}
		return 1
	}
}

// Java: CommandLineRunner.Flags.BooleanOptionHandler.getDefaultMetaVariable()
func (h *commandLineRunner_Flags_BooleanOptionHandler) GetDefaultMetaVariable() *string {
	return nil
}

// commandLineRunner_Flags_WarningGuardErrorOptionHandler is our own parser
// for warning guards that preserves the original order of the flags.
//
// Java: CommandLineRunner.Flags.WarningGuardErrorOptionHandler
type commandLineRunner_Flags_WarningGuardErrorOptionHandler struct {
	args4j.StringOptionHandler
}

// Handlers are used reflectively by args4j
//
// Java: CommandLineRunner.Flags.WarningGuardErrorOptionHandler(CmdLineParser, OptionDef, Setter<? super String>)
func newCommandLineRunner_Flags_WarningGuardErrorOptionHandler(
	parser *args4j.CmdLineParser, option args4j.OptionDef, setter args4j.Setter) args4j.OptionHandler {
	h := &commandLineRunner_Flags_WarningGuardErrorOptionHandler{}
	h.InitStringOptionHandler(h, parser, option, newCommandLineRunner_Flags_MultiFlagSetter(setter, CheckLevel_ERROR, &commandLineRunner_Flags_guardLevels))
	return h
}

// Java: CommandLineRunner.Flags.WarningGuardWarningOptionHandler
type commandLineRunner_Flags_WarningGuardWarningOptionHandler struct {
	args4j.StringOptionHandler
}

// Handlers are used reflectively by args4j
//
// Java: CommandLineRunner.Flags.WarningGuardWarningOptionHandler(CmdLineParser, OptionDef, Setter<? super String>)
func newCommandLineRunner_Flags_WarningGuardWarningOptionHandler(
	parser *args4j.CmdLineParser, option args4j.OptionDef, setter args4j.Setter) args4j.OptionHandler {
	h := &commandLineRunner_Flags_WarningGuardWarningOptionHandler{}
	h.InitStringOptionHandler(h, parser, option, newCommandLineRunner_Flags_MultiFlagSetter(setter, CheckLevel_WARNING, &commandLineRunner_Flags_guardLevels))
	return h
}

// Java: CommandLineRunner.Flags.WarningGuardOffOptionHandler
type commandLineRunner_Flags_WarningGuardOffOptionHandler struct {
	args4j.StringOptionHandler
}

// Handlers are used reflectively by args4j
//
// Java: CommandLineRunner.Flags.WarningGuardOffOptionHandler(CmdLineParser, OptionDef, Setter<? super String>)
func newCommandLineRunner_Flags_WarningGuardOffOptionHandler(
	parser *args4j.CmdLineParser, option args4j.OptionDef, setter args4j.Setter) args4j.OptionHandler {
	h := &commandLineRunner_Flags_WarningGuardOffOptionHandler{}
	h.InitStringOptionHandler(h, parser, option, newCommandLineRunner_Flags_MultiFlagSetter(setter, CheckLevel_OFF, &commandLineRunner_Flags_guardLevels))
	return h
}

// Java: CommandLineRunner.Flags.JsOptionHandler
type commandLineRunner_Flags_JsOptionHandler struct {
	args4j.StringOptionHandler
}

// Handlers are used reflectively by args4j
//
// Java: CommandLineRunner.Flags.JsOptionHandler(CmdLineParser, OptionDef, Setter<? super String>)
func newCommandLineRunner_Flags_JsOptionHandler(
	parser *args4j.CmdLineParser, option args4j.OptionDef, setter args4j.Setter) args4j.OptionHandler {
	h := &commandLineRunner_Flags_JsOptionHandler{}
	h.InitStringOptionHandler(h, parser, option, newCommandLineRunner_Flags_MultiFlagSetter(setter, AbstractCommandLineRunner_JsSourceType_JS, &commandLineRunner_Flags_mixedJsSources))
	return h
}

// Java: CommandLineRunner.Flags.JsZipOptionHandler
type commandLineRunner_Flags_JsZipOptionHandler struct {
	args4j.StringOptionHandler
}

// Handlers are used reflectively by args4j
//
// Java: CommandLineRunner.Flags.JsZipOptionHandler(CmdLineParser, OptionDef, Setter<? super String>)
func newCommandLineRunner_Flags_JsZipOptionHandler(
	parser *args4j.CmdLineParser, option args4j.OptionDef, setter args4j.Setter) args4j.OptionHandler {
	h := &commandLineRunner_Flags_JsZipOptionHandler{}
	h.InitStringOptionHandler(h, parser, option, newCommandLineRunner_Flags_MultiFlagSetter(setter, AbstractCommandLineRunner_JsSourceType_JS_ZIP, &commandLineRunner_Flags_mixedJsSources))
	return h
}

// Java: CommandLineRunner.Flags.MultiFlagSetter<T>
type commandLineRunner_Flags_MultiFlagSetter[T abstractCommandLineRunner_flag] struct {
	proxy   args4j.Setter
	flag    T
	entries *[]*AbstractCommandLineRunner_FlagEntry[T]
}

// Java: CommandLineRunner.Flags.MultiFlagSetter(Setter<? super String>, T, List<FlagEntry<T>>)
func newCommandLineRunner_Flags_MultiFlagSetter[T abstractCommandLineRunner_flag](
	proxy args4j.Setter, flag T, entries *[]*AbstractCommandLineRunner_FlagEntry[T]) *commandLineRunner_Flags_MultiFlagSetter[T] {
	return &commandLineRunner_Flags_MultiFlagSetter[T]{proxy: proxy, flag: flag, entries: entries}
}

// Java: CommandLineRunner.Flags.MultiFlagSetter.isMultiValued()
func (s *commandLineRunner_Flags_MultiFlagSetter[T]) IsMultiValued() bool {
	return s.proxy.IsMultiValued()
}

// Java: CommandLineRunner.Flags.MultiFlagSetter.getType()
func (s *commandLineRunner_Flags_MultiFlagSetter[T]) GetType() *args4j.Class {
	return s.proxy.GetType()
}

// AddStringValue implements args4j.StringSetter without interface boxing.
func (s *commandLineRunner_Flags_MultiFlagSetter[T]) AddStringValue(value string) {
	normalizedValue := value
	if len(value) >= 1 && (value[0] == '\'' || value[0] == '"') && value[len(value)-1] == value[0] {
		if len(value) == 1 {
			normalizedValue = ""
		} else {
			normalizedValue = value[1 : len(value)-1]
		}
	}
	if ss, ok := s.proxy.(args4j.StringSetter); ok {
		ss.AddStringValue(normalizedValue)
	} else {
		s.proxy.AddValue(normalizedValue)
	}
	*s.entries = append(*s.entries, NewAbstractCommandLineRunner_FlagEntry(s.flag, normalizedValue))
}

// Java: CommandLineRunner.Flags.MultiFlagSetter.addValue(String)
func (s *commandLineRunner_Flags_MultiFlagSetter[T]) AddValue(v any) {
	if v == nil {
		s.proxy.AddValue("")
		*s.entries = append(*s.entries, NewAbstractCommandLineRunner_FlagEntry(s.flag, ""))
		return
	}
	value, _ := v.(string)
	s.AddStringValue(value)
}

// Java: CommandLineRunner.Flags.MultiFlagSetter.asFieldSetter()
func (s *commandLineRunner_Flags_MultiFlagSetter[T]) AsFieldSetter() *args4j.FieldSetter {
	return s.proxy.AsFieldSetter()
}

// Java: CommandLineRunner.Flags.MultiFlagSetter.asAnnotatedElement()
func (s *commandLineRunner_Flags_MultiFlagSetter[T]) AsAnnotatedElement() args4j.AnnotatedElement {
	return s.proxy.AsAnnotatedElement()
}

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
func (o CommandLineRunner_FormattingOption) ApplyToOptions(options *CompilerOptions) {
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

// Java: CommandLineRunner.processArgs(String[])
func commandLineRunner_processArgs(args []string) []string {
	// Args4j has a different format that the old command-line parser.
	// So we use some voodoo to get the args into the format that args4j
	// expects.
	quotesPattern := jutil.Pattern_Compile("^['\"](.*)['\"]$")
	var processedArgs []string

	for _, arg := range args {
		if name, value, ok := commandLineRunner_matchArgPattern(arg); ok {
			processedArgs = append(processedArgs, name)

			quotesMatcher := quotesPattern.Matcher(value)
			if quotesMatcher.Matches() {
				processedArgs = append(processedArgs, *quotesMatcher.Group(1))
			} else {
				processedArgs = append(processedArgs, value)
			}
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
	flagFile := jutil.Path_Of(flagFileString)

	buffer := jutil.Files_NewBufferedReader(flagFile, jutil.StandardCharsets_UTF_8)
	// Builds the tokens.
	builder := jutil.NewStringBuilder()
	// Stores the built tokens.
	var tokens []string
	// Indicates if we are in a "quoted" token.
	quoted := false
	// Indicates if the char being processed has been escaped.
	escaped := false
	// Indicates whether this is the beginning of the file.
	isFirstCharacter := true

	cbuf := make([]uint16, 1)
	for {
		// Java: c = buffer.read()
		if buffer.Read(cbuf, 0, 1) < 0 {
			break
		}
		c := int(cbuf[0])

		// Ignoring the BOM.
		if isFirstCharacter {
			isFirstCharacter = false
			if c == CommandLineRunner_UTF8_BOM_CODE {
				continue
			}
		}

		if c == 32 || c == 9 || c == 10 || c == 13 {
			if quoted {
				builder.AppendChar(uint16(c))
			} else if builder.Length() != 0 {
				tokens = append(tokens, builder.String())
				builder.SetLength(0)
			}
		} else if c == 34 {
			if escaped {
				if quoted {
					builder.SetCharAt(builder.Length()-1, uint16(c))
				} else {
					builder.AppendChar(uint16(c))
				}
			} else {
				quoted = !quoted
			}
		} else {
			builder.AppendChar(uint16(c))
		}

		escaped = c == 92
	}

	buffer.Close()

	if builder.Length() != 0 {
		tokens = append(tokens, builder.String())
	}

	r.flags.flagFiles = []string{}

	tokens = commandLineRunner_processArgs(tokens)

	// Command-line warning levels should override flag file settings,
	// which means they should go last.
	previousGuardLevels := slices.Clone(commandLineRunner_Flags_guardLevels)
	previousMixedJsSources := slices.Clone(commandLineRunner_Flags_mixedJsSources)
	commandLineRunner_Flags_guardLevels = commandLineRunner_Flags_guardLevels[:0]
	commandLineRunner_Flags_mixedJsSources = commandLineRunner_Flags_mixedJsSources[:0]
	r.flags.Parse(tokens)
	commandLineRunner_Flags_guardLevels = append(commandLineRunner_Flags_guardLevels, previousGuardLevels...)
	commandLineRunner_Flags_mixedJsSources = append(commandLineRunner_Flags_mixedJsSources, previousMixedJsSources...)

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
		"Version: " + CommandLineRunnerVersion_COMPILER_VERSION}, "\n")
}

// Java: CommandLineRunner.initConfigFromFlags(String[], PrintStream, PrintStream)
func (r *CommandLineRunner) InitConfigFromFlags(args []string, out *jutil.PrintStream, err *jutil.PrintStream) {
	flags := r.flags

	r.errorStream = err
	processedArgs := commandLineRunner_processArgs(args)

	commandLineRunner_Flags_guardLevels = commandLineRunner_Flags_guardLevels[:0]
	commandLineRunner_Flags_mixedJsSources = commandLineRunner_Flags_mixedJsSources[:0]

	var mixedSources []*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]
	var mappings []SourceMap_LocationMapping
	var sourceMapInputs *jutil.LinkedHashMap[string, string]
	parseInlineSourceMaps := false
	applyInputSourceMaps := false
	jutil.Try(func() {
		flags.Parse(processedArgs)

		r.ProcessFlagFiles()

		mixedSources = flags.GetMixedJsSources()
		mappings = flags.GetSourceMapLocationMappings()
		sourceMapInputs = flags.GetSourceMapInputs()
		parseInlineSourceMaps = *flags.parseInlineSourceMaps
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
		panic(NewFlagUsageException(
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

	if flags.outputWrapper != "" && flags.isolationMode != CompilerOptions_IsolationMode_NONE {
		r.ReportError("--output_wrapper and --isolation_mode may not be used together.")
	}

	if flags.isolationMode == CompilerOptions_IsolationMode_IIFE {
		flags.outputWrapper = "(function(){%output%}).call(this);"
	}

	// Handle --compilation_level=BUNDLE
	bundleFiles := []string{}
	skipNormalOutputs := false
	if flags.compilationLevelParsed != nil && *flags.compilationLevelParsed == CompilationLevel_BUNDLE {
		if flags.jsOutputFile == "" {
			r.ReportError("--compilation_level=BUNDLE cannot be used without a --js_output_file.")
		} else {
			bundleFiles = []string{flags.jsOutputFile}
			flags.jsOutputFile = ""
			skipNormalOutputs = true
		}
	}

	var conv CodingConvention
	if flags.thirdParty {
		conv = CodingConventions_GetDefault()
	} else if flags.chromePass {
		conv = NewChromeCodingConvention()
	} else {
		conv = NewClosureCodingConvention()
	}

	if !flags.renaming &&
		flags.compilationLevelParsed != nil && *flags.compilationLevelParsed == CompilationLevel_ADVANCED_OPTIMIZATIONS {
		r.ReportError("ERROR - renaming cannot be disabled when ADVANCED_OPTIMIZATIONS is used.")
	}

	var dependencyOptions *DependencyOptions
	jutil.Try(func() {
		dependencyOptions = DependencyOptions_FromFlags(
			flags.dependencyMode, flags.entryPoint, []string{}, nil, false, false)
	}, func(e IFlagUsageException) {
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
			SetWarningGuards(commandLineRunner_Flags_guardLevels).
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

		segmentOfCompilationToRun := *flags.segmentOfCompilationToRun
		switch segmentOfCompilationToRun {
		case CompilerOptions_SegmentOfCompilationToRun_CHECKS:
			jutil.CheckStateMsg(
				flags.filenameToRestoreFrom == nil,
				"Cannot restore and run CHECKS segment of compilation")
			config.SetSaveCompilationStateToFilename(flags.filenameToSaveTo, 1)
		case CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS:
			config.SetContinueSavedCompilationFileName(flags.filenameToRestoreFrom, 1)
			config.SetSaveCompilationStateToFilename(flags.filenameToSaveTo, 2)
		case CompilerOptions_SegmentOfCompilationToRun_FINALIZATIONS:
			jutil.CheckStateMsg(
				flags.filenameToSaveTo == nil,
				"Cannot run FINALIZATIONS segment of compilation and then save the result")
			config.SetContinueSavedCompilationFileName(flags.filenameToRestoreFrom, 2)
		case CompilerOptions_SegmentOfCompilationToRun_ENTIRE_COMPILATION:
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Cannot run %s segment of compilation: " + segmentOfCompilationToRun.String()))
		}
	}

	r.errorStream = nil
}

// Java: CommandLineRunner.addAllowlistWarningsGuard(CompilerOptions, File)
func (r *CommandLineRunner) AddAllowlistWarningsGuard(options *CompilerOptions, allowlistFile *jutil.File) {
	options.AddWarningsGuard(AllowlistWarningsGuard_FromFile(allowlistFile))
}

// Java: CommandLineRunner.checkChunkName(String)
func (r *CommandLineRunner) CheckChunkName(name string) {
	if !rhino.TokenStream_IsJSIdentifier(commandLineRunner_extraChunkNameChars.Matcher(name).ReplaceAll("_")) {
		panic(NewFlagUsageException("Invalid chunk name: '" + name + "'"))
	}
}

// Java: CommandLineRunner.createOptions()
func (r *CommandLineRunner) CreateOptions() *CompilerOptions {
	flags := r.flags
	options := NewCompilerOptions()

	if flags.languageIn != "" {
		languageMode := CompilerOptions_LanguageMode_FromString(&flags.languageIn)
		if languageMode != nil && *languageMode == CompilerOptions_LanguageMode_UNSUPPORTED {
			panic(NewFlagUsageException(
				"Cannot specify the unsupported set of features for language_in."))
		}
		if languageMode != nil {
			options.SetLanguageIn(*languageMode)
		} else {
			panic(NewFlagUsageException(
				"Unknown language `" +
					flags.languageIn +
					"' specified. Expected one of: " +
					jutil.CollectionToString(CompilerOptions_LanguageMode_ValidCommandLineNames())))
		}
	}

	languageMode := CompilerOptions_LanguageMode_FromString(&flags.languageOut)
	if languageMode != nil && *languageMode == CompilerOptions_LanguageMode_UNSUPPORTED {
		panic(NewFlagUsageException(
			"Cannot specify the unsupported set of features for language_out."))
	}
	if languageMode != nil {
		options.SetLanguageOut(*languageMode)
	} else {
		panic(NewFlagUsageException(
			"Unknown language `" +
				flags.languageOut +
				"' specified. Expected one of: " +
				jutil.CollectionToString(CompilerOptions_LanguageMode_ValidCommandLineNames())))
	}

	options.SetCodingConvention(NewClosureCodingConvention())

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
		options.SetOutputJs(CompilerOptions_OutputJs_NONE)
	}

	options.SetIncrementalChecks(flags.incrementalCheckMode)

	options.SetContinueAfterErrors(flags.continueAfterErrors)

	if flags.useTypesForOptimization {
		level.SetTypeBasedOptimizationOptions(options)
	}

	if flags.assumeFunctionWrapper ||
		flags.isolationMode == CompilerOptions_IsolationMode_IIFE ||
		flags.chunkOutputType == CompilerOptions_ChunkOutputType_ES_MODULES {
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
			j2clPassMode := commandLineRunner_j2clPassModeValueOf(jutil.Ascii_ToUpperCase(flags.j2clPassMode))
			options.SetJ2clPass(j2clPassMode)
		}, func(ex jutil.IIllegalArgumentException) {
			panic(NewFlagUsageException(
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
				NewXtbMessageBundle(
					jutil.NewFileInputStream(flags.translationsFile), flags.translationsProject))
		}, func(e jutil.IIOException) {
			panic(jutil.NewRuntimeExceptionWithMessageAndCause("Reading XTB file", e))
		})
	} else if CompilationLevel_ADVANCED_OPTIMIZATIONS == level {
		// In SIMPLE or WHITESPACE mode, if the user hasn't specified a
		// translations file, they might reasonably try to write their own
		// implementation of goog.getMsg that makes the substitution at
		// run-time.
		//
		// In ADVANCED mode, goog.getMsg is going to be renamed anyway,
		// so we might as well inline it. But shut off the i18n warnings,
		// because the user didn't really ask for i18n.
		options.SetMessageBundle(NewEmptyMessageBundle())
		options.SetWarningLevel(DiagnosticGroups_MSG_CONVENTIONS, CheckLevel_OFF)
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
		options.SetVariableRenaming(VariableRenamingPolicy_OFF)
		options.SetPropertyRenaming(PropertyRenamingPolicy_OFF)
	}

	if flags.instrumentCodeParsed == CompilerOptions_InstrumentOption_PRODUCTION &&
		flags.instrumentationMappingOutputFile == "" {
		panic(NewFlagUsageException(
			"Expected --instrument_mapping_report to be set when " +
				"--instrument_for_coverage_option is set to Production"))
	}

	if flags.instrumentationMappingOutputFile != "" &&
		flags.instrumentCodeParsed != CompilerOptions_InstrumentOption_PRODUCTION {
		panic(NewFlagUsageException(
			"Expected --instrument_for_coverage_option to be passed with PRODUCTION " +
				"when --instrument_mapping_report is set"))
	}

	if flags.productionInstrumentationArrayName == "" &&
		flags.instrumentCodeParsed == CompilerOptions_InstrumentOption_PRODUCTION {
		panic(NewFlagUsageException(
			"Expected --production_instrumentation_array_name to be set when " +
				"--instrument_for_coverage_option is set to Production"))
	}
	options.SetInstrumentForCoverageOption(flags.instrumentCodeParsed)
	options.SetProductionInstrumentationArrayName(flags.productionInstrumentationArrayName)
	options.SetAllowDynamicImport(flags.allowDynamicImport)
	options.SetDynamicImportAlias(flags.dynamicImportAlias)
	options.SetAssumeStaticInheritanceIsNotUsed(flags.assumeStaticInheritanceIsNotUsed)
	options.SetCrossChunkCodeMotionNoStubMethods(flags.assumeNoPrototypeMethodEnumeration)

	if flags.chunkOutputType == CompilerOptions_ChunkOutputType_ES_MODULES {
		if flags.renamePrefixNamespace != nil {
			panic(NewFlagUsageException(
				"Expected --rename_prefix_namespace not to be specified when " +
					"--chunk_output_type is set to ES_MODULES."))
		}

		if flags.emitUseStrict {
			panic(NewFlagUsageException(
				"Expected --emit_use_strict should not be specified when " +
					"--chunk_output_type is set to ES_MODULES."))
		}

		options.SetChunkOutputType(flags.chunkOutputType)
		options.SetEmitUseStrict(false)

		if level == CompilationLevel_ADVANCED_OPTIMIZATIONS {
			options.SetExtractPrototypeMemberDeclarationsOfMode(
				CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_CHUNK_TEMP)
		}
	}

	return options
}

// commandLineRunner_j2clPassModeValueOf is J2clPassMode.valueOf(String):
// the constant of that name, or IllegalArgumentException.
//
// Go-only: language-adapter — Enum.valueOf.
func commandLineRunner_j2clPassModeValueOf(name string) CompilerOptions_J2clPassMode {
	for _, m := range CompilerOptions_J2clPassMode_values() {
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
func (r *CommandLineRunner) CreateCompiler() *Compiler {
	return NewCompilerWithOutStream(r.GetErrorPrintStream())
}

// Java: CommandLineRunner.getBundler()
func (r *CommandLineRunner) GetBundler() *ClosureBundler {
	if r.bundler != nil {
		return r.bundler
	}

	var moduleRoots []string
	if len(r.flags.moduleRoot) != 0 {
		moduleRoots = slices.Clone(r.flags.moduleRoot)
	} else {
		moduleRoots = []string{ModuleLoader_DEFAULT_FILENAME_PREFIX}
	}

	options := r.CreateOptions()
	r.bundler = NewClosureBundlerWithEs6ModuleTranspiler(
		Transpile_Transpiler_NULL,
		NewBaseTranspiler(
			NewBaseTranspiler_CompilerSupplierWithModuleResolution(
				CompilerOptions_LanguageMode_ECMASCRIPT_NEXT.ToFeatureSet().WithoutFeatures(parser.FeatureSet_Feature_MODULES),
				options.GetModuleResolutionMode(),
				moduleRoots,
				options.GetBrowserResolverPrefixReplacements()),
			/* runtimeLibraryName= */ ""))
	return r.bundler
}

// Java: CommandLineRunner.prepForBundleAndAppendTo(Appendable, CompilerInput, String)
func (r *CommandLineRunner) PrepForBundleAndAppendTo(out jutil.Appendable, input *CompilerInput, content string) {
	r.GetBundler().WithPath(input.GetName()).AppendTo(out, input, content)
}

// Java: CommandLineRunner.appendRuntimeTo(Appendable)
func (r *CommandLineRunner) AppendRuntimeTo(out jutil.Appendable) {
	r.GetBundler().AppendRuntimeTo(out)
}

// Java: CommandLineRunner.createExterns(CompilerOptions)
func (r *CommandLineRunner) CreateExterns(options *CompilerOptions) []*SourceFile {
	externs := r.AbstractCommandLineRunnerBase.CreateExterns(options)
	if r.IsInTestMode() {
		return externs
	}
	builtinExterns := AbstractCommandLineRunner_GetBuiltinExterns(options.GetEnvironment())
	builtinExterns = append(builtinExterns, externs...)
	return builtinExterns
}

// Java: CommandLineRunner.loadConformanceConfigs(List<String>)
func commandLineRunner_loadConformanceConfigs(configPaths []string) []*ConformanceConfig {
	configs := []*ConformanceConfig{}

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
func commandLineRunner_loadConformanceConfig(configFile string) *ConformanceConfig {
	textProto := jutil.Files_AsCharSource(jutil.NewFile(configFile), jutil.StandardCharsets_UTF_8).Read()

	builder := ConformanceConfig_NewBuilder()

	// Looking for BOM.
	if textProto != "" && jutil.StringCharAt(textProto, 0) == CommandLineRunner_UTF8_BOM_CODE {
		// Stripping the BOM.
		textProto = jutil.StringSubstringFrom(textProto, 1)
	}

	jutil.Try(func() {
		TextFormat_Merge(textProto, builder)
	}, func(e jutil.IException) {
		panic(jutil.NewRuntimeExceptionWithCause(e))
	})
	return builder.Build()
}

// Deprecated: use AbstractCommandLineRunner_GetBuiltinExterns(
// CompilerOptions_Environment_BROWSER).
//
// Java: CommandLineRunner.getDefaultExterns()
func CommandLineRunner_GetDefaultExterns() []*SourceFile {
	return AbstractCommandLineRunner_GetBuiltinExterns(CompilerOptions_Environment_BROWSER)
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

// commandLineRunner_jsInputs is the Map<String, String> of findJsFiles: a
// TreeMap or a LinkedHashMap.
//
// Go-only: language-adapter — the two Map implementations behind one
// interface.
type commandLineRunner_jsInputs interface {
	Put(k string, v string) string
	Remove(k string) string
	Values() func(yield func(string) bool)
}

type commandLineRunner_linkedJsInputs struct {
	*jutil.LinkedHashMap[string, string]
}

func (m commandLineRunner_linkedJsInputs) Values() func(yield func(string) bool) {
	return m.LinkedHashMap.Values()
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
	var allJsInputs commandLineRunner_jsInputs
	if sortAlphabetically {
		allJsInputs = jutil.NewTreeMap[string, string](jutil.StringCompareTo)
	} else {
		allJsInputs = commandLineRunner_linkedJsInputs{jutil.NewLinkedHashMap[string, string]()}
	}
	excludes := jutil.NewLinkedHashSet[string]()
	var wd *jutil.Path
	for _, pattern := range patterns {
		if !strings.Contains(pattern, "*") && !strings.HasPrefix(pattern, "!") {
			if jutil.File_IsDirectoryPath(pattern) {
				commandLineRunner_matchPaths(jutil.NewFileWithParent(jutil.NewFile(pattern), "**.js").String(), allJsInputs, excludes)
			} else {
				original := jutil.Path_Of(pattern)
				norm := original.Normalize()
				if !norm.IsAbsolute() && wd == nil {
					wd = jutil.Path_WorkingDir()
				}
				pathStringAbsolute := norm.ToAbsolutePathWithWd(wd).String()
				if !excludes.Contains(pathStringAbsolute) {
					allJsInputs.Put(pathStringAbsolute, original.String())
				}
			}
		} else {
			commandLineRunner_matchPaths(pattern, allJsInputs, excludes)
		}
	}

	var result []string
	for v := range allJsInputs.Values() {
		result = append(result, v)
	}
	return result
}

// Java: CommandLineRunner.matchPaths(String, Map<String, String>, Set<String>)
func commandLineRunner_matchPaths(pattern string, allJsInputs commandLineRunner_jsInputs, excludes *jutil.LinkedHashSet[string]) {
	remove := strings.IndexByte(pattern, '!') == 0
	if remove {
		pattern = pattern[1:]
	}

	separator := jutil.File_separator

	// Split the pattern into two pieces: the globbing part
	// and the non-globbing prefix.
	// Go-only: language-adapter — Guava Splitter.on / Joiner.on with the
	// ASCII separator are strings.Split / strings.Join.
	patternParts := strings.Split(pattern, jutil.File_separator)
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

	matcher := jutil.FileSystem_GetPathMatcher("glob:" + prefix + separator + pattern)
	jutil.Files_WalkFileTree(
		jutil.Path_Of(prefix),
		/* FOLLOW_LINKS */ true,
		math.MaxInt32,
		&commandLineRunner_matchPathsVisitor{matcher: matcher, remove: remove, allJsInputs: allJsInputs, excludes: excludes})
}

// commandLineRunner_matchPathsVisitor is matchPaths' anonymous
// SimpleFileVisitor<Path>.
type commandLineRunner_matchPathsVisitor struct {
	jutil.SimpleFileVisitor
	matcher     *jutil.PathMatcher
	remove      bool
	allJsInputs commandLineRunner_jsInputs
	excludes    *jutil.LinkedHashSet[string]
}

func (v *commandLineRunner_matchPathsVisitor) VisitFile(p *jutil.Path) jutil.FileVisitResult {
	if v.matcher.Matches(p) || v.matcher.Matches(p.Normalize()) {
		pathStringAbsolute := p.Normalize().ToAbsolutePath().String()
		if v.remove {
			v.excludes.Add(pathStringAbsolute)
			v.allJsInputs.Remove(pathStringAbsolute)
		} else if !v.excludes.Contains(pathStringAbsolute) {
			v.allJsInputs.Put(pathStringAbsolute, p.String())
		}
	}
	return jutil.FileVisitResult_CONTINUE
}

func (v *commandLineRunner_matchPathsVisitor) VisitFileFailed(file *jutil.Path, e jutil.IIOException) jutil.FileVisitResult {
	return jutil.FileVisitResult_SKIP_SUBTREE
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
