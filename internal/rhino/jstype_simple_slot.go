// Port of com/google/javascript/rhino/jstype/SimpleSlot.java @ closure-compiler v20261021

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

// SimpleSlot is the minimum implementation of StaticTypedSlot.
//
// Java: class SimpleSlot (subclassed by jscomp's SymbolTable.Symbol).
type SimpleSlot interface {
	StaticTypedSlot
}

// SimpleSlotBase holds SimpleSlot's fields and methods; plain slots are
// *SimpleSlotBase values. The base methods make no virtual calls, so no self
// field is needed.
type SimpleSlotBase struct {
	name     string
	type_    JSType
	inferred bool
}

// NewSimpleSlot constructs a plain SimpleSlot.
//
// Java: SimpleSlot(String, JSType, boolean)
func NewSimpleSlot(name string, type_ JSType, inferred bool) *SimpleSlotBase {
	s := &SimpleSlotBase{}
	s.InitSimpleSlotBase(name, type_, inferred)
	return s
}

// InitSimpleSlotBase is the Go form of the constructor for subclasses.
//
// Java: SimpleSlot(String, JSType, boolean)
func (s *SimpleSlotBase) InitSimpleSlotBase(name string, type_ JSType, inferred bool) {
	s.name = name
	s.type_ = type_
	s.inferred = inferred
}

// Java: SimpleSlot.getName()
func (s *SimpleSlotBase) GetName() string {
	return s.name
}

// Java: SimpleSlot.getType()
func (s *SimpleSlotBase) GetType() JSType {
	return s.type_
}

// Java: SimpleSlot.isTypeInferred()
func (s *SimpleSlotBase) IsTypeInferred() bool {
	return s.inferred
}

// Java: SimpleSlot.getDeclaration()
func (s *SimpleSlotBase) GetDeclaration() StaticRef {
	return nil
}

// Java: SimpleSlot.getJSDocInfo()
func (s *SimpleSlotBase) GetJSDocInfo() *JSDocInfo {
	return nil
}

// Java: SimpleSlot.getScope()
func (s *SimpleSlotBase) GetScope() StaticScope {
	return nil
}
