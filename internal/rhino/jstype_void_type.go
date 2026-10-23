// Port of com/google/javascript/rhino/jstype/VoidType.java @ closure-compiler v20261015

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

import "github.com/paullepoulpe/jsgo/internal/base"

// VoidType is the void type.
type VoidType struct {
	ValueTypeBase
}

// Java: VoidType(JSTypeRegistry)
func NewVoidType(registry *JSTypeRegistry) *VoidType {
	t := &VoidType{}
	t.objFlags = jsKindVoid
	t.InitValueTypeBase(t, registry)
	return t
}

// Java: VoidType.getTypeClass()
func (t *VoidType) GetTypeClass() JSTypeClass {
	return JSTypeClass_VOID
}

// Java: VoidType.restrictByNotNullOrUndefined()
func (t *VoidType) RestrictByNotNullOrUndefined() JSType {
	return t.GetRegistry().GetNativeType(JSTypeNative_NO_TYPE)
}

// Java: VoidType.restrictByNotUndefined()
func (t *VoidType) RestrictByNotUndefined() JSType {
	return t.GetRegistry().GetNativeType(JSTypeNative_NO_TYPE)
}

// Java: VoidType.testForEquality(JSType)
func (t *VoidType) TestForEquality(that JSType) base.Tri {
	if base.Tri_UNKNOWN == t.ValueTypeBase.TestForEquality(that) {
		return base.Tri_UNKNOWN
	}
	if that.IsSubtypeOf(t) || that.IsSubtypeOf(t.GetNativeType(JSTypeNative_NULL_TYPE)) {
		return base.Tri_TRUE
	}
	return base.Tri_FALSE
}

// Java: VoidType.matchesNumberContext()
func (t *VoidType) MatchesNumberContext() bool {
	return false
}

// Java: VoidType.matchesObjectContext()
func (t *VoidType) MatchesObjectContext() bool {
	return false
}

// Java: VoidType.matchesStringContext()
func (t *VoidType) MatchesStringContext() bool {
	return true
}

// Java: VoidType.isVoidType()
func (t *VoidType) IsVoidType() bool {
	return true
}

// Java: VoidType.isVoidable()
func (t *VoidType) IsVoidable() bool {
	return true
}

// Java: VoidType.isExplicitlyVoidable()
func (t *VoidType) IsExplicitlyVoidable() bool {
	return true
}

// Java: VoidType.getDisplayName()
func (t *VoidType) GetDisplayName() *string {
	return &jsType_displayName_undefined
}

// Java: VoidType.getPossibleToBooleanOutcomes()
func (t *VoidType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_FALSE
}

// Java: VoidType.visit(Visitor<T>)
func (t *VoidType) Visit(visitor Visitor) any {
	return visitor.CaseVoidType()
}

var jsType_displayName_undefined = "undefined"
