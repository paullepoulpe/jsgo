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
// Port of com/google/javascript/jscomp/AbstractCommandLineRunner.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.AbstractCommandLineRunner
type AbstractCommandLineRunner interface {
	Implements_jscomp_AbstractCommandLineRunner()
	IsJavaNil_() bool
	GetAbstractCommandLineRunnerBase() *AbstractCommandLineRunnerBase
	EnableTestMode(j2gort.Ext_com_google_common_base_Supplier, j2gort.Ext_com_google_common_base_Supplier, j2gort.Ext_com_google_common_base_Supplier, j2gort.Ext_com_google_common_base_Function)
	SetExitCodeReceiver(j2gort.Ext_com_google_common_base_Function)
	IsInTestMode() bool
	GetCommandLineConfig() *AbstractCommandLineRunner_CommandLineConfig
	CreateCompiler() *Compiler
	PrepForBundleAndAppendTo(j2gort.Ext_java_lang_Appendable, *CompilerInput, string)
	AppendRuntimeTo(j2gort.Ext_java_lang_Appendable)
	CreateOptions() *CompilerOptions
	GetDiagnosticGroups() *DiagnosticGroups
	AddAllowlistWarningsGuard(*CompilerOptions, *j2gort.Ext_java_io_File)
	SetRunOptions(*CompilerOptions)
	GetCompiler() *Compiler
	GetVersionText() string
	Run()
	GetErrorPrintStream() *j2gort.Ext_java_io_PrintStream
	ParseJsonFilesFromInputStream() j2gort.Ext_java_util_List
	CreateInputs_4(j2gort.Ext_java_util_List, j2gort.Ext_java_util_List, bool, j2gort.Ext_java_util_List) j2gort.Ext_java_util_List
	CheckChunkName(string)
	GetChunkOutputFileName(*JSChunk) string
	GetChunkConformanceFileName(*JSChunk) string
	WriteChunkOutput(string, j2gort.Ext_java_lang_Appendable, CodePrinter_LicenseTracker, *JSChunk)
	WriteOutput_8(j2gort.Ext_java_lang_Appendable, *Compiler, CodePrinter_LicenseTracker, *JSChunk, string, string, j2gort.Ext_com_google_common_base_Function, string)
	WriteOutput_7(j2gort.Ext_java_lang_Appendable, *Compiler, string, string, string, j2gort.Ext_com_google_common_base_Function, string)
	DoRun() int32
	GetCompileMetricsRecorder() CompileMetricsRecorderInterface
	ProcessResults(*Result, j2gort.Ext_java_util_List, *CompilerOptions) int32
	GetJavascriptEscaper() j2gort.Ext_com_google_common_base_Function
	OutputSingleBinary(*CompilerOptions)
	CreateJsonFile(*CompilerOptions, string, j2gort.Ext_com_google_common_base_Function) *AbstractCommandLineRunner_JsonFileSpec
	OutputJsonStream()
	CreateExterns(*CompilerOptions) j2gort.Ext_java_util_List
	ShouldGenerateMapPerChunk(*CompilerOptions) bool
	ExpandSourceMapPath(*CompilerOptions, *JSChunk) string
	FilenameToOutputStream(string) j2gort.Ext_java_io_OutputStream
	CreateWriter(j2gort.Ext_java_io_OutputStream, j2gort.Ext_java_nio_charset_Charset) j2gort.Ext_java_io_Writer
	PrintChunkGraphJsonTo(j2gort.Ext_java_lang_Appendable)
	PrintChunkGraphManifestOrBundleTo(*JSChunkGraph, j2gort.Ext_java_lang_Appendable, bool)
	PrintManifestTo(*JSChunk, j2gort.Ext_java_lang_Appendable)
	PrintBundleTo(*JSChunk, j2gort.Ext_java_lang_Appendable)
}

func ToIface_AbstractCommandLineRunner[T interface {
	AbstractCommandLineRunner
	comparable
}](v T) AbstractCommandLineRunner {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.AbstractCommandLineRunner
type AbstractCommandLineRunnerBase struct {
	self_                       AbstractCommandLineRunner
	config                      *AbstractCommandLineRunner_CommandLineConfig      `gson:"config"`
	in                          j2gort.Ext_java_io_InputStream                    `gson:"in"`
	defaultJsOutput             *j2gort.Ext_java_io_PrintStream                   `gson:"defaultJsOutput"`
	err                         *j2gort.Ext_java_io_PrintStream                   `gson:"err"`
	compiler                    *Compiler                                         `gson:"compiler"`
	inputCharset                j2gort.Ext_java_nio_charset_Charset               `gson:"inputCharset"`
	outputCharset2              j2gort.Ext_java_nio_charset_Charset               `gson:"outputCharset2"`
	legacyOutputCharset         j2gort.Ext_java_nio_charset_Charset               `gson:"legacyOutputCharset"`
	testMode                    bool                                              `gson:"testMode"`
	externsSupplierForTesting   j2gort.Ext_com_google_common_base_Supplier        `gson:"externsSupplierForTesting"`
	inputsSupplierForTesting    j2gort.Ext_com_google_common_base_Supplier        `gson:"inputsSupplierForTesting"`
	chunksSupplierForTesting    j2gort.Ext_com_google_common_base_Supplier        `gson:"chunksSupplierForTesting"`
	exitCodeReceiver            j2gort.Ext_com_google_common_base_Function        `gson:"exitCodeReceiver"`
	rootRelativePathsMap        j2gort.Ext_java_util_Map                          `gson:"rootRelativePathsMap"`
	parsedChunkWrappers         j2gort.Ext_java_util_Map                          `gson:"parsedChunkWrappers"`
	parsedChunkOutputFiles      j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"parsedChunkOutputFiles"`
	parsedChunkConformanceFiles j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"parsedChunkConformanceFiles"`
	gson                        *j2gort.Ext_com_google_gson_Gson                  `gson:"gson"`
	filesToStreamOut            []*AbstractCommandLineRunner_JsonFileSpec         `gson:"filesToStreamOut"`
	idHash_                     atomic.Int32
}

func (this_ *AbstractCommandLineRunnerBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(AbstractCommandLineRunner)
}
func (this_ *AbstractCommandLineRunnerBase) IsJavaNil_() bool { return this_ == nil }
func (*AbstractCommandLineRunnerBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractCommandLineRunner", "AbstractCommandLineRunner"
}
func (this_ *AbstractCommandLineRunnerBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractCommandLineRunner", "com.google.javascript.jscomp.AbstractCommandLineRunner", func(v_ any) bool {
		return j2gort.InstanceOf[AbstractCommandLineRunner](v_)
	})
}
func (b *AbstractCommandLineRunnerBase) GetAbstractCommandLineRunnerBase() *AbstractCommandLineRunnerBase {
	return b
}
func (*AbstractCommandLineRunnerBase) Implements_jscomp_AbstractCommandLineRunner() {}

var static_AbstractCommandLineRunner_OUTPUT_SAME_AS_INPUT_ERROR *DiagnosticType

func StaticGet_AbstractCommandLineRunner_OUTPUT_SAME_AS_INPUT_ERROR() *DiagnosticType {
	clinit_AbstractCommandLineRunner()
	return static_AbstractCommandLineRunner_OUTPUT_SAME_AS_INPUT_ERROR
}
func StaticSet_AbstractCommandLineRunner_OUTPUT_SAME_AS_INPUT_ERROR(v *DiagnosticType) {
	clinit_AbstractCommandLineRunner()
	static_AbstractCommandLineRunner_OUTPUT_SAME_AS_INPUT_ERROR = v
}

var static_AbstractCommandLineRunner_COULD_NOT_SERIALIZE_AST *DiagnosticType

func StaticGet_AbstractCommandLineRunner_COULD_NOT_SERIALIZE_AST() *DiagnosticType {
	clinit_AbstractCommandLineRunner()
	return static_AbstractCommandLineRunner_COULD_NOT_SERIALIZE_AST
}
func StaticSet_AbstractCommandLineRunner_COULD_NOT_SERIALIZE_AST(v *DiagnosticType) {
	clinit_AbstractCommandLineRunner()
	static_AbstractCommandLineRunner_COULD_NOT_SERIALIZE_AST = v
}

var static_AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST *DiagnosticType

func StaticGet_AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST() *DiagnosticType {
	clinit_AbstractCommandLineRunner()
	return static_AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST
}
func StaticSet_AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST(v *DiagnosticType) {
	clinit_AbstractCommandLineRunner()
	static_AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST = v
}

var static_AbstractCommandLineRunner_NO_TREE_GENERATED_ERROR *DiagnosticType

func StaticGet_AbstractCommandLineRunner_NO_TREE_GENERATED_ERROR() *DiagnosticType {
	clinit_AbstractCommandLineRunner()
	return static_AbstractCommandLineRunner_NO_TREE_GENERATED_ERROR
}
func StaticSet_AbstractCommandLineRunner_NO_TREE_GENERATED_ERROR(v *DiagnosticType) {
	clinit_AbstractCommandLineRunner()
	static_AbstractCommandLineRunner_NO_TREE_GENERATED_ERROR = v
}

var static_AbstractCommandLineRunner_INVALID_CHUNK_SOURCEMAP_PATTERN *DiagnosticType

func StaticGet_AbstractCommandLineRunner_INVALID_CHUNK_SOURCEMAP_PATTERN() *DiagnosticType {
	clinit_AbstractCommandLineRunner()
	return static_AbstractCommandLineRunner_INVALID_CHUNK_SOURCEMAP_PATTERN
}
func StaticSet_AbstractCommandLineRunner_INVALID_CHUNK_SOURCEMAP_PATTERN(v *DiagnosticType) {
	clinit_AbstractCommandLineRunner()
	static_AbstractCommandLineRunner_INVALID_CHUNK_SOURCEMAP_PATTERN = v
}

var static_AbstractCommandLineRunner_EXPECTED_DIAGNOSTIC_NOT_FOUND *DiagnosticType

func StaticGet_AbstractCommandLineRunner_EXPECTED_DIAGNOSTIC_NOT_FOUND() *DiagnosticType {
	clinit_AbstractCommandLineRunner()
	return static_AbstractCommandLineRunner_EXPECTED_DIAGNOSTIC_NOT_FOUND
}
func StaticSet_AbstractCommandLineRunner_EXPECTED_DIAGNOSTIC_NOT_FOUND(v *DiagnosticType) {
	clinit_AbstractCommandLineRunner()
	static_AbstractCommandLineRunner_EXPECTED_DIAGNOSTIC_NOT_FOUND = v
}

var static_AbstractCommandLineRunner_AMBIGUOUS_EXPECTATION *DiagnosticType

func StaticGet_AbstractCommandLineRunner_AMBIGUOUS_EXPECTATION() *DiagnosticType {
	clinit_AbstractCommandLineRunner()
	return static_AbstractCommandLineRunner_AMBIGUOUS_EXPECTATION
}
func StaticSet_AbstractCommandLineRunner_AMBIGUOUS_EXPECTATION(v *DiagnosticType) {
	clinit_AbstractCommandLineRunner()
	static_AbstractCommandLineRunner_AMBIGUOUS_EXPECTATION = v
}
func StaticGet_AbstractCommandLineRunner_WAITING_FOR_INPUT_WARNING() string {
	return "The compiler is waiting for input via stdin."
}
func StaticGet_AbstractCommandLineRunner_GZIPPED_TYPEDAST_BUFFER_SIZE() int32 { return int32(8388608) }
func StaticGet_AbstractCommandLineRunner_OUTPUT_MARKER() string               { return "%output%" }
func StaticGet_AbstractCommandLineRunner_OUTPUT_MARKER_JS_STRING() string     { return "%output|jsstring%" }

var clinitDone_AbstractCommandLineRunner uint32
var clinitState_AbstractCommandLineRunner j2gort.ClinitState

func clinit_AbstractCommandLineRunner() {
	if atomic.LoadUint32(&clinitDone_AbstractCommandLineRunner) != 0 {
		return
	}
	clinitSlow_AbstractCommandLineRunner()
}

func clinitSlow_AbstractCommandLineRunner() {
	if !j2gort.EnterClinit(&clinitState_AbstractCommandLineRunner, &clinitDone_AbstractCommandLineRunner) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_AbstractCommandLineRunner, &clinitDone_AbstractCommandLineRunner)
	static_AbstractCommandLineRunner_OUTPUT_SAME_AS_INPUT_ERROR = DiagnosticType_Error("JSC_OUTPUT_SAME_AS_INPUT_ERROR", "Bad output file (already listed as input file): {0}")
	static_AbstractCommandLineRunner_COULD_NOT_SERIALIZE_AST = DiagnosticType_Error("JSC_COULD_NOT_SERIALIZE_AST", "Could not serialize ast to: {0}")
	static_AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST = DiagnosticType_Error("JSC_COULD_NOT_DESERIALIZE_AST", "Could not deserialize ast from: {0}")
	static_AbstractCommandLineRunner_NO_TREE_GENERATED_ERROR = DiagnosticType_Error("JSC_NO_TREE_GENERATED_ERROR", "Code contains errors. No tree was generated.")
	static_AbstractCommandLineRunner_INVALID_CHUNK_SOURCEMAP_PATTERN = DiagnosticType_Error("JSC_INVALID_CHUNK_SOURCEMAP_PATTERN", "When using --chunk or --module flags, the --create_source_map flag must contain %outname% in the value.")
	static_AbstractCommandLineRunner_EXPECTED_DIAGNOSTIC_NOT_FOUND = DiagnosticType_Error("JSC_EXPECTED_DIAGNOSTIC_NOT_FOUND", "Expected diagnostic not found: {0}")
	static_AbstractCommandLineRunner_AMBIGUOUS_EXPECTATION = DiagnosticType_Error("JSC_AMBIGUOUS_EXPECTATION", "Multiple expected diagnostics matched the error: \"{0}\". Matches: \"{1}\"")
}

func (this_ *AbstractCommandLineRunnerBase) initFields_AbstractCommandLineRunner() {
	_ = this_

	this_.exitCodeReceiver = j2gort.ToIface_Ext_com_google_common_base_Function(StaticGet_AbstractCommandLineRunner_SystemExitCodeReceiver_INSTANCE())

	this_.filesToStreamOut = []*AbstractCommandLineRunner_JsonFileSpec{}
}

func InitAbstractCommandLineRunner(this_ *AbstractCommandLineRunnerBase) {
	InitAbstractCommandLineRunner_3(this_, j2gort.Ext_java_lang_System_F_In, j2gort.Ext_java_lang_System_F_Out, j2gort.Ext_java_lang_System_F_Err)
}

func InitAbstractCommandLineRunner_2(this_ *AbstractCommandLineRunnerBase, out *j2gort.Ext_java_io_PrintStream, err *j2gort.Ext_java_io_PrintStream) {
	_ = out
	_ = err
	InitAbstractCommandLineRunner_3(this_, j2gort.Ext_java_lang_System_F_In, out, err)
}

func InitAbstractCommandLineRunner_3(this_ *AbstractCommandLineRunnerBase, in j2gort.Ext_java_io_InputStream, out *j2gort.Ext_java_io_PrintStream, err *j2gort.Ext_java_io_PrintStream) {
	_ = in
	_ = out
	_ = err
	this_.initFields_AbstractCommandLineRunner()
	this_.config = NewAbstractCommandLineRunner_CommandLineConfig()
	this_.in = j2gort.CheckNotNull(in)
	this_.defaultJsOutput = j2gort.CheckNotNullPtr(out)
	this_.err = j2gort.CheckNotNullPtr(err)
	this_.gson = j2gort.NewExt_com_google_gson_Gson()
}

// Java method: AbstractCommandLineRunner.enableTestMode(com.google.common.base.Supplier<java.util.List<com.google.javascript.jscomp.SourceFile>> externsSupplier, com.google.common.base.Supplier<java.util.List<com.google.javascript.jscomp.SourceFile>> inputsSupplier, com.google.common.base.Supplier<java.util.List<com.google.javascript.jscomp.JSChunk>> chunksSupplier, com.google.common.base.Function<java.lang.Integer,java.lang.Void> exitCodeReceiver)
func (this_ *AbstractCommandLineRunnerBase) EnableTestMode(externsSupplier j2gort.Ext_com_google_common_base_Supplier, inputsSupplier j2gort.Ext_com_google_common_base_Supplier, chunksSupplier j2gort.Ext_com_google_common_base_Supplier, exitCodeReceiver j2gort.Ext_com_google_common_base_Function) {
	_ = externsSupplier
	_ = inputsSupplier
	_ = chunksSupplier
	_ = exitCodeReceiver
	if !((inputsSupplier == nil) != (chunksSupplier == nil)) {
		j2gort.CheckArgument(false)
	}
	this_.testMode = true
	this_.externsSupplierForTesting = externsSupplier
	this_.inputsSupplierForTesting = inputsSupplier
	this_.chunksSupplierForTesting = chunksSupplier
	this_.exitCodeReceiver = exitCodeReceiver
}

// Java method: AbstractCommandLineRunner.setExitCodeReceiver(com.google.common.base.Function<java.lang.Integer,java.lang.Void> newExitCodeReceiver)
func (this_ *AbstractCommandLineRunnerBase) SetExitCodeReceiver(newExitCodeReceiver j2gort.Ext_com_google_common_base_Function) {
	_ = newExitCodeReceiver
	this_.exitCodeReceiver = j2gort.CheckNotNull(newExitCodeReceiver)
}

// Java method: AbstractCommandLineRunner.isInTestMode()
func (this_ *AbstractCommandLineRunnerBase) IsInTestMode() bool {
	return this_.testMode
}

// Java method: AbstractCommandLineRunner.isOutputInJson()
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_isOutputInJson_138() bool {
	return ((this_.config.jsonStreamMode == StaticGet_CompilerOptions_JsonStreamMode_OUT()) || (this_.config.jsonStreamMode == StaticGet_CompilerOptions_JsonStreamMode_BOTH()))
}

// Java method: AbstractCommandLineRunner.getCommandLineConfig()
func (this_ *AbstractCommandLineRunnerBase) GetCommandLineConfig() *AbstractCommandLineRunner_CommandLineConfig {
	return this_.config
}

// Java method: AbstractCommandLineRunner.createCompiler()
func (this_ *AbstractCommandLineRunnerBase) CreateCompiler() *Compiler {
	panic("abstract or native method invoked")
}

// Java method: AbstractCommandLineRunner.prepForBundleAndAppendTo(java.lang.Appendable out, com.google.javascript.jscomp.CompilerInput input, java.lang.String content)
func (this_ *AbstractCommandLineRunnerBase) PrepForBundleAndAppendTo(out j2gort.Ext_java_lang_Appendable, input *CompilerInput, content string) {
	_ = out
	_ = input
	_ = content
	panic("abstract or native method invoked")
}

// Java method: AbstractCommandLineRunner.appendRuntimeTo(java.lang.Appendable out)
func (this_ *AbstractCommandLineRunnerBase) AppendRuntimeTo(out j2gort.Ext_java_lang_Appendable) {
	_ = out
	panic("abstract or native method invoked")
}

// Java method: AbstractCommandLineRunner.createOptions()
func (this_ *AbstractCommandLineRunnerBase) CreateOptions() *CompilerOptions {
	panic("abstract or native method invoked")
}

// Java method: AbstractCommandLineRunner.getDiagnosticGroups()
func (this_ *AbstractCommandLineRunnerBase) GetDiagnosticGroups() *DiagnosticGroups {
	if this_.compiler == nil {
		return NewDiagnosticGroups()
	}
	return this_.compiler.GetDiagnosticGroups()
}

// Java method: AbstractCommandLineRunner.addAllowlistWarningsGuard(com.google.javascript.jscomp.CompilerOptions options, java.io.File allowlistFile)
func (this_ *AbstractCommandLineRunnerBase) AddAllowlistWarningsGuard(options *CompilerOptions, allowlistFile *j2gort.Ext_java_io_File) {
	_ = options
	_ = allowlistFile
	panic("abstract or native method invoked")
}

// Java static method: AbstractCommandLineRunner.setWarningGuardOptions(com.google.javascript.jscomp.CompilerOptions options, java.util.ArrayList<com.google.javascript.jscomp.AbstractCommandLineRunner.FlagEntry<com.google.javascript.jscomp.CheckLevel>> warningGuards, com.google.javascript.jscomp.DiagnosticGroups diagnosticGroups)
func AbstractCommandLineRunner_SetWarningGuardOptions(options *CompilerOptions, warningGuards j2gort.Ext_java_util_ArrayList, diagnosticGroups *DiagnosticGroups) {
	clinit_AbstractCommandLineRunner()
	_ = options
	_ = warningGuards
	_ = diagnosticGroups
	if warningGuards != nil {
		var groupNames j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](DiagnosticGroups_GetRegisteredGroups().KeySet_0())
		_ = groupNames
		{
			cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(warningGuards))
			for cur_1.HasNext() {
				rawElem_2 := cur_1.Next()
				entry := j2gort.CastPtr[AbstractCommandLineRunner_FlagEntry](rawElem_2)
				_ = entry
				{
					var iifeCond_1 bool
					{

						sr_3 := "*"
						se_4 := entry.value
						iifeCond_1 = j2gort.CheckNotNullString(sr_3) == se_4
					}
					if iifeCond_1 {
						{
							cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(groupNames))
							for cur_5.HasNext() {
								groupName := cur_5.NextStr()
								_ = groupName
								if !(j2gort.SetContainsStr(StaticGet_DiagnosticGroups_WildcardExcludedGroups(), groupName)) {
									diagnosticGroups.SetWarningLevel(options, groupName, j2gort.CastPtr[CheckLevel](entry.flag))
								}
							}
						}
					} else if j2gort.SetContainsStr(groupNames, entry.value) {
						diagnosticGroups.SetWarningLevel(options, entry.value, j2gort.CastPtr[CheckLevel](entry.flag))
					} else {
						panic(NewFlagUsageException(("Unknown diagnostic group: '" + j2gort.StrValString(entry.value) + "'")))
					}
				}
			}
		}
	}
}

// Java method: AbstractCommandLineRunner.setRunOptions(com.google.javascript.jscomp.CompilerOptions options)
func (this_ *AbstractCommandLineRunnerBase) SetRunOptions(options *CompilerOptions) {
	_ = options
	var diagnosticGroups *DiagnosticGroups = this_.GetDiagnosticGroups()
	_ = diagnosticGroups
	if this_.config.ShouldSaveAfterStage1() || this_.config.priv_CommandLineConfig_shouldContinueCompilation_225() {
		if options.IsChecksOnly() {
			panic(NewFlagUsageException("checks_only mode is incompatible with multi-stage compilation"))
		}
		if options.GetExternExportsPath() != j2gort.NullString {
			panic(NewFlagUsageException("generating externs from exports is incompatible with multi-stage compilation"))
		}
	}
	AbstractCommandLineRunner_SetWarningGuardOptions(options, this_.config.warningGuards, diagnosticGroups)
	if !(j2gort.CheckNotNullString(this_.config.warningsAllowFile) == "") {
		this_.self_.AddAllowlistWarningsGuard(options, j2gort.NewExt_java_io_File_1_string(this_.config.warningsAllowFile))
	}
	if !(this_.config.hideWarningsFor.IsEmpty()) {
		options.AddWarningsGuard(NewShowByPathWarningsGuard_2_slice_string_ShowByPathWarningsGuard_ShowType(func() []string {
			arr_9 := this_.config.hideWarningsFor.ToArray_1_slice_any(func() []any {
				arr_6 := []string{}
				if arr_6 == nil {
					return nil
				}
				out_7 := make([]any, len(arr_6))
				for i_, el_8 := range arr_6 {
					out_7[i_] = j2gort.BoxString(el_8)
				}
				return out_7
			}())
			if arr_9 == nil {
				return nil
			}
			out_10 := make([]string, len(arr_9))
			for i_, el_11 := range arr_9 {
				out_10[i_] = j2gort.CastString(el_11)
			}
			return out_10
		}(), StaticGet_ShowByPathWarningsGuard_ShowType_EXCLUDE()))
	}
	var define j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.config.define))
	_ = define
	if this_.config.browserFeaturesetYear != int32(0) {
		func() {
			defer func() {
				if rec_12 := recover(); rec_12 != nil {
					jex_13 := javalang.RecoverException(rec_12)
					_ = jex_13
					if j2gort.InstanceOf[javalang.IIllegalStateException](jex_13) {
						e := j2gort.AsIface[javalang.IIllegalStateException](jex_13)
						_ = e
						panic(NewFlagUsageException(j2gort.ThrowableGetMessage(e)))
					} else {
						panic(rec_12)
					}
				}
			}()
			options.SetBrowserFeaturesetYear_1_int32(this_.config.browserFeaturesetYear)
		}()
	}
	AbstractCommandLineRunner_CreateDefineReplacements(define, options)
	options.SetTweakProcessing(this_.config.tweakProcessing)
	if this_.config.dependencyOptions != nil {
		options.SetDependencyOptions(this_.config.dependencyOptions)
	}
	options.SetDevMode(this_.config.jscompDevMode)
	options.SetCodingConvention(this_.config.codingConvention)
	options.SetSummaryDetailLevel(this_.config.summaryDetailLevel)
	options.SetTrustedStrings(true)
	this_.legacyOutputCharset = this_.priv_AbstractCommandLineRunner_getLegacyOutputCharset_183()
	options.SetOutputCharset(this_.legacyOutputCharset)
	this_.outputCharset2 = this_.priv_AbstractCommandLineRunner_getOutputCharset2_184()
	this_.inputCharset = this_.priv_AbstractCommandLineRunner_getInputCharset_182()
	if j2gort.JStrLen(this_.config.jsOutputFile) > int32(0) {
		if this_.config.skipNormalOutputs {
			panic(NewFlagUsageException("skip_normal_outputs and js_output_file cannot be used together."))
		}
	}
	if this_.config.skipNormalOutputs && this_.config.printAst {
		panic(NewFlagUsageException("skip_normal_outputs and print_ast cannot be used together."))
	}
	if this_.config.skipNormalOutputs && this_.config.printTree {
		panic(NewFlagUsageException("skip_normal_outputs and print_tree cannot be used together."))
	}
	if this_.config.skipNormalOutputs && this_.config.printTreeJson {
		panic(NewFlagUsageException("skip_normal_outputs and print_tree_json_path cannot be used together."))
	}
	if j2gort.JStrLen(this_.config.createSourceMap) > int32(0) {
		options.SetSourceMapOutputPath(this_.config.createSourceMap)
	} else if this_.priv_AbstractCommandLineRunner_isOutputInJson_138() {
		options.SetSourceMapOutputPath("%outname%")
	}
	options.SetSourceMapDetailLevel(this_.config.sourceMapDetailLevel)
	options.SetSourceMapFormat(this_.config.sourceMapFormat)
	options.SetSourceMapLocationMappings(j2gort.AsIface[j2gort.Ext_java_util_List](this_.config.sourceMapLocationMappings))
	options.SetParseInlineSourceMaps(this_.config.parseInlineSourceMaps)
	options.SetApplyInputSourceMaps(this_.config.applyInputSourceMaps)
	var inputSourceMaps j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.NewExt_com_google_common_collect_ImmutableMap_Builder()
	_ = inputSourceMaps
	{
		cur_14 := j2gort.NewMapEntryCursor(this_.config.sourceMapInputFiles)
		for cur_14.NextEntry() {
			var sourceMap *SourceFile = SourceFile_S_Builder().WithKind(pkg_rhino.StaticGet_StaticSourceFile_SourceKind_NON_CODE()).WithPath_1_string(j2gort.CastString(cur_14.GetValue())).Build()
			_ = sourceMap
			inputSourceMaps.Put_2(j2gort.BoxString(cur_14.KeyStr()), NewSourceMapInput(sourceMap))
		}
	}
	options.SetInputSourceMaps(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](inputSourceMaps.BuildOrThrow()))
	if !(j2gort.CheckNotNullString(this_.config.variableMapInputFile) == "") {
		options.SetInputVariableMap(VariableMap_Load(this_.config.variableMapInputFile))
	}
	if !(j2gort.CheckNotNullString(this_.config.propertyMapInputFile) == "") {
		options.SetInputPropertyMap(VariableMap_Load(this_.config.propertyMapInputFile))
	}
	if !(this_.config.outputManifests.IsEmpty()) {
		var uniqueNames j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
		_ = uniqueNames
		{
			cur_15 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.config.outputManifests))
			for cur_15.HasNext() {
				filename := cur_15.NextStr()
				_ = filename
				if !(uniqueNames.Add(filename)) {
					panic(NewFlagUsageException(("output_manifest flags specify duplicate file names: " + j2gort.StrValString(filename))))
				}
			}
		}
	}
	if !(this_.config.outputBundles.IsEmpty()) {
		var uniqueNames_2 j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
		_ = uniqueNames_2
		{
			cur_16 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.config.outputBundles))
			for cur_16.HasNext() {
				filename_2 := cur_16.NextStr()
				_ = filename_2
				if !(uniqueNames_2.Add(filename_2)) {
					panic(NewFlagUsageException(("output_bundle flags specify duplicate file names: " + j2gort.StrValString(filename_2))))
				}
			}
		}
	}
	options.SetProcessCommonJSModules(this_.config.processCommonJSModules)
	options.SetModuleRoots(this_.config.moduleRoots)
	options.SetAngularPass(this_.config.angularPass)
	if !(j2gort.CheckNotNullString(this_.config.jsonWarningsFile) == "") {
		options.AddReportGenerator(NewJsonErrorReportGenerator(j2gort.NewExt_java_io_PrintStream_1_string(this_.config.jsonWarningsFile), ToIface_SourceExcerptProvider(this_.compiler)))
	}
	if this_.config.errorFormat == StaticGet_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_JSON() {
		clinit_JsonErrorReportGenerator()
		var val_17 JsonErrorReportGenerator
		InitJsonErrorReportGenerator(&val_17, this_.GetErrorPrintStream(), ToIface_SourceExcerptProvider(this_.compiler))
		var errorGenerator *JsonErrorReportGenerator = &val_17
		_ = errorGenerator
		this_.compiler.SetErrorManager(j2gort.AsIface[ErrorManager](NewSortingErrorManager(j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_1(errorGenerator)))))
	}
	if this_.config.printTree || this_.config.printTreeJson {
		options.SetParseJsDocDocumentation(pkg_parsing.StaticGet_Config_JsDocParsing_INCLUDE_ALL_COMMENTS())
	}
	if this_.config.skipNormalOutputs {
		this_.compiler.SetPreferRegexParser(true)
	}
}

// Java method: AbstractCommandLineRunner.getCompiler()
func (this_ *AbstractCommandLineRunnerBase) GetCompiler() *Compiler {
	return this_.compiler
}

// Java static method: AbstractCommandLineRunner.getBuiltinExterns(com.google.javascript.jscomp.CompilerOptions.Environment env)
func AbstractCommandLineRunner_GetBuiltinExterns(env CompilerOptions_Environment) j2gort.Ext_java_util_List {
	clinit_AbstractCommandLineRunner()
	_ = env
	var input j2gort.Ext_java_io_InputStream = AbstractCommandLineRunner_GetExternsInput()
	_ = input
	defer func() {
		if input != nil {
			input.Close_0()
		}
	}()
	{
		var zip *j2gort.Ext_java_util_zip_ZipInputStream = j2gort.NewExt_java_util_zip_ZipInputStream_1(input)
		_ = zip
		var envPrefix string = (j2gort.StrValString(j2gort.Ext_com_google_common_base_Ascii_ToLowerCase_1_string(j2gort.CheckNotNullEnum(env).String())) + "/")
		_ = envPrefix
		var mapFromExternsZip j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
		_ = mapFromExternsZip
		{
			var entry *j2gort.Ext_java_util_zip_ZipEntry = nil
			_ = entry
			for (func() *j2gort.Ext_java_util_zip_ZipEntry {
				var asgn_18 *j2gort.Ext_java_util_zip_ZipEntry = zip.GetNextEntry()
				entry = asgn_18
				return asgn_18
			}()) != nil {
				var filename string = entry.GetName()
				_ = filename
				if j2gort.JStrIndexOf(filename, "/") >= 0 {
					if !(javalang.StringStartsWith(j2gort.CheckNotNullString(filename), j2gort.CheckNotNullString(envPrefix))) {
						continue
					}
					filename = j2gort.JStrSubFrom(filename, j2gort.JStrLen(envPrefix))
				}
				var entryStream *j2gort.Ext_java_io_BufferedInputStream = j2gort.NewExt_java_io_BufferedInputStream_1(j2gort.Ext_com_google_common_io_ByteStreams_Limit(zip, entry.GetSize()))
				_ = entryStream
				j2gort.MapPutStr(mapFromExternsZip, filename, j2gort.PtrToAny(SourceFile_S_Builder().WithPath_1_string(("externs.zip//" + j2gort.StrValString(filename))).WithContent_1_Ext_java_io_InputStream(entryStream).Build()))
			}
		}
		return DefaultExterns_PrepareExterns(env, mapFromExternsZip)
	}
	panic("unreachable")
}

// Java method: AbstractCommandLineRunner.getVersionText()
func (this_ *AbstractCommandLineRunnerBase) GetVersionText() string {
	panic("abstract or native method invoked")
}

// Java static method: AbstractCommandLineRunner.getExternsInput()
func AbstractCommandLineRunner_GetExternsInput() j2gort.Ext_java_io_InputStream {
	clinit_AbstractCommandLineRunner()
	var input j2gort.Ext_java_io_InputStream = func() j2gort.Ext_java_io_InputStream {
		if b_, ok_ := j2gort.ClassResourceBytes(j2gort.ClassLiteral("com.google.javascript.jscomp.AbstractCommandLineRunner", "AbstractCommandLineRunner"), "/externs.zip"); ok_ {
			s_ := &j2gort.Ext_java_io_InputStreamBase{Core_: b_}
			s_.Self_ = s_
			return s_
		}
		return nil
	}()
	_ = input
	if input == nil {
		input = func() j2gort.Ext_java_io_InputStream {
			if b_, ok_ := j2gort.ClassResourceBytes(j2gort.ClassLiteral("com.google.javascript.jscomp.AbstractCommandLineRunner", "AbstractCommandLineRunner"), "externs.zip"); ok_ {
				s_ := &j2gort.Ext_java_io_InputStreamBase{Core_: b_}
				s_.Self_ = s_
				return s_
			}
			return nil
		}()
	}
	j2gort.CheckNotNull(input)
	return input
}

// Java method: AbstractCommandLineRunner.run()
func (this_ *AbstractCommandLineRunnerBase) Run() {
	var result int32 = int32(0)
	_ = result
	func() {
		defer func() {
			if rec_19 := recover(); rec_19 != nil {
				jex_20 := javalang.RecoverException(rec_19)
				_ = jex_20
				if j2gort.InstanceOf[*FlagUsageException](jex_20) {
					e := j2gort.AsIface[*FlagUsageException](jex_20)
					_ = e
					this_.err.Println_1_string(j2gort.ThrowableGetMessage(e))
					result = int32(-1)
				} else if j2gort.InstanceOf[javalang.Throwable](jex_20) {
					t := j2gort.AsIface[javalang.Throwable](jex_20)
					_ = t
					j2gort.ThrowablePrintStackTrace(t, this_.err)
					result = int32(-2)
				} else {
					panic(rec_19)
				}
			}
		}()
		result = this_.DoRun()
	}()
	this_.exitCodeReceiver.Apply(j2gort.ValueOfInt(result))
}

// Java method: AbstractCommandLineRunner.getErrorPrintStream()
func (this_ *AbstractCommandLineRunnerBase) GetErrorPrintStream() *j2gort.Ext_java_io_PrintStream {
	return this_.err
}

// Java method: AbstractCommandLineRunner.parseJsonFilesFromInputStream()
func (this_ *AbstractCommandLineRunnerBase) ParseJsonFilesFromInputStream() j2gort.Ext_java_util_List {
	var jsonFiles []*AbstractCommandLineRunner_JsonFileSpec = []*AbstractCommandLineRunner_JsonFileSpec{}
	_ = jsonFiles
	func() {
		var reader *j2gort.Ext_com_google_gson_stream_JsonReader = j2gort.NewExt_com_google_gson_stream_JsonReader(j2gort.NewExt_java_io_InputStreamReader_2_Ext_java_io_InputStream_Ext_java_nio_charset_Charset(this_.in, this_.inputCharset))
		_ = reader
		defer j2gort.CloseResource(reader)
		reader.BeginArray()
		for reader.HasNext() {
			var jsonFile *AbstractCommandLineRunner_JsonFileSpec = AbstractCommandLineRunner_JsonFileSpec_ReadJson(reader)
			_ = jsonFile
			jsonFiles = append(jsonFiles, jsonFile)
		}
		reader.EndArray()
	}()
	return j2gort.WrapTypedArrayList(jsonFiles)
}

// Java method: AbstractCommandLineRunner.createInputs(java.util.List<com.google.javascript.jscomp.AbstractCommandLineRunner.FlagEntry<com.google.javascript.jscomp.AbstractCommandLineRunner.JsSourceType>> files, boolean allowStdIn, java.util.List<com.google.javascript.jscomp.AbstractCommandLineRunner.JsChunkSpec> jsChunkSpecs)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_createInputs_152(files j2gort.Ext_java_util_List, allowStdIn bool, jsChunkSpecs j2gort.Ext_java_util_List) j2gort.Ext_java_util_List {
	_ = files
	_ = allowStdIn
	_ = jsChunkSpecs
	return this_.CreateInputs_4(files, nil, allowStdIn, jsChunkSpecs)
}

// Java method: AbstractCommandLineRunner.createInputs(java.util.List<com.google.javascript.jscomp.AbstractCommandLineRunner.FlagEntry<com.google.javascript.jscomp.AbstractCommandLineRunner.JsSourceType>> files, java.util.List<com.google.javascript.jscomp.AbstractCommandLineRunner.JsonFileSpec> jsonFiles, java.util.List<com.google.javascript.jscomp.AbstractCommandLineRunner.JsChunkSpec> jsChunkSpecs)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_createInputs_153(files j2gort.Ext_java_util_List, jsonFiles j2gort.Ext_java_util_List, jsChunkSpecs j2gort.Ext_java_util_List) j2gort.Ext_java_util_List {
	_ = files
	_ = jsonFiles
	_ = jsChunkSpecs
	return this_.CreateInputs_4(files, jsonFiles, false, jsChunkSpecs)
}

// Java method: AbstractCommandLineRunner.createInputs(java.util.List<com.google.javascript.jscomp.AbstractCommandLineRunner.FlagEntry<com.google.javascript.jscomp.AbstractCommandLineRunner.JsSourceType>> files, java.util.@org.jspecify.annotations.Nullable List<com.google.javascript.jscomp.AbstractCommandLineRunner.JsonFileSpec> jsonFiles, boolean allowStdIn, java.util.List<com.google.javascript.jscomp.AbstractCommandLineRunner.JsChunkSpec> jsChunkSpecs)
func (this_ *AbstractCommandLineRunnerBase) CreateInputs_4(files j2gort.Ext_java_util_List, jsonFiles j2gort.Ext_java_util_List, allowStdIn bool, jsChunkSpecs j2gort.Ext_java_util_List) j2gort.Ext_java_util_List {
	_ = files
	_ = jsonFiles
	_ = allowStdIn
	_ = jsChunkSpecs
	var inputs []*SourceFile = j2gort.MakeNonNilSlice[*SourceFile](files.Size())
	_ = inputs
	var usingStdin bool = false
	_ = usingStdin
	var jsChunkIndex int32 = int32(0)
	_ = jsChunkIndex
	var jsChunkSpec *AbstractCommandLineRunner_JsChunkSpec = j2gort.CastPtr[AbstractCommandLineRunner_JsChunkSpec](j2gort.Ext_com_google_common_collect_Iterables_GetFirst(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](jsChunkSpecs), nil))
	_ = jsChunkSpec
	var cumulatedInputFilesExpected int32
	_ = cumulatedInputFilesExpected
	if jsChunkSpec == nil {
		cumulatedInputFilesExpected = int32(2147483647)
	} else {
		cumulatedInputFilesExpected = jsChunkSpec.GetNumInputs()
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < files.Size(); i = i + 1 {
			var file *AbstractCommandLineRunner_FlagEntry = j2gort.CastPtr[AbstractCommandLineRunner_FlagEntry](files.Get(i))
			_ = file
			var filename string = file.value
			_ = filename
			if j2gort.Cast[AbstractCommandLineRunner_JsSourceType](file.flag) == StaticGet_AbstractCommandLineRunner_JsSourceType_JS_ZIP() {
				if this_.config.typedAstListInputFilename != j2gort.NullString {
					panic(NewFlagUsageException("Can't use TypedASTs with --zip."))
				}
				if !("-" == filename) {
					var newFiles j2gort.Ext_java_util_List = SourceFile_FromZipFile(filename, this_.inputCharset)
					_ = newFiles
					if j2gort.MapContainsKeyStr(this_.rootRelativePathsMap, filename) {
						var rootFilename string = j2gort.CastString(j2gort.MapGetStr(this_.rootRelativePathsMap, filename))
						_ = rootFilename
						{
							cur_21 := j2gort.NewIterCursor(j2gort.CheckNotNull(newFiles))
							for cur_21.HasNext() {
								rawElem_22 := cur_21.Next()
								zipEntry := j2gort.CastPtr[SourceFile](rawElem_22)
								_ = zipEntry
								var zipEntryName string = zipEntry.GetName()
								_ = zipEntryName
								if !(j2gort.JStrIndexOf(zipEntryName, j2gort.CheckNotNullString(filename)) >= 0) {
									j2gort.CheckState(false)
								}
								var zipmap string = javalang.StringReplace(zipEntryName, j2gort.CheckNotNullString(filename), j2gort.CheckNotNullString(rootFilename))
								_ = zipmap
								j2gort.MapPutStr(this_.rootRelativePathsMap, zipEntryName, j2gort.BoxString(zipmap))
							}
						}
					}
					inputs = j2gort.AppendCollection[*SourceFile](inputs, j2gort.AsIface[j2gort.Ext_java_util_Collection](newFiles))
					if jsChunkSpec != nil {
						r_23 := jsChunkSpec
						r_23.numJsFiles = (r_23.numJsFiles + (newFiles.Size() - int32(1)))
					}
				}
			} else if !("-" == filename) {
				var kind pkg_rhino.StaticSourceFile_SourceKind
				_ = kind
				if j2gort.Cast[AbstractCommandLineRunner_JsSourceType](file.flag) == StaticGet_AbstractCommandLineRunner_JsSourceType_WEAKDEP() {
					kind = pkg_rhino.StaticGet_StaticSourceFile_SourceKind_WEAK()
				} else {
					kind = pkg_rhino.StaticGet_StaticSourceFile_SourceKind_STRONG()
				}
				var newFile *SourceFile = SourceFile_S_Builder().WithPath_1_string(filename).WithCharset(this_.inputCharset).WithKind(kind).Build()
				_ = newFile
				inputs = append(inputs, newFile)
			} else {
				if this_.config.typedAstListInputFilename != j2gort.NullString {
					panic(NewFlagUsageException("Can't use TypedASTs with stdin."))
				}
				if !(this_.config.defaultToStdin) {
					panic(NewFlagUsageException("Can't specify stdin."))
				}
				if usingStdin {
					panic(NewFlagUsageException("Can't specify stdin twice."))
				}
				if !(this_.config.outputManifests.IsEmpty()) {
					panic(NewFlagUsageException("Manifest files cannot be generated when the input is from stdin."))
				}
				if !(this_.config.outputBundles.IsEmpty()) {
					panic(NewFlagUsageException("Bundle files cannot be generated when the input is from stdin."))
				}
				this_.err.Println_1_string("The compiler is waiting for input via stdin.")
				inputs = append(inputs, SourceFile_S_Builder().WithPath_1_string("stdin").WithContent_1_Ext_java_io_InputStream(this_.in).WithCharset(this_.inputCharset).Build())
				usingStdin = true
			}
			if i >= (cumulatedInputFilesExpected - int32(1)) {
				jsChunkIndex = jsChunkIndex + 1
				if jsChunkIndex < jsChunkSpecs.Size() {
					jsChunkSpec = j2gort.CastPtr[AbstractCommandLineRunner_JsChunkSpec](jsChunkSpecs.Get(jsChunkIndex))
					cumulatedInputFilesExpected = (cumulatedInputFilesExpected + jsChunkSpec.GetNumInputs())
				}
			}
		}
	}
	if jsonFiles != nil {
		{
			cur_24 := j2gort.NewIterCursor(j2gort.CheckNotNull(jsonFiles))
			for cur_24.HasNext() {
				rawElem_25 := cur_24.Next()
				jsonFile := j2gort.CastPtr[AbstractCommandLineRunner_JsonFileSpec](rawElem_25)
				_ = jsonFile
				inputs = append(inputs, SourceFile_FromCode_2(jsonFile.GetPath(), jsonFile.GetSrc()))
			}
		}
	}
	return j2gort.WrapTypedArrayList(inputs)
}

// Java method: AbstractCommandLineRunner.deduplicateIjsFiles(java.util.List<com.google.javascript.jscomp.AbstractCommandLineRunner.FlagEntry<com.google.javascript.jscomp.AbstractCommandLineRunner.JsSourceType>> files, java.util.List<java.lang.String> moduleRoots, boolean hasChunkSpecs)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_deduplicateIjsFiles_155(files j2gort.Ext_java_util_List, moduleRoots j2gort.Ext_java_util_List, hasChunkSpecs bool) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = files
	_ = moduleRoots
	_ = hasChunkSpecs
	var errors []*JSError = []*JSError{}
	_ = errors
	var relativeToAbsoluteName j2gort.OrderedMap[string, string] = j2gort.NewOrderedMap[string, string]()
	_ = relativeToAbsoluteName
	{
		cur_26 := j2gort.NewIterCursor(j2gort.CheckNotNull(files))
		for cur_26.HasNext() {
			rawElem_27 := cur_26.Next()
			file := j2gort.CastPtr[AbstractCommandLineRunner_FlagEntry](rawElem_27)
			_ = file
			if (j2gort.Cast[AbstractCommandLineRunner_JsSourceType](file.flag) == StaticGet_AbstractCommandLineRunner_JsSourceType_JS()) || (j2gort.Cast[AbstractCommandLineRunner_JsSourceType](file.flag) == StaticGet_AbstractCommandLineRunner_JsSourceType_WEAKDEP()) {
				var absoluteName string = file.value
				_ = absoluteName
				var relativeName string = this_.priv_AbstractCommandLineRunner_getModuleRootRelativeName_156(absoluteName, moduleRoots)
				_ = relativeName
				relativeToAbsoluteName.Put(relativeName, absoluteName)
			}
		}
	}
	var iterator j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](files.Iterator())
	_ = iterator
	for iterator.HasNext() {
		var file_2 *AbstractCommandLineRunner_FlagEntry = j2gort.CastPtr[AbstractCommandLineRunner_FlagEntry](iterator.Next())
		_ = file_2
		if j2gort.Cast[AbstractCommandLineRunner_JsSourceType](file_2.flag) == StaticGet_AbstractCommandLineRunner_JsSourceType_IJS() {
			if hasChunkSpecs {
				panic(NewFlagUsageException("--ijs is incompatible with --chunk or --module."))
			}
			var absoluteName_2 string = file_2.value
			_ = absoluteName_2
			if !(javalang.StringEndsWith(j2gort.CheckNotNullString(absoluteName_2), ".i.js")) {
				errors = append(errors, j2gort.CheckNotNullElem(JSError_Make_2(StaticGet_IjsErrors_BAD_IJS_FILE_NAME(), []string{absoluteName_2})))
				continue
			}
			var relativeName_2 string = this_.priv_AbstractCommandLineRunner_getModuleRootRelativeName_156(absoluteName_2, moduleRoots)
			_ = relativeName_2
			var relativeNonIjsName string = j2gort.JStrSub(relativeName_2, int32(0), (j2gort.JStrLen(relativeName_2) - j2gort.JStrLen(".i.js")))
			_ = relativeNonIjsName
			if relativeToAbsoluteName.ContainsKey(relativeNonIjsName) {
				errors = append(errors, j2gort.CheckNotNullElem(JSError_Make_2(StaticGet_IjsErrors_CONFLICTING_IJS_FILE(), []string{j2gort.OrderedMapGetStr(relativeToAbsoluteName, relativeNonIjsName), absoluteName_2})))
				iterator.Remove()
			}
		}
	}
	return j2gort.ImmListFromTypedSlice(errors)
}

// Java method: AbstractCommandLineRunner.getModuleRootRelativeName(java.lang.String filename, java.util.List<java.lang.String> moduleRoots)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_getModuleRootRelativeName_156(filename string, moduleRoots j2gort.Ext_java_util_List) string {
	_ = filename
	_ = moduleRoots
	{
		cur_28 := j2gort.NewIterCursor(j2gort.CheckNotNull(moduleRoots))
		for cur_28.HasNext() {
			moduleRoot := cur_28.NextStr()
			_ = moduleRoot
			if javalang.StringStartsWith(j2gort.CheckNotNullString(filename), j2gort.StrValString(moduleRoot)+"/") {
				return j2gort.JStrSubFrom(filename, (j2gort.JStrLen(moduleRoot) + int32(1)))
			}
		}
	}
	return filename
}

// Java method: AbstractCommandLineRunner.createSourceInputs(java.util.List<com.google.javascript.jscomp.AbstractCommandLineRunner.JsChunkSpec> jsChunkSpecs, java.util.List<com.google.javascript.jscomp.AbstractCommandLineRunner.FlagEntry<com.google.javascript.jscomp.AbstractCommandLineRunner.JsSourceType>> files, java.util.List<com.google.javascript.jscomp.AbstractCommandLineRunner.JsonFileSpec> jsonFiles, java.util.List<java.lang.String> moduleRoots)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_createSourceInputs_157(jsChunkSpecs j2gort.Ext_java_util_List, files j2gort.Ext_java_util_List, jsonFiles j2gort.Ext_java_util_List, moduleRoots j2gort.Ext_java_util_List) j2gort.Ext_java_util_List {
	_ = jsChunkSpecs
	_ = files
	_ = jsonFiles
	_ = moduleRoots
	if this_.IsInTestMode() {
		if this_.inputsSupplierForTesting != nil {
			return j2gort.AsIface[j2gort.Ext_java_util_List](this_.inputsSupplierForTesting.Get())
		} else {
			return nil
		}
	}
	if (files.IsEmpty() && (jsonFiles == nil)) && this_.config.defaultToStdin {
		files = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(NewAbstractCommandLineRunner_FlagEntry(j2gort.BoxNullableVal(StaticGet_AbstractCommandLineRunner_JsSourceType_JS()), "-")))
	}
	{
		cur_29 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.priv_AbstractCommandLineRunner_deduplicateIjsFiles_155(files, moduleRoots, !(jsChunkSpecs.IsEmpty()))))
		for cur_29.HasNext() {
			rawElem_30 := cur_29.Next()
			error_ := j2gort.CastPtr[JSError](rawElem_30)
			_ = error_
			this_.compiler.Report_1(error_)
		}
	}
	var didRet_31 bool
	var retVal_32 j2gort.Ext_java_util_List
	func() {
		defer func() {
			if rec_33 := recover(); rec_33 != nil {
				jex_34 := javalang.RecoverException(rec_33)
				_ = jex_34
				if j2gort.InstanceOf[*FlagUsageException](jex_34) {
					e := j2gort.AsIface[*FlagUsageException](jex_34)
					_ = e
					panic(NewFlagUsageException(("Bad --js flag. " + j2gort.StrValString(j2gort.ThrowableGetMessage(e)))))
				} else {
					panic(rec_33)
				}
			}
		}()
		if jsonFiles != nil {
			retVal_32 = this_.priv_AbstractCommandLineRunner_createInputs_153(files, jsonFiles, jsChunkSpecs)
			didRet_31 = true
			return
		} else {
			retVal_32 = this_.priv_AbstractCommandLineRunner_createInputs_152(files, true, jsChunkSpecs)
			didRet_31 = true
			return
		}
	}()
	if didRet_31 {
		return retVal_32
	}
	panic("unreachable")
}

// Java method: AbstractCommandLineRunner.createExternInputs(java.util.List<java.lang.String> files)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_createExternInputs_158(files *[]string) j2gort.Ext_java_util_List {
	_ = files
	var externFiles j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = externFiles
	for _, file := range *files {
		_ = file
		externFiles.Add_1(NewAbstractCommandLineRunner_FlagEntry(j2gort.BoxNullableVal(StaticGet_AbstractCommandLineRunner_JsSourceType_EXTERN()), file))
	}

	var didRet_36 bool
	var retVal_37 j2gort.Ext_java_util_List
	func() {
		defer func() {
			if rec_38 := recover(); rec_38 != nil {
				jex_39 := javalang.RecoverException(rec_38)
				_ = jex_39
				if j2gort.InstanceOf[*FlagUsageException](jex_39) {
					e := j2gort.AsIface[*FlagUsageException](jex_39)
					_ = e
					panic(NewFlagUsageException(("Bad --externs flag. " + j2gort.StrValString(j2gort.ThrowableGetMessage(e)))))
				} else {
					panic(rec_38)
				}
			}
		}()
		retVal_37 = this_.priv_AbstractCommandLineRunner_createInputs_152(externFiles, false, j2gort.NewExt_java_util_ArrayList())
		didRet_36 = true
		return
	}()
	if didRet_36 {
		return retVal_37
	}
	panic("unreachable")
}

// Java static method: AbstractCommandLineRunner.createJsChunks(java.util.List<com.google.javascript.jscomp.AbstractCommandLineRunner.JsChunkSpec> specs, java.util.List<com.google.javascript.jscomp.CompilerInput> inputs)
func AbstractCommandLineRunner_CreateJsChunks(specs j2gort.Ext_java_util_List, inputs j2gort.Ext_java_util_List) j2gort.Ext_java_util_List {
	clinit_AbstractCommandLineRunner()
	_ = specs
	_ = inputs
	if !(specs != nil) {
		j2gort.CheckState(false)
	}
	if !(!(specs.IsEmpty())) {
		j2gort.CheckState(false)
	}
	if !(inputs != nil) {
		j2gort.CheckState(false)
	}
	var chunkNames []string = j2gort.MakeNonNilSlice[string](specs.Size())
	_ = chunkNames
	var chunksByName j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = chunksByName
	var chunksFileCountMap j2gort.OrderedMap[string, int32] = j2gort.NewOrderedMap[string, int32]()
	_ = chunksFileCountMap
	var numJsFilesExpected int32 = int32(0)
	_ = numJsFilesExpected
	var minJsFilesRequired int32 = int32(0)
	_ = minJsFilesRequired
	{
		cur_40 := j2gort.NewIterCursor(j2gort.CheckNotNull(specs))
		for cur_40.HasNext() {
			rawElem_41 := cur_40.Next()
			spec := j2gort.CastPtr[AbstractCommandLineRunner_JsChunkSpec](rawElem_41)
			_ = spec
			if j2gort.MapContainsKeyStr(chunksByName, spec.name) {
				panic(NewFlagUsageException(("Duplicate chunk name: " + j2gort.StrValString(spec.name))))
			}
			var chunk *JSChunk = NewJSChunk(spec.name)
			_ = chunk
			{

				for _, dep := range spec.deps {

					_ = dep
					var other *JSChunk = j2gort.CastPtr[JSChunk](j2gort.MapGetStr(chunksByName, dep))
					_ = other
					if other == nil {
						panic(NewFlagUsageException(("Chunk '" + j2gort.StrValString(spec.name) + "' depends on unknown chunk '" + j2gort.StrValString(dep) + "'. Be sure to list chunks in dependency order.")))
					}
					chunk.AddDependency(other)
				}
			}
			if spec.numJsFiles < int32(0) {
				numJsFilesExpected = int32(-1)
			} else {
				minJsFilesRequired = (minJsFilesRequired + spec.numJsFiles)
			}
			if numJsFilesExpected >= int32(0) {
				numJsFilesExpected = (numJsFilesExpected + spec.numJsFiles)
			}
			chunkNames = j2gort.SliceInsertAt(chunkNames, int32(0), spec.name)
			chunksFileCountMap.Put(spec.name, spec.numJsFiles)
			j2gort.MapPutStr(chunksByName, spec.name, chunk)
		}
	}
	var totalNumJsFiles int32 = inputs.Size()
	_ = totalNumJsFiles
	if (numJsFilesExpected >= int32(0)) || (minJsFilesRequired > totalNumJsFiles) {
		if minJsFilesRequired > totalNumJsFiles {
			numJsFilesExpected = minJsFilesRequired
		}
		if numJsFilesExpected > totalNumJsFiles {
			panic(NewFlagUsageException(("Not enough JS files specified. Expected " + j2gort.StrValInt32(numJsFilesExpected) + " but found " + j2gort.StrValInt32(totalNumJsFiles))))
		} else if numJsFilesExpected < totalNumJsFiles {
			panic(NewFlagUsageException(("Too many JS files specified. Expected " + j2gort.StrValInt32(numJsFilesExpected) + " but found " + j2gort.StrValInt32(totalNumJsFiles))))
		}
	}
	var numJsFilesLeft int32 = totalNumJsFiles
	_ = numJsFilesLeft
	var chunkIndex int32 = int32(0)
	_ = chunkIndex
	for _, chunkName := range chunkNames {
		_ = chunkName
		var numJsFiles int32 = func() int32 {
			if v_, ok_ := chunksFileCountMap.GetOk(chunkName); ok_ {
				return v_
			}
			j2gort.ThrowNPE()
			return 0
		}()
		_ = numJsFiles
		var chunk_2 *JSChunk = j2gort.CastPtr[JSChunk](j2gort.MapGetStr(chunksByName, chunkName))
		_ = chunk_2
		if (chunkIndex == (int32(len(chunkNames)) - int32(1))) && (numJsFiles == int32(-1)) {
			numJsFiles = numJsFilesLeft
		}
		var chunkFiles j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](inputs.SubList((numJsFilesLeft - numJsFiles), numJsFilesLeft))
		_ = chunkFiles
		{
			cur_44 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunkFiles))
			for cur_44.HasNext() {
				rawElem_45 := cur_44.Next()
				input := j2gort.CastPtr[CompilerInput](rawElem_45)
				_ = input
				chunk_2.Add_1_CompilerInput(input)
			}
		}
		numJsFilesLeft = (numJsFilesLeft - numJsFiles)
		chunkIndex = chunkIndex + 1
	}

	return j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](chunksByName.Values()))
}

// Java method: AbstractCommandLineRunner.checkChunkName(java.lang.String name)
func (this_ *AbstractCommandLineRunnerBase) CheckChunkName(name string) {
	_ = name
	if !(pkg_rhino.TokenStream_IsJSIdentifier(name)) {
		panic(NewFlagUsageException(("Invalid chunk name: '" + j2gort.StrValString(name) + "'")))
	}
}

// Java static method: AbstractCommandLineRunner.parseChunkWrappers(java.util.List<java.lang.String> specs, java.lang.Iterable<com.google.javascript.jscomp.JSChunk> chunks)
func AbstractCommandLineRunner_ParseChunkWrappers(specs j2gort.Ext_java_util_List, chunks j2gort.Ext_java_lang_Iterable) j2gort.Ext_java_util_Map {
	clinit_AbstractCommandLineRunner()
	_ = specs
	_ = chunks
	if !(specs != nil) {
		j2gort.CheckState(false)
	}
	var wrappers j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = wrappers
	{
		cur_46 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunks))
		for cur_46.HasNext() {
			rawElem_47 := cur_46.Next()
			c := j2gort.CastPtr[JSChunk](rawElem_47)
			_ = c
			j2gort.MapPutStr(wrappers, c.GetName(), "")
		}
	}
	{
		cur_48 := j2gort.NewIterCursor(j2gort.CheckNotNull(specs))
		for cur_48.HasNext() {
			spec := cur_48.NextStr()
			_ = spec
			var pos int32 = j2gort.JStrIndexOfCp(spec, int32(uint16(58)))
			_ = pos
			if pos == int32(-1) {
				panic(NewFlagUsageException(("Expected chunk wrapper to have <name>:<wrapper> format: " + j2gort.StrValString(spec))))
			}
			var name string = j2gort.JStrSub(spec, int32(0), pos)
			_ = name
			if !(j2gort.MapContainsKeyStr(wrappers, name)) {
				panic(NewFlagUsageException(("Unknown chunk: '" + j2gort.StrValString(name) + "'")))
			}
			var wrapper string = j2gort.JStrSubFrom(spec, (pos + int32(1)))
			_ = wrapper
			wrapper = javalang.StringReplace(javalang.StringReplace(wrapper, "%output%", "%s"), "%n%", "\n")
			if !(j2gort.JStrIndexOf(wrapper, "%s") >= 0) {
				panic(NewFlagUsageException(("No %s placeholder in chunk wrapper: '" + j2gort.StrValString(wrapper) + "'")))
			}
			j2gort.MapPutStr(wrappers, name, j2gort.BoxString(wrapper))
		}
	}
	return wrappers
}

// Java static method: AbstractCommandLineRunner.parseChunkOutputFiles(java.util.List<java.lang.String> specs)
func AbstractCommandLineRunner_ParseChunkOutputFiles(specs j2gort.Ext_java_util_List) j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_AbstractCommandLineRunner()
	_ = specs
	if !(specs != nil) {
		j2gort.CheckArgument(false)
	}
	var outputFilesBuilder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = outputFilesBuilder
	{
		cur_49 := j2gort.NewIterCursor(j2gort.CheckNotNull(specs))
		for cur_49.HasNext() {
			spec := cur_49.NextStr()
			_ = spec
			var pos int32 = j2gort.JStrIndexOfCp(spec, int32(uint16(58)))
			_ = pos
			if pos == int32(-1) {
				panic(NewFlagUsageException(("Expected chunk_output_file to have <name>:<output_file> format: " + j2gort.StrValString(spec))))
			}
			var name string = j2gort.JStrSub(spec, int32(0), pos)
			_ = name
			var filename string = j2gort.JStrSubFrom(spec, (pos + int32(1)))
			_ = filename
			outputFilesBuilder.Put_2(j2gort.BoxString(name), j2gort.BoxString(filename))
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](outputFilesBuilder.BuildOrThrow())
}

// Java method: AbstractCommandLineRunner.getChunkOutputFileName(com.google.javascript.jscomp.JSChunk m)
func (this_ *AbstractCommandLineRunnerBase) GetChunkOutputFileName(m *JSChunk) string {
	_ = m
	if this_.parsedChunkOutputFiles == nil {
		this_.parsedChunkOutputFiles = AbstractCommandLineRunner_ParseChunkOutputFiles(this_.config.chunkOutputFiles)
	}
	var outputFile string = j2gort.CastString(j2gort.MapGetStr(this_.parsedChunkOutputFiles, m.GetName()))
	_ = outputFile
	if outputFile != j2gort.NullString {
		return (j2gort.StrValString(this_.config.chunkOutputPathPrefix) + j2gort.StrValString(outputFile))
	}
	return (j2gort.StrValString(this_.config.chunkOutputPathPrefix) + j2gort.StrValString(m.GetName()) + ".js")
}

// Java method: AbstractCommandLineRunner.getChunkConformanceFileName(com.google.javascript.jscomp.JSChunk m)
func (this_ *AbstractCommandLineRunnerBase) GetChunkConformanceFileName(m *JSChunk) string {
	_ = m
	if this_.parsedChunkConformanceFiles == nil {
		this_.parsedChunkConformanceFiles = AbstractCommandLineRunner_ParseChunkOutputFiles(this_.config.chunkConformanceFiles)
	}
	return j2gort.CastString(j2gort.MapGetStr(this_.parsedChunkConformanceFiles, m.GetName()))
}

// Java method: AbstractCommandLineRunner.writeChunkOutput(java.lang.String fileName, java.lang.Appendable out, com.google.javascript.jscomp.CodePrinter.LicenseTracker lt, com.google.javascript.jscomp.JSChunk m)
func (this_ *AbstractCommandLineRunnerBase) WriteChunkOutput(fileName string, out j2gort.Ext_java_lang_Appendable, lt CodePrinter_LicenseTracker, m *JSChunk) {
	_ = fileName
	_ = out
	_ = lt
	_ = m
	if this_.parsedChunkWrappers == nil {
		this_.parsedChunkWrappers = AbstractCommandLineRunner_ParseChunkWrappers(this_.config.chunkWrapper, j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable(this_.compiler.GetChunkGraph().GetAllChunks())))
	}
	if !(this_.priv_AbstractCommandLineRunner_isOutputInJson_138()) {
		AbstractCommandLineRunner_MaybeCreateDirsForPath(fileName)
	}
	var baseName string = j2gort.NewExt_java_io_File_1_string(fileName).GetName()
	_ = baseName
	this_.WriteOutput_8(out, this_.compiler, lt, m, javalang.StringReplace(j2gort.CastString(j2gort.MapGetStr(this_.parsedChunkWrappers, m.GetName())), "%basename%", j2gort.CheckNotNullString(baseName)), "%s", nil, fileName)
}

// Java method: AbstractCommandLineRunner.writeOutput(java.lang.Appendable out, com.google.javascript.jscomp.Compiler compiler, com.google.javascript.jscomp.CodePrinter.LicenseTracker licenseTracker, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable JSChunk chunk, java.lang.String wrapper, java.lang.String codePlaceholder, com.google.common.base.@org.jspecify.annotations.Nullable Function<java.lang.String,java.lang.String> escaper, java.lang.String filename)
func (this_ *AbstractCommandLineRunnerBase) WriteOutput_8(out j2gort.Ext_java_lang_Appendable, compiler *Compiler, licenseTracker CodePrinter_LicenseTracker, chunk *JSChunk, wrapper string, codePlaceholder string, escaper j2gort.Ext_com_google_common_base_Function, filename string) {
	_ = out
	_ = compiler
	_ = licenseTracker
	_ = chunk
	_ = wrapper
	_ = codePlaceholder
	_ = escaper
	_ = filename
	if compiler.GetOptions().GetOutputJs() == StaticGet_CompilerOptions_OutputJs_SENTINEL() {
		out.Append_1_any("// No JS output because the compiler was run in checks-only mode.\n")
		return
	}
	if !(compiler.GetOptions().GetOutputJs() == StaticGet_CompilerOptions_OutputJs_NORMAL()) {
		j2gort.CheckState(false)
	}
	var code string
	_ = code
	if chunk == nil {
		code = compiler.ToSource_0()
	} else {
		code = compiler.ToSource_2(licenseTracker, chunk)
	}
	this_.WriteOutput_7(out, compiler, code, wrapper, codePlaceholder, escaper, filename)
}

// Java method: AbstractCommandLineRunner.writeOutput(java.lang.Appendable out, com.google.javascript.jscomp.Compiler compiler, java.lang.String code, java.lang.String wrapper, java.lang.String codePlaceholder, com.google.common.base.@org.jspecify.annotations.Nullable Function<java.lang.String,java.lang.String> escaper, java.lang.String filename)
func (this_ *AbstractCommandLineRunnerBase) WriteOutput_7(out j2gort.Ext_java_lang_Appendable, compiler *Compiler, code string, wrapper string, codePlaceholder string, escaper j2gort.Ext_com_google_common_base_Function, filename string) {
	_ = out
	_ = compiler
	_ = code
	_ = wrapper
	_ = codePlaceholder
	_ = escaper
	_ = filename
	var pos int32 = j2gort.JStrIndexOf(wrapper, codePlaceholder)
	_ = pos
	if pos != int32(-1) {
		var prefix string = ""
		_ = prefix
		if pos > int32(0) {
			prefix = j2gort.JStrSub(wrapper, int32(0), pos)
			out.Append_1_any(j2gort.BoxString(prefix))
		}
		out.Append_1_any(func() any {
			if escaper == nil {
				return j2gort.BoxString(code)
			}
			return j2gort.BoxString(j2gort.CastString(escaper.Apply(j2gort.BoxString(code))))
		}())
		var suffixStart int32 = (pos + j2gort.JStrLen(codePlaceholder))
		_ = suffixStart
		if suffixStart != j2gort.JStrLen(wrapper) {
			out.Append_1_any(j2gort.BoxString(j2gort.JStrSubFrom(wrapper, suffixStart)))
		}
		if this_.GetCommandLineConfig().includeTrailingNewline {
			out.Append_1_uint16(uint16(10))
		}
		if (compiler != nil) && (compiler.GetSourceMap() != nil) {
			compiler.GetSourceMap().SetWrapperPrefix(prefix)
		}
	} else {
		out.Append_1_any(j2gort.BoxString(code))
		if this_.GetCommandLineConfig().includeTrailingNewline {
			out.Append_1_uint16(uint16(10))
		}
	}
}

// Java static method: AbstractCommandLineRunner.maybeCreateDirsForPath(java.lang.String pathPrefix)
func AbstractCommandLineRunner_MaybeCreateDirsForPath(pathPrefix string) {
	clinit_AbstractCommandLineRunner()
	_ = pathPrefix
	if !(j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(pathPrefix)) {
		var parent *j2gort.Ext_java_io_File = j2gort.NewExt_java_io_File_1_string(pathPrefix).GetParentFile()
		_ = parent
		if parent != nil {
			parent.Mkdirs()
		}
	}
}

// Java method: AbstractCommandLineRunner.createDefaultOutput()
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_createDefaultOutput_165() j2gort.Ext_java_lang_Appendable {
	var writeOutputToFile bool = !(j2gort.CheckNotNullString(this_.config.jsOutputFile) == "")
	_ = writeOutputToFile
	if writeOutputToFile {
		return j2gort.AsIface[j2gort.Ext_java_lang_Appendable](this_.priv_AbstractCommandLineRunner_fileNameToLegacyOutputWriter_190(this_.config.jsOutputFile))
	} else {
		return j2gort.AsIface[j2gort.Ext_java_lang_Appendable](this_.priv_AbstractCommandLineRunner_streamToLegacyOutputWriter_193(j2gort.AsIface[j2gort.Ext_java_io_OutputStream](this_.defaultJsOutput)))
	}
}

// Java static method: AbstractCommandLineRunner.closeAppendable(java.lang.Appendable output)
func AbstractCommandLineRunner_CloseAppendable(output j2gort.Ext_java_lang_Appendable) {
	clinit_AbstractCommandLineRunner()
	_ = output
	var flushable j2gort.Ext_java_io_Flushable
	_ = flushable
	var ok_50 bool
	flushable, ok_50 = j2gort.TryCast[j2gort.Ext_java_io_Flushable](output)
	if ok_50 {
		flushable.Flush()
	}
	var closeable j2gort.Ext_java_io_Closeable
	_ = closeable
	var ok_51 bool
	closeable, ok_51 = j2gort.TryCast[j2gort.Ext_java_io_Closeable](output)
	if ok_51 {
		closeable.Close_0()
	}
}

// Java method: AbstractCommandLineRunner.doRun()
func (this_ *AbstractCommandLineRunnerBase) DoRun() int32 {
	var metricsRecorder CompileMetricsRecorderInterface = this_.GetCompileMetricsRecorder()
	_ = metricsRecorder
	metricsRecorder.RecordActionStart()
	if this_.config.printVersion {
		this_.defaultJsOutput.Println_1_string(this_.self_.GetVersionText())
		this_.defaultJsOutput.Flush()
		return int32(0)
	}
	Compiler_SetLoggingLevel(j2gort.Ext_java_util_logging_Level_Parse(this_.config.loggingLevel))
	this_.compiler = this_.self_.CreateCompiler()
	if !(this_.config.expectedDiagnostics.IsEmpty()) {
		this_.compiler.SetErrorManager(NewVerifyingErrorManager(this_.compiler.GetErrorManager(), this_.config.expectedDiagnostics))
	}
	var options *CompilerOptions = this_.self_.CreateOptions()
	_ = options
	this_.SetRunOptions(options)
	var typedAstListInputFilename string = this_.config.typedAstListInputFilename
	_ = typedAstListInputFilename
	var externs j2gort.Ext_java_util_List = this_.self_.CreateExterns(options)
	_ = externs
	var chunks j2gort.Ext_java_util_List = nil
	_ = chunks
	this_.rootRelativePathsMap = this_.priv_AbstractCommandLineRunner_constructRootRelativePathsMap_210()
	var writeOutputToFile bool = !(j2gort.CheckNotNullString(this_.config.jsOutputFile) == "")
	_ = writeOutputToFile
	var outputFileNames []string = []string{}
	_ = outputFileNames
	if writeOutputToFile {
		outputFileNames = append(outputFileNames, this_.config.jsOutputFile)
	}
	var createCommonJsModules bool = false
	_ = createCommonJsModules
	if options.GetProcessCommonJSModules() && ((int32(len(this_.config.chunk)) == int32(1)) && (this_.config.chunk[int32(0)] == "auto")) {
		createCommonJsModules = true
		this_.config.chunk = j2gort.SliceRemoveAt(this_.config.chunk, int32(0))
	}
	var jsChunkSpecs j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = jsChunkSpecs
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(this_.config.chunk)); i = i + 1 {
			jsChunkSpecs.Add_1(j2gort.PtrToAny(AbstractCommandLineRunner_JsChunkSpec_Create(this_.config.chunk[i], (i == int32(0)))))
		}
	}
	var jsonFiles j2gort.Ext_java_util_List = nil
	_ = jsonFiles
	if (this_.config.jsonStreamMode == StaticGet_CompilerOptions_JsonStreamMode_IN()) || (this_.config.jsonStreamMode == StaticGet_CompilerOptions_JsonStreamMode_BOTH()) {
		jsonFiles = this_.ParseJsonFilesFromInputStream()
		var inputSourceMaps j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.NewExt_com_google_common_collect_ImmutableMap_Builder()
		_ = inputSourceMaps
		var inputPathByWebpackId j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.NewExt_com_google_common_collect_ImmutableMap_Builder()
		_ = inputPathByWebpackId
		var foundJsonInputSourceMap bool = false
		_ = foundJsonInputSourceMap
		{
			cur_52 := j2gort.NewIterCursor(j2gort.CheckNotNull(jsonFiles))
			for cur_52.HasNext() {
				rawElem_53 := cur_52.Next()
				jsonFile := j2gort.CastPtr[AbstractCommandLineRunner_JsonFileSpec](rawElem_53)
				_ = jsonFile
				if (jsonFile.GetSourceMap() != j2gort.NullString) && (j2gort.JStrLen(jsonFile.GetSourceMap()) > int32(0)) {
					var sourceMapPath string = (j2gort.StrValString(jsonFile.GetPath()) + ".map")
					_ = sourceMapPath
					var sourceMap *SourceFile = SourceFile_FromCode_2(sourceMapPath, jsonFile.GetSourceMap())
					_ = sourceMap
					inputSourceMaps.Put_2(j2gort.BoxString(jsonFile.GetPath()), NewSourceMapInput(sourceMap))
					foundJsonInputSourceMap = true
				}
				if jsonFile.GetWebpackId() != j2gort.NullString {
					inputPathByWebpackId.Put_2(j2gort.BoxString(jsonFile.GetWebpackId()), j2gort.BoxString(jsonFile.GetPath()))
				}
			}
		}
		if foundJsonInputSourceMap {
			inputSourceMaps.PutAll_1_Ext_java_util_Map(j2gort.AsIface[j2gort.Ext_java_util_Map](options.GetInputSourceMaps()))
			options.SetInputSourceMaps(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](inputSourceMaps.BuildOrThrow()))
		}
		this_.compiler.InitWebpackMap(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](inputPathByWebpackId.BuildOrThrow()))
	} else {
		var emptyMap j2gort.Ext_com_google_common_collect_ImmutableMap = j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()
		_ = emptyMap
		this_.compiler.InitWebpackMap(emptyMap)
	}
	options.SetDoLateLocalization(this_.config.ShouldDoLateLocalization())
	options.SetAlwaysGatherSourceMapInfo(this_.config.ShouldAlwaysGatherSourceMapInfo())
	this_.compiler.InitOptions(options)
	var sources j2gort.Ext_java_util_List = this_.priv_AbstractCommandLineRunner_createSourceInputs_157(jsChunkSpecs, this_.config.mixedJsSources, jsonFiles, this_.config.moduleRoots)
	_ = sources
	if !(jsChunkSpecs.IsEmpty()) {
		if this_.IsInTestMode() {
			chunks = j2gort.AsIface[j2gort.Ext_java_util_List](this_.chunksSupplierForTesting.Get())
		} else {
			{
				var iifeCond_1 bool
				{

					re_54 := StaticGet_CompilerOptions_JsonStreamMode_IN()
					ae_55 := this_.config.jsonStreamMode
					if re_54 == 0 {
						j2gort.ThrowNPE()
					}
					iifeCond_1 = re_54 == ae_55
				}
				if iifeCond_1 || func() bool {
					re_56 := StaticGet_CompilerOptions_JsonStreamMode_NONE()
					ae_57 := this_.config.jsonStreamMode
					if re_56 == 0 {
						j2gort.ThrowNPE()
					}
					return re_56 == ae_57
				}() {
					{
						cur_58 := j2gort.NewIterCursor(j2gort.CheckNotNull(jsChunkSpecs))
						for cur_58.HasNext() {
							rawElem_59 := cur_58.Next()
							m := j2gort.CastPtr[AbstractCommandLineRunner_JsChunkSpec](rawElem_59)
							_ = m
							this_.self_.CheckChunkName(m.GetName())
						}
					}
				}
			}
			var inputs []*CompilerInput = []*CompilerInput{}
			_ = inputs
			{
				cur_60 := j2gort.NewIterCursor(j2gort.CheckNotNull(sources))
				for cur_60.HasNext() {
					rawElem_61 := cur_60.Next()
					source := j2gort.CastPtr[SourceFile](rawElem_61)
					_ = source
					inputs = append(inputs, NewCompilerInput_2_SourceFile_bool(source, false))
				}
			}
			chunks = AbstractCommandLineRunner_CreateJsChunks(jsChunkSpecs, j2gort.ImmListFromTypedSlice(inputs))
		}
		{
			cur_62 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunks))
			for cur_62.HasNext() {
				rawElem_63 := cur_62.Next()
				m_2 := j2gort.CastPtr[JSChunk](rawElem_63)
				_ = m_2
				outputFileNames = append(outputFileNames, this_.GetChunkOutputFileName(m_2))
			}
		}
		if typedAstListInputFilename != j2gort.NullString {
			this_.priv_AbstractCommandLineRunner_initChunksWithTypedAstFilesystem_171(externs, chunks, options, typedAstListInputFilename)
		} else {
			this_.compiler.InitChunks(externs, chunks, options)
		}
	} else {
		if typedAstListInputFilename != j2gort.NullString {
			this_.priv_AbstractCommandLineRunner_initWithTypedAstFilesystem_170(externs, sources, options, typedAstListInputFilename)
		} else {
			this_.compiler.Init(externs, sources, options)
		}
	}
	jsChunkSpecs = nil
	jsonFiles = nil
	externs = nil
	sources = nil
	var result *Result = nil
	_ = result
	if this_.config.skipNormalOutputs {
		metricsRecorder.RecordActionName("skip normal outputs")
		this_.compiler.OrderInputsWithLargeStack()
		result = nil
	} else if this_.compiler.HasErrors() {
		metricsRecorder.RecordActionName("initialization errors occurred")
		this_.compiler.GenerateReport()
		result = this_.compiler.GetResult()
	} else {
		func() {
			defer func() {
				this_.compiler.GenerateReport()
			}()
			this_.priv_AbstractCommandLineRunner_performCompilation_167(metricsRecorder)
			var refresh bool = (((chunks != nil) && (this_.config.restoredCompilationStage != int32(-1))) && (this_.config.saveAfterCompilationStage == int32(-1)))
			_ = refresh
			if refresh {
				chunks = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable(this_.compiler.GetChunks()))
			}
		}()
		result = this_.compiler.GetResult()
	}
	if createCommonJsModules {
		chunks = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable(this_.compiler.GetChunks()))
		{
			cur_64 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunks))
			for cur_64.HasNext() {
				rawElem_65 := cur_64.Next()
				c := j2gort.CastPtr[JSChunk](rawElem_65)
				_ = c
				outputFileNames = append(outputFileNames, this_.GetChunkOutputFileName(c))
			}
		}
	}
	for _, outputFileName := range outputFileNames {
		_ = outputFileName
		if this_.compiler.GetSourceFileByName(outputFileName) != nil {
			this_.compiler.Report_1(JSError_Make_2(StaticGet_AbstractCommandLineRunner_OUTPUT_SAME_AS_INPUT_ERROR(), []string{outputFileName}))
			return int32(1)
		}
	}

	var shouldProcessResults bool = (this_.config.GetSaveCompilationStateToFilename() == j2gort.NullString)
	_ = shouldProcessResults
	var exitStatus int32
	_ = exitStatus
	if shouldProcessResults {
		exitStatus = this_.ProcessResults(result, chunks, options)
	} else {
		exitStatus = AbstractCommandLineRunner_GetExitStatusForResult(result)
	}
	metricsRecorder.RecordResultMetrics(ToIface_AbstractCompiler(this_.compiler), result)
	return exitStatus
}

// Java method: AbstractCommandLineRunner.performCompilation(com.google.javascript.jscomp.CompileMetricsRecorderInterface metricsRecorder)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_performCompilation_167(metricsRecorder CompileMetricsRecorderInterface) {
	_ = metricsRecorder
	this_.priv_AbstractCommandLineRunner_initializeStateBeforeCompilation_172()
	this_.priv_AbstractCommandLineRunner_runCompilerPasses_168(metricsRecorder)
	if !(this_.compiler.HasErrors()) {
		this_.priv_AbstractCommandLineRunner_saveState_174()
	}
	this_.compiler.PerformPostCompilationTasks()
}

// Java method: AbstractCommandLineRunner.runCompilerPasses(com.google.javascript.jscomp.CompileMetricsRecorderInterface metricsRecorder)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_runCompilerPasses_168(metricsRecorder CompileMetricsRecorderInterface) {
	_ = metricsRecorder
	var runStage1 bool = false
	_ = runStage1
	var runStage2 bool = false
	_ = runStage2
	var runStage3 bool = false
	_ = runStage3
	var instrumentForCoverage bool = false
	_ = instrumentForCoverage
	var actionMetricsName string = j2gort.NullString
	_ = actionMetricsName
	if this_.compiler.GetOptions().GetInstrumentForCoverageOnly() {
		actionMetricsName = "instrument for coverage"
		instrumentForCoverage = true
	} else if this_.config.ShouldSaveAfterStage1() {
		actionMetricsName = "stage 1"
		runStage1 = true
	} else if this_.config.ShouldRestoreTypedAstsPerformStage2AndSave() {
		actionMetricsName = "parse & optimize"
		runStage2 = true
	} else if this_.config.ShouldRestoreTypedAstsPerformStages2And3() {
		actionMetricsName = "skip-checks compile"
		runStage2 = true
		runStage3 = true
	} else if this_.config.ShouldRestoreAndPerformStage2AndSave() {
		actionMetricsName = "stage 2/3"
		runStage2 = true
	} else if this_.config.ShouldRestoreAndPerformStages2And3() {
		actionMetricsName = "stage 2/2"
		runStage2 = true
		runStage3 = true
	} else if this_.config.ShouldRestoreAndPerformStage3() {
		actionMetricsName = "stage 3/3"
		runStage3 = true
	} else {
		if this_.compiler.GetOptions().IsChecksOnly() {
			actionMetricsName = "checks-only"
		} else {
			actionMetricsName = "full compile"
		}
		runStage1 = true
		runStage2 = true
		runStage3 = true
	}
	metricsRecorder.RecordActionName(actionMetricsName)
	if this_.compiler.HasErrors() {
		return
	}
	metricsRecorder.RecordStartState(ToIface_AbstractCompiler(this_.compiler))
	if runStage1 {
		this_.compiler.Stage1Passes()
		if this_.compiler.HasErrors() {
			return
		}
	}
	if runStage2 {
		this_.compiler.Stage2Passes(StaticGet_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS())
		if this_.compiler.HasErrors() {
			return
		}
	}
	if runStage3 {
		this_.compiler.Stage3Passes()
		if this_.compiler.HasErrors() {
			return
		}
	}
	if instrumentForCoverage {
		this_.compiler.InstrumentForCoverage()
	}
}

// Java method: AbstractCommandLineRunner.getCompileMetricsRecorder()
func (this_ *AbstractCommandLineRunnerBase) GetCompileMetricsRecorder() CompileMetricsRecorderInterface {
	return NewDummyCompileMetricsRecorder()
}

// Java method: AbstractCommandLineRunner.initWithTypedAstFilesystem(java.util.List<com.google.javascript.jscomp.SourceFile> externs, java.util.List<com.google.javascript.jscomp.SourceFile> sources, com.google.javascript.jscomp.CompilerOptions options, java.lang.String filename)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_initWithTypedAstFilesystem_170(externs j2gort.Ext_java_util_List, sources j2gort.Ext_java_util_List, options *CompilerOptions, filename string) {
	_ = externs
	_ = sources
	_ = options
	_ = filename
	func() {
		defer func() {
			if rec_67 := recover(); rec_67 != nil {
				jex_68 := javalang.RecoverException(rec_67)
				_ = jex_68
				if j2gort.InstanceOf[javaio.IIOException](jex_68) {
					e := j2gort.AsIface[javaio.IIOException](jex_68)
					_ = e
					this_.compiler.Report_1(JSError_Make_2(StaticGet_AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST(), []string{filename}))
				} else {
					panic(rec_67)
				}
			}
		}()
		var typedAstListStream *j2gort.Ext_java_util_zip_GZIPInputStream = j2gort.NewExt_java_util_zip_GZIPInputStream_2(j2gort.NewExt_java_io_FileInputStream_1_string(filename), int32(8388608))
		_ = typedAstListStream
		defer j2gort.CloseResource(typedAstListStream)
		this_.compiler.InitWithTypedAstFilesystem(externs, sources, options, typedAstListStream)
	}()
}

// Java method: AbstractCommandLineRunner.initChunksWithTypedAstFilesystem(java.util.List<com.google.javascript.jscomp.SourceFile> externs, java.util.List<com.google.javascript.jscomp.JSChunk> chunks, com.google.javascript.jscomp.CompilerOptions options, java.lang.String filename)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_initChunksWithTypedAstFilesystem_171(externs j2gort.Ext_java_util_List, chunks j2gort.Ext_java_util_List, options *CompilerOptions, filename string) {
	_ = externs
	_ = chunks
	_ = options
	_ = filename
	func() {
		defer func() {
			if rec_69 := recover(); rec_69 != nil {
				jex_70 := javalang.RecoverException(rec_69)
				_ = jex_70
				if j2gort.InstanceOf[javaio.IIOException](jex_70) {
					e := j2gort.AsIface[javaio.IIOException](jex_70)
					_ = e
					this_.compiler.Report_1(JSError_Make_2(StaticGet_AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST(), []string{filename}))
				} else {
					panic(rec_69)
				}
			}
		}()
		var typedAstListStream *j2gort.Ext_java_util_zip_GZIPInputStream = j2gort.NewExt_java_util_zip_GZIPInputStream_2(j2gort.NewExt_java_io_FileInputStream_1_string(filename), int32(8388608))
		_ = typedAstListStream
		defer j2gort.CloseResource(typedAstListStream)
		this_.compiler.InitChunksWithTypedAstFilesystem(externs, chunks, options, typedAstListStream)
	}()
}

// Java method: AbstractCommandLineRunner.initializeStateBeforeCompilation()
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_initializeStateBeforeCompilation_172() {
	if this_.config.restoredCompilationStage != int32(-1) {
		this_.priv_AbstractCommandLineRunner_restoreState_173(this_.config.GetContinueSavedCompilationFileName())
	} else if this_.config.typedAstListInputFilename != j2gort.NullString {
	} else {
		this_.compiler.ParseForCompilation()
	}
}

// Java method: AbstractCommandLineRunner.restoreState(java.lang.String filename)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_restoreState_173(filename string) {
	_ = filename
	func() {
		defer func() {
			if rec_71 := recover(); rec_71 != nil {
				jex_72 := javalang.RecoverException(rec_71)
				_ = jex_72
				if j2gort.InstanceOf[javaio.IIOException](jex_72) || j2gort.InstanceOf[j2gort.Ext_java_lang_ClassNotFoundException](jex_72) {
					e := jex_72
					_ = e
					this_.compiler.Report_1(JSError_Make_2(StaticGet_AbstractCommandLineRunner_COULD_NOT_DESERIALIZE_AST(), []string{filename}))
				} else {
					panic(rec_71)
				}
			}
		}()
		var serializedInputStream *j2gort.Ext_java_io_BufferedInputStream = j2gort.NewExt_java_io_BufferedInputStream_1(j2gort.NewExt_java_io_FileInputStream_1_string(filename))
		_ = serializedInputStream
		defer j2gort.CloseResource(serializedInputStream)
		this_.compiler.RestoreState(serializedInputStream)
	}()
}

// Java method: AbstractCommandLineRunner.saveState()
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_saveState_174() {
	if this_.config.GetSaveCompilationStateToFilename() == j2gort.NullString {
		return
	}
	var filename string = this_.config.GetSaveCompilationStateToFilename()
	_ = filename
	func() {
		defer func() {
			if rec_73 := recover(); rec_73 != nil {
				jex_74 := javalang.RecoverException(rec_73)
				_ = jex_74
				if j2gort.InstanceOf[javaio.IIOException](jex_74) {
					e := j2gort.AsIface[javaio.IIOException](jex_74)
					_ = e
					this_.compiler.Report_1(JSError_Make_2(StaticGet_AbstractCommandLineRunner_COULD_NOT_SERIALIZE_AST(), []string{filename}))
				} else {
					panic(rec_73)
				}
			}
		}()
		var serializedOutputStream *j2gort.Ext_java_io_BufferedOutputStream = j2gort.NewExt_java_io_BufferedOutputStream_1(j2gort.NewExt_java_io_FileOutputStream_1_string(filename))
		_ = serializedOutputStream
		defer j2gort.CloseResource(serializedOutputStream)
		this_.compiler.SaveState(serializedOutputStream)
	}()
}

// Java method: AbstractCommandLineRunner.processResults(com.google.javascript.jscomp.Result result, java.util.List<com.google.javascript.jscomp.JSChunk> chunks, B options)
func (this_ *AbstractCommandLineRunnerBase) ProcessResults(result *Result, chunks j2gort.Ext_java_util_List, options *CompilerOptions) int32 {
	_ = result
	_ = chunks
	_ = options
	if this_.config.printAst {
		if this_.compiler.GetRoot() == 0 {
			return int32(1)
		} else {
			var jsOutput j2gort.Ext_java_lang_Appendable = this_.priv_AbstractCommandLineRunner_createDefaultOutput_165()
			_ = jsOutput
			var cfg ControlFlowGraph[pkg_rhino.Node] = this_.compiler.ComputeCFG()
			_ = cfg
			if cfg != nil && !cfg.IsJavaNil_() {
				if r_, ok_ := any(cfg).(interface{ ReleaseSlabs_() }); ok_ {
					defer r_.ReleaseSlabs_()
				}
			}

			DotFormatter_AppendDot(this_.compiler.GetRoot().GetLastChild(), cfg, jsOutput)
			jsOutput.Append_1_uint16(uint16(10))
			AbstractCommandLineRunner_CloseAppendable(jsOutput)
			return int32(0)
		}
	}
	if this_.config.printTree {
		if this_.compiler.GetRoot() == 0 {
			this_.compiler.Report_1(JSError_Make_2(StaticGet_AbstractCommandLineRunner_NO_TREE_GENERATED_ERROR(), []string{}))
			return int32(1)
		} else {
			var jsOutput_2 j2gort.Ext_java_lang_Appendable = this_.priv_AbstractCommandLineRunner_createDefaultOutput_165()
			_ = jsOutput_2
			this_.compiler.GetJsRoot().AppendStringTree(jsOutput_2)
			jsOutput_2.Append_1_any("\n")
			AbstractCommandLineRunner_CloseAppendable(jsOutput_2)
			return int32(0)
		}
	}
	if this_.config.printTreeJson {
		if this_.compiler.GetRoot() == 0 {
			this_.compiler.Report_1(JSError_Make_2(StaticGet_AbstractCommandLineRunner_NO_TREE_GENERATED_ERROR(), []string{}))
			return int32(1)
		} else {
			var jsOutput_3 j2gort.Ext_java_lang_Appendable = this_.priv_AbstractCommandLineRunner_createDefaultOutput_165()
			_ = jsOutput_3
			this_.compiler.GetJsRoot().AppendJsonTree(jsOutput_3)
			AbstractCommandLineRunner_CloseAppendable(jsOutput_3)
			return int32(0)
		}
	}
	if this_.config.skipNormalOutputs {
		this_.priv_AbstractCommandLineRunner_outputManifest_201()
		this_.priv_AbstractCommandLineRunner_outputBundle_202()
		this_.priv_AbstractCommandLineRunner_outputChunkGraphJson_205()
		return int32(0)
	} else if (options.GetOutputJs() != StaticGet_CompilerOptions_OutputJs_NONE()) && result.success {
		this_.priv_AbstractCommandLineRunner_outputChunkGraphJson_205()
		if chunks == nil {
			this_.OutputSingleBinary(options)
			if !(this_.priv_AbstractCommandLineRunner_isOutputInJson_138()) {
				this_.priv_AbstractCommandLineRunner_outputSourceMap_196(options, this_.config.jsOutputFile)
			}
		} else {
			var error_ *DiagnosticType = this_.priv_AbstractCommandLineRunner_outputChunkBinaryAndSourceMaps_180(this_.compiler.GetChunkGraph(), options)
			_ = error_
			if error_ != nil {
				this_.compiler.Report_1(JSError_Make_2(error_, []string{}))
				return int32(1)
			}
		}
		if options.GetExternExportsPath() != j2gort.NullString {
			func() {
				var eeOut j2gort.Ext_java_io_Writer = this_.priv_AbstractCommandLineRunner_openExternExportsStream_187(options, this_.config.jsOutputFile)
				_ = eeOut
				defer j2gort.CloseResource(eeOut)
				eeOut.Append_1_any(j2gort.BoxString(result.externExport))
			}()
		}
		this_.priv_AbstractCommandLineRunner_outputNameMaps_198()
		this_.priv_AbstractCommandLineRunner_outputStringMap_199()
		this_.priv_AbstractCommandLineRunner_outputManifest_201()
		this_.priv_AbstractCommandLineRunner_outputBundle_202()
		this_.priv_AbstractCommandLineRunner_outputInstrumentationMapping_203()
		if this_.priv_AbstractCommandLineRunner_isOutputInJson_138() {
			this_.OutputJsonStream()
		}
	}
	return AbstractCommandLineRunner_GetExitStatusForResult(result)
}

// Java static method: AbstractCommandLineRunner.getExitStatusForResult(com.google.javascript.jscomp.Result result)
func AbstractCommandLineRunner_GetExitStatusForResult(result *Result) int32 {
	clinit_AbstractCommandLineRunner()
	_ = result
	return min(result.errors.Size(), int32(127))
}

// Java method: AbstractCommandLineRunner.getJavascriptEscaper()
func (this_ *AbstractCommandLineRunnerBase) GetJavascriptEscaper() j2gort.Ext_com_google_common_base_Function {
	return SourceCodeEscapers_JavascriptEscaper().AsFunction()
}

// Java method: AbstractCommandLineRunner.outputSingleBinary(B options)
func (this_ *AbstractCommandLineRunnerBase) OutputSingleBinary(options *CompilerOptions) {
	_ = options
	var escaper j2gort.Ext_com_google_common_base_Function = nil
	_ = escaper
	var marker string = "%output%"
	_ = marker
	if j2gort.JStrIndexOf(this_.config.outputWrapper, "%output|jsstring%") >= 0 {
		marker = "%output|jsstring%"
		escaper = this_.GetJavascriptEscaper()
	}
	if this_.priv_AbstractCommandLineRunner_isOutputInJson_138() {
		this_.filesToStreamOut = append(this_.filesToStreamOut, this_.CreateJsonFile(options, marker, escaper))
	} else {
		if !(j2gort.CheckNotNullString(this_.config.jsOutputFile) == "") {
			AbstractCommandLineRunner_MaybeCreateDirsForPath(this_.config.jsOutputFile)
		}
		var jsOutput j2gort.Ext_java_lang_Appendable = this_.priv_AbstractCommandLineRunner_createDefaultOutput_165()
		_ = jsOutput
		this_.WriteOutput_8(jsOutput, this_.compiler, NewCompiler_SingleBinaryLicenseTracker(ToIface_AbstractCompiler(this_.compiler)), nil, this_.config.outputWrapper, marker, escaper, this_.config.jsOutputFile)
		AbstractCommandLineRunner_CloseAppendable(jsOutput)
	}
}

// Java method: AbstractCommandLineRunner.createJsonFile(B options, java.lang.String outputMarker, com.google.common.base.Function<java.lang.String,java.lang.String> escaper)
func (this_ *AbstractCommandLineRunnerBase) CreateJsonFile(options *CompilerOptions, outputMarker string, escaper j2gort.Ext_com_google_common_base_Function) *AbstractCommandLineRunner_JsonFileSpec {
	_ = options
	_ = outputMarker
	_ = escaper
	var jsOutput j2gort.Ext_java_lang_Appendable = j2gort.NewExt_java_lang_StringBuilder()
	_ = jsOutput
	this_.WriteOutput_8(jsOutput, this_.compiler, NewCompiler_SingleBinaryLicenseTracker(ToIface_AbstractCompiler(this_.compiler)), nil, this_.config.outputWrapper, outputMarker, escaper, this_.config.jsOutputFile)
	var val_75 AbstractCommandLineRunner_JsonFileSpec
	InitAbstractCommandLineRunner_JsonFileSpec_2(&val_75, javalang.ToJavaString(jsOutput), func() string {
		if j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(this_.config.jsOutputFile) {
			return "compiled.js"
		}
		return this_.config.jsOutputFile
	}())
	var jsonOutput *AbstractCommandLineRunner_JsonFileSpec = &val_75
	_ = jsonOutput
	if options.ShouldGatherSourceMapInfo() {
		var sourcemap *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = sourcemap
		this_.compiler.GetSourceMap().AppendTo(sourcemap, jsonOutput.GetPath())
		jsonOutput.SetSourceMap(sourcemap.String())
	}
	return jsonOutput
}

// Java method: AbstractCommandLineRunner.outputJsonStream()
func (this_ *AbstractCommandLineRunnerBase) OutputJsonStream() {
	func() {
		var jsonWriter *j2gort.Ext_com_google_gson_stream_JsonWriter = j2gort.NewExt_com_google_gson_stream_JsonWriter(j2gort.NewExt_java_io_BufferedWriter_1(j2gort.NewExt_java_io_OutputStreamWriter_2_Ext_java_io_OutputStream_Ext_java_nio_charset_Charset(j2gort.AsIface[j2gort.Ext_java_io_OutputStream](this_.defaultJsOutput), j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8)))
		_ = jsonWriter
		defer j2gort.CloseResource(jsonWriter)
		jsonWriter.SetHtmlSafe(false)
		jsonWriter.BeginArray()
		for _, spec := range this_.filesToStreamOut {
			_ = spec
			spec.WriteJson(jsonWriter)
		}

		jsonWriter.EndArray()
	}()
}

// Java method: AbstractCommandLineRunner.outputChunkBinaryAndSourceMaps(com.google.javascript.jscomp.JSChunkGraph chunkGraph, B options)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_outputChunkBinaryAndSourceMaps_180(chunkGraph *JSChunkGraph, options *CompilerOptions) *DiagnosticType {
	_ = chunkGraph
	_ = options
	var chunks j2gort.Ext_java_lang_Iterable = chunkGraph.GetAllChunks()
	_ = chunks
	this_.parsedChunkWrappers = AbstractCommandLineRunner_ParseChunkWrappers(this_.config.chunkWrapper, chunks)
	if !(this_.priv_AbstractCommandLineRunner_isOutputInJson_138()) {
		AbstractCommandLineRunner_MaybeCreateDirsForPath((j2gort.StrValString(this_.config.chunkOutputPathPrefix) + "dummy"))
	}
	var mapFileOut j2gort.Ext_java_io_Writer = nil
	_ = mapFileOut
	if !(((this_.ShouldGenerateMapPerChunk(options) || !(options.ShouldGatherSourceMapInfo())) || (this_.config.jsonStreamMode == StaticGet_CompilerOptions_JsonStreamMode_OUT())) || (this_.config.jsonStreamMode == StaticGet_CompilerOptions_JsonStreamMode_BOTH())) {
		return StaticGet_AbstractCommandLineRunner_INVALID_CHUNK_SOURCEMAP_PATTERN()
	}
	var val_78 Compiler_ChunkGraphAwareLicenseTracker
	InitCompiler_ChunkGraphAwareLicenseTracker(&val_78, ToIface_AbstractCompiler(this_.compiler))
	var mlicenseTracker *Compiler_ChunkGraphAwareLicenseTracker = &val_78
	_ = mlicenseTracker
	{
		cur_79 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunks))
		for cur_79.HasNext() {
			rawElem_80 := cur_79.Next()
			m := j2gort.CastPtr[JSChunk](rawElem_80)
			_ = m
			{
				var iifeCond_1 bool
				{

					sr_81 := m.GetName()
					se_82 := "$weak$"
					iifeCond_1 = j2gort.CheckNotNullString(sr_81) == se_82
				}
				if iifeCond_1 {
					continue
				}
			}
			if this_.priv_AbstractCommandLineRunner_isOutputInJson_138() {
				this_.filesToStreamOut = append(this_.filesToStreamOut, this_.priv_AbstractCommandLineRunner_createJsonFileFromChunk_181(m))
			} else {
				if this_.ShouldGenerateMapPerChunk(options) {
					mapFileOut = this_.priv_AbstractCommandLineRunner_fileNameToOutputWriter2_191(this_.ExpandSourceMapPath(options, m))
				}
				var chunkFilename string = this_.GetChunkOutputFileName(m)
				_ = chunkFilename
				AbstractCommandLineRunner_MaybeCreateDirsForPath(chunkFilename)
				func() {
					var writer j2gort.Ext_java_io_Writer = this_.priv_AbstractCommandLineRunner_fileNameToLegacyOutputWriter_190(chunkFilename)
					_ = writer
					defer j2gort.CloseResource(writer)
					if options.ShouldGatherSourceMapInfo() {
						this_.compiler.ResetAndIntitializeSourceMap()
					}
					mlicenseTracker.SetCurrentChunkContext(m)
					this_.WriteChunkOutput(chunkFilename, j2gort.AsIface[j2gort.Ext_java_lang_Appendable](writer), mlicenseTracker, m)
					if options.ShouldGatherSourceMapInfo() {
						this_.compiler.GetSourceMap().AppendTo(j2gort.AsIface[j2gort.Ext_java_lang_Appendable](mapFileOut), chunkFilename)
					}
				}()
				if this_.ShouldGenerateMapPerChunk(options) && (mapFileOut != nil) {
					mapFileOut.Close_0()
					mapFileOut = nil
				}
			}
		}
	}
	if mapFileOut != nil {
		mapFileOut.Close_0()
	}
	return nil
}

// Java method: AbstractCommandLineRunner.createJsonFileFromChunk(com.google.javascript.jscomp.JSChunk chunk)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_createJsonFileFromChunk_181(chunk *JSChunk) *AbstractCommandLineRunner_JsonFileSpec {
	_ = chunk
	this_.compiler.ResetAndIntitializeSourceMap()
	var filename string = this_.GetChunkOutputFileName(chunk)
	_ = filename
	var output *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = output
	this_.WriteChunkOutput(filename, output, NewCompiler_ScriptNodeLicensesOnlyTracker(ToIface_AbstractCompiler(this_.compiler)), chunk)
	var val_83 AbstractCommandLineRunner_JsonFileSpec
	InitAbstractCommandLineRunner_JsonFileSpec_2(&val_83, output.String(), filename)
	var jsonFile *AbstractCommandLineRunner_JsonFileSpec = &val_83
	_ = jsonFile
	var chunkSourceMap *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = chunkSourceMap
	this_.compiler.GetSourceMap().AppendTo(chunkSourceMap, this_.GetChunkOutputFileName(chunk))
	jsonFile.SetSourceMap(chunkSourceMap.String())
	return jsonFile
}

// Java method: AbstractCommandLineRunner.getInputCharset()
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_getInputCharset_182() j2gort.Ext_java_nio_charset_Charset {
	if !(j2gort.CheckNotNullString(this_.config.charset) == "") {
		if !(j2gort.Ext_java_nio_charset_Charset_IsSupported(this_.config.charset)) {
			panic(NewFlagUsageException((j2gort.StrValString(this_.config.charset) + " is not a valid charset name.")))
		}
		return j2gort.Ext_java_nio_charset_Charset_ForName_1(this_.config.charset)
	}
	return j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8
}

// Java method: AbstractCommandLineRunner.getLegacyOutputCharset()
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_getLegacyOutputCharset_183() j2gort.Ext_java_nio_charset_Charset {
	if !(j2gort.CheckNotNullString(this_.config.charset) == "") {
		if !(j2gort.Ext_java_nio_charset_Charset_IsSupported(this_.config.charset)) {
			panic(NewFlagUsageException((j2gort.StrValString(this_.config.charset) + " is not a valid charset name.")))
		}
		return j2gort.Ext_java_nio_charset_Charset_ForName_1(this_.config.charset)
	}
	return j2gort.Ext_java_nio_charset_StandardCharsets_F_US_ASCII
}

// Java method: AbstractCommandLineRunner.getOutputCharset2()
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_getOutputCharset2_184() j2gort.Ext_java_nio_charset_Charset {
	if !(j2gort.CheckNotNullString(this_.config.charset) == "") {
		if !(j2gort.Ext_java_nio_charset_Charset_IsSupported(this_.config.charset)) {
			panic(NewFlagUsageException((j2gort.StrValString(this_.config.charset) + " is not a valid charset name.")))
		}
		return j2gort.Ext_java_nio_charset_Charset_ForName_1(this_.config.charset)
	}
	return j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8
}

// Java method: AbstractCommandLineRunner.createExterns(com.google.javascript.jscomp.CompilerOptions options)
func (this_ *AbstractCommandLineRunnerBase) CreateExterns(options *CompilerOptions) j2gort.Ext_java_util_List {
	_ = options
	if this_.IsInTestMode() {
		return j2gort.AsIface[j2gort.Ext_java_util_List](this_.externsSupplierForTesting.Get())
	} else {
		return this_.priv_AbstractCommandLineRunner_createExternInputs_158(&this_.config.externs)
	}
}

// Java method: AbstractCommandLineRunner.shouldGenerateMapPerChunk(B options)
func (this_ *AbstractCommandLineRunnerBase) ShouldGenerateMapPerChunk(options *CompilerOptions) bool {
	_ = options
	return (options.ShouldGatherSourceMapInfo() && (j2gort.JStrIndexOf(options.GetSourceMapOutputPath(), "%outname%") >= 0))
}

// Java method: AbstractCommandLineRunner.openExternExportsStream(B options, java.lang.String path)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_openExternExportsStream_187(options *CompilerOptions, path string) j2gort.Ext_java_io_Writer {
	_ = options
	_ = path
	var externExportsPath string = options.GetExternExportsPath()
	_ = externExportsPath
	if externExportsPath == j2gort.NullString {
		return nil
	}
	var exPath string = externExportsPath
	_ = exPath
	if !(j2gort.JStrIndexOf(exPath, j2gort.CheckNotNullString(j2gort.Ext_java_io_File_F_Separator)) >= 0) {
		var outputFile *j2gort.Ext_java_io_File = j2gort.NewExt_java_io_File_1_string(path)
		_ = outputFile
		exPath = (j2gort.StrValString(outputFile.GetParent()) + j2gort.StrValChar(j2gort.Ext_java_io_File_F_SeparatorChar) + j2gort.StrValString(exPath))
	}
	return this_.priv_AbstractCommandLineRunner_fileNameToOutputWriter2_191(exPath)
}

// Java method: AbstractCommandLineRunner.expandCommandLinePath(java.lang.String path, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable JSChunk forChunk)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_expandCommandLinePath_188(path string, forChunk *JSChunk) string {
	_ = path
	_ = forChunk
	var sub string = j2gort.NullString
	_ = sub
	if forChunk != nil {
		sub = this_.GetChunkOutputFileName(forChunk)
	} else if len(this_.config.chunk) != 0 {
		sub = this_.config.chunkOutputPathPrefix
	} else {
		sub = this_.config.jsOutputFile
	}
	return javalang.StringReplace(path, "%outname%", j2gort.CheckNotNullString(sub))
}

// Java method: AbstractCommandLineRunner.expandSourceMapPath(B options, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable JSChunk forChunk)
func (this_ *AbstractCommandLineRunnerBase) ExpandSourceMapPath(options *CompilerOptions, forChunk *JSChunk) string {
	_ = options
	_ = forChunk
	if !(options.ShouldGatherSourceMapInfo()) {
		return j2gort.NullString
	}
	return this_.priv_AbstractCommandLineRunner_expandCommandLinePath_188(options.GetSourceMapOutputPath(), forChunk)
}

// Java method: AbstractCommandLineRunner.fileNameToLegacyOutputWriter(java.lang.String fileName)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_fileNameToLegacyOutputWriter_190(fileName string) j2gort.Ext_java_io_Writer {
	_ = fileName
	if fileName == j2gort.NullString {
		return nil
	}
	if this_.IsInTestMode() {
		return j2gort.NewExt_java_io_StringWriter()
	}
	return this_.priv_AbstractCommandLineRunner_streamToLegacyOutputWriter_193(this_.FilenameToOutputStream(fileName))
}

// Java method: AbstractCommandLineRunner.fileNameToOutputWriter2(java.lang.String fileName)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_fileNameToOutputWriter2_191(fileName string) j2gort.Ext_java_io_Writer {
	_ = fileName
	if fileName == j2gort.NullString {
		return nil
	}
	if this_.IsInTestMode() {
		return j2gort.NewExt_java_io_StringWriter()
	}
	return this_.priv_AbstractCommandLineRunner_streamToOutputWriter2_194(this_.FilenameToOutputStream(fileName))
}

// Java method: AbstractCommandLineRunner.filenameToOutputStream(java.lang.String fileName)
func (this_ *AbstractCommandLineRunnerBase) FilenameToOutputStream(fileName string) j2gort.Ext_java_io_OutputStream {
	_ = fileName
	if fileName == j2gort.NullString {
		return nil
	}
	return j2gort.NewExt_java_io_FileOutputStream_1_string(fileName)
}

// Java method: AbstractCommandLineRunner.streamToLegacyOutputWriter(java.io.OutputStream stream)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_streamToLegacyOutputWriter_193(stream j2gort.Ext_java_io_OutputStream) j2gort.Ext_java_io_Writer {
	_ = stream
	if this_.legacyOutputCharset == nil {
		return this_.CreateWriter(stream, j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8)
	} else {
		return this_.CreateWriter(stream, this_.legacyOutputCharset)
	}
}

// Java method: AbstractCommandLineRunner.streamToOutputWriter2(java.io.OutputStream stream)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_streamToOutputWriter2_194(stream j2gort.Ext_java_io_OutputStream) j2gort.Ext_java_io_Writer {
	_ = stream
	if this_.outputCharset2 == nil {
		return this_.CreateWriter(stream, j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8)
	} else {
		return this_.CreateWriter(stream, this_.outputCharset2)
	}
}

// Java method: AbstractCommandLineRunner.createWriter(java.io.OutputStream stream, java.nio.charset.Charset charset)
func (this_ *AbstractCommandLineRunnerBase) CreateWriter(stream j2gort.Ext_java_io_OutputStream, charset j2gort.Ext_java_nio_charset_Charset) j2gort.Ext_java_io_Writer {
	_ = stream
	_ = charset
	return j2gort.NewExt_java_io_BufferedWriter_1(j2gort.NewExt_java_io_OutputStreamWriter_2_Ext_java_io_OutputStream_Ext_java_nio_charset_Charset(stream, charset))
}

// Java method: AbstractCommandLineRunner.outputSourceMap(B options, java.lang.String associatedName)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_outputSourceMap_196(options *CompilerOptions, associatedName string) {
	_ = options
	_ = associatedName
	if !(options.ShouldGatherSourceMapInfo()) || (j2gort.CheckNotNullString(options.GetSourceMapOutputPath()) == "/dev/null") {
		return
	}
	var outName string = this_.ExpandSourceMapPath(options, nil)
	_ = outName
	AbstractCommandLineRunner_MaybeCreateDirsForPath(outName)
	func() {
		var out j2gort.Ext_java_io_Writer = this_.priv_AbstractCommandLineRunner_fileNameToOutputWriter2_191(outName)
		_ = out
		defer j2gort.CloseResource(out)
		this_.compiler.GetSourceMap().AppendTo(j2gort.AsIface[j2gort.Ext_java_lang_Appendable](out), associatedName)
	}()
}

// Java method: AbstractCommandLineRunner.getMapPath(java.lang.String outputFile)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_getMapPath_197(outputFile string) string {
	_ = outputFile
	var basePath string = ""
	_ = basePath
	if j2gort.CheckNotNullString(outputFile) == "" {
		if !(j2gort.CheckNotNullString(this_.config.chunkOutputPathPrefix) == "") {
			basePath = this_.config.chunkOutputPathPrefix
		} else {
			basePath = "jscompiler"
		}
	} else {
		var file *j2gort.Ext_java_io_File = j2gort.NewExt_java_io_File_1_string(outputFile)
		_ = file
		var outputFileName string = file.GetName()
		_ = outputFileName
		if javalang.StringEndsWith(j2gort.CheckNotNullString(outputFileName), ".js") {
			outputFileName = j2gort.JStrSub(outputFileName, int32(0), (j2gort.JStrLen(outputFileName) - int32(3)))
		}
		var fileParent string = file.GetParent()
		_ = fileParent
		if fileParent == j2gort.NullString {
			basePath = outputFileName
		} else {
			basePath = (j2gort.StrValString(file.GetParent()) + j2gort.StrValChar(j2gort.Ext_java_io_File_F_SeparatorChar) + j2gort.StrValString(outputFileName))
		}
	}
	return basePath
}

// Java method: AbstractCommandLineRunner.outputNameMaps()
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_outputNameMaps_198() {
	var propertyMapOutputPath string = j2gort.NullString
	_ = propertyMapOutputPath
	var variableMapOutputPath string = j2gort.NullString
	_ = variableMapOutputPath
	if this_.config.createNameMapFiles {
		var basePath string = this_.priv_AbstractCommandLineRunner_getMapPath_197(this_.config.jsOutputFile)
		_ = basePath
		propertyMapOutputPath = (j2gort.StrValString(basePath) + "_props_map.out")
		variableMapOutputPath = (j2gort.StrValString(basePath) + "_vars_map.out")
	}
	if !(j2gort.CheckNotNullString(this_.config.variableMapOutputFile) == "") {
		if variableMapOutputPath != j2gort.NullString {
			panic(NewFlagUsageException("The flags variable_map_output_file and create_name_map_files cannot both be used simultaneously."))
		}
		variableMapOutputPath = this_.config.variableMapOutputFile
	}
	if !(j2gort.CheckNotNullString(this_.config.propertyMapOutputFile) == "") {
		if propertyMapOutputPath != j2gort.NullString {
			panic(NewFlagUsageException("The flags property_map_output_file and create_name_map_files cannot both be used simultaneously."))
		}
		propertyMapOutputPath = this_.config.propertyMapOutputFile
	}
	if (variableMapOutputPath != j2gort.NullString) && (this_.compiler.GetVariableMap() != nil) {
		this_.compiler.GetVariableMap().Save(variableMapOutputPath)
	}
	if (propertyMapOutputPath != j2gort.NullString) && (this_.compiler.GetPropertyMap() != nil) {
		this_.compiler.GetPropertyMap().Save(propertyMapOutputPath)
	}
}

// Java method: AbstractCommandLineRunner.outputStringMap()
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_outputStringMap_199() {
	if !(j2gort.CheckNotNullString(this_.config.stringMapOutputPath) == "") {
		if this_.compiler.GetStringMap() == nil {
			if !(j2gort.NewExt_java_io_File_1_string(this_.config.stringMapOutputPath).CreateNewFile()) {
				panic(j2gort.AsIface[javaio.IIOException](j2gort.ConstructShimThrowable("java.io.IOException", "Could not create file: "+j2gort.StrValString(this_.config.stringMapOutputPath))))
			}
		} else {
			this_.compiler.GetStringMap().Save(this_.config.stringMapOutputPath)
		}
	}
}

// Java static method: AbstractCommandLineRunner.createDefineReplacements(java.util.List<java.lang.String> definitions, com.google.javascript.jscomp.CompilerOptions options)
func AbstractCommandLineRunner_CreateDefineReplacements(definitions j2gort.Ext_java_util_List, options *CompilerOptions) {
	clinit_AbstractCommandLineRunner()
	_ = definitions
	_ = options
	{
		cur_84 := j2gort.NewIterCursor(j2gort.CheckNotNull(definitions))
		for cur_84.HasNext() {
			override := cur_84.NextStr()
			_ = override
			var assignment []string = j2gort.JStrSplitLimit(override, "=", int32(2))
			_ = assignment
			var defName string = assignment[int32(0)]
			_ = defName
			if j2gort.JStrLen(defName) > int32(0) {
				var defValue string
				_ = defValue
				if int32(len(assignment)) == int32(1) {
					defValue = "true"
				} else {
					defValue = assignment[int32(1)]
				}
				var isTrue bool = (j2gort.CheckNotNullString(defValue) == "true")
				_ = isTrue
				var isFalse bool = (j2gort.CheckNotNullString(defValue) == "false")
				_ = isFalse
				if isTrue || isFalse {
					options.SetDefineToBooleanLiteral(defName, isTrue)
					continue
				} else if (j2gort.JStrLen(defValue) > int32(1)) && (((int32(j2gort.JStrCharAt(defValue, int32(0))) == int32(uint16(39))) && (int32(j2gort.JStrCharAt(defValue, (j2gort.JStrLen(defValue)-int32(1)))) == int32(uint16(39)))) || ((int32(j2gort.JStrCharAt(defValue, int32(0))) == int32(uint16(34))) && (int32(j2gort.JStrCharAt(defValue, (j2gort.JStrLen(defValue)-int32(1)))) == int32(uint16(34))))) {
					var maybeStringVal string = j2gort.JStrSub(defValue, int32(1), (j2gort.JStrLen(defValue) - int32(1)))
					_ = maybeStringVal
					if j2gort.JStrIndexOfCp(maybeStringVal, int32(j2gort.JStrCharAt(defValue, int32(0)))) == int32(-1) {
						options.SetDefineToStringLiteral(defName, maybeStringVal)
						continue
					}
				} else {
					var tryFlow_85 int
					func() {
						defer func() {
							if rec_86 := recover(); rec_86 != nil {
								jex_87 := javalang.RecoverException(rec_86)
								_ = jex_87
								if j2gort.InstanceOf[javalang.INumberFormatException](jex_87) {
									e := j2gort.AsIface[javalang.INumberFormatException](jex_87)
									_ = e
								} else {
									panic(rec_86)
								}
							}
						}()
						var value float64 = j2gort.ParseDouble(defValue)
						_ = value
						options.SetDefineToDoubleLiteral(defName, value)
						tryFlow_85 = 2
						return
					}()
					if tryFlow_85 == 2 {
						continue
					}
					if j2gort.JStrLen(defValue) > int32(0) {
						options.SetDefineToStringLiteral(defName, defValue)
						continue
					}
				}
			}
			panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "--define flag syntax invalid: "+j2gort.StrValString(override))))
		}
	}
}

// Java method: AbstractCommandLineRunner.shouldGenerateOutputPerChunk(java.lang.String output)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_shouldGenerateOutputPerChunk_200(output string) bool {
	_ = output
	return ((len(this_.config.chunk) != 0 && (output != j2gort.NullString)) && (j2gort.JStrIndexOf(output, "%outname%") >= 0))
}

// Java method: AbstractCommandLineRunner.outputManifest()
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_outputManifest_201() {
	this_.priv_AbstractCommandLineRunner_outputManifestOrBundle_204(this_.config.outputManifests, true)
}

// Java method: AbstractCommandLineRunner.outputBundle()
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_outputBundle_202() {
	this_.priv_AbstractCommandLineRunner_outputManifestOrBundle_204(this_.config.outputBundles, false)
}

// Java method: AbstractCommandLineRunner.outputInstrumentationMapping()
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_outputInstrumentationMapping_203() {
	if !(j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(this_.config.instrumentationMappingFile)) {
		var path string = this_.priv_AbstractCommandLineRunner_expandCommandLinePath_188(this_.config.instrumentationMappingFile, nil)
		_ = path
		this_.compiler.GetInstrumentationMapping().Save(path)
	}
}

// Java method: AbstractCommandLineRunner.outputManifestOrBundle(java.util.List<java.lang.String> outputFiles, boolean isManifest)
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_outputManifestOrBundle_204(outputFiles j2gort.Ext_java_util_List, isManifest bool) {
	_ = outputFiles
	_ = isManifest
	if outputFiles.IsEmpty() {
		return
	}
	{
		cur_88 := j2gort.NewIterCursor(j2gort.CheckNotNull(outputFiles))
		for cur_88.HasNext() {
			output := cur_88.NextStr()
			_ = output
			if j2gort.CheckNotNullString(output) == "" {
				continue
			}
			if this_.priv_AbstractCommandLineRunner_shouldGenerateOutputPerChunk_200(output) {
				var chunks j2gort.Ext_java_lang_Iterable = this_.compiler.GetChunkGraph().GetAllChunks()
				_ = chunks
				{
					cur_89 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunks))
					for cur_89.HasNext() {
						rawElem_90 := cur_89.Next()
						chunk := j2gort.CastPtr[JSChunk](rawElem_90)
						_ = chunk
						func() {
							var out j2gort.Ext_java_io_Writer = this_.priv_AbstractCommandLineRunner_fileNameToOutputWriter2_191(this_.priv_AbstractCommandLineRunner_expandCommandLinePath_188(output, chunk))
							_ = out
							defer j2gort.CloseResource(out)
							if isManifest {
								this_.PrintManifestTo(chunk, j2gort.AsIface[j2gort.Ext_java_lang_Appendable](out))
							} else {
								this_.PrintBundleTo(chunk, j2gort.AsIface[j2gort.Ext_java_lang_Appendable](out))
							}
						}()
					}
				}
			} else {
				func() {
					var out_2 j2gort.Ext_java_io_Writer = this_.priv_AbstractCommandLineRunner_fileNameToOutputWriter2_191(this_.priv_AbstractCommandLineRunner_expandCommandLinePath_188(output, nil))
					_ = out_2
					defer j2gort.CloseResource(out_2)
					if len(this_.config.chunk) == 0 {
						var chunk_2 *JSChunk = this_.compiler.GetChunkGraph().GetChunkByName("$strong$")
						_ = chunk_2
						if isManifest {
							this_.PrintManifestTo(chunk_2, j2gort.AsIface[j2gort.Ext_java_lang_Appendable](out_2))
						} else {
							this_.PrintBundleTo(chunk_2, j2gort.AsIface[j2gort.Ext_java_lang_Appendable](out_2))
						}
					} else {
						this_.PrintChunkGraphManifestOrBundleTo(this_.compiler.GetChunkGraph(), j2gort.AsIface[j2gort.Ext_java_lang_Appendable](out_2), isManifest)
					}
				}()
			}
		}
	}
}

// Java method: AbstractCommandLineRunner.outputChunkGraphJson()
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_outputChunkGraphJson_205() {
	if (this_.config.outputChunkDependencies != j2gort.NullString) && (j2gort.JStrLen(this_.config.outputChunkDependencies) != int32(0)) {
		func() {
			var out j2gort.Ext_java_io_Writer = this_.priv_AbstractCommandLineRunner_fileNameToOutputWriter2_191(this_.config.outputChunkDependencies)
			_ = out
			defer j2gort.CloseResource(out)
			this_.PrintChunkGraphJsonTo(j2gort.AsIface[j2gort.Ext_java_lang_Appendable](out))
		}()
	}
}

// Java method: AbstractCommandLineRunner.printChunkGraphJsonTo(java.lang.Appendable out)
func (this_ *AbstractCommandLineRunnerBase) PrintChunkGraphJsonTo(out j2gort.Ext_java_lang_Appendable) {
	_ = out
	out.Append_1_any(j2gort.BoxString(this_.compiler.GetChunkGraph().ToJson().String()))
}

// Java method: AbstractCommandLineRunner.printChunkGraphManifestOrBundleTo(com.google.javascript.jscomp.JSChunkGraph graph, java.lang.Appendable out, boolean isManifest)
func (this_ *AbstractCommandLineRunnerBase) PrintChunkGraphManifestOrBundleTo(graph *JSChunkGraph, out j2gort.Ext_java_lang_Appendable, isManifest bool) {
	_ = graph
	_ = out
	_ = isManifest
	var commas *j2gort.Ext_com_google_common_base_Joiner = j2gort.Ext_com_google_common_base_Joiner_On_1_string(",")
	_ = commas
	var requiresNewline bool = false
	_ = requiresNewline
	{
		cur_91 := j2gort.NewIterCursor(j2gort.CheckNotNull(graph.GetAllChunks()))
		for cur_91.HasNext() {
			rawElem_92 := cur_91.Next()
			chunk := j2gort.CastPtr[JSChunk](rawElem_92)
			_ = chunk
			if !(isManifest) && chunk.IsWeak() {
				continue
			}
			if requiresNewline {
				out.Append_1_any("\n")
			}
			if isManifest {
				var dependencies string = commas.Join_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](chunk.GetSortedDependencyNames()))
				_ = dependencies
				out.Append_1_any(j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("{%s%s}\n", []any{j2gort.BoxString(chunk.GetName()), func() any {
					if j2gort.CheckNotNullString(dependencies) == "" {
						return ""
					}
					return ":" + j2gort.StrValString(dependencies)
				}()})))
				this_.PrintManifestTo(chunk, out)
			} else {
				this_.PrintBundleTo(chunk, out)
			}
			requiresNewline = true
		}
	}
}

// Java method: AbstractCommandLineRunner.printManifestTo(com.google.javascript.jscomp.JSChunk chunk, java.lang.Appendable out)
func (this_ *AbstractCommandLineRunnerBase) PrintManifestTo(chunk *JSChunk, out j2gort.Ext_java_lang_Appendable) {
	_ = chunk
	_ = out
	{
		cur_93 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunk.GetInputs()))
		for cur_93.HasNext() {
			rawElem_94 := cur_93.Next()
			input := j2gort.CastPtr[CompilerInput](rawElem_94)
			_ = input
			var name string = input.GetName()
			_ = name
			if AbstractCompiler_IsFillFileName(name) {
				continue
			}
			var rootRelativePath string = j2gort.CastString(j2gort.MapGetStr(this_.rootRelativePathsMap, name))
			_ = rootRelativePath
			var displayName string
			_ = displayName
			if rootRelativePath != j2gort.NullString {
				displayName = rootRelativePath
			} else {
				displayName = name
			}
			out.Append_1_any(j2gort.BoxString(displayName))
			out.Append_1_any("\n")
		}
	}
}

// Java method: AbstractCommandLineRunner.printBundleTo(com.google.javascript.jscomp.JSChunk chunk, java.lang.Appendable out)
func (this_ *AbstractCommandLineRunnerBase) PrintBundleTo(chunk *JSChunk, out j2gort.Ext_java_lang_Appendable) {
	_ = chunk
	_ = out
	var inputs []*CompilerInput = j2gort.ImmListToSliceNullable[*CompilerInput](chunk.GetInputs())
	_ = inputs
	var options *CompilerOptions = this_.compiler.GetOptions()
	_ = options
	if options.GetNumParallelThreads() > int32(1) {
		{
			var val_95 PrebuildAst
			InitPrebuildAst(&val_95, ToIface_AbstractCompiler(this_.compiler), this_.compiler.GetOptions().GetNumParallelThreads())
			(&val_95).Prebuild(j2gort.ImmListFromTypedSlice(inputs))
		}
	}
	if options.GetRuntimeLibraryMode() == StaticGet_RuntimeJsLibManager_RuntimeLibraryMode_INJECT() {
		{

			for _, input := range inputs {

				_ = input
				if input.IsEs6Module() {
					this_.self_.AppendRuntimeTo(out)
					break
				}
			}
		}
	}
	{

		for _, input_2 := range inputs {

			_ = input_2
			var name string = input_2.GetName()
			_ = name
			var code string = input_2.GetSourceFile().GetCode()
			_ = code
			if AbstractCompiler_IsFillFileName(name) && (j2gort.CheckNotNullString(code) == "") {
				continue
			}
			var rootRelativePath string = j2gort.CastString(j2gort.MapGetStr(this_.rootRelativePathsMap, name))
			_ = rootRelativePath
			var displayName string
			_ = displayName
			if rootRelativePath != j2gort.NullString {
				displayName = rootRelativePath
			} else {
				displayName = input_2.GetName()
			}
			out.Append_1_any("//")
			out.Append_1_any(j2gort.BoxString(displayName))
			out.Append_1_any("\n")
			this_.self_.PrepForBundleAndAppendTo(out, input_2, code)
			out.Append_1_any("\n")
		}
	}
}

// Java method: AbstractCommandLineRunner.constructRootRelativePathsMap()
func (this_ *AbstractCommandLineRunnerBase) priv_AbstractCommandLineRunner_constructRootRelativePathsMap_210() j2gort.Ext_java_util_Map {
	var rootRelativePathsMap j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = rootRelativePathsMap
	{
		cur_100 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.config.manifestMaps))
		for cur_100.HasNext() {
			mapString := cur_100.NextStr()
			_ = mapString
			var colonIndex int32 = j2gort.JStrIndexOfCp(mapString, int32(uint16(58)))
			_ = colonIndex
			if !(colonIndex > int32(0)) {
				j2gort.CheckState(false)
			}
			var execPath string = j2gort.JStrSub(mapString, int32(0), colonIndex)
			_ = execPath
			var rootRelativePath string = j2gort.JStrSubFrom(mapString, (colonIndex + int32(1)))
			_ = rootRelativePath
			if !(j2gort.JStrIndexOfCp(rootRelativePath, int32(uint16(58))) == int32(-1)) {
				j2gort.CheckState(false)
			}
			j2gort.MapPutStr(rootRelativePathsMap, execPath, j2gort.BoxString(rootRelativePath))
		}
	}
	return rootRelativePathsMap
}

// Java class: com.google.javascript.jscomp.AbstractCommandLineRunner.CommandLineConfig
type AbstractCommandLineRunner_CommandLineConfig struct {
	printVersion                     bool                                                          `gson:"printVersion"`
	printTree                        bool                                                          `gson:"printTree"`
	printTreeJson                    bool                                                          `gson:"printTreeJson"`
	printAst                         bool                                                          `gson:"printAst"`
	jscompDevMode                    CompilerOptions_DevMode                                       `gson:"jscompDevMode"`
	loggingLevel                     string                                                        `gson:"loggingLevel"`
	externs                          []string                                                      `gson:"externs"`
	mixedJsSources                   j2gort.Ext_java_util_List                                     `gson:"mixedJsSources"`
	defaultToStdin                   bool                                                          `gson:"defaultToStdin"`
	jsOutputFile                     string                                                        `gson:"jsOutputFile"`
	continueSavedCompilationFileName string                                                        `gson:"continueSavedCompilationFileName"`
	restoredCompilationStage         int32                                                         `gson:"restoredCompilationStage"`
	saveAfterCompilationStage        int32                                                         `gson:"saveAfterCompilationStage"`
	typedAstListInputFilename        string                                                        `gson:"typedAstListInputFilename"`
	saveCompilationStateToFilename   string                                                        `gson:"saveCompilationStateToFilename"`
	chunk                            []string                                                      `gson:"chunk"`
	sourceMapInputFiles              j2gort.Ext_java_util_Map                                      `gson:"sourceMapInputFiles"`
	parseInlineSourceMaps            bool                                                          `gson:"parseInlineSourceMaps"`
	expectedDiagnostics              j2gort.Ext_java_util_List                                     `gson:"expectedDiagnostics"`
	variableMapInputFile             string                                                        `gson:"variableMapInputFile"`
	propertyMapInputFile             string                                                        `gson:"propertyMapInputFile"`
	variableMapOutputFile            string                                                        `gson:"variableMapOutputFile"`
	createNameMapFiles               bool                                                          `gson:"createNameMapFiles"`
	propertyMapOutputFile            string                                                        `gson:"propertyMapOutputFile"`
	stringMapOutputPath              string                                                        `gson:"stringMapOutputPath"`
	instrumentationMappingFile       string                                                        `gson:"instrumentationMappingFile"`
	codingConvention                 CodingConvention                                              `gson:"codingConvention"`
	summaryDetailLevel               int32                                                         `gson:"summaryDetailLevel"`
	outputWrapper                    string                                                        `gson:"outputWrapper"`
	chunkWrapper                     j2gort.Ext_java_util_List                                     `gson:"chunkWrapper"`
	chunkOutputPathPrefix            string                                                        `gson:"chunkOutputPathPrefix"`
	chunkOutputFiles                 j2gort.Ext_java_util_List                                     `gson:"chunkOutputFiles"`
	chunkConformanceFiles            j2gort.Ext_java_util_List                                     `gson:"chunkConformanceFiles"`
	createSourceMap                  string                                                        `gson:"createSourceMap"`
	sourceMapDetailLevel             *SourceMap_DetailLevel                                        `gson:"sourceMapDetailLevel"`
	sourceMapFormat                  *SourceMap_Format                                             `gson:"sourceMapFormat"`
	sourceMapLocationMappings        j2gort.Ext_com_google_common_collect_ImmutableList            `gson:"sourceMapLocationMappings"`
	applyInputSourceMaps             bool                                                          `gson:"applyInputSourceMaps"`
	warningGuards                    j2gort.Ext_java_util_ArrayList                                `gson:"warningGuards"`
	define                           j2gort.Ext_java_util_List                                     `gson:"define"`
	browserFeaturesetYear            int32                                                         `gson:"browserFeaturesetYear"`
	tweakProcessing                  CompilerOptions_TweakProcessing                               `gson:"tweakProcessing"`
	charset                          string                                                        `gson:"charset"`
	dependencyOptions                *DependencyOptions                                            `gson:"dependencyOptions"`
	outputManifests                  j2gort.Ext_java_util_List                                     `gson:"outputManifests"`
	outputChunkDependencies          string                                                        `gson:"outputChunkDependencies"`
	outputBundles                    j2gort.Ext_java_util_List                                     `gson:"outputBundles"`
	skipNormalOutputs                bool                                                          `gson:"skipNormalOutputs"`
	manifestMaps                     j2gort.Ext_java_util_List                                     `gson:"manifestMaps"`
	processCommonJSModules           bool                                                          `gson:"processCommonJSModules"`
	moduleRoots                      j2gort.Ext_java_util_List                                     `gson:"moduleRoots"`
	warningsAllowFile                string                                                        `gson:"warningsAllowFile"`
	hideWarningsFor                  j2gort.Ext_java_util_List                                     `gson:"hideWarningsFor"`
	angularPass                      bool                                                          `gson:"angularPass"`
	jsonStreamMode                   CompilerOptions_JsonStreamMode                                `gson:"jsonStreamMode"`
	errorFormat                      AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption `gson:"errorFormat"`
	jsonWarningsFile                 string                                                        `gson:"jsonWarningsFile"`
	includeTrailingNewline           bool                                                          `gson:"includeTrailingNewline"`
	idHash_                          atomic.Int32
}

func (this_ *AbstractCommandLineRunner_CommandLineConfig) BindSelf_(s any) {
	_ = s
}
func (this_ *AbstractCommandLineRunner_CommandLineConfig) IsJavaNil_() bool { return this_ == nil }
func (*AbstractCommandLineRunner_CommandLineConfig) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractCommandLineRunner$CommandLineConfig", "CommandLineConfig"
}
func (this_ *AbstractCommandLineRunner_CommandLineConfig) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AbstractCommandLineRunner_CommandLineConfig) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractCommandLineRunner.CommandLineConfig", "com.google.javascript.jscomp.AbstractCommandLineRunner$CommandLineConfig", func(v_ any) bool {
		return j2gort.InstanceOf[*AbstractCommandLineRunner_CommandLineConfig](v_)
	})
}

func (this_ *AbstractCommandLineRunner_CommandLineConfig) initFields_AbstractCommandLineRunner_CommandLineConfig() {
	_ = this_

	this_.jscompDevMode = StaticGet_CompilerOptions_DevMode_OFF()
	this_.loggingLevel = j2gort.Ext_java_util_logging_Level_F_WARNING.GetName()
	this_.externs = []string{}
	this_.mixedJsSources = j2gort.NewExt_java_util_ArrayList()

	this_.continueSavedCompilationFileName = j2gort.NullString
	this_.restoredCompilationStage = int32(-1)
	this_.saveAfterCompilationStage = int32(-1)
	this_.typedAstListInputFilename = j2gort.NullString
	this_.saveCompilationStateToFilename = j2gort.NullString
	this_.chunk = []string{}
	this_.sourceMapInputFiles = j2gort.NewExt_java_util_LinkedHashMap()

	this_.expectedDiagnostics = j2gort.NewExt_java_util_ArrayList()

	this_.codingConvention = CodingConventions_GetDefault()
	this_.summaryDetailLevel = int32(1)

	this_.chunkWrapper = j2gort.NewExt_java_util_ArrayList()

	this_.chunkOutputFiles = j2gort.NewExt_java_util_ArrayList()
	this_.chunkConformanceFiles = j2gort.NewExt_java_util_ArrayList()

	this_.sourceMapDetailLevel = StaticGet_SourceMap_DetailLevel_ALL()
	this_.sourceMapFormat = StaticGet_SourceMap_Format_DEFAULT()
	this_.sourceMapLocationMappings = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()

	this_.warningGuards = j2gort.NewExt_java_util_ArrayList()
	this_.define = j2gort.NewExt_java_util_ArrayList()

	this_.tweakProcessing = StaticGet_CompilerOptions_TweakProcessing_OFF()

	this_.outputManifests = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	this_.outputChunkDependencies = j2gort.NullString
	this_.outputBundles = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())

	this_.manifestMaps = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())

	this_.moduleRoots = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfStr1("./"))

	this_.hideWarningsFor = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())

	this_.jsonStreamMode = StaticGet_CompilerOptions_JsonStreamMode_NONE()
	this_.errorFormat = StaticGet_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_STANDARD()

	this_.includeTrailingNewline = true
}

// Java constructor: CommandLineConfig()
func NewAbstractCommandLineRunner_CommandLineConfig() *AbstractCommandLineRunner_CommandLineConfig {
	this_ := &AbstractCommandLineRunner_CommandLineConfig{}
	InitAbstractCommandLineRunner_CommandLineConfig(this_)
	return this_
}

func InitAbstractCommandLineRunner_CommandLineConfig(this_ *AbstractCommandLineRunner_CommandLineConfig) {
	this_.initFields_AbstractCommandLineRunner_CommandLineConfig()
}

// Java method: CommandLineConfig.setPrintVersion(boolean x)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetPrintVersion(x bool) *AbstractCommandLineRunner_CommandLineConfig {
	_ = x
	this_.printVersion = x
	return this_
}

// Java method: CommandLineConfig.setPrintTree(boolean printTree)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetPrintTree(printTree bool) *AbstractCommandLineRunner_CommandLineConfig {
	_ = printTree
	this_.printTree = printTree
	return this_
}

// Java method: CommandLineConfig.setPrintTreeJson(boolean printTreeJson)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetPrintTreeJson(printTreeJson bool) *AbstractCommandLineRunner_CommandLineConfig {
	_ = printTreeJson
	this_.printTreeJson = printTreeJson
	return this_
}

// Java method: CommandLineConfig.setPrintAst(boolean printAst)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetPrintAst(printAst bool) *AbstractCommandLineRunner_CommandLineConfig {
	_ = printAst
	this_.printAst = printAst
	return this_
}

// Java method: CommandLineConfig.setJscompDevMode(com.google.javascript.jscomp.CompilerOptions.DevMode jscompDevMode)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetJscompDevMode(jscompDevMode CompilerOptions_DevMode) *AbstractCommandLineRunner_CommandLineConfig {
	_ = jscompDevMode
	this_.jscompDevMode = jscompDevMode
	return this_
}

// Java method: CommandLineConfig.setLoggingLevel(java.lang.String loggingLevel)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetLoggingLevel(loggingLevel string) *AbstractCommandLineRunner_CommandLineConfig {
	_ = loggingLevel
	this_.loggingLevel = loggingLevel
	return this_
}

// Java method: CommandLineConfig.setExterns(java.util.List<java.lang.String> externs)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetExterns(externs j2gort.Ext_java_util_List) *AbstractCommandLineRunner_CommandLineConfig {
	_ = externs
	clear(this_.externs)
	this_.externs = this_.externs[:0]
	this_.externs = j2gort.AppendStrCol(this_.externs, j2gort.AsIface[j2gort.Ext_java_util_Collection](externs))
	return this_
}

// Java method: CommandLineConfig.setMixedJsSources(java.util.List<com.google.javascript.jscomp.AbstractCommandLineRunner.FlagEntry<com.google.javascript.jscomp.AbstractCommandLineRunner.JsSourceType>> mixedJsSources)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetMixedJsSources(mixedJsSources j2gort.Ext_java_util_List) *AbstractCommandLineRunner_CommandLineConfig {
	_ = mixedJsSources
	this_.mixedJsSources.Clear()
	this_.mixedJsSources.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](mixedJsSources))
	return this_
}

// Java method: CommandLineConfig.setDefaultToStdin()
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetDefaultToStdin() *AbstractCommandLineRunner_CommandLineConfig {
	this_.defaultToStdin = true
	return this_
}

// Java method: CommandLineConfig.setJsOutputFile(java.lang.String jsOutputFile)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetJsOutputFile(jsOutputFile string) *AbstractCommandLineRunner_CommandLineConfig {
	_ = jsOutputFile
	this_.jsOutputFile = jsOutputFile
	return this_
}

// Java method: CommandLineConfig.setContinueSavedCompilationFileName(java.lang.String fileName, int restoredCompilationStage)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetContinueSavedCompilationFileName(fileName string, restoredCompilationStage int32) *AbstractCommandLineRunner_CommandLineConfig {
	_ = fileName
	_ = restoredCompilationStage
	if fileName != j2gort.NullString {
		if !((restoredCompilationStage > int32(0)) && (restoredCompilationStage < int32(3))) {
			j2gort.CheckArgumentFmt(false, "invalid compilation stage: %s", j2gort.ValueOfInt(restoredCompilationStage))
		}
		this_.continueSavedCompilationFileName = fileName
		this_.restoredCompilationStage = restoredCompilationStage
	}
	return this_
}

// Java method: CommandLineConfig.shouldSaveAfterStage1()
func (this_ *AbstractCommandLineRunner_CommandLineConfig) ShouldSaveAfterStage1() bool {
	if (this_.saveCompilationStateToFilename != j2gort.NullString) && (this_.saveAfterCompilationStage == int32(1)) {
		if !(this_.restoredCompilationStage < int32(0)) {
			j2gort.CheckStateFmt(false, "cannot perform stage 1 after restoring from stage %s", j2gort.ValueOfInt(this_.restoredCompilationStage))
		}
		if !(this_.continueSavedCompilationFileName == j2gort.NullString) {
			j2gort.CheckStateStr(false, "cannot restore a saved compilation and also save after stage 1")
		}
		return true
	} else {
		return false
	}
}

// Java method: CommandLineConfig.getContinueSavedCompilationFileName()
func (this_ *AbstractCommandLineRunner_CommandLineConfig) GetContinueSavedCompilationFileName() string {
	return this_.continueSavedCompilationFileName
}

// Java method: CommandLineConfig.shouldRestoreAndPerformStages2And3()
func (this_ *AbstractCommandLineRunner_CommandLineConfig) ShouldRestoreAndPerformStages2And3() bool {
	return ((this_.priv_CommandLineConfig_shouldContinueCompilation_225() && (this_.restoredCompilationStage == int32(1))) && (this_.saveAfterCompilationStage < int32(0)))
}

// Java method: CommandLineConfig.shouldContinueCompilation()
func (this_ *AbstractCommandLineRunner_CommandLineConfig) priv_CommandLineConfig_shouldContinueCompilation_225() bool {
	if this_.continueSavedCompilationFileName != j2gort.NullString {
		if !((this_.restoredCompilationStage > int32(0)) && (this_.restoredCompilationStage < int32(3))) {
			j2gort.CheckStateFmt(false, "invalid restored compilation stage: %s", j2gort.ValueOfInt(this_.restoredCompilationStage))
		}
		return true
	} else {
		return false
	}
}

// Java method: CommandLineConfig.shouldRestoreTypedAstsPerformStage2AndSave()
func (this_ *AbstractCommandLineRunner_CommandLineConfig) ShouldRestoreTypedAstsPerformStage2AndSave() bool {
	return ((this_.typedAstListInputFilename != j2gort.NullString) && (this_.saveAfterCompilationStage == int32(2)))
}

// Java method: CommandLineConfig.shouldRestoreTypedAstsPerformStages2And3()
func (this_ *AbstractCommandLineRunner_CommandLineConfig) ShouldRestoreTypedAstsPerformStages2And3() bool {
	return ((this_.typedAstListInputFilename != j2gort.NullString) && (this_.saveAfterCompilationStage == int32(-1)))
}

// Java method: CommandLineConfig.shouldRestoreAndPerformStage2AndSave()
func (this_ *AbstractCommandLineRunner_CommandLineConfig) ShouldRestoreAndPerformStage2AndSave() bool {
	return ((this_.priv_CommandLineConfig_shouldContinueCompilation_225() && (this_.restoredCompilationStage == int32(1))) && (this_.saveAfterCompilationStage == int32(2)))
}

// Java method: CommandLineConfig.shouldRestoreAndPerformStage3()
func (this_ *AbstractCommandLineRunner_CommandLineConfig) ShouldRestoreAndPerformStage3() bool {
	if this_.priv_CommandLineConfig_shouldContinueCompilation_225() && (this_.restoredCompilationStage == int32(2)) {
		if !(this_.saveAfterCompilationStage < int32(0)) {
			j2gort.CheckStateFmt(false, "request to save after stage %s is invalid when restoring from stage 2", j2gort.ValueOfInt(this_.saveAfterCompilationStage))
		}
		return true
	} else {
		return false
	}
}

// Java method: CommandLineConfig.shouldDoLateLocalization()
func (this_ *AbstractCommandLineRunner_CommandLineConfig) ShouldDoLateLocalization() bool {
	return ((this_.ShouldRestoreAndPerformStage2AndSave() || this_.ShouldRestoreAndPerformStage3()) || this_.ShouldRestoreTypedAstsPerformStage2AndSave())
}

// Java method: CommandLineConfig.shouldAlwaysGatherSourceMapInfo()
func (this_ *AbstractCommandLineRunner_CommandLineConfig) ShouldAlwaysGatherSourceMapInfo() bool {
	return ((this_.ShouldSaveAfterStage1() || this_.ShouldRestoreAndPerformStage2AndSave()) || this_.ShouldRestoreTypedAstsPerformStage2AndSave())
}

// Java method: CommandLineConfig.setTypedAstListInputFilename(java.lang.@org.jspecify.annotations.Nullable String fileName)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetTypedAstListInputFilename(fileName string) *AbstractCommandLineRunner_CommandLineConfig {
	_ = fileName
	this_.typedAstListInputFilename = fileName
	return this_
}

// Java method: CommandLineConfig.setSaveCompilationStateToFilename(java.lang.String fileName, int saveAfterCompilationStage)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetSaveCompilationStateToFilename(fileName string, saveAfterCompilationStage int32) *AbstractCommandLineRunner_CommandLineConfig {
	_ = fileName
	_ = saveAfterCompilationStage
	this_.saveCompilationStateToFilename = fileName
	this_.saveAfterCompilationStage = saveAfterCompilationStage
	return this_
}

// Java method: CommandLineConfig.getSaveCompilationStateToFilename()
func (this_ *AbstractCommandLineRunner_CommandLineConfig) GetSaveCompilationStateToFilename() string {
	return this_.saveCompilationStateToFilename
}

// Java method: CommandLineConfig.setChunk(java.util.List<java.lang.String> chunk)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetChunk(chunk j2gort.Ext_java_util_List) *AbstractCommandLineRunner_CommandLineConfig {
	_ = chunk
	clear(this_.chunk)
	this_.chunk = this_.chunk[:0]
	this_.chunk = j2gort.AppendStrCol(this_.chunk, j2gort.AsIface[j2gort.Ext_java_util_Collection](chunk))
	return this_
}

// Java method: CommandLineConfig.setSourceMapInputFiles(java.util.Map<java.lang.String,java.lang.String> sourceMapInputFiles)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetSourceMapInputFiles(sourceMapInputFiles j2gort.Ext_java_util_Map) *AbstractCommandLineRunner_CommandLineConfig {
	_ = sourceMapInputFiles
	this_.sourceMapInputFiles = sourceMapInputFiles
	return this_
}

// Java method: CommandLineConfig.setParseInlineSourceMaps(boolean parseInlineSourceMaps)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetParseInlineSourceMaps(parseInlineSourceMaps bool) *AbstractCommandLineRunner_CommandLineConfig {
	_ = parseInlineSourceMaps
	this_.parseInlineSourceMaps = parseInlineSourceMaps
	return this_
}

// Java method: CommandLineConfig.setExpectedDiagnostics(java.util.List<java.lang.String> expectedDiagnostics)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetExpectedDiagnostics(expectedDiagnostics j2gort.Ext_java_util_List) *AbstractCommandLineRunner_CommandLineConfig {
	_ = expectedDiagnostics
	this_.expectedDiagnostics.Clear()
	this_.expectedDiagnostics.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](expectedDiagnostics))
	return this_
}

// Java method: CommandLineConfig.setVariableMapInputFile(java.lang.String variableMapInputFile)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetVariableMapInputFile(variableMapInputFile string) *AbstractCommandLineRunner_CommandLineConfig {
	_ = variableMapInputFile
	this_.variableMapInputFile = variableMapInputFile
	return this_
}

// Java method: CommandLineConfig.setPropertyMapInputFile(java.lang.String propertyMapInputFile)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetPropertyMapInputFile(propertyMapInputFile string) *AbstractCommandLineRunner_CommandLineConfig {
	_ = propertyMapInputFile
	this_.propertyMapInputFile = propertyMapInputFile
	return this_
}

// Java method: CommandLineConfig.setVariableMapOutputFile(java.lang.String variableMapOutputFile)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetVariableMapOutputFile(variableMapOutputFile string) *AbstractCommandLineRunner_CommandLineConfig {
	_ = variableMapOutputFile
	this_.variableMapOutputFile = variableMapOutputFile
	return this_
}

// Java method: CommandLineConfig.setCreateNameMapFiles(boolean createNameMapFiles)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetCreateNameMapFiles(createNameMapFiles bool) *AbstractCommandLineRunner_CommandLineConfig {
	_ = createNameMapFiles
	this_.createNameMapFiles = createNameMapFiles
	return this_
}

// Java method: CommandLineConfig.setPropertyMapOutputFile(java.lang.String propertyMapOutputFile)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetPropertyMapOutputFile(propertyMapOutputFile string) *AbstractCommandLineRunner_CommandLineConfig {
	_ = propertyMapOutputFile
	this_.propertyMapOutputFile = propertyMapOutputFile
	return this_
}

// Java method: CommandLineConfig.setStringMapOutputFile(java.lang.String stringMapOutputPath)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetStringMapOutputFile(stringMapOutputPath string) *AbstractCommandLineRunner_CommandLineConfig {
	_ = stringMapOutputPath
	this_.stringMapOutputPath = stringMapOutputPath
	return this_
}

// Java method: CommandLineConfig.setInstrumentationMappingFile(java.lang.String instrumentationMappingFile)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetInstrumentationMappingFile(instrumentationMappingFile string) *AbstractCommandLineRunner_CommandLineConfig {
	_ = instrumentationMappingFile
	this_.instrumentationMappingFile = instrumentationMappingFile
	return this_
}

// Java method: CommandLineConfig.setCodingConvention(com.google.javascript.jscomp.CodingConvention codingConvention)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetCodingConvention(codingConvention CodingConvention) *AbstractCommandLineRunner_CommandLineConfig {
	_ = codingConvention
	this_.codingConvention = codingConvention
	return this_
}

// Java method: CommandLineConfig.setSummaryDetailLevel(int summaryDetailLevel)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetSummaryDetailLevel(summaryDetailLevel int32) *AbstractCommandLineRunner_CommandLineConfig {
	_ = summaryDetailLevel
	this_.summaryDetailLevel = summaryDetailLevel
	return this_
}

// Java method: CommandLineConfig.setOutputWrapper(java.lang.String outputWrapper)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetOutputWrapper(outputWrapper string) *AbstractCommandLineRunner_CommandLineConfig {
	_ = outputWrapper
	this_.outputWrapper = outputWrapper
	return this_
}

// Java method: CommandLineConfig.setChunkWrapper(java.util.List<java.lang.String> chunkWrapper)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetChunkWrapper(chunkWrapper j2gort.Ext_java_util_List) *AbstractCommandLineRunner_CommandLineConfig {
	_ = chunkWrapper
	this_.chunkWrapper.Clear()
	this_.chunkWrapper.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](chunkWrapper))
	return this_
}

// Java method: CommandLineConfig.setChunkOutputPathPrefix(java.lang.String chunkOutputPathPrefix)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetChunkOutputPathPrefix(chunkOutputPathPrefix string) *AbstractCommandLineRunner_CommandLineConfig {
	_ = chunkOutputPathPrefix
	this_.chunkOutputPathPrefix = chunkOutputPathPrefix
	return this_
}

// Java method: CommandLineConfig.setChunkOutputFiles(java.util.List<java.lang.String> chunkOutputFiles)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetChunkOutputFiles(chunkOutputFiles j2gort.Ext_java_util_List) *AbstractCommandLineRunner_CommandLineConfig {
	_ = chunkOutputFiles
	this_.chunkOutputFiles.Clear()
	this_.chunkOutputFiles.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](chunkOutputFiles))
	return this_
}

// Java method: CommandLineConfig.setChunkConformanceFiles(java.util.List<java.lang.String> chunkConformanceFiles)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetChunkConformanceFiles(chunkConformanceFiles j2gort.Ext_java_util_List) *AbstractCommandLineRunner_CommandLineConfig {
	_ = chunkConformanceFiles
	this_.chunkConformanceFiles.Clear()
	this_.chunkConformanceFiles.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](chunkConformanceFiles))
	return this_
}

// Java method: CommandLineConfig.setCreateSourceMap(java.lang.String createSourceMap)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetCreateSourceMap(createSourceMap string) *AbstractCommandLineRunner_CommandLineConfig {
	_ = createSourceMap
	this_.createSourceMap = createSourceMap
	return this_
}

// Java method: CommandLineConfig.setSourceMapDetailLevel(com.google.javascript.jscomp.SourceMap.DetailLevel level)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetSourceMapDetailLevel(level *SourceMap_DetailLevel) *AbstractCommandLineRunner_CommandLineConfig {
	_ = level
	this_.sourceMapDetailLevel = level
	return this_
}

// Java method: CommandLineConfig.setSourceMapFormat(com.google.javascript.jscomp.SourceMap.Format format)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetSourceMapFormat(format *SourceMap_Format) *AbstractCommandLineRunner_CommandLineConfig {
	_ = format
	this_.sourceMapFormat = format
	return this_
}

// Java method: CommandLineConfig.setSourceMapLocationMappings(java.util.List<com.google.javascript.jscomp.SourceMap.LocationMapping> locationMappings)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetSourceMapLocationMappings(locationMappings j2gort.Ext_java_util_List) *AbstractCommandLineRunner_CommandLineConfig {
	_ = locationMappings
	this_.sourceMapLocationMappings = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](locationMappings))
	return this_
}

// Java method: CommandLineConfig.setApplyInputSourceMaps(boolean applyInputSourceMaps)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetApplyInputSourceMaps(applyInputSourceMaps bool) *AbstractCommandLineRunner_CommandLineConfig {
	_ = applyInputSourceMaps
	this_.applyInputSourceMaps = applyInputSourceMaps
	return this_
}

// Java method: CommandLineConfig.setWarningGuards(java.util.List<com.google.javascript.jscomp.AbstractCommandLineRunner.FlagEntry<com.google.javascript.jscomp.CheckLevel>> warningGuards)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetWarningGuards(warningGuards j2gort.Ext_java_util_List) *AbstractCommandLineRunner_CommandLineConfig {
	_ = warningGuards
	this_.warningGuards.Clear()
	this_.warningGuards.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](warningGuards))
	return this_
}

// Java method: CommandLineConfig.setDefine(java.util.List<java.lang.String> define)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetDefine(define j2gort.Ext_java_util_List) *AbstractCommandLineRunner_CommandLineConfig {
	_ = define
	this_.define.Clear()
	this_.define.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](define))
	return this_
}

// Java method: CommandLineConfig.setBrowserFeaturesetYear(java.lang.Integer browserFeaturesetYear)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetBrowserFeaturesetYear(browserFeaturesetYear *j2gort.BoxInt) *AbstractCommandLineRunner_CommandLineConfig {
	_ = browserFeaturesetYear
	this_.browserFeaturesetYear = j2gort.UnboxInt(browserFeaturesetYear)
	return this_
}

// Java method: CommandLineConfig.setTweakProcessing(com.google.javascript.jscomp.CompilerOptions.TweakProcessing tweakProcessing)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetTweakProcessing(tweakProcessing CompilerOptions_TweakProcessing) *AbstractCommandLineRunner_CommandLineConfig {
	_ = tweakProcessing
	this_.tweakProcessing = tweakProcessing
	return this_
}

// Java method: CommandLineConfig.setCharset(java.lang.String charset)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetCharset(charset string) *AbstractCommandLineRunner_CommandLineConfig {
	_ = charset
	this_.charset = charset
	return this_
}

// Java method: CommandLineConfig.setDependencyOptions(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable DependencyOptions dependencyOptions)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetDependencyOptions(dependencyOptions *DependencyOptions) *AbstractCommandLineRunner_CommandLineConfig {
	_ = dependencyOptions
	this_.dependencyOptions = dependencyOptions
	return this_
}

// Java method: CommandLineConfig.setOutputManifest(java.util.List<java.lang.String> outputManifests)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetOutputManifest(outputManifests j2gort.Ext_java_util_List) *AbstractCommandLineRunner_CommandLineConfig {
	_ = outputManifests
	this_.outputManifests = j2gort.NewExt_java_util_ArrayList()
	{
		cur_101 := j2gort.NewIterCursor(j2gort.CheckNotNull(outputManifests))
		for cur_101.HasNext() {
			manifestName := cur_101.NextStr()
			_ = manifestName
			if !(j2gort.CheckNotNullString(manifestName) == "") {
				this_.outputManifests.Add_1(j2gort.BoxString(manifestName))
			}
		}
	}
	this_.outputManifests = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.outputManifests)))
	return this_
}

// Java method: CommandLineConfig.setOutputChunkDependencies(java.lang.String outputChunkDependencies)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetOutputChunkDependencies(outputChunkDependencies string) *AbstractCommandLineRunner_CommandLineConfig {
	_ = outputChunkDependencies
	this_.outputChunkDependencies = outputChunkDependencies
	return this_
}

// Java method: CommandLineConfig.setOutputBundle(java.util.List<java.lang.String> outputBundles)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetOutputBundle(outputBundles j2gort.Ext_java_util_List) *AbstractCommandLineRunner_CommandLineConfig {
	_ = outputBundles
	this_.outputBundles = outputBundles
	return this_
}

// Java method: CommandLineConfig.setSkipNormalOutputs(boolean skipNormalOutputs)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetSkipNormalOutputs(skipNormalOutputs bool) *AbstractCommandLineRunner_CommandLineConfig {
	_ = skipNormalOutputs
	this_.skipNormalOutputs = skipNormalOutputs
	return this_
}

// Java method: CommandLineConfig.setManifestMaps(java.util.List<java.lang.String> manifestMaps)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetManifestMaps(manifestMaps j2gort.Ext_java_util_List) *AbstractCommandLineRunner_CommandLineConfig {
	_ = manifestMaps
	this_.manifestMaps = manifestMaps
	return this_
}

// Java method: CommandLineConfig.setProcessCommonJSModules(boolean processCommonJSModules)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetProcessCommonJSModules(processCommonJSModules bool) *AbstractCommandLineRunner_CommandLineConfig {
	_ = processCommonJSModules
	this_.processCommonJSModules = processCommonJSModules
	return this_
}

// Java method: CommandLineConfig.setModuleRoots(java.util.List<java.lang.String> jsModuleRoots)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetModuleRoots(jsModuleRoots j2gort.Ext_java_util_List) *AbstractCommandLineRunner_CommandLineConfig {
	_ = jsModuleRoots
	this_.moduleRoots = jsModuleRoots
	return this_
}

// Java method: CommandLineConfig.setWarningsAllowlistFile(java.lang.String fileName)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetWarningsAllowlistFile(fileName string) *AbstractCommandLineRunner_CommandLineConfig {
	_ = fileName
	this_.warningsAllowFile = fileName
	return this_
}

// Java method: CommandLineConfig.setHideWarningsFor(java.util.List<java.lang.String> hideWarningsFor)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetHideWarningsFor(hideWarningsFor j2gort.Ext_java_util_List) *AbstractCommandLineRunner_CommandLineConfig {
	_ = hideWarningsFor
	this_.hideWarningsFor = hideWarningsFor
	return this_
}

// Java method: CommandLineConfig.setAngularPass(boolean angularPass)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetAngularPass(angularPass bool) *AbstractCommandLineRunner_CommandLineConfig {
	_ = angularPass
	this_.angularPass = angularPass
	return this_
}

// Java method: CommandLineConfig.setJsonStreamMode(com.google.javascript.jscomp.CompilerOptions.JsonStreamMode mode)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetJsonStreamMode(mode CompilerOptions_JsonStreamMode) *AbstractCommandLineRunner_CommandLineConfig {
	_ = mode
	this_.jsonStreamMode = mode
	return this_
}

// Java method: CommandLineConfig.setErrorFormat(com.google.javascript.jscomp.AbstractCommandLineRunner.CommandLineConfig.ErrorFormatOption errorFormat)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetErrorFormat(errorFormat AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) *AbstractCommandLineRunner_CommandLineConfig {
	_ = errorFormat
	this_.errorFormat = errorFormat
	return this_
}

// Java method: CommandLineConfig.setJsonWarningsFile(java.lang.String jsonWarningsFile)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetJsonWarningsFile(jsonWarningsFile string) *AbstractCommandLineRunner_CommandLineConfig {
	_ = jsonWarningsFile
	this_.jsonWarningsFile = jsonWarningsFile
	return this_
}

// Java method: CommandLineConfig.setIncludeTrailingNewline(boolean includeTrailingNewline)
func (this_ *AbstractCommandLineRunner_CommandLineConfig) SetIncludeTrailingNewline(includeTrailingNewline bool) *AbstractCommandLineRunner_CommandLineConfig {
	_ = includeTrailingNewline
	this_.includeTrailingNewline = includeTrailingNewline
	return this_
}

// Java enum: com.google.javascript.jscomp.AbstractCommandLineRunner.CommandLineConfig.ErrorFormatOption
type AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption uint8

const (
	enum_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_STANDARD AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption = 1
	enum_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_JSON     AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption = 2
)

var enumBaseTable_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption = [...]j2gort.EnumBase{
	{Name_: "STANDARD", Ordinal_: 0},
	{Name_: "JSON", Ordinal_: 1},
}
var enumValues_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption = []AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption{1, 2}

func StaticGet_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_STANDARD() AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption {
	return enum_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_STANDARD
}
func StaticGet_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_JSON() AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption {
	return enum_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_JSON
}

func AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_Values() []AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption {
	return j2gort.CloneSlice(enumValues_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption)
}

func AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption_ValueOf(name string) AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption {
		if enumBaseTable_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption[i_].Name_ == name {
			return AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.AbstractCommandLineRunner.CommandLineConfig.ErrorFormatOption." + name))
}

func (this_ AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) IsJavaNil_() bool {
	return this_ == 0
}
func (AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) HasIdentityEquals_() {}
func (AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) Implements_j2gort_Ext_java_lang_Enum() {
}
func (AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractCommandLineRunner$CommandLineConfig$ErrorFormatOption", "ErrorFormatOption"
}
func (this_ AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption[this_-1]
}
func (this_ AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) EnumOrdinal_() int32 {
	return int32(this_) - 1
}
func (this_ AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption[this_-1].Name_
}
func (this_ AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption[this_-1].Name_
}
func (this_ AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) Equals(other_ any) bool {
	o_, ok_ := other_.(AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption[this_-1].HashCode()
}
func (this_ AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption[this_-1].CompareTo(other_)
}
func (this_ AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption[this_-1].Name_)
}
func init() {
	vals_ := []any{AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption(1), AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AbstractCommandLineRunner.CommandLineConfig.ErrorFormatOption", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AbstractCommandLineRunner$CommandLineConfig$ErrorFormatOption", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractCommandLineRunner.CommandLineConfig.ErrorFormatOption", "com.google.javascript.jscomp.AbstractCommandLineRunner$CommandLineConfig$ErrorFormatOption", func(v_ any) bool {
		return j2gort.InstanceOf[AbstractCommandLineRunner_CommandLineConfig_ErrorFormatOption](v_)
	})
}

// Java class: com.google.javascript.jscomp.AbstractCommandLineRunner.JsonFileSpec
type AbstractCommandLineRunner_JsonFileSpec struct {
	src       string `gson:"src"`
	path      string `gson:"path"`
	sourceMap string `gson:"source_map"`
	webpackId string `gson:"webpack_id"`
	idHash_   atomic.Int32
}

func (this_ *AbstractCommandLineRunner_JsonFileSpec) BindSelf_(s any) {
	_ = s
}
func (this_ *AbstractCommandLineRunner_JsonFileSpec) IsJavaNil_() bool { return this_ == nil }
func (*AbstractCommandLineRunner_JsonFileSpec) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractCommandLineRunner$JsonFileSpec", "JsonFileSpec"
}
func (this_ *AbstractCommandLineRunner_JsonFileSpec) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AbstractCommandLineRunner_JsonFileSpec) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractCommandLineRunner.JsonFileSpec", "com.google.javascript.jscomp.AbstractCommandLineRunner$JsonFileSpec", func(v_ any) bool {
		return j2gort.InstanceOf[*AbstractCommandLineRunner_JsonFileSpec](v_)
	})
}

func (this_ *AbstractCommandLineRunner_JsonFileSpec) initFields_AbstractCommandLineRunner_JsonFileSpec() {
	_ = this_
	this_.src = j2gort.NullString
	this_.path = j2gort.NullString
	this_.sourceMap = j2gort.NullString
	this_.webpackId = j2gort.NullString
}

// Java constructor: JsonFileSpec(java.lang.String src, java.lang.String path)
func NewAbstractCommandLineRunner_JsonFileSpec_2(src string, path string) *AbstractCommandLineRunner_JsonFileSpec {
	this_ := &AbstractCommandLineRunner_JsonFileSpec{}
	InitAbstractCommandLineRunner_JsonFileSpec_2(this_, src, path)
	return this_
}

func InitAbstractCommandLineRunner_JsonFileSpec_2(this_ *AbstractCommandLineRunner_JsonFileSpec, src string, path string) {
	_ = src
	_ = path
	InitAbstractCommandLineRunner_JsonFileSpec_4(this_, src, path, j2gort.NullString, j2gort.NullString)
}

// Java constructor: JsonFileSpec(java.lang.String src, java.lang.String path, java.lang.String sourceMap)
func NewAbstractCommandLineRunner_JsonFileSpec_3(src string, path string, sourceMap string) *AbstractCommandLineRunner_JsonFileSpec {
	this_ := &AbstractCommandLineRunner_JsonFileSpec{}
	InitAbstractCommandLineRunner_JsonFileSpec_3(this_, src, path, sourceMap)
	return this_
}

func InitAbstractCommandLineRunner_JsonFileSpec_3(this_ *AbstractCommandLineRunner_JsonFileSpec, src string, path string, sourceMap string) {
	_ = src
	_ = path
	_ = sourceMap
	InitAbstractCommandLineRunner_JsonFileSpec_4(this_, src, path, sourceMap, j2gort.NullString)
}

// Java constructor: JsonFileSpec(java.lang.@org.jspecify.annotations.Nullable String src, java.lang.@org.jspecify.annotations.Nullable String path, java.lang.@org.jspecify.annotations.Nullable String sourceMap, java.lang.@org.jspecify.annotations.Nullable String webpackId)
func NewAbstractCommandLineRunner_JsonFileSpec_4(src string, path string, sourceMap string, webpackId string) *AbstractCommandLineRunner_JsonFileSpec {
	this_ := &AbstractCommandLineRunner_JsonFileSpec{}
	InitAbstractCommandLineRunner_JsonFileSpec_4(this_, src, path, sourceMap, webpackId)
	return this_
}

func InitAbstractCommandLineRunner_JsonFileSpec_4(this_ *AbstractCommandLineRunner_JsonFileSpec, src string, path string, sourceMap string, webpackId string) {
	_ = src
	_ = path
	_ = sourceMap
	_ = webpackId
	this_.initFields_AbstractCommandLineRunner_JsonFileSpec()
	this_.src = src
	this_.path = path
	this_.sourceMap = sourceMap
	this_.webpackId = webpackId
}

// Java static method: JsonFileSpec.readJson(com.google.gson.stream.JsonReader reader)
func AbstractCommandLineRunner_JsonFileSpec_ReadJson(reader *j2gort.Ext_com_google_gson_stream_JsonReader) *AbstractCommandLineRunner_JsonFileSpec {
	_ = reader
	reader.BeginObject()
	var src string = j2gort.NullString
	_ = src
	var path string = j2gort.NullString
	_ = path
	var sourceMap string = j2gort.NullString
	_ = sourceMap
	var webpackId string = j2gort.NullString
	_ = webpackId
	for reader.HasNext() {
		var name string = reader.NextName()
		_ = name
		switch j2gort.CheckNotNullString(name) {
		case "src":
			src = reader.NextString()
		case "path":
			path = reader.NextString()
		case "source_map", "sourceMap":
			sourceMap = reader.NextString()
		case "webpack_id", "webpackId":
			webpackId = reader.NextString()
		default:
			reader.SkipValue()
		}
	}
	reader.EndObject()
	return NewAbstractCommandLineRunner_JsonFileSpec_4(src, path, sourceMap, webpackId)
}

// Java method: JsonFileSpec.writeJson(com.google.gson.stream.JsonWriter writer)
func (this_ *AbstractCommandLineRunner_JsonFileSpec) WriteJson(writer *j2gort.Ext_com_google_gson_stream_JsonWriter) {
	_ = writer
	writer.BeginObject()
	if this_.src != j2gort.NullString {
		writer.Name("src").Value_1_string(this_.src)
	}
	if this_.path != j2gort.NullString {
		writer.Name("path").Value_1_string(this_.path)
	}
	if this_.sourceMap != j2gort.NullString {
		writer.Name("source_map").Value_1_string(this_.sourceMap)
	}
	if this_.webpackId != j2gort.NullString {
		writer.Name("webpack_id").Value_1_string(this_.webpackId)
	}
	writer.EndObject()
}

// Java method: JsonFileSpec.getSrc()
func (this_ *AbstractCommandLineRunner_JsonFileSpec) GetSrc() string {
	return this_.src
}

// Java method: JsonFileSpec.getPath()
func (this_ *AbstractCommandLineRunner_JsonFileSpec) GetPath() string {
	return this_.path
}

// Java method: JsonFileSpec.getSourceMap()
func (this_ *AbstractCommandLineRunner_JsonFileSpec) GetSourceMap() string {
	return this_.sourceMap
}

// Java method: JsonFileSpec.getWebpackId()
func (this_ *AbstractCommandLineRunner_JsonFileSpec) GetWebpackId() string {
	return this_.webpackId
}

// Java method: JsonFileSpec.setSourceMap(java.lang.String map)
func (this_ *AbstractCommandLineRunner_JsonFileSpec) SetSourceMap(map_ string) {
	_ = map_
	this_.sourceMap = map_
}

// Java enum: com.google.javascript.jscomp.AbstractCommandLineRunner.JsSourceType
type AbstractCommandLineRunner_JsSourceType uint8

const (
	enum_AbstractCommandLineRunner_JsSourceType_EXTERN  AbstractCommandLineRunner_JsSourceType = 1
	enum_AbstractCommandLineRunner_JsSourceType_JS      AbstractCommandLineRunner_JsSourceType = 2
	enum_AbstractCommandLineRunner_JsSourceType_JS_ZIP  AbstractCommandLineRunner_JsSourceType = 3
	enum_AbstractCommandLineRunner_JsSourceType_WEAKDEP AbstractCommandLineRunner_JsSourceType = 4
	enum_AbstractCommandLineRunner_JsSourceType_IJS     AbstractCommandLineRunner_JsSourceType = 5
)

var enumBaseTable_AbstractCommandLineRunner_JsSourceType = [...]j2gort.EnumBase{
	{Name_: "EXTERN", Ordinal_: 0},
	{Name_: "JS", Ordinal_: 1},
	{Name_: "JS_ZIP", Ordinal_: 2},
	{Name_: "WEAKDEP", Ordinal_: 3},
	{Name_: "IJS", Ordinal_: 4},
}
var enumValues_AbstractCommandLineRunner_JsSourceType = []AbstractCommandLineRunner_JsSourceType{1, 2, 3, 4, 5}

var enumField_AbstractCommandLineRunner_JsSourceType_flagName = [...]string{
	"extern",
	"js",
	"jszip",
	"weakdep",
	"ijs",
}

func (this_ AbstractCommandLineRunner_JsSourceType) FieldGet_flagName() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_AbstractCommandLineRunner_JsSourceType_flagName[this_-1]
}

func StaticGet_AbstractCommandLineRunner_JsSourceType_EXTERN() AbstractCommandLineRunner_JsSourceType {
	return enum_AbstractCommandLineRunner_JsSourceType_EXTERN
}
func StaticGet_AbstractCommandLineRunner_JsSourceType_JS() AbstractCommandLineRunner_JsSourceType {
	return enum_AbstractCommandLineRunner_JsSourceType_JS
}
func StaticGet_AbstractCommandLineRunner_JsSourceType_JS_ZIP() AbstractCommandLineRunner_JsSourceType {
	return enum_AbstractCommandLineRunner_JsSourceType_JS_ZIP
}
func StaticGet_AbstractCommandLineRunner_JsSourceType_WEAKDEP() AbstractCommandLineRunner_JsSourceType {
	return enum_AbstractCommandLineRunner_JsSourceType_WEAKDEP
}
func StaticGet_AbstractCommandLineRunner_JsSourceType_IJS() AbstractCommandLineRunner_JsSourceType {
	return enum_AbstractCommandLineRunner_JsSourceType_IJS
}

func AbstractCommandLineRunner_JsSourceType_Values() []AbstractCommandLineRunner_JsSourceType {
	return j2gort.CloneSlice(enumValues_AbstractCommandLineRunner_JsSourceType)
}

func AbstractCommandLineRunner_JsSourceType_ValueOf(name string) AbstractCommandLineRunner_JsSourceType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_AbstractCommandLineRunner_JsSourceType {
		if enumBaseTable_AbstractCommandLineRunner_JsSourceType[i_].Name_ == name {
			return AbstractCommandLineRunner_JsSourceType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.AbstractCommandLineRunner.JsSourceType." + name))
}

func (this_ AbstractCommandLineRunner_JsSourceType) IsJavaNil_() bool                              { return this_ == 0 }
func (AbstractCommandLineRunner_JsSourceType) HasIdentityEquals_()                                 {}
func (AbstractCommandLineRunner_JsSourceType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (AbstractCommandLineRunner_JsSourceType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ AbstractCommandLineRunner_JsSourceType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (AbstractCommandLineRunner_JsSourceType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractCommandLineRunner$JsSourceType", "JsSourceType"
}
func (this_ AbstractCommandLineRunner_JsSourceType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_AbstractCommandLineRunner_JsSourceType[this_-1]
}
func (this_ AbstractCommandLineRunner_JsSourceType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ AbstractCommandLineRunner_JsSourceType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ AbstractCommandLineRunner_JsSourceType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_AbstractCommandLineRunner_JsSourceType[this_-1].Name_
}
func (this_ AbstractCommandLineRunner_JsSourceType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_AbstractCommandLineRunner_JsSourceType[this_-1].Name_
}
func (this_ AbstractCommandLineRunner_JsSourceType) Equals(other_ any) bool {
	o_, ok_ := other_.(AbstractCommandLineRunner_JsSourceType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ AbstractCommandLineRunner_JsSourceType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_AbstractCommandLineRunner_JsSourceType[this_-1].HashCode()
}
func (this_ AbstractCommandLineRunner_JsSourceType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_AbstractCommandLineRunner_JsSourceType[this_-1].CompareTo(other_)
}
func (this_ AbstractCommandLineRunner_JsSourceType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_AbstractCommandLineRunner_JsSourceType[this_-1].Name_)
}
func init() {
	vals_ := []any{AbstractCommandLineRunner_JsSourceType(1), AbstractCommandLineRunner_JsSourceType(2), AbstractCommandLineRunner_JsSourceType(3), AbstractCommandLineRunner_JsSourceType(4), AbstractCommandLineRunner_JsSourceType(5)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AbstractCommandLineRunner.JsSourceType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AbstractCommandLineRunner$JsSourceType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractCommandLineRunner.JsSourceType", "com.google.javascript.jscomp.AbstractCommandLineRunner$JsSourceType", func(v_ any) bool {
		return j2gort.InstanceOf[AbstractCommandLineRunner_JsSourceType](v_)
	})
}

// Java class: com.google.javascript.jscomp.AbstractCommandLineRunner.FlagEntry
type AbstractCommandLineRunner_FlagEntry struct {
	flag  any    `gson:"flag"`
	value string `gson:"value"`
}

func (this_ *AbstractCommandLineRunner_FlagEntry) BindSelf_(s any) {
	_ = s
}
func (this_ *AbstractCommandLineRunner_FlagEntry) IsJavaNil_() bool { return this_ == nil }
func (*AbstractCommandLineRunner_FlagEntry) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractCommandLineRunner$FlagEntry", "FlagEntry"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractCommandLineRunner.FlagEntry", "com.google.javascript.jscomp.AbstractCommandLineRunner$FlagEntry", func(v_ any) bool {
		return j2gort.InstanceOf[*AbstractCommandLineRunner_FlagEntry](v_)
	})
}

func (this_ *AbstractCommandLineRunner_FlagEntry) initFields_AbstractCommandLineRunner_FlagEntry() {
	_ = this_
	this_.value = j2gort.NullString
}

// Java constructor: FlagEntry(T flag, java.lang.String value)
func NewAbstractCommandLineRunner_FlagEntry(flag any, value string) *AbstractCommandLineRunner_FlagEntry {
	this_ := &AbstractCommandLineRunner_FlagEntry{}
	InitAbstractCommandLineRunner_FlagEntry(this_, flag, value)
	return this_
}

func InitAbstractCommandLineRunner_FlagEntry(this_ *AbstractCommandLineRunner_FlagEntry, flag any, value string) {
	_ = flag
	_ = value
	this_.initFields_AbstractCommandLineRunner_FlagEntry()
	this_.flag = flag
	this_.value = value
}

// Java method: FlagEntry.equals(java.lang.Object o)
func (this_ *AbstractCommandLineRunner_FlagEntry) Equals(arg_102 any) bool {
	o := arg_102
	_ = o
	if j2gort.InstanceOf[*AbstractCommandLineRunner_FlagEntry](o) {
		var that *AbstractCommandLineRunner_FlagEntry = j2gort.CastPtr[AbstractCommandLineRunner_FlagEntry](o)
		_ = that
		return (javalang.ObjectsEquals(j2gort.CheckNotNull(that.flag), this_.flag) && func() bool {
			sr_103 := that.value
			se_104 := this_.value
			return j2gort.CheckNotNullString(sr_103) == se_104
		}())
	}
	return false
}

// Java method: FlagEntry.hashCode()
func (this_ *AbstractCommandLineRunner_FlagEntry) HashCode() int32 {
	return (javalang.ObjectHashCode(j2gort.CheckNotNull(this_.flag)) + javalang.StringHashCode(j2gort.CheckNotNullString(this_.value)))
}

// Java method: FlagEntry.toString()
func (this_ *AbstractCommandLineRunner_FlagEntry) String() string {
	return (j2gort.StrVal(this_.flag) + "=" + j2gort.StrValString(this_.value))
}

// Java method: FlagEntry.getFlag()
func (this_ *AbstractCommandLineRunner_FlagEntry) GetFlag() any {
	return this_.flag
}

// Java method: FlagEntry.getValue()
func (this_ *AbstractCommandLineRunner_FlagEntry) GetValue() string {
	return this_.value
}

// Java class: com.google.javascript.jscomp.AbstractCommandLineRunner.JsChunkSpec
type AbstractCommandLineRunner_JsChunkSpec struct {
	name       string   `gson:"name"`
	numInputs  int32    `gson:"numInputs"`
	deps       []string `gson:"deps"`
	numJsFiles int32    `gson:"numJsFiles"`
	idHash_    atomic.Int32
}

func (this_ *AbstractCommandLineRunner_JsChunkSpec) BindSelf_(s any) {
	_ = s
}
func (this_ *AbstractCommandLineRunner_JsChunkSpec) IsJavaNil_() bool { return this_ == nil }
func (*AbstractCommandLineRunner_JsChunkSpec) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractCommandLineRunner$JsChunkSpec", "JsChunkSpec"
}
func (this_ *AbstractCommandLineRunner_JsChunkSpec) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AbstractCommandLineRunner_JsChunkSpec) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractCommandLineRunner.JsChunkSpec", "com.google.javascript.jscomp.AbstractCommandLineRunner$JsChunkSpec", func(v_ any) bool {
		return j2gort.InstanceOf[*AbstractCommandLineRunner_JsChunkSpec](v_)
	})
}

func (this_ *AbstractCommandLineRunner_JsChunkSpec) initFields_AbstractCommandLineRunner_JsChunkSpec() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: JsChunkSpec(java.lang.String name, int numInputs, com.google.common.collect.ImmutableList<java.lang.String> deps)
func NewAbstractCommandLineRunner_JsChunkSpec(name string, numInputs int32, deps j2gort.Ext_com_google_common_collect_ImmutableList) *AbstractCommandLineRunner_JsChunkSpec {
	this_ := &AbstractCommandLineRunner_JsChunkSpec{}
	InitAbstractCommandLineRunner_JsChunkSpec(this_, name, numInputs, deps)
	return this_
}

func InitAbstractCommandLineRunner_JsChunkSpec(this_ *AbstractCommandLineRunner_JsChunkSpec, name string, numInputs int32, deps j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = name
	_ = numInputs
	_ = deps
	this_.initFields_AbstractCommandLineRunner_JsChunkSpec()
	this_.name = name
	this_.numInputs = numInputs
	this_.deps = j2gort.ImmListToSliceNullable[string](deps)
	this_.numJsFiles = numInputs
}

// Java static method: JsChunkSpec.create(java.lang.String specString, boolean isFirstChunk)
func AbstractCommandLineRunner_JsChunkSpec_Create(specString string, isFirstChunk bool) *AbstractCommandLineRunner_JsChunkSpec {
	_ = specString
	_ = isFirstChunk
	var parts []string = j2gort.JStrSplit(specString, ":")
	_ = parts
	if (int32(len(parts)) < int32(2)) || (int32(len(parts)) > int32(4)) {
		panic(NewFlagUsageException(("Expected 2-4 colon-delimited parts in chunk spec: " + j2gort.StrValString(specString))))
	}
	var name string = parts[int32(0)]
	_ = name
	var deps []string
	_ = deps
	if (int32(len(parts)) > int32(2)) && (j2gort.JStrLen(parts[int32(2)]) > int32(0)) {
		deps = j2gort.JStrSplit(parts[int32(2)], ",")
	} else {
		deps = j2gort.NewStringArray(int32(0))
	}
	var numInputs int32 = int32(-1)
	_ = numInputs
	func() {
		defer func() {
			if rec_105 := recover(); rec_105 != nil {
				jex_106 := javalang.RecoverException(rec_105)
				_ = jex_106
				if j2gort.InstanceOf[javalang.INumberFormatException](jex_106) {
					ignored := j2gort.AsIface[javalang.INumberFormatException](jex_106)
					_ = ignored
					numInputs = int32(-1)
				} else {
					panic(rec_105)
				}
			}
		}()
		numInputs = j2gort.ParseInt(parts[int32(1)], 10)
	}()
	if numInputs < int32(0) {
		if (int32(len(parts)) == int32(2)) && func() bool {
			sr_107 := "auto"
			se_108 := parts[int32(1)]
			return j2gort.CheckNotNullString(sr_107) == se_108
		}() {
			if !(isFirstChunk) {
				panic(NewFlagUsageException(("Invalid JS file count '" + j2gort.StrValString(parts[int32(1)]) + "' for chunk: " + j2gort.StrValString(name) + ". Only the first chunk may specify a size of 'auto' and it must have no dependencies.")))
			}
		} else {
			panic(NewFlagUsageException(("Invalid JS file count '" + j2gort.StrValString(parts[int32(1)]) + "' for chunk: " + j2gort.StrValString(name))))
		}
	}
	return NewAbstractCommandLineRunner_JsChunkSpec(name, numInputs, j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_slice_any(func() []any {
		arr_109 := deps
		if arr_109 == nil {
			return nil
		}
		out_110 := make([]any, len(arr_109))
		for i_, el_111 := range arr_109 {
			out_110[i_] = j2gort.BoxString(el_111)
		}
		return out_110
	}()))
}

// Java method: JsChunkSpec.getName()
func (this_ *AbstractCommandLineRunner_JsChunkSpec) GetName() string {
	return this_.name
}

// Java method: JsChunkSpec.getNumInputs()
func (this_ *AbstractCommandLineRunner_JsChunkSpec) GetNumInputs() int32 {
	return this_.numInputs
}

// Java method: JsChunkSpec.getDeps()
func (this_ *AbstractCommandLineRunner_JsChunkSpec) GetDeps() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.ImmListFromTypedSlice(this_.deps)
}

// Java method: JsChunkSpec.getNumJsFiles()
func (this_ *AbstractCommandLineRunner_JsChunkSpec) GetNumJsFiles() int32 {
	return this_.numJsFiles
}

// Java class: com.google.javascript.jscomp.AbstractCommandLineRunner.SystemExitCodeReceiver
type AbstractCommandLineRunner_SystemExitCodeReceiver struct {
	idHash_ atomic.Int32
}

func (this_ *AbstractCommandLineRunner_SystemExitCodeReceiver) BindSelf_(s any) {
	_ = s
}
func (this_ *AbstractCommandLineRunner_SystemExitCodeReceiver) IsJavaNil_() bool { return this_ == nil }
func (*AbstractCommandLineRunner_SystemExitCodeReceiver) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractCommandLineRunner$SystemExitCodeReceiver", "SystemExitCodeReceiver"
}
func (this_ *AbstractCommandLineRunner_SystemExitCodeReceiver) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AbstractCommandLineRunner_SystemExitCodeReceiver) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractCommandLineRunner.SystemExitCodeReceiver", "com.google.javascript.jscomp.AbstractCommandLineRunner$SystemExitCodeReceiver", func(v_ any) bool {
		return j2gort.InstanceOf[*AbstractCommandLineRunner_SystemExitCodeReceiver](v_)
	})
}
func (*AbstractCommandLineRunner_SystemExitCodeReceiver) Implements_j2gort_Ext_com_google_common_base_Function() {
}
func (*AbstractCommandLineRunner_SystemExitCodeReceiver) Implements_j2gort_Ext_java_util_function_Function() {
}

var static_AbstractCommandLineRunner_SystemExitCodeReceiver_INSTANCE *AbstractCommandLineRunner_SystemExitCodeReceiver

func StaticGet_AbstractCommandLineRunner_SystemExitCodeReceiver_INSTANCE() *AbstractCommandLineRunner_SystemExitCodeReceiver {
	clinit_AbstractCommandLineRunner_SystemExitCodeReceiver()
	return static_AbstractCommandLineRunner_SystemExitCodeReceiver_INSTANCE
}
func StaticSet_AbstractCommandLineRunner_SystemExitCodeReceiver_INSTANCE(v *AbstractCommandLineRunner_SystemExitCodeReceiver) {
	clinit_AbstractCommandLineRunner_SystemExitCodeReceiver()
	static_AbstractCommandLineRunner_SystemExitCodeReceiver_INSTANCE = v
}

var clinitDone_AbstractCommandLineRunner_SystemExitCodeReceiver uint32
var clinitState_AbstractCommandLineRunner_SystemExitCodeReceiver j2gort.ClinitState

func clinit_AbstractCommandLineRunner_SystemExitCodeReceiver() {
	if atomic.LoadUint32(&clinitDone_AbstractCommandLineRunner_SystemExitCodeReceiver) != 0 {
		return
	}
	clinitSlow_AbstractCommandLineRunner_SystemExitCodeReceiver()
}

func clinitSlow_AbstractCommandLineRunner_SystemExitCodeReceiver() {
	if !j2gort.EnterClinit(&clinitState_AbstractCommandLineRunner_SystemExitCodeReceiver, &clinitDone_AbstractCommandLineRunner_SystemExitCodeReceiver) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_AbstractCommandLineRunner_SystemExitCodeReceiver, &clinitDone_AbstractCommandLineRunner_SystemExitCodeReceiver)
	static_AbstractCommandLineRunner_SystemExitCodeReceiver_INSTANCE = NewAbstractCommandLineRunner_SystemExitCodeReceiver()
}

func (this_ *AbstractCommandLineRunner_SystemExitCodeReceiver) initFields_AbstractCommandLineRunner_SystemExitCodeReceiver() {
	_ = this_
}

// Java constructor: SystemExitCodeReceiver()
func NewAbstractCommandLineRunner_SystemExitCodeReceiver() *AbstractCommandLineRunner_SystemExitCodeReceiver {
	clinit_AbstractCommandLineRunner_SystemExitCodeReceiver()
	this_ := &AbstractCommandLineRunner_SystemExitCodeReceiver{}
	InitAbstractCommandLineRunner_SystemExitCodeReceiver(this_)
	return this_
}

func InitAbstractCommandLineRunner_SystemExitCodeReceiver(this_ *AbstractCommandLineRunner_SystemExitCodeReceiver) {
	this_.initFields_AbstractCommandLineRunner_SystemExitCodeReceiver()
}

// Java method: SystemExitCodeReceiver.apply(java.lang.Integer exitCode)
func (this_ *AbstractCommandLineRunner_SystemExitCodeReceiver) Apply(arg_112 any) any {
	exitCode := j2gort.CastPtr[j2gort.BoxInt](arg_112)
	_ = exitCode
	var exitCodeValue int32 = j2gort.UnboxInt(j2gort.CheckNotNullPtr(exitCode))
	_ = exitCodeValue
	var exitCodeByte int8 = int8(exitCodeValue)
	_ = exitCodeByte
	if (int32(exitCodeByte) == int32(0)) && (exitCodeValue != int32(0)) {
		exitCodeByte = int8(-1)
	}
	j2gort.Ext_java_lang_System_Exit(int32(exitCodeByte))
	return nil
}

func (this_ *AbstractCommandLineRunner_SystemExitCodeReceiver) Equals(p0_ any) bool {
	return j2gort.RefEquals(this_, p0_)
}
func (this_ *AbstractCommandLineRunner_JsChunkSpec) Slice_GetDeps() []string {
	return this_.deps
}
