// Port of com/google/javascript/jscomp/disambiguate/UseSiteRenamer.java @ closure-compiler v20261021

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
	"maps"
	"slices"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// UseSiteRenamer applies renaming to property use sites following cluster
// computation.
type UseSiteRenamer struct {
	mutationCb func(rhino.Node)

	// Keys and values in insertion order, populated lazily when trackRenamingIndex is true.
	renamingEntries     map[string][]string
	renamingOrder       []string
	trackRenamingIndex  bool
	scratchReps         []*ColorGraphNode
	scratchClusterNames map[*ColorGraphNode]string
}

const useSiteRenamer_INVALIDATED_NAME_VALUE = "<INVALIDATED>"

// Java: UseSiteRenamer(Consumer<Node>)
func NewUseSiteRenamer(mutationCb func(rhino.Node)) *UseSiteRenamer {
	return &UseSiteRenamer{
		mutationCb:         mutationCb,
		trackRenamingIndex: true,
	}
}

func (r *UseSiteRenamer) recordRenaming(propName, newName string) {
	if r.renamingEntries == nil {
		r.renamingEntries = make(map[string][]string)
	}
	vals, exists := r.renamingEntries[propName]
	if !exists {
		r.renamingOrder = append(r.renamingOrder, propName)
		r.renamingEntries[propName] = []string{newName}
		return
	}
	if !slices.Contains(vals, newName) {
		r.renamingEntries[propName] = append(vals, newName)
	}
}

func (r *UseSiteRenamer) allocClusterName(origRep *ColorGraphNode, propName string, rep *ColorGraphNode) string {
	if origRep == rep {
		return propName
	}
	return "JSC$" + strconv.Itoa(int(rep.GetIndex())) + "_" + propName
}

// RenameUses renames all references to prop.
//
// If prop is invalid or should otherwise not be renamed, the AST will not
// be changed.
//
// Java: UseSiteRenamer.renameUses(PropertyClustering)
func (r *UseSiteRenamer) RenameUses(prop *PropertyClustering) UseSiteRenamer_RenameUsesResult {
	if prop.IsInvalidated() {
		if r.trackRenamingIndex {
			r.recordRenaming(prop.GetName(), useSiteRenamer_INVALIDATED_NAME_VALUE)
		}
		return UseSiteRenamer_RenameUsesResult_INVALIDATED
	}

	clusters := prop.GetClusters()
	if clusters.NumSets() <= 1 {
		// Don't bother renaming clusters with a single element. Renaming
		// won't actaully disambiguate anything in this case, so skip the
		// work.
		if r.trackRenamingIndex {
			propName := prop.GetName()
			r.recordRenaming(propName, propName)
		}
		return UseSiteRenamer_RenameUsesResult_ONLY_ONE_CLUSTER
	}

	propName := prop.GetName()
	reps := clusters.AppendRepresentatives(r.scratchReps[:0])
	r.scratchReps = reps
	if r.scratchClusterNames == nil {
		r.scratchClusterNames = make(map[*ColorGraphNode]string, max(16, len(reps)))
	} else {
		clear(r.scratchClusterNames)
	}
	origRep := prop.GetOriginalNameClusterRep()
	for _, rep := range reps {
		name := r.allocClusterName(origRep, propName, rep)
		r.scratchClusterNames[rep] = name
		if r.trackRenamingIndex {
			r.recordRenaming(propName, name)
		}
	}

	sites := prop.useSites.sites
	origID := prop.nameID
	for i := range sites {
		site := sites[i].site
		flatRep := clusters.Find(sites[i].owner)
		if flatRep == origRep && site.GetStringID() == origID {
			continue
		}
		newName := r.scratchClusterNames[flatRep]
		if newName != site.GetString() {
			site.SetString(newName)
			r.mutationCb(site)
		}
	}
	return UseSiteRenamer_RenameUsesResult_DISAMBIGUATED
}

// UseSiteRenamer_RenameUsesResult is UseSiteRenamer.RenameUsesResult.
type UseSiteRenamer_RenameUsesResult int32

const (
	UseSiteRenamer_RenameUsesResult_INVALIDATED UseSiteRenamer_RenameUsesResult = iota
	UseSiteRenamer_RenameUsesResult_ONLY_ONE_CLUSTER
	UseSiteRenamer_RenameUsesResult_DISAMBIGUATED
)

var useSiteRenamer_RenameUsesResult_names = [...]string{"INVALIDATED", "ONLY_ONE_CLUSTER", "DISAMBIGUATED"}

// Name is Enum.name().
func (r UseSiteRenamer_RenameUsesResult) Name() string {
	return useSiteRenamer_RenameUsesResult_names[r]
}

// String is Enum.toString().
func (r UseSiteRenamer_RenameUsesResult) String() string { return r.Name() }

// Ordinal is Enum.ordinal().
func (r UseSiteRenamer_RenameUsesResult) Ordinal() int { return int(r) }

// UseSiteNameSlice holds the renamed cluster names for a single property.
type UseSiteNameSlice []string

// ToSlice returns a copy of the renamed cluster names.
func (s UseSiteNameSlice) ToSlice() []string { return slices.Clone(s) }

// UseSiteRenamingIndex maps property names to their renamed cluster names in insertion order.
type UseSiteRenamingIndex struct {
	entries map[string][]string
	order   []string
}

// Get returns the renamed cluster names for name (empty when absent).
func (idx *UseSiteRenamingIndex) Get(name string) UseSiteNameSlice {
	if idx == nil {
		return nil
	}
	return idx.entries[name]
}

// Entries iterates (propertyName, renamedName) pairs in insertion order.
func (idx *UseSiteRenamingIndex) Entries() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		if idx == nil {
			return
		}
		for _, k := range idx.order {
			for _, v := range idx.entries[k] {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// GetRenamingIndex is renamingIndex.build().
//
// Java: UseSiteRenamer.getRenamingIndex()
func (r *UseSiteRenamer) GetRenamingIndex() *UseSiteRenamingIndex {
	if len(r.renamingOrder) == 0 {
		return &UseSiteRenamingIndex{}
	}
	entries := maps.Clone(r.renamingEntries)
	for _, k := range r.renamingOrder {
		entries[k] = slices.Clone(entries[k])
	}
	return &UseSiteRenamingIndex{
		entries: entries,
		order:   slices.Clone(r.renamingOrder),
	}
}

// useSiteRenamer_createAllClusterNames creates a unique name for each
// cluster in prop and maps it to the cluster representative.
//
// Java: UseSiteRenamer.createAllClusterNames(PropertyClustering)
func useSiteRenamer_createAllClusterNames(prop *PropertyClustering) *jutil.LinkedHashMap[*ColorGraphNode, string] {
	out := jutil.NewLinkedHashMap[*ColorGraphNode, string]()
	for r := range prop.GetClusters().AllRepresentatives().All() {
		out.Put(r, useSiteRenamer_createClusterName(prop, r))
	}
	return out
}

// Java: UseSiteRenamer.createClusterName(PropertyClustering, ColorGraphNode)
func useSiteRenamer_createClusterName(prop *PropertyClustering, rep *ColorGraphNode) string {
	if prop.GetOriginalNameClusterRep() == rep {
		return prop.GetName()
	}

	return "JSC$" + strconv.Itoa(rep.GetIndex()) + "_" + prop.GetName()
}
