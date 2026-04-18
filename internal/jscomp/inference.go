// Port of com/google/javascript/jscomp/TypeInference.java @ closure-compiler v20260415
//
// The Java pass does forward dataflow analysis assigning a JSType to every
// expression node. Tier-0 only needs inference for the initializer of a
// single var-decl, which is always a literal — so we implement the one
// case that matters.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

type TypeInference struct {
	scope *TypedScope
}

func NewTypeInference(scope *TypedScope) *TypeInference { return &TypeInference{scope: scope} }

// InferType mirrors the result of TypeInference for expression nodes.
// Returns Unknown for tokens this scaffold doesn't handle yet.
func (t *TypeInference) InferType(n *rhino.Node) jstype.JSType {
	switch n.Token() {
	case rhino.STRINGLIT:
		return jstype.String
	case rhino.NUMBER:
		return jstype.Number
	case rhino.TRUE, rhino.FALSE:
		return jstype.Boolean
	case rhino.NULL:
		return jstype.Null
	case rhino.UNDEFINED:
		return jstype.Void
	case rhino.NAME:
		// Port of TypeInference.java's "N_NAME" case: a bare identifier
		// resolves to the declared type of its binding. If the binding
		// isn't declared or has no @type JSDoc, fall back to Unknown
		// (conservative — no warning).
		if declared := t.scope.GetDeclaredType(n.String_()); declared != nil {
			return declared
		}
		return jstype.Unknown
	case rhino.THIS:
		// Port of TypeInference.java's THIS handling — the type bound
		// to `this` in the lexically nearest function/constructor scope.
		if inst := t.scope.ThisType(); inst != nil {
			return inst
		}
		return jstype.Unknown
	case rhino.GETPROP:
		// GETPROP(receiver, STRINGLIT propName). The receiver's type
		// must expose the property, otherwise we stay Unknown rather
		// than emit (TypeCheck.visitGetProp owns the diagnostic).
		kids := n.Children()
		if len(kids) < 2 {
			return jstype.Unknown
		}
		if obj, ok := t.InferType(kids[0]).(*jstype.ObjectType); ok {
			if pt := obj.GetPropertyType(kids[1].String_()); pt != nil {
				return pt
			}
		}
		return jstype.Unknown
	case rhino.NEW:
		// `new F(args)` evaluates to F's instance type when F is a
		// known constructor; Unknown otherwise.
		if callee := n.FirstChild(); callee != nil && callee.Token() == rhino.NAME {
			if fn, ok := t.scope.GetDeclaredType(callee.String_()).(*jstype.FunctionType); ok && fn.Instance != nil {
				return fn.Instance
			}
		}
		return jstype.Unknown
	case rhino.CALL:
		// Direct call: the result type is the callee's @return.
		if callee := n.FirstChild(); callee != nil && callee.Token() == rhino.NAME {
			if fn, ok := t.scope.GetDeclaredType(callee.String_()).(*jstype.FunctionType); ok {
				return fn.Return
			}
		}
		return jstype.Unknown
	}
	return jstype.Unknown
}
