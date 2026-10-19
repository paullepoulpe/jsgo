package jutil

import (
	"reflect"
	"strings"
	"unsafe"
)

// Go-only: language-adapter — Guava com.google.common.base.Preconditions.

// CheckState is Preconditions.checkState(boolean).
func CheckState(expression bool) {
	if !expression {
		panic(NewIllegalStateException())
	}
}

// CheckStateMsg is Preconditions.checkState(boolean, String, Object...).
func CheckStateMsg(expression bool, template string, args ...any) {
	if !expression {
		panic(NewIllegalStateExceptionMsg(LenientFormat(template, args...)))
	}
}

// Fixed-arity variants (CheckStateMsg1..3, CheckArgumentMsg1..3,
// CheckNotNullMsg1..2) mirror Guava's own non-varargs overloads such as
// checkArgument(boolean, String, Object p1, Object p2), which exist to avoid
// allocating on the success path. In Go the variadic forms box every argument
// into an interface on every call, even when the check passes, so ported code
// uses the fixed-arity form whenever there are 1-3 template arguments
// (PORTING.md §5).

// CheckStateMsg1 is Preconditions.checkState(boolean, String, Object).
func CheckStateMsg1[A any](expression bool, template string, p1 A) {
	if !expression {
		failState(template, p1)
	}
}

// CheckStateMsg2 is Preconditions.checkState(boolean, String, Object, Object).
func CheckStateMsg2[A, B any](expression bool, template string, p1 A, p2 B) {
	if !expression {
		failState(template, p1, p2)
	}
}

// CheckStateMsg3 is Preconditions.checkState(boolean, String, Object, Object, Object).
func CheckStateMsg3[A, B, C any](expression bool, template string, p1 A, p2 B, p3 C) {
	if !expression {
		failState(template, p1, p2, p3)
	}
}

//go:noinline
func failState(template string, args ...any) {
	panic(NewIllegalStateExceptionMsg(LenientFormat(template, args...)))
}

// CheckArgument is Preconditions.checkArgument(boolean).
func CheckArgument(expression bool) {
	if !expression {
		panic(NewIllegalArgumentException())
	}
}

// CheckArgumentMsg is Preconditions.checkArgument(boolean, String, Object...).
func CheckArgumentMsg(expression bool, template string, args ...any) {
	if !expression {
		panic(NewIllegalArgumentExceptionMsg(LenientFormat(template, args...)))
	}
}

// CheckArgumentMsg1 is Preconditions.checkArgument(boolean, String, Object).
func CheckArgumentMsg1[A any](expression bool, template string, p1 A) {
	if !expression {
		failArgument(template, p1)
	}
}

// CheckArgumentMsg2 is Preconditions.checkArgument(boolean, String, Object, Object).
func CheckArgumentMsg2[A, B any](expression bool, template string, p1 A, p2 B) {
	if !expression {
		failArgument(template, p1, p2)
	}
}

// CheckArgumentMsg3 is Preconditions.checkArgument(boolean, String, Object, Object, Object).
func CheckArgumentMsg3[A, B, C any](expression bool, template string, p1 A, p2 B, p3 C) {
	if !expression {
		failArgument(template, p1, p2, p3)
	}
}

//go:noinline
func failArgument(template string, args ...any) {
	panic(NewIllegalArgumentExceptionMsg(LenientFormat(template, args...)))
}

// IsNil reports whether v is Java null: a nil interface, or a nil pointer,
// map, slice, func or chan stored in an interface.
func IsNil(v any) bool {
	if v == nil {
		return true
	}
	if (*[2]unsafe.Pointer)(unsafe.Pointer(&v))[1] != nil && reflect.ValueOf(v).Kind() != reflect.Slice {
		return false
	}
	return isNilSlow(v)
}

//go:noinline
func isNilSlow(v any) bool {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan, reflect.Interface:
		return rv.IsNil()
	}
	return false
}

// CheckNotNull is Preconditions.checkNotNull(T).
func CheckNotNull[T any](reference T) T {
	if IsNil(reference) {
		failNull()
	}
	return reference
}

//go:noinline
func failNull() {
	panic(NewNullPointerException())
}

// CheckNotNullMsg is Preconditions.checkNotNull(T, String, Object...).
func CheckNotNullMsg[T any](reference T, template string, args ...any) T {
	if IsNil(reference) {
		failNotNull(template, args...)
	}
	return reference
}

// CheckNotNullMsg1 is Preconditions.checkNotNull(T, String, Object).
func CheckNotNullMsg1[T, A any](reference T, template string, p1 A) T {
	if IsNil(reference) {
		failNotNull(template, p1)
	}
	return reference
}

// CheckNotNullMsg2 is Preconditions.checkNotNull(T, String, Object, Object).
func CheckNotNullMsg2[T, A, B any](reference T, template string, p1 A, p2 B) T {
	if IsNil(reference) {
		failNotNull(template, p1, p2)
	}
	return reference
}

//go:noinline
func failNotNull(template string, args ...any) {
	panic(NewNullPointerExceptionMsg(LenientFormat(template, args...)))
}

// CheckElementIndex is Preconditions.checkElementIndex(int, int).
func CheckElementIndex(index, size int) int {
	if index < 0 || index >= size {
		panic(NewIndexOutOfBoundsException(LenientFormat("index (%s) must be less than size (%s)", index, size)))
	}
	return index
}

// LenientFormat is Guava's Strings.lenientFormat: each "%s" is replaced by
// the next argument's Java string form; leftover arguments are appended as
// " [a, b]".
func LenientFormat(template string, args ...any) string {
	var sb strings.Builder
	i := 0
	for i < len(args) {
		p := strings.Index(template, "%s")
		if p < 0 {
			break
		}
		sb.WriteString(template[:p])
		sb.WriteString(ToJavaString(args[i]))
		template = template[p+2:]
		i++
	}
	sb.WriteString(template)
	if i < len(args) {
		sb.WriteString(" [")
		for j := i; j < len(args); j++ {
			if j > i {
				sb.WriteString(", ")
			}
			sb.WriteString(ToJavaString(args[j]))
		}
		sb.WriteString("]")
	}
	return sb.String()
}
