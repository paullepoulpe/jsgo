// Port of com/google/javascript/rhino/jstype/EquivalenceMethod.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/JSTypeClass.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/BooleanLiteralSet.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/JSTypeIterations.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/NominalTypeBuilder.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/TypeStringBuilder.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/ContainsUpperBoundSuperTypeVisitor.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/CanCastToVisitor.java @ closure-compiler v20261021

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
 *   John Lenz
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
	"strconv"
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
	panic(jutil.NewIllegalArgumentExceptionMsg("Ordinal: " + strconv.Itoa(ordinal)))
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
	panic(jutil.NewIndexOutOfBoundsException("Ordinal: " + strconv.Itoa(s.Ordinal())))
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
	for _, t := range types {
		if predicate(t) {
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
	for _, t := range types {
		if !predicate(t) {
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
	for _, t := range types {
		builder = append(builder, mapper(t))
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
	indentDepth      int
	isForAnnotations bool
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
func (b *TypeStringBuilder) Append(x string) {
	b.builder.WriteString(x)
}

// Java: TypeStringBuilder.append(JSType)
func (b *TypeStringBuilder) AppendType(x JSType) {
	switch t := x.(type) {
	case *BooleanType:
		b.builder.WriteString("boolean")
	case *NumberType:
		b.builder.WriteString("number")
	case *StringType:
		b.builder.WriteString("string")
	case *BigIntType:
		b.builder.WriteString("bigint")
	case *NullType:
		b.builder.WriteString("null")
	case *VoidType:
		b.builder.WriteString("undefined")
	case *AllType:
		b.builder.WriteString("*")
	case *UnknownType:
		if t.isChecked {
			b.builder.WriteString("??")
		} else {
			b.builder.WriteString("?")
		}
	case *NoType:
		t.AppendTo(b)
	case *NoObjectTypeBase:
		t.AppendTo(b)
	case *NoResolvedType:
		t.AppendTo(b)
	case *KnownSymbolType:
		b.builder.WriteString(t.name)
	case *SymbolTypeBase:
		t.AppendTo(b)
	case *InstanceObjectType:
		t.AppendTo(b)
	case *RecordType:
		t.AppendTo(b)
	case *PrototypeObjectTypeBase:
		t.AppendTo(b)
	case *FunctionTypeBase:
		t.AppendTo(b)
	case *UnionType:
		t.AppendTo(b)
	case *TemplatizedType:
		t.AppendTo(b)
	case *TemplateType:
		t.AppendTo(b)
	case *NamedType:
		t.AppendTo(b)
	case *EnumType:
		t.AppendTo(b)
	case *EnumElementType:
		t.AppendTo(b)
	case *ProxyObjectTypeBase:
		t.AppendTo(b)
	default:
		tmp := &TypeStringBuilder{
			isForAnnotations: b.isForAnnotations,
			indentDepth:      b.indentDepth,
		}
		x.AppendTo(tmp)
		b.builder.WriteString(tmp.Build())
	}
}

// AppendTypes appends the JSTypes separated by separator.
//
// Java: TypeStringBuilder.appendAll(Iterable<?>, String)
func (b *TypeStringBuilder) AppendTypes(elements []JSType, separator string) {
	for i, e := range elements {
		if i > 0 {
			b.builder.WriteString(separator)
		}
		b.AppendType(e)
	}
}

// AppendStrings appends the strings separated by separator.
//
// Java: TypeStringBuilder.appendAll(Iterable<?>, String)
func (b *TypeStringBuilder) AppendStrings(elements []string, separator string) {
	for i, s := range elements {
		if i > 0 {
			b.builder.WriteString(separator)
		}
		b.builder.WriteString(s)
	}
}

// Java: TypeStringBuilder.appendNonNull(JSType)
func (b *TypeStringBuilder) AppendNonNull(type_ JSType) {
	if b.isForAnnotations &&
		type_.IsObject() &&
		!type_.IsUnknownType() &&
		!type_.IsTemplateType() &&
		!type_.IsRecordType() &&
		!type_.IsFunctionType() &&
		!type_.IsUnionType() &&
		!type_.IsLiteralObject() {
		b.builder.WriteString("!")
	}
	b.AppendType(type_)
}

const typeStringIndentSpaces = "                                "

// Java: TypeStringBuilder.breakLineAndIndent()
func (b *TypeStringBuilder) BreakLineAndIndent() {
	b.builder.WriteByte('\n')
	spaces := b.indentDepth * 2
	for spaces > len(typeStringIndentSpaces) {
		b.builder.WriteString(typeStringIndentSpaces)
		spaces -= len(typeStringIndentSpaces)
	}
	if spaces > 0 {
		b.builder.WriteString(typeStringIndentSpaces[:spaces])
	}
}

// Java: TypeStringBuilder.indent(Runnable)
func (b *TypeStringBuilder) Indent(cb func()) {
	b.indentDepth++
	cb()
	b.indentDepth--
}

// Java: TypeStringBuilder.build()
func (b *TypeStringBuilder) Build() string {
	return b.builder.String()
}

// Java: TypeStringBuilder.cloneWithConfig()
func (b *TypeStringBuilder) CloneWithConfig() *TypeStringBuilder {
	return &TypeStringBuilder{
		isForAnnotations: b.isForAnnotations,
		indentDepth:      b.indentDepth,
	}
}

// containsUpperBoundSuperTypePresent reports whether target appears in the
// reference chain of forwarding types reachable from type_.
//
// Java: new ContainsUpperBoundSuperTypeVisitor(target) == Result.PRESENT
func containsUpperBoundSuperTypePresent(type_ JSType, target JSType) bool {
	return containsUpperBoundSuperType(type_, target) == ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
}

// containsUpperBoundSuperTypeCycle reports whether the reference chain of
// forwarding types reachable from type_ contains a cycle.
//
// Java: new ContainsUpperBoundSuperTypeVisitor(null) == Result.CYCLE
func containsUpperBoundSuperTypeCycle(type_ JSType) bool {
	return containsUpperBoundSuperType(type_, nil) == ContainsUpperBoundSuperTypeVisitor_Result_CYCLE
}

func containsUpperBoundSuperType(type_ JSType, target JSType) ContainsUpperBoundSuperTypeVisitor_Result {
	var seenBuf [8]JSType
	seen := seenBuf[:0]
	return containsUpperBoundSuperTypeRec(type_, target, &seen)
}

func containsUpperBoundSuperTypeRec(type_ JSType, target JSType, seen *[]JSType) ContainsUpperBoundSuperTypeVisitor_Result {
	switch x := type_.(type) {
	case *TemplateType:
		return containsUpperBoundForwarding(x, x.GetBound(), target, seen)
	case *NamedType:
		return containsUpperBoundForwarding(x, x.GetReferencedType(), target, seen)
	case *TemplatizedType:
		return containsUpperBoundForwarding(x, x.GetReferencedType(), target, seen)
	case *UnionType:
		if JSType(x) == target {
			return ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
		}
		for _, alt := range x.GetAlternates() {
			foundInAlt := containsUpperBoundSuperTypeRec(alt, target, seen)
			if foundInAlt != ContainsUpperBoundSuperTypeVisitor_Result_ABSENT {
				return foundInAlt
			}
		}
		return ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
	case *AllType, *BooleanType, *NoObjectTypeBase, *UnknownType, *NullType,
		*NumberType, *BigIntType, *StringType, *VoidType:
		// Visitor.WithDefaultCase passes nil for zero-arg cases, so CaseDefault
		// always returns ABSENT for these types.
		return ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
	case *NoType, *EnumElementType, *FunctionTypeBase,
		*InstanceObjectType, *PrototypeObjectTypeBase, *RecordType, *EnumType, *NoResolvedType,
		*ProxyObjectTypeBase, *SymbolTypeBase, *KnownSymbolType:
		if target != nil && type_ == target {
			return ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
		}
		return ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
	default:
		if _, ok := type_.(NoObjectType); ok {
			return ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
		}
		switch type_.(type) {
		case FunctionType, ProxyObjectType, ObjectType, SymbolType:
			if target != nil && type_ == target {
				return ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
			}
			return ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
		}
		panic(jutil.NewUnsupportedOperationException())
	}
}

// Java: ContainsUpperBoundSuperTypeVisitor.caseForwardingType(JSType, JSType)
func containsUpperBoundForwarding(type_ JSType, reference JSType, target JSType, seen *[]JSType) ContainsUpperBoundSuperTypeVisitor_Result {
	if type_ == target {
		return ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
	} else if slices.Contains(*seen, type_) {
		return ContainsUpperBoundSuperTypeVisitor_Result_CYCLE
	} else {
		*seen = append(*seen, type_)
		return containsUpperBoundSuperTypeRec(reference, target, seen)
	}
}

// ContainsUpperBoundSuperTypeVisitor_Result represents the outcome of a
// visitation of the ContainsUpperBoundSuperTypeVisitor.
type ContainsUpperBoundSuperTypeVisitor_Result int32

const (
	ContainsUpperBoundSuperTypeVisitor_Result_PRESENT ContainsUpperBoundSuperTypeVisitor_Result = iota
	ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
	// Containment is undecidable due to a reference cycle.
	ContainsUpperBoundSuperTypeVisitor_Result_CYCLE
)

var containsUpperBoundSuperTypeVisitor_Result_names = [...]string{"PRESENT", "ABSENT", "CYCLE"}

// String is Enum.name().
func (r ContainsUpperBoundSuperTypeVisitor_Result) String() string {
	return containsUpperBoundSuperTypeVisitor_Result_names[r]
}

// canCastTo reports whether values of thisType can be cast to thatType.
//
// Java: CanCastToVisitor + JSType.visit(RelationshipVisitor, JSType)
func canCastTo(thisType JSType, thatType JSType) bool {
	switch x := thisType.(type) {
	case *UnknownType, *NoType, *NoObjectTypeBase, *AllType, *TemplateType:
		return true
	case *BigIntType, *BooleanType, *NullType, *NumberType, *StringType,
		*SymbolTypeBase, *KnownSymbolType, *VoidType,
		*FunctionTypeBase,
		*InstanceObjectType, *PrototypeObjectTypeBase, *RecordType, *EnumType, *NoResolvedType:
		return canCastToHelper(thisType, thatType)
	case *UnionType:
		return canCastUnionFrom(x, thatType)
	case *TemplatizedType:
		// TODO(johnlenz): once the templated type work is finished, restrict the
		// type parameters.
		return canCastTo(x.GetReferencedTypeInternal(), thatType)
	case *EnumElementType:
		return canCastTo(x.GetPrimitiveType(), thatType)
	case *NamedType:
		return canCastTo(x.referencedType, thatType)
	case *ProxyObjectTypeBase:
		return canCastTo(x.referencedType, thatType)
	default:
		if _, ok := thisType.(NoObjectType); ok {
			return true // TODO(johnlenz): restrict to objects
		}
		if fn, ok := thisType.(FunctionType); ok {
			return canCastToHelper(fn, thatType)
		}
		if proxy, ok := thisType.(ProxyObjectType); ok {
			return canCastTo(proxy.GetReferencedTypeInternal(), thatType)
		}
		if obj, ok := thisType.(ObjectType); ok {
			return canCastToHelper(obj, thatType)
		}
		if val, ok := thisType.(ValueType); ok {
			return canCastToHelper(val, thatType)
		}
		panic(jutil.NewUnsupportedOperationException())
	}
}

// Java: CanCastToVisitor.canCastToUnion(JSType, UnionType)
func canCastToUnion(thisType JSType, unionType *UnionType) bool {
	for _, type_ := range unionType.GetAlternates() {
		if canCastTo(thisType, type_) {
			return true
		}
	}
	return false
}

// Java: CanCastToVisitor.canCastToFunction(JSType, FunctionType)
func canCastToFunction(thisType JSType, functionType FunctionType) bool {
	if thisType.IsFunctionType() {
		// TODO(johnlenz): visit function parts
		return true
	}
	return thisType.IsSubtypeOf(functionType) ||
		functionType.IsSubtype(thisType)
}

// Java: CanCastToVisitor.isInterface(JSType)
func canCastIsInterface(type_ JSType) bool {
	objType := type_.ToObjectType()
	if objType != nil {
		constructor := objType.GetConstructor()
		return constructor != nil && constructor.IsInterface()
	}
	return false
}

// Java: CanCastToVisitor.castCastToHelper(JSType, JSType)
func canCastToHelper(thisType JSType, thatType JSType) bool {
	if thisType.IsNoResolvedType() {
		return true
	} else if thatType.IsUnknownType() ||
		thatType.IsAllType() ||
		thatType.IsNoObjectType() || // TODO(johnlenz): restrict to objects
		thatType.IsNoResolvedType() ||
		thatType.IsNoType() {
		return true
	} else if thisType.IsRecordType() || thatType.IsRecordType() {
		return true // TODO(johnlenz): are there any misuses we can catch?
	} else if canCastIsInterface(thisType) || canCastIsInterface(thatType) {
		return true // TODO(johnlenz): are there any misuses we can catch?
	} else if thatType.IsEnumElementType() {
		return canCastTo(thisType, thatType.ToMaybeEnumElementType().GetPrimitiveType())
	} else if thatType.IsUnionType() {
		return canCastToUnion(thisType, thatType.ToMaybeUnionType())
	} else if thatType.IsFunctionType() {
		return canCastToFunction(thisType, thatType.ToMaybeFunctionType())
	} else if thatType.IsTemplatizedType() {
		// TODO(johnlenz): once the templated type work is finished, restrict
		// the type parameters.
		return canCastTo(thisType, thatType.ToMaybeTemplatizedType().GetReferencedTypeInternal())
	}

	return thisType.IsSubtypeOf(thatType) || thatType.IsSubtypeOf(thisType)
}

// Java: CanCastToVisitor.caseUnionType(UnionType, JSType)
func canCastUnionFrom(thisType *UnionType, thatType JSType) bool {
	visited := false
	for _, type_ := range thisType.GetAlternates() {
		if type_.IsVoidType() || type_.IsNullType() {
			// Don't allow if the only match between the types is null or void,
			// otherwise any nullable type would be castable to any other nullable
			// type and we don't want that.
		} else {
			visited = true
			if canCastTo(type_, thatType) {
				return true
			}
		}
	}

	// Special case the "null|undefined" union and allow it to be cast to any
	// type allowing either null|undefined.
	if !visited {
		nullType := thisType.GetNativeType(JSTypeNative_NULL_TYPE)
		voidType := thisType.GetNativeType(JSTypeNative_VOID_TYPE)
		return canCastTo(nullType, thatType) || canCastTo(voidType, thatType)
	}

	return false
}
