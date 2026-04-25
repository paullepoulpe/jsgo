// Port of com/google/javascript/jscomp/TypeCheck.java
// @ closure-compiler v20260415 (subset).
//
// Strict-mode operand checks for bitwise and shift operators. Java
// implements these inline in TypeCheck.visitBinaryOperator (the
// bitwise / shift cases at TypeCheck.java:2966-3015). Each operand is
// run through TypeValidator.expectBitwiseable / expectNumberStrict
// which fires JSC_INVALID_OPERAND_TYPE when the operand is bitwiseable
// (so default-mode JSC_TYPE_MISMATCH did not fire) but not a strict
// subtype of the native NUMBER type.
//
// Go-only: language-adapter — extracts the strict bitwise/shift case
// out of visit() into its own helper so the existing GroupStrictCheckTypes
// switch in TypeCheck.Visit only needs a one-line dispatch.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/diag"
	rhino "github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// visitStrictBitwise mirrors the strict path inside Java's
// TypeValidator.expectBitwiseable / expectNumberStrict for the
// BIT* and SHIFT* operators (and their compound-assignment variants).
//
// Behavior:
//   - Operands whose type is unknown / nil are skipped — Java's
//     isSubtypeOf(NUMBER_TYPE) accepts unknown.
//   - Operands whose type is NOT bitwiseable are skipped here: the
//     default-mode visitBitOp already reported JSC_TYPE_MISMATCH for
//     them, and Java's expectBitwiseable short-circuits before
//     expectNumberStrict in that case.
//   - Otherwise, if the type is not a subtype of NUMBER, fire
//     JSC_INVALID_OPERAND_TYPE.
//
// Diagnostic placement matches Java:
//   - For BITAND/BITOR/BITXOR (and ASSIGN_*): point at the operand
//     node (Java's expectBitwiseable passes the operand as `n`).
//   - For LSH/RSH/URSH (and ASSIGN_*): point at the operator node
//     (Java's strict shift path passes `n` — the operator — to
//     expectNumberStrict).
func (c *TypeCheck) visitStrictBitwise(n *rhino.Node, scope *TypedScope) []diag.Diagnostic {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	infer := NewTypeInference(scope)
	var out []diag.Diagnostic
	isShift := isShiftToken(n.Token())
	op := bitOpSymbol(n.Token())
	for i, k := range kids {
		t := infer.InferType(k)
		if t == nil || isUnknownOrNil(t) {
			continue
		}
		// Java's expectBitwiseable short-circuits to mismatch (already
		// reported by default-mode visitBitOp) before reaching the
		// strict check; replicate that gate here for BIT* operators.
		if !isShift && !isBitwiseable(t) {
			continue
		}
		// Java's shift case skips strict reporting when matchesNumberContext
		// is false (the default-mode BIT_OPERATION error was reported instead).
		if isShift && isDefinitelyNotNumber(t) {
			continue
		}
		if jstype.IsSubtype(t, jstype.Number) {
			continue
		}
		var msg string
		var line, col int
		if isShift {
			msg = "operator " + op
			line, col = n.Lineno(), n.Charno()
		} else {
			if i == 0 {
				msg = "bad left operand to bitwise operator"
			} else {
				msg = "bad right operand to bitwise operator"
			}
			line, col = k.Lineno(), k.Charno()
		}
		out = append(out, diag.Diagnostic{
			Path: c.path, Line: line, Col: col,
			Code: diag.JSC_INVALID_OPERAND_TYPE, Level: diag.Warning,
			Message: msg + "\nfound   : " + t.TypeName() + "\nrequired: number",
		})
	}
	return out
}

// Go-only: language-adapter — token classifier for visitStrictBitwise; Java's
// switch in visitBinaryOperator groups SHIFT tokens lexically.
func isShiftToken(t rhino.Token) bool {
	switch t {
	case rhino.LSH, rhino.RSH, rhino.URSH,
		rhino.ASSIGN_LSH, rhino.ASSIGN_RSH, rhino.ASSIGN_URSH:
		return true
	}
	return false
}
