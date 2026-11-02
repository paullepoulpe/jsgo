// Port of com/google/javascript/jscomp/disambiguate/ClusterPropagatorTest.java @ closure-compiler v20261028

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
	"iter"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
)

// ColorGraphNodePropsView adapts ColorGraphNode's flat property list to the
// map-like API expected by tests.
type ColorGraphNodePropsView struct {
	n *ColorGraphNode
}

func (v ColorGraphNodePropsView) Size() int { return len(v.n.props) }

func (v ColorGraphNodePropsView) IsEmpty() bool { return len(v.n.props) == 0 }

func (v ColorGraphNodePropsView) ContainsKey(prop *PropertyClustering) bool {
	return v.n.HasProp(prop)
}

func (v ColorGraphNodePropsView) Put(prop *PropertyClustering, _ ColorGraphNode_PropAssociation) {
	v.n.AddProp(prop)
}

func (v ColorGraphNodePropsView) PutIfAbsent(prop *PropertyClustering, _ ColorGraphNode_PropAssociation) {
	v.n.AddProp(prop)
}

func (v ColorGraphNodePropsView) Keys() iter.Seq[*PropertyClustering] {
	return func(yield func(*PropertyClustering) bool) {
		for _, p := range v.n.props {
			if !yield(p) {
				return
			}
		}
	}
}

// GetAssociatedProps is the set of properties that might be accessed
// from this type.
//
// Java: ColorGraphNode.getAssociatedProps()
func (n *ColorGraphNode) GetAssociatedProps() ColorGraphNodePropsView {
	return ColorGraphNodePropsView{n: n}
}

// ColorNodeUFElementsView provides a read-only set view over the elements of a
// ColorNodeUnionFind.
type ColorNodeUFElementsView struct {
	u *ColorNodeUnionFind
}

func (v ColorNodeUFElementsView) Size() int {
	return len(v.u.nodes)
}

func (v ColorNodeUFElementsView) IsEmpty() bool {
	return len(v.u.nodes) == 0
}

func (v ColorNodeUFElementsView) Contains(e *ColorGraphNode) bool {
	if e == nil {
		return false
	}
	_, ok := v.u.lookupSlot(e)
	return ok
}

func (v ColorNodeUFElementsView) All() iter.Seq[*ColorGraphNode] {
	return func(yield func(*ColorGraphNode) bool) {
		for _, el := range v.u.nodes {
			if !yield(el) {
				return
			}
		}
	}
}

func (u *ColorNodeUnionFind) Elements() ColorNodeUFElementsView {
	return ColorNodeUFElementsView{u: u}
}

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
	c.result = c.propagator.TraverseEdge(c.src, 0, c.dest)
}
