package jutil

import (
	"math"
	"reflect"
	"strconv"
	"strings"
)

// Go-only: language-adapter — java.lang.Object equality/hashing,
// java.util.Objects and Guava MoreObjects.

// hashable is implemented by ported classes that override Java's
// Object.equals(Object) and Object.hashCode(). Code that compares or hashes
// arbitrary values (collections, Objects.equals) dispatches to it.
type hashable interface {
	Equals(o any) bool
	HashCode() int32
}

// Equaler is implemented by ported classes that override Object.equals but
// whose hashCode is never observed.
type Equaler interface {
	Equals(o any) bool
}

// ObjectsEquals is java.util.Objects.equals(Object, Object): null-safe, and
// dispatching to Equals(any) when the value overrides equals. Values without
// an Equals method compare like their Java counterparts: strings, booleans
// and integers by value (boxed types of different kinds are never equal),
// doubles by bit pattern (Double.equals: NaN equals NaN, 0.0 != -0.0), and
// everything else (pointers, enums) by identity.
func ObjectsEquals(a, b any) bool {
	if a == nil {
		return IsNil(b)
	}
	if b == nil {
		return IsNil(a)
	}
	switch x := a.(type) {
	case string:
		y, ok := b.(string)
		return ok && x == y
	case int:
		y, ok := b.(int)
		return ok && x == y
	case int32:
		y, ok := b.(int32)
		return ok && x == y
	case int64:
		y, ok := b.(int64)
		return ok && x == y
	case int16:
		y, ok := b.(int16)
		return ok && x == y
	case int8:
		y, ok := b.(int8)
		return ok && x == y
	case uint:
		y, ok := b.(uint)
		return ok && x == y
	case uint64:
		y, ok := b.(uint64)
		return ok && x == y
	case uint32:
		y, ok := b.(uint32)
		return ok && x == y
	case uint16:
		y, ok := b.(uint16)
		return ok && x == y
	case uint8:
		y, ok := b.(uint8)
		return ok && x == y
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case float64:
		y, ok := b.(float64)
		return ok && doubleBits(x) == doubleBits(y)
	case float32:
		y, ok := b.(float32)
		return ok && floatBits(x) == floatBits(y)
	case *string:
		if x == nil {
			return IsNil(b)
		}
		y, ok := b.(*string)
		return ok && y != nil && *x == *y
	}
	va := reflect.ValueOf(a)
	switch va.Kind() {
	case reflect.Pointer:
		if a == b {
			return true
		}
		if va.IsNil() {
			return IsNil(b)
		}
		if IsNil(b) {
			return false
		}
		if e, ok := a.(Equaler); ok {
			return e.Equals(b)
		}
		return false
	case reflect.Uint32:
		if va.Uint() == 0 {
			if _, ok := a.(interface{ IsNil() bool }); ok {
				return IsNil(b)
			}
		}
		if a == b {
			return true
		}
		if IsNil(b) {
			return false
		}
		if e, ok := a.(Equaler); ok {
			return e.Equals(b)
		}
		return false
	}
	aNil, bNil := IsNil(a), IsNil(b)
	if aNil || bNil {
		return aNil && bNil
	}
	if e, ok := a.(Equaler); ok {
		return e.Equals(b)
	}
	ta := va.Type()
	tb := reflect.TypeOf(b)
	if ta != tb {
		return false
	}
	if !ta.Comparable() {
		// Arrays and other reference types without equals(): identity.
		vb := reflect.ValueOf(b)
		switch va.Kind() {
		case reflect.Slice:
			return va.Len() == vb.Len() && (va.Len() == 0 || va.Pointer() == vb.Pointer())
		case reflect.Map, reflect.Func:
			return va.Pointer() == vb.Pointer()
		}
		return false
	}
	return a == b
}

// ObjectsEqualsOf is ObjectsEquals for two values of a type parameter.
//
// Go-only: performance — in generic code, scalars, pointers, structs and node
// IDs avoid boxing/reflection fast-paths while matching ObjectsEquals exactly.
// Generic collections (HamtPMap, ObjectLinkedHashMap) use it.
func ObjectsEqualsOf[T any](a, b T) bool {
	t := reflect.TypeFor[T]()
	switch t.Kind() {
	case reflect.Pointer:
		if any(a) == any(b) {
			return true
		}
		var zero T
		if any(a) == any(zero) || any(b) == any(zero) {
			return false
		}
		if x, ok := any(a).(*string); ok {
			return *x == *any(b).(*string)
		}
		if e, ok := any(a).(Equaler); ok {
			return e.Equals(b)
		}
		return false
	case reflect.Interface:
		return ObjectsEquals(a, b)
	case reflect.String:
		if x, ok := any(a).(string); ok {
			return any(b) == x
		}
		if e, ok := any(a).(Equaler); ok {
			return e.Equals(b)
		}
		return any(a) == any(b)
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8,
		reflect.Uint, reflect.Uint64, reflect.Uint16, reflect.Uint8, reflect.Bool:
		if t.PkgPath() == "" {
			return any(a) == any(b)
		}
		if e, ok := any(a).(Equaler); ok {
			return e.Equals(b)
		}
		return any(a) == any(b)
	case reflect.Float64:
		if x, ok := any(a).(float64); ok {
			return doubleBits(x) == doubleBits(any(b).(float64))
		}
		if e, ok := any(a).(Equaler); ok {
			return e.Equals(b)
		}
		return any(a) == any(b)
	case reflect.Float32:
		if x, ok := any(a).(float32); ok {
			return floatBits(x) == floatBits(any(b).(float32))
		}
		if e, ok := any(a).(Equaler); ok {
			return e.Equals(b)
		}
		return any(a) == any(b)
	case reflect.Struct:
		if e, ok := any(a).(Equaler); ok {
			return e.Equals(b)
		}
		if t.Comparable() {
			return any(a) == any(b)
		}
		return false
	case reflect.Uint32:
		if any(a) == any(b) {
			return true
		}
		if t.PkgPath() == "" {
			return false
		}
		var zero T
		if any(a) == any(zero) || any(b) == any(zero) {
			if _, ok := any(a).(interface{ IsNil() bool }); ok {
				return false
			}
		}
		if e, ok := any(a).(Equaler); ok {
			return e.Equals(b)
		}
		return false
	}
	return ObjectsEquals(a, b)
}

// objectHashCodeOf is ObjectHashCode for a value of a type parameter.
//
// Go-only: performance — see ObjectsEqualsOf.
func objectHashCodeOf[T any](v T) int32 {
	t := reflect.TypeFor[T]()
	switch t.Kind() {
	case reflect.Pointer:
		var zero T
		if any(v) == any(zero) {
			return 0
		}
		if s, ok := any(v).(*string); ok {
			return StringHashCode(*s)
		}
		if h, ok := any(v).(hashable); ok {
			return h.HashCode()
		}
		p := uint64(reflect.ValueOf(v).Pointer())
		return int32(p ^ (p >> 32))
	case reflect.Interface:
		return ObjectHashCode(v)
	case reflect.String:
		if s, ok := any(v).(string); ok {
			return StringHashCode(s)
		}
		if h, ok := any(v).(hashable); ok {
			return h.HashCode()
		}
		return StringHashCode(reflect.ValueOf(v).String())
	case reflect.Int:
		if x, ok := any(v).(int); ok {
			return int32(x)
		}
		if h, ok := any(v).(hashable); ok {
			return h.HashCode()
		}
		return int32(reflect.ValueOf(v).Int())
	case reflect.Int32:
		if x, ok := any(v).(int32); ok {
			return x
		}
		if h, ok := any(v).(hashable); ok {
			return h.HashCode()
		}
		return int32(reflect.ValueOf(v).Int())
	case reflect.Int64:
		if x, ok := any(v).(int64); ok {
			return int32(x ^ int64(uint64(x)>>32))
		}
		if h, ok := any(v).(hashable); ok {
			return h.HashCode()
		}
		x := reflect.ValueOf(v).Int()
		return int32(x ^ int64(uint64(x)>>32))
	case reflect.Int16, reflect.Int8:
		if t.PkgPath() == "" {
			if x, ok := any(v).(int16); ok {
				return int32(x)
			}
			return int32(any(v).(int8))
		}
		if h, ok := any(v).(hashable); ok {
			return h.HashCode()
		}
		return int32(reflect.ValueOf(v).Int())
	case reflect.Uint, reflect.Uint64:
		if t.PkgPath() == "" {
			if x, ok := any(v).(uint64); ok {
				return int32(x ^ (x >> 32))
			}
			u := uint64(any(v).(uint))
			return int32(u ^ (u >> 32))
		}
		if h, ok := any(v).(hashable); ok {
			return h.HashCode()
		}
		x := reflect.ValueOf(v).Uint()
		return int32(x ^ (x >> 32))
	case reflect.Uint16, reflect.Uint8:
		if t.PkgPath() == "" {
			if x, ok := any(v).(uint8); ok {
				return int32(x)
			}
			return int32(any(v).(uint16))
		}
		if h, ok := any(v).(hashable); ok {
			return h.HashCode()
		}
		return int32(reflect.ValueOf(v).Uint())
	case reflect.Bool:
		if x, ok := any(v).(bool); ok {
			if x {
				return 1231
			}
			return 1237
		}
		if h, ok := any(v).(hashable); ok {
			return h.HashCode()
		}
		if reflect.ValueOf(v).Bool() {
			return 1231
		}
		return 1237
	case reflect.Float64:
		if x, ok := any(v).(float64); ok {
			bits := doubleBits(x)
			return int32(bits ^ (bits >> 32))
		}
	case reflect.Float32:
		if x, ok := any(v).(float32); ok {
			return int32(floatBits(x))
		}
	case reflect.Struct:
		if h, ok := any(v).(hashable); ok {
			return h.HashCode()
		}
		return 0
	case reflect.Uint32:
		if x, ok := any(v).(uint32); ok {
			return int32(x)
		}
		var zero T
		if any(v) == any(zero) {
			if _, ok := any(v).(interface{ IsNil() bool }); ok {
				return 0
			}
		}
		if h, ok := any(v).(hashable); ok {
			return h.HashCode()
		}
		return int32(reflect.ValueOf(v).Uint())
	}
	return ObjectHashCode(v)
}

// doubleBits is Double.doubleToLongBits (all NaNs collapse to one pattern).
func doubleBits(d float64) uint64 {
	if math.IsNaN(d) {
		return 0x7ff8000000000000
	}
	return math.Float64bits(d)
}

func floatBits(f float32) uint32 {
	if f != f {
		return 0x7fc00000
	}
	return math.Float32bits(f)
}

// ObjectHashCode is java.util.Objects.hashCode(Object): 0 for null, the
// value's HashCode() if it overrides hashCode, and otherwise the hashCode of
// the corresponding Java boxed type (String, Boolean, Integer, Long, Double,
// Character), or an identity hash for reference types.
func ObjectHashCode(v any) int32 {
	if v == nil {
		return 0
	}
	switch x := v.(type) {
	case string:
		return StringHashCode(x)
	case *string:
		if x == nil {
			return 0
		}
		return StringHashCode(*x)
	case bool:
		if x {
			return 1231
		}
		return 1237
	case int:
		// Go int stands for Java int.
		return int32(x)
	case int32:
		return x
	case int16:
		return int32(x)
	case int8:
		return int32(x)
	case uint16:
		return int32(x)
	case uint32:
		return int32(x)
	case uint64:
		return int32(x ^ (x >> 32))
	case uint:
		return int32(uint64(x) ^ (uint64(x) >> 32))
	case uint8:
		return int32(x)
	case int64:
		return int32(x ^ int64(uint64(x)>>32))
	case float64:
		bits := doubleBits(x)
		return int32(bits ^ (bits >> 32))
	case float32:
		return int32(floatBits(x))
	}
	if IsNil(v) {
		return 0
	}
	if h, ok := v.(hashable); ok {
		return h.HashCode()
	}
	return identityHashCodeNonNil(v)
}

// IdentityHashCode is System.identityHashCode(Object). Java's identity hash
// is arbitrary (and differs from run to run); here it is derived from the
// object's address, which is stable for the object's lifetime. Values that
// are not references (Go value types standing in for Java enums or value
// classes) hash by value.
func IdentityHashCode(v any) int32 {
	if IsNil(v) {
		return 0
	}
	return identityHashCodeNonNil(v)
}

func identityHashCodeNonNil(v any) int32 {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan, reflect.UnsafePointer:
		p := uint64(rv.Pointer())
		return int32(p ^ (p >> 32))
	case reflect.Int64:
		x := rv.Int()
		return int32(x ^ int64(uint64(x)>>32))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return int32(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		x := rv.Uint()
		return int32(x ^ (x >> 32))
	case reflect.String:
		return StringHashCode(rv.String())
	case reflect.Bool:
		if rv.Bool() {
			return 1231
		}
		return 1237
	}
	return 0
}

// IdentityHashCodeOfPointer is IdentityHashCode(p) for a typed pointer.
func IdentityHashCodeOfPointer[T any](p *T) int32 {
	if p == nil {
		return 0
	}
	x := uint64(reflect.ValueOf(p).Pointer())
	return int32(x ^ (x >> 32))
}

// ObjectToString is java.lang.Object.toString() for classes that do not
// override it: "<fully qualified class name>@<hex identity hash>".
func ObjectToString(className string, v any) string {
	return className + "@" + strconv.FormatUint(uint64(uint32(IdentityHashCode(v))), 16)
}

// IntegerToHexString is Integer.toHexString(int): the unsigned hex digits.
func IntegerToHexString(i int32) string {
	return strconv.FormatUint(uint64(uint32(i)), 16)
}

// LongToHexString is Long.toHexString(long): the unsigned hex digits.
func LongToHexString(i int64) string {
	return strconv.FormatUint(uint64(i), 16)
}

// StringPtr returns a pointer to a copy of s. Java String values that may be
// null are *string in the port; this builds a non-null one from a literal or
// expression.
func StringPtr(s string) *string { return &s }

// ---------------------------------------------------------------------------
// MoreObjects.ToStringHelper

type toStringHolder struct {
	name          string
	hasName       bool
	value         any
	unconditional bool
}

// ToStringHelper is Guava's MoreObjects.ToStringHelper.
type ToStringHelper struct {
	className      string
	holders        []toStringHolder
	omitNullValues bool
}

// MoreObjects_ToStringHelper is MoreObjects.toStringHelper(Object self) /
// toStringHelper(String className); pass the simple class name.
func MoreObjects_ToStringHelper(className string) *ToStringHelper {
	return &ToStringHelper{className: className}
}

// OmitNullValues is ToStringHelper.omitNullValues().
func (h *ToStringHelper) OmitNullValues() *ToStringHelper {
	h.omitNullValues = true
	return h
}

// Add is ToStringHelper.add(String, Object).
func (h *ToStringHelper) Add(name string, value any) *ToStringHelper {
	h.holders = append(h.holders, toStringHolder{name: name, hasName: true, value: value})
	return h
}

// AddBool is ToStringHelper.add(String, boolean), which (like the other
// primitive overloads) is never omitted.
func (h *ToStringHelper) AddBool(name string, value bool) *ToStringHelper {
	s := "false"
	if value {
		s = "true"
	}
	h.holders = append(h.holders, toStringHolder{name: name, hasName: true, value: s, unconditional: true})
	return h
}

// AddValue is ToStringHelper.addValue(Object).
func (h *ToStringHelper) AddValue(value any) *ToStringHelper {
	h.holders = append(h.holders, toStringHolder{value: value})
	return h
}

// String is ToStringHelper.toString().
func (h *ToStringHelper) String() string {
	var sb strings.Builder
	sb.WriteString(h.className)
	sb.WriteByte('{')
	sep := ""
	for _, vh := range h.holders {
		v := vh.value
		if !vh.unconditional && h.omitNullValues && IsNil(v) {
			continue
		}
		sb.WriteString(sep)
		sep = ", "
		if vh.hasName {
			sb.WriteString(vh.name)
			sb.WriteByte('=')
		}
		sb.WriteString(ToJavaString(v))
	}
	sb.WriteByte('}')
	return sb.String()
}
