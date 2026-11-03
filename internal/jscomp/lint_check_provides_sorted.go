// Port of com/google/javascript/jscomp/lint/CheckProvidesSorted.java @ closure-compiler v20261028

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

	"github.com/paullepoulpe/jsgo/internal/ast"
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

	firstNode rhino.Node
	lastNode  rhino.Node
	finished  bool

	replacement string
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

var checkProvidesSortedModeNames = [...]string{"COLLECT_ONLY", "COLLECT_AND_REPORT"}

// Name is Enum.name().
func (m CheckProvidesSorted_Mode) Name() string { return checkProvidesSortedModeNames[m] }

// String is Enum.toString().
func (m CheckProvidesSorted_Mode) String() string { return m.Name() }

// NewCheckProvidesSorted constructs a CheckProvidesSorted callback.
//
// Java: CheckProvidesSorted(CheckProvidesSorted.Mode)
func NewCheckProvidesSorted(mode CheckProvidesSorted_Mode) *CheckProvidesSorted {
	return &CheckProvidesSorted{
		firstNode: rhino.NilNode,
		lastNode:  rhino.NilNode,
		mode:      mode,
	}
}

// GetFirstNode returns the node for the first recognized provide statement.
//
// Java: CheckProvidesSorted.getFirstNode()
func (c *CheckProvidesSorted) GetFirstNode() rhino.Node { return c.firstNode }

// GetLastNode returns the node for the last recognized provide statement.
//
// Java: CheckProvidesSorted.getLastNode()
func (c *CheckProvidesSorted) GetLastNode() rhino.Node { return c.lastNode }

// Replacement returns the canonical provides replacement when NeedsFix() is true.
//
// Java: CheckProvidesSorted.getReplacement()
func (c *CheckProvidesSorted) Replacement() (string, bool) {
	if !c.needsFix {
		return "", false
	}
	return c.replacement, true
}

// NeedsFix returns whether the provides need to be fixed, i.e., whether they
// are *not* already canonical.
//
// Java: CheckProvidesSorted.needsFix()
func (c *CheckProvidesSorted) NeedsFix() bool { return c.needsFix }

// Java: CheckProvidesSorted.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckProvidesSorted) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// Traverse top-level statements until a block of contiguous provides is found.
	return !c.finished &&
		(parent.IsNil() || parent.IsRoot() || parent.IsScript() || parent.IsModuleBody())
}

// Java: CheckProvidesSorted.visit(NodeTraversal, Node, Node)
func (c *CheckProvidesSorted) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsScript() {
		c.checkCanonical(t)
		return
	}

	if namespace, ok := validGoogProvideNamespace(n); ok {
		c.originalProvides = append(c.originalProvides, namespace)
		if c.firstNode.IsNil() {
			c.firstNode = n
		}
		c.lastNode = n
	} else if len(c.originalProvides) != 0 {
		c.finished = true
	}
}

func validGoogProvideNamespace(stmt rhino.Node) (string, bool) {
	if !NodeUtil_IsGoogProvideCall(stmt) {
		return "", false
	}
	call := ast.MustCall(ast.MustExprStmt(stmt).Expr().Node())
	onlyArg, ok := call.OnlyArg()
	if !ok {
		return "", false
	}
	ns, ok := ast.AsStringLit(onlyArg.Node())
	if !ok {
		return "", false
	}
	return ns.Value(), true
}

func formatProvide(namespace string) string {
	return "goog.provide('" + namespace + "');"
}

func (c *CheckProvidesSorted) checkCanonical(t *NodeTraversal) {
	sorted := slices.Clone(c.originalProvides)
	slices.SortStableFunc(sorted, jutil.StringCompareTo)
	canonicalProvides := slices.Compact(sorted)
	if !slices.Equal(c.originalProvides, canonicalProvides) {
		c.needsFix = true
		formatted := make([]string, len(canonicalProvides))
		for i, p := range canonicalProvides {
			formatted[i] = formatProvide(p)
		}
		c.replacement = strings.Join(formatted, "\n")
		if c.mode == CheckProvidesSorted_Mode_COLLECT_AND_REPORT {
			t.Report(c.firstNode, CheckProvidesSorted_PROVIDES_NOT_SORTED, c.replacement)
		}
	}
}
