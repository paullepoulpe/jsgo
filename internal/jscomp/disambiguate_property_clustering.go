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

type colorUFNodeMeta struct {
	parent uint32
	rank   uint8
}

// ColorNodeUnionFind is a compact, pointer-minimal union-find specialized for
// *ColorGraphNode with exact StandardUnionFind rank and representative
// tie-breaking semantics.
type ColorNodeUnionFind struct {
	elems        []*ColorGraphNode
	meta         []colorUFNodeMeta
	slotByIdx    map[uint32]uint32
	inlineElems  [4]*ColorGraphNode
	inlineParent [4]uint8
	inlineRank   [4]uint8
	inlineLen    uint8
	numSets      uint32
}

func (u *ColorNodeUnionFind) findInlineRoot(slot uint8) uint8 {
	p := u.inlineParent[slot]
	if p == slot {
		return slot
	}
	root := p
	for u.inlineParent[root] != root {
		root = u.inlineParent[root]
	}
	for slot != root {
		next := u.inlineParent[slot]
		u.inlineParent[slot] = root
		slot = next
	}
	return root
}

func (u *ColorNodeUnionFind) findSpilledRoot(slot uint32) uint32 {
	p := u.meta[slot].parent
	if p == slot {
		return slot
	}
	root := p
	for u.meta[root].parent != root {
		root = u.meta[root].parent
	}
	for slot != root {
		next := u.meta[slot].parent
		u.meta[slot].parent = root
		slot = next
	}
	return root
}

func (u *ColorNodeUnionFind) spillToHeap(fifth *ColorGraphNode) uint32 {
	elems := make([]*ColorGraphNode, 5, 8)
	meta := make([]colorUFNodeMeta, 5, 8)
	slotByIdx := make(map[uint32]uint32, 8)
	for i := uint8(0); i < 4; i++ {
		el := u.inlineElems[i]
		elems[i] = el
		meta[i] = colorUFNodeMeta{
			parent: uint32(u.inlineParent[i]),
			rank:   u.inlineRank[i],
		}
		slotByIdx[uint32(el.GetIndex())] = uint32(i)
	}
	elems[4] = fifth
	meta[4] = colorUFNodeMeta{parent: 4, rank: 0}
	slotByIdx[uint32(fifth.GetIndex())] = 4

	u.inlineElems = [4]*ColorGraphNode{}
	u.elems = elems
	u.meta = meta
	u.slotByIdx = slotByIdx
	u.numSets++
	return 4
}

func (u *ColorNodeUnionFind) findRootOrCreate(e *ColorGraphNode) uint32 {
	if u.slotByIdx != nil {
		key := uint32(e.GetIndex())
		if slot, ok := u.slotByIdx[key]; ok {
			return u.findSpilledRoot(slot)
		}
		slot := uint32(len(u.elems))
		u.elems = append(u.elems, e)
		u.meta = append(u.meta, colorUFNodeMeta{parent: slot, rank: 0})
		u.slotByIdx[key] = slot
		u.numSets++
		return slot
	}
	for i := uint8(0); i < u.inlineLen; i++ {
		if u.inlineElems[i] == e {
			return uint32(u.findInlineRoot(i))
		}
	}
	if u.inlineLen < 4 {
		slot := u.inlineLen
		u.inlineElems[slot] = e
		u.inlineParent[slot] = slot
		u.inlineRank[slot] = 0
		u.inlineLen++
		u.numSets++
		return uint32(slot)
	}
	return u.spillToHeap(e)
}

func (u *ColorNodeUnionFind) Add(e *ColorGraphNode) {
	u.findRootOrCreate(e)
}

func (u *ColorNodeUnionFind) Union(a, b *ColorGraphNode) *ColorGraphNode {
	if a == b {
		root := u.findRootOrCreate(a)
		if u.slotByIdx == nil {
			return u.inlineElems[root]
		}
		return u.elems[root]
	}
	rootA := u.findRootOrCreate(a)
	rootB := u.findRootOrCreate(b)
	if u.slotByIdx == nil {
		ra := uint8(rootA)
		rb := uint8(rootB)
		if ra == rb {
			return u.inlineElems[ra]
		}
		u.numSets--
		if u.inlineRank[ra] >= u.inlineRank[rb] {
			u.inlineParent[rb] = ra
			if u.inlineRank[ra] == u.inlineRank[rb] {
				u.inlineRank[ra]++
			}
			return u.inlineElems[ra]
		}
		u.inlineParent[ra] = rb
		u.inlineElems[rb], u.inlineElems[ra] = u.inlineElems[ra], u.inlineElems[rb]
		return u.inlineElems[rb]
	}
	if rootA == rootB {
		return u.elems[rootA]
	}
	u.numSets--
	rankA := u.meta[rootA].rank
	rankB := u.meta[rootB].rank
	if rankA >= rankB {
		u.meta[rootB].parent = rootA
		if rankA == rankB {
			u.meta[rootA].rank = rankA + 1
		}
		return u.elems[rootA]
	}
	u.meta[rootA].parent = rootB
	elemA := u.elems[rootA]
	elemB := u.elems[rootB]
	u.elems[rootB] = elemA
	u.elems[rootA] = elemB
	u.slotByIdx[uint32(elemA.GetIndex())] = rootB
	u.slotByIdx[uint32(elemB.GetIndex())] = rootA
	return elemA
}

func (u *ColorNodeUnionFind) Find(e *ColorGraphNode) *ColorGraphNode {
	if u.slotByIdx == nil {
		for i := uint8(0); i < u.inlineLen; i++ {
			if u.inlineElems[i] == e {
				return u.inlineElems[u.findInlineRoot(i)]
			}
		}
		jutil.CheckArgumentMsg1(false, "Element does not exist: %s", e)
		return nil
	}
	slot, ok := u.slotByIdx[uint32(e.GetIndex())]
	if !ok || u.elems[slot] != e {
		jutil.CheckArgumentMsg1(false, "Element does not exist: %s", e)
		return nil
	}
	return u.elems[u.findSpilledRoot(slot)]
}

func (u *ColorNodeUnionFind) AreEquivalent(a, b *ColorGraphNode) bool {
	return u.Find(a) == u.Find(b)
}

func (u *ColorNodeUnionFind) NumSets() int {
	return int(u.numSets)
}

func (u *ColorNodeUnionFind) AppendRepresentatives(dst []*ColorGraphNode) []*ColorGraphNode {
	if u.slotByIdx == nil {
		for i := uint8(0); i < u.inlineLen; i++ {
			if u.inlineParent[i] == i {
				dst = append(dst, u.inlineElems[i])
			}
		}
		return dst
	}
	for i := range u.meta {
		if u.meta[i].parent == uint32(i) {
			dst = append(dst, u.elems[i])
		}
	}
	return dst
}

func (u *ColorNodeUnionFind) AllRepresentatives() *jutil.LinkedHashSet[*ColorGraphNode] {
	result := jutil.NewLinkedHashSet[*ColorGraphNode]()
	if u.slotByIdx == nil {
		for i := uint8(0); i < u.inlineLen; i++ {
			if u.inlineParent[i] == i {
				result.Add(u.inlineElems[i])
			}
		}
		return result
	}
	for i := range u.meta {
		if u.meta[i].parent == uint32(i) {
			result.Add(u.elems[i])
		}
	}
	return result
}

// ColorNodeUFElementsView provides a read-only set view over the elements of a
// ColorNodeUnionFind.
type ColorNodeUFElementsView struct {
	u *ColorNodeUnionFind
}

func (v ColorNodeUFElementsView) Size() int {
	if v.u.slotByIdx == nil {
		return int(v.u.inlineLen)
	}
	return len(v.u.elems)
}

func (v ColorNodeUFElementsView) IsEmpty() bool {
	return v.Size() == 0
}

func (v ColorNodeUFElementsView) Contains(e *ColorGraphNode) bool {
	if e == nil {
		return false
	}
	if v.u.slotByIdx == nil {
		for i := uint8(0); i < v.u.inlineLen; i++ {
			if v.u.inlineElems[i] == e {
				return true
			}
		}
		return false
	}
	slot, ok := v.u.slotByIdx[uint32(e.GetIndex())]
	return ok && v.u.elems[slot] == e
}

func (v ColorNodeUFElementsView) All() iter.Seq[*ColorGraphNode] {
	return func(yield func(*ColorGraphNode) bool) {
		if v.u.slotByIdx == nil {
			for i := uint8(0); i < v.u.inlineLen; i++ {
				if !yield(v.u.inlineElems[i]) {
					return
				}
			}
			return
		}
		for _, el := range v.u.elems {
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
	name string

	useSites PropertyUseSites

	clusters ColorNodeUnionFind

	originalNameClusterRep *ColorGraphNode

	lastInvalidation Invalidation
}

// Java: PropertyClustering(String)
func NewPropertyClustering(name string) *PropertyClustering {
	// checkNotNull(name): a Go string is never null.
	return &PropertyClustering{
		name: name,
	}
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
	return jutil.MoreObjects_ToStringHelper("PropertyClustering").Add("name", p.name).String()
}
