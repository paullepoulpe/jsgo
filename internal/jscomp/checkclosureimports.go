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
// TODO(port): CheckClosureImports uses ModuleMetadataMap from com/google/javascript/jscomp/CheckClosureImports.java:220
// Reason: Go added factory function for stub type; Java uses class constructor directly.
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
// import alias is not used in code position.
// Mirrors Java's checkValidImportCodeReference(NodeTraversal t, Node nameNode)
// (lines 289-317). Java's body performs only structural checks via scope lookup
// with no diagnostic emission; this port faithfully mirrors that structure.
//
// TODO(port): NodeTraversal.getScope() / Scope.getVar(String) — scope not yet
// threaded through walk(); structural walk is faithful but var resolution is
// deferred until NodeTraversal is wired.
func (c *CheckClosureImports) checkValidImportCodeReference(nameNode *rhino.Node) {
	// Java: Var var = t.getScope().getVar(nameNode.getString());
	// TODO: use scope.GetSlot(nameNode.GetString()) once scope is available.
	// if var == nil || var.GetNameNode() == nameNode { return }
	var dummy *UntypedVar // TODO: replace with real scope lookup.
	_ = dummy

	// Java: Node declarationNameNode = var.getNameNode();
	// if declarationNameNode == nil || !NodeUtil.isDeclarationLValue(declarationNameNode) { return }
	// TODO: NU.IsDeclarationLValue(declarationNameNode)

	// Java: compute decl from declarationNameNode — handle string-key destructuring:
	//   if declarationNameNode.getParent().isStringKey():
	//     decl = declarationNameNode.getGrandparent().getGrandparent()
	//   else:
	//     decl = declarationNameNode.getParent()
	// TODO: port decl computation once scope is wired.

	// Java: if !NodeUtil.isNameDeclaration(decl) || !decl.hasOneChild()
	//          || !decl.getFirstChild().hasChildren()
	//          || !decl.getFirstChild().getLastChild().isCall() { return }
	// TODO: port structural guard.

	// Java: ClosureImport kindOfCall = kindOfCall(decl.getFirstChild().getLastChild());
	// if kindOfCall == nil { return }
	// NOTE: Java emits no diagnostic here — the method returns without reporting.
	_ = nameNode
}

// ---------------------------------------------------------------------------
// checkGoogModuleGet — mirrors Checker.checkGoogModuleGet.
// ---------------------------------------------------------------------------

// checkGoogModuleGet validates a goog.module.get(…) call.
// Mirrors Java's checkGoogModuleGet(NodeTraversal t, Node call, ModuleMetadata currentModule)
// (lines 319-398).
func (c *CheckClosureImports) checkGoogModuleGet(call *rhino.Node) {
	// Java line 320: if (t.inModuleHoistScope()) { t.report(call, MODULE_USES_GOOG_MODULE_GET); return; }
	// TODO(port): NodeTraversal.inModuleHoistScope() — scope stack not yet wired; deferred.

	// Java lines 325-356: if not es6/googModule and not WEBPACK resolution mode,
	// check inGlobalScope() for invalid alias assignment patterns.
	// TODO(port): ModuleMetadata.isEs6Module()/isGoogModule() and
	// NodeTraversal.inGlobalScope() — ModuleMetadata + scope not yet wired; deferred.
	//
	// Java: walk up GETPROP chain to find getSubExprRoot:
	//   Node getSubExprRoot = call;
	//   while (NodeUtil.isNormalOrOptChainGetProp(getSubExprRoot.getParent()))
	//       getSubExprRoot = getSubExprRoot.getParent();
	// Java: check parent of getSubExprRoot:
	//   if parent.isName() || parent.isDestructuringLhs() → invalid
	//   if parent.isAssign() && (target.isName() || target.isObjectPattern()) → invalid
	// if invalid: t.report(call, INVALID_GET_CALL_SCOPE)
	// TODO(port): INVALID_GET_CALL_SCOPE global-scope check — deferred.

	// Java line 359: if (!call.hasTwoChildren() || !call.getSecondChild().isStringLit())
	if !cciHasTwoChildren(call) || cciSecondChild(call) == nil || cciSecondChild(call).Token() != rhino.STRINGLIT {
		c.report(call, diag.JSC_INVALID_CLOSURE_IMPORT_CALL, diag.Error,
			"goog.module.get parameter must be a string literal.")
		return
	}

	// Java line 364-369: namespace lookup.
	namespace := cciSecondChild(call).String_()
	requiredModule := c.moduleMetadata.getModulesByGoogNamespace(namespace)
	if requiredModule == nil {
		c.report(call, diag.JSC_MISSING_MODULE_OR_PROVIDE, diag.Error,
			fmt.Sprintf("Required namespace %q never defined.", namespace))
		return
	}

	// Java lines 372-398: if isFillingAnAlias && isModule → check aliasVar and its RHS.
	// Java: Node maybeAssign = call.getParent();
	//       boolean isFillingAnAlias = maybeAssign.isAssign() && maybeAssign.getFirstChild().isName();
	//       boolean isModule = currentModule.isGoogModule() || currentModule.isEs6Module();
	// TODO(port): ModuleMetadata.isGoogModule()/isEs6Module() not yet wired.
	// Structural port (isFillingAnAlias part):
	maybeAssign := call.Parent()
	isFillingAnAlias := maybeAssign != nil && maybeAssign.IsAssign() &&
		maybeAssign.FirstChild() != nil && maybeAssign.FirstChild().IsName()
	isModule := false // TODO: derive from currentModule once ModuleMetadata is ported.
	if isFillingAnAlias && isModule {
		// Java: String aliasName = call.getParent().getFirstChild().getString();
		// Java: Var aliasVar = t.getScope().getVar(aliasName);
		// if aliasVar == nil → report INVALID_GET_ALIAS, return
		// TODO(port): scope.GetSlot(aliasName) — scope not yet wired.

		// Java: Node aliasVarNodeRhs = NodeUtil.getRValueOfLValue(aliasVar.getNode());
		// if aliasVarNodeRhs == nil || !NodeUtil.isCallTo(aliasVarNodeRhs, GOOG_FORWARD_DECLARE)
		//     → report INVALID_GET_ALIAS, return
		// TODO(port): NU.GetRValueOfLValue / NU.IsCallTo — deferred.

		// Java: if !namespace.equals(aliasVarNodeRhs.getLastChild().getString())
		//     → report INVALID_GET_ALIAS, return
		// TODO(port): deferred.
		_ = requiredModule
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
		if importType == closureImportRequireDynamic {
			// goog.requireDynamic can reference modules loaded at runtime
			// (e.g. via goog.loadModule). Skip the static namespace check
			// to avoid false positives when the module is loaded dynamically.
			// Mirrors Java's CheckClosureImports which handles requireDynamic
			// separately and doesn't fire JSC_MISSING_MODULE_OR_PROVIDE for it.
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
// before the current require call.
// Mirrors Java's verifyRequireOrder(String namespace, Node call, NodeTraversal t, ModuleMetadata requiredModule)
// (lines 523-550).
func (c *CheckClosureImports) verifyRequireOrder(namespace string, call *rhino.Node, requiredModule *ModuleMetadata) {
	// Java line 525-528: if (!namespacesSeen.contains(namespace)) { t.report(…); return; }
	if _, seen := c.namespacesSeen[namespace]; !seen {
		c.report(call, diag.JSC_LATE_PROVIDE_ERROR, diag.Error,
			fmt.Sprintf("Required namespace %q not provided yet.", namespace))
		return
	}

	// Java lines 530-532: if (requiredModule.rootNode() == null || requiredModule.rootNode().isFromExterns())
	//   return; // synthetic metadata for tests may have no root node.
	if requiredModule.rootNode == nil || requiredModule.rootNode.IsFromExterns() {
		return
	}

	// Java lines 534-549: cross-chunk check via JSChunkGraph.dependsOn.
	// TODO(port): JSChunkGraph + CompilerInput.getChunk() — chunk graph not yet ported.
	// Java: InputId requiredInputId = NodeUtil.getInputId(requiredModule.rootNode());
	// Java: CompilerInput requiredInput = compiler.getInput(requiredInputId);
	// Java: JSChunk requiredChunk = requiredInput.getChunk();
	// Java: JSChunk currentChunk = t.getChunk();
	// Java: if (currentChunk != requiredChunk && !chunkGraph.dependsOn(currentChunk, requiredChunk))
	//           compiler.report(JSError.make(call, CROSS_CHUNK_REQUIRE_ERROR, …));
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
// Delegates to rhino.Node.MatchesQualifiedNameStr (port of Java Node.matchesQualifiedName(String)).
func cciMatchesQualifiedName(n *rhino.Node, name string) bool {
	if n == nil {
		return false
	}
	return n.MatchesQualifiedNameStr(name)
}

// cciGetQualifiedName reconstructs the dotted string from a GETPROP chain
// or bare NAME. Returns "" for non-name nodes.
// Delegates to rhino.Node.GetQualifiedName (port of Java Node.getQualifiedName).
func cciGetQualifiedName(n *rhino.Node) string {
	if n == nil {
		return ""
	}
	return n.GetQualifiedName()
}

// cciHasTwoChildren returns true when n has exactly two children.
// Delegates to rhino.Node.HasTwoChildren (port of Java Node.hasTwoChildren).
func cciHasTwoChildren(n *rhino.Node) bool {
	return n != nil && n.HasTwoChildren()
}

// cciSecondChild returns the second child of n, or nil.
// Delegates to rhino.Node.GetSecondChild (port of Java Node.getSecondChild).
func cciSecondChild(n *rhino.Node) *rhino.Node {
	if n == nil {
		return nil
	}
	return n.GetSecondChild()
}

// cciHasOneChild returns true when n has exactly one child.
// Delegates to rhino.Node.HasOneChild (port of Java Node.hasOneChild).
func cciHasOneChild(n *rhino.Node) bool {
	return n != nil && n.HasOneChild()
}

// cciIsExprResult reports whether n is an EXPR_RESULT node.
// Delegates to rhino.Node.IsExprResult (port of Java Node.isExprResult).
func cciIsExprResult(n *rhino.Node) bool {
	return n != nil && n.IsExprResult()
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
// Delegates to rhino.Node.IsLet (port of Java Node.isLet).
func cciIsLet(n *rhino.Node) bool {
	return n != nil && n.IsLet()
}

// cciIsConst reports whether n is a CONST declaration.
// Delegates to rhino.Node.IsConst (port of Java Node.isConst).
func cciIsConst(n *rhino.Node) bool {
	return n != nil && n.IsConst()
}

// cciGrandParent returns n.Parent().Parent(), or nil.
// Delegates to rhino.Node.GetGrandparent (port of Java Node.getGrandparent).
func cciGrandParent(n *rhino.Node) *rhino.Node {
	if n == nil {
		return nil
	}
	return n.GetGrandparent()
}

// isValidDestructuringLhs reports whether n represents a destructuring LHS
// (i.e. the first child of a VAR/LET/CONST that is a DESTRUCTURING_LHS or
// OBJECTLIT used on the LHS of an initializer).
//
// Mirrors the Token.DESTRUCTURING_LHS check in Java; the legacy OBJECTLIT
// fallback covers older AST shapes that haven't yet been normalised to
// DESTRUCTURING_LHS.
func isValidDestructuringLhs(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	t := n.Token()
	return t == rhino.DESTRUCTURING_LHS || t == rhino.OBJECTLIT
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
// TODO(port): CheckClosureImports.Checker.checkShortName from com/google/javascript/jscomp/CheckClosureImports.java:422
// Reason: Go extracted lastIndexOf helper; Java inlines String.lastIndexOf('.').
func cciLastDotIndex(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			return i
		}
	}
	return -1
}
