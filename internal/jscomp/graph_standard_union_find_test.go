// Port of com/google/javascript/jscomp/graph/StandardUnionFindTest.java @ closure-compiler v20261015

/*
 * Copyright 2008 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func setToSlice[E any](s jutil.Set[E]) []E {
	return slices.Collect(s.All())
}

func TestStandardUnionFind_Empty(t *testing.T) {
	union := NewStandardUnionFind[string]()
	assertEqual(t, "classes", len(union.AllEquivalenceClasses()), 0)
}

func TestStandardUnionFind_Add(t *testing.T) {
	union := NewStandardUnionFind[string]()
	union.Add("foo")
	union.Add("bar")
	union.Find("foo")
	assertEqual(t, "classes", len(union.AllEquivalenceClasses()), 2)
}

func TestStandardUnionFind_Union(t *testing.T) {
	union := NewStandardUnionFind[string]()
	union.Union("A", "B")
	union.Union("C", "D")
	assertEqual(t, "B", union.Find("B"), union.Find("A"))
	assertEqual(t, "D", union.Find("D"), union.Find("C"))
	assertEqual(t, "A != D", union.Find("A") == union.Find("D"), false)
}

func TestStandardUnionFind_SetSize(t *testing.T) {
	union := NewStandardUnionFind[string]()
	union.Union("A", "B")
	union.Union("B", "C")
	union.Union("D", "E")
	union.Union("F", "F")

	assertEqual(t, "A", union.FindAll("A").Size(), 3)
	assertEqual(t, "B", union.FindAll("B").Size(), 3)
	assertEqual(t, "C", union.FindAll("C").Size(), 3)
	assertEqual(t, "D", union.FindAll("D").Size(), 2)
	assertEqual(t, "F", union.FindAll("F").Size(), 1)
}

func TestStandardUnionFind_Find(t *testing.T) {
	union := NewStandardUnionFind[string]()
	union.Add("A")
	union.Add("B")
	assertEqual(t, "A", union.Find("A"), "A")
	assertEqual(t, "B", union.Find("B"), "B")

	union.Union("A", "B")
	assertEqual(t, "B", union.Find("B"), union.Find("A"))

	assertThrows[jutil.IIllegalArgumentException](t, func() { union.Find("Z") })
}

func TestStandardUnionFind_AllEquivalenceClasses(t *testing.T) {
	union := NewStandardUnionFind[string]()
	union.Union("A", "B")
	union.Union("A", "B")
	union.Union("B", "A")
	union.Union("B", "C")
	union.Union("D", "E")
	union.Union("F", "F")

	classes := union.AllEquivalenceClasses()
	var got []string
	for _, c := range classes {
		got = append(got, c.String())
	}
	assertSameElements(t, "classes", got, []string{"[A, B, C]", "[D, E]", "[F]"})
}

func TestStandardUnionFind_FindAll(t *testing.T) {
	union := NewStandardUnionFind[string]()
	union.Union("A", "B")
	union.Union("A", "B")
	union.Union("B", "A")
	union.Union("D", "E")
	union.Union("F", "F")

	aSet := union.FindAll("A")
	assertSameElements(t, "aSet", setToSlice(aSet), []string{"A", "B"})

	union.Union("B", "C")
	assertEqual(t, "contains C", aSet.Contains("C"), true)
	assertEqual(t, "size", aSet.Size(), 3)

	assertThrows[jutil.IIllegalArgumentException](t, func() { union.FindAll("Z") })
}

func TestStandardUnionFind_FindAllIterator(t *testing.T) {
	union := NewStandardUnionFind[string]()
	union.Union("A", "B")
	union.Union("B", "C")
	union.Union("A", "B")
	union.Union("D", "E")

	aSet := union.FindAll("A")
	assertSliceEqual(t, "aSet", setToSlice(aSet), []string{"A", "B", "C"})

	dSet := union.FindAll("D")
	assertSliceEqual(t, "dSet", setToSlice(dSet), []string{"D", "E"})
}

func TestStandardUnionFind_FindAllSize(t *testing.T) {
	union := NewStandardUnionFind[string]()
	union.Union("A", "B")
	union.Union("B", "C")
	assertEqual(t, "A", union.FindAll("A").Size(), 3)
	assertEqual(t, "B", union.FindAll("B").Size(), 3)
	assertEqual(t, "C", union.FindAll("C").Size(), 3)
	union.Union("D", "E")
	assertEqual(t, "C", union.FindAll("C").Size(), 3)
	assertEqual(t, "D", union.FindAll("D").Size(), 2)
	union.Union("B", "E")
	assertEqual(t, "C", union.FindAll("C").Size(), 5)
	assertEqual(t, "D", union.FindAll("D").Size(), 5)
}

func TestStandardUnionFind_Elements(t *testing.T) {
	union := NewStandardUnionFind[string]()
	union.Union("A", "B")
	union.Union("B", "C")
	union.Union("A", "B")
	union.Union("D", "E")

	elements := union.Elements()
	assertSliceEqual(t, "elements", setToSlice(elements), []string{"A", "B", "C", "D", "E"})
	assertEqual(t, "no F", elements.Contains("F"), false)
}

func TestStandardUnionFind_Copy(t *testing.T) {
	union := NewStandardUnionFind[string]()
	union.Union("A", "B")
	union.Union("B", "Z")
	union.Union("X", "Y")
	var copy UnionFind[string] = NewStandardUnionFindCopyOf[string](union)
	assertSameElements(t, "Z", setToSlice(copy.FindAll("Z")), []string{"A", "B", "Z"})
	assertSameElements(t, "X", setToSlice(copy.FindAll("X")), []string{"X", "Y"})
}

func TestStandardUnionFind_ChangesToCopyDontAffectOriginal(t *testing.T) {
	union := NewStandardUnionFind[string]()
	union.Union("A", "B")
	union.Union("X", "Y")
	union.Union("A", "C")
	var copy UnionFind[string] = NewStandardUnionFindCopyOf[string](union)
	copy.Union("A", "D")
	assertSameElements(t, "copy D", setToSlice(copy.FindAll("D")), []string{"A", "B", "C", "D"})
	assertSameElements(t, "union A", setToSlice(union.FindAll("A")), []string{"A", "B", "C"})
	assertSameElements(t, "copy X", setToSlice(copy.FindAll("X")), []string{"X", "Y"})
	assertThrows[jutil.IIllegalArgumentException](t, func() { union.FindAll("D") })
}

func TestStandardUnionFind_CheckEquivalent(t *testing.T) {
	union := NewStandardUnionFind[string]()
	union.Union("A", "B")
	union.Add("C")
	assertEqual(t, "A~B", union.AreEquivalent("A", "B"), true)
	assertEqual(t, "C~A", union.AreEquivalent("C", "A"), false)
	assertEqual(t, "C~B", union.AreEquivalent("C", "B"), false)
	assertThrows[jutil.IIllegalArgumentException](t, func() { union.AreEquivalent("A", "F") })
}
