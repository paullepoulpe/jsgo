// Port of com/google/javascript/jscomp/TypeCheck.java's
// checkDuplicateVarDeclaration branch @ closure-compiler v20260415
// (narrow; namespace-property slot only).
//
// Fires JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH when a namespace
// property is re-declared with a different type. Java tracks two
// shapes:
//
//   /** @type {T1} */ goog.foo;
//   /** @type {T2} */ goog.foo;
//
// and the assignment-form mirror. jsgo records every such decl
// on the scope's namespaceProps map; this pass walks the AST a
// second time, flagging each subsequent declaration against the
// first one's type.

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

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

// sameDeclaredType reports whether two JSTypes declared on the
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

// walkFunctionAssigns collects every `owner.prop = function(){…}`
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

// isFreshFunctionDecl reports whether the ASSIGN carries JSDoc
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

// walkTypedefStubs collects every `/** @typedef {T} */ owner.prop;`
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

// walkTemplateParmCollisions detects functions whose @template names
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

// walk recurses over the AST calling `emit` for every namespace-
// property type declaration it finds.
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
