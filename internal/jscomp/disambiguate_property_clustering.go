// Port of com/google/javascript/jscomp/disambiguate/PropertyClustering.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PropertyClustering is the disambiguation clusters for a given property
// name.
//
// This is a struct used to aggregate information to be processed by other
// classes in this package. It is intentionally mutable and doesn't attempt
// to enforce invariants in the contents of its datastructures. Instances
// trust that sibling classes make mutations correctly.
type PropertyClustering struct {
	name string

	useSites    jutil.LinkedHashMap[*rhino.Node, *ColorGraphNode]
	useSitesBuf [8]jutil.LhmSlot[*rhino.Node, *ColorGraphNode]

	clusters    StandardUnionFind[*ColorGraphNode]
	clustersMap jutil.LinkedHashMap[*ColorGraphNode, *standardUnionFind_Node[*ColorGraphNode]]
	clustersBuf [4]jutil.LhmSlot[*ColorGraphNode, *standardUnionFind_Node[*ColorGraphNode]]

	originalNameClusterRep *ColorGraphNode
	lastOwner              *ColorGraphNode

	lastInvalidation Invalidation
}

func (p *PropertyClustering) init(name string) {
	p.name = name
	p.useSites.InitWithBuffer(p.useSitesBuf[:])
	p.clustersMap.InitWithBuffer(p.clustersBuf[:])
	p.clusters.elmap = &p.clustersMap
}

// Java: PropertyClustering(String)
func NewPropertyClustering(name string) *PropertyClustering {
	// checkNotNull(name): a Go string is never null.
	p := &PropertyClustering{}
	p.init(name)
	return p
}

// Java: PropertyClustering.getName()
func (p *PropertyClustering) GetName() string {
	return p.name
}

// GetUseSites returns the locations properties with this name were used in
// the program, mapping to their receiver type.
//
// This index allows property references to be efficiently renamed once all
// clusters have been found. It prevents us from re-traversing the code.
//
// Java: PropertyClustering.getUseSites()
func (p *PropertyClustering) GetUseSites() *jutil.LinkedHashMap[*rhino.Node, *ColorGraphNode] {
	return &p.useSites
}

// Java: PropertyClustering.getClusters()
func (p *PropertyClustering) GetClusters() *StandardUnionFind[*ColorGraphNode] {
	return &p.clusters
}

// GetOriginalNameClusterRep is the current representative of the cluster
// of types whose properties must keep their original name.
//
// The following refers to all computations with respect to a single
// property name. Since extern and enum properties cannot be renamed, all
// other types in a cluster with an extern type or enum cannot rename their
// property either. In theory, there could be many such clusters containing
// an extern or enum; however, in practice we conflate them into one. This
// is equivalent because all of those clusters would end up using the same,
// unchanged, property name. This representation also simplifies tracking
// of such clusters.
//
// In practice the types in this cluster include
//   - externs, whose properties cannot be renamed without breaking
//     references to external code
//   - enums, e.g. for const Actions = {STOP: 0, GO: 1}; the disambiguator
//     will not rename STOP or GO
//   - boxable scalars, like string and number
//
// Note that enum properties could probably be safely renamed, but this
// would require cleaning up code depending on the legacy no-renaming
// behavior.
//
// Java: PropertyClustering.getOriginalNameClusterRep()
func (p *PropertyClustering) GetOriginalNameClusterRep() *ColorGraphNode {
	jutil.CheckState(!p.IsInvalidated())
	if p.originalNameClusterRep == nil {
		return nil
	}
	return p.clusters.Find(p.originalNameClusterRep)
}

// Java: PropertyClustering.isInvalidated()
func (p *PropertyClustering) IsInvalidated() bool {
	return p.lastInvalidation != nil
}

// Java: PropertyClustering.invalidate(Invalidation)
func (p *PropertyClustering) Invalidate(invalidation Invalidation) {
	p.lastInvalidation = jutil.CheckNotNull(invalidation)
}

// Java: PropertyClustering.getLastInvalidation()
func (p *PropertyClustering) GetLastInvalidation() Invalidation {
	jutil.CheckState(p.IsInvalidated())
	return p.lastInvalidation
}

// RegisterOriginalNameType indicates that all property references off this
// ColorGraphNode must keep their original name.
//
// See GetOriginalNameClusterRep for more details.
//
// Java: PropertyClustering.registerOriginalNameType(ColorGraphNode)
func (p *PropertyClustering) RegisterOriginalNameType(type_ *ColorGraphNode) {
	jutil.CheckState(!p.IsInvalidated())
	if p.originalNameClusterRep == nil {
		p.originalNameClusterRep = type_
		p.clusters.Add(type_)
	} else if p.originalNameClusterRep != type_ {
		p.clusters.Union(p.originalNameClusterRep, type_)
	}
}

// String is for debugging only (@DoNotCall).
//
// Java: PropertyClustering.toString()
func (p *PropertyClustering) String() string {
	return jutil.MoreObjects_ToStringHelper("PropertyClustering").Add("name", p.name).String()
}
