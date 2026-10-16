package jutil

// PriorityQueue is java.util.PriorityQueue with a comparator: the
// array-backed binary heap of OpenJDK, with its sift-up / sift-down rules,
// so that elements that compare equal come out in the same order as in
// Java.
type PriorityQueue[E any] struct {
	queue      []E
	comparator func(a, b E) int
}

// NewPriorityQueue is new PriorityQueue<>(comparator) (and
// new PriorityQueue<>(initialCapacity, comparator): the capacity is not
// observable).
func NewPriorityQueue[E any](comparator func(a, b E) int) *PriorityQueue[E] {
	return &PriorityQueue[E]{comparator: comparator}
}

// Add is PriorityQueue.add(E) / offer(E): throws NullPointerException for
// null; always returns true.
func (q *PriorityQueue[E]) Add(e E) bool {
	if IsNil(e) {
		panic(NewNullPointerException())
	}
	k := len(q.queue)
	var zero E
	q.queue = append(q.queue, zero)
	// siftUpUsingComparator(k, e, queue, comparator)
	for k > 0 {
		parent := (k - 1) >> 1
		p := q.queue[parent]
		if q.comparator(e, p) >= 0 {
			break
		}
		q.queue[k] = p
		k = parent
	}
	q.queue[k] = e
	return true
}

// Poll is PriorityQueue.poll(): removes and returns the head of the queue,
// or the zero value (Java null) if it is empty.
func (q *PriorityQueue[E]) Poll() E {
	var zero E
	if len(q.queue) == 0 {
		return zero
	}
	result := q.queue[0]
	n := len(q.queue) - 1
	x := q.queue[n]
	q.queue[n] = zero
	q.queue = q.queue[:n]
	if n > 0 {
		// siftDownUsingComparator(0, x, queue, n, comparator)
		k := 0
		half := n >> 1
		for k < half {
			child := (k << 1) + 1
			c := q.queue[child]
			right := child + 1
			if right < n && q.comparator(c, q.queue[right]) > 0 {
				child = right
				c = q.queue[child]
			}
			if q.comparator(x, c) <= 0 {
				break
			}
			q.queue[k] = c
			k = child
		}
		q.queue[k] = x
	}
	return result
}

// Peek is PriorityQueue.peek(): the head of the queue, or the zero value
// (Java null) if it is empty.
func (q *PriorityQueue[E]) Peek() E {
	if len(q.queue) == 0 {
		var zero E
		return zero
	}
	return q.queue[0]
}

// Size is PriorityQueue.size().
func (q *PriorityQueue[E]) Size() int { return len(q.queue) }

// IsEmpty is AbstractCollection.isEmpty().
func (q *PriorityQueue[E]) IsEmpty() bool { return len(q.queue) == 0 }

// Clear is PriorityQueue.clear().
func (q *PriorityQueue[E]) Clear() {
	clear(q.queue)
	q.queue = q.queue[:0]
}
