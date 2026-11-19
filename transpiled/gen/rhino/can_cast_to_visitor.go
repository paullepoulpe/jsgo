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
// Port of com/google/javascript/rhino/jstype/CanCastToVisitor.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.CanCastToVisitor
type CanCastToVisitor struct {
	idHash_ atomic.Int32
}

func (this_ *CanCastToVisitor) BindSelf_(s any) {
	_ = s
}
func (this_ *CanCastToVisitor) IsJavaNil_() bool { return this_ == nil }
func (*CanCastToVisitor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.CanCastToVisitor", "CanCastToVisitor"
}
func (this_ *CanCastToVisitor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CanCastToVisitor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.CanCastToVisitor", "com.google.javascript.rhino.jstype.CanCastToVisitor", func(v_ any) bool {
		return j2gort.InstanceOf[*CanCastToVisitor](v_)
	})
}
func (*CanCastToVisitor) Implements_rhino_RelationshipVisitor() {}

func (this_ *CanCastToVisitor) initFields_CanCastToVisitor() {
	_ = this_
}

// Java constructor: CanCastToVisitor()
func NewCanCastToVisitor() *CanCastToVisitor {
	this_ := &CanCastToVisitor{}
	InitCanCastToVisitor(this_)
	return this_
}

func InitCanCastToVisitor(this_ *CanCastToVisitor) {
	this_.initFields_CanCastToVisitor()
}

// Java method: CanCastToVisitor.caseUnknownType(com.google.javascript.rhino.jstype.JSType thisType, com.google.javascript.rhino.jstype.JSType thatType)
func (this_ *CanCastToVisitor) CaseUnknownType(thisType JSType, thatType JSType) any {
	_ = thisType
	_ = thatType
	return j2gort.ValueOfBoolean(true)
}

// Java method: CanCastToVisitor.caseNoType(com.google.javascript.rhino.jstype.JSType thatType)
func (this_ *CanCastToVisitor) CaseNoType(thatType JSType) any {
	_ = thatType
	return j2gort.ValueOfBoolean(true)
}

// Java method: CanCastToVisitor.caseNoObjectType(com.google.javascript.rhino.jstype.JSType thatType)
func (this_ *CanCastToVisitor) CaseNoObjectType(thatType JSType) any {
	_ = thatType
	return j2gort.ValueOfBoolean(true)
}

// Java method: CanCastToVisitor.caseAllType(com.google.javascript.rhino.jstype.JSType thatType)
func (this_ *CanCastToVisitor) CaseAllType(thatType JSType) any {
	_ = thatType
	return j2gort.ValueOfBoolean(true)
}

// Java method: CanCastToVisitor.canCastToUnion(com.google.javascript.rhino.jstype.JSType thisType, com.google.javascript.rhino.jstype.UnionType unionType)
func (this_ *CanCastToVisitor) CanCastToUnion(thisType JSType, unionType *UnionType) bool {
	_ = thisType
	_ = unionType
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = unionType.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			if j2gort.UnboxBoolean((j2gort.CastPtr[j2gort.BoxBoolean](thisType.Visit_2(this_, unionType.GetAlternateAt(i))))) {
				return true
			}
		}
	}
	return false
}

// Java method: CanCastToVisitor.canCastToFunction(com.google.javascript.rhino.jstype.JSType thisType, com.google.javascript.rhino.jstype.FunctionType functionType)
func (this_ *CanCastToVisitor) CanCastToFunction(thisType JSType, functionType *FunctionType) bool {
	_ = thisType
	_ = functionType
	if thisType.IsFunctionType() {
		return true
	} else {
		return (thisType.IsSubtypeOf_1(ToIface_JSType(functionType)) || functionType.IsSubtype_1(thisType))
	}
}

// Java method: CanCastToVisitor.isInterface(com.google.javascript.rhino.jstype.JSType type)
func (this_ *CanCastToVisitor) priv_CanCastToVisitor_isInterface_8822(type_ JSType) bool {
	_ = type_
	var objType ObjectType = type_.ToObjectType()
	_ = objType
	if objType != nil && !objType.IsJavaNil_() {
		var constructor JSType = ToIface_JSType(objType.GetConstructor())
		_ = constructor
		return ((constructor != nil && !constructor.IsJavaNil_()) && constructor.IsInterface())
	}
	return false
}

// Java method: CanCastToVisitor.castCastToHelper(com.google.javascript.rhino.jstype.JSType thisType, com.google.javascript.rhino.jstype.JSType thatType)
func (this_ *CanCastToVisitor) CastCastToHelper(thisType JSType, thatType JSType) *j2gort.BoxBoolean {
	_ = thisType
	_ = thatType
	if thisType.IsNoResolvedType() {
		return j2gort.ValueOfBoolean(true)
	} else if (((thatType.IsUnknownType() || thatType.IsAllType()) || thatType.IsNoObjectType()) || thatType.IsNoResolvedType()) || thatType.IsNoType() {
		return j2gort.ValueOfBoolean(true)
	} else if thisType.IsRecordType() || thatType.IsRecordType() {
		return j2gort.ValueOfBoolean(true)
	} else if this_.priv_CanCastToVisitor_isInterface_8822(thisType) || this_.priv_CanCastToVisitor_isInterface_8822(thatType) {
		return j2gort.ValueOfBoolean(true)
	} else if thatType.IsEnumElementType() {
		return j2gort.CastPtr[j2gort.BoxBoolean](thisType.Visit_2(this_, thatType.ToMaybeEnumElementType().GetPrimitiveType()))
	} else if thatType.IsUnionType() {
		return j2gort.ValueOfBoolean(this_.CanCastToUnion(thisType, thatType.ToMaybeUnionType()))
	} else if thatType.IsFunctionType() {
		return j2gort.ValueOfBoolean(this_.CanCastToFunction(thisType, thatType.ToMaybeFunctionType()))
	} else if thatType.IsTemplatizedType() {
		return j2gort.CastPtr[j2gort.BoxBoolean](thisType.Visit_2(this_, thatType.ToMaybeTemplatizedType().GetReferencedTypeInternal()))
	}
	return j2gort.ValueOfBoolean((thisType.IsSubtypeOf_1(thatType) || thatType.IsSubtypeOf_1(thisType)))
}

// Java method: CanCastToVisitor.caseValueType(com.google.javascript.rhino.jstype.ValueType thisType, com.google.javascript.rhino.jstype.JSType thatType)
func (this_ *CanCastToVisitor) CaseValueType(thisType ValueType, thatType JSType) any {
	_ = thisType
	_ = thatType
	return j2gort.PtrToAny(this_.CastCastToHelper(j2gort.AsIface[JSType](thisType), thatType))
}

// Java method: CanCastToVisitor.caseObjectType(com.google.javascript.rhino.jstype.ObjectType thisType, com.google.javascript.rhino.jstype.JSType thatType)
func (this_ *CanCastToVisitor) CaseObjectType(thisType ObjectType, thatType JSType) any {
	_ = thisType
	_ = thatType
	return j2gort.PtrToAny(this_.CastCastToHelper(j2gort.AsIface[JSType](thisType), thatType))
}

// Java method: CanCastToVisitor.caseFunctionType(com.google.javascript.rhino.jstype.FunctionType thisType, com.google.javascript.rhino.jstype.JSType thatType)
func (this_ *CanCastToVisitor) CaseFunctionType(thisType *FunctionType, thatType JSType) any {
	_ = thisType
	_ = thatType
	return j2gort.PtrToAny(this_.CastCastToHelper(ToIface_JSType(thisType), thatType))
}

// Java method: CanCastToVisitor.caseUnionType(com.google.javascript.rhino.jstype.UnionType thisType, com.google.javascript.rhino.jstype.JSType thatType)
func (this_ *CanCastToVisitor) CaseUnionType(thisType *UnionType, thatType JSType) any {
	_ = thisType
	_ = thatType
	var visited bool = false
	_ = visited
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = thisType.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			var type_ JSType = thisType.GetAlternateAt(i)
			_ = type_
			if type_.IsVoidType() || type_.IsNullType() {
			} else {
				visited = true
				if j2gort.UnboxBoolean((j2gort.CastPtr[j2gort.BoxBoolean](type_.Visit_2(this_, thatType)))) {
					return j2gort.ValueOfBoolean(true)
				}
			}
		}
	}
	if !(visited) {
		var l_NULL_TYPE JSType = thisType.GetNativeType(StaticGet_JSTypeNative_NULL_TYPE())
		_ = l_NULL_TYPE
		var l_VOID_TYPE JSType = thisType.GetNativeType(StaticGet_JSTypeNative_VOID_TYPE())
		_ = l_VOID_TYPE
		return j2gort.ValueOfBoolean((j2gort.UnboxBoolean(j2gort.CastPtr[j2gort.BoxBoolean](l_NULL_TYPE.Visit_2(this_, thatType))) || j2gort.UnboxBoolean(j2gort.CastPtr[j2gort.BoxBoolean](l_VOID_TYPE.Visit_2(this_, thatType)))))
	}
	return j2gort.ValueOfBoolean(false)
}

// Java method: CanCastToVisitor.caseTemplatizedType(com.google.javascript.rhino.jstype.TemplatizedType thisType, com.google.javascript.rhino.jstype.JSType thatType)
func (this_ *CanCastToVisitor) CaseTemplatizedType(thisType *TemplatizedType, thatType JSType) any {
	_ = thisType
	_ = thatType
	return j2gort.PtrToAny(j2gort.CastPtr[j2gort.BoxBoolean](thisType.GetReferencedTypeInternal().Visit_2(this_, thatType)))
}

// Java method: CanCastToVisitor.caseTemplateType(com.google.javascript.rhino.jstype.TemplateType thisType, com.google.javascript.rhino.jstype.JSType thatType)
func (this_ *CanCastToVisitor) CaseTemplateType(thisType *TemplateType, thatType JSType) any {
	_ = thisType
	_ = thatType
	return j2gort.ValueOfBoolean(true)
}

// Java method: CanCastToVisitor.caseEnumElementType(com.google.javascript.rhino.jstype.EnumElementType typeType, com.google.javascript.rhino.jstype.JSType thatType)
func (this_ *CanCastToVisitor) CaseEnumElementType(typeType *EnumElementType, thatType JSType) any {
	_ = typeType
	_ = thatType
	return j2gort.PtrToAny(j2gort.CastPtr[j2gort.BoxBoolean](typeType.GetPrimitiveType().Visit_2(this_, thatType)))
}
