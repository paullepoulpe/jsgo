package jutil

// Go-only: language-adapter — java.util.Objects.hash / Arrays.hashCode and
// List.hashCode / Set.hashCode for Go slices.

// ObjectsHash is java.util.Objects.hash(Object...), i.e.
// Arrays.hashCode(Object[]): 31-based polynomial over ObjectHashCode.
func ObjectsHash(values ...any) int32 {
	var result int32 = 1
	for _, v := range values {
		result = 31*result + ObjectHashCode(v)
	}
	return result
}

// ListHashCode is List.hashCode() for a slice standing in for a Java list.
func ListHashCode[T any](elems []T) int32 {
	var h int32 = 1
	for _, e := range elems {
		h = 31*h + ObjectHashCode(e)
	}
	return h
}

// SetHashCode is Set.hashCode() for a slice standing in for a Java set: the
// sum of the element hash codes.
func SetHashCode[T any](elems []T) int32 {
	var h int32
	for _, e := range elems {
		h += ObjectHashCode(e)
	}
	return h
}
