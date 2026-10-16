// Port of com/google/javascript/jscomp/PassFactory.java
// @ closure-compiler v20260415

package jscomp

// CompilerPass is a single phase of the compiler pipeline.
// Mirrors Java's CompilerPass interface.
type CompilerPass interface {
	// Process runs the pass on the given externs and root nodes.
	Process(externs, root interface{})
}

// CompilerPassFunc is a function adapter for CompilerPass.
type CompilerPassFunc func(externs, root interface{})

func (f CompilerPassFunc) Process(externs, root interface{}) { f(externs, root) }

// PreconditionResult captures the outcome of a PassFactory precondition check.
// Mirrors Java's PassFactory.PreconditionResult record.
type PreconditionResult struct {
	Success bool
	Message string
}

// PreconditionResultSuccess is the success sentinel. Mirrors Java's
// PassFactory.PreconditionResult.SUCCESS.
var PreconditionResultSuccess = PreconditionResult{Success: true}

// PassFactory is a factory for creating JSCompiler passes based on the
// CompilerOptions injected. Mirrors Java's PassFactory @AutoValue class.
type PassFactory struct {
	// Name is the pass name as it will appear in logs.
	Name string
	// RunInFixedPointLoop indicates whether this factory must appear in a
	// PhaseOptimizer loop.
	RunInFixedPointLoop bool
	// Condition is a function that decides (based on options) whether the pass
	// should be added to the compilation. Defaults to always-true.
	Condition func(*CompilerOptions) bool
	// PreconditionCheck validates preconditions before the pass runs.
	PreconditionCheck func(*CompilerOptions) PreconditionResult
	// InternalFactory creates the actual CompilerPass instance.
	InternalFactory func(compiler *AbstractCompiler) CompilerPass
}

// PassFactoryBuilder supports the builder pattern for PassFactory,
// mirroring Java's PassFactory.Builder inner class.
type PassFactoryBuilder struct {
	f PassFactory
}

// NewPassFactoryBuilder returns a builder with sensible defaults.
// Mirrors Java's PassFactory.builder().
func NewPassFactoryBuilder() *PassFactoryBuilder {
	return &PassFactoryBuilder{f: PassFactory{
		Condition:         func(*CompilerOptions) bool { return true },
		PreconditionCheck: func(*CompilerOptions) PreconditionResult { return PreconditionResultSuccess },
	}}
}

func (b *PassFactoryBuilder) SetName(name string) *PassFactoryBuilder {
	b.f.Name = name
	return b
}

func (b *PassFactoryBuilder) SetRunInFixedPointLoop(v bool) *PassFactoryBuilder {
	b.f.RunInFixedPointLoop = v
	return b
}

func (b *PassFactoryBuilder) SetCondition(cond func(*CompilerOptions) bool) *PassFactoryBuilder {
	b.f.Condition = cond
	return b
}

func (b *PassFactoryBuilder) SetPreconditionCheck(check func(*CompilerOptions) PreconditionResult) *PassFactoryBuilder {
	b.f.PreconditionCheck = check
	return b
}

func (b *PassFactoryBuilder) SetInternalFactory(factory func(*AbstractCompiler) CompilerPass) *PassFactoryBuilder {
	b.f.InternalFactory = factory
	return b
}

// Build validates and returns the PassFactory.
// Mirrors Java's PassFactory.Builder.build().
func (b *PassFactoryBuilder) Build() *PassFactory {
	if b.f.Name == "" {
		panic("PassFactory name must not be empty")
	}
	result := b.f
	return &result
}

// CreateEmptyPass creates a no-op pass with the given name.
// Mirrors Java's PassFactory.createEmptyPass(String).
func CreateEmptyPass(name string) *PassFactory {
	return NewPassFactoryBuilder().
		SetName(name).
		SetInternalFactory(func(c *AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(externs, root interface{}) {})
		}).
		Build()
}

// ValidatePreconditions runs the PreconditionCheck and panics if it fails.
// Mirrors Java's PassFactory.validatePreconditions.
func (pf *PassFactory) ValidatePreconditions(options *CompilerOptions) {
	result := pf.PreconditionCheck(options)
	if !result.Success {
		panic("Precondition for pass " + pf.Name + " failed: " + result.Message)
	}
}

// Create instantiates the compiler pass. Mirrors Java's PassFactory.create.
func (pf *PassFactory) Create(compiler *AbstractCompiler) CompilerPass {
	return pf.InternalFactory(compiler)
}
