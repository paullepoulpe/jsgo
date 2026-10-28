// Port of com/google/javascript/rhino/JSDocInfo.java @ closure-compiler v20261021

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

package rhino

import (
	"math/bits"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JSDocInfo_Builder is a builder for JSDocInfo objects. This builder is
// required because JSDocInfo instances have immutable structure. It provides
// early incompatibility detection among properties stored on the JSDocInfo
// object being created.
type JSDocInfo_Builder struct {
	propMask   int64
	entries    []any
	entriesBuf [4]any
	bits       int64

	// whether the current JSDocInfo has valuable information
	populated bool
	// whether a current marker is active at the tail of JSDocInfo_MARKERS
	hasCurrentMarker bool
	// the unique license texts added so far
	licenseTexts []string
}

// NewJSDocInfo_Builder is the implicit Builder constructor.
//
// Java: JSDocInfo.Builder()
func NewJSDocInfo_Builder() *JSDocInfo_Builder {
	b := &JSDocInfo_Builder{}
	b.entries = b.entriesBuf[:0]
	return b
}

// Reset resets the builder for reuse without allocating.
func (b *JSDocInfo_Builder) Reset() {
	clear(b.entries)
	b.entries = b.entriesBuf[:0]
	b.propMask = 0
	b.bits = 0
	b.populated = false
	b.hasCurrentMarker = false
	clear(b.licenseTexts)
	b.licenseTexts = b.licenseTexts[:0]
}

// Java: JSDocInfo.Builder.copyFrom(JSDocInfo)
func JSDocInfo_Builder_CopyFrom(info *JSDocInfo) *JSDocInfo_Builder {
	return info.ToBuilder()
}

// Java: JSDocInfo.Builder.maybeCopyFrom(JSDocInfo)
func JSDocInfo_Builder_MaybeCopyFrom(info *JSDocInfo) *JSDocInfo_Builder {
	if info != nil {
		return info.ToBuilder()
	}
	return JSDocInfo_Builder_().ParseDocumentation()
}

// JSDocInfo_Builder_MaybeCopyFromWithNewType returns a JSDocInfo.Builder
// that contains a copy of the given JSDocInfo in which only the `@type`
// field of the JSDocInfo is replaced with the given typeExpression. This is
// done to prevent generating code in the client module which references
// local variables from another module.
//
// Java: JSDocInfo.Builder.maybeCopyFromWithNewType(JSDocInfo, JSTypeExpression)
func JSDocInfo_Builder_MaybeCopyFromWithNewType(info *JSDocInfo, typeExpression *JSTypeExpression) *JSDocInfo_Builder {
	if info == nil {
		return JSDocInfo_Builder_().ParseDocumentation().SetType(typeExpression)
	}
	return info.ToBuilder().SetType(typeExpression)
}

// Java: JSDocInfo.Builder.copyFromWithNewType(JSDocInfo, JSTypeExpression)
func JSDocInfo_Builder_CopyFromWithNewType(info *JSDocInfo, typeExpression *JSTypeExpression) *JSDocInfo_Builder {
	return info.ToBuilder().SetType(typeExpression)
}

// JSDocInfo_Builder_MaybeCopyFromAndReplaceNames returns a JSDocInfo.Builder
// that contains a JSDoc in which all module local types (which may be inside
// `@param`, `@type` or `@returns` are replaced with unknown. This is done to
// prevent generating code in the client module which references local
// variables from another module.
//
// Java: JSDocInfo.Builder.maybeCopyFromAndReplaceNames(JSDocInfo, Set<String>)
func JSDocInfo_Builder_MaybeCopyFromAndReplaceNames(info *JSDocInfo, moduleLocalNamesToReplace *jutil.LinkedHashSet[string]) *JSDocInfo_Builder {
	if info != nil {
		return JSDocInfo_Builder_CopyFromAndReplaceNames(info, moduleLocalNamesToReplace)
	}
	return JSDocInfo_Builder_().ParseDocumentation()
}

// Java: JSDocInfo.Builder.copyFromAndReplaceNames(JSDocInfo, Set<String>)
func JSDocInfo_Builder_CopyFromAndReplaceNames(info *JSDocInfo, oldNames *jutil.LinkedHashSet[string]) *JSDocInfo_Builder {
	return info.CloneAndReplaceTypeNames(oldNames).ToBuilder() // TODO - populated
}

// ParseDocumentation configures the builder to parse documentation. This
// should be called immediately after instantiating the builder if
// documentation should be included, since it enables various operations to
// do work that would otherwise be no-ops.
//
// Java: JSDocInfo.Builder.parseDocumentation()
func (b *JSDocInfo_Builder) ParseDocumentation() *JSDocInfo_Builder {
	b.SetBit(JSDocInfo_Bit_INCLUDE_DOCUMENTATION, true)
	return b
}

// Java: JSDocInfo.Builder.shouldParseDocumentation()
func (b *JSDocInfo_Builder) ShouldParseDocumentation() bool {
	return b.CheckBit(JSDocInfo_Bit_INCLUDE_DOCUMENTATION)
}

// RecordOriginalCommentString sets the original JSDoc comment string. This
// is a no-op if the builder isn't configured to record documentation.
//
// Java: JSDocInfo.Builder.recordOriginalCommentString(String)
func (b *JSDocInfo_Builder) RecordOriginalCommentString(sourceComment string) {
	if b.ShouldParseDocumentation() {
		b.populated = true
		b.SetProp(JSDocInfo_SOURCE_COMMENT, sourceComment)
	}
}

// RecordOriginalCommentPosition sets the position of original JSDoc
// comment.
//
// Java: JSDocInfo.Builder.recordOriginalCommentPosition(int)
func (b *JSDocInfo_Builder) RecordOriginalCommentPosition(position int) {
	if b.ShouldParseDocumentation() {
		b.populated = true
		b.SetProp(JSDocInfo_ORIGINAL_COMMENT_POSITION, position)
	}
}

// IsPopulatedWithFileOverview returns whether this builder is populated with
// information that can be used to build a JSDocInfo object that has a
// fileoverview tag.
//
// Java: JSDocInfo.Builder.isPopulatedWithFileOverview()
func (b *JSDocInfo_Builder) IsPopulatedWithFileOverview() bool {
	return b.populated &&
		((b.bits&
			(JSDocInfo_Bit_FILEOVERVIEW.Mask()|
				JSDocInfo_Bit_EXTERNS.Mask()|
				JSDocInfo_Bit_NOCOMPILE.Mask()|
				JSDocInfo_Bit_NOCOVERAGE.Mask()|
				JSDocInfo_Bit_TYPE_SUMMARY.Mask()|
				JSDocInfo_Bit_ENHANCED_NAMESPACE.Mask())) != 0 ||
			b.IsModsRecorded())
}

// IsDescriptionRecorded returns whether this builder recorded a
// description.
//
// Java: JSDocInfo.Builder.isDescriptionRecorded()
func (b *JSDocInfo_Builder) IsDescriptionRecorded() bool {
	return b.GetProp(JSDocInfo_DESCRIPTION) != nil
}

// IsNoInline returns whether the `@noinline` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.Builder.isNoInline()
func (b *JSDocInfo_Builder) IsNoInline() bool {
	return b.CheckBit(JSDocInfo_Bit_NOINLINE)
}

// IsRequireInlining returns whether the `@requireInlining` annotation is
// present on this JSDocInfo.
//
// Java: JSDocInfo.Builder.isRequireInlining()
func (b *JSDocInfo_Builder) IsRequireInlining() bool {
	return b.CheckBit(JSDocInfo_Bit_REQUIRE_INLINING)
}

// IsEncourageInlining returns whether the `@encourageInlining` annotation is
// present on this JSDocInfo.
//
// Java: JSDocInfo.Builder.isEncourageInlining()
func (b *JSDocInfo_Builder) IsEncourageInlining() bool {
	return b.CheckBit(JSDocInfo_Bit_ENCOURAGE_INLINING)
}

// Build builds a JSDocInfo object based on the populated information and
// returns it. If no value was populated, this method simply returns null.
//
// Java: JSDocInfo.Builder.build()
func (b *JSDocInfo_Builder) Build() *JSDocInfo {
	return b.BuildWithAlways( /* always= */ false)
}

// BuildWithAlways builds a JSDocInfo object based on the populated
// information and returns it. If no value was populated and `always` is
// false, returns null. If `always` is true, returns a default JSDocInfo.
//
// Java: JSDocInfo.Builder.build(boolean)
func (b *JSDocInfo_Builder) BuildWithAlways(always bool) *JSDocInfo {
	if b.populated || always {
		info := newJSDocInfoFromBuilder(b.bits, b.propMask, b.entries)
		b.populated = false
		return info
	}
	return nil
}

// BuildAndReset builds a JSDocInfo object based on the populated information
// and returns it. Once this method is called, the builder can be reused to
// build another JSDocInfo object. If no value was populated, this method
// simply returns null.
//
// Java: JSDocInfo.Builder.buildAndReset()
func (b *JSDocInfo_Builder) BuildAndReset() *JSDocInfo {
	info := b.Build()
	b.bits &= JSDocInfo_Bit_INCLUDE_DOCUMENTATION.Mask() // only keep this one flag
	clear(b.entries)
	b.entries = b.entries[:0]
	b.propMask = 0
	b.populated = false
	b.hasCurrentMarker = false
	clear(b.licenseTexts)
	b.licenseTexts = b.licenseTexts[:0]
	return info
}

func (b *JSDocInfo_Builder) currentMarkerPtr() *JSDocInfo_Marker {
	if !b.hasCurrentMarker {
		return nil
	}
	markers := jsdocGet[[]JSDocInfo_Marker](b.GetProp(JSDocInfo_MARKERS))
	if len(markers) == 0 {
		return nil
	}
	return &markers[len(markers)-1]
}

// MarkAnnotation adds a marker to the current JSDocInfo and populates the
// marker with the annotation information.
//
// Java: JSDocInfo.Builder.markAnnotation(String, int, int)
func (b *JSDocInfo_Builder) MarkAnnotation(annotation string, lineno int, charno int) {
	marker := b.AddMarker()
	if marker != nil {
		SetTrimmedItem(&marker.annotation, annotation)
		marker.annotation.SetPositionInformation(lineno, charno, lineno, charno+jutil.StringLength(annotation))
		b.populated = true
		b.hasCurrentMarker = true
	} else {
		b.hasCurrentMarker = false
	}
}

// Java: JSDocInfo.Builder.addMarker()
func (b *JSDocInfo_Builder) AddMarker() *JSDocInfo_Marker {
	if b.ShouldParseDocumentation() {
		markers := jsdocGet[[]JSDocInfo_Marker](b.GetProp(JSDocInfo_MARKERS))
		markers = append(markers, JSDocInfo_Marker{})
		b.SetProp(JSDocInfo_MARKERS, markers)
		return &markers[len(markers)-1]
	}
	return nil
}

// MarkText adds a textual block to the current marker.
//
// Java: JSDocInfo.Builder.markText(String, int, int, int, int)
func (b *JSDocInfo_Builder) MarkText(text string, startLineno int, startCharno int, endLineno int, endCharno int) {
	if marker := b.currentMarkerPtr(); marker != nil {
		marker.description.SetItem(text)
		marker.description.SetPositionInformation(startLineno, startCharno, endLineno, endCharno)
	}
}

// MarkTypeNode adds a type declaration to the current marker.
//
// Java: JSDocInfo.Builder.markTypeNode(Node, int, int, int, int, boolean)
func (b *JSDocInfo_Builder) MarkTypeNode(typeNode Node, lineno int, startCharno int, endLineno int, endCharno int, hasLC bool) {
	if marker := b.currentMarkerPtr(); marker != nil {
		marker.type_.SetItem(typeNode)
		marker.type_.SetHasBrackets(hasLC)
		marker.type_.SetPositionInformation(lineno, startCharno, endLineno, endCharno)
	}
}

// MarkName adds a name declaration to the current marker.
//
// Java: JSDocInfo.Builder.markName(String, Node, int, int)
func (b *JSDocInfo_Builder) MarkName(name string, templateNode Node, lineno int, charno int) {
	if marker := b.currentMarkerPtr(); marker != nil {
		checkTrimmedStringItem(name)
		nameLen := jutil.StringLength(name)
		node := Node_NewStringWithToken(Token_NAME, name).SetLinenoCharno(lineno, charno)
		node.SetLength(nameLen)
		if !templateNode.IsNil() {
			node.SetStaticSourceFileFrom(templateNode)
		}
		marker.nameNode.SetItem(node)
		marker.nameNode.SetPositionInformation(lineno, charno, lineno, charno+nameLen)
	}
}

// RecordBlockDescription records a block-level description; returns true
// if the description was recorded.
//
// Java: JSDocInfo.Builder.recordBlockDescription(String)
func (b *JSDocInfo_Builder) RecordBlockDescription(description string) bool {
	b.populated = true
	if !b.ShouldParseDocumentation() {
		return true
	}
	return b.PopulateProp(JSDocInfo_BLOCK_DESCRIPTION, description)
}

// RecordVisibility records a visibility; returns true if the visibility was
// recorded and false if it was already defined.
//
// Java: JSDocInfo.Builder.recordVisibility(Visibility)
func (b *JSDocInfo_Builder) RecordVisibility(visibility JSDocInfo_Visibility) bool {
	if b.GetProp(JSDocInfo_VISIBILITY) == nil {
		b.populated = true
		b.SetProp(JSDocInfo_VISIBILITY, visibility)
		return true
	}
	return false
}

// Java: JSDocInfo.Builder.overwriteVisibility(Visibility)
func (b *JSDocInfo_Builder) OverwriteVisibility(visibility JSDocInfo_Visibility) {
	b.populated = true
	b.SetProp(JSDocInfo_VISIBILITY, visibility)
}

// RecordParameter records a typed parameter; returns true if the typed
// parameter was recorded and false if a parameter with the same name was
// already defined.
//
// Java: JSDocInfo.Builder.recordParameter(String, JSTypeExpression)
func (b *JSDocInfo_Builder) RecordParameter(parameterName string, type_ *JSTypeExpression) bool {
	return !b.HasAnySingletonTypeTags() &&
		b.populateParamEntry(JSDocInfo_PARAMETERS, RhinoStringPool_AddOrGet(parameterName), type_)
}

// RecordParameterDescription records a parameter's description; returns
// true if the parameter's description was recorded and false if a parameter
// with the same name was already defined.
//
// Java: JSDocInfo.Builder.recordParameterDescription(String, String)
func (b *JSDocInfo_Builder) RecordParameterDescription(parameterName string, description string) bool {
	if !b.ShouldParseDocumentation() {
		return true
	}
	return b.populateStringMapEntry(
		JSDocInfo_PARAMETER_DESCRIPTIONS, RhinoStringPool_AddOrGet(parameterName), description)
}

// RecordTemplateTypeName records a template type name; returns true if the
// template type name was recorded and false if the input template type name
// was already defined.
//
// Java: JSDocInfo.Builder.recordTemplateTypeName(String)
func (b *JSDocInfo_Builder) RecordTemplateTypeName(name string) bool {
	return b.RecordTemplateTypeNameWithBound(name, nil)
}

// Java: JSDocInfo.Builder.recordTemplateTypeName(String, JSTypeExpression)
func (b *JSDocInfo_Builder) RecordTemplateTypeNameWithBound(name string, bound *JSTypeExpression) bool {
	if bound == nil {
		bound = JSTypeExpression_IMPLICIT_TEMPLATE_BOUND
	}
	transformations := jsdocGet[*jutil.LinkedHashMap[string, Node]](b.GetProp(JSDocInfo_TYPE_TRANSFORMATIONS))
	if (transformations != nil && transformations.ContainsKey(name)) ||
		b.hasPropKey(JSDocInfo_TYPEDEF_TYPE) {
		return false
	}
	return b.populateParamEntry(JSDocInfo_TEMPLATE_TYPE_NAMES, RhinoStringPool_AddOrGet(name), bound)
}

// RecordTypeTransformation records a type transformation expression
// together with its template type name.
//
// Java: JSDocInfo.Builder.recordTypeTransformation(String, Node)
func (b *JSDocInfo_Builder) RecordTypeTransformation(name string, expr Node) bool {
	for _, p := range jsdocGet[[]JSDocParam](b.GetProp(JSDocInfo_TEMPLATE_TYPE_NAMES)) {
		if p.Name == name {
			return false
		}
	}
	return b.populateNodeMapEntry(JSDocInfo_TYPE_TRANSFORMATIONS, RhinoStringPool_AddOrGet(name), expr)
}

func (b *JSDocInfo_Builder) appendStringProp(prop *JSDocInfo_Property, s string) {
	list := jsdocGet[[]string](b.GetProp(prop))
	b.SetProp(prop, append(list, s))
}

// RecordThrowsAnnotation records a throw annotation description; returns
// true if the type's description was recorded. The description of a throw
// annotation is the text including the type.
//
// Java: JSDocInfo.Builder.recordThrowsAnnotation(String)
func (b *JSDocInfo_Builder) RecordThrowsAnnotation(annotation string) bool {
	b.populated = true
	if b.CheckBit(JSDocInfo_Bit_NOSIDEEFFECTS) {
		return false
	}
	// TODO(user): Does it make sense to check for singleton tags here?
	// Note that if the @throws annotation appears before a singleton tag like
	// @type, the throws annotation is preserved, but if it appears after the
	// singleton tag, it gets dropped.
	if !b.HasAnySingletonTypeTags() {
		throwsAnnotations := jsdocGet[[]string](b.GetProp(JSDocInfo_THROWS_ANNOTATIONS))
		if b.ShouldParseDocumentation() {
			b.SetProp(JSDocInfo_THROWS_ANNOTATIONS, append(throwsAnnotations, annotation))
		} else if len(throwsAnnotations) == 0 {
			// Add at least one annotation so that PureFunctionIdentifier knows
			// about the side effect.
			b.SetProp(JSDocInfo_THROWS_ANNOTATIONS, append(throwsAnnotations, ""))
		}
	}
	return true
}

// RecordAuthor adds an author to the current information.
//
// Java: JSDocInfo.Builder.recordAuthor(String)
func (b *JSDocInfo_Builder) RecordAuthor(author string) bool {
	b.populated = true
	if b.ShouldParseDocumentation() {
		b.appendStringProp(JSDocInfo_AUTHORS, author)
	}
	// NOTE: this could be removed, since it's always true.
	return true
}

// RecordReference adds a reference ("@see") to the current information.
//
// Java: JSDocInfo.Builder.recordReference(String)
func (b *JSDocInfo_Builder) RecordReference(reference string) bool {
	b.populated = true
	if b.ShouldParseDocumentation() {
		b.appendStringProp(JSDocInfo_SEES, reference)
	}
	// NOTE: this could be removed, since it's always true.
	return true
}

// RecordConsistentIdGenerator records that the JSDocInfo being built should
// have its JSDocInfo#isConsistentIdGenerator() flag set to true; returns
// true if the consistentIdGenerator flag was recorded and false if it was
// already recorded.
//
// Java: JSDocInfo.Builder.recordConsistentIdGenerator()
func (b *JSDocInfo_Builder) RecordConsistentIdGenerator() bool {
	return b.PopulateProp(JSDocInfo_ID_GENERATOR, JSDocInfo_IdGenerator_CONSISTENT)
}

// RecordStableIdGenerator records that the JSDocInfo being built should have
// its JSDocInfo#isStableIdGenerator() flag set to true; returns true if the
// stableIdGenerator flag was recorded and false if it was already recorded.
//
// Java: JSDocInfo.Builder.recordStableIdGenerator()
func (b *JSDocInfo_Builder) RecordStableIdGenerator() bool {
	return b.PopulateProp(JSDocInfo_ID_GENERATOR, JSDocInfo_IdGenerator_STABLE)
}

// RecordXidGenerator records that the JSDocInfo being built should have its
// JSDocInfo#isXidGenerator() flag set to true; returns true if the
// isXidGenerator flag was recorded and false if it was already recorded.
//
// Java: JSDocInfo.Builder.recordXidGenerator()
func (b *JSDocInfo_Builder) RecordXidGenerator() bool {
	return b.PopulateProp(JSDocInfo_ID_GENERATOR, JSDocInfo_IdGenerator_XID)
}

// RecordMappedIdGenerator records that the JSDocInfo being built should have
// its JSDocInfo#isStableIdGenerator() flag set to true; returns true if the
// stableIdGenerator flag was recorded and false if it was already recorded.
//
// Java: JSDocInfo.Builder.recordMappedIdGenerator()
func (b *JSDocInfo_Builder) RecordMappedIdGenerator() bool {
	return b.PopulateProp(JSDocInfo_ID_GENERATOR, JSDocInfo_IdGenerator_MAPPED)
}

// RecordIdGenerator records that the JSDocInfo being built should have its
// JSDocInfo#isIdGenerator() flag set to true; returns true if the
// idGenerator flag was recorded and false if it was already recorded.
//
// Java: JSDocInfo.Builder.recordIdGenerator()
func (b *JSDocInfo_Builder) RecordIdGenerator() bool {
	return b.PopulateProp(JSDocInfo_ID_GENERATOR, JSDocInfo_IdGenerator_UNIQUE)
}

// RecordDeprecationReason records the deprecation reason.
//
// Java: JSDocInfo.Builder.recordDeprecationReason(String)
func (b *JSDocInfo_Builder) RecordDeprecationReason(reason string) bool {
	return b.PopulateProp(JSDocInfo_DEPRECATION_REASON, reason)
}

// IsDeprecationReasonRecorded returns whether a deprecation reason has been
// recorded.
//
// Java: JSDocInfo.Builder.isDeprecationReasonRecorded()
func (b *JSDocInfo_Builder) IsDeprecationReasonRecorded() bool {
	return b.GetProp(JSDocInfo_DEPRECATION_REASON) != nil
}

func (b *JSDocInfo_Builder) recordSuppressionWarnings(warnings []string, description string) {
	b.populated = true
	current := jsdocGet[[]JSDocSuppression](b.GetProp(JSDocInfo_SUPPRESSIONS))
	for _, existing := range current {
		if sameWarningSet(existing.Warnings, warnings) {
			// Exact @suppress warning exists already. Return without recording.
			return
		}
	}
	b.SetProp(JSDocInfo_SUPPRESSIONS, append(current, JSDocSuppression{
		Warnings:    warnings,
		Description: description,
	}))
}

// RecordSuppressionSliceWithDescription records a parsed @suppress slice.
func (b *JSDocInfo_Builder) RecordSuppressionSliceWithDescription(rawNames []string, description string) {
	warnings := make([]string, 0, len(rawNames))
	for _, raw := range rawNames {
		s := RhinoStringPool_AddOrGet(raw)
		if !slices.Contains(warnings, s) {
			warnings = append(warnings, s)
		}
	}
	b.recordSuppressionWarnings(warnings, description)
}

// Java: JSDocInfo.Builder.recordSuppressions(ImmutableSet<String>, String)
func (b *JSDocInfo_Builder) RecordSuppressionsWithDescription(suppressions *jutil.LinkedHashSet[string], description string) {
	warnings := make([]string, 0, suppressions.Size())
	for raw := range suppressions.All() {
		s := RhinoStringPool_AddOrGet(raw)
		if !slices.Contains(warnings, s) {
			warnings = append(warnings, s)
		}
	}
	b.recordSuppressionWarnings(warnings, description)
}

// RecordSuppressions records the list of suppressed warnings, possibly
// adding to the set of already configured warnings.
//
// Java: JSDocInfo.Builder.recordSuppressions(Set<String>)
func (b *JSDocInfo_Builder) RecordSuppressions(suppressions *jutil.LinkedHashSet[string]) {
	jutil.CheckNotNull(suppressions)
	b.RecordSuppressionsWithDescription(suppressions, "")
}

// Java: JSDocInfo.Builder.recordSuppression(String)
func (b *JSDocInfo_Builder) RecordSuppression(suppression string) {
	b.recordSuppressionWarnings([]string{RhinoStringPool_AddOrGet(suppression)}, "")
}

// RecordModifies records the list of modifies warnings.
//
// Java: JSDocInfo.Builder.recordModifies(Set<String>)
func (b *JSDocInfo_Builder) RecordModifies(modifies *jutil.LinkedHashSet[string]) bool {
	return !b.HasAnySingletonSideEffectTags() && b.PopulateProp(JSDocInfo_MODIFIES, JSDocInfo_InternSet(modifies))
}

// RecordType records a type; returns true if the type was recorded and false
// if it is invalid or was already defined.
//
// Java: JSDocInfo.Builder.recordType(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordType(type_ *JSTypeExpression) bool {
	return type_ != nil && !b.HasAnyTypeRelatedTags() && b.PopulateProp(JSDocInfo_TYPE, type_)
}

// Java: JSDocInfo.Builder.recordInlineType()
func (b *JSDocInfo_Builder) RecordInlineType() {
	b.PopulateBit(JSDocInfo_Bit_INLINE_TYPE, true)
}

// RecordTypedef records that the JSDocInfo being built should be populated
// with a `typedef`'d type.
//
// Java: JSDocInfo.Builder.recordTypedef(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordTypedef(type_ *JSTypeExpression) bool {
	return type_ != nil &&
		!b.HasAnyTypeRelatedTags() &&
		b.GetProp(JSDocInfo_TEMPLATE_TYPE_NAMES) == nil &&
		b.PopulateProp(JSDocInfo_TYPEDEF_TYPE, type_)
}

// RecordReturnType records a return type; returns true if the return type
// was recorded and false if it is invalid or was already defined.
//
// Java: JSDocInfo.Builder.recordReturnType(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordReturnType(type_ *JSTypeExpression) bool {
	return type_ != nil && !b.HasAnySingletonTypeTags() && b.PopulateProp(JSDocInfo_RETURN_TYPE, type_)
}

// RecordReturnDescription records a return description; returns true if the
// return description was recorded and false if it is invalid or was already
// defined.
//
// Java: JSDocInfo.Builder.recordReturnDescription(String)
func (b *JSDocInfo_Builder) RecordReturnDescription(description string) bool {
	if !b.ShouldParseDocumentation() {
		return true
	}
	return b.PopulateProp(JSDocInfo_RETURN_DESCRIPTION, description)
}

// RecordDefineType records the type of a define.
//
// 'Define' values are special constants that may be manipulated by the
// compiler. They are designed to mimic the #define command in the C
// preprocessor.
//
// Java: JSDocInfo.Builder.recordDefineType(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordDefineType(type_ *JSTypeExpression) bool {
	if type_ != nil && !b.CheckBit(JSDocInfo_Bit_CONST) && !b.CheckBit(JSDocInfo_Bit_DEFINE) && b.RecordType(type_) {
		return b.PopulateBit(JSDocInfo_Bit_DEFINE, true)
	}
	return false
}

// RecordEnumParameterType records a parameter type to an enum; returns true
// if the enum's parameter type was recorded and false if it was invalid or
// already defined.
//
// Java: JSDocInfo.Builder.recordEnumParameterType(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordEnumParameterType(type_ *JSTypeExpression) bool {
	if type_ != nil && !b.HasAnyTypeRelatedTags() {
		b.SetProp(JSDocInfo_ENUM_PARAMETER_TYPE, type_)
		b.populated = true
		return true
	}
	return false
}

// TODO(tbreisacher): Disallow nullable types here. If someone writes "@this
// {Foo}" in their JS we automatically treat it as though they'd written
// "@this {!Foo}". But, if the type node is created in the compiler (e.g. in
// the WizPass) we should explicitly add the '!'

// RecordThisType records a type for `@this` annotation; returns true if the
// type was recorded and false if it is invalid or if it collided with
// `@enum` or `@type` annotations.
//
// Java: JSDocInfo.Builder.recordThisType(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordThisType(type_ *JSTypeExpression) bool {
	return type_ != nil && !b.HasAnySingletonTypeTags() && b.PopulateProp(JSDocInfo_THIS_TYPE, type_)
}

// RecordBaseType records a base type; returns true if the base type was
// recorded and false if it was already defined.
//
// Java: JSDocInfo.Builder.recordBaseType(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordBaseType(type_ *JSTypeExpression) bool {
	return type_ != nil && !b.HasAnySingletonTypeTags() && b.PopulateProp(JSDocInfo_BASE_TYPE, type_)
}

// ChangeBaseType changes a base type, even if one has already been set on
// currentInfo; returns true if the base type was changed successfully.
//
// Java: JSDocInfo.Builder.changeBaseType(JSTypeExpression)
func (b *JSDocInfo_Builder) ChangeBaseType(type_ *JSTypeExpression) bool {
	if type_ != nil && !b.HasAnySingletonTypeTags() {
		b.SetProp(JSDocInfo_BASE_TYPE, type_)
		b.populated = true
		return true
	}
	return false
}

// RecordConstancy records that the JSDocInfo being built should have its
// JSDocInfo#isConstant() flag set to true; returns true if the constancy was
// recorded and false if it was already defined.
//
// Java: JSDocInfo.Builder.recordConstancy()
func (b *JSDocInfo_Builder) RecordConstancy() bool {
	return b.PopulateBit(JSDocInfo_Bit_CONST, true)
}

// RecordMutable records that the JSDocInfo being built should have its
// JSDocInfo#isConstant() flag set to false; returns true if the mutability
// was recorded and false if it was already defined.
//
// Java: JSDocInfo.Builder.recordMutable()
func (b *JSDocInfo_Builder) RecordMutable() bool {
	return b.PopulateBit(JSDocInfo_Bit_CONST, false)
}

// RecordFinality records that the JSDocInfo being built should have its
// JSDocInfo#isFinal() flag set to true; returns true if the finality was
// recorded and false if it was already defined.
//
// Java: JSDocInfo.Builder.recordFinality()
func (b *JSDocInfo_Builder) RecordFinality() bool {
	return b.PopulateBit(JSDocInfo_Bit_FINAL, true)
}

// RecordDescription records a description giving context for translation
// (i18n); returns true if the description was recorded and false if the
// description was invalid or was already defined.
//
// Java: JSDocInfo.Builder.recordDescription(String)
func (b *JSDocInfo_Builder) RecordDescription(description string) bool {
	return b.PopulateProp(JSDocInfo_DESCRIPTION, description)
}

// RecordTsType records a tsType giving context for .d.ts generation.
//
// Java: JSDocInfo.Builder.recordTsType(String)
func (b *JSDocInfo_Builder) RecordTsType(tsType string) {
	b.populated = true
	b.appendStringProp(JSDocInfo_TS_TYPES, tsType)
}

// RecordMeaning records a meaning giving context for translation (i18n).
// Different meanings will result in different translations; returns true if
// the meaning was successfully updated.
//
// Java: JSDocInfo.Builder.recordMeaning(String)
func (b *JSDocInfo_Builder) RecordMeaning(meaning string) bool {
	return b.PopulateProp(JSDocInfo_MEANING, meaning)
}

// RecordAlternateMessageId records an ID for an alternate message to be used
// if this message is not yet translated; returns true if the alternate
// message ID was successfully updated.
//
// Java: JSDocInfo.Builder.recordAlternateMessageId(String)
func (b *JSDocInfo_Builder) RecordAlternateMessageId(alternateMessageId string) bool {
	return b.PopulateProp(JSDocInfo_ALTERNATE_MESSAGE_ID, alternateMessageId)
}

// RecordClosurePrimitiveId records an identifier for a Closure Primitive.
// function; returns true if the id was successfully updated.
//
// Java: JSDocInfo.Builder.recordClosurePrimitiveId(String)
func (b *JSDocInfo_Builder) RecordClosurePrimitiveId(closurePrimitiveId string) bool {
	return b.PopulateProp(JSDocInfo_CLOSURE_PRIMITIVE_ID, RhinoStringPool_AddOrGet(closurePrimitiveId))
}

// RecordFileOverview records a fileoverview description; returns true if the
// description was recorded and false if the description was invalid or was
// already defined.
//
// Java: JSDocInfo.Builder.recordFileOverview(String)
func (b *JSDocInfo_Builder) RecordFileOverview(description string) bool {
	b.SetBit(JSDocInfo_Bit_FILEOVERVIEW, true)
	b.populated = true
	if !b.ShouldParseDocumentation() {
		return true
	}
	return b.PopulateProp(JSDocInfo_FILEOVERVIEW_DESCRIPTION, description)
}

// RecordEnhance records enhanced namespace; returns true if the enhanced
// namespace was recorded.
//
// Java: JSDocInfo.Builder.recordEnhance(String)
func (b *JSDocInfo_Builder) RecordEnhance(namespace string) bool {
	b.SetBit(JSDocInfo_Bit_ENHANCED_NAMESPACE, true)
	b.populated = true
	return b.PopulateProp(JSDocInfo_ENHANCED_NAMESPACE, namespace)
}

// IsModsRecorded returns whether this builder recorded a modded namespace.
//
// Java: JSDocInfo.Builder.isModsRecorded()
func (b *JSDocInfo_Builder) IsModsRecorded() bool {
	return b.GetProp(JSDocInfo_MODS) != nil
}

// RecordMods records modded namespace; returns true if the modded namespace
// was recorded.
//
// Java: JSDocInfo.Builder.recordMods(String)
func (b *JSDocInfo_Builder) RecordMods(namespace string) bool {
	return b.PopulateProp(JSDocInfo_MODS, RhinoStringPool_AddOrGet(namespace))
}

// Java: JSDocInfo.Builder.recordLicense(String)
func (b *JSDocInfo_Builder) RecordLicense(license string) bool {
	b.SetProp(JSDocInfo_LICENSE, RhinoStringPool_AddOrGet(license))
	b.populated = true
	return true
}

// Java: JSDocInfo.Builder.addLicense(String)
func (b *JSDocInfo_Builder) AddLicense(license string) bool {
	if slices.Contains(b.licenseTexts, license) {
		return false
	}
	b.licenseTexts = append(b.licenseTexts, license)
	txt := jsdocGet[string](b.GetProp(JSDocInfo_LICENSE)) // Strings.nullToEmpty(txt)
	return b.RecordLicense(RhinoStringPool_AddOrGet(txt + license))
}

// RecordNoCompile records that the JSDocInfo being built should have its
// JSDocInfo#isNoCompile() flag set to `true`. Returns `true` if the no
// compile flag was recorded and `false` if it was already recorded
//
// Java: JSDocInfo.Builder.recordNoCompile()
func (b *JSDocInfo_Builder) RecordNoCompile() bool {
	return b.PopulateBit(JSDocInfo_Bit_NOCOMPILE, true)
}

// RecordNoDts records that the JSDocInfo being built should have its
// JSDocInfo#isNoDtsOutput() flag set to `true`. Returns `true` if the no
// compile flag was recorded and `false` if it was already recorded
//
// Java: JSDocInfo.Builder.recordNoDts()
func (b *JSDocInfo_Builder) RecordNoDts() bool {
	return b.PopulateBit(JSDocInfo_Bit_NODTS, true)
}

// RecordNoCollapse records that the JSDocInfo being built should have its
// JSDocInfo#isNoCollapse() flag set to `true`. Returns `true` if the no
// collapse flag was recorded and `false` if it was already recorded
//
// Java: JSDocInfo.Builder.recordNoCollapse()
func (b *JSDocInfo_Builder) RecordNoCollapse() bool {
	return b.PopulateBit(JSDocInfo_Bit_NOCOLLAPSE, true)
}

// RecordNoInline records that the JSDocInfo being built should have its
// JSDocInfo#isNoInline() flag set to `true`. Returns `true` if the no inline
// flag was recorded and `false` if it was already recorded
//
// Java: JSDocInfo.Builder.recordNoInline()
func (b *JSDocInfo_Builder) RecordNoInline() bool {
	return b.PopulateBit(JSDocInfo_Bit_NOINLINE, true)
}

// RecordEncourageInlining records that the JSDocInfo being built should have
// its JSDocInfo#isEncourageInlining() flag set to `true`. Returns `true` if
// the encourageInlining flag was recorded and `false` if it was already
// recorded
//
// Java: JSDocInfo.Builder.recordEncourageInlining()
func (b *JSDocInfo_Builder) RecordEncourageInlining() bool {
	return b.PopulateBit(JSDocInfo_Bit_ENCOURAGE_INLINING, true)
}

// RecordRequireInlining records that the JSDocInfo being built should have
// its JSDocInfo#isRequireInlining() flag set to `true`. Returns `true` if the
// requireInlining flag was recorded and `false` if it was already recorded
//
// Java: JSDocInfo.Builder.recordRequireInlining()
func (b *JSDocInfo_Builder) RecordRequireInlining() bool {
	return b.PopulateBit(JSDocInfo_Bit_REQUIRE_INLINING, true)
}

// RecordPureOrBreakMyCode records that the JSDocInfo being built should have
// its JSDocInfo#isPureOrBreakMyCode() flag set to `true`. Returns `true` if
// the no pureOrBreakMyCode flag was recorded and `false` if it was already
// recorded
//
// Java: JSDocInfo.Builder.recordPureOrBreakMyCode()
func (b *JSDocInfo_Builder) RecordPureOrBreakMyCode() bool {
	return b.PopulateBit(JSDocInfo_Bit_PURE_OR_BREAK_MY_CODE, true)
}

// RecordCollapsibleOrBreakMyCode records that the JSDocInfo being built
// should have its JSDocInfo#isCollapsibleOrBreakMyCode() flag set to `true`.
// Returns `true` if the no collapsibleOrBreakMyCode flag was recorded and
// `false` if it was already recorded
//
// Java: JSDocInfo.Builder.recordCollapsibleOrBreakMyCode()
func (b *JSDocInfo_Builder) RecordCollapsibleOrBreakMyCode() bool {
	return b.PopulateBit(JSDocInfo_Bit_COLLAPSIBLE_OR_BREAK_MY_CODE, true)
}

// RecordConstructor records that the JSDocInfo being built should have its
// JSDocInfo#isConstructor() flag set to `true`. Returns `true` if the
// constructor was recorded and `false` if it was already defined or it was
// incompatible with the existing flags
//
// Java: JSDocInfo.Builder.recordConstructor()
func (b *JSDocInfo_Builder) RecordConstructor() bool {
	return !b.HasAnySingletonTypeTags() &&
		!b.IsConstructorOrInterface() &&
		b.PopulateBit(JSDocInfo_Bit_CONSTRUCTOR, true)
}

// RecordImplicitMatch records that the JSDocInfo being built should have its
// JSDocInfo#usesImplicitMatch() flag set to `true`. Returns `true` if the
// `@record` tag was recorded and `false` if it was already defined or it was
// incompatible with the existing flags
//
// Java: JSDocInfo.Builder.recordImplicitMatch()
func (b *JSDocInfo_Builder) RecordImplicitMatch() bool {
	return !b.HasAnySingletonTypeTags() &&
		!b.IsConstructorOrInterface() &&
		b.PopulateBit(JSDocInfo_Bit_RECORD, true) &&
		b.PopulateBit(JSDocInfo_Bit_INTERFACE, true)
}

// Java: JSDocInfo.Builder.recordProvideGoog()
func (b *JSDocInfo_Builder) RecordProvideGoog() bool {
	return b.PopulateBit(JSDocInfo_Bit_PROVIDE_GOOG, true)
}

// Java: JSDocInfo.Builder.recordProvideAlreadyProvided()
func (b *JSDocInfo_Builder) RecordProvideAlreadyProvided() bool {
	return b.PopulateBit(JSDocInfo_Bit_PROVIDE_ALREADY_PROVIDED, true)
}

// IsConstructorRecorded reports whether the JSDocInfo being built will have
// its JSDocInfo#isConstructor() flag set to true.
//
// Java: JSDocInfo.Builder.isConstructorRecorded()
func (b *JSDocInfo_Builder) IsConstructorRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_CONSTRUCTOR)
}

// RecordUnrestricted records that the JSDocInfo being built should have its
// JSDocInfo#makesUnrestricted() flag set to true; returns true if annotation
// was recorded and false if it was already defined or it was incompatible
// with the existing flags.
//
// Java: JSDocInfo.Builder.recordUnrestricted()
func (b *JSDocInfo_Builder) RecordUnrestricted() bool {
	return !b.HasAnySingletonTypeTags() &&
		((b.bits & (JSDocInfo_Bit_INTERFACE.Mask() | JSDocInfo_Bit_DICT.Mask() | JSDocInfo_Bit_STRUCT.Mask())) == 0) &&
		b.PopulateBit(JSDocInfo_Bit_UNRESTRICTED, true)
}

// Java: JSDocInfo.Builder.isUnrestrictedRecorded()
func (b *JSDocInfo_Builder) IsUnrestrictedRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_UNRESTRICTED)
}

// RecordAbstract records that the JSDocInfo being built should have its
// JSDocInfo#isAbstract() flag set to true; returns true if the flag was
// recorded and false if it was already defined or it was incompatible with
// the existing flags.
//
// Java: JSDocInfo.Builder.recordAbstract()
func (b *JSDocInfo_Builder) RecordAbstract() bool {
	visibility := b.GetProp(JSDocInfo_VISIBILITY)
	isPrivate := visibility != nil && jsdocGet[JSDocInfo_Visibility](visibility) == JSDocInfo_Visibility_PRIVATE
	return !b.HasAnySingletonTypeTags() &&
		((b.bits & (JSDocInfo_Bit_INTERFACE.Mask() | JSDocInfo_Bit_FINAL.Mask())) == 0) &&
		!isPrivate &&
		b.PopulateBit(JSDocInfo_Bit_ABSTRACT, true)
}

// RecordStruct records that the JSDocInfo being built should have its
// JSDocInfo#makesStructs() flag set to true; returns true if the struct was
// recorded and false if it was already defined or it was incompatible with
// the existing flags.
//
// Java: JSDocInfo.Builder.recordStruct()
func (b *JSDocInfo_Builder) RecordStruct() bool {
	return !b.HasAnySingletonTypeTags() &&
		((b.bits & (JSDocInfo_Bit_DICT.Mask() | JSDocInfo_Bit_UNRESTRICTED.Mask())) == 0) &&
		b.PopulateBit(JSDocInfo_Bit_STRUCT, true)
}

// Java: JSDocInfo.Builder.isStructRecorded()
func (b *JSDocInfo_Builder) IsStructRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_STRUCT)
}

// RecordDict records that the JSDocInfo being built should have its
// JSDocInfo#makesDicts() flag set to true; returns true if the dict was
// recorded and false if it was already defined or it was incompatible with
// the existing flags.
//
// Java: JSDocInfo.Builder.recordDict()
func (b *JSDocInfo_Builder) RecordDict() bool {
	return !b.HasAnySingletonTypeTags() &&
		((b.bits & (JSDocInfo_Bit_STRUCT.Mask() | JSDocInfo_Bit_UNRESTRICTED.Mask())) == 0) &&
		b.PopulateBit(JSDocInfo_Bit_DICT, true)
}

// Java: JSDocInfo.Builder.isDictRecorded()
func (b *JSDocInfo_Builder) IsDictRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_DICT)
}

// RecordOverride records that the JSDocInfo being built should have its
// JSDocInfo#isOverride() flag set to true.
//
// Java: JSDocInfo.Builder.recordOverride()
func (b *JSDocInfo_Builder) RecordOverride() bool {
	return b.PopulateBit(JSDocInfo_Bit_OVERRIDE, true)
}

// RecordDeprecated records that the JSDocInfo being built should have its
// JSDocInfo#isDeprecated() flag set to true.
//
// Java: JSDocInfo.Builder.recordDeprecated()
func (b *JSDocInfo_Builder) RecordDeprecated() bool {
	return b.PopulateBit(JSDocInfo_Bit_DEPRECATED, true)
}

// RecordInterface records that the JSDocInfo being built should have its
// JSDocInfo#isInterface() flag set to true; returns true if the flag was
// recorded and false if it was already defined or it was incompatible with
// the existing flags.
//
// Java: JSDocInfo.Builder.recordInterface()
func (b *JSDocInfo_Builder) RecordInterface() bool {
	return !b.HasAnySingletonTypeTags() &&
		((b.bits & (JSDocInfo_Bit_CONSTRUCTOR.Mask() | JSDocInfo_Bit_ABSTRACT.Mask())) == 0) &&
		b.PopulateBit(JSDocInfo_Bit_INTERFACE, true)
}

// RecordExport records that the JSDocInfo being built should have its
// JSDocInfo#isExport() flag set to true.
//
// Java: JSDocInfo.Builder.recordExport()
func (b *JSDocInfo_Builder) RecordExport() bool {
	return b.PopulateBit(JSDocInfo_Bit_EXPORT, true)
}

// RemoveExport records that the JSDocInfo being built should have its
// JSDocInfo#isExport() flag set to false.
//
// Java: JSDocInfo.Builder.removeExport()
func (b *JSDocInfo_Builder) RemoveExport() bool {
	return b.PopulateBit(JSDocInfo_Bit_EXPORT, false)
}

// RecordImplicitCast records that the JSDocInfo being built should have its
// JSDocInfo#isImplicitCast() flag set to true.
//
// Java: JSDocInfo.Builder.recordImplicitCast()
func (b *JSDocInfo_Builder) RecordImplicitCast() bool {
	return b.PopulateBit(JSDocInfo_Bit_IMPLICITCAST, true)
}

// RecordNoSideEffects records that the JSDocInfo being built should have its
// JSDocInfo#isNoSideEffects() flag set to true.
//
// Java: JSDocInfo.Builder.recordNoSideEffects()
func (b *JSDocInfo_Builder) RecordNoSideEffects() bool {
	return !b.HasAnySingletonSideEffectTags() &&
		b.IsPropEmpty(JSDocInfo_THROWS_ANNOTATIONS) &&
		b.PopulateBit(JSDocInfo_Bit_NOSIDEEFFECTS, true)
}

// IsNoSideEffectsRecorded returns whether current JSDoc is annotated with
// @nosideeffects.
//
// Java: JSDocInfo.Builder.isNoSideEffectsRecorded()
func (b *JSDocInfo_Builder) IsNoSideEffectsRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_NOSIDEEFFECTS)
}

// IsModifiesRecorded returns whether current JSDoc is annotated with
// @modifies.
//
// Java: JSDocInfo.Builder.isModifiesRecorded()
func (b *JSDocInfo_Builder) IsModifiesRecorded() bool {
	return !b.IsPropEmpty(JSDocInfo_MODIFIES)
}

// IsThrowsRecorded returns whether current JSDoc is annotated with @throws.
//
// Java: JSDocInfo.Builder.isThrowsRecorded()
func (b *JSDocInfo_Builder) IsThrowsRecorded() bool {
	return !b.IsPropEmpty(JSDocInfo_THROWS_ANNOTATIONS)
}

// RecordExterns records that the JSDocInfo being built should have its
// JSDocInfo#isExterns() flag set to true.
//
// Java: JSDocInfo.Builder.recordExterns()
func (b *JSDocInfo_Builder) RecordExterns() bool {
	return b.PopulateBit(JSDocInfo_Bit_EXTERNS, true)
}

// RecordNoCoverage records that the JSDocInfo being built should have its
// JSDocInfo#isNoCoverage() flag set to true.
//
// Java: JSDocInfo.Builder.recordNoCoverage()
func (b *JSDocInfo_Builder) RecordNoCoverage() bool {
	return b.PopulateBit(JSDocInfo_Bit_NOCOVERAGE, true)
}

// RecordTypeSummary records that the JSDocInfo being built should have its
// JSDocInfo#isTypeSummary() flag set to true.
//
// Java: JSDocInfo.Builder.recordTypeSummary()
func (b *JSDocInfo_Builder) RecordTypeSummary() bool {
	return b.PopulateBit(JSDocInfo_Bit_TYPE_SUMMARY, true)
}

// IsInterfaceRecorded reports whether the JSDocInfo being built will have
// its JSDocInfo#isInterface() flag set to true.
//
// Java: JSDocInfo.Builder.isInterfaceRecorded()
func (b *JSDocInfo_Builder) IsInterfaceRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_INTERFACE)
}

// HasParameter returns whether a parameter of the given name has already
// been recorded.
//
// Java: JSDocInfo.Builder.hasParameter(String)
func (b *JSDocInfo_Builder) HasParameter(name string) bool {
	for _, p := range jsdocGet[[]JSDocParam](b.GetProp(JSDocInfo_PARAMETERS)) {
		if p.Name == name {
			return true
		}
	}
	return false
}

// RecordImplementedInterface records an implemented interface.
//
// Java: JSDocInfo.Builder.recordImplementedInterface(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordImplementedInterface(interfaceType *JSTypeExpression) bool {
	return interfaceType != nil && b.AddUnique(JSDocInfo_IMPLEMENTED_INTERFACES, interfaceType)
}

// RecordExtendedInterface records an extended interface type.
//
// Java: JSDocInfo.Builder.recordExtendedInterface(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordExtendedInterface(interfaceType *JSTypeExpression) bool {
	return interfaceType != nil && b.AddUnique(JSDocInfo_EXTENDED_INTERFACES, interfaceType)
}

// Java: JSDocInfo.Builder.addUnique(Property<ArrayList<JSTypeExpression>>, JSTypeExpression)
func (b *JSDocInfo_Builder) AddUnique(prop *JSDocInfo_Property, elem *JSTypeExpression) bool {
	list := jsdocGet[[]*JSTypeExpression](b.GetProp(prop))
	for _, e := range list {
		if elem.IsEquivalentTo(e) {
			return false
		}
	}
	b.SetProp(prop, append(list, elem))
	b.populated = true
	return true
}

// RecordLends records that we're lending to another name.
//
// Java: JSDocInfo.Builder.recordLends(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordLends(name *JSTypeExpression) bool {
	return !b.HasAnyTypeRelatedTags() && b.PopulateProp(JSDocInfo_LENDS_NAME, name)
}

// IsNgInjectRecorded returns whether current JSDoc is annotated with
// `@ngInject`.
//
// Java: JSDocInfo.Builder.isNgInjectRecorded()
func (b *JSDocInfo_Builder) IsNgInjectRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_NG_INJECT)
}

// RecordNgInject records that we'd like to add `$inject` property inferred
// from parameters.
//
// Java: JSDocInfo.Builder.recordNgInject(boolean)
func (b *JSDocInfo_Builder) RecordNgInject(ngInject bool) bool {
	return b.PopulateBit(JSDocInfo_Bit_NG_INJECT, true)
}

// IsWizactionRecorded returns whether current JSDoc is annotated with
// `@wizaction`.
//
// Java: JSDocInfo.Builder.isWizactionRecorded()
func (b *JSDocInfo_Builder) IsWizactionRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_WIZ_ACTION)
}

// RecordWizaction records that this method is to be exposed as a wizaction.
//
// Java: JSDocInfo.Builder.recordWizaction()
func (b *JSDocInfo_Builder) RecordWizaction() bool {
	return b.PopulateBit(JSDocInfo_Bit_WIZ_ACTION, true)
}

// IsWizcallbackRecorded returns if current JSDoc is annotated with
// `@wizcallback`.
//
// Java: JSDocInfo.Builder.isWizcallbackRecorded()
func (b *JSDocInfo_Builder) IsWizcallbackRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_WIZ_CALLBACK)
}

// RecordWizcallback records that this method is to be exposed as a
// wizcallback.
//
// Java: JSDocInfo.Builder.recordWizcallback()
func (b *JSDocInfo_Builder) RecordWizcallback() bool {
	return b.PopulateBit(JSDocInfo_Bit_WIZ_CALLBACK, true)
}

// IsPolymerBehaviorRecorded returns whether current JSDoc is annotated with
// `@polymerBehavior`.
//
// Java: JSDocInfo.Builder.isPolymerBehaviorRecorded()
func (b *JSDocInfo_Builder) IsPolymerBehaviorRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_POLYMER_BEHAVIOR)
}

// RecordPolymerBehavior records that this method is to be exposed as a
// polymerBehavior.
//
// Java: JSDocInfo.Builder.recordPolymerBehavior()
func (b *JSDocInfo_Builder) RecordPolymerBehavior() bool {
	return b.PopulateBit(JSDocInfo_Bit_POLYMER_BEHAVIOR, true)
}

// IsPolymerRecorded returns whether current JSDoc is annotated with
// `@polymer`.
//
// Java: JSDocInfo.Builder.isPolymerRecorded()
func (b *JSDocInfo_Builder) IsPolymerRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_POLYMER)
}

// RecordPolymer records that this method is to be exposed as a polymer
// element.
//
// Java: JSDocInfo.Builder.recordPolymer()
func (b *JSDocInfo_Builder) RecordPolymer() bool {
	return b.PopulateBit(JSDocInfo_Bit_POLYMER, true)
}

// IsCustomElementRecorded returns whether current JSDoc is annotated with
// `@customElement`.
//
// Java: JSDocInfo.Builder.isCustomElementRecorded()
func (b *JSDocInfo_Builder) IsCustomElementRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_CUSTOM_ELEMENT)
}

// RecordCustomElement records that this method is to be exposed as a
// customElement.
//
// Java: JSDocInfo.Builder.recordCustomElement()
func (b *JSDocInfo_Builder) RecordCustomElement() bool {
	return b.PopulateBit(JSDocInfo_Bit_CUSTOM_ELEMENT, true)
}

// IsMixinClassRecorded returns whether current JSDoc is annotated with
// `@mixinClass`.
//
// Java: JSDocInfo.Builder.isMixinClassRecorded()
func (b *JSDocInfo_Builder) IsMixinClassRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_MIXIN_CLASS)
}

// RecordMixinClass records that this method is to be exposed as a mixinClass.
//
// Java: JSDocInfo.Builder.recordMixinClass()
func (b *JSDocInfo_Builder) RecordMixinClass() bool {
	return b.PopulateBit(JSDocInfo_Bit_MIXIN_CLASS, true)
}

// IsMixinFunctionRecorded returns whether current JSDoc is annotated with
// `@mixinFunction`.
//
// Java: JSDocInfo.Builder.isMixinFunctionRecorded()
func (b *JSDocInfo_Builder) IsMixinFunctionRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_MIXIN_FUNCTION)
}

// RecordMixinFunction records that this method is to be exposed as a
// mixinFunction.
//
// Java: JSDocInfo.Builder.recordMixinFunction()
func (b *JSDocInfo_Builder) RecordMixinFunction() bool {
	return b.PopulateBit(JSDocInfo_Bit_MIXIN_FUNCTION, true)
}

// IsSassGeneratedCssTsRecorded returns whether current JSDoc is annotated
// with `@sassGeneratedCssTs`.
//
// Java: JSDocInfo.Builder.isSassGeneratedCssTsRecorded()
func (b *JSDocInfo_Builder) IsSassGeneratedCssTsRecorded() bool {
	return b.CheckBit(JSDocInfo_Bit_SASS_GENERATED_CSS_TS)
}

// RecordSassGeneratedCssTs records that this script was generated by Sass
// from a.css.ts file.
//
// Java: JSDocInfo.Builder.recordSassGeneratedCssTs()
func (b *JSDocInfo_Builder) RecordSassGeneratedCssTs() {
	b.PopulateBit(JSDocInfo_Bit_SASS_GENERATED_CSS_TS, true)
}

// LogTypeInCompiler returns whether we should log the type of values passed
// to this function.
//
// Java: JSDocInfo.Builder.logTypeInCompiler()
func (b *JSDocInfo_Builder) LogTypeInCompiler() bool {
	return b.CheckBit(JSDocInfo_Bit_LOG_TYPE_IN_COMPILER)
}

// RecordLogTypeInCompiler records that the types of values passed to this
// method should be logged in the compiler.
//
// Java: JSDocInfo.Builder.recordLogTypeInCompiler()
func (b *JSDocInfo_Builder) RecordLogTypeInCompiler() bool {
	return b.PopulateBit(JSDocInfo_Bit_LOG_TYPE_IN_COMPILER, true)
}

// IsClosureUnawareCode returns whether JSDoc is annotated with the
// `@closureUnaware` annotation.
//
// Java: JSDocInfo.Builder.isClosureUnawareCode()
func (b *JSDocInfo_Builder) IsClosureUnawareCode() bool {
	return b.GetProp(JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE) != nil
}

// RecordClosureUnawareCode records that this JSDoc was annotated with the
// `@closureUnaware` annotation.
//
// Java: JSDocInfo.Builder.recordClosureUnawareCode()
func (b *JSDocInfo_Builder) RecordClosureUnawareCode() bool {
	return b.PopulateProp(JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE, JSDocInfo_PerFileClosureUnawareMode_UNSPECIFIED)
}

// RecordClosureUnawareCodeWithMode records that this JSDoc was annotated
// with the `@closureUnaware` annotation.
//
// Java: JSDocInfo.Builder.recordClosureUnawareCode(PerFileClosureUnawareMode)
func (b *JSDocInfo_Builder) RecordClosureUnawareCodeWithMode(mode JSDocInfo_PerFileClosureUnawareMode) bool {
	return b.PopulateProp(JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE, mode)
}

// RemoveClosureUnawareCode removes the `@closureUnaware` annotation from
// this JSDoc, returning true the annotation was present before removal.
//
// Java: JSDocInfo.Builder.removeClosureUnawareCode()
func (b *JSDocInfo_Builder) RemoveClosureUnawareCode() bool {
	existing := b.removeProp(JSDocInfo_CLOSURE_UNAWARE_CONFIG_VALUE)
	return existing != nil
}

// RecordUsedViaDotConstructor records that this JSDoc was annotated with the
// `@usedViaDotConstructor` annotation.
//
// Java: JSDocInfo.Builder.recordUsedViaDotConstructor()
func (b *JSDocInfo_Builder) RecordUsedViaDotConstructor() bool {
	return b.PopulateBit(JSDocInfo_Bit_USED_VIA_DOT_CONSTRUCTOR, true)
}

// SetType force-sets the type, while recordType backs off. This is useful
// for (e.g.) copyFromWithNewType.
//
// TODO(sdh): this is a new method - consider removing it in favor of
// recordType? The main difference is that this force-sets the type, while
// recordType backs off.
//
// Java: JSDocInfo.Builder.setType(JSTypeExpression)
func (b *JSDocInfo_Builder) SetType(type_ *JSTypeExpression) *JSDocInfo_Builder {
	b.removeProp(JSDocInfo_RETURN_TYPE)
	b.removeProp(JSDocInfo_ENUM_PARAMETER_TYPE)
	b.removeProp(JSDocInfo_TYPEDEF_TYPE)
	if type_ == nil {
		b.removeProp(JSDocInfo_TYPE)
	} else {
		b.SetProp(JSDocInfo_TYPE, type_)
	}
	return b
}

// HasAnyTypeRelatedTags reports whether the current doc info has other type
// tags, like `@param` or `@return` or `@type` or etc.
//
// Java: JSDocInfo.Builder.hasAnyTypeRelatedTags()
func (b *JSDocInfo_Builder) HasAnyTypeRelatedTags() bool {
	return (b.bits&(JSDocInfo_Bit_CONSTRUCTOR.Mask()|JSDocInfo_Bit_INTERFACE.Mask()|JSDocInfo_Bit_ABSTRACT.Mask())) != 0 ||
		b.HasAnyParameters() ||
		b.GetProp(JSDocInfo_RETURN_TYPE) != nil ||
		b.GetProp(JSDocInfo_BASE_TYPE) != nil ||
		!b.IsPropEmpty(JSDocInfo_EXTENDED_INTERFACES) ||
		b.GetProp(JSDocInfo_LENDS_NAME) != nil ||
		b.GetProp(JSDocInfo_THIS_TYPE) != nil ||
		b.HasAnySingletonTypeTags()
}

// Java: JSDocInfo.Builder.hasAnyParameters()
func (b *JSDocInfo_Builder) HasAnyParameters() bool {
	return len(jsdocGet[[]JSDocParam](b.GetProp(JSDocInfo_PARAMETERS))) != 0
}

// Java: JSDocInfo.Builder.isPropEmpty(Property<? extends Collection<?>>)
func (b *JSDocInfo_Builder) IsPropEmpty(prop *JSDocInfo_Property) bool {
	return prop.IsDefault(b.GetProp(prop))
}

// HasAnySingletonTypeTags reports whether the current doc info has any of
// the singleton type tags that may not appear with other type tags, like
// `@type` or `@typedef`.
//
// Java: JSDocInfo.Builder.hasAnySingletonTypeTags()
func (b *JSDocInfo_Builder) HasAnySingletonTypeTags() bool {
	return b.GetProp(JSDocInfo_TYPE) != nil ||
		b.GetProp(JSDocInfo_TYPEDEF_TYPE) != nil ||
		b.GetProp(JSDocInfo_ENUM_PARAMETER_TYPE) != nil
}

// HasAnySingletonSideEffectTags reports whether the current doc info has
// any of the singleton type tags that may not appear with other type tags,
// like `@type` or `@typedef`.
//
// Java: JSDocInfo.Builder.hasAnySingletonSideEffectTags()
func (b *JSDocInfo_Builder) HasAnySingletonSideEffectTags() bool {
	return b.CheckBit(JSDocInfo_Bit_NOSIDEEFFECTS) || !b.IsPropEmpty(JSDocInfo_MODIFIES)
}

// Java: JSDocInfo.Builder.isConstructorOrInterface()
func (b *JSDocInfo_Builder) IsConstructorOrInterface() bool {
	return (b.bits & (JSDocInfo_Bit_CONSTRUCTOR.Mask() | JSDocInfo_Bit_INTERFACE.Mask())) != 0
}

func (b *JSDocInfo_Builder) hasPropKey(prop *JSDocInfo_Property) bool {
	return (b.propMask & prop.mask) != 0
}

func (b *JSDocInfo_Builder) removeProp(prop *JSDocInfo_Property) any {
	if (b.propMask & prop.mask) == 0 {
		return nil
	}
	idx := bits.OnesCount64(uint64(b.propMask & (prop.mask - 1)))
	old := b.entries[idx]
	copy(b.entries[idx:], b.entries[idx+1:])
	b.entries[len(b.entries)-1] = nil
	b.entries = b.entries[:len(b.entries)-1]
	b.propMask &^= prop.mask
	return old
}

// Java: JSDocInfo.Builder.setProp(Property<T>, T)
func (b *JSDocInfo_Builder) SetProp(prop *JSDocInfo_Property, value any) {
	idx := bits.OnesCount64(uint64(b.propMask & (prop.mask - 1)))
	if (b.propMask & prop.mask) != 0 {
		b.entries[idx] = value
		return
	}
	b.propMask |= prop.mask
	if b.entries == nil {
		b.entries = b.entriesBuf[:0]
	}
	b.entries = append(b.entries, nil)
	copy(b.entries[idx+1:], b.entries[idx:])
	b.entries[idx] = value
}

// Java: JSDocInfo.Builder.getProp(Property<T>)
func (b *JSDocInfo_Builder) GetProp(prop *JSDocInfo_Property) any {
	if (b.propMask & prop.mask) == 0 {
		return nil
	}
	idx := bits.OnesCount64(uint64(b.propMask & (prop.mask - 1)))
	return b.entries[idx]
}

func (b *JSDocInfo_Builder) populateParamEntry(prop *JSDocInfo_Property, key string, typ *JSTypeExpression) bool {
	params := jsdocGet[[]JSDocParam](b.GetProp(prop))
	for i := range params {
		if params[i].Name == key {
			if params[i].Type == nil {
				params[i].Type = typ
				b.populated = true
				return true
			}
			return false
		}
	}
	b.SetProp(prop, append(params, JSDocParam{Name: key, Type: typ}))
	b.populated = true
	return true
}

func (b *JSDocInfo_Builder) populateStringMapEntry(prop *JSDocInfo_Property, key string, val string) bool {
	kvs := jsdocGet[[]jsDocKV](b.GetProp(prop))
	for _, kv := range kvs {
		if kv.key == key {
			return false
		}
	}
	b.SetProp(prop, append(kvs, jsDocKV{key: key, val: val}))
	b.populated = true
	return true
}

func (b *JSDocInfo_Builder) populateNodeMapEntry(prop *JSDocInfo_Property, key string, expr Node) bool {
	m := jsdocGet[*jutil.LinkedHashMap[string, Node]](b.GetProp(prop))
	if m == nil {
		m = jutil.NewLinkedHashMap[string, Node]()
		b.SetProp(prop, m)
	} else if prev, exists := m.Lookup(key); exists && !prev.IsNil() {
		return false
	}
	m.Put(key, expr)
	b.populated = true
	return true
}

// Java: JSDocInfo.Builder.populateProp(Property<T>, T)
func (b *JSDocInfo_Builder) PopulateProp(prop *JSDocInfo_Property, value any) bool {
	b.populated = true
	idx := bits.OnesCount64(uint64(b.propMask & (prop.mask - 1)))
	if (b.propMask & prop.mask) != 0 {
		if isNilPropValue(prop.kind, b.entries[idx]) {
			b.entries[idx] = value
			return true
		}
		return false
	}
	b.propMask |= prop.mask
	if b.entries == nil {
		b.entries = b.entriesBuf[:0]
	}
	b.entries = append(b.entries, nil)
	copy(b.entries[idx+1:], b.entries[idx:])
	b.entries[idx] = value
	return true
}

// Java: JSDocInfo.Builder.checkBit(Bit)
func (b *JSDocInfo_Builder) CheckBit(bit JSDocInfo_Bit) bool {
	return (b.bits & bit.Mask()) != 0
}

// Java: JSDocInfo.Builder.setBit(Bit, boolean)
func (b *JSDocInfo_Builder) SetBit(bit JSDocInfo_Bit, value bool) {
	if value {
		b.bits |= bit.Mask()
	} else {
		b.bits &^= bit.Mask()
	}
}

// Java: JSDocInfo.Builder.populateBit(Bit, boolean)
func (b *JSDocInfo_Builder) PopulateBit(bit JSDocInfo_Bit, value bool) bool {
	if b.CheckBit(bit) != value {
		b.SetBit(bit, value)
		b.populated = true
		return b.populated
	}
	return false
}

// Java: JSDocInfo.internSet(Set<String>)
func JSDocInfo_InternSet(strings_ *jutil.LinkedHashSet[string]) *jutil.LinkedHashSet[string] {
	interned := jutil.NewLinkedHashSet[string]()
	for s := range strings_.All() {
		interned.Add(RhinoStringPool_AddOrGet(s))
	}
	return interned
}
