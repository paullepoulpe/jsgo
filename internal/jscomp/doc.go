// Package jscomp implements the Closure JavaScript compiler pipeline, passes,
// typechecker, module resolvers, bundler, and AST/type state serialization.
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
//   - refactoring_*     — automated suggested-fix and refactoring passes
//   - serialization_*   — TypedAST and color pool serializers/deserializers
//   - transpile_*       — standalone per-file transpiler API
//   - type_*            — reverse abstract interpreters and type validators
//
// Generic graph data structures and algorithms (com.google.javascript.jscomp.graph)
// live in package [github.com/paullepoulpe/jsgo/internal/graph].
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
