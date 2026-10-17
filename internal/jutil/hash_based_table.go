package jutil

import "iter"

// Go-only: language-adapter — Guava HashBasedTable (via StandardTable): a
// LinkedHashMap from row key to a LinkedHashMap from column key to value,
// so rows and columns iterate in insertion order. Row and column keys use
// Go equality (see LinkedHashMap).

// HashBasedTable is com.google.common.collect.HashBasedTable.
type HashBasedTable[R comparable, C comparable, V any] struct {
	backingMap *LinkedHashMap[R, *LinkedHashMap[C, V]]
	size       int
}

// HashBasedTable_Create is HashBasedTable.create().
func HashBasedTable_Create[R comparable, C comparable, V any]() *HashBasedTable[R, C, V] {
	return &HashBasedTable[R, C, V]{backingMap: NewLinkedHashMap[R, *LinkedHashMap[C, V]]()}
}

// Put is Table.put(R, C, V): it returns the previous value (zero if none).
func (t *HashBasedTable[R, C, V]) Put(rowKey R, columnKey C, value V) V {
	CheckNotNull(any(rowKey))
	CheckNotNull(any(columnKey))
	CheckNotNull(any(value))
	row, ok := t.backingMap.Lookup(rowKey)
	if !ok {
		row = NewLinkedHashMap[C, V]()
		t.backingMap.Put(rowKey, row)
	}
	old, had := row.Lookup(columnKey)
	row.Put(columnKey, value)
	if !had {
		t.size++
	}
	return old
}

// Get is Table.get(Object, Object).
func (t *HashBasedTable[R, C, V]) Get(rowKey R, columnKey C) (V, bool) {
	if row, ok := t.backingMap.Lookup(rowKey); ok {
		return row.Lookup(columnKey)
	}
	var zero V
	return zero, false
}

// Contains is Table.contains(Object, Object).
func (t *HashBasedTable[R, C, V]) Contains(rowKey R, columnKey C) bool {
	_, ok := t.Get(rowKey, columnKey)
	return ok
}

// ContainsRow is Table.containsRow(Object).
func (t *HashBasedTable[R, C, V]) ContainsRow(rowKey R) bool {
	return t.backingMap.ContainsKey(rowKey)
}

// IsEmpty is Table.isEmpty().
func (t *HashBasedTable[R, C, V]) IsEmpty() bool { return t.backingMap.IsEmpty() }

// Size is Table.size(): the number of cells.
func (t *HashBasedTable[R, C, V]) Size() int { return t.size }

// Row is Table.row(R). The result must be treated as read-only; for a
// missing row it is a new empty map (Guava returns a view that would see
// later insertions).
func (t *HashBasedTable[R, C, V]) Row(rowKey R) *LinkedHashMap[C, V] {
	if row, ok := t.backingMap.Lookup(rowKey); ok {
		return row
	}
	return NewLinkedHashMap[C, V]()
}

// RowKeySet is Table.rowKeySet(), in insertion order.
func (t *HashBasedTable[R, C, V]) RowKeySet() iter.Seq[R] { return t.backingMap.Keys() }

// RowMap is Table.rowMap(); it must be treated as read-only.
func (t *HashBasedTable[R, C, V]) RowMap() *LinkedHashMap[R, *LinkedHashMap[C, V]] {
	return t.backingMap
}

// String is StandardTable.toString(): rowMap().toString().
func (t *HashBasedTable[R, C, V]) String() string { return t.backingMap.String() }

// Remove is Table.remove(Object, Object): it removes the cell and returns
// its value (zero and false if there was none). A row left empty is
// removed, so it no longer appears in RowKeySet or RowMap.
func (t *HashBasedTable[R, C, V]) Remove(rowKey R, columnKey C) (V, bool) {
	row, ok := t.backingMap.Lookup(rowKey)
	if !ok {
		var zero V
		return zero, false
	}
	old, had := row.Lookup(columnKey)
	if !had {
		var zero V
		return zero, false
	}
	row.Remove(columnKey)
	t.size--
	if row.IsEmpty() {
		t.backingMap.Remove(rowKey)
	}
	return old, true
}

// Values is Table.values(): every cell's value, row by row in row insertion
// order and, within a row, in column insertion order.
func (t *HashBasedTable[R, C, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, row := range t.backingMap.All() {
			for _, v := range row.All() {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// ValueSlice is a copy of Values() as a slice (ImmutableList.copyOf(values())).
func (t *HashBasedTable[R, C, V]) ValueSlice() []V {
	out := make([]V, 0, t.size)
	for v := range t.Values() {
		out = append(out, v)
	}
	return out
}
