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
	site  rhino.Node
	owner *ColorGraphNode
}

// PropertyUseSites stores the (AST node, receiver ColorGraphNode) pairs for a
// PropertyClustering in insertion order.
type PropertyUseSites struct {
	sites  []PropertyUseSite
	inline [2]PropertyUseSite
}

func (u *PropertyUseSites) Size() int {
	return len(u.sites)
}

func (u *PropertyUseSites) IsEmpty() bool {
	return len(u.sites) == 0
}

func (u *PropertyUseSites) Clear() {
	clear(u.sites)
	u.sites = u.inline[:0]
}

func (u *PropertyUseSites) PutNew(site rhino.Node, owner *ColorGraphNode) {
	if u.sites == nil {
		u.sites = u.inline[:0]
	}
	u.sites = append(u.sites, PropertyUseSite{site: site, owner: owner})
}

func (u *PropertyUseSites) Put(site rhino.Node, owner *ColorGraphNode) *ColorGraphNode {
	for i := range u.sites {
		if u.sites[i].site == site {
			old := u.sites[i].owner
			u.sites[i].owner = owner
			return old
		}
	}
	u.PutNew(site, owner)
	return nil
}

func (u *PropertyUseSites) Get(site rhino.Node) *ColorGraphNode {
	for i := range u.sites {
		if u.sites[i].site == site {
			return u.sites[i].owner
		}
	}
	return nil
}

func (u *PropertyUseSites) All() iter.Seq2[rhino.Node, *ColorGraphNode] {
	return func(yield func(rhino.Node, *ColorGraphNode) bool) {
		for _, e := range u.sites {
			if !yield(e.site, e.owner) {
				return
			}
		}
	}
}

func (u *PropertyUseSites) Keys() iter.Seq[rhino.Node] {
	return func(yield func(rhino.Node) bool) {
		for _, e := range u.sites {
			if !yield(e.site) {
				return
			}
		}
	}
}

func (u *PropertyUseSites) KeySlice() []rhino.Node {
	if len(u.sites) == 0 {
		return nil
	}
	out := make([]rhino.Node, len(u.sites))
	for i, e := range u.sites {
		out[i] = e.site
	}
	return out
}

// ColorNodeUnionFind is a compact union-find specialized for *ColorGraphNode
// with exact StandardUnionFind rank and representative tie-breaking semantics.
type ColorNodeUnionFind struct {
	nodes     []*ColorGraphNode
	parent    []int32
	rank      []uint8
	slotByIdx map[uint32]int32
	nodesBuf  [4]*ColorGraphNode
	parentBuf [4]int32
	rankBuf   [4]uint8
	numSets   uint32
}

func (u *ColorNodeUnionFind) findRoot(slot int32) int32 {
	root := u.parent[slot]
	if root == slot {
		return slot
	}
	for u.parent[root] != root {
		root = u.parent[root]
	}
	for slot != root {
		next := u.parent[slot]
		u.parent[slot] = root
		slot = next
	}
	return root
}

func (u *ColorNodeUnionFind) lookupSlot(e *ColorGraphNode) (int32, bool) {
	if u.slotByIdx != nil {
		slot, ok := u.slotByIdx[uint32(e.GetIndex())]
		return slot, ok && u.nodes[slot] == e
	}
	for i, n := range u.nodes {
		if n == e {
			return int32(i), true
		}
	}
	return -1, false
}

func (u *ColorNodeUnionFind) findRootOrCreate(e *ColorGraphNode) int32 {
	if slot, ok := u.lookupSlot(e); ok {
		return u.findRoot(slot)
	}
	if u.nodes == nil {
		u.nodes = u.nodesBuf[:0]
		u.parent = u.parentBuf[:0]
		u.rank = u.rankBuf[:0]
	} else if u.slotByIdx == nil && len(u.nodes) == 8 {
		u.slotByIdx = make(map[uint32]int32, 16)
		for i, n := range u.nodes {
			u.slotByIdx[uint32(n.GetIndex())] = int32(i)
		}
	}
	slot := int32(len(u.nodes))
	u.nodes = append(u.nodes, e)
	u.parent = append(u.parent, slot)
	u.rank = append(u.rank, 0)
	if u.slotByIdx != nil {
		u.slotByIdx[uint32(e.GetIndex())] = slot
	}
	u.numSets++
	return slot
}

func (u *ColorNodeUnionFind) Add(e *ColorGraphNode) {
	u.findRootOrCreate(e)
}

func (u *ColorNodeUnionFind) Union(a, b *ColorGraphNode) *ColorGraphNode {
	if a == b {
		return u.nodes[u.findRootOrCreate(a)]
	}
	rootA := u.findRootOrCreate(a)
	rootB := u.findRootOrCreate(b)
	if rootA == rootB {
		return u.nodes[rootA]
	}
	u.numSets--
	rankA, rankB := u.rank[rootA], u.rank[rootB]
	if rankA >= rankB {
		u.parent[rootB] = rootA
		if rankA == rankB {
			u.rank[rootA] = rankA + 1
		}
		return u.nodes[rootA]
	}
	u.parent[rootA] = rootB
	elemA, elemB := u.nodes[rootA], u.nodes[rootB]
	u.nodes[rootB], u.nodes[rootA] = elemA, elemB
	if u.slotByIdx != nil {
		u.slotByIdx[uint32(elemA.GetIndex())] = rootB
		u.slotByIdx[uint32(elemB.GetIndex())] = rootA
	}
	return elemA
}

func (u *ColorNodeUnionFind) Find(e *ColorGraphNode) *ColorGraphNode {
	slot, ok := u.lookupSlot(e)
	if !ok {
		jutil.CheckArgumentMsg1(false, "Element does not exist: %s", e)
		return nil
	}
	return u.nodes[u.findRoot(slot)]
}

func (u *ColorNodeUnionFind) AreEquivalent(a, b *ColorGraphNode) bool {
	return u.Find(a) == u.Find(b)
}

func (u *ColorNodeUnionFind) NumSets() int {
	return int(u.numSets)
}

func (u *ColorNodeUnionFind) AppendRepresentatives(dst []*ColorGraphNode) []*ColorGraphNode {
	for i, p := range u.parent {
		if p == int32(i) {
			dst = append(dst, u.nodes[i])
		}
	}
	return dst
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

// PropertyClustering is the disambiguation clusters for a given property
// name.
//
// This is a struct used to aggregate information to be processed by other
// classes in this package. It is intentionally mutable and doesn't attempt
// to enforce invariants in the contents of its datastructures. Instances
// trust that sibling classes make mutations correctly.
type PropertyClustering struct {
	name   string
	nameID uint32

	useSites PropertyUseSites

	clusters ColorNodeUnionFind

	originalNameClusterRep *ColorGraphNode

	lastInvalidation Invalidation
}

// Java: PropertyClustering(String)
func NewPropertyClustering(name string) *PropertyClustering {
	// checkNotNull(name): a Go string is never null.
	return &PropertyClustering{
		name:   name,
		nameID: rhino.RhinoStringPool_InternID(name),
	}
}

// Java: PropertyClustering.getName()
func (p *PropertyClustering) GetName() string {
	if p.name == "" && p.nameID != 0 {
		p.name = rhino.RhinoStringPool_StringFromID(p.nameID)
	}
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
func (p *PropertyClustering) GetClusters() *ColorNodeUnionFind {
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
	return p.lastInvalidation.IsValid()
}

// Java: PropertyClustering.invalidate(Invalidation)
func (p *PropertyClustering) Invalidate(invalidation Invalidation) {
	jutil.CheckState(invalidation.IsValid())
	p.lastInvalidation = invalidation
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
	return jutil.MoreObjects_ToStringHelper("PropertyClustering").Add("name", p.GetName()).String()
}

type propertyClusteringIndex struct {
	slab  jutil.Slab[PropertyClustering]
	order []*PropertyClustering
	byID  map[uint32]uint32
}

func newPropertyClusteringIndex(capHint int) *propertyClusteringIndex {
	return &propertyClusteringIndex{
		order: make([]*PropertyClustering, 0, capHint),
		byID:  make(map[uint32]uint32, capHint),
	}
}

func (idx *propertyClusteringIndex) getOrCreateByID(id uint32) *PropertyClustering {
	if pos, ok := idx.byID[id]; ok {
		return idx.order[pos]
	}
	p := idx.slab.New()
	p.nameID = id
	pos := uint32(len(idx.order))
	idx.order = append(idx.order, p)
	idx.byID[id] = pos
	return p
}

func (idx *propertyClusteringIndex) getOrCreateByName(name string) *PropertyClustering {
	id := rhino.RhinoStringPool_InternID(name)
	if pos, ok := idx.byID[id]; ok {
		return idx.order[pos]
	}
	p := idx.slab.New()
	p.name = name
	p.nameID = id
	pos := uint32(len(idx.order))
	idx.order = append(idx.order, p)
	idx.byID[id] = pos
	return p
}

func (idx *propertyClusteringIndex) getByName(name string) *PropertyClustering {
	id := rhino.RhinoStringPool_LookupID(name)
	if id == 0 && name != "" {
		return nil
	}
	if pos, ok := idx.byID[id]; ok {
		return idx.order[pos]
	}
	return nil
}

func (idx *propertyClusteringIndex) values() []*PropertyClustering {
	return idx.order
}
