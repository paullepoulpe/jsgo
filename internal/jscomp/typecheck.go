// Port of com/google/javascript/jscomp/TypeCheck.java @ closure-compiler v20260415
//
// TypeCheck walks the AST after TypedScopeCreator + TypeInference and
// emits JSErrors. It handles:
//   - "initializing variable" JSC_TYPE_MISMATCH on var/let/const
//   - "assignment" JSC_TYPE_MISMATCH on bare reassignment
//   - "actual parameter N of F does not match formal parameter"
//     JSC_TYPE_MISMATCH on call + `new` sites
//   - "inconsistent return type" JSC_TYPE_MISMATCH on return statements
//
// The walk mirrors TypeCheck.java's NodeTraversal: every statement is
// checked, then we descend into child expressions so call/new arg
// checks fire even inside initializers.

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

type TypeCheck struct {
	creator *TypedScopeCreator
	global  *TypedScope
	infer   *TypeInference
	path    string
}

func NewTypeCheck(path string, creator *TypedScopeCreator, global *TypedScope) *TypeCheck {
	return &TypeCheck{
		creator: creator,
		global:  global,
		infer:   NewTypeInference(global),
		path:    path,
	}
}

// Check walks the SCRIPT and returns accumulated diagnostics.
func (c *TypeCheck) Check(script *rhino.Node) []diag.Diagnostic {
	return c.checkBlock(script, c.global)
}

func (c *TypeCheck) checkBlock(block *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	var out []diag.Diagnostic
	for _, child := range block.Children() {
		out = append(out, c.checkStatement(child, scope)...)
	}
	return out
}

func (c *TypeCheck) checkStatement(stmt *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	var out []diag.Diagnostic
	switch stmt.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		out = append(out, c.checkVar(stmt, scope)...)
		// Init expressions may contain nested calls whose args should
		// still be validated.
		if name := stmt.FirstChild(); name != nil {
			if init := name.FirstChild(); init != nil {
				out = append(out, c.walkExpr(init, scope)...)
			}
		}
	case rhino.EXPR_RESULT:
		if a := stmt.FirstChild(); a != nil {
			switch a.Token() {
			case rhino.ASSIGN:
				out = append(out, c.checkAssign(a, scope)...)
				if kids := a.Children(); len(kids) >= 2 {
					out = append(out, c.walkExpr(kids[1], scope)...)
				}
			case rhino.CALL:
				out = append(out, c.checkCall(a, scope)...)
				out = append(out, c.walkExprChildren(a, scope)...)
			case rhino.NEW:
				out = append(out, c.checkNew(a, scope)...)
				out = append(out, c.walkExprChildren(a, scope)...)
			default:
				out = append(out, c.walkExpr(a, scope)...)
			}
		}
	case rhino.FUNCTION:
		out = append(out, c.checkFunction(stmt)...)
	case rhino.RETURN:
		// Inside function bodies, return is handled by checkFunction's
		// per-statement loop. At script level, ignore.
	case rhino.BLOCK:
		out = append(out, c.checkBlock(stmt, scope)...)
	}
	return out
}

// walkExpr recursively visits an expression subtree, checking any
// CALL / NEW nodes it finds along the way.
func (c *TypeCheck) walkExpr(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	if n == nil {
		return nil
	}
	var out []diag.Diagnostic
	switch n.Token() {
	case rhino.CALL:
		out = append(out, c.checkCall(n, scope)...)
	case rhino.NEW:
		out = append(out, c.checkNew(n, scope)...)
	case rhino.FUNCTION:
		// Arrow or function-expression: descend into its body via
		// checkFunction so return-type checks fire.
		out = append(out, c.checkFunction(n)...)
		return out
	}
	out = append(out, c.walkExprChildren(n, scope)...)
	return out
}

func (c *TypeCheck) walkExprChildren(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	var out []diag.Diagnostic
	for _, k := range n.Children() {
		out = append(out, c.walkExpr(k, scope)...)
	}
	return out
}

func (c *TypeCheck) checkVar(v *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	name := v.FirstChild()
	if name == nil || name.Token() != rhino.NAME {
		return nil
	}
	init := name.FirstChild()
	if init == nil {
		return nil
	}
	declared := scope.GetDeclaredType(name.String_())
	if declared == nil {
		return nil
	}
	infer := NewTypeInference(scope)
	found := infer.InferType(init)
	if jstype.IsSubtype(found, declared) {
		return nil
	}
	msg := fmt.Sprintf("initializing variable\nfound   : %s\nrequired: %s",
		found.TypeName(), declared.TypeName())
	return []diag.Diagnostic{{
		Path:    c.path,
		Line:    init.Lineno(),
		Col:     init.Charno(),
		Code:    diag.JSC_TYPE_MISMATCH,
		Level:   diag.Error,
		Message: msg,
	}}
}

func (c *TypeCheck) checkAssign(a *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := a.Children()
	if len(kids) < 2 {
		return nil
	}
	name, rhs := kids[0], kids[1]
	if name.Token() != rhino.NAME {
		return nil
	}
	declared := scope.GetDeclaredType(name.String_())
	if declared == nil {
		return nil
	}
	infer := NewTypeInference(scope)
	found := infer.InferType(rhs)
	if jstype.IsSubtype(found, declared) {
		return nil
	}
	msg := fmt.Sprintf("assignment\nfound   : %s\nrequired: %s",
		found.TypeName(), declared.TypeName())
	return []diag.Diagnostic{{
		Path:    c.path,
		Line:    name.Lineno(),
		Col:     name.Charno(),
		Code:    diag.JSC_TYPE_MISMATCH,
		Level:   diag.Error,
		Message: msg,
	}}
}

// checkCall and checkNew share arg-list machinery; NEW reports under
// the function's name-shape even though the call comes through the new
// keyword, matching TypeCheck.java's behaviour.
func (c *TypeCheck) checkCall(call *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	return c.checkArgList(call, scope)
}

func (c *TypeCheck) checkNew(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	return c.checkArgList(n, scope)
}

func (c *TypeCheck) checkArgList(node *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := node.Children()
	if len(kids) < 1 {
		return nil
	}
	callee := kids[0]
	args := kids[1:]
	if callee.Token() != rhino.NAME {
		return nil
	}
	fn, ok := scope.GetDeclaredType(callee.String_()).(*jstype.FunctionType)
	if !ok {
		return nil
	}
	// Arity mismatch is out of Tier-1 scope (needs @param optional /
	// variadic parsing). Skip to avoid false positives.
	if len(args) != len(fn.Params) {
		return nil
	}
	var out []diag.Diagnostic
	infer := NewTypeInference(scope)
	for i, arg := range args {
		param := fn.Params[i]
		found := infer.InferType(arg)
		if jstype.IsSubtype(found, param) {
			continue
		}
		msg := fmt.Sprintf("actual parameter %d of %s does not match formal parameter\nfound   : %s\nrequired: %s",
			i+1, callee.String_(), found.TypeName(), param.TypeName())
		out = append(out, diag.Diagnostic{
			Path:    c.path,
			Line:    arg.Lineno(),
			Col:     arg.Charno(),
			Code:    diag.JSC_TYPE_MISMATCH,
			Level:   diag.Error,
			Message: msg,
		})
	}
	return out
}

// checkFunction validates @return consistency and recurses into the
// body with the function's local scope.
func (c *TypeCheck) checkFunction(fn *rhino.Node) []diag.Diagnostic {
	kids := fn.Children()
	if len(kids) < 3 {
		return nil
	}
	body := kids[2]
	localScope := c.creator.Functions[fn]
	if localScope == nil {
		localScope = c.global
	}

	var declaredReturn jstype.JSType
	if info := fn.JSDocInfo(); info != nil && info.ReturnType() != nil {
		declaredReturn = c.global.resolveRegistry(info.ReturnType().TypeName)
	}

	var out []diag.Diagnostic
	for _, stmt := range body.Children() {
		if stmt.Token() == rhino.RETURN {
			if declaredReturn != nil {
				if retExpr := stmt.FirstChild(); retExpr != nil {
					infer := NewTypeInference(localScope)
					found := infer.InferType(retExpr)
					if !jstype.IsSubtype(found, declaredReturn) {
						msg := fmt.Sprintf("inconsistent return type\nfound   : %s\nrequired: %s",
							found.TypeName(), declaredReturn.TypeName())
						out = append(out, diag.Diagnostic{
							Path:    c.path,
							Line:    retExpr.Lineno(),
							Col:     retExpr.Charno(),
							Code:    diag.JSC_TYPE_MISMATCH,
							Level:   diag.Error,
							Message: msg,
						})
					}
				}
			}
			// Still recurse into the returned expression for call-arg
			// checks.
			if retExpr := stmt.FirstChild(); retExpr != nil {
				out = append(out, c.walkExpr(retExpr, localScope)...)
			}
			continue
		}
		out = append(out, c.checkStatement(stmt, localScope)...)
	}
	return out
}

// resolveRegistry looks up a type expression through a fresh registry.
func (s *TypedScope) resolveRegistry(expr string) jstype.JSType {
	return jstype.NewRegistry().ResolveExpr(expr)
}
