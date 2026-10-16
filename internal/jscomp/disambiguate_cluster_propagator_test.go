// Port of com/google/javascript/jscomp/disambiguate/ClusterPropagatorTest.java @ closure-compiler v20260415

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

	"github.com/paullepoulpe/jsgo/internal/colors"
)

var clusterPropagatorTest_TEST_COLOR = colors.StandardColors_NUMBER

// clusterPropagatorTest holds the per-test fields of ClusterPropagatorTest.
type clusterPropagatorTest struct {
	t          *testing.T
	propagator *ClusterPropagator
	src        *ColorGraphNode
	dest       *ColorGraphNode
	prop       *PropertyClustering
	result     bool
}

// runClusterPropagatorTest runs body on fresh fields, then the @After
// checks.
func runClusterPropagatorTest(t *testing.T, body func(c *clusterPropagatorTest)) {
	c := &clusterPropagatorTest{
		t:          t,
		propagator: NewClusterPropagator(),
		src:        ColorGraphNode_CreateForTestingWithColor(clusterPropagatorTest_TEST_COLOR, -1),
		dest:       ColorGraphNode_CreateForTestingWithColor(clusterPropagatorTest_TEST_COLOR, -2),
		prop:       NewPropertyClustering("prop"),
	}
	body(c)
	c.verifyPropertyFlow()
	c.verifyClusterMerging()
}

func TestClusterPropagator_Propagation_mergesProperties_fromSrcToDest(t *testing.T) {
	runClusterPropagatorTest(t, func(c *clusterPropagatorTest) {
		// Given
		clusterPropagatorTest_associate(c.prop, c.src)

		// When
		c.propagateFromSrcToDest()

		// Then
		assertTrue(t, "result", c.result)
	})
}

func TestClusterPropagator_Propagation_mergesProperties_fromSrc_thatExistOnDest(t *testing.T) {
	runClusterPropagatorTest(t, func(c *clusterPropagatorTest) {
		// Given
		clusterPropagatorTest_associate(c.prop, c.src)
		clusterPropagatorTest_associate(c.prop, c.dest)

		// When
		c.propagateFromSrcToDest()

		// Then
		assertFalse(t, "result", c.result)
	})
}

func TestClusterPropagator_Propagation_doesNotMergeProperties_fromDestToSrc(t *testing.T) {
	runClusterPropagatorTest(t, func(c *clusterPropagatorTest) {
		// Given
		clusterPropagatorTest_associate(c.prop, c.dest)

		// When
		c.propagateFromSrcToDest()

		// Then
		assertFalse(t, "result", c.result)
		assertTrue(t, "src props empty", c.src.GetAssociatedProps().IsEmpty())
	})
}

func TestClusterPropagator_Propagation_doesNotMergeProperties_thatAreInvalidated(t *testing.T) {
	runClusterPropagatorTest(t, func(c *clusterPropagatorTest) {
		// Given
		clusterPropagatorTest_associate(c.prop, c.src)
		c.prop.Invalidate(Invalidation_WellKnownProperty())

		// When
		c.propagateFromSrcToDest()

		// Then
		assertFalse(t, "result", c.result)
		assertTrue(t, "dest props empty", c.dest.GetAssociatedProps().IsEmpty())
	})
}

// Java: @After verifyPropertyFlow()
func (c *clusterPropagatorTest) verifyPropertyFlow() {
	c.t.Helper()
	for p := range c.src.GetAssociatedProps().Keys() {
		if !p.IsInvalidated() && !c.dest.GetAssociatedProps().ContainsKey(p) {
			c.t.Errorf("dest is missing valid src prop %v", p)
		}
	}
}

// Java: @After verifyClusterMerging()
func (c *clusterPropagatorTest) verifyClusterMerging() {
	c.t.Helper()
	if c.prop.IsInvalidated() {
		return
	}

	if clusterPropagatorTest_areAssociated(c.t, c.prop, c.src) && clusterPropagatorTest_areAssociated(c.t, c.prop, c.dest) {
		assertTrue(c.t, "areEquivalent(src, dest)", c.prop.GetClusters().AreEquivalent(c.src, c.dest))
	}
}

func clusterPropagatorTest_areAssociated(t *testing.T, prop *PropertyClustering, flat *ColorGraphNode) bool {
	t.Helper()
	if prop.GetClusters().Elements().Contains(flat) {
		assertTrue(t, "associated props contain prop", flat.GetAssociatedProps().ContainsKey(prop))
		return true
	}

	assertFalse(t, "associated props contain prop", flat.GetAssociatedProps().ContainsKey(prop))
	return false
}

func clusterPropagatorTest_associate(prop *PropertyClustering, node *ColorGraphNode) {
	node.GetAssociatedProps().Put(prop, ColorGraphNode_PropAssociation_AST)
	prop.GetClusters().Add(node)
}

func (c *clusterPropagatorTest) propagateFromSrcToDest() {
	c.result = c.propagator.TraverseEdge(c.src, nil, c.dest)
}
