// Port of com/google/javascript/jscomp/AccessorSummary.java @ closure-compiler v20261021

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
)

// AccessorSummary is a strongly typed view of information about getters and
// setters collected from the AST.
type AccessorSummary struct {
	assumeAlwaysGetterAndSetter bool
	accessors                   *jutil.LinkedHashMap[string, AccessorSummary_PropertyAccessKind]
	byName                      map[string]AccessorSummary_PropertyAccessKind
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

// accessorSummary_PropertyAccessKind_flags is the `flags` field of each
// constant.
var accessorSummary_PropertyAccessKind_flags = [...]int8{0, 1, 2, 3}

// AccessorSummary_PropertyAccessKind_values is PropertyAccessKind.values().
func AccessorSummary_PropertyAccessKind_values() []AccessorSummary_PropertyAccessKind {
	return []AccessorSummary_PropertyAccessKind{
		AccessorSummary_PropertyAccessKind_NORMAL,
		AccessorSummary_PropertyAccessKind_GETTER_ONLY,
		AccessorSummary_PropertyAccessKind_SETTER_ONLY,
		AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER,
	}
}

// Name is Enum.name().
func (k AccessorSummary_PropertyAccessKind) Name() string {
	return accessorSummary_PropertyAccessKind_names[k]
}

// String is Enum.toString().
func (k AccessorSummary_PropertyAccessKind) String() string { return k.Name() }

// Ordinal is Enum.ordinal().
func (k AccessorSummary_PropertyAccessKind) Ordinal() int { return int(k) }

// Java: AccessorSummary.PropertyAccessKind.hasGetter()
func (k AccessorSummary_PropertyAccessKind) HasGetter() bool {
	return (accessorSummary_PropertyAccessKind_flags[k] & 1) != 0
}

// Java: AccessorSummary.PropertyAccessKind.hasSetter()
func (k AccessorSummary_PropertyAccessKind) HasSetter() bool {
	return (accessorSummary_PropertyAccessKind_flags[k] & 2) != 0
}

// Java: AccessorSummary.PropertyAccessKind.hasGetterOrSetter()
func (k AccessorSummary_PropertyAccessKind) HasGetterOrSetter() bool {
	return (accessorSummary_PropertyAccessKind_flags[k] & 3) != 0
}

// UnionWith is used to combine information from externs and from sources.
//
// Java: AccessorSummary.PropertyAccessKind.unionWith(PropertyAccessKind)
func (k AccessorSummary_PropertyAccessKind) UnionWith(other AccessorSummary_PropertyAccessKind) AccessorSummary_PropertyAccessKind {
	combinedFlags := int32(accessorSummary_PropertyAccessKind_flags[k] | accessorSummary_PropertyAccessKind_flags[other])
	switch combinedFlags {
	case 0:
		return AccessorSummary_PropertyAccessKind_NORMAL
	case 1:
		return AccessorSummary_PropertyAccessKind_GETTER_ONLY
	case 2:
		return AccessorSummary_PropertyAccessKind_SETTER_ONLY
	case 3:
		return AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER
	default:
		panic(jutil.NewIllegalStateExceptionMsg("unexpected value: " + jutil.ToJavaString(combinedFlags)))
	}
}

// AccessorSummary_Create is AccessorSummary.create(Map).
//
// TODO(nickreid): Efficiently verify that no entry in `accessor` is
// `NORMAL`.
//
// Java: AccessorSummary.create(Map<String, PropertyAccessKind>)
func AccessorSummary_Create(accessors *jutil.LinkedHashMap[string, AccessorSummary_PropertyAccessKind]) *AccessorSummary {
	return newAccessorSummary(jutil.NewLinkedHashMapFrom(accessors))
}

// Java: AccessorSummary(ImmutableMap<String, PropertyAccessKind>)
func newAccessorSummary(accessors *jutil.LinkedHashMap[string, AccessorSummary_PropertyAccessKind]) *AccessorSummary {
	var byName map[string]AccessorSummary_PropertyAccessKind
	if !accessors.IsEmpty() {
		byName = make(map[string]AccessorSummary_PropertyAccessKind, accessors.Size())
		for name, kind := range accessors.All() {
			byName[name] = kind
		}
	}
	return &AccessorSummary{accessors: accessors, byName: byName, assumeAlwaysGetterAndSetter: false}
}

// Java: AccessorSummary(boolean)
func newAccessorSummaryAssuming(assumeAlwaysGetterAndSetter bool) *AccessorSummary {
	return &AccessorSummary{
		accessors:                   jutil.NewLinkedHashMap[string, AccessorSummary_PropertyAccessKind](),
		assumeAlwaysGetterAndSetter: assumeAlwaysGetterAndSetter,
	}
}

// Java: AccessorSummary.getAccessors()
func (s *AccessorSummary) GetAccessors() *jutil.LinkedHashMap[string, AccessorSummary_PropertyAccessKind] {
	return s.accessors
}

// Java: AccessorSummary.getKind(String)
func (s *AccessorSummary) GetKind(name string) AccessorSummary_PropertyAccessKind {
	if s.assumeAlwaysGetterAndSetter {
		return AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER
	}
	return s.byName[name]
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
	builder :=
		AccessorSummaryProto_NewBuilder().
			SetAssumeAlwaysGetterAndSetter(s.assumeAlwaysGetterAndSetter)
	if s.assumeAlwaysGetterAndSetter {
		jutil.CheckState(s.accessors.IsEmpty())
	} else {
		for name, kind := range s.accessors.All() {
			builder.AddAccessors(
				AccessorSummaryEntryProto_NewBuilder().
					SetName(name).
					SetKind(AccessorSummary_ToProto(kind)).
					Build())
		}
	}
	return builder.Build()
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
	builder := jutil.NewImmutableMapBuilder[string, AccessorSummary_PropertyAccessKind]()
	for _, entry := range proto.GetAccessorsList() {
		builder.Put(entry.GetName(), AccessorSummary_FromProtoKind(entry.GetKind()))
	}
	return AccessorSummary_Create(builder.BuildOrThrow())
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
