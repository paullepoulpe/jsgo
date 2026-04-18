// Port of com/google/javascript/jscomp/TypeCheck.java
// @ closure-compiler v20260415 (subset).
//
// TypeCheck implements the NodeTraversal.Callback contract and
// dispatches per-node in Visit — the same post-order visitor shape
// Java's class uses (extends AbstractPostOrderCallback, overrides
// visit(NodeTraversal, Node, Node)). Each visit-method covers one
// node kind and is a direct analogue of Java's visitCall /
// visitNew / visitAssign / visitGetProp / visitFunction /
// visitClassField dispatch.
//
// Diagnostics currently emitted:
//   - "initializing variable" JSC_TYPE_MISMATCH on var/let/const
//   - "assignment" / "assignment to property X of Y" JSC_TYPE_MISMATCH
//   - "actual parameter N of F does not match formal parameter"
//     JSC_TYPE_MISMATCH on call + `new` sites
//   - "inconsistent return type" JSC_TYPE_MISMATCH on return statements
//   - JSC_INEXISTENT_PROPERTY on structural receivers with closed
//     property sets

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
	path    string
	diags   []diag.Diagnostic
}

func NewTypeCheck(path string, creator *TypedScopeCreator, global *TypedScope) *TypeCheck {
	return &TypeCheck{
		creator: creator,
		global:  global,
		path:    path,
	}
}

// Check drives a NodeTraversal over the script and returns every
// diagnostic collected during Visit. Matches Java's
// TypeCheck#process(…) entry-point shape at a narrow scope.
func (c *TypeCheck) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	trav := NewNodeTraversal(c, c.creator, c.global)
	trav.Traverse(script)
	return c.diags
}

// ShouldTraverse is the pre-order hook. TypeCheck always descends —
// diagnostic emission happens in Visit.
func (c *TypeCheck) ShouldTraverse(_ *NodeTraversal, _, _ *rhino.Node) bool {
	return true
}

// Visit is the post-order dispatch. Matches Java's visit(t,n,parent)
// override on AbstractPostOrderCallback; each branch mirrors a
// TypeCheck.java visitXxx helper.
func (c *TypeCheck) Visit(t *NodeTraversal, n, parent *rhino.Node) {
	scope := t.CurrentScope()
	switch n.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		c.emit(c.visitVar(n, scope))
	case rhino.ASSIGN:
		c.emit(c.visitAssign(n, scope))
	case rhino.CALL:
		c.emit(c.visitCall(n, scope))
	case rhino.NEW:
		c.emit(c.visitCall(n, scope)) // share arg-list machinery
	case rhino.GETPROP:
		c.emit(c.visitGetProp(n, scope))
	case rhino.CLASS_FIELD_DEF:
		c.emit(c.visitClassField(n, parent, scope))
	case rhino.RETURN:
		c.emit(c.visitReturn(n, scope))
	}
}

func (c *TypeCheck) emit(d []diag.Diagnostic) {
	c.diags = append(c.diags, d...)
}

// ---- per-node visitors ------------------------------------------ //

// visitVar — port of TypeCheck.java's handling of VAR / LET / CONST:
// compare the initializer's inferred type against the @type
// declaration. Reports at the initializer's position to match Java's
// "initializing variable" variant.
func (c *TypeCheck) visitVar(v *rhino.Node, scope *TypedScope) []diag.Diagnostic {
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
		Path: c.path, Line: init.Lineno(), Col: init.Charno(),
		Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
	}}
}

// visitAssign handles bare `IDENT = …` and `obj.prop = …` assignment
// forms. Reports at the LHS position (matching Java's "assignment"
// and "assignment to property X of Y" variants).
func (c *TypeCheck) visitAssign(a *rhino.Node, scope *TypedScope) []diag.Diagnostic {
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
		if _, isUnion := found.(*jstype.UnionType); isUnion {
			return nil
		}
		if jstype.IsSubtype(found, declared) {
			return nil
		}
		msg := fmt.Sprintf("assignment\nfound   : %s\nrequired: %s",
			found.TypeName(), declared.TypeName())
		return []diag.Diagnostic{{
			Path: c.path, Line: lhs.Lineno(), Col: lhs.Charno(),
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
		}}
	case rhino.GETPROP:
		propKids := lhs.Children()
		if len(propKids) < 2 {
			return nil
		}
		propName := propKids[1].String_()
		recv := propKids[0]
		var declared jstype.JSType
		var ownerName string
		if recv.Token() == rhino.NAME {
			if fn, ok := scope.GetDeclaredType(recv.String_()).(*jstype.FunctionType); ok && fn.Instance != nil {
				if t := fn.Instance.GetStaticPropertyType(propName); t != nil {
					declared = t
					ownerName = recv.String_()
				}
			}
		}
		if declared == nil {
			obj, ok := infer.InferType(recv).(*jstype.ObjectType)
			if !ok {
				return nil
			}
			declared = obj.GetPropertyType(propName)
			if declared == nil {
				return nil
			}
			ownerName = obj.TypeName()
		}
		found := infer.InferType(rhs)
		if _, isUnion := found.(*jstype.UnionType); isUnion {
			return nil
		}
		if jstype.IsSubtype(found, declared) {
			return nil
		}
		msg := fmt.Sprintf("assignment to property %s of %s\nfound   : %s\nrequired: %s",
			propName, ownerName, found.TypeName(), declared.TypeName())
		return []diag.Diagnostic{{
			Path: c.path, Line: lhs.Lineno(), Col: lhs.Charno(),
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
		}}
	}
	return nil
}

// visitCall handles both CALL and NEW nodes. Reports per-argument
// JSC_TYPE_MISMATCH at the argument's position, matching
// TypeCheck.java's visitCall / visitNew.
func (c *TypeCheck) visitCall(node *rhino.Node, scope *TypedScope) []diag.Diagnostic {
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
			Path: c.path, Line: arg.Lineno(), Col: arg.Charno(),
			Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
		})
	}
	return out
}

// visitGetProp — port of TypeCheck.java's visitGetProp. Fires
// JSC_INEXISTENT_PROPERTY only on NAMED structural receivers;
// anonymous inline records are widened through unions in Java and
// firing on them false-positives.
func (c *TypeCheck) visitGetProp(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
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
	if obj.Name == "" {
		return nil
	}
	if obj.HasProperty(propName) {
		return nil
	}
	msg := fmt.Sprintf("Property %s never defined on %s", propName, obj.TypeName())
	return []diag.Diagnostic{{
		Path: c.path, Line: propNode.Lineno(), Col: propNode.Charno(),
		Code: diag.JSC_INEXISTENT_PROPERTY, Level: diag.Error, Message: msg,
	}}
}

// visitReturn validates `return <expr>;` against the enclosing
// function's @return annotation. The enclosing function is found by
// walking Node parent links — mirroring how Java's TypeCheck tracks
// the current function through NodeUtil.getEnclosingFunction.
func (c *TypeCheck) visitReturn(ret *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	enclosing := findEnclosingFunction(ret)
	if enclosing == nil {
		return nil
	}
	info := enclosing.JSDocInfo()
	if info == nil || info.ReturnType() == nil {
		return nil
	}
	declared := scope.ResolveExpr(info.ReturnType().TypeName)
	retExpr := ret.FirstChild()
	if retExpr == nil {
		return nil
	}
	infer := NewTypeInference(scope)
	found := infer.InferType(retExpr)
	if _, isUnion := found.(*jstype.UnionType); isUnion {
		return nil
	}
	if jstype.IsSubtype(found, declared) {
		return nil
	}
	msg := fmt.Sprintf("inconsistent return type\nfound   : %s\nrequired: %s",
		found.TypeName(), declared.TypeName())
	return []diag.Diagnostic{{
		Path: c.path, Line: retExpr.Lineno(), Col: retExpr.Charno(),
		Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
	}}
}

// visitClassField — port of the class-field branch inside Java's
// visitClass. Emits "assignment to property X of C" at the field
// name's position when the initializer type doesn't match the
// declared @type.
func (c *TypeCheck) visitClassField(field *rhino.Node, parent *rhino.Node, scope *TypedScope) []diag.Diagnostic {
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
	className := ""
	if parent != nil && parent.Parent() != nil {
		gp := parent.Parent()
		if gp.Token() == rhino.CLASS && gp.FirstChild() != nil {
			className = gp.FirstChild().String_()
		}
	}
	msg := fmt.Sprintf("assignment to property %s of %s\nfound   : %s\nrequired: %s",
		name.String_(), className, found.TypeName(), declared.TypeName())
	return []diag.Diagnostic{{
		Path: c.path, Line: name.Lineno(), Col: name.Charno(),
		Code: diag.JSC_TYPE_MISMATCH, Level: diag.Error, Message: msg,
	}}
}

// ---- helpers ----------------------------------------------------- //

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

// findEnclosingFunction walks Node parent links to the nearest
// FUNCTION ancestor. Mirrors NodeUtil.getEnclosingFunction in Java.
func findEnclosingFunction(n *rhino.Node) *rhino.Node {
	for cur := n.Parent(); cur != nil; cur = cur.Parent() {
		if cur.Token() == rhino.FUNCTION {
			return cur
		}
	}
	return nil
}
