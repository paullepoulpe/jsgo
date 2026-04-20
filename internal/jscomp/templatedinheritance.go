// Port of com/google/javascript/jscomp/TypedScopeCreator.java's
// handling of `@extends {Foo<T>}` on constructors
// @ closure-compiler v20260415 (narrow; post-scope wiring only).
//
// The Java TypedScopeCreator resolves a templatized @extends clause at
// constructor-type creation time, stashing both the nominal prototype
// link and the full TemplatizedType on the InstanceType so subsequent
// property lookups apply the substitution. Our scope.go currently does
// only the nominal half (and only when @extends carries no brackets);
// this pass fills in the templatized half across the script so
// TypeCheck, MissingReturnCheck, and SuperclassCheck all observe a
// complete inheritance graph.
//
// It also seeds InstanceType.TemplateParams from the constructor's
// FunctionType.TemplateTypes — again, something scope.go's Pass-2 does
// on the FunctionType but not the Instance, so the substitution map
// built inside TemplatizedType.GetPropertyType has the required keys.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// WireTemplatizedInheritance walks every constructor declaration in
// `script` and, when its @extends carries template args, wires the
// Instance's templatized supertype + nominal super link. Idempotent:
// when scope.go already wired the super the call is a no-op.
func WireTemplatizedInheritance(script *rhino.Node, scope *TypedScope) {
	if script == nil {
		return
	}
	// First pass: seed every generic class's Instance.TemplateParams
	// from its constructor's FunctionType.TemplateTypes.
	walkConstructors(script, func(name string, _ *rhino.JSDocInfo) {
		if fn, ok := scope.GetRawDeclaredType(name).(*jstype.FunctionType); ok && fn.Instance != nil {
			if len(fn.Instance.TemplateParams) == 0 && len(fn.TemplateTypes) > 0 {
				fn.Instance.TemplateParams = fn.TemplateTypes
			}
		}
	})
	// Second pass: wire @extends. Runs after the seed so the base's
	// TemplateParams are already populated when we build the
	// TemplatizedType.
	walkConstructors(script, func(name string, info *rhino.JSDocInfo) {
		if info == nil || info.BaseType() == nil {
			return
		}
		fn, ok := scope.GetRawDeclaredType(name).(*jstype.FunctionType)
		if !ok || fn.Instance == nil {
			return
		}
		resolved := scope.ResolveExpr(info.BaseType().TypeName)
		tpl := unwrapTemplatizedSuper(resolved)
		if tpl == nil || tpl.Base == nil {
			return
		}
		// Fill the nominal super link when scope.go's
		// `LookupNominalType(base.TypeName)` lookup failed (brackets
		// present). Never overwrite an existing wiring.
		if fn.Instance.Super() == nil {
			fn.Instance.SetSuper(tpl.Base)
		}
		if fn.Instance.TemplatizedSuper() == nil {
			fn.Instance.SetTemplatizedSuper(tpl)
		}
	})
}

// unwrapTemplatizedSuper extracts the TemplatizedType out of a
// @extends-resolved JSType, which is typically wrapped in a legacy-
// nullable union (e.g. `Foo<string>` resolves to `(Foo<string>|null)`).
func unwrapTemplatizedSuper(t jstype.JSType) *jstype.TemplatizedType {
	switch v := t.(type) {
	case *jstype.TemplatizedType:
		return v
	case *jstype.UnionType:
		for _, m := range v.Members() {
			if tpl, ok := m.(*jstype.TemplatizedType); ok {
				return tpl
			}
		}
	}
	return nil
}

// walkConstructors visits every top-level `function Name(…)` and
// `var Name = function(…)` node, invoking `fn` with the declared name
// and the owning JSDocInfo. Class declarations are out of scope: the
// templatized @extends shape only shows up on legacy ctor-functions in
// the extracted corpus.
func walkConstructors(n *rhino.Node, fn func(name string, info *rhino.JSDocInfo)) {
	if n == nil {
		return
	}
	for _, c := range n.Children() {
		switch c.Token() {
		case rhino.FUNCTION:
			if first := c.FirstChild(); first != nil && first.String_() != "" {
				fn(first.String_(), c.JSDocInfo())
			}
		case rhino.VAR, rhino.LET, rhino.CONST:
			nm := c.FirstChild()
			if nm != nil && nm.String_() != "" {
				if init := nm.FirstChild(); init != nil && init.Token() == rhino.FUNCTION {
					info := c.JSDocInfo()
					if info == nil {
						info = init.JSDocInfo()
					}
					fn(nm.String_(), info)
				}
			}
		}
		walkConstructors(c, fn)
	}
}
