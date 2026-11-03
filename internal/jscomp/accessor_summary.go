// Port of com/google/javascript/jscomp/AccessorSummary.java @ closure-compiler v20261028

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// AccessorSummary is a strongly typed view of information about getters and
// setters collected from the AST.
type AccessorSummary struct {
	alwaysGetterAndSetter bool
	hasAnyAccessors       bool
	kindsByID             []uint64 // 2 bits per strpool.ID (32 IDs per uint64)
	accessorIDs           []uint32 // insertion-ordered strpool.IDs
}

// AccessorSummary_PropertyAccessKind indicates whether a property has a
// getter or a setter, or both.
type AccessorSummary_PropertyAccessKind int32

// To save space properties without getters or setters won't appear in the
// maps at all, but NORMAL will be returned by some methods.
const (
	AccessorSummary_PropertyAccessKind_NORMAL AccessorSummary_PropertyAccessKind = iota
	AccessorSummary_PropertyAccessKind_GETTER_ONLY
	AccessorSummary_PropertyAccessKind_SETTER_ONLY
	AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER
)

var accessorSummary_PropertyAccessKind_names = [...]string{
	"NORMAL", "GETTER_ONLY", "SETTER_ONLY", "GETTER_AND_SETTER",
}

// Name is Enum.name().
func (k AccessorSummary_PropertyAccessKind) Name() string {
	return accessorSummary_PropertyAccessKind_names[k]
}

// String is Enum.toString().
func (k AccessorSummary_PropertyAccessKind) String() string { return k.Name() }

// Java: AccessorSummary.PropertyAccessKind.hasGetter()
func (k AccessorSummary_PropertyAccessKind) HasGetter() bool {
	return (k & AccessorSummary_PropertyAccessKind_GETTER_ONLY) != 0
}

// Java: AccessorSummary.PropertyAccessKind.hasSetter()
func (k AccessorSummary_PropertyAccessKind) HasSetter() bool {
	return (k & AccessorSummary_PropertyAccessKind_SETTER_ONLY) != 0
}

// Java: AccessorSummary.PropertyAccessKind.hasGetterOrSetter()
func (k AccessorSummary_PropertyAccessKind) HasGetterOrSetter() bool {
	return (k & AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER) != 0
}

// UnionWith is used to combine information from externs and from sources.
//
// Java: AccessorSummary.PropertyAccessKind.unionWith(PropertyAccessKind)
func (k AccessorSummary_PropertyAccessKind) UnionWith(other AccessorSummary_PropertyAccessKind) AccessorSummary_PropertyAccessKind {
	return k | other
}

// AccessorSummary_Create is AccessorSummary.create(Map).
//
// TODO(nickreid): Efficiently verify that no entry in `accessor` is
// `NORMAL`.
//
// Java: AccessorSummary.create(Map<String, PropertyAccessKind>)
func AccessorSummary_Create(accessors jutil.MapSource[string, AccessorSummary_PropertyAccessKind]) *AccessorSummary {
	jutil.CheckNotNull(accessors)
	n := accessors.Size()
	if n == 0 {
		return &AccessorSummary{}
	}
	entries := make([]accessorIDEntry, 0, n)
	for name, kind := range accessors.All() {
		entries = append(entries, accessorIDEntry{
			id:   uint32(strpool.Intern(name)),
			kind: kind,
		})
	}
	return newAccessorSummaryFromIDs([][]accessorIDEntry{entries})
}

type accessorIDEntry struct {
	id   uint32
	kind AccessorSummary_PropertyAccessKind
}

func newAccessorSummaryFromIDs(batches [][]accessorIDEntry) *AccessorSummary {
	var maxID uint32
	totalEntries := 0
	hasAny := false
	for _, batch := range batches {
		totalEntries += len(batch)
		for _, e := range batch {
			if e.kind != AccessorSummary_PropertyAccessKind_NORMAL {
				hasAny = true
				if e.id > maxID {
					maxID = e.id
				}
			}
		}
	}
	if !hasAny {
		return &AccessorSummary{}
	}
	kindsByID := make([]uint64, int(maxID>>5)+1)
	capHint := totalEntries
	if int(maxID)+1 < capHint {
		capHint = int(maxID) + 1
	}
	accessorIDs := make([]uint32, 0, capHint)
	for _, batch := range batches {
		for _, e := range batch {
			if e.kind == AccessorSummary_PropertyAccessKind_NORMAL {
				continue
			}
			wordIdx := e.id >> 5
			shift := (e.id & 31) * 2
			word := kindsByID[wordIdx]
			if (word>>shift)&3 == 0 {
				accessorIDs = append(accessorIDs, e.id)
			}
			kindsByID[wordIdx] = word | ((uint64(e.kind) & 3) << shift)
		}
	}
	return &AccessorSummary{
		hasAnyAccessors: true,
		kindsByID:       kindsByID,
		accessorIDs:     accessorIDs,
	}
}

// Java: AccessorSummary(boolean)
func newAccessorSummaryAssuming(assumeAlwaysGetterAndSetter bool) *AccessorSummary {
	return &AccessorSummary{
		alwaysGetterAndSetter: assumeAlwaysGetterAndSetter,
	}
}

// AccessorIDs returns the interned strpool.IDs of recorded accessors in insertion order.
func (s *AccessorSummary) AccessorIDs() []uint32 {
	return s.accessorIDs
}

// Java: AccessorSummary.getAccessors()
func (s *AccessorSummary) GetAccessors() *OrderedStringMap[AccessorSummary_PropertyAccessKind] {
	m := NewOrderedStringMapWithCapacity[AccessorSummary_PropertyAccessKind](len(s.accessorIDs))
	for _, id := range s.accessorIDs {
		m.PutNew(strpool.String(strpool.ID(id)), s.GetKindByID(id))
	}
	return m
}

// GetKindByID returns the PropertyAccessKind for an interned strpool.ID.
func (s *AccessorSummary) GetKindByID(id uint32) AccessorSummary_PropertyAccessKind {
	if s.alwaysGetterAndSetter {
		return AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER
	}
	if !s.hasAnyAccessors {
		return AccessorSummary_PropertyAccessKind_NORMAL
	}
	wordIdx := int(id >> 5)
	if wordIdx >= len(s.kindsByID) {
		return AccessorSummary_PropertyAccessKind_NORMAL
	}
	return AccessorSummary_PropertyAccessKind((s.kindsByID[wordIdx] >> ((id & 31) * 2)) & 3)
}

// Java: AccessorSummary.getKind(String)
func (s *AccessorSummary) GetKind(name string) AccessorSummary_PropertyAccessKind {
	if s.alwaysGetterAndSetter {
		return AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER
	}
	if !s.hasAnyAccessors {
		return AccessorSummary_PropertyAccessKind_NORMAL
	}
	id := strpool.Lookup(name)
	if id == 0 {
		return AccessorSummary_PropertyAccessKind_NORMAL
	}
	return s.GetKindByID(uint32(id))
}

// AccessorSummary_CreateAssumingAlwaysGetterAndSetter returns an accessor
// summary that assumes every access is a potential getter or setter.
//
// Java: AccessorSummary.createAssumingAlwaysGetterAndSetter()
func AccessorSummary_CreateAssumingAlwaysGetterAndSetter() *AccessorSummary {
	return newAccessorSummaryAssuming(true)
}

// Java: AccessorSummary.toProto()
func (s *AccessorSummary) ToProto() *AccessorSummaryProto {
	if s.alwaysGetterAndSetter {
		jutil.CheckState(len(s.accessorIDs) == 0)
		return &AccessorSummaryProto{assumeAlwaysGetterAndSetter_: true}
	}
	var accessors []*AccessorSummaryEntryProto
	if len(s.accessorIDs) > 0 {
		accessors = make([]*AccessorSummaryEntryProto, len(s.accessorIDs))
		for i, id := range s.accessorIDs {
			accessors[i] = &AccessorSummaryEntryProto{
				name_: strpool.String(strpool.ID(id)),
				kind_: int32(AccessorSummary_ToProto(s.GetKindByID(id))),
			}
		}
	}
	return &AccessorSummaryProto{accessors_: accessors}
}

// Java: AccessorSummary.toProto(PropertyAccessKind)
func AccessorSummary_ToProto(kind AccessorSummary_PropertyAccessKind) PropertyAccessKindProto {
	switch kind {
	case AccessorSummary_PropertyAccessKind_NORMAL:
		return PropertyAccessKindProto_KIND_NORMAL
	case AccessorSummary_PropertyAccessKind_GETTER_ONLY:
		return PropertyAccessKindProto_KIND_GETTER_ONLY
	case AccessorSummary_PropertyAccessKind_SETTER_ONLY:
		return PropertyAccessKindProto_KIND_SETTER_ONLY
	case AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER:
		return PropertyAccessKindProto_KIND_GETTER_AND_SETTER
	}
	// Go-only: language-adapter — unreachable for a valid enum value: the
	// Java switch expression is exhaustive.
	panic(jutil.NewIllegalStateException())
}

// Java: AccessorSummary.fromProto(AccessorSummaryProto)
func AccessorSummary_FromProto(proto *AccessorSummaryProto) *AccessorSummary {
	if proto.GetAssumeAlwaysGetterAndSetter() {
		return AccessorSummary_CreateAssumingAlwaysGetterAndSetter()
	}
	protoList := proto.GetAccessorsList()
	if len(protoList) == 0 {
		return &AccessorSummary{}
	}
	entries := make([]accessorIDEntry, len(protoList))
	for i, entry := range protoList {
		entries[i] = accessorIDEntry{
			id:   uint32(strpool.Intern(entry.GetName())),
			kind: AccessorSummary_FromProtoKind(entry.GetKind()),
		}
	}
	return newAccessorSummaryFromIDs([][]accessorIDEntry{entries})
}

// Java: AccessorSummary.fromProto(PropertyAccessKindProto)
func AccessorSummary_FromProtoKind(kind PropertyAccessKindProto) AccessorSummary_PropertyAccessKind {
	switch kind {
	case PropertyAccessKindProto_KIND_NORMAL:
		return AccessorSummary_PropertyAccessKind_NORMAL
	case PropertyAccessKindProto_KIND_GETTER_ONLY:
		return AccessorSummary_PropertyAccessKind_GETTER_ONLY
	case PropertyAccessKindProto_KIND_SETTER_ONLY:
		return AccessorSummary_PropertyAccessKind_SETTER_ONLY
	case PropertyAccessKindProto_KIND_GETTER_AND_SETTER:
		return AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("Unknown kind: " + kind.String()))
	}
}
