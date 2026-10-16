// Port of com/google/javascript/rhino/jstype/RelationshipVisitor.java
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
 *   John Lenz
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
 * ***** END LICENSE BLOCK ***** */

package jstype

// RelationshipVisitor calculates a value of type T from two JSTypes
// based on the structural relationship between them.
//
// Java uses a generic type parameter T on the interface.
// Go interfaces cannot carry generic methods, so each case method
// returns `any`; callers type-assert the concrete result.
//
// Port of RelationshipVisitor<T> — T is erased to any.
// @author johnlenz@google.com (John Lenz)
type RelationshipVisitor interface {
	// CaseUnknownType is the unknown type's case.
	CaseUnknownType(thisType JSType, thatType JSType) any

	// CaseNoType is the bottom type's case.
	CaseNoType(thatType JSType) any

	// CaseNoObjectType is the bottom object type's case.
	CaseNoObjectType(thatType JSType) any

	// CaseAllType is the all type's case.
	CaseAllType(thatType JSType) any

	// CaseValueType is the value type's case.
	// TODO: replace JSType with *ValueType once ValueType is ported.
	CaseValueType(thisType JSType, thatType JSType) any

	// CaseObjectType is the object type's case.
	CaseObjectType(thisType *ObjectType, thatType JSType) any

	// CaseFunctionType is the function type's case.
	CaseFunctionType(thisType *FunctionType, thatType JSType) any

	// CaseUnionType is the union type's case.
	CaseUnionType(thisType *UnionType, thatType JSType) any

	// CaseTemplatizedType is the templatized type's case.
	CaseTemplatizedType(thisType *TemplatizedType, thatType JSType) any

	// CaseTemplateType is the template type's case.
	CaseTemplateType(thisType *TemplateType, thatType JSType) any

	// CaseEnumElementType is the enum element type's case.
	CaseEnumElementType(thisType *EnumElementType, thatType JSType) any
}
