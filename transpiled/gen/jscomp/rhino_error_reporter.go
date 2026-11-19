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
// Port of com/google/javascript/jscomp/RhinoErrorReporter.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java polymorphic class interface: com.google.javascript.jscomp.RhinoErrorReporter
type RhinoErrorReporter interface {
	Implements_jscomp_RhinoErrorReporter()
	IsJavaNil_() bool
	GetRhinoErrorReporterBase() *RhinoErrorReporterBase
	WarningAtLine(string, string, int32, int32)
	ErrorAtLine(string, string, int32, int32)
}

func ToIface_RhinoErrorReporter[T interface {
	RhinoErrorReporter
	comparable
}](v T) RhinoErrorReporter {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_RhinoErrorReporter[T RhinoErrorReporter](v T) RhinoErrorReporter {
	var vi RhinoErrorReporter = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.RhinoErrorReporter
type RhinoErrorReporterBase struct {
	internalReporter ErrorHandler `gson:"internalReporter"`
	idHash_          atomic.Int32
}

func (this_ *RhinoErrorReporterBase) BindSelf_(s any) {
	_ = s
}
func (this_ *RhinoErrorReporterBase) IsJavaNil_() bool { return this_ == nil }
func (*RhinoErrorReporterBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RhinoErrorReporter", "RhinoErrorReporter"
}
func (this_ *RhinoErrorReporterBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RhinoErrorReporter", "com.google.javascript.jscomp.RhinoErrorReporter", func(v_ any) bool {
		return j2gort.InstanceOf[RhinoErrorReporter](v_)
	})
}
func (b *RhinoErrorReporterBase) GetRhinoErrorReporterBase() *RhinoErrorReporterBase { return b }
func (*RhinoErrorReporterBase) Implements_jscomp_RhinoErrorReporter()                {}

var static_RhinoErrorReporter_PARSE_ERROR *DiagnosticType

func StaticGet_RhinoErrorReporter_PARSE_ERROR() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_PARSE_ERROR
}
func StaticSet_RhinoErrorReporter_PARSE_ERROR(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_PARSE_ERROR = v
}

var static_RhinoErrorReporter_TYPE_PARSE_ERROR *DiagnosticType

func StaticGet_RhinoErrorReporter_TYPE_PARSE_ERROR() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_TYPE_PARSE_ERROR
}
func StaticSet_RhinoErrorReporter_TYPE_PARSE_ERROR(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_TYPE_PARSE_ERROR = v
}

var static_RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR *DiagnosticType

func StaticGet_RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR
}
func StaticSet_RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR = v
}

var static_RhinoErrorReporter_UNRECOGNIZED_TYPEOF_ERROR *DiagnosticType

func StaticGet_RhinoErrorReporter_UNRECOGNIZED_TYPEOF_ERROR() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_UNRECOGNIZED_TYPEOF_ERROR
}
func StaticSet_RhinoErrorReporter_UNRECOGNIZED_TYPEOF_ERROR(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_UNRECOGNIZED_TYPEOF_ERROR = v
}

var static_RhinoErrorReporter_CYCLIC_INHERITANCE_ERROR *DiagnosticType

func StaticGet_RhinoErrorReporter_CYCLIC_INHERITANCE_ERROR() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_CYCLIC_INHERITANCE_ERROR
}
func StaticSet_RhinoErrorReporter_CYCLIC_INHERITANCE_ERROR(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_CYCLIC_INHERITANCE_ERROR = v
}

var static_RhinoErrorReporter_JSDOC_MISSING_BRACES_WARNING *DiagnosticType

func StaticGet_RhinoErrorReporter_JSDOC_MISSING_BRACES_WARNING() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_JSDOC_MISSING_BRACES_WARNING
}
func StaticSet_RhinoErrorReporter_JSDOC_MISSING_BRACES_WARNING(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_JSDOC_MISSING_BRACES_WARNING = v
}

var static_RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING *DiagnosticType

func StaticGet_RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING
}
func StaticSet_RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING = v
}

var static_RhinoErrorReporter_JSDOC_IMPORT_TYPE_WARNING *DiagnosticType

func StaticGet_RhinoErrorReporter_JSDOC_IMPORT_TYPE_WARNING() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_JSDOC_IMPORT_TYPE_WARNING
}
func StaticSet_RhinoErrorReporter_JSDOC_IMPORT_TYPE_WARNING(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_JSDOC_IMPORT_TYPE_WARNING = v
}

var static_RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS *DiagnosticType

func StaticGet_RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS
}
func StaticSet_RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS = v
}

var static_RhinoErrorReporter_TRAILING_COMMA *DiagnosticType

func StaticGet_RhinoErrorReporter_TRAILING_COMMA() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_TRAILING_COMMA
}
func StaticSet_RhinoErrorReporter_TRAILING_COMMA(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_TRAILING_COMMA = v
}

var static_RhinoErrorReporter_DUPLICATE_PARAM *DiagnosticType

func StaticGet_RhinoErrorReporter_DUPLICATE_PARAM() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_DUPLICATE_PARAM
}
func StaticSet_RhinoErrorReporter_DUPLICATE_PARAM(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_DUPLICATE_PARAM = v
}

var static_RhinoErrorReporter_DUPLICATE_VISIBILITY *DiagnosticType

func StaticGet_RhinoErrorReporter_DUPLICATE_VISIBILITY() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_DUPLICATE_VISIBILITY
}
func StaticSet_RhinoErrorReporter_DUPLICATE_VISIBILITY(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_DUPLICATE_VISIBILITY = v
}

var static_RhinoErrorReporter_UNNECESSARY_ESCAPE *DiagnosticType

func StaticGet_RhinoErrorReporter_UNNECESSARY_ESCAPE() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_UNNECESSARY_ESCAPE
}
func StaticSet_RhinoErrorReporter_UNNECESSARY_ESCAPE(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_UNNECESSARY_ESCAPE = v
}

var static_RhinoErrorReporter_INVALID_PARAM *DiagnosticType

func StaticGet_RhinoErrorReporter_INVALID_PARAM() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_INVALID_PARAM
}
func StaticSet_RhinoErrorReporter_INVALID_PARAM(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_INVALID_PARAM = v
}

var static_RhinoErrorReporter_BAD_JSDOC_ANNOTATION *DiagnosticType

func StaticGet_RhinoErrorReporter_BAD_JSDOC_ANNOTATION() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_BAD_JSDOC_ANNOTATION
}
func StaticSet_RhinoErrorReporter_BAD_JSDOC_ANNOTATION(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_BAD_JSDOC_ANNOTATION = v
}

var static_RhinoErrorReporter_INVALID_ES3_PROP_NAME *DiagnosticType

func StaticGet_RhinoErrorReporter_INVALID_ES3_PROP_NAME() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_INVALID_ES3_PROP_NAME
}
func StaticSet_RhinoErrorReporter_INVALID_ES3_PROP_NAME(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_INVALID_ES3_PROP_NAME = v
}

var static_RhinoErrorReporter_PARSE_TREE_TOO_DEEP *DiagnosticType

func StaticGet_RhinoErrorReporter_PARSE_TREE_TOO_DEEP() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_PARSE_TREE_TOO_DEEP
}
func StaticSet_RhinoErrorReporter_PARSE_TREE_TOO_DEEP(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_PARSE_TREE_TOO_DEEP = v
}

var static_RhinoErrorReporter_INVALID_OCTAL_LITERAL *DiagnosticType

func StaticGet_RhinoErrorReporter_INVALID_OCTAL_LITERAL() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_INVALID_OCTAL_LITERAL
}
func StaticSet_RhinoErrorReporter_INVALID_OCTAL_LITERAL(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_INVALID_OCTAL_LITERAL = v
}

var static_RhinoErrorReporter_STRING_CONTINUATION *DiagnosticType

func StaticGet_RhinoErrorReporter_STRING_CONTINUATION() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_STRING_CONTINUATION
}
func StaticSet_RhinoErrorReporter_STRING_CONTINUATION(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_STRING_CONTINUATION = v
}

var static_RhinoErrorReporter_LANGUAGE_FEATURE *DiagnosticType

func StaticGet_RhinoErrorReporter_LANGUAGE_FEATURE() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_LANGUAGE_FEATURE
}
func StaticSet_RhinoErrorReporter_LANGUAGE_FEATURE(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_LANGUAGE_FEATURE = v
}

var static_RhinoErrorReporter_UNSUPPORTED_LANGUAGE_FEATURE *DiagnosticType

func StaticGet_RhinoErrorReporter_UNSUPPORTED_LANGUAGE_FEATURE() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_UNSUPPORTED_LANGUAGE_FEATURE
}
func StaticSet_RhinoErrorReporter_UNSUPPORTED_LANGUAGE_FEATURE(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_UNSUPPORTED_LANGUAGE_FEATURE = v
}

var static_RhinoErrorReporter_UNSUPPORTED_BOUNDED_GENERIC_TYPES *DiagnosticType

func StaticGet_RhinoErrorReporter_UNSUPPORTED_BOUNDED_GENERIC_TYPES() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_UNSUPPORTED_BOUNDED_GENERIC_TYPES
}
func StaticSet_RhinoErrorReporter_UNSUPPORTED_BOUNDED_GENERIC_TYPES(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_UNSUPPORTED_BOUNDED_GENERIC_TYPES = v
}

var static_RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR *DiagnosticType

func StaticGet_RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR
}
func StaticSet_RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR = v
}

var static_RhinoErrorReporter_CLOSURE_UNAWARE_ANNOTATION_PRESENT *DiagnosticType

func StaticGet_RhinoErrorReporter_CLOSURE_UNAWARE_ANNOTATION_PRESENT() *DiagnosticType {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_CLOSURE_UNAWARE_ANNOTATION_PRESENT
}
func StaticSet_RhinoErrorReporter_CLOSURE_UNAWARE_ANNOTATION_PRESENT(v *DiagnosticType) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_CLOSURE_UNAWARE_ANNOTATION_PRESENT = v
}

var static_RhinoErrorReporter_typeMap j2gort.Ext_com_google_common_collect_ImmutableMap

func StaticGet_RhinoErrorReporter_TypeMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_RhinoErrorReporter()
	return static_RhinoErrorReporter_typeMap
}
func StaticSet_RhinoErrorReporter_TypeMap(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_RhinoErrorReporter()
	static_RhinoErrorReporter_typeMap = v
}

var clinitDone_RhinoErrorReporter uint32
var clinitState_RhinoErrorReporter j2gort.ClinitState

func clinit_RhinoErrorReporter() {
	if atomic.LoadUint32(&clinitDone_RhinoErrorReporter) != 0 {
		return
	}
	clinitSlow_RhinoErrorReporter()
}

func clinitSlow_RhinoErrorReporter() {
	if !j2gort.EnterClinit(&clinitState_RhinoErrorReporter, &clinitDone_RhinoErrorReporter) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RhinoErrorReporter, &clinitDone_RhinoErrorReporter)
	static_RhinoErrorReporter_PARSE_ERROR = DiagnosticType_Error("JSC_PARSE_ERROR", "Parse error. {0}")
	static_RhinoErrorReporter_TYPE_PARSE_ERROR = DiagnosticType_Warning("JSC_TYPE_PARSE_ERROR", "{0}")
	static_RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR = DiagnosticType_Warning("JSC_UNRECOGNIZED_TYPE_ERROR", "{0}")
	static_RhinoErrorReporter_UNRECOGNIZED_TYPEOF_ERROR = DiagnosticType_Warning("JSC_UNRECOGNIZED_TYPEOF_ERROR", "{0}")
	static_RhinoErrorReporter_CYCLIC_INHERITANCE_ERROR = DiagnosticType_Warning("JSC_CYCLIC_INHERITANCE_ERROR", "{0}")
	static_RhinoErrorReporter_JSDOC_MISSING_BRACES_WARNING = DiagnosticType_Disabled("JSC_JSDOC_MISSING_BRACES_WARNING", "{0}")
	static_RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING = DiagnosticType_Disabled("JSC_JSDOC_MISSING_TYPE_WARNING", "{0}")
	static_RhinoErrorReporter_JSDOC_IMPORT_TYPE_WARNING = DiagnosticType_Disabled("JSC_JSDOC_IMPORT_TYPE_WARNING", "{0}")
	static_RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS = DiagnosticType_Warning("JSC_TOO_MANY_TEMPLATE_PARAMS", "{0}")
	static_RhinoErrorReporter_TRAILING_COMMA = DiagnosticType_Error("JSC_TRAILING_COMMA", "Parse error. IE8 (and below) will parse trailing commas in array and object literals incorrectly. If you are targeting newer versions of JS, set the appropriate language_in option.")
	static_RhinoErrorReporter_DUPLICATE_PARAM = DiagnosticType_Error("JSC_DUPLICATE_PARAM", "Parse error. {0}")
	static_RhinoErrorReporter_DUPLICATE_VISIBILITY = DiagnosticType_Warning("JSC_DUPLICATE_VISIBILITY", "{0}")
	static_RhinoErrorReporter_UNNECESSARY_ESCAPE = DiagnosticType_Disabled("JSC_UNNECESSARY_ESCAPE", "Parse error. {0}")
	static_RhinoErrorReporter_INVALID_PARAM = DiagnosticType_Warning("JSC_INVALID_PARAM", "Parse error. {0}")
	static_RhinoErrorReporter_BAD_JSDOC_ANNOTATION = DiagnosticType_Warning("JSC_BAD_JSDOC_ANNOTATION", "Parse error. {0}")
	static_RhinoErrorReporter_INVALID_ES3_PROP_NAME = DiagnosticType_Warning("JSC_INVALID_ES3_PROP_NAME", "Keywords and reserved words are not allowed as unquoted property names in older versions of JavaScript. If you are targeting newer versions of JavaScript, set the appropriate language_in option.")
	static_RhinoErrorReporter_PARSE_TREE_TOO_DEEP = DiagnosticType_Error("JSC_PARSE_TREE_TOO_DEEP", "Parse tree too deep.")
	static_RhinoErrorReporter_INVALID_OCTAL_LITERAL = DiagnosticType_Warning("JSC_INVALID_OCTAL_LITERAL", "This style of octal literal is not supported in strict mode.")
	static_RhinoErrorReporter_STRING_CONTINUATION = DiagnosticType_Disabled("JSC_STRING_CONTINUATION", "{0}")
	static_RhinoErrorReporter_LANGUAGE_FEATURE = DiagnosticType_Error("JSC_LANGUAGE_FEATURE", "{0}.")
	static_RhinoErrorReporter_UNSUPPORTED_LANGUAGE_FEATURE = DiagnosticType_Error("JSC_UNSUPPORTED_LANGUAGE_FEATURE", "{0}.")
	static_RhinoErrorReporter_UNSUPPORTED_BOUNDED_GENERIC_TYPES = DiagnosticType_Error("JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES", "Bounded generic semantics are currently still in development")
	static_RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR = DiagnosticType_Error("JSC_BOUNDED_GENERIC_TYPE_ERROR", "Bounded generic type error. {0} assigned to template type {1} is not a subtype of bound {2}")
	static_RhinoErrorReporter_CLOSURE_UNAWARE_ANNOTATION_PRESENT = DiagnosticType_Disabled("JSC_CLOSURE_UNAWARE_ANNOTATION_PRESENT", pkg_rhino.StaticGet_Msg_JSDOC_CLOSURE_UNAWARE_CODE_INVALID().Format_0())
	static_RhinoErrorReporter_typeMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder().Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("Trailing comma is not legal in an ECMA-262 object initializer")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_TRAILING_COMMA()))).Put_2(j2gort.PtrToAny(RhinoErrorReporter_ReplacePlaceHolders("Duplicate parameter name \"{0}\"")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_DUPLICATE_PARAM()))).Put_2(j2gort.PtrToAny(RhinoErrorReporter_ReplacePlaceHolders(pkg_rhino.StaticGet_Msg_JSDOC_EXTRA_VISIBILITY().Format_0())), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_DUPLICATE_VISIBILITY()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("Unnecessary escape:.*")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_UNNECESSARY_ESCAPE()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("^invalid param name.*")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_INVALID_PARAM()))).Put_2(j2gort.PtrToAny(RhinoErrorReporter_ReplacePlaceHolders(pkg_rhino.StaticGet_Msg_BAD_JSDOC_TAG().Format_0())), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_BAD_JSDOC_ANNOTATION()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("^Keywords and reserved words are not allowed as unquoted property.*")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_INVALID_ES3_PROP_NAME()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("^Too many template parameters\n.*")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1(".*Type annotations should have curly braces.*")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_JSDOC_MISSING_BRACES_WARNING()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("Missing type declaration\\.")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1(".*Unknown type.*")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1(".*Unknown type.*\n.*")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("^Missing type for `typeof` value.*")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_UNRECOGNIZED_TYPEOF_ERROR()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("^Cycle detected in inheritance chain of type .*")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_CYCLIC_INHERITANCE_ERROR()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("^Bad type annotation. Import in typedef.*")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_JSDOC_IMPORT_TYPE_WARNING()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("^Bad type annotation.*")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_TYPE_PARSE_ERROR()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("constructed type must be an object type")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_TYPE_PARSE_ERROR()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("Too deep recursion while parsing")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_PARSE_TREE_TOO_DEEP()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("^Octal .*literal.*")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_INVALID_OCTAL_LITERAL()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("^String continuations.*")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_STRING_CONTINUATION()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("^This language feature is only supported for .*")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_LANGUAGE_FEATURE()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("^This language feature is not currently supported by the compiler: .*")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_UNSUPPORTED_LANGUAGE_FEATURE()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("Bounded generic semantics are currently still in development")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_UNSUPPORTED_BOUNDED_GENERIC_TYPES()))).Put_2(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("^Bounded generic type error.*")), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR()))).Put_2(j2gort.PtrToAny(RhinoErrorReporter_ReplacePlaceHolders(pkg_rhino.StaticGet_Msg_JSDOC_CLOSURE_UNAWARE_CODE_INVALID().Format_0())), j2gort.PtrToAny(StaticGet_RhinoErrorReporter_CLOSURE_UNAWARE_ANNOTATION_PRESENT()))).BuildOrThrow())
}

func (this_ *RhinoErrorReporterBase) initFields_RhinoErrorReporter() {
	_ = this_
}

// Java constructor: RhinoErrorReporter(com.google.javascript.jscomp.ErrorHandler internalReporter)
func NewRhinoErrorReporter(internalReporter ErrorHandler) RhinoErrorReporter {
	clinit_RhinoErrorReporter()
	this_ := &RhinoErrorReporterBase{}
	InitRhinoErrorReporter(this_, internalReporter)
	return this_
}

func InitRhinoErrorReporter(this_ *RhinoErrorReporterBase, internalReporter ErrorHandler) {
	_ = internalReporter
	this_.initFields_RhinoErrorReporter()
	this_.internalReporter = internalReporter
}

// Java static method: RhinoErrorReporter.replacePlaceHolders(java.lang.String s)
func RhinoErrorReporter_ReplacePlaceHolders(s string) *j2gort.Ext_java_util_regex_Pattern {
	clinit_RhinoErrorReporter()
	_ = s
	s = j2gort.Ext_java_util_regex_Pattern_Quote(s)
	return j2gort.Ext_java_util_regex_Pattern_Compile_1(j2gort.JStrReplaceAll(s, "\\{\\d+\\}", "\\\\E.*\\\\Q"))
}

// Java static method: RhinoErrorReporter.forOldRhino(com.google.javascript.jscomp.ErrorHandler internalReporter)
func RhinoErrorReporter_ForOldRhino(internalReporter ErrorHandler) pkg_rhino.ErrorReporter {
	clinit_RhinoErrorReporter()
	_ = internalReporter
	return NewRhinoErrorReporter_OldRhinoErrorReporter(internalReporter)
}

// Java method: RhinoErrorReporter.warningAtLine(java.lang.String message, java.lang.String sourceName, int line, int lineOffset)
func (this_ *RhinoErrorReporterBase) WarningAtLine(message string, sourceName string, line int32, lineOffset int32) {
	_ = message
	_ = sourceName
	_ = line
	_ = lineOffset
	this_.internalReporter.Report_2(nil, RhinoErrorReporter_MakeError(message, sourceName, line, lineOffset, StaticGet_CheckLevel_WARNING()))
}

// Java method: RhinoErrorReporter.errorAtLine(java.lang.String message, java.lang.String sourceName, int line, int lineOffset)
func (this_ *RhinoErrorReporterBase) ErrorAtLine(message string, sourceName string, line int32, lineOffset int32) {
	_ = message
	_ = sourceName
	_ = line
	_ = lineOffset
	this_.internalReporter.Report_2(nil, RhinoErrorReporter_MakeError(message, sourceName, line, lineOffset, StaticGet_CheckLevel_ERROR()))
}

// Java static method: RhinoErrorReporter.mapError(java.lang.String message)
func RhinoErrorReporter_MapError(message string) *DiagnosticType {
	clinit_RhinoErrorReporter()
	_ = message
	{
		cur_1 := j2gort.NewMapEntryCursor(StaticGet_RhinoErrorReporter_TypeMap())
		for cur_1.NextEntry() {
			if j2gort.CastPtr[j2gort.Ext_java_util_regex_Pattern](cur_1.GetKey()).Matcher(j2gort.BoxString(message)).Matches() {
				return j2gort.CastPtr[DiagnosticType](cur_1.GetValue())
			}
		}
	}
	return nil
}

// Java static method: RhinoErrorReporter.makeError(java.lang.String message, java.lang.String sourceName, int line, int lineOffset, com.google.javascript.jscomp.CheckLevel defaultLevel)
func RhinoErrorReporter_MakeError(message string, sourceName string, line int32, lineOffset int32, defaultLevel *CheckLevel) *JSError {
	clinit_RhinoErrorReporter()
	_ = message
	_ = sourceName
	_ = line
	_ = lineOffset
	_ = defaultLevel
	var type_ *DiagnosticType = RhinoErrorReporter_MapError(message)
	_ = type_
	var builder *JSError_Builder = JSError_S_Builder(func() *DiagnosticType {
		if type_ != nil {
			return type_
		}
		return StaticGet_RhinoErrorReporter_PARSE_ERROR()
	}(), []string{message}).SetSourceLocation(sourceName, line, lineOffset)
	_ = builder
	if type_ == nil {
		builder.SetLevel(defaultLevel)
	}
	return builder.Build()
}

// Java class: com.google.javascript.jscomp.RhinoErrorReporter.OldRhinoErrorReporter
type RhinoErrorReporter_OldRhinoErrorReporter struct {
	RhinoErrorReporterBase
}

func (this_ *RhinoErrorReporter_OldRhinoErrorReporter) BindSelf_(s any) {
	_ = s
}
func (this_ *RhinoErrorReporter_OldRhinoErrorReporter) IsJavaNil_() bool { return this_ == nil }
func (*RhinoErrorReporter_OldRhinoErrorReporter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RhinoErrorReporter$OldRhinoErrorReporter", "OldRhinoErrorReporter"
}
func (this_ *RhinoErrorReporter_OldRhinoErrorReporter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.RhinoErrorReporterBase.IdentityHash_()
}
func (*RhinoErrorReporter_OldRhinoErrorReporter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RhinoErrorReporter.OldRhinoErrorReporter", "com.google.javascript.jscomp.RhinoErrorReporter$OldRhinoErrorReporter", func(v_ any) bool {
		return j2gort.InstanceOf[*RhinoErrorReporter_OldRhinoErrorReporter](v_)
	})
}
func (*RhinoErrorReporter_OldRhinoErrorReporter) Implements_jscomp_RhinoErrorReporter() {}
func (*RhinoErrorReporter_OldRhinoErrorReporter) Implements_rhino_ErrorReporter()       {}

func (this_ *RhinoErrorReporter_OldRhinoErrorReporter) initFields_RhinoErrorReporter_OldRhinoErrorReporter() {
	_ = this_
}

// Java constructor: OldRhinoErrorReporter(com.google.javascript.jscomp.ErrorHandler internalReporter)
func NewRhinoErrorReporter_OldRhinoErrorReporter(internalReporter ErrorHandler) *RhinoErrorReporter_OldRhinoErrorReporter {
	this_ := &RhinoErrorReporter_OldRhinoErrorReporter{}
	InitRhinoErrorReporter_OldRhinoErrorReporter(this_, internalReporter)
	return this_
}

func InitRhinoErrorReporter_OldRhinoErrorReporter(this_ *RhinoErrorReporter_OldRhinoErrorReporter, internalReporter ErrorHandler) {
	_ = internalReporter
	InitRhinoErrorReporter(&this_.RhinoErrorReporterBase, internalReporter)
	this_.initFields_RhinoErrorReporter_OldRhinoErrorReporter()
}

// Java method: OldRhinoErrorReporter.error(java.lang.String message, java.lang.String sourceName, int line, int lineOffset)
func (this_ *RhinoErrorReporter_OldRhinoErrorReporter) M_Error(message string, sourceName string, line int32, lineOffset int32) {
	_ = message
	_ = sourceName
	_ = line
	_ = lineOffset
	this_.RhinoErrorReporterBase.ErrorAtLine(message, sourceName, line, lineOffset)
}

// Java method: OldRhinoErrorReporter.warning(java.lang.String message, java.lang.String sourceName, int line, int lineOffset)
func (this_ *RhinoErrorReporter_OldRhinoErrorReporter) Warning(message string, sourceName string, line int32, lineOffset int32) {
	_ = message
	_ = sourceName
	_ = line
	_ = lineOffset
	this_.RhinoErrorReporterBase.WarningAtLine(message, sourceName, line, lineOffset)
}
