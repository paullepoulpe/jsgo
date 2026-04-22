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
// undeclared NAME), then TypeCheck, then MissingReturnCheck.
//
// An externs scope is built from defaultExternsSrc and used as the
// parent scope of the user source — matches Compiler#compile which
// accepts both externs and sources.
//
// Options default to NewCompilerOptions() (checkTypes=ERROR,
// strictCheckTypes=OFF) so the oracle's default-mode diagnostics
// still match. Callers wanting strict behaviour use
// CheckSourceWithOptions.
func CheckSource(path, src string) ([]diag.Diagnostic, error) {
	return CheckSourceWithOptions(path, src, NewCompilerOptions())
}

// CheckSourceWithOptions is the per-options variant. Mirrors how
// Compiler#compile accepts a CompilerOptions instance and threads
// it into each pass's constructor.
func CheckSourceWithOptions(path, src string, opts *CompilerOptions) ([]diag.Diagnostic, error) {
	if opts == nil {
		opts = NewCompilerOptions()
	}
	script, err := parser.Parse(src)
	if err != nil {
		return nil, err
	}
	reg := jstype.NewRegistry()
	creator := NewTypedScopeCreator(reg)
	externScope := creator.createExternScope()
	scope := creator.createScopeInWithParent(script, externScope)
	var all []diag.Diagnostic
	all = append(all, NewVarCheck(path, creator, scope).Check(script)...)
	all = append(all, NewVariableReferenceCheck(path, externScope).Check(script)...)
	all = append(all, NewCheckJsDocTypes(path, src, scope).Check()...)
	all = append(all, NewTypeCheckWithOptions(path, creator, scope, opts).Check(script)...)
	all = append(all, NewMissingReturnCheck(path, scope).Check(script)...)
	all = append(all, NewSuperclassCheck(path, creator, scope).Check(script)...)
	all = append(all, NewCheckDupVarDecl(path, scope).Check(script)...)
	all = append(all, NewCheckDuplicateMember(path).Check(script)...)
	all = append(all, NewCheckSideEffects(path).Check(script)...)
	all = append(all, NewCheckDebuggerStatement(path, opts).Check(script)...)
	return all, nil
}
