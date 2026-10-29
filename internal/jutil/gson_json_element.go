// Port of com/google/gson/JsonElement.java (Gson 2.9.1), with
// com/google/gson/JsonObject.java, JsonArray.java, JsonPrimitive.java,
// JsonNull.java and com/google/gson/internal/LazilyParsedNumber.java, and
// the LinkedTreeMap.Node.hashCode of com/google/gson/internal/LinkedTreeMap.java
// (in JsonObject.hashCode)

/*
 * Copyright (C) 2008 Google Inc. (JsonElement.java, JsonObject.java, JsonArray.java,
 *   JsonPrimitive.java, JsonNull.java)
 * Copyright (C) 2011 Google Inc. (LazilyParsedNumber.java)
 * Copyright (C) 2010 The Android Open Source Project, Copyright (C) 2012 Google Inc.
 *   (LinkedTreeMap.java)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jutil

import (
	"iter"
	"slices"
	"strings"
)

// Go-only: language-adapter — Gson 2.9.1's JSON tree (com.google.gson.
// JsonElement, JsonObject, JsonArray, JsonPrimitive, JsonNull and
// com.google.gson.internal.LazilyParsedNumber), as shaded into the pinned
// closure-compiler.jar. Ported from the Gson 2.9.1 sources; the getters
// with floating-point or BigInteger/BigDecimal results (getAsDouble,
// getAsFloat, getAsBigDecimal, ...) are not provided.

// JavaNumber is java.lang.Number as seen by Gson.
type JavaNumber interface {
	// Java: Number.intValue()
	IntValue() int32
	// Java: Number.longValue()
	LongValue() int64
	// Java: Number.doubleValue()
	DoubleValue() float64
	// Java: Number.toString()
	String() string
}

// JsonElement is com.google.gson.JsonElement.
type JsonElement interface {
	// Java: JsonElement.deepCopy()
	DeepCopy() JsonElement
	// Java: JsonElement.isJsonArray()
	IsJsonArray() bool
	// Java: JsonElement.isJsonObject()
	IsJsonObject() bool
	// Java: JsonElement.isJsonPrimitive()
	IsJsonPrimitive() bool
	// Java: JsonElement.isJsonNull()
	IsJsonNull() bool
	// Java: JsonElement.getAsJsonObject()
	GetAsJsonObject() *JsonObject
	// Java: JsonElement.getAsJsonArray()
	GetAsJsonArray() *JsonArray
	// Java: JsonElement.getAsJsonPrimitive()
	GetAsJsonPrimitive() *JsonPrimitive
	// Java: JsonElement.getAsBoolean()
	GetAsBoolean() bool
	// Java: JsonElement.getAsNumber()
	GetAsNumber() JavaNumber
	// Java: JsonElement.getAsString()
	GetAsString() string
	// Java: JsonElement.getAsLong()
	GetAsLong() int64
	// Java: JsonElement.getAsInt()
	GetAsInt() int32
	// Go-only: non-panicking accessors for callers that propagate errors.
	TryGetAsInt() (int32, Throwable)
	TryGetAsString() (string, Throwable)
	// Java: JsonElement.toString()
	String() string
	// Java: Object.equals(Object), overridden by every JsonElement subclass
	Equals(o any) bool
	// Java: Object.hashCode(), overridden by every JsonElement subclass
	HashCode() int32

	isJsonElement()
}

// jsonElementSimpleName is getClass().getSimpleName().
func jsonElementSimpleName(e JsonElement) string {
	switch e.(type) {
	case *JsonObject:
		return "JsonObject"
	case *JsonArray:
		return "JsonArray"
	case *JsonPrimitive:
		return "JsonPrimitive"
	case *JsonNull:
		return "JsonNull"
	}
	panic(NewIllegalStateException())
}

// jsonElementClassName is getClass().getName() (the shaded class name).
func jsonElementClassName(e JsonElement) string {
	return GsonPackage + "." + jsonElementSimpleName(e)
}

// jsonElementToString is the compact JSON form of e, written by a lenient
// JsonWriter.
func jsonElementToString(e JsonElement) string {
	var stringWriter strings.Builder
	jsonWriter := NewJsonWriter(&stringWriter)
	jsonWriter.SetLenient(true)
	Streams_Write(e, jsonWriter)
	return stringWriter.String()
}

// ---------------------------------------------------------------------------
// JsonObject

// JsonObject is com.google.gson.JsonObject: members in insertion order (a
// LinkedTreeMap in Gson; re-adding a member keeps its position).
type JsonObject struct {
	members *LinkedHashMap[string, JsonElement]
}

func (*JsonObject) isJsonElement() {}

// NewJsonObject is new JsonObject().
//
// Java: JsonObject()
func NewJsonObject() *JsonObject {
	return &JsonObject{members: NewLinkedHashMap[string, JsonElement]()}
}

func (*JsonObject) IsJsonArray() bool     { return false }
func (*JsonObject) IsJsonObject() bool    { return true }
func (*JsonObject) IsJsonPrimitive() bool { return false }
func (*JsonObject) IsJsonNull() bool      { return false }

func (o *JsonObject) GetAsJsonObject() *JsonObject { return o }
func (o *JsonObject) GetAsJsonArray() *JsonArray {
	panic(NewIllegalStateExceptionMsg("Not a JSON Array: " + o.String()))
}
func (o *JsonObject) GetAsJsonPrimitive() *JsonPrimitive {
	panic(NewIllegalStateExceptionMsg("Not a JSON Primitive: " + o.String()))
}

func (*JsonObject) GetAsBoolean() bool      { panic(NewUnsupportedOperationExceptionMsg("JsonObject")) }
func (*JsonObject) GetAsNumber() JavaNumber { panic(NewUnsupportedOperationExceptionMsg("JsonObject")) }
func (*JsonObject) GetAsString() string     { panic(NewUnsupportedOperationExceptionMsg("JsonObject")) }
func (*JsonObject) GetAsLong() int64        { panic(NewUnsupportedOperationExceptionMsg("JsonObject")) }
func (*JsonObject) GetAsInt() int32         { panic(NewUnsupportedOperationExceptionMsg("JsonObject")) }

func (*JsonObject) TryGetAsInt() (int32, Throwable) {
	return 0, NewUnsupportedOperationExceptionMsg("JsonObject")
}
func (*JsonObject) TryGetAsString() (string, Throwable) {
	return "", NewUnsupportedOperationExceptionMsg("JsonObject")
}

func (o *JsonObject) String() string { return jsonElementToString(o) }

// Java: JsonObject.deepCopy()
func (o *JsonObject) DeepCopy() JsonElement {
	result := NewJsonObject()
	for k, v := range o.members.All() {
		result.Add(k, v.DeepCopy())
	}
	return result
}

// Add adds a member; a nil value is stored as JsonNull.INSTANCE.
//
// Java: JsonObject.add(String, JsonElement)
func (o *JsonObject) Add(property string, value JsonElement) {
	if value == nil {
		value = JsonNull_INSTANCE
	}
	o.members.Put(property, value)
}

// Remove removes a member and returns it (nil if absent).
//
// Java: JsonObject.remove(String)
func (o *JsonObject) Remove(property string) JsonElement {
	return o.members.Remove(property)
}

// AddPropertyString adds a string member.
//
// Java: JsonObject.addProperty(String, String)
func (o *JsonObject) AddPropertyString(property string, value *string) {
	if value == nil {
		o.Add(property, JsonNull_INSTANCE)
	} else {
		o.Add(property, NewJsonPrimitiveString(*value))
	}
}

// EntrySet iterates the members in insertion order.
//
// Java: JsonObject.entrySet()
func (o *JsonObject) EntrySet() *LinkedHashMap[string, JsonElement] {
	return o.members
}

// Size is the number of members.
//
// Java: JsonObject.size()
func (o *JsonObject) Size() int {
	return o.members.Size()
}

// Has reports whether the member exists.
//
// Java: JsonObject.has(String)
func (o *JsonObject) Has(memberName string) bool {
	return o.members.ContainsKey(memberName)
}

// Get returns the member, or nil (Java null) if absent.
//
// Java: JsonObject.get(String)
func (o *JsonObject) Get(memberName string) JsonElement {
	return o.members.Get(memberName)
}

// Equals reports whether other is a JsonObject with equal members.
//
// Java: JsonObject.equals(Object)
func (o *JsonObject) Equals(other any) bool {
	if other == any(o) {
		return true
	}
	if x, ok := other.(*JsonObject); ok && x != nil {
		return jsonObject_membersEqual(x.members, o.members)
	}
	return false
}

// jsonObject_membersEqual is LinkedTreeMap.equals (Map.equals): the same
// keys, each mapped to equal values.
func jsonObject_membersEqual(a, b *LinkedHashMap[string, JsonElement]) bool {
	if a.Size() != b.Size() {
		return false
	}
	for k, v := range a.All() {
		w, ok := b.Lookup(k)
		if !ok || !v.Equals(w) {
			return false
		}
	}
	return true
}

// HashCode is the members' Map.hashCode(): the sum of the entries' hash
// codes, each LinkedTreeMap.Node.hashCode(), key.hashCode() ^
// value.hashCode().
//
// Java: JsonObject.hashCode()
func (o *JsonObject) HashCode() int32 {
	h := int32(0)
	for k, v := range o.members.All() {
		h += StringHashCode(k) ^ v.HashCode()
	}
	return h
}

// ---------------------------------------------------------------------------
// JsonArray

// JsonArray is com.google.gson.JsonArray.
type JsonArray struct {
	elements []JsonElement
}

func (*JsonArray) isJsonElement() {}

// NewJsonArray is new JsonArray().
//
// Java: JsonArray()
func NewJsonArray() *JsonArray {
	return &JsonArray{}
}

func (*JsonArray) IsJsonArray() bool     { return true }
func (*JsonArray) IsJsonObject() bool    { return false }
func (*JsonArray) IsJsonPrimitive() bool { return false }
func (*JsonArray) IsJsonNull() bool      { return false }

func (a *JsonArray) GetAsJsonObject() *JsonObject {
	panic(NewIllegalStateExceptionMsg("Not a JSON Object: " + a.String()))
}
func (a *JsonArray) GetAsJsonArray() *JsonArray { return a }
func (a *JsonArray) GetAsJsonPrimitive() *JsonPrimitive {
	panic(NewIllegalStateExceptionMsg("Not a JSON Primitive: " + a.String()))
}

func (a *JsonArray) String() string { return jsonElementToString(a) }

// Java: JsonArray.deepCopy()
func (a *JsonArray) DeepCopy() JsonElement {
	if len(a.elements) == 0 {
		return NewJsonArray()
	}
	elems := make([]JsonElement, len(a.elements))
	for i, e := range a.elements {
		elems[i] = e.DeepCopy()
	}
	return &JsonArray{elements: elems}
}

// Add appends an element; nil is stored as JsonNull.INSTANCE.
//
// Java: JsonArray.add(JsonElement)
func (a *JsonArray) Add(element JsonElement) {
	if element == nil {
		element = JsonNull_INSTANCE
	}
	a.elements = append(a.elements, element)
}

// AddString appends a string element.
//
// Java: JsonArray.add(String)
func (a *JsonArray) AddString(s *string) {
	if s == nil {
		a.elements = append(a.elements, JsonNull_INSTANCE)
	} else {
		a.elements = append(a.elements, NewJsonPrimitiveString(*s))
	}
}

// Size is the number of elements.
//
// Java: JsonArray.size()
func (a *JsonArray) Size() int {
	return len(a.elements)
}

// IsEmpty reports whether there are no elements.
//
// Java: JsonArray.isEmpty()
func (a *JsonArray) IsEmpty() bool {
	return len(a.elements) == 0
}

// All iterates the elements in order.
//
// Java: JsonArray.iterator()
func (a *JsonArray) All() iter.Seq[JsonElement] {
	return slices.Values(a.elements)
}

// Get returns the i-th element.
//
// Java: JsonArray.get(int)
func (a *JsonArray) Get(i int) JsonElement {
	if uint(i) >= uint(len(a.elements)) {
		panic(NewIndexOutOfBoundsException(LenientFormat("Index %s out of bounds for length %s", i, len(a.elements))))
	}
	return a.elements[i]
}

// Java: JsonArray.getAsNumber()
func (a *JsonArray) GetAsNumber() JavaNumber {
	if len(a.elements) == 1 {
		return a.elements[0].GetAsNumber()
	}
	panic(NewIllegalStateException())
}

// Java: JsonArray.getAsString()
func (a *JsonArray) GetAsString() string {
	if len(a.elements) == 1 {
		return a.elements[0].GetAsString()
	}
	panic(NewIllegalStateException())
}

func (a *JsonArray) TryGetAsString() (string, Throwable) {
	if len(a.elements) == 1 {
		return a.elements[0].TryGetAsString()
	}
	return "", NewIllegalStateException()
}

// Java: JsonArray.getAsLong()
func (a *JsonArray) GetAsLong() int64 {
	if len(a.elements) == 1 {
		return a.elements[0].GetAsLong()
	}
	panic(NewIllegalStateException())
}

// Java: JsonArray.getAsInt()
func (a *JsonArray) GetAsInt() int32 {
	if len(a.elements) == 1 {
		return a.elements[0].GetAsInt()
	}
	panic(NewIllegalStateException())
}

func (a *JsonArray) TryGetAsInt() (int32, Throwable) {
	if len(a.elements) == 1 {
		return a.elements[0].TryGetAsInt()
	}
	return 0, NewIllegalStateException()
}

// Java: JsonArray.getAsBoolean()
func (a *JsonArray) GetAsBoolean() bool {
	if len(a.elements) == 1 {
		return a.elements[0].GetAsBoolean()
	}
	panic(NewIllegalStateException())
}

// Equals reports whether other is a JsonArray with equal elements, in the
// same order (ArrayList.equals).
//
// Java: JsonArray.equals(Object)
func (a *JsonArray) Equals(o any) bool {
	if o == any(a) {
		return true
	}
	x, ok := o.(*JsonArray)
	if !ok || x == nil || len(x.elements) != len(a.elements) {
		return false
	}
	for i, e := range a.elements {
		if !x.elements[i].Equals(e) {
			return false
		}
	}
	return true
}

// HashCode is the elements' List.hashCode().
//
// Java: JsonArray.hashCode()
func (a *JsonArray) HashCode() int32 {
	h := int32(1)
	for _, e := range a.elements {
		h = 31*h + e.HashCode()
	}
	return h
}

// ---------------------------------------------------------------------------
// JsonPrimitive

// JsonPrimitive is com.google.gson.JsonPrimitive: a string, a boolean or a
// number (value holds a Go string, a Go bool or a JavaNumber).
type JsonPrimitive struct {
	value any
}

func (*JsonPrimitive) isJsonElement() {}

// NewJsonPrimitiveBoolean is new JsonPrimitive(Boolean).
//
// Java: JsonPrimitive(Boolean)
func NewJsonPrimitiveBoolean(b bool) *JsonPrimitive {
	return &JsonPrimitive{value: b}
}

// NewJsonPrimitiveNumber is new JsonPrimitive(Number).
//
// Java: JsonPrimitive(Number)
func NewJsonPrimitiveNumber(number JavaNumber) *JsonPrimitive {
	return &JsonPrimitive{value: CheckNotNull(number)}
}

// NewJsonPrimitiveString is new JsonPrimitive(String).
//
// Java: JsonPrimitive(String)
func NewJsonPrimitiveString(s string) *JsonPrimitive {
	return &JsonPrimitive{value: s}
}

func (*JsonPrimitive) IsJsonArray() bool     { return false }
func (*JsonPrimitive) IsJsonObject() bool    { return false }
func (*JsonPrimitive) IsJsonPrimitive() bool { return true }
func (*JsonPrimitive) IsJsonNull() bool      { return false }

func (p *JsonPrimitive) GetAsJsonObject() *JsonObject {
	panic(NewIllegalStateExceptionMsg("Not a JSON Object: " + p.String()))
}
func (p *JsonPrimitive) GetAsJsonArray() *JsonArray {
	panic(NewIllegalStateExceptionMsg("Not a JSON Array: " + p.String()))
}
func (p *JsonPrimitive) GetAsJsonPrimitive() *JsonPrimitive { return p }

func (p *JsonPrimitive) String() string { return jsonElementToString(p) }

// Java: JsonPrimitive.deepCopy()
func (p *JsonPrimitive) DeepCopy() JsonElement {
	return p
}

// IsBoolean reports whether the value is a boolean.
//
// Java: JsonPrimitive.isBoolean()
func (p *JsonPrimitive) IsBoolean() bool {
	_, ok := p.value.(bool)
	return ok
}

// Java: JsonPrimitive.getAsBoolean()
func (p *JsonPrimitive) GetAsBoolean() bool {
	if b, ok := p.value.(bool); ok {
		return b
	}
	return booleanParseBoolean(p.GetAsString())
}

// booleanParseBoolean is Boolean.parseBoolean(String): "true", ignoring case.
func booleanParseBoolean(s string) bool {
	return len(s) == 4 && StringToLowerCaseASCII(s) == "true"
}

// IsNumber reports whether the value is a number.
//
// Java: JsonPrimitive.isNumber()
func (p *JsonPrimitive) IsNumber() bool {
	_, ok := p.value.(JavaNumber)
	return ok
}

// Java: JsonPrimitive.getAsNumber()
func (p *JsonPrimitive) GetAsNumber() JavaNumber {
	if s, ok := p.value.(string); ok {
		return NewLazilyParsedNumber(s)
	}
	return p.value.(JavaNumber)
}

// IsString reports whether the value is a string.
//
// Java: JsonPrimitive.isString()
func (p *JsonPrimitive) IsString() bool {
	_, ok := p.value.(string)
	return ok
}

// Java: JsonPrimitive.getAsString()
func (p *JsonPrimitive) GetAsString() string {
	if num, ok := p.value.(JavaNumber); ok {
		return num.String()
	}
	if b, ok := p.value.(bool); ok {
		return ToJavaString(b)
	}
	return p.value.(string)
}

func (p *JsonPrimitive) TryGetAsString() (string, Throwable) {
	return p.GetAsString(), nil
}

// Java: JsonPrimitive.getAsLong()
func (p *JsonPrimitive) GetAsLong() int64 {
	if num, ok := p.value.(JavaNumber); ok {
		return num.LongValue()
	}
	return LongParseLong(p.GetAsString(), 10)
}

// Java: JsonPrimitive.getAsInt()
func (p *JsonPrimitive) GetAsInt() int32 {
	if num, ok := p.value.(JavaNumber); ok {
		return num.IntValue()
	}
	return IntegerParseInt(p.GetAsString(), 10)
}

func (p *JsonPrimitive) TryGetAsInt() (int32, Throwable) {
	if lpn, ok := p.value.(*LazilyParsedNumber); ok {
		return lpn.tryIntValue()
	}
	if num, ok := p.value.(JavaNumber); ok {
		return num.IntValue(), nil
	}
	s := p.GetAsString()
	if v, ok := TryIntegerParseInt(s, 10); ok {
		return v, nil
	}
	return 0, badNumberInput(s, 10)
}

// Java: JsonPrimitive.hashCode()
func (p *JsonPrimitive) HashCode() int32 {
	if p.value == nil {
		return 31
	}
	// Using recommended hashing algorithm from Effective Java for longs and doubles
	if jsonPrimitive_isIntegral(p) {
		value := p.GetAsNumber().LongValue()
		return int32(value ^ Ushr64(value, 32))
	}
	if _, ok := p.value.(JavaNumber); ok {
		value := DoubleToLongBits(p.GetAsNumber().DoubleValue())
		return int32(value ^ Ushr64(value, 32))
	}
	return ObjectHashCode(p.value)
}

// Java: JsonPrimitive.equals(Object)
func (p *JsonPrimitive) Equals(obj any) bool {
	if obj == any(p) {
		return true
	}
	other, ok := obj.(*JsonPrimitive)
	if !ok || other == nil {
		return false
	}
	if p.value == nil {
		return other.value == nil
	}
	if jsonPrimitive_isIntegral(p) && jsonPrimitive_isIntegral(other) {
		return p.GetAsNumber().LongValue() == other.GetAsNumber().LongValue()
	}
	_, numberA := p.value.(JavaNumber)
	_, numberB := other.value.(JavaNumber)
	if numberA && numberB {
		a := p.GetAsNumber().DoubleValue()
		// Java standard types other than double return true for two NaN. So, need
		// special handling for double.
		b := other.GetAsNumber().DoubleValue()
		return a == b || (a != a && b != b)
	}
	// value.equals(other.value): a String, a Boolean, or a Number that is
	// not equal to the String or Boolean other.value.
	if numberA {
		return false
	}
	return p.value == other.value
}

// jsonPrimitive_isIntegral returns true if the specified number is an
// integral type (Long, Integer, Short, Byte, BigInteger).
//
// Java: JsonPrimitive.isIntegral(JsonPrimitive)
func jsonPrimitive_isIntegral(primitive *JsonPrimitive) bool {
	switch primitive.value.(type) {
	case Integer, gson_javaInteger:
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// JsonNull

// JsonNull is com.google.gson.JsonNull.
type JsonNull struct{}

func (*JsonNull) isJsonElement() {}

// JsonNull_INSTANCE is JsonNull.INSTANCE.
var JsonNull_INSTANCE = &JsonNull{}

func (*JsonNull) IsJsonArray() bool     { return false }
func (*JsonNull) IsJsonObject() bool    { return false }
func (*JsonNull) IsJsonPrimitive() bool { return false }
func (*JsonNull) IsJsonNull() bool      { return true }

func (n *JsonNull) GetAsJsonObject() *JsonObject {
	panic(NewIllegalStateExceptionMsg("Not a JSON Object: " + n.String()))
}
func (n *JsonNull) GetAsJsonArray() *JsonArray {
	panic(NewIllegalStateExceptionMsg("Not a JSON Array: " + n.String()))
}
func (n *JsonNull) GetAsJsonPrimitive() *JsonPrimitive {
	panic(NewIllegalStateExceptionMsg("Not a JSON Primitive: " + n.String()))
}

func (*JsonNull) GetAsBoolean() bool      { panic(NewUnsupportedOperationExceptionMsg("JsonNull")) }
func (*JsonNull) GetAsNumber() JavaNumber { panic(NewUnsupportedOperationExceptionMsg("JsonNull")) }
func (*JsonNull) GetAsString() string     { panic(NewUnsupportedOperationExceptionMsg("JsonNull")) }
func (*JsonNull) GetAsLong() int64        { panic(NewUnsupportedOperationExceptionMsg("JsonNull")) }
func (*JsonNull) GetAsInt() int32         { panic(NewUnsupportedOperationExceptionMsg("JsonNull")) }

func (*JsonNull) TryGetAsInt() (int32, Throwable) {
	return 0, NewUnsupportedOperationExceptionMsg("JsonNull")
}
func (*JsonNull) TryGetAsString() (string, Throwable) {
	return "", NewUnsupportedOperationExceptionMsg("JsonNull")
}

func (n *JsonNull) String() string { return jsonElementToString(n) }

// Java: JsonNull.deepCopy()
func (n *JsonNull) DeepCopy() JsonElement {
	return JsonNull_INSTANCE
}

// HashCode is JsonNull.class.hashCode().
//
// Go-only: language-adapter — Java returns the identity hash of the
// JsonNull class object, which is arbitrary; this is a fixed value (the
// hash of the class name).
//
// Java: JsonNull.hashCode()
func (n *JsonNull) HashCode() int32 {
	return StringHashCode(jsonElementClassName(n))
}

// Equals reports whether other is a JsonNull: all instances of JsonNull
// are the same.
//
// Java: JsonNull.equals(Object)
func (n *JsonNull) Equals(other any) bool {
	if other == any(n) {
		return true
	}
	x, ok := other.(*JsonNull)
	return ok && x != nil
}

// ---------------------------------------------------------------------------
// LazilyParsedNumber

// LazilyParsedNumber is com.google.gson.internal.LazilyParsedNumber: a
// number kept as its source text.
type LazilyParsedNumber struct {
	value string
}

// NewLazilyParsedNumber is new LazilyParsedNumber(String).
//
// Java: LazilyParsedNumber(String)
func NewLazilyParsedNumber(value string) *LazilyParsedNumber {
	return &LazilyParsedNumber{value: value}
}

// Java: LazilyParsedNumber.intValue()
func (n *LazilyParsedNumber) IntValue() int32 {
	if v, ok := TryIntegerParseInt(n.value, 10); ok {
		return v
	}
	if v, ok := TryParseLong(n.value, 10); ok {
		return int32(v)
	}
	return NewBigDecimalFromString(n.value).IntValue()
}

func (n *LazilyParsedNumber) tryIntValue() (result int32, err Throwable) {
	if v, ok := TryIntegerParseInt(n.value, 10); ok {
		return v, nil
	}
	if v, ok := TryParseLong(n.value, 10); ok {
		return int32(v), nil
	}
	Try(func() {
		result = NewBigDecimalFromString(n.value).IntValue()
	}, func(nfe INumberFormatException) {
		err = nfe
	})
	return result, err
}

// Java: LazilyParsedNumber.longValue()
func (n *LazilyParsedNumber) LongValue() int64 {
	if v, ok := TryParseLong(n.value, 10); ok {
		return v
	}
	return NewBigDecimalFromString(n.value).LongValue()
}

// Java: LazilyParsedNumber.doubleValue()
func (n *LazilyParsedNumber) DoubleValue() float64 {
	return DoubleParseDouble(n.value)
}

// Java: LazilyParsedNumber.toString()
func (n *LazilyParsedNumber) String() string {
	return n.value
}

// Java: LazilyParsedNumber.hashCode()
func (n *LazilyParsedNumber) HashCode() int32 {
	return StringHashCode(n.value)
}

// Java: LazilyParsedNumber.equals(Object)
func (n *LazilyParsedNumber) Equals(obj any) bool {
	if obj == any(n) {
		return true
	}
	if other, ok := obj.(*LazilyParsedNumber); ok && other != nil {
		return n.value == other.value
	}
	return false
}
