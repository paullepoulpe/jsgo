// Port of com/google/javascript/rhino/testing/AbstractStaticScope.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/testing/Asserts.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/testing/MapBasedScope.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/testing/TemplateTypeMapSubject.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/testing/TestErrorReporter.java @ closure-compiler v20261021

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
 *   Nick Santos
 *   Google Inc.
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
	"slices"
	"strconv"
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

// Helper methods for making assertions about the validity of types.

// Java: Asserts.assertResolvesToSame(JSType)
func Asserts_AssertResolvesToSame(t testing.TB, type_ rhino.JSType) rhino.JSType {
	t.Helper()
	Truth{T: t}.AssertThat(Asserts_AssertValidResolve(t, type_)).IsSameInstanceAs(type_)
	return type_
}

// Returns the resolved type.
//
// Java: Asserts.assertValidResolve(JSType)
func Asserts_AssertValidResolve(t testing.TB, type_ rhino.JSType) rhino.JSType {
	t.Helper()
	reporter := NewTestErrorReporter()
	resolvedType := type_.Resolve(reporter)
	Truth{T: t}.AssertWithMessage("JSType#resolve should not affect object equality").
		AboutTypes().
		That(resolvedType).
		IsEqualTo(type_)
	reporter.VerifyHasEncounteredAllWarningsAndErrors(t)
	return resolvedType
}

// Java: Asserts.assertTypeCollectionEquals(Iterable<T>, Iterable<S>)
func Asserts_AssertTypeCollectionEquals[T rhino.JSType, S rhino.JSType](t testing.TB, a []T, b []S) {
	t.Helper()
	tr := Truth{T: t}
	tr.AssertThat(b).HasSize(len(a))
	for i := range a {
		tr.AssertType(b[i]).IsEqualTo(a[i])
	}
}

// For the given equivalent types, run all type operations that should have
// trivially predictable results. Proves that the two types are equivalent.
//
// Java: Asserts.assertEquivalenceOperations(JSType, JSType)
func Asserts_AssertEquivalenceOperations(t testing.TB, a rhino.JSType, b rhino.JSType) {
	t.Helper()
	tr := Truth{T: t}
	tr.AssertType(a).IsEqualTo(a)
	tr.AssertType(b).IsEqualTo(a)
	tr.AssertType(b).IsEqualTo(b)

	tr.AssertThat(a.IsSubtypeOf(b)).IsTrue()
	tr.AssertThat(a.IsSubtypeOf(a)).IsTrue()
	tr.AssertThat(b.IsSubtypeOf(b)).IsTrue()
	tr.AssertThat(b.IsSubtypeOf(a)).IsTrue()

	tr.AssertType(a.GetGreatestSubtype(b)).IsEqualTo(a)
	tr.AssertType(a.GetGreatestSubtype(a)).IsEqualTo(a)
	tr.AssertType(b.GetGreatestSubtype(b)).IsEqualTo(a)
	tr.AssertType(b.GetGreatestSubtype(a)).IsEqualTo(a)

	tr.AssertType(a.GetLeastSupertype(b)).IsEqualTo(a)
	tr.AssertType(a.GetLeastSupertype(a)).IsEqualTo(a)
	tr.AssertType(b.GetLeastSupertype(b)).IsEqualTo(a)
	tr.AssertType(b.GetLeastSupertype(a)).IsEqualTo(a)

	tr.AssertThat(a.CanCastTo(b)).IsTrue()
	tr.AssertThat(a.CanCastTo(a)).IsTrue()
	tr.AssertThat(b.CanCastTo(b)).IsTrue()
	tr.AssertThat(b.CanCastTo(a)).IsTrue()
}

// TestErrorReporter is an error reporter for testing that verifies that
// messages reported to the reporter are expected.
//
// Use expectAllErrors/expectAllWarnings before the code that reports them,
// and VerifyHasEncounteredAllWarningsAndErrors to check that all were seen.
type TestErrorReporter struct {
	expectedErrors   []string
	expectedWarnings []string
	// Seen diagnostics are stored in a LinkedHashMap to deduplicate
	// identical diagnostics. The default
	// com.google.javascript.jscomp.SortingErrorManager behaves in this way.
	seenErrors   *jutil.LinkedHashMap[string, string]
	seenWarnings *jutil.LinkedHashMap[string, string]
}

// NewTestErrorReporter returns a reporter that expects no diagnostics.
//
// Java: TestErrorReporter()
func NewTestErrorReporter() *TestErrorReporter {
	return &TestErrorReporter{
		expectedErrors:   []string{},
		expectedWarnings: []string{},
		seenErrors:       jutil.NewLinkedHashMap[string, string](),
		seenWarnings:     jutil.NewLinkedHashMap[string, string](),
	}
}

// Java: TestErrorReporter.error(String, String, int, int)
func (r *TestErrorReporter) Error(message string, sourceName string, line int, lineOffset int) {
	r.seenErrors.Put(r.fmtDiagnosticKey(message, sourceName, line, lineOffset), message)
}

// Java: TestErrorReporter.warning(String, String, int, int)
func (r *TestErrorReporter) Warning(message string, sourceName string, line int, lineOffset int) {
	r.seenWarnings.Put(r.fmtDiagnosticKey(message, sourceName, line, lineOffset), message)
}

// Java: TestErrorReporter.expectAllErrors(String...)
func (r *TestErrorReporter) ExpectAllErrors(errors ...string) *TestErrorReporter {
	r.expectedErrors = append(r.expectedErrors, errors...)
	return r
}

// Java: TestErrorReporter.expectAllWarnings(String...)
func (r *TestErrorReporter) ExpectAllWarnings(warnings ...string) *TestErrorReporter {
	r.expectedWarnings = append(r.expectedWarnings, warnings...)
	return r
}

// VerifyHasEncounteredAllWarningsAndErrors checks that exactly the expected
// diagnostics were reported, in order. Go-only: language-adapter — takes
// the test that Truth reports to.
//
// Java: TestErrorReporter.verifyHasEncounteredAllWarningsAndErrors()
func (r *TestErrorReporter) VerifyHasEncounteredAllWarningsAndErrors(t testing.TB) {
	t.Helper()
	tr := Truth{T: t}
	tr.AssertThat(slices.Collect(r.seenWarnings.Values())).ContainsExactlyElementsIn(r.expectedWarnings).InOrder()
	tr.AssertThat(slices.Collect(r.seenErrors.Values())).ContainsExactlyElementsIn(r.expectedErrors).InOrder()
}

// Java: TestErrorReporter.fmtDiagnosticKey(String, String, int, int)
func (r *TestErrorReporter) fmtDiagnosticKey(message string, sourceName string, line int, lineOffset int) string {
	return message + ":" + sourceName + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(lineOffset)
}
