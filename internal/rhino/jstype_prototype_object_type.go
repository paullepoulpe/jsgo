// Port of com/google/javascript/rhino/jstype/PrototypeObjectType.java @ closure-compiler v20261021

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

// PrototypeObjectTypeBase holds PrototypeObjectType's fields and methods;
// plain prototype object types are *PrototypeObjectTypeBase values.
//
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
type PrototypeObjectTypeBase struct {
	ObjectTypeBase

	// NOTE(nicksantos): The implicit prototype can change over time.
	// Modeling this is a bear. Always call getImplicitPrototype(), because
	// some subclasses override this to do special resolution handling.
	//
	// Go-only: performance — the type's base (8 bytes) rather than an
	// ObjectType interface (16).
	implicitPrototypeFallback *ObjectTypeBase

	// If this is a function prototype, then this is the owner. A
	// PrototypeObjectType can only be the prototype of one function. If we
	// try to do this for multiple functions, then we'll have to create a new
	// one.
	ownerFunction *FunctionTypeBase

	properties PropertyMap

	templateParamCount int32
	classNameID        rhinoStringID
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
func NewPrototypeObjectType(builder *PrototypeObjectType_Builder) *PrototypeObjectTypeBase {
	t := &PrototypeObjectTypeBase{}
	t.objFlags = objKindPrototype
	t.InitPrototypeObjectTypeBase(t, builder)
	return t
}

// InitPrototypeObjectTypeBase is the Go form of the constructor for
// subclasses.
//
// Java: PrototypeObjectType(Builder<?>)
func (p *PrototypeObjectTypeBase) InitPrototypeObjectTypeBase(self ObjectType, builder *PrototypeObjectType_Builder) {
	p.InitObjectTypeBaseWithTemplateTypeMap(self, builder.registry, builder.templateTypeMap)

	p.classNameID = builder.classNameID
	p.templateParamCount = int32(builder.templateParamCount)
	if builder.nativeType {
		p.protoFlags |= protoFlagNative
	}
	if builder.anonymousType {
		p.protoFlags |= protoFlagAnonymous
	}

	p.properties.sorted = true
	p.properties.parentSource = &p.ObjectTypeBase
	p.properties.cachedKeySetCounter = 1

	if builder.nativeType || builder.implicitPrototype != nil {
		p.SetImplicitPrototype(builder.implicitPrototype)
	} else {
		p.SetImplicitPrototype(p.GetRegistry().GetNativeObjectType(JSTypeNative_OBJECT_TYPE))
	}

	if builder.anonymousType {
		jutil.CheckState(p.classNameID == 0)
	}
	jutil.CheckNotNull(p.templateTypeMap)
	// Also guarantees `templateParamCount >= 0`.
	jutil.CheckState(p.templateTypeMap.Size() >= int(p.templateParamCount))

	p.GetRegistry().GetResolver().ResolveIfClosed(self, PrototypeObjectType_TYPE_CLASS)
}

// Java: PrototypeObjectType.getTypeClass()
func (p *PrototypeObjectTypeBase) GetTypeClass() JSTypeClass {
	return PrototypeObjectType_TYPE_CLASS
}

// PrototypeObjectType_Builder holds the fields and fluent setters of
// PrototypeObjectType.Builder.
type PrototypeObjectType_Builder struct {
	registry *JSTypeRegistry

	implicitPrototype  ObjectType
	templateTypeMap    *TemplateTypeMap
	templateParamCount int
	classNameID        rhinoStringID
	nativeType         bool
	anonymousType      bool
}

// InitPrototypeObjectType_Builder initializes the embedded base builder.
//
// Java: PrototypeObjectType.Builder(JSTypeRegistry)
func (b *PrototypeObjectType_Builder) InitPrototypeObjectType_Builder(registry *JSTypeRegistry) {
	b.registry = registry
	b.templateTypeMap = registry.GetEmptyTemplateTypeMap()
}

// Java: PrototypeObjectType.Builder(JSTypeRegistry)
func NewPrototypeObjectType_Builder(registry *JSTypeRegistry) *PrototypeObjectType_Builder {
	b := &PrototypeObjectType_Builder{}
	b.InitPrototypeObjectType_Builder(registry)
	return b
}

// Java: PrototypeObjectType.Builder.setName(String)
func (b *PrototypeObjectType_Builder) SetName(x *string) *PrototypeObjectType_Builder {
	b.classNameID = rhinoStringIDFromPtr(x)
	return b
}

// SetNameString sets the non-null class name without allocating a *string.
func (b *PrototypeObjectType_Builder) SetNameString(x string) *PrototypeObjectType_Builder {
	b.classNameID = rhinoStringPoolIntern(x)
	return b
}

// Java: PrototypeObjectType.Builder.setImplicitPrototype(ObjectType)
func (b *PrototypeObjectType_Builder) SetImplicitPrototype(x ObjectType) *PrototypeObjectType_Builder {
	b.implicitPrototype = x
	return b
}

// Java: PrototypeObjectType.Builder.setNative(boolean)
func (b *PrototypeObjectType_Builder) SetNative(x bool) *PrototypeObjectType_Builder {
	b.nativeType = x
	return b
}

// Java: PrototypeObjectType.Builder.setAnonymous(boolean)
func (b *PrototypeObjectType_Builder) SetAnonymous(x bool) *PrototypeObjectType_Builder {
	b.anonymousType = x
	return b
}

// Java: PrototypeObjectType.Builder.setTemplateTypeMap(TemplateTypeMap)
func (b *PrototypeObjectType_Builder) SetTemplateTypeMap(x *TemplateTypeMap) *PrototypeObjectType_Builder {
	b.templateTypeMap = x
	return b
}

// Java: PrototypeObjectType.Builder.setTemplateParamCount(int)
func (b *PrototypeObjectType_Builder) SetTemplateParamCount(x int) *PrototypeObjectType_Builder {
	b.templateParamCount = x
	return b
}

// Java: PrototypeObjectType.Builder.build()
func (b *PrototypeObjectType_Builder) Build() *PrototypeObjectTypeBase {
	return NewPrototypeObjectType(b)
}

func prototypeObjectTypeBaseOf(t JSType) *PrototypeObjectTypeBase {
	switch x := t.(type) {
	case *PrototypeObjectTypeBase:
		return x
	case *InstanceObjectType:
		return &x.PrototypeObjectTypeBase
	case *FunctionTypeBase:
		return &x.PrototypeObjectTypeBase
	case *RecordType:
		return &x.PrototypeObjectTypeBase
	case *EnumType:
		return &x.PrototypeObjectTypeBase
	default:
		return nil
	}
}

// Java: PrototypeObjectType.builder(JSTypeRegistry)
func PrototypeObjectType_Builder_(registry *JSTypeRegistry) *PrototypeObjectType_Builder {
	return NewPrototypeObjectType_Builder(registry)
}

// Java: PrototypeObjectType.getPropertyMap()
func (p *PrototypeObjectTypeBase) GetPropertyMap() *PropertyMap {
	return &p.properties
}

// Java: PrototypeObjectType.defineProperty(Property.Key, JSType, boolean, Node)
func (p *PrototypeObjectTypeBase) DefinePropertyByKey(name Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
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
			self.DefineInferredPropertyByKey(propertyName, self.GetPropertyTypeByKey(propertyName), NilNode)
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
func (p *PrototypeObjectTypeBase) SetPropertyNodeByKey(propertyName Property_Key, defSite Node) {
	property := p.properties.GetOwnPropertyByKey(propertyName)
	if property != nil {
		property.SetNode(defSite)
	}
}

// Java: PrototypeObjectType.matchesNumberContext()
func (p *PrototypeObjectTypeBase) MatchesNumberContext() bool {
	// BigInt is intentionally left out here. It cannot be coerced to a Number.
	return p.self().IsNumberObjectType() ||
		p.self().IsDateType() ||
		p.self().IsBooleanObjectType() ||
		p.self().IsStringObjectType() ||
		p.HasOverriddenNativeProperty("valueOf")
}

// Java: PrototypeObjectType.matchesStringContext()
func (p *PrototypeObjectTypeBase) MatchesStringContext() bool {
	return p.self().IsTheObjectType() ||
		p.self().IsStringObjectType() ||
		p.self().IsDateType() ||
		p.self().IsRegexpType() ||
		p.self().IsArrayType() ||
		p.self().IsNumberObjectType() ||
		p.self().IsBigIntObjectType() ||
		p.self().IsBooleanObjectType() ||
		p.HasOverriddenNativeProperty("toString")
}

// Java: PrototypeObjectType.matchesSymbolContext()
func (p *PrototypeObjectTypeBase) MatchesSymbolContext() bool {
	return p.self().IsSymbolObjectType()
}

// hasOverriddenNativeProperty: given the name of a native object property,
// checks whether the property is present on the object and different from
// the native one.
//
// Java: PrototypeObjectType.HasOverriddenNativeProperty(String)
func (p *PrototypeObjectTypeBase) HasOverriddenNativeProperty(propertyName string) bool {
	if p.self().IsNativeObjectType() {
		return false
	}

	propertyType := p.selfObj().GetPropertyType(propertyName)
	var nativeType ObjectType
	if p.self().IsFunctionType() {
		nativeType = p.GetRegistry().GetNativeObjectType(JSTypeNative_FUNCTION_PROTOTYPE)
	} else {
		nativeType = p.GetRegistry().GetNativeObjectType(JSTypeNative_OBJECT_PROTOTYPE)
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
	if r.key.symbol == nil {
		sb.Append(r.key.HumanReadableName())
	} else {
		sb.Append("[")
		sb.Append(r.key.HumanReadableName())
		sb.Append("]")
	}
	sb.Append(": ")
	sb.AppendNonNull(r.type_)
}

// Java: PrototypeObjectType.PROP_COMPARATOR
func PrototypeObjectType_PROP_COMPARATOR(thisProp, otherProp PrototypeObjectType_PropertyForPrettyPrinting) int {
	// Handle the easy case where the property names are distinct.
	thisName := thisProp.key.HumanReadableName()
	otherName := otherProp.key.HumanReadableName()
	if thisName != otherName {
		return jutil.StringCompareTo(thisName, otherName)
	}
	if thisProp.key.symbol == nil {
		if otherProp.key.symbol == nil {
			return jutil.StringCompareTo(thisName, otherName)
		}
		return -1 // order string keys before symbol keys
	}
	if otherProp.key.symbol == nil {
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
			sb.Append(self.GetNormalizedReferenceName())
		} else {
			sb.Append(*self.GetReferenceName())
		}
		return
	}

	if !p.IsPrettyPrint() || sb.printing.contains(p.GetJSTypeBase()) {
		if sb.IsForAnnotations() {
			sb.Append("?")
		} else {
			sb.Append("{...}")
		}
		return
	}

	// Don't pretty print recursively. It would cause infinite recursion.
	sb.printing.push(p.GetJSTypeBase())

	// Sort and deduplicate properties across the prototype chain (subtype wins).
	var propsBuf [8]PrototypeObjectType_PropertyForPrettyPrinting
	properties := propsBuf[:0]
	propEq := func(a, b PrototypeObjectType_PropertyForPrettyPrinting) bool {
		return PrototypeObjectType_PROP_COMPARATOR(a, b) == 0
	}
	for current := self; current != nil; current = current.GetImplicitPrototype() {
		if current.IsNativeObjectType() || len(properties) > PrototypeObjectType_MAX_PRETTY_PRINTED_PROPERTIES {
			break
		}

		for _, stringKey := range current.GetOwnPropertyNames() {
			key := NewProperty_StringKey(stringKey)
			properties = append(properties, PrototypeObjectType_PropertyForPrettyPrinting{key, current.GetPropertyTypeByKey(key)})
		}
		for _, symbolKey := range current.GetOwnPropertyKnownSymbols() {
			key := NewProperty_SymbolKey(symbolKey)
			properties = append(properties, PrototypeObjectType_PropertyForPrettyPrinting{key, current.GetPropertyTypeByKey(key)})
		}
		slices.SortStableFunc(properties, PrototypeObjectType_PROP_COMPARATOR)
		properties = slices.CompactFunc(properties, propEq)
	}

	multiline := !sb.IsForAnnotations() && len(properties) > 1
	sb.Append("{")
	sb.indentDepth++
	if multiline {
		sb.BreakLineAndIndent()
	}

	for idx, property := range properties {
		i := idx + 1

		if !sb.IsForAnnotations() && i > PrototypeObjectType_MAX_PRETTY_PRINTED_PROPERTIES {
			sb.Append("...")
			break
		}

		property.AppendTo(sb)
		if i < len(properties) {
			sb.Append(",")
			if multiline {
				sb.BreakLineAndIndent()
			} else {
				sb.Append(" ")
			}
		}
	}
	sb.indentDepth--
	if multiline {
		sb.BreakLineAndIndent()
	}
	sb.Append("}")

	sb.printing.pop()
}

// Java: PrototypeObjectType.setPrettyPrint(boolean)
func (p *PrototypeObjectTypeBase) SetPrettyPrint(prettyPrint bool) {
	if prettyPrint {
		p.protoFlags |= protoFlagPrettyPrint
	} else {
		p.protoFlags &^= protoFlagPrettyPrint
	}
}

// Java: PrototypeObjectType.isPrettyPrint()
func (p *PrototypeObjectTypeBase) IsPrettyPrint() bool {
	return (p.protoFlags & protoFlagPrettyPrint) != 0
}

// Java: PrototypeObjectType.getConstructor()
func (p *PrototypeObjectTypeBase) GetConstructor() FunctionType {
	return nil
}

// Java: PrototypeObjectType.getImplicitPrototype()
func (p *PrototypeObjectTypeBase) GetImplicitPrototype() ObjectType {
	return objectTypeOf(p.implicitPrototypeFallback)
}

// SetImplicitPrototype: this should only be reset on the
// FunctionPrototypeType, only to fix an incorrectly established prototype
// chain due to the user having a mismatch in super class declaration, and
// only before properties on that type are processed.
//
// Java: PrototypeObjectType.setImplicitPrototype(ObjectType)
func (p *PrototypeObjectTypeBase) SetImplicitPrototype(implicitPrototype ObjectType) {
	jutil.CheckState(!p.selfObj().HasCachedValues())
	p.implicitPrototypeFallback = objTypeBase(implicitPrototype)
	if reg := p.GetRegistry(); reg != nil && !reg.parallelReadOnly {
		// Go-only: in parallel read-only mode the cache is bypassed and the
		// generation is bumped when the mode ends (SetParallelReadOnly).
		reg.leastSupertypeGen++
	}
	if implicitPrototype != nil {
		p.self().MaybeLoosenTypecheckingDueToForwardReferencedSupertype(implicitPrototype)
	}
}

// Java: PrototypeObjectType.getTemplateParamCount()
func (p *PrototypeObjectTypeBase) GetTemplateParamCount() int {
	return int(p.templateParamCount)
}

// Java: ObjectType.hasReferenceName()
func (p *PrototypeObjectTypeBase) HasReferenceName() bool {
	return p.classNameID != 0 || p.ownerFunction != nil
}

func (p *PrototypeObjectTypeBase) referenceNameString() (string, bool) {
	if p.classNameID != 0 {
		return rhinoStringFromID(p.classNameID), true
	}
	if p.ownerFunction != nil {
		if ownerName, ok := p.ownerFunction.referenceNameString(); ok {
			return ownerName + ".prototype", true
		}
		return "null.prototype", true
	}
	return "", false
}

// Java: PrototypeObjectType.getReferenceName()
func (p *PrototypeObjectTypeBase) GetReferenceName() *string {
	if p.classNameID != 0 {
		return rhinoStringPtrFromID(p.classNameID)
	}
	if name, ok := p.referenceNameString(); ok {
		return &name
	}
	return nil
}

// Java: PrototypeObjectType.isAnonymous()
func (p *PrototypeObjectTypeBase) IsAnonymous() bool {
	return (p.protoFlags & protoFlagAnonymous) != 0
}

// IsNativeObjectType reports whether this is a built-in object.
//
// Java: PrototypeObjectType.isNativeObjectType()
func (p *PrototypeObjectTypeBase) IsNativeObjectType() bool {
	return (p.protoFlags & protoFlagNative) != 0
}

// Java: PrototypeObjectType.setOwnerFunction(FunctionType)
func (p *PrototypeObjectTypeBase) SetOwnerFunction(type_ FunctionType) {
	jutil.CheckState(p.ownerFunction == nil || type_ == nil)
	p.ownerFunction = type_
}

// Java: PrototypeObjectType.getOwnerFunction()
func (p *PrototypeObjectTypeBase) GetOwnerFunction() FunctionType {
	return p.ownerFunction
}

// Java: PrototypeObjectType.getCtorImplementedInterfaces()
func (p *PrototypeObjectTypeBase) GetCtorImplementedInterfaces() []ObjectType {
	if p.self().IsFunctionPrototypeType() {
		return p.selfObj().GetOwnerFunction().GetImplementedInterfaces()
	}
	return []ObjectType{}
}

// Java: PrototypeObjectType.getCtorExtendedInterfaces()
func (p *PrototypeObjectTypeBase) GetCtorExtendedInterfaces() []ObjectType {
	if p.self().IsFunctionPrototypeType() {
		return p.selfObj().GetOwnerFunction().GetExtendedInterfaces()
	}
	return []ObjectType{}
}

// Java: PrototypeObjectType.resolveInternal(ErrorReporter)
func (p *PrototypeObjectTypeBase) ResolveInternal(reporter ErrorReporter) JSType {
	implicitPrototype := p.selfObj().GetImplicitPrototype()
	if implicitPrototype != nil {
		p.implicitPrototypeFallback = objTypeBase(implicitPrototype.Resolve(reporter).(ObjectType))
	}
	p.properties.ensureSorted()
	for i := 0; i < len(p.properties.entries()); i++ {
		prop := p.properties.entries()[i]
		prop.SetType(JSType_SafeResolve(prop.GetType(), reporter))
	}
	return p.self()
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
	if constraintObj.GetObjectTypeBase().objFlags&objKindMask == objKindRecord {
		pm := constraintObj.GetPropertyMap()
		pm.ensureSorted()
		for _, prop := range pm.entries() {
			key := NewProperty_StringKeyByID(uint32(prop.nameID))
			propType := prop.GetType()
			if !self.IsPropertyTypeDeclaredByKey(key) {
				typeToInfer := propType
				if !p.self().HasPropertyByKey(key) {
					typeToInfer = p.GetNativeType(JSTypeNative_VOID_TYPE).GetLeastSupertype(propType)
				}
				self.DefineInferredPropertyByKey(key, typeToInfer, NilNode)
			}
		}
		return
	}
	for _, prop := range constraintObj.GetOwnPropertyNames() {
		propType := constraintObj.GetPropertyType(prop)
		if !self.IsPropertyTypeDeclared(prop) {
			typeToInfer := propType
			if !p.self().HasProperty(prop) {
				typeToInfer = p.GetNativeType(JSTypeNative_VOID_TYPE).GetLeastSupertype(propType)
			}
			self.DefineInferredProperty(prop, typeToInfer, NilNode)
		}
	}
}

// Java: PrototypeObjectType.recursionUnsafeHashCode()
func (p *PrototypeObjectTypeBase) RecursionUnsafeHashCode(*JSTypeStack) int32 {
	if p.self().IsStructuralType() {
		var classHash int32
		if p.classNameID != 0 {
			classHash = jutil.StringHashCode(rhinoStringFromID(p.classNameID))
		}
		return 31*(31+classHash) + p.properties.HashCode()
	} else {
		return jutil.IdentityHashCode(p.self())
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
