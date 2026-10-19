# Post-colors native-design activation diagnostics

**Untimed and NON-RACE only.** These overlays add diagnostic mutexes and numeric counters. Never merge them into the plain focused race overlay or a measured production binary. All previous overlays/build attempts remain unchanged.

Prepared against `/workspace/jsgo-post-colors-design-integration` at `9d190a2b20d05bc33b203b4370176129ac80dc80`.

The combined overlay is:

```
/workspace/jsgo-results/round2/post-colors-activation-overlay/all-six-generated/overlay.json
```

It combines corrected native liveness, first-scope coallocation, compact inverter, ASCII output, FSIV quotient, and native RUC task/FIFO numeric observations. The existing Array diagnostic work is separate. The generator is root-parameterized and composable:

```
python3 prepare.py --root <checkout> --base-overlay <plain-or-diagnostic-overlay.json> \
  --out <new-directory> --jobs <pinned-jobs.json>
```

It reads effective production source from a base overlay before adding narrow hooks. Original/effective-source hashes, generated-file hashes, root/revision and pinned-job hash are recorded in the generated manifest. Existing production is never edited. Use a new output directory when a prior build is active.

## What is observed

- `tryNativeLiveVariables`: attach a local existing numeric probe only if the caller supplied none; aggregate after the original call/exception cleanup. Report actual tries, successful solver runs, CFG nodes, solver steps, transfers, transfer-cache hits and escape transitions. No observer callback is installed in the algorithm. This counts actual admitted native solver calls, not mere compiler eligibility.
- `declareFirstCoallocated`: increment only after successful original constructor/legality/publication work has assigned the original map.
- Compact inverter: temporarily attach a private probe only when no original probe exists; report actual arena selections, completed cleanup, records/variables/shadow links/name observations and capacity retirements. An overlay-only probe flag suppresses NameOrder allocation for diagnostic probes. Original correctness-test probes retain their original NameOrder behavior. No AST, Var, names or Scope are retained by the collector.
- ASCII output: the dedicated harness sets the existing runner-owned output probe before calling real `DoRun`. Report actual printers, ASCII/printer promotions, UTF-16 units, join builders and post-colors printers. Printer/buffer selection guards remain original.
- RUC: if the base includes `remove_unused_activation.py` instrumentation, the harness resets/takes its original numeric phase records after the real job. Otherwise this field is absent. Task and FIFO experiment flags remain independent.
- FSIV: if the base includes the quotient diagnostic layer, `DoRun` writes its private aggregate to a separate per-job JSON file and the combined report embeds it. Queries are recorded after the function's original candidate/rewrite stage, including cohort worker passes.

No flag, AST, admission condition, pass factory, compiler option, reporter, priority, error path or CLI callback is changed by these hooks. Global collectors carry numeric data only and are enabled by the dedicated test harness. Their synchronization can mask races, which is why race validation must always use the separate plain overlay.

## Exact pinned-workload harness

The default pinned manifest is the existing five-task Array diagnostic manifest. The test reads its exact workspace/CLI arguments/Java artifact hashes. It neither enables experimental flags nor overrides thread counts. Use a corresponding jobs manifest for a different N; the report records the actual args and process GOMAXPROCS separately.

```
JSGO_POST_COLORS_DIAGNOSTIC_JOBS=/workspace/jsgo-results/round2/array-recipe-diagnostic-overlay/jobs.json \
  <diagnostic-test-binary> -test.run '^TestPostColorsDesignDiagnosticPinnedTasks$' -test.v
```

Set the intended factors in the environment before running. Each job uses the stock runner and actual loader/ownership lifecycle, checks JS/stdout/stderr/status/source-map bytes against pinned Java, writes `post-colors-activation.json` in its own workspace, and logs short actual-dispatch totals. A zero count remains visible rather than being treated as performance evidence for an unexecuted implementation. The harness writes no production stdout/stderr diagnostics.

The prepared combined chain was generated source-only as follows: original RUC generator on the frozen root → merge its two replacements into the preserved plain focused overlay → a private copy of the FSIV diagnostic generator → this generator. No Go build or test was executed by this agent.
