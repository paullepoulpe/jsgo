package jutil

import (
	"math"
	"reflect"
	"strconv"
	"strings"
)

// Go-only: language-adapter — java.lang.Object equality/hashing,
// java.util.Objects and Guava MoreObjects.

// Hashable is implemented by ported classes that override Java's
// Object.equals(Object) and Object.hashCode(). Code that compares or hashes
// arbitrary values (collections, Objects.equals) dispatches to it.
type Hashable interface {
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
	aNil, bNil := IsNil(a), IsNil(b)
	if aNil || bNil {
		return aNil && bNil
	}
	if e, ok := a.(Equaler); ok {
		return e.Equals(b)
	}
	switch x := a.(type) {
	case float64:
		y, ok := b.(float64)
		return ok && doubleBits(x) == doubleBits(y)
	case float32:
		y, ok := b.(float32)
		return ok && floatBits(x) == floatBits(y)
	}
	ta, tb := reflect.TypeOf(a), reflect.TypeOf(b)
	if ta != tb {
		return false
	}
	if !ta.Comparable() {
		// Arrays and other reference types without equals(): identity.
		va, vb := reflect.ValueOf(a), reflect.ValueOf(b)
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
	if IsNil(v) {
		return 0
	}
	if h, ok := v.(Hashable); ok {
		return h.HashCode()
	}
	switch x := v.(type) {
	case string:
		return StringHashCode(x)
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
	case int64:
		return int32(x ^ int64(uint64(x)>>32))
	case float64:
		bits := doubleBits(x)
		return int32(bits ^ (bits >> 32))
	case float32:
		return int32(floatBits(x))
	}
	return IdentityHashCode(v)
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
	omitEmpty      bool
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

// OmitEmptyValues is ToStringHelper.omitEmptyValues().
func (h *ToStringHelper) OmitEmptyValues() *ToStringHelper {
	h.omitEmpty = true
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
	h.holders = append(h.holders, toStringHolder{name: name, hasName: true, value: ToJavaString(value), unconditional: true})
	return h
}

// AddValue is ToStringHelper.addValue(Object).
func (h *ToStringHelper) AddValue(value any) *ToStringHelper {
	h.holders = append(h.holders, toStringHolder{value: value})
	return h
}

func toStringHelperIsEmpty(v any) bool {
	switch x := v.(type) {
	case string:
		return x == ""
	case JavaCollection:
		return x.IsEmpty()
	case JavaMap:
		return x.IsEmpty()
	}
	return false
}

// String is ToStringHelper.toString().
func (h *ToStringHelper) String() string {
	var sb strings.Builder
	sb.WriteString(h.className)
	sb.WriteByte('{')
	sep := ""
	for _, vh := range h.holders {
		v := vh.value
		if vh.unconditional || (func() bool {
			if IsNil(v) {
				return !h.omitNullValues
			}
			return !h.omitEmpty || !toStringHelperIsEmpty(v)
		})() {
			sb.WriteString(sep)
			sep = ", "
			if vh.hasName {
				sb.WriteString(vh.name)
				sb.WriteByte('=')
			}
			sb.WriteString(ToJavaString(v))
		}
	}
	sb.WriteByte('}')
	return sb.String()
}
