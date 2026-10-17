// Port of com/google/javascript/rhino/jstype/StaticTypedScope.java @ closure-compiler v20261015

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

// StaticTypedScope must be implemented by any object that defines variables
// for the purposes of static analysis. It is distinguished from the
// Scriptable class that Rhino normally uses to represent a run-time scope.
//
// Go-only: language-adapter — Java narrows the return types of
// getParentScope() (to StaticTypedScope), getSlot(String) and
// getOwnSlot(String) (to StaticTypedSlot). Go interfaces have no covariant
// overrides, so those methods keep StaticScope's signatures; every
// StaticTypedScope implementation returns StaticTypedScope / StaticTypedSlot
// values from them, and callers that rely on the narrower Java type
// type-assert (`slot.(StaticTypedSlot)`).
type StaticTypedScope interface {
	StaticScope

	// GetTypeOfThis returns the expected type of `this` in the current scope.
	GetTypeOfThis() JSType

	// LookupQualifiedName looks up a given qualified name in the scope. if
	// that fails, looks up the component properties off of any owner types
	// that are in scope.
	//
	// This is always more or equally expensive as calling getSlot(String
	// name), so should only be used when necessary.
	//
	// This only returns declared qualified names and known properties. It
	// returns null given an inferred name.
	//
	// A Java default method: implementations that do not override it call
	// StaticTypedScope_LookupQualifiedName.
	LookupQualifiedName(qname QualifiedName) JSType
}

// StaticTypedScope_LookupQualifiedName is the default body of
// StaticTypedScope.lookupQualifiedName.
//
// Java: StaticTypedScope.lookupQualifiedName(QualifiedName) (default method)
func StaticTypedScope_LookupQualifiedName(s StaticTypedScope, qname QualifiedName) JSType {
	slot, _ := s.GetSlot(qname.Join()).(StaticTypedSlot)
	if slot != nil && !slot.IsTypeInferred() {
		return slot.GetType()
	} else if !qname.IsSimple() {
		type_ := s.LookupQualifiedName(qname.GetOwner())
		if type_ != nil && !type_.IsUnknownType() {
			return type_.FindPropertyType(qname.GetComponent())
		}
	}
	return nil
}
