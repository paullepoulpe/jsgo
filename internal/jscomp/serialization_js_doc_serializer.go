// Port of com/google/javascript/jscomp/serialization/JSDocSerializer.java @ closure-compiler v20261015

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
	stringPool := StringPool_Builder_()
	return JSDocSerializer_DeserializeJsdoc(JSDocSerializer_SerializeJsdoc(jsdoc, stringPool), stringPool.Build())
}

// JSDocSerializer_SerializeJsdoc may return nil.
//
// Java: JSDocSerializer.serializeJsdoc(JSDocInfo, StringPool.Builder)
func JSDocSerializer_SerializeJsdoc(jsdoc *rhino.JSDocInfo, stringPool *StringPool_Builder) *OptimizationJsdoc {
	if jsdoc == nil {
		return nil
	}

	builder := OptimizationJsdoc_NewBuilder()

	if jsdoc.HasFileOverview() {
		builder.AddKind(JsdocTag_JSDOC_FILEOVERVIEW)
	}

	if jsdoc.GetLicense() != nil {
		builder.SetLicenseTextPointer(int32(stringPool.Put(*jsdoc.GetLicense())))
	}

	if jsdoc.IsSassGeneratedCssTs() {
		builder.AddKind(JsdocTag_JSDOC_SASS_GENERATED_CSS_TS)
	}

	// Used by CoverageInstrumentationCallback
	if jsdoc.IsNoCoverage() {
		builder.AddKind(JsdocTag_JSDOC_NO_COVERAGE)
	}

	if jsdoc.IsNoInline() {
		builder.AddKind(JsdocTag_JSDOC_NO_INLINE)
	}

	if jsdoc.IsEncourageInlining() {
		builder.AddKind(JsdocTag_JSDOC_ENCOURAGE_INLINING)
	}

	if jsdoc.IsRequireInlining() {
		builder.AddKind(JsdocTag_JSDOC_REQUIRE_INLINING)
	}

	if jsdoc.IsNoCollapse() {
		builder.AddKind(JsdocTag_JSDOC_NO_COLLAPSE)
	}

	if jsdoc.IsProvideGoog() {
		builder.AddKind(JsdocTag_JSDOC_PROVIDE_GOOG)
	}

	if jsdoc.IsProvideAlreadyProvided() {
		builder.AddKind(JsdocTag_JSDOC_PROVIDE_ALREADY_PROVIDED)
	}

	if jsdoc.IsTypeSummary() {
		builder.AddKind(JsdocTag_JSDOC_TYPE_SUMMARY_FILE)
	}

	if jsdoc.IsPureOrBreakMyCode() {
		builder.AddKind(JsdocTag_JSDOC_PURE_OR_BREAK_MY_CODE)
	}
	if jsdoc.IsCollapsibleOrBreakMyCode() {
		builder.AddKind(JsdocTag_JSDOC_COLLAPSIBLE_OR_BREAK_MY_CODE)
	}
	if jsdoc.HasThisType() {
		builder.AddKind(JsdocTag_JSDOC_THIS)
	}
	if jsdoc.HasEnumParameterType() {
		builder.AddKind(JsdocTag_JSDOC_ENUM)
	}
	if jsdoc.IsDefine() {
		builder.AddKind(JsdocTag_JSDOC_DEFINE)
	}
	if jsdoc.HasConstAnnotation() {
		builder.AddKind(JsdocTag_JSDOC_CONST)
	}
	if jsdoc.IsAnyIdGenerator() {
		builder.AddKind(jsDocSerializer_serializeIdGenerator(jsdoc))
	}

	// Used by PureFunctionIdentifier
	if jsdoc.IsNoSideEffects() {
		builder.AddKind(JsdocTag_JSDOC_NO_SIDE_EFFECTS)
	}
	if jsdoc.HasModifies() {
		for modifies := range jsdoc.GetModifies().All() {
			switch modifies {
			case "this":
				builder.AddKind(JsdocTag_JSDOC_MODIFIES_THIS)
			case "arguments":
				builder.AddKind(JsdocTag_JSDOC_MODIFIES_ARGUMENTS)
			default:
				// Currently, anything other than "this" is considered a modification to arguments
				builder.AddKind(JsdocTag_JSDOC_MODIFIES_ARGUMENTS)
			}
		}
	}
	if len(jsdoc.GetThrowsAnnotations()) != 0 {
		builder.AddKind(JsdocTag_JSDOC_THROWS)
	}

	// Used by DevirtualizeMethods and CollapseProperties
	if jsdoc.IsConstructor() {
		builder.AddKind(JsdocTag_JSDOC_CONSTRUCTOR)
	}
	if jsdoc.IsInterface() {
		builder.AddKind(JsdocTag_JSDOC_INTERFACE)
	}
	if jsdoc.GetSuppressions().Contains("partialAlias") {
		builder.AddKind(JsdocTag_JSDOC_SUPPRESS_PARTIAL_ALIAS)
	}

	// Used by ClosureCodeRemoval
	if jsdoc.IsAbstract() {
		builder.AddKind(JsdocTag_JSDOC_ABSTRACT)
	}

	// Used by ReplaceMessages
	if jsdoc.GetDescription() != nil {
		builder.SetDescriptionPointer(int32(stringPool.Put(*jsdoc.GetDescription())))
	}
	if jsdoc.GetAlternateMessageId() != nil {
		builder.SetAlternateMessageIdPointer(int32(stringPool.Put(*jsdoc.GetAlternateMessageId())))
	}
	if jsdoc.GetMeaning() != nil {
		builder.SetMeaningPointer(int32(stringPool.Put(*jsdoc.GetMeaning())))
	}
	if jsdoc.GetSuppressions().Contains("messageConventions") {
		builder.AddKind(JsdocTag_JSDOC_SUPPRESS_MESSAGE_CONVENTION)
	}
	if jsdoc.GetSuppressions().Contains("untranspilableFeatures") {
		builder.AddKind(JsdocTag_JSDOC_SUPPRESS_UNTRANSPILABLE_FEATURES)
	}
	if jsdoc.IsUsedViaDotConstructor() {
		builder.AddKind(JsdocTag_JSDOC_USED_VIA_DOT_CONSTRUCTOR)
	}
	if jsdoc.GetPerFileClosureUnawareMode() != nil {
		builder.SetClosureUnawareModePointer(
			int32(stringPool.Put(jsdoc.GetPerFileClosureUnawareMode().String())))
	}

	result := builder.Build()
	if OptimizationJsdoc_GetDefaultInstance().Equals(result) {
		return nil
	}
	return jutil.CheckNotNull(result)
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

// JSDocSerializer_DeserializeJsdoc may return nil.
//
// Java: JSDocSerializer.deserializeJsdoc(OptimizationJsdoc, StringPool)
func JSDocSerializer_DeserializeJsdoc(serializedJsdoc *OptimizationJsdoc, stringPool *StringPool) *rhino.JSDocInfo {
	if serializedJsdoc == nil {
		return nil
	}
	builder := rhino.JSDocInfo_Builder_()

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

	// lazily create these sets to save a few hundred ms for some large projects
	var modifies *jutil.TreeSet[string]
	var suppressions *jutil.TreeSet[string]

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
			if modifies == nil {
				modifies = jsDocSerializer_newStringTreeSet()
			}
			modifies.Add("this")
			continue
		case JsdocTag_JSDOC_MODIFIES_ARGUMENTS:
			if modifies == nil {
				modifies = jsDocSerializer_newStringTreeSet()
			}
			modifies.Add("arguments")
			continue
		case JsdocTag_JSDOC_THROWS:
			builder.RecordThrowsAnnotation("{!" + jsDocSerializer_PLACEHOLDER_TYPE_NAME + "}")
			continue
		case JsdocTag_JSDOC_CONSTRUCTOR:
			builder.RecordConstructor()
			continue
		case JsdocTag_JSDOC_INTERFACE:
			builder.RecordInterface()
			continue
		case JsdocTag_JSDOC_SUPPRESS_PARTIAL_ALIAS:
			if suppressions == nil {
				suppressions = jsDocSerializer_newStringTreeSet()
			}
			suppressions.Add("partialAlias")
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
			if suppressions == nil {
				suppressions = jsDocSerializer_newStringTreeSet()
			}
			suppressions.Add("messageConventions")
			continue
			// ReportUntranspilableFeatures pass will run in stage2 since it uses languageOut
			// information. It reports diagnostic {@code UNTRANSPILABLE_FEATURE_PRESENT} that can be
			// supppressed using `untranspilableFeatures` suppression tag.
			// Hence we must propagate it.
		case JsdocTag_JSDOC_SUPPRESS_UNTRANSPILABLE_FEATURES:
			if suppressions == nil {
				suppressions = jsDocSerializer_newStringTreeSet()
			}
			suppressions.Add("untranspilableFeatures")
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
	if modifies != nil {
		builder.RecordModifies(jutil.NewLinkedHashSet(modifies.ToSlice()...))
	}
	if suppressions != nil {
		builder.RecordSuppressions(jutil.NewLinkedHashSet(suppressions.ToSlice()...))
	}

	return builder.Build()
}
