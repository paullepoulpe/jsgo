// Port of com/google/javascript/rhino/jstype/EnumType.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// EnumType is an enum type representing a branded collection of elements.
// Each element is referenced by its name, and has an EnumElementType type.
type EnumType struct {
	PrototypeObjectTypeBase

	// the type of the individual elements
	elementType *EnumElementType
	// the elements' interned string IDs (zero GC pointers inside the backing array)
	elementIDs []rhinoStringID
	// the node representing rhs of the enum
	source       Node
	googModuleID rhinoStringID
}

// Java: EnumType.TYPE_CLASS
const EnumType_TYPE_CLASS = JSTypeClass_ENUM

// Java: EnumType.builder(JSTypeRegistry)
func EnumType_Builder_(registry *JSTypeRegistry) *EnumType_Builder {
	return NewEnumType_Builder(registry)
}

// EnumType_Builder is EnumType.Builder.
type EnumType_Builder struct {
	PrototypeObjectType_Builder

	elementType   JSType
	source        Node
	elementNameID rhinoStringID
	googModuleID  rhinoStringID
}

// Java: EnumType.Builder(JSTypeRegistry)
func NewEnumType_Builder(registry *JSTypeRegistry) *EnumType_Builder {
	b := &EnumType_Builder{}
	b.InitPrototypeObjectType_Builder(registry)
	return b
}

// Java: EnumType.Builder.setName(String)
func (b *EnumType_Builder) SetName(x *string) *EnumType_Builder {
	b.classNameID = rhinoStringPoolIntern("enum{" + jutil.ToJavaString(x) + "}")
	b.elementNameID = rhinoStringIDFromPtr(x)
	return b
}

// Java: PrototypeObjectType.Builder.setImplicitPrototype(ObjectType)
func (b *EnumType_Builder) SetImplicitPrototype(x ObjectType) *EnumType_Builder {
	b.implicitPrototype = x
	return b
}

// Java: PrototypeObjectType.Builder.setNative(boolean)
func (b *EnumType_Builder) SetNative(x bool) *EnumType_Builder {
	b.nativeType = x
	return b
}

// Java: PrototypeObjectType.Builder.setAnonymous(boolean)
func (b *EnumType_Builder) SetAnonymous(x bool) *EnumType_Builder {
	b.anonymousType = x
	return b
}

// Java: PrototypeObjectType.Builder.setTemplateTypeMap(TemplateTypeMap)
func (b *EnumType_Builder) SetTemplateTypeMap(x *TemplateTypeMap) *EnumType_Builder {
	b.templateTypeMap = x
	return b
}

// Java: PrototypeObjectType.Builder.setTemplateParamCount(int)
func (b *EnumType_Builder) SetTemplateParamCount(x int) *EnumType_Builder {
	b.templateParamCount = x
	return b
}

// SetSource sets the object literal that creates the enum, a reference to
// another enum, or null.
//
// Java: EnumType.Builder.setSource(Node)
func (b *EnumType_Builder) SetSource(x Node) *EnumType_Builder {
	b.source = x
	return b
}

// SetGoogModuleId sets the ID of the goog.module in which this enum was
// declared.
//
// Java: EnumType.Builder.setGoogModuleId(String)
func (b *EnumType_Builder) SetGoogModuleId(x *string) *EnumType_Builder {
	b.googModuleID = rhinoStringIDFromPtr(x)
	return b
}

// SetOptionalGoogModuleId sets or clears the goog.module ID without allocating a *string.
func (b *EnumType_Builder) SetOptionalGoogModuleId(id string, ok bool) *EnumType_Builder {
	if !ok {
		b.googModuleID = 0
	} else {
		b.googModuleID = rhinoStringPoolIntern(id)
	}
	return b
}

// SetElementType sets the base type of the individual elements.
//
// Java: EnumType.Builder.setElementType(JSType)
func (b *EnumType_Builder) SetElementType(x JSType) *EnumType_Builder {
	b.elementType = x
	return b
}

// Java: EnumType.Builder.build()
func (b *EnumType_Builder) Build() *EnumType {
	return NewEnumType(b)
}

// Java: EnumType(Builder)
func NewEnumType(builder *EnumType_Builder) *EnumType {
	t := &EnumType{}
	t.objFlags = objKindEnum
	t.InitPrototypeObjectTypeBase(t, &builder.PrototypeObjectType_Builder)
	t.elementType = newEnumElementTypeByID(builder.registry, builder.elementType, builder.elementNameID, t)
	t.source = builder.source
	t.googModuleID = builder.googModuleID

	t.GetRegistry().GetResolver().ResolveIfClosed(t, EnumType_TYPE_CLASS)
	return t
}

// Java: EnumType.getTypeClass()
func (t *EnumType) GetTypeClass() JSTypeClass {
	return EnumType_TYPE_CLASS
}

// Java: JSType.isEnumType()
func (t *EnumType) IsEnumType() bool {
	return true
}

// Java: EnumType.toMaybeEnumType()
func (t *EnumType) ToMaybeEnumType() *EnumType {
	return t
}

// Java: EnumType.getImplicitPrototype()
func (t *EnumType) GetImplicitPrototype() ObjectType {
	return t.GetRegistry().GetNativeObjectType(JSTypeNative_OBJECT_TYPE)
}

// GetElements gets the elements defined on this enum. The returned set is
// immutable.
//
// Java: EnumType.getElements()
func (t *EnumType) GetElements() []string {
	out := make([]string, len(t.elementIDs))
	for i, id := range t.elementIDs {
		out[i] = rhinoStringFromID(id)
	}
	return out
}

// DefineElement defines a new element on this enum.
//
// name: the name of the new element. definingNode: the Node that defines
// this new element. Returns true iff the new element is added successfully
//
// Java: EnumType.defineElement(String, Node)
func (t *EnumType) DefineElement(name string, definingNode Node) bool {
	id := rhinoStringPoolIntern(name)
	if !slices.Contains(t.elementIDs, id) {
		t.elementIDs = append(t.elementIDs, id)
	}
	return t.DefineDeclaredPropertyByKey(NewProperty_StringKeyByID(uint32(id)), t.elementType, definingNode)
}

// GetElementsType gets the elements' type, which is a subtype of the
// enumerated type.
//
// Java: EnumType.getElementsType()
func (t *EnumType) GetElementsType() *EnumElementType {
	return t.elementType
}

// GetEnumeratedTypeOfEnumObject gets the enumerated type.
//
// Java: EnumType.getEnumeratedTypeOfEnumObject()
func (t *EnumType) GetEnumeratedTypeOfEnumObject() JSType {
	return t.elementType.GetPrimitiveType()
}

// Java: EnumType.testForEquality(JSType)
func (t *EnumType) TestForEquality(that JSType) base.Tri {
	result := t.PrototypeObjectTypeBase.TestForEquality(that)
	if result != jsType_TRI_NULL {
		return result
	}
	if t.Equals(that) {
		return base.Tri_TRUE
	}
	return base.Tri_FALSE
}

// Java: EnumType.appendTo(TypeStringBuilder)
func (t *EnumType) AppendTo(sb *TypeStringBuilder) {
	if sb.IsForAnnotations() {
		sb.Append("!Object")
	} else {
		sb.Append(jutil.ToJavaString(t.GetReferenceName()))
	}
}

// Java: EnumType.getDisplayName()
func (t *EnumType) GetDisplayName() *string {
	return t.elementType.GetDisplayName()
}

// Java: EnumType.getConstructor()
func (t *EnumType) GetConstructor() FunctionType {
	return nil
}

// Java: EnumType.matchesNumberContext()
func (t *EnumType) MatchesNumberContext() bool {
	return false
}

// Java: EnumType.matchesStringContext()
func (t *EnumType) MatchesStringContext() bool {
	return true
}

// Java: EnumType.matchesObjectContext()
func (t *EnumType) MatchesObjectContext() bool {
	return true
}

// Java: EnumType.getSource()
func (t *EnumType) GetSource() Node {
	return t.source
}

// Java: EnumType.getGoogModuleId()
func (t *EnumType) GetGoogModuleId() *string {
	return rhinoStringPtrFromID(t.googModuleID)
}

// Java: EnumType.resolveInternal(ErrorReporter)
func (t *EnumType) ResolveInternal(reporter ErrorReporter) JSType {
	t.elementType = t.elementType.Resolve(reporter).(*EnumElementType)
	return t.PrototypeObjectTypeBase.ResolveInternal(reporter)
}
