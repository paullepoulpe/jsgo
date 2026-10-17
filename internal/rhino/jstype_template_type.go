// Port of com/google/javascript/rhino/jstype/TemplateType.java @ closure-compiler v20261015

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

// TemplateType is a placeholder type, used as keys in TemplateTypeMaps.
//
// For functions with function(this: T, ...) and T as arguments, type
// inference will set the type of this on a function literal argument to the
// actual type of T.
type TemplateType struct {
	ProxyObjectTypeBase

	name               string
	bound              JSType
	typeTransformation *Node
}

// Java: TemplateType.TYPE_CLASS
const TemplateType_TYPE_CLASS = JSTypeClass_TEMPLATE

// Java: TemplateType(JSTypeRegistry, String)
func NewTemplateType(registry *JSTypeRegistry, name string) *TemplateType {
	return NewTemplateTypeWithBoundAndTransformation(registry, name, nil, nil)
}

// Java: TemplateType(JSTypeRegistry, String, JSType)
func NewTemplateTypeWithBound(registry *JSTypeRegistry, name string, bound JSType) *TemplateType {
	return NewTemplateTypeWithBoundAndTransformation(registry, name, bound, nil)
}

// Java: TemplateType(JSTypeRegistry, String, Node)
func NewTemplateTypeWithTransformation(registry *JSTypeRegistry, name string, typeTransformation *Node) *TemplateType {
	return NewTemplateTypeWithBoundAndTransformation(registry, name, nil, typeTransformation)
}

// Java: TemplateType(JSTypeRegistry, String, JSType, Node)
func NewTemplateTypeWithBoundAndTransformation(
	registry *JSTypeRegistry, name string, bound JSType, typeTransformation *Node) *TemplateType {
	t := &TemplateType{}
	referenced := bound
	if bound == nil {
		referenced = registry.GetNativeObjectType(JSTypeNative_UNKNOWN_TYPE)
	}
	t.InitProxyObjectTypeBase(t, registry, referenced)
	t.name = name
	if bound == nil {
		t.bound = registry.GetNativeObjectType(JSTypeNative_UNKNOWN_TYPE)
	} else {
		t.bound = bound
	}
	t.typeTransformation = typeTransformation

	registry.GetResolver().ResolveIfClosed(t, TemplateType_TYPE_CLASS)
	return t
}

// Java: TemplateType.getTypeClass()
func (t *TemplateType) GetTypeClass() JSTypeClass {
	return TemplateType_TYPE_CLASS
}

// Java: TemplateType.getReferenceName()
func (t *TemplateType) GetReferenceName() *string {
	return &t.name
}

// Java: TemplateType.appendTo(TypeStringBuilder)
func (t *TemplateType) AppendTo(sb *TypeStringBuilder) {
	if t.bound == JSType(t.registry.GetNativeObjectType(JSTypeNative_UNKNOWN_TYPE)) {
		// This is an unbound template, don't print it's bound
		sb.Append(t.name)
	} else {
		sb.Append(t.name).Append(" extends ").AppendType(t.bound)
	}
}

// Java: TemplateType.toMaybeTemplateType()
func (t *TemplateType) ToMaybeTemplateType() *TemplateType {
	return t
}

// Java: TemplateType.hasAnyTemplateTypesInternal()
func (t *TemplateType) HasAnyTemplateTypesInternal() bool {
	return true
}

// Java: TemplateType.visit(Visitor<T>)
func (t *TemplateType) Visit(visitor Visitor) any {
	return visitor.CaseTemplateType(t)
}

// Java: TemplateType.visit(RelationshipVisitor<T>, JSType)
func (t *TemplateType) VisitRelationship(visitor RelationshipVisitor, that JSType) any {
	return visitor.CaseTemplateType(t, that)
}

// Java: TemplateType.isTypeTransformation()
func (t *TemplateType) IsTypeTransformation() bool {
	return t.typeTransformation != nil
}

// Java: TemplateType.getTypeTransformation()
func (t *TemplateType) GetTypeTransformation() *Node {
	return t.typeTransformation
}

// Java: TemplateType.getBound()
func (t *TemplateType) GetBound() JSType {
	return t.bound
}

// Java: TemplateType.setBound(JSType)
func (t *TemplateType) SetBound(bound JSType) {
	t.bound = bound
	t.SetReferencedType(bound)
}

// Java: TemplateType.setValidator(Predicate<JSType>)
func (t *TemplateType) SetValidator(validator func(JSType) bool) bool {
	// ProxyObjectType will delegate to the unknown type's setValidator, which
	// we don't want. Some validator functions special-case template types.
	// Note that this.isUnknownType() is still true, so this override only
	// affects validators that treat TemplateTypes differently from a random
	// type for which isUnknownType() is true. (e.g.
	// FunctionTypeBuilder#ExtendedTypeValidator)
	return validator(t)
}

// ContainsCycle returns whether or not there is a cycle in the reference
// chain of this type.
//
// Java: TemplateType.containsCycle()
func (t *TemplateType) ContainsCycle() bool {
	// By passing in an unreachable type `null` as the target, the visitor will
	// only return if a cycle is found or a terminating type is reached.
	typeVisitor := NewContainsUpperBoundSuperTypeVisitor(nil)
	return t.Visit(typeVisitor) == ContainsUpperBoundSuperTypeVisitor_Result_CYCLE
}

// Java: TemplateType.recursionUnsafeHashCode()
func (t *TemplateType) RecursionUnsafeHashCode() int32 {
	// Since TemplateType has identity semantics for equality (see
	// EqualityChecker), use identity hash code to avoid spurious hash
	// conflicts
	return jutil.IdentityHashCode(t)
}
