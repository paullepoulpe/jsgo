// Port of com/google/javascript/jscomp/Compiler.java
// @ closure-compiler v20260415 (shim; pipeline wrapper only).
//
// Closure's Compiler is a ~3000-line orchestrator (PhaseOptimizer,
// pass registry, Result, module graph, warning manager, …). This
// file is a ~30-line adapter that wires our three independent
// passes (parse → TypedScopeCreator → TypeCheck) for the CLI +
// harness. A real Compiler port — pass registry, phase ordering,
// Results — lands when we have more than these three passes.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// CheckSource runs the Tier-0 pipeline on one source file and returns
// accumulated diagnostics. A parse error is returned as err; the caller
// decides whether to surface it as a diagnostic or as a fatal.
//
// Pass order mirrors Java's DefaultPassConfig: VarCheck first (so
// UNDEFINED_VARIABLE fires before downstream passes inspect an
// undeclared NAME), then TypeCheck.
func CheckSource(path, src string) ([]diag.Diagnostic, error) {
	script, err := parser.Parse(src)
	if err != nil {
		return nil, err
	}
	reg := jstype.NewRegistry()
	creator := NewTypedScopeCreator(reg)
	scope := creator.CreateScope(script)
	var all []diag.Diagnostic
	all = append(all, NewVarCheck(path, creator, scope).Check(script)...)
	all = append(all, NewTypeCheck(path, creator, scope).Check(script)...)
	return all, nil
}
