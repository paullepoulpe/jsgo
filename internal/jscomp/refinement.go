// Port of fragment of com/google/javascript/jscomp/TypeInference.java +
// com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java +
// com/google/javascript/jscomp/type/ChainableReverseAbstractInterpreter.java
// @ closure-compiler v20260415 (narrow; no aliasing).
//
// Java's SemanticReverseAbstractInterpreter walks the boolean
// expression tree of an if/hook condition and produces a refined
// FlowScope for the truthy/falsy branch. Our port turns the same
// condition shapes into a flat `name → JSType` map that TypeCheck
// pushes onto the active TypedScope while traversing the branch.
// Recognised condition forms:
//
//   x                                     – truthy branch strips falsy from x
//   obj.prop                               – truthy narrows obj.prop slot
//   x == null / x === null                – truthy keeps null/undef on x
//   x != null / x !== null                – truthy strips null/undef on x
//   typeof x === 'string'/'number'/...    – truthy narrows x to the primitive
//   x instanceof Foo                       – truthy narrows x to Foo's instance
//   !cond                                   – flips the refinement direction
//   cond1 && cond2 / cond1 || cond2        – sub-refinement of both sides
//
// Nested AND / OR: when refining the right operand of an `a && b`
// truthy, `a` must be truthy (short-circuit). We push `a`'s truthy
// refinements before inspecting `b` so shapes like `a && a.foo` see
// the narrowed `a` when deciding whether `a.foo` is narrowable.
// Mirrors Java's ChainableReverseAbstractInterpreter.caseAndOrMaybeShortCircuiting.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// buildRefinements inspects `cond` and returns the name → type map
// that should be active on the truthy or falsy branch.
func buildRefinements(scope *TypedScope, cond *rhino.Node, whenTruthy bool) map[string]jstype.JSType {
	out := map[string]jstype.JSType{}
	collectRefinements(scope, cond, whenTruthy, out)
	return out
}

func collectRefinements(scope *TypedScope, cond *rhino.Node, whenTruthy bool, out map[string]jstype.JSType) {
	if cond == nil {
		return
	}
	switch cond.Token() {
	case rhino.NOT:
		// `!c` truthy ↔ `c` falsy.
		collectRefinements(scope, cond.FirstChild(), !whenTruthy, out)
		return
	case rhino.AND:
		kids := cond.Children()
		if len(kids) != 2 {
			return
		}
		// `a && b` on truthy = BOTH hold, on falsy = EITHER falsy.
		// `a || b` on truthy = EITHER truthy, on falsy = BOTH falsy.
		// `bothHold` is true when we're on the branch that requires
		// both operands to agree, false when we're on the union branch.
		bothHold := (cond.Token() == rhino.AND) == whenTruthy
		collectLogicRefinements(scope, kids[0], kids[1], whenTruthy, bothHold, out)
		return
	case rhino.OR:
		kids := cond.Children()
		if len(kids) != 2 {
			return
		}
		bothHold := (cond.Token() == rhino.AND) == whenTruthy
		collectLogicRefinements(scope, kids[0], kids[1], whenTruthy, bothHold, out)
		return
	case rhino.NAME:
		if t := scope.GetDeclaredType(cond.String_()); t != nil {
			if whenTruthy {
				if narrowed := stripFalsy(t); narrowed != nil {
					out[cond.String_()] = narrowed
				}
			} else {
				// On the falsy branch of `if (a)` with a : T|null,
				// the value must be one of the "falsy" members —
				// null / undefined / false. Non-union and
				// primitive-only types can't be precisely narrowed
				// here (an unknown `number` could be 0 or not), so
				// we only refine unions that contain null / void.
				if narrowed := keepFalsyRefined(t); narrowed != nil {
					out[cond.String_()] = narrowed
				}
			}
		}
		return
	case rhino.GETPROP:
		// `if (obj.prop)` narrows the obj.prop slot inside the then-
		// branch (and keeps only its falsy members in the else).
		// Refinements are keyed by the qualified-name string so
		// inferGetProp can consult GetQualifiedRefinement. Port of
		// SemanticReverseAbstractInterpreter's GETPROP case.
		qname := qualifiedName(cond)
		if qname == "" {
			return
		}
		t := resolveQualifiedType(scope, cond)
		if t == nil {
			return
		}
		if whenTruthy {
			if narrowed := stripFalsy(t); narrowed != nil {
				out[qname] = narrowed
			}
		} else {
			if narrowed := keepFalsyRefined(t); narrowed != nil {
				out[qname] = narrowed
			}
		}
		return
	}
	if isNullCheck(cond) {
		name, equalToNull := nullCheckTarget(cond)
		if name == "" {
			return
		}
		t := scope.GetDeclaredType(name)
		if t == nil {
			return
		}
		isNull := equalToNull == whenTruthy
		if isNull {
			out[name] = keepNullUndef(t)
		} else if narrowed := stripNullUndef(t); narrowed != nil {
			out[name] = narrowed
		}
		return
	}
	if name, primName, equal := typeofCheckTarget(cond); name != "" {
		prim := primitiveByName(primName)
		if prim == nil {
			return
		}
		t := scope.GetDeclaredType(name)
		if t == nil {
			return
		}
		matchesPositive := equal == whenTruthy
		if matchesPositive {
			out[name] = keepMember(t, prim)
		} else if narrowed := removeMember(t, prim); narrowed != nil {
			out[name] = narrowed
		}
		return
	}
	if name, target := instanceofCheckTarget(scope, cond); name != "" {
		// `x instanceof Foo` truthy → x narrows to Foo's instance.
		// On the falsy branch we can't easily express "anything except
		// Foo" so we leave x alone.
		if whenTruthy {
			out[name] = target
		}
		return
	}
}

// collectLogicRefinements folds the refinement rules for `a && b`
// and `a || b` into one helper. `whenTruthy` is the arm of the
// surrounding condition; `bothHold` tells us whether that arm
// requires both operands to agree (AND-truthy, OR-falsy) or either
// to suffice (AND-falsy, OR-truthy). Right-operand analysis always
// runs under the refinement that short-circuit guarantees:
//   - bothHold: right is evaluated only when left already took the
//     `whenTruthy` value, so push left's `whenTruthy` refinements.
//   - !bothHold: right is evaluated only when left took the opposite
//     value, so push left's `!whenTruthy` refinements.
// Mirrors Java's ChainableReverseAbstractInterpreter.caseAndOrMaybe
// ShortCircuiting + RestrictByAnd / RestrictByOr combined.
func collectLogicRefinements(scope *TypedScope, leftNode, rightNode *rhino.Node, whenTruthy, bothHold bool, out map[string]jstype.JSType) {
	leftMatch := map[string]jstype.JSType{}
	collectRefinements(scope, leftNode, whenTruthy, leftMatch)
	if bothHold {
		// Both operands agree with the outer branch. Push left's
		// refinements before inspecting right so nested chains like
		// `a && a.foo` see `a` narrowed. The combined refinement is
		// the per-side union (right wins on collisions — matches
		// Java's incremental overlay semantics).
		for k, v := range leftMatch {
			out[k] = v
		}
		scope.PushRefinement(leftMatch)
		defer scope.PopRefinement()
		collectRefinements(scope, rightNode, whenTruthy, out)
		return
	}
	// Either operand can satisfy the outer branch. Right runs only
	// when left took the opposite value, so push that before
	// inspecting right; then merge the two sides via union+declared.
	leftOpposite := buildRefinements(scope, leftNode, !whenTruthy)
	scope.PushRefinement(leftOpposite)
	rightMatch := buildRefinements(scope, rightNode, whenTruthy)
	scope.PopRefinement()
	mergeBranchRefinements(scope, leftMatch, rightMatch, out)
}

// mergeBranchRefinements combines two branch-refinement maps (e.g.
// the truthy-a and truthy-b maps for `a || b` on the truthy arm)
// into `out`. For each name refined on both sides, the result is
// the union of the two refinements; for names refined on one side
// only we widen to (refined ∪ declared) so the overall narrowing
// still covers the "other side runs with declared type" case.
// Mirrors Java's LinkedFlowScope.joinOp.apply for branch merges.
func mergeBranchRefinements(scope *TypedScope, left, right, out map[string]jstype.JSType) {
	for name, lt := range left {
		if rt, ok := right[name]; ok {
			out[name] = jstype.NewUnion(lt, rt)
			continue
		}
		if declared := scope.GetDeclaredType(name); declared != nil {
			out[name] = jstype.NewUnion(lt, declared)
		}
	}
	for name, rt := range right {
		if _, seen := left[name]; seen {
			continue
		}
		if declared := scope.GetDeclaredType(name); declared != nil {
			out[name] = jstype.NewUnion(rt, declared)
		}
	}
}

// qualifiedName returns a dotted string for a GETPROP chain rooted
// at a NAME or THIS, e.g. `t.x.y` → "t.x.y", `this.x` → "this.x".
// Returns "" for any chain that doesn't resolve to a pure
// name/this-rooted access (GETELEM, CALL receivers, etc.).
func qualifiedName(n *rhino.Node) string {
	if n == nil {
		return ""
	}
	switch n.Token() {
	case rhino.NAME:
		return n.String_()
	case rhino.THIS:
		return "this"
	case rhino.GETPROP:
		kids := n.Children()
		if len(kids) < 2 {
			return ""
		}
		left := qualifiedName(kids[0])
		if left == "" {
			return ""
		}
		return left + "." + kids[1].String_()
	}
	return ""
}

// resolveQualifiedType looks up the declared type for a GETPROP
// chain — walks from the NAME / THIS root through each property
// step. Returns nil if any step is unknown. Mirrors how Java's
// SemanticReverseAbstractInterpreter reads the current property
// type off the receiver's ObjectType before computing a refinement.
func resolveQualifiedType(scope *TypedScope, n *rhino.Node) jstype.JSType {
	if n == nil {
		return nil
	}
	switch n.Token() {
	case rhino.NAME:
		return scope.GetDeclaredType(n.String_())
	case rhino.THIS:
		if inst := scope.ThisType(); inst != nil {
			return inst
		}
		return nil
	case rhino.GETPROP:
		kids := n.Children()
		if len(kids) < 2 {
			return nil
		}
		recv := resolveQualifiedType(scope, kids[0])
		if recv == nil {
			return nil
		}
		// Receiver might be a constructor — prefer its static slot
		// over the instance's prototype property.
		if fn, ok := recv.(*jstype.FunctionType); ok && fn.Instance != nil {
			if pt := fn.Instance.GetStaticPropertyType(kids[1].String_()); pt != nil {
				return pt
			}
			recv = fn.Instance
		}
		if tplt, ok := recv.(*jstype.TemplatizedType); ok {
			recv = tplt.Base
		}
		if obj, ok := recv.(*jstype.ObjectType); ok {
			return obj.GetPropertyType(kids[1].String_())
		}
	}
	return nil
}

// instanceofCheckTarget recognises `NAME instanceof Ctor`. Returns
// the NAME and the constructor's instance type so the truthy branch
// can narrow to that instance. Returns ("", nil) on no match.
func instanceofCheckTarget(scope *TypedScope, cond *rhino.Node) (string, *jstype.ObjectType) {
	if cond == nil || cond.Token() != rhino.INSTANCEOF {
		return "", nil
	}
	kids := cond.Children()
	if len(kids) != 2 {
		return "", nil
	}
	left, right := kids[0], kids[1]
	if left.Token() != rhino.NAME || right.Token() != rhino.NAME {
		return "", nil
	}
	if inst := scope.LookupNominalType(right.String_()); inst != nil {
		return left.String_(), inst
	}
	return "", nil
}
