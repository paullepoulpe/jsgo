// Package jscomp implements the Closure JavaScript compiler pipeline, passes,
// typechecker, module resolvers, and AST/type state serialization.
//
// # Architecture & Compilation Stages
//
// [Compiler] owns the compilation state (inputs, externs, chunks, symbol/type
// tables, diagnostics, and the [rhino.NodeStore] AST arena) configured by
// [CompilerOptions]. [DefaultPassConfig] defines the pass schedule executed by
// [PhaseOptimizer] across five stages:
//
//  1. INITIALIZATIONS — input parsing, dependency/chunk sorting, synthetic
//     externs, and module metadata gathering.
//  2. CHECKS — syntactic checks, ES/CommonJS/goog.module rewriting,
//     [TypedScopeCreator], [TypeInferencePass], [TypeCheck], conformance, and
//     lint checks.
//  3. TRANSPILATIONS — feature-driven downleveling (classes, generators,
//     async/await, destructuring, optional chaining, polyfill injection).
//  4. OPTIMIZATIONS — color conversion, property disambiguation/inlining,
//     fixed-point peephole/inlining/dead-code loops, cross-chunk code motion,
//     variable coalescing, and variable/property renaming.
//  5. FINALIZATIONS — denormalization, [CodePrinter] / [CodeGenerator] output,
//     and V3 source map emission.
//
// # File Naming Conventions
//
// Files porting classes from Java subpackages of com.google.javascript.jscomp
// carry subpackage prefixes so they group naturally in directory listings:
//   - closure_unaware_* — @closureUnaware shadow-AST transpilation and optimization
//   - deps_*            — dependency scanners, goog.provide/require, module resolvers
//   - disambiguate_*    — type/color-based property disambiguation and ambiguation
//   - es6_*             — ES2015+ transpilation passes
//   - ijs_*             — .i.js type-summary generation and checking
//   - instrumentation_* — branch/statement coverage instrumentation
//   - j2cl_*            — J2CL-specific optimization passes
//   - lint_check_*      — optional style and lint checks
//   - modules_*         — ES module and goog.module processors and metadata
//   - peephole_*        — local AST peephole optimizations
//   - polymer_*         — Polymer component rewriting
//   - proto_*           — TypedAST, conformance, and state protobuf codecs
//   - serialization_*   — TypedAST and color pool serializers/deserializers
//   - type_*            — reverse abstract interpreters and type validators
//
// Higher-level tools and standalone subpackages live in their own Go packages:
//   - [github.com/paullepoulpe/jsgo/internal/graph]       — generic graph data structures and algorithms
//   - [github.com/paullepoulpe/jsgo/internal/refactoring] — AST pattern matching, refactoring drivers, and SuggestedFix
//   - [github.com/paullepoulpe/jsgo/internal/linter]      — standalone single-file Linter, LinterMain, and LintPassConfig
//   - [github.com/paullepoulpe/jsgo/internal/bundle]      — per-source transformers, caching, and coverage instrumentation
//   - [github.com/paullepoulpe/jsgo/internal/transpile]   — standalone per-file ES transpiler API and ClosureBundler
//   - [github.com/paullepoulpe/jsgo/internal/cli]         — CommandLineRunner CLI front end, flags, and output writers
//
// # Determinism & AST Conventions
//
// Output and diagnostics must be 100% deterministic and byte-identical to
// Closure Compiler regardless of --num_parallel_threads:
//   - Use [jutil.LinkedHashMap] / [jutil.LinkedHashSet] whenever iteration
//     order can affect AST mutations, symbol IDs, variable/property renaming,
//     or diagnostic ordering. Every Go map range in non-test files is checked
//     by tools/maprangecheck and must carry an // order-independent: comment.
//   - Prefer the typed AST views in [github.com/paullepoulpe/jsgo/internal/ast]
//     over positional [rhino.Node] child navigation (FirstChild/Next/ChildAtIndex)
//     when writing or refactoring passes.
package jscomp
