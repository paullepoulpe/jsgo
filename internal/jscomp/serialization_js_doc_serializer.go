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
	if jsdoc == nil {
		return nil
	}

	var builder rhino.JSDocInfo_Builder
	var init bool
	ensureBuilder := func() *rhino.JSDocInfo_Builder {
		init = true
		return &builder
	}

	if mode := jsdoc.GetPerFileClosureUnawareMode(); mode != nil {
		ensureBuilder().RecordClosureUnawareCodeWithMode(*mode)
	}
	if license := jsdoc.GetLicense(); license != nil && *license != "" {
		ensureBuilder().AddLicense(rhino.RhinoStringPool_AddOrGet(*license))
	}
	if meaning := jsdoc.GetMeaning(); meaning != nil && *meaning != "" {
		ensureBuilder().RecordMeaning(rhino.RhinoStringPool_AddOrGet(*meaning))
	}
	if description := jsdoc.GetDescription(); description != nil && *description != "" {
		ensureBuilder().RecordDescription(rhino.RhinoStringPool_AddOrGet(*description))
	}
	if alternateMessageId := jsdoc.GetAlternateMessageId(); alternateMessageId != nil && *alternateMessageId != "" {
		ensureBuilder().RecordAlternateMessageId(rhino.RhinoStringPool_AddOrGet(*alternateMessageId))
	}

	if jsdoc.HasFileOverview() {
		ensureBuilder().RecordFileOverview("")
	}
	if jsdoc.IsSassGeneratedCssTs() {
		ensureBuilder().RecordSassGeneratedCssTs()
	}
	if jsdoc.IsNoCoverage() {
		ensureBuilder().RecordNoCoverage()
	}
	if jsdoc.IsNoInline() {
		ensureBuilder().RecordNoInline()
	}
	if jsdoc.IsEncourageInlining() {
		ensureBuilder().RecordEncourageInlining()
	}
	if jsdoc.IsRequireInlining() {
		ensureBuilder().RecordRequireInlining()
	}
	if jsdoc.IsNoCollapse() {
		ensureBuilder().RecordNoCollapse()
	}
	if jsdoc.IsProvideGoog() {
		ensureBuilder().RecordProvideGoog()
	}
	if jsdoc.IsProvideAlreadyProvided() {
		ensureBuilder().RecordProvideAlreadyProvided()
	}
	if jsdoc.IsTypeSummary() {
		ensureBuilder().RecordTypeSummary()
	}
	if jsdoc.IsPureOrBreakMyCode() {
		ensureBuilder().RecordPureOrBreakMyCode()
	}
	if jsdoc.IsCollapsibleOrBreakMyCode() {
		ensureBuilder().RecordCollapsibleOrBreakMyCode()
	}
	if jsdoc.HasThisType() {
		ensureBuilder().RecordThisType(jsDocSerializer_placeholderType)
	}
	if jsdoc.HasEnumParameterType() {
		ensureBuilder().RecordEnumParameterType(jsDocSerializer_placeholderType)
	}
	if jsdoc.IsDefine() {
		ensureBuilder().RecordDefineType(jsDocSerializer_placeholderType)
	}
	if jsdoc.HasConstAnnotation() {
		ensureBuilder().RecordConstancy()
	}
	if jsdoc.IsAnyIdGenerator() {
		if jsdoc.IsConsistentIdGenerator() {
			ensureBuilder().RecordConsistentIdGenerator()
		} else if jsdoc.IsStableIdGenerator() {
			ensureBuilder().RecordStableIdGenerator()
		} else if jsdoc.IsXidGenerator() {
			ensureBuilder().RecordXidGenerator()
		} else if jsdoc.IsMappedIdGenerator() {
			ensureBuilder().RecordMappedIdGenerator()
		} else if jsdoc.IsIdGenerator() {
			ensureBuilder().RecordIdGenerator()
		} else {
			panic(jutil.NewIllegalStateExceptionMsg("Failed to identify idGenerator inside JSDoc: " + jsdoc.String()))
		}
	}
	if jsdoc.IsNoSideEffects() {
		ensureBuilder().RecordNoSideEffects()
	}

	var hasModifiesArguments, hasModifiesThis bool
	if jsdoc.HasModifies() {
		for modifies := range jsdoc.GetModifies().All() {
			if modifies == "this" {
				hasModifiesThis = true
			} else {
				hasModifiesArguments = true
			}
		}
	}

	if len(jsdoc.GetThrowsAnnotations()) != 0 {
		ensureBuilder().RecordThrowsAnnotation("{!" + jsDocSerializer_PLACEHOLDER_TYPE_NAME + "}")
	}
	if jsdoc.IsConstructor() {
		ensureBuilder().RecordConstructor()
	}
	if jsdoc.IsInterface() {
		ensureBuilder().RecordInterface()
	}

	var hasSuppressMessageConvention, hasSuppressPartialAlias, hasSuppressUntranspilableFeatures bool
	jsdoc.ForEachSuppression(func(s string) bool {
		switch s {
		case "messageConventions":
			hasSuppressMessageConvention = true
		case "partialAlias":
			hasSuppressPartialAlias = true
		case "untranspilableFeatures":
			hasSuppressUntranspilableFeatures = true
		}
		return true
	})

	if jsdoc.IsAbstract() {
		ensureBuilder().RecordAbstract()
	}
	if jsdoc.IsUsedViaDotConstructor() {
		_ = ensureBuilder().RecordUsedViaDotConstructor()
	}

	if hasModifiesArguments || hasModifiesThis {
		modSet := jutil.NewLinkedHashSet[string]()
		if hasModifiesArguments {
			modSet.Add("arguments")
		}
		if hasModifiesThis {
			modSet.Add("this")
		}
		ensureBuilder().RecordModifies(modSet)
	}
	if hasSuppressMessageConvention || hasSuppressPartialAlias || hasSuppressUntranspilableFeatures {
		suppSet := jutil.NewLinkedHashSet[string]()
		if hasSuppressMessageConvention {
			suppSet.Add("messageConventions")
		}
		if hasSuppressPartialAlias {
			suppSet.Add("partialAlias")
		}
		if hasSuppressUntranspilableFeatures {
			suppSet.Add("untranspilableFeatures")
		}
		ensureBuilder().RecordSuppressionsWithDescription(suppSet, "")
	}

	if !init {
		return nil
	}
	return builder.Build()
}

// JSDocSerializer_SerializeJsdoc may return nil.
//
// Java: JSDocSerializer.serializeJsdoc(JSDocInfo, StringPool.Builder)
func JSDocSerializer_SerializeJsdoc(jsdoc *rhino.JSDocInfo, stringPool *StringPool_Builder) *OptimizationJsdoc {
	var m OptimizationJsdoc
	var kindBuf [8]int32
	if !jsDocSerializer_fillJsdocProto(jsdoc, stringPool, &m, kindBuf[:0]) {
		return nil
	}
	m.kind_ = append([]int32(nil), m.kind_...)
	m.frozen = true
	m.memoizedSize = m.scalarsSerializedSize()
	return &m
}

func jsDocSerializer_fillJsdocTagsOnly(jsdoc *rhino.JSDocInfo, kindBuf []int32) []int32 {
	addKind := func(v JsdocTag) {
		kindBuf = append(kindBuf, v.GetNumber())
	}

	if jsdoc.HasFileOverview() {
		addKind(JsdocTag_JSDOC_FILEOVERVIEW)
	}
	if jsdoc.IsSassGeneratedCssTs() {
		addKind(JsdocTag_JSDOC_SASS_GENERATED_CSS_TS)
	}
	// Used by CoverageInstrumentationCallback
	if jsdoc.IsNoCoverage() {
		addKind(JsdocTag_JSDOC_NO_COVERAGE)
	}
	if jsdoc.IsNoInline() {
		addKind(JsdocTag_JSDOC_NO_INLINE)
	}
	if jsdoc.IsEncourageInlining() {
		addKind(JsdocTag_JSDOC_ENCOURAGE_INLINING)
	}
	if jsdoc.IsRequireInlining() {
		addKind(JsdocTag_JSDOC_REQUIRE_INLINING)
	}
	if jsdoc.IsNoCollapse() {
		addKind(JsdocTag_JSDOC_NO_COLLAPSE)
	}
	if jsdoc.IsProvideGoog() {
		addKind(JsdocTag_JSDOC_PROVIDE_GOOG)
	}
	if jsdoc.IsProvideAlreadyProvided() {
		addKind(JsdocTag_JSDOC_PROVIDE_ALREADY_PROVIDED)
	}
	if jsdoc.IsTypeSummary() {
		addKind(JsdocTag_JSDOC_TYPE_SUMMARY_FILE)
	}
	if jsdoc.IsPureOrBreakMyCode() {
		addKind(JsdocTag_JSDOC_PURE_OR_BREAK_MY_CODE)
	}
	if jsdoc.IsCollapsibleOrBreakMyCode() {
		addKind(JsdocTag_JSDOC_COLLAPSIBLE_OR_BREAK_MY_CODE)
	}
	if jsdoc.HasThisType() {
		addKind(JsdocTag_JSDOC_THIS)
	}
	if jsdoc.HasEnumParameterType() {
		addKind(JsdocTag_JSDOC_ENUM)
	}
	if jsdoc.IsDefine() {
		addKind(JsdocTag_JSDOC_DEFINE)
	}
	if jsdoc.HasConstAnnotation() {
		addKind(JsdocTag_JSDOC_CONST)
	}
	if jsdoc.IsAnyIdGenerator() {
		addKind(jsDocSerializer_serializeIdGenerator(jsdoc))
	}

	// Used by PureFunctionIdentifier
	if jsdoc.IsNoSideEffects() {
		addKind(JsdocTag_JSDOC_NO_SIDE_EFFECTS)
	}
	if jsdoc.HasModifies() {
		for modifies := range jsdoc.GetModifies().All() {
			switch modifies {
			case "this":
				addKind(JsdocTag_JSDOC_MODIFIES_THIS)
			case "arguments":
				addKind(JsdocTag_JSDOC_MODIFIES_ARGUMENTS)
			default:
				// Currently, anything other than "this" is considered a modification to arguments
				addKind(JsdocTag_JSDOC_MODIFIES_ARGUMENTS)
			}
		}
	}
	if len(jsdoc.GetThrowsAnnotations()) != 0 {
		addKind(JsdocTag_JSDOC_THROWS)
	}

	// Used by DevirtualizeMethods and CollapseProperties
	if jsdoc.IsConstructor() {
		addKind(JsdocTag_JSDOC_CONSTRUCTOR)
	}
	if jsdoc.IsInterface() {
		addKind(JsdocTag_JSDOC_INTERFACE)
	}
	if jsdoc.HasSuppression("partialAlias") {
		addKind(JsdocTag_JSDOC_SUPPRESS_PARTIAL_ALIAS)
	}

	// Used by ClosureCodeRemoval
	if jsdoc.IsAbstract() {
		addKind(JsdocTag_JSDOC_ABSTRACT)
	}

	// Used by ReplaceMessages
	if jsdoc.HasSuppression("messageConventions") {
		addKind(JsdocTag_JSDOC_SUPPRESS_MESSAGE_CONVENTION)
	}
	if jsdoc.HasSuppression("untranspilableFeatures") {
		addKind(JsdocTag_JSDOC_SUPPRESS_UNTRANSPILABLE_FEATURES)
	}
	if jsdoc.IsUsedViaDotConstructor() {
		addKind(JsdocTag_JSDOC_USED_VIA_DOT_CONSTRUCTOR)
	}
	return kindBuf
}

func jsDocSerializer_fillJsdocProto(jsdoc *rhino.JSDocInfo, stringPool *StringPool_Builder, m *OptimizationJsdoc, kindBuf []int32) bool {
	if jsdoc == nil {
		return false
	}

	m.kind_ = jsDocSerializer_fillJsdocTagsOnly(jsdoc, kindBuf)
	if jsdoc.GetLicense() != nil {
		m.licenseTextPointer_ = int32(stringPool.Put(*jsdoc.GetLicense()))
	}
	if jsdoc.GetDescription() != nil {
		m.descriptionPointer_ = int32(stringPool.Put(*jsdoc.GetDescription()))
	}
	if jsdoc.GetAlternateMessageId() != nil {
		m.alternateMessageIdPointer_ = int32(stringPool.Put(*jsdoc.GetAlternateMessageId()))
	}
	if jsdoc.GetMeaning() != nil {
		m.meaningPointer_ = int32(stringPool.Put(*jsdoc.GetMeaning()))
	}
	if jsdoc.GetPerFileClosureUnawareMode() != nil {
		m.closureUnawareModePointer_ =
			int32(stringPool.Put(jsdoc.GetPerFileClosureUnawareMode().String()))
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
	for _, m := range rhino.JSDocInfo_PerFileClosureUnawareMode_values() {
		if m.String() == name {
			return m
		}
	}
	panic(jutil.NewIllegalArgumentExceptionMsg(
		"No enum constant com.google.javascript.rhino.JSDocInfo.PerFileClosureUnawareMode." + name))
}

// jsDocSerializer_newStringTreeSet is new TreeSet<String>() (natural
// String order).
func jsDocSerializer_newStringTreeSet() *jutil.TreeSet[string] {
	return jutil.NewTreeSet(jutil.StringCompareTo)
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
	builder.Reset()
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
		builder.AddLicense(stringPool.Get(int(serializedJsdoc.GetLicenseTextPointer())))
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

	// Bitmasks avoid allocating TreeSets/LinkedHashSets on the deserialization hot path
	// while preserving exact alphabetical TreeSet<String> order:
	// modifies: bit 0 = "arguments", bit 1 = "this"
	// suppressions: bit 0 = "messageConventions", bit 1 = "partialAlias", bit 2 = "untranspilableFeatures"
	var modMask uint8
	var suppMask uint8

	for i := 0; i < serializedJsdoc.GetKindCount(); i++ {
		tag := serializedJsdoc.GetKind(i)
		switch tag {
		case JsdocTag_JSDOC_CONST:
			builder.RecordConstancy()
			continue
		case JsdocTag_JSDOC_ENUM:
			builder.RecordEnumParameterType(jsDocSerializer_placeholderType)
			continue
		case JsdocTag_JSDOC_THIS:
			builder.RecordThisType(jsDocSerializer_placeholderType)
			continue
		case JsdocTag_JSDOC_NO_COLLAPSE:
			builder.RecordNoCollapse()
			continue
		case JsdocTag_JSDOC_NO_INLINE:
			builder.RecordNoInline()
			continue
		case JsdocTag_JSDOC_REQUIRE_INLINING:
			builder.RecordRequireInlining()
			continue
		case JsdocTag_JSDOC_ENCOURAGE_INLINING:
			builder.RecordEncourageInlining()
			continue
		case JsdocTag_JSDOC_PROVIDE_GOOG:
			builder.RecordProvideGoog()
			continue
		case JsdocTag_JSDOC_PROVIDE_ALREADY_PROVIDED:
			builder.RecordProvideAlreadyProvided()
			continue
		case JsdocTag_JSDOC_TYPE_SUMMARY_FILE:
			builder.RecordTypeSummary()
			continue
		case JsdocTag_JSDOC_PURE_OR_BREAK_MY_CODE:
			builder.RecordPureOrBreakMyCode()
			continue
		case JsdocTag_JSDOC_COLLAPSIBLE_OR_BREAK_MY_CODE:
			builder.RecordCollapsibleOrBreakMyCode()
			continue
		case JsdocTag_JSDOC_DEFINE:
			builder.RecordDefineType(jsDocSerializer_placeholderType)
			continue
		case JsdocTag_JSDOC_NO_SIDE_EFFECTS:
			builder.RecordNoSideEffects()
			continue
		case JsdocTag_JSDOC_MODIFIES_THIS:
			modMask |= 2
			continue
		case JsdocTag_JSDOC_MODIFIES_ARGUMENTS:
			modMask |= 1
			continue
		case JsdocTag_JSDOC_THROWS:
			builder.RecordThrowsAnnotation(jsDocSerializer_throwsPlaceholder)
			continue
		case JsdocTag_JSDOC_CONSTRUCTOR:
			builder.RecordConstructor()
			continue
		case JsdocTag_JSDOC_INTERFACE:
			builder.RecordInterface()
			continue
		case JsdocTag_JSDOC_SUPPRESS_PARTIAL_ALIAS:
			suppMask |= 2
			continue
		case JsdocTag_JSDOC_ID_GENERATOR_CONSISTENT:
			builder.RecordConsistentIdGenerator()
			continue
		case JsdocTag_JSDOC_ID_GENERATOR_STABLE:
			builder.RecordStableIdGenerator()
			continue
		case JsdocTag_JSDOC_ID_GENERATOR_MAPPED:
			builder.RecordMappedIdGenerator()
			continue
		case JsdocTag_JSDOC_ID_GENERATOR_XID:
			builder.RecordXidGenerator()
			continue
		case JsdocTag_JSDOC_ID_GENERATOR_INCONSISTENT:
			builder.RecordIdGenerator()
			continue
		case JsdocTag_JSDOC_ABSTRACT:
			builder.RecordAbstract()
			continue
			// TODO(lharker): stage 2 passes ideally shouldn't report diagnostics, so this could be
			// moved to stage 1.
		case JsdocTag_JSDOC_SUPPRESS_MESSAGE_CONVENTION:
			suppMask |= 1
			continue
			// ReportUntranspilableFeatures pass will run in stage2 since it uses languageOut
			// information. It reports diagnostic {@code UNTRANSPILABLE_FEATURE_PRESENT} that can be
			// supppressed using `untranspilableFeatures` suppression tag.
			// Hence we must propagate it.
		case JsdocTag_JSDOC_SUPPRESS_UNTRANSPILABLE_FEATURES:
			suppMask |= 4
			continue
		case JsdocTag_JSDOC_FILEOVERVIEW:
			builder.RecordFileOverview("")
			continue
		case JsdocTag_JSDOC_NO_COVERAGE:
			builder.RecordNoCoverage()
			continue
		case JsdocTag_JSDOC_SASS_GENERATED_CSS_TS:
			builder.RecordSassGeneratedCssTs()
			continue
		case JsdocTag_JSDOC_USED_VIA_DOT_CONSTRUCTOR:
			_ = builder.RecordUsedViaDotConstructor()
			continue
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

	return builder.Build()
}
