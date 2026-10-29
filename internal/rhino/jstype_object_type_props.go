// Port of com/google/javascript/rhino/jstype/ObjectType.java @ closure-compiler v20261021

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

// GetPropertyMap returns the property map that manages the set of
// properties for an object.
//
// Java: ObjectType.getPropertyMap()
func (o *ObjectTypeBase) GetPropertyMap() *PropertyMap {
	return PropertyMap_ImmutableEmptyMap()
}

// GetSlotByKey is the default getSlot implementation. This gets overridden
// by FunctionType for lazily-resolved prototypes.
//
// Java: ObjectType.getSlot(Property.Key)
func (o *ObjectTypeBase) GetSlotByKey(name Property_Key) *Property {
	_, prop := o.selfObj().GetPropertyMap().findClosestPairByKey(name)
	return prop
}

// Java: ObjectType.getSlot(String)
func (o *ObjectTypeBase) GetSlot(name string) *Property {
	return o.selfObj().GetSlotByKey(propertyKeyLookup(name))
}

// Java: ObjectType.getOwnSlot(String)
func (o *ObjectTypeBase) GetOwnSlot(name string) *Property {
	return o.selfObj().GetPropertyMap().GetOwnProperty(name)
}

// Java: ObjectType.getOwnSlot(Property.Key)
func (o *ObjectTypeBase) GetOwnSlotByKey(name Property_Key) *Property {
	return o.selfObj().GetPropertyMap().GetOwnPropertyByKey(name)
}

// GetOwnSlotByID looks up an own property slot by its interned string ID.
func (o *ObjectTypeBase) GetOwnSlotByID(id uint32) *Property {
	return o.selfObj().GetPropertyMap().GetOwnPropertyByID(rhinoStringID(id))
}

// GetClosestDefiningType returns the closest ancestor that defines the
// property including this type itself.
//
// Java: ObjectType.getClosestDefiningType(String)
func (o *ObjectTypeBase) GetClosestDefiningType(propertyName string) ObjectType {
	owner, _ := o.selfObj().GetPropertyMap().findClosestPairByKey(propertyKeyLookup(propertyName))
	return owner
}

// Java: ObjectType.findClosestDefinition(Property.Key)
func (o *ObjectTypeBase) FindClosestDefinitionByKey(propertyName Property_Key) *Property_OwnedProperty {
	return o.selfObj().GetPropertyMap().FindClosestByKey(propertyName)
}

func propertyKeyForNameAndNode(propertyName string, propertyNode Node) Property_Key {
	if !propertyNode.IsNil() {
		rec := propertyNode.rec()
		if rec.kindFlags&kindMask == uint8(node_kindString) {
			if id := rhinoStringID(rec.payload); id != 0 && rhinoStringFromID(id) == propertyName {
				return Property_Key{stringID: id}
			}
		}
	}
	return Property_Key{stringID: rhinoStringPoolIntern(propertyName)}
}

// DefineDeclaredProperty defines a property whose type is explicitly
// declared by the programmer.
//
// propertyName: the property's name. type: the type. propertyNode: the node
// corresponding to the declaration of property which might later be
// accessed using getPropertyNode.
//
// Java: ObjectType.defineDeclaredProperty(String, JSType, Node)
func (o *ObjectTypeBase) DefineDeclaredProperty(propertyName string, type_ JSType, propertyNode Node) bool {
	return o.DefineDeclaredPropertyByKey(propertyKeyForNameAndNode(propertyName, propertyNode), type_, propertyNode)
}

// Java: ObjectType.defineDeclaredProperty(Property.Key, JSType, Node)
func (o *ObjectTypeBase) DefineDeclaredPropertyByKey(propertyName Property_Key, type_ JSType, propertyNode Node) bool {
	self := o.selfObj()
	result := self.DefinePropertyByKey(propertyName, type_ /* inferred= */, false, propertyNode)
	if propertyName.symbol == nil {
		// All property definitions go through this method or
		// defineInferredProperty. Because the properties defined an an object
		// can affect subtyping, it's slightly more efficient to register this
		// after defining the property.
		if propertyName.stringID != 0 {
			o.GetRegistry().registerPropertyOnTypeByID(propertyName.stringID, self)
		} else {
			o.GetRegistry().RegisterPropertyOnType("", self)
		}
	}
	return result
}

// DefineSynthesizedProperty defines a property whose type is on a
// synthesized object. These objects don't actually exist in the user's
// program. They're just used for bookkeeping in the type system.
//
// Java: ObjectType.defineSynthesizedProperty(String, JSType, Node)
func (o *ObjectTypeBase) DefineSynthesizedProperty(propertyName string, type_ JSType, propertyNode Node) bool {
	return o.selfObj().DefinePropertyByKey(propertyKeyForNameAndNode(propertyName, propertyNode), type_, false, propertyNode)
}

// DefineInferredPropertyByKey defines a property whose type is inferred.
//
// Java: ObjectType.defineInferredProperty(Property.Key, JSType, Node)
func (o *ObjectTypeBase) DefineInferredPropertyByKey(propertyName Property_Key, type_ JSType, propertyNode Node) bool {
	// Java: checkArgument(propertyName.string() != null, ...); string() throws
	// UnsupportedOperationException for a symbol key and is never null
	// otherwise.
	if propertyName.symbol != nil {
		panic(jutil.NewUnsupportedOperationException())
	}
	if propertyName.stringID != 0 {
		return o.defineInferredPropertyWithKey(propertyName, type_, propertyNode)
	}
	return o.DefineInferredProperty("", type_, propertyNode)
}

// DefineInferredProperty defines a property whose type is inferred.
//
// propertyName: the property's name. type: the type. propertyNode: the node
// corresponding to the inferred definition of property that might later be
// accessed using getPropertyNode.
//
// Java: ObjectType.defineInferredProperty(String, JSType, Node)
func (o *ObjectTypeBase) DefineInferredProperty(propertyName string, type_ JSType, propertyNode Node) bool {
	return o.defineInferredPropertyWithKey(propertyKeyForNameAndNode(propertyName, propertyNode), type_, propertyNode)
}

func (o *ObjectTypeBase) defineInferredPropertyWithKey(key Property_Key, type_ JSType, propertyNode Node) bool {
	self := o.selfObj()
	if o.self().HasPropertyByKey(key) {
		if o.IsPropertyTypeDeclaredByKey(key) {
			// We never want to hide a declared property with an inferred
			// property.
			return true
		}
		originalType := self.GetPropertyTypeByKey(key)
		if originalType == nil {
			failJSTypeNull()
		}
		type_ = originalType.GetLeastSupertype(type_)
	}
	// TODO(b/140764208): verify that if isResolved() then type.isResolved().
	// Defining unresolved properties on resolved types is dangerous because
	// the property type may never be resolved.

	result := self.DefinePropertyByKey(key, type_, true, propertyNode)

	// All property definitions go through this method or
	// defineDeclaredProperty. Because the properties defined an an object can
	// affect subtyping, it's slightly more efficient to register this after
	// defining the property.
	o.GetRegistry().registerPropertyOnTypeByID(key.stringID, self)

	return result
}

// DefineProperty defines a property.
//
// For clarity, callers should prefer #defineDeclaredProperty and
// #defineInferredProperty.
//
// propertyName: the property's name. type: the type. inferred: true if this
// property's type is inferred. propertyNode: the node that represents the
// definition of property. Depending on the actual sub-type the node type
// might be different. The general idea is to have an estimate of where in
// the source code this property is defined.
//
// Returns true if the property was registered successfully, false if this
// conflicts with a previous property type declaration.
//
// Java: ObjectType.defineProperty(String, JSType, boolean, Node)
func (o *ObjectTypeBase) DefineProperty(propertyName string, type_ JSType, inferred bool, propertyNode Node) bool {
	return o.selfObj().DefinePropertyByKey(NewProperty_StringKey(propertyName), type_, inferred, propertyNode)
}

// GetPropertyNode gets the node corresponding to the definition of the
// specified property. This could be the node corresponding to declaration of
// the property or the node corresponding to the first reference to this
// property, e.g., "this.propertyName" in a constructor. Note this is mainly
// intended to be an estimate of where in the source code a property is
// defined. Sometime the returned node is not even part of the global AST but
// in the AST of the JsDoc that defines a type.
//
// Returns the Node corresponding to the property or null.
//
// Java: ObjectType.getPropertyNode(String)
func (o *ObjectTypeBase) GetPropertyNode(propertyName string) Node {
	p := o.selfObj().GetSlotByKey(propertyKeyLookup(propertyName))
	if p == nil {
		return NilNode
	}
	return p.GetNode()
}

// Java: ObjectType.getPropertyNode(Property.Key)
func (o *ObjectTypeBase) GetPropertyNodeByKey(propertyName Property_Key) Node {
	p := o.selfObj().GetSlotByKey(propertyName)
	if p == nil {
		return NilNode
	}
	return p.GetNode()
}

// Java: ObjectType.getPropertyDefSite(String)
func (o *ObjectTypeBase) GetPropertyDefSite(propertyName string) Node {
	p := o.selfObj().GetSlotByKey(propertyKeyLookup(propertyName))
	if p == nil {
		return NilNode
	}
	return p.GetNode()
}

// Java: ObjectType.getPropertyDefSite(Property.Key)
func (o *ObjectTypeBase) GetPropertyDefSiteByKey(propertyName Property_Key) Node {
	p := o.selfObj().GetSlotByKey(propertyName)
	if p == nil {
		return NilNode
	}
	return p.GetNode()
}

// Java: ObjectType.getPropertyJSDocInfo(String)
func (o *ObjectTypeBase) GetPropertyJSDocInfo(propertyName string) *JSDocInfo {
	return o.GetPropertyJSDocInfoByKey(propertyKeyLookup(propertyName))
}

// Java: ObjectType.getPropertyJSDocInfo(Property.Key)
func (o *ObjectTypeBase) GetPropertyJSDocInfoByKey(propertyName Property_Key) *JSDocInfo {
	p := o.selfObj().GetSlotByKey(propertyName)
	if p == nil {
		return nil
	}
	return p.GetJSDocInfo()
}

// GetOwnPropertyJSDocInfo gets the docInfo on the specified property on
// this type. This should not be implemented recursively, as you generally
// need to know exactly on which type in the prototype chain the JSDocInfo
// exists.
//
// Java: ObjectType.getOwnPropertyJSDocInfo(String)
func (o *ObjectTypeBase) GetOwnPropertyJSDocInfo(propertyName string) *JSDocInfo {
	return o.GetOwnPropertyJSDocInfoByKey(propertyKeyLookup(propertyName))
}

// Java: ObjectType.getOwnPropertyJSDocInfo(Property.Key)
func (o *ObjectTypeBase) GetOwnPropertyJSDocInfoByKey(propertyName Property_Key) *JSDocInfo {
	p := o.selfObj().GetPropertyMap().GetOwnPropertyByKey(propertyName)
	if p == nil {
		return nil
	}
	return p.GetJSDocInfo()
}

// Java: ObjectType.getOwnPropertyDefSite(String)
func (o *ObjectTypeBase) GetOwnPropertyDefSite(propertyName string) Node {
	return o.GetOwnPropertyDefSiteByKey(propertyKeyLookup(propertyName))
}

// Java: ObjectType.getOwnPropertyDefSite(Property.Key)
func (o *ObjectTypeBase) GetOwnPropertyDefSiteByKey(propertyName Property_Key) Node {
	p := o.selfObj().GetPropertyMap().GetOwnPropertyByKey(propertyName)
	if p == nil {
		return NilNode
	}
	return p.GetNode()
}

// SetPropertyJSDocInfo sets the docInfo for the specified property from the
// JSDocInfo on its definition. info may be null.
//
// Java: ObjectType.setPropertyJSDocInfo(String, JSDocInfo)
func (o *ObjectTypeBase) SetPropertyJSDocInfo(propertyName string, info *JSDocInfo) {
	o.selfObj().SetPropertyJSDocInfoByKey(NewProperty_StringKey(propertyName), info)
}

// Java: ObjectType.setPropertyJSDocInfo(Property.Key, JSDocInfo)
func (o *ObjectTypeBase) SetPropertyJSDocInfoByKey(propertyName Property_Key, info *JSDocInfo) {
	// by default, do nothing
}

// Java: ObjectType.setPropertyNode(Property.Key, Node)
func (o *ObjectTypeBase) SetPropertyNodeByKey(propertyName Property_Key, defSite Node) {
	// by default, do nothing
}

// Java: ObjectType.findPropertyTypeWithoutConsideringTemplateTypes(Property.Key)
func (o *ObjectTypeBase) FindPropertyTypeWithoutConsideringTemplateTypes(propertyName Property_Key) JSType {
	self := o.selfObj()
	if _, isFn := self.(FunctionType); isFn || o.objFlags&objKindMask == objKindFallback || o.objFlags&objKindMask == objKindTemplatized {
		if self.HasPropertyByKey(propertyName) {
			return self.GetPropertyTypeByKey(propertyName)
		}
		return nil
	}
	if self.IsEmptyType() || self.IsUnknownType() {
		return o.GetPropertyTypeFromSlot(self.GetSlotByKey(propertyName))
	}
	if slot := self.GetSlotByKey(propertyName); slot != nil {
		return slot.GetType()
	}
	return nil
}

// HasPropertyByID reports whether t has a property with interned string ID id,
// without re-looking up id in the string pool.
func HasPropertyByID(t JSType, id uint32) bool {
	return t.HasPropertyByKey(NewProperty_StringKeyByID(id))
}

// GetPropertyKindByID returns the HasPropertyKind of t for interned string ID id.
func GetPropertyKindByID(t JSType, id uint32) JSType_HasPropertyKind {
	return t.GetPropertyKindByKey(NewProperty_StringKeyByID(id))
}

// FindPropertyTypeByID coerces t to an ObjectType and looks up the property
// type for interned string ID id, returning nil if absent.
func FindPropertyTypeByID(t JSType, id uint32) JSType {
	return t.FindPropertyTypeByKey(NewProperty_StringKeyByID(id))
}

// GetPropertyType gets the property type of the property whose name is
// given. If the underlying object does not have this property, the Unknown
// type is returned to indicate that no information is available on this
// property.
//
// This gets overridden by FunctionType for lazily-resolved call() and
// bind() functions.
//
// Returns the property's type or UnknownType. This method never returns
// null.
//
// Java: ObjectType.getPropertyType(String)
func (o *ObjectTypeBase) GetPropertyType(propertyName string) JSType {
	return o.selfObj().GetPropertyTypeByKey(propertyKeyLookup(propertyName))
}

// Java: ObjectType.getPropertyType(Property.Key)
func (o *ObjectTypeBase) GetPropertyTypeByKey(propertyName Property_Key) JSType {
	slot := o.selfObj().GetSlotByKey(propertyName)
	return o.GetPropertyTypeFromSlot(slot)
}

// Java: ObjectType.GetPropertyTypeFromSlot(StaticTypedSlot)
func (o *ObjectTypeBase) GetPropertyTypeFromSlot(slot *Property) JSType {
	if slot == nil {
		if o.self().IsNoResolvedType() || o.self().IsCheckedUnknownType() {
			return o.GetNativeType(JSTypeNative_CHECKED_UNKNOWN_TYPE)
		} else if o.self().IsEmptyType() {
			return o.GetNativeType(JSTypeNative_NO_TYPE)
		}
		return o.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	}
	return slot.GetType()
}

// Java: ObjectType.getPropertyKind(Property.Key, boolean)
func (o *ObjectTypeBase) GetPropertyKindByKeyWithAutobox(propertyName Property_Key, autobox bool) JSType_HasPropertyKind {
	// Unknown types have all properties.
	return JSType_HasPropertyKind_Of(o.self().IsEmptyType() || o.self().IsUnknownType() || o.selfObj().GetSlotByKey(propertyName) != nil)
}

// Java: ObjectType.getOwnPropertyKind(Property.Key)
func (o *ObjectTypeBase) GetOwnPropertyKindByKey(propertyName Property_Key) JSType_HasPropertyKind {
	if o.selfObj().GetPropertyMap().GetOwnPropertyByKey(propertyName) != nil {
		return JSType_HasPropertyKind_KNOWN_PRESENT
	}
	return JSType_HasPropertyKind_ABSENT
}

// HasOwnProperty checks whether the property whose name is given is present
// directly on the object. Returns false even if it is declared on a
// supertype.
//
// Java: ObjectType.hasOwnProperty(String)
func (o *ObjectTypeBase) HasOwnProperty(propertyName string) bool {
	return o.selfObj().GetPropertyMap().GetOwnProperty(propertyName) != nil
}

// Java: ObjectType.hasOwnProperty(Property.Key)
func (o *ObjectTypeBase) HasOwnPropertyByKey(propertyName Property_Key) bool {
	return o.selfObj().GetPropertyMap().GetOwnPropertyByKey(propertyName) != nil
}

// Java: ObjectType.getOwnPropertyKnownSymbols()
func (o *ObjectTypeBase) GetOwnPropertyKnownSymbols() []*KnownSymbolType {
	return o.selfObj().GetPropertyMap().GetOwnKnownSymbols()
}

// GetOwnPropertyNames returns the names of all the properties directly on
// this type.
//
// Overridden by FunctionType to add "prototype".
//
// Java: ObjectType.getOwnPropertyNames()
func (o *ObjectTypeBase) GetOwnPropertyNames() []string {
	// TODO(sdh): ObjectType specifies that this should include prototype
	// properties, but currently it does not. Check if this is a constructor
	// and add them, but this could possibly break things so it should be done
	// separately.
	return o.selfObj().GetPropertyMap().GetOwnPropertyNames()
}

// ForEachOwnPropertyName calls fn for each own property name directly on this type.
func (o *ObjectTypeBase) ForEachOwnPropertyName(fn func(string)) {
	o.selfObj().GetPropertyMap().ForEachOwnPropertyName(fn)
}

// GetOwnPropertyKeys returns the names of all the properties directly on
// this type.
//
// Java: ObjectType.getOwnPropertyKeys()
func (o *ObjectTypeBase) GetOwnPropertyKeys() []Property_Key {
	self := o.selfObj()
	ownSymbols := self.GetPropertyMap().GetOwnKnownSymbols()
	if _, isFn := self.(FunctionType); !isFn && o.objFlags&objKindMask != objKindFallback {
		pm := self.GetPropertyMap()
		pm.ensureSorted()
		entries := pm.entries()
		out := make([]Property_Key, len(entries)+len(ownSymbols))
		for i, prop := range entries {
			out[i] = NewProperty_StringKeyByID(uint32(prop.nameID))
		}
		for i, symbol := range ownSymbols {
			out[len(entries)+i] = NewProperty_SymbolKey(symbol)
		}
		return out
	}
	ownNames := self.GetOwnPropertyNames()
	out := make([]Property_Key, len(ownNames)+len(ownSymbols))
	for i, propertyName := range ownNames {
		out[i] = NewProperty_StringKey(propertyName)
	}
	for i, symbol := range ownSymbols {
		out[len(ownNames)+i] = NewProperty_SymbolKey(symbol)
	}
	return out
}

// IsPropertyTypeInferred checks whether the property's type is inferred.
//
// Java: ObjectType.isPropertyTypeInferred(String)
func (o *ObjectTypeBase) IsPropertyTypeInferred(propertyName string) bool {
	return o.IsPropertyTypeInferredByKey(propertyKeyLookup(propertyName))
}

// Java: ObjectType.isPropertyTypeInferred(Property.Key)
func (o *ObjectTypeBase) IsPropertyTypeInferredByKey(propertyName Property_Key) bool {
	slot := o.selfObj().GetSlotByKey(propertyName)
	if slot == nil {
		return false
	}
	return slot.IsTypeInferred()
}

// IsPropertyTypeDeclared checks whether the property's type is declared.
//
// Java: ObjectType.isPropertyTypeDeclared(String)
func (o *ObjectTypeBase) IsPropertyTypeDeclared(propertyName string) bool {
	return o.IsPropertyTypeDeclaredByKey(propertyKeyLookup(propertyName))
}

// Java: ObjectType.isPropertyTypeDeclared(Property.Key)
func (o *ObjectTypeBase) IsPropertyTypeDeclaredByKey(propertyName Property_Key) bool {
	slot := o.selfObj().GetSlotByKey(propertyName)
	if slot == nil {
		return false
	}
	return !slot.IsTypeInferred()
}

// HasOwnDeclaredProperty reports whether the given property is declared on
// this object.
//
// Java: ObjectType.hasOwnDeclaredProperty(String)
func (o *ObjectTypeBase) HasOwnDeclaredProperty(name string) bool {
	return o.HasOwnProperty(name) && o.IsPropertyTypeDeclared(name)
}

// Java: ObjectType.hasOwnDeclaredProperty(Property.Key)
func (o *ObjectTypeBase) HasOwnDeclaredPropertyByKey(name Property_Key) bool {
	return o.HasOwnPropertyByKey(name) && o.IsPropertyTypeDeclaredByKey(name)
}

// IsPropertyInExterns checks whether the property was defined in the
// externs.
//
// Java: ObjectType.isPropertyInExterns(String)
func (o *ObjectTypeBase) IsPropertyInExterns(propertyName string) bool {
	p := o.selfObj().GetSlotByKey(propertyKeyLookup(propertyName))
	if p == nil {
		return false
	}
	return p.IsFromExterns()
}

// GetPropertiesCount gets the number of properties of this object.
//
// Java: ObjectType.getPropertiesCount()
func (o *ObjectTypeBase) GetPropertiesCount() int {
	return o.selfObj().GetPropertyMap().GetPropertiesCount()
}

// GetPropertyNames returns a list of properties defined or inferred on this
// type and any of its supertypes.
//
// Java: ObjectType.getPropertyNames()
func (o *ObjectTypeBase) GetPropertyNames() []string {
	return o.selfObj().GetPropertyMap().GetAllKeys().StringKeys()
}

// GetAllKeys returns a set of all KnownSymbolTypes properties defined on
// this type and any of its supertypes.
//
// Java: ObjectType.getAllKeys()
func (o *ObjectTypeBase) GetAllKeys() *PropertyMap_AllKeys {
	return o.selfObj().GetPropertyMap().GetAllKeys()
}
