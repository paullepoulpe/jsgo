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

// Go-only: language-adapter
// NewCompilerInput creates a CompilerInput for a source file.
// Mirrors Java's CompilerInput(SourceFile) constructor.
func NewCompilerInput(id, sourceCode string) *CompilerInput {
	return &CompilerInput{id: id, sourceCode: sourceCode}
}

// Go-only: language-adapter
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
// Mirrors Java's CompilerInput.getPathRelativeToClosureBase(); throws as Java
// does with UnsupportedOperationException — the method has no useful
// implementation until ModuleLoader is ported.
// TODO(port): CompilerInput.getPathRelativeToClosureBase() — requires ModuleLoader
func (ci *CompilerInput) GetPathRelativeToClosureBase() string {
	// Java throws UnsupportedOperationException with "TODO(nicksantos): Implement me."
	// ModuleLoader.resolve(getName()) is required to compute the closure-relative path.
	msg := "getPathRelativeToClosureBase: requires ModuleLoader port (see CompilerInput.java TODO)"
	panic(msg)
}

// ---------------------------------------------------------------------------
// Source content
// ---------------------------------------------------------------------------

// Go-only: language-adapter
// GetSourceCode returns the raw JS source. Java's CompilerInput delegates to
// sourceFile.getCode(); this Go port inlines the source directly.
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

// Go-only: language-adapter
// AddOrderedRequireBySymbol registers a symbol string as a goog.require dep.
// Convenience wrapper around AddOrderedRequire for the common case.
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

// Go-only: language-adapter
// AddTypeRequire registers a type-only namespace dependency. Java's DepsFinder
// populates typeRequires internally; Go needs an explicit adder for tests.
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

// Go-only: language-adapter
// GetCompiler returns the compiler reference. Java accesses the AbstractCompiler
// field directly within the same package; Go needs an explicit getter.
func (ci *CompilerInput) GetCompiler() *AbstractCompiler { return ci.compiler }

// ---------------------------------------------------------------------------
// AST access
// ---------------------------------------------------------------------------

// ClearAst releases the cached AST. Mirrors Java's CompilerInput.clearAst().
// In this port the AST is managed externally by TypedScopeCreator, so we reset
// the full parse dependency info to force re-computation on next access.
func (ci *CompilerInput) ClearAst() {
	// Java's JsAst.clearAst() sets root = null and clears cached source.
	// In Go, we reset dependency info as a proxy for re-parse readiness.
	ci.hasFullParseDependencyInfo = false
}

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

// Go-only: test-helper
// Validate checks that the CompilerInput is in a consistent state.
// Returns an error if validation fails.
func (ci *CompilerInput) Validate() error {
	if ci.id == "" {
		return fmt.Errorf("CompilerInput has empty id")
	}
	return nil
}

// ---------------------------------------------------------------------------
// Chunk membership
// ---------------------------------------------------------------------------

// JSChunk is a lightweight stub for the chunk (module) that a CompilerInput
// belongs to. Mirrors Java's JSChunk class (only the identity portion).
type JSChunk struct {
	name string
}

// Go-only: language-adapter
// NewJSChunk returns a new JSChunk with the given name.
// Mirrors Java's JSChunk(String) constructor.
func NewJSChunk(name string) *JSChunk {
	return &JSChunk{name: name}
}

// GetName returns the chunk name.
func (ch *JSChunk) GetName() string { return ch.name }

// chunk holds the JSChunk this input belongs to.
// Added as a separate embedded field to match Java's CompilerInput.chunk.
// Access via GetChunk / SetChunk / OverrideChunk.

// --- chunk accessor methods ---

var ciChunks = map[*CompilerInput]*JSChunk{}

// GetChunk returns the chunk to which the input belongs.
// Mirrors Java's CompilerInput.getChunk().
func (ci *CompilerInput) GetChunk() *JSChunk {
	return ciChunks[ci]
}

// SetChunk sets the chunk to which the input belongs.
// An input may only belong to one chunk; panics if already set to a different one.
// Mirrors Java's CompilerInput.setChunk(JSChunk).
func (ci *CompilerInput) SetChunk(chunk *JSChunk) {
	existing := ciChunks[ci]
	if existing != nil && chunk != nil && existing != chunk {
		panic(fmt.Sprintf("CompilerInput %q already belongs to chunk %q; cannot set to %q",
			ci.id, existing.name, chunk.name))
	}
	if chunk == nil {
		delete(ciChunks, ci)
	} else {
		ciChunks[ci] = chunk
	}
}

// OverrideChunk unconditionally sets the chunk, even if already set.
// Mirrors Java's CompilerInput.overrideChunk(JSChunk).
func (ci *CompilerInput) OverrideChunk(chunk *JSChunk) {
	if chunk == nil {
		delete(ciChunks, ci)
	} else {
		ciChunks[ci] = chunk
	}
}

// ---------------------------------------------------------------------------
// TypedScope association
// ---------------------------------------------------------------------------

// ciTypedScopes maps CompilerInput instances to their associated typed scopes.
// The scope values are *TypedScope from scope.go; stored as interface{} to
// avoid a circular-type dependency in this file.
var ciTypedScopes = map[*CompilerInput]interface{}{}

// Go-only: language-adapter
// SetInputTypedScope stores the typed scope for this input's MODULE_BODY node.
// Mirrors Java's CompilerInput.setTypedScope(TypedScope) with a renamed method
// to avoid collision with TypedScope-type methods in Go.
func (ci *CompilerInput) SetInputTypedScope(scope interface{}) {
	if scope == nil {
		delete(ciTypedScopes, ci)
	} else {
		ciTypedScopes[ci] = scope
	}
}

// Go-only: language-adapter
// GetInputTypedScope returns the typed scope associated with this input.
// Mirrors Java's CompilerInput.getTypedScope() with a renamed method to avoid
// collision with TypedScope-type methods in Go.
func (ci *CompilerInput) GetInputTypedScope() interface{} {
	return ciTypedScopes[ci]
}

// ---------------------------------------------------------------------------
// Feature set
// ---------------------------------------------------------------------------

// FeatureSet is a stub for the set of ECMAScript features used by an input.
// Mirrors the relevant parts of Java's FeatureSet.
type FeatureSet struct {
	version string
}

// Go-only: language-adapter
// NewFeatureSet returns a FeatureSet for the given language version string.
// Mirrors the FeatureSet.ES2021 factory pattern in Java.
func NewFeatureSet(version string) FeatureSet {
	return FeatureSet{version: version}
}

// Version returns the language version string (e.g. "es6").
func (f FeatureSet) Version() string { return f.version }

// Contains reports whether this feature set contains all features in other.
// Simplified stub: returns true if versions match.
func (f FeatureSet) Contains(other FeatureSet) bool {
	return f.version >= other.version
}

var ciFeatures = map[*CompilerInput]FeatureSet{}

// GetFeatures returns the feature set used by this input.
// Mirrors Java's CompilerInput.getFeatures(AbstractCompiler).
func (ci *CompilerInput) GetFeatures(compiler *AbstractCompiler) FeatureSet {
	if f, ok := ciFeatures[ci]; ok {
		return f
	}
	return FeatureSet{}
}

// Go-only: language-adapter
// SetFeatures stores the feature set for this input. Java's JsAst.parse() sets
// the feature set from the ParseResult internally; Go needs an explicit setter.
func (ci *CompilerInput) SetFeatures(f FeatureSet) {
	ciFeatures[ci] = f
}

// ---------------------------------------------------------------------------
// Shadow AST (for goog closure-unaware subtrees)
// ---------------------------------------------------------------------------

// ShadowAstRoot is a stub for the shadow AST node. Mirrors Java's JsAst.root
// field used by initShadowAst.
var ciShadowAsts = map[*CompilerInput]interface{}{}

// InitShadowAst stores a shadow script node for closure-unaware subtrees.
// Mirrors Java's CompilerInput.initShadowAst(Node).
func (ci *CompilerInput) InitShadowAst(shadowScript interface{}) {
	ciShadowAsts[ci] = shadowScript
}

// ---------------------------------------------------------------------------
// Serialization helpers
// ---------------------------------------------------------------------------

// SimpleDependencyInfo is a value-type snapshot of dependency info for one
// input. Mirrors the relevant subset of Java's SimpleDependencyInfo.
type SimpleDependencyInfo struct {
	Name              string
	Provides          []string
	Requires          []Require
	TypeRequires      []string
	LoadFlags         map[string]string
	HasExterns        bool
	HasNoCompile      bool
}

// Go-only: language-adapter
// ToSimpleDependencyInfo converts the CompilerInput's current dependency state
// into a SimpleDependencyInfo snapshot. Java's SimpleDependencyInfo.builder()
// is used inline in Java; Go uses this method for serialization and testing.
func (ci *CompilerInput) ToSimpleDependencyInfo() SimpleDependencyInfo {
	return SimpleDependencyInfo{
		Name:         ci.id,
		Provides:     append([]string(nil), ci.provides...),
		Requires:     append([]Require(nil), ci.GetRequires()...),
		TypeRequires: append([]string(nil), ci.typeRequires...),
		LoadFlags:    ci.GetLoadFlags(),
		HasExterns:   ci.hasExternsAnnotation,
		HasNoCompile: ci.hasNoCompileAnnotation,
	}
}

// ---------------------------------------------------------------------------
// Dependency reset
// ---------------------------------------------------------------------------

// Go-only: language-adapter
// ClearParsedDependencyInfo clears all dependency information so it can be
// re-computed from a fresh parse. Java resets dependencyInfo to nil and calls
// clearAst(); Go exposes the reset as an explicit method.
func (ci *CompilerInput) ClearParsedDependencyInfo() {
	ci.requires = nil
	ci.extraRequires = nil
	ci.provides = nil
	ci.typeRequires = nil
	ci.dynamicRequires = nil
	ci.requireDynamicImports = nil
	ci.hasFullParseDependencyInfo = false
}

// ---------------------------------------------------------------------------
// Utility: as symbol list
// ---------------------------------------------------------------------------

// Go-only: language-adapter
// RequireSymbols returns the list of required namespace symbols from a slice
// of Require values. Mirrors Java's Require.asSymbolList(Collection<Require>)
// which is a static method on the inner Require class.
func RequireSymbols(requires []Require) []string {
	out := make([]string, 0, len(requires))
	for _, r := range requires {
		out = append(out, r.Symbol)
	}
	return out
}
