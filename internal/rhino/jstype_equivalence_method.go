// Port of com/google/javascript/rhino/jstype/EquivalenceMethod.java @ closure-compiler v20261021

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
 *   Nick Santos
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

// EquivalenceMethod represents different ways for comparing equality among
// types.
type EquivalenceMethod int32

const (
	// EquivalenceMethod_IDENTITY indicates that the two types should behave
	// exactly the same under all type operations.
	//
	// Thus, {string} != {?} and {Unresolved} != {?}
	EquivalenceMethod_IDENTITY EquivalenceMethod = iota
	// EquivalenceMethod_DATA_FLOW indicates that the two types are almost
	// exactly the same, and that a data flow analysis algorithm comparing them
	// should consider them equal.
	//
	// In traditional type inference, the types form a finite lattice, and
	// this ensures that type inference will terminate.
	//
	// In our type system, the unknown types do not obey the lattice rules. So
	// if we continue to perform inference over the unknown types, we may never
	// terminate.
	//
	// By treating all unknown types as equivalent for the purposes of data
	// flow analysis, we ensure that the algorithm will terminate.
	//
	// Thus, {string} != {?} and {Unresolved} ~= {?}
	EquivalenceMethod_DATA_FLOW
)

var equivalenceMethod_names = [...]string{"IDENTITY", "DATA_FLOW"}

// String is Enum.name().
func (m EquivalenceMethod) String() string { return equivalenceMethod_names[m] }
