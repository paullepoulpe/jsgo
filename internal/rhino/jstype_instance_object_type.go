// Port of com/google/javascript/rhino/jstype/InstanceObjectType.java @ closure-compiler v20261021

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

// InstanceObjectType is an object type that is an instance of some function
// constructor.
type InstanceObjectType struct {
	PrototypeObjectTypeBase

	constructor *FunctionTypeBase
}

// Java: InstanceObjectType.TYPE_CLASS
const InstanceObjectType_TYPE_CLASS = JSTypeClass_INSTANCE_OBJECT

// Java: InstanceObjectType(Builder)
func NewInstanceObjectType(builder *InstanceObjectType_Builder) *InstanceObjectType {
	t := &InstanceObjectType{}
	t.objFlags = objKindInstance
	t.InitPrototypeObjectTypeBase(t, &builder.PrototypeObjectType_Builder)
	t.constructor = jutil.CheckNotNull(builder.constructor)

	t.GetRegistry().GetResolver().ResolveIfClosed(t, InstanceObjectType_TYPE_CLASS)
	return t
}

// InstanceObjectType_Builder is InstanceObjectType.Builder.
type InstanceObjectType_Builder struct {
	PrototypeObjectType_Builder

	constructor FunctionType
}

// Java: InstanceObjectType.Builder(JSTypeRegistry)
func NewInstanceObjectType_Builder(registry *JSTypeRegistry) *InstanceObjectType_Builder {
	b := &InstanceObjectType_Builder{}
	b.InitPrototypeObjectType_Builder(registry)
	return b
}

// Java: PrototypeObjectType.Builder.setName(String)
func (b *InstanceObjectType_Builder) SetName(x *string) *InstanceObjectType_Builder {
	b.classNameID = rhinoStringIDFromPtr(x)
	return b
}

// Java: PrototypeObjectType.Builder.setImplicitPrototype(ObjectType)
func (b *InstanceObjectType_Builder) SetImplicitPrototype(x ObjectType) *InstanceObjectType_Builder {
	b.implicitPrototype = x
	return b
}

// Java: PrototypeObjectType.Builder.setNative(boolean)
func (b *InstanceObjectType_Builder) SetNative(x bool) *InstanceObjectType_Builder {
	b.nativeType = x
	return b
}

// Java: PrototypeObjectType.Builder.setAnonymous(boolean)
func (b *InstanceObjectType_Builder) SetAnonymous(x bool) *InstanceObjectType_Builder {
	b.anonymousType = x
	return b
}

// Java: PrototypeObjectType.Builder.setTemplateTypeMap(TemplateTypeMap)
func (b *InstanceObjectType_Builder) SetTemplateTypeMap(x *TemplateTypeMap) *InstanceObjectType_Builder {
	b.templateTypeMap = x
	return b
}

// Java: PrototypeObjectType.Builder.setTemplateParamCount(int)
func (b *InstanceObjectType_Builder) SetTemplateParamCount(x int) *InstanceObjectType_Builder {
	b.templateParamCount = x
	return b
}

// Java: InstanceObjectType.Builder.setConstructor(FunctionType)
func (b *InstanceObjectType_Builder) SetConstructor(x FunctionType) *InstanceObjectType_Builder {
	b.constructor = x
	return b
}

// Java: InstanceObjectType.Builder.build()
func (b *InstanceObjectType_Builder) Build() *InstanceObjectType {
	return NewInstanceObjectType(b)
}

// Java: InstanceObjectType.builderForCtor(FunctionType)
func InstanceObjectType_BuilderForCtor(ctor FunctionType) *InstanceObjectType_Builder {
	return &InstanceObjectType_Builder{
		PrototypeObjectType_Builder: PrototypeObjectType_Builder{
			registry:          ctor.GetJSTypeBase().GetRegistry(),
			classNameID:       rhinoStringIDFromPtr(ctor.GetReferenceName()),
			implicitPrototype: nil, // This isn't shared by a function and it's instances.
			nativeType:        ctor.IsNativeObjectType(),
			anonymousType:     ctor.IsAnonymous(),
			templateTypeMap:   ctor.GetTemplateTypeMap(),
			// Recall that in ES5 code, instance and constructor template parameters
			// were indistinguishable. That asumption is maintained here by deault,
			// but later code may may overwrite the template parameter count.
			templateParamCount: ctor.GetTemplateParamCount(),
		},
		constructor: ctor,
	}
}

// Java: InstanceObjectType.getTypeClass()
func (t *InstanceObjectType) GetTypeClass() JSTypeClass {
	return InstanceObjectType_TYPE_CLASS
}

// Java: ObjectType.hasReferenceName()
func (t *InstanceObjectType) HasReferenceName() bool {
	return t.constructor.HasReferenceName()
}

// Java: InstanceObjectType.getReferenceName()
func (t *InstanceObjectType) GetReferenceName() *string {
	return t.GetConstructor().GetReferenceName()
}

// Java: InstanceObjectType.getImplicitPrototype()
func (t *InstanceObjectType) GetImplicitPrototype() ObjectType {
	return t.GetConstructor().GetPrototype()
}

// Java: InstanceObjectType.getConstructor()
func (t *InstanceObjectType) GetConstructor() FunctionType {
	return t.constructor
}

// Java: InstanceObjectType.defineProperty(Property.Key, JSType, boolean, Node)
func (t *InstanceObjectType) DefinePropertyByKey(name Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	proto := t.GetImplicitPrototype()
	if proto != nil && proto.HasOwnDeclaredPropertyByKey(name) {
		return false
	}
	return t.PrototypeObjectTypeBase.DefinePropertyByKey(name, type_, inferred, propertyNode)
}

// Java: InstanceObjectType.appendTo(TypeStringBuilder)
func (t *InstanceObjectType) AppendTo(sb *TypeStringBuilder) {
	ctor := t.constructor
	if !ctor.HasReferenceName() {
		t.PrototypeObjectTypeBase.AppendTo(sb)
		return
	} else if sb.IsForAnnotations() {
		sb.Append(ctor.GetNormalizedReferenceName())
		return
	}

	name := *ctor.GetReferenceName()
	if name == "" {
		n := ctor.GetSource()
		sourceFileName := "unknown"
		lineno := 0
		if !n.IsNil() {
			sourceFileName = jutil.ToJavaString(n.GetSourceFileName())
			lineno = n.GetLineno()
		}
		sb.Append("<anonymous@")
		sb.Append(sourceFileName)
		sb.Append(":")
		sb.Append(jutil.ToJavaString(lineno))
		sb.Append(">")
	}
	sb.Append(name)
}

func (t *InstanceObjectType) isNativeClass(name string) bool {
	c := t.constructor
	return (c.protoFlags&protoFlagNative) != 0 && c.classNameID != 0 && rhinoStringFromID(c.classNameID) == name
}

// Java: InstanceObjectType.isTheObjectType()
func (t *InstanceObjectType) IsTheObjectType() bool {
	return t.isNativeClass("Object")
}

// Java: InstanceObjectType.isInstanceType()
func (t *InstanceObjectType) IsInstanceType() bool {
	return true
}

// Java: InstanceObjectType.isArrayType()
func (t *InstanceObjectType) IsArrayType() bool {
	return t.isNativeClass("Array")
}

// Java: InstanceObjectType.isReadonlyArrayType()
func (t *InstanceObjectType) IsReadonlyArrayType() bool {
	return t.isNativeClass("ReadonlyArray")
}

// Java: InstanceObjectType.isBigIntObjectType()
func (t *InstanceObjectType) IsBigIntObjectType() bool {
	return t.isNativeClass("BigInt")
}

// Java: InstanceObjectType.isStringObjectType()
func (t *InstanceObjectType) IsStringObjectType() bool {
	return t.isNativeClass("String")
}

// Java: InstanceObjectType.isSymbolObjectType()
func (t *InstanceObjectType) IsSymbolObjectType() bool {
	return t.isNativeClass("Symbol")
}

// Java: InstanceObjectType.isBooleanObjectType()
func (t *InstanceObjectType) IsBooleanObjectType() bool {
	return t.isNativeClass("Boolean")
}

// Java: InstanceObjectType.isNumberObjectType()
func (t *InstanceObjectType) IsNumberObjectType() bool {
	return t.isNativeClass("Number")
}

// Java: InstanceObjectType.isDateType()
func (t *InstanceObjectType) IsDateType() bool {
	return t.isNativeClass("Date")
}

// Java: InstanceObjectType.isRegexpType()
func (t *InstanceObjectType) IsRegexpType() bool {
	return t.isNativeClass("RegExp")
}

// Java: InstanceObjectType.isNominalType()
func (t *InstanceObjectType) IsNominalType() bool {
	return t.HasReferenceName()
}

// Java: InstanceObjectType.recursionUnsafeHashCode()
func (t *InstanceObjectType) RecursionUnsafeHashCode(hs *JSTypeStack) int32 {
	if t.HasReferenceName() {
		return NamedType_NominalHashCode(t)
	} else {
		// Java: super.hashCode(), the final JSType.hashCode(): this type is in
		// progress, so -1.
		return t.PrototypeObjectTypeBase.GetJSTypeBase().hashCodeIn(hs)
	}
}

// Java: InstanceObjectType.getCtorImplementedInterfaces()
func (t *InstanceObjectType) GetCtorImplementedInterfaces() []ObjectType {
	return t.GetConstructor().GetImplementedInterfaces()
}

// Java: InstanceObjectType.getCtorExtendedInterfaces()
func (t *InstanceObjectType) GetCtorExtendedInterfaces() []ObjectType {
	return t.GetConstructor().GetExtendedInterfaces()
}

// Java: InstanceObjectType.resolveInternal(ErrorReporter)
func (t *InstanceObjectType) ResolveInternal(reporter ErrorReporter) JSType {
	resolved := t.PrototypeObjectTypeBase.ResolveInternal(reporter)
	// No need to update the constructor with the result of safeResolve, since
	// the constructor will always resolve to itself - it's never a forward
	// reference because to create this InstanceObjectType, we must have the
	// constructor first. But we do need to call resolve on the constructor to
	// resolve its fields and implemented/extended interfaces, if any.
	_ = JSType_SafeResolve(t.constructor, reporter)
	return resolved
}

// jsTypeStringEquals is literal.equals(s) for a nullable String s.
func jsTypeStringEquals(literal string, s *string) bool {
	return s != nil && *s == literal
}
