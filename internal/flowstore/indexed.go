// Copyright 2026 The jsgo Authors.
// Licensed under the Apache License, Version 2.0.

package flowstore

import (
	"iter"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const pageBits = 2
const pageSize = 1 << pageBits

// MaxID bounds persistent directory copying. Larger scopes use the original HAMT.
const MaxID = 256

type IDs struct {
	mu     sync.RWMutex
	ids    map[string]int
	hashes map[int32]string
	names  []string
	order  []int
}

func (ids *IDs) lookup(name string) (int, bool) {
	ids.mu.RLock()
	id, ok := ids.ids[name]
	ids.mu.RUnlock()
	return id, ok
}

func (ids *IDs) intern(name string) int {
	if id, ok := ids.lookup(name); ok {
		return id
	}
	ids.mu.Lock()
	defer ids.mu.Unlock()
	if id, ok := ids.ids[name]; ok {
		return id
	}
	// Do not retain arbitrary qualified names once a scope exceeds the bound.
	if len(ids.ids) >= MaxID {
		return MaxID
	}
	// HamtPMap's collision pivot behavior is observable. Keep collided symbols
	// on that exact backend rather than replacing its operation history with a
	// mathematically equivalent map whose join results can differ.
	hash := jutil.StringHashCode(name)
	if previous, ok := ids.hashes[hash]; ok && previous != name {
		return MaxID
	}
	id := len(ids.ids)
	ids.hashes[hash] = name
	ids.ids[name] = id
	ids.names = append(ids.names, name)
	// Publish an immutable canonical HAMT key order. No per-join sorting/copying.
	order := make([]int, len(ids.order)+1)
	position := 0
	orderedHash := uint32(rhino.HamtPMap_Hash(name))
	for position < len(ids.order) && uint32(rhino.HamtPMap_Hash(ids.names[ids.order[position]])) < orderedHash {
		position++
	}
	copy(order, ids.order[:position])
	order[position] = id
	copy(order[position+1:], ids.order[position:])
	ids.order = order
	return id
}

func (ids *IDs) name(id int) string { ids.mu.RLock(); defer ids.mu.RUnlock(); return ids.names[id] }

func (ids *IDs) ordered() []int { ids.mu.RLock(); defer ids.mu.RUnlock(); return ids.order }

type page[V comparable] [pageSize]V

// Store is a persistent string-keyed store for identity-comparable nonzero values.
// It implements the PMap subset needed by inferred slots (whose pointer identity
// is their Object.equals behavior). Zero values represent absence.
type Store[V comparable] struct {
	ids   *IDs
	pages []*page[V]
	count int
	// Non-nil only after the persistent-vector eligibility bound is exceeded.
	// Keeping the wrapper allows old vector versions and new HAMT versions to join.
	fallback *rhino.HamtPMap[string, V]
}

// New creates an empty immutable version sharing the given scope-local symbol IDs.
func New[V comparable](ids *IDs) *Store[V] { return &Store[V]{ids: ids} }

// NewIDs creates an owner whose stable IDs live as long as its scope/versions.
func NewIDs() *IDs { return &IDs{ids: make(map[string]int), hashes: make(map[int32]string)} }

func (m *Store[V]) IsEmpty() bool {
	if m.fallback != nil {
		return m.fallback.IsEmpty()
	}
	return m.count == 0
}

func (m *Store[V]) Values() iter.Seq[V] {
	if m.fallback != nil {
		return m.fallback.Values()
	}
	return func(yield func(V) bool) {
		for _, page := range m.pages {
			if page == nil {
				continue
			}
			for _, slot := range page {
				if slot != *new(V) && !yield(slot) {
					return
				}
			}
		}
	}
}

func (m *Store[V]) Keys() iter.Seq[string] {
	if m.fallback != nil {
		return m.fallback.Keys()
	}
	return func(yield func(string) bool) {
		for p, page := range m.pages {
			if page == nil {
				continue
			}
			for k, value := range page {
				if value != *new(V) && !yield(m.ids.name(p*pageSize+k)) {
					return
				}
			}
		}
	}
}

func (m *Store[V]) getID(id int) V {
	p := id >> pageBits
	if p >= len(m.pages) || m.pages[p] == nil {
		var zero V
		return zero
	}
	return m.pages[p][id&(pageSize-1)]
}

func (m *Store[V]) Get(name string) V {
	if m.fallback != nil {
		return m.fallback.Get(name)
	}
	id, ok := m.ids.lookup(name)
	if !ok {
		var zero V
		return zero
	}
	return m.getID(id)
}

func (m *Store[V]) withID(id int, slot V) *Store[V] {
	old := m.getID(id)
	if old == slot {
		return m
	}
	pageID := id >> pageBits
	length := max(len(m.pages), pageID+1)
	pages := make([]*page[V], length)
	copy(pages, m.pages)
	page := new(page[V])
	if pageID < len(m.pages) && m.pages[pageID] != nil {
		*page = *m.pages[pageID]
	}
	page[id&(pageSize-1)] = slot
	pages[pageID] = page
	count := m.count
	if old == *new(V) {
		count++
	}
	if slot == *new(V) {
		count--
	}
	return &Store[V]{ids: m.ids, pages: pages, count: count}
}

func (m *Store[V]) asHAMT() *rhino.HamtPMap[string, V] {
	if m.fallback != nil {
		return m.fallback
	}
	result := rhino.HamtPMap_Empty[string, V]()
	for key := range m.Keys() {
		result = result.PlusMap(key, m.Get(key))
	}
	return result
}

func referenceMap[V comparable](other rhino.PMap[string, V]) *rhino.HamtPMap[string, V] {
	switch other := other.(type) {
	case *Store[V]:
		return other.asHAMT()
	case *rhino.HamtPMap[string, V]:
		return other
	default:
		result := rhino.HamtPMap_Empty[string, V]()
		for key := range other.Keys() {
			result = result.PlusMap(key, other.Get(key))
		}
		return result
	}
}

func (m *Store[V]) Plus(name string, slot V) rhino.PMap[string, V] {
	jutil.CheckNotNull(slot)
	jutil.CheckArgument(slot != *new(V))
	if m.fallback != nil {
		next := m.fallback.PlusMap(name, slot)
		if next == m.fallback {
			return m
		}
		return &Store[V]{ids: m.ids, fallback: next}
	}
	id := m.ids.intern(name)
	if id >= MaxID {
		return &Store[V]{ids: m.ids, fallback: m.asHAMT().PlusMap(name, slot)}
	}
	return m.withID(id, slot)
}

func (m *Store[V]) Minus(name string) rhino.PMap[string, V] {
	if m.fallback != nil {
		next := m.fallback.MinusMap(name)
		if next == m.fallback {
			return m
		}
		return &Store[V]{ids: m.ids, fallback: next}
	}
	id, ok := m.ids.lookup(name)
	if !ok || m.getID(id) == *new(V) {
		return m
	}
	return m.withID(id, *new(V))
}

func (m *Store[V]) Reconcile(other rhino.PMap[string, V], join rhino.PMap_Reconciler[string, V]) rhino.PMap[string, V] {
	that, ok := other.(*Store[V])
	if !ok || m.ids != that.ids || m.fallback != nil || that.fallback != nil {
		reference := referenceMap(other)
		joined := m.asHAMT().Reconcile(reference, join).(*rhino.HamtPMap[string, V])
		if joined == m.fallback {
			return m
		}
		return &Store[V]{ids: m.ids, fallback: joined}
	}
	if m == that {
		return m
	}

	// HAMT callbacks can resolve lazy types. Use its canonical unsigned hash
	// order, rather than ID encounter order, even though the storage is indexed.
	pages := make([]*page[V], max(len(m.pages), len(that.pages)))
	copy(pages, m.pages)
	var copied [MaxID / pageSize]bool
	count := m.count
	sameLeft, sameRight := true, true
	for _, id := range m.ids.ordered() {
		a, b := m.getID(id), that.getID(id)
		if a == b {
			continue
		}
		result := jutil.CheckNotNull(join(m.ids.name(id), a, b))
		sameLeft = sameLeft && result == a
		sameRight = sameRight && result == b
		if result == a {
			continue
		}
		p := id >> pageBits
		if !copied[p] {
			next := new(page[V])
			if pages[p] != nil {
				*next = *pages[p]
			}
			pages[p] = next
			copied[p] = true
		}
		pages[p][id&(pageSize-1)] = result
		if a == *new(V) {
			count++
		}
	}
	if sameLeft {
		return m
	}
	if sameRight {
		return that
	}
	return &Store[V]{ids: m.ids, pages: pages, count: count}
}

func (m *Store[V]) Equivalent(other rhino.PMap[string, V], equal func(V, V) bool) bool {
	that, ok := other.(*Store[V])
	if !ok || m.ids != that.ids || m.fallback != nil || that.fallback != nil {
		reference := referenceMap(other)
		return m.asHAMT().Equivalent(reference, equal)
	}
	if m == that {
		return true
	}

	// Different key sets can stop the reference traversal before or after some
	// equality callbacks. Delegate that exact case rather than changing lazy-type
	// resolution side effects by choosing another early-exit order.
	differentKeys := m.count != that.count
	if !differentKeys {
		for _, id := range m.ids.ordered() {
			if (m.getID(id) == *new(V)) != (that.getID(id) == *new(V)) {
				differentKeys = true
				break
			}
		}
	}
	if differentKeys {
		return m.asHAMT().Equivalent(that.asHAMT(), equal)
	}
	for _, id := range m.ids.ordered() {
		a, b := m.getID(id), that.getID(id)
		// EqualSlots itself returns immediately on pointer identity. Its only type
		// operations therefore occur on changed values, in the same HAMT order.
		if a != b && !equal(a, b) {
			return false
		}
	}
	return true
}
