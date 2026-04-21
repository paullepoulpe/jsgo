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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// extractMultiDeclName scrapes the variable name out of a
// JSC_VAR_MULTIPLY_DECLARED_ERROR message. The formatter writes
// `Variable <name> declared more than once. First occurrence: …`.
func extractMultiDeclName(msg string) string {
	const prefix = "Variable "
	if !strings.HasPrefix(msg, prefix) {
		return ""
	}
	rest := msg[len(prefix):]
	if i := strings.Index(rest, " declared more than once"); i >= 0 {
		return rest[:i]
	}
	return ""
}

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
	// Post-scope wiring: fill in templatized @extends links that
	// scope.go's LookupNominalType-based path can't see.
	WireTemplatizedInheritance(script, scope)
	var all []diag.Diagnostic
	all = append(all, NewVarCheck(path, creator, scope).Check(script)...)
	varRefDiags := NewVariableReferenceCheck(path, externScope).Check(script)
	all = append(all, varRefDiags...)
	all = append(all, NewJSTypeResolver(path, src, creator, scope, externScope).Check(script)...)
	// Names that already triggered JSC_VAR_MULTIPLY_DECLARED_ERROR
	// (e.g. a script-local @interface that shadows an extern of the
	// same name) get suppressed downstream so we don't double-report
	// derived diagnostics like JSC_INTERFACE_METHOD_NOT_IMPLEMENTED —
	// Java's TypeCheck skips those when the interface slot was
	// rejected by the earlier VarCheck pass.
	multiDecl := map[string]bool{}
	for _, d := range varRefDiags {
		if d.Code == diag.JSC_VAR_MULTIPLY_DECLARED_ERROR {
			if name := extractMultiDeclName(d.Message); name != "" {
				multiDecl[name] = true
			}
		}
	}
	tc := NewTypeCheckWithOptions(path, creator, scope, opts)
	tc.SuppressedNames = multiDecl
	all = append(all, tc.Check(script)...)
	// CheckGlobalThis runs after TypeCheck, matching the Java
	// DefaultPassConfig order. It reports `this` references that
	// resolve to the global object and is structurally independent
	// of the inference pipeline.
	all = append(all, NewCheckGlobalThis(path).Check(script)...)
	all = append(all, NewMissingReturnCheck(path, scope).Check(script)...)
	all = append(all, NewCheckDuplicateMember(path).Check(script)...)
	all = append(all, NewCheckSideEffects(path).Check(script)...)
	// Java runs CheckAccessControls after TypeInference + InferJSDocInfo
	// stamp per-property visibility. Ordering here matches: TypeCheck has
	// already populated InstanceType property slots; CheckAccessControls
	// reads them + the AST's JSDoc in a single extra walk.
	all = append(all, NewCheckAccessControls(path, creator, scope).Check(script)...)
	return all, nil
}
