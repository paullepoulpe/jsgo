// Port of com/google/gson/JsonElement.java (Gson 2.9.1), with
// com/google/gson/JsonObject.java, JsonArray.java, JsonPrimitive.java,
// JsonNull.java and com/google/gson/internal/LazilyParsedNumber.java

/*
 * Copyright (C) 2008 Google Inc. (JsonElement.java, JsonObject.java, JsonArray.java,
 *   JsonPrimitive.java, JsonNull.java)
 * Copyright (C) 2011 Google Inc. (LazilyParsedNumber.java)
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
	"strings"
)

// Go-only: language-adapter — Gson 2.9.1's JSON tree (com.google.gson.
// JsonElement, JsonObject, JsonArray, JsonPrimitive, JsonNull and
// com.google.gson.internal.LazilyParsedNumber), as shaded into the pinned
// closure-compiler.jar. Ported from the Gson 2.9.1 sources; members that
// need Java's floating-point parsing (getAsDouble, getAsFloat, ...),
// BigInteger/BigDecimal results, or equals/hashCode are not provided.

// JavaNumber is java.lang.Number as seen by Gson.
type JavaNumber interface {
	// Java: Number.intValue()
	IntValue() int32
	// Java: Number.longValue()
	LongValue() int64
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
	// Java: JsonElement.getAsJsonNull()
	GetAsJsonNull() *JsonNull
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
	// Java: JsonElement.toString()
	String() string

	// Go-only: language-adapter — field access through the interface.
	GetJsonElementBase() *JsonElementBase
}

// JsonElementBase holds the base implementations of JsonElement.
type JsonElementBase struct {
	// Go-only: language-adapter — the dynamic receiver for virtual calls.
	self JsonElement
}

// GetJsonElementBase returns the base.
func (b *JsonElementBase) GetJsonElementBase() *JsonElementBase { return b }

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

// IsJsonArray reports whether this element is an array.
//
// Java: JsonElement.isJsonArray()
func (b *JsonElementBase) IsJsonArray() bool {
	_, ok := b.self.(*JsonArray)
	return ok
}

// IsJsonObject reports whether this element is an object.
//
// Java: JsonElement.isJsonObject()
func (b *JsonElementBase) IsJsonObject() bool {
	_, ok := b.self.(*JsonObject)
	return ok
}

// IsJsonPrimitive reports whether this element is a primitive.
//
// Java: JsonElement.isJsonPrimitive()
func (b *JsonElementBase) IsJsonPrimitive() bool {
	_, ok := b.self.(*JsonPrimitive)
	return ok
}

// IsJsonNull reports whether this element is the JSON null.
//
// Java: JsonElement.isJsonNull()
func (b *JsonElementBase) IsJsonNull() bool {
	_, ok := b.self.(*JsonNull)
	return ok
}

// GetAsJsonObject returns this element as an object, or throws.
//
// Java: JsonElement.getAsJsonObject()
func (b *JsonElementBase) GetAsJsonObject() *JsonObject {
	if b.self.IsJsonObject() {
		return b.self.(*JsonObject)
	}
	panic(NewIllegalStateExceptionMsg("Not a JSON Object: " + b.self.String()))
}

// GetAsJsonArray returns this element as an array, or throws.
//
// Java: JsonElement.getAsJsonArray()
func (b *JsonElementBase) GetAsJsonArray() *JsonArray {
	if b.self.IsJsonArray() {
		return b.self.(*JsonArray)
	}
	panic(NewIllegalStateExceptionMsg("Not a JSON Array: " + b.self.String()))
}

// GetAsJsonPrimitive returns this element as a primitive, or throws.
//
// Java: JsonElement.getAsJsonPrimitive()
func (b *JsonElementBase) GetAsJsonPrimitive() *JsonPrimitive {
	if b.self.IsJsonPrimitive() {
		return b.self.(*JsonPrimitive)
	}
	panic(NewIllegalStateExceptionMsg("Not a JSON Primitive: " + b.self.String()))
}

// GetAsJsonNull returns this element as the JSON null, or throws.
//
// Java: JsonElement.getAsJsonNull()
func (b *JsonElementBase) GetAsJsonNull() *JsonNull {
	if b.self.IsJsonNull() {
		return b.self.(*JsonNull)
	}
	panic(NewIllegalStateExceptionMsg("Not a JSON Null: " + b.self.String()))
}

// Java: JsonElement.getAsBoolean()
func (b *JsonElementBase) GetAsBoolean() bool {
	panic(NewUnsupportedOperationExceptionMsg(jsonElementSimpleName(b.self)))
}

// Java: JsonElement.getAsNumber()
func (b *JsonElementBase) GetAsNumber() JavaNumber {
	panic(NewUnsupportedOperationExceptionMsg(jsonElementSimpleName(b.self)))
}

// Java: JsonElement.getAsString()
func (b *JsonElementBase) GetAsString() string {
	panic(NewUnsupportedOperationExceptionMsg(jsonElementSimpleName(b.self)))
}

// Java: JsonElement.getAsLong()
func (b *JsonElementBase) GetAsLong() int64 {
	panic(NewUnsupportedOperationExceptionMsg(jsonElementSimpleName(b.self)))
}

// Java: JsonElement.getAsInt()
func (b *JsonElementBase) GetAsInt() int32 {
	panic(NewUnsupportedOperationExceptionMsg(jsonElementSimpleName(b.self)))
}

// String is the compact JSON form of this element, written by a lenient
// JsonWriter.
//
// Java: JsonElement.toString()
func (b *JsonElementBase) String() string {
	var stringWriter strings.Builder
	jsonWriter := NewJsonWriter(&stringWriter)
	jsonWriter.SetLenient(true)
	Streams_Write(b.self, jsonWriter)
	return stringWriter.String()
}

// ---------------------------------------------------------------------------
// JsonObject

// JsonObject is com.google.gson.JsonObject: members in insertion order (a
// LinkedTreeMap in Gson; re-adding a member keeps its position).
type JsonObject struct {
	JsonElementBase
	members *LinkedHashMap[string, JsonElement]
}

// NewJsonObject is new JsonObject().
//
// Java: JsonObject()
func NewJsonObject() *JsonObject {
	o := &JsonObject{members: NewLinkedHashMap[string, JsonElement]()}
	o.self = o
	return o
}

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

// AddPropertyBoolean adds a boolean member.
//
// Java: JsonObject.addProperty(String, Boolean)
func (o *JsonObject) AddPropertyBoolean(property string, value bool) {
	o.Add(property, NewJsonPrimitiveBoolean(value))
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

// ---------------------------------------------------------------------------
// JsonArray

// JsonArray is com.google.gson.JsonArray.
type JsonArray struct {
	JsonElementBase
	elements *ArrayList[JsonElement]
}

// NewJsonArray is new JsonArray().
//
// Java: JsonArray()
func NewJsonArray() *JsonArray {
	a := &JsonArray{elements: NewArrayList[JsonElement]()}
	a.self = a
	return a
}

// Java: JsonArray.deepCopy()
func (a *JsonArray) DeepCopy() JsonElement {
	result := NewJsonArray()
	for e := range a.elements.All() {
		result.Add(e.DeepCopy())
	}
	return result
}

// Add appends an element; nil is stored as JsonNull.INSTANCE.
//
// Java: JsonArray.add(JsonElement)
func (a *JsonArray) Add(element JsonElement) {
	if element == nil {
		element = JsonNull_INSTANCE
	}
	a.elements.Add(element)
}

// AddString appends a string element.
//
// Java: JsonArray.add(String)
func (a *JsonArray) AddString(s *string) {
	if s == nil {
		a.elements.Add(JsonNull_INSTANCE)
	} else {
		a.elements.Add(NewJsonPrimitiveString(*s))
	}
}

// Size is the number of elements.
//
// Java: JsonArray.size()
func (a *JsonArray) Size() int {
	return a.elements.Size()
}

// IsEmpty reports whether there are no elements.
//
// Java: JsonArray.isEmpty()
func (a *JsonArray) IsEmpty() bool {
	return a.elements.IsEmpty()
}

// All iterates the elements in order.
//
// Java: JsonArray.iterator()
func (a *JsonArray) All() iter.Seq[JsonElement] {
	return a.elements.All()
}

// Get returns the i-th element.
//
// Java: JsonArray.get(int)
func (a *JsonArray) Get(i int) JsonElement {
	return a.elements.Get(i)
}

// Java: JsonArray.getAsNumber()
func (a *JsonArray) GetAsNumber() JavaNumber {
	if a.elements.Size() == 1 {
		return a.elements.Get(0).GetAsNumber()
	}
	panic(NewIllegalStateException())
}

// Java: JsonArray.getAsString()
func (a *JsonArray) GetAsString() string {
	if a.elements.Size() == 1 {
		return a.elements.Get(0).GetAsString()
	}
	panic(NewIllegalStateException())
}

// Java: JsonArray.getAsLong()
func (a *JsonArray) GetAsLong() int64 {
	if a.elements.Size() == 1 {
		return a.elements.Get(0).GetAsLong()
	}
	panic(NewIllegalStateException())
}

// Java: JsonArray.getAsInt()
func (a *JsonArray) GetAsInt() int32 {
	if a.elements.Size() == 1 {
		return a.elements.Get(0).GetAsInt()
	}
	panic(NewIllegalStateException())
}

// Java: JsonArray.getAsBoolean()
func (a *JsonArray) GetAsBoolean() bool {
	if a.elements.Size() == 1 {
		return a.elements.Get(0).GetAsBoolean()
	}
	panic(NewIllegalStateException())
}

// ---------------------------------------------------------------------------
// JsonPrimitive

// JsonPrimitive is com.google.gson.JsonPrimitive: a string, a boolean or a
// number (value holds a Go string, a Go bool or a JavaNumber).
type JsonPrimitive struct {
	JsonElementBase
	value any
}

// NewJsonPrimitiveBoolean is new JsonPrimitive(Boolean).
//
// Java: JsonPrimitive(Boolean)
func NewJsonPrimitiveBoolean(b bool) *JsonPrimitive {
	p := &JsonPrimitive{value: b}
	p.self = p
	return p
}

// NewJsonPrimitiveNumber is new JsonPrimitive(Number).
//
// Java: JsonPrimitive(Number)
func NewJsonPrimitiveNumber(number JavaNumber) *JsonPrimitive {
	p := &JsonPrimitive{value: CheckNotNull(number)}
	p.self = p
	return p
}

// NewJsonPrimitiveString is new JsonPrimitive(String).
//
// Java: JsonPrimitive(String)
func NewJsonPrimitiveString(s string) *JsonPrimitive {
	p := &JsonPrimitive{value: s}
	p.self = p
	return p
}

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
	// Check to see if the value as a String is "true" in any case.
	return booleanParseBoolean(p.GetAsString())
}

// booleanParseBoolean is Boolean.parseBoolean(String): "true", ignoring
// case.
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
	if p.IsNumber() {
		return p.GetAsNumber().String()
	} else if b, ok := p.value.(bool); ok {
		return ToJavaString(b)
	} else {
		return p.value.(string)
	}
}

// Java: JsonPrimitive.getAsLong()
func (p *JsonPrimitive) GetAsLong() int64 {
	if p.IsNumber() {
		return p.GetAsNumber().LongValue()
	}
	return LongParseLong(p.GetAsString(), 10)
}

// Java: JsonPrimitive.getAsInt()
func (p *JsonPrimitive) GetAsInt() int32 {
	if p.IsNumber() {
		return p.GetAsNumber().IntValue()
	}
	return IntegerParseInt(p.GetAsString(), 10)
}

// ---------------------------------------------------------------------------
// JsonNull

// JsonNull is com.google.gson.JsonNull.
type JsonNull struct {
	JsonElementBase
}

// JsonNull_INSTANCE is JsonNull.INSTANCE.
var JsonNull_INSTANCE = newJsonNull()

func newJsonNull() *JsonNull {
	n := &JsonNull{}
	n.self = n
	return n
}

// Java: JsonNull.deepCopy()
func (n *JsonNull) DeepCopy() JsonElement {
	return JsonNull_INSTANCE
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
func (n *LazilyParsedNumber) IntValue() (result int32) {
	Try(func() {
		result = IntegerParseInt(n.value, 10)
	}, func(e INumberFormatException) {
		Try(func() {
			result = int32(LongParseLong(n.value, 10))
		}, func(nfe INumberFormatException) {
			result = NewBigDecimalFromString(n.value).IntValue()
		})
	})
	return result
}

// Java: LazilyParsedNumber.longValue()
func (n *LazilyParsedNumber) LongValue() (result int64) {
	Try(func() {
		result = LongParseLong(n.value, 10)
	}, func(e INumberFormatException) {
		result = NewBigDecimalFromString(n.value).LongValue()
	})
	return result
}

// Java: LazilyParsedNumber.toString()
func (n *LazilyParsedNumber) String() string {
	return n.value
}
