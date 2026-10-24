// Port of com/google/javascript/rhino/jstype/KnownSymbolType.java @ closure-compiler v20261021

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

// KnownSymbolType is the type of a well-known (or unique) symbol, such as
// Symbol.iterator.
type KnownSymbolType struct {
	SymbolTypeBase

	name string
}

// Java: KnownSymbolType.TYPE_CLASS
const KnownSymbolType_TYPE_CLASS = JSTypeClass_WELL_KNOWN_SYMBOL

// Java: KnownSymbolType(JSTypeRegistry, String)
func NewKnownSymbolType(registry *JSTypeRegistry, name string) *KnownSymbolType {
	t := &KnownSymbolType{}
	t.objFlags = jsKindKnownSymbol
	t.InitSymbolTypeBase(t, registry)
	t.name = name
	return t
}

// Java: KnownSymbolType.getTypeClass()
func (t *KnownSymbolType) GetTypeClass() JSTypeClass {
	return KnownSymbolType_TYPE_CLASS
}

// Java: KnownSymbolType.getDisplayName()
func (t *KnownSymbolType) GetDisplayName() *string {
	return &t.name
}

// Java: KnownSymbolType.isKnownSymbolValueType()
func (t *KnownSymbolType) IsKnownSymbolValueType() bool {
	return true
}

// Java: KnownSymbolType.toMaybeKnownSymbolType()
func (t *KnownSymbolType) ToMaybeKnownSymbolType() *KnownSymbolType {
	return t
}
