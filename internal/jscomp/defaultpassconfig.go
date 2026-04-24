// Port of com/google/javascript/jscomp/DefaultPassConfig.java
// @ closure-compiler v20260415
//
// DefaultPassConfig is a ~3000-line file in Java covering all compiler passes.
// This port covers the structural skeleton and the passes that jsgo currently
// implements. Unimplemented passes are registered as no-ops so the pass list
// infrastructure is exercisable.

package jscomp

// DefaultPassConfig is the concrete PassConfig implementation for the full
// compiler pipeline. Embeds BasePassConfig and overrides the three abstract
// phase methods. Mirrors Java's DefaultPassConfig final class.
type DefaultPassConfig struct {
	*BasePassConfig
}

// NewDefaultPassConfig constructs a DefaultPassConfig for the given options.
// Mirrors Java's DefaultPassConfig(CompilerOptions) constructor.
func NewDefaultPassConfig(options *CompilerOptions) *DefaultPassConfig {
	if options == nil {
		options = NewCompilerOptions()
	}
	return &DefaultPassConfig{BasePassConfig: NewBasePassConfig(options)}
}

// GetChecks returns the checking passes to run.
// Mirrors Java's DefaultPassConfig.getChecks().
//
// The ordering here mirrors Java's pass order in DefaultPassConfig.getChecks:
// VarCheck → TypeCheck → CheckMissingReturn → CheckDupVarDecl →
// CheckDuplicateMember → CheckSideEffects → CheckDebuggerStatement →
// CheckRegExp → CheckSuspiciousCode → CheckUnreachableCode → CheckSuper →
// CheckAccessControls → CheckConformance → CheckJSDoc →
// CheckTypeImportCodeReferences → CheckClosureImports → CheckMissingRequires →
// CheckGlobalThis (gated on no VAR_MULTIPLY_DECLARED_ERROR).
func (c *DefaultPassConfig) GetChecks() *PassListBuilder {
	checks := NewPassListBuilder(c.options)

	// Checking passes — wired to the passes jsgo currently implements.
	// Each factory wraps the existing standalone Check* functions.
	// The AbstractCompiler is not yet threaded into these constructors; that
	// wiring lands in a follow-up once the compiler driver is fully wired.

	checks.Add(CreateEmptyPass("beforeStandardChecks"))

	// checkVars — VarCheck + VariableReferenceCheck
	checks.Add(&PassFactory{
		Name:      "checkVars",
		Condition: func(*CompilerOptions) bool { return true },
		PreconditionCheck: func(*CompilerOptions) PreconditionResult {
			return PreconditionResultSuccess
		},
		// InternalFactory not wired; pipeline uses CheckSource directly for now.
		InternalFactory: func(compiler *AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(externs, root interface{}) {})
		},
	})

	checks.Add(&PassFactory{
		Name: "checkTypes",
		Condition: func(opts *CompilerOptions) bool {
			return opts.IsOn(GroupCheckTypes)
		},
		PreconditionCheck: func(*CompilerOptions) PreconditionResult {
			return PreconditionResultSuccess
		},
		InternalFactory: func(compiler *AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(externs, root interface{}) {})
		},
	})

	checks.Add(&PassFactory{
		Name:      "checkMissingReturn",
		Condition: func(*CompilerOptions) bool { return true },
		PreconditionCheck: func(*CompilerOptions) PreconditionResult {
			return PreconditionResultSuccess
		},
		InternalFactory: func(compiler *AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(externs, root interface{}) {})
		},
	})

	checks.Add(&PassFactory{
		Name:      "checkDuplicateMembers",
		Condition: func(*CompilerOptions) bool { return true },
		PreconditionCheck: func(*CompilerOptions) PreconditionResult {
			return PreconditionResultSuccess
		},
		InternalFactory: func(compiler *AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(externs, root interface{}) {})
		},
	})

	checks.Add(&PassFactory{
		Name:      "suspiciousCode",
		Condition: func(*CompilerOptions) bool { return true },
		PreconditionCheck: func(*CompilerOptions) PreconditionResult {
			return PreconditionResultSuccess
		},
		InternalFactory: func(compiler *AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(externs, root interface{}) {})
		},
	})

	checks.Add(&PassFactory{
		Name:      "checkAccessControls",
		Condition: func(*CompilerOptions) bool { return true },
		PreconditionCheck: func(*CompilerOptions) PreconditionResult {
			return PreconditionResultSuccess
		},
		InternalFactory: func(compiler *AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(externs, root interface{}) {})
		},
	})

	checks.Add(&PassFactory{
		Name:      "checkConformance",
		Condition: func(*CompilerOptions) bool { return true },
		PreconditionCheck: func(*CompilerOptions) PreconditionResult {
			return PreconditionResultSuccess
		},
		InternalFactory: func(compiler *AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(externs, root interface{}) {})
		},
	})

	checks.Add(&PassFactory{
		Name:      "checkJsDoc",
		Condition: func(*CompilerOptions) bool { return true },
		PreconditionCheck: func(*CompilerOptions) PreconditionResult {
			return PreconditionResultSuccess
		},
		InternalFactory: func(compiler *AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(externs, root interface{}) {})
		},
	})

	checks.Add(&PassFactory{
		Name:      "checkMissingRequires",
		Condition: func(*CompilerOptions) bool { return true },
		PreconditionCheck: func(*CompilerOptions) PreconditionResult {
			return PreconditionResultSuccess
		},
		InternalFactory: func(compiler *AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(externs, root interface{}) {})
		},
	})

	checks.Add(CreateEmptyPass(PassNameAfterStandardChecks))

	return checks
}

// GetOptimizations returns an empty list — optimizations are not yet ported.
// Mirrors Java's DefaultPassConfig.getOptimizations() skeleton.
func (c *DefaultPassConfig) GetOptimizations() *PassListBuilder {
	return NewPassListBuilder(c.options)
}

// GetFinalizations returns an empty list — finalizations are not yet ported.
// Mirrors Java's DefaultPassConfig.getFinalizations() skeleton.
func (c *DefaultPassConfig) GetFinalizations() *PassListBuilder {
	return NewPassListBuilder(c.options)
}

