package jscomp

// Go-only: seeded faults, the deliberate bugs that check the behavioural
// oracle (tools/behave/behave.py faults, docs/behavioural-oracle.md).
//
// They exist only in a binary built with `-tags jsgo_faults`
// (seeded_faults_on.go): there, JSGO_FAST_FAULT=<name> turns on one fault, in
// fast mode only, and the oracle rebuilds its tests' fast builds with it and
// reports whether running them catches the bug (the mutation score). In every
// other build (seeded_faults_off.go) seededFault is the constant faultNone, so
// each fault check is false at compile time and the environment is never read.

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
	// Benign variants: output changes that keep behaviour, which the oracle
	// must pass (or report as blind spots), never fail.
	//
	// RenameVars and RenameProperties skip the first name they would assign,
	// so every renamed name shifts: a pure renaming, which the coverage rule's
	// diff must call alpha-isomorphic.
	faultBenignRename
	// InlineFunctions inlines nothing: a different, correct decision at every
	// call site it would have inlined.
	faultBenignNoInline
)

var seededFaultNames = map[string]seededFaultKind{
	"inline-args":      faultInlineArgs,
	"drop-unused-init": faultDropUnusedInit,
	"fold-negzero":     faultFoldNegZero,
	"hook-negation":    faultHookNegation,
	"collapse-unsafe":  faultCollapseUnsafe,
	"inline-vars":      faultInlineVars,
	"swap-args":        faultSwapArgs,
	"benign-rename":    faultBenignRename,
	"benign-no-inline": faultBenignNoInline,
}

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

// fastFault reports whether fault f is on for a compilation with options o
// (always false without the jsgo_faults build tag).
func fastFault(o *CompilerOptions, f seededFaultKind) bool {
	return seededFault == f && o.isFastMode()
}
