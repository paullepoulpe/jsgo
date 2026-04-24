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
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
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
	line   int
	col    int
	isFunc bool // true when declared by a function/class declaration (not var/let/const)
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
	// A file-level `@fileoverview @suppress {duplicate}` suppresses all
	// duplicate-declaration diagnostics for the entire script. Mirrors
	// Java's hasDuplicateDeclarationSuppression / CheckLevel.OFF gate.
	// The parser propagates @fileoverview suppressions into the first
	// statement's JSDoc via mergeSuppressions; we also walk the whole
	// script to catch any node that carries the suppress annotation.
	if v.hasFileSuppressDuplicate(script) {
		v.suppressDepth++
		defer func() { v.suppressDepth-- }()
	}
	// Pre-scan script-level var declarations for the early-reference check.
	v.preScanVarDecls(script, v.earlyRefStack[0])
	v.walk(script, nil)
	return v.diags
}

// Go-only: hasFileSuppressDuplicate reports whether any node in the script
// carries a @suppress {duplicate} annotation. This catches both the
// SCRIPT node itself (if the parser attaches @fileoverview there)
// and the first-statement node (where mergeSuppressions propagates
// the @fileoverview suppress). Mirrors Java's
// collectFileSuppressions / CompilerPass.skipForSuppressedPath logic.
func (v *VariableReferenceCheck) hasFileSuppressDuplicate(script *rhino.Node) bool {
	var check func(n *rhino.Node) bool
	check = func(n *rhino.Node) bool {
		if n == nil {
			return false
		}
		if info := n.JSDocInfo(); info != nil && info.HasSuppression("duplicate") {
			return true
		}
		for _, c := range n.Children() {
			if check(c) {
				return true
			}
		}
		return false
	}
	return check(script)
}

// Go-only: preScanVarDecls collects every `var`-keyword declaration in `n`'s
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

// Go-only: walk is the AST traversal driver. Java uses NodeTraversal callbacks
// (shouldTraverse / visit); Go uses a direct recursive walk function.
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
					v.recordFuncDecl(target, name, name.String_())
				} else {
					v.recordDecl(name, name.String_())
				}
			}
		} else if n.Token() == rhino.FUNCTION {
			// Named function expression (parent is not a block-level
			// statement). Check if the function name bleeds into the
			// enclosing scope and masks an existing non-function variable.
			// Mirrors TypeCheck.java shouldTraverse's FUNCTION_MASKS_VARIABLE
			// check (IE bleeding-name bug detection).
			v.checkFunctionMasksVariable(n)
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

// Go-only: isInsideVarInitializer reports whether NAME reference node n is
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

// Go-only: isVarSelfInitOr returns true for the `var x = x || {}` pattern where
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

// Go-only: isBindingSiteVRC classifies NAME nodes that are declarations (param,
// VAR/LET/CONST child, function/class name, class-field-def name,
// STRING_KEY child inside destructuring, ARRAY_PATTERN child, etc.)
// rather than value references. Mirrors the logic in VarCheck.isBindingSite,
// kept separate to avoid a package-level coupling. Both are 1:1 ports of
// the same NodeUtil helpers extended for ES2015+ destructuring patterns.
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

// Go-only: recordDecl dispatches to recordDeclInto for the current scope.
func (v *VariableReferenceCheck) recordDecl(nameNode *rhino.Node, name string) {
	v.recordDeclInto(v.scopes[len(v.scopes)-1], nameNode, name, false)
}

// Go-only: recordFuncDecl records a function declaration into a specific scope map.
func (v *VariableReferenceCheck) recordFuncDecl(scope map[string]declSite, nameNode *rhino.Node, name string) {
	v.recordDeclInto(scope, nameNode, name, true)
}

// Go-only: isMemberContext reports whether `parent` (the node above a FUNCTION /
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

// Go-only: recordDeclInto records a declaration into the given scope map, emitting
// JSC_VAR_MULTIPLY_DECLARED_ERROR on conflict. Java's equivalent is
// ReferenceCheckingBehavior.checkRedeclaration.
func (v *VariableReferenceCheck) recordDeclInto(scope map[string]declSite, nameNode *rhino.Node, name string, isFunc bool) {
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
	scope[name] = declSite{line: nameNode.Lineno(), col: nameNode.Charno(), isFunc: isFunc}
}

// Go-only: checkFunctionMasksVariable mirrors TypeCheck.java's shouldTraverse
// FUNCTION_MASKS_VARIABLE path. Fires when a named function expression
// (not a function declaration) has the same name as an existing
// non-function variable in the enclosing scope — the so-called IE
// "bleeding name" bug where the named function expression's name leaks
// into the outer scope under IE9 quirks mode.
//
// Java condition (TypeCheck.shouldTraverse):
//   var = outerScope.getVar(funcNameString)
//   var != null
//     && var.getScope().hasSameContainerScope(outerScope)
//     && !(var.getType() instanceof FunctionType)
//     && !hasDuplicateDeclarationSuppression(var.getNameNode())
//
// We approximate !(var.getType() instanceof FunctionType) by checking
// whether the first declaration of the name was recorded as isFunc=false
// (i.e. it came from a var/let/const, not a function/class declaration).
func (v *VariableReferenceCheck) checkFunctionMasksVariable(fn *rhino.Node) {
	kids := fn.Children()
	if len(kids) < 1 {
		return
	}
	funcName := kids[0].String_()
	if funcName == "" {
		return // anonymous function expression — no bleeding
	}
	// The enclosing scope is scopes[-2] because we already pushed the
	// function's own scope frame when entering the FUNCTION node.
	if len(v.scopes) < 2 {
		return
	}
	enclosing := v.scopes[len(v.scopes)-2]
	prev, found := enclosing[funcName]
	if !found {
		return
	}
	// Only fire when the existing declaration is NOT a function/class
	// declaration (mirrors Java's !(var.getType() instanceof FunctionType)).
	if prev.isFunc {
		return
	}
	if v.suppressDepth > 0 {
		return
	}
	v.diags = append(v.diags, diag.Diagnostic{
		Path:    v.path,
		Line:    fn.Lineno(),
		Col:     fn.Charno(),
		Code:    diag.JSC_FUNCTION_MASKS_VARIABLE,
		Level:   diag.Error, // upgraded to error by --jscomp_error=checkTypes (DiagnosticGroups TYPE_CHECK)
		Message: fmt.Sprintf("function %s masks variable (IE bug)", funcName),
	})
}

// ---------------------------------------------------------------------------
// CheckDupVarDecl — inlined from checkdupvardecl.go.
// Port of com/google/javascript/jscomp/TypeCheck.java's
// checkDuplicateVarDeclaration branch @ closure-compiler v20260415
// (narrow; namespace-property slot only).
//
// Fires JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH when a namespace
// property is re-declared with a different type. Java tracks two
// shapes:
//
//	/** @type {T1} */ goog.foo;
//	/** @type {T2} */ goog.foo;
//
// and the assignment-form mirror. jsgo records every such decl
// on the scope's namespaceProps map; this pass walks the AST a
// second time, flagging each subsequent declaration against the
// first one's type.
// ---------------------------------------------------------------------------

// CheckDupVarDecl emits JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH
// when a namespace property's @type is redeclared with a
// different type. Runs after TypedScopeCreator has populated
// the scope's namespaceProps map with the *final* type — the
// diagnostic walks the AST to find each declaration's position
// and compares the text-level @type against the stored value.
type CheckDupVarDecl struct {
	path  string
	scope *TypedScope
	diags []diag.Diagnostic
}

func NewCheckDupVarDecl(path string, scope *TypedScope) *CheckDupVarDecl {
	return &CheckDupVarDecl{path: path, scope: scope}
}

// Check walks the script collecting every
// `/** @type {T} */ owner.prop[ = … ]` pattern. For each
// (owner, prop) pair, records the FIRST declaration and compares
// every subsequent one — firing when the resolved types don't
// match.
func (c *CheckDupVarDecl) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	type decl struct {
		at *rhino.Node
		t  jstype.JSType
	}
	seen := map[string]map[string]*decl{}
	c.walk(script, func(owner, prop string, at *rhino.Node, t jstype.JSType) {
		if seen[owner] == nil {
			seen[owner] = map[string]*decl{}
		}
		first := seen[owner][prop]
		if first == nil {
			seen[owner][prop] = &decl{at: at, t: t}
			return
		}
		if sameDeclaredType(t, first.t) {
			return
		}
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: at.Lineno(), Col: at.Charno(),
			Code: diag.JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH, Level: diag.Error,
			Message: fmt.Sprintf(
				"variable %s.%s redefined with type %s, original definition at %s:%d with type %s",
				owner, prop, t.TypeName(),
				c.path, first.at.Lineno(), first.t.TypeName()),
		})
	})
	// Pass 2: function-assign-to-namespace — fires
	// JSC_DUP_VAR_DECLARATION for the
	// `ns.foo = function(){}; ns.foo = function(){}` shape when
	// any subsequent assign carries a declaration JSDoc (@param,
	// @return, @constructor, @interface, @record) and isn't
	// @suppress {duplicate}.
	type decl2 struct{ at *rhino.Node }
	fnSeen := map[string]map[string]*decl2{}
	c.walkFunctionAssigns(script, func(owner, prop string, at *rhino.Node, fresh bool) {
		if fnSeen[owner] == nil {
			fnSeen[owner] = map[string]*decl2{}
		}
		first := fnSeen[owner][prop]
		if first == nil {
			fnSeen[owner][prop] = &decl2{at: at}
			return
		}
		if !fresh {
			return
		}
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: at.Lineno(), Col: at.Charno(),
			Code: diag.JSC_DUP_VAR_DECLARATION, Level: diag.Error,
			Message: fmt.Sprintf(
				"variable %s.%s redefined, original definition at %s:%d",
				owner, prop, c.path, first.at.Lineno()),
		})
	})
	// Pass 3: @typedef stubs that redeclare a previous @constructor
	// function-assign. Java fires DUP_VAR_DECLARATION_TYPE_MISMATCH from
	// TypedScopeCreator.expectUndeclaredVariable when the typedef's "None"
	// type clashes with the constructor's FunctionType. Mirrors that shape
	// without needing access to the resolved types.
	c.walkTypedefStubs(script, func(owner, prop string, at *rhino.Node) {
		if fnSeen[owner] == nil || fnSeen[owner][prop] == nil {
			return
		}
		first := fnSeen[owner][prop]
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: at.Lineno(), Col: at.Charno(),
			Code: diag.JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH, Level: diag.Error,
			Message: fmt.Sprintf(
				"variable %s.%s redefined with type None, original definition at %s:%d with type (typeof %s.%s)",
				owner, prop, c.path, first.at.Lineno(), owner, prop),
		})
	})
	// Pass 4: @template-vs-@param name collisions inside function bodies.
	// Java's TypedScopeCreator fires DUP_VAR_DECLARATION_TYPE_MISMATCH via
	// expectUndeclaredVariable when it tries to declare a @template name T
	// into the function scope but finds a same-named @param already there.
	// The diagnostic points at the FUNCTION node (col 0 of the fn keyword).
	c.walkTemplateParmCollisions(script, func(fn *rhino.Node, name string, paramTypeName string) {
		c.diags = append(c.diags, diag.Diagnostic{
			Path: c.path, Line: fn.Lineno(), Col: fn.Charno(),
			Code: diag.JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH, Level: diag.Error,
			Message: fmt.Sprintf(
				"variable %s redefined with type undefined, original definition at %s:%d with type %s",
				name, c.path, fn.Lineno(), paramTypeName),
		})
	})
	return c.diags
}

// Go-only: sameDeclaredType reports whether two JSTypes declared on the
// same slot should be treated as equal by TypeCheck's duplicate-
// declaration check. Java compares by nominal identity: two
// @record-backed ObjectTypes with the same structural surface
// are still different types. Primitives compare by value.
func sameDeclaredType(a, b jstype.JSType) bool {
	if a == nil || b == nil {
		return a == b
	}
	if oa, ok := a.(*jstype.ObjectType); ok {
		// Named (nominal) ObjectTypes compare by pointer identity —
		// two different @constructor types are always unequal even if
		// they have the same surface. Anonymous inline record types
		// ({x: number}) are equal when their rendered TypeName matches,
		// mirroring Java's equals() for anonymous structural types.
		if oa.Name == "" {
			if ob, ok2 := b.(*jstype.ObjectType); ok2 && ob.Name == "" {
				return a.TypeName() == b.TypeName()
			}
		}
		return a == b
	}
	if _, ok := a.(*jstype.FunctionType); ok {
		return a == b
	}
	if _, ok := a.(*jstype.EnumType); ok {
		return a == b
	}
	if _, ok := a.(*jstype.TemplatizedType); ok {
		return a == b
	}
	// Primitives / unions: round-trip subtype test is sufficient.
	return jstype.IsSubtype(a, b) && jstype.IsSubtype(b, a)
}

// Go-only: walkFunctionAssigns collects every `owner.prop = function(){…}`
// assignment whose JSDoc signals a fresh function declaration
// (@param / @return / @constructor / @interface / @record),
// skipping @suppress {duplicate}. Mirrors Java's
// TypeCheck.checkDuplicateVarDeclaration firing shape.
func (c *CheckDupVarDecl) walkFunctionAssigns(n *rhino.Node, emit func(owner, prop string, at *rhino.Node, fresh bool)) {
	if n == nil {
		return
	}
	if n.Token() == rhino.EXPR_RESULT {
		if inner := n.FirstChild(); inner != nil && inner.Token() == rhino.ASSIGN {
			kids := inner.Children()
			if len(kids) >= 2 {
				lhs, rhs := kids[0], kids[1]
				if rhs.Token() == rhino.FUNCTION && lhs.Token() == rhino.GETPROP {
					lhsKids := lhs.Children()
					if len(lhsKids) >= 2 {
						recv := lhsKids[0]
						if recv.Token() == rhino.NAME {
							if _, ok := c.scope.GetDeclaredType(recv.String_()).(*jstype.FunctionType); !ok {
								propNode := lhsKids[1]
								emit(recv.String_(), propNode.String_(), propNode, isFreshFunctionDecl(n, inner))
							}
						}
					}
				}
			}
		}
	}
	for _, ch := range n.Children() {
		c.walkFunctionAssigns(ch, emit)
	}
}

// Go-only: isFreshFunctionDecl reports whether the ASSIGN carries JSDoc
// Java considers a function-declaration marker (as opposed to a
// mere type annotation). Also returns false when the author
// explicitly suppressed `duplicate`.
func isFreshFunctionDecl(stmt, a *rhino.Node) bool {
	info := a.JSDocInfo()
	if info == nil {
		info = stmt.JSDocInfo()
	}
	if info == nil {
		return false
	}
	if info.HasSuppression("duplicate") {
		return false
	}
	if info.IsConstructor() || info.IsInterface() || info.IsRecord() {
		return true
	}
	if info.ReturnType() != nil {
		return true
	}
	if len(info.ParameterNames()) > 0 {
		return true
	}
	return false
}

// Go-only: walkTypedefStubs collects every `/** @typedef {T} */ owner.prop;`
// GETPROP stub (EXPR_RESULT > GETPROP with JSDoc @typedef). Used to
// detect when a @typedef re-declares a name that was previously a
// @constructor function assign. Mirrors TypedScopeCreator's
// expectUndeclaredVariable path when newType is the typedef "None".
func (c *CheckDupVarDecl) walkTypedefStubs(n *rhino.Node, emit func(owner, prop string, at *rhino.Node)) {
	if n == nil {
		return
	}
	if n.Token() == rhino.EXPR_RESULT {
		if inner := n.FirstChild(); inner != nil && inner.Token() == rhino.GETPROP {
			info := n.JSDocInfo()
			if info == nil {
				info = inner.JSDocInfo()
			}
			if info != nil && info.TypedefType() != nil {
				kids := inner.Children()
				if len(kids) >= 2 {
					recv := kids[0]
					if recv.Token() == rhino.NAME {
						propNode := kids[1]
						emit(recv.String_(), propNode.String_(), propNode)
					}
				}
			}
		}
	}
	for _, ch := range n.Children() {
		c.walkTypedefStubs(ch, emit)
	}
}

// Go-only: walkTemplateParmCollisions detects functions whose @template names
// collide with a parameter name in the same function. Mirrors
// TypedScopeCreator's declareTemplateTypeNames loop which calls
// expectUndeclaredVariable (firing DUP_VAR_DECLARATION_TYPE_MISMATCH)
// when `!currentScope.canDeclare(name)`.
//
// emit is called with the FUNCTION node, the colliding name, and the
// @param's declared type name (for the error message).
func (c *CheckDupVarDecl) walkTemplateParmCollisions(n *rhino.Node, emit func(fn *rhino.Node, name string, paramTypeName string)) {
	if n == nil {
		return
	}
	if n.Token() == rhino.FUNCTION {
		info := n.JSDocInfo()
		tnames := info.TemplateNames()
		if len(tnames) > 0 {
			kids := n.Children()
			if len(kids) >= 2 {
				params := kids[1]
				paramSet := map[string]bool{}
				for _, p := range params.Children() {
					paramSet[p.String_()] = true
				}
				for _, tn := range tnames {
					if paramSet[tn] {
						// Determine the @param's declared type name for the message.
						paramTypeName := "unknown"
						if pt := info.ParameterType(tn); pt != nil {
							paramTypeName = pt.TypeName
						}
						emit(n, tn, paramTypeName)
					}
				}
			}
		}
	}
	for _, ch := range n.Children() {
		c.walkTemplateParmCollisions(ch, emit)
	}
}

// Go-only: walk recurses over the AST calling `emit` for every namespace-
// property type declaration it finds. Java uses NodeTraversal.
func (c *CheckDupVarDecl) walk(n *rhino.Node, emit func(owner, prop string, at *rhino.Node, t jstype.JSType)) {
	if n == nil {
		return
	}
	if n.Token() == rhino.EXPR_RESULT {
		inner := n.FirstChild()
		if inner != nil && inner.Token() == rhino.ASSIGN {
			if owner, prop, at, t := c.declFromAssign(n, inner); owner != "" {
				emit(owner, prop, at, t)
			}
		}
		if inner != nil && inner.Token() == rhino.GETPROP {
			if owner, prop, at, t := c.declFromStub(n, inner); owner != "" {
				emit(owner, prop, at, t)
			}
		}
	}
	for _, ch := range n.Children() {
		c.walk(ch, emit)
	}
}

// Go-only: declFromAssign extracts owner/prop/type from an assignment expression.
func (c *CheckDupVarDecl) declFromAssign(stmt, a *rhino.Node) (string, string, *rhino.Node, jstype.JSType) {
	info := a.JSDocInfo()
	if info == nil {
		info = stmt.JSDocInfo()
	}
	if info == nil || !info.HasType() {
		return "", "", nil, nil
	}
	kids := a.Children()
	if len(kids) < 2 {
		return "", "", nil, nil
	}
	lhs := kids[0]
	if lhs.Token() != rhino.GETPROP {
		return "", "", nil, nil
	}
	lhsKids := lhs.Children()
	if len(lhsKids) < 2 {
		return "", "", nil, nil
	}
	recv := lhsKids[0]
	if recv.Token() != rhino.NAME {
		return "", "", nil, nil
	}
	ownerType := c.scope.GetDeclaredType(recv.String_())
	// Skip when the owner is undeclared (nil) — undeclared owners
	// can't have typed namespace properties. Mirrors Java's
	// expectUndeclaredVariable gate: the variable must exist in
	// scope for the duplicate check to trigger.
	if ownerType == nil {
		return "", "", nil, nil
	}
	if _, ok := ownerType.(*jstype.FunctionType); ok {
		return "", "", nil, nil
	}
	propNode := lhsKids[1]
	return recv.String_(), propNode.String_(), propNode, c.scope.ResolveExpr(info.Type().TypeName)
}

// Go-only: declFromStub extracts owner/prop/type from a stub expression.
func (c *CheckDupVarDecl) declFromStub(stmt, getprop *rhino.Node) (string, string, *rhino.Node, jstype.JSType) {
	info := stmt.JSDocInfo()
	if info == nil {
		info = getprop.JSDocInfo()
	}
	if info == nil || !info.HasType() {
		return "", "", nil, nil
	}
	kids := getprop.Children()
	if len(kids) < 2 {
		return "", "", nil, nil
	}
	recv := kids[0]
	if recv.Token() != rhino.NAME {
		return "", "", nil, nil
	}
	ownerType := c.scope.GetDeclaredType(recv.String_())
	// Skip when the owner is undeclared — mirrors Java's scope guard.
	if ownerType == nil {
		return "", "", nil, nil
	}
	if _, ok := ownerType.(*jstype.FunctionType); ok {
		return "", "", nil, nil
	}
	propNode := kids[1]
	return recv.String_(), propNode.String_(), propNode, c.scope.ResolveExpr(info.Type().TypeName)
}
