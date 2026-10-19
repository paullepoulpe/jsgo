# Runtime-library table admission audit

Source-only audit, 2026-10-19. Frozen base c7998d45 and first merged fundamental candidate 4bbbb8a9. No builds, tests, profiles or timings executed for this audit. Initial failing logs/overlays are preserved. This identifies invalid negative experiments; it does not claim gains from re-enabling them.

## Confirmed cause and correction

The initial FSIV quotient guard reused scopeFirstAllocationEligible. Its runtimeLibraryTypedAsts != nil exclusion rejects ordinary optimized CLI compiles. Actual diagnostic attempts were typed 1352, bundle 16494 and downlevel-simple 3697, all stock-rejected, admitted zero. The complete SIMPLE CLI fixture failed with the same reason while all 17 full-set state/decision oracles passed. Such timings would compare the original full-set solver with itself.

DefaultPassConfig.getEarlyOptimizationPasses calls addNonTypedAstNormalizationPasses for ordinary non-precompiled-library compilation (default_pass_config.go2655). The latter schedules typesToColors (2955). ConvertTypesToColors.Process always calls InitRuntimeLibraryTypedAsts before normalization: line143 without type checking, line155 with checking. Compiler.InitRuntimeLibraryTypedAsts stores a nonnil map (compiler.go912). No ordinary pipeline clears it. The table contains deferred built-in runtime-library AST suppliers, not proof of a public/injected typed AST filesystem. Its mere presence does not imply any particular private analysis invokes those suppliers.

Quotient correction 3df44712 gives only that analysis a dedicated source/observer proof. It preserves normalized lifecycle, exact compiler/facade, fresh CLI source ownership, no injected typedAstFilesystem, stock warnings/reporters/conventions/observers and debug/step fallbacks. It admits an already-loaded runtime-library table because its transfer/program preparation never reads/invokes that table. Tests retain positive real-CLI admission/query checks and now assert the genuine post-typesToColors state; an unused throwing runtime supplier proves non-use. No other guard was changed.

## Other affected candidates

These are source-proven phase exclusions in the audited revision. Archived exact-build provenance and real-task dispatch should be checked before revising individual historical conclusions.

| Candidate | Guard | Normal optimized CLI implication | Historical screen evidence |
|---|---|---|---|
| Native liveness and later two-word solver | native_live_variables.go24 | All private DAE/coalesce entries happen after typesToColors, so no activation. Does not establish representation performance. | native-live-words-screen-command.json explicitly enables NATIVE_LIVE_VARIABLES for prior_live and words; archived /workspace/jsgo-structural-native-analysis has the same table exclusion. |
| Compact contextual inverter | contextual_inverter_compact.go21 | Finalization after colors always rejects. Transpile-only can differ; do not reject whitespace results categorically. | native-storage-runtime-flags.json enables COMPACT_INVERTER. |
| First Scope/Var/map coallocation | scope_first_allocation.go13, invoked by stock_scope_insertion.go41 | All ordinary IV/RUC private creators, including early IV in early *optimization*, see the loaded table and reject coallocation. | native-storage-runtime-flags.json enables SCOPE_FIRST_ALLOCATION. |
| Native RUC task storage | remove_unused_tasks.go27 | Ordinary early/main/final RUC after colors rejects. | pass-storage-screen-command.json and code-design-screen-command.json include tasks factor. Archived /workspace/jsgo-structural-traversal-next has identical guard. |
| Intrusive RUC FIFO | same nativeTasksSafe path | Also rejects before selecting the intrusive queue; negative task/queue screen does not evaluate this algorithm. | Both screen commands enable REMOVE_UNUSED_TASKS plus REMOVE_UNUSED_INTRUSIVE_FIFO for intrusive factor. |
| ASCII output buffer | native_output_buffer.go158 | Output after optimized compilation rejects. WHITESPACE/transpile-only may activate and remains separately meaningful. Checks-only may have no JS output. | native-storage-runtime-flags.json enables ASCII_OUTPUT_BUFFER. |

The ordinary stock-insertion shortcut is independent of the first-coallocation flag and does not call this runtime-table proof. Its earlier results are not invalidated by this finding. Likewise, the original inverter no-op proof, generic native CFG/local-flow/reference paths, recipes and bitset changes must not be conflated with the named factors above merely because they share a screen.

## Similar guards that are not automatically broken

fresh_template_values.go43 runs during initial typed inference, before colors; its table-absent guard can be satisfied. This also protects the new Array recipe. access_control_property_facts.go72 runs in the checks pass before optimization conversion; its table absence is not inherently contradictory. Both still require actual positive runtime counters to validate a particular workload. CommonJS runs before colors and has an independently audited fresh-CLI proof.

Transpile-only schedules normalize/transpilation/inverter without the ordinary typesToColors optimization conversion (default_pass_config.go1909 onward). Therefore compact inverter/output-buffer activation in downlevel-whitespace cannot be inferred from optimized SIMPLE/ADVANCED. Real dispatch must be measured separately by task and N.

## Required next review and validation

Do not globally weaken scopeFirstAllocationEligible. Each affected candidate must audit whether its algorithm actually touches runtime-library supplier invocation, type/color-dependent callbacks or mutable declaration graphs. Keep injected typedAstFilesystem/custom reporters/conventions/observers and public paths excluded unless independently proved. Coallocation additionally retains its specific read-only declaration-map/lifetime proof.

For each separately corrected factor: a complete stock CLI lifecycle fixture must assert an actual positive target-dispatch count after typesToColors; fallback fixtures must retain custom/injected behavior. Separate untimed diagnostics should count admission and weighted work on the real pinned jobs/N, then ordinary no-diagnostic normal/race/oracle gates must run. Only after those pass should matched binary/profile/time/RSS screens be repeated. Diagnostic synchronization must remain outside race and timed binaries.

The prior negative screens cannot be used to claim these named designs were exhausted. This finding does not demonstrate that they will help; it reopens their evaluation with the required activation evidence.
