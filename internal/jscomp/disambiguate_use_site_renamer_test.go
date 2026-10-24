// Port of com/google/javascript/jscomp/disambiguate/UseSiteRenamerTest.java @ closure-compiler v20261015

/*
 * Copyright 2019 The Closure Compiler Authors.
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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const useSiteRenamerTest_PROP_NAME = "prop"

// useSiteRenamerTest holds the per-test state of UseSiteRenamerTest; its
// @After checks run in finish.
type useSiteRenamerTest struct {
	t                 *testing.T
	reportedMutations *jutil.LinkedHashSet[rhino.Node]
	prop              *PropertyClustering
	renamingIndex     *jutil.SetMultimap[string, string]
}

func newUseSiteRenamerTest(t *testing.T) *useSiteRenamerTest {
	return &useSiteRenamerTest{
		t:                 t,
		reportedMutations: jutil.NewLinkedHashSet[rhino.Node](),
		prop:              NewPropertyClustering(useSiteRenamerTest_PROP_NAME),
	}
}

func (s *useSiteRenamerTest) finish() {
	s.verifyRenamingIndexContainsExactlyNewNames()
	s.verifyNewNamesAreReportedAsMutations()
}

func (s *useSiteRenamerTest) verifyRenamingIndexContainsExactlyNewNames() {
	t := s.t
	t.Helper()
	if s.prop.IsInvalidated() {
		var entries []string
		for k, v := range s.renamingIndex.Entries() {
			entries = append(entries, k+"="+v)
		}
		assertSliceEqual(t, "renaming index", entries, []string{useSiteRenamerTest_PROP_NAME + "=<INVALIDATED>"})
		return
	}

	var values []string
	for _, v := range s.renamingIndex.Entries() {
		values = append(values, v)
	}
	assertSameElements(t, "renaming index values", values, useSiteRenamerTest_namesOf(s.prop.GetUseSites().KeySlice()))
}

func (s *useSiteRenamerTest) verifyNewNamesAreReportedAsMutations() {
	t := s.t
	t.Helper()
	if s.prop.IsInvalidated() {
		return
	}

	mutatedNodes := jutil.NewLinkedHashSet[rhino.Node]()
	for n := range s.prop.GetUseSites().Keys() {
		if n.GetString() != useSiteRenamerTest_PROP_NAME {
			mutatedNodes.Add(n)
		}
	}

	assertSameElements(t, "mutations", s.reportedMutations.ToSlice(), mutatedNodes.ToSlice())
}

func (s *useSiteRenamerTest) runRename() {
	renamer := NewUseSiteRenamer(func(n rhino.Node) { s.reportedMutations.Add(n) })
	renamer.RenameUses(s.prop)

	s.renamingIndex = renamer.GetRenamingIndex()
}

func useSiteRenamerTest_namesOf(useSites []rhino.Node) []string {
	out := jutil.NewLinkedHashSet[string]()
	for _, n := range useSites {
		out.Add(n.GetString())
	}
	return out.ToSlice()
}

func TestUseSiteRenamer_RenameUsesRenamesConsistentlyWithinEachCluster(t *testing.T) {
	s := newUseSiteRenamerTest(t)
	defer s.finish()
	// Given
	name1a := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)
	name1b := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)
	name1c := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)
	name2 := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)

	type1a := ColorGraphNode_CreateForTesting(-1)
	type1b := ColorGraphNode_CreateForTesting(-2)
	type1c := ColorGraphNode_CreateForTesting(-3)
	type2 := ColorGraphNode_CreateForTesting(-4)

	s.prop.GetUseSites().Put(name1a, type1a)
	s.prop.GetUseSites().Put(name1b, type1b)
	s.prop.GetUseSites().Put(name1c, type1c)
	s.prop.GetUseSites().Put(name2, type2)

	s.prop.GetClusters().Union(type1a, type1b)
	s.prop.GetClusters().Union(type1a, type1c)
	s.prop.GetClusters().Add(type2)

	// When
	s.runRename()

	// Then
	assertEqual(t, "1b", name1b.GetString(), name1a.GetString())
	assertEqual(t, "1c", name1c.GetString(), name1a.GetString())
	if name1a.GetString() == useSiteRenamerTest_PROP_NAME {
		t.Error("1a not renamed")
	}
}

func TestUseSiteRenamer_RenameUsesRenamesDistinctlyBetweenEachCluster(t *testing.T) {
	s := newUseSiteRenamerTest(t)
	defer s.finish()
	// Given
	name1 := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)
	name2 := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)
	name3 := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)

	type1 := ColorGraphNode_CreateForTesting(-1)
	type2 := ColorGraphNode_CreateForTesting(-2)
	type3 := ColorGraphNode_CreateForTesting(-3)

	s.prop.GetUseSites().Put(name1, type1)
	s.prop.GetUseSites().Put(name2, type2)
	s.prop.GetUseSites().Put(name3, type3)

	s.prop.GetClusters().Add(type1)
	s.prop.GetClusters().Add(type2)
	s.prop.GetClusters().Add(type3)

	// When
	s.runRename()

	// Then
	if name1.GetString() == name2.GetString() || name1.GetString() == name3.GetString() ||
		name2.GetString() == name3.GetString() {
		t.Errorf("names not distinct: %s %s %s", name1.GetString(), name2.GetString(), name3.GetString())
	}
}

func TestUseSiteRenamer_RenameUsesDoesntRenameExternsCluster(t *testing.T) {
	s := newUseSiteRenamerTest(t)
	defer s.finish()
	// Given
	externName1 := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)
	externName2 := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)
	externName3 := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)
	srcName := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)

	externType1 := ColorGraphNode_CreateForTesting(-1)
	externType2 := ColorGraphNode_CreateForTesting(-2)
	externType3 := ColorGraphNode_CreateForTesting(-3)
	srcType := ColorGraphNode_CreateForTesting(-4)

	s.prop.GetUseSites().Put(externName1, externType1)
	s.prop.GetUseSites().Put(externName2, externType2)
	s.prop.GetUseSites().Put(externName3, externType3)
	s.prop.GetUseSites().Put(srcName, srcType)

	s.prop.RegisterOriginalNameType(externType1)
	s.prop.RegisterOriginalNameType(externType2)
	s.prop.RegisterOriginalNameType(externType3)

	s.prop.GetClusters().Add(srcType)

	// When
	s.runRename()

	// Then
	assertEqual(t, "extern 1", externName1.GetString(), useSiteRenamerTest_PROP_NAME)
	assertEqual(t, "extern 2", externName2.GetString(), useSiteRenamerTest_PROP_NAME)
	assertEqual(t, "extern 3", externName3.GetString(), useSiteRenamerTest_PROP_NAME)
	if srcName.GetString() == useSiteRenamerTest_PROP_NAME {
		t.Error("src not renamed")
	}
}

func TestUseSiteRenamer_RenameUsesDoesntRenameInvalidatedProps(t *testing.T) {
	s := newUseSiteRenamerTest(t)
	defer s.finish()
	// Given
	name1 := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)
	name2 := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)
	name3 := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)

	type1 := ColorGraphNode_CreateForTesting(-1)
	type2 := ColorGraphNode_CreateForTesting(-2)
	type3 := ColorGraphNode_CreateForTesting(-3)

	s.prop.GetUseSites().Put(name1, type1)
	s.prop.GetUseSites().Put(name2, type2)
	s.prop.GetUseSites().Put(name3, type3)

	s.prop.Invalidate(Invalidation_WellKnownProperty())

	// When
	s.runRename()

	// Then
	assertEqual(t, "1", name1.GetString(), useSiteRenamerTest_PROP_NAME)
	assertEqual(t, "1", name1.GetString(), useSiteRenamerTest_PROP_NAME)
	assertEqual(t, "2", name2.GetString(), useSiteRenamerTest_PROP_NAME)
}

func TestUseSiteRenamer_RenameUsesDoesntRenamePropWhenThereIsOnlyOneCluster(t *testing.T) {
	s := newUseSiteRenamerTest(t)
	defer s.finish()
	// Given
	name1a := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)
	name1b := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)
	name1c := rhino.IR_Name(useSiteRenamerTest_PROP_NAME)

	type1a := ColorGraphNode_CreateForTesting(-1)
	type1b := ColorGraphNode_CreateForTesting(-2)
	type1c := ColorGraphNode_CreateForTesting(-3)

	s.prop.GetUseSites().Put(name1a, type1a)
	s.prop.GetUseSites().Put(name1b, type1b)
	s.prop.GetUseSites().Put(name1c, type1c)

	s.prop.GetClusters().Union(type1a, type1b)
	s.prop.GetClusters().Union(type1a, type1c)

	// When
	s.runRename()

	// Then
	assertEqual(t, "1b", name1b.GetString(), useSiteRenamerTest_PROP_NAME)
	assertEqual(t, "1c", name1c.GetString(), useSiteRenamerTest_PROP_NAME)
	assertEqual(t, "1a", name1a.GetString(), useSiteRenamerTest_PROP_NAME)
}
