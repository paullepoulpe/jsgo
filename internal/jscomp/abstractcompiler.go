// Port of com/google/javascript/jscomp/AbstractCompiler.java
// @ closure-compiler v20260415

package jscomp

import "github.com/paullepoulpe/jsgo/internal/diag"

// LifeCycleStage represents what point in optimizations the compiler is at.
// Mirrors Java's AbstractCompiler.LifeCycleStage enum.
type LifeCycleStage int

const (
	// LifeCycleRaw is the initial state: no normalization, no color.
	LifeCycleRaw LifeCycleStage = iota
	// LifeCycleColorsAndSimplifiedJSDoc after ConvertTypesToColors.
	LifeCycleColorsAndSimplifiedJSDoc
	// LifeCycleNormalized after Normalize.java.
	LifeCycleNormalized
	// LifeCycleNormalizedObfuscated after variable renaming.
	LifeCycleNormalizedObfuscated
)

func (s LifeCycleStage) IsNormalized() bool {
	return s == LifeCycleNormalized || s == LifeCycleNormalizedObfuscated
}

func (s LifeCycleStage) IsNormalizedUnobfuscated() bool {
	return s == LifeCycleNormalized
}

func (s LifeCycleStage) IsNormalizedObfuscated() bool {
	return s == LifeCycleNormalizedObfuscated
}

func (s LifeCycleStage) HasColorAndSimplifiedJSDoc() bool {
	return s != LifeCycleRaw
}

// ConfigContext distinguishes parse configuration contexts.
// Mirrors Java's AbstractCompiler.ConfigContext enum.
type ConfigContext int

const (
	ConfigContextDefault ConfigContext = iota
	ConfigContextExterns
)

// AbstractCompiler is the base type for the compiler. Passes receive a
// *AbstractCompiler to access compiler state without depending on the full
// Compiler implementation. Mirrors Java's AbstractCompiler abstract class.
//
// This Go version is a struct rather than an interface so that passes can
// embed it and override only the methods they need. Methods that haven't been
// wired to a real implementation yet panic with "not implemented" rather than
// silently returning zero values.
type AbstractCompiler struct {
	options       *CompilerOptions
	lifeCycle     LifeCycleStage
	diagnostics   []diag.Diagnostic
	passName      string
}

// NewAbstractCompiler returns an AbstractCompiler with default options.
func NewAbstractCompiler(opts *CompilerOptions) *AbstractCompiler {
	if opts == nil {
		opts = NewCompilerOptions()
	}
	return &AbstractCompiler{options: opts}
}

// GetOptions returns the compiler options. Mirrors Java's AbstractCompiler.getOptions.
func (c *AbstractCompiler) GetOptions() *CompilerOptions {
	return c.options
}

// GetLifeCycleStage returns the current AST life-cycle stage.
// Mirrors Java's AbstractCompiler.getLifeCycleStage.
func (c *AbstractCompiler) GetLifeCycleStage() LifeCycleStage {
	return c.lifeCycle
}

// SetLifeCycleStage updates the current AST life-cycle stage.
// Mirrors Java's AbstractCompiler.setLifeCycleStage.
func (c *AbstractCompiler) SetLifeCycleStage(s LifeCycleStage) {
	c.lifeCycle = s
}

// BeforePass is called before each pass runs.
// Mirrors Java's AbstractCompiler.beforePass.
func (c *AbstractCompiler) BeforePass(passName string) {
	c.passName = passName
}

// AfterPass is called after each pass finishes.
// Mirrors Java's AbstractCompiler.afterPass.
func (c *AbstractCompiler) AfterPass(passName string) {}

// GetLastPassName returns the name of the most-recently started pass.
// Mirrors Java's AbstractCompiler.getLastPassName.
func (c *AbstractCompiler) GetLastPassName() string {
	return c.passName
}

// Report records a diagnostic. Mirrors Java's AbstractCompiler.report.
func (c *AbstractCompiler) Report(d diag.Diagnostic) {
	c.diagnostics = append(c.diagnostics, d)
}

// Diagnostics returns accumulated diagnostics.
func (c *AbstractCompiler) Diagnostics() []diag.Diagnostic {
	return c.diagnostics
}

// HasTypeCheckingRun reports whether type checking has run.
// Mirrors Java's AbstractCompiler.hasTypeCheckingRun.
func (c *AbstractCompiler) HasTypeCheckingRun() bool {
	return c.lifeCycle >= LifeCycleColorsAndSimplifiedJSDoc
}

// SetTypeCheckingHasRun sets whether type checking has completed.
func (c *AbstractCompiler) SetTypeCheckingHasRun(hasRun bool) {
	if hasRun && c.lifeCycle < LifeCycleColorsAndSimplifiedJSDoc {
		c.lifeCycle = LifeCycleColorsAndSimplifiedJSDoc
	}
}

// createFillFileName mirrors Java's AbstractCompiler.createFillFileName.
func createFillFileName(chunkName string) string {
	return chunkName + "$fillFile"
}

// IsFillFileName mirrors Java's AbstractCompiler.isFillFileName.
func IsFillFileName(fileName string) bool {
	const suffix = "$fillFile"
	return len(fileName) >= len(suffix) && fileName[len(fileName)-len(suffix):] == suffix
}
