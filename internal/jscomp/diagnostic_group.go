// Port of com/google/javascript/jscomp/DiagnosticGroup.java @ closure-compiler v20261021

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
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

const diagnosticGroupMapThreshold = 8

var nextDiagnosticGroupID atomic.Uint32

// DiagnosticTypeSlice is an ordered, deduplicated slice of *DiagnosticType
// values belonging to a DiagnosticGroup.
type DiagnosticTypeSlice []*DiagnosticType

func (s DiagnosticTypeSlice) All() iter.Seq[*DiagnosticType] {
	return func(yield func(*DiagnosticType) bool) {
		for _, t := range s {
			if !yield(t) {
				return
			}
		}
	}
}

func (s DiagnosticTypeSlice) Size() int {
	return len(s)
}

func (s DiagnosticTypeSlice) IsEmpty() bool {
	return len(s) == 0
}

func (s DiagnosticTypeSlice) Contains(type_ *DiagnosticType) bool {
	if type_ == nil {
		return false
	}
	for _, t := range s {
		if t == type_ || t.key == type_.key {
			return true
		}
	}
	return false
}

// DiagnosticGroup groups a set of related diagnostic types together, so that
// they can be toggled on and off as one unit.
type DiagnosticGroup struct {
	// The deduplicated types represented by this group, in insertion order.
	types DiagnosticTypeSlice

	// Allocated only when len(types) > diagnosticGroupMapThreshold; keyed by DiagnosticType.key.
	typeKeySet map[string]struct{}

	// A human-readable name for the group. May be nil.
	name *string

	id uint32
}

const diagnosticGroup_serialVersionUID int64 = 1

// NewDiagnosticGroupWithName creates a group that matches all errors of the
// given types.
//
// Java: DiagnosticGroup(String, DiagnosticType...)
func NewDiagnosticGroupWithName(name *string, types ...*DiagnosticType) *DiagnosticGroup {
	if len(types) <= diagnosticGroupMapThreshold {
		dedup := make(DiagnosticTypeSlice, 0, len(types))
		for _, t := range types {
			jutil.CheckNotNull(t)
			dup := false
			for _, existing := range dedup {
				if existing == t || existing.key == t.key {
					dup = true
					break
				}
			}
			if !dup {
				dedup = append(dedup, t)
			}
		}
		return &DiagnosticGroup{
			id:    nextDiagnosticGroupID.Add(1),
			name:  name,
			types: dedup,
		}
	}

	keySet := make(map[string]struct{}, len(types))
	dedup := make(DiagnosticTypeSlice, 0, len(types))
	for _, t := range types {
		jutil.CheckNotNull(t)
		if _, exists := keySet[t.key]; !exists {
			keySet[t.key] = struct{}{}
			dedup = append(dedup, t)
		}
	}
	var storedSet map[string]struct{}
	if len(dedup) > diagnosticGroupMapThreshold {
		storedSet = keySet
	}
	return &DiagnosticGroup{
		id:         nextDiagnosticGroupID.Add(1),
		name:       name,
		types:      dedup,
		typeKeySet: storedSet,
	}
}

// NewDiagnosticGroup creates a group that matches all errors of the given
// types.
//
// Java: DiagnosticGroup(DiagnosticType...)
func NewDiagnosticGroup(types ...*DiagnosticType) *DiagnosticGroup {
	return NewDiagnosticGroupWithName(nil, types...)
}

// newDiagnosticGroupForSingleType creates a diagnostic group with no name
// that only matches the given type.
//
// Java: DiagnosticGroup(DiagnosticType)
func newDiagnosticGroupForSingleType(type_ *DiagnosticType) *DiagnosticGroup {
	jutil.CheckNotNull(type_)
	return &DiagnosticGroup{
		id:    nextDiagnosticGroupID.Add(1),
		types: DiagnosticTypeSlice{type_},
	}
}

// NewDiagnosticGroupFromGroups creates a composite group.
//
// Java: DiagnosticGroup(DiagnosticGroup...)
func NewDiagnosticGroupFromGroups(groups ...*DiagnosticGroup) *DiagnosticGroup {
	return NewDiagnosticGroupWithNameFromGroups(nil, groups...)
}

// NewDiagnosticGroupWithNameFromGroups creates a composite group.
//
// Java: DiagnosticGroup(String, DiagnosticGroup...)
func NewDiagnosticGroupWithNameFromGroups(name *string, groups ...*DiagnosticGroup) *DiagnosticGroup {
	total := 0
	for _, group := range groups {
		jutil.CheckNotNull(group)
		total += len(group.types)
	}
	if total <= diagnosticGroupMapThreshold {
		dedup := make(DiagnosticTypeSlice, 0, total)
		for _, group := range groups {
			for _, t := range group.types {
				dup := false
				for _, existing := range dedup {
					if existing == t || existing.key == t.key {
						dup = true
						break
					}
				}
				if !dup {
					dedup = append(dedup, t)
				}
			}
		}
		return &DiagnosticGroup{
			id:    nextDiagnosticGroupID.Add(1),
			name:  name,
			types: dedup,
		}
	}

	keySet := make(map[string]struct{}, total)
	dedup := make(DiagnosticTypeSlice, 0, total)
	for _, group := range groups {
		for _, t := range group.types {
			if _, exists := keySet[t.key]; !exists {
				keySet[t.key] = struct{}{}
				dedup = append(dedup, t)
			}
		}
	}
	var storedSet map[string]struct{}
	if len(dedup) > diagnosticGroupMapThreshold {
		storedSet = keySet
	}
	return &DiagnosticGroup{
		id:         nextDiagnosticGroupID.Add(1),
		name:       name,
		types:      dedup,
		typeKeySet: storedSet,
	}
}

// DiagnosticGroups with only a single DiagnosticType, keyed by DiagnosticType.key.
var diagnosticGroup_singletons = make(map[string]*DiagnosticGroup)

// Go-only: language-adapter — guards diagnosticGroup_singletons.
var diagnosticGroup_mu sync.Mutex

// DiagnosticGroup_ForType creates a diagnostic group that matches only the
// given type.
//
// Java: DiagnosticGroup.forType(DiagnosticType)
func DiagnosticGroup_ForType(type_ *DiagnosticType) *DiagnosticGroup {
	jutil.CheckNotNull(type_)
	if g := type_.singletonGroup.Load(); g != nil {
		return g
	}
	diagnosticGroup_mu.Lock()
	g := diagnosticGroup_singletons[type_.key]
	if g == nil {
		g = newDiagnosticGroupForSingleType(type_)
		diagnosticGroup_singletons[type_.key] = g
	}
	diagnosticGroup_mu.Unlock()
	type_.singletonGroup.Store(g)
	return g
}

// Matches returns whether the given error's type matches a type in this
// group.
//
// Java: DiagnosticGroup.matches(JSError)
func (g *DiagnosticGroup) Matches(error *JSError) bool {
	return g.MatchesType(error.Type())
}

// MatchesType returns whether the given type matches a type in this group.
//
// Java: DiagnosticGroup.matches(DiagnosticType)
func (g *DiagnosticGroup) MatchesType(type_ *DiagnosticType) bool {
	if type_ == nil {
		return false
	}
	if g.typeKeySet != nil {
		_, ok := g.typeKeySet[type_.key]
		return ok
	}
	for _, t := range g.types {
		if t == type_ || t.key == type_.key {
			return true
		}
	}
	return false
}

// GetTypes returns all the types in this group.
//
// Java: DiagnosticGroup.getTypes()
func (g *DiagnosticGroup) GetTypes() DiagnosticTypeSlice {
	return g.types
}

// GetName may return nil.
//
// Java: DiagnosticGroup.getName()
func (g *DiagnosticGroup) GetName() *string {
	return g.name
}

// Java: DiagnosticGroup.toString()
func (g *DiagnosticGroup) String() string {
	if g.name == nil {
		return jutil.ObjectToString("com.google.javascript.jscomp.DiagnosticGroup", g)
	}
	return "DiagnosticGroup<" + *g.name + ">"
}

// Equals reports whether other refers to the same DiagnosticGroup, including
// across value copies that preserve the group's creation ID.
func (g *DiagnosticGroup) Equals(other any) bool {
	switch o := other.(type) {
	case *DiagnosticGroup:
		if g == nil || o == nil {
			return g == o
		}
		return g == o || (g.id != 0 && g.id == o.id)
	case DiagnosticGroup:
		if g == nil {
			return false
		}
		return g.id != 0 && g.id == o.id
	default:
		return false
	}
}

// HashCode returns a stable hash code derived from the group's creation ID so
// value-copied DiagnosticGroups hash identically to the original.
func (g *DiagnosticGroup) HashCode() int32 {
	if g == nil {
		return 0
	}
	if g.id != 0 {
		return int32(g.id)
	}
	return jutil.IdentityHashCodeOfPointer(g)
}
