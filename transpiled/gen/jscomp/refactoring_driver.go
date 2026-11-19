/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/refactoring/RefactoringDriver.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

var lambdaSingleton_refactoring_driver_1 j2gort.Ext_com_google_common_base_Function

func init() {
	lambdaSingleton_refactoring_driver_1 = j2gort.NewLambda_Ext_com_google_common_base_Function(func(lp_3 any) any {
		file := j2gort.CastString(lp_3)
		_ = file
		return j2gort.PtrToAny(SourceFile_S_Builder().WithPath_1_string(file).Build())
	})
}

// Java class: com.google.javascript.refactoring.RefactoringDriver
type RefactoringDriver struct {
	compiler *Compiler      `gson:"compiler"`
	rootNode pkg_rhino.Node `gson:"rootNode"`
	idHash_  atomic.Int32
}

func (this_ *RefactoringDriver) BindSelf_(s any) {
	_ = s
}
func (this_ *RefactoringDriver) IsJavaNil_() bool { return this_ == nil }
func (*RefactoringDriver) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.RefactoringDriver", "RefactoringDriver"
}
func (this_ *RefactoringDriver) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RefactoringDriver) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.RefactoringDriver", "com.google.javascript.refactoring.RefactoringDriver", func(v_ any) bool {
		return j2gort.InstanceOf[*RefactoringDriver](v_)
	})
}

func (this_ *RefactoringDriver) initFields_RefactoringDriver() {
	_ = this_
}

// Java constructor: RefactoringDriver(java.util.List<com.google.javascript.jscomp.SourceFile> inputs, java.util.List<com.google.javascript.jscomp.SourceFile> externs, com.google.javascript.jscomp.CompilerOptions compilerOptions)
func NewRefactoringDriver(inputs j2gort.Ext_java_util_List, externs j2gort.Ext_java_util_List, compilerOptions *CompilerOptions) *RefactoringDriver {
	this_ := &RefactoringDriver{}
	InitRefactoringDriver(this_, inputs, externs, compilerOptions)
	return this_
}

func InitRefactoringDriver(this_ *RefactoringDriver, inputs j2gort.Ext_java_util_List, externs j2gort.Ext_java_util_List, compilerOptions *CompilerOptions) {
	_ = inputs
	_ = externs
	_ = compilerOptions
	this_.initFields_RefactoringDriver()
	this_.compiler = RefactoringDriver_CreateCompiler(inputs, externs, compilerOptions)
	this_.rootNode = this_.compiler.GetRoot()
}

// Java method: RefactoringDriver.drive(com.google.javascript.refactoring.Scanner scanner, java.util.regex.@org.jspecify.annotations.Nullable Pattern includeFilePattern)
func (this_ *RefactoringDriver) Drive_2(scanner Scanner, includeFilePattern *j2gort.Ext_java_util_regex_Pattern) j2gort.Ext_java_util_List {
	_ = scanner
	_ = includeFilePattern
	var val_1 JsFlumeCallback
	InitJsFlumeCallback(&val_1, scanner, includeFilePattern)
	var callback *JsFlumeCallback = &val_1
	_ = callback
	NodeTraversal_Traverse(ToIface_AbstractCompiler(this_.compiler), this_.rootNode, callback)
	var fixes j2gort.Ext_java_util_List = callback.GetFixes()
	_ = fixes
	fixes.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](scanner.ProcessAllMatches(j2gort.AsIface[j2gort.Ext_java_util_Collection](callback.GetMatches()))))
	return fixes
}

// Java method: RefactoringDriver.drive(com.google.javascript.refactoring.Scanner scanner)
func (this_ *RefactoringDriver) Drive_1(scanner Scanner) j2gort.Ext_java_util_List {
	_ = scanner
	return this_.Drive_2(scanner, nil)
}

// Java method: RefactoringDriver.getCompiler()
func (this_ *RefactoringDriver) GetCompiler() *Compiler {
	return this_.compiler
}

// Java static method: RefactoringDriver.createCompiler(java.util.List<com.google.javascript.jscomp.SourceFile> inputs, java.util.List<com.google.javascript.jscomp.SourceFile> externs, com.google.javascript.jscomp.CompilerOptions compilerOptions)
func RefactoringDriver_CreateCompiler(inputs j2gort.Ext_java_util_List, externs j2gort.Ext_java_util_List, compilerOptions *CompilerOptions) *Compiler {
	_ = inputs
	_ = externs
	_ = compilerOptions
	var compiler *Compiler = NewCompiler_1_ErrorManager(NewBlackHoleErrorManager())
	_ = compiler
	compiler.DisableThreads()
	compiler.Compile_3_Ext_java_util_List_Ext_java_util_List_CompilerOptions(externs, inputs, compilerOptions)
	return compiler
}

// Java static method: RefactoringDriver.getCompilerOptions()
func RefactoringDriver_GetCompilerOptions() *CompilerOptions {
	clinit_CompilerOptions()
	var val_2 CompilerOptions
	InitCompilerOptions(&val_2)
	var options *CompilerOptions = &val_2
	_ = options
	options.SetLanguageIn(StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_NEXT())
	options.SetLanguageOut(StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT5())
	options.SetSummaryDetailLevel(int32(0))
	options.SetDependencyOptions(DependencyOptions_SortOnly())
	options.SetChecksOnly(true)
	options.SetContinueAfterErrors(true)
	options.SetParseJsDocDocumentation(pkg_parsing.StaticGet_Config_JsDocParsing_INCLUDE_ALL_COMMENTS())
	options.SetCheckSuspiciousCode(true)
	options.SetCheckSymbols(true)
	options.SetCheckTypes(true)
	options.SetClosurePass(true)
	options.SetGenerateExports(true)
	options.SetPreserveClosurePrimitives(true)
	options.SetWarningLevel(StaticGet_DiagnosticGroups_EXTRA_REQUIRE(), StaticGet_CheckLevel_WARNING())
	return options
}

// Java class: com.google.javascript.refactoring.RefactoringDriver.Builder
type RefactoringDriver_Builder struct {
	inputs          *j2gort.Ext_com_google_common_collect_ImmutableList_Builder `gson:"inputs"`
	externs         *j2gort.Ext_com_google_common_collect_ImmutableList_Builder `gson:"externs"`
	compilerOptions *CompilerOptions                                            `gson:"compilerOptions"`
	idHash_         atomic.Int32
}

func (this_ *RefactoringDriver_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *RefactoringDriver_Builder) IsJavaNil_() bool { return this_ == nil }
func (*RefactoringDriver_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.RefactoringDriver$Builder", "Builder"
}
func (this_ *RefactoringDriver_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RefactoringDriver_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.RefactoringDriver.Builder", "com.google.javascript.refactoring.RefactoringDriver$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*RefactoringDriver_Builder](v_)
	})
}

var static_RefactoringDriver_Builder_TO_SOURCE_FILE_FN j2gort.Ext_com_google_common_base_Function

func StaticGet_RefactoringDriver_Builder_TO_SOURCE_FILE_FN() j2gort.Ext_com_google_common_base_Function {
	clinit_RefactoringDriver_Builder()
	return static_RefactoringDriver_Builder_TO_SOURCE_FILE_FN
}
func StaticSet_RefactoringDriver_Builder_TO_SOURCE_FILE_FN(v j2gort.Ext_com_google_common_base_Function) {
	clinit_RefactoringDriver_Builder()
	static_RefactoringDriver_Builder_TO_SOURCE_FILE_FN = v
}

var clinitDone_RefactoringDriver_Builder uint32
var clinitState_RefactoringDriver_Builder j2gort.ClinitState

func clinit_RefactoringDriver_Builder() {
	if atomic.LoadUint32(&clinitDone_RefactoringDriver_Builder) != 0 {
		return
	}
	clinitSlow_RefactoringDriver_Builder()
}

func clinitSlow_RefactoringDriver_Builder() {
	if !j2gort.EnterClinit(&clinitState_RefactoringDriver_Builder, &clinitDone_RefactoringDriver_Builder) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RefactoringDriver_Builder, &clinitDone_RefactoringDriver_Builder)
	static_RefactoringDriver_Builder_TO_SOURCE_FILE_FN = lambdaSingleton_refactoring_driver_1
}

func (this_ *RefactoringDriver_Builder) initFields_RefactoringDriver_Builder() {
	_ = this_
	this_.inputs = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	this_.externs = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	this_.compilerOptions = RefactoringDriver_GetCompilerOptions()
}

// Java constructor: Builder()
func NewRefactoringDriver_Builder() *RefactoringDriver_Builder {
	this_ := &RefactoringDriver_Builder{}
	InitRefactoringDriver_Builder(this_)
	return this_
}

func InitRefactoringDriver_Builder(this_ *RefactoringDriver_Builder) {
	this_.initFields_RefactoringDriver_Builder()
}

// Java method: Builder.addExternsFromFile(java.lang.String filename)
func (this_ *RefactoringDriver_Builder) AddExternsFromFile_1_string(filename string) *RefactoringDriver_Builder {
	_ = filename
	this_.externs.AddVoid_(j2gort.PtrToAny(SourceFile_FromFile_1(filename)))
	return this_
}

// Java method: Builder.addExternsFromFile(java.lang.Iterable<java.lang.String> externs)
func (this_ *RefactoringDriver_Builder) AddExternsFromFile_1_Ext_java_lang_Iterable(externs j2gort.Ext_java_lang_Iterable) *RefactoringDriver_Builder {
	_ = externs
	this_.externs.AddAllVoid_(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_Lists_Transform(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable(externs)), StaticGet_RefactoringDriver_Builder_TO_SOURCE_FILE_FN())))
	return this_
}

// Java method: Builder.addExternsFromCode(java.lang.String code)
func (this_ *RefactoringDriver_Builder) AddExternsFromCode(code string) *RefactoringDriver_Builder {
	_ = code
	this_.externs.AddVoid_(j2gort.PtrToAny(SourceFile_FromCode_2("externs", code)))
	return this_
}

// Java method: Builder.addExterns(java.lang.Iterable<com.google.javascript.jscomp.SourceFile> externs)
func (this_ *RefactoringDriver_Builder) AddExterns(externs j2gort.Ext_java_lang_Iterable) *RefactoringDriver_Builder {
	_ = externs
	this_.externs.AddAllVoid_(externs)
	return this_
}

// Java method: Builder.addInputsFromFile(java.lang.String filename)
func (this_ *RefactoringDriver_Builder) AddInputsFromFile_1_string(filename string) *RefactoringDriver_Builder {
	_ = filename
	this_.inputs.AddVoid_(j2gort.PtrToAny(SourceFile_FromFile_1(filename)))
	return this_
}

// Java method: Builder.addInputsFromFile(java.lang.Iterable<java.lang.String> inputs)
func (this_ *RefactoringDriver_Builder) AddInputsFromFile_1_Ext_java_lang_Iterable(inputs j2gort.Ext_java_lang_Iterable) *RefactoringDriver_Builder {
	_ = inputs
	this_.inputs.AddAllVoid_(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_Lists_Transform(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable(inputs)), StaticGet_RefactoringDriver_Builder_TO_SOURCE_FILE_FN())))
	return this_
}

// Java method: Builder.addInputsFromCode(java.lang.String code)
func (this_ *RefactoringDriver_Builder) AddInputsFromCode_1(code string) *RefactoringDriver_Builder {
	_ = code
	return this_.AddInputsFromCode_2(code, "input")
}

// Java method: Builder.addInputsFromCode(java.lang.String code, java.lang.String filename)
func (this_ *RefactoringDriver_Builder) AddInputsFromCode_2(code string, filename string) *RefactoringDriver_Builder {
	_ = code
	_ = filename
	this_.inputs.AddVoid_(j2gort.PtrToAny(SourceFile_FromCode_2(filename, code)))
	return this_
}

// Java method: Builder.addInputs(java.lang.Iterable<com.google.javascript.jscomp.SourceFile> inputs)
func (this_ *RefactoringDriver_Builder) AddInputs(inputs j2gort.Ext_java_lang_Iterable) *RefactoringDriver_Builder {
	_ = inputs
	this_.inputs.AddAllVoid_(inputs)
	return this_
}

// Java method: Builder.withCompilerOptions(com.google.javascript.jscomp.CompilerOptions compilerOptions)
func (this_ *RefactoringDriver_Builder) WithCompilerOptions(compilerOptions *CompilerOptions) *RefactoringDriver_Builder {
	_ = compilerOptions
	this_.compilerOptions = j2gort.CheckNotNullPtr(compilerOptions)
	return this_
}

// Java method: Builder.build()
func (this_ *RefactoringDriver_Builder) Build() *RefactoringDriver {
	return NewRefactoringDriver(j2gort.AsIface[j2gort.Ext_java_util_List](this_.inputs.BuildImmList_()), j2gort.AsIface[j2gort.Ext_java_util_List](this_.externs.BuildImmList_()), this_.compilerOptions)
}
