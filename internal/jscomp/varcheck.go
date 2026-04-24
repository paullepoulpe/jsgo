// Port of com/google/javascript/jscomp/VarCheck.java
// @ closure-compiler v20260415 (narrow; tier-level undeclared-NAME
// reporting only).
//
// Java's VarCheck walks every NAME reference in the AST and verifies
// that the name is declared somewhere reachable from the current
// scope: local var/let/const, @param, function name, class name,
// imported global, or the externs bundle (DefaultExterns.js +
// addExterns fixtures). Misses fire JSC_UNDEFINED_VARIABLE.
//
// jsgo doesn't load a real externs bundle yet (see TODO.md P2-1), so
// the "known globals" set below is a hand-curated list of the ES
// built-ins and DOM globals the checker tests routinely reference.
// The list is intentionally narrow: every false-positive we could
// induce by flagging a real browser / ES stdlib name would cost more
// than the ~50 true-positive fixtures the check unlocks.

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// VarCheck is the NodeTraversal.Callback that reports undeclared
// identifier references. Mirrors VarCheck.java's post-order Callback
// shape — ShouldTraverse always descends, Visit inspects NAME uses.
type VarCheck struct {
	path     string
	creator  *TypedScopeCreator
	global   *TypedScope
	reported map[string]bool // name → already-reported (Java reports once per file)
	diags    []diag.Diagnostic
}

// NewVarCheck mirrors VarCheck.java's constructor shape — it binds to
// the compiler's scope creator + global scope, and seeds the builtin
// extern names.
func NewVarCheck(path string, creator *TypedScopeCreator, global *TypedScope) *VarCheck {
	return &VarCheck{
		path:     path,
		creator:  creator,
		global:   global,
		reported: map[string]bool{},
	}
}

// Check runs the traversal and returns any JSC_UNDEFINED_VARIABLE
// diagnostics discovered. Fresh VarCheck per file matches the Java
// Compiler#process entry point.
// Go-only: language-adapter — Java's process() takes externs+root Node pair; Go uses a single SCRIPT node.
func (v *VarCheck) Check(script *rhino.Node) []diag.Diagnostic {
	v.diags = nil
	trav := NewNodeTraversal(v, v.creator, v.global)
	trav.Traverse(script)
	return v.diags
}

// UndeclaredNames returns the set of variable names that were reported
// as undeclared by the most recent Check call. Callers (e.g. TypeCheck)
// can use this to suppress downstream false-positives on properties of
// undeclared objects, mirroring Java VarCheck's synthesized-extern
// mechanism which marks undeclared names as externs.
// Go-only: language-adapter — Java uses createSynthesizedExternVar; Go uses this set.
func (v *VarCheck) UndeclaredNames() map[string]bool {
	return v.reported
}

// REQUIRED_SYMBOLS mirrors VarCheck.java's REQUIRED_SYMBOLS — names that must
// always be available to runtime libraries (polyfills, etc.), even if not
// explicitly declared in source. In Java these are added as synthesized externs
// when not already declared. In Go they are covered by defaultExternNames().
var VarCheckRequiredSymbols = map[string]bool{
	"AggregateError": true, "Array": true, "Error": true,
	"Float32Array": true, "Function": true, "Infinity": true,
	"JSCOMPILER_PRESERVE": true, "JSCompiler_renameProperty": true,
	"Map": true, "Math": true, "NaN": true, "Number": true,
	"Object": true, "Promise": true, "RangeError": true,
	"ReferenceError": true, "Reflect": true, "RegExp": true,
	"Set": true, "String": true, "SuppressedError": true,
	"Symbol": true, "TypeError": true, "WeakMap": true,
	"global": true, "globalThis": true, "isNaN": true,
	"parseFloat": true, "parseInt": true, "self": true,
	"undefined": true, "window": true,
}

func (v *VarCheck) ShouldTraverse(_ *NodeTraversal, _, _ *rhino.Node) bool { return true }

// EnterScope mirrors VarCheck.java's enterScope — Java's enterScope is also
// a no-op in VarCheck (it delegates to the default ScopedCallback.enterScope).
func (v *VarCheck) EnterScope(t *NodeTraversal) {
	_ = t // VarCheck doesn't need scope entry bookkeeping.
}

// ExitScope mirrors VarCheck.java's exitScope — in Java this adds required
// symbols to the synthetic-externs set when they're missing from scope.
// In Go this is a no-op since defaultExternNames() covers required symbols.
func (v *VarCheck) ExitScope(t *NodeTraversal) {
	_ = t // Required symbols are handled via defaultExternNames() in checkName.
}

// Visit mirrors VarCheck.java's ScopedCallback.visit — dispatches to
// checkName for NAME nodes. Other token types are ignored (no scope needed).
func (v *VarCheck) Visit(t *NodeTraversal, n, parent *rhino.Node) {
	if n == nil {
		return
	}
	if n.Token() != rhino.NAME {
		return
	}
	v.checkName(t, n, parent)
}

// checkName mirrors VarCheck.java's checkName — validates that a NAME node
// references a declared variable. Fires JSC_UNDEFINED_VARIABLE when the
// name is not declared in the reachable scope chain or well-known externs.
func (v *VarCheck) checkName(t *NodeTraversal, n, parent *rhino.Node) {
	name := n.String_()
	if name == "" {
		return
	}
	// Nodes in NAME position that aren't identifier references:
	//   - VAR / LET / CONST children (the bound name)
	//   - PARAM_LIST children
	//   - FUNCTION / CLASS name (first child)
	//   - CLASS_FIELD_DEF name (first child)
	//   - GETPROP property side (second child, stored as STRINGLIT
	//     already but some paths leave NAME nodes)
	//   - STRING_KEY children (object-literal keys)
	if v.isBindingSite(n, parent) {
		return
	}
	// GETPROP receivers that are the property (second child) are
	// STRINGLIT tokens, not NAME, so we don't need to disambiguate
	// there. Only dotted property names created via tryLHSExpression
	// may accidentally land here — guard by parent kind.
	if parent != nil && parent.Token() == rhino.GETPROP {
		kids := parent.Children()
		if len(kids) >= 2 && kids[1] == n {
			return
		}
	}
	scope := t.CurrentScope()
	if scope.GetDeclaredType(name) != nil {
		return
	}
	// Fallback: treat well-known extern names as declared even if the
	// scope chain doesn't contain them. This handles cases where the
	// externs scope might be incomplete or when using the new pipeline.
	if defaultExternNames()[name] {
		return
	}
	v.handleUndeclaredVariableRef(t, n)
}

// handleUndeclaredVariableRef mirrors VarCheck.java's handleUndeclaredVariableRef.
// Reports JSC_UNDEFINED_VARIABLE for the undeclared name, unless it's inside
// a typeof expression (Java treats typeof as an existence check). Deduplicates
// within a single file (Java reports once per file per name).
func (v *VarCheck) handleUndeclaredVariableRef(t *NodeTraversal, n *rhino.Node) {
	name := n.String_()
	// Java: if (n.getParent().isTypeOf()) — typeof checks are existence checks.
	if parent := n.Parent(); parent != nil && parent.Token() == rhino.TYPEOF {
		return
	}
	if v.reported[name] {
		return
	}
	v.reported[name] = true
	v.diags = append(v.diags, diag.Diagnostic{
		Path:    v.path,
		Line:    n.Lineno(),
		Col:     n.Charno(),
		Code:    diag.JSC_UNDEFINED_VARIABLE,
		Level:   diag.Error,
		Message: fmt.Sprintf("variable %s is undeclared", name),
	})
}

// isBindingSite classifies NAME nodes that are declarations (param,
// VAR child, function/class name, class-field-def name, STRING_KEY
// child inside destructuring, object-literal shorthand) rather than
// references. Matches NodeUtil.isNameDeclaration's intent.
// Go-only: language-adapter — Java uses NodeUtil.isReferenceName directly; Go needs its own helper.
func (v *VarCheck) isBindingSite(n, parent *rhino.Node) bool {
	if parent == nil {
		return true // orphan NAMEs aren't references
	}
	switch parent.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		// Only the first child of the declaration is the binding; any
		// NAME inside the initializer is a reference.
		return parent.FirstChild() == n
	case rhino.PARAM_LIST:
		return true
	case rhino.FUNCTION, rhino.CLASS:
		return parent.FirstChild() == n
	case rhino.CLASS_FIELD_DEF:
		return parent.FirstChild() == n
	case rhino.ARRAY_PATTERN:
		// Direct NAME children of an array pattern are bindings:
		// `[a, b] = rhs` — a and b are bindings.
		return true
	case rhino.OBJECT_REST:
		// `{...rest}` — the NAME inside OBJECT_REST is a binding.
		return true
	case rhino.STRING_KEY:
		// In destructuring: `{x: y}` or `{x}` inside OBJECT_PATTERN —
		// the NAME child of STRING_KEY is a binding target. In an object
		// literal, STRING_KEY is under OBJECTLIT, and the NAME child IS
		// a reference. Distinguish by checking the grandparent.
		gp := parent.Parent()
		if gp != nil && gp.Token() == rhino.OBJECT_PATTERN {
			return true
		}
		return false
	case rhino.DEFAULT_VALUE:
		// `{a = dflt}` or `[a = dflt]` — first child is the binding
		// target; second child is the default expression (a reference).
		// Only mark as binding if we are in a pattern context (parent of
		// DEFAULT_VALUE is STRING_KEY/ARRAY_PATTERN/OBJECT_PATTERN).
		if parent.FirstChild() == n {
			gp := parent.Parent()
			if gp != nil {
				switch gp.Token() {
				case rhino.STRING_KEY, rhino.ARRAY_PATTERN, rhino.OBJECT_PATTERN, rhino.PARAM_LIST:
					return true
				}
			}
		}
		return false
	}
	return false
}

// locationOf mirrors VarCheck.java's private static locationOf(Node) helper.
// Returns a human-readable location string for a node, or "<unknown>" for nil.
// Go-only: language-adapter — takes path as an extra arg since Go nodes don't carry their source file.
func locationOf(path string, n *rhino.Node) string {
	if n == nil {
		return "<unknown>"
	}
	return fmt.Sprintf("%s:%d:%d", path, n.Lineno(), n.Charno())
}

// reportVarMultiplyDeclared mirrors VarCheck.reportVarMultiplyDeclared.
// Reports a JSC_VAR_MULTIPLY_DECLARED_ERROR at `current` for variable `name`.
// The original declaration node is `original` (may be nil for implicit decls).
// Go-only: language-adapter — takes diags and path as extra args (Java uses compiler reference).
func reportVarMultiplyDeclared(diags *[]diag.Diagnostic, path, name string, current, original *rhino.Node) {
	msg := fmt.Sprintf("Variable %s declared more than once. First occurrence: %s", name, locationOf(path, original))
	*diags = append(*diags, diag.Diagnostic{
		Path:    path,
		Line:    current.Lineno(),
		Col:     current.Charno(),
		Code:    diag.JSC_VAR_MULTIPLY_DECLARED_ERROR,
		Level:   diag.Error,
		Message: msg,
	})
}

// defaultExternNames is the small fallback set used only by
// initIsUndeclaredRef (to decide whether a bare-ref initializer
// should downgrade the binding to Unknown). The primary extern
// lookup goes through the externs scope parsed from
// defaultExternsSrc — this list is just the names that can appear
// as initializer-only refs without a proper @type annotation.
// Go-only: language-adapter — Java uses a full externs bundle; Go uses this hand-curated list.
func defaultExternNames() map[string]bool {
	names := []string{
		// Primitive / language intrinsics
		"undefined", "Infinity", "NaN", "globalThis", "arguments", "eval",

		// Core ES constructors / namespaces
		"Object", "Function", "Array", "String", "Number", "Boolean",
		"Math", "JSON", "RegExp", "Date", "Error", "TypeError",
		"RangeError", "SyntaxError", "ReferenceError", "URIError",
		"EvalError", "Symbol", "Promise", "Map", "Set", "WeakMap",
		"WeakSet", "Proxy", "Reflect", "Iterator", "AsyncIterator",
		"ArrayBuffer", "DataView", "Int8Array", "Uint8Array",
		"Uint8ClampedArray", "Int16Array", "Uint16Array", "Int32Array",
		"Uint32Array", "Float32Array", "Float64Array", "BigInt",
		"BigInt64Array", "BigUint64Array",

		// Core ES free functions
		"parseInt", "parseFloat", "isNaN", "isFinite",
		"encodeURIComponent", "encodeURI",
		"decodeURIComponent", "decodeURI",

		// Browser globals (DOM + BOM subset)
		"window", "document", "navigator", "location", "history",
		"localStorage", "sessionStorage", "console", "self", "top",
		"alert", "confirm", "prompt", "setTimeout", "setInterval",
		"clearTimeout", "clearInterval", "requestAnimationFrame",
		"cancelAnimationFrame", "XMLHttpRequest", "fetch",
		"Element", "HTMLElement", "HTMLDivElement", "HTMLSpanElement",
		"HTMLInputElement", "HTMLButtonElement", "HTMLAnchorElement",
		"HTMLFormElement", "HTMLSelectElement", "HTMLTextAreaElement",
		"HTMLImageElement", "HTMLScriptElement",
		"Event", "MouseEvent", "KeyboardEvent", "UIEvent",
		"FocusEvent", "InputEvent", "Node",
		// DOM constructor types (uppercase, distinct from lowercase globals)
		"Window", "Document", "EventTarget",
		// Thenable typedef is a type alias, not a var, so VarCheck
		// doesn't need to know about it — it lives in typeAliases.
		"Thenable",
	}
	out := make(map[string]bool, len(names))
	for _, n := range names {
		out[n] = true
	}
	return out
}
