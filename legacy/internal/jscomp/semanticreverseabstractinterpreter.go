// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
//
// Java's SemanticReverseAbstractInterpreter walks the boolean
// expression tree of an if/hook condition and produces a refined
// FlowScope for the truthy/falsy branch. Our port turns the same
// condition shapes into a flat `name → JSType` map that TypeCheck
// pushes onto the active TypedScope while traversing the branch.
// Recognised condition forms:
//
//   x                                     – truthy branch strips falsy from x
//   x == null / x === null               – truthy keeps null/undef (== loose) or null (=== strict)
//   x != null / x !== null               – truthy strips null/undef (!=) or null (!==)
//   typeof x === 'string'/'number'/...   – truthy narrows x to the named primitive
//   x instanceof Foo                      – truthy narrows x to Foo's instance type;
//                                          falsy removes Foo from a union (SRAI caseInstanceOf)
//   'prop' in obj                         – truthy records obj.prop as Unknown (SRAI caseIn)
//   !cond                                  – flips the refinement direction
//   cond1 && cond2                        – combines cond1 and cond2 in the truthy branch

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino/jstype"
)

// SemanticReverseAbstractInterpreter mirrors the Java class of the same name.
// Java's class is a stateful refinement engine that holds a JSTypeRegistry
// reference and produces refined FlowScopes. The Go port operates over a
// per-call refinement map, so the receiver is intentionally minimal — it
// carries the active TypedScope so the caseXxx methods can match Java's
// signatures without threading the scope through every call.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
type SemanticReverseAbstractInterpreter struct {
	scope *TypedScope
}

// Go-only: language-adapter — Java's SRAI constructor takes a JSTypeRegistry
// and is invoked once per Compiler. The Go port wires SRAI per scope so a
// factory makes more sense than a stateful registry.
func NewSemanticReverseAbstractInterpreter(scope *TypedScope) *SemanticReverseAbstractInterpreter {
	return &SemanticReverseAbstractInterpreter{scope: scope}
}

// getPreciserScopeKnowingConditionOutcome produces the refinement map active
// on the truthy or falsy branch of `condition`. Mirrors Java's identically-
// named entry point. The Go port returns a flat name → JSType map instead of
// a FlowScope.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func (s *SemanticReverseAbstractInterpreter) getPreciserScopeKnowingConditionOutcome(
	condition *rhino.Node, whenTruthy bool,
) map[string]jstype.JSType {
	return buildRefinements(s.scope, condition, whenTruthy)
}

// caseEquality dispatches the equality refinement for a binary EQ/NE/SHEQ/SHNE
// condition. Wraps the null-check branch of collectRefinements so it matches
// the Java method signature.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func (s *SemanticReverseAbstractInterpreter) caseEquality(
	cond *rhino.Node, whenTruthy bool,
) map[string]jstype.JSType {
	out := map[string]jstype.JSType{}
	if cond == nil || !isNullCheck(cond) {
		return out
	}
	collectRefinements(s.scope, cond, whenTruthy, out)
	return out
}

// caseTypeOf narrows the operand of a `typeof x === 'prim'` comparison.
// Wraps the typeof branch of collectRefinements.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func (s *SemanticReverseAbstractInterpreter) caseTypeOf(
	cond *rhino.Node, whenTruthy bool,
) map[string]jstype.JSType {
	out := map[string]jstype.JSType{}
	if cond == nil {
		return out
	}
	if name, _, _ := typeofCheckTarget(cond); name == "" {
		return out
	}
	collectRefinements(s.scope, cond, whenTruthy, out)
	return out
}

// caseInstanceOf narrows the LHS of `x instanceof Foo`.
// Wraps the instanceof branch of collectRefinements.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func (s *SemanticReverseAbstractInterpreter) caseInstanceOf(
	cond *rhino.Node, whenTruthy bool,
) map[string]jstype.JSType {
	out := map[string]jstype.JSType{}
	if cond == nil || cond.Token() != rhino.INSTANCEOF {
		return out
	}
	collectRefinements(s.scope, cond, whenTruthy, out)
	return out
}

// caseIn handles the `'prop' in obj` truthy refinement. Mirrors Java's
// identically-named method that records obj.prop with Unknown type when the
// LHS is a string literal.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func (s *SemanticReverseAbstractInterpreter) caseIn(
	cond *rhino.Node, whenTruthy bool,
) map[string]jstype.JSType {
	out := map[string]jstype.JSType{}
	if cond == nil || cond.Token() != rhino.IN || !whenTruthy {
		return out
	}
	collectRefinements(s.scope, cond, whenTruthy, out)
	return out
}

// caseNameOrGetProp refines an `if (x)` or `if (obj.p)` truthy/falsy branch
// for a NAME or GETPROP condition. Wraps the NAME / GETPROP arms of
// collectRefinements.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func (s *SemanticReverseAbstractInterpreter) caseNameOrGetProp(
	cond *rhino.Node, whenTruthy bool,
) map[string]jstype.JSType {
	out := map[string]jstype.JSType{}
	if cond == nil {
		return out
	}
	if cond.Token() != rhino.NAME && cond.Token() != rhino.GETPROP {
		return out
	}
	collectRefinements(s.scope, cond, whenTruthy, out)
	return out
}

// caseAndOrNotShortCircuiting handles AND truthy / OR falsy refinements where
// both operands' constraints apply. Wraps the AND / OR-falsy branches of
// collectRefinements.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func (s *SemanticReverseAbstractInterpreter) caseAndOrNotShortCircuiting(
	left, right *rhino.Node, andOutcomeTruthy bool,
) map[string]jstype.JSType {
	out := map[string]jstype.JSType{}
	collectRefinements(s.scope, left, andOutcomeTruthy, out)
	collectRefinements(s.scope, right, andOutcomeTruthy, out)
	return out
}

// caseAndOrMaybeShortCircuiting handles AND falsy / OR truthy refinements
// where either operand may have determined the outcome. The Go port unions
// the per-arm refinements and widens to declared types when only one side
// refines a name. Wraps the OR-truthy branch of collectRefinements.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func (s *SemanticReverseAbstractInterpreter) caseAndOrMaybeShortCircuiting(
	left, right *rhino.Node, orOutcomeTruthy bool,
) map[string]jstype.JSType {
	leftRefs := buildRefinements(s.scope, left, orOutcomeTruthy)
	rightRefs := buildRefinements(s.scope, right, orOutcomeTruthy)
	out := map[string]jstype.JSType{}
	for name, lt := range leftRefs {
		if rt, ok := rightRefs[name]; ok {
			out[name] = jstype.NewUnion(lt, rt)
			continue
		}
		if declared := s.scope.GetDeclaredType(name); declared != nil {
			out[name] = jstype.NewUnion(lt, declared)
		}
	}
	for name, rt := range rightRefs {
		if _, seen := leftRefs[name]; seen {
			continue
		}
		if declared := s.scope.GetDeclaredType(name); declared != nil {
			out[name] = jstype.NewUnion(rt, declared)
		}
	}
	return out
}

// caseIsArray narrows a value passed to `Array.isArray(value)`.
// jsgo does not yet implement the Array.isArray refinement; method is a
// stub mirroring Java's caseIsArray signature so the SRAI surface is
// 1:1 with the Java class.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func (s *SemanticReverseAbstractInterpreter) caseIsArray(
	value *rhino.Node, whenTruthy bool,
) map[string]jstype.JSType {
	return map[string]jstype.JSType{}
}

// maybeRestrictName records a refinement for `node` when restrictedType
// differs from originalType. Mirrors Java's identically-named helper that
// branches the FlowScope on a restricted type. The Go port writes directly
// into the refinement map.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func (s *SemanticReverseAbstractInterpreter) maybeRestrictName(
	out map[string]jstype.JSType, node *rhino.Node, originalType, restrictedType jstype.JSType,
) {
	if restrictedType == nil || restrictedType == originalType {
		return
	}
	qname := qualifiedName(node)
	if qname == "" {
		return
	}
	out[qname] = restrictedType
}

// maybeRestrictTwoNames is the binary form of maybeRestrictName, used by
// caseEquality to refine both sides of `a == b`.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func (s *SemanticReverseAbstractInterpreter) maybeRestrictTwoNames(
	out map[string]jstype.JSType,
	left *rhino.Node, originalLeftType, restrictedLeftType jstype.JSType,
	right *rhino.Node, originalRightType, restrictedRightType jstype.JSType,
) {
	s.maybeRestrictName(out, left, originalLeftType, restrictedLeftType)
	s.maybeRestrictName(out, right, originalRightType, restrictedRightType)
}

// eq mirrors Java's eq merging function for `==`. The Go port stores a single
// JSType per name (no TypePair) so this returns the type that survives the
// equality check on the truthy branch.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func eq(a, b jstype.JSType) (jstype.JSType, jstype.JSType) {
	if a == nil || b == nil {
		return a, b
	}
	// `==` truthy: a and b are loosely equal — keep null/undefined on both.
	return a, b
}

// ne mirrors Java's ne merging function for `!=`.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func ne(a, b jstype.JSType) (jstype.JSType, jstype.JSType) {
	if a == nil || b == nil {
		return a, b
	}
	return refineNonNull(a), refineNonNull(b)
}

// sheq mirrors Java's sheq (strict equality) merging function for `===`.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func sheq(a, b jstype.JSType) (jstype.JSType, jstype.JSType) {
	return a, b
}

// shne mirrors Java's shne (strict inequality) merging function for `!==`.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func shne(a, b jstype.JSType) (jstype.JSType, jstype.JSType) {
	if a == nil || b == nil {
		return a, b
	}
	return stripNull(a), stripNull(b)
}

// ineq mirrors Java's ineq merging function for `<`, `<=`, `>`, `>=`.
// Strips undefined from both operands on the truthy branch.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func ineq(a, b jstype.JSType) (jstype.JSType, jstype.JSType) {
	return stripVoid(a), stripVoid(b)
}

// unwrap is a no-op in the Go port: Java unwraps RefinementTrackingFlowScope
// layers, but jsgo represents refinements as a plain map and never wraps the
// scope. Method is included for surface parity with the Java class.
//
// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
func unwrap(refs map[string]jstype.JSType) map[string]jstype.JSType {
	return refs
}

// Go-only: language-adapter — Java's SRAI returns a FlowScope; the Go port
// returns a flat name → JSType map that TypeCheck pushes onto the active
// TypedScope. buildRefinements is the entry point for that map-returning
// shape, with no direct Java analog.
func buildRefinements(scope *TypedScope, cond *rhino.Node, whenTruthy bool) map[string]jstype.JSType {
	out := map[string]jstype.JSType{}
	collectRefinements(scope, cond, whenTruthy, out)
	return out
}

// Go-only: language-adapter — flat dispatch routine that fans out into the
// equivalent of Java's caseEquality / caseTypeOf / caseInstanceOf / caseIn
// / caseAndOrNotShortCircuiting / caseAndOrMaybeShortCircuiting / caseNameOrGetProp
// branches. The Java SRAI splits these into separate methods on a stateful
// class; the Go port keeps a single dispatch so the refinement map can be
// mutated in place.
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
		kids := cond.Children()
		if len(kids) != 2 {
			return
		}
		if !whenTruthy {
			// Falsy branch of `a || b` — both sides falsy. Intersection
			// by taking sequential refinements (later overwrites earlier).
			collectRefinements(scope, kids[0], false, out)
			collectRefinements(scope, kids[1], false, out)
			return
		}
		// Truthy branch of `a || b` — EITHER side could be truthy.
		// Union the per-side truthy refinements for each common name.
		leftRefs := buildRefinements(scope, kids[0], true)
		rightRefs := buildRefinements(scope, kids[1], true)
		for name, lt := range leftRefs {
			if rt, ok := rightRefs[name]; ok {
				out[name] = jstype.NewUnion(lt, rt)
				continue
			}
			// Right side didn't refine `name` — on that branch `name`
			// keeps its declared type, so the overall union widens to
			// (lt | declared).
			if declared := scope.GetDeclaredType(name); declared != nil {
				out[name] = jstype.NewUnion(lt, declared)
			}
		}
		for name, rt := range rightRefs {
			if _, seen := leftRefs[name]; seen {
				continue
			}
			if declared := scope.GetDeclaredType(name); declared != nil {
				out[name] = jstype.NewUnion(rt, declared)
			}
		}
		return
	case rhino.NAME:
		if t := scope.GetDeclaredType(cond.String_()); t != nil {
			if whenTruthy {
				// Port of SRAI.caseNameShortCircuiting truthy branch:
				// restrict to the members that can produce a truthy
				// value. Uses jstype.GetRestrictedTypeGivenOutcome
				// which applies GetPossibleToBooleanOutcomes per member
				// — more precise than the ad-hoc stripFalsy.
				if narrowed := refineTruthy(t); narrowed != nil {
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
		// `if (obj.prop)` narrows obj.prop to its truthy form
		// inside the then-branch (and falsy on the else). Keyed
		// by the qualified-name string "obj.prop" so inferGetProp
		// can consult the refinement.
		qname := qualifiedName(cond)
		if qname == "" {
			return
		}
		t := resolveQualifiedType(scope, cond)
		if t == nil {
			return
		}
		if whenTruthy {
			if narrowed := refineTruthy(t); narrowed != nil {
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
		name, equalToNull, strict, nullLit := nullCheckTarget(cond)
		if name == "" {
			return
		}
		// For GETPROP qualified names (e.g. "x.y"), we can't look up the
		// declared type from the scope chain. Instead, record a refinement
		// that suppresses inexistent-property errors when the check shows
		// the property is non-null/non-undefined. Port of SRAI caseEquality
		// for GETPROP via inferQualifiedSlot → Unknown.
		if strings.Contains(name, ".") {
			isNull := equalToNull == whenTruthy
			if !isNull {
				// `x.y != null` or `x.y !== null` truthy: property exists.
				out[name] = jstype.Unknown
			}
			return
		}
		t := scope.GetDeclaredType(name)
		if t == nil {
			return
		}
		// `x === null` truthy → keep only null; `x !== null` truthy → strip null.
		// `x == null` truthy → keep null+undefined; `x != null` truthy → strip both.
		// Port of SemanticReverseAbstractInterpreter.caseEquality / eq / sheq.
		isNull := equalToNull == whenTruthy
		if strict {
			if isNull {
				if nullLit {
					out[name] = keepNull(t)
				} else {
					out[name] = keepVoid(t)
				}
			} else {
				var narrowed jstype.JSType
				if nullLit {
					narrowed = stripNull(t)
				} else {
					narrowed = stripVoid(t)
				}
				if narrowed != nil {
					out[name] = narrowed
				}
			}
		} else {
			if isNull {
				out[name] = keepNullUndef(t)
			} else {
				// `x != null` truthy branch: restrict to non-null/non-undefined.
				// Uses jstype.RestrictByNotNullOrUndefined — the proper Java
				// port of JSType.restrictByNotNullOrUndefined, handling union
				// filtering and returning None for entirely-null types.
				narrowed := refineNonNull(t)
				if narrowed != nil {
					out[name] = narrowed
				}
			}
		}
		return
	}
	if name, primName, equal := typeofCheckTarget(cond); name != "" {
		matchesPositive := equal == whenTruthy
		// For GETPROP qualified names (e.g. "x.foo"), the declared type
		// is not in the scope chain; instead we use the refinement to
		// record that the property exists (Unknown) on the falsy branch
		// of `typeof x.foo == 'undefined'` (i.e. after the guard throws).
		// Port of SRAI caseTypeOf for GETPROP targets.
		if strings.Contains(name, ".") {
			if primName == "undefined" && !matchesPositive {
				// Negated: typeof x.foo != 'undefined' — property is defined.
				out[name] = jstype.Unknown
			}
			return
		}
		t := scope.GetDeclaredType(name)
		if t == nil {
			return
		}
		// `typeof x == 'object'` is a special case — narrows to
		// null plus object-ish members on the truthy branch, and
		// strips those on the falsy branch. Matches Java's
		// SemanticReverseAbstractInterpreter's 'object' case.
		if primName == "object" {
			if matchesPositive {
				out[name] = keepObjectish(t)
			} else if narrowed := stripObjectish(t); narrowed != nil {
				out[name] = narrowed
			}
			return
		}
		if matchesPositive {
			// `typeof x === 'prim'` truthy: restrict to members matching
			// the typeof string. Uses jstype.RestrictByTypeof — the Java
			// port from ClosureReverseAbstractInterpreter which handles
			// "function", "symbol", bigint, etc. in addition to the basic
			// primitives that primitiveByName covers.
			if narrowed := refineTypeof(t, primName); narrowed != nil {
				out[name] = narrowed
			}
		} else {
			// `typeof x !== 'prim'` truthy: remove members matching the
			// typeof string from the union.
			prim := primitiveByName(primName)
			if prim == nil {
				return
			}
			if narrowed := removeMember(t, prim); narrowed != nil {
				out[name] = narrowed
			}
		}
		return
	}
	if name, target := instanceofCheckTarget(scope, cond); name != "" {
		// `x instanceof Foo` truthy → x narrows to Foo's instance type.
		// `x instanceof Foo` falsy → remove Foo's instance type from the
		// union when x is a union. Port of caseInstanceOf +
		// RestrictByFalseInstanceOfResultVisitor from SRAI.
		if whenTruthy {
			out[name] = target
		} else {
			t := scope.GetDeclaredType(name)
			if t != nil {
				if narrowed := removeFromUnion(t, target); narrowed != nil {
					out[name] = narrowed
				}
			}
		}
		return
	}
	// `'prop' in obj` (truthy) — add obj.prop as a refinement with
	// Unknown type so JSC_INEXISTENT_PROPERTY is suppressed inside the
	// branch. Port of SemanticReverseAbstractInterpreter.caseIn.
	if cond.Token() == rhino.IN && whenTruthy {
		kids := cond.Children()
		if len(kids) == 2 && kids[0].Token() == rhino.STRINGLIT {
			propName := kids[0].String_()
			obj := kids[1]
			if qname := qualifiedName(obj); qname != "" {
				out[qname+"."+propName] = jstype.Unknown
			}
		}
	}
}

// Go-only: language-adapter — Java calls
// JSType.getRestrictedTypeGivenOutcome(Outcome.TRUE) directly on the type.
// jsgo's jstype.GetRestrictedTypeGivenOutcome wraps None in NewUnion which
// widens rather than narrows, so we delegate to stripFalsy as a sound but
// less precise stand-in. Mirrors Java's caseNameOrGetProp truthy path.
func refineTruthy(t jstype.JSType) jstype.JSType {
	return stripFalsy(t)
}

// Go-only: language-adapter — wraps jstype.RestrictByNotNullOrUndefined,
// converting None → nil so callers can use nil as "no refinement". Java
// inlines `JSType.restrictByNotNullOrUndefined()` at the call site.
func refineNonNull(t jstype.JSType) jstype.JSType {
	if t == nil {
		return nil
	}
	result := jstype.RestrictByNotNullOrUndefined(t)
	if result == nil || result == jstype.None {
		return nil
	}
	return result
}

// Go-only: language-adapter — wraps jstype.RestrictByTypeof with union
// filtering and an Unknown fallback. Java's SRAI calls
// getRestrictedByTypeOfResult on the type directly; jsgo's typeof-restriction
// machinery is split between RestrictByTypeof and primitiveByName, and we
// stitch them together here.
func refineTypeof(t jstype.JSType, typeofValue string) jstype.JSType {
	if t == nil {
		return nil
	}
	// For union types: use RestrictByTypeof and filter None members.
	if u, ok := t.(*jstype.UnionType); ok {
		var kept []jstype.JSType
		for _, m := range u.Members() {
			r := jstype.RestrictByTypeof(m, typeofValue)
			if r == nil || r == jstype.None {
				continue
			}
			kept = append(kept, r)
		}
		if len(kept) == 0 {
			return nil
		}
		if len(kept) == 1 {
			return kept[0]
		}
		return jstype.NewUnion(kept...)
	}
	// For non-union (including Unknown): fall back to keepMember so that
	// Unknown narrows to the primitive type (matching SRAI's legacy path
	// for untyped variables). jstype.RestrictByTypeof would return t
	// unchanged for Unknown, losing the narrowing.
	prim := primitiveByName(typeofValue)
	if prim != nil {
		return keepMember(t, prim)
	}
	// For "function" and "symbol" (not in primitiveByName): use
	// RestrictByTypeof which handles FunctionType / SymbolType.
	result := jstype.RestrictByTypeof(t, typeofValue)
	if result == nil || result == jstype.None {
		return nil
	}
	return result
}

// qualifiedName returns a dotted string for a GETPROP chain rooted
// at a NAME or THIS, e.g. `t.x.y` → "t.x.y", `this.x` → "this.x".
// Returns "" for any chain that doesn't resolve to a pure
// name/this-rooted access (GETELEM, CALL receivers, etc.).
//
// Go-only: language-adapter — helper shared between
// semanticreverseabstractinterpreter.go and inference.go. Java uses
// Node.getQualifiedName() directly on rhino.Node; jsgo's rhino.Node
// has no equivalent method so we provide a free function.
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
// step, consulting scope's namespace-prop map for flat two-segment
// chains. Returns nil if any step is unknown.
//
// Go-only: language-adapter — helper used by collectRefinements for
// GETPROP condition narrowing. Java's SRAI uses StaticTypedScope.getSlot
// with a qualified name string; jsgo's TypedScope walks namespace-prop
// records so we resolve the chain manually.
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
		// Two-segment GETPROP rooted at a NAME that's been
		// registered as a namespace (via recordNamespaceAssign or
		// recordNamespaceStub) — consult the namespace-prop map
		// before descending into ObjectType properties.
		if kids[0].Token() == rhino.NAME {
			if t := scope.LookupNamespaceProp(kids[0].String_(), kids[1].String_()); t != nil {
				return t
			}
		}
		recv := resolveQualifiedType(scope, kids[0])
		if recv == nil {
			return nil
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
//
// Go-only: language-adapter — pattern-matching helper for caseInstanceOf.
// Java inlines the INSTANCEOF dispatch in getPreciserScopeKnowingConditionOutcome;
// jsgo extracts the pattern match into a typed helper so the call sites stay
// flat.
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
