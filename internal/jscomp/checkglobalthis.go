// Port of com/google/javascript/jscomp/CheckGlobalThis.java @ closure-compiler v20260415
//
// Emits JSC_USED_GLOBAL_THIS for `this` references that resolve to the
// global object: script top-level, non-@constructor / non-@this
// function bodies that can't carry a @this annotation, and arrow
// functions inheriting the global `this`.
//
// A `this` is reported only when it's observed — on the LHS of an
// assignment or as the receiver of a property access — matching the
// Java class's `shouldReportThis` narrowing.
//
// Structural skips (mirroring the Java ShouldTraverse branches):
//
//   - class static initialization blocks,
//   - class member / computed fields,
//   - functions annotated with @constructor / @interface / @this /
//     @override, or whose declared @type pins `function(this:T,…)`,
//   - object-literal methods being lent to `Foo.prototype` via @lends,
//   - the RHS of a prototype-slot assignment
//     (`Foo.prototype.m = function() {…}`).
//
// The pass deliberately does _not_ follow `this` through aliases —
// `var a = this; a.x = 1;` is silently accepted, matching Java.

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckGlobalThis is the jsgo port of Java's final class of the same
// name. The walker runs ShouldTraverse / Visit directly — no external
// NodeTraversal driver is needed since CheckGlobalThis never consults
// the scope stack.
type CheckGlobalThis struct {
	path  string
	diags []diag.Diagnostic

	// assignLhsChild, when non-nil, is the first child of the ASSIGN
	// whose subtree we're currently inside. Nested assigns must not
	// overwrite it so `(a = this).b = c` still reports the outer
	// `this`. Mirrors Java's field of the same name.
	assignLhsChild *rhino.Node
}

func NewCheckGlobalThis(path string) *CheckGlobalThis {
	return &CheckGlobalThis{path: path}
}

func (c *CheckGlobalThis) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.assignLhsChild = nil
	c.traverse(script, nil)
	return c.diags
}

// traverse walks `n` with the ShouldTraverse / Visit dispatch order
// from NodeTraversal: pre-order guard, then children, then post-order
// visit. Matches the Java driver's semantics without pulling in the
// full NodeTraversal machinery.
func (c *CheckGlobalThis) traverse(n, parent *rhino.Node) {
	if n == nil || !c.shouldTraverse(n, parent) {
		return
	}
	for _, child := range n.Children() {
		c.traverse(child, n)
	}
	c.visit(n)
}

func (c *CheckGlobalThis) shouldTraverse(n, parent *rhino.Node) bool {
	if n.Token() == rhino.FUNCTION && !c.shouldEnterFunction(n, parent) {
		return false
	}
	if isClassStaticBlock(n) {
		return false
	}
	// `this` inside a class field definition refers to the instance
	// (non-static) or the class itself (static) — never global.
	if n.Token() == rhino.CLASS_FIELD_DEF {
		return false
	}
	if parent != nil && parent.Token() == rhino.ASSIGN {
		return c.enterAssignChild(n, parent)
	}
	return true
}

// shouldEnterFunction decides whether to descend into a FUNCTION.
// Arrow functions always descend (they inherit `this`). Non-arrow
// functions are skipped when their `this` binding is otherwise
// established via JSDoc, via the declared function type, or via the
// surrounding declaration context.
func (c *CheckGlobalThis) shouldEnterFunction(fn, parent *rhino.Node) bool {
	if fn.GetBooleanProp(rhino.PropArrowFn) {
		return true
	}
	jsDoc := functionJSDoc(fn)
	if jsDoc != nil {
		if jsDoc.IsConstructor() || jsDoc.IsInterface() ||
			hasThisType(jsDoc) || isOverride(jsDoc) {
			return false
		}
		if functionTypeBindsThis(jsDoc) {
			return false
		}
	}
	// Only descend into functions that could legally carry a @this
	// annotation — e.g. not functions passed as call arguments.
	if parent == nil {
		return true
	}
	switch parent.Token() {
	case rhino.BLOCK, rhino.SCRIPT, rhino.NAME, rhino.ASSIGN,
		rhino.STRING_KEY:
		// annotatable
	default:
		return false
	}
	// Skip methods being lent to a prototype via @lends on the
	// enclosing object literal.
	if mayBeObjectLitKey(parent) {
		if grand := parent.Parent(); grand != nil {
			if lends := grand.JSDocInfo(); lends != nil && lendsIsPrototype(lends) {
				return false
			}
		}
	}
	return true
}

// enterAssignChild routes traversal for a node whose parent is
// ASSIGN. Records the LHS in assignLhsChild so `visit` knows to
// report bare `this` there, and prunes the RHS of prototype-slot
// assignments (where the RHS function establishes its own `this`).
func (c *CheckGlobalThis) enterAssignChild(n, parent *rhino.Node) bool {
	kids := parent.Children()
	if len(kids) == 0 {
		return true
	}
	lhs := kids[0]
	if n == lhs {
		if c.assignLhsChild == nil {
			c.assignLhsChild = lhs
		}
		return true
	}
	// RHS: skip `X.prototype.m = …` — the function sets its own
	// `this`.
	if !isNormalGet(lhs) {
		return true
	}
	if lhs.Token() == rhino.GETPROP && propName(lhs) == "prototype" {
		return false
	}
	if inner := lhs.FirstChild(); inner != nil &&
		inner.Token() == rhino.GETPROP && propName(inner) == "prototype" {
		return false
	}
	return true
}

func (c *CheckGlobalThis) visit(n *rhino.Node) {
	if n.Token() == rhino.THIS && c.shouldReportThis(n) {
		c.diags = append(c.diags, diag.Diagnostic{
			Path:    c.path,
			Line:    n.Lineno(),
			Col:     n.Charno(),
			Code:    diag.JSC_USED_GLOBAL_THIS,
			Level:   diag.Warning,
			Message: "dangerous use of the global this object",
		})
	}
	if n == c.assignLhsChild {
		c.assignLhsChild = nil
	}
}

func (c *CheckGlobalThis) shouldReportThis(n *rhino.Node) bool {
	if c.assignLhsChild != nil {
		return true
	}
	parent := n.Parent()
	return parent != nil && isNormalGet(parent)
}

// ---- helpers ----------------------------------------------------- //

// hasThisType approximates JSDocInfo.hasThisType(). jsgo doesn't parse
// `@this` into JSDocInfo yet; once it does, swap this for a direct
// field read.
func hasThisType(*rhino.JSDocInfo) bool { return false }

// isOverride approximates JSDocInfo.isOverride(). `@override` isn't
// parsed yet.
func isOverride(*rhino.JSDocInfo) bool { return false }

// functionTypeBindsThis reports whether `info.getType()` is a
// `function(this:T, …)` expression. jsgo stores type expressions as
// a flat string rather than a Node tree, so we recognise the intent
// by a textual prefix match — adequate for the oracles we track.
func functionTypeBindsThis(info *rhino.JSDocInfo) bool {
	if info.Type() == nil {
		return false
	}
	name := strings.TrimPrefix(info.Type().TypeName, "function(")
	if len(name) == len(info.Type().TypeName) {
		return false
	}
	trimmed := strings.TrimLeft(name, " \t")
	return strings.HasPrefix(trimmed, "this:") || strings.HasPrefix(trimmed, "this :")
}

// lendsIsPrototype reports whether a `@lends {Foo.prototype}` JSDoc
// targets a prototype. jsgo doesn't parse @lends yet.
func lendsIsPrototype(*rhino.JSDocInfo) bool { return false }

// mayBeObjectLitKey mirrors NodeUtil.mayBeObjectLitKey.
func mayBeObjectLitKey(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.STRING_KEY {
		return false
	}
	p := n.Parent()
	return p != nil && p.Token() == rhino.OBJECTLIT
}

// isClassStaticBlock reports whether `n` is a class's `static {}`
// initialization block. jsgo's Token enum doesn't distinguish static
// blocks from ordinary BLOCKs yet, so this always returns false —
// the guard shape is preserved for when that lands.
func isClassStaticBlock(*rhino.Node) bool { return false }

// isNormalGet mirrors NodeUtil.isNormalGet — property access that's
// not an optional chain. jsgo's parser doesn't tokenise `?.` as a
// distinct GETPROP variant, so the GETPROP / GETELEM pair covers the
// full space.
func isNormalGet(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.GETPROP, rhino.GETELEM:
		return true
	}
	return false
}

// propName returns the identifier on a GETPROP's `.prop` segment.
func propName(n *rhino.Node) string {
	kids := n.Children()
	if len(kids) < 2 {
		return ""
	}
	return kids[1].String_()
}
