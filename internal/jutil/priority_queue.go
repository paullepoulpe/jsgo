package jutil

import (
	"iter"
	"slices"
)

// Go-only: language-adapter — java.util.PriorityQueue with a comparator,
// written from its Javadoc ("an unbounded priority queue based on a priority
// heap") and the textbook array binary heap (element i has children 2i+1 and
// 2i+2). Which of several equal elements leaves first, and the iteration
// order (the heap array, from index 0), are observable; the tie rules below
// were determined by black-box experiments against JDK 21 and are checked by
// testdata/priority_queue.tsv:
//   - offer appends the element and sifts it up while it is strictly smaller
//     than its parent (an equal parent stays above it); adding to an empty
//     queue calls no comparator;
//   - poll moves the last element to the root and sifts it down: it follows
//     the smaller child, the left one when they are equal, and stops as soon
//     as it is not greater than that child;
//   - remove(Object) removes the first equal element in array order: the
//     last element takes its slot and sifts down, or, if it did not move,
//     sifts up.

// PriorityQueue is java.util.PriorityQueue<E> with a comparator.
type PriorityQueue[E comparable] struct {
	compare  func(a, b E) int
	queue    []E
	modCount int // structural changes, for the fail-fast iteration
}

// NewPriorityQueue is new PriorityQueue<>(comparator).
func NewPriorityQueue[E comparable](compare func(a, b E) int) *PriorityQueue[E] {
	return &PriorityQueue[E]{compare: compare}
}

// NewPriorityQueueWithCapacity is new PriorityQueue<>(initialCapacity,
// comparator); the capacity must be at least 1.
func NewPriorityQueueWithCapacity[E comparable](initialCapacity int, compare func(a, b E) int) *PriorityQueue[E] {
	if initialCapacity < 1 {
		panic(NewIllegalArgumentException())
	}
	return &PriorityQueue[E]{compare: compare, queue: make([]E, 0, min(initialCapacity, 1<<10))}
}

// Add is Queue.add: offer, which always succeeds.
func (q *PriorityQueue[E]) Add(e E) bool { return q.Offer(e) }

// Offer is Queue.offer: inserts e (NullPointerException for null).
func (q *PriorityQueue[E]) Offer(e E) bool {
	if IsNil(any(e)) {
		panic(NewNullPointerException())
	}
	q.modCount++
	q.queue = append(q.queue, e)
	q.siftUp(len(q.queue)-1, e)
	return true
}

// siftUp places x, destined for slot i, above every strictly greater
// ancestor.
func (q *PriorityQueue[E]) siftUp(i int, x E) {
	for i > 0 {
		parent := (i - 1) / 2
		p := q.queue[parent]
		if q.compare(x, p) >= 0 {
			break
		}
		q.queue[i] = p
		i = parent
	}
	q.queue[i] = x
}

// siftDown places x, destined for slot i, below every smaller descendant on
// the path of smaller children (the left child on ties).
func (q *PriorityQueue[E]) siftDown(i int, x E) {
	n := len(q.queue)
	for {
		child := 2*i + 1
		if child >= n {
			break
		}
		c := q.queue[child]
		if right := child + 1; right < n && q.compare(c, q.queue[right]) > 0 {
			child = right
			c = q.queue[right]
		}
		if q.compare(x, c) <= 0 {
			break
		}
		q.queue[i] = c
		i = child
	}
	q.queue[i] = x
}

// AddAll is AbstractQueue.addAll: adds the elements one by one; reports
// whether the queue changed.
func (q *PriorityQueue[E]) AddAll(c []E) bool {
	for _, e := range c {
		q.Add(e)
	}
	return len(c) > 0
}

// Peek is Queue.peek: the head, or the zero value (null) when empty.
func (q *PriorityQueue[E]) Peek() E {
	if len(q.queue) == 0 {
		var zero E
		return zero
	}
	return q.queue[0]
}

// Element is Queue.element: the head; NoSuchElementException when empty.
func (q *PriorityQueue[E]) Element() E {
	if len(q.queue) == 0 {
		panic(NewNoSuchElementException())
	}
	return q.queue[0]
}

// Poll is Queue.poll: removes and returns the head, or the zero value (null)
// when empty.
func (q *PriorityQueue[E]) Poll() E {
	var zero E
	n := len(q.queue)
	if n == 0 {
		return zero
	}
	q.modCount++
	head := q.queue[0]
	last := q.queue[n-1]
	q.queue[n-1] = zero
	q.queue = q.queue[:n-1]
	if n > 1 {
		q.siftDown(0, last)
	}
	return head
}

// Remove is Queue.remove(): the head; NoSuchElementException when empty.
func (q *PriorityQueue[E]) Remove() E {
	if len(q.queue) == 0 {
		panic(NewNoSuchElementException())
	}
	return q.Poll()
}

// indexOf is the slot of the first element o equals (o.equals(e)), or -1.
func (q *PriorityQueue[E]) indexOf(o any) int {
	if IsNil(o) {
		return -1
	}
	for i, e := range q.queue {
		if ObjectsEquals(o, any(e)) {
			return i
		}
	}
	return -1
}

// RemoveObject is Collection.remove(Object): removes one element equal to o.
func (q *PriorityQueue[E]) RemoveObject(o any) bool {
	i := q.indexOf(o)
	if i < 0 {
		return false
	}
	q.modCount++
	var zero E
	last := len(q.queue) - 1
	moved := q.queue[last]
	q.queue[last] = zero
	q.queue = q.queue[:last]
	if i < last {
		q.siftDown(i, moved)
		if q.queue[i] == moved {
			q.siftUp(i, moved)
		}
	}
	return true
}

// Contains is Collection.contains(Object).
func (q *PriorityQueue[E]) Contains(o any) bool { return q.indexOf(o) >= 0 }

// Size is Collection.size().
func (q *PriorityQueue[E]) Size() int { return len(q.queue) }

// IsEmpty is Collection.isEmpty().
func (q *PriorityQueue[E]) IsEmpty() bool { return len(q.queue) == 0 }

// Clear is Collection.clear().
func (q *PriorityQueue[E]) Clear() {
	q.modCount++
	clear(q.queue)
	q.queue = q.queue[:0]
}

// All iterates in the iterator's order: the heap array from the head. Like
// Java's iterator it throws ConcurrentModificationException when the queue
// changed and another element is due (the end is checked against the
// current size).
func (q *PriorityQueue[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		expected := q.modCount
		for i := 0; i < len(q.queue); i++ {
			if q.modCount != expected {
				panic(NewConcurrentModificationException())
			}
			if !yield(q.queue[i]) {
				return
			}
		}
	}
}

// ToSlice is toArray(): the elements in iteration order.
func (q *PriorityQueue[E]) ToSlice() []E { return slices.Clone(q.queue) }

// String is AbstractCollection.toString(): "[a, b]".
func (q *PriorityQueue[E]) String() string { return javaCollectionString(any(q), q.All()) }
