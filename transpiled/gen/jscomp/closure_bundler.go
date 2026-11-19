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
// Port of com/google/javascript/jscomp/deps/ClosureBundler.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.deps.ClosureBundler
type ClosureBundler struct {
	transpiler          Transpile_Transpiler     `gson:"transpiler"`
	es6ModuleTranspiler Transpile_Transpiler     `gson:"es6ModuleTranspiler"`
	mode                *ClosureBundler_EvalMode `gson:"mode"`
	sourceUrl           string                   `gson:"sourceUrl"`
	path                string                   `gson:"path"`
	embedSourcemap      bool                     `gson:"embedSourcemap"`
	sourceMapCache      j2gort.Ext_java_util_Map `gson:"sourceMapCache"`
	minifier            any                      `gson:"minifier"`
	idHash_             atomic.Int32
}

func (this_ *ClosureBundler) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureBundler) IsJavaNil_() bool { return this_ == nil }
func (*ClosureBundler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.ClosureBundler", "ClosureBundler"
}
func (this_ *ClosureBundler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureBundler) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.ClosureBundler", "com.google.javascript.jscomp.deps.ClosureBundler", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureBundler](v_)
	})
}

func (this_ *ClosureBundler) initFields_ClosureBundler() {
	_ = this_
	this_.sourceUrl = j2gort.NullString
	this_.path = j2gort.NullString
}

// Java constructor: ClosureBundler()
func NewClosureBundler() *ClosureBundler {
	this_ := &ClosureBundler{}
	InitClosureBundler(this_)
	return this_
}

func InitClosureBundler(this_ *ClosureBundler) {
	InitClosureBundler_1(this_, StaticGet_Transpile_Transpiler_NULL())
}

// Java constructor: ClosureBundler(com.google.javascript.jscomp.transpile.Transpiler transpiler)
func NewClosureBundler_1(transpiler Transpile_Transpiler) *ClosureBundler {
	this_ := &ClosureBundler{}
	InitClosureBundler_1(this_, transpiler)
	return this_
}

func InitClosureBundler_1(this_ *ClosureBundler, transpiler Transpile_Transpiler) {
	_ = transpiler
	InitClosureBundler_2(this_, transpiler, ToIface_Transpile_Transpiler(StaticGet_BaseTranspiler_LATEST_TRANSPILER()))
}

// Java constructor: ClosureBundler(com.google.javascript.jscomp.transpile.Transpiler transpiler, com.google.javascript.jscomp.transpile.Transpiler es6ModuleTranspiler)
func NewClosureBundler_2(transpiler Transpile_Transpiler, es6ModuleTranspiler Transpile_Transpiler) *ClosureBundler {
	this_ := &ClosureBundler{}
	InitClosureBundler_2(this_, transpiler, es6ModuleTranspiler)
	return this_
}

func InitClosureBundler_2(this_ *ClosureBundler, transpiler Transpile_Transpiler, es6ModuleTranspiler Transpile_Transpiler) {
	_ = transpiler
	_ = es6ModuleTranspiler
	InitClosureBundler_8(this_, transpiler, es6ModuleTranspiler, StaticGet_ClosureBundler_EvalMode_NORMAL(), j2gort.NullString, "unknown_source", nil, j2gort.NewExt_java_util_concurrent_ConcurrentHashMap(), false)
}

// Java constructor: ClosureBundler(com.google.javascript.jscomp.transpile.Transpiler transpiler, com.google.javascript.jscomp.transpile.Transpiler es6ModuleTranspiler, com.google.javascript.jscomp.deps.ClosureBundler.EvalMode mode, java.lang.@org.jspecify.annotations.Nullable String sourceUrl, java.lang.String path, java.lang.Object minifier, java.util.Map<java.lang.String,java.lang.String> sourceMapCache, boolean embedSourcemap)
func NewClosureBundler_8(transpiler Transpile_Transpiler, es6ModuleTranspiler Transpile_Transpiler, mode *ClosureBundler_EvalMode, sourceUrl string, path string, minifier any, sourceMapCache j2gort.Ext_java_util_Map, embedSourcemap bool) *ClosureBundler {
	this_ := &ClosureBundler{}
	InitClosureBundler_8(this_, transpiler, es6ModuleTranspiler, mode, sourceUrl, path, minifier, sourceMapCache, embedSourcemap)
	return this_
}

func InitClosureBundler_8(this_ *ClosureBundler, transpiler Transpile_Transpiler, es6ModuleTranspiler Transpile_Transpiler, mode *ClosureBundler_EvalMode, sourceUrl string, path string, minifier any, sourceMapCache j2gort.Ext_java_util_Map, embedSourcemap bool) {
	_ = transpiler
	_ = es6ModuleTranspiler
	_ = mode
	_ = sourceUrl
	_ = path
	_ = minifier
	_ = sourceMapCache
	_ = embedSourcemap
	this_.initFields_ClosureBundler()
	this_.transpiler = transpiler
	this_.mode = mode
	this_.sourceUrl = sourceUrl
	this_.path = path
	this_.sourceMapCache = sourceMapCache
	this_.es6ModuleTranspiler = es6ModuleTranspiler
	this_.minifier = minifier
	this_.embedSourcemap = embedSourcemap
}

// Java method: ClosureBundler.withTranspilers(com.google.javascript.jscomp.transpile.Transpiler newTranspiler, com.google.javascript.jscomp.transpile.Transpiler newEs6ModuleTranspiler)
func (this_ *ClosureBundler) WithTranspilers(newTranspiler Transpile_Transpiler, newEs6ModuleTranspiler Transpile_Transpiler) *ClosureBundler {
	_ = newTranspiler
	_ = newEs6ModuleTranspiler
	return NewClosureBundler_8(newTranspiler, newEs6ModuleTranspiler, this_.mode, this_.sourceUrl, this_.path, this_.minifier, this_.sourceMapCache, this_.embedSourcemap)
}

// Java method: ClosureBundler.withTranspiler(com.google.javascript.jscomp.transpile.Transpiler newTranspiler)
func (this_ *ClosureBundler) WithTranspiler(newTranspiler Transpile_Transpiler) *ClosureBundler {
	_ = newTranspiler
	return this_.WithTranspilers(newTranspiler, this_.es6ModuleTranspiler)
}

// Java method: ClosureBundler.withEs6ModuleTranspiler(com.google.javascript.jscomp.transpile.Transpiler newEs6ModuleTranspiler)
func (this_ *ClosureBundler) WithEs6ModuleTranspiler(newEs6ModuleTranspiler Transpile_Transpiler) *ClosureBundler {
	_ = newEs6ModuleTranspiler
	return this_.WithTranspilers(this_.transpiler, newEs6ModuleTranspiler)
}

// Java method: ClosureBundler.disableJ2clMinifier()
func (this_ *ClosureBundler) DisableJ2clMinifier() *ClosureBundler {
	return NewClosureBundler_8(this_.transpiler, this_.es6ModuleTranspiler, this_.mode, this_.sourceUrl, this_.path, nil, this_.sourceMapCache, this_.embedSourcemap)
}

// Java method: ClosureBundler.useEval(boolean useEval)
func (this_ *ClosureBundler) UseEval(useEval bool) *ClosureBundler {
	_ = useEval
	var newMode *ClosureBundler_EvalMode
	_ = newMode
	if useEval {
		newMode = StaticGet_ClosureBundler_EvalMode_EVAL()
	} else {
		newMode = StaticGet_ClosureBundler_EvalMode_NORMAL()
	}
	return NewClosureBundler_8(this_.transpiler, this_.es6ModuleTranspiler, newMode, this_.sourceUrl, this_.path, this_.minifier, this_.sourceMapCache, this_.embedSourcemap)
}

// Java method: ClosureBundler.withSourceUrl(java.lang.String newSourceUrl)
func (this_ *ClosureBundler) WithSourceUrl(newSourceUrl string) *ClosureBundler {
	_ = newSourceUrl
	return NewClosureBundler_8(this_.transpiler, this_.es6ModuleTranspiler, this_.mode, newSourceUrl, this_.path, this_.minifier, this_.sourceMapCache, this_.embedSourcemap)
}

// Java method: ClosureBundler.withPath(java.lang.String newPath)
func (this_ *ClosureBundler) WithPath(newPath string) *ClosureBundler {
	_ = newPath
	return NewClosureBundler_8(this_.transpiler, this_.es6ModuleTranspiler, this_.mode, this_.sourceUrl, newPath, this_.minifier, this_.sourceMapCache, this_.embedSourcemap)
}

// Java method: ClosureBundler.embedSourcemap()
func (this_ *ClosureBundler) EmbedSourcemap() *ClosureBundler {
	return NewClosureBundler_8(this_.transpiler, this_.es6ModuleTranspiler, this_.mode, this_.sourceUrl, this_.path, this_.minifier, this_.sourceMapCache, true)
}

// Java static method: ClosureBundler.appendInput(java.lang.Appendable out, com.google.javascript.jscomp.deps.DependencyInfo info, java.lang.String contents)
func ClosureBundler_AppendInput(out j2gort.Ext_java_lang_Appendable, info DependencyInfo, contents string) {
	_ = out
	_ = info
	_ = contents
	{
		var val_1 ClosureBundler
		InitClosureBundler(&val_1)
		(&val_1).AppendTo_3_Ext_java_lang_Appendable_DependencyInfo_string(out, info, contents)
	}
}

// Java method: ClosureBundler.appendTo(java.lang.Appendable out, com.google.javascript.jscomp.deps.DependencyInfo info, java.lang.String content)
func (this_ *ClosureBundler) AppendTo_3_Ext_java_lang_Appendable_DependencyInfo_string(out j2gort.Ext_java_lang_Appendable, info DependencyInfo, content string) {
	_ = out
	_ = info
	_ = content
	this_.AppendTo_3_Ext_java_lang_Appendable_DependencyInfo_Ext_com_google_common_io_CharSource(out, info, j2gort.Ext_com_google_common_io_CharSource_Wrap(j2gort.BoxString(content)))
}

// Java method: ClosureBundler.appendTo(java.lang.Appendable out, com.google.javascript.jscomp.deps.DependencyInfo info, java.io.File content, java.nio.charset.Charset contentCharset)
func (this_ *ClosureBundler) AppendTo_4(out j2gort.Ext_java_lang_Appendable, info DependencyInfo, content *j2gort.Ext_java_io_File, contentCharset j2gort.Ext_java_nio_charset_Charset) {
	_ = out
	_ = info
	_ = content
	_ = contentCharset
	this_.AppendTo_3_Ext_java_lang_Appendable_DependencyInfo_Ext_com_google_common_io_CharSource(out, info, j2gort.Ext_com_google_common_io_Files_AsCharSource(content, contentCharset))
}

// Java method: ClosureBundler.appendTo(java.lang.Appendable out, com.google.javascript.jscomp.deps.DependencyInfo info, com.google.common.io.CharSource content)
func (this_ *ClosureBundler) AppendTo_3_Ext_java_lang_Appendable_DependencyInfo_Ext_com_google_common_io_CharSource(out j2gort.Ext_java_lang_Appendable, info DependencyInfo, content j2gort.Ext_com_google_common_io_CharSource) {
	_ = out
	_ = info
	_ = content
	var code string = content.Read()
	_ = code
	if info.IsGoogModule() {
		this_.mode.AppendGoogModule(this_.priv_ClosureBundler_transpile_6031(code), out, this_.sourceUrl)
	} else if info.IsEs6Module() && j2gort.IfaceEquals[Transpile_Transpiler](this_.transpiler, StaticGet_Transpile_Transpiler_NULL()) {
		this_.mode.AppendTraditional(this_.priv_ClosureBundler_transpileEs6Module_6032(code), out, this_.sourceUrl)
	} else {
		this_.mode.AppendTraditional(this_.priv_ClosureBundler_transpile_6031(code), out, this_.sourceUrl)
	}
}

// Java method: ClosureBundler.appendRuntimeTo(java.lang.Appendable out)
func (this_ *ClosureBundler) AppendRuntimeTo(out j2gort.Ext_java_lang_Appendable) {
	_ = out
	var runtime string = this_.transpiler.Runtime()
	_ = runtime
	if !(j2gort.CheckNotNullString(runtime) == "") {
		this_.mode.AppendTraditional(runtime, out, j2gort.NullString)
	}
	if j2gort.IfaceEquals[Transpile_Transpiler](this_.transpiler, StaticGet_Transpile_Transpiler_NULL()) {
		this_.mode.AppendTraditional(this_.es6ModuleTranspiler.Runtime(), out, j2gort.NullString)
	}
	this_.mode.AppendTraditional("this.CLOSURE_EVAL_PREFILTER = function(s) { return s; };", out, j2gort.NullString)
	this_.mode.AppendTraditional("(function(thisValue){", out, j2gort.NullString)
	this_.mode.AppendTraditional("var isChrome87 = false; try {isChrome87 =  eval(trustedTypes.emptyScript) !== trustedTypes.emptyScript } catch (e) {} if (typeof trustedTypes !== 'undefined' && trustedTypes.createPolicy &&isChrome87 ) {", out, j2gort.NullString)
	this_.mode.AppendTraditional("  var policy = trustedTypes.createPolicy('goog#devserver',{ createScript: function(s){ return s; }});", out, j2gort.NullString)
	this_.mode.AppendTraditional("  thisValue.CLOSURE_EVAL_PREFILTER = policy.createScript.bind(policy);", out, j2gort.NullString)
	this_.mode.AppendTraditional("}", out, j2gort.NullString)
	this_.mode.AppendTraditional("})(this);", out, j2gort.NullString)
}

// Java method: ClosureBundler.getSourceMap(java.lang.String path)
func (this_ *ClosureBundler) GetSourceMap(path string) string {
	_ = path
	return j2gort.Ext_com_google_common_base_Strings_NullToEmpty(j2gort.CastString(j2gort.MapGetStr(this_.sourceMapCache, path)))
}

// Java method: ClosureBundler.transpile(java.lang.String s, com.google.javascript.jscomp.transpile.Transpiler t)
func (this_ *ClosureBundler) priv_ClosureBundler_transpile_6030(s string, t Transpile_Transpiler) string {
	_ = s
	_ = t
	var result *TranspileResult = nil
	_ = result
	func() {
		defer func() {
			if rec_2 := recover(); rec_2 != nil {
				jex_3 := javalang.RecoverException(rec_2)
				_ = jex_3
				if j2gort.InstanceOf[j2gort.Ext_java_net_URISyntaxException](jex_3) {
					e := j2gort.AsIface[j2gort.Ext_java_net_URISyntaxException](jex_3)
					_ = e
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
				} else {
					panic(rec_2)
				}
			}
		}()
		result = t.Transpile(j2gort.NewExt_java_net_URI_1(this_.path), s)
	}()
	j2gort.MapPutStr(this_.sourceMapCache, this_.path, j2gort.BoxString(result.SourceMap()))
	if this_.embedSourcemap {
		return result.EmbedSourcemapBase64().Transpiled()
	} else {
		return result.Transpiled()
	}
}

// Java method: ClosureBundler.transpile(java.lang.String s)
func (this_ *ClosureBundler) priv_ClosureBundler_transpile_6031(s string) string {
	_ = s
	return this_.priv_ClosureBundler_transpile_6030(s, this_.transpiler)
}

// Java method: ClosureBundler.transpileEs6Module(java.lang.String s)
func (this_ *ClosureBundler) priv_ClosureBundler_transpileEs6Module_6032(s string) string {
	_ = s
	return this_.priv_ClosureBundler_transpile_6030(s, this_.es6ModuleTranspiler)
}

// Java static method: ClosureBundler.appendSourceUrl(java.lang.Appendable out, com.google.javascript.jscomp.deps.ClosureBundler.EscapeMode mode, java.lang.String sourceUrl)
func ClosureBundler_AppendSourceUrl(out j2gort.Ext_java_lang_Appendable, mode *ClosureBundler_EscapeMode, sourceUrl string) {
	_ = out
	_ = mode
	_ = sourceUrl
	if sourceUrl == j2gort.NullString {
		return
	}
	var toAppend string = ("\n//# sourceURL=" + j2gort.StrValString(sourceUrl) + "\n")
	_ = toAppend
	mode.Append(toAppend, out)
}

// Java enum: com.google.javascript.jscomp.deps.ClosureBundler.EvalMode
type ClosureBundler_EvalMode struct {
	j2gort.EnumBase
	fn_AppendTraditional func(*ClosureBundler_EvalMode, string, j2gort.Ext_java_lang_Appendable, string)
	fn_AppendGoogModule  func(*ClosureBundler_EvalMode, string, j2gort.Ext_java_lang_Appendable, string)
}

func (this_ *ClosureBundler_EvalMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *ClosureBundler_EvalMode) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*ClosureBundler_EvalMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*ClosureBundler_EvalMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *ClosureBundler_EvalMode) initFields_ClosureBundler_EvalMode() {
	_ = this_
}

func initEnumClosureBundler_EvalMode(this_ *ClosureBundler_EvalMode, enumName_ string, enumOrd_ int32) *ClosureBundler_EvalMode {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	this_.initFields_ClosureBundler_EvalMode()
	return this_
}

var enum_ClosureBundler_EvalMode_EVAL *ClosureBundler_EvalMode

func StaticGet_ClosureBundler_EvalMode_EVAL() *ClosureBundler_EvalMode {
	clinit_ClosureBundler_EvalMode()
	return enum_ClosureBundler_EvalMode_EVAL
}

var enum_ClosureBundler_EvalMode_NORMAL *ClosureBundler_EvalMode

func StaticGet_ClosureBundler_EvalMode_NORMAL() *ClosureBundler_EvalMode {
	clinit_ClosureBundler_EvalMode()
	return enum_ClosureBundler_EvalMode_NORMAL
}

var clinitDone_ClosureBundler_EvalMode uint32
var clinitState_ClosureBundler_EvalMode j2gort.ClinitState

func clinit_ClosureBundler_EvalMode() {
	if atomic.LoadUint32(&clinitDone_ClosureBundler_EvalMode) != 0 {
		return
	}
	clinitSlow_ClosureBundler_EvalMode()
}

func clinitSlow_ClosureBundler_EvalMode() {
	if !j2gort.EnterClinit(&clinitState_ClosureBundler_EvalMode, &clinitDone_ClosureBundler_EvalMode) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ClosureBundler_EvalMode, &clinitDone_ClosureBundler_EvalMode)
	enum_ClosureBundler_EvalMode_EVAL = &ClosureBundler_EvalMode{}
	enum_ClosureBundler_EvalMode_EVAL.fn_AppendTraditional = func(this_ *ClosureBundler_EvalMode, s string, out j2gort.Ext_java_lang_Appendable, sourceUrl string) {
		_ = s
		_ = out
		_ = sourceUrl
		out.Append_1_any("eval(this.CLOSURE_EVAL_PREFILTER(\"")
		StaticGet_ClosureBundler_EscapeMode_ESCAPED().Append(s, out)
		ClosureBundler_AppendSourceUrl(out, StaticGet_ClosureBundler_EscapeMode_ESCAPED(), sourceUrl)
		out.Append_1_any("\"));\n")
	}
	enum_ClosureBundler_EvalMode_EVAL.fn_AppendGoogModule = func(this_ *ClosureBundler_EvalMode, s string, out j2gort.Ext_java_lang_Appendable, sourceUrl string) {
		_ = s
		_ = out
		_ = sourceUrl
		out.Append_1_any("goog.loadModule(\"")
		StaticGet_ClosureBundler_EscapeMode_ESCAPED().Append(s, out)
		ClosureBundler_AppendSourceUrl(out, StaticGet_ClosureBundler_EscapeMode_ESCAPED(), sourceUrl)
		out.Append_1_any("\");\n")
	}
	initEnumClosureBundler_EvalMode(enum_ClosureBundler_EvalMode_EVAL, "EVAL", int32(0))
	enum_ClosureBundler_EvalMode_NORMAL = &ClosureBundler_EvalMode{}
	enum_ClosureBundler_EvalMode_NORMAL.fn_AppendTraditional = func(this_ *ClosureBundler_EvalMode, s string, out j2gort.Ext_java_lang_Appendable, sourceUrl string) {
		_ = s
		_ = out
		_ = sourceUrl
		StaticGet_ClosureBundler_EscapeMode_NORMAL().Append(s, out)
		ClosureBundler_AppendSourceUrl(out, StaticGet_ClosureBundler_EscapeMode_NORMAL(), sourceUrl)
	}
	enum_ClosureBundler_EvalMode_NORMAL.fn_AppendGoogModule = func(this_ *ClosureBundler_EvalMode, s string, out j2gort.Ext_java_lang_Appendable, sourceUrl string) {
		_ = s
		_ = out
		_ = sourceUrl
		out.Append_1_any("goog.loadModule(function(exports) {'use strict';")
		StaticGet_ClosureBundler_EscapeMode_NORMAL().Append(s, out)
		out.Append_1_any("\n;return exports;});\n")
		ClosureBundler_AppendSourceUrl(out, StaticGet_ClosureBundler_EscapeMode_NORMAL(), sourceUrl)
	}
	initEnumClosureBundler_EvalMode(enum_ClosureBundler_EvalMode_NORMAL, "NORMAL", int32(1))
	enumValues_ClosureBundler_EvalMode = []*ClosureBundler_EvalMode{enum_ClosureBundler_EvalMode_EVAL, enum_ClosureBundler_EvalMode_NORMAL}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.ClosureBundler.EvalMode", []any{enum_ClosureBundler_EvalMode_EVAL, enum_ClosureBundler_EvalMode_NORMAL})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.ClosureBundler$EvalMode", []any{enum_ClosureBundler_EvalMode_EVAL, enum_ClosureBundler_EvalMode_NORMAL})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.deps.ClosureBundler.EvalMode", "com.google.javascript.jscomp.deps.ClosureBundler$EvalMode", func() {
		clinit_ClosureBundler_EvalMode()
	})
}

var enumValues_ClosureBundler_EvalMode []*ClosureBundler_EvalMode

func ClosureBundler_EvalMode_Values() []*ClosureBundler_EvalMode {
	clinit_ClosureBundler_EvalMode()
	return j2gort.CloneSlice(enumValues_ClosureBundler_EvalMode)
}

func ClosureBundler_EvalMode_ValueOf(name string) *ClosureBundler_EvalMode {
	clinit_ClosureBundler_EvalMode()
	for _, c := range enumValues_ClosureBundler_EvalMode {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.deps.ClosureBundler.EvalMode." + name))
}

// Java method: EvalMode.appendTraditional(java.lang.String s, java.lang.Appendable out, java.lang.@org.jspecify.annotations.Nullable String sourceUrl)
func (this_ *ClosureBundler_EvalMode) AppendTraditional(s string, out j2gort.Ext_java_lang_Appendable, sourceUrl string) {
	if this_.fn_AppendTraditional != nil {
		this_.fn_AppendTraditional(this_, s, out, sourceUrl)
		return
	}
	this_.baseEnum_AppendTraditional(s, out, sourceUrl)
}

func (this_ *ClosureBundler_EvalMode) baseEnum_AppendTraditional(s string, out j2gort.Ext_java_lang_Appendable, sourceUrl string) {
	_ = s
	_ = out
	_ = sourceUrl
	panic("abstract or native method invoked")
}

// Java method: EvalMode.appendGoogModule(java.lang.String s, java.lang.Appendable out, java.lang.String sourceUrl)
func (this_ *ClosureBundler_EvalMode) AppendGoogModule(s string, out j2gort.Ext_java_lang_Appendable, sourceUrl string) {
	if this_.fn_AppendGoogModule != nil {
		this_.fn_AppendGoogModule(this_, s, out, sourceUrl)
		return
	}
	this_.baseEnum_AppendGoogModule(s, out, sourceUrl)
}

func (this_ *ClosureBundler_EvalMode) baseEnum_AppendGoogModule(s string, out j2gort.Ext_java_lang_Appendable, sourceUrl string) {
	_ = s
	_ = out
	_ = sourceUrl
	panic("abstract or native method invoked")
}

func (this_ *ClosureBundler_EvalMode) IsJavaNil_() bool { return this_ == nil }
func (*ClosureBundler_EvalMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.ClosureBundler$EvalMode", "EvalMode"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.ClosureBundler.EvalMode", "com.google.javascript.jscomp.deps.ClosureBundler$EvalMode", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureBundler_EvalMode](v_)
	})
}

// Java enum: com.google.javascript.jscomp.deps.ClosureBundler.EscapeMode
type ClosureBundler_EscapeMode struct {
	j2gort.EnumBase
	fn_Append func(*ClosureBundler_EscapeMode, string, j2gort.Ext_java_lang_Appendable)
}

func (this_ *ClosureBundler_EscapeMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *ClosureBundler_EscapeMode) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*ClosureBundler_EscapeMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*ClosureBundler_EscapeMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *ClosureBundler_EscapeMode) initFields_ClosureBundler_EscapeMode() {
	_ = this_
}

func initEnumClosureBundler_EscapeMode(this_ *ClosureBundler_EscapeMode, enumName_ string, enumOrd_ int32) *ClosureBundler_EscapeMode {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	this_.initFields_ClosureBundler_EscapeMode()
	return this_
}

var enum_ClosureBundler_EscapeMode_ESCAPED *ClosureBundler_EscapeMode

func StaticGet_ClosureBundler_EscapeMode_ESCAPED() *ClosureBundler_EscapeMode {
	clinit_ClosureBundler_EscapeMode()
	return enum_ClosureBundler_EscapeMode_ESCAPED
}

var enum_ClosureBundler_EscapeMode_NORMAL *ClosureBundler_EscapeMode

func StaticGet_ClosureBundler_EscapeMode_NORMAL() *ClosureBundler_EscapeMode {
	clinit_ClosureBundler_EscapeMode()
	return enum_ClosureBundler_EscapeMode_NORMAL
}

var clinitDone_ClosureBundler_EscapeMode uint32
var clinitState_ClosureBundler_EscapeMode j2gort.ClinitState

func clinit_ClosureBundler_EscapeMode() {
	if atomic.LoadUint32(&clinitDone_ClosureBundler_EscapeMode) != 0 {
		return
	}
	clinitSlow_ClosureBundler_EscapeMode()
}

func clinitSlow_ClosureBundler_EscapeMode() {
	if !j2gort.EnterClinit(&clinitState_ClosureBundler_EscapeMode, &clinitDone_ClosureBundler_EscapeMode) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ClosureBundler_EscapeMode, &clinitDone_ClosureBundler_EscapeMode)
	enum_ClosureBundler_EscapeMode_ESCAPED = &ClosureBundler_EscapeMode{}
	enum_ClosureBundler_EscapeMode_ESCAPED.fn_Append = func(this_ *ClosureBundler_EscapeMode, s string, out j2gort.Ext_java_lang_Appendable) {
		_ = s
		_ = out
		SourceCodeEscapers_AppendWithJavascriptEscaper(j2gort.BoxString(s), out)
	}
	initEnumClosureBundler_EscapeMode(enum_ClosureBundler_EscapeMode_ESCAPED, "ESCAPED", int32(0))
	enum_ClosureBundler_EscapeMode_NORMAL = &ClosureBundler_EscapeMode{}
	enum_ClosureBundler_EscapeMode_NORMAL.fn_Append = func(this_ *ClosureBundler_EscapeMode, s string, out j2gort.Ext_java_lang_Appendable) {
		_ = s
		_ = out
		out.Append_1_any(j2gort.BoxString(s))
	}
	initEnumClosureBundler_EscapeMode(enum_ClosureBundler_EscapeMode_NORMAL, "NORMAL", int32(1))
	enumValues_ClosureBundler_EscapeMode = []*ClosureBundler_EscapeMode{enum_ClosureBundler_EscapeMode_ESCAPED, enum_ClosureBundler_EscapeMode_NORMAL}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.ClosureBundler.EscapeMode", []any{enum_ClosureBundler_EscapeMode_ESCAPED, enum_ClosureBundler_EscapeMode_NORMAL})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.ClosureBundler$EscapeMode", []any{enum_ClosureBundler_EscapeMode_ESCAPED, enum_ClosureBundler_EscapeMode_NORMAL})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.deps.ClosureBundler.EscapeMode", "com.google.javascript.jscomp.deps.ClosureBundler$EscapeMode", func() {
		clinit_ClosureBundler_EscapeMode()
	})
}

var enumValues_ClosureBundler_EscapeMode []*ClosureBundler_EscapeMode

func ClosureBundler_EscapeMode_Values() []*ClosureBundler_EscapeMode {
	clinit_ClosureBundler_EscapeMode()
	return j2gort.CloneSlice(enumValues_ClosureBundler_EscapeMode)
}

func ClosureBundler_EscapeMode_ValueOf(name string) *ClosureBundler_EscapeMode {
	clinit_ClosureBundler_EscapeMode()
	for _, c := range enumValues_ClosureBundler_EscapeMode {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.deps.ClosureBundler.EscapeMode." + name))
}

// Java method: EscapeMode.append(java.lang.String s, java.lang.Appendable out)
func (this_ *ClosureBundler_EscapeMode) Append(s string, out j2gort.Ext_java_lang_Appendable) {
	if this_.fn_Append != nil {
		this_.fn_Append(this_, s, out)
		return
	}
	this_.baseEnum_Append(s, out)
}

func (this_ *ClosureBundler_EscapeMode) baseEnum_Append(s string, out j2gort.Ext_java_lang_Appendable) {
	_ = s
	_ = out
	panic("abstract or native method invoked")
}

func (this_ *ClosureBundler_EscapeMode) IsJavaNil_() bool { return this_ == nil }
func (*ClosureBundler_EscapeMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.ClosureBundler$EscapeMode", "EscapeMode"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.ClosureBundler.EscapeMode", "com.google.javascript.jscomp.deps.ClosureBundler$EscapeMode", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureBundler_EscapeMode](v_)
	})
}
