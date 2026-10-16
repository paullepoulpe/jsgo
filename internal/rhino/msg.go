// Port of com/google/javascript/rhino/Msg.java @ closure-compiler v20260415

package rhino

import (
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Msg holds the error message constants.
type Msg int32

const (
	Msg_BAD_FILEOVERVIEW_VISIBIIITY_ANNOTATION Msg = iota
	Msg_BAD_JSDOC_TAG
	Msg_DUP_VARIABLE_NAME
	Msg_END_ANNOTATION_EXPECTED
	Msg_INVALID_VARIABLE_NAME
	Msg_JSDOC_ALTERNATEMESSAGEID_EXTRA
	Msg_JSDOC_REQUIRE_INLINING
	Msg_JSDOC_ENCOURAGE_INLINING
	Msg_JSDOC_AUTHORMISSING
	Msg_JSDOC_CLOSUREPRIMITIVE_EXTRA
	Msg_JSDOC_CLOSUREPRIMITIVE_INVALID
	Msg_JSDOC_CLOSUREPRIMITIVE_MISSING
	Msg_JSDOC_CLOSURE_UNAWARE_CODE_EXTRA
	Msg_JSDOC_CLOSURE_UNAWARE_CONFIG_INVALID_VALUE
	Msg_JSDOC_CLOSURE_UNAWARE_CODE_INVALID
	Msg_JSDOC_COLLAPSIBLEORBREAKMYCODE
	Msg_JSDOC_CONST
	Msg_JSDOC_CUSTOMELEMENT_EXTRA
	Msg_JSDOC_DEFINE
	Msg_JSDOC_DEPRECATED
	Msg_JSDOC_DESC_EXTRA
	Msg_JSDOC_EXPORT
	Msg_JSDOC_EXTENDS_DUPLICATE
	Msg_JSDOC_EXTERNS
	Msg_JSDOC_EXTRAVERSION
	Msg_JSDOC_EXTRA_VISIBILITY
	Msg_JSDOC_FILEOVERVIEW_EXTRA
	Msg_JSDOC_FINAL
	Msg_JSDOC_FUNCTION_NEWNOTOBJECT
	Msg_JSDOC_FUNCTION_VARARGS
	Msg_JSDOC_IDGEN_BAD
	Msg_JSDOC_IDGEN_DUPLICATE
	Msg_JSDOC_IDGEN_UNKNOWN
	Msg_JSDOC_IMPLEMENTS_DUPLICATE
	Msg_JSDOC_IMPLEMENTS_EXTRAQUALIFIER
	Msg_JSDOC_IMPLICITCAST
	Msg_JSDOC_IMPORT
	Msg_JSDOC_INCOMPAT_TYPE
	Msg_JSDOC_INCOMPAT_INLINING
	Msg_JSDOC_INTERFACE_CONSTRUCTOR
	Msg_JSDOC_LENDS_INCOMPATIBLE
	Msg_JSDOC_LENDS_MISSING
	Msg_JSDOC_LOCALEFILE
	Msg_JSDOC_LOCALEOBJECT
	Msg_JSDOC_LOCALESELECT
	Msg_JSDOC_LOCALEVALUE
	Msg_JSDOC_PROVIDE_GOOG
	Msg_JSDOC_PROVIDE_ALREADY_PROVIDED
	Msg_JSDOC_MEANING_EXTRA
	Msg_JSDOC_MISSING_BRACES
	Msg_JSDOC_MISSING_COLON
	Msg_JSDOC_MISSING_GT
	Msg_JSDOC_MISSING_LC
	Msg_JSDOC_MISSING_LP
	Msg_JSDOC_MISSING_RB
	Msg_JSDOC_MISSING_RC
	Msg_JSDOC_MISSING_RP
	Msg_JSDOC_MISSING_TYPE_DECLARATION
	Msg_JSDOC_MIXINCLASS_EXTRA
	Msg_JSDOC_MIXINFUNCTION_EXTRA
	Msg_JSDOC_MODIFIES
	Msg_JSDOC_MODIFIES_DUPLICATE
	Msg_JSDOC_MODIFIES_UNKNOWN
	Msg_JSDOC_MODS
	Msg_JSDOC_MODS_EXTRA
	Msg_JSDOC_NAME_SYNTAX
	Msg_JSDOC_NGINJECT_EXTRA
	Msg_JSDOC_NOCOLLAPSE
	Msg_JSDOC_NOCOMPILE
	Msg_JSDOC_NOCOVERAGE
	Msg_JSDOC_NODTS
	Msg_JSDOC_NOINLINE
	Msg_JSDOC_NOSIDEEFFECTS
	Msg_JSDOC_OVERRIDE
	Msg_JSDOC_POLYMERBEHAVIOR_EXTRA
	Msg_JSDOC_POLYMER_EXTRA
	Msg_JSDOC_PUREORBREAKMYCODE
	Msg_JSDOC_RECORD
	Msg_JSDOC_SASS_GENERATED_CSS_TS
	Msg_JSDOC_SEEMISSING
	Msg_JSDOC_SUPPRESS
	Msg_JSDOC_SUPPRESS_UNKNOWN
	Msg_JSDOC_TEMPLATE_BOUNDEDGENERICS_USED
	Msg_JSDOC_TEMPLATE_BOUNDSWITHTTL
	Msg_JSDOC_TEMPLATE_MULTIPLEDECLARATION
	Msg_JSDOC_TEMPLATE_NAME_MISSING
	Msg_JSDOC_TEMPLATE_NAME_REDECLARATION
	Msg_JSDOC_TEMPLATE_TYPETRANSFORMATION_EXPRESSIONMISSING
	Msg_JSDOC_TEMPLATE_TYPETRANSFORMATION_MISSINGDELIMIIER
	Msg_JSDOC_TYPE
	Msg_JSDOC_TYPESUMMARY
	Msg_JSDOC_TYPETRANSFORMATION_EXTRA_PARAM
	Msg_JSDOC_TYPETRANSFORMATION_INVALID
	Msg_JSDOC_TYPETRANSFORMATION_INVALID_EXPRESSION
	Msg_JSDOC_TYPETRANSFORMATION_INVALID_INSIDE
	Msg_JSDOC_TYPETRANSFORMATION_MISSING_PARAM
	Msg_JSDOC_TYPE_RECORD_DUPLICATE
	Msg_JSDOC_TYPE_SYNTAX
	Msg_JSDOC_USEDVIADOTCONSTRUCTOR
	Msg_JSDOC_UNNECESSARY_BRACES
	Msg_JSDOC_VERSIONMISSING
	Msg_JSDOC_WIZACTION
	Msg_MISSING_VARIABLE_NAME
	Msg_NO_TYPE_NAME
	Msg_UNEXPECTED_EOF
	Msg_JSDOC_WIZCALLBACK
)

var msg_names = [...]string{
	"BAD_FILEOVERVIEW_VISIBIIITY_ANNOTATION",
	"BAD_JSDOC_TAG",
	"DUP_VARIABLE_NAME",
	"END_ANNOTATION_EXPECTED",
	"INVALID_VARIABLE_NAME",
	"JSDOC_ALTERNATEMESSAGEID_EXTRA",
	"JSDOC_REQUIRE_INLINING",
	"JSDOC_ENCOURAGE_INLINING",
	"JSDOC_AUTHORMISSING",
	"JSDOC_CLOSUREPRIMITIVE_EXTRA",
	"JSDOC_CLOSUREPRIMITIVE_INVALID",
	"JSDOC_CLOSUREPRIMITIVE_MISSING",
	"JSDOC_CLOSURE_UNAWARE_CODE_EXTRA",
	"JSDOC_CLOSURE_UNAWARE_CONFIG_INVALID_VALUE",
	"JSDOC_CLOSURE_UNAWARE_CODE_INVALID",
	"JSDOC_COLLAPSIBLEORBREAKMYCODE",
	"JSDOC_CONST",
	"JSDOC_CUSTOMELEMENT_EXTRA",
	"JSDOC_DEFINE",
	"JSDOC_DEPRECATED",
	"JSDOC_DESC_EXTRA",
	"JSDOC_EXPORT",
	"JSDOC_EXTENDS_DUPLICATE",
	"JSDOC_EXTERNS",
	"JSDOC_EXTRAVERSION",
	"JSDOC_EXTRA_VISIBILITY",
	"JSDOC_FILEOVERVIEW_EXTRA",
	"JSDOC_FINAL",
	"JSDOC_FUNCTION_NEWNOTOBJECT",
	"JSDOC_FUNCTION_VARARGS",
	"JSDOC_IDGEN_BAD",
	"JSDOC_IDGEN_DUPLICATE",
	"JSDOC_IDGEN_UNKNOWN",
	"JSDOC_IMPLEMENTS_DUPLICATE",
	"JSDOC_IMPLEMENTS_EXTRAQUALIFIER",
	"JSDOC_IMPLICITCAST",
	"JSDOC_IMPORT",
	"JSDOC_INCOMPAT_TYPE",
	"JSDOC_INCOMPAT_INLINING",
	"JSDOC_INTERFACE_CONSTRUCTOR",
	"JSDOC_LENDS_INCOMPATIBLE",
	"JSDOC_LENDS_MISSING",
	"JSDOC_LOCALEFILE",
	"JSDOC_LOCALEOBJECT",
	"JSDOC_LOCALESELECT",
	"JSDOC_LOCALEVALUE",
	"JSDOC_PROVIDE_GOOG",
	"JSDOC_PROVIDE_ALREADY_PROVIDED",
	"JSDOC_MEANING_EXTRA",
	"JSDOC_MISSING_BRACES",
	"JSDOC_MISSING_COLON",
	"JSDOC_MISSING_GT",
	"JSDOC_MISSING_LC",
	"JSDOC_MISSING_LP",
	"JSDOC_MISSING_RB",
	"JSDOC_MISSING_RC",
	"JSDOC_MISSING_RP",
	"JSDOC_MISSING_TYPE_DECLARATION",
	"JSDOC_MIXINCLASS_EXTRA",
	"JSDOC_MIXINFUNCTION_EXTRA",
	"JSDOC_MODIFIES",
	"JSDOC_MODIFIES_DUPLICATE",
	"JSDOC_MODIFIES_UNKNOWN",
	"JSDOC_MODS",
	"JSDOC_MODS_EXTRA",
	"JSDOC_NAME_SYNTAX",
	"JSDOC_NGINJECT_EXTRA",
	"JSDOC_NOCOLLAPSE",
	"JSDOC_NOCOMPILE",
	"JSDOC_NOCOVERAGE",
	"JSDOC_NODTS",
	"JSDOC_NOINLINE",
	"JSDOC_NOSIDEEFFECTS",
	"JSDOC_OVERRIDE",
	"JSDOC_POLYMERBEHAVIOR_EXTRA",
	"JSDOC_POLYMER_EXTRA",
	"JSDOC_PUREORBREAKMYCODE",
	"JSDOC_RECORD",
	"JSDOC_SASS_GENERATED_CSS_TS",
	"JSDOC_SEEMISSING",
	"JSDOC_SUPPRESS",
	"JSDOC_SUPPRESS_UNKNOWN",
	"JSDOC_TEMPLATE_BOUNDEDGENERICS_USED",
	"JSDOC_TEMPLATE_BOUNDSWITHTTL",
	"JSDOC_TEMPLATE_MULTIPLEDECLARATION",
	"JSDOC_TEMPLATE_NAME_MISSING",
	"JSDOC_TEMPLATE_NAME_REDECLARATION",
	"JSDOC_TEMPLATE_TYPETRANSFORMATION_EXPRESSIONMISSING",
	"JSDOC_TEMPLATE_TYPETRANSFORMATION_MISSINGDELIMIIER",
	"JSDOC_TYPE",
	"JSDOC_TYPESUMMARY",
	"JSDOC_TYPETRANSFORMATION_EXTRA_PARAM",
	"JSDOC_TYPETRANSFORMATION_INVALID",
	"JSDOC_TYPETRANSFORMATION_INVALID_EXPRESSION",
	"JSDOC_TYPETRANSFORMATION_INVALID_INSIDE",
	"JSDOC_TYPETRANSFORMATION_MISSING_PARAM",
	"JSDOC_TYPE_RECORD_DUPLICATE",
	"JSDOC_TYPE_SYNTAX",
	"JSDOC_USEDVIADOTCONSTRUCTOR",
	"JSDOC_UNNECESSARY_BRACES",
	"JSDOC_VERSIONMISSING",
	"JSDOC_WIZACTION",
	"MISSING_VARIABLE_NAME",
	"NO_TYPE_NAME",
	"UNEXPECTED_EOF",
	"JSDOC_WIZCALLBACK",
}

var msg_texts = [...]string{
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

// Java: Msg.text
func (m Msg) Text() string { return msg_texts[m] }

// Msg_values is Msg.values().
func Msg_values() []Msg {
	out := make([]Msg, len(msg_names))
	for i := range out {
		out[i] = Msg(i)
	}
	return out
}

// String is Enum.name().
func (m Msg) String() string { return msg_names[m] }

// Java: Msg.format()
func (m Msg) Format() string {
	return m.Text()
}

// FormatArgs substitutes "{i}" with String.valueOf(args[i]).
//
// Java: Msg.format(Object...)
func (m Msg) FormatArgs(args ...any) string {
	// Note that this doesn't handle single-quote hence not compatible with
	// MessageFormat.
	s := m.Text()
	for i := 0; i < len(args); i++ {
		toReplace := "{" + strconv.Itoa(i) + "}"
		s = strings.ReplaceAll(s, toReplace, jutil.ToJavaString(args[i]))
	}
	return s
}
