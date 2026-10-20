package jscomp

import (
	"iter"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

const (
	maxMustDefLayoutVars  = 512
	maxMustDefLayoutWords = maxMustDefLayoutVars / 64
)

// Go-only: performance — a function's canonical Vars have compact local IDs.
// Larger functions retain the sparse representation to bound lattice memory.
// The fallback also preserves ScopedName equality for non-canonical Vars.
type mustDefLayout struct {
	vars     []*Var
	ids      map[*Var]uint16
	numWords int
	wide     bool
}

var emptyMustDefLayout = &mustDefLayout{}

func newMustDefLayout(vars *jutil.LinkedHashMap[string, *Var]) *mustDefLayout {
	if vars == nil || vars.Size() > maxMustDefLayoutVars {
		return nil
	}
	n := vars.Size()
	if n == 0 {
		return emptyMustDefLayout
	}
	l := &mustDefLayout{
		vars:     make([]*Var, 0, n),
		ids:      make(map[*Var]uint16, n),
		numWords: (n + 63) >> 6,
		wide:     n > 64,
	}
	for v := range vars.Values() {
		l.ids[v] = uint16(len(l.vars))
		l.vars = append(l.vars, v)
	}
	return l
}

// mustDefMap represents the same insertion-ordered partial map as the Java
// lattice. A clear presence bit is TOP; a present nil value is BOTTOM. The
// order slice preserves encounter order independently of the variable IDs.
type mustDefMap struct {
	layout       *mustDefLayout
	values       []*MustBeReachingVariableDef_Definition
	present      uint64
	presentWords [maxMustDefLayoutWords]uint64
	order        []uint16
	valuesBuf    [8]*MustBeReachingVariableDef_Definition
	orderBuf     [8]uint16
	fallback     *jutil.ObjectLinkedHashMap[*Var, *MustBeReachingVariableDef_Definition]
}

func newMustDefMap(layout *mustDefLayout) mustDefMap {
	m := mustDefMap{layout: layout}
	if layout == nil {
		m.fallback = jutil.NewObjectLinkedHashMap[*Var, *MustBeReachingVariableDef_Definition]()
	}
	return m
}

func (m *mustDefMap) hasID(id uint16) bool {
	if !m.layout.wide {
		return m.present&(uint64(1)<<id) != 0
	}
	return m.presentWords[id>>6]&(uint64(1)<<(id&63)) != 0
}

func (m *mustDefMap) init(layout *mustDefLayout) {
	m.layout = layout
	m.values = nil
	m.present = 0
	if layout != nil && layout.wide {
		clear(m.presentWords[:layout.numWords])
	}
	m.order = nil
	if layout == nil {
		m.fallback = jutil.NewObjectLinkedHashMap[*Var, *MustBeReachingVariableDef_Definition]()
	} else {
		m.fallback = nil
	}
}

func (m *mustDefMap) clone() mustDefMap {
	if m.layout != nil {
		return mustDefMap{
			layout:       m.layout,
			values:       slices.Clone(m.values),
			present:      m.present,
			presentWords: m.presentWords,
			order:        slices.Clone(m.order),
		}
	}
	n := newMustDefMap(nil)
	for v, d := range m.fallback.All() {
		n.fallback.Put(v, d)
	}
	return n
}

func (m *mustDefMap) cloneInto(dst *mustDefMap) {
	if m.layout != nil {
		dst.layout = m.layout
		if !m.layout.wide {
			dst.present = m.present
		} else {
			copy(dst.presentWords[:m.layout.numWords], m.presentWords[:m.layout.numWords])
		}
		dst.fallback = nil
		if m.values == nil {
			dst.values = nil
			dst.order = nil
			return
		}
		if cap(dst.values) >= len(m.values) {
			dst.values = dst.values[:len(m.values)]
		} else if len(m.values) <= len(dst.valuesBuf) {
			dst.values = dst.valuesBuf[:len(m.values)]
		} else {
			dst.values = make([]*MustBeReachingVariableDef_Definition, len(m.values))
		}
		copy(dst.values, m.values)
		if cap(dst.order) >= len(m.values) {
			dst.order = dst.order[:len(m.order)]
		} else if len(m.values) <= len(dst.orderBuf) {
			dst.order = dst.orderBuf[:len(m.order):len(dst.orderBuf)]
		} else {
			dst.order = make([]uint16, len(m.order), len(m.values))
		}
		copy(dst.order, m.order)
		return
	}
	dst.layout = nil
	dst.values = nil
	dst.present = 0
	dst.order = nil
	dst.fallback = jutil.NewObjectLinkedHashMap[*Var, *MustBeReachingVariableDef_Definition]()
	for v, d := range m.fallback.All() {
		dst.fallback.Put(v, d)
	}
}

func (m *mustDefMap) materialize() {
	f := m.fallback
	if f == nil {
		f = jutil.NewObjectLinkedHashMap[*Var, *MustBeReachingVariableDef_Definition]()
	} else {
		f.Clear()
	}
	for _, id := range m.order {
		f.Put(m.layout.vars[id], m.values[id])
	}
	if m.layout != nil && m.layout.wide {
		clear(m.presentWords[:m.layout.numWords])
	}
	m.fallback = f
	m.layout, m.values, m.order = nil, nil, nil
	m.present = 0
}

func (m *mustDefMap) Put(v *Var, d *MustBeReachingVariableDef_Definition) {
	if m.layout == nil {
		m.fallback.Put(v, d)
		return
	}
	id, ok := m.layout.ids[v]
	if !ok {
		m.materialize()
		m.fallback.Put(v, d)
		return
	}
	m.putID(id, d)
}

func (m *mustDefMap) putID(id uint16, d *MustBeReachingVariableDef_Definition) {
	if m.values == nil {
		n := len(m.layout.vars)
		if cap(m.values) >= n && cap(m.order) >= n {
			clear(m.values[:n])
			m.values = m.values[:n]
			m.order = m.order[:0]
		} else if n <= len(m.valuesBuf) {
			clear(m.valuesBuf[:n])
			m.values = m.valuesBuf[:n]
			m.order = m.orderBuf[:0:len(m.orderBuf)]
		} else {
			m.values = make([]*MustBeReachingVariableDef_Definition, n)
			m.order = make([]uint16, 0, n)
		}
	}
	if !m.layout.wide {
		bit := uint64(1) << id
		if m.present&bit == 0 {
			m.order = append(m.order, id)
			m.present |= bit
		}
	} else {
		w := id >> 6
		bit := uint64(1) << (id & 63)
		if m.presentWords[w]&bit == 0 {
			m.order = append(m.order, id)
			m.presentWords[w] |= bit
		}
	}
	m.values[id] = d
}

func (m *mustDefMap) invalidateDependent(v *Var) {
	if m.layout != nil {
		for _, id := range m.order {
			if otherDef := m.values[id]; otherDef != nil && !otherDef.depends.IsEmpty() && otherDef.depends.Contains(v) {
				m.values[id] = nil
			}
		}
		return
	}
	for key, otherDef := range m.fallback.All() {
		if otherDef != nil && !otherDef.depends.IsEmpty() && otherDef.depends.Contains(v) {
			m.fallback.Put(key, nil)
		}
	}
}

func (m *mustDefMap) invalidateParameters() {
	if m.layout != nil {
		for _, id := range m.order {
			value := m.values[id]
			if value == nil || value.depends.IsEmpty() {
				continue
			}
			for dep := range value.depends.All() {
				if MustBeReachingVariableDef_IsParameter(dep) {
					m.values[id] = nil
					break
				}
			}
		}
		return
	}
	for key, value := range m.fallback.All() {
		if value == nil || value.depends.IsEmpty() {
			continue
		}
		for dep := range value.depends.All() {
			if MustBeReachingVariableDef_IsParameter(dep) {
				m.fallback.Put(key, nil)
				break
			}
		}
	}
}

func (m *mustDefMap) lookup(v *Var) (*MustBeReachingVariableDef_Definition, bool) {
	if m.layout == nil {
		return m.fallback.Lookup(v)
	}
	if id, ok := m.layout.ids[v]; ok {
		if !m.hasID(id) {
			return nil, false
		}
		return m.values[id], true
	}
	// Preserve value equality when queried with another Var for the same
	// scoped name. Production transfers normally use canonical Vars.
	for _, id := range m.order {
		if jutil.ObjectsEquals(m.layout.vars[id], v) {
			return m.values[id], true
		}
	}
	return nil, false
}

func (m *mustDefMap) Get(v *Var) *MustBeReachingVariableDef_Definition {
	d, _ := m.lookup(v)
	return d
}

func (m *mustDefMap) ContainsKey(v *Var) bool {
	_, ok := m.lookup(v)
	return ok
}

func (m *mustDefMap) All() iter.Seq2[*Var, *MustBeReachingVariableDef_Definition] {
	if m.layout == nil {
		return m.fallback.All()
	}
	return func(yield func(*Var, *MustBeReachingVariableDef_Definition) bool) {
		for i := 0; i < len(m.order); i++ {
			id := m.order[i]
			if !yield(m.layout.vars[id], m.values[id]) {
				return
			}
		}
	}
}

func (m *mustDefMap) size() int {
	if m.layout != nil {
		return len(m.order)
	}
	return m.fallback.Size()
}

func (m *mustDefMap) equals(other *mustDefMap) bool {
	if m.size() != other.size() {
		return false
	}
	if m.layout != nil && m.layout == other.layout {
		if !m.layout.wide {
			if m.present != other.present {
				return false
			}
		} else {
			for i := 0; i < m.layout.numWords; i++ {
				if m.presentWords[i] != other.presentWords[i] {
					return false
				}
			}
		}
		for _, id := range m.order {
			a, b := m.values[id], other.values[id]
			if a != b && (a == nil || b == nil || a.node != b.node) {
				return false
			}
		}
		return true
	}
	if m.layout != nil {
		for _, id := range m.order {
			v, d := m.layout.vars[id], m.values[id]
			x, ok := other.lookup(v)
			if !ok || !jutil.ObjectsEquals(d, x) {
				return false
			}
		}
		return true
	}
	for v, d := range m.fallback.All() {
		x, ok := other.lookup(v)
		if !ok || !jutil.ObjectsEquals(d, x) {
			return false
		}
	}
	return true
}

func (m *mustDefMap) HashCode() int32 {
	var h int32
	if m.layout != nil {
		for _, id := range m.order {
			h += jutil.ObjectHashCode(m.layout.vars[id]) ^ jutil.ObjectHashCode(m.values[id])
		}
		return h
	}
	for v, d := range m.fallback.All() {
		h += jutil.ObjectHashCode(v) ^ jutil.ObjectHashCode(d)
	}
	return h
}
