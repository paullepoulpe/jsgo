// Port of com/google/javascript/rhino/jstype/ValueType.java @ closure-compiler v20261015

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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// ValueType is the class of value types (null, void, number, boolean,
// string).
//
// Java: abstract class ValueType extends JSType.
type ValueType interface {
	JSType

	// GetValueTypeBase gives access to the base through the interface.
	//
	// Go-only: language-adapter — field access through the interface
	// (PORTING.md §4).
	GetValueTypeBase() *ValueTypeBase
}

// ValueTypeBase holds ValueType's base method implementations; the virtual
// calls go through JSTypeBase.self.
type ValueTypeBase struct {
	JSTypeBase
}

// InitValueTypeBase is the Go form of the ValueType constructor.
//
// Java: ValueType(JSTypeRegistry)
func (v *ValueTypeBase) InitValueTypeBase(self ValueType, registry *JSTypeRegistry) {
	v.InitJSTypeBase(self, registry)
	v.EagerlyResolveToSelf()
}

// GetValueTypeBase gives access to the base through the interface.
//
// Go-only: language-adapter — field access through the interface (PORTING.md §4).
func (v *ValueTypeBase) GetValueTypeBase() *ValueTypeBase { return v }

// Java: ValueType.resolveInternal(ErrorReporter)
func (v *ValueTypeBase) ResolveInternal(reporter ErrorReporter) JSType {
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: ValueType.appendTo(TypeStringBuilder)
func (v *ValueTypeBase) AppendTo(sb *TypeStringBuilder) {
	sb.Append(jutil.ToJavaString(v.self.GetDisplayName()))
}

// Java: ValueType.hasDisplayName()
func (v *ValueTypeBase) HasDisplayName() bool {
	return true
}

// Java: ValueType.visit(RelationshipVisitor<T>, JSType)
func (v *ValueTypeBase) VisitRelationship(visitor RelationshipVisitor, that JSType) any {
	return visitor.CaseValueType(v.self.(ValueType), that)
}

// Java: ValueType.getPropertyKind(Property.Key, boolean)
func (v *ValueTypeBase) GetPropertyKindByKeyWithAutobox(propertyName Property_Key, autobox bool) JSType_HasPropertyKind {
	if autobox && v.self.IsBoxableScalar() {
		return v.self.AutoboxesTo().GetPropertyKindByKeyWithAutobox(propertyName, autobox)
	} else {
		return v.JSTypeBase.GetPropertyKindByKeyWithAutobox(propertyName, autobox)
	}
}

// Java: ValueType.recursionUnsafeHashCode()
func (v *ValueTypeBase) RecursionUnsafeHashCode() int32 {
	// Subclasses of this type are unique within a JSTypeRegistry.
	return jutil.IdentityHashCode(v.self)
}
