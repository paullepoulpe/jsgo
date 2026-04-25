// Port of com/google/javascript/rhino/jstype/JSTypeClass.java
// @ closure-compiler v20260415

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1/GPL 2.0
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
 * Alternatively, the contents of this file may be used under the terms of
 * the GNU General Public License Version 2 or later (the "GPL"), in which
 * case the provisions of the GPL are applicable instead of those above. If
 * you wish to allow use of your version of this file only under the terms of
 * the GPL and not to allow others to use your version of this file under the
 * MPL, indicate your decision by deleting the provisions above and replacing
 * them with the notice and other provisions required by the GPL. If you do
 * not delete the provisions above, a recipient may use your version of this
 * file under either the MPL or the GPL.
 *
 * ***** END LICENSE BLOCK *****
 */

package jstype

// JSTypeClass is the set of concrete subtypes of JSType.
//
// Each class reports the corresponding element as its instance type-class
// discriminator. These values are not intended to be compatible with Liskov
// substitution; JSTypes have exactly one associated JSTypeClass and it must
// not affect the outcome of operations such as equality. Use of this type is
// tantamount to reflection (e.g. getClass or instanceof) and should only be
// used where such a comparison is appropriate and desirable.
//
// Use of this type is preferred to true reflection, when possible, for
// simplicity. This is safe because this package controls all JSType
// subclasses; they can all be known a priori.
//
// Abstract classes are not listed here because there would be no way to use
// those elements.
//
// TODO: wire existing JSType implementations to carry a JSTypeClass field.
type JSTypeClass int

const (
	JSTypeClassAll JSTypeClass = iota
	JSTypeClassArrow
	JSTypeClassBoolean
	JSTypeClassBigint
	JSTypeClassEnum
	JSTypeClassEnumElement
	JSTypeClassFunction
	JSTypeClassInstanceObject
	JSTypeClassNamed
	JSTypeClassNo
	JSTypeClassNoObject
	JSTypeClassNoResolved
	JSTypeClassNull
	JSTypeClassNumber
	JSTypeClassPrototypeObject
	JSTypeClassProxyObject
	JSTypeClassRecord
	JSTypeClassString
	JSTypeClassSymbol
	JSTypeClassWellKnownSymbol
	JSTypeClassTemplate
	JSTypeClassTemplatized
	JSTypeClassUnion
	JSTypeClassUnknown
	JSTypeClassVoid
)
