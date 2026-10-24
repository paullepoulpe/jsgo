// Port of com/google/javascript/rhino/jstype/RecordType.java @ closure-compiler v20261015

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

// RecordType is a record (structural) type.
//
// Subtyping: The subtyping of a record type is defined via structural
// comparison of a record type's properties. For example, a record type of
// the form { a : TYPE_1 } is a supertype of a record type of the form { b :
// TYPE_2, a : TYPE_1 } because B can be assigned to A and matches all
// constraints. Similarly, a defined type can be assigned to a record type so
// long as that defined type matches all property constraints of the record
// type. A record type of the form { a : A, b : B } can be assigned to a
// record of type { a : A }.
type RecordType struct {
	PrototypeObjectTypeBase

	declared bool
	isFrozen bool
}

// Java: RecordType.TYPE_CLASS
const RecordType_TYPE_CLASS = JSTypeClass_RECORD

// NewRecordType creates a record type.
//
// registry: the type registry under which this type lives. properties: a
// map of all the properties of this record type. declared: whether this is a
// declared or synthesized type. A synthesized record type is just used for
// bookkeeping in the type system. A declared record type was actually used
// in the user's program.
//
// Throws IllegalStateException if the RecordProperty associated with a
// property is null.
//
// Java: RecordType(JSTypeRegistry, Map<String, RecordProperty>, boolean)
func NewRecordType(registry *JSTypeRegistry, properties *jutil.TreeMap[string, *RecordTypeBuilder_RecordProperty], declared bool) *RecordType {
	t := &RecordType{}
	t.objFlags = objKindRecord
	builderFields := PrototypeObjectType_BuilderFields{
		registry:        registry,
		templateTypeMap: registry.GetEmptyTemplateTypeMap(),
	}
	t.InitPrototypeObjectTypeBase(t, &builderFields)
	t.SetPrettyPrint(true)
	t.declared = declared

	for property := range properties.Keys() {
		prop := properties.Get(property)
		if prop == nil {
			panic(jutil.NewIllegalStateExceptionMsg(
				"RecordProperty associated with a property should not be null!"))
		}
		if declared {
			t.DefineDeclaredProperty(property, prop.GetType(), prop.GetPropertyNode())
		} else {
			t.DefineSynthesizedProperty(property, prop.GetType(), prop.GetPropertyNode())
		}
	}
	// Freeze the record type.
	t.isFrozen = true

	registry.GetResolver().ResolveIfClosed(t, RecordType_TYPE_CLASS)
	return t
}

// Java: RecordType.getTypeClass()
func (t *RecordType) GetTypeClass() JSTypeClass {
	return RecordType_TYPE_CLASS
}

// IsSynthetic reports whether this is synthesized for internal bookkeeping.
//
// Java: RecordType.isSynthetic()
func (t *RecordType) IsSynthetic() bool {
	return !t.declared
}

// Java: RecordType.getImplicitPrototype()
func (t *RecordType) GetImplicitPrototype() ObjectType {
	return t.GetRegistry().GetNativeObjectType(JSTypeNative_OBJECT_TYPE)
}

// Java: RecordType.defineProperty(Property.Key, JSType, boolean, Node)
func (t *RecordType) DefinePropertyByKey(propertyName Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	if t.isFrozen {
		return false
	}

	return t.PrototypeObjectTypeBase.DefinePropertyByKey(propertyName, type_, inferred, propertyNode)
}

// Java: RecordType.getGreatestSubtypeHelper(JSType)
func (t *RecordType) GetGreatestSubtypeHelper(that JSType) JSType {
	if that.IsRecordType() {
		thatRecord := that.ToMaybeRecordType()
		builder := NewRecordTypeBuilder(t.GetRegistry())
		builder.SetSynthesized(true)

		noType := t.GetRegistry().GetNativeObjectType(JSTypeNative_NO_TYPE)

		// The greatest subtype consists of those *unique* properties of both
		// record types. If any property conflicts, then the NO_TYPE type is
		// returned.
		for _, property := range t.GetOwnPropertyNames() {
			thisPropertyType := t.GetPropertyType(property)
			var propType JSType
			if thatRecord.HasProperty(property) {
				thatPropertyType := thatRecord.GetPropertyType(property)
				propType = thisPropertyType.GetGreatestSubtype(thatPropertyType)
				if propType.Equals(noType) {
					return noType
				}
			} else {
				propType = thisPropertyType
			}
			builder.AddProperty(property, propType, t.GetPropertyNode(property))
		}

		for _, property := range thatRecord.GetOwnPropertyNames() {
			if !t.HasProperty(property) {
				builder.AddProperty(property, thatRecord.GetPropertyType(property),
					thatRecord.GetPropertyNode(property))
			}
		}
		return builder.Build()
	}

	greatestSubtype := t.GetRegistry().GetNativeType(JSTypeNative_NO_OBJECT_TYPE)
	thatRestrictedToObj := t.GetRegistry().GetNativeType(JSTypeNative_OBJECT_TYPE).GetGreatestSubtype(that)
	if !thatRestrictedToObj.IsEmptyType() {
		// In this branch, the other type is some object type. We find the
		// greatest subtype with the following algorithm:
		// 1) For each property "x" of this record type, take the union of all
		//    classes with a property "x" with a compatible property type. and
		//    which are a subtype of that.
		// 2) Take the intersection of all of these unions.
		for _, propName := range t.GetOwnPropertyNames() {
			propType := t.GetPropertyType(propName)
			builder := UnionType_Builder_(t.GetRegistry())
			for _, alt := range t.GetRegistry().GetEachReferenceTypeWithProperty(propName) {
				altPropType := alt.GetPropertyType(propName)
				if altPropType != nil &&
					!alt.Equals(t) &&
					alt.IsSubtypeOf(that) &&
					altPropType.IsSubtypeOf(propType) {
					builder.AddAlternate(alt)
				}
			}
			greatestSubtype = greatestSubtype.GetLeastSupertype(builder.Build())
		}
	}
	return greatestSubtype
}

// Java: JSType.isRecordType()
func (t *RecordType) IsRecordType() bool {
	return true
}

// Java: RecordType.toMaybeRecordType()
func (t *RecordType) ToMaybeRecordType() *RecordType {
	return t
}

// Java: RecordType.isStructuralType()
func (t *RecordType) IsStructuralType() bool {
	return true
}
