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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// jsDocKV is an ordered (key, val) string pair for parameter descriptions.
type jsDocKV struct {
	key string
	val string
}

// JSDocSuppression represents a single @suppress annotation's deduplicated
// warning list and optional trailing description.
type JSDocSuppression struct {
	Warnings    []string
	Description string
}

// FormatJSDocSuppressions formats suppressions identically to Java's
// ImmutableMap<ImmutableSet<String>, String>.toString().
func FormatJSDocSuppressions(s []JSDocSuppression) string {
	var b strings.Builder
	b.WriteByte('{')
	for i, sup := range s {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteByte('[')
		for j, w := range sup.Warnings {
			if j > 0 {
				b.WriteString(", ")
			}
			b.WriteString(w)
		}
		b.WriteString("]=")
		b.WriteString(sup.Description)
	}
	b.WriteByte('}')
	return b.String()
}

func sameWarningSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, w := range a {
		if !slices.Contains(b, w) {
			return false
		}
	}
	return true
}

func formatJsDocKVs(s []jsDocKV) string {
	var b strings.Builder
	b.WriteByte('{')
	for i, kv := range s {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(kv.key)
		b.WriteByte('=')
		b.WriteString(kv.val)
	}
	b.WriteByte('}')
	return b.String()
}

func jsDocStringSetFromSlice(s []string) *jutil.LinkedHashSet[string] {
	if len(s) == 0 {
		return jsDocInfo_noSuppressions
	}
	return jutil.NewLinkedHashSet(s...)
}

var (
	// jsDocInfo_noSuppressions is GetSuppressions' result for a JSDocInfo
	// without suppressions (ImmutableSet.of()). Never modified.
	jsDocInfo_noSuppressions = jutil.NewLinkedHashSet[string]()

	emptyStringSlice              = []string{}
	jsDocModifiesThis             = []string{RhinoStringPool_AddOrGet("this")}
	jsDocModifiesArguments        = []string{RhinoStringPool_AddOrGet("arguments")}
	jsDocModifiesThisAndArguments = []string{RhinoStringPool_AddOrGet("this"), RhinoStringPool_AddOrGet("arguments")}
	jsDocModifiesArgumentsAndThis = []string{RhinoStringPool_AddOrGet("arguments"), RhinoStringPool_AddOrGet("this")}
)

func (i *JSDocInfo) getStringProp(p *JSDocInfo_Property) (string, bool) {
	v := p.getRaw(i)
	if v == nil {
		return "", false
	}
	return jsdocGet[string](v), true
}

func (i *JSDocInfo) getStringPtrProp(p *JSDocInfo_Property) *string {
	s, ok := i.getStringProp(p)
	if !ok {
		return nil
	}
	return &s
}

func (i *JSDocInfo) suppressions() []JSDocSuppression {
	return jsdocGet[[]JSDocSuppression](JSDocInfo_SUPPRESSIONS.getRaw(i))
}

// GetSuppressions returns the set of suppressed warnings.
//
// Java: JSDocInfo.getSuppressions()
func (i *JSDocInfo) GetSuppressions() *jutil.LinkedHashSet[string] {
	suppressions := i.suppressions()
	if len(suppressions) == 0 {
		// Go-only: performance — ImmutableSet.of() is a shared singleton;
		// callers only read the result, so share one empty set instead of
		// allocating one per call (this runs for every JSDoc, externs
		// included).
		return jsDocInfo_noSuppressions
	}
	names := jutil.NewLinkedHashSet[string]()
	for _, sup := range suppressions {
		for _, name := range sup.Warnings {
			names.Add(name)
		}
	}
	return names
}

// HasAnySuppressions reports whether i has at least one @suppress warning name
// without allocating a merged set when multiple @suppress tags are present.
func (i *JSDocInfo) HasAnySuppressions() bool {
	for _, sup := range i.suppressions() {
		if len(sup.Warnings) > 0 {
			return true
		}
	}
	return false
}

// HasSuppression reports whether warning name is in i's @suppress set without
// allocating a merged set when multiple @suppress tags are present.
func (i *JSDocInfo) HasSuppression(name string) bool {
	for _, sup := range i.suppressions() {
		if slices.Contains(sup.Warnings, name) {
			return true
		}
	}
	return false
}

// ForEachSuppression visits each suppressed warning name on i (may visit a
// duplicate name across multiple @suppress tags if the caller is idempotent).
func (i *JSDocInfo) ForEachSuppression(visit func(string) bool) {
	for _, sup := range i.suppressions() {
		for _, name := range sup.Warnings {
			if !visit(name) {
				return
			}
		}
	}
}

// GetSuppressionsAndTheirDescription returns the slice of @suppress warning
// groups and their corresponding descriptions.
//
// Java: JSDocInfo.getSuppressionsAndTheirDescription()
func (i *JSDocInfo) GetSuppressionsAndTheirDescription() []JSDocSuppression {
	if suppressions := i.suppressions(); suppressions != nil {
		return suppressions
	}
	return []JSDocSuppression{}
}

// GetModifiesSlice returns the deduplicated @modifies annotations without
// allocating a LinkedHashSet.
func (i *JSDocInfo) GetModifiesSlice() []string {
	if i == nil {
		return nil
	}
	return jsdocGet[[]string](JSDocInfo_MODIFIES.getRaw(i))
}

// GetModifies returns the set of sideeffect notations.
//
// Java: JSDocInfo.getModifies()
func (i *JSDocInfo) GetModifies() *jutil.LinkedHashSet[string] {
	return jsDocStringSetFromSlice(i.GetModifiesSlice())
}

// Java: JSDocInfo.hasModifies()
func (i *JSDocInfo) HasModifies() bool {
	return (i.propertyKeysBitset & JSDocInfo_MODIFIES.mask) != 0
}

// ModifiesThis gets the value of the @modifies{this} annotation stored in
// the doc info.
//
// Java: JSDocInfo.modifiesThis()
func (i *JSDocInfo) ModifiesThis() bool {
	return slices.Contains(i.GetModifiesSlice(), "this")
}

// HasSideEffectsArgumentsAnnotation returns whether the @modifies annotation
// includes "arguments" or any named parameters.
//
// Java: JSDocInfo.hasSideEffectsArgumentsAnnotation()
func (i *JSDocInfo) HasSideEffectsArgumentsAnnotation() bool {
	modifies := i.GetModifiesSlice()
	// TODO(johnlenz): if we start tracking parameters individually this should
	// simply be a check for "arguments".
	return len(modifies) > 1 || (len(modifies) == 1 && modifies[0] != "this")
}

// HasThrowsAnnotations reports whether any @throws annotations are present.
func (i *JSDocInfo) HasThrowsAnnotations() bool {
	return (i.propertyKeysBitset & JSDocInfo_THROWS_ANNOTATIONS.mask) != 0
}

// GetThrowsAnnotations returns the list of thrown types and descriptions as
// text.
//
// Java: JSDocInfo.getThrowsAnnotations()
func (i *JSDocInfo) GetThrowsAnnotations() []string {
	if annotations := jsdocGet[[]string](JSDocInfo_THROWS_ANNOTATIONS.getRaw(i)); annotations != nil {
		return annotations
	}
	return []string{}
}

// GetClosurePrimitiveId returns the `@closurePrimitive {id}` identifier.
//
// Java: JSDocInfo.getClosurePrimitiveId()
func (i *JSDocInfo) GetClosurePrimitiveId() *string {
	return i.getStringPtrProp(JSDocInfo_CLOSURE_PRIMITIVE_ID)
}

// GetClosurePrimitiveIdString returns the `@closurePrimitive` id without allocating a *string.
func (i *JSDocInfo) GetClosurePrimitiveIdString() (string, bool) {
	return i.getStringProp(JSDocInfo_CLOSURE_PRIMITIVE_ID)
}

// HasClosurePrimitiveId reports whether this JSDoc is annotated with
// `@closurePrimitive`.
//
// Java: JSDocInfo.hasClosurePrimitiveId()
func (i *JSDocInfo) HasClosurePrimitiveId() bool {
	return (i.propertyKeysBitset & JSDocInfo_CLOSURE_PRIMITIVE_ID.mask) != 0
}

// GetDescription gets the description specified by the `@desc` annotation.
//
// Java: JSDocInfo.getDescription()
func (i *JSDocInfo) GetDescription() *string {
	return i.getStringPtrProp(JSDocInfo_DESCRIPTION)
}

// GetDescriptionString returns the `@desc` description without allocating a *string.
func (i *JSDocInfo) GetDescriptionString() (string, bool) {
	return i.getStringProp(JSDocInfo_DESCRIPTION)
}

// GetMeaning gets the meaning specified by the `@meaning` annotation.
//
// In localization systems, two messages with the same content but different
// "meanings" may be translated differently. By default, we use the name of
// the variable that the message is initialized to as the "meaning" of the
// message.
//
// But some code generators (like Closure Templates) inject their own meaning
// with the jsdoc `@meaning` annotation.
//
// Java: JSDocInfo.getMeaning()
func (i *JSDocInfo) GetMeaning() *string {
	return i.getStringPtrProp(JSDocInfo_MEANING)
}

// GetMeaningString returns the `@meaning` string without allocating a *string.
func (i *JSDocInfo) GetMeaningString() (string, bool) {
	return i.getStringProp(JSDocInfo_MEANING)
}

// GetAlternateMessageId gets the alternate message ID specified by the
// `@alternateMessageId` annotation.
//
// In localization systems, if we migrate from one message ID algorithm to
// another, we can specify the old one via `@alternateMessageId`. This allows
// the product to use the previous translation while waiting for the new one
// to be translated.
//
// Some code generators (like Closure Templates) inject this.
//
// Java: JSDocInfo.getAlternateMessageId()
func (i *JSDocInfo) GetAlternateMessageId() *string {
	return i.getStringPtrProp(JSDocInfo_ALTERNATE_MESSAGE_ID)
}

// GetAlternateMessageIdString returns the `@alternateMessageId` string without allocating a *string.
func (i *JSDocInfo) GetAlternateMessageIdString() (string, bool) {
	return i.getStringProp(JSDocInfo_ALTERNATE_MESSAGE_ID)
}

// GetDeprecationReason returns the deprecation reason or null if none
// specified.
//
// Java: JSDocInfo.getDeprecationReason()
func (i *JSDocInfo) GetDeprecationReason() *string {
	return i.getStringPtrProp(JSDocInfo_DEPRECATION_REASON)
}

// GetDeprecationReasonString returns the `@deprecated` reason without allocating a *string.
func (i *JSDocInfo) GetDeprecationReasonString() (string, bool) {
	return i.getStringProp(JSDocInfo_DEPRECATION_REASON)
}

// GetLicense gets the description specified by the `@license` annotation.
//
// Java: JSDocInfo.getLicense()
func (i *JSDocInfo) GetLicense() *string {
	return i.getStringPtrProp(JSDocInfo_LICENSE)
}

// GetLicenseString returns the `@license` text without allocating a *string.
func (i *JSDocInfo) GetLicenseString() (string, bool) {
	return i.getStringProp(JSDocInfo_LICENSE)
}

func (i *JSDocInfo) parameterDescriptions() []jsDocKV {
	return jsdocGet[[]jsDocKV](JSDocInfo_PARAMETER_DESCRIPTIONS.getRaw(i))
}

// HasDescriptionForParameter returns whether a description exists for the
// parameter with the specified name.
//
// Java: JSDocInfo.hasDescriptionForParameter(String)
func (i *JSDocInfo) HasDescriptionForParameter(name string) bool {
	_, ok := i.GetDescriptionForParameterString(name)
	return ok
}

// GetDescriptionForParameter returns the description for the parameter with
// the given name, if its exists.
//
// Java: JSDocInfo.getDescriptionForParameter(String)
func (i *JSDocInfo) GetDescriptionForParameter(name string) *string {
	d, ok := i.GetDescriptionForParameterString(name)
	if !ok {
		return nil
	}
	return &d
}

// GetDescriptionForParameterString returns the description for the parameter
// with the given name without allocating a *string.
func (i *JSDocInfo) GetDescriptionForParameterString(name string) (string, bool) {
	for _, kv := range i.parameterDescriptions() {
		if kv.key == name {
			return kv.val, true
		}
	}
	return "", false
}

// GetAuthors returns the list of authors or null if none.
//
// Java: JSDocInfo.getAuthors()
func (i *JSDocInfo) GetAuthors() []string {
	return jsdocGet[[]string](JSDocInfo_AUTHORS.getRaw(i))
}

// GetReferences returns the list of references or null if none.
//
// Java: JSDocInfo.getReferences()
func (i *JSDocInfo) GetReferences() []string {
	return jsdocGet[[]string](JSDocInfo_SEES.getRaw(i))
}

// GetReturnDescription returns the description of the returned object or
// null if none specified.
//
// Java: JSDocInfo.getReturnDescription()
func (i *JSDocInfo) GetReturnDescription() *string {
	return i.getStringPtrProp(JSDocInfo_RETURN_DESCRIPTION)
}

// GetReturnDescriptionString returns the `@return` description without allocating a *string.
func (i *JSDocInfo) GetReturnDescriptionString() (string, bool) {
	return i.getStringProp(JSDocInfo_RETURN_DESCRIPTION)
}

// GetBlockDescription returns the block-level description or null if none
// specified.
//
// Java: JSDocInfo.getBlockDescription()
func (i *JSDocInfo) GetBlockDescription() *string {
	return i.getStringPtrProp(JSDocInfo_BLOCK_DESCRIPTION)
}

// GetBlockDescriptionString returns the block-level description without allocating a *string.
func (i *JSDocInfo) GetBlockDescriptionString() (string, bool) {
	return i.getStringProp(JSDocInfo_BLOCK_DESCRIPTION)
}

// HasFileOverview returns whether this has a fileoverview flag.
//
// Java: JSDocInfo.hasFileOverview()
func (i *JSDocInfo) HasFileOverview() bool {
	return i.CheckBit(JSDocInfo_Bit_FILEOVERVIEW)
}

// GetFileOverview returns the file overview or null if none specified.
//
// Java: JSDocInfo.getFileOverview()
func (i *JSDocInfo) GetFileOverview() *string {
	return i.getStringPtrProp(JSDocInfo_FILEOVERVIEW_DESCRIPTION)
}

// GetFileOverviewString returns the `@fileoverview` description without allocating a *string.
func (i *JSDocInfo) GetFileOverviewString() (string, bool) {
	return i.getStringProp(JSDocInfo_FILEOVERVIEW_DESCRIPTION)
}

// HasEnhance returns whether this has an enhanced namespace.
//
// Java: JSDocInfo.hasEnhance()
func (i *JSDocInfo) HasEnhance() bool {
	return i.CheckBit(JSDocInfo_Bit_ENHANCED_NAMESPACE)
}

// GetEnhance returns the enhanced namespace or null if none is specified.
//
// Java: JSDocInfo.getEnhance()
func (i *JSDocInfo) GetEnhance() *string {
	return i.getStringPtrProp(JSDocInfo_ENHANCED_NAMESPACE)
}

// GetEnhanceString returns the `@enhance` namespace without allocating a *string.
func (i *JSDocInfo) GetEnhanceString() (string, bool) {
	return i.getStringProp(JSDocInfo_ENHANCED_NAMESPACE)
}

// HasMods returns whether this has a modded namespace.
//
// Java: JSDocInfo.hasMods()
func (i *JSDocInfo) HasMods() bool {
	return (i.propertyKeysBitset & JSDocInfo_MODS.mask) != 0
}

// GetMods returns the modded namespace or null if none is specified.
//
// Java: JSDocInfo.getMods()
func (i *JSDocInfo) GetMods() *string {
	return i.getStringPtrProp(JSDocInfo_MODS)
}

// GetModsString returns the `@mods` namespace without allocating a *string.
func (i *JSDocInfo) GetModsString() (string, bool) {
	return i.getStringProp(JSDocInfo_MODS)
}

// GetMarkers gets the list of all markers for the documentation in this
// JSDoc.
//
// Java: JSDocInfo.getMarkers()
func (i *JSDocInfo) GetMarkers() []JSDocInfo_Marker {
	if markers := jsdocGet[[]JSDocInfo_Marker](JSDocInfo_MARKERS.getRaw(i)); markers != nil {
		return markers
	}
	return []JSDocInfo_Marker{}
}

// GetOriginalCommentString returns the original JSDoc comment string.
// Returns null unless parseJsDocDocumentation is enabled via the
// ParserConfig.
//
// Java: JSDocInfo.getOriginalCommentString()
func (i *JSDocInfo) GetOriginalCommentString() *string {
	return i.getStringPtrProp(JSDocInfo_SOURCE_COMMENT)
}

// GetOriginalCommentStringValue returns the original JSDoc comment string without allocating a *string.
func (i *JSDocInfo) GetOriginalCommentStringValue() (string, bool) {
	return i.getStringProp(JSDocInfo_SOURCE_COMMENT)
}

// Java: JSDocInfo.getOriginalCommentPosition()
func (i *JSDocInfo) GetOriginalCommentPosition() int {
	pos := JSDocInfo_ORIGINAL_COMMENT_POSITION.getRaw(i)
	if pos != nil {
		return jsdocGet[int](pos)
	}
	return 0
}

// IsAtSignCodePresent returns whether the `@code` is present within this
// JSDocInfo.
//
// Java: JSDocInfo.isAtSignCodePresent()
func (i *JSDocInfo) IsAtSignCodePresent() bool {
	entireComment, ok := i.GetOriginalCommentStringValue()
	return ok && strings.Contains(entireComment, "@code")
}

// Java: JSDocInfo.isDocumentationIncluded()
func (i *JSDocInfo) IsDocumentationIncluded() bool {
	return i.CheckBit(JSDocInfo_Bit_INCLUDE_DOCUMENTATION)
}

// IsConstant returns whether this JSDocInfo implies that annotated value is
// constant.
//
// Java: JSDocInfo.isConstant()
func (i *JSDocInfo) IsConstant() bool {
	// @desc is used with goog.getMsg to define messages to be translated, and
	// thus must be @const in order for translation to work correctly.
	return (i.propertyBits&(JSDocInfo_Bit_CONST.Mask()|JSDocInfo_Bit_DEFINE.Mask()|JSDocInfo_Bit_FINAL.Mask())) != 0 ||
		(i.propertyKeysBitset&JSDocInfo_DESCRIPTION.mask) != 0
}

// HasConstAnnotation returns whether the `@const` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.hasConstAnnotation()
func (i *JSDocInfo) HasConstAnnotation() bool {
	return i.CheckBit(JSDocInfo_Bit_CONST)
}

// IsFinal returns whether the `@final` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isFinal()
func (i *JSDocInfo) IsFinal() bool {
	return i.CheckBit(JSDocInfo_Bit_FINAL)
}

// IsDefine returns whether the `@define` annotation is present on this
// JSDocInfo. If this annotation is present, then the #getType() method will
// retrieve the define type.
//
// Java: JSDocInfo.isDefine()
func (i *JSDocInfo) IsDefine() bool {
	return i.CheckBit(JSDocInfo_Bit_DEFINE)
}

// IsOverride returns whether the `@override` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isOverride()
func (i *JSDocInfo) IsOverride() bool {
	return i.CheckBit(JSDocInfo_Bit_OVERRIDE)
}

// IsDeprecated returns whether the `@deprecated` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.isDeprecated()
func (i *JSDocInfo) IsDeprecated() bool {
	return i.CheckBit(JSDocInfo_Bit_DEPRECATED)
}

// IsExport returns whether the `@export` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isExport()
func (i *JSDocInfo) IsExport() bool {
	return i.CheckBit(JSDocInfo_Bit_EXPORT)
}

// IsImplicitCast returns whether the `@implicitCast` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.isImplicitCast()
func (i *JSDocInfo) IsImplicitCast() bool {
	return i.CheckBit(JSDocInfo_Bit_IMPLICITCAST)
}

// IsNoSideEffects returns whether the `@nosideeffects` annotation is present
// on this JSDocInfo.
//
// Java: JSDocInfo.isNoSideEffects()
func (i *JSDocInfo) IsNoSideEffects() bool {
	return i.CheckBit(JSDocInfo_Bit_NOSIDEEFFECTS)
}

// IsExterns returns whether the `@externs` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isExterns()
func (i *JSDocInfo) IsExterns() bool {
	return i.CheckBit(JSDocInfo_Bit_EXTERNS)
}

// IsNoCoverage returns whether the `@nocoverage` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.isNoCoverage()
func (i *JSDocInfo) IsNoCoverage() bool {
	return i.CheckBit(JSDocInfo_Bit_NOCOVERAGE)
}

// IsTypeSummary returns whether the `@typeSummary` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.isTypeSummary()
func (i *JSDocInfo) IsTypeSummary() bool {
	return i.CheckBit(JSDocInfo_Bit_TYPE_SUMMARY)
}

// IsNoCompile returns whether the `@nocompile` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isNoCompile()
func (i *JSDocInfo) IsNoCompile() bool {
	return i.CheckBit(JSDocInfo_Bit_NOCOMPILE)
}

// IsNoDts returns whether the `@nodts` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isNoDts()
func (i *JSDocInfo) IsNoDts() bool {
	return i.CheckBit(JSDocInfo_Bit_NODTS)
}

// IsNoCollapse returns whether the `@nocollapse` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.isNoCollapse()
func (i *JSDocInfo) IsNoCollapse() bool {
	return i.CheckBit(JSDocInfo_Bit_NOCOLLAPSE)
}

// IsNoInline returns whether the `@noinline` annotation is present on this
// JSDocInfo.
//
// Java: JSDocInfo.isNoInline()
func (i *JSDocInfo) IsNoInline() bool {
	return i.CheckBit(JSDocInfo_Bit_NOINLINE)
}

// IsRequireInlining returns whether the `@requireInlining` annotation is
// present on this JSDocInfo.
//
// Java: JSDocInfo.isRequireInlining()
func (i *JSDocInfo) IsRequireInlining() bool {
	return i.CheckBit(JSDocInfo_Bit_REQUIRE_INLINING)
}

// IsEncourageInlining returns whether the `@encourageInlining` annotation is
// present on this JSDocInfo.
//
// Java: JSDocInfo.isEncourageInlining()
func (i *JSDocInfo) IsEncourageInlining() bool {
	return i.CheckBit(JSDocInfo_Bit_ENCOURAGE_INLINING)
}

// IsCollapsibleOrBreakMyCode returns whether the `@collapsibleOrBreakMyCode`
// annotation is present on this JSDocInfo.
//
// Java: JSDocInfo.isCollapsibleOrBreakMyCode()
func (i *JSDocInfo) IsCollapsibleOrBreakMyCode() bool {
	return i.CheckBit(JSDocInfo_Bit_COLLAPSIBLE_OR_BREAK_MY_CODE)
}

// IsPureOrBreakMyCode returns whether the `@pureOrBreakMyCode` annotation is
// present on this JSDocInfo.
//
// Java: JSDocInfo.isPureOrBreakMyCode()
func (i *JSDocInfo) IsPureOrBreakMyCode() bool {
	return i.CheckBit(JSDocInfo_Bit_PURE_OR_BREAK_MY_CODE)
}

// IsProvideGoog returns whether the `@provideGoog` annotation is present on
// this JSDocInfo.
//
// Java: JSDocInfo.isProvideGoog()
func (i *JSDocInfo) IsProvideGoog() bool {
	return i.CheckBit(JSDocInfo_Bit_PROVIDE_GOOG)
}

// IsProvideAlreadyProvided returns whether the `@provideAlreadyProvided`
// annotation is present on this JSDocInfo.
//
// Java: JSDocInfo.isProvideAlreadyProvided()
func (i *JSDocInfo) IsProvideAlreadyProvided() bool {
	return i.CheckBit(JSDocInfo_Bit_PROVIDE_ALREADY_PROVIDED)
}

// IsNgInject returns whether JSDoc is annotated with `@ngInject` annotation.
//
// Java: JSDocInfo.isNgInject()
func (i *JSDocInfo) IsNgInject() bool {
	return i.CheckBit(JSDocInfo_Bit_NG_INJECT)
}

// IsWizaction returns whether JSDoc is annotated with `@wizaction`
// annotation.
//
// Java: JSDocInfo.isWizaction()
func (i *JSDocInfo) IsWizaction() bool {
	return i.CheckBit(JSDocInfo_Bit_WIZ_ACTION)
}

// IsWizcallback returns if JSDoc is annotated with `@wizcallback` annotation.
//
// Java: JSDocInfo.isWizcallback()
func (i *JSDocInfo) IsWizcallback() bool {
	return i.CheckBit(JSDocInfo_Bit_WIZ_CALLBACK)
}

// IsPolymerBehavior returns whether JSDoc is annotated with
// `@polymerBehavior` annotation.
//
// Java: JSDocInfo.isPolymerBehavior()
func (i *JSDocInfo) IsPolymerBehavior() bool {
	return i.CheckBit(JSDocInfo_Bit_POLYMER_BEHAVIOR)
}

// IsPolymer returns whether JSDoc is annotated with `@polymer` annotation.
//
// Java: JSDocInfo.isPolymer()
func (i *JSDocInfo) IsPolymer() bool {
	return i.CheckBit(JSDocInfo_Bit_POLYMER)
}

// IsCustomElement returns whether JSDoc is annotated with `@customElement`
// annotation.
//
// Java: JSDocInfo.isCustomElement()
func (i *JSDocInfo) IsCustomElement() bool {
	return i.CheckBit(JSDocInfo_Bit_CUSTOM_ELEMENT)
}

// IsMixinClass returns whether JSDoc is annotated with `@mixinClass`
// annotation.
//
// Java: JSDocInfo.isMixinClass()
func (i *JSDocInfo) IsMixinClass() bool {
	return i.CheckBit(JSDocInfo_Bit_MIXIN_CLASS)
}

// IsMixinFunction returns whether JSDoc is annotated with `@mixinFunction`
// annotation.
//
// Java: JSDocInfo.isMixinFunction()
func (i *JSDocInfo) IsMixinFunction() bool {
	return i.CheckBit(JSDocInfo_Bit_MIXIN_FUNCTION)
}

// Java: JSDocInfo.isSassGeneratedCssTs()
func (i *JSDocInfo) IsSassGeneratedCssTs() bool {
	return i.CheckBit(JSDocInfo_Bit_SASS_GENERATED_CSS_TS)
}

// IsUsedViaDotConstructor returns whether JSDoc is annotated with the
// `@usedViaDotConstructor` annotation.
//
// Java: JSDocInfo.isUsedViaDotConstructor()
func (i *JSDocInfo) IsUsedViaDotConstructor() bool {
	return i.CheckBit(JSDocInfo_Bit_USED_VIA_DOT_CONSTRUCTOR)
}

// GetLogTypeInCompiler returns whether we should log the type of values
// passed to this function.
//
// Java: JSDocInfo.getLogTypeInCompiler()
func (i *JSDocInfo) GetLogTypeInCompiler() bool {
	return i.CheckBit(JSDocInfo_Bit_LOG_TYPE_IN_COMPILER)
}
