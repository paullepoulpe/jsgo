// Port of com/google/javascript/jscomp/CheckMissingRequires.java @ closure-compiler v20260415
//
// Detects references to fully qualified Closure namespaces that
// should be imported via goog.require / goog.requireType.
//
// ModuleMetadata / ModuleMetadataMap are stubbed below with a single-file
// gatherer (cmrGatherModuleMetadata) that mirrors the relevant subset of
// GatherModuleMetadata.java — enough to power CheckMissingRequires for
// single-file analysis. The QualifiedName class and GoogRequire / Var helpers
// remain stubbed pending a full scope-system port.

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ---------------------------------------------------------------------------
// Stubbed module-system types
// ---------------------------------------------------------------------------

// cmrModuleMetadata is a stub for
// com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata.
//
// TODO: replace with a real port of ModuleMetadata once the module
// system (modules/ package) is ported. The real class tracks:
//   - isModule() / isGoogModule() / isNonLegacyGoogModule()
//   - googNamespaces() — the set of goog.provide / goog.module names
//   - stronglyRequiredGoogNamespaces() / weaklyRequiredGoogNamespaces()
//   - hasLegacyGoogNamespaces() — true when any namespace is a legacy provide
//   - rootNode() — the SCRIPT node for the file
type cmrModuleMetadata struct {
	// sourceFileName is the file path for this module.
	sourceFileName string
	// googNamespaces is the set of namespaces declared by this file.
	googNamespaces map[string]bool
	// strongRequires is the set of namespaces strongly required by this file.
	strongRequires map[string]bool
	// weakRequires is the set of namespaces weakly required by this file.
	weakRequires map[string]bool
	// isModule reports whether this file is a goog.module file.
	isModuleFlag bool
	// hasLegacyNS reports whether any namespace is a legacy provide.
	hasLegacyNS bool
	// isNonLegacy reports true for goog.module without declareLegacyNamespace.
	isNonLegacy bool
	// rootNode is the SCRIPT node.
	rootNode *rhino.Node
}

func (m *cmrModuleMetadata) isModule() bool { return m != nil && m.isModuleFlag }

func (m *cmrModuleMetadata) hasLegacyGoogNamespaces() bool { return m != nil && m.hasLegacyNS }

func (m *cmrModuleMetadata) isNonLegacyGoogModule() bool { return m != nil && m.isNonLegacy }

func (m *cmrModuleMetadata) containsNamespace(ns string) bool {
	return m != nil && m.googNamespaces[ns]
}

func (m *cmrModuleMetadata) containsStrongRequire(ns string) bool {
	return m != nil && m.strongRequires[ns]
}

func (m *cmrModuleMetadata) containsWeakRequire(ns string) bool {
	return m != nil && m.weakRequires[ns]
}

func (m *cmrModuleMetadata) getRootSourceFileName() string {
	if m == nil {
		return ""
	}
	return m.sourceFileName
}

// cmrModuleMetadataMap is a stub for
// com.google.javascript.jscomp.modules.ModuleMetadataMap.
//
// TODO: replace with a real port of ModuleMetadataMap. The real class
// exposes getModulesByGoogNamespace() returning an ImmutableMap from
// namespace string to ModuleMetadata.
type cmrModuleMetadataMap struct {
	// byNamespace maps a Closure namespace to the metadata of the file that provides it.
	byNamespace map[string]*cmrModuleMetadata
}

func (m *cmrModuleMetadataMap) get(ns string) *cmrModuleMetadata {
	if m == nil {
		return nil
	}
	return m.byNamespace[ns]
}

// ---------------------------------------------------------------------------
// cmrGatherModuleMetadata — single-file port of GatherModuleMetadata.java
// ---------------------------------------------------------------------------

// cmrGatherModuleMetadata walks `script` and builds a cmrModuleMetadataMap
// for the single file identified by `path`. It mirrors the subset of
// GatherModuleMetadata.java needed by CheckMissingRequires:
//   - goog.provide(ns)                → GOOG_PROVIDE; googNamespaces += ns
//   - goog.module(ns)                 → GOOG_MODULE;  googNamespaces += ns
//   - goog.module.declareLegacyNamespace() → upgrades to LEGACY_GOOG_MODULE
//   - goog.require(ns)                → strongRequires += ns
//   - goog.requireType(ns)            → weakRequires += ns
//
// The returned map contains exactly one entry per declared namespace,
// mirroring ModuleMetadataMap.getModulesByGoogNamespace().
func cmrGatherModuleMetadata(path string, script *rhino.Node) *cmrModuleMetadataMap {
	md := &cmrModuleMetadata{
		sourceFileName: path,
		googNamespaces: map[string]bool{},
		strongRequires: map[string]bool{},
		weakRequires:   map[string]bool{},
		rootNode:       script,
	}
	// moduleType tracks SCRIPT(0) / GOOG_PROVIDE(1) / GOOG_MODULE(2).
	// 0 = SCRIPT, 1 = GOOG_PROVIDE, 2 = GOOG_MODULE/LEGACY_GOOG_MODULE
	moduleType := 0
	declareLegacy := false
	cmrScanGoogCalls(script, func(callee, arg string) {
		switch callee {
		case "goog.provide":
			if arg != "" {
				md.googNamespaces[arg] = true
				if moduleType < 1 {
					moduleType = 1
				}
			}
		case "goog.module":
			if arg != "" {
				md.googNamespaces[arg] = true
				if moduleType < 2 {
					moduleType = 2
				}
			}
		case "goog.module.declareLegacyNamespace":
			declareLegacy = true
		case "goog.require":
			if arg != "" {
				md.strongRequires[arg] = true
			}
		case "goog.requireType":
			if arg != "" {
				md.weakRequires[arg] = true
			}
		}
	})
	switch moduleType {
	case 1:
		md.hasLegacyNS = true
	case 2:
		md.isModuleFlag = true
		md.isNonLegacy = !declareLegacy
		if declareLegacy {
			md.hasLegacyNS = true
		}
	}
	// Build the byNamespace map: one entry per declared namespace.
	byNS := map[string]*cmrModuleMetadata{}
	for ns := range md.googNamespaces {
		// Each entry is a clone with exactly one namespace, matching Java's
		// getModulesByGoogNamespace() contract (each value has size-1 googNamespaces).
		entry := &cmrModuleMetadata{
			sourceFileName: md.sourceFileName,
			googNamespaces: map[string]bool{ns: true},
			strongRequires: md.strongRequires,
			weakRequires:   md.weakRequires,
			isModuleFlag:   md.isModuleFlag,
			hasLegacyNS:    md.hasLegacyNS,
			isNonLegacy:    md.isNonLegacy,
			rootNode:       md.rootNode,
		}
		byNS[ns] = entry
	}
	return &cmrModuleMetadataMap{byNamespace: byNS}
}

// cmrScanGoogCalls walks `n` in pre-order and calls `fn` for every CALL
// node whose callee is a goog.* qualified name. `arg` is the string value
// of the first argument (or "" when absent / non-string).
func cmrScanGoogCalls(n *rhino.Node, fn func(callee, arg string)) {
	if n == nil {
		return
	}
	if n.Token() == rhino.CALL {
		callee := n.FirstChild()
		if callee != nil {
			qn := cmrQualifiedNameFromNode(callee)
			if qn != nil && strings.HasPrefix(qn.join(), "goog") {
				arg := ""
				kids := n.Children()
				if len(kids) >= 2 && kids[1].Token() == rhino.STRINGLIT {
					arg = kids[1].String_()
				}
				fn(qn.join(), arg)
			}
		}
	}
	for _, ch := range n.Children() {
		cmrScanGoogCalls(ch, fn)
	}
}

// cmrQualifiedName is a stub for com.google.javascript.rhino.QualifiedName.
//
// TODO: replace with a real port of QualifiedName. The real class
// represents dotted names ("a.b.c") and supports:
//   - getRoot() — first component ("a")
//   - join() — the full dotted string
//   - isSimple() — no dots
//   - getOwner() — the parent name (nil for simple names)
//   - getprop(s) — append a component
//   - components() — iterator over individual parts
//   - getComponentCount() — number of parts
//   - matches(Node) — structural match against a NAME/GETPROP chain
type cmrQualifiedName struct {
	components []string
}

func cmrQualifiedNameOf(dotted string) *cmrQualifiedName {
	if dotted == "" {
		return &cmrQualifiedName{}
	}
	return &cmrQualifiedName{components: strings.Split(dotted, ".")}
}

func (q *cmrQualifiedName) getRoot() string {
	if q == nil || len(q.components) == 0 {
		return ""
	}
	return q.components[0]
}

func (q *cmrQualifiedName) join() string {
	if q == nil {
		return ""
	}
	return strings.Join(q.components, ".")
}

func (q *cmrQualifiedName) isSimple() bool {
	return q != nil && len(q.components) <= 1
}

func (q *cmrQualifiedName) getOwner() *cmrQualifiedName {
	if q == nil || len(q.components) <= 1 {
		return nil
	}
	return &cmrQualifiedName{components: q.components[:len(q.components)-1]}
}

func (q *cmrQualifiedName) getprop(s string) *cmrQualifiedName {
	if q == nil {
		return cmrQualifiedNameOf(s)
	}
	newComponents := make([]string, len(q.components)+1)
	copy(newComponents, q.components)
	newComponents[len(q.components)] = s
	return &cmrQualifiedName{components: newComponents}
}

func (q *cmrQualifiedName) getComponentCount() int {
	if q == nil {
		return 0
	}
	return len(q.components)
}

// cmrQualifiedNameFromNode converts a NAME or GETPROP chain into a
// cmrQualifiedName. Returns nil when `n` is not a qualified name chain.
func cmrQualifiedNameFromNode(n *rhino.Node) *cmrQualifiedName {
	if n == nil {
		return nil
	}
	switch n.Token() {
	case rhino.NAME:
		if n.String_() == "" {
			return nil
		}
		return cmrQualifiedNameOf(n.String_())
	case rhino.GETPROP:
		kids := n.Children()
		if len(kids) < 2 {
			return nil
		}
		owner := cmrQualifiedNameFromNode(kids[0])
		if owner == nil {
			return nil
		}
		return owner.getprop(kids[1].String_())
	}
	return nil
}

// cmrMatchesQualifiedName reports whether a GETPROP/NAME chain matches
// the given qualified name (mirrors QualifiedName.matches(Node)).
func cmrMatchesQualifiedName(n *rhino.Node, q *cmrQualifiedName) bool {
	if q == nil || n == nil {
		return false
	}
	qn := cmrQualifiedNameFromNode(n)
	if qn == nil {
		return false
	}
	return qn.join() == q.join()
}

// cmrGoogRequire is a stub for NodeUtil.GoogRequire (an inner record
// in the Java NodeUtil class).
//
// TODO: replace with a real port of NodeUtil.GoogRequire. The real
// record carries:
//   - namespace() — the required namespace string
//   - property() — the destructured property name, or nil for whole-namespace imports
//   - isStrongRequire() — true for goog.require, false for goog.requireType
type cmrGoogRequire struct {
	namespace        string
	property         string // empty string means nil in Java
	hasProperty      bool
	isStrongRequire_ bool
}

func (r *cmrGoogRequire) isStrongRequire() bool {
	return r != nil && r.isStrongRequire_
}

// cmrVar is a stub for com.google.javascript.jscomp.Var.
//
// TODO: replace with a real port of Var. The real class tracks:
//   - getScope() — the TypedScope that declares this variable
//   - getDeclaration() — the declaring Var (for aliased vars)
//   - getNode() — the NAME node of the declaration
//   - getInitialValue() — the initializer expression node
//   - getScopeRoot() — the root node of the declaring scope
type cmrVar struct {
	name         string
	declNode     *rhino.Node
	initialValue *rhino.Node
	scopeRoot    *rhino.Node
	isLocal      bool
}

func (v *cmrVar) getInitialValue() *rhino.Node {
	if v == nil {
		return nil
	}
	return v.initialValue
}

func (v *cmrVar) getScopeRoot() *rhino.Node {
	if v == nil {
		return nil
	}
	return v.scopeRoot
}

// ---------------------------------------------------------------------------
// cmrStrength — mirrors the Java Strength enum
// ---------------------------------------------------------------------------

// cmrStrength represents the strength of a namespace reference: whether
// it requires a goog.require (strong) or goog.requireType (weak/type-only).
// Mirrors CheckMissingRequires.Strength in Java.
type cmrStrength int

const (
	// cmrStrengthCode — a runtime code reference (strong, not type-only).
	cmrStrengthCode cmrStrength = iota
	// cmrStrengthImplementsExtends — @implements / @extends (strong + type-only).
	cmrStrengthImplementsExtends
	// cmrStrengthWeakType — other JSDoc annotations (weak / type-only).
	cmrStrengthWeakType
)

func (s cmrStrength) isStrong() bool {
	return s == cmrStrengthCode || s == cmrStrengthImplementsExtends
}

func (s cmrStrength) isTypeOnly() bool {
	return s == cmrStrengthImplementsExtends || s == cmrStrengthWeakType
}

// ---------------------------------------------------------------------------
// Diagnostic codes
// ---------------------------------------------------------------------------

// Diagnostic codes mirror the Java constants in CheckMissingRequires.
// Their string values must match what the Java compiler emits.
// Disabled codes (DiagnosticType.disabled) are registered here for
// completeness but are not currently emitted by the pass.
const (
	// cmrCodeMissingRequire — JSC_MISSING_REQUIRE (warning)
	// Fired when a file references a fully qualified Closure namespace that
	// should be aliased via goog.require.
	cmrCodeMissingRequire = diag.JSC_MISSING_REQUIRE

	// cmrCodeMissingRequireType — JSC_MISSING_REQUIRE_TYPE (disabled in Java)
	// Type-only variant of JSC_MISSING_REQUIRE; emitted for type-only references.
	cmrCodeMissingRequireType = diag.JSC_MISSING_REQUIRE_TYPE

	// cmrCodeIncorrectNSAliasRequire — JSC_INCORRECT_NAMESPACE_ALIAS_REQUIRE (disabled)
	// Fired when a destructured import creates an alias for a sub-namespace that
	// has its own goog.require namespace.
	cmrCodeIncorrectNSAliasRequire = diag.JSC_INCORRECT_NAMESPACE_ALIAS_REQUIRE

	// cmrCodeIncorrectNSAliasRequireType — JSC_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE (disabled)
	cmrCodeIncorrectNSAliasRequireType = diag.JSC_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE

	// cmrCodeIndirectNSRefRequire — JSC_INDIRECT_NAMESPACE_REF_REQUIRE (disabled)
	// Fired when a reference goes through an alias that covers a parent namespace
	// while the sub-namespace has its own goog.require.
	cmrCodeIndirectNSRefRequire = diag.JSC_INDIRECT_NAMESPACE_REF_REQUIRE

	// cmrCodeIndirectNSRefRequireType — JSC_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE (disabled)
	cmrCodeIndirectNSRefRequireType = diag.JSC_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE

	// cmrCodeMissingRequireInProvidesFile — JSC_MISSING_REQUIRE_IN_PROVIDES_FILE (warning)
	// Fired in goog.provide files where a referenced namespace is not required.
	cmrCodeMissingRequireInProvidesFile = diag.JSC_MISSING_REQUIRE_IN_PROVIDES_FILE

	// cmrCodeMissingRequireTypeInProvidesFile — JSC_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE (disabled)
	cmrCodeMissingRequireTypeInProvidesFile = diag.JSC_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE

	// cmrCodeNonLegacyGoogModuleReference — JSC_NON_LEGACY_GOOG_MODULE_REFERENCE (error)
	// Fired when a name defined only in a non-legacy goog.module is used
	// as a fully-qualified runtime reference.
	cmrCodeNonLegacyGoogModuleReference = diag.JSC_NON_LEGACY_GOOG_MODULE_REFERENCE

	// cmrCodeMissingRequireInGoogScope — JSC_MISSING_REQUIRE_IN_GOOG_SCOPE (warning)
	// Fired inside goog.scope bodies when a namespace alias is used without a require.
	cmrCodeMissingRequireInGoogScope = diag.JSC_MISSING_REQUIRE_IN_GOOG_SCOPE

	// cmrCodeMissingRequireTypeInGoogScope — JSC_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE (warning)
	cmrCodeMissingRequireTypeInGoogScope = diag.JSC_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE

	// cmrCodeMissingRequireForGoogModuleGet — JSC_MISSING_REQUIRE_FOR_GOOG_MODULE_GET (warning)
	// Fired when goog.module.get is called at script load time for a namespace
	// that is not strongly required by the current file.
	cmrCodeMissingRequireForGoogModuleGet = diag.JSC_MISSING_REQUIRE_FOR_GOOG_MODULE_GET
)

// Ensure all code constants are referenced so the compiler does not
// complain about unused identifiers. These variables are not exported.
var _ = cmrCodeMissingRequireType
var _ = cmrCodeIncorrectNSAliasRequire
var _ = cmrCodeIncorrectNSAliasRequireType
var _ = cmrCodeIndirectNSRefRequire
var _ = cmrCodeIndirectNSRefRequireType
var _ = cmrCodeMissingRequireTypeInProvidesFile
var _ = cmrCodeMissingRequireTypeInGoogScope

// ---------------------------------------------------------------------------
// CheckMissingRequires pass
// ---------------------------------------------------------------------------

// cmrGoogScope is the QualifiedName for "goog.scope" — used to identify
// goog.scope call bodies. Mirrors the Java static field GOOG_SCOPE.
var cmrGoogScope = cmrQualifiedNameOf("goog.scope")

// CheckMissingRequires detects references to fully qualified Closure
// namespaces in files that should import them via goog.require /
// goog.requireType. Mirrors CheckMissingRequires.java.
//
// The pass is disabled until the module-system stubs are replaced with
// real implementations (cmrModuleMetadata / cmrModuleMetadataMap). At
// that point, wire NewCheckMissingRequires into CheckSourceWithOptions.
type CheckMissingRequires struct {
	path string

	// moduleByNamespace maps a Closure namespace to the ModuleMetadata
	// of the file that provides it. Mirrors the Java field of the same name.
	// TODO: populate once ModuleMetadataMap is ported.
	moduleByNamespace map[string]*cmrModuleMetadata

	// templateParamNames is the set of @template parameter names seen so far
	// in the file currently being checked. Cleared when a file's root node is
	// visited in post-order. Mirrors the Java LinkedHashSet.
	templateParamNames map[string]bool

	// controlFlowScopeDepth mirrors the Java field: tracks how many
	// "control flow scopes" we have entered beyond the global scope.
	// IIFEs and goog.scope bodies at the global level do not increment it.
	controlFlowScopeDepth int

	diags []diag.Diagnostic
}

// NewCheckMissingRequires constructs the pass, mirroring the Java constructor
// signature (compiler + moduleMetadataMap). Pass a map built by
// cmrGatherModuleMetadata; nil map causes the pass to no-op.
func NewCheckMissingRequires(path string, moduleMetadataMap *cmrModuleMetadataMap) *CheckMissingRequires {
	c := &CheckMissingRequires{
		path:               path,
		templateParamNames: map[string]bool{},
	}
	if moduleMetadataMap != nil {
		c.moduleByNamespace = moduleMetadataMap.byNamespace
	}
	return c
}

// Check runs the pass over `script` and returns any diagnostics found.
// Mirrors the Java CompilerPass.process entry point; externs are ignored
// here (consistent with the narrow-port contract).
func (c *CheckMissingRequires) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.templateParamNames = map[string]bool{}
	c.controlFlowScopeDepth = 0

	// Without a real ModuleMetadataMap we cannot determine which namespaces
	// exist or which files provide/require them. The caller should supply one
	// via cmrGatherModuleMetadata; if nil, no-op to avoid false positives.
	if c.moduleByNamespace == nil {
		return nil
	}

	c.walk(script, nil, nil)
	return c.diags
}

// ---------------------------------------------------------------------------
// Tree walker
// ---------------------------------------------------------------------------

// walk is a pre-order / post-order recursive traversal mirroring the Java
// AbstractModuleCallback contract (shouldTraverse pre-order, visit post-order).
// `currentModule` is the ModuleMetadata for the enclosing file; for the
// top-level SCRIPT node it is resolved from the moduleByNamespace map by
// finding the metadata whose rootNode matches, mirroring AbstractModuleCallback
// enterModule behaviour.
func (c *CheckMissingRequires) walk(n *rhino.Node, parent *rhino.Node, currentModule *cmrModuleMetadata) {
	if n == nil {
		return
	}

	// Resolve currentModule for SCRIPT nodes (AbstractModuleCallback.enterModule).
	// We look for any cmrModuleMetadata whose rootNode is this SCRIPT node.
	if n.Token() == rhino.SCRIPT && currentModule == nil {
		for _, md := range c.moduleByNamespace {
			if md.rootNode == n {
				currentModule = md
				break
			}
		}
		// If no namespace was declared (plain script), use a synthetic metadata
		// with the rootNode set so the walker can still detect the file type.
		if currentModule == nil {
			currentModule = &cmrModuleMetadata{
				sourceFileName: c.path,
				googNamespaces: map[string]bool{},
				strongRequires: map[string]bool{},
				weakRequires:   map[string]bool{},
				rootNode:       n,
			}
		}
	}

	// --- shouldTraverse equivalent (pre-order) ---
	if cmrIsValidCfgRoot(n) && !c.cmrIsInGlobalControlFlowScope(n) {
		c.controlFlowScopeDepth++
	}

	if currentModule != nil {
		c.cmrVisitNode(n, currentModule)
	}

	// Descend into children
	for _, ch := range n.Children() {
		c.walk(ch, n, currentModule)
	}

	// --- visit equivalent (post-order) ---
	if cmrIsValidCfgRoot(n) && !c.cmrIsInGlobalControlFlowScope(n) {
		c.controlFlowScopeDepth--
	}
	if currentModule != nil && n == currentModule.rootNode {
		c.templateParamNames = map[string]bool{}
	}
}

// ---------------------------------------------------------------------------
// cmrVisitNode — mirrors visitNode(NodeTraversal, Node, ModuleMetadata)
// ---------------------------------------------------------------------------

func (c *CheckMissingRequires) cmrVisitNode(n *rhino.Node, currentModule *cmrModuleMetadata) {
	info := n.JSDocInfo()
	if info != nil {
		c.cmrVisitJsDocInfo(n, currentModule, info)
	}

	// Check non-root GETPROP chains and NAME nodes for fully-qualified references.
	if n.Parent() != nil && n.Parent().Token() != rhino.GETPROP {
		if qn := cmrQualifiedNameFromNode(n); qn != nil {
			root := qn.getRoot()
			if root != "this" && root != "super" {
				c.cmrVisitQualifiedName(n, currentModule, qn, cmrStrengthCode)
			}
		}
	}

	// visitMaybeDeclaration for NAME nodes.
	if n.Token() == rhino.NAME && n.String_() != "" {
		c.cmrVisitMaybeDeclaration(n, currentModule)
	}

	// visitMaybeGoogModuleGet for CALL nodes at top-level control flow.
	if n.Token() == rhino.CALL && c.controlFlowScopeDepth == 0 {
		c.cmrVisitMaybeGoogModuleGet(n, currentModule)
	}
}

// ---------------------------------------------------------------------------
// cmrVisitJsDocInfo — mirrors visitJsDocInfo
// ---------------------------------------------------------------------------

func (c *CheckMissingRequires) cmrVisitJsDocInfo(n *rhino.Node, currentModule *cmrModuleMetadata, info *rhino.JSDocInfo) {
	// Collect template names before checking.
	for _, name := range info.TemplateNames() {
		c.templateParamNames[name] = true
	}

	if info.HasType() {
		c.cmrVisitJsDocExpr(n, currentModule, info.Type(), cmrStrengthWeakType)
	}
	for _, param := range info.ParameterNames() {
		if pt := info.ParameterType(param); pt != nil {
			c.cmrVisitJsDocExpr(n, currentModule, pt, cmrStrengthWeakType)
		}
	}
	if rt := info.ReturnType(); rt != nil {
		c.cmrVisitJsDocExpr(n, currentModule, rt, cmrStrengthWeakType)
	}
	if ep := info.EnumParameterType(); ep != nil {
		c.cmrVisitJsDocExpr(n, currentModule, ep, cmrStrengthWeakType)
	}
	if td := info.TypedefType(); td != nil {
		c.cmrVisitJsDocExpr(n, currentModule, td, cmrStrengthWeakType)
	}
	// @extends requires a strong require, not just a requireType.
	for _, ext := range info.BaseTypes() {
		c.cmrVisitJsDocExpr(n, currentModule, ext, cmrStrengthImplementsExtends)
	}
	// @implements requires a strong require.
	for _, impl := range info.Implementations() {
		c.cmrVisitJsDocExpr(n, currentModule, impl, cmrStrengthImplementsExtends)
	}
}

// cmrVisitJsDocExpr mirrors visitJsDocExpr.
// JSTypeExpression in our port only carries a TypeName string (Tier-0
// limitation). The Java side parses a full type AST and calls getAllTypeNodes().
// Here we parse the dotted name directly from the TypeName field.
//
// TODO: when JSTypeExpression is upgraded to a Node-rooted parse tree,
// iterate over all type nodes in the expression as Java does.
func (c *CheckMissingRequires) cmrVisitJsDocExpr(n *rhino.Node, currentModule *cmrModuleMetadata, expr *rhino.JSTypeExpression, strength cmrStrength) {
	if expr == nil || expr.TypeName == "" {
		return
	}
	// Strip template arguments, e.g. "Array<string>" → "Array".
	typeName := expr.TypeName
	if idx := strings.IndexByte(typeName, '<'); idx >= 0 {
		typeName = typeName[:idx]
	}
	// Strip leading "!" / "?" nullable markers.
	typeName = strings.TrimLeft(typeName, "!?")
	if typeName == "" {
		return
	}
	qn := cmrQualifiedNameOf(typeName)
	c.cmrVisitQualifiedName(n, currentModule, qn, strength)
}

// ---------------------------------------------------------------------------
// cmrVisitMaybeDeclaration — mirrors visitMaybeDeclaration
// ---------------------------------------------------------------------------

// cmrVisitMaybeDeclaration checks destructured goog.require imports for
// cases where the destructured property corresponds to its own Closure
// namespace (e.g., `const {Bar} = goog.require('foo')` when 'foo.Bar' is
// its own namespace).
//
// TODO: implement once Var / GoogRequire stubs are replaced with real
// scope-based implementations. The logic (Java lines 301–381):
//   1. Look up rootName in the current scope as a Var.
//   2. If it is a local destructured goog.require with a property,
//      check whether (namespace + "." + property) is its own namespace.
//   3. If so, fire INCORRECT_NAMESPACE_ALIAS_REQUIRE /
//      INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE.
func (c *CheckMissingRequires) cmrVisitMaybeDeclaration(_ *rhino.Node, currentModule *cmrModuleMetadata) {
	if !currentModule.isModule() {
		return
	}
	// Skip TypeScript-generated sources (mirrors Java isTypeScriptSource check).
	if cmrIsTypeScriptSource(c.path) {
		return
	}
	// Scope-based Var lookup not yet implemented; returns without firing.
}

// ---------------------------------------------------------------------------
// cmrVisitMaybeGoogModuleGet — mirrors visitMaybeGoogModuleGet
// ---------------------------------------------------------------------------

// cmrVisitMaybeGoogModuleGet checks goog.module.get calls at script-load time
// to ensure the required namespace has an acceptable goog.require.
func (c *CheckMissingRequires) cmrVisitMaybeGoogModuleGet(googModuleGet *rhino.Node, currentFile *cmrModuleMetadata) {
	if !cmrIsGoogModuleGetCall(googModuleGet) {
		return
	}
	kids := googModuleGet.Children()
	if len(kids) < 2 {
		return
	}
	secondChild := kids[1]
	if secondChild.Token() != rhino.STRINGLIT {
		return
	}
	importedNamespace := secondChild.String_()
	requiredFile := c.moduleByNamespace[importedNamespace]
	if requiredFile == nil {
		return
	}
	qn := cmrQualifiedNameOf(importedNamespace)
	if !cmrHasAcceptableRequire(currentFile, qn, requiredFile, cmrStrengthWeakType) {
		c.cmrReport(googModuleGet, cmrCodeMissingRequireForGoogModuleGet, importedNamespace, "")
	}
}

// ---------------------------------------------------------------------------
// cmrVisitQualifiedName — mirrors visitQualifiedName
// ---------------------------------------------------------------------------

func (c *CheckMissingRequires) cmrVisitQualifiedName(n *rhino.Node, currentFile *cmrModuleMetadata, qualifiedName *cmrQualifiedName, referenceStrength cmrStrength) {
	rootName := qualifiedName.getRoot()

	if qualifiedName.isSimple() {
		if c.templateParamNames[rootName] {
			return
		}
	}

	if cmrIsDeclarationLValue(n) {
		return
	}

	// TODO: look up rootName in the current scope (requires Var stub to be
	// populated). Java branches on whether the var is local to call
	// checkMissingRequireThroughShortName vs checkMissingRequireThroughFullyQualifiedName.
	// For now, always take the fully-qualified path (handles global references).
	c.cmrCheckMissingRequireThroughFullyQualifiedName(n, currentFile, qualifiedName, referenceStrength)
}

// ---------------------------------------------------------------------------
// cmrCheckMissingRequireThroughShortName — mirrors checkMissingRequireThroughShortName
// ---------------------------------------------------------------------------

// cmrCheckMissingRequireThroughShortName is called when the root of the
// qualified name is a local variable that may be an alias for an imported namespace.
//
// TODO: implement once Var / GoogRequire / scope stubs are real. Java logic:
//   1. Skip TypeScript sources.
//   2. Determine inModule / inGoogScope.
//   3. For simple names (no dots), return — already validated in visitMaybeDeclaration.
//   4. Dispatch to checkMissingRequireThroughShortNameInGoogScope or
//      checkMissingRequireThroughShortNameInModule.
func (c *CheckMissingRequires) cmrCheckMissingRequireThroughShortName(
	n *rhino.Node,
	localVar *cmrVar,
	currentFile *cmrModuleMetadata,
	qualifiedName *cmrQualifiedName,
	referenceStrength cmrStrength,
) {
	// Skip TypeScript-generated sources (mirrors Java isTypeScriptSource check).
	if cmrIsTypeScriptSource(c.path) {
		return
	}
	// Scope-based Var / GoogRequire resolution not yet implemented.
	_ = n
	_ = localVar
	_ = currentFile
	_ = qualifiedName
	_ = referenceStrength
}

// ---------------------------------------------------------------------------
// cmrCheckMissingRequireThroughShortNameInGoogScope
// ---------------------------------------------------------------------------

// cmrCheckMissingRequireThroughShortNameInGoogScope handles the goog.scope case.
//
// TODO: implement once Var stubs are real. Java logic handles two cases:
//  1. Initial value is a qualified name (legacy namespace + optional extra props).
//  2. Initial value is a goog.module.get call (+ optional extra props).
func (c *CheckMissingRequires) cmrCheckMissingRequireThroughShortNameInGoogScope(
	n *rhino.Node,
	localVar *cmrVar,
	currentFile *cmrModuleMetadata,
	qualifiedName *cmrQualifiedName,
	referenceStrength cmrStrength,
) {
	_ = n
	_ = localVar
	_ = currentFile
	_ = qualifiedName
	_ = referenceStrength
}

// ---------------------------------------------------------------------------
// cmrCheckMissingRequireThroughShortNameInModule
// ---------------------------------------------------------------------------

// cmrCheckMissingRequireThroughShortNameInModule handles the goog.module case.
//
// TODO: implement once GoogRequire stubs are real.
func (c *CheckMissingRequires) cmrCheckMissingRequireThroughShortNameInModule(
	n *rhino.Node,
	require *cmrGoogRequire,
	currentFile *cmrModuleMetadata,
	qualifiedName *cmrQualifiedName,
	referenceStrength cmrStrength,
) {
	_ = n
	_ = require
	_ = currentFile
	_ = qualifiedName
	_ = referenceStrength
}

// ---------------------------------------------------------------------------
// cmrCheckNormalizedShortNameImport
// ---------------------------------------------------------------------------

// cmrCheckNormalizedShortNameImport is the shared inner logic for the two
// "through short name" methods above. It resolves the longest-prefix match
// of `normalizedQualifiedName` against known namespaces and fires the
// appropriate diagnostic.
//
// TODO: call this from cmrCheckMissingRequireThroughShortNameInGoogScope and
// cmrCheckMissingRequireThroughShortNameInModule once those stubs are real.
func (c *CheckMissingRequires) cmrCheckNormalizedShortNameImport(
	n *rhino.Node,
	currentFile *cmrModuleMetadata,
	originalImport *cmrModuleMetadata,
	originalQualifiedName *cmrQualifiedName,
	normalizedQualifiedName *cmrQualifiedName,
	referenceStrength cmrStrength,
) {
	for subName := normalizedQualifiedName; subName != nil; subName = subName.getOwner() {
		namespace := subName.join()
		if cmrIsAllowedNamespace(currentFile, namespace) {
			return
		}
		requiredFile := c.moduleByNamespace[namespace]
		if requiredFile == nil {
			continue
		}

		if !cmrHasAcceptableRequire(currentFile, subName, requiredFile, referenceStrength) ||
			originalImport == nil ||
			originalImport != requiredFile {

			var code string
			if currentFile.isModule() {
				if referenceStrength.isStrong() {
					code = cmrCodeIndirectNSRefRequire
				} else {
					code = cmrCodeIndirectNSRefRequireType
				}
			} else {
				if referenceStrength.isStrong() {
					code = cmrCodeMissingRequireInGoogScope
				} else {
					code = cmrCodeMissingRequireTypeInGoogScope
				}
			}

			// Compute the error node: walk up the GETPROP chain by the
			// difference in component count.
			diff := normalizedQualifiedName.getComponentCount() - subName.getComponentCount()
			errorNode := n
			if n.Token() == rhino.GETPROP || n.Token() == rhino.NAME {
				for i := 0; i < diff && errorNode.Token() == rhino.GETPROP; i++ {
					errorNode = errorNode.FirstChild()
				}
				c.cmrReport(errorNode, code, namespace, "")
			} else {
				// JSDoc reference — trim the original qualified name by diff levels.
				errorQName := originalQualifiedName
				for i := 0; i < diff && errorQName != nil; i++ {
					errorQName = errorQName.getOwner()
				}
				if errorQName != nil {
					c.cmrReport(n, code, errorQName.join(), "")
				}
			}
		}
		return
	}
}

// ---------------------------------------------------------------------------
// cmrCheckMissingRequireThroughFullyQualifiedName
// ---------------------------------------------------------------------------

// cmrCheckMissingRequireThroughFullyQualifiedName mirrors the Java method of
// the same name. It looks for the longest prefix of `qualifiedName` that
// matches a known Closure namespace and fires the appropriate diagnostic.
func (c *CheckMissingRequires) cmrCheckMissingRequireThroughFullyQualifiedName(
	n *rhino.Node,
	currentFile *cmrModuleMetadata,
	qualifiedName *cmrQualifiedName,
	referenceStrength cmrStrength,
) {
	for subName := qualifiedName; subName != nil; subName = subName.getOwner() {
		namespace := subName.join()
		if cmrIsAllowedNamespace(currentFile, namespace) {
			return
		}
		requiredFile := c.moduleByNamespace[namespace]
		if requiredFile == nil {
			continue
		}
		fileName := requiredFile.getRootSourceFileName()

		var code string
		if currentFile.isModule() {
			// In module files, any FQN reference (even with a require) is disallowed.
			if referenceStrength.isStrong() {
				code = cmrCodeMissingRequire
			} else {
				code = cmrCodeMissingRequireType
			}
		} else if !cmrHasAcceptableRequire(currentFile, subName, requiredFile, referenceStrength) {
			// In non-module files, fire only when there is no acceptable require.
			if referenceStrength.isStrong() {
				code = cmrCodeMissingRequireInProvidesFile
			} else {
				code = cmrCodeMissingRequireTypeInProvidesFile
			}
		} else if !referenceStrength.isTypeOnly() && requiredFile.isNonLegacyGoogModule() {
			// Has a valid require but the name won't be defined at runtime
			// (non-legacy goog.module).
			code = cmrCodeNonLegacyGoogModuleReference
		} else {
			return
		}

		c.cmrReport(n, code, namespace, fileName)
		return
	}
}

// ---------------------------------------------------------------------------
// cmrNormalizeQualifiedNamePlusImport — mirrors normalizeQualifiedNamePlusImport
// ---------------------------------------------------------------------------

// cmrNormalizeQualifiedNamePlusImport builds the fully-qualified name from a
// local alias + its GoogRequire origin. Mirrors the Java helper.
func cmrNormalizeQualifiedNamePlusImport(qualifiedName *cmrQualifiedName, imported *cmrGoogRequire) *cmrQualifiedName {
	newQName := cmrQualifiedNameOf(imported.namespace)
	if imported.hasProperty {
		newQName = newQName.getprop(imported.property)
	}
	return cmrSwapRootOfQualifiedName(qualifiedName, newQName)
}

// ---------------------------------------------------------------------------
// cmrSwapRootOfQualifiedName — mirrors swapRootOfQualifiedName
// ---------------------------------------------------------------------------

// cmrSwapRootOfQualifiedName replaces the root component of `original` with
// `newRoot` (which may itself be a dotted name). Mirrors the Java static helper.
//
// Example: original="oldRoot.a.b", newRoot="x.y" → "x.y.a.b"
func cmrSwapRootOfQualifiedName(original, newRoot *cmrQualifiedName) *cmrQualifiedName {
	result := newRoot
	if original == nil || len(original.components) <= 1 {
		return result
	}
	// Skip the root component of original and append the rest to newRoot.
	for _, comp := range original.components[1:] {
		result = result.getprop(comp)
	}
	return result
}

// ---------------------------------------------------------------------------
// Utility helpers — mirrors of Java static methods
// ---------------------------------------------------------------------------

// cmrIsAllowedNamespace mirrors CheckMissingRequires.isAllowedNamespace.
func cmrIsAllowedNamespace(currentFile *cmrModuleMetadata, namespace string) bool {
	if namespace == "goog.module" {
		// Special-cased: goog.module is both a primitive keyword and a
		// Closure Library namespace; references to the namespace itself are OK.
		return true
	}
	return currentFile.containsNamespace(namespace)
}

// cmrHasAcceptableRequire mirrors CheckMissingRequires.hasAcceptableRequire.
// It checks whether `rdep` has a require (strong or weak, depending on
// `referenceStrength`) that covers any prefix of `namespace` declared by `dep`.
func cmrHasAcceptableRequire(rdep *cmrModuleMetadata, namespace *cmrQualifiedName, dep *cmrModuleMetadata, referenceStrength cmrStrength) bool {
	if rdep == nil || dep == nil {
		return false
	}
	for parent := namespace; parent != nil; parent = parent.getOwner() {
		ns := parent.join()
		if !dep.containsNamespace(ns) {
			continue
		}
		if rdep.containsStrongRequire(ns) {
			return true
		}
		if !referenceStrength.isStrong() && rdep.containsWeakRequire(ns) {
			return true
		}
	}
	return false
}

// cmrIsValidCfgRoot reports whether `n` is a valid control-flow graph root.
// Mirrors NodeUtil.isValidCfgRoot: DEFINITE_CFG_ROOTS = {FUNCTION, SCRIPT,
// MODULE_BODY, ROOT}. (Class static blocks are excluded here as they are not
// produced by our parser tier.)
func cmrIsValidCfgRoot(n *rhino.Node) bool {
	switch n.Token() {
	case rhino.SCRIPT, rhino.FUNCTION, rhino.MODULE_BODY, rhino.ROOT:
		return true
	}
	return false
}

// cmrIsInGlobalControlFlowScope mirrors CheckMissingRequires.isInGlobalControlFlowScope.
// Returns true for nodes that are within the global control-flow scope (SCRIPT, ROOT,
// IIFEs, goog.scope bodies) and false for nodes that introduce a new scope (MODULE_BODY,
// non-IIFE FUNCTIONs). Matches the Java switch exactly.
func (c *CheckMissingRequires) cmrIsInGlobalControlFlowScope(n *rhino.Node) bool {
	switch n.Token() {
	case rhino.SCRIPT, rhino.ROOT:
		return true
	case rhino.MODULE_BODY:
		// MODULE_BODY creates a new cfg root but is not in the global scope.
		return false
	case rhino.FUNCTION:
		if c.controlFlowScopeDepth > 0 {
			return false
		}
		// IIFEs and goog.scope callees are in the global control-flow scope.
		body := cmrGetFunctionBody(n)
		if body != nil && cmrIsGoogScopeBody(body) {
			return true
		}
		return cmrIsIIFE(n)
	}
	return false
}

// cmrGetFunctionBody returns the BLOCK child of a FUNCTION node (the body),
// mirroring NodeUtil.getFunctionBody.
func cmrGetFunctionBody(fn *rhino.Node) *rhino.Node {
	if fn == nil || fn.Token() != rhino.FUNCTION {
		return nil
	}
	// Java: FUNCTION children are [NAME, PARAM_LIST, BLOCK].
	kids := fn.Children()
	for _, ch := range kids {
		if ch.Token() == rhino.BLOCK {
			return ch
		}
	}
	return nil
}

// cmrIsGoogScopeBody mirrors the Java static helper of the same name.
// It returns true when `body` is the BLOCK of a function passed directly to goog.scope(…).
func cmrIsGoogScopeBody(body *rhino.Node) bool {
	if body == nil || body.Token() != rhino.BLOCK {
		return false
	}
	fn := body.Parent()
	if fn == nil || fn.Token() != rhino.FUNCTION {
		return false
	}
	call := fn.Parent()
	if call == nil || call.Token() != rhino.CALL {
		return false
	}
	callee := call.FirstChild()
	return cmrMatchesQualifiedName(callee, cmrGoogScope)
}

// cmrIsIIFE mirrors the Java static isIIFE helper in CheckMissingRequires.
func cmrIsIIFE(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.FUNCTION {
		return false
	}
	parent := n.Parent()
	if parent == nil || parent.Token() != rhino.CALL {
		return false
	}
	return parent.FirstChild() == n
}

// cmrIsGoogModuleGetCall mirrors NodeUtil.isGoogModuleGetCall.
// It returns true when `n` is a CALL node of the form goog.module.get("…").
func cmrIsGoogModuleGetCall(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.CALL {
		return false
	}
	kids := n.Children()
	if len(kids) < 2 {
		return false
	}
	callee := kids[0]
	qn := cmrQualifiedNameFromNode(callee)
	if qn == nil || qn.join() != "goog.module.get" {
		return false
	}
	// Second child must be a string literal.
	return kids[1].Token() == rhino.STRINGLIT
}

// cmrTocci converts a cmrModuleMetadataMap into a ModuleMetadataMap so
// that CheckClosureImports can use the same single-file module graph built
// by cmrGatherModuleMetadata. Only the googNamespaces set and module type
// are carried across (CheckClosureImports only queries isGoogModule/isModule
// and getModulesByGoogNamespace).
func cmrTocci(m *cmrModuleMetadataMap) *ModuleMetadataMap {
	out := newModuleMetadataMap()
	if m == nil {
		return out
	}
	for ns, md := range m.byNamespace {
		var mt cciModuleType
		if md.isModuleFlag {
			mt = cciModuleTypeGoogModule
		} else {
			mt = cciModuleTypeScript
		}
		gns := map[string]struct{}{ns: {}}
		out.byGoogNamespace[ns] = &ModuleMetadata{
			googNamespaces: gns,
			moduleType:     mt,
			rootNode:       md.rootNode,
		}
	}
	return out
}

// cmrIsTypeScriptSource mirrors StaticSourceFile.isTypeScriptSource().
// Returns true for .closure.js and .tsx.cl.js files, which are
// tsickle-generated TypeScript transpilations.
func cmrIsTypeScriptSource(path string) bool {
	return strings.HasSuffix(path, ".closure.js") || strings.HasSuffix(path, ".tsx.cl.js")
}

// cmrIsDeclarationLValue mirrors NodeUtil.isDeclarationLValue.
// Returns true when `n` is the left-hand side of a VAR/LET/CONST declaration
// or the name being declared in a FUNCTION/CLASS statement.
func cmrIsDeclarationLValue(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	parent := n.Parent()
	if parent == nil {
		return false
	}
	switch parent.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		return true
	case rhino.FUNCTION, rhino.CLASS:
		// The NAME child at position 0 is the declared name.
		kids := parent.Children()
		if len(kids) > 0 && kids[0] == n {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// cmrReport — diagnostic emission helper
// ---------------------------------------------------------------------------

// cmrReport emits a diagnostic. The `code` constant determines the level
// (WARNING or ERROR) consistent with the Java DiagnosticType registration.
// `arg0` and `arg1` are the format arguments; `arg1` is ignored when empty.
func (c *CheckMissingRequires) cmrReport(n *rhino.Node, code, arg0, arg1 string) {
	if n == nil {
		return
	}
	level := diag.Warning
	if code == cmrCodeNonLegacyGoogModuleReference {
		level = diag.Error
	}
	msg := cmrFormatMessage(code, arg0, arg1)
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    n.Lineno(),
		Col:     n.Charno(),
		Code:    code,
		Level:   level,
		Message: msg,
	})
}

// cmrFormatMessage returns a human-readable message for the given code,
// substituting arg0 and arg1 into the Java format template.
func cmrFormatMessage(code, arg0, arg1 string) string {
	switch code {
	case cmrCodeMissingRequire:
		return "'" + arg0 + "' references a fully qualified namespace, which is disallowed by the style guide.\n" +
			"Please add a goog.require, assign or destructure it into an alias, and use the alias instead.\n" +
			"Namespace defined in " + arg1 + "."
	case cmrCodeMissingRequireType:
		return "'" + arg0 + "' references a fully qualified namespace, which is disallowed by the style guide.\n" +
			"Please add a goog.requireType, assign or destructure it into an alias, and use the alias instead.\n" +
			"Namespace defined in " + arg1 + "."
	case cmrCodeIncorrectNSAliasRequire:
		return "'" + arg0 + "' is its own namespace.\nPlease add a separate goog.require and use that alias instead."
	case cmrCodeIncorrectNSAliasRequireType:
		return "'" + arg0 + "' is its own namespace.\nPlease add a separate goog.requireType and use that alias instead."
	case cmrCodeIndirectNSRefRequire:
		return "'" + arg0 + "' should have its own goog.require.\nPlease add a separate goog.require and use that alias instead."
	case cmrCodeIndirectNSRefRequireType:
		return "'" + arg0 + "' should have its own goog.requireType.\nPlease add a separate goog.requireType and use that alias instead."
	case cmrCodeMissingRequireInProvidesFile:
		return "'" + arg0 + "' references a namespace which was not required by this file.\n" +
			"Please add a goog.require.\nNamespace defined in " + arg1 + "."
	case cmrCodeMissingRequireTypeInProvidesFile:
		return "'" + arg0 + "' references a namespace which was not required by this file.\n" +
			"Please add a goog.requireType.\nNamespace defined in " + arg1 + "."
	case cmrCodeNonLegacyGoogModuleReference:
		return "'" + arg0 + "' references the name of a module without goog.declareLegacyNamespace(), which " +
			"is not actually defined. Use goog.module.get() instead."
	case cmrCodeMissingRequireInGoogScope:
		return "'" + arg0 + "' is its own namespace. Please add a goog.require and reference that alias."
	case cmrCodeMissingRequireTypeInGoogScope:
		return "'" + arg0 + "' is its own namespace. Please add a goog.requireType and reference that alias."
	case cmrCodeMissingRequireForGoogModuleGet:
		return "'" + arg0 + "' references a namespace which was not required by this file.\nPlease add a goog.require."
	}
	return arg0
}
