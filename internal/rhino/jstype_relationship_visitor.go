// Port of com/google/javascript/rhino/jstype/RelationshipVisitor.java @ closure-compiler v20261021

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
 *   John Lenz
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

// RelationshipVisitor is a vistor for relationships between types: the
// visited type is `this` and the second argument is `that`.
//
// Go-only: language-adapter — RelationshipVisitor<T> with T erased to `any`
// (see Visitor).
type RelationshipVisitor interface {
	// CaseUnknownType is the unknown type's case.
	CaseUnknownType(thisType JSType, thatType JSType) any

	// CaseNoType is the bottom type's case.
	CaseNoType(thatType JSType) any

	// CaseNoObjectType is the bottom Object type's case.
	CaseNoObjectType(thatType JSType) any

	// CaseAllType is the all type's case.
	CaseAllType(thatType JSType) any

	// CaseValueType is the value types' case.
	CaseValueType(thisType ValueType, thatType JSType) any

	// CaseObjectType is the object type's case.
	CaseObjectType(thisType ObjectType, thatType JSType) any

	// CaseFunctionType is the function type's case.
	CaseFunctionType(thisType FunctionType, thatType JSType) any

	// CaseUnionType is the union type's case.
	CaseUnionType(thisType *UnionType, thatType JSType) any

	// CaseTemplatizedType is the templatized type's case.
	CaseTemplatizedType(thisType *TemplatizedType, thatType JSType) any

	// CaseTemplateType is the template type's case.
	CaseTemplateType(thisType *TemplateType, thatType JSType) any

	// CaseEnumElementType is the enum element type's case.
	CaseEnumElementType(typeType *EnumElementType, thatType JSType) any
}
