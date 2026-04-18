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
	case rhino.CLASS:
		// Descend into class body: methods check via checkFunction
		// (they run in a this-bound local scope), fields get a
		// JSC_TYPE_MISMATCH when their initializer doesn't match the
		// @type declaration. Mirrors TypeCheck.java's visitClass.
		if kids := stmt.Children(); len(kids) >= 3 {
			className := kids[0].String_()
			for _, m := range kids[2].Children() {
				switch m.Token() {
				case rhino.FUNCTION:
					out = append(out, c.checkFunction(m)...)
				case rhino.CLASS_FIELD_DEF:
					out = append(out, c.checkClassField(m, scope, className)...)
				}
			}
		}
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
	case rhino.GETPROP:
		out = append(out, c.checkGetProp(n, scope)...)
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
	lhs, rhs := kids[0], kids[1]
	infer := NewTypeInference(scope)
	switch lhs.Token() {
	case rhino.NAME:
		declared := scope.GetDeclaredType(lhs.String_())
		if declared == nil {
			return nil
		}
		found := infer.InferType(rhs)
		// Conservative: dataflow narrowing on union-typed sources isn't
		// implemented; suppress the check to avoid false positives.
		if _, isUnion := found.(*jstype.UnionType); isUnion {
			return nil
		}
		if jstype.IsSubtype(found, declared) {
			return nil
		}
		msg := fmt.Sprintf("assignment\nfound   : %s\nrequired: %s",
			found.TypeName(), declared.TypeName())
		return []diag.Diagnostic{{
			Path:    c.path,
			Line:    lhs.Lineno(),
			Col:     lhs.Charno(),
			Code:    diag.JSC_TYPE_MISMATCH,
			Level:   diag.Error,
			Message: msg,
		}}
	case rhino.GETPROP:
		// Check `obj.prop = rhs` against obj's declared property type.
		propKids := lhs.Children()
		if len(propKids) < 2 {
			return nil
		}
		obj, ok := infer.InferType(propKids[0]).(*jstype.ObjectType)
		if !ok {
			return nil
		}
		propName := propKids[1].String_()
		declared := obj.GetPropertyType(propName)
		if declared == nil {
			return nil
		}
		found := infer.InferType(rhs)
		if _, isUnion := found.(*jstype.UnionType); isUnion {
			return nil
		}
		if jstype.IsSubtype(found, declared) {
			return nil
		}
		msg := fmt.Sprintf("assignment to property %s of %s\nfound   : %s\nrequired: %s",
			propName, obj.TypeName(), found.TypeName(), declared.TypeName())
		return []diag.Diagnostic{{
			Path:    c.path,
			Line:    lhs.Lineno(),
			Col:     lhs.Charno(),
			Code:    diag.JSC_TYPE_MISMATCH,
			Level:   diag.Error,
			Message: msg,
		}}
	}
	return nil
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
	infer := NewTypeInference(scope)
	fn, ok := infer.InferType(callee).(*jstype.FunctionType)
	if !ok {
		return nil
	}
	// Arity mismatch is out of Tier-1 scope (needs @param optional /
	// variadic parsing). Skip to avoid false positives.
	if len(args) != len(fn.Params) {
		return nil
	}
	var out []diag.Diagnostic
	calleeName := describeCallee(callee)
	for i, arg := range args {
		param := fn.Params[i]
		found := infer.InferType(arg)
		if _, isUnion := found.(*jstype.UnionType); isUnion {
			continue
		}
		if jstype.IsSubtype(found, param) {
			continue
		}
		msg := fmt.Sprintf("actual parameter %d of %s does not match formal parameter\nfound   : %s\nrequired: %s",
			i+1, calleeName, found.TypeName(), param.TypeName())
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
					// Java narrows union types at the return point via
					// dataflow analysis (TypeInference + FlowScope). We
					// don't, so when `found` is a UnionType we'd issue
					// false positives on code that's in fact safe.
					// Stay conservative: skip the check in that shape.
					if _, isUnion := found.(*jstype.UnionType); !isUnion &&
						!jstype.IsSubtype(found, declaredReturn) {
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

// checkClassField validates that a class field's initializer matches
// its @type annotation, reporting at the NAME position to match Java's
// "assignment to property X of C" shape.
func (c *TypeCheck) checkClassField(field *rhino.Node, scope *TypedScope, className string) []diag.Diagnostic {
	kids := field.Children()
	if len(kids) < 2 {
		return nil
	}
	name := kids[0]
	init := kids[1]
	info := field.JSDocInfo()
	if info == nil || !info.HasType() {
		return nil
	}
	declared := scope.ResolveExpr(info.Type().TypeName)
	infer := NewTypeInference(scope)
	found := infer.InferType(init)
	if _, isUnion := found.(*jstype.UnionType); isUnion {
		return nil
	}
	if jstype.IsSubtype(found, declared) {
		return nil
	}
	msg := fmt.Sprintf("assignment to property %s of %s\nfound   : %s\nrequired: %s",
		name.String_(), className, found.TypeName(), declared.TypeName())
	return []diag.Diagnostic{{
		Path:    c.path,
		Line:    name.Lineno(),
		Col:     name.Charno(),
		Code:    diag.JSC_TYPE_MISMATCH,
		Level:   diag.Error,
		Message: msg,
	}}
}

// describeCallee renders a NAME or GETPROP chain as a dotted string
// for the "actual parameter N of X" diagnostic, matching how Java
// formats the callee in TypeCheck.java's visitCall message.
func describeCallee(n *rhino.Node) string {
	switch n.Token() {
	case rhino.NAME:
		return n.String_()
	case rhino.GETPROP:
		kids := n.Children()
		if len(kids) >= 2 {
			return describeCallee(kids[0]) + "." + kids[1].String_()
		}
	case rhino.THIS:
		return "this"
	}
	return "<fn>"
}

// checkGetProp mirrors TypeCheck.java's visitGetProp: if the receiver
// type is a declared structural ObjectType (a @record or @interface)
// whose property set is closed and the accessed property isn't
// present anywhere up the @extends chain, emit
// JSC_INEXISTENT_PROPERTY.
//
// Gated strictly to structural types because Closure's default
// ObjectType is open (extra properties can be attached dynamically
// unless @struct is set). Firing on nominal types would produce many
// false positives on valid late-bound-property code like
// `f().a = 3; g(y.a);`.
func (c *TypeCheck) checkGetProp(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	propNode := kids[1]
	propName := propNode.String_()
	infer := NewTypeInference(scope)
	recvType := infer.InferType(kids[0])
	obj, ok := recvType.(*jstype.ObjectType)
	if !ok {
		return nil
	}
	if !obj.IsStructural {
		return nil
	}
	if obj.HasProperty(propName) {
		return nil
	}
	msg := fmt.Sprintf("Property %s never defined on %s", propName, obj.TypeName())
	return []diag.Diagnostic{{
		Path:    c.path,
		Line:    propNode.Lineno(),
		Col:     propNode.Charno(),
		Code:    diag.JSC_INEXISTENT_PROPERTY,
		Level:   diag.Error,
		Message: msg,
	}}
}

// resolveRegistry looks up a type expression through a fresh registry.
func (s *TypedScope) resolveRegistry(expr string) jstype.JSType {
	return jstype.NewRegistry().ResolveExpr(expr)
}
