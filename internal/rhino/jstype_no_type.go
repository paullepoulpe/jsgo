// Port of com/google/javascript/rhino/jstype/NoType.java @ closure-compiler v20261021

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

// NoType is the bottom type, representing the subclass of any value or
// object.
//
// Although JavaScript programmers can't explicitly denote the bottom type,
// it comes up in static analysis. For example, if we have:
//
//	var x = null;
//	if (x) {
//	  f(x);
//	}
//
// We need to be able to assign x a type within the f(x) call. Since it has
// no possible type, we assign x the NoType, so that f(x) is legal no matter
// what the type of f's first argument is.
type NoType struct {
	NoObjectTypeBase
}

// Java: NoType(JSTypeRegistry)
func NewNoType(registry *JSTypeRegistry) *NoType {
	t := &NoType{}
	t.objFlags = objKindNo
	t.InitNoObjectTypeBase(t, registry)
	return t
}

// Java: NoType.getTypeClass()
func (t *NoType) GetTypeClass() JSTypeClass {
	return JSTypeClass_NO
}

// Java: NoType.isNoObjectType()
func (t *NoType) IsNoObjectType() bool {
	return false
}

// Java: NoType.isNoType()
func (t *NoType) IsNoType() bool {
	return true
}

// Java: NoType.isNullable()
func (t *NoType) IsNullable() bool {
	return true
}

// Java: NoType.isVoidable()
func (t *NoType) IsVoidable() bool {
	return true
}

// Java: NoType.getPossibleToBooleanOutcomes()
func (t *NoType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

// Java: NoType.appendTo(TypeStringBuilder)
func (t *NoType) AppendTo(sb *TypeStringBuilder) {
	if sb.IsForAnnotations() {
		sb.Append("?")
	} else {
		sb.Append("None")
	}
}
