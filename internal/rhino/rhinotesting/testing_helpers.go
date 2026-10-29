// Port of com/google/javascript/rhino/testing/AbstractStaticScope.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/testing/MapBasedScope.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/testing/TemplateTypeMapSubject.java @ closure-compiler v20261021

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
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhinotesting

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AbstractStaticScopeBase holds default StaticTypedScope methods for test
// scopes; the embedding test scope provides GetSlot.
type AbstractStaticScopeBase struct {
	self rhino.StaticTypedScope
}

// InitAbstractStaticScopeBase sets the most-derived receiver.
func (s *AbstractStaticScopeBase) InitAbstractStaticScopeBase(self rhino.StaticTypedScope) {
	s.self = self
}

// Java: AbstractStaticScope.getRootNode()
func (s *AbstractStaticScopeBase) GetRootNode() rhino.Node {
	return rhino.NilNode
}

// Java: AbstractStaticScope.getParentScope()
func (s *AbstractStaticScopeBase) GetParentScope() rhino.StaticScope {
	return nil
}

// Java: AbstractStaticScope.getOwnSlot(String)
func (s *AbstractStaticScopeBase) GetOwnSlot(name string) rhino.StaticSlot {
	return s.self.GetSlot(name)
}

// Java: AbstractStaticScope.getTypeOfThis()
func (s *AbstractStaticScopeBase) GetTypeOfThis() rhino.JSType {
	return nil
}

// Java: StaticScope.getTopmostScopeOfEventualDeclaration(String) (default method)
func (s *AbstractStaticScopeBase) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	return rhino.StaticScope_GetTopmostScopeOfEventualDeclaration(s.self, name)
}

// Java: StaticTypedScope.lookupQualifiedName(QualifiedName) (default method)
func (s *AbstractStaticScopeBase) LookupQualifiedName(qname rhino.QualifiedName) rhino.JSType {
	return rhino.StaticTypedScope_LookupQualifiedName(s.self, qname)
}

// MapBasedScope is a scope based on a simple map of names to types.
type MapBasedScope struct {
	slots *jutil.LinkedHashMap[string, rhino.StaticTypedSlot]
}

// NewMapBasedScope returns a scope with one non-inferred slot per entry of
// namesToTypes, in order.
//
// Java: MapBasedScope(Map<String, ? extends JSType>)
func NewMapBasedScope(namesToTypes *jutil.LinkedHashMap[string, rhino.JSType]) *MapBasedScope {
	s := &MapBasedScope{slots: jutil.NewLinkedHashMap[string, rhino.StaticTypedSlot]()}
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

// Java: AbstractStaticScope.getRootNode()
func (s *MapBasedScope) GetRootNode() rhino.Node { return rhino.NilNode }

// Java: AbstractStaticScope.getParentScope()
func (s *MapBasedScope) GetParentScope() rhino.StaticScope { return nil }

// Java: AbstractStaticScope.getTypeOfThis()
func (s *MapBasedScope) GetTypeOfThis() rhino.JSType { return nil }

// Java: MapBasedScope.getSlot(String)
func (s *MapBasedScope) GetSlot(name string) rhino.StaticSlot {
	if slot := s.slots.Get(name); slot != nil {
		return slot
	}
	return nil
}

// Java: AbstractStaticScope.getOwnSlot(String)
func (s *MapBasedScope) GetOwnSlot(name string) rhino.StaticSlot {
	return s.GetSlot(name)
}

// Java: StaticScope.getTopmostScopeOfEventualDeclaration(String) (default method)
func (s *MapBasedScope) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	return rhino.StaticScope_GetTopmostScopeOfEventualDeclaration(s, name)
}

// Java: StaticTypedScope.lookupQualifiedName(QualifiedName) (default method)
func (s *MapBasedScope) LookupQualifiedName(qname rhino.QualifiedName) rhino.JSType {
	return rhino.StaticTypedScope_LookupQualifiedName(s, qname)
}

// EqualsTester checks equals/hashCode: items in the same group are equal
// with equal hash codes, items in different groups are not equal, and no
// item equals null or an instance of an unrelated class.
//
// Go-only: language-adapter — Guava's com.google.common.testing.EqualsTester,
// used by the ported tests.
type EqualsTester struct {
	groups [][]any
}

// NewEqualsTester returns an empty tester.
func NewEqualsTester() *EqualsTester {
	return &EqualsTester{}
}

// AddEqualityGroup adds a group of mutually equal items.
func (e *EqualsTester) AddEqualityGroup(items ...any) *EqualsTester {
	e.groups = append(e.groups, items)
	return e
}

type equalsTesterNotAnInstance struct{}

// TestEquals runs the checks.
func (e *EqualsTester) TestEquals(t testing.TB) *EqualsTester {
	t.Helper()
	for gi, group := range e.groups {
		for _, item := range group {
			if jutil.ObjectsEquals(item, nil) {
				t.Fatalf("%s must not be Object#equals to null", Show(item))
			}
			if jutil.ObjectsEquals(item, &equalsTesterNotAnInstance{}) {
				t.Fatalf("%s must not be Object#equals to an arbitrary object of another class", Show(item))
			}
			for gj, other := range e.groups {
				for _, related := range other {
					eq := jutil.ObjectsEquals(item, related)
					if gi == gj {
						if !eq {
							t.Fatalf("%s [group %d] must be Object#equals to %s [group %d]", Show(item), gi+1, Show(related), gj+1)
						}
						if jutil.ObjectHashCode(item) != jutil.ObjectHashCode(related) {
							t.Fatalf("the Object#hashCode of %s must be equal to the Object#hashCode of %s", Show(item), Show(related))
						}
					} else if eq {
						t.Fatalf("%s [group %d] must not be Object#equals to %s [group %d]", Show(item), gi+1, Show(related), gj+1)
					}
				}
			}
		}
	}
	return e
}

// TemplateTypeMapSubject is a Truth Subject for TemplateTypeMap.
type TemplateTypeMapSubject struct {
	t      testing.TB
	actual *rhino.TemplateTypeMap
}

// Java: TemplateTypeMapSubject.assertThat(TemplateTypeMap)
func TemplateTypeMapSubject_AssertThat(t testing.TB, actual *rhino.TemplateTypeMap) *TemplateTypeMapSubject {
	return &TemplateTypeMapSubject{t: t, actual: actual}
}

// Java: TemplateTypeMapSubject.hasKeysAndValues(ImmutableList<TemplateType>, ImmutableList<? extends JSType>)
func (s *TemplateTypeMapSubject) HasKeysAndValues(keys []*rhino.TemplateType, values []rhino.JSType) {
	s.t.Helper()
	tr := Truth{T: s.t}
	tr.AssertWithMessage("value of: templateTypeMap.getTemplateKeys()").That(s.actual.GetTemplateKeys()).IsEqualTo(keys)
	tr.AssertWithMessage("value of: templateTypeMap.getTemplateValues()").That(s.actual.GetTemplateValues()).IsEqualTo(values)
}

// Java: TemplateTypeMapSubject.hasSubmapsAt(int...)
func (s *TemplateTypeMapSubject) HasSubmapsAt(indices ...int) {
	s.t.Helper()
	Truth{T: s.t}.AssertWithMessage("value of: templateTypeMap.getIndicesOfSubmapsForTesting()").
		That(s.actual.GetIndicesOfSubmapsForTesting()).
		ContainsExactlyElementsIn(indices)
}
