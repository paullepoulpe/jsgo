package jscomp

import (
	"iter"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Go-only: performance — a function's canonical Vars have compact local IDs.
// Larger functions retain the sparse representation to bound lattice memory.
// The fallback also preserves ScopedName equality for non-canonical Vars.
type mustDefLayout struct {
	vars []*Var
	ids  map[*Var]uint8
}

var emptyMustDefLayout = &mustDefLayout{}

func newMustDefLayout(vars *jutil.LinkedHashMap[string, *Var]) *mustDefLayout {
	if vars == nil || vars.Size() > 64 {
		return nil
	}
	if vars.Size() == 0 {
		return emptyMustDefLayout
	}
	l := &mustDefLayout{
		vars: make([]*Var, 0, vars.Size()),
		ids:  make(map[*Var]uint8, vars.Size()),
	}
	for v := range vars.Values() {
		l.ids[v] = uint8(len(l.vars))
		l.vars = append(l.vars, v)
	}
	return l
}

// mustDefMap represents the same insertion-ordered partial map as the Java
// lattice. A clear presence bit is TOP; a present nil value is BOTTOM. The
// order slice preserves encounter order independently of the variable IDs.
type mustDefMap struct {
	layout    *mustDefLayout
	values    []*MustBeReachingVariableDef_Definition
	present   uint64
	order     []uint8
	valuesBuf [8]*MustBeReachingVariableDef_Definition
	orderBuf  [8]uint8
	fallback  *jutil.ObjectLinkedHashMap[*Var, *MustBeReachingVariableDef_Definition]
}

func newMustDefMap(layout *mustDefLayout) mustDefMap {
	m := mustDefMap{layout: layout}
	if layout == nil {
		m.fallback = jutil.NewObjectLinkedHashMap[*Var, *MustBeReachingVariableDef_Definition]()
	}
	return m
}

func (m *mustDefMap) init(layout *mustDefLayout) {
	m.layout = layout
	m.values = nil
	m.present = 0
	m.order = nil
	if layout == nil {
		m.fallback = jutil.NewObjectLinkedHashMap[*Var, *MustBeReachingVariableDef_Definition]()
	} else {
		m.fallback = nil
	}
}

func (m *mustDefMap) clone() mustDefMap {
	if m.layout != nil {
		return mustDefMap{layout: m.layout, values: slices.Clone(m.values),
			present: m.present, order: slices.Clone(m.order)}
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
		dst.present = m.present
		dst.fallback = nil
		if m.values == nil {
			dst.values = nil
			dst.order = nil
			return
		}
		if len(m.values) <= len(dst.valuesBuf) {
			dst.values = dst.valuesBuf[:len(m.values)]
		} else {
			dst.values = make([]*MustBeReachingVariableDef_Definition, len(m.values))
		}
		copy(dst.values, m.values)
		if len(m.values) <= len(dst.orderBuf) {
			dst.order = dst.orderBuf[:len(m.order):len(dst.orderBuf)]
		} else {
			dst.order = make([]uint8, len(m.order), len(m.values))
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
	f := jutil.NewObjectLinkedHashMap[*Var, *MustBeReachingVariableDef_Definition]()
	for _, id := range m.order {
		f.Put(m.layout.vars[id], m.values[id])
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

func (m *mustDefMap) putID(id uint8, d *MustBeReachingVariableDef_Definition) {
	if m.values == nil {
		n := len(m.layout.vars)
		if n <= len(m.valuesBuf) {
			clear(m.valuesBuf[:n])
			m.values = m.valuesBuf[:n]
			m.order = m.orderBuf[:0:len(m.orderBuf)]
		} else {
			m.values = make([]*MustBeReachingVariableDef_Definition, n)
			m.order = make([]uint8, 0, n)
		}
	}
	bit := uint64(1) << id
	if m.present&bit == 0 {
		m.order = append(m.order, id)
		m.present |= bit
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
		if m.present&(uint64(1)<<id) == 0 {
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
		if m.present != other.present {
			return false
		}
		for _, id := range m.order {
			a, b := m.values[id], other.values[id]
			if a != b && (a == nil || b == nil || a.node != b.node) {
				return false
			}
		}
		return true
	}
	for v, d := range m.All() {
		x, ok := other.lookup(v)
		if !ok || !jutil.ObjectsEquals(d, x) {
			return false
		}
	}
	return true
}

func (m *mustDefMap) HashCode() int32 {
	var h int32
	for v, d := range m.All() {
		h += jutil.ObjectHashCode(v) ^ jutil.ObjectHashCode(d)
	}
	return h
}
