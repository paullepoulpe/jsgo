// Clean-room specification based on Java Language Specification & java.lang.Object Javadoc.

package javalang

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

// NullString is a WTF-8/UTF-8 impossible byte sentinel representing Java `(String) null`
// inside a value-typed Go `string` with zero heap allocation.
const NullString = "\xff\x00j2go_null"

// Equaler is implemented by types overriding Object.equals(Object).
type Equaler interface {
	Equals(other any) bool
}

// HashCoder is implemented by types overriding Object.hashCode().
type HashCoder interface {
	HashCode() int32
}

// Hashable is implemented by types overriding Object.equals and Object.hashCode.
type Hashable interface {
	Equaler
	HashCode() int32
}

// JavaNilChecker is implemented by generated structs and runtime wrappers for reflection-free null checks.
type JavaNilChecker interface {
	IsJavaNil_() bool
}

// IdentityHasher is implemented by generated structs that cache or compute their identity hash without reflection.
type IdentityHasher interface {
	IdentityHash_() int32
}

// IdentityEqualer is implemented by generated classes inheriting default Object.equals pointer/value identity.
type IdentityEqualer interface {
	HasIdentityEquals_()
	IsJavaNil_() bool
}

// ObjectsEquals reproduces java.util.Objects.equals(a, b).
// It handles nils, typed nil interfaces, Equaler dispatch, primitive comparisons,
// pointer reference identity, and double IEEE-754 bit-pattern equality (where NaN equals NaN).
func ObjectsEquals(a, b any) bool {
	if a == nil {
		return IsNil(b)
	}
	if sa, ok := a.(string); ok {
		if sb, ok := b.(string); ok {
			return sa == sb
		}
		if sa == NullString {
			return IsNil(b)
		}
		return false
	}
	if eq, ok := a.(Equaler); ok {
		if nc, ok2 := a.(JavaNilChecker); ok2 {
			if nc.IsJavaNil_() {
				return IsNil(b)
			}
		} else if IsNil(a) {
			return IsNil(b)
		}
		if b == nil {
			return false
		}
		if ncB, ok2 := b.(JavaNilChecker); ok2 {
			if ncB.IsJavaNil_() {
				return false
			}
		} else if sB, ok2 := b.(string); ok2 {
			if sB == NullString {
				return false
			}
		} else if IsNil(b) {
			return false
		}
		return eq.Equals(b)
	}
	if ia, ok := a.(IdentityEqualer); ok {
		if ia.IsJavaNil_() {
			return IsNil(b)
		}
		return a == b
	}
	if b == nil {
		return IsNil(a)
	}
	switch x := a.(type) {
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
	}
	if IsNil(a) {
		return IsNil(b)
	}
	if IsNil(b) {
		return false
	}
	if eq, ok := b.(Equaler); ok {
		return eq.Equals(a)
	}
	if _, ok := a.(interface{ HasIdentityEquals_() }); ok {
		return a == b
	}
	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)
	if ta != tb {
		return false
	}
	if ta.Kind() == reflect.Pointer {
		return a == b
	}
	if ta.Comparable() && !typeMayContainInterface(ta) {
		return a == b
	}
	if ta.Comparable() && reflect.ValueOf(a).Comparable() && reflect.ValueOf(b).Comparable() {
		return a == b
	}
	return reflect.DeepEqual(a, b)
}

func typeMayContainInterface(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Interface:
		return true
	case reflect.Struct:
		for i, n := 0, t.NumField(); i < n; i++ {
			if typeMayContainInterface(t.Field(i).Type) {
				return true
			}
		}
		return false
	case reflect.Array:
		return typeMayContainInterface(t.Elem())
	default:
		return false
	}
}

func doubleBits(d float64) uint64 {
	if math.IsNaN(d) {
		return 0x7ff8000000000000
	}
	return math.Float64bits(d)
}

func floatBits(f float32) uint32 {
	if math.IsNaN(float64(f)) {
		return 0x7fc00000
	}
	return math.Float32bits(f)
}

func fold64(u uint64) int32 {
	return int32(u ^ (u >> 32))
}

// ObjectHashCode reproduces java.util.Objects.hashCode(o).
func ObjectHashCode(v any) int32 {
	if v == nil {
		return 0
	}
	if s, ok := v.(string); ok {
		if s == NullString {
			return 0
		}
		return StringHashCode(s)
	}
	if h, ok := v.(HashCoder); ok {
		if nc, ok2 := v.(JavaNilChecker); ok2 {
			if nc.IsJavaNil_() {
				return 0
			}
		} else if IsNil(v) {
			return 0
		}
		return h.HashCode()
	}
	if ih, ok := v.(IdentityHasher); ok {
		if nc, ok2 := v.(JavaNilChecker); ok2 {
			if nc.IsJavaNil_() {
				return 0
			}
		} else if IsNil(v) {
			return 0
		}
		return ih.IdentityHash_()
	}
	switch x := v.(type) {
	case bool:
		if x {
			return 1231
		}
		return 1237
	case int:
		return int32(x)
	case int32:
		return x
	case int64:
		return fold64(uint64(x))
	case int16:
		return int32(x)
	case int8:
		return int32(x)
	case uint32:
		return int32(x)
	case uint64:
		return fold64(x)
	case uint:
		return fold64(uint64(x))
	case uint16:
		return int32(x)
	case uint8:
		return int32(x)
	case float64:
		return fold64(doubleBits(x))
	case float32:
		return int32(floatBits(x))
	}
	if IsNil(v) {
		return 0
	}
	return IdentityHashCode(v)
}

// IdentityHashCode reproduces System.identityHashCode(o).
func IdentityHashCode(v any) int32 {
	if v == nil {
		return 0
	}
	if ih, ok := v.(IdentityHasher); ok {
		if nc, ok2 := v.(JavaNilChecker); ok2 {
			if nc.IsJavaNil_() {
				return 0
			}
		} else if IsNil(v) {
			return 0
		}
		return ih.IdentityHash_()
	}
	if s, ok := v.(string); ok {
		if s == NullString {
			return 0
		}
		return StringHashCode(s)
	}
	switch x := v.(type) {
	case int:
		return int32(x)
	case int32:
		return x
	case int64:
		return fold64(uint64(x))
	case int16:
		return int32(x)
	case int8:
		return int32(x)
	case uint:
		return fold64(uint64(x))
	case uint64:
		return fold64(x)
	case uint32:
		return int32(x)
	case uint16:
		return int32(x)
	case uint8:
		return int32(x)
	case bool:
		if x {
			return 1231
		}
		return 1237
	case float64:
		return fold64(doubleBits(x))
	case float32:
		return int32(floatBits(x))
	}
	if nc, ok := v.(JavaNilChecker); ok && nc.IsJavaNil_() {
		return 0
	}
	if n, ok := v.(interface{ IsNil() bool }); ok && n.IsNil() {
		return 0
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan, reflect.UnsafePointer:
		if rv.IsNil() {
			return 0
		}
		return fold64(uint64(rv.Pointer()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return int32(rv.Int())
	case reflect.Int64:
		return fold64(uint64(rv.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fold64(rv.Uint())
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

// ObjectToString formats the default Java Object.toString: "ClassName@hex(hash)".
func ObjectToString(className string, v any) string {
	return className + "@" + strconv.FormatUint(uint64(uint32(IdentityHashCode(v))), 16)
}

// IsNil reports whether an interface value is untyped nil or contains a typed nil pointer/slice/map.
func IsNil(v any) bool {
	if v == nil {
		return true
	}
	if nc, ok := v.(JavaNilChecker); ok {
		return nc.IsJavaNil_()
	}
	if s, ok := v.(string); ok {
		return s == NullString
	}
	switch x := v.(type) {
	case int, int32, int64, int16, int8, uint, uint64, uint32, uint16, uint8, bool, float64, float32:
		return false
	case []string:
		return x == nil
	case []byte:
		return x == nil
	case []int:
		return x == nil
	case []int32:
		return x == nil
	case []int64:
		return x == nil
	case []float64:
		return x == nil
	case []bool:
		return x == nil
	case []uint16:
		return x == nil
	case []any:
		return x == nil
	case *string:
		return x == nil
	case *int:
		return x == nil
	case *int32:
		return x == nil
	case *int64:
		return x == nil
	case *bool:
		return x == nil
	case *float64:
		return x == nil
	case *float32:
		return x == nil
	}
	if n, ok := v.(interface{ IsNil() bool }); ok {
		return n.IsNil()
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan, reflect.Interface, reflect.UnsafePointer:
		return rv.IsNil()
	}
	return false
}

// ToJavaString converts any value to its Java string representation.
func ToJavaString(v any) string {
	if v == nil {
		return "null"
	}
	switch x := v.(type) {
	case string:
		if x == NullString {
			return "null"
		}
		return x
	case float64:
		return DoubleToString(x)
	case float32:
		return FloatToString(x)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case int:
		return strconv.Itoa(x)
	case bool:
		if x {
			return "true"
		}
		return "false"
	case uint16:
		return StringFromChar(x)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case uint8:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	}
	if IsNil(v) {
		return "null"
	}
	if str, ok := v.(fmt.Stringer); ok {
		return str.String()
	}
	return fmt.Sprint(v)
}
