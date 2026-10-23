// Port of com/google/javascript/rhino/jstype/SymbolType.java @ closure-compiler v20261015

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

// SymbolType is the symbol type.
//
// Java: class SymbolType extends ValueType (subclassed by KnownSymbolType).
type SymbolType interface {
	ValueType

	// GetSymbolTypeBase gives access to the base through the interface.
	//
	// Go-only: language-adapter — field access through the interface
	// (PORTING.md §4).
	GetSymbolTypeBase() *SymbolTypeBase
}

// SymbolTypeBase holds SymbolType's methods; plain symbol types are
// *SymbolTypeBase values.
type SymbolTypeBase struct {
	ValueTypeBase
}

// NewSymbolType constructs the plain symbol type.
//
// Java: SymbolType(JSTypeRegistry)
func NewSymbolType(registry *JSTypeRegistry) *SymbolTypeBase {
	t := &SymbolTypeBase{}
	t.objFlags = jsKindScalarValue
	t.InitSymbolTypeBase(t, registry)
	return t
}

// InitSymbolTypeBase is the Go form of the constructor for subclasses.
//
// Java: SymbolType(JSTypeRegistry)
func (t *SymbolTypeBase) InitSymbolTypeBase(self SymbolType, registry *JSTypeRegistry) {
	t.InitValueTypeBase(self, registry)
}

// GetSymbolTypeBase gives access to the base through the interface.
//
// Go-only: language-adapter — field access through the interface (PORTING.md §4).
func (t *SymbolTypeBase) GetSymbolTypeBase() *SymbolTypeBase { return t }

// Java: SymbolType.getTypeClass()
func (t *SymbolTypeBase) GetTypeClass() JSTypeClass {
	return JSTypeClass_SYMBOL
}

// Java: SymbolType.testForEquality(JSType)
func (t *SymbolTypeBase) TestForEquality(that JSType) base.Tri {
	result := t.ValueTypeBase.TestForEquality(that)
	if result != jsType_TRI_NULL {
		return result
	}

	if that.CanCastTo(t.GetNativeType(JSTypeNative_SYMBOL_TYPE)) ||
		that.CanCastTo(t.GetNativeType(JSTypeNative_SYMBOL_OBJECT_TYPE)) {
		return base.Tri_UNKNOWN
	}
	return base.Tri_FALSE
}

// Java: SymbolType.isSymbolValueType()
func (t *SymbolTypeBase) IsSymbolValueType() bool {
	return true
}

// Java: SymbolType.matchesNumberContext()
func (t *SymbolTypeBase) MatchesNumberContext() bool {
	return false
}

// Java: SymbolType.matchesStringContext()
func (t *SymbolTypeBase) MatchesStringContext() bool {
	return false
}

// Java: SymbolType.matchesSymbolContext()
func (t *SymbolTypeBase) MatchesSymbolContext() bool {
	return true
}

// Java: SymbolType.matchesObjectContext()
func (t *SymbolTypeBase) MatchesObjectContext() bool {
	return true
}

// Java: SymbolType.getDisplayName()
func (t *SymbolTypeBase) GetDisplayName() *string {
	return &jsType_displayName_symbol
}

var jsType_displayName_symbol = "symbol"

// Java: SymbolType.autoboxesTo()
func (t *SymbolTypeBase) AutoboxesTo() JSType {
	return t.GetNativeType(JSTypeNative_SYMBOL_OBJECT_TYPE)
}

// Java: SymbolType.getPossibleToBooleanOutcomes()
func (t *SymbolTypeBase) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_TRUE
}

// Java: SymbolType.visit(Visitor<T>)
func (t *SymbolTypeBase) Visit(visitor Visitor) any {
	return visitor.CaseSymbolType(t.self.(SymbolType))
}
