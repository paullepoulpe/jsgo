// Port of com/google/javascript/rhino/jstype/StaticTypedSlot.java @ closure-compiler v20260415

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

// StaticTypedSlot must be implemented by variables that can appear as
// members of a StaticTypedScope.
//
// Go-only: language-adapter — getDeclaration() (StaticTypedRef in Java) and
// getScope() (StaticTypedScope in Java) keep StaticSlot's signatures; see
// StaticTypedScope.
type StaticTypedSlot interface {
	StaticSlot

	// GetType returns the type information, if any, for this slot, or null
	// if no type is declared for it.
	GetType() JSType

	// IsTypeInferred returns whether the type has been inferred (as opposed
	// to declared).
	IsTypeInferred() bool
}
