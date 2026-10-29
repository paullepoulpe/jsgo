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

	// Maps each namespace into the existing import statements for that namespace.
	// Use an ArrayListMultimap so that values for a key are iterated in a deterministic order.
	//
	// Go-only: language-adapter — ArrayListMultimap.create(): a map from each
	// key to the list of its values in insertion order. Java's keys iterate in
	// java.util.HashMap order; here they iterate in insertion order. The only
	// iteration (canonicalizeImports) is followed by a sort on the formatted
	// text, which contains the namespace, so the key order is not observable.
	importsByNamespace *jutil.LinkedHashMap[string, []*CheckRequiresSorted_ImportStatement]

	// The import statements in the order they appear.
	originalImports []*CheckRequiresSorted_ImportStatement

	// The import statements in canonical order.
	firstNode rhino.Node
	lastNode  rhino.Node
	finished  bool

	needsFix    bool
	replacement *string
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
	for _, primitive := range CheckRequiresSorted_ImportPrimitive_values() {
		if checkRequiresSorted_ImportPrimitive_name[primitive] == name {
			return primitive
		}
	}
	panic(jutil.NewIllegalArgumentExceptionMsg("Invalid primitive name " + name))
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
	if p1.Ordinal() < p2.Ordinal() {
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
	// Go-only: language-adapter — requireNonNull(exportedName/localName): Go
	// strings cannot be null.
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
	// ComparisonChain.start().compare(exportedName).compare(localName).result()
	if c := jutil.StringCompareTo(b.ExportedName(), other.ExportedName()); c != 0 {
		if c < 0 {
			return -1
		}
		return 1
	}
	if c := jutil.StringCompareTo(b.LocalName(), other.LocalName()); c != 0 {
		if c < 0 {
			return -1
		}
		return 1
	}
	return 0
}

// IsCanonical returns true if the destructuring binding is not canonical.
//
// For example:
//   - `{Foo}` is canonical
//   - `{Foo: Bar}` is canonical
//   - `{Foo: Foo}` is not canonical
//
// Java: CheckRequiresSorted.DestructuringBinding.isCanonical()
func (b *CheckRequiresSorted_DestructuringBinding) IsCanonical() bool {
	return b.ExportedName() != b.LocalName() || b.IsShorthandProperty()
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
		b.ExportedName(), b.LocalName() /* isShorthandProperty= */, true)
}

// Equals is the record's component-wise equals.
func (b *CheckRequiresSorted_DestructuringBinding) Equals(o any) bool {
	other, ok := o.(*CheckRequiresSorted_DestructuringBinding)
	if !ok || other == nil {
		return false
	}
	return b.exportedName == other.exportedName &&
		b.localName == other.localName &&
		b.isShorthandProperty == other.isShorthandProperty
}

// CheckRequiresSorted_ImportStatement is an import statement, which may have
// been merged from several import statements for the same namespace in the
// original code.
//
// An import statement has exactly one of three shapes:
//
//   - Standalone: has no LHS, as in `goog.require('namespace')`.
//   - Aliasing: has an LHS with an alias, as in `const alias =
//     goog.require('namespace')`.
//   - Destructuring: has an LHS with a destructuring pattern, as in `const
//     {name: localName} = goog.require('namespace')`.
//
// nodes: the nodes this import statement was merged from, in source order.
// primitive: the import primitive being called. namespace: the namespace
// being imported. alias: the alias for an aliasing import (when hasAlias is
// true). destructures: the destructures for a destructuring import in source
// order, or nil if the import isn't destructuring. If the import is
// destructuring but the pattern is empty, the value is non-nil but empty.
type CheckRequiresSorted_ImportStatement struct {
	nodes        []rhino.Node
	namespace    string
	alias        string
	destructures []*CheckRequiresSorted_DestructuringBinding
	primitive    CheckRequiresSorted_ImportPrimitive
	hasAlias     bool
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
	return !s.IsAliasing() && !s.IsDestructuring()
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
	if CheckRequiresSorted_ImportPrimitive_Stronger(s.Primitive(), otherPrimitive) != s.Primitive() {
		return newCheckRequiresSorted_ImportStatement(s.nodes, otherPrimitive, s.namespace, s.alias, s.hasAlias, s.destructures)
	}
	return s
}

// Java: CheckRequiresSorted.ImportStatement.formatWithoutDoc()
func (s *CheckRequiresSorted_ImportStatement) FormatWithoutDoc() string {
	var sb strings.Builder

	if !s.IsStandalone() {
		sb.WriteString("const ")
	}
	if s.hasAlias {
		sb.WriteString(s.alias)
	}
	if s.destructures != nil {
		sb.WriteString("{")
		first := true
		for _, binding := range s.destructures {
			exportedName := binding.ExportedName()
			localName := binding.LocalName()
			if first {
				first = false
			} else {
				sb.WriteString(", ")
			}
			sb.WriteString(exportedName)
			if exportedName != localName {
				sb.WriteString(": ")
				sb.WriteString(localName)
			}
		}
		sb.WriteString("}")
	}
	if !s.IsStandalone() {
		sb.WriteString(" = ")
	}
	sb.WriteString(s.Primitive().String())
	sb.WriteString("('")
	sb.WriteString(s.Namespace())
	sb.WriteString("');")

	return sb.String()
}

// Format formats the import statement into code.
//
// Java: CheckRequiresSorted.ImportStatement.format()
func (s *CheckRequiresSorted_ImportStatement) Format() string {
	var sb strings.Builder

	for _, node := range s.Nodes() {
		comment := node.GetNonJSDocCommentString()
		if comment != "" {
			sb.WriteString(comment)
			sb.WriteString("\n")
		}
		jsDoc := NodeUtil_GetBestJSDocInfo(node)
		if jsDoc != nil {
			// StringBuilder.append(String) appends "null" for a null string.
			if original := jsDoc.GetOriginalCommentString(); original != nil {
				sb.WriteString(*original)
			} else {
				sb.WriteString("null")
			}
			sb.WriteString("\n")
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

// Equals is the record's component-wise equals: the lists compare element by
// element (Node by identity), alias and destructures with Objects.equals.
func (s *CheckRequiresSorted_ImportStatement) Equals(o any) bool {
	other, ok := o.(*CheckRequiresSorted_ImportStatement)
	if !ok || other == nil {
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
		func(a, b *CheckRequiresSorted_DestructuringBinding) bool { return a.Equals(b) })
}

var checkRequiresSorted_GOOG_REQUIRE = rhino.QualifiedName_Of("goog.require")

var checkRequiresSorted_GOOG_REQUIRETYPE = rhino.QualifiedName_Of("goog.requireType")

var checkRequiresSorted_GOOG_FORWARDDECLARE = rhino.QualifiedName_Of("goog.forwardDeclare")

// NewCheckRequiresSorted is new CheckRequiresSorted(mode).
//
// Java: CheckRequiresSorted(CheckRequiresSorted.Mode)
func NewCheckRequiresSorted(mode CheckRequiresSorted_Mode) *CheckRequiresSorted {
	return &CheckRequiresSorted{
		importsByNamespace: jutil.NewLinkedHashMap[string, []*CheckRequiresSorted_ImportStatement](),
		firstNode:          rhino.NilNode,
		lastNode:           rhino.NilNode,
		finished:           false,
		needsFix:           false,
		replacement:        nil,
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
	return c.replacement
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
	if n.IsScript() {
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
		// importsByNamespace.put(stmt.namespace(), stmt)
		c.importsByNamespace.Put(stmt.Namespace(), append(c.importsByNamespace.Get(stmt.Namespace()), stmt))
		if c.firstNode.IsNil() {
			c.lastNode = n
			c.firstNode = c.lastNode
		} else {
			c.lastNode = n
		}
	} else if !c.importsByNamespace.IsEmpty() {
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
						// {a: a} and {a: b} both yield false
						// {a} yields true
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
	canonicalImports := CheckRequiresSorted_CanonicalizeImports(c.importsByNamespace)
	if !slices.EqualFunc(c.originalImports, canonicalImports,
		func(a, b *CheckRequiresSorted_ImportStatement) bool { return a.Equals(b) }) {
		c.needsFix = true
		formatted := make([]string, len(canonicalImports))
		for i, stmt := range canonicalImports {
			formatted[i] = stmt.Format()
		}
		c.replacement = jutil.StringPtr(strings.Join(formatted, "\n"))
		if c.mode == CheckRequiresSorted_Mode_COLLECT_AND_REPORT {
			t.Report(c.firstNode, CheckRequiresSorted_REQUIRES_NOT_SORTED, *c.replacement)
		}
	}
}

// CheckRequiresSorted_CanonicalizeImports canonicalizes a list of import
// statements by deduplicating and merging imports for the same namespace,
// and sorting the result.
//
// Java: CheckRequiresSorted.canonicalizeImports(Multimap<String, ImportStatement>)
func CheckRequiresSorted_CanonicalizeImports(
	importsByNamespace *jutil.LinkedHashMap[string, []*CheckRequiresSorted_ImportStatement]) []*CheckRequiresSorted_ImportStatement {
	var canonicalImports []*CheckRequiresSorted_ImportStatement

	for namespace := range importsByNamespace.Keys() {
		allImports := importsByNamespace.Get(namespace)

		// Find the strongest primitive across all existing imports. Every emitted import for this
		// namespace will use this primitive. This makes the logic simpler and cannot change runtime
		// behavior, but may produce spurious changes when multiple aliasing imports of differing
		// strength exist (which are already in violation of the style guide).
		strongestPrimitive := CheckRequiresSorted_ImportPrimitive_WEAKEST
		for _, i := range allImports {
			strongestPrimitive = CheckRequiresSorted_ImportPrimitive_Stronger(strongestPrimitive, i.Primitive())
		}

		// Emit each aliasing import separately, as deduplicating them would require code references
		// to be rewritten.
		hasAliasing := false
		for _, stmt := range allImports {
			if !stmt.IsAliasing() {
				continue
			}
			canonicalImports = append(canonicalImports, stmt.Upgrade(strongestPrimitive))
			hasAliasing = true
		}

		// Emit a single destructuring import with a non-empty pattern, merged from the existing
		// destructuring imports.
		hasDestructuring := false
		destructuringNodes := []rhino.Node{}
		for _, i := range allImports {
			if i.IsDestructuring() {
				destructuringNodes = append(destructuringNodes, i.Nodes()...)
			}
		}
		// .flatMap(destructures).map(canonicalizeShorthandProperties).sorted().distinct()
		var allDestructures []*CheckRequiresSorted_DestructuringBinding
		for _, i := range allImports {
			if i.IsDestructuring() {
				for _, d := range i.Destructures() {
					allDestructures = append(allDestructures, d.CanonicalizeShorthandProperties())
				}
			}
		}
		slices.SortStableFunc(allDestructures, func(a, b *CheckRequiresSorted_DestructuringBinding) int {
			return a.CompareTo(b)
		})
		destructures := []*CheckRequiresSorted_DestructuringBinding{}
		for _, d := range allDestructures {
			if !slices.ContainsFunc(destructures, func(e *CheckRequiresSorted_DestructuringBinding) bool { return e.Equals(d) }) {
				destructures = append(destructures, d)
			}
		}
		if len(destructures) != 0 {
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
					standaloneNodes = append(standaloneNodes, i.Nodes()...)
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
	slices.SortStableFunc(canonicalImports, func(a, b *CheckRequiresSorted_ImportStatement) int {
		return a.CompareTo(b)
	})

	return canonicalImports
}
