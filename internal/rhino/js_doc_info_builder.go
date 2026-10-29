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

// IsDescriptionRecorded returns whether this builder recorded a
// description.
//
// Java: JSDocInfo.Builder.isDescriptionRecorded()
func (b *JSDocInfo_Builder) IsDescriptionRecorded() bool {
	return b.GetProp(JSDocInfo_DESCRIPTION) != nil
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
	if b.hasPropKey(JSDocInfo_TYPEDEF_TYPE) {
		return false
	}
	for _, tr := range jsdocGet[[]JSDocTypeTransform](b.GetProp(JSDocInfo_TYPE_TRANSFORMATIONS)) {
		if tr.Name == name {
			return false
		}
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
	if b.HasAnySingletonSideEffectTags() {
		return false
	}
	return b.PopulateProp(JSDocInfo_MODIFIES, internModifiesSet(modifies))
}

// RecordModifiesSlice records a parsed @modifies slice without allocating a LinkedHashSet.
func (b *JSDocInfo_Builder) RecordModifiesSlice(modifies []string) bool {
	if b.HasAnySingletonSideEffectTags() {
		return false
	}
	return b.PopulateProp(JSDocInfo_MODIFIES, internModifiesSlice(modifies))
}

// RecordType records a type; returns true if the type was recorded and false
// if it is invalid or was already defined.
//
// Java: JSDocInfo.Builder.recordType(JSTypeExpression)
func (b *JSDocInfo_Builder) RecordType(type_ *JSTypeExpression) bool {
	return type_ != nil && !b.HasAnyTypeRelatedTags() && b.PopulateProp(JSDocInfo_TYPE, type_)
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
	transforms := jsdocGet[[]JSDocTypeTransform](b.GetProp(prop))
	for i := range transforms {
		if transforms[i].Name == key {
			if transforms[i].TypeNode.IsNil() {
				transforms[i].TypeNode = expr
				b.populated = true
				return true
			}
			return false
		}
	}
	b.SetProp(prop, append(transforms, JSDocTypeTransform{Name: key, TypeNode: expr}))
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

// Java: JSDocInfo.internSet(Set<String>)
func JSDocInfo_InternSet(strings_ *jutil.LinkedHashSet[string]) *jutil.LinkedHashSet[string] {
	interned := jutil.NewLinkedHashSet[string]()
	for s := range strings_.All() {
		interned.Add(RhinoStringPool_AddOrGet(s))
	}
	return interned
}

func internModifiesSet(modifies *jutil.LinkedHashSet[string]) []string {
	if modifies == nil || modifies.IsEmpty() {
		return emptyStringSlice
	}
	switch modifies.Size() {
	case 1:
		for s := range modifies.All() {
			return internSingleModifies(s)
		}
	case 2:
		var a, b string
		idx := 0
		for s := range modifies.All() {
			if idx == 0 {
				a = s
			} else {
				b = s
			}
			idx++
		}
		return internTwoModifies(a, b)
	}
	out := make([]string, 0, modifies.Size())
	for s := range modifies.All() {
		out = append(out, RhinoStringPool_AddOrGet(s))
	}
	return out
}

func internModifiesSlice(modifies []string) []string {
	if len(modifies) == 0 {
		return emptyStringSlice
	}
	if len(modifies) == 1 {
		return internSingleModifies(modifies[0])
	}
	if len(modifies) == 2 && modifies[0] != modifies[1] {
		return internTwoModifies(modifies[0], modifies[1])
	}
	out := make([]string, 0, len(modifies))
	for _, raw := range modifies {
		s := RhinoStringPool_AddOrGet(raw)
		if !slices.Contains(out, s) {
			out = append(out, s)
		}
	}
	if len(out) == 1 {
		return internSingleModifies(out[0])
	}
	if len(out) == 2 {
		return internTwoModifies(out[0], out[1])
	}
	return out
}

func internSingleModifies(s string) []string {
	switch s {
	case "this":
		return jsDocModifiesThis
	case "arguments":
		return jsDocModifiesArguments
	default:
		return []string{RhinoStringPool_AddOrGet(s)}
	}
}

func internTwoModifies(a, b string) []string {
	if a == "this" && b == "arguments" {
		return jsDocModifiesThisAndArguments
	}
	if a == "arguments" && b == "this" {
		return jsDocModifiesArgumentsAndThis
	}
	return []string{RhinoStringPool_AddOrGet(a), RhinoStringPool_AddOrGet(b)}
}
