package jutil

import (
	"reflect"
	"strings"
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

// IsNil reports whether v is Java null: a nil interface, or a nil pointer,
// map, slice, func or chan stored in an interface.
func IsNil(v any) bool {
	if v == nil {
		return true
	}
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
		panic(NewNullPointerException())
	}
	return reference
}

// CheckNotNullMsg is Preconditions.checkNotNull(T, String, Object...).
func CheckNotNullMsg[T any](reference T, template string, args ...any) T {
	if IsNil(reference) {
		panic(NewNullPointerExceptionMsg(LenientFormat(template, args...)))
	}
	return reference
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
