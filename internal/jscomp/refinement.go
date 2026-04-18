// Port of fragment of com/google/javascript/jscomp/TypeInference.java +
// com/google/javascript/jscomp/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415 (narrow; no AND/OR sub-refinement,
// no aliasing, no .property refinement).
//
// Java's SemanticReverseAbstractInterpreter walks the boolean
// expression tree of an if/hook condition and produces a refined
// FlowScope for the truthy/falsy branch. Our port turns the same
// condition shapes into a flat `name → JSType` map that TypeCheck
// pushes onto the active TypedScope while traversing the branch.
// Recognised condition forms:
//
//   x                                     – truthy branch strips falsy from x
//   x == null / x === null               – truthy keeps null/undef on x
//   x != null / x !== null               – truthy strips null/undef on x
//   typeof x === 'string'/'number'/...   – truthy narrows x to the named primitive
//   x instanceof Foo                      – truthy narrows x to Foo's instance type
//   !cond                                  – flips the refinement direction
//   cond1 && cond2                        – combines cond1 and cond2 in the truthy branch
//
// More elaborate AND / OR / parenthesised-conjunction handling lands
// with the full SRAI port; the current set covers the most common
// guard shapes the corpus exercises.

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
		// On the truthy branch BOTH conditions hold; on the falsy
		// branch we don't know which one failed, so we conservatively
		// emit nothing.
		if !whenTruthy {
			return
		}
		kids := cond.Children()
		if len(kids) != 2 {
			return
		}
		collectRefinements(scope, kids[0], true, out)
		collectRefinements(scope, kids[1], true, out)
		return
	case rhino.OR:
		// On the falsy branch BOTH sides are falsy.
		if whenTruthy {
			return
		}
		kids := cond.Children()
		if len(kids) != 2 {
			return
		}
		collectRefinements(scope, kids[0], false, out)
		collectRefinements(scope, kids[1], false, out)
		return
	case rhino.NAME:
		if whenTruthy {
			if t := scope.GetDeclaredType(cond.String_()); t != nil {
				if narrowed := stripFalsy(t); narrowed != nil {
					out[cond.String_()] = narrowed
				}
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
