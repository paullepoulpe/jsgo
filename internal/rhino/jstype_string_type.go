// Port of com/google/javascript/rhino/jstype/StringType.java @ closure-compiler v20261021

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

// StringType is the string type.
type StringType struct {
	ValueTypeBase
}

// Java: StringType(JSTypeRegistry)
func NewStringType(registry *JSTypeRegistry) *StringType {
	t := &StringType{}
	t.objFlags = jsKindScalarValue
	t.InitValueTypeBase(t, registry)
	return t
}

// Java: StringType.getTypeClass()
func (t *StringType) GetTypeClass() JSTypeClass {
	return JSTypeClass_STRING
}

// Java: StringType.testForEquality(JSType)
func (t *StringType) TestForEquality(that JSType) base.Tri {
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

// Java: StringType.isStringValueType()
func (t *StringType) IsStringValueType() bool {
	return true
}

// Java: StringType.matchesNumberContext()
func (t *StringType) MatchesNumberContext() bool {
	return true
}

// Java: StringType.matchesStringContext()
func (t *StringType) MatchesStringContext() bool {
	return true
}

// Java: StringType.matchesObjectContext()
func (t *StringType) MatchesObjectContext() bool {
	// TODO(user): Revisit this for ES4, which is stricter.
	return true
}

// Java: StringType.getDisplayName()
func (t *StringType) GetDisplayName() *string {
	return &jsType_displayName_string
}

// Java: StringType.autoboxesTo()
func (t *StringType) AutoboxesTo() JSType {
	return t.GetNativeType(JSTypeNative_STRING_OBJECT_TYPE)
}

// Java: StringType.getPossibleToBooleanOutcomes()
func (t *StringType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

var jsType_displayName_string = "string"
