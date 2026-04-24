// Package runtime provides minimal Go re-implementations of Guava/JDK
// collection helpers used in generated javaport skeletons.
//
// Generated files that reference ImmutableList, ImmutableMap, ImmutableSet,
// Optional, Iterables.concat, or Preconditions.checkState should import this
// package as:
//
//	import "github.com/paullepoulpe/jsgo/tools/javaport/runtime"
//
// All functions are thin wrappers that mirror the Java semantics closely
// enough for ported code to compile and run.
package runtime

import "fmt"

// ImmutableList returns a copy of the given elements as a Go slice.
// Mirrors Guava ImmutableList.of(...)  and ImmutableList.copyOf(...).
func ImmutableList[T any](elems ...T) []T {
	out := make([]T, len(elems))
	copy(out, elems)
	return out
}

// ImmutableSet returns a map[T]struct{} from the given elements.
// Mirrors Guava ImmutableSet.of(...) and ImmutableSet.copyOf(...).
func ImmutableSet[T comparable](elems ...T) map[T]struct{} {
	out := make(map[T]struct{}, len(elems))
	for _, e := range elems {
		out[e] = struct{}{}
	}
	return out
}

// ImmutableMap returns a map built from alternating key, value arguments.
// Mirrors Guava ImmutableMap.of(k1,v1, k2,v2, ...).
func ImmutableMap[K comparable, V any](pairs ...any) map[K]V {
	if len(pairs)%2 != 0 {
		panic("runtime.ImmutableMap: odd number of arguments")
	}
	out := make(map[K]V, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		out[pairs[i].(K)] = pairs[i+1].(V)
	}
	return out
}

// Optional wraps a nullable pointer, mirroring Java Optional<T>.
// Use Optional[T](nil) for Optional.empty() and Optional(&v) for
// Optional.of(v).
func Optional[T any](v *T) *T { return v }

// OptionalOf wraps a value as a non-nil pointer (Optional.of).
func OptionalOf[T any](v T) *T { return &v }

// Concat concatenates multiple slices, mirroring Iterables.concat.
func Concat[T any](slices ...[]T) []T {
	total := 0
	for _, s := range slices {
		total += len(s)
	}
	out := make([]T, 0, total)
	for _, s := range slices {
		out = append(out, s...)
	}
	return out
}

// CheckState panics with msg if cond is false.
// Mirrors Preconditions.checkState(cond, msg).
func CheckState(cond bool, format string, args ...any) {
	if !cond {
		panic(fmt.Sprintf(format, args...))
	}
}

// CheckNotNil panics if v is nil.
// Mirrors Preconditions.checkNotNull(v).
func CheckNotNil[T any](v *T, msg string) *T {
	if v == nil {
		panic("checkNotNull: " + msg)
	}
	return v
}
