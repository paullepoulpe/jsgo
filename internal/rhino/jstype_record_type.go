// Port of com/google/javascript/rhino/jstype/RecordType.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/RecordTypeBuilder.java @ closure-compiler v20261021

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

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

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
}

// Java: RecordType.TYPE_CLASS
const RecordType_TYPE_CLASS = JSTypeClass_RECORD

func newRecordTypeFromProps(registry *JSTypeRegistry, props []recordProp, declared bool) *RecordType {
	t := &RecordType{}
	t.objFlags = objKindRecord
	builder := PrototypeObjectType_Builder{
		registry:        registry,
		templateTypeMap: registry.GetEmptyTemplateTypeMap(),
	}
	t.InitPrototypeObjectTypeBase(t, &builder)
	t.SetPrettyPrint(true)
	if declared {
		t.protoFlags |= protoFlagRecordDeclared
	}

	for i := range props {
		prop := &props[i]
		if declared {
			t.DefineDeclaredProperty(prop.name, prop.type_, prop.propertyNode)
		} else {
			t.DefineSynthesizedProperty(prop.name, prop.type_, prop.propertyNode)
		}
	}
	t.protoFlags |= protoFlagRecordFrozen

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
	return (t.protoFlags & protoFlagRecordDeclared) == 0
}

// Java: RecordType.getImplicitPrototype()
func (t *RecordType) GetImplicitPrototype() ObjectType {
	return t.GetRegistry().GetNativeObjectType(JSTypeNative_OBJECT_TYPE)
}

// Java: RecordType.defineProperty(Property.Key, JSType, boolean, Node)
func (t *RecordType) DefinePropertyByKey(propertyName Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	if (t.protoFlags & protoFlagRecordFrozen) != 0 {
		return false
	}

	return t.PrototypeObjectTypeBase.DefinePropertyByKey(propertyName, type_, inferred, propertyNode)
}

// Java: RecordType.getGreatestSubtypeHelper(JSType)
func (t *RecordType) GetGreatestSubtypeHelper(that JSType) JSType {
	if that.IsRecordType() {
		thatRecord := that.ToMaybeRecordType()
		var builder RecordTypeBuilder
		builder.init(t.GetRegistry())
		builder.SetSynthesized(true)

		noType := t.GetRegistry().GetNativeObjectType(JSTypeNative_NO_TYPE)

		// The greatest subtype consists of those *unique* properties of both
		// record types. If any property conflicts, then the NO_TYPE type is
		// returned.
		thisMap := t.GetPropertyMap()
		thisMap.ensureSorted()
		for _, prop := range thisMap.entries() {
			propKey := NewProperty_StringKeyByID(uint32(prop.nameID))
			thisPropertyType := prop.GetType()
			var propType JSType
			if thatRecord.HasPropertyByKey(propKey) {
				thatPropertyType := thatRecord.GetPropertyTypeByKey(propKey)
				propType = thisPropertyType.GetGreatestSubtype(thatPropertyType)
				if propType.Equals(noType) {
					return noType
				}
			} else {
				propType = thisPropertyType
			}
			builder.AddProperty(propKey.String_(), propType, prop.GetNode())
		}

		thatMap := thatRecord.GetPropertyMap()
		thatMap.ensureSorted()
		for _, prop := range thatMap.entries() {
			propKey := NewProperty_StringKeyByID(uint32(prop.nameID))
			if !t.HasPropertyByKey(propKey) {
				builder.AddProperty(propKey.String_(), prop.GetType(), prop.GetNode())
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
		thisMap := t.GetPropertyMap()
		thisMap.ensureSorted()
		for _, prop := range thisMap.entries() {
			propID := prop.nameID
			propKey := NewProperty_StringKeyByID(uint32(propID))
			propType := prop.GetType()
			builder := UnionType_Builder_(t.GetRegistry())
			t.GetRegistry().forEachRefTypeWithPropertyByID(propID, func(alt ObjectType) {
				altPropType := alt.GetPropertyTypeByKey(propKey)
				if altPropType != nil &&
					!alt.Equals(t) &&
					alt.IsSubtypeOf(that) &&
					altPropType.IsSubtypeOf(propType) {
					builder.AddAlternate(alt)
				}
			})
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

type recordProp struct {
	name         string
	type_        JSType
	propertyNode Node
}

func compareRecordProp(a, b recordProp) int {
	return jutil.StringCompareTo(a.name, b.name)
}

// RecordTypeBuilder is a builder for record types.
type RecordTypeBuilder struct {
	registry   *JSTypeRegistry
	props      []recordProp
	propsBuf   [8]recordProp
	isEmpty    bool
	isDeclared bool
}

func (b *RecordTypeBuilder) init(registry *JSTypeRegistry) {
	b.registry = registry
	b.props = b.propsBuf[:0]
	b.isEmpty = true
	b.isDeclared = true
}

// Java: RecordTypeBuilder(JSTypeRegistry)
func NewRecordTypeBuilder(registry *JSTypeRegistry) *RecordTypeBuilder {
	b := &RecordTypeBuilder{}
	b.init(registry)
	return b
}

// SetSynthesized: see the comments on RecordType about synthetic types.
//
// Java: RecordTypeBuilder.setSynthesized(boolean)
func (b *RecordTypeBuilder) SetSynthesized(synthesized bool) {
	b.isDeclared = !synthesized
}

// AddProperty adds a property with the given name and type to the record
// type.
//
// name: the name of the new property. type: the JSType of the new property.
// propertyNode: the node that holds this property definition. Returns the
// builder itself for chaining purposes.
//
// Java: RecordTypeBuilder.addProperty(String, JSType, Node)
func (b *RecordTypeBuilder) AddProperty(name string, type_ JSType, propertyNode Node) *RecordTypeBuilder {
	b.isEmpty = false
	for i := range b.props {
		if b.props[i].name == name {
			b.props[i].type_ = type_
			b.props[i].propertyNode = propertyNode
			return b
		}
	}
	b.props = append(b.props, recordProp{name: name, type_: type_, propertyNode: propertyNode})
	return b
}

// Build creates a record.
//
// Java: RecordTypeBuilder.build()
func (b *RecordTypeBuilder) Build() JSType {
	// If we have an empty record, simply return the object type.
	if b.isEmpty {
		return b.registry.GetNativeObjectType(JSTypeNative_OBJECT_TYPE)
	}
	var sortedBuf [8]recordProp
	var sorted []recordProp
	if len(b.props) <= len(sortedBuf) {
		sorted = sortedBuf[:len(b.props)]
	} else {
		sorted = make([]recordProp, len(b.props))
	}
	copy(sorted, b.props)
	slices.SortFunc(sorted, compareRecordProp)
	return newRecordTypeFromProps(b.registry, sorted, b.isDeclared)
}

// RecordTypeBuilder_RecordProperty is RecordTypeBuilder.RecordProperty.
type RecordTypeBuilder_RecordProperty struct {
	type_        JSType
	propertyNode Node
}

// Java: RecordTypeBuilder.RecordProperty(JSType, Node)
func NewRecordTypeBuilder_RecordProperty(type_ JSType, propertyNode Node) *RecordTypeBuilder_RecordProperty {
	return &RecordTypeBuilder_RecordProperty{type_: type_, propertyNode: propertyNode}
}

// Java: RecordTypeBuilder.RecordProperty.getType()
func (p *RecordTypeBuilder_RecordProperty) GetType() JSType {
	return p.type_
}

// Java: RecordTypeBuilder.RecordProperty.getPropertyNode()
func (p *RecordTypeBuilder_RecordProperty) GetPropertyNode() Node {
	return p.propertyNode
}

// Java: RecordTypeBuilder.RecordProperty.toString()
func (p *RecordTypeBuilder_RecordProperty) String() string {
	return "RecordProperty{type: " + jutil.ToJavaString(p.type_) + ", node: " + jutil.ToJavaString(p.propertyNode) + "}"
}
