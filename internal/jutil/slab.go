package jutil

// Slab is a transitional shim: New is new(T). Callers should call new(T)
// themselves; the remaining call sites will be removed.
//
// It used to allocate values in chunks, but a pointer to any element kept
// its whole chunk alive: on typed-full that held 7% more live heap (and
// 15-20% more RSS) than one allocation per value, which costs 7-10% more CPU.
type Slab[T any] struct{}

// NewFixedSlab returns a Slab; the chunk size is ignored.
func NewFixedSlab[T any](int) Slab[T] { return Slab[T]{} }

// New returns a pointer to a new zero T.
func (*Slab[T]) New() *T { return new(T) }
