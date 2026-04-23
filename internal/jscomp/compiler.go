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

// hasVarMultiplyError reports whether the diagnostic slice contains any
// JSC_VAR_MULTIPLY_DECLARED_ERROR entry. CheckGlobalThis is skipped when
// such errors are present, mirroring Java's behaviour where the
// suspiciousCode phase (which hosts CheckGlobalThis) does not run for
// files that fail declaration-uniqueness checks.
func hasVarMultiplyError(diags []diag.Diagnostic) bool {
	for _, d := range diags {
		if d.Code == diag.JSC_VAR_MULTIPLY_DECLARED_ERROR {
			return true
		}
	}
	return false
}

// CheckSource runs the Tier-0 pipeline on one source file and returns
// accumulated diagnostics. A parse error is returned as err; the caller
// decides whether to surface it as a diagnostic or as a fatal.
//
// Pass order mirrors Java's DefaultPassConfig: VarCheck first (so
// UNDEFINED_VARIABLE fires before downstream passes inspect an
// undeclared NAME), then TypeCheck, then CheckMissingReturn.
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
	// InferConsts annotates NAME nodes before diagnostic passes run.
	NewInferConsts().Check(script)
	var all []diag.Diagnostic
	all = append(all, NewVarCheck(path, creator, scope).Check(script)...)
	all = append(all, NewVariableReferenceCheck(path, externScope).Check(script)...)
	all = append(all, NewCheckJsDocTypes(path, src, scope).Check()...)
	all = append(all, NewTypeCheckWithOptions(path, creator, scope, opts).Check(script)...)
	all = append(all, NewCheckMissingReturn(path, scope).Check(script)...)
	all = append(all, NewSuperclassCheck(path, creator, scope).Check(script)...)
	all = append(all, NewCheckDupVarDecl(path, scope).Check(script)...)
	all = append(all, NewCheckDuplicateMember(path).Check(script)...)
	all = append(all, NewCheckSideEffects(path).Check(script)...)
	all = append(all, NewCheckDebuggerStatement(path, opts).Check(script)...)
	all = append(all, NewCheckRegExp(path, true, true).Check(script)...)
	all = append(all, NewCheckSuspiciousCode(path).Check(script)...)
	all = append(all, NewCheckUnreachableCode(path).Check(script)...)
	all = append(all, NewCheckSuper(path).Check(script)...)
	all = append(all, NewCheckConformance(path).Check(script)...)
	all = append(all, NewCheckJSDoc(path).Check(script)...)
	// CheckClosureImports is registered here but runs with a stub
	// ModuleMetadataMap (no namespace metadata). Full enforcement fires
	// only after ModuleMetadataMap + ModuleLoader are ported. Until then
	// the pass produces no diagnostics so it does not disturb the corpus.
	all = append(all, NewCheckClosureImports(path, nil).Check(script)...)
	// CheckMissingRequires requires a ModuleMetadataMap (module system not yet
	// ported). Pass nil so the pass no-ops until the module system lands.
	all = append(all, NewCheckMissingRequires(path, nil).Check(script)...)
	// CheckGlobalThis is skipped when VarCheck / VariableReferenceCheck
	// already detected a JSC_VAR_MULTIPLY_DECLARED_ERROR. Java's compiler
	// aborts the suspiciousCode phase (which hosts CheckGlobalThis) when
	// declaration errors are present, so the oracle never shows both errors
	// for the same file.
	if !hasVarMultiplyError(all) {
		all = append(all, NewCheckGlobalThis(path).Check(script)...)
	}
	return all, nil
}
