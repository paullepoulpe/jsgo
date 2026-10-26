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

// CanCastToVisitor is a "can cast to" relationship visitor.
//
// Java: implements RelationshipVisitor<Boolean>
type CanCastToVisitor struct{}

// Java: CanCastToVisitor()
func NewCanCastToVisitor() *CanCastToVisitor {
	return &CanCastToVisitor{}
}

// Java: CanCastToVisitor.caseUnknownType(JSType, JSType)
func (v *CanCastToVisitor) CaseUnknownType(thisType JSType, thatType JSType) bool {
	return true
}

// Java: CanCastToVisitor.caseNoType(JSType)
func (v *CanCastToVisitor) CaseNoType(thatType JSType) bool {
	return true
}

// Java: CanCastToVisitor.caseNoObjectType(JSType)
func (v *CanCastToVisitor) CaseNoObjectType(thatType JSType) bool {
	return true // TODO(johnlenz): restrict to objects
}

// Java: CanCastToVisitor.caseAllType(JSType)
func (v *CanCastToVisitor) CaseAllType(thatType JSType) bool {
	return true
}

// Java: CanCastToVisitor.canCastToUnion(JSType, UnionType)
func (v *CanCastToVisitor) CanCastToUnion(thisType JSType, unionType *UnionType) bool {
	for _, type_ := range unionType.GetAlternates() {
		if VisitJSTypeRelationship[bool](thisType, v, type_) {
			return true
		}
	}
	return false
}

// Java: CanCastToVisitor.canCastToFunction(JSType, FunctionType)
func (v *CanCastToVisitor) CanCastToFunction(thisType JSType, functionType FunctionType) bool {
	if thisType.IsFunctionType() {
		// TODO(johnlenz): visit function parts
		return true
	} else {
		return thisType.IsSubtypeOf(functionType) ||
			functionType.IsSubtype(thisType)
	}
}

// Java: CanCastToVisitor.isInterface(JSType)
func (v *CanCastToVisitor) IsInterface(type_ JSType) bool {
	objType := type_.ToObjectType()
	if objType != nil {
		constructor := objType.GetConstructor()
		return constructor != nil && constructor.IsInterface()
	}
	return false
}

// Java: CanCastToVisitor.castCastToHelper(JSType, JSType)
func (v *CanCastToVisitor) CastCastToHelper(thisType JSType, thatType JSType) bool {
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
	} else if v.IsInterface(thisType) || v.IsInterface(thatType) {
		return true // TODO(johnlenz): are there any misuses we can catch?
	} else if thatType.IsEnumElementType() {
		return VisitJSTypeRelationship[bool](thisType, v,
			thatType.ToMaybeEnumElementType().GetPrimitiveType())
	} else if thatType.IsUnionType() {
		return v.CanCastToUnion(thisType, thatType.ToMaybeUnionType())
	} else if thatType.IsFunctionType() {
		return v.CanCastToFunction(thisType, thatType.ToMaybeFunctionType())
	} else if thatType.IsTemplatizedType() {
		// TODO(johnlenz): once the templated type work is finished, restrict
		// the type parameters.
		return VisitJSTypeRelationship[bool](thisType, v,
			thatType.ToMaybeTemplatizedType().GetReferencedTypeInternal())
	}

	return thisType.IsSubtypeOf(thatType) || thatType.IsSubtypeOf(thisType)
}

// Java: CanCastToVisitor.caseValueType(ValueType, JSType)
func (v *CanCastToVisitor) CaseValueType(thisType ValueType, thatType JSType) bool {
	return v.CastCastToHelper(thisType, thatType)
}

// Java: CanCastToVisitor.caseObjectType(ObjectType, JSType)
func (v *CanCastToVisitor) CaseObjectType(thisType ObjectType, thatType JSType) bool {
	return v.CastCastToHelper(thisType, thatType)
}

// Java: CanCastToVisitor.caseFunctionType(FunctionType, JSType)
func (v *CanCastToVisitor) CaseFunctionType(thisType FunctionType, thatType JSType) bool {
	return v.CastCastToHelper(thisType, thatType)
}

// Java: CanCastToVisitor.caseUnionType(UnionType, JSType)
func (v *CanCastToVisitor) CaseUnionType(thisType *UnionType, thatType JSType) bool {
	visited := false
	for _, type_ := range thisType.GetAlternates() {
		if type_.IsVoidType() || type_.IsNullType() {
			// Don't allow if the only match between the types is null or void,
			// otherwise any nullable type would be castable to any other nullable
			// type and we don't want that.
		} else {
			visited = true
			if VisitJSTypeRelationship[bool](type_, v, thatType) {
				return true
			}
		}
	}

	// Special case the "null|undefined" union and allow it to be cast to any
	// cast to any type containing allowing either null|undefined.
	if !visited {
		NULL_TYPE := thisType.GetNativeType(JSTypeNative_NULL_TYPE)
		VOID_TYPE := thisType.GetNativeType(JSTypeNative_VOID_TYPE)
		return VisitJSTypeRelationship[bool](NULL_TYPE, v, thatType) || VisitJSTypeRelationship[bool](VOID_TYPE, v, thatType)
	}

	return false
}

// Java: CanCastToVisitor.caseTemplatizedType(TemplatizedType, JSType)
func (v *CanCastToVisitor) CaseTemplatizedType(thisType *TemplatizedType, thatType JSType) bool {
	// TODO(johnlenz): once the templated type work is finished, restrict the
	// type parameters.
	return VisitJSTypeRelationship[bool](thisType.GetReferencedTypeInternal(), v, thatType)
}

// Java: CanCastToVisitor.caseTemplateType(TemplateType, JSType)
func (v *CanCastToVisitor) CaseTemplateType(thisType *TemplateType, thatType JSType) bool {
	return true
}

// Java: CanCastToVisitor.caseEnumElementType(EnumElementType, JSType)
func (v *CanCastToVisitor) CaseEnumElementType(typeType *EnumElementType, thatType JSType) bool {
	return VisitJSTypeRelationship[bool](typeType.GetPrimitiveType(), v, thatType)
}
