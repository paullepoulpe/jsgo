// Port of com/google/javascript/jscomp/parsing/Annotation.java @ closure-compiler v20261021

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

package parsing

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

// Name is Enum.name().
//
// Java: Annotation.name()
func (a Annotation) Name() string { return annotation_names[a] }

// String is Enum.toString().
//
// Java: Annotation.toString()
func (a Annotation) String() string { return a.Name() }

// Java: Annotation.recognizedAnnotations
var Annotation_recognizedAnnotations = map[string]Annotation{
	"ngInject":                 Annotation_NG_INJECT,
	"abstract":                 Annotation_ABSTRACT,
	"alternateMessageId":       Annotation_ALTERNATE_MESSAGE_ID,
	"argument":                 Annotation_PARAM,
	"author":                   Annotation_AUTHOR,
	"closurePrimitive":         Annotation_CLOSURE_PRIMITIVE,
	"closureUnaware":           Annotation_CLOSURE_UNAWARE_CODE,
	"const":                    Annotation_CONSTANT,
	"collapsibleOrBreakMyCode": Annotation_COLLAPSIBLE_OR_BREAK_MY_CODE,
	"constant":                 Annotation_CONSTANT,
	"constructor":              Annotation_CONSTRUCTOR,
	"customElement":            Annotation_CUSTOM_ELEMENT,
	"copyright":                Annotation_LICENSE,
	"define":                   Annotation_DEFINE,
	"deprecated":               Annotation_DEPRECATED,
	"desc":                     Annotation_DESC,
	"dict":                     Annotation_DICT,
	"requireInlining":          Annotation_REQUIRE_INLINING,
	"encourageInlining":        Annotation_ENCOURAGE_INLINING,
	"enum":                     Annotation_ENUM,
	"enhance":                  Annotation_ENHANCE,
	"export":                   Annotation_EXPORT,
	"extends":                  Annotation_EXTENDS,
	"externs":                  Annotation_EXTERNS,
	"fileoverview":             Annotation_FILE_OVERVIEW,
	"final":                    Annotation_FINAL,
	"idGenerator":              Annotation_IDGENERATOR,
	"implements":               Annotation_IMPLEMENTS,
	"implicitCast":             Annotation_IMPLICIT_CAST,
	"inheritDoc":               Annotation_INHERIT_DOC,
	"interface":                Annotation_INTERFACE,
	"record":                   Annotation_RECORD,
	"lends":                    Annotation_LENDS,
	"license":                  Annotation_LICENSE,
	"logTypeInCompiler":        Annotation_LOG_TYPE_IN_COMPILER,
	"mayhaveextraedge":         Annotation_MAY_HAVE_EXTRA_EDGE,
	"meaning":                  Annotation_MEANING,
	"mixinClass":               Annotation_MIXIN_CLASS,
	"mixinFunction":            Annotation_MIXIN_FUNCTION,
	"modifies":                 Annotation_MODIFIES,
	"mods":                     Annotation_MODS,
	"nocollapse":               Annotation_NO_COLLAPSE,
	"nocompile":                Annotation_NO_COMPILE,
	"nocoverage":               Annotation_NO_COVERAGE,
	"nodts":                    Annotation_NO_DTS,
	"noinline":                 Annotation_NO_INLINE,
	"nosideeffects":            Annotation_NO_SIDE_EFFECTS,
	"override":                 Annotation_OVERRIDE,
	"owner":                    Annotation_AUTHOR,
	"package":                  Annotation_PACKAGE,
	"param":                    Annotation_PARAM,
	"polymer":                  Annotation_POLYMER,
	"polymerBehavior":          Annotation_POLYMER_BEHAVIOR,
	"preserve":                 Annotation_PRESERVE,
	"private":                  Annotation_PRIVATE,
	"protected":                Annotation_PROTECTED,
	"provideGoog":              Annotation_PROVIDE_GOOG,
	"provideAlreadyProvided":   Annotation_PROVIDE_ALREADY_PROVIDED,
	"public":                   Annotation_PUBLIC,
	"pureOrBreakMyCode":        Annotation_PURE_OR_BREAK_MY_CODE,
	"return":                   Annotation_RETURN,
	"returns":                  Annotation_RETURN,
	"sassGeneratedCssTs":       Annotation_SASS_GENERATED_CSS_TS,
	"see":                      Annotation_SEE,
	"soyModule":                Annotation_SOY_MODULE,
	"soyTemplate":              Annotation_SOY_TEMPLATE,
	"struct":                   Annotation_STRUCT,
	"suppress":                 Annotation_SUPPRESS,
	"template":                 Annotation_TEMPLATE,
	"this":                     Annotation_THIS,
	"throws":                   Annotation_THROWS,
	"type":                     Annotation_TYPE,
	"typedef":                  Annotation_TYPEDEF,
	"typeSummary":              Annotation_TYPE_SUMMARY,
	"unrestricted":             Annotation_UNRESTRICTED,
	"usedViaDotConstructor":    Annotation_USED_VIA_DOT_CONSTRUCTOR,
	"wizaction":                Annotation_WIZACTION,
	"tsType":                   Annotation_TS_TYPE,
	"wizAnalyzer":              Annotation_WIZ_ANALYZER,
	"wizcallback":              Annotation_WIZCALLBACK,
}

var annotation_recognizedOrder = []string{
	"ngInject",
	"abstract",
	"alternateMessageId",
	"argument",
	"author",
	"closurePrimitive",
	"closureUnaware",
	"const",
	"collapsibleOrBreakMyCode",
	"constant",
	"constructor",
	"customElement",
	"copyright",
	"define",
	"deprecated",
	"desc",
	"dict",
	"requireInlining",
	"encourageInlining",
	"enum",
	"enhance",
	"export",
	"extends",
	"externs",
	"fileoverview",
	"final",
	"idGenerator",
	"implements",
	"implicitCast",
	"inheritDoc",
	"interface",
	"record",
	"lends",
	"license",
	"logTypeInCompiler",
	"mayhaveextraedge",
	"meaning",
	"mixinClass",
	"mixinFunction",
	"modifies",
	"mods",
	"nocollapse",
	"nocompile",
	"nocoverage",
	"nodts",
	"noinline",
	"nosideeffects",
	"override",
	"owner",
	"package",
	"param",
	"polymer",
	"polymerBehavior",
	"preserve",
	"private",
	"protected",
	"provideGoog",
	"provideAlreadyProvided",
	"public",
	"pureOrBreakMyCode",
	"return",
	"returns",
	"sassGeneratedCssTs",
	"see",
	"soyModule",
	"soyTemplate",
	"struct",
	"suppress",
	"template",
	"this",
	"throws",
	"type",
	"typedef",
	"typeSummary",
	"unrestricted",
	"usedViaDotConstructor",
	"wizaction",
	"tsType",
	"wizAnalyzer",
	"wizcallback",
}
