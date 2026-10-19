# RUC task storage after typesToColors

Source-only correction based on `c7998d45`: remove the
`runtimeLibraryTypedAsts != nil` rejection from `RemoveUnusedCode.nativeTasksSafe`.
No other phase's guard changes here. The existing native task and intrusive
FIFO flags remain independent:

- `JSGO_EXPERIMENT_REMOVE_UNUSED_TASKS=1` selects private continuation storage
  only inside an admitted builtin pass.
- `JSGO_EXPERIMENT_REMOVE_UNUSED_INTRUSIVE_FIFO=1` selects intrusive queueing
  inside that storage. By itself it does not activate native tasks.

Stock DefaultPassConfig runs typesToColors before normalized RUC. Both its
typed and untyped branches initialize the runtime-library table, including
ordinary SIMPLE builds. Consequently the previous table guard rejected
normal optimized CLI runs. Representation fixtures prepared directly for RUC
did not expose that false negative. Old negative measurements without positive
dispatch evidence cannot establish the performance of an active experiment.

## Reachable callback audit

Legacy continuation Apply and native applyNativeTask both dispatch the same
TraverseFunction/TraverseNode methods. RUC's reachable compiler operations
obtain the analyzer, coding convention, extern property set and accessor
summary; create syntactic scopes and no-op diagnostic logs; or report changes.
AstAnalyzer is a concrete data holder without a compiler/runtime loader.
PolyfillUsageFinder queries the polyfill symbol table and traverses guarded
uses; it does not inject a library. Removable.Remove implementations rewrite
AST nodes and report scope changes. None invokes InitRuntimeLibraryTypedAsts,
EnsureLibraryInjected, loadResourceContents, or a runtime AST supplier.

An unused private table therefore introduces no callback/reentry boundary for
this phase. It requires no new global provenance marker. The correction keeps
the independent typedAstFilesystem, fresh CLI ownership, exact builtin/public
builder, compiler/convention, custom pass, warning/error/change observer,
redeclaration handler and diagnostic logger fences.

The new normal `TestRemoveUnusedTaskStockColors*` fixtures establish genuine
stock fresh-CLI ownership, execute the actual default typesToColors factory,
normalize, and invoke the actual default RUC factory. They require real native
construction/application counts, compare full pass snapshots across all four
flag combinations, and insert a throwing unused supplier that must remain
uninvoked. Public, outside-ownership, custom-convention and restored typed-AST
paths still require their independent fallback.

## Complete fresh CLI activation overlay

`tools/bench/remove_unused_activation.py` generates diagnostic copies with
`--root MERGED_CHECKOUT --out NEW_DIRECTORY`. It does not edit source, build,
run Go, download inputs or relax admission. It attaches existing private task
probes only in the overlaid builtin factory and adds coordinator-only records
after Process. Never include this overlay in a timed binary.

```sh
python3 -I tools/bench/remove_unused_activation.py \
  --root /workspace/INTEGRATION --out /tmp/ruc-activation
```

Add its `-overlay=/tmp/ruc-activation/overlay.json` to the next centrally
scheduled `go test -c` build. Overlay replacement paths are disjoint from the
CommonJS activation overlay and can be merged into the same build.
Run the following generated test alone:

```sh
/path/to/diagnostic-jscomp.test \
  -test.run='^TestRemoveUnusedTaskActualCLIPostColorsActivation$' -test.v
```

It calls actual fresh stock CommandLineRunner.DoRun for SIMPLE and ADVANCED,
across tasks-off, native tasks, native tasks plus intrusive FIFO, and FIFO-only.
Every native-enabled mode must show actual post-colors factory activation,
continuation construction and application; every mode must reproduce the
original exit status, output, source map and diagnostics. This would fail with
the old blanket table guard even if output alone matched.

The generated `TestRemoveUnusedTaskBenchmarkActivation` can also consume the
five fixed task JSON cases from `commonjs_activation.py` without changing their
arguments or manufacturing pass execution:

```sh
JSGO_EXPERIMENT_REMOVE_UNUSED_TASKS=1 \
JSGO_EXPERIMENT_REMOVE_UNUSED_INTRUSIVE_FIFO=1 \
JSGO_RUC_ACTIVATION_CASE=/path/to/cases/bundle/threads-2/case-0.json \
  /path/to/diagnostic-jscomp.test \
  -test.run='^TestRemoveUnusedTaskBenchmarkActivation$' -test.v
```

The result goes to the case's report path plus `.ruc.json`, containing the
actual CLI exit/diagnostics and one post-colors/probe record per default RUC
invocation. Use a distinct case/report directory per flag configuration.
No invocations or zero Activations establishes no evidence for an active
optimization. There are no new timing, RSS or performance conclusions here.

Normal tests and these actual-CLI overlay tests have not been built or run by
the source owner. Only source review, generator execution and Go syntax
formatting were performed; central normal/race/artifact gates and same-binary
default-GC factor measurements remain required.
