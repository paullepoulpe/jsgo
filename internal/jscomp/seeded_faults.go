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
	// FunctionArgumentInjector ignores an argument's side effects when it
	// decides whether the argument needs a temporary, so a call's argument is
	// substituted into the inlined body as it is: dropped when the parameter
	// is unused, evaluated twice when it is used twice (a flipped inlining
	// condition).
	faultInlineArgs
	// RemoveUnusedCode deletes `var x = f()` without keeping the call (a
	// dropped side effect).
	faultDropUnusedInit
	// PeepholeFoldConstants folds an arithmetic result of -0 to 0.
	faultFoldNegZero
	// PeepholeMinimizeConditions drops the leading NOT of a ?: condition
	// without swapping the branches.
	faultHookNegation
	// CollapseProperties flattens names whose collapse is unsafe; this check
	// sits in the fast-mode-only early flattening.
	faultCollapseUnsafe
	// InlineVariables moves a variable's value to its only use without
	// checking that nothing in between has side effects (LocalVarMotion).
	faultInlineVars
	// FunctionToBlockMutator declares the temporaries of an inlined call's
	// arguments in reverse order (a swapped evaluation order).
	faultSwapArgs
)

var seededFaultNames = map[string]seededFaultKind{
	"inline-args":      faultInlineArgs,
	"drop-unused-init": faultDropUnusedInit,
	"fold-negzero":     faultFoldNegZero,
	"hook-negation":    faultHookNegation,
	"collapse-unsafe":  faultCollapseUnsafe,
	"inline-vars":      faultInlineVars,
	"swap-args":        faultSwapArgs,
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
