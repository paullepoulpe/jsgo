// Port of com/google/javascript/rhino/testing/MapBasedScope.java @ closure-compiler v20261015

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Nick Santos
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

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
