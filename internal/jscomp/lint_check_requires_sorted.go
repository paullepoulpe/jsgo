// Port of com/google/javascript/jscomp/lint/CheckRequiresSorted.java @ closure-compiler v20261021

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

	// Maps each namespace into the existing import statements for that namespace,
	// accompanied by namespaceOrder for deterministic iteration order.
	importsByNamespace map[string][]*CheckRequiresSorted_ImportStatement
	namespaceOrder     []string

	// The import statements in the order they appear.
	originalImports []*CheckRequiresSorted_ImportStatement

	// The import statements in canonical order.
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

var checkRequiresSorted_Mode_names = [...]string{"COLLECT_ONLY", "COLLECT_AND_REPORT"}

// Name is Enum.name().
func (m CheckRequiresSorted_Mode) Name() string { return checkRequiresSorted_Mode_names[m] }

// String is Enum.toString().
func (m CheckRequiresSorted_Mode) String() string { return m.Name() }

// Ordinal is Enum.ordinal().
func (m CheckRequiresSorted_Mode) Ordinal() int { return int(m) }

// CheckRequiresSorted_ImportPrimitive: primitives that may be called in an
// import statement.
type CheckRequiresSorted_ImportPrimitive int32

const (
	CheckRequiresSorted_ImportPrimitive_REQUIRE CheckRequiresSorted_ImportPrimitive = iota
	CheckRequiresSorted_ImportPrimitive_REQUIRE_TYPE
	CheckRequiresSorted_ImportPrimitive_FORWARD_DECLARE
)

var checkRequiresSorted_ImportPrimitive_enumNames = [...]string{"REQUIRE", "REQUIRE_TYPE", "FORWARD_DECLARE"}

// The name field of each constant, indexed by ordinal.
var checkRequiresSorted_ImportPrimitive_name = [...]string{
	"goog.require",
	"goog.requireType",
	"goog.forwardDeclare",
}

// CheckRequiresSorted_ImportPrimitive_values is ImportPrimitive.values().
func CheckRequiresSorted_ImportPrimitive_values() []CheckRequiresSorted_ImportPrimitive {
	return []CheckRequiresSorted_ImportPrimitive{
		CheckRequiresSorted_ImportPrimitive_REQUIRE,
		CheckRequiresSorted_ImportPrimitive_REQUIRE_TYPE,
		CheckRequiresSorted_ImportPrimitive_FORWARD_DECLARE,
	}
}

// Name is Enum.name().
func (p CheckRequiresSorted_ImportPrimitive) Name() string {
	return checkRequiresSorted_ImportPrimitive_enumNames[p]
}

// Ordinal is Enum.ordinal().
func (p CheckRequiresSorted_ImportPrimitive) Ordinal() int { return int(p) }

// CheckRequiresSorted_ImportPrimitive_FromName returns the primitive with the
// given name.
//
// Java: CheckRequiresSorted.ImportPrimitive.fromName(String)
func CheckRequiresSorted_ImportPrimitive_FromName(name string) CheckRequiresSorted_ImportPrimitive {
	switch name {
	case "goog.require":
		return CheckRequiresSorted_ImportPrimitive_REQUIRE
	case "goog.requireType":
		return CheckRequiresSorted_ImportPrimitive_REQUIRE_TYPE
	case "goog.forwardDeclare":
		return CheckRequiresSorted_ImportPrimitive_FORWARD_DECLARE
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("Invalid primitive name " + name))
	}
}

var CheckRequiresSorted_ImportPrimitive_WEAKEST = CheckRequiresSorted_ImportPrimitive_FORWARD_DECLARE

// CheckRequiresSorted_ImportPrimitive_Stronger returns the stronger of two
// primitives.
//
// `goog.require` is stronger than `goog.requireType`, which is stronger than
// `goog.forwardDeclare`.
//
// Java: CheckRequiresSorted.ImportPrimitive.stronger(ImportPrimitive, ImportPrimitive)
func CheckRequiresSorted_ImportPrimitive_Stronger(p1 CheckRequiresSorted_ImportPrimitive, p2 CheckRequiresSorted_ImportPrimitive) CheckRequiresSorted_ImportPrimitive {
	if p1 < p2 {
		return p1
	}
	return p2
}

// String is the overridden toString(): the primitive's name.
//
// Java: CheckRequiresSorted.ImportPrimitive.toString()
func (p CheckRequiresSorted_ImportPrimitive) String() string {
	return checkRequiresSorted_ImportPrimitive_name[p]
}

// CheckRequiresSorted_DestructuringBinding is one of the bindings of a
// destructuring pattern.
//
// exportedName and localName are equal in the case where the binding does
// not explicitly specify a local name.
type CheckRequiresSorted_DestructuringBinding struct {
	exportedName        string
	localName           string
	isShorthandProperty bool
}

// newCheckRequiresSorted_DestructuringBinding is the record's canonical
// constructor.
//
// Java: CheckRequiresSorted.DestructuringBinding(String, String, boolean)
func newCheckRequiresSorted_DestructuringBinding(
	exportedName string, localName string, isShorthandProperty bool) *CheckRequiresSorted_DestructuringBinding {
	jutil.CheckArgument(!isShorthandProperty || exportedName == localName)
	return &CheckRequiresSorted_DestructuringBinding{
		exportedName:        exportedName,
		localName:           localName,
		isShorthandProperty: isShorthandProperty,
	}
}

// Java: CheckRequiresSorted.DestructuringBinding.of(String, String, boolean)
func CheckRequiresSorted_DestructuringBinding_Of(
	exportedName string, localName string, isShorthandProperty bool) *CheckRequiresSorted_DestructuringBinding {
	return newCheckRequiresSorted_DestructuringBinding(exportedName, localName, isShorthandProperty)
}

// ExportedName is the record accessor exportedName().
func (b *CheckRequiresSorted_DestructuringBinding) ExportedName() string { return b.exportedName }

// LocalName is the record accessor localName().
func (b *CheckRequiresSorted_DestructuringBinding) LocalName() string { return b.localName }

// IsShorthandProperty is the record accessor isShorthandProperty().
func (b *CheckRequiresSorted_DestructuringBinding) IsShorthandProperty() bool {
	return b.isShorthandProperty
}

// CompareTo compares two bindings according to the style guide sort order.
//
// Java: CheckRequiresSorted.DestructuringBinding.compareTo(DestructuringBinding)
func (b *CheckRequiresSorted_DestructuringBinding) CompareTo(other *CheckRequiresSorted_DestructuringBinding) int {
	if c := jutil.StringCompareTo(b.exportedName, other.exportedName); c != 0 {
		if c < 0 {
			return -1
		}
		return 1
	}
	if c := jutil.StringCompareTo(b.localName, other.localName); c != 0 {
		if c < 0 {
			return -1
		}
		return 1
	}
	return 0
}

// IsCanonical returns true if the destructuring binding is canonical.
//
// For example:
//   - `{Foo}` is canonical
//   - `{Foo: Bar}` is canonical
//   - `{Foo: Foo}` is not canonical
//
// Java: CheckRequiresSorted.DestructuringBinding.isCanonical()
func (b *CheckRequiresSorted_DestructuringBinding) IsCanonical() bool {
	return b.exportedName != b.localName || b.isShorthandProperty
}

// CanonicalizeShorthandProperties canonicalizes the destructuring to a
// shorthand property when applicable.
//
// In practice, `{Foo: Foo}` gets simplified to `{Foo}`.
//
// Java: CheckRequiresSorted.DestructuringBinding.canonicalizeShorthandProperties()
func (b *CheckRequiresSorted_DestructuringBinding) CanonicalizeShorthandProperties() *CheckRequiresSorted_DestructuringBinding {
	if b.IsCanonical() {
		return b
	}
	return CheckRequiresSorted_DestructuringBinding_Of(
		b.exportedName, b.localName /* isShorthandProperty= */, true)
}

// Equals is the record's component-wise equals.
func (b *CheckRequiresSorted_DestructuringBinding) Equals(other *CheckRequiresSorted_DestructuringBinding) bool {
	if other == nil {
		return false
	}
	return b.exportedName == other.exportedName &&
		b.localName == other.localName &&
		b.isShorthandProperty == other.isShorthandProperty
}

// CheckRequiresSorted_ImportStatement is an import statement, which may have
// been merged from several import statements for the same namespace in the
// original code.
type CheckRequiresSorted_ImportStatement struct {
	nodes               []rhino.Node
	namespace           string
	alias               string
	formattedWithoutDoc string
	destructures        []*CheckRequiresSorted_DestructuringBinding
	primitive           CheckRequiresSorted_ImportPrimitive
	hasAlias            bool
}

// newCheckRequiresSorted_ImportStatement is the record's canonical
// constructor.
//
// Java: CheckRequiresSorted.ImportStatement(ImmutableList<Node>, ImportPrimitive, String, String, ImmutableList<DestructuringBinding>)
func newCheckRequiresSorted_ImportStatement(
	nodes []rhino.Node,
	primitive CheckRequiresSorted_ImportPrimitive,
	namespace string,
	alias string,
	hasAlias bool,
	destructures []*CheckRequiresSorted_DestructuringBinding) *CheckRequiresSorted_ImportStatement {
	if nodes == nil {
		panic(jutil.NewNullPointerExceptionMsg("nodes"))
	}
	jutil.CheckArgumentMsg(
		!hasAlias || destructures == nil,
		"Import statement cannot be simultaneously aliasing and destructuring")
	return &CheckRequiresSorted_ImportStatement{
		nodes:        nodes,
		primitive:    primitive,
		namespace:    namespace,
		alias:        alias,
		hasAlias:     hasAlias,
		destructures: destructures,
	}
}

// CheckRequiresSorted_ImportStatement_Of creates a new import statement.
//
// Java: CheckRequiresSorted.ImportStatement.of(ImmutableList<Node>, ImportPrimitive, String, String, ImmutableList<DestructuringBinding>)
func CheckRequiresSorted_ImportStatement_Of(
	nodes []rhino.Node,
	primitive CheckRequiresSorted_ImportPrimitive,
	namespace string,
	alias string,
	hasAlias bool,
	destructures []*CheckRequiresSorted_DestructuringBinding) *CheckRequiresSorted_ImportStatement {
	return newCheckRequiresSorted_ImportStatement(nodes, primitive, namespace, alias, hasAlias, destructures)
}

// Nodes is the record accessor nodes().
func (s *CheckRequiresSorted_ImportStatement) Nodes() []rhino.Node { return s.nodes }

// Primitive is the record accessor primitive().
func (s *CheckRequiresSorted_ImportStatement) Primitive() CheckRequiresSorted_ImportPrimitive {
	return s.primitive
}

// Namespace is the record accessor namespace().
func (s *CheckRequiresSorted_ImportStatement) Namespace() string { return s.namespace }

// Alias is the record accessor alias().
func (s *CheckRequiresSorted_ImportStatement) Alias() string { return s.alias }

// Destructures is the record accessor destructures().
func (s *CheckRequiresSorted_ImportStatement) Destructures() []*CheckRequiresSorted_DestructuringBinding {
	return s.destructures
}

// IsStandalone returns whether the import is standalone.
//
// Java: CheckRequiresSorted.ImportStatement.isStandalone()
func (s *CheckRequiresSorted_ImportStatement) IsStandalone() bool {
	return !s.hasAlias && s.destructures == nil
}

// IsAliasing returns whether the import is aliasing.
//
// Java: CheckRequiresSorted.ImportStatement.isAliasing()
func (s *CheckRequiresSorted_ImportStatement) IsAliasing() bool {
	return s.hasAlias
}

// IsDestructuring returns whether the import is destructuring.
//
// Java: CheckRequiresSorted.ImportStatement.isDestructuring()
func (s *CheckRequiresSorted_ImportStatement) IsDestructuring() bool {
	return s.destructures != nil
}

// Upgrade returns an import statement identical to the current one, except
// for its primitive, which is upgraded to the given one if stronger.
//
// Java: CheckRequiresSorted.ImportStatement.upgrade(ImportPrimitive)
func (s *CheckRequiresSorted_ImportStatement) Upgrade(otherPrimitive CheckRequiresSorted_ImportPrimitive) *CheckRequiresSorted_ImportStatement {
	if CheckRequiresSorted_ImportPrimitive_Stronger(s.primitive, otherPrimitive) != s.primitive {
		return newCheckRequiresSorted_ImportStatement(s.nodes, otherPrimitive, s.namespace, s.alias, s.hasAlias, s.destructures)
	}
	return s
}

// Java: CheckRequiresSorted.ImportStatement.formatWithoutDoc()
func (s *CheckRequiresSorted_ImportStatement) FormatWithoutDoc() string {
	if s.formattedWithoutDoc != "" {
		return s.formattedWithoutDoc
	}
	var sb strings.Builder

	if !s.IsStandalone() {
		sb.WriteString("const ")
	}
	if s.hasAlias {
		sb.WriteString(s.alias)
	}
	if s.destructures != nil {
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
	if !s.IsStandalone() {
		sb.WriteString(" = ")
	}
	sb.WriteString(s.primitive.String())
	sb.WriteString("('")
	sb.WriteString(s.namespace)
	sb.WriteString("');")

	s.formattedWithoutDoc = sb.String()
	return s.formattedWithoutDoc
}

// Format formats the import statement into code.
//
// Java: CheckRequiresSorted.ImportStatement.format()
func (s *CheckRequiresSorted_ImportStatement) Format() string {
	var sb strings.Builder

	for _, node := range s.nodes {
		comment := node.GetNonJSDocCommentString()
		if comment != "" {
			sb.WriteString(comment)
			sb.WriteByte('\n')
		}
		jsDoc := NodeUtil_GetBestJSDocInfo(node)
		if jsDoc != nil {
			// StringBuilder.append(String) appends "null" for a null string.
			if original := jsDoc.GetOriginalCommentString(); original != nil {
				sb.WriteString(*original)
			} else {
				sb.WriteString("null")
			}
			sb.WriteByte('\n')
		}
	}

	sb.WriteString(s.FormatWithoutDoc())

	return sb.String()
}

// CompareTo compares two import statements according to the style guide sort
// order.
//
// Java: CheckRequiresSorted.ImportStatement.compareTo(ImportStatement)
func (s *CheckRequiresSorted_ImportStatement) CompareTo(other *CheckRequiresSorted_ImportStatement) int {
	return jutil.StringCompareTo(s.FormatWithoutDoc(), other.FormatWithoutDoc())
}

// Equals is the record's component-wise equals.
func (s *CheckRequiresSorted_ImportStatement) Equals(other *CheckRequiresSorted_ImportStatement) bool {
	if other == nil {
		return false
	}
	if !slices.Equal(s.nodes, other.nodes) ||
		s.primitive != other.primitive ||
		s.namespace != other.namespace ||
		s.hasAlias != other.hasAlias ||
		s.alias != other.alias ||
		(s.destructures == nil) != (other.destructures == nil) {
		return false
	}
	return slices.EqualFunc(s.destructures, other.destructures,
		(*CheckRequiresSorted_DestructuringBinding).Equals)
}

var checkRequiresSorted_GOOG_REQUIRE = rhino.QualifiedName_Of("goog.require")

var checkRequiresSorted_GOOG_REQUIRETYPE = rhino.QualifiedName_Of("goog.requireType")

var checkRequiresSorted_GOOG_FORWARDDECLARE = rhino.QualifiedName_Of("goog.forwardDeclare")

// NewCheckRequiresSorted is new CheckRequiresSorted(mode).
//
// Java: CheckRequiresSorted(CheckRequiresSorted.Mode)
func NewCheckRequiresSorted(mode CheckRequiresSorted_Mode) *CheckRequiresSorted {
	return &CheckRequiresSorted{
		importsByNamespace: make(map[string][]*CheckRequiresSorted_ImportStatement),
		firstNode:          rhino.NilNode,
		lastNode:           rhino.NilNode,
		mode:               mode,
	}
}

// GetFirstNode returns the node for the first recognized import statement.
//
// Java: CheckRequiresSorted.getFirstNode()
func (c *CheckRequiresSorted) GetFirstNode() rhino.Node {
	return c.firstNode
}

// GetLastNode returns the node for the last recognized import statement.
//
// Java: CheckRequiresSorted.getLastNode()
func (c *CheckRequiresSorted) GetLastNode() rhino.Node {
	return c.lastNode
}

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
func (c *CheckRequiresSorted) NeedsFix() bool {
	return c.needsFix
}

// Java: CheckRequiresSorted.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckRequiresSorted) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// Traverse top-level statements until a block of contiguous requires is found.
	return !c.finished &&
		(parent.IsNil() || parent.IsRoot() || parent.IsScript() || parent.IsModuleBody())
}

// Java: CheckRequiresSorted.visit(NodeTraversal, Node, Node)
func (c *CheckRequiresSorted) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if _, ok := ast.AsScript(n); ok {
		c.CheckCanonical(t)
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
	if !callNode.IsNil() && CheckRequiresSorted_IsValidImportCall(callNode) {
		stmt := CheckRequiresSorted_ParseImport(callNode)
		c.originalImports = append(c.originalImports, stmt)
		ns := stmt.namespace
		existing, seen := c.importsByNamespace[ns]
		if !seen {
			c.namespaceOrder = append(c.namespaceOrder, ns)
		}
		c.importsByNamespace[ns] = append(existing, stmt)
		if c.firstNode.IsNil() {
			c.firstNode = n
		}
		c.lastNode = n
	} else if len(c.importsByNamespace) > 0 {
		c.finished = true
	}
}

// Java: CheckRequiresSorted.isValidImportCall(Node)
func CheckRequiresSorted_IsValidImportCall(n rhino.Node) bool {
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

// Java: CheckRequiresSorted.parseImport(Node)
func CheckRequiresSorted_ParseImport(callNode rhino.Node) *CheckRequiresSorted_ImportStatement {
	call := ast.MustCall(callNode)
	primitive := CheckRequiresSorted_ImportPrimitive_FromName(*call.Callee().Node().GetQualifiedName())
	namespace := ast.MustStringLit(call.MustFirstArg().Node()).Value()

	parent := callNode.GetParent()
	if parent.IsExprResult() {
		// goog.require('a');
		return CheckRequiresSorted_ImportStatement_Of(
			[]rhino.Node{parent},
			primitive,
			namespace,
			"",
			false,
			/* destructures= */ nil)
	}

	grandparent := parent.GetParent()
	if name, ok := ast.AsName(parent); ok {
		// const a = goog.require('a');
		return CheckRequiresSorted_ImportStatement_Of(
			[]rhino.Node{grandparent}, primitive, namespace, name.Ident(), true /* destructures= */, nil)
	}
	// const {a: b, c} = goog.require('a');
	// ImmutableList.builder(): non-null even when empty.
	destructures := []*CheckRequiresSorted_DestructuringBinding{}
	if lhs, ok := ast.AsDestructuringLhs(parent); ok {
		if objPattern, ok := lhs.Pattern().AsObjectPattern(); ok {
			for prop := range objPattern.Props() {
				if key, ok := prop.AsStringKey(); ok {
					if local, ok := key.Value().AsName(); ok {
						destructures = append(
							destructures,
							CheckRequiresSorted_DestructuringBinding_Of(
								key.Key(), local.Ident(), key.Node().IsShorthandProperty()))
						continue
					}
				}
				name := prop.Node()
				destructures = append(
					destructures,
					CheckRequiresSorted_DestructuringBinding_Of(
						name.GetString(), name.GetFirstChild().GetString(), name.IsShorthandProperty()))
			}
			return CheckRequiresSorted_ImportStatement_Of(
				[]rhino.Node{grandparent},
				primitive,
				namespace,
				"",
				false,
				destructures)
		}
	}
	for name := parent.GetFirstFirstChild(); !name.IsNil(); name = name.GetNext() {
		exportedName := name.GetString()
		localName := name.GetFirstChild().GetString()
		isShorthandProperty := name.IsShorthandProperty()
		destructures = append(destructures, CheckRequiresSorted_DestructuringBinding_Of(exportedName, localName, isShorthandProperty))
	}
	return CheckRequiresSorted_ImportStatement_Of(
		[]rhino.Node{grandparent},
		primitive,
		namespace,
		"",
		false,
		destructures)
}

// Java: CheckRequiresSorted.checkCanonical(NodeTraversal)
func (c *CheckRequiresSorted) CheckCanonical(t *NodeTraversal) {
	canonicalImports := CheckRequiresSorted_CanonicalizeImports(c.importsByNamespace, c.namespaceOrder)
	if !slices.EqualFunc(c.originalImports, canonicalImports, (*CheckRequiresSorted_ImportStatement).Equals) {
		c.needsFix = true
		var sb strings.Builder
		for i, stmt := range canonicalImports {
			if i > 0 {
				sb.WriteByte('\n')
			}
			sb.WriteString(stmt.Format())
		}
		c.replacement = sb.String()
		c.hasReplacement = true
		if c.mode == CheckRequiresSorted_Mode_COLLECT_AND_REPORT {
			t.Report(c.firstNode, CheckRequiresSorted_REQUIRES_NOT_SORTED, c.replacement)
		}
	}
}

// CheckRequiresSorted_CanonicalizeImports canonicalizes a list of import
// statements by deduplicating and merging imports for the same namespace,
// and sorting the result.
//
// Java: CheckRequiresSorted.canonicalizeImports(Multimap<String, ImportStatement>)
func CheckRequiresSorted_CanonicalizeImports(
	importsByNamespace map[string][]*CheckRequiresSorted_ImportStatement,
	namespaceOrder []string) []*CheckRequiresSorted_ImportStatement {
	if len(namespaceOrder) == 0 {
		return nil
	}
	canonicalImports := make([]*CheckRequiresSorted_ImportStatement, 0, len(namespaceOrder))

	for _, namespace := range namespaceOrder {
		allImports := importsByNamespace[namespace]

		// Find the strongest primitive across all existing imports. Every emitted import for this
		// namespace will use this primitive.
		strongestPrimitive := CheckRequiresSorted_ImportPrimitive_WEAKEST
		for _, i := range allImports {
			strongestPrimitive = CheckRequiresSorted_ImportPrimitive_Stronger(strongestPrimitive, i.primitive)
		}

		// Emit each aliasing import separately, as deduplicating them would require code references
		// to be rewritten.
		hasAliasing := false
		for _, stmt := range allImports {
			if !stmt.hasAlias {
				continue
			}
			canonicalImports = append(canonicalImports, stmt.Upgrade(strongestPrimitive))
			hasAliasing = true
		}

		// Emit a single destructuring import with a non-empty pattern, merged from the existing
		// destructuring imports.
		hasDestructuring := false
		destructuringNodes := []rhino.Node{}
		var allDestructures []*CheckRequiresSorted_DestructuringBinding
		for _, i := range allImports {
			if i.destructures != nil {
				destructuringNodes = append(destructuringNodes, i.nodes...)
				for _, d := range i.destructures {
					allDestructures = append(allDestructures, d.CanonicalizeShorthandProperties())
				}
			}
		}
		if len(allDestructures) > 0 {
			slices.SortFunc(allDestructures, (*CheckRequiresSorted_DestructuringBinding).CompareTo)
			destructures := slices.CompactFunc(allDestructures, (*CheckRequiresSorted_DestructuringBinding).Equals)
			canonicalImports = append(canonicalImports,
				CheckRequiresSorted_ImportStatement_Of(
					destructuringNodes,
					strongestPrimitive,
					namespace,
					"",
					false,
					destructures))
			hasDestructuring = true
		}

		// Emit a standalone import unless an aliasing or destructuring one already exists.
		if !hasAliasing && !hasDestructuring {
			standaloneNodes := []rhino.Node{}
			for _, i := range allImports {
				if i.IsStandalone() {
					standaloneNodes = append(standaloneNodes, i.nodes...)
				}
			}
			canonicalImports = append(canonicalImports,
				CheckRequiresSorted_ImportStatement_Of(
					standaloneNodes,
					strongestPrimitive,
					namespace,
					"",
					false,
					/* destructures= */ nil))
		}
	}

	// Sorting by natural order yields the correct result due to the implementation of
	// ImportStatement#compareTo.
	slices.SortStableFunc(canonicalImports, (*CheckRequiresSorted_ImportStatement).CompareTo)

	return canonicalImports
}
