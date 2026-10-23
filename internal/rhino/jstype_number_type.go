// Port of com/google/javascript/rhino/jstype/NumberType.java @ closure-compiler v20261015

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

// NumberType is the number type.
type NumberType struct {
	ValueTypeBase
}

// Java: NumberType(JSTypeRegistry)
func NewNumberType(registry *JSTypeRegistry) *NumberType {
	t := &NumberType{}
	t.objFlags = jsKindScalarValue
	t.InitValueTypeBase(t, registry)
	return t
}

// Java: NumberType.getTypeClass()
func (t *NumberType) GetTypeClass() JSTypeClass {
	return JSTypeClass_NUMBER
}

// Java: NumberType.testForEquality(JSType)
func (t *NumberType) TestForEquality(that JSType) base.Tri {
	result := t.ValueTypeBase.TestForEquality(that)
	if result != jsType_TRI_NULL {
		return result
	}
	if that.IsUnknownType() ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_OBJECT_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_NUMBER_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_STRING_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_BOOLEAN_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_BIGINT_TYPE)) {
		return base.Tri_UNKNOWN
	}
	return base.Tri_FALSE
}

// Java: NumberType.isNumberValueType()
func (t *NumberType) IsNumberValueType() bool {
	return true
}

// Java: NumberType.matchesNumberContext()
func (t *NumberType) MatchesNumberContext() bool {
	return true
}

// Java: NumberType.matchesStringContext()
func (t *NumberType) MatchesStringContext() bool {
	return true
}

// Java: NumberType.matchesObjectContext()
func (t *NumberType) MatchesObjectContext() bool {
	// TODO(user): Revisit this for ES4, which is stricter.
	return true
}

// Java: NumberType.getDisplayName()
func (t *NumberType) GetDisplayName() *string {
	return &jsType_displayName_number
}

// Java: NumberType.getPossibleToBooleanOutcomes()
func (t *NumberType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

// Java: NumberType.visit(Visitor<T>)
func (t *NumberType) Visit(visitor Visitor) any {
	return visitor.CaseNumberType()
}

// Java: NumberType.autoboxesTo()
func (t *NumberType) AutoboxesTo() JSType {
	return t.GetNativeType(JSTypeNative_NUMBER_OBJECT_TYPE)
}

var jsType_displayName_number = "number"
