// Port of com/google/javascript/rhino/jstype/JSTypeRegistryTest.java @ closure-compiler v20261021

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
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino_test

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	. "github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// TODO(user): extend this class with more tests, as JSTypeRegistry is
// now much larger
type jsTypeRegistryTest struct {
	rhinotesting.Truth
	registry *JSTypeRegistry
}

func newJsTypeRegistryTest(t *testing.T) *jsTypeRegistryTest {
	s := &jsTypeRegistryTest{Truth: rhinotesting.NewTruth(t)}
	s.registry = NewJSTypeRegistryWithForwardDeclaredTypes(nil, nil)
	s.SetUp()
	return s
}

func (s *jsTypeRegistryTest) SetUp() {
	unused := s.registry.GetResolver().OpenForDefinition()
	_ = unused
}

func TestJSTypeRegistry_GetBuiltInType_boolean(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	s.AssertType(s.registry.GetType(nil, "boolean")).
		IsEqualTo(s.registry.GetNativeType(JSTypeNative_BOOLEAN_TYPE))
}

func TestJSTypeRegistry_GetBuiltInType_bigint(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	s.AssertType(s.registry.GetType(nil, "bigint")).
		IsEqualTo(s.registry.GetNativeType(JSTypeNative_BIGINT_TYPE))
}

func TestJSTypeRegistry_GetBuiltInType_iterable(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	s.AssertType(s.registry.GetGlobalType("Iterable")).IsEqualTo(s.registry.GetNativeType(JSTypeNative_ITERABLE_TYPE))
}

func TestJSTypeRegistry_GetBuiltInType_iterator(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	s.AssertType(s.registry.GetGlobalType("Iterator")).IsEqualTo(s.registry.GetNativeType(JSTypeNative_ITERATOR_TYPE))
}

func TestJSTypeRegistry_GetBuiltInType_iteratorLike(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	s.AssertType(s.registry.GetGlobalType("IteratorLike")).
		IsEqualTo(s.registry.GetNativeType(JSTypeNative_ITERATOR_LIKE_TYPE))
}

func TestJSTypeRegistry_IteratorExtendsIteratorLike(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	s.AssertThat(
		s.registry.
			GetNativeType(JSTypeNative_ITERATOR_TYPE).
			IsSubtypeOf(s.registry.GetNativeType(JSTypeNative_ITERATOR_LIKE_TYPE))).
		IsTrue()
}

func TestJSTypeRegistry_GetBuiltInType_generator(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	s.AssertType(s.registry.GetGlobalType("Generator")).
		IsEqualTo(s.registry.GetNativeType(JSTypeNative_GENERATOR_TYPE))
}

func TestJSTypeRegistry_GetBuiltInType_async_iterable(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	s.AssertType(s.registry.GetGlobalType("AsyncIterable")).
		IsEqualTo(s.registry.GetNativeType(JSTypeNative_ASYNC_ITERABLE_TYPE))
}

func TestJSTypeRegistry_GetBuiltInType_async_iterator(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	s.AssertType(s.registry.GetGlobalType("AsyncIterator")).
		IsEqualTo(s.registry.GetNativeType(JSTypeNative_ASYNC_ITERATOR_TYPE))
}

func TestJSTypeRegistry_GetBuiltInType_async_generator(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	s.AssertType(s.registry.GetGlobalType("AsyncGenerator")).
		IsEqualTo(s.registry.GetNativeType(JSTypeNative_ASYNC_GENERATOR_TYPE))
}

func TestJSTypeRegistry_GetBuiltInType_gbigint(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	s.AssertType(s.registry.GetGlobalType("gbigint")).
		IsEqualTo(s.registry.GetNativeType(JSTypeNative_GBIGINT_TYPE))
}

func TestJSTypeRegistry_GetBuiltInType_gbigint_unification(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	// This simulates what happens when the compiler parses: /** @interface */ function gbigint() {}
	interfaceType :=
		FunctionType_Builder_(s.registry).ForInterface().WithName(jutil.StringPtr("gbigint")).Build()
	s.registry.DeclareType(nil, "gbigint", interfaceType.GetInstanceType())

	s.AssertType(s.registry.GetGlobalType("gbigint")).
		IsEqualTo(s.registry.GetNativeType(JSTypeNative_GBIGINT_TYPE))
}

func TestJSTypeRegistry_GetBuildInType_iTemplateArray(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	s.AssertType(s.registry.GetGlobalType("ITemplateArray")).
		IsEqualTo(s.registry.GetNativeType(JSTypeNative_I_TEMPLATE_ARRAY_TYPE))
}

func TestJSTypeRegistry_GetBuiltInType_Promise(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	promiseType := s.registry.GetNativeObjectType(JSTypeNative_PROMISE_TYPE)
	s.AssertType(s.registry.GetGlobalType("Promise")).IsEqualTo(promiseType)

	// Test that it takes one parameter of type
	// function(function((IThenable<TYPE>|TYPE|null|{then: ?})=): ?, function(*=): ?): ?
	promiseCtor := promiseType.GetConstructor()
	paramList := promiseCtor.GetParameters()
	s.AssertThat(paramList).HasSize(1)
	firstParameter := paramList[0]
	paramType := firstParameter.GetJSType().ToMaybeFunctionType()
	s.AssertThat(paramType.String()).
		IsEqualTo(
			"function(function((IThenable<TYPE>|TYPE|null|{then: ?})=): ?, function(*=): ?): ?")
}

func TestJSTypeRegistry_GetDeclaredType(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	type_ := s.registry.CreateAnonymousObjectType(nil)
	name := "Foo"
	s.registry.DeclareType(nil, name, type_)
	s.AssertType(s.registry.GetType(nil, name)).IsEqualTo(type_)

	// Ensure different instances are independent.
	typeRegistry2 := NewJSTypeRegistry(nil)
	s.AssertThat(typeRegistry2.GetType(nil, name)).IsNull()
	s.AssertType(s.registry.GetType(nil, name)).IsEqualTo(type_)
}

func TestJSTypeRegistry_ReadableTypeName(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	registry := s.registry

	s.AssertThat(s.getReadableTypeNameHelperNative(registry, JSTypeNative_ALL_TYPE)).IsEqualTo("*")

	s.AssertThat(s.getReadableTypeNameHelperNative(registry, JSTypeNative_BOOLEAN_TYPE)).IsEqualTo("boolean")
	s.AssertThat(s.getReadableTypeNameHelperNative(registry, JSTypeNative_BOOLEAN_OBJECT_TYPE)).IsEqualTo("Boolean")
	s.AssertThat(s.getReadableTypeNameHelperNative(registry, JSTypeNative_BOOLEAN_OBJECT_FUNCTION_TYPE)).
		IsEqualTo("function")

	s.AssertThat(s.getReadableTypeNameHelperNative(registry, JSTypeNative_NULL_VOID)).IsEqualTo("(null|undefined)")
	s.AssertThat(s.getReadableTypeNameHelperNativeDeref(registry, JSTypeNative_NULL_VOID, true)).IsEqualTo("(null|undefined)")

	s.AssertThat(
		s.getReadableTypeNameHelper(
			registry, s.union(registry, JSTypeNative_NUMBER_TYPE, JSTypeNative_STRING_TYPE, JSTypeNative_NULL_TYPE))).
		IsEqualTo("(number|string|null)")

	s.AssertThat(
		s.getReadableTypeNameHelperDeref(
			registry, s.union(registry, JSTypeNative_NUMBER_TYPE, JSTypeNative_STRING_TYPE, JSTypeNative_NULL_TYPE), true)).
		IsEqualTo("(Number|String)")
}

func TestJSTypeRegistry_CreateTypeFromCommentNode_createsNamedTypeIfNameIsUndefined(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	// Create an empty global scope.
	globalScope := rhinotesting.NewMapBasedScope(mapOf())

	type_ :=
		s.registry.CreateTypeFromCommentNodeWithScope(
			NewNodeWithChild(Token_BANG, IR_String("a.b.c")), "srcfile.js", globalScope)

	// The type registry assumes that 'a.b.c' will be defined later.
	rhinotesting.IsInstanceOf[*NamedType](s.AssertThat(type_))
}

func TestJSTypeRegistry_CreateTypeFromCommentNode_multipleLookupsInSameScope(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	// Create an empty global scope.
	globalScope := rhinotesting.NewMapBasedScope(mapOf())

	typeOne :=
		s.registry.CreateTypeFromCommentNodeWithScope(
			NewNodeWithChild(Token_BANG, IR_String("Foo")), "srcfile.js", globalScope)
	typeTwo :=
		s.registry.CreateTypeFromCommentNodeWithScope(
			NewNodeWithChild(Token_BANG, IR_String("Foo")), "srcfile.js", globalScope)

	// Different lookups of Foo create equivalent instances of NamedType.
	rhinotesting.IsInstanceOf[*NamedType](s.AssertThat(typeOne))
	rhinotesting.IsInstanceOf[*NamedType](s.AssertThat(typeTwo))
	s.AssertType(typeOne).IsEqualTo(typeTwo)
}

func TestJSTypeRegistry_CreateTypeFromCommentNode_usesGlobalTypeIfExists(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	// Create a global scope and a global type 'Foo'.
	globalScope :=
		s.createStaticTypedScope(IR_Root(), nil, jutil.NewLinkedHashMap[string, StaticTypedSlot](), jutil.NewLinkedHashSet[string]())
	s.registry.DeclareType(globalScope, "Foo", s.registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE))

	// Create an empty child scope.
	emptyLocalScope :=
		s.createStaticTypedScope(IR_Block(), globalScope, jutil.NewLinkedHashMap[string, StaticTypedSlot](), jutil.NewLinkedHashSet[string]())

	// Looking up 'Foo' in the child scope resolves to the global type `Foo`, not a NamedType.
	type_ :=
		s.registry.CreateTypeFromCommentNodeWithScope(
			NewNodeWithChild(Token_BANG, IR_String("Foo")), "srcfile.js", emptyLocalScope)
	rhinotesting.IsNotInstanceOf[*NamedType](s.AssertThat(type_))
}

// localSlot is the anonymous SimpleSlot subclass of
// testCreateTypeFromCommentNode_createsNamedTypeIfLocalShadowsGlobalType.
type localSlot struct {
	*SimpleSlotBase
	scope StaticTypedScope
}

func (l *localSlot) GetScope() StaticScope { return l.scope }

func TestJSTypeRegistry_CreateTypeFromCommentNode_createsNamedTypeIfLocalShadowsGlobalType(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	// Create a global scope and a global type 'Foo'.
	unknownType := s.registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	globalScope :=
		s.createStaticTypedScope(IR_Root(), nil, jutil.NewLinkedHashMap[string, StaticTypedSlot](), jutil.NewLinkedHashSet[string]())
	s.registry.DeclareType(globalScope, "Foo", unknownType)

	// Create a child scope containing a slot for 'Foo'.
	slots := jutil.NewLinkedHashMap[string, StaticTypedSlot]()
	localScopeWithFoo :=
		s.createStaticTypedScope(IR_Block(), globalScope, slots, jutil.NewLinkedHashSet[string]())
	slots.Put(
		"Foo",
		&localSlot{SimpleSlotBase: NewSimpleSlot("Foo", unknownType /* inferred= */, true), scope: localScopeWithFoo})

	// Looking up 'Foo' in the local scope creates a NamedType; even though `Foo` is not in the
	// JSTypeRegistry in the local scope, the registry assumes it will be defined later.
	type_ :=
		s.registry.CreateTypeFromCommentNodeWithScope(
			NewNodeWithChild(Token_BANG, IR_String("Foo")), "srcfile.js", localScopeWithFoo)
	rhinotesting.IsInstanceOf[*NamedType](s.AssertThat(type_))
}

func TestJSTypeRegistry_NativeTypesAreUnique(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	for _, n1 := range JSTypeNative_values() {
		for _, n2 := range JSTypeNative_values() {
			t1 := s.registry.GetNativeType(n1)
			t2 := s.registry.GetNativeType(n2)
			if !t1.Equals(t2) {
				continue
			}
			s.AssertThat(t1).IsSameInstanceAs(t2)
			s.AssertThat(n1).IsEqualTo(n2)
		}
	}
}

func TestJSTypeRegistry_CreateTypeFromCommentNode_usesTopMostScopeOfName(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	// Create a global scope and a global type 'Foo'.
	globalScope :=
		s.createStaticTypedScope(IR_Root(), nil, jutil.NewLinkedHashMap[string, StaticTypedSlot](), jutil.NewLinkedHashSet[string]())
	s.registry.DeclareType(globalScope, "Foo", s.registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE))

	// Create a child scope with no actual variables, but that returns the child scope when asked
	// for the eventual scope of the name "Foo".
	emptyLocalScope :=
		s.createStaticTypedScope(IR_Block(), globalScope, jutil.NewLinkedHashMap[string, StaticTypedSlot](), jutil.NewLinkedHashSet("Foo"))

	// Looking up 'Foo' in the child scope resolves to a NamedType.
	type_ :=
		s.registry.CreateTypeFromCommentNodeWithScope(
			NewNodeWithChild(Token_BANG, IR_String("Foo")), "srcfile.js", emptyLocalScope)
	rhinotesting.IsInstanceOf[*NamedType](s.AssertThat(type_))
}

func TestJSTypeRegistry_GetBuiltInType_ReadonlyMap(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	readonlyMapType := s.registry.GetNativeObjectType(JSTypeNative_READONLY_MAP_TYPE)
	s.AssertType(s.registry.GetGlobalType("ReadonlyMap")).IsEqualTo(readonlyMapType)
}

func TestJSTypeRegistry_GetBuiltInType_Map(t *testing.T) {
	s := newJsTypeRegistryTest(t)
	// Ensure we have a Map.
	mapType := s.registry.GetNativeObjectType(JSTypeNative_MAP_TYPE)
	s.AssertType(s.registry.GetGlobalType("Map")).IsEqualTo(mapType)

	// Ensure it implements ReadonlyMap.
	readonlyMapType := s.registry.GetNativeObjectType(JSTypeNative_READONLY_MAP_TYPE)
	s.AssertThat(mapType.GetCtorImplementedInterfaces()).HasSize(1)
	s.AssertThat(append([]ObjectType{}, mapType.GetCtorImplementedInterfaces()...)[0].GetRawType()).
		IsEqualTo(readonlyMapType)

	// Ensure it takes one optional parameter of the appropriate type.
	paramList := mapType.GetConstructor().GetParameters()
	s.AssertThat(paramList).HasSize(1)
	s.AssertThat(paramList[0].IsOptional()).IsTrue()
	s.AssertThat(paramList[0].GetJSType().String()).
		IsEqualTo("(Array<Array<(K|V)>>|Iterable<Array<(K|V)>,?,?>|null|undefined)")
}

// jsTypeRegistryTestScope is the anonymous AbstractStaticScope of
// createStaticTypedScope.
type jsTypeRegistryTestScope struct {
	rhinotesting.AbstractStaticScopeBase
	root          Node
	parentScope   StaticTypedScope
	slots         *jutil.LinkedHashMap[string, StaticTypedSlot]
	reservedNames *jutil.LinkedHashSet[string]
}

func (sc *jsTypeRegistryTestScope) GetRootNode() Node {
	return sc.root
}

func (sc *jsTypeRegistryTestScope) GetParentScope() StaticScope {
	if sc.parentScope == nil {
		return nil
	}
	return sc.parentScope
}

func (sc *jsTypeRegistryTestScope) GetSlot(name string) StaticSlot {
	if slot := sc.slots.Get(name); slot != nil {
		return slot
	}
	return nil
}

func (sc *jsTypeRegistryTestScope) GetTopmostScopeOfEventualDeclaration(name string) StaticScope {
	if sc.slots.ContainsKey(name) || sc.reservedNames.Contains(name) {
		return sc
	}
	if sc.parentScope != nil {
		return sc.parentScope.GetTopmostScopeOfEventualDeclaration(name)
	}
	return nil
}

// Returns a scope that overrides a few methods from {@link AbstractStaticScope}
func (s *jsTypeRegistryTest) createStaticTypedScope(
	root Node,
	parentScope StaticTypedScope,
	slots *jutil.LinkedHashMap[string, StaticTypedSlot],
	reservedNames *jutil.LinkedHashSet[string]) StaticTypedScope {
	sc := &jsTypeRegistryTestScope{root: root, parentScope: parentScope, slots: slots, reservedNames: reservedNames}
	sc.InitAbstractStaticScopeBase(sc)
	return sc
}

func (s *jsTypeRegistryTest) union(registry *JSTypeRegistry, types ...JSTypeNative) JSType {
	return registry.CreateUnionTypeOfNatives(types...)
}

func (s *jsTypeRegistryTest) getReadableTypeNameHelperNative(registry *JSTypeRegistry, type_ JSTypeNative) string {
	return s.getReadableTypeNameHelperDeref(registry, registry.GetNativeType(type_), false)
}

func (s *jsTypeRegistryTest) getReadableTypeNameHelperNativeDeref(
	registry *JSTypeRegistry, type_ JSTypeNative, deref bool) string {
	return s.getReadableTypeNameHelperDeref(registry, registry.GetNativeType(type_), deref)
}

func (s *jsTypeRegistryTest) getReadableTypeNameHelper(registry *JSTypeRegistry, type_ JSType) string {
	return s.getReadableTypeNameHelperDeref(registry, type_, false)
}

func (s *jsTypeRegistryTest) getReadableTypeNameHelperDeref(registry *JSTypeRegistry, type_ JSType, deref bool) string {
	n := NewNode(Token_ADD)
	n.SetJSType(type_)
	return registry.GetReadableJSTypeName(n, deref)
}
