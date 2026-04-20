// Port of com/google/javascript/jscomp/TypeCheck.java#visitCall
// (arity-mismatch branch) + visitNew (non-constructor branch)
// @ closure-compiler v20260415 (narrow).
//
// TypeCheck in Java fires JSC_WRONG_ARGUMENT_COUNT when a CALL's
// argument count falls outside the callee's [minArity, maxArity]
// window (maxArity is +∞ when the last param is rest) and
// JSC_NOT_A_CONSTRUCTOR when `new X()` targets a non-constructible
// callee. The two diagnostics live together in visitCall/visitNew
// alongside per-argument JSC_TYPE_MISMATCH reporting.
//
// This file runs the arity + constructibility check as a separate
// NodeTraversal pass so it can ship without touching TypeCheck. A
// pre-pass stamps ParamSpecs (Optional/Rest) onto every FunctionType
// the JSDoc parser annotated with modifiers — scope.go leaves those
// flags on the stored JSTypeExpression without propagating them
// further, so we patch them on here before walking call sites.

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// CheckArity implements JSC_WRONG_ARGUMENT_COUNT + JSC_NOT_A_CONSTRUCTOR.
// Mirrors the visitCall / visitNew dispatch in TypeCheck.java at a
// narrow scope.
type CheckArity struct {
	creator *TypedScopeCreator
	global  *TypedScope
	path    string
	diags   []diag.Diagnostic

	// userFns is the set of FunctionTypes whose source FUNCTION node
	// lives in the user script (as opposed to the externs bundle or
	// the builtins stub). Only call sites whose callee resolves into
	// this set get an arity check — externs declarations rely on
	// Closure's `opt_*` / `var_args` naming convention which jsgo
	// doesn't model, and false-positiving against the JS stdlib
	// would outweigh the true-positives unlocked elsewhere.
	userFns map[*jstype.FunctionType]bool
}

func NewCheckArity(path string, creator *TypedScopeCreator, global *TypedScope) *CheckArity {
	return &CheckArity{creator: creator, global: global, path: path, userFns: map[*jstype.FunctionType]bool{}}
}

// Check runs the traversal and returns the collected diagnostics.
// A pre-walk first populates ParamSpecs on any FunctionType whose
// source function carries `{T=}` / `{...T}` @param modifiers.
func (c *CheckArity) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.populateParamSpecs(script, c.global)
	trav := NewNodeTraversal(c, c.creator, c.global)
	trav.Traverse(script)
	return c.diags
}

func (c *CheckArity) ShouldTraverse(_ *NodeTraversal, _, _ *rhino.Node) bool { return true }

// Visit dispatches CALL and NEW nodes. Mirrors TypeCheck#visit.
func (c *CheckArity) Visit(t *NodeTraversal, n, parent *rhino.Node) {
	if n == nil {
		return
	}
	switch n.Token() {
	case rhino.CALL:
		c.visitCall(n, t.CurrentScope())
	case rhino.NEW:
		c.visitNew(n, t.CurrentScope())
	}
}

// visitCall fires JSC_WRONG_ARGUMENT_COUNT when the call site's arg
// count is outside the callee FunctionType's [MinArity, MaxArity]
// window. Mirrors the arity branch of TypeCheck.java#visitCall.
func (c *CheckArity) visitCall(node *rhino.Node, scope *TypedScope) {
	kids := node.Children()
	if len(kids) < 1 {
		return
	}
	callee := kids[0]
	args := kids[1:]
	fn := c.resolveCallee(callee, scope)
	// Only arity-check user-source functions — see userFns docs.
	if fn == nil || !c.userFns[fn] {
		return
	}
	minA := fn.MinArity()
	maxA := fn.MaxArity()
	n := len(args)
	if n >= minA && (maxA < 0 || n <= maxA) {
		return
	}
	// Java emits the diagnostic at the CALL node position, not the
	// callee — match that.
	calleeName := describeCallee(callee)
	maxText := fmt.Sprintf("%d", maxA)
	if maxA < 0 {
		maxText = "infinity"
	}
	msg := fmt.Sprintf("Function %s: called with %d argument(s). Function requires at least %d argument(s) and no more than %s argument(s).",
		calleeName, n, minA, maxText)
	c.diags = append(c.diags, diag.Diagnostic{
		Path: c.path, Line: node.Lineno(), Col: node.Charno(),
		Code: diag.JSC_WRONG_ARGUMENT_COUNT, Level: diag.Error, Message: msg,
	})
}

// visitNew fires JSC_NOT_A_CONSTRUCTOR when the receiver of `new` is
// a callable that isn't marked constructible. Interfaces / records
// carry IsConstructor=true but their Instance.IsStructural guards
// them from instantiation — we reject those too, matching Java.
func (c *CheckArity) visitNew(node *rhino.Node, scope *TypedScope) {
	kids := node.Children()
	if len(kids) < 1 {
		return
	}
	callee := kids[0]
	ctype := c.inferCalleeType(callee, scope)
	if ctype == nil {
		return
	}
	// Only complain about FunctionType receivers, and only when the
	// source FUNCTION lives in user code. The Function / Object
	// built-in nominal types are callable-plus-constructible; a
	// `@param {Function}` binding should be instantiable without
	// warning.
	fn, ok := ctype.(*jstype.FunctionType)
	if !ok || !c.userFns[fn] || !isNotConstructor(fn) {
		return
	}
	msg := fmt.Sprintf("cannot instantiate non-constructor, found type: %s", ctype.TypeName())
	// Java reports the diagnostic at the `new` keyword's column,
	// which sits 4 chars (`new `) before the callee's NAME/GETPROP.
	// Our parser records the NEW node at the callee's position, so
	// we back up by len("new ") to match the oracle's column.
	col := node.Charno() - 4
	if col < 0 {
		col = node.Charno()
	}
	c.diags = append(c.diags, diag.Diagnostic{
		Path: c.path, Line: node.Lineno(), Col: col,
		Code: diag.JSC_NOT_A_CONSTRUCTOR, Level: diag.Error, Message: msg,
	})
}

// resolveCallee returns the FunctionType backing the call's callee,
// or nil when the target isn't a concrete function (unknown / any /
// non-function type). Stays conservative — nil means "skip the
// arity check" to match Java's "unknown skips all checks" stance.
func (c *CheckArity) resolveCallee(callee *rhino.Node, scope *TypedScope) *jstype.FunctionType {
	t := c.inferCalleeType(callee, scope)
	if t == nil {
		return nil
	}
	fn, ok := t.(*jstype.FunctionType)
	if !ok {
		return nil
	}
	return fn
}

// inferCalleeType returns the callee's type using the declared scope
// lookup for NAME references and the standard TypeInference for
// everything else. Returns nil when the result is Unknown/any so the
// callers can skip the check.
func (c *CheckArity) inferCalleeType(callee *rhino.Node, scope *TypedScope) jstype.JSType {
	var t jstype.JSType
	if callee.Token() == rhino.NAME {
		t = scope.GetDeclaredType(callee.String_())
	} else {
		t = NewTypeInference(scope).InferType(callee)
	}
	if t == nil || t == jstype.Unknown || t == jstype.All {
		return nil
	}
	return t
}

// isNotConstructor reports whether `fn` is a callable value that
// cannot be used with `new`. A plain FunctionType without
// IsConstructor is rejected; a @interface / @record FunctionType
// (IsConstructor=true but Instance.IsStructural=true) is rejected
// too because interfaces aren't instantiable.
func isNotConstructor(fn *jstype.FunctionType) bool {
	return !fn.IsConstructor || (fn.Instance != nil && fn.Instance.IsStructural)
}

// populateParamSpecs walks the user script pre-order, marking every
// user-source FunctionType in userFns and stamping ParamSpecs on the
// ones whose JSDoc carries per-param `{T=}` / `{...T}` modifiers.
// scope.go itself already resolved the JSTypeExpression to a JSType
// and stored it in fn.Params — the modifiers are the only piece it
// dropped, so we re-apply them here without re-running type
// resolution. Externs-declared functions are skipped because this
// walk only visits the user script root.
func (c *CheckArity) populateParamSpecs(script *rhino.Node, scope *TypedScope) {
	var walk func(n *rhino.Node, scope *TypedScope)
	walk = func(n *rhino.Node, scope *TypedScope) {
		if n == nil {
			return
		}
		if n.Token() == rhino.FUNCTION {
			c.recordUserFunction(n, scope)
		}
		// Push into the FUNCTION's local scope when we have one so
		// nested functions resolve against their outer @param decls.
		childScope := scope
		if c.creator != nil {
			if local, ok := c.creator.Functions[n]; ok && local != nil {
				childScope = local
			}
		}
		for _, child := range n.Children() {
			walk(child, childScope)
		}
	}
	walk(script, scope)
}

// recordUserFunction locates the FunctionType for a user-source
// FUNCTION node, adds it to userFns, and stamps ParamSpecs when
// the JSDoc carries `{T=}` / `{...T}` modifiers.
func (c *CheckArity) recordUserFunction(fn *rhino.Node, scope *TypedScope) {
	ft := c.findFunctionType(fn, scope)
	if ft != nil {
		c.userFns[ft] = true
	}
	c.stampParamSpecs(fn, ft)
}

// stampParamSpecs attaches Optional/Rest flags from a FUNCTION
// node's JSDoc to the matching FunctionType. No-op when no
// modifiers are present, when the FunctionType lookup failed, or
// when the param counts mismatch.
func (c *CheckArity) stampParamSpecs(fn *rhino.Node, ft *jstype.FunctionType) {
	if ft == nil {
		return
	}
	info := fn.JSDocInfo()
	if info == nil || len(info.ParameterNames()) == 0 {
		return
	}
	kids := fn.Children()
	if len(kids) < 2 {
		return
	}
	paramNode := kids[1]
	names := make([]string, 0, len(paramNode.Children()))
	for _, p := range paramNode.Children() {
		names = append(names, p.String_())
	}
	if len(ft.Params) != len(names) {
		return
	}
	specs := make([]jstype.ParamSpec, len(names))
	any := false
	for i, pn := range names {
		specs[i].Type = ft.Params[i]
		expr := info.ParameterType(pn)
		if expr == nil {
			continue
		}
		if expr.Optional {
			specs[i].Optional = true
			any = true
		}
		if expr.Rest {
			specs[i].Rest = true
			any = true
		}
	}
	if !any {
		return
	}
	ft.ParamSpecs = specs
}

// findFunctionType locates the FunctionType a FUNCTION AST node
// ended up declared as. Checks the named-function slot first
// (`function F(){}` / `var F = function(){}`), then the prototype
// attachment (`Foo.prototype.M = function(){}`).
func (c *CheckArity) findFunctionType(fn *rhino.Node, scope *TypedScope) *jstype.FunctionType {
	kids := fn.Children()
	if len(kids) < 1 {
		return nil
	}
	// Case 1: `function F(params){}` — name is the first child.
	if name := kids[0].String_(); name != "" {
		if t := scope.GetDeclaredType(name); t != nil {
			if ft, ok := t.(*jstype.FunctionType); ok {
				return ft
			}
		}
	}
	// Case 2: function expression on the RHS of a declaration or
	// assignment. Walk up the parent chain.
	parent := fn.Parent()
	if parent == nil {
		return nil
	}
	switch parent.Token() {
	case rhino.NAME:
		// `var X = function(){}` — parent NAME is the VAR binding.
		if t := scope.GetDeclaredType(parent.String_()); t != nil {
			if ft, ok := t.(*jstype.FunctionType); ok {
				return ft
			}
		}
	case rhino.ASSIGN:
		// `X.prototype.M = function(){}` or `X.M = function(){}`.
		lhs := parent.FirstChild()
		if lhs == nil {
			return nil
		}
		if lhs.Token() == rhino.GETPROP {
			lhsKids := lhs.Children()
			if len(lhsKids) >= 2 {
				// Prototype case: `owner.prototype.method`.
				inner := lhsKids[0]
				propName := lhsKids[1].String_()
				if inner.Token() == rhino.GETPROP {
					innerKids := inner.Children()
					if len(innerKids) >= 2 &&
						innerKids[0].Token() == rhino.NAME &&
						innerKids[1].String_() == "prototype" {
						ownerName := innerKids[0].String_()
						if ot, ok := scope.GetDeclaredType(ownerName).(*jstype.FunctionType); ok && ot.Instance != nil {
							if pt := ot.Instance.GetPropertyType(propName); pt != nil {
								if ft, ok := pt.(*jstype.FunctionType); ok {
									return ft
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}
