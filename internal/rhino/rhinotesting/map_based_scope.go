// Port of com/google/javascript/rhino/testing/MapBasedScope.java @ closure-compiler v20260415

package rhinotesting

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// MapBasedScope is a scope based on a simple map of names to types.
type MapBasedScope struct {
	AbstractStaticScopeBase
	slots *jutil.LinkedHashMap[string, rhino.StaticTypedSlot]
}

// NewMapBasedScope returns a scope with one non-inferred slot per entry of
// namesToTypes, in order.
//
// Java: MapBasedScope(Map<String, ? extends JSType>)
func NewMapBasedScope(namesToTypes *jutil.LinkedHashMap[string, rhino.JSType]) *MapBasedScope {
	s := &MapBasedScope{slots: jutil.NewLinkedHashMap[string, rhino.StaticTypedSlot]()}
	s.InitAbstractStaticScopeBase(s)
	if namesToTypes != nil {
		for name, type_ := range namesToTypes.All() {
			s.slots.Put(name, rhino.NewSimpleSlot(name, type_, false))
		}
	}
	return s
}

// NewMapBasedScopeOf builds a MapBasedScope from alternating names and
// types, like ImmutableMap.of(k1, v1, k2, v2, ...).
func NewMapBasedScopeOf(namesAndTypes ...any) *MapBasedScope {
	m := jutil.NewLinkedHashMap[string, rhino.JSType]()
	for i := 0; i < len(namesAndTypes); i += 2 {
		m.Put(namesAndTypes[i].(string), namesAndTypes[i+1].(rhino.JSType))
	}
	return NewMapBasedScope(m)
}

// Java: MapBasedScope.emptyScope()
func MapBasedScope_EmptyScope() *MapBasedScope {
	return NewMapBasedScope(jutil.NewLinkedHashMap[string, rhino.JSType]())
}

// Java: MapBasedScope.getSlot(String)
func (s *MapBasedScope) GetSlot(name string) rhino.StaticSlot {
	if slot := s.slots.Get(name); slot != nil {
		return slot
	}
	return nil
}
