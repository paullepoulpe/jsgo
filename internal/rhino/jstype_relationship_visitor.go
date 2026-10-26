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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// RelationshipVisitor is a generic visitor for relationships between types:
// the visited type is `this` and the second argument is `that`.
//
// Java: interface RelationshipVisitor<T>; dispatched via
// VisitJSTypeRelationship since Go methods cannot declare their own type
// parameters.
type RelationshipVisitor[T any] interface {
	// CaseUnknownType is the unknown type's case.
	CaseUnknownType(thisType JSType, thatType JSType) T

	// CaseNoType is the bottom type's case.
	CaseNoType(thatType JSType) T

	// CaseNoObjectType is the bottom Object type's case.
	CaseNoObjectType(thatType JSType) T

	// CaseAllType is the all type's case.
	CaseAllType(thatType JSType) T

	// CaseValueType is the value types' case.
	CaseValueType(thisType ValueType, thatType JSType) T

	// CaseObjectType is the object type's case.
	CaseObjectType(thisType ObjectType, thatType JSType) T

	// CaseFunctionType is the function type's case.
	CaseFunctionType(thisType FunctionType, thatType JSType) T

	// CaseUnionType is the union type's case.
	CaseUnionType(thisType *UnionType, thatType JSType) T

	// CaseTemplatizedType is the templatized type's case.
	CaseTemplatizedType(thisType *TemplatizedType, thatType JSType) T

	// CaseTemplateType is the template type's case.
	CaseTemplateType(thisType *TemplateType, thatType JSType) T

	// CaseEnumElementType is the enum element type's case.
	CaseEnumElementType(typeType *EnumElementType, thatType JSType) T
}

// VisitJSTypeRelationship dispatches t and that to visitor's matching case
// method.
//
// Java: JSType.visit(RelationshipVisitor<T>, JSType)
func VisitJSTypeRelationship[T any](t JSType, visitor RelationshipVisitor[T], that JSType) T {
	switch x := t.(type) {
	case *UnknownType:
		return visitor.CaseUnknownType(x, that)
	case *NoType:
		return visitor.CaseNoType(that)
	case *NoObjectTypeBase:
		return visitor.CaseNoObjectType(that)
	case *AllType:
		return visitor.CaseAllType(that)
	case *BigIntType:
		return visitor.CaseValueType(x, that)
	case *BooleanType:
		return visitor.CaseValueType(x, that)
	case *NullType:
		return visitor.CaseValueType(x, that)
	case *NumberType:
		return visitor.CaseValueType(x, that)
	case *StringType:
		return visitor.CaseValueType(x, that)
	case *SymbolTypeBase:
		return visitor.CaseValueType(x, that)
	case *KnownSymbolType:
		return visitor.CaseValueType(x, that)
	case *VoidType:
		return visitor.CaseValueType(x, that)
	case *FunctionTypeBase:
		return visitor.CaseFunctionType(x, that)
	case *InstanceObjectType:
		return visitor.CaseObjectType(x, that)
	case *PrototypeObjectTypeBase:
		return visitor.CaseObjectType(x, that)
	case *RecordType:
		return visitor.CaseObjectType(x, that)
	case *EnumType:
		return visitor.CaseObjectType(x, that)
	case *NoResolvedType:
		return visitor.CaseObjectType(x, that)
	case *UnionType:
		return visitor.CaseUnionType(x, that)
	case *TemplatizedType:
		return visitor.CaseTemplatizedType(x, that)
	case *TemplateType:
		return visitor.CaseTemplateType(x, that)
	case *EnumElementType:
		return visitor.CaseEnumElementType(x, that)
	case *NamedType:
		return VisitJSTypeRelationship(x.referencedType, visitor, that)
	case *ProxyObjectTypeBase:
		return VisitJSTypeRelationship(x.referencedType, visitor, that)
	default:
		if _, ok := t.(NoObjectType); ok {
			return visitor.CaseNoObjectType(that)
		}
		if fn, ok := t.(FunctionType); ok {
			return visitor.CaseFunctionType(fn, that)
		}
		if proxy, ok := t.(ProxyObjectType); ok {
			return VisitJSTypeRelationship(proxy.GetReferencedTypeInternal(), visitor, that)
		}
		if obj, ok := t.(ObjectType); ok {
			return visitor.CaseObjectType(obj, that)
		}
		if val, ok := t.(ValueType); ok {
			return visitor.CaseValueType(val, that)
		}
		panic(jutil.NewUnsupportedOperationException())
	}
}
