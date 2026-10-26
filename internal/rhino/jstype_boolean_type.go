// Port of com/google/javascript/rhino/jstype/BooleanType.java @ closure-compiler v20261021

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

// BooleanType is the boolean type.
type BooleanType struct {
	ValueTypeBase
}

// Java: BooleanType(JSTypeRegistry)
func NewBooleanType(registry *JSTypeRegistry) *BooleanType {
	t := &BooleanType{}
	t.objFlags = jsKindScalarValue
	t.InitValueTypeBase(t, registry)
	return t
}

// Java: BooleanType.getTypeClass()
func (t *BooleanType) GetTypeClass() JSTypeClass {
	return JSTypeClass_BOOLEAN
}

// Java: BooleanType.testForEquality(JSType)
func (t *BooleanType) TestForEquality(that JSType) base.Tri {
	result := t.ValueTypeBase.TestForEquality(that)
	if result != jsType_TRI_NULL {
		return result
	}
	if that.IsUnknownType() ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_NUMBER_STRING_BOOLEAN)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_BIGINT_TYPE)) ||
		that.IsObject() {
		return base.Tri_UNKNOWN
	}
	return base.Tri_FALSE
}

// Java: BooleanType.isBooleanValueType()
func (t *BooleanType) IsBooleanValueType() bool {
	return true
}

// Java: BooleanType.matchesNumberContext()
func (t *BooleanType) MatchesNumberContext() bool {
	return true
}

// Java: BooleanType.matchesStringContext()
func (t *BooleanType) MatchesStringContext() bool {
	return true
}

// Java: BooleanType.matchesObjectContext()
func (t *BooleanType) MatchesObjectContext() bool {
	// TODO(user): Revisit this for ES4, which is stricter.
	return true
}

// Java: BooleanType.autoboxesTo()
func (t *BooleanType) AutoboxesTo() JSType {
	return t.GetNativeType(JSTypeNative_BOOLEAN_OBJECT_TYPE)
}

// Java: BooleanType.getDisplayName()
func (t *BooleanType) GetDisplayName() *string {
	return &jsType_displayName_boolean
}

// Java: BooleanType.getPossibleToBooleanOutcomes()
func (t *BooleanType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

var jsType_displayName_boolean = "boolean"
