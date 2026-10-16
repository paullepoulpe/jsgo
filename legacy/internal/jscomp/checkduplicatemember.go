// Port of com/google/javascript/jscomp/StrictModeCheck.java @ closure-compiler v20260415
//
// This file ports the duplicate-member subset of StrictModeCheck:
//
//  Object literals: fires JSC_DUPLICATE_MEMBER + JSC_DUP_VAR_DECLARATION
//  when the same key appears twice in an OBJECTLIT — `{BB: 1, BB: 2}`.
//  Java emits a WARNING at the first occurrence and an ERROR at the second.
//
//  Class bodies: fires JSC_DUPLICATE_MEMBER (WARNING only) when a
//  class declares two members with the same name, mirroring Java's
//  StrictModeCheck.checkObjectLiteralOrClass. Getter + setter for the
//  same name are compatible and do not trigger the diagnostic. Static
//  and instance slots are tracked separately.
//
// The NonExternChecks inner class (eval/arguments redeclaration, callee/caller
// property access) is not yet ported because it requires scope infrastructure.
// TODO: port USE_OF_WITH, EVAL_DECLARATION, EVAL_ASSIGNMENT, ARGUMENTS_DECLARATION,
//   ARGUMENTS_ASSIGNMENT, ARGUMENTS_CALLEE_FORBIDDEN, ARGUMENTS_CALLER_FORBIDDEN,
//   FUNCTION_CALLER_FORBIDDEN, FUNCTION_ARGUMENTS_PROP_FORBIDDEN, DELETE_VARIABLE.
// TODO: port NonExternChecks.checkDeclaration, NonExternChecks.checkGetProp.

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/legacy/internal/diag"
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
)

// CheckDuplicateMember mirrors the duplicate-member checking subset of
// Java's StrictModeCheck class (StrictModeCheck extends AbstractPostOrderCallback
// implements CompilerPass).
type CheckDuplicateMember struct {
	path  string
	diags []diag.Diagnostic
}

// NewCheckDuplicateMember mirrors the Java constructor
// StrictModeCheck(AbstractCompiler compiler, CheckLevel defaultLevel).
// Go-only: language-adapter — Go constructor shim (Java uses class constructor directly; Go has no class constructors).
func NewCheckDuplicateMember(path string) *CheckDuplicateMember {
	return &CheckDuplicateMember{path: path}
}

// Check walks the AST flagging duplicate members in every OBJECTLIT
// and every class body (CLASS → BLOCK).
// Mirrors Java's process(Node externs, Node root): calls
// NodeTraversal.traverseRoots and NodeTraversal.traverse(compiler, root, new NonExternChecks()).
func (c *CheckDuplicateMember) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.walk(script, "")
	return c.diags
}

// walk drives the traversal, calling visit() for each node and recursing into children.
// Pre-order walk mirrors the order of Java's AbstractPostOrderCallback for this pass.
// TODO(port): NodeTraversal.traverseRoots(AbstractCompiler, Node, Callback, Node...) — NodeTraversal adapter not yet ported.
func (c *CheckDuplicateMember) walk(n *rhino.Node, owner string) {
	if n == nil {
		return
	}
	c.visit(n, owner)
	for _, ch := range n.Children() {
		c.walk(ch, owner)
	}
}

// visit checks a single node. Mirrors Java's StrictModeCheck.visit(NodeTraversal t, Node n, Node parent):
//
//	if (n.isAssign())          checkAssignment(n)
//	else if (n.isDelProp())   checkDelete(t, n)
//	else if (n.isObjectLit()) checkObjectLiteralOrClass(n)
//	else if (n.isClass())     checkObjectLiteralOrClass(n.getLastChild())
//	else if (n.isWith())      checkWith(n)
func (c *CheckDuplicateMember) visit(n *rhino.Node, owner string) {
	switch n.Token() {
	case rhino.OBJECTLIT:
		c.checkObjectLiteralOrClass(n, owner)
	case rhino.CLASS:
		// Java: checkObjectLiteralOrClass(n.getLastChild())
		// CLASS children: [0]=name, [1]=superName, [2]=body BLOCK.
		if kids := n.Children(); len(kids) >= 3 {
			c.checkObjectLiteralOrClass(kids[2], owner)
		}
	}
}

// checkObjectLiteralOrClass ports StrictModeCheck.checkObjectLiteralOrClass(Node n).
// Checks both object literals and class bodies (n is the body BLOCK for classes).
// Uses four sets (getters / setters × instance / static) to track names,
// iterating backwards for class bodies so the last duplicate "wins" (sloppy-mode semantics).
//
// Java (lines 163-202):
//
//	private void checkObjectLiteralOrClass(Node n) {
//	  LinkedHashSet<String> getters = new LinkedHashSet<>();
//	  LinkedHashSet<String> setters = new LinkedHashSet<>();
//	  LinkedHashSet<String> staticGetters = new LinkedHashSet<>();
//	  LinkedHashSet<String> staticSetters = new LinkedHashSet<>();
//	  for (Node key = n.getLastChild(); key != null; key = key.getPrevious()) {
//	    if (key.isEmpty() || key.isComputedProp() || key.isSpread() || …) continue;
//	    String keyName = key.getString();
//	    if (!key.isSetterDef()) { set = …getters; if (!set.add(keyName)) report(key, DUPLICATE_MEMBER, keyName); }
//	    if (!key.isGetterDef()) { set = …setters; if (!set.add(keyName)) report(key, DUPLICATE_MEMBER, keyName); }
//	  }
//	}
func (c *CheckDuplicateMember) checkObjectLiteralOrClass(n *rhino.Node, owner string) {
	kids := n.Children()

	if n.Token() == rhino.OBJECTLIT {
		// Object literal: check STRING_KEY children for duplicates.
		// Iterate forward, tracking first-seen by name; fires diagnostics
		// at first (JSC_DUPLICATE_MEMBER) and second (JSC_DUP_VAR_DECLARATION) occurrences.
		scope := enclosingOwnerName(n, owner)
		seen := map[string]*rhino.Node{}
		for _, kv := range kids {
			if kv.Token() != rhino.STRING_KEY {
				continue
			}
			name := kv.String_()
			if first, ok := seen[name]; ok {
				c.diags = append(c.diags,
					diag.Diagnostic{
						Path: c.path, Line: first.Lineno(), Col: first.Charno(),
						Code:  diag.JSC_DUPLICATE_MEMBER,
						Level: diag.Warning,
						Message: fmt.Sprintf(
							"Class or object literal contains duplicate member %q. In non-strict code, the last duplicate will overwrite the others.",
							name),
					},
					diag.Diagnostic{
						Path: c.path, Line: kv.Lineno(), Col: kv.Charno(),
						Code:  diag.JSC_DUP_VAR_DECLARATION,
						Level: diag.Error,
						Message: fmt.Sprintf(
							"variable %s.%s redefined, original definition at %s:%d",
							scope, name, c.path, first.Lineno()),
					},
				)
				continue
			}
			seen[name] = kv
		}
		return
	}

	// Class body: iterate backwards (Java: for (Node key = n.getLastChild(); ...));
	// skip computed, empty, static-block nodes.
	// Four independent name-sets track instance-getter, instance-setter,
	// static-getter, and static-setter slots — matching the Java 4-set design.
	// A getter and a setter with the same name are compatible and do NOT trigger
	// the diagnostic; two getters (or two non-setter members) DO.
	getters := map[string]*rhino.Node{}
	setters := map[string]*rhino.Node{}
	staticGetters := map[string]*rhino.Node{}
	staticSetters := map[string]*rhino.Node{}

	for i := len(kids) - 1; i >= 0; i-- {
		m := kids[i]
		switch m.Token() {
		case rhino.FUNCTION, rhino.CLASS_FIELD_DEF:
			// proceed
		default:
			continue // skip COMPUTED_PROP, EMPTY, etc.
		}
		first := m.FirstChild()
		if first == nil {
			continue
		}
		name := first.String_()
		if name == "" {
			continue
		}
		isStatic := m.JSDocInfo().IsStatic()
		isGetter := m.IsGetterDef()
		isSetter := m.IsSetterDef()

		if !isSetter {
			gs := getters
			if isStatic {
				gs = staticGetters
			}
			if _, dup := gs[name]; dup {
				c.report(m, diag.JSC_DUPLICATE_MEMBER, diag.Warning,
					fmt.Sprintf("Class or object literal contains duplicate member %q. In non-strict code, the last duplicate will overwrite the others.", name))
			} else {
				gs[name] = m
			}
		}
		if !isGetter {
			ss := setters
			if isStatic {
				ss = staticSetters
			}
			if _, dup := ss[name]; dup {
				c.report(m, diag.JSC_DUPLICATE_MEMBER, diag.Warning,
					fmt.Sprintf("Class or object literal contains duplicate member %q. In non-strict code, the last duplicate will overwrite the others.", name))
			} else {
				ss[name] = m
			}
		}
	}
}

// report emits a diagnostic. Mirrors Java's
// void report(Node n, DiagnosticType diagnostic, String... args).
func (c *CheckDuplicateMember) report(n *rhino.Node, code string, level diag.Level, msg string) {
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    n.Lineno(),
		Col:     n.Charno(),
		Code:    code,
		Level:   level,
		Message: msg,
	})
}

// enclosingOwnerName walks parents looking for the VAR / ASSIGN
// LHS that owns this OBJECTLIT, returning the dotted name used to
// qualify the JSC_DUP_VAR_DECLARATION diagnostic. Falls back to
// `<anonymous>` when the literal isn't directly assigned to a
// named slot.
// Go-only: language-adapter — helper for JSC_DUP_VAR_DECLARATION message formatting (Java uses NodeUtil.getEnclosingName inline; Go extracts for clarity).
func enclosingOwnerName(n *rhino.Node, fallback string) string {
	for cur := n.Parent(); cur != nil; cur = cur.Parent() {
		switch cur.Token() {
		case rhino.NAME:
			// VAR name — treat as the owner when VAR-initialised.
			if p := cur.Parent(); p != nil {
				switch p.Token() {
				case rhino.VAR, rhino.LET, rhino.CONST:
					return cur.String_()
				}
			}
		case rhino.ASSIGN:
			if kids := cur.Children(); len(kids) >= 1 {
				return flattenGetProp(kids[0])
			}
		}
	}
	if fallback != "" {
		return fallback
	}
	return "<anonymous>"
}

// flattenGetProp renders a NAME / GETPROP chain as a dotted path
// (`goog.a.BB`). Unknown tokens fall through to `<anonymous>`.
// Go-only: language-adapter — helper for JSC_DUP_VAR_DECLARATION message formatting (Java uses NodeUtil.getEnclosingName inline; Go extracts for clarity).
func flattenGetProp(n *rhino.Node) string {
	if n == nil {
		return "<anonymous>"
	}
	switch n.Token() {
	case rhino.NAME:
		return n.String_()
	case rhino.GETPROP:
		kids := n.Children()
		if len(kids) >= 2 {
			return flattenGetProp(kids[0]) + "." + kids[1].String_()
		}
	}
	return "<anonymous>"
}
