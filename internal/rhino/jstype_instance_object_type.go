// Port of com/google/javascript/rhino/jstype/InstanceObjectType.java @ closure-compiler v20261015

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
	t := builder.registry.allocInstanceObjectType()
	t.InitPrototypeObjectTypeBase(t, builder.GetBuilderFields())
	t.constructor = jutil.CheckNotNull(builder.constructor).GetFunctionTypeBase()

	t.registry.GetResolver().ResolveIfClosed(t, InstanceObjectType_TYPE_CLASS)
	return t
}

// InstanceObjectType_Builder is InstanceObjectType.Builder.
type InstanceObjectType_Builder struct {
	PrototypeObjectType_Builder[*InstanceObjectType_Builder]

	constructor FunctionType
}

// Java: InstanceObjectType.Builder(JSTypeRegistry)
func NewInstanceObjectType_Builder(registry *JSTypeRegistry) *InstanceObjectType_Builder {
	b := &InstanceObjectType_Builder{}
	b.InitPrototypeObjectType_Builder(b, registry)
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
	return NewInstanceObjectType_Builder(ctor.GetJSTypeBase().registry).
		SetName(ctor.GetReferenceName()).
		SetImplicitPrototype(nil). // This isn't shared by a function and it's instances.
		SetNative(ctor.IsNativeObjectType()).
		SetAnonymous(ctor.IsAnonymous()).
		SetTemplateTypeMap(ctor.GetTemplateTypeMap()).
		// Recall that in ES5 code, instance and constructor template parameters
		// were indistinguishable. That asumption is maintained here by deault,
		// but later code may may overwrite the template parameter count.
		SetTemplateParamCount(ctor.GetTemplateParamCount()).
		SetConstructor(ctor)
}

// Java: InstanceObjectType.getTypeClass()
func (t *InstanceObjectType) GetTypeClass() JSTypeClass {
	return InstanceObjectType_TYPE_CLASS
}

// Java: ObjectType.hasReferenceName()
func (t *InstanceObjectType) HasReferenceName() bool {
	return t.constructor.selfFunc().HasReferenceName()
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
	return t.constructor.selfFunc()
}

// Java: InstanceObjectType.defineProperty(Property.Key, JSType, boolean, Node)
func (t *InstanceObjectType) DefinePropertyByKey(name Property_Key, type_ JSType, inferred bool, propertyNode *Node) bool {
	proto := t.GetImplicitPrototype()
	if proto != nil && proto.HasOwnDeclaredPropertyByKey(name) {
		return false
	}
	return t.PrototypeObjectTypeBase.DefinePropertyByKey(name, type_, inferred, propertyNode)
}

// Java: InstanceObjectType.appendTo(TypeStringBuilder)
func (t *InstanceObjectType) AppendTo(sb *TypeStringBuilder) {
	ctor := t.constructor.selfFunc()
	if !ctor.HasReferenceName() {
		t.PrototypeObjectTypeBase.AppendTo(sb)
		return
	} else if sb.IsForAnnotations() {
		sb.Append(*ctor.GetNormalizedReferenceName())
		return
	}

	name := *ctor.GetReferenceName()
	if name == "" {
		n := ctor.GetSource()
		sourceFileName := "unknown"
		lineno := 0
		if n != nil {
			sourceFileName = jutil.ToJavaString(n.GetSourceFileName())
			lineno = n.GetLineno()
		}
		sb.Append("<anonymous@").
			Append(sourceFileName).
			Append(":").
			Append(jutil.ToJavaString(lineno)).
			Append(">")
	}
	sb.Append(name)
}

// Java: InstanceObjectType.isTheObjectType()
func (t *InstanceObjectType) IsTheObjectType() bool {
	return t.GetConstructor().IsNativeObjectType() && jsTypeStringEquals("Object", t.GetReferenceName())
}

// Java: InstanceObjectType.isInstanceType()
func (t *InstanceObjectType) IsInstanceType() bool {
	return true
}

// Java: InstanceObjectType.isArrayType()
func (t *InstanceObjectType) IsArrayType() bool {
	return t.GetConstructor().IsNativeObjectType() && jsTypeStringEquals("Array", t.GetReferenceName())
}

// Java: InstanceObjectType.isReadonlyArrayType()
func (t *InstanceObjectType) IsReadonlyArrayType() bool {
	return t.GetConstructor().IsNativeObjectType() && jsTypeStringEquals("ReadonlyArray", t.GetReferenceName())
}

// Java: InstanceObjectType.isBigIntObjectType()
func (t *InstanceObjectType) IsBigIntObjectType() bool {
	return t.GetConstructor().IsNativeObjectType() && jsTypeStringEquals("BigInt", t.GetReferenceName())
}

// Java: InstanceObjectType.isStringObjectType()
func (t *InstanceObjectType) IsStringObjectType() bool {
	return t.GetConstructor().IsNativeObjectType() && jsTypeStringEquals("String", t.GetReferenceName())
}

// Java: InstanceObjectType.isSymbolObjectType()
func (t *InstanceObjectType) IsSymbolObjectType() bool {
	return t.GetConstructor().IsNativeObjectType() && jsTypeStringEquals("Symbol", t.GetReferenceName())
}

// Java: InstanceObjectType.isBooleanObjectType()
func (t *InstanceObjectType) IsBooleanObjectType() bool {
	return t.GetConstructor().IsNativeObjectType() && jsTypeStringEquals("Boolean", t.GetReferenceName())
}

// Java: InstanceObjectType.isNumberObjectType()
func (t *InstanceObjectType) IsNumberObjectType() bool {
	return t.GetConstructor().IsNativeObjectType() && jsTypeStringEquals("Number", t.GetReferenceName())
}

// Java: InstanceObjectType.isDateType()
func (t *InstanceObjectType) IsDateType() bool {
	return t.GetConstructor().IsNativeObjectType() && jsTypeStringEquals("Date", t.GetReferenceName())
}

// Java: InstanceObjectType.isRegexpType()
func (t *InstanceObjectType) IsRegexpType() bool {
	return t.GetConstructor().IsNativeObjectType() && jsTypeStringEquals("RegExp", t.GetReferenceName())
}

// Java: InstanceObjectType.isNominalType()
func (t *InstanceObjectType) IsNominalType() bool {
	return t.HasReferenceName()
}

// Java: InstanceObjectType.recursionUnsafeHashCode()
func (t *InstanceObjectType) RecursionUnsafeHashCode() int32 {
	if t.HasReferenceName() {
		return NamedType_NominalHashCode(t)
	} else {
		return t.PrototypeObjectTypeBase.HashCode()
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
	_ = JSType_SafeResolve(t.constructor.selfFunc(), reporter)
	return resolved
}

// jsTypeStringEquals is literal.equals(s) for a nullable String s.
func jsTypeStringEquals(literal string, s *string) bool {
	return s != nil && *s == literal
}
