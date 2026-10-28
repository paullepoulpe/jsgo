// Port of com/google/javascript/rhino/jstype/CanCastToVisitor.java @ closure-compiler v20261021

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

// canCastTo reports whether values of thisType can be cast to thatType.
//
// Java: CanCastToVisitor + JSType.visit(RelationshipVisitor, JSType)
func canCastTo(thisType JSType, thatType JSType) bool {
	switch x := thisType.(type) {
	case *UnknownType, *NoType, *NoObjectTypeBase, *AllType, *TemplateType:
		return true
	case *BigIntType, *BooleanType, *NullType, *NumberType, *StringType,
		*SymbolTypeBase, *KnownSymbolType, *VoidType,
		*FunctionTypeBase,
		*InstanceObjectType, *PrototypeObjectTypeBase, *RecordType, *EnumType, *NoResolvedType:
		return canCastToHelper(thisType, thatType)
	case *UnionType:
		return canCastUnionFrom(x, thatType)
	case *TemplatizedType:
		// TODO(johnlenz): once the templated type work is finished, restrict the
		// type parameters.
		return canCastTo(x.GetReferencedTypeInternal(), thatType)
	case *EnumElementType:
		return canCastTo(x.GetPrimitiveType(), thatType)
	case *NamedType:
		return canCastTo(x.referencedType, thatType)
	case *ProxyObjectTypeBase:
		return canCastTo(x.referencedType, thatType)
	default:
		if _, ok := thisType.(NoObjectType); ok {
			return true // TODO(johnlenz): restrict to objects
		}
		if fn, ok := thisType.(FunctionType); ok {
			return canCastToHelper(fn, thatType)
		}
		if proxy, ok := thisType.(ProxyObjectType); ok {
			return canCastTo(proxy.GetReferencedTypeInternal(), thatType)
		}
		if obj, ok := thisType.(ObjectType); ok {
			return canCastToHelper(obj, thatType)
		}
		if val, ok := thisType.(ValueType); ok {
			return canCastToHelper(val, thatType)
		}
		panic(jutil.NewUnsupportedOperationException())
	}
}

// Java: CanCastToVisitor.canCastToUnion(JSType, UnionType)
func canCastToUnion(thisType JSType, unionType *UnionType) bool {
	for _, type_ := range unionType.GetAlternates() {
		if canCastTo(thisType, type_) {
			return true
		}
	}
	return false
}

// Java: CanCastToVisitor.canCastToFunction(JSType, FunctionType)
func canCastToFunction(thisType JSType, functionType FunctionType) bool {
	if thisType.IsFunctionType() {
		// TODO(johnlenz): visit function parts
		return true
	}
	return thisType.IsSubtypeOf(functionType) ||
		functionType.IsSubtype(thisType)
}

// Java: CanCastToVisitor.isInterface(JSType)
func canCastIsInterface(type_ JSType) bool {
	objType := type_.ToObjectType()
	if objType != nil {
		constructor := objType.GetConstructor()
		return constructor != nil && constructor.IsInterface()
	}
	return false
}

// Java: CanCastToVisitor.castCastToHelper(JSType, JSType)
func canCastToHelper(thisType JSType, thatType JSType) bool {
	if thisType.IsNoResolvedType() {
		return true
	} else if thatType.IsUnknownType() ||
		thatType.IsAllType() ||
		thatType.IsNoObjectType() || // TODO(johnlenz): restrict to objects
		thatType.IsNoResolvedType() ||
		thatType.IsNoType() {
		return true
	} else if thisType.IsRecordType() || thatType.IsRecordType() {
		return true // TODO(johnlenz): are there any misuses we can catch?
	} else if canCastIsInterface(thisType) || canCastIsInterface(thatType) {
		return true // TODO(johnlenz): are there any misuses we can catch?
	} else if thatType.IsEnumElementType() {
		return canCastTo(thisType, thatType.ToMaybeEnumElementType().GetPrimitiveType())
	} else if thatType.IsUnionType() {
		return canCastToUnion(thisType, thatType.ToMaybeUnionType())
	} else if thatType.IsFunctionType() {
		return canCastToFunction(thisType, thatType.ToMaybeFunctionType())
	} else if thatType.IsTemplatizedType() {
		// TODO(johnlenz): once the templated type work is finished, restrict
		// the type parameters.
		return canCastTo(thisType, thatType.ToMaybeTemplatizedType().GetReferencedTypeInternal())
	}

	return thisType.IsSubtypeOf(thatType) || thatType.IsSubtypeOf(thisType)
}

// Java: CanCastToVisitor.caseUnionType(UnionType, JSType)
func canCastUnionFrom(thisType *UnionType, thatType JSType) bool {
	visited := false
	for _, type_ := range thisType.GetAlternates() {
		if type_.IsVoidType() || type_.IsNullType() {
			// Don't allow if the only match between the types is null or void,
			// otherwise any nullable type would be castable to any other nullable
			// type and we don't want that.
		} else {
			visited = true
			if canCastTo(type_, thatType) {
				return true
			}
		}
	}

	// Special case the "null|undefined" union and allow it to be cast to any
	// type allowing either null|undefined.
	if !visited {
		nullType := thisType.GetNativeType(JSTypeNative_NULL_TYPE)
		voidType := thisType.GetNativeType(JSTypeNative_VOID_TYPE)
		return canCastTo(nullType, thatType) || canCastTo(voidType, thatType)
	}

	return false
}
