// Port of com/google/javascript/jscomp/VariableReferenceCheck.java +
// the duplicate-decl branch of com/google/javascript/jscomp/VarCheck.java
// @ closure-compiler v20260415 (narrow; intra-file scope only).
//
// Java's VariableReferenceCheck walks the AST per-scope and records
// every value-namespace declaration (var/let/const/function/class/
// @param). A second declaration of the same name in the same scope
// fires JSC_VAR_MULTIPLY_DECLARED_ERROR, with a message pointing at
// the first occurrence's file:line:col. jsgo doesn't load the real
// externs bundle yet, so we can only detect duplicates that both
// live in the source file under check (the externs-vs-source
// collisions — ~26 extracted fixtures against es6.js's Disposable
// etc. — need P2-1 to land first).
//
// REFERENCE_BEFORE_DECLARE (JSC_REFERENCE_BEFORE_DECLARE): fires as a
// WARNING when a `var` binding is referenced before its declaration
// within the same container scope (script or function body). Mirrors
// Java's checkEarlyReference / EARLY_REFERENCE path in
// ReferenceCheckingBehavior.checkVar. Note: `let`/`const` early
// references fire EARLY_REFERENCE_ERROR (a harder error) in Java —
// that code path is not yet ported.

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// VariableReferenceCheck is the NodeTraversal.Callback that reports
// intra-file duplicate `var` declarations and early-reference warnings.
// Mirrors the pass of the same name on the Java side.
type VariableReferenceCheck struct {
	path  string
	diags []diag.Diagnostic

	// Stack of scope maps — one entry per enclosing function /
	// script. Each map key is a declared name and the value is the
	// first occurrence's source position, used to build the "First
	// occurrence: <path>:<line>:<col>" message.
	scopes []map[string]declSite

	// earlyRefStack mirrors scopes: earlyRefStack[i] holds the
	// complete set of `var`-declared names (pre-scanned) in the
	// container scope at depth i, used to detect early references.
	// Only `var` declarations populate this map (not let/const/
	// function — hoisted functions are always "declared" from the
	// start; let/const fire a different error not yet ported).
	earlyRefStack []map[string]earlyRefDecl

	// suppressDepth is non-zero whenever we're inside a function
	// whose JSDoc carries `@suppress {duplicate}`. While > 0,
	// recordDecl skips emitting the duplicate diagnostic. Mirrors
	// Java's per-pass suppression-stack behaviour.
	suppressDepth int
}

type declSite struct {
	line int
	col  int
}

// earlyRefDecl records the source position of a `var` declaration for
// the purpose of early-reference detection.
type earlyRefDecl struct {
	line int
	col  int
}

func NewVariableReferenceCheck(path string, externScope *TypedScope) *VariableReferenceCheck {
	v := &VariableReferenceCheck{
		path:          path,
		scopes:        []map[string]declSite{{}},
		earlyRefStack: []map[string]earlyRefDecl{{}},
	}
	// Seed the global frame with every name declared by the externs
	// scope. Re-declaring one of these in user source fires
	// JSC_VAR_MULTIPLY_DECLARED_ERROR at the source-side position,
	// which matches Java's oracle.
	if externScope != nil {
		for name := range externScope.declared {
			v.scopes[0][name] = declSite{line: 0, col: 0}
		}
	}
	return v
}

// Check runs the traversal and returns the collected diagnostics.
// Callers are expected to build a fresh VariableReferenceCheck per
// file, matching Java's Compiler pass-instance pattern.
func (v *VariableReferenceCheck) Check(script *rhino.Node) []diag.Diagnostic {
	v.diags = nil
	// Pre-scan script-level var declarations for the early-reference check.
	v.preScanVarDecls(script, v.earlyRefStack[0])
	v.walk(script, nil)
	return v.diags
}

// preScanVarDecls collects every `var`-keyword declaration in `n`'s
// immediate container scope (i.e. it does not descend into nested
// function bodies) and records the name → position in `dest`. This
// mirrors Java's two-pass behaviour where references are collected
// in parse-tree order and then checked against the declaration after
// the scope exits.
func (v *VariableReferenceCheck) preScanVarDecls(n *rhino.Node, dest map[string]earlyRefDecl) {
	if n == nil {
		return
	}
	// Do not cross function boundaries — each function has its own
	// container scope that is pre-scanned separately before walking.
	if n.Token() == rhino.FUNCTION {
		return
	}
	if n.Token() == rhino.VAR {
		for _, ch := range n.Children() {
			if ch.Token() == rhino.NAME && ch.String_() != "" {
				name := ch.String_()
				if _, exists := dest[name]; !exists {
					dest[name] = earlyRefDecl{line: ch.Lineno(), col: ch.Charno()}
				}
			}
		}
	}
	for _, ch := range n.Children() {
		v.preScanVarDecls(ch, dest)
	}
}

func (v *VariableReferenceCheck) walk(n, parent *rhino.Node) {
	if n == nil {
		return
	}
	// Enter a new scope frame for every FUNCTION (including arrows);
	// CLASS bodies don't create a value-namespace scope of their own
	// in the Java model, so we don't push one here.
	if n.Token() == rhino.FUNCTION {
		v.scopes = append(v.scopes, map[string]declSite{})
		earlyRefs := map[string]earlyRefDecl{}
		v.earlyRefStack = append(v.earlyRefStack, earlyRefs)
		// `@suppress {duplicate}` on the enclosing function silences
		// duplicate-declaration diagnostics for everything inside its
		// body — matches the Java pass's gate.
		suppressed := false
		if info := n.JSDocInfo(); info != nil && info.HasSuppression("duplicate") {
			v.suppressDepth++
			suppressed = true
		}
		defer func() {
			v.scopes = v.scopes[:len(v.scopes)-1]
			v.earlyRefStack = v.earlyRefStack[:len(v.earlyRefStack)-1]
			if suppressed {
				v.suppressDepth--
			}
		}()
		// Function parameters are declarations in the new scope.
		kids := n.Children()
		if len(kids) >= 2 {
			for _, p := range kids[1].Children() {
				v.recordDecl(p, p.String_())
			}
			// Pre-scan the function body (third child = BLOCK) for var
			// declarations so early-reference checks work for references
			// that textually precede the var statement inside the function.
			if len(kids) >= 3 {
				v.preScanVarDecls(kids[2], earlyRefs)
			}
		}
	}
	switch n.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		// Multi-binding forms like `var a, b, c` each fire their own
		// duplicate check. Java treats every NAME child as a
		// declaration site.
		for _, name := range n.Children() {
			if name.Token() == rhino.NAME {
				v.recordDecl(name, name.String_())
			}
		}
	case rhino.FUNCTION, rhino.CLASS:
		// Named function / class declarations create a binding in the
		// enclosing scope. Skip members — class methods (FUNCTION
		// children of a CLASS-owned BLOCK) and object-literal methods
		// (FUNCTION under STRING_KEY under OBJECTLIT) are properties,
		// not value-namespace declarations.
		if isMemberContext(parent) {
			break
		}
		if parent != nil && (parent.Token() == rhino.BLOCK || parent.Token() == rhino.SCRIPT) {
			if kids := n.Children(); len(kids) >= 1 && kids[0].String_() != "" {
				name := kids[0]
				// The declaration lands in the enclosing scope (the
				// one above the just-pushed function scope). Use
				// scopes[-2] when we just pushed a function frame.
				if n.Token() == rhino.FUNCTION {
					target := v.scopes[len(v.scopes)-2]
					v.recordDeclInto(target, name, name.String_())
				} else {
					v.recordDecl(name, name.String_())
				}
			}
		}
	case rhino.NAME:
		// Check for an early reference: a NAME used as a value reference
		// (not a declaration binding) before its `var` declaration in
		// the same container scope. Mirrors Java's checkEarlyReference
		// / EARLY_REFERENCE (JSC_REFERENCE_BEFORE_DECLARE) warning for
		// var-declared names.
		v.checkEarlyReference(n, parent)
	}
	for _, c := range n.Children() {
		v.walk(c, n)
	}
}

// checkEarlyReference mirrors Java's checkEarlyReference for the
// EARLY_REFERENCE (JSC_REFERENCE_BEFORE_DECLARE) warning. It fires
// when a NAME node is a value reference (not a binding site) to a
// `var` that is declared later in the same container scope.
//
// Java conditions ported here:
//   - The reference is in the same container scope as the declaration
//     (hasSameContainerScope). We approximate this: the name must be
//     in the current scope's earlyRefStack entry (pre-scanned vars),
//     and we must not be inside a nested function (depth > the depth
//     where the var was pre-scanned).
//   - For `var` declarations: the self-referential `var x = x || {}`
//     pattern is exempt (Java walks up through OR nodes to detect it).
//     The `??` (COALESCE) form is NOT exempt — it fires the warning.
//   - The reference must come before the declaration position OR be
//     inside the initializer of the same-line declaration
//     (e.g. `var x = x ?? {}` — the reference is in the RHS).
func (v *VariableReferenceCheck) checkEarlyReference(n, parent *rhino.Node) {
	name := n.String_()
	if name == "" {
		return
	}
	// Only check NAME nodes that are value references, not binding sites.
	if isBindingSiteVRC(n, parent) {
		return
	}
	// Look up the pre-scanned var declaration in the current
	// container scope only (the top of earlyRefStack).
	decl, found := v.earlyRefStack[len(v.earlyRefStack)-1][name]
	if !found {
		return
	}
	refLine, refCol := n.Lineno(), n.Charno()
	// Fire if the reference textually precedes the declaration.
	if refLine > decl.line {
		return
	}
	if refLine == decl.line && refCol >= decl.col {
		// Same line, reference col >= decl col: only fire when the
		// reference is inside the var's own initializer expression
		// (e.g. `var x = x ?? {}`). Java handles this naturally via
		// post-order traversal; we detect it by walking up the
		// parent chain to find the declaration NAME node.
		if !isInsideVarInitializer(n, name) {
			return
		}
	}
	// Exempt the `var x = x || {}` self-init pattern. Java walks
	// up through OR-as-LHS chains and checks whether the eventual
	// parent is the declaration NAME itself. We replicate the same
	// traversal here.
	if isVarSelfInitOr(n, name) {
		return
	}
	v.diags = append(v.diags, diag.Diagnostic{
		Path:    v.path,
		Line:    refLine,
		Col:     refCol,
		Code:    diag.JSC_REFERENCE_BEFORE_DECLARE,
		Level:   diag.Warning,
		Message: fmt.Sprintf("Variable referenced before declaration: %s", name),
	})
}

// isInsideVarInitializer reports whether NAME reference node n is
// inside the initializer expression of a var declaration for the same
// variable name. This detects the `var x = f(x)` family where the
// reference in the RHS is semantically before initialization.
// Mirrors the intuition behind Java's post-order reference collection:
// in post-order, the initializer expression is visited before the
// declaration NAME, so hasSeenDeclaration is false.
func isInsideVarInitializer(n *rhino.Node, name string) bool {
	curr := n.Parent()
	for curr != nil {
		if curr.Token() == rhino.NAME && curr.String_() == name {
			if p := curr.Parent(); p != nil && p.Token() == rhino.VAR {
				return true
			}
		}
		curr = curr.Parent()
	}
	return false
}

// isVarSelfInitOr returns true for the `var x = x || {}` pattern where
// the reference `n` (a NAME node with the given name) sits on the LHS
// of a chain of OR expressions whose ultimate parent is the declaration
// NAME node for the same variable. This mirrors Java's special case in
// checkEarlyReference:
//
//	Node curr = reference.getParent();
//	while (curr.isOr() && curr.getParent().getFirstChild() == curr) { curr = curr.getParent(); }
//	if (curr.isName() && curr.getString().equals(v.getName())) { return false; }
func isVarSelfInitOr(n *rhino.Node, varName string) bool {
	curr := n.Parent()
	for curr != nil && curr.Token() == rhino.OR {
		p := curr.Parent()
		if p == nil || p.FirstChild() != curr {
			break
		}
		curr = p
	}
	if curr != nil && curr.Token() == rhino.NAME && curr.String_() == varName {
		return true
	}
	return false
}

// isBindingSiteVRC classifies NAME nodes that are declarations (param,
// VAR/LET/CONST child, function/class name, class-field-def name,
// STRING_KEY child) rather than value references. Mirrors the logic in
// VarCheck.isBindingSite, kept separate to avoid a package-level
// coupling. The name "_VRC" suffix distinguishes it from VarCheck's
// copy; both are 1:1 ports of the same NodeUtil helpers.
func isBindingSiteVRC(n, parent *rhino.Node) bool {
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
	case rhino.GETPROP:
		// Property access: obj.prop — the second child is a property name,
		// not a variable reference.
		kids := parent.Children()
		return len(kids) >= 2 && kids[1] == n
	}
	return false
}

func (v *VariableReferenceCheck) recordDecl(nameNode *rhino.Node, name string) {
	v.recordDeclInto(v.scopes[len(v.scopes)-1], nameNode, name)
}

// isMemberContext reports whether `parent` (the node above a FUNCTION /
// CLASS) is a class-body or object-literal slot, where the nested
// definition is a member rather than a var-like declaration.
func isMemberContext(parent *rhino.Node) bool {
	if parent == nil {
		return false
	}
	// Class methods live under the CLASS's BLOCK body (the third
	// child of a CLASS node).
	if parent.Token() == rhino.BLOCK && parent.Parent() != nil && parent.Parent().Token() == rhino.CLASS {
		return true
	}
	// Object-literal method shorthand wraps the FUNCTION under a
	// STRING_KEY.
	if parent.Token() == rhino.STRING_KEY {
		return true
	}
	return false
}

func (v *VariableReferenceCheck) recordDeclInto(scope map[string]declSite, nameNode *rhino.Node, name string) {
	if name == "" {
		return
	}
	if prev, seen := scope[name]; seen {
		if v.suppressDepth > 0 {
			return
		}
		v.diags = append(v.diags, diag.Diagnostic{
			Path:  v.path,
			Line:  nameNode.Lineno(),
			Col:   nameNode.Charno(),
			Code:  diag.JSC_VAR_MULTIPLY_DECLARED_ERROR,
			Level: diag.Error,
			Message: fmt.Sprintf("Variable %s declared more than once. First occurrence: %s:%d:%d",
				name, v.path, prev.line, prev.col),
		})
		return
	}
	scope[name] = declSite{line: nameNode.Lineno(), col: nameNode.Charno()}
}
