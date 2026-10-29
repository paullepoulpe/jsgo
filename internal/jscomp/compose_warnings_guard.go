// Port of com/google/javascript/jscomp/ComposeWarningsGuard.java @ closure-compiler v20261021

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
	"cmp"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type warningsGuardSequenced interface {
	ensureGuardSeq() uint64
}

func sameWarningsGuard(a, b WarningsGuard) bool {
	if sa, ok := a.(warningsGuardSequenced); ok {
		if sb, ok := b.(warningsGuardSequenced); ok {
			return sa.ensureGuardSeq() == sb.ensureGuardSeq()
		}
	}
	return a == b
}

type composeGuardEntry struct {
	guard WarningsGuard
	seq   uint64
}

func composeGuardEntryCompare(a, b composeGuardEntry) int {
	priorityDiff := a.guard.GetPriority() - b.guard.GetPriority()
	if priorityDiff != 0 {
		return priorityDiff
	}
	// If the warnings guards have the same priority, the one that
	// was added last wins.
	return cmp.Compare(b.seq, a.seq)
}

// ComposeWarningsGuard is a WarningsGuard that represents just a chain of
// other guards. For example we could have following chain 1) all warnings
// outside of /foo/ should be suppressed 2) errors with key JSC_BAR should be
// marked as warning 3) the rest should be reported as error
//
// This class is designed for such behavior.
type ComposeWarningsGuard struct {
	WarningsGuardBase

	// Guards ordered in application priority (priority asc, addition order desc).
	entries      []composeGuardEntry
	nextSeq      uint64
	demoteErrors bool
}

var _ WarningsGuard = (*ComposeWarningsGuard)(nil)

const composeWarningsGuard_serialVersionUID int64 = 1

// Java: ComposeWarningsGuard(List<WarningsGuard>)
func NewComposeWarningsGuardFromList(guards []WarningsGuard) *ComposeWarningsGuard {
	g := &ComposeWarningsGuard{}
	g.AddGuards(guards)
	return g
}

// Java: ComposeWarningsGuard(WarningsGuard...)
func NewComposeWarningsGuard(guards ...WarningsGuard) *ComposeWarningsGuard {
	for _, guard := range guards {
		jutil.CheckNotNull(guard) // ImmutableList.copyOf
	}
	g := &ComposeWarningsGuard{}
	g.AddGuards(guards)
	return g
}

// Java: ComposeWarningsGuard.addGuard(WarningsGuard)
func (g *ComposeWarningsGuard) AddGuard(guard WarningsGuard) {
	if composeGuard, ok := guard.(*ComposeWarningsGuard); ok {
		if composeGuard.demoteErrors {
			g.demoteErrors = true
		}
		// Snapshot in reverse so self-composition (g.AddGuard(g)) preserves order.
		reversed := make([]WarningsGuard, len(composeGuard.entries))
		for i, entry := range composeGuard.entries {
			reversed[len(composeGuard.entries)-1-i] = entry.guard
		}
		g.AddGuards(reversed)
		return
	}

	if s, ok := guard.(warningsGuardSequenced); ok {
		s.ensureGuardSeq()
	}
	for i, existing := range g.entries {
		if sameWarningsGuard(existing.guard, guard) {
			g.entries = slices.Delete(g.entries, i, i+1)
			break
		}
	}
	g.nextSeq++
	entry := composeGuardEntry{guard: guard, seq: g.nextSeq}
	pos := len(g.entries)
	for i, existing := range g.entries {
		if composeGuardEntryCompare(entry, existing) < 0 {
			pos = i
			break
		}
	}
	g.entries = slices.Insert(g.entries, pos, entry)
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
	for _, entry := range g.entries {
		newLevel := entry.guard.Level(error)
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
	for _, type_ := range group.GetTypes() {
		singleton := DiagnosticGroup_ForType(type_)
		for _, entry := range g.entries {
			switch entry.guard.MustRunChecks(singleton) {
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
	for _, entry := range g.entries {
		switch entry.guard.MustRunChecks(group) {
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
	res := make([]WarningsGuard, len(g.entries))
	for i, entry := range g.entries {
		res[i] = entry.guard
	}
	return res
}

// Java: ComposeWarningsGuard.toString()
func (g *ComposeWarningsGuard) String() string {
	if len(g.entries) == 0 {
		return ""
	}
	parts := make([]string, len(g.entries))
	for i, entry := range g.entries {
		parts[i] = jutil.ToJavaString(entry.guard)
	}
	return strings.Join(parts, ", ")
}
