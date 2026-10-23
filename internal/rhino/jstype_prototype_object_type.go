// Port of com/google/javascript/rhino/jstype/PrototypeObjectType.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// PrototypeObjectType is the object type that represents instances of
// JavaScript objects such as Object, Date, Function.
//
// Objects in JavaScript are unordered collections of properties. Each
// property consists of a name, a value and a set of attributes.
//
// Each instance has an implicit prototype property ([[Prototype]]) pointing
// to an object instance, which itself has an implicit property, thus forming
// a chain.
//
// A class begins life with no name. Later, a name may be provided once it
// can be inferred. Note that the name in this case is strictly for debugging
// purposes. Looking up type name references goes through the
// JSTypeRegistry.
//
// Java: class PrototypeObjectType extends ObjectType (subclassed by
// EnumType, FunctionType, InstanceObjectType and RecordType).
type PrototypeObjectType interface {
	ObjectType

	SetPrettyPrint(prettyPrint bool)
	IsPrettyPrint() bool
	SetImplicitPrototype(implicitPrototype ObjectType)
	IsAnonymous() bool
	MatchRecordTypeConstraint(constraintObj ObjectType)

	// GetPrototypeObjectTypeBase gives access to the base through the
	// interface.
	//
	// Go-only: language-adapter — field access through the interface
	// (PORTING.md §4).
	GetPrototypeObjectTypeBase() *PrototypeObjectTypeBase
}

// PrototypeObjectTypeBase holds PrototypeObjectType's fields and methods;
// plain prototype object types are *PrototypeObjectTypeBase values.
type PrototypeObjectTypeBase struct {
	ObjectTypeBase

	className *string

	// NOTE(nicksantos): The implicit prototype can change over time.
	// Modeling this is a bear. Always call getImplicitPrototype(), because
	// some subclasses override this to do special resolution handling.
	implicitPrototypeFallback ObjectType

	// If this is a function prototype, then this is the owner. A
	// PrototypeObjectType can only be the prototype of one function. If we
	// try to do this for multiple functions, then we'll have to create a new
	// one.
	ownerFunction *FunctionTypeBase

	properties PropertyMap

	templateParamCount int32
	nativeType         bool
	anonymousType      bool

	// Whether the toString representation of this should be pretty-printed,
	// by printing all properties.
	prettyPrint bool
}

// Java: PrototypeObjectType.TYPE_CLASS
const PrototypeObjectType_TYPE_CLASS = JSTypeClass_PROTOTYPE_OBJECT

// Java: PrototypeObjectType.MAX_PRETTY_PRINTED_PROPERTIES
const PrototypeObjectType_MAX_PRETTY_PRINTED_PROPERTIES = 10

// NewPrototypeObjectType creates an object type, allowing specification of
// the implicit prototype, whether the object is native, and any templatized
// types.
//
// Java: PrototypeObjectType(Builder<?>)
func NewPrototypeObjectType(builder *PrototypeObjectType_BuilderFields) *PrototypeObjectTypeBase {
	t := &PrototypeObjectTypeBase{}
	t.objFlags = objKindPrototype
	t.InitPrototypeObjectTypeBase(t, builder)
	return t
}

// InitPrototypeObjectTypeBase is the Go form of the constructor for
// subclasses.
//
// Java: PrototypeObjectType(Builder<?>)
func (p *PrototypeObjectTypeBase) InitPrototypeObjectTypeBase(self PrototypeObjectType, builder *PrototypeObjectType_BuilderFields) {
	p.InitObjectTypeBaseWithTemplateTypeMap(self, builder.registry, builder.templateTypeMap)

	p.className = builder.className
	p.templateParamCount = int32(builder.templateParamCount)
	p.nativeType = builder.nativeType
	p.anonymousType = builder.anonymousType

	p.properties.sorted = true
	p.properties.embeddedInProto = true
	p.properties.cachedKeySetCounter = 1

	if p.nativeType || builder.implicitPrototype != nil {
		p.SetImplicitPrototype(builder.implicitPrototype)
	} else {
		p.SetImplicitPrototype(p.registry.GetNativeObjectType(JSTypeNative_OBJECT_TYPE))
	}

	if p.anonymousType {
		jutil.CheckState(p.className == nil)
	}
	jutil.CheckNotNull(p.templateTypeMap)
	// Also guarantees `templateParamCount >= 0`.
	jutil.CheckState(p.templateTypeMap.Size() >= int(p.templateParamCount))

	p.registry.GetResolver().ResolveIfClosed(self, PrototypeObjectType_TYPE_CLASS)
}

// GetPrototypeObjectTypeBase gives access to the base through the interface.
//
// Go-only: language-adapter — field access through the interface (PORTING.md §4).
func (p *PrototypeObjectTypeBase) GetPrototypeObjectTypeBase() *PrototypeObjectTypeBase { return p }

// Java: PrototypeObjectType.getTypeClass()
func (p *PrototypeObjectTypeBase) GetTypeClass() JSTypeClass {
	return PrototypeObjectType_TYPE_CLASS
}

// PrototypeObjectType_BuilderFields holds the fields of
// PrototypeObjectType.Builder.
//
// Go-only: language-adapter — Builder<T extends Builder<T>> is generic in the
// type its setters return; the fields live in this non-generic struct so the
// PrototypeObjectType constructor can read any builder.
type PrototypeObjectType_BuilderFields struct {
	registry *JSTypeRegistry

	className         *string
	implicitPrototype ObjectType

	nativeType    bool
	anonymousType bool

	templateTypeMap    *TemplateTypeMap
	templateParamCount int
}

// PrototypeObjectType_Builder is PrototypeObjectType.Builder<T>: its
// setters return the concrete builder T (Java's castThis()).
type PrototypeObjectType_Builder[T any] struct {
	PrototypeObjectType_BuilderFields
	self T
}

// InitPrototypeObjectType_Builder is the Go form of the Builder
// constructor; self is the concrete builder.
//
// Java: PrototypeObjectType.Builder(JSTypeRegistry)
func (b *PrototypeObjectType_Builder[T]) InitPrototypeObjectType_Builder(self T, registry *JSTypeRegistry) {
	b.self = self
	b.registry = registry

	b.templateTypeMap = registry.GetEmptyTemplateTypeMap()
}

// GetBuilderFields gives access to the builder's fields.
//
// Go-only: language-adapter — see PrototypeObjectType_BuilderFields.
func (b *PrototypeObjectType_Builder[T]) GetBuilderFields() *PrototypeObjectType_BuilderFields {
	return &b.PrototypeObjectType_BuilderFields
}

// Java: PrototypeObjectType.Builder.setName(String)
func (b *PrototypeObjectType_Builder[T]) SetName(x *string) T {
	b.className = x
	return b.CastThis()
}

// Java: PrototypeObjectType.Builder.setImplicitPrototype(ObjectType)
func (b *PrototypeObjectType_Builder[T]) SetImplicitPrototype(x ObjectType) T {
	b.implicitPrototype = x
	return b.CastThis()
}

// Java: PrototypeObjectType.Builder.setNative(boolean)
func (b *PrototypeObjectType_Builder[T]) SetNative(x bool) T {
	b.nativeType = x
	return b.CastThis()
}

// Java: PrototypeObjectType.Builder.setAnonymous(boolean)
func (b *PrototypeObjectType_Builder[T]) SetAnonymous(x bool) T {
	b.anonymousType = x
	return b.CastThis()
}

// Java: PrototypeObjectType.Builder.setTemplateTypeMap(TemplateTypeMap)
func (b *PrototypeObjectType_Builder[T]) SetTemplateTypeMap(x *TemplateTypeMap) T {
	b.templateTypeMap = x
	return b.CastThis()
}

// Java: PrototypeObjectType.Builder.setTemplateParamCount(int)
func (b *PrototypeObjectType_Builder[T]) SetTemplateParamCount(x int) T {
	b.templateParamCount = x
	return b.CastThis()
}

// Java: PrototypeObjectType.Builder.castThis()
func (b *PrototypeObjectType_Builder[T]) CastThis() T {
	return b.self
}

// PrototypeObjectType_PlainBuilder is the Builder<?> returned by
// PrototypeObjectType.builder(registry), which builds plain
// PrototypeObjectTypes.
//
// Go-only: language-adapter — names the instantiation of the generic
// Builder<T> whose T is itself.
type PrototypeObjectType_PlainBuilder struct {
	PrototypeObjectType_Builder[*PrototypeObjectType_PlainBuilder]
}

// Java: PrototypeObjectType.Builder(JSTypeRegistry)
func NewPrototypeObjectType_PlainBuilder(registry *JSTypeRegistry) *PrototypeObjectType_PlainBuilder {
	b := &PrototypeObjectType_PlainBuilder{}
	b.InitPrototypeObjectType_Builder(b, registry)
	return b
}

// Java: PrototypeObjectType.Builder.build()
func (b *PrototypeObjectType_PlainBuilder) Build() PrototypeObjectType {
	return NewPrototypeObjectType(b.GetBuilderFields())
}

// Java: PrototypeObjectType.builder(JSTypeRegistry)
func PrototypeObjectType_Builder_(registry *JSTypeRegistry) *PrototypeObjectType_PlainBuilder {
	return NewPrototypeObjectType_PlainBuilder(registry)
}

// Java: PrototypeObjectType.getPropertyMap()
func (p *PrototypeObjectTypeBase) GetPropertyMap() *PropertyMap {
	return &p.properties
}

// Java: PrototypeObjectType.defineProperty(Property.Key, JSType, boolean, Node)
func (p *PrototypeObjectTypeBase) DefinePropertyByKey(name Property_Key, type_ JSType, inferred bool, propertyNode *Node) bool {
	if p.selfObj().HasOwnDeclaredPropertyByKey(name) {
		return false
	}
	newProp := NewPropertyByKey(name, type_, inferred, propertyNode)
	p.properties.PutPropertyByKey(name, newProp)
	return true
}

// Java: PrototypeObjectType.setPropertyJSDocInfo(Property.Key, JSDocInfo)
func (p *PrototypeObjectTypeBase) SetPropertyJSDocInfoByKey(propertyName Property_Key, info *JSDocInfo) {
	if info != nil {
		property := p.properties.GetOwnPropertyByKey(propertyName)
		if property == nil {
			// If docInfo was attached, but the type of the property was not
			// defined anywhere, then we consider this an explicit declaration of
			// the property.
			self := p.selfObj()
			self.DefineInferredPropertyByKey(propertyName, self.GetPropertyTypeByKey(propertyName), nil)
			property = p.properties.GetOwnPropertyByKey(propertyName)
		}

		// The prototype property is not represented as a normal Property. We
		// probably don't want to attach any JSDoc to it anyway.
		if property != nil {
			property.SetJSDocInfo(info)
		}
	}
}

// Java: PrototypeObjectType.setPropertyNode(Property.Key, Node)
func (p *PrototypeObjectTypeBase) SetPropertyNodeByKey(propertyName Property_Key, defSite *Node) {
	property := p.properties.GetOwnPropertyByKey(propertyName)
	if property != nil {
		property.SetNode(defSite)
	}
}

// Java: PrototypeObjectType.matchesNumberContext()
func (p *PrototypeObjectTypeBase) MatchesNumberContext() bool {
	// BigInt is intentionally left out here. It cannot be coerced to a Number.
	return p.self.IsNumberObjectType() ||
		p.self.IsDateType() ||
		p.self.IsBooleanObjectType() ||
		p.self.IsStringObjectType() ||
		p.HasOverriddenNativeProperty("valueOf")
}

// Java: PrototypeObjectType.matchesStringContext()
func (p *PrototypeObjectTypeBase) MatchesStringContext() bool {
	return p.self.IsTheObjectType() ||
		p.self.IsStringObjectType() ||
		p.self.IsDateType() ||
		p.self.IsRegexpType() ||
		p.self.IsArrayType() ||
		p.self.IsNumberObjectType() ||
		p.self.IsBigIntObjectType() ||
		p.self.IsBooleanObjectType() ||
		p.HasOverriddenNativeProperty("toString")
}

// Java: PrototypeObjectType.matchesSymbolContext()
func (p *PrototypeObjectTypeBase) MatchesSymbolContext() bool {
	return p.self.IsSymbolObjectType()
}

// hasOverriddenNativeProperty: given the name of a native object property,
// checks whether the property is present on the object and different from
// the native one.
//
// Java: PrototypeObjectType.HasOverriddenNativeProperty(String)
func (p *PrototypeObjectTypeBase) HasOverriddenNativeProperty(propertyName string) bool {
	if p.self.IsNativeObjectType() {
		return false
	}

	propertyType := p.selfObj().GetPropertyType(propertyName)
	var nativeType ObjectType
	if p.self.IsFunctionType() {
		nativeType = p.registry.GetNativeObjectType(JSTypeNative_FUNCTION_PROTOTYPE)
	} else {
		nativeType = p.registry.GetNativeObjectType(JSTypeNative_OBJECT_PROTOTYPE)
	}
	nativePropertyType := nativeType.GetPropertyType(propertyName)
	return propertyType != nativePropertyType
}

// Java: PrototypeObjectType.matchesObjectContext()
func (p *PrototypeObjectTypeBase) MatchesObjectContext() bool {
	return true
}

// PrototypeObjectType_PropertyForPrettyPrinting is the private record
// PropertyForPrettyPrinting(Property.Key key, JSType type).
type PrototypeObjectType_PropertyForPrettyPrinting struct {
	key   Property_Key
	type_ JSType
}

// Java: PrototypeObjectType.PropertyForPrettyPrinting.appendTo(TypeStringBuilder)
func (r PrototypeObjectType_PropertyForPrettyPrinting) AppendTo(sb *TypeStringBuilder) {
	switch r.key.Kind() {
	case Property_KeyKind_STRING:
		sb.Append(r.key.HumanReadableName())
	case Property_KeyKind_SYMBOL:
		sb.Append("[")
		sb.Append(r.key.HumanReadableName())
		sb.Append("]")
	}
	sb.Append(": ").AppendNonNull(r.type_)
}

// Java: PrototypeObjectType.PROP_COMPARATOR
func PrototypeObjectType_PROP_COMPARATOR(thisProp, otherProp PrototypeObjectType_PropertyForPrettyPrinting) int {
	// Handle the easy case where the property names are distinct.
	if thisProp.key.HumanReadableName() != otherProp.key.HumanReadableName() {
		return jutil.StringCompareTo(thisProp.key.HumanReadableName(), otherProp.key.HumanReadableName())
	}
	switch thisProp.key.Kind() {
	case Property_KeyKind_STRING:
		switch otherProp.key.Kind() {
		case Property_KeyKind_STRING:
			return jutil.StringCompareTo(thisProp.key.String_(), otherProp.key.String_())
		}
		return -1 // order string keys before symbol keys
	}
	switch otherProp.key.Kind() {
	case Property_KeyKind_STRING:
		return 1 // order string keys before symbol keys.
	}
	// Given two symbol keys with the same human-readable name, compare them by
	// their type string. This should be rare, so it's not a big performance
	// hit to do the extra toString()s.
	return jutil.StringCompareTo(thisProp.type_.String(), otherProp.type_.String())
}

// Java: PrototypeObjectType.appendTo(TypeStringBuilder)
func (p *PrototypeObjectTypeBase) AppendTo(sb *TypeStringBuilder) {
	self := p.selfObj()
	if self.HasReferenceName() {
		if sb.IsForAnnotations() {
			sb.Append(*self.GetNormalizedReferenceName())
		} else {
			sb.Append(*self.GetReferenceName())
		}
		return
	}

	if !p.prettyPrint {
		if sb.IsForAnnotations() {
			sb.Append("?")
		} else {
			sb.Append("{...}")
		}
		return
	}

	// Don't pretty print recursively. It would cause infinite recursion.
	p.prettyPrint = false

	// Use a tree set so that the properties are sorted.
	properties := jutil.NewTreeSet[PrototypeObjectType_PropertyForPrettyPrinting](PrototypeObjectType_PROP_COMPARATOR)
	for current := self; current != nil; current = current.GetImplicitPrototype() {
		if current.IsNativeObjectType() || properties.Size() > PrototypeObjectType_MAX_PRETTY_PRINTED_PROPERTIES {
			break
		}

		for _, stringKey := range current.GetOwnPropertyNames() {
			key := NewProperty_StringKey(stringKey)
			properties.Add(PrototypeObjectType_PropertyForPrettyPrinting{key, current.GetPropertyTypeByKey(key)})
		}
		for _, symbolKey := range current.GetOwnPropertyKnownSymbols() {
			key := NewProperty_SymbolKey(symbolKey)
			properties.Add(PrototypeObjectType_PropertyForPrettyPrinting{key, current.GetPropertyTypeByKey(key)})
		}
	}

	multiline := !sb.IsForAnnotations() && properties.Size() > 1
	sb.Append("{").
		Indent(func() {
			if multiline {
				sb.BreakLineAndIndent()
			}

			i := 0
			for property := range properties.All() {
				i++

				if !sb.IsForAnnotations() && i > PrototypeObjectType_MAX_PRETTY_PRINTED_PROPERTIES {
					sb.Append("...")
					break
				}

				property.AppendTo(sb)
				if i < properties.Size() {
					sb.Append(",")
					if multiline {
						sb.BreakLineAndIndent()
					} else {
						sb.Append(" ")
					}
				}
			}
		})
	if multiline {
		sb.BreakLineAndIndent()
	}
	sb.Append("}")

	p.prettyPrint = true
}

// Java: PrototypeObjectType.setPrettyPrint(boolean)
func (p *PrototypeObjectTypeBase) SetPrettyPrint(prettyPrint bool) {
	p.prettyPrint = prettyPrint
}

// Java: PrototypeObjectType.isPrettyPrint()
func (p *PrototypeObjectTypeBase) IsPrettyPrint() bool {
	return p.prettyPrint
}

// Java: PrototypeObjectType.getConstructor()
func (p *PrototypeObjectTypeBase) GetConstructor() FunctionType {
	return nil
}

// Java: PrototypeObjectType.getImplicitPrototype()
func (p *PrototypeObjectTypeBase) GetImplicitPrototype() ObjectType {
	return p.implicitPrototypeFallback
}

// SetImplicitPrototype: this should only be reset on the
// FunctionPrototypeType, only to fix an incorrectly established prototype
// chain due to the user having a mismatch in super class declaration, and
// only before properties on that type are processed.
//
// Java: PrototypeObjectType.setImplicitPrototype(ObjectType)
func (p *PrototypeObjectTypeBase) SetImplicitPrototype(implicitPrototype ObjectType) {
	jutil.CheckState(!p.selfObj().HasCachedValues())
	p.implicitPrototypeFallback = implicitPrototype
	if p.registry != nil {
		p.registry.leastSupertypeGen++
	}
	if implicitPrototype != nil {
		p.self.MaybeLoosenTypecheckingDueToForwardReferencedSupertype(implicitPrototype)
	}
}

// Java: PrototypeObjectType.getTemplateParamCount()
func (p *PrototypeObjectTypeBase) GetTemplateParamCount() int {
	return int(p.templateParamCount)
}

// Java: ObjectType.hasReferenceName()
func (p *PrototypeObjectTypeBase) HasReferenceName() bool {
	return p.className != nil || p.ownerFunction != nil
}

// Java: PrototypeObjectType.getReferenceName()
func (p *PrototypeObjectTypeBase) GetReferenceName() *string {
	if p.className != nil {
		return p.className
	} else if p.ownerFunction != nil {
		name := jutil.ToJavaString(p.ownerFunction.GetReferenceName()) + ".prototype"
		return &name
	} else {
		return nil
	}
}

// Java: PrototypeObjectType.isAnonymous()
func (p *PrototypeObjectTypeBase) IsAnonymous() bool {
	return p.anonymousType
}

// IsNativeObjectType reports whether this is a built-in object.
//
// Java: PrototypeObjectType.isNativeObjectType()
func (p *PrototypeObjectTypeBase) IsNativeObjectType() bool {
	return p.nativeType
}

// Java: PrototypeObjectType.setOwnerFunction(FunctionType)
func (p *PrototypeObjectTypeBase) SetOwnerFunction(type_ FunctionType) {
	jutil.CheckState(p.ownerFunction == nil || type_ == nil)
	if type_ == nil {
		p.ownerFunction = nil
	} else {
		p.ownerFunction = type_.GetFunctionTypeBase()
	}
}

// Java: PrototypeObjectType.getOwnerFunction()
func (p *PrototypeObjectTypeBase) GetOwnerFunction() FunctionType {
	if p.ownerFunction == nil {
		return nil
	}
	return p.ownerFunction.selfFunc()
}

// Java: PrototypeObjectType.getCtorImplementedInterfaces()
func (p *PrototypeObjectTypeBase) GetCtorImplementedInterfaces() []ObjectType {
	if p.self.IsFunctionPrototypeType() {
		return p.selfObj().GetOwnerFunction().GetImplementedInterfaces()
	}
	return []ObjectType{}
}

// Java: PrototypeObjectType.getCtorExtendedInterfaces()
func (p *PrototypeObjectTypeBase) GetCtorExtendedInterfaces() []ObjectType {
	if p.self.IsFunctionPrototypeType() {
		return p.selfObj().GetOwnerFunction().GetExtendedInterfaces()
	}
	return []ObjectType{}
}

// Java: PrototypeObjectType.resolveInternal(ErrorReporter)
func (p *PrototypeObjectTypeBase) ResolveInternal(reporter ErrorReporter) JSType {
	implicitPrototype := p.selfObj().GetImplicitPrototype()
	if implicitPrototype != nil {
		p.implicitPrototypeFallback = implicitPrototype.Resolve(reporter).(ObjectType)
	}
	p.properties.ensureSorted()
	for i := 0; i < len(p.properties.entries); i++ {
		prop := p.properties.entries[i]
		prop.SetType(JSType_SafeResolve(prop.GetType(), reporter))
	}
	return p.self
}

// Java: PrototypeObjectType.matchConstraint(JSType)
func (p *PrototypeObjectTypeBase) MatchConstraint(constraint JSType) {
	// We only want to match constraints on anonymous types.
	if p.selfObj().HasReferenceName() {
		return
	}

	// Handle the case where the constraint object is a record type.
	//
	// param constraint {{prop: (number|undefined)}}
	// function f(constraint) {}
	// f({});
	//
	// We want to modify the object literal to match the constraint, by taking
	// any each property on the record and trying to match properties on this
	// object.
	if constraint.IsRecordType() {
		p.MatchRecordTypeConstraint(constraint.ToObjectType())
	} else if constraint.IsUnionType() {
		for _, alt := range constraint.ToMaybeUnionType().GetAlternates() {
			if alt.IsRecordType() {
				p.MatchRecordTypeConstraint(alt.ToObjectType())
			}
		}
	}
}

// Java: PrototypeObjectType.matchRecordTypeConstraint(ObjectType)
func (p *PrototypeObjectTypeBase) MatchRecordTypeConstraint(constraintObj ObjectType) {
	self := p.selfObj()
	for _, prop := range constraintObj.GetOwnPropertyNames() {
		propType := constraintObj.GetPropertyType(prop)
		if !self.IsPropertyTypeDeclared(prop) {
			typeToInfer := propType
			if !p.self.HasProperty(prop) {
				typeToInfer = p.GetNativeType(JSTypeNative_VOID_TYPE).GetLeastSupertype(propType)
			}
			self.DefineInferredProperty(prop, typeToInfer, nil)
		}
	}
}

// Java: PrototypeObjectType.recursionUnsafeHashCode()
func (p *PrototypeObjectTypeBase) RecursionUnsafeHashCode() int32 {
	if p.self.IsStructuralType() {
		var classHash int32
		if p.className != nil {
			classHash = jutil.StringHashCode(*p.className)
		}
		return 31*(31+classHash) + p.properties.HashCode()
	} else {
		return jutil.IdentityHashCode(p.self)
	}
}

// jsTypeNullableString converts a nullable Java String to the value Java
// would pass as an Object: nil for null, the string otherwise.
//
// Go-only: language-adapter — *string stands for a nullable String.
func jsTypeNullableString(s *string) any {
	if s == nil {
		return nil
	}
	return *s
}
