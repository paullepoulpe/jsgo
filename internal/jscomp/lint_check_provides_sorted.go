// Port of com/google/javascript/jscomp/lint/CheckProvidesSorted.java @ closure-compiler v20261021

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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckProvidesSorted checks that goog.provide statements are sorted and
// deduplicated, exposing the necessary information to produce a suggested
// fix.
type CheckProvidesSorted struct {
	mode CheckProvidesSorted_Mode

	// The provided namespaces in the order they appear.
	originalProvides []string

	// The provided namespaces in canonical order.
	firstNode rhino.Node
	lastNode  rhino.Node
	finished  bool

	replacement *string
	needsFix    bool
}

var CheckProvidesSorted_PROVIDES_NOT_SORTED = DiagnosticType_Warning(
	"JSC_PROVIDES_NOT_SORTED",
	"goog.provide() statements are not sorted."+
		" The correct order is:\n\n{0}\n")

// CheckProvidesSorted_Mode: operation modes.
type CheckProvidesSorted_Mode int32

const (
	// CheckProvidesSorted_Mode_COLLECT_ONLY: collect information to
	// determine whether a fix is required, but do not report a warning.
	CheckProvidesSorted_Mode_COLLECT_ONLY CheckProvidesSorted_Mode = iota
	// CheckProvidesSorted_Mode_COLLECT_AND_REPORT: additionally report a
	// warning.
	CheckProvidesSorted_Mode_COLLECT_AND_REPORT
)

var checkProvidesSorted_Mode_names = [...]string{"COLLECT_ONLY", "COLLECT_AND_REPORT"}

// CheckProvidesSorted_Mode_values is Mode.values().
func CheckProvidesSorted_Mode_values() []CheckProvidesSorted_Mode {
	return []CheckProvidesSorted_Mode{CheckProvidesSorted_Mode_COLLECT_ONLY, CheckProvidesSorted_Mode_COLLECT_AND_REPORT}
}

// Name is Enum.name().
func (m CheckProvidesSorted_Mode) Name() string { return checkProvidesSorted_Mode_names[m] }

// String is Enum.toString().
func (m CheckProvidesSorted_Mode) String() string { return m.Name() }

// Ordinal is Enum.ordinal().
func (m CheckProvidesSorted_Mode) Ordinal() int { return int(m) }

// NewCheckProvidesSorted is new CheckProvidesSorted(mode).
//
// Java: CheckProvidesSorted(CheckProvidesSorted.Mode)
func NewCheckProvidesSorted(mode CheckProvidesSorted_Mode) *CheckProvidesSorted {
	return &CheckProvidesSorted{
		firstNode:   rhino.NilNode,
		lastNode:    rhino.NilNode,
		finished:    false,
		replacement: nil,
		needsFix:    false,
		mode:        mode,
	}
}

// GetFirstNode returns the node for the first recognized provide statement.
//
// Java: CheckProvidesSorted.getFirstNode()
func (c *CheckProvidesSorted) GetFirstNode() rhino.Node {
	return c.firstNode
}

// GetLastNode returns the node for the last recognized provide statement.
//
// Java: CheckProvidesSorted.getLastNode()
func (c *CheckProvidesSorted) GetLastNode() rhino.Node {
	return c.lastNode
}

// GetReplacement returns a textual replacement yielding a canonical version
// of the provides.
//
// Java: CheckProvidesSorted.getReplacement()
func (c *CheckProvidesSorted) GetReplacement() *string {
	return c.replacement
}

// NeedsFix returns whether the provides need to be fixed, i.e., whether they
// are *not* already canonical.
//
// Java: CheckProvidesSorted.needsFix()
func (c *CheckProvidesSorted) NeedsFix() bool {
	return c.needsFix
}

// Java: CheckProvidesSorted.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckProvidesSorted) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// Traverse top-level statements until a block of contiguous provides is found.
	return !c.finished &&
		(parent.IsNil() || parent.IsRoot() || parent.IsScript() || parent.IsModuleBody())
}

// Java: CheckProvidesSorted.visit(NodeTraversal, Node, Node)
func (c *CheckProvidesSorted) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsScript() {
		c.CheckCanonical(t)
		return
	}

	if NodeUtil_IsGoogProvideCall(n) && CheckProvidesSorted_AreProvideArgumentsValid(n.GetFirstChild()) {
		c.originalProvides = append(c.originalProvides, CheckProvidesSorted_GetNamespace(n))
		if c.firstNode.IsNil() {
			c.lastNode = n
			c.firstNode = c.lastNode
		} else {
			c.lastNode = n
		}
	} else if len(c.originalProvides) != 0 {
		c.finished = true
	}
}

// Java: CheckProvidesSorted.areProvideArgumentsValid(Node)
func CheckProvidesSorted_AreProvideArgumentsValid(call rhino.Node) bool {
	return call.HasTwoChildren() && call.GetSecondChild().IsStringLit()
}

// Java: CheckProvidesSorted.getNamespace(Node)
func CheckProvidesSorted_GetNamespace(n rhino.Node) string {
	return n.GetFirstChild().GetSecondChild().GetString()
}

// CheckProvidesSorted_FormatProvide returns the code for a correctly
// formatted provide call.
//
// Java: CheckProvidesSorted.formatProvide(String)
func CheckProvidesSorted_FormatProvide(namespace string) string {
	var sb strings.Builder

	sb.WriteString("goog.provide('")
	sb.WriteString(namespace)
	sb.WriteString("');")

	return sb.String()
}

// Java: CheckProvidesSorted.checkCanonical(NodeTraversal)
func (c *CheckProvidesSorted) CheckCanonical(t *NodeTraversal) {
	// originalProvides.stream().sorted().distinct().collect(toImmutableList())
	sorted := slices.Clone(c.originalProvides)
	slices.SortStableFunc(sorted, jutil.StringCompareTo)
	canonicalProvides := slices.Compact(sorted)
	if !slices.Equal(c.originalProvides, canonicalProvides) {
		c.needsFix = true
		formatted := make([]string, len(canonicalProvides))
		for i, p := range canonicalProvides {
			formatted[i] = CheckProvidesSorted_FormatProvide(p)
		}
		c.replacement = jutil.StringPtr(strings.Join(formatted, "\n"))
		if c.mode == CheckProvidesSorted_Mode_COLLECT_AND_REPORT {
			t.Report(c.firstNode, CheckProvidesSorted_PROVIDES_NOT_SORTED, *c.replacement)
		}
	}
}
