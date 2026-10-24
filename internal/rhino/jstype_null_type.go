// Port of com/google/javascript/rhino/jstype/NullType.java @ closure-compiler v20261021

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

// NullType is the null type.
type NullType struct {
	ValueTypeBase
}

// Java: NullType(JSTypeRegistry)
func NewNullType(registry *JSTypeRegistry) *NullType {
	t := &NullType{}
	t.objFlags = jsKindNull
	t.InitValueTypeBase(t, registry)
	return t
}

// Java: NullType.getTypeClass()
func (t *NullType) GetTypeClass() JSTypeClass {
	return JSTypeClass_NULL
}

// Java: NullType.isNullType()
func (t *NullType) IsNullType() bool {
	return true
}

// Java: NullType.isNullable()
func (t *NullType) IsNullable() bool {
	return true
}

// Java: NullType.matchesNumberContext()
func (t *NullType) MatchesNumberContext() bool {
	return true
}

// Java: NullType.matchesObjectContext()
func (t *NullType) MatchesObjectContext() bool {
	return false
}

// Java: NullType.matchesStringContext()
func (t *NullType) MatchesStringContext() bool {
	return true
}

// Java: NullType.restrictByNotNullOrUndefined()
func (t *NullType) RestrictByNotNullOrUndefined() JSType {
	return t.GetRegistry().GetNativeType(JSTypeNative_NO_TYPE)
}

// Java: NullType.restrictByNotNull()
func (t *NullType) RestrictByNotNull() JSType {
	return t.GetRegistry().GetNativeType(JSTypeNative_NO_TYPE)
}

// Java: NullType.testForEquality(JSType)
func (t *NullType) TestForEquality(that JSType) base.Tri {
	result := t.ValueTypeBase.TestForEquality(that)
	if result != jsType_TRI_NULL {
		return result
	}
	if that.IsNullType() || that.IsVoidType() {
		return base.Tri_TRUE
	}
	if that.IsUnknownType() || that.IsNullable() {
		return base.Tri_UNKNOWN
	}
	return base.Tri_FALSE
}

// Java: NullType.getDisplayName()
func (t *NullType) GetDisplayName() *string {
	return &jsType_displayName_null
}

// Java: NullType.getPossibleToBooleanOutcomes()
func (t *NullType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_FALSE
}

// Java: NullType.visit(Visitor<T>)
func (t *NullType) Visit(visitor Visitor) any {
	return visitor.CaseNullType()
}

var jsType_displayName_null = "null"
