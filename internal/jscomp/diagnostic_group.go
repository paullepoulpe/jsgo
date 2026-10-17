// Port of com/google/javascript/jscomp/DiagnosticGroup.java @ closure-compiler v20260415

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
)

// DiagnosticGroup groups a set of related diagnostic types together, so that
// they can be toggled on and off as one unit.
type DiagnosticGroup struct {
	// The set of types represented by this group, hashed by key.
	types *jutil.ObjectLinkedHashSet[*DiagnosticType]

	// A human-readable name for the group. May be nil.
	name *string
}

const diagnosticGroup_serialVersionUID int64 = 1

// NewDiagnosticGroupWithName creates a group that matches all errors of the
// given types.
//
// Java: DiagnosticGroup(String, DiagnosticType...)
func NewDiagnosticGroupWithName(name *string, types ...*DiagnosticType) *DiagnosticGroup {
	g := &DiagnosticGroup{name: name, types: jutil.NewObjectLinkedHashSet[*DiagnosticType]()}
	for _, t := range types {
		g.types.Add(jutil.CheckNotNull(t))
	}
	return g
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
	return &DiagnosticGroup{name: nil, types: jutil.NewObjectLinkedHashSet(jutil.CheckNotNull(type_))}
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
	set := jutil.NewObjectLinkedHashSet[*DiagnosticType]()

	for _, group := range groups {
		set.AddAll(group.types.All())
	}

	return &DiagnosticGroup{name: name, types: set}
}

// DiagnosticGroups with only a single DiagnosticType.
var diagnosticGroup_singletons = jutil.NewObjectLinkedHashMap[*DiagnosticType, *DiagnosticGroup]()

// Go-only: language-adapter — the monitor of DiagnosticGroup.class, which
// its static synchronized methods hold.
var diagnosticGroup_classMonitor jutil.Monitor

// DiagnosticGroup_ForType creates a diagnostic group that matches only the
// given type.
//
// Java: DiagnosticGroup.forType(DiagnosticType)
func DiagnosticGroup_ForType(type_ *DiagnosticType) *DiagnosticGroup {
	diagnosticGroup_classMonitor.Lock()
	defer diagnosticGroup_classMonitor.Unlock()
	if g, ok := diagnosticGroup_singletons.Lookup(type_); !ok || g == nil {
		diagnosticGroup_singletons.Put(type_, newDiagnosticGroupForSingleType(type_))
	}
	return diagnosticGroup_singletons.Get(type_)
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
	return g.types.Contains(type_)
}

// GetTypes returns all the types in this group.
//
// Java: DiagnosticGroup.getTypes()
func (g *DiagnosticGroup) GetTypes() *jutil.ObjectLinkedHashSet[*DiagnosticType] {
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
