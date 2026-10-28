// Port of com/google/javascript/jscomp/serialization/JSDocSerializer.java @ closure-compiler v20261021

/*
 * Copyright 2020 The Closure Compiler Authors.
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

package jscomp

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JSDocSerializer holds utilities for serializing and deserializing JSDoc
// necessary for optimzations.
type JSDocSerializer struct{}

// JSDocSerializer_ConvertJSDocInfoForOptimizations returns a variant of
// input JSDocInfo where fields not needed for optimizations are removed.
//
// This uses the serialization / deserialization logic for JSDoc and ensures
// optimizations can't accidentally depend on fields that we don't
// serialize.
//
// It returns a new JSDocInfo object or nil if no serializable fields are
// found.
//
// Java: JSDocSerializer.convertJSDocInfoForOptimizations(JSDocInfo)
func JSDocSerializer_ConvertJSDocInfoForOptimizations(jsdoc *rhino.JSDocInfo) *rhino.JSDocInfo {
	return jsDocSerializer_convertJSDocInfoWithBuilder(jsdoc, nil)
}

func jsDocSerializer_convertJSDocInfoWithBuilder(
	jsdoc *rhino.JSDocInfo,
	builder *rhino.JSDocInfo_Builder,
) *rhino.JSDocInfo {
	if jsdoc == nil {
		return nil
	}

	var kindBuf [8]int32
	kinds := jsDocSerializer_fillJsdocTagsOnly(jsdoc, kindBuf[:0])
	license, description, alternateMessageId, meaning, mode, hasMode := jsDocSerializer_stringFields(jsdoc)
	if len(kinds) == 0 &&
		license == "" &&
		description == "" &&
		alternateMessageId == "" &&
		meaning == "" &&
		!hasMode {
		return nil
	}
	return jsDocSerializer_buildConvertedJSDoc(
		kinds, license, description, alternateMessageId, meaning, mode, hasMode, builder)
}

func jsDocSerializer_buildConvertedJSDoc(
	kinds []int32,
	license, description, alternateMessageId, meaning string,
	mode rhino.JSDocInfo_PerFileClosureUnawareMode,
	hasMode bool,
	builder *rhino.JSDocInfo_Builder,
) *rhino.JSDocInfo {
	var localBuilder rhino.JSDocInfo_Builder
	if builder == nil {
		builder = &localBuilder
	}
	builder.Reset()

	if hasMode {
		builder.RecordClosureUnawareCodeWithMode(mode)
	}
	if license != "" {
		builder.RecordLicense(license)
	}
	if meaning != "" {
		builder.RecordMeaning(rhino.RhinoStringPool_AddOrGet(meaning))
	}
	if description != "" {
		builder.RecordDescription(rhino.RhinoStringPool_AddOrGet(description))
	}
	if alternateMessageId != "" {
		builder.RecordAlternateMessageId(rhino.RhinoStringPool_AddOrGet(alternateMessageId))
	}
	jsDocSerializer_applyTagsToBuilder(kinds, builder)
	return builder.Build()
}

// JSDocSerializer_SerializeJsdoc may return nil.
//
// Java: JSDocSerializer.serializeJsdoc(JSDocInfo, StringPool.Builder)
func JSDocSerializer_SerializeJsdoc(jsdoc *rhino.JSDocInfo, stringPool *StringPool_Builder) *OptimizationJsdoc {
	if jsdoc == nil {
		return nil
	}
	var kindBuf [8]int32
	kinds := jsDocSerializer_fillJsdocTagsOnly(jsdoc, kindBuf[:0])
	license, description, alternateMessageId, meaning, mode, hasMode := jsDocSerializer_stringFields(jsdoc)
	if len(kinds) == 0 &&
		license == "" &&
		description == "" &&
		alternateMessageId == "" &&
		meaning == "" &&
		!hasMode {
		return nil
	}
	return jsDocSerializer_allocProto(
		kinds, license, description, alternateMessageId, meaning, mode, hasMode, stringPool)
}

func jsDocSerializer_allocProto(
	kinds []int32,
	license, description, alternateMessageId, meaning string,
	mode rhino.JSDocInfo_PerFileClosureUnawareMode,
	hasMode bool,
	stringPool *StringPool_Builder,
) *OptimizationJsdoc {
	m := &OptimizationJsdoc{
		frozen: true,
	}
	if len(kinds) > 0 {
		m.kind_ = slices.Clone(kinds)
	}
	if license != "" {
		m.licenseTextPointer_ = int32(stringPool.Put(license))
	}
	if description != "" {
		m.descriptionPointer_ = int32(stringPool.Put(description))
	}
	if alternateMessageId != "" {
		m.alternateMessageIdPointer_ = int32(stringPool.Put(alternateMessageId))
	}
	if meaning != "" {
		m.meaningPointer_ = int32(stringPool.Put(meaning))
	}
	if hasMode {
		m.closureUnawareModePointer_ = int32(stringPool.Put(mode.String()))
	}
	m.memoizedSize = m.scalarsSerializedSize()
	return m
}

func jsDocSerializer_stringFields(jsdoc *rhino.JSDocInfo) (
	license, description, alternateMessageId, meaning string,
	mode rhino.JSDocInfo_PerFileClosureUnawareMode,
	hasMode bool,
) {
	license, _ = jsdoc.GetLicenseString()
	description, _ = jsdoc.GetDescriptionString()
	alternateMessageId, _ = jsdoc.GetAlternateMessageIdString()
	meaning, _ = jsdoc.GetMeaningString()
	mode, hasMode = jsdoc.GetPerFileClosureUnawareModeValue()
	return
}

func jsDocSerializer_fillJsdocTagsOnly(jsdoc *rhino.JSDocInfo, kindBuf []int32) []int32 {
	if jsdoc.HasFileOverview() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_FILEOVERVIEW_VALUE)
	}
	if jsdoc.IsSassGeneratedCssTs() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_SASS_GENERATED_CSS_TS_VALUE)
	}
	// Used by CoverageInstrumentationCallback
	if jsdoc.IsNoCoverage() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_NO_COVERAGE_VALUE)
	}
	if jsdoc.IsNoInline() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_NO_INLINE_VALUE)
	}
	if jsdoc.IsEncourageInlining() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_ENCOURAGE_INLINING_VALUE)
	}
	if jsdoc.IsRequireInlining() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_REQUIRE_INLINING_VALUE)
	}
	if jsdoc.IsNoCollapse() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_NO_COLLAPSE_VALUE)
	}
	if jsdoc.IsProvideGoog() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_PROVIDE_GOOG_VALUE)
	}
	if jsdoc.IsProvideAlreadyProvided() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_PROVIDE_ALREADY_PROVIDED_VALUE)
	}
	if jsdoc.IsTypeSummary() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_TYPE_SUMMARY_FILE_VALUE)
	}
	if jsdoc.IsPureOrBreakMyCode() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_PURE_OR_BREAK_MY_CODE_VALUE)
	}
	if jsdoc.IsCollapsibleOrBreakMyCode() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_COLLAPSIBLE_OR_BREAK_MY_CODE_VALUE)
	}
	if jsdoc.HasThisType() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_THIS_VALUE)
	}
	if jsdoc.HasEnumParameterType() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_ENUM_VALUE)
	}
	if jsdoc.IsDefine() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_DEFINE_VALUE)
	}
	if jsdoc.HasConstAnnotation() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_CONST_VALUE)
	}
	if jsdoc.IsAnyIdGenerator() {
		kindBuf = append(kindBuf, jsDocSerializer_serializeIdGenerator(jsdoc).GetNumber())
	}

	// Used by PureFunctionIdentifier
	if jsdoc.IsNoSideEffects() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_NO_SIDE_EFFECTS_VALUE)
	}
	if jsdoc.HasModifies() {
		for modifies := range jsdoc.GetModifies().All() {
			if modifies == "this" {
				kindBuf = append(kindBuf, JsdocTag_JSDOC_MODIFIES_THIS_VALUE)
			} else {
				// Currently, anything other than "this" is considered a modification to arguments
				kindBuf = append(kindBuf, JsdocTag_JSDOC_MODIFIES_ARGUMENTS_VALUE)
			}
		}
	}
	if jsdoc.HasThrowsAnnotations() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_THROWS_VALUE)
	}

	// Used by DevirtualizeMethods and CollapseProperties
	if jsdoc.IsConstructor() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_CONSTRUCTOR_VALUE)
	}
	if jsdoc.IsInterface() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_INTERFACE_VALUE)
	}
	hasSuppressions := jsdoc.HasAnySuppressions()
	if hasSuppressions && jsdoc.HasSuppression("partialAlias") {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_SUPPRESS_PARTIAL_ALIAS_VALUE)
	}

	// Used by ClosureCodeRemoval
	if jsdoc.IsAbstract() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_ABSTRACT_VALUE)
	}

	// Used by ReplaceMessages
	if hasSuppressions {
		if jsdoc.HasSuppression("messageConventions") {
			kindBuf = append(kindBuf, JsdocTag_JSDOC_SUPPRESS_MESSAGE_CONVENTION_VALUE)
		}
		if jsdoc.HasSuppression("untranspilableFeatures") {
			kindBuf = append(kindBuf, JsdocTag_JSDOC_SUPPRESS_UNTRANSPILABLE_FEATURES_VALUE)
		}
	}
	if jsdoc.IsUsedViaDotConstructor() {
		kindBuf = append(kindBuf, JsdocTag_JSDOC_USED_VIA_DOT_CONSTRUCTOR_VALUE)
	}
	return kindBuf
}

func jsDocSerializer_fillJsdocProto(jsdoc *rhino.JSDocInfo, stringPool *StringPool_Builder, m *OptimizationJsdoc, kindBuf []int32) bool {
	if jsdoc == nil {
		return false
	}

	m.kind_ = jsDocSerializer_fillJsdocTagsOnly(jsdoc, kindBuf)
	license, description, alternateMessageId, meaning, mode, hasMode := jsDocSerializer_stringFields(jsdoc)
	if license != "" {
		m.licenseTextPointer_ = int32(stringPool.Put(license))
	}
	if description != "" {
		m.descriptionPointer_ = int32(stringPool.Put(description))
	}
	if alternateMessageId != "" {
		m.alternateMessageIdPointer_ = int32(stringPool.Put(alternateMessageId))
	}
	if meaning != "" {
		m.meaningPointer_ = int32(stringPool.Put(meaning))
	}
	if hasMode {
		m.closureUnawareModePointer_ = int32(stringPool.Put(mode.String()))
	}

	return len(m.kind_) != 0 ||
		m.licenseTextPointer_ != 0 ||
		m.descriptionPointer_ != 0 ||
		m.meaningPointer_ != 0 ||
		m.alternateMessageIdPointer_ != 0 ||
		m.closureUnawareModePointer_ != 0
}

// Java: JSDocSerializer.serializeIdGenerator(JSDocInfo)
func jsDocSerializer_serializeIdGenerator(doc *rhino.JSDocInfo) JsdocTag {
	if doc.IsConsistentIdGenerator() {
		return JsdocTag_JSDOC_ID_GENERATOR_CONSISTENT
	} else if doc.IsStableIdGenerator() {
		return JsdocTag_JSDOC_ID_GENERATOR_STABLE
	} else if doc.IsXidGenerator() {
		return JsdocTag_JSDOC_ID_GENERATOR_XID
	} else if doc.IsMappedIdGenerator() {
		return JsdocTag_JSDOC_ID_GENERATOR_MAPPED
	} else if doc.IsIdGenerator() {
		return JsdocTag_JSDOC_ID_GENERATOR_INCONSISTENT
	}
	panic(jutil.NewIllegalStateExceptionMsg("Failed to identify idGenerator inside JSDoc: " + doc.String()))
}

// Optimizations shouldn't care about the contents of JSTypeExpressions but some JSDoc APIs
// expect their presence, so create a placeholder type.
var jsDocSerializer_SYNTHETIC_SOURCE = SourceFile_FromCode("JSDocSerializer_placeholder_source", "")

const jsDocSerializer_PLACEHOLDER_TYPE_NAME = "JSDocSerializer_placeholder_type"

// Java: JSDocSerializer.createPlaceholderType()
func jsDocSerializer_createPlaceholderType() *rhino.JSTypeExpression {
	// the BANG (!) token makes unit testing easier, as the JSDoc parser implicitly adds "!"
	// to some JSTypeExpressions
	name := rhino.IR_String(jsDocSerializer_PLACEHOLDER_TYPE_NAME)
	bang := rhino.NewNodeWithChild(rhino.Token_BANG, name)
	name.SetStaticSourceFile(jsDocSerializer_SYNTHETIC_SOURCE)
	bang.SetStaticSourceFile(jsDocSerializer_SYNTHETIC_SOURCE)
	return rhino.NewJSTypeExpression(bang, jsDocSerializer_SYNTHETIC_SOURCE.GetName())
}

var jsDocSerializer_placeholderType = jsDocSerializer_createPlaceholderType()

// jsDocSerializer_perFileClosureUnawareModeValueOf is
// PerFileClosureUnawareMode.valueOf(String).
//
// Go-only: language-adapter — Enum.valueOf.
func jsDocSerializer_perFileClosureUnawareModeValueOf(name string) rhino.JSDocInfo_PerFileClosureUnawareMode {
	switch name {
	case "UNSPECIFIED":
		return rhino.JSDocInfo_PerFileClosureUnawareMode_UNSPECIFIED
	case "SIMPLE":
		return rhino.JSDocInfo_PerFileClosureUnawareMode_SIMPLE
	case "WHITESPACE":
		return rhino.JSDocInfo_PerFileClosureUnawareMode_WHITESPACE
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"No enum constant com.google.javascript.rhino.JSDocInfo.PerFileClosureUnawareMode." + name))
	}
}

var (
	jsDocSerializer_modifiesArguments        = rhino.JSDocInfo_InternSet(jutil.NewLinkedHashSet("arguments"))
	jsDocSerializer_modifiesThis             = rhino.JSDocInfo_InternSet(jutil.NewLinkedHashSet("this"))
	jsDocSerializer_modifiesArgumentsAndThis = rhino.JSDocInfo_InternSet(jutil.NewLinkedHashSet("arguments", "this"))
	jsDocSerializer_throwsPlaceholder        = "{!" + jsDocSerializer_PLACEHOLDER_TYPE_NAME + "}"
)

// JSDocSerializer_DeserializeJsdoc may return nil.
//
// Java: JSDocSerializer.deserializeJsdoc(OptimizationJsdoc, StringPool)
func JSDocSerializer_DeserializeJsdoc(serializedJsdoc *OptimizationJsdoc, stringPool *StringPool) *rhino.JSDocInfo {
	if serializedJsdoc == nil || serializedJsdoc == OptimizationJsdoc_GetDefaultInstance() {
		return nil
	}
	var builder rhino.JSDocInfo_Builder
	return jsDocSerializer_deserializeJsdocWithBuilder(serializedJsdoc, stringPool, &builder)
}

func jsDocSerializer_deserializeJsdocWithBuilder(
	serializedJsdoc *OptimizationJsdoc,
	stringPool *StringPool,
	builder *rhino.JSDocInfo_Builder,
) *rhino.JSDocInfo {
	if serializedJsdoc == nil || serializedJsdoc == OptimizationJsdoc_GetDefaultInstance() {
		return nil
	}
	builder.Reset()

	if serializedJsdoc.GetClosureUnawareModePointer() != 0 {
		builder.RecordClosureUnawareCodeWithMode(
			jsDocSerializer_perFileClosureUnawareModeValueOf(
				stringPool.Get(int(serializedJsdoc.GetClosureUnawareModePointer()))))
	}
	if serializedJsdoc.GetLicenseTextPointer() != 0 {
		builder.RecordLicense(stringPool.Get(int(serializedJsdoc.GetLicenseTextPointer())))
	}
	if serializedJsdoc.GetMeaningPointer() != 0 {
		builder.RecordMeaning(stringPool.Get(int(serializedJsdoc.GetMeaningPointer())))
	}
	if serializedJsdoc.GetDescriptionPointer() != 0 {
		builder.RecordDescription(stringPool.Get(int(serializedJsdoc.GetDescriptionPointer())))
	}
	if serializedJsdoc.GetAlternateMessageIdPointer() != 0 {
		builder.RecordAlternateMessageId(
			stringPool.Get(int(serializedJsdoc.GetAlternateMessageIdPointer())))
	}

	jsDocSerializer_applyTagsToBuilder(serializedJsdoc.kind_, builder)
	return builder.Build()
}

func jsDocSerializer_applyTagsToBuilder(kinds []int32, builder *rhino.JSDocInfo_Builder) {
	// Bitmasks avoid allocating TreeSets/LinkedHashSets on the deserialization hot path
	// while preserving exact alphabetical TreeSet<String> order:
	// modifies: bit 0 = "arguments", bit 1 = "this"
	// suppressions: bit 0 = "messageConventions", bit 1 = "partialAlias", bit 2 = "untranspilableFeatures"
	var modMask uint8
	var suppMask uint8

	for i := 0; i < len(kinds); i++ {
		tag := jsdocTag_fromNumber(kinds[i])
		switch tag {
		case JsdocTag_JSDOC_CONST:
			builder.RecordConstancy()
		case JsdocTag_JSDOC_ENUM:
			builder.RecordEnumParameterType(jsDocSerializer_placeholderType)
		case JsdocTag_JSDOC_THIS:
			builder.RecordThisType(jsDocSerializer_placeholderType)
		case JsdocTag_JSDOC_NO_COLLAPSE:
			builder.RecordNoCollapse()
		case JsdocTag_JSDOC_NO_INLINE:
			builder.RecordNoInline()
		case JsdocTag_JSDOC_REQUIRE_INLINING:
			builder.RecordRequireInlining()
		case JsdocTag_JSDOC_ENCOURAGE_INLINING:
			builder.RecordEncourageInlining()
		case JsdocTag_JSDOC_PROVIDE_GOOG:
			builder.RecordProvideGoog()
		case JsdocTag_JSDOC_PROVIDE_ALREADY_PROVIDED:
			builder.RecordProvideAlreadyProvided()
		case JsdocTag_JSDOC_TYPE_SUMMARY_FILE:
			builder.RecordTypeSummary()
		case JsdocTag_JSDOC_PURE_OR_BREAK_MY_CODE:
			builder.RecordPureOrBreakMyCode()
		case JsdocTag_JSDOC_COLLAPSIBLE_OR_BREAK_MY_CODE:
			builder.RecordCollapsibleOrBreakMyCode()
		case JsdocTag_JSDOC_DEFINE:
			builder.RecordDefineType(jsDocSerializer_placeholderType)
		case JsdocTag_JSDOC_NO_SIDE_EFFECTS:
			builder.RecordNoSideEffects()
		case JsdocTag_JSDOC_MODIFIES_THIS:
			modMask |= 2
		case JsdocTag_JSDOC_MODIFIES_ARGUMENTS:
			modMask |= 1
		case JsdocTag_JSDOC_THROWS:
			builder.RecordThrowsAnnotation(jsDocSerializer_throwsPlaceholder)
		case JsdocTag_JSDOC_CONSTRUCTOR:
			builder.RecordConstructor()
		case JsdocTag_JSDOC_INTERFACE:
			builder.RecordInterface()
		case JsdocTag_JSDOC_SUPPRESS_PARTIAL_ALIAS:
			suppMask |= 2
		case JsdocTag_JSDOC_ID_GENERATOR_CONSISTENT:
			builder.RecordConsistentIdGenerator()
		case JsdocTag_JSDOC_ID_GENERATOR_STABLE:
			builder.RecordStableIdGenerator()
		case JsdocTag_JSDOC_ID_GENERATOR_MAPPED:
			builder.RecordMappedIdGenerator()
		case JsdocTag_JSDOC_ID_GENERATOR_XID:
			builder.RecordXidGenerator()
		case JsdocTag_JSDOC_ID_GENERATOR_INCONSISTENT:
			builder.RecordIdGenerator()
		case JsdocTag_JSDOC_ABSTRACT:
			builder.RecordAbstract()
		// TODO(lharker): stage 2 passes ideally shouldn't report diagnostics, so this could be
		// moved to stage 1.
		case JsdocTag_JSDOC_SUPPRESS_MESSAGE_CONVENTION:
			suppMask |= 1
		// ReportUntranspilableFeatures pass will run in stage2 since it uses languageOut
		// information. It reports diagnostic {@code UNTRANSPILABLE_FEATURE_PRESENT} that can be
		// supppressed using `untranspilableFeatures` suppression tag.
		// Hence we must propagate it.
		case JsdocTag_JSDOC_SUPPRESS_UNTRANSPILABLE_FEATURES:
			suppMask |= 4
		case JsdocTag_JSDOC_FILEOVERVIEW:
			builder.RecordFileOverview("")
		case JsdocTag_JSDOC_NO_COVERAGE:
			builder.RecordNoCoverage()
		case JsdocTag_JSDOC_SASS_GENERATED_CSS_TS:
			builder.RecordSassGeneratedCssTs()
		case JsdocTag_JSDOC_USED_VIA_DOT_CONSTRUCTOR:
			_ = builder.RecordUsedViaDotConstructor()
		case JsdocTag_JSDOC_UNSPECIFIED, JsdocTag_UNRECOGNIZED:
			panic(NewMalformedTypedAstException(
				"Unsupported JSDoc tag can't be deserialized: " + tag.String()))
		}
	}

	switch modMask {
	case 1:
		builder.RecordModifies(jsDocSerializer_modifiesArguments)
	case 2:
		builder.RecordModifies(jsDocSerializer_modifiesThis)
	case 3:
		builder.RecordModifies(jsDocSerializer_modifiesArgumentsAndThis)
	}
	if suppMask != 0 {
		var suppBuf [3]string
		suppSlice := suppBuf[:0]
		if (suppMask & 1) != 0 {
			suppSlice = append(suppSlice, "messageConventions")
		}
		if (suppMask & 2) != 0 {
			suppSlice = append(suppSlice, "partialAlias")
		}
		if (suppMask & 4) != 0 {
			suppSlice = append(suppSlice, "untranspilableFeatures")
		}
		builder.RecordSuppressionSliceWithDescription(suppSlice, "")
	}
}
