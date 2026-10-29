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

import (
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ObjectType is the object type.
//
// In JavaScript, all object types have properties, and each of those
// properties has a type. Property types may be DECLARED, INFERRED, or
// UNKNOWN.
//
// DECLARED properties have an explicit type annotation, as in:
//
//	/** @type {number} */
//	Foo.prototype.bar = 1;
//
// This property may only hold number values, and an assignment to any other
// type of value is an error.
//
// INFERRED properties do not have an explicit type annotation. Rather, we
// try to find all the possible types that this property can hold.
//
//	Foo.prototype.bar = 1;
//
// If the programmer assigns other types of values to this property, the
// property will take on the union of all these types.
//
// UNKNOWN properties are properties on the UNKNOWN type. The UNKNOWN type has
// all properties, but we do not know whether they are declared or inferred.
//
// Java: abstract class ObjectType extends JSType.
type ObjectType interface {
	JSType

	GetPropertyMap() *PropertyMap
	// Java: ObjectType.getSlot(Property.Key)
	GetSlotByKey(name Property_Key) *Property
	// Java: ObjectType.getSlot(String)
	GetSlot(name string) *Property
	// Java: ObjectType.getOwnSlot(String)
	GetOwnSlot(name string) *Property
	// Java: ObjectType.getOwnSlot(Property.Key)
	GetOwnSlotByKey(name Property_Key) *Property
	GetOwnSlotByID(id uint32) *Property
	GetTypeOfThis() JSType
	GetTemplateTypes() []JSType
	SetJSDocInfo(info *JSDocInfo)
	DetectImplicitPrototypeCycle() bool
	DetectInheritanceCycle() bool
	GetReferenceName() *string
	HasReferenceName() bool
	GetNormalizedReferenceName() string
	GetRawType() ObjectType
	GetConstructor() FunctionType
	GetSuperClassConstructor() FunctionType
	GetClosestDefiningType(propertyName string) ObjectType
	FindClosestDefinitionByKey(propertyName Property_Key) *Property_OwnedProperty
	GetImplicitPrototype() ObjectType
	DefineDeclaredProperty(propertyName string, type_ JSType, propertyNode Node) bool
	DefineDeclaredPropertyByKey(propertyName Property_Key, type_ JSType, propertyNode Node) bool
	DefineSynthesizedProperty(propertyName string, type_ JSType, propertyNode Node) bool
	DefineInferredPropertyByKey(propertyName Property_Key, type_ JSType, propertyNode Node) bool
	DefineInferredProperty(propertyName string, type_ JSType, propertyNode Node) bool
	DefineProperty(propertyName string, type_ JSType, inferred bool, propertyNode Node) bool
	DefinePropertyByKey(propertyName Property_Key, type_ JSType, inferred bool, propertyNode Node) bool
	GetPropertyNode(propertyName string) Node
	GetPropertyNodeByKey(propertyName Property_Key) Node
	GetPropertyDefSite(propertyName string) Node
	GetPropertyDefSiteByKey(propertyName Property_Key) Node
	GetPropertyJSDocInfo(propertyName string) *JSDocInfo
	GetPropertyJSDocInfoByKey(propertyName Property_Key) *JSDocInfo
	GetOwnPropertyJSDocInfo(propertyName string) *JSDocInfo
	GetOwnPropertyJSDocInfoByKey(propertyName Property_Key) *JSDocInfo
	GetOwnPropertyDefSite(propertyName string) Node
	GetOwnPropertyDefSiteByKey(propertyName Property_Key) Node
	SetPropertyJSDocInfo(propertyName string, info *JSDocInfo)
	SetPropertyJSDocInfoByKey(propertyName Property_Key, info *JSDocInfo)
	SetPropertyNodeByKey(propertyName Property_Key, defSite Node)
	GetPropertyType(propertyName string) JSType
	GetPropertyTypeByKey(propertyName Property_Key) JSType
	GetOwnPropertyKindByKey(propertyName Property_Key) JSType_HasPropertyKind
	HasOwnProperty(propertyName string) bool
	HasOwnPropertyByKey(propertyName Property_Key) bool
	GetOwnPropertyKnownSymbols() []*KnownSymbolType
	GetOwnPropertyNames() []string
	ForEachOwnPropertyName(fn func(string))
	GetOwnPropertyKeys() []Property_Key
	IsPropertyTypeInferred(propertyName string) bool
	IsPropertyTypeInferredByKey(propertyName Property_Key) bool
	IsPropertyTypeDeclared(propertyName string) bool
	IsPropertyTypeDeclaredByKey(propertyName Property_Key) bool
	HasOwnDeclaredProperty(name string) bool
	HasOwnDeclaredPropertyByKey(name Property_Key) bool
	IsPropertyInExterns(propertyName string) bool
	GetPropertiesCount() int
	GetPropertyNames() []string
	GetAllKeys() *PropertyMap_AllKeys
	IsImplicitPrototypeOf(other ObjectType) bool
	HasCachedValues() bool
	ClearCachedValues()
	GetOwnerFunction() FunctionType
	SetOwnerFunction(type_ FunctionType)
	GetCtorImplementedInterfaces() []ObjectType
	GetCtorExtendedInterfaces() []ObjectType
	GetEnumeratedTypeOfEnumObject() JSType

	// GetObjectTypeBase gives access to ObjectType's fields through the
	// interface.
	//
	// Go-only: language-adapter — field access through the interface
	// (PORTING.md §4).
	GetObjectTypeBase() *ObjectTypeBase
}

// ObjectTypeBase holds the fields and base method implementations of
// ObjectType.
type ObjectTypeBase struct {
	JSTypeBase
}

func objTypeBase(x ObjectType) *ObjectTypeBase {
	if x == nil {
		return nil
	}
	return x.GetObjectTypeBase()
}

// objectTypeOf returns the ObjectType whose base is b, or nil for nil.
func objectTypeOf(b *ObjectTypeBase) ObjectType {
	if b == nil {
		return nil
	}
	return b.self_.(ObjectType)
}

func (o *ObjectTypeBase) selfObj() ObjectType {
	return o.self_.(ObjectType)
}

// InitObjectTypeBase is the Go form of the ObjectType(JSTypeRegistry)
// constructor.
//
// Java: ObjectType(JSTypeRegistry)
func (o *ObjectTypeBase) InitObjectTypeBase(self ObjectType, registry *JSTypeRegistry) {
	o.unknown = true
	o.InitJSTypeBase(self, registry)
}

// InitObjectTypeBaseWithTemplateTypeMap is the Go form of the
// ObjectType(JSTypeRegistry, TemplateTypeMap) constructor.
//
// Java: ObjectType(JSTypeRegistry, TemplateTypeMap)
func (o *ObjectTypeBase) InitObjectTypeBaseWithTemplateTypeMap(
	self ObjectType, registry *JSTypeRegistry, templateTypeMap *TemplateTypeMap) {
	o.unknown = true
	o.InitJSTypeBaseWithTemplateTypeMap(self, registry, templateTypeMap)
}

// GetObjectTypeBase gives access to the base through the interface.
//
// Go-only: language-adapter — field access through the interface (PORTING.md §4).
func (o *ObjectTypeBase) GetObjectTypeBase() *ObjectTypeBase { return o }

// Java: ObjectType.getTypeOfThis()
func (o *ObjectTypeBase) GetTypeOfThis() JSType {
	return nil
}

// GetTemplateTypes gets the declared default element type.
//
// Java: ObjectType.getTemplateTypes()
func (o *ObjectTypeBase) GetTemplateTypes() []JSType {
	return []JSType{}
}

// GetJSDocInfo gets the docInfo for this type.
//
// Java: ObjectType.getJSDocInfo()
func (o *ObjectTypeBase) GetJSDocInfo() *JSDocInfo {
	if o.objFlags&objHasDocFlag == 0 {
		return nil
	}
	return o.GetRegistry().getObjectDocInfo(o)
}

// SetJSDocInfo sets the docInfo for this type from the given JSDocInfo. The
// JSDocInfo may be null.
//
// Java: ObjectType.setJSDocInfo(JSDocInfo)
func (o *ObjectTypeBase) SetJSDocInfo(info *JSDocInfo) {
	if info == nil {
		if o.objFlags&objHasDocFlag != 0 {
			o.objFlags &^= objHasDocFlag
			o.GetRegistry().deleteObjectDocInfo(o)
		}
		return
	}
	o.objFlags |= objHasDocFlag
	o.GetRegistry().setObjectDocInfo(o, info)
}

// DetectImplicitPrototypeCycle detects a cycle in the implicit prototype
// chain. This method accesses the #getImplicitPrototype() method and must
// therefore be invoked only after the object is sufficiently initialized to
// respond to calls to this method.
//
// Returns true iff an implicit prototype cycle was detected.
//
// Java: ObjectType.detectImplicitPrototypeCycle()
func (o *ObjectTypeBase) DetectImplicitPrototypeCycle() bool {
	// detecting cycle
	self := o.selfObj()
	o.flags |= jsTypeFlagVisited
	p := self.GetImplicitPrototype()
	for p != nil {
		if p.GetObjectTypeBase().flags&jsTypeFlagVisited != 0 {
			return true
		} else {
			p.GetObjectTypeBase().flags |= jsTypeFlagVisited
		}
		p = p.GetImplicitPrototype()
	}

	// clean up
	p = self
	for {
		p.GetObjectTypeBase().flags &^= jsTypeFlagVisited
		p = p.GetImplicitPrototype()
		if p == nil {
			break
		}
	}
	return false
}

// DetectInheritanceCycle detects cycles in either the implicit prototype
// chain, or the implemented/extended interfaces.
//
// Returns true iff a cycle was detected.
//
// Java: ObjectType.detectInheritanceCycle()
func (o *ObjectTypeBase) DetectInheritanceCycle() bool {
	self := o.selfObj()
	if self.DetectImplicitPrototypeCycle() ||
		objectTypeIterablesContains(self.GetCtorImplementedInterfaces(), self) {
		return true
	}
	fnType := self.GetConstructor()
	return fnType != nil && fnType.CheckExtendsLoop() != nil
}

// objectTypeIterablesContains is Iterables.contains(iterable, element):
// Objects.equals over the elements.
func objectTypeIterablesContains(iterable []ObjectType, element ObjectType) bool {
	for _, e := range iterable {
		if jutil.ObjectsEquals(e, element) {
			return true
		}
	}
	return false
}

// HasReferenceName: INVARIANT: hasReferenceName() is true if and only if
// getReferenceName() returns a non-null string.
//
// Returns true if the object is named, false if it is anonymous.
//
// Java: ObjectType.hasReferenceName()
func (o *ObjectTypeBase) HasReferenceName() bool {
	return o.selfObj().GetReferenceName() != nil
}

func normalizeReferenceName(s string) string {
	start := jutil.StringIndexOfChar(s, '(')
	if start != -1 {
		end := jutil.StringLastIndexOfCharFrom(s, ')', jutil.StringLength(s)-1)
		prefix := jutil.StringSubstring(s, 0, start)
		// Java: `end + 1 % name.length() == 0` (% binds tighter than +).
		if end+1%jutil.StringLength(s) == 0 {
			return prefix
		}
		return prefix + jutil.StringSubstringFrom(s, end+1)
	}
	return s
}

// GetNormalizedReferenceName: due to the complexity of some of our internal
// type systems, sometimes we have different types constructed by the same
// constructor. In other parts of the type system, these are called
// delegates. We construct these types by appending suffixes to the
// constructor name.
//
// The normalized reference name does not have these suffixes, and as such,
// recollapses these implicit types back to their real type. Note that
// suffixes such as ".prototype" can be added after the delegate suffix, so
// anything after the parentheses must still be retained.
//
// Java: ObjectType.getNormalizedReferenceName()
func (o *ObjectTypeBase) GetNormalizedReferenceName() string {
	name := o.selfObj().GetReferenceName()
	if name == nil {
		return ""
	}
	return normalizeReferenceName(*name)
}

// Java: ObjectType.getDisplayName()
func (o *ObjectTypeBase) GetDisplayName() *string {
	name := o.selfObj().GetReferenceName()
	if name == nil {
		return nil
	}
	if jutil.StringIndexOfChar(*name, '(') == -1 {
		return name
	}
	norm := normalizeReferenceName(*name)
	return &norm
}

// ObjectType_CreateDelegateSuffix creates a suffix for a proxy delegate.
//
// Java: ObjectType.createDelegateSuffix(String)
func ObjectType_CreateDelegateSuffix(suffix string) string {
	return "(" + suffix + ")"
}

// Java: ObjectType.getRawType()
func (o *ObjectTypeBase) GetRawType() ObjectType {
	t := o.self().ToMaybeTemplatizedType()
	if t == nil {
		return o.selfObj()
	}
	return t.GetReferencedType()
}

// Java: ObjectType.testForEquality(JSType)
func (o *ObjectTypeBase) TestForEquality(that JSType) base.Tri {
	// super
	result := o.JSTypeBase.TestForEquality(that)
	if result != jsType_TRI_NULL {
		return result
	}

	// TODO: consider tighten "testForEquality" for subtypes of Object: if Foo
	// and Bar are not related we don't want to allow "==" on them
	// (similiarly we should disallow number == for non-number context values,
	// etc).

	if that.IsUnknownType() ||
		that.IsSubtypeOf(o.GetNativeType(JSTypeNative_OBJECT_TYPE)) ||
		that.IsSubtypeOf(o.GetNativeType(JSTypeNative_NUMBER_TYPE)) ||
		that.IsSubtypeOf(o.GetNativeType(JSTypeNative_STRING_TYPE)) ||
		that.IsSubtypeOf(o.GetNativeType(JSTypeNative_BOOLEAN_TYPE)) ||
		that.IsSubtypeOf(o.GetNativeType(JSTypeNative_SYMBOL_TYPE)) ||
		that.IsSubtypeOf(o.GetNativeType(JSTypeNative_BIGINT_TYPE)) {
		return base.Tri_UNKNOWN
	}
	return base.Tri_FALSE
}

// Java: ObjectType.getSuperClassConstructor()
func (o *ObjectTypeBase) GetSuperClassConstructor() FunctionType {
	iproto := o.selfObj().GetImplicitPrototype()
	if iproto == nil {
		return nil
	}
	iproto = iproto.GetImplicitPrototype()
	if iproto == nil {
		return nil
	}
	return iproto.GetConstructor()
}

// Java: ObjectType.isStructuralType()
func (o *ObjectTypeBase) IsStructuralType() bool {
	constructor := o.selfObj().GetConstructor()
	return constructor != nil && constructor.IsStructuralInterface()
}

// IsImplicitPrototypeOf returns whether this is on the implicit prototype
// chain of other.
//
// this need not be the immediate prototype of other.
//
// Java: ObjectType.isImplicitPrototypeOf(ObjectType)
func (o *ObjectTypeBase) IsImplicitPrototypeOf(other ObjectType) bool {
	unwrappedThisBase := objTypeBase(ObjectType_DeeplyUnwrap(o.selfObj()))
	if unwrappedThisBase == nil {
		return false
	}

	for other = ObjectType_DeeplyUnwrap(other); ; {
		ob := objTypeBase(other)
		if ob == nil {
			return false
		}
		// The prototype should match exactly. NOTE: the use of "==" here rather
		// than equals is deliberate. This method is very hot in the type
		// checker and relying on identity improves performance of both type
		// checking/type inferrence and property disambiguation.
		if unwrappedThisBase == ob {
			return true
		}
		other = ObjectType_DeeplyUnwrap(other.GetImplicitPrototype())
	}
}

// ObjectType_DeeplyUnwrap recursively unwraps all ProxyObjectTypes,
// including unwrapping the raw type out of a templatized type.
//
// Guaranteed to return a non-proxy type (exception: will return an
// unresolved NamedType).
//
// Java: ObjectType.deeplyUnwrap(ObjectType)
func ObjectType_DeeplyUnwrap(original ObjectType) ObjectType {
	current := original
	for {
		cb := objTypeBase(current)
		if cb == nil {
			return nil
		}
		switch cb.objFlags & objKindMask {
		case objKindInstance, objKindPrototype, objKindFunction, objKindRecord, objKindEnum,
			objKindUnknown, objKindNoObject, objKindNo, objKindNoResolved, objKindEnumElement:
			return current
		case objKindTemplatized:
			current = current.(*TemplatizedType).GetReferencedType()
		case objKindNamed:
			nt := current.(*NamedType)
			if !nt.IsSuccessfullyResolved() {
				return current
			}
			current = nt.referencedObjType
		case objKindTemplate, objKindProxy:
			current = current.(ProxyObjectType).GetProxyObjectTypeBase().referencedObjType
		default:
			proxy, ok := current.(ProxyObjectType)
			if !ok {
				return current
			}
			// TODO(lharker): remove this case and instead fail. Only the Rhino unit
			// tests are triggering this by creating new ProxyObjectTypes.
			current = proxy.GetReferencedObjTypeInternal()
		}
	}
}

// Java: ObjectType.getPossibleToBooleanOutcomes()
func (o *ObjectTypeBase) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_TRUE
}

// IsUnknownType: we treat this as the unknown type if any of its implicit
// prototype properties is unknown.
//
// Java: ObjectType.isUnknownType()
func (o *ObjectTypeBase) IsUnknownType() bool {
	// If the object is unknown now, check the supertype again, because it
	// might have been resolved since the last check.
	if o.unknown {
		self := o.selfObj()
		implicitProto := self.GetImplicitPrototype()
		if implicitProto == nil || implicitProto.IsNativeObjectType() {
			o.unknown = false
			for _, interfaceType := range self.GetCtorExtendedInterfaces() {
				if interfaceType.IsUnknownType() {
					o.unknown = true
					break
				}
			}
		} else {
			o.unknown = implicitProto.IsUnknownType()
		}
	}
	return o.unknown
}

// Java: ObjectType.isObject()
func (o *ObjectTypeBase) IsObject() bool {
	return true
}

// Java: JSType.isObjectType()
func (o *ObjectTypeBase) IsObjectType() bool {
	return true
}

// HasCachedValues returns true if any cached values have been set for this
// type. If true, then the prototype chain should not be changed, as it might
// invalidate the cached values.
//
// Java: ObjectType.hasCachedValues()
func (o *ObjectTypeBase) HasCachedValues() bool {
	return !o.unknown
}

// ClearCachedValues clears cached values. Should be called before making
// changes to a prototype that may have been changed since creation.
//
// Java: ObjectType.clearCachedValues()
func (o *ObjectTypeBase) ClearCachedValues() {
	o.unknown = true
}

// IsNativeObjectType reports whether this is a built-in object.
//
// Java: ObjectType.isNativeObjectType()
func (o *ObjectTypeBase) IsNativeObjectType() bool {
	return false
}

// ObjectType_Cast is a null-safe version of JSType#toObjectType.
//
// Java: ObjectType.cast(JSType)
func ObjectType_Cast(type_ JSType) ObjectType {
	if type_ == nil {
		return nil
	}
	return type_.ToObjectType()
}

// Java: ObjectType.isFunctionPrototypeType()
func (o *ObjectTypeBase) IsFunctionPrototypeType() bool {
	return o.selfObj().GetOwnerFunction() != nil
}

// Java: ObjectType.getOwnerFunction()
func (o *ObjectTypeBase) GetOwnerFunction() FunctionType {
	return nil
}

// SetOwnerFunction sets the owner function. By default, does nothing.
//
// Java: ObjectType.setOwnerFunction(FunctionType)
func (o *ObjectTypeBase) SetOwnerFunction(type_ FunctionType) {}

// GetCtorImplementedInterfaces gets the interfaces implemented by the ctor
// associated with this type. Intended to be overridden by subclasses.
//
// Java: ObjectType.getCtorImplementedInterfaces()
func (o *ObjectTypeBase) GetCtorImplementedInterfaces() []ObjectType {
	return []ObjectType{}
}

// GetCtorExtendedInterfaces gets the interfaces extended by the interface
// associated with this type. Intended to be overridden by subclasses.
//
// Java: ObjectType.getCtorExtendedInterfaces()
func (o *ObjectTypeBase) GetCtorExtendedInterfaces() []ObjectType {
	return []ObjectType{}
}

// Java: ObjectType.getEnumeratedTypeOfEnumObject()
func (o *ObjectTypeBase) GetEnumeratedTypeOfEnumObject() JSType {
	return nil
}
