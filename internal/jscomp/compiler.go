// Port of com/google/javascript/jscomp/Compiler.java
// @ closure-compiler v20260415 (shim; pipeline wrapper only).
//
// Closure's Compiler is a ~3000-line orchestrator (PhaseOptimizer,
// pass registry, Result, module graph, warning manager, …). This
// file is a ~30-line adapter that wires our three independent
// passes (parse → TypedScopeCreator → TypeCheck) for the CLI +
// harness. A real Compiler port — pass registry, phase ordering,
// Results — lands when we have more than these three passes.
//
// Substantive Java methods not yet ported (requires unported types):
//   TODO(port): Compiler() — constructor
//   TODO(port): Compiler(PrintStream) — constructor with output stream
//   TODO(port): Compiler(ErrorManager) — constructor with error manager
//   TODO(port): Compiler.compile(SourceFile, SourceFile, CompilerOptions)
//   TODO(port): Compiler.compile(List, List, CompilerOptions)
//   TODO(port): Compiler.compileModules(List, List, CompilerOptions)
//   TODO(port): Compiler.compileChunks(List, List, CompilerOptions)
//   TODO(port): Compiler.init(List, List, CompilerOptions)
//   TODO(port): Compiler.initChunks(List, List, CompilerOptions)
//   TODO(port): Compiler.initOptions(CompilerOptions)
//   TODO(port): Compiler.parseInputs()
//   TODO(port): Compiler.parse()
//   TODO(port): Compiler.check()
//   TODO(port): Compiler.optimize()
//   TODO(port): Compiler.stage1Passes()
//   TODO(port): Compiler.stage2Passes()
//   TODO(port): Compiler.toSource()
//   TODO(port): Compiler.toSource(Node)
//   TODO(port): Compiler.report(JSError)
//   TODO(port): Compiler.getInput(InputId)
//   TODO(port): Compiler.getChunkGraph()
//   TODO(port): Compiler.getJsRoot()
//   TODO(port): Compiler.getRoot()
//   TODO(port): Compiler.getErrorManager()
//   TODO(port): Compiler.getCodingConvention()
//   TODO(port): Compiler.getReverseAbstractInterpreter()
//   TODO(port): Compiler.getTypeRegistry()
//   TODO(port): Compiler.getTopScope()
//   TODO(port): Compiler.getTypedScopeCreator()
//   TODO(port): Compiler.getParserConfig(ConfigContext)
//   TODO(port): Compiler.getUniqueIdSupplier()
//   TODO(port): Compiler.getUniqueNameIdSupplier()
//   TODO(port): Compiler.getChangeTracker()
//   TODO(port): Compiler.beforePass(String)
//   TODO(port): Compiler.afterPass(String)
//   TODO(port): Compiler.process(CompilerPass)
//   TODO(port): Compiler.setErrorManager(ErrorManager)
//   TODO(port): Compiler.printConfig()
//   TODO(port): Compiler.reconcileOptionsWithGuards()
//   TODO(port): Compiler.getNodeForCodeInsertion(JSChunk)
//   TODO(port): Compiler.getTypeValidator()
//   TODO(port): Compiler.getTypeMismatches()
//   TODO(port): Compiler.getSyntheticExternsInput()
//   TODO(port): Compiler.getSynthesizedTypeSummaryInput()
//   TODO(port): Compiler.getModuleLoader()
//   TODO(port): Compiler.getRuntimeJsLibManager()
//   TODO(port): Compiler.getTranspilationNamespace()
//   TODO(port): Compiler.initWithTypedAstFilesystem(List, List, CompilerOptions, InputStream)
//   TODO(port): Compiler.initChunksWithTypedAstFilesystem(List, List, CompilerOptions, InputStream)
//   TODO(port): Compiler.getInputsInOrder()
//   TODO(port): Compiler.clearTypedScopeCreator()
//   TODO(port): Compiler.hasScopeChanged(Node)
//   TODO(port): Compiler.hasHaltingErrors()
//   TODO(port): Compiler.addIndexProvider(IndexProvider)
//   TODO(port): Compiler.getIndex(Class)
//   TODO(port): Compiler.getDefaultErrorReporter()
//   TODO(port): Compiler.clearJSTypeRegistry()
//   TODO(port): Compiler.getColorRegistry()
//   TODO(port): Compiler.setColorRegistry(ColorRegistry)
//   TODO(port): Compiler.initRuntimeLibraryTypedAsts(Optional)
//   TODO(port): Compiler.parseTestCode(String)
//   TODO(port): Compiler.createOrReopenLog(Class, String, String...)
//   TODO(port): Compiler.createOrReopenIndexedLog(Class, String, String...)
//   TODO(port): Compiler.createDefaultExpressionDecomposer()
//   TODO(port): Compiler.createExpressionDecomposer(Supplier, ImmutableSet, Scope)

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
