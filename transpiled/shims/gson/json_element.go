// Clean-room Gson 2.9.1 JSON DOM tree (Apache 2.0).

package gson

import (
	"iter"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javautil"
)

// JavaNumber represents java.lang.Number values stored in JsonPrimitive.
type JavaNumber interface {
	IntValue() int32
	LongValue() int64
	DoubleValue() float64
	String() string
}

// LazilyParsedNumber implements com.google.gson.internal.LazilyParsedNumber.
type LazilyParsedNumber struct {
	value string
}

func NewLazilyParsedNumber(value string) *LazilyParsedNumber {
	return &LazilyParsedNumber{value: value}
}

func (n *LazilyParsedNumber) IntValue() int32 {
	if v, err := strconv.ParseInt(n.value, 10, 32); err == nil {
		return int32(v)
	}
	if v, err := strconv.ParseInt(n.value, 10, 64); err == nil {
		return int32(v)
	}
	return javalang.D2I(javalang.DoubleParseDouble(n.value))
}

func (n *LazilyParsedNumber) LongValue() int64 {
	if v, err := strconv.ParseInt(n.value, 10, 64); err == nil {
		return v
	}
	return javalang.D2L(javalang.DoubleParseDouble(n.value))
}

func (n *LazilyParsedNumber) DoubleValue() float64 {
	return javalang.DoubleParseDouble(n.value)
}

func (n *LazilyParsedNumber) String() string { return n.value }
func (n *LazilyParsedNumber) HashCode() int32 {
	return javalang.StringHashCode(n.value)
}
func (n *LazilyParsedNumber) Equals(other any) bool {
	o, ok := other.(*LazilyParsedNumber)
	return ok && o != nil && o.value == n.value
}

type javaInteger int64

func (n javaInteger) IntValue() int32      { return int32(n) }
func (n javaInteger) LongValue() int64     { return int64(n) }
func (n javaInteger) DoubleValue() float64 { return float64(n) }
func (n javaInteger) String() string       { return strconv.FormatInt(int64(n), 10) }

type javaDouble float64

func (n javaDouble) IntValue() int32      { return javalang.D2I(float64(n)) }
func (n javaDouble) LongValue() int64     { return javalang.D2L(float64(n)) }
func (n javaDouble) DoubleValue() float64 { return float64(n) }
func (n javaDouble) String() string       { return javalang.DoubleToString(float64(n)) }

// JsonElement implements com.google.gson.JsonElement.
type JsonElement interface {
	DeepCopy() JsonElement
	IsJsonArray() bool
	IsJsonObject() bool
	IsJsonPrimitive() bool
	IsJsonNull() bool
	GetAsJsonObject() *JsonObject
	GetAsJsonArray() *JsonArray
	GetAsJsonPrimitive() *JsonPrimitive
	GetAsBoolean() bool
	GetAsNumber() JavaNumber
	GetAsString() string
	GetAsLong() int64
	GetAsInt() int32
	GetAsDouble() float64
	String() string
	Equals(o any) bool
	HashCode() int32
	isJsonElement()
}

func jsonElementToString(e JsonElement) string {
	var sb strings.Builder
	jw := NewJsonWriterSink(&sb)
	jw.SetLenient(true)
	Streams_Write(e, jw)
	return sb.String()
}

// ---------------------------------------------------------------------------
// JsonObject
// ---------------------------------------------------------------------------

// JsonObject implements com.google.gson.JsonObject preserving insertion order.
type JsonObject struct {
	members *javautil.LinkedHashMap[string, JsonElement]
}

func (*JsonObject) isJsonElement() {}

func NewJsonObject() *JsonObject {
	return &JsonObject{members: javautil.NewLinkedHashMap[string, JsonElement]()}
}

func (o *JsonObject) ensureMap() *javautil.LinkedHashMap[string, JsonElement] {
	if o.members == nil {
		o.members = javautil.NewLinkedHashMap[string, JsonElement]()
	}
	return o.members
}

func (*JsonObject) IsJsonArray() bool     { return false }
func (*JsonObject) IsJsonObject() bool    { return true }
func (*JsonObject) IsJsonPrimitive() bool { return false }
func (*JsonObject) IsJsonNull() bool      { return false }

func (o *JsonObject) GetAsJsonObject() *JsonObject { return o }
func (o *JsonObject) GetAsJsonArray() *JsonArray {
	panic(javalang.NewIllegalStateException("Not a JSON Array: " + o.String()))
}
func (o *JsonObject) GetAsJsonPrimitive() *JsonPrimitive {
	panic(javalang.NewIllegalStateException("Not a JSON Primitive: " + o.String()))
}

func (*JsonObject) GetAsBoolean() bool      { panic(javalang.NewUnsupportedOperationException("JsonObject")) }
func (*JsonObject) GetAsNumber() JavaNumber { panic(javalang.NewUnsupportedOperationException("JsonObject")) }
func (*JsonObject) GetAsString() string     { panic(javalang.NewUnsupportedOperationException("JsonObject")) }
func (*JsonObject) GetAsLong() int64        { panic(javalang.NewUnsupportedOperationException("JsonObject")) }
func (*JsonObject) GetAsInt() int32         { panic(javalang.NewUnsupportedOperationException("JsonObject")) }
func (*JsonObject) GetAsDouble() float64    { panic(javalang.NewUnsupportedOperationException("JsonObject")) }

func (o *JsonObject) String() string { return jsonElementToString(o) }

func (o *JsonObject) DeepCopy() JsonElement {
	cp := NewJsonObject()
	if o.members != nil {
		for k, v := range o.members.All() {
			cp.Add(k, v.DeepCopy())
		}
	}
	return cp
}

func (o *JsonObject) Add(property string, value JsonElement) {
	if value == nil {
		value = JsonNull_INSTANCE
	}
	o.ensureMap().Put(property, value)
}

func (o *JsonObject) Remove(property string) JsonElement {
	if o.members == nil {
		return nil
	}
	v, _ := o.members.Remove(property)
	return v
}

func (o *JsonObject) AddPropertyString(property, value string) {
	o.Add(property, NewJsonPrimitiveString(value))
}

func (o *JsonObject) AddPropertyBoolean(property string, value bool) {
	o.Add(property, NewJsonPrimitiveBoolean(value))
}

func (o *JsonObject) AddPropertyLong(property string, value int64) {
	o.Add(property, NewJsonPrimitiveNumber(javaInteger(value)))
}

func (o *JsonObject) AddPropertyDouble(property string, value float64) {
	o.Add(property, NewJsonPrimitiveNumber(javaDouble(value)))
}

func (o *JsonObject) EntrySet() *javautil.LinkedHashMap[string, JsonElement] {
	return o.ensureMap()
}

func (o *JsonObject) Size() int {
	if o == nil || o.members == nil {
		return 0
	}
	return o.members.Size()
}

func (o *JsonObject) Has(memberName string) bool {
	return o != nil && o.members != nil && o.members.ContainsKey(memberName)
}

func (o *JsonObject) Get(memberName string) JsonElement {
	if o == nil || o.members == nil {
		return nil
	}
	v, ok := o.members.Get(memberName)
	if !ok {
		return nil
	}
	return v
}

func (o *JsonObject) Equals(other any) bool {
	if other == any(o) {
		return true
	}
	x, ok := other.(*JsonObject)
	if !ok || x == nil || o == nil {
		return false
	}
	if o.Size() != x.Size() {
		return false
	}
	if o.Size() == 0 {
		return true
	}
	for k, v := range o.members.All() {
		w, exists := x.members.Get(k)
		if !exists || !v.Equals(w) {
			return false
		}
	}
	return true
}

func (o *JsonObject) HashCode() int32 {
	if o == nil || o.members == nil {
		return 0
	}
	var h int32
	for k, v := range o.members.All() {
		h += javalang.StringHashCode(k) ^ v.HashCode()
	}
	return h
}

// ---------------------------------------------------------------------------
// JsonArray
// ---------------------------------------------------------------------------

// JsonArray implements com.google.gson.JsonArray.
type JsonArray struct {
	elements []JsonElement
}

func (*JsonArray) isJsonElement() {}

func NewJsonArray() *JsonArray { return &JsonArray{} }

func (*JsonArray) IsJsonArray() bool     { return true }
func (*JsonArray) IsJsonObject() bool    { return false }
func (*JsonArray) IsJsonPrimitive() bool { return false }
func (*JsonArray) IsJsonNull() bool      { return false }

func (a *JsonArray) GetAsJsonObject() *JsonObject {
	panic(javalang.NewIllegalStateException("Not a JSON Object: " + a.String()))
}
func (a *JsonArray) GetAsJsonArray() *JsonArray { return a }
func (a *JsonArray) GetAsJsonPrimitive() *JsonPrimitive {
	panic(javalang.NewIllegalStateException("Not a JSON Primitive: " + a.String()))
}

func (a *JsonArray) String() string { return jsonElementToString(a) }

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

func (a *JsonArray) Add(element JsonElement) {
	if element == nil {
		element = JsonNull_INSTANCE
	}
	a.elements = append(a.elements, element)
}

func (a *JsonArray) AddString(s string) {
	a.elements = append(a.elements, NewJsonPrimitiveString(s))
}

func (a *JsonArray) Size() int     { return len(a.elements) }
func (a *JsonArray) IsEmpty() bool { return len(a.elements) == 0 }
func (a *JsonArray) All() iter.Seq[JsonElement] {
	return slices.Values(a.elements)
}

func (a *JsonArray) Get(i int) JsonElement {
	if i < 0 || i >= len(a.elements) {
		panic(javalang.NewIndexOutOfBoundsExceptionForIndex(i, len(a.elements)))
	}
	return a.elements[i]
}

func (a *JsonArray) singleElem() JsonElement {
	if len(a.elements) == 1 {
		return a.elements[0]
	}
	panic(javalang.NewIllegalStateException("JsonArray size != 1"))
}

func (a *JsonArray) GetAsBoolean() bool      { return a.singleElem().GetAsBoolean() }
func (a *JsonArray) GetAsNumber() JavaNumber { return a.singleElem().GetAsNumber() }
func (a *JsonArray) GetAsString() string     { return a.singleElem().GetAsString() }
func (a *JsonArray) GetAsLong() int64        { return a.singleElem().GetAsLong() }
func (a *JsonArray) GetAsInt() int32         { return a.singleElem().GetAsInt() }
func (a *JsonArray) GetAsDouble() float64    { return a.singleElem().GetAsDouble() }

func (a *JsonArray) Equals(o any) bool {
	if o == any(a) {
		return true
	}
	x, ok := o.(*JsonArray)
	if !ok || x == nil || a == nil || len(x.elements) != len(a.elements) {
		return false
	}
	for i, e := range a.elements {
		if !x.elements[i].Equals(e) {
			return false
		}
	}
	return true
}

func (a *JsonArray) HashCode() int32 {
	var h int32 = 1
	for _, e := range a.elements {
		h = 31*h + e.HashCode()
	}
	return h
}

// ---------------------------------------------------------------------------
// JsonPrimitive
// ---------------------------------------------------------------------------

// JsonPrimitive implements com.google.gson.JsonPrimitive.
type JsonPrimitive struct {
	value any
}

func (*JsonPrimitive) isJsonElement() {}

func NewJsonPrimitiveBoolean(b bool) *JsonPrimitive { return &JsonPrimitive{value: b} }
func NewJsonPrimitiveString(s string) *JsonPrimitive { return &JsonPrimitive{value: s} }
func NewJsonPrimitiveInt(v int32) *JsonPrimitive {
	return &JsonPrimitive{value: javaInteger(v)}
}
func NewJsonPrimitiveLong(v int64) *JsonPrimitive {
	return &JsonPrimitive{value: javaInteger(v)}
}
func NewJsonPrimitiveDouble(v float64) *JsonPrimitive {
	return &JsonPrimitive{value: javaDouble(v)}
}
func NewJsonPrimitiveNumber(num JavaNumber) *JsonPrimitive {
	if num == nil {
		panic(javalang.NewNullPointerException())
	}
	return &JsonPrimitive{value: num}
}

func (*JsonPrimitive) IsJsonArray() bool     { return false }
func (*JsonPrimitive) IsJsonObject() bool    { return false }
func (*JsonPrimitive) IsJsonPrimitive() bool { return true }
func (*JsonPrimitive) IsJsonNull() bool      { return false }

func (p *JsonPrimitive) GetAsJsonObject() *JsonObject {
	panic(javalang.NewIllegalStateException("Not a JSON Object: " + p.String()))
}
func (p *JsonPrimitive) GetAsJsonArray() *JsonArray {
	panic(javalang.NewIllegalStateException("Not a JSON Array: " + p.String()))
}
func (p *JsonPrimitive) GetAsJsonPrimitive() *JsonPrimitive { return p }
func (p *JsonPrimitive) DeepCopy() JsonElement              { return p }
func (p *JsonPrimitive) String() string                     { return jsonElementToString(p) }

func (p *JsonPrimitive) IsBoolean() bool { _, ok := p.value.(bool); return ok }
func (p *JsonPrimitive) IsNumber() bool  { _, ok := p.value.(JavaNumber); return ok }
func (p *JsonPrimitive) IsString() bool  { _, ok := p.value.(string); return ok }

func (p *JsonPrimitive) GetAsBoolean() bool {
	if b, ok := p.value.(bool); ok {
		return b
	}
	return strings.EqualFold(p.GetAsString(), "true")
}

func (p *JsonPrimitive) GetAsNumber() JavaNumber {
	if s, ok := p.value.(string); ok {
		return NewLazilyParsedNumber(s)
	}
	return p.value.(JavaNumber)
}

func (p *JsonPrimitive) GetAsString() string {
	switch v := p.value.(type) {
	case JavaNumber:
		return v.String()
	case bool:
		if v {
			return "true"
		}
		return "false"
	case string:
		return v
	default:
		return javalang.ToJavaString(v)
	}
}

func (p *JsonPrimitive) GetAsLong() int64 {
	if num, ok := p.value.(JavaNumber); ok {
		return num.LongValue()
	}
	return javalang.LongParseLong(p.GetAsString())
}

func (p *JsonPrimitive) GetAsInt() int32 {
	if num, ok := p.value.(JavaNumber); ok {
		return num.IntValue()
	}
	return javalang.IntegerParseInt(p.GetAsString())
}

func (p *JsonPrimitive) GetAsDouble() float64 {
	if num, ok := p.value.(JavaNumber); ok {
		return num.DoubleValue()
	}
	return javalang.DoubleParseDouble(p.GetAsString())
}

func isIntegralPrimitive(p *JsonPrimitive) bool {
	_, ok := p.value.(javaInteger)
	return ok
}

func (p *JsonPrimitive) HashCode() int32 {
	if p.value == nil {
		return 31
	}
	if isIntegralPrimitive(p) {
		v := p.GetAsNumber().LongValue()
		return int32(v ^ javalang.Ushr64(v, 32))
	}
	if _, ok := p.value.(JavaNumber); ok {
		bits := int64(math.Float64bits(p.GetAsNumber().DoubleValue()))
		return int32(bits ^ javalang.Ushr64(bits, 32))
	}
	return javalang.ObjectHashCode(p.value)
}

func (p *JsonPrimitive) Equals(obj any) bool {
	if obj == any(p) {
		return true
	}
	other, ok := obj.(*JsonPrimitive)
	if !ok || other == nil || p == nil {
		return false
	}
	if p.value == nil {
		return other.value == nil
	}
	if isIntegralPrimitive(p) && isIntegralPrimitive(other) {
		return p.GetAsNumber().LongValue() == other.GetAsNumber().LongValue()
	}
	_, numA := p.value.(JavaNumber)
	_, numB := other.value.(JavaNumber)
	if numA && numB {
		a := p.GetAsNumber().DoubleValue()
		b := other.GetAsNumber().DoubleValue()
		return a == b || (math.IsNaN(a) && math.IsNaN(b))
	}
	if numA != numB {
		return false
	}
	return p.value == other.value
}

// ---------------------------------------------------------------------------
// JsonNull
// ---------------------------------------------------------------------------

// JsonNull implements com.google.gson.JsonNull.
type JsonNull struct{}

func (*JsonNull) isJsonElement() {}

var JsonNull_INSTANCE = &JsonNull{}

func (*JsonNull) IsJsonArray() bool     { return false }
func (*JsonNull) IsJsonObject() bool    { return false }
func (*JsonNull) IsJsonPrimitive() bool { return false }
func (*JsonNull) IsJsonNull() bool      { return true }

func (n *JsonNull) GetAsJsonObject() *JsonObject {
	panic(javalang.NewIllegalStateException("Not a JSON Object: " + n.String()))
}
func (n *JsonNull) GetAsJsonArray() *JsonArray {
	panic(javalang.NewIllegalStateException("Not a JSON Array: " + n.String()))
}
func (n *JsonNull) GetAsJsonPrimitive() *JsonPrimitive {
	panic(javalang.NewIllegalStateException("Not a JSON Primitive: " + n.String()))
}

func (*JsonNull) GetAsBoolean() bool      { panic(javalang.NewUnsupportedOperationException("JsonNull")) }
func (*JsonNull) GetAsNumber() JavaNumber { panic(javalang.NewUnsupportedOperationException("JsonNull")) }
func (*JsonNull) GetAsString() string     { panic(javalang.NewUnsupportedOperationException("JsonNull")) }
func (*JsonNull) GetAsLong() int64        { panic(javalang.NewUnsupportedOperationException("JsonNull")) }
func (*JsonNull) GetAsInt() int32         { panic(javalang.NewUnsupportedOperationException("JsonNull")) }
func (*JsonNull) GetAsDouble() float64    { panic(javalang.NewUnsupportedOperationException("JsonNull")) }

func (n *JsonNull) String() string        { return "null" }
func (n *JsonNull) DeepCopy() JsonElement { return JsonNull_INSTANCE }
func (n *JsonNull) HashCode() int32       { return javalang.StringHashCode("com.google.gson.JsonNull") }
func (n *JsonNull) Equals(other any) bool {
	_, ok := other.(*JsonNull)
	return ok && other != nil
}
