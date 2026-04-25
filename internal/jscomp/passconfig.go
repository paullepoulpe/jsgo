// Port of com/google/javascript/jscomp/PassConfig.java
// @ closure-compiler v20260415

package jscomp

// PassConfig declares the abstract API for pass configuration.
// Concrete implementations return lists of PassFactory instances for each
// phase. Mirrors Java's PassConfig abstract class.
type PassConfig interface {
	// GetWhitespaceOnlyPasses returns checking passes that run even in
	// whitespace-only mode. Mirrors Java's getWhitespaceOnlyPasses.
	GetWhitespaceOnlyPasses() *PassListBuilder

	// GetTranspileOnlyPasses returns transpilation-only passes.
	// Mirrors Java's getTranspileOnlyPasses.
	GetTranspileOnlyPasses() *PassListBuilder

	// GetChecks returns the checking passes to run.
	// Mirrors Java's getChecks.
	GetChecks() *PassListBuilder

	// GetOptimizations returns the optimization passes to run.
	// Mirrors Java's getOptimizations.
	GetOptimizations() *PassListBuilder

	// GetFinalizations returns the finalization passes to run.
	// Mirrors Java's getFinalizations.
	GetFinalizations() *PassListBuilder

	// Options returns the CompilerOptions this config was constructed with.
	Options() *CompilerOptions
}

// BasePassConfig provides default no-op implementations of the whitespace-only
// and transpile-only passes that return empty builders. Concrete sub-configs
// embed BasePassConfig and override only GetChecks / GetOptimizations /
// GetFinalizations. Mirrors Java's PassConfig partial implementations.
type BasePassConfig struct {
	options *CompilerOptions
}

// NewBasePassConfig constructs a BasePassConfig for the given options.
func NewBasePassConfig(options *CompilerOptions) *BasePassConfig {
	return &BasePassConfig{options: options}
}

func (c *BasePassConfig) Options() *CompilerOptions {
	return c.options
}

func (c *BasePassConfig) GetWhitespaceOnlyPasses() *PassListBuilder {
	return NewPassListBuilder(c.options)
}

func (c *BasePassConfig) GetTranspileOnlyPasses() *PassListBuilder {
	return NewPassListBuilder(c.options)
}

// GetChecks is abstract — subclasses must override.
func (c *BasePassConfig) GetChecks() *PassListBuilder {
	return NewPassListBuilder(c.options)
}

// GetOptimizations is abstract — subclasses must override.
func (c *BasePassConfig) GetOptimizations() *PassListBuilder {
	return NewPassListBuilder(c.options)
}

// GetFinalizations is abstract — subclasses must override.
func (c *BasePassConfig) GetFinalizations() *PassListBuilder {
	return NewPassListBuilder(c.options)
}

// GetBasePassConfig walks a chain of PassConfigDelegate wrappers to find the
// innermost (non-delegate) PassConfig. Mirrors Java's
// PassConfig.getBasePassConfig.
func GetBasePassConfig(cfg PassConfig) PassConfig {
	for {
		if delegate, ok := cfg.(*PassConfigDelegate); ok {
			cfg = delegate.delegate
		} else {
			return cfg
		}
	}
}

// PassConfigDelegate proxies all PassConfig calls to an inner delegate.
// Mirrors Java's PassConfig.PassConfigDelegate inner class.
type PassConfigDelegate struct {
	delegate PassConfig
}

// NewPassConfigDelegate wraps the given PassConfig in a delegate.
func NewPassConfigDelegate(delegate PassConfig) *PassConfigDelegate {
	return &PassConfigDelegate{delegate: delegate}
}

func (d *PassConfigDelegate) Options() *CompilerOptions {
	return d.delegate.Options()
}

func (d *PassConfigDelegate) GetWhitespaceOnlyPasses() *PassListBuilder {
	return d.delegate.GetWhitespaceOnlyPasses()
}

func (d *PassConfigDelegate) GetTranspileOnlyPasses() *PassListBuilder {
	return d.delegate.GetTranspileOnlyPasses()
}

func (d *PassConfigDelegate) GetChecks() *PassListBuilder {
	return d.delegate.GetChecks()
}

func (d *PassConfigDelegate) GetOptimizations() *PassListBuilder {
	return d.delegate.GetOptimizations()
}

func (d *PassConfigDelegate) GetFinalizations() *PassListBuilder {
	return d.delegate.GetFinalizations()
}
