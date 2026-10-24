// Port of com/google/javascript/jscomp/disambiguate/PropertyClustering.java @ closure-compiler v20261021

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
	"iter"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type PropertyUseSite struct {
	site  *rhino.Node
	owner *ColorGraphNode
}

type propertyUseSiteChunk struct {
	entries []PropertyUseSite
	next    *propertyUseSiteChunk
}

// PropertyUseSites stores the (AST node, receiver ColorGraphNode) pairs for a
// PropertyClustering in insertion order.
type PropertyUseSites struct {
	inline   [2]PropertyUseSite
	head     *propertyUseSiteChunk
	tail     *propertyUseSiteChunk
	count    int32
	tailUsed int32
}

func (u *PropertyUseSites) Size() int {
	return int(u.count)
}

func (u *PropertyUseSites) IsEmpty() bool {
	return u.count == 0
}

func (u *PropertyUseSites) Clear() {
	u.inline = [2]PropertyUseSite{}
	u.head = nil
	u.tail = nil
	u.count = 0
	u.tailUsed = 0
}

func (u *PropertyUseSites) PutNew(site *rhino.Node, owner *ColorGraphNode) {
	c := u.count
	u.count = c + 1
	if c < int32(len(u.inline)) {
		u.inline[c] = PropertyUseSite{site: site, owner: owner}
		return
	}
	tail := u.tail
	used := u.tailUsed
	if tail != nil && int(used) < len(tail.entries) {
		tail.entries[used] = PropertyUseSite{site: site, owner: owner}
		u.tailUsed = used + 1
		return
	}
	chunkSize := 4
	if c >= 6 {
		chunkSize = 16
	}
	ch := &propertyUseSiteChunk{
		entries: make([]PropertyUseSite, chunkSize),
	}
	ch.entries[0] = PropertyUseSite{site: site, owner: owner}
	if u.tail == nil {
		u.head = ch
	} else {
		u.tail.next = ch
	}
	u.tail = ch
	u.tailUsed = 1
}

func (u *PropertyUseSites) Put(site *rhino.Node, owner *ColorGraphNode) *ColorGraphNode {
	nInline := int(u.count)
	if nInline > len(u.inline) {
		nInline = len(u.inline)
	}
	for i := 0; i < nInline; i++ {
		if u.inline[i].site == site {
			old := u.inline[i].owner
			u.inline[i].owner = owner
			return old
		}
	}
	for ch := u.head; ch != nil; ch = ch.next {
		entries := ch.entries
		if ch == u.tail {
			entries = entries[:u.tailUsed]
		}
		for i := range entries {
			if entries[i].site == site {
				old := entries[i].owner
				entries[i].owner = owner
				return old
			}
		}
	}
	u.PutNew(site, owner)
	return nil
}

func (u *PropertyUseSites) Get(site *rhino.Node) *ColorGraphNode {
	nInline := int(u.count)
	if nInline > len(u.inline) {
		nInline = len(u.inline)
	}
	for i := 0; i < nInline; i++ {
		if u.inline[i].site == site {
			return u.inline[i].owner
		}
	}
	for ch := u.head; ch != nil; ch = ch.next {
		entries := ch.entries
		if ch == u.tail {
			entries = entries[:u.tailUsed]
		}
		for i := range entries {
			if entries[i].site == site {
				return entries[i].owner
			}
		}
	}
	return nil
}

func (u *PropertyUseSites) All() iter.Seq2[*rhino.Node, *ColorGraphNode] {
	return func(yield func(*rhino.Node, *ColorGraphNode) bool) {
		nInline := int(u.count)
		if nInline > len(u.inline) {
			nInline = len(u.inline)
		}
		for i := 0; i < nInline; i++ {
			e := u.inline[i]
			if !yield(e.site, e.owner) {
				return
			}
		}
		for ch := u.head; ch != nil; ch = ch.next {
			entries := ch.entries
			if ch == u.tail {
				entries = entries[:u.tailUsed]
			}
			for i := range entries {
				e := entries[i]
				if !yield(e.site, e.owner) {
					return
				}
			}
		}
	}
}

func (u *PropertyUseSites) Keys() iter.Seq[*rhino.Node] {
	return func(yield func(*rhino.Node) bool) {
		for k := range u.All() {
			if !yield(k) {
				return
			}
		}
	}
}

func (u *PropertyUseSites) KeySlice() []*rhino.Node {
	if u.count == 0 {
		return nil
	}
	out := make([]*rhino.Node, 0, u.count)
	for k := range u.All() {
		out = append(out, k)
	}
	return out
}

// PropertyClustering is the disambiguation clusters for a given property
// name.
//
// This is a struct used to aggregate information to be processed by other
// classes in this package. It is intentionally mutable and doesn't attempt
// to enforce invariants in the contents of its datastructures. Instances
// trust that sibling classes make mutations correctly.
type PropertyClustering struct {
	name string

	useSites PropertyUseSites

	clusters    StandardUnionFind[*ColorGraphNode]
	clustersMap jutil.LinkedHashMap[*ColorGraphNode, *standardUnionFind_Node[*ColorGraphNode]]
	clustersBuf [4]jutil.LhmSlot[*ColorGraphNode, *standardUnionFind_Node[*ColorGraphNode]]

	originalNameClusterRep *ColorGraphNode
	lastOwner              *ColorGraphNode
	lastOwner2             *ColorGraphNode

	lastInvalidation Invalidation
}

func (p *PropertyClustering) init(name string) {
	p.name = name
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
func (p *PropertyClustering) GetUseSites() *PropertyUseSites {
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
