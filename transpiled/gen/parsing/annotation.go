/*
 * Copyright 2010 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/Annotation.java (generated mechanically by j2go).

package parsing

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.jscomp.parsing.Annotation
type Annotation uint8

const (
	enum_Annotation_ABSTRACT                     Annotation = 1
	enum_Annotation_ALTERNATE_MESSAGE_ID         Annotation = 2
	enum_Annotation_AUTHOR                       Annotation = 3
	enum_Annotation_CLOSURE_PRIMITIVE            Annotation = 4
	enum_Annotation_COLLAPSIBLE_OR_BREAK_MY_CODE Annotation = 5
	enum_Annotation_CONSTANT                     Annotation = 6
	enum_Annotation_CONSTRUCTOR                  Annotation = 7
	enum_Annotation_CUSTOM_ELEMENT               Annotation = 8
	enum_Annotation_RECORD                       Annotation = 9
	enum_Annotation_DEFINE                       Annotation = 10
	enum_Annotation_DEPRECATED                   Annotation = 11
	enum_Annotation_DESC                         Annotation = 12
	enum_Annotation_DICT                         Annotation = 13
	enum_Annotation_ENCOURAGE_INLINING           Annotation = 14
	enum_Annotation_ENHANCE                      Annotation = 15
	enum_Annotation_ENUM                         Annotation = 16
	enum_Annotation_EXTENDS                      Annotation = 17
	enum_Annotation_EXTERNS                      Annotation = 18
	enum_Annotation_EXPORT                       Annotation = 19
	enum_Annotation_FILE_OVERVIEW                Annotation = 20
	enum_Annotation_FINAL                        Annotation = 21
	enum_Annotation_IDGENERATOR                  Annotation = 22
	enum_Annotation_IMPLEMENTS                   Annotation = 23
	enum_Annotation_IMPLICIT_CAST                Annotation = 24
	enum_Annotation_INHERIT_DOC                  Annotation = 25
	enum_Annotation_INTERFACE                    Annotation = 26
	enum_Annotation_JSX                          Annotation = 27
	enum_Annotation_JSX_FRAGMENT                 Annotation = 28
	enum_Annotation_LEGACY_NAMED_TS_EXPORT       Annotation = 29
	enum_Annotation_LENDS                        Annotation = 30
	enum_Annotation_LICENSE                      Annotation = 31
	enum_Annotation_LOG_TYPE_IN_COMPILER         Annotation = 32
	enum_Annotation_MAY_HAVE_EXTRA_EDGE          Annotation = 33
	enum_Annotation_MEANING                      Annotation = 34
	enum_Annotation_MIXIN_CLASS                  Annotation = 35
	enum_Annotation_MIXIN_FUNCTION               Annotation = 36
	enum_Annotation_MODIFIES                     Annotation = 37
	enum_Annotation_MODS                         Annotation = 38
	enum_Annotation_NG_INJECT                    Annotation = 39
	enum_Annotation_NO_COLLAPSE                  Annotation = 40
	enum_Annotation_NO_COMPILE                   Annotation = 41
	enum_Annotation_NO_COVERAGE                  Annotation = 42
	enum_Annotation_NO_DTS                       Annotation = 43
	enum_Annotation_NO_INLINE                    Annotation = 44
	enum_Annotation_NO_SIDE_EFFECTS              Annotation = 45
	enum_Annotation_NOT_IMPLEMENTED              Annotation = 46
	enum_Annotation_OVERRIDE                     Annotation = 47
	enum_Annotation_PACKAGE                      Annotation = 48
	enum_Annotation_PARAM                        Annotation = 49
	enum_Annotation_POLYMER                      Annotation = 50
	enum_Annotation_POLYMER_BEHAVIOR             Annotation = 51
	enum_Annotation_PRESERVE                     Annotation = 52
	enum_Annotation_PRIVATE                      Annotation = 53
	enum_Annotation_PROTECTED                    Annotation = 54
	enum_Annotation_PROVIDE_GOOG                 Annotation = 55
	enum_Annotation_PROVIDE_ALREADY_PROVIDED     Annotation = 56
	enum_Annotation_PUBLIC                       Annotation = 57
	enum_Annotation_PURE_OR_BREAK_MY_CODE        Annotation = 58
	enum_Annotation_REQUIRE_INLINING             Annotation = 59
	enum_Annotation_RETURN                       Annotation = 60
	enum_Annotation_SASS_GENERATED_CSS_TS        Annotation = 61
	enum_Annotation_SEE                          Annotation = 62
	enum_Annotation_SOY_MODULE                   Annotation = 63
	enum_Annotation_SOY_TEMPLATE                 Annotation = 64
	enum_Annotation_STRUCT                       Annotation = 65
	enum_Annotation_SUPPRESS                     Annotation = 66
	enum_Annotation_TEMPLATE                     Annotation = 67
	enum_Annotation_CLOSURE_UNAWARE_CODE         Annotation = 68
	enum_Annotation_THIS                         Annotation = 69
	enum_Annotation_THROWS                       Annotation = 70
	enum_Annotation_TYPE                         Annotation = 71
	enum_Annotation_TYPEDEF                      Annotation = 72
	enum_Annotation_TYPE_SUMMARY                 Annotation = 73
	enum_Annotation_UNRESTRICTED                 Annotation = 74
	enum_Annotation_USED_VIA_DOT_CONSTRUCTOR     Annotation = 75
	enum_Annotation_WIZACTION                    Annotation = 76
	enum_Annotation_TS_TYPE                      Annotation = 77
	enum_Annotation_WIZ_ANALYZER                 Annotation = 78
	enum_Annotation_WIZCALLBACK                  Annotation = 79
)

var enumBaseTable_Annotation = [...]j2gort.EnumBase{
	{Name_: "ABSTRACT", Ordinal_: 0},
	{Name_: "ALTERNATE_MESSAGE_ID", Ordinal_: 1},
	{Name_: "AUTHOR", Ordinal_: 2},
	{Name_: "CLOSURE_PRIMITIVE", Ordinal_: 3},
	{Name_: "COLLAPSIBLE_OR_BREAK_MY_CODE", Ordinal_: 4},
	{Name_: "CONSTANT", Ordinal_: 5},
	{Name_: "CONSTRUCTOR", Ordinal_: 6},
	{Name_: "CUSTOM_ELEMENT", Ordinal_: 7},
	{Name_: "RECORD", Ordinal_: 8},
	{Name_: "DEFINE", Ordinal_: 9},
	{Name_: "DEPRECATED", Ordinal_: 10},
	{Name_: "DESC", Ordinal_: 11},
	{Name_: "DICT", Ordinal_: 12},
	{Name_: "ENCOURAGE_INLINING", Ordinal_: 13},
	{Name_: "ENHANCE", Ordinal_: 14},
	{Name_: "ENUM", Ordinal_: 15},
	{Name_: "EXTENDS", Ordinal_: 16},
	{Name_: "EXTERNS", Ordinal_: 17},
	{Name_: "EXPORT", Ordinal_: 18},
	{Name_: "FILE_OVERVIEW", Ordinal_: 19},
	{Name_: "FINAL", Ordinal_: 20},
	{Name_: "IDGENERATOR", Ordinal_: 21},
	{Name_: "IMPLEMENTS", Ordinal_: 22},
	{Name_: "IMPLICIT_CAST", Ordinal_: 23},
	{Name_: "INHERIT_DOC", Ordinal_: 24},
	{Name_: "INTERFACE", Ordinal_: 25},
	{Name_: "JSX", Ordinal_: 26},
	{Name_: "JSX_FRAGMENT", Ordinal_: 27},
	{Name_: "LEGACY_NAMED_TS_EXPORT", Ordinal_: 28},
	{Name_: "LENDS", Ordinal_: 29},
	{Name_: "LICENSE", Ordinal_: 30},
	{Name_: "LOG_TYPE_IN_COMPILER", Ordinal_: 31},
	{Name_: "MAY_HAVE_EXTRA_EDGE", Ordinal_: 32},
	{Name_: "MEANING", Ordinal_: 33},
	{Name_: "MIXIN_CLASS", Ordinal_: 34},
	{Name_: "MIXIN_FUNCTION", Ordinal_: 35},
	{Name_: "MODIFIES", Ordinal_: 36},
	{Name_: "MODS", Ordinal_: 37},
	{Name_: "NG_INJECT", Ordinal_: 38},
	{Name_: "NO_COLLAPSE", Ordinal_: 39},
	{Name_: "NO_COMPILE", Ordinal_: 40},
	{Name_: "NO_COVERAGE", Ordinal_: 41},
	{Name_: "NO_DTS", Ordinal_: 42},
	{Name_: "NO_INLINE", Ordinal_: 43},
	{Name_: "NO_SIDE_EFFECTS", Ordinal_: 44},
	{Name_: "NOT_IMPLEMENTED", Ordinal_: 45},
	{Name_: "OVERRIDE", Ordinal_: 46},
	{Name_: "PACKAGE", Ordinal_: 47},
	{Name_: "PARAM", Ordinal_: 48},
	{Name_: "POLYMER", Ordinal_: 49},
	{Name_: "POLYMER_BEHAVIOR", Ordinal_: 50},
	{Name_: "PRESERVE", Ordinal_: 51},
	{Name_: "PRIVATE", Ordinal_: 52},
	{Name_: "PROTECTED", Ordinal_: 53},
	{Name_: "PROVIDE_GOOG", Ordinal_: 54},
	{Name_: "PROVIDE_ALREADY_PROVIDED", Ordinal_: 55},
	{Name_: "PUBLIC", Ordinal_: 56},
	{Name_: "PURE_OR_BREAK_MY_CODE", Ordinal_: 57},
	{Name_: "REQUIRE_INLINING", Ordinal_: 58},
	{Name_: "RETURN", Ordinal_: 59},
	{Name_: "SASS_GENERATED_CSS_TS", Ordinal_: 60},
	{Name_: "SEE", Ordinal_: 61},
	{Name_: "SOY_MODULE", Ordinal_: 62},
	{Name_: "SOY_TEMPLATE", Ordinal_: 63},
	{Name_: "STRUCT", Ordinal_: 64},
	{Name_: "SUPPRESS", Ordinal_: 65},
	{Name_: "TEMPLATE", Ordinal_: 66},
	{Name_: "CLOSURE_UNAWARE_CODE", Ordinal_: 67},
	{Name_: "THIS", Ordinal_: 68},
	{Name_: "THROWS", Ordinal_: 69},
	{Name_: "TYPE", Ordinal_: 70},
	{Name_: "TYPEDEF", Ordinal_: 71},
	{Name_: "TYPE_SUMMARY", Ordinal_: 72},
	{Name_: "UNRESTRICTED", Ordinal_: 73},
	{Name_: "USED_VIA_DOT_CONSTRUCTOR", Ordinal_: 74},
	{Name_: "WIZACTION", Ordinal_: 75},
	{Name_: "TS_TYPE", Ordinal_: 76},
	{Name_: "WIZ_ANALYZER", Ordinal_: 77},
	{Name_: "WIZCALLBACK", Ordinal_: 78},
}
var enumValues_Annotation = []Annotation{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79}

func StaticGet_Annotation_ABSTRACT() Annotation {
	return enum_Annotation_ABSTRACT
}
func StaticGet_Annotation_ALTERNATE_MESSAGE_ID() Annotation {
	return enum_Annotation_ALTERNATE_MESSAGE_ID
}
func StaticGet_Annotation_AUTHOR() Annotation {
	return enum_Annotation_AUTHOR
}
func StaticGet_Annotation_CLOSURE_PRIMITIVE() Annotation {
	return enum_Annotation_CLOSURE_PRIMITIVE
}
func StaticGet_Annotation_COLLAPSIBLE_OR_BREAK_MY_CODE() Annotation {
	return enum_Annotation_COLLAPSIBLE_OR_BREAK_MY_CODE
}
func StaticGet_Annotation_CONSTANT() Annotation {
	return enum_Annotation_CONSTANT
}
func StaticGet_Annotation_CONSTRUCTOR() Annotation {
	return enum_Annotation_CONSTRUCTOR
}
func StaticGet_Annotation_CUSTOM_ELEMENT() Annotation {
	return enum_Annotation_CUSTOM_ELEMENT
}
func StaticGet_Annotation_RECORD() Annotation {
	return enum_Annotation_RECORD
}
func StaticGet_Annotation_DEFINE() Annotation {
	return enum_Annotation_DEFINE
}
func StaticGet_Annotation_DEPRECATED() Annotation {
	return enum_Annotation_DEPRECATED
}
func StaticGet_Annotation_DESC() Annotation {
	return enum_Annotation_DESC
}
func StaticGet_Annotation_DICT() Annotation {
	return enum_Annotation_DICT
}
func StaticGet_Annotation_ENCOURAGE_INLINING() Annotation {
	return enum_Annotation_ENCOURAGE_INLINING
}
func StaticGet_Annotation_ENHANCE() Annotation {
	return enum_Annotation_ENHANCE
}
func StaticGet_Annotation_ENUM() Annotation {
	return enum_Annotation_ENUM
}
func StaticGet_Annotation_EXTENDS() Annotation {
	return enum_Annotation_EXTENDS
}
func StaticGet_Annotation_EXTERNS() Annotation {
	return enum_Annotation_EXTERNS
}
func StaticGet_Annotation_EXPORT() Annotation {
	return enum_Annotation_EXPORT
}
func StaticGet_Annotation_FILE_OVERVIEW() Annotation {
	return enum_Annotation_FILE_OVERVIEW
}
func StaticGet_Annotation_FINAL() Annotation {
	return enum_Annotation_FINAL
}
func StaticGet_Annotation_IDGENERATOR() Annotation {
	return enum_Annotation_IDGENERATOR
}
func StaticGet_Annotation_IMPLEMENTS() Annotation {
	return enum_Annotation_IMPLEMENTS
}
func StaticGet_Annotation_IMPLICIT_CAST() Annotation {
	return enum_Annotation_IMPLICIT_CAST
}
func StaticGet_Annotation_INHERIT_DOC() Annotation {
	return enum_Annotation_INHERIT_DOC
}
func StaticGet_Annotation_INTERFACE() Annotation {
	return enum_Annotation_INTERFACE
}
func StaticGet_Annotation_JSX() Annotation {
	return enum_Annotation_JSX
}
func StaticGet_Annotation_JSX_FRAGMENT() Annotation {
	return enum_Annotation_JSX_FRAGMENT
}
func StaticGet_Annotation_LEGACY_NAMED_TS_EXPORT() Annotation {
	return enum_Annotation_LEGACY_NAMED_TS_EXPORT
}
func StaticGet_Annotation_LENDS() Annotation {
	return enum_Annotation_LENDS
}
func StaticGet_Annotation_LICENSE() Annotation {
	return enum_Annotation_LICENSE
}
func StaticGet_Annotation_LOG_TYPE_IN_COMPILER() Annotation {
	return enum_Annotation_LOG_TYPE_IN_COMPILER
}
func StaticGet_Annotation_MAY_HAVE_EXTRA_EDGE() Annotation {
	return enum_Annotation_MAY_HAVE_EXTRA_EDGE
}
func StaticGet_Annotation_MEANING() Annotation {
	return enum_Annotation_MEANING
}
func StaticGet_Annotation_MIXIN_CLASS() Annotation {
	return enum_Annotation_MIXIN_CLASS
}
func StaticGet_Annotation_MIXIN_FUNCTION() Annotation {
	return enum_Annotation_MIXIN_FUNCTION
}
func StaticGet_Annotation_MODIFIES() Annotation {
	return enum_Annotation_MODIFIES
}
func StaticGet_Annotation_MODS() Annotation {
	return enum_Annotation_MODS
}
func StaticGet_Annotation_NG_INJECT() Annotation {
	return enum_Annotation_NG_INJECT
}
func StaticGet_Annotation_NO_COLLAPSE() Annotation {
	return enum_Annotation_NO_COLLAPSE
}
func StaticGet_Annotation_NO_COMPILE() Annotation {
	return enum_Annotation_NO_COMPILE
}
func StaticGet_Annotation_NO_COVERAGE() Annotation {
	return enum_Annotation_NO_COVERAGE
}
func StaticGet_Annotation_NO_DTS() Annotation {
	return enum_Annotation_NO_DTS
}
func StaticGet_Annotation_NO_INLINE() Annotation {
	return enum_Annotation_NO_INLINE
}
func StaticGet_Annotation_NO_SIDE_EFFECTS() Annotation {
	return enum_Annotation_NO_SIDE_EFFECTS
}
func StaticGet_Annotation_NOT_IMPLEMENTED() Annotation {
	return enum_Annotation_NOT_IMPLEMENTED
}
func StaticGet_Annotation_OVERRIDE() Annotation {
	return enum_Annotation_OVERRIDE
}
func StaticGet_Annotation_PACKAGE() Annotation {
	return enum_Annotation_PACKAGE
}
func StaticGet_Annotation_PARAM() Annotation {
	return enum_Annotation_PARAM
}
func StaticGet_Annotation_POLYMER() Annotation {
	return enum_Annotation_POLYMER
}
func StaticGet_Annotation_POLYMER_BEHAVIOR() Annotation {
	return enum_Annotation_POLYMER_BEHAVIOR
}
func StaticGet_Annotation_PRESERVE() Annotation {
	return enum_Annotation_PRESERVE
}
func StaticGet_Annotation_PRIVATE() Annotation {
	return enum_Annotation_PRIVATE
}
func StaticGet_Annotation_PROTECTED() Annotation {
	return enum_Annotation_PROTECTED
}
func StaticGet_Annotation_PROVIDE_GOOG() Annotation {
	return enum_Annotation_PROVIDE_GOOG
}
func StaticGet_Annotation_PROVIDE_ALREADY_PROVIDED() Annotation {
	return enum_Annotation_PROVIDE_ALREADY_PROVIDED
}
func StaticGet_Annotation_PUBLIC() Annotation {
	return enum_Annotation_PUBLIC
}
func StaticGet_Annotation_PURE_OR_BREAK_MY_CODE() Annotation {
	return enum_Annotation_PURE_OR_BREAK_MY_CODE
}
func StaticGet_Annotation_REQUIRE_INLINING() Annotation {
	return enum_Annotation_REQUIRE_INLINING
}
func StaticGet_Annotation_RETURN() Annotation {
	return enum_Annotation_RETURN
}
func StaticGet_Annotation_SASS_GENERATED_CSS_TS() Annotation {
	return enum_Annotation_SASS_GENERATED_CSS_TS
}
func StaticGet_Annotation_SEE() Annotation {
	return enum_Annotation_SEE
}
func StaticGet_Annotation_SOY_MODULE() Annotation {
	return enum_Annotation_SOY_MODULE
}
func StaticGet_Annotation_SOY_TEMPLATE() Annotation {
	return enum_Annotation_SOY_TEMPLATE
}
func StaticGet_Annotation_STRUCT() Annotation {
	return enum_Annotation_STRUCT
}
func StaticGet_Annotation_SUPPRESS() Annotation {
	return enum_Annotation_SUPPRESS
}
func StaticGet_Annotation_TEMPLATE() Annotation {
	return enum_Annotation_TEMPLATE
}
func StaticGet_Annotation_CLOSURE_UNAWARE_CODE() Annotation {
	return enum_Annotation_CLOSURE_UNAWARE_CODE
}
func StaticGet_Annotation_THIS() Annotation {
	return enum_Annotation_THIS
}
func StaticGet_Annotation_THROWS() Annotation {
	return enum_Annotation_THROWS
}
func StaticGet_Annotation_TYPE() Annotation {
	return enum_Annotation_TYPE
}
func StaticGet_Annotation_TYPEDEF() Annotation {
	return enum_Annotation_TYPEDEF
}
func StaticGet_Annotation_TYPE_SUMMARY() Annotation {
	return enum_Annotation_TYPE_SUMMARY
}
func StaticGet_Annotation_UNRESTRICTED() Annotation {
	return enum_Annotation_UNRESTRICTED
}
func StaticGet_Annotation_USED_VIA_DOT_CONSTRUCTOR() Annotation {
	return enum_Annotation_USED_VIA_DOT_CONSTRUCTOR
}
func StaticGet_Annotation_WIZACTION() Annotation {
	return enum_Annotation_WIZACTION
}
func StaticGet_Annotation_TS_TYPE() Annotation {
	return enum_Annotation_TS_TYPE
}
func StaticGet_Annotation_WIZ_ANALYZER() Annotation {
	return enum_Annotation_WIZ_ANALYZER
}
func StaticGet_Annotation_WIZCALLBACK() Annotation {
	return enum_Annotation_WIZCALLBACK
}

var static_Annotation_recognizedAnnotations j2gort.Ext_com_google_common_collect_ImmutableMap

func StaticGet_Annotation_RecognizedAnnotations() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_Annotation()
	return static_Annotation_recognizedAnnotations
}
func StaticSet_Annotation_RecognizedAnnotations(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_Annotation()
	static_Annotation_recognizedAnnotations = v
}

var clinitDone_Annotation uint32
var clinitState_Annotation j2gort.ClinitState

func clinit_Annotation() {
	if atomic.LoadUint32(&clinitDone_Annotation) != 0 {
		return
	}
	clinitSlow_Annotation()
}

func clinitSlow_Annotation() {
	if !j2gort.EnterClinit(&clinitState_Annotation, &clinitDone_Annotation) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Annotation, &clinitDone_Annotation)
	static_Annotation_recognizedAnnotations = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.NewExt_com_google_common_collect_ImmutableMap_Builder().Put_2("ngInject", j2gort.BoxNullableVal(StaticGet_Annotation_NG_INJECT()))).Put_2("abstract", j2gort.BoxNullableVal(StaticGet_Annotation_ABSTRACT()))).Put_2("alternateMessageId", j2gort.BoxNullableVal(StaticGet_Annotation_ALTERNATE_MESSAGE_ID()))).Put_2("argument", j2gort.BoxNullableVal(StaticGet_Annotation_PARAM()))).Put_2("author", j2gort.BoxNullableVal(StaticGet_Annotation_AUTHOR()))).Put_2("closurePrimitive", j2gort.BoxNullableVal(StaticGet_Annotation_CLOSURE_PRIMITIVE()))).Put_2("closureUnaware", j2gort.BoxNullableVal(StaticGet_Annotation_CLOSURE_UNAWARE_CODE()))).Put_2("const", j2gort.BoxNullableVal(StaticGet_Annotation_CONSTANT()))).Put_2("collapsibleOrBreakMyCode", j2gort.BoxNullableVal(StaticGet_Annotation_COLLAPSIBLE_OR_BREAK_MY_CODE()))).Put_2("constant", j2gort.BoxNullableVal(StaticGet_Annotation_CONSTANT()))).Put_2("constructor", j2gort.BoxNullableVal(StaticGet_Annotation_CONSTRUCTOR()))).Put_2("customElement", j2gort.BoxNullableVal(StaticGet_Annotation_CUSTOM_ELEMENT()))).Put_2("copyright", j2gort.BoxNullableVal(StaticGet_Annotation_LICENSE()))).Put_2("define", j2gort.BoxNullableVal(StaticGet_Annotation_DEFINE()))).Put_2("deprecated", j2gort.BoxNullableVal(StaticGet_Annotation_DEPRECATED()))).Put_2("desc", j2gort.BoxNullableVal(StaticGet_Annotation_DESC()))).Put_2("dict", j2gort.BoxNullableVal(StaticGet_Annotation_DICT()))).Put_2("requireInlining", j2gort.BoxNullableVal(StaticGet_Annotation_REQUIRE_INLINING()))).Put_2("encourageInlining", j2gort.BoxNullableVal(StaticGet_Annotation_ENCOURAGE_INLINING()))).Put_2("enum", j2gort.BoxNullableVal(StaticGet_Annotation_ENUM()))).Put_2("enhance", j2gort.BoxNullableVal(StaticGet_Annotation_ENHANCE()))).Put_2("export", j2gort.BoxNullableVal(StaticGet_Annotation_EXPORT()))).Put_2("extends", j2gort.BoxNullableVal(StaticGet_Annotation_EXTENDS()))).Put_2("externs", j2gort.BoxNullableVal(StaticGet_Annotation_EXTERNS()))).Put_2("fileoverview", j2gort.BoxNullableVal(StaticGet_Annotation_FILE_OVERVIEW()))).Put_2("final", j2gort.BoxNullableVal(StaticGet_Annotation_FINAL()))).Put_2("idGenerator", j2gort.BoxNullableVal(StaticGet_Annotation_IDGENERATOR()))).Put_2("implements", j2gort.BoxNullableVal(StaticGet_Annotation_IMPLEMENTS()))).Put_2("implicitCast", j2gort.BoxNullableVal(StaticGet_Annotation_IMPLICIT_CAST()))).Put_2("inheritDoc", j2gort.BoxNullableVal(StaticGet_Annotation_INHERIT_DOC()))).Put_2("interface", j2gort.BoxNullableVal(StaticGet_Annotation_INTERFACE()))).Put_2("record", j2gort.BoxNullableVal(StaticGet_Annotation_RECORD()))).Put_2("legacyNamedTsExport", j2gort.BoxNullableVal(StaticGet_Annotation_LEGACY_NAMED_TS_EXPORT()))).Put_2("lends", j2gort.BoxNullableVal(StaticGet_Annotation_LENDS()))).Put_2("license", j2gort.BoxNullableVal(StaticGet_Annotation_LICENSE()))).Put_2("logTypeInCompiler", j2gort.BoxNullableVal(StaticGet_Annotation_LOG_TYPE_IN_COMPILER()))).Put_2("mayhaveextraedge", j2gort.BoxNullableVal(StaticGet_Annotation_MAY_HAVE_EXTRA_EDGE()))).Put_2("meaning", j2gort.BoxNullableVal(StaticGet_Annotation_MEANING()))).Put_2("mixinClass", j2gort.BoxNullableVal(StaticGet_Annotation_MIXIN_CLASS()))).Put_2("mixinFunction", j2gort.BoxNullableVal(StaticGet_Annotation_MIXIN_FUNCTION()))).Put_2("modifies", j2gort.BoxNullableVal(StaticGet_Annotation_MODIFIES()))).Put_2("mods", j2gort.BoxNullableVal(StaticGet_Annotation_MODS()))).Put_2("nocollapse", j2gort.BoxNullableVal(StaticGet_Annotation_NO_COLLAPSE()))).Put_2("nocompile", j2gort.BoxNullableVal(StaticGet_Annotation_NO_COMPILE()))).Put_2("nocoverage", j2gort.BoxNullableVal(StaticGet_Annotation_NO_COVERAGE()))).Put_2("nodts", j2gort.BoxNullableVal(StaticGet_Annotation_NO_DTS()))).Put_2("noinline", j2gort.BoxNullableVal(StaticGet_Annotation_NO_INLINE()))).Put_2("nosideeffects", j2gort.BoxNullableVal(StaticGet_Annotation_NO_SIDE_EFFECTS()))).Put_2("override", j2gort.BoxNullableVal(StaticGet_Annotation_OVERRIDE()))).Put_2("owner", j2gort.BoxNullableVal(StaticGet_Annotation_AUTHOR()))).Put_2("package", j2gort.BoxNullableVal(StaticGet_Annotation_PACKAGE()))).Put_2("param", j2gort.BoxNullableVal(StaticGet_Annotation_PARAM()))).Put_2("polymer", j2gort.BoxNullableVal(StaticGet_Annotation_POLYMER()))).Put_2("polymerBehavior", j2gort.BoxNullableVal(StaticGet_Annotation_POLYMER_BEHAVIOR()))).Put_2("preserve", j2gort.BoxNullableVal(StaticGet_Annotation_PRESERVE()))).Put_2("private", j2gort.BoxNullableVal(StaticGet_Annotation_PRIVATE()))).Put_2("protected", j2gort.BoxNullableVal(StaticGet_Annotation_PROTECTED()))).Put_2("provideGoog", j2gort.BoxNullableVal(StaticGet_Annotation_PROVIDE_GOOG()))).Put_2("provideAlreadyProvided", j2gort.BoxNullableVal(StaticGet_Annotation_PROVIDE_ALREADY_PROVIDED()))).Put_2("public", j2gort.BoxNullableVal(StaticGet_Annotation_PUBLIC()))).Put_2("pureOrBreakMyCode", j2gort.BoxNullableVal(StaticGet_Annotation_PURE_OR_BREAK_MY_CODE()))).Put_2("return", j2gort.BoxNullableVal(StaticGet_Annotation_RETURN()))).Put_2("returns", j2gort.BoxNullableVal(StaticGet_Annotation_RETURN()))).Put_2("sassGeneratedCssTs", j2gort.BoxNullableVal(StaticGet_Annotation_SASS_GENERATED_CSS_TS()))).Put_2("see", j2gort.BoxNullableVal(StaticGet_Annotation_SEE()))).Put_2("soyModule", j2gort.BoxNullableVal(StaticGet_Annotation_SOY_MODULE()))).Put_2("soyTemplate", j2gort.BoxNullableVal(StaticGet_Annotation_SOY_TEMPLATE()))).Put_2("struct", j2gort.BoxNullableVal(StaticGet_Annotation_STRUCT()))).Put_2("suppress", j2gort.BoxNullableVal(StaticGet_Annotation_SUPPRESS()))).Put_2("template", j2gort.BoxNullableVal(StaticGet_Annotation_TEMPLATE()))).Put_2("this", j2gort.BoxNullableVal(StaticGet_Annotation_THIS()))).Put_2("throws", j2gort.BoxNullableVal(StaticGet_Annotation_THROWS()))).Put_2("type", j2gort.BoxNullableVal(StaticGet_Annotation_TYPE()))).Put_2("typedef", j2gort.BoxNullableVal(StaticGet_Annotation_TYPEDEF()))).Put_2("typeSummary", j2gort.BoxNullableVal(StaticGet_Annotation_TYPE_SUMMARY()))).Put_2("unrestricted", j2gort.BoxNullableVal(StaticGet_Annotation_UNRESTRICTED()))).Put_2("usedViaDotConstructor", j2gort.BoxNullableVal(StaticGet_Annotation_USED_VIA_DOT_CONSTRUCTOR()))).Put_2("wizaction", j2gort.BoxNullableVal(StaticGet_Annotation_WIZACTION()))).Put_2("tsType", j2gort.BoxNullableVal(StaticGet_Annotation_TS_TYPE()))).Put_2("wizAnalyzer", j2gort.BoxNullableVal(StaticGet_Annotation_WIZ_ANALYZER()))).Put_2("wizcallback", j2gort.BoxNullableVal(StaticGet_Annotation_WIZCALLBACK()))).BuildOrThrow())
}

func Annotation_Values() []Annotation {
	clinit_Annotation()
	return j2gort.CloneSlice(enumValues_Annotation)
}

func Annotation_ValueOf(name string) Annotation {
	clinit_Annotation()
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Annotation {
		if enumBaseTable_Annotation[i_].Name_ == name {
			return Annotation(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.Annotation." + name))
}

func (this_ Annotation) IsJavaNil_() bool                              { return this_ == 0 }
func (Annotation) HasIdentityEquals_()                                 {}
func (Annotation) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Annotation) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Annotation) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Annotation) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.Annotation", "Annotation"
}
func (this_ Annotation) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Annotation[this_-1]
}
func (this_ Annotation) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Annotation) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Annotation) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Annotation[this_-1].Name_
}
func (this_ Annotation) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Annotation[this_-1].Name_
}
func (this_ Annotation) Equals(other_ any) bool {
	o_, ok_ := other_.(Annotation)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Annotation) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Annotation[this_-1].HashCode()
}
func (this_ Annotation) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Annotation[this_-1].CompareTo(other_)
}
func (this_ Annotation) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Annotation[this_-1].Name_)
}
func init() {
	vals_ := []any{Annotation(1), Annotation(2), Annotation(3), Annotation(4), Annotation(5), Annotation(6), Annotation(7), Annotation(8), Annotation(9), Annotation(10), Annotation(11), Annotation(12), Annotation(13), Annotation(14), Annotation(15), Annotation(16), Annotation(17), Annotation(18), Annotation(19), Annotation(20), Annotation(21), Annotation(22), Annotation(23), Annotation(24), Annotation(25), Annotation(26), Annotation(27), Annotation(28), Annotation(29), Annotation(30), Annotation(31), Annotation(32), Annotation(33), Annotation(34), Annotation(35), Annotation(36), Annotation(37), Annotation(38), Annotation(39), Annotation(40), Annotation(41), Annotation(42), Annotation(43), Annotation(44), Annotation(45), Annotation(46), Annotation(47), Annotation(48), Annotation(49), Annotation(50), Annotation(51), Annotation(52), Annotation(53), Annotation(54), Annotation(55), Annotation(56), Annotation(57), Annotation(58), Annotation(59), Annotation(60), Annotation(61), Annotation(62), Annotation(63), Annotation(64), Annotation(65), Annotation(66), Annotation(67), Annotation(68), Annotation(69), Annotation(70), Annotation(71), Annotation(72), Annotation(73), Annotation(74), Annotation(75), Annotation(76), Annotation(77), Annotation(78), Annotation(79)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.Annotation", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.Annotation", "com.google.javascript.jscomp.parsing.Annotation", func(v_ any) bool {
		return j2gort.InstanceOf[Annotation](v_)
	})
}
