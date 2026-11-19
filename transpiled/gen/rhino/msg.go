/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/Msg.java (generated mechanically by j2go).

package rhino

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.rhino.Msg
type Msg uint8

const (
	enum_Msg_BAD_FILEOVERVIEW_VISIBIIITY_ANNOTATION              Msg = 1
	enum_Msg_BAD_JSDOC_TAG                                       Msg = 2
	enum_Msg_DUP_VARIABLE_NAME                                   Msg = 3
	enum_Msg_END_ANNOTATION_EXPECTED                             Msg = 4
	enum_Msg_INVALID_VARIABLE_NAME                               Msg = 5
	enum_Msg_JSDOC_ALTERNATEMESSAGEID_EXTRA                      Msg = 6
	enum_Msg_JSDOC_REQUIRE_INLINING                              Msg = 7
	enum_Msg_JSDOC_ENCOURAGE_INLINING                            Msg = 8
	enum_Msg_JSDOC_AUTHORMISSING                                 Msg = 9
	enum_Msg_JSDOC_CLOSUREPRIMITIVE_EXTRA                        Msg = 10
	enum_Msg_JSDOC_CLOSUREPRIMITIVE_INVALID                      Msg = 11
	enum_Msg_JSDOC_CLOSUREPRIMITIVE_MISSING                      Msg = 12
	enum_Msg_JSDOC_CLOSURE_UNAWARE_CODE_EXTRA                    Msg = 13
	enum_Msg_JSDOC_CLOSURE_UNAWARE_CONFIG_INVALID_VALUE          Msg = 14
	enum_Msg_JSDOC_CLOSURE_UNAWARE_CODE_INVALID                  Msg = 15
	enum_Msg_JSDOC_COLLAPSIBLEORBREAKMYCODE                      Msg = 16
	enum_Msg_JSDOC_CONST                                         Msg = 17
	enum_Msg_JSDOC_CUSTOMELEMENT_EXTRA                           Msg = 18
	enum_Msg_JSDOC_DEFINE                                        Msg = 19
	enum_Msg_JSDOC_DEPRECATED                                    Msg = 20
	enum_Msg_JSDOC_DESC_EXTRA                                    Msg = 21
	enum_Msg_JSDOC_EXPORT                                        Msg = 22
	enum_Msg_JSDOC_EXTENDS_DUPLICATE                             Msg = 23
	enum_Msg_JSDOC_EXTERNS                                       Msg = 24
	enum_Msg_JSDOC_EXTRAVERSION                                  Msg = 25
	enum_Msg_JSDOC_EXTRA_VISIBILITY                              Msg = 26
	enum_Msg_JSDOC_FILEOVERVIEW_EXTRA                            Msg = 27
	enum_Msg_JSDOC_FINAL                                         Msg = 28
	enum_Msg_JSDOC_FUNCTION_NEWNOTOBJECT                         Msg = 29
	enum_Msg_JSDOC_FUNCTION_VARARGS                              Msg = 30
	enum_Msg_JSDOC_IDGEN_BAD                                     Msg = 31
	enum_Msg_JSDOC_IDGEN_DUPLICATE                               Msg = 32
	enum_Msg_JSDOC_IDGEN_UNKNOWN                                 Msg = 33
	enum_Msg_JSDOC_IMPLEMENTS_DUPLICATE                          Msg = 34
	enum_Msg_JSDOC_IMPLEMENTS_EXTRAQUALIFIER                     Msg = 35
	enum_Msg_JSDOC_IMPLICITCAST                                  Msg = 36
	enum_Msg_JSDOC_IMPORT                                        Msg = 37
	enum_Msg_JSDOC_INCOMPAT_TYPE                                 Msg = 38
	enum_Msg_JSDOC_INCOMPAT_INLINING                             Msg = 39
	enum_Msg_JSDOC_INTERFACE_CONSTRUCTOR                         Msg = 40
	enum_Msg_JSDOC_LEGACY_TS_EXPORT                              Msg = 41
	enum_Msg_JSDOC_LENDS_INCOMPATIBLE                            Msg = 42
	enum_Msg_JSDOC_LENDS_MISSING                                 Msg = 43
	enum_Msg_JSDOC_LOCALEFILE                                    Msg = 44
	enum_Msg_JSDOC_LOCALEOBJECT                                  Msg = 45
	enum_Msg_JSDOC_LOCALESELECT                                  Msg = 46
	enum_Msg_JSDOC_LOCALEVALUE                                   Msg = 47
	enum_Msg_JSDOC_PROVIDE_GOOG                                  Msg = 48
	enum_Msg_JSDOC_PROVIDE_ALREADY_PROVIDED                      Msg = 49
	enum_Msg_JSDOC_MEANING_EXTRA                                 Msg = 50
	enum_Msg_JSDOC_MISSING_BRACES                                Msg = 51
	enum_Msg_JSDOC_MISSING_COLON                                 Msg = 52
	enum_Msg_JSDOC_MISSING_GT                                    Msg = 53
	enum_Msg_JSDOC_MISSING_LC                                    Msg = 54
	enum_Msg_JSDOC_MISSING_LP                                    Msg = 55
	enum_Msg_JSDOC_MISSING_RB                                    Msg = 56
	enum_Msg_JSDOC_MISSING_RC                                    Msg = 57
	enum_Msg_JSDOC_MISSING_RP                                    Msg = 58
	enum_Msg_JSDOC_MISSING_TYPE_DECLARATION                      Msg = 59
	enum_Msg_JSDOC_MIXINCLASS_EXTRA                              Msg = 60
	enum_Msg_JSDOC_MIXINFUNCTION_EXTRA                           Msg = 61
	enum_Msg_JSDOC_MODIFIES                                      Msg = 62
	enum_Msg_JSDOC_MODIFIES_DUPLICATE                            Msg = 63
	enum_Msg_JSDOC_MODIFIES_UNKNOWN                              Msg = 64
	enum_Msg_JSDOC_MODS                                          Msg = 65
	enum_Msg_JSDOC_MODS_EXTRA                                    Msg = 66
	enum_Msg_JSDOC_NAME_SYNTAX                                   Msg = 67
	enum_Msg_JSDOC_NGINJECT_EXTRA                                Msg = 68
	enum_Msg_JSDOC_NOCOLLAPSE                                    Msg = 69
	enum_Msg_JSDOC_NOCOMPILE                                     Msg = 70
	enum_Msg_JSDOC_NOCOVERAGE                                    Msg = 71
	enum_Msg_JSDOC_NODTS                                         Msg = 72
	enum_Msg_JSDOC_NOINLINE                                      Msg = 73
	enum_Msg_JSDOC_NOSIDEEFFECTS                                 Msg = 74
	enum_Msg_JSDOC_NOSIDEEFFECTS_WITH_MODIFIES                   Msg = 75
	enum_Msg_JSDOC_NOSIDEEFFECTS_WITH_THROWS                     Msg = 76
	enum_Msg_JSDOC_OVERRIDE                                      Msg = 77
	enum_Msg_JSDOC_POLYMERBEHAVIOR_EXTRA                         Msg = 78
	enum_Msg_JSDOC_POLYMER_EXTRA                                 Msg = 79
	enum_Msg_JSDOC_PUREORBREAKMYCODE                             Msg = 80
	enum_Msg_JSDOC_RECORD                                        Msg = 81
	enum_Msg_JSDOC_SASS_GENERATED_CSS_TS                         Msg = 82
	enum_Msg_JSDOC_SEEMISSING                                    Msg = 83
	enum_Msg_JSDOC_SUPPRESS                                      Msg = 84
	enum_Msg_JSDOC_SUPPRESS_UNKNOWN                              Msg = 85
	enum_Msg_JSDOC_TEMPLATE_BOUNDEDGENERICS_USED                 Msg = 86
	enum_Msg_JSDOC_TEMPLATE_BOUNDSWITHTTL                        Msg = 87
	enum_Msg_JSDOC_TEMPLATE_MULTIPLEDECLARATION                  Msg = 88
	enum_Msg_JSDOC_TEMPLATE_NAME_MISSING                         Msg = 89
	enum_Msg_JSDOC_TEMPLATE_NAME_REDECLARATION                   Msg = 90
	enum_Msg_JSDOC_TEMPLATE_TYPETRANSFORMATION_EXPRESSIONMISSING Msg = 91
	enum_Msg_JSDOC_TEMPLATE_TYPETRANSFORMATION_MISSINGDELIMIIER  Msg = 92
	enum_Msg_JSDOC_THROWS                                        Msg = 93
	enum_Msg_JSDOC_TYPE                                          Msg = 94
	enum_Msg_JSDOC_TYPESUMMARY                                   Msg = 95
	enum_Msg_JSDOC_TYPETRANSFORMATION_EXTRA_PARAM                Msg = 96
	enum_Msg_JSDOC_TYPETRANSFORMATION_INVALID                    Msg = 97
	enum_Msg_JSDOC_TYPETRANSFORMATION_INVALID_EXPRESSION         Msg = 98
	enum_Msg_JSDOC_TYPETRANSFORMATION_INVALID_INSIDE             Msg = 99
	enum_Msg_JSDOC_TYPETRANSFORMATION_MISSING_PARAM              Msg = 100
	enum_Msg_JSDOC_TYPE_RECORD_DUPLICATE                         Msg = 101
	enum_Msg_JSDOC_TYPE_SYNTAX                                   Msg = 102
	enum_Msg_JSDOC_USEDVIADOTCONSTRUCTOR                         Msg = 103
	enum_Msg_JSDOC_UNNECESSARY_BRACES                            Msg = 104
	enum_Msg_JSDOC_VERSIONMISSING                                Msg = 105
	enum_Msg_JSDOC_WIZACTION                                     Msg = 106
	enum_Msg_MISSING_VARIABLE_NAME                               Msg = 107
	enum_Msg_NO_TYPE_NAME                                        Msg = 108
	enum_Msg_UNEXPECTED_EOF                                      Msg = 109
	enum_Msg_JSDOC_WIZCALLBACK                                   Msg = 110
)

var enumBaseTable_Msg = [...]j2gort.EnumBase{
	{Name_: "BAD_FILEOVERVIEW_VISIBIIITY_ANNOTATION", Ordinal_: 0},
	{Name_: "BAD_JSDOC_TAG", Ordinal_: 1},
	{Name_: "DUP_VARIABLE_NAME", Ordinal_: 2},
	{Name_: "END_ANNOTATION_EXPECTED", Ordinal_: 3},
	{Name_: "INVALID_VARIABLE_NAME", Ordinal_: 4},
	{Name_: "JSDOC_ALTERNATEMESSAGEID_EXTRA", Ordinal_: 5},
	{Name_: "JSDOC_REQUIRE_INLINING", Ordinal_: 6},
	{Name_: "JSDOC_ENCOURAGE_INLINING", Ordinal_: 7},
	{Name_: "JSDOC_AUTHORMISSING", Ordinal_: 8},
	{Name_: "JSDOC_CLOSUREPRIMITIVE_EXTRA", Ordinal_: 9},
	{Name_: "JSDOC_CLOSUREPRIMITIVE_INVALID", Ordinal_: 10},
	{Name_: "JSDOC_CLOSUREPRIMITIVE_MISSING", Ordinal_: 11},
	{Name_: "JSDOC_CLOSURE_UNAWARE_CODE_EXTRA", Ordinal_: 12},
	{Name_: "JSDOC_CLOSURE_UNAWARE_CONFIG_INVALID_VALUE", Ordinal_: 13},
	{Name_: "JSDOC_CLOSURE_UNAWARE_CODE_INVALID", Ordinal_: 14},
	{Name_: "JSDOC_COLLAPSIBLEORBREAKMYCODE", Ordinal_: 15},
	{Name_: "JSDOC_CONST", Ordinal_: 16},
	{Name_: "JSDOC_CUSTOMELEMENT_EXTRA", Ordinal_: 17},
	{Name_: "JSDOC_DEFINE", Ordinal_: 18},
	{Name_: "JSDOC_DEPRECATED", Ordinal_: 19},
	{Name_: "JSDOC_DESC_EXTRA", Ordinal_: 20},
	{Name_: "JSDOC_EXPORT", Ordinal_: 21},
	{Name_: "JSDOC_EXTENDS_DUPLICATE", Ordinal_: 22},
	{Name_: "JSDOC_EXTERNS", Ordinal_: 23},
	{Name_: "JSDOC_EXTRAVERSION", Ordinal_: 24},
	{Name_: "JSDOC_EXTRA_VISIBILITY", Ordinal_: 25},
	{Name_: "JSDOC_FILEOVERVIEW_EXTRA", Ordinal_: 26},
	{Name_: "JSDOC_FINAL", Ordinal_: 27},
	{Name_: "JSDOC_FUNCTION_NEWNOTOBJECT", Ordinal_: 28},
	{Name_: "JSDOC_FUNCTION_VARARGS", Ordinal_: 29},
	{Name_: "JSDOC_IDGEN_BAD", Ordinal_: 30},
	{Name_: "JSDOC_IDGEN_DUPLICATE", Ordinal_: 31},
	{Name_: "JSDOC_IDGEN_UNKNOWN", Ordinal_: 32},
	{Name_: "JSDOC_IMPLEMENTS_DUPLICATE", Ordinal_: 33},
	{Name_: "JSDOC_IMPLEMENTS_EXTRAQUALIFIER", Ordinal_: 34},
	{Name_: "JSDOC_IMPLICITCAST", Ordinal_: 35},
	{Name_: "JSDOC_IMPORT", Ordinal_: 36},
	{Name_: "JSDOC_INCOMPAT_TYPE", Ordinal_: 37},
	{Name_: "JSDOC_INCOMPAT_INLINING", Ordinal_: 38},
	{Name_: "JSDOC_INTERFACE_CONSTRUCTOR", Ordinal_: 39},
	{Name_: "JSDOC_LEGACY_TS_EXPORT", Ordinal_: 40},
	{Name_: "JSDOC_LENDS_INCOMPATIBLE", Ordinal_: 41},
	{Name_: "JSDOC_LENDS_MISSING", Ordinal_: 42},
	{Name_: "JSDOC_LOCALEFILE", Ordinal_: 43},
	{Name_: "JSDOC_LOCALEOBJECT", Ordinal_: 44},
	{Name_: "JSDOC_LOCALESELECT", Ordinal_: 45},
	{Name_: "JSDOC_LOCALEVALUE", Ordinal_: 46},
	{Name_: "JSDOC_PROVIDE_GOOG", Ordinal_: 47},
	{Name_: "JSDOC_PROVIDE_ALREADY_PROVIDED", Ordinal_: 48},
	{Name_: "JSDOC_MEANING_EXTRA", Ordinal_: 49},
	{Name_: "JSDOC_MISSING_BRACES", Ordinal_: 50},
	{Name_: "JSDOC_MISSING_COLON", Ordinal_: 51},
	{Name_: "JSDOC_MISSING_GT", Ordinal_: 52},
	{Name_: "JSDOC_MISSING_LC", Ordinal_: 53},
	{Name_: "JSDOC_MISSING_LP", Ordinal_: 54},
	{Name_: "JSDOC_MISSING_RB", Ordinal_: 55},
	{Name_: "JSDOC_MISSING_RC", Ordinal_: 56},
	{Name_: "JSDOC_MISSING_RP", Ordinal_: 57},
	{Name_: "JSDOC_MISSING_TYPE_DECLARATION", Ordinal_: 58},
	{Name_: "JSDOC_MIXINCLASS_EXTRA", Ordinal_: 59},
	{Name_: "JSDOC_MIXINFUNCTION_EXTRA", Ordinal_: 60},
	{Name_: "JSDOC_MODIFIES", Ordinal_: 61},
	{Name_: "JSDOC_MODIFIES_DUPLICATE", Ordinal_: 62},
	{Name_: "JSDOC_MODIFIES_UNKNOWN", Ordinal_: 63},
	{Name_: "JSDOC_MODS", Ordinal_: 64},
	{Name_: "JSDOC_MODS_EXTRA", Ordinal_: 65},
	{Name_: "JSDOC_NAME_SYNTAX", Ordinal_: 66},
	{Name_: "JSDOC_NGINJECT_EXTRA", Ordinal_: 67},
	{Name_: "JSDOC_NOCOLLAPSE", Ordinal_: 68},
	{Name_: "JSDOC_NOCOMPILE", Ordinal_: 69},
	{Name_: "JSDOC_NOCOVERAGE", Ordinal_: 70},
	{Name_: "JSDOC_NODTS", Ordinal_: 71},
	{Name_: "JSDOC_NOINLINE", Ordinal_: 72},
	{Name_: "JSDOC_NOSIDEEFFECTS", Ordinal_: 73},
	{Name_: "JSDOC_NOSIDEEFFECTS_WITH_MODIFIES", Ordinal_: 74},
	{Name_: "JSDOC_NOSIDEEFFECTS_WITH_THROWS", Ordinal_: 75},
	{Name_: "JSDOC_OVERRIDE", Ordinal_: 76},
	{Name_: "JSDOC_POLYMERBEHAVIOR_EXTRA", Ordinal_: 77},
	{Name_: "JSDOC_POLYMER_EXTRA", Ordinal_: 78},
	{Name_: "JSDOC_PUREORBREAKMYCODE", Ordinal_: 79},
	{Name_: "JSDOC_RECORD", Ordinal_: 80},
	{Name_: "JSDOC_SASS_GENERATED_CSS_TS", Ordinal_: 81},
	{Name_: "JSDOC_SEEMISSING", Ordinal_: 82},
	{Name_: "JSDOC_SUPPRESS", Ordinal_: 83},
	{Name_: "JSDOC_SUPPRESS_UNKNOWN", Ordinal_: 84},
	{Name_: "JSDOC_TEMPLATE_BOUNDEDGENERICS_USED", Ordinal_: 85},
	{Name_: "JSDOC_TEMPLATE_BOUNDSWITHTTL", Ordinal_: 86},
	{Name_: "JSDOC_TEMPLATE_MULTIPLEDECLARATION", Ordinal_: 87},
	{Name_: "JSDOC_TEMPLATE_NAME_MISSING", Ordinal_: 88},
	{Name_: "JSDOC_TEMPLATE_NAME_REDECLARATION", Ordinal_: 89},
	{Name_: "JSDOC_TEMPLATE_TYPETRANSFORMATION_EXPRESSIONMISSING", Ordinal_: 90},
	{Name_: "JSDOC_TEMPLATE_TYPETRANSFORMATION_MISSINGDELIMIIER", Ordinal_: 91},
	{Name_: "JSDOC_THROWS", Ordinal_: 92},
	{Name_: "JSDOC_TYPE", Ordinal_: 93},
	{Name_: "JSDOC_TYPESUMMARY", Ordinal_: 94},
	{Name_: "JSDOC_TYPETRANSFORMATION_EXTRA_PARAM", Ordinal_: 95},
	{Name_: "JSDOC_TYPETRANSFORMATION_INVALID", Ordinal_: 96},
	{Name_: "JSDOC_TYPETRANSFORMATION_INVALID_EXPRESSION", Ordinal_: 97},
	{Name_: "JSDOC_TYPETRANSFORMATION_INVALID_INSIDE", Ordinal_: 98},
	{Name_: "JSDOC_TYPETRANSFORMATION_MISSING_PARAM", Ordinal_: 99},
	{Name_: "JSDOC_TYPE_RECORD_DUPLICATE", Ordinal_: 100},
	{Name_: "JSDOC_TYPE_SYNTAX", Ordinal_: 101},
	{Name_: "JSDOC_USEDVIADOTCONSTRUCTOR", Ordinal_: 102},
	{Name_: "JSDOC_UNNECESSARY_BRACES", Ordinal_: 103},
	{Name_: "JSDOC_VERSIONMISSING", Ordinal_: 104},
	{Name_: "JSDOC_WIZACTION", Ordinal_: 105},
	{Name_: "MISSING_VARIABLE_NAME", Ordinal_: 106},
	{Name_: "NO_TYPE_NAME", Ordinal_: 107},
	{Name_: "UNEXPECTED_EOF", Ordinal_: 108},
	{Name_: "JSDOC_WIZCALLBACK", Ordinal_: 109},
}
var enumValues_Msg = []Msg{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110}

var enumField_Msg_text = [...]string{
	"{0} visibility not allowed in @fileoverview block",
	"illegal use of unknown JSDoc tag \"{0}\"; ignoring it. Place another character before the @ to stop JSCompiler from parsing it as an annotation.",
	"duplicate variable name \"{0}\"",
	"expected end of line or comment.",
	"invalid param name \"{0}\"",
	"extra @alternateMessageId tag",
	"extra @requireInlining tag",
	"extra @encourageInlining tag",
	"@author tag missing author",
	"conflicting @closurePrimitive tag",
	"invalid id in @closurePrimitive tag.",
	"missing id in @closurePrimitive tag.",
	"extra @closureUnaware tag",
	"invalid value for @closureUnaware: {0}",
	"@closureUnaware annotation is not allowed in this compilation",
	"extra @collapsibleOrBreakMyCode tag",
	"conflicting @const tag",
	"extra @customElement tag",
	"conflicting @define tag",
	"extra @deprecated tag",
	"extra @desc tag",
	"extra @export tag",
	"duplicate @extends tag",
	"extra @externs tag",
	"conflicting @version tag",
	"extra visibility tag",
	"extra @fileoverview tag",
	"extra @final tag.",
	"constructed type must be an object type",
	"variable length argument must be last.",
	"malformed @idGenerator tag",
	"extra @idGenerator tag",
	"unknown @idGenerator parameter: {0}",
	"duplicate @implements tag.",
	"@implements/@extends requires a bare interface/record name without ! or ?.",
	"extra @implicitCast tag.",
	"Import in typedef is not supported.",
	"type annotation incompatible with other annotations.",
	"@requireInlining, @encourageInlining, and @noinline are mutually exclusive.",
	"cannot be both an interface and a constructor.",
	"extra @legacyNamedTsExport tag",
	"@lends tag incompatible with other annotations.",
	"missing object name in @lends tag.",
	"extra @localeFile tag",
	"extra @localeObject tag",
	"extra @localeSelect tag",
	"extra @localeValue tag",
	"extra @provideGoog tag",
	"extra @provideAlreadyProvided tag",
	"extra @meaning tag",
	"Type annotations should have curly braces.",
	"expecting colon after this",
	"missing closing >",
	"missing opening {",
	"missing opening (",
	"missing closing ]",
	"expected closing }",
	"missing closing )",
	"Missing type declaration.",
	"extra @mixinClass tag",
	"extra @mixinFunction tag",
	"malformed @modifies tag",
	"conflicting @modifies tag",
	"unknown @modifies parameter: {0}",
	"malformed @mods tag",
	"extra @mods tag",
	"name not recognized due to syntax error.",
	"extra @ngInject tag",
	"extra @nocollapse tag",
	"extra @nocompile tag",
	"extra @nocoverage tag",
	"extra @nodts tag",
	"extra @noinline tag",
	"conflicting @nosideeffects tag",
	"@nosideeffects functions cannot have @modifies (modifying arguments/this is a side effect)",
	"@nosideeffects functions cannot have @throws (throwing is a side effect)",
	"extra @override/@inheritDoc tag.",
	"extra @polymerBehavior tag",
	"extra @polymer tag",
	"extra @pureOrBreakMyCode tag",
	"conflicting @record tag.",
	"extra @sassGeneratedCssTs tag",
	"@see tag missing description",
	"malformed @suppress tag",
	"unknown @suppress parameter: {0}",
	"Bounded generic semantics are currently still in development",
	"Template types cannot combine bounds and TTL.",
	"Multiple template names cannot be declared with bounds or TTL.",
	"@template tag missing type name.",
	"Type name(s) for @template annotation declared twice.",
	"Missing type transformation expression.",
	"Expected end delimiter for a type transformation.",
	"conflicting @throws tag",
	"conflicting @type tag",
	"extra @typeSummary tag",
	"Found extra parameter in {0}",
	"Invalid {0}",
	"Invalid {0} expression",
	"Invalid expression inside {0}",
	"Missing parameter in {0}",
	"Duplicate record field {0}.",
	"type not recognized due to syntax error.",
	"extra @usedViaDotConstructor tag",
	"braces are not required here",
	"@version tag missing version information",
	"extra @wizaction tag",
	"expecting a variable name in a @param tag.",
	"expecting a type name.",
	"Unexpected end of file",
	"extra @wizcallback tag",
}

func (this_ Msg) FieldGet_text() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_Msg_text[this_-1]
}

func StaticGet_Msg_BAD_FILEOVERVIEW_VISIBIIITY_ANNOTATION() Msg {
	return enum_Msg_BAD_FILEOVERVIEW_VISIBIIITY_ANNOTATION
}
func StaticGet_Msg_BAD_JSDOC_TAG() Msg {
	return enum_Msg_BAD_JSDOC_TAG
}
func StaticGet_Msg_DUP_VARIABLE_NAME() Msg {
	return enum_Msg_DUP_VARIABLE_NAME
}
func StaticGet_Msg_END_ANNOTATION_EXPECTED() Msg {
	return enum_Msg_END_ANNOTATION_EXPECTED
}
func StaticGet_Msg_INVALID_VARIABLE_NAME() Msg {
	return enum_Msg_INVALID_VARIABLE_NAME
}
func StaticGet_Msg_JSDOC_ALTERNATEMESSAGEID_EXTRA() Msg {
	return enum_Msg_JSDOC_ALTERNATEMESSAGEID_EXTRA
}
func StaticGet_Msg_JSDOC_REQUIRE_INLINING() Msg {
	return enum_Msg_JSDOC_REQUIRE_INLINING
}
func StaticGet_Msg_JSDOC_ENCOURAGE_INLINING() Msg {
	return enum_Msg_JSDOC_ENCOURAGE_INLINING
}
func StaticGet_Msg_JSDOC_AUTHORMISSING() Msg {
	return enum_Msg_JSDOC_AUTHORMISSING
}
func StaticGet_Msg_JSDOC_CLOSUREPRIMITIVE_EXTRA() Msg {
	return enum_Msg_JSDOC_CLOSUREPRIMITIVE_EXTRA
}
func StaticGet_Msg_JSDOC_CLOSUREPRIMITIVE_INVALID() Msg {
	return enum_Msg_JSDOC_CLOSUREPRIMITIVE_INVALID
}
func StaticGet_Msg_JSDOC_CLOSUREPRIMITIVE_MISSING() Msg {
	return enum_Msg_JSDOC_CLOSUREPRIMITIVE_MISSING
}
func StaticGet_Msg_JSDOC_CLOSURE_UNAWARE_CODE_EXTRA() Msg {
	return enum_Msg_JSDOC_CLOSURE_UNAWARE_CODE_EXTRA
}
func StaticGet_Msg_JSDOC_CLOSURE_UNAWARE_CONFIG_INVALID_VALUE() Msg {
	return enum_Msg_JSDOC_CLOSURE_UNAWARE_CONFIG_INVALID_VALUE
}
func StaticGet_Msg_JSDOC_CLOSURE_UNAWARE_CODE_INVALID() Msg {
	return enum_Msg_JSDOC_CLOSURE_UNAWARE_CODE_INVALID
}
func StaticGet_Msg_JSDOC_COLLAPSIBLEORBREAKMYCODE() Msg {
	return enum_Msg_JSDOC_COLLAPSIBLEORBREAKMYCODE
}
func StaticGet_Msg_JSDOC_CONST() Msg {
	return enum_Msg_JSDOC_CONST
}
func StaticGet_Msg_JSDOC_CUSTOMELEMENT_EXTRA() Msg {
	return enum_Msg_JSDOC_CUSTOMELEMENT_EXTRA
}
func StaticGet_Msg_JSDOC_DEFINE() Msg {
	return enum_Msg_JSDOC_DEFINE
}
func StaticGet_Msg_JSDOC_DEPRECATED() Msg {
	return enum_Msg_JSDOC_DEPRECATED
}
func StaticGet_Msg_JSDOC_DESC_EXTRA() Msg {
	return enum_Msg_JSDOC_DESC_EXTRA
}
func StaticGet_Msg_JSDOC_EXPORT() Msg {
	return enum_Msg_JSDOC_EXPORT
}
func StaticGet_Msg_JSDOC_EXTENDS_DUPLICATE() Msg {
	return enum_Msg_JSDOC_EXTENDS_DUPLICATE
}
func StaticGet_Msg_JSDOC_EXTERNS() Msg {
	return enum_Msg_JSDOC_EXTERNS
}
func StaticGet_Msg_JSDOC_EXTRAVERSION() Msg {
	return enum_Msg_JSDOC_EXTRAVERSION
}
func StaticGet_Msg_JSDOC_EXTRA_VISIBILITY() Msg {
	return enum_Msg_JSDOC_EXTRA_VISIBILITY
}
func StaticGet_Msg_JSDOC_FILEOVERVIEW_EXTRA() Msg {
	return enum_Msg_JSDOC_FILEOVERVIEW_EXTRA
}
func StaticGet_Msg_JSDOC_FINAL() Msg {
	return enum_Msg_JSDOC_FINAL
}
func StaticGet_Msg_JSDOC_FUNCTION_NEWNOTOBJECT() Msg {
	return enum_Msg_JSDOC_FUNCTION_NEWNOTOBJECT
}
func StaticGet_Msg_JSDOC_FUNCTION_VARARGS() Msg {
	return enum_Msg_JSDOC_FUNCTION_VARARGS
}
func StaticGet_Msg_JSDOC_IDGEN_BAD() Msg {
	return enum_Msg_JSDOC_IDGEN_BAD
}
func StaticGet_Msg_JSDOC_IDGEN_DUPLICATE() Msg {
	return enum_Msg_JSDOC_IDGEN_DUPLICATE
}
func StaticGet_Msg_JSDOC_IDGEN_UNKNOWN() Msg {
	return enum_Msg_JSDOC_IDGEN_UNKNOWN
}
func StaticGet_Msg_JSDOC_IMPLEMENTS_DUPLICATE() Msg {
	return enum_Msg_JSDOC_IMPLEMENTS_DUPLICATE
}
func StaticGet_Msg_JSDOC_IMPLEMENTS_EXTRAQUALIFIER() Msg {
	return enum_Msg_JSDOC_IMPLEMENTS_EXTRAQUALIFIER
}
func StaticGet_Msg_JSDOC_IMPLICITCAST() Msg {
	return enum_Msg_JSDOC_IMPLICITCAST
}
func StaticGet_Msg_JSDOC_IMPORT() Msg {
	return enum_Msg_JSDOC_IMPORT
}
func StaticGet_Msg_JSDOC_INCOMPAT_TYPE() Msg {
	return enum_Msg_JSDOC_INCOMPAT_TYPE
}
func StaticGet_Msg_JSDOC_INCOMPAT_INLINING() Msg {
	return enum_Msg_JSDOC_INCOMPAT_INLINING
}
func StaticGet_Msg_JSDOC_INTERFACE_CONSTRUCTOR() Msg {
	return enum_Msg_JSDOC_INTERFACE_CONSTRUCTOR
}
func StaticGet_Msg_JSDOC_LEGACY_TS_EXPORT() Msg {
	return enum_Msg_JSDOC_LEGACY_TS_EXPORT
}
func StaticGet_Msg_JSDOC_LENDS_INCOMPATIBLE() Msg {
	return enum_Msg_JSDOC_LENDS_INCOMPATIBLE
}
func StaticGet_Msg_JSDOC_LENDS_MISSING() Msg {
	return enum_Msg_JSDOC_LENDS_MISSING
}
func StaticGet_Msg_JSDOC_LOCALEFILE() Msg {
	return enum_Msg_JSDOC_LOCALEFILE
}
func StaticGet_Msg_JSDOC_LOCALEOBJECT() Msg {
	return enum_Msg_JSDOC_LOCALEOBJECT
}
func StaticGet_Msg_JSDOC_LOCALESELECT() Msg {
	return enum_Msg_JSDOC_LOCALESELECT
}
func StaticGet_Msg_JSDOC_LOCALEVALUE() Msg {
	return enum_Msg_JSDOC_LOCALEVALUE
}
func StaticGet_Msg_JSDOC_PROVIDE_GOOG() Msg {
	return enum_Msg_JSDOC_PROVIDE_GOOG
}
func StaticGet_Msg_JSDOC_PROVIDE_ALREADY_PROVIDED() Msg {
	return enum_Msg_JSDOC_PROVIDE_ALREADY_PROVIDED
}
func StaticGet_Msg_JSDOC_MEANING_EXTRA() Msg {
	return enum_Msg_JSDOC_MEANING_EXTRA
}
func StaticGet_Msg_JSDOC_MISSING_BRACES() Msg {
	return enum_Msg_JSDOC_MISSING_BRACES
}
func StaticGet_Msg_JSDOC_MISSING_COLON() Msg {
	return enum_Msg_JSDOC_MISSING_COLON
}
func StaticGet_Msg_JSDOC_MISSING_GT() Msg {
	return enum_Msg_JSDOC_MISSING_GT
}
func StaticGet_Msg_JSDOC_MISSING_LC() Msg {
	return enum_Msg_JSDOC_MISSING_LC
}
func StaticGet_Msg_JSDOC_MISSING_LP() Msg {
	return enum_Msg_JSDOC_MISSING_LP
}
func StaticGet_Msg_JSDOC_MISSING_RB() Msg {
	return enum_Msg_JSDOC_MISSING_RB
}
func StaticGet_Msg_JSDOC_MISSING_RC() Msg {
	return enum_Msg_JSDOC_MISSING_RC
}
func StaticGet_Msg_JSDOC_MISSING_RP() Msg {
	return enum_Msg_JSDOC_MISSING_RP
}
func StaticGet_Msg_JSDOC_MISSING_TYPE_DECLARATION() Msg {
	return enum_Msg_JSDOC_MISSING_TYPE_DECLARATION
}
func StaticGet_Msg_JSDOC_MIXINCLASS_EXTRA() Msg {
	return enum_Msg_JSDOC_MIXINCLASS_EXTRA
}
func StaticGet_Msg_JSDOC_MIXINFUNCTION_EXTRA() Msg {
	return enum_Msg_JSDOC_MIXINFUNCTION_EXTRA
}
func StaticGet_Msg_JSDOC_MODIFIES() Msg {
	return enum_Msg_JSDOC_MODIFIES
}
func StaticGet_Msg_JSDOC_MODIFIES_DUPLICATE() Msg {
	return enum_Msg_JSDOC_MODIFIES_DUPLICATE
}
func StaticGet_Msg_JSDOC_MODIFIES_UNKNOWN() Msg {
	return enum_Msg_JSDOC_MODIFIES_UNKNOWN
}
func StaticGet_Msg_JSDOC_MODS() Msg {
	return enum_Msg_JSDOC_MODS
}
func StaticGet_Msg_JSDOC_MODS_EXTRA() Msg {
	return enum_Msg_JSDOC_MODS_EXTRA
}
func StaticGet_Msg_JSDOC_NAME_SYNTAX() Msg {
	return enum_Msg_JSDOC_NAME_SYNTAX
}
func StaticGet_Msg_JSDOC_NGINJECT_EXTRA() Msg {
	return enum_Msg_JSDOC_NGINJECT_EXTRA
}
func StaticGet_Msg_JSDOC_NOCOLLAPSE() Msg {
	return enum_Msg_JSDOC_NOCOLLAPSE
}
func StaticGet_Msg_JSDOC_NOCOMPILE() Msg {
	return enum_Msg_JSDOC_NOCOMPILE
}
func StaticGet_Msg_JSDOC_NOCOVERAGE() Msg {
	return enum_Msg_JSDOC_NOCOVERAGE
}
func StaticGet_Msg_JSDOC_NODTS() Msg {
	return enum_Msg_JSDOC_NODTS
}
func StaticGet_Msg_JSDOC_NOINLINE() Msg {
	return enum_Msg_JSDOC_NOINLINE
}
func StaticGet_Msg_JSDOC_NOSIDEEFFECTS() Msg {
	return enum_Msg_JSDOC_NOSIDEEFFECTS
}
func StaticGet_Msg_JSDOC_NOSIDEEFFECTS_WITH_MODIFIES() Msg {
	return enum_Msg_JSDOC_NOSIDEEFFECTS_WITH_MODIFIES
}
func StaticGet_Msg_JSDOC_NOSIDEEFFECTS_WITH_THROWS() Msg {
	return enum_Msg_JSDOC_NOSIDEEFFECTS_WITH_THROWS
}
func StaticGet_Msg_JSDOC_OVERRIDE() Msg {
	return enum_Msg_JSDOC_OVERRIDE
}
func StaticGet_Msg_JSDOC_POLYMERBEHAVIOR_EXTRA() Msg {
	return enum_Msg_JSDOC_POLYMERBEHAVIOR_EXTRA
}
func StaticGet_Msg_JSDOC_POLYMER_EXTRA() Msg {
	return enum_Msg_JSDOC_POLYMER_EXTRA
}
func StaticGet_Msg_JSDOC_PUREORBREAKMYCODE() Msg {
	return enum_Msg_JSDOC_PUREORBREAKMYCODE
}
func StaticGet_Msg_JSDOC_RECORD() Msg {
	return enum_Msg_JSDOC_RECORD
}
func StaticGet_Msg_JSDOC_SASS_GENERATED_CSS_TS() Msg {
	return enum_Msg_JSDOC_SASS_GENERATED_CSS_TS
}
func StaticGet_Msg_JSDOC_SEEMISSING() Msg {
	return enum_Msg_JSDOC_SEEMISSING
}
func StaticGet_Msg_JSDOC_SUPPRESS() Msg {
	return enum_Msg_JSDOC_SUPPRESS
}
func StaticGet_Msg_JSDOC_SUPPRESS_UNKNOWN() Msg {
	return enum_Msg_JSDOC_SUPPRESS_UNKNOWN
}
func StaticGet_Msg_JSDOC_TEMPLATE_BOUNDEDGENERICS_USED() Msg {
	return enum_Msg_JSDOC_TEMPLATE_BOUNDEDGENERICS_USED
}
func StaticGet_Msg_JSDOC_TEMPLATE_BOUNDSWITHTTL() Msg {
	return enum_Msg_JSDOC_TEMPLATE_BOUNDSWITHTTL
}
func StaticGet_Msg_JSDOC_TEMPLATE_MULTIPLEDECLARATION() Msg {
	return enum_Msg_JSDOC_TEMPLATE_MULTIPLEDECLARATION
}
func StaticGet_Msg_JSDOC_TEMPLATE_NAME_MISSING() Msg {
	return enum_Msg_JSDOC_TEMPLATE_NAME_MISSING
}
func StaticGet_Msg_JSDOC_TEMPLATE_NAME_REDECLARATION() Msg {
	return enum_Msg_JSDOC_TEMPLATE_NAME_REDECLARATION
}
func StaticGet_Msg_JSDOC_TEMPLATE_TYPETRANSFORMATION_EXPRESSIONMISSING() Msg {
	return enum_Msg_JSDOC_TEMPLATE_TYPETRANSFORMATION_EXPRESSIONMISSING
}
func StaticGet_Msg_JSDOC_TEMPLATE_TYPETRANSFORMATION_MISSINGDELIMIIER() Msg {
	return enum_Msg_JSDOC_TEMPLATE_TYPETRANSFORMATION_MISSINGDELIMIIER
}
func StaticGet_Msg_JSDOC_THROWS() Msg {
	return enum_Msg_JSDOC_THROWS
}
func StaticGet_Msg_JSDOC_TYPE() Msg {
	return enum_Msg_JSDOC_TYPE
}
func StaticGet_Msg_JSDOC_TYPESUMMARY() Msg {
	return enum_Msg_JSDOC_TYPESUMMARY
}
func StaticGet_Msg_JSDOC_TYPETRANSFORMATION_EXTRA_PARAM() Msg {
	return enum_Msg_JSDOC_TYPETRANSFORMATION_EXTRA_PARAM
}
func StaticGet_Msg_JSDOC_TYPETRANSFORMATION_INVALID() Msg {
	return enum_Msg_JSDOC_TYPETRANSFORMATION_INVALID
}
func StaticGet_Msg_JSDOC_TYPETRANSFORMATION_INVALID_EXPRESSION() Msg {
	return enum_Msg_JSDOC_TYPETRANSFORMATION_INVALID_EXPRESSION
}
func StaticGet_Msg_JSDOC_TYPETRANSFORMATION_INVALID_INSIDE() Msg {
	return enum_Msg_JSDOC_TYPETRANSFORMATION_INVALID_INSIDE
}
func StaticGet_Msg_JSDOC_TYPETRANSFORMATION_MISSING_PARAM() Msg {
	return enum_Msg_JSDOC_TYPETRANSFORMATION_MISSING_PARAM
}
func StaticGet_Msg_JSDOC_TYPE_RECORD_DUPLICATE() Msg {
	return enum_Msg_JSDOC_TYPE_RECORD_DUPLICATE
}
func StaticGet_Msg_JSDOC_TYPE_SYNTAX() Msg {
	return enum_Msg_JSDOC_TYPE_SYNTAX
}
func StaticGet_Msg_JSDOC_USEDVIADOTCONSTRUCTOR() Msg {
	return enum_Msg_JSDOC_USEDVIADOTCONSTRUCTOR
}
func StaticGet_Msg_JSDOC_UNNECESSARY_BRACES() Msg {
	return enum_Msg_JSDOC_UNNECESSARY_BRACES
}
func StaticGet_Msg_JSDOC_VERSIONMISSING() Msg {
	return enum_Msg_JSDOC_VERSIONMISSING
}
func StaticGet_Msg_JSDOC_WIZACTION() Msg {
	return enum_Msg_JSDOC_WIZACTION
}
func StaticGet_Msg_MISSING_VARIABLE_NAME() Msg {
	return enum_Msg_MISSING_VARIABLE_NAME
}
func StaticGet_Msg_NO_TYPE_NAME() Msg {
	return enum_Msg_NO_TYPE_NAME
}
func StaticGet_Msg_UNEXPECTED_EOF() Msg {
	return enum_Msg_UNEXPECTED_EOF
}
func StaticGet_Msg_JSDOC_WIZCALLBACK() Msg {
	return enum_Msg_JSDOC_WIZCALLBACK
}

func Msg_Values() []Msg {
	return j2gort.CloneSlice(enumValues_Msg)
}

func Msg_ValueOf(name string) Msg {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Msg {
		if enumBaseTable_Msg[i_].Name_ == name {
			return Msg(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.Msg." + name))
}

// Java method: Msg.format()
func (this_ Msg) Format_0() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return this_.FieldGet_text()
}

// Java method: Msg.format(java.lang.Object[] args)
func (this_ Msg) Format_1(args []any) string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = args
	var s string = this_.FieldGet_text()
	_ = s
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(args)); i = i + 1 {
			var toReplace string = ("{" + j2gort.StrValInt32(i) + "}")
			_ = toReplace
			s = javalang.StringReplace(s, j2gort.CheckNotNullString(toReplace), j2gort.StrVal(args[i]))
		}
	}
	return s
}

func (this_ Msg) IsJavaNil_() bool                              { return this_ == 0 }
func (Msg) HasIdentityEquals_()                                 {}
func (Msg) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Msg) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Msg) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Msg) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.Msg", "Msg" }
func (this_ Msg) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Msg[this_-1]
}
func (this_ Msg) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Msg) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Msg) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Msg[this_-1].Name_
}
func (this_ Msg) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Msg[this_-1].Name_
}
func (this_ Msg) Equals(other_ any) bool {
	o_, ok_ := other_.(Msg)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Msg) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Msg[this_-1].HashCode()
}
func (this_ Msg) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Msg[this_-1].CompareTo(other_)
}
func (this_ Msg) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Msg[this_-1].Name_)
}
func init() {
	vals_ := []any{Msg(1), Msg(2), Msg(3), Msg(4), Msg(5), Msg(6), Msg(7), Msg(8), Msg(9), Msg(10), Msg(11), Msg(12), Msg(13), Msg(14), Msg(15), Msg(16), Msg(17), Msg(18), Msg(19), Msg(20), Msg(21), Msg(22), Msg(23), Msg(24), Msg(25), Msg(26), Msg(27), Msg(28), Msg(29), Msg(30), Msg(31), Msg(32), Msg(33), Msg(34), Msg(35), Msg(36), Msg(37), Msg(38), Msg(39), Msg(40), Msg(41), Msg(42), Msg(43), Msg(44), Msg(45), Msg(46), Msg(47), Msg(48), Msg(49), Msg(50), Msg(51), Msg(52), Msg(53), Msg(54), Msg(55), Msg(56), Msg(57), Msg(58), Msg(59), Msg(60), Msg(61), Msg(62), Msg(63), Msg(64), Msg(65), Msg(66), Msg(67), Msg(68), Msg(69), Msg(70), Msg(71), Msg(72), Msg(73), Msg(74), Msg(75), Msg(76), Msg(77), Msg(78), Msg(79), Msg(80), Msg(81), Msg(82), Msg(83), Msg(84), Msg(85), Msg(86), Msg(87), Msg(88), Msg(89), Msg(90), Msg(91), Msg(92), Msg(93), Msg(94), Msg(95), Msg(96), Msg(97), Msg(98), Msg(99), Msg(100), Msg(101), Msg(102), Msg(103), Msg(104), Msg(105), Msg(106), Msg(107), Msg(108), Msg(109), Msg(110)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.Msg", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Msg", "com.google.javascript.rhino.Msg", func(v_ any) bool {
		return j2gort.InstanceOf[Msg](v_)
	})
}
