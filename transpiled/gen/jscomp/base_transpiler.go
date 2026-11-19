/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/transpile/BaseTranspiler.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_base_transpiler_1 SourceMap_LocationMapping

func init() {
	lambdaSingleton_base_transpiler_1 = NewLambda_SourceMap_LocationMapping(func(location string) string {
		_ = location
		var didRet_4 bool
		var retVal_5 string
		func() {
			defer func() {
				if rec_6 := recover(); rec_6 != nil {
					jex_7 := javalang.RecoverException(rec_6)
					_ = jex_7
					if j2gort.InstanceOf[j2gort.Ext_java_net_URISyntaxException](jex_7) {
						e := j2gort.AsIface[j2gort.Ext_java_net_URISyntaxException](jex_7)
						_ = e
					} else {
						panic(rec_6)
					}
				}
			}()
			if j2gort.NewExt_java_net_URI_1(location).GetScheme() != j2gort.NullString {
				retVal_5 = location
				didRet_4 = true
				return
			}
		}()
		if didRet_4 {
			return retVal_5
		}
		return NewSourceMap_PrefixLocationMapping("", "/").Map(location)
	})
}

// Java class: com.google.javascript.jscomp.transpile.BaseTranspiler
type BaseTranspiler struct {
	compilerSupplier   *BaseTranspiler_CompilerSupplier `gson:"compilerSupplier"`
	runtimeLibraryName string                           `gson:"runtimeLibraryName"`
	idHash_            atomic.Int32
}

func (this_ *BaseTranspiler) BindSelf_(s any) {
	_ = s
}
func (this_ *BaseTranspiler) IsJavaNil_() bool { return this_ == nil }
func (*BaseTranspiler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.transpile.BaseTranspiler", "BaseTranspiler"
}
func (this_ *BaseTranspiler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*BaseTranspiler) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.transpile.BaseTranspiler", "com.google.javascript.jscomp.transpile.BaseTranspiler", func(v_ any) bool {
		return j2gort.InstanceOf[*BaseTranspiler](v_)
	})
}
func (*BaseTranspiler) Implements_jscomp_Transpile_Transpiler() {}

var static_BaseTranspiler_LATEST_TRANSPILER *BaseTranspiler

func StaticGet_BaseTranspiler_LATEST_TRANSPILER() *BaseTranspiler {
	clinit_BaseTranspiler()
	return static_BaseTranspiler_LATEST_TRANSPILER
}
func StaticSet_BaseTranspiler_LATEST_TRANSPILER(v *BaseTranspiler) {
	clinit_BaseTranspiler()
	static_BaseTranspiler_LATEST_TRANSPILER = v
}

var static_BaseTranspiler_ES5_TRANSPILER *BaseTranspiler

func StaticGet_BaseTranspiler_ES5_TRANSPILER() *BaseTranspiler {
	clinit_BaseTranspiler()
	return static_BaseTranspiler_ES5_TRANSPILER
}
func StaticSet_BaseTranspiler_ES5_TRANSPILER(v *BaseTranspiler) {
	clinit_BaseTranspiler()
	static_BaseTranspiler_ES5_TRANSPILER = v
}

var clinitDone_BaseTranspiler uint32
var clinitState_BaseTranspiler j2gort.ClinitState

func clinit_BaseTranspiler() {
	if atomic.LoadUint32(&clinitDone_BaseTranspiler) != 0 {
		return
	}
	clinitSlow_BaseTranspiler()
}

func clinitSlow_BaseTranspiler() {
	if !j2gort.EnterClinit(&clinitState_BaseTranspiler, &clinitDone_BaseTranspiler) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_BaseTranspiler, &clinitDone_BaseTranspiler)
	static_BaseTranspiler_LATEST_TRANSPILER = BaseTranspiler_To_2(pkg_parser.FeatureSet_Latest(), "")
	static_BaseTranspiler_ES5_TRANSPILER = BaseTranspiler_To_1(StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT5().ToFeatureSet())
}

func (this_ *BaseTranspiler) initFields_BaseTranspiler() {
	_ = this_
	this_.runtimeLibraryName = j2gort.NullString
}

// Java constructor: BaseTranspiler(com.google.javascript.jscomp.transpile.BaseTranspiler.CompilerSupplier compilerSupplier, java.lang.String runtimeLibraryName)
func NewBaseTranspiler(compilerSupplier *BaseTranspiler_CompilerSupplier, runtimeLibraryName string) *BaseTranspiler {
	clinit_BaseTranspiler()
	this_ := &BaseTranspiler{}
	InitBaseTranspiler(this_, compilerSupplier, runtimeLibraryName)
	return this_
}

func InitBaseTranspiler(this_ *BaseTranspiler, compilerSupplier *BaseTranspiler_CompilerSupplier, runtimeLibraryName string) {
	_ = compilerSupplier
	_ = runtimeLibraryName
	this_.initFields_BaseTranspiler()
	this_.compilerSupplier = j2gort.CheckNotNullPtr(compilerSupplier)
	this_.runtimeLibraryName = j2gort.CheckNotNullString(runtimeLibraryName)
}

// Java method: BaseTranspiler.transpile(java.net.URI path, java.lang.String code)
func (this_ *BaseTranspiler) Transpile(path *j2gort.Ext_java_net_URI, code string) *TranspileResult {
	_ = path
	_ = code
	var result *BaseTranspiler_CompileResult = this_.compilerSupplier.Compile(path, code)
	_ = result
	if !(result.transpiled) {
		return NewTranspileResult(path, code, code, "")
	}
	return NewTranspileResult(path, code, result.source, result.sourceMap)
}

// Java method: BaseTranspiler.runtime()
func (this_ *BaseTranspiler) Runtime() string {
	var sb j2gort.LocalStringBuilder

	if !(j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(this_.runtimeLibraryName)) {
		sb.AppendStr(this_.compilerSupplier.Runtime(this_.runtimeLibraryName))
	}
	sb.AppendStr(this_.compilerSupplier.Runtime("modules"))
	return sb.String()
}

// Java static method: BaseTranspiler.to(com.google.javascript.jscomp.parsing.parser.FeatureSet featureSet, java.lang.String runtime)
func BaseTranspiler_To_2(featureSet *pkg_parser.FeatureSet, runtime string) *BaseTranspiler {
	clinit_BaseTranspiler()
	_ = featureSet
	_ = runtime
	return NewBaseTranspiler(NewBaseTranspiler_CompilerSupplier_1(featureSet), runtime)
}

// Java static method: BaseTranspiler.to(com.google.javascript.jscomp.parsing.parser.FeatureSet featureSet)
func BaseTranspiler_To_1(featureSet *pkg_parser.FeatureSet) *BaseTranspiler {
	clinit_BaseTranspiler()
	_ = featureSet
	return BaseTranspiler_To_2(featureSet, "es6_runtime")
}

// Java class: com.google.javascript.jscomp.transpile.BaseTranspiler.CompilerSupplier
type BaseTranspiler_CompilerSupplier struct {
	moduleRoots        j2gort.Ext_com_google_common_collect_ImmutableList `gson:"moduleRoots"`
	prefixReplacements j2gort.Ext_com_google_common_collect_ImmutableMap  `gson:"prefixReplacements"`
	outputFeatureSet   *pkg_parser.FeatureSet                             `gson:"outputFeatureSet"`
	idHash_            atomic.Int32
	moduleResolution   ModuleLoader_ResolutionMode `gson:"moduleResolution"`
}

func (this_ *BaseTranspiler_CompilerSupplier) BindSelf_(s any) {
	_ = s
}
func (this_ *BaseTranspiler_CompilerSupplier) IsJavaNil_() bool { return this_ == nil }
func (*BaseTranspiler_CompilerSupplier) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.transpile.BaseTranspiler$CompilerSupplier", "CompilerSupplier"
}
func (this_ *BaseTranspiler_CompilerSupplier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*BaseTranspiler_CompilerSupplier) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.transpile.BaseTranspiler.CompilerSupplier", "com.google.javascript.jscomp.transpile.BaseTranspiler$CompilerSupplier", func(v_ any) bool {
		return j2gort.InstanceOf[*BaseTranspiler_CompilerSupplier](v_)
	})
}

func (this_ *BaseTranspiler_CompilerSupplier) initFields_BaseTranspiler_CompilerSupplier() {
	_ = this_
}

// Java constructor: CompilerSupplier()
func NewBaseTranspiler_CompilerSupplier() *BaseTranspiler_CompilerSupplier {
	this_ := &BaseTranspiler_CompilerSupplier{}
	InitBaseTranspiler_CompilerSupplier(this_)
	return this_
}

func InitBaseTranspiler_CompilerSupplier(this_ *BaseTranspiler_CompilerSupplier) {
	InitBaseTranspiler_CompilerSupplier_1(this_, StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT5().ToFeatureSet())
}

// Java constructor: CompilerSupplier(com.google.javascript.jscomp.parsing.parser.FeatureSet outputFeatureSet)
func NewBaseTranspiler_CompilerSupplier_1(outputFeatureSet *pkg_parser.FeatureSet) *BaseTranspiler_CompilerSupplier {
	this_ := &BaseTranspiler_CompilerSupplier{}
	InitBaseTranspiler_CompilerSupplier_1(this_, outputFeatureSet)
	return this_
}

func InitBaseTranspiler_CompilerSupplier_1(this_ *BaseTranspiler_CompilerSupplier, outputFeatureSet *pkg_parser.FeatureSet) {
	_ = outputFeatureSet
	InitBaseTranspiler_CompilerSupplier_4(this_, outputFeatureSet, NewCompilerOptions().GetModuleResolutionMode(), j2gort.Ext_com_google_common_collect_ImmutableList_Of_0(), j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0())
}

// Java constructor: CompilerSupplier(com.google.javascript.jscomp.parsing.parser.FeatureSet outputFeatureSet, com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode moduleResolution, com.google.common.collect.ImmutableList<java.lang.String> moduleRoots, com.google.common.collect.ImmutableMap<java.lang.String,java.lang.String> prefixReplacements)
func NewBaseTranspiler_CompilerSupplier_4(outputFeatureSet *pkg_parser.FeatureSet, moduleResolution ModuleLoader_ResolutionMode, moduleRoots j2gort.Ext_com_google_common_collect_ImmutableList, prefixReplacements j2gort.Ext_com_google_common_collect_ImmutableMap) *BaseTranspiler_CompilerSupplier {
	this_ := &BaseTranspiler_CompilerSupplier{}
	InitBaseTranspiler_CompilerSupplier_4(this_, outputFeatureSet, moduleResolution, moduleRoots, prefixReplacements)
	return this_
}

func InitBaseTranspiler_CompilerSupplier_4(this_ *BaseTranspiler_CompilerSupplier, outputFeatureSet *pkg_parser.FeatureSet, moduleResolution ModuleLoader_ResolutionMode, moduleRoots j2gort.Ext_com_google_common_collect_ImmutableList, prefixReplacements j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = outputFeatureSet
	_ = moduleResolution
	_ = moduleRoots
	_ = prefixReplacements
	this_.initFields_BaseTranspiler_CompilerSupplier()
	this_.outputFeatureSet = outputFeatureSet
	this_.moduleResolution = moduleResolution
	this_.moduleRoots = moduleRoots
	this_.prefixReplacements = prefixReplacements
}

// Java method: CompilerSupplier.compile(java.net.URI path, java.lang.String code)
func (this_ *BaseTranspiler_CompilerSupplier) Compile(path *j2gort.Ext_java_net_URI, code string) *BaseTranspiler_CompileResult {
	_ = path
	_ = code
	var compiler *Compiler = this_.Compiler()
	_ = compiler
	var result *Result = compiler.Compile_3_SourceFile_SourceFile_CompilerOptions(this_.CreateTrivialExterns(), SourceFile_FromCode_2(path.String(), code), this_.Options())
	_ = result
	var source string = compiler.ToSource_0()
	_ = source
	var sourceMap *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sourceMap
	if result.sourceMap != nil {
		func() {
			defer func() {
				if rec_1 := recover(); rec_1 != nil {
					jex_2 := javalang.RecoverException(rec_1)
					_ = jex_2
					if j2gort.InstanceOf[javaio.IIOException](jex_2) {
						e := j2gort.AsIface[javaio.IIOException](jex_2)
						_ = e
					} else {
						panic(rec_1)
					}
				}
			}()
			result.sourceMap.AppendTo(sourceMap, path.String())
		}()
	}
	var transpiled bool = result.transpiledFiles
	_ = transpiled
	if !(result.errors.IsEmpty()) {
		panic(NewTranspilationException_3(ToIface_SourceExcerptProvider(compiler), result.errors, result.warnings))
	}
	return NewBaseTranspiler_CompileResult(source, transpiled, func() string {
		if transpiled {
			return sourceMap.String()
		}
		return ""
	}())
}

// Java method: CompilerSupplier.runtime(java.lang.String library)
func (this_ *BaseTranspiler_CompilerSupplier) Runtime(library string) string {
	_ = library
	var compiler *Compiler = this_.Compiler()
	_ = compiler
	var options *CompilerOptions = this_.Options()
	_ = options
	options.SetForceLibraryInjection(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.ImmListOfStr1(library)))
	compiler.Compile_3_SourceFile_SourceFile_CompilerOptions(this_.CreateTrivialExterns(), this_.CreateEmptySource(), options)
	return compiler.ToSource_0()
}

// Java method: CompilerSupplier.compiler()
func (this_ *BaseTranspiler_CompilerSupplier) Compiler() *Compiler {
	return NewCompiler()
}

// Java method: CompilerSupplier.options()
func (this_ *BaseTranspiler_CompilerSupplier) Options() *CompilerOptions {
	clinit_CompilerOptions()
	var val_3 CompilerOptions
	InitCompilerOptions(&val_3)
	var options *CompilerOptions = &val_3
	_ = options
	this_.SetOptions(options)
	return options
}

// Java method: CompilerSupplier.setOptions(com.google.javascript.jscomp.CompilerOptions options)
func (this_ *BaseTranspiler_CompilerSupplier) SetOptions(options *CompilerOptions) {
	_ = options
	options.SetLanguageIn(StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_NEXT())
	options.LegacySetOutputFeatureSet(this_.outputFeatureSet.Without_2(pkg_parser.StaticGet_FeatureSet_Feature_MODULES(), []*pkg_parser.FeatureSet_Feature{}))
	options.SetEmitUseStrict(false)
	options.SetQuoteKeywordProperties(true)
	options.SetSkipNonTranspilationPasses(true)
	options.SetVariableRenaming(StaticGet_VariableRenamingPolicy_OFF())
	options.SetPropertyRenaming(StaticGet_PropertyRenamingPolicy_OFF())
	options.SetWrapGoogModulesForWhitespaceOnly(false)
	options.SetPrettyPrint(true)
	options.SetWarningLevel(StaticGet_DiagnosticGroups_NON_STANDARD_JSDOC(), StaticGet_CheckLevel_OFF())
	options.SetEs6ModuleTranspilation(StaticGet_CompilerOptions_Es6ModuleTranspilation_TO_COMMON_JS_LIKE_MODULES())
	options.SetModuleResolutionMode(this_.moduleResolution)
	options.SetModuleRoots(j2gort.AsIface[j2gort.Ext_java_util_List](this_.moduleRoots))
	options.SetBrowserResolverPrefixReplacements(this_.prefixReplacements)
	options.SetWarningLevel(StaticGet_DiagnosticGroups_MODULE_LOAD(), StaticGet_CheckLevel_OFF())
	options.SetPathEscaper(StaticGet_ModuleLoader_PathEscaper_CANONICALIZE_ONLY())
	options.SetParseInlineSourceMaps(true)
	options.SetApplyInputSourceMaps(true)
	options.SetSourceMapOutputPath("/dev/null")
	options.SetSourceMapIncludeSourcesContent(true)
	options.SetSourceMapLocationMappings(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(lambdaSingleton_base_transpiler_1)))
}

// Java method: CompilerSupplier.createTrivialExterns()
func (this_ *BaseTranspiler_CompilerSupplier) CreateTrivialExterns() *SourceFile {
	return SourceFile_FromCode_2("externs.js", "function Symbol() {}")
}

// Java method: CompilerSupplier.createEmptySource()
func (this_ *BaseTranspiler_CompilerSupplier) CreateEmptySource() *SourceFile {
	return SourceFile_FromCode_2("empty.js", "")
}

// Java class: com.google.javascript.jscomp.transpile.BaseTranspiler.CompileResult
type BaseTranspiler_CompileResult struct {
	source     string `gson:"source"`
	sourceMap  string `gson:"sourceMap"`
	idHash_    atomic.Int32
	transpiled bool `gson:"transpiled"`
}

func (this_ *BaseTranspiler_CompileResult) BindSelf_(s any) {
	_ = s
}
func (this_ *BaseTranspiler_CompileResult) IsJavaNil_() bool { return this_ == nil }
func (*BaseTranspiler_CompileResult) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.transpile.BaseTranspiler$CompileResult", "CompileResult"
}
func (this_ *BaseTranspiler_CompileResult) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*BaseTranspiler_CompileResult) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.transpile.BaseTranspiler.CompileResult", "com.google.javascript.jscomp.transpile.BaseTranspiler$CompileResult", func(v_ any) bool {
		return j2gort.InstanceOf[*BaseTranspiler_CompileResult](v_)
	})
}

func (this_ *BaseTranspiler_CompileResult) initFields_BaseTranspiler_CompileResult() {
	_ = this_
	this_.source = j2gort.NullString
	this_.sourceMap = j2gort.NullString
}

// Java constructor: CompileResult(java.lang.String source, boolean transpiled, java.lang.String sourceMap)
func NewBaseTranspiler_CompileResult(source string, transpiled bool, sourceMap string) *BaseTranspiler_CompileResult {
	this_ := &BaseTranspiler_CompileResult{}
	InitBaseTranspiler_CompileResult(this_, source, transpiled, sourceMap)
	return this_
}

func InitBaseTranspiler_CompileResult(this_ *BaseTranspiler_CompileResult, source string, transpiled bool, sourceMap string) {
	_ = source
	_ = transpiled
	_ = sourceMap
	this_.initFields_BaseTranspiler_CompileResult()
	this_.source = j2gort.CheckNotNullString(source)
	this_.transpiled = transpiled
	this_.sourceMap = j2gort.CheckNotNullString(sourceMap)
}
