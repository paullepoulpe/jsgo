// Port of com/google/javascript/jscomp/disambiguate/UseSiteRenamer.java @ closure-compiler v20261015

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
	"strconv"
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// UseSiteRenamer applies renaming to property use sites following cluster
// computation.
type UseSiteRenamer struct {
	mutationCb func(*rhino.Node)

	// ImmutableSetMultimap.Builder: keys and values in insertion order.
	renamingIndex       *jutil.SetMultimap[string, string]
	trackRenamingIndex  bool
	scratchReps         []*ColorGraphNode
	scratchClusterNames map[*ColorGraphNode]string
	nameByteSlab        []byte
}

const useSiteRenamer_INVALIDATED_NAME_VALUE = "<INVALIDATED>"

// Java: UseSiteRenamer(Consumer<Node>)
func NewUseSiteRenamer(mutationCb func(*rhino.Node)) *UseSiteRenamer {
	return &UseSiteRenamer{
		mutationCb:         mutationCb,
		renamingIndex:      jutil.NewSetMultimap[string, string](),
		trackRenamingIndex: true,
	}
}

func (r *UseSiteRenamer) allocClusterName(origRep *ColorGraphNode, propName string, rep *ColorGraphNode) string {
	if origRep == rep {
		return propName
	}
	// Format "JSC$" + strconv.Itoa(rep.GetIndex()) + "_" + propName into noscan slab.
	maxLen := 4 + 20 + 1 + len(propName)
	if len(r.nameByteSlab) < maxLen {
		r.nameByteSlab = make([]byte, max(maxLen, 16384))
	}
	buf := r.nameByteSlab[:0]
	buf = append(buf, "JSC$"...)
	buf = strconv.AppendInt(buf, int64(rep.GetIndex()), 10)
	buf = append(buf, '_')
	buf = append(buf, propName...)
	n := len(buf)
	s := unsafe.String(&r.nameByteSlab[0], n)
	r.nameByteSlab = r.nameByteSlab[n:]
	return s
}

// RenameUses renames all references to prop.
//
// If prop is invalid or should otherwise not be renamed, the AST will not
// be changed.
//
// Java: UseSiteRenamer.renameUses(PropertyClustering)
func (r *UseSiteRenamer) RenameUses(prop *PropertyClustering) UseSiteRenamer_RenameUsesResult {
	propName := prop.GetName()
	if prop.IsInvalidated() {
		if r.trackRenamingIndex {
			r.renamingIndex.Put(propName, useSiteRenamer_INVALIDATED_NAME_VALUE)
		}
		return UseSiteRenamer_RenameUsesResult_INVALIDATED
	}

	clusters := prop.GetClusters()
	if clusters.NumSets() <= 1 {
		// Don't bother renaming clusters with a single element. Renaming
		// won't actaully disambiguate anything in this case, so skip the
		// work.
		if r.trackRenamingIndex {
			r.renamingIndex.Put(propName, propName)
		}
		return UseSiteRenamer_RenameUsesResult_ONLY_ONE_CLUSTER
	}

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
			r.renamingIndex.Put(propName, name)
		}
	}

	useSites := prop.GetUseSites()
	nInline := int(useSites.count)
	if nInline > len(useSites.inline) {
		nInline = len(useSites.inline)
	}
	for i := 0; i < nInline; i++ {
		e := useSites.inline[i]
		flatRep := clusters.Find(e.owner)
		newName := r.scratchClusterNames[flatRep]
		if newName != e.site.GetString() {
			e.site.SetString(newName)
			r.mutationCb(e.site)
		}
	}
	for ch := useSites.head; ch != nil; ch = ch.next {
		entries := ch.entries
		if ch == useSites.tail {
			entries = entries[:useSites.tailUsed]
		}
		for i := range entries {
			e := entries[i]
			flatRep := clusters.Find(e.owner)
			newName := r.scratchClusterNames[flatRep]
			if newName != e.site.GetString() {
				e.site.SetString(newName)
				r.mutationCb(e.site)
			}
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

// UseSiteRenamer_RenameUsesResult_values is RenameUsesResult.values().
func UseSiteRenamer_RenameUsesResult_values() []UseSiteRenamer_RenameUsesResult {
	return []UseSiteRenamer_RenameUsesResult{
		UseSiteRenamer_RenameUsesResult_INVALIDATED,
		UseSiteRenamer_RenameUsesResult_ONLY_ONE_CLUSTER,
		UseSiteRenamer_RenameUsesResult_DISAMBIGUATED,
	}
}

// Name is Enum.name().
func (r UseSiteRenamer_RenameUsesResult) Name() string {
	return useSiteRenamer_RenameUsesResult_names[r]
}

// String is Enum.toString().
func (r UseSiteRenamer_RenameUsesResult) String() string { return r.Name() }

// Ordinal is Enum.ordinal().
func (r UseSiteRenamer_RenameUsesResult) Ordinal() int { return int(r) }

// GetRenamingIndex is renamingIndex.build().
//
// Java: UseSiteRenamer.getRenamingIndex()
func (r *UseSiteRenamer) GetRenamingIndex() *jutil.SetMultimap[string, string] {
	return jutil.SetMultimapCopyOf(r.renamingIndex)
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

	return "JSC$" + jutil.ToJavaString(rep.GetIndex()) + "_" + prop.GetName()
}
