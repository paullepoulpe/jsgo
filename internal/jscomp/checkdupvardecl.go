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
		if jstype.IsSubtype(t, first.t) && jstype.IsSubtype(first.t, t) {
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
	return c.diags
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
	if _, ok := c.scope.GetDeclaredType(recv.String_()).(*jstype.FunctionType); ok {
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
	if _, ok := c.scope.GetDeclaredType(recv.String_()).(*jstype.FunctionType); ok {
		return "", "", nil, nil
	}
	propNode := kids[1]
	return recv.String_(), propNode.String_(), propNode, c.scope.ResolveExpr(info.Type().TypeName)
}
