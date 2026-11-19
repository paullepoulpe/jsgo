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
// Port of com/google/javascript/jscomp/CommandLineRunner.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.CommandLineRunner
type CommandLineRunner struct {
	AbstractCommandLineRunnerBase
	flags       *CommandLineRunner_Flags        `gson:"flags"`
	errors      bool                            `gson:"errors"`
	runCompiler bool                            `gson:"runCompiler"`
	errorStream *j2gort.Ext_java_io_PrintStream `gson:"errorStream"`
	bundler     *ClosureBundler                 `gson:"bundler"`
}

func (this_ *CommandLineRunner) BindSelf_(s any) {
	_ = s
	this_.AbstractCommandLineRunnerBase.BindSelf_(s)
}
func (this_ *CommandLineRunner) IsJavaNil_() bool { return this_ == nil }
func (*CommandLineRunner) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CommandLineRunner", "CommandLineRunner"
}
func (this_ *CommandLineRunner) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractCommandLineRunnerBase.IdentityHash_()
}
func (*CommandLineRunner) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CommandLineRunner", "com.google.javascript.jscomp.CommandLineRunner", func(v_ any) bool {
		return j2gort.InstanceOf[*CommandLineRunner](v_)
	})
}
func (*CommandLineRunner) Implements_jscomp_AbstractCommandLineRunner() {}

func StaticGet_CommandLineRunner_OUTPUT_MARKER() string { return "%output%" }
func StaticGet_CommandLineRunner_UTF8_BOM_CODE() int32  { return int32(65279) }

var static_CommandLineRunner_extraChunkNameChars *j2gort.Ext_java_util_regex_Pattern

func StaticGet_CommandLineRunner_ExtraChunkNameChars() *j2gort.Ext_java_util_regex_Pattern {
	clinit_CommandLineRunner()
	return static_CommandLineRunner_extraChunkNameChars
}
func StaticSet_CommandLineRunner_ExtraChunkNameChars(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_CommandLineRunner()
	static_CommandLineRunner_extraChunkNameChars = v
}

var static_CommandLineRunner_phaseLogger *j2gort.Ext_java_util_logging_Logger

func StaticGet_CommandLineRunner_PhaseLogger() *j2gort.Ext_java_util_logging_Logger {
	clinit_CommandLineRunner()
	return static_CommandLineRunner_phaseLogger
}
func StaticSet_CommandLineRunner_PhaseLogger(v *j2gort.Ext_java_util_logging_Logger) {
	clinit_CommandLineRunner()
	static_CommandLineRunner_phaseLogger = v
}

var clinitDone_CommandLineRunner uint32
var clinitState_CommandLineRunner j2gort.ClinitState

func clinit_CommandLineRunner() {
	if atomic.LoadUint32(&clinitDone_CommandLineRunner) != 0 {
		return
	}
	clinitSlow_CommandLineRunner()
}

func clinitSlow_CommandLineRunner() {
	if !j2gort.EnterClinit(&clinitState_CommandLineRunner, &clinitDone_CommandLineRunner) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CommandLineRunner, &clinitDone_CommandLineRunner)
	static_CommandLineRunner_extraChunkNameChars = j2gort.Ext_java_util_regex_Pattern_Compile_1("[-.]+")
	static_CommandLineRunner_phaseLogger = j2gort.Ext_java_util_logging_Logger_GetLogger_1(j2gort.CheckNotNullPtr(j2gort.ClassLiteral("com.google.javascript.jscomp.PhaseOptimizer", "PhaseOptimizer")).GetName())
}

func (this_ *CommandLineRunner) initFields_CommandLineRunner() {
	_ = this_
	this_.flags = NewCommandLineRunner_Flags()

}

// Java constructor: CommandLineRunner(java.lang.String[] args)
func NewCommandLineRunner_1(args []string) *CommandLineRunner {
	clinit_CommandLineRunner()
	this_ := &CommandLineRunner{}
	this_.BindSelf_(this_)
	InitCommandLineRunner_1(this_, args)
	return this_
}

func InitCommandLineRunner_1(this_ *CommandLineRunner, args []string) {
	_ = args
	InitAbstractCommandLineRunner(&this_.AbstractCommandLineRunnerBase)
	this_.initFields_CommandLineRunner()
	this_.priv_CommandLineRunner_initConfigFromFlags_1406(args, j2gort.Ext_java_lang_System_F_Out, j2gort.Ext_java_lang_System_F_Err)
}

// Java constructor: CommandLineRunner(java.lang.String[] args, java.io.PrintStream out, java.io.PrintStream err)
func NewCommandLineRunner_3(args []string, out *j2gort.Ext_java_io_PrintStream, err *j2gort.Ext_java_io_PrintStream) *CommandLineRunner {
	clinit_CommandLineRunner()
	this_ := &CommandLineRunner{}
	this_.BindSelf_(this_)
	InitCommandLineRunner_3(this_, args, out, err)
	return this_
}

func InitCommandLineRunner_3(this_ *CommandLineRunner, args []string, out *j2gort.Ext_java_io_PrintStream, err *j2gort.Ext_java_io_PrintStream) {
	_ = args
	_ = out
	_ = err
	InitAbstractCommandLineRunner_2(&this_.AbstractCommandLineRunnerBase, out, err)
	this_.initFields_CommandLineRunner()
	this_.priv_CommandLineRunner_initConfigFromFlags_1406(args, out, err)
}

// Java constructor: CommandLineRunner(java.lang.String[] args, java.io.InputStream in, java.io.PrintStream out, java.io.PrintStream err)
func NewCommandLineRunner_4(args []string, in j2gort.Ext_java_io_InputStream, out *j2gort.Ext_java_io_PrintStream, err *j2gort.Ext_java_io_PrintStream) *CommandLineRunner {
	clinit_CommandLineRunner()
	this_ := &CommandLineRunner{}
	this_.BindSelf_(this_)
	InitCommandLineRunner_4(this_, args, in, out, err)
	return this_
}

func InitCommandLineRunner_4(this_ *CommandLineRunner, args []string, in j2gort.Ext_java_io_InputStream, out *j2gort.Ext_java_io_PrintStream, err *j2gort.Ext_java_io_PrintStream) {
	_ = args
	_ = in
	_ = out
	_ = err
	InitAbstractCommandLineRunner_3(&this_.AbstractCommandLineRunnerBase, in, out, err)
	this_.initFields_CommandLineRunner()
	this_.priv_CommandLineRunner_initConfigFromFlags_1406(args, out, err)
}

// Java static method: CommandLineRunner.processArgs(java.lang.String[] args)
func CommandLineRunner_ProcessArgs(args []string) j2gort.Ext_java_util_List {
	clinit_CommandLineRunner()
	_ = args
	var argPattern *j2gort.Ext_java_util_regex_Pattern = j2gort.Ext_java_util_regex_Pattern_Compile_2("(--?[a-zA-Z_]+)=(.*)", int32(32))
	_ = argPattern
	var quotesPattern *j2gort.Ext_java_util_regex_Pattern = j2gort.Ext_java_util_regex_Pattern_Compile_1("^['\"](.*)['\"]$")
	_ = quotesPattern
	var processedArgs []string = []string{}
	_ = processedArgs
	for _, el_1 := range j2gort.CheckNotNullSlice(args) {
		arg := el_1
		_ = arg
		var matcher *j2gort.Ext_java_util_regex_Matcher = argPattern.Matcher(j2gort.BoxString(arg))
		_ = matcher
		if matcher.Matches() {
			processedArgs = append(processedArgs, matcher.Group_1_int32(int32(1)))
			var value string = matcher.Group_1_int32(int32(2))
			_ = value
			var quotesMatcher *j2gort.Ext_java_util_regex_Matcher = quotesPattern.Matcher(j2gort.BoxString(value))
			_ = quotesMatcher
			if quotesMatcher.Matches() {
				processedArgs = append(processedArgs, quotesMatcher.Group_1_int32(int32(1)))
			} else {
				processedArgs = append(processedArgs, value)
			}
		} else {
			processedArgs = append(processedArgs, arg)
		}
	}
	return j2gort.WrapTypedArrayList(processedArgs)
}

// Java method: CommandLineRunner.reportError(java.lang.String message)
func (this_ *CommandLineRunner) priv_CommandLineRunner_reportError_1403(message string) {
	_ = message
	this_.errors = true
	this_.errorStream.Println_1_string(message)
	this_.errorStream.Flush()
}

// Java method: CommandLineRunner.processFlagFiles()
func (this_ *CommandLineRunner) priv_CommandLineRunner_processFlagFiles_1404() {
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.flags.flagFiles))
		for cur_2.HasNext() {
			flagFile := cur_2.NextStr()
			_ = flagFile
			func() {
				defer func() {
					if rec_3 := recover(); rec_3 != nil {
						jex_4 := javalang.RecoverException(rec_3)
						_ = jex_4
						if j2gort.InstanceOf[javaio.IIOException](jex_4) {
							ioErr := j2gort.AsIface[javaio.IIOException](jex_4)
							_ = ioErr
							this_.priv_CommandLineRunner_reportError_1403(("ERROR - " + j2gort.StrValString(flagFile) + " read error."))
						} else {
							panic(rec_3)
						}
					}
				}()
				this_.priv_CommandLineRunner_processFlagFile_1405(flagFile)
			}()
		}
	}
}

// Java method: CommandLineRunner.processFlagFile(java.lang.String flagFileString)
func (this_ *CommandLineRunner) priv_CommandLineRunner_processFlagFile_1405(flagFileString string) {
	_ = flagFileString
	var flagFile j2gort.Ext_java_nio_file_Path = j2gort.Ext_java_nio_file_Path_Of_2(flagFileString, []string{})
	_ = flagFile
	var buffer *j2gort.Ext_java_io_BufferedReader = j2gort.Ext_java_nio_file_Files_NewBufferedReader_2(flagFile, j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8)
	_ = buffer
	var builder *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = builder
	var tokens j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = tokens
	var quoted bool = false
	_ = quoted
	var escaped bool = false
	_ = escaped
	var isFirstCharacter bool = true
	_ = isFirstCharacter
	var c int32 = int32(0)
	_ = c
	for (func() int32 { var asgn_5 int32 = buffer.Read_0(); c = asgn_5; return asgn_5 }()) != int32(-1) {
		if isFirstCharacter {
			isFirstCharacter = false
			if c == int32(65279) {
				continue
			}
		}
		if (((c == int32(32)) || (c == int32(9))) || (c == int32(10))) || (c == int32(13)) {
			if quoted {
				builder.Append_1_uint16(uint16(c))
			} else if builder.Length() != int32(0) {
				tokens.Add_1(builder.String())
				builder.SetLength(int32(0))
			}
		} else if c == int32(34) {
			if escaped {
				if quoted {
					builder.SetCharAt((builder.Length() - int32(1)), uint16(c))
				} else {
					builder.Append_1_uint16(uint16(c))
				}
			} else {
				quoted = !(quoted)
			}
		} else {
			builder.Append_1_uint16(uint16(c))
		}
		escaped = (c == int32(92))
	}
	buffer.Close_0()
	if builder.Length() != int32(0) {
		tokens.Add_1(builder.String())
	}
	{
		r_6 := this_.flags
		r_6.flagFiles = j2gort.NewExt_java_util_ArrayList()
	}
	tokens = CommandLineRunner_ProcessArgs(j2gort.ColToNewStringArray(tokens))
	var previousGuardLevels j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](StaticGet_CommandLineRunner_Flags_GuardLevels()))
	_ = previousGuardLevels
	var previousMixedJsSources j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](StaticGet_CommandLineRunner_Flags_MixedJsSources()))
	_ = previousMixedJsSources
	StaticGet_CommandLineRunner_Flags_GuardLevels().Clear()
	StaticGet_CommandLineRunner_Flags_MixedJsSources().Clear()
	this_.flags.priv_Flags_parse_1411(tokens)
	StaticGet_CommandLineRunner_Flags_GuardLevels().AddAll_1(previousGuardLevels)
	StaticGet_CommandLineRunner_Flags_MixedJsSources().AddAll_1(previousMixedJsSources)
	if !(this_.flags.flagFiles.IsEmpty()) {
		this_.priv_CommandLineRunner_reportError_1403("ERROR - Arguments in the file cannot contain --flagfile option.")
	}
}

// Java method: CommandLineRunner.getVersionText()
func (this_ *CommandLineRunner) GetVersionText() string {
	return j2gort.Ext_java_lang_String_Join_2_any_slice_any("\n", []any{"Closure Compiler (http://github.com/google/closure-compiler)", "Version: " + "v20261104"})
}

// Java method: CommandLineRunner.initConfigFromFlags(java.lang.String[] args, java.io.PrintStream out, java.io.PrintStream err)
func (this_ *CommandLineRunner) priv_CommandLineRunner_initConfigFromFlags_1406(args []string, out *j2gort.Ext_java_io_PrintStream, err *j2gort.Ext_java_io_PrintStream) {
	_ = args
	_ = out
	_ = err
	this_.errorStream = err
	var processedArgs j2gort.Ext_java_util_List = CommandLineRunner_ProcessArgs(args)
	_ = processedArgs
	StaticGet_CommandLineRunner_Flags_GuardLevels().Clear()
	StaticGet_CommandLineRunner_Flags_MixedJsSources().Clear()
	var mixedSources j2gort.Ext_java_util_List = nil
	_ = mixedSources
	var mappings j2gort.Ext_java_util_List = nil
	_ = mappings
	var sourceMapInputs j2gort.Ext_com_google_common_collect_ImmutableMap = nil
	_ = sourceMapInputs
	var parseInlineSourceMaps bool = false
	_ = parseInlineSourceMaps
	var applyInputSourceMaps bool = false
	_ = applyInputSourceMaps
	func() {
		defer func() {
			if rec_7 := recover(); rec_7 != nil {
				jex_8 := javalang.RecoverException(rec_7)
				_ = jex_8
				if j2gort.InstanceOf[j2gort.Ext_org_kohsuke_args4j_CmdLineException](jex_8) {
					e := j2gort.AsIface[j2gort.Ext_org_kohsuke_args4j_CmdLineException](jex_8)
					_ = e
					this_.priv_CommandLineRunner_reportError_1403(j2gort.ThrowableGetMessage(e))
				} else if j2gort.InstanceOf[javaio.IIOException](jex_8) {
					ioErr := j2gort.AsIface[javaio.IIOException](jex_8)
					_ = ioErr
					this_.priv_CommandLineRunner_reportError_1403(("ERROR - ioException: " + j2gort.StrVal(ioErr)))
				} else {
					panic(rec_7)
				}
			}
		}()
		this_.flags.priv_Flags_parse_1411(processedArgs)
		this_.priv_CommandLineRunner_processFlagFiles_1404()
		mixedSources = this_.flags.GetMixedJsSources()
		mappings = this_.flags.GetSourceMapLocationMappings()
		sourceMapInputs = this_.flags.GetSourceMapInputs()
		parseInlineSourceMaps = j2gort.UnboxBoolean(this_.flags.parseInlineSourceMaps)
		applyInputSourceMaps = this_.flags.applyInputSourceMaps
	}()
	if this_.flags.processCommonJsModules {
		{
			r_9 := this_.flags
			r_9.processClosurePrimitives = true
		}
	}
	if (j2gort.UnboxInt(this_.flags.browserFeaturesetYear) != int32(0)) && (this_.flags.languageOut != "ECMASCRIPT_NEXT") {
		panic(NewFlagUsageException("ERROR - both flags `--browser_featureset_year` and `--language_out` specified."))
	}
	if this_.flags.outputWrapper == j2gort.NullString {
		{
			r_10 := this_.flags
			r_10.outputWrapper = ""
		}
	}
	if (this_.flags.outputWrapperFile != j2gort.NullString) && !(j2gort.CheckNotNullString(this_.flags.outputWrapperFile) == "") {
		func() {
			defer func() {
				if rec_11 := recover(); rec_11 != nil {
					jex_12 := javalang.RecoverException(rec_11)
					_ = jex_12
					if j2gort.InstanceOf[javalang.IException](jex_12) {
						e_2 := j2gort.AsIface[javalang.IException](jex_12)
						_ = e_2
						this_.priv_CommandLineRunner_reportError_1403("ERROR - invalid output_wrapper_file specified.")
					} else {
						panic(rec_11)
					}
				}
			}()
			{
				r_13 := this_.flags
				r_13.outputWrapper = j2gort.Ext_com_google_common_io_Files_AsCharSource(j2gort.NewExt_java_io_File_1_string(this_.flags.outputWrapperFile), j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8).Read()
			}
		}()
	}
	if !(j2gort.CheckNotNullString(this_.flags.outputWrapper) == "") && !(j2gort.JStrIndexOf(this_.flags.outputWrapper, "%output%") >= 0) {
		this_.priv_CommandLineRunner_reportError_1403(("ERROR - invalid output_wrapper specified. Missing '" + "%output%" + "'."))
	}
	if !(j2gort.CheckNotNullString(this_.flags.outputWrapper) == "") && (this_.flags.isolationMode != StaticGet_CompilerOptions_IsolationMode_NONE()) {
		this_.priv_CommandLineRunner_reportError_1403("--output_wrapper and --isolation_mode may not be used together.")
	}
	if this_.flags.isolationMode == StaticGet_CompilerOptions_IsolationMode_IIFE() {
		{
			r_14 := this_.flags
			r_14.outputWrapper = "(function(){%output%}).call(this);"
		}
	}
	var bundleFiles j2gort.Ext_com_google_common_collect_ImmutableList = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	_ = bundleFiles
	var skipNormalOutputs bool = false
	_ = skipNormalOutputs
	if this_.flags.compilationLevelParsed == StaticGet_CompilationLevel_BUNDLE() {
		if j2gort.CheckNotNullString(this_.flags.jsOutputFile) == "" {
			this_.priv_CommandLineRunner_reportError_1403("--compilation_level=BUNDLE cannot be used without a --js_output_file.")
		} else {
			bundleFiles = j2gort.ImmListOfStr1(this_.flags.jsOutputFile)
			{
				r_15 := this_.flags
				r_15.jsOutputFile = ""
			}
			skipNormalOutputs = true
		}
	}
	var conv CodingConvention = nil
	_ = conv
	if this_.flags.thirdParty {
		conv = CodingConventions_GetDefault()
	} else if this_.flags.chromePass {
		conv = NewChromeCodingConvention()
	} else {
		conv = NewClosureCodingConvention()
	}
	if !(this_.flags.renaming) && (this_.flags.compilationLevelParsed == StaticGet_CompilationLevel_ADVANCED_OPTIMIZATIONS()) {
		this_.priv_CommandLineRunner_reportError_1403("ERROR - renaming cannot be disabled when ADVANCED_OPTIMIZATIONS is used.")
	}
	var dependencyOptions *DependencyOptions = nil
	_ = dependencyOptions
	func() {
		defer func() {
			if rec_16 := recover(); rec_16 != nil {
				jex_17 := javalang.RecoverException(rec_16)
				_ = jex_17
				if j2gort.InstanceOf[*FlagUsageException](jex_17) {
					e_3 := j2gort.AsIface[*FlagUsageException](jex_17)
					_ = e_3
					this_.priv_CommandLineRunner_reportError_1403(j2gort.ThrowableGetMessage(e_3))
				} else {
					panic(rec_16)
				}
			}
		}()
		dependencyOptions = DependencyOptions_FromFlags(this_.flags.dependencyMode, this_.flags.entryPoint, j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()), j2gort.NullString, false, false)
	}()
	if this_.errors {
		CommandLineRunner_Flags_PrintShortUsageAfterErrors(this_.errorStream)
	} else if this_.flags.displayHelp || this_.flags.helpMarkdown {
		this_.flags.priv_Flags_printUsage_1412(out)
	} else {
		this_.runCompiler = true
		var config *AbstractCommandLineRunner_CommandLineConfig = this_.GetCommandLineConfig()
		_ = config
		config.SetPrintVersion(this_.flags.version).SetPrintTree(this_.flags.printTree).SetPrintTreeJson(this_.flags.printTreeJson).SetPrintAst(this_.flags.printAst).SetJscompDevMode(this_.flags.jscompDevMode).SetLoggingLevel(this_.flags.loggingLevel).SetExterns(this_.flags.externs).SetMixedJsSources(mixedSources).SetDefaultToStdin().SetJsOutputFile(this_.flags.jsOutputFile).SetChunk(this_.flags.chunk).SetVariableMapOutputFile(this_.flags.variableMapOutputFile).SetCreateNameMapFiles(this_.flags.createNameMapFiles).SetPropertyMapOutputFile(this_.flags.propertyMapOutputFile).SetExpectedDiagnostics(this_.flags.expectedDiagnostics).SetPropertyMapInputFile(this_.flags.propertyMapInputFile).SetVariableMapInputFile(this_.flags.variableMapInputFile).SetInstrumentationMappingFile(this_.flags.instrumentationMappingOutputFile).SetCodingConvention(conv).SetSummaryDetailLevel(this_.flags.summaryDetailLevel).SetOutputWrapper(this_.flags.outputWrapper).SetChunkWrapper(this_.flags.chunkWrapper).SetChunkOutputPathPrefix(this_.flags.chunkOutputPathPrefix).SetCreateSourceMap(this_.flags.createSourceMap).SetSourceMapFormat(this_.flags.sourceMapFormat).SetSourceMapLocationMappings(mappings).SetSourceMapInputFiles(j2gort.AsIface[j2gort.Ext_java_util_Map](sourceMapInputs)).SetParseInlineSourceMaps(parseInlineSourceMaps).SetApplyInputSourceMaps(applyInputSourceMaps).SetWarningGuards(StaticGet_CommandLineRunner_Flags_GuardLevels()).SetDefine(this_.flags.define).SetBrowserFeaturesetYear(this_.flags.browserFeaturesetYear).SetCharset(this_.flags.charset).SetDependencyOptions(dependencyOptions).SetOutputManifest(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfStr1(this_.flags.outputManifest))).SetOutputBundle(j2gort.AsIface[j2gort.Ext_java_util_List](bundleFiles)).SetSkipNormalOutputs(skipNormalOutputs).SetOutputChunkDependencies(this_.flags.outputChunkDependencies).SetProcessCommonJSModules(this_.flags.processCommonJsModules).SetModuleRoots(this_.flags.moduleRoot).SetWarningsAllowlistFile(this_.flags.warningsAllowlistFile).SetHideWarningsFor(this_.flags.hideWarningsFor).SetAngularPass(this_.flags.angularPass).SetJsonStreamMode(this_.flags.jsonStreamMode).SetErrorFormat(this_.flags.errorFormat)
		var segmentOfCompilationToRun CompilerOptions_SegmentOfCompilationToRun = this_.flags.segmentOfCompilationToRun
		_ = segmentOfCompilationToRun
		switch (segmentOfCompilationToRun).Ordinal() {
		case int32(1) /* SegmentOfCompilationToRun.CHECKS */ :
			if !(this_.flags.filenameToRestoreFrom == j2gort.NullString) {
				j2gort.CheckStateStr(false, "Cannot restore and run CHECKS segment of compilation")
			}
			config.SetSaveCompilationStateToFilename(this_.flags.filenameToSaveTo, int32(1))
		case int32(4) /* SegmentOfCompilationToRun.OPTIMIZATIONS */ :
			config.SetContinueSavedCompilationFileName(this_.flags.filenameToRestoreFrom, int32(1))
			config.SetSaveCompilationStateToFilename(this_.flags.filenameToSaveTo, int32(2))
		case int32(6) /* SegmentOfCompilationToRun.FINALIZATIONS */ :
			if !(this_.flags.filenameToSaveTo == j2gort.NullString) {
				j2gort.CheckStateStr(false, "Cannot run FINALIZATIONS segment of compilation and then save the result")
			}
			config.SetContinueSavedCompilationFileName(this_.flags.filenameToRestoreFrom, int32(2))
		case int32(0) /* SegmentOfCompilationToRun.ENTIRE_COMPILATION */ :
		default:
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Cannot run %s segment of compilation: "+j2gort.StrVal(this_.flags.segmentOfCompilationToRun))))
		}
	}
	this_.errorStream = nil
}

// Java method: CommandLineRunner.addAllowlistWarningsGuard(com.google.javascript.jscomp.CompilerOptions options, java.io.File allowlistFile)
func (this_ *CommandLineRunner) AddAllowlistWarningsGuard(options *CompilerOptions, allowlistFile *j2gort.Ext_java_io_File) {
	_ = options
	_ = allowlistFile
	options.AddWarningsGuard(ToIface_WarningsGuard(AllowlistWarningsGuard_FromFile(allowlistFile)))
}

// Java method: CommandLineRunner.checkChunkName(java.lang.String name)
func (this_ *CommandLineRunner) CheckChunkName(name string) {
	_ = name
	if !(pkg_rhino.TokenStream_IsJSIdentifier(StaticGet_CommandLineRunner_ExtraChunkNameChars().Matcher(j2gort.BoxString(name)).ReplaceAll_1_string("_"))) {
		panic(NewFlagUsageException(("Invalid chunk name: '" + j2gort.StrValString(name) + "'")))
	}
}

// Java method: CommandLineRunner.createOptions()
func (this_ *CommandLineRunner) CreateOptions() *CompilerOptions {
	clinit_CompilerOptions()
	var val_18 CompilerOptions
	InitCompilerOptions(&val_18)
	var options *CompilerOptions = &val_18
	_ = options
	if !(j2gort.CheckNotNullString(this_.flags.languageIn) == "") {
		var languageMode CompilerOptions_LanguageMode = CompilerOptions_LanguageMode_FromString(this_.flags.languageIn)
		_ = languageMode
		if languageMode == StaticGet_CompilerOptions_LanguageMode_UNSUPPORTED() {
			panic(NewFlagUsageException("Cannot specify the unsupported set of features for language_in."))
		}
		if languageMode != 0 {
			options.SetLanguageIn(languageMode)
		} else {
			panic(NewFlagUsageException(("Unknown language `" + j2gort.StrValString(this_.flags.languageIn) + "' specified. Expected one of: " + j2gort.StrVal(CompilerOptions_LanguageMode_ValidCommandLineNames()))))
		}
	}
	var languageMode_2 CompilerOptions_LanguageMode = CompilerOptions_LanguageMode_FromString(this_.flags.languageOut)
	_ = languageMode_2
	if languageMode_2 == StaticGet_CompilerOptions_LanguageMode_UNSUPPORTED() {
		panic(NewFlagUsageException("Cannot specify the unsupported set of features for language_out."))
	}
	if languageMode_2 != 0 {
		options.SetLanguageOut(languageMode_2)
	} else {
		panic(NewFlagUsageException(("Unknown language `" + j2gort.StrValString(this_.flags.languageOut) + "' specified. Expected one of: " + j2gort.StrVal(CompilerOptions_LanguageMode_ValidCommandLineNames()))))
	}
	options.SetCodingConvention(NewClosureCodingConvention())
	options.SetExtraAnnotationNames(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.flags.extraAnnotationName))
	var level CompilationLevel = this_.flags.compilationLevelParsed
	_ = level
	level.SetOptionsForCompilationLevel(options)
	if this_.flags.debug {
		level.SetDebugOptionsForCompilationLevel(options)
	}
	options.SetNumParallelThreads(this_.flags.numParallelThreads)
	options.SetEnvironment(this_.flags.environment)
	options.SetChecksOnly(this_.flags.checksOnly)
	if this_.flags.checksOnly {
		options.SetOutputJs(StaticGet_CompilerOptions_OutputJs_NONE())
	}
	options.SetIncrementalChecks(this_.flags.incrementalCheckMode)
	options.SetContinueAfterErrors(this_.flags.continueAfterErrors)
	if this_.flags.useTypesForOptimization {
		level.SetTypeBasedOptimizationOptions(options)
	}
	if (this_.flags.assumeFunctionWrapper || (this_.flags.isolationMode == StaticGet_CompilerOptions_IsolationMode_IIFE())) || (this_.flags.chunkOutputType == StaticGet_CompilerOptions_ChunkOutputType_ES_MODULES()) {
		level.SetWrappedOutputOptimizations(options)
		options.SetAssumeGlobalScopeIsIsolated(true)
	}
	if this_.flags.typedAstOutputFile != j2gort.NullString {
		options.SetTypedAstOutputFile(j2gort.Ext_java_nio_file_Path_Of_2(this_.flags.typedAstOutputFile, []string{}))
	}
	options.SetGenerateExports(this_.flags.generateExports)
	options.SetExportLocalPropertyDefinitions(this_.flags.exportLocalPropertyDefinitions)
	var wLevel WarningLevel = this_.flags.warningLevel
	_ = wLevel
	wLevel.SetOptionsForWarningLevel(options)
	{
		cur_19 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.flags.formatting))
		for cur_19.HasNext() {
			rawElem_20 := cur_19.Next()
			formattingOption := j2gort.Cast[CommandLineRunner_FormattingOption](rawElem_20)
			_ = formattingOption
			formattingOption.priv_FormattingOption_applyToOptions_1427(options)
		}
	}
	options.SetClosurePass(this_.flags.processClosurePrimitives)
	options.SetAngularPass(this_.flags.angularPass)
	options.SetPolymerVersion(this_.flags.polymerVersion)
	options.SetChromePass(this_.flags.chromePass)
	if !(j2gort.CheckNotNullString(this_.flags.j2clPassMode) == "") {
		func() {
			defer func() {
				if rec_21 := recover(); rec_21 != nil {
					jex_22 := javalang.RecoverException(rec_21)
					_ = jex_22
					if j2gort.InstanceOf[javalang.IIllegalArgumentException](jex_22) {
						ex := j2gort.AsIface[javalang.IIllegalArgumentException](jex_22)
						_ = ex
						panic(NewFlagUsageException(("Unknown J2clPassMode `" + j2gort.StrValString(this_.flags.j2clPassMode) + "' specified.")))
					} else {
						panic(rec_21)
					}
				}
			}()
			var j2clPassMode CompilerOptions_J2clPassMode = CompilerOptions_J2clPassMode_ValueOf(j2gort.Ext_com_google_common_base_Ascii_ToUpperCase_1_string(this_.flags.j2clPassMode))
			_ = j2clPassMode
			options.SetJ2clPass(j2clPassMode)
		}()
	}
	options.SetRemoveJ2clAsserts(this_.flags.removeJ2cLAsserts)
	options.SetRenamePrefix(this_.flags.renamePrefix)
	options.SetRenamePrefixNamespace(this_.flags.renamePrefixNamespace)
	options.SetPreserveTypeAnnotations(this_.flags.preserveTypeAnnotations)
	options.SetPreventLibraryInjection(!(this_.flags.injectLibraries))
	if !(this_.flags.forceInjectLibraries.IsEmpty()) {
		options.SetForceLibraryInjection(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.flags.forceInjectLibraries))
	}
	options.SetRewritePolyfills((this_.flags.rewritePolyfills && options.GetLanguageIn().ToFeatureSet().Contains_1_FeatureSet(pkg_parser.StaticGet_FeatureSet_ES2015())))
	options.SetIsolatePolyfills(this_.flags.isolatePolyfills)
	if !(j2gort.CheckNotNullString(this_.flags.translationsFile) == "") {
		func() {
			defer func() {
				if rec_23 := recover(); rec_23 != nil {
					jex_24 := javalang.RecoverException(rec_23)
					_ = jex_24
					if j2gort.InstanceOf[javaio.IIOException](jex_24) {
						e := j2gort.AsIface[javaio.IIOException](jex_24)
						_ = e
						panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "Reading XTB file", e)))
					} else {
						panic(rec_23)
					}
				}
			}()
			options.SetMessageBundle(NewXtbMessageBundle(j2gort.NewExt_java_io_FileInputStream_1_string(this_.flags.translationsFile), this_.flags.translationsProject))
		}()
	} else if StaticGet_CompilationLevel_ADVANCED_OPTIMIZATIONS() == level {
		options.SetMessageBundle(NewEmptyMessageBundle())
		options.SetWarningLevel(StaticGet_DiagnosticGroups_MSG_CONVENTIONS(), StaticGet_CheckLevel_OFF())
	}
	options.SetConformanceConfigs(j2gort.AsIface[j2gort.Ext_java_util_List](CommandLineRunner_LoadConformanceConfigs(this_.flags.conformanceConfigs)))
	options.SetPrintSourceAfterEachPass(this_.flags.printSourceAfterEachPass)
	options.SetTracerMode(this_.flags.tracerMode)
	options.SetStrictModeInput(this_.flags.strictModeInput)
	options.SetEmitUseStrict(this_.flags.emitUseStrict)
	options.SetSourceMapIncludeSourcesContent(this_.flags.sourceMapIncludeSourcesContent)
	options.SetModuleResolutionMode(this_.flags.moduleResolutionMode)
	options.SetBrowserResolverPrefixReplacements(j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(this_.flags.browserResolverPrefixReplacements))
	if this_.flags.packageJsonEntryNames != j2gort.NullString {
		func() {
			defer func() {
				if rec_25 := recover(); rec_25 != nil {
					jex_26 := javalang.RecoverException(rec_25)
					_ = jex_26
					if j2gort.InstanceOf[j2gort.Ext_org_kohsuke_args4j_CmdLineException](jex_26) {
						e_2 := j2gort.AsIface[j2gort.Ext_org_kohsuke_args4j_CmdLineException](jex_26)
						_ = e_2
						this_.priv_CommandLineRunner_reportError_1403("ERROR - invalid package_json_entry_names format specified.")
					} else {
						panic(rec_25)
					}
				}
			}()
			var packageJsonEntryNames j2gort.Ext_java_util_List = this_.flags.GetPackageJsonEntryNames()
			_ = packageJsonEntryNames
			options.SetPackageJsonEntryNames(packageJsonEntryNames)
		}()
	}
	if !(this_.flags.renaming) {
		options.SetVariableRenaming(StaticGet_VariableRenamingPolicy_OFF())
		options.SetPropertyRenaming(StaticGet_PropertyRenamingPolicy_OFF())
	}
	if (this_.flags.instrumentCodeParsed == StaticGet_CompilerOptions_InstrumentOption_PRODUCTION()) && j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(this_.flags.instrumentationMappingOutputFile) {
		panic(NewFlagUsageException("Expected --instrument_mapping_report to be set when --instrument_for_coverage_option is set to Production"))
	}
	if !(j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(this_.flags.instrumentationMappingOutputFile)) && (this_.flags.instrumentCodeParsed != StaticGet_CompilerOptions_InstrumentOption_PRODUCTION()) {
		panic(NewFlagUsageException("Expected --instrument_for_coverage_option to be passed with PRODUCTION when --instrument_mapping_report is set"))
	}
	if j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(this_.flags.productionInstrumentationArrayName) && (this_.flags.instrumentCodeParsed == StaticGet_CompilerOptions_InstrumentOption_PRODUCTION()) {
		panic(NewFlagUsageException("Expected --production_instrumentation_array_name to be set when --instrument_for_coverage_option is set to Production"))
	}
	options.SetInstrumentForCoverageOption(this_.flags.instrumentCodeParsed)
	options.SetProductionInstrumentationArrayName(this_.flags.productionInstrumentationArrayName)
	options.SetAllowDynamicImport(this_.flags.allowDynamicImport)
	options.SetDynamicImportAlias(this_.flags.dynamicImportAlias)
	options.SetAssumeStaticInheritanceIsNotUsed(this_.flags.assumeStaticInheritanceIsNotUsed)
	options.SetCrossChunkCodeMotionNoStubMethods(this_.flags.assumeNoPrototypeMethodEnumeration)
	if this_.flags.chunkOutputType == StaticGet_CompilerOptions_ChunkOutputType_ES_MODULES() {
		if this_.flags.renamePrefixNamespace != j2gort.NullString {
			panic(NewFlagUsageException("Expected --rename_prefix_namespace not to be specified when --chunk_output_type is set to ES_MODULES."))
		}
		if this_.flags.emitUseStrict {
			panic(NewFlagUsageException("Expected --emit_use_strict should not be specified when --chunk_output_type is set to ES_MODULES."))
		}
		options.SetChunkOutputType(this_.flags.chunkOutputType)
		options.SetEmitUseStrict(false)
		if level == StaticGet_CompilationLevel_ADVANCED_OPTIMIZATIONS() {
			options.SetExtractPrototypeMemberDeclarations_1_CompilerOptions_ExtractPrototypeMemberDeclarationsMode(StaticGet_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_CHUNK_TEMP())
		}
	}
	return options
}

// Java method: CommandLineRunner.createCompiler()
func (this_ *CommandLineRunner) CreateCompiler() *Compiler {
	return NewCompiler_1_Ext_java_io_PrintStream(this_.GetErrorPrintStream())
}

// Java method: CommandLineRunner.getBundler()
func (this_ *CommandLineRunner) priv_CommandLineRunner_getBundler_1407() *ClosureBundler {
	if this_.bundler != nil {
		return this_.bundler
	}
	var moduleRoots j2gort.Ext_com_google_common_collect_ImmutableList = nil
	_ = moduleRoots
	if !(this_.flags.moduleRoot.IsEmpty()) {
		moduleRoots = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.flags.moduleRoot))
	} else {
		moduleRoots = j2gort.ImmListOfStr1("./")
	}
	var options *CompilerOptions = this_.CreateOptions()
	_ = options
	return func() *ClosureBundler {
		var asgn_27 *ClosureBundler = NewClosureBundler_2(StaticGet_Transpile_Transpiler_NULL(), NewBaseTranspiler(NewBaseTranspiler_CompilerSupplier_4(StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_NEXT().ToFeatureSet().Without_2(pkg_parser.StaticGet_FeatureSet_Feature_MODULES(), []*pkg_parser.FeatureSet_Feature{}), options.GetModuleResolutionMode(), moduleRoots, options.GetBrowserResolverPrefixReplacements()), ""))
		this_.bundler = asgn_27
		return asgn_27
	}()
}

// Java method: CommandLineRunner.prepForBundleAndAppendTo(java.lang.Appendable out, com.google.javascript.jscomp.CompilerInput input, java.lang.String content)
func (this_ *CommandLineRunner) PrepForBundleAndAppendTo(out j2gort.Ext_java_lang_Appendable, input *CompilerInput, content string) {
	_ = out
	_ = input
	_ = content
	this_.priv_CommandLineRunner_getBundler_1407().WithPath(input.GetName()).AppendTo_3_Ext_java_lang_Appendable_DependencyInfo_string(out, ToIface_DependencyInfo(input), content)
}

// Java method: CommandLineRunner.appendRuntimeTo(java.lang.Appendable out)
func (this_ *CommandLineRunner) AppendRuntimeTo(out j2gort.Ext_java_lang_Appendable) {
	_ = out
	this_.priv_CommandLineRunner_getBundler_1407().AppendRuntimeTo(out)
}

// Java method: CommandLineRunner.createExterns(com.google.javascript.jscomp.CompilerOptions options)
func (this_ *CommandLineRunner) CreateExterns(options *CompilerOptions) j2gort.Ext_java_util_List {
	_ = options
	var externs j2gort.Ext_java_util_List = this_.AbstractCommandLineRunnerBase.CreateExterns(options)
	_ = externs
	if this_.IsInTestMode() {
		return externs
	} else {
		var builtinExterns j2gort.Ext_java_util_List = AbstractCommandLineRunner_GetBuiltinExterns(options.GetEnvironment())
		_ = builtinExterns
		builtinExterns.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](externs))
		return builtinExterns
	}
}

// Java static method: CommandLineRunner.loadConformanceConfigs(java.util.List<java.lang.String> configPaths)
func CommandLineRunner_LoadConformanceConfigs(configPaths j2gort.Ext_java_util_List) j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_CommandLineRunner()
	_ = configPaths
	var configs *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	_ = configs
	{
		cur_28 := j2gort.NewIterCursor(j2gort.CheckNotNull(configPaths))
		for cur_28.HasNext() {
			configPath := cur_28.NextStr()
			_ = configPath
			func() {
				defer func() {
					if rec_29 := recover(); rec_29 != nil {
						jex_30 := javalang.RecoverException(rec_29)
						_ = jex_30
						if j2gort.InstanceOf[javaio.IIOException](jex_30) {
							e := j2gort.AsIface[javaio.IIOException](jex_30)
							_ = e
							panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "Error loading conformance config", e)))
						} else {
							panic(rec_29)
						}
					}
				}()
				configs.AddVoid_(j2gort.PtrToAny(CommandLineRunner_LoadConformanceConfig(configPath)))
			}()
		}
	}
	return configs.BuildImmList_()
}

// Java static method: CommandLineRunner.loadConformanceConfig(java.lang.String configFile)
func CommandLineRunner_LoadConformanceConfig(configFile string) *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig {
	clinit_CommandLineRunner()
	_ = configFile
	var textProto string = j2gort.Ext_com_google_common_io_Files_AsCharSource(j2gort.NewExt_java_io_File_1_string(configFile), j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8).Read()
	_ = textProto
	var builder *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_Builder = j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_NewBuilder_0()
	_ = builder
	if !(j2gort.CheckNotNullString(textProto) == "") && (int32(j2gort.JStrCharAt(textProto, int32(0))) == int32(65279)) {
		textProto = j2gort.JStrSubFrom(textProto, int32(1))
	}
	func() {
		defer func() {
			if rec_31 := recover(); rec_31 != nil {
				jex_32 := javalang.RecoverException(rec_31)
				_ = jex_32
				if j2gort.InstanceOf[javalang.IException](jex_32) {
					e := j2gort.AsIface[javalang.IException](jex_32)
					_ = e
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
				} else {
					panic(rec_31)
				}
			}
		}()
		j2gort.Ext_com_google_protobuf_TextFormat_Merge_2_any_Ext_com_google_protobuf_Message_Builder(j2gort.BoxString(textProto), j2gort.AsIface[j2gort.Ext_com_google_protobuf_Message_Builder](builder))
	}()
	return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig](builder.Build())
}

// Java static method: CommandLineRunner.getDefaultExterns()
func CommandLineRunner_GetDefaultExterns() j2gort.Ext_java_util_List {
	clinit_CommandLineRunner()
	return AbstractCommandLineRunner_GetBuiltinExterns(StaticGet_CompilerOptions_Environment_BROWSER())
}

// Java static method: CommandLineRunner.findJsFiles(java.util.Collection<java.lang.String> patterns)
func CommandLineRunner_FindJsFiles_1(patterns j2gort.Ext_java_util_Collection) j2gort.Ext_java_util_List {
	clinit_CommandLineRunner()
	_ = patterns
	return CommandLineRunner_FindJsFiles_2(patterns, false)
}

// Java static method: CommandLineRunner.findJsFiles(java.util.Collection<java.lang.String> patterns, boolean sortAlphabetically)
func CommandLineRunner_FindJsFiles_2(patterns j2gort.Ext_java_util_Collection, sortAlphabetically bool) j2gort.Ext_java_util_List {
	clinit_CommandLineRunner()
	_ = patterns
	_ = sortAlphabetically
	var allJsInputs j2gort.Ext_java_util_Map
	_ = allJsInputs
	if sortAlphabetically {
		allJsInputs = j2gort.NewExt_java_util_TreeMap()
	} else {
		allJsInputs = j2gort.NewExt_java_util_LinkedHashMap()
	}
	var excludes j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	_ = excludes
	{
		cur_33 := j2gort.NewIterCursor(j2gort.CheckNotNull(patterns))
		for cur_33.HasNext() {
			pattern := cur_33.NextStr()
			_ = pattern
			if !(j2gort.JStrIndexOf(pattern, "*") >= 0) && !(javalang.StringStartsWith(j2gort.CheckNotNullString(pattern), "!")) {
				var matchedFile *j2gort.Ext_java_io_File = j2gort.NewExt_java_io_File_1_string(pattern)
				_ = matchedFile
				if matchedFile.IsDirectory() {
					CommandLineRunner_MatchPaths(j2gort.NewExt_java_io_File_2_Ext_java_io_File_string(matchedFile, "**.js").String(), allJsInputs, excludes)
				} else {
					var original j2gort.Ext_java_nio_file_Path = j2gort.Ext_java_nio_file_Path_Of_2(pattern, []string{})
					_ = original
					var pathStringAbsolute string = original.Normalize().ToAbsolutePath().String()
					_ = pathStringAbsolute
					if !(j2gort.SetContainsStr(excludes, pathStringAbsolute)) {
						j2gort.MapPutStr(allJsInputs, pathStringAbsolute, j2gort.BoxString(original.String()))
					}
				}
			} else {
				CommandLineRunner_MatchPaths(pattern, allJsInputs, excludes)
			}
		}
	}
	return j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](allJsInputs.Values()))
}

// Java static method: CommandLineRunner.matchPaths(java.lang.String pattern, java.util.Map<java.lang.String,java.lang.String> allJsInputs, java.util.Set<java.lang.String> excludes)
func CommandLineRunner_MatchPaths(pattern string, allJsInputs j2gort.Ext_java_util_Map, excludes j2gort.Ext_java_util_Set) {
	clinit_CommandLineRunner()
	_ = pattern
	_ = allJsInputs
	_ = excludes
	var fs j2gort.Ext_java_nio_file_FileSystem = j2gort.Ext_java_nio_file_FileSystems_GetDefault()
	_ = fs
	var remove bool = (j2gort.JStrIndexOfCp(pattern, int32(uint16(33))) == int32(0))
	_ = remove
	if remove {
		pattern = j2gort.JStrSubFrom(pattern, int32(1))
	}
	var separator string
	_ = separator
	if j2gort.CheckNotNullString(j2gort.Ext_java_io_File_F_Separator) == "\\" {
		separator = "\\\\"
	} else {
		separator = j2gort.Ext_java_io_File_F_Separator
	}
	var patternParts j2gort.Ext_java_util_List = j2gort.SplitterSplitToListStr(j2gort.Ext_com_google_common_base_Splitter_On_1_string(j2gort.Ext_java_io_File_F_Separator), pattern)
	_ = patternParts
	var prefix string = "."
	_ = prefix
	{
		var i int32 = int32(0)
		_ = i
		for ; i < patternParts.Size(); i = i + 1 {
			if j2gort.JStrIndexOf(j2gort.ListGetStr(patternParts, i), "*") >= 0 {
				if i > int32(0) {
					prefix = j2gort.Ext_com_google_common_base_Joiner_On_1_string(separator).Join_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_util_List](patternParts.SubList(int32(0), i)))
					pattern = j2gort.Ext_com_google_common_base_Joiner_On_1_string(separator).Join_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_util_List](patternParts.SubList(i, patternParts.Size())))
				}
				break
			}
		}
	}
	var matcher j2gort.Ext_java_nio_file_PathMatcher = fs.GetPathMatcher(("glob:" + j2gort.StrValString(prefix) + j2gort.StrValString(separator) + j2gort.StrValString(pattern)))
	_ = matcher
	j2gort.Ext_java_nio_file_Files_WalkFileTree_4(fs.GetPath(prefix, []string{}), j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_java_util_EnumSet_Of_1(j2gort.AsIface[j2gort.Ext_java_lang_Enum](j2gort.Ext_java_nio_file_FileVisitOption_F_FOLLOW_LINKS))), int32(2147483647), NewCommandLineRunner__Anon_11(matcher, remove, excludes, allJsInputs))
}

// Java method: CommandLineRunner.shouldRunCompiler()
func (this_ *CommandLineRunner) ShouldRunCompiler() bool {
	return this_.runCompiler
}

// Java method: CommandLineRunner.hasErrors()
func (this_ *CommandLineRunner) HasErrors() bool {
	return this_.errors
}

// Java static method: CommandLineRunner.main(java.lang.String[] args)
func CommandLineRunner_Main(args []string) {
	clinit_CommandLineRunner()
	_ = args
	if StaticGet_CommandLineRunner_PhaseLogger() != nil {
		StaticGet_CommandLineRunner_PhaseLogger().SetLevel(j2gort.Ext_java_util_logging_Level_F_OFF)
	}
	var runner *CommandLineRunner = NewCommandLineRunner_1(args)
	_ = runner
	if runner.ShouldRunCompiler() {
		runner.Run()
	}
	if runner.HasErrors() {
		j2gort.Ext_java_lang_System_Exit(int32(-1))
	}
}

// Java class: com.google.javascript.jscomp.CommandLineRunner.Flags
type CommandLineRunner_Flags struct {
	browserFeaturesetYear              *j2gort.BoxInt                                                `gson:"browserFeaturesetYear"`
	displayHelp                        bool                                                          `gson:"displayHelp"`
	printTree                          bool                                                          `gson:"printTree"`
	printTreeJson                      bool                                                          `gson:"printTreeJson"`
	printAst                           bool                                                          `gson:"printAst"`
	emitUseStrict                      bool                                                          `gson:"emitUseStrict"`
	strictModeInput                    bool                                                          `gson:"strictModeInput"`
	jscompDevMode                      CompilerOptions_DevMode                                       `gson:"jscompDevMode"`
	loggingLevel                       string                                                        `gson:"loggingLevel"`
	externs                            j2gort.Ext_java_util_List                                     `gson:"externs"`
	js                                 j2gort.Ext_java_util_List                                     `gson:"js"`
	unusedJsZip                        j2gort.Ext_java_util_List                                     `gson:"unusedJsZip"`
	jsOutputFile                       string                                                        `gson:"jsOutputFile"`
	chunk                              j2gort.Ext_java_util_List                                     `gson:"chunk"`
	filenameToSaveTo                   string                                                        `gson:"filenameToSaveTo"`
	filenameToRestoreFrom              string                                                        `gson:"filenameToRestoreFrom"`
	segmentOfCompilationToRun          CompilerOptions_SegmentOfCompilationToRun                     `gson:"segmentOfCompilationToRun"`
	variableMapOutputFile              string                                                        `gson:"variableMapOutputFile"`
	instrumentationMappingOutputFile   string                                                        `gson:"instrumentationMappingOutputFile"`
	createNameMapFiles                 bool                                                          `gson:"createNameMapFiles"`
	sourceMapIncludeSourcesContent     bool                                                          `gson:"sourceMapIncludeSourcesContent"`
	propertyMapOutputFile              string                                                        `gson:"propertyMapOutputFile"`
	thirdParty                         bool                                                          `gson:"thirdParty"`
	summaryDetailLevel                 int32                                                         `gson:"summaryDetailLevel"`
	isolationMode                      CompilerOptions_IsolationMode                                 `gson:"isolationMode"`
	outputWrapper                      string                                                        `gson:"outputWrapper"`
	outputWrapperFile                  string                                                        `gson:"outputWrapperFile"`
	chunkWrapper                       j2gort.Ext_java_util_List                                     `gson:"chunkWrapper"`
	chunkOutputPathPrefix              string                                                        `gson:"chunkOutputPathPrefix"`
	createSourceMap                    string                                                        `gson:"createSourceMap"`
	sourceMapFormat                    *SourceMap_Format                                             `gson:"sourceMapFormat"`
	sourceMapLocationMapping           j2gort.Ext_java_util_List                                     `gson:"sourceMapLocationMapping"`
	sourceMapInputs                    j2gort.Ext_java_util_List                                     `gson:"sourceMapInputs"`
	parseInlineSourceMaps              *j2gort.BoxBoolean                                            `gson:"parseInlineSourceMaps"`
	applyInputSourceMaps               bool                                                          `gson:"applyInputSourceMaps"`
	jscompError                        j2gort.Ext_java_util_List                                     `gson:"jscompError"`
	jscompWarning                      j2gort.Ext_java_util_List                                     `gson:"jscompWarning"`
	jscompOff                          j2gort.Ext_java_util_List                                     `gson:"jscompOff"`
	define                             j2gort.Ext_java_util_List                                     `gson:"define"`
	charset                            string                                                        `gson:"charset"`
	compilationLevel                   string                                                        `gson:"compilationLevel"`
	compilationLevelParsed             CompilationLevel                                              `gson:"compilationLevelParsed"`
	numParallelThreads                 int32                                                         `gson:"numParallelThreads"`
	checksOnly                         bool                                                          `gson:"checksOnly"`
	incrementalCheckMode               CompilerOptions_IncrementalCheckMode                          `gson:"incrementalCheckMode"`
	continueAfterErrors                bool                                                          `gson:"continueAfterErrors"`
	useTypesForOptimization            bool                                                          `gson:"useTypesForOptimization"`
	assumeFunctionWrapper              bool                                                          `gson:"assumeFunctionWrapper"`
	warningLevel                       WarningLevel                                                  `gson:"warningLevel"`
	debug                              bool                                                          `gson:"debug"`
	typedAstOutputFile                 string                                                        `gson:"typedAstOutputFile"`
	generateExports                    bool                                                          `gson:"generateExports"`
	exportLocalPropertyDefinitions     bool                                                          `gson:"exportLocalPropertyDefinitions"`
	formatting                         j2gort.Ext_java_util_List                                     `gson:"formatting"`
	processCommonJsModules             bool                                                          `gson:"processCommonJsModules"`
	moduleRoot                         j2gort.Ext_java_util_List                                     `gson:"moduleRoot"`
	processClosurePrimitives           bool                                                          `gson:"processClosurePrimitives"`
	angularPass                        bool                                                          `gson:"angularPass"`
	polymerVersion                     *j2gort.BoxInt                                                `gson:"polymerVersion"`
	chromePass                         bool                                                          `gson:"chromePass"`
	j2clPassMode                       string                                                        `gson:"j2clPassMode"`
	removeJ2cLAsserts                  bool                                                          `gson:"removeJ2cLAsserts"`
	outputManifest                     string                                                        `gson:"outputManifest"`
	outputChunkDependencies            string                                                        `gson:"outputChunkDependencies"`
	languageIn                         string                                                        `gson:"languageIn"`
	languageOut                        string                                                        `gson:"languageOut"`
	version                            bool                                                          `gson:"version"`
	translationsFile                   string                                                        `gson:"translationsFile"`
	translationsProject                string                                                        `gson:"translationsProject"`
	flagFiles                          j2gort.Ext_java_util_List                                     `gson:"flagFiles"`
	warningsAllowlistFile              string                                                        `gson:"warningsAllowlistFile"`
	hideWarningsFor                    j2gort.Ext_java_util_List                                     `gson:"hideWarningsFor"`
	extraAnnotationName                j2gort.Ext_java_util_List                                     `gson:"extraAnnotationName"`
	tracerMode                         CompilerOptions_TracerMode                                    `gson:"tracerMode"`
	renamePrefix                       string                                                        `gson:"renamePrefix"`
	renamePrefixNamespace              string                                                        `gson:"renamePrefixNamespace"`
	conformanceConfigs                 j2gort.Ext_java_util_List                                     `gson:"conformanceConfigs"`
	environment                        CompilerOptions_Environment                                   `gson:"environment"`
	jsonStreamMode                     CompilerOptions_JsonStreamMode                                `gson:"jsonStreamMode"`
	preserveTypeAnnotations            bool                                                          `gson:"preserveTypeAnnotations"`
	injectLibraries                    bool                                                          `gson:"injectLibraries"`
	forceInjectLibraries               j2gort.Ext_java_util_List                                     `gson:"forceInjectLibraries"`
	dependencyMode                     DependencyOptions_DependencyMode                              `gson:"dependencyMode"`
	entryPoint                         j2gort.Ext_java_util_List                                     `gson:"entryPoint"`
	rewritePolyfills                   bool                                                          `gson:"rewritePolyfills"`
	isolatePolyfills                   bool                                                          `gson:"isolatePolyfills"`
	printSourceAfterEachPass           bool                                                          `gson:"printSourceAfterEachPass"`
	moduleResolutionMode               ModuleLoader_ResolutionMode                                   `gson:"moduleResolutionMode"`
	browserResolverPrefixReplacements  j2gort.Ext_java_util_Map                                      `gson:"browserResolverPrefixReplacements"`
	packageJsonEntryNames              string                                                        `gson:"packageJsonEntryNames"`
	errorFormat                        AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption `gson:"errorFormat"`
	renaming                           bool                                                          `gson:"renaming"`
	helpMarkdown                       bool                                                          `gson:"helpMarkdown"`
	instrumentForCoverageOption        string                                                        `gson:"instrumentForCoverageOption"`
	productionInstrumentationArrayName string                                                        `gson:"productionInstrumentationArrayName"`
	chunkOutputType                    CompilerOptions_ChunkOutputType                               `gson:"chunkOutputType"`
	instrumentCodeParsed               CompilerOptions_InstrumentOption                              `gson:"instrumentCodeParsed"`
	allowDynamicImport                 bool                                                          `gson:"allowDynamicImport"`
	dynamicImportAlias                 string                                                        `gson:"dynamicImportAlias"`
	assumeStaticInheritanceIsNotUsed   bool                                                          `gson:"assumeStaticInheritanceIsNotUsed"`
	assumeNoPrototypeMethodEnumeration bool                                                          `gson:"assumeNoPrototypeMethodEnumeration"`
	variableMapInputFile               string                                                        `gson:"variableMapInputFile"`
	propertyMapInputFile               string                                                        `gson:"propertyMapInputFile"`
	expectedDiagnostics                j2gort.Ext_java_util_List                                     `gson:"expectedDiagnostics"`
	arguments                          j2gort.Ext_java_util_List                                     `gson:"arguments"`
	parser                             *j2gort.Ext_org_kohsuke_args4j_CmdLineParser                  `gson:"parser"`
	maxLineLength                      int32                                                         `gson:"maxLineLength"`
	idHash_                            atomic.Int32
}

func (this_ *CommandLineRunner_Flags) BindSelf_(s any) {
	_ = s
}
func (this_ *CommandLineRunner_Flags) IsJavaNil_() bool { return this_ == nil }
func (*CommandLineRunner_Flags) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CommandLineRunner$Flags", "Flags"
}
func (this_ *CommandLineRunner_Flags) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CommandLineRunner_Flags) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CommandLineRunner.Flags", "com.google.javascript.jscomp.CommandLineRunner$Flags", func(v_ any) bool {
		return j2gort.InstanceOf[*CommandLineRunner_Flags](v_)
	})
}

var static_CommandLineRunner_Flags_guardLevels j2gort.Ext_java_util_List

func StaticGet_CommandLineRunner_Flags_GuardLevels() j2gort.Ext_java_util_List {
	clinit_CommandLineRunner_Flags()
	return static_CommandLineRunner_Flags_guardLevels
}
func StaticSet_CommandLineRunner_Flags_GuardLevels(v j2gort.Ext_java_util_List) {
	clinit_CommandLineRunner_Flags()
	static_CommandLineRunner_Flags_guardLevels = v
}

var static_CommandLineRunner_Flags_mixedJsSources j2gort.Ext_java_util_List

func StaticGet_CommandLineRunner_Flags_MixedJsSources() j2gort.Ext_java_util_List {
	clinit_CommandLineRunner_Flags()
	return static_CommandLineRunner_Flags_mixedJsSources
}
func StaticSet_CommandLineRunner_Flags_MixedJsSources(v j2gort.Ext_java_util_List) {
	clinit_CommandLineRunner_Flags()
	static_CommandLineRunner_Flags_mixedJsSources = v
}

var static_CommandLineRunner_Flags_categories j2gort.Ext_com_google_common_collect_ImmutableMultimap

func StaticGet_CommandLineRunner_Flags_Categories() j2gort.Ext_com_google_common_collect_ImmutableMultimap {
	clinit_CommandLineRunner_Flags()
	return static_CommandLineRunner_Flags_categories
}
func StaticSet_CommandLineRunner_Flags_Categories(v j2gort.Ext_com_google_common_collect_ImmutableMultimap) {
	clinit_CommandLineRunner_Flags()
	static_CommandLineRunner_Flags_categories = v
}
func StaticGet_CommandLineRunner_Flags_BOLD_PREFIX() string   { return "\x1b[1m" }
func StaticGet_CommandLineRunner_Flags_NORMAL_PREFIX() string { return "\x1b[0m" }
func StaticGet_CommandLineRunner_Flags_MARKDOWN_CHARS_TO_ESCAPE() string {
	return "[-*\\`\\[\\]{}\\(\\)#+\\.!<>]"
}

var static_CommandLineRunner_Flags_WHITESPACE_PATTERN *j2gort.Ext_java_util_regex_Pattern

func StaticGet_CommandLineRunner_Flags_WHITESPACE_PATTERN() *j2gort.Ext_java_util_regex_Pattern {
	clinit_CommandLineRunner_Flags()
	return static_CommandLineRunner_Flags_WHITESPACE_PATTERN
}
func StaticSet_CommandLineRunner_Flags_WHITESPACE_PATTERN(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_CommandLineRunner_Flags()
	static_CommandLineRunner_Flags_WHITESPACE_PATTERN = v
}

var clinitDone_CommandLineRunner_Flags uint32
var clinitState_CommandLineRunner_Flags j2gort.ClinitState

func clinit_CommandLineRunner_Flags() {
	if atomic.LoadUint32(&clinitDone_CommandLineRunner_Flags) != 0 {
		return
	}
	clinitSlow_CommandLineRunner_Flags()
}

func clinitSlow_CommandLineRunner_Flags() {
	if !j2gort.EnterClinit(&clinitState_CommandLineRunner_Flags, &clinitDone_CommandLineRunner_Flags) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CommandLineRunner_Flags, &clinitDone_CommandLineRunner_Flags)
	static_CommandLineRunner_Flags_guardLevels = j2gort.Ext_java_util_Collections_SynchronizedList(j2gort.NewExt_java_util_ArrayList())
	static_CommandLineRunner_Flags_mixedJsSources = j2gort.Ext_java_util_Collections_SynchronizedList(j2gort.NewExt_java_util_ArrayList())
	static_CommandLineRunner_Flags_categories = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.NewExt_com_google_common_collect_ImmutableMultimap_Builder().PutAll_2_any_Ext_java_lang_Iterable("Basic Usage", j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_8("compilation_level", "env", "externs", "js", "js_output_file", "language_in", "language_out", "warning_level")))).PutAll_2_any_Ext_java_lang_Iterable("Warning and Error Management", j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_9("conformance_configs", "error_format", "extra_annotation_name", "hide_warnings_for", "jscomp_error", "jscomp_off", "jscomp_warning", "strict_mode_input", "warnings_allowlist_file")))).PutAll_2_any_Ext_java_lang_Iterable("Output", j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_10("assume_function_wrapper", "debug", "emit_use_strict", "export_local_property_definitions", "formatting", "generate_exports", "isolation_mode", "output_wrapper", "output_wrapper_file", "rename_variable_prefix")))).PutAll_2_any_Ext_java_lang_Iterable("Dependency Management", j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.ImmListOfStr2("dependency_mode", "entry_point")))).PutAll_2_any_Ext_java_lang_Iterable("JS Modules", j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_5("dynamic_import_alias", "js_module_root", "module_resolution", "process_common_js_modules", "package_json_entry_names")))).PutAll_2_any_Ext_java_lang_Iterable("Library and Framework Specific", j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_7("angular_pass", "force_inject_library", "inject_libraries", "polymer_version", "process_closure_primitives", "rewrite_polyfills", "isolate_polyfills")))).PutAll_2_any_Ext_java_lang_Iterable("Code Splitting", j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_5("chunk", "chunk_output_path_prefix", "chunk_output_type", "chunk_wrapper", "rename_prefix_namespace")))).PutAll_2_any_Ext_java_lang_Iterable("Reports", j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_8("create_source_map", "output_manifest", "output_chunk_dependencies", "property_renaming_report", "source_map_input", "source_map_include_content", "source_map_location_mapping", "variable_renaming_report")))).PutAll_2_any_Ext_java_lang_Iterable("Miscellaneous", j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_11("assume_static_inheritance_is_not_used", "browser_featureset_year", "charset", "checks_only", "define", "flagfile", "help", "json_streams", "third_party", "use_types_for_optimization", "version")))).Build())
	static_CommandLineRunner_Flags_WHITESPACE_PATTERN = j2gort.Ext_java_util_regex_Pattern_Compile_1("\\s")
}

func (this_ *CommandLineRunner_Flags) initFields_CommandLineRunner_Flags() {
	_ = this_
	this_.browserFeaturesetYear = j2gort.ValueOfInt(int32(0))

	this_.strictModeInput = true
	this_.jscompDevMode = StaticGet_CompilerOptions_DevMode_OFF()
	this_.loggingLevel = j2gort.Ext_java_util_logging_Level_F_WARNING.GetName()
	this_.externs = j2gort.NewExt_java_util_ArrayList()
	this_.js = j2gort.NewExt_java_util_ArrayList()

	this_.chunk = j2gort.NewExt_java_util_ArrayList()
	this_.filenameToSaveTo = j2gort.NullString
	this_.filenameToRestoreFrom = j2gort.NullString
	this_.segmentOfCompilationToRun = StaticGet_CompilerOptions_SegmentOfCompilationToRun_ENTIRE_COMPILATION()

	this_.summaryDetailLevel = int32(1)
	this_.isolationMode = StaticGet_CompilerOptions_IsolationMode_NONE()

	this_.chunkWrapper = j2gort.NewExt_java_util_ArrayList()
	this_.chunkOutputPathPrefix = "./"

	this_.sourceMapFormat = StaticGet_SourceMap_Format_DEFAULT()
	this_.sourceMapLocationMapping = j2gort.NewExt_java_util_ArrayList()
	this_.sourceMapInputs = j2gort.NewExt_java_util_ArrayList()
	this_.parseInlineSourceMaps = j2gort.ValueOfBoolean(true)
	this_.applyInputSourceMaps = true
	this_.jscompError = j2gort.NewExt_java_util_ArrayList()
	this_.jscompWarning = j2gort.NewExt_java_util_ArrayList()
	this_.jscompOff = j2gort.NewExt_java_util_ArrayList()
	this_.define = j2gort.NewExt_java_util_ArrayList()

	this_.compilationLevel = "SIMPLE"

	this_.numParallelThreads = int32(1)

	this_.incrementalCheckMode = StaticGet_CompilerOptions_IncrementalCheckMode_OFF()

	this_.useTypesForOptimization = true

	this_.warningLevel = StaticGet_WarningLevel_DEFAULT()

	this_.typedAstOutputFile = j2gort.NullString
	this_.generateExports = true
	this_.exportLocalPropertyDefinitions = true
	this_.formatting = j2gort.NewExt_java_util_ArrayList()

	this_.moduleRoot = j2gort.NewExt_java_util_ArrayList()
	this_.processClosurePrimitives = true

	this_.j2clPassMode = "AUTO"
	this_.removeJ2cLAsserts = true

	this_.languageIn = "STABLE"
	this_.languageOut = "ECMASCRIPT_NEXT"

	this_.translationsProject = j2gort.NullString
	this_.flagFiles = j2gort.NewExt_java_util_ArrayList()

	this_.hideWarningsFor = j2gort.NewExt_java_util_ArrayList()
	this_.extraAnnotationName = j2gort.NewExt_java_util_ArrayList()
	this_.tracerMode = StaticGet_CompilerOptions_TracerMode_OFF()
	this_.renamePrefix = j2gort.NullString
	this_.renamePrefixNamespace = j2gort.NullString
	this_.conformanceConfigs = j2gort.NewExt_java_util_ArrayList()
	this_.environment = StaticGet_CompilerOptions_Environment_BROWSER()
	this_.jsonStreamMode = StaticGet_CompilerOptions_JsonStreamMode_NONE()

	this_.injectLibraries = true
	this_.forceInjectLibraries = j2gort.NewExt_java_util_ArrayList()

	this_.entryPoint = j2gort.NewExt_java_util_ArrayList()
	this_.rewritePolyfills = true

	this_.moduleResolutionMode = StaticGet_ModuleLoader_ResolutionMode_BROWSER()
	this_.browserResolverPrefixReplacements = j2gort.NewExt_java_util_LinkedHashMap()
	this_.packageJsonEntryNames = j2gort.NullString
	this_.errorFormat = StaticGet_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_STANDARD()
	this_.renaming = true

	this_.instrumentForCoverageOption = "NONE"

	this_.chunkOutputType = StaticGet_CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE()
	this_.instrumentCodeParsed = StaticGet_CompilerOptions_InstrumentOption_NONE()
	this_.allowDynamicImport = true
	this_.dynamicImportAlias = j2gort.NullString
	this_.assumeStaticInheritanceIsNotUsed = true

	this_.expectedDiagnostics = j2gort.NewExt_java_util_ArrayList()
	this_.arguments = j2gort.NewExt_java_util_ArrayList()
	this_.maxLineLength = int32(80)
}

// Java constructor: Flags()
func NewCommandLineRunner_Flags() *CommandLineRunner_Flags {
	clinit_CommandLineRunner_Flags()
	this_ := &CommandLineRunner_Flags{}
	InitCommandLineRunner_Flags(this_)
	return this_
}

func InitCommandLineRunner_Flags(this_ *CommandLineRunner_Flags) {
	this_.initFields_CommandLineRunner_Flags()
	this_.parser = j2gort.NewExt_org_kohsuke_args4j_CmdLineParser_1(j2gort.PtrToAny(this_))
}

// Java method: Flags.setOption(java.lang.String name, java.lang.String value)
func (this_ *CommandLineRunner_Flags) SetOption(name string, value string) bool {
	_ = name
	_ = value
	switch j2gort.CheckNotNullString(name) {
	case "--js":
		this_.js.Add_1(j2gort.BoxString(value))
		return true
	case "--compilation_level":
		this_.compilationLevel = value
		return true
	case "--warning_level":
		func() {
			defer func() {
				if rec_34 := recover(); rec_34 != nil {
					jex_35 := javalang.RecoverException(rec_34)
					_ = jex_35
					if j2gort.InstanceOf[javalang.IIllegalArgumentException](jex_35) {
						e := j2gort.AsIface[javalang.IIllegalArgumentException](jex_35)
						_ = e
						panic(j2gort.NewExt_org_kohsuke_args4j_CmdLineException_2_Ext_org_kohsuke_args4j_CmdLineParser_Throwable(this_.parser, j2gort.AsIface[javalang.Throwable](e)))
					} else {
						panic(rec_34)
					}
				}
			}()
			this_.warningLevel = WarningLevel_ValueOf(value)
		}()
		return true
	case "--language_in":
		this_.languageIn = value
		return true
	case "--language_out":
		this_.languageOut = value
		return true
	case "--js_output_file":
		this_.jsOutputFile = value
		return true
	default:
		return false
	}
	panic("unreachable")
}

// Java method: Flags.parse(java.util.List<java.lang.String> args)
func (this_ *CommandLineRunner_Flags) priv_Flags_parse_1411(args j2gort.Ext_java_util_List) {
	_ = args
	this_.parser.ParseArgument_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](args))
	this_.compilationLevelParsed = CompilationLevel_FromString(j2gort.Ext_com_google_common_base_Ascii_ToUpperCase_1_string(this_.compilationLevel))
	if this_.compilationLevelParsed == 0 {
		panic(j2gort.NewExt_org_kohsuke_args4j_CmdLineException_2_Ext_org_kohsuke_args4j_CmdLineParser_string(this_.parser, ("Bad value for --compilation_level: " + j2gort.StrValString(this_.compilationLevel))))
	}
	this_.instrumentCodeParsed = CompilerOptions_InstrumentOption_FromString(j2gort.Ext_com_google_common_base_Ascii_ToUpperCase_1_string(this_.instrumentForCoverageOption))
	if this_.instrumentCodeParsed == 0 {
		panic(j2gort.NewExt_org_kohsuke_args4j_CmdLineException_2_Ext_org_kohsuke_args4j_CmdLineParser_string(this_.parser, ("Bad value for --instrument_for_coverage_option: " + j2gort.StrValString(this_.instrumentForCoverageOption))))
	}
}

// Java method: Flags.printUsage(java.io.PrintStream ps)
func (this_ *CommandLineRunner_Flags) priv_Flags_printUsage_1412(ps *j2gort.Ext_java_io_PrintStream) {
	_ = ps
	var outputStream *j2gort.Ext_java_io_OutputStreamWriter = j2gort.NewExt_java_io_OutputStreamWriter_2_Ext_java_io_OutputStream_Ext_java_nio_charset_Charset(j2gort.AsIface[j2gort.Ext_java_io_OutputStream](ps), j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8)
	_ = outputStream
	var isFirst bool = true
	_ = isFirst
	{
		cur_36 := j2gort.NewMapEntryCursor(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](StaticGet_CommandLineRunner_Flags_Categories().AsMap()))
		for cur_36.NextEntry() {
			var prefix string = "\n\n"
			_ = prefix
			var suffix string = ""
			_ = suffix
			if isFirst {
				isFirst = false
				prefix = ""
			}
			if j2gort.CheckNotNullString(cur_36.KeyStr()) == "Warning and Error Management" {
				if this_.helpMarkdown {
					suffix = ("\n## Available Error Groups\n\n  - " + j2gort.StrValString(javalang.StringReplace("accessControls, checkPrototypalTypes, checkRegExp, checkTypes, checkVars, conformanceViolations, const, constantProperty, deprecated, deprecatedAnnotations, defineWithoutGoogDefine, duplicateMessage, es5Strict, externsValidation, functionParams, globalThis, invalidCasts, lintVarDeclarations, misplacedTypeAnnotation, missingCallToSuper, missingOverride, missingPolyfill, missingProperties, missingProvide, missingRequire, missingReturn, missingSourcesWarnings, moduleLoad, moduleImport, msgDescriptions, nonStandardJsDocs, partialAlias, polymer, reportUnknownTypes, strictCheckTypes, strictMissingProperties, strictModuleDepCheck, strictPrimitiveOperators, suspiciousCode, typeInvalidation, undefinedVars, underscore, unknownDefines, unusedLocalVariables, uselessCode, untranspilableFeatures,visibility, weakModuleGet", ", ", "\n  - ")))
				} else {
					suffix = ("\n" + "\x1b[1m" + "Available Error Groups: " + "\x1b[0m" + "accessControls, checkPrototypalTypes, checkRegExp, checkTypes, checkVars, conformanceViolations, const, constantProperty, deprecated, deprecatedAnnotations, defineWithoutGoogDefine, duplicateMessage, es5Strict, externsValidation, functionParams, globalThis, invalidCasts, lintVarDeclarations, misplacedTypeAnnotation, missingCallToSuper, missingOverride, missingPolyfill, missingProperties, missingProvide, missingRequire, missingReturn, missingSourcesWarnings, moduleLoad, moduleImport, msgDescriptions, nonStandardJsDocs, partialAlias, polymer, reportUnknownTypes, strictCheckTypes, strictMissingProperties, strictModuleDepCheck, strictPrimitiveOperators, suspiciousCode, typeInvalidation, undefinedVars, underscore, unknownDefines, unusedLocalVariables, uselessCode, untranspilableFeatures,visibility, weakModuleGet")
				}
			}
			if this_.helpMarkdown {
				this_.maxLineLength = int32(5000)
				this_.parser.SetUsageWidth(this_.maxLineLength)
			}
			this_.priv_Flags_printCategoryUsage_1413(cur_36.KeyStr(), j2gort.AsIface[j2gort.Ext_java_util_Collection](cur_36.GetValue()), outputStream, prefix, suffix)
		}
	}
	ps.Flush()
}

// Java method: Flags.printCategoryUsage(java.lang.String categoryName, java.util.Collection<java.lang.String> options, java.io.OutputStreamWriter outputStream, java.lang.String prefix, java.lang.String suffix)
func (this_ *CommandLineRunner_Flags) priv_Flags_printCategoryUsage_1413(categoryName string, options j2gort.Ext_java_util_Collection, outputStream *j2gort.Ext_java_io_OutputStreamWriter, prefix string, suffix string) {
	_ = categoryName
	_ = options
	_ = outputStream
	_ = prefix
	_ = suffix
	func() {
		defer func() {
			if rec_37 := recover(); rec_37 != nil {
				jex_38 := javalang.RecoverException(rec_37)
				_ = jex_38
				if j2gort.InstanceOf[javaio.IIOException](jex_38) {
					e := j2gort.AsIface[javaio.IIOException](jex_38)
					_ = e
				} else {
					panic(rec_37)
				}
			}
		}()
		if prefix != j2gort.NullString {
			this_.priv_Flags_printStringLineWrapped_1414(prefix, outputStream)
		}
		if this_.helpMarkdown {
			outputStream.Write_1_string(("# " + j2gort.StrValString(categoryName) + "\n"))
			{
				cur_39 := j2gort.NewIterCursor(j2gort.CheckNotNull(options))
				for cur_39.HasNext() {
					optionName := cur_39.NextStr()
					_ = optionName
					var stringWriter *j2gort.Ext_java_io_StringWriter = j2gort.NewExt_java_io_StringWriter()
					_ = stringWriter
					this_.parser.PrintUsage_3(stringWriter, nil, j2gort.NewLambda_Ext_org_kohsuke_args4j_OptionHandlerFilter(func(optionHandler j2gort.Ext_org_kohsuke_args4j_spi_OptionHandler) bool {
						_ = optionHandler
						var namedOptionDef *j2gort.Ext_org_kohsuke_args4j_NamedOptionDef
						_ = namedOptionDef
						var ok_40 bool
						namedOptionDef, ok_40 = j2gort.TryCast[*j2gort.Ext_org_kohsuke_args4j_NamedOptionDef](optionHandler.GetExt_org_kohsuke_args4j_spi_OptionHandlerBase().F_Option)
						if ok_40 {
							return (!(optionHandler.GetExt_org_kohsuke_args4j_spi_OptionHandlerBase().F_Option.Hidden()) && func() bool {
								sr_41 := optionName
								se_42 := j2gort.JStrReplaceFirst(namedOptionDef.Name(), "^--", "")
								return j2gort.CheckNotNullString(sr_41) == se_42
							}())
						}
						return false
					}))
					stringWriter.Flush()
					var rawOptionUsage string = stringWriter.String()
					_ = rawOptionUsage
					var delimiterIndex int32 = j2gort.JStrIndexOf(rawOptionUsage, " : ")
					_ = delimiterIndex
					if delimiterIndex > int32(0) {
						outputStream.Write_1_string(("\n**" + j2gort.StrValString(javalang.StringTrim(j2gort.CheckNotNullString(j2gort.JStrSub(rawOptionUsage, int32(0), delimiterIndex)))) + "**  \n"))
						var optionDescription string = javalang.StringTrim(j2gort.CheckNotNullString(j2gort.JStrReplaceAll(j2gort.JStrSubFrom(rawOptionUsage, (delimiterIndex+int32(3))), "[-*\\`\\[\\]{}\\(\\)#+\\.!<>]", "\\\\$0")))
						_ = optionDescription
						outputStream.Write_1_string((j2gort.StrValString(optionDescription) + "\n"))
					} else {
						outputStream.Write_1_string(j2gort.JStrReplaceAll(rawOptionUsage, "[-*\\`\\[\\]{}\\(\\)#+\\.!<>]", "\\\\$0"))
					}
					outputStream.Flush()
				}
			}
		} else {
			outputStream.Write_1_string(("\x1b[1m" + j2gort.StrValString(categoryName) + ":\n" + "\x1b[0m"))
			this_.parser.PrintUsage_3(j2gort.AsIface[j2gort.Ext_java_io_Writer](outputStream), nil, j2gort.NewLambda_Ext_org_kohsuke_args4j_OptionHandlerFilter(func(optionHandler j2gort.Ext_org_kohsuke_args4j_spi_OptionHandler) bool {
				_ = optionHandler
				var namedOptionDef *j2gort.Ext_org_kohsuke_args4j_NamedOptionDef
				_ = namedOptionDef
				var ok_43 bool
				namedOptionDef, ok_43 = j2gort.TryCast[*j2gort.Ext_org_kohsuke_args4j_NamedOptionDef](optionHandler.GetExt_org_kohsuke_args4j_spi_OptionHandlerBase().F_Option)
				if ok_43 {
					return (!(optionHandler.GetExt_org_kohsuke_args4j_spi_OptionHandlerBase().F_Option.Hidden()) && options.Contains(j2gort.BoxString(j2gort.JStrReplaceFirst(namedOptionDef.Name(), "^--", ""))))
				}
				return false
			}))
		}
		if suffix != j2gort.NullString {
			this_.priv_Flags_printStringLineWrapped_1414(suffix, outputStream)
		}
	}()
}

// Java method: Flags.printStringLineWrapped(java.lang.String input, java.io.OutputStreamWriter outputStream)
func (this_ *CommandLineRunner_Flags) priv_Flags_printStringLineWrapped_1414(input string, outputStream *j2gort.Ext_java_io_OutputStreamWriter) {
	_ = input
	_ = outputStream
	if j2gort.JStrLen(input) < this_.maxLineLength {
		outputStream.Write_1_string(input)
		return
	}
	var endIndex int32 = this_.maxLineLength
	_ = endIndex
	var subString string = j2gort.JStrSub(input, int32(0), this_.maxLineLength)
	_ = subString
	var whitespaceMatcher *j2gort.Ext_java_util_regex_Matcher = StaticGet_CommandLineRunner_Flags_WHITESPACE_PATTERN().Matcher(j2gort.BoxString(subString))
	_ = whitespaceMatcher
	var foundMatch bool = false
	_ = foundMatch
	for whitespaceMatcher.Find_0() {
		endIndex = whitespaceMatcher.Start_0()
		foundMatch = true
	}
	outputStream.Write_1_string((j2gort.StrValString(j2gort.JStrSub(input, int32(0), endIndex)) + "\n"))
	this_.priv_Flags_printStringLineWrapped_1414(("    " + j2gort.StrValString(j2gort.JStrSubFrom(input, func() int32 {
		if foundMatch {
			return (endIndex + int32(1))
		}
		return endIndex
	}()))), outputStream)
}

// Java static method: Flags.printShortUsageAfterErrors(java.io.PrintStream ps)
func CommandLineRunner_Flags_PrintShortUsageAfterErrors(ps *j2gort.Ext_java_io_PrintStream) {
	clinit_CommandLineRunner_Flags()
	_ = ps
	ps.Print_1_string("Sample usage: ")
	ps.Println_1_string("--compilation_level (-O) VAL --externs VAL --js VAL --js_output_file VAL --warning_level (-W) [QUIET | DEFAULT | VERBOSE]")
	ps.Println_1_string("Run with --help for all options and details")
	ps.Flush()
}

// Java method: Flags.getMixedJsSources()
func (this_ *CommandLineRunner_Flags) GetMixedJsSources() j2gort.Ext_java_util_List {
	var mixedSources j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = mixedSources
	var excludes j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
	_ = excludes
	{
		cur_44 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_CommandLineRunner_Flags_MixedJsSources()))
		for cur_44.HasNext() {
			rawElem_45 := cur_44.Next()
			source := j2gort.CastPtr[AbstractCommandLineRunner_FlagEntry](rawElem_45)
			_ = source
			if javalang.StringEndsWith(j2gort.CheckNotNullString(source.GetValue()), ".zip") {
				mixedSources.Add_1(j2gort.PtrToAny(source))
			} else if javalang.StringStartsWith(j2gort.CheckNotNullString(source.GetValue()), "!") {
				{
					cur_46 := j2gort.NewIterCursor(j2gort.CheckNotNull(CommandLineRunner_FindJsFiles_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.ImmListOfStr1(j2gort.JStrSubFrom(source.GetValue(), int32(1)))))))
					for cur_46.HasNext() {
						filename := cur_46.NextStr()
						_ = filename
						excludes.Add(filename)
						mixedSources.Remove_1_any(NewAbstractCommandLineRunner_FlagEntry(j2gort.BoxNullableVal(StaticGet_AbstractCommandLineRunner_JsSourceType_JS()), filename))
					}
				}
			} else {
				{
					cur_47 := j2gort.NewIterCursor(j2gort.CheckNotNull(CommandLineRunner_FindJsFiles_2(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_java_util_Collections_SingletonList(j2gort.BoxString(source.GetValue()))), true)))
					for cur_47.HasNext() {
						filename_2 := cur_47.NextStr()
						_ = filename_2
						if !(excludes.Contains(filename_2)) {
							mixedSources.Add_1(NewAbstractCommandLineRunner_FlagEntry(j2gort.BoxNullableVal(StaticGet_AbstractCommandLineRunner_JsSourceType_JS()), filename_2))
						}
					}
				}
			}
		}
	}
	var fromArguments j2gort.Ext_java_util_List = CommandLineRunner_FindJsFiles_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.arguments))
	_ = fromArguments
	{
		cur_48 := j2gort.NewIterCursor(j2gort.CheckNotNull(fromArguments))
		for cur_48.HasNext() {
			filename_3 := cur_48.NextStr()
			_ = filename_3
			mixedSources.Add_1(NewAbstractCommandLineRunner_FlagEntry(j2gort.BoxNullableVal(StaticGet_AbstractCommandLineRunner_JsSourceType_JS()), filename_3))
		}
	}
	if (!(StaticGet_CommandLineRunner_Flags_MixedJsSources().IsEmpty()) && !(this_.arguments.IsEmpty())) && mixedSources.IsEmpty() {
		panic(j2gort.NewExt_org_kohsuke_args4j_CmdLineException_2_Ext_org_kohsuke_args4j_CmdLineParser_string(this_.parser, "No inputs matched"))
	}
	return mixedSources
}

// Java method: Flags.getSourceMapLocationMappings()
func (this_ *CommandLineRunner_Flags) GetSourceMapLocationMappings() j2gort.Ext_java_util_List {
	var locationMappings []*SourceMap_PrefixLocationMapping = []*SourceMap_PrefixLocationMapping{}
	_ = locationMappings
	var split j2gort.Ext_com_google_common_collect_ImmutableMap = this_.priv_Flags_splitPipeParts_1418(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.sourceMapLocationMapping), "--source_map_location_mapping")
	_ = split
	{
		cur_49 := j2gort.NewMapEntryCursor(split)
		for cur_49.NextEntry() {
			locationMappings = append(locationMappings, NewSourceMap_PrefixLocationMapping(cur_49.KeyStr(), j2gort.CastString(cur_49.GetValue())))
		}
	}
	return j2gort.ImmListFromTypedSlice(locationMappings)
}

// Java method: Flags.getSourceMapInputs()
func (this_ *CommandLineRunner_Flags) GetSourceMapInputs() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.priv_Flags_splitPipeParts_1418(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.sourceMapInputs), "--source_map_input")
}

// Java method: Flags.splitPipeParts(java.lang.Iterable<java.lang.String> input, java.lang.String flagName)
func (this_ *CommandLineRunner_Flags) priv_Flags_splitPipeParts_1418(input j2gort.Ext_java_lang_Iterable, flagName string) j2gort.Ext_com_google_common_collect_ImmutableMap {
	_ = input
	_ = flagName
	var result j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.NewExt_com_google_common_collect_ImmutableMap_Builder()
	_ = result
	var splitter *j2gort.Ext_com_google_common_base_Splitter = j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(124)).Limit(int32(2))
	_ = splitter
	{
		cur_50 := j2gort.NewIterCursor(j2gort.CheckNotNull(input))
		for cur_50.HasNext() {
			inputSourceMap := cur_50.NextStr()
			_ = inputSourceMap
			var parts j2gort.Ext_java_util_List = j2gort.SplitterSplitToListStr(splitter, inputSourceMap)
			_ = parts
			if parts.Size() != int32(2) {
				panic(j2gort.NewExt_org_kohsuke_args4j_CmdLineException_2_Ext_org_kohsuke_args4j_CmdLineParser_string(this_.parser, ("Bad value for " + j2gort.StrValString(flagName) + " (duplicate key): " + j2gort.StrVal(input))))
			}
			result.Put_2(j2gort.BoxString(j2gort.ListGetStr(parts, int32(0))), j2gort.BoxString(j2gort.ListGetStr(parts, int32(1))))
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](result.BuildOrThrow())
}

// Java method: Flags.getPackageJsonEntryNames()
func (this_ *CommandLineRunner_Flags) GetPackageJsonEntryNames() j2gort.Ext_java_util_List {
	return j2gort.SplitterSplitToListStr(j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(44)), this_.packageJsonEntryNames)
}

func (this_ *CommandLineRunner_Flags) RegisterArgs4jOptions(p *j2gort.Ext_org_kohsuke_args4j_CmdLineParser) {
	_ = this_
	_ = p
	{
		opt_0_ := j2gort.NewArgs4jNamedOptionDef("--browser_featureset_year", nil, "shortcut for defining goog.FEATURESET_YEAR=YYYY. The minimum valid value of the browser year is 2012", "", false, false, false, false)
		set_0_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.Integer", "Integer"), func(val_ any) {
			this_.browserFeaturesetYear = j2gort.ValueOfInt(j2gort.Args4jAsInt(p, val_))
		}, func() any {
			return j2gort.ToAny(this_.browserFeaturesetYear)
		})
		j2gort.Args4jRegisterHandler(p, opt_0_, j2gort.NewArgs4jBuiltinHandler(p, opt_0_, set_0_, "int", ""), false)
	}
	{
		opt_1_ := j2gort.NewArgs4jNamedOptionDef("--help", nil, "Displays this message on stdout and exit", "", false, false, false, false)
		set_1_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.displayHelp = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.displayHelp)
		})
		j2gort.Args4jRegisterHandler(p, opt_1_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_1_, set_1_), false)
	}
	{
		opt_2_ := j2gort.NewArgs4jNamedOptionDef("--print_tree", nil, "Prints out the parse tree and exits", "", false, false, true, false)
		set_2_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.printTree = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.printTree)
		})
		j2gort.Args4jRegisterHandler(p, opt_2_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_2_, set_2_), false)
	}
	{
		opt_3_ := j2gort.NewArgs4jNamedOptionDef("--print_tree_json", nil, "Prints out the parse tree as json and exits", "", false, false, true, false)
		set_3_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.printTreeJson = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.printTreeJson)
		})
		j2gort.Args4jRegisterHandler(p, opt_3_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_3_, set_3_), false)
	}
	{
		opt_4_ := j2gort.NewArgs4jNamedOptionDef("--print_ast", nil, "Prints a dot file describing the internal abstract syntax tree and exits", "", false, false, true, false)
		set_4_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.printAst = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.printAst)
		})
		j2gort.Args4jRegisterHandler(p, opt_4_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_4_, set_4_), false)
	}
	{
		opt_5_ := j2gort.NewArgs4jNamedOptionDef("--emit_use_strict", nil, "Start output with \"'use strict';\".", "", false, false, false, false)
		set_5_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.emitUseStrict = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.emitUseStrict)
		})
		j2gort.Args4jRegisterHandler(p, opt_5_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_5_, set_5_), false)
	}
	{
		opt_6_ := j2gort.NewArgs4jNamedOptionDef("--strict_mode_input", nil, "Assume input sources are to run in strict mode.", "", false, false, false, false)
		set_6_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.strictModeInput = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.strictModeInput)
		})
		j2gort.Args4jRegisterHandler(p, opt_6_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_6_, set_6_), false)
	}
	{
		opt_7_ := j2gort.NewArgs4jNamedOptionDef("--jscomp_dev_mode", []string{"--dev_mode"}, "", "", false, false, true, false)
		set_7_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("com.google.javascript.jscomp.CompilerOptions.DevMode", "DevMode"), func(val_ any) {
			this_.jscompDevMode = func() CompilerOptions_DevMode {
				_ = CompilerOptions_DevMode_Values()
				return j2gort.Cast[CompilerOptions_DevMode](j2gort.Args4jParseEnum(p, "com.google.javascript.jscomp.CompilerOptions.DevMode", val_))
			}()
		}, func() any {
			return this_.jscompDevMode
		})
		j2gort.Args4jRegisterHandler(p, opt_7_, func() j2gort.Ext_org_kohsuke_args4j_spi_OptionHandler {
			_ = CompilerOptions_DevMode_Values()
			return j2gort.NewArgs4jBuiltinHandler(p, opt_7_, set_7_, "enum", "com.google.javascript.jscomp.CompilerOptions.DevMode")
		}(), false)
	}
	{
		opt_8_ := j2gort.NewArgs4jNamedOptionDef("--logging_level", nil, "The logging level (standard java.util.logging.Level values) for Compiler progress. Does not control errors or warnings for the JavaScript code under compilation", "", false, false, true, false)
		set_8_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.loggingLevel = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.loggingLevel == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.loggingLevel)
		})
		j2gort.Args4jRegisterHandler(p, opt_8_, j2gort.NewArgs4jBuiltinHandler(p, opt_8_, set_8_, "string", ""), false)
	}
	{
		opt_9_ := j2gort.NewArgs4jNamedOptionDef("--externs", nil, "The file containing JavaScript externs. You may specify multiple", "", false, false, false, true)
		set_9_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.externs == nil {
				this_.externs = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.externs, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_9_, j2gort.NewArgs4jBuiltinHandler(p, opt_9_, set_9_, "string", ""), false)
	}
	{
		opt_10_ := j2gort.NewArgs4jNamedOptionDef("--js", nil, "The JavaScript filename. You may specify multiple. The flag name is optional, because args are interpreted as files by default. You may also use minimatch-style glob patterns. For example, use --js='**.js' --js='!**_test.js' to recursively include all js files that do not end in _test.js", "", false, false, false, true)
		set_10_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.js == nil {
				this_.js = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.js, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_10_, NewCommandLineRunner_Flags_JsOptionHandler(p, opt_10_, set_10_), false)
	}
	{
		opt_11_ := j2gort.NewArgs4jNamedOptionDef("--jszip", nil, "The JavaScript zip filename. You may specify multiple.", "", false, false, true, true)
		set_11_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.unusedJsZip == nil {
				this_.unusedJsZip = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.unusedJsZip, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_11_, NewCommandLineRunner_Flags_JsZipOptionHandler(p, opt_11_, set_11_), false)
	}
	{
		opt_12_ := j2gort.NewArgs4jNamedOptionDef("--js_output_file", nil, "Primary output filename. If not specified, output is written to stdout", "", false, false, false, false)
		set_12_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.jsOutputFile = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.jsOutputFile == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.jsOutputFile)
		})
		j2gort.Args4jRegisterHandler(p, opt_12_, j2gort.NewArgs4jBuiltinHandler(p, opt_12_, set_12_, "string", ""), false)
	}
	{
		opt_13_ := j2gort.NewArgs4jNamedOptionDef("--chunk", nil, "A JavaScript chunk specification. The format is <name>:<num-js-files>[:[<dep>,...][:]]]. Chunk names must be unique. Each dep is the name of a chunk that this chunk depends on. Chunks must be listed in dependency order, and JS source files must be listed in the corresponding order. Where --chunk flags occur in relation to --js flags is unimportant. <num-js-files> may be set to 'auto' for the first chunk if it has no dependencies. Provide the value 'auto' to trigger chunk creation from CommonJS modules.", "", false, false, false, true)
		set_13_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.chunk == nil {
				this_.chunk = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.chunk, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_13_, j2gort.NewArgs4jBuiltinHandler(p, opt_13_, set_13_, "string", ""), false)
	}
	{
		opt_14_ := j2gort.NewArgs4jNamedOptionDef("--filename_to_save_to", nil, "Filename to save state so that the compilation can be resumed later.", "", false, false, true, false)
		set_14_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.filenameToSaveTo = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.filenameToSaveTo == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.filenameToSaveTo)
		})
		j2gort.Args4jRegisterHandler(p, opt_14_, j2gort.NewArgs4jBuiltinHandler(p, opt_14_, set_14_, "string", ""), false)
	}
	{
		opt_15_ := j2gort.NewArgs4jNamedOptionDef("--filename_to_restore_from", nil, "Filename where a compilation state was previously saved.", "", false, false, true, false)
		set_15_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.filenameToRestoreFrom = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.filenameToRestoreFrom == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.filenameToRestoreFrom)
		})
		j2gort.Args4jRegisterHandler(p, opt_15_, j2gort.NewArgs4jBuiltinHandler(p, opt_15_, set_15_, "string", ""), false)
	}
	{
		opt_16_ := j2gort.NewArgs4jNamedOptionDef("--segment_of_compilation_to_run", nil, "Which segment of the compilation to run.", "", false, false, true, false)
		set_16_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("com.google.javascript.jscomp.CompilerOptions.SegmentOfCompilationToRun", "SegmentOfCompilationToRun"), func(val_ any) {
			this_.segmentOfCompilationToRun = func() CompilerOptions_SegmentOfCompilationToRun {
				_ = CompilerOptions_SegmentOfCompilationToRun_Values()
				return j2gort.Cast[CompilerOptions_SegmentOfCompilationToRun](j2gort.Args4jParseEnum(p, "com.google.javascript.jscomp.CompilerOptions.SegmentOfCompilationToRun", val_))
			}()
		}, func() any {
			return this_.segmentOfCompilationToRun
		})
		j2gort.Args4jRegisterHandler(p, opt_16_, func() j2gort.Ext_org_kohsuke_args4j_spi_OptionHandler {
			_ = CompilerOptions_SegmentOfCompilationToRun_Values()
			return j2gort.NewArgs4jBuiltinHandler(p, opt_16_, set_16_, "enum", "com.google.javascript.jscomp.CompilerOptions.SegmentOfCompilationToRun")
		}(), false)
	}
	{
		opt_17_ := j2gort.NewArgs4jNamedOptionDef("--variable_renaming_report", nil, "File where the serialized version of the variable renaming map produced should be saved", "", false, false, false, false)
		set_17_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.variableMapOutputFile = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.variableMapOutputFile == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.variableMapOutputFile)
		})
		j2gort.Args4jRegisterHandler(p, opt_17_, j2gort.NewArgs4jBuiltinHandler(p, opt_17_, set_17_, "string", ""), false)
	}
	{
		opt_18_ := j2gort.NewArgs4jNamedOptionDef("--instrument_mapping_report", nil, "File where the encoded parameters created by Production Instrumentation are mapped to their pre-encoded values. The %outname% placeholder will expand to the name of the output file that the source map corresponds to. Must be used in tandem with --instrument_for_coverage_option=PRODUCTION", "", false, false, false, false)
		set_18_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.instrumentationMappingOutputFile = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.instrumentationMappingOutputFile == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.instrumentationMappingOutputFile)
		})
		j2gort.Args4jRegisterHandler(p, opt_18_, j2gort.NewArgs4jBuiltinHandler(p, opt_18_, set_18_, "string", ""), false)
	}
	{
		opt_19_ := j2gort.NewArgs4jNamedOptionDef("--create_renaming_reports", nil, "If true, variable renaming and property renaming report files will be produced as {binary name}_vars_renaming_report.out and {binary name}_props_renaming_report.out. Note that this flag cannot be used in conjunction with either variable_renaming_report or property_renaming_report", "", false, false, true, false)
		set_19_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.createNameMapFiles = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.createNameMapFiles)
		})
		j2gort.Args4jRegisterHandler(p, opt_19_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_19_, set_19_), false)
	}
	{
		opt_20_ := j2gort.NewArgs4jNamedOptionDef("--source_map_include_content", nil, "Includes sources content into source map. Greatly increases the size of source maps but offers greater portability", "", false, false, false, false)
		set_20_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.sourceMapIncludeSourcesContent = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.sourceMapIncludeSourcesContent)
		})
		j2gort.Args4jRegisterHandler(p, opt_20_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_20_, set_20_), false)
	}
	{
		opt_21_ := j2gort.NewArgs4jNamedOptionDef("--property_renaming_report", nil, "File where the serialized version of the property renaming map produced should be saved", "", false, false, false, false)
		set_21_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.propertyMapOutputFile = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.propertyMapOutputFile == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.propertyMapOutputFile)
		})
		j2gort.Args4jRegisterHandler(p, opt_21_, j2gort.NewArgs4jBuiltinHandler(p, opt_21_, set_21_, "string", ""), false)
	}
	{
		opt_22_ := j2gort.NewArgs4jNamedOptionDef("--third_party", nil, "Check source validity but do not enforce Closure style rules and conventions", "", false, false, false, false)
		set_22_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.thirdParty = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.thirdParty)
		})
		j2gort.Args4jRegisterHandler(p, opt_22_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_22_, set_22_), false)
	}
	{
		opt_23_ := j2gort.NewArgs4jNamedOptionDef("--summary_detail_level", nil, "Controls how detailed the compilation summary is. Values: 0 (never print summary), 1 (print summary only if there are errors or warnings), 2 (print summary if the 'checkTypes' diagnostic  group is enabled, see --jscomp_warning), 3 (always print summary). The default level is 1", "", false, false, true, false)
		set_23_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("int", "int"), func(val_ any) {
			this_.summaryDetailLevel = j2gort.Args4jAsInt(p, val_)
		}, func() any {
			return j2gort.ValueOfInt(this_.summaryDetailLevel)
		})
		j2gort.Args4jRegisterHandler(p, opt_23_, j2gort.NewArgs4jBuiltinHandler(p, opt_23_, set_23_, "int", ""), false)
	}
	{
		opt_24_ := j2gort.NewArgs4jNamedOptionDef("--isolation_mode", nil, "If set to IIFE the compiler output will follow the form:\n  (function(){%output%}).call(this);\nOptions: NONE, IIFE", "", false, false, false, false)
		set_24_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("com.google.javascript.jscomp.CompilerOptions.IsolationMode", "IsolationMode"), func(val_ any) {
			this_.isolationMode = func() CompilerOptions_IsolationMode {
				_ = CompilerOptions_IsolationMode_Values()
				return j2gort.Cast[CompilerOptions_IsolationMode](j2gort.Args4jParseEnum(p, "com.google.javascript.jscomp.CompilerOptions.IsolationMode", val_))
			}()
		}, func() any {
			return this_.isolationMode
		})
		j2gort.Args4jRegisterHandler(p, opt_24_, func() j2gort.Ext_org_kohsuke_args4j_spi_OptionHandler {
			_ = CompilerOptions_IsolationMode_Values()
			return j2gort.NewArgs4jBuiltinHandler(p, opt_24_, set_24_, "enum", "com.google.javascript.jscomp.CompilerOptions.IsolationMode")
		}(), false)
	}
	{
		opt_25_ := j2gort.NewArgs4jNamedOptionDef("--output_wrapper", nil, "Interpolate output into this string at the place denoted by the marker token %output%. Use marker token %output|jsstring% to do js string escaping on the output. Consider using the --isolation_mode flag instead.", "", false, false, false, false)
		set_25_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.outputWrapper = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.outputWrapper == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.outputWrapper)
		})
		j2gort.Args4jRegisterHandler(p, opt_25_, j2gort.NewArgs4jBuiltinHandler(p, opt_25_, set_25_, "string", ""), false)
	}
	{
		opt_26_ := j2gort.NewArgs4jNamedOptionDef("--output_wrapper_file", nil, "Loads the specified file and passes the file contents to the --output_wrapper flag, replacing the value if it exists. This is useful if you want special characters like newline in the wrapper.", "", false, false, false, false)
		set_26_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.outputWrapperFile = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.outputWrapperFile == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.outputWrapperFile)
		})
		j2gort.Args4jRegisterHandler(p, opt_26_, j2gort.NewArgs4jBuiltinHandler(p, opt_26_, set_26_, "string", ""), false)
	}
	{
		opt_27_ := j2gort.NewArgs4jNamedOptionDef("--chunk_wrapper", nil, "An output wrapper for a JavaScript chunk (optional). The format is <name>:<wrapper>. The chunk name must correspond with a chunk specified using --chunk. The wrapper must contain %s as the code placeholder. Alternately, %output% can be used in place of %s. %n% can be used to represent a newline. The %basename% placeholder can also be used to substitute the base name of the chunk output file.", "", false, false, false, true)
		set_27_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.chunkWrapper == nil {
				this_.chunkWrapper = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.chunkWrapper, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_27_, j2gort.NewArgs4jBuiltinHandler(p, opt_27_, set_27_, "string", ""), false)
	}
	{
		opt_28_ := j2gort.NewArgs4jNamedOptionDef("--chunk_output_path_prefix", nil, "Prefix for filenames of compiled JS chunks. <chunk-name>.js will be appended to this prefix. Directories will be created as needed. Use with --chunk", "", false, false, false, false)
		set_28_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.chunkOutputPathPrefix = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.chunkOutputPathPrefix == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.chunkOutputPathPrefix)
		})
		j2gort.Args4jRegisterHandler(p, opt_28_, j2gort.NewArgs4jBuiltinHandler(p, opt_28_, set_28_, "string", ""), false)
	}
	{
		opt_29_ := j2gort.NewArgs4jNamedOptionDef("--create_source_map", nil, "If specified, a source map file mapping the generated source files back to the original source file will be output to the specified path. The %outname% placeholder will expand to the name of the output file that the source map corresponds to.", "", false, false, false, false)
		set_29_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.createSourceMap = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.createSourceMap == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.createSourceMap)
		})
		j2gort.Args4jRegisterHandler(p, opt_29_, j2gort.NewArgs4jBuiltinHandler(p, opt_29_, set_29_, "string", ""), false)
	}
	{
		opt_30_ := j2gort.NewArgs4jNamedOptionDef("--source_map_format", nil, "The source map format to produce. Options are V3 and DEFAULT, which are equivalent.", "", false, false, true, false)
		set_30_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("com.google.javascript.jscomp.SourceMap.Format", "Format"), func(val_ any) {
			this_.sourceMapFormat = func() *SourceMap_Format {
				_ = SourceMap_Format_Values()
				return j2gort.Cast[*SourceMap_Format](j2gort.Args4jParseEnum(p, "com.google.javascript.jscomp.SourceMap.Format", val_))
			}()
		}, func() any {
			return j2gort.ToAny(this_.sourceMapFormat)
		})
		j2gort.Args4jRegisterHandler(p, opt_30_, func() j2gort.Ext_org_kohsuke_args4j_spi_OptionHandler {
			_ = SourceMap_Format_Values()
			return j2gort.NewArgs4jBuiltinHandler(p, opt_30_, set_30_, "enum", "com.google.javascript.jscomp.SourceMap.Format")
		}(), false)
	}
	{
		opt_31_ := j2gort.NewArgs4jNamedOptionDef("--source_map_location_mapping", nil, "Source map location mapping separated by a '|' (i.e. filesystem-path|webserver-path)", "", false, false, false, true)
		set_31_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.sourceMapLocationMapping == nil {
				this_.sourceMapLocationMapping = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.sourceMapLocationMapping, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_31_, j2gort.NewArgs4jBuiltinHandler(p, opt_31_, set_31_, "string", ""), false)
	}
	{
		opt_32_ := j2gort.NewArgs4jNamedOptionDef("--source_map_input", nil, "Source map locations for input files, separated by a '|', (i.e. input-file-path|input-source-map)", "", false, false, false, true)
		set_32_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.sourceMapInputs == nil {
				this_.sourceMapInputs = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.sourceMapInputs, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_32_, j2gort.NewArgs4jBuiltinHandler(p, opt_32_, set_32_, "string", ""), false)
	}
	{
		opt_33_ := j2gort.NewArgs4jNamedOptionDef("--parse_inline_source_maps", nil, "Parse inline source maps (//# sourceMappingURL=data:...)", "", false, false, false, false)
		set_33_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.Boolean", "Boolean"), func(val_ any) {
			this_.parseInlineSourceMaps = j2gort.ValueOfBoolean(j2gort.Args4jAsBool(p, val_))
		}, func() any {
			return j2gort.ToAny(this_.parseInlineSourceMaps)
		})
		j2gort.Args4jRegisterHandler(p, opt_33_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_33_, set_33_), false)
	}
	{
		opt_34_ := j2gort.NewArgs4jNamedOptionDef("--apply_input_source_maps", nil, "Apply input source maps to the output source map, i.e. have the result map back to original inputs.  Input sourcemaps can be located in 2 ways:\n 1) by the//# sourceMappingURL=<url>. \n 2) using the--source_map_location_mapping flag.\nsourceMappingURL=<url> can read both paths and inline Base64 encoded sourcemaps. For inline Base64 encoded sourcemaps, see --parse_inline_source_maps.", "", false, false, false, false)
		set_34_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.applyInputSourceMaps = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.applyInputSourceMaps)
		})
		j2gort.Args4jRegisterHandler(p, opt_34_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_34_, set_34_), false)
	}
	{
		opt_35_ := j2gort.NewArgs4jNamedOptionDef("--jscomp_error", nil, "Make the named class of warnings an error. Must be one of the error group items. '*' adds all supported.", "", false, false, false, true)
		set_35_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.jscompError == nil {
				this_.jscompError = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.jscompError, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_35_, NewCommandLineRunner_Flags_WarningGuardErrorOptionHandler(p, opt_35_, set_35_), false)
	}
	{
		opt_36_ := j2gort.NewArgs4jNamedOptionDef("--jscomp_warning", nil, "Make the named class of warnings a normal warning. Must be one of the error group items. '*' adds all supported.", "", false, false, false, true)
		set_36_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.jscompWarning == nil {
				this_.jscompWarning = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.jscompWarning, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_36_, NewCommandLineRunner_Flags_WarningGuardWarningOptionHandler(p, opt_36_, set_36_), false)
	}
	{
		opt_37_ := j2gort.NewArgs4jNamedOptionDef("--jscomp_off", nil, "Turn off the named class of warnings. Must be one of the error group items. '*' adds all supported.", "", false, false, false, true)
		set_37_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.jscompOff == nil {
				this_.jscompOff = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.jscompOff, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_37_, NewCommandLineRunner_Flags_WarningGuardOffOptionHandler(p, opt_37_, set_37_), false)
	}
	{
		opt_38_ := j2gort.NewArgs4jNamedOptionDef("--define", []string{"--D", "-D"}, "Override the value of a variable annotated @define. The format is <name>[=<val>], where <name> is the name of a @define variable and <val> is a boolean, number, or a single-quoted string that contains no single quotes. If [=<val>] is omitted, the variable is marked true", "", false, false, false, true)
		set_38_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.define == nil {
				this_.define = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.define, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_38_, j2gort.NewArgs4jBuiltinHandler(p, opt_38_, set_38_, "string", ""), false)
	}
	{
		opt_39_ := j2gort.NewArgs4jNamedOptionDef("--charset", nil, "Input and output charset for all files. By default, we accept UTF-8 as input and output US_ASCII", "", false, false, false, false)
		set_39_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.charset = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.charset == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.charset)
		})
		j2gort.Args4jRegisterHandler(p, opt_39_, j2gort.NewArgs4jBuiltinHandler(p, opt_39_, set_39_, "string", ""), false)
	}
	{
		opt_40_ := j2gort.NewArgs4jNamedOptionDef("--compilation_level", []string{"-O"}, "Specifies the compilation level to use. Options: BUNDLE, WHITESPACE_ONLY, SIMPLE (default), TRANSPILE_ONLY,ADVANCED", "", false, false, false, false)
		set_40_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.compilationLevel = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.compilationLevel == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.compilationLevel)
		})
		j2gort.Args4jRegisterHandler(p, opt_40_, j2gort.NewArgs4jBuiltinHandler(p, opt_40_, set_40_, "string", ""), false)
	}
	{
		opt_41_ := j2gort.NewArgs4jNamedOptionDef("--num_parallel_threads", nil, "Use multiple threads to parallelize parts of the compilation.", "", false, false, true, false)
		set_41_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("int", "int"), func(val_ any) {
			this_.numParallelThreads = j2gort.Args4jAsInt(p, val_)
		}, func() any {
			return j2gort.ValueOfInt(this_.numParallelThreads)
		})
		j2gort.Args4jRegisterHandler(p, opt_41_, j2gort.NewArgs4jBuiltinHandler(p, opt_41_, set_41_, "int", ""), false)
	}
	{
		opt_42_ := j2gort.NewArgs4jNamedOptionDef("--checks_only", []string{"--checks-only"}, "Don't generate output. Run checks, but no optimization passes.", "", false, false, false, false)
		set_42_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.checksOnly = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.checksOnly)
		})
		j2gort.Args4jRegisterHandler(p, opt_42_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_42_, set_42_), false)
	}
	{
		opt_43_ := j2gort.NewArgs4jNamedOptionDef("--incremental_check_mode", nil, "Generate or check externs-like .i.js files representing individual libraries.", "", false, false, false, false)
		set_43_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("com.google.javascript.jscomp.CompilerOptions.IncrementalCheckMode", "IncrementalCheckMode"), func(val_ any) {
			this_.incrementalCheckMode = func() CompilerOptions_IncrementalCheckMode {
				_ = CompilerOptions_IncrementalCheckMode_Values()
				return j2gort.Cast[CompilerOptions_IncrementalCheckMode](j2gort.Args4jParseEnum(p, "com.google.javascript.jscomp.CompilerOptions.IncrementalCheckMode", val_))
			}()
		}, func() any {
			return this_.incrementalCheckMode
		})
		j2gort.Args4jRegisterHandler(p, opt_43_, func() j2gort.Ext_org_kohsuke_args4j_spi_OptionHandler {
			_ = CompilerOptions_IncrementalCheckMode_Values()
			return j2gort.NewArgs4jBuiltinHandler(p, opt_43_, set_43_, "enum", "com.google.javascript.jscomp.CompilerOptions.IncrementalCheckMode")
		}(), false)
	}
	{
		opt_44_ := j2gort.NewArgs4jNamedOptionDef("--continue_after_errors", nil, "Continue trying to compile after an error is encountered.", "", false, false, false, false)
		set_44_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.continueAfterErrors = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.continueAfterErrors)
		})
		j2gort.Args4jRegisterHandler(p, opt_44_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_44_, set_44_), false)
	}
	{
		opt_45_ := j2gort.NewArgs4jNamedOptionDef("--use_types_for_optimization", nil, "Enable or disable the optimizations based on available type information. Inaccurate type annotations may result in incorrect results.", "", false, false, false, false)
		set_45_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.useTypesForOptimization = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.useTypesForOptimization)
		})
		j2gort.Args4jRegisterHandler(p, opt_45_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_45_, set_45_), false)
	}
	{
		opt_46_ := j2gort.NewArgs4jNamedOptionDef("--assume_function_wrapper", nil, "Enable additional optimizations based on the assumption that the output will be wrapped with a function wrapper.  This flag is used to indicate that \"global\" declarations will not actually be global but instead isolated to the compilation unit. This enables additional optimizations.", "", false, false, false, false)
		set_46_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.assumeFunctionWrapper = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.assumeFunctionWrapper)
		})
		j2gort.Args4jRegisterHandler(p, opt_46_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_46_, set_46_), false)
	}
	{
		opt_47_ := j2gort.NewArgs4jNamedOptionDef("--warning_level", []string{"-W"}, "Specifies the warning level to use. Options: QUIET, DEFAULT, VERBOSE", "", false, false, false, false)
		set_47_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("com.google.javascript.jscomp.WarningLevel", "WarningLevel"), func(val_ any) {
			this_.warningLevel = func() WarningLevel {
				_ = WarningLevel_Values()
				return j2gort.Cast[WarningLevel](j2gort.Args4jParseEnum(p, "com.google.javascript.jscomp.WarningLevel", val_))
			}()
		}, func() any {
			return this_.warningLevel
		})
		j2gort.Args4jRegisterHandler(p, opt_47_, func() j2gort.Ext_org_kohsuke_args4j_spi_OptionHandler {
			_ = WarningLevel_Values()
			return j2gort.NewArgs4jBuiltinHandler(p, opt_47_, set_47_, "enum", "com.google.javascript.jscomp.WarningLevel")
		}(), false)
	}
	{
		opt_48_ := j2gort.NewArgs4jNamedOptionDef("--debug", nil, "Enable debugging options. Property renaming uses long mangled names which can be mapped back to the original name.", "", false, false, false, false)
		set_48_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.debug = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.debug)
		})
		j2gort.Args4jRegisterHandler(p, opt_48_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_48_, set_48_), false)
	}
	{
		opt_49_ := j2gort.NewArgs4jNamedOptionDef("--typed_ast_output_file", nil, "Sets file to output in-progress typedAST format. DO NOT USE!", "", false, false, true, false)
		set_49_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.typedAstOutputFile = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.typedAstOutputFile == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.typedAstOutputFile)
		})
		j2gort.Args4jRegisterHandler(p, opt_49_, j2gort.NewArgs4jBuiltinHandler(p, opt_49_, set_49_, "string", ""), false)
	}
	{
		opt_50_ := j2gort.NewArgs4jNamedOptionDef("--generate_exports", nil, "Generates export code for those marked with @export", "", false, false, false, false)
		set_50_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.generateExports = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.generateExports)
		})
		j2gort.Args4jRegisterHandler(p, opt_50_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_50_, set_50_), false)
	}
	{
		opt_51_ := j2gort.NewArgs4jNamedOptionDef("--export_local_property_definitions", nil, "Generates export code for local properties marked with @export", "", false, false, false, false)
		set_51_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.exportLocalPropertyDefinitions = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.exportLocalPropertyDefinitions)
		})
		j2gort.Args4jRegisterHandler(p, opt_51_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_51_, set_51_), false)
	}
	{
		opt_52_ := j2gort.NewArgs4jNamedOptionDef("--formatting", nil, "Specifies which formatting options, if any, should be applied to the output JS. Options: PRETTY_PRINT, PRINT_INPUT_DELIMITER, SINGLE_QUOTES", "", false, false, false, true)
		set_52_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("com.google.javascript.jscomp.CommandLineRunner.FormattingOption", "FormattingOption"), func(val_ any) {
			if this_.formatting == nil {
				this_.formatting = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.formatting, func() CommandLineRunner_FormattingOption {
				_ = CommandLineRunner_FormattingOption_Values()
				return j2gort.Cast[CommandLineRunner_FormattingOption](j2gort.Args4jParseEnum(p, "com.google.javascript.jscomp.CommandLineRunner.FormattingOption", val_))
			}())
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_52_, func() j2gort.Ext_org_kohsuke_args4j_spi_OptionHandler {
			_ = CommandLineRunner_FormattingOption_Values()
			return j2gort.NewArgs4jBuiltinHandler(p, opt_52_, set_52_, "enum", "com.google.javascript.jscomp.CommandLineRunner.FormattingOption")
		}(), false)
	}
	{
		opt_53_ := j2gort.NewArgs4jNamedOptionDef("--process_common_js_modules", nil, "Process CommonJS modules to a concatenable form.", "", false, false, false, false)
		set_53_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.processCommonJsModules = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.processCommonJsModules)
		})
		j2gort.Args4jRegisterHandler(p, opt_53_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_53_, set_53_), false)
	}
	{
		opt_54_ := j2gort.NewArgs4jNamedOptionDef("--js_module_root", nil, "Path prefixes to be removed from ES6 & CommonJS modules.", "", false, false, false, true)
		set_54_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.moduleRoot == nil {
				this_.moduleRoot = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.moduleRoot, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_54_, j2gort.NewArgs4jBuiltinHandler(p, opt_54_, set_54_, "string", ""), false)
	}
	{
		opt_55_ := j2gort.NewArgs4jNamedOptionDef("--process_closure_primitives", nil, "Processes built-ins from the Closure library, such as goog.require(), goog.provide(), and goog.exportSymbol(). True by default.", "", false, false, false, false)
		set_55_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.processClosurePrimitives = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.processClosurePrimitives)
		})
		j2gort.Args4jRegisterHandler(p, opt_55_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_55_, set_55_), false)
	}
	{
		opt_56_ := j2gort.NewArgs4jNamedOptionDef("--angular_pass", nil, "Generate $inject properties for AngularJS for functions annotated with @ngInject", "", false, false, false, false)
		set_56_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.angularPass = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.angularPass)
		})
		j2gort.Args4jRegisterHandler(p, opt_56_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_56_, set_56_), false)
	}
	{
		opt_57_ := j2gort.NewArgs4jNamedOptionDef("--polymer_version", nil, "Which version of Polymer is being used (1 or 2).", "", false, false, false, false)
		set_57_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.Integer", "Integer"), func(val_ any) {
			this_.polymerVersion = j2gort.ValueOfInt(j2gort.Args4jAsInt(p, val_))
		}, func() any {
			return j2gort.ToAny(this_.polymerVersion)
		})
		j2gort.Args4jRegisterHandler(p, opt_57_, j2gort.NewArgs4jBuiltinHandler(p, opt_57_, set_57_, "int", ""), false)
	}
	{
		opt_58_ := j2gort.NewArgs4jNamedOptionDef("--chrome_pass", nil, "Enable Chrome-specific options for handling cr.* functions.", "", false, false, true, false)
		set_58_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.chromePass = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.chromePass)
		})
		j2gort.Args4jRegisterHandler(p, opt_58_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_58_, set_58_), false)
	}
	{
		opt_59_ := j2gort.NewArgs4jNamedOptionDef("--j2cl_pass", nil, "Rewrite J2CL output to be compiler-friendly if enabled (ON or AUTO). Options:OFF, ON, AUTO(default)", "", false, false, true, false)
		set_59_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.j2clPassMode = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.j2clPassMode == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.j2clPassMode)
		})
		j2gort.Args4jRegisterHandler(p, opt_59_, j2gort.NewArgs4jBuiltinHandler(p, opt_59_, set_59_, "string", ""), false)
	}
	{
		opt_60_ := j2gort.NewArgs4jNamedOptionDef("--remove_j2cl_asserts", nil, "Remove calls to J2CL assertions.", "", false, false, true, false)
		set_60_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.removeJ2cLAsserts = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.removeJ2cLAsserts)
		})
		j2gort.Args4jRegisterHandler(p, opt_60_, j2gort.NewArgs4jBuiltinHandler(p, opt_60_, set_60_, "boolean", ""), false)
	}
	{
		opt_61_ := j2gort.NewArgs4jNamedOptionDef("--output_manifest", nil, "Prints out a list of all the files in the compilation. If --dependency_mode=PRUNE or PRUNE_LEGACY is specified, this will not include files that got dropped because they were not required. The %outname% placeholder expands to the JS output file. If you're using modularization, using %outname% will create a manifest for each module.", "", false, false, false, false)
		set_61_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.outputManifest = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.outputManifest == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.outputManifest)
		})
		j2gort.Args4jRegisterHandler(p, opt_61_, j2gort.NewArgs4jBuiltinHandler(p, opt_61_, set_61_, "string", ""), false)
	}
	{
		opt_62_ := j2gort.NewArgs4jNamedOptionDef("--output_chunk_dependencies", nil, "Prints out a JSON file of dependencies between chunks.", "", false, false, false, false)
		set_62_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.outputChunkDependencies = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.outputChunkDependencies == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.outputChunkDependencies)
		})
		j2gort.Args4jRegisterHandler(p, opt_62_, j2gort.NewArgs4jBuiltinHandler(p, opt_62_, set_62_, "string", ""), false)
	}
	{
		opt_63_ := j2gort.NewArgs4jNamedOptionDef("--language_in", nil, "Sets the language spec to which input sources should conform. Options: ECMASCRIPT3, ECMASCRIPT5, ECMASCRIPT5_STRICT, ECMASCRIPT_2015, ECMASCRIPT_2016, ECMASCRIPT_2017, ECMASCRIPT_2018, ECMASCRIPT_2019, ECMASCRIPT_2020,ECMASCRIPT_2021, STABLE, ECMASCRIPT_NEXT (latest features supported),UNSTABLE (for testing only)", "", false, false, false, false)
		set_63_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.languageIn = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.languageIn == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.languageIn)
		})
		j2gort.Args4jRegisterHandler(p, opt_63_, j2gort.NewArgs4jBuiltinHandler(p, opt_63_, set_63_, "string", ""), false)
	}
	{
		opt_64_ := j2gort.NewArgs4jNamedOptionDef("--language_out", nil, "Sets the language spec to which output should conform. Options: ECMASCRIPT3, ECMASCRIPT5, ECMASCRIPT_2015, ECMASCRIPT_2016, ECMASCRIPT_2017, ECMASCRIPT_2018, ECMASCRIPT_2019, ECMASCRIPT_2020, ECMASCRIPT_2021, STABLE, ECMASCRIPT_NEXT (latest features supported)", "", false, false, false, false)
		set_64_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.languageOut = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.languageOut == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.languageOut)
		})
		j2gort.Args4jRegisterHandler(p, opt_64_, j2gort.NewArgs4jBuiltinHandler(p, opt_64_, set_64_, "string", ""), false)
	}
	{
		opt_65_ := j2gort.NewArgs4jNamedOptionDef("--version", nil, "Prints the compiler version to stdout and exit.", "", false, false, false, false)
		set_65_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.version = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.version)
		})
		j2gort.Args4jRegisterHandler(p, opt_65_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_65_, set_65_), false)
	}
	{
		opt_66_ := j2gort.NewArgs4jNamedOptionDef("--translations_file", nil, "Source of translated messages. Currently only supports XTB.", "", false, false, true, false)
		set_66_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.translationsFile = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.translationsFile == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.translationsFile)
		})
		j2gort.Args4jRegisterHandler(p, opt_66_, j2gort.NewArgs4jBuiltinHandler(p, opt_66_, set_66_, "string", ""), false)
	}
	{
		opt_67_ := j2gort.NewArgs4jNamedOptionDef("--translations_project", nil, "Scopes all translations to the specified project.When specified, we will use different message ids so that messages in different projects can have different translations.", "", false, false, true, false)
		set_67_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.translationsProject = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.translationsProject == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.translationsProject)
		})
		j2gort.Args4jRegisterHandler(p, opt_67_, j2gort.NewArgs4jBuiltinHandler(p, opt_67_, set_67_, "string", ""), false)
	}
	{
		opt_68_ := j2gort.NewArgs4jNamedOptionDef("--flagfile", nil, "A file (or files) containing additional command-line options.", "", false, false, true, true)
		set_68_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.flagFiles == nil {
				this_.flagFiles = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.flagFiles, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_68_, j2gort.NewArgs4jBuiltinHandler(p, opt_68_, set_68_, "string", ""), false)
	}
	{
		opt_69_ := j2gort.NewArgs4jNamedOptionDef("--warnings_allowlist_file", []string{"--warnings_whitelist_file"}, "A file containing warnings to suppress. Each line should be of the form\n<file-name>:<line-number>?  <warning-description>", "", false, false, false, false)
		set_69_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.warningsAllowlistFile = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.warningsAllowlistFile == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.warningsAllowlistFile)
		})
		j2gort.Args4jRegisterHandler(p, opt_69_, j2gort.NewArgs4jBuiltinHandler(p, opt_69_, set_69_, "string", ""), false)
	}
	{
		opt_70_ := j2gort.NewArgs4jNamedOptionDef("--hide_warnings_for", nil, "If specified, files whose path contains this string will have their warnings hidden. You may specify multiple.", "", false, false, false, true)
		set_70_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.hideWarningsFor == nil {
				this_.hideWarningsFor = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.hideWarningsFor, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_70_, j2gort.NewArgs4jBuiltinHandler(p, opt_70_, set_70_, "string", ""), false)
	}
	{
		opt_71_ := j2gort.NewArgs4jNamedOptionDef("--extra_annotation_name", nil, "A allowlist of tag names in JSDoc. You may specify multiple", "", false, false, false, true)
		set_71_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.extraAnnotationName == nil {
				this_.extraAnnotationName = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.extraAnnotationName, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_71_, j2gort.NewArgs4jBuiltinHandler(p, opt_71_, set_71_, "string", ""), false)
	}
	{
		opt_72_ := j2gort.NewArgs4jNamedOptionDef("--tracer_mode", nil, "Shows the duration of each compiler pass and the impact to the compiled output size. Options: ALL, AST_SIZE, RAW_SIZE, TIMING_ONLY, OFF", "", false, false, true, false)
		set_72_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("com.google.javascript.jscomp.CompilerOptions.TracerMode", "TracerMode"), func(val_ any) {
			this_.tracerMode = func() CompilerOptions_TracerMode {
				_ = CompilerOptions_TracerMode_Values()
				return j2gort.Cast[CompilerOptions_TracerMode](j2gort.Args4jParseEnum(p, "com.google.javascript.jscomp.CompilerOptions.TracerMode", val_))
			}()
		}, func() any {
			return this_.tracerMode
		})
		j2gort.Args4jRegisterHandler(p, opt_72_, func() j2gort.Ext_org_kohsuke_args4j_spi_OptionHandler {
			_ = CompilerOptions_TracerMode_Values()
			return j2gort.NewArgs4jBuiltinHandler(p, opt_72_, set_72_, "enum", "com.google.javascript.jscomp.CompilerOptions.TracerMode")
		}(), false)
	}
	{
		opt_73_ := j2gort.NewArgs4jNamedOptionDef("--rename_variable_prefix", nil, "Specifies a prefix that will be prepended to all variables.", "", false, false, false, false)
		set_73_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.renamePrefix = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.renamePrefix == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.renamePrefix)
		})
		j2gort.Args4jRegisterHandler(p, opt_73_, j2gort.NewArgs4jBuiltinHandler(p, opt_73_, set_73_, "string", ""), false)
	}
	{
		opt_74_ := j2gort.NewArgs4jNamedOptionDef("--rename_prefix_namespace", nil, "Specifies the name of an object that will be used to store all non-extern globals", "", false, false, false, false)
		set_74_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.renamePrefixNamespace = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.renamePrefixNamespace == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.renamePrefixNamespace)
		})
		j2gort.Args4jRegisterHandler(p, opt_74_, j2gort.NewArgs4jBuiltinHandler(p, opt_74_, set_74_, "string", ""), false)
	}
	{
		opt_75_ := j2gort.NewArgs4jNamedOptionDef("--conformance_configs", nil, "A list of JS Conformance configurations in text protocol buffer format.", "", false, false, false, true)
		set_75_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.conformanceConfigs == nil {
				this_.conformanceConfigs = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.conformanceConfigs, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_75_, j2gort.NewArgs4jBuiltinHandler(p, opt_75_, set_75_, "string", ""), false)
	}
	{
		opt_76_ := j2gort.NewArgs4jNamedOptionDef("--env", nil, "Determines the set of builtin externs to load. Options: BROWSER, CUSTOM. Defaults to BROWSER.", "", false, false, false, false)
		set_76_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("com.google.javascript.jscomp.CompilerOptions.Environment", "Environment"), func(val_ any) {
			this_.environment = func() CompilerOptions_Environment {
				_ = CompilerOptions_Environment_Values()
				return j2gort.Cast[CompilerOptions_Environment](j2gort.Args4jParseEnum(p, "com.google.javascript.jscomp.CompilerOptions.Environment", val_))
			}()
		}, func() any {
			return this_.environment
		})
		j2gort.Args4jRegisterHandler(p, opt_76_, func() j2gort.Ext_org_kohsuke_args4j_spi_OptionHandler {
			_ = CompilerOptions_Environment_Values()
			return j2gort.NewArgs4jBuiltinHandler(p, opt_76_, set_76_, "enum", "com.google.javascript.jscomp.CompilerOptions.Environment")
		}(), false)
	}
	{
		opt_77_ := j2gort.NewArgs4jNamedOptionDef("--json_streams", nil, "Specifies whether standard input and output streams will be a JSON array of sources. Each source will be an object of the form {path: filename, src: file_contents, source_map: srcmap_contents }. Intended for use by stream-based build systems such as gulpjs. Options: NONE, IN, OUT, BOTH. Defaults to NONE.", "", false, false, false, false)
		set_77_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("com.google.javascript.jscomp.CompilerOptions.JsonStreamMode", "JsonStreamMode"), func(val_ any) {
			this_.jsonStreamMode = func() CompilerOptions_JsonStreamMode {
				_ = CompilerOptions_JsonStreamMode_Values()
				return j2gort.Cast[CompilerOptions_JsonStreamMode](j2gort.Args4jParseEnum(p, "com.google.javascript.jscomp.CompilerOptions.JsonStreamMode", val_))
			}()
		}, func() any {
			return this_.jsonStreamMode
		})
		j2gort.Args4jRegisterHandler(p, opt_77_, func() j2gort.Ext_org_kohsuke_args4j_spi_OptionHandler {
			_ = CompilerOptions_JsonStreamMode_Values()
			return j2gort.NewArgs4jBuiltinHandler(p, opt_77_, set_77_, "enum", "com.google.javascript.jscomp.CompilerOptions.JsonStreamMode")
		}(), false)
	}
	{
		opt_78_ := j2gort.NewArgs4jNamedOptionDef("--preserve_type_annotations", nil, "Preserves type annotations.", "", false, false, true, false)
		set_78_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.preserveTypeAnnotations = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.preserveTypeAnnotations)
		})
		j2gort.Args4jRegisterHandler(p, opt_78_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_78_, set_78_), false)
	}
	{
		opt_79_ := j2gort.NewArgs4jNamedOptionDef("--inject_libraries", nil, "Allow injecting runtime libraries.", "", false, false, false, false)
		set_79_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.injectLibraries = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.injectLibraries)
		})
		j2gort.Args4jRegisterHandler(p, opt_79_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_79_, set_79_), false)
	}
	{
		opt_80_ := j2gort.NewArgs4jNamedOptionDef("--force_inject_library", nil, "Force injection of named runtime libraries. The format is <name> where <name> is the name of a runtime library. Possible libraries include: base, es6_runtime", "", false, false, false, true)
		set_80_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.forceInjectLibraries == nil {
				this_.forceInjectLibraries = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.forceInjectLibraries, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_80_, j2gort.NewArgs4jBuiltinHandler(p, opt_80_, set_80_, "string", ""), false)
	}
	{
		opt_81_ := j2gort.NewArgs4jNamedOptionDef("--dependency_mode", nil, "Specifies how the compiler should determine the set and order of files for a compilation. Options: NONE the compiler will include all src files in the order listed, SORT_ONLY the compiler will include all source files in dependency order, PRUNE files will only be included if they are transitive dependencies of files listed in the --entry_point flag and then sorted in dependency order, PRUNE_LEGACY same as PRUNE but files that do not goog.provide a namespace and are not modules will be automatically added as --entry_point entries. Defaults to PRUNE_LEGACY if entry points are defined, otherwise to NONE.", "", false, false, false, false)
		set_81_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("com.google.javascript.jscomp.DependencyOptions.DependencyMode", "DependencyMode"), func(val_ any) {
			this_.dependencyMode = func() DependencyOptions_DependencyMode {
				_ = DependencyOptions_DependencyMode_Values()
				return j2gort.Cast[DependencyOptions_DependencyMode](j2gort.Args4jParseEnum(p, "com.google.javascript.jscomp.DependencyOptions.DependencyMode", val_))
			}()
		}, func() any {
			return this_.dependencyMode
		})
		j2gort.Args4jRegisterHandler(p, opt_81_, func() j2gort.Ext_org_kohsuke_args4j_spi_OptionHandler {
			_ = DependencyOptions_DependencyMode_Values()
			return j2gort.NewArgs4jBuiltinHandler(p, opt_81_, set_81_, "enum", "com.google.javascript.jscomp.DependencyOptions.DependencyMode")
		}(), false)
	}
	{
		opt_82_ := j2gort.NewArgs4jNamedOptionDef("--entry_point", nil, "A file or namespace to use as the starting point for determining which src files to include in the compilation. ES6 and CommonJS modules are specified as file paths (without the extension). Closure-library namespaces are specified with a \"goog:\" prefix. Example: --entry_point=goog:goog.Promise", "", false, false, false, true)
		set_82_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.entryPoint == nil {
				this_.entryPoint = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.entryPoint, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_82_, j2gort.NewArgs4jBuiltinHandler(p, opt_82_, set_82_, "string", ""), false)
	}
	{
		opt_83_ := j2gort.NewArgs4jNamedOptionDef("--rewrite_polyfills", nil, "Injects polyfills for ES2015+ library classes and methods used in source. See also the \"Polyfills\" GitHub Wiki page.", "", false, false, false, false)
		set_83_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.rewritePolyfills = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.rewritePolyfills)
		})
		j2gort.Args4jRegisterHandler(p, opt_83_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_83_, set_83_), false)
	}
	{
		opt_84_ := j2gort.NewArgs4jNamedOptionDef("--isolate_polyfills", nil, "Hides injected polyfills from the global scope and any external code. See the the \"Polyfills\" GitHub Wiki page for details.", "", false, false, false, false)
		set_84_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.isolatePolyfills = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.isolatePolyfills)
		})
		j2gort.Args4jRegisterHandler(p, opt_84_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_84_, set_84_), false)
	}
	{
		opt_85_ := j2gort.NewArgs4jNamedOptionDef("--print_source_after_each_pass", nil, "Whether to iteratively print resulting JS source per pass.", "", false, false, true, false)
		set_85_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.printSourceAfterEachPass = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.printSourceAfterEachPass)
		})
		j2gort.Args4jRegisterHandler(p, opt_85_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_85_, set_85_), false)
	}
	{
		opt_86_ := j2gort.NewArgs4jNamedOptionDef("--module_resolution", nil, "Specifies how the compiler locates modules. BROWSER requires all module imports to begin with a '.' or '/' and have a file extension. NODE uses the node module rules. WEBPACK looks up modules from a special lookup map.", "", false, false, false, false)
		set_86_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode", "ResolutionMode"), func(val_ any) {
			this_.moduleResolutionMode = func() ModuleLoader_ResolutionMode {
				_ = ModuleLoader_ResolutionMode_Values()
				return j2gort.Cast[ModuleLoader_ResolutionMode](j2gort.Args4jParseEnum(p, "com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode", val_))
			}()
		}, func() any {
			return this_.moduleResolutionMode
		})
		j2gort.Args4jRegisterHandler(p, opt_86_, func() j2gort.Ext_org_kohsuke_args4j_spi_OptionHandler {
			_ = ModuleLoader_ResolutionMode_Values()
			return j2gort.NewArgs4jBuiltinHandler(p, opt_86_, set_86_, "enum", "com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode")
		}(), false)
	}
	{
		opt_87_ := j2gort.NewArgs4jNamedOptionDef("--browser_resolver_prefix_replacements", nil, "Prefixes to replace in ES6 import paths before resolving. module_resolution must be BROWSER_WITH_TRANSFORMED_PREFIXES to take effect.", "", false, false, false, true)
		set_87_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.util.Map", "Map"), func(val_ any) {
			if this_.browserResolverPrefixReplacements == nil {
				this_.browserResolverPrefixReplacements = j2gort.NewArgs4jStringMap()
			}
			k_, v_ := j2gort.Args4jSplitMapEntry(p, val_)
			j2gort.Args4jMapPut(this_.browserResolverPrefixReplacements, k_, v_)
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_87_, j2gort.NewArgs4jBuiltinHandler(p, opt_87_, set_87_, "map", ""), false)
	}
	{
		opt_88_ := j2gort.NewArgs4jNamedOptionDef("--package_json_entry_names", nil, "Ordered list of entries to look for in package.json files when processing modules with the NODE module resolution strategy (i.e. esnext:main,browser,main). Defaults to a list with the following entries: \"browser\", \"module\", \"main\".", "", false, false, false, false)
		set_88_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.packageJsonEntryNames = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.packageJsonEntryNames == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.packageJsonEntryNames)
		})
		j2gort.Args4jRegisterHandler(p, opt_88_, j2gort.NewArgs4jBuiltinHandler(p, opt_88_, set_88_, "string", ""), false)
	}
	{
		opt_89_ := j2gort.NewArgs4jNamedOptionDef("--error_format", nil, "Specifies format for error messages.", "", false, false, false, false)
		set_89_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("com.google.javascript.jscomp.AbstractCommandLineRunner.CommandLineConfig.ErrorFormatOption", "ErrorFormatOption"), func(val_ any) {
			this_.errorFormat = func() AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption {
				_ = AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_Values()
				return j2gort.Cast[AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption](j2gort.Args4jParseEnum(p, "com.google.javascript.jscomp.AbstractCommandLineRunner.CommandLineConfig.ErrorFormatOption", val_))
			}()
		}, func() any {
			return this_.errorFormat
		})
		j2gort.Args4jRegisterHandler(p, opt_89_, func() j2gort.Ext_org_kohsuke_args4j_spi_OptionHandler {
			_ = AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_Values()
			return j2gort.NewArgs4jBuiltinHandler(p, opt_89_, set_89_, "enum", "com.google.javascript.jscomp.AbstractCommandLineRunner.CommandLineConfig.ErrorFormatOption")
		}(), false)
	}
	{
		opt_90_ := j2gort.NewArgs4jNamedOptionDef("--renaming", nil, "Disables variable renaming. Cannot be used with ADVANCED optimizations.", "", false, false, false, false)
		set_90_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.renaming = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.renaming)
		})
		j2gort.Args4jRegisterHandler(p, opt_90_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_90_, set_90_), false)
	}
	{
		opt_91_ := j2gort.NewArgs4jNamedOptionDef("--help_markdown", nil, "Prints markdown formatted flag usage", "", false, false, true, false)
		set_91_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.helpMarkdown = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.helpMarkdown)
		})
		j2gort.Args4jRegisterHandler(p, opt_91_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_91_, set_91_), false)
	}
	{
		opt_92_ := j2gort.NewArgs4jNamedOptionDef("--instrument_for_coverage_option", nil, "Enable code instrumentation to perform code coverage analysis. Options are:\n 1. NONE (default)\n 2. LINE - Instrument code by line.\n 3. BRANCH - Instrument code by branch.\n 4. PRODUCTION - Function Instrumentation on compiled JS code.\n", "", false, false, false, false)
		set_92_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.instrumentForCoverageOption = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.instrumentForCoverageOption == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.instrumentForCoverageOption)
		})
		j2gort.Args4jRegisterHandler(p, opt_92_, j2gort.NewArgs4jBuiltinHandler(p, opt_92_, set_92_, "string", ""), false)
	}
	{
		opt_93_ := j2gort.NewArgs4jNamedOptionDef("--production_instrumentation_array_name", nil, "Name of the global array used by production instrumentation. The array name should be declared as an extern so it is not renamed by the compiler. A function that parses the global array should also be included. This flag is to be used in tandem with --instrument_code=PRODUCTION", "", false, false, false, false)
		set_93_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.productionInstrumentationArrayName = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.productionInstrumentationArrayName == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.productionInstrumentationArrayName)
		})
		j2gort.Args4jRegisterHandler(p, opt_93_, j2gort.NewArgs4jBuiltinHandler(p, opt_93_, set_93_, "string", ""), false)
	}
	{
		opt_94_ := j2gort.NewArgs4jNamedOptionDef("--chunk_output_type", nil, "Indicates what format the compiler should use for output chunks. GLOBAL_NAMESPACE is typically used in conjunction with --rename_prefix_namespace. ES_MODULES outputs chunks as proper modules with 'import' and 'export' statements.", "", false, false, false, false)
		set_94_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("com.google.javascript.jscomp.CompilerOptions.ChunkOutputType", "ChunkOutputType"), func(val_ any) {
			this_.chunkOutputType = func() CompilerOptions_ChunkOutputType {
				_ = CompilerOptions_ChunkOutputType_Values()
				return j2gort.Cast[CompilerOptions_ChunkOutputType](j2gort.Args4jParseEnum(p, "com.google.javascript.jscomp.CompilerOptions.ChunkOutputType", val_))
			}()
		}, func() any {
			return this_.chunkOutputType
		})
		j2gort.Args4jRegisterHandler(p, opt_94_, func() j2gort.Ext_org_kohsuke_args4j_spi_OptionHandler {
			_ = CompilerOptions_ChunkOutputType_Values()
			return j2gort.NewArgs4jBuiltinHandler(p, opt_94_, set_94_, "enum", "com.google.javascript.jscomp.CompilerOptions.ChunkOutputType")
		}(), false)
	}
	{
		opt_95_ := j2gort.NewArgs4jNamedOptionDef("--allow_dynamic_import", nil, "Indicates that the compiler should allow dynamic import expressions. Dynamic import expressions are not yet fully supported and may lead to broken output code.", "", false, false, false, false)
		set_95_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.allowDynamicImport = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.allowDynamicImport)
		})
		j2gort.Args4jRegisterHandler(p, opt_95_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_95_, set_95_), false)
	}
	{
		opt_96_ := j2gort.NewArgs4jNamedOptionDef("--dynamic_import_alias", nil, "Instructs the compiler to replace dynamic imports expressions with a function call using the specified name. Allows dynamic import expressions to be externally polyfilled when the output language level does not natively support them. An alias of 'import' is allowed.", "", false, false, false, false)
		set_96_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.dynamicImportAlias = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.dynamicImportAlias == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.dynamicImportAlias)
		})
		j2gort.Args4jRegisterHandler(p, opt_96_, j2gort.NewArgs4jBuiltinHandler(p, opt_96_, set_96_, "string", ""), false)
	}
	{
		opt_97_ := j2gort.NewArgs4jNamedOptionDef("--assume_static_inheritance_is_not_used", nil, "Assume that static (class-side) inheritance is not being used and that static methods will not be referenced via `this` or through subclasses. This enables optimizations that could break code that did those things.", "", false, false, false, false)
		set_97_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.assumeStaticInheritanceIsNotUsed = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.assumeStaticInheritanceIsNotUsed)
		})
		j2gort.Args4jRegisterHandler(p, opt_97_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_97_, set_97_), false)
	}
	{
		opt_98_ := j2gort.NewArgs4jNamedOptionDef("--assume_no_prototype_method_enumeration", nil, "Assume that prototype method enumeration is not being used. This allows the compiler to move a prototype method declaration into a deeper chunk without creating stub functions in a parent chunk.", "", false, false, false, false)
		set_98_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("boolean", "boolean"), func(val_ any) {
			this_.assumeNoPrototypeMethodEnumeration = j2gort.Args4jAsBool(p, val_)
		}, func() any {
			return j2gort.ValueOfBoolean(this_.assumeNoPrototypeMethodEnumeration)
		})
		j2gort.Args4jRegisterHandler(p, opt_98_, NewCommandLineRunner_Flags_BooleanOptionHandler(p, opt_98_, set_98_), false)
	}
	{
		opt_99_ := j2gort.NewArgs4jNamedOptionDef("--variable_map_input_file", nil, "File containing the serialized version of the variable renaming map produced by a previous compilation", "", false, false, false, false)
		set_99_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.variableMapInputFile = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.variableMapInputFile == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.variableMapInputFile)
		})
		j2gort.Args4jRegisterHandler(p, opt_99_, j2gort.NewArgs4jBuiltinHandler(p, opt_99_, set_99_, "string", ""), false)
	}
	{
		opt_100_ := j2gort.NewArgs4jNamedOptionDef("--property_map_input_file", nil, "File containing the serialized version of the property renaming map produced by a previous compilation", "", false, false, false, false)
		set_100_ := j2gort.NewArgs4jFieldSetter(false, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			this_.propertyMapInputFile = j2gort.Args4jAsString(val_)
		}, func() any {
			if this_.propertyMapInputFile == j2gort.NullString {
				return nil
			}
			return j2gort.BoxString(this_.propertyMapInputFile)
		})
		j2gort.Args4jRegisterHandler(p, opt_100_, j2gort.NewArgs4jBuiltinHandler(p, opt_100_, set_100_, "string", ""), false)
	}
	{
		opt_101_ := j2gort.NewArgs4jNamedOptionDef("--expected_diagnostics", nil, "Expected diagnostics in the format [(line,col)]CODE:regex. CODE is the JSCompiler DiagnosticType key. regex matches the description.", "", false, false, false, true)
		set_101_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.expectedDiagnostics == nil {
				this_.expectedDiagnostics = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.expectedDiagnostics, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_101_, j2gort.NewArgs4jBuiltinHandler(p, opt_101_, set_101_, "string", ""), false)
	}
	{
		opt_102_ := j2gort.NewArgs4jArgumentDef("", "", false, false, true, 0)
		set_102_ := j2gort.NewArgs4jFieldSetter(true, j2gort.ClassLiteral("java.lang.String", "String"), func(val_ any) {
			if this_.arguments == nil {
				this_.arguments = j2gort.NewArgs4jList()
			}
			j2gort.Args4jListAdd(this_.arguments, j2gort.BoxString(j2gort.Args4jAsString(val_)))
		}, func() any {
			return nil
		})
		j2gort.Args4jRegisterHandler(p, opt_102_, j2gort.NewArgs4jBuiltinHandler(p, opt_102_, set_102_, "string", ""), true)
	}
}

// Java class: com.google.javascript.jscomp.CommandLineRunner.Flags.BooleanOptionHandler
type CommandLineRunner_Flags_BooleanOptionHandler struct {
	j2gort.Ext_org_kohsuke_args4j_spi_OptionHandlerBase
	idHash_ atomic.Int32
}

func (this_ *CommandLineRunner_Flags_BooleanOptionHandler) BindSelf_(s any) {
	_ = s
	this_.Ext_org_kohsuke_args4j_spi_OptionHandlerBase.BindExtSelf_(s)
}
func (this_ *CommandLineRunner_Flags_BooleanOptionHandler) IsJavaNil_() bool { return this_ == nil }
func (*CommandLineRunner_Flags_BooleanOptionHandler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CommandLineRunner$Flags$BooleanOptionHandler", "BooleanOptionHandler"
}
func (this_ *CommandLineRunner_Flags_BooleanOptionHandler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CommandLineRunner.Flags.BooleanOptionHandler", "com.google.javascript.jscomp.CommandLineRunner$Flags$BooleanOptionHandler", func(v_ any) bool {
		return j2gort.InstanceOf[*CommandLineRunner_Flags_BooleanOptionHandler](v_)
	})
}
func (*CommandLineRunner_Flags_BooleanOptionHandler) Implements_j2gort_Ext_org_kohsuke_args4j_spi_OptionHandler() {
}

var static_CommandLineRunner_Flags_BooleanOptionHandler_TRUES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_CommandLineRunner_Flags_BooleanOptionHandler_TRUES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_CommandLineRunner_Flags_BooleanOptionHandler()
	return static_CommandLineRunner_Flags_BooleanOptionHandler_TRUES
}
func StaticSet_CommandLineRunner_Flags_BooleanOptionHandler_TRUES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_CommandLineRunner_Flags_BooleanOptionHandler()
	static_CommandLineRunner_Flags_BooleanOptionHandler_TRUES = v
}

var static_CommandLineRunner_Flags_BooleanOptionHandler_FALSES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_CommandLineRunner_Flags_BooleanOptionHandler_FALSES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_CommandLineRunner_Flags_BooleanOptionHandler()
	return static_CommandLineRunner_Flags_BooleanOptionHandler_FALSES
}
func StaticSet_CommandLineRunner_Flags_BooleanOptionHandler_FALSES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_CommandLineRunner_Flags_BooleanOptionHandler()
	static_CommandLineRunner_Flags_BooleanOptionHandler_FALSES = v
}

var clinitDone_CommandLineRunner_Flags_BooleanOptionHandler uint32
var clinitState_CommandLineRunner_Flags_BooleanOptionHandler j2gort.ClinitState

func clinit_CommandLineRunner_Flags_BooleanOptionHandler() {
	if atomic.LoadUint32(&clinitDone_CommandLineRunner_Flags_BooleanOptionHandler) != 0 {
		return
	}
	clinitSlow_CommandLineRunner_Flags_BooleanOptionHandler()
}

func clinitSlow_CommandLineRunner_Flags_BooleanOptionHandler() {
	if !j2gort.EnterClinit(&clinitState_CommandLineRunner_Flags_BooleanOptionHandler, &clinitDone_CommandLineRunner_Flags_BooleanOptionHandler) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CommandLineRunner_Flags_BooleanOptionHandler, &clinitDone_CommandLineRunner_Flags_BooleanOptionHandler)
	static_CommandLineRunner_Flags_BooleanOptionHandler_TRUES = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_4("true", "on", "yes", "1")
	static_CommandLineRunner_Flags_BooleanOptionHandler_FALSES = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_4("false", "off", "no", "0")
}

func (this_ *CommandLineRunner_Flags_BooleanOptionHandler) initFields_CommandLineRunner_Flags_BooleanOptionHandler() {
	_ = this_
}

// Java constructor: BooleanOptionHandler(org.kohsuke.args4j.CmdLineParser parser, org.kohsuke.args4j.OptionDef option, org.kohsuke.args4j.spi.Setter<? super java.lang.Boolean> setter)
func NewCommandLineRunner_Flags_BooleanOptionHandler(parser *j2gort.Ext_org_kohsuke_args4j_CmdLineParser, option j2gort.Ext_org_kohsuke_args4j_OptionDef, setter j2gort.Ext_org_kohsuke_args4j_spi_Setter) *CommandLineRunner_Flags_BooleanOptionHandler {
	this_ := &CommandLineRunner_Flags_BooleanOptionHandler{}
	this_.BindSelf_(this_)
	InitCommandLineRunner_Flags_BooleanOptionHandler(this_, parser, option, setter)
	return this_
}

func InitCommandLineRunner_Flags_BooleanOptionHandler(this_ *CommandLineRunner_Flags_BooleanOptionHandler, parser *j2gort.Ext_org_kohsuke_args4j_CmdLineParser, option j2gort.Ext_org_kohsuke_args4j_OptionDef, setter j2gort.Ext_org_kohsuke_args4j_spi_Setter) {
	_ = parser
	_ = option
	_ = setter
	j2gort.InitExt_org_kohsuke_args4j_spi_OptionHandler(&this_.Ext_org_kohsuke_args4j_spi_OptionHandlerBase, parser, option, setter)
	this_.initFields_CommandLineRunner_Flags_BooleanOptionHandler()
}

// Java method: BooleanOptionHandler.parseArguments(org.kohsuke.args4j.spi.Parameters params)
func (this_ *CommandLineRunner_Flags_BooleanOptionHandler) ParseArguments(params j2gort.Ext_org_kohsuke_args4j_spi_Parameters) int32 {
	_ = params
	var param string = j2gort.NullString
	_ = param
	func() {
		defer func() {
			if rec_51 := recover(); rec_51 != nil {
				jex_52 := javalang.RecoverException(rec_51)
				_ = jex_52
				if j2gort.InstanceOf[j2gort.Ext_org_kohsuke_args4j_CmdLineException](jex_52) {
					e := j2gort.AsIface[j2gort.Ext_org_kohsuke_args4j_CmdLineException](jex_52)
					_ = e
					param = j2gort.NullString
				} else {
					panic(rec_51)
				}
			}
		}()
		param = params.GetParameter(int32(0))
	}()
	if param == j2gort.NullString {
		this_.GetExt_org_kohsuke_args4j_spi_OptionHandlerBase().F_Setter.AddValue(j2gort.ValueOfBoolean(true))
		return int32(0)
	} else {
		var lowerParam string = j2gort.Ext_com_google_common_base_Ascii_ToLowerCase_1_string(param)
		_ = lowerParam
		if j2gort.SetContainsStr(StaticGet_CommandLineRunner_Flags_BooleanOptionHandler_TRUES(), lowerParam) {
			this_.GetExt_org_kohsuke_args4j_spi_OptionHandlerBase().F_Setter.AddValue(j2gort.ValueOfBoolean(true))
		} else if j2gort.SetContainsStr(StaticGet_CommandLineRunner_Flags_BooleanOptionHandler_FALSES(), lowerParam) {
			this_.GetExt_org_kohsuke_args4j_spi_OptionHandlerBase().F_Setter.AddValue(j2gort.ValueOfBoolean(false))
		} else {
			this_.GetExt_org_kohsuke_args4j_spi_OptionHandlerBase().F_Setter.AddValue(j2gort.ValueOfBoolean(true))
			return int32(0)
		}
		return int32(1)
	}
}

// Java method: BooleanOptionHandler.getDefaultMetaVariable()
func (this_ *CommandLineRunner_Flags_BooleanOptionHandler) GetDefaultMetaVariable() string {
	return j2gort.NullString
}

// Java class: com.google.javascript.jscomp.CommandLineRunner.Flags.WarningGuardErrorOptionHandler
type CommandLineRunner_Flags_WarningGuardErrorOptionHandler struct {
	j2gort.Ext_org_kohsuke_args4j_spi_StringOptionHandlerBase
	idHash_ atomic.Int32
}

func (this_ *CommandLineRunner_Flags_WarningGuardErrorOptionHandler) BindSelf_(s any) {
	_ = s
	this_.Ext_org_kohsuke_args4j_spi_StringOptionHandlerBase.BindExtSelf_(s)
}
func (this_ *CommandLineRunner_Flags_WarningGuardErrorOptionHandler) IsJavaNil_() bool {
	return this_ == nil
}
func (*CommandLineRunner_Flags_WarningGuardErrorOptionHandler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CommandLineRunner$Flags$WarningGuardErrorOptionHandler", "WarningGuardErrorOptionHandler"
}
func (this_ *CommandLineRunner_Flags_WarningGuardErrorOptionHandler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CommandLineRunner.Flags.WarningGuardErrorOptionHandler", "com.google.javascript.jscomp.CommandLineRunner$Flags$WarningGuardErrorOptionHandler", func(v_ any) bool {
		return j2gort.InstanceOf[*CommandLineRunner_Flags_WarningGuardErrorOptionHandler](v_)
	})
}
func (*CommandLineRunner_Flags_WarningGuardErrorOptionHandler) Implements_j2gort_Ext_org_kohsuke_args4j_spi_StringOptionHandler() {
}
func (*CommandLineRunner_Flags_WarningGuardErrorOptionHandler) Implements_j2gort_Ext_org_kohsuke_args4j_spi_OptionHandler() {
}

func (this_ *CommandLineRunner_Flags_WarningGuardErrorOptionHandler) initFields_CommandLineRunner_Flags_WarningGuardErrorOptionHandler() {
	_ = this_
}

// Java constructor: WarningGuardErrorOptionHandler(org.kohsuke.args4j.CmdLineParser parser, org.kohsuke.args4j.OptionDef option, org.kohsuke.args4j.spi.Setter<? super java.lang.String> setter)
func NewCommandLineRunner_Flags_WarningGuardErrorOptionHandler(parser *j2gort.Ext_org_kohsuke_args4j_CmdLineParser, option j2gort.Ext_org_kohsuke_args4j_OptionDef, setter j2gort.Ext_org_kohsuke_args4j_spi_Setter) *CommandLineRunner_Flags_WarningGuardErrorOptionHandler {
	this_ := &CommandLineRunner_Flags_WarningGuardErrorOptionHandler{}
	this_.BindSelf_(this_)
	InitCommandLineRunner_Flags_WarningGuardErrorOptionHandler(this_, parser, option, setter)
	return this_
}

func InitCommandLineRunner_Flags_WarningGuardErrorOptionHandler(this_ *CommandLineRunner_Flags_WarningGuardErrorOptionHandler, parser *j2gort.Ext_org_kohsuke_args4j_CmdLineParser, option j2gort.Ext_org_kohsuke_args4j_OptionDef, setter j2gort.Ext_org_kohsuke_args4j_spi_Setter) {
	_ = parser
	_ = option
	_ = setter
	j2gort.InitExt_org_kohsuke_args4j_spi_StringOptionHandler(&this_.Ext_org_kohsuke_args4j_spi_StringOptionHandlerBase, parser, option, NewCommandLineRunner_Flags_MultiFlagSetter(setter, j2gort.PtrToAny(StaticGet_CheckLevel_ERROR()), StaticGet_CommandLineRunner_Flags_GuardLevels()))
	this_.initFields_CommandLineRunner_Flags_WarningGuardErrorOptionHandler()
}

// Java class: com.google.javascript.jscomp.CommandLineRunner.Flags.WarningGuardWarningOptionHandler
type CommandLineRunner_Flags_WarningGuardWarningOptionHandler struct {
	j2gort.Ext_org_kohsuke_args4j_spi_StringOptionHandlerBase
	idHash_ atomic.Int32
}

func (this_ *CommandLineRunner_Flags_WarningGuardWarningOptionHandler) BindSelf_(s any) {
	_ = s
	this_.Ext_org_kohsuke_args4j_spi_StringOptionHandlerBase.BindExtSelf_(s)
}
func (this_ *CommandLineRunner_Flags_WarningGuardWarningOptionHandler) IsJavaNil_() bool {
	return this_ == nil
}
func (*CommandLineRunner_Flags_WarningGuardWarningOptionHandler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CommandLineRunner$Flags$WarningGuardWarningOptionHandler", "WarningGuardWarningOptionHandler"
}
func (this_ *CommandLineRunner_Flags_WarningGuardWarningOptionHandler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CommandLineRunner.Flags.WarningGuardWarningOptionHandler", "com.google.javascript.jscomp.CommandLineRunner$Flags$WarningGuardWarningOptionHandler", func(v_ any) bool {
		return j2gort.InstanceOf[*CommandLineRunner_Flags_WarningGuardWarningOptionHandler](v_)
	})
}
func (*CommandLineRunner_Flags_WarningGuardWarningOptionHandler) Implements_j2gort_Ext_org_kohsuke_args4j_spi_StringOptionHandler() {
}
func (*CommandLineRunner_Flags_WarningGuardWarningOptionHandler) Implements_j2gort_Ext_org_kohsuke_args4j_spi_OptionHandler() {
}

func (this_ *CommandLineRunner_Flags_WarningGuardWarningOptionHandler) initFields_CommandLineRunner_Flags_WarningGuardWarningOptionHandler() {
	_ = this_
}

// Java constructor: WarningGuardWarningOptionHandler(org.kohsuke.args4j.CmdLineParser parser, org.kohsuke.args4j.OptionDef option, org.kohsuke.args4j.spi.Setter<? super java.lang.String> setter)
func NewCommandLineRunner_Flags_WarningGuardWarningOptionHandler(parser *j2gort.Ext_org_kohsuke_args4j_CmdLineParser, option j2gort.Ext_org_kohsuke_args4j_OptionDef, setter j2gort.Ext_org_kohsuke_args4j_spi_Setter) *CommandLineRunner_Flags_WarningGuardWarningOptionHandler {
	this_ := &CommandLineRunner_Flags_WarningGuardWarningOptionHandler{}
	this_.BindSelf_(this_)
	InitCommandLineRunner_Flags_WarningGuardWarningOptionHandler(this_, parser, option, setter)
	return this_
}

func InitCommandLineRunner_Flags_WarningGuardWarningOptionHandler(this_ *CommandLineRunner_Flags_WarningGuardWarningOptionHandler, parser *j2gort.Ext_org_kohsuke_args4j_CmdLineParser, option j2gort.Ext_org_kohsuke_args4j_OptionDef, setter j2gort.Ext_org_kohsuke_args4j_spi_Setter) {
	_ = parser
	_ = option
	_ = setter
	j2gort.InitExt_org_kohsuke_args4j_spi_StringOptionHandler(&this_.Ext_org_kohsuke_args4j_spi_StringOptionHandlerBase, parser, option, NewCommandLineRunner_Flags_MultiFlagSetter(setter, j2gort.PtrToAny(StaticGet_CheckLevel_WARNING()), StaticGet_CommandLineRunner_Flags_GuardLevels()))
	this_.initFields_CommandLineRunner_Flags_WarningGuardWarningOptionHandler()
}

// Java class: com.google.javascript.jscomp.CommandLineRunner.Flags.WarningGuardOffOptionHandler
type CommandLineRunner_Flags_WarningGuardOffOptionHandler struct {
	j2gort.Ext_org_kohsuke_args4j_spi_StringOptionHandlerBase
	idHash_ atomic.Int32
}

func (this_ *CommandLineRunner_Flags_WarningGuardOffOptionHandler) BindSelf_(s any) {
	_ = s
	this_.Ext_org_kohsuke_args4j_spi_StringOptionHandlerBase.BindExtSelf_(s)
}
func (this_ *CommandLineRunner_Flags_WarningGuardOffOptionHandler) IsJavaNil_() bool {
	return this_ == nil
}
func (*CommandLineRunner_Flags_WarningGuardOffOptionHandler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CommandLineRunner$Flags$WarningGuardOffOptionHandler", "WarningGuardOffOptionHandler"
}
func (this_ *CommandLineRunner_Flags_WarningGuardOffOptionHandler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CommandLineRunner.Flags.WarningGuardOffOptionHandler", "com.google.javascript.jscomp.CommandLineRunner$Flags$WarningGuardOffOptionHandler", func(v_ any) bool {
		return j2gort.InstanceOf[*CommandLineRunner_Flags_WarningGuardOffOptionHandler](v_)
	})
}
func (*CommandLineRunner_Flags_WarningGuardOffOptionHandler) Implements_j2gort_Ext_org_kohsuke_args4j_spi_StringOptionHandler() {
}
func (*CommandLineRunner_Flags_WarningGuardOffOptionHandler) Implements_j2gort_Ext_org_kohsuke_args4j_spi_OptionHandler() {
}

func (this_ *CommandLineRunner_Flags_WarningGuardOffOptionHandler) initFields_CommandLineRunner_Flags_WarningGuardOffOptionHandler() {
	_ = this_
}

// Java constructor: WarningGuardOffOptionHandler(org.kohsuke.args4j.CmdLineParser parser, org.kohsuke.args4j.OptionDef option, org.kohsuke.args4j.spi.Setter<? super java.lang.String> setter)
func NewCommandLineRunner_Flags_WarningGuardOffOptionHandler(parser *j2gort.Ext_org_kohsuke_args4j_CmdLineParser, option j2gort.Ext_org_kohsuke_args4j_OptionDef, setter j2gort.Ext_org_kohsuke_args4j_spi_Setter) *CommandLineRunner_Flags_WarningGuardOffOptionHandler {
	this_ := &CommandLineRunner_Flags_WarningGuardOffOptionHandler{}
	this_.BindSelf_(this_)
	InitCommandLineRunner_Flags_WarningGuardOffOptionHandler(this_, parser, option, setter)
	return this_
}

func InitCommandLineRunner_Flags_WarningGuardOffOptionHandler(this_ *CommandLineRunner_Flags_WarningGuardOffOptionHandler, parser *j2gort.Ext_org_kohsuke_args4j_CmdLineParser, option j2gort.Ext_org_kohsuke_args4j_OptionDef, setter j2gort.Ext_org_kohsuke_args4j_spi_Setter) {
	_ = parser
	_ = option
	_ = setter
	j2gort.InitExt_org_kohsuke_args4j_spi_StringOptionHandler(&this_.Ext_org_kohsuke_args4j_spi_StringOptionHandlerBase, parser, option, NewCommandLineRunner_Flags_MultiFlagSetter(setter, j2gort.PtrToAny(StaticGet_CheckLevel_OFF()), StaticGet_CommandLineRunner_Flags_GuardLevels()))
	this_.initFields_CommandLineRunner_Flags_WarningGuardOffOptionHandler()
}

// Java class: com.google.javascript.jscomp.CommandLineRunner.Flags.JsOptionHandler
type CommandLineRunner_Flags_JsOptionHandler struct {
	j2gort.Ext_org_kohsuke_args4j_spi_StringOptionHandlerBase
	idHash_ atomic.Int32
}

func (this_ *CommandLineRunner_Flags_JsOptionHandler) BindSelf_(s any) {
	_ = s
	this_.Ext_org_kohsuke_args4j_spi_StringOptionHandlerBase.BindExtSelf_(s)
}
func (this_ *CommandLineRunner_Flags_JsOptionHandler) IsJavaNil_() bool { return this_ == nil }
func (*CommandLineRunner_Flags_JsOptionHandler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CommandLineRunner$Flags$JsOptionHandler", "JsOptionHandler"
}
func (this_ *CommandLineRunner_Flags_JsOptionHandler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CommandLineRunner.Flags.JsOptionHandler", "com.google.javascript.jscomp.CommandLineRunner$Flags$JsOptionHandler", func(v_ any) bool {
		return j2gort.InstanceOf[*CommandLineRunner_Flags_JsOptionHandler](v_)
	})
}
func (*CommandLineRunner_Flags_JsOptionHandler) Implements_j2gort_Ext_org_kohsuke_args4j_spi_StringOptionHandler() {
}
func (*CommandLineRunner_Flags_JsOptionHandler) Implements_j2gort_Ext_org_kohsuke_args4j_spi_OptionHandler() {
}

func (this_ *CommandLineRunner_Flags_JsOptionHandler) initFields_CommandLineRunner_Flags_JsOptionHandler() {
	_ = this_
}

// Java constructor: JsOptionHandler(org.kohsuke.args4j.CmdLineParser parser, org.kohsuke.args4j.OptionDef option, org.kohsuke.args4j.spi.Setter<? super java.lang.String> setter)
func NewCommandLineRunner_Flags_JsOptionHandler(parser *j2gort.Ext_org_kohsuke_args4j_CmdLineParser, option j2gort.Ext_org_kohsuke_args4j_OptionDef, setter j2gort.Ext_org_kohsuke_args4j_spi_Setter) *CommandLineRunner_Flags_JsOptionHandler {
	this_ := &CommandLineRunner_Flags_JsOptionHandler{}
	this_.BindSelf_(this_)
	InitCommandLineRunner_Flags_JsOptionHandler(this_, parser, option, setter)
	return this_
}

func InitCommandLineRunner_Flags_JsOptionHandler(this_ *CommandLineRunner_Flags_JsOptionHandler, parser *j2gort.Ext_org_kohsuke_args4j_CmdLineParser, option j2gort.Ext_org_kohsuke_args4j_OptionDef, setter j2gort.Ext_org_kohsuke_args4j_spi_Setter) {
	_ = parser
	_ = option
	_ = setter
	j2gort.InitExt_org_kohsuke_args4j_spi_StringOptionHandler(&this_.Ext_org_kohsuke_args4j_spi_StringOptionHandlerBase, parser, option, NewCommandLineRunner_Flags_MultiFlagSetter(setter, j2gort.BoxNullableVal(StaticGet_AbstractCommandLineRunner_JsSourceType_JS()), StaticGet_CommandLineRunner_Flags_MixedJsSources()))
	this_.initFields_CommandLineRunner_Flags_JsOptionHandler()
}

// Java class: com.google.javascript.jscomp.CommandLineRunner.Flags.JsZipOptionHandler
type CommandLineRunner_Flags_JsZipOptionHandler struct {
	j2gort.Ext_org_kohsuke_args4j_spi_StringOptionHandlerBase
	idHash_ atomic.Int32
}

func (this_ *CommandLineRunner_Flags_JsZipOptionHandler) BindSelf_(s any) {
	_ = s
	this_.Ext_org_kohsuke_args4j_spi_StringOptionHandlerBase.BindExtSelf_(s)
}
func (this_ *CommandLineRunner_Flags_JsZipOptionHandler) IsJavaNil_() bool { return this_ == nil }
func (*CommandLineRunner_Flags_JsZipOptionHandler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CommandLineRunner$Flags$JsZipOptionHandler", "JsZipOptionHandler"
}
func (this_ *CommandLineRunner_Flags_JsZipOptionHandler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CommandLineRunner.Flags.JsZipOptionHandler", "com.google.javascript.jscomp.CommandLineRunner$Flags$JsZipOptionHandler", func(v_ any) bool {
		return j2gort.InstanceOf[*CommandLineRunner_Flags_JsZipOptionHandler](v_)
	})
}
func (*CommandLineRunner_Flags_JsZipOptionHandler) Implements_j2gort_Ext_org_kohsuke_args4j_spi_StringOptionHandler() {
}
func (*CommandLineRunner_Flags_JsZipOptionHandler) Implements_j2gort_Ext_org_kohsuke_args4j_spi_OptionHandler() {
}

func (this_ *CommandLineRunner_Flags_JsZipOptionHandler) initFields_CommandLineRunner_Flags_JsZipOptionHandler() {
	_ = this_
}

// Java constructor: JsZipOptionHandler(org.kohsuke.args4j.CmdLineParser parser, org.kohsuke.args4j.OptionDef option, org.kohsuke.args4j.spi.Setter<? super java.lang.String> setter)
func NewCommandLineRunner_Flags_JsZipOptionHandler(parser *j2gort.Ext_org_kohsuke_args4j_CmdLineParser, option j2gort.Ext_org_kohsuke_args4j_OptionDef, setter j2gort.Ext_org_kohsuke_args4j_spi_Setter) *CommandLineRunner_Flags_JsZipOptionHandler {
	this_ := &CommandLineRunner_Flags_JsZipOptionHandler{}
	this_.BindSelf_(this_)
	InitCommandLineRunner_Flags_JsZipOptionHandler(this_, parser, option, setter)
	return this_
}

func InitCommandLineRunner_Flags_JsZipOptionHandler(this_ *CommandLineRunner_Flags_JsZipOptionHandler, parser *j2gort.Ext_org_kohsuke_args4j_CmdLineParser, option j2gort.Ext_org_kohsuke_args4j_OptionDef, setter j2gort.Ext_org_kohsuke_args4j_spi_Setter) {
	_ = parser
	_ = option
	_ = setter
	j2gort.InitExt_org_kohsuke_args4j_spi_StringOptionHandler(&this_.Ext_org_kohsuke_args4j_spi_StringOptionHandlerBase, parser, option, NewCommandLineRunner_Flags_MultiFlagSetter(setter, j2gort.BoxNullableVal(StaticGet_AbstractCommandLineRunner_JsSourceType_JS_ZIP()), StaticGet_CommandLineRunner_Flags_MixedJsSources()))
	this_.initFields_CommandLineRunner_Flags_JsZipOptionHandler()
}

// Java class: com.google.javascript.jscomp.CommandLineRunner.Flags.MultiFlagSetter
type CommandLineRunner_Flags_MultiFlagSetter struct {
	proxy   j2gort.Ext_org_kohsuke_args4j_spi_Setter `gson:"proxy"`
	flag    any                                      `gson:"flag"`
	entries j2gort.Ext_java_util_List                `gson:"entries"`
	idHash_ atomic.Int32
}

func (this_ *CommandLineRunner_Flags_MultiFlagSetter) BindSelf_(s any) {
	_ = s
}
func (this_ *CommandLineRunner_Flags_MultiFlagSetter) IsJavaNil_() bool { return this_ == nil }
func (*CommandLineRunner_Flags_MultiFlagSetter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CommandLineRunner$Flags$MultiFlagSetter", "MultiFlagSetter"
}
func (this_ *CommandLineRunner_Flags_MultiFlagSetter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CommandLineRunner_Flags_MultiFlagSetter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CommandLineRunner.Flags.MultiFlagSetter", "com.google.javascript.jscomp.CommandLineRunner$Flags$MultiFlagSetter", func(v_ any) bool {
		return j2gort.InstanceOf[*CommandLineRunner_Flags_MultiFlagSetter](v_)
	})
}
func (*CommandLineRunner_Flags_MultiFlagSetter) Implements_j2gort_Ext_org_kohsuke_args4j_spi_Setter() {
}

func (this_ *CommandLineRunner_Flags_MultiFlagSetter) initFields_CommandLineRunner_Flags_MultiFlagSetter() {
	_ = this_
}

// Java constructor: MultiFlagSetter(org.kohsuke.args4j.spi.Setter<? super java.lang.String> proxy, T flag, java.util.List<com.google.javascript.jscomp.AbstractCommandLineRunner.FlagEntry<T>> entries)
func NewCommandLineRunner_Flags_MultiFlagSetter(proxy j2gort.Ext_org_kohsuke_args4j_spi_Setter, flag any, entries j2gort.Ext_java_util_List) *CommandLineRunner_Flags_MultiFlagSetter {
	this_ := &CommandLineRunner_Flags_MultiFlagSetter{}
	InitCommandLineRunner_Flags_MultiFlagSetter(this_, proxy, flag, entries)
	return this_
}

func InitCommandLineRunner_Flags_MultiFlagSetter(this_ *CommandLineRunner_Flags_MultiFlagSetter, proxy j2gort.Ext_org_kohsuke_args4j_spi_Setter, flag any, entries j2gort.Ext_java_util_List) {
	_ = proxy
	_ = flag
	_ = entries
	this_.initFields_CommandLineRunner_Flags_MultiFlagSetter()
	this_.proxy = proxy
	this_.flag = flag
	this_.entries = entries
}

// Java method: MultiFlagSetter.isMultiValued()
func (this_ *CommandLineRunner_Flags_MultiFlagSetter) IsMultiValued() bool {
	return this_.proxy.IsMultiValued()
}

// Java method: MultiFlagSetter.getType()
func (this_ *CommandLineRunner_Flags_MultiFlagSetter) GetType() *javalang.Class {
	return this_.proxy.GetType()
}

// Java method: MultiFlagSetter.addValue(java.lang.String value)
func (this_ *CommandLineRunner_Flags_MultiFlagSetter) AddValue(arg_53 any) {
	value := j2gort.CastString(arg_53)
	_ = value
	var normalizedValue string = value
	_ = normalizedValue
	if (((value != j2gort.NullString) && (j2gort.JStrLen(value) > int32(0))) && ((j2gort.CheckNotNullString(j2gort.JStrSub(value, int32(0), int32(1))) == "'") || (j2gort.CheckNotNullString(j2gort.JStrSub(value, int32(0), int32(1))) == "\""))) && func() bool {
		sr_54 := j2gort.JStrSubFrom(value, (j2gort.JStrLen(value) - int32(1)))
		se_55 := j2gort.JStrSub(value, int32(0), int32(1))
		return j2gort.CheckNotNullString(sr_54) == se_55
	}() {
		normalizedValue = j2gort.JStrSub(value, int32(1), (j2gort.JStrLen(value) - int32(1)))
	}
	this_.proxy.AddValue(j2gort.BoxString(normalizedValue))
	this_.entries.Add_1(NewAbstractCommandLineRunner_FlagEntry(this_.flag, normalizedValue))
}

// Java method: MultiFlagSetter.asFieldSetter()
func (this_ *CommandLineRunner_Flags_MultiFlagSetter) AsFieldSetter() *j2gort.Ext_org_kohsuke_args4j_spi_FieldSetter {
	return this_.proxy.AsFieldSetter()
}

// Java method: MultiFlagSetter.asAnnotatedElement()
func (this_ *CommandLineRunner_Flags_MultiFlagSetter) AsAnnotatedElement() j2gort.Ext_java_lang_reflect_AnnotatedElement {
	return this_.proxy.AsAnnotatedElement()
}

// Java enum: com.google.javascript.jscomp.CommandLineRunner.FormattingOption
type CommandLineRunner_FormattingOption uint8

const (
	enum_CommandLineRunner_FormattingOption_PRETTY_PRINT          CommandLineRunner_FormattingOption = 1
	enum_CommandLineRunner_FormattingOption_PRINT_INPUT_DELIMITER CommandLineRunner_FormattingOption = 2
	enum_CommandLineRunner_FormattingOption_SINGLE_QUOTES         CommandLineRunner_FormattingOption = 3
)

var enumBaseTable_CommandLineRunner_FormattingOption = [...]j2gort.EnumBase{
	{Name_: "PRETTY_PRINT", Ordinal_: 0},
	{Name_: "PRINT_INPUT_DELIMITER", Ordinal_: 1},
	{Name_: "SINGLE_QUOTES", Ordinal_: 2},
}
var enumValues_CommandLineRunner_FormattingOption = []CommandLineRunner_FormattingOption{1, 2, 3}

func StaticGet_CommandLineRunner_FormattingOption_PRETTY_PRINT() CommandLineRunner_FormattingOption {
	return enum_CommandLineRunner_FormattingOption_PRETTY_PRINT
}
func StaticGet_CommandLineRunner_FormattingOption_PRINT_INPUT_DELIMITER() CommandLineRunner_FormattingOption {
	return enum_CommandLineRunner_FormattingOption_PRINT_INPUT_DELIMITER
}
func StaticGet_CommandLineRunner_FormattingOption_SINGLE_QUOTES() CommandLineRunner_FormattingOption {
	return enum_CommandLineRunner_FormattingOption_SINGLE_QUOTES
}

func CommandLineRunner_FormattingOption_Values() []CommandLineRunner_FormattingOption {
	return j2gort.CloneSlice(enumValues_CommandLineRunner_FormattingOption)
}

func CommandLineRunner_FormattingOption_ValueOf(name string) CommandLineRunner_FormattingOption {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CommandLineRunner_FormattingOption {
		if enumBaseTable_CommandLineRunner_FormattingOption[i_].Name_ == name {
			return CommandLineRunner_FormattingOption(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CommandLineRunner.FormattingOption." + name))
}

// Java method: FormattingOption.applyToOptions(com.google.javascript.jscomp.CompilerOptions options)
func (this_ CommandLineRunner_FormattingOption) priv_FormattingOption_applyToOptions_1427(options *CompilerOptions) {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = options
	switch (this_).Ordinal() {
	case int32(0) /* FormattingOption.PRETTY_PRINT */ :
		options.SetPrettyPrint(true)
	case int32(1) /* FormattingOption.PRINT_INPUT_DELIMITER */ :
		options.SetPrintInputDelimiter(true)
	case int32(2) /* FormattingOption.SINGLE_QUOTES */ :
		options.SetPreferSingleQuotes(true)
	}
}

func (this_ CommandLineRunner_FormattingOption) IsJavaNil_() bool                              { return this_ == 0 }
func (CommandLineRunner_FormattingOption) HasIdentityEquals_()                                 {}
func (CommandLineRunner_FormattingOption) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CommandLineRunner_FormattingOption) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CommandLineRunner_FormattingOption) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CommandLineRunner_FormattingOption) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CommandLineRunner$FormattingOption", "FormattingOption"
}
func (this_ CommandLineRunner_FormattingOption) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CommandLineRunner_FormattingOption[this_-1]
}
func (this_ CommandLineRunner_FormattingOption) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CommandLineRunner_FormattingOption) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CommandLineRunner_FormattingOption) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CommandLineRunner_FormattingOption[this_-1].Name_
}
func (this_ CommandLineRunner_FormattingOption) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CommandLineRunner_FormattingOption[this_-1].Name_
}
func (this_ CommandLineRunner_FormattingOption) Equals(other_ any) bool {
	o_, ok_ := other_.(CommandLineRunner_FormattingOption)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CommandLineRunner_FormattingOption) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CommandLineRunner_FormattingOption[this_-1].HashCode()
}
func (this_ CommandLineRunner_FormattingOption) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CommandLineRunner_FormattingOption[this_-1].CompareTo(other_)
}
func (this_ CommandLineRunner_FormattingOption) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CommandLineRunner_FormattingOption[this_-1].Name_)
}
func init() {
	vals_ := []any{CommandLineRunner_FormattingOption(1), CommandLineRunner_FormattingOption(2), CommandLineRunner_FormattingOption(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CommandLineRunner.FormattingOption", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CommandLineRunner$FormattingOption", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CommandLineRunner.FormattingOption", "com.google.javascript.jscomp.CommandLineRunner$FormattingOption", func(v_ any) bool {
		return j2gort.InstanceOf[CommandLineRunner_FormattingOption](v_)
	})
}

// Java class:
type CommandLineRunner__Anon_11 struct {
	j2gort.Ext_java_nio_file_SimpleFileVisitorBase
	cap_0_matcher     j2gort.Ext_java_nio_file_PathMatcher
	cap_2_excludes    j2gort.Ext_java_util_Set
	cap_3_allJsInputs j2gort.Ext_java_util_Map
	idHash_           atomic.Int32
	cap_1_remove      bool
}

func (this_ *CommandLineRunner__Anon_11) BindSelf_(s any) {
	_ = s
	this_.Ext_java_nio_file_SimpleFileVisitorBase.BindExtSelf_(s)
}
func (this_ *CommandLineRunner__Anon_11) IsJavaNil_() bool { return this_ == nil }
func (*CommandLineRunner__Anon_11) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CommandLineRunner$1", ""
}
func (this_ *CommandLineRunner__Anon_11) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CommandLineRunner__Anon_11) Implements_j2gort_Ext_java_nio_file_SimpleFileVisitor() {}
func (*CommandLineRunner__Anon_11) Implements_j2gort_Ext_java_nio_file_FileVisitor()       {}

func (this_ *CommandLineRunner__Anon_11) initFields_CommandLineRunner__Anon_11() {
	_ = this_
}

func NewCommandLineRunner__Anon_11(cap_p_matcher j2gort.Ext_java_nio_file_PathMatcher, cap_p_remove bool, cap_p_excludes j2gort.Ext_java_util_Set, cap_p_allJsInputs j2gort.Ext_java_util_Map) *CommandLineRunner__Anon_11 {
	this_ := &CommandLineRunner__Anon_11{}
	this_.BindSelf_(this_)
	this_.cap_0_matcher = cap_p_matcher
	this_.cap_1_remove = cap_p_remove
	this_.cap_2_excludes = cap_p_excludes
	this_.cap_3_allJsInputs = cap_p_allJsInputs
	j2gort.InitExt_java_nio_file_SimpleFileVisitor(&this_.Ext_java_nio_file_SimpleFileVisitorBase)
	this_.initFields_CommandLineRunner__Anon_11()
	return this_
}

// Java method: .visitFile(java.nio.file.Path p, java.nio.file.attribute.BasicFileAttributes attrs)
func (this_ *CommandLineRunner__Anon_11) VisitFile(arg_56 any, attrs j2gort.Ext_java_nio_file_attribute_BasicFileAttributes) *j2gort.Ext_java_nio_file_FileVisitResult {
	p := j2gort.AsIface[j2gort.Ext_java_nio_file_Path](arg_56)
	_ = p
	_ = attrs
	if this_.cap_0_matcher.Matches(p) || this_.cap_0_matcher.Matches(p.Normalize()) {
		var pathStringAbsolute string = p.Normalize().ToAbsolutePath().String()
		_ = pathStringAbsolute
		if this_.cap_1_remove {
			j2gort.SetAddStr(this_.cap_2_excludes, pathStringAbsolute)
			j2gort.MapRemoveStr(this_.cap_3_allJsInputs, pathStringAbsolute)
		} else if !(j2gort.SetContainsStr(this_.cap_2_excludes, pathStringAbsolute)) {
			j2gort.MapPutStr(this_.cap_3_allJsInputs, pathStringAbsolute, j2gort.BoxString(p.String()))
		}
	}
	return j2gort.Ext_java_nio_file_FileVisitResult_F_CONTINUE
}

// Java method: .visitFileFailed(java.nio.file.Path file, java.io.IOException e)
func (this_ *CommandLineRunner__Anon_11) VisitFileFailed(arg_57 any, e javaio.IIOException) *j2gort.Ext_java_nio_file_FileVisitResult {
	file := j2gort.AsIface[j2gort.Ext_java_nio_file_Path](arg_57)
	_ = file
	_ = e
	return j2gort.Ext_java_nio_file_FileVisitResult_F_SKIP_SUBTREE
}
