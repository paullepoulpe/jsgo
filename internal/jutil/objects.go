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
	if na, ok := a.(interface{ IsNil() bool }); ok {
		if na.IsNil() {
			return IsNil(b)
		}
	}
	va := reflect.ValueOf(a)
	if va.Kind() == reflect.Pointer {
		if va.IsNil() {
			return IsNil(b)
		}
		if a == b {
			return true
		}
		e, ok := a.(Equaler)
		if !ok || IsNil(b) {
			return false
		}
		return e.Equals(b)
	}
	if reflectKindIsNil(va) {
		return IsNil(b)
	}
	if IsNil(b) {
		return false
	}
	if e, ok := a.(Equaler); ok {
		return e.Equals(b)
	}
	vb := reflect.ValueOf(b)
	ta := va.Type()
	if ta != vb.Type() {
		return false
	}
	if !ta.Comparable() {
		// Arrays and other reference types without equals(): identity.
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

func fold64(u uint64) int32 {
	return int32(u ^ (u >> 32))
}

func boolHash(b bool) int32 {
	if b {
		return 1231
	}
	return 1237
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
		return boolHash(x)
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
		return fold64(x)
	case uint:
		return fold64(uint64(x))
	case uint8:
		return int32(x)
	case int64:
		return fold64(uint64(x))
	case float64:
		return fold64(doubleBits(x))
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
		return fold64(uint64(rv.Pointer()))
	case reflect.Int64:
		return fold64(uint64(rv.Int()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return int32(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fold64(rv.Uint())
	case reflect.String:
		return StringHashCode(rv.String())
	case reflect.Bool:
		return boolHash(rv.Bool())
	}
	return 0
}

// IdentityHashCodeOfPointer is IdentityHashCode(p) for a typed pointer.
func IdentityHashCodeOfPointer[T any](p *T) int32 {
	if p == nil {
		return 0
	}
	return fold64(uint64(reflect.ValueOf(p).Pointer()))
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
