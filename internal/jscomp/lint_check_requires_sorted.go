// Port of com/google/javascript/jscomp/lint/CheckRequiresSorted.java @ closure-compiler v20261028

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

// CheckRequiresSorted checks that Closure import statements (goog.require,
// goog.requireType, and goog.forwardDeclare) are sorted and deduplicated,
// exposing the necessary information to produce a suggested fix.
type CheckRequiresSorted struct {
	mode CheckRequiresSorted_Mode

	importsByNamespace OrderedStringMap[[]requiresSortedImport]
	originalImports    []requiresSortedImport

	firstNode rhino.Node
	lastNode  rhino.Node
	finished  bool

	needsFix       bool
	hasReplacement bool
	replacement    string
}

var CheckRequiresSorted_REQUIRES_NOT_SORTED = DiagnosticType_Warning(
	"JSC_REQUIRES_NOT_SORTED",
	"goog.require() and goog.requireType() statements are not in recommended format."+
		" The correct order is:\n\n{0}\n")

// CheckRequiresSorted_Mode: operation modes.
type CheckRequiresSorted_Mode int32

const (
	// CheckRequiresSorted_Mode_COLLECT_ONLY: collect information to determine
	// whether a fix is required, but do not report a warning.
	CheckRequiresSorted_Mode_COLLECT_ONLY CheckRequiresSorted_Mode = iota
	// CheckRequiresSorted_Mode_COLLECT_AND_REPORT: additionally report a
	// warning.
	CheckRequiresSorted_Mode_COLLECT_AND_REPORT
)

type requiresSortedPrimitive uint8

const (
	requiresSortedRequire requiresSortedPrimitive = iota
	requiresSortedRequireType
	requiresSortedForwardDeclare
)

var requiresSortedPrimitiveNames = [...]string{
	"goog.require",
	"goog.requireType",
	"goog.forwardDeclare",
}

func (p requiresSortedPrimitive) String() string { return requiresSortedPrimitiveNames[p] }

func requiresSortedPrimitiveFromName(name string) requiresSortedPrimitive {
	switch name {
	case "goog.require":
		return requiresSortedRequire
	case "goog.requireType":
		return requiresSortedRequireType
	case "goog.forwardDeclare":
		return requiresSortedForwardDeclare
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("Invalid primitive name " + name))
	}
}

type requiresSortedBinding struct {
	exportedName        string
	localName           string
	isShorthandProperty bool
}

func (b requiresSortedBinding) canonicalizeShorthand() requiresSortedBinding {
	if b.exportedName == b.localName {
		b.isShorthandProperty = true
	}
	return b
}

type requiresSortedImport struct {
	nodes               []rhino.Node
	destructures        []requiresSortedBinding
	namespace           string
	alias               string
	formattedWithoutDoc string
	primitive           requiresSortedPrimitive
	hasAlias            bool
	isDestructuring     bool
}

func (s requiresSortedImport) isStandalone() bool {
	return !s.hasAlias && !s.isDestructuring
}

func (s requiresSortedImport) upgrade(other requiresSortedPrimitive) requiresSortedImport {
	if other < s.primitive {
		s.primitive = other
		s.formattedWithoutDoc = ""
	}
	return s
}

func (s *requiresSortedImport) formatWithoutDoc() string {
	if s.formattedWithoutDoc != "" {
		return s.formattedWithoutDoc
	}
	var sb strings.Builder
	if !s.isStandalone() {
		sb.WriteString("const ")
	}
	if s.hasAlias {
		sb.WriteString(s.alias)
	}
	if s.isDestructuring {
		sb.WriteByte('{')
		for i, binding := range s.destructures {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(binding.exportedName)
			if binding.exportedName != binding.localName {
				sb.WriteString(": ")
				sb.WriteString(binding.localName)
			}
		}
		sb.WriteByte('}')
	}
	if !s.isStandalone() {
		sb.WriteString(" = ")
	}
	sb.WriteString(s.primitive.String())
	sb.WriteString("('")
	sb.WriteString(s.namespace)
	sb.WriteString("');")
	s.formattedWithoutDoc = sb.String()
	return s.formattedWithoutDoc
}

func (s *requiresSortedImport) format() string {
	var sb strings.Builder
	for _, node := range s.nodes {
		if comment := node.GetNonJSDocCommentString(); comment != "" {
			sb.WriteString(comment)
			sb.WriteByte('\n')
		}
		if jsDoc := NodeUtil_GetBestJSDocInfo(node); jsDoc != nil {
			// StringBuilder.append(String) appends "null" for a null string.
			if original := jsDoc.GetOriginalCommentString(); original != nil {
				sb.WriteString(*original)
			} else {
				sb.WriteString("null")
			}
			sb.WriteByte('\n')
		}
	}
	sb.WriteString(s.formatWithoutDoc())
	return sb.String()
}

func equalRequiresSortedImport(a, b requiresSortedImport) bool {
	return a.primitive == b.primitive &&
		a.namespace == b.namespace &&
		a.hasAlias == b.hasAlias &&
		a.alias == b.alias &&
		a.isDestructuring == b.isDestructuring &&
		slices.Equal(a.nodes, b.nodes) &&
		slices.Equal(a.destructures, b.destructures)
}

var (
	checkRequiresSorted_GOOG_REQUIRE        = rhino.QualifiedName_Of("goog.require")
	checkRequiresSorted_GOOG_REQUIRETYPE    = rhino.QualifiedName_Of("goog.requireType")
	checkRequiresSorted_GOOG_FORWARDDECLARE = rhino.QualifiedName_Of("goog.forwardDeclare")
)

// NewCheckRequiresSorted is new CheckRequiresSorted(mode).
//
// Java: CheckRequiresSorted(CheckRequiresSorted.Mode)
func NewCheckRequiresSorted(mode CheckRequiresSorted_Mode) *CheckRequiresSorted {
	return &CheckRequiresSorted{
		firstNode: rhino.NilNode,
		lastNode:  rhino.NilNode,
		mode:      mode,
	}
}

// GetFirstNode returns the node for the first recognized import statement.
//
// Java: CheckRequiresSorted.getFirstNode()
func (c *CheckRequiresSorted) GetFirstNode() rhino.Node { return c.firstNode }

// GetLastNode returns the node for the last recognized import statement.
//
// Java: CheckRequiresSorted.getLastNode()
func (c *CheckRequiresSorted) GetLastNode() rhino.Node { return c.lastNode }

// GetReplacement returns a textual replacement yielding a canonical version
// of the imports.
//
// Java: CheckRequiresSorted.getReplacement()
func (c *CheckRequiresSorted) GetReplacement() *string {
	if !c.hasReplacement {
		return nil
	}
	s := c.replacement
	return &s
}

// NeedsFix returns whether the imports need to be fixed, i.e., whether they
// are *not* already canonical.
//
// Java: CheckRequiresSorted.needsFix()
func (c *CheckRequiresSorted) NeedsFix() bool { return c.needsFix }

// Java: CheckRequiresSorted.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckRequiresSorted) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// Traverse top-level statements until a block of contiguous requires is found.
	return !c.finished &&
		(parent.IsNil() || parent.IsRoot() || parent.IsScript() || parent.IsModuleBody())
}

// Java: CheckRequiresSorted.visit(NodeTraversal, Node, Node)
func (c *CheckRequiresSorted) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if _, ok := ast.AsScript(n); ok {
		c.checkCanonical(t)
		return
	}

	var callNode rhino.Node
	if exprStmt, ok := ast.AsExprStmt(n); ok {
		callNode = exprStmt.Expr().Node()
	} else if decl, ok := ast.AsNameDeclaration(n); ok {
		if first, ok := decl.FirstDeclarator(); ok {
			callNode = first.Node().GetLastChild()
		}
	}
	if !callNode.IsNil() && isValidRequireImportCall(callNode) {
		stmt := parseRequireImport(callNode)
		c.originalImports = append(c.originalImports, stmt)
		c.importsByNamespace.Put(stmt.namespace, append(c.importsByNamespace.Get(stmt.namespace), stmt))
		if c.firstNode.IsNil() {
			c.firstNode = n
		}
		c.lastNode = n
	} else if !c.importsByNamespace.IsEmpty() {
		c.finished = true
	}
}

func isValidRequireImportCall(n rhino.Node) bool {
	call, ok := ast.AsCall(n)
	if !ok || !n.HasTwoChildren() {
		return false
	}
	callee := call.Callee().Node()
	firstArg, ok := call.FirstArg()
	return ok &&
		firstArg.Token() == rhino.Token_STRINGLIT &&
		(checkRequiresSorted_GOOG_REQUIRE.Matches(callee) ||
			checkRequiresSorted_GOOG_REQUIRETYPE.Matches(callee) ||
			checkRequiresSorted_GOOG_FORWARDDECLARE.Matches(callee))
}

func parseRequireImport(callNode rhino.Node) requiresSortedImport {
	call := ast.MustCall(callNode)
	primitive := requiresSortedPrimitiveFromName(*call.Callee().Node().GetQualifiedName())
	namespace := ast.MustStringLit(call.MustFirstArg().Node()).Value()

	parent := callNode.GetParent()
	if parent.IsExprResult() {
		return requiresSortedImport{
			nodes:     []rhino.Node{parent},
			primitive: primitive,
			namespace: namespace,
		}
	}

	grandparent := parent.GetParent()
	if name, ok := ast.AsName(parent); ok {
		return requiresSortedImport{
			nodes:     []rhino.Node{grandparent},
			primitive: primitive,
			namespace: namespace,
			alias:     name.Ident(),
			hasAlias:  true,
		}
	}

	var destructures []requiresSortedBinding
	if lhs, ok := ast.AsDestructuringLhs(parent); ok {
		if objPattern, ok := lhs.Pattern().AsObjectPattern(); ok {
			for prop := range objPattern.Props() {
				if key, ok := prop.AsStringKey(); ok {
					if local, ok := key.Value().AsName(); ok {
						destructures = append(destructures, requiresSortedBinding{
							exportedName:        key.Key(),
							localName:           local.Ident(),
							isShorthandProperty: key.Node().IsShorthandProperty(),
						})
						continue
					}
				}
				propNode := prop.Node()
				destructures = append(destructures, requiresSortedBinding{
					exportedName:        propNode.GetString(),
					localName:           propNode.GetFirstChild().GetString(),
					isShorthandProperty: propNode.IsShorthandProperty(),
				})
			}
		} else {
			for child := range lhs.Pattern().Node().Children() {
				destructures = append(destructures, requiresSortedBinding{
					exportedName:        child.GetString(),
					localName:           child.GetFirstChild().GetString(),
					isShorthandProperty: child.IsShorthandProperty(),
				})
			}
		}
	}
	return requiresSortedImport{
		nodes:           []rhino.Node{grandparent},
		primitive:       primitive,
		namespace:       namespace,
		destructures:    destructures,
		isDestructuring: true,
	}
}

func (c *CheckRequiresSorted) checkCanonical(t *NodeTraversal) {
	canonicalImports := canonicalizeRequireImports(&c.importsByNamespace)
	if !slices.EqualFunc(c.originalImports, canonicalImports, equalRequiresSortedImport) {
		c.needsFix = true
		var sb strings.Builder
		for i := range canonicalImports {
			if i > 0 {
				sb.WriteByte('\n')
			}
			sb.WriteString(canonicalImports[i].format())
		}
		c.replacement = sb.String()
		c.hasReplacement = true
		if c.mode == CheckRequiresSorted_Mode_COLLECT_AND_REPORT {
			t.Report(c.firstNode, CheckRequiresSorted_REQUIRES_NOT_SORTED, c.replacement)
		}
	}
}

func canonicalizeRequireImports(importsByNamespace *OrderedStringMap[[]requiresSortedImport]) []requiresSortedImport {
	entries := importsByNamespace.Entries()
	if len(entries) == 0 {
		return nil
	}
	canonicalImports := make([]requiresSortedImport, 0, len(entries))

	for _, entry := range entries {
		namespace, allImports := entry.Key, entry.Value

		strongest := requiresSortedForwardDeclare
		for i := range allImports {
			strongest = min(strongest, allImports[i].primitive)
		}

		hasAliasing := false
		for i := range allImports {
			if allImports[i].hasAlias {
				canonicalImports = append(canonicalImports, allImports[i].upgrade(strongest))
				hasAliasing = true
			}
		}

		hasDestructuring := false
		var destructuringNodes []rhino.Node
		var allDestructures []requiresSortedBinding
		for i := range allImports {
			if allImports[i].isDestructuring {
				destructuringNodes = append(destructuringNodes, allImports[i].nodes...)
				for _, d := range allImports[i].destructures {
					allDestructures = append(allDestructures, d.canonicalizeShorthand())
				}
			}
		}
		if len(allDestructures) > 0 {
			slices.SortFunc(allDestructures, func(a, b requiresSortedBinding) int {
				if c := jutil.StringCompareTo(a.exportedName, b.exportedName); c != 0 {
					return c
				}
				return jutil.StringCompareTo(a.localName, b.localName)
			})
			canonicalImports = append(canonicalImports, requiresSortedImport{
				nodes:           destructuringNodes,
				primitive:       strongest,
				namespace:       namespace,
				destructures:    slices.Compact(allDestructures),
				isDestructuring: true,
			})
			hasDestructuring = true
		}

		if !hasAliasing && !hasDestructuring {
			var standaloneNodes []rhino.Node
			for i := range allImports {
				if allImports[i].isStandalone() {
					standaloneNodes = append(standaloneNodes, allImports[i].nodes...)
				}
			}
			canonicalImports = append(canonicalImports, requiresSortedImport{
				nodes:     standaloneNodes,
				primitive: strongest,
				namespace: namespace,
			})
		}
	}

	for i := range canonicalImports {
		canonicalImports[i].formatWithoutDoc()
	}
	slices.SortStableFunc(canonicalImports, func(a, b requiresSortedImport) int {
		return jutil.StringCompareTo(a.formattedWithoutDoc, b.formattedWithoutDoc)
	})
	return canonicalImports
}
