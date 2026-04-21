// Port of com/google/javascript/jscomp/CheckGlobalThis.java @ closure-compiler v20260415
//
// Checks for uses of `this` that are likely to reference the global
// `this` object unintentionally. A `this` reference is considered
// unsafe only when it appears in a context that observes the binding
// (on the LHS of an assignment or as the receiver of a property
// access) and isn't inside a function whose `this` binding is
// otherwise established:
//
//   - a class static initialization block,
//   - a prototype method / object-method literal,
//   - a function annotated with `@constructor` / `@interface` /
//     `@this` / `@override`,
//   - a `function(this:T, …)` JSDoc-typed function expression,
//   - a context where attaching a `@this` annotation isn't natural
//     (e.g. a function passed as a call argument).
//
// The traversal tracks `assignLhsChild`, the first child of the
// nearest enclosing ASSIGN whose descendants we're walking — nested
// assigns don't override it (`(a = this).b = c` still reports the
// outer `this`).
//
// Note that this pass deliberately does _not_ track `this` through
// aliases (`var a = this; a.x = 1;` is silently accepted), matching
// the Java implementation.

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckGlobalThis mirrors Java's final class of the same name.
// Implements the NodeTraversal.Callback contract via ShouldTraverse /
// Visit, and walks the script directly (no external NodeTraversal
// driver needed — a simple recursive walker is sufficient because
// CheckGlobalThis doesn't consult the scope stack).
type CheckGlobalThis struct {
	path  string
	diags []diag.Diagnostic

	// assignLhsChild, when non-nil, is the first child of the ASSIGN
	// whose subtree we're currently inside. Mirrors Java's field of
	// the same name — see the class-level comment for the nested-
	// assignment invariant.
	assignLhsChild *rhino.Node
}

func NewCheckGlobalThis(path string) *CheckGlobalThis {
	return &CheckGlobalThis{path: path}
}

// Check walks the script and emits JSC_USED_GLOBAL_THIS at every
// `this` reference that resolves to the global object.
func (c *CheckGlobalThis) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.assignLhsChild = nil
	c.traverse(script, nil)
	return c.diags
}

// traverse is a post-order walker that mirrors NodeTraversal's
// ShouldTraverse / Visit dispatch. Runs ShouldTraverse first (may
// mutate assignLhsChild), descends on true, then runs Visit.
func (c *CheckGlobalThis) traverse(n, parent *rhino.Node) {
	if n == nil {
		return
	}
	if !c.shouldTraverse(n, parent) {
		return
	}
	for _, child := range n.Children() {
		c.traverse(child, n)
	}
	c.visit(n, parent)
}

// shouldTraverse mirrors CheckGlobalThis.shouldTraverse in Java.
// Returning false prunes the subtree; the key cases are functions
// whose `this` binding is established elsewhere, class static
// blocks / member fields, and the RHS of a prototype-slot assign.
func (c *CheckGlobalThis) shouldTraverse(n, parent *rhino.Node) bool {
	if n.Token() == rhino.FUNCTION {
		// Arrow functions inherit `this` from the enclosing scope —
		// always descend.
		if n.GetBooleanProp(rhino.PropArrowFn) {
			return true
		}
		jsDoc := bestJSDocInfo(n)
		if jsDoc != nil &&
			(jsDoc.IsConstructor() ||
				jsDoc.IsInterface() ||
				hasThisType(jsDoc) ||
				isOverride(jsDoc)) {
			return false
		}
		// `function(this:T, …)` in the declared type — `this` is
		// bound, so skip the body.
		if jsDoc != nil && functionTypeBindsThis(jsDoc) {
			return false
		}
		// Only traverse functions that could legally carry a @this
		// annotation. Functions passed as call arguments or wedged
		// inside expressions aren't normally annotatable, so Java's
		// CheckGlobalThis ignores their bodies.
		if parent != nil {
			switch parent.Token() {
			case rhino.BLOCK, rhino.SCRIPT, rhino.NAME, rhino.ASSIGN,
				rhino.STRING_KEY:
				// annotatable context — fall through
			default:
				return false
			}
			// Skip functions being lent to a prototype via
			// `@lends {Foo.prototype}` on the enclosing object
			// literal.
			if mayBeObjectLitKey(parent) {
				if grand := parent.Parent(); grand != nil {
					if lends := grand.JSDocInfo(); lends != nil && lendsIsPrototype(lends) {
						return false
					}
				}
			}
		}
	}

	// Class static initialization blocks never see the global `this`.
	if isClassStaticBlock(n) {
		return false
	}

	// `this` inside a class field definition refers to the instance
	// (non-static field) or the class itself (static field).
	if isMemberOrComputedField(n) {
		return false
	}

	if parent != nil && parent.Token() == rhino.ASSIGN {
		kids := parent.Children()
		if len(kids) >= 1 {
			lhs := kids[0]
			if n == lhs {
				// Always traverse the LHS. Nested assigns must not
				// overwrite assignLhsChild so we still report the
				// outer `this` in `(a = this).b = c`.
				if c.assignLhsChild == nil {
					c.assignLhsChild = lhs
				}
			} else {
				// RHS: skip assignments that target a prototype slot
				// — `X.prototype.m = function(){…}` establishes the
				// function's `this` binding.
				if isNormalGet(lhs) {
					if lhs.Token() == rhino.GETPROP && propName(lhs) == "prototype" {
						return false
					}
					if llhs := lhs.FirstChild(); llhs != nil &&
						llhs.Token() == rhino.GETPROP && propName(llhs) == "prototype" {
						return false
					}
				}
			}
		}
	}

	return true
}

// visit is the post-order callback. Emits JSC_USED_GLOBAL_THIS when a
// `this` keyword resolves to the global object, and clears
// assignLhsChild on the way back up through the ASSIGN LHS.
func (c *CheckGlobalThis) visit(n, _ *rhino.Node) {
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

// shouldReportThis mirrors Java's helper of the same name. Fires
// either inside the LHS of an assignment or when the immediate
// parent is a property access.
func (c *CheckGlobalThis) shouldReportThis(n *rhino.Node) bool {
	if c.assignLhsChild != nil {
		return true
	}
	parent := n.Parent()
	return parent != nil && isNormalGet(parent)
}

// ---- helpers ----------------------------------------------------- //

// bestJSDocInfo mirrors NodeUtil.getBestJSDocInfo: prefers the JSDoc
// directly attached to `fn`, otherwise walks up through the typical
// declaration wrappers (ASSIGN / VAR / LET / CONST / EXPR_RESULT /
// NAME) looking for the nearest comment. Stops at the first
// non-wrapper ancestor so we don't attach an unrelated surrounding
// comment.
func bestJSDocInfo(fn *rhino.Node) *rhino.JSDocInfo {
	if fn == nil {
		return nil
	}
	if info := fn.JSDocInfo(); info != nil {
		return info
	}
	for cur := fn.Parent(); cur != nil; cur = cur.Parent() {
		if info := cur.JSDocInfo(); info != nil {
			return info
		}
		switch cur.Token() {
		case rhino.ASSIGN, rhino.EXPR_RESULT, rhino.VAR, rhino.LET,
			rhino.CONST, rhino.NAME:
			continue
		}
		break
	}
	return nil
}

// hasThisType approximates JSDocInfo.hasThisType() — jsgo's JSDoc
// surface doesn't record the `@this` tag directly, so this always
// returns false until the annotation is plumbed through
// JsDocInfoParser. Kept as a named helper so the call-site reads
// like the Java source.
func hasThisType(*rhino.JSDocInfo) bool { return false }

// isOverride approximates JSDocInfo.isOverride(). `@override` isn't
// parsed yet; stub to false so the call-site structure matches the
// Java source.
func isOverride(*rhino.JSDocInfo) bool { return false }

// functionTypeBindsThis checks whether a JSDoc `@type
// {function(this:T, …)}` declares a concrete `this`. jsgo stores
// type expressions as a flat TypeName string rather than a Node
// tree, so this recovers the intent by scanning for the `this:`
// marker — good enough for the fixtures we have.
func functionTypeBindsThis(info *rhino.JSDocInfo) bool {
	if info == nil || info.Type() == nil {
		return false
	}
	name := info.Type().TypeName
	if !strings.HasPrefix(name, "function(") {
		return false
	}
	// Strip `function(` prefix, look at the first argument spelling.
	inner := strings.TrimPrefix(name, "function(")
	trimmed := strings.TrimLeft(inner, " \t")
	return strings.HasPrefix(trimmed, "this:") || strings.HasPrefix(trimmed, "this :")
}

// lendsIsPrototype reports whether a `@lends {Foo.prototype}` JSDoc
// targets a prototype. Matches Java's
// `jsDoc.hasLendsName() &&
//  jsDoc.getLendsName().getRoot().getString().endsWith(".prototype")`.
// jsgo doesn't parse `@lends` yet; the helper is stubbed to false so
// the call-site shape mirrors the Java source.
func lendsIsPrototype(*rhino.JSDocInfo) bool { return false }

// mayBeObjectLitKey reports whether `n` is an object-literal key
// (shorthand or STRING_KEY). Matches NodeUtil.mayBeObjectLitKey.
func mayBeObjectLitKey(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	if n.Token() != rhino.STRING_KEY {
		return false
	}
	p := n.Parent()
	return p != nil && p.Token() == rhino.OBJECTLIT
}

// isClassStaticBlock reports whether `n` is a class's `static {}`
// initialization block. jsgo's Token enum doesn't yet distinguish
// static blocks from ordinary BLOCKs, so this is always false —
// matches the Java guard's shape without filtering anything.
func isClassStaticBlock(*rhino.Node) bool { return false }

// isMemberOrComputedField reports whether `n` is a class member
// field (`x = 1`) or computed field (`[key] = 1`) attached to the
// class body. jsgo only tokenises CLASS_FIELD_DEF today, which
// covers the non-computed case; the computed variant is recognised
// when Token.COMPUTED_FIELD_DEF lands.
func isMemberOrComputedField(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.CLASS_FIELD_DEF {
		return false
	}
	// In Java this also checks that the parent is CLASS_MEMBERS;
	// jsgo hangs CLASS_FIELD_DEF directly off CLASS, so the parent
	// test degrades to CLASS.
	p := n.Parent()
	return p != nil && p.Token() == rhino.CLASS
}

// isNormalGet reports whether `n` is a property access (`x.y` /
// `x[y]`). Matches NodeUtil.isNormalGet, which excludes optional
// chaining — jsgo doesn't parse `?.` yet so the GETPROP / GETELEM
// pair covers the whole space.
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

// propName returns the property identifier on a GETPROP. jsgo stores
// `x.y`'s `y` as the second child's String_ (mirrors Java's
// `Node.getString()` on a GETPROP).
func propName(n *rhino.Node) string {
	kids := n.Children()
	if len(kids) < 2 {
		return ""
	}
	return kids[1].String_()
}
