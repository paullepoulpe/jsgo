// Port of com/google/javascript/jscomp/CompilerInput.java
// @ closure-compiler v20260415

package jscomp

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
	// provides is the set of namespaces declared by this input.
	provides []string
	// requires is the ordered list of required namespaces.
	requires []string
	// typeRequires is the list of type-only required namespaces.
	typeRequires []string
	// hasFullParseDependencyInfo is set once full AST-based deps have been computed.
	hasFullParseDependencyInfo bool

	// compiler is the AbstractCompiler used for lazy parsing.
	compiler *AbstractCompiler
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

// GetInputID returns the stable name/identifier string for this input.
// Mirrors Java's CompilerInput.getInputId().
func (ci *CompilerInput) GetInputID() string { return ci.id }

// GetName returns the input name. Mirrors Java's CompilerInput.getName().
func (ci *CompilerInput) GetName() string { return ci.id }

// GetSourceCode returns the raw JS source. Mirrors Java's SourceFile.getCode().
func (ci *CompilerInput) GetSourceCode() string { return ci.sourceCode }

// IsExtern reports whether this is an externs file.
// Mirrors Java's CompilerInput.isExtern().
func (ci *CompilerInput) IsExtern() bool { return ci.isExtern }

// SetIsExtern marks this input as an externs file.
// Mirrors Java's CompilerInput.setIsExtern().
func (ci *CompilerInput) SetIsExtern() { ci.isExtern = true }

// GetJsModuleType returns the detected module type.
// Mirrors Java's CompilerInput.getJsModuleType().
func (ci *CompilerInput) GetJsModuleType() ModuleType { return ci.jsModuleType }

// SetJsModuleType sets the module type.
// Mirrors Java's CompilerInput.setJsModuleType().
func (ci *CompilerInput) SetJsModuleType(t ModuleType) { ci.jsModuleType = t }

// AddProvide registers a namespace that this input provides.
// Mirrors Java's CompilerInput.addProvide().
func (ci *CompilerInput) AddProvide(ns string) { ci.provides = append(ci.provides, ns) }

// GetProvides returns the list of provided namespaces.
// Mirrors Java's CompilerInput.getProvides().
func (ci *CompilerInput) GetProvides() []string { return ci.provides }

// AddOrderedRequire registers a namespace that this input requires.
// Mirrors Java's CompilerInput.addOrderedRequire().
func (ci *CompilerInput) AddOrderedRequire(ns string) {
	for _, r := range ci.requires {
		if r == ns {
			return
		}
	}
	ci.requires = append(ci.requires, ns)
}

// GetRequires returns the list of required namespaces in order.
// Mirrors Java's CompilerInput.getRequires().
func (ci *CompilerInput) GetRequires() []string { return ci.requires }

// AddTypeRequire registers a type-only namespace dependency.
func (ci *CompilerInput) AddTypeRequire(ns string) {
	ci.typeRequires = append(ci.typeRequires, ns)
}

// GetTypeRequires returns the type-only required namespaces.
// Mirrors Java's CompilerInput.getTypeRequires().
func (ci *CompilerInput) GetTypeRequires() []string { return ci.typeRequires }

// SetHasFullParseDependencyInfo marks that full AST-based dep info is available.
// Mirrors Java's CompilerInput.setHasFullParseDependencyInfo().
func (ci *CompilerInput) SetHasFullParseDependencyInfo(v bool) {
	ci.hasFullParseDependencyInfo = v
}

// HasFullParseDependencyInfo reports whether AST-based dep info has been computed.
func (ci *CompilerInput) HasFullParseDependencyInfo() bool {
	return ci.hasFullParseDependencyInfo
}

// SetCompiler stores the compiler reference used for lazy parsing.
// Mirrors Java's CompilerInput.setCompiler().
func (ci *CompilerInput) SetCompiler(compiler *AbstractCompiler) {
	ci.compiler = compiler
}

// IsEs6Module reports whether this input uses ES module syntax.
// Mirrors Java's CompilerInput.isEs6Module().
func (ci *CompilerInput) IsEs6Module() bool { return ci.jsModuleType == ModuleTypeES6 }

// IsGoogModule reports whether this input is a goog.module file.
// Mirrors Java's CompilerInput.isGoogModule().
func (ci *CompilerInput) IsGoogModule() bool { return ci.jsModuleType == ModuleTypeGoog }

func (ci *CompilerInput) String() string { return ci.id }
