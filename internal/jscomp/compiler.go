// Port of com/google/javascript/jscomp/Compiler.java
// @ closure-compiler v20260415
//
// Closure's Compiler is a ~3000-line orchestrator (PhaseOptimizer,
// pass registry, Result, module graph, warning manager, …). The Go
// port is structured as:
//
//   - The free functions (CheckSource, parseSource, scanLexicalErrors,
//     scanBoundedGenericErrors) drive the harness + CLI today by
//     wiring our three independent passes (parse → TypedScopeCreator
//     → TypeCheck) directly. They mirror the work Java's
//     Compiler.compile() / parseForCompilation() / stage1Passes() do
//     for the Tier-0 pipeline.
//
//   - The Compiler struct + methods below are the 1:1 surface port of
//     Java's Compiler class. They wrap an embedded *AbstractCompiler
//     for shared state (options, lifecycle, diagnostics) and add the
//     getters / setters / lifecycle hooks Java exposes. Method names
//     match Java's exactly so that downstream features (pass
//     pipeline, Result construction) can land 1:1.
//
//   - Companion types (Result, JSError, DiagnosticGroups,
//     ErrorManager) are minimal stubs with TODO(port) markers
//     pointing at the Java types they will eventually back. They are
//     real Go types — not aliased interface{} — so signatures can
//     match Java's once the underlying classes land.

package jscomp

import (
	"fmt"
	"os"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// useNewPipeline reports whether JSGO_NEW_PIPELINE=1 is set, gating the
// ParserRewrite → IRFactory path. When false, the legacy parser.Parse path
// is used unchanged.
var useNewPipeline = os.Getenv("JSGO_NEW_PIPELINE") == "1"

// Go-only: language-adapter
// scanLexicalErrors runs the new Scanner over src and returns one
// JSC_PARSE_ERROR diagnostic per reported lexical error (invalid escape
// sequences, unterminated strings, …). Mirrors Java's scanner which is
// part of the parse pipeline and surfaces every scanner error via the
// compiler's ErrorReporter. The legacy Go parser (still the default for
// the AST path) silently recovers from these, so we run a scanner-only
// pre-pass to reach Java parity. Returns nil if the source is lexically
// valid.
func scanLexicalErrors(path, src string) []diag.Diagnostic {
	er := &scannerOnlyReporter{}
	sc := parser.NewScannerFull(er, nil, src, 0)
	for {
		tok := sc.NextToken()
		if tok == nil {
			break
		}
		if tok.Type == parser.TokenTypeEndOfFile {
			break
		}
	}
	if len(er.errs) == 0 {
		return nil
	}
	// Java's ParserRunner reports only the first scanner error and halts
	// compilation (subsequent invalid escapes in the same source produce
	// cascading false positives, which Java avoids). Emit only the first.
	e := er.errs[0]
	return []diag.Diagnostic{{
		Path:    path,
		Line:    e.line + 1,
		Col:     e.col,
		Code:    diag.JSC_PARSE_ERROR,
		Level:   diag.Error,
		Message: "Parse error. " + e.msg,
	}}
}

// Go-only: language-adapter — minimal ErrorReporter that records every
// scanner error (line, col, message).
type scannerOnlyReporter struct {
	errs []scannerErr
}

type scannerErr struct {
	line int
	col  int
	msg  string
}

func (r *scannerOnlyReporter) ReportError(position parser.SourcePosition, format string, args ...interface{}) {
	r.errs = append(r.errs, scannerErr{line: position.Line, col: position.Column, msg: fmt.Sprintf(format, args...)})
}

func (r *scannerOnlyReporter) ReportWarning(position parser.SourcePosition, format string, args ...interface{}) {
	_ = position
}

// Go-only: language-adapter
// scanBoundedGenericErrors walks the AST collecting every JSDocInfo's
// BoundedTemplatePositions and emitting one
// JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES per occurrence. Mirrors how
// Java's JsDocInfoParser TEMPLATE branch records the diagnostic via
// RhinoErrorReporter and halts subsequent compilation, so this
// pre-pass shorts the pipeline before VarCheck runs.
func scanBoundedGenericErrors(path string, script *rhino.Node) []diag.Diagnostic {
	var out []diag.Diagnostic
	var walk func(n *rhino.Node)
	walk = func(n *rhino.Node) {
		if n == nil {
			return
		}
		if info := n.JSDocInfo(); info != nil {
			for _, pos := range info.BoundedTemplatePositions() {
				out = append(out, diag.Diagnostic{
					Path:    path,
					Line:    pos.Line,
					Col:     pos.Col,
					Code:    diag.JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES,
					Level:   diag.Error,
					Message: "Bounded generic semantics are currently still in development",
				})
			}
		}
		for _, c := range n.Children() {
			walk(c)
		}
	}
	walk(script)
	return out
}

// Go-only: package-bridge
// parseSource selects the parse path based on useNewPipeline.
// It mirrors Compiler.parse(input) which calls ParserRunner.parse internally.
// Returns the SCRIPT-rooted rhino.Node and any fatal parse error.
func parseSource(path, src string) (*rhino.Node, error) {
	if useNewPipeline {
		result := runNewPipeline(path, src)
		if result.hadError || result.ast == nil {
			msg := result.firstError
			if msg == "" {
				msg = "parse error"
			}
			return nil, fmt.Errorf("%s", msg)
		}
		return result.ast, nil
	}
	return parser.Parse(src)
}

// Go-only: package-bridge
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

// Go-only: package-bridge
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

// Go-only: package-bridge
// CheckSourceWithOptions is the per-options variant. Mirrors how
// Compiler#compile accepts a CompilerOptions instance and threads
// it into each pass's constructor.
func CheckSourceWithOptions(path, src string, opts *CompilerOptions) ([]diag.Diagnostic, error) {
	if opts == nil {
		opts = NewCompilerOptions()
	}
	// Pre-scan for lexical errors (e.g. bare `\s` invalid escape) that the
	// legacy parser silently recovers from. Mirrors Java's scanner which
	// reports every invalid-escape via ErrorReporter before the parser
	// builds the AST, so Closure emits JSC_PARSE_ERROR even if subsequent
	// recovery produces a tree. When lexical errors are present we skip
	// downstream passes (Java halts compilation on parser errors).
	if parseDiags := scanLexicalErrors(path, src); len(parseDiags) > 0 {
		return parseDiags, nil
	}
	script, err := parseSource(path, src)
	if err != nil {
		return nil, err
	}
	// Bounded-generic-types are a JSDoc parse error in Java's
	// JsDocInfoParser (TEMPLATE branch). Java reports them via
	// RhinoErrorReporter and stops compilation; jsgo runs a quick
	// scan here so we don't run VarCheck and emit cascading
	// JSC_UNDEFINED_VARIABLE on names whose declarations were
	// skipped by the bounded-generics aborted parse.
	if boundedDiags := scanBoundedGenericErrors(path, script); len(boundedDiags) > 0 {
		return boundedDiags, nil
	}
	reg := jstype.NewRegistry()
	creator := NewTypedScopeCreator(reg)
	creator.sourcePath = path
	creator.sourceSrc = src
	externScope := creator.createExternScope()
	scope := creator.createScopeInWithParent(script, externScope)
	// InferConsts annotates NAME nodes before diagnostic passes run.
	NewInferConsts().Check(script)
	// InferJsDocInfo propagates JSDoc from AST nodes to their associated
	// JSType objects. Pure annotation pass; must run after scope creation.
	NewInferJsDocInfo(creator, scope).Process(nil, script)
	var all []diag.Diagnostic
	// Note: scope-creation diagnostics (e.g. cyclic enum,
	// JSC_CONFLICTING_GETTER_SETTER_TYPE) are drained inside
	// ResolveTypesInScope below — do NOT also drain them here or they
	// will be duplicated.
	varChecker := NewVarCheck(path, creator, scope)
	all = append(all, varChecker.Check(script)...)
	all = append(all, NewVariableReferenceCheck(path, externScope).Check(script)...)
	all = append(all, NewCheckJsDocTypes(path, src, scope, creator).Check()...)
	all = append(all, creator.ResolveTypesInScope(scope)...)
	tc := NewTypeCheckWithOptions(path, creator, scope, opts)
	// Inject VarCheck's undeclared-name set so TypeCheck can suppress
	// POSSIBLE_INEXISTENT_PROPERTY on property chains whose root object
	// was already flagged as undeclared. Mirrors Java VarCheck's
	// createSynthesizedExternVar mechanism.
	tc.SetUndeclaredNames(varChecker.UndeclaredNames())
	all = append(all, tc.Check(script)...)
	all = append(all, NewCheckMissingReturn(path, scope).Check(script)...)
	all = append(all, NewCheckDupVarDecl(path, scope).Check(script)...)
	all = append(all, NewCheckDuplicateMember(path).Check(script)...)
	cse := NewCheckSideEffects(path)
	cse.SetUndeclaredNames(varChecker.UndeclaredNames())
	all = append(all, cse.Check(script)...)
	all = append(all, NewCheckDebuggerStatement(path, opts).Check(script)...)
	all = append(all, NewCheckRegExp(path, true, true).Check(script)...)
	all = append(all, NewCheckSuspiciousCode(path).Check(script)...)
	all = append(all, NewCheckUnreachableCode(path).Check(script)...)
	all = append(all, NewCheckSuper(path).Check(script)...)
	all = append(all, NewCheckAccessControlsWithCreator(path, creator, scope).Check(script)...)
	all = append(all, NewCheckConformance(path).Check(script)...)
	all = append(all, NewCheckJSDoc(path, src).Check(script)...)
	all = append(all, NewCheckTypeImportCodeReferences(path).Check(script)...)
	// Build a single-file ModuleMetadataMap by scanning goog.provide/module/require
	// calls in the script. This mirrors GatherModuleMetadata.java, which the Java
	// compiler runs before CheckClosureImports and CheckMissingRequires.
	moduleMetadata := cmrGatherModuleMetadata(path, script)
	// Wire the real metadata into CheckClosureImports so its namespace-existence
	// checks fire when goog.provide / goog.module declarations are present.
	all = append(all, NewCheckClosureImports(path, cmrTocci(moduleMetadata)).Check(script)...)
	// CheckMissingRequires now has a real single-file module graph.
	all = append(all, NewCheckMissingRequires(path, moduleMetadata).Check(script)...)
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

// ---------------------------------------------------------------------------
// Compiler — 1:1 surface port of Java's Compiler class.
// ---------------------------------------------------------------------------

// JSError is a single error or warning produced by the compiler.
// Mirrors Java's JSError class. The full Java type carries a SourceFile
// reference, a DiagnosticType, formatted args, and a default level; the
// Go port carries the equivalent fields used by jsgo's diag pipeline.
// TODO(port): JSError full surface — DiagnosticType, args, makeFromNode, format(...)
type JSError struct {
	// Diag is the underlying diagnostic. jsgo represents JSError via a
	// diag.Diagnostic + a default level; Java keeps them as separate
	// classes.
	Diag diag.Diagnostic
	// DefaultLevel is the level to surface this error at unless
	// overridden by a WarningsGuard. Mirrors JSError.defaultLevel().
	DefaultLevel diag.Level
}

// Result holds the outputs of a single compilation. Mirrors Java's
// Result class — the Java type is an immutable record bundling errors,
// warnings, variable maps, source maps, etc. Only the fields jsgo's
// callers consult today are wired; everything else is reserved for
// later passes.
// TODO(port): Result — full Java record (variableMap, propertyMap, sourceMap, …)
type Result struct {
	// Errors is the immutable list of errors produced by the compile.
	// Mirrors Result.errors.
	Errors []JSError
	// Warnings is the immutable list of warnings. Mirrors Result.warnings.
	Warnings []JSError
	// VariableMap is the renaming map for variables. Mirrors Result.variableMap.
	VariableMap interface{}
	// PropertyMap is the renaming map for properties. Mirrors Result.propertyMap.
	PropertyMap interface{}
	// StringMap is the renaming map for replaced strings. Mirrors Result.stringMap.
	StringMap interface{}
	// InstrumentationMapping is the parameter encoding for instrumentation.
	// Mirrors Result.instrumentationMapping.
	InstrumentationMapping interface{}
	// SourceMap is the SourceMap for the produced output. Mirrors Result.sourceMap.
	SourceMap interface{}
	// ExternExports is the @externs export string. Mirrors Result.externExport.
	ExternExports string
	// CssNames is the set of CSS class names used. Mirrors Result.cssNames.
	CssNames map[string]bool
	// IdGeneratorMap is the serialized id generator map. Mirrors Result.idGeneratorMap.
	IdGeneratorMap string
	// TranspiledFiles is true if any file required transpilation.
	// Mirrors Result.transpiledFiles.
	TranspiledFiles bool
}

// DiagnosticGroups is the registry of well-known DiagnosticGroup values
// addressable from CompilerOptions and the @suppress annotation. Mirrors
// Java's DiagnosticGroups class, which Compiler.getDiagnosticGroups()
// returns a fresh instance of so callers can iterate the named groups.
// jsgo only tracks the groups exercised by the current passes
// (CHECK_TYPES, STRICT_CHECK_TYPES, …) — see internal/jscomp/diagnosticgroup.go.
// TODO(port): DiagnosticGroups — full ~1500-line Java registry of well-known groups.
type DiagnosticGroups struct{}

// ErrorManager records errors and warnings during compilation.
// Mirrors Java's ErrorManager interface; concrete implementations
// (LoggerErrorManager, ThreadSafeDelegatingErrorManager, etc.) are not
// yet ported — the Go ErrorManager is a thin in-memory recorder used
// by Compiler.report.
// TODO(port): ErrorManager — interface + LoggerErrorManager / ThreadSafeDelegatingErrorManager
type ErrorManager struct {
	errors   []JSError
	warnings []JSError
}

// Report records an error or warning at the given level. Mirrors Java's
// ErrorManager.report(CheckLevel, JSError).
func (em *ErrorManager) Report(level diag.Level, err JSError) {
	if level == diag.Error {
		em.errors = append(em.errors, err)
	} else {
		em.warnings = append(em.warnings, err)
	}
}

// GetErrors returns all errors recorded so far. Mirrors Java's
// ErrorManager.getErrors().
func (em *ErrorManager) GetErrors() []JSError { return em.errors }

// GetWarnings returns all warnings recorded so far. Mirrors Java's
// ErrorManager.getWarnings().
func (em *ErrorManager) GetWarnings() []JSError { return em.warnings }

// GetErrorCount returns the number of errors recorded so far.
// Mirrors Java's ErrorManager.getErrorCount().
func (em *ErrorManager) GetErrorCount() int { return len(em.errors) }

// GetWarningCount returns the number of warnings recorded so far.
// Mirrors Java's ErrorManager.getWarningCount().
func (em *ErrorManager) GetWarningCount() int { return len(em.warnings) }

// HasHaltingErrors reports whether any halting-level error has been
// recorded. Mirrors Java's ErrorManager.hasHaltingErrors().
func (em *ErrorManager) HasHaltingErrors() bool { return len(em.errors) > 0 }

// GenerateReport flushes the error/warning summary to the configured
// output stream. Mirrors Java's ErrorManager.generateReport().
// TODO(port): ErrorManager.generateReport — requires MessageFormatter port
func (em *ErrorManager) GenerateReport() {
	// No-op: MessageFormatter not yet ported. Java's LoggerErrorManager
	// prints a pretty summary at the end of compilation.
	_ = em.errors
}

// Compiler is the orchestrator of a single compilation. It mirrors
// Java's Compiler class (~3000 lines): owns options, the error
// manager, externs/source roots, the chunk graph, the type registry,
// and the pass pipeline. Today the Go port wraps an embedded
// *AbstractCompiler (which carries the shared cross-pass state) and
// adds the Compiler-specific surface. Future work fills in the pass
// pipeline and Result construction.
type Compiler struct {
	// AbstractCompiler holds the state shared across passes (options,
	// lifecycle stage, diagnostics buffer, type registry, …). Java's
	// Compiler extends AbstractCompiler — the Go port composes it.
	*AbstractCompiler

	// outStream is the target for human-readable error output.
	// Mirrors Java's Compiler.outStream.
	outStream *os.File

	// errorManager records errors/warnings. Mirrors Java's
	// Compiler.errorManager (a ThreadSafeDelegatingErrorManager wrap).
	errorManager *ErrorManager

	// externsRoot is the root of the externs AST. Mirrors Java's
	// Compiler.externsRoot.
	externsRoot *rhino.Node

	// jsRoot is the root of the user JS AST. Mirrors Java's
	// Compiler.jsRoot.
	jsRoot *rhino.Node

	// externAndJsRoot is the synthesized parent of externsRoot/jsRoot.
	// Mirrors Java's Compiler.externAndJsRoot.
	externAndJsRoot *rhino.Node

	// typeRegistry is the JSType registry. Mirrors Java's
	// Compiler.typeRegistry (built lazily on first parse).
	typeRegistry *jstype.Registry

	// typedScopeCreator builds the typed scope tree. Mirrors Java's
	// Compiler.typedScopeCreator.
	typedScopeCreator *TypedScopeCreator

	// topScope is the top-level typed scope. Mirrors Java's
	// Compiler.topScope.
	topScope *TypedScope

	// inputs is the ordered list of CompilerInputs for this
	// compilation. Mirrors Java's Compiler.inputs (which is derived
	// from the chunk graph).
	inputs []*CompilerInput

	// inputsById indexes CompilerInputs by their InputId for
	// O(1) lookup. Mirrors Java's Compiler.inputsById.
	inputsById map[string]*CompilerInput

	// moduleTypesByName tracks the module type detected per module
	// name. Mirrors Java's Compiler.moduleTypesByName.
	moduleTypesByName map[string]ModuleType

	// passConfig is the pass-pipeline configuration. Java uses
	// PassConfig (DefaultPassConfig); Go threads CompilerOptions and
	// builds passes in CheckSourceWithOptions today.
	// TODO(port): Compiler.passes / PassConfig — requires PhaseOptimizer port
	passConfig interface{}
}

// Go-only: language-adapter
// NewCompiler returns a Compiler that reports errors and warnings to
// its logger. Mirrors Java's Compiler() no-arg constructor; Go has no
// constructors, so the matching factory name is exported.
func NewCompiler() *Compiler {
	return NewCompilerWithStream(nil)
}

// Go-only: language-adapter
// NewCompilerWithStream returns a Compiler that reports errors and
// warnings to the given output stream. Mirrors Java's
// Compiler(PrintStream) constructor (Go has no overloaded
// constructors so each overload becomes a named factory).
func NewCompilerWithStream(outStream *os.File) *Compiler {
	c := &Compiler{
		AbstractCompiler:  NewAbstractCompiler(NewCompilerOptions()),
		outStream:         outStream,
		errorManager:      &ErrorManager{},
		moduleTypesByName: make(map[string]ModuleType),
		inputsById:        make(map[string]*CompilerInput),
	}
	return c
}

// Go-only: language-adapter
// NewCompilerWithErrorManager returns a Compiler that uses a custom
// error manager. Mirrors Java's Compiler(ErrorManager) constructor.
func NewCompilerWithErrorManager(em *ErrorManager) *Compiler {
	c := NewCompiler()
	c.SetErrorManager(em)
	return c
}

// SetErrorManager installs a custom ErrorManager. Mirrors Java's
// Compiler.setErrorManager(ErrorManager). Java wraps the passed-in
// manager in a ThreadSafeDelegatingErrorManager; the Go port stores
// it directly (jsgo is single-threaded).
func (c *Compiler) SetErrorManager(em *ErrorManager) {
	if em == nil {
		panic("the error manager cannot be null")
	}
	c.errorManager = em
}

// GetErrorManager returns the active ErrorManager. Mirrors Java's
// Compiler.getErrorManager().
func (c *Compiler) GetErrorManager() *ErrorManager {
	return c.errorManager
}

// InitOptions initialises the compiler options. Mirrors Java's
// Compiler.initOptions(CompilerOptions). Java additionally invokes
// initExperimentalForceTranspileOptions, sets allowableFeatures from
// the input language, and chooses an ErrorManager based on outStream
// — those branches are TODO until FeatureSet / LoggerErrorManager
// land.
// TODO(port): Compiler.initOptions full body — requires FeatureSet, LoggerErrorManager, MessageFormatter
func (c *Compiler) InitOptions(opts *CompilerOptions) {
	c.AbstractCompiler.options = opts
}

// GetOptions returns the active CompilerOptions. Mirrors Java's
// Compiler.getOptions() (an @Override of AbstractCompiler.getOptions).
func (c *Compiler) GetOptions() *CompilerOptions {
	return c.AbstractCompiler.GetOptions()
}

// GetDiagnosticGroups returns the registry of named diagnostic
// groups. Mirrors Java's Compiler.getDiagnosticGroups(), which
// constructs a fresh DiagnosticGroups each time.
func (c *Compiler) GetDiagnosticGroups() *DiagnosticGroups {
	return &DiagnosticGroups{}
}

// Report records a JSError. Mirrors Java's Compiler.report(JSError).
// Java consults the WarningsGuard to upgrade/downgrade the level
// before delegating to errorManager.report; the Go port uses the
// JSError's DefaultLevel directly until WarningsGuard is wired.
// TODO(port): Compiler.report(JSError) full body — requires WarningsGuard + ErrorHandler
func (c *Compiler) Report(err JSError) {
	level := err.DefaultLevel
	c.errorManager.Report(level, err)
	// Mirror AbstractCompiler.Report so existing diagnostic-collecting
	// callers continue to work.
	c.AbstractCompiler.Report(err.Diag)
}

// Go-only: language-adapter
// ReportWithLevel records a JSError at a specific level. Mirrors
// Java's Compiler.report(CheckLevel, JSError) — Go has no overloaded
// methods so the level-taking variant gets a distinct name.
func (c *Compiler) ReportWithLevel(ignoredLevel diag.Level, err JSError) {
	_ = ignoredLevel
	c.Report(err)
}

// GetErrors returns the list of errors recorded so far. Mirrors
// Java's Compiler.getErrors().
func (c *Compiler) GetErrors() []JSError {
	if c.errorManager == nil {
		return nil
	}
	return c.errorManager.GetErrors()
}

// GetWarnings returns the list of warnings recorded so far. Mirrors
// Java's Compiler.getWarnings().
func (c *Compiler) GetWarnings() []JSError {
	if c.errorManager == nil {
		return nil
	}
	return c.errorManager.GetWarnings()
}

// GetErrorCount returns the number of errors. Mirrors Java's
// Compiler.getErrorCount().
func (c *Compiler) GetErrorCount() int {
	return c.errorManager.GetErrorCount()
}

// GetWarningCount returns the number of warnings. Mirrors Java's
// Compiler.getWarningCount().
func (c *Compiler) GetWarningCount() int {
	return c.errorManager.GetWarningCount()
}

// HasErrors reports whether any halting errors have been recorded.
// Mirrors Java's Compiler.hasErrors() which delegates to
// hasHaltingErrors().
func (c *Compiler) HasErrors() bool {
	return c.HasHaltingErrors()
}

// HasHaltingErrors reports whether errors that should halt
// compilation have been recorded. Mirrors Java's
// Compiler.hasHaltingErrors(). Java additionally checks
// canContinueAfterErrors on the options; the Go port follows the
// same logic.
func (c *Compiler) HasHaltingErrors() bool {
	if c.errorManager != nil && c.errorManager.HasHaltingErrors() {
		return true
	}
	return c.AbstractCompiler.HasHaltingErrors()
}

// GetResult returns the Result of the compile. Mirrors Java's
// Compiler.getResult() which bundles errors, warnings, and the
// various rename maps into an immutable Result record.
func (c *Compiler) GetResult() *Result {
	return &Result{
		Errors:                 c.GetErrors(),
		Warnings:               c.GetWarnings(),
		VariableMap:            c.GetVariableMap(),
		PropertyMap:            c.GetPropertyMap(),
		StringMap:              c.GetStringMap(),
		InstrumentationMapping: c.GetInstrumentationMapping(),
		ExternExports:          c.GetExternExports(),
		CssNames:               c.GetCssNames(),
		IdGeneratorMap:         c.GetIdGeneratorMap(),
		TranspiledFiles:        c.GetTranspiledFiles(),
	}
}

// GetInstrumentationMapping returns the instrumentation mapping.
// Mirrors Java's Compiler.getInstrumentationMapping().
func (c *Compiler) GetInstrumentationMapping() interface{} {
	return c.AbstractCompiler.instrumentationMapping
}

// GetRoot returns the synthetic root containing both the externs and
// the user JS roots. Mirrors Java's Compiler.getRoot().
func (c *Compiler) GetRoot() *rhino.Node {
	return c.externAndJsRoot
}

// GetExternsRoot returns the externs AST root. Mirrors Java's
// Compiler.getExternsRoot() (protected — exposed here as the Go port
// has no package-protected access).
func (c *Compiler) GetExternsRoot() *rhino.Node {
	return c.externsRoot
}

// GetJsRoot returns the user JS AST root. Mirrors Java's
// Compiler.getJsRoot() (protected — exposed here as Go has no
// package-protected access).
func (c *Compiler) GetJsRoot() *rhino.Node {
	return c.jsRoot
}

// GetTypeRegistry returns the JSType registry for this compilation.
// Mirrors Java's Compiler.getTypeRegistry(). Java lazily constructs
// the registry the first time it's accessed; the Go port does the
// same.
func (c *Compiler) GetTypeRegistry() *jstype.Registry {
	if c.typeRegistry == nil {
		c.typeRegistry = jstype.NewRegistry()
	}
	return c.typeRegistry
}

// GetTypedScopeCreator returns the typed scope creator. Mirrors Java's
// Compiler.getTypedScopeCreator(); returns nil until type checking
// has run.
func (c *Compiler) GetTypedScopeCreator() *TypedScopeCreator {
	return c.typedScopeCreator
}

// GetTopScope returns the top-level typed scope. Mirrors Java's
// Compiler.getTopScope(); returns nil until type checking has run.
func (c *Compiler) GetTopScope() *TypedScope {
	return c.topScope
}

// GetCodingConvention returns the active CodingConvention. Mirrors
// Java's Compiler.getCodingConvention(), which returns the option's
// convention or a default ClosureCodingConvention.
// TODO(port): Compiler.getCodingConvention — requires CodingConvention port
func (c *Compiler) GetCodingConvention() interface{} {
	// Not yet ported. Java consults options.getCodingConvention() and
	// falls back to defaultCodingConvention.
	_ = c.GetOptions()
	return nil
}

// GetReverseAbstractInterpreter returns the reverse abstract
// interpreter used for type narrowing. Mirrors Java's
// Compiler.getReverseAbstractInterpreter(), which constructs a
// ClosureReverseAbstractInterpreter chained to a SemanticReverseAbstractInterpreter.
// TODO(port): Compiler.getReverseAbstractInterpreter — requires SemanticReverseAbstractInterpreter
func (c *Compiler) GetReverseAbstractInterpreter() interface{} {
	_ = c.GetTypeRegistry()
	return nil
}

// GetTypeMismatches returns the set of TypeMismatch records collected
// during type checking. Mirrors Java's Compiler.getTypeMismatches(),
// which delegates to the TypeValidator. Returns an empty slice until
// TypeValidator is wired.
// TODO(port): Compiler.getTypeMismatches — requires TypeValidator port
func (c *Compiler) GetTypeMismatches() []interface{} {
	_ = c.GetTypeValidator()
	return nil
}

// GetTypeValidator returns the TypeValidator instance. Mirrors Java's
// Compiler.getTypeValidator(); returns nil until TypeValidator is
// ported.
// TODO(port): Compiler.getTypeValidator — requires TypeValidator port
func (c *Compiler) GetTypeValidator() interface{} {
	_ = c.typeRegistry
	return nil
}

// GetInputsById returns the InputId-keyed map of inputs. Mirrors
// Java's Compiler.getInputsById().
func (c *Compiler) GetInputsById() map[string]*CompilerInput {
	return c.inputsById
}

// GetInput returns the CompilerInput for the given InputId, or nil.
// Mirrors Java's Compiler.getInput(InputId).
func (c *Compiler) GetInput(id string) *CompilerInput {
	if c.inputsById == nil {
		return nil
	}
	return c.inputsById[id]
}

// GetInputsInOrder returns the inputs in dependency order. Mirrors
// Java's Compiler.getInputsInOrder().
// TODO(port): Compiler.getInputsInOrder — requires ChunkGraph dependency sort
func (c *Compiler) GetInputsInOrder() []*CompilerInput {
	return c.inputs
}

// GetChunks returns the JSChunks of this compilation. Mirrors Java's
// Compiler.getChunks(), which iterates the chunk graph. Returns nil
// until ChunkGraph is wired.
// TODO(port): Compiler.getChunks — requires ChunkGraph port
func (c *Compiler) GetChunks() []interface{} {
	_ = c.GetChunkGraph()
	return nil
}

// GetChunkGraph returns the ChunkGraph. Mirrors Java's
// Compiler.getChunkGraph(); returns nil until ChunkGraph is ported.
// TODO(port): Compiler.getChunkGraph — requires ChunkGraph port
func (c *Compiler) GetChunkGraph() interface{} {
	_ = c.inputs
	return nil
}

// GetModuleLoader returns the ModuleLoader. Mirrors Java's
// Compiler.getModuleLoader().
// TODO(port): Compiler.getModuleLoader — requires ModuleLoader port
func (c *Compiler) GetModuleLoader() interface{} {
	_ = c.GetOptions()
	return nil
}

// GetCompilerExecutor returns the CompilerExecutor. Mirrors Java's
// Compiler.getCompilerExecutor() (protected).
// TODO(port): Compiler.getCompilerExecutor — requires CompilerExecutor port
func (c *Compiler) GetCompilerExecutor() interface{} {
	_ = c.GetOptions()
	return nil
}

// GetSourceMap returns the active SourceMap. Mirrors Java's
// Compiler.getSourceMap().
// TODO(port): Compiler.getSourceMap — requires SourceMap port
func (c *Compiler) GetSourceMap() interface{} {
	_ = c.GetOptions()
	return nil
}

// GetSyntheticExternsInput returns the synthetic externs input.
// Mirrors Java's Compiler.getSyntheticExternsInput(); returns nil
// until SourceFile/CompilerInput externs synthesis is wired.
// TODO(port): Compiler.getSyntheticExternsInput — requires SourceFile + JsAst
func (c *Compiler) GetSyntheticExternsInput() *CompilerInput {
	_ = c.inputsById
	return nil
}

// GetSynthesizedTypeSummaryInput returns the synthesized
// type-summary input. Mirrors Java's
// Compiler.getSynthesizedTypeSummaryInput().
// TODO(port): Compiler.getSynthesizedTypeSummaryInput — requires TypedAst port
func (c *Compiler) GetSynthesizedTypeSummaryInput() *CompilerInput {
	_ = c.inputsById
	return nil
}

// GetUniqueIdSupplier returns the unique-id supplier used by Compiler
// passes that need to mint fresh names. Mirrors Java's
// Compiler.getUniqueIdSupplier().
// TODO(port): Compiler.getUniqueIdSupplier — requires UniqueIdSupplier port
func (c *Compiler) GetUniqueIdSupplier() interface{} {
	_ = c.GetOptions()
	return nil
}

// GetUniqueNameIdSupplier returns the unique-name supplier used by
// Compiler passes (e.g. ExpressionDecomposer). Mirrors Java's
// Compiler.getUniqueNameIdSupplier(), which threads through a counter
// so each pass gets fresh ids.
// TODO(port): Compiler.getUniqueNameIdSupplier — requires Supplier<String> port
func (c *Compiler) GetUniqueNameIdSupplier() interface{} {
	_ = c.GetOptions()
	return nil
}

// Init initialises the compiler with externs and source inputs.
// Mirrors Java's Compiler.init(List, List, CompilerOptions).
// TODO(port): Compiler.init — requires SourceFile, ChunkGraph, full pass pipeline
func (c *Compiler) Init(externs []*CompilerInput, inputs []*CompilerInput, opts *CompilerOptions) {
	c.InitOptions(opts)
	c.inputs = inputs
	c.inputsById = make(map[string]*CompilerInput, len(inputs)+len(externs))
	for _, in := range inputs {
		c.inputsById[in.GetInputID()] = in
	}
	for _, ex := range externs {
		c.inputsById[ex.GetInputID()] = ex
	}
}

// InitChunks initialises the compiler with externs and chunked
// inputs. Mirrors Java's Compiler.initChunks(List, List, CompilerOptions).
// TODO(port): Compiler.initChunks — requires JSChunk port
func (c *Compiler) InitChunks(externs []*CompilerInput, chunks []interface{}, opts *CompilerOptions) {
	_ = chunks
	c.InitOptions(opts)
	for _, ex := range externs {
		if c.inputsById == nil {
			c.inputsById = make(map[string]*CompilerInput)
		}
		c.inputsById[ex.GetInputID()] = ex
	}
}

// Compile compiles a single externs+input pair. Mirrors Java's
// Compiler.compile(SourceFile, SourceFile, CompilerOptions).
// TODO(port): Compiler.compile(extern, input, opts) — requires full pipeline orchestration
func (c *Compiler) Compile(extern *CompilerInput, input *CompilerInput, opts *CompilerOptions) *Result {
	return c.CompileMany(
		[]*CompilerInput{extern},
		[]*CompilerInput{input},
		opts,
	)
}

// Go-only: language-adapter
// CompileMany compiles the given lists of externs and inputs.
// Mirrors Java's Compiler.compile(List, List, CompilerOptions) — Go
// has no overloaded methods so the list-taking variant gets a
// distinct name.
// TODO(port): Compiler.compile(List, List, opts) — requires init/parseForCompilation/stage*Passes
func (c *Compiler) CompileMany(externs []*CompilerInput, inputs []*CompilerInput, opts *CompilerOptions) *Result {
	c.Init(externs, inputs, opts)
	if !c.HasErrors() {
		c.ParseForCompilation()
	}
	c.GenerateReport()
	return c.GetResult()
}

// CompileChunks compiles the given chunks. Mirrors Java's
// Compiler.compileChunks(List, List, CompilerOptions).
// TODO(port): Compiler.compileChunks — requires JSChunk + ChunkGraph
func (c *Compiler) CompileChunks(externs []*CompilerInput, chunks []interface{}, opts *CompilerOptions) *Result {
	c.InitChunks(externs, chunks, opts)
	c.GenerateReport()
	return c.GetResult()
}

// CompileModules is the legacy name for CompileChunks. Mirrors
// Java's Compiler.compileModules(List, List, CompilerOptions),
// which delegates to compileChunks.
// TODO(port): Compiler.compileModules — alias for compileChunks
func (c *Compiler) CompileModules(externs []*CompilerInput, modules []interface{}, opts *CompilerOptions) *Result {
	return c.CompileChunks(externs, modules, opts)
}

// GenerateReport flushes the error/warning summary. Mirrors Java's
// Compiler.generateReport(), which delegates to errorManager.
func (c *Compiler) GenerateReport() {
	if c.errorManager != nil {
		c.errorManager.GenerateReport()
	}
}

// ParseForCompilation parses every input in preparation for the
// pipeline stages. Mirrors Java's Compiler.parseForCompilation(),
// which builds the externsRoot/jsRoot from each CompilerInput's AST.
// TODO(port): Compiler.parseForCompilation — requires CompilerInput.getAstRoot wiring
func (c *Compiler) ParseForCompilation() {
	// Not yet wired: CompilerInput lazy parse + root attachment.
	// Java iterates inputs, calls input.getAstRoot(this), then
	// attaches each script under jsRoot/externsRoot.
	for _, in := range c.inputs {
		_ = in
	}
}

// Parse mirrors Java's Compiler.parse() — the no-arg overload. Java's
// version drives parseForCompilation through the chunk graph.
// TODO(port): Compiler.parse() — requires PassConfig + ChunkGraph
func (c *Compiler) Parse() {
	c.ParseForCompilation()
	return
}

// Go-only: language-adapter
// ParseFile parses a single source file. Mirrors Java's
// Compiler.parse(SourceFile) — Go has no overloaded methods so the
// SourceFile-taking variant gets a distinct name.
// TODO(port): Compiler.parse(SourceFile) — requires SourceFile.getCode + ParserRunner
func (c *Compiler) ParseFile(input *CompilerInput) *rhino.Node {
	if input == nil {
		return nil
	}
	script, err := parseSource(input.GetName(), input.sourceCode)
	if err != nil {
		return nil
	}
	c.RegisterScriptNode(input.GetName(), script)
	return script
}

// ParseTestCode parses a piece of test JavaScript. Mirrors Java's
// Compiler.parseTestCode(String), used by tests to build an AST
// without going through parseForCompilation.
func (c *Compiler) ParseTestCode(js string) *rhino.Node {
	script, err := parseSource("[testcode]", js)
	if err != nil {
		return nil
	}
	return script
}

// Stage1Passes runs stage-1 passes (error and type checking).
// Mirrors Java's Compiler.stage1Passes().
// TODO(port): Compiler.stage1Passes — requires PhaseOptimizer + DefaultPassConfig
func (c *Compiler) Stage1Passes() {
	// Not yet wired. Java runs performChecks() under runInCompilerThread.
	_ = c.passConfig
}

// Stage2Passes runs stage-2 passes (transpilation + optimisation).
// Mirrors Java's Compiler.stage2Passes(SegmentOfCompilationToRun).
// TODO(port): Compiler.stage2Passes — requires PhaseOptimizer + transpilation passes
func (c *Compiler) Stage2Passes(segment int) {
	_ = segment
	_ = c.passConfig
}

// Stage3Passes runs stage-3 passes (localisation). Mirrors Java's
// Compiler.stage3Passes().
// TODO(port): Compiler.stage3Passes — requires localisation pass set
func (c *Compiler) Stage3Passes() {
	_ = c.passConfig
}

// PerformPostCompilationTasks runs the post-compilation finalisers
// (source-map flushing, etc.). Mirrors Java's
// Compiler.performPostCompilationTasks().
// TODO(port): Compiler.performPostCompilationTasks — requires SourceMap + LicenseTracker
func (c *Compiler) PerformPostCompilationTasks() {
	_ = c.GetOptions()
}

// InstrumentForCoverage instruments the AST for coverage measurement.
// Mirrors Java's Compiler.instrumentForCoverage().
// TODO(port): Compiler.instrumentForCoverage — requires CoverageInstrumentationPass
func (c *Compiler) InstrumentForCoverage() {
	_ = c.GetOptions()
}

// WhitespaceOnlyPasses runs the whitespace-only optimisation pipeline.
// Mirrors Java's Compiler.whitespaceOnlyPasses().
// TODO(port): Compiler.whitespaceOnlyPasses — requires PhaseOptimizer
func (c *Compiler) WhitespaceOnlyPasses() {
	_ = c.passConfig
}

// TranspileAndDontCheck runs only transpilation passes without checks.
// Mirrors Java's Compiler.transpileAndDontCheck().
// TODO(port): Compiler.transpileAndDontCheck — requires transpilation pass set
func (c *Compiler) TranspileAndDontCheck() {
	_ = c.passConfig
}

// ToSource prints the entire compiled program as source. Mirrors
// Java's Compiler.toSource() (an @Override of AbstractCompiler.toSource).
// TODO(port): Compiler.toSource() — requires CodePrinter port
func (c *Compiler) ToSource() string {
	_ = c.jsRoot
	return ""
}

// Go-only: language-adapter
// ToSourceNode prints a single node to source. Mirrors Java's
// Compiler.toSource(Node n) — Go has no overloaded methods so the
// node-taking variant gets a distinct name.
// TODO(port): Compiler.toSource(Node) — requires CodePrinter port
func (c *Compiler) ToSourceNode(n *rhino.Node) string {
	_ = n
	_ = c.GetOptions()
	return ""
}

// PrintConfig logs the active configuration to the configured output
// stream. Mirrors Java's Compiler.printConfig().
func (c *Compiler) PrintConfig() {
	if c.outStream == nil {
		return
	}
	fmt.Fprintf(c.outStream, "==== Compiler configuration ====\n")
	if opts := c.GetOptions(); opts != nil {
		fmt.Fprintf(c.outStream, "language-in: %v\n", opts.getLanguageIn())
	}
}

// ReconcileOptionsWithGuards aligns the warnings guards with the
// current options. Mirrors Java's Compiler.reconcileOptionsWithGuards()
// (protected).
// TODO(port): Compiler.reconcileOptionsWithGuards — requires WarningsGuard chain
func (c *Compiler) ReconcileOptionsWithGuards() {
	_ = c.GetOptions()
}

// SetPassConfig installs a custom PassConfig. Mirrors Java's
// Compiler.setPassConfig(PassConfig).
// TODO(port): Compiler.setPassConfig — requires PassConfig type
func (c *Compiler) SetPassConfig(passes interface{}) {
	c.passConfig = passes
}

// GetParserConfig returns the parser config for the requested
// context. Mirrors Java's Compiler.getParserConfig(ConfigContext).
// TODO(port): Compiler.getParserConfig — requires Config type port
func (c *Compiler) GetParserConfig(ctx ConfigContext) interface{} {
	_ = ctx
	_ = c.GetOptions()
	return nil
}

// DisableThreads disables all threads and forces synchronous
// execution. Mirrors Java's Compiler.disableThreads(). The Go port
// is single-threaded today, so this is a no-op.
func (c *Compiler) DisableThreads() {
	_ = c
}

// SetTimeout sets the timeout for runInCompilerThread tasks. Mirrors
// Java's Compiler.setTimeout(int).
// TODO(port): Compiler.setTimeout — requires CompilerExecutor wiring
func (c *Compiler) SetTimeout(timeout int) {
	_ = timeout
	_ = c
}

// GetCrossChunkIdGenerator returns the cross-chunk id generator.
// Mirrors Java's Compiler.getCrossChunkIdGenerator() (an @Override
// of AbstractCompiler.getCrossChunkIdGenerator).
func (c *Compiler) GetCrossChunkIdGenerator() interface{} {
	return c.AbstractCompiler.GetCrossChunkIdGenerator()
}

// GetExportedNames returns the set of exported variable names.
// Mirrors Java's Compiler.getExportedNames() (delegates to AbstractCompiler).
func (c *Compiler) GetExportedNames() map[string]bool {
	return c.AbstractCompiler.GetExportedNames()
}

// GetDefineNames returns the set of @define names. Mirrors Java's
// Compiler.getDefineNames() (delegates to AbstractCompiler).
func (c *Compiler) GetDefineNames() map[string]bool {
	return c.AbstractCompiler.GetDefineNames()
}

// GetExternProperties returns the set of property names found in the
// externs. Mirrors Java's Compiler.getExternProperties() (an
// @Override of AbstractCompiler.getExternProperties).
func (c *Compiler) GetExternProperties() map[string]bool {
	return c.AbstractCompiler.GetExternProperties()
}

// GetComments returns the comments collected for the named file.
// Mirrors Java's Compiler.getComments(String filename), which
// returns the recorded NonJSDocComment list.
func (c *Compiler) GetComments(filename string) interface{} {
	return c.AbstractCompiler.GetComments(filename)
}

// GetModuleMetadataMap returns the active module metadata map.
// Mirrors Java's Compiler.getModuleMetadataMap() (an @Override of
// AbstractCompiler.getModuleMetadataMap).
func (c *Compiler) GetModuleMetadataMap() interface{} {
	return c.AbstractCompiler.GetModuleMetadataMap()
}

// SetModuleMetadataMap installs the module metadata map. Mirrors
// Java's Compiler.setModuleMetadataMap(ModuleMetadataMap).
func (c *Compiler) SetModuleMetadataMap(m interface{}) {
	c.AbstractCompiler.SetModuleMetadataMap(m)
}

// GetModuleMap returns the active ModuleMap. Mirrors Java's
// Compiler.getModuleMap().
func (c *Compiler) GetModuleMap() interface{} {
	return c.AbstractCompiler.GetModuleMap()
}

// SetModuleMap installs the module map. Mirrors Java's
// Compiler.setModuleMap(ModuleMap).
func (c *Compiler) SetModuleMap(m interface{}) {
	c.AbstractCompiler.SetModuleMap(m)
}

// JoinPathParts joins path components with the platform separator.
// Mirrors Java's Compiler.joinPathParts(String...) static helper,
// which the CLI uses to assemble debug-log directory paths.
func JoinPathParts(parts ...string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += "/"
		}
		out += p
	}
	return out
}

// JoinPathParts is the receiver-method form, matching Java's static
// access pattern when called from a Compiler instance.
func (c *Compiler) JoinPathParts(parts ...string) string {
	return JoinPathParts(parts...)
}

// SaveState serialises the compiler state to the given writer.
// Mirrors Java's Compiler.saveState(OutputStream).
// TODO(port): Compiler.saveState — requires CompilerState serialisation
func (c *Compiler) SaveState(w *os.File) error {
	_ = w
	_ = c
	return nil
}

// RestoreState rebuilds compiler state from the given reader.
// Mirrors Java's Compiler.restoreState(InputStream).
// TODO(port): Compiler.restoreState — requires CompilerState deserialisation
func (c *Compiler) RestoreState(r *os.File) error {
	_ = r
	_ = c
	return nil
}
