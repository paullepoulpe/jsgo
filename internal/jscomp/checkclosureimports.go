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
//
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

// ModuleMetadata is a minimal stub for
// com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata.
// TODO: replace with the real ModuleMetadata port.
type ModuleMetadata struct {
	// googNamespaces is the set of goog.provide / goog.module / goog.declareModuleId
	// namespaces declared by the module's root node.
	googNamespaces map[string]struct{}
	moduleType     cciModuleType
	// rootNode is the SCRIPT node for the module; nil for synthetic/test metadata.
	rootNode *rhino.Node
}

// isGoogModule mirrors ModuleMetadata.isGoogModule().
func (m *ModuleMetadata) isGoogModule() bool { return m.moduleType == cciModuleTypeGoogModule }

// isEs6Module mirrors ModuleMetadata.isEs6Module().
func (m *ModuleMetadata) isEs6Module() bool { return m.moduleType == cciModuleTypeEs6Module }

// isModule mirrors ModuleMetadata.isModule() — true for goog.module and ES6 modules.
func (m *ModuleMetadata) isModule() bool {
	return m.isGoogModule() || m.isEs6Module()
}

// ModuleMetadataMap is a minimal stub for
// com.google.javascript.jscomp.modules.ModuleMetadataMap.
// TODO: replace with the real ModuleMetadataMap port.
type ModuleMetadataMap struct {
	// byGoogNamespace maps goog.provide / goog.module / goog.declareModuleId
	// namespace strings to their owning module metadata.
	byGoogNamespace map[string]*ModuleMetadata
}

// newModuleMetadataMap constructs an empty ModuleMetadataMap.
// Go-only: language-adapter — factory constructor for the stub type (Java uses class constructor directly; Go needs factory function).
func newModuleMetadataMap() *ModuleMetadataMap {
	return &ModuleMetadataMap{byGoogNamespace: map[string]*ModuleMetadata{}}
}

// getModulesByGoogNamespace mirrors
// ModuleMetadataMap.getModulesByGoogNamespace().get(namespace).
func (m *ModuleMetadataMap) getModulesByGoogNamespace(ns string) *ModuleMetadata {
	return m.byGoogNamespace[ns]
}

// ---------------------------------------------------------------------------
// ClosureImport enum — mirrors the ClosureImport enum with methods.
// ---------------------------------------------------------------------------

// closureImport mirrors the ClosureImport enum in Java's CheckClosureImports.
// Each constant represents one of the recognised goog.* import call types.
type closureImport int

const (
	closureImportRequire        closureImport = iota // goog.require
	closureImportForwardDeclare                      // goog.forwardDeclare
	closureImportRequireType                         // goog.requireType
	closureImportRequireDynamic                      // goog.requireDynamic
)

// allowDestructuring mirrors ClosureImport.allowDestructuring().
// Returns true for REQUIRE, REQUIRE_TYPE, REQUIRE_DYNAMIC.
func (ci closureImport) allowDestructuring() bool {
	switch ci {
	case closureImportRequire, closureImportRequireType, closureImportRequireDynamic:
		return true
	}
	return false
}

// aliasMustBeConstant mirrors ClosureImport.aliasMustBeConstant().
// Returns true for REQUIRE, REQUIRE_TYPE, REQUIRE_DYNAMIC.
func (ci closureImport) aliasMustBeConstant() bool {
	switch ci {
	case closureImportRequire, closureImportRequireType, closureImportRequireDynamic:
		return true
	}
	return false
}

// mustBeOrdered mirrors ClosureImport.mustBeOrdered().
// Only REQUIRE enforces ordering (provided namespace must appear before the call).
func (ci closureImport) mustBeOrdered() bool {
	return ci == closureImportRequire
}

// callName mirrors ClosureImport.callName().
// Returns the human-readable name of the import call for error messages.
func (ci closureImport) callName() string {
	switch ci {
	case closureImportRequire:
		return "goog.require"
	case closureImportForwardDeclare:
		return "goog.forwardDeclare"
	case closureImportRequireType:
		return "goog.requireType"
	case closureImportRequireDynamic:
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
// Java: final class CheckClosureImports implements CompilerPass.
// The Checker inner class (extends AbstractModuleCallback) is inlined here.
//
// Because ModuleMetadataMap is not yet ported, the pass accumulates
// diagnostics via a minimal stub map; namespacesSeen tracks which goog
// namespaces have been declared prior to a require (ordering check).
type CheckClosureImports struct {
	path           string
	moduleMetadata *ModuleMetadataMap // TODO: replace with real ModuleMetadataMap port
	namespacesSeen map[string]struct{}
	diags          []diag.Diagnostic
}

// NewCheckClosureImports mirrors the Java constructor:
//
//	CheckClosureImports(AbstractCompiler compiler, ModuleMetadataMap)
//
// The path argument is the source file path used in diagnostics.
// The metadata map is stubbed; pass nil to use an empty map.
// Go-only: language-adapter — Go constructor shim (Java uses class constructor directly; Go has no class constructors).
func NewCheckClosureImports(path string, metadata *ModuleMetadataMap) *CheckClosureImports {
	if metadata == nil {
		metadata = newModuleMetadataMap()
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

	// enterModule: seed namespacesSeen via AbstractModuleCallback.enterModule.
	// TODO: once ModuleMetadataMap is ported, derive from currentModule.googNamespaces().
	// For now walk the AST to collect goog.provide / goog.module namespace declarations
	// encountered before each require, so the ordering check (LATE_PROVIDE_ERROR) fires.
	c.enterModule(script)

	// Walk the AST and check every call that looks like a Closure import.
	c.walk(script, nil)

	return c.diags
}

// enterModule is called when the traversal enters a module scope.
// Mirrors AbstractModuleCallback.enterModule(ModuleMetadata currentModule, Node moduleScopeRoot).
// Java: namespacesSeen.addAll(currentModule.googNamespaces().elementSet()).
// Go stub: pre-populates namespacesSeen by walking the AST for goog.provide/module calls.
func (c *CheckClosureImports) enterModule(n *rhino.Node) {
	if n == nil {
		return
	}
	if n.Token() == rhino.CALL {
		if ns := cciGoogProvideOrModuleNamespace(n); ns != "" {
			c.namespacesSeen[ns] = struct{}{}
		}
	}
	for _, ch := range n.Children() {
		c.enterModule(ch)
	}
}

// walk visits every node depth-first. Mirrors the Checker inner-class
// visit() method, which receives both CALL and NAME nodes.
// TODO(port): NodeTraversal.traverse(AbstractCompiler, Node, Callback) — NodeTraversal adapter not yet ported.
func (c *CheckClosureImports) walk(n, parent *rhino.Node) {
	if n == nil {
		return
	}
	if n.Token() == rhino.CALL {
		importType := c.kindOfCall(n)
		if importType >= 0 {
			c.checkImport(n, parent, importType)
		} else if cciMatchesQualifiedName(n.FirstChild(), "goog.module.get") {
			c.checkGoogModuleGet(n)
		}
	} else if n.Token() == rhino.DYNAMIC_IMPORT {
		c.checkDynamicImportPath(n)
	} else if n.Token() == rhino.NAME {
		c.checkValidImportCodeReference(n)
	}
	for _, ch := range n.Children() {
		c.walk(ch, n)
	}
}

// ---------------------------------------------------------------------------
// kindOfCall — mirrors Checker.kindOfCall(Node callNode).
// ---------------------------------------------------------------------------

// kindOfCall returns the closureImport kind for callNode, or -1 if not recognized.
// Mirrors Java's private ClosureImport kindOfCall(Node callNode).
func (c *CheckClosureImports) kindOfCall(callNode *rhino.Node) closureImport {
	callee := callNode.FirstChild()
	switch {
	case cciMatchesQualifiedName(callee, "goog.require"):
		return closureImportRequire
	case cciMatchesQualifiedName(callee, "goog.forwardDeclare"):
		return closureImportForwardDeclare
	case cciMatchesQualifiedName(callee, "goog.requireType"):
		return closureImportRequireType
	case cciMatchesQualifiedName(callee, "goog.requireDynamic"):
		return closureImportRequireDynamic
	}
	return -1
}

// ---------------------------------------------------------------------------
// checkValidImportCodeReference — mirrors Checker.checkValidImportCodeReference.
// ---------------------------------------------------------------------------

// checkValidImportCodeReference validates that a NAME reference to a closure
// import alias is not used in code position (mirrors Java lines 289-317).
// The full check requires scope infrastructure (Var lookup); this stub
// only performs the structural check that is possible without it.
//
// TODO: implement full scope-based check once TypedScopeCreator is wired.
func (c *CheckClosureImports) checkValidImportCodeReference(nameNode *rhino.Node) {
	// TODO: use t.getScope().getVar(nameNode.getString()) to find the declaration.
	// TODO: check NodeUtil.isDeclarationLValue(declarationNameNode).
	// TODO: check kindOfCall on the RHS call and fire JSC_TYPE_IMPORT_CODE_REFERENCE.
	_ = nameNode
}

// ---------------------------------------------------------------------------
// checkGoogModuleGet — mirrors Checker.checkGoogModuleGet.
// ---------------------------------------------------------------------------

// checkGoogModuleGet validates a goog.module.get(…) call.
// Mirrors Java's checkGoogModuleGet(NodeTraversal t, Node call, ModuleMetadata currentModule)
// (lines 319-398).
// Full scope/chunk checks require ModuleMetadata + scope stack; those
// are stubbed until the module system is ported.
func (c *CheckClosureImports) checkGoogModuleGet(call *rhino.Node) {
	// TODO: check t.inModuleHoistScope() → MODULE_USES_GOOG_MODULE_GET (Java line 320).
	// TODO: check INVALID_GET_CALL_SCOPE for goog.provide files in global scope (lines 325-357).
	// TODO: check INVALID_GET_ALIAS when filling a forwardDeclare alias (lines 372-398).

	// Validate argument: must have two children, second must be string literal.
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
// checkDynamicImportPath — mirrors BrowserModuleResolver.resolveJsModule.
// ---------------------------------------------------------------------------

// checkDynamicImportPath validates the module specifier in a dynamic
// import(…) expression. In BROWSER resolution mode only absolute paths
// (starting with "/") and relative paths (starting with "./" or "../")
// are valid. Any other specifier — e.g. bare "foo.js" — is "ambiguous"
// and triggers JSC_INVALID_MODULE_PATH.
// TODO(port): BrowserModuleResolver.resolve(String) — BrowserModuleResolver not yet ported; logic inlined here.
func (c *CheckClosureImports) checkDynamicImportPath(n *rhino.Node) {
	specifier := n.FirstChild()
	if specifier == nil || specifier.Token() != rhino.STRINGLIT {
		return
	}
	path := specifier.String_()
	// isAmbiguousIdentifier: not absolute ("/…") and not relative ("./…" or "../…").
	isAbsolute := len(path) > 0 && path[0] == '/'
	isRelative := len(path) >= 2 && path[:2] == "./"
	if !isRelative && len(path) >= 3 {
		isRelative = path[:3] == "../"
	}
	if !isAbsolute && !isRelative {
		c.report(specifier, diag.JSC_INVALID_MODULE_PATH, diag.Error,
			fmt.Sprintf("Invalid module path %q for resolution mode \"BROWSER\"", path))
	}
}

// ---------------------------------------------------------------------------
// checkImport — mirrors Checker.checkImport.
// ---------------------------------------------------------------------------

// checkImport validates a single Closure-import call (require / forwardDeclare /
// requireType / requireDynamic). `call` is the CALL node; `parent` is its
// parent in the AST.
// Mirrors Java's checkImport(NodeTraversal t, Node call, Node parent,
//
//	ModuleMetadata currentModule, ClosureImport importType) (lines 437-520).
func (c *CheckClosureImports) checkImport(call, parent *rhino.Node, importType closureImport) {
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

	if (!atTopLevel || !validAssignment) && importType != closureImportRequireDynamic {
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
	isAliased := parent != nil && cciIsNameDecl(cciGrandParent(call))
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
			if lhs != nil && isValidDestructuringLhs(lhs) {
				if importType.allowDestructuring() {
					if !isValidDestructuringImport(lhs) {
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
				c.checkShortName(lhs, ns)
			}
		}
	}

	// -- namespace existence check ------------------------------------------
	namespace := cciSecondChild(call).String_()
	requiredModule := c.moduleMetadata.getModulesByGoogNamespace(namespace)

	if requiredModule == nil {
		if importType == closureImportForwardDeclare {
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
		c.verifyRequireOrder(namespace, call, requiredModule)
	}

	// -- ES6 module mutual-import check -------------------------------------
	// TODO: fire JSC_SHOULD_IMPORT_ES6_MODULE when both currentModule and
	// requiredModule are ES6 modules. Requires ModuleMetadata.isEs6Module().
}

// ---------------------------------------------------------------------------
// verifyRequireOrder — mirrors CheckClosureImports.verifyRequireOrder.
// ---------------------------------------------------------------------------

// verifyRequireOrder checks that the required namespace was declared
// before the current require call. Mirrors Java's
// verifyRequireOrder(String namespace, Node call, NodeTraversal t, ModuleMetadata requiredModule)
// (lines 523-550).
func (c *CheckClosureImports) verifyRequireOrder(namespace string, call *rhino.Node, _ *ModuleMetadata) {
	if _, seen := c.namespacesSeen[namespace]; !seen {
		c.report(call, diag.JSC_LATE_PROVIDE_ERROR, diag.Error,
			fmt.Sprintf("Required namespace %q not provided yet.", namespace))
		return
	}
	// TODO: cross-chunk check (JSC_XMODULE_REQUIRE_ERROR) — requires
	// JSChunkGraph + CompilerInput.getChunk(), neither of which is ported yet.
}

// ---------------------------------------------------------------------------
// checkShortName — mirrors Checker.checkShortName.
// ---------------------------------------------------------------------------

// checkShortName validates that the short alias name for a require has
// matching capitalization with the last segment of the namespace.
// Mirrors Java's checkShortName(NodeTraversal t, Node shortNameNode, String namespace)
// (lines 420-435).
func (c *CheckClosureImports) checkShortName(shortNameNode *rhino.Node, namespace string) {
	if shortNameNode == nil {
		return
	}
	shortName := shortNameNode.String_()
	lastSegment := namespace[cciLastDotIndex(namespace)+1:]
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
// isValidDestructuringImport — mirrors Checker.isValidDestructuringImport.
// ---------------------------------------------------------------------------

// isValidDestructuringImport mirrors Java's isValidDestructuringImport(Node destructuringLhs).
// A valid destructuring import is `const { Foo, Bar } = goog.require(…)` —
// every key of the object pattern must be a plain string key with a NAME value.
// Java (lines 401-418): checks objectPattern.isObjectPattern(), then each stringKey.isStringKey()
// and stringKey.getFirstChild().isName().
func isValidDestructuringImport(lhs *rhino.Node) bool {
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

// ---------------------------------------------------------------------------
// report helper.
// ---------------------------------------------------------------------------

// report emits a diagnostic. Mirrors Java's t.report(node, diagnosticType, args...).
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
// TODO(port): Node.matchesQualifiedName(Node) — qualified-name matcher (Java uses Node.matchesQualifiedName; rhino.Node does not have it yet).
func cciMatchesQualifiedName(n *rhino.Node, name string) bool {
	if n == nil {
		return false
	}
	return cciGetQualifiedName(n) == name
}

// cciGetQualifiedName reconstructs the dotted string from a GETPROP chain
// or bare NAME. Returns "" for non-name nodes.
// TODO(port): Node.getQualifiedName() — qualified-name builder (Java uses Node.getQualifiedName; rhino.Node does not have it yet).
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
// TODO(port): Node.hasTwoChildren() — rhino.Node does not yet expose this method.
func cciHasTwoChildren(n *rhino.Node) bool {
	return n != nil && len(n.Children()) == 2
}

// cciSecondChild returns the second child of n, or nil.
// Mirrors Node.getSecondChild() in Java.
// TODO(port): Node.getSecondChild() — rhino.Node does not yet expose this method.
func cciSecondChild(n *rhino.Node) *rhino.Node {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	return kids[1]
}

// cciHasOneChild returns true when n has exactly one child.
// Mirrors Node.hasOneChild() in Java.
// TODO(port): Node.hasOneChild() — rhino.Node does not yet expose this method.
func cciHasOneChild(n *rhino.Node) bool {
	return n != nil && len(n.Children()) == 1
}

// cciIsExprResult reports whether n is an EXPR_RESULT node.
// Mirrors Node.isExprResult() in Java.
// TODO(port): Node.isExprResult() — rhino.Node does not yet expose this method.
func cciIsExprResult(n *rhino.Node) bool {
	return n != nil && n.Token() == rhino.EXPR_RESULT
}

// cciIsNameDecl reports whether n is VAR / LET / CONST.
// Mirrors NodeUtil.isNameDeclaration in Java.
// TODO(port): NodeUtil.isNameDeclaration(Node) — NodeUtil method not fully ported; adapter here due to package-level name clash.
func cciIsNameDecl(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	t := n.Token()
	return t == rhino.VAR || t == rhino.LET || t == rhino.CONST
}

// cciIsLet reports whether n is a LET declaration.
// Mirrors Node.isLet() in Java.
// TODO(port): Node.isLet() — rhino.Node does not yet expose this method.
func cciIsLet(n *rhino.Node) bool {
	return n != nil && n.Token() == rhino.LET
}

// cciIsConst reports whether n is a CONST declaration.
// Mirrors Node.isConst() in Java.
// TODO(port): Node.isConst() — rhino.Node does not yet expose this method.
func cciIsConst(n *rhino.Node) bool {
	return n != nil && n.Token() == rhino.CONST
}

// cciGrandParent returns n.Parent().Parent(), or nil.
// Mirrors Node.getGrandparent() in Java.
// TODO(port): Node.getGrandparent() — rhino.Node does not yet expose this method.
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

// isValidDestructuringLhs reports whether n represents a destructuring LHS
// (i.e. the first child of a VAR/LET/CONST that is an OBJECTLIT or
// array pattern used on the LHS of an initializer).
//
// Java uses Token.DESTRUCTURING_LHS which our rhino subset doesn't carry;
// we approximate by checking that the first child of a name-declaration
// child is an OBJECTLIT.
//
// TODO: add DESTRUCTURING_LHS token to rhino when rhino.Node is extended.
// TODO(port): Token.DESTRUCTURING_LHS — rhino token not yet ported; approximating with OBJECTLIT check.
func isValidDestructuringLhs(n *rhino.Node) bool {
	// Approximate: if the node itself is an OBJECTLIT it is acting as
	// a destructuring LHS inside a var/let/const initializer.
	if n == nil {
		return false
	}
	return n.Token() == rhino.OBJECTLIT
}

// cciHasFunctionAncestor returns true when any ancestor of n is a FUNCTION
// node. Used as a proxy for "not at top-level scope".
// TODO: replace with real NodeTraversal scope-stack query.
// TODO(port): NodeTraversal.inGlobalScope() / Scope.isGlobal() — scope infrastructure not yet ported; using stack proxy.
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
// Used by enterModule to seed namespacesSeen.
// TODO(port): ModuleMetadata.getGoogNamespace() — namespace extractor for enterModule (ModuleMetadata not fully ported).
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

// cciLastDotIndex returns the index of the last '.' in s, or -1.
// Go-only: language-adapter — strings.LastIndex helper for checkShortName (Java uses String.lastIndexOf('.') inline; Go extracts for clarity).
func cciLastDotIndex(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			return i
		}
	}
	return -1
}
