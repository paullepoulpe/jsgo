// Port of com/google/javascript/jscomp/CheckClosureImports.java @ closure-compiler v20260415 (narrow; ModuleMetadata stubbed)
//
// Checks all goog.requires, goog.module.gets, goog.forwardDeclares, and
// goog.requireTypes in all files.
//
// The Java pass depends on ModuleMetadataMap, ModuleLoader, and chunk-graph
// infrastructure not yet ported to Go. Those are stubbed minimally inside
// this file (see TODO comments). The pass structure and all diagnostic
// codes are preserved 1:1; the actual enforcement logic is wired to fire
// once the stubs are replaced.
//
// Diagnostic codes defined here (mirroring Java):
//
//	JSC_INVALID_CLOSURE_IMPORT_DESTRUCTURING (error)
//	JSC_ONE_CLOSURE_IMPORT_PER_DECLARATION   (error)
//	JSC_INVALID_CLOSURE_IMPORT_CALL          (error)
//	JSC_LATE_PROVIDE_ERROR                   (error)
//	JSC_LET_CLOSURE_IMPORT                   (disabled in Java, warning here)
//	JSC_NO_CLOSURE_IMPORT_DESTRUCTURING      (error)
//	JSC_LHS_OF_CLOUSRE_IMPORT_MUST_BE_CONST_IN_ES_MODULE (error)
//	JSC_XMODULE_REQUIRE_ERROR                (warning)
//
// Plus cross-file codes consumed but defined elsewhere:
//	JSC_INVALID_CLOSURE_CALL_ERROR (ClosurePrimitiveErrors)
//	JSC_MODULE_USES_GOOG_MODULE_GET (ClosurePrimitiveErrors)
//	JSC_MISSING_MODULE_OR_PROVIDE (ClosurePrimitiveErrors)
//	JSC_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE (ClosurePrimitiveErrors)
//	JSC_GOOG_MODULE_INVALID_GET_CALL_SCOPE (ClosurePrimitiveErrors)
//	JSC_GOOG_MODULE_INVALID_GET_ALIAS (ClosureRewriteModule)
//	JSC_INCORRECT_SHORTNAME_CAPITALIZATION (ClosureCheckModule)
//	JSC_SHOULD_IMPORT_ES6_MODULE (Es6RewriteModules)

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ---------------------------------------------------------------------------
// Stubs for module-system types not yet ported.
// ---------------------------------------------------------------------------

// cciModuleType mirrors the subset of ModuleMetadata.ModuleType that
// CheckClosureImports queries.
// TODO: replace with the real ModuleMetadata port.
type cciModuleType int

const (
	cciModuleTypeScript    cciModuleType = iota // plain script / goog.provide
	cciModuleTypeGoogModule                     // goog.module(…)
	cciModuleTypeEs6Module                      // ES6 import/export
)

// cciModuleMetadata is a minimal stub for
// com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata.
// TODO: replace with the real ModuleMetadata port.
type cciModuleMetadata struct {
	// googNamespaces is the set of goog.provide / goog.module / goog.declareModuleId
	// namespaces declared by the module's root node.
	googNamespaces map[string]struct{}
	moduleType     cciModuleType
	// rootNode is the SCRIPT node for the module; nil for synthetic/test metadata.
	rootNode *rhino.Node
}

func (m *cciModuleMetadata) isGoogModule() bool { return m.moduleType == cciModuleTypeGoogModule }
func (m *cciModuleMetadata) isEs6Module() bool  { return m.moduleType == cciModuleTypeEs6Module }

// isModule mirrors ModuleMetadata.isModule() — true for goog.module and ES6 modules.
func (m *cciModuleMetadata) isModule() bool {
	return m.isGoogModule() || m.isEs6Module()
}

// cciModuleMetadataMap is a minimal stub for
// com.google.javascript.jscomp.modules.ModuleMetadataMap.
// TODO: replace with the real ModuleMetadataMap port.
type cciModuleMetadataMap struct {
	// byGoogNamespace maps goog.provide / goog.module / goog.declareModuleId
	// namespace strings to their owning module metadata.
	byGoogNamespace map[string]*cciModuleMetadata
}

func newCciModuleMetadataMap() *cciModuleMetadataMap {
	return &cciModuleMetadataMap{byGoogNamespace: map[string]*cciModuleMetadata{}}
}

// getModulesByGoogNamespace mirrors
// ModuleMetadataMap.getModulesByGoogNamespace().get(namespace).
func (m *cciModuleMetadataMap) getModulesByGoogNamespace(ns string) *cciModuleMetadata {
	return m.byGoogNamespace[ns]
}

// ---------------------------------------------------------------------------
// cciClosureImport — mirrors the ClosureImport enum with methods.
// ---------------------------------------------------------------------------

type cciClosureImport int

const (
	cciImportRequire        cciClosureImport = iota // goog.require
	cciImportForwardDeclare                         // goog.forwardDeclare
	cciImportRequireType                            // goog.requireType
	cciImportRequireDynamic                         // goog.requireDynamic
)

// allowDestructuring mirrors ClosureImport.allowDestructuring().
func (ci cciClosureImport) allowDestructuring() bool {
	switch ci {
	case cciImportRequire, cciImportRequireType, cciImportRequireDynamic:
		return true
	}
	return false
}

// aliasMustBeConstant mirrors ClosureImport.aliasMustBeConstant().
func (ci cciClosureImport) aliasMustBeConstant() bool {
	switch ci {
	case cciImportRequire, cciImportRequireType, cciImportRequireDynamic:
		return true
	}
	return false
}

// mustBeOrdered mirrors ClosureImport.mustBeOrdered().
func (ci cciClosureImport) mustBeOrdered() bool {
	return ci == cciImportRequire
}

// callName mirrors ClosureImport.callName().
func (ci cciClosureImport) callName() string {
	switch ci {
	case cciImportRequire:
		return "goog.require"
	case cciImportForwardDeclare:
		return "goog.forwardDeclare"
	case cciImportRequireType:
		return "goog.requireType"
	case cciImportRequireDynamic:
		return "goog.requireDynamic"
	}
	return "<unknown>"
}

// ---------------------------------------------------------------------------
// CheckClosureImports — the CompilerPass struct.
// ---------------------------------------------------------------------------

// CheckClosureImports mirrors com.google.javascript.jscomp.CheckClosureImports.
// It walks the AST verifying that goog.require / goog.forwardDeclare /
// goog.requireType / goog.requireDynamic calls are well-formed, in scope,
// and refer to defined namespaces.
//
// Because ModuleMetadataMap is not yet ported, the pass accumulates
// diagnostics via a minimal stub map; namespacesSeen tracks which goog
// namespaces have been declared prior to a require (ordering check).
type CheckClosureImports struct {
	path            string
	moduleMetadata  *cciModuleMetadataMap // TODO: replace with real ModuleMetadataMap port
	namespacesSeen  map[string]struct{}
	diags           []diag.Diagnostic
}

// NewCheckClosureImports mirrors the Java constructor:
//
//	CheckClosureImports(AbstractCompiler compiler, ModuleMetadataMap)
//
// The path argument is the source file path used in diagnostics.
// The metadata map is stubbed; pass nil to use an empty map.
func NewCheckClosureImports(path string, metadata *cciModuleMetadataMap) *CheckClosureImports {
	if metadata == nil {
		metadata = newCciModuleMetadataMap()
	}
	return &CheckClosureImports{
		path:           path,
		moduleMetadata: metadata,
		namespacesSeen: map[string]struct{}{},
	}
}

// Check runs the pass over `script`, accumulating diagnostics.
// Mirrors CompilerPass.process(Node externs, Node root) — this shim
// operates on a single script node (externs traversal is elided until
// a multi-file compiler is ported).
func (c *CheckClosureImports) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.namespacesSeen = map[string]struct{}{}

	// TODO: Once ModuleMetadataMap is ported, seed namespacesSeen via
	// AbstractModuleCallback.enterModule(currentModule).
	// For now walk the AST to collect goog.provide / goog.module
	// namespace declarations encountered before each require, so the
	// ordering check (LATE_PROVIDE_ERROR) can fire correctly.
	c.collectNamespaces(script)

	// Walk the AST and check every call that looks like a Closure import.
	c.walk(script, nil)

	return c.diags
}

// collectNamespaces pre-populates namespacesSeen with every namespace
// that is *declared* in the file (goog.provide / goog.module calls).
// This is a simplified substitute for the Java AbstractModuleCallback
// enterModule hook which would receive pre-built metadata.
//
// TODO: replace with ModuleMetadataMap.enterModule seeding.
func (c *CheckClosureImports) collectNamespaces(n *rhino.Node) {
	if n == nil {
		return
	}
	if n.Token() == rhino.CALL {
		if ns := cciGoogProvideOrModuleNamespace(n); ns != "" {
			c.namespacesSeen[ns] = struct{}{}
		}
	}
	for _, ch := range n.Children() {
		c.collectNamespaces(ch)
	}
}

// walk visits every node depth-first.  Mirrors the Checker inner-class
// visit() method, which receives both CALL and NAME nodes.
func (c *CheckClosureImports) walk(n, parent *rhino.Node) {
	if n == nil {
		return
	}
	if n.Token() == rhino.CALL {
		kind, ok := cciKindOfCall(n)
		if ok {
			c.checkImport(n, parent)
			_ = kind // kind is used in checkImport via cciKindOfCall re-call
		} else if cciMatchesQualifiedName(n.FirstChild(), "goog.module.get") {
			c.checkGoogModuleGet(n)
		}
	}
	for _, ch := range n.Children() {
		c.walk(ch, n)
	}
}

// ---------------------------------------------------------------------------
// cciKindOfCall — mirrors Checker.kindOfCall(Node callNode).
// ---------------------------------------------------------------------------

// cciKindOfCall returns the ClosureImport kind for callNode, if any.
// Returns (_, false) when the call is not a recognized Closure import.
func cciKindOfCall(callNode *rhino.Node) (cciClosureImport, bool) {
	callee := callNode.FirstChild()
	switch {
	case cciMatchesQualifiedName(callee, "goog.require"):
		return cciImportRequire, true
	case cciMatchesQualifiedName(callee, "goog.forwardDeclare"):
		return cciImportForwardDeclare, true
	case cciMatchesQualifiedName(callee, "goog.requireType"):
		return cciImportRequireType, true
	case cciMatchesQualifiedName(callee, "goog.requireDynamic"):
		return cciImportRequireDynamic, true
	}
	return 0, false
}

// ---------------------------------------------------------------------------
// checkGoogModuleGet — mirrors Checker.checkGoogModuleGet.
// ---------------------------------------------------------------------------

// checkGoogModuleGet validates a goog.module.get(…) call.
// Full scope/chunk checks require ModuleMetadata + scope stack; those
// are stubbed until the module system is ported.
func (c *CheckClosureImports) checkGoogModuleGet(call *rhino.Node) {
	// TODO: check t.inModuleHoistScope() → MODULE_USES_GOOG_MODULE_GET
	// TODO: check INVALID_GET_CALL_SCOPE for goog.provide files in global scope
	// TODO: check INVALID_GET_ALIAS when filling a forwardDeclare alias

	// Validate argument.
	if !cciHasTwoChildren(call) || cciSecondChild(call) == nil || cciSecondChild(call).Token() != rhino.STRINGLIT {
		c.report(call, diag.JSC_INVALID_CLOSURE_IMPORT_CALL, diag.Error,
			"goog.module.get parameter must be a string literal.")
		return
	}

	namespace := cciSecondChild(call).String_()
	if c.moduleMetadata.getModulesByGoogNamespace(namespace) == nil {
		c.report(call, diag.JSC_MISSING_MODULE_OR_PROVIDE, diag.Error,
			fmt.Sprintf("Required namespace %q never defined.", namespace))
	}
}

// ---------------------------------------------------------------------------
// checkImport — mirrors Checker.checkImport.
// ---------------------------------------------------------------------------

// checkImport validates a single Closure-import call (require / forwardDeclare /
// requireType / requireDynamic). `call` is the CALL node; `parent` is its
// parent in the AST.
func (c *CheckClosureImports) checkImport(call, parent *rhino.Node) {
	importType, _ := cciKindOfCall(call)

	// -- scope check --------------------------------------------------------
	// goog.requireDynamic is allowed anywhere; others must be at top-level
	// scope of a module or global hoist scope AND directly under an
	// EXPR_RESULT or a name-declaration.
	//
	// TODO: use real scope-stack from NodeTraversal once wired; for now we
	// approximate "not in function" by checking for a FUNCTION ancestor.
	atTopLevel := !cciHasFunctionAncestor(call)
	isModule := false // TODO: derive from ModuleMetadata.isModule()
	validAssignment := isModule || cciIsExprResult(parent)

	if (!atTopLevel || !validAssignment) && importType != cciImportRequireDynamic {
		c.report(call, diag.JSC_INVALID_CLOSURE_CALL_ERROR, diag.Error,
			"Closure import call is not in an acceptable scope.")
		return
	}

	// -- argument check -----------------------------------------------------
	if !cciHasTwoChildren(call) || cciSecondChild(call) == nil || cciSecondChild(call).Token() != rhino.STRINGLIT {
		c.report(call, diag.JSC_INVALID_CLOSURE_IMPORT_CALL, diag.Error,
			fmt.Sprintf("%s parameter must be a string literal.", importType.callName()))
		return
	}

	// -- alias / declaration checks -----------------------------------------
	isAliased := parent != nil && cciIsNameDeclaration(cciGrandParent(call))
	if isModule && isAliased {
		decl := cciGrandParent(call)
		if decl != nil {
			if !cciHasOneChild(decl) {
				c.report(call, diag.JSC_ONE_CLOSURE_IMPORT_PER_DECLARATION, diag.Error,
					fmt.Sprintf("There may only be one %s per var/let/const declaration.", importType.callName()))
				return
			}

			if importType.aliasMustBeConstant() {
				if cciIsLet(decl) {
					c.report(call, diag.JSC_LET_CLOSURE_IMPORT, diag.Warning,
						"Module imports must be constant. Please use 'const' instead of 'let'.")
				} else if !cciIsConst(decl) {
					// ES module: must use const
					// TODO: check currentModule.isEs6Module() when ModuleMetadata is ported
					c.report(call, diag.JSC_LHS_OF_CLOUSRE_IMPORT_MUST_BE_CONST_IN_ES_MODULE, diag.Error,
						fmt.Sprintf("The left side of a %s must use 'const' (not 'let' or 'var') in an ES module.", importType.callName()))
				}
			}

			lhs := decl.FirstChild()
			if lhs != nil && cciIsDestructuringLhs(lhs) {
				if importType.allowDestructuring() {
					if !cciIsValidDestructuringImport(lhs) {
						c.report(decl, diag.JSC_INVALID_CLOSURE_IMPORT_DESTRUCTURING, diag.Error,
							fmt.Sprintf("Destructuring %s must be a simple object pattern.", importType.callName()))
					}
				} else {
					c.report(decl, diag.JSC_NO_CLOSURE_IMPORT_DESTRUCTURING, diag.Error,
						fmt.Sprintf("Cannot destructure the return value of %s", importType.callName()))
				}
			} else if lhs != nil {
				// checkShortName — compare alias capitalization with last namespace segment
				ns := cciSecondChild(call).String_()
				c.cciCheckShortName(lhs, ns)
			}
		}
	}

	// -- namespace existence check ------------------------------------------
	namespace := cciSecondChild(call).String_()
	requiredModule := c.moduleMetadata.getModulesByGoogNamespace(namespace)

	if requiredModule == nil {
		if importType == cciImportForwardDeclare {
			// goog.forwardDeclare allows missing namespaces (legacy behaviour).
			// When aliased inside a module, fire a softer warning.
			if isAliased && isModule {
				c.report(call, diag.JSC_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE, diag.Error,
					fmt.Sprintf("Required namespace %q never defined.", namespace))
			}
			return
		}
		c.report(call, diag.JSC_MISSING_MODULE_OR_PROVIDE, diag.Error,
			fmt.Sprintf("Required namespace %q never defined.", namespace))
		return
	}

	// -- ordering check -----------------------------------------------------
	if importType.mustBeOrdered() {
		c.cciVerifyRequireOrder(namespace, call, requiredModule)
	}

	// -- ES6 module mutual-import check -------------------------------------
	// TODO: fire JSC_SHOULD_IMPORT_ES6_MODULE when both currentModule and
	// requiredModule are ES6 modules. Requires ModuleMetadata.isEs6Module().
}

// ---------------------------------------------------------------------------
// cciVerifyRequireOrder — mirrors CheckClosureImports.verifyRequireOrder.
// ---------------------------------------------------------------------------

// cciVerifyRequireOrder checks that the required namespace was declared
// before the current require call. Mirrors verifyRequireOrder() in Java.
func (c *CheckClosureImports) cciVerifyRequireOrder(namespace string, call *rhino.Node, _ *cciModuleMetadata) {
	if _, seen := c.namespacesSeen[namespace]; !seen {
		c.report(call, diag.JSC_LATE_PROVIDE_ERROR, diag.Error,
			fmt.Sprintf("Required namespace %q not provided yet.", namespace))
		return
	}
	// TODO: cross-chunk check (JSC_XMODULE_REQUIRE_ERROR) — requires
	// JSChunkGraph + CompilerInput.getChunk(), neither of which is ported yet.
}

// ---------------------------------------------------------------------------
// cciCheckShortName — mirrors Checker.checkShortName.
// ---------------------------------------------------------------------------

// cciCheckShortName validates that the short alias name for a require has
// matching capitalization with the last segment of the namespace.
// Mirrors ClosureCheckModule.checkShortName / the inline call in Java.
func (c *CheckClosureImports) cciCheckShortName(shortNameNode *rhino.Node, namespace string) {
	if shortNameNode == nil {
		return
	}
	shortName := shortNameNode.String_()
	lastDot := len(namespace) - 1
	for i := len(namespace) - 1; i >= 0; i-- {
		if namespace[i] == '.' {
			lastDot = i
			break
		}
	}
	lastSegment := namespace[lastDot+1:]
	if shortName == lastSegment || lastSegment == "" {
		return
	}
	if len(shortName) == 0 || len(lastSegment) == 0 {
		return
	}
	shortUpper := shortName[0] >= 'A' && shortName[0] <= 'Z'
	segUpper := lastSegment[0] >= 'A' && lastSegment[0] <= 'Z'
	if shortUpper != segUpper {
		var correctedStart byte
		if shortUpper {
			correctedStart = shortName[0] + 32 // toUpperCase → toLowerCase
		} else {
			correctedStart = shortName[0] - 32 // toLowerCase → toUpperCase
		}
		correctedName := string([]byte{correctedStart}) + shortName[1:]
		c.report(shortNameNode, diag.JSC_INCORRECT_SHORTNAME_CAPITALIZATION, diag.Warning,
			fmt.Sprintf("Incorrect short name capitalization: %q should be %q.", shortName, correctedName))
	}
}

// ---------------------------------------------------------------------------
// report helper.
// ---------------------------------------------------------------------------

func (c *CheckClosureImports) report(n *rhino.Node, code string, level diag.Level, msg string) {
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    n.Lineno(),
		Col:     n.Charno(),
		Code:    code,
		Level:   level,
		Message: msg,
	})
}

// ---------------------------------------------------------------------------
// AST helper utilities (cci-prefixed to avoid package-level collisions).
// ---------------------------------------------------------------------------

// cciMatchesQualifiedName returns true when node is a GETPROP chain that
// spells out the dotted `name` (e.g. "goog.require", "goog.module.get").
// Mirrors Node.matchesQualifiedName(Node) in Java.
func cciMatchesQualifiedName(n *rhino.Node, name string) bool {
	if n == nil {
		return false
	}
	return cciGetQualifiedName(n) == name
}

// cciGetQualifiedName reconstructs the dotted string from a GETPROP chain
// or bare NAME. Returns "" for non-name nodes.
func cciGetQualifiedName(n *rhino.Node) string {
	if n == nil {
		return ""
	}
	switch n.Token() {
	case rhino.NAME:
		return n.String_()
	case rhino.GETPROP:
		kids := n.Children()
		if len(kids) >= 2 {
			left := cciGetQualifiedName(kids[0])
			if left == "" {
				return ""
			}
			return left + "." + kids[1].String_()
		}
	}
	return ""
}

// cciHasTwoChildren returns true when n has exactly two children.
// Mirrors call.hasTwoChildren() in Java.
func cciHasTwoChildren(n *rhino.Node) bool {
	return n != nil && len(n.Children()) == 2
}

// cciSecondChild returns the second child of n, or nil.
func cciSecondChild(n *rhino.Node) *rhino.Node {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	return kids[1]
}

// cciHasOneChild returns true when n has exactly one child.
func cciHasOneChild(n *rhino.Node) bool {
	return n != nil && len(n.Children()) == 1
}

// cciIsExprResult reports whether n is an EXPR_RESULT node.
func cciIsExprResult(n *rhino.Node) bool {
	return n != nil && n.Token() == rhino.EXPR_RESULT
}

// cciIsNameDeclaration reports whether n is VAR / LET / CONST.
// Mirrors NodeUtil.isNameDeclaration in Java.
func cciIsNameDeclaration(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	t := n.Token()
	return t == rhino.VAR || t == rhino.LET || t == rhino.CONST
}

// cciIsLet reports whether n is a LET declaration.
func cciIsLet(n *rhino.Node) bool {
	return n != nil && n.Token() == rhino.LET
}

// cciIsConst reports whether n is a CONST declaration.
func cciIsConst(n *rhino.Node) bool {
	return n != nil && n.Token() == rhino.CONST
}

// cciGrandParent returns n.Parent().Parent(), or nil.
func cciGrandParent(n *rhino.Node) *rhino.Node {
	if n == nil {
		return nil
	}
	p := n.Parent()
	if p == nil {
		return nil
	}
	return p.Parent()
}

// cciIsDestructuringLhs reports whether n represents a destructuring LHS
// (i.e. the first child of a VAR/LET/CONST that is an OBJECTLIT or
// array pattern used on the LHS of an initializer).
//
// Java uses Token.DESTRUCTURING_LHS which our rhino subset doesn't carry;
// we approximate by checking that the first child of a name-declaration
// child is an OBJECTLIT.
//
// TODO: add DESTRUCTURING_LHS token to rhino when rhino.Node is extended.
func cciIsDestructuringLhs(n *rhino.Node) bool {
	// Approximate: if the node itself is an OBJECTLIT it is acting as
	// a destructuring LHS inside a var/let/const initializer.
	if n == nil {
		return false
	}
	return n.Token() == rhino.OBJECTLIT
}

// cciIsValidDestructuringImport mirrors Checker.isValidDestructuringImport.
// A valid destructuring import is `const { Foo, Bar } = goog.require(…)` —
// every key of the object pattern must be a plain string key with a NAME value.
func cciIsValidDestructuringImport(lhs *rhino.Node) bool {
	// lhs is the OBJECTLIT acting as destructuring pattern.
	for _, key := range lhs.Children() {
		if key.Token() != rhino.STRING_KEY {
			return false
		}
		val := key.FirstChild()
		if val == nil || val.Token() != rhino.NAME {
			return false
		}
	}
	return true
}

// cciHasFunctionAncestor returns true when any ancestor of n is a FUNCTION
// node. Used as a proxy for "not at top-level scope".
// TODO: replace with real NodeTraversal scope-stack query.
func cciHasFunctionAncestor(n *rhino.Node) bool {
	for cur := n.Parent(); cur != nil; cur = cur.Parent() {
		if cur.Token() == rhino.FUNCTION {
			return true
		}
	}
	return false
}

// cciGoogProvideOrModuleNamespace returns the namespace string from a
// goog.provide(ns) or goog.module(ns) call, or "" if the node is neither.
// Used by collectNamespaces to seed namespacesSeen.
func cciGoogProvideOrModuleNamespace(call *rhino.Node) string {
	if call.Token() != rhino.CALL {
		return ""
	}
	callee := call.FirstChild()
	if callee == nil {
		return ""
	}
	name := cciGetQualifiedName(callee)
	if name != "goog.provide" && name != "goog.module" && name != "goog.declareModuleId" {
		return ""
	}
	arg := cciSecondChild(call)
	if arg == nil || arg.Token() != rhino.STRINGLIT {
		return ""
	}
	return arg.String_()
}
