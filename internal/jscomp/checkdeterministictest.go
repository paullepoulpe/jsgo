// Port of JSC_DETERMINISTIC_TEST branch of
// com/google/javascript/jscomp/TypeCheck.java + TypeValidator.testEquality
// @ closure-compiler v20260415.
//
// Fires JSC_DETERMINISTIC_TEST on `==` / `!=` / `===` / `!==` when the
// two operand types can never compare equal (e.g. `number == string` under
// `===`, or `null == number`).
//
// Java's TypeCheck.visitCompare dispatches differently for the two
// operator groups:
//
//   EQ / NE   — TypeCheck calls leftType.testForEquality(rightType),
//               which honours JS coercion rules (number ↔ string ↔ boolean).
//
//   SHEQ / SHNE — TypeCheck calls
//               leftType.canTestForShallowEqualityWith(rightType),
//               the strict-equality variant that requires a non-empty
//               greatest common subtype.
//
// Both paths strip null/undefined before comparing (Java's
// restrictByNotNullOrUndefined) so `(string|null) === number` and
// `string === number` both fire.

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// CheckDeterministicTest walks the AST and fires JSC_DETERMINISTIC_TEST
// at EQ / NE / SHEQ / SHNE nodes whose operand types make the
// comparison always yield the same boolean.
type CheckDeterministicTest struct {
	path    string
	creator *TypedScopeCreator // to look up function-local scopes
	scope   *TypedScope        // current active scope (updated per function)
	diags   []diag.Diagnostic
}

// NewCheckDeterministicTest constructs a CheckDeterministicTest pass.
func NewCheckDeterministicTest(path string, creator *TypedScopeCreator, scope *TypedScope) *CheckDeterministicTest {
	return &CheckDeterministicTest{path: path, creator: creator, scope: scope}
}

// Check walks script and returns accumulated diagnostics.
func (c *CheckDeterministicTest) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.walk(script, c.scope)
	return c.diags
}

func (c *CheckDeterministicTest) walk(n *rhino.Node, scope *TypedScope) {
	if n == nil {
		return
	}
	// Mirror NodeTraversal: when descending into a FUNCTION, switch to its
	// local scope (which holds the declared parameter types).
	if n.Token() == rhino.FUNCTION {
		if local, ok := c.creator.Functions[n]; ok && local != nil {
			scope = local
		}
	}
	switch n.Token() {
	case rhino.EQ, rhino.NE, rhino.SHEQ, rhino.SHNE:
		c.diags = append(c.diags, c.visitEquality(n, scope)...)
	}
	for _, ch := range n.Children() {
		c.walk(ch, scope)
	}
}

func (c *CheckDeterministicTest) visitEquality(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	infer := NewTypeInference(scope)
	leftType := infer.InferType(kids[0])
	rightType := infer.InferType(kids[1])

	// Java strips null/undefined before comparing (restrictByNotNullOrUndefined).
	leftRestricted := cdtRestrictByNotNullOrUndefined(leftType)
	rightRestricted := cdtRestrictByNotNullOrUndefined(rightType)

	var result cdtTri
	switch n.Token() {
	case rhino.EQ:
		result = cdtTestForEquality(leftRestricted, rightRestricted)
	case rhino.NE:
		result = cdtTriNot(cdtTestForEquality(leftRestricted, rightRestricted))
	case rhino.SHEQ:
		result = cdtTestForShallowEquality(leftRestricted, rightRestricted)
	case rhino.SHNE:
		result = cdtTriNot(cdtTestForShallowEquality(leftRestricted, rightRestricted))
	default:
		return nil
	}

	if result == cdtUnknown {
		return nil
	}
	valueStr := "false"
	if result == cdtTrue {
		valueStr = "true"
	}
	return []diag.Diagnostic{{
		Path:  c.path,
		Line:  n.Lineno(),
		Col:   n.Charno(),
		Code:  diag.JSC_DETERMINISTIC_TEST,
		Level: diag.Error,
		Message: fmt.Sprintf(
			"condition always evaluates to %s\nleft : %s\nright: %s",
			valueStr, leftType.TypeName(), rightType.TypeName()),
	}}
}

// cdtTri mirrors Java's com.google.javascript.rhino.jstype.TernaryValue (Tri).
type cdtTri int8

const (
	cdtUnknown cdtTri = 0
	cdtTrue    cdtTri = 1
	cdtFalse   cdtTri = -1
)

// cdtTriNot flips TRUE↔FALSE; UNKNOWN stays UNKNOWN.
func cdtTriNot(t cdtTri) cdtTri {
	switch t {
	case cdtTrue:
		return cdtFalse
	case cdtFalse:
		return cdtTrue
	}
	return cdtUnknown
}

// cdtRestrictByNotNullOrUndefined ports JSType.restrictByNotNullOrUndefined.
//
// Java's per-type overrides:
//   - NullType  → returns NO_TYPE (the empty / bottom type)
//   - VoidType  → returns NO_TYPE
//   - UnionType → strips null/undefined members
//   - everything else → returns self unchanged
//
// We represent NO_TYPE as nil; callers that receive nil must treat both
// sides as if they were the empty type (→ UNKNOWN, conservative).
func cdtRestrictByNotNullOrUndefined(t jstype.JSType) jstype.JSType {
	if t == nil {
		return nil // already empty
	}
	// Null and Undefined restrict to the empty type.
	if t == jstype.Null || t == jstype.Void {
		return nil
	}
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
	switch len(kept) {
	case 0:
		return nil // all members stripped → empty type
	case 1:
		return kept[0]
	default:
		return jstype.NewUnion(kept...)
	}
}

// cdtTestForEquality ports JSType.testForEquality and its per-type
// overrides for the `==` / `!=` operators.
//
// Java's logic (simplified):
//   - All, Unknown, NoResolved → UNKNOWN (conservative)
//   - Both empty types → TRUE; one empty → UNKNOWN
//   - Union on either side → reduce over members
//   - FunctionType vs non-Object → FALSE; vs Object → UNKNOWN
//   - symbol vs (symbol|Symbol object) → UNKNOWN; vs anything else → FALSE
//   - NullType  : null==null/undefined → TRUE; else unknown if nullable, else FALSE
//   - VoidType  : undefined==null/undefined → TRUE; else unknown if nullable, else FALSE
//   - Number/String/Boolean: can coerce with each other and objects → UNKNOWN; else FALSE
//   - EnumElementType: delegates to underlying primitive
//   - ObjectType / default → null from helper (fall through)
func cdtTestForEquality(a, b jstype.JSType) cdtTri {
	// nil represents the empty type (NO_TYPE in Java), produced by
	// restrictByNotNullOrUndefined when the type IS null or undefined.
	// Java's testForEqualityHelper handles empty types:
	//   both empty → TRUE (both sides are the empty set: vacuously equal)
	//   one empty  → UNKNOWN (conservative)
	aEmpty := a == nil
	bEmpty := b == nil
	if aEmpty || bEmpty {
		if aEmpty && bEmpty {
			return cdtTrue // e.g. null == null → both restrict to NO_TYPE → TRUE
		}
		return cdtUnknown
	}
	// All, Unknown → conservative.
	if cdtIsWide(a) || cdtIsWide(b) {
		return cdtUnknown
	}
	// Union: reduce to the merge of each member's result.
	if au, ok := a.(*jstype.UnionType); ok {
		return cdtUnionTestForEquality(au, b)
	}
	if bu, ok := b.(*jstype.UnionType); ok {
		return cdtUnionTestForEquality(bu, a)
	}
	// EnumElementType: unwrap to underlying and re-evaluate.
	if ae, ok := a.(*jstype.EnumElementType); ok {
		return cdtTestForEquality(ae.Enum.Underlying, b)
	}
	if be, ok := b.(*jstype.EnumElementType); ok {
		return cdtTestForEquality(a, be.Enum.Underlying)
	}

	// FunctionType or the built-in Function ObjectType: comparable to
	// objects; incompatible with primitives. Mirrors Java's FunctionType
	// testForEquality path; the ObjectType{Function} case covers the
	// @type {Function} / @param {Function} annotation form where the
	// scope creator resolves `Function` to the built-in type rather than
	// a concrete function signature.
	if cdtIsFunctionLike(a) {
		return cdtFunctionVsType(b)
	}
	if cdtIsFunctionLike(b) {
		return cdtFunctionVsType(a)
	}

	// symbol primitive (Java SymbolType.testForEquality).
	if a == jstype.Symbol {
		return cdtSymbolVsType(b)
	}
	if b == jstype.Symbol {
		return cdtSymbolVsType(a)
	}

	// null (Java NullType.testForEquality).
	if a == jstype.Null {
		return cdtNullVsType(b)
	}
	if b == jstype.Null {
		return cdtNullVsType(a)
	}

	// undefined/void (Java VoidType.testForEquality).
	if a == jstype.Void {
		return cdtVoidVsType(b)
	}
	if b == jstype.Void {
		return cdtVoidVsType(a)
	}

	// Number/String/Boolean can coerce with each other and objects → UNKNOWN;
	// not with symbol (already handled above).
	if cdtIsCoerciblePrimitive(a) {
		return cdtCoercibleVsType(b)
	}
	if cdtIsCoerciblePrimitive(b) {
		return cdtCoercibleVsType(a)
	}

	// ObjectType / unknown kind → conservative.
	return cdtUnknown
}

// cdtUnionTestForEquality reduces testForEquality over union members:
// all members must agree on the same Tri; any disagreement → UNKNOWN.
func cdtUnionTestForEquality(u *jstype.UnionType, other jstype.JSType) cdtTri {
	var result cdtTri
	first := true
	for _, m := range u.Members() {
		r := cdtTestForEquality(m, other)
		if first {
			result = r
			first = false
		} else if r != result {
			return cdtUnknown
		}
	}
	if first {
		return cdtUnknown
	}
	return result
}

// cdtTestForShallowEquality ports JSType.canTestForShallowEqualityWith
// for the `===` / `!==` operators: returns TRUE/FALSE/UNKNOWN based on
// whether the two types share a non-empty greatest common subtype.
//
// Rule: two types can be === only if they could both hold the same
// runtime value.  Primitive and object types are always disjoint (no
// value is both a primitive and an object).  Same-type primitives are
// never disjoint.  Object types may share instances via inheritance.
//
// Returns:
//
//	cdtUnknown — could be equal or not (non-deterministic)
//	cdtFalse   — can never be === (fires the warning)
//	cdtTrue    — always === (not currently produced by this path; Java
//	             doesn't produce it here either for SHEQ)
func cdtTestForShallowEquality(a, b jstype.JSType) cdtTri {
	// nil = empty type (NO_TYPE). Java's canTestForShallowEqualityWith:
	//   if (isEmptyType() || that.isEmptyType()) {
	//       return isSubtypeOf(that) || that.isSubtypeOf(this);
	//   }
	// NO_TYPE.isSubtypeOf(X) = true always → returns true → can test →
	// UNKNOWN (no diagnostic).
	if a == nil || b == nil {
		return cdtUnknown
	}
	// All, Unknown → conservative.
	if cdtIsWide(a) || cdtIsWide(b) {
		return cdtUnknown
	}
	// Empty types — mirror Java: subtype check.
	if cdtIsEmpty(a) || cdtIsEmpty(b) {
		if jstype.IsSubtype(a, b) || jstype.IsSubtype(b, a) {
			return cdtUnknown
		}
		return cdtFalse
	}
	// EnumElementType: unwrap to underlying for strict comparison.
	if ae, ok := a.(*jstype.EnumElementType); ok {
		return cdtTestForShallowEquality(ae.Enum.Underlying, b)
	}
	if be, ok := b.(*jstype.EnumElementType); ok {
		return cdtTestForShallowEquality(a, be.Enum.Underlying)
	}
	// Union: reduce over members.
	if au, ok := a.(*jstype.UnionType); ok {
		return cdtUnionTestForShallowEquality(au, b)
	}
	if bu, ok := b.(*jstype.UnionType); ok {
		return cdtUnionTestForShallowEquality(bu, a)
	}

	// symbol primitive vs object → disjoint (strict types, no coercion).
	if a == jstype.Symbol {
		return cdtShallowSymbolVsType(b)
	}
	if b == jstype.Symbol {
		return cdtShallowSymbolVsType(a)
	}

	// Primitives vs object (or function): no common subtype → FALSE.
	aPrim := cdtIsPrimitive(a)
	bPrim := cdtIsPrimitive(b)
	if aPrim && !bPrim {
		return cdtFalse
	}
	if bPrim && !aPrim {
		return cdtFalse
	}
	// Two different primitives → no common subtype.
	if aPrim && bPrim && a != b {
		return cdtFalse
	}
	// Two object types (or two functions): conservatively UNKNOWN.
	return cdtUnknown
}

// cdtUnionTestForShallowEquality reduces cdtTestForShallowEquality over
// union members: any FALSE member that disagrees collapses to UNKNOWN.
func cdtUnionTestForShallowEquality(u *jstype.UnionType, other jstype.JSType) cdtTri {
	var result cdtTri
	first := true
	for _, m := range u.Members() {
		r := cdtTestForShallowEquality(m, other)
		if first {
			result = r
			first = false
		} else if r != result {
			return cdtUnknown
		}
	}
	if first {
		return cdtUnknown
	}
	return result
}

// --- type-classification helpers ----------------------------------------

// cdtIsWide reports whether t is the All (*) or Unknown (?) sentinel,
// which silences all determinism checks (conservative).
func cdtIsWide(t jstype.JSType) bool {
	if t == nil {
		return true
	}
	n := t.TypeName()
	return n == "*" || n == "?"
}

// cdtIsEmpty reports whether t is an uninhabited type (NoType). jsgo
// doesn't materialise a dedicated NoType; treat as false (conservative).
func cdtIsEmpty(_ jstype.JSType) bool { return false }

// cdtIsPrimitive reports whether t is a non-object, non-function primitive
// (number, string, boolean, null, undefined, symbol).
func cdtIsPrimitive(t jstype.JSType) bool {
	switch t {
	case jstype.Number, jstype.String, jstype.Boolean,
		jstype.Null, jstype.Void, jstype.Symbol:
		return true
	}
	return false
}

// cdtIsCoerciblePrimitive reports whether t is one of number/string/boolean —
// the primitives that JS `==` coerces between each other and with objects.
func cdtIsCoerciblePrimitive(t jstype.JSType) bool {
	return t == jstype.Number || t == jstype.String || t == jstype.Boolean
}

// cdtFunctionVsType implements the FunctionType branch of
// JSType.testForEquality.  Java computes
//   otherType.getGreatestSubtype(OBJECT_TYPE)
// and returns FALSE when the result is the empty type (i.e. `other` has
// no overlap with Object).  Primitives, null, undefined, and symbol all
// have no overlap with Object → FALSE.  Object / FunctionType / anything
// that IS or extends Object → UNKNOWN.
func cdtFunctionVsType(other jstype.JSType) cdtTri {
	if cdtIsWide(other) {
		return cdtUnknown
	}
	// Primitives (including null, undefined, symbol) have no common
	// subtype with Object → greatest-subtype is empty → FALSE.
	if cdtIsPrimitive(other) {
		return cdtFalse
	}
	// Object types and functions → UNKNOWN.
	return cdtUnknown
}

// cdtSymbolVsType implements Java SymbolType.testForEquality for `==`.
// A symbol primitive can == another symbol or a Symbol object (UNKNOWN);
// it cannot == any other primitive or null/undefined (FALSE).
func cdtSymbolVsType(other jstype.JSType) cdtTri {
	if cdtIsWide(other) {
		return cdtUnknown
	}
	if other == jstype.Symbol {
		return cdtUnknown
	}
	// Object type or function could hold a Symbol (canCastTo logic) → UNKNOWN.
	if cdtIsObject(other) {
		return cdtUnknown
	}
	// Primitive (not symbol) or null/undefined → FALSE.
	return cdtFalse
}

// cdtShallowSymbolVsType implements the symbol branch for `===`.
// Under strict equality a symbol primitive can only === another symbol
// primitive; object types (including Symbol wrapper) are disjoint.
func cdtShallowSymbolVsType(other jstype.JSType) cdtTri {
	if cdtIsWide(other) {
		return cdtUnknown
	}
	if other == jstype.Symbol {
		return cdtUnknown // same primitive type → might be ===
	}
	// Any object type (including Symbol wrapper, Object, …) is disjoint
	// from the symbol primitive under strict equality.
	if cdtIsObject(other) {
		return cdtFalse
	}
	// Different primitives → disjoint.
	return cdtFalse
}

// cdtNullVsType implements Java NullType.testForEquality.
// null == null/undefined → TRUE; null == nullable → UNKNOWN; else FALSE.
func cdtNullVsType(other jstype.JSType) cdtTri {
	if cdtIsWide(other) {
		return cdtUnknown
	}
	if other == jstype.Null {
		return cdtTrue
	}
	if other == jstype.Void {
		return cdtTrue // null == undefined is true in JS
	}
	if cdtIsNullable(other) {
		return cdtUnknown
	}
	return cdtFalse
}

// cdtVoidVsType implements Java VoidType.testForEquality.
// undefined == null/undefined → TRUE; else if nullable → UNKNOWN; else FALSE.
func cdtVoidVsType(other jstype.JSType) cdtTri {
	if cdtIsWide(other) {
		return cdtUnknown
	}
	if other == jstype.Void || other == jstype.Null {
		return cdtTrue
	}
	if cdtIsNullable(other) {
		return cdtUnknown
	}
	return cdtFalse
}

// cdtCoercibleVsType handles number/string/boolean testForEquality.
//
// Java's StringType (and NumberType, BooleanType) return UNKNOWN if `that`
// is a subtype of Object — i.e. they can coerce with objects via JS `==`
// coercion (e.g. `number == new Number(1)` is allowed). They return FALSE
// for symbol, and for null/undefined (which are handled via the
// restrict-by-not-null step before we get here, so in practice they arrive
// as nil).
//
// Ports Java NumberType/StringType/BooleanType.testForEquality exactly:
//   - Wide (All, Unknown) → UNKNOWN
//   - null / undefined / symbol → FALSE
//   - Object subtypes (including wrapper objects) → UNKNOWN
//   - same-or-other coercible primitives → UNKNOWN
func cdtCoercibleVsType(other jstype.JSType) cdtTri {
	if cdtIsWide(other) {
		return cdtUnknown
	}
	switch other {
	case jstype.Null, jstype.Void, jstype.Symbol:
		return cdtFalse
	}
	if cdtIsObject(other) {
		// Java: that.isSubtypeOf(OBJECT_TYPE) → UNKNOWN. Coercible primitives
		// can == wrapper objects (e.g. `number == new Number(x)`).
		return cdtUnknown
	}
	return cdtUnknown
}

// cdtIsObject reports whether t is an object type (ObjectType,
// PrototypeObjectType, InstanceObjectType, FunctionType, TemplatizedType,
// EnumElementType wrapping an object, …) — i.e. not a primitive.
// Mirrors Java JSType.isObject() which returns true for any ObjectType
// subclass. Used to decide whether a coercible primitive can == the
// other side (UNKNOWN) or whether it's disjoint (FALSE).
func cdtIsObject(t jstype.JSType) bool {
	if t == nil {
		return false
	}
	if _, ok := t.(*jstype.ObjectType); ok {
		return true
	}
	if _, ok := t.(*jstype.PrototypeObjectType); ok {
		return true
	}
	if _, ok := t.(*jstype.InstanceObjectType); ok {
		return true
	}
	if _, ok := t.(*jstype.FunctionType); ok {
		return true
	}
	if _, ok := t.(*jstype.TemplatizedType); ok {
		return true
	}
	if ee, ok := t.(*jstype.EnumElementType); ok {
		return cdtIsObject(ee.Enum.Underlying)
	}
	return false
}

// cdtIsNullable reports whether t (a non-union type) could hold null.
// For object types this is always true in Closure's legacy mode (nullable
// by default).  For primitives it is false.
func cdtIsNullable(t jstype.JSType) bool {
	return cdtIsObject(t)
}

// cdtIsFunctionLike reports whether t behaves like a function for
// deterministic-test purposes: either it IS a FunctionType or it is
// the built-in Function ObjectType (the @type {Function} form).
// Mirrors Java's path where both function signatures and the top-level
// Function type go through FunctionType.testForEquality.
func cdtIsFunctionLike(t jstype.JSType) bool {
	if _, ok := t.(*jstype.FunctionType); ok {
		return true
	}
	if obj, ok := t.(*jstype.ObjectType); ok && obj.Name == "Function" {
		return true
	}
	return false
}
