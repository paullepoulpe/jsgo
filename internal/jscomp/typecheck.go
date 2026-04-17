// Port of com/google/javascript/jscomp/TypeCheck.java @ closure-compiler v20260415
//
// TypeCheck walks the AST after TypedScopeCreator + TypeInference and
// emits JSErrors. Tier-0 implemented the "initializing variable" and
// "assignment" variants of JSC_TYPE_MISMATCH. Tier-1 adds:
//   - "actual parameter N of F does not match formal parameter" at call
//     sites (JSC_TYPE_MISMATCH)
//   - "inconsistent return type" inside function bodies whose @return
//     disagrees with the returned expression (JSC_TYPE_MISMATCH)
// Messages mirror Java's phrasing so dogfooding the CLI feels right,
// but only (path, line, col, code) is part of the conformance check.

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
	var out []diag.Diagnostic
	out = append(out, c.checkBlock(script, c.global)...)
	return out
}

// checkBlock walks a SCRIPT or BLOCK with the given active scope.
func (c *TypeCheck) checkBlock(block *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	var out []diag.Diagnostic
	for _, child := range block.Children() {
		out = append(out, c.checkStatement(child, scope)...)
	}
	return out
}

func (c *TypeCheck) checkStatement(stmt *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	switch {
	case isDeclToken(stmt.Token()):
		return c.checkVar(stmt, scope)
	case stmt.Token() == rhino.EXPR_RESULT:
		if a := stmt.FirstChild(); a != nil {
			switch a.Token() {
			case rhino.ASSIGN:
				return c.checkAssign(a, scope)
			case rhino.CALL:
				return c.checkCall(a, scope)
			}
		}
	case stmt.Token() == rhino.FUNCTION:
		return c.checkFunction(stmt)
	case stmt.Token() == rhino.RETURN:
		// RETURN at script level is nonsense; ignore. Inside a function
		// body, the enclosing checkFunction drives return-type checks.
	}
	return nil
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

// checkAssign handles bare `IDENT = LITERAL;` statements. Reports at
// the LHS position to match Java's "assignment" variant.
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

// checkCall handles `CALL(callee, arg...)`. Reports per-argument
// JSC_TYPE_MISMATCH diagnostics at each argument's position, matching
// TypeCheck.java's visitCall.
func (c *TypeCheck) checkCall(call *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := call.Children()
	if len(kids) < 1 {
		return nil
	}
	callee := kids[0]
	args := kids[1:]
	if callee.Token() != rhino.NAME {
		return nil
	}
	declared := scope.GetDeclaredType(callee.String_())
	fn, ok := declared.(*jstype.FunctionType)
	if !ok {
		return nil
	}
	// Arity mismatch is a structural issue jsgo treats as out-of-scope
	// for Tier-1 (Java's diagnostic depends on @param optional/variadic
	// annotations we don't yet parse). Skip checks when arity differs
	// rather than risk false positives.
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

// checkFunction validates the body of a function declaration: @return
// consistency with every `return EXPR;`, and recursively checks nested
// statements.
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
		// Re-resolve through the same registry the scope creator used.
		declaredReturn = c.global.resolveRegistry(info.ReturnType().TypeName)
	}

	var out []diag.Diagnostic
	for _, stmt := range body.Children() {
		if stmt.Token() == rhino.RETURN && declaredReturn != nil {
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
			continue
		}
		out = append(out, c.checkStatement(stmt, localScope)...)
	}
	return out
}

// resolveRegistry looks up a type name through the same registry the
// scope creator used. Wired on the global scope for simplicity — all
// primitive names resolve identically everywhere.
func (s *TypedScope) resolveRegistry(name string) jstype.JSType {
	return jstype.NewRegistry().ResolveName(name)
}
