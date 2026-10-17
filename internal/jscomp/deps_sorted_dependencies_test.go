// Port of test/com/google/javascript/jscomp/deps/SortedDependenciesTest.java @ closure-compiler v20260415

/*
 * Copyright 2015 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"slices"
	"testing"
)

func TestSortedDependencies_Sort(t *testing.T) {
	a := SimpleDependencyInfo_Builder_("a", "a").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("b"), DependencyInfo_Require_GoogRequireSymbol("c")).Build()
	b := SimpleDependencyInfo_Builder_("b", "b").SetProvidesArray("b").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("d")).Build()
	c := SimpleDependencyInfo_Builder_("c", "c").SetProvidesArray("c").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("d")).Build()
	d := SimpleDependencyInfo_Builder_("d", "d").SetProvidesArray("d").Build()

	jsaAssertSortedInputs(t, jsaSdis(d, b, c, a), jsaSdis(a, b, c, d))
	jsaAssertSortedInputs(t, jsaSdis(d, b, c, a), jsaSdis(d, b, c, a))
	jsaAssertSortedInputs(t, jsaSdis(d, c, b, a), jsaSdis(d, c, b, a))
	jsaAssertSortedInputs(t, jsaSdis(d, b, c, a), jsaSdis(d, a, b, c))

	jsaAssertSortedDeps(t, jsaSdis(d, b, c, a), jsaSdis(d, b, c, a), jsaSdis(a))
	jsaAssertSortedDeps(t, jsaSdis(d, c), jsaSdis(d, c, b, a), jsaSdis(c))
	jsaAssertSortedDeps(t, jsaSdis(d), jsaSdis(d, c, b, a), jsaSdis(d))
}

func TestSortedDependencies_Sort2(t *testing.T) {
	ab := SimpleDependencyInfo_Builder_("ab", "ab").SetProvidesArray("a", "b").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("d"), DependencyInfo_Require_GoogRequireSymbol("f")).Build()
	c := SimpleDependencyInfo_Builder_("c", "c").SetProvidesArray("c").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("h")).Build()
	d := SimpleDependencyInfo_Builder_("d", "d").SetProvidesArray("d").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("e"), DependencyInfo_Require_GoogRequireSymbol("f")).Build()
	ef := SimpleDependencyInfo_Builder_("ef", "ef").SetProvidesArray("e", "f").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("g"), DependencyInfo_Require_GoogRequireSymbol("c")).Build()
	g := SimpleDependencyInfo_Builder_("g", "g").SetProvidesArray("g").Build()
	hi := SimpleDependencyInfo_Builder_("hi", "hi").SetProvidesArray("h", "i").Build()

	jsaAssertSortedInputs(t, jsaSdis(g, hi, c, ef, d, ab), jsaSdis(ab, c, d, ef, g, hi))

	jsaAssertSortedDeps(t, jsaSdis(g), jsaSdis(ab, c, d, ef, g, hi), jsaSdis(g))
	jsaAssertSortedDeps(t, jsaSdis(g, hi, c, ef, d), jsaSdis(ab, c, d, ef, g, hi), jsaSdis(d, hi))
}

func TestSortedDependencies_Sort3(t *testing.T) {
	a := SimpleDependencyInfo_Builder_("a", "a").SetProvidesArray("a").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("c")).Build()
	b := SimpleDependencyInfo_Builder_("b", "b").SetProvidesArray("b").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("a")).Build()
	c := SimpleDependencyInfo_Builder_("c", "c").SetProvidesArray("c").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("b")).Build()

	jsaAssertOrder(t, jsaSdis(a, b, c), jsaSdis(b, c, a))
}

func TestSortedDependencies_Sort4(t *testing.T) {
	// Check the degenerate case.
	a := SimpleDependencyInfo_Builder_("a", "a").SetProvidesArray("a").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("a")).Build()
	jsaAssertSortedDeps(t, jsaSdis(a), jsaSdis(a), jsaSdis(a))
}

func TestSortedDependencies_Sort5(t *testing.T) {
	a := SimpleDependencyInfo_Builder_("a", "a").SetProvidesArray("a").Build()
	b := SimpleDependencyInfo_Builder_("b", "b").SetProvidesArray("b").Build()
	c := SimpleDependencyInfo_Builder_("c", "c").SetProvidesArray("c").Build()

	jsaAssertSortedInputs(t, jsaSdis(a, b, c), jsaSdis(a, b, c))
	jsaAssertSortedInputs(t, jsaSdis(c, b, a), jsaSdis(c, b, a))
}

func TestSortedDependencies_Sort6(t *testing.T) {
	a := SimpleDependencyInfo_Builder_("gin", "gin").SetProvidesArray("gin").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("tonic")).Build()
	b := SimpleDependencyInfo_Builder_("tonic", "tonic").SetProvidesArray("tonic").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("gin2")).Build()
	c := SimpleDependencyInfo_Builder_("gin2", "gin2").SetProvidesArray("gin2").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("gin")).Build()
	d := SimpleDependencyInfo_Builder_("gin3", "gin3").SetProvidesArray("gin3").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("gin")).Build()

	jsaAssertOrder(t, jsaSdis(a, b, c, d), jsaSdis(c, b, a, d))
}

func TestSortedDependencies_Sort7(t *testing.T) {
	a := SimpleDependencyInfo_Builder_("gin", "gin").SetProvidesArray("gin").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("tonic")).Build()
	b := SimpleDependencyInfo_Builder_("tonic", "tonic").SetProvidesArray("tonic").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("gin")).Build()
	c := SimpleDependencyInfo_Builder_("gin2", "gin2").SetProvidesArray("gin2").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("gin")).Build()
	d := SimpleDependencyInfo_Builder_("gin3", "gin3").SetProvidesArray("gin3").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("gin")).Build()

	jsaAssertOrder(t, jsaSdis(a, b, c, d), jsaSdis(b, a, c, d))
}

func TestSortedDependencies_Sort8(t *testing.T) {
	a := SimpleDependencyInfo_Builder_("A", "A").SetProvidesArray("A").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("B")).Build()
	b := SimpleDependencyInfo_Builder_("B", "B").SetProvidesArray("B").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("C")).Build()
	c := SimpleDependencyInfo_Builder_("C", "C").SetProvidesArray("C").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("D")).Build()
	d := SimpleDependencyInfo_Builder_("D", "D").SetProvidesArray("D").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("A")).Build()

	jsaAssertOrder(t, jsaSdis(a, b, c, d), jsaSdis(d, c, b, a))
}

func TestSortedDependencies_Sort9(t *testing.T) {
	a := SimpleDependencyInfo_Builder_("A", "A").SetProvidesArray("A").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("B")).Build()
	a2 := SimpleDependencyInfo_Builder_("A", "A").SetProvidesArray("A").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("B1")).Build()
	b := SimpleDependencyInfo_Builder_("B", "B").SetProvidesArray("B").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("C")).Build()
	c := SimpleDependencyInfo_Builder_("C", "C").SetProvidesArray("C").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("E")).Build()
	d := SimpleDependencyInfo_Builder_("D", "D").SetProvidesArray("D").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("A")).Build()
	e := SimpleDependencyInfo_Builder_("B1", "B1").SetProvidesArray("B1").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("C1")).Build()
	f := SimpleDependencyInfo_Builder_("C1", "C1").SetProvidesArray("C1").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("D1")).Build()
	g := SimpleDependencyInfo_Builder_("D1", "D1").SetProvidesArray("D1").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("A")).Build()

	jsaAssertOrder(t, jsaSdis(a, a2, b, c, d, e, f, g), jsaSdis(c, b, a, g, f, e, a2, d))
}

func TestSortedDependencies_Sort10(t *testing.T) {
	a := SimpleDependencyInfo_Builder_("A", "A").SetProvidesArray("A").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("C")).Build()
	b := SimpleDependencyInfo_Builder_("B", "B").SetProvidesArray("B").Build()
	c := SimpleDependencyInfo_Builder_("C", "C").SetProvidesArray("C").Build()

	sorted := NewSortedDependencies(jsaSdis(a, b, c))

	jsaDeepEq(t, sorted.GetSortedList(), jsaSdis(c, a, b))
}

func TestSortedDependencies_WeakSort(t *testing.T) {
	a := SimpleDependencyInfo_Builder_("a", "a").SetProvidesArray("a").SetTypeRequiresArray("b").Build()
	b := SimpleDependencyInfo_Builder_("b", "b").SetProvidesArray("b").SetTypeRequiresArray("c").Build()
	c := SimpleDependencyInfo_Builder_("c", "c").SetProvidesArray("c").SetTypeRequiresArray("d").Build()
	d := SimpleDependencyInfo_Builder_("d", "d").SetProvidesArray("d").Build()

	// The order of weak edges is based on user input.
	jsaAssertSortedWeakDeps(t, jsaSdis(b, c, d), jsaSdis(a, b, c, d), jsaSdis(a))

	jsaAssertSortedWeakDeps(t, jsaSdis(c, b, d), jsaSdis(a, c, b, d), jsaSdis(a))

	jsaAssertSortedWeakDeps(t, jsaSdis(d, b, c), jsaSdis(d, b, c, a), jsaSdis(a))

	jsaAssertSortedWeakDeps(t, jsaSdis(c, d, b), jsaSdis(c, d, b, a), jsaSdis(a))
}

func TestSortedDependencies_WeakSortIncludesStrongEdgesFromWeakSources(t *testing.T) {
	a := SimpleDependencyInfo_Builder_("a", "a").SetProvidesArray("a").SetTypeRequiresArray("b").Build()
	b := SimpleDependencyInfo_Builder_("b", "b").SetProvidesArray("b").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("c")).Build()
	c := SimpleDependencyInfo_Builder_("c", "c").SetProvidesArray("c").Build()

	// Order is input order invariant due to the strong edge.
	jsaAssertSortedWeakDeps(t, jsaSdis(c, b), jsaSdis(a, b, c), jsaSdis(a))

	jsaAssertSortedWeakDeps(t, jsaSdis(c, b), jsaSdis(a, c, b), jsaSdis(a))
}

func TestSortedDependencies_WeakAndStrongIsStrong(t *testing.T) {
	a := SimpleDependencyInfo_Builder_("a", "a").SetProvidesArray("a").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("c")).Build()
	b := SimpleDependencyInfo_Builder_("b", "b").SetProvidesArray("b").SetTypeRequiresArray("c").Build()
	c := SimpleDependencyInfo_Builder_("c", "c").SetProvidesArray("c").Build()

	jsaAssertSortedWeakDeps(t, jsaSdis(), jsaSdis(a, b, c), jsaSdis(a, b))

	jsaAssertSortedWeakDeps(t, jsaSdis(), jsaSdis(a, b, c), jsaSdis(b, a))

	jsaAssertSortedWeakDeps(t, jsaSdis(), jsaSdis(c, b, a), jsaSdis(a, b))

	jsaAssertSortedWeakDeps(t, jsaSdis(), jsaSdis(c, a, b), jsaSdis(b, a))
}

func TestSortedDependencies_SortCircularWeakStrongReference(t *testing.T) {
	a := SimpleDependencyInfo_Builder_("a", "a").SetProvidesArray("a").SetRequiresArray(DependencyInfo_Require_GoogRequireSymbol("b")).Build()
	b := SimpleDependencyInfo_Builder_("b", "b").SetProvidesArray("b").SetTypeRequiresArray("a").Build()

	jsaAssertSortedInputs(t, jsaSdis(b, a), jsaSdis(a, b))
	jsaAssertSortedInputs(t, jsaSdis(b, a), jsaSdis(b, a))
}

func TestSortedDependencies_SortCircularWeakReference(t *testing.T) {
	a := SimpleDependencyInfo_Builder_("a", "a").SetProvidesArray("a").SetTypeRequiresArray("b").Build()
	b := SimpleDependencyInfo_Builder_("b", "b").SetProvidesArray("b").SetTypeRequiresArray("c").Build()
	c := SimpleDependencyInfo_Builder_("c", "c").SetProvidesArray("c").SetTypeRequiresArray("b").Build()

	jsaAssertSortedWeakDeps(t, jsaSdis(b, c), jsaSdis(b, c, a), jsaSdis(a))
	jsaAssertSortedWeakDeps(t, jsaSdis(c, b), jsaSdis(c, b, a), jsaSdis(a))
	jsaAssertSortedWeakDeps(t, jsaSdis(b, c), jsaSdis(a, b, c), jsaSdis(a))
}

func jsaSdis(infos ...*SimpleDependencyInfo) []*SimpleDependencyInfo { return infos }

func jsaSdiNames(infos []*SimpleDependencyInfo) []string {
	var names []string
	for _, i := range infos {
		names = append(names, i.GetName())
	}
	return names
}

func jsaSdiCheck(t *testing.T, got, expected []*SimpleDependencyInfo) {
	t.Helper()
	if !slices.Equal(got, expected) {
		t.Errorf("got %v, want %v", jsaSdiNames(got), jsaSdiNames(expected))
	}
}

func jsaAssertSortedInputs(t *testing.T, expected, shuffled []*SimpleDependencyInfo) {
	t.Helper()
	sorted := NewSortedDependencies(shuffled)
	jsaSdiCheck(t, sorted.GetSortedList(), expected)
}

func jsaAssertSortedDeps(t *testing.T, expected, shuffled, roots []*SimpleDependencyInfo) {
	t.Helper()
	sorted := NewSortedDependencies(shuffled)
	jsaSdiCheck(t, sorted.GetSortedStrongDependenciesOf(roots), expected)
}

func jsaAssertSortedWeakDeps(t *testing.T, expected, shuffled, roots []*SimpleDependencyInfo) {
	t.Helper()
	sorted := NewSortedDependencies(shuffled)
	jsaSdiCheck(t, sorted.GetSortedWeakDependenciesOf(roots), expected)
}

func jsaAssertOrder(t *testing.T, shuffle, expected []*SimpleDependencyInfo) {
	t.Helper()
	sorted := NewSortedDependencies(shuffle)
	jsaSdiCheck(t, sorted.GetSortedList(), expected)
}
