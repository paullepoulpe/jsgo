// Port of com/google/javascript/jscomp/CheckAccessControls.java @ closure-compiler v20260415 (subset)
//
// Closure's CheckAccessControls is a ~1400-line visitor that gates
// `@private` / `@protected` / `@public` / `@package` / `@deprecated`
// / `@final` / `@const` access rules across the AST, consulting
// CollectFileOverviewVisibility for per-file defaults. This port
// covers the two diagnostics the tier-2 / tier-3 fixtures exercise:
//
//   - JSC_BAD_PRIVATE_PROPERTY_ACCESS — read of a `@private`-tagged
//     property from outside its declaring class's methods.
//   - JSC_BAD_PROTECTED_PROPERTY_ACCESS — read of a `@protected`
//     property from outside the declaring class's inheritance chain.
//
// Package / public / deprecated / final / const-property diagnostics
// and fileoverview visibility are not yet hooked up — they arrive
// when tier-3+ fixtures need them. The port deliberately keeps the
// pass single-file (one CheckSource invocation = one source + an
// externs scope); the Java "declaring file" concept reduces to "the
// declaring class for the property" in our world.
//
// Phase 1 (collect) walks every property declaration in the AST,
// maps the property to the ObjectType it belongs on (via the
// already-built TypedScope), and stamps per-property visibility on
// that ObjectType. Mirrors what Java's InferJSDocInfo does before
// CheckAccessControls runs — stamping `JSDocInfo.getVisibility()` on
// the FunctionType / ObjectType property slot.
//
// Phase 2 (check) walks every GETPROP, computes the receiver's
// ObjectType via TypeInference, and consults the visibility stamped
// in Phase 1. The current class is tracked on a stack that pushes
// when we enter a FUNCTION / CLASS and pops on exit — mirrors Java's
// `currentClassStack` with `ThisType` on the local scope as the
// analog for "receiver class of the enclosing method".

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// CheckAccessControls implements the @private / @protected property
// access-control diagnostics. Mirrors the Java class of the same name.
type CheckAccessControls struct {
	path    string
	creator *TypedScopeCreator
	scope   *TypedScope
	diags   []diag.Diagnostic

	// declaringClass records the ObjectType on which each
	// (receiver ObjectType, property name) pair was originally
	// declared. Populated in Phase 1. Used by the protected-access
	// check to compare the current class against the declaring
	// class's inheritance chain. Mirrors Java's
	// `AccessControlUtils.getObjectType` which returns the type that
	// actually declares the property (walking super chains).
	declaringClass map[propKey]*jstype.ObjectType

	// enclosingClass is the stack of ObjectTypes bound to `this`
	// across nested FUNCTION / CLASS scopes. Push on enter, pop on
	// exit. Mirrors Java's `currentClassStack` in CheckAccessControls.
	enclosingClass []*jstype.ObjectType
}

type propKey struct {
	owner *jstype.ObjectType
	name  string
}

// NewCheckAccessControls constructs the pass with a handle to the
// already-built TypedScope; CheckAccessControls assumes TypedScope
// population has already populated constructor-instance types.
func NewCheckAccessControls(path string, creator *TypedScopeCreator, scope *TypedScope) *CheckAccessControls {
	return &CheckAccessControls{
		path:           path,
		creator:        creator,
		scope:          scope,
		declaringClass: map[propKey]*jstype.ObjectType{},
	}
}

// Check runs the two-phase walk and returns the accumulated
// diagnostics.
func (c *CheckAccessControls) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.collectVisibilities(script, c.scope)
	c.walk(script, c.scope, nil)
	return c.diags
}

// --------------------------------------------------------------
// Phase 1: collect per-property visibility from JSDoc.
// --------------------------------------------------------------

// collectVisibilities walks the AST and stamps visibility bits on
// ObjectType properties. Covers the declaration shapes Tier-2
// TypedScopeCreator actually registers:
//
//   - `class C { /** @private */ method() {…} }`
//   - `class C { /** @private */ field = …; }`
//   - `/** @private */ C.prototype.method = …`
//   - `/** @constructor */ function C(){ /** @private */ this.x = …; }`
//
// For @private-at-variable-declaration (e.g. `/** @private */ var x = …`)
// there is no per-property visibility to stamp — those are scoped to
// a Var slot and are reported via BAD_PRIVATE_GLOBAL_ACCESS in Java
// (not covered by this port).
func (c *CheckAccessControls) collectVisibilities(n *rhino.Node, scope *TypedScope) {
	if n == nil {
		return
	}
	switch n.Token() {
	case rhino.CLASS:
		c.collectClassBody(n, scope)
	case rhino.EXPR_RESULT:
		if inner := n.FirstChild(); inner != nil && inner.Token() == rhino.ASSIGN {
			c.collectPrototypeAssign(inner, scope)
			c.collectThisPropertyAssign(inner, scope)
		}
	case rhino.FUNCTION:
		// Constructor bodies contain `this.x = …` writes that may
		// carry @private / @protected; walk the body with the
		// constructor's TypedScope active so `this` resolves to the
		// instance.
		if local, ok := c.creator.Functions[n]; ok && local != nil {
			for _, child := range n.Children() {
				c.collectVisibilities(child, local)
			}
			return
		}
	}
	for _, child := range n.Children() {
		c.collectVisibilities(child, scope)
	}
}

// collectClassBody stamps visibility for each class member definition.
// Mirrors Java's InferJSDocInfo pass on the CLASS_MEMBERS arm of
// createPropertyReference.
func (c *CheckAccessControls) collectClassBody(cls *rhino.Node, scope *TypedScope) {
	kids := cls.Children()
	if len(kids) < 3 {
		return
	}
	nameNode := kids[0]
	body := kids[2]
	if nameNode.String_() == "" {
		return
	}
	fn, ok := scope.GetDeclaredType(nameNode.String_()).(*jstype.FunctionType)
	if !ok || fn.Instance == nil {
		return
	}
	inst := fn.Instance
	for _, m := range body.Children() {
		info := m.JSDocInfo()
		if info == nil || info.Visibility() == rhino.VisibilityInherited {
			continue
		}
		memberName := ""
		switch m.Token() {
		case rhino.FUNCTION:
			if first := m.FirstChild(); first != nil {
				memberName = first.String_()
			}
		case rhino.CLASS_FIELD_DEF:
			if first := m.FirstChild(); first != nil {
				memberName = first.String_()
			}
		}
		if memberName == "" || memberName == "constructor" {
			continue
		}
		v := mapVisibility(info.Visibility())
		inst.SetPropertyVisibility(memberName, v)
		c.declaringClass[propKey{inst, memberName}] = inst
	}
}

// collectPrototypeAssign handles the `Foo.prototype.name = value`
// declaration form (ES5 class syntax). Mirrors
// recordPrototypeAssign in TypedScopeCreator — we read the same
// JSDocInfo on the ASSIGN's parent or the ASSIGN itself.
func (c *CheckAccessControls) collectPrototypeAssign(a *rhino.Node, scope *TypedScope) {
	kids := a.Children()
	if len(kids) < 2 {
		return
	}
	lhs := kids[0]
	if lhs.Token() != rhino.GETPROP {
		return
	}
	lhsKids := lhs.Children()
	if len(lhsKids) < 2 {
		return
	}
	// Expect lhs to be of shape `(GETPROP (GETPROP Name "prototype") propName)`
	inner := lhsKids[0]
	if inner.Token() != rhino.GETPROP {
		return
	}
	innerKids := inner.Children()
	if len(innerKids) < 2 {
		return
	}
	if innerKids[0].Token() != rhino.NAME || innerKids[1].String_() != "prototype" {
		return
	}
	ownerName := innerKids[0].String_()
	propName := lhsKids[1].String_()
	info := a.JSDocInfo()
	if info == nil || info.Visibility() == rhino.VisibilityInherited {
		return
	}
	fn, ok := scope.GetDeclaredType(ownerName).(*jstype.FunctionType)
	if !ok || fn.Instance == nil {
		return
	}
	v := mapVisibility(info.Visibility())
	fn.Instance.SetPropertyVisibility(propName, v)
	c.declaringClass[propKey{fn.Instance, propName}] = fn.Instance
}

// collectThisPropertyAssign handles `/** @private */ this.x = …`
// inside a constructor body. The enclosing TypedScope's ThisType is
// the instance we stamp. Mirrors collectThisPropertiesIn in
// TypedScopeCreator.
func (c *CheckAccessControls) collectThisPropertyAssign(a *rhino.Node, scope *TypedScope) {
	kids := a.Children()
	if len(kids) < 2 {
		return
	}
	lhs := kids[0]
	if lhs.Token() != rhino.GETPROP {
		return
	}
	lhsKids := lhs.Children()
	if len(lhsKids) < 2 {
		return
	}
	if lhsKids[0].Token() != rhino.THIS {
		return
	}
	info := a.JSDocInfo()
	if info == nil || info.Visibility() == rhino.VisibilityInherited {
		return
	}
	inst := scope.ThisType()
	if inst == nil {
		return
	}
	propName := lhsKids[1].String_()
	v := mapVisibility(info.Visibility())
	inst.SetPropertyVisibility(propName, v)
	c.declaringClass[propKey{inst, propName}] = inst
}

// mapVisibility converts the rhino JSDoc enum to the jstype property-
// slot enum. The two enums mirror each other (INHERITED == 0, …) so
// this is a direct cast; the function exists so a future shift of
// either enum doesn't silently miscompile.
func mapVisibility(v rhino.Visibility) jstype.Visibility {
	switch v {
	case rhino.VisibilityPrivate:
		return jstype.VisibilityPrivate
	case rhino.VisibilityProtected:
		return jstype.VisibilityProtected
	case rhino.VisibilityPublic:
		return jstype.VisibilityPublic
	case rhino.VisibilityPackage:
		return jstype.VisibilityPackage
	}
	return jstype.VisibilityInherited
}

// --------------------------------------------------------------
// Phase 2: walk GETPROPs and emit diagnostics.
// --------------------------------------------------------------

// walk walks the AST with an enclosing-class stack, pushing when we
// enter a FUNCTION / CLASS (via the scope's ThisType binding) and
// popping on exit. Mirrors CheckAccessControls.shouldTraverse /
// exitAccessControlScope in Java.
func (c *CheckAccessControls) walk(n *rhino.Node, scope *TypedScope, parent *rhino.Node) {
	if n == nil {
		return
	}
	pushed := false
	local := scope
	if c.creator != nil {
		if s, ok := c.creator.Functions[n]; ok && s != nil {
			local = s
			if inst := s.ThisType(); inst != nil {
				c.enclosingClass = append(c.enclosingClass, inst)
				pushed = true
			}
		}
	}
	if n.Token() == rhino.GETPROP {
		c.checkPropertyVisibility(n, parent, local)
	}
	for _, child := range n.Children() {
		c.walk(child, local, n)
	}
	if pushed {
		c.enclosingClass = c.enclosingClass[:len(c.enclosingClass)-1]
	}
}

// checkPropertyVisibility emits the @private / @protected diagnostic
// when the GETPROP's receiver's property slot has a restrictive
// visibility and the current class-stack position doesn't satisfy
// the access rule.
//
// Skipped cases (mirrors Java's early-outs):
//
//   - The GETPROP is itself the LHS of the declaring assignment
//     (`this.x = …` is the definition site, not an access).
//   - The receiver isn't a nominal ObjectType (union / unknown / …).
//   - The property carries VisibilityInherited (no tag).
func (c *CheckAccessControls) checkPropertyVisibility(n, parent *rhino.Node, scope *TypedScope) {
	kids := n.Children()
	if len(kids) < 2 {
		return
	}
	recv := kids[0]
	propNode := kids[1]
	propName := propNode.String_()
	if propName == "" {
		return
	}
	// Skip the declaration site: `/** @private */ Foo.prototype.x = …`
	// is a write, not a read. Same for `this.x = …` in the ctor.
	if parent != nil && parent.Token() == rhino.ASSIGN {
		if parent.FirstChild() == n {
			return
		}
	}
	recvType := NewTypeInference(scope).InferType(recv)
	obj := dereferencedObject(recvType)
	if obj == nil {
		return
	}
	v := obj.GetPropertyVisibility(propName)
	if v == jstype.VisibilityInherited || v == jstype.VisibilityPublic {
		return
	}
	declClass := c.findDeclaringClass(obj, propName)
	if declClass == nil {
		return
	}
	switch v {
	case jstype.VisibilityPrivate:
		if !c.insideDeclaringClass(declClass) {
			c.emitPrivate(propNode, propName, declClass)
		}
	case jstype.VisibilityProtected:
		if !c.insideClassOrSubclass(declClass) {
			c.emitProtected(propNode, propName, declClass)
		}
	}
}

// findDeclaringClass walks the super chain of `obj` to find the
// class that actually declared `name`. Mirrors
// AccessControlUtils.getObjectType in Java (which walks the implicit
// prototype chain via getOwnSlot).
func (c *CheckAccessControls) findDeclaringClass(obj *jstype.ObjectType, name string) *jstype.ObjectType {
	for cur := obj; cur != nil; cur = cur.Super() {
		if dc, ok := c.declaringClass[propKey{cur, name}]; ok {
			return dc
		}
	}
	return nil
}

// insideDeclaringClass reports whether the enclosing-class stack
// contains `declClass`. Mirrors Java's same-file-of-declaration
// check, reduced to same-class for our single-file model: access
// to a `@private` property is OK iff it's inside the class that
// declared it.
func (c *CheckAccessControls) insideDeclaringClass(declClass *jstype.ObjectType) bool {
	for _, cls := range c.enclosingClass {
		if cls == declClass {
			return true
		}
	}
	return false
}

// insideClassOrSubclass reports whether the enclosing-class stack
// contains either the declaring class or one of its subclasses.
// Mirrors Java's `scopeType.isSubtypeOf(ownerType)` check in
// checkProtectedPropertyVisibility.
func (c *CheckAccessControls) insideClassOrSubclass(declClass *jstype.ObjectType) bool {
	for _, cls := range c.enclosingClass {
		if cls == nil {
			continue
		}
		if cls.IsSubtypeOf(declClass) {
			return true
		}
	}
	return false
}

// emitPrivate records a JSC_BAD_PRIVATE_PROPERTY_ACCESS diagnostic.
// Message shape matches Java's format string
// "Access to private property {0} of {1} not allowed here." so CLI
// output remains recognisable.
func (c *CheckAccessControls) emitPrivate(propNode *rhino.Node, propName string, declClass *jstype.ObjectType) {
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    propNode.Lineno(),
		Col:     propNode.Charno(),
		Code:    diag.JSC_BAD_PRIVATE_PROPERTY_ACCESS,
		Level:   diag.Warning,
		Message: fmt.Sprintf("Access to private property %s of %s not allowed here.", propName, declClass.TypeName()),
	})
}

// emitProtected records a JSC_BAD_PROTECTED_PROPERTY_ACCESS
// diagnostic. Message shape matches Java's format string.
func (c *CheckAccessControls) emitProtected(propNode *rhino.Node, propName string, declClass *jstype.ObjectType) {
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    propNode.Lineno(),
		Col:     propNode.Charno(),
		Code:    diag.JSC_BAD_PROTECTED_PROPERTY_ACCESS,
		Level:   diag.Warning,
		Message: fmt.Sprintf("Access to protected property %s of %s not allowed here.", propName, declClass.TypeName()),
	})
}
