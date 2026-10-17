package jutil

import "iter"

// Go-only: language-adapter — java.util.ArrayDeque.

// ArrayDeque is java.util.ArrayDeque: a resizable ring buffer used as a
// stack (Push/Pop/Peek at the front) or a queue (Add at the back, Poll at
// the front). Iteration (All) runs from the first (front) element to the
// last. As in Java, null elements are rejected with a NullPointerException.
type ArrayDeque[T any] struct {
	buf  []T
	head int
	size int
}

// NewArrayDeque is new ArrayDeque<>().
func NewArrayDeque[T any]() *ArrayDeque[T] { return &ArrayDeque[T]{} }

// NewArrayDequeWithCapacity is new ArrayDeque<>(numElements). The capacity
// is only a hint; it does not change the behaviour.
func NewArrayDequeWithCapacity[T any](numElements int) *ArrayDeque[T] {
	if numElements < 1 {
		numElements = 1
	}
	return &ArrayDeque[T]{buf: make([]T, numElements)}
}

// NewArrayDequeFrom is new ArrayDeque<>(collection): the elements in the
// collection's iteration order.
func NewArrayDequeFrom[T any](elems iter.Seq[T]) *ArrayDeque[T] {
	d := NewArrayDeque[T]()
	for e := range elems {
		d.AddLast(e)
	}
	return d
}

func (d *ArrayDeque[T]) grow() {
	n := len(d.buf) * 2
	if n < 8 {
		n = 8
	}
	nb := make([]T, n)
	for i := 0; i < d.size; i++ {
		nb[i] = d.buf[(d.head+i)%len(d.buf)]
	}
	d.buf = nb
	d.head = 0
}

func arrayDequeCheckNotNull[T any](e T) {
	if IsNil(e) {
		panic(NewNullPointerException())
	}
}

// Size is ArrayDeque.size().
func (d *ArrayDeque[T]) Size() int { return d.size }

// IsEmpty is ArrayDeque.isEmpty().
func (d *ArrayDeque[T]) IsEmpty() bool { return d.size == 0 }

// AddFirst is ArrayDeque.addFirst(E).
func (d *ArrayDeque[T]) AddFirst(e T) {
	arrayDequeCheckNotNull(e)
	if d.size == len(d.buf) {
		d.grow()
	}
	d.head = (d.head - 1 + len(d.buf)) % len(d.buf)
	d.buf[d.head] = e
	d.size++
}

// AddLast is ArrayDeque.addLast(E).
func (d *ArrayDeque[T]) AddLast(e T) {
	arrayDequeCheckNotNull(e)
	if d.size == len(d.buf) {
		d.grow()
	}
	d.buf[(d.head+d.size)%len(d.buf)] = e
	d.size++
}

// Add is ArrayDeque.add(E) (addLast); always true.
func (d *ArrayDeque[T]) Add(e T) bool {
	d.AddLast(e)
	return true
}

// AddAll is ArrayDeque.addAll(Collection): adds each element at the back;
// reports whether the deque changed.
func (d *ArrayDeque[T]) AddAll(elems iter.Seq[T]) bool {
	changed := false
	for e := range elems {
		d.AddLast(e)
		changed = true
	}
	return changed
}

// Push is ArrayDeque.push(E) (addFirst).
func (d *ArrayDeque[T]) Push(e T) { d.AddFirst(e) }

// PeekFirst is ArrayDeque.peekFirst(): the zero value (Java null) when
// empty.
func (d *ArrayDeque[T]) PeekFirst() T {
	if d.size == 0 {
		var zero T
		return zero
	}
	return d.buf[d.head]
}

// PeekLast is ArrayDeque.peekLast(): the zero value (Java null) when empty.
func (d *ArrayDeque[T]) PeekLast() T {
	if d.size == 0 {
		var zero T
		return zero
	}
	return d.buf[(d.head+d.size-1)%len(d.buf)]
}

// Peek is ArrayDeque.peek() (peekFirst).
func (d *ArrayDeque[T]) Peek() T { return d.PeekFirst() }

// PollFirst is ArrayDeque.pollFirst(): the zero value (Java null) when
// empty.
func (d *ArrayDeque[T]) PollFirst() T {
	var zero T
	if d.size == 0 {
		return zero
	}
	e := d.buf[d.head]
	d.buf[d.head] = zero
	d.head = (d.head + 1) % len(d.buf)
	d.size--
	return e
}

// PollLast is ArrayDeque.pollLast(): the zero value (Java null) when empty.
func (d *ArrayDeque[T]) PollLast() T {
	var zero T
	if d.size == 0 {
		return zero
	}
	i := (d.head + d.size - 1) % len(d.buf)
	e := d.buf[i]
	d.buf[i] = zero
	d.size--
	return e
}

// Poll is ArrayDeque.poll() (pollFirst).
func (d *ArrayDeque[T]) Poll() T { return d.PollFirst() }

// RemoveFirst is ArrayDeque.removeFirst(): NoSuchElementException when
// empty.
func (d *ArrayDeque[T]) RemoveFirst() T {
	if d.size == 0 {
		panic(NewNoSuchElementException())
	}
	return d.PollFirst()
}

// RemoveLast is ArrayDeque.removeLast(): NoSuchElementException when
// empty.
func (d *ArrayDeque[T]) RemoveLast() T {
	if d.size == 0 {
		panic(NewNoSuchElementException())
	}
	return d.PollLast()
}

// Pop is ArrayDeque.pop() (removeFirst).
func (d *ArrayDeque[T]) Pop() T { return d.RemoveFirst() }

// GetFirst is ArrayDeque.getFirst(): NoSuchElementException when empty.
func (d *ArrayDeque[T]) GetFirst() T {
	if d.size == 0 {
		panic(NewNoSuchElementException())
	}
	return d.PeekFirst()
}

// GetLast is ArrayDeque.getLast(): NoSuchElementException when empty.
func (d *ArrayDeque[T]) GetLast() T {
	if d.size == 0 {
		panic(NewNoSuchElementException())
	}
	return d.PeekLast()
}

// Clear is ArrayDeque.clear().
func (d *ArrayDeque[T]) Clear() {
	clear(d.buf)
	d.head = 0
	d.size = 0
}

// Contains is ArrayDeque.contains(Object), comparing with ObjectsEquals.
func (d *ArrayDeque[T]) Contains(o T) bool {
	for i := 0; i < d.size; i++ {
		if ObjectsEquals(o, d.buf[(d.head+i)%len(d.buf)]) {
			return true
		}
	}
	return false
}

// All iterates from the first element to the last (ArrayDeque.iterator()).
func (d *ArrayDeque[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < d.size; i++ {
			if !yield(d.buf[(d.head+i)%len(d.buf)]) {
				return
			}
		}
	}
}

// Descending iterates from the last element to the first
// (ArrayDeque.descendingIterator()).
func (d *ArrayDeque[T]) Descending() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := d.size - 1; i >= 0; i-- {
			if !yield(d.buf[(d.head+i)%len(d.buf)]) {
				return
			}
		}
	}
}

// ToSlice returns the elements from first to last (ArrayDeque.toArray()).
func (d *ArrayDeque[T]) ToSlice() []T {
	out := make([]T, 0, d.size)
	for e := range d.All() {
		out = append(out, e)
	}
	return out
}

// String is AbstractCollection.toString().
func (d *ArrayDeque[T]) String() string { return CollectionToString(d.ToSlice()) }
