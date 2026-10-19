# Untimed Array recipe dispatch diagnostics

`prepare.py` creates an overlay and a focused Go test. It never changes a checkout,
builds a compiler, downloads inputs, forces an eligibility guard, or sets an
experimental flag. The frozen production candidate remains `cbd3d91c`.

Use the next merged source checkout and, preferably, its already-prepared focused
factor overlay. This lets the scheduled focused test binary carry the diagnostic
harness; no extra giant CLI build is needed:

```sh
python3 /workspace/jsgo-results/round2/array-recipe-diagnostic-overlay/prepare.py \
  --root /workspace/CHOSEN-MERGED-CHECKOUT \
  --base-overlay /workspace/CHOSEN-FOCUSED-OVERLAY/overlay.json
```

The base overlay must already target the chosen root and must not replace the
two production files that this diagnostic hooks. Without `--base-overlay`, the
generator copies the existing focused-fresh-template preparation into its own
`base/` directory and adds the recipe fixture. `--keep-test FILE` can retain
additional merged factor fixtures in that fallback preparation. It does not
mutate the original focused overlay directory.

The generated `manifest.json` records the root commit, checksums, exact production
changes and pinned reference. `jobs.json` contains the exact five baseline task
argument lists, separate artifact workspaces and pinned Java artifact digests.
Workspace creation only installs symlinks to existing pinned cache inputs and
the selected checkout's benchmark app; it runs no install/fetch/build command.

Build ONLY in a root-scheduled slot, using the generated overlay in place of the
base focused overlay:

```sh
env -u GOGC -u GOMEMLIMIT -u GODEBUG GOMAXPROCS=1 \
  GOCACHE=/workspace/toolchains/cache GOPATH=/workspace/toolchains/gopath GOTOOLCHAIN=local \
  /workspace/toolchains/go1.27.2/bin/go test -race -p 1 \
  -overlay /workspace/jsgo-results/round2/array-recipe-diagnostic-overlay/overlay.json \
  -c -o /workspace/jsgo-results/round2/bin/jscomp-array-recipe-diagnostic.test ./internal/jscomp
```

This is an UNTIMED diagnostic binary. Never use it for elapsed-time, CPU or RSS
comparisons. Build the ordinary production CLI separately without this overlay.
The dedicated diagnostic test skips unless its manifest environment is supplied,
so the same binary can run the normal focused factor correctness fixtures first.

Run it only when the root releases CPU for untimed full tasks. Keep the scheduled
candidate's existing baseline experiment environment and add the independent
recipe flag; leave temporary vectors off. The two explicit prerequisites are:

```sh
JSGO_EXPERIMENT_CLI_SOURCE_LOCK=1 \
JSGO_EXPERIMENT_ARRAY_UNKNOWN_RECIPE=1 \
JSGO_EXPERIMENT_FRESH_TEMPLATE_VALUES= \
JSGO_ARRAY_RECIPE_DIAGNOSTIC_JOBS=/workspace/jsgo-results/round2/array-recipe-diagnostic-overlay/jobs.json \
  /workspace/jsgo-results/round2/bin/jscomp-array-recipe-diagnostic.test \
  -test.run '^TestArrayRecipeDiagnosticPinnedTasks$' -test.v -test.timeout=15m \
  > /workspace/jsgo-results/round2/array-recipe-dispatch.log 2>&1
```

A single task can be selected as a fresh test process, for example:
`-test.run '^TestArrayRecipeDiagnosticPinnedTasks$/^typed-app$'`. Every invocation
constructs a new real `CommandLineRunner` and calls its actual `DoRun` with the
pinned arguments. It does not call EnableTestMode, synthesize source-loader
ownership, inject roots/types, bypass private entry points, or alter compiler
thread flags. Each compiler's ordinary default pipeline chooses whether inference
and private recipe admission run. A zero-pass/zero-hit row is meaningful (e.g. a
task that never performs type inference), and is not forced into activation.

The overlay attaches a caller-owned numeric probe only when the dedicated test
has installed its collector and no original probe is present. Full admission
and cheap V2 checks remain intact; the legacy guard bit stays false. It omits
the probe's per-scope AST observation append so large diagnostic tasks retain
only numeric totals. The collector runs after the private inference mode and
creator window have been restored, including on panic. It does not register a
compiler observer/reporter or alter warning guards. A mutex protects the
optional diagnostic-only collector pointer; each CLI's totals are test-owned.

The test verifies raw status/stdout/stderr/emitted JS/source-map hashes against
the pinned Java reference. Numeric diagnostic reports go through `t.Log`, outside
the captured CLI streams, and therefore do not change the compatibility gate.
No output normalization is applied.

Each task logs one `ARRAY_RECIPE_DIAGNOSTIC` JSON object with:

- invocation/pass and admitted-window counts;
- constructor calls, original recipe seeds, fresh coallocated hits and native
  mutation fallbacks;
- full/fast scope and boundary checks;
- the actual recipe/vector/source-ownership environment flags.

These are activation/coverage diagnostics, not a performance result. No tracing
is added to the timed CLI source or binary. The generator and emitted diagnostic
Go sources were prepared and formatted but have not been compiled or executed
yet; root owns that scheduled validation.
