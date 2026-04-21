// Port of com/google/javascript/jscomp/TypeCheck.java's
// deterministic-test branch @ closure-compiler v20260415 (typecheck
// helper split out for readability).
//
// TypeCheck helper: JSC_DETERMINISTIC_TEST emission.
//
// Upstream Closure does NOT ship a standalone `CheckDeterministicTest`
// class — the diagnostic is emitted inline from TypeCheck.java's
// binary-comparison branches (EQ / NE / SHEQ / SHNE). This file hosts
// the helper logic called from typecheck.go's visit dispatch; it's
// split out for readability, not because it's a separate pass.
//
// Diagnostic:
//
//   JSC_DETERMINISTIC_TEST
//     condition always evaluates to <truth>
//     left : <leftType>
//     right: <rightType>
//
// Fires when the two operands of `==` / `!=` / `===` / `!==` are
// known at compile time to always compare the same way. For `==`
// / `!=` we defer to Tri-valued `testForEquality` on the
// restrict-by-not-null-or-undefined'd operand types (mirrors
// TypeCheck.java's `leftTypeRestricted.testForEquality(...)` call).
// For `===` / `!==` we defer to `canTestForShallowEqualityWith`:
// disjoint primitive slots → always false under `===` (always true
// under `!==`).

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// deterministicTestState is the scoped receiver for the helper
// below. One instance is built per binary-comparison visit inside
// TypeCheck; callers capture the returned diagnostics via
// VisitDeterministicTest.
type deterministicTestState struct {
	path  string
	scope *TypedScope
	diags []diag.Diagnostic
}

// VisitDeterministicTest is the entry point invoked from TypeCheck's
// EQ / NE / SHEQ / SHNE case. Matches Java's inline
// checkForDeterministicTest call inside TypeCheck.visitBinaryOp.
func VisitDeterministicTest(path string, scope *TypedScope, n *rhino.Node) []diag.Diagnostic {
	c := &deterministicTestState{path: path, scope: scope}
	c.visitEquality(n)
	return c.diags
}

// visitEquality ports the JSC_DETERMINISTIC_TEST branch of
// TypeCheck.visitBinaryOp: when the two operands of `==` / `!=` /
// `===` / `!==` are guaranteed to always compare the same way (the
// Java Tri-valued `testForEquality` / `canTestForShallowEqualityWith`
// results collapse to TRUE or FALSE), we emit at the whole
// expression's source position. Unknown / All / union operand types
// stay silent — those correspond to the `Tri.UNKNOWN` result in
// Java, which is not reported.
func (c *deterministicTestState) visitEquality(n *rhino.Node) {
	kids := n.Children()
	if len(kids) < 2 {
		return
	}
	infer := NewTypeInference(c.scope)
	l := infer.InferType(kids[0])
	r := infer.InferType(kids[1])
	strict := n.Token() == rhino.SHEQ || n.Token() == rhino.SHNE
	det, equalValue := detDisjointForEquality(l, r, strict)
	if !det {
		return
	}
	// Compute the final boolean for this operator. `equalValue` is
	// the result under `==` / `===`; `!=` / `!==` flip it. Matches
	// TypeCheck.java's `if (n.isNE()) { result = result.not(); }`.
	value := equalValue
	if n.Token() == rhino.NE || n.Token() == rhino.SHNE {
		value = !value
	}
	valueStr := "false"
	if value {
		valueStr = "true"
	}
	c.diags = append(c.diags, diag.Diagnostic{
		Path: c.path, Line: n.Lineno(), Col: n.Charno(),
		Code: diag.JSC_DETERMINISTIC_TEST, Level: diag.Error,
		Message: fmt.Sprintf(
			"condition always evaluates to %s\nleft : %s\nright: %s",
			valueStr, l.TypeName(), r.TypeName()),
	})
}

// detDisjointForEquality reports whether types `a` and `b` always
// produce the same result under `==` / `===` / `!=` / `!==`. JS's
// `==` coerces across number / string / boolean; function ↔ string
// / number, symbol ↔ anything-else, etc. are always-false. Single-
// inhabitant types (null, undefined) compared against themselves
// always succeed. Conservative: stay silent on Unknown / All / unions
// that might overlap — that matches Java's Tri.UNKNOWN fallback.
//
// Returns (deterministic, valueIfEqualOperator):
//   - `deterministic` true when the pair is disjoint OR both sides
//     are the same single-inhabitant type.
//   - `valueIfEqualOperator` is the comparison result assuming the
//     operator is `==` or `===` (`true` for same-singleton, `false`
//     for disjoint). Callers flip for `!=` / `!==`.
func detDisjointForEquality(a, b jstype.JSType, strict bool) (bool, bool) {
	if a == nil || b == nil {
		return false, false
	}
	if isUnknownOrNil(a) || isUnknownOrNil(b) {
		return false, false
	}
	// Unions: if one side is a nullable wrapper (e.g. `(!Symbol|null)`),
	// strip the null / undefined projection so the underlying type
	// participates in the disjointness check. Mirrors Java's
	// restrictByNotNullOrUndefined before testForEquality. Guarded:
	// don't strip `a`'s null arm when `b` itself is null / undefined,
	// because the comparison `X|null == null` is genuinely ambiguous
	// (the null arm might match) and Java's testForEquality reflects
	// that by returning UNKNOWN. Same guard on the symmetric side.
	if !isNullOrVoid(b) {
		a = detStripNullUndefIfWrapper(a)
	}
	if !isNullOrVoid(a) {
		b = detStripNullUndefIfWrapper(b)
	}
	if _, ok := a.(*jstype.UnionType); ok {
		return false, false
	}
	if _, ok := b.(*jstype.UnionType); ok {
		return false, false
	}
	ka := detClassifyEquality(a)
	kb := detClassifyEquality(b)
	if ka == detEqUnknown || kb == detEqUnknown {
		return false, false
	}
	if ka == kb {
		// Same single-inhabitant kind → always equal.
		if ka == detEqNull || ka == detEqVoid {
			return true, true
		}
		return false, false
	}
	// `null == undefined` is true; keep them quiet since `==` still
	// succeeds across this pair.
	if (ka == detEqNull && kb == detEqVoid) || (ka == detEqVoid && kb == detEqNull) {
		return false, false
	}
	// number / string / boolean coerce to each other under `==`.
	if detIsPrimitiveKind(ka) && detIsPrimitiveKind(kb) {
		return false, false
	}
	// Symbol primitive vs a wrapper Object/Function: under `==`,
	// Object can box a primitive via `valueOf`, so the comparison
	// might succeed. Only `===` is deterministically false. Matches
	// Java's testForEquality which returns UNKNOWN for `symbol == Object`
	// but `canTestForShallowEqualityWith` returns false for both sides.
	if (ka == detEqSymbol && (kb == detEqObject || kb == detEqFunction)) ||
		(kb == detEqSymbol && (ka == detEqObject || ka == detEqFunction)) {
		if !strict {
			return false, false
		}
	}
	// Everything else — function / object vs primitive, or distinct
	// function-vs-object, or symbol vs anything else — is always
	// false under `==` / `===`.
	return true, false
}

func detIsPrimitiveKind(k detEqualityKind) bool {
	switch k {
	case detEqNumber, detEqString, detEqBoolean:
		return true
	}
	return false
}

// detEqualityKind is the coarse equality-lattice bucket used for the
// disjointness check. Mirrors JSType.testForEquality's case analysis
// at just enough granularity to reproduce Closure's always-true /
// always-false decisions.
type detEqualityKind int

const (
	detEqUnknown detEqualityKind = iota
	detEqNumber
	detEqString
	detEqBoolean
	detEqNull
	detEqVoid
	detEqFunction
	detEqObject
	detEqSymbol
)

// detClassifyEquality returns the detEqualityKind for `t`. Unknown /
// All / templates over unknown map to detEqUnknown so the caller
// declines to emit.
func detClassifyEquality(t jstype.JSType) detEqualityKind {
	switch t {
	case jstype.Number:
		return detEqNumber
	case jstype.String:
		return detEqString
	case jstype.Boolean:
		return detEqBoolean
	case jstype.Null:
		return detEqNull
	case jstype.Void:
		return detEqVoid
	case jstype.Symbol:
		return detEqSymbol
	}
	if _, ok := t.(*jstype.FunctionType); ok {
		return detEqFunction
	}
	if obj, ok := t.(*jstype.ObjectType); ok {
		return detEqObjectKind(obj)
	}
	if tt, ok := t.(*jstype.TemplatizedType); ok && tt.Base != nil {
		return detEqObjectKind(tt.Base)
	}
	if et, ok := t.(*jstype.EnumElementType); ok && et.Enum != nil && et.Enum.Underlying != nil {
		// `E<Symbol>` where E is an enum of Symbol values behaves
		// like its underlying type for equality disjointness. Matches
		// how Java's testForEquality unwraps an EnumElementType to its
		// primitiveType before comparing.
		return detClassifyEquality(et.Enum.Underlying)
	}
	return detEqUnknown
}

// detEqObjectKind collapses every distinct nominal ObjectType to a
// single `detEqObject` bucket for coarse disjointness checks — good
// enough for `number == Date` / `Function == string` style
// always-false comparisons without enumerating every extern.
func detEqObjectKind(_ *jstype.ObjectType) detEqualityKind { return detEqObject }

// isNullOrVoid reports whether `t` is exactly the Null or Void
// primitive. Used to guard nullable-wrapper stripping — when the
// other operand is already null / undefined, the union's null arm
// might match at runtime so the comparison isn't deterministic.
func isNullOrVoid(t jstype.JSType) bool { return t == jstype.Null || t == jstype.Void }

// detStripNullUndefIfWrapper returns the underlying non-null /
// non-undefined member of `t` when `t` is a nullable wrapper like
// `!Symbol|null` / `T|undefined`. Returns `t` unchanged for
// non-union types and for unions that have more than one non-null
// member — the comparison stays silent in that case to match
// Java's UNKNOWN result for overlapping multi-member unions.
func detStripNullUndefIfWrapper(t jstype.JSType) jstype.JSType {
	u, ok := t.(*jstype.UnionType)
	if !ok {
		return t
	}
	var kept []jstype.JSType
	for _, m := range u.Members() {
		if m == jstype.Null || m == jstype.Void {
			continue
		}
		kept = append(kept, m)
	}
	if len(kept) == 1 {
		return kept[0]
	}
	return t
}
