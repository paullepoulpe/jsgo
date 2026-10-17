// Port of com/google/javascript/jscomp/ComposeWarningsGuard.java @ closure-compiler v20260415

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ComposeWarningsGuard is a WarningsGuard that represents just a chain of
// other guards. For example we could have following chain 1) all warnings
// outside of /foo/ should be suppressed 2) errors with key JSC_BAR should be
// marked as warning 3) the rest should be reported as error
//
// This class is designed for such behavior.
type ComposeWarningsGuard struct {
	WarningsGuardBase

	// The order that the guards were added in.
	orderOfAddition *jutil.LinkedHashMap[WarningsGuard, int]
	numberOfAdds    int
	guardComparator func(a, b WarningsGuard) int
	demoteErrors    bool

	// The order that the guards are applied in.
	guards *jutil.TreeSet[WarningsGuard]
}

var _ WarningsGuard = (*ComposeWarningsGuard)(nil)

const composeWarningsGuard_serialVersionUID int64 = 1

// composeWarningsGuard_GuardComparator is ComposeWarningsGuard.GuardComparator.
type composeWarningsGuard_GuardComparator struct {
	orderOfAddition *jutil.LinkedHashMap[WarningsGuard, int]
}

// Java: ComposeWarningsGuard.GuardComparator.compare(WarningsGuard, WarningsGuard)
func (c *composeWarningsGuard_GuardComparator) Compare(a, b WarningsGuard) int {
	priorityDiff := a.GetPriority() - b.GetPriority()
	if priorityDiff != 0 {
		return priorityDiff
	}
	// If the warnings guards have the same priority, the one that
	// was added last wins.
	return composeWarningsGuard_intValue(c.orderOfAddition, b) - composeWarningsGuard_intValue(c.orderOfAddition, a)
}

// composeWarningsGuard_intValue is orderOfAddition.get(g).intValue(), which
// throws NullPointerException for an unknown guard.
func composeWarningsGuard_intValue(m *jutil.LinkedHashMap[WarningsGuard, int], g WarningsGuard) int {
	v, ok := m.Lookup(g)
	if !ok {
		panic(jutil.NewNullPointerException())
	}
	return v
}

func newComposeWarningsGuardBase() *ComposeWarningsGuard {
	g := &ComposeWarningsGuard{orderOfAddition: jutil.NewLinkedHashMap[WarningsGuard, int]()}
	cmp := &composeWarningsGuard_GuardComparator{orderOfAddition: g.orderOfAddition}
	g.guardComparator = cmp.Compare
	g.guards = jutil.NewTreeSet(g.guardComparator)
	return g
}

// Java: ComposeWarningsGuard(List<WarningsGuard>)
func NewComposeWarningsGuardFromList(guards []WarningsGuard) *ComposeWarningsGuard {
	g := newComposeWarningsGuardBase()
	g.AddGuards(guards)
	return g
}

// Java: ComposeWarningsGuard(WarningsGuard...)
func NewComposeWarningsGuard(guards ...WarningsGuard) *ComposeWarningsGuard {
	for _, guard := range guards {
		jutil.CheckNotNull(guard) // ImmutableList.copyOf
	}
	return NewComposeWarningsGuardFromList(append([]WarningsGuard(nil), guards...))
}

// Java: ComposeWarningsGuard.addGuard(WarningsGuard)
func (g *ComposeWarningsGuard) AddGuard(guard WarningsGuard) {
	if composeGuard, ok := guard.(*ComposeWarningsGuard); ok {
		if composeGuard.demoteErrors {
			g.demoteErrors = composeGuard.demoteErrors
		}
		// Reverse the guards, so that they have the same order in the result.
		var reversed []WarningsGuard
		for x := range composeGuard.guards.Descending() {
			reversed = append(reversed, x)
		}
		g.AddGuards(reversed)
	} else {
		g.numberOfAdds++
		g.orderOfAddition.Put(guard, g.numberOfAdds)
		g.guards.Remove(guard)
		g.guards.Add(guard)
	}
}

// Java: ComposeWarningsGuard.addGuards(Iterable<WarningsGuard>)
func (g *ComposeWarningsGuard) AddGuards(guards []WarningsGuard) {
	for _, guard := range guards {
		g.AddGuard(guard)
	}
}

// Level may return nil.
//
// Java: ComposeWarningsGuard.level(JSError)
func (g *ComposeWarningsGuard) Level(error *JSError) *CheckLevel {
	for guard := range g.guards.All() {
		newLevel := guard.Level(error)
		if newLevel != nil {
			if g.demoteErrors && *newLevel == CheckLevel_ERROR {
				return CheckLevel_WARNING.Ptr()
			}
			return newLevel
		}
	}
	return nil
}

// Java: ComposeWarningsGuard.mustRunChecks(DiagnosticGroup)
func (g *ComposeWarningsGuard) MustRunChecks(group *DiagnosticGroup) base.Tri {
	// TODO(b/189635620): Merge these helper methods. Why are they asymmetric?
	enable := g.Enables(group)
	disable := g.Disables(group)
	jutil.CheckStateMsg2(!enable || !disable, "%s applied to %s", g, group)
	if enable {
		return base.Tri_TRUE
	} else if disable {
		return base.Tri_FALSE
	} else {
		return base.Tri_UNKNOWN
	}
}

// Java: ComposeWarningsGuard.disables(DiagnosticGroup)
func (g *ComposeWarningsGuard) Disables(group *DiagnosticGroup) bool {
nextSingleton:
	for type_ := range group.GetTypes().All() {
		singleton := DiagnosticGroup_ForType(type_)
		for guard := range g.guards.All() {
			switch guard.MustRunChecks(singleton) {
			case base.Tri_TRUE:
				return false
			case base.Tri_FALSE:
				continue nextSingleton
			case base.Tri_UNKNOWN:
			}
		}
		return false
	}
	return true
}

// Java: ComposeWarningsGuard.enables(DiagnosticGroup)
func (g *ComposeWarningsGuard) Enables(group *DiagnosticGroup) bool {
	for guard := range g.guards.All() {
		switch guard.MustRunChecks(group) {
		case base.Tri_TRUE:
			return true
		case base.Tri_FALSE:
			return false
		case base.Tri_UNKNOWN:
		}
	}
	return false
}

// GetGuards returns the guards in the order they are applied (a snapshot of
// the unmodifiable sorted set).
//
// Java: ComposeWarningsGuard.getGuards()
func (g *ComposeWarningsGuard) GetGuards() []WarningsGuard {
	return g.guards.ToSlice()
}

// Java: ComposeWarningsGuard.toString()
func (g *ComposeWarningsGuard) String() string {
	var parts []string
	for guard := range g.guards.All() {
		parts = append(parts, jutil.ToJavaString(guard))
	}
	return strings.Join(parts, ", ")
}
