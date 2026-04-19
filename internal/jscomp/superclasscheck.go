// Port of com/google/javascript/jscomp/TypeCheck.java's
// checkPropertyVisibility / validateMethodOverride branches
// @ closure-compiler v20260415 (narrow; no access-modifier check,
// no `@override`-required diagnostic).
//
// Fires JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH when an override of
// a superclass property has an incompatible type. Two override
// shapes are recognised:
//
//   X.prototype.method = function(...) {...};   // if X extends Y
//   /** @override */
//   Y.prototype.method = …
//
//   class Derived extends Base { method() { ... } }
//
// The check is scope-driven: it walks class / prototype-assign AST
// nodes, resolves the override's type via TypeInference, and
// compares against the superclass's slot via ObjectType's super
// chain.

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

type SuperclassCheck struct {
	path       string
	creator    *TypedScopeCreator
	scope      *TypedScope
	scriptRoot *rhino.Node
	diags      []diag.Diagnostic
}

func NewSuperclassCheck(path string, creator *TypedScopeCreator, scope *TypedScope) *SuperclassCheck {
	return &SuperclassCheck{path: path, creator: creator, scope: scope}
}

func (s *SuperclassCheck) Check(script *rhino.Node) []diag.Diagnostic {
	s.diags = nil
	s.scriptRoot = script
	s.walk(script)
	return s.diags
}

func (s *SuperclassCheck) walk(n *rhino.Node) {
	if n == nil {
		return
	}
	switch n.Token() {
	case rhino.CLASS:
		s.checkClass(n)
	case rhino.FUNCTION:
		s.checkFunctionImplements(n)
	case rhino.VAR, rhino.LET, rhino.CONST:
		s.checkDeclaredConstructorImplements(n)
	case rhino.EXPR_RESULT:
		if a := n.FirstChild(); a != nil && a.Token() == rhino.ASSIGN {
			s.checkPrototypeAssign(n, a)
		}
	}
	for _, c := range n.Children() {
		s.walk(c)
	}
}

// checkFunctionImplements walks an AST FUNCTION declaration and,
// when it carries @constructor + @implements {Foo, Bar, …}, emits
// JSC_INTERFACE_METHOD_NOT_IMPLEMENTED for each interface property
// the constructor's instance doesn't implement (directly or
// through its @extends chain). Mirrors TypeCheck's
// validateImplements loop.
func (s *SuperclassCheck) checkFunctionImplements(fn *rhino.Node) {
	info := fn.JSDocInfo()
	if !info.IsConstructor() {
		return
	}
	kids := fn.Children()
	if len(kids) < 1 {
		return
	}
	nameNode := kids[0]
	if nameNode.String_() == "" {
		return
	}
	s.validateImplements(fn, nameNode.String_(), info)
}

// checkDeclaredConstructorImplements handles the
// `/** @constructor @implements {Foo} */ var Bar = function(){}`
// form — JSDoc is on the VAR, not on the inner FUNCTION.
func (s *SuperclassCheck) checkDeclaredConstructorImplements(v *rhino.Node) {
	info := v.JSDocInfo()
	if !info.IsConstructor() {
		return
	}
	nameNode := v.FirstChild()
	if nameNode == nil || nameNode.Token() != rhino.NAME || nameNode.String_() == "" {
		return
	}
	s.validateImplements(v, nameNode.String_(), info)
}

// validateImplements verifies `ctorName` implements every property
// each of its @implements interfaces declares. Emits one
// JSC_INTERFACE_METHOD_NOT_IMPLEMENTED per missing property. The
// diagnostic position matches Java's — at the constructor
// declaration node's starting column.
func (s *SuperclassCheck) validateImplements(reportAt *rhino.Node, ctorName string, info *rhino.JSDocInfo) {
	impls := info.Implementations()
	if len(impls) == 0 {
		return
	}
	ctor, ok := s.scope.GetRawDeclaredType(ctorName).(*jstype.FunctionType)
	if !ok || ctor.Instance == nil {
		return
	}
	inst := ctor.Instance
	for _, expr := range impls {
		iface := s.scope.LookupNominalType(expr.TypeName)
		if iface == nil {
			continue
		}
		// Walk the interface's full property set (its own members
		// plus any inherited via @extends {Other}).
		for _, propName := range allInterfaceProperties(iface) {
			if inst.GetPropertyType(propName) != nil {
				continue
			}
			s.diags = append(s.diags, diag.Diagnostic{
				Path: s.path, Line: reportAt.Lineno(), Col: reportAt.Charno(),
				Code: diag.JSC_INTERFACE_METHOD_NOT_IMPLEMENTED, Level: diag.Error,
				Message: fmt.Sprintf(
					"property %s on interface %s is not implemented by type %s",
					propName, iface.TypeName(), ctorName),
			})
		}
	}
}

// allInterfaceProperties returns every property name declared on
// `iface` or reachable through its super chain, deduplicated.
func allInterfaceProperties(iface *jstype.ObjectType) []string {
	seen := map[string]bool{}
	var out []string
	for cur := iface; cur != nil; cur = cur.Super() {
		for _, name := range cur.MemberNames() {
			if seen[name] {
				continue
			}
			seen[name] = true
			out = append(out, name)
		}
	}
	return out
}

// checkPrototypeAssign handles `X.prototype.name = rhs;`. When X has
// a super (via @extends) and super declares `name`, verify the rhs
// type matches — else fire JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH.
// Same shape runs against each @implemented interface's slots,
// firing JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH instead.
func (s *SuperclassCheck) checkPrototypeAssign(stmt *rhino.Node, assign *rhino.Node) {
	kids := assign.Children()
	if len(kids) < 2 {
		return
	}
	lhs, rhs := kids[0], kids[1]
	if lhs.Token() != rhino.GETPROP {
		return
	}
	lhsKids := lhs.Children()
	if len(lhsKids) < 2 {
		return
	}
	inner := lhsKids[0]
	propName := lhsKids[1].String_()
	propNode := lhsKids[1]
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
	fn, ok := s.scope.GetRawDeclaredType(ownerName).(*jstype.FunctionType)
	if !ok || fn.Instance == nil {
		return
	}
	overrideType := s.inferAssignValue(stmt, assign, rhs, fn.Instance)
	if overrideType == nil {
		return
	}
	// Superclass slot (via @extends chain).
	if super := fn.Instance.Super(); super != nil {
		if slot := super.GetPropertyType(propName); slot != nil {
			if !jstype.IsSubtype(overrideType, slot) {
				s.diags = append(s.diags, diag.Diagnostic{
					Path:  s.path,
					Line:  stmt.Lineno(),
					Col:   stmt.Charno(),
					Code:  diag.JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH,
					Level: diag.Error,
					Message: fmt.Sprintf("mismatch of the %s property type and the type of the property it overrides from superclass %s\noriginal: %s\noverride: %s",
						propName, super.TypeName(), slot.TypeName(), overrideType.TypeName()),
				})
				return
			}
		}
	}
	// @implements interfaces.
	for _, iface := range s.implementedInstances(ownerName) {
		slot := iface.GetPropertyType(propName)
		if slot == nil {
			continue
		}
		if jstype.IsSubtype(overrideType, slot) {
			continue
		}
		s.diags = append(s.diags, diag.Diagnostic{
			Path:  s.path,
			Line:  propNode.Lineno(),
			Col:   propNode.Charno(),
			Code:  diag.JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH,
			Level: diag.Error,
			Message: fmt.Sprintf("mismatch of the %s property on type %s and the type of the property it overrides from interface %s\noriginal: %s\noverride: %s",
				propName, ownerName, iface.TypeName(), slot.TypeName(), overrideType.TypeName()),
		})
		return
	}
}

// implementedInstances returns the interface instance types declared
// via `@implements {Foo}` on the constructor named by ownerName.
func (s *SuperclassCheck) implementedInstances(ownerName string) []*jstype.ObjectType {
	// The scope doesn't currently expose the JSDocInfo that produced
	// the FunctionType, so we walk the AST once to collect
	// @implements declarations. Caching would be worthwhile if this
	// became hot.
	info := s.findConstructorJSDoc(ownerName)
	if info == nil {
		return nil
	}
	var out []*jstype.ObjectType
	for _, e := range info.Implementations() {
		if inst := s.scope.LookupNominalType(e.TypeName); inst != nil {
			out = append(out, inst)
		}
	}
	return out
}

// findConstructorJSDoc scans the AST for a `function Foo(…)` or
// `var Foo = function(…)` constructor named `name` and returns its
// JSDocInfo. Returns nil when no matching declaration exists.
func (s *SuperclassCheck) findConstructorJSDoc(name string) *rhino.JSDocInfo {
	return s.scanConstructorJSDoc(s.scriptRoot, name)
}

func (s *SuperclassCheck) scanConstructorJSDoc(n *rhino.Node, name string) *rhino.JSDocInfo {
	if n == nil {
		return nil
	}
	for _, c := range n.Children() {
		if c.Token() == rhino.FUNCTION {
			if first := c.FirstChild(); first != nil && first.String_() == name {
				return c.JSDocInfo()
			}
		}
		if c.Token() == rhino.VAR || c.Token() == rhino.LET || c.Token() == rhino.CONST {
			nm := c.FirstChild()
			if nm != nil && nm.String_() == name {
				if init := nm.FirstChild(); init != nil && init.Token() == rhino.FUNCTION {
					if info := c.JSDocInfo(); info != nil {
						return info
					}
					return init.JSDocInfo()
				}
			}
		}
		if found := s.scanConstructorJSDoc(c, name); found != nil {
			return found
		}
	}
	return nil
}

// inferAssignValue resolves the type attached to `X.prototype.m = …`.
// Prefers an explicit @type / @param-based FunctionType from the
// enclosing ASSIGN's JSDoc; falls back to inferring the RHS via
// TypeInference.
func (s *SuperclassCheck) inferAssignValue(stmt, assign, rhs *rhino.Node, inst *jstype.ObjectType) jstype.JSType {
	info := assign.JSDocInfo()
	if info == nil {
		info = stmt.JSDocInfo()
	}
	if info != nil && info.HasType() {
		return s.scope.ResolveExpr(info.Type().TypeName)
	}
	if rhs.Token() == rhino.FUNCTION {
		return functionTypeFromJSDoc(info, rhs, s.scope)
	}
	ti := NewTypeInference(s.scope)
	return ti.InferType(rhs)
}

// checkClass verifies each method / field of a `class C extends B`
// declaration against B's slot.
func (s *SuperclassCheck) checkClass(cls *rhino.Node) {
	kids := cls.Children()
	if len(kids) < 3 {
		return
	}
	nameNode := kids[0]
	body := kids[2]
	preStub, ok := s.scope.GetRawDeclaredType(nameNode.String_()).(*jstype.FunctionType)
	if !ok || preStub.Instance == nil || preStub.Instance.Super() == nil {
		return
	}
	super := preStub.Instance.Super()
	for _, m := range body.Children() {
		switch m.Token() {
		case rhino.FUNCTION:
			if len(m.Children()) < 1 {
				continue
			}
			methodName := m.FirstChild().String_()
			if methodName == "" || methodName == "constructor" {
				continue
			}
			superSlot := super.GetPropertyType(methodName)
			if superSlot == nil {
				continue
			}
			info := m.JSDocInfo()
			ownSlot := functionTypeFromJSDoc(info, m, s.scope)
			if ownSlot == nil {
				continue
			}
			if jstype.IsSubtype(ownSlot, superSlot) {
				continue
			}
			s.diags = append(s.diags, diag.Diagnostic{
				Path:  s.path,
				Line:  m.Lineno(),
				Col:   m.Charno(),
				Code:  diag.JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH,
				Level: diag.Error,
				Message: fmt.Sprintf("mismatch of the %s property type and the type of the property it overrides from superclass %s\noriginal: %s\noverride: %s",
					methodName, super.TypeName(), superSlot.TypeName(), ownSlot.TypeName()),
			})
		case rhino.CLASS_FIELD_DEF:
			if len(m.Children()) < 1 {
				continue
			}
			fieldName := m.FirstChild().String_()
			if fieldName == "" {
				continue
			}
			superSlot := super.GetPropertyType(fieldName)
			if superSlot == nil {
				continue
			}
			info := m.JSDocInfo()
			if info == nil || !info.HasType() {
				continue
			}
			ownType := s.scope.ResolveExpr(info.Type().TypeName)
			if jstype.IsSubtype(ownType, superSlot) {
				continue
			}
			s.diags = append(s.diags, diag.Diagnostic{
				Path:  s.path,
				Line:  m.Lineno(),
				Col:   m.Charno(),
				Code:  diag.JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH,
				Level: diag.Error,
				Message: fmt.Sprintf("mismatch of the %s property type and the type of the property it overrides from superclass %s\noriginal: %s\noverride: %s",
					fieldName, super.TypeName(), superSlot.TypeName(), ownType.TypeName()),
			})
		}
	}
}

// functionTypeFromJSDoc builds a FunctionType from a FUNCTION node's
// JSDoc (@param / @return). Returns nil when the node doesn't carry
// enough JSDoc to infer a full signature.
func functionTypeFromJSDoc(info *rhino.JSDocInfo, fn *rhino.Node, scope *TypedScope) *jstype.FunctionType {
	kids := fn.Children()
	if len(kids) < 2 {
		return nil
	}
	params := kids[1]
	var paramTypes []jstype.JSType
	for _, p := range params.Children() {
		var pt jstype.JSType = jstype.Unknown
		if info != nil {
			if expr := info.ParameterType(p.String_()); expr != nil {
				pt = scope.ResolveExpr(expr.TypeName)
			}
		}
		paramTypes = append(paramTypes, pt)
	}
	var ret jstype.JSType = jstype.Unknown
	if info != nil && info.ReturnType() != nil {
		ret = scope.ResolveExpr(info.ReturnType().TypeName)
	}
	return jstype.NewFunctionType(paramTypes, ret)
}
