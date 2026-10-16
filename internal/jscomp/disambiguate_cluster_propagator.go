// Port of com/google/javascript/jscomp/disambiguate/ClusterPropagator.java @ closure-compiler v20260415

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
	"github.com/paullepoulpe/jsgo/internal/colors"
)

// ClusterPropagator is a callback to propagate clusterings across a type
// graph.
type ClusterPropagator struct{}

// Java: ClusterPropagator()
func NewClusterPropagator() *ClusterPropagator {
	return &ClusterPropagator{}
}

// Java: ClusterPropagator.traverseEdge(ColorGraphNode, Object, ColorGraphNode)
func (p *ClusterPropagator) TraverseEdge(src *ColorGraphNode, unused any, dest *ColorGraphNode) bool {
	if src.GetColor() == colors.StandardColors_TOP_OBJECT {
		// Performance optimization. The for loop below already skips
		// invalidated properties. Colors that are inherently invalidating
		// will have all their properties invalidated already.
		// For example - in a sample large project, TOP_OBJECT had ~4.7k
		// properties and it had ~700k direct subtypes. So the for loop below
		// had to iterate over 4.7k * 700k = ~3.3B properties, just for the
		// TOP_OBJECT node, that were all already invalidated.
		return false
	}

	startDestPropCount := dest.GetAssociatedProps().Size()

	for prop := range src.GetAssociatedProps().Keys() {
		if prop.IsInvalidated() {
			continue
		}

		dest.GetAssociatedProps().PutIfAbsent(prop, ColorGraphNode_PropAssociation_SUPERTYPE)
		prop.GetClusters().Union(src, dest)
	}

	// Were any properties added to dest?
	return startDestPropCount < dest.GetAssociatedProps().Size()
}

var _ FixedPointGraphTraversal_EdgeCallback[*ColorGraphNode, any] = (*ClusterPropagator)(nil)
