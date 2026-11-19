/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/bundle/CompilerBasedTransformer.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.bundle.CompilerBasedTransformer
type CompilerBasedTransformer interface {
	Implements_jscomp_CompilerBasedTransformer()
	IsJavaNil_() bool
	GetCompilerBasedTransformerBase() *CompilerBasedTransformerBase
	Source_Transformer
	GetRuntime() *j2gort.Ext_java_util_Optional
	GetTranformationName() string
}

func ToIface_CompilerBasedTransformer[T interface {
	CompilerBasedTransformer
	comparable
}](v T) CompilerBasedTransformer {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_CompilerBasedTransformer[T CompilerBasedTransformer](v T) CompilerBasedTransformer {
	var vi CompilerBasedTransformer = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.bundle.CompilerBasedTransformer
type CompilerBasedTransformerBase struct {
	self_            CompilerBasedTransformer
	compilerSupplier CompilerBasedTransformer_CompilerSupplier `gson:"compilerSupplier"`
	idHash_          atomic.Int32
}

func (this_ *CompilerBasedTransformerBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(CompilerBasedTransformer)
}
func (this_ *CompilerBasedTransformerBase) IsJavaNil_() bool { return this_ == nil }
func (*CompilerBasedTransformerBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.bundle.CompilerBasedTransformer", "CompilerBasedTransformer"
}
func (this_ *CompilerBasedTransformerBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.bundle.CompilerBasedTransformer", "com.google.javascript.jscomp.bundle.CompilerBasedTransformer", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerBasedTransformer](v_)
	})
}
func (b *CompilerBasedTransformerBase) GetCompilerBasedTransformerBase() *CompilerBasedTransformerBase {
	return b
}
func (*CompilerBasedTransformerBase) Implements_jscomp_CompilerBasedTransformer() {}
func (*CompilerBasedTransformerBase) Implements_jscomp_Source_Transformer()       {}

func (this_ *CompilerBasedTransformerBase) initFields_CompilerBasedTransformer() {
	_ = this_
}

func InitCompilerBasedTransformer(this_ *CompilerBasedTransformerBase, compilerSupplier CompilerBasedTransformer_CompilerSupplier) {
	_ = compilerSupplier
	this_.initFields_CompilerBasedTransformer()
	this_.compilerSupplier = j2gort.CheckNotNullIface(compilerSupplier)
}

// Java method: CompilerBasedTransformer.getRuntime()
func (this_ *CompilerBasedTransformerBase) GetRuntime() *j2gort.Ext_java_util_Optional {
	panic("abstract or native method invoked")
}

// Java method: CompilerBasedTransformer.getTranformationName()
func (this_ *CompilerBasedTransformerBase) GetTranformationName() string {
	panic("abstract or native method invoked")
}

// Java method: CompilerBasedTransformer.transform(com.google.javascript.jscomp.bundle.Source input)
func (this_ *CompilerBasedTransformerBase) Transform(input Source) Source {
	_ = input
	var result *CompilerBasedTransformer_CompileResult = this_.compilerSupplier.GetCompilerBasedTransformer_CompilerSupplierBase().Compile(input.Path(), input.GetSourceBase().Code())
	_ = result
	if len(result.errors) != 0 {
		var formatter MessageFormatter = StaticGet_ErrorFormat_SOURCELESS().ToFormatter(nil, false)
		_ = formatter
		var message j2gort.LocalStringBuilder
		message.AppendStr(this_.self_.GetTranformationName()).AppendStr(" failed.\n")

		{

			for _, error_ := range result.errors {

				_ = error_
				message.AppendStr(formatter.FormatError(error_))
			}
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(message.String()))))
	}
	if !(result.transformed) {
		return input
	}
	var builder Source_Builder = input.ToBuilder().GetSource_BuilderBase().SetCode_1_string(result.source).SetSourceMap(result.sourceMap)
	_ = builder
	if this_.self_.GetRuntime().IsPresent() {
		builder.GetSource_BuilderBase().AddRuntime([]string{this_.compilerSupplier.GetCompilerBasedTransformer_CompilerSupplierBase().Runtime(j2gort.CastString(this_.self_.GetRuntime().Get()))})
	}
	return builder.GetSource_BuilderBase().Build()
}

func (this_ *CompilerBasedTransformerBase) AsFunction() j2gort.Ext_java_util_function_Function {
	return IfaceDefault_Source_Transformer_AsFunction(this_.self_)
}

func (this_ *CompilerBasedTransformerBase) AndThen(p0_ Source_Transformer) Source_Transformer {
	return IfaceDefault_Source_Transformer_AndThen(this_.self_, p0_)
}

func (this_ *CompilerBasedTransformerBase) Compose(p0_ Source_Transformer) Source_Transformer {
	return IfaceDefault_Source_Transformer_Compose(this_.self_, p0_)
}

// Java polymorphic class interface: com.google.javascript.jscomp.bundle.CompilerBasedTransformer.CompilerSupplier
type CompilerBasedTransformer_CompilerSupplier interface {
	Implements_jscomp_CompilerBasedTransformer_CompilerSupplier()
	IsJavaNil_() bool
	GetCompilerBasedTransformer_CompilerSupplierBase() *CompilerBasedTransformer_CompilerSupplierBase
	Compile(j2gort.Ext_java_nio_file_Path, string) *CompilerBasedTransformer_CompileResult
	Transformed(*Result) bool
	Runtime(string) string
	Compiler() *Compiler
	Options() *CompilerOptions
	SetOptions(*CompilerOptions)
}

func ToIface_CompilerBasedTransformer_CompilerSupplier[T interface {
	CompilerBasedTransformer_CompilerSupplier
	comparable
}](v T) CompilerBasedTransformer_CompilerSupplier {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_CompilerBasedTransformer_CompilerSupplier[T CompilerBasedTransformer_CompilerSupplier](v T) CompilerBasedTransformer_CompilerSupplier {
	var vi CompilerBasedTransformer_CompilerSupplier = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.bundle.CompilerBasedTransformer.CompilerSupplier
type CompilerBasedTransformer_CompilerSupplierBase struct {
	self_   CompilerBasedTransformer_CompilerSupplier
	idHash_ atomic.Int32
}

func (this_ *CompilerBasedTransformer_CompilerSupplierBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(CompilerBasedTransformer_CompilerSupplier)
}
func (this_ *CompilerBasedTransformer_CompilerSupplierBase) IsJavaNil_() bool { return this_ == nil }
func (*CompilerBasedTransformer_CompilerSupplierBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.bundle.CompilerBasedTransformer$CompilerSupplier", "CompilerSupplier"
}
func (this_ *CompilerBasedTransformer_CompilerSupplierBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.bundle.CompilerBasedTransformer.CompilerSupplier", "com.google.javascript.jscomp.bundle.CompilerBasedTransformer$CompilerSupplier", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerBasedTransformer_CompilerSupplier](v_)
	})
}
func (b *CompilerBasedTransformer_CompilerSupplierBase) GetCompilerBasedTransformer_CompilerSupplierBase() *CompilerBasedTransformer_CompilerSupplierBase {
	return b
}
func (*CompilerBasedTransformer_CompilerSupplierBase) Implements_jscomp_CompilerBasedTransformer_CompilerSupplier() {
}

var static_CompilerBasedTransformer_CompilerSupplier_EXTERNS *SourceFile

func StaticGet_CompilerBasedTransformer_CompilerSupplier_EXTERNS() *SourceFile {
	clinit_CompilerBasedTransformer_CompilerSupplier()
	return static_CompilerBasedTransformer_CompilerSupplier_EXTERNS
}
func StaticSet_CompilerBasedTransformer_CompilerSupplier_EXTERNS(v *SourceFile) {
	clinit_CompilerBasedTransformer_CompilerSupplier()
	static_CompilerBasedTransformer_CompilerSupplier_EXTERNS = v
}

var static_CompilerBasedTransformer_CompilerSupplier_EMPTY *SourceFile

func StaticGet_CompilerBasedTransformer_CompilerSupplier_EMPTY() *SourceFile {
	clinit_CompilerBasedTransformer_CompilerSupplier()
	return static_CompilerBasedTransformer_CompilerSupplier_EMPTY
}
func StaticSet_CompilerBasedTransformer_CompilerSupplier_EMPTY(v *SourceFile) {
	clinit_CompilerBasedTransformer_CompilerSupplier()
	static_CompilerBasedTransformer_CompilerSupplier_EMPTY = v
}

var static_CompilerBasedTransformer_CompilerSupplier_ES5_WARNINGS *DiagnosticGroup

func StaticGet_CompilerBasedTransformer_CompilerSupplier_ES5_WARNINGS() *DiagnosticGroup {
	clinit_CompilerBasedTransformer_CompilerSupplier()
	return static_CompilerBasedTransformer_CompilerSupplier_ES5_WARNINGS
}
func StaticSet_CompilerBasedTransformer_CompilerSupplier_ES5_WARNINGS(v *DiagnosticGroup) {
	clinit_CompilerBasedTransformer_CompilerSupplier()
	static_CompilerBasedTransformer_CompilerSupplier_ES5_WARNINGS = v
}

var clinitDone_CompilerBasedTransformer_CompilerSupplier uint32
var clinitState_CompilerBasedTransformer_CompilerSupplier j2gort.ClinitState

func clinit_CompilerBasedTransformer_CompilerSupplier() {
	if atomic.LoadUint32(&clinitDone_CompilerBasedTransformer_CompilerSupplier) != 0 {
		return
	}
	clinitSlow_CompilerBasedTransformer_CompilerSupplier()
}

func clinitSlow_CompilerBasedTransformer_CompilerSupplier() {
	if !j2gort.EnterClinit(&clinitState_CompilerBasedTransformer_CompilerSupplier, &clinitDone_CompilerBasedTransformer_CompilerSupplier) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CompilerBasedTransformer_CompilerSupplier, &clinitDone_CompilerBasedTransformer_CompilerSupplier)
	static_CompilerBasedTransformer_CompilerSupplier_EXTERNS = SourceFile_FromCode_2("externs.js", "function Symbol() {}")
	static_CompilerBasedTransformer_CompilerSupplier_EMPTY = SourceFile_FromCode_2("empty.js", "")
	static_CompilerBasedTransformer_CompilerSupplier_ES5_WARNINGS = NewDiagnosticGroup_1_DiagnosticType_2([]*DiagnosticType{DiagnosticType_Error("JSC_CANNOT_CONVERT", "")})
}

func (this_ *CompilerBasedTransformer_CompilerSupplierBase) initFields_CompilerBasedTransformer_CompilerSupplier() {
	_ = this_
}

// Java constructor: CompilerSupplier()
func NewCompilerBasedTransformer_CompilerSupplier() CompilerBasedTransformer_CompilerSupplier {
	clinit_CompilerBasedTransformer_CompilerSupplier()
	this_ := &CompilerBasedTransformer_CompilerSupplierBase{}
	this_.BindSelf_(this_)
	InitCompilerBasedTransformer_CompilerSupplier(this_)
	return this_
}

func InitCompilerBasedTransformer_CompilerSupplier(this_ *CompilerBasedTransformer_CompilerSupplierBase) {
	this_.initFields_CompilerBasedTransformer_CompilerSupplier()
}

// Java method: CompilerSupplier.compile(java.nio.file.Path path, java.lang.String code)
func (this_ *CompilerBasedTransformer_CompilerSupplierBase) Compile(path j2gort.Ext_java_nio_file_Path, code string) *CompilerBasedTransformer_CompileResult {
	_ = path
	_ = code
	var compiler *Compiler = this_.Compiler()
	_ = compiler
	var result *Result = compiler.Compile_3_SourceFile_SourceFile_CompilerOptions(StaticGet_CompilerBasedTransformer_CompilerSupplier_EXTERNS(), SourceFile_FromCode_2(path.String(), code), this_.Options())
	_ = result
	var source string = compiler.ToSource_0()
	_ = source
	var sourceMap *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sourceMap
	if result.sourceMap != nil {
		func() {
			defer func() {
				if rec_3 := recover(); rec_3 != nil {
					jex_4 := javalang.RecoverException(rec_3)
					_ = jex_4
					if j2gort.InstanceOf[javaio.IIOException](jex_4) {
						e := j2gort.AsIface[javaio.IIOException](jex_4)
						_ = e
					} else {
						panic(rec_3)
					}
				}
			}()
			result.sourceMap.AppendTo(sourceMap, path.String())
		}()
	}
	var transformed bool = this_.self_.Transformed(result)
	_ = transformed
	return NewCompilerBasedTransformer_CompileResult(source, result.errors, transformed, func() string {
		if transformed {
			return sourceMap.String()
		}
		return ""
	}())
}

// Java method: CompilerSupplier.transformed(com.google.javascript.jscomp.Result result)
func (this_ *CompilerBasedTransformer_CompilerSupplierBase) Transformed(result *Result) bool {
	_ = result
	return result.transpiledFiles
}

// Java method: CompilerSupplier.runtime(java.lang.String library)
func (this_ *CompilerBasedTransformer_CompilerSupplierBase) Runtime(library string) string {
	_ = library
	var compiler *Compiler = this_.Compiler()
	_ = compiler
	var options *CompilerOptions = this_.Options()
	_ = options
	options.SetForceLibraryInjection(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.ImmListOfStr1(library)))
	compiler.Compile_3_SourceFile_SourceFile_CompilerOptions(StaticGet_CompilerBasedTransformer_CompilerSupplier_EXTERNS(), StaticGet_CompilerBasedTransformer_CompilerSupplier_EMPTY(), options)
	return compiler.ToSource_0()
}

// Java method: CompilerSupplier.compiler()
func (this_ *CompilerBasedTransformer_CompilerSupplierBase) Compiler() *Compiler {
	return NewCompiler()
}

// Java method: CompilerSupplier.options()
func (this_ *CompilerBasedTransformer_CompilerSupplierBase) Options() *CompilerOptions {
	clinit_CompilerOptions()
	var val_5 CompilerOptions
	InitCompilerOptions(&val_5)
	var options *CompilerOptions = &val_5
	_ = options
	this_.self_.SetOptions(options)
	return options
}

// Java method: CompilerSupplier.setOptions(com.google.javascript.jscomp.CompilerOptions options)
func (this_ *CompilerBasedTransformer_CompilerSupplierBase) SetOptions(options *CompilerOptions) {
	_ = options
	options.SetLanguageOut(StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT5())
	options.SetQuoteKeywordProperties(true)
	options.SetSkipNonTranspilationPasses(true)
	options.SetVariableRenaming(StaticGet_VariableRenamingPolicy_OFF())
	options.SetPropertyRenaming(StaticGet_PropertyRenamingPolicy_OFF())
	options.SetWrapGoogModulesForWhitespaceOnly(false)
	options.SetPrettyPrint(true)
	options.SetSourceMapOutputPath("/dev/null")
	options.SetSourceMapIncludeSourcesContent(true)
	options.SetWarningLevel(StaticGet_CompilerBasedTransformer_CompilerSupplier_ES5_WARNINGS(), StaticGet_CheckLevel_OFF())
}

// Java class: com.google.javascript.jscomp.bundle.CompilerBasedTransformer.CompileResult
type CompilerBasedTransformer_CompileResult struct {
	source      string     `gson:"source"`
	errors      []*JSError `gson:"errors"`
	sourceMap   string     `gson:"sourceMap"`
	idHash_     atomic.Int32
	transformed bool `gson:"transformed"`
}

func (this_ *CompilerBasedTransformer_CompileResult) BindSelf_(s any) {
	_ = s
}
func (this_ *CompilerBasedTransformer_CompileResult) IsJavaNil_() bool { return this_ == nil }
func (*CompilerBasedTransformer_CompileResult) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.bundle.CompilerBasedTransformer$CompileResult", "CompileResult"
}
func (this_ *CompilerBasedTransformer_CompileResult) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CompilerBasedTransformer_CompileResult) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.bundle.CompilerBasedTransformer.CompileResult", "com.google.javascript.jscomp.bundle.CompilerBasedTransformer$CompileResult", func(v_ any) bool {
		return j2gort.InstanceOf[*CompilerBasedTransformer_CompileResult](v_)
	})
}

func (this_ *CompilerBasedTransformer_CompileResult) initFields_CompilerBasedTransformer_CompileResult() {
	_ = this_
	this_.source = j2gort.NullString
	this_.sourceMap = j2gort.NullString
}

// Java constructor: CompileResult(java.lang.String source, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.JSError> errors, boolean transformed, java.lang.String sourceMap)
func NewCompilerBasedTransformer_CompileResult(source string, errors j2gort.Ext_com_google_common_collect_ImmutableList, transformed bool, sourceMap string) *CompilerBasedTransformer_CompileResult {
	this_ := &CompilerBasedTransformer_CompileResult{}
	InitCompilerBasedTransformer_CompileResult(this_, source, errors, transformed, sourceMap)
	return this_
}

func InitCompilerBasedTransformer_CompileResult(this_ *CompilerBasedTransformer_CompileResult, source string, errors j2gort.Ext_com_google_common_collect_ImmutableList, transformed bool, sourceMap string) {
	_ = source
	_ = errors
	_ = transformed
	_ = sourceMap
	this_.initFields_CompilerBasedTransformer_CompileResult()
	this_.source = j2gort.CheckNotNullString(source)
	this_.errors = j2gort.ImmListToSliceNullable[*JSError](j2gort.CheckNotNull(errors))
	this_.transformed = transformed
	this_.sourceMap = j2gort.CheckNotNullString(sourceMap)
}
