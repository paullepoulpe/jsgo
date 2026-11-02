package jscomp

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const (
	smallSetInlineCap = 16
	smallMapInlineCap = 8
	inlineSmallCap    = smallSetInlineCap
)

// smallSet is a reusable insertion-ordered set backed by an inline buffer for
// <= smallSetInlineCap elements and a lazily populated hash index above that
// threshold.
type smallSet[K comparable] struct {
	items []K
	index map[K]struct{}
	buf   [smallSetInlineCap]K
}

func (s *smallSet[K]) clear() {
	if s == nil {
		return
	}
	if s.items == nil {
		s.items = s.buf[:0]
	} else {
		clear(s.items)
		s.items = s.items[:0]
	}
	if len(s.index) > 0 {
		clear(s.index)
	}
}

func (s *smallSet[K]) isEmpty() bool {
	return s == nil || len(s.items) == 0
}

func (s *smallSet[K]) size() int {
	if s == nil {
		return 0
	}
	return len(s.items)
}

func (s *smallSet[K]) contains(k K) bool {
	if s == nil || len(s.items) == 0 {
		return false
	}
	if len(s.items) <= smallSetInlineCap {
		return slices.Contains(s.items, k)
	}
	_, ok := s.index[k]
	return ok
}

func (s *smallSet[K]) add(k K) bool {
	if len(s.items) <= smallSetInlineCap {
		if slices.Contains(s.items, k) {
			return false
		}
		if s.items == nil {
			s.items = s.buf[:0]
		}
		if len(s.items) == smallSetInlineCap {
			if s.index == nil {
				s.index = make(map[K]struct{}, smallSetInlineCap*2)
			} else {
				clear(s.index)
			}
			for _, existing := range s.items {
				s.index[existing] = struct{}{}
			}
			s.index[k] = struct{}{}
		}
		s.items = append(s.items, k)
		return true
	}
	if _, exists := s.index[k]; exists {
		return false
	}
	s.index[k] = struct{}{}
	s.items = append(s.items, k)
	return true
}

// smallIDSet is a reusable insertion-ordered set of non-zero 32-bit strpool IDs.
type smallIDSet struct {
	smallSet[uint32]
}

type inlineIDSet = smallIDSet

func (s *smallIDSet) clear() {
	if s == nil {
		return
	}
	s.smallSet.clear()
}

func (s *smallIDSet) isEmpty() bool {
	return s == nil || s.smallSet.isEmpty()
}

func (s *smallIDSet) size() int {
	if s == nil {
		return 0
	}
	return s.smallSet.size()
}

func (s *smallIDSet) contains(id uint32) bool {
	if s == nil || id == 0 {
		return false
	}
	return s.smallSet.contains(id)
}

func (s *smallIDSet) add(id uint32) bool {
	if id == 0 {
		return false
	}
	return s.smallSet.add(id)
}

func (s *smallIDSet) addAll(other *smallIDSet) {
	if other == nil {
		return
	}
	for _, id := range other.items {
		s.add(id)
	}
}

func (s *smallIDSet) copyFrom(other *smallIDSet) {
	s.clear()
	if other != nil {
		for _, id := range other.items {
			s.add(id)
		}
	}
}

// smallMap is a reusable map backed by parallel inline key/value arrays for
// <= smallMapInlineCap entries and a lazily populated hash index above that
// threshold. Insertion order is preserved unless removeSwap is called.
type smallMap[K comparable, V any] struct {
	keys   []K
	vals   []V
	index  map[K]int
	keyBuf [smallMapInlineCap]K
	valBuf [smallMapInlineCap]V
}

func (m *smallMap[K, V]) clear() {
	if m == nil {
		return
	}
	if m.keys == nil {
		m.keys = m.keyBuf[:0]
		m.vals = m.valBuf[:0]
	} else {
		clear(m.keys)
		clear(m.vals)
		m.keys = m.keys[:0]
		m.vals = m.vals[:0]
	}
	if len(m.index) > 0 {
		clear(m.index)
	}
}

func (m *smallMap[K, V]) release() {
	if m == nil {
		return
	}
	if len(m.keys) > 0 {
		clear(m.keys)
		clear(m.vals)
	}
	m.keys = nil
	m.vals = nil
	m.index = nil
}

func (m *smallMap[K, V]) isEmpty() bool {
	return m == nil || len(m.keys) == 0
}

func (m *smallMap[K, V]) size() int {
	if m == nil {
		return 0
	}
	return len(m.keys)
}

func (m *smallMap[K, V]) values() []V {
	if m == nil {
		return nil
	}
	return m.vals
}

func (m *smallMap[K, V]) indexOf(k K) int {
	if m == nil || len(m.keys) == 0 {
		return -1
	}
	if len(m.keys) <= smallMapInlineCap {
		return slices.Index(m.keys, k)
	}
	if idx, ok := m.index[k]; ok {
		return idx
	}
	return -1
}

func (m *smallMap[K, V]) get(k K) (V, bool) {
	if idx := m.indexOf(k); idx >= 0 {
		return m.vals[idx], true
	}
	var zero V
	return zero, false
}

func (m *smallMap[K, V]) getOrZero(k K) V {
	v, _ := m.get(k)
	return v
}

func (m *smallMap[K, V]) rebuildIndex() {
	if m.index == nil {
		m.index = make(map[K]int, len(m.keys)*2)
	} else {
		clear(m.index)
	}
	for i, k := range m.keys {
		m.index[k] = i
	}
}

func (m *smallMap[K, V]) appendNew(k K, v V) int {
	if m.keys == nil {
		m.keys = m.keyBuf[:0]
		m.vals = m.valBuf[:0]
	}
	newIdx := len(m.keys)
	m.keys = append(m.keys, k)
	m.vals = append(m.vals, v)
	if len(m.keys) == smallMapInlineCap+1 {
		m.rebuildIndex()
	} else if len(m.keys) > smallMapInlineCap+1 {
		m.index[k] = newIdx
	}
	return newIdx
}

func (m *smallMap[K, V]) put(k K, v V) (prev V, replaced bool) {
	if idx := m.indexOf(k); idx >= 0 {
		prev = m.vals[idx]
		m.vals[idx] = v
		return prev, true
	}
	m.appendNew(k, v)
	return prev, false
}

func (m *smallMap[K, V]) getOrInsert(k K, init V) (val *V, existed bool) {
	if idx := m.indexOf(k); idx >= 0 {
		return &m.vals[idx], true
	}
	idx := m.appendNew(k, init)
	return &m.vals[idx], false
}

func (m *smallMap[K, V]) removeSwap(k K) bool {
	idx := m.indexOf(k)
	if idx < 0 {
		return false
	}
	last := len(m.keys) - 1
	wasIndexed := len(m.keys) > smallMapInlineCap
	if idx != last {
		movedKey := m.keys[last]
		m.keys[idx] = movedKey
		m.vals[idx] = m.vals[last]
		if wasIndexed {
			m.index[movedKey] = idx
		}
	}
	var zeroK K
	var zeroV V
	m.keys[last] = zeroK
	m.vals[last] = zeroV
	m.keys = m.keys[:last]
	m.vals = m.vals[:last]
	if wasIndexed {
		if len(m.keys) <= smallMapInlineCap {
			clear(m.index)
		} else {
			delete(m.index, k)
		}
	}
	return true
}

func (m *smallMap[K, V]) retainFunc(keep func(K, V) bool) {
	if m == nil || len(m.keys) == 0 {
		return
	}
	w := 0
	for i, v := range m.vals {
		k := m.keys[i]
		if keep(k, v) {
			m.keys[w] = k
			m.vals[w] = v
			w++
		}
	}
	clear(m.keys[w:])
	clear(m.vals[w:])
	m.keys = m.keys[:w]
	m.vals = m.vals[:w]
	if w > smallMapInlineCap {
		m.rebuildIndex()
	} else if len(m.index) > 0 {
		clear(m.index)
	}
}

// smallIDNodeMap maps non-zero 32-bit strpool IDs to replacement Nodes.
type smallIDNodeMap struct {
	smallMap[uint32, rhino.Node]
}

type inlineIDNodeMap = smallIDNodeMap

func (m *smallIDNodeMap) clear() {
	if m == nil {
		return
	}
	m.smallMap.clear()
}

func (m *smallIDNodeMap) release() {
	if m == nil {
		return
	}
	m.smallMap.release()
}

func (m *smallIDNodeMap) isEmpty() bool {
	return m == nil || m.smallMap.isEmpty()
}

func (m *smallIDNodeMap) size() int {
	if m == nil {
		return 0
	}
	return m.smallMap.size()
}

func (m *smallIDNodeMap) get(id uint32) rhino.Node {
	if m == nil || id == 0 {
		return rhino.NilNode
	}
	return m.smallMap.getOrZero(id)
}

func (m *smallIDNodeMap) put(id uint32, node rhino.Node) rhino.Node {
	if id == 0 {
		return rhino.NilNode
	}
	prev, _ := m.smallMap.put(id, node)
	return prev
}

func (m *smallIDNodeMap) remove(id uint32) {
	if m == nil || id == 0 {
		return
	}
	m.smallMap.removeSwap(id)
}
