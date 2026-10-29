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

// Java: JSDocInfo.Builder.recordInlineType()
func (b *JSDocInfo_Builder) RecordInlineType() {
	b.PopulateBit(JSDocInfo_Bit_INLINE_TYPE, true)
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

// RecordClosurePrimitiveId records an identifier for a Closure Primitive.
// function; returns true if the id was successfully updated.
//
// Java: JSDocInfo.Builder.recordClosurePrimitiveId(String)
func (b *JSDocInfo_Builder) RecordClosurePrimitiveId(closurePrimitiveId string) bool {
	return b.PopulateProp(JSDocInfo_CLOSURE_PRIMITIVE_ID, RhinoStringPool_AddOrGet(closurePrimitiveId))
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
		return true
	}
	return false
}
