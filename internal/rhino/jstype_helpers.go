// Port of com/google/javascript/rhino/jstype/EquivalenceMethod.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/JSTypeClass.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/BooleanLiteralSet.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/JSTypeIterations.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/NominalTypeBuilder.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/TypeStringBuilder.java @ closure-compiler v20261021

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
 *   Nick Santos
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

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
 *   Ben Lickly
 *   Dimitris Vardoulakis
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// EquivalenceMethod represents different ways for comparing equality among
// types.
type EquivalenceMethod int32

const (
	// EquivalenceMethod_IDENTITY indicates that the two types should behave
	// exactly the same under all type operations.
	//
	// Thus, {string} != {?} and {Unresolved} != {?}
	EquivalenceMethod_IDENTITY EquivalenceMethod = iota
	// EquivalenceMethod_DATA_FLOW indicates that the two types are almost
	// exactly the same, and that a data flow analysis algorithm comparing them
	// should consider them equal.
	//
	// In traditional type inference, the types form a finite lattice, and
	// this ensures that type inference will terminate.
	//
	// In our type system, the unknown types do not obey the lattice rules. So
	// if we continue to perform inference over the unknown types, we may never
	// terminate.
	//
	// By treating all unknown types as equivalent for the purposes of data
	// flow analysis, we ensure that the algorithm will terminate.
	//
	// Thus, {string} != {?} and {Unresolved} ~= {?}
	EquivalenceMethod_DATA_FLOW
)

var equivalenceMethod_names = [...]string{"IDENTITY", "DATA_FLOW"}

// String is Enum.name().
func (m EquivalenceMethod) String() string { return equivalenceMethod_names[m] }

// JSTypeClass is the set of Java classes that are subtypes of JSType.
//
// Each class reports the corresponding element as its instance
// JSType#getTypeClass method.
//
// These values are not intended to be compatible with Liskov substitution;
// JSTypes have exactly one associated JSTypeClass and it must not affect the
// outcome of operations such as Object#equals. Use of this enum is
// tantamount to reflection (e.g. Object#getClass or instanceof) and should
// only be used where such a comparison is appropriate and desirable.
//
// Use of this enum is preferred to true reflection, when possible, for
// simplicity. This is safe because this package controls all JSType
// subclasses; they can all be known a priori.
//
// Abstract classes are not listed here because there would be no way to use
// those elements.
type JSTypeClass int32

const (
	JSTypeClass_ALL JSTypeClass = iota
	JSTypeClass_ARROW
	JSTypeClass_BOOLEAN
	JSTypeClass_BIGINT
	JSTypeClass_ENUM
	JSTypeClass_ENUM_ELEMENT
	JSTypeClass_FUNCTION
	JSTypeClass_INSTANCE_OBJECT
	JSTypeClass_NAMED
	JSTypeClass_NO
	JSTypeClass_NO_OBJECT
	JSTypeClass_NO_RESOLVED
	JSTypeClass_NULL
	JSTypeClass_NUMBER
	JSTypeClass_PROTOTYPE_OBJECT
	JSTypeClass_PROXY_OBJECT
	JSTypeClass_RECORD
	JSTypeClass_STRING
	JSTypeClass_SYMBOL
	JSTypeClass_WELL_KNOWN_SYMBOL
	JSTypeClass_TEMPLATE
	JSTypeClass_TEMPLATIZED
	JSTypeClass_UNION
	JSTypeClass_UNKNOWN
	JSTypeClass_VOID
)

var jsTypeClass_names = [...]string{
	"ALL", "ARROW", "BOOLEAN", "BIGINT", "ENUM", "ENUM_ELEMENT", "FUNCTION", "INSTANCE_OBJECT",
	"NAMED", "NO", "NO_OBJECT", "NO_RESOLVED", "NULL", "NUMBER", "PROTOTYPE_OBJECT",
	"PROXY_OBJECT", "RECORD", "STRING", "SYMBOL", "WELL_KNOWN_SYMBOL", "TEMPLATE", "TEMPLATIZED",
	"UNION", "UNKNOWN", "VOID",
}

// String is Enum.name().
func (c JSTypeClass) String() string { return jsTypeClass_names[c] }

// Java: JSTypeClass.isTypeOf(JSType)
func (c JSTypeClass) IsTypeOf(t JSType) bool {
	return c == t.GetTypeClass()
}

// BooleanLiteralSet is a set in the domain {true,false}. There are four
// possible sets: {}, {true}, {false}, {true,false}.
type BooleanLiteralSet int32

const (
	BooleanLiteralSet_EMPTY BooleanLiteralSet = iota
	BooleanLiteralSet_TRUE
	BooleanLiteralSet_FALSE
	BooleanLiteralSet_BOTH
)

var booleanLiteralSet_names = [...]string{"EMPTY", "TRUE", "FALSE", "BOTH"}

// String is Enum.name().
func (s BooleanLiteralSet) String() string { return booleanLiteralSet_names[s] }

// Ordinal is Enum.ordinal().
func (s BooleanLiteralSet) Ordinal() int { return int(s) }

// Java: BooleanLiteralSet.fromOrdinal(int)
func (s BooleanLiteralSet) FromOrdinal(ordinal int) BooleanLiteralSet {
	switch ordinal {
	case 0:
		return BooleanLiteralSet_EMPTY
	case 1:
		return BooleanLiteralSet_TRUE
	case 2:
		return BooleanLiteralSet_FALSE
	case 3:
		return BooleanLiteralSet_BOTH
	}
	panic(jutil.NewIllegalArgumentExceptionMsg("Ordinal: " + jutil.ToJavaString(ordinal)))
}

// Intersection computes the intersection of this set and that.
//
// Java: BooleanLiteralSet.intersection(BooleanLiteralSet)
func (s BooleanLiteralSet) Intersection(that BooleanLiteralSet) BooleanLiteralSet {
	return s.FromOrdinal(s.Ordinal() & that.Ordinal())
}

// Union computes the union of this set and that.
//
// Java: BooleanLiteralSet.union(BooleanLiteralSet)
func (s BooleanLiteralSet) Union(that BooleanLiteralSet) BooleanLiteralSet {
	return s.FromOrdinal(s.Ordinal() | that.Ordinal())
}

// Contains returns whether this contains the given literal value.
//
// Java: BooleanLiteralSet.contains(boolean)
func (s BooleanLiteralSet) Contains(literalValue bool) bool {
	switch s.Ordinal() {
	case 0:
		return false
	case 1:
		return literalValue
	case 2:
		return !literalValue
	case 3:
		return true
	}
	panic(jutil.NewIndexOutOfBoundsException("Ordinal: " + jutil.ToJavaString(s.Ordinal())))
}

// BooleanLiteralSet_Get returns the singleton set {literalValue}.
//
// Java: BooleanLiteralSet.get(boolean)
func BooleanLiteralSet_Get(literalValue bool) BooleanLiteralSet {
	if literalValue {
		return BooleanLiteralSet_TRUE
	}
	return BooleanLiteralSet_FALSE
}

// ToTri converts to a Tri.
//
// Java: BooleanLiteralSet.toTri()
func (s BooleanLiteralSet) ToTri() base.Tri {
	if s == BooleanLiteralSet_TRUE {
		return base.Tri_TRUE
	} else if s == BooleanLiteralSet_FALSE {
		return base.Tri_FALSE
	}
	return base.Tri_UNKNOWN
}

// JSTypeIterations holds Stream-like functions for efficient querying of
// collections of JSTypes.
//
// These methods are designed to minimize allocations since they are
// expected to be called very often. That's why they don't:
//
//   - instantiate Iterators
//   - instantiate Streams
//   - (un)box primitives
//   - expect closure-generating lambdas
type JSTypeIterations struct{}

// JSTypeIterations_AnyTypeMatches returns whether any element of types
// matches predicate.
//
// Java: JSTypeIterations.anyTypeMatches(Predicate<? super JSType>, ImmutableList<? extends JSType>)
func JSTypeIterations_AnyTypeMatches[S JSType](predicate func(JSType) bool, types []S) bool {
	for i := 0; i < len(types); i++ {
		if predicate(types[i]) {
			return true
		}
	}
	return false
}

// JSTypeIterations_AnyTypeMatchesInUnion returns whether any element of
// union matches predicate.
//
// Java: JSTypeIterations.anyTypeMatches(Predicate<? super JSType>, UnionType)
func JSTypeIterations_AnyTypeMatchesInUnion(predicate func(JSType) bool, union *UnionType) bool {
	return JSTypeIterations_AnyTypeMatches(predicate, union.GetAlternates())
}

// JSTypeIterations_AllTypesMatch returns whether all elements of types
// match predicate.
//
// Java: JSTypeIterations.allTypesMatch(Predicate<? super JSType>, ImmutableList<? extends JSType>)
func JSTypeIterations_AllTypesMatch[S JSType](predicate func(JSType) bool, types []S) bool {
	for i := 0; i < len(types); i++ {
		if !predicate(types[i]) {
			return false
		}
	}
	return true
}

// JSTypeIterations_AllTypesMatchInUnion returns whether all elements of
// union match predicate.
//
// Java: JSTypeIterations.allTypesMatch(Predicate<? super JSType>, UnionType)
func JSTypeIterations_AllTypesMatchInUnion(predicate func(JSType) bool, union *UnionType) bool {
	return JSTypeIterations_AllTypesMatch(predicate, union.GetAlternates())
}

// JSTypeIterations_MapTypes returns the result of applying mapper to all
// elements of types.
//
// Java: JSTypeIterations.mapTypes(Function<? super JSType, T>, ImmutableList<? extends JSType>)
func JSTypeIterations_MapTypes[S JSType, T any](mapper func(S) T, types []S) []T {
	builder := make([]T, 0, len(types))
	for i := 0; i < len(types); i++ {
		builder = append(builder, mapper(types[i]))
	}
	return builder
}

// JSTypeIterations_MapTypesInUnion returns the union of applying mapper to
// all elements of union.
//
// Java: JSTypeIterations.mapTypes(Function<? super JSType, ? extends JSType>, UnionType)
func JSTypeIterations_MapTypesInUnion(mapper func(JSType) JSType, union *UnionType) JSType {
	b := UnionType_Builder_(union.GetRegistry())
	for _, alt := range union.GetAlternates() {
		b.AddAlternate(mapper(alt))
	}
	return b.Build()
}

// NominalTypeBuilder is the builder interface for declaring properties on
// class-like (nominal) types. Nominal types consist primarily of three
// separate object types, which may each have their own properties declared
// on them separately: (1) the constructor function, (2) the instance type,
// and (3) the prototype object.
//
// This builder is used during the first part of type checking, while the
// type checker builds up the set of known types. Thus, this interface serves
// as a placeholder to allow operating on (i.e. both assigning properties to,
// and referencing the yet-to-exist type) these not-yet-available types.
type NominalTypeBuilder struct {
	constructor FunctionType
	instance    ObjectType
	prototype   ObjectType
}

// Java: NominalTypeBuilder(FunctionType, ObjectType)
func NewNominalTypeBuilder(constructor FunctionType, instance ObjectType) *NominalTypeBuilder {
	return &NominalTypeBuilder{
		constructor: constructor,
		instance:    instance,
		prototype:   constructor.GetPrototypeProperty(),
	}
}

// DeclarePrototypeProperty declares a property on the nominal type's
// prototype.
//
// Java: NominalTypeBuilder.declarePrototypeProperty(String, JSType, Node)
func (b *NominalTypeBuilder) DeclarePrototypeProperty(name string, type_ JSType, defSite Node) {
	b.prototype.DefineDeclaredProperty(name, type_, defSite)
}

// DeclareInstanceProperty declares an instance property on the nominal
// type.
//
// Java: NominalTypeBuilder.declareInstanceProperty(String, JSType, Node)
func (b *NominalTypeBuilder) DeclareInstanceProperty(name string, type_ JSType, defSite Node) {
	b.instance.DefineDeclaredProperty(name, type_, defSite)
}

// DeclareConstructorProperty declares a static property on the nominal
// type's constructor.
//
// Java: NominalTypeBuilder.declareConstructorProperty(String, JSType, Node)
func (b *NominalTypeBuilder) DeclareConstructorProperty(name string, type_ JSType, defSite Node) {
	b.constructor.DefineDeclaredProperty(name, type_, defSite)
}

// SuperClass returns a NominalTypeBuilder for this type's superclass.
//
// Java: NominalTypeBuilder.superClass()
func (b *NominalTypeBuilder) SuperClass() *NominalTypeBuilder {
	ctor := b.instance.GetSuperClassConstructor()
	if ctor == nil {
		return nil
	}
	return NewNominalTypeBuilder(ctor, ctor.GetInstanceType())
}

// Constructor returns the constructor as a JSType.
//
// Java: NominalTypeBuilder.constructor()
func (b *NominalTypeBuilder) Constructor() FunctionType {
	return b.constructor
}

// Instance returns the instance type as a JSType.
//
// Java: NominalTypeBuilder.instance()
func (b *NominalTypeBuilder) Instance() ObjectType {
	return b.instance
}

// PrototypeOrInstance returns the type of the prototype object (OTI).
//
// TODO(sdh): See if we can just delete this entirely and use instance()
// instead?
//
// Java: NominalTypeBuilder.prototypeOrInstance()
func (b *NominalTypeBuilder) PrototypeOrInstance() ObjectType {
	return b.prototype
}

// TypeStringBuilder visits types to assemble an associated stringification.
//
// Instances are single use. They irreversibly accumulate state required
// during traversal.
type TypeStringBuilder struct {
	builder          strings.Builder
	isForAnnotations bool
	indentation      string
}

// Java: TypeStringBuilder(boolean)
func NewTypeStringBuilder(isForAnnotations bool) *TypeStringBuilder {
	return &TypeStringBuilder{isForAnnotations: isForAnnotations}
}

// Java: TypeStringBuilder.isForAnnotations()
func (b *TypeStringBuilder) IsForAnnotations() bool {
	return b.isForAnnotations
}

// Java: TypeStringBuilder.append(String)
func (b *TypeStringBuilder) Append(x string) *TypeStringBuilder {
	b.builder.WriteString(x)
	return b
}

// Java: TypeStringBuilder.append(JSType)
func (b *TypeStringBuilder) AppendType(x JSType) *TypeStringBuilder {
	x.AppendTo(b)
	return b
}

// AppendTypes appends the JSTypes separated by separator.
//
// Java: TypeStringBuilder.appendAll(Iterable<?>, String)
func (b *TypeStringBuilder) AppendTypes(elements []JSType, separator string) *TypeStringBuilder {
	for i, e := range elements {
		if i > 0 {
			b.Append(separator)
		}
		b.AppendType(e)
	}
	return b
}

// AppendStrings appends the strings separated by separator.
//
// Java: TypeStringBuilder.appendAll(Iterable<?>, String)
func (b *TypeStringBuilder) AppendStrings(elements []string, separator string) *TypeStringBuilder {
	for i, s := range elements {
		if i > 0 {
			b.Append(separator)
		}
		b.Append(s)
	}
	return b
}

// Java: TypeStringBuilder.appendNonNull(JSType)
func (b *TypeStringBuilder) AppendNonNull(type_ JSType) *TypeStringBuilder {
	if b.isForAnnotations &&
		type_.IsObject() &&
		!type_.IsUnknownType() &&
		!type_.IsTemplateType() &&
		!type_.IsRecordType() &&
		!type_.IsFunctionType() &&
		!type_.IsUnionType() &&
		!type_.IsLiteralObject() {
		b.Append("!")
	}
	return b.AppendType(type_)
}

// Java: TypeStringBuilder.breakLineAndIndent()
func (b *TypeStringBuilder) BreakLineAndIndent() *TypeStringBuilder {
	b.builder.WriteString("\n")
	b.builder.WriteString(b.indentation)
	return b
}

// Java: TypeStringBuilder.indent(Runnable)
func (b *TypeStringBuilder) Indent(cb func()) *TypeStringBuilder {
	lastIndent := b.indentation
	b.indentation = lastIndent + "  "
	cb()
	b.indentation = lastIndent
	return b
}

// Java: TypeStringBuilder.build()
func (b *TypeStringBuilder) Build() string {
	return b.builder.String()
}

// Java: TypeStringBuilder.cloneWithConfig()
func (b *TypeStringBuilder) CloneWithConfig() *TypeStringBuilder {
	clone := NewTypeStringBuilder(b.isForAnnotations)
	clone.indentation = b.indentation
	return clone
}
