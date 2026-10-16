// Port of com/google/javascript/rhino/jstype/ObjectType.java @ closure-compiler v20260415

package rhino

import (
	"iter"

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
	GetTypeOfThis() JSType
	GetTemplateTypes() []JSType
	SetJSDocInfo(info *JSDocInfo)
	DetectImplicitPrototypeCycle() bool
	DetectInheritanceCycle() bool
	GetReferenceName() *string
	HasReferenceName() bool
	GetNormalizedReferenceName() *string
	GetRawType() ObjectType
	GetConstructor() FunctionType
	GetSuperClassConstructor() FunctionType
	GetClosestDefiningType(propertyName string) ObjectType
	FindClosestDefinition(propertyName string) *Property_OwnedProperty
	FindClosestDefinitionByKey(propertyName Property_Key) *Property_OwnedProperty
	GetImplicitPrototype() ObjectType
	GetImplicitPrototypeChain() iter.Seq[ObjectType]
	DefineDeclaredProperty(propertyName string, type_ JSType, propertyNode *Node) bool
	DefineDeclaredPropertyByKey(propertyName Property_Key, type_ JSType, propertyNode *Node) bool
	DefineSynthesizedProperty(propertyName string, type_ JSType, propertyNode *Node) bool
	DefineInferredPropertyByKey(propertyName Property_Key, type_ JSType, propertyNode *Node) bool
	DefineInferredProperty(propertyName string, type_ JSType, propertyNode *Node) bool
	DefineProperty(propertyName string, type_ JSType, inferred bool, propertyNode *Node) bool
	DefinePropertyByKey(propertyName Property_Key, type_ JSType, inferred bool, propertyNode *Node) bool
	GetPropertyNode(propertyName string) *Node
	GetPropertyNodeByKey(propertyName Property_Key) *Node
	GetPropertyDefSite(propertyName string) *Node
	GetPropertyDefSiteByKey(propertyName Property_Key) *Node
	GetPropertyJSDocInfo(propertyName string) *JSDocInfo
	GetPropertyJSDocInfoByKey(propertyName Property_Key) *JSDocInfo
	GetOwnPropertyJSDocInfo(propertyName string) *JSDocInfo
	GetOwnPropertyJSDocInfoByKey(propertyName Property_Key) *JSDocInfo
	GetOwnPropertyDefSite(propertyName string) *Node
	GetOwnPropertyDefSiteByKey(propertyName Property_Key) *Node
	SetPropertyJSDocInfo(propertyName string, info *JSDocInfo)
	SetPropertyJSDocInfoByKey(propertyName Property_Key, info *JSDocInfo)
	SetPropertyNode(propertyName string, defSite *Node)
	SetPropertyNodeByKey(propertyName Property_Key, defSite *Node)
	GetPropertyType(propertyName string) JSType
	GetPropertyTypeByKey(propertyName Property_Key) JSType
	GetOwnPropertyKind(propertyName string) JSType_HasPropertyKind
	GetOwnPropertyKindByKey(propertyName Property_Key) JSType_HasPropertyKind
	HasOwnProperty(propertyName string) bool
	HasOwnPropertyByKey(propertyName Property_Key) bool
	GetOwnPropertyKnownSymbols() []*KnownSymbolType
	GetOwnPropertyNames() []string
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
	GetPropertyTypeMap() *jutil.LinkedHashMap[string, JSType]
	GetEnumeratedTypeOfEnumObject() JSType

	// GetObjectTypeBase gives access to ObjectType's fields through the
	// interface.
	//
	// Go-only: language-adapter — field access through the interface
	// (PORTING.md §4).
	GetObjectTypeBase() *ObjectTypeBase
}

// ObjectTypeBase holds the fields and base method implementations of
// ObjectType; self (shadowing JSTypeBase.self) dispatches ObjectType's
// virtual methods.
type ObjectTypeBase struct {
	JSTypeBase
	self ObjectType

	visited bool
	docInfo *JSDocInfo
	unknown bool
}

// InitObjectTypeBase is the Go form of the ObjectType(JSTypeRegistry)
// constructor.
//
// Java: ObjectType(JSTypeRegistry)
func (o *ObjectTypeBase) InitObjectTypeBase(self ObjectType, registry *JSTypeRegistry) {
	o.self = self
	o.unknown = true
	o.InitJSTypeBase(self, registry)
}

// InitObjectTypeBaseWithTemplateTypeMap is the Go form of the
// ObjectType(JSTypeRegistry, TemplateTypeMap) constructor.
//
// Java: ObjectType(JSTypeRegistry, TemplateTypeMap)
func (o *ObjectTypeBase) InitObjectTypeBaseWithTemplateTypeMap(
	self ObjectType, registry *JSTypeRegistry, templateTypeMap *TemplateTypeMap) {
	o.self = self
	o.unknown = true
	o.InitJSTypeBaseWithTemplateTypeMap(self, registry, templateTypeMap)
}

// GetObjectTypeBase gives access to the base through the interface.
//
// Go-only: language-adapter — field access through the interface (PORTING.md §4).
func (o *ObjectTypeBase) GetObjectTypeBase() *ObjectTypeBase { return o }

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
	property := o.self.GetPropertyMap().FindClosestByKey(name)
	if property == nil {
		return nil
	}
	return property.GetValue()
}

// Java: ObjectType.getSlot(String)
func (o *ObjectTypeBase) GetSlot(name string) *Property {
	return o.self.GetSlotByKey(NewProperty_StringKey(name))
}

// Java: ObjectType.getOwnSlot(String)
func (o *ObjectTypeBase) GetOwnSlot(name string) *Property {
	return o.self.GetPropertyMap().GetOwnProperty(name)
}

// Java: ObjectType.getOwnSlot(Property.Key)
func (o *ObjectTypeBase) GetOwnSlotByKey(name Property_Key) *Property {
	return o.self.GetPropertyMap().GetOwnPropertyByKey(name)
}

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
	return o.docInfo
}

// SetJSDocInfo sets the docInfo for this type from the given JSDocInfo. The
// JSDocInfo may be null.
//
// Java: ObjectType.setJSDocInfo(JSDocInfo)
func (o *ObjectTypeBase) SetJSDocInfo(info *JSDocInfo) {
	o.docInfo = info
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
	o.visited = true
	p := o.self.GetImplicitPrototype()
	for p != nil {
		if p.GetObjectTypeBase().visited {
			return true
		} else {
			p.GetObjectTypeBase().visited = true
		}
		p = p.GetImplicitPrototype()
	}

	// clean up
	p = o.self
	for {
		p.GetObjectTypeBase().visited = false
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
	if o.self.DetectImplicitPrototypeCycle() ||
		objectTypeIterablesContains(o.self.GetCtorImplementedInterfaces(), o.self) {
		return true
	}
	fnType := o.self.GetConstructor()
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
	return o.self.GetReferenceName() != nil
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
func (o *ObjectTypeBase) GetNormalizedReferenceName() *string {
	name := o.self.GetReferenceName()
	if name != nil {
		start := jutil.StringIndexOfChar(*name, '(')
		if start != -1 {
			end := jutil.StringLastIndexOfCharFrom(*name, ')', jutil.StringLength(*name)-1)
			prefix := jutil.StringSubstring(*name, 0, start)
			// Java: `end + 1 % name.length() == 0` (% binds tighter than +).
			var result string
			if end+1%jutil.StringLength(*name) == 0 {
				result = prefix
			} else {
				result = prefix + jutil.StringSubstringFrom(*name, end+1)
			}
			return &result
		}
	}
	return name
}

// Java: ObjectType.getDisplayName()
func (o *ObjectTypeBase) GetDisplayName() *string {
	return o.self.GetNormalizedReferenceName()
}

// ObjectType_CreateDelegateSuffix creates a suffix for a proxy delegate.
//
// Java: ObjectType.createDelegateSuffix(String)
func ObjectType_CreateDelegateSuffix(suffix string) string {
	return "(" + suffix + ")"
}

// Java: ObjectType.getRawType()
func (o *ObjectTypeBase) GetRawType() ObjectType {
	t := o.self.ToMaybeTemplatizedType()
	if t == nil {
		return o.self
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
	iproto := o.self.GetImplicitPrototype()
	if iproto == nil {
		return nil
	}
	iproto = iproto.GetImplicitPrototype()
	if iproto == nil {
		return nil
	}
	return iproto.GetConstructor()
}

// GetClosestDefiningType returns the closest ancestor that defines the
// property including this type itself.
//
// Java: ObjectType.getClosestDefiningType(String)
func (o *ObjectTypeBase) GetClosestDefiningType(propertyName string) ObjectType {
	property := o.self.GetPropertyMap().FindClosestByKey(NewProperty_StringKey(propertyName))
	if property == nil {
		return nil
	}
	return property.GetOwner()
}

// FindClosestDefinition returns the closest definition of the property
// including this type itself.
//
// Java: ObjectType.findClosestDefinition(String)
func (o *ObjectTypeBase) FindClosestDefinition(propertyName string) *Property_OwnedProperty {
	return o.self.FindClosestDefinitionByKey(NewProperty_StringKey(propertyName))
}

// Java: ObjectType.findClosestDefinition(Property.Key)
func (o *ObjectTypeBase) FindClosestDefinitionByKey(propertyName Property_Key) *Property_OwnedProperty {
	return o.self.GetPropertyMap().FindClosestByKey(propertyName)
}

// GetImplicitPrototypeChain returns a lazy, dynamic Iterable for the types
// forming the implicit prototype chain of this type.
//
// The chain is iterated bottom to top; from the nearest ancestor to the most
// distant. Iteration stops when the next ancestor would be a null reference.
//
// The created Iterators will not reflect changes to the prototype chain of
// elements it has already iterated past, but will reflect those of upcoming
// elements. Neither the Iterable nor its Iterator support mutation.
//
// Java: ObjectType.getImplicitPrototypeChain()
func (o *ObjectTypeBase) GetImplicitPrototypeChain() iter.Seq[ObjectType] {
	self := o.self

	return func(yield func(ObjectType) bool) {
		next := self // We increment past this type before first access.
		for {
			next = next.GetImplicitPrototype()
			if next == nil || !yield(next) {
				return
			}
		}
	}
}

// DefineDeclaredProperty defines a property whose type is explicitly
// declared by the programmer.
//
// propertyName: the property's name. type: the type. propertyNode: the node
// corresponding to the declaration of property which might later be
// accessed using getPropertyNode.
//
// Java: ObjectType.defineDeclaredProperty(String, JSType, Node)
func (o *ObjectTypeBase) DefineDeclaredProperty(propertyName string, type_ JSType, propertyNode *Node) bool {
	return o.self.DefineDeclaredPropertyByKey(NewProperty_StringKey(propertyName), type_, propertyNode)
}

// Java: ObjectType.defineDeclaredProperty(Property.Key, JSType, Node)
func (o *ObjectTypeBase) DefineDeclaredPropertyByKey(propertyName Property_Key, type_ JSType, propertyNode *Node) bool {
	result := o.self.DefinePropertyByKey(propertyName, type_ /* inferred= */, false, propertyNode)
	if propertyName.Kind() == Property_KeyKind_STRING {
		// All property definitions go through this method or
		// defineInferredProperty. Because the properties defined an an object
		// can affect subtyping, it's slightly more efficient to register this
		// after defining the property.
		o.registry.RegisterPropertyOnType(propertyName.String_(), o.self)
	}
	return result
}

// DefineSynthesizedProperty defines a property whose type is on a
// synthesized object. These objects don't actually exist in the user's
// program. They're just used for bookkeeping in the type system.
//
// Java: ObjectType.defineSynthesizedProperty(String, JSType, Node)
func (o *ObjectTypeBase) DefineSynthesizedProperty(propertyName string, type_ JSType, propertyNode *Node) bool {
	return o.self.DefineProperty(propertyName, type_, false, propertyNode)
}

// DefineInferredPropertyByKey defines a property whose type is inferred.
//
// Java: ObjectType.defineInferredProperty(Property.Key, JSType, Node)
func (o *ObjectTypeBase) DefineInferredPropertyByKey(propertyName Property_Key, type_ JSType, propertyNode *Node) bool {
	// Java: checkArgument(propertyName.string() != null, ...); string() throws
	// UnsupportedOperationException for a symbol key and is never null
	// otherwise.
	name := propertyName.String_()
	jutil.CheckArgumentMsg(
		true, "Symbol-based property keys are not supported for inferred properties.")
	return o.self.DefineInferredProperty(name, type_, propertyNode)
}

// DefineInferredProperty defines a property whose type is inferred.
//
// propertyName: the property's name. type: the type. propertyNode: the node
// corresponding to the inferred definition of property that might later be
// accessed using getPropertyNode.
//
// Java: ObjectType.defineInferredProperty(String, JSType, Node)
func (o *ObjectTypeBase) DefineInferredProperty(propertyName string, type_ JSType, propertyNode *Node) bool {
	if o.self.HasProperty(propertyName) {
		if o.self.IsPropertyTypeDeclared(propertyName) {
			// We never want to hide a declared property with an inferred
			// property.
			return true
		}
		originalType := jutil.CheckNotNull(o.self.GetPropertyType(propertyName))
		type_ = originalType.GetLeastSupertype(type_)
	}
	// TODO(b/140764208): verify that if isResolved() then type.isResolved().
	// Defining unresolved properties on resolved types is dangerous because
	// the property type may never be resolved.

	result := o.self.DefineProperty(propertyName, type_, true, propertyNode)

	// All property definitions go through this method or
	// defineDeclaredProperty. Because the properties defined an an object can
	// affect subtyping, it's slightly more efficient to register this after
	// defining the property.
	o.registry.RegisterPropertyOnType(propertyName, o.self)

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
func (o *ObjectTypeBase) DefineProperty(propertyName string, type_ JSType, inferred bool, propertyNode *Node) bool {
	return o.self.DefinePropertyByKey(NewProperty_StringKey(propertyName), type_, inferred, propertyNode)
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
func (o *ObjectTypeBase) GetPropertyNode(propertyName string) *Node {
	return o.self.GetPropertyNodeByKey(NewProperty_StringKey(propertyName))
}

// Java: ObjectType.getPropertyNode(Property.Key)
func (o *ObjectTypeBase) GetPropertyNodeByKey(propertyName Property_Key) *Node {
	p := o.self.GetSlotByKey(propertyName)
	if p == nil {
		return nil
	}
	return p.GetNode()
}

// Java: ObjectType.getPropertyDefSite(String)
func (o *ObjectTypeBase) GetPropertyDefSite(propertyName string) *Node {
	return o.self.GetPropertyDefSiteByKey(NewProperty_StringKey(propertyName))
}

// Java: ObjectType.getPropertyDefSite(Property.Key)
func (o *ObjectTypeBase) GetPropertyDefSiteByKey(propertyName Property_Key) *Node {
	return o.self.GetPropertyNodeByKey(propertyName)
}

// Java: ObjectType.getPropertyJSDocInfo(String)
func (o *ObjectTypeBase) GetPropertyJSDocInfo(propertyName string) *JSDocInfo {
	return o.self.GetPropertyJSDocInfoByKey(NewProperty_StringKey(propertyName))
}

// Java: ObjectType.getPropertyJSDocInfo(Property.Key)
func (o *ObjectTypeBase) GetPropertyJSDocInfoByKey(propertyName Property_Key) *JSDocInfo {
	p := o.self.GetSlotByKey(propertyName)
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
	return o.self.GetOwnPropertyJSDocInfoByKey(NewProperty_StringKey(propertyName))
}

// Java: ObjectType.getOwnPropertyJSDocInfo(Property.Key)
func (o *ObjectTypeBase) GetOwnPropertyJSDocInfoByKey(propertyName Property_Key) *JSDocInfo {
	p := o.self.GetOwnSlotByKey(propertyName)
	if p == nil {
		return nil
	}
	return p.GetJSDocInfo()
}

// Java: ObjectType.getOwnPropertyDefSite(String)
func (o *ObjectTypeBase) GetOwnPropertyDefSite(propertyName string) *Node {
	return o.self.GetOwnPropertyDefSiteByKey(NewProperty_StringKey(propertyName))
}

// Java: ObjectType.getOwnPropertyDefSite(Property.Key)
func (o *ObjectTypeBase) GetOwnPropertyDefSiteByKey(propertyName Property_Key) *Node {
	p := o.self.GetOwnSlotByKey(propertyName)
	if p == nil {
		return nil
	}
	return p.GetNode()
}

// SetPropertyJSDocInfo sets the docInfo for the specified property from the
// JSDocInfo on its definition. info may be null.
//
// Java: ObjectType.setPropertyJSDocInfo(String, JSDocInfo)
func (o *ObjectTypeBase) SetPropertyJSDocInfo(propertyName string, info *JSDocInfo) {
	o.self.SetPropertyJSDocInfoByKey(NewProperty_StringKey(propertyName), info)
}

// Java: ObjectType.setPropertyJSDocInfo(Property.Key, JSDocInfo)
func (o *ObjectTypeBase) SetPropertyJSDocInfoByKey(propertyName Property_Key, info *JSDocInfo) {
	// by default, do nothing
}

// SetPropertyNode sets the node where the property was defined.
//
// Java: ObjectType.setPropertyNode(String, Node)
func (o *ObjectTypeBase) SetPropertyNode(propertyName string, defSite *Node) {
	o.self.SetPropertyNodeByKey(NewProperty_StringKey(propertyName), defSite)
}

// Java: ObjectType.setPropertyNode(Property.Key, Node)
func (o *ObjectTypeBase) SetPropertyNodeByKey(propertyName Property_Key, defSite *Node) {
	// by default, do nothing
}

// Java: ObjectType.findPropertyTypeWithoutConsideringTemplateTypes(Property.Key)
func (o *ObjectTypeBase) FindPropertyTypeWithoutConsideringTemplateTypes(propertyName Property_Key) JSType {
	if o.self.HasPropertyByKey(propertyName) {
		return o.self.GetPropertyTypeByKey(propertyName)
	}
	return nil
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
	return o.self.GetPropertyTypeByKey(NewProperty_StringKey(propertyName))
}

// Java: ObjectType.getPropertyType(Property.Key)
func (o *ObjectTypeBase) GetPropertyTypeByKey(propertyName Property_Key) JSType {
	slot := o.self.GetSlotByKey(propertyName)
	return o.GetPropertyTypeFromSlot(slot)
}

// Java: ObjectType.GetPropertyTypeFromSlot(StaticTypedSlot)
func (o *ObjectTypeBase) GetPropertyTypeFromSlot(slot *Property) JSType {
	if slot == nil {
		if o.self.IsNoResolvedType() || o.self.IsCheckedUnknownType() {
			return o.GetNativeType(JSTypeNative_CHECKED_UNKNOWN_TYPE)
		} else if o.self.IsEmptyType() {
			return o.GetNativeType(JSTypeNative_NO_TYPE)
		}
		return o.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	}
	return slot.GetType()
}

// Java: ObjectType.getPropertyKind(Property.Key, boolean)
func (o *ObjectTypeBase) GetPropertyKindByKeyWithAutobox(propertyName Property_Key, autobox bool) JSType_HasPropertyKind {
	// Unknown types have all properties.
	return JSType_HasPropertyKind_Of(o.self.IsEmptyType() || o.self.IsUnknownType() || o.self.GetSlotByKey(propertyName) != nil)
}

// GetOwnPropertyKind checks whether the property whose name is given is
// present directly on the object. Returns false even if it is declared on a
// supertype.
//
// Java: ObjectType.getOwnPropertyKind(String)
func (o *ObjectTypeBase) GetOwnPropertyKind(propertyName string) JSType_HasPropertyKind {
	return o.self.GetOwnPropertyKindByKey(NewProperty_StringKey(propertyName))
}

// Java: ObjectType.getOwnPropertyKind(Property.Key)
func (o *ObjectTypeBase) GetOwnPropertyKindByKey(propertyName Property_Key) JSType_HasPropertyKind {
	if o.self.GetOwnSlotByKey(propertyName) != nil {
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
	return o.self.HasOwnPropertyByKey(NewProperty_StringKey(propertyName))
}

// Java: ObjectType.hasOwnProperty(Property.Key)
func (o *ObjectTypeBase) HasOwnPropertyByKey(propertyName Property_Key) bool {
	return o.self.GetOwnPropertyKindByKey(propertyName) != JSType_HasPropertyKind_ABSENT
}

// Java: ObjectType.getOwnPropertyKnownSymbols()
func (o *ObjectTypeBase) GetOwnPropertyKnownSymbols() []*KnownSymbolType {
	return o.self.GetPropertyMap().GetOwnKnownSymbols()
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
	return o.self.GetPropertyMap().GetOwnPropertyNames()
}

// GetOwnPropertyKeys returns the names of all the properties directly on
// this type.
//
// Java: ObjectType.getOwnPropertyKeys()
func (o *ObjectTypeBase) GetOwnPropertyKeys() []Property_Key {
	ownPropertyNames := o.self.GetOwnPropertyNames()
	ownSymbols := o.self.GetOwnPropertyKnownSymbols()
	ownPropertyKeys := jutil.NewObjectLinkedHashSet[Property_Key]()
	for _, propertyName := range ownPropertyNames {
		ownPropertyKeys.Add(NewProperty_StringKey(propertyName))
	}
	for _, symbol := range ownSymbols {
		ownPropertyKeys.Add(NewProperty_SymbolKey(symbol))
	}
	return ownPropertyKeys.ToSlice()
}

// IsPropertyTypeInferred checks whether the property's type is inferred.
//
// Java: ObjectType.isPropertyTypeInferred(String)
func (o *ObjectTypeBase) IsPropertyTypeInferred(propertyName string) bool {
	return o.self.IsPropertyTypeInferredByKey(NewProperty_StringKey(propertyName))
}

// Java: ObjectType.isPropertyTypeInferred(Property.Key)
func (o *ObjectTypeBase) IsPropertyTypeInferredByKey(propertyName Property_Key) bool {
	slot := o.self.GetSlotByKey(propertyName)
	if slot == nil {
		return false
	}
	return slot.IsTypeInferred()
}

// IsPropertyTypeDeclared checks whether the property's type is declared.
//
// Java: ObjectType.isPropertyTypeDeclared(String)
func (o *ObjectTypeBase) IsPropertyTypeDeclared(propertyName string) bool {
	return o.self.IsPropertyTypeDeclaredByKey(NewProperty_StringKey(propertyName))
}

// Java: ObjectType.isPropertyTypeDeclared(Property.Key)
func (o *ObjectTypeBase) IsPropertyTypeDeclaredByKey(propertyName Property_Key) bool {
	slot := o.self.GetSlotByKey(propertyName)
	if slot == nil {
		return false
	}
	return !slot.IsTypeInferred()
}

// Java: ObjectType.isStructuralType()
func (o *ObjectTypeBase) IsStructuralType() bool {
	constructor := o.self.GetConstructor()
	return constructor != nil && constructor.IsStructuralInterface()
}

// HasOwnDeclaredProperty reports whether the given property is declared on
// this object.
//
// Java: ObjectType.hasOwnDeclaredProperty(String)
func (o *ObjectTypeBase) HasOwnDeclaredProperty(name string) bool {
	return o.self.HasOwnProperty(name) && o.self.IsPropertyTypeDeclared(name)
}

// Java: ObjectType.hasOwnDeclaredProperty(Property.Key)
func (o *ObjectTypeBase) HasOwnDeclaredPropertyByKey(name Property_Key) bool {
	return o.self.HasOwnPropertyByKey(name) && o.self.IsPropertyTypeDeclaredByKey(name)
}

// IsPropertyInExterns checks whether the property was defined in the
// externs.
//
// Java: ObjectType.isPropertyInExterns(String)
func (o *ObjectTypeBase) IsPropertyInExterns(propertyName string) bool {
	p := o.self.GetSlot(propertyName)
	if p == nil {
		return false
	}
	return p.IsFromExterns()
}

// GetPropertiesCount gets the number of properties of this object.
//
// Java: ObjectType.getPropertiesCount()
func (o *ObjectTypeBase) GetPropertiesCount() int {
	return o.self.GetPropertyMap().GetPropertiesCount()
}

// GetPropertyNames returns a list of properties defined or inferred on this
// type and any of its supertypes.
//
// Java: ObjectType.getPropertyNames()
func (o *ObjectTypeBase) GetPropertyNames() []string {
	return o.self.GetPropertyMap().GetAllKeys().StringKeys()
}

// GetAllKeys returns a set of all KnownSymbolTypes properties defined on
// this type and any of its supertypes.
//
// Java: ObjectType.getAllKeys()
func (o *ObjectTypeBase) GetAllKeys() *PropertyMap_AllKeys {
	return o.self.GetPropertyMap().GetAllKeys()
}

// Java: ObjectType.visit(Visitor<T>)
func (o *ObjectTypeBase) Visit(visitor Visitor) any {
	return visitor.CaseObjectType(o.self)
}

// Java: ObjectType.visit(RelationshipVisitor<T>, JSType)
func (o *ObjectTypeBase) VisitRelationship(visitor RelationshipVisitor, that JSType) any {
	return visitor.CaseObjectType(o.self, that)
}

// IsImplicitPrototypeOf returns whether this is on the implicit prototype
// chain of other.
//
// this need not be the immediate prototype of other.
//
// Java: ObjectType.isImplicitPrototypeOf(ObjectType)
func (o *ObjectTypeBase) IsImplicitPrototypeOf(other ObjectType) bool {
	unwrappedThis := ObjectType_DeeplyUnwrap(o.self)

	for other = ObjectType_DeeplyUnwrap(other); other != nil; other = ObjectType_DeeplyUnwrap(other.GetImplicitPrototype()) {
		// The prototype should match exactly. NOTE: the use of "==" here rather
		// than equals is deliberate. This method is very hot in the type
		// checker and relying on identity improves performance of both type
		// checking/type inferrence and property disambiguation.
		if unwrappedThis == other {
			return true
		}
	}

	return false
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
		proxy, ok := current.(ProxyObjectType)
		if !ok {
			break
		}
		if current.IsTemplatizedType() {
			current = current.ToMaybeTemplatizedType().GetReferencedType()
		} else if current.IsNamedType() {
			if !current.IsSuccessfullyResolved() {
				break
			}
			current = current.ToMaybeNamedType().GetReferencedObjTypeInternal()
		} else {
			// TODO(lharker): remove this case and instead fail. Only the Rhino unit
			// tests are triggering this by creating new ProxyObjectTypes.
			current = proxy.GetReferencedObjTypeInternal()
		}
	}
	return current
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
		implicitProto := o.self.GetImplicitPrototype()
		if implicitProto == nil || implicitProto.IsNativeObjectType() {
			o.unknown = false
			for _, interfaceType := range o.self.GetCtorExtendedInterfaces() {
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
	return o.self.GetOwnerFunction() != nil
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

// GetPropertyTypeMap gets the map of properties to types covered in an
// object type. Returns a Map that maps the property's name to the
// property's type
//
// Java: ObjectType.getPropertyTypeMap()
func (o *ObjectTypeBase) GetPropertyTypeMap() *jutil.LinkedHashMap[string, JSType] {
	propTypeMap := jutil.NewLinkedHashMap[string, JSType]()
	for _, name := range o.self.GetPropertyNames() {
		if propTypeMap.ContainsKey(name) {
			// ImmutableMap.Builder.buildOrThrow() rejects duplicate keys.
			panic(jutil.NewIllegalArgumentExceptionMsg("Multiple entries with same key: " + name))
		}
		propTypeMap.Put(name, o.self.GetPropertyType(name))
	}
	return propTypeMap
}

// Java: ObjectType.getEnumeratedTypeOfEnumObject()
func (o *ObjectTypeBase) GetEnumeratedTypeOfEnumObject() JSType {
	return nil
}
