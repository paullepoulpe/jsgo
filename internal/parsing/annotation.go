// Port of com/google/javascript/jscomp/parsing/Annotation.java @ closure-compiler v20260415

package parsing

import "github.com/paullepoulpe/jsgo/internal/jutil"

// Annotation is the set of all natively recognized JSDoc annotations.
type Annotation int32

const (
	Annotation_ABSTRACT Annotation = iota
	Annotation_ALTERNATE_MESSAGE_ID
	Annotation_AUTHOR
	Annotation_CLOSURE_PRIMITIVE
	Annotation_COLLAPSIBLE_OR_BREAK_MY_CODE
	Annotation_CONSTANT
	Annotation_CONSTRUCTOR
	Annotation_CUSTOM_ELEMENT
	Annotation_RECORD
	Annotation_DEFINE
	Annotation_DEPRECATED
	Annotation_DESC
	Annotation_DICT
	Annotation_ENCOURAGE_INLINING
	Annotation_ENHANCE
	Annotation_ENUM
	Annotation_EXTENDS
	Annotation_EXTERNS
	Annotation_EXPORT
	Annotation_FILE_OVERVIEW
	Annotation_FINAL
	Annotation_IDGENERATOR
	Annotation_IMPLEMENTS
	Annotation_IMPLICIT_CAST
	Annotation_INHERIT_DOC
	Annotation_INTERFACE
	Annotation_JSX
	Annotation_JSX_FRAGMENT
	Annotation_LENDS
	Annotation_LICENSE // same as preserve
	Annotation_LOG_TYPE_IN_COMPILER
	Annotation_MAY_HAVE_EXTRA_EDGE
	Annotation_MEANING
	Annotation_MIXIN_CLASS
	Annotation_MIXIN_FUNCTION
	Annotation_MODIFIES
	Annotation_MODS
	Annotation_NG_INJECT
	Annotation_NO_COLLAPSE
	Annotation_NO_COMPILE
	Annotation_NO_COVERAGE
	// Annotation_NO_DTS is a tag to suppress clutz's d.ts generation for classes
	// and method. This is specifically for the use of J2CL.
	//
	// Annotating classes, class methods and fields with @nodts has a side
	// effect of not triggering a hard error on the code which the extended
	// subclass unintentionally reuses the same property from base class. The
	// author of the code using this tag should be specifically aware of this
	// and should be able to communicate this to the consumers of their code.
	Annotation_NO_DTS
	Annotation_NO_INLINE
	Annotation_NO_SIDE_EFFECTS
	Annotation_NOT_IMPLEMENTED
	Annotation_OVERRIDE
	Annotation_PACKAGE
	Annotation_PARAM
	Annotation_POLYMER
	Annotation_POLYMER_BEHAVIOR
	Annotation_PRESERVE // same as license
	Annotation_PRIVATE
	Annotation_PROTECTED
	Annotation_PROVIDE_GOOG // @provideGoog - appears only in base.js
	Annotation_PROVIDE_ALREADY_PROVIDED
	Annotation_PUBLIC
	Annotation_PURE_OR_BREAK_MY_CODE
	Annotation_REQUIRE_INLINING
	Annotation_RETURN
	Annotation_SASS_GENERATED_CSS_TS
	Annotation_SEE
	Annotation_SOY_MODULE
	Annotation_SOY_TEMPLATE
	Annotation_STRUCT
	Annotation_SUPPRESS
	Annotation_TEMPLATE
	Annotation_CLOSURE_UNAWARE_CODE
	Annotation_THIS
	Annotation_THROWS
	Annotation_TYPE
	Annotation_TYPEDEF
	Annotation_TYPE_SUMMARY
	Annotation_UNRESTRICTED
	Annotation_USED_VIA_DOT_CONSTRUCTOR
	Annotation_WIZACTION
	Annotation_TS_TYPE
	Annotation_WIZ_ANALYZER
	Annotation_WIZCALLBACK
)

var annotation_names = [...]string{
	"ABSTRACT",
	"ALTERNATE_MESSAGE_ID",
	"AUTHOR",
	"CLOSURE_PRIMITIVE",
	"COLLAPSIBLE_OR_BREAK_MY_CODE",
	"CONSTANT",
	"CONSTRUCTOR",
	"CUSTOM_ELEMENT",
	"RECORD",
	"DEFINE",
	"DEPRECATED",
	"DESC",
	"DICT",
	"ENCOURAGE_INLINING",
	"ENHANCE",
	"ENUM",
	"EXTENDS",
	"EXTERNS",
	"EXPORT",
	"FILE_OVERVIEW",
	"FINAL",
	"IDGENERATOR",
	"IMPLEMENTS",
	"IMPLICIT_CAST",
	"INHERIT_DOC",
	"INTERFACE",
	"JSX",
	"JSX_FRAGMENT",
	"LENDS",
	"LICENSE",
	"LOG_TYPE_IN_COMPILER",
	"MAY_HAVE_EXTRA_EDGE",
	"MEANING",
	"MIXIN_CLASS",
	"MIXIN_FUNCTION",
	"MODIFIES",
	"MODS",
	"NG_INJECT",
	"NO_COLLAPSE",
	"NO_COMPILE",
	"NO_COVERAGE",
	"NO_DTS",
	"NO_INLINE",
	"NO_SIDE_EFFECTS",
	"NOT_IMPLEMENTED",
	"OVERRIDE",
	"PACKAGE",
	"PARAM",
	"POLYMER",
	"POLYMER_BEHAVIOR",
	"PRESERVE",
	"PRIVATE",
	"PROTECTED",
	"PROVIDE_GOOG",
	"PROVIDE_ALREADY_PROVIDED",
	"PUBLIC",
	"PURE_OR_BREAK_MY_CODE",
	"REQUIRE_INLINING",
	"RETURN",
	"SASS_GENERATED_CSS_TS",
	"SEE",
	"SOY_MODULE",
	"SOY_TEMPLATE",
	"STRUCT",
	"SUPPRESS",
	"TEMPLATE",
	"CLOSURE_UNAWARE_CODE",
	"THIS",
	"THROWS",
	"TYPE",
	"TYPEDEF",
	"TYPE_SUMMARY",
	"UNRESTRICTED",
	"USED_VIA_DOT_CONSTRUCTOR",
	"WIZACTION",
	"TS_TYPE",
	"WIZ_ANALYZER",
	"WIZCALLBACK",
}

// Annotation_values is Annotation.values().
//
// Java: Annotation.values()
func Annotation_values() []Annotation {
	values := make([]Annotation, len(annotation_names))
	for i := range values {
		values[i] = Annotation(i)
	}
	return values
}

// Name is Enum.name().
//
// Java: Annotation.name()
func (a Annotation) Name() string { return annotation_names[a] }

// String is Enum.toString().
//
// Java: Annotation.toString()
func (a Annotation) String() string { return a.Name() }

// Java: Annotation.recognizedAnnotations
var Annotation_recognizedAnnotations = func() *jutil.LinkedHashMap[string, Annotation] {
	b := newConfig_immutableMapBuilder()
	b.put("ngInject", Annotation_NG_INJECT)
	b.put("abstract", Annotation_ABSTRACT)
	b.put("alternateMessageId", Annotation_ALTERNATE_MESSAGE_ID)
	b.put("argument", Annotation_PARAM)
	b.put("author", Annotation_AUTHOR)
	b.put("closurePrimitive", Annotation_CLOSURE_PRIMITIVE)
	b.put("closureUnaware", Annotation_CLOSURE_UNAWARE_CODE)
	b.put("const", Annotation_CONSTANT)
	b.put("collapsibleOrBreakMyCode", Annotation_COLLAPSIBLE_OR_BREAK_MY_CODE)
	b.put("constant", Annotation_CONSTANT)
	b.put("constructor", Annotation_CONSTRUCTOR)
	b.put("customElement", Annotation_CUSTOM_ELEMENT)
	b.put("copyright", Annotation_LICENSE)
	b.put("define", Annotation_DEFINE)
	b.put("deprecated", Annotation_DEPRECATED)
	b.put("desc", Annotation_DESC)
	b.put("dict", Annotation_DICT)
	b.put("requireInlining", Annotation_REQUIRE_INLINING)
	b.put("encourageInlining", Annotation_ENCOURAGE_INLINING)
	b.put("enum", Annotation_ENUM)
	b.put("enhance", Annotation_ENHANCE)
	b.put("export", Annotation_EXPORT)
	b.put("extends", Annotation_EXTENDS)
	b.put("externs", Annotation_EXTERNS)
	b.put("fileoverview", Annotation_FILE_OVERVIEW)
	b.put("final", Annotation_FINAL)
	b.put("idGenerator", Annotation_IDGENERATOR)
	b.put("implements", Annotation_IMPLEMENTS)
	b.put("implicitCast", Annotation_IMPLICIT_CAST)
	b.put("inheritDoc", Annotation_INHERIT_DOC)
	b.put("interface", Annotation_INTERFACE)
	b.put("record", Annotation_RECORD)
	b.put("lends", Annotation_LENDS)
	b.put("license", Annotation_LICENSE)
	b.put("logTypeInCompiler", Annotation_LOG_TYPE_IN_COMPILER)
	b.put("mayhaveextraedge", Annotation_MAY_HAVE_EXTRA_EDGE)
	b.put("meaning", Annotation_MEANING)
	b.put("mixinClass", Annotation_MIXIN_CLASS)
	b.put("mixinFunction", Annotation_MIXIN_FUNCTION)
	b.put("modifies", Annotation_MODIFIES)
	b.put("mods", Annotation_MODS)
	b.put("nocollapse", Annotation_NO_COLLAPSE)
	b.put("nocompile", Annotation_NO_COMPILE)
	b.put("nocoverage", Annotation_NO_COVERAGE)
	b.put("nodts", Annotation_NO_DTS)
	b.put("noinline", Annotation_NO_INLINE)
	b.put("nosideeffects", Annotation_NO_SIDE_EFFECTS)
	b.put("override", Annotation_OVERRIDE)
	b.put("owner", Annotation_AUTHOR)
	b.put("package", Annotation_PACKAGE)
	b.put("param", Annotation_PARAM)
	b.put("polymer", Annotation_POLYMER)
	b.put("polymerBehavior", Annotation_POLYMER_BEHAVIOR)
	b.put("preserve", Annotation_PRESERVE)
	b.put("private", Annotation_PRIVATE)
	b.put("protected", Annotation_PROTECTED)
	b.put("provideGoog", Annotation_PROVIDE_GOOG)
	b.put("provideAlreadyProvided", Annotation_PROVIDE_ALREADY_PROVIDED)
	b.put("public", Annotation_PUBLIC)
	b.put("pureOrBreakMyCode", Annotation_PURE_OR_BREAK_MY_CODE)
	b.put("return", Annotation_RETURN)
	b.put("returns", Annotation_RETURN)
	b.put("sassGeneratedCssTs", Annotation_SASS_GENERATED_CSS_TS)
	b.put("see", Annotation_SEE)
	b.put("soyModule", Annotation_SOY_MODULE)
	b.put("soyTemplate", Annotation_SOY_TEMPLATE)
	b.put("struct", Annotation_STRUCT)
	b.put("suppress", Annotation_SUPPRESS)
	b.put("template", Annotation_TEMPLATE)
	b.put("this", Annotation_THIS)
	b.put("throws", Annotation_THROWS)
	b.put("type", Annotation_TYPE)
	b.put("typedef", Annotation_TYPEDEF)
	b.put("typeSummary", Annotation_TYPE_SUMMARY)
	b.put("unrestricted", Annotation_UNRESTRICTED)
	b.put("usedViaDotConstructor", Annotation_USED_VIA_DOT_CONSTRUCTOR)
	b.put("wizaction", Annotation_WIZACTION)
	b.put("tsType", Annotation_TS_TYPE)
	b.put("wizAnalyzer", Annotation_WIZ_ANALYZER)
	b.put("wizcallback", Annotation_WIZCALLBACK)
	return b.buildOrThrow()
}()
