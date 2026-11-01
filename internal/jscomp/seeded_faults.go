package jscomp

// Go-only: seeded faults, the deliberate bugs that check the behavioural
// oracle (tools/behave/behave.py faults, docs/behavioural-oracle.md).
//
// JSGO_FAST_FAULT=<name> turns on one fault, in fast mode only: the oracle
// rebuilds its tests' fast builds with it and reports whether running them
// catches the bug (the mutation score). Unset, the default, seededFault is
// faultNone: compat mode never reaches a fault, and in fast mode each fault
// site costs one comparison of a package variable that is never written after
// initialisation.

import "os"

type seededFaultKind uint8

const (
	faultNone seededFaultKind = iota
	// FunctionArgumentInjector substitutes an argument directly even when the
	// callee's own side effects may change its value first (a flipped
	// inlining condition).
	faultInlineArgs
	// RemoveUnusedCode deletes `var x = f()` without keeping the call (a
	// dropped side effect).
	faultDropUnusedInit
	// PeepholeSubstituteAlternateSyntax rotates `a * (b / c)` to `b / c * a`
	// without checking for side effects (a swapped evaluation order).
	faultSwapOperands
	// PeepholeFoldConstants folds an arithmetic result of -0 to 0.
	faultFoldNegZero
	// PeepholeMinimizeConditions drops the leading NOT of a ?: condition
	// without swapping the branches.
	faultHookNegation
	// CollapseProperties flattens names whose collapse is unsafe; this check
	// sits in the fast-mode-only early flattening.
	faultCollapseUnsafe
)

var seededFaultNames = map[string]seededFaultKind{
	"inline-args":      faultInlineArgs,
	"drop-unused-init": faultDropUnusedInit,
	"swap-operands":    faultSwapOperands,
	"fold-negzero":     faultFoldNegZero,
	"hook-negation":    faultHookNegation,
	"collapse-unsafe":  faultCollapseUnsafe,
}

var seededFault = parseSeededFault(os.Getenv("JSGO_FAST_FAULT"))

func parseSeededFault(name string) seededFaultKind {
	if name == "" {
		return faultNone
	}
	f, ok := seededFaultNames[name]
	if !ok {
		panic("JSGO_FAST_FAULT: unknown fault " + name)
	}
	return f
}

// fastFault reports whether fault f is on for a compilation with options o.
func fastFault(o *CompilerOptions, f seededFaultKind) bool {
	return seededFault == f && o.isFastMode()
}
