package jutil

import (
	"reflect"
	"strings"
)

// Go-only: language-adapter — Guava com.google.common.base.Preconditions.

// CheckState is Preconditions.checkState(boolean).
func CheckState(expression bool) {
	if !expression {
		panic(&IllegalStateException{RuntimeException{ThrowableBase{className: "java.lang.IllegalStateException"}}})
	}
}

// CheckStateMsg is Preconditions.checkState(boolean, String, Object...).
func CheckStateMsg(expression bool, template string, args ...any) {
	if !expression {
		failState(template, args...)
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
		panic(&IllegalArgumentException{RuntimeException{ThrowableBase{className: "java.lang.IllegalArgumentException"}}})
	}
}

// CheckArgumentMsg is Preconditions.checkArgument(boolean, String, Object...).
func CheckArgumentMsg(expression bool, template string, args ...any) {
	if !expression {
		failArgument(template, args...)
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

// IsNil reports whether v is Java null: a nil interface, a nil pointer,
// map, slice, func or chan stored in an interface, or a node id whose IsNil
// method says so (rhino.Node, which is nil when 0).
func IsNil(v any) bool {
	if v == nil {
		return true
	}
	switch v.(type) {
	case string, int, int32, int64, uint16, float64, float32, bool:
		return false
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return rv.IsNil()
	case reflect.Uint32:
		if h, ok := v.(interface{ IsNil() bool }); ok {
			return h.IsNil()
		}
	}
	return false
}

//go:noinline
func failNull() {
	panic(&NullPointerException{RuntimeException{ThrowableBase{className: "java.lang.NullPointerException"}}})
}

// CheckNotNullPtr is Preconditions.checkNotNull(T) for a typed pointer without reflection.
func CheckNotNullPtr[T any](p *T) *T {
	if p == nil {
		failNull()
	}
	return p
}

// CheckNotNull is Preconditions.checkNotNull(T).
func CheckNotNull[T any](reference T) T {
	if IsNil(reference) {
		failNull()
	}
	return reference
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
		failElementIndex(index, size)
	}
	return index
}

//go:noinline
func failElementIndex(index, size int) {
	panic(NewIndexOutOfBoundsException(LenientFormat("index (%s) must be less than size (%s)", index, size)))
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
