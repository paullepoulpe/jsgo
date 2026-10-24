// Port of com/google/javascript/rhino/jstype/PropertyMapTest.java @ closure-compiler v20261021

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

package rhino_test

import (
	"testing"

	. "github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

type propertyMapTest struct {
	rhinotesting.Truth
	symbolKey  Property_Key
	registry   *JSTypeRegistry
	numberType JSType
}

func newPropertyMapTest(t *testing.T) *propertyMapTest {
	s := &propertyMapTest{Truth: rhinotesting.NewTruth(t)}
	s.SetUp()
	return s
}

func (s *propertyMapTest) SetUp() {
	s.registry = NewJSTypeRegistry(rhinotesting.NewTestErrorReporter())
	s.symbolKey = NewProperty_SymbolKey(NewKnownSymbolType(s.registry, "Symbol.foo"))
	s.numberType = s.registry.GetNativeType(JSTypeNative_NUMBER_TYPE)
}

func TestPropertyMap_EmptyMap_hasNoProperties(t *testing.T) {
	s := newPropertyMapTest(t)
	empty := PropertyMap_ImmutableEmptyMap()

	s.AssertThat(empty.GetPropertiesCount()).IsEqualTo(0)
	s.AssertThat(empty.GetOwnPropertyNames()).IsEmpty()
	s.AssertThat(empty.GetOwnKnownSymbols()).IsEmpty()
	s.AssertThat(empty.GetAllKeys()).
		IsEqualTo(NewPropertyMap_AllKeys([]string{}, []*KnownSymbolType{}))
	s.AssertThat(empty.GetPrimaryParent()).IsNull()
	s.AssertThat(empty.GetOwnProperty("x")).IsNull()
	s.AssertThat(empty.GetOwnPropertyByKey(s.symbolKey)).IsNull()
}

func TestPropertyMap_GetAllKeys_returnsBothSymbolAndStringKeys(t *testing.T) {
	s := newPropertyMapTest(t)
	properties := NewPropertyMap()

	properties.PutProperty("x", NewProperty("x", s.numberType /* inferred= */, false, nil))
	properties.PutPropertyByKey(
		s.symbolKey, NewPropertyByKey(s.symbolKey, s.numberType /* inferred= */, false, nil))

	s.AssertThat(properties.GetAllKeys().StringKeys()).ContainsExactly("x")
	s.AssertThat(properties.GetAllKeys().KnownSymbolKeys()).ContainsExactly(s.symbolKey.Symbol())
}

func TestPropertyMap_GetAllKeys_cachesStringKeyResult(t *testing.T) {
	s := newPropertyMapTest(t)
	properties := NewPropertyMap()

	properties.PutProperty("x", NewProperty("x", s.numberType /* inferred= */, false, nil))

	result1 := properties.GetAllKeys().StringKeys()
	result2 := properties.GetAllKeys().StringKeys()
	s.AssertThat(result1).IsSameInstanceAs(result2)
}

func TestPropertyMap_GetAllKeys_cachesSymbolKeyResult(t *testing.T) {
	s := newPropertyMapTest(t)
	properties := NewPropertyMap()

	properties.PutPropertyByKey(
		s.symbolKey, NewPropertyByKey(s.symbolKey, s.numberType /* inferred= */, false, nil))

	result1 := properties.GetAllKeys().KnownSymbolKeys()
	result2 := properties.GetAllKeys().KnownSymbolKeys()
	s.AssertThat(result1).IsSameInstanceAs(result2)
}

func TestPropertyMap_GetAllKeys_beforeAndAfterMutation_stringKey_invalidatesCache(t *testing.T) {
	s := newPropertyMapTest(t)
	// PropertyMap caches the result of getAllKeys(), as it can be expensive to compute for
	// maps with many ancestors. Verify that recalculating getAllKeys() after a mutation returns the
	// correct result.
	properties := NewPropertyMap()

	// Request getAllKeys() once, which will cache the empty results.
	s.AssertThat(properties.GetAllKeys().StringKeys()).IsEmpty()

	// Invalidate the cache.
	properties.PutProperty("x", NewProperty("x", s.numberType /* inferred= */, false, nil))

	// getAllKeys() should recalculate the result.
	s.AssertThat(properties.GetAllKeys().StringKeys()).ContainsExactly("x")
}

func TestPropertyMap_GetAllKeys_beforeAndAfterMutation_symbolKey_invalidatesCache(t *testing.T) {
	s := newPropertyMapTest(t)
	// PropertyMap caches the result of getAllKeys(), as it can be expensive to compute for
	// maps with many ancestors. Verify that recalculating getAllKeys() after a mutation returns the
	// correct result.
	properties := NewPropertyMap()

	// Request getAllKeys() once, which will cache the empty results.
	s.AssertThat(properties.GetAllKeys().KnownSymbolKeys()).IsEmpty()

	// Invalidate the cache.
	properties.PutPropertyByKey(
		s.symbolKey, NewPropertyByKey(s.symbolKey, s.numberType /* inferred= */, false, nil))

	// getAllKeys() should recalculate the result.
	s.AssertThat(properties.GetAllKeys().KnownSymbolKeys()).ContainsExactly(s.symbolKey.Symbol())
}

func TestPropertyMap_LookUpKeysFromParent_stringKey(t *testing.T) {
	s := newPropertyMapTest(t)
	parentType :=
		s.registry.CreateRecordType(mapOf("x", s.numberType)).AssertObjectType()
	childType := s.registry.CreateObjectType(ptr("Child"), parentType)
	child := NewPropertyMap()
	child.SetParentSource(childType)

	s.AssertThat(child.GetOwnProperty("x")).IsNull()
	s.AssertThat(child.FindClosest("x").GetValue().GetType()).IsEqualTo(s.numberType)
	s.AssertThat(child.FindClosestByKey(NewProperty_StringKey("x")).GetValue().GetType()).
		IsEqualTo(s.numberType)
}

func TestPropertyMap_LookUpKeysFromParent_symbolKey(t *testing.T) {
	s := newPropertyMapTest(t)
	parentType := s.registry.CreateRecordType(mapOf()).AssertObjectType()
	parentType.DefineDeclaredPropertyByKey(s.symbolKey, s.numberType, nil)
	childType := s.registry.CreateObjectType(ptr("Child"), parentType)
	child := NewPropertyMap()
	child.SetParentSource(childType)

	s.AssertThat(child.GetOwnPropertyByKey(s.symbolKey)).IsNull()
	s.AssertThat(child.FindClosestByKey(s.symbolKey).GetValue().GetType()).IsEqualTo(s.numberType)
}
