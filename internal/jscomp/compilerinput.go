// Port of com/google/javascript/jscomp/CompilerInput.java
// @ closure-compiler v20260415

package jscomp

import "fmt"

// ModuleType classifies the JavaScript module system used by an input file.
// Mirrors Java's CompilerInput.ModuleType enum.
type ModuleType int

const (
	// ModuleTypeNone is plain script — no module system.
	ModuleTypeNone ModuleType = iota
	// ModuleTypeGoog is a goog.module() file.
	ModuleTypeGoog
	// ModuleTypeES6 is an ES2015+ module (import/export).
	ModuleTypeES6
	// ModuleTypeCommonJS is a CommonJS module (require/exports).
	ModuleTypeCommonJS
	// ModuleTypeJSON is a JSON file.
	ModuleTypeJSON
	// ModuleTypeImportedScript is a script imported by another module.
	ModuleTypeImportedScript
)

// Require represents a single namespace dependency from a CompilerInput.
// Mirrors Java's CompilerInput DependencyInfo.Require.
type Require struct {
	// Symbol is the required namespace.
	Symbol string
	// Type classifies the require call style.
	Type RequireType
}

// RequireType classifies how a namespace is required.
// Mirrors Java's DependencyInfo.Require.Type enum.
type RequireType int

const (
	// RequireTypeGoogRequireSymbol — goog.require('sym')
	RequireTypeGoogRequireSymbol RequireType = iota
	// RequireTypeEs6Import — ES import statement
	RequireTypeEs6Import
	// RequireTypeCommonJS — CommonJS require()
	RequireTypeCommonJS
	// RequireTypeBase — implicit goog base dependency
	RequireTypeBase
)

// GoogRequireSymbol creates a Require for a goog.require() call.
// Mirrors Java's Require.googRequireSymbol(String).
func GoogRequireSymbol(symbol string) Require {
	return Require{Symbol: symbol, Type: RequireTypeGoogRequireSymbol}
}

// Es6Import creates a Require for an ES import statement.
// Mirrors Java's Require.es6Import(String, String).
func Es6Import(moduleName, importedName string) Require {
	return Require{Symbol: moduleName, Type: RequireTypeEs6Import}
}

// RequireBase is the implicit base dependency added when goog.* is referenced.
// Mirrors Java's Require.BASE.
var RequireBase = Require{Symbol: "goog", Type: RequireTypeBase}

// DependencyInfo is the minimal interface for input dependency metadata.
// Mirrors the relevant parts of Java's DependencyInfo interface.
type DependencyInfo interface {
	// GetName returns the canonical name of this input.
	GetName() string
	// GetProvides returns the namespaces this input provides.
	GetProvides() []string
	// GetRequires returns the ordered list of required namespaces.
	GetRequires() []Require
	// GetTypeRequires returns type-only required namespaces.
	GetTypeRequires() []string
	// IsEs6Module reports whether this input is an ES module.
	IsEs6Module() bool
	// IsGoogModule reports whether this input is a goog.module.
	IsGoogModule() bool
}

// CompilerInput is the internal representation of a single input file fed to
// the compiler. It holds the source path, raw content, dependency metadata
// (provides / requires), and a parsed AST root once parsing has run.
//
// Mirrors Java's CompilerInput class. The full Java version delegates lazy AST
// construction to a JsAst inner class and resolves module paths via
// ModuleLoader; those pieces are not yet wired in Go but the public surface
// is preserved for structural parity.
type CompilerInput struct {
	// id is the stable name/identifier for this input (file path).
	id string
	// sourceCode is the raw JS source text.
	sourceCode string
	// isExtern reports whether this input is an externs file.
	isExtern bool
	// jsModuleType classifies the module system used.
	jsModuleType ModuleType
	// provides is the set of namespaces declared by this input (extraProvides).
	provides []string
	// requires is the ordered list of required namespaces (orderedRequires).
	requires []Require
	// extraRequires holds additional requires added by addRequire.
	extraRequires []Require
	// typeRequires is the list of type-only required namespaces.
	typeRequires []string
	// dynamicRequires holds namespaces loaded via goog.requireDynamic.
	dynamicRequires []string
	// requireDynamicImports holds namespaces loaded via goog.requireDynamic as imports.
	requireDynamicImports []string
	// hasFullParseDependencyInfo is set once full AST-based deps have been computed.
	hasFullParseDependencyInfo bool

	// compiler is the AbstractCompiler used for lazy parsing.
	compiler *AbstractCompiler

	// hasExternsAnnotation records whether @externs was found in the JSDoc.
	hasExternsAnnotation bool
	// hasNoCompileAnnotation records whether @nocompile was found in the JSDoc.
	hasNoCompileAnnotation bool
}

// NewCompilerInput creates a CompilerInput for a source file.
// Mirrors Java's CompilerInput(SourceFile) constructor.
func NewCompilerInput(id, sourceCode string) *CompilerInput {
	return &CompilerInput{id: id, sourceCode: sourceCode}
}

// NewExternCompilerInput creates a CompilerInput for an externs file.
// Mirrors Java's CompilerInput(SourceFile, boolean) constructor with isExtern=true.
func NewExternCompilerInput(id, sourceCode string) *CompilerInput {
	return &CompilerInput{id: id, sourceCode: sourceCode, isExtern: true}
}

// ---------------------------------------------------------------------------
// Identity / naming
// ---------------------------------------------------------------------------

// GetInputID returns the stable name/identifier string for this input.
// Mirrors Java's CompilerInput.getInputId().
func (ci *CompilerInput) GetInputID() string { return ci.id }

// GetName returns the input name. Mirrors Java's CompilerInput.getName().
func (ci *CompilerInput) GetName() string { return ci.id }

// GetPathRelativeToClosureBase returns the path relative to the closure base.
// Mirrors Java's CompilerInput.getPathRelativeToClosureBase(); panics as not
// yet implemented (matching Java's UnsupportedOperationException).
func (ci *CompilerInput) GetPathRelativeToClosureBase() string {
	panic("getPathRelativeToClosureBase not implemented")
}

// ---------------------------------------------------------------------------
// Source content
// ---------------------------------------------------------------------------

// GetSourceCode returns the raw JS source. Mirrors Java's SourceFile.getCode().
func (ci *CompilerInput) GetSourceCode() string { return ci.sourceCode }

// GetCode returns the raw JS source. Mirrors Java's CompilerInput.getCode().
func (ci *CompilerInput) GetCode() string { return ci.sourceCode }

// ---------------------------------------------------------------------------
// Extern flag
// ---------------------------------------------------------------------------

// IsExtern reports whether this is an externs file.
// Mirrors Java's CompilerInput.isExtern().
func (ci *CompilerInput) IsExtern() bool { return ci.isExtern }

// SetIsExtern marks this input as an externs file.
// Mirrors Java's CompilerInput.setIsExtern().
func (ci *CompilerInput) SetIsExtern() { ci.isExtern = true }

// ---------------------------------------------------------------------------
// Module type
// ---------------------------------------------------------------------------

// GetJsModuleType returns the detected module type.
// Mirrors Java's CompilerInput.getJsModuleType().
func (ci *CompilerInput) GetJsModuleType() ModuleType { return ci.jsModuleType }

// SetJsModuleType sets the module type.
// Mirrors Java's CompilerInput.setJsModuleType().
func (ci *CompilerInput) SetJsModuleType(t ModuleType) { ci.jsModuleType = t }

// IsEs6Module reports whether this input uses ES module syntax.
// Mirrors Java's CompilerInput.isEs6Module().
func (ci *CompilerInput) IsEs6Module() bool { return ci.jsModuleType == ModuleTypeES6 }

// IsGoogModule reports whether this input is a goog.module file.
// Mirrors Java's CompilerInput.isGoogModule().
func (ci *CompilerInput) IsGoogModule() bool { return ci.jsModuleType == ModuleTypeGoog }

// ---------------------------------------------------------------------------
// Provides
// ---------------------------------------------------------------------------

// AddProvide registers a namespace that this input provides (extraProvides).
// Mirrors Java's CompilerInput.addProvide().
func (ci *CompilerInput) AddProvide(ns string) { ci.provides = append(ci.provides, ns) }

// GetProvides returns the list of provided namespaces.
// Mirrors Java's CompilerInput.getProvides().
func (ci *CompilerInput) GetProvides() []string { return ci.provides }

// GetKnownProvides returns the known provides without triggering a full parse.
// Mirrors Java's CompilerInput.getKnownProvides().
func (ci *CompilerInput) GetKnownProvides() []string { return ci.provides }

// ---------------------------------------------------------------------------
// Requires
// ---------------------------------------------------------------------------

// AddOrderedRequire registers a namespace that this input requires in order.
// Deduplicates on insertion.
// Mirrors Java's CompilerInput.addOrderedRequire(Require).
func (ci *CompilerInput) AddOrderedRequire(req Require) bool {
	for _, r := range ci.requires {
		if r == req {
			return false
		}
	}
	ci.requires = append(ci.requires, req)
	return true
}

// AddOrderedRequireBySymbol registers a symbol string as a goog.require dep.
// Convenience wrapper around AddOrderedRequire.
func (ci *CompilerInput) AddOrderedRequireBySymbol(ns string) {
	ci.AddOrderedRequire(GoogRequireSymbol(ns))
}

// AddRequire registers an additional require (extraRequires).
// Mirrors Java's CompilerInput.addRequire(Require).
func (ci *CompilerInput) AddRequire(req Require) {
	ci.extraRequires = append(ci.extraRequires, req)
}

// GetRequires returns the list of required namespaces in order.
// Mirrors Java's CompilerInput.getRequires().
func (ci *CompilerInput) GetRequires() []Require {
	if ci.hasFullParseDependencyInfo {
		return ci.requires
	}
	// Merge orderedRequires + extraRequires for a combined view.
	combined := make([]Require, 0, len(ci.requires)+len(ci.extraRequires))
	combined = append(combined, ci.requires...)
	combined = append(combined, ci.extraRequires...)
	return combined
}

// GetKnownRequires returns the known requires without triggering a full parse.
// Mirrors Java's CompilerInput.getKnownRequires().
func (ci *CompilerInput) GetKnownRequires() []Require {
	combined := make([]Require, 0, len(ci.requires)+len(ci.extraRequires))
	combined = append(combined, ci.requires...)
	combined = append(combined, ci.extraRequires...)
	return combined
}

// GetKnownRequiredSymbols returns the symbols of all known requires.
// Mirrors Java's CompilerInput.getKnownRequiredSymbols().
func (ci *CompilerInput) GetKnownRequiredSymbols() []string {
	reqs := ci.GetKnownRequires()
	symbols := make([]string, 0, len(reqs))
	for _, r := range reqs {
		symbols = append(symbols, r.Symbol)
	}
	return symbols
}

// ---------------------------------------------------------------------------
// Type requires
// ---------------------------------------------------------------------------

// AddTypeRequire registers a type-only namespace dependency.
// Mirrors Java's CompilerInput.addTypeRequire (via DependencyInfo).
func (ci *CompilerInput) AddTypeRequire(ns string) {
	ci.typeRequires = append(ci.typeRequires, ns)
}

// GetTypeRequires returns the type-only required namespaces.
// Mirrors Java's CompilerInput.getTypeRequires().
func (ci *CompilerInput) GetTypeRequires() []string { return ci.typeRequires }

// ---------------------------------------------------------------------------
// Dynamic requires
// ---------------------------------------------------------------------------

// AddDynamicRequire registers a dynamically-required namespace.
// Mirrors Java's CompilerInput.addDynamicRequire(String).
func (ci *CompilerInput) AddDynamicRequire(ns string) bool {
	for _, r := range ci.dynamicRequires {
		if r == ns {
			return false
		}
	}
	ci.dynamicRequires = append(ci.dynamicRequires, ns)
	return true
}

// GetDynamicRequires returns the list of dynamically required namespaces.
// Mirrors Java's CompilerInput.getDynamicRequires().
func (ci *CompilerInput) GetDynamicRequires() []string { return ci.dynamicRequires }

// AddRequireDynamicImports registers a goog.requireDynamic() import.
// Mirrors Java's CompilerInput.addRequireDynamicImports(String).
func (ci *CompilerInput) AddRequireDynamicImports(ns string) {
	for _, r := range ci.requireDynamicImports {
		if r == ns {
			return
		}
	}
	ci.requireDynamicImports = append(ci.requireDynamicImports, ns)
}

// GetRequireDynamicImports returns the list of goog.requireDynamic() imports.
// Mirrors Java's CompilerInput.getRequireDynamicImports().
func (ci *CompilerInput) GetRequireDynamicImports() []string { return ci.requireDynamicImports }

// ---------------------------------------------------------------------------
// Full-parse dependency info
// ---------------------------------------------------------------------------

// SetHasFullParseDependencyInfo marks that full AST-based dep info is available.
// Mirrors Java's CompilerInput.setHasFullParseDependencyInfo(boolean).
func (ci *CompilerInput) SetHasFullParseDependencyInfo(v bool) {
	ci.hasFullParseDependencyInfo = v
}

// HasFullParseDependencyInfo reports whether AST-based dep info has been computed.
// Mirrors Java's CompilerInput.hasFullParseDependencyInfo().
func (ci *CompilerInput) HasFullParseDependencyInfo() bool {
	return ci.hasFullParseDependencyInfo
}

// ---------------------------------------------------------------------------
// Compiler reference for lazy parsing
// ---------------------------------------------------------------------------

// SetCompiler stores the compiler reference used for lazy parsing.
// Mirrors Java's CompilerInput.setCompiler(AbstractCompiler).
func (ci *CompilerInput) SetCompiler(compiler *AbstractCompiler) {
	ci.compiler = compiler
}

// GetCompiler returns the compiler reference. Convenience accessor.
func (ci *CompilerInput) GetCompiler() *AbstractCompiler { return ci.compiler }

// ---------------------------------------------------------------------------
// AST access
// ---------------------------------------------------------------------------

// ClearAst releases the cached AST. Mirrors Java's CompilerInput.clearAst().
// In this port the AST is not held by CompilerInput, so this is a no-op.
func (ci *CompilerInput) ClearAst() {}

// ---------------------------------------------------------------------------
// Load flags (module type map mirroring DependencyInfo.getLoadFlags)
// ---------------------------------------------------------------------------

// GetLoadFlags returns the load flags map for this input.
// Mirrors Java's DependencyInfo.getLoadFlags().
func (ci *CompilerInput) GetLoadFlags() map[string]string {
	flags := make(map[string]string)
	switch ci.jsModuleType {
	case ModuleTypeGoog:
		flags["module"] = "goog"
	case ModuleTypeES6:
		flags["module"] = "es6"
	case ModuleTypeCommonJS:
		flags["module"] = "commonjs"
	}
	return flags
}

// ---------------------------------------------------------------------------
// Externs / nocompile annotations
// ---------------------------------------------------------------------------

// GetHasExternsAnnotation reports whether this input has @externs in its JSDoc.
// Mirrors Java's DependencyInfo.getHasExternsAnnotation().
func (ci *CompilerInput) GetHasExternsAnnotation() bool { return ci.hasExternsAnnotation }

// SetHasExternsAnnotation sets the @externs annotation flag.
func (ci *CompilerInput) SetHasExternsAnnotation(v bool) { ci.hasExternsAnnotation = v }

// GetHasNoCompileAnnotation reports whether this input has @nocompile.
// Mirrors Java's DependencyInfo.getHasNoCompileAnnotation().
func (ci *CompilerInput) GetHasNoCompileAnnotation() bool { return ci.hasNoCompileAnnotation }

// SetHasNoCompileAnnotation sets the @nocompile annotation flag.
func (ci *CompilerInput) SetHasNoCompileAnnotation(v bool) { ci.hasNoCompileAnnotation = v }

// ---------------------------------------------------------------------------
// Line offset for error reporting
// ---------------------------------------------------------------------------

// GetLineOffset returns the character offset of the given (1-indexed) line.
// Mirrors Java's CompilerInput.getLineOffset(int).
// This is a simplified implementation that counts newlines in the source.
func (ci *CompilerInput) GetLineOffset(lineno int) int {
	if lineno <= 1 {
		return 0
	}
	offset := 0
	lines := 1
	for i, ch := range ci.sourceCode {
		if ch == '\n' {
			lines++
			if lines == lineno {
				return i + 1
			}
		}
		_ = offset
	}
	return len(ci.sourceCode)
}

// ---------------------------------------------------------------------------
// String representation
// ---------------------------------------------------------------------------

func (ci *CompilerInput) String() string { return ci.id }

// ---------------------------------------------------------------------------
// Helpers / validation
// ---------------------------------------------------------------------------

// Validate checks that the CompilerInput is in a consistent state.
// Returns an error if validation fails.
func (ci *CompilerInput) Validate() error {
	if ci.id == "" {
		return fmt.Errorf("CompilerInput has empty id")
	}
	return nil
}
